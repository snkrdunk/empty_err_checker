package empty_err_checker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "empty_err_checker is ..."

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

	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		switch current := n.(type) {
		case *ast.IfStmt:
			if !push {
				return false
			}
			parentIfStmtSlice := []*ast.IfStmt{}
			for i := range stack[:len(stack)-1] {
				if s, ok := stack[i].(*ast.IfStmt); ok && i != len(stack)-1 {
					parentIfStmtSlice = append(parentIfStmtSlice, s)
				}
			}
			res := getEmptyErrReturnStmt(pass, current, parentIfStmtSlice)
			if res != nil {
				pass.Reportf(res.Pos(), "returned error is not checked.")
			}
		}
		return true
	})

	return nil, nil
}

func getEmptyErrReturnStmt(pass *analysis.Pass, current *ast.IfStmt, parents []*ast.IfStmt) *ast.ReturnStmt {
	for _, b := range current.Body.List {
		r, ok := b.(*ast.ReturnStmt)
		if !ok {
			continue
		}
		if !isReturnErr(pass, r) {
			return nil
		}
		if isCheckedIfErrIsNil(current) {
			return nil
		}
		if isAssignedErr(current) {
			return nil
		}
		for i := range parents {
			if isCheckedIfErrIsNil(parents[i]) {
				return nil
			}
			if isAssignedErr(parents[i]) {
				return nil
			}
		}
		return r
	}
	return nil
}

func isReturnErr(pass *analysis.Pass, r *ast.ReturnStmt) bool {
	var isReturnErr bool
	for _, result := range r.Results {
		ident, ok := result.(*ast.Ident)
		if !ok {
			continue
		}
		// NOTE: returnにerrがあるかどうかチェック
		if ident.Name == "err" && pass.TypesInfo.Types[ident].Type.String() == "error" {
			isReturnErr = true
		}
	}
	return isReturnErr
}

func isCheckedIfErrIsNil(ifStmt *ast.IfStmt) bool {
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	xIdent, ok := binaryExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	yIdent, ok := binaryExpr.Y.(*ast.Ident)
	if !ok {
		return false
	}
	return xIdent.Name == "err" && yIdent.Name == "nil" && binaryExpr.Op.String() == "!="
}

func isAssignedErr(ifStmt *ast.IfStmt) bool {
	var isAssignedErr bool
	for _, b := range ifStmt.Body.List {
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
	}
	return isAssignedErr
}
