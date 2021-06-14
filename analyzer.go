package nilerr

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var NilAnalyzer = &analysis.Analyzer{
	Name: "nil_error_struct",
	Doc:  "Nil Error Struct\n\n",
	Run:  run,
}

var errTyp = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			fd, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if !funcReturnsError(pass, fd.Type.Results.List) {
				return false
			}

			errorStatementIndex := 0
			errorStatementName := ""

			funcStmts := fd.Body.List
			for i, stmt := range funcStmts {

				// if there is a var statement
				genDecl, found := isVarStatment(stmt)
				if !found {
					continue
				}

				for _, spec := range genDecl.Specs {
					foundName, found := isErrorStruct(pass, spec)
					if !found {
						return true
					}

					errorStatementName = foundName
					errorStatementIndex = i + 1
				}

				// this is sequential and weird
				for i := errorStatementIndex; i <= len(funcStmts); i++ {
					if retStmt, ok := funcStmts[i].(*ast.ReturnStmt); ok {
						if errorStatementPos, ok := returnsInvalidError(pass, errorStatementName, retStmt); ok {
							// we found an invalid return
							pass.Reportf(errorStatementPos.NamePos, "uninitialized custom error returned %q",
								render(pass.Fset, errorStatementPos))
							return false
						}
					}
				}

			}

			return true
		})
	}

	return nil, nil
}

func isVarStatment(stmt ast.Stmt) (*ast.GenDecl, bool) {
	declStmt, ok := stmt.(*ast.DeclStmt)
	if !ok {
		return nil, false
	}

	genDecl, ok := declStmt.Decl.(*ast.GenDecl)
	if !ok {
		return nil, false

	}

	if genDecl.Tok == token.VAR {
		return genDecl, true
	}

	return nil, false
}

func funcReturnsError(pass *analysis.Pass, list []*ast.Field) bool {
	isError := func(v ast.Expr) bool {
		if n, ok := pass.TypesInfo.TypeOf(v).(*types.Named); ok {
			o := n.Obj()
			return o != nil && o.Pkg() == nil && o.Name() == "error"
		}
		return false
	}

	for _, returnArg := range list {
		if isError(returnArg.Type) {
			return true
		}
	}

	return false
}

func returnsInvalidError(pass *analysis.Pass, name string, returnStmt *ast.ReturnStmt) (*ast.Ident, bool) {
	for _, result := range returnStmt.Results {
		if identifier, ok := result.(*ast.Ident); ok {
			if identifier.Name == name {
				return identifier, true
			}
		}
	}
	return nil, false
}

func isErrorStruct(pass *analysis.Pass, spec ast.Spec) (string, bool) {
	v, ok := spec.(*ast.ValueSpec)
	if !ok {
		return "", false
	}

	if v.Type == nil {
		return "", false
	}

	// if the initial values are non-nil
	if v.Values != nil {
		return "", false
	}

	// TODO: check if it not an error first
	p, ok := pass.TypesInfo.TypeOf(v.Type).Underlying().(*types.Pointer)
	if !ok {
		return "", false
	}

	// check if it's implements the error type
	if n, ok := p.Elem().(*types.Named); ok {
		if !types.Implements(n, errTyp) {
			return "", false
		}
	}

	// not sure why there is more than one name here
	return v.Names[0].Name, true
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
