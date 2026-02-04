# Docker Deployment Guide

## Tổng quan

Repository này cung cấp nhiều cách để containerize và chạy ứng dụng Payment Gateway:

1. **Full Stack** - Chạy tất cả services (app + postgres + kafka + akhq)
2. **App Only** - Chỉ chạy app, kết nối tới external services

## 📋 Yêu cầu

- Docker 20.10+
- Docker Compose 2.0+

## 🚀 Cách sử dụng

### 1. Full Stack (Tất cả services)

Chạy toàn bộ stack bao gồm PostgreSQL, Kafka, AKHQ và Payment Gateway:

```bash
# Build và chạy tất cả services
docker-compose up -d

# Xem logs
docker-compose logs -f payment-gateway

# Dừng tất cả services
docker-compose down

# Dừng và xóa volumes
docker-compose down -v
```

**Services được expose:**
- Payment Gateway API: `http://localhost:8085`
- PostgreSQL: `localhost:5433`
- Kafka: `localhost:9092`
- AKHQ (Kafka UI): `http://localhost:8080`

---

### 2. App Only (Kết nối external services)

Chỉ chạy app container, kết nối tới PostgreSQL và Kafka đã chạy sẵn bên ngoài.

#### Option A: Sử dụng Host Network Mode (Khuyến nghị cho local development)

```bash
# Build và chạy app với host network
docker-compose -f docker-compose.app-only.yml up -d

# Xem logs
docker-compose -f docker-compose.app-only.yml logs -f

# Dừng app
docker-compose -f docker-compose.app-only.yml down
```

**Ưu điểm:**
- App có thể kết nối trực tiếp tới services trên localhost
- Không cần cấu hình network phức tạp
- API có thể truy cập từ bên ngoài qua `http://localhost:8085`

**Lưu ý:** Cần đảm bảo PostgreSQL và Kafka đã chạy trước:
```bash
# Chạy chỉ PostgreSQL và Kafka từ docker-compose.yml gốc
docker-compose up -d postgres broker
```

#### Option B: Sử dụng Bridge Network (Production)

Nếu services bên ngoài cũng nằm trong Docker network:

```bash
# 1. Tạo external network (chỉ cần làm 1 lần)
docker network create payment-network

# 2. Chạy external services và connect vào network
docker run -d --name external-postgres \
  --network payment-network \
  -e POSTGRES_PASSWORD=123456 \
  -p 5433:5432 \
  postgres:16.11-alpine3.23

docker run -d --name external-kafka \
  --network payment-network \
  -p 9092:9092 \
  apache/kafka:3.7.0

# 3. Build và chạy app
docker build -t payment-gateway .

docker run -d --name payment-gateway \
  --network payment-network \
  -p 8085:8085 \
  -e DB_HOST=external-postgres \
  -e DB_PORT=5432 \
  -e KAFKA_BROKER=external-kafka:9092 \
  payment-gateway
```

#### Option C: Kết nối tới services trên máy host

Nếu services chạy ngoài Docker trên máy host:

**macOS:**
```bash
docker run -d --name payment-gateway \
  -p 8085:8085 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5433 \
  -e KAFKA_BROKER=host.docker.internal:9092 \
  payment-gateway
```

**Linux:**
```bash
docker run -d --name payment-gateway \
  --add-host=host.docker.internal:host-gateway \
  -p 8085:8085 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5433 \
  -e KAFKA_BROKER=host.docker.internal:9092 \
  payment-gateway
```

---

## 🔧 Cấu hình

### Environment Variables

Tạo file `.env` từ `.env.example`:

```bash
cp .env.example .env
```

Các biến quan trọng:

#### Database
- `DB_HOST` - Database host (default: `localhost` hoặc `postgres`)
- `DB_PORT` - Database port (default: `5433` hoặc `5432`)
- `DB_USER` - Database user (default: `postgres`)
- `DB_PASSWORD` - Database password (default: `123456`)
- `DB_NAME` - Database name (default: `payment_db`)

#### Kafka
- `KAFKA_BROKER` - Kafka broker address (default: `localhost:9092` hoặc `broker:29092`)
- `KAFKA_PRODUCER_TOPIC` - Producer topic (default: `transaction_create`)
- `KAFKA_CONSUMER_TOPIC` - Consumer topic (default: `consumer_topic`)
- `KAFKA_CONSUMER_GROUP_ID` - Consumer group ID (default: `consumer_group`)

