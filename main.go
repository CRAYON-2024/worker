package main

import (
	"context"

	"github.com/CRAYON-2024/worker/bootstrap"
	"github.com/CRAYON-2024/worker/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {
	conn := bootstrap.NewContainer()

	conn.GetDBRead()

	root := &cobra.Command{}

	root.AddCommand(cmd.WorkerCommand())

	root.ExecuteContext(context.Background())
}
