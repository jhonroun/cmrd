package grpcapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jhonroun/cmrd/internal/grpcapi/pb"
	"github.com/jhonroun/cmrd/pkg/cmrd"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serviceClient interface {
	Resolve(ctx context.Context, links []string) ([]cmrd.FileTask, error)
	Download(ctx context.Context, links []string, onProgress cmrd.ProgressHandler) error
}

type jobState struct {
	JobID    string
	Phase    string
	Percent  float64
	Message  string
	Done     bool
	ErrText  string
	Cancel   context.CancelFunc
	Started  time.Time
	Finished time.Time
}

// Server implements CMRD gRPC service.
type Server struct {
	pb.UnimplementedCMRDServiceServer

	baseConfig    cmrd.Config
	clientFactory func(cmrd.Config) (serviceClient, error)

	mu   sync.RWMutex
	jobs map[string]*jobState
	subs map[string]map[uint64]chan jobState
}

var (
	jobCounter atomic.Uint64
	subCounter atomic.Uint64
)

// NewServer creates gRPC service instance.
func NewServer(cfg cmrd.Config) *Server {
	return NewServerWithFactory(cfg, func(config cmrd.Config) (serviceClient, error) {
		return cmrd.New(config)
	})
}

// NewServerWithFactory creates gRPC server with custom client factory.
func NewServerWithFactory(cfg cmrd.Config, factory func(cmrd.Config) (serviceClient, error)) *Server {
	return &Server{
		baseConfig:    cfg,
		clientFactory: factory,
		jobs:          make(map[string]*jobState),
		subs:          make(map[string]map[uint64]chan jobState),
	}
}

func (s *Server) ResolveLinks(ctx context.Context, req *pb.ResolveLinksRequest) (*pb.ResolveLinksResponse, error) {
	if req == nil || len(req.Links) == 0 {
		return nil, status.Error(codes.InvalidArgument, "links are required")
	}

	client, err := s.clientFactory(s.baseConfig)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create client: %v", err)
	}

	files, err := client.Resolve(ctx, req.Links)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resolve: %v", err)
	}

	response := &pb.ResolveLinksResponse{
		Files: make([]*pb.ResolvedFile, 0, len(files)),
	}
	for _, file := range files {
		response.Files = append(response.Files, &pb.ResolvedFile{
			URL:    file.URL,
			Output: file.Output,
		})
	}
	return response, nil
}

func (s *Server) StartDownload(_ context.Context, req *pb.StartDownloadRequest) (*pb.StartDownloadResponse, error) {
	if req == nil || len(req.Links) == 0 {
		return nil, status.Error(codes.InvalidArgument, "links are required")
	}

	cfg := s.baseConfig
	if strings.TrimSpace(req.DownloadDir) != "" {
		cfg.DownloadDir = strings.TrimSpace(req.DownloadDir)
	}

	client, err := s.clientFactory(cfg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create client: %v", err)
	}

	jobID := nextJobID()
	jobCtx, cancel := context.WithCancel(context.Background())
	s.setJob(&jobState{
		JobID:   jobID,
		Phase:   "created",
		Message: "job created",
		Cancel:  cancel,
		Started: time.Now(),
	})

	go func() {
		err := client.Download(jobCtx, req.Links, func(event cmrd.ProgressEvent) {
			s.updateJob(jobID, func(state *jobState) {
				state.Phase = fallback(event.Phase, state.Phase)
				state.Percent = event.Percent
				state.Message = fallback(event.Message, state.Message)
				state.Done = event.Done
				if event.Done && state.Finished.IsZero() {
					state.Finished = time.Now()
				}
				if event.Err != nil {
					state.ErrText = event.Err.Error()
					state.Done = true
					state.Finished = time.Now()
				}
			})
		})

		if err != nil {
			s.updateJob(jobID, func(state *jobState) {
				if errors.Is(err, context.Canceled) && state.Phase == "canceled" {
					if state.Finished.IsZero() {
						state.Finished = time.Now()
					}
					return
				}
				state.Phase = "failed"
				state.Done = true
				state.ErrText = err.Error()
				state.Message = "download failed"
				state.Finished = time.Now()
			})
			return
		}

		s.updateJob(jobID, func(state *jobState) {
			state.Phase = "done"
			state.Percent = 100
			state.Done = true
			state.Message = "download completed"
			state.Finished = time.Now()
		})
	}()

	return &pb.StartDownloadResponse{JobID: jobID}, nil
}

