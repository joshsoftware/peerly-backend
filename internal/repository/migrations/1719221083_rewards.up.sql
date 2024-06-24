CREATE TABLE IF NOT EXISTS rewards (
    id SERIAL PRIMARY KEY,
    appreciation_id INT NOT NULL REFERENCES appreciations(id),
    point INT NOT NULL,
    sender BIGINT NOT NULL REFERENCES users(id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT
);
