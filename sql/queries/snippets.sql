-- name: CreateSnippet :one
WITH inserted_snippet AS (
INSERT INTO snippets(id, created_at, updated_at, language_id, user_id, snippet_text, snippet_description, snippet_title)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5)
RETURNING *
)
SELECT inserted_snippet.*, users.username
FROM inserted_snippet
INNER JOIN users ON users.id = inserted_snippet.user_id;

-- name: GetSnippetsByCreatedAt :many
SELECT COUNT(*) OVER () AS total_count,
 snippets.id, snippets.created_at, snippets.updated_at, snippets.user_id, snippet_text, users.username, snippet_description, snippet_title, languages.name AS language
FROM snippets
INNER JOIN languages ON snippets.language_id = languages.id
INNER JOIN users ON snippets.user_id = users.id
WHERE
(languages.name = sqlc.narg('language') OR sqlc.narg('language') IS NULL)
AND (users.username = sqlc.narg('username') OR sqlc.narg('username') IS NULL)
AND (search_vector @@ to_tsquery('simple', sqlc.narg('search')) OR sqlc.narg('search') IS NULL)
ORDER BY snippets.created_at DESC
LIMIT $1 OFFSET $2;


-- name: GetSnippetById :one
SELECT snippets.id, snippets.created_at, snippets.updated_at, snippets.user_id, snippet_text, users.username, snippet_description, snippet_title, languages.name AS language
FROM snippets
INNER JOIN languages ON snippets.language_id = languages.id
INNER JOIN users ON snippets.user_id = users.id
WHERE snippets.id = $1;

-- name: GetSnippetBySearch :many
SELECT snippets.id, snippets.created_at, snippets.updated_at, snippets.user_id, snippet_text, users.username, snippet_description, snippet_title, languages.name AS language
FROM snippets
INNER JOIN languages ON snippets.language_id = languages.id
INNER JOIN users ON snippets.user_id = users.id
WHERE search_vector @@ to_tsquery('simple', $1);


