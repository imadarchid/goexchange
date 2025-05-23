// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const CreateUser = `-- name: CreateUser :one
INSERT INTO users (first_name, last_name, dob, email, balance)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type CreateUserParams struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Dob       time.Time `json:"dob"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, CreateUser,
		arg.FirstName,
		arg.LastName,
		arg.Dob,
		arg.Email,
		arg.Balance,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const DeleteUser = `-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, DeleteUser, id)
	err := row.Scan(&id)
	return id, err
}

const GetAllUsers = `-- name: GetAllUsers :many
SELECT id, first_name, last_name, email, dob, balance, created_at FROM users
ORDER BY created_at DESC
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, GetAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Dob,
			&i.Balance,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateUser = `-- name: UpdateUser :one
UPDATE users
SET first_name = $1, last_name = $2, dob = $3, balance = $4
WHERE id = $5
RETURNING id
`

type UpdateUserParams struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Dob       time.Time `json:"dob"`
	Balance   float64   `json:"balance"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, UpdateUser,
		arg.FirstName,
		arg.LastName,
		arg.Dob,
		arg.Balance,
		arg.ID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const UpdateUserBalance = `-- name: UpdateUserBalance :one
UPDATE users
SET balance = $1
WHERE id = $2
RETURNING id
`

type UpdateUserBalanceParams struct {
	Balance float64   `json:"balance"`
	ID      uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserBalance(ctx context.Context, arg UpdateUserBalanceParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, UpdateUserBalance, arg.Balance, arg.ID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
