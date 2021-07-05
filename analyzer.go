package analyzer

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
	Doc:  "Nil Error Struct \n\n",
	Run:  run,
}

var errTyp = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

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

	errorStatementNames := []string{}
	errorStatementFound := false

	// search for an error declaration
	funcStmts := fd.Body.List
	for i, stmt := range funcStmts {
		// if we've found an error
		if errorStatementFound {
			//	search for the returned error

			if retStmt, ok := funcStmts[i].(*ast.ReturnStmt); ok {
				if errorStatementPos, ok := returnsInvalidError(pass, errorStatementNames, retStmt); ok {
					// we found an invalid return
					pass.Reportf(errorStatementPos.NamePos, "uninitialized custom error returned %q",
						render(pass.Fset, errorStatementPos))
					return false
				}
			}

		} else {
			//	search for an error
			// if there is a var statement
			genDecl, found := isVarStatment(stmt)
			if !found {
				continue
			}

			// if the var statement declares an error struct
			for _, spec := range genDecl.Specs {
				foundName, found := isErrorStruct(pass, spec)
				if !found {
					continue
				}

				errorStatementNames = foundName
				errorStatementFound = true
				continue
			}

		}

	}

	return true
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

func returnsInvalidError(pass *analysis.Pass, names []string, returnStmt *ast.ReturnStmt) (*ast.Ident, bool) {
	for _, result := range returnStmt.Results {
		if identifier, ok := result.(*ast.Ident); ok {
			for _, name := range names {
				if identifier.Name == name {
					return identifier, true
				}
			}
		}
	}
	return nil, false
}

func isErrorStruct(pass *analysis.Pass, spec ast.Spec) ([]string, bool) {
	v, ok := spec.(*ast.ValueSpec)
	if !ok {
		return []string{}, false
	}

	if v.Type == nil {
		return []string{}, false
	}

	// if the initial values are non-nil
	if v.Values != nil {
		return []string{}, false
	}

	specType := pass.TypesInfo.TypeOf(v.Type)

	// if it is a pointer, get the elment type
	if pointer, ok := specType.(*types.Pointer); ok {
		specType = pointer.Elem()
	}

	// check if it's implements the error type
	if n, ok := specType.(*types.Named); ok {
		if !types.Implements(n, errTyp) {
			return []string{}, false
		}
		// not sure why there is more than one name here ??
		return getNamesFromNames(v.Names), true
	}
	return []string{}, false
}

func getNamesFromNames(nameList []*ast.Ident) []string {
	stringNames := make([]string, len(nameList))

	for i, ident := range nameList {
		stringNames[i] = ident.Name
	}

	return stringNames
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
