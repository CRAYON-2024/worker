package cmd

import (
	"github.com/CRAYON-2024/worker/bootstrap"
	"github.com/spf13/cobra"
)

func Exec() {
	var (
		root      = &cobra.Command{}
		container = bootstrap.NewContainer()
	)

	root.AddCommand(
		WorkerCommand(container),
		RestCommand(container),
		ProcessCommand(container),
	)

	root.Execute()
}
