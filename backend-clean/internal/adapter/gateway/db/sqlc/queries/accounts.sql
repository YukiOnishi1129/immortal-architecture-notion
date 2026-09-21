-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: GetAccountByProvider :one
SELECT *
FROM accounts
WHERE provider = $1
  AND provider_account_id = $2;

-- name: GetAccountByEmail :one
SELECT *
FROM accounts
WHERE email = $1;

-- name: UpsertAccount :one
INSERT INTO accounts (
    email,
    first_name,
    last_name,
    provider,
    provider_account_id,
    thumbnail,
    last_login_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (provider, provider_account_id)
DO UPDATE SET
    email = EXCLUDED.email,
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    thumbnail = EXCLUDED.thumbnail,
    last_login_at = EXCLUDED.last_login_at,
    updated_at = NOW()
RETURNING *;
