package di

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type Command struct {
	dig.Out
	Cmd *cobra.Command `group:"cmd"`
}

func NewCommand(cmd *cobra.Command) Command {
	return Command{
		Cmd: cmd,
	}
}

func ProvideCommand(cmd *cobra.Command) func() Command {
	return func() Command {
		return NewCommand(cmd)
	}
}
