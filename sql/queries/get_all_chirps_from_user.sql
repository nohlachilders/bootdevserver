-- name: GetAllChirpsFromUser :many
select * from chirps
where id = $1
order by created_at asc;
