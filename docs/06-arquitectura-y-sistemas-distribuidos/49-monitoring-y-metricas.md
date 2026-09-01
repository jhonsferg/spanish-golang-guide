# Capítulo 49: Monitoring y métricas

## Introducción

En sistemas en producción, no basta con escribir código correcto. Necesitas saber qué está pasando en tus aplicaciones: cuántas requests procesa por segundo, cuánta memoria consume, cuáles son los tiempos de respuesta, y detectar anomalías antes de que afecten a los usuarios. Esta es la misión de la **observabilidad** y el **monitoring**.

Go, con su excelente rendimiento y modelo de concurrencia, requiere herramientas sofisticadas para monitorear comportamientos complejos. Prometheus se ha convertido en el estándar de facto para recolección de métricas en arquitecturas modernas, especialmente en ecosistemas containerizados y Kubernetes.

Este capítulo te guía a través de:

- Los tres pilares de la observabilidad
- Prometheus y su arquitectura pull
- Instrumentación de aplicaciones Go
- Alertas y dashboards con Grafana
- PromQL para análisis avanzado
- Best practices y antipatterns comunes

---

## 49.1 Observabilidad: Los Tres Pilares

La observabilidad moderna se construye sobre tres componentes interdependientes: **logs**, **métricas** y **traces**. Juntos proporcionan una visión completa del comportamiento de tu aplicación.

### 49.1.1 Conceptos Fundamentales

**Logs (Registros):**

- Registros detallados de eventos discretos
- Incluyen contexto completo del evento
- Típicamente no agregables a escala
- Ideales para debugging

```go
// Logs: detalles de eventos específicos
log.Println("Usuario autenticado: user_id=123, ip=192.168.1.1, duration_ms=45")
```

**Métricas:**

- Mediciones numéricas con respecto al tiempo
- Agregables y queryables
- Bajo overhead incluso con alto volumen
- Ideales para alertas y dashboards

```go
// Métrica: agregable en tiempo, consulta histórica
requestsTotal.WithLabelValues("POST", "200").Inc()
requestDuration.WithLabelValues("POST").Observe(0.045)
```

**Traces (Trazas):**

- Seguimiento de una request a través de múltiples servicios
- Visibilidad de latencias distribuidas
- Contexto causal en sistemas distribuidos
- Ideales para debugging de aplicaciones microservicios

```
request: GET /api/users/123
  ├─ auth service (45ms)
  │  └─ cache lookup (2ms)
  ├─ database service (200ms)
  └─ response formatting (10ms)
Total: 255ms
```

### 49.1.2 Triángulo de Observabilidad

```
                  LOGS
                  /  \
                 /    \
                /      \
         ┌─────/────────\─────┐
         |                    |
         |    APLICACIÓN      |
         |                    |
         └────────┬───────────┘
                  |
           METRICS   TRACES
```

**Cómo se complementan:**

- Métrica sube → logs del evento relevante → trace del request
- Trace identifica servicio lento → métricas de ese servicio → logs detallados

### 49.1.3 Niveles de Observabilidad

**Nivel 1: Básico**

- Uptime/availability
- Error rate
- Request rate
- Response time (p50, p95, p99)

**Nivel 2: Intermedio**

- Métricas por endpoint
- Utilización de recursos (CPU, memoria)
- Queue depths
- Cache hit rates

**Nivel 3: Avanzado**

- Trazas distribuidas completas
- Profiling continuo
- Anomaly detection
- Correlación de eventos

---

## 49.2 Prometheus: Arquitectura y Conceptos

Prometheus es un sistema de monitoreo de series temporales diseñado para aplicaciones modernas. A diferencia de sistemas push tradicionales, Prometheus **pull** (tira) las métricas de tus aplicaciones.

### 49.2.1 Arquitectura Pull vs Push

**Modelo Pull (Prometheus):**

```
Prometheus
    ↓ (scrape every 15s)
  /metrics endpoint
    ↑
   Tu aplicación
```

Ventajas:

- Prometheus controla la frecuencia de scraping
- Más fácil escalar múltiples instancias de Prometheus
- Menos configuración en aplicaciones
- Mejor detección de aplicaciones down

**Modelo Push (Traditonal):**

```
Tu aplicación
    ↓ (push)
  Monitoring system
```

### 49.2.2 Series Temporales

Una métrica en Prometheus es una **serie temporal**: secuencia de valores timestamped.

```
nombre_metrica{label1="valor1", label2="valor2"} valor timestamp
```

Ejemplo real:

```
http_requests_total{method="GET", endpoint="/api/users", status="200"} 1523 1609459200000
http_requests_total{method="GET", endpoint="/api/users", status="200"} 1847 1609459260000
http_requests_total{method="GET", endpoint="/api/users", status="200"} 2102 1609459320000
```

### 49.2.3 Tipos de Métricas

**Counter (Contador):**

- Solo incrementa (nunca disminuye)
- Se reinicia en 0 al reiniciar la aplicación
- Uso: requests totales, errores totales, bytes procesados

```
requests_total{service="api"} 45230
errors_total{service="api"} 123
```

**Gauge (Indicador):**

- Puede subir o bajar
- Captura estado instantáneo
- Uso: memoria usada, conexiones activas, temperatura

```
memory_usage_bytes 536870912
active_connections 523
```

**Histogram (Histograma):**

- Distribuición de valores en buckets
- Genera automáticamente quantiles
- Uso: latencias, tamaños de payload

```
request_duration_seconds_bucket{le="0.1"} 100
request_duration_seconds_bucket{le="0.5"} 345
request_duration_seconds_bucket{le="1.0"} 450
request_duration_seconds_bucket{le="+Inf"} 500
request_duration_seconds_sum 1250.5
request_duration_seconds_count 500
```

**Summary (Resumen):**

- Similar a Histogram pero con quantiles calculados en cliente
- Uso: cuando necesitas exactitud en percentiles

```
request_duration_seconds{quantile="0.5"} 0.25
request_duration_seconds{quantile="0.9"} 0.85
request_duration_seconds_sum 1250.5
request_duration_seconds_count 500
```

### 49.2.4 Labels (Etiquetas)

Las etiquetas permiten dimensionar tus métricas. Son cruciales pero pueden causar problemas si se usan mal.

```go
// Buena etiquetización
requestDuration.WithLabelValues("GET", "users", "200").Observe(0.045)
// Dimensiones: método HTTP, endpoint, código de respuesta
// Cardinalidad predecible

// Mala etiquetización
requestDuration.WithLabelValues("GET", user_id, "200").Observe(0.045)
// Si hay millones de usuarios → cardinalidad explosiva → problemas
```

### 49.2.5 Flujo Completo en Prometheus

```
┌─────────────────────────────────────┐
│   Tu Aplicación Go                  │
│  (expone /metrics)                  │
└──────────────┬──────────────────────┘
               ↑ (scrape)
┌──────────────┴──────────────────────┐
│   Prometheus Server                 │
│  (time-series database)             │
└──────────────┬──────────────────────┘
               ↓ (query)
    ┌──────────┴──────────┐
    │                     │
    ↓                     ↓
┌─────────┐         ┌──────────┐
│ Grafana │         │ Alerting │
│         │         │ Engine   │
└─────────┘         └──────────┘
```

