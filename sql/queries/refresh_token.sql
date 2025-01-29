-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, expires_at, user_id)
values (
    $1,
    now(),
    now(),
    $2,
    $3
    )
returning *;
