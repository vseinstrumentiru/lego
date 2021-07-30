package service

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cfg := new(config)

	cmd := &cobra.Command{
		Use:  "service [name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.Name = args[0]
			_, err := create(cfg)

			return err
		},
	}

	cmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "show logs")

	return cmd
}
