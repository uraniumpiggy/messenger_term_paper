package messages

import "context"

type Storage interface {
	SaveMessage(context.Context, *Message) error
	GetMessages(context.Context, uint32, uint32, uint32) ([]*Message, error)
}
