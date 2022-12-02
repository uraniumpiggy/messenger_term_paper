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
	router.HandleFunc("/chats/{chatId}", middleware.ErrorMiddleware(middleware.AuthMiddleware(h.DeleteChat))).Methods("DELETE")
	router.HandleFunc("/chats/{chatId}/{username}", middleware.ErrorMiddleware(middleware.AuthMiddleware(h.RemoveUser))).Methods("DELETE")
	router.HandleFunc("/chats/{chatId}/{username}", middleware.ErrorMiddleware(middleware.AuthMiddleware(h.AddUser))).Methods("POST")
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
	w.WriteHeader(201)
	return nil
}

func (h *handler) AuthUser(w http.ResponseWriter, r *http.Request) error {
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
	w.WriteHeader(201)
	return err
}

func (h *handler) GetAllUsernames(w http.ResponseWriter, r *http.Request, userId uint32) error {
	prefix := r.URL.Query().Get("prefix")
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

func (h *handler) DeleteChat(w http.ResponseWriter, r *http.Request, userId uint32) error {
	params := mux.Vars(r)
	chatId := params["chatId"]
	if chatId == "" {
		return apperror.ErrBadRequest
	}
	err := h.service.DeleteChat(context.Background(), chatId)
	if err != nil {
		return err
	}
	w.WriteHeader(204)
	return nil
}

func (h *handler) AddUser(w http.ResponseWriter, r *http.Request, userId uint32) error {
	params := mux.Vars(r)
	chatId := params["chatId"]
	username := params["username"]
	if username == "" || chatId == "" {
		return apperror.ErrBadRequest
	}
	err := h.service.AddUser(context.Background(), chatId, username)
	if err != nil {
		return err
	}
	w.WriteHeader(200)
	return nil

}

func (h *handler) RemoveUser(w http.ResponseWriter, r *http.Request, userId uint32) error {
	params := mux.Vars(r)
	chatId := params["chatId"]
	username := params["username"]
	if username == "" || chatId == "" {
		return apperror.ErrBadRequest
	}
	err := h.service.RemoveUser(context.Background(), chatId, username)
	if err != nil {
		return err
	}
	w.WriteHeader(204)
	return nil

}
