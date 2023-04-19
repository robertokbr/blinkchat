package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	controller_factories "github.com/robertokbr/blinkchat/src/infrastructure/controllers/factories"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/usecases"
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

	pool := models.NewPool()
	poolWorker := usecases.NewPoolWorker(pool)
	matchPoolPairs := usecases.NewMatchPoolPairs(pool)
	poolManager := usecases.NewPoolManager(poolWorker, matchPoolPairs)
	numOfCPUs := runtime.NumCPU()
	go poolManager.Execute(numOfCPUs)

	wsConnectionsController := controller_factories.MakeWebsocketConnectionsController(pool)

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/ws", wsConnectionsController.Create)
	http.HandleFunc("/ws/connections", wsConnectionsController.FindAll)

	logger.Info("Starting server on port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Infof("error starting server: %v", err)
	}
}
