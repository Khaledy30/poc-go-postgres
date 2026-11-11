package user

import (
	"context"
	"errors"
	"os"
	db "poc-go/internal/database"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewDB() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func toUser(user db.User) *User {
	var deletedAt *time.Time

	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return &User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		DeletedAt: deletedAt,
	}
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, queries: db.New(pool)}
}

func (r *Repository) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

func (r *Repository) List(ctx context.Context, page, limit int) ([]User, error) {
	offset := (page - 1) * limit

	users, err := r.queries.ListActiveUsers(ctx, db.ListActiveUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]User, 0, len(users))

	for _, user := range users {
		result = append(result, *toUser(user))
	}

	return result, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

var (
	ErrNotFound       = errors.New("usuário não encontrado")
	ErrAlreadyDeleted = errors.New("usuário já foi deletado")
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	user, err := r.queries.GetUserByID(ctx, parsedID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if user.DeletedAt.Valid {
		return ErrAlreadyDeleted
	}

	err = r.queries.SoftDeleteUser(ctx, parsedID)
	if err != nil {
		return err
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