---

## 49.3 Instrumentación de Aplicaciones Go

Go proporciona el cliente oficial `prometheus/client_golang` que hace instrumentación sencilla y eficiente.

### 49.3.1 Instalación y Setup Básico

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

### 49.3.2 Counter: Registro de Eventos Totales

```go
package main

import (
 "net/http"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
 // promauto registra automáticamente
 requestsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_requests_total",
   Help: "Total de requests HTTP procesados",
  },
  []string{"method", "endpoint", "status"},
 )

 errorsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "errors_total",
   Help: "Total de errores",
  },
  []string{"type", "severity"},
 )
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
 // Simular procesamiento
 status := 200

 // Registrar métrica
 requestsTotal.WithLabelValues(r.Method, "/api/users",
  fmt.Sprint(status)).Inc()

 w.WriteHeader(status)
 w.Write([]byte("users list"))
}

func main() {
 http.HandleFunc("/api/users", handleUsers)

 // Exponemos endpoint de métricas
 http.Handle("/metrics", promhttp.Handler())

 http.ListenAndServe(":8080", nil)
}
```

**Curl para probar:**

```bash
# Request normal
curl http://localhost:8080/api/users

# Ver métricas
curl http://localhost:8080/metrics | grep http_requests_total
```

### 49.3.3 Gauge: Estado Instantáneo

```go
var (
 activeConnections = promauto.NewGaugeVec(
  prometheus.GaugeOpts{
   Name: "active_connections",
   Help: "Conexiones activas por servicio",
  },
  []string{"service"},
 )

 memoryUsageBytes = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "memory_usage_bytes",
   Help: "Memoria usada en bytes",
  },
 )
)

func monitorConnections() {
 // Incrementar
 activeConnections.WithLabelValues("database").Inc()

 // Decrementar
 activeConnections.WithLabelValues("database").Dec()

 // Establecer valor exacto
 memoryUsageBytes.Set(float64(getMemoryUsage()))
}
```

### 49.3.4 Histogram: Distribuciones y Latencias

```go
var (
 requestDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_duration_seconds",
   Help:    "Duración de requests HTTP",
   // Buckets: rangos de latencia para análisis
   Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
  },
  []string{"method", "endpoint"},
 )

 payloadSizeBytes = promauto.NewHistogram(
  prometheus.HistogramOpts{
   Name:    "http_payload_size_bytes",
   Help:    "Tamaño de payloads HTTP",
   Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100 a 1GB
  },
 )
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
 start := time.Now()

 // ... procesar ...

 duration := time.Since(start).Seconds()
 requestDuration.WithLabelValues(r.Method, "/api/data").Observe(duration)
 payloadSizeBytes.Observe(float64(len(responsePayload)))
}
```

**Buckets automáticos:**

```go
// Lineal: 0, 1, 2, 3, 4, 5
prometheus.LinearBuckets(0, 1, 5)

// Exponencial: 0.1, 1, 10, 100, 1000, 10000
prometheus.ExponentialBuckets(0.1, 10, 5)

// Personalizado
[]float64{.001, .01, .1, 1}
```

### 49.3.5 Summary: Quantiles Precisos

```go
var (
 dbQueryDuration = promauto.NewSummaryVec(
  prometheus.SummaryOpts{
   Name:       "db_query_duration_seconds",
   Help:       "Duración de queries a BD",
   Objectives: map[float64]float64{
    0.5:   0.05,  // Median con error de 5%
    0.9:   0.01,  // p90 con error de 1%
    0.99:  0.001, // p99 con error de 0.1%
    0.999: 0.0001,// p999 con error de 0.01%
   },
  },
  []string{"query_type"},
 )
)

func queryDatabase(query string) {
 start := time.Now()
 // ... ejecutar query ...
 duration := time.Since(start).Seconds()
 dbQueryDuration.WithLabelValues("select").Observe(duration)
}
```

**Histogram vs Summary:**

| Aspecto | Histogram | Summary |
|---------|-----------|---------|
| Cálculo de percentiles | Servidor (PromQL) | Cliente |
| Overhead | Bajo | Medio |
| Precisión | Aproximada | Exacta |
| Agregación | Excelente | Limitada |
| Uso recomendado | Request times | Métricas críticas |

### 49.3.6 Comparación: Prometheus en Diferentes Lenguajes

**Go (prometheus/client_golang):**

```go
counter := promauto.NewCounter(
 prometheus.CounterOpts{Name: "requests_total"},
)
counter.Inc()
```

**Python (prometheus_client):**

```python
from prometheus_client import Counter
counter = Counter('requests_total', 'Total requests')
counter.inc()
```

**Java (Micrometer):**

```java
MeterRegistry registry = new SimpleMeterRegistry();
Counter counter = Counter.builder("requests.total")
    .register(registry);
counter.increment();
```

---

## 49.4 Instrumentación Avanzada: Middleware y Autoinstrumentación

La instrumentación manual es tediosa y propensa a errores. El middleware automático captura métricas sin tocar la lógica de negocio.

### 49.4.1 Middleware HTTP

```go
package main

import (
 "net/http"
 "strconv"
 "time"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
 requestsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_requests_total",
   Help: "Total requests",
  },
  []string{"method", "endpoint", "status"},
 )

 requestDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_duration_seconds",
   Help:    "Request duration",
   Buckets: []float64{.001, .01, .1, 1},
  },
  []string{"method", "endpoint"},
 )

 requestSize = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_size_bytes",
   Help:    "Request size",
   Buckets: prometheus.ExponentialBuckets(100, 10, 5),
  },
  []string{"method"},
 )

 responseSize = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_response_size_bytes",
   Help:    "Response size",
   Buckets: prometheus.ExponentialBuckets(100, 10, 5),
  },
  []string{"method", "status"},
 )
)

// ResponseWriter wrapper para capturar status y tamaño
type MetricsWriter struct {
 http.ResponseWriter
 statusCode int
 size       int
}

func (m *MetricsWriter) WriteHeader(statusCode int) {
 m.statusCode = statusCode
 m.ResponseWriter.WriteHeader(statusCode)
}

func (m *MetricsWriter) Write(b []byte) (int, error) {
 m.size += len(b)
 return m.ResponseWriter.Write(b)
}

// Middleware que instrumenta requests
func metricsMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  start := time.Now()

  // Capturar tamaño del request
  requestSize.WithLabelValues(r.Method).Observe(
   float64(r.ContentLength),
  )

  // Wrapper para capturar response
  mw := &MetricsWriter{ResponseWriter: w, statusCode: 200}

  // Ejecutar handler
  next.ServeHTTP(mw, r)

  // Registrar métricas
  duration := time.Since(start).Seconds()
  status := strconv.Itoa(mw.statusCode)
  endpoint := r.URL.Path

  requestsTotal.WithLabelValues(r.Method, endpoint, status).Inc()
  requestDuration.WithLabelValues(r.Method, endpoint).Observe(duration)
  responseSize.WithLabelValues(r.Method, status).Observe(float64(mw.size))
 })
}

func main() {
 mux := http.NewServeMux()

 // Handler protegido por middleware
 mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("user list"))
 })

 // Aplicar middleware globalmente
 handler := metricsMiddleware(mux)

 http.ListenAndServe(":8080", handler)
}
```

