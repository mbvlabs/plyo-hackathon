-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE sources (
    id TEXT PRIMARY KEY,
    research_brief_id TEXT NOT NULL,
    source_url TEXT NOT NULL,
    FOREIGN KEY (research_brief_id) REFERENCES researchbriefs(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXIST sources;
-- +goose StatementEnd
