# Developer Guide - Добавление пользовательских правил / Adding Custom Rules

Этот руководство поможет вам добавить собственные правила для проверки логирования.

This guide will help you add your own rules for log checking.

## Текущая реализация / Current Implementation

linterlog в настоящее время реализует 4 правила:

linterlog currently implements 4 rules:

1. **Lowercase first letter** (capital_letter): Сообщения должны начинаться со строчной буквы
   - ✨ Поддерживает SuggestedFix для автоматического исправления
2. **English-only content** (only_english): Только английские буквы, цифры и пробелы
3. **No special symbols** (special_symbols): Запрет `!`, `:`, `;`, `...`, эмодзи и не-ASCII символов (кроме букв)
4. **Sensitive data** (sensitive_data): Обнаружение `password`, `token`, `api_key`, `apikey`, `secret`, `credential`

Все правила настраиваются через `config/config.yml`.

All rules are configurable via `config/config.yml`.

## Архитектура / Architecture

Линтер построен на базе фреймворка `go/analysis` и состоит из нескольких ключевых компонентов:

The linter is built on the `go/analysis` framework and consists of several key components:

```
linterlog.go
├── Analyzer          - главный анализатор / main analyzer
├── run()             - точка входа / entry point
├── isLogCall()       - определяет лог-вызовы / identifies log calls
├── checkLogCall()    - проверяет лог-вызов / checks log call
└── checkLogMessage() - проверяет сообщение / checks message
```

## Шаг 1: Определение функций логирования / Step 1: Define Logging Functions

Текущая реализация `isLogCall()` распознает следующие методы:

Current `isLogCall()` implementation recognizes these methods:

```go
func isLogCall(call *ast.CallExpr) bool {
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return false
    }

    // Поддерживаемые функции логирования
    // Supported logging functions
    logFuncs := map[string]bool{
        "Print": true, "Println": true, "Printf": true,
        "Fatal": true, "Fatalf": true, "Fatalln": true,
        "Panic": true, "Panicf": true, "Panicln": true,
        "Error": true, "Errorf": true, "Errorln": true,
        "Warn": true, "Warnf": true, "Warnln": true,
        "Warning": true,
        "Info":    true, "Infof": true, "Infoln": true,
        "Debug": true, "Debugf": true, "Debugln": true,
        "Log": true, "Logf": true,
    }

    return logFuncs[selector.Sel.Name]
}
```

Эти методы работают для стандартной библиотеки `log`, `log/slog`, `zap`, `logrus` и других.

These methods work for standard library `log`, `log/slog`, `zap`, `logrus`, and others.

Чтобы добавить новые функции, просто добавьте их в map `logFuncs`.

To add new functions, simply add them to the `logFuncs` map.

## Шаг 2: Добавление правил проверки / Step 2: Add Validation Rules

### Текущая реализация / Current Implementation

Функция `checkLogMessage()` реализует 4 правила:

The `checkLogMessage()` function implements 4 rules:

```go
func checkLogMessage(pass *analysis.Pass, lit *ast.BasicLit, cfg *config.Config) {
    // ... извлечение сообщения / extract message ...
    message := []rune(unquoted)

    // Rule 1: Lowercase first letter (с SuggestedFix / with SuggestedFix)
    if cfg.Rules.CapitalLetter {
        if len(message) > 0 && unicode.IsLetter(message[0]) && unicode.IsUpper(message[0]) {
            message[0] = unicode.ToLower(message[0])
            pass.Report(
                analysis.Diagnostic{
                    Pos:     lit.Pos(),
                    End:     lit.End(),
                    Message: "log message should not start with a capital letter",
                    SuggestedFixes: []analysis.SuggestedFix{
                        {
                            Message: "lowercase first letter",
                            TextEdits: []analysis.TextEdit{
                                {
                                    Pos:     lit.Pos(),
                                    End:     lit.End(),
                                    NewText: []byte(string(message)),
                                },
                            },
                        },
                    },
                },
            )
        }
    }

    // Rule 3: Check for special symbols (!, :, ;, ..., emojis)
    if cfg.Rules.SpecialSymbols {
        if hasSpecialSymbols(message) {
            pass.Reportf(lit.Pos(), "log message should not contain special symbols or emojis")
            return
        }
    }

    // Rule 2: English-only + digits + spaces
    if cfg.Rules.OnlyEnglish {
        if !isEnglishOnly(message) {
            pass.Reportf(lit.Pos(), "log message should contain only english symbols")
            return
        }
    }

    // Rule 4: Sensitive data
    if cfg.Rules.SensitiveData {
        if err := checkSensitiveData(unquoted); err != nil {
            pass.Reportf(lit.Pos(), "log message should not contain sensitive data")
        }
    }
}
```

