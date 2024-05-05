package courier

import (
	"context"

	"github.com/gojek/courier-go"
)

// Publish publishes a message to the specified topic using the Courier client.
func (c *Client) Publish(ctx context.Context, topic string, message []byte) error {
	if err := c.Client.Publish(ctx, topic, message, courier.QOSOne, courier.Retained(true)); err != nil {
		return err
	}

	return nil
}

// Subscribe allows to subscribe to messages from an MQTT broker
func (c *Client) Subscribe(ctx context.Context, topic string, fn func(ctx context.Context, data []byte) error) {
	// TODO: implement subscribe with custom marshaler
	// c.Client.Subscribe(ctx, topic, func(ctx context.Context, ps courier.PubSub, m *courier.Message) {
	// fn(ctx,m.DecodePayload())
	// })
}
