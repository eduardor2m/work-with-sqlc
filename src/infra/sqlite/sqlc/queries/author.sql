-- name: GetAuthor :one
SELECT * FROM author WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM author ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO author (name, bio) VALUES (?, ?) RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM author WHERE id = ?;

-- name: DeleteAllAuthors :exec
DELETE FROM author;