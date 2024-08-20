alter table badges
add updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT;

alter table badges
add updated_by BIGINT;