### Добавление нового правила / Adding a New Rule

Добавьте проверку в `checkLogMessage()`:

Add a check to `checkLogMessage()`:

```go
func checkLogMessage(pass *analysis.Pass, call *ast.CallExpr, lit *ast.BasicLit) {
    message := lit.Value
    if len(message) < 3 {
        return
    }
    
    // Удаляем кавычки / Remove quotes
    msg := message[1:len(message)-1]
    
    // Правило: сообщение должно начинаться с заглавной буквы
    // Rule: message should start with capital letter
    if len(msg) > 0 && msg[0] >= 'a' && msg[0] <= 'z' {
        pass.Reportf(lit.Pos(), "log message should start with a capital letter")
    }
}
```

### Сложное правило / Complex Rule

Для более сложных проверок создайте отдельную функцию:

For more complex checks, create a separate function:

```go
func checkForSensitiveData(pass *analysis.Pass, call *ast.CallExpr) {
    // Проверяем все аргументы / Check all arguments
    for _, arg := range call.Args {
        // Обрабатываем строковые литералы / Handle string literals
        if lit, ok := arg.(*ast.BasicLit); ok {
            msg := strings.ToLower(lit.Value)
            
            sensitivePatterns := []string{
                "password", "token", "secret", "api_key",
            }
            
            for _, pattern := range sensitivePatterns {
                if strings.Contains(msg, pattern) {
                    pass.Reportf(lit.Pos(), 
                        "log message may contain sensitive data: %s", pattern)
                }
            }
        }
        
        // Обрабатываем переменные / Handle variables
        if ident, ok := arg.(*ast.Ident); ok {
            varName := strings.ToLower(ident.Name)
            if strings.Contains(varName, "password") {
                pass.Reportf(ident.Pos(), 
                    "avoid logging variables that may contain sensitive data")
            }
        }
    }
}
```

## Шаг 3: Интеграция правила / Step 3: Integrate the Rule

Добавьте вызов вашей функции в `run()` или `checkLogCall()`:

Add a call to your function in `run()` or `checkLogCall()`:

```go
func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
    if len(call.Args) == 0 {
        pass.Reportf(call.Pos(), "log call should have at least one argument")
        return
    }

    if lit, ok := call.Args[0].(*ast.BasicLit); ok {
        checkLogMessage(pass, call, lit)
    }
    
    // Добавьте ваши проверки / Add your checks
    checkForSensitiveData(pass, call)
}
```

## Шаг 4: Создание тестов / Step 4: Create Tests

Добавьте тестовые случаи в `testdata/src/a/a.go`:

Add test cases to `testdata/src/a/a.go`:

```go
package a

import "log"

func ExampleSensitiveData() {
    password := "secret123"
    log.Println(password) // want "avoid logging variables that may contain sensitive data"
    
    log.Println("password: abc") // want "log message may contain sensitive data: password"
    
    log.Println("Valid log message") // OK
}
```

Формат комментариев `// want "..."` указывает ожидаемые сообщения об ошибках.

The `// want "..."` comment format specifies expected error messages.

## Шаг 5: Запуск тестов / Step 5: Run Tests

```bash
make test
```

Если тест не проходит, проверьте:
- Совпадает ли сообщение об ошибке с `want`
- Правильно ли определяется позиция ошибки
- Корректно ли работает логика проверки

If the test fails, check:
- Does the error message match the `want` comment?
- Is the error position correctly identified?
- Does the validation logic work correctly?

## Шаг 6: Отладка / Step 6: Debugging

Используйте флаг `-debug` для детальной информации:

Use the `-debug` flag for detailed information:

```bash
./bin/linterlog -debug=fpstv ./testdata/src/a
```

Добавьте временные `fmt.Printf` для отладки:

Add temporary `fmt.Printf` for debugging:

```go
func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
    fmt.Printf("Checking call at %s\n", pass.Fset.Position(call.Pos()))
    fmt.Printf("Args count: %d\n", len(call.Args))
    // ... ваш код
}
```

## Работа с AST / Working with AST

### Основные типы узлов / Basic Node Types

```go
// Литералы / Literals
*ast.BasicLit       // "string", 123, true

// Идентификаторы / Identifiers
*ast.Ident          // переменные, функции / variables, functions

// Вызовы функций / Function calls
*ast.CallExpr       // func(args...)

// Селекторы / Selectors
*ast.SelectorExpr   // pkg.Function, obj.Method

// Бинарные операции / Binary operations
*ast.BinaryExpr     // a + b, x == y

// Композитные литералы / Composite literals
*ast.CompositeLit   // []int{1, 2}, struct{}{field: value}
```

