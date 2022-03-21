CREATE TABLE "account"
(
    id         uuid primary key,
    amount     float     not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

CREATE TYPE "operation_type" AS enum ('deposit', 'withdraw', 'payment', 'accrual', 'commission');

CREATE TABLE "operation"
(
    id          uuid primary key,
    account_id  uuid           not null references account (id),
    "type"      operation_type not null,
    amount      float          not null,
    description text           not null,
    created_at  timestamp      not null default now()
);
