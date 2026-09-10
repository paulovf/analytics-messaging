package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{conn: conn, channel: ch}, nil
}

func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) DeclareExchange(name, kind string) error {
	return c.channel.ExchangeDeclare(
		name,  // name
		kind,  // type (topic, direct, fanout)
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

func (c *Client) Publish(ctx context.Context, exchange, routingKey string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal event body: %w", err)
	}

	return c.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

func (c *Client) Consume(queue, consumerName string, handler func(body []byte) error) error {
	msgs, err := c.channel.Consume(
		queue,        // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("[RABBITMQ CONSUMER ERROR] Error processing message: %v", err)
			}
		}
	}()

	return nil
}

func (c *Client) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (c *Client) BindQueue(queue, exchange, routingKey string) error {
	return c.channel.QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // args
	)
}
