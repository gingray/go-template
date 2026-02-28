package kafka

import (
	"context"
	"fmt"

	"github.com/gingray/go-template/pkg/app"
	"github.com/gingray/go-template/pkg/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Kafka struct {
	client *kgo.Client
	logger config.Logger
}

func NewKafka(app *app.App) *Kafka {
	return &Kafka{client: app.Kafka, logger: app.Logger}
}

func (k *Kafka) Name() string {
	return "kafka"
}

func (k *Kafka) Ready(ctx context.Context) error {
	return nil
}

func (k *Kafka) Run(ctx context.Context) error {
	for {
		fetches := k.client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return nil
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			return fmt.Errorf("fetch error: %v", errs)
		}

		fetches.EachRecord(func(r *kgo.Record) {
			k.logger.Info("kafka msg", "topic", r.Topic, "key", string(r.Key), "value", string(r.Value))
		})

		// commit offsets if using consumer group
		if err := k.client.CommitUncommittedOffsets(ctx); err != nil {
			return fmt.Errorf("failed to commit offsets: %w", err)
		}
	}
}

func (k *Kafka) Shutdown(ctx context.Context) error {
	k.client.Close()
	return nil
}
