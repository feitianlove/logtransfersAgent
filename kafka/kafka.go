package kafka

import (
	"github.com/feitianlove/golib/kafka"
)

func InitProduct(config kafka.Kafka) (*kafka.KProduct, error) {
	product, err := kafka.NewKafkaProduct(&config)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func InitConsumer(config kafka.Kafka) (*kafka.KConsumer, error) {
	consumer, err := kafka.NewKafkaConsumer(&config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
