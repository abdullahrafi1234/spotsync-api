# SpotSync API 🚗⚡

Smart Parking & EV Charging Reservation System — built with Go, Echo, GORM, and PostgreSQL.

## 🔗 Live URL

https://spotsync-api-896o.onrender.com/

## ✨ Features

- JWT authentication with role-based access control (driver/admin)
- Parking zone management with dynamic available-spot calculation
- Concurrency-safe reservation system using database row-locking
- Clean Architecture (DTO → Handler → Service → Repository → Model)
- Centralized error handling

## 🛠️ Tech Stack

- Go 1.26
- Echo (web framework)
- GORM + PostgreSQL (NeonDB)
- JWT (golang-jwt/jwt/v5), bcrypt
- go-playground/validator

## 🏛️ Architecture

Handler → Service → Repository → Model

- **Handler**: binds & validates requests, extracts JWT claims, returns JSON
- **Service**: business logic (password hashing, JWT generation, capacity rules)
- **Repository**: all GORM database operations, including the row-locking transaction
- **Model**: GORM structs mapped to database tables

Dependencies are wired manually in `main.go` (Repository → Service → Handler).

## ⚙️ Setup (Local)

1. Clone the repo
   \`\`\`bash
   git clone https://github.com/abdullahrafi1234/spotsync-api.git
   cd spotsync-api
   \`\`\`
2. Install dependencies
   \`\`\`bash
   go mod tidy
   \`\`\`
3. Create a `.env` file:
   \`\`\`env
   DATABASE_URL=postgresql://user:pass@host/dbname?sslmode=require
   JWT_SECRET=your-secret-key
   PORT=8080
   \`\`\`
4. Run
   \`\`\`bash
   go run main.go
   \`\`\`

## 📡 API Endpoints

| Method | Endpoint                             | Access                            |
| ------ | ------------------------------------ | --------------------------------- |
| POST   | /api/v1/auth/register                | Public                            |
| POST   | /api/v1/auth/login                   | Public                            |
| GET    | /api/v1/zones                        | Public                            |
| GET    | /api/v1/zones/:id                    | Public                            |
| POST   | /api/v1/zones                        | Admin                             |
| POST   | /api/v1/reservations                 | Authenticated                     |
| GET    | /api/v1/reservations/my-reservations | Authenticated                     |
| DELETE | /api/v1/reservations/:id             | Authenticated (own) / Admin (any) |
| GET    | /api/v1/reservations                 | Admin                             |

## 🔒 Concurrency Handling

Reservation creation uses a GORM transaction with `SELECT ... FOR UPDATE` row-locking on the parking zone to prevent overbooking when multiple users reserve the last available spot simultaneously.
