-- +goose Up
alter table users
    add column is_red bool default false not null;

-- +goose Down
alter table users
    drop column is_red;
