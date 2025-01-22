-- +goose Up
-- +goose StatementBegin
ALTER TABLE snippets ADD COLUMN search_vector tsvector;
UPDATE snippets SET search_vector = to_tsvector('simple', snippet_title || ' ' || snippet_description || ' ' || snippet_text);
CREATE OR REPLACE FUNCTION update_search_vector() RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('simple', NEW.snippet_title || ' ' || NEW.snippet_description || ' ' || NEW.snippet_text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER tsvector_update
BEFORE INSERT OR UPDATE ON snippets
FOR EACH ROW EXECUTE FUNCTION update_search_vector();
CREATE INDEX idx_snippets_search ON snippets USING gin(search_vector);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER tsvector_update ON snippets;
DROP FUNCTION update_search_vector();
DROP INDEX idx_snippets_search;
ALTER TABLE snippets DROP COLUMN search_vector;
-- +goose StatementEnd
