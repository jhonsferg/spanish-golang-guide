# Capítulo 37: Logging y observabilidad

## Introducción

El logging es la tercera columna del trípode de observabilidad en sistemas modernos. Mientras que métricas nos dicen qué está pasando y trazas nos muestran cómo sucede, los logs nos cuentan la historia con contexto y detalles. En Go, el ecosistema de logging ha evolucionado desde el simple paquete `log` hasta `slog` (Go 1.21+), que trae structured logging de primera clase al lenguaje.

Este capítulo explora cómo instrumentar aplicaciones Go con logging efectivo, context propagation, distributed tracing y observabilidad integral. Aprenderás cuándo loguear, qué loguear, y cómo hacerlo sin contaminar tus aplicaciones con millones de eventos innecesarios.

**Principio fundamental:** Los logs son para humanos; las métricas son para máquinas; las trazas conectan ambas.

---

## 37.1 ¿Qué es Observabilidad? - La Tríada

### 37.1.1 Los Tres Pilares

La observabilidad moderna descansa sobre tres pilares interconectados:

```
┌─────────────────────────────────────────────────────────┐
│              OBSERVABILIDAD MODERNA                      │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  [LOGS]           [MÉTRICAS]         [TRAZAS]          │
│  ├─ Eventos      ├─ Números         ├─ Contexto       │
│  ├─ Contexto     ├─ Tendencias      ├─ Causales       │
│  ├─ Debugging    ├─ Alerting        ├─ Dependencias   │
│  └─ Debugging    └─ Performance     └─ Duración       │
│                                                           │
├─────────────────────────────────────────────────────────┤
│   Integración: Request IDs, Correlation IDs,            │
│   Context Propagation, Structured Logging               │
└─────────────────────────────────────────────────────────┘
```

| Pilar | Propósito | Pregunta |
|---|---|---|
| **LOGS** | Contexto y detalles | ¿Qué pasó exactamente? |
| **MÉTRICAS** | Tendencias y alertas | ¿Qué está pasando ahora? |
| **TRAZAS** | Flujo y causalidad | ¿Cómo se relacionan los eventos? |

### 37.1.2 Evolución del Logging en Go

**Generación 1: El paquete `log` estándar (Go 1.0+)**

```go
// Simple, pero limitado
log.Println("User created:", userID)
log.Fatal("Critical error occurred")
```

**Generación 2: Librerías externas (Logrus, Zap, etc.)**

```go
// Structured logging, pero requería dependencias externas
logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action": "created",
}).Info("User created")
```

**Generación 3: `slog` en Go 1.21+**

```go
// Structured logging en la stdlib
logger.InfoContext(ctx, "User created",
    slog.Int("user_id", userID),
    slog.String("action", "created"),
)
```

### 37.1.3 Antipatrones Comunes

❌ **Logging sin contexto**

```go
// Malo: No sabes a qué request pertenece
logger.Info("Database query executed")
```

✅ **Logging con contexto**

```go
// Bueno: Contexto claro
logger.InfoContext(ctx, "Database query executed",
    slog.String("request_id", requestID),
    slog.Int("rows", count),
)
```

❌ **Loguear datos sensibles**

```go
// Malo: Expone credenciales
logger.Info("Login attempt", slog.String("password", pwd))
```

✅ **Loguear información relevante**

```go
// Bueno: Información segura
logger.InfoContext(ctx, "Login attempt",
    slog.String("user", email),
    slog.Bool("success", err == nil),
)
```

---

## 37.2 El Paquete `log` Estándar - Fundamentos

### 37.2.1 Logger Simple

Go incluye un paquete `log` minimalista pero funcional:

```go
package main

import (
    "fmt"
    "log"
    "os"
)

func main() {
    // Logger estándar
    log.Println("Aplicación iniciada")
    log.Print("Mensaje sin newline:")
    log.Printf("Formato: %v\n", 42)

    // Log Fatal detiene la ejecución
    if err := connectDB(); err != nil {
        log.Fatal("No se pudo conectar a BD:", err)
    }
}

func connectDB() error {
    return fmt.Errorf("conexión rechazada")
}
```

### 37.2.2 Flags del Logger

```go
package main

import (
    "log"
    "os"
)

func main() {
    // Sin flags: solo mensaje
    logger1 := log.New(os.Stdout, "APP: ", 0)
    logger1.Println("Mensaje simple")
    // Output: APP: Mensaje simple

    // Con timestamp
    logger2 := log.New(os.Stdout, "", log.LstdFlags)
    logger2.Println("Con timestamp")
    // Output: 2024/01/15 14:30:45 Con timestamp

    // Con ruta del archivo y línea
    logger3 := log.New(os.Stdout, "[DEBUG] ", log.Lshortfile|log.LstdFlags)
    logger3.Println("Con ubicación")
    // Output: [DEBUG] 2024/01/15 14:30:45 main.go:20 Con ubicación

    // Flags disponibles:
    // log.Ldate      - Fecha
    // log.Ltime      - Hora
    // log.Lmicroseconds - Microsegundos
    // log.Llongfile  - Ruta completa
    // log.Lshortfile - Solo archivo:línea
    // log.LUTC       - Zona UTC
}
```

### 37.2.3 Writer Personalizado

```go
package main

import (
    "io"
    "log"
    "os"
    "strings"
)

// CustomWriter filtra logs por nivel
type CustomWriter struct {
    w io.Writer
}

func (cw *CustomWriter) Write(p []byte) (n int, err error) {
    msg := string(p)
    // Filtrar logs innecesarios
    if strings.Contains(msg, "DEBUG") && shouldSkipDebug() {
        return len(p), nil
    }
    return cw.w.Write(p)
}

func shouldSkipDebug() bool {
    return os.Getenv("LOG_DEBUG") != "1"
}

func main() {
    writer := &CustomWriter{w: os.Stdout}
    logger := log.New(writer, "[APP] ", log.LstdFlags)
    logger.Println("Este mensaje se registra")
    logger.Println("[DEBUG] Este puede ser filtrado")
}
```

### 37.2.2 Limitaciones del Paquete `log`

| Limitación | Impacto |
|---|---|
| Sin structured logging | Difícil de parsear en producción |
| Sin niveles | TODO es INFO |
| Sin contexto | Imposible trazar requests |
| Sin handlers customizables | Código boilerplate |

---

## 37.3 Structured Logging con `slog` (Go 1.21+)

### 37.3.1 Introducción a `slog`

