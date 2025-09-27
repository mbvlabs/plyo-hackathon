-- name: QuerySourcesByID :one
select * from sources where id=?;

-- name: QuerySourcesByResearchBriefID :many
select * from sources where research_brief_id=?;

-- name: QuerySourcess :many
select * from sources;

-- name: QueryAllSourcess :many
select * from sources;

-- name: InsertSources :one
insert into
    sources (id, research_brief_id, source_url)
values
    (?, ?, ?)
returning *;

-- name: UpdateSources :one
update sources
    set research_brief_id=?, source_url=?
where id = ?
returning *;

-- name: DeleteSources :exec
delete from sources where id=?;

-- name: QueryPaginatedSourcess :many
select * from sources 
order by created_at desc 
limit ? offset ?;

-- name: CountSourcess :one
select count(*) from sources;

