# Employee Attendance System Design

Sistem ini dirancang untuk mendukung kebutuhan 25.000+ karyawan dari sebuah perusahaan ritel yang menerapkan kebijakan Work From Anywhere (WFA). Karyawan tersebar di seluruh wilayah Indonesia dan sebagian di luar negeri. Tujuan dari sistem ini adalah mencatat dan memantau kehadiran karyawan dengan dukungan lokasi, waktu, dan validasi yang kuat.

## Assumptions

- Perusahaan memiliki +- 25.000 karyawan tersebar di seluruh Indonesia dan beberapa bekerja di luar negeri (WFA).
- Karyawan melakukan absensi melalui aplikasi mobile (Android/iOS) atau web.
- Absensi menggunakan geolokasi dan verifikasi selfie (deteksi wajah) untuk menghindari kecurangan.
- Sistem memiliki kesadaran zona waktu (jam lokal karyawan sesuai lokasi).
- Tersedia fitur cuti, izin, dan absensi manual jika terjadi kesalahan teknis.
- Setiap karyawan hanya bisa absen jika berada di lokasi yang diizinkan (bisa didefinisikan oleh admin).

---

## Functional Requirements

- Karyawan dapat:
  - Melakukan absensi masuk & keluar
  - Melihat histori absensi
  - Mengajukan cuti/izin
- Admin dapat:
  - Melihat dashboard laporan absensi
  - Melihat anomali (tidak absen, telat, lokasi tidak valid, dsb)
  - Menyusun shift dan jadwal kerja
  - Mengekspor laporan absensi
- Sistem dapat:
  - Mendeteksi pemalsuan (absensi palsu) via GPS/Pengecekan Wajah
  - Menyesuaikan waktu berdasarkan zona waktu pengguna
  - Mengirim notifikasi saat lupa absen keluar

---

## Non-Functional Requirements

| Kebutuhan               | Detail                                                         |
| ----------------------- | -------------------------------------------------------------- |
| Ketersediaan            | 99.9% waktu aktif, ketersediaan tinggi via infrastruktur cloud |
| Skalabilitas            | Mampu menangani 25.000+ karyawan aktif secara paralel          |
| Keamanan                | Otentikasi JWT, SSL/TLS, enkripsi foto wajah, anti-pemalsuan   |
| Performa                | Waktu respons < 300ms untuk permintaan absensi                 |
| Kompatibilitas          | Aplikasi Mobile (Android/iOS), Browser Web                     |
| Pencatatan & Pemantauan | Tersedia log aktivitas, pelacakan kesalahan, dan sistem audit  |
| Cadangan                | Cadangan harian database absensi dan metadata selfie           |

---

## Core Entities (ERD)

1. Users

   - id: CHAR(36) PRIMARY KEY
   - email: VARCHAR(255) UNIQUE NOT NULL
   - password: TEXT NOT NULL
   - full_name: VARCHAR(255) NOT NULL
   - role: ENUM('admin', 'hr', 'manager', 'employee')
   - location: VARCHAR(255)
   - timezone: VARCHAR(100)
   - photo_url: TEXT
   - status: ENUM('active', 'inactive', 'terminated', 'suspended')
   - email_verified_at: DATETIME
   - created_at, updated_at, deleted_at: DATETIME

2. Employees

   - id: CHAR(36) PRIMARY KEY
   - user_id: CHAR(36) FOREIGN KEY
   - department_id: CHAR(36) FOREIGN KEY
   - name: VARCHAR(255)
   - email: VARCHAR(255) UNIQUE
   - role: VARCHAR(100)
   - location: VARCHAR(255)
   - timezone: VARCHAR(100)
   - photo_url: TEXT
   - status: ENUM('active', 'inactive', 'terminated', 'suspended')
   - join_date: DATE
   - reporting_to: CHAR(36)

3. Attendances

   - id: CHAR(36) PRIMARY KEY
   - employee_id: CHAR(36) FOREIGN KEY
   - date: DATE
   - check_in_time: DATETIME
   - check_out_time: DATETIME
   - type: VARCHAR(50)
   - location: VARCHAR(255)
   - work_location_id: CHAR(36)
   - selfie_url: TEXT
   - status: ENUM('present', 'absent', 'late', 'early_leave')
   - device_info: JSON
   - ip_address: VARCHAR(45)
   - created_at, updated_at: DATETIME

