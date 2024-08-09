CREATE TABLE IF NOT EXISTS notification_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    notification_token TEXT NOT NULL,
    UNIQUE (user_id, notification_token)
);