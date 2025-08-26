# API Documentation - Sistem Absensi NFC

## Base URL
```
http://localhost:1323
```

## Authentication
API ini menggunakan JWT (JSON Web Token) untuk autentikasi. Token harus disertakan dalam header `Authorization` dengan format:
```
Authorization: Bearer <your_jwt_token>
```

## Response Format
Semua response menggunakan format JSON dengan struktur berikut:

### Success Response
```json
{
  "success": true,
  "message": "Success message",
  "data": {}
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```

## Endpoints

### 1. Authentication

#### POST /auth/register
Mendaftarkan user baru ke sistem.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "role": "user" // optional: user, admin, super_admin (default: user)
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid-string",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "jwt-access-token",
    "refresh_token": "jwt-refresh-token"
  }
}
```

#### POST /auth/login
Login user ke sistem.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid-string",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "is_active": true
    },
    "access_token": "jwt-access-token",
    "refresh_token": "jwt-refresh-token"
  }
}
```

#### POST /auth/refresh
Memperbarui access token menggunakan refresh token.

**Request Body:**
```json
{
  "refresh_token": "jwt-refresh-token"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "new-jwt-access-token",
    "refresh_token": "new-jwt-refresh-token"
  }
}
```

#### GET /auth/profile
Mendapatkan profil user yang sedang login.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "user": {
      "id": "uuid-string",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 2. Attendance Management

#### POST /attendance/record
Mencatat absensi siswa menggunakan NFC.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "nfc_uid": "04:52:3A:B2:C1:90:80"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Attendance recorded successfully",
  "data": {
    "attendance": {
      "id": "uuid-string",
      "student_id": "uuid-string",
      "date": "2024-01-01",
      "time_in": "2024-01-01T08:00:00Z",
      "time_out": null,
      "status": "present",
      "created_at": "2024-01-01T08:00:00Z",
      "updated_at": "2024-01-01T08:00:00Z"
    },
    "student": {
      "id": "uuid-string",
      "name": "Jane Doe",
      "class": "10A",
      "student_id": "2024001"
    },
    "action": "check_in" // atau "check_out"
  }
}
```

#### GET /attendance/history/:student_id
Mendapatkan riwayat absensi siswa.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `start_date` (optional): Format YYYY-MM-DD
- `end_date` (optional): Format YYYY-MM-DD
- `limit` (optional): Default 50
- `offset` (optional): Default 0

**Response:**
```json
{
  "success": true,
  "message": "Attendance history retrieved successfully",
  "data": {
    "student": {
      "id": "uuid-string",
      "name": "Jane Doe",
      "class": "10A",
      "student_id": "2024001"
    },
    "attendances": [
      {
        "id": "uuid-string",
        "date": "2024-01-01",
        "time_in": "2024-01-01T08:00:00Z",
        "time_out": "2024-01-01T15:00:00Z",
        "status": "present",
        "created_at": "2024-01-01T08:00:00Z",
        "updated_at": "2024-01-01T15:00:00Z"
      }
    ],
    "total": 1,
    "limit": 50,
    "offset": 0
  }
}
```

#### GET /attendance/today
Mendapatkan absensi semua siswa hari ini.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Today's attendance retrieved successfully",
  "data": {
    "date": "2024-01-01",
    "attendances": [
      {
        "id": "uuid-string",
        "student": {
          "id": "uuid-string",
          "name": "Jane Doe",
          "class": "10A",
          "student_id": "2024001"
        },
        "time_in": "2024-01-01T08:00:00Z",
        "time_out": "2024-01-01T15:00:00Z",
        "status": "present"
      }
    ],
    "summary": {
      "total_students": 100,
      "present": 85,
      "late": 10,
      "absent": 5
    }
  }
}
```

### 3. NFC Card Management (Admin Only)

#### POST /admin/nfc/register
Mendaftarkan kartu NFC baru untuk siswa.

**Headers:**
```
Authorization: Bearer <admin_access_token>
```

**Request Body:**
```json
{
  "nfc_uid": "04:52:3A:B2:C1:90:80",
  "student_name": "Jane Doe",
  "student_id": "2024001",
  "class": "10A",
  "school_id": "uuid-string"
}
```

**Response:**
```json
{
  "success": true,
  "message": "NFC card registered successfully",
  "data": {
    "student": {
      "id": "uuid-string",
      "nfc_uid": "04:52:3A:B2:C1:90:80",
      "name": "Jane Doe",
      "class": "10A",
      "student_id": "2024001",
      "school_id": "uuid-string",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

## Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | BAD_REQUEST | Request tidak valid atau parameter hilang |
| 401 | UNAUTHORIZED | Token tidak valid atau expired |
| 403 | FORBIDDEN | Tidak memiliki permission untuk akses resource |
| 404 | NOT_FOUND | Resource tidak ditemukan |
| 409 | CONFLICT | Data sudah ada (duplicate) |
| 500 | INTERNAL_ERROR | Server error |

## Data Models

### User
```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "role": "user|admin|super_admin",
  "is_active": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Student
```json
{
  "id": "uuid",
  "nfc_uid": "string",
  "name": "string",
  "class": "string",
  "student_id": "string",
  "school_id": "uuid",
  "is_active": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### School
```json
{
  "id": "uuid",
  "name": "string",
  "address": "string",
  "phone": "string",
  "email": "string",
  "is_active": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Attendance
```json
{
  "id": "uuid",
  "student_id": "uuid",
  "date": "date",
  "time_in": "datetime|null",
  "time_out": "datetime|null",
  "status": "present|late|absent",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Testing dengan cURL

### Register User
```bash
curl -X POST http://localhost:1323/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123",
    "name": "Admin User",
    "role": "admin"
  }'
```

### Login
```bash
curl -X POST http://localhost:1323/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

### Get Profile
```bash
curl -X GET http://localhost:1323/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Record Attendance
```bash
curl -X POST http://localhost:1323/attendance/record \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "nfc_uid": "04:52:3A:B2:C1:90:80"
  }'
```

## Rate Limiting
API ini menggunakan rate limiting untuk mencegah abuse:
- 100 requests per menit untuk endpoint authentication
- 1000 requests per menit untuk endpoint lainnya

## CORS
API mendukung CORS untuk semua origin dalam mode development. Untuk production, pastikan mengkonfigurasi CORS dengan domain yang tepat.

## Security
- Semua password di-hash menggunakan bcrypt
- JWT token memiliki expiry time (default: 24 jam untuk access token, 7 hari untuk refresh token)
- Middleware authentication melindungi endpoint yang memerlukan login
- Role-based access control untuk endpoint admin dan super admin

## Environment Variables
Pastikan file `.env` sudah dikonfigurasi dengan benar:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=dimasrio
DB_SSLMODE=disable
JWT_SECRET=your-super-secret-jwt-key
PORT=1323
ENVIRONMENT=development
```