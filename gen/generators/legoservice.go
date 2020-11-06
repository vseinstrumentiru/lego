// +build gen

package generators

import (
	. "github.com/dave/jennifer/jen"

	"github.com/vseinstrumentiru/lego/v2/gen/helpers"
)

func NewLegoService(path string) error {
	if err := helpers.MkDir(path); err != nil {
		return err
	}
	if err := LegoServiceStore(path); err != nil {
		return err
	}
	if err := LegoService(path); err != nil {
		return err
	}
	if err := EmptyGoFile("models", path); err != nil {
		return err
	}
	if err := EmptyGoFile("events", path); err != nil {
		return err
	}
	if err := EmptyGoFile("handlers", path); err != nil {
		return err
	}
	if err := EmptyGoFile("middleware", path); err != nil {
		return err
	}

	return nil
}

func LegoService(path string) error {
	pkg := helpers.PackageName(path)
	f := NewFile(pkg)

	f.Comment("Service contract")
	f.Comment("lego:service:contract")
	f.Type().Id("Service").Interface(
		Comment("add service methods here..."),
	)

	f.Comment("Service constructor")
	f.Comment("lego:service:provider")
	f.Func().Id("NewService").Params(Id("store").Id("Store")).Id("Service").
		Block(
			Return(
				Id("service").Values(
					Id("store").Op(":").Id("store"),
				),
			),
		)

	f.Comment("Service realization")
	f.Comment("lego:service")
	f.Type().Id("service").Struct(
		Id("store").Id("Store"),
	)

	return f.Save(helpers.Path(path, "service.go"))
}

func LegoServiceStore(path string) error {
	pkg := helpers.PackageName(path)
	f := NewFile(pkg)
	f.Comment("Store contract")
	f.Comment("lego:store:contract")
	f.Type().Id("Store").Interface(
		Comment("add store methods here..."),
	)

	if err := helpers.MkDir(path, pkg+"store"); err != nil {
		return err
	}

	if err := GitKeep(helpers.Path(path, pkg+"store")); err != nil {
		return err
	}

	return f.Save(helpers.Path(path, "store.go"))
}
