-- name: GetCategorySummary :many
SELECT
    c.id,
    c.name,
    c.description,
    COUNT(pc.product_id) as product_count
FROM
    categories c
LEFT JOIN
    product_categories pc ON c.id = pc.category_id
GROUP BY
    c.id, c.name, c.description
ORDER BY
    c.name;

-- name: GetTotalCategories :one
SELECT
    COUNT(*) as total_categories
FROM
    categories;

-- name: GetTotalProducts :one
SELECT
    COUNT(*) as total_products
FROM
    products;

-- name: GetProductStatusSummary :many
SELECT
    status,
    COUNT(*) as count
FROM
    products
GROUP BY
    status;
