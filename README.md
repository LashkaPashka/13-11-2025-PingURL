# Go Link Checker

## Описание проекта

**Go Link Checker** — это веб-сервис на Go, который позволяет:

- Проверять доступность интернет-ресурсов по ссылкам.
- Отправлять как одну ссылку, так и сразу несколько.
- Формировать PDF-отчёт по статусам ссылок для выбранных задач.
- Сохранять состояние задач между перезапусками сервера, чтобы незавершённые задачи продолжали выполняться.

Все данные сохраняются в локальном JSON-файле (`storage/storage.json`).

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/yourusername/13-11-2025-PingURL.git
cd 13-11-2025-PingURL
```

2. Установите переменные окружения в файле (`config/local.yaml`).
```yaml
```

3. Установка зависимостей.
```bash
go mod tidy
```

4. Запуска сервера.
```bash
CONFIG="./config/local.yaml" go run cmd/app/main.go
```

## REST API
1. POST `/url-check` — отправка ссылок для проверки

Request:
```json
{
	"links": [
		"wfsfe.gg",
		"google.com"
	]
}
```

Response:
```json
{
		{
			"status_code": 404,
			"status": "failed",
			"url_name": "wfsfe.gg",
			"available": "down"
		},
		{
			"status_code": 200,
			"status": "done",
			"url_name": "google.com",
			"available": "up"
		}
	],
	"status": "done",
	"links_num": "fab2ea7e-b50c-4c25-bec3-355b3a7bd370"
}
```

2. POST `get-info-links` - получение информации о ссылках в pdf файле

Request:
```json
{
	"links_list": [
		"fab2ea7e-b50c-4c25-bec3-355b3a7bd370"
	]
}
```

Response:
```pdf
Link Check Report
=======================
https://ya.ru - up
https://wfsfe.gg - down
https://google.com - up
```

#№ Устойчивость к сбоям

##№ Сохранение прогресса

Все изменения записываются в `tasks.json` сразу после проверки ссылки.

##№ Продолжение незавершённых задач

При старте сервера все ссылки со статусом `"pending"` или `"in-progress"` ставятся обратно в очередь на проверку.

##№ Атомарная запись

Перед записью создаётся временный файл, который затем заменяет основной (`os.WriteFile` + `os.Rename`).

##№ Graceful shutdown (рекомендация)

Можно добавить обработку сигнала `SIGINT`/`SIGTERM` для корректного завершения всех воркеров.


