-- name: QueryResearchBriefByID :one
select * from researchbriefs where id=?;

-- name: QueryResearchBriefs :many
select * from researchbriefs;

-- name: QueryAllResearchBriefs :many
select * from researchbriefs;

-- name: InsertResearchBrief :one
insert into
    researchbriefs (id, identification_status, company_name, official_domain, headquarters, industry, company_type, status, geographic_scope, research_depth, confidence_score, last_updated)
values
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
returning *;

-- name: UpdateResearchBrief :one
update researchbriefs
    set identification_status=?, company_name=?, official_domain=?, headquarters=?, industry=?, company_type=?, status=?, geographic_scope=?, research_depth=?, confidence_score=?, last_updated=?
where id = ?
returning *;

-- name: DeleteResearchBrief :exec
delete from researchbriefs where id=?;

-- name: QueryPaginatedResearchBriefs :many
select * from researchbriefs 
order by created_at desc 
limit ? offset ?;

-- name: CountResearchBriefs :one
select count(*) from researchbriefs;

