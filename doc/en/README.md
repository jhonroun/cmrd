# CMRD: documentation (EN)

## Navigation
- Main guide: `doc/en/README.md`
- CLI: `doc/en/CLI.md`
- gRPC API: `doc/en/GRPC.md`

## 1. Purpose
CMRD rewrites the reference `reference/cloud_mail_downloader.php` into Go and provides:
- CLI mode.
- Bubble Tea TUI mode with live screen refresh.
- gRPC API for WEB UI / desktop GUI clients.
- library API for embedding into other Go projects.

## 2. Architecture
- `cmd/cmrd`: CLI executable.
- `pkg/cmrd`: public library API.
- `internal/cloudmail`: Cloud.Mail link resolution logic.
- `internal/aria2`: aria2 input generation and process runner.
- `internal/tui`: Bubble Tea progress UI.
- `internal/grpcapi`: gRPC server implementation.
- `api/proto/cmrd/v1/cmrd.proto`: gRPC contract.

## 3. Operating Modes
- CLI:
  - `cmrd resolve`
  - `cmrd download`
  - `cmrd serve-grpc`
- TUI:
  - enabled via `--tui=true` in `download`.
  - keys: `h`/`?` (help), `q`/`Ctrl+C` (quit).
- gRPC:
  - start server and control jobs externally.
- Library:
  - import `github.com/jhonroun/cmrd/pkg/cmrd`.

## 4. Installation
1. Go 1.24+ is required.
2. `aria2c` is required (from `PATH` or via env var).
3. Build:
   - `go build -o cmrd ./cmd/cmrd`

## 5. Configuration
- Environment:
  - `CMRD_ARIA2C_PATH`: absolute path to `aria2c`.
- CLI flags:
  - `--links`: links file path.
  - `--dir`: destination directory.
  - `--proxy`, `--proxy-auth`: proxy configuration.
  - `--timeout`: HTTP timeout.
  - `--tui`: enable/disable TUI mode.

## 6. Quick Start (EN)
1. Create `links.txt` (one Cloud.Mail public link per line).
2. Resolve links:
   - `cmrd resolve --links links.txt`
3. Download with TUI:
   - `cmrd download --links links.txt --dir downloads --tui=true`
4. Start gRPC API:
   - `cmrd serve-grpc --listen :50051`

## 7. First Run (EN)
1. Use `resolve` first and validate mapped output paths.
2. Ensure `aria2c` is available (`CMRD_ARIA2C_PATH` or `PATH`).
3. Run `download` in TUI mode.
4. Start `serve-grpc` if you need an external GUI client.

## 8. gRPC API (short)
Methods:
- `ResolveLinks`
- `StartDownload`
- `GetProgress`
- `SubscribeProgress` (server stream)
- `StopJob`
- `ListJobs`

Protocol:
- `api/proto/cmrd/v1/cmrd.proto`

## 9. Library Usage
```go
package main

import (
	"context"
	"log"

	"github.com/jhonroun/cmrd/pkg/cmrd"
)

func main() {
	cfg := cmrd.DefaultConfig()
	client, err := cmrd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	links := []string{"https://cloud.mail.ru/public/9bFs/gVzxjU5uC"}
	files, err := client.Resolve(context.Background(), links)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resolved %d files", len(files))
}
```

## 10. License
Project license model follows aria2 licensing (`GPL-2.0-or-later`).
