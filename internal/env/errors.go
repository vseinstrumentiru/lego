package env

import "github.com/spf13/viper"

type ErrConfigFileNotFound = viper.ConfigFileNotFoundError

func IsFileNotFound(err error) bool {
	_, ok := err.(viper.ConfigFileNotFoundError)

	return ok
}
