package mysql

import "github.com/vseinstrumentiru/lego/v2/transport/mysql"

// Deprecated: use mysql.Config.
type Config struct {
	Config mysql.Config `mapstructure:",squash"`
}

// Deprecated: already sets.
func SetLogger(interface{}) {}
