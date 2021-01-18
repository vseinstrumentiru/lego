package generators

import (
	. "github.com/dave/jennifer/jen"

	"github.com/vseinstrumentiru/lego/v2/gen/internal/helpers"
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
	if err := LegoEvents(path); err != nil {
		return err
	}
	if err := EmptyGoFile("models", path); err != nil {
		return err
	}

	return nil
}

func LegoService(path string) error {
	pkg := helpers.PackageName(path)
	f := NewFile(pkg)

	f.Comment("+kit:endpoint")
	f.Comment("Service contract")
	f.Type().Id("Service").Interface(
		Comment("add service methods here..."),
	)

	f.Comment("Service constructor")
	f.Func().Id("NewService").Params(Id("store").Id("Store")).Id("Service").
		Block(
			Return(
				Id("service").Values(
					Id("store").Op(":").Id("store"),
				),
			),
		)

	f.Comment("Service realization")
	f.Type().Id("service").Struct(
		Id("store").Id("Store"),
	)

	return f.Save(helpers.Path(path, "service.go"))
}

func LegoServiceStore(path string) error {
	pkg := helpers.PackageName(path)
	f := NewFile(pkg)
	f.Comment("Store contract")
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

func LegoEvents(path string) error {
	pkg := helpers.PackageName(path)
	f := NewFile(pkg)
	f.Comment("+mga:event:dispatcher")
	f.Comment("Events contract")
	f.Type().Id("Events").Interface(
		Comment("add dispatcher methods here..."),
	)

	return f.Save(helpers.Path(path, "events.go"))
}
