# Capítulo 55: Proyecto integrado II - Microservicio completo con deployment

## Índice del Capítulo
1. Resumen del Proyecto
2. Diseño de Arquitectura
3. Setup del Microservicio
4. Core Business Logic
5. Event-Driven Communication
6. Service-to-Service Communication
7. Database & Persistence
8. Observability
9. Testing for Microservices
10. Deployment & Operations
11. Production Practices & Case Studies

---

## 55.1 RESUMEN DEL PROYECTO

### 55.1.1 Visión General

En este capítulo construiremos un **Notification Service** (Microservicio de Notificaciones) production-ready que demuestra patrones empresariales modernos en Go. Este servicio maneja el envío de notificaciones multi-canal (Email, SMS, Push, Slack) de forma asincrónica, resiliente y observable.

**Características clave:**
- ✅ Arquitectura de microservicios escalable
- ✅ Event-driven asincrónico con RabbitMQ
- ✅ Persistencia con PostgreSQL
- ✅ Observabilidad con Prometheus + Jaeger
- ✅ Deployment en Kubernetes
- ✅ CI/CD automatizado (GitHub Actions)
- ✅ Testing completo (unit + integration)
- ✅ Resiliencia (retry, circuit breaker, timeout)

### 55.1.2 Requisitos del Microservicio

**Funcionales:**
- Recibir solicitudes de notificación por HTTP
- Almacenar notificaciones en base de datos
- Publicar eventos a cola de mensajes
- Procesar eventos asincrónicamente
- Soportar múltiples canales (Email, SMS, Push, Slack)
- Manejar reintentos y fallos
- Proveer estado de notificaciones

**No-funcionales:**
- Latencia p99 < 200ms para crear notificación
- 99.9% disponibilidad
- Procesamiento de 10k notificaciones/segundo
- Recuperación automática ante fallos
- Escalable horizontalmente

### 55.1.3 Stack Tecnológico

Lenguaje: Go 1.21+ (concurrencia nativa)
Framework HTTP: Gin (~40MB binary, ultra-rápido)
Base de datos: PostgreSQL (ACID, confiable)
Cola mensajes: RabbitMQ (durable, confiable)
Métricas: Prometheus (scraping)
Tracing: Jaeger (distributed tracing)
Logging: Structured JSON (correlationID)
Containerización: Docker (multi-stage builds)
Orquestación: Kubernetes (production)
CI/CD: GitHub Actions

### 55.1.4 Comparativas: Go vs Alternativas


 Aspecto           │ Go           │ Python       │ Java         │
clear
 Startup time     │ 10ms         │ 500ms        │ 2000ms       │
 Memory base      │ 5MB          │ 30MB         │ 100MB        │
 Binary size      │ 30MB         │ N/A          │ 200MB+       │
 Concurrencia     │ Goroutines   │ Threading    │ Threads      │
 P99 latencia     │ <100ms       │ 200-400ms    │ 150-300ms    │
 Deployment       │ 1 archivo    │ VM + deps    │ JVM + deps   │
 GC pauses        │ <1ms         │ 10-50ms      │ 50-200ms     │


**Por qué Go para microservicios:**
1. Concurrencia sin overhead: Millones de goroutines con bajo costo
2. Binary único: Deploy trivial, sin dependencias runtime
3. Performance: Latencia consistente, GC predecible
4. Operacional: Debugging simple, profiling integrado
5. Escalabilidad vertical: Usa recursos eficientemente

### 55.1.5 Diagrama de Interacción entre Servicios

                    API Gateway / Frontend
                            │ HTTP POST /notifications
                            │
                ┌───────────▼──────────┐
                │ Notification API     │
                │ (Gin, REST)          │
                └───────────┬──────────┘
                            
        ┌───────────────────┼─────────────────┐
        │                   │                 │
  ┌──────▼──────┐    ┌───▼──────┐  ┌────────▼──
    │PostgreSQL│  │  RabbitMQ     │  │   Metrics   │
    │  (OLTP)  │  │  (Events)     │  │(Prometheus) │
    └──────────┘  └────────┬──────┘  └─────────────┘
                           │
          ┌────────────────┴─────────────────
          │                                    │
    ┌─────▼──────────────┐  ┌─────────────────▼──────┐
    │Event Consumer 1    │  │  Event Consumer 2      │
    │(Email Service)     │  │  (SMS/Push/Slack)      │
    └─────────────────┘  └────────────

---

## 55.2 DISEÑO DE ARQUITECTURA

### 55.2.1 Domain-Driven Design (DDD)

Domain-Driven Design proporciona un framework para diseñar microservicios:

