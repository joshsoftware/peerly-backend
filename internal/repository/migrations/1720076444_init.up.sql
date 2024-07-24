CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(25) NOT NULL
);

CREATE TABLE IF NOT EXISTS grades (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    points INT
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    employee_id varchar(225) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(255) , 
    profile_image_url TEXT,
    designation VARCHAR(255) NOT NULL, 
    reward_quota_balance INT NOT NULL,
    status INT DEFAULT 1,
    role_id INT NOT NULL REFERENCES roles(id),
    grade_id INT NOT NULL REFERENCES grades(id),
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    created_by BIGINT ,
    updated_at BIGINT,
    updated_by BIGINT 
);

CREATE TABLE IF NOT EXISTS user_blacklisted_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token TEXT NOT NULL UNIQUE,
    expires_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS core_values (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    parent_core_value_id INT REFERENCES core_values(id)
);

CREATE TABLE IF NOT EXISTS badges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(45) NOT NULL,
    reward_points INT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_badges (
    id SERIAL PRIMARY KEY,
    badge_id INT NOT NULL REFERENCES badges(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT
);

CREATE TABLE IF NOT EXISTS appreciations (
    id SERIAL PRIMARY KEY,
    core_value_id INT NOT NULL REFERENCES core_values(id),
    description TEXT NOT NULL, -- reason
    is_valid BOOLEAN NOT NULL DEFAULT true,
    total_reward_points INT DEFAULT 0,
    quarter INT NOT NULL,
    sender BIGINT NOT NULL REFERENCES users(id),
    receiver BIGINT NOT NULL REFERENCES users(id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT
);

CREATE TABLE IF NOT EXISTS rewards (
    id SERIAL PRIMARY KEY,
    appreciation_id INT NOT NULL REFERENCES appreciations(id),
    point INT NOT NULL,
    sender BIGINT NOT NULL REFERENCES users(id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT
);

CREATE TABLE IF NOT EXISTS resolutions (
    id SERIAL PRIMARY KEY,
    appreciation_id INT NOT NULL REFERENCES appreciations(id),
    reporting_comment VARCHAR NOT NULL,
    reported_by BIGINT NOT NULL REFERENCES users(id),
    reported_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    moderator_action INT,
    moderator_comment VARCHAR(45),
    moderated_by BIGINT  REFERENCES users(id),
    moderated_at BIGINT 
);

CREATE TABLE IF NOT EXISTS organization_config (
    id BIGSERIAL PRIMARY KEY,
    reward_multiplier INT,
    reward_quota_renewal_frequency INT, -- Assuming month is only a unit for this integer value
    timezone VARCHAR(100) DEFAULT 'UTC',
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    created_by BIGINT REFERENCES users(id),
    updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    updated_by BIGINT REFERENCES users(id)
);