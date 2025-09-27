-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE agentguidances (
    id TEXT PRIMARY KEY,
    research_brief_id TEXT NOT NULL,
    guidance_key TEXT NOT NULL,
    guidance_value TEXT NOT NULL,
    FOREIGN KEY (research_brief_id) REFERENCES researchbriefs(id) ON DELETE CASCADE,
    UNIQUE (research_brief_id, guidance_key)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXIST agentguidance;
-- +goose StatementEnd
