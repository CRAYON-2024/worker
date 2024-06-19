package bootstrap

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func (c *Container) connectConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	// NewConsumer creates a new consumer using the given broker addresses and configuration
	conn, err := sarama.NewConsumer(viper.GetStringSlice("kafka.consumer.broker"), config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Container) GetKafkaConsumer() sarama.Consumer {
	if c.consumer != nil {
		return c.consumer
	}

	cons, err := c.connectConsumer()
	if err != nil {
		log.Fatalln("failed to connect to consumer", err)
	}
	c.consumer = cons
	return c.consumer
}
