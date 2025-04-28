-- name: CreateUser :one
INSERT INTO users (first_name, last_name, dob, email, balance)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;
-- name: UpdateUser :one
UPDATE users
SET first_name = $1, last_name = $2, dob = $3, balance = $4
WHERE id = $5
RETURNING id;

-- name: UpdateUserBalance :one
UPDATE users
SET balance = $1
WHERE id = $2
RETURNING id;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at DESC;
