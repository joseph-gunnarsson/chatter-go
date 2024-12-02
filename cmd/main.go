package main

import (
	"flag"
	"log"
	"os"

	"github.com/joseph-gunnarsson/chatter-go/internal/client"
	"github.com/joseph-gunnarsson/chatter-go/internal/errors"
	"github.com/joseph-gunnarsson/chatter-go/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Error: No command provided. Please specify the command.")
	}
	switch os.Args[1] {
	case "host":
		hostCommand()
	case "connect":
		connectCommand()
	default:
		log.Println("Error: Command not found")
		os.Exit(1)
	}
}

func hostCommand() {
	host := flag.String("host", "localhost", "Host address to bind the server")
	port := flag.String("port", "8888", "Port for the server (without ':')")
	password := flag.String("password", "", "Password for server authentication")
	flag.CommandLine.Parse(os.Args[2:])

	s := server.NewServer(*host, *port, *password)
	s.Run()
}

func connectCommand() {
	host := flag.String("host", "localhost", "Server host address")
	port := flag.String("port", "8888", "Server port to connect to")
	password := flag.String("password", "", "Password for server authentication")
	flag.CommandLine.Parse(os.Args[2:])

	client := client.NewClient(*host, *port, *password)
	err := client.Run()
	errors.HandleError(err)
}
