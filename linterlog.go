package linterlog

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

const Doc = `linterlog is a linter for analyzing log statements in Go code

This linter checks log statements in your Go code against configurable rules.
You can define custom rules to ensure consistent logging practices across your codebase.

Example rules you can implement:
- Ensure log messages start with a capital letter
- Check for proper log level usage
- Validate log message format
- Ensure structured logging fields follow conventions
- Check for sensitive data in log messages
`

// Analyzer is the main analyzer for linterlog
var Analyzer = &analysis.Analyzer{
	Name: "linterlog",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Visit all call expressions to find logging calls
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Check if this is a logging function call
			if isLogCall(call) {
				checkLogCall(pass, call)
			}

			return true
		})
	}
	return nil, nil
}

// isLogCall checks if a call expression is a logging function call
// This is where you can define which functions are considered logging functions
func isLogCall(call *ast.CallExpr) bool {
	// Get the function selector
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Check for common logging function names
	// Extend this list based on your logging library
	logFuncs := map[string]bool{
		"Print":   true,
		"Println": true,
		"Printf":  true,
		"Fatal":   true,
		"Fatalf":  true,
		"Fatalln": true,
		"Panic":   true,
		"Panicf":  true,
		"Panicln": true,
		"Error":   true,
		"Errorf":  true,
		"Errorln": true,
		"Warn":    true,
		"Warnf":   true,
		"Warnln":  true,
		"Warning": true,
		"Info":    true,
		"Infof":   true,
		"Infoln":  true,
		"Debug":   true,
		"Debugf":  true,
		"Debugln": true,
		"Log":     true,
		"Logf":    true,
	}

	return logFuncs[selector.Sel.Name]
}

// checkLogCall performs checks on a logging call
// This is where you implement your custom rules
func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	// Example rule: Check if log call has at least one argument
	if len(call.Args) == 0 {
		pass.Reportf(call.Pos(), "log call should have at least one argument")
		return
	}

	// Example rule: Check if the first argument is a string literal
	// You can add more sophisticated checks here
	if lit, ok := call.Args[0].(*ast.BasicLit); ok {
		checkLogMessage(pass, call, lit)
	}
}

// checkLogMessage performs checks on the log message
// Implement your custom rules here
func checkLogMessage(pass *analysis.Pass, call *ast.CallExpr, lit *ast.BasicLit) {
	// Example: You can add rules like:
	// - Check message format
	// - Check for sensitive data patterns
	// - Ensure messages start with capital letter
	// - Check for proper structured logging format
	
	// This is a placeholder for your custom rules
	// Uncomment and modify the example below:
	
	/*
	message := lit.Value
	if len(message) > 2 && message[1] >= 'a' && message[1] <= 'z' {
		pass.Reportf(lit.Pos(), "log message should start with a capital letter")
	}
	*/
}
