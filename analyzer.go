package analyzer

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var NilAnalyzer = &analysis.Analyzer{
	Name: "nil_error_struct",
	Doc:  "Nil Error Struct \n\n",
	Run:  run,
}

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
