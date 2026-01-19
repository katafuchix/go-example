-- name: GetUser :one
SELECT * FROM users 
WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
WHERE deleted_at IS NULL 
ORDER BY id DESC;

-- name: CreateUser :execresult
INSERT INTO users (name) VALUES (?);

-- name: UpdateUser :exec
UPDATE users SET name = ? 
WHERE id = ? AND deleted_at IS NULL;

-- name: DeleteUser :exec
-- 物理削除ではなく、現在時刻を入れて「論理削除」にする
UPDATE users SET deleted_at = NOW(3) 
WHERE id = ?;