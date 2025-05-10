CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role ENUM('admin', 'hr', 'manager', 'employee') NOT NULL DEFAULT 'employee',
    location VARCHAR(255),
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    photo_url TEXT,
    status ENUM('active', 'inactive', 'terminated', 'suspended') NOT NULL DEFAULT 'active',
    email_verified_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE TABLE departments (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    head_employee_id CHAR(36),
    location VARCHAR(255),
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    wfa_policy JSON,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE employees (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36),
    department_id CHAR(36),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(100),
    location VARCHAR(255),
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    photo_url TEXT,
    status ENUM('active', 'inactive', 'terminated', 'suspended') NOT NULL DEFAULT 'active',
    join_date DATE,
    reporting_to CHAR(36),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
    FOREIGN KEY (reporting_to) REFERENCES employees(id) ON DELETE SET NULL
);

CREATE TABLE work_locations (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    country_code VARCHAR(10),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    radius INTEGER,
    status ENUM('active', 'inactive') NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE employee_work_locations (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36),
    work_location_id CHAR(36),
    is_primary BOOLEAN DEFAULT false,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (work_location_id) REFERENCES work_locations(id) ON DELETE CASCADE
);

CREATE TABLE attendances (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36),
    date DATE NOT NULL,
    check_in_time DATETIME,
    check_out_time DATETIME,
    type VARCHAR(50),
    location VARCHAR(255),
    work_location_id CHAR(36),
    selfie_url TEXT,
    status ENUM('present', 'absent', 'late', 'early_leave') NOT NULL DEFAULT 'present',
    device_info JSON,
    ip_address VARCHAR(45),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_attendance (employee_id, date),
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (work_location_id) REFERENCES work_locations(id) ON DELETE SET NULL
);

CREATE TABLE schedules (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36),
    date DATE NOT NULL,
    shift_start TIME NOT NULL,
    shift_end TIME NOT NULL,
    break_start TIME,
    break_end TIME,
    work_location_id CHAR(36),
    schedule_type ENUM('regular', 'shift', 'flexible') NOT NULL DEFAULT 'regular',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (work_location_id) REFERENCES work_locations(id) ON DELETE SET NULL
);

CREATE TABLE leave_requests (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    type ENUM('annual', 'sick', 'unpaid', 'maternity', 'paternity') NOT NULL,
    reason TEXT,
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
    reviewed_by VARCHAR(255),
    reviewed_at TIMESTAMP,
    note TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE absence_anomalies (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36),
    date DATE NOT NULL,
    type ENUM('late', 'not_present', 'left_early', 'forgot_checkin') NOT NULL,
    note TEXT,
    verified BOOLEAN DEFAULT false,
    verified_by VARCHAR(255),
    verified_at TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE notifications (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36),
    type ENUM('reminder', 'warning', 'info') NOT NULL,
    message TEXT NOT NULL,
    send_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE holidays (
    id CHAR(36) PRIMARY KEY,
    date DATE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    country_code VARCHAR(10),
    is_national BOOLEAN DEFAULT true,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE devices (
    id CHAR(36) PRIMARY KEY,
    employee_id CHAR(36),
    device_id VARCHAR(255) UNIQUE NOT NULL,
    device_type VARCHAR(50),
    os_version VARCHAR(50),
    app_version VARCHAR(50),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE device_logs (
    id CHAR(36) PRIMARY KEY,
    device_id CHAR(36),
    ip_address VARCHAR(45),
    user_agent TEXT,
    event VARCHAR(50), 
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);
