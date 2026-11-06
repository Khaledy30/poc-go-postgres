package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db *sql.DB
}

func NewDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL não está definido")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	query := `
        INSERT INTO users (name, email)
        VALUES ($1, $2)
        RETURNING id, name, email, created_at
    `

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, input.Name, input.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) List(ctx context.Context, page, limit int) ([]User, error) {
	offset := (page - 1) * limit
	query := `
        SELECT id, name, email, created_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
        SELECT id, name, email, created_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

var (
	ErrNotFound       = errors.New("usuário não encontrado")
	ErrAlreadyDeleted = errors.New("usuário já foi deletado")
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}

	// Nenhuma linha afetada: não existe ou já estava deletado
	var deletedAt sql.NullTime
	err = r.db.QueryRowContext(ctx,
		`SELECT deleted_at FROM users WHERE id = $1`, id,
	).Scan(&deletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	return ErrAlreadyDeleted
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
