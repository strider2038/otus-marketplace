-- name: FindUser :one
SELECT id, email, password, role, first_name, last_name, phone, created_at, updated_at
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: FindUserByEmail :one
SELECT id, email, password, role, first_name, last_name, phone, created_at, updated_at
FROM "user"
WHERE email = $1
LIMIT 1;

-- name: CountUsersByEmail :one
SELECT count(id)
FROM "user"
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO "user" (id, email, password, role, first_name, last_name, phone)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, email, password, role, first_name, last_name, phone, created_at, updated_at;

-- name: UpdateUser :one
UPDATE "user"
SET email = $2, role = $3, first_name = $4, last_name = $5, phone = $6, updated_at = now()
WHERE id = $1
RETURNING id, email, password, role, first_name, last_name, phone, created_at, updated_at;

-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE id = $1;
