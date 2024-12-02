package server

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/joseph-gunnarsson/chatter-go/internal/errors"
)

type chatter struct {
	name      string
	conn      net.Conn
	ip        string
	publicKey rsa.PublicKey
}

func (c *chatter) handleCommunication(server *Server) {
	defer c.conn.Close()
	buf := make([]byte, 1024)
	commands := defineCommands()

	for {
		n, err := c.conn.Read(buf)

		if err != nil {
			errors.HandleError(err)
			break
		}

		log.Printf("Message recived from %s", c.ip)

		requestString := strings.TrimSpace(string(buf[:n]))

		if strings.HasPrefix(requestString, "/") {
			parts := strings.Split(requestString, " ")

			cmd, isCommand := commands[parts[0]]

			if isCommand {
				cmd.Handler(server, c, parts[1:])
				continue
			}
		}

		server.broadcastChatterMessage(c, requestString)
	}
}

func (c *chatter) writeMessage(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	formattedMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	_, err := c.conn.Write([]byte(formattedMessage))
	errors.HandleError(err)
}
