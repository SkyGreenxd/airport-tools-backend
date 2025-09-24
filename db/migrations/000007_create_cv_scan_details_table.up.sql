CREATE TABLE IF NOT EXISTS cv_scan_details (
    id BIGSERIAL PRIMARY KEY,
    cv_scan_id BIGINT NOT NULL REFERENCES cv_scans(id) ON DELETE CASCADE,
    detected_tool_type_id BIGINT NOT NULL REFERENCES tool_types(id),
    image_hash VARCHAR(256),
    embedding VECTOR(512)
);