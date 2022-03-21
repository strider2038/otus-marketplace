-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    id         uuid primary key,
    email      text      not null unique,
    first_name text      not null,
    last_name  text      not null,
    phone      text      not null default '',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

CREATE TABLE "notification"
(
    id         uuid primary key,
    user_id    uuid      not null references "user" (id),
    message    text      not null,
    created_at timestamp not null default now()
);

INSERT INTO "user" (id, email, first_name, last_name)
VALUES ('2a104e66-1c78-4577-ab15-3ae935180c17',
        'broker@marketplace.com',
        'Broker',
        'Account');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
DROP TABLE "notification";
-- +goose StatementEnd
