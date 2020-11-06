// +build gen

package service

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"io"

	"github.com/dave/jennifer/jen"
	"sagikazarmark.dev/mga/pkg/gentypes"
	"sagikazarmark.dev/mga/pkg/genutils"
	"sagikazarmark.dev/mga/pkg/jenutils"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/vseinstrumentiru/lego/v2/gen/helpers/parser"
)

// nolint: gochecknoglobals
var (
	marker = markers.Must(markers.MakeDefinition("lego:service:contract", markers.DescribesType, struct{}{}))
)

type Generator struct {
}

func (g Generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(marker); err != nil {
		return err
	}

	into.AddHelp(
		marker,
		markers.SimpleHelp("LeGo", "enables service generator for interface"),
	)

	return nil
}

func (g Generator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		outContents := g.generatePackage(ctx, root)
		if outContents == nil {
			continue
		}

		writeOut(ctx, root, outContents)
	}

	return nil
}

func writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	outputFile, err := ctx.Open(root, "service.go")
	if err != nil {
		root.AddError(err)

		return
	}
	defer outputFile.Close()
	n, err := outputFile.Write(outBytes)
	if err != nil {
		root.AddError(err)

		return
	}
	if n < len(outBytes) {
		root.AddError(io.ErrShortWrite)
	}
}

func (g Generator) generatePackage(ctx *genall.GenerationContext, root *loader.Package) []byte {
	ctx.Checker.Check(root, func(node ast.Node) bool {
		// ignore non-interfaces
		_, isStruct := node.(*ast.StructType)

		return isStruct
	})

	root.NeedTypesInfo()

	var serviceContract *parser.Object

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		if marker := info.Markers.Get(marker.Name); marker == nil {
			return
		}

		if info.Name != "Service" {
			return
		}

		typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
		if typeInfo == types.Typ[types.Invalid] {
			root.AddError(loader.ErrFromNode(fmt.Errorf("unknown type %s", info.Name), info.RawSpec))

			return
		}

		i, err := parser.GetObject(root.TypesInfo.TypeOf(info.RawSpec.Name))

		if err != nil {
			root.AddError(err)

			return
		}

		serviceContract = &i
	})

	if err != nil {
		root.AddError(err)

		return nil
	}

	if serviceContract == nil {
		return nil
	}

	packageName, packagePath := root.Name, root.PkgPath
	if pkgRefer, ok := ctx.OutputRule.(genutils.PackageRefer); ok {
		packageName, packagePath = pkgRefer.PackageRef(root)
	}

	file := parser.File{
		File: gentypes.File{
			Package: gentypes.PackageRef{
				Name: packageName,
				Path: packagePath,
			},
		},
		Object: *serviceContract,
	}

	outContents, err := generateCode(file)
	if err != nil {
		root.AddError(err)

		return nil
	}

	return outContents
}

func generateCode(file parser.File) ([]byte, error) {
	code := jen.NewFilePathName(file.Package.Path, file.Package.Name)

	const (
		serviceName = "service"
	)

	code.Type().Id(serviceName).Struct()

	for i := 0; i < len(file.Object.Methods); i++ {
		method := file.Object.Methods[i]

		var params []jen.Code
		for j := 0; j < len(method.Params); j++ {
			param := method.Params[j]

			params = append(params, jenutils.Type(jen.Id(param.Name), param.Type))
		}

		var results []jen.Code
		for j := 0; j < len(method.Results); j++ {
			param := method.Results[j]

			results = append(results, jenutils.Type(jen.Id(param.Name), param.Type))
		}

		code.Func().Params(
			jen.Id("s").Id(serviceName),
		).
			Id(method.Name).Params(params...).Params(results...).
			Block(
				jen.Panic(jen.Lit("not implemented")),
			).Line()
	}

	var buf bytes.Buffer
	err := code.Render(&buf)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
