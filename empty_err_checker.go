package empty_err_checker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "empty_err_checker is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "empty_err_checker",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.IfStmt:
			var isAssignedErr bool
			for _, b := range n.Body.List {
				// NOTE: check if err is assigned
				if assignStmt, ok := b.(*ast.AssignStmt); ok {
					for _, expr := range assignStmt.Lhs {
						ident, ok := expr.(*ast.Ident)
						if !ok {
							continue
						}
						if ident.Name == "err" {
							isAssignedErr = true
						}
					}
				}

				if returnStmt, ok := b.(*ast.ReturnStmt); ok {
					for _, result := range returnStmt.Results {
						ident, ok := result.(*ast.Ident)
						if !ok {
							continue
						}
						if ident.Name == "err" && pass.TypesInfo.Types[ident].Type.String() == "error" {
							binaryExpr, ok := n.Cond.(*ast.BinaryExpr)
							if !ok {
								if !isAssignedErr {
									pass.Reportf(returnStmt.Pos(), "returned error is not checked.")
								}
								continue
							}
							xIdent, ok := binaryExpr.X.(*ast.Ident)
							if !ok {
								if !isAssignedErr {
									pass.Reportf(returnStmt.Pos(), "returned error is not checked.")
								}
								continue
							}
							yIdnet, ok := binaryExpr.Y.(*ast.Ident)
							if !ok {
								if !isAssignedErr {
									pass.Reportf(returnStmt.Pos(), "returned error is not checked.")
								}
								continue
							}
							if !(xIdent.Name == "err" && yIdnet.Name == "nil" && binaryExpr.Op.String() == "!=") && !isAssignedErr {
								pass.Reportf(returnStmt.Pos(), "returned error is not checked.")
							}
						}
					}
				}
			}
		}
	})

	return nil, nil
}
