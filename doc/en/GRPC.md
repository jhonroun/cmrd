# CMRD gRPC API (EN)

## Protocol
- Contract file: `api/proto/cmrd/v1/cmrd.proto`
- Service: `cmrd.v1.CMRDService`

## Methods
- `ResolveLinks(ResolveLinksRequest) returns (ResolveLinksResponse)`
- `StartDownload(StartDownloadRequest) returns (StartDownloadResponse)`
- `GetProgress(GetProgressRequest) returns (GetProgressResponse)`
- `SubscribeProgress(GetProgressRequest) returns (stream GetProgressResponse)`
- `StopJob(StopJobRequest) returns (StopJobResponse)`
- `ListJobs(ListJobsRequest) returns (ListJobsResponse)`

## Method Intent
- `ResolveLinks`: resolve links without running download.
- `StartDownload`: create and start a background download job, returns `job_id`.
- `GetProgress`: polling progress for a specific `job_id`.
- `SubscribeProgress`: live progress updates over server stream.
- `StopJob`: cancel a running job.
- `ListJobs`: list known jobs in server memory.

## Typical Client Flow
1. Call `StartDownload` and keep returned `job_id`.
2. Subscribe using `SubscribeProgress`.
3. Read stream updates until `done=true`.
4. Optionally call `StopJob` to cancel.

## Minimal Go Example
```go
conn, err := grpc.Dial(
    "127.0.0.1:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil { panic(err) }
defer conn.Close()

client := pb.NewCMRDServiceClient(conn)

start, err := client.StartDownload(ctx, &pb.StartDownloadRequest{
    Links: []string{"https://cloud.mail.ru/public/9bFs/gVzxjU5uC"},
})
if err != nil { panic(err) }

stream, err := client.SubscribeProgress(ctx, &pb.GetProgressRequest{JobID: start.JobID})
if err != nil { panic(err) }

for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        panic(err)
    }
    if msg.Done {
        break
    }
}
```
