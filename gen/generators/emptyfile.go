package generators

import (
	"os"

	"github.com/vseinstrumentiru/lego/v2/gen/helpers"
)

func EmptyFile(name string, path string) error {
	file, err := os.OpenFile(helpers.Path(path, name), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func GitKeep(path string) error {
	return EmptyFile(".gitkeep", path)
}
