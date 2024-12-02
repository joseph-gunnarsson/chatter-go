package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/joseph-gunnarsson/chatter-go/internal/errors"
)

type Server struct {
	name     string
	port     string
	password string
	chatters map[net.Addr]*chatter
}

func NewServer(server_name, server_port, server_password string) *Server {
	formattedPort := fmt.Sprintf(":%s", server_port)

	return &Server{
		name:     server_name,
		port:     formattedPort,
		password: server_password,
		chatters: make(map[net.Addr]*chatter),
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.port)

	errors.HandleError(err)

	log.Printf("Server is listening on port %s...", s.port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			errors.HandleError(errors.ServerError{Message: "failed to accept connection: " + err.Error()})
			continue
		}

		go s.handleConnection(conn)

	}
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Printf("New connection from %s", conn.RemoteAddr())

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		log.Printf("Failed to read password: %v", err)
		conn.Close()
		return
	}

	receivedPassword := strings.TrimSpace(string(buf[:n]))

	if receivedPassword != s.password {
		log.Printf("Invalid password from %s. Received: '%s', Expected: '%s'",
			conn.RemoteAddr(), receivedPassword, s.password)
		conn.Write([]byte(""))
		conn.Close()
		return
	}
	conn.Write([]byte("OK"))
	chatter := &chatter{
		name: fmt.Sprintf("Anonymous%d", len(s.chatters)+1),
		conn: conn,
		ip:   conn.RemoteAddr().String(),
	}

	s.chatters[conn.RemoteAddr()] = chatter
	go chatter.handleCommunication(s)
}

func (s *Server) broadcastChatterMessage(sender *chatter, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	formattedMessage := fmt.Sprintf("[%s] %s: %s\n",
		timestamp,
		sender.name,
		message,
	)

	messageBytes := []byte(formattedMessage)

	for addr, chatter := range s.chatters {
		if addr.String() != sender.ip {
			chatter.conn.Write(messageBytes)
		}
	}
}

func (s *Server) broadcastServerMessage(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	formattedMessage := fmt.Sprintf("[%s] %s\n",
		timestamp,
		message,
	)

	messageBytes := []byte(formattedMessage)

	for _, chatter := range s.chatters {
		chatter.conn.Write(messageBytes)
	}
}

func (s *Server) updateName(sender *chatter, name string) {
	oldName := sender.name
	sender.name = name

	updateMessage := fmt.Sprintf("%s has changed their name to %s.", oldName, name)
	s.broadcastServerMessage(updateMessage)
}

func (s *Server) quit(sender *chatter) {
	err := sender.conn.Close()

	errors.HandleError(err)
	delete(s.chatters, sender.conn.RemoteAddr())
	updateMessage := fmt.Sprintf("%s has left the server.", sender.name)
	s.broadcastServerMessage(updateMessage)
}

func (s *Server) kickChatter(chatter *chatter) {
	err := chatter.conn.Close()

	errors.HandleError(err)
	delete(s.chatters, chatter.conn.RemoteAddr())

	updateMessage := fmt.Sprintf("%s has been kicked from the server.", chatter.name)
	s.broadcastServerMessage(updateMessage)
}
