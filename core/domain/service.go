package domain

import "fmt"

type EventStore interface {
	Save(event interface{}) error
	Load() ([]interface{}, error)
}

type UserService struct {
	eventStore EventStore
	users      map[string]*User
}

func NewUserService(eventStore EventStore) *UserService {
	return &UserService{
		eventStore: eventStore,
		users:      make(map[string]*User),
	}
}

func (s *UserService) RegisterUser(user RegisterUser) error {
	event := UserRegisteredEvent(user)
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
