package main

import (
	"log"
	"net/http"

	"github.com/robertokbr/blinkchat/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/infrastructure/database"
	"github.com/robertokbr/blinkchat/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/infrastructure/pkg/websocket"
	"github.com/robertokbr/blinkchat/usecases"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
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

	for i := 0; i < 10; i++ {
		go pool.Start(i)
	}

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/ws", websocketConnectionsController.Create)
	http.HandleFunc("/connections", websocketConnectionsController.FindAll)

	log.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)
}