`slog` (Structured Log) es el estándar oficial de Go para logging estructurado:

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // Logger por defecto (JSON a stderr)
    logger := slog.Default()

    // Logging básico
    logger.Info("Aplicación iniciada")
    logger.Debug("Mensaje de depuración")
    logger.Warn("Advertencia")
    logger.Error("Error ocurrido")

    // Con atributos
    logger.InfoContext(context.Background(), "Usuario autenticado",
        slog.String("email", "user@example.com"),
        slog.Int("user_id", 42),
        slog.Bool("premium", true),
    )
}
```

### 37.3.2 Creando Loggers Personalizados

```go
package main

import (
    "io"
    "log/slog"
    "os"
)

func main() {
    // Handler de texto (legible para desarrollo)
    textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slog.LevelDebug,
        AddSource: true,
    })
    textLogger := slog.New(textHandler)

    // Handler JSON (para producción)
    jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    })
    jsonLogger := slog.New(jsonHandler)

    // Cambiar el logger por defecto
    slog.SetDefault(jsonLogger)

    // Usar
    textLogger.Info("Esto va a stdout en texto",
        slog.String("env", "dev"),
    )
}
```

### 37.3.3 Niveles de Log

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // Crear logger solo con nivel INFO
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    })
    logger := slog.New(handler)

    logger.Debug("Este NO se registra (nivel INFO)")
    logger.Info("Este SÍ se registra")
    logger.Warn("Esto también")
    logger.Error("Y esto también")

    // Niveles numéricos (puedes crear custom levels)
    // DEBUG  = -4
    // INFO   = 0
    // WARN   = 4
    // ERROR  = 8
}
```

### 37.3.4 Handler Personalizado

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "os"
)

// CustomHandler escribe logs filtrados
type CustomHandler struct {
    handler slog.Handler
    filter  func(slog.Level) bool
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
    // Filtrar por nivel
    if !h.filter(r.Level) {
        return nil
    }
    // Delegar al handler real
    return h.handler.Handle(ctx, r)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return &CustomHandler{
        handler: h.handler.WithAttrs(attrs),
        filter:  h.filter,
    }
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
    return &CustomHandler{
        handler: h.handler.WithGroup(name),
        filter:  h.filter,
    }
}

func main() {
    baseHandler := slog.NewJSONHandler(os.Stdout, nil)

    // Filtro: solo WARNING y ERROR en producción
    customHandler := &CustomHandler{
        handler: baseHandler,
        filter: func(level slog.Level) bool {
            return level >= slog.LevelWarn
        },
    }

    logger := slog.New(customHandler)
    logger.Info("No se registra")
    logger.Warn("Esto SÍ se registra")
}
```

### 37.3.5 Atributos y Contexto

```go
package main

import (
    "context"
    "log/slog"
    "os"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    // Atributos directos
    logger.Info("Usuario creado",
        slog.Int("user_id", 123),
        slog.String("email", "user@example.com"),
        slog.Float64("balance", 99.99),
    )

    // Grupos de atributos
    logger.Info("Transacción completada",
        slog.Group("user",
            slog.Int("id", 123),
            slog.String("email", "user@example.com"),
        ),
        slog.Group("transaction",
            slog.String("id", "TXN-456"),
            slog.Float64("amount", 50.00),
        ),
    )

    // Con contexto
    ctx := context.Background()
    ctx = context.WithValue(ctx, "request_id", "REQ-789")

    logger.InfoContext(ctx, "Request procesado",
        slog.String("path", "/api/users"),
        slog.String("method", "POST"),
    )
}
```

---

## 37.4 Niveles de Log - Cuándo Usar Cada Uno

### 37.4.1 Guía Completa de Niveles

```
DEBUG (-4)  → Información muy detallada para desarrolladores
INFO (0)    → Eventos importantes del negocio/aplicación
WARN (4)    → Algo inesperado pero recuperable
ERROR (8)   → Error que impide una operación pero no detiene app
FATAL (12)  → Error irrecuperable, aplicación se detiene
```

### 37.4.2 Ejemplos Prácticos

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

func main() {
    ctx := context.Background()

    // DEBUG: Detalles internos
    logger.DebugContext(ctx, "Iniciando conexión a BD",
        slog.String("host", "localhost"),
        slog.Int("port", 5432),
    )

    // INFO: Eventos importantes
    logger.InfoContext(ctx, "Conexión a BD establecida",
        slog.Float64("latency_ms", 45.3),
        slog.String("version", "PostgreSQL 14"),
    )

    // WARN: Situaciones inesperadas
    if isHighLoad() {
        logger.WarnContext(ctx, "Sistema bajo carga alta",
            slog.Int("pending_requests", 500),
            slog.Int("warning_threshold", 100),
        )
    }

    // ERROR: Operaciones fallidas
    if err := processPayment(ctx); err != nil {
        logger.ErrorContext(ctx, "Pago rechazado",
            slog.Any("error", err),
            slog.String("payment_id", "PAY-123"),
        )
    }
}

func isHighLoad() bool {
    return true
}

func processPayment(ctx context.Context) error {
    return errors.New("fondos insuficientes")
}
```

### 37.4.3 Antipatrones de Niveles

❌ **Loguear todo en INFO**

```go
// Malo: Ruido excesivo
logger.Info("Variable x =", x)
logger.Info("Iteración 1000")
logger.Info("Query ejecutado")
```

✅ **Usar niveles apropiados**

```go
// Bueno: Información clara
logger.Debug("Variable x =", x)         // Solo en debug
logger.Debug("Iteración 1000")          // Detalle interno
logger.Info("Query ejecutado en 45ms")  // Evento importante
```

---

## 37.5 Handlers y Formatters - Controlando la Salida

