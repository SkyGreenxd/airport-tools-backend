CREATE TABLE IF NOT EXISTS tool_set_items (
    id BIGSERIAL PRIMARY KEY,
    tool_set_id BIGINT NOT NULL REFERENCES tool_sets(id) ON DELETE CASCADE,
    tool_type_id BIGINT NOT NULL REFERENCES tool_types(id) ON DELETE RESTRICT,
    UNIQUE (tool_set_id, tool_type_id)
);