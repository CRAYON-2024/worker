package bootstrap

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	cont *Container
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("failed to : ", err)
	}
}

type Container struct {
	ctx      context.Context
	dBR      *pgx.Conn
	dBW      *pgx.Conn
	consumer sarama.Consumer
	producer sarama.SyncProducer
	redis    *redis.Client
	trace    *sdktrace.TracerProvider
}

func NewContainer() (res *Container) {
	res = cont
	if res == nil {
		res = &Container{
			ctx: context.Background(),
		}
	}

	res.initTracer()
	return res
}

func (c *Container) Terminate() error {
	if c.consumer != nil {
		log.Println("closing kafka consumer", c.consumer.Close())
	}

	return nil
}
