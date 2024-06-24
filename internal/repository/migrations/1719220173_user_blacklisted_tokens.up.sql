CREATE TABLE IF NOT EXISTS user_blacklisted_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token TEXT NOT NULL UNIQUE,
    expires_at BIGINT NOT NULL
);