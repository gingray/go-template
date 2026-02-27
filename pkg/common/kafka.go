package common

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConfig struct {
	Brokers       []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	ConsumerGroup string   `env:"KAFKA_CONSUMER_GROUP" envDefault:"test"`
	ConsumeTopics []string `env:"KAFKA_CONSUME_TOPICS" envDefault:"test"`
}

func (a *App) WithKafka(cfg *KafkaConfig) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(cfg.ConsumeTopics...),
	)
	a.Kafka = client
	a.AddReadyHandler(func(ctx context.Context) error {
		return client.Ping(ctx)
	})
	a.AddShutdownHandler(func(ctx context.Context) error {
		closeCh := make(chan struct{})
		go func() {
			client.Close()
			close(closeCh)
		}()
	
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-closeCh:
			return nil
		}
	})
	return err
}
