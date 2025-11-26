package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var panicCheckAnalyzer = &analysis.Analyzer{
	Name: "paniclint",
	Doc:  "check for panic(), log.Fatal() and os.Exit() outside of main package",
	Run:  run,
}

func checkExitOutsideMain(node ast.Node, pass *analysis.Pass, isMainPkg bool) bool {
	decl, ok := node.(*ast.FuncDecl)
	if ok {
		return !isMainPkg || decl.Name.Name != "main"
	}

	call, ok := node.(*ast.CallExpr)
	if !ok {
		return true
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return true
	}

	if pkg.Name == "log" && sel.Sel.Name == "Fatal" {
		pass.Reportf(call.Pos(), "should not use log.Fatal() outside main() func")
	}

	if pkg.Name == "os" && sel.Sel.Name == "Exit" {
		pass.Reportf(call.Pos(), "should not use os.Exit() outside main() func")
	}

	return true
}

func checkPanic(node ast.Node, pass *analysis.Pass) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return true
	}

	id, ok := call.Fun.(*ast.Ident)
	if !ok {
		return true
	}

	if id.Name == "panic" {
		pass.Reportf(call.Pos(), "should not use panic()")
	}

	return true
}

func run(pass *analysis.Pass) (interface{}, error) {
	isMainPkg := pass.Pkg.Name() == "main"
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			return checkExitOutsideMain(node, pass, isMainPkg)
		})
		ast.Inspect(file, func(node ast.Node) bool {
			return checkPanic(node, pass)
		})
	}
	return nil, nil
}

func main() {
	singlechecker.Main(panicCheckAnalyzer)
}
