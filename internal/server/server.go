package server

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/bestpilotingalaxy/ws-chat/internal/transport"
)

// Server ...
type Server struct {
	ClientPool map[*transport.WSAdapter]struct{}
	Register   chan *transport.WSAdapter
	// TODO: specify broadcast message type
	Broadcast  chan string
	Unregister chan *transport.WSAdapter
	cfg        *config.Config
}

// New gets config and creates server instance
func New(cfg interface{}) *Server {

	return &Server{
		ClientPool: make(map[*transport.WSAdapter]struct{}),
		Register:   make(chan *transport.WSAdapter),
		Broadcast:  make(chan string),
		Unregister: make(chan *transport.WSAdapter),
	}
}

// Run ...
func (s *Server) Run() {
	for {
		select {

		case adapter := <-s.Register:
			s.register(adapter)

		case message := <-s.Broadcast:
			s.broadcastAll(message)

		case adapter := <-s.Unregister:
			s.unregister(adapter)
		}
	}
}

// Notify chat members with text msg
func (s *Server) notify(notification string) {
	go func() {
		s.Broadcast <- notification
	}()
}

// Remove the client from the pool
func (s *Server) unregister(adapter *transport.WSAdapter) {
	if _, ok := s.ClientPool[adapter]; ok {
		delete(s.ClientPool, adapter)
	}
}

// Add new client adapter to pool
func (s *Server) register(adapter *transport.WSAdapter) {
	s.ClientPool[adapter] = struct{}{}
	log.Infof("Client [%s] joined the pool", adapter)
	s.notify(fmt.Sprintf("%s joined", adapter.Conn.RemoteAddr()))
}

// Send the message to all clients in pool
func (s *Server) broadcastAll(message string) {
	log.Debug("message received: ", message)
	for adapter := range s.ClientPool {
		go adapter.Write(message)
	}
}
