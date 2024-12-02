package server

import (
	"fmt"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Handler     func(s *Server, c *chatter, args []string)
}

func defineCommands() map[string]Command {
	return map[string]Command{
		"/name": {
			Name:        "/name",
			Description: "Change your username",
			Handler:     handleNameChange,
		},
		"/list": {
			Name:        "/list",
			Description: "List all connected chatters",
			Handler:     handleListChatters,
		},
		"/kick": {
			Name:        "/kick",
			Description: "Kick a chatter from the server (admin only)",
			Handler:     handleKick,
		},
		"/help": {
			Name:        "/help",
			Description: "Show available commands",
			Handler:     handleHelp,
		},
		"/quit": {
			Name:        "/quit",
			Description: "Disconnect from the server",
			Handler:     handleQuit,
		},
	}
}

func handleNameChange(s *Server, c *chatter, args []string) {
	if len(args) == 0 {
		c.writeMessage("usage: /name <new_username>")
		return
	}
	newName := args[0]
	s.updateName(c, newName)
}

func handleListChatters(s *Server, c *chatter, args []string) {
	var chattersList strings.Builder
	chattersList.WriteString("Connected Chatters:\n")
	for _, chatter := range s.chatters {
		chattersList.WriteString(fmt.Sprintf("- %s (%s)\n", chatter.name, chatter.ip))
	}
	c.conn.Write([]byte(chattersList.String()))
}

func handleKick(s *Server, c *chatter, args []string) {
	if len(args) == 0 {
		c.writeMessage("usage: /kick <ip_address>")
		return
	}
	kickIP := args[0]
	for addr, chatter := range s.chatters {
		if addr.String() == kickIP {
			s.kickChatter(chatter)
			c.writeMessage(fmt.Sprintf("Kicked chatter with IP: %s", kickIP))
		}
	}
}

func handleHelp(s *Server, c *chatter, args []string) {
	commands := defineCommands()
	var helpText strings.Builder
	helpText.WriteString("Available Commands:\n")
	for _, cmd := range commands {
		helpText.WriteString(fmt.Sprintf("%s - %s\n", cmd.Name, cmd.Description))
	}

	c.conn.Write([]byte(helpText.String()))
}

func handleQuit(s *Server, c *chatter, args []string) {
	c.writeMessage("Disconnecting from the server")
	s.quit(c)
}
