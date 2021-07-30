package store

import (
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"

	"github.com/vseinstrumentiru/lego/v2/gen/internal/pkg"
)

const servicePath = "./internal/app/"

func create(c *config) ([]string, error) {
	var err error

	parts := strings.Split(c.Name, ":")
	if len(parts) != 2 {
		return nil, errors.New("name must be contains two parts ([service name]:[store name])")
	}

	fileName := parts[1]
	c.Name = parts[0]

	if c.Name, c.workDir, err = pkg.GetPath(servicePath + c.Name); err != nil {
		return nil, err
	}
	c.workDir += "/" + c.Name

	cmd := pkg.NewExecutor(c.workDir, c.Verbose, c.DryRun)

	if _, err = cmd.Mkdir(fmt.Sprintf("/../%s", c.Name)); err != nil {
		return nil, err
	}

	if ok, err := cmd.Exist("./store.go"); err != nil {
		return nil, err
	} else if !ok {
		f := jen.NewFile(c.Name)
		f.Comment("Store contract")
		f.Type().Id("Store").Interface(
			jen.Comment("add store methods here..."),
		)

		err = cmd.Exec(func() error {
			return f.Save(fmt.Sprintf("%s/%s.go", c.workDir, "store"))
		}, "store interface created")

		if err != nil {
			return nil, err
		}
	}

	if fileName == "store" {
		return nil, nil
	}

	storePkg := fmt.Sprintf("%sstore", c.Name)
	storePath := fmt.Sprintf("/%s/%s.go", storePkg, fileName)

	if _, err = cmd.Mkdir("./" + storePkg); err != nil {
		return nil, err
	}

	rootPkg, err := pkg.GetRootPkg(c.workDir + "/../../..")
	if err != nil {
		return nil, err
	}

	fullPkg := fmt.Sprintf("%s/internal/app/%s", rootPkg, c.Name)

	f := jen.NewFile(storePkg)

	f.ImportName(fullPkg, c.Name)
	f.Comment(fmt.Sprintf("New%s store constructor", strcase.ToCamel(fileName)))
	f.Func().Id(fmt.Sprintf("New%s", strcase.ToCamel(fileName))).Params().Qual(fullPkg, "Store").
		Block(
			jen.Return(
				jen.Id(fileName).Values(),
			),
		)

	f.Comment(fileName + " is Store implementation")
	f.Type().Id(fileName).Struct()

	err = cmd.Exec(func() error {
		return f.Save(c.workDir + storePath)
	}, fmt.Sprintf("store %s created", fileName))

	if err != nil {
		return nil, err
	}

	return cmd.Result(), nil
}
