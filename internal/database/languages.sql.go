// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: languages.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createLanguage = `-- name: CreateLanguage :one
INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), $1)
returning id, name
`

func (q *Queries) CreateLanguage(ctx context.Context, name string) (Language, error) {
	row := q.db.QueryRowContext(ctx, createLanguage, name)
	var i Language
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getLanguageByName = `-- name: GetLanguageByName :one
SELECT id FROM languages
WHERE name = $1
`

func (q *Queries) GetLanguageByName(ctx context.Context, name string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getLanguageByName, name)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
