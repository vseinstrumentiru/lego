package config

type Runtime interface {
	On(key string, callback interface{}) (ok bool)
	Is(key string) bool
	Not(key string) bool
	Run(apps ...interface{})
}

type FlagBinder interface {
	To(key string)
}

type Env interface {
	SetDefault(key string, value interface{})
	SetAlias(alias string, originalKey string)
	SetFlag(name string, value interface{}, usage string) FlagBinder
}

type Validatable interface {
	Validate() (err error)
}

type WithDefaults interface {
	SetDefaults(env Env)
}
