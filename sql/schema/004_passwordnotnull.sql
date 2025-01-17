
-- +goose Up
alter table users
    alter column hashed_password set not null;

-- +goose Down
alter table users
    alter column hashed_password drop not null;
