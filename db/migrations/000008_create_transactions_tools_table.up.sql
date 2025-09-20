CREATE TABLE IF NOT EXISTS transactions_tools (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id),
    tool_id BIGINT NOT NULL REFERENCES tools(id)
);