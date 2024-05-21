package analyzer

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

//nolint:gochecknoglobals
var (
	skipForTest bool
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "exportes",
		Doc:  "Doc: used to find the private fields and methods of structures",
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
			if strings.Contains(ff.Name.String(), "New") {
				return true
			}
			if skipForTest {
				return true
			}
			ast.Inspect(ff.Body, checkExport(pass, ff.Recv))

			return true
		})
	}
	return nil, nil
}

func checkExport(pass *analysis.Pass, receiver *ast.FieldList) func(ast.Node) bool {
	return func(node ast.Node) bool {
		if checkExportStructFields(pass, node) {
			return true
		}

		selectExpr, isSelector := node.(*ast.SelectorExpr)
		if !isSelector {
			return true
		}

		selection, ok := pass.TypesInfo.Selections[selectExpr]
		if !ok {
			return true
		}

		if selection.Kind() != types.MethodVal && selection.Kind() != types.FieldVal {
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
	_, ok = pass.TypesInfo.TypeOf(lit.Type).Underlying().(*types.Struct)
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
