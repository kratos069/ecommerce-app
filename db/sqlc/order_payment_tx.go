package db

import (
	"context"
	"fmt"
)

type CreateOrderAndPaymentTxParams struct {
	UserID        int64  `json:"user_id"`
	PaymentMethod string `json:"payment_method"`
}

type CreateOrderAndPaymentTxResult struct {
	Order   Order   `json:"order"`
	Payment Payment `json:"payment"`
}

func (store *SQLStore) CreateOrderAndPaymentTx(
	ctx context.Context,
	arg CreateOrderAndPaymentTxParams,
) (CreateOrderAndPaymentTxResult, error) {

	var result CreateOrderAndPaymentTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// Fetch cart items
		cartItems, err := q.GetCartItemsByUser(ctx, arg.UserID)
		if err != nil {
			return fmt.Errorf("failed to fetch cart: %w", err)
		}
		if len(cartItems) == 0 {
			return fmt.Errorf("cart is empty")
		}

		// Calculate total amount
		var totalAmount float64
		for _, item := range cartItems {
			product, err := q.GetProduct(ctx, item.ProductID)
			if err != nil {
				return fmt.Errorf("failed to fetch product price: %w", err)
			}
			totalAmount += float64(item.Quantity) * product.Price
		}

		// Create order
		order, err := q.CreateOrder(ctx, CreateOrderParams{
			UserID:      arg.UserID,
			TotalAmount: totalAmount,
			Status:      "completed",
		})
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
		result.Order = order

		// Create payment
		payment, err := q.CreatePayment(ctx, CreatePaymentParams{
			OrderID:       order.OrderID,
			PaymentMethod: arg.PaymentMethod,
			Amount:        totalAmount,
			Status:        "paid",
		})
		if err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}
		result.Payment = payment

		// Clear cart
		for _, item := range cartItems {
			err := q.DeleteCartItem(ctx, item.CartItemID)
			if err != nil {
				return fmt.Errorf("failed to clear cart item %d: %w", item.CartItemID, err)
			}
		}

		return nil
	})

	return result, err
}
