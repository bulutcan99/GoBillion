package consumer

import (
	"errors"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

const (
	_retryTimes     = 2
	_backOffSeconds = 2
)

type KafkaBrokers []string

var (
	ErrFailedToStartConsumer = errors.New("failed to start Kafka consumer")
)

func NewKafkaConsumer(brokers KafkaBrokers, groupID string) (*sarama.ConsumerGroup, error) {
	var (
		consumer        *sarama.ConsumerGroup
		connectionCount int64
	)
	config := sarama.NewConfig()

	for {
		conn, err := sarama.NewConsumerGroup(brokers, groupID, config)
		if err != nil {
			zap.S().Error("Failed to attach to producer...", err)
			connectionCount++
		} else {
			consumer = &conn
			break
		}

		if connectionCount > _retryTimes {
			zap.S().Error("Failed to connect after retries", err)
			return nil, ErrFailedToStartConsumer
		}

		zap.S().Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)
		continue

	}

	zap.S().Info("Kafka consumer started successfully")
	return consumer, nil
}
