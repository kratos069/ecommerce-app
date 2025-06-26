-- name: CreateReview :one
INSERT INTO "Reviews" (user_id, product_id, rating, comment)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetReviewsForProduct :many
SELECT * FROM "Reviews" 
WHERE product_id = $1 ORDER BY created_at DESC;
