CREATE TABLE IF NOT EXISTS cv_scans (
    id BIGSERIAL PRIMARY KEY,
    scan_type VARCHAR(16) NOT NULL, -- 'checkout' или 'checkin'
    image_url VARCHAR(256) NOT NULL,
    image_hash VARCHAR(256) NOT NULL
);