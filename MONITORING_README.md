# Monitoring Setup with Grafana

Proyek ini menggunakan stack monitoring lengkap dengan Prometheus, Grafana, dan Node Exporter untuk memantau performa sistem QRIS Payment.

## 🏗️ Arsitektur Monitoring

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   QRIS App      │    │ Metrics Server  │    │   Prometheus     │
│   (API:3000)    │    │ (Metrics:8080)  │───▶│   (Metrics)      │
│                 │    │                 │    │                 │
└─────────────────┘    │ /metrics        │    │ - App metrics   │
                       └─────────────────┘    │ - System metrics│
┌─────────────────┐                         │ - DB metrics    │
│  PostgreSQL     │    ┌─────────────────┐    └─────────────────┘
│  (Database)     │    │  Node Exporter  │             ▲
└─────────────────┘    │  (System:9100)  │             │
                       └─────────────────┘             │
┌─────────────────┐                                     │
│    Redis        │                                     │
│  (Cache)        │                                     │
└─────────────────┘                                     ▼
                                              ┌─────────────────┐
                                              │    Grafana       │
                                              │  (Dashboard)     │
                                              │  (Port:3001)     │
                                              └─────────────────┘
```

## 🚀 Quick Start

### Development Setup

```bash
# Jalankan semua services termasuk monitoring
docker-compose up -d

# Atau jalankan hanya monitoring services
docker-compose up -d prometheus grafana node-exporter metrics
```

### Production Setup

```bash
# Jalankan dengan production overrides
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Set environment variables untuk production
export GRAFANA_ADMIN_USER=admin
export GRAFANA_ADMIN_PASSWORD=secure_password_here
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 2. Akses Services

- **API Server**: http://localhost:3000
- **Metrics Server**: http://localhost:8080/metrics
- **Grafana**: http://localhost:3001 (admin/admin123)
- **Prometheus**: http://localhost:9090

### 3. Import Dashboard

1. Buka Grafana (http://localhost:3001)
2. Login dengan admin/admin123
3. Pilih menu "Dashboards" → "Import"
4. Upload file `monitoring/grafana/dashboards/qris-dashboard.json`
5. Dashboard akan otomatis ter-import dengan data source Prometheus

## 📊 Metrics yang Dikumpulkan

### Application Metrics (Auto-generated)
- HTTP request count per endpoint, method, status
- Response time histograms (p50, p95, p99)
- Active goroutines
- Memory usage (heap, stack)
- GC statistics

### System Metrics (Node Exporter)
- CPU usage per core
- Memory usage (RAM, swap)
- Disk I/O statistics
- Network interface traffic
- System load average
- Filesystem usage

### Database Metrics
- Active connections
- Query execution time
- Connection pool statistics
- Database size

### Cache Metrics (Redis)
- Memory usage
- Hit/miss ratios
- Connected clients
- Keyspace statistics

## 🔧 Konfigurasi

### Prometheus Configuration
File: `monitoring/prometheus.yml`

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'qris-app'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'
```

### Grafana Provisioning
- Datasources: `monitoring/grafana/provisioning/datasources/`
- Dashboards: `monitoring/grafana/provisioning/dashboards/`

## 📈 Custom Metrics

Untuk menambahkan custom business metrics:

```go
import "github.com/prometheus/client_golang/prometheus"

// Define custom metrics
var (
    qrisPaymentsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "qris_payments_total",
            Help: "Total number of QRIS payments processed",
        },
        []string{"status", "merchant_id", "payment_method"},
    )

    qrisPaymentAmount = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "qris_payment_amount",
            Help: "QRIS payment amount distribution",
            Buckets: prometheus.DefBuckets, // 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
        },
        []string{"currency"},
    )
)

// Register metrics (biasanya di init function)
func init() {
    prometheus.MustRegister(qrisPaymentsTotal)
    prometheus.MustRegister(qrisPaymentAmount)
}

// Use metrics in your code
func (u *PaymentUseCase) Pay(ctx context.Context, accountID string, req *model.PaymentRequest) (*model.PaymentResponse, error) {
    // ... existing code ...

    if payment.Status == "success" {
        qrisPaymentsTotal.WithLabelValues("success", req.MerchantID, req.PaymentMethod).Inc()
        qrisPaymentAmount.WithLabelValues("IDR").Observe(float64(req.Amount))
    } else {
        qrisPaymentsTotal.WithLabelValues("failed", req.MerchantID, req.PaymentMethod).Inc()
    }

    return resp, nil
}
```

## 🚨 Alerting (Optional)

Untuk menambahkan alerting, tambahkan AlertManager:

```yaml
# docker-compose.yml
alertmanager:
  image: prom/alertmanager:v0.26.0
  ports:
    - "9093:9093"
  volumes:
    - ./monitoring/alertmanager.yml:/etc/alertmanager/alertmanager.yml
```

## 📋 Dashboard Panels

Dashboard default mencakup:

1. **HTTP Requests** - Rate request per endpoint
2. **Response Time** - 95th percentile latency
3. **Database Connections** - Active PostgreSQL connections
4. **Redis Memory** - Cache memory usage
5. **System CPU** - CPU utilization
6. **System Memory** - RAM usage
7. **Disk I/O** - Storage performance
8. **Network Traffic** - Bandwidth usage

## 🔍 Troubleshooting

### Metrics tidak muncul
1. Pastikan metrics server running: `docker-compose ps`
2. Check metrics endpoint: `curl http://localhost:8080/metrics`
3. Verify Prometheus targets: http://localhost:9090/targets
4. Check metrics server logs: `docker-compose logs metrics`

### Grafana tidak bisa connect ke Prometheus
1. Check network: `docker network ls`
2. Verify service names di docker-compose.yml
3. Check Prometheus logs: `docker-compose logs prometheus`

### Dashboard kosong
1. Pastikan data source Prometheus configured
2. Check time range di dashboard
3. Verify metrics names match query

## 📚 Referensi

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Node Exporter](https://github.com/prometheus/node_exporter)
- [Fiber Prometheus Middleware](https://github.com/ansrivas/fiberprometheus)