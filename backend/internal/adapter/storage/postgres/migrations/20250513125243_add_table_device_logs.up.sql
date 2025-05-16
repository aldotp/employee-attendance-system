CREATE TABLE device_logs (
    id UUID PRIMARY KEY,
    device_id UUID,
    ip_address VARCHAR(45),
    user_agent TEXT,
    event VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES devices (id) ON DELETE CASCADE
);