package users

import (
	"context"
	"encoding/json"
	"messenger/internal/apperror"
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
	router.HandleFunc("/chats/get", middleware.ErrorMiddleware(middleware.AuthMiddleware(h.GetUserChats))).Methods("GET")
	router.HandleFunc("/users", middleware.ErrorMiddleware(middleware.AuthMiddleware(h.GetAllUsernames))).Methods("GET")
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	var credentials UserRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return apperror.ErrBadRequest
	}
	err = h.service.RegisterUser(context.Background(), &credentials)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) AuthUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var credentials UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return apperror.ErrBadRequest
	}
	s, err := h.service.AuthUser(context.Background(), &credentials)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		return apperror.ErrInternalError
	}
	return nil
}

func (h *handler) CreateChat(w http.ResponseWriter, r *http.Request, userId uint32) error {
	var data CreateChatRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return apperror.ErrBadRequest
	}

	err = h.service.CreateChat(context.Background(), &data, userId)
	return err
}

func (h *handler) GetAllUsernames(w http.ResponseWriter, r *http.Request, userId uint32) error {
	w.Header().Set("Content-Type", "application/json")
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		return apperror.ErrBadRequest
	}
	ur, err := h.service.GetAllUsernames(context.Background(), prefix)
	if err != nil {
		return apperror.ErrInternalError
	}
	err = json.NewEncoder(w).Encode(ur)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) GetUserChats(w http.ResponseWriter, r *http.Request, userId uint32) error {
	w.Header().Set("Content-Type", "application/json")
	chats, err := h.service.GetUserChats(context.Background(), userId)
	if err != nil {
		return apperror.ErrInternalError
	}
	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		return err
	}
	return nil
}
