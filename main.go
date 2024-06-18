package main

import (
	"context"

	"github.com/CRAYON-2024/worker/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{}

	root.AddCommand(cmd.WorkerCommand())

	root.ExecuteContext(context.Background())
}
