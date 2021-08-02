package wstools

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

// AddRoutes add websocket routes to server engine router
func (hub *ChatHub) AddRoutes(app *fiber.App) *fiber.App {
	// middleware chacks if connection upgraded
	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})
	// websocket connection hadler
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		// creates adaptor for websocket connection to serve it
		adaptor := NewAdaptor(conn, hub.Broadcast, hub.Unregister)
		// Register the client in hub
		hub.Register <- adaptor
		adaptor.Listen()
	}))
	return app
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
