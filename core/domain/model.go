package domain

type RegisterUser struct {
	Email    string
	Password string
}

type User struct {
	Email    string
	Password string
}

type UserRegisteredEvent struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
