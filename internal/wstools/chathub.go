package wstools

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/bestpilotingalaxy/ws-chat/internal/packets/jrpc"
)

// ChatHub is hub for all connections
type ChatHub struct {
	JRPCRouter *jrpc.Router
	// pool of adapters binded to websock connections
	ClientPool map[*Adaptor]struct{}
	Register   chan *Adaptor
	// TODO: specify broadcast message type
	Broadcast  chan string
	JRPCchan   chan *jrpc.Request
	Unregister chan *Adaptor
	// using to stop hub Run() loop and close channels
	Interrupt chan struct{}
	// TODO: add hub configuration options
	cfg *config.Config
}

// NewHub gets config and creates ChatHub instance
func NewHub(cfg interface{}, jrpcR *jrpc.Router) *ChatHub {
	return &ChatHub{
		JRPCRouter: jrpcR,
		ClientPool: make(map[*Adaptor]struct{}),
		Register:   make(chan *Adaptor),
		Broadcast:  make(chan string),
		JRPCchan:   make(chan *jrpc.Request),
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
			hub.BroadcastToAll(message)

		case adapter := <-hub.Unregister:
			hub.unregister(adapter)

		case jrpcReq := <-hub.JRPCchan:
			hub.processJRPC(jrpcReq)

		case <-hub.Interrupt:
			hub.close()
			break SERVING_LOOP
		}
	}
}

// BroadcastToAll send the message to all clients in pool
func (hub *ChatHub) BroadcastToAll(message string) {
	log.Debug("message received: ", message)
	for adapter := range hub.ClientPool {
		go adapter.Write(message)
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

// Process jrpc request, check if called method exists and run it
func (hub *ChatHub) processJRPC(req *jrpc.Request) {
	log.Debug("recieved JRPC call: ", req)
	if err := hub.JRPCRouter.CheckMethodExist(req); err != nil {
		log.Error("unsupported function called: %s", err)
		return
	}
	go hub.JRPCRouter.Process(req)
}

// Close hub unbuffered channels
func (hub *ChatHub) close() {
	close(hub.Broadcast)
	close(hub.Register)
	close(hub.Unregister)
	log.Info("hub channels closed...")
}
