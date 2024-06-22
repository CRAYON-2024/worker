package cmd

import (
	"github.com/CRAYON-2024/worker/bootstrap"
	"github.com/CRAYON-2024/worker/internal/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RestCommand(container *bootstrap.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "Run a web server service",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.WithField("component", "rest").Info("running rest")

			api.APIRouter(container)

			return nil
		},
	}
}
