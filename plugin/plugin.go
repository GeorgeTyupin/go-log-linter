package main

import (
	"golang.org/x/tools/go/analysis"

	"github.com/GeorgeTyupin/go-log-linter/internal/analyzer"
	"github.com/GeorgeTyupin/go-log-linter/internal/config"
)

func New(_ any) ([]*analysis.Analyzer, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return []*analysis.Analyzer{analyzer.NewAnalyzer(cfg)}, nil
}

func main() {}
