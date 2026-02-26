# CMRD CLI (EN)

## Commands
- `cmrd help`
- `cmrd version`
- `cmrd resolve`
- `cmrd download`
- `cmrd serve-grpc`

## cmrd resolve
Resolves Cloud.Mail public links into direct file URLs without downloading.

Example:
```bash
cmrd resolve --links links.txt --json
```

Flags:
- `--links` links file path (default `links.txt`).
- `--json` print JSON output.
- `--timeout` HTTP timeout, e.g. `45s`.
- `--proxy` proxy URL or host:port.
- `--proxy-auth` proxy auth in `user:pass` format.

## cmrd download
Resolves links and runs aria2c downloader.

Example:
```bash
cmrd download --links links.txt --dir downloads --tui=true
```

Flags:
- `--links` links file path.
- `--dir` destination directory.
- `--aria2c` explicit aria2c binary path.
- `--timeout` HTTP timeout.
- `--proxy` proxy URL or host:port.
- `--proxy-auth` proxy auth.
- `--tui` enable/disable Bubble Tea TUI.
- `--keep-input` keep temporary aria2 input file after completion.

## cmrd serve-grpc
Starts the gRPC API server for WEB UI/GUI clients.

Example:
```bash
cmrd serve-grpc --listen :50051 --dir downloads
```

Flags:
- `--listen` bind address.
- `--dir` default download destination directory.
- `--aria2c` aria2c path.
- `--timeout` HTTP timeout.
- `--proxy` proxy configuration.
- `--proxy-auth` proxy auth.
- `--keep-input` keep aria2 input file.

## Environment Variables
- `CMRD_ARIA2C_PATH` path to aria2c binary when `--aria2c` is not set.

## links.txt Format
- One Cloud.Mail public link per line.
- Empty lines are ignored.
- Lines starting with `#` are treated as comments.

Example:
```text
# public links
https://cloud.mail.ru/public/9bFs/gVzxjU5uC
https://cloud.mail.ru/public/3umo/mCi4k2ZTs
```
