package configs

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

func NewKafkaConsumer(topics []string) (*KafkaConsumer, error) {

	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9093",
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{Consumer: consumer}, nil
}

func (kc *KafkaConsumer) ConsumeMessages() {
	// Listening messages from topics
	for {
		msg, err := kc.Consumer.ReadMessage(-1)
		if err == nil {
			// Handling received message
			log.Printf("Received message from topic %s: %s\n", *msg.TopicPartition.Topic, string(msg.Value))
		} else {
			log.Printf("Error reading message: %v\n", err)
		}
	}
}
