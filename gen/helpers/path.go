package helpers

import (
	"os"
	"strings"
)

func Path(path ...string) string {
	return strings.Join(path, string(os.PathSeparator))
}

func MkDir(segments ...string) error {
	path := Path(segments...)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	return nil
}

func PackageName(path string) string {
	s := strings.Split(path, string(os.PathSeparator))
	return s[len(s)-1]
}
