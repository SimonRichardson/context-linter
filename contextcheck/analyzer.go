package contextcheck

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "contextcheck",
	Doc:  "reports context aware methods that are not used with context",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if isGenerated(file) {
			continue
		}

		// Check if the file imports context, then build up a map of context
		// imports, which also includes aliased imports.
		contextPaths := make(map[string]struct{})
		for _, imp := range file.Imports {
			if imp.Path.Value != `"context"` {
				continue
			}

			// Aliased import.
			if imp.Name != nil {
				contextPaths[imp.Name.Name] = struct{}{}
				continue
			}

			contextPaths[imp.Path.Value] = struct{}{}
		}

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			ident, ok := selExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			// If the first argument is a context aware call, then we can
			// assume that the context is being passed in, so we don't need
			// to check it.
			if isContextAware(callExpr, contextPaths) {
				return true
			}

			selectorName := selExpr.Sel
			fmt.Printf("%T %+v\n", ident, ident)
			fmt.Printf("%T %+v\n", selectorName, selectorName)

			return true
		})
	}

	return nil, nil
}

func isGenerated(file *ast.File) bool {
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "// Code generated ") && strings.HasSuffix(c.Text, " DO NOT EDIT.") {
				return true
			}
		}
	}

	return false
}

func isContextAware(callExpr *ast.CallExpr, imports map[string]struct{}) bool {
	if num := len(callExpr.Args); num == 0 {
		return false
	}

	first := callExpr.Args[0]
	if first == nil {
		return false
	}
	if callExpr, ok := first.(*ast.CallExpr); ok {
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selExpr.X.(*ast.Ident); ok {
				if _, ok := imports[ident.Name]; ok {
					return true
				}
			}
		}
	}
	return false
}
