package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"project-survey-stats-writer/internal/configuration"
	"project-survey-stats-writer/internal/events/model"
)

type Consumer struct {
	Ready    chan bool
	Messages chan model.Message

	consumer sarama.ConsumerGroup
	config   *configuration.EventsConfiguration
	handler  *handler
}

func NewConsumer(config *configuration.EventsConfiguration) *Consumer {
	return &Consumer{config: config}
}

func (c *Consumer) Init() error {
	var err error
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false

	c.consumer, err = sarama.NewConsumerGroup([]string{c.config.BrokerUrl}, c.config.Group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	c.Ready = make(chan bool)
	c.Messages = make(chan model.Message)
	c.handler = &handler{ready: &c.Ready, messages: &c.Messages}

	return nil
}

func (c *Consumer) Consume(ctx context.Context) error {
	topics := []string{c.config.TrackingEventsTopic, c.config.CompletionEventsTopic}
	for {
		if err := c.consumer.Consume(ctx, topics, c.handler); err != nil {
			close(c.Ready)
			return err
		}
		if ctx.Err() != nil {
			return nil
		}
		c.Ready = make(chan bool)
	}
}

func (c *Consumer) CloseConnection() error {
	return c.consumer.Close()
}
