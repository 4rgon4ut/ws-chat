package api

import (
	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Router ...
type Router struct {
	*fiber.App
	Config *config.Server
}

// NewRouter ...
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

// RunAPI ...
func (r *Router) RunAPI() {
	if err := r.Listen("0.0.0.0:" + r.Config.Port); err != nil {
		log.Fatalf("Cant Start server due: %s", err)
	}
}
