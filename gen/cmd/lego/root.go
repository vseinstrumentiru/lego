package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func register(cmds ...*cobra.Command) []interface{} {
	var res []interface{}
	for _, cmd := range cmds {
		res = append(res, registerCommand(cmd))
	}

	return res
}

func registerCommand(cmd *cobra.Command) interface{} {
	return func() Command {
		return Command{
			Command: cmd,
		}
	}
}

type Command struct {
	dig.Out
	Command *cobra.Command `group:"cmd"`
}
