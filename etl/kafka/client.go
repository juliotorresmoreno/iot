package kafka

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct {
	topic    string
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func NewKafkaClient(topic string) *KafkaClient {
	return &KafkaClient{
		topic: topic,
	}
}

func (k *KafkaClient) MakeProducer() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "host1:9092,host2:9092",
		"client.id":         hostname,
		"acks":              "all",
	})

	if err != nil {
		return err
	}

	k.producer = producer

	return nil
}

func (k *KafkaClient) MakeConsumer() error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "host1:9092,host2:9092",
		"group.id":          "foo",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		return err
	}

	k.consumer = consumer

	return nil
}

func (k KafkaClient) Pub(payload any) error {
	return nil
}

func (k *KafkaClient) Sub(fn func(ch string, data any)) error {
	var err error
	if k.consumer == nil {
		err = k.MakeConsumer()
	}
	if err != nil {
		return err
	}

	err = k.consumer.Subscribe(k.topic, nil)
	if err != nil {
		return err
	}
	c := k.consumer

	for {
		ev := c.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			// application-specific processing
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}
}
