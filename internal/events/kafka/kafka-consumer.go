package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"project-survey-stats-writer/internal/events/contracts"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
}

func InitConsumer(url string) (contracts.IMessageConsumer, error) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  url,
		"group.id":           "my-kafkaConsumer-group",
		"auto.offset.reset":  "latest",
		"enable.auto.commit": "false",
	})
	if err != nil {
		return nil, err
	}

	topic := "my-topic" // Replace with your Kafka topic
	err = kafkaConsumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, fmt.Errorf("error subscribing to topic: %v", err)
	}

	consumer := &KafkaConsumer{kafkaConsumer}
	return consumer, nil
}

func (k *KafkaConsumer) ConsumeMessage() ([]byte, error) {
	ev := k.consumer.Poll(100)
	switch e := ev.(type) {
	case *kafka.Message:
		return e.Value, nil
	case kafka.Error:
		return nil, e
	}

	return nil, nil
}

func (k *KafkaConsumer) CloseConnection() error {
	return k.consumer.Close()
}
