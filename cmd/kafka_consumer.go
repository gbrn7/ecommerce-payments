package cmd

import (
	"ecommerce-payments/helpers"
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM/sarama"
)

func ServeKafkaConsumerPaymentInit() {
	brokers := strings.Split(helpers.GetEnv("KAFKA_BROKERS", "localhost:29092,localhost:29093,localhost:29094"), ",")
	topic := helpers.GetEnv("KAFKA_TOPIC_PAYMENT_INITIATE", "refund-topic")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		helpers.Logger.Error("failed to connect with kafka consumer payment init ", err)
		return
	}

	partitionNumberStr := helpers.GetEnv("KAFKA_TOPIC_PAYMENT_INITIATE_PARTITION", "3")
	partitionNumber, _ := strconv.Atoi(partitionNumberStr)
	for i := int32(0); i < int32(partitionNumber); i++ {
		partitionConsumer, err := consumer.ConsumePartition(topic, i, sarama.OffsetOldest)
		if err != nil {
			helpers.Logger.Errorf("failed to create consumer partition %d %s", i, err)
			return
		}

		for msg := range partitionConsumer.Messages() {
			fmt.Printf("Receive message: %s from partition %d\n", string(msg.Value), msg.Partition)
		}
	}

}

func ServeKafkaConsumerRefund() {
	brokers := strings.Split(helpers.GetEnv("KAFKA_BROKERS", "localhost:29092,localhost:29093,localhost:29094"), ",")
	topic := helpers.GetEnv("KAFKA_TOPIC_REFUND", "refund-topic")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		helpers.Logger.Error("failed to connect with kafka consumer refund ", err)
		return
	}

	partitionNumberStr := helpers.GetEnv("KAFKA_TOPIC_REFUND", "3")
	partitionNumber, _ := strconv.Atoi(partitionNumberStr)
	for i := int32(0); i < int32(partitionNumber); i++ {
		partitionConsumer, err := consumer.ConsumePartition(topic, i, sarama.OffsetOldest)
		if err != nil {
			helpers.Logger.Errorf("failed to create consumer partition %d %s", i, err)
			return
		}

		for msg := range partitionConsumer.Messages() {
			fmt.Printf("Receive message: %s from partition %d\n", string(msg.Value), msg.Partition)
		}
	}
}
