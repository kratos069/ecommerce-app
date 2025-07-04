// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: category.sql

package db

import (
	"context"
)

const getCategory = `-- name: GetCategory :one
SELECT category_id, name FROM "Categories" 
WHERE category_id = $1
`

func (q *Queries) GetCategory(ctx context.Context, categoryID int64) (Category, error) {
	row := q.db.QueryRow(ctx, getCategory, categoryID)
	var i Category
	err := row.Scan(&i.CategoryID, &i.Name)
	return i, err
}

const listCategories = `-- name: ListCategories :many
SELECT category_id, name FROM "Categories" 
ORDER BY category_id
`

func (q *Queries) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.Query(ctx, listCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.CategoryID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
