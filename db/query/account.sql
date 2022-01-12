-- name: CreateAccount :one
INSERT INTO accounts (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1;

-- name: GetAccountByEmail :one
SELECT * FROM accounts WHERE email = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateAccountEmailVerificationByID :one
UPDATE accounts SET is_email_verified = $2 WHERE id = $1 RETURNING *;

-- name: UpdateAccountEmailVerificationByEmail :one
UPDATE accounts SET is_email_verified = $2 WHERE email = $1 RETURNING *;

-- name: UpdateAccountPasswordByEmail :exec
UPDATE accounts SET password = $2 WHERE email = $1;