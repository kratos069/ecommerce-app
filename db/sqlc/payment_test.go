package db

import (
	"context"
	"testing"

	"github.com/e-commerce/util"
	"github.com/stretchr/testify/require"
)

func createRandomPayment(t *testing.T) Payment {
	order := createRandomOrder(t)

	arg := CreatePaymentParams{
		OrderID:       order.OrderID,
		PaymentMethod: util.RandomPaymentMethod(),
		Amount:        util.RandomPrice(),
		Status:        util.RandomStatus(),
	}

	payment, err := testStore.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.OrderID, payment.OrderID)
	require.Equal(t, arg.PaymentMethod, payment.PaymentMethod)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.Status, payment.Status)

	return payment
}

func TestCreatePayment(t *testing.T) {
	createRandomPayment(t)
}

func TestGetPaymentByOrder(t *testing.T) {
	payment1 := createRandomPayment(t)

	payment2, err := testStore.GetPaymentByOrder(context.Background(),
		payment1.OrderID)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.OrderID, payment2.OrderID)
	require.Equal(t, payment1.PaymentID, payment2.PaymentID)
	require.Equal(t, payment1.Amount, payment2.Amount)
	require.Equal(t, payment1.Status, payment2.Status)
}
