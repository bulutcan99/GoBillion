package producer

import (
	"errors"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

type KafkaBrokers []string

var (
	ErrCannotConnectKafka = errors.New("cannot connect to Kafka")
)

func NewKafkaProducer(brokers KafkaBrokers) (*sarama.SyncProducer, error) {
	var (
		producer        *sarama.SyncProducer
		connectionCount int64
	)

	zap.S().Info("Kafka brokers: ", brokers)

	config := sarama.NewConfig()
	config.Producer.Retry.Max = _retryTimes
	config.Producer.Retry.Backoff = _backOffSeconds * time.Second

	for {
		conn, err := sarama.NewSyncProducer(brokers, config)
		if err != nil {
			zap.S().Error("Failed to connect to Kafka...", err)
			connectionCount++
		} else {
			producer = &conn
			break
		}

		if connectionCount > _retryTimes {
			zap.S().Error("Failed to connect after retries", err)
			return nil, ErrCannotConnectKafka
		}

		zap.S().Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)
		continue
	}

	zap.S().Info("Connected to Kafka ")
	return producer, nil
}
