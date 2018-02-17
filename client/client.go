package client

import (
	"fmt"
)

type Client struct {
	data     []byte
	Sections map[string][]string
}

func New() *Client {
	return new(Client)
}

func (c *Client) Write(data []byte) (int, error) {
	c.data = data
	err := c.Parse(data)
	return len(data), err
}

func (c *Client) Parse(data []byte) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) Tree() map[string][]string {
	return c.Sections
}
