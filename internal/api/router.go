package api

import (
	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/gofiber/fiber/v2"
)

// Router ...
type Router struct {
	*fiber.App
	Config *config.Server
}

// NewRouter ...
func NewRouter(c *config.Server) *Router {
	r := fiber.New()
	return &Router{
		r,
		c,
	}
}

// RunAPI ...
func (r *Router) RunAPI() {
	if err := r.Listen("0.0.0.0:" + r.Config.Port); err != nil {
		log.Fatalf("Cant Start server due: %s", err)
	}

}
