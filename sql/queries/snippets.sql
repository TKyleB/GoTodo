-- name: CreateSnippet :one
INSERT INTO snippets(id, created_at, updated_at, language_id, author_id, snippet_text)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3)
RETURNING *;