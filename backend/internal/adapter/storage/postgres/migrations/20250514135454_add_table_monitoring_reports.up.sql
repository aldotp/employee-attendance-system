CREATE TABLE monitoring_reports (
    id UUID,
    report_type VARCHAR(255) NOT NULL,
    data TEXT NOT NULL,
    generated_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (id)
);

CREATE INDEX idx_monitoring_reports_generated_at ON monitoring_reports (generated_at);