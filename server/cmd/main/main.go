package main

import (
	"context"
	"fmt"
	"messenger/internal/messages"
	messagesdb "messenger/internal/messages/message_db"

	"messenger/internal/config"
	"messenger/internal/users"
	"messenger/internal/users/db"
	"messenger/pkg/client/postgres"
	"messenger/pkg/logging"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	logger := logging.NewLogger()

	logger.Info("Create router")
	router := mux.NewRouter()

	cfg := config.GetConfig()

	dbClient, err := postgres.NewClient(
		context.Background(),
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Database)
	if err != nil {
		panic(err)
	}

	logger.Info("Connected to database")

	userStorage := db.NewStorage(dbClient)
	userService := users.NewService(userStorage)
	userHandler := users.NewHandler(userService)
	userHandler.Register(router)

	logger.Info("User handler ready")

	chatStorage := messagesdb.NewStorage(dbClient)
	chatService := messages.NewService(chatStorage)
	chatHandler := messages.NewHandler(chatService)
	chatHandler.Register(router)

	logger.Info("Message handler ready")

	logger.Info("Server started...")
	err2 := http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port), router)
	if err2 != nil {
		panic(err2)
	}
}
