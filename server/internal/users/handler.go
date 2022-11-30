package users

import (
	"context"
	"encoding/json"
	"messenger/internal/handlers"
	"messenger/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
	service *Service
	mid     *middleware.Middleware
}

func NewHandler(service *Service, mid *middleware.Middleware) handlers.Handler {
	return &handler{
		service: service,
		mid:     mid,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/register", h.mid.ErrorMiddleware(h.RegisterUser)).Methods("POST")
	router.HandleFunc("/login", h.mid.ErrorMiddleware(h.AuthUser)).Methods("GET")
	router.HandleFunc("/chats/create", h.mid.ErrorMiddleware(h.mid.AuthMiddleware(h.CreateChat))).Methods("POST")
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	var credentials UserRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return err
	}
	err = h.service.RegisterUser(context.Background(), &credentials)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) AuthUser(w http.ResponseWriter, r *http.Request) error {
	var credentials UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return err
	}
	s, err := h.service.AuthUser(context.Background(), &credentials)
	if err != nil {
		return err
	}
	w.Write([]byte(s))
	return nil
}

func (h *handler) CreateChat(w http.ResponseWriter, r *http.Request) error {
	return nil
}
