-- name: GetCategory :one
SELECT * FROM "Categories" 
WHERE category_id = $1;

-- name: ListCategories :many
SELECT * FROM "Categories" 
ORDER BY category_id;
