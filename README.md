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