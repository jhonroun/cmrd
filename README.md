# CMRD

[![CI](https://github.com/jhonroun/cmrd/actions/workflows/ci.yml/badge.svg)](https://github.com/jhonroun/cmrd/actions/workflows/ci.yml)
![Version](https://img.shields.io/badge/version-1.0.0-blue)

Cloud.Mail public link downloader rewritten from PHP to idiomatic Go with CLI, TUI (Bubble Tea), gRPC API and library mode.

Загрузчик публичных ссылок Cloud.Mail, переписанный с PHP на идиоматичный Go, с режимами CLI, TUI (Bubble Tea), gRPC API и библиотечным использованием.

## Documentation / Документация
- English docs: `doc/en/README.md`
- Русская документация: `doc/ru/README.md`
- English CLI: `doc/en/CLI.md`
- Русский CLI: `doc/ru/CLI.md`
- English gRPC: `doc/en/GRPC.md`
- Русский gRPC: `doc/ru/GRPC.md`

## English Summary
CMRD resolves Cloud.Mail public links into direct file URLs, prepares aria2c input, and runs multi-file downloads with resumable behavior.

### Quick Start (EN)
1. Build CLI:
   - `go build -o cmrd ./cmd/cmrd`
2. Set aria2c path (optional, if aria2c is not in PATH):
   - Linux/macOS: `export CMRD_ARIA2C_PATH=/usr/bin/aria2c`
   - Windows PowerShell: `$env:CMRD_ARIA2C_PATH="C:\tools\aria2c.exe"`
3. Create `links.txt` with one Cloud.Mail public link per line.
4. Run:
   - `./cmrd download --links links.txt --dir downloads --tui=true`

### First Run (EN)
1. Start with resolve-only mode:
   - `./cmrd resolve --links links.txt`
2. Verify resolved paths.
3. Start download:
   - `./cmrd download --links links.txt --dir downloads`
4. Optional: start gRPC API for external UI:
   - `./cmrd serve-grpc --listen :50051`

## Краткое описание на русском
CMRD преобразует публичные ссылки Cloud.Mail в прямые URL файлов, формирует вход для aria2c и запускает многопоточную загрузку с поддержкой докачки.

### Quick Start (RU)
1. Соберите CLI:
   - `go build -o cmrd ./cmd/cmrd`
2. Укажите путь к aria2c (опционально, если нет в PATH):
   - Linux/macOS: `export CMRD_ARIA2C_PATH=/usr/bin/aria2c`
   - Windows PowerShell: `$env:CMRD_ARIA2C_PATH="C:\tools\aria2c.exe"`
3. Создайте `links.txt` (по одной публичной ссылке Cloud.Mail в строке).
4. Запустите:
   - `./cmrd download --links links.txt --dir downloads --tui=true`

### First Run (RU)
1. Сначала проверьте резолв ссылок:
   - `./cmrd resolve --links links.txt`
2. Проверьте сформированные пути файлов.
3. Запустите скачивание:
   - `./cmrd download --links links.txt --dir downloads`
4. Опционально поднимите gRPC API для внешнего UI:
   - `./cmrd serve-grpc --listen :50051`