### 37.5.1 Text vs JSON

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // Handler de TEXTO (legible)
    textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,
        Level:     slog.LevelDebug,
    })
    textLogger := slog.New(textHandler)
    textLogger.Info("Usuario conectado", slog.String("email", "user@example.com"))
    /* Output:
    time=2024-01-15T14:30:45.123Z level=INFO source=main.go:14 msg="Usuario conectado" email=user@example.com
    */

    // Handler JSON (máquina-legible)
    jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    })
    jsonLogger := slog.New(jsonHandler)
    jsonLogger.Info("Usuario conectado", slog.String("email", "user@example.com"))
    /* Output:
    {"time":"2024-01-15T14:30:45.123Z","level":"INFO","msg":"Usuario conectado","email":"user@example.com"}
    */
}
```

### 37.5.2 Formatter Personalizado para Producción

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "os"
    "time"
)

// ProductionLogFormat estructura de log para producción
type ProductionLogFormat struct {
    Timestamp   string                 `json:"@timestamp"`
    Level       string                 `json:"level"`
    Message     string                 `json:"message"`
    Source      string                 `json:"source,omitempty"`
    RequestID   string                 `json:"request_id,omitempty"`
    UserID      string                 `json:"user_id,omitempty"`
    Duration    int64                  `json:"duration_ms,omitempty"`
    Error       string                 `json:"error,omitempty"`
    Attrs       map[string]interface{} `json:"attrs,omitempty"`
}

// ProductionHandler escribe logs en formato producción
type ProductionHandler struct {
    jsonHandler slog.Handler
}

func (h *ProductionHandler) Handle(ctx context.Context, r slog.Record) error {
    record := ProductionLogFormat{
        Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
        Level:     r.Level.String(),
        Message:   r.Message,
        Attrs:     make(map[string]interface{}),
    }

    // Extraer atributos conocidos
    r.Attrs(func(a slog.Attr) bool {
        switch a.Key {
        case "request_id":
            record.RequestID = a.Value.String()
        case "user_id":
            record.UserID = a.Value.String()
        case "duration_ms":
            record.Duration = a.Value.Int64()
        case "error":
            record.Error = a.Value.String()
        default:
            record.Attrs[a.Key] = a.Value.Any()
        }
        return true
    })

    data, _ := json.Marshal(record)
    fmt.Println(string(data))
    return nil
}

func (h *ProductionHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return h
}

func (h *ProductionHandler) WithGroup(name string) slog.Handler {
    return h
}

func main() {
    handler := &ProductionHandler{}
    logger := slog.New(handler)

    logger.Info("Usuario registrado",
        slog.String("request_id", "REQ-123"),
        slog.String("user_id", "USER-456"),
        slog.Int64("duration_ms", 125),
    )
}
```

### 37.5.3 Múltiples Handlers (Multiplex)

```go
package main

import (
    "io"
    "log/slog"
    "os"
)

// MultiHandler envía logs a múltiples destinos
type MultiHandler struct {
    handlers []slog.Handler
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
    for _, h := range m.handlers {
        if err := h.Handle(ctx, r); err != nil {
            return err
        }
    }
    return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    handlers := make([]slog.Handler, len(m.handlers))
    for i, h := range m.handlers {
        handlers[i] = h.WithAttrs(attrs)
    }
    return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
    handlers := make([]slog.Handler, len(m.handlers))
    for i, h := range m.handlers {
        handlers[i] = h.WithGroup(name)
    }
    return &MultiHandler{handlers: handlers}
}

func main() {
    // Logs a stdout en texto
    textFile, _ := os.Create("app.log")
    textHandler := slog.NewTextHandler(textFile, nil)

    // Logs a archivo en JSON
    jsonFile, _ := os.Create("app.json")
    jsonHandler := slog.NewJSONHandler(jsonFile, nil)

    // Multiplex
    multi := &MultiHandler{handlers: []slog.Handler{
        textHandler,
        jsonHandler,
    }}

    logger := slog.New(multi)
    logger.Info("Este mensaje va a ambos archivos")

    defer textFile.Close()
    defer jsonFile.Close()
}
```

---

## 37.6 Context Propagation - Rastreando Requests

### 37.6.1 Request IDs y Correlation IDs

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "time"

    "github.com/google/uuid"
)

type contextKey string

const (
    requestIDKey contextKey = "request_id"
    userIDKey    contextKey = "user_id"
    traceIDKey   contextKey = "trace_id"
)

// Middleware que inyecta Request ID
func RequestIDMiddleware(next func(context.Context)) func(context.Context) {
    return func(ctx context.Context) {
        requestID := uuid.New().String()
        ctx = context.WithValue(ctx, requestIDKey, requestID)
        next(ctx)
    }
}

func getRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return "unknown"
}

func logWithContext(ctx context.Context, logger *slog.Logger, msg string, attrs ...slog.Attr) {
    // Añadir automáticamente el request ID
    allAttrs := append([]slog.Attr{
        slog.String("request_id", getRequestID(ctx)),
    }, slog.Attr{})

    // Copiar atributos
    for _, attr := range attrs {
        allAttrs = append(allAttrs, attr)
    }

    logger.InfoContext(ctx, msg, allAttrs...)
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    // Simular request
    ctx := context.Background()
    ctx = context.WithValue(ctx, requestIDKey, "REQ-"+uuid.New().String())

    logWithContext(ctx, logger, "Request iniciado",
        slog.String("path", "/api/users"),
        slog.String("method", "GET"),
    )

    time.Sleep(100 * time.Millisecond)

    logWithContext(ctx, logger, "Query a BD",
        slog.String("table", "users"),
        slog.Int("rows", 42),
    )

    logWithContext(ctx, logger, "Response enviada",
        slog.Int("status", 200),
    )
}
```

### 37.6.2 Logger con Contexto Inyectado

```go
package main

import (
    "context"
    "log/slog"
    "os"
)

// ContextLogger envuelve slog.Logger con inyección de contexto
type ContextLogger struct {
    logger *slog.Logger
    attrs  []slog.Attr
}

func (cl *ContextLogger) InfoContext(ctx context.Context, msg string, attrs ...slog.Attr) {
    allAttrs := append(cl.attrs, attrs...)
    cl.logger.InfoContext(ctx, msg, allAttrs...)
}

func (cl *ContextLogger) ErrorContext(ctx context.Context, msg string, attrs ...slog.Attr) {
    allAttrs := append(cl.attrs, attrs...)
    cl.logger.ErrorContext(ctx, msg, allAttrs...)
}

func (cl *ContextLogger) DebugContext(ctx context.Context, msg string, attrs ...slog.Attr) {
    allAttrs := append(cl.attrs, attrs...)
    cl.logger.DebugContext(ctx, msg, allAttrs...)
}

func (cl *ContextLogger) WithAttrs(attrs ...slog.Attr) *ContextLogger {
    newAttrs := append(cl.attrs, attrs...)
    return &ContextLogger{logger: cl.logger, attrs: newAttrs}
}

