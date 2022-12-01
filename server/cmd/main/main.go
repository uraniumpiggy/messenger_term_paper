package main

import (
	"context"
	"log"
	"messenger/internal/messages"
	messagesdb "messenger/internal/messages/message_db"

	// "messenger/internal/messages/messagesdb"
	"messenger/internal/users"
	"messenger/internal/users/db"
	"messenger/pkg/client/postgres"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	dbClient, err := postgres.NewClient(
		context.Background(),
		"localhost",
		"5432",
		"user",
		"secret",
		"service-db")
	if err != nil {
		panic(err)
	}

	log.Println("Connected to database")

	userStorage := db.NewStorage(dbClient)
	userService := users.NewService(userStorage)
	userHandler := users.NewHandler(userService)
	userHandler.Register(router)

	chatStorage := messagesdb.NewStorage(dbClient)
	chatService := messages.NewService(chatStorage)
	chatHandler := messages.NewHandler(chatService)
	chatHandler.Register(router)

	http.ListenAndServe(":8080", router)
}
