-- name: AddCartItem :one
INSERT INTO "CartItems" (user_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCartItemsByUser :many
SELECT * FROM "CartItems" WHERE user_id = $1;

-- name: DeleteCartItem :exec
DELETE FROM "CartItems" WHERE cart_item_id = $1;

-- name: ClearUserCart :exec
DELETE FROM "CartItems" WHERE user_id = $1;

-- name: GetCartItemByUserAndProduct :one
SELECT * FROM "CartItems"
WHERE user_id = $1 AND product_id = $2
LIMIT 1;

-- name: UpdateCartItemQuantity :one
UPDATE "CartItems"
SET quantity = $1
WHERE user_id = $2 AND product_id = $3
RETURNING cart_item_id, user_id, product_id, quantity, price;

