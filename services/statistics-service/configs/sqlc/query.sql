-- name: FindDailyDeals :many
SELECT "date", item_id, item_name, count, amount
FROM "daily_deals"
WHERE "date" >= now() - interval '1 week'
ORDER BY "date", item_name DESC;

-- name: AddDailyDeals :exec
INSERT INTO "daily_deals" AS d ("date", item_id, item_name, count, amount)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ("date")
DO UPDATE SET count = d.count + EXCLUDED.count, amount = d.amount + EXCLUDED.amount;

-- name: FindTotalDailyDeals :many
SELECT "date", count, amount
FROM "total_daily_deals"
WHERE "date" >= now() - interval '1 week'
ORDER BY "date" DESC;

-- name: AddTotalDailyDeals :exec
INSERT INTO "total_daily_deals" AS d ("date", count, amount)
VALUES ($1, $2, $3)
ON CONFLICT ("date")
DO UPDATE SET count = d.count + EXCLUDED.count, amount = d.amount + EXCLUDED.amount;

-- name: FindTop10Deals :many
SELECT item_id, item_name, count, amount
FROM "top10_deals"
ORDER BY amount DESC
LIMIT 10;

-- name: AddTop10Deals :exec
INSERT INTO "top10_deals" AS d (item_id, item_name, count, amount)
VALUES ($1, $2, $3, $4)
ON CONFLICT (item_id)
DO UPDATE SET count = d.count + EXCLUDED.count, amount = d.amount + EXCLUDED.amount;
