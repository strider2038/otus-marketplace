-- name: FindUser :one
SELECT id, email, first_name, last_name, phone, created_at, updated_at
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (id, email, first_name, last_name, phone)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, first_name, last_name, phone, created_at, updated_at;

-- name: UpdateUser :one
UPDATE "user"
SET email = $2, first_name = $3, last_name = $4, phone = $5, updated_at = now()
WHERE id = $1
RETURNING id, email, first_name, last_name, phone, created_at, updated_at;

-- name: FindNotificationsByUser :many
SELECT id, user_id, message, created_at
FROM "notification"
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateNotification :one
INSERT INTO "notification" (id, user_id, message)
VALUES ($1, $2, $3)
RETURNING id, user_id, message, created_at;
