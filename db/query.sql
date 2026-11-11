-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id, name, email, created_at, deleted_at;

-- name: GetUserByID :one
SELECT id, name, email, created_at, deleted_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListActiveUsers :many
SELECT id, name, email, created_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;