func main() {
    baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    // Logger con atributos base
    logger := &ContextLogger{
        logger: baseLogger,
        attrs: []slog.Attr{
            slog.String("service", "user-api"),
            slog.String("version", "1.0.0"),
        },
    }

    ctx := context.Background()

    // Todos los logs incluirán service y version
    logger.InfoContext(ctx, "Usuario autenticado",
        slog.String("user_id", "123"),
    )

    // Extender logger para un request específico
    requestLogger := logger.WithAttrs(
        slog.String("request_id", "REQ-456"),
        slog.String("ip", "192.168.1.1"),
    )

    requestLogger.InfoContext(ctx, "Accediendo a recurso",
        slog.String("resource", "/api/profile"),
    )
}
```

### 37.6.3 Propagación a través de Goroutines

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "sync"
    "time"

    "github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func handleRequest(ctx context.Context, logger *slog.Logger, requestID string) {
    // Crear nuevo contexto con request ID
    ctx = context.WithValue(ctx, requestIDKey, requestID)

    logger.InfoContext(ctx, "Request iniciado")

    var wg sync.WaitGroup

    // Procesar en paralelo, pero manteniendo contexto
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(taskNum int) {
            defer wg.Done()
            processTask(ctx, logger, taskNum)
        }(i)
    }

    wg.Wait()
    logger.InfoContext(ctx, "Request completado")
}

func processTask(ctx context.Context, logger *slog.Logger, taskNum int) {
    requestID := ctx.Value(requestIDKey).(string)

    logger.InfoContext(ctx, "Tarea iniciada",
        slog.Int("task", taskNum),
        slog.String("request_id", requestID),
    )

    time.Sleep(100 * time.Millisecond)

    logger.InfoContext(ctx, "Tarea completada",
        slog.Int("task", taskNum),
    )
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    ctx := context.Background()
    requestID := "REQ-" + uuid.New().String()

    handleRequest(ctx, logger, requestID)
}
```

---

## 37.7 Sampling y Rate Limiting - Control de Volumen

### 37.7.1 Problema: Logging Explosivo

En aplicaciones de alto volumen, loguear todo causa:

- Degradación de performance
- Desbordamiento de almacenamiento
- Logs no indexables

### 37.7.2 Sampling Probabilístico

```go
package main

import (
    "context"
    "log/slog"
    "math/rand"
    "os"
)

// SamplingHandler registra solo una fracción de eventos
type SamplingHandler struct {
    handler   slog.Handler
    sampleRate float64 // 0.0 a 1.0
}

func (sh *SamplingHandler) Handle(ctx context.Context, r slog.Record) error {
    // Para logs WARN y ERROR, siempre registrar
    if r.Level >= slog.LevelWarn {
        return sh.handler.Handle(ctx, r)
    }

    // Para DEBUG/INFO, aplicar sampling
    if rand.Float64() > sh.sampleRate {
        return nil
    }

    return sh.handler.Handle(ctx, r)
}

func (sh *SamplingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return &SamplingHandler{
        handler:    sh.handler.WithAttrs(attrs),
        sampleRate: sh.sampleRate,
    }
}

func (sh *SamplingHandler) WithGroup(name string) slog.Handler {
    return &SamplingHandler{
        handler:    sh.handler.WithGroup(name),
        sampleRate: sh.sampleRate,
    }
}

func main() {
    jsonHandler := slog.NewJSONHandler(os.Stdout, nil)

    // Solo 10% de logs INFO/DEBUG
    samplingHandler := &SamplingHandler{
        handler:    jsonHandler,
        sampleRate: 0.1,
    }

    logger := slog.New(samplingHandler)

    for i := 0; i < 100; i++ {
        logger.Info("Evento de alto volumen", slog.Int("iteration", i))
    }
    // ~90 eventos registrados (no 100)
}
```

### 37.7.3 Token Bucket Rate Limiter

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "sync"
    "time"
)

// RateLimitingHandler usa Token Bucket para limitar logs
type RateLimitingHandler struct {
    handler slog.Handler
    rate    int           // eventos por segundo
    mu      sync.Mutex
    tokens  float64
    lastRefill time.Time
}

func NewRateLimitingHandler(handler slog.Handler, rate int) *RateLimitingHandler {
    return &RateLimitingHandler{
        handler:    handler,
        rate:       rate,
        tokens:     float64(rate),
        lastRefill: time.Now(),
    }
}

func (rl *RateLimitingHandler) Handle(ctx context.Context, r slog.Record) error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    // Refill tokens basado en tiempo transcurrido
    now := time.Now()
    elapsed := now.Sub(rl.lastRefill).Seconds()
    rl.tokens = min(float64(rl.rate), rl.tokens+elapsed*float64(rl.rate))
    rl.lastRefill = now

    // Si no hay tokens, descartar
    if rl.tokens < 1 {
        return nil
    }

    rl.tokens--
    return rl.handler.Handle(ctx, r)
}

func (rl *RateLimitingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return &RateLimitingHandler{
        handler:    rl.handler.WithAttrs(attrs),
        rate:       rl.rate,
        tokens:     rl.tokens,
        lastRefill: rl.lastRefill,
    }
}

func (rl *RateLimitingHandler) WithGroup(name string) slog.Handler {
    return &RateLimitingHandler{
        handler:    rl.handler.WithGroup(name),
        rate:       rl.rate,
        tokens:     rl.tokens,
        lastRefill: rl.lastRefill,
    }
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

func main() {
    jsonHandler := slog.NewJSONHandler(os.Stdout, nil)

    // Máximo 1000 logs por segundo
    limitedHandler := NewRateLimitingHandler(jsonHandler, 1000)
    logger := slog.New(limitedHandler)

    // Intentar loguear 2000 eventos en 1 segundo
    start := time.Now()
    count := 0
    for i := 0; i < 2000 && time.Since(start) < time.Second; i++ {
        logger.Info("Evento", slog.Int("id", i))
        count++
    }
    // ~1000 eventos registrados (no 2000)
}
```

### 37.7.4 Adaptive Sampling

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "sync"
    "sync/atomic"
    "time"
)

// AdaptiveSamplingHandler ajusta sampling basado en volumen
type AdaptiveSamplingHandler struct {
    handler slog.Handler
    counter atomic.Int64
    mu      sync.RWMutex
    sampleRate float64
    maxPerSecond int64
}

func (as *AdaptiveSamplingHandler) Handle(ctx context.Context, r slog.Record) error {
    // Siempre registrar errores
    if r.Level >= slog.LevelError {
        return as.handler.Handle(ctx, r)
    }

    as.counter.Add(1)

    // Ajustar sampling si es necesario
    as.mu.RLock()
    rate := as.sampleRate
    as.mu.RUnlock()

    if as.counter.Load()%int64(1.0/rate) == 0 {
        return as.handler.Handle(ctx, r)
    }

    return nil
}

func (as *AdaptiveSamplingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return as
}

func (as *AdaptiveSamplingHandler) WithGroup(name string) slog.Handler {
    return as
}

func main() {
    jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
    adaptiveHandler := &AdaptiveSamplingHandler{
        handler:      jsonHandler,
        sampleRate:   0.1,
        maxPerSecond: 1000,
    }

    logger := slog.New(adaptiveHandler)

    // Simular carga variable
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for i := 0; i < 5; i++ {
        for j := 0; j < 1000; j++ {
            logger.Info("High volume event", slog.Int("batch", i), slog.Int("item", j))
        }
        <-ticker.C
    }
}
```

---

## 37.8 Distributed Tracing - Conectando la Tríada

### 37.8.1 OpenTelemetry Basics

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/attribute"
)