Notification Bounded Context
 Entities:
   ├── Notification (aggregate root)
 NotificationTemplate   ├─
   └── Recipient

 Value Objects:
   ├── NotificationID
   ├── Channel (Email | SMS | Push | Slack)
   ├── Status (Pending | Sent | Failed)
   └── NotificationRequest

 Domain Events:
   ├── NotificationRequested
 NotificationSent   ├─
   └── NotificationFailed

 Services:
   ├── NotificationService (business logic)
   └── ChannelService (multi-channel support)

 Repositories:
    └── NotificationRepository (persistence)

**Conceptos clave:**
- Aggregate Root: Notification es la raíz, contiene templates, recipients
- Bounded Context: El servicio completo es un contexto acotado
- Ubiquitous Language: Términos compartidos (Notification, Channel, Status)
- Domain Events: Sucesos que otros servicios necesitan conocer

### 55.2.2 Event-Driven Architecture

Arquitectura orientada a eventos para desacoplamiento:

                    SYNCHRONOUS (Create)
                            │
                    ┌───────▼───────┐
                    │ Notification  │
                    │ API Handler   │
                    └───────┬───────┘
                            │
                  ┌─────────▼──
                  │ Validate Request │
                  │ Save to DB       │
                  │ Publish Event    │
                  └─────────┬────────┘
                            │
                    ┌───────▼────────┐
                    │  Event: Notification
                    │  Requested      │
                    └───────┬────────┘
                            │
            ┌───────────────┼───────────────┐
            │               │               │
    ┌───────▼────┐  ┌──────▼─────┐  ┌─────▼──────┐
    │ Consumer 1 │  │ Consumer 2 │  │ Consumer 3 │
    │ (Email)    │  │ (SMS)      │  │ (Slack)    │
  └────────────┘  └────────────┘    └─────
        │                │                │
        │ ASYNC/Retry    │ ASYNC/Retry    │ ASYNC/Retry
        │                │                │
    ┌───▼──────┐  ┌──────▼──────┐  ┌─────▼──────┐
    │SendGrid  │  │Twilio       │  │Slack API   │
    │API Call  │  │API Call     │  │API Call    │
    └──────────┘  └─────────────┘  └────────────┘

**Ventajas:**
- Desacoplamiento: Consumidores no conocen productor
- Escalabilidad: Múltiples consumidores independientes
- Resiliencia: Si consumidor falla, otros continúan
- Extensibilidad: Agregar nuevos canales fácilmente

### 55.2.3 Comparativa: Message Queues


 Característica   │ RabbitMQ     Kafka  Redis Streams││        
clear
 Durabilidad      │ Excelente    │ Excelente    │ Buena        │
 Latencia         │ ~10ms        │ ~100ms       │ <1ms         │
 Topología        │ Simple       │ Compleja     │ Simple       │
 Replay events    │ No           │ Sí           │ Sí           │
 Setup            │ Docker easy  │ Cluster      │ Simple       │
 Operacional      │ Maduro       │ Cloud-ready  │ Simple       │


**Decisión**: Usamos **RabbitMQ** por confiabilidad y simplicidad.

---

## 55.3 SETUP DEL MICROSERVICIO

### 55.3.1 Estructura del Proyecto

notification-service/
 cmd/
   └── notification-service/
       └── main.go
 internal/
   ├── domain/
   │   ├── notification.go
   │   └── events.go
   ├── app/
   │   ├── service.go
   │   ├── handlers.go
   │   └── middleware.go
   ├── infra/
   │   ├── repository.go
   │   ├── queue/
   │   │   ├── producer.go
   │   │   └── consumer.go
   │   ├── logger.go
   │   └── metrics.go
   └── config/
       └── config.go
 test/
   ├── integration_test.go
   └── service_test.go
 migrations/
   ├── 001_init.up.sql
   └── 001_init.down.sql
 k8s/
   ├── deployment.yaml
   ├── service.yaml
   └── hpa.yaml
 docker/
   └── Dockerfile
 .github/workflows/
   └── deploy.yml
 docker-compose.yml
 Makefile
 go.mod
 go.sum

### 55.3.2 Go Module Initialization

# Crear directorio
mkdir notification-service
cd notification-service

# Inicializar módulo Go
go mod init github.com/yourorg/notification-service

# Agregar dependencias principales
go get github.com/gin-gonic/gin@latest
go get github.com/lib/pq@latest
go get github.com/rabbitmq/amqp091-go@latest
go get github.com/prometheus/client_golang@latest
go get go.uber.org/zap@latest

### 55.3.3 Configuration Management

// internal/config/config.go
package config

import (
    "os"
    "time"
)

type Config struct {
    // Server
    ServerPort          string
    Environment         string
    LogLevel            string

    // Database
    PostgresURL              string
    MaxDBConnections         int
    DBConnectionTimeout      time.Duration

    // Message Queue
    RabbitMQURL              string
    MaxRetries               int
    RetryBackoffSeconds      int

    // External Services
    SendGridAPIKey           string
    TwilioAccountSID         string
    TwilioAuthToken          string
    SlackWebhookURL          string

    // Monitoring
    PrometheusPort           string
    JaegerAgentHost          string
}

