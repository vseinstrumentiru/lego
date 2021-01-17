package generators

import (
	. "github.com/dave/jennifer/jen"

	"github.com/vseinstrumentiru/lego/v2/gen/helpers"
)

func EmptyGoFile(name string, path string) error {
	f := NewFile(helpers.PackageName(path))
	return f.Save(helpers.Path(path, name+".go"))
}
