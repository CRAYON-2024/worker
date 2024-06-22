package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/CRAYON-2024/worker/bootstrap"
	"github.com/CRAYON-2024/worker/internal/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ProcessCommand(container *bootstrap.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "process",
		Short: "Consume data from server",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				wg          = &sync.WaitGroup{}
				chn         = make(chan []byte)
				wrkr        = worker.NewWorker(container.GetKafkaConsumer(), wg, chn, viper.GetString("kafka.topic.handson"))
				ctx, cancel = context.WithCancel(context.Background())
			)

			go wrkr.Process(ctx)
			if err := wrkr.Dispatch(ctx); err != nil {
				log.Fatalln("failed to run dispatcher", err)
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
			<-quit

			cancel()
			return nil
		},
	}
}