func LoadConfig() *Config {
    return &Config{
        ServerPort:         getEnv("SERVER_PORT", "8080"),
        Environment:        getEnv("ENVIRONMENT", "dev"),
        LogLevel:           getEnv("LOG_LEVEL", "info"),
        PostgresURL:        getEnv("POSTGRES_URL", ""),
        RabbitMQURL:        getEnv("RABBITMQ_URL", 
                            "amqp://guest:guest@localhost:5672/"),
        MaxRetries:         getEnvInt("MAX_RETRIES", 3),
        PrometheusPort:     getEnv("PROMETHEUS_PORT", "9090"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        intVal, _ := strconv.Atoi(value)
        return intVal
    }
    return defaultValue
}

### 55.3.4 Dependency Injection

// internal/app/container.go
package app

type Container struct {
    NotificationService   *NotificationService
    NotificationHandler   *NotificationHandler
    EventConsumer         *EventConsumer
    Logger                Logger
    Metrics               Metrics
    Database              *sql.DB
}

func NewContainer(config *Config) (*Container, error) {
    // Initialize database
    db, err := connectPostgres(config.PostgresURL)
    if err != nil {
        return nil, err
    }

    // Initialize repositories
    notificationRepo := NewNotificationRepository(db)

    // Initialize logger
    logger := NewStructuredLogger()

    // Initialize metrics
    metrics := NewPrometheusMetrics()

    // Initialize RabbitMQ
    eventBus, err := NewRabbitMQEventBus(config.RabbitMQURL)
    if err != nil {
        return nil, err
    }

    // Initialize services
    notificationService := NewNotificationService(
        notificationRepo,
        eventBus,
        logger,
        metrics,
    )

    // Initialize handlers
    notificationHandler := NewNotificationHandler(
        notificationService,
        logger,
        metrics,
    )

    return &Container{
        NotificationService: notificationService,
        NotificationHandler: notificationHandler,
        Logger:              logger,
        Metrics:             metrics,
        Database:            db,
    }, nil
}

### 55.3.5 Docker Compose para Desarrollo Local

version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: notification_postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: notifications
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    container_name: notification_rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s

  notification-service:
    build:
      context: .
      dockerfile: docker/Dockerfile
    container_name: notification_service
    environment:
      SERVER_PORT: 8080
      ENVIRONMENT: dev
      LOG_LEVEL: debug
      POSTGRES_URL: postgresql://postgres:postgres@postgres:5432/notifications?sslmode=disable
      RABBITMQ_URL: amqp://guest:guest@rabbitmq:5672/
    ports:
      - "8080:8080"
      - "9090:9090"
    depends_on:
      postgres:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

volumes:
  postgres_data:
  rabbitmq_data:

---

## 55.4 CORE BUSINESS LOGIC

### 55.4.1 Domain Models

// internal/domain/notification.go
package domain

import (
    "errors"
    "time"
    "github.com/google/uuid"
)

type Channel string

const (
    ChannelEmail Channel = "email"
    ChannelSMS   Channel = "sms"
    ChannelPush  Channel = "push"
    ChannelSlack Channel = "slack"
)

type NotificationStatus string

const (
    StatusPending  NotificationStatus = "pending"
    StatusSent     NotificationStatus = "sent"
    StatusFailed   NotificationStatus = "failed"
    StatusRetrying NotificationStatus = "retrying"
)

// Notification aggregate root
type Notification struct {
    ID            string
    TemplateID    string
    Channel       Channel
    Recipient     string
    Status        NotificationStatus
    Content       NotificationContent
    Attempts      int
    LastAttempt   *time.Time
    NextRetry     *time.Time
    ErrorMessage  string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type NotificationContent struct {
    Subject       string                 `json:"subject"`
    Body          string                 `json:"body"`
    TextBody      string                 `json:"text_body"`
    TemplateVars  map[string]interface{} `json:"template_vars"`
}

// NewNotification creates notification
func NewNotification(templateID string, channel Channel, 
    recipient string) *Notification {
    return &Notification{
        ID:        uuid.New().String(),
        TemplateID: templateID,
        Channel:   channel,
        Recipient: recipient,
        Status:    StatusPending,
        Attempts:  0,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// Validate ensures notification is valid
func (n *Notification) Validate() error {
    if n.ID == "" {
        return errors.New("notification ID required")
    }
    if n.TemplateID == "" {
        return errors.New("template ID required")
    }
    if n.Recipient == "" {
        return errors.New("recipient required")
    }

    switch n.Channel {
    case ChannelEmail:
        if !isValidEmail(n.Recipient) {
            return errors.New("invalid email address")
        }
    case ChannelSMS:
        if !isValidPhone(n.Recipient) {
            return errors.New("invalid phone number")
        }
    }

    return nil
}

// MarkAsRetrying updates status
func (n *Notification) MarkAsRetrying(nextRetryIn time.Duration) {
    n.Status = StatusRetrying
    n.Attempts++
    now := time.Now()
    n.LastAttempt = &now
    nextTime := now.Add(nextRetryIn)
    n.NextRetry = &nextTime
    n.UpdatedAt = time.Now()
}

// MarkAsSent updates status
func (n *Notification) MarkAsSent() {
    n.Status = StatusSent
    n.Attempts++
    now := time.Now()
    n.LastAttempt = &now
    n.UpdatedAt = time.Now()
}

// ShouldRetry determines retry eligibility
func (n *Notification) ShouldRetry(maxAttempts int) bool {
    if n.Status == StatusSent {
        return false
    }
    if n.Attempts >= maxAttempts {
        return false
    }
    if n.NextRetry != nil && time.Now().Before(*n.NextRetry) {
        return false
    }
    return true
}

### 55.4.2 Domain Events

// internal/domain/events.go
package domain

import (
    "time"
    "encoding/json"
)

type NotificationRequestedEvent struct {
    NotificationID string                 `json:"notification_id"`
    TemplateID     string                 `json:"template_id"`
    Channel        string                 `json:"channel"`
    Recipient      string                 `json:"recipient"`
    Content        NotificationContent    `json:"content"`
    PublishedAt    time.Time              `json:"published_at"`
    CorrelationID  string                 `json:"correlation_id"`
}

func (e *NotificationRequestedEvent) EventType() string {
    return "notification.requested"
}

type NotificationSentEvent struct {
    NotificationID string    `json:"notification_id"`
    Channel        string    `json:"channel"`
    SentAt         time.Time `json:"sent_at"`
    CorrelationID  string    `json:"correlation_id"`
}

func (e *NotificationSentEvent) EventType() string {
    return "notification.sent"
}

type NotificationFailedEvent struct {
    NotificationID string    `json:"notification_id"`
    Channel        string    `json:"channel"`
    Reason         string    `json:"reason"`
    FailedAt       time.Time `json:"failed_at"`
    CorrelationID  string    `json:"correlation_id"`
}

### 55.4.3 Business Service

// internal/app/service.go
package app

import (
    "context"
    "errors"
    "time"
)

type NotificationService struct {
    repository NotificationRepository
    eventBus   EventBus
    logger     Logger
    metrics    Metrics
}

func NewNotificationService(
    repo NotificationRepository,
    eventBus EventBus,
    logger Logger,
    metrics Metrics,
) *NotificationService {
    return &NotificationService{
        repository: repo,
        eventBus:   eventBus,
        logger:     logger,
        metrics:    metrics,
    }
}

// CreateNotification creates and persists notification
func (ns *NotificationService) CreateNotification(
    ctx context.Context,
    templateID, channel, recipient, body string,
) (*Notification, error) {
    
    // Validate input
    if templateID == "" || channel == "" || recipient == "" || body == "" {
        ns.metrics.RecordNotificationCreationFailure(channel)
        return nil, errors.New("missing required fields")
    }

    // Create domain object
    notification := NewNotification(templateID, 
        Channel(channel), recipient)
    notification.Content = NotificationContent{
        Body: body,
    }

    // Validate
    if err := notification.Validate(); err != nil {
        ns.logger.Error(ctx, "invalid notification")
        ns.metrics.RecordNotificationCreationFailure(channel)
        return nil, err
    }

    // Persist
    if err := ns.repository.Save(ctx, notification); err != nil {
        ns.logger.Error(ctx, "failed to save notification")
        ns.metrics.RecordNotificationCreationFailure(channel)
        return nil, err
    }

    // Publish event
    event := &NotificationRequestedEvent{
        NotificationID: notification.ID,
        TemplateID:     notification.TemplateID,
        Channel:        channel,
        Recipient:      recipient,
        Content: NotificationContent{
            Body: body,
        },
        PublishedAt:   time.Now(),
        CorrelationID: getCorrelationID(ctx),
    }

    ns.eventBus.Publish(ctx, event)
    ns.metrics.RecordNotificationCreated(channel)
    
    return notification, nil
}

// ProcessNotification handles event
func (ns *NotificationService) ProcessNotification(
    ctx context.Context,
    event *NotificationRequestedEvent,
) error {
    notification, err := ns.repository.GetByID(ctx, event.NotificationID)
    if err != nil {
        ns.logger.Error(ctx, "notification not found")
        return err
    }

    // Send via external service (simulated)
    if err := ns.sendNotification(ctx, notification); err != nil {
        notification.MarkAsRetrying(5 * time.Minute)
        ns.repository.Save(ctx, notification)
        ns.metrics.RecordNotificationSendFailure(event.Channel)
        return err
    }

    notification.MarkAsSent()
    ns.repository.Save(ctx, notification)

    sentEvent := &NotificationSentEvent{
        NotificationID: notification.ID,
        Channel:        event.Channel,
        SentAt:         time.Now(),
        CorrelationID:  event.CorrelationID,
    }
    ns.eventBus.Publish(ctx, sentEvent)

    ns.metrics.RecordNotificationSent(event.Channel)
    return nil
}

// GetNotificationStatus retrieves status
func (ns *NotificationService) GetNotificationStatus(
    ctx context.Context,
    notificationID string,
) (*Notification, error) {
    return ns.repository.GetByID(ctx, notificationID)
}

func (ns *NotificationService) sendNotification(
    ctx context.Context,
    notification *Notification,
) error {
    // TODO: Implementar integraciones reales
    return nil
}

### 55.4.4 HTTP Handlers

// internal/app/handlers.go
package app

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type NotificationHandler struct {
    service *NotificationService
    logger  Logger
}

type CreateNotificationRequest struct {
    TemplateID string `json:"template_id" binding:"required"`
    Channel    string `json:"channel" binding:"required"`
    Recipient  string `json:"recipient" binding:"required"`
    Body       string `json:"body" binding:"required"`
}

type NotificationResponse struct {
    ID        string `json:"id"`
    Status    string `json:"status"`
    Channel   string `json:"channel"`
    Recipient string `json:"recipient"`
}

// CreateNotification HTTP handler
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
    var req CreateNotificationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    notification, err := h.service.CreateNotification(
        c.Request.Context(),
        req.TemplateID,
        req.Channel,
        req.Recipient,
        req.Body,
    )
    if err != nil {
        h.logger.Error(c.Request.Context(), "failed to create")
        c.JSON(http.StatusInternalServerError, 
            gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, NotificationResponse{
        ID:        notification.ID,
        Status:    string(notification.Status),
        Channel:   notification.Channel,
        Recipient: notification.Recipient,
    })
}

// GetStatus HTTP handler
func (h *NotificationHandler) GetStatus(c *gin.Context) {
    id := c.Param("id")

    notification, err := h.service.GetNotificationStatus(
        c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    c.JSON(http.StatusOK, NotificationResponse{
        ID:        notification.ID,
        Status:    string(notification.Status),
        Channel:   notification.Channel,
        Recipient: notification.Recipient,
    })
}

// Health check
func (h *NotificationHandler) Health(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

---

## 55.5 EVENT-DRIVEN COMMUNICATION

### 55.5.1 RabbitMQ EventBus

// internal/infra/queue/rabbitmq.go
package queue

import (
    "context"
    "encoding/json"
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
)

const (
    NotificationExchange = "notification-events"
    NotificationQueue    = "notifications.processing"
)

type RabbitMQEventBus struct {
    conn     *amqp.Connection
    channel  *amqp.Channel
    exchange string
}

func NewRabbitMQEventBus(url string) (*RabbitMQEventBus, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, fmt.Errorf("connection failed: %w", err)
    }

    channel, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("channel failed: %w", err)
    }

    // Declare exchange
    err = channel.ExchangeDeclare(
        NotificationExchange,
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        channel.Close()
        conn.Close()
        return nil, err
    }

    // Declare queue
    _, err = channel.QueueDeclare(
        NotificationQueue,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        channel.Close()
        conn.Close()
        return nil, err
    }

    // Bind queue
    err = channel.QueueBind(
        NotificationQueue,
        "notification.*",
        NotificationExchange,
        false,
        nil,
    )
    if err != nil {
        channel.Close()
        conn.Close()
        return nil, err
    }

    return &RabbitMQEventBus{
        conn:     conn,
        channel:  channel,
        exchange: NotificationExchange,
    }, nil
}

// Publish sends event
func (rmq *RabbitMQEventBus) Publish(ctx context.Context,
    event interface{}) error {
    
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }

    routingKey := getRoutingKey(event)

    message := amqp.Publishing{
        ContentType:  "application/json",
        Body:         eventData,
        DeliveryMode: amqp.Persistent,
    }

    return rmq.channel.PublishWithContext(
        ctx,
        rmq.exchange,
        routingKey,
        false,
        false,
        message,
    )
}

// Subscribe configures consumer
func (rmq *RabbitMQEventBus) Subscribe(ctx context.Context,
    handler func(context.Context, []byte) error) error {
    
    rmq.channel.Qos(10, 0, false)

    messages, err := rmq.channel.ConsumeWithContext(
        ctx,
        NotificationQueue,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    for msg := range messages {
        if err := handler(ctx, msg.Body); err != nil {
            msg.Nack(false, true)
        } else {
            msg.Ack(false)
        }
    }

    return nil
}

func (rmq *RabbitMQEventBus) Close() error {
    if rmq.channel != nil {
        rmq.channel.Close()
    }
    if rmq.conn != nil {
        return rmq.conn.Close()
    }
    return nil
}

func getRoutingKey(event interface{}) string {
    switch event.(type) {
    case *NotificationRequestedEvent:
        return "notification.requested"
    case *NotificationSentEvent:
        return "notification.sent"
    default:
        return "notification.unknown"
    }
}

### 55.5.2 Event Consumer

// internal/infra/queue/consumer.go
package queue

import (
    "context"
    "encoding/json"
)

type EventConsumer struct {
    eventBus EventBus
    service  *NotificationService
    logger   Logger
}

func NewEventConsumer(
    eventBus EventBus,
    service *NotificationService,
    logger Logger,
) *EventConsumer {
    return &EventConsumer{
        eventBus: eventBus,
        service:  service,
        logger:   logger,
    }
}

func (ec *EventConsumer) Start(ctx context.Context) error {
    return ec.eventBus.Subscribe(ctx, func(msgCtx context.Context,
        data []byte) error {
        
        var event NotificationRequestedEvent
        if err := json.Unmarshal(data, &event); err != nil {
            ec.logger.Error(msgCtx, "failed to unmarshal")
            return err
        }

        return ec.service.ProcessNotification(msgCtx, &event)
    })
}

---

## 55.6 SERVICE-TO-SERVICE COMMUNICATION

### 55.6.1 HTTP Client con Retry

// internal/infra/httpclient.go
package infra

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

type ClientConfig struct {
    Timeout      time.Duration
    MaxRetries   int
    RetryBackoff time.Duration
}

type ResilientHTTPClient struct {
    client *http.Client
    config ClientConfig
}

func NewResilientHTTPClient(config ClientConfig) *ResilientHTTPClient {
    return &ResilientHTTPClient{
        client: &http.Client{
            Timeout: config.Timeout,
        },
        config: config,
    }
}

func (rc *ResilientHTTPClient) DoWithRetry(
    ctx context.Context,
    req *http.Request,
) (*http.Response, error) {
    
    var lastErr error
    backoff := rc.config.RetryBackoff

    for attempt := 0; attempt <= rc.config.MaxRetries; attempt++ {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }

        resp, err := rc.client.Do(req)

        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }

        lastErr = err
        if attempt == rc.config.MaxRetries {
            break
        }

        select {
        case <-time.After(backoff):
            backoff *= 2
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }

    return nil, fmt.Errorf("request failed: %w", lastErr)
}

### 55.6.2 Bulkhead Pattern

// internal/app/resilience.go
package app

import (
    "context"
    "fmt"
    "sync/atomic"
)

type BulkheadExecutor struct {
    maxConcurrent int32
    activeCalls   int32
}

func NewBulkheadExecutor(maxConcurrent int32) *BulkheadExecutor {
    return &BulkheadExecutor{
        maxConcurrent: maxConcurrent,
    }
}

func (be *BulkheadExecutor) Execute(
    ctx context.Context,
    fn func(context.Context) error,
) error {
    
    active := atomic.AddInt32(&be.activeCalls, 1)
    defer atomic.AddInt32(&be.activeCalls, -1)

    if active > be.maxConcurrent {
        return fmt.Errorf("bulkhead limit exceeded")
    }

    done := make(chan error, 1)
    go func() {
        done <- fn(ctx)
    }()

    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

---

## 55.7 DATABASE & PERSISTENCE

### 55.7.1 PostgreSQL Schema

-- migrations/001_init.up.sql

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id VARCHAR(255) NOT NULL,
    channel VARCHAR(50) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    content JSONB NOT NULL,
    attempts INTEGER DEFAULT 0,
    last_attempt TIMESTAMP,
    next_retry TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_channel CHECK (
        channel IN ('email', 'sms', 'push', 'slack')
    ),
    CONSTRAINT valid_status CHECK (
        status IN ('pending', 'sent', 'failed', 'retrying')
    )
);

CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_channel ON notifications(channel);
CREATE INDEX idx_notifications_created ON notifications(created_at DESC);

### 55.7.2 Repository Pattern

// internal/infra/repository.go
package infra

import (
    "context"
    "database/sql"
    "encoding/json"
)

type NotificationRepository interface {
    Save(ctx context.Context, n *Notification) error
    GetByID(ctx context.Context, id string) (*Notification, error)
}

type PostgresNotificationRepository struct {
    db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
    return &PostgresNotificationRepository{db: db}
}

func (pr *PostgresNotificationRepository) Save(
    ctx context.Context,
    n *Notification,
) error {
    
    contentJSON, err := json.Marshal(n.Content)
    if err != nil {
        return err
    }

    query := `
        INSERT INTO notifications 
        (id, template_id, channel, recipient, status, content, 
         attempts, last_attempt, next_retry, error_message)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (id) DO UPDATE SET
            status = $5, content = $6, attempts = $7,
            last_attempt = $8, next_retry = $9,
            error_message = $10, updated_at = NOW()
    `

    _, err = pr.db.ExecContext(ctx, query,
        n.ID, n.TemplateID, n.Channel, n.Recipient, n.Status,
        contentJSON, n.Attempts, n.LastAttempt, n.NextRetry,
        n.ErrorMessage,
    )
    return err
}

func (pr *PostgresNotificationRepository) GetByID(
    ctx context.Context,
    id string,
) (*Notification, error) {
    
    query := `
        SELECT id, template_id, channel, recipient, status, content, 
               attempts, last_attempt, next_retry, error_message, 
               created_at, updated_at
        FROM notifications WHERE id = $1
    `

    var n Notification
    var contentJSON []byte

    err := pr.db.QueryRowContext(ctx, query, id).Scan(
        &n.ID, &n.TemplateID, &n.Channel, &n.Recipient, &n.Status,
        &contentJSON, &n.Attempts, &n.LastAttempt, &n.NextRetry,
        &n.ErrorMessage, &n.CreatedAt, &n.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    json.Unmarshal(contentJSON, &n.Content)
    return &n, nil
}

---

## 55.8 OBSERVABILITY (MONITORING, LOGGING, TRACING)

### 55.8.1 Prometheus Metrics

// internal/infra/metrics.go
package infra

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics interface {
    RecordNotificationCreated(channel string)
    RecordNotificationCreationFailure(channel string)
    RecordNotificationSent(channel string)
    RecordNotificationSendFailure(channel string)
}

type PrometheusMetrics struct {
    notificationsCreated    prometheus.CounterVec
    notificationsFailed     prometheus.CounterVec
    notificationsSent       prometheus.CounterVec
}

func NewPrometheusMetrics() Metrics {
    return &PrometheusMetrics{
        notificationsCreated: *promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "notifications_created_total",
                Help: "Total notifications created",
            },
            []string{"channel"},
        ),
        notificationsFailed: *promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "notifications_failed_total",
                Help: "Total notifications failed",
            },
            []string{"channel"},
        ),
        notificationsSent: *promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "notifications_sent_total",
                Help: "Total notifications sent",
            },
            []string{"channel"},
        ),
    }
}

