package app

import (
	"context"

	"github.com/gingray/go-template/pkg/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (a *App) WithKafka(cfg *config.KafkaConfig) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.AllowAutoTopicCreation(),
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
