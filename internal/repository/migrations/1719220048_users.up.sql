CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
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
    created_by BIGINT REFERENCES users(id),
    updated_at BIGINT,
    updated_by BIGINT REFERENCES users(id)
);