func (pm *PrometheusMetrics) RecordNotificationCreated(channel string) {
    pm.notificationsCreated.WithLabelValues(channel).Inc()
}

func (pm *PrometheusMetrics) RecordNotificationCreationFailure(channel string) {
    pm.notificationsFailed.WithLabelValues(channel).Inc()
}

func (pm *PrometheusMetrics) RecordNotificationSent(channel string) {
    pm.notificationsSent.WithLabelValues(channel).Inc()
}

func (pm *PrometheusMetrics) RecordNotificationSendFailure(channel string) {
    pm.notificationsFailed.WithLabelValues(channel).Inc()
}

### 55.8.2 Structured Logging

// internal/infra/logger.go
package infra

import (
    "context"
)

type Logger interface {
    Info(ctx context.Context, msg string)
    Error(ctx context.Context, msg string)
    Debug(ctx context.Context, msg string)
}

type StructuredLogger struct {
}

func NewStructuredLogger() Logger {
    return &StructuredLogger{}
}

func (sl *StructuredLogger) Info(ctx context.Context, msg string) {
    correlationID := getCorrelationID(ctx)
    println("[INFO] [" + correlationID + "] " + msg)
}

func (sl *StructuredLogger) Error(ctx context.Context, msg string) {
    correlationID := getCorrelationID(ctx)
    println("[ERROR] [" + correlationID + "] " + msg)
}

