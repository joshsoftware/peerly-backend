CREATE TABLE IF NOT EXISTS appreciations (
    id SERIAL PRIMARY KEY,
    core_value_id INT NOT NULL REFERENCES core_values(id),
    description TEXT NOT NULL, -- reason
    is_valid BOOLEAN NOT NULL DEFAULT true,
    total_rewards INT DEFAULT 0,
    quarter INT NOT NULL,
    sender BIGINT NOT NULL REFERENCES users(id),
    receiver BIGINT NOT NULL REFERENCES users(id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT
);
