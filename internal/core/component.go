package core

import "reflect"

type LeGo struct {
	component
}

type component interface {
	Component()
}

type Typed interface {
	Type()
}

func typ(t interface{}) string {
	rt := reflect.TypeOf(t)

	if rt.Kind() != reflect.Func || rt.NumIn() <= 0 {
		return ""
	}

	return rt.In(0).Name()
}

func componentName(t interface{}) string {
	rt := reflect.TypeOf(t)

	if rt.Kind() != reflect.Func || rt.NumIn() <= 0 {
		return ""
	}

	c := rt.In(0)

	if !c.Implements(reflect.TypeOf((*component)(nil)).Elem()) {
		return ""
	}

	return c.PkgPath()
}
