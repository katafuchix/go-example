-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (name, created_at, updated_at) 
VALUES (?, NOW(), NOW());

-- name: UpdateUser :exec
UPDATE users SET name = ?, updated_at = NOW() WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id DESC;