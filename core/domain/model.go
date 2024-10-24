package domain

type EventStore interface {
	Save(event interface{}) error
	Load() ([]interface{}, error)
}

type RegisterUser struct {
	Email    string
	Password string
}
type RegisterThing struct {
	Id       string
	Password string
}

type User struct {
	Email    string
	Password string
}

type Thing struct {
	Id       string
	Password string
}

type ThingRegisteredEvent struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
type UserRegisteredEvent struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type IsAdminEvent struct {
	Email string `json:"email"`
}
