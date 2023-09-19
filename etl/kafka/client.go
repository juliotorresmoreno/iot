package kafka

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/juliotorresmoreno/iot/etl/config"
	"github.com/juliotorresmoreno/iot/etl/log"
)

var logger = log.Log

type KafkaClient struct {
	topic    string
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func NewKafkaClient(topic string) (*KafkaClient, error) {
	return &KafkaClient{
		topic: topic,
	}, nil
}

func (k *KafkaClient) MakeProducer() error {
	conf, _ := config.GetConfig()
	producer, err := kafka.NewProducer(conf.Kaftka.Producer)

	if err != nil {
		return err
	}

	k.producer = producer

	return nil
}

func (k *KafkaClient) MakeConsumer() error {
	conf, _ := config.GetConfig()
	consumer, err := kafka.NewConsumer(conf.Kaftka.Consumer)
	if err != nil {
		return err
	}

	k.consumer = consumer

	return nil
}

func (k *KafkaClient) Pub(payload any) error {
	var err error
	if k.producer == nil {
		err = k.MakeProducer()
	}
	if err != nil {
		return err
	}

	buff := bytes.NewBufferString("")
	json.NewEncoder(buff).Encode(payload)

	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: kafka.PartitionAny,
		},
		Value: buff.Bytes(),
	}, nil)
	return err
}

func (k *KafkaClient) Sub(fn func(ch string, data any) error) error {
	var err error
	if k.consumer == nil {
		err = k.MakeConsumer()
	}
	if err != nil {
		return err
	}

	logger.Info("Connecting to topic: " + k.topic)
	err = k.consumer.Subscribe(k.topic, func(c *kafka.Consumer, e kafka.Event) error {
		return nil
	})
	if err != nil {
		return err
	}

	for {
		msg, err := k.consumer.ReadMessage(100)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
