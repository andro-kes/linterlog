package linterlog_test

import (
	"testing"

	"github.com/andro-kes/linterlog"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, linterlog.Analyzer, "a")
}
