package wstools

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
)

// ChatHub is hub for all connections
type ChatHub struct {
	// pool of adapters binded to websock connections
	ClientPool map[*Adaptor]struct{}
	Register   chan *Adaptor
	// TODO: specify broadcast message type
	Broadcast  chan string
	Unregister chan *Adaptor
	// using to stop hub Run() loop and close channels
	Interrupt chan struct{}
	// TODO: add hub configuration options
	cfg *config.Config
}

// NewHub gets config and creates ChatHub instance
func NewHub(cfg interface{}) *ChatHub {
	return &ChatHub{
		ClientPool: make(map[*Adaptor]struct{}),
		Register:   make(chan *Adaptor),
		Broadcast:  make(chan string),
		Unregister: make(chan *Adaptor),
		Interrupt:  make(chan struct{}, 1),
	}
}

// Run makes hub serving all connected clients in infinite loop
func (hub *ChatHub) Run() {
SERVING_LOOP:
	for {
		select {

		case adapter := <-hub.Register:
			hub.register(adapter)

		case message := <-hub.Broadcast:
			hub.broadcastAll(message)

		case adapter := <-hub.Unregister:
			hub.unregister(adapter)

		case <-hub.Interrupt:
			hub.close()
			break SERVING_LOOP
		}
	}
}

// Notify chat members with text msg
func (hub *ChatHub) notify(notification string) {
	go func() {
		hub.Broadcast <- notification
	}()
}

// Remove the client from the pool
func (hub *ChatHub) unregister(adaptor *Adaptor) {
	if _, ok := hub.ClientPool[adaptor]; ok {
		delete(hub.ClientPool, adaptor)
	}
}

// Add new client adapter to pool
func (hub *ChatHub) register(adaptor *Adaptor) {
	hub.ClientPool[adaptor] = struct{}{}
	log.Infof("client [%s] joined the pool", adaptor.Conn.RemoteAddr())
	hub.notify(fmt.Sprintf("%s joined", adaptor.Conn.RemoteAddr()))
}

// Send the message to all clients in pool
func (hub *ChatHub) broadcastAll(message string) {
	log.Debug("message received: ", message)
	for adapter := range hub.ClientPool {
		go adapter.Write(message)
	}
}

// Close hub unbuffered channels
func (hub *ChatHub) close() {
	close(hub.Broadcast)
	close(hub.Register)
	close(hub.Unregister)
	log.Info("hub channels closed...")

}
