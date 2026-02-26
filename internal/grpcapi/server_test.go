package grpcapi

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/jhonroun/cmrd/internal/grpcapi/pb"
	"github.com/jhonroun/cmrd/pkg/cmrd"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type mockServiceClient struct {
	resolveResult []cmrd.FileTask
	resolveErr    error
	downloadErr   error
	downloadFn    func(context.Context, []string, cmrd.ProgressHandler) error
}

func (m *mockServiceClient) Resolve(_ context.Context, _ []string) ([]cmrd.FileTask, error) {
	if m.resolveErr != nil {
		return nil, m.resolveErr
	}
	return m.resolveResult, nil
}

func (m *mockServiceClient) Download(ctx context.Context, links []string, onProgress cmrd.ProgressHandler) error {
	if m.downloadFn != nil {
		return m.downloadFn(ctx, links, onProgress)
	}
	return m.downloadErr
}

type progressStreamStub struct {
	ctx      context.Context
	messages []*pb.GetProgressResponse
}

func newProgressStreamStub(ctx context.Context) *progressStreamStub {
	return &progressStreamStub{ctx: ctx}
}

func (s *progressStreamStub) Send(msg *pb.GetProgressResponse) error {
	s.messages = append(s.messages, msg)
	return nil
}

func (s *progressStreamStub) SetHeader(metadata.MD) error  { return nil }
func (s *progressStreamStub) SendHeader(metadata.MD) error { return nil }
func (s *progressStreamStub) SetTrailer(metadata.MD)       {}
func (s *progressStreamStub) Context() context.Context     { return s.ctx }
func (s *progressStreamStub) SendMsg(interface{}) error    { return nil }
func (s *progressStreamStub) RecvMsg(interface{}) error    { return io.EOF }

func TestResolveLinksInvalidArgument(t *testing.T) {
	server := NewServerWithFactory(cmrd.DefaultConfig(), func(cmrd.Config) (serviceClient, error) {
		return &mockServiceClient{}, nil
	})

	_, err := server.ResolveLinks(context.Background(), &pb.ResolveLinksRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %s", status.Code(err))
	}
}

func TestStartDownloadAndSubscribeProgress(t *testing.T) {
	server := NewServerWithFactory(cmrd.DefaultConfig(), func(cmrd.Config) (serviceClient, error) {
		return &mockServiceClient{
			downloadFn: func(ctx context.Context, _ []string, onProgress cmrd.ProgressHandler) error {
				onProgress(cmrd.ProgressEvent{Phase: "resolve", Message: "resolve complete", TotalFiles: 2})
				onProgress(cmrd.ProgressEvent{Phase: "download", Percent: 25, Message: "downloading", TotalFiles: 2})
				onProgress(cmrd.ProgressEvent{Phase: "download", Percent: 100, Message: "done", Done: true, TotalFiles: 2})
				return nil
			},
		}, nil
	})

	start, err := server.StartDownload(context.Background(), &pb.StartDownloadRequest{
		Links: []string{"https://cloud.mail.ru/public/9bFs/gVzxjU5uC"},
	})
	if err != nil {
		t.Fatalf("start download: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream := newProgressStreamStub(ctx)
	if err := server.SubscribeProgress(&pb.GetProgressRequest{JobID: start.JobID}, stream); err != nil {
		t.Fatalf("subscribe progress: %v", err)
	}
	if len(stream.messages) == 0 {
		t.Fatalf("expected at least one progress message")
	}

	last := stream.messages[len(stream.messages)-1]
	if !last.Done {
		t.Fatalf("expected final done state")
	}
}

func TestStopJob(t *testing.T) {
	server := NewServerWithFactory(cmrd.DefaultConfig(), func(cmrd.Config) (serviceClient, error) {
		return &mockServiceClient{
			downloadFn: func(ctx context.Context, _ []string, _ cmrd.ProgressHandler) error {
				<-ctx.Done()
				return ctx.Err()
			},
		}, nil
	})

	start, err := server.StartDownload(context.Background(), &pb.StartDownloadRequest{
		Links: []string{"https://cloud.mail.ru/public/9bFs/gVzxjU5uC"},
	})
	if err != nil {
		t.Fatalf("start download: %v", err)
	}

	stop, err := server.StopJob(context.Background(), &pb.StopJobRequest{JobID: start.JobID})
	if err != nil {
		t.Fatalf("stop job: %v", err)
	}
	if !stop.Stopped {
		t.Fatalf("expected stopped=true")
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		progress, progressErr := server.GetProgress(context.Background(), &pb.GetProgressRequest{JobID: start.JobID})
		if progressErr != nil {
			t.Fatalf("get progress: %v", progressErr)
		}
		if progress.Done {
			if progress.Phase != "canceled" && progress.Phase != "failed" {
				t.Fatalf("unexpected phase after stop: %s", progress.Phase)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting canceled state")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestResolveLinksInternalError(t *testing.T) {
	server := NewServerWithFactory(cmrd.DefaultConfig(), func(cmrd.Config) (serviceClient, error) {
		return &mockServiceClient{
			resolveErr: errors.New("resolve failure"),
		}, nil
	})

	_, err := server.ResolveLinks(context.Background(), &pb.ResolveLinksRequest{
		Links: []string{"https://cloud.mail.ru/public/9bFs/gVzxjU5uC"},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal, got %s", status.Code(err))
	}
}
