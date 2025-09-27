-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE companycandidates (
    id TEXT PRIMARY KEY,
    research_brief_id TEXT NOT NULL,
    name TEXT NOT NULL,
    domain TEXT NOT NULL,
    description TEXT NOT NULL,
    industry TEXT NOT NULL,
    location TEXT NOT NULL,
    FOREIGN KEY (research_brief_id) REFERENCES researchbriefs(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXIST companycandidates;
-- +goose StatementEnd
