-- name: QueryReportByID :one
select * from reports where id=?;

-- name: QueryReports :many
select * from reports;

-- name: QueryAllReports :many
select * from reports;

-- name: InsertReport :one
insert into
    reports (id, created_at, updated_at, compay_candidate_id, company_name, status, progress_percentage, preliminary_research_completed, company_intelligence_completed, competitive_intelligence_completed, market_dynamics_completed, trend_analysis_completed, company_intelligence_data, competitive_intelligence_data, market_dynamics_data, trend_analysis_data, final_report, completed_at)
values
    (?, datetime('now'), datetime('now'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
returning *;

-- name: UpdateReport :one
update reports
    set updated_at=datetime('now'), compay_candidate_id=?, company_name=?, status=?, progress_percentage=?, preliminary_research_completed=?, company_intelligence_completed=?, competitive_intelligence_completed=?, market_dynamics_completed=?, trend_analysis_completed=?, company_intelligence_data=?, competitive_intelligence_data=?, market_dynamics_data=?, trend_analysis_data=?, final_report=?, completed_at=?
where id = ?
returning *;

-- name: DeleteReport :exec
delete from reports where id=?;

-- name: QueryPaginatedReports :many
select * from reports 
order by created_at desc 
limit ? offset ?;

-- name: CountReports :one
select count(*) from reports;

-- name: UpdateCompanyIntelligence :exec
UPDATE reports
SET company_intelligence_data = ?,
    company_intelligence_completed = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateCompetitiveIntelligence :exec
UPDATE reports
SET competitive_intelligence_data = ?,
    competitive_intelligence_completed = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateMarketDynamics :exec
UPDATE reports
SET market_dynamics_data = ?,
    market_dynamics_completed = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateTrendAnalysis :exec
UPDATE reports
SET trend_analysis_data = ?,
    trend_analysis_completed = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateReportProgress :exec
UPDATE reports
SET progress_percentage = ?,
    status = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateFinalReport :exec
UPDATE reports
SET final_report = ?,
    status = 'completed',
    completed_at = datetime('now'),
    updated_at = datetime('now')
WHERE id = ?;

