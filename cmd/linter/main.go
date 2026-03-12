package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/GeorgeTyupin/go-log-linter/internal/analyzer"

	_ "go.uber.org/zap"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