### 49.4.2 Instrumentación de Base de Datos

```go
type InstrumentedDB struct {
 db *sql.DB
}

var (
 dbQueryDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "db_query_duration_seconds",
   Help:    "DB query duration",
   Buckets: []float64{.001, .01, .1, 1},
  },
  []string{"operation", "table"},
 )

 dbErrors = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "db_errors_total",
   Help: "DB errors",
  },
  []string{"operation", "error_type"},
 )

 dbConnections = promauto.NewGaugeVec(
  prometheus.GaugeOpts{
   Name: "db_connections_active",
   Help: "Active DB connections",
  },
  []string{"type"},
 )
)

func (idb *InstrumentedDB) QueryRow(ctx context.Context,
 query string, args ...interface{}) *sql.Row {

 start := time.Now()
 row := idb.db.QueryRowContext(ctx, query, args...)

 duration := time.Since(start).Seconds()
 dbQueryDuration.WithLabelValues("query", "users").Observe(duration)

 return row
}

func (idb *InstrumentedDB) Exec(ctx context.Context,
 query string, args ...interface{}) (sql.Result, error) {

 start := time.Now()
 result, err := idb.db.ExecContext(ctx, query, args...)

 duration := time.Since(start).Seconds()
 dbQueryDuration.WithLabelValues("exec", "users").Observe(duration)

 if err != nil {
  dbErrors.WithLabelValues("exec", "unknown").Inc()
 }

 return result, err
}

func (idb *InstrumentedDB) UpdateConnectionMetrics() {
 stats := idb.db.Stats()
 dbConnections.WithLabelValues("open").Set(float64(stats.OpenConnections))
 dbConnections.WithLabelValues("in_use").Set(float64(stats.InUse))
 dbConnections.WithLabelValues("idle").Set(float64(stats.Idle))
}
```

### 49.4.3 Instrumentación de Goroutines

```go
var (
 goroutinesActive = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "goroutines_active",
   Help: "Número de goroutines activas",
  },
 )

 taskQueueLength = promauto.NewGaugeVec(
  prometheus.GaugeOpts{
   Name: "task_queue_length",
   Help: "Tamaño de colas de tareas",
  },
  []string{"queue_name"},
 )
)

// Pool de goroutines con instrumentación
type InstrumentedWorkerPool struct {
 taskQueue chan func()
 workers   int
 done      chan struct{}
}

func NewWorkerPool(workers int) *InstrumentedWorkerPool {
 pool := &InstrumentedWorkerPool{
  taskQueue: make(chan func(), 1000),
  workers:   workers,
  done:      make(chan struct{}),
 }

 // Iniciar workers
 for i := 0; i < workers; i++ {
  go pool.worker()
 }

 // Monitorear métricas cada segundo
 go pool.monitorMetrics()

 return pool
}

func (p *InstrumentedWorkerPool) Submit(task func()) {
 p.taskQueue <- task
}

func (p *InstrumentedWorkerPool) worker() {
 for task := range p.taskQueue {
  task()
  goroutinesActive.Dec()
 }
}

func (p *InstrumentedWorkerPool) monitorMetrics() {
 ticker := time.NewTicker(time.Second)
 defer ticker.Stop()

 for range ticker.C {
  goroutinesActive.Set(float64(runtime.NumGoroutine()))
  taskQueueLength.WithLabelValues("default").Set(
   float64(len(p.taskQueue)),
  )
 }
}
```

### 49.4.4 Auto-instrumentación de HTTP Requests Salientes

```go
type InstrumentedHTTPClient struct {
 client *http.Client
}

var (
 outboundRequests = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_client_requests_total",
   Help: "Outbound HTTP requests",
  },
  []string{"method", "status", "host"},
 )

 outboundDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_client_duration_seconds",
   Help:    "Outbound HTTP duration",
   Buckets: []float64{.001, .01, .1, 1, 5},
  },
  []string{"method", "host"},
 )
)

func (ic *InstrumentedHTTPClient) Do(req *http.Request) (*http.Response, error) {
 start := time.Now()

 resp, err := ic.client.Do(req)

 duration := time.Since(start).Seconds()
 method := req.Method
 host := req.Host
 status := "unknown"

 if resp != nil {
  status = strconv.Itoa(resp.StatusCode)
 }

 outboundRequests.WithLabelValues(method, status, host).Inc()
 outboundDuration.WithLabelValues(method, host).Observe(duration)

 return resp, err
}
```

---

## 49.5 Alerting: Detectar Problemas Automáticamente

Las métricas sin alertas son solo gráficas bonitas. El alerting es lo que permite accionar sobre problemas reales.

### 49.5.1 Concepto de Alertas

Una alerta tiene tres partes:

1. **Condición**: regla PromQL que evalúa
2. **Umbral**: cuándo consideramos que es crítico
3. **Acción**: dónde notificar (email, Slack, etc)

### 49.5.2 Alert Rules en Prometheus

Archivo `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['localhost:9093']

rule_files:
  - "alert_rules.yml"
  - "recording_rules.yml"

scrape_configs:
  - job_name: 'app'
    static_configs:
      - targets: ['localhost:8080']
```

Archivo `alert_rules.yml`:

```yaml
groups:
  - name: application_alerts
    interval: 15s
    rules:
      # Alerta: tasa de error > 5%
      - alert: HighErrorRate
        expr: |
          (
            sum(rate(http_requests_total{status=~"5.."}[5m]))
            /
            sum(rate(http_requests_total[5m]))
          ) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Tasa de error alta"
          description: "Tasa de error es {{ $value | humanizePercentage }}"

      # Alerta: latencia p95 > 500ms
      - alert: HighLatency
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket[5m]))
            by (le)
          ) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Latencia alta detectada"
          description: "p95 latency: {{ $value }}s"

      # Alerta: instancia caída
      - alert: InstanceDown
        expr: up{job="app"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Instancia {{ $labels.instance }} está caída"
          description: "No se detecta heartbeat por 1 minuto"

      # Alerta: memoria > 80%
      - alert: HighMemoryUsage
        expr: |
          (memory_usage_bytes / memory_limit_bytes) > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Uso alto de memoria"
          description: "Memoria: {{ $value | humanizePercentage }}"

      # Alerta: queue backlog
      - alert: QueueBacklog
        expr: task_queue_length > 1000
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Queue backlog alto"
          description: "{{ $value }} tareas en cola"
```

### 49.5.3 Ciclo de Vida de Alertas

```
INACTIVE (condición false)
    ↓ (expr becomes true)
PENDING (waiting 'for' duration)
    ↓ (duration met)
FIRING (alert active)
    ├─ Send to Alertmanager
    ├─ Repeat according to 'repeat_interval'
    └─ Until condition becomes false
         ↓
    INACTIVE
```

### 49.5.4 Severidad y Contexto

Las alertas deben tener contexto útil:

