package main

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var panicCheckAnalyzer = &analysis.Analyzer{
	Name: "dontpanic",
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

	checkUsage(pass, pkg, sel.Sel, "log", "Fatal")
	checkUsage(pass, pkg, sel.Sel, "os", "Exit")
	return true
}

func checkUsage(pass *analysis.Pass, pkgID, funID *ast.Ident, pkg, funName string) {
	importedPkg, ok := pass.TypesInfo.Uses[pkgID].(*types.PkgName)
	if !ok {
		return
	}
	if importedPkg.Imported().Path() == pkg && funID.Name == funName {
		pass.Reportf(funID.Pos(), "should not use %s.%s() outside main() func", pkg, funName)
	}
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
