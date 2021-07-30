package env

import (
	"errors"

	"github.com/spf13/viper"
)

type ErrConfigFileNotFound = viper.ConfigFileNotFoundError

func IsFileNotFound(err error) bool {
	return errors.As(err, &viper.ConfigFileNotFoundError{})
}
