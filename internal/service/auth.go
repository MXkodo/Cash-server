package service

import (
	"context"
	"errors"
	"time"

	"github.com/MXkodo/cash-server/internal/repo"
	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repo.UserRepo
	jwtSecret []byte
	rdb       *redis.Client
}

func NewAuthService(userRepo *repo.UserRepo, jwtSecret string, rdb *redis.Client) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		rdb:       rdb,
	}
}

func (s *AuthService) Register(login, pswd string) (map[string]string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.CreateUser(login, string(hashedPassword)); err != nil {
		return nil, err
	}
	return map[string]string{"login": login}, nil
}

func (s *AuthService) Authenticate(login, pswd string) (string, error) {
	user, err := s.userRepo.FindUserByLogin(login)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pswd)) != nil {
		return "", errors.New("неверный логин или пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID, 
		"login":    user.Login,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), 
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.rdb.Set(ctx, tokenString, user.Login, 72*time.Hour).Err(); err != nil {
		return "", err
	}

	return tokenString, nil
}


func (s *AuthService) Logout(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.rdb.Del(ctx, token).Err(); err != nil {
		return err 
	}

	return nil
}
