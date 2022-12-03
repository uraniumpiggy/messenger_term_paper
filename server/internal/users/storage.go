package users

import (
	"context"
)

type Storage interface {
	RegisterUser(context.Context, *UserRegisterRequest) error
	AuthUser(context.Context, *UserLoginRequest) (string, error)
	GetUserInfo(context.Context, *UserLoginRequest) (*UserInfo, error)
	CreateChat(context.Context, *CreateChatRequest, uint32) error
	GetAllUsernames(ctx context.Context, prefix string) ([]string, error)
	GetUserChats(ctx context.Context, userId uint32) ([]*ChatInfo, error)
	DeleteChat(ctx context.Context, chatId uint32) error
	AddUserToChat(ctx context.Context, username string, chatId uint32) error
	RemoveUserFromChat(ctx context.Context, username string, chatId uint32) error
	IsUserInChat(ctx context.Context, userId uint32, chatId uint32) (bool, error)
}
