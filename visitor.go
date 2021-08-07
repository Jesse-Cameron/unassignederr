package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var errTyp = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

type errorVisitor struct {
	pass              *analysis.Pass
	ErrStatementNames []string
	ErrStatementFound bool
	ErrStatement      *ast.Ident
}

func findUnassignedErrors(pass *analysis.Pass, fd *ast.FuncDecl) (*ast.Ident, bool) {
	visitor := errorVisitor{pass: pass}
	ast.Walk(&visitor, fd)
	return visitor.ErrStatement, visitor.ErrStatementFound
}

func (v *errorVisitor) Visit(node ast.Node) ast.Visitor {
	// if we have already found an error in our walking
	if v.ErrStatementFound {
		switch n := node.(type) {
		case *ast.ReturnStmt:
			// if we are returning the error
			errorStatement, ok := returnsInvalidError(v.ErrStatementNames, n)
			if ok {
				v.ErrStatement = errorStatement
				return nil // early return as we've found invalid return
			}

		case *ast.AssignStmt:
			// if we find an assignment statement
			// check if the error is assigned
			errorsLeftOver := isErrAssigned(n.Lhs, v.ErrStatementNames)
			// if we have reassigned all of the errors, then restart the search
			if len(errorsLeftOver) == 0 {
				v.ErrStatementFound = false
			}
			v.ErrStatementNames = errorsLeftOver

		default:
			// no nothing
		}
	} else {
		// check if the node declares an error struct
		errNames, found := foundErrorStruct(v.pass, node)
		if found {
			v.ErrStatementFound = found
			v.ErrStatementNames = errNames
		}
	}

	// keep searching
	return v
}

func foundErrorStruct(pass *analysis.Pass, node ast.Node) ([]string, bool) {
	// if there is a var statement
	genDecl, found := isVarStatment(node)
	if !found {
		return []string{}, false
	}

	// if the var statement declares an error struct
	for _, spec := range genDecl.Specs {
		foundName, found := isErrorStruct(pass, spec)
		if found {
			return foundName, true
		}
	}

	return []string{}, false
}

func isVarStatment(stmt ast.Node) (*ast.GenDecl, bool) {
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

func returnsInvalidError(names []string, returnStmt *ast.ReturnStmt) (*ast.Ident, bool) {
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

// isErrAssigned checks the expression to see if assigned to any of the errname
// removes them if they do
func isErrAssigned(exprs []ast.Expr, errNames []string) []string {
	assignedNames := []string{}
	for _, assignedExpr := range exprs {
		if ident, ok := assignedExpr.(*ast.Ident); ok {
			assignedNames = append(assignedNames, ident.Name)
		}
	}

	// get all of the assigned names that are in the errNames
	subset := sliceSubset(errNames, assignedNames)

	return subset
}

// remove any items from a that are in b
func sliceSubset(a, b []string) []string {
	results := []string{}

	for _, aValue := range a {
		if !existsInList(b, aValue) {
			results = append(results, aValue)
		}
	}

	return results
}

func existsInList(list []string, search string) bool {
	for _, item := range list {
		if item == search {
			return true
		}
	}
	return false
}
