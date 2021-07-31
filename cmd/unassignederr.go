package main

import (
	analyzer "github.com/Jesse-Cameron/unassignederr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.UnassignedErrAnalyzer)
}
