-- name: createAccount :one
INSERT INTO accounts (owner, balance, currency)
    VALUES ($1, $2, $3) RETURNING *;

-- name: getAccount :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts 
ORDER by id 
LIMIT $1
OFFSET $2;

-- name: updateAccount :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING *;

-- name: deleteAccount :exec
DELETE FROM accounts
WHERE id = $1;