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
