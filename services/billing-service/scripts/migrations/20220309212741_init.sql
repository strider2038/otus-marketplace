-- +goose Up
-- +goose StatementBegin
CREATE TABLE "account"
(
    id         uuid primary key,
    amount     float     not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

CREATE TYPE "operation_type" AS enum('deposit', 'withdraw', 'payment', 'accrual');

CREATE TABLE "operation"
(
    id          uuid primary key,
    account_id  uuid           not null references account (id),
    "type"      operation_type not null,
    amount      float          not null,
    description text           not null,
    created_at  timestamp      not null default now()
);

INSERT INTO "account" (id) VALUES ('2a104e66-1c78-4577-ab15-3ae935180c17');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "operation";
DROP TABLE "account";
DROP TYPE "operation_type";
-- +goose StatementEnd
