-- +goose Up
-- +goose StatementBegin
CREATE TABLE snippets(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    language_id uuid NOT NULL,
    author_id uuid NOT NULL,
    snippet_text TEXT NOT NULL,
    FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE snippets;
-- +goose StatementEnd
