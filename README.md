
# Chatter Go

Chatter go is a TCP-based chat application built with Go, offering real-time communication between connected clients. The project consists of two core components: a server for handling connections and managing communication, and a client for user interaction and message exchange.



## Features

- Host a chat server with ease
- Connect to chat servers securely
- Optional password authentication
- Flexible host and port configuration


## Installation

To install chatter go, first, clone the repository and build the tool:

```bash
    git clone https://github.com/joseph-gunnarsson/chatter-go.git
    cd chatter-go
    go build -o chatter-go cmd/main.go
```

For easier access, you can move the compiled binary to your system's bin directory
 ```bash
    sudo mv chatter-go /usr/local/bin/
```   
This allows you to run the tool from anywhere in your terminal by simply typing:

```
    chatter-go
```
## Usage

Host a server

```golang
    ./chatter-go host [flags]
```

 #### Host flags
  - ```-host```: Server binding address (default: localhost)
  - ```-port```: Server port (default: 8888)
  - ```-password```: Optional server authentication password

Connect to server

```golang
    ./chatter-go connect [flags]
```

   #### Host flags
  - ```-host```: Target server address (default: localhost)
  - ```-port```: Server connection port (default: 8888)
  - ```-password```: Server password

Example:

```golang
    ./chatter-go connect -host example.com -port 9999 -password secretkey
```