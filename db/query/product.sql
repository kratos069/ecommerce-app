-- name: CreateProduct :one
INSERT INTO "Products" 
(name, description, product_image, price, stock_quantity, category_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM "Products" WHERE product_id = $1;

-- name: ListProducts :many
SELECT * FROM "Products" 
ORDER BY product_id 
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE "Products"
SET name = $2, description = $3, product_image = $4,
stock_quantity = $5, price = $6, category_id = $7
WHERE product_id = $1
RETURNING *;

-- name: ListProductsByCategory :many
SELECT * FROM "Products" 
WHERE category_id = $1 
ORDER BY product_id;

-- name: GetProductForUpdate :one
SELECT * FROM "Products"
WHERE product_id = $1
FOR NO KEY UPDATE;

-- name: DeleteProduct :exec
DELETE FROM "Products"
WHERE product_id = $1;

-- name: DecreaseProductStock :exec
UPDATE "Products"
SET stock_quantity = stock_quantity - $1
WHERE product_id = $2 AND stock_quantity >= $1;
