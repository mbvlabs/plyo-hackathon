-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE reports (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    compay_candidate_id TEXT NOT NULL,
    company_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    progress_percentage INTEGER DEFAULT 0,

    preliminary_research_completed BOOLEAN DEFAULT FALSE,
    company_intelligence_completed BOOLEAN DEFAULT FALSE,
    competitive_intelligence_completed BOOLEAN DEFAULT FALSE,
    market_dynamics_completed BOOLEAN DEFAULT FALSE,
    trend_analysis_completed BOOLEAN DEFAULT FALSE,

    company_intelligence_data TEXT,
    competitive_intelligence_data TEXT,
    market_dynamics_data TEXT,
    trend_analysis_data TEXT,

    final_report TEXT,

    completed_at DATETIME,

    FOREIGN KEY (compay_candidate_id) REFERENCES companycandidates(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXIST reports; 
-- +goose StatementEnd
