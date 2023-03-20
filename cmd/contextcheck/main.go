package main

import (
	"github.com/SimonRichardson/context-linter/contextcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(contextcheck.Analyzer)
}
