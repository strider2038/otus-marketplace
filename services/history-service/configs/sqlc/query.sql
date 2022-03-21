-- name: FindDealsByUser :many
SELECT id, user_id, item_id, item_name, "type", amount, commission, completed_at
FROM "deal"
WHERE user_id = $1
ORDER BY completed_at DESC;

-- name: AddDeal :exec
INSERT INTO "deal" (id, user_id, item_id, item_name, "type", amount, commission, completed_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