---

## 🏗️ Build Manual

### Build Docker Image

```bash
# Build image
docker build -t payment-gateway:latest .

# Build với custom tag
docker build -t payment-gateway:v1.0.0 .

# Build với build args
docker build \
  --build-arg GO_VERSION=1.24 \
  -t payment-gateway:latest .
```

### Run Container Manual

```bash
# Run với default config
docker run -d \
  --name payment-gateway \
  -p 8085:8085 \
  payment-gateway:latest

# Run với custom environment
docker run -d \
  --name payment-gateway \
  -p 8085:8085 \
  -e DB_HOST=postgres-server \
  -e DB_PORT=5432 \
  -e DB_PASSWORD=secure_password \
  -e KAFKA_BROKER=kafka-server:9092 \
  payment-gateway:latest

# Run với env file
docker run -d \
  --name payment-gateway \
  -p 8085:8085 \
  --env-file .env \
  payment-gateway:latest
```

---

## 🧪 Testing

### Health Check

```bash
# Kiểm tra health status
curl http://localhost:8085/health

# Hoặc sử dụng docker
docker exec payment-gateway wget -qO- http://localhost:8085/health
```

### Logs

```bash
# View logs (docker-compose)
docker-compose logs -f payment-gateway

# View logs (manual docker)
docker logs -f payment-gateway

# View last 100 lines
docker logs --tail 100 payment-gateway
```

---

## 🛠️ Troubleshooting

### App không kết nối được tới Database/Kafka

1. **Kiểm tra network connectivity:**
```bash
# Từ container, ping tới service
docker exec payment-gateway ping postgres
docker exec payment-gateway ping broker

# Hoặc check bằng telnet
docker exec payment-gateway telnet postgres 5432
docker exec payment-gateway telnet broker 9092
```

2. **Kiểm tra environment variables:**
```bash
docker exec payment-gateway env | grep -E "DB_|KAFKA_"
```

3. **Kiểm tra services external có đang chạy:**
```bash
# Check PostgreSQL
docker-compose ps postgres
psql -h localhost -p 5433 -U postgres -d payment_db

# Check Kafka
docker-compose ps broker
docker exec broker opt/kafka/bin/kafka-topics.sh --bootstrap-server broker:29092 --list
```

### Port conflict

Nếu port 8085 đã được sử dụng:

```bash
# Thay đổi port mapping
docker run -d -p 8086:8085 payment-gateway

# Hoặc trong docker-compose.yml
ports:
  - "8086:8085"
```

### Container restart liên tục

```bash
# Xem logs để debug
docker logs payment-gateway

# Kiểm tra health check
docker inspect payment-gateway | grep -A 10 Health
```

---

## 📊 Production Considerations

### 1. Resource Limits

Thêm resource limits vào docker-compose.yml:

```yaml
services:
  payment-gateway:
    # ... other config
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 2. Secrets Management

Không commit file `.env` chứa credentials. Sử dụng Docker secrets hoặc external secret management:

```bash
# Sử dụng Docker secrets
echo "secure_password" | docker secret create db_password -

docker service create \
  --name payment-gateway \
  --secret db_password \
  payment-gateway:latest
```

### 3. Monitoring

Thêm health check endpoint và monitoring tools:

```yaml
healthcheck:
  test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8085/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### 4. Logging

Configure log driver cho production:

```yaml
services:
  payment-gateway:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

---

## 🔗 Quick Commands Reference

```bash
# Full stack
docker-compose up -d                    # Start all services
docker-compose down                     # Stop all services
docker-compose down -v                  # Stop and remove volumes
docker-compose logs -f payment-gateway  # View app logs

# App only
docker-compose -f docker-compose.app-only.yml up -d    # Start app only
docker-compose -f docker-compose.app-only.yml down     # Stop app

# Manual operations
docker build -t payment-gateway .                      # Build image
docker run -d -p 8085:8085 payment-gateway            # Run container
docker logs -f payment-gateway                         # View logs
docker exec -it payment-gateway sh                     # Shell access
docker stop payment-gateway                            # Stop container
docker rm payment-gateway                              # Remove container
```

---

## 📝 Notes

- App lắng nghe trên port `8085`
- Default timezone: `Asia/Ho_Chi_Minh`
- Healthcheck endpoint: `/health` (cần implement nếu chưa có)
- Consumer worker chạy trong goroutine riêng
- Sử dụng multi-stage build để optimize image size
