package interfaces

import (
	"bytes"
	"fmt"
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

type implementationMark struct {
	Pkg    string `marker:"pkg,optional"`
	Parent string `marker:"parent,optional"`
}

func NewGenerator(name string, mark string, filename string, implementation string) genall.Generator {
	gen := generator{
		name:     name,
		mark:     markers.Must(markers.MakeDefinition(mark, markers.DescribesType, struct{}{})),
		filename: filename,
	}

	if implementation != "" {
		gen.implementation = markers.Must(markers.MakeDefinition(implementation, markers.DescribesType, implementationMark{}))
	}

	return gen
}

type generator struct {
	mark           *markers.Definition
	implementation *markers.Definition
	name           string
	filename       string
}

func (g generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(g.mark); err != nil {
		return err
	}

	into.AddHelp(
		g.mark,
		markers.SimpleHelp("LeGo", "enables service generator for interface"),
	)
	if g.implementation != nil {
		if err := into.Register(g.implementation); err != nil {
			return err
		}

		into.AddHelp(
			g.implementation,
			markers.SimpleHelp("LeGo", "enables service generator for interface"),
		)
	}

	return nil
}

func (g generator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		outContents := g.generatePackage(ctx, root)
		if outContents == nil {
			continue
		}

		g.writeOut(ctx, root, outContents)
	}

	return nil
}

func (g generator) writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	outputFile, err := ctx.Open(root, g.filename)
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

func (g generator) generatePackage(ctx *genall.GenerationContext, root *loader.Package) []byte {
	root.NeedTypesInfo()

	var serviceContract *parser.Object
	var implementations []parser.Object

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		var isImplementation bool
		if marker := info.Markers.Get(g.mark.Name); marker == nil {
			if marker := info.Markers.Get(g.implementation.Name); marker == nil {
				return
			}

			isImplementation = true
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

		if isImplementation {
			implementations = append(implementations, i)
		} else {
			serviceContract = &i
		}
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

	outContents, err := g.generateCode(file)
	if err != nil {
		root.AddError(err)

		return nil
	}

	return outContents
}

func (g generator) generateCode(file parser.File) ([]byte, error) {
	code := jen.NewFilePathName(file.Package.Path, file.Package.Name)

	code.Comment(fmt.Sprintf("+lego:service:parent=%s:pkg=%s", file.Object.Name, file.Object.Package.Path))
	code.Type().Id(g.name).Struct()

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
			jen.Id("s").Id(g.name),
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
