package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type appErrorHandler func(w http.ResponseWriter, r *http.Request) error

func AuthMiddleware(h appErrorHandler) appErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Header["Token"] == nil {
			return fmt.Errorf("Unauth")
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return []byte("SecretYouShouldHide"), nil
		})

		if err != nil {
			return err
		}

		if token.Valid {
			h(w, r)
		} else {
			return fmt.Errorf("unauth")
		}

		return nil
	}
}

func ErrorMiddleware(h appErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}
}

func WsMiddleware() {

}
