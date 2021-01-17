package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

type Type1 struct {
	Val int
}

func (t Type1) Instance() Type1 {
	return t
}

func NewType1Pointer() *Type1 {
	return &Type1{Val: 1}
}

func NewType1() Type1 {
	return Type1{Val: 1}
}

func Test_Singleton(t *testing.T) {
	c := New()
	assert.NoError(t, c.Register(NewType1Pointer))
	var type1 *Type1
	assert.NoError(t, c.Execute(func(a *Type1) { type1 = a }))
	assert.NotNil(t, type1)
	assert.Equal(t, 1, type1.Val)
	type1.Val = 2
	var type2 *Type1
	assert.NoError(t, c.Execute(func(a *Type1) { type2 = a }))
	assert.Equal(t, 2, type2.Val)
}

func Test_Not_Singleton(t *testing.T) {
	c := New()
	assert.NoError(t, c.Register(NewType1))
	var type1 Type1
	assert.NoError(t, c.Execute(func(a Type1) { type1 = a }))
	assert.Equal(t, 1, type1.Val)
	type1.Val = 2
	assert.Equal(t, 2, type1.Val)
	var type2 Type1
	assert.NoError(t, c.Execute(func(a Type1) { type2 = a }))
	assert.Equal(t, 1, type2.Val)
}

func Test_Instance(t *testing.T) {
	c := New()
	{
		type1 := Type1{Val: 5}
		assert.NoError(t, c.Instance(type1))
	}
	var type1 Type1
	assert.NoError(t, c.Execute(func(a Type1) { type1 = a }))
	assert.Equal(t, 5, type1.Val)
	type1.Val = 2
	assert.Equal(t, 2, type1.Val)
	var type2 Type1
	assert.NoError(t, c.Execute(func(a Type1) { type2 = a }))
	assert.Equal(t, 5, type2.Val)
}

func Test_Instance_Pointer(t *testing.T) {
	c := New()
	{
		type1 := &Type1{Val: 5}
		assert.NoError(t, c.Instance(type1))
	}
	var type1 *Type1
	assert.NoError(t, c.Execute(func(a *Type1) { type1 = a }))
	assert.Equal(t, 5, type1.Val)
	type1.Val = 2
	assert.Equal(t, 2, type1.Val)
	var type2 *Type1
	assert.NoError(t, c.Execute(func(a *Type1) { type2 = a }))
	assert.Equal(t, 2, type2.Val)
}

type Type2 struct {
	Val int
}

type Type3 struct {
	Val int
}

type App struct {
	T1 Type1
	T2 Type2
	T3 Type3
}

func (app *App) ProvideType3(t1 Type1, t2 Type2) Type3 {
	return Type3{Val: t1.Val + t2.Val}
}

func (app *App) ProvideType2() Type2 {
	return Type2{Val: 2}
}

func TestContainer_MakeApp(t *testing.T) {
	c := New()
	assert.NoError(t, c.Register(NewType1))
	app := App{}

	assert.NoError(t, c.Make(&app))

	ass := assert.New(t)
	ass.Equal(1, app.T1.Val)
	ass.Equal(2, app.T2.Val)
	ass.Equal(3, app.T3.Val)
}

func TestContainer_SameInstances(t *testing.T) {
	c := New()
	assert.NoError(t, c.Instance(Type1{Val: 1}))
	assert.Error(t, c.Instance(Type1{Val: 1}))
	assert.NoError(t, c.Instance(Type1{Val: 2}, dig.Name("second")))

	err := c.Execute(func(in struct {
		dig.In
		First  Type1
		Second Type1 `name:"second"`
	}) {
		assert.Equal(t, 1, in.First.Val)
		assert.Equal(t, 2, in.Second.Val)
	})

	assert.NoError(t, err)
}
