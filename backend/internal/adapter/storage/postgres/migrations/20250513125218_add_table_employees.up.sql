CREATE TABLE employees (
    id UUID PRIMARY KEY,
    user_id UUID,
    department_id UUID,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    photo_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (
        status IN (
            'active',
            'inactive',
            'terminated',
            'suspended'
        )
    ),
    join_date DATE,
    reporting_to UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);