package main

import (
	"golang.org/x/tools/go/analysis"

	"github.com/GeorgeTyupin/go-log-linter/internal/analyzer"
)

func New(_ any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func main() {}
