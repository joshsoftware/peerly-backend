alter table grades
add updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT;

alter table grades
add updated_by BIGINT;

