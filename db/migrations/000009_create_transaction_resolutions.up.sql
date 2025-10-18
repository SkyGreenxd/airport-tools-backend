CREATE TABLE IF NOT EXISTS transaction_resolutions (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id),
    qa_employee_id BIGINT NOT NULL REFERENCES users(id),
    reason VARCHAR(32) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);