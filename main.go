package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/CRAYON-2024/worker/internal/cmd"
)

func main() {
	root := &cobra.Command{}

	root.AddCommand(cmd.WorkerCommand())

	root.ExecuteContext(context.Background())
}
