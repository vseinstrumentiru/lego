// +build gen

package interfaces

import (
	"golang.org/x/tools/go/packages"
	"sagikazarmark.dev/mga/pkg/genutils"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type Config struct {
	Name     string
	Mark     string
	ImplMark string
	Paths    []string
	Output   string
	FileName string
}

func Generate(cfg Config) (*genall.Runtime, error) {
	generator := NewGenerator(cfg.Name, cfg.Mark, cfg.FileName, cfg.ImplMark)
	generators := genall.Generators{&generator}

	runtime, err := forRoots(generators, cfg.Paths...)
	if err != nil {
		return nil, err
	}

	outputRule, err := genutils.LookupOutput(cfg.Output)
	if err != nil {
		return nil, err
	}

	runtime.OutputRules.Default = outputRule

	return runtime, nil
}

// copied from genall package to override package loader configuration.
//
// required for supporting various types (basic type aliases, imports from other packages).
func forRoots(g genall.Generators, rootPaths ...string) (*genall.Runtime, error) {
	roots, err := loader.LoadRootsWithConfig(
		&packages.Config{
			Mode: packages.NeedDeps | packages.NeedTypes,
		},
		rootPaths...,
	)
	if err != nil {
		return nil, err
	}
	rt := &genall.Runtime{
		Generators: g,
		GenerationContext: genall.GenerationContext{
			Collector: &markers.Collector{
				Registry: &markers.Registry{},
			},
			Roots:     roots,
			InputRule: genall.InputFromFileSystem,
			Checker:   &loader.TypeChecker{},
		},
		OutputRules: genall.OutputRules{Default: genall.OutputToNothing},
	}
	if err := rt.Generators.RegisterMarkers(rt.Collector.Registry); err != nil {
		return nil, err
	}

	return rt, nil
}
