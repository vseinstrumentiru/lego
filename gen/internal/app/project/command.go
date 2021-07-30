package project

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cfg := new(config)

	cmd := &cobra.Command{
		Use:  "new [name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.Name = args[0]
			_, err := create(cfg)

			return err
		},
	}

	cmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "show logs")
	cmd.Flags().StringVarP(&cfg.GitRemotePath, "remote", "r", "", "git remote repository url")
	cmd.Flags().BoolVarP(&cfg.HasGit, "use-git", "g", false, "init git repository")
	cmd.Flags().BoolVarP(&cfg.UseProtobuf, "use-proto", "p", false, "use protobufs")
	cmd.Flags().BoolVarP(&cfg.UseGraphql, "use-graphql", "q", false, "use graphql schema")

	return cmd
}
