package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddToCartTx_NewItem(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t, 10)

	arg := AddToCartTxParams{
		UserID:    user.UserID,
		ProductID: product.ProductID,
		Quantity:  3,
	}

	result, err := testStore.AddToCartTx(context.Background(), arg)
	require.NoError(t, err)

	// Verify returned cart item
	cart := result.CartItem
	require.NotZero(t, cart.CartItemID)
	require.Equal(t, arg.UserID, cart.UserID)
	require.Equal(t, arg.ProductID, cart.ProductID)
	require.Equal(t, arg.Quantity, cart.Quantity)
	require.Equal(t, product.Price, cart.Price)

	// Verify stock was reduced
	updatedProduct, err := testStore.GetProduct(
		context.Background(), cart.ProductID)
	require.NoError(t, err)
	require.Equal(
		t, product.StockQuantity-arg.Quantity,
		updatedProduct.StockQuantity)
}

func TestAddToCartTx_UpdateExisting(t *testing.T) {
	user := createRandomUser(t)

	initialStock := int64(10)
	product := createRandomProduct(t, initialStock)

	// First add to cart
	initialQty := int64(5)
	_, err := testStore.AddToCartTx(
		context.Background(), AddToCartTxParams{
			UserID:    user.UserID,
			ProductID: product.ProductID,
			Quantity:  initialQty,
		})
	require.NoError(t, err)

	// Second add to cart
	additionalQty := int64(4)
	result, err := testStore.AddToCartTx(context.Background(),
		AddToCartTxParams{
			UserID:    user.UserID,
			ProductID: product.ProductID,
			Quantity:  additionalQty,
		})
	require.NoError(t, err)

	// Verify total quantity in cart
	expectedQty := initialQty + additionalQty
	require.Equal(t, expectedQty, result.CartItem.Quantity)

	// Verify total stock remaining
	expectedStock := initialStock - expectedQty

	updatedProduct, err := testStore.GetProduct(context.Background(),
		result.CartItem.ProductID)
	require.NoError(t, err)
	require.Equal(t, expectedStock, updatedProduct.StockQuantity)
}

func TestAddToCartTx_InsufficientStock(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t, 2)

	// check for only 2 in stock

	arg := AddToCartTxParams{
		UserID:    user.UserID,
		ProductID: product.ProductID,
		Quantity:  5, // more than stock
	}

	result, err := testStore.AddToCartTx(context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not enough stock")
	require.Empty(t, result.CartItem)
}

func TestAddToCartTx_ProductNotFound(t *testing.T) {
	arg := AddToCartTxParams{
		UserID:    createRandomUser(t).UserID,
		ProductID: 999999, // non-existent
		Quantity:  1,
	}

	result, err := testStore.AddToCartTx(context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "product not found")
	require.Empty(t, result.CartItem)
}
