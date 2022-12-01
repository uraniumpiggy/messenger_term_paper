package messages

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Service struct {
	mutex       sync.RWMutex
	connections map[uint32][]*websocket.Conn
	storage     Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		mutex:       sync.RWMutex{},
		connections: make(map[uint32][]*websocket.Conn),
		storage:     storage,
	}
}

func (s *Service) SendMessageToChat(ctx context.Context, conn *websocket.Conn, chatId, userId uint32) error {
	disconnect := make(chan struct{})
	s.mutex.Lock()
	s.connections[chatId] = append(s.connections[chatId], conn)
	s.mutex.Unlock()
	fmt.Println(s.connections)
	defer func() {
		s.mutex.Lock()
		for idx, val := range s.connections[chatId] {
			if val == conn {
				s.connections[chatId] = append(s.connections[chatId][:idx], s.connections[chatId][idx+1:]...)
				break
			}
		}
		if len(s.connections[chatId]) == 0 {
			delete(s.connections, chatId)
		}
		fmt.Println(s.connections)
		s.mutex.Unlock()
	}()
	go func() {
		for {
			t, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("err 1: %s", err)
				disconnect <- struct{}{}
				break
			}
			for _, c := range s.connections[chatId] {
				if c != conn {
					err = c.WriteMessage(t, message)
					if err != nil {
						fmt.Printf("err 2: %s", err)
						break
					}
				}
			}
			msg := &Message{
				ChatId: chatId,
				UserId: userId,
				Body:   string(message),
			}
			s.storage.SaveMessage(ctx, msg)
		}
	}()

	select {
	case <-disconnect:
		return nil
	}

}

func (s *Service) GetMessages(ctx context.Context, page, limit, chatId string) ([]*Message, error) {
	nPage, err1 := strconv.ParseUint(page, 10, 32)
	nLimit, err2 := strconv.ParseUint(limit, 10, 32)
	nChatId, err3 := strconv.ParseUint(chatId, 10, 32)

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("Error with string convertion %s, %s, %s", err1, err2, err3)
	}

	if nPage < 0 || nLimit < 0 || nChatId < 0 {
		return nil, fmt.Errorf("Negative number(s)")
	}

	offset := (nPage - 1) * nLimit
	messages, err := s.storage.GetMessages(ctx, uint32(nLimit), uint32(offset), uint32(nChatId))
	if err != nil {
		return nil, err
	}
	return messages, nil
}