```yaml
- alert: DatabaseConnectionPoolExhausted
  expr: db_connections_active / db_connections_max > 0.9
  for: 2m
  labels:
    severity: critical      # critical, warning, info
    team: backend          # equipo responsable
    runbook: db-pool-exhausted  # referencia a runbook
    slo: database-availability
  annotations:
    summary: "Connection pool exhausted on {{ $labels.instance }}"
    description: |
      Database {{ $labels.db_name }} has {{ $value | humanizePercentage }}
      connections in use. Check: https://runbooks.example.com/db-pool-exhausted
    dashboard: "http://grafana.example.com/d/abc123"
    action: "Restart app instances or scale database"
```

### 49.5.5 Testing de Alertas

```go
// alert_test.go
package main

import (
 "testing"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHighErrorRateAlert(t *testing.T) {
 // Reset counters
 requestsTotal.Reset()

 // Simular 95% de errores (alerta debe activarse)
 for i := 0; i < 100; i++ {
  if i < 95 {
   requestsTotal.WithLabelValues("GET", "/api", "500").Inc()
  } else {
   requestsTotal.WithLabelValues("GET", "/api", "200").Inc()
  }
 }

 // Verificar que métrica tiene el valor esperado
 // Nota: se verifica contra archivo Prometheus scrape
}
```

---

## 49.6 Grafana: Visualización y Dashboards

Prometheus almacena datos, pero Grafana es donde ocurre la visualización y comprensión.

### 49.6.1 Conceptos de Grafana

**Datasources**: conexiones a Prometheus/Loki/etc
**Dashboards**: colecciones de paneles
**Panels**: gráficos individuales
**Queries**: búsquedas PromQL
**Alerts**: notificaciones basadas en condiciones

### 49.6.2 Configuración de Datasource

JSON de Grafana (`provisioning/datasources/prometheus.yml`):

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    url: http://prometheus:9090
    access: proxy
    isDefault: true
    jsonData:
      timeInterval: 15s
```

### 49.6.3 Dashboard JSON

```json
{
  "dashboard": {
    "title": "API Monitoring",
    "uid": "api-monitoring",
    "timezone": "UTC",
    "panels": [
      {
        "id": 1,
        "title": "Requests per second",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total[5m])) by (status)",
            "legendFormat": "{{ status }}"
          }
        ],
        "yaxes": [
          {"label": "req/s"}
        ]
      },
      {
        "id": 2,
        "title": "Error Rate",
        "type": "stat",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))",
            "legendFormat": "error rate"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "unit": "percentunit",
            "thresholds": {
              "mode": "absolute",
              "steps": [
                {"color": "green", "value": null},
                {"color": "yellow", "value": 0.01},
                {"color": "red", "value": 0.05}
              ]
            }
          }
        }
      },
      {
        "id": 3,
        "title": "Latency Percentiles",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "p50"
          },
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "p95"
          },
          {
            "expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "p99"
          }
        ]
      }
    ]
  }
}
```

### 49.6.4 Dashboard Avanzados con Variables

```json
{
  "dashboard": {
    "templating": {
      "list": [
        {
          "name": "job",
          "type": "query",
          "datasource": "Prometheus",
          "query": "label_values(up, job)",
          "current": {"text": "app", "value": "app"},
          "multi": false
        },
        {
          "name": "instance",
          "type": "query",
          "datasource": "Prometheus",
          "query": "label_values(up{job=\"$job\"}, instance)",
          "multi": true
        },
        {
          "name": "interval",
          "type": "interval",
          "value": "5m",
          "options": ["1m", "5m", "10m", "30m", "1h"]
        }
      ]
    },
    "panels": [
      {
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{job=\"$job\", instance=~\"$instance\"}[$interval]))",
            "legendFormat": "{{instance}}"
          }
        ]
      }
    ]
  }
}
```

### 49.6.5 Alertas en Grafana

```json
{
  "panel": {
    "alert": {
      "name": "High Error Rate",
      "conditions": [
        {
          "evaluator": {"params": [0.05], "type": "gt"},
          "operator": {"type": "and"},
          "query": {
            "params": ["A", "5m", "now"]
          },
          "reducer": {"type": "avg"},
          "type": "query"
        }
      ],
      "executionErrorState": "alerting",
      "frequency": "1m",
      "handler": 1,
      "message": "Error rate exceeded threshold",
      "noDataState": "no_data",
      "notifications": [
        {"uid": "slack-channel"}
      ]
    }
  }
}
```

---

## 49.7 PromQL: Language de Queries

PromQL es el lenguaje para consultar datos en Prometheus. Dominarla es essential.

### 49.7.1 Tipos de Datos

**Instant Vector**: conjunto de series en un momento

```
http_requests_total{job="app"}
→ 1234 @ timestamp1
→ 5678 @ timestamp2
```

**Range Vector**: series en un rango de tiempo

```
http_requests_total[5m]
→ [1000, 1050, 1100, 1150, 1200]
```

**Scalar**: número simple

```
5 or 2.5
```

**String**: texto (raro)

```
"hello"
```

### 49.7.2 Selectores (Matching)

```promql
# Exact match
http_requests_total{job="api"}

# Regex match
http_requests_total{endpoint=~"/api/.*"}

# Negación
http_requests_total{status!="200"}

# Regex negation
http_requests_total{status!~"2.."}

# Múltiples condiciones
http_requests_total{job="api", status="500"}

# Combinación
http_requests_total{
  method=~"GET|POST",
  status!~"2..",
  endpoint=~"/api/users.*"
}
```

### 49.7.3 Operadores Aritméticos

```promql
# Suma de dos métricas
requests_total + errors_total

# Multiplicación (escala)
memory_bytes * 1024  # convertir a bits

# División
(errors_total / requests_total) * 100  # porcentaje

# Divisiones especiales
rate(requests_total[5m]) / 60  # por segundo

# Módulo
total_items % 10
```

### 49.7.4 Funciones de Rango

Estas funciones operan sobre range vectors:

```promql
# Tasa (requests por segundo en 5 minutos)
rate(http_requests_total[5m])

# Aumento (cambio total en rango)
increase(http_requests_total[5m])

# Derivada (aproximación de pendiente)
deriv(temperature_celsius[1h])

# Delta (cambio entre primero y último punto)
delta(cpu_temp_celsius[30m])

# Interpolación lineal
absent_over_time(up[5m])
```

### 49.7.5 Funciones de Agregación

```promql
# Suma total
sum(http_requests_total)

# Máximo
max(memory_usage_bytes)

# Promedio
avg(http_request_duration_seconds)

# Mínimo
min(active_connections)

# Cuantil (percentil)
quantile(0.95, http_request_duration_seconds)

# Contar series
count(http_requests_total)

# Desviación estándar
stddev(cpu_load)

# Con agrupación
sum by (status) (http_requests_total)
# → total requests by status code
# → {status="200"}: 5000
# → {status="500"}: 100

# Sin dimensiones específicas
sum without (endpoint) (http_requests_total)
# → agrupa por todo EXCEPTO endpoint
```

### 49.7.6 Subqueries y Ventanas

```promql
# Subconsulta: máximo de rates en última hora
max_over_time(rate(requests_total[5m])[1h:5m])
# Evalúa rate(requests_total[5m]) cada 5 min durante 1h, luego max()

