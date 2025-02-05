-- name: UpdateUserRed :one
update users
set
    is_red = true,
    updated_at = now()
where id = $1
returning *;
