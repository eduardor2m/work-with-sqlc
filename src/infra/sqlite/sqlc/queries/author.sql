-- name: GetAuthor :one
SELECT * FROM author
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM author
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO author (
    name, bio
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM author
WHERE id = $1;

-- name: DeleteAllAuthors :exec
DELETE FROM author;