func main() {
    // Crear exportador de trazas
    exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
    )
    defer tp.Shutdown(context.Background())

    otel.SetTracerProvider(tp)
    tracer := otel.Tracer("my-app")

    // Crear span
    ctx, span := tracer.Start(context.Background(), "ProcessOrder")
    defer span.End()

    // Añadir atributos
    span.SetAttributes(
        attribute.String("user_id", "123"),
        attribute.Int("order_id", 456),
    )

    // Log con contexto
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger.InfoContext(ctx, "Procesando orden",
        slog.String("order_id", "456"),
    )

    // Hacer trabajo
    time.Sleep(100 * time.Millisecond)

    span.AddEvent("order_validated",
        trace.WithAttributes(attribute.Bool("valid", true)),
    )

    logger.InfoContext(ctx, "Orden validada")
}
```

### 37.8.2 Spans y Context Propagation

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
    tracer = otel.Tracer("my-app")
}

func handleRequest(ctx context.Context, logger *slog.Logger, requestID string) {
    ctx, span := tracer.Start(ctx, "handleRequest")
    defer span.End()

    span.SetAttributes(attribute.String("request_id", requestID))
    logger.InfoContext(ctx, "Request iniciado")

    // Llamada a BD
    ctx, dbSpan := tracer.Start(ctx, "queryDatabase")
    result := queryDatabase(ctx)
    dbSpan.End()

    logger.InfoContext(ctx, "BD query completado",
        slog.Int("rows", result),
    )

    // Enviar respuesta
    ctx, sendSpan := tracer.Start(ctx, "sendResponse")
    time.Sleep(10 * time.Millisecond)
    sendSpan.End()

    logger.InfoContext(ctx, "Response enviada")
}

func queryDatabase(ctx context.Context) int {
    _, span := tracer.Start(ctx, "executeQuery")
    defer span.End()

    time.Sleep(50 * time.Millisecond)
    return 42
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    handleRequest(context.Background(), logger, "REQ-123")
}
```

### 37.8.3 Jaeger Integration (Exportador Real)

```go
package main

import (
    "context"
    "log/slog"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger/jaegerexporter"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initJaeger(serviceName string) func(context.Context) error {
    // Crear exportador Jaeger
    exporter, err := jaegerexporter.New(
        jaegerexporter.WithAgentHost("localhost"),
        jaegerexporter.WithAgentPort(6831),
    )
    if err != nil {
        panic(err)
    }

    // Crear recurso
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceName(serviceName),
            semconv.ServiceVersion("1.0.0"),
        ),
    )
    if err != nil {
        panic(err)
    }

    // Crear TracerProvider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )

    otel.SetTracerProvider(tp)

    return tp.Shutdown
}

func main() {
    shutdown := initJaeger("user-api")
    defer shutdown(context.Background())

    tracer := otel.Tracer("user-api")
    ctx, span := tracer.Start(context.Background(), "ProcessUserSignup")
    defer span.End()

    logger := slog.Default()
    logger.InfoContext(ctx, "Usuario registrado")
}
```

---

## 37.9 Métricas - Datos Observables

### 37.9.1 Contadores (Counters)

```go
package main

import (
    "context"

    "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/otel"
)

var (
    requestCounter metric.Int64Counter
    errorCounter   metric.Int64Counter
)

func init() {
    meter := otel.GetMeterProvider().Meter("my-app")

    requestCounter, _ = meter.Int64Counter(
        "requests_total",
        metric.WithDescription("Número total de requests"),
    )

    errorCounter, _ = meter.Int64Counter(
        "errors_total",
        metric.WithDescription("Número total de errores"),
    )
}

func processRequest(ctx context.Context) error {
    // Incrementar contador de requests
    requestCounter.Add(ctx, 1)

    err := someOperation()
    if err != nil {
        errorCounter.Add(ctx, 1)
        return err
    }

    return nil
}

func someOperation() error {
    return nil
}

func main() {
    processRequest(context.Background())
}
```

### 37.9.2 Histogramas (Histograms)

```go
package main

import (
    "context"
    "time"

    "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/otel"
)

var requestDuration metric.Float64Histogram

func init() {
    meter := otel.GetMeterProvider().Meter("my-app")

    requestDuration, _ = meter.Float64Histogram(
        "request_duration_seconds",
        metric.WithDescription("Duración de requests en segundos"),
        metric.WithUnit("s"),
    )
}

func handleRequest(ctx context.Context) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        requestDuration.Record(ctx, duration)
    }()

    // Hacer trabajo
    time.Sleep(100 * time.Millisecond)
}

func main() {
    for i := 0; i < 100; i++ {
        handleRequest(context.Background())
    }
}
```

### 37.9.3 Gauges - Medidores

```go
package main

import (
    "context"
    "runtime"

    "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/otel"
)

var memoryGauge metric.Int64ObservableGauge

func init() {
    meter := otel.GetMeterProvider().Meter("my-app")

    memoryGauge, _ = meter.Int64ObservableGauge(
        "memory_usage_bytes",
        metric.WithDescription("Uso de memoria en bytes"),
    )

    meter.RegisterCallback(
        func(ctx context.Context, obs metric.Observer) error {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            obs.ObserveInt64(memoryGauge, int64(m.Alloc))
            return nil
        },
        memoryGauge,
    )
}

func main() {
    // Metrics se recopilan automáticamente
    ctx := context.Background()
    _ = ctx
}
```

### 37.9.4 Prometheus Integration

