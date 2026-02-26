# CMRD gRPC API (RU)

## Протокол
- Файл контракта: `api/proto/cmrd/v1/cmrd.proto`
- Сервис: `cmrd.v1.CMRDService`

## Методы
- `ResolveLinks(ResolveLinksRequest) returns (ResolveLinksResponse)`
- `StartDownload(StartDownloadRequest) returns (StartDownloadResponse)`
- `GetProgress(GetProgressRequest) returns (GetProgressResponse)`
- `SubscribeProgress(GetProgressRequest) returns (stream GetProgressResponse)`
- `StopJob(StopJobRequest) returns (StopJobResponse)`
- `ListJobs(ListJobsRequest) returns (ListJobsResponse)`

## Назначение методов
- `ResolveLinks`: резолв ссылок без запуска скачивания.
- `StartDownload`: запуск новой фоновой задачи скачивания, возвращает `job_id`.
- `GetProgress`: polling-состояние задачи по `job_id`.
- `SubscribeProgress`: live-обновления состояния задачи по stream.
- `StopJob`: остановка задачи по `job_id`.
- `ListJobs`: список известных задач в памяти сервера.

## Типовой сценарий клиента
1. Вызвать `StartDownload` и получить `job_id`.
2. Подписаться через `SubscribeProgress` на `job_id`.
3. Читать stream до `done=true`.
4. При необходимости вызвать `StopJob`.

## Минимальный пример (Go)
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
