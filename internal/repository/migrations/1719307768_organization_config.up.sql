CREATE TABLE organization_config (
    id BIGSERIAL PRIMARY KEY,
    reward_multiplier INT,
    reward_quota_renewal_frequency INT, -- Assuming month is just a unit for this integer value
    timezone VARCHAR(100) DEFAULT 'UTC',
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    created_by BIGINT REFERENCES users(id),
    updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    updated_by BIGINT REFERENCES users(id)
);