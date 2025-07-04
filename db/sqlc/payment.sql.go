// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: payment.sql

package db

import (
	"context"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO "Payments" (order_id, payment_method, amount, status)
VALUES ($1, $2, $3, $4)
RETURNING payment_id, order_id, payment_method, amount, status, created_at
`

type CreatePaymentParams struct {
	OrderID       int64   `json:"order_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRow(ctx, createPayment,
		arg.OrderID,
		arg.PaymentMethod,
		arg.Amount,
		arg.Status,
	)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.OrderID,
		&i.PaymentMethod,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getPaymentByOrder = `-- name: GetPaymentByOrder :one
SELECT payment_id, order_id, payment_method, amount, status, created_at FROM "Payments" WHERE order_id = $1
`

func (q *Queries) GetPaymentByOrder(ctx context.Context, orderID int64) (Payment, error) {
	row := q.db.QueryRow(ctx, getPaymentByOrder, orderID)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.OrderID,
		&i.PaymentMethod,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}
