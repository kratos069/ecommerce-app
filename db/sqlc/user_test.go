package db

import (
	"context"
	"testing"
	"time"

	"github.com/e-commerce/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:     util.RandomOwner(),
		PasswordHash: hashedPassword,
		Email:        util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUserByEmail(context.Background(),
		user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestGetUserByID(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUserByID(context.Background(),
		user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	err := testStore.DeleteUser(context.Background(), user.UserID)
	require.NoError(t, err)

	user2, err := testStore.GetUserByID(context.Background(), user.UserID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListRandomUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testStore.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
}