```go
package main

import (
    "context"
    "fmt"
    "net/http"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/prometheus"
    "go.opentelemetry.io/otel/sdk/metric"
)

func initPrometheus() {
    // Crear exportador Prometheus
    exporter, err := prometheus.New()
    if err != nil {
        panic(err)
    }

    // Crear MeterProvider
    mp := metric.NewMeterProvider(metric.WithReader(exporter))
    otel.SetMeterProvider(mp)

    // Iniciar servidor HTTP para Prometheus
    http.Handle("/metrics", exporter)
    go func() {
        if err := http.ListenAndServe(":9090", nil); err != nil {
            panic(err)
        }
    }()

    fmt.Println("Prometheus listening on :9090/metrics")
}

func main() {
    initPrometheus()

    meter := otel.GetMeterProvider().Meter("my-app")
    counter, _ := meter.Int64Counter("my_counter")

    for i := 0; i < 100; i++ {
        counter.Add(context.Background(), 1)
    }

    // Métricas disponibles en http://localhost:9090/metrics
    select {}
}
```

---

## 37.10 Error Tracking - Captura de Errores

### 37.10.1 Logging de Errores con Stack Traces

```go
package main

import (
    "fmt"
    "log/slog"
    "os"
    "runtime/debug"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    err := someFailedOperation()
    if err != nil {
        // Opción 1: Con stack completo
        logger.Error("Operación falló",
            slog.Any("error", err),
            slog.String("stack", string(debug.Stack())),
        )

        // Opción 2: Con error wrapping
        logger.Error("Operación falló",
            slog.Any("error", fmt.Errorf("wrapped: %w", err)),
        )
    }
}

func someFailedOperation() error {
    return fmt.Errorf("base de datos no disponible")
}
```

### 37.10.2 Panic Recovery con Logging

```go
package main

import (
    "fmt"
    "log/slog"
    "os"
    "runtime/debug"
)

// SafeExecute ejecuta una función y captura panics
func SafeExecute(logger *slog.Logger, name string, fn func()) {
    defer func() {
        if r := recover(); r != nil {
            logger.Error("Panic recuperado",
                slog.String("function", name),
                slog.Any("panic", r),
                slog.String("stack", string(debug.Stack())),
            )
        }
    }()

    fn()
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    SafeExecute(logger, "riskyOperation", func() {
        panic("Algo terrible pasó!")
    })

    logger.Info("Programa continúa después del panic")
}
```

### 37.10.3 Error Aggregation

```go
package main

import (
    "fmt"
    "log/slog"
    "os"
    "sync"
)

// ErrorCollector recopila errores de múltiples goroutines
type ErrorCollector struct {
    mu     sync.Mutex
    errors []error
}

func (ec *ErrorCollector) Add(err error) {
    if err == nil {
        return
    }
    ec.mu.Lock()
    ec.errors = append(ec.errors, err)
    ec.mu.Unlock()
}

func (ec *ErrorCollector) Log(logger *slog.Logger) {
    ec.mu.Lock()
    defer ec.mu.Unlock()

    if len(ec.errors) == 0 {
        logger.Info("Sin errores")
        return
    }

    for i, err := range ec.errors {
        logger.Error("Error recopilado",
            slog.Int("index", i),
            slog.Any("error", err),
        )
    }
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    collector := &ErrorCollector{}

    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            err := someTask(id)
            collector.Add(err)
        }(i)
    }

    wg.Wait()
    collector.Log(logger)
}

func someTask(id int) error {
    if id == 2 {
        return fmt.Errorf("tarea %d falló", id)
    }
    return nil
}
```

### 37.10.4 Sentry Integration (Error Reporting Real)

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"

    "github.com/getsentry/sentry-go"
)

func initSentry() {
    err := sentry.Init(sentry.ClientOptions{
        Dsn: "https://example@sentry.io/project",
        Environment: "production",
        TracesSampleRate: 1.0,
    })
    if err != nil {
        panic(err)
    }
    defer sentry.Flush(2 * time.Second)
}

func handleError(ctx context.Context, logger *slog.Logger, err error) {
    logger.ErrorContext(ctx, "Error ocurrido", slog.Any("error", err))

    // Enviar a Sentry
    eventID := sentry.CaptureException(err)
    if eventID != nil {
        logger.InfoContext(ctx, "Error reportado a Sentry",
            slog.String("event_id", string(*eventID)),
        )
    }
}

func main() {
    initSentry()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    ctx := context.Background()

    handleError(ctx, logger, fmt.Errorf("algo falló"))
}
```

---

## 37.11 Buenas Prácticas y Patterns - Mastering Observabilidad

### 37.11.1 Log Aggregation Pipeline

```
Aplicación
    ↓
┌─────────────┐
│  Structured │
│   Logging   │
└──────┬──────┘
       ↓
┌─────────────────┐
│ JSON/Text File  │
│ o STDOUT        │
└──────┬──────────┘
       ↓
┌──────────────────────────┐
│ Log Shipper              │
│ (Filebeat, Fluentd)      │
└──────┬───────────────────┘
       ↓
┌──────────────────────────┐
│ Log Aggregator           │
│ (Elasticsearch, Loki)    │
└──────┬───────────────────┘
       ↓
┌──────────────────────────┐
│ Visualization            │
│ (Kibana, Grafana)        │
└──────────────────────────┘
```

### 37.11.2 The 12-Factor App Logging

```go
package main

import (
    "fmt"
    "log/slog"
    "os"
)

// Principio 12FA: Logs a stdout, nunca a archivo
func setupLogging() *slog.Logger {
    // SIEMPRE JSON en producción
    if os.Getenv("ENV") == "production" {
        return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        }))
    }

    // Texto legible en desarrollo
    return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
}

func main() {
    logger := setupLogging()

    logger.Info("Aplicación iniciada",
        slog.String("version", os.Getenv("APP_VERSION")),
        slog.String("env", os.Getenv("ENV")),
    )
}
```

### 37.11.3 Structured Logging Best Practices

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "time"
)

// LogEvent estructura eventos de forma consistente
type LogEvent struct {
    Timestamp   time.Time
    RequestID   string
    UserID      string
    Action      string
    Status      string
    Duration    int64
    Error       string
}

func logEvent(logger *slog.Logger, ctx context.Context, evt LogEvent) {
    attrs := []slog.Attr{
        slog.Time("timestamp", evt.Timestamp),
        slog.String("request_id", evt.RequestID),
        slog.String("user_id", evt.UserID),
        slog.String("action", evt.Action),
        slog.String("status", evt.Status),
        slog.Int64("duration_ms", evt.Duration),
    }

    if evt.Error != "" {
        attrs = append(attrs, slog.String("error", evt.Error))
        logger.ErrorContext(ctx, evt.Action, attrs...)
    } else {
        logger.InfoContext(ctx, evt.Action, attrs...)
    }
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    ctx := context.Background()

    logEvent(logger, ctx, LogEvent{
        Timestamp: time.Now(),
        RequestID: "REQ-123",
        UserID:    "USER-456",
        Action:    "login",
        Status:    "success",
        Duration:  145,
    })
}
```

