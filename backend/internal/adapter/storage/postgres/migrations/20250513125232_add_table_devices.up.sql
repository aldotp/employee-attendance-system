CREATE TABLE devices (
    id UUID PRIMARY KEY,
    user_id UUID,
    device_id VARCHAR(255) UNIQUE NOT NULL,
    device_type VARCHAR(50),
    os_version VARCHAR(50),
    app_version VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES employees (id) ON DELETE CASCADE
);