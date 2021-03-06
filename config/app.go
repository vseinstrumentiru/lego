package config

const Undefined = "undefined"

func IsUndefined(s string) bool {
	return s == "" || s == Undefined
}

func UndefinedApplication() *Application {
	return &Application{
		Name:       Undefined,
		DataCenter: Undefined,
	}
}

type Application struct {
	Name       string
	DataCenter string
	DebugMode  bool
	LocalMode  bool
}

func (c Application) FullName() string {
	name := c.Name
	if c.DataCenter != "" {
		name += "-" + c.DataCenter
	}

	return name
}
