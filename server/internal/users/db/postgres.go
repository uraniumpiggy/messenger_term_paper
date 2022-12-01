package db

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/internal/users"
)

type db struct {
	*sql.DB
}

func NewStorage(database *sql.DB) users.Storage {
	return &db{database}
}

func (d *db) RegisterUser(ctx context.Context, data *users.UserRegisterRequest) error {
	var id int
	err := d.QueryRowContext(ctx, `select id from users where login = $1`, data.Login).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if id != 0 {
		return fmt.Errorf("User already exists")
	}
	res, err := d.ExecContext(ctx, `insert into users (login, username, password) values ($1, $2, $3)`, data.Login, data.Username, data.Password)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return fmt.Errorf("This username already exist")
		}
		return err
	}
	affected, err := res.RowsAffected()
	if affected == 0 || err != nil {
		return fmt.Errorf("err in db")
	}
	return nil

}

func (d *db) AuthUser(ctx context.Context, data *users.UserLoginRequest) (string, error) {
	var password string
	err := d.QueryRowContext(ctx, `select password from users where login = $1`, data.Login).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("User not found")
		}
		return "", err
	}
	return password, nil
}

func (d *db) GetUserInfo(ctx context.Context, data *users.UserLoginRequest) (*users.UserInfo, error) {
	res := &users.UserInfo{}
	err := d.QueryRowContext(ctx, `select id, username from users where login = $1`, data.Login).Scan(&res.UserID, &res.Username)
	if err != nil {
		return nil, err
	}

	res.ChatIDs = make([]uint32, 0)
	res.ChatNames = make([]string, 0)
	rows, err := d.QueryContext(ctx, `select id, chat_name from chats as c, users_chats as uc where uc.user_id = $1 and c.id = uc.chat_id`, res.UserID)
	defer rows.Close()

	for rows.Next() {
		var id uint32
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		res.ChatIDs = append(res.ChatIDs, id)
		res.ChatNames = append(res.ChatNames, name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (d *db) CreateChat(ctx context.Context, data *users.CreateChatRequest, userId uint32) error {
	ids := make([]uint32, 0)
	var members bytes.Buffer
	for i, val := range data.ChatMemberNames {
		members.WriteString(fmt.Sprintf("'%s'", val))
		if i != len(data.ChatMemberNames)-1 {
			members.WriteString(",")
		}
	}

	rows, err := d.QueryContext(ctx, fmt.Sprintf("select id from users where username in (%s)", members.String()))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item uint32
		if err := rows.Scan(&item); err != nil {
			return err
		}
		ids = append(ids, item)
	}

	t, err := d.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	var chatId uint32
	err = t.QueryRowContext(ctx, `insert into chats (chat_name) values ($1) returning id`, data.ChatName).Scan(&chatId)
	if err != nil {
		t.Rollback()
		return fmt.Errorf("Error")
	}

	var buffer bytes.Buffer
	buffer.WriteString("insert into users_chats (user_id, chat_id) values ")
	for _, val := range ids {
		buffer.WriteString(fmt.Sprintf("(%d, %d),", val, chatId))
	}
	buffer.WriteString(fmt.Sprintf("(%d, %d)", userId, chatId))

	_, err = t.ExecContext(ctx, buffer.String())
	if err != nil {
		t.Rollback()
		return fmt.Errorf("Error")
	}

	t.Commit()

	return nil
}

func (d *db) GetAllUsernames(ctx context.Context, prefix string) ([]string, error) {
	res := make([]string, 0)
	r, err := d.QueryContext(ctx, `select username from users where username like $1`, "%"+prefix+"%")
	defer r.Close()
	if err != nil {
		return nil, err
	}
	for r.Next() {
		var name string
		if err := r.Scan(&name); err != nil {
			return nil, err
		}
		res = append(res, name)
	}
	return res, nil
}

func (d *db) GetUserChats(ctx context.Context, userId uint32) ([]*users.ChatInfo, error) {
	res := make([]*users.ChatInfo, 0)
	r, err := d.QueryContext(ctx, `select c.id, c.chat_name from chats as c left join users_chats as uc on uc.user_id = $1 and c.id = uc.chat_id`, userId)
	defer r.Close()
	if err != nil {
		return nil, err
	}

	for r.Next() {
		ci := &users.ChatInfo{}
		if err := r.Scan(&ci.ChatId, &ci.ChatName); err != nil {
			return nil, err
		}
		ci.MemeberNames = make([]string, 0)
		rows, err := d.QueryContext(ctx, `select u.username from users_chats as uc, users as u where uc.chat_id = $1 and uc.user_id = u.id`, ci.ChatId)
		defer rows.Close()
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var username string
			if err := rows.Scan(&username); err != nil {
				return nil, err
			}
			ci.MemeberNames = append(ci.MemeberNames, username)
		}
		res = append(res, ci)
	}

	return res, nil
}
