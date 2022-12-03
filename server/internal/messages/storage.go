package messages

import "context"

type Storage interface {
	IsUserInChat(ctx context.Context, userId, chatId uint32) (bool, error)
	SaveMessage(context.Context, *Message) error
	GetMessages(context.Context, uint32, uint32, uint32) ([]*Message, error)
}
