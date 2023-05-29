package natsclient

import (
	"github.com/nats-io/nats.go"
	"log"
)

type Client struct {
	nc *nats.Conn
}

func NewClient() *Client {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS server: %v", err)
	}

	return &Client{
		nc: nc,
	}
}

func (c *Client) Publish(subject string, data []byte) error {
	return c.nc.Publish(subject, data)
}

func (c *Client) Subscribe(subject string, callback func(msg *nats.Msg)) (*nats.Subscription, error) {
	return c.nc.Subscribe(subject, callback)
}
