-- name: CreateProduct :one
INSERT INTO products (
  name, sku, price
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :exec
UPDATE products
SET name = $1, price = $2
WHERE id = $3;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;