# Ventana deslizante
max_over_time(memory_usage_bytes[10m])  # máximo en última ventana
min_over_time(memory_usage_bytes[10m])  # mínimo
avg_over_time(memory_usage_bytes[10m])  # promedio
```

### 49.7.7 Histograms y Percentiles

```promql
# Percentil 95 de duración de requests
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
)

# p99 por endpoint
histogram_quantile(0.99,
  sum(rate(http_request_duration_seconds_bucket[5m]))
  by (endpoint, le)
)

# Cantidad de requests en bucket < 100ms
sum(increase(http_request_duration_seconds_bucket{le="0.1"}[5m]))
```

### 49.7.8 Ejemplos Prácticos de PromQL

```promql
# 1. Tasa de error en porcentaje
(
  sum(rate(http_requests_total{status=~"5.."}[5m]))
  /
  sum(rate(http_requests_total[5m]))
) * 100

# 2. Requests pendientes por servicio
sum by (service) (task_queue_length)

# 3. Cambio de memoria en últimas 2 horas
delta(memory_usage_bytes[2h])

# 4. Top 5 endpoints por tráfico
topk(5, sum by (endpoint) (rate(http_requests_total[5m])))

# 5. Bottom 3 endpoints (menos tráfico)
bottomk(3, sum by (endpoint) (rate(http_requests_total[5m])))

# 6. Instancias con latencia > 200ms
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m]))
  by (instance, le)
) > 0.2

# 7. Crecimiento de conexiones activas
rate(active_connections[1h])

# 8. Volatilidad de latencia
stddev(rate(http_request_duration_seconds[5m]))

# 9. Cambio porcentual vs ayer
(
  rate(http_requests_total[5m])
  /
  rate(http_requests_total offset 24h [5m])
) - 1

# 10. Predicción simple (2 horas)
predict_linear(requests_total[1h], 7200)
```

---

## 49.8 Recording Rules: Pre-computar Queries

Queries complejas pueden ser costosas. Recording rules pre-computan resultados.

### 49.8.1 Concepto

Sin recording rules:

```
Usuario consulta:
  histogram_quantile(0.95,
    sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
  )
  ↓
Prometheus computa cada vez
  ↓
Resultado
```

Con recording rules:

```
Prometheus automáticamente cada minuto:
  job:request_latency:p95 =
    histogram_quantile(0.95,
      sum(rate(http_request_duration_seconds_bucket[1m])) by (le)
    )
  ↓
Usuario consulta:
  job:request_latency:p95
  ↓
Resultado (pre-computado)
```

### 49.8.2 Recording Rules YAML

```yaml
# recording_rules.yml
groups:
  - name: performance
    interval: 1m
    rules:
      # Latencias: pre-computar percentiles
      - record: job:request_latency:p50
        expr: histogram_quantile(0.50,
          sum(rate(http_request_duration_seconds_bucket[1m])) by (le, job)
        )

      - record: job:request_latency:p95
        expr: histogram_quantile(0.95,
          sum(rate(http_request_duration_seconds_bucket[1m])) by (le, job)
        )

      - record: job:request_latency:p99
        expr: histogram_quantile(0.99,
          sum(rate(http_request_duration_seconds_bucket[1m])) by (le, job)
        )

      # Tasa de error
      - record: job:error_rate:1m
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[1m])) by (job)
          /
          sum(rate(http_requests_total[1m])) by (job)

      # Throughput
      - record: job:requests_per_second:1m
        expr: sum(rate(http_requests_total[1m])) by (job)

      # Requests exitosos
      - record: job:success_requests:1m
        expr: sum(rate(http_requests_total{status=~"2.."}[1m])) by (job)

      # Conexiones activas
      - record: instance:connections:active
        expr: active_connections

      # Memoria como porcentaje del límite
      - record: instance:memory_usage:percent
        expr: (memory_usage_bytes / memory_limit_bytes) * 100

  - name: availability
    interval: 5m
    rules:
      # Disponibilidad en 5 minutos
      - record: job:availability:5m
        expr: (sum(up{job!=""}) by (job) / count(up{job!=""}) by (job)) * 100

      # Alertas basadas en availability
      - alert: ServiceUnavailable
        expr: job:availability:5m < 50
        for: 5m
        annotations:
          summary: "{{ $labels.job }} availability below 50%"
```

### 49.8.3 Convención de Naming

La convención es: `level:metric:aggregation_period`

```
job:request_latency:p95
  ├─ level: "job" (agregación en esta dimensión)
  ├─ metric: "request_latency" (qué se mide)
  └─ aggregation: "p95" (qué se calcula)

instance:memory_usage:percent
  ├─ level: "instance"
  ├─ metric: "memory_usage"
  └─ aggregation: "percent"

cluster:disk_free:bytes
  ├─ level: "cluster"
  ├─ metric: "disk_free"
  └─ aggregation: "bytes"
```

### 49.8.4 Performance Impact

Recording rules reducen:

- **CPU en Prometheus**: pre-computa vs compute on demand
- **Latencia de queries**: resultado ya está listo
- **Complejidad**: queries simples en dashboards

Incrementan:

- **Storage**: pre-computados son nuevas series
- **Network I/O**: Prometheus deve evaluar cada intervalo

**Regla**: usar recording rules para queries que se ejecutan frecuentemente (dashboards) pero no para alertas únicas.

---

## 49.9 Best Practices en Observabilidad

### 49.9.1 Cardinalidad: El Mayor Problema

**Alta cardinalidad** = muchos valores únicos de labels = disaster.

```go
// ❌ MALO: user_id como label
requestDuration.WithLabelValues(r.Method, user_id, status).Observe(duration)
// Si hay 1M usuarios → 1M series × 3 métodos × 50 status = EXPLOSIÓN

// ✅ BIEN: solo dimensiones con cardinalidad baja
requestDuration.WithLabelValues(r.Method, endpoint, status).Observe(duration)
// 10 métodos × 100 endpoints × 10 status = 10,000 series

// ✅ Alternativa: tracking en logs
log.WithFields(log.Fields{
 "user_id": user_id,
 "request_id": request_id,
}).Info("Request completed")
```

**Guía de Cardinalidad Segura:**

| Tipo | Cardinalidad Máxima | Ejemplo |
|------|-------------------|---------|
| HTTP Method | 10 | GET, POST, PUT, DELETE, ... |
| Status Code | 50 | 200, 201, 400, 404, 500, ... |
| Endpoint | 100-500 | /api/users, /api/posts, ... |
| Database | 10-20 | primary, replica1, replica2, ... |
| Environment | 5-10 | prod, staging, dev, test |
| Service | 20-50 | auth, api, db, cache, ... |
| ❌ User ID | 1M+ | 12345, 67890, 99999, ... |
| ❌ Request ID | ∞ | abc123def456, ... |
| ❌ Session ID | ∞ | sess_xyz, ... |

### 49.9.2 Label Design

```go
// ❌ MALO
requestDuration.WithLabelValues(
 fmt.Sprintf("%s_%s_%s", method, endpoint, user_id)).Observe(duration)
// Label como string complejo, imposible de queryar

// ✅ BIEN
requestDuration.WithLabelValues(method, endpoint).Observe(duration)
// Cada dimensión es su propio label

