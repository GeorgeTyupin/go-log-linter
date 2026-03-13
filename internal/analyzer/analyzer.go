package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/GeorgeTyupin/go-log-linter/internal/analyzer/rules"
	"github.com/GeorgeTyupin/go-log-linter/internal/config"
)

func NewAnalyzer(cfg *config.Config) *analysis.Analyzer {
	Analyzer := &analysis.Analyzer{
		Name:     "loglinter",
		Doc:      "Checks log messages for style violations",
		Run:      newRun(cfg),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	return Analyzer
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

func newRun(cfg *config.Config) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		// Разбиваем кастомные паттерны из флага или конфига
		var extraPatterns []string
		if cfg.ExtraSensitivePatterns != "" {
			for _, p := range strings.Split(cfg.ExtraSensitivePatterns, ",") {
				p = strings.TrimSpace(p)
				if p != "" {
					extraPatterns = append(extraPatterns, p)
				}
			}
		}

		nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

		insp.Preorder(nodeFilter, func(n ast.Node) {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return
			}

			msg, lit, ok := extractLogMessage(pass, callExpr)
			if !ok {
				return
			}

			// 1. Проверка на секреты (Приоритет: прерываем остальные проверки)
			if !cfg.DisableNoSensitiveData {
				if diagnostic, violated := rules.CheckNoSensitiveData(msg, extraPatterns); violated {
					pass.Reportf(lit.Pos(), "%s", diagnostic)
					return
				}
			}

			// 2. Проверка на строчные буквы
			if !cfg.DisableLowercase {
				if diagnostic, violated := rules.CheckLowercase(msg); violated {
					pass.Report(analysis.Diagnostic{
						Pos:     lit.Pos(),
						Message: diagnostic,
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: "convert first letter to lowercase",
								TextEdits: []analysis.TextEdit{
									{
										// Заменяем первый символ внутри кавычек на строчный.
										// lit.Pos() указывает на открывающую кавычку, поэтому +1.
										Pos:     lit.Pos() + 1,
										End:     lit.Pos() + 1 + token.Pos(utf8.RuneLen([]rune(msg)[0])),
										NewText: []byte(string(unicode.ToLower([]rune(msg)[0]))),
									},
								},
							},
						},
					})
				}
			}

			// 3. Проверка на английский язык
			if !cfg.DisableEnglishOnly {
				if diagnostic, violated := rules.CheckEnglishOnly(msg); violated {
					pass.Reportf(lit.Pos(), "%s", diagnostic)
				}
			}

			// 4. Проверка на спецсимволы
			if !cfg.DisableNoSpecialChars {
				if diagnostic, violated := rules.CheckNoSpecialChars(msg); violated {
					pass.Report(analysis.Diagnostic{
						Pos:     lit.Pos(),
						Message: diagnostic,
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: "remove special characters and emojis",
								TextEdits: []analysis.TextEdit{
									{
										Pos:     lit.Pos() + 1,
										End:     lit.End() - 1,
										NewText: []byte(sanitizeMessage(msg)),
									},
								},
							},
						},
					})
				}
			}

		})

		return nil, nil
	}
}

// sanitizeMessage удаляет из сообщения все недопустимые символы (спецсимволы и эмодзи).
func sanitizeMessage(msg string) string {
	var b strings.Builder
	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || isAllowedPunct(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

var allowedPunct = map[rune]bool{
	' ': true, '-': true, '_': true, '/': true,
	',': true, ':': true,
	'(': true, ')': true, '[': true, ']': true,
	'\'': true,
}

func isAllowedPunct(r rune) bool {
	return allowedPunct[r]
}

func extractLogMessage(pass *analysis.Pass, call *ast.CallExpr) (string, *ast.BasicLit, bool) {
	if len(call.Args) == 0 {
		return "", nil, false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", nil, false
	}

	if !logMethods[sel.Sel.Name] {
		return "", nil, false
	}

	if !isLoggerCall(pass, sel) {
		return "", nil, false
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", nil, false
	}

	msg := strings.Trim(lit.Value, `"`+"`")
	return msg, lit, true
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
