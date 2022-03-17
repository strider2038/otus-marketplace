-- name: FindUser :one
SELECT id, email, password, first_name, last_name, phone, created_at, updated_at
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: FindUserByEmail :one
SELECT id, email, password, first_name, last_name, phone, created_at, updated_at
FROM "user"
WHERE email = $1
LIMIT 1;

-- name: CountUsersByEmail :one
SELECT count(id)
FROM "user"
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO "user" (id, email, password, first_name, last_name, phone)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, password, first_name, last_name, phone, created_at, updated_at;

-- name: UpdateUser :one
UPDATE "user"
SET email = $2, first_name = $3, last_name = $4, phone = $5, updated_at = now()
WHERE id = $1
RETURNING id, email, password, first_name, last_name, phone, created_at, updated_at;

-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE id = $1;
