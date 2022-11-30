package middleware

import (
	"fmt"
	"net/http"
)

type appErrorHandler func(w http.ResponseWriter, r *http.Request) error

type Middleware struct {
}

func (m *Middleware) AuthMiddleware(h appErrorHandler) appErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("Auth...")
		h(w, r)
		return nil
	}
}

func (m *Middleware) ErrorMiddleware(h appErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Error handling...")
		h(w, r)
	}
}

func (m *Middleware) WsMiddleware() {

}
