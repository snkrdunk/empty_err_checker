package empty_err_checker

import (
	"go/ast"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "empty_err_checker",
	Doc:  "empty_err_checker is checking whether the return value 'err' is nil.",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func init() {
	register.Plugin(Analyzer.Name, New)
}

type PluginImpl struct{}

func New(settings any) (register.LinterPlugin, error) {
	return &PluginImpl{}, nil
}

func (f *PluginImpl) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

func (f *PluginImpl) GetLoadMode() string {
	return register.LoadModeSyntax
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
				if s, ok := stack[i].(*ast.IfStmt); ok {
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
		if isCheckedIfErrIsNil(current.Cond, nil) {
			return nil
		}
		if isAssignedErr(current) {
			return nil
		}
		for i := range parents {
			if isCheckedIfErrIsNil(parents[i].Cond, nil) {
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

// NOTE: if err != nil && !isValid() || isChecked() {} のような条件を再起的にチェックしてる
func isCheckedIfErrIsNil(cond any, preBinaryExpr *ast.BinaryExpr) bool {
	switch n := cond.(type) {
	case *ast.BinaryExpr:
		xRes := isCheckedIfErrIsNil(n.X, cond.(*ast.BinaryExpr))
		yRes := isCheckedIfErrIsNil(n.Y, cond.(*ast.BinaryExpr))
		return xRes || yRes
	case *ast.Ident:
		if preBinaryExpr == nil {
			return false
		}
		xIdent, ok := preBinaryExpr.X.(*ast.Ident)
		if !ok {
			return false
		}
		yIdent, ok := preBinaryExpr.Y.(*ast.Ident)
		if !ok {
			return false
		}
		return xIdent.Name == "err" && yIdent.Name == "nil" && preBinaryExpr.Op.String() == "!="
	}
	return false
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
