CREATE TABLE IF NOT EXISTS tools (
    id BIGSERIAL PRIMARY KEY,
    type_tool_id BIGINT NOT NULL REFERENCES tool_types(id),
    toir_id BIGINT UNIQUE NOT NULL,
    location_id BIGINT NOT NULL REFERENCES locations(id),
    sn_bn VARCHAR(32) NOT NULL,
    expires_at DATE
);