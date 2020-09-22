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
