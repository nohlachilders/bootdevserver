-- name: GetAllChirpsFromUser :many
select * from chirps
where user_id = $1
order by created_at asc;
