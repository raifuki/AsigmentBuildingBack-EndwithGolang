# 🚀 Task Management API

REST API hoàn chỉnh cho hệ thống quản lý công việc, được xây dựng bằng **Golang + Gin + PostgreSQL + Redis**.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/Gin-v1.10-00ADD8)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-316192?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 📋 Mục lục

- [Tính năng](#-tính-năng)
- [Tech Stack](#-tech-stack)
- [Kiến trúc](#-kiến-trúc)
- [Cài đặt](#-cài-đặt)
- [Chạy ứng dụng](#-chạy-ứng-dụng)
- [API Endpoints](#-api-endpoints)
- [Testing](#-testing)
- [Cấu trúc thư mục](#-cấu-trúc-thư-mục)

---

## ✨ Tính năng

### 🔐 Authentication & Authorization
- Đăng ký / Đăng nhập với JWT
- Mã hóa mật khẩu bằng **bcrypt**
- Phân quyền **admin** / **user**
- Middleware xác thực JWT bảo vệ private routes

### 📦 CRUD Resources
- **User**: Register, Login, Get, Update, Delete
- **Project**: Full CRUD + pagination
- **Task**: Full CRUD + filter theo status, assign user
- **Comment**: Full CRUD + phân quyền owner

### ⚡ Performance
- **Redis caching** cho GET endpoints (5-10 phút TTL)
- **Cache invalidation** tự động khi có write operation
- **Rate limiting** theo Token Bucket pattern (Lua script atomic)
- **Pagination + filtering** cho tất cả list endpoints

### 🛡️ Security
- **Rate limiting 4 tầng**: global / auth / user / write
- **CORS** middleware
- **Centralized error handling** với recovery từ panic
- **Input validation** với `go-playground/validator`

### 🐳 DevOps
- **Docker multi-stage build** (image ~15MB)
- **docker-compose** với PostgreSQL + Redis
- **Graceful shutdown**
- **Health check** endpoint
- **CI/CD** với GitHub Actions

---

## 🛠 Tech Stack

| Layer | Technology |
|-------|-----------|
| **Language** | Go 1.24 |
| **Framework** | Gin v1.10 |
| **ORM** | GORM v1.25 |
| **Database** | PostgreSQL 16 |
| **Cache** | Redis 7 |
| **Auth** | JWT (golang-jwt/v5) + bcrypt |
| **Validation** | go-playground/validator/v10 |
| **Testing** | testify + sqlmock |
| **Container** | Docker + docker-compose |
| **CI/CD** | GitHub Actions |

---

## 🏗 Kiến trúc
┌─────────────┐
│ Client │
└──────┬──────┘
│ HTTP
▼
┌─────────────────────┐
│ Gin Router │
├─────────────────────┤
│ Middlewares │
│ ├─ Logger │
│ ├─ CORS │
│ ├─ Rate Limiter │
│ ├─ Auth (JWT) │
│ └─ Error Handler │
├─────────────────────┤
│ Handlers │ ← Parse request, call service
├─────────────────────┤
│ Services │ ← Business logic + Cache
├─────────────────────┤
│ Repositories │ ← Data access (GORM)
└──────┬──────────────┘
│
┌───┴────┐
▼ ▼
┌──────┐ ┌───────┐
│ PG │ │ Redis │
└──────┘ └───────┘

text

---

## 💻 Cài đặt

### Yêu cầu hệ thống

- **Go** ≥ 1.24 — [Download](https://go.dev/dl/)
- **Docker** + **Docker Compose** — [Download](https://www.docker.com/products/docker-desktop/)
- **Git**

### Clone repository

```bash
git clone https://github.com/yourname/task-management.git
cd task-management
Cấu hình môi trường
Copy file .env.example thành .env:

bash
cp .env.example .env
Sửa các biến trong .env:

env
APP_PORT=8080
APP_ENV=development

# Local dev (nếu chạy Docker thì không cần)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskdb
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Bắt buộc đổi trong production
JWT_SECRET=change-me-to-random-32-chars
JWT_EXPIRED_HOURS=24
🚀 Chạy ứng dụng
Cách 1: Docker Compose (khuyên dùng)
bash
# Chạy toàn bộ stack (API + PostgreSQL + Redis)
docker compose up --build

# Chạy nền
docker compose up --build -d

# Xem log
docker compose logs -f api

# Dừng
docker compose down

# Dừng và xóa data
docker compose down -v
API sẽ chạy tại: http://localhost:8080

Cách 2: Chạy local (cần PG + Redis sẵn)
bash
# Chỉ chạy database + redis bằng Docker
docker run -d --name task_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=taskdb \
  -p 5432:5432 postgres:16-alpine

docker run -d --name task_redis -p 6379:6379 redis:7-alpine

# Chạy app
go run ./cmd/api
Kiểm tra
bash
curl http://localhost:8080/health
# → {"status":"ok"}
📚 API Endpoints
Base URL
text
http://localhost:8080/api/v1
🔓 Public Endpoints
Health Check
http
GET /health
Register
http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "123456"
}
Login
http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "123456"
}
Response:

json
{
  "success": true,
  "message": "Login success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { "id": 1, "name": "John Doe", "email": "john@example.com", "role": "user" }
  }
}
🔒 Protected Endpoints
Tất cả request sau cần header:

text
Authorization: Bearer <JWT_TOKEN>
Projects
Method	Endpoint	Mô tả
POST	/projects	Tạo project
GET	/projects?page=1&limit=10	List projects (có pagination)
GET	/projects/:id	Chi tiết project
PUT	/projects/:id	Cập nhật project
DELETE	/projects/:id	Xóa project (soft delete)
Ví dụ tạo project:

http
POST /projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Website Redesign",
  "description": "Redesign company website"
}
Tasks
Method	Endpoint	Mô tả
POST	/tasks/project/:projectId	Tạo task
GET	/tasks?project_id=1&status=todo&page=1&limit=10	List tasks (có filter)
GET	/tasks/:id	Chi tiết task
PUT	/tasks/:id	Cập nhật task
DELETE	/tasks/:id	Xóa task
Ví dụ tạo task:

http
POST /tasks/project/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Design homepage",
  "description": "Create mockup for homepage",
  "status": "todo",
  "priority": 2,
  "due_date": "2026-10-15T00:00:00Z",
  "assignee_id": 2
}
Comments
Method	Endpoint	Mô tả
POST	/comments/task/:taskId	Tạo comment
GET	/comments/task/:taskId	List comments của task
GET	/comments/:id	Chi tiết comment
PUT	/comments/:id	Cập nhật comment (owner only)
DELETE	/comments/:id	Xóa comment
📊 Response Format
Success:

json
{
  "success": true,
  "message": "OK",
  "data": { ... }
}
Error:

json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "Email": "email",
    "Password": "min"
  }
}
📈 Rate Limiting Headers
Mỗi response có:

text
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
Retry-After: 12               (chỉ khi bị 429)
X-RateLimit-Reset: 1727889234 (chỉ khi bị 429)
⚠️ Status Codes
Code	Ý nghĩa
200	OK
201	Created
400	Bad Request (validation)
401	Unauthorized (token thiếu/sai)
403	Forbidden (không có quyền)
404	Not Found
429	Too Many Requests (rate limit)
500	Internal Server Error
🧪 Testing
Chạy test
bash
# Chạy tất cả test
go test ./... -v

# Test với coverage
go test ./... -v -race -coverprofile=coverage.out

# Xem coverage report
go tool cover -html=coverage.out

# Chỉ test 1 package
go test ./internal/services/... -v
Code quality
bash
# Format code
gofmt -w .

# Static analysis
go vet ./...

# Lint (cần cài golangci-lint)
golangci-lint run ./...
Test API bằng cURL
bash
# Health
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"123456"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# Save token
TOKEN="eyJhbGc..."

# Create project
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Demo","description":"Test project"}'

Runtime: Docker

Root Directory: task-management (nếu code trong thư mục con)

Thêm Environment Variables:

text
APP_ENV=production
JWT_SECRET=<random-32-chars>
JWT_EXPIRED_HOURS=24
DATABASE_URL=<postgres-internal-url>
REDIS_URL=<redis-internal-url>
Click Create Web Service

Deploy lên Fly.io
bash
fly launch
fly secrets set \
  APP_ENV=production \
  JWT_SECRET=<random> \
  DATABASE_URL=<postgres-url> \
  REDIS_URL=<redis-url>
fly deploy
Docker image size
bash
docker build -t task-api .
docker images task-api
# → ~15MB (multi-stage build)
📁 Cấu trúc thư mục
text
task-management/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── cache/                      # Redis cache layer
│   ├── config/                     # Config loader
│   ├── database/                   # DB connections
│   ├── handlers/                   # HTTP handlers
│   ├── middlewares/                # Gin middlewares
│   ├── models/                     # GORM models
│   ├── ratelimit/                  # Rate limiter
│   ├── repositories/               # Data access
│   ├── routes/                     # Route definitions
│   └── services/                   # Business logic
├── pkg/
│   ├── hash/                       # bcrypt wrapper
│   ├── jwt/                        # JWT manager
│   ├── response/                   # Response helpers
│   └── validator/                  # Validation
├── .env.example                    # Env template
├── .golangci.yml                   # Lint config
├── Dockerfile                      # Multi-stage build
├── docker-compose.yml              # Local stack
├── go.mod
├── go.sum
└── README.md
🔑 Environment Variables
Biến	Mô tả	Default
APP_PORT	Port server	8080
APP_ENV	Môi trường	development
DATABASE_URL	Postgres connection string (Render/Fly.io)	-
REDIS_URL	Redis connection string (Render/Fly.io)	-
DB_HOST	Postgres host (local dev)	localhost
DB_PORT	Postgres port	5432
DB_USER	Postgres user	postgres
DB_PASSWORD	Postgres password	postgres
DB_NAME	Database name	taskdb
DB_SSLMODE	SSL mode	disable
REDIS_HOST	Redis host	localhost
REDIS_PORT	Redis port	6379
REDIS_PASSWORD	Redis password	``
REDIS_DB	Redis DB index	0
JWT_SECRET	Bắt buộc đổi trong production	-
JWT_EXPIRED_HOURS	Token TTL	24
⚠️ Production: luôn dùng DATABASE_URL và REDIS_URL thay vì các biến riêng lẻ.

🎯 Roadmap
☑ CRUD User, Project, Task, Comment
☑ JWT Authentication + bcrypt
☑ Redis caching + invalidation
☑ Rate limiting (token bucket)
☑ Docker multi-stage + docker-compose
☑ CI/CD với GitHub Actions
☑ Deploy lên Render
□ Swagger documentation
□ Prometheus metrics
□ WebSocket notifications
□ Background jobs (email)
🤝 Contributing
Fork repo

Tạo branch: git checkout -b feature/awesome-feature

Commit: git commit -m "Add awesome feature"

Push: git push origin feature/awesome-feature

Tạo Pull Request

📄 License
MIT License — xem file LICENSE để biết chi tiết.

👤 Author
Your Name

GitHub: @yourname

Email: your.email@example.com

🙏 Acknowledgments
Gin Web Framework

GORM

go-redis

golang-jwt