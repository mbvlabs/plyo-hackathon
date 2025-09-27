-- name: QueryCompanyCandidatesByID :one
select * from companycandidates where id=?;

-- name: QueryCompanyCandidatesByResearchBriefID :many
select * from companycandidates where research_brief_id=?;

-- name: QueryCompanyCandidatess :many
select * from companycandidates;

-- name: QueryAllCompanyCandidatess :many
select * from companycandidates;

-- name: InsertCompanyCandidates :one
insert into
    companycandidates (id, research_brief_id, name, domain, description, industry, location)
values
    (?, ?, ?, ?, ?, ?, ?)
returning *;

-- name: UpdateCompanyCandidates :one
update companycandidates
    set research_brief_id=?, name=?, domain=?, description=?, industry=?, location=?
where id = ?
returning *;

-- name: DeleteCompanyCandidates :exec
delete from companycandidates where id=?;

-- name: QueryPaginatedCompanyCandidatess :many
select * from companycandidates 
order by created_at desc 
limit ? offset ?;

-- name: CountCompanyCandidatess :one
select count(*) from companycandidates;

