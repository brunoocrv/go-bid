-- name: CreateProduct :one
INSERT INTO products (seller_id, name, description, base_price, auction_end)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetProductById :one
SELECT id, seller_id, name, description, base_price, auction_end
FROM products
WHERE id = $1;
