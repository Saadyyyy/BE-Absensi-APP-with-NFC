# Attendance System Backend

Backend API untuk sistem absensi sekolah menggunakan kartu NFC yang dibangun dengan Go, Echo framework, dan GORM.

## 🚀 Fitur

- **Authentication & Authorization**: JWT-based authentication dengan role-based access control
- **NFC Card Management**: Registrasi dan manajemen kartu NFC untuk siswa
- **Attendance Tracking**: Pencatatan absensi masuk dan keluar menggunakan NFC
- **User Management**: Manajemen pengguna dengan role (user, admin, super_admin)
- **Database Integration**: PostgreSQL dengan GORM ORM

## 📁 Struktur Folder

```
be/
├── config/          # Konfigurasi database
├── controllers/     # HTTP handlers
├── middleware/      # Custom middleware (JWT, CORS, dll)
├── models/          # Database models
├── routes/          # Route definitions
├── utils/           # Utility functions
├── .env.example     # Environment variables template
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── server.go        # Main application entry point
└── README.md        # Dokumentasi
```

## 🛠 Teknologi

- **Go 1.23.3**
- **Echo v4** - Web framework
- **GORM** - ORM untuk database operations
- **PostgreSQL** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing
- **UUID** - Unique identifiers

## 📦 Dependencies

```go
// Framework & Database
github.com/labstack/echo/v4
gorm.io/gorm
gorm.io/driver/postgres

// Authentication & Security
github.com/golang-jwt/jwt/v5
golang.org/x/crypto

// Utilities
github.com/google/uuid
```

## 🚀 Quick Start

### 1. Clone & Setup
```bash
cd be/
cp .env.example .env
# Edit .env dengan konfigurasi database Anda
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Setup Database
Pastikan PostgreSQL sudah running dan buat database:
```sql
CREATE DATABASE attendance_db;
```

### 4. Run Application
```bash
go run server.go
```

Server akan berjalan di `http://localhost:1323`

## 📚 API Endpoints

### Health Check
```
GET /api/v1/health
```

### Authentication
```
POST /api/v1/auth/login
POST /api/v1/auth/register
POST /api/v1/auth/refresh
```

### User Profile (Protected)
```
GET /api/v1/profile
```

### Attendance (Protected)
```
POST /api/v1/attendance/record
GET /api/v1/attendance/today
GET /api/v1/attendance/history/:student_id
```

### Admin (Admin Role Required)
```
POST /api/v1/admin/nfc/register
```

### Super Admin (Super Admin Role Required)
```
GET /api/v1/super-admin/users
```

## 🔐 Authentication

### Login Request
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Login Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "user@example.com",
    "role": "user"
  }
}
```

### Authorization Header
```
Authorization: Bearer <your-jwt-token>
```

## 📊 Database Models

### User
- ID (UUID)
- Email (unique)
- Password (hashed)
- Name
- Role (user/admin/super_admin)
- IsActive
- Timestamps

### Student
- ID (UUID)
- NFC UID (unique)
- Name
- Class
- Student ID (unique)
- School ID
- IsActive
- Timestamps

### Attendance
- ID (UUID)
- Student ID (foreign key)
- Date
- Time In
- Time Out
- Status (present/late/absent)
- Timestamps

### School
- ID (UUID)
- Name
- Address
- Phone
- Email
- IsActive
- Timestamps

## 🔧 Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=attendance_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Server
PORT=1323
ENVIRONMENT=development
```

## 🎯 NFC Attendance Flow

1. **Registrasi Kartu**: Admin mendaftarkan kartu NFC ke siswa
2. **Check-in**: Siswa tap kartu → sistem catat waktu masuk
3. **Check-out**: Siswa tap kartu lagi → sistem catat waktu keluar
4. **Status**: Otomatis menentukan status (present/late) berdasarkan waktu

## 🔒 Security Features

- **Password Hashing**: bcrypt untuk hash password
- **JWT Authentication**: Secure token-based auth
- **Role-based Access**: Different access levels
- **Input Validation**: Request validation
- **CORS Protection**: Cross-origin request handling

## 🚀 Deployment

### Build Binary
```bash
go build -o attendance-server server.go
```

### Run Binary
```bash
./attendance-server
```

### Docker (Optional)
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

## 📝 Development Notes

- Database auto-migration dijalankan saat startup
- JWT token expire dalam 24 jam
- Refresh token expire dalam 7 hari
- Default school start time: 07:30 (untuk menentukan status late)
- Semua UUID menggunakan `github.com/google/uuid`

## 🤝 Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## 📄 License

MIT License