package bootstrap

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

func (c *Container) GetKafkaWorkerConsumer() *kafka.Reader {
	if c.consumer != nil {
		return c.consumer
	}

	c.consumer = c.newWorkerConsumer()

	return c.consumer
}

func (c *Container) newWorkerConsumer() *kafka.Reader {
	fmt.Println(viper.GetStringSlice("kafka.consumer.broker"))
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: viper.GetStringSlice("kafka.consumer.broker"),
		GroupID: viper.GetString("kafka.consumer.workerGroupID"),
		Topic:   viper.GetString("kafka.topic.worker"),
	})

	return consumer
}

func (c *Container) GetKafkaConsumer() *kafka.Reader {
	if c.consumer != nil {
		return c.consumer
	}

	c.consumer = c.newConsumer()

	return c.consumer
}

func (c *Container) newConsumer() *kafka.Reader {
	fmt.Println(viper.GetStringSlice("kafka.consumer.broker"))
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: viper.GetStringSlice("kafka.consumer.broker"),
		GroupID: viper.GetString("kafka.consumer.groupID"),
		Topic:   viper.GetString("kafka.topic.handson"),
	})

	return consumer
}

func (c *Container) GetKafkaProducer() *kafka.Writer {
	if c.producer != nil {
		return c.producer
	}

	c.producer = c.newProducer()

	return c.producer
}

func (c *Container) newProducer() *kafka.Writer {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: viper.GetStringSlice("kafka.consumer.broker"),
	})

	return producer
}
