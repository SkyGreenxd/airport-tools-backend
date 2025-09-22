CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    location_id BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'OPEN',
    checkout_scan_id BIGINT REFERENCES cv_scans(id),
    checkin_scan_id BIGINT REFERENCES cv_scans(id),
    reason TEXT
);