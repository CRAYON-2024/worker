package worker

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

type Worker struct {
	Consumer *kafka.Reader
	chn      chan []byte
	topic    string
	wg       *sync.WaitGroup
}

func NewWorker(
	consumer *kafka.Reader,
	wg *sync.WaitGroup,
	chn chan []byte,
	topic string,
) *Worker {
	return &Worker{
		Consumer: consumer,
		wg:       wg,
		chn:      chn,
		topic:    topic,
	}
}
