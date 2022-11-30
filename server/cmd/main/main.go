package main

import (
	"context"
	"log"
	"messenger/internal/messeges"
	"messenger/internal/middleware"
	"messenger/internal/users"
	"messenger/internal/users/db"
	"messenger/pkg/client/postgres"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	mid := &middleware.Middleware{}

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

	userHandler := users.NewHandler(userService, mid)

	userHandler.Register(router)

	chatHandler := messeges.NewHandler()
	chatHandler.Register(router)

	http.ListenAndServe(":8080", router)
}
