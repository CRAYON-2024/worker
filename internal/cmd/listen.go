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

var (
	topic string
)

func ListenCommand(container *bootstrap.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listen",
		Short: "Consume data from server",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				wg          = &sync.WaitGroup{}
				chn         = make(chan []byte)
				ctx, cancel = context.WithCancel(context.Background())
				wkr         *worker.Worker
			)



			switch topic {
			case "handson":
				wkr = worker.NewWorker(container.GetKafkaConsumer(), wg, chn, viper.GetString("kafka.topic.handson"))
			case "worker":
				wkr = worker.NewWorker(container.GetKafkaWorkerConsumer(), wg, chn, viper.GetString("kafka.topic.worker"))
			}

			go wkr.Process(ctx)
			if err := wkr.Dispatch(ctx); err != nil {
				log.Fatalln("failed to run dispatcher", err)
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
			<-quit

			cancel()
			return nil
		},
	}

	cmd.Flags().StringVarP(&topic, "topic", "t", "handson", "Define which topic you want to listen")

	return cmd
}
