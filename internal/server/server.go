package server

import (
	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/internal/client"
	"github.com/gofiber/websocket/v2"
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
			log.Info("connection registered")

		case message := <-s.Broadcast:
			log.Debug("message received: ", message)

			// Send the message to all clients
			for connection := range s.ClientPool {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Error("write error: ", err)

					s.Unregister <- connection
					log.Error("Closing connection: ", connection)
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}
		case connection := <-s.Unregister:
			// Remove the client from the hub
			delete(s.ClientPool, connection)

			log.Info("connection unregistered")
		}
	}
}