### Обход AST / AST Traversal

```go
ast.Inspect(file, func(n ast.Node) bool {
    // Возвращает true для продолжения обхода
    // Returns true to continue traversal
    
    switch node := n.(type) {
    case *ast.CallExpr:
        // Обработка вызова функции
        // Handle function call
        
    case *ast.FuncDecl:
        // Обработка объявления функции
        // Handle function declaration
        
    case *ast.IfStmt:
        // Обработка if-оператора
        // Handle if statement
    }
    
    return true
})
```

### Получение информации о типах / Getting Type Information

```go
func checkWithTypeInfo(pass *analysis.Pass, call *ast.CallExpr) {
    // Получить тип выражения / Get expression type
    typeInfo := pass.TypesInfo.TypeOf(call.Fun)
    
    // Получить объект (функцию, переменную) / Get object (function, variable)
    if ident, ok := call.Fun.(*ast.Ident); ok {
        obj := pass.TypesInfo.ObjectOf(ident)
        if obj != nil {
            fmt.Printf("Object: %s, Type: %s\n", obj.Name(), obj.Type())
        }
    }
}
```

## Примеры реальных правил / Real-World Rule Examples

### Правило 1: Запрет fmt.Println в production коде / Rule 1: Forbid fmt.Println in production

```go
func checkNoFmtPrintln(pass *analysis.Pass, call *ast.CallExpr) {
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return
    }
    
    // Проверяем, что это fmt.Println
    if ident, ok := selector.X.(*ast.Ident); ok {
        if ident.Name == "fmt" && selector.Sel.Name == "Println" {
            // Проверяем, что это не тестовый файл
            filename := pass.Fset.Position(call.Pos()).Filename
            if !strings.HasSuffix(filename, "_test.go") {
                pass.Reportf(call.Pos(), 
                    "use proper logger instead of fmt.Println")
            }
        }
    }
}
```

### Правило 2: Обязательный контекст для ошибок / Rule 2: Required context for errors

```go
func checkErrorLogging(pass *analysis.Pass, call *ast.CallExpr) {
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return
    }
    
    // Только для Error и Fatal уровней
    if selector.Sel.Name != "Error" && selector.Sel.Name != "Fatal" {
        return
    }
    
    // Проверяем наличие error в аргументах
    hasError := false
    for _, arg := range call.Args {
        typeInfo := pass.TypesInfo.TypeOf(arg)
        if typeInfo != nil && typeInfo.String() == "error" {
            hasError = true
            break
        }
    }
    
    if !hasError {
        pass.Reportf(call.Pos(), 
            "error-level log should include error object")
    }
}
```

### Правило 3: Структурированные поля / Rule 3: Structured fields

```go
func checkStructuredLogging(pass *analysis.Pass, call *ast.CallExpr) {
    // Для логгеров типа logrus или zap
    // For loggers like logrus or zap
    
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return
    }
    
    // Проверяем метод логирования
    if selector.Sel.Name == "Info" || selector.Sel.Name == "Error" {
        // Проверяем, что перед ним есть WithFields
        if _, ok := selector.X.(*ast.CallExpr); !ok {
            pass.Reportf(call.Pos(), 
                "consider using structured logging with WithFields")
        }
    }
}
```

## Полезные ссылки / Useful Links

- [go/analysis documentation](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [go/ast documentation](https://pkg.go.dev/go/ast)
- [go/types documentation](https://pkg.go.dev/go/types)
- [AST Explorer for Go](https://yuroyoro.github.io/goast-viewer/)
- [Writing Custom Linters](https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/)

## Поддержка / Support

Если у вас возникли вопросы:
1. Проверьте файл `EXAMPLES.md` для примеров
2. Изучите существующий код в `linterlog.go`
3. Запустите тесты с флагом `-v` для детальной информации
4. Откройте issue в GitHub репозитории

If you have questions:
1. Check `EXAMPLES.md` file for examples
2. Study existing code in `linterlog.go`
3. Run tests with `-v` flag for detailed information
4. Open an issue in the GitHub repository

## AI Assistance Acknowledgment

Эта документация была создана и обновлена с использованием AI-помощников (GitHub Copilot и LLM Assistant) для ускорения процесса разработки.

This documentation was created and updated with the assistance of AI tools (GitHub Copilot and LLM Assistant) to accelerate the development process.
