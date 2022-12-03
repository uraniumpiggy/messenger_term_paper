package messagesdb

import (
	"context"
	"database/sql"
	"messenger/internal/messages"
)

type db struct {
	*sql.DB
}

func NewStorage(database *sql.DB) *db {
	return &db{database}
}

func (d *db) IsUserInChat(ctx context.Context, userId, chatId uint32) (bool, error) {
	var count int
	err := d.QueryRowContext(ctx, `select count(user_id) from users_chats where user_id = $1 and chat_id = $2`, userId, chatId).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (d *db) SaveMessage(ctx context.Context, data *messages.Message) error {
	_, err := d.ExecContext(ctx, `insert into messages (user_id, chat_id, body) values ($1, $2, $3)`, data.UserId, data.ChatId, data.Body)
	return err
}

func (d *db) GetMessages(ctx context.Context, limit, offset, chat_id uint32) ([]*messages.Message, error) {
	res := make([]*messages.Message, 0)
	rows, err := d.QueryContext(ctx, `select user_id, chat_id, body, created_at from messages where chat_id = $1 order by created_at desc limit $2 offset $3`, chat_id, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		message := &messages.Message{}
		if err := rows.Scan(&message.UserId, &message.ChatId, &message.Body, &message.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, message)
	}

	return res, nil
}
