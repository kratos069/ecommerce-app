// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: product.sql

package db

import (
	"context"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO "Products" 
(name, description, product_image, price, stock_quantity, category_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING product_id, name, description, product_image, price, stock_quantity, category_id, created_at
`

type CreateProductParams struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ProductImage  string  `json:"product_image"`
	Price         float64 `json:"price"`
	StockQuantity int64   `json:"stock_quantity"`
	CategoryID    int64   `json:"category_id"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		arg.Name,
		arg.Description,
		arg.ProductImage,
		arg.Price,
		arg.StockQuantity,
		arg.CategoryID,
	)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.ProductImage,
		&i.Price,
		&i.StockQuantity,
		&i.CategoryID,
		&i.CreatedAt,
	)
	return i, err
}

const decreaseProductStock = `-- name: DecreaseProductStock :exec
UPDATE "Products"
SET stock_quantity = stock_quantity - $1
WHERE product_id = $2 AND stock_quantity >= $1
`

type DecreaseProductStockParams struct {
	StockQuantity int64 `json:"stock_quantity"`
	ProductID     int64 `json:"product_id"`
}

func (q *Queries) DecreaseProductStock(ctx context.Context, arg DecreaseProductStockParams) error {
	_, err := q.db.Exec(ctx, decreaseProductStock, arg.StockQuantity, arg.ProductID)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM "Products"
WHERE product_id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, productID int64) error {
	_, err := q.db.Exec(ctx, deleteProduct, productID)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT product_id, name, description, product_image, price, stock_quantity, category_id, created_at FROM "Products" WHERE product_id = $1
`

func (q *Queries) GetProduct(ctx context.Context, productID int64) (Product, error) {
	row := q.db.QueryRow(ctx, getProduct, productID)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.ProductImage,
		&i.Price,
		&i.StockQuantity,
		&i.CategoryID,
		&i.CreatedAt,
	)
	return i, err
}

const getProductForUpdate = `-- name: GetProductForUpdate :one
SELECT product_id, name, description, product_image, price, stock_quantity, category_id, created_at FROM "Products"
WHERE product_id = $1
FOR NO KEY UPDATE
`

func (q *Queries) GetProductForUpdate(ctx context.Context, productID int64) (Product, error) {
	row := q.db.QueryRow(ctx, getProductForUpdate, productID)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.ProductImage,
		&i.Price,
		&i.StockQuantity,
		&i.CategoryID,
		&i.CreatedAt,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT product_id, name, description, product_image, price, stock_quantity, category_id, created_at FROM "Products" 
ORDER BY product_id 
LIMIT $1 OFFSET $2
`

type ListProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.db.Query(ctx, listProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ProductID,
			&i.Name,
			&i.Description,
			&i.ProductImage,
			&i.Price,
			&i.StockQuantity,
			&i.CategoryID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductsByCategory = `-- name: ListProductsByCategory :many
SELECT product_id, name, description, product_image, price, stock_quantity, category_id, created_at FROM "Products" 
WHERE category_id = $1 
ORDER BY product_id
`

func (q *Queries) ListProductsByCategory(ctx context.Context, categoryID int64) ([]Product, error) {
	rows, err := q.db.Query(ctx, listProductsByCategory, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ProductID,
			&i.Name,
			&i.Description,
			&i.ProductImage,
			&i.Price,
			&i.StockQuantity,
			&i.CategoryID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :one
UPDATE "Products"
SET name = $2, description = $3, product_image = $4,
stock_quantity = $5, price = $6, category_id = $7
WHERE product_id = $1
RETURNING product_id, name, description, product_image, price, stock_quantity, category_id, created_at
`

type UpdateProductParams struct {
	ProductID     int64   `json:"product_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ProductImage  string  `json:"product_image"`
	StockQuantity int64   `json:"stock_quantity"`
	Price         float64 `json:"price"`
	CategoryID    int64   `json:"category_id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.ProductID,
		arg.Name,
		arg.Description,
		arg.ProductImage,
		arg.StockQuantity,
		arg.Price,
		arg.CategoryID,
	)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.ProductImage,
		&i.Price,
		&i.StockQuantity,
		&i.CategoryID,
		&i.CreatedAt,
	)
	return i, err
}
