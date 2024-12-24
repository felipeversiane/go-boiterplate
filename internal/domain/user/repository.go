package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/felipeversiane/go-boiterplate/internal/domain"
	"github.com/felipeversiane/go-boiterplate/internal/infra/database"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db database.DatabaseInterface
}

type UserRepositoryInterface interface {
	Create(user domain.User, ctx context.Context) (string, error)
	BulkCreate(users []domain.User, ctx context.Context) error
	Retrieve(id string, ctx context.Context) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Update(id string, user domain.User, ctx context.Context) error
	Delete(id string, ctx context.Context) error
}

func NewUserRepository(db database.DatabaseInterface) UserRepositoryInterface {
	return &userRepository{db}
}

func (repository *userRepository) Create(domain domain.User, ctx context.Context) (string, error) {
	query := `INSERT INTO users (id, email, first_name, last_name, password, created_at, updated_at, deleted) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	args := []interface{}{
		domain.ID,
		domain.Email,
		domain.FirstName,
		domain.LastName,
		domain.Password,
		domain.CreatedAt,
		domain.UpdatedAt,
		domain.Deleted,
	}
	var id string
	err := repository.db.GetDB().QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("unable to insert user: %w", err)
	}
	return id, nil
}

func (repository *userRepository) BulkCreate(users []domain.User, ctx context.Context) error {
	query := `
		INSERT INTO users (id, email, first_name, last_name, password, created_at, updated_at, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	batch := &pgx.Batch{}
	for _, user := range users {
		batch.Queue(query,
			user.ID,
			user.Email,
			user.FirstName,
			user.LastName,
			user.Password,
			user.CreatedAt,
			user.UpdatedAt,
			user.Deleted,
		)
	}

	results := repository.db.GetDB().SendBatch(ctx, batch)
	defer results.Close()

	for _, user := range users {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				slog.Info("user %s already exists", user.Email)
				continue
			}

			return fmt.Errorf("unable to insert row: %w", err)
		}
	}

	return results.Close()
}

func (repository *userRepository) Retrieve(id string, ctx context.Context) (*domain.User, error) {
	query := `
		SELECT id, email, first_name, last_name, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted = false`

	var user domain.User
	err := repository.db.GetDB().QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error querying user by ID: %w", err)
	}

	return &user, nil
}

func (repository *userRepository) List(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, email, first_name, last_name, created_at, updated_at
		FROM users
		WHERE deleted = false`

	rows, err := repository.db.GetDB().Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		user := domain.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *userRepository) Update(id string, domain domain.User, ctx context.Context) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, updated_at = $3 
		WHERE id = $4 AND deleted = false
	`

	args := []interface{}{
		domain.FirstName,
		domain.LastName,
		domain.UpdatedAt,
		id,
	}

	_, err := repository.db.GetDB().Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (repository *userRepository) Delete(id string, ctx context.Context) error {
	query := `
		UPDATE users 
		SET deleted = true, updated_at = $1 
		WHERE id = $2 AND deleted = false
	`

	args := []interface{}{
		time.Now().UTC(),
		id,
	}

	result, err := repository.db.GetDB().Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}
