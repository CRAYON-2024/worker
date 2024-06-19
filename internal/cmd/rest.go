package cmd

import (
	"github.com/CRAYON-2024/worker/bootstrap"
	"github.com/spf13/cobra"
)

func restCommand(cfg *bootstrap.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "Fetch API",
		Run: func(cmd *cobra.Command, args []string) {
			
		},
	}
}
