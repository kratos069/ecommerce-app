package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type AddToCartTxParams struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type AddToCartTxResult struct {
	CartItem CartItem `json:"cart_item"`
}

func (store *SQLStore) AddToCartTx(ctx context.Context,
	arg AddToCartTxParams) (AddToCartTxResult, error) {
	var result AddToCartTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// Lock product row
		product, err := q.GetProductForUpdate(ctx, arg.ProductID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) ||
				strings.Contains(err.Error(), "no rows") {
				return fmt.Errorf("product not found")
			}
			return err
		}

		// Check stock
		if product.StockQuantity < arg.Quantity {
			return fmt.Errorf(
				"not enough stock: available %d, requested %d",
				product.StockQuantity, arg.Quantity)
		}

		// Decrease stock
		if err := q.DecreaseProductStock(ctx,
			DecreaseProductStockParams{
				StockQuantity: arg.Quantity,
				ProductID:     arg.ProductID,
			}); err != nil {
			return fmt.Errorf("failed to decrease stock: %w", err)
		}

		// Try to fetch existing cart item
		existing, err := q.GetCartItemByUserAndProduct(ctx,
			GetCartItemByUserAndProductParams{
				UserID:    arg.UserID,
				ProductID: arg.ProductID,
			})
		// considering any “no rows” as “not existing”
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) ||
				strings.Contains(err.Error(), "no rows") {
				// insert path
				item, err := q.AddCartItem(ctx, AddCartItemParams{
					UserID:    arg.UserID,
					ProductID: arg.ProductID,
					Quantity:  arg.Quantity,
					Price:     product.Price,
				})
				if err != nil {
					return fmt.Errorf("failed to add cart item: %w", err)
				}
				result.CartItem = item
				return nil
			}
			return fmt.Errorf("failed to check existing cart: %w", err)
		}

		// Update existing cart item
		newQty := existing.Quantity + arg.Quantity
		if _, err := q.UpdateCartItemQuantity(
			ctx, UpdateCartItemQuantityParams{
				Quantity:  newQty,
				UserID:    arg.UserID,
				ProductID: arg.ProductID,
			}); err != nil {
			return fmt.Errorf(
				"failed to update cart quantity: %w", err)
		}

		updatedCartItem, err := q.GetCartItemByUserAndProduct(
			ctx, GetCartItemByUserAndProductParams{
				UserID:    arg.UserID,
				ProductID: arg.ProductID,
			})
		if err != nil {
			return fmt.Errorf(
				"failed to fetch updated cart item: %w", err)
		}
		result.CartItem = updatedCartItem
		return nil
	})

	return result, err
}
