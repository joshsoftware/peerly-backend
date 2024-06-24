CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    contact_email VARCHAR(255) NOT NULL,
    is_email_verified BOOLEAN DEFAULT FALSE,
    domain_name VARCHAR(255) NOT NULL,
    status INT DEFAULT 1,
    subscription_status INT DEFAULT 1,
    subscription_valid_upto BIGINT,
    reward_multiplier INT NOT NULL,
    reward_quota_renewal_frequency INT NOT NULL,
    timezone VARCHAR(100) DEFAULT 'UTC',
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    created_by BIGINT ,
    updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    updated_by BIGINT 
);