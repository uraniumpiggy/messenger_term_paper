package middleware

import (
	"errors"
	"fmt"
	"messenger/internal/apperror"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type appErrorHandler func(w http.ResponseWriter, r *http.Request) error
type authHandler func(w http.ResponseWriter, r *http.Request, userId uint32) error
type appWsHandler func(*websocket.Conn, uint32, uint32) error

func AuthMiddleware(h authHandler) appErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Header["Token"] == nil {
			return apperror.ErrUnauthorized
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return []byte("SecretYouShouldHide"), nil
		})

		if err != nil {
			return apperror.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			return apperror.ErrUnauthorized
		}

		userId := claims["user_id"].(float64)

		if token.Valid {
			h(w, r, uint32(userId))
		} else {
			return apperror.ErrUnauthorized
		}

		return nil
	}
}

func ErrorMiddleware(h appErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var appError *apperror.AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appError) {
				if errors.Is(err, apperror.ErrNotFound) {
					w.WriteHeader(404)
					w.Write(apperror.ErrNotFound.Marshal())
					return
				}
				if errors.Is(err, apperror.ErrBadRequest) {
					w.WriteHeader(400)
					w.Write(apperror.ErrBadRequest.Marshal())
					return
				}
				if errors.Is(err, apperror.ErrUnauthorized) {
					w.WriteHeader(401)
					w.Write(apperror.ErrBadRequest.Marshal())
					return
				}
				if errors.Is(err, apperror.ErrInternalError) {
					w.WriteHeader(500)
					w.Write(apperror.ErrBadRequest.Marshal())
					return
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsMiddleware(h appWsHandler) authHandler {
	return func(w http.ResponseWriter, r *http.Request, uiserId uint32) error {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return err
		}
		params := mux.Vars(r)
		id := params["chatId"]
		if id == "" {
			return fmt.Errorf("Some err")
		}
		chatId, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		err = h(conn, uint32(chatId), uiserId)
		return err
	}
}
