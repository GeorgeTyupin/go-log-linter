package main

import (
	"log"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/GeorgeTyupin/go-log-linter/internal/analyzer"
	"github.com/GeorgeTyupin/go-log-linter/internal/config"

	_ "go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	singlechecker.Main(analyzer.NewAnalyzer(cfg))
}
