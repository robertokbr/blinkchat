package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	controller_factories "github.com/robertokbr/blinkchat/src/infrastructure/controllers/factories"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
)

func init() {
	godotenv.Load()

	logger.Debug("Starting app with debug mode on...")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	if _, err := database.NewDatabase().Connect(); err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	wsConnections := controller_factories.MakeWebsocketConnectionsController()

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/ws", wsConnections.Create)
	http.HandleFunc("/ws/connections", wsConnections.FindAll)

	logger.Info("Starting server on port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Infof("error starting server: %v", err)
	}
}