func (s *Server) GetProgress(_ context.Context, req *pb.GetProgressRequest) (*pb.GetProgressResponse, error) {
	if req == nil || strings.TrimSpace(req.JobID) == "" {
		return nil, status.Error(codes.InvalidArgument, "job_id is required")
	}

	state, ok := s.getJob(req.JobID)
	if !ok {
		return nil, status.Error(codes.NotFound, "job not found")
	}
	return toProgressResponse(state), nil
}

func (s *Server) SubscribeProgress(req *pb.GetProgressRequest, stream pb.CMRDService_SubscribeProgressServer) error {
	if req == nil || strings.TrimSpace(req.JobID) == "" {
		return status.Error(codes.InvalidArgument, "job_id is required")
	}

	subID, updates, err := s.subscribe(req.JobID)
	if err != nil {
		if errors.Is(err, errJobNotFound) {
			return status.Error(codes.NotFound, "job not found")
		}
		return status.Errorf(codes.Internal, "subscribe: %v", err)
	}
	defer s.unsubscribe(req.JobID, subID)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case state, ok := <-updates:
			if !ok {
				return nil
			}
			if err := stream.Send(toProgressResponse(&state)); err != nil {
				return err
			}
			if state.Done {
				return nil
			}
		}
	}
}

func (s *Server) StopJob(_ context.Context, req *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	if req == nil || strings.TrimSpace(req.JobID) == "" {
		return nil, status.Error(codes.InvalidArgument, "job_id is required")
	}

	state, ok := s.getJob(req.JobID)
	if !ok {
		return nil, status.Error(codes.NotFound, "job not found")
	}

	if state.Cancel != nil {
		state.Cancel()
	}

	s.updateJob(req.JobID, func(current *jobState) {
		current.Phase = "canceled"
		current.Done = true
		current.Message = "job canceled"
		current.Finished = time.Now()
	})

	return &pb.StopJobResponse{Stopped: true}, nil
}

func (s *Server) ListJobs(_ context.Context, _ *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	response := &pb.ListJobsResponse{
		Jobs: make([]*pb.JobInfo, 0, len(s.jobs)),
	}
	for _, state := range s.jobs {
		response.Jobs = append(response.Jobs, &pb.JobInfo{
			JobID:   state.JobID,
			Phase:   state.Phase,
			Percent: float32(state.Percent),
			Message: state.Message,
			Done:    state.Done,
			Error:   state.ErrText,
		})
	}
	return response, nil
}

var errJobNotFound = errors.New("job not found")

func nextJobID() string {
	value := jobCounter.Add(1)
	return fmt.Sprintf("job-%d-%06d", time.Now().Unix(), value)
}

func (s *Server) setJob(state *jobState) {
	s.mu.Lock()
	s.jobs[state.JobID] = state
	s.mu.Unlock()
}

func (s *Server) getJob(jobID string) (*jobState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.jobs[jobID]
	if !ok {
		return nil, false
	}
	clone := *state
	return &clone, true
}

func (s *Server) updateJob(jobID string, update func(*jobState)) {
	s.mu.Lock()
	state, ok := s.jobs[jobID]
	if !ok {
		s.mu.Unlock()
		return
	}
	update(state)
	snapshot := *state
	subscribers := make([]chan jobState, 0, len(s.subs[jobID]))
	for _, ch := range s.subs[jobID] {
		subscribers = append(subscribers, ch)
	}
	s.mu.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- snapshot:
		default:
		}
	}
}

func (s *Server) subscribe(jobID string) (uint64, <-chan jobState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.jobs[jobID]
	if !ok {
		return 0, nil, errJobNotFound
	}
	if s.subs[jobID] == nil {
		s.subs[jobID] = make(map[uint64]chan jobState)
	}

	id := subCounter.Add(1)
	ch := make(chan jobState, 16)
	s.subs[jobID][id] = ch
	ch <- *state
	return id, ch, nil
}

func (s *Server) unsubscribe(jobID string, subID uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobSubs, ok := s.subs[jobID]
	if !ok {
		return
	}
	ch, ok := jobSubs[subID]
	if !ok {
		return
	}
	delete(jobSubs, subID)
	close(ch)
	if len(jobSubs) == 0 {
		delete(s.subs, jobID)
	}
}

func toProgressResponse(state *jobState) *pb.GetProgressResponse {
	return &pb.GetProgressResponse{
		JobID:   state.JobID,
		Phase:   state.Phase,
		Percent: float32(state.Percent),
		Message: state.Message,
		Done:    state.Done,
		Error:   state.ErrText,
	}
}

func fallback(value string, fallbackValue string) string {
	if strings.TrimSpace(value) == "" {
		return fallbackValue
	}
	return value
}
