package domain

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, user User) (User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

var ErrEmailDuplicate = errors.New("Email sudah terdaftar. Silakan gunakan email lain.")
var ErrUserNotFound = errors.New("Pengguna tidak ditemukan.")
var ErrInvalidCredentials = errors.New("email atau password anda salah.")
