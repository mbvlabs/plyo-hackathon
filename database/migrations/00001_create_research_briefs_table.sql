-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE researchbriefs (
    id TEXT PRIMARY KEY,
    identification_status TEXT NOT NULL,
    company_name TEXT NOT NULL,
    official_domain TEXT NOT NULL,
    headquarters TEXT NOT NULL,
    industry TEXT NOT NULL,
    company_type TEXT NOT NULL,
    status TEXT NOT NULL,
    geographic_scope TEXT NOT NULL,
    research_depth TEXT NOT NULL,
    confidence_score REAL NOT NULL CHECK (confidence_score >= 0.0 AND confidence_score <= 1.0),
    last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXIST researchbriefs; 
-- +goose StatementEnd
