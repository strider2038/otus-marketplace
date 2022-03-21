CREATE TYPE "user_role" AS enum (
    'trader',
    'broker',
    'admin'
);

CREATE TABLE "user"
(
    id         uuid primary key,
    email      text      not null unique,
    password   text      not null,
    role       user_role not null default 'trader',
    first_name text      not null,
    last_name  text      not null,
    phone      text      not null default '',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
)
