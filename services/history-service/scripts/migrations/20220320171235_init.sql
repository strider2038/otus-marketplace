-- +goose Up
-- +goose StatementBegin
CREATE TYPE "deal_type" AS enum ('purchase', 'sale');

CREATE TABLE "deal"
(
    id           uuid primary key,
    user_id      uuid             not null,
    item_id      uuid             not null,
    item_name    text             not null,
    "type"       deal_type        not null,
    amount       double precision not null,
    commission   double precision not null,
    completed_at timestamp        not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "deal";
DROP TYPE "deal_type";
-- +goose StatementEnd
