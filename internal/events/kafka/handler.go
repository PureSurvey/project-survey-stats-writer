package kafka

import (
	"github.com/IBM/sarama"
	"project-survey-stats-writer/internal/events/model"
)

type handler struct {
	ready    *chan bool
	messages *chan model.Message
}

func (h *handler) Setup(sarama.ConsumerGroupSession) error {
	close(*h.ready)
	return nil
}

func (h *handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			*h.messages <- model.Message{Topic: message.Topic, Value: message.Value}
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
