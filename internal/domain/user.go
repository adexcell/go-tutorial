package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int64 `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type UserService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type UserCache interface {
	Set(ctx context.Context, user *User, ttl time.Duration) error
	Get(ctx context.Context, userID int64) (*User, error)
}