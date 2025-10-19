CREATE TABLE IF NOT EXISTS model_err_items (
    resolution_id BIGINT NOT NULL REFERENCES transaction_resolutions(id) ON DELETE CASCADE,
    tool_type_id BIGINT NOT NULL REFERENCES tool_types(id) ON DELETE CASCADE,
    PRIMARY KEY (resolution_id, tool_type_id)
);