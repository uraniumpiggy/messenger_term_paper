package middleware

// import (
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// type AppWsHandlerFunc func (websocket.Conn)

// func WebSocketMiddleware(h AppWsHandlerFunc) http.HandlerFunc {
// 	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ws, err
// }
