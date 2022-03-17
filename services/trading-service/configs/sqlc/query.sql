-- name: FindItem :one
SELECT id, name, initial_count, initial_price, commission_percent, created_at
FROM "item"
WHERE id = $1
LIMIT 1;

-- name: FindAllItems :many
SELECT id, name, initial_count, initial_price, commission_percent, created_at
FROM "item"
ORDER BY name;

-- name: AddItem :one
INSERT INTO "item" (id, name, initial_count, initial_price, commission_percent, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, initial_count, initial_price, commission_percent, created_at;


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
WHERE id = $1
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
ORDER BY price DESC
LIMIT 1 FOR UPDATE;

-- name: CreatePurchaseOrder :one
INSERT INTO "purchase_order" (id, user_id, item_id, payment_id, deal_id, price, commission, status, created_at,
                              updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, user_id, item_id, payment_id, deal_id, price, commission, status, created_at, updated_at;

-- name: UpdatePurchaseOrder :one
UPDATE "purchase_order"
SET status     = $2,
    payment_id = $3,
    deal_id    = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, item_id, payment_id, deal_id, price, commission, status, created_at, updated_at;


-- name: FindSellOrder :one
SELECT id,
       user_id,
       item_id,
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
       accrual_id,
       deal_id,
       price,
       commission,
       status,
       created_at,
       updated_at
FROM "sell_order"
WHERE id = $1
ORDER BY created_at DESC;

-- name: FindSellOrderForDeal :one
SELECT id,
       user_id,
       item_id,
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
ORDER BY price DESC
LIMIT 1 FOR UPDATE;

-- name: CreateSellOrder :one
INSERT INTO "sell_order" (id, user_id, item_id, accrual_id, deal_id, price, commission, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, user_id, item_id, accrual_id, deal_id, price, commission, status, created_at, updated_at;

-- name: UpdateSellOrder :one
UPDATE "sell_order"
SET status     = $2,
    accrual_id = $3,
    deal_id    = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, item_id, accrual_id, deal_id, price, commission, status, created_at, updated_at;