// ✅ Usar labels solo para dimensiones de consulta
userIDLabel := metrics.NewLabel("user_id")  // NO: cardinalidad infinita
methodLabel := metrics.NewLabel("method")    // OK: cardinalidad ~10
regionLabel := metrics.NewLabel("region")    // OK: cardinalidad ~10
```

### 49.9.3 SLOs y SLIs: Definir lo que Importa

**SLI (Service Level Indicator)**: métrica medible de experiencia del usuario
**SLO (Service Level Objective)**: meta para ese SLI

```go
// SLI 1: Disponibilidad
// Métrica: uptime
// SLO: 99.9% (9 horas downtime/mes máximo)

// SLI 2: Latencia
// Métrica: p95 request time < 500ms
// SLO: 99% de requests cumplen esto

// SLI 3: Error rate
// Métrica: errores 5xx
// SLO: < 0.1% (menos de 1 por 1000 requests)

// Implementación
var (
 sliAvailability = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "sli:availability:percent",
   Help: "SLI: Availability percentage",
  },
 )

 sliLatency = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "sli:latency:p95_seconds",
   Help: "SLI: p95 request latency",
  },
 )

 sliErrorRate = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "sli:error_rate:percent",
   Help: "SLI: Error rate percentage",
  },
 )
)

// Alert si SLO no se cumple
var alertSLOBreach = `
  - alert: SLOBreach
    expr: |
      (sli:latency:p95_seconds > 0.5) or
      (sli:error_rate:percent > 0.1) or
      (sli:availability:percent < 99.9)
    for: 5m
    annotations:
      summary: "SLO breach detected"
`
```

### 49.9.4 Instrumentación Mínima Necesaria

No necesitas medir todo. Comienza con:

**Mínimo obligatorio (Tier 1):**

```go
// 1. Tasa de requests
requestsTotal.WithLabelValues(method, status).Inc()

// 2. Latencia
requestDuration.WithLabelValues(method).Observe(duration)

// 3. Disponibilidad
up.Set(1)  // o 0 si está caída
```

**Útil (Tier 2):**

```go
// 4. Tamaño de request/response
requestSizeBytes.Observe(float64(req.ContentLength))

// 5. Error details
errorsTotal.WithLabelValues(errorType, severity).Inc()

// 6. Recursos
memoryUsageBytes.Set(float64(getMemory()))
```

**Avanzado (Tier 3):**

```go
// 7. Trazas distribuidas
tracer.Start(ctx, "operation")

// 8. Logs estructurados
log.WithContext(ctx).WithFields(fields).Info("event")

// 9. Profiling
go func() {
 ticker := time.NewTicker(1 * time.Minute)
 for range ticker.C {
  pprof.WriteHeapProfile(file)
 }
}()
```

### 49.9.5 Evitar Antipatterns

```go
// ❌ ANTIPATTERN 1: Métrica de identidad
// Si cada request es único, métrica no sirve
requestID := uuid.New().String()
metricWithRequestID.WithLabelValues(requestID).Inc()
// → millones de series, base de datos explota

// ✅ CORRECCIÓN: Usar logs
log.WithField("request_id", requestID).Info("Processing")

// ❌ ANTIPATTERN 2: Métrica muy agregada
// Pierdes información útil
totalMetric.Inc()  // ¿cuándo ocurrió? ¿por qué?
// → no puedes responder "cuál endpoint falla?"

// ✅ CORRECCIÓN: Labels apropiados
requestDuration.WithLabelValues(endpoint, method).Observe(duration)

// ❌ ANTIPATTERN 3: Métrica sin contexto
errorCount.Inc()  // ¿qué tipo de error?
// → no sabes si es config error o bug runtime

// ✅ CORRECCIÓN: Semántica clara
errorsTotal.WithLabelValues("validation", "user_input").Inc()
errorsTotal.WithLabelValues("runtime", "null_pointer").Inc()

// ❌ ANTIPATTERN 4: Métrica nunca consultada
someRandomMetric.Inc()  // te pareció que podría ser útil
// → ruido en base de datos, nadie la consulta

// ✅ CORRECCIÓN: Solo medir lo que respondas
// "¿Qué es SLI del servicio?" → captura eso
// "¿Cuál es la latencia p95?" → captura eso
```

---

## 49.10 Escalabilidad: De 1 Servidor a Millones de Series

### 49.10.1 Limitaciones de Prometheus Single-Node

Un Prometheus standalone típicamente maneja:

- 1-2 millones de series activas
- 500k-1M scrape targets
- Queries en ~100ms

Problemas a escala:

- **Storage**: serie por (métrica, labels) = explosión
- **Ingesta**: scraping de 10k targets simultáneamente
- **Query**: evaluar millones de series cada segundo
- **HA**: no hay replicación automática

### 49.10.2 Remote Storage

Delega storage histórico a sistemas externos:

```yaml
# prometheus.yml
remote_write:
  - url: http://thanos-receive:19291/api/v1/receive
    queue_config:
      max_shards: 10
      min_shards: 1
      capacity: 10000
      min_backoff: 30ms
      max_backoff: 100ms

remote_read:
  - url: http://thanos-query:10901/api/v1/read
    read_recent: true
```

Sistemas populares:

- **Thanos**: Open source, almacenamiento en S3/GCS/Azure
- **Cortex**: Multi-tenant, ideal para SaaS
- **VictoriaMetrics**: Performance ultra-optimizado

### 49.10.3 Prometheus Federation

Múltiples Prometheus en jerarquía:

```
┌─────────────────────────┐
│   Global Prometheus     │
│   (queries all)         │
└────────┬────────────────┘
         │ scrape
    ┌────┴────────────────────┐
    ↓                         ↓
┌──────────────┐      ┌──────────────┐
│ Regional     │      │ Regional     │
│ Prometheus 1 │      │ Prometheus 2 │
└──────────────┘      └──────────────┘
    ↓ scrape              ↓ scrape
┌──────────────┐      ┌──────────────┐
│ Local apps   │      │ Local apps   │
└──────────────┘      └──────────────┘
```

Configuración de federation:

```yaml
# global-prometheus.yml
scrape_configs:
  - job_name: 'regional-1'
    static_configs:
      - targets: ['regional1-prometheus:9090']
    metrics_path: '/federate'
    params:
      'match[]':
        - 'up'
        - '{job!=""}'

# regional-prometheus.yml
scrape_configs:
  - job_name: 'local-app'
    static_configs:
      - targets: ['app1:8080']
```

### 49.10.4 Sharding por Label

Distribuye carga de ingesta por múltiples Prometheus:

```go
// Shard ID basado en label
func getShardID(labelValue string) int {
 hash := fnv.New32a()
 hash.Write([]byte(labelValue))
 return int(hash.Sum32()) % numShards
}

