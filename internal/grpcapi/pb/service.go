package pb

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const CMRDServiceServiceName = "cmrd.v1.CMRDService"

type CMRDServiceClient interface {
	ResolveLinks(ctx context.Context, in *ResolveLinksRequest, opts ...grpc.CallOption) (*ResolveLinksResponse, error)
	StartDownload(ctx context.Context, in *StartDownloadRequest, opts ...grpc.CallOption) (*StartDownloadResponse, error)
	GetProgress(ctx context.Context, in *GetProgressRequest, opts ...grpc.CallOption) (*GetProgressResponse, error)
	SubscribeProgress(ctx context.Context, in *GetProgressRequest, opts ...grpc.CallOption) (CMRDService_SubscribeProgressClient, error)
	StopJob(ctx context.Context, in *StopJobRequest, opts ...grpc.CallOption) (*StopJobResponse, error)
	ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error)
}

type cmrdServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCMRDServiceClient(cc grpc.ClientConnInterface) CMRDServiceClient {
	return &cmrdServiceClient{cc: cc}
}

func (c *cmrdServiceClient) ResolveLinks(ctx context.Context, in *ResolveLinksRequest, opts ...grpc.CallOption) (*ResolveLinksResponse, error) {
	out := new(ResolveLinksResponse)
	err := c.cc.Invoke(ctx, "/"+CMRDServiceServiceName+"/ResolveLinks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cmrdServiceClient) StartDownload(ctx context.Context, in *StartDownloadRequest, opts ...grpc.CallOption) (*StartDownloadResponse, error) {
	out := new(StartDownloadResponse)
	err := c.cc.Invoke(ctx, "/"+CMRDServiceServiceName+"/StartDownload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cmrdServiceClient) GetProgress(ctx context.Context, in *GetProgressRequest, opts ...grpc.CallOption) (*GetProgressResponse, error) {
	out := new(GetProgressResponse)
	err := c.cc.Invoke(ctx, "/"+CMRDServiceServiceName+"/GetProgress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cmrdServiceClient) SubscribeProgress(ctx context.Context, in *GetProgressRequest, opts ...grpc.CallOption) (CMRDService_SubscribeProgressClient, error) {
	stream, err := c.cc.NewStream(ctx, &CMRDServiceServiceDesc.Streams[0], "/"+CMRDServiceServiceName+"/SubscribeProgress", opts...)
	if err != nil {
		return nil, err
	}
	client := &cmrdServiceSubscribeProgressClient{ClientStream: stream}
	if err := client.SendMsg(in); err != nil {
		return nil, err
	}
	if err := client.CloseSend(); err != nil {
		return nil, err
	}
	return client, nil
}

type CMRDService_SubscribeProgressClient interface {
	Recv() (*GetProgressResponse, error)
	grpc.ClientStream
}

type cmrdServiceSubscribeProgressClient struct {
	grpc.ClientStream
}

func (x *cmrdServiceSubscribeProgressClient) Recv() (*GetProgressResponse, error) {
	msg := new(GetProgressResponse)
	if err := x.ClientStream.RecvMsg(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *cmrdServiceClient) StopJob(ctx context.Context, in *StopJobRequest, opts ...grpc.CallOption) (*StopJobResponse, error) {
	out := new(StopJobResponse)
	err := c.cc.Invoke(ctx, "/"+CMRDServiceServiceName+"/StopJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cmrdServiceClient) ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error) {
	out := new(ListJobsResponse)
	err := c.cc.Invoke(ctx, "/"+CMRDServiceServiceName+"/ListJobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type CMRDServiceServer interface {
	ResolveLinks(context.Context, *ResolveLinksRequest) (*ResolveLinksResponse, error)
	StartDownload(context.Context, *StartDownloadRequest) (*StartDownloadResponse, error)
	GetProgress(context.Context, *GetProgressRequest) (*GetProgressResponse, error)
	SubscribeProgress(*GetProgressRequest, CMRDService_SubscribeProgressServer) error
	StopJob(context.Context, *StopJobRequest) (*StopJobResponse, error)
	ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error)
	mustEmbedUnimplementedCMRDServiceServer()
}

type UnimplementedCMRDServiceServer struct{}

func (UnimplementedCMRDServiceServer) ResolveLinks(context.Context, *ResolveLinksRequest) (*ResolveLinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveLinks not implemented")
}

func (UnimplementedCMRDServiceServer) StartDownload(context.Context, *StartDownloadRequest) (*StartDownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartDownload not implemented")
}

func (UnimplementedCMRDServiceServer) GetProgress(context.Context, *GetProgressRequest) (*GetProgressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProgress not implemented")
}

func (UnimplementedCMRDServiceServer) SubscribeProgress(*GetProgressRequest, CMRDService_SubscribeProgressServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeProgress not implemented")
}

func (UnimplementedCMRDServiceServer) StopJob(context.Context, *StopJobRequest) (*StopJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopJob not implemented")
}

func (UnimplementedCMRDServiceServer) ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListJobs not implemented")
}

func (UnimplementedCMRDServiceServer) mustEmbedUnimplementedCMRDServiceServer() {}

func RegisterCMRDServiceServer(s grpc.ServiceRegistrar, srv CMRDServiceServer) {
	s.RegisterService(&CMRDServiceServiceDesc, srv)
}

func _CMRDService_ResolveLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CMRDServiceServer).ResolveLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + CMRDServiceServiceName + "/ResolveLinks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CMRDServiceServer).ResolveLinks(ctx, req.(*ResolveLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CMRDService_StartDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartDownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CMRDServiceServer).StartDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + CMRDServiceServiceName + "/StartDownload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CMRDServiceServer).StartDownload(ctx, req.(*StartDownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CMRDService_GetProgress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProgressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CMRDServiceServer).GetProgress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + CMRDServiceServiceName + "/GetProgress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CMRDServiceServer).GetProgress(ctx, req.(*GetProgressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CMRDService_SubscribeProgress_Handler(srv interface{}, stream grpc.ServerStream) error {
	req := new(GetProgressRequest)
	if err := stream.RecvMsg(req); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return srv.(CMRDServiceServer).SubscribeProgress(req, &cmrdServiceSubscribeProgressServer{ServerStream: stream})
}

type CMRDService_SubscribeProgressServer interface {
	Send(*GetProgressResponse) error
	grpc.ServerStream
}

type cmrdServiceSubscribeProgressServer struct {
	grpc.ServerStream
}

func (x *cmrdServiceSubscribeProgressServer) Send(msg *GetProgressResponse) error {
	return x.ServerStream.SendMsg(msg)
}

func _CMRDService_StopJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CMRDServiceServer).StopJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + CMRDServiceServiceName + "/StopJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CMRDServiceServer).StopJob(ctx, req.(*StopJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CMRDService_ListJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CMRDServiceServer).ListJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + CMRDServiceServiceName + "/ListJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CMRDServiceServer).ListJobs(ctx, req.(*ListJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var CMRDServiceServiceDesc = grpc.ServiceDesc{
	ServiceName: CMRDServiceServiceName,
	HandlerType: (*CMRDServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResolveLinks",
			Handler:    _CMRDService_ResolveLinks_Handler,
		},
		{
			MethodName: "StartDownload",
			Handler:    _CMRDService_StartDownload_Handler,
		},
		{
			MethodName: "GetProgress",
			Handler:    _CMRDService_GetProgress_Handler,
		},
		{
			MethodName: "StopJob",
			Handler:    _CMRDService_StopJob_Handler,
		},
		{
			MethodName: "ListJobs",
			Handler:    _CMRDService_ListJobs_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeProgress",
			Handler:       _CMRDService_SubscribeProgress_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/proto/cmrd/v1/cmrd.proto",
}
