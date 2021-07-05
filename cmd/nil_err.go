package main

import (
	nilerr "github.com/Jesse-Cameron/golang-nil-error-struct"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(nilerr.NilAnalyzer)
}