4. Schedules

   - id: CHAR(36) PRIMARY KEY
   - employee_id: CHAR(36) FOREIGN KEY
   - date: DATE
   - shift_start: TIME
   - shift_end: TIME
   - break_start: TIME
   - break_end: TIME
   - work_location_id: CHAR(36)
   - schedule_type: ENUM('regular', 'shift', 'flexible')
   - created_at, updated_at: DATETIME

5. LeaveRequests

   - id: CHAR(36) PRIMARY KEY
   - employee_id: CHAR(36) FOREIGN KEY
   - start_date: DATE
   - end_date: DATE
   - type: ENUM('annual', 'sick', 'unpaid', 'maternity', 'paternity')
   - reason: TEXT
   - status: ENUM('pending', 'approved', 'rejected')
   - reviewed_by: VARCHAR(255)
   - reviewed_at: TIMESTAMP
   - note: TEXT
   - created_at, updated_at: DATETIME

6. AbsenceAnomalies

   - id: CHAR(36) PRIMARY KEY
   - employee_id: CHAR(36) FOREIGN KEY
   - date: DATE
   - type: ENUM('late', 'not_present', 'left_early', 'forgot_checkin')
   - note: TEXT
   - verified: BOOLEAN
   - verified_by: VARCHAR(255)
   - verified_at: TIMESTAMP
   - created_at: DATETIME

7. Notifications

   - id: CHAR(36) PRIMARY KEY
   - employee_id: CHAR(36) FOREIGN KEY
   - type: ENUM('reminder', 'warning', 'info')
   - message: TEXT
   - send_at: DATETIME
   - created_at: DATETIME

8. Departments

   - id: CHAR(36) PRIMARY KEY
   - name: VARCHAR(255)
   - head_employee_id: CHAR(36)
   - location: VARCHAR(255)
   - timezone: VARCHAR(100)
   - wfa_policy: JSON
   - created_at, updated_at: DATETIME

9. WorkLocations

   - id: CHAR(36) PRIMARY KEY
   - name: VARCHAR(255)
   - address: TEXT
   - country_code: VARCHAR(10)
   - latitude: DECIMAL(10,8)
   - longitude: DECIMAL(11,8)
   - radius: INTEGER
   - status: ENUM('active', 'inactive')
   - created_at, updated_at: DATETIME

10. EmployeeWorkLocations

    - id: CHAR(36) PRIMARY KEY
    - employee_id: CHAR(36) FOREIGN KEY
    - work_location_id: CHAR(36) FOREIGN KEY
    - is_primary: BOOLEAN
    - created_at, updated_at: DATETIME

11. Holidays

    - id: CHAR(36) PRIMARY KEY
    - date: DATE
    - name: VARCHAR(255)
    - description: TEXT
    - country_code: VARCHAR(10)
    - is_national: BOOLEAN
    - created_at, updated_at: DATETIME

12. Devices

    - id: CHAR(36) PRIMARY KEY
    - employee_id: CHAR(36) FOREIGN KEY
    - device_id: VARCHAR(255) UNIQUE
    - device_type: VARCHAR(50)
    - os_version: VARCHAR(50)
    - app_version: VARCHAR(50)
    - created_at, updated_at: DATETIME

13. DeviceLogs
    - id: CHAR(36) PRIMARY KEY
    - device_id: CHAR(36) FOREIGN KEY
    - ip_address: VARCHAR(45)
    - user_agent: TEXT
    - event: VARCHAR(50)
    - created_at: DATETIME

---

## Use Cases

1. **Pendaftaran Pengguna**
   - Pengguna dapat mendaftar dengan email dan kata sandi.
   - Pengguna dapat mendaftar sebagai karyawan atau admin.
   - Pengguna dapat mendaftar dengan foto profil.
2. **Masuk Pengguna**
   - Pengguna dapat masuk dengan email dan kata sandi.
   - Pengguna dapat masuk dengan akun media sosial.
3. **Profil Pengguna**
   - Pengguna dapat melihat profil mereka.
   - Pengguna dapat memperbarui profil mereka.
   - Pengguna dapat mengunggah foto profil.
4. **Absensi**
   - Pengguna dapat melihat riwayat absensi mereka.
   - Pengguna dapat melihat ringkasan absensi mereka.
   - Pengguna dapat melihat peta absensi mereka.
   - Pengguna dapat melihat laporan absensi mereka.
5. **Manajemen Cuti**
   - Pengguna dapat melihat riwayat cuti mereka.
   - Pengguna dapat mengajukan cuti.
   - Pengguna dapat melihat saldo cuti mereka.
   - Pengguna dapat melihat status cuti mereka.
