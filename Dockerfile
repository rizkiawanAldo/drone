FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install make and other build dependencies
RUN apk add --no-cache make git

# Configure Go environment for better network resilience
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Copy Makefile to ensure we can run make commands
COPY Makefile ./
COPY api.yaml ./

# Copy the rest of the source code
COPY . .

# Initialize project and generate code
RUN make init
RUN make generate

# Build the application
RUN make build

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies if needed
RUN apk add --no-cache make

# Copy the binary from builder stage
COPY --from=builder /app/server /app/

# Copy Makefile (for the run command) and any other needed files
COPY --from=builder /app/Makefile /app/

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["make", "run"]