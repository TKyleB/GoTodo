-- +goose Up
-- +goose StatementBegin
INSERT INTO languages (id, name) 
VALUES (gen_random_uuid (), 'text');

INSERT INTO
    languages (id, name)
VALUES (gen_random_uuid (), 'java');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM languages WHERE name='text';
DELETE FROM languages WHERE name='java';
-- +goose StatementEnd
