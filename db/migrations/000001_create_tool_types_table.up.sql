CREATE TABLE IF NOT EXISTS tool_types (
    id BIGSERIAL PRIMARY KEY,
    part_number VARCHAR(32) NOT NULL,
    description TEXT,
    co VARCHAR(5) NOT NULL,
    mc VARCHAR(10) NOT NULL
);