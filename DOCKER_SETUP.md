# Docker Setup Guide

This application uses Docker and Docker Compose for containerization and orchestration.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

## Quick Start

### 1. Build and Run with Docker Compose

```bash
# Navigate to the project directory
cd d:\xampp\htdocs\app-master

# Start the application and PostgreSQL database
docker-compose up -d

# View logs
docker-compose logs -f app
```

### 2. Access the Application

- **Application**: http://localhost:3000
- **PostgreSQL**: localhost:5432

### Database Credentials

- **User**: itio_user
- **Password**: Admin@95555
- **Database**: itio_golang
- **Host**: postgres (from within Docker), localhost (from host machine)

## Build Configuration

### Dockerfile Details

The Dockerfile uses a **multi-stage build** strategy:

#### Stage 1: Builder
- Base Image: `golang:1.22.3-alpine`
- Installs build dependencies (gcc, musl-dev, pkg-config, libpq-dev)
- Downloads Go modules
- Compiles the Go application with CGO enabled for PostgreSQL support

#### Stage 2: Runtime
- Base Image: `alpine:latest` (minimal, ~130MB vs ~900MB)
- Only includes necessary runtime dependencies:
  - `ca-certificates` - SSL/TLS support
  - `postgresql-client` - Database client tools
  - `libpq` - PostgreSQL client library
  - `tzdata` - Timezone data
- Copies binary, views, and assets from builder stage

### Key Features

- **Health Checks**: Both app and PostgreSQL have health checks
- **Volume Mounts**: Allows hot-reload of views and assets during development
- **Environment Variables**: Supports external configuration via `.env` file
- **Network Isolation**: Services communicate via Docker network `app-network`
- **Optimized Size**: ~200MB final image (optimized with Alpine)

## Docker Compose Services

### app Service
- **Build**: Builds from local Dockerfile
- **Ports**: 3000:3000
- **Dependencies**: Waits for PostgreSQL to be healthy
- **Volumes**: Views, assets, and logs directories
- **Restart**: unless-stopped

### postgres Service
- **Image**: postgres:15-alpine
- **Ports**: 5432:5432
- **Data Persistence**: Named volume `postgres_data`
- **Performance Tuning**:
  - `max_connections=100`
  - `shared_buffers=256MB`
  - `effective_cache_size=1GB`

## Important Configurations

### Database Connection

The Go application connects to PostgreSQL using:
```
Host: postgres
Port: 5432
User: itio_user
Password: Admin@95555
Database: itio_golang
```

### Static Files

- `/views` - HTML templates
- `/assets` - CSS, JavaScript, images

Both are mounted as volumes for development convenience.

## Common Commands

```bash
# Start services in background
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f app
docker-compose logs -f postgres

# Rebuild images
docker-compose build --no-cache

# Execute command in container
docker exec -it crypto-payment-app ./main
docker exec -it crypto-payment-db psql -U itio_user -d itio_golang

# Check service status
docker-compose ps

# View specific service logs
docker-compose logs -f app --tail 100
```

## Environment Variables

Create/update `.env` file with:

```env
PORT=3000
DB_USER=itio_user
DB_PASSWORD=Admin@95555
DB_NAME=itio_golang
CommonURL=http://localhost:3000
FileURL=http://localhost:3000
GOOGLE_CLIENT_ID=your_value
GOOGLE_CLIENT_SECRET=your_value
FACEBOOK_APP_ID=your_value
FACEBOOK_APP_SECRET=your_value
LINKEDIN_CLIENT_ID=your_value
LINKEDIN_CLIENT_SECRET=your_value
SMTPusername=your_email
SMTPpassword=your_password
SMTPhost=smtppro.zoho.in
SMTPport=587
```

## Dependencies Included

### Go Modules
- **Fiber v2**: Modern, fast web framework
- **GORM**: Object-Relational Mapping
- **PostgreSQL Driver**: jackc/pgx
- **Crypto Libraries**:
  - ethereum/go-ethereum
  - btcsuite/btcd
  - skip2/go-qrcode
- **Authentication**: 
  - OAuth2 (Google, Facebook, LinkedIn)
  - JWT (golang-jwt)
  - OTP (pquerna/otp)
- **Data Processing**:
  - PDF Generation (gofpdf)
  - Excel (xuri/excelize)
  - Charts (go-chart, go-echarts)
  - Email (SMTP)

### System Dependencies
- Build: gcc, build-base, musl-dev, pkg-config, libpq-dev
- Runtime: ca-certificates, postgresql-client, libpq, tzdata

## Troubleshooting

### Container won't start
```bash
# Check logs
docker-compose logs app

# Rebuild from scratch
docker-compose build --no-cache
docker-compose up -d
```

### Database connection issues
```bash
# Verify PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Test connection
docker exec -it crypto-payment-db psql -U itio_user -d itio_golang
```

### Port already in use
Modify `docker-compose.yml` to use different ports:
```yaml
ports:
  - "3001:3000"  # Change 3001 to desired port
```

### Health check failures
Wait longer for service startup:
```bash
# Monitor startup
docker-compose logs -f app
```

## Performance Optimization

### For Production

1. Use specific tag for images (not `latest`)
2. Set proper resource limits in `docker-compose.yml`:
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

3. Enable log rotation:
```yaml
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Cleanup

```bash
# Remove containers
docker-compose down

# Remove containers and volumes (WARNING: deletes database data)
docker-compose down -v

# Remove unused images
docker image prune -a

# Remove all unused resources
docker system prune -a
```

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/compose-file/)
- [Fiber Framework](https://gofiber.io/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
