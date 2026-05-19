# Project Tracker API

Backend RESTful API for Project Registration and Tracking Platform.

## Stack

- **Go 1.22+**
- **Gin** (HTTP framework)
- **GORM** (ORM)
- **PostgreSQL 16**
- **Redis 7**
- **Docker & Docker Compose**

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.22+ (for local development)

### Run with Docker

```bash
# Start all services (PostgreSQL, Redis, API)
docker-compose up -d --build

# Stop all services
docker-compose down
```

The API will be available at `http://localhost:8080/api/v1`.

### Local Development

1. Copy `.env.example` to `.env` and adjust values if needed:
```bash
cp .env.example .env
```

2. Start PostgreSQL and Redis:
```bash
docker-compose up -d postgres redis
```

3. Run the application:
```bash
go run cmd/server/main.go
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | 8080 | HTTP server port |
| `SERVER_HOST` | 0.0.0.0 | HTTP server host |
| `GIN_MODE` | debug | Gin mode (debug/release) |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | PostgreSQL user |
| `DB_PASSWORD` | postgres | PostgreSQL password |
| `DB_NAME` | project_tracker | PostgreSQL database |
| `DB_SSLMODE` | disable | PostgreSQL SSL mode |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `JWT_SECRET` | default-secret | JWT signing secret |
| `JWT_EXPIRATION_HOURS` | 24 | JWT expiration time |
| `CORS_ALLOWED_ORIGINS` | http://localhost:3000 | Allowed CORS origins |
| `LOG_LEVEL` | debug | Log level |
| `LOG_FORMAT` | json | Log format (json/console) |

## API Endpoints

### Authentication
| Method | Endpoint | Auth |
|---|---|---|
| POST | `/api/v1/auth/register` | No |
| POST | `/api/v1/auth/login` | No |
| GET | `/api/v1/auth/me` | Yes |

### Users
| Method | Endpoint | Auth |
|---|---|---|
| GET | `/api/v1/users` | Yes (admin, manager) |
| GET | `/api/v1/users/:id` | Yes |
| PUT | `/api/v1/users/:id` | Yes |
| DELETE | `/api/v1/users/:id` | Yes (admin) |

### Projects
| Method | Endpoint | Auth |
|---|---|---|
| GET | `/api/v1/projects` | Yes |
| POST | `/api/v1/projects` | Yes |
| GET | `/api/v1/projects/:id` | Yes |
| PUT | `/api/v1/projects/:id` | Yes |
| DELETE | `/api/v1/projects/:id` | Yes |
| PATCH | `/api/v1/projects/:id/status` | Yes |

### Tasks
| Method | Endpoint | Auth |
|---|---|---|
| GET | `/api/v1/projects/:projectId/tasks` | Yes |
| POST | `/api/v1/projects/:projectId/tasks` | Yes |
| GET | `/api/v1/tasks/:id` | Yes |
| PUT | `/api/v1/tasks/:id` | Yes |
| DELETE | `/api/v1/tasks/:id` | Yes |
| PATCH | `/api/v1/tasks/:id/status` | Yes |

### Members
| Method | Endpoint | Auth |
|---|---|---|
| GET | `/api/v1/projects/:projectId/members` | Yes |
| POST | `/api/v1/projects/:projectId/members` | Yes |
| PATCH | `/api/v1/projects/:projectId/members/:userId` | Yes |
| DELETE | `/api/v1/projects/:projectId/members/:userId` | Yes |

### Reports
| Method | Endpoint | Auth |
|---|---|---|
| GET | `/api/v1/reports/dashboard` | Yes |
| GET | `/api/v1/reports/projects/:id` | Yes |
| GET | `/api/v1/reports/projects/by-status` | Yes |
| GET | `/api/v1/reports/tasks/by-status` | Yes |

## Makefile Commands

| Command | Description |
|---|---|
| `make run` | Run the application locally |
| `make build` | Build the binary |
| `make test` | Run unit tests |
| `make docker-up` | Start Docker services |
| `make docker-down` | Stop Docker services |
| `make fmt` | Format Go code |

## Project Structure

```
api/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/           # Configuration loading
│   ├── domain/
│   │   ├── model/        # GORM entities
│   │   ├── repository/   # Repository interfaces
│   │   └── service/      # Business logic
│   ├── infrastructure/
│   │   ├── database/     # DB & Redis connections, migrations
│   │   ├── repository/   # Repository implementations
│   │   └── middleware/   # Auth, CORS, Logger, Role, ProjectMember
│   └── transport/http/
│       ├── handler/      # HTTP handlers
│       ├── dto/          # Data Transfer Objects
│       └── router.go     # Route definitions
├── pkg/
│   ├── auth/             # JWT utilities
│   ├── hash/             # Bcrypt utilities
│   └── response/         # Standardized HTTP responses
├── docker-compose.yml    # Docker orchestration
├── Dockerfile            # Go application image
└── go.mod                # Go dependencies
```
