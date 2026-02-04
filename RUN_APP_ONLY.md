# 🚀 Hướng dẫn chạy riêng Payment Gateway App

> **Trường hợp**: PostgreSQL và Kafka đã chạy sẵn bằng Docker trên cùng server

## 📋 Yêu cầu

- PostgreSQL container đang chạy với tên: `payment-postgres`
- Kafka container đang chạy với tên: `broker`
- Hai container này phải trong cùng Docker network

---

## ✅ Bước 1: Kiểm tra các service đã chạy

```bash
# Kiểm tra PostgreSQL
docker ps | grep payment-postgres

# Kiểm tra Kafka
docker ps | grep broker

# Kiểm tra network name (thường là <foldername>_default)
docker network ls | grep go-payment-gateway
```

**Output mẫu:**
```
CONTAINER ID   IMAGE                      STATUS         NAMES
abc123def456   postgres:16.11-alpine...   Up 2 hours     payment-postgres
def789ghi012   apache/kafka:3.7.0         Up 2 hours     broker

NETWORK ID     NAME                         DRIVER
xyz789abc456   go-payment-gateway_default   bridge
```

---

## ✅ Bước 2: Cấu hình file .env

Tạo hoặc chỉnh sửa file `.env`:

```bash
cp .env.example .env
nano .env  # hoặc vim .env
```

**Nội dung file `.env`:**
```env
# Database Configuration
# Sử dụng container name và internal port
DB_HOST=payment-postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123456
DB_NAME=payment_db
DB_SSLMODE=disable

# Kafka Configuration
# Sử dụng container name và internal port (29092 là internal listener)
KAFKA_BROKER=broker:29092
KAFKA_PRODUCER_TOPIC=transaction_create
KAFKA_CONSUMER_TOPIC=consumer_topic
KAFKA_CONSUMER_GROUP_ID=consumer_group
```

---

## ✅ Bước 3: Build và chạy app

### Cách 1: Sử dụng docker-compose.app-only.yml (Recommended)

```bash
# Build image
docker-compose -f docker-compose.app-only.yml build

# Start app
docker-compose -f docker-compose.app-only.yml up -d

# Xem logs
docker-compose -f docker-compose.app-only.yml logs -f payment-gateway
```

### Cách 2: Sử dụng Docker commands trực tiếp

```bash
# Build image
docker build -t payment-gateway:latest .

# Lấy network name
NETWORK_NAME=$(docker inspect payment-postgres --format='{{range $k, $v := .NetworkSettings.Networks}}{{$k}}{{end}}')

# Run container và join vào cùng network
docker run -d \
  --name payment-gateway \
  --network $NETWORK_NAME \
  -p 8085:8085 \
  --env-file .env \
  --restart unless-stopped \
  payment-gateway:latest
```

---

## ✅ Bước 4: Kiểm tra kết nối

### 1. Check logs để xem app đã khởi động thành công

```bash
docker logs -f payment-gateway
```

**Logs mong đợi:**
```
[GIN-debug] Listening and serving HTTP on :8085
Successfully connected to database
Kafka consumer started...
```

### 2. Kiểm tra kết nối database

```bash
docker exec -it payment-gateway sh -c "wget -qO- http://localhost:8085/api/transactions"
```

### 3. Test API từ bên ngoài

```bash
# Từ máy local hoặc bất kỳ đâu
curl -X POST http://YOUR_SERVER_IP:8085/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100000,
    "currency": "VND",
    "merchant_id": "MERCHANT_001"
  }'
```

---

## 🔧 Troubleshooting

### ❌ Lỗi: "network go-payment-gateway_default not found"

**Nguyên nhân:** Network name không đúng

**Giải pháp:**

```bash
# Tìm network thực tế
docker network ls

# Update file docker-compose.app-only.yml
# Thay đổi network name trong section networks:
networks:
  go-payment-gateway_default:
    external: true
    name: <your_actual_network_name>  # Thay tên thực tế ở đây
```

### ❌ Lỗi: "dial tcp: lookup payment-postgres: no such host"

**Nguyên nhân:** App không thể resolve tên container PostgreSQL

**Giải pháp:**

```bash
# Kiểm tra tên container PostgreSQL
docker ps | grep postgres

# Update .env với tên container chính xác
DB_HOST=<actual_postgres_container_name>

# Restart app
docker-compose -f docker-compose.app-only.yml restart
```

### ❌ Lỗi: "connection refused" khi kết nối Kafka

**Nguyên nhân:** Sử dụng sai listener port

**Giải pháp:**

```bash
# Trong Docker network, phải dùng internal listener: broker:29092
# KHÔNG dùng localhost:9092 (đó là external listener)

# Update .env
KAFKA_BROKER=broker:29092

# Restart app
docker-compose -f docker-compose.app-only.yml restart
```

### ❌ Không thể gọi API từ bên ngoài

**Kiểm tra:**

```bash
# 1. Kiểm tra port đã được expose
docker ps | grep payment-gateway
# Phải thấy: 0.0.0.0:8085->8085/tcp

# 2. Kiểm tra firewall (nếu trên server)
sudo ufw status
sudo ufw allow 8085/tcp

# 3. Test từ local server trước
curl http://localhost:8085/api/transactions
```

---

## 🛑 Stop và Remove

```bash
# Stop app
docker-compose -f docker-compose.app-only.yml down

# Hoặc
docker stop payment-gateway
docker rm payment-gateway
```

---

## 📊 Monitoring

### Xem logs real-time
```bash
docker logs -f payment-gateway
```

### Xem resource usage
```bash
docker stats payment-gateway
```

### Xem network connections
```bash
docker exec payment-gateway netstat -an | grep ESTABLISHED
```

---

## 🔐 Security Notes

1. **Đổi password mặc định** trong production:
   ```env
   DB_PASSWORD=your_secure_password_here
   ```

2. **Không expose port không cần thiết**

3. **Sử dụng TLS/SSL** cho production:
   ```env
   DB_SSLMODE=require
   ```

4. **Giới hạn network access** bằng firewall rules

---

## 📞 Support

Nếu gặp vấn đề, check:
1. Logs của app: `docker logs payment-gateway`
2. Logs của PostgreSQL: `docker logs payment-postgres`
3. Logs của Kafka: `docker logs broker`
4. Network connectivity: `docker exec payment-gateway ping payment-postgres`

---

**✅ Hoàn tất!** App đã sẵn sàng nhận request từ bên ngoài qua `http://YOUR_SERVER_IP:8085`
