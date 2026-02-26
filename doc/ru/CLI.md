# CMRD CLI (RU)

## Команды
- `cmrd help`
- `cmrd version`
- `cmrd resolve`
- `cmrd download`
- `cmrd serve-grpc`

## cmrd resolve
Преобразует публичные ссылки Cloud.Mail в прямые URL файлов без запуска скачивания.

Пример:
```bash
cmrd resolve --links links.txt --json
```

Флаги:
- `--links` путь к файлу ссылок (по умолчанию `links.txt`).
- `--json` вывод в JSON.
- `--timeout` таймаут HTTP, например `45s`.
- `--proxy` прокси URL или host:port.
- `--proxy-auth` авторизация прокси в формате `user:pass`.

## cmrd download
Резолвит ссылки и запускает aria2c.

Пример:
```bash
cmrd download --links links.txt --dir downloads --tui=true
```

Флаги:
- `--links` путь к файлу ссылок.
- `--dir` каталог назначения.
- `--aria2c` путь к бинарнику aria2c.
- `--timeout` таймаут HTTP.
- `--proxy` прокси URL или host:port.
- `--proxy-auth` авторизация прокси.
- `--tui` включить/выключить Bubble Tea TUI.
- `--keep-input` не удалять временный input-файл aria2 после завершения.

## cmrd serve-grpc
Запускает gRPC API-сервер для WEB UI/GUI клиентов.

Пример:
```bash
cmrd serve-grpc --listen :50051 --dir downloads
```

Флаги:
- `--listen` адрес прослушивания.
- `--dir` каталог скачивания по умолчанию.
- `--aria2c` путь к aria2c.
- `--timeout` таймаут HTTP.
- `--proxy` прокси.
- `--proxy-auth` авторизация прокси.
- `--keep-input` сохранять input-файл aria2.

## Переменные окружения
- `CMRD_ARIA2C_PATH` путь к бинарнику aria2c, если флаг `--aria2c` не задан.

## Формат файла links.txt
- Одна публичная ссылка Cloud.Mail в строке.
- Пустые строки игнорируются.
- Строки, начинающиеся с `#`, считаются комментариями.

Пример:
```text
# public links
https://cloud.mail.ru/public/9bFs/gVzxjU5uC
https://cloud.mail.ru/public/3umo/mCi4k2ZTs
```
