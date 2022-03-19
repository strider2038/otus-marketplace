-- name: FindItem :one
SELECT id, name, initial_count, initial_price, commission_percent, created_at
FROM "item"
WHERE id = $1
LIMIT 1;

-- name: FindAllItems :many
SELECT id, name, initial_count, initial_price, commission_percent, created_at
FROM "item"
ORDER BY name;

-- name: AddItem :exec
INSERT INTO "item" (id, name, initial_count, initial_price, commission_percent, created_at)
VALUES ($1, $2, $3, $4, $5, $6);


-- name: FindUserItemByIDForUpdate :one
SELECT id, item_id, user_id, is_on_sale, created_at, updated_at
FROM "user_item"
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: FindUserItemForSale :one
SELECT id, item_id, user_id, is_on_sale, created_at, updated_at
FROM "user_item"
WHERE user_id = $1 AND is_on_sale IS FALSE
ORDER BY created_at
LIMIT 1
FOR UPDATE;

-- name: FindAggregatedUserItems :many
SELECT i.id AS id, i.name AS name, count(*) AS count, count(CASE WHEN ui.is_on_sale THEN 1 END) as on_sale_count
FROM "user_item" ui
INNER JOIN "item" i ON (ui.item_id = i.id)
WHERE ui.user_id = $1
GROUP BY i.id
ORDER BY i.name;

-- name: CreateUserItem :exec
INSERT INTO "user_item" (id, item_id, user_id, is_on_sale)
VALUES ($1, $2, $3, $4);

-- name: UpdateUserItem :exec
UPDATE "user_item"
SET is_on_sale = $2,
    updated_at = now()
WHERE id = $1;

-- name: DeleteUserItem :exec
DELETE FROM "user_item"
WHERE id = $1;


-- name: FindPurchaseOrder :one
SELECT id,
       user_id,
       item_id,
       payment_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "purchase_order"
WHERE id = $1
LIMIT 1;

-- name: FindPurchaseOrderForUpdate :one
SELECT id,
       user_id,
       item_id,
       payment_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "purchase_order"
WHERE id = $1
LIMIT 1 FOR UPDATE;

-- name: FindPurchaseOrderByPaymentForUpdate :one
SELECT id,
       user_id,
       item_id,
       payment_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "purchase_order"
WHERE payment_id = $1
LIMIT 1 FOR UPDATE;

-- name: FindPurchaseOrderByDealForUpdate :one
SELECT id,
       user_id,
       item_id,
       payment_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "purchase_order"
WHERE deal_id = $1
LIMIT 1 FOR UPDATE;

-- name: FindPurchaseOrdersByUser :many
SELECT id,
       user_id,
       item_id,
       payment_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "purchase_order"
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: FindPurchaseOrderForDeal :one
SELECT id,
       user_id,
       item_id,
       payment_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "purchase_order"
WHERE item_id = $1
  AND price <= $2
  AND user_id != $3
ORDER BY price DESC
LIMIT 1 FOR UPDATE;

-- name: CreatePurchaseOrder :exec
INSERT INTO "purchase_order" (id, user_id, item_id, payment_id, deal_id, price, commission, status, created_at,
                              updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: UpdatePurchaseOrder :exec
UPDATE "purchase_order"
SET status     = $2,
    payment_id = $3,
    deal_id    = $4,
    updated_at = now()
WHERE id = $1;

-- name: DeletePurchaseOrder :exec
DELETE FROM "purchase_order"
WHERE id = $1;


-- name: FindSellOrder :one
SELECT id,
       user_id,
       item_id,
       user_item_id,
       accrual_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "sell_order"
WHERE id = $1
LIMIT 1;

-- name: FindSellOrderForUpdate :one
SELECT id,
       user_id,
       item_id,
       user_item_id,
       accrual_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "sell_order"
WHERE id = $1
LIMIT 1 FOR UPDATE;

-- name: FindSellOrderByAccrualForUpdate :one
SELECT id,
       user_id,
       item_id,
       user_item_id,
       accrual_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "sell_order"
WHERE accrual_id = $1
LIMIT 1 FOR UPDATE;

-- name: FindSellOrderByDealForUpdate :one
SELECT id,
       user_id,
       item_id,
       user_item_id,
       accrual_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "sell_order"
WHERE deal_id = $1
LIMIT 1 FOR UPDATE;

-- name: FindSellOrdersByUser :many
SELECT id,
       user_id,
       item_id,
       user_item_id,
       accrual_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "sell_order"
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: FindSellOrderForDeal :one
SELECT id,
       user_id,
       item_id,
       user_item_id,
       accrual_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "sell_order"
WHERE item_id = $1
  AND price <= $2
  AND user_id != $3
ORDER BY price DESC
LIMIT 1 FOR UPDATE;

-- name: CreateSellOrder :exec
INSERT INTO "sell_order" (id, user_id, item_id, user_item_id, accrual_id, deal_id, price, commission, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: UpdateSellOrder :exec
UPDATE "sell_order"
SET status     = $2,
    accrual_id = $3,
    deal_id    = $4,
    updated_at = now()
WHERE id = $1;

-- name: DeleteSellOrder :exec
DELETE FROM "sell_order"
WHERE id = $1;
