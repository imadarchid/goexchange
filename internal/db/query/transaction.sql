-- name: CreateTransaction :one
INSERT INTO transactions (price, amount, buyer_order, seller_order, asset, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetTransactionById :one
SELECT * FROM transactions
WHERE id = $1;

-- name: GetTransactionsByUser :many
SELECT * FROM transactions
JOIN orders ON transactions.buyer_order = orders.id OR transactions.seller_order = orders.id
WHERE orders.created_by = $1
ORDER BY orders.created_at DESC;

-- name: GetAllTransactions :many
SELECT * FROM transactions
ORDER BY created_at DESC;
