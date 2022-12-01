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
}

func NewHandler(service *Service) handlers.Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/register", middleware.ErrorMiddleware(h.RegisterUser)).Methods("POST")
	router.HandleFunc("/login", middleware.ErrorMiddleware(h.AuthUser)).Methods("GET")
	router.HandleFunc("/chats/create", middleware.ErrorMiddleware(middleware.AuthMiddleware(h.CreateChat))).Methods("POST")
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

	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) CreateChat(w http.ResponseWriter, r *http.Request, userId uint32) error {
	var data CreateChatRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	err = h.service.CreateChat(context.Background(), &data, userId)
	return nil
}
