package bootstrap

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
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
	Ctx            context.Context
	dBR            *pgx.Conn
	dBW            *pgx.Conn
	redis          *redis.Client
	producer       *kafka.Writer
	consumer       *kafka.Reader
	workerConsumer *kafka.Reader
	trace          *sdktrace.TracerProvider
}

func NewContainer() (res *Container) {
	res = cont
	if res == nil {
		res = &Container{
			Ctx: context.Background(),
		}
	}

	res.initTracer()

	return res
}

func (c *Container) Terminate() error {
	if c.consumer != nil {
		log.Println("closing kafka consumer", c.consumer.Close())
	}

	if c.producer != nil {
		log.Println("closing kafka producer", c.producer.Close())
	}

	return nil
}
