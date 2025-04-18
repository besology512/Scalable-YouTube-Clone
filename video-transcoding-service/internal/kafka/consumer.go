package kafka

import (
	"context"
	"encoding/json"
	"log"

	"video-transcoding-service/internal/types"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	messageChan   chan *types.VideoMessage
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "transcode-group", config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		messageChan:   make(chan *types.VideoMessage),
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, topics []string) error {
	handler := &ConsumerGroupHandler{consumer: c}

	go func() {
		for err := range c.consumerGroup.Errors() {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	return c.consumerGroup.Consume(ctx, topics, handler)
}

func (c *Consumer) Messages() <-chan *types.VideoMessage {
	return c.messageChan
}

func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}

type ConsumerGroupHandler struct {
	consumer *Consumer
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var videoMsg types.VideoMessage
		if err := json.Unmarshal(message.Value, &videoMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		log.Printf("Received valid video message: %s", videoMsg.VideoURL)
		h.consumer.messageChan <- &videoMsg
		session.MarkMessage(message, "")
	}
	return nil
}
