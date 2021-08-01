package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/bestpilotingalaxy/ws-chat/internal/server"
	"github.com/bestpilotingalaxy/ws-chat/internal/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	cfg := config.New()
	config.SetupLogger("DEBUG")
	srv := server.New(cfg.Server)

	router := fiber.New()
	router.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go srv.Run()

	router.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		adapter := transport.NewWSAdapter(conn, srv.Broadcast, srv.Unregister)
		// Register the client
		srv.Register <- adapter
		adapter.Listen()
	}))
	log.Fatal(router.Listen("0.0.0.0:3000"))
}
