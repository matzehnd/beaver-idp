package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	eventStore EventStore
	users      map[string]*User
	admins     map[string]bool
}

func NewUserService(eventStore EventStore) *UserService {
	return &UserService{
		eventStore: eventStore,
		users:      make(map[string]*User),
		admins:     make(map[string]bool),
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
	if len(s.admins) == 0 {
		adminEvent := IsAdminEvent{Email: user.Email}
		if err := s.eventStore.Save(adminEvent); err != nil {
			return err
		}
		s.apply(adminEvent)
	}
	return nil
}

func (s *UserService) GetUser(email string) (*User, error) {
	user, exists := s.users[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) PasswordIsValid(user User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) UserIsAdmin(user User) bool {
	isAdmin, exists := s.admins[user.Email]
	return exists && isAdmin
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
		s.admins[e.Email] = false
	case IsAdminEvent:
		s.admins[e.Email] = true
	}
}
