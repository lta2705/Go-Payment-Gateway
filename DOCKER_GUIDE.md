# 🐳 Hướng dẫn Docker cho Go Payment Gateway

## 📋 Tổng quan

Ứng dụng và tất cả services (PostgreSQL, Kafka, AKHQ) đều chạy trong Docker containers. Có 2 cách triển khai:

1. **Full Stack** - Một docker-compose chạy tất cả (khuyến nghị)
2. **App Only** - App riêng biệt, kết nối tới services Docker đã chạy sẵn

---

## 🚀 Option 1: Full Stack (Khuyến nghị)

### Chạy tất cả cùng lúc

```bash
# Build và start tất cả services
docker-compose up -d

# Kiểm tra status
docker-compose ps

# Xem logs
docker-compose logs -f payment-gateway

# Stop tất cả
docker-compose down
```

### Services được expose:

| Service | Container Name | Internal Address | External Access |
|---------|---------------|------------------|-----------------|
| Payment Gateway | payment-gateway | payment-gateway:8085 | http://localhost:8085 |
| PostgreSQL | payment-postgres | payment-postgres:5432 | localhost:5433 |
| Kafka | broker | broker:29092 | localhost:9092 |
| AKHQ | payment-akhq | payment-akhq:8080 | http://localhost:8080 |

### Kiến trúc kết nối:

```
┌─────────────────────────────────────────────────┐
│         Docker Network (default)                │
│                                                 │
│  ┌──────────────┐      ┌──────────────┐       │
│  │   Payment    │─────▶│  PostgreSQL  │       │
│  │   Gateway    │      │payment-postgres      │
│  │              │      │   :5432      │       │
│  │   :8085      │      └──────────────┘       │
│  └──────────────┘                              │
│         │                                      │
│         │                                      │
│         │         ┌──────────────┐            │
│         └────────▶│    Kafka     │            │
│                   │   broker     │            │
│                   │   :29092     │            │
│                   └──────────────┘            │
│                                               │
└─────────────────────────────────────────────────┘
         │                │
         │ :8085          │ :9092, :5433
         ▼                ▼
    Host Machine     Host Machine
   (External API)  (External Access)
```

---

## 🔧 Option 2: App Only

Khi PostgreSQL và Kafka đã chạy sẵn trong Docker (từ docker-compose.yml hoặc manually).

### Bước 1: Chạy services trước (nếu chưa có)

```bash
# Chỉ chạy PostgreSQL và Kafka
docker-compose up -d postgres broker akhq

# Kiểm tra
docker-compose ps
```

### Bước 2: Chạy app kết nối vào cùng network

```bash
# App sẽ tự động kết nối vào network của docker-compose chính
docker-compose -f docker-compose.app-only.yml up -d

# Xem logs
docker-compose -f docker-compose.app-only.yml logs -f

# Stop app
docker-compose -f docker-compose.app-only.yml down
```

### ⚠️ Lưu ý quan trọng:

App sử dụng **external network** `go-payment-gateway_default` được tạo bởi docker-compose.yml chính.

**Network name format:** `<folder-name>_default`

Nếu folder project của bạn khác `Go-Payment-Gateway`, cần update network name:

```bash
# Kiểm tra network name hiện tại
docker network ls | grep payment

# Ví dụ output: my-payment_default
```

