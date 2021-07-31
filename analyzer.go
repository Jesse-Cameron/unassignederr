package analyzer

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var UnassignedErrAnalyzer = &analysis.Analyzer{
	Name: "unassignederr",
	Doc:  DocString,
	Run:  run,
}

const DocString = "A tool for identifying when a uninitialised error struct is being incorrectly returned.\n\nunassignederr checks for functions where there is a returned an error struct that hasn't been initialized. Usually this can be resolved by ensuring that the variable is initialiased. Either in the same statement, or in subsequent lines."

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			return parseNode(pass, node)
		})
	}

	return nil, nil
}

// parseNode checks a node to see if it contains a given error
func parseNode(pass *analysis.Pass, node ast.Node) bool {
	// if is a function declaration
	fd, ok := node.(*ast.FuncDecl)
	if !ok {
		return true
	}

	if !funcReturnsError(pass, fd.Type.Results.List) {
		return true
	}

	err, found := findUnassignedErrors(pass, fd)
	if found {
		pass.Reportf(err.NamePos, "uninitialized custom error returned %q",
			render(pass.Fset, err))
		return false
	}

	return true
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
