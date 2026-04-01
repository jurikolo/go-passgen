# go-passgen
Sample Go app to try out modern AI features

## Overview

A simple password generator web application written in Go. It provides a web interface and API for generating secure passwords.

## Features

- Web interface for generating passwords
- REST API endpoint (`/api/generate`)
- Health check endpoint (`/health`)
- Configurable port via `PORT` environment variable (default: 8080)

## Building and Running

### Using Go directly

```bash
go build -o passgen .
./passgen
```

### Using Makefile

The project includes a Makefile with convenient targets:

```bash
make build      # Build the Go binary
make test       # Run all tests
make docker-build  # Build Docker image
make docker-run    # Run Docker container locally
make docker-push   # Push Docker image to registry
make clean      # Clean up built artifacts
make all        # Run tests, build binary, and build Docker image
make local      # Build and run locally
```

### Using Docker

Build the Docker image:

```bash
docker build -t go-passgen .
```

Run the container:

```bash
docker run -p 8080:8080 --name passgen go-passgen
```

## Docker Image Details

The application uses a multi-stage Docker build:

1. **Build stage**: Uses `golang:1.26-alpine` to compile the application and run tests
2. **Final stage**: Uses `alpine:3.23` with a non-root user for security
3. **Port**: Exposes port 8080
4. **Health check**: Includes a health check endpoint

## Deployment on AWS EC2 with Podman and Traefik

This section explains how to deploy the go-passgen application on an AWS EC2 Ubuntu Linux server that already runs Podman and Traefik with proxy-network enabled.

### Prerequisites

- AWS EC2 Ubuntu instance (22.04 LTS or later)
- Podman installed and configured
- Traefik reverse proxy running with proxy-network enabled
- Docker Hub or other container registry access

### Step 1: Pull the Docker Image

If you have built and pushed the image to a registry:

```bash
podman pull jurikolo/go-passgen:latest
```

Or build locally on the EC2 instance:

```bash
git clone https://github.com/jurikolo/go-passgen.git
cd go-passgen
podman build -t go-passgen .
```

### Step 2: Create a Podman Pod (Optional but Recommended)

Create a pod to group your containers and share network namespace:

```bash
podman pod create --name web-apps -p 8080:8080
```

### Step 3: Run the Container

Run the go-passgen container in the pod:

```bash
podman run -d \
  --pod web-apps \
  --name passgen \
  -e PORT=8080 \
  --restart unless-stopped \
  go-passgen:latest
```

Or without a pod:

```bash
podman run -d \
  --name passgen \
  -p 8080:8080 \
  -e PORT=8080 \
  --restart unless-stopped \
  go-passgen:latest
```

### Step 4: Configure Traefik

Assuming Traefik is already running with proxy-network enabled, you need to add labels to the container for Traefik to discover it. Update the run command with Traefik labels:

```bash
podman run -d \
  --name passgen \
  --network proxy-network \
  -l "traefik.enable=true" \
  -l "traefik.http.routers.passgen.rule=Host(\`passgen.yourdomain.com\`)" \
  -l "traefik.http.routers.passgen.entrypoints=websecure" \
  -l "traefik.http.services.passgen.loadbalancer.server.port=8080" \
  -e PORT=8080 \
  --restart unless-stopped \
  go-passgen:latest
```

### Step 5: Verify Deployment

Check if the container is running:

```bash
podman ps
```

Test the application:

```bash
curl http://localhost:8080/health
```

### Step 6: Running Alongside Other Containers

To run go-passgen alongside other containers on the same host:

1. Ensure each container uses unique host ports or uses Traefik routing
2. Use Podman pods for related services
3. Configure Traefik with different hostnames or path prefixes for each service

Example Traefik configuration for multiple services:

```bash
# Service 1: go-passgen
podman run -d \
  --name passgen \
  --network proxy-network \
  -l "traefik.enable=true" \
  -l "traefik.http.routers.passgen.rule=Host(\`passgen.example.com\`)" \
  -l "traefik.http.services.passgen.loadbalancer.server.port=8080" \
  go-passgen:latest

# Service 2: Another app
podman run -d \
  --name another-app \
  --network proxy-network \
  -l "traefik.enable=true" \
  -l "traefik.http.routers.anotherapp.rule=Host(\`app.example.com\`)" \
  -l "traefik.http.services.anotherapp.loadbalancer.server.port=3000" \
  another-app:latest
```

### Step 7: Persistence and Updates

For production deployment, consider:

1. **Volume mounts** for any persistent data
2. **Environment variables** for configuration
3. **Automated updates** using watchtower or systemd timers
4. **Logging** with journald or external log aggregation

### Troubleshooting

- **Container not starting**: Check logs with `podman logs passgen`
- **Traefik not routing**: Verify container labels and Traefik configuration
- **Port conflicts**: Ensure no other service is using port 8080
- **Network issues**: Confirm containers are on the same network (`proxy-network`)

## Development

### Project Structure

```
.
├── main.go              # Application entry point
├── generator/           # Password generation logic
├── handler/             # HTTP handlers and templates
├── Dockerfile          # Multi-stage Docker build
├── Makefile            # Build automation
└── README.md           # This file
```

### Running Tests

```bash
go test ./...
```

### API Endpoints

- `GET /` - Web interface
- `POST /api/generate` - Generate a password (accepts `length` and `complexity` query parameters)
- `GET /health` - Health check endpoint
