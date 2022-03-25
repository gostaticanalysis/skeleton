package internal

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

// copy from golang.org/x/tools/analysis/passes/buildssa
func BuildSSA(pkg *packages.Package, mode ssa.BuilderMode) (*ssa.Program, []*ssa.Function, error) {

	prog := ssa.NewProgram(pkg.Fset, mode)

	// Create SSA packages for all imports.
	// Order is not significant.
	created := make(map[*types.Package]bool)
	var createAll func(pkgs []*types.Package)
	createAll = func(pkgs []*types.Package) {
		for _, p := range pkgs {
			if !created[p] {
				created[p] = true
				prog.CreatePackage(p, nil, nil, true)
				createAll(p.Imports())
			}
		}
	}
	createAll(pkg.Types.Imports())

	// Create and build the primary package.
	ssapkg := prog.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, false)
	ssapkg.Build()

	// Compute list of source functions, including literals,
	// in source order.
	var funcs []*ssa.Function
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			if fdecl, ok := decl.(*ast.FuncDecl); ok {

				// SSA will not build a Function
				// for a FuncDecl named blank.
				// That's arguably too strict but
				// relaxing it would break uniqueness of
				// names of package members.
				if fdecl.Name.Name == "_" {
					continue
				}

				// (init functions have distinct Func
				// objects named "init" and distinct
				// ssa.Functions named "init#1", ...)

				fn := pkg.TypesInfo.Defs[fdecl.Name].(*types.Func)
				if fn == nil {
					panic(fn)
				}

				f := ssapkg.Prog.FuncValue(fn)
				if f == nil {
					panic(fn)
				}

				var addAnons func(f *ssa.Function)
				addAnons = func(f *ssa.Function) {
					funcs = append(funcs, f)
					for _, anon := range f.AnonFuncs {
						addAnons(anon)
					}
				}
				addAnons(f)
			}
		}
	}

	return prog, funcs, nil
}
