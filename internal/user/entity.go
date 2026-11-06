package user

import "time"

type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // soft delete
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
