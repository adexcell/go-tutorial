package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/adexcell/go-tutorial/internal/domain"
	"golang.org/x/crypto/bcrypt"
)



type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password string) error {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return domain.ErrEmailAlreadyRegistered
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("ошибка при проверке email: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return fmt.Errorf("пароль слишком длинный: %w", err)
	} else if err != nil {
		return fmt.Errorf("неожиданная ошибка: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	return s.repo.Create(ctx, user)
}
