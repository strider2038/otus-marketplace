-- +goose Up
-- +goose StatementBegin
CREATE TABLE "daily_deals"
(
    "date"    date primary key,
    item_id   uuid             not null,
    item_name text             not null,
    count     bigint           not null,
    amount    double precision not null
);

CREATE TABLE "total_daily_deals"
(
    "date" date primary key,
    count  bigint           not null,
    amount double precision not null
);

CREATE TABLE "top10_deals"
(
    item_id   uuid             not null primary key,
    item_name text             not null,
    count     bigint           not null,
    amount    double precision not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "top10_deals";
DROP TABLE "total_daily_deals";
DROP TABLE "daily_deals";
-- +goose StatementEnd
