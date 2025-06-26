package db

import (
	"context"
	"testing"
	"time"

	"github.com/e-commerce/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T) Order {
	user := createRandomUser(t)

	arg := CreateOrderParams{
		UserID:      user.UserID,
		TotalAmount: util.RandomPrice(),
		Status:      util.RandomStatus(),
	}
	order, err := testStore.CreateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	require.Equal(t, arg.UserID, order.UserID)
	require.Equal(t, arg.TotalAmount, order.TotalAmount)
	require.Equal(t, arg.Status, order.Status)

	return order
}

func TestCreateOrder(t *testing.T) {
	createRandomOrder(t)
}

func TestGetOrder(t *testing.T) {
	order1 := createRandomOrder(t)
	order2, err := testStore.GetOrder(context.Background(), order1.OrderID)
	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.OrderID, order2.OrderID)
	require.Equal(t, order1.UserID, order2.UserID)
	require.Equal(t, order1.TotalAmount, order2.TotalAmount)
	require.Equal(t, order1.Status, order2.Status)
	require.WithinDuration(t, order1.CreatedAt, order2.CreatedAt, time.Second)
}

func TestListOrdersByUser(t *testing.T) {
	var orders Order
	for i := 0; i < 5; i++ {
		orders = createRandomOrder(t)
	}

	retOrders, err := testStore.ListOrdersByUser(context.Background(),
		orders.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, retOrders)
}
