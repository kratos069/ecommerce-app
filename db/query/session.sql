-- name: CreateSession :one
INSERT INTO "Sessions" (
  id, username, refresh_token, user_agent, client_ip, is_blocked, expired_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM "Sessions"
WHERE id = $1 LIMIT 1;