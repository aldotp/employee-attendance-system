CREATE TABLE attendances (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (
        type IN ('check_in', 'check_out')
    ),
    status VARCHAR(20) NOT NULL CHECK (
        status IN (
            'present',
            'absent',
            'late',
            'leave'
        )
    ),
    notes TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    selfie_url TEXT,
    hours_worked NUMERIC(5, 2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_attendances_user_id ON attendances (user_id);