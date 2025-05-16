CREATE TABLE leave_requests (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (
        type IN (
            'annual',
            'sick',
            'unpaid',
            'maternity',
            'paternity'
        )
    ),
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'approved',
            'rejected'
        )
    ),
    reviewed_by VARCHAR(255),
    reviewed_at TIMESTAMPTZ,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_leave_requests_users_id ON leave_requests (user_id);

CREATE INDEX idx_leave_requests_status ON leave_requests (status);