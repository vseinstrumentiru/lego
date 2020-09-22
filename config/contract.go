package config

type FlagBinder interface {
	To(key string)
}

type Env interface {
	SetDefault(key string, value interface{})
	SetAlias(alias string, originalKey string)
	SetFlag(name string, value interface{}, usage string) FlagBinder
}

type Validateable interface {
	Validate() error
}

type ConfigWithDefaults interface {
	SetDefaults(env Env)
}
