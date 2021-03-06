// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package database

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

const addItem = `-- name: AddItem :exec
INSERT INTO "item" (id, name, initial_count, initial_price, commission_percent, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
`

type AddItemParams struct {
	ID                uuid.UUID
	Name              string
	InitialCount      int32
	InitialPrice      float64
	CommissionPercent float64
	CreatedAt         time.Time
}

func (q *Queries) AddItem(ctx context.Context, arg AddItemParams) error {
	_, err := q.db.Exec(ctx, addItem,
		arg.ID,
		arg.Name,
		arg.InitialCount,
		arg.InitialPrice,
		arg.CommissionPercent,
		arg.CreatedAt,
	)
	return err
}

const createPurchaseOrder = `-- name: CreatePurchaseOrder :exec
INSERT INTO "purchase_order" (id, user_id, item_id, payment_id, deal_id, price, commission, status, created_at,
                              updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

type CreatePurchaseOrderParams struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	ItemID     uuid.UUID
	PaymentID  uuid.NullUUID
	DealID     uuid.NullUUID
	Price      float64
	Commission float64
	Status     PurchaseStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) CreatePurchaseOrder(ctx context.Context, arg CreatePurchaseOrderParams) error {
	_, err := q.db.Exec(ctx, createPurchaseOrder,
		arg.ID,
		arg.UserID,
		arg.ItemID,
		arg.PaymentID,
		arg.DealID,
		arg.Price,
		arg.Commission,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const createSellOrder = `-- name: CreateSellOrder :exec
INSERT INTO "sell_order" (id, user_id, item_id, user_item_id, accrual_id, deal_id, price, commission, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

type CreateSellOrderParams struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	ItemID     uuid.UUID
	UserItemID uuid.UUID
	AccrualID  uuid.NullUUID
	DealID     uuid.NullUUID
	Price      float64
	Commission float64
	Status     SellStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) CreateSellOrder(ctx context.Context, arg CreateSellOrderParams) error {
	_, err := q.db.Exec(ctx, createSellOrder,
		arg.ID,
		arg.UserID,
		arg.ItemID,
		arg.UserItemID,
		arg.AccrualID,
		arg.DealID,
		arg.Price,
		arg.Commission,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const createUserItem = `-- name: CreateUserItem :exec
INSERT INTO "user_item" (id, item_id, user_id, is_on_sale)
VALUES ($1, $2, $3, $4)
`

type CreateUserItemParams struct {
	ID       uuid.UUID
	ItemID   uuid.UUID
	UserID   uuid.UUID
	IsOnSale bool
}

func (q *Queries) CreateUserItem(ctx context.Context, arg CreateUserItemParams) error {
	_, err := q.db.Exec(ctx, createUserItem,
		arg.ID,
		arg.ItemID,
		arg.UserID,
		arg.IsOnSale,
	)
	return err
}

const deletePurchaseOrder = `-- name: DeletePurchaseOrder :exec
DELETE FROM "purchase_order"
WHERE id = $1
`

func (q *Queries) DeletePurchaseOrder(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePurchaseOrder, id)
	return err
}

const deleteSellOrder = `-- name: DeleteSellOrder :exec
DELETE FROM "sell_order"
WHERE id = $1
`

func (q *Queries) DeleteSellOrder(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteSellOrder, id)
	return err
}

const deleteUserItem = `-- name: DeleteUserItem :exec
DELETE FROM "user_item"
WHERE id = $1
`

func (q *Queries) DeleteUserItem(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUserItem, id)
	return err
}

const findAggregatedUserItems = `-- name: FindAggregatedUserItems :many
SELECT i.id AS id, i.name AS name, count(*) AS count, count(CASE WHEN ui.is_on_sale THEN 1 END) as on_sale_count
FROM "user_item" ui
INNER JOIN "item" i ON (ui.item_id = i.id)
WHERE ui.user_id = $1
GROUP BY i.id
ORDER BY i.name
`

type FindAggregatedUserItemsRow struct {
	ID          uuid.UUID
	Name        string
	Count       int64
	OnSaleCount int64
}

