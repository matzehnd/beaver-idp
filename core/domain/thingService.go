package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type ThingService struct {
	eventStore EventStore
	things     map[string]*Thing
}

func NewThingService(eventStore EventStore) *ThingService {
	return &ThingService{
		eventStore: eventStore,
		things:     make(map[string]*Thing),
	}
}

func (s *ThingService) RegisterThing(thing RegisterThing) error {
	_, exists := s.things[thing.Id]
	if exists {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(thing.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash password: %T", err)
	}
	event := ThingRegisteredEvent{Id: thing.Id, Password: string(hash)}
	if err := s.eventStore.Save(event); err != nil {
		return err
	}
	s.apply(event)

	return nil
}

func (s *ThingService) GetThing(id string) (*Thing, error) {
	thing, exists := s.things[id]
	if !exists {
		return nil, fmt.Errorf("thing not found")
	}
	return thing, nil
}

func (s *ThingService) PasswordIsValid(thing Thing, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(thing.Password), []byte(password))
	return err == nil
}

func (s *ThingService) RebuildEventState() error {
	events, err := s.eventStore.Load()
	if err != nil {
		return err
	}
	for _, event := range events {
		s.apply(event)
	}
	return nil
}

func (s *ThingService) apply(event interface{}) {
	switch e := event.(type) {
	case ThingRegisteredEvent:
		s.things[e.Id] = &Thing{
			Id:       e.Id,
			Password: e.Password,
		}
	}
}
