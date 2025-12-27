package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/adexcell/go-tutorial/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		insert into users (email, password_hash)
		values ($1, $2)
		returning id, created_at`

	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `select id, email, password_hash, created_at from users where email=$1`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	} 
	if err != nil {
		return nil, fmt.Errorf("не удалось получить данные: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `select id, email, password_hash, created_at from users where id=$1`
	
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	} 
	if err != nil {
		return nil, fmt.Errorf("не удалось получить данные: %w", err)
	}

	return user, nil
}
