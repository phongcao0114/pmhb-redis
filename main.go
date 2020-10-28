package main

import (
	"flag"
	"pmhb-redis/internal/app"
	"pmhb-redis/internal/app/config"

	"log"
	"os"
	"strings"
)

func main() {
	// 1. Parse flag
	configPath := flag.String("config", "configs", "set configs path, default as: 'configs'")
	state := flag.String("state", "dev", "set working environment")
	port := flag.Int("port", 8082, "port number")
	flag.Parse()

	// 2. Allow override state of the app via environment variable
	appState := os.Getenv("APP_STATE")
	if len(strings.TrimSpace(appState)) > 0 {
		*state = appState
	}

	// 3. Get new configuration
	conf, err := config.New(*configPath, *state, *port)
	if err != nil {
		log.Fatal(err)
	}

	server := app.New()
	server.Use(
		app.SetConfig(conf),
		app.UseRedis(),
		app.UseRouter(),
	)

	server.Start()
}
