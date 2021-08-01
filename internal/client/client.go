package client

import (
	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
)

// Client ...
type Client struct {
	// The websocket connection.
	conn *websocket.Conn
}

// New ...
func New() *Client {
	log.Info()
	return &Client{}
}
