package worker

import (
	"context"
	"log"
)

func (worker *Worker) Dispatch(ctx context.Context) error {
	log.Println("listening to topic", worker.topic)

	defer func() {
		if err := worker.Consumer.Close(); err != nil {
			log.Fatalf("Failed to close consumer: %v", err)
		}
		log.Println("shutting down consumer")
	}()

	For:
	for {
		select {
		case <-ctx.Done():
			close(worker.chn)
			log.Println("dispatcher closing chn")
			break For
		default:
			msg, err := worker.Consumer.FetchMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					break For
				}
				log.Printf("Failed to fetch message: %v", err)
				continue
			}
			worker.chn <- msg.Value
			if err := worker.Consumer.CommitMessages(ctx, msg); err != nil {
				log.Printf("Failed to commit message: %v", err)
			}
		}
	}

	return nil
}
