package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	controller_factories "github.com/robertokbr/blinkchat/src/infrastructure/controllers/factories"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/middlewares"
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

	mux := http.NewServeMux()

	jobs := make(chan models.Message)
	pool := models.NewPool()
	poolWorkerManagerUsecase := usecases.NewPoolWorkerManager(pool, jobs)

	go poolWorkerManagerUsecase.Start()
	defer poolWorkerManagerUsecase.Stop()

	websocketClientsController := controller_factories.MakeWebsocketClientsController(pool, jobs)
	usersController := controller_factories.MakeUsersController()

	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/ws", websocketClientsController.Create)
	mux.HandleFunc("/ws/clients", websocketClientsController.FindAll)
	mux.HandleFunc("/users", usersController.Create)

	logger.Info("Starting server on port 8080...")

	api := middlewares.HttpCors(middlewares.LogRequest(mux))

	if err := http.ListenAndServe(":8080", api); err != nil {
		logger.Infof("error starting server: %v", err)
	}
}
