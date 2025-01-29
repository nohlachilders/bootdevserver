-- +goose Up
create table refresh_tokens(
    token text primary key, 
    created_at timestamp not null,
    updated_at timestamp not null,
    expires_at timestamp not null,
    revoked_at timestamp,
    user_id uuid not null references users on delete cascade
);

-- +goose Down
drop table refresh_tokens;
