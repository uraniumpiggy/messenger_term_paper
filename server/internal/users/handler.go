package users

import (
	"messenger/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/", h.ServeChat).Methods("GET")
}

func (h *handler) ServeChat(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User hi"))
}
