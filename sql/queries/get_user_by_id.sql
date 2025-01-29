-- name: GetUserByID :one
select * from users
where id = $1;
