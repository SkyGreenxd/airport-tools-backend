CREATE TABLE IF NOT EXISTS cv_scans (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    scan_type VARCHAR(16) NOT NULL,
    image_url VARCHAR(256) NOT NULL
);