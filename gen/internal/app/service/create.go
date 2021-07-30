package service

import (
	"fmt"

	jen "github.com/dave/jennifer/jen"

	"github.com/vseinstrumentiru/lego/v2/gen/internal/pkg"
)

const servicePath = "./internal/app/"

func create(c *config) ([]string, error) {
	var err error
	if c.Name, c.workDir, err = pkg.GetPath(servicePath + c.Name); err != nil {
		return nil, err
	}

	c.workDir += "/" + c.Name
	cmd := pkg.NewExecutor(c.workDir, c.Verbose, c.DryRun)

	if _, err = cmd.Mkdir(fmt.Sprintf("/../%s", c.Name)); err != nil {
		return nil, err
	}

	f := jen.NewFile(c.Name)

	f.Comment("+kit:endpoint:errorStrategy=service")
	f.Line()
	f.Comment("Service contract")
	f.Type().Id("Service").Interface(
		jen.Comment("add service methods here..."),
	)

	f.Comment("NewService Service constructor")
	f.Func().Id("NewService").Params(jen.Id("store").Id("Store")).Id("Service").
		Block(
			jen.Return(
				jen.Id("service").Values(
					jen.Id("store").Op(":").Id("store"),
				),
			),
		)

	f.Comment("service implementation of Service")
	f.Type().Id("service").Struct(
		jen.Id("store").Id("Store"),
	)

	err = cmd.Exec(func() error {
		return f.Save(c.workDir + "/service.go")
	}, fmt.Sprintf("service %s created", c.Name))

	if err != nil {
		return nil, err
	}

	return cmd.Result(), nil
}
