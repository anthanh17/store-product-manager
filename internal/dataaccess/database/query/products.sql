-- name: CreateProduct :one
INSERT INTO products (
  name,
  description,
  price,
  stock_quantity,
  status,
  image_url
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: GetProduct :one
SELECT p.*,
       COALESCE(
           json_agg(
               json_build_object(
                   'id', c.id,
                   'name', c.name
               )
           ) FILTER (WHERE c.id IS NOT NULL), '[]'
       ) as categories
FROM products p
LEFT JOIN product_categories pc ON p.id = pc.product_id
LEFT JOIN categories c ON pc.category_id = c.id
WHERE p.id = $1
GROUP BY p.id
LIMIT 1;

-- name: GetProductReviews :many
SELECT r.id, r.product_id, r.user_id, u.username, r.rating, r.comment, r.created_at
FROM reviews r
JOIN users u ON r.user_id = u.id
WHERE r.product_id = $1
ORDER BY r.created_at DESC;


-- name: UpdateProduct :one
UPDATE products
SET
  name = $2,
  description = $3,
  price = $4,
  stock_quantity = $5,
  status = $6,
  image_url = $7,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: DeleteProductCategories :exec
DELETE FROM product_categories
WHERE product_id = $1;

-- name: GetProductCategories :many
SELECT c.*
FROM categories c
JOIN product_categories pc ON c.id = pc.category_id
WHERE pc.product_id = $1;

-- name: AddProductCategory :exec
INSERT INTO product_categories (
  product_id,
  category_id
) VALUES (
  $1, $2
);

-- name: RemoveProductCategory :exec
DELETE FROM product_categories
WHERE product_id = $1 AND category_id = $2;

-- name: RemoveAllProductCategories :exec
DELETE FROM product_categories
WHERE product_id = $1;

-- name: ListProductsWithFilters :many
SELECT p.*,
       COALESCE(
           json_agg(
               json_build_object(
                   'id', c.id,
                   'name', c.name
               )
           ) FILTER (WHERE c.id IS NOT NULL), '[]'
       ) as categories
FROM products p
LEFT JOIN product_categories pc ON p.id = pc.product_id
LEFT JOIN categories c ON pc.category_id = c.id
WHERE
    (@status::text = '' OR p.status = @status::text)
    AND
    (@search_product_name::text = '' OR p.name ILIKE @search_product_name)
GROUP BY p.id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;



-- name: CountProductsWithFilters :one
SELECT COUNT(DISTINCT p.id)
FROM products p
WHERE
    (@status::text = '' OR p.status = @status::text)
    AND
    (@search_product_name::text = '' OR p.name ILIKE @search_product_name);
