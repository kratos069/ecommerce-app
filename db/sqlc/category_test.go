package db

import (
	"context"
	"testing"

	"github.com/e-commerce/util"
	"github.com/stretchr/testify/require"
)

func TestGetCategory(t *testing.T) {
	category, err := testStore.GetCategory(context.Background(),
		util.RandomCategoryID())
	require.NoError(t, err)
	require.NotEmpty(t, category)
}

func TestListCategories(t *testing.T) {
	categories, err := testStore.ListCategories(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, categories)
	require.Equal(t, len(categories), 5)
}
