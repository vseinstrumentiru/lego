package provide

import "github.com/spf13/cobra"

func RootCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "lego",
		Hidden: true,
	}
}
