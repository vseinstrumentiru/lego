package app

import "github.com/vseinstrumentiru/lego/v2/internal/execute"

func command(r *runtime) {
	r.container.Execute(execute.RunCommands)
}
