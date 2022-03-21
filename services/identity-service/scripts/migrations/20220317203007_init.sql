-- +goose Up
-- +goose StatementBegin
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
);

INSERT INTO "user" (id, email, password, role, first_name, last_name)
VALUES ('2a104e66-1c78-4577-ab15-3ae935180c17',
        'broker@marketplace.com',
        '$argon2id$v=19$m=65536,t=1,p=4$1UwsYpQM2UTHa+cIWdEC1A$qYHFJja4b/JdOzfqDPCyj74r5hQo1yYozLA6khwRkPQ',
        'broker',
        'Broker',
        'Account');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
DROP TYPE "user_role";
-- +goose StatementEnd
