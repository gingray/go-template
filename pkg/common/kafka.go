package common

import "github.com/twmb/franz-go/pkg/kgo"

type KafkaConfig struct {
	Brokers       []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	ConsumerGroup string   `env:"KAFKA_CONSUMER_GROUP" envDefault:"test"`
	ConsumeTopics []string `env:"KAFKA_CONSUME_TOPICS" envDefault:"test"`
}

func NewKafka(cfg *KafkaConfig) (*kgo.Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(cfg.ConsumeTopics...),
	)
	return client, err
}
