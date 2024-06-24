CREATE TABLE IF NOT EXISTS core_values (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    parent_core_value_id INT REFERENCES core_values(id)
);