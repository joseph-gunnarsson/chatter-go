package client

import (
	"fmt"
	"net"

	"github.com/joseph-gunnarsson/chatter-go/internal/errors"
)

type Client struct {
	host           string
	port           string
	conn           net.Conn
	serverPassowrd string
}

var messages chan string = make(chan string)

func NewClient(host, port, password string) *Client {
	return &Client{
		host:           host,
		port:           port,
		serverPassowrd: password,
	}
}

func (c *Client) Run() error {
	address := fmt.Sprintf("%s:%s", c.host, c.port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		errors.HandleError(errors.ClientErrorExit{Message: "Failed when trying to connect to server"})
	}

	defer conn.Close()

	c.conn = conn
	c.verifyPassword()

	go c.sendMessages()

	newGUI(c)

	return nil
}

func (c *Client) verifyPassword() {
	_, err := c.conn.Write([]byte(c.serverPassowrd))
	if err != nil {
		errors.HandleError(errors.ClientErrorExit{Message: "Failed when trying to send password hash to server"})
		return
	}

	buf := make([]byte, 1024)
	_, err = c.conn.Read(buf)
	if err != nil {
		errors.HandleError(errors.ClientErrorExit{Message: "Failed to read server authentication response"})
		return
	}
}

func (c *Client) sendMessages() {
	for message := range messages {
		_, err := c.conn.Write([]byte(message))
		if err != nil {
			errors.HandleError(errors.ClientErrorExit{Message: "Failed to send message to server"})
		}
	}
}
