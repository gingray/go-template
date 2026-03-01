package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gingray/go-template/pkg/app"
	"github.com/gingray/go-template/pkg/config"
)

type BasicConsumer struct {
	logger config.Logger
	topic  string
}

func NewBasicConsumer(app *app.App, topic string) *BasicConsumer {
	return &BasicConsumer{logger: app.Logger, topic: topic}
}

func (b *BasicConsumer) Topic() string {
	return b.topic
}

func (b *BasicConsumer) Consume(ctx context.Context, key string, value []byte) error {
	var jsonValue interface{}
	err := json.Unmarshal(value, &jsonValue)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	b.logger.Info("basic consumer", "key", key, "value", jsonValue)

	return nil
}