### 37.11.4 Performance Considerations

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "sync"
)

// BufferedLogger asincroniza logging para mejor performance
type BufferedLogger struct {
    logger *slog.Logger
    ch     chan logMessage
    wg     sync.WaitGroup
}

type logMessage struct {
    ctx   context.Context
    level slog.Level
    msg   string
    attrs []slog.Attr
}

func (bl *BufferedLogger) InfoContext(ctx context.Context, msg string, attrs ...slog.Attr) {
    bl.ch <- logMessage{ctx: ctx, level: slog.LevelInfo, msg: msg, attrs: attrs}
}

func (bl *BufferedLogger) start() {
    bl.wg.Add(1)
    go func() {
        defer bl.wg.Done()
        for m := range bl.ch {
            switch m.level {
            case slog.LevelInfo:
                bl.logger.InfoContext(m.ctx, m.msg, m.attrs...)
            case slog.LevelError:
                bl.logger.ErrorContext(m.ctx, m.msg, m.attrs...)
            }
        }
    }()
}

func (bl *BufferedLogger) close() {
    close(bl.ch)
    bl.wg.Wait()
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    bl := &BufferedLogger{
        logger: logger,
        ch:     make(chan logMessage, 1000),
    }
    bl.start()

    for i := 0; i < 10000; i++ {
        bl.InfoContext(context.Background(), "Event", slog.Int("id", i))
    }

    bl.close()
}
```

### 37.11.5 Alerting Rules

```
# PrometheusAlerts.yml
groups:
  - name: application_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(errors_total[5m]) > 0.05
        for: 5m
        annotations:
          summary: "Error rate > 5% in {{ $labels.service }}"

      - alert: SlowRequests
        expr: histogram_quantile(0.95, request_duration_seconds) > 1.0
        for: 5m
        annotations:
          summary: "P95 latency > 1s in {{ $labels.service }}"

      - alert: DiskSpaceLow
        expr: node_filesystem_avail_bytes < 10737418240  # 10GB
        for: 5m
        annotations:
          summary: "Disk space low on {{ $labels.device }}"
```

### 37.11.6 Retention Policy

```go
package main

import (
    "time"
)

// LogRetentionPolicy define cuánto retener
type LogRetentionPolicy struct {
    DebugRetention time.Duration // 7 días
    InfoRetention  time.Duration // 30 días
    ErrorRetention time.Duration // 90 días
}

var DefaultPolicy = LogRetentionPolicy{
    DebugRetention: 7 * 24 * time.Hour,
    InfoRetention:  30 * 24 * time.Hour,
    ErrorRetention: 90 * 24 * time.Hour,
}

// PurgeOldLogs elimina logs antiguos
func PurgeOldLogs(policy LogRetentionPolicy, now time.Time) {
    // Implementar lógica de purga basada en política
    debugCutoff := now.Add(-policy.DebugRetention)
    infoCutoff := now.Add(-policy.InfoRetention)
    errorCutoff := now.Add(-policy.ErrorRetention)

    _ = debugCutoff
    _ = infoCutoff
    _ = errorCutoff
    // Ejecutar queries SQL para eliminar registros antiguos
}
```

---

## 37.12 Ejercicios Prácticos

### Ejercicio 1: Logger Wrapper Simple

**Objetivo:** Crear un wrapper sobre el logger estándar que inyecte automáticamente contexto global.

```go
// Archivo: logger_wrapper.go
package main

import (
    "context"
    "log/slog"
    "os"
    "sync"
)

type AppLogger struct {
    logger *slog.Logger
    mu     sync.RWMutex
    attrs  []slog.Attr
}

// NewAppLogger crea un nuevo logger con atributos base
func NewAppLogger(serviceName, version string) *AppLogger {
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })

    return &AppLogger{
        logger: slog.New(handler),
        attrs: []slog.Attr{
            slog.String("service", serviceName),
            slog.String("version", version),
        },
    }
}

// Info registra un mensaje informativo
func (al *AppLogger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
    al.mu.RLock()
    allAttrs := append(al.attrs, attrs...)
    al.mu.RUnlock()

    al.logger.InfoContext(ctx, msg, allAttrs...)
}

// Error registra un error
func (al *AppLogger) Error(ctx context.Context, msg string, attrs ...slog.Attr) {
    al.mu.RLock()
    allAttrs := append(al.attrs, attrs...)
    al.mu.RUnlock()

    al.logger.ErrorContext(ctx, msg, allAttrs...)
}

// WithAttrs retorna un nuevo logger con atributos adicionales
func (al *AppLogger) WithAttrs(attrs ...slog.Attr) *AppLogger {
    al.mu.RLock()
    newAttrs := append(al.attrs, attrs...)
    al.mu.RUnlock()

    return &AppLogger{
        logger: al.logger,
        attrs:  newAttrs,
    }
}

// main para testing
func ExerciseOne() {
    logger := NewAppLogger("user-service", "1.0.0")
    ctx := context.Background()

    logger.Info(ctx, "Servicio iniciado")

    requestLogger := logger.WithAttrs(
        slog.String("request_id", "REQ-001"),
    )

    requestLogger.Info(ctx, "Procesando request")
}
```

### Ejercicio 2: Structured Logging con Validación

**Objetivo:** Implementar validación de entrada antes de loguear.

```go
// Archivo: structured_logging.go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "time"
)

type LogEntry struct {
    UserID    int
    Action    string
    Resource  string
    Status    string
    Duration  time.Duration
}

// ValidateAndLog valida la entrada antes de loguear
func ValidateAndLog(logger *slog.Logger, ctx context.Context, entry LogEntry) error {
    // Validación
    if entry.UserID <= 0 {
        return fmt.Errorf("invalid user_id: %d", entry.UserID)
    }
    if entry.Action == "" {
        return fmt.Errorf("action cannot be empty")
    }
    if entry.Status == "" {
        entry.Status = "unknown"
    }

    // Log estructurado
    logger.InfoContext(ctx, entry.Action,
        slog.Int("user_id", entry.UserID),
        slog.String("resource", entry.Resource),
        slog.String("status", entry.Status),
        slog.Int64("duration_ms", entry.Duration.Milliseconds()),
    )

    return nil
}

