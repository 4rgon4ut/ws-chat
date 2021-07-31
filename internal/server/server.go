package server

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/internal/client"
	"github.com/gofiber/websocket/v2"
)

var (
	wg  sync.WaitGroup
	mux sync.Mutex
)

// Server ...
type Server struct {
	ClientPool map[*websocket.Conn]*client.Client
	Register   chan *websocket.Conn
	// TODO: specify broadcast message type
	Broadcast  chan string
	Unregister chan *websocket.Conn
}

// New gets config and creates server instance
func New(cfg interface{}) *Server {
	return &Server{}
}

// Run ...
func (s *Server) Run() {
	for {
		select {
		case connection := <-s.Register:
			s.ClientPool[connection] = client.New()
			log.Info("connection registered:", connection)

		case message := <-s.Broadcast:
			log.Debug("message received: ", message)
			// Send the message to all clients
			for connection := range s.ClientPool {
				wg.Add(1)
				go s.sendMsg(connection, message)
			}
		case connection := <-s.Unregister:
			// Remove the client from the pool
			delete(s.ClientPool, connection)
			log.Info("connection unregistered:", connection)
		}
	}
}

//
func (s *Server) sendMsg(c *websocket.Conn, message string) {
	defer wg.Done()
	if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Error("write error: ", err)
		s.Unregister <- c
		log.Error("closing connection: ", c)
		c.WriteMessage(websocket.CloseMessage, []byte{})
		c.Close()
	}
}
