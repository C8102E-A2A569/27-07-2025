## Download files and archive them into ZIP

**gRPC-сервис для скачивания файлов по URL и упаковки их в ZIP архив**

Функциональность: 
- CreateTask -- создание задачи, в которую потом можно добавить файлы
- AddFiles -- скачивает файлы по URL в заданном формате(.pdf, .jpeg, .jpg) и добавляет их в задачу
- GetTaskStatus -- получить статус задачи: created, processing, completed, failed
- DownloadArchive -- выдает готовый ZIP архив

Ограничения: 
- возможно добавить только 3 файла в задачу
- одновременно обрабатываться могут только 3 задачи
- расширение файлов: .pdf, .jpeg(.jpg)

Структура: 
- `cmd/main.go` — запуск gRPC-сервера
- `internal/config` — конфигурация сервиса
- `internal/model` — логика архивации, скачивания, структура task
- `internal/service` — реализация gRPC-сервиса(основные методы системы)
- `proto/archive.proto`
- `test/integration-test.go` — интеграционный тест

Запуск: 
1. `go run cmd/main.go` 
2. `go run test/integration-test.go`