// ExerciseTwo para testing
func ExerciseTwo() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    ctx := context.Background()

    // Caso válido
    ValidateAndLog(logger, ctx, LogEntry{
        UserID:   123,
        Action:   "create_post",
        Resource: "/api/posts",
        Status:   "success",
        Duration: 250 * time.Millisecond,
    })

    // Caso inválido
    err := ValidateAndLog(logger, ctx, LogEntry{
        UserID:  -1,
        Action:  "update_profile",
    })
    if err != nil {
        logger.Error("Log validation failed", slog.Any("error", err))
    }
}
```

### Ejercicio 3: Request ID Tracking

**Objetivo:** Propagar request IDs a través de múltiples layers.

```go
// Archivo: request_tracking.go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "os"

    "github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// RequestIDMiddleware inyecta request ID en contexto
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), requestIDKey, requestID)

        logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
        logger.InfoContext(ctx, "Request received",
            slog.String("method", r.Method),
            slog.String("path", r.URL.Path),
        )

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// LogWithRequestID registra un mensaje con request ID automático
func LogWithRequestID(logger *slog.Logger, ctx context.Context, msg string, attrs ...slog.Attr) {
    requestID := ctx.Value(requestIDKey).(string)

    allAttrs := append([]slog.Attr{
        slog.String("request_id", requestID),
    }, attrs...)

    logger.InfoContext(ctx, msg, allAttrs...)
}

// ExerciseThree simula tracking
func ExerciseThree() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    // Simular request
    ctx := context.Background()
    ctx = context.WithValue(ctx, requestIDKey, fmt.Sprintf("REQ-%s", uuid.New().String()))

    LogWithRequestID(logger, ctx, "Accessing database",
        slog.String("table", "users"),
    )

    LogWithRequestID(logger, ctx, "Query completed",
        slog.Int("rows", 42),
    )
}
```

### Ejercicio 4: Error Logging con Panic Recovery

**Objetivo:** Capturar y loguear panics de forma segura.

```go
// Archivo: error_logging.go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "runtime/debug"
)

type SafeOperation struct {
    logger *slog.Logger
}

// Execute ejecuta una función con protección contra panics
func (so *SafeOperation) Execute(ctx context.Context, name string, fn func() error) {
    defer func() {
        if r := recover(); r != nil {
            so.logger.ErrorContext(ctx, "Panic recovered",
                slog.String("operation", name),
                slog.Any("panic_value", r),
                slog.String("stack", string(debug.Stack())),
            )
        }
    }()

    if err := fn(); err != nil {
        so.logger.ErrorContext(ctx, "Operation failed",
            slog.String("operation", name),
            slog.Any("error", err),
        )
    } else {
        so.logger.InfoContext(ctx, "Operation succeeded",
            slog.String("operation", name),
        )
    }
}

// ExerciseFour para testing
func ExerciseFour() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    ctx := context.Background()

    safeOp := &SafeOperation{logger: logger}

    // Operación normal
    safeOp.Execute(ctx, "fetch_user", func() error {
        return nil
    })

    // Operación con error
    safeOp.Execute(ctx, "update_profile", func() error {
        return fmt.Errorf("database connection failed")
    })

    // Operación con panic
    safeOp.Execute(ctx, "delete_data", func() error {
        panic("unexpected condition!")
    })
}
```

### Ejercicio 5: Observability Stack Completo

**Objetivo:** Integrar logging, métricas y trazas en una aplicación real.

```go
// Archivo: observability_stack.go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "math/rand"
    "os"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

// UserService simula un servicio con observabilidad completa
type UserService struct {
    logger *slog.Logger
    tracer trace.Tracer
}

func NewUserService(logger *slog.Logger, tracer trace.Tracer) *UserService {
    return &UserService{
        logger: logger,
        tracer: tracer,
    }
}

func (us *UserService) CreateUser(ctx context.Context, email string) (int, error) {
    ctx, span := us.tracer.Start(ctx, "CreateUser")
    defer span.End()

    us.logger.InfoContext(ctx, "Creating user", slog.String("email", email))

    span.SetAttributes(attribute.String("email", email))

    // Simular inserción en BD
    time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

    // Simular fallo ocasional
    if rand.Float64() < 0.1 {
        err := fmt.Errorf("database error")
        us.logger.ErrorContext(ctx, "User creation failed", slog.Any("error", err))
        span.RecordError(err)
        return 0, err
    }

    userID := rand.Intn(10000)
    us.logger.InfoContext(ctx, "User created successfully",
        slog.Int("user_id", userID),
    )

    span.AddEvent("user_created", trace.WithAttributes(
        attribute.Int("user_id", userID),
    ))

    return userID, nil
}

func (us *UserService) GetUser(ctx context.Context, userID int) (string, error) {
    ctx, span := us.tracer.Start(ctx, "GetUser")
    defer span.End()

    span.SetAttributes(attribute.Int("user_id", userID))

    us.logger.DebugContext(ctx, "Fetching user", slog.Int("user_id", userID))

    time.Sleep(50 * time.Millisecond)

    return fmt.Sprintf("user%d@example.com", userID), nil
}

// ExerciseFive para testing
func ExerciseFive() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))

    tracer := otel.Tracer("user-service")

    service := NewUserService(logger, tracer)
    ctx := context.Background()

    // Crear usuarios
    for i := 0; i < 5; i++ {
        email := fmt.Sprintf("user%d@example.com", i)
        userID, err := service.CreateUser(ctx, email)
        if err == nil {
            // Obtener usuario
            _, _ = service.GetUser(ctx, userID)
        }
    }

    logger.Info("Ejercicio completado")
}
```

---

## Conclusión

La observabilidad no es un bolsillo de la arquitectura de software: **es fundamental**. Go, con su filosofía de sencillez y su ecosistema robusto, proporciona herramientas de primer nivel para instrumentar aplicaciones.

**Recordar:**

- 📊 **Logs** para contexto y debugging
- 📈 **Métricas** para tendencias y alertas
- 🔗 **Trazas** para causalidad y correlación
- 🔀 **Contexto** para conectarlo todo

Domina estos conceptos y tus aplicaciones serán transparentes, debuggeables y observables. En producción, eso es invaluable.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/37-logging-y-observabilidad/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/37-logging-y-observabilidad):

```bash
cd examples/37-logging-y-observabilidad
go run .
```
