-- name: CreateAccount :one
INSERT INTO accounts (
    email,
    password_hash,
    phone
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByEmail :one
SELECT * FROM accounts
WHERE email = $1 LIMIT 1;

-- name: GetAccountByPhone :one
SELECT * FROM accounts
WHERE phone = $1 LIMIT 1;


-- name: UpdateAccount :one
UPDATE accounts
SET
    phone = COALESCE(sqlc.narg('phone'), phone)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdatePassword :exec
UPDATE accounts
SET password_hash = $2
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAccounts :one
SELECT COUNT(*) FROM accounts;
