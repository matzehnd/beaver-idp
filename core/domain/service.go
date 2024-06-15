package domain

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type EventStore interface {
	Save(event interface{}) error
	Load() ([]interface{}, error)
}

type UserService struct {
	eventStore EventStore
	users      map[string]*User
	privateKey []byte
}

func NewUserService(eventStore EventStore, privateKey []byte) *UserService {
	return &UserService{
		eventStore: eventStore,
		users:      make(map[string]*User),
		privateKey: privateKey,
	}
}

func (s *UserService) RegisterUser(user RegisterUser) error {
	_, exists := s.users[user.Email]
	if exists {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash password: %T", err)
	}
	event := UserRegisteredEvent{Email: user.Email, Password: string(hash)}
	if err := s.eventStore.Save(event); err != nil {
		return err
	}
	s.apply(event)
	return nil
}

func (s *UserService) GetUser(email string) (*User, error) {
	user, exists := s.users[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) CreateToken(tokenRequest CreateToken) (string, error) {
	user, exists := s.users[tokenRequest.Email]
	if !exists {
		return "", fmt.Errorf("user not found")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tokenRequest.Password))

	if err != nil {
		return "", fmt.Errorf("pw wrong")
	}

	token, err := tokenFromUser(*user, s.privateKey)

	if err != nil {
		return "", fmt.Errorf("unable to get token: %T", err)
	}

	return token, nil
}

func (s *UserService) RebuildEventStream() error {
	events, err := s.eventStore.Load()
	if err != nil {
		return err
	}
	for _, event := range events {
		s.apply(event)
	}
	return nil
}

func (s *UserService) apply(event interface{}) {
	switch e := event.(type) {
	case UserRegisteredEvent:
		s.users[e.Email] = &User{
			Email:    e.Email,
			Password: e.Password,
		}
	}
}

func tokenFromUser(user User, privateKey []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	fmt.Println(err)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"sub": user.Email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
