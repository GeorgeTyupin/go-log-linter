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

## Сборка

```bash
go build -buildmode=plugin -o plugin/loglinter.so plugin/plugin.go
```

## Интеграция с golangci-lint

Добавить в `.golangci.yml`:

```yaml
linters:
  settings:
    custom:
      loglinter:
        path: ./plugin/loglinter.so
        description: Checks log messages for style violations
        original-url: github.com/GeorgeTyupin/go-log-linter
```

## Запуск без golangci-lint

```bash
go run ./cmd/linter ./...
```

## Тесты

```bash
go test ./...
```