// Enviar a Prometheus correcto
shard := getShardID(user_id)
prometheusInstance := fmt.Sprintf("prometheus-%d:9090", shard)
pushMetric(prometheusInstance, metric)
```

### 49.10.5 Thanos: La Solución Moderna

Thanos convierte múltiples Prometheus en sistema distribuido:

```yaml
# docker-compose.yml con Thanos
services:
  prometheus-1:
    image: prom/prometheus
    volumes:
      - ./prometheus-1.yml:/etc/prometheus/prometheus.yml
    command:
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=2h'

  thanos-sidecar-1:
    image: quay.io/thanos/thanos
    command:
      - sidecar
      - --tsdb.path=/prometheus
      - --prometheus.url=http://prometheus-1:9090
      - --objstore.config-file=/etc/thanos/objstore.yml
      - --grpc-address=0.0.0.0:10901

  thanos-querier:
    image: quay.io/thanos/thanos
    command:
      - query
      - --store=dns+thanos-sidecar-1:10901
      - --store=dns+thanos-sidecar-2:10901
    ports:
      - "9090:10902"
```

Beneficios de Thanos:

- HA automático
- Retención ilimitada (S3/GCS)
- Query distribuida
- Downsampling automático

```

---

## 49.11 Buenas Prácticas y Antipatterns Comunes

### 49.11.1 Checklist de Implementación

```

□ Métrica de disponibilidad (up) implementada
□ Latencias (request duration) medidas
□ Errores contados por tipo
□ Logs correlacionados con request IDs
□ Trazas distribuidas en servicios críticos
□ Dashboards de SLIs principales
□ Alertas para SLOs
□ Cardinality monitoreada
□ Retención de datos configurada (mínimo 15 días)
□ Alertas testeadas (¿disparan cuando deben?)
□ Runbooks vinculados en alertas
□ Equipo capacitado en PromQL básico

```

### 49.11.2 Patrones de Éxito

**Patrón 1: El Triángulo de Observabilidad**
```go
// Implementar los 3 pilares juntos
// No solo métricas, también logs y trazas

// Métrica detecta problema
if errorRate > 0.05 {
    // Log proporciona contexto
    log.Error("High error rate",
        "error_rate", errorRate,
        "affected_endpoint", endpoint)

    // Trace muestra dónde falló
    span.AddEvent("error_detected",
        trace.WithAttributes(
            attribute.String("endpoint", endpoint),
            attribute.Float64("error_rate", errorRate),
        ))
}
```

**Patrón 2: Histogramas + Recording Rules + Alertas**

```yaml
# Histograma: capta distribución
http_request_duration_seconds_bucket

# Recording rule: pre-computa percentiles
job:request_latency:p95

# Alerta: usa recording rule (eficiente)
- alert: HighLatency
  expr: job:request_latency:p95 > 0.5
```

**Patrón 3: Labels por Contexto de Consulta**

```go
// Labels que responden: "¿cuál es?"
requestDuration.WithLabelValues(
    method,      // ¿cuál HTTP method?
    endpoint,    // ¿cuál endpoint?
    status,      // ¿cuál status?
).Observe(duration)

// NO labels de identidad única
// (user_id, session_id, request_id)
```

### 49.11.3 Debugging de Problemas Comunes

**Problema: "Prometheus está lento"**

```promql
# 1. Verificar cardinality
count(ALERTS)  # alertas activas

# 2. Verificar series potenciales
sum(count({__name__=~".+"})) by (__name__)
# Buscar spikes en cardinality

# 3. Queries lentas
topk(10, prometheus_tsdb_symbol_table_size_bytes)

# Solución: remover labels de alta cardinalidad
```

**Problema: "Alertas no disparan"**

```yaml
# Verificar:
# 1. Expresión PromQL correcta
# → Ir a Prometheus -> Alerts, ver si expresa correctamente

# 2. for: duration suficiente
# → Si es muy corta, fluctuaciones disparan falsas alarmas

# 3. Alertmanager conectado
# → curl http://prometheus:9090/api/v1/alerts | grep firing

# 4. Receiver configurado
# → alertmanager.yml tiene notificaciones?
```

**Problema: "Storage lleno"**

```promql
# Investigar qué ocupa espacio
topk(20, sum by (__name__)
  (increases(storage_bucket_chunks_total[1h])))

# Reducir:
# 1. Retention time más baja
# 2. Remote storage + downsampling
# 3. Remover labels innecesarios
```

### 49.11.4 Evolución Recomendada

**Semana 1-2: Foundation**

- Prometheus + Grafana deployed
- Métricas básicas (requests, latency, errors)
- Dashboard simple con 5-10 gráficos

**Semana 3-4: Alerting**

- Alertas para SLOs críticos
- Alertmanager + notificaciones
- Runbooks para alertas

**Mes 2: Escalabilidad**

- Evaluación de cardinality
- Optimization de labels
- Remote storage si crece

**Mes 3+: Madurez**

- Trazas distribuidas
- Logs centralizados
- Profiling continuo
- Correlación automática (traces ↔ logs ↔ metrics)

---

## Ejercicios Progresivos

### Ejercicio 1: Prometheus Básico

**Objetivo**: Crear app que expone métricas simples y scrapearlas con Prometheus

**Código** (`main.go`):

```go
package main

import (
 "fmt"
 "net/http"
 "time"

 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
 requestsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "api_requests_total",
   Help: "Total API requests",
  },
  []string{"method", "endpoint", "status"},
 )

 requestDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "api_request_duration_seconds",
   Help:    "Request latency in seconds",
   Buckets: []float64{.001, .01, .1, 1},
  },
  []string{"method", "endpoint"},
 )
)

func handleAPI(w http.ResponseWriter, r *http.Request) {
 start := time.Now()

 // Simular trabajo
 time.Sleep(time.Duration(100+int(time.Now().Unix()%50)) * time.Millisecond)

 duration := time.Since(start).Seconds()
 status := 200

 requestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprint(status)).Inc()
 requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

 w.Write([]byte("OK"))
}

func main() {
 http.HandleFunc("/api/data", handleAPI)
 http.Handle("/metrics", promhttp.Handler())

 fmt.Println("Server running on :8080")
 http.ListenAndServe(":8080", nil)
}
```

**Configuración Prometheus** (`prometheus.yml`):

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'api'
    static_configs:
      - targets: ['localhost:8080']
```

**Verificación**:

```bash
# App en terminal 1
go run main.go

# Prometheus en terminal 2
./prometheus --config.file=prometheus.yml

# Verificar métricas (terminal 3)
curl http://localhost:8080/metrics | grep api

# Ver en Prometheus UI
# http://localhost:9090/graph
# Query: rate(api_requests_total[1m])
```

---

### Ejercicio 2: Instrumentación con Middleware

**Objetivo**: Aplicar middleware automático de métricas sin tocar handlers

**Código** (`main.go`):

```go
package main

import (
 "fmt"
 "net/http"
 "strconv"
 "time"

 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
 httpRequests = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_requests_total",
  },
  []string{"method", "path", "status"},
 )

 httpDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_duration_seconds",
   Buckets: []float64{.001, .01, .1, 1},
  },
  []string{"method", "path"},
 )

 activeRequests = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "http_requests_active",
  },
 )
)

type ResponseWriter struct {
 http.ResponseWriter
 statusCode int
 size       int
}

func (rw *ResponseWriter) WriteHeader(code int) {
 rw.statusCode = code
 rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
 rw.size += len(b)
 return rw.ResponseWriter.Write(b)
}

func metricsMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  start := time.Now()
  activeRequests.Inc()
  defer activeRequests.Dec()

  rw := &ResponseWriter{ResponseWriter: w, statusCode: 200}
  next.ServeHTTP(rw, r)

  duration := time.Since(start).Seconds()
  status := strconv.Itoa(rw.statusCode)

  httpRequests.WithLabelValues(r.Method, r.URL.Path, status).Inc()
  httpDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
 })
}

func main() {
 mux := http.NewServeMux()

 // Handlers sin instrumentación
 mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
  time.Sleep(50 * time.Millisecond)
  w.Write([]byte("users list"))
 })

 mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
  time.Sleep(100 * time.Millisecond)
  w.Write([]byte("posts list"))
 })

 mux.Handle("/metrics", promhttp.Handler())

 // Aplicar middleware global
 handler := metricsMiddleware(mux)

 fmt.Println("Server with middleware on :8080")
 http.ListenAndServe(":8080", handler)
}
```

**Testing**:

```bash
# Generar tráfico
for i in {1..100}; do
  curl http://localhost:8080/users
  curl http://localhost:8080/posts
done

# Ver métricas
curl http://localhost:8080/metrics | grep http
```

---

### Ejercicio 3: Dashboards Grafana

**Objetivo**: Crear dashboard con 5+ gráficos

**docker-compose.yml**:

```yaml
version: '3'
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - ./grafana-dashboards:/etc/grafana/provisioning/dashboards
```

**Dashboard JSON** (`grafana-dashboards/dashboard.json`):

```json
{
  "dashboard": {
    "title": "API Performance",
    "panels": [
      {
        "id": 1,
        "title": "Requests/sec",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total[1m]))"
          }
        ]
      },
      {
        "id": 2,
        "title": "Error Rate %",
        "targets": [
          {
            "expr": "(sum(rate(http_requests_total{status=~\"5..\"}[1m])) / sum(rate(http_requests_total[1m]))) * 100"
          }
        ]
      },
      {
        "id": 3,
        "title": "Latency (p95)",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[1m])) by (le))"
          }
        ]
      },
      {
        "id": 4,
        "title": "Active Requests",
        "targets": [
          {
            "expr": "http_requests_active"
          }
        ]
      },
      {
        "id": 5,
        "title": "Requests by Path",
        "targets": [
          {
            "expr": "sum by (path) (rate(http_requests_total[1m]))"
          }
        ]
      }
    ]
  }
}
```

---

### Ejercicio 4: Alertas Basadas en Métricas

**Objetivo**: Definir alertas realistas con Alertmanager

**alert_rules.yml**:

```yaml
groups:
  - name: api_alerts
    interval: 15s
    rules:
      - alert: HighErrorRate
        expr: |
          (sum(rate(http_requests_total{status=~"5.."}[5m]))
           / sum(rate(http_requests_total[5m]))) > 0.05
        for: 5m
        annotations:
          summary: "Error rate > 5%"

      - alert: HighLatency
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
          ) > 1
        for: 5m
        annotations:
          summary: "p95 latency > 1s"

      - alert: TooManyActiveRequests
        expr: http_requests_active > 1000
        for: 2m
        annotations:
          summary: "> 1000 active requests"
```

**alertmanager.yml**:

```yaml
global:
  resolve_timeout: 5m

route:
  receiver: 'console'

receivers:
  - name: 'console'
    webhook_configs:
      - url: 'http://localhost:5000/'
```

---

### Ejercicio 5: Full Stack - App Completa

**Objetivo**: Aplicación completa con Prometheus, Grafana, alertas y trazas

**main.go**:

```go
package main

import (
 "context"
 "database/sql"
 "fmt"
 "net/http"
 "strconv"
 "time"

 _ "github.com/mattn/go-sqlite3"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
 "github.com/prometheus/client_golang/prometheus/promhttp"
 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/exporters/jaeger"
 "go.opentelemetry.io/otel/sdk/trace"
)

var (
 // Métricas
 dbQueryDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "db_query_duration_seconds",
   Buckets: []float64{.001, .01, .1, 1},
  },
  []string{"operation"},
 )

 activeConnections = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "active_connections",
  },
 )

 requestLatency = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_latency_seconds",
   Buckets: []float64{.001, .01, .1},
  },
  []string{"endpoint"},
 )
)

type Database struct {
 db *sql.DB
}

func (d *Database) QueryUser(ctx context.Context, id string) (string, error) {
 start := time.Now()
 defer func() {
  dbQueryDuration.WithLabelValues("select_user").Observe(
   time.Since(start).Seconds(),
  )
 }()

 // Tracing
 ctx, span := otel.Tracer("").Start(ctx, "queryUser")
 defer span.End()

 activeConnections.Inc()
 defer activeConnections.Dec()

 row := d.db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = ?", id)
 var name string
 err := row.Scan(&name)
 return name, err
}

func handleGetUser(db *Database) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  start := time.Now()
  defer func() {
   requestLatency.WithLabelValues("/users").Observe(
    time.Since(start).Seconds(),
   )
  }()

  id := r.URL.Query().Get("id")
  name, err := db.QueryUser(r.Context(), id)

  if err != nil {
   w.WriteHeader(404)
   return
  }

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, `{"id":"%s","name":"%s"}`, id, name)
 }
}

func main() {
 // Tracing setup
 exp, _ := jaeger.New(jaeger.WithCollectorEndpoint(
  jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
 ))
 tp := trace.NewTracerProvider(trace.WithBatcher(exp))
 otel.SetTracerProvider(tp)

 // Database
 sqldb, _ := sql.Open("sqlite3", ":memory:")
 db := &Database{db: sqldb}

 // Initialize
 sqldb.Exec("CREATE TABLE users (id TEXT, name TEXT)")
 sqldb.Exec("INSERT INTO users VALUES ('1', 'Alice'), ('2', 'Bob')")

 // Routes
 http.HandleFunc("/user", handleGetUser(db))
 http.Handle("/metrics", promhttp.Handler())

 fmt.Println("Full stack app on :8080")
 http.ListenAndServe(":8080", nil)
}
```

---

## Conclusión

La observabilidad es un pilar fundamental de operaciones modernas. Go, combinado con Prometheus y Grafana, proporciona una plataforma excepcional para construir sistemas visibles y confiables.

**Recuerda:**

1. **Empieza simple**: Métricas básicas (requests, latency, errors)
2. **Crece gradualmente**: Agrega alertas, dashboards, trazas
3. **Mantén cardinality baja**: El mayor problema en Prometheus
4. **Automatiza todo**: Middleware > instrumentación manual
5. **Correlaciona pilares**: Logs + Métricas + Trazas juntos

Con estos principios, construirás sistemas que no solo funcionan, sino que se entienden a sí mismos.

---

## Referencias

- **Prometheus**: <https://prometheus.io/docs>
- **Client Go**: <https://pkg.go.dev/github.com/prometheus/client_golang>
- **PromQL**: <https://prometheus.io/docs/prometheus/latest/querying/basics>
- **Grafana**: <https://grafana.com/docs>
- **Thanos**: <https://thanos.io/>
- **Observability**: <https://o11y.rocks/>

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/49-monitoring-y-metricas/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/49-monitoring-y-metricas):

```bash
cd examples/49-monitoring-y-metricas
go run .
```
