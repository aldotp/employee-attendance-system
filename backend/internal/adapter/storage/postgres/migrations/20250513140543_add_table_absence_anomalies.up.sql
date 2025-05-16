CREATE TABLE absence_anomalies (
    id UUID PRIMARY KEY,
    user_id UUID,
    date DATE NOT NULL,
    type VARCHAR(20) CHECK (
        type IN (
            'late',
            'not_present',
            'left_early',
            'forgot_checkin'
        )
    ) NOT NULL,
    note TEXT,
    verified BOOLEAN DEFAULT false,
    verified_by VARCHAR(255),
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);