// +build gen

package main

import (
	"fmt"
	"os"

	"emperror.dev/emperror"
	"github.com/spf13/cobra"

	"github.com/vseinstrumentiru/lego/v2/gen/commands/service"
	"github.com/vseinstrumentiru/lego/v2/gen/commands/structure"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	defer emperror.HandleRecover(emperror.ErrorHandlerFunc(func(err error) {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}))

	root := &cobra.Command{
		Use:     "LeGo",
		Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commitHash, buildDate),
	}

	root.AddCommand(
		structure.Command,
		service.Command,
	)

	emperror.Panic(root.Execute())
}
