package config

func Undefined() *Application {
	return &Application{
		Name:       "undefined",
		DataCenter: "undefined",
	}
}

type Application struct {
	Name       string
	DataCenter string
}

func (c Application) FullName() string {
	name := c.Name
	if c.DataCenter != "" {
		name += "-" + c.DataCenter
	}

	return name
}
