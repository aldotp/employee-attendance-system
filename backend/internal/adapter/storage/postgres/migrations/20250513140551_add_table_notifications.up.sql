CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID,
    type VARCHAR(10) CHECK (
        type IN ('reminder', 'warning', 'info')
    ) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    send_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);