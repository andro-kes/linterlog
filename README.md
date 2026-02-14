# linterlog

Тестовое задание для компании Selectel
Линтер для логов Go, совместимый с golangci-lint.

## Описание

linterlog - это настраиваемый линтер для анализа лог-записей в Go коде. Он разработан для интеграции с golangci-lint и позволяет вам определять собственные правила для проверки логирования в вашем проекте.

## Особенности

- ✅ Полная совместимость с golangci-lint
- ✅ Использует стандартный фреймворк go/analysis
- ✅ Легко настраивается и расширяется
- ✅ Поддерживает популярные библиотеки логирования (log, log/slog, zap и др.)
- ✅ Готовая структура для добавления пользовательских правил
- ✅ Автоматическое исправление (SuggestedFix) для правила lowercase-first-letter

## Реализованные правила проверки

linterlog проверяет лог-сообщения по следующим правилам (настраиваются через `config/config.yml`):

1. **Lowercase first letter**: Сообщения должны начинаться со строчной буквы
   - ✨ Поддерживает автоматическое исправление (SuggestedFix)
   
2. **English-only content**: Сообщения должны содержать только английские буквы, цифры и пробелы
   
3. **No special symbols or emojis**: Сообщения не должны содержать специальные символы
   - Запрещены: `!`, `:`, `;`, `...` (многоточие)
   - Запрещены: эмодзи и не-ASCII символы (если они не буквы)
   
4. **No sensitive data**: Сообщения не должны содержать чувствительные данные
   - Запрещенные паттерны: `password`, `token`, `api_key`, `apikey`, `secret`, `credential`
   - Для конкатенации строк используется расширенный список, включающий `token`

## Установка

### Как standalone инструмент

```bash
go install github.com/andro-kes/linterlog/cmd/linterlog@latest
```

### С использованием make

```bash
make install
```

### Как плагин для golangci-lint

1. Соберите плагин:
```bash
make plugin
```

2. Настройте `.golangci.yml`:
```yaml
linters-settings:
  custom:
    linterlog:
      path: ./bin/linterlog.so
      description: Linter for analyzing log statements
      original-url: github.com/andro-kes/linterlog

linters:
  enable:
    - linterlog
```

## Использование

### Standalone режим

```bash
# Анализ текущего пакета
linterlog ./...

# Анализ конкретного пакета
linterlog ./pkg/mypackage
```

### С golangci-lint

```bash
golangci-lint run
```

### С make

```bash
# Сборка бинарника
make build

# Запуск тестов
make test

# Запуск на примерах
make example

# Показать все доступные команды
make help
```

## Разработка пользовательских правил

Основная логика линтера находится в файле `linterlog.go`. Вот как добавить свои правила:

### 1. Определение лог-функций

Отредактируйте функцию `isLogCall()` чтобы добавить функции логирования из вашей библиотеки:

```go
func isLogCall(call *ast.CallExpr) bool {
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return false
    }

    // Добавьте свои функции логирования
    logFuncs := map[string]bool{
        "Info":  true,
        "Error": true,
        "Debug": true,
        // ... добавьте ваши функции
    }

    return logFuncs[selector.Sel.Name]
}
```

### 2. Добавление правил проверки

Реализуйте свои правила в функции `checkLogMessage()`:

```go
func checkLogMessage(pass *analysis.Pass, call *ast.CallExpr, lit *ast.BasicLit) {
    message := lit.Value
    
    // Пример правила: сообщение должно начинаться с заглавной буквы
    if len(message) > 2 && message[1] >= 'a' && message[1] <= 'z' {
        pass.Reportf(lit.Pos(), "log message should start with a capital letter")
    }
    
    // Пример: проверка на наличие конфиденциальных данных
    if strings.Contains(message, "password") {
        pass.Reportf(lit.Pos(), "log message should not contain sensitive data")
    }
    
    // Добавьте свои правила здесь
}
```

### 3. Создание тестов

Добавьте тестовые случаи в `testdata/src/a/a.go`:

```go
package a

import "log"

func Example() {
    log.Print() // want "log call should have at least one argument"
    log.Println("lowercase message") // want "log message should start with a capital letter"
    log.Println("Valid message") // OK
}
```

### 4. Запуск тестов

```bash
make test
```

## Структура проекта

```
.
├── linterlog.go              # Основная логика линтера
├── linterlog_test.go          # Тесты линтера
├── cmd/
│   └── linterlog/
│       └── main.go            # CLI точка входа
├── plugin/
│   └── plugin.go              # Плагин для golangci-lint
├── testdata/
│   └── src/a/
│       └── a.go               # Тестовые данные
├── Makefile                   # Автоматизация сборки
├── .golangci.yml              # Конфигурация golangci-lint
├── go.mod                     # Go модуль
└── README.md                  # Эта документация
```

## Примеры правил

Вот несколько идей для пользовательских правил:

### Формат сообщений
- Проверка, что сообщения начинаются с заглавной буквы
- Проверка, что сообщения не заканчиваются точкой
- Проверка максимальной длины сообщения

### Безопасность
- Обнаружение конфиденциальных данных (пароли, токены, ключи)
- Проверка на PII (personally identifiable information)

### Структурное логирование
- Проверка использования структурированных полей
- Валидация имен полей (например, snake_case или camelCase)
- Проверка типов значений полей

### Уровни логирования
- Проверка соответствия уровня логирования контексту
- Обеспечение использования правильных уровней в разных пакетах

### Производительность
- Обнаружение дорогих операций в логах (например, сериализация больших объектов)
- Проверка условного логирования для Debug уровня

## Поддерживаемые пакеты логирования

Линтер распознает следующие функции логирования:

- **Стандартная библиотека** (`log`, `log/slog`): `Print`, `Println`, `Printf`, `Fatal`, `Fatalf`, `Fatalln`, `Panic`, `Panicf`, `Panicln`
- **Уровни логирования**: `Error`, `Errorf`, `Errorln`, `Warn`, `Warnf`, `Warnln`, `Warning`, `Info`, `Infof`, `Infoln`, `Debug`, `Debugf`, `Debugln`, `Log`, `Logf`
- **zap, logrus и другие**: Любой логгер, чьи методы соответствуют указанным выше именам

## Конфигурация

Правила можно включать/отключать через файл `config/config.yml`:

```yaml
rules:
  capital_letter: true    # Проверка lowercase first letter
  only_english: true      # Проверка English-only content
  special_symbols: true   # Проверка на специальные символы
  sensitive_data: true    # Проверка на чувствительные данные
```

**Примечание**: Пользовательские паттерны для sensitive data пока не поддерживаются.

## Разработка

### Требования

- Go 1.22+ (тестируется на Go 1.22, 1.23, 1.24)
- make (опционально)

### Команды для разработки

```bash
# Установка зависимостей
make deps

# Сборка
make build

# Запуск тестов
make test

# Линтинг
make lint

# Очистка
make clean
```

## Интеграция с CI/CD

### GitHub Actions

Репозиторий использует GitHub Actions для CI/CD (см. `.github/workflows/ci.yml`):

- **Test**: Запуск тестов на Go 1.22, 1.23, 1.24
  - Включает race detector и coverage
  - Загрузка coverage в Codecov (для Go 1.23)
- **Build**: Сборка бинарника и smoke-test
- **Lint**: Запуск golangci-lint

Пример интеграции в ваш проект:

```yaml
name: Lint
on: [push, pull_request]
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Build plugin
        run: make plugin
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
```

## Статус выполнения задания

### Реализованные этапы

- ✅ **Структура проекта и analyzer**: Создана структура на базе `go/analysis`
- ✅ **Правила проверки**: Реализованы все 4 правила с конфигурацией
- ✅ **Тесты и testdata**: Полное покрытие тестами с использованием `analysistest`
- ✅ **Плагин для golangci-lint**: Готов к использованию плагин

### Текущие возможности

- Анализ литералов и конкатенации строк в лог-сообщениях
- Настраиваемые правила через YAML конфигурацию
- Автоматическое исправление для правила lowercase-first-letter
- Поддержка популярных библиотек логирования
- Интеграция с golangci-lint

## Лицензия

MIT

## Вклад

Приветствуются вклады! Пожалуйста, откройте issue или pull request.

## Автор

andro-kes

## AI Assistance Acknowledgment

Этот проект разработан с использованием AI-помощников:
- **GitHub Copilot**: Использовался для ускорения разработки кода
- **LLM Assistant**: Использовался для создания и обновления документации, а также для review pull requests

AI-инструменты помогли повысить продуктивность разработки, но весь код был проверен и протестирован вручную.

## Ссылки

- [go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) - Фреймворк для создания линтеров
- [golangci-lint](https://golangci-lint.run/) - Агрегатор линтеров для Go
- [Как писать линтеры](https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/)
