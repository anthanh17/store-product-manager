-- name: CreateCategory :one
INSERT INTO categories (
  name,
  description
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1 LIMIT 1;

-- name: ListCategoriesWithPagination :many
SELECT * FROM categories
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories;

-- name: UpdateCategory :one
UPDATE categories
SET
  name = $2,
  description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: CountCategoriesWithName :one
SELECT COUNT(*) FROM categories
WHERE name = $1;
