CREATE TABLE IF NOT EXISTS transactions_tools (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    tool_id BIGINT NOT NULL REFERENCES tools(id) ON DELETE RESTRICT,
    UNIQUE(transaction_id, tool_id)
);