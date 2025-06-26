package db

import (
	"context"
	"testing"

	"github.com/e-commerce/util"
	"github.com/stretchr/testify/require"
)

func createRandomCartItem(t *testing.T) CartItem {
	user := createRandomUser(t)
	product := createRandomProduct(t, 12)

	arg := AddCartItemParams{
		UserID:    user.UserID,
		ProductID: product.ProductID,
		Quantity:  util.RandomInt(1, 10),
		Price:     util.RandomPrice(),
	}

	cartItem, err := testStore.AddCartItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cartItem)

	require.Equal(t, arg.UserID, cartItem.UserID)
	require.Equal(t, arg.ProductID, cartItem.ProductID)
	require.Equal(t, arg.Quantity, cartItem.Quantity)

	return cartItem
}

func TestAddCartItem(t *testing.T) {
	createRandomCartItem(t)
}

func TestGetCartItemsByUser(t *testing.T) {
	cartItem1 := createRandomCartItem(t)
	cartItems, err := testStore.GetCartItemsByUser(context.Background(),
		cartItem1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, cartItems)
}

func TestDeleteCartItem(t *testing.T) {
	carItem := createRandomCartItem(t)
	err := testStore.DeleteCartItem(context.Background(), carItem.CartItemID)
	require.NoError(t, err)
}
