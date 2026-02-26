# CMRD: документация (RU)

## Навигация
- Основное руководство: `doc/ru/README.md`
- CLI: `doc/ru/CLI.md`
- gRPC API: `doc/ru/GRPC.md`

## 1. Назначение
CMRD переписывает референсный `reference/cloud_mail_downloader.php` на Go и предоставляет:
- CLI режим.
- TUI режим на Bubble Tea с обновляемым экраном.
- gRPC API для WEB UI/GUI клиентов.
- библиотечный API для встраивания в другие Go-проекты.

## 2. Архитектура
- `cmd/cmrd`: исполняемый файл CLI.
- `pkg/cmrd`: публичный библиотечный API.
- `internal/cloudmail`: резолв ссылок Cloud.Mail через API.
- `internal/aria2`: генерация input и запуск aria2c.
- `internal/tui`: Bubble Tea интерфейс с прогрессом.
- `internal/grpcapi`: gRPC сервер и сервисные методы.
- `api/proto/cmrd/v1/cmrd.proto`: описание gRPC контракта.

## 3. Режимы работы
- CLI:
  - `cmrd resolve`
  - `cmrd download`
  - `cmrd serve-grpc`
- TUI:
  - активируется флагом `--tui=true` в `download`.
  - клавиши: `h`/`?` (помощь), `q`/`Ctrl+C` (выход).
- gRPC:
  - сервер для запуска задач и чтения прогресса внешними клиентами.
- Library:
  - пакет `github.com/jhonroun/cmrd/pkg/cmrd`.

## 4. Установка
1. Нужен Go 1.24+.
2. Нужен `aria2c` (в `PATH` или через переменную окружения).
3. Сборка:
   - `go build -o cmrd ./cmd/cmrd`

## 5. Конфигурация
- Переменные окружения:
  - `CMRD_ARIA2C_PATH`: путь к бинарнику aria2c.
- CLI-флаги:
  - `--links`: путь к файлу со ссылками.
  - `--dir`: папка назначения.
  - `--proxy`, `--proxy-auth`: прокси-настройки.
  - `--timeout`: HTTP timeout.
  - `--tui`: включение/выключение TUI.

## 6. Quick Start (RU)
1. Подготовьте `links.txt`:
   - одна публичная ссылка Cloud.Mail на строку.
2. Проверка резолва:
   - `cmrd resolve --links links.txt`
3. Запуск скачивания:
   - `cmrd download --links links.txt --dir downloads --tui=true`
4. Запуск gRPC:
   - `cmrd serve-grpc --listen :50051`

## 7. First Run (RU)
1. Запустите `resolve` и убедитесь, что ссылки корректно раскрываются.
2. Убедитесь, что `aria2c` доступен (`CMRD_ARIA2C_PATH` или `PATH`).
3. Запустите `download` в TUI-режиме.
4. Для интеграции GUI поднимите `serve-grpc`.

## 8. gRPC API (кратко)
Методы:
- `ResolveLinks`
- `StartDownload`
- `GetProgress`
- `SubscribeProgress` (server stream)
- `StopJob`
- `ListJobs`

Протокол:
- `api/proto/cmrd/v1/cmrd.proto`

## 9. Использование как библиотека
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

## 10. Лицензия
Тип лицензии проекта синхронизирован с моделью лицензирования aria2 (`GPL-2.0-or-later`).
