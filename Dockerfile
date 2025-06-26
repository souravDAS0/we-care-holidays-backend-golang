# Builder stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./


# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Build the seeder binary (optional)
RUN CGO_ENABLED=0 GOOS=linux go build -o wecare-holidays-seeder ./cmd/seeder

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/wecare-holidays-seeder .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/.env ./.env

# Copy seed data (optional, ensure it exists)
COPY --from=builder /app/internal/seeder/data ./internal/seeder/data


# Set execution permissions
RUN chmod +x main wecare-holidays-seeder

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]