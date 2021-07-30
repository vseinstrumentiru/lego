package parser

import (
	"go/types"

	"emperror.dev/errors"
	"sagikazarmark.dev/mga/pkg/gentypes"
)

type Object struct {
	gentypes.TypeRef
	Methods []ObjectMethod
}

type ObjectParam struct {
	Name string
	Type types.Type
}

type ObjectMethod struct {
	Name    string
	Params  []ObjectParam
	Results []ObjectParam
}

//nolint:interfacer
// GetObject ParseEvent parses an object as an event.
func GetObject(in types.Type) (Object, error) {
	named, ok := in.(*types.Named)
	if !ok {
		return Object{}, errors.Errorf("%q is not a named type", in.String())
	}

	object := Object{
		TypeRef: gentypes.TypeRef{
			Name: named.Obj().Name(),
			Package: gentypes.PackageRef{
				Name: named.Obj().Pkg().Name(),
				Path: named.Obj().Pkg().Path(),
			},
		},
	}

	obj, _ := named.Underlying().(*types.Interface)

	for i := 0; i < obj.NumMethods(); i++ {
		m := obj.Method(i)

		method := ObjectMethod{
			Name: m.Name(),
		}

		sig, ok := m.Type().(*types.Signature)
		if !ok {
			continue
		}

		for j := 0; j < sig.Params().Len(); j++ {
			rawParam := sig.Params().At(j)

			method.Params = append(method.Params, ObjectParam{
				Name: rawParam.Name(),
				Type: rawParam.Type(),
			})
		}

		for j := 0; j < sig.Results().Len(); j++ {
			rawParam := sig.Results().At(j)
			method.Results = append(method.Results, ObjectParam{
				Name: rawParam.Name(),
				Type: rawParam.Type(),
			})
		}

		object.Methods = append(object.Methods, method)
	}

	return object, nil
}
