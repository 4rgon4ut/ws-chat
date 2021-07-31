package config

import (
	"context"
	"log"
	"sync"

	"github.com/sethvargo/go-envconfig"
)

var once sync.Once

// Server configuration
type server struct {
	Port string `env:"SERVER_PORT"`
}

// Config ...
type Config struct {
	Server   server
	LogLevel string `env:"LOG_LEVEL"`
}

// New creates config struct and fills with env variables
func New() *Config {
	ctx := context.Background()
	c := &Config{}
	proccess := func() {
		if err := envconfig.Process(ctx, c); err != nil {
			log.Fatal(err)
		}
	}
	once.Do(proccess)
	return c
}
