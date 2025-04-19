# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Run stage
FROM alpine:3.18

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata curl bash postgresql-client

# Install migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/internal/dataaccess/database/migration ./migration

# Copy startup script
COPY deployments/start.sh /app/start.sh
RUN chmod +x /app/start.sh

# Expose the application port
EXPOSE 8081

# Command to run the application
CMD ["/app/start.sh"]
