package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	value := fmt.Sprintf("message-%d", 0)
	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(value),
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
	err = k.consumer.SubscribeTopics([]string{k.topic}, func(c *kafka.Consumer, e kafka.Event) error {
		fmt.Println(e)
		return nil
	})
	if err != nil {
		return err
	}
	c := k.consumer

	for {
		ev := <-c.Events()
		fmt.Println("event:", ev)
	}
}
