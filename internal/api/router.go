package api

import (
	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/bestpilotingalaxy/ws-chat/internal/wstools"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

// Router ...
type Router struct {
	*fiber.App
	Config *config.Server
}

// NewRouter new fiber app with middlewares
func NewRouter(c *config.Server) *Router {
	app := fiber.New()

	// Default configuration fiber middlewares
	// https://docs.gofiber.io/api/middleware/recover
	app.Use(recover.New())
	// https://docs.gofiber.io/api/middleware/logger
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Europe/Moscow",
	}))

	return &Router{
		app,
		c,
	}
}

// RunAPI start listen specified <addr:port>
func (r *Router) RunAPI() {
	if err := r.Listen("0.0.0.0:" + r.Config.Port); err != nil {
		log.Fatalf("cant Start server due: %s", err)
	}
}

// AddWSRoutes add websocket routes to server engine router
func (r *Router) AddWSRoutes(hub *wstools.ChatHub) {
	// middleware checks if connection upgraded
	r.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})
	// websocket connection hadler
	r.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		// creates adaptor for websocket connection to serve it
		adaptor := wstools.NewAdaptor(conn, hub.Broadcast, hub.JRPCchan, hub.Unregister)
		// Register the client in hub
		hub.Register <- adaptor
		adaptor.Listen()
	}))
}
