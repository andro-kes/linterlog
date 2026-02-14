# Примеры пользовательских правил

Этот файл содержит примеры того, как можно расширить linterlog для проверки различных аспектов логирования.

## Текущая реализация / Current Implementation

linterlog в настоящее время реализует следующие правила:

linterlog currently implements the following rules:

### Правило 1: Lowercase First Letter

```go
// Rule 1: Check if message starts with uppercase letter
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
```

**✨ Особенность**: Это правило поддерживает автоматическое исправление (SuggestedFix).

**✨ Feature**: This rule supports auto-fix (SuggestedFix).

### Правило 2: English-Only Content

```go
func isEnglishOnly(message []rune) bool {
    for _, r := range message {
        if isASCIILetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
            continue
        }
        return false
    }
    return true
}

func isASCIILetter(r rune) bool {
    return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}
```

Сообщения могут содержать только английские буквы (A-Z, a-z), цифры (0-9) и пробелы.

Messages can only contain English letters (A-Z, a-z), digits (0-9), and spaces.

### Правило 3: No Special Symbols or Emojis

```go
func hasSpecialSymbols(message []rune) bool {
    for i, r := range message {
        // Check for ellipsis (...)
        if r == '.' && i+2 < len(message) && message[i+1] == '.' && message[i+2] == '.' {
            return true
        }
        // Check for specific special symbols
        if r == '!' || r == ':' || r == ';' {
            return true
        }
        // Check for non-ASCII symbols that are not letters
        if r > 127 && !unicode.IsLetter(r) {
            return true
        }
    }
    return false
}
```

Запрещены: `!`, `:`, `;`, `...` (многоточие), эмодзи и не-ASCII символы (кроме букв).

Forbidden: `!`, `:`, `;`, `...` (ellipsis), emojis, and non-ASCII symbols (except letters).

### Правило 4: Sensitive Data Detection

```go
func checkSensitiveData(msg string) error {
    alwaysSensitive := []string{
        "password",
        "api_key",
        "apikey",
        "secret",
        "credential",
    }
    lower := strings.ToLower(msg)
    for _, sd := range alwaysSensitive {
        if strings.Contains(lower, sd) {
            return errors.New("sensitive data was detected")
        }
    }
    return nil
}

// For concatenated strings, includes "token" pattern
func checkSensitiveDataWithToken(msg string) error {
    sensitiveData := []string{
        "password",
        "token",      // Only used for concatenation checks
        "api_key",
        "apikey",
        "secret",
        "credential",
    }
    lower := strings.ToLower(msg)
    for _, sd := range sensitiveData {
        if strings.Contains(lower, sd) {
            return errors.New("sensitive data was detected")
        }
    }
    return nil
}
```

**Примечание**: Функция `checkSensitiveData()` (для литералов) не включает `token`, тогда как `checkSensitiveDataWithToken()` (для конкатенаций) включает его, поскольку конкатенации вида "token" + значение считаются подозрительными.

**Note**: The `checkSensitiveData()` function (for literals) does not include `token`, while `checkSensitiveDataWithToken()` (for concatenations) includes it, since concatenations like "token" + value are considered suspicious.

---

## Примеры расширения / Extension Examples

Ниже приведены примеры того, как можно добавить дополнительные правила.

Below are examples of how to add additional rules.

## Пример 1: Проверка формата сообщений

```go
func checkLogMessage(pass *analysis.Pass, call *ast.CallExpr, lit *ast.BasicLit) {
    message := lit.Value
    if len(message) < 3 {
        return
    }
    
    // Удаляем кавычки
    msg := message[1:len(message)-1]
    
    // Дополнительное правило: сообщение не должно заканчиваться точкой
    // Additional rule: message should not end with a period
    if len(msg) > 0 && msg[len(msg)-1] == '.' {
        pass.Reportf(lit.Pos(), "log message should not end with a period")
    }
    
    // Дополнительное правило: максимальная длина сообщения
    // Additional rule: maximum message length
    if len(msg) > 100 {
        pass.Reportf(lit.Pos(), "log message is too long (max 100 characters)")
    }
}
```

## Пример 2: Расширенная проверка конфиденциальных данных

Пример расширения текущего правила:

```go
func checkExtendedSensitiveData(pass *analysis.Pass, lit *ast.BasicLit) {
    message := strings.ToLower(lit.Value)
    
    // Расширение текущего списка
    // Extension of current list
    sensitivePatterns := []string{
        // Current patterns:
        "password",
        "token",  // used in concatenation checks
        "secret",
        "api_key",
        "apikey",
        "credential",
        
        // Additional patterns:
        "auth",
        "bearer",
        "private_key",
        "access_key",
    }
    
    for _, pattern := range sensitivePatterns {
        if strings.Contains(message, pattern) {
            pass.Reportf(lit.Pos(), 
                "log message may contain sensitive data: %s", pattern)
        }
    }
}
```

## Пример 3: Проверка структурированного логирования (logrus/zap)

```go
func checkStructuredLogging(pass *analysis.Pass, call *ast.CallExpr) {
    // Проверка для logrus: log.WithFields(log.Fields{...}).Info(...)
    if len(call.Args) > 0 {
        // Проверяем, что используются WithFields
        selector, ok := call.Fun.(*ast.SelectorExpr)
        if !ok {
            return
        }
        
        // Проверяем, что перед методом логирования есть WithFields
        prevCall, ok := selector.X.(*ast.CallExpr)
        if ok {
            if prevSelector, ok := prevCall.Fun.(*ast.SelectorExpr); ok {
                if prevSelector.Sel.Name != "WithFields" && prevSelector.Sel.Name != "With" {
                    pass.Reportf(call.Pos(), 
                        "consider using structured logging with WithFields")
                }
            }
        }
    }
}
```

