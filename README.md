# go-log-linter

Линтер для Go, который проверяет лог-сообщения на соответствие правилам. Совместим с `golangci-lint`.

## Правила

| Правило | Пример нарушения |
|---|---|
| Сообщение начинается со строчной буквы | `slog.Info("Starting server")` |
| Только английский язык | `slog.Info("запуск сервера")` |
| Нет спецсимволов и эмодзи | `slog.Info("done! 🚀")` |
| Нет чувствительных данных | `slog.Info("token: " + t)` |

Поддерживаемые логгеры: `log`, `log/slog`, `go.uber.org/zap`.

## Сборка и запуск

### 1. Как плагин для golangci-lint (Рекомендуется)

Сначала соберите плагин:
```bash
go build -buildmode=plugin -o plugin/loglinter.so plugin/plugin.go
```

Запуск проверки:
```bash
golangci-lint run ./testdata/sample.go
```
*Примечание: Ваш `golangci-lint` должен быть собран той же версией Go, что и плагин. Если они не совпадают, используйте `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.0`.*

#### Настройка .golangci.yml
Линтер уже настроен в проекте, но для справки — вот что добавлено в `.golangci.yml`:
```yaml
linters:
  enable:
    - loglinter
  custom:
    loglinter:
      path: ./plugin/loglinter.so
      description: Checks log messages for style violations
```

### 2. Без golangci-lint
Запуск напрямую через `go run`:
```bash
go run ./cmd/linter ./testdata/
```
*(Мы указываем путь `./testdata/` явно. Вы можете указать любой другой путь с нуждающимися в проверке файлами)*

## Тесты
```bash
go test ./...
```
