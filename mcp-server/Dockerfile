# Build stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o passgen-mcp .

# Final image
FROM alpine:3.23

RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=appuser:appgroup /app/passgen-mcp /app/passgen-mcp

USER appuser

# Set environment variable defaults
ENV PASSGEN_URL=http://localhost:8080
ENV MCP_PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/passgen-mcp"]
CMD ["-http", "-port", "8080"]