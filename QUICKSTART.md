# Быстрый старт / Quick Start

## Русский

### Требования

- Go 1.22 или выше (тестируется на Go 1.22, 1.23, 1.24)
- make (опционально)

### Шаг 1: Клонирование репозитория
```bash
git clone https://github.com/andro-kes/linterlog.git
cd linterlog
```

### Шаг 2: Установка зависимостей
```bash
make deps
```

### Шаг 3: Запуск тестов
```bash
make test
```

### Шаг 4: Сборка
```bash
make build
```

### Шаг 5: Запуск на примере
```bash
make example
```

Вы увидите вывод:
```
Running linter on testdata...
/path/to/testdata/src/a/a.go:7:2: log call should have at least one argument
```

### Шаг 6: Добавление собственных правил

1. Откройте `linterlog.go`
2. Найдите функцию `checkLogMessage()`
3. Раскомментируйте пример или добавьте свои правила
4. Запустите `make test` для проверки
5. Запустите `make build && make example` для тестирования

### Примеры правил

См. файл `EXAMPLES.md` для подробных примеров различных правил.

---

## English

### Requirements

- Go 1.22 or higher (tested on Go 1.22, 1.23, 1.24)
- make (optional)

### Step 1: Clone the repository
```bash
git clone https://github.com/andro-kes/linterlog.git
cd linterlog
```

### Step 2: Install dependencies
```bash
make deps
```

### Step 3: Run tests
```bash
make test
```

### Step 4: Build
```bash
make build
```

### Step 5: Run on example
```bash
make example
```

You will see output like:
```
Running linter on testdata...
/path/to/testdata/src/a/a.go:7:2: log call should have at least one argument
```

### Step 6: Add your own rules

1. Open `linterlog.go`
2. Find the `checkLogMessage()` function
3. Add your custom rules (all 4 rules are already implemented)
4. Run `make test` to verify
5. Run `make build && make example` to test

### Current Implemented Rules

См. `README.md` для подробной информации о текущих правилах:
- Lowercase first letter (с автоисправлением / with auto-fix)
- English-only content
- No special symbols or emojis
- Sensitive data detection

See `README.md` for detailed information about current rules.

---

## Интеграция с golangci-lint / Integration with golangci-lint

### Русский

1. Соберите плагин:
```bash
make plugin
```

2. Добавьте в `.golangci.yml` вашего проекта:
```yaml
linters-settings:
  custom:
    linterlog:
      path: /path/to/linterlog/bin/linterlog.so
      description: Linter for log statements
      original-url: github.com/andro-kes/linterlog

linters:
  enable:
    - linterlog
```

3. Запустите:
```bash
golangci-lint run
```

### English

1. Build the plugin:
```bash
make plugin
```

2. Add to your project's `.golangci.yml`:
```yaml
linters-settings:
  custom:
    linterlog:
      path: /path/to/linterlog/bin/linterlog.so
      description: Linter for log statements
      original-url: github.com/andro-kes/linterlog

linters:
  enable:
    - linterlog
```

3. Run:
```bash
golangci-lint run
```

---

## Доступные команды / Available Commands

Запустите `make help` для просмотра всех доступных команд.

Run `make help` to see all available commands:

```bash
$ make help
Usage: make [target]

Available targets:
  build           Build the standalone linter binary
  plugin          Build the golangci-lint plugin
  test            Run tests
  lint            Run linter on the project itself
  clean           Clean build artifacts
  install         Install the linter to $GOPATH/bin
  deps            Download dependencies
  example         Run the linter on example code
  help            Show this help message
```


---

## AI Assistance Acknowledgment

Эта документация была создана с использованием AI-помощников (GitHub Copilot и LLM Assistant).

This documentation was created with the assistance of AI tools (GitHub Copilot and LLM Assistant).
