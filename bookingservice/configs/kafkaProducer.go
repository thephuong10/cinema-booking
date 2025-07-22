package configs

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

const (
	MaxRetries = 3
	BackOff    = 2 * time.Second
)

type KafkaProducer struct {
	Producer *kafka.Producer
}

func NewKafkaProducer() (*KafkaProducer, error) {

	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9093",
		"acks":              "all",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{Producer: producer}, nil
}

func (kp *KafkaProducer) SendMessage(topic string, message string) error {
	deliveryChan := make(chan kafka.Event)

	var err error
	for retries := 0; retries <= MaxRetries; retries++ {
		err = kp.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, deliveryChan)

		if err != nil {
			log.Printf("Error producing message (attempt %d): %v", retries+1, err)
			if retries == MaxRetries {
				return fmt.Errorf("failed after %d attempts: %v", retries+1, err)
			}
			// Wait before retrying
			time.Sleep(BackOff) // 2 seconds backoff
			continue
		}

		// If message was successfully sent
		ev := <-deliveryChan
		msg := ev.(*kafka.Message)
		if msg.TopicPartition.Error != nil {
			log.Printf("Message delivery failed: %v", msg.TopicPartition.Error)
			if retries == MaxRetries {
				return msg.TopicPartition.Error
			}
			// Retry on failure
			time.Sleep(BackOff) // 2 seconds backoff
		} else {
			log.Printf("Message delivered to topic %s [%d] at offset %v\n",
				*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
			break
		}
	}

	return nil
}
