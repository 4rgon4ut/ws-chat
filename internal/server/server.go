package server

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/bestpilotingalaxy/ws-chat/internal/transport"
	"github.com/gofiber/websocket/v2"
)

var (
	wg  sync.WaitGroup
	mux sync.Mutex
)

// Server ...
type Server struct {
	ClientPool map[*transport.WSAdapter]struct{}
	Register   chan *websocket.Conn
	// TODO: specify broadcast message type
	Broadcast  chan []byte
	Unregister chan *transport.WSAdapter
	cfg        *config.Config
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
			s.register(connection)
		case message := <-s.Broadcast:
			log.Debug("message received: ", message)
			// Send the message to all clients
			for adapter := range s.ClientPool {
				adapter.SendCh <- message
			}
		case adapter := <-s.Unregister:
			// Remove the client from the pool
			delete(s.ClientPool, adapter)
			log.Info("connection unregistered:", adapter)
		}
	}
}

func (s *Server) register(conn *websocket.Conn) {
	adapterWriteCh := make(chan []byte)
	adapter := transport.NewWSAdapter(conn, s.Broadcast, adapterWriteCh)
	s.ClientPool[adapter] = struct{}{}
	log.Info("connection registered:", conn)
	adapter.Serve(s.Unregister)
}
