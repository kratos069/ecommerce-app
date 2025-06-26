package db

import (
	"context"
	"testing"

	"github.com/e-commerce/util"
	"github.com/stretchr/testify/require"
)

func createRandomReview(t *testing.T) Review {
	user := createRandomUser(t)
	product := createRandomProduct(t, 5)

	arg := CreateReviewParams{
		UserID:    user.UserID,
		ProductID: product.ProductID,
		Rating:    util.RandomRating(),
		Comment:   util.RandomComment(),
	}

	review, err := testStore.CreateReview(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, review)

	require.Equal(t, arg.UserID, review.UserID)
	require.Equal(t, arg.ProductID, review.ProductID)
	require.Equal(t, arg.Rating, review.Rating)
	require.Equal(t, arg.Comment, review.Comment)

	return review
}

func TestCreateReview(t *testing.T) {
	createRandomReview(t)
}

func TestGetReviewsForProduct(t *testing.T) {
	review := createRandomReview(t)

	reviews, err := testStore.GetReviewsForProduct(context.Background(), review.ProductID)
	require.NoError(t, err)
	require.NotEmpty(t, reviews)

	// found := false
	// for _, r := range reviews {
	// 	if r.ReviewID == review.ReviewID {
	// 		found = true
	// 		require.Equal(t, review.UserID, r.UserID)
	// 		require.Equal(t, review.Rating, r.Rating)
	// 		require.Equal(t, review.Comment, r.Comment)
	// 		break
	// 	}
	// }
	// require.True(t, found, "Created review not found in fetched reviews")
}
