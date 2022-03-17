-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    id         uuid primary key,
    email      text not null unique,
    password   text not null,
    first_name text not null,
    last_name  text not null,
    phone      text not null default '',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user"
-- +goose StatementEnd
