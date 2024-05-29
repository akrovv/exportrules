package analyzer

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "exportrules",
		Doc:  "Doc: uses for verification using private fields / methods",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			ff, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if strings.HasPrefix(ff.Name.Name, "New") {
				ast.Inspect(ff.Body, checkNewFunc(pass))
				return false
			}

			ast.Inspect(ff.Body, checkExport(pass, ff.Recv))
			return false
		})
	}
	return nil, nil
}

func checkExport(pass *analysis.Pass, receiver *ast.FieldList) func(ast.Node) bool {
	return func(node ast.Node) bool {
		if checkExportStructFields(pass, node) {
			return true
		}

		selectExpr, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		selection, ok := pass.TypesInfo.Selections[selectExpr]
		if !ok {
			return true
		}

		var message string
		if selection.Kind() == types.FieldVal {
			message = "access to private field"
		} else {
			message = "call to private method"
		}

		if compareStruct(pass, receiver, pass.TypesInfo.TypeOf(selectExpr.X)) {
			return true
		}

		if !selectExpr.Sel.IsExported() {
			pass.Report(analysis.Diagnostic{
				Pos:     selectExpr.Pos(),
				Message: message,
			})
		}
		return true
	}
}

func checkExportStructFields(pass *analysis.Pass, node ast.Node) bool {
	lit, ok := node.(*ast.CompositeLit)
	if !ok {
		return false
	}

	t := pass.TypesInfo.TypeOf(lit.Type)
	if t == nil {
		return false
	}
	_, ok = t.Underlying().(*types.Struct)
	if !ok {
		return false
	}

	for _, elem := range lit.Elts {
		kv, ok := elem.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		field, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		if !field.IsExported() {
			pass.Report(analysis.Diagnostic{
				Pos:     kv.Pos(),
				Message: "initialization of private field: " + field.Name,
			})
		}
	}

	return true
}

func checkNewFunc(pass *analysis.Pass) func(ast.Node) bool {
	return func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CompositeLit:
			for _, elt := range n.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				checkKeyValue(pass, kv)
			}
		case *ast.AssignStmt:
			if len(n.Lhs) != len(n.Rhs) {
				return true
			}
			for i, rhs := range n.Rhs {
				cl, ok := rhs.(*ast.CompositeLit)
				if !ok {
					continue
				}

				ident, ok := cl.Type.(*ast.Ident)
				if !ok || ident.IsExported() {
					continue
				}

				if selExpr, ok := n.Lhs[i].(*ast.SelectorExpr); ok {
					if xIdent, ok := selExpr.X.(*ast.Ident); ok {
						pass.Reportf(selExpr.Pos(), "private field %q should not be initialized with a composite literal", xIdent.Name+"."+selExpr.Sel.Name)
					}
				}
			}
		}

		return true
	}
}

func checkKeyValue(pass *analysis.Pass, kv *ast.KeyValueExpr) {
	key, ok := kv.Key.(*ast.Ident)
	if !ok {
		return
	}

	value, ok := kv.Value.(*ast.CompositeLit)
	if !ok {
		return
	}

	valueIdent, ok := value.Type.(*ast.Ident)
	if !ok || valueIdent.IsExported() {
		return
	}

	pass.Reportf(kv.Pos(), "private field %q should not be initialized with a composite literal", key.Name)
}

func compareStruct(pass *analysis.Pass, receiver *ast.FieldList, typeInfo types.Type) bool {
	if receiver == nil {
		return false
	}
	if len(receiver.List) == 0 {
		return false
	}

	methodStruct := pass.TypesInfo.TypeOf(receiver.List[0].Type)
	if methodStruct == nil {
		return false
	}

	return types.Identical(methodStruct, typeInfo)
}
