package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adexcell/go-tutorial/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         domain.UserRepository
	tokenManager domain.TokenManager
	cache        domain.UserCache
	tokenTTL     time.Duration
	cacheTTL     time.Duration
}


func NewUserService(
	repo domain.UserRepository,
	tokenManager domain.TokenManager,
	cache domain.UserCache,
	tokenTTL time.Duration,
	cacheTTL time.Duration,
) *UserService {
	return &UserService{
		repo:         repo,
		tokenManager: tokenManager,
		cache:        cache,
		tokenTTL:     tokenTTL,
		cacheTTL:     cacheTTL,
	}
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

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	if errors.Is(err, domain.ErrNotFound) {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.tokenManager.NewJWT(user.ID, s.tokenTTL)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}
	user, err := s.cache.Get(ctx, id)

	// Cache hit
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("ошибка при обращении к кешу: %w", err)
	}

	user, err = s.repo.GetByID(ctx, id)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка при обращении к бд: %w", err)
	}

	if err := s.cache.Set(ctx, user, s.cacheTTL); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении данных в кеш: %w", err)
	}

	return user, nil
}
