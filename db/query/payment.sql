-- name: CreatePayment :one
INSERT INTO "Payments" (order_id, payment_method, amount, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPaymentByOrder :one
SELECT * FROM "Payments" WHERE order_id = $1;
