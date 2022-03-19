CREATE TABLE "item"
(
    id                 uuid primary key,
    name               text             not null,
    initial_count      int              not null,
    initial_price      double precision not null,
    commission_percent double precision not null,
    created_at         timestamp        not null default now()
);

CREATE TABLE "user_item"
(
    id         uuid primary key,
    item_id    uuid      not null references item (id),
    user_id    uuid      not null,
    is_on_sale boolean   not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

CREATE TYPE "purchase_status" AS enum(
    'pending',
    'canceled',
    'paymentPending',
    'paymentSucceeded',
    'paymentFailed',
    'approved'
);

CREATE TYPE "sell_status" AS enum(
    'pending',
    'canceled',
    'dealPending',
    'accrualPending',
    'approved'
);

CREATE TABLE "sell_order"
(
    id           uuid primary key,
    user_id      uuid             not null,
    item_id      uuid             not null references item (id),
    user_item_id uuid             not null references user_item (id),
    accrual_id   uuid unique,
    deal_id      uuid unique,
    price        double precision not null,
    commission   double precision not null,
    status       sell_status      not null default 'pending',
    created_at   timestamp        not null default now(),
    updated_at   timestamp        not null default now()
);

CREATE TABLE "purchase_order"
(
    id         uuid primary key,
    user_id    uuid             not null,
    item_id    uuid             not null references item (id),
    payment_id uuid unique,
    deal_id    uuid unique,
    price      double precision not null,
    commission double precision not null,
    status     purchase_status  not null default 'pending',
    created_at timestamp        not null default now(),
    updated_at timestamp        not null default now()
);
