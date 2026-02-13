package main

import (
	"github.com/andro-kes/linterlog"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin is required by golangci-lint
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

// GetAnalyzers returns the analyzers to be used by golangci-lint
func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		linterlog.Analyzer,
	}
}
