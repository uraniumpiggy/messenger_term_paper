package main

import (
	"messenger/internal/messeges"
	"messenger/internal/users"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	chatHandler := messeges.NewHandler()
	userHandler := users.NewHandler()

	chatHandler.Register(router)
	userHandler.Register(router)

	http.ListenAndServe(":8080", router)

}
