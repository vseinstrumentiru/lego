package main

import (
	"emperror.dev/emperror"
	"github.com/spf13/cobra"

	"github.com/vseinstrumentiru/lego/gen/commands/service"
	"github.com/vseinstrumentiru/lego/gen/commands/structure"
)

func main() {
	root := &cobra.Command{}

	root.AddCommand(
		structure.Command,
		service.Command,
	)

	emperror.Panic(root.Execute())
}