Sau đó update trong [docker-compose.app-only.yml](docker-compose.app-only.yml#L53-L55):

```yaml
networks:
  go-payment-gateway_default:  # ← Đổi thành network name thực tế
    external: true
```

---

## 🔑 Cấu hình kết nối

### Khi app chạy TRONG Docker network:

```bash
# Environment variables cho app container
DB_HOST=payment-postgres   # Container name
DB_PORT=5432              # Internal port
KAFKA_BROKER=broker:29092 # Container name + internal port
```

### Khi app chạy NGOÀI Docker (local development):

```bash
# Environment variables cho app chạy trực tiếp
DB_HOST=localhost
DB_PORT=5433              # Mapped port
KAFKA_BROKER=localhost:9092 # Mapped port
```

---

## 📝 Các lệnh thường dùng

### Build và Deploy

```bash
# Build lại image
docker-compose build payment-gateway

# Build không dùng cache
docker-compose build --no-cache payment-gateway

# Chạy với build mới
docker-compose up -d --build payment-gateway
```

### Debugging

```bash
# Xem logs realtime
docker-compose logs -f payment-gateway

# Xem logs 100 dòng cuối
docker-compose logs --tail 100 payment-gateway

# Shell vào container
docker exec -it payment-gateway sh

# Kiểm tra network connectivity từ container
docker exec payment-gateway ping payment-postgres
docker exec payment-gateway ping broker

# Kiểm tra environment variables
docker exec payment-gateway env | grep -E "DB_|KAFKA_"
```

### Network Management

```bash
# List networks
docker network ls

# Inspect network
docker network inspect go-payment-gateway_default

# Xem containers trong network
docker network inspect go-payment-gateway_default | grep Name
```

### Cleanup

```bash
# Stop tất cả
docker-compose down

# Stop và xóa volumes (⚠️ Mất data)
docker-compose down -v

# Xóa images cũ
docker image prune -a

# Xóa toàn bộ (⚠️ Careful)
docker system prune -a --volumes
```

---

## 🧪 Testing kết nối

### 1. Kiểm tra app health

```bash
curl http://localhost:8085/health
```

### 2. Test database connection

```bash
# Từ host
psql -h localhost -p 5433 -U postgres -d payment_db

# Từ app container
docker exec -it payment-gateway sh
# Trong container (nếu có psql):
psql -h payment-postgres -p 5432 -U postgres -d payment_db
```

### 3. Test Kafka connection

```bash
# List topics từ host
docker exec broker opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 --list

# Từ app container
docker exec payment-gateway wget -qO- http://broker:29092
```

---

## 🐛 Troubleshooting

### ❌ Error: "network go-payment-gateway_default declared as external, but could not be found"

**Nguyên nhân:** Network chưa được tạo

**Giải pháp:**
```bash
# Option 1: Chạy docker-compose.yml chính trước
docker-compose up -d postgres broker

# Option 2: Tạo network manually
docker network create go-payment-gateway_default

# Option 3: Kiểm tra tên folder và update network name
pwd  # Xem tên folder hiện tại
docker network ls | grep default
```

### ❌ Error: "dial tcp: lookup payment-postgres: no such host"

**Nguyên nhân:** Container không thể resolve DNS

**Giải pháp:**
```bash
# Kiểm tra containers có cùng network không
docker network inspect go-payment-gateway_default

# Đảm bảo postgres container đang chạy
docker-compose ps postgres

# Restart app container
docker-compose -f docker-compose.app-only.yml restart
```

### ❌ Error: Port 8085 already in use

**Giải pháp:**
```bash
# Option 1: Stop container đang dùng port
docker ps | grep 8085
docker stop <container-id>

# Option 2: Đổi port mapping
# Trong docker-compose.yml:
ports:
  - "8086:8085"  # Map sang port khác
```

### ❌ App restart liên tục

**Kiểm tra:**
```bash
# Xem logs
docker logs payment-gateway

# Kiểm tra health check
docker inspect payment-gateway | grep -A 20 State
```

**Nguyên nhân thường gặp:**
- Database chưa ready → Đợi healthcheck
- Wrong credentials → Kiểm tra .env
- Missing dependencies → Kiểm tra go.mod

---

## 📊 Production Checklist

- [ ] Sử dụng secrets thay vì .env file
- [ ] Configure resource limits (CPU, memory)
- [ ] Enable monitoring (Prometheus, Grafana)
- [ ] Configure log aggregation
- [ ] Set up backup cho volumes
- [ ] Use specific image tags (không dùng `:latest`)
- [ ] Configure restart policies
- [ ] Set up health checks cho tất cả services
- [ ] Use read-only filesystem khi có thể
- [ ] Configure security scanning

---

## 🔗 Quick Reference

```bash
# Start everything
docker-compose up -d

# Start only app (services already running)
docker-compose -f docker-compose.app-only.yml up -d

# View logs
docker-compose logs -f payment-gateway

# Restart app
docker-compose restart payment-gateway

# Stop everything
docker-compose down

# Rebuild and restart
docker-compose up -d --build

# Shell access
docker exec -it payment-gateway sh

# Check health
curl http://localhost:8085/health
```

---

## 📚 Tài liệu liên quan

- [Dockerfile](Dockerfile) - Multi-stage build configuration
- [docker-compose.yml](docker-compose.yml) - Full stack setup
- [docker-compose.app-only.yml](docker-compose.app-only.yml) - App only setup
- [.env.example](.env.example) - Environment variables template
