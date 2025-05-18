# Plantation Management API

A Go API for plantation management and drone patrol planning.

## Project Structure

```
.
├── api.yaml            # OpenAPI v3 specification
├── cmd                 # Application entry points
│   └── server          # HTTP server
│       └── main.go     # Main server entry point
├── database.sql        # Database schema definition
├── docker-compose.yml  # Docker-compose configuration
├── Dockerfile          # Docker build file
├── generated           # Generated code from OpenAPI spec (created at build time)
├── go.mod              # Go modules file
├── go.sum              # Go modules checksums
├── internal            # Internal application code
│   ├── api             # API handlers
│   ├── config          # Configuration management
│   └── repository      # Data access layer
├── Makefile            # Build and development scripts
└── README.md           # Project documentation
```

## Requirements

- Go 1.22 or higher
- Docker & Docker Compose (or Colima)
- Make
- PostgreSQL 14 or higher

## Database Setup

### Using Docker Compose

The recommended way to set up the database is through Docker Compose:

```bash
docker compose up -d postgres
```

### Using Local PostgreSQL Installation

1. Start PostgreSQL service:

```bash
# For macOS with Homebrew
brew services start postgresql

# For Linux
sudo systemctl start postgresql
```

2. Create the database:

```bash
createdb plantation
```

3. Initialize the database schema:

```bash
psql -d plantation -f database.sql
```

### Database Configuration

The application uses the following environment variables for database configuration:

- `DB_HOST`: Database host (default: `localhost`)
- `DB_PORT`: Database port (default: `5432`)
- `DB_USER`: Database user (default: `postgres`)
- `DB_PASSWORD`: Database password (default: `postgres`)
- `DB_NAME`: Database name (default: `plantation`)
- `DB_SSLMODE`: SSL mode (default: `disable`)

### Connecting with pgAdmin

To connect to the database using pgAdmin:

1. Open pgAdmin in your browser
2. Right-click on "Servers" in the left panel and select "Create" > "Server..."
3. In the "General" tab, give your connection a name (e.g., "Plantation Project")
4. Switch to the "Connection" tab and enter these details:
   - Host: localhost (or the value of DB_HOST environment variable)
   - Port: 5432 (or the value of DB_PORT environment variable)
   - Maintenance database: plantation (or the value of DB_NAME environment variable)
   - Username: postgres (or the value of DB_USER environment variable)
   - Password: postgres (or the value of DB_PASSWORD environment variable)
   - SSL Mode: disable (or the value of DB_SSLMODE environment variable)

## Getting Started

1. Initialize the project:

```
make init
```

2. Build the project:

```
make build
```

3. Run the tests:

```
make test
```

4. Start the application:

```
make run
```

## Docker

1. Build the Docker image:

```
make docker-build
```

2. Run the application with Docker:

```
make docker-run
```

## API Endpoints

- `POST /estate` - Create a new estate
- `POST /estate/{id}/tree` - Add a tree to an estate
- `GET /estate/{id}/stats` - Get stats about trees in an estate
- `GET /estate/{id}/drone-plan` - Get drone monitoring travel plan

## License

[License information] 