-- name: CreateSnippet :one
INSERT INTO snippets(id, created_at, updated_at, language_id, user_id, snippet_text)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3)
RETURNING *;

-- name: GetSnippetsByCreatedAt :many
SELECT snippets.id, snippets.created_at, snippets.updated_at, snippets.user_id, snippet_text, languages.name as language
FROM snippets
INNER JOIN languages ON snippets.language_id = languages.id
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: GetSnippetCount :one
SELECT COUNT(*) FROM snippets;
