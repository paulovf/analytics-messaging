package rabbitmq_test

import (
	"context"
	"testing"

	"github.com/paulovf/analytics-messaging/rabbitmq"
	"github.com/stretchr/testify/assert"
)

func TestNewClient_ConnectionFailure(t *testing.T) {
	client, err := rabbitmq.NewClient("amqp://invalid_user:invalid_pass@localhost:9999/")

	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to connect to RabbitMQ")
}

func TestPublish_UninitializedClient(t *testing.T) {
	client := &rabbitmq.Client{}

	err := client.Publish(context.Background(), "test_exchange", "test_key", map[string]string{"payload": "test"})
	assert.Error(t, err)
}
