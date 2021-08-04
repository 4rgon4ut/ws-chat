package main

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/bestpilotingalaxy/ws-chat/config"
	"github.com/bestpilotingalaxy/ws-chat/internal/api"
	"github.com/bestpilotingalaxy/ws-chat/internal/wstools"
)

func main() {
	cfg := config.New()
	config.SetupLogger(cfg.LogLevel)

	apiServer := api.NewRouter(&cfg.Server)
	wsHub := wstools.NewHub(cfg.Server)
	apiServer.AddWSRoutes(wsHub)

	go wsHub.Run()
	go apiServer.RunAPI()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	wsHub.Interrupt <- struct{}{}
	apiServer.Shutdown()
	log.Info("Good bye!")
	os.Exit(0)
}
