package grpcapi

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/jhonroun/cmrd/internal/grpcapi/pb"
	"github.com/jhonroun/cmrd/pkg/cmrd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestGRPCWireSubscribeProgress(t *testing.T) {
	const bufferSize = 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	server := grpc.NewServer()

	service := NewServerWithFactory(cmrd.DefaultConfig(), func(cmrd.Config) (serviceClient, error) {
		return &mockServiceClient{
			downloadFn: func(ctx context.Context, _ []string, onProgress cmrd.ProgressHandler) error {
				onProgress(cmrd.ProgressEvent{Phase: "resolve", Message: "resolve complete", TotalFiles: 1})
				onProgress(cmrd.ProgressEvent{Phase: "download", Percent: 100, Message: "done", Done: true, TotalFiles: 1})
				return nil
			},
		}, nil
	})
	pb.RegisterCMRDServiceServer(server, service)

	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		server.Stop()
	})

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewCMRDServiceClient(conn)
	start, err := client.StartDownload(ctx, &pb.StartDownloadRequest{
		Links: []string{"https://cloud.mail.ru/public/9bFs/gVzxjU5uC"},
	})
	if err != nil {
		t.Fatalf("start download: %v", err)
	}

	stream, err := client.SubscribeProgress(ctx, &pb.GetProgressRequest{JobID: start.JobID})
	if err != nil {
		t.Fatalf("subscribe progress: %v", err)
	}

	count := 0
	done := false
	for {
		msg, recvErr := stream.Recv()
		if recvErr == io.EOF {
			break
		}
		if recvErr != nil {
			t.Fatalf("recv: %v", recvErr)
		}
		count++
		done = msg.Done
	}

	if count == 0 {
		t.Fatalf("expected at least one message")
	}
	if !done {
		t.Fatalf("expected final done message")
	}
}
