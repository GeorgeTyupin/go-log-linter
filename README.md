# go-log-linter

[![CI](https://github.com/GeorgeTyupin/go-log-linter/actions/workflows/ci.yml/badge.svg)](https://github.com/GeorgeTyupin/go-log-linter/actions/workflows/ci.yml)

Линтер для Go, который проверяет лог-сообщения на соответствие правилам стиля и безопасности. Полностью совместим с `golangci-lint` как плагин.

## Правила

| Правило | Пример нарушения | Авто-исправление |
|---|---|:---:|
| Сообщение со строчной буквы | `slog.Info("Starting server")` | ✅ |
| Только английский язык | `slog.Info("запуск сервера")` | ❌ |
| Нет спецсимволов и эмодзи | `slog.Info("done! 🚀")` | ✅ |
| Нет чувствительных данных | `slog.Info("password: " + p)` | ❌ |

Поддерживаемые логгеры: `log`, `log/slog`, `go.uber.org/zap`.

## Конфигурация

Все настройки линтера теперь хранятся в центральном файле [configs/config.yaml](configs/config.yaml). Это позволяет один раз настроить правила для всех способов запуска.

### Параметры configs/config.yaml:
```yaml
disable-lowercase: false       # Отключить проверку строчной буквы
disable-english-only: false     # Отключить проверку английского языка
disable-no-special-chars: false  # Отключить проверку спецсимволов
disable-no-sensitive-data: false # Отключить проверку секретов
extra-sensitive-patterns: "mytoken,cvv" # Свои паттерны для поиска секретов
```

## Сборка и запуск

### 1. Как плагин для golangci-lint (Рекомендуется)

Сначала соберите плагин:
```bash
go build -buildmode=plugin -o plugin/loglinter.so plugin/plugin.go
```

Запуск проверки:
```bash
# Обычный запуск (не забудьте очистить кэш при смене конфига!)
golangci-lint cache clean && golangci-lint run ./path/to/your/code

# С автоматическим исправлением ошибок
golangci-lint run --fix ./path/to/your/code
```

#### Настройка .golangci.yml:
```yaml
linters:
  enable:
    - loglinter
linters-settings:
  custom:
    loglinter:
      path: ./plugin/loglinter.so
      description: Checks log messages for style violations
```

### 2. Автономный режим (без golangci-lint)

Вы можете запустить линтер как обычную консольную утилиту. Настройки в этом случае также берутся из [configs/config.yaml](configs/config.yaml).

#### Запуск через go run:
```bash
# Обычная проверка
go run ./cmd/linter ./path/to/your/code

# Применить авто-исправления
go run ./cmd/linter -fix ./path/to/your/code
```

#### Сборка и запуск бинарного файла:
```bash
# 1. Собрать бинарник
go build -o loglinter ./cmd/linter

# 2. Запустить проверку
./loglinter ./path/to/your/code

# 3. Применить авто-исправления
./loglinter -fix ./path/to/your/code
```

## Авто-исправление (SuggestedFixes)

Линтер поддерживает умное авто-исправление для правил **lowercase** и **спецсимволы**.

- В **golangci-lint** используйте флаг `--fix`.
- В **автономном режиме** (go run / бинарник) используйте флаг `-fix`.

> [!TIP]
> В современных IDE (VS Code, GoLand) исправления доступны через меню **Quick Fix** (лампочка) прямо во время написания кода.

## Разработка и тесты

Для запуска тестов используйте стандартную команду Go:
```bash
go test ./...
```

Проект настроен на автоматическое тестирование через **GitHub Actions** при каждом пуше.
