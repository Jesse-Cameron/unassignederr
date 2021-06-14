package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Jesse-Cameron/golang-nil-err-pointer/nilerr"
)

func main() {
	singlechecker.Main(nilerr.NilAnalyzer)
}
