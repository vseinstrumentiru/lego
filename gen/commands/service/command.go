package service

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vseinstrumentiru/lego/v2/gen/generators/interfaces"
)

var Command = &cobra.Command{
	Use:  "service",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path []string
		if len(args) == 0 {
			path = append(path, ".")
		} else {
			path = strings.Split(args[0], ",")
		}

		runtime, err := interfaces.Generate(interfaces.Config{
			Name:     "service",
			Mark:     "lego:service:contract",
			ImplMark: "lego:service",
			Paths:    path,
			Output:   "subpkg:suffix=service",
			FileName: "service.go",
		})

		if err != nil {
			return err
		}

		if hadErrs := runtime.Run(); hadErrs {
			os.Exit(1)
		}

		return nil
	},
}
