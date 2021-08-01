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

	return &Server{
		ClientPool: make(map[*transport.WSAdapter]struct{}),
		Register:   make(chan *websocket.Conn),
		Broadcast:  make(chan []byte),
		Unregister: make(chan *transport.WSAdapter),
	}
}

// Run ...
func (s *Server) Run() {
	for {
		select {
		//
		case connection := <-s.Register:
			log.Info("Recieved connection")
			s.register(connection)
		//
		case message := <-s.Broadcast:
			log.Debug("message received: ", message)
			// Send the message to all clients
			for adapter := range s.ClientPool {
				select {
				case adapter.SendCh <- message:
				default:
					close(adapter.SendCh)
					delete(s.ClientPool, adapter)
				}
			}
		case adapter := <-s.Unregister:
			// Remove the client from the pool
			if _, ok := s.ClientPool[adapter]; ok {
				delete(s.ClientPool, adapter)
				close(adapter.SendCh)
			}
		}
	}
}

func (s *Server) register(conn *websocket.Conn) {
	adapterWriteCh := make(chan []byte)
	adapter := transport.NewWSAdapter(conn, s.Broadcast, adapterWriteCh, s.Unregister)
	s.ClientPool[adapter] = struct{}{}
	log.Info("connection registered:", conn)
	adapter.Serve()
}
