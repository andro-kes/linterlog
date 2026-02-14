package linterlog

import (
	"errors"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

const Doc = `linterlog - статический анализатор для Go кода
Правила:
	- Лог-сообщения должны начинаться со строчной буквы
	- Лог-сообщения должны быть только на английском языке
	- Лог-сообщения не должны содержать спецсимволы или эмодзи
	- Лог-сообщения не должны содержать чувствительные данные
`

var Analyzer = &analysis.Analyzer{
	Name: "linterlog",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if isLogCall(call) {
				checkLogCall(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

func isLogCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	logFuncs := map[string]bool{
		"Print": true, "Println": true, "Printf": true,
		"Fatal": true, "Fatalf": true, "Fatalln": true,
		"Panic": true, "Panicf": true, "Panicln": true,
		"Error": true, "Errorf": true, "Errorln": true,
		"Warn": true, "Warnf": true, "Warnln": true,
		"Warning": true,
		"Info": true, "Infof": true, "Infoln": true,
		"Debug": true, "Debugf": true, "Debugln": true,
		"Log": true, "Logf": true,
	}

	return logFuncs[selector.Sel.Name]
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}
	msgExpr := call.Args[0]

	switch e := msgExpr.(type) {
	case *ast.BasicLit:
		checkLogMessage(pass, e)

	case *ast.BinaryExpr:
		checkBinaryMessage(pass, e)

	default:
	}
}

func checkLogMessage(pass *analysis.Pass, lit *ast.BasicLit) {
	if lit.Kind != token.STRING {
		return
	}

	unquoted, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	message := []rune(unquoted)

	// Rule 1: Check if message starts with uppercase letter
	if len(message) > 0 && unicode.IsLetter(message[0]) && unicode.IsUpper(message[0]) {
		pass.Reportf(lit.Pos(), "log message should not start with a capital letter")
	}

	// Rule 3: Check for special symbols first (!, :, ;, ..., emojis)
	if hasSpecialSymbols(message) {
		pass.Reportf(lit.Pos(), "log message should not contain special symbols or emojis")
		return
	}

	// Rule 2: Check for English-only + digits
	if !isEnglishOnly(message) {
		pass.Reportf(lit.Pos(), "log message should contain only english symbols")
		return
	}

	// Rule 4: Check for sensitive data
	if err := checkSensitiveData(unquoted); err != nil {
		pass.Reportf(lit.Pos(), "log message should not contain sensitive data")
	}
}

func checkBinaryMessage(pass *analysis.Pass, expr *ast.BinaryExpr) {
	parts := extractStringLiteralsFromConcat(expr)
	joined := strings.Join(parts, "")
	message := []rune(joined)

	// Rule 1: Check if message starts with uppercase letter
	if len(message) > 0 && unicode.IsLetter(message[0]) && unicode.IsUpper(message[0]) {
		pass.Reportf(expr.Pos(), "log message should not start with a capital letter")
	}

	// Rule 3: Check for special symbols first (!, :, ;, ..., emojis)
	if hasSpecialSymbols(message) {
		pass.Reportf(expr.Pos(), "log message should not contain special symbols or emojis")
		return
	}

	// Rule 2: Check for English-only + digits
	if !isEnglishOnly(message) {
		pass.Reportf(expr.Pos(), "log message should contain only english symbols")
		return
	}

	// Rule 4: Check for sensitive data in parts and joined message
	for _, p := range parts {
		if err := checkSensitiveDataWithToken(p); err != nil {
			pass.Reportf(expr.Pos(), "log message should not contain sensitive data")
			return
		}
	}

	if err := checkSensitiveDataWithToken(joined); err != nil {
		pass.Reportf(expr.Pos(), "log message should not contain sensitive data")
		return
	}
}

func extractStringLiteralsFromConcat(e ast.Expr) []string {
	var out []string

	var walk func(ast.Expr)
	walk = func(x ast.Expr) {
		switch n := x.(type) {
		case *ast.BinaryExpr:
			if n.Op == token.ADD {
				walk(n.X)
				walk(n.Y)
			}
		case *ast.ParenExpr:
			walk(n.X)
		case *ast.BasicLit:
			if n.Kind == token.STRING {
				if unq, err := strconv.Unquote(n.Value); err == nil {
					out = append(out, unq)
				}
			}
		}
	}

	walk(e)
	return out
}

func isASCIILetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

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

// checkSensitiveDataWithToken checks for sensitive data including "token" pattern
// Used for concatenated strings where "token" + value is suspicious
func checkSensitiveDataWithToken(msg string) error {
	sensitiveData := []string{
		"password",
		"token",
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

func isEnglishOnly(message []rune) bool {
	for _, r := range message {
		if isASCIILetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		return false
	}
	return true
}

func hasSpecialSymbols(message []rune) bool {
	for i, r := range message {
		if r == '.' && i+2 < len(message) && message[i+1] == '.' && message[i+2] == '.' {
			return true
		}
		if r == '!' || r == ':' || r == ';' {
			return true
		}
		if r > 127 && !unicode.IsLetter(r) {
			return true
		}
	}
	return false
}