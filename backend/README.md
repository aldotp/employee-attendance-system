# 📱 Employee Attendance System

Because of lack of time, i only create for the monolith not split into microservices

## 📋 Table of Contents

- [About System](#about-system)
- [Assumptions](#assumptions)
- [System Requirements](#system-requirements)
- [Database Entities](#database-entities)
- [System Features](#system-features)
- [System Architecture](#system-architecture)
- [Process Flow](#process-flow)

## 🎯 About System

This system is designed to support 25,000+ employees of a retail company implementing Work From Anywhere (WFA) policy. Employees are spread across Indonesia and some work overseas. The system aims to record and monitor employee attendance with strong location, time, and validation support.

## 📌 Assumptions

- ✅ 25,000+ employees across Indonesia and overseas
- 📱 Attendance via mobile or web application
- 🌍 GPS validation
- 🕒 Multi-timezone support
- 📝 Integrated leave and permission features
- 📍 Work location validation

## 💻 System Requirements

### Functional Requirements

#### 👥 Employee

- Check-in & check-out
- View attendance history
- Submit leave requests

#### 👨‍💼 Admin

- Attendance report dashboard
- Monitor anomalies
- Work schedule management
- Export reports

#### 🤖 System

- Fraud detection
- Timezone management
- Notification system

### Non-Functional Requirements

| Aspect           | Specification            |
| ---------------- | ------------------------ |
| 🔄 Availability  | 99.9% uptime             |
| 📈 Scalability   | 25,000+ active users     |
| 🔒 Security      | JWT, SSL/TLS, Encryption |
| ⚡ Performance   | Response < 300ms         |
| 📱 Compatibility | Mobile & Web             |
| 📊 Monitoring    | Log & Audit Trail        |
| 💾 Backup        | Daily Backup             |

## 📚 Database Entities

[Database schema remains unchanged as it contains technical specifications]

## 🎉 System Features

### 1. User Management

- Registration & Login
- User Profile
- Access Management

### 2. Attendance

- Check-in/Check-out
- Location Validation
- Face Recognition
- Attendance History

### 3. Leave Management

- Leave Application
- Approval Workflow
- Leave Balance
- Leave History

### 4. Schedule Management

- Shift Scheduling
- Work Rotation
- Work Calendar

### 5. Monitoring & Reports

- Dashboard Analytics
- Attendance Reports
- Anomaly Detection
- Data Export

## 🏗 System Architecture

### Technology Stack

- **Frontend**: React.js/VueJS/NextJS, Flutter/React Native
- **Backend**: Golang/NodeJS/Laravel
- **Database**: PostgreSQL/MySQL
- **Storage**: Amazon S3/Google Cloud
- **Security**: JWT/OAuth2/Firebase Auth
- **Queue**: RabbitMQ/Kafka
- **Monitoring**: Prometheus + Grafana
- **DevOps**: Docker, Kubernetes, CI/CD

## 🔄 Process Flow

### 1. Authentication Flow

- Login via mobile/web
- Credential validation
- Role-based access
- Session storage

### 2. Attendance Flow

- Open attendance feature
- Schedule validation
- GPS location check
- Radius validation
- Take selfie
- Face recognition
- Record time & location
- Send notification

### 3. Leave Flow

- Submit leave request
- Balance validation
- Notify approver
- Review submission
- Update status
- Update attendance
- Send notification

### 4. Schedule Flow

- Create schedule
- Conflict validation
- Notify employee
- Schedule confirmation
- Update database

### 5. Anomaly Flow

- Detect anomaly
- Record type
- Notify admin
- Verification
- Update status

### 6. Notification Flow

- Check trigger events
- Prepare notification template
- Determine recipient
- Send notification
- Record sending status

### 7. Report Generation Flow

- Select report type
- Gather data
- Process data
- Generate report
- Store report history

### How tu run

#### Run the Rest Api

```shell
go run cmd/main.go rest
```

### Run Worker Generate Reporting

```shell
# Up the migration
go run cmd/main.go consumer generate_reporting
```

## Run Migrations

```shell
# Up the migration
migrate -path ./internal/adapter/storage/postgres/migrations \
        -database "postgres://user:12345678a@127.0.0.1:5432/user-service?sslmode=disable" \
        -verbose up


# Down the migration
migrate -path ./internal/adapter/storage/postgres/migrations \
        -database "postgres://user:12345678a@127.0.0.1:5432/user-service?sslmode=disable" \
        -verbose down
```

## New Migrartions

```shell
// add new migration
migrate create -ext sql -dir  ./internal/adapter/storage/postgres/migrations -format "20060102150405" add_table_carts
```

#### Generate Swagger

```shell
swag init -g cmd/main.go
```
