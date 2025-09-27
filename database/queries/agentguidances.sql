-- name: QueryAgentGuidanceByID :one
select * from agentguidances where id=?;

-- name: QueryAgentGuidancesByResearchBriefID :many
select * from agentguidances where research_brief_id=?;

-- name: QueryAgentGuidances :many
select * from agentguidances;

-- name: QueryAllAgentGuidances :many
select * from agentguidances;

-- name: InsertAgentGuidance :one
insert into
    agentguidances (id, research_brief_id, guidance_key, guidance_value)
values
    (?, ?, ?, ?)
returning *;

-- name: UpdateAgentGuidance :one
update agentguidances
    set research_brief_id=?, guidance_key=?, guidance_value=?
where id = ?
returning *;

-- name: DeleteAgentGuidance :exec
delete from agentguidances where id=?;

-- name: QueryPaginatedAgentGuidances :many
select * from agentguidances 
order by created_at desc 
limit ? offset ?;

-- name: CountAgentGuidances :one
select count(*) from agentguidances;

