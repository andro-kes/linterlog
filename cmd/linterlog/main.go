package main

import (
	"github.com/andro-kes/linterlog"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(linterlog.Analyzer)
}
