-- name: CreateUser :one
INSERT INTO "Users" (username, password_hash, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM "Users" WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "Users" WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM "Users" ORDER BY user_id LIMIT $1 OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM "Users"
WHERE user_id = $1;
