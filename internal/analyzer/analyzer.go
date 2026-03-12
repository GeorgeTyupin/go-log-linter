package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/GeorgeTyupin/go-log-linter/internal/analyzer/rules"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "Checks log messages for style violations",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var logMethods = map[string]bool{
	"Debug": true, "Info": true, "Warn": true, "Warning": true,
	"Error": true, "Fatal": true, "Panic": true, "Print": true,
	"Printf": true, "Println": true,
}

var logPackages = map[string]bool{
	"log":             true,
	"log/slog":        true,
	"go.uber.org/zap": true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		msg, pos, ok := extractLogMessage(pass, callExpr)
		if !ok {
			return
		}

		for _, check := range []func(string) (string, bool){
			rules.CheckLowercase,
			rules.CheckEnglishOnly,
			rules.CheckNoSpecialChars,
			rules.CheckNoSensitiveData,
		} {
			if diagnostic, violated := check(msg); violated {
				pass.Reportf(pos, "%s", diagnostic)
			}
		}
	})

	return nil, nil
}

func extractLogMessage(pass *analysis.Pass, call *ast.CallExpr) (string, token.Pos, bool) {
	if len(call.Args) == 0 {
		return "", 0, false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", 0, false
	}

	if !logMethods[sel.Sel.Name] {
		return "", 0, false
	}

	if !isLoggerCall(pass, sel) {
		return "", 0, false
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", 0, false
	}

	msg := strings.Trim(lit.Value, `"`+"`")
	return msg, lit.Pos(), true
}

func isLoggerCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	if ident, ok := sel.X.(*ast.Ident); ok {
		obj := pass.TypesInfo.ObjectOf(ident)
		if obj == nil {
			return false
		}
		if pkgName, ok := obj.(*types.PkgName); ok {
			return logPackages[pkgName.Imported().Path()]
		}
	}

	t := pass.TypesInfo.TypeOf(sel.X)
	if t == nil {
		return false
	}

	t = deref(t)

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	return logPackages[named.Obj().Pkg().Path()]
}

func deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}
