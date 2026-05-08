# 🚀 Golang Production-Ready REST API Authentication

API Otentikasi yang dibangun menggunakan **Golang**, **Gin Framework**, dan **GORM**. Project ini mengimplementasikan *Clean Architecture* dan *Repository Pattern* untuk memastikan kode mudah dikelola, diuji, dan dikembangkan.

## ✨ Fitur Utama
- **User Management**: Register, Login, Get Profile.
- **Security**: Password hashing dengan `bcrypt`.
- **JWT Token**: Access Token (15m) & Refresh Token (7d).
- **Token Management**: Refresh token disimpan di DB & mendukung fitur Revoke (Logout).
- **Security Audit**: Pencatatan log login (IP Address, User Agent, Status).
- **Password Recovery**: Fitur Forgot & Reset Password dengan token unik.
- **Authorization**: Role-based access control (Admin & Guest).
- **Documentation**: Swagger UI terintegrasi.

## 🛠️ Tech Stack
- **Language**: [Golang](https://go.dev/)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: MySQL / MariaDB
- **Token**: JWT (JSON Web Token)
- **Documentation**: Swagger (Swaggo)

## 📂 Struktur Folder
```
.
├── cmd/
│   └── main.go             # Entry point aplikasi
├── config/                 # Konfigurasi Database & Env
├── controller/             # HTTP Handlers (Menerima input user)
├── docs/                   # Swagger Documentation (Auto-generated)
├── middleware/             # JWT & Role Authentication
├── model/                  # Struct Database & DTO
├── repository/             # Layer Database (Query)
├── routes/                 # Definisi Endpoint API
├── service/                # Logika Bisnis (Core Logic)
├── utils/                  # Helper (JWT, Bcrypt, Response)
└── .env                    # Environment Variables
.
```

## 🚀 Cara Menjalankan Project
1. Clone Project
```
git clone https://github.com/pudinazhar/go-auth-api.git
```
cd go-auth-api

2. Install Dependensi
```
go mod tidy
```
3. Konfigurasi Database

Buat database di MySQL/MariaDB bernama go_auth_db. Kemudian buat file .env di root folder:
```
PORT=8080
DB_USER=root
DB_PASSWORD=password_kamu
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=go_auth_db

JWT_SECRET=SangatRahasia123!
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```
4. Jalankan Aplikasi
```
go run cmd/main.go
```
Database akan otomatis ter-migrasi saat aplikasi dijalankan.

## 📖 Dokumentasi API
Setelah aplikasi berjalan, Anda dapat mengakses dokumentasi interaktif (Swagger) di:
```
http://localhost:8080/swagger/index.html
```
Perbaharui Doc Dokumentasi
```
swag init -g cmd/main.go
```

## 🧪 Menjalankan Unit Test
```
go test ./utils -v
```

## 🔐 Keamanan
- Password di-hash menggunakan bcrypt dengan cost default.
- Header Authorization menggunakan format: Bearer <token>.
- Refresh Token di-rotate setiap kali digunakan untuk mencegah pembajakan sesi.

## Error
Jika mengalami error saat menggunakan docker, itu karena mysql membutuhkan waktu beberapa saat untuk hidup, kita cukup melakukan perintah ini saja,
```
docker-compose start app
```

## Kontak Saya
- [Telegram](https://t.me/pudin_ira)
- [Website](https://italazhar.com)