// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: snippets.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createSnippet = `-- name: CreateSnippet :one
INSERT INTO snippets(id, created_at, updated_at, language_id, user_id, snippet_text)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3)
RETURNING id, created_at, updated_at, language_id, user_id, snippet_text
`

type CreateSnippetParams struct {
	LanguageID  uuid.UUID
	UserID      uuid.UUID
	SnippetText string
}

func (q *Queries) CreateSnippet(ctx context.Context, arg CreateSnippetParams) (Snippet, error) {
	row := q.db.QueryRowContext(ctx, createSnippet, arg.LanguageID, arg.UserID, arg.SnippetText)
	var i Snippet
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LanguageID,
		&i.UserID,
		&i.SnippetText,
	)
	return i, err
}

const getSnippetCount = `-- name: GetSnippetCount :one
SELECT COUNT(*)
FROM snippets
`

func (q *Queries) GetSnippetCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSnippetCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getSnippetsByCreatedAt = `-- name: GetSnippetsByCreatedAt :many
SELECT COUNT(*) OVER () AS total_count,
 snippets.id, snippets.created_at, snippets.updated_at, snippets.user_id, snippet_text, languages.name AS language
FROM snippets
INNER JOIN languages ON snippets.language_id = languages.id
WHERE (languages.name = $3 OR $3 IS NULL)
ORDER BY created_at
LIMIT $1 OFFSET $2
`

type GetSnippetsByCreatedAtParams struct {
	Limit    int32
	Offset   int32
	Language sql.NullString
}

type GetSnippetsByCreatedAtRow struct {
	TotalCount  int64
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID
	SnippetText string
	Language    string
}

func (q *Queries) GetSnippetsByCreatedAt(ctx context.Context, arg GetSnippetsByCreatedAtParams) ([]GetSnippetsByCreatedAtRow, error) {
	rows, err := q.db.QueryContext(ctx, getSnippetsByCreatedAt, arg.Limit, arg.Offset, arg.Language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSnippetsByCreatedAtRow
	for rows.Next() {
		var i GetSnippetsByCreatedAtRow
		if err := rows.Scan(
			&i.TotalCount,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.SnippetText,
			&i.Language,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
