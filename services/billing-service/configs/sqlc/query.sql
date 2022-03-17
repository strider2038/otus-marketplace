-- name: FindAccount :one
SELECT id, amount, created_at, updated_at
FROM "account"
WHERE id = $1
LIMIT 1;

-- name: FindAccountForUpdate :one
SELECT id, amount, created_at, updated_at
FROM "account"
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: CreateAccount :one
INSERT INTO "account" (id)
VALUES ($1)
RETURNING id, amount, created_at, updated_at;

-- name: UpdateAccount :one
UPDATE "account"
SET amount = $2, updated_at = now()
WHERE id = $1
RETURNING id, amount, created_at, updated_at;

-- name: FindOperationsByAccount :many
SELECT id, account_id, amount, "type", description, created_at
FROM "operation"
WHERE account_id = $1
ORDER BY created_at DESC;

-- name: CreateOperation :one
INSERT INTO "operation" (id, account_id, amount, "type", description)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, account_id, amount, "type", description, created_at;
