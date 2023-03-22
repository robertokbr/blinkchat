package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/infrastructure/pkg/websocket"
	"github.com/robertokbr/blinkchat/src/usecases"
)

func init() {
	godotenv.Load()
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	logger.Debug("Starting app with debug mode on...")

	connection, err := database.NewDatabase().Connect()

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	usersRepo := repositories.NewUsersRepository(connection)

	createUserUC := usecases.NewCreateUser(usersRepo)

	var pool = websocket.NewPool()

	var websocketConnectionsController = controllers.WebsocketConnections{
		Pool:              pool,
		CreateUserUsecase: createUserUC,
	}

	threads := runtime.NumCPU()

	go func() {
		for i := 0; i < threads; i++ {
			go pool.Start(i)
		}

		pool.MatchPairs()
	}()

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/ws", websocketConnectionsController.Create)
	http.HandleFunc("/ws/connections", websocketConnectionsController.FindAll)

	logger.Info("Starting server on port 8080...")

	if err = http.ListenAndServe(":8080", nil); err != nil {
		logger.Infof("error starting server: %v", err)
	}
}
