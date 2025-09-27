-- name: QuerySpecialConsiderationsByID :one
select * from specialconsiderations where id=?;

-- name: QuerySpecialConsiderationsByResearchBriefID :many
select * from specialconsiderations where research_brief_id=?;

-- name: QuerySpecialConsiderationss :many
select * from specialconsiderations;

-- name: QueryAllSpecialConsiderationss :many
select * from specialconsiderations;

-- name: InsertSpecialConsiderations :one
insert into
    specialconsiderations (id, research_brief_id, consideration)
values
    (?, ?, ?)
returning *;

-- name: UpdateSpecialConsiderations :one
update specialconsiderations
    set research_brief_id=?, consideration=?
where id = ?
returning *;

-- name: DeleteSpecialConsiderations :exec
delete from specialconsiderations where id=?;

-- name: QueryPaginatedSpecialConsiderationss :many
select * from specialconsiderations 
order by created_at desc 
limit ? offset ?;

-- name: CountSpecialConsiderationss :one
select count(*) from specialconsiderations;

