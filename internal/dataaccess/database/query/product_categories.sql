-- name: CreateProductCategory :one
INSERT INTO product_categories (
  product_id,
  category_id
) VALUES (
  $1, $2
) RETURNING *;
