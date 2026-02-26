package pb

import "github.com/golang/protobuf/proto"

type ResolveLinksRequest struct {
	Links []string `protobuf:"bytes,1,rep,name=links,proto3" json:"links,omitempty"`
}

func (m *ResolveLinksRequest) Reset()         { *m = ResolveLinksRequest{} }
func (m *ResolveLinksRequest) String() string { return proto.CompactTextString(m) }
func (*ResolveLinksRequest) ProtoMessage()    {}

type ResolvedFile struct {
	URL    string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Output string `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
}

func (m *ResolvedFile) Reset()         { *m = ResolvedFile{} }
func (m *ResolvedFile) String() string { return proto.CompactTextString(m) }
func (*ResolvedFile) ProtoMessage()    {}

type ResolveLinksResponse struct {
	Files []*ResolvedFile `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
}

func (m *ResolveLinksResponse) Reset()         { *m = ResolveLinksResponse{} }
func (m *ResolveLinksResponse) String() string { return proto.CompactTextString(m) }
func (*ResolveLinksResponse) ProtoMessage()    {}

type StartDownloadRequest struct {
	Links       []string `protobuf:"bytes,1,rep,name=links,proto3" json:"links,omitempty"`
	DownloadDir string   `protobuf:"bytes,2,opt,name=download_dir,json=downloadDir,proto3" json:"download_dir,omitempty"`
}

func (m *StartDownloadRequest) Reset()         { *m = StartDownloadRequest{} }
func (m *StartDownloadRequest) String() string { return proto.CompactTextString(m) }
func (*StartDownloadRequest) ProtoMessage()    {}

type StartDownloadResponse struct {
	JobID string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (m *StartDownloadResponse) Reset()         { *m = StartDownloadResponse{} }
func (m *StartDownloadResponse) String() string { return proto.CompactTextString(m) }
func (*StartDownloadResponse) ProtoMessage()    {}

type GetProgressRequest struct {
	JobID string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (m *GetProgressRequest) Reset()         { *m = GetProgressRequest{} }
func (m *GetProgressRequest) String() string { return proto.CompactTextString(m) }
func (*GetProgressRequest) ProtoMessage()    {}

type GetProgressResponse struct {
	JobID   string  `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Phase   string  `protobuf:"bytes,2,opt,name=phase,proto3" json:"phase,omitempty"`
	Percent float32 `protobuf:"fixed32,3,opt,name=percent,proto3" json:"percent,omitempty"`
	Message string  `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	Done    bool    `protobuf:"varint,5,opt,name=done,proto3" json:"done,omitempty"`
	Error   string  `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *GetProgressResponse) Reset()         { *m = GetProgressResponse{} }
func (m *GetProgressResponse) String() string { return proto.CompactTextString(m) }
func (*GetProgressResponse) ProtoMessage()    {}

type StopJobRequest struct {
	JobID string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (m *StopJobRequest) Reset()         { *m = StopJobRequest{} }
func (m *StopJobRequest) String() string { return proto.CompactTextString(m) }
func (*StopJobRequest) ProtoMessage()    {}

type StopJobResponse struct {
	Stopped bool `protobuf:"varint,1,opt,name=stopped,proto3" json:"stopped,omitempty"`
}

func (m *StopJobResponse) Reset()         { *m = StopJobResponse{} }
func (m *StopJobResponse) String() string { return proto.CompactTextString(m) }
func (*StopJobResponse) ProtoMessage()    {}

type ListJobsRequest struct{}

func (m *ListJobsRequest) Reset()         { *m = ListJobsRequest{} }
func (m *ListJobsRequest) String() string { return proto.CompactTextString(m) }
func (*ListJobsRequest) ProtoMessage()    {}

type JobInfo struct {
	JobID   string  `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Phase   string  `protobuf:"bytes,2,opt,name=phase,proto3" json:"phase,omitempty"`
	Percent float32 `protobuf:"fixed32,3,opt,name=percent,proto3" json:"percent,omitempty"`
	Message string  `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	Done    bool    `protobuf:"varint,5,opt,name=done,proto3" json:"done,omitempty"`
	Error   string  `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *JobInfo) Reset()         { *m = JobInfo{} }
func (m *JobInfo) String() string { return proto.CompactTextString(m) }
func (*JobInfo) ProtoMessage()    {}

type ListJobsResponse struct {
	Jobs []*JobInfo `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
}

func (m *ListJobsResponse) Reset()         { *m = ListJobsResponse{} }
func (m *ListJobsResponse) String() string { return proto.CompactTextString(m) }
func (*ListJobsResponse) ProtoMessage()    {}

var (
	_ proto.Message = (*ResolveLinksRequest)(nil)
	_ proto.Message = (*ResolveLinksResponse)(nil)
	_ proto.Message = (*StartDownloadRequest)(nil)
	_ proto.Message = (*StartDownloadResponse)(nil)
	_ proto.Message = (*GetProgressRequest)(nil)
	_ proto.Message = (*GetProgressResponse)(nil)
	_ proto.Message = (*StopJobRequest)(nil)
	_ proto.Message = (*StopJobResponse)(nil)
	_ proto.Message = (*ListJobsRequest)(nil)
	_ proto.Message = (*ListJobsResponse)(nil)
)
