// +build gen

package generators

import (
	. "github.com/dave/jennifer/jen"

	"github.com/vseinstrumentiru/lego/gen/helpers"
)

func EmptyGoFile(name string, path string) error {
	f := NewFile(helpers.PackageName(path))
	return f.Save(helpers.Path(path, name+".go"))
}
