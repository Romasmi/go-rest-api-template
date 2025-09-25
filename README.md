# Go REST API Template

A comprehensive template for building RESTful APIs in Go. This template includes everything you need to get started with building a production-ready API.

## Features

- **Routing**: Using [Chi](https://github.com/go-chi/chi) for flexible and fast HTTP routing
- **Authentication**: JWT-based authentication with middleware
- **Database**: PostgreSQL integration with connection pooling
- **Migrations**: Database migrations using [golang-migrate](https://github.com/golang-migrate/migrate)
- **CRUD Operations**: Example CRUD operations for a User entity
- **Validation**: Request validation using [validator](https://github.com/go-playground/validator)
- **Documentation**: API documentation using Swagger/OpenAPI
- **Configuration**: Environment-based configuration with sensible defaults
- **Logging**: Structured logging
- **Error Handling**: Consistent error handling and responses
- **Middleware**: Common middleware for request ID, real IP, logging, recovery, timeouts, and CORS
- **Graceful Shutdown**: Graceful shutdown of the HTTP server

## Project Structure

```
.
├── api/
│   └── swagger/         # Swagger/OpenAPI specifications
├── cmd/
│   └── api/             # Application entry points
│       └── main.go      # Main application
├── docs/                # Documentation
├── internal/            # Private application code
│   ├── config/          # Configuration package
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   └── services/        # Business logic
├── migrations/          # Database migrations
├── pkg/                 # Public packages
│   ├── logger/          # Logging package
│   └── validator/       # Validation package
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Docker (optional)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-rest-api-template.git
cd go-rest-api-template
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file:

```bash
cp .env.example .env
```

4. Update the `overried.ini` file with your configuration or specify environment variables. 
For example database -> host in yaml is DATABASE_HOST in env.  :.

### Running the API

1. Start PostgreSQL:

```bash
# Using Docker
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=go_rest_api -p 5432:5432 -d postgres
```

2. Run database migrations:

```bash
go run cmd/api/main.go migrate
```

3. Start the API:

```bash
go run cmd/api/main.go
```

4. Access the API at http://localhost:8080

5. Access the Swagger documentation at http://localhost:8080/swagger/

### API Endpoints

#### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login a user

#### Users

- `GET /api/v1/users` - List users (requires authentication)
- `GET /api/v1/users/{id}` - Get a user by ID (requires authentication)
- `PUT /api/v1/users/{id}` - Update a user (requires authentication)
- `DELETE /api/v1/users/{id}` - Delete a user (requires authentication)

#### Health Check

- `GET /health` - Health check endpoint

## Configuration

The application can be configured using environment variables. See the `.env.example` file for available options.

## Development

### Adding a New Entity

1. Create a new model in `internal/models/`
2. Create a new repository in `internal/repository/`
3. Create a new service in `internal/services/`
4. Create a new handler in `internal/handlers/`
5. Register the new handler in `cmd/api/main.go`

### Adding a New Migration

1. Create a new migration file in `migrations/`:

```bash
# Create a migration to add a new table
echo "CREATE TABLE IF NOT EXISTS items (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL);" > migrations/000002_create_items_table.up.sql
echo "DROP TABLE IF EXISTS items;" > migrations/000002_create_items_table.down.sql
```

2. Run the migration:

```bash
go run cmd/api/main.go migrate
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.