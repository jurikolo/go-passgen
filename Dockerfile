# Multi-stage Dockerfile for go-passgen
# Stage 1: Build and test
FROM golang:1.26-alpine AS builder

# Install git (required for go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum (if exists) to download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Run tests before building
RUN go test ./...

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o passgen .

# Stage 2: Final minimal image
FROM alpine:3.23

# Create a non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder --chown=appuser:appgroup /app/passgen /app/passgen

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Health check (optional but good practice)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the application
CMD ["/app/passgen"]