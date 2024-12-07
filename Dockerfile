# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/gohex cmd/api/main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/bin/gohex .
COPY --from=builder /app/config/config.yaml ./config/

# Create non-root user
RUN adduser -D -g '' gohex
USER gohex

# Expose port
EXPOSE 8080

# Set environment variables
ENV GOHEX_APP_ENVIRONMENT=production

# Command to run the application
ENTRYPOINT ["./gohex"]
CMD ["-config", "config/config.yaml"] 