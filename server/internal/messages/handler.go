package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"messenger/internal/handlers"
	"messenger/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) handlers.Handler {
	return &handler{service: service}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/ws/{chatId}", middleware.ErrorMiddleware(middleware.AuthMiddleware(middleware.WsMiddleware(h.ServeChat)))).Methods("GET")
	router.HandleFunc("/chats/{chatId}", middleware.ErrorMiddleware(middleware.AuthMiddleware(h.GetChatHistory))).Methods("GET")
}

func (h *handler) ServeChat(conn *websocket.Conn, chatId, userId uint32) error {
	err := h.service.SendMessageToChat(context.Background(), conn, chatId, userId)
	return err
}

func (h *handler) GetChatHistory(w http.ResponseWriter, r *http.Request, userId uint32) error {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	params := mux.Vars(r)
	chatId := params["chatId"]
	if page == "" || limit == "" || chatId == "" {
		return fmt.Errorf("Bad request")
	}
	msgs, err := h.service.GetMessages(context.Background(), page, limit, chatId)
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(msgs)
	return nil
}
