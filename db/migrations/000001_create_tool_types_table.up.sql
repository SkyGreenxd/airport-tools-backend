CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS tool_types (
    id BIGSERIAL PRIMARY KEY,
    part_number VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(256),
    reference_image_hash VARCHAR(256) UNIQUE NOT NULL,
    reference_embedding VECTOR(1280) NOT NULL
);