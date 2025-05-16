CREATE TABLE schedules (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    date DATE NOT NULL,
    shift_start TIME NOT NULL,
    shift_end TIME NOT NULL,
    break_start TIME,
    break_end TIME,
    work_location_id UUID,
    schedule_type VARCHAR(20) NOT NULL DEFAULT 'regular' CHECK (
        schedule_type IN (
            'regular',
            'shift',
            'flexible'
        )
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (work_location_id) REFERENCES work_locations (id) ON DELETE SET NULL
);

CREATE INDEX idx_schedules_users_id ON schedules (user_id);

CREATE INDEX idx_schedules_date ON schedules (date);