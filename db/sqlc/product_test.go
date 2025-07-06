package db

import (
	"context"
	"testing"
	"time"

	"github.com/e-commerce/util"
	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T, stockQuantity int64) Product {
	arg := CreateProductParams{
		Name:          util.RandomProductTitle(),
		Description:   util.RandomDescription(),
		ProductImage:  util.RandomProductImageURL(),
		Price:         util.RandomPrice(),
		StockQuantity: stockQuantity,
		CategoryID:    util.RandomCategoryID(),
	}

	product, err := testStore.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.Description, product.Description)
	require.Equal(t, arg.ProductImage, product.ProductImage)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.StockQuantity, product.StockQuantity)
	require.Equal(t, arg.CategoryID, product.CategoryID)

	return product
}

func TestCreateProduct(t *testing.T) {
	createRandomProduct(t, 10)
}

func TestGetProduct(t *testing.T) {
	product1 := createRandomProduct(t, 10)
	product2, err := testStore.GetProduct(context.Background(),
		product1.ProductID)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ProductID, product2.ProductID)
	require.Equal(t, product1.CategoryID, product2.CategoryID)
	require.Equal(t, product1.Description, product2.Description)
	require.Equal(t, product1.ProductImage, product2.ProductImage)
	require.Equal(t, product1.StockQuantity, product2.StockQuantity)
	require.Equal(t, product1.Price, product2.Price)
	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
}

func TestListProducts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomProduct(t, 1)
	}

	arg := ListProductsParams{
		Limit:  5,
		Offset: 5,
	}

	products, err := testStore.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, products)
}

func TestListProductsByCategoryID(t *testing.T) {
	product := createRandomProduct(t, 10)
	products, err := testStore.ListProductsByCategory(context.Background(),
		product.CategoryID)
	require.NoError(t, err)
	require.NotEmpty(t, products)
}

func TestUpdateProduct(t *testing.T) {
	product1 := createRandomProduct(t, 10)

	arg := UpdateProductParams{
		ProductID:     product1.ProductID,
		Name:          util.RandomOwner(),
		ProductImage:  util.RandomProductImageURL(),
		StockQuantity: util.RandomInt(1, 50),
		Description:   util.RandomDescription(),
		Price:         util.RandomPrice(),
		CategoryID:    util.RandomCategoryID(),
	}

	product2, err := testStore.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ProductID, product2.ProductID)
}

func TestDeleteMovie(t *testing.T) {
	product1 := createRandomProduct(t, 10)
	err := testStore.DeleteProduct(context.Background(), product1.ProductID)
	require.NoError(t, err)

	product2, err := testStore.GetProduct(context.Background(), product1.ProductID)
	require.Error(t, err)
	require.Empty(t, product2)
}