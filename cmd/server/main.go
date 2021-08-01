package main

import (
	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/bestpilotingalaxy/ws-chat/internal/api"
	"github.com/bestpilotingalaxy/ws-chat/internal/wstools"
)

func main() {
	cfg := config.New()
	config.SetupLogger(cfg.LogLevel)

	api := api.NewRouter(&cfg.Server)
	wsHub := wstools.NewHub(cfg.Server)
	wsHub.AddRoutes(api.App)

	go wsHub.Run()

	api.RunAPI()
}
