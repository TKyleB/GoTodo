-- +goose Up
-- +goose StatementBegin
CREATE TABLE languages(
    id uuid PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), 'python');

INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), 'javascript');

INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), 'go');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE languages;
-- +goose StatementEnd
