FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate code from OpenAPI spec
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
RUN mkdir -p generated
RUN oapi-codegen -package generated -generate types,server,spec api.yaml > generated/api.gen.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server cmd/server/main.go

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/server /app/server

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/server"] 