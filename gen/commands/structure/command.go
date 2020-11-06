// +build gen

package structure

import (
	"os"
	"path/filepath"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vseinstrumentiru/lego/v2/gen/generators"
	"github.com/vseinstrumentiru/lego/v2/gen/helpers"
)

var Command = &cobra.Command{
	Use:  "init [path]",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetDefault("structure", map[string]interface{}{
			"cmd": map[string]interface{}{
				"server": "legostarter",
			},
			"internal": map[string]interface{}{
				"app": map[string]interface{}{
					"example": "legoservice",
				},
				"platform": map[string]interface{}{
					".gitkeep": "emptyfile",
				},
				"common": map[string]interface{}{
					".gitkeep": "emptyfile",
				},
			},
			"pkg": map[string]interface{}{
				".gitkeep": "emptyfile",
			},
		})

		if len(args) == 0 {
			args = append(args, "")
		}

		if path, err := cmd.Flags().GetString("config"); err == nil && path != "" {
			viper.AddConfigPath(path)
			emperror.Panic(viper.ReadInConfig())
		}

		path, err := os.Getwd()
		emperror.Panic(err)

		path, err = filepath.Abs(path + string(os.PathSeparator) + args[0])
		emperror.Panic(err)

		makeStruct(path, viper.GetStringMap("structure"))
	},
}

func init() {
	Command.Flags().StringP("config", "c", "", "config path")
}

func makeStruct(path string, structure map[string]interface{}) {
	emperror.Panic(helpers.MkDir(path))

	for name, sub := range structure {
		switch t := sub.(type) {
		case string:
			switch t {
			case "emptyfile":
				emperror.Panic(generators.EmptyFile(name, path))
			case "gofile":
				emperror.Panic(generators.EmptyGoFile(name, path))
			case "legostarter":
				emperror.Panic(generators.NewLeGoStarter(helpers.Path(path, name)))
			case "legoservice":
				emperror.Panic(generators.NewLegoService(helpers.Path(path, name)))
			default:
				emperror.Panic(errors.New(t + ": file command not found"))
			}
		case map[string]interface{}:
			makeStruct(helpers.Path(path, name), t)
		}
	}
}
