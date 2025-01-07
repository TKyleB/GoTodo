-- name: CreateLanguage :one
INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), $1)
returning *;