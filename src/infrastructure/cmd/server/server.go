package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/infrastructure/controllers/factories"
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
	_, err := database.NewDatabase().Connect()

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	wsConnections := factories.MakeWebsocketConnectionsController()

	threads := runtime.NumCPU()

	go func() {
		for i := 0; i < threads; i++ {
			go wsConnections.Pool.Start(i)
		}

		wsConnections.Pool.MatchPairs()
	}()

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/ws", wsConnections.Create)
	http.HandleFunc("/ws/connections", wsConnections.FindAll)

	logger.Info("Starting server on port 8080...")

	if err = http.ListenAndServe(":8080", nil); err != nil {
		logger.Infof("error starting server: %v", err)
	}
}
