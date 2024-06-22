package worker

import (
	"context"
	"log"
)

func (worker *Worker) Process(context context.Context) {
	for v := range worker.chn {
		log.Println(string(v))
	}
}
