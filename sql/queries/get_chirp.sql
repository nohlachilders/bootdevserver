-- name: GetChirp :one
select * from chirps
where id = $1;
