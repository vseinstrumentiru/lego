package pkg

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func GetPath(newDir string) (name, path string, err error) {
	if path, err = os.Getwd(); err != nil {
		return "", "", err
	}

	path += "/" + newDir

	if path, err = filepath.Abs(path); err != nil {
		return "", "", err
	}

	lastSlash := strings.LastIndex(path, "/")
	name = path[lastSlash+1:]
	path = path[0:lastSlash]

	return name, path, nil
}

func GetRootPkg(rootPath string) (name string, err error) {
	file, err := os.Open(rootPath + "/go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	return strings.TrimSpace(line[7:]), nil
}
