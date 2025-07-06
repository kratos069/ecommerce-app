package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOrderAndPaymentTx_Success(t *testing.T) {
	user := createRandomUser(t)
	product1 := createRandomProduct(t, 10)
	product2 := createRandomProduct(t, 5)

	// Add to cart
	_, err := testStore.AddToCartTx(context.Background(), AddToCartTxParams{
		UserID:    user.UserID,
		ProductID: product1.ProductID,
		Quantity:  2,
	})
	require.NoError(t, err)

	_, err = testStore.AddToCartTx(context.Background(), AddToCartTxParams{
		UserID:    user.UserID,
		ProductID: product2.ProductID,
		Quantity:  1,
	})
	require.NoError(t, err)

	arg := CreateOrderAndPaymentTxParams{
		UserID:        user.UserID,
		PaymentMethod: "stripe",
	}

	result, err := testStore.CreateOrderAndPaymentTx(
		context.Background(), arg)
	require.NoError(t, err)

	// Check Order
	order := result.Order
	require.NotZero(t, order.OrderID)
	require.Equal(t, user.UserID, order.UserID)
	require.Equal(t, "completed", order.Status)
	require.True(t, order.TotalAmount > 0)

	// Check Payment
	payment := result.Payment
	require.Equal(t, order.OrderID, payment.OrderID)
	require.Equal(t, "paid", payment.Status)
	require.Equal(t, "stripe", payment.PaymentMethod)
	require.Equal(t, order.TotalAmount, payment.Amount)

	// Check cart is empty
	cart, err := testStore.GetCartItemsByUser(
		context.Background(), user.UserID)
	require.NoError(t, err)
	require.Len(t, cart, 0)
}

func TestCreateOrderAndPaymentTx_EmptyCart(t *testing.T) {
	user := createRandomUser(t)

	arg := CreateOrderAndPaymentTxParams{
		UserID:        user.UserID,
		PaymentMethod: "card",
	}

	result, err := testStore.CreateOrderAndPaymentTx(
		context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cart is empty")
	require.Empty(t, result.Order)
	require.Empty(t, result.Payment)
}

func TestCreateOrderAndPaymentTx_ProductDeleted(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t, 5)

	// Add product to cart
	addResult, err := testStore.AddToCartTx(
		context.Background(), AddToCartTxParams{
			UserID:    user.UserID,
			ProductID: product.ProductID,
			Quantity:  2,
		})
	require.NoError(t, err)

	// Manually delete cart item (to allow product deletion)
	err = testStore.DeleteCartItem(
		context.Background(), addResult.CartItem.CartItemID)
	require.NoError(t, err)

	// Now delete the product safely
	err = testStore.DeleteProduct(
		context.Background(), product.ProductID)
	require.NoError(t, err)

	// Re-add a cart item pointing to now-nonexistent product 
	// (simulate data corruption / race)
	_, err = testStore.AddCartItem(
		context.Background(), AddCartItemParams{
		UserID:    user.UserID,
		ProductID: product.ProductID, // no longer exists
		Quantity:  1,
		Price:     product.Price,
	})
	require.Error(t, err)

	// Try creating order + payment
	arg := CreateOrderAndPaymentTxParams{
		UserID:        user.UserID,
		PaymentMethod: "cash",
	}

	result, err := testStore.CreateOrderAndPaymentTx(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, result.Order)
	require.Empty(t, result.Payment)
}