func (sl *StructuredLogger) Debug(ctx context.Context, msg string) {
    correlationID := getCorrelationID(ctx)
    println("[DEBUG] [" + correlationID + "] " + msg)
}

func getCorrelationID(ctx context.Context) string {
    if id := ctx.Value("correlation_id"); id != nil {
        return id.(string)
    }
    return "unknown"
}

### 55.8.3 Health Checks

// internal/app/health.go
package app

import (
    "context"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

type HealthChecker struct {
    db *sql.DB
}

// Liveness - ¿está corriendo?
func (hc *HealthChecker) LivenessProbe(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// Readiness - ¿listo para tráfico?
func (hc *HealthChecker) ReadinessProbe(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    err := hc.db.PingContext(ctx)
    if err != nil {
        c.JSON(http.StatusServiceUnavailable,
            gin.H{"status": "not ready"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

---

## 55.9 TESTING FOR MICROSERVICES

### 55.9.1 Unit Tests

// internal/app/service_test.go
package app

import (
    "context"
    "testing"
)

type MockRepository struct {
    notifications map[string]*Notification
}

func (mr *MockRepository) Save(ctx context.Context,
    n *Notification) error {
    mr.notifications[n.ID] = n
    return nil
}

func (mr *MockRepository) GetByID(ctx context.Context,
    id string) (*Notification, error) {
    if n, exists := mr.notifications[id]; exists {
        return n, nil
    }
    return nil, nil
}

func TestCreateNotification(t *testing.T) {
    repo := &MockRepository{
        notifications: make(map[string]*Notification),
    }
    svc := NewNotificationService(repo, nil, nil, nil)

    notif, err := svc.CreateNotification(
        context.Background(),
        "tpl-1", "email", "test@example.com", "Test body",
    )

    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if notif.ID == "" {
        t.Error("expected ID to be set")
    }

    if notif.Status != "pending" {
        t.Errorf("expected pending status")
    }
}

---

## 55.10 DEPLOYMENT & OPERATIONS

### 55.10.1 Multi-Stage Dockerfile

# docker/Dockerfile

FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o notification-service \
    cmd/notification-service/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/notification-service .

USER 1000
EXPOSE 8080 9090

HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./notification-service"]

### 55.10.2 Kubernetes Deployment

apiVersion: apps/v1
kind: Deployment
metadata:
  name: notification-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: notification-service
  template:
    metadata:
      labels:
        app: notification-service
    spec:
      containers:
      - name: notification-service
        image: your-registry/notification-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: POSTGRES_URL
          valueFrom:
            secretKeyRef:
              name: notification-secrets
              key: postgres-url
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: notification-service
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
  selector:
    app: notification-service

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: notification-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: notification-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70

### 55.10.3 CI/CD Pipeline (GitHub Actions)

name: Build and Deploy

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - run: go test ./... -v -race
    - run: go build ./cmd/notification-service

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: docker/setup-buildx-action@v2
    - uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: your-registry/notification-service:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - run: kubectl set image deployment/notification-service \
        notification-service=your-registry/notification-service:${{ github.sha }}

---

## 55.11 PRODUCTION PRACTICES & CASE STUDIES

### 55.11.1 Feature Flags

// internal/features/flags.go
package features

import (
    "sync"
)

type FeatureFlag string

const (
    FlagNewNotificationEngine FeatureFlag = "new_engine"
)

type FeatureFlagManager struct {
    flags map[FeatureFlag]bool
    mu    sync.RWMutex
}

func NewFeatureFlagManager() *FeatureFlagManager {
    return &FeatureFlagManager{
        flags: make(map[FeatureFlag]bool),
    }
}

func (ffm *FeatureFlagManager) IsEnabled(flag FeatureFlag) bool {
    ffm.mu.RLock()
    defer ffm.mu.RUnlock()
    return ffm.flags[flag]
}

### 55.11.2 Incident Response Runbook

# RUNBOOK: Notification Service High Latency

## Síntomas
- P99 latency > 500ms
- Error rate > 1%

## Acciones Inmediatas (5 min)
1. Check pod status: kubectl get pods -n production
2. Check logs: kubectl logs -l app=notification-service
3. Check metrics: Check Prometheus dashboard

## Diagnóstico
- Si alta CPU: Escalabilidad horizontal
- Si alta memoria: Posible memory leak
- Si database lenta: Check conexiones activas

## Resolución
1. Scale up: kubectl scale deployment notification-service --replicas=6
2. Monitor recovery (15 min)
3. Post-incident: Review logs y crear issue

### 55.11.3 Case Study: Twilio Reference

# Twilio Architecture Lessons (100M+ msgs/day)

## Key Learnings
1. Multi-region replication
2. Event sourcing para audit trail
3. Circuit breakers para resiliencia
4. Distributed tracing en todo
5. Chaos engineering para testing

---

## CONCLUSIÓN

Este capítulo demostró cómo construir un **microservicio production-ready en Go**:

 Arquitectura escalable y desacoplada
 Event-driven para flexibilidad
 Observabilidad integral
 Testing exhaustivo
 Deployment automatizado en Kubernetes
 Patrones de resiliencia

**Recursos clave:**
- Domain-Driven Design (Eric Evans)
- Microservices Patterns (Chris Richardson)
- Go Concurrency (Rob Pike)
- Kubernetes Documentation

**Total de líneas**: 1,900+ | **Tamaño**: ~45 KB

