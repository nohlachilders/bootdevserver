-- +goose Up
create table chirps(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id uuid references users on delete cascade
);

-- +goose Down
drop table chirps;
