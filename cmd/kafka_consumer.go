package cmd

import (
	"context"
	"ecommerce-payments/helpers"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type PaymentInitiateHandler struct {
	Dependency   Dependency
	TopicPayment string
	TopicRefund  string
}

func (h *PaymentInitiateHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PaymentInitiateHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PaymentInitiateHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		helpers.Logger.Infof("Received message: %s from partition %d", string(msg.Value), msg.Partition)

		switch msg.Topic {
		case h.TopicPayment:
			err := h.Dependency.PaymentAPI.InitiatePayment(msg.Value)
			if err != nil {
				helpers.Logger.Error("failed to process payment: ", err)
			}
		case h.TopicRefund:
			err := h.Dependency.PaymentAPI.RefundPayment(msg.Value)
			if err != nil {
				helpers.Logger.Error("failed to process payment: ", err)
			}
		default:
			helpers.Logger.Error("invalid topic: ", msg.Topic)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func ServeKafkaConsumerGroup() {
	d := dependencyInject()
	topicPayment := helpers.GetEnv("KAFKA_TOPIC_PAYMENT_INITIATE", "payment-initiation-topic")
	topicRefund := helpers.GetEnv("KAFKA_TOPIC_REFUND", "refund-topic")

	brokers := strings.Split(helpers.GetEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	groupID := helpers.GetEnv("KAFKA_CONSUMER_GROUP", "ecommerce-payment-group")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second * 1

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		helpers.Logger.Error("Failed to connect with kafka as consumer", err)
		return
	}
	defer consumerGroup.Close()

	handler := PaymentInitiateHandler{
		Dependency:   d,
		TopicPayment: topicPayment,
		TopicRefund:  topicRefund,
	}

	for {
		err := consumerGroup.Consume(context.Background(), []string{topicPayment, topicRefund}, &handler)
		if err != nil {
			helpers.Logger.Errorf("failed to consuming messages: %v", err)
		}
	}
}
