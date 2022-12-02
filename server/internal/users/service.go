package users

import (
	"context"
	"fmt"
	"messenger/internal/apperror"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func (s *Service) RegisterUser(ctx context.Context, data *UserRegisterRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return err
	}
	data.Password = string(hash)
	err = s.storage.RegisterUser(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) AuthUser(ctx context.Context, data *UserLoginRequest) (*UserWithToken, error) {
	passHash, err := s.storage.AuthUser(ctx, data)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(data.Password))
	if err != nil {
		return nil, fmt.Errorf("Incorrect credentials")
	}

	uwt, err := s.storage.GetUserInfo(ctx, data)
	if err != nil {
		return nil, err
	}

	userWithToken := &UserWithToken{}

	userWithToken.UserID = uwt.UserID
	userWithToken.Username = uwt.Username
	userWithToken.ChatIDs = uwt.ChatIDs
	userWithToken.ChatNames = uwt.ChatNames

	token, err := s.generateJWT(userWithToken.UserID)
	if err != nil {
		return nil, err
	}

	userWithToken.Token = token

	return userWithToken, nil
}

func (s *Service) CreateChat(ctx context.Context, data *CreateChatRequest, userId uint32) error {
	if len(data.ChatMemberNames) == 0 {
		return fmt.Errorf("Chat must contain at least two members")
	}
	err := s.storage.CreateChat(ctx, data, userId)
	return err
}

func (s *Service) GetUserChats(ctx context.Context, userId uint32) ([]*ChatInfo, error) {
	ci, err := s.storage.GetUserChats(ctx, userId)
	if err != nil {
		return nil, err
	}
	return ci, nil
}

func (s *Service) GetAllUsernames(ctx context.Context, prefix string) (*UsernameResponse, error) {
	usernames, err := s.storage.GetAllUsernames(ctx, prefix)
	if err != nil {
		return nil, err
	}
	res := &UsernameResponse{Usernames: usernames}
	return res, nil
}

func (s *Service) generateJWT(userId uint32) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user_id"] = userId

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) DeleteChat(ctx context.Context, chatId string) error {
	chId, err := strconv.ParseUint(chatId, 10, 32)
	if err != nil {
		return apperror.ErrBadRequest
	}
	if chId < 1 {
		return apperror.ErrBadRequest
	}
	err = s.storage.DeleteChat(ctx, uint32(chId))
	return err
}

func (s *Service) AddUser(ctx context.Context, chatId, username string) error {
	chId, err := strconv.ParseUint(chatId, 10, 32)
	if err != nil {
		return apperror.ErrBadRequest
	}
	if chId < 1 {
		return apperror.ErrBadRequest
	}
	err = s.storage.AddUserToChat(ctx, username, uint32(chId))
	return err
}

func (s *Service) RemoveUser(ctx context.Context, chatId, username string) error {
	chId, err := strconv.ParseUint(chatId, 10, 32)
	if err != nil {
		return apperror.ErrBadRequest
	}
	if chId < 1 {
		return apperror.ErrBadRequest
	}
	err = s.storage.RemoveUserFromChat(ctx, username, uint32(chId))
	return err
}
