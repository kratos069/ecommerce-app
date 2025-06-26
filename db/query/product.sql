-- name: CreateProduct :one
INSERT INTO "Products" 
(name, description, price, stock_quantity, category_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM "Products" WHERE product_id = $1;

-- name: ListProducts :many
SELECT * FROM "Products" 
ORDER BY product_id 
LIMIT $1 OFFSET $2;

-- name: ListProductsByCategory :many
SELECT * FROM "Products" 
WHERE category_id = $1 
ORDER BY product_id;

-- name: GetProductForUpdate :one
SELECT * FROM "Products"
WHERE product_id = $1
FOR UPDATE;

-- name: DecreaseProductStock :exec
UPDATE "Products"
SET stock_quantity = stock_quantity - $1
WHERE product_id = $2 AND stock_quantity >= $1;
