package newproject

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vseinstrumentiru/lego/v2/gen/internal/generators"
	"github.com/vseinstrumentiru/lego/v2/gen/internal/helpers"
)

func getDefaultStructure() map[string]interface{} {
	return map[string]interface{}{
		"api": map[string]interface{}{
			".gitkeep": "emptyfile",
		},
		"cmd": map[string]interface{}{
			"server": "legostarter",
		},
		"internal": map[string]interface{}{
			"app": map[string]interface{}{
				"newservice": "legoservice",
			},
			"platform": map[string]interface{}{
				".gitkeep": "emptyfile",
			},
		},
	}
}

func NewCommand(box *packr.Box) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "new [name] [path]",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetDefault("structure", getDefaultStructure())

			modPath := strings.Trim(args[0], "/")
			if len(modPath) == 0 {
				return errors.New("package name required")
			}
			projectName := modPath
			{
				segments := strings.Split(modPath, "/")
				if len(segments) > 1 {
					projectName = segments[len(segments)-1]
				}
			}

			if len(args) == 1 {
				args = append(args, projectName)
			}

			var err error
			{
				var configPath string
				if configPath, err = cmd.Flags().GetString("template"); err == nil && configPath != "" {
					viper.AddConfigPath(configPath)
					err = viper.ReadInConfig()
				}

				if err != nil {
					return err
				}
			}

			var path string
			if path, err = os.Getwd(); err != nil {
				return err
			}

			if path, err = filepath.Abs(path + string(os.PathSeparator) + args[1]); err != nil {
				return err
			}

			makeStruct(path, viper.GetStringMap("structure"))

			{
				file, err := box.Find("newproject/main.mk")
				if err != nil {
					return err
				}
				if err = helpers.CopyPaste(file, helpers.Path(path, "main.mk")); err != nil {
					return err
				}
			}

			{
				file, err := box.Find("newproject/Makefile")
				if err != nil {
					return err
				}
				if err = helpers.CopyPaste(file, helpers.Path(path, "Makefile")); err != nil {
					return err
				}
			}

			{
				file, err := box.Find("newproject/gitignore.dist")
				if err != nil {
					return err
				}
				if err = helpers.CopyPaste(file, helpers.Path(path, ".gitignore")); err != nil {
					return err
				}
			}

			{
				file, err := box.Find("newproject/golangci.yml")
				if err != nil {
					return err
				}
				if err = helpers.CopyPaste(file, helpers.Path(path, ".golangci.yml")); err != nil {
					return err
				}
			}

			{
				file, err := box.Find("newproject/config.tpl")
				if err != nil {
					return err
				}
				if err = helpers.CopyPasteTemplate(file, helpers.Path(path+"/cmd/server", "config.yaml"), struct{ Name string }{Name: projectName}); err != nil {
					return err
				}
			}

			{
				file, err := box.Find("newproject/config.tpl")
				if err != nil {
					return err
				}
				if err = helpers.CopyPasteTemplate(file, helpers.Path(path+"/cmd/server", "config.yaml.dist"), struct{ Name string }{Name: projectName}); err != nil {
					return err
				}
			}

			{
				file, err := box.Find("newproject/mod.tpl")
				if err != nil {
					return err
				}
				if err = helpers.CopyPasteTemplate(file, helpers.Path(path, "go.mod"), struct{ Path string }{Path: modPath}); err != nil {
					return err
				}
			}

			fmt.Printf("Project \"%s\" created. Make something great!\n", projectName)

			return nil
		},
	}

	cmd.Flags().StringP("template", "s", "", "path to config.yaml with project template")

	return cmd
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
