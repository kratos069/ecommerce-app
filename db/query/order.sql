-- name: CreateOrder :one
INSERT INTO "Orders" (user_id, total_amount, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM "Orders" WHERE order_id = $1;

-- name: ListOrdersByUser :many
SELECT * FROM "Orders" WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CountUserProductOrders :one
SELECT COUNT(*) FROM "Orders" o
JOIN "OrderItems" oi ON o.order_id = oi.order_id
WHERE o.user_id = $1 AND oi.product_id = $2;

-- name: UpdateOrderStatus :one
UPDATE "Orders"
SET status = $2
WHERE order_id = $1
RETURNING *;
