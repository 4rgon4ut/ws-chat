package config

import (
	"context"
	"os"

	"sync"

	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
)

var once sync.Once

// Server configuration
type Server struct {
	Port string `env:"SERVER_PORT"`
}

// Config global config struct
type Config struct {
	Server
	LogLevel string `env:"LOG_LEVEL"`
}

// New creates config struct and fills with env variables
func New() *Config {
	ctx := context.Background()
	c := &Config{}
	process := func() {
		if err := envconfig.Process(ctx, c); err != nil {
			log.Fatal(err)
		}
	}
	once.Do(process)
	return c
}

// SetupLogger ...
func SetupLogger(level string) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	switch level {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}
