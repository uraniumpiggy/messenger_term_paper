package users

import (
	"context"
	"fmt"
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

func (s *Service) AuthUser(ctx context.Context, data *UserLoginRequest) (string, error) {
	passHash, err := s.storage.AuthUser(ctx, data)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(data.Password))
	if err != nil {
		return "", fmt.Errorf("Incorrect credentials")
	}

	token, err := s.generateJWT("username")
	if err != nil {
		return "", err
	}
	return token, nil

}

func (s *Service) generateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = username

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
