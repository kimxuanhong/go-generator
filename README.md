# Go Project Generator API

API để generate Go project với Clean Architecture structure.

## Endpoints

### Health Check
```bash
GET /health
```

### Generate Project
```bash
POST /generate
```

## Request Body

```json
{
  "projectName": "string",        // Required: Tên project
  "moduleName": "string",         // Required: Module name (e.g., github.com/user/project)
  "framework": "string",          // Required: Framework (gin | fiber | echo)
  "libs": ["string"],             // Optional: List of libraries (redis | postgres | mysql | resty | cron | rabbitmq | kafka | activemq | mapstructure | validator | opentelemetry)
  "includeExample": boolean       // Optional: Include example code (default: false)
}
```

## Examples

### 1. Generate project without example code (Minimal)

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": ["redis", "postgres"]
  }' \
  --output my-project.zip
```

### 2. Generate project with example code

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": ["redis", "postgres"],
    "includeExample": true
  }' \
  --output my-project.zip
```

### 3. Generate project with Fiber framework

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "fiber",
    "libs": ["redis", "postgres"],
    "includeExample": true
  }' \
  --output my-project.zip
```

### 4. Generate project with only Redis

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": ["redis"],
    "includeExample": true
  }' \
  --output my-project.zip
```

### 5. Generate project with only Postgres

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "fiber",
    "libs": ["postgres"],
    "includeExample": true
  }' \
  --output my-project.zip
```

### 6. Generate minimal project (no libs, no example)

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": []
  }' \
  --output my-project.zip
```

### 7. Generate project with MySQL

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": ["mysql"],
    "includeExample": true
  }' \
  --output my-project.zip
```

### 8. Generate project with Resty

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": ["resty"],
    "includeExample": true
  }' \
  --output my-project.zip
```

### 9. Generate project with all libraries

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": ["redis", "postgres", "mysql", "resty", "cron", "rabbitmq", "kafka"],
    "includeExample": true
  }' \
  --output my-project.zip
```

### 10. Generate project with Cron

```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-project",
    "moduleName": "github.com/user/my-project",
    "framework": "gin",
    "libs": ["cron"],
    "includeExample": true
  }' \
  --output my-project.zip
```

## Response

### Success
- **Status Code**: 200 OK
- **Content-Type**: application/zip
- **Body**: ZIP file containing the generated project

### Error
- **Status Code**: 400 Bad Request (validation error) or 500 Internal Server Error
- **Content-Type**: text/plain
- **Body**: Error message

## Project Structure

### Without Example Code
```
project-name/
├── cmd/
│   └── main.go
├── internal/
│   ├── app/
│   │   └── server.go          # Simple server with health check
│   ├── deps/
│   │   ├── config.go
│   │   └── deps.go
│   └── infrastructure/
│       ├── redis/
│       └── postgres/
└── config/
    └── config.json
```

### With Example Code (includeExample: true)
```
project-name/
├── cmd/
│   └── main.go
├── internal/
│   ├── app/
│   │   └── server.go          # Server with example routes
│   ├── domain/
│   │   └── entity.go          # User entity and repository interfaces
│   ├── usecase/
│   │   └── user_usecase.go    # Business logic
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── cache_repository.go
│   ├── handler/
│   │   └── user_handler.go    # HTTP handlers
│   ├── deps/
│   │   ├── config.go
│   │   └── deps.go
│   └── infrastructure/
│       ├── redis/
│       └── postgres/
└── config/
    └── config.json
```

## Frameworks

### Gin
- Framework: `gin`
- Config: `gin` section in config.json
- Port: 8080 (default)

### Fiber
- Framework: `fiber`
- Config: `fiber` section in config.json
- Port: 8080 (default)

### Echo
- Framework: `echo`
- Config: `echo` section in config.json
- Debug: true/false
- Port: 8080 (default)

## Libraries

### Redis
- Library: `redis`
- Config: `redis` section in config.json
- Default: `localhost:6379`

### Postgres
- Library: `postgres`
- Config: `postgres` section in config.json
- Uses GORM for database operations
- Default: `postgres://postgres:postgres@localhost:5432/example`

### MySQL
- Library: `mysql`
- Config: `mysql` section in config.json
- Uses GORM for database operations
- Default: `user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local`

### Logging (Built-in)
- **Logrus is built-in** - automatically included in all generated projects
- Config: `log` section in config.json
- Configurable: log level (debug, info, warn, error, fatal, panic)
- Default: JSON formatter with info level
- Cannot be disabled

### Resty
- Library: `resty`
- Config: `resty` section in config.json
- HTTP client library
- Configurable: baseURL, timeout, retryCount, retryWaitTime, debug, proxyURL, followRedirect

### Cron (Quartz-style)
- Library: `cron`
- Config: `cron` section in config.json
- Cron scheduler using github.com/robfig/cron/v3
- Configurable: location (timezone), jobs (name, schedule, description, enabled)
- Supports standard cron expressions (e.g., "0 */5 * * * *" - every 5 minutes)
- Example job included when `includeExample: true`

## Quick Start

1. Start the server:
```bash
go run main.go
```

2. Generate a project:
```bash
curl -X POST "http://localhost:8080/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "my-api",
    "moduleName": "github.com/user/my-api",
    "framework": "gin",
    "libs": ["redis", "postgres"],
    "includeExample": true
  }' \
  --output my-api.zip
```

3. Extract and run:
```bash
unzip my-api.zip
cd my-api
go mod tidy
go run cmd/main.go
```

## Testing

Run the test script:
```bash
chmod +x test-examples.sh
./test-examples.sh
```

## Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "service": "go-project-generator"
}
```