6. **Manajemen Jadwal**
   - Pengguna dapat melihat jadwal mereka.
   - Pengguna dapat melihat detail jadwal mereka.
   - Pengguna dapat memperbarui jadwal mereka.
7. **Deteksi Anomali**
   - Pengguna dapat melihat riwayat anomali mereka.
   - Pengguna dapat melihat detail anomali mereka.
   - Pengguna dapat melihat laporan anomali mereka.
8. **Notifikasi**
   - Pengguna dapat melihat riwayat notifikasi mereka.
   - Pengguna dapat melihat detail notifikasi mereka.
   - Pengguna dapat menandai notifikasi sebagai telah dibaca.
9. **Pembuatan Laporan**
   - Pengguna dapat membuat laporan absensi.
   - Pengguna dapat membuat laporan cuti.
   - Pengguna dapat membuat laporan ringkasan absensi.
10. **Dasbor Admin**
    - Admin dapat melihat ringkasan absensi.
    - Admin dapat melihat riwayat absensi.
    - Admin dapat melihat peta absensi.
    - Admin dapat melihat laporan absensi.
    - Admin dapat melihat riwayat cuti.
    - Admin dapat melihat laporan cuti.
    - Admin dapat melihat riwayat anomali.
    - Admin dapat melihat laporan anomali.
    - Admin dapat melihat profil pengguna.
    - Admin dapat melihat riwayat absensi pengguna.
    - Admin dapat melihat riwayat cuti pengguna.
    - Admin dapat melihat riwayat anomali pengguna.

## System Architecture

**Komponen:**

Komponen Teknologi

- **Frontend:** React.js / VueJS / NextJS (Web), Flutter / React Native (Mobile)
- **Backend Golang:** NodeJS / Golang / Laravel (Express, Gin, Lumen)
- **API Gateway:** Kong / NGINX Reverse Proxy
- **Database:** PostgreSQL / MySQL
- **File Storage:** Amazon S3 / Google Cloud Storage
- **Authentication:** JWT / OAuth2 + LDAP / Firebase Authentication
- **Message Queue:** RabbitMQ / Apache Kafka (untuk notifikasi dan pemrosesan async)
- **Monitoring Prometheus:** + Grafana
- **CI/CD:** GitHub Actions / GitLab CI
- **Containerization:** Docker + Kubernetes (EKS / GKE / Self-hosted)

---

## High Level Flow

1. Alur Otentikasi

   - Pengguna melakukan login melalui aplikasi mobile/web
   - Sistem memvalidasi kredensial
   - Pengguna mendapatkan akses sesuai peran
   - Sistem menyimpan informasi sesi

2. Alur Absensi

   - Karyawan membuka fitur absensi
   - Sistem memvalidasi jadwal kerja
   - Sistem mengecek lokasi via GPS
   - Sistem memvalidasi radius lokasi kerja
   - Karyawan melakukan selfie
   - Sistem melakukan pengenalan wajah
   - Sistem mencatat waktu & lokasi
   - Sistem mengirim notifikasi sukses/gagal

3. Alur Manajemen Cuti

   - Karyawan mengajukan cuti/izin
   - Sistem memvalidasi sisa cuti
   - Sistem mengirim notifikasi ke pemberi persetujuan
   - Pemberi persetujuan mereview pengajuan
   - Sistem memperbarui status cuti
   - Sistem mengupdate catatan absensi
   - Sistem mengirim notifikasi hasil

4. Alur Manajemen Jadwal

   - Admin membuat jadwal kerja
   - Sistem memvalidasi konflik jadwal
   - Sistem mengirim notifikasi ke karyawan
   - Karyawan menerima jadwal baru
   - Sistem mengupdate basis data

5. Alur Deteksi Anomali

   - Sistem mendeteksi anomali absensi
   - Sistem mencatat jenis anomali
   - Sistem mengirim notifikasi ke admin/HR
   - Admin/HR melakukan verifikasi
   - Sistem mengupdate status anomali

6. Alur Notifikasi

   - Sistem mengecek pemicu kejadian
   - Sistem mempersiapkan templat notifikasi
   - Sistem menentukan penerima
   - Sistem mengirim notifikasi
   - Sistem mencatat status pengiriman

7. Alur Pembuatan Laporan
   - Admin memilih jenis laporan
   - Sistem mengumpulkan data
   - Sistem memproses data
   - Sistem menghasilkan laporan
   - Sistem menyimpan riwayat laporan
