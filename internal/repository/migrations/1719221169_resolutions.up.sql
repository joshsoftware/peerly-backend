CREATE TABLE IF NOT EXISTS resolutions (
    id SERIAL PRIMARY KEY,
    appreciation_id INT NOT NULL REFERENCES appreciations(id),
    reporting_action INT NOT NULL,
    reporting_comment VARCHAR NOT NULL,
    reported_by BIGINT NOT NULL REFERENCES users(id),
    reported_at BIGINT DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    moderator_action INT,
    moderator_comment VARCHAR(45),
    moderated_by BIGINT NOT NULL REFERENCES users(id),
    moderated_at BIGINT 
);

ALTER TABLE organizations ADD CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE organizations ADD CONSTRAINT fk_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);