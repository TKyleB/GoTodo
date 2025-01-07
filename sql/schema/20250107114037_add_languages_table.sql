-- +goose Up
-- +goose StatementBegin
CREATE TABLE languages(
    id uuid PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), 'Python');

INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), 'Javascript');

INSERT INTO languages(id, name)
VALUES(gen_random_uuid(), 'Go');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE languages;
-- +goose StatementEnd
