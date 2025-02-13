package external

import (
	"context"
	"ecommerce-payments/helpers"
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

func (e External) ProduceKafkaMessage(ctx context.Context, topic string, data []byte) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	brokers := strings.Split(helpers.GetEnv("KAFKA_BROKERS", "localhost:29092,localhost:29093,localhost:29094"), ",")

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return errors.Wrap(err, "failed to comunicate with kafka brokers")
	}

	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := producer.SendMessage(message)

	if err != nil {
		return errors.Wrap(err, "failed to produce message to kafka")
	}

	helpers.Logger.Info(fmt.Sprintf("successfuly to produce on topic %s, partition %d, offset %d", topic, partition, offset))

	return nil
}