func (q *Queries) FindAggregatedUserItems(ctx context.Context, userID uuid.UUID) ([]FindAggregatedUserItemsRow, error) {
	rows, err := q.db.Query(ctx, findAggregatedUserItems, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindAggregatedUserItemsRow
	for rows.Next() {
		var i FindAggregatedUserItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Count,
			&i.OnSaleCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findAllItems = `-- name: FindAllItems :many
SELECT id, name, initial_count, initial_price, commission_percent, created_at
FROM "item"
ORDER BY name
`

func (q *Queries) FindAllItems(ctx context.Context) ([]Item, error) {
	rows, err := q.db.Query(ctx, findAllItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.InitialCount,
			&i.InitialPrice,
			&i.CommissionPercent,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findItem = `-- name: FindItem :one
SELECT id, name, initial_count, initial_price, commission_percent, created_at
FROM "item"
WHERE id = $1
LIMIT 1
`

func (q *Queries) FindItem(ctx context.Context, id uuid.UUID) (Item, error) {
	row := q.db.QueryRow(ctx, findItem, id)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InitialCount,
		&i.InitialPrice,
		&i.CommissionPercent,
		&i.CreatedAt,
	)
	return i, err
}

const findPurchaseOrder = `-- name: FindPurchaseOrder :one
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
LIMIT 1
`

func (q *Queries) FindPurchaseOrder(ctx context.Context, id uuid.UUID) (PurchaseOrder, error) {
	row := q.db.QueryRow(ctx, findPurchaseOrder, id)
	var i PurchaseOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.PaymentID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findPurchaseOrderByDealForUpdate = `-- name: FindPurchaseOrderByDealForUpdate :one
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
LIMIT 1 FOR UPDATE
`

func (q *Queries) FindPurchaseOrderByDealForUpdate(ctx context.Context, dealID uuid.NullUUID) (PurchaseOrder, error) {
	row := q.db.QueryRow(ctx, findPurchaseOrderByDealForUpdate, dealID)
	var i PurchaseOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.PaymentID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findPurchaseOrderByPaymentForUpdate = `-- name: FindPurchaseOrderByPaymentForUpdate :one
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
LIMIT 1 FOR UPDATE
`

func (q *Queries) FindPurchaseOrderByPaymentForUpdate(ctx context.Context, paymentID uuid.NullUUID) (PurchaseOrder, error) {
	row := q.db.QueryRow(ctx, findPurchaseOrderByPaymentForUpdate, paymentID)
	var i PurchaseOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.PaymentID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findPurchaseOrderForDeal = `-- name: FindPurchaseOrderForDeal :one
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
LIMIT 1 FOR UPDATE SKIP LOCKED
`

type FindPurchaseOrderForDealParams struct {
	ItemID uuid.UUID
	Price  float64
	UserID uuid.UUID
}

func (q *Queries) FindPurchaseOrderForDeal(ctx context.Context, arg FindPurchaseOrderForDealParams) (PurchaseOrder, error) {
	row := q.db.QueryRow(ctx, findPurchaseOrderForDeal, arg.ItemID, arg.Price, arg.UserID)
	var i PurchaseOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.PaymentID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findPurchaseOrderForUpdate = `-- name: FindPurchaseOrderForUpdate :one
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
LIMIT 1 FOR UPDATE
`

func (q *Queries) FindPurchaseOrderForUpdate(ctx context.Context, id uuid.UUID) (PurchaseOrder, error) {
	row := q.db.QueryRow(ctx, findPurchaseOrderForUpdate, id)
	var i PurchaseOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.PaymentID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findPurchaseOrdersByUser = `-- name: FindPurchaseOrdersByUser :many
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
ORDER BY created_at DESC
`

func (q *Queries) FindPurchaseOrdersByUser(ctx context.Context, userID uuid.UUID) ([]PurchaseOrder, error) {
	rows, err := q.db.Query(ctx, findPurchaseOrdersByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PurchaseOrder
	for rows.Next() {
		var i PurchaseOrder
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ItemID,
			&i.PaymentID,
			&i.DealID,
			&i.Price,
			&i.Commission,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findSellOrder = `-- name: FindSellOrder :one
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
LIMIT 1
`

func (q *Queries) FindSellOrder(ctx context.Context, id uuid.UUID) (SellOrder, error) {
	row := q.db.QueryRow(ctx, findSellOrder, id)
	var i SellOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.UserItemID,
		&i.AccrualID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findSellOrderByAccrualForUpdate = `-- name: FindSellOrderByAccrualForUpdate :one
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
LIMIT 1 FOR UPDATE
`

func (q *Queries) FindSellOrderByAccrualForUpdate(ctx context.Context, accrualID uuid.NullUUID) (SellOrder, error) {
	row := q.db.QueryRow(ctx, findSellOrderByAccrualForUpdate, accrualID)
	var i SellOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.UserItemID,
		&i.AccrualID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findSellOrderByDealForUpdate = `-- name: FindSellOrderByDealForUpdate :one
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
LIMIT 1 FOR UPDATE
`

func (q *Queries) FindSellOrderByDealForUpdate(ctx context.Context, dealID uuid.NullUUID) (SellOrder, error) {
	row := q.db.QueryRow(ctx, findSellOrderByDealForUpdate, dealID)
	var i SellOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.UserItemID,
		&i.AccrualID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findSellOrderForDeal = `-- name: FindSellOrderForDeal :one
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
LIMIT 1 FOR UPDATE SKIP LOCKED
`

type FindSellOrderForDealParams struct {
	ItemID uuid.UUID
	Price  float64
	UserID uuid.UUID
}

func (q *Queries) FindSellOrderForDeal(ctx context.Context, arg FindSellOrderForDealParams) (SellOrder, error) {
	row := q.db.QueryRow(ctx, findSellOrderForDeal, arg.ItemID, arg.Price, arg.UserID)
	var i SellOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.UserItemID,
		&i.AccrualID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findSellOrderForUpdate = `-- name: FindSellOrderForUpdate :one
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
LIMIT 1 FOR UPDATE
`

func (q *Queries) FindSellOrderForUpdate(ctx context.Context, id uuid.UUID) (SellOrder, error) {
	row := q.db.QueryRow(ctx, findSellOrderForUpdate, id)
	var i SellOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ItemID,
		&i.UserItemID,
		&i.AccrualID,
		&i.DealID,
		&i.Price,
		&i.Commission,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findSellOrdersByUser = `-- name: FindSellOrdersByUser :many
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
ORDER BY created_at DESC
`

func (q *Queries) FindSellOrdersByUser(ctx context.Context, userID uuid.UUID) ([]SellOrder, error) {
	rows, err := q.db.Query(ctx, findSellOrdersByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SellOrder
	for rows.Next() {
		var i SellOrder
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ItemID,
			&i.UserItemID,
			&i.AccrualID,
			&i.DealID,
			&i.Price,
			&i.Commission,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUserItemByIDForUpdate = `-- name: FindUserItemByIDForUpdate :one
SELECT id, item_id, user_id, is_on_sale, created_at, updated_at
FROM "user_item"
WHERE id = $1
LIMIT 1
FOR UPDATE
`

func (q *Queries) FindUserItemByIDForUpdate(ctx context.Context, id uuid.UUID) (UserItem, error) {
	row := q.db.QueryRow(ctx, findUserItemByIDForUpdate, id)
	var i UserItem
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.UserID,
		&i.IsOnSale,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserItemForSale = `-- name: FindUserItemForSale :one
SELECT id, item_id, user_id, is_on_sale, created_at, updated_at
FROM "user_item"
WHERE user_id = $1 AND is_on_sale IS FALSE
ORDER BY created_at
LIMIT 1
FOR UPDATE
`

func (q *Queries) FindUserItemForSale(ctx context.Context, userID uuid.UUID) (UserItem, error) {
	row := q.db.QueryRow(ctx, findUserItemForSale, userID)
	var i UserItem
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.UserID,
		&i.IsOnSale,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPurchaseOrdersStateOfUser = `-- name: GetPurchaseOrdersStateOfUser :one
SELECT COALESCE(max(updated_at)::text, '')::text
FROM purchase_order
WHERE user_id = $1
`

func (q *Queries) GetPurchaseOrdersStateOfUser(ctx context.Context, userID uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, getPurchaseOrdersStateOfUser, userID)
	var column_1 string
	err := row.Scan(&column_1)
	return column_1, err
}

const getSellOrdersStateOfUser = `-- name: GetSellOrdersStateOfUser :one
SELECT COALESCE(max(updated_at)::text, '')::text
FROM sell_order
WHERE user_id = $1
`

func (q *Queries) GetSellOrdersStateOfUser(ctx context.Context, userID uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, getSellOrdersStateOfUser, userID)
	var column_1 string
	err := row.Scan(&column_1)
	return column_1, err
}

const updatePurchaseOrder = `-- name: UpdatePurchaseOrder :exec
UPDATE "purchase_order"
SET status     = $2,
    payment_id = $3,
    deal_id    = $4,
    updated_at = now()
WHERE id = $1
`

type UpdatePurchaseOrderParams struct {
	ID        uuid.UUID
	Status    PurchaseStatus
	PaymentID uuid.NullUUID
	DealID    uuid.NullUUID
}

func (q *Queries) UpdatePurchaseOrder(ctx context.Context, arg UpdatePurchaseOrderParams) error {
	_, err := q.db.Exec(ctx, updatePurchaseOrder,
		arg.ID,
		arg.Status,
		arg.PaymentID,
		arg.DealID,
	)
	return err
}

const updateSellOrder = `-- name: UpdateSellOrder :exec
UPDATE "sell_order"
SET status     = $2,
    accrual_id = $3,
    deal_id    = $4,
    updated_at = now()
WHERE id = $1
`

type UpdateSellOrderParams struct {
	ID        uuid.UUID
	Status    SellStatus
	AccrualID uuid.NullUUID
	DealID    uuid.NullUUID
}

func (q *Queries) UpdateSellOrder(ctx context.Context, arg UpdateSellOrderParams) error {
	_, err := q.db.Exec(ctx, updateSellOrder,
		arg.ID,
		arg.Status,
		arg.AccrualID,
		arg.DealID,
	)
	return err
}

const updateUserItem = `-- name: UpdateUserItem :exec
UPDATE "user_item"
SET is_on_sale = $2,
    updated_at = now()
WHERE id = $1
`

type UpdateUserItemParams struct {
	ID       uuid.UUID
	IsOnSale bool
}

func (q *Queries) UpdateUserItem(ctx context.Context, arg UpdateUserItemParams) error {
	_, err := q.db.Exec(ctx, updateUserItem, arg.ID, arg.IsOnSale)
	return err
}
