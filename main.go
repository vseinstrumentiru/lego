package main

import (
	"emperror.dev/emperror"
	"github.com/spf13/cobra"

	"github.com/vseinstrumentiru/lego/gen/commands/bootstrap"
	"github.com/vseinstrumentiru/lego/gen/commands/service"
)

func main() {
	root := &cobra.Command{}

	root.AddCommand(
		bootstrap.Command,
		service.Command,
	)

	emperror.Panic(root.Execute())
}