## Пример 4: Проверка уровней логирования в разных пакетах

```go
func checkLogLevel(pass *analysis.Pass, call *ast.CallExpr) {
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return
    }
    
    // Получаем имя пакета
    pkgPath := pass.Pkg.Path()
    
    // Правило: в production коде не должно быть Debug логов
    if strings.Contains(pkgPath, "/internal/") || strings.Contains(pkgPath, "/pkg/") {
        if selector.Sel.Name == "Debug" || selector.Sel.Name == "Debugf" {
            pass.Reportf(call.Pos(), 
                "Debug logs should not be used in production packages")
        }
    }
    
    // Правило: в тестах можно использовать только Info и Debug
    if strings.HasSuffix(pass.Fset.File(call.Pos()).Name(), "_test.go") {
        if selector.Sel.Name == "Error" || selector.Sel.Name == "Fatal" {
            pass.Reportf(call.Pos(), 
                "Error/Fatal logs should not be used in tests, use t.Error/t.Fatal instead")
        }
    }
}
```

## Пример 5: Проверка использования форматирования

```go
func checkFormatting(pass *analysis.Pass, call *ast.CallExpr) {
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return
    }
    
    funcName := selector.Sel.Name
    
    // Проверяем использование Printf с одним аргументом
    if strings.HasSuffix(funcName, "f") && len(call.Args) == 1 {
        pass.Reportf(call.Pos(), 
            "using formatted logging function with only format string, consider using non-formatted version")
    }
    
    // Проверяем использование Print с форматированием
    if !strings.HasSuffix(funcName, "f") && len(call.Args) > 0 {
        if lit, ok := call.Args[0].(*ast.BasicLit); ok {
            if strings.Contains(lit.Value, "%") {
                pass.Reportf(call.Pos(), 
                    "format specifiers found in non-formatted logging function, use Printf instead")
            }
        }
    }
}
```

## Пример 6: Проверка контекста

```go
func checkContext(pass *analysis.Pass, call *ast.CallExpr) {
    // Проверяем, что первый аргумент - context.Context для определенных функций
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return
    }
    
    // Для некоторых логгеров требуется передавать контекст
    contextRequiredFuncs := []string{"InfoContext", "ErrorContext", "DebugContext"}
    
    for _, funcName := range contextRequiredFuncs {
        if selector.Sel.Name == funcName {
            if len(call.Args) == 0 {
                pass.Reportf(call.Pos(), 
                    "%s requires context as first argument", funcName)
                return
            }
            
            // Проверяем тип первого аргумента
            // Это упрощенная проверка, для production кода нужна более точная проверка типа
            if ident, ok := call.Args[0].(*ast.Ident); ok {
                if ident.Name == "nil" {
                    pass.Reportf(call.Pos(), 
                        "context should not be nil")
                }
            }
        }
    }
}
```

## Пример 7: Расширение для разных библиотек логирования

```go
// Поддержка различных популярных библиотек
func isLogCall(call *ast.CallExpr) bool {
    selector, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return false
    }
    
    // Стандартная библиотека log
    stdLogFuncs := map[string]bool{
        "Print": true, "Printf": true, "Println": true,
        "Fatal": true, "Fatalf": true, "Fatalln": true,
        "Panic": true, "Panicf": true, "Panicln": true,
    }
    
    // logrus
    logrusLevels := map[string]bool{
        "Trace": true, "Debug": true, "Info": true,
        "Warn": true, "Warning": true, "Error": true,
        "Fatal": true, "Panic": true,
        "Tracef": true, "Debugf": true, "Infof": true,
        "Warnf": true, "Warningf": true, "Errorf": true,
        "Fatalf": true, "Panicf": true,
    }
    
    // zap
    zapLevels := map[string]bool{
        "Debug": true, "Info": true, "Warn": true, "Error": true,
        "DPanic": true, "Panic": true, "Fatal": true,
        "Debugw": true, "Infow": true, "Warnw": true, "Errorw": true,
        "DPanicw": true, "Panicw": true, "Fatalw": true,
    }
    
    // zerolog
    zerologLevels := map[string]bool{
        "Debug": true, "Info": true, "Warn": true, "Error": true,
        "Fatal": true, "Panic": true, "Trace": true,
        "Log": true, "Print": true,
    }
    
    funcName := selector.Sel.Name
    
    return stdLogFuncs[funcName] || 
           logrusLevels[funcName] || 
           zapLevels[funcName] || 
           zerologLevels[funcName]
}
```

## Интеграция примеров в ваш проект

Чтобы использовать эти примеры:

1. Скопируйте нужную функцию в `linterlog.go`
2. Вызовите её из функции `run()` или `checkLogCall()`
3. Добавьте соответствующие тесты в `testdata/src/a/a.go`
4. Запустите `make test` для проверки

Пример интеграции:

```go
func run(pass *analysis.Pass) (interface{}, error) {
    for _, file := range pass.Files {
        ast.Inspect(file, func(n ast.Node) bool {
            call, ok := n.(*ast.CallExpr)
            if !ok {
                return true
            }

            if isLogCall(call) {
                checkLogCall(pass, call)
                checkLogLevel(pass, call)        // Добавлено
                checkFormatting(pass, call)      // Добавлено
                checkContext(pass, call)         // Добавлено
            }

            return true
        })
    }
    return nil, nil
}
```


## AI Assistance Acknowledgment

Эта документация была создана и обновлена с использованием AI-помощников (GitHub Copilot и LLM Assistant) для ускорения процесса разработки и обеспечения полноты примеров.

This documentation was created and updated with the assistance of AI tools (GitHub Copilot and LLM Assistant) to accelerate the development process and ensure comprehensive examples.
