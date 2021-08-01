package client

import (
	"github.com/bestpilotingalaxy/ws-chat/internal/transport"
)

// Client ...
type Client struct {
	// The websocket connection.
	wsAdapter *transport.WSAdapter
}

// New ...
func New() *Client {
	return &Client{}
}

// func Connect() {
// 	websocket.New()

// }
