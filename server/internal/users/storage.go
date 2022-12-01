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
}
