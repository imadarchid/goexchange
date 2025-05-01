-- name: CreateOrder :one
INSERT INTO orders (price, amount, side, order_type, asset, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetOrderById :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrdersByUser :many
SELECT * FROM orders
WHERE created_by = $1
ORDER BY created_at DESC;

-- name: GetAllOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET order_status = $1
WHERE id = $2
RETURNING id;

-- name: GetSubmittedOrders :many
SELECT * FROM orders
WHERE order_status = 'SUBMITTED' OR order_status = 'PARTIALLY_FILLED'
ORDER BY created_at DESC;