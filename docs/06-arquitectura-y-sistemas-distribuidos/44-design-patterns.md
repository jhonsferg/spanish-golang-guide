# Capítulo 44: Design patterns en Go

## Tabla de Contenidos
1. [Descripción General de Patrones de Diseño](#descripción-general)
2. [Patrón Singleton](#patrón-singleton)
3. [Patrón Factory](#patrón-factory)
4. [Patrón Builder](#patrón-builder)
5. [Patrón Decorator](#patrón-decorator)
6. [Patrón Strategy](#patrón-strategy)
7. [Patrón Observer](#patrón-observer)
8. [Patrón Adapter](#patrón-adapter)
9. [Patrón Repository](#patrón-repository)
10. [Inyección de Dependencias](#inyección-de-dependencias)
11. [Buenas Prácticas y Antipatterns](#buenas-prácticas-y-antipatterns)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 44.1 - Descripción General de Patrones de Diseño {#descripción-general}

### ¿Qué son los Patrones de Diseño?

Los patrones de diseño son soluciones probadas a problemas comunes en el desarrollo de software. Fueron popularizados por la obra "Design Patterns: Elements of Reusable Object-Oriented Software" (Gang of Four). En Go, algunos patrones se adaptan de manera idiomatic aprovechando características del lenguaje.

**Definición según GoF:**
> "Un patrón de diseño es una solución general y reutilizable a un problema común dentro de un contexto particular del diseño de software."

### Clasificación de Patrones

Los 23 patrones originales se dividen en tres categorías:

#### **1. Patrones Creacionales** (5 patrones)
Controlan la creación de objetos para evitar la complejidad y mantener la flexibilidad:

- **Singleton**: Una única instancia controlada globalmente
- **Factory Method**: Crear objetos sin especificar sus clases concretas
- **Abstract Factory**: Crear familias de objetos relacionados
- **Builder**: Construir objetos complejos paso a paso
- **Prototype**: Clonar objetos existentes en lugar de crear nuevos

#### **2. Patrones Estructurales** (7 patrones)
Tratan con composición de clases y objetos para formar estructuras más grandes:

- **Adapter**: Convertir interfaz de una clase en otra esperada
- **Bridge**: Desacoplar abstracción de su implementación
- **Composite**: Componer objetos en árboles para representar jerarquías
- **Decorator**: Añadir responsabilidades dinámicamente
- **Facade**: Proporcionar interfaz unificada simplificada
- **Flyweight**: Compartir objetos comunes eficientemente
- **Proxy**: Proporcionar sustituto para otro objeto

#### **3. Patrones Comportamentales** (11 patrones)
Definen patrones de comunicación e interacción entre objetos:

- **Chain of Responsibility**: Pasar solicitud por cadena de manejadores
- **Command**: Encapsular solicitud como objeto
- **Iterator**: Acceder elementos de colección secuencialmente
- **Mediator**: Centralizar comunicación entre objetos
- **Memento**: Capturar/restaurar estado de objeto
- **Observer**: Notificar cambios a múltiples observadores
- **State**: Cambiar comportamiento según estado interno
- **Strategy**: Usar algoritmos intercambiables
- **Template Method**: Esqueleto de algoritmo en clase base
- **Visitor**: Ejecutar operaciones en elementos de estructura
- **Interpreter**: Interpretar lenguaje específico del dominio

### Go: Lenguaje Pragmático con Patrones

Go rechaza ciertos patrones por filosofía de diseño:

| Patrón | Java | Go | Razón |
|--------|------|----|----- |
| Singleton | ✓ Común | ⚠️ Raro | Go prefiere inicialización simple |
| Factory | ✓ Común | ✓ Común | Adaptado con funciones |
| Builder | ⚠️ Workaround | ✓ Natural | Go se adapta mejor con structs |
| Visitor | ✓ Usado | ✗ Evitar | Type assertions son más simples |
| Abstract Factory | ✓ Común | ⚠️ Raro | Innecesario en Go |

### Filosofía Go sobre Patrones

La comunidad Go sigue estos principios:

```
"Patrones: Sí, si son idiomáticos
         No, si son innecesarios
```

**Principios Go:**
1. **Interfaces pequeñas**: Una interfaz = un método
2. **Composición sobre herencia**: Embedding vs herencia
3. **Explícito sobre implícito**: Código claro preferible a magic
4. **Funciones first-class**: Muchas "pseudo-patterns" usan solo funciones

---

## 44.2 - Patrón Singleton {#patrón-singleton}

### Concepto

El patrón Singleton asegura que una clase tenga una única instancia y proporciona un punto global de acceso a ella.

**Problema que resuelve:**
```
❌ Múltiples instancias de recursos costosos (BD, logger)
❌ Inconsistencia de estado global
❌ Acceso no controlado desde cualquier lugar
```

### Implementación Básica en Go

#### **Versión Simple (Unsafe - Solo para Single-threaded)**

```go
package singleton

import "fmt"

type Logger struct {
    messages []string
}

var logger *Logger

// Inicialización global - NO es thread-safe
func GetLogger() *Logger {
    if logger == nil {
        logger = &Logger{
            messages: make([]string, 0),
        }
    }
    return logger
}

func (l *Logger) Log(msg string) {
    l.messages = append(l.messages, msg)
    fmt.Println(msg)
}
```

**Problema:** En aplicaciones concurrentes, dos goroutines pueden crear instancias distintas.

#### **Versión Sync.Once (Thread-Safe - Recomendada)**

```go
package singleton

import (
    "fmt"
    "sync"
)

type Logger struct {
    messages []string
    mu       sync.Mutex
}

var (
    logger *Logger
    once   sync.Once
)

// Thread-safe gracias a sync.Once
func GetLogger() *Logger {
    once.Do(func() {
        logger = &Logger{
            messages: make([]string, 0),
        }
    })
    return logger
}

func (l *Logger) Log(msg string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    l.messages = append(l.messages, msg)
    fmt.Println("[LOG]", msg)
}

func (l *Logger) GetMessages() []string {
    l.mu.Lock()
    defer l.mu.Unlock()
    return append([]string{}, l.messages...)
}
```

#### **Versión con Inicialización Global (Más Idiomático)**

```go
package logger

import (
    "fmt"
    "sync"
)

type Logger struct {
    mu  sync.Mutex
    out []string
}

// Singleton global
var global = &Logger{
    out: make([]string, 0),
}

func Info(msg string) {
    global.mu.Lock()
    defer global.mu.Unlock()
    
    global.out = append(global.out, fmt.Sprintf("INFO: %s", msg))
    fmt.Println(msg)
}

func GetAll() []string {
    global.mu.Lock()
    defer global.mu.Unlock()
    return append([]string{}, global.out...)
}
```

### Comparación con Otros Lenguajes

#### **Python: Decorator o Metaclass**
```python
# Python usa decorador para Singleton
@singleton
class DatabaseConnection:
    def __init__(self):
        self.connected = False
```

#### **Java: Eager o Lazy Initialization**
```java
// Java requiere sincronización manual
public class Singleton {
    private static Singleton instance;
    
    public synchronized static Singleton getInstance() {
        if (instance == null) {
            instance = new Singleton();
        }
        return instance;
    }
}
```

#### **Go: sync.Once (Elegante y Seguro)**
```go
// Go simplifica con sync.Once
var instance *MyType
var once sync.Once

func GetInstance() *MyType {
    once.Do(func() {
        instance = &MyType{}
    })
    return instance
}
```

### Casos de Uso en Go Real

#### **1. Logger Global (logging package)**
```go
package logging

import (
    "log"
    "os"
    "sync"
)

type Logger struct {
    *log.Logger
    level int
}

var (
    logInstance *Logger
    logOnce     sync.Once
)

func Instance() *Logger {
    logOnce.Do(func() {
        logInstance = &Logger{
            Logger: log.New(os.Stdout, "[APP] ", log.LstdFlags),
            level:  1,
        }
    })
    return logInstance
}

func Debug(msg string) {
    Instance().Println("DEBUG:", msg)
}
```

#### **2. Database Connection Pool**
```go
package database

import (
    "database/sql"
    "sync"
)

type DB struct {
    *sql.DB
}

var (
    dbInstance *DB
    dbOnce     sync.Once
)

func GetDB() *DB {
    dbOnce.Do(func() {
        sqlDB, _ := sql.Open("postgres", "...")
        dbInstance = &DB{sqlDB}
    })
    return dbInstance
}
```

#### **3. Configuration Manager**
```go
package config

import (
    "encoding/json"
    "os"
    "sync"
)

type Config struct {
    Database string
    Port     int
    Debug    bool
}

var (
    cfg  *Config
    once sync.Once
)

func Load() *Config {
    once.Do(func() {
        data, _ := os.ReadFile("config.json")
        cfg = &Config{}
        json.Unmarshal(data, cfg)
    })
    return cfg
}
```

### Antipatterns: Cuándo NO Usar Singleton

❌ **Acoplamiento Global**: Dificulta testing y refactoring
❌ **Hidden Dependencies**: Difícil ver qué depende de qué
❌ **Overkill para Simple Vars**: Go permite simplemente `var db *Database`

```go
// ❌ MALO: Demasiada complejidad para simple var
func GetDatabase() *Database {
    // 100 líneas de singleton boilerplate...
}

// ✅ BIEN: Inicialización directa
var db *Database

func init() {
    db = &Database{...}
}
```

---

## 44.3 - Patrón Factory {#patrón-factory}

### Concepto

El patrón Factory proporciona una interfaz para crear objetos sin especificar sus clases concretas exactas.

**Ventajas:**
✓ Desacoplamiento de creación de objetos
✓ Centralización de lógica de creación
✓ Flexibilidad para cambiar implementaciones

### Factory Method (Simple)

#### **Ejemplo: Diferentes Tipos de Transporte**

```go
package transport

type Vehicle interface {
    Start() string
    Stop() string
    Capacity() int
}

type Car struct {
    brand string
}

func (c *Car) Start() string { return "Car engine starts" }
func (c *Car) Stop() string  { return "Car stops" }
func (c *Car) Capacity() int { return 4 }

type Truck struct {
    capacity int
}

func (t *Truck) Start() string { return "Truck engine roars" }
func (t *Truck) Stop() string  { return "Truck stops" }
func (t *Truck) Capacity() int { return t.capacity }

type Bus struct {
    seats int
}

func (b *Bus) Start() string { return "Bus door opens" }
func (b *Bus) Stop() string  { return "Bus stops" }
func (b *Bus) Capacity() int { return b.seats }

// Factory Method
func CreateVehicle(vehicleType string) (Vehicle, error) {
    switch vehicleType {
    case "car":
        return &Car{brand: "Toyota"}, nil
    case "truck":
        return &Truck{capacity: 1000}, nil
    case "bus":
        return &Bus{seats: 50}, nil
    default:
        return nil, fmt.Errorf("unknown vehicle type: %s", vehicleType)
    }
}

// Uso
func main() {
    vehicle, _ := CreateVehicle("car")
    fmt.Println(vehicle.Start())
    fmt.Println("Capacity:", vehicle.Capacity())
}
```

### Abstract Factory (Complejo)

#### **Ejemplo: UI Elements por Sistema Operativo**

```go
package ui

// Interfaces de productos
type Button interface {
    Render() string
}

type TextBox interface {
    Render() string
}

// Factory abstracta
type UIFactory interface {
    CreateButton() Button
    CreateTextBox() TextBox
}

// Implementación Windows
type WindowsButton struct{}
func (wb *WindowsButton) Render() string { return "[Windows Button]" }

type WindowsTextBox struct{}
func (wtb *WindowsTextBox) Render() string { return "[Windows TextBox]" }

type WindowsFactory struct{}
func (wf *WindowsFactory) CreateButton() Button { return &WindowsButton{} }
func (wf *WindowsFactory) CreateTextBox() TextBox { return &WindowsTextBox{} }

// Implementación macOS
type MacButton struct{}
func (mb *MacButton) Render() string { return "[Mac Button]" }

type MacTextBox struct{}
func (mtb *MacTextBox) Render() string { return "[Mac TextBox]" }

type MacFactory struct{}
func (mf *MacFactory) CreateButton() Button { return &MacButton{} }
func (mf *MacFactory) CreateTextBox() TextBox { return &MacTextBox{} }

// Factory selector
func GetUIFactory(os string) UIFactory {
    switch os {
    case "windows":
        return &WindowsFactory{}
    case "mac":
        return &MacFactory{}
    default:
        return &WindowsFactory{}
    }
}

// Uso
func main() {
    factory := GetUIFactory("mac")
    button := factory.CreateButton()
    textbox := factory.CreateTextBox()
    
    fmt.Println(button.Render())
    fmt.Println(textbox.Render())
}
```

### Functional Factory Pattern

Go permite usar funciones como factories:

```go
package handler

type Handler interface {
    Handle(req interface{}) interface{}
}

type GetHandler struct {
    data map[string]string
}

func (h *GetHandler) Handle(req interface{}) interface{} {
    return "GET response"
}

type PostHandler struct {
    logger interface{}
}

func (h *PostHandler) Handle(req interface{}) interface{} {
    return "POST response"
}

// Factory como type (muy idiomático Go)
type HandlerFactory func(config map[string]interface{}) (Handler, error)

var handlers = map[string]HandlerFactory{
    "GET": func(config map[string]interface{}) (Handler, error) {
        return &GetHandler{
            data: config["data"].(map[string]string),
        }, nil
    },
    "POST": func(config map[string]interface{}) (Handler, error) {
        return &PostHandler{
            logger: config["logger"],
        }, nil
    },
}

func CreateHandler(method string, config map[string]interface{}) (Handler, error) {
    if factory, ok := handlers[method]; ok {
        return factory(config)
    }
    return nil, fmt.Errorf("unknown handler: %s", method)
}
```

### Factory en Packages Estándar de Go

#### **Ejemplo: Decoder Factory**

```go
package main

import (
    "encoding/json"
    "encoding/xml"
    "io"
)

type Decoder interface {
    Decode(v interface{}) error
}

// Factory para decoders
func NewDecoder(format string, r io.Reader) Decoder {
    switch format {
    case "json":
        return json.NewDecoder(r)
    case "xml":
        return xml.NewDecoder(r)
    default:
        return nil
    }
}
```

#### **Ejemplo: Database Drivers**

```go
// Go implementa Factory implícitamente con driver registración:
import "database/sql"
import _ "github.com/lib/pq" // PostgreSQL driver

func main() {
    // Factory implícita en sql.Open
    db, err := sql.Open("postgres", "...")
    // Internamente, sql.Open busca driver registrado
}
```

---

## 44.4 - Patrón Builder {#patrón-builder}

### Concepto

El patrón Builder separa la construcción de un objeto complejo de su representación, permitiendo crear diferentes representaciones paso a paso.

**Problema que resuelve:**
```
❌ Constructores con muchos parámetros
❌ Parámetros opcionales
❌ Objetos complejos con múltiples dependencias
```

### Builder Básico

#### **Ejemplo: Configuration Object**

```go
package config

type ServerConfig struct {
    Host         string
    Port         int
    TLS          bool
    MaxConnections int
    Timeout      int
    LogLevel     string
}

type ServerConfigBuilder struct {
    config *ServerConfig
}

func NewServerConfigBuilder() *ServerConfigBuilder {
    return &ServerConfigBuilder{
        config: &ServerConfig{
            Host: "localhost",
            Port: 8080,
            TLS: false,
            MaxConnections: 100,
            Timeout: 30,
            LogLevel: "info",
        },
    }
}

// Métodos fluidos
func (b *ServerConfigBuilder) WithHost(host string) *ServerConfigBuilder {
    b.config.Host = host
    return b
}

func (b *ServerConfigBuilder) WithPort(port int) *ServerConfigBuilder {
    b.config.Port = port
    return b
}

func (b *ServerConfigBuilder) WithTLS(enabled bool) *ServerConfigBuilder {
    b.config.TLS = enabled
    return b
}

func (b *ServerConfigBuilder) WithMaxConnections(max int) *ServerConfigBuilder {
    b.config.MaxConnections = max
    return b
}

func (b *ServerConfigBuilder) WithTimeout(seconds int) *ServerConfigBuilder {
    b.config.Timeout = seconds
    return b
}

func (b *ServerConfigBuilder) WithLogLevel(level string) *ServerConfigBuilder {
    b.config.LogLevel = level
    return b
}

// Construir objeto final
func (b *ServerConfigBuilder) Build() *ServerConfig {
    return b.config
}

// Uso
func main() {
    config := NewServerConfigBuilder().
        WithHost("api.example.com").
        WithPort(443).
        WithTLS(true).
        WithMaxConnections(500).
        WithLogLevel("debug").
        Build()
    
    fmt.Printf("%+v\n", config)
}
```

### Builder con Validación

```go
package database

import "fmt"

type DatabaseConfig struct {
    Driver   string
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type DBConfigBuilder struct {
    config *DatabaseConfig
    errors []string
}

func NewDBConfigBuilder() *DBConfigBuilder {
    return &DBConfigBuilder{
        config: &DatabaseConfig{Port: 5432},
        errors: []string{},
    }
}

func (b *DBConfigBuilder) WithDriver(driver string) *DBConfigBuilder {
    if driver == "" {
        b.errors = append(b.errors, "driver cannot be empty")
    } else {
        b.config.Driver = driver
    }
    return b
}

func (b *DBConfigBuilder) WithHost(host string) *DBConfigBuilder {
    if host == "" {
        b.errors = append(b.errors, "host cannot be empty")
    } else {
        b.config.Host = host
    }
    return b
}

func (b *DBConfigBuilder) WithPort(port int) *DBConfigBuilder {
    if port <= 0 || port > 65535 {
        b.errors = append(b.errors, "port must be between 1 and 65535")
    } else {
        b.config.Port = port
    }
    return b
}

func (b *DBConfigBuilder) WithCredentials(user, pass string) *DBConfigBuilder {
    if user == "" {
        b.errors = append(b.errors, "username cannot be empty")
    } else {
        b.config.Username = user
    }
    if pass == "" {
        b.errors = append(b.errors, "password cannot be empty")
    } else {
        b.config.Password = pass
    }
    return b
}

func (b *DBConfigBuilder) WithDatabase(db string) *DBConfigBuilder {
    if db == "" {
        b.errors = append(b.errors, "database name cannot be empty")
    } else {
        b.config.Database = db
    }
    return b
}

// Validar antes de construir
func (b *DBConfigBuilder) Build() (*DatabaseConfig, error) {
    if len(b.errors) > 0 {
        return nil, fmt.Errorf("validation errors: %v", b.errors)
    }
    return b.config, nil
}

// Uso con manejo de errores
func main() {
    config, err := NewDBConfigBuilder().
        WithDriver("postgres").
        WithHost("db.example.com").
        WithPort(5432).
        WithCredentials("user", "pass").
        WithDatabase("mydb").
        Build()
    
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("%+v\n", config)
}
```

### Builder Alternativo: Opciones Pattern (Más Idiomático)

```go
package http

type ClientOptions struct {
    Timeout    int
    RetryCount int
    UserAgent  string
    Headers    map[string]string
}

// Tipo para función opción
type ClientOption func(*ClientOptions)

// Funciones opción
func WithTimeout(timeout int) ClientOption {
    return func(opts *ClientOptions) {
        opts.Timeout = timeout
    }
}

func WithRetryCount(count int) ClientOption {
    return func(opts *ClientOptions) {
        opts.RetryCount = count
    }
}

func WithUserAgent(ua string) ClientOption {
    return func(opts *ClientOptions) {
        opts.UserAgent = ua
    }
}

func WithHeaders(headers map[string]string) ClientOption {
    return func(opts *ClientOptions) {
        opts.Headers = headers
    }
}

type Client struct {
    opts ClientOptions
}

func NewClient(options ...ClientOption) *Client {
    opts := ClientOptions{
        Timeout:    30,
        RetryCount: 3,
        UserAgent:  "Go-Client/1.0",
        Headers:    make(map[string]string),
    }
    
    for _, opt := range options {
        opt(&opts)
    }
    
    return &Client{opts: opts}
}

// Uso
func main() {
    client := NewClient(
        WithTimeout(60),
        WithRetryCount(5),
        WithUserAgent("MyApp/2.0"),
        WithHeaders(map[string]string{
            "Authorization": "Bearer token",
        }),
    )
    
    fmt.Printf("%+v\n", client.opts)
}
```

### Builder en Packages Estándar

#### **Ejemplo: HTTP Request Builder**

```go
package main

import (
    "net/http"
    "bytes"
)

func main() {
    // Builder patrón implícito en http.Client
    client := &http.Client{}
    
    // Request builder pattern
    req, _ := http.NewRequest(
        http.MethodPost,
        "https://api.example.com/data",
        bytes.NewBufferString(`{"key":"value"}`),
    )
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer token")
    
    resp, _ := client.Do(req)
    defer resp.Body.Close()
}
```

---

## 44.5 - Patrón Decorator {#patrón-decorator}

### Concepto

El patrón Decorator añade dinámicamente nuevas responsabilidades a objetos sin modificar su estructura.

**Ventajas:**
✓ Añadir funcionalidad sin herencia
✓ Composición sobre extensión
✓ Responsabilidad única

### Decorator con Funciones

```go
package decorators

import (
    "fmt"
    "log"
    "time"
)

// Tipo para función que decora
type Handler func(string) string

// Decorator: Logging
func WithLogging(h Handler) Handler {
    return func(input string) string {
        log.Printf("[LOG] Iniciando con: %s", input)
        result := h(input)
        log.Printf("[LOG] Completado: %s", result)
        return result
    }
}

// Decorator: Timing
func WithTiming(h Handler) Handler {
    return func(input string) string {
        start := time.Now()
        result := h(input)
        fmt.Printf("[TIMING] Duración: %v\n", time.Since(start))
        return result
    }
}

// Decorator: Caching
func WithCaching(h Handler) Handler {
    cache := make(map[string]string)
    return func(input string) string {
        if result, ok := cache[input]; ok {
            fmt.Println("[CACHE] Hit:", input)
            return result
        }
        result := h(input)
        cache[input] = result
        fmt.Println("[CACHE] Miss:", input)
        return result
    }
}

// Función base
func ProcessData(data string) string {
    return fmt.Sprintf("Processed: %s", data)
}

// Composición de decorators
func main() {
    handler := WithLogging(WithTiming(WithCaching(ProcessData)))
    
    result1 := handler("data1")
    result2 := handler("data1") // Desde cache
    result3 := handler("data2")
    
    fmt.Println(result1, result2, result3)
}
```

### Decorator con Interfaces

```go
package middleware

import (
    "fmt"
    "net/http"
    "time"
)

// Interfaz base
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type LoggingHandler struct {
    next http.Handler
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    fmt.Printf("[LOG] %s %s\n", r.Method, r.URL.Path)
    h.next.ServeHTTP(w, r)
    fmt.Printf("[LOG] Duration: %v\n", time.Since(start))
}

type AuthHandler struct {
    next http.Handler
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("Authorization")
    if token == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    fmt.Println("[AUTH] Token:", token)
    h.next.ServeHTTP(w, r)
}

type RateLimitHandler struct {
    next  http.Handler
    limit int
}

func (h *RateLimitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("[RATE_LIMIT] Requests allowed: %d\n", h.limit)
    h.next.ServeHTTP(w, r)
}

// Composición
func WrapHandler(h http.Handler) http.Handler {
    h = &LoggingHandler{next: h}
    h = &AuthHandler{next: h}
    h = &RateLimitHandler{next: h, limit: 100}
    return h
}

// Uso
func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Response"))
    })
    
    wrapped := WrapHandler(handler)
    
    http.Handle("/api/data", wrapped)
    http.ListenAndServe(":8080", nil)
}
```

### Decorator Real: Database Wrapper

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"
)

// Interfaz de base de datos
type DB interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    Exec(query string, args ...interface{}) (sql.Result, error)
}

// Implementación real
type RealDB struct {
    *sql.DB
}

// Decorator: Logging
type LoggingDB struct {
    db DB
}

func (ldb *LoggingDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
    log.Printf("[DB QUERY] %s %v", query, args)
    return ldb.db.Query(query, args...)
}

func (ldb *LoggingDB) Exec(query string, args ...interface{}) (sql.Result, error) {
    log.Printf("[DB EXEC] %s %v", query, args)
    return ldb.db.Exec(query, args...)
}

// Decorator: Timing
type TimingDB struct {
    db DB
}

func (tdb *TimingDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
    start := time.Now()
    result, err := tdb.db.Query(query, args...)
    fmt.Printf("[DB TIMING] Query took: %v\n", time.Since(start))
    return result, err
}

func (tdb *TimingDB) Exec(query string, args ...interface{}) (sql.Result, error) {
    start := time.Now()
    result, err := tdb.db.Exec(query, args...)
    fmt.Printf("[DB TIMING] Exec took: %v\n", time.Since(start))
    return result, err
}

// Composición
func NewDecoratedDB(sqlDB *sql.DB) DB {
    realDB := &RealDB{sqlDB}
    decorated := DB(realDB)
    decorated = &LoggingDB{db: decorated}
    decorated = &TimingDB{db: decorated}
    return decorated
}
```

---

## 44.6 - Patrón Strategy {#patrón-strategy}

### Concepto

El patrón Strategy encapsula un conjunto de algoritmos intercambiables permitiendo seleccionar el algoritmo en tiempo de ejecución.

**Problema que resuelve:**
```
❌ Switch/case gigantes con lógica diferente
❌ Cambiar algoritmo en tiempo de ejecución
❌ Testear algoritmos independientemente
```

### Strategy Básico

#### **Ejemplo: Ordenamiento Intercambiable**

```go
package sorting

type SortStrategy interface {
    Sort(data []int)
}

// BubbleSort Strategy
type BubbleSort struct{}

func (bs *BubbleSort) Sort(data []int) {
    n := len(data)
    for i := 0; i < n; i++ {
        for j := 0; j < n-i-1; j++ {
            if data[j] > data[j+1] {
                data[j], data[j+1] = data[j+1], data[j]
            }
        }
    }
}

// QuickSort Strategy
type QuickSort struct{}

func (qs *QuickSort) Sort(data []int) {
    quickSort(data, 0, len(data)-1)
}

func quickSort(data []int, low, high int) {
    if low < high {
        pi := partition(data, low, high)
        quickSort(data, low, pi-1)
        quickSort(data, pi+1, high)
    }
}

func partition(data []int, low, high int) int {
    pivot := data[high]
    i := low - 1
    for j := low; j < high; j++ {
        if data[j] < pivot {
            i++
            data[i], data[j] = data[j], data[i]
        }
    }
    data[i+1], data[high] = data[high], data[i+1]
    return i + 1
}

// Contexto
type Sorter struct {
    strategy SortStrategy
}

func (s *Sorter) SetStrategy(strategy SortStrategy) {
    s.strategy = strategy
}

func (s *Sorter) SortData(data []int) {
    if s.strategy != nil {
        s.strategy.Sort(data)
    }
}

// Uso
func main() {
    data1 := []int{5, 2, 8, 1, 9}
    sorter := &Sorter{strategy: &BubbleSort{}}
    sorter.SortData(data1)
    fmt.Println("BubbleSort:", data1)
    
    data2 := []int{5, 2, 8, 1, 9}
    sorter.SetStrategy(&QuickSort{})
    sorter.SortData(data2)
    fmt.Println("QuickSort:", data2)
}
```

### Strategy con Funciones

```go
package payment

// Estrategia de pago como función
type PaymentStrategy func(amount float64) bool

// Diferentes métodos de pago
func CreditCard(amount float64) bool {
    fmt.Printf("Pagando $%.2f con Tarjeta de Crédito\n", amount)
    return true
}

func PayPal(amount float64) bool {
    fmt.Printf("Pagando $%.2f con PayPal\n", amount)
    return true
}

func Bitcoin(amount float64) bool {
    if amount > 10000 {
        fmt.Println("Bitcoin no disponible para montos > $10000")
        return false
    }
    fmt.Printf("Pagando $%.2f con Bitcoin\n", amount)
    return true
}

// Procesador de pagos
type PaymentProcessor struct {
    strategy PaymentStrategy
}

func (p *PaymentProcessor) SetPaymentMethod(strategy PaymentStrategy) {
    p.strategy = strategy
}

func (p *PaymentProcessor) Process(amount float64) bool {
    return p.strategy(amount)
}

// Uso
func main() {
    processor := &PaymentProcessor{}
    
    processor.SetPaymentMethod(CreditCard)
    processor.Process(100.00)
    
    processor.SetPaymentMethod(Bitcoin)
    processor.Process(5000.00)
    processor.Process(15000.00)
}
```

### Strategy en Operaciones Complejas

```go
package dataprocessing

import (
    "fmt"
    "sort"
)

type CompressionStrategy interface {
    Compress(data []byte) ([]byte, error)
}

type ZipCompression struct{}
func (zc *ZipCompression) Compress(data []byte) ([]byte, error) {
    return []byte("zip_compressed"), nil
}

type GzipCompression struct{}
func (gc *GzipCompression) Compress(data []byte) ([]byte, error) {
    return []byte("gzip_compressed"), nil
}

type EncryptionStrategy interface {
    Encrypt(data []byte) ([]byte, error)
}

type AESEncryption struct{}
func (ae *AESEncryption) Encrypt(data []byte) ([]byte, error) {
    return []byte("aes_encrypted"), nil
}

type RSAEncryption struct{}
func (re *RSAEncryption) Encrypt(data []byte) ([]byte, error) {
    return []byte("rsa_encrypted"), nil
}

// Pipeline con múltiples strategies
type DataPipeline struct {
    compression CompressionStrategy
    encryption  EncryptionStrategy
}

func (dp *DataPipeline) SetCompression(cs CompressionStrategy) {
    dp.compression = cs
}

func (dp *DataPipeline) SetEncryption(es EncryptionStrategy) {
    dp.encryption = es
}

func (dp *DataPipeline) Process(data []byte) ([]byte, error) {
    compressed, err := dp.compression.Compress(data)
    if err != nil {
        return nil, err
    }
    
    encrypted, err := dp.encryption.Encrypt(compressed)
    if err != nil {
        return nil, err
    }
    
    return encrypted, nil
}

// Uso
func main() {
    pipeline := &DataPipeline{
        compression: &GzipCompression{},
        encryption:  &AESEncryption{},
    }
    
    result, _ := pipeline.Process([]byte("sensitive data"))
    fmt.Printf("Processed: %v\n", result)
}
```

---

## 44.7 - Patrón Observer {#patrón-observer}

### Concepto

El patrón Observer define una dependencia uno-a-muchos entre objetos. Cuando un objeto cambia de estado, todos sus observadores son notificados automáticamente.

**Casos de uso:**
✓ Event listeners
✓ Pub/Sub systems
✓ Live updates
✓ Change notifications

### Observer Básico

```go
package events

import (
    "fmt"
    "sync"
)

// Observer interface
type Observer interface {
    Update(data interface{})
}

// Subject que es observado
type Subject struct {
    mu        sync.RWMutex
    observers []Observer
    state     interface{}
}

// Registrar observador
func (s *Subject) Attach(obs Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.observers = append(s.observers, obs)
}

// Desregistrar observador
func (s *Subject) Detach(obs Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for i, o := range s.observers {
        if o == obs {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            break
        }
    }
}

// Notificar a todos los observadores
func (s *Subject) Notify() {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    for _, obs := range s.observers {
        obs.Update(s.state)
    }
}

// Cambiar estado
func (s *Subject) SetState(state interface{}) {
    s.mu.Lock()
    s.state = state
    s.mu.Unlock()
    s.Notify()
}

// Implementación de observador
type ConcreteObserver struct {
    name string
}

func (co *ConcreteObserver) Update(data interface{}) {
    fmt.Printf("[%s] Observé cambio: %v\n", co.name, data)
}

// Uso
func main() {
    subject := &Subject{}
    
    obs1 := &ConcreteObserver{name: "Observer1"}
    obs2 := &ConcreteObserver{name: "Observer2"}
    
    subject.Attach(obs1)
    subject.Attach(obs2)
    
    subject.SetState("Nuevo estado")
    
    subject.Detach(obs1)
    subject.SetState("Otro estado")
}
```

### Observer con Funciones (Más Idiomático)

```go
package eventbus

import (
    "fmt"
    "sync"
)

// Callback para eventos
type EventCallback func(event interface{})

// Bus de eventos
type EventBus struct {
    mu        sync.RWMutex
    callbacks map[string][]EventCallback
}

func NewEventBus() *EventBus {
    return &EventBus{
        callbacks: make(map[string][]EventCallback),
    }
}

// Subscribirse a un evento
func (eb *EventBus) Subscribe(eventType string, callback EventCallback) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    eb.callbacks[eventType] = append(eb.callbacks[eventType], callback)
}

// Publicar evento
func (eb *EventBus) Publish(eventType string, data interface{}) {
    eb.mu.RLock()
    callbacks, ok := eb.callbacks[eventType]
    eb.mu.RUnlock()
    
    if !ok {
        return
    }
    
    for _, callback := range callbacks {
        go callback(data)
    }
}

// Uso
func main() {
    bus := NewEventBus()
    
    // Suscribirse a eventos
    bus.Subscribe("user.created", func(event interface{}) {
        fmt.Println("Usuario creado:", event)
    })
    
    bus.Subscribe("user.deleted", func(event interface{}) {
        fmt.Println("Usuario eliminado:", event)
    })
    
    // Publicar eventos
    bus.Publish("user.created", map[string]string{"id": "123", "name": "John"})
    bus.Publish("user.deleted", map[string]string{"id": "123"})
}
```

### Observer Real: Temperature Monitor

```go
package monitoring

import (
    "fmt"
    "math"
)

// Datos del evento
type TemperatureEvent struct {
    Value     float64
    Timestamp int64
    Sensor    string
}

// Observer
type TemperatureObserver interface {
    OnTemperatureChange(event TemperatureEvent)
}

// Sensor de temperatura
type TemperatureSensor struct {
    observers []TemperatureObserver
    current   float64
}

func (ts *TemperatureSensor) RegisterObserver(obs TemperatureObserver) {
    ts.observers = append(ts.observers, obs)
}

func (ts *TemperatureSensor) NotifyObservers(event TemperatureEvent) {
    for _, obs := range ts.observers {
        obs.OnTemperatureChange(event)
    }
}

func (ts *TemperatureSensor) SetTemperature(temp float64) {
    ts.current = temp
    event := TemperatureEvent{
        Value:     temp,
        Timestamp: 1234567890,
        Sensor:    "Main",
    }
    ts.NotifyObservers(event)
}

// Observador: Logger
type TemperatureLogger struct{}

func (tl *TemperatureLogger) OnTemperatureChange(event TemperatureEvent) {
    fmt.Printf("[LOG] Temperatura: %.1f°C\n", event.Value)
}

// Observador: Alarma
type TemperatureAlarm struct {
    threshold float64
}

func (ta *TemperatureAlarm) OnTemperatureChange(event TemperatureEvent) {
    if event.Value > ta.threshold {
        fmt.Printf("[ALARM] ¡Temperatura alta! %.1f°C\n", event.Value)
    }
}

// Observador: Estadísticas
type TemperatureStats struct {
    readings []float64
}

func (ts *TemperatureStats) OnTemperatureChange(event TemperatureEvent) {
    ts.readings = append(ts.readings, event.Value)
    
    avg := 0.0
    for _, r := range ts.readings {
        avg += r
    }
    avg /= float64(len(ts.readings))
    
    fmt.Printf("[STATS] Promedio: %.1f°C, Lecturas: %d\n", avg, len(ts.readings))
}

// Uso
func main() {
    sensor := &TemperatureSensor{}
    
    sensor.RegisterObserver(&TemperatureLogger{})
    sensor.RegisterObserver(&TemperatureAlarm{threshold: 30.0})
    sensor.RegisterObserver(&TemperatureStats{})
    
    sensor.SetTemperature(25.5)
    sensor.SetTemperature(32.0)
    sensor.SetTemperature(28.0)
}
```

---

## 44.8 - Patrón Adapter {#patrón-adapter}

### Concepto

El patrón Adapter convierte la interfaz de una clase en otra interfaz que el cliente espera.

**Ventajas:**
✓ Integrar APIs incompatibles
✓ Mantener compatibilidad hacia atrás
✓ Trabajar con librerías de terceros

### Adapter Simple

#### **Ejemplo: Adaptador de Formatos**

```go
package adapters

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
)

// Interfaz esperada por el cliente
type DataWriter interface {
    WriteData(data interface{}) ([]byte, error)
}

// Implementación para JSON
type JSONWriter struct{}

func (jw *JSONWriter) WriteData(data interface{}) ([]byte, error) {
    return json.Marshal(data)
}

// Sistema legado: XML
type LegacyXMLSystem struct{}

func (lxs *LegacyXMLSystem) ExportToXML(obj interface{}) ([]byte, error) {
    return xml.Marshal(obj)
}

// Adaptador: XML a interfaz DataWriter
type XMLAdapter struct {
    legacySystem *LegacyXMLSystem
}

func (xa *XMLAdapter) WriteData(data interface{}) ([]byte, error) {
    return xa.legacySystem.ExportToXML(data)
}

// Uso
func main() {
    data := map[string]string{"name": "Alice", "age": "30"}
    
    var writer DataWriter
    
    writer = &JSONWriter{}
    result, _ := writer.WriteData(data)
    fmt.Println("JSON:", string(result))
    
    writer = &XMLAdapter{legacySystem: &LegacyXMLSystem{}}
    result, _ = writer.WriteData(data)
    fmt.Println("XML:", string(result))
}
```

### Adapter para Diferentes Interfaces

```go
package httpclients

import "fmt"

// Interfaz esperada: HTTP Client
type HTTPClient interface {
    Get(url string) (string, error)
    Post(url string, data string) (string, error)
}

// Cliente estándar
type StandardClient struct{}

func (sc *StandardClient) Get(url string) (string, error) {
    return fmt.Sprintf("GET %s", url), nil
}

func (sc *StandardClient) Post(url string, data string) (string, error) {
    return fmt.Sprintf("POST %s: %s", url, data), nil
}

// Sistema legado: usa métodos diferentes
type LegacyHTTPSystem struct{}

func (lhs *LegacyHTTPSystem) HttpRequest(method, url, body string) string {
    return fmt.Sprintf("%s request to %s with %s", method, url, body)
}

// Adaptador
type LegacyHTTPAdapter struct {
    legacySystem *LegacyHTTPSystem
}

func (lha *LegacyHTTPAdapter) Get(url string) (string, error) {
    return lha.legacySystem.HttpRequest("GET", url, ""), nil
}

func (lha *LegacyHTTPAdapter) Post(url string, data string) (string, error) {
    return lha.legacySystem.HttpRequest("POST", url, data), nil
}

// Función que usa HTTPClient
func MakeRequest(client HTTPClient, url string) {
    result, _ := client.Get(url)
    fmt.Println(result)
}

// Uso
func main() {
    MakeRequest(&StandardClient{}, "https://api.example.com")
    
    adapter := &LegacyHTTPAdapter{
        legacySystem: &LegacyHTTPSystem{},
    }
    MakeRequest(adapter, "https://api.example.com")
}
```

### Adapter Chain

```go
package middleware

import (
    "fmt"
    "net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request)

// Sistema 1: Maneja logging
type LoggingSystem struct{}

func (ls *LoggingSystem) LogRequest(h Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[LOG] %s %s\n", r.Method, r.URL)
        h(w, r)
    }
}

// Sistema 2: Maneja autenticación
type AuthSystem struct{}

func (as *AuthSystem) CheckAuth(h Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        h(w, r)
    }
}

// Adaptador para composición
type MiddlewareAdapter struct {
    logger *LoggingSystem
    auth   *AuthSystem
}

func (ma *MiddlewareAdapter) Chain(h Handler) Handler {
    h = ma.logger.LogRequest(h)
    h = ma.auth.CheckAuth(h)
    return h
}

// Handler base
func ApiHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Success"))
}

// Uso
func main() {
    adapter := &MiddlewareAdapter{
        logger: &LoggingSystem{},
        auth:   &AuthSystem{},
    }
    
    http.HandleFunc("/api/data", adapter.Chain(ApiHandler))
    http.ListenAndServe(":8080", nil)
}
```

---

## 44.9 - Patrón Repository {#patrón-repository}

### Concepto

El patrón Repository abstrae la lógica de acceso a datos, proporcionando una colección similar a interfaz para acceder a objetos.

**Ventajas:**
✓ Separación de concerns
✓ Facilita testing
✓ Cambiar BD fácilmente
✓ Centralizar queries

### Repository Básico

```go
package repository

import (
    "database/sql"
    "errors"
    "fmt"
)

type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

// Interfaz del repositorio
type UserRepository interface {
    FindByID(id int) (*User, error)
    FindAll() ([]*User, error)
    Save(user *User) error
    Delete(id int) error
    Update(user *User) error
}

// Implementación SQL
type SQLUserRepository struct {
    db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
    return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) FindByID(id int) (*User, error) {
    user := &User{}
    err := r.db.QueryRow(
        "SELECT id, name, email, age FROM users WHERE id = $1",
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    return user, err
}

func (r *SQLUserRepository) FindAll() ([]*User, error) {
    rows, err := r.db.Query("SELECT id, name, email, age FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*User
    for rows.Next() {
        user := &User{}
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func (r *SQLUserRepository) Save(user *User) error {
    err := r.db.QueryRow(
        "INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id",
        user.Name, user.Email, user.Age,
    ).Scan(&user.ID)
    return err
}

func (r *SQLUserRepository) Update(user *User) error {
    _, err := r.db.Exec(
        "UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4",
        user.Name, user.Email, user.Age, user.ID,
    )
    return err
}

func (r *SQLUserRepository) Delete(id int) error {
    _, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
    return err
}
```

### Repository en Memoria (Para Testing)

```go
package repository

import (
    "errors"
    "sync"
)

// Implementación en memoria para testing
type InMemoryUserRepository struct {
    mu    sync.RWMutex
    users map[int]*User
    nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users:  make(map[int]*User),
        nextID: 1,
    }
}

func (r *InMemoryUserRepository) FindByID(id int) (*User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    user, ok := r.users[id]
    if !ok {
        return nil, errors.New("user not found")
    }
    return user, nil
}

func (r *InMemoryUserRepository) FindAll() ([]*User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var users []*User
    for _, user := range r.users {
        users = append(users, user)
    }
    return users, nil
}

func (r *InMemoryUserRepository) Save(user *User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    user.ID = r.nextID
    r.users[r.nextID] = user
    r.nextID++
    return nil
}

func (r *InMemoryUserRepository) Update(user *User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.users[user.ID]; !ok {
        return errors.New("user not found")
    }
    r.users[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) Delete(id int) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.users[id]; !ok {
        return errors.New("user not found")
    }
    delete(r.users, id)
    return nil
}
```

### Service que usa Repository

```go
package service

import "fmt"

// Service usa Repository
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (us *UserService) CreateUser(name, email string, age int) (*User, error) {
    user := &User{
        Name:  name,
        Email: email,
        Age:   age,
    }
    
    if err := us.repo.Save(user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    return user, nil
}

func (us *UserService) GetUser(id int) (*User, error) {
    return us.repo.FindByID(id)
}

func (us *UserService) ListUsers() ([]*User, error) {
    return us.repo.FindAll()
}

func (us *UserService) DeleteUser(id int) error {
    return us.repo.Delete(id)
}

// Testing
func TestUserService() {
    repo := NewInMemoryUserRepository()
    service := NewUserService(repo)
    
    user, _ := service.CreateUser("Alice", "alice@example.com", 30)
    fmt.Printf("Created: %+v\n", user)
    
    retrieved, _ := service.GetUser(user.ID)
    fmt.Printf("Retrieved: %+v\n", retrieved)
}
```

---

## 44.10 - Inyección de Dependencias {#inyección-de-dependencias}

### Concepto

La Inyección de Dependencias (DI) es un patrón donde las dependencias de un objeto se proporcionan externamente en lugar de crearlas internamente.

**Beneficios:**
✓ Desacoplamiento
✓ Testabilidad
✓ Flexibilidad
✓ Inversión de Control

### Constructor Injection

```go
package di

import "fmt"

// Dependencias
type Logger interface {
    Log(msg string)
}

type Database interface {
    Query(query string) string
}

// Implementaciones
type ConsoleLogger struct{}

func (cl *ConsoleLogger) Log(msg string) {
    fmt.Println("[LOG]", msg)
}

type MockDatabase struct{}

func (md *MockDatabase) Query(query string) string {
    return "mock result"
}

// Servicio que recibe dependencias
type UserService struct {
    logger Logger
    db     Database
}

// Constructor Injection
func NewUserService(logger Logger, db Database) *UserService {
    return &UserService{
        logger: logger,
        db:     db,
    }
}

func (us *UserService) GetUser(id string) {
    us.logger.Log(fmt.Sprintf("Fetching user %s", id))
    result := us.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %s", id))
    us.logger.Log(fmt.Sprintf("Result: %s", result))
}

// Uso
func main() {
    logger := &ConsoleLogger{}
    db := &MockDatabase{}
    
    service := NewUserService(logger, db)
    service.GetUser("123")
}
```

### Container de DI

```go
package container

import (
    "fmt"
    "sync"
)

// Contenedor simple
type Container struct {
    mu        sync.RWMutex
    services  map[string]interface{}
    factories map[string]func() interface{}
}

func NewContainer() *Container {
    return &Container{
        services:  make(map[string]interface{}),
        factories: make(map[string]func() interface{}),
    }
}

// Registrar servicio singleton
func (c *Container) Register(name string, service interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.services[name] = service
}

// Registrar factory
func (c *Container) RegisterFactory(name string, factory func() interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.factories[name] = factory
}

// Obtener servicio
func (c *Container) Get(name string) (interface{}, error) {
    c.mu.RLock()
    
    if service, ok := c.services[name]; ok {
        c.mu.RUnlock()
        return service, nil
    }
    
    factory, ok := c.factories[name]
    c.mu.RUnlock()
    
    if !ok {
        return nil, fmt.Errorf("service not found: %s", name)
    }
    
    return factory(), nil
}

// Uso
func main() {
    container := NewContainer()
    
    // Registrar singleton
    logger := &ConsoleLogger{}
    container.Register("logger", logger)
    
    // Registrar factory
    container.RegisterFactory("db", func() interface{} {
        return &MockDatabase{}
    })
    
    // Usar
    logger, _ := container.Get("logger")
    db, _ := container.Get("db")
    
    service := NewUserService(logger.(Logger), db.(Database))
    service.GetUser("123")
}
```

### Functional Options Pattern

```go
package options

type Server struct {
    host    string
    port    int
    timeout int
    logger  Logger
}

// Opción como función
type ServerOption func(*Server)

func WithHost(host string) ServerOption {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout int) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func WithLogger(logger Logger) ServerOption {
    return func(s *Server) {
        s.logger = logger
    }
}

// Constructor con opciones
func NewServer(options ...ServerOption) *Server {
    s := &Server{
        host:    "localhost",
        port:    8080,
        timeout: 30,
        logger:  &ConsoleLogger{},
    }
    
    for _, opt := range options {
        opt(s)
    }
    
    return s
}

// Uso
func main() {
    server := NewServer(
        WithHost("api.example.com"),
        WithPort(443),
        WithTimeout(60),
        WithLogger(&ConsoleLogger{}),
    )
    
    fmt.Printf("%+v\n", server)
}
```

---

## 44.11 - Buenas Prácticas y Antipatterns {#buenas-prácticas-y-antipatterns}

### Cuándo Usar Patrones

#### **✓ Cuándo SÍ usar:**

1. **El problema es complejo**: Múltiples casos, lógica intrincada
2. **Beneficio es claro**: Testabilidad, mantenibilidad, flexibilidad
3. **Es idiomático**: Se adapta al lenguaje
4. **El equipo entiende**: Todos conocen el patrón

#### **✗ Cuándo NO usar:**

1. **YAGNI (You Aren't Gonna Need It)**: No lo necesitas ahora
2. **Over-engineering**: Demasiada complejidad para poco beneficio
3. **No es idiomático**: Fuerza lenguaje a otros paradigmas
4. **Confunde el código**: Más difícil de leer que beneficio da

### Antipatterns Comunes

#### **1. Singleton Innecesario**

```go
// ❌ MALO: Singleton cuando no lo necesitas
var db *Database
var once sync.Once

func GetDB() *Database {
    once.Do(func() {
        db = &Database{}
    })
    return db
}

// ✅ BIEN: Simplemente crea la instancia
var db = &Database{}

func init() {
    db = NewDatabase()
}
```

#### **2. Factory Innecesaria**

```go
// ❌ MALO: Factory para algo simple
func CreateUser(name string) *User {
    switch {
    case name == "":
        return &User{}
    default:
        return &User{Name: name}
    }
}

// ✅ BIEN: Constructor directo
func NewUser(name string) *User {
    return &User{Name: name}
}
```

#### **3. Over-Engineering con Interfaces**

```go
// ❌ MALO: Interfaz innecesaria
type StringWriterInterface interface {
    WriteString(s string) error
}

type StringWriter struct{}

func (sw *StringWriter) WriteString(s string) error {
    fmt.Println(s)
    return nil
}

// ✅ BIEN: Usa io.Writer que ya existe
func WriteString(w io.Writer, s string) error {
    _, err := w.WriteString(s)
    return err
}
```

#### **4. Estrategia Innecesaria**

```go
// ❌ MALO: Strategy para un solo caso
type SortStrategy interface {
    Sort(data []int)
}

type QuickSortStrategy struct{}

func (qss *QuickSortStrategy) Sort(data []int) {
    sort.Ints(data)
}

// ✅ BIEN: Usa la función directamente
sort.Ints(data)
```

### Go Antipatterns Específicos

#### **1. Goroutines sin Control**

```go
// ❌ MALO: Goroutines que pueden huir
func Process(items []string) {
    for _, item := range items {
        go func(i string) {
            // Sin forma de saber cuándo termina
            doSomething(i)
        }(item)
    }
}

// ✅ BIEN: Usar WaitGroup
func Process(items []string) {
    var wg sync.WaitGroup
    for _, item := range items {
        wg.Add(1)
        go func(i string) {
            defer wg.Done()
            doSomething(i)
        }(item)
    }
    wg.Wait()
}
```

#### **2. Context Ignorado**

```go
// ❌ MALO: Ignorar context
func FetchData(ctx context.Context, url string) string {
    // Ignora ctx completamente
    return makeRequest(url)
}

// ✅ BIEN: Respetar context
func FetchData(ctx context.Context, url string) (string, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    return httpClient.Do(req)
}
```

#### **3. Errors Ignorados**

```go
// ❌ MALO: Ignorar errores
func SaveUser(user *User) {
    db.Save(user) // ¿Y si falla?
    fmt.Println("User saved")
}

// ✅ BIEN: Manejar errores
func SaveUser(user *User) error {
    if err := db.Save(user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    fmt.Println("User saved")
    return nil
}
```

### Comparación: Go vs Otros Lenguajes

| Aspecto | Go | Java | Python |
|--------|----|----|--------|
| Patrones comunes | ~15 | ~23 | ~15 |
| Singleton necesario | Raro | Común | Raro |
| Factory Pattern | Sí (funciones) | Sí | Sí (factories) |
| Visitor Pattern | No | Sí | Sí |
| Herencia profunda | No | Común | Raro |
| Interfaces mínimas | Sí | No | Sí |
| Over-engineering | Bajo | Alto | Medio |

### Principios de Diseño en Go

#### **1. Simplicidad Primero**
```go
// Busca la solución más simple primero
// Agrega complejidad solo si es necesario
```

#### **2. Interfaces Pequeñas**
```go
// ✓ Una responsabilidad
type Reader interface {
    Read(p []byte) (n int, err error)
}

// ✗ Muchas responsabilidades
type DataProcessor interface {
    Read() []byte
    Write() error
    Validate() bool
    Transform() interface{}
}
```

#### **3. Composición sobre Herencia**
```go
// Go no tiene herencia, usa embedding
type Logger struct {
    output io.Writer
}

type Server struct {
    Logger // Composición
    port   int
}
```

#### **4. Testing First**
```go
// Escribe tests que deberían pasar
// Esto guía el diseño automáticamente
type Repository interface {
    Save(item Item) error
}
// Fácil testear con mock
```

---

## 44.12 - Ejercicios Progresivos {#ejercicios-progresivos}

### Ejercicio 1: Singleton - Logger Única

**Objetivo:** Implementar un logger singleton thread-safe que pueda ser usado desde múltiples goroutines.

**Requisitos:**
- Usar `sync.Once`
- Métodos: `Debug()`, `Info()`, `Warn()`, `Error()`
- Thread-safe
- Guarde historial de últimos 100 logs
- No duplicar logs idénticos consecutivos

**Solución esperada:**
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
)

type Logger struct {
    mu       sync.Mutex
    logs     []string
    level    LogLevel
    maxLogs  int
}

var (
    instance *Logger
    once     sync.Once
)

func GetLogger() *Logger {
    once.Do(func() {
        instance = &Logger{
            logs:    make([]string, 0),
            level:   INFO,
            maxLogs: 100,
        }
    })
    return instance
}

func (l *Logger) log(level string, msg string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    entry := fmt.Sprintf("[%s] %s - %s", time.Now().Format("15:04:05"), level, msg)
    
    if len(l.logs) > 0 && l.logs[len(l.logs)-1] == entry {
        return // No duplicar
    }
    
    l.logs = append(l.logs, entry)
    if len(l.logs) > l.maxLogs {
        l.logs = l.logs[1:]
    }
    
    fmt.Println(entry)
}

func (l *Logger) Debug(msg string) { l.log("DEBUG", msg) }
func (l *Logger) Info(msg string)  { l.log("INFO", msg) }
func (l *Logger) Warn(msg string)  { l.log("WARN", msg) }
func (l *Logger) Error(msg string) { l.log("ERROR", msg) }

func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            logger := GetLogger()
            logger.Info(fmt.Sprintf("Goroutine %d", id))
        }(i)
    }
    
    wg.Wait()
}
```

### Ejercicio 2: Factory - Creador de Handlers

**Objetivo:** Crear una factory que genere diferentes tipos de handlers HTTP.

**Requisitos:**
- Tipos: GET, POST, PUT, DELETE
- Cada tipo tiene comportamiento distinto
- Registry de handlers extensible
- Manejar tipo desconocido

**Solución esperada:**
```go
package main

import (
    "fmt"
    "net/http"
)

type Handler interface {
    Handle(data string) string
}

type GetHandler struct{}
func (gh *GetHandler) Handle(data string) string {
    return fmt.Sprintf("GET: %s", data)
}

type PostHandler struct{}
func (ph *PostHandler) Handle(data string) string {
    return fmt.Sprintf("POST: %s", data)
}

type PutHandler struct{}
func (ph *PutHandler) Handle(data string) string {
    return fmt.Sprintf("PUT: %s", data)
}

type DeleteHandler struct{}
func (dh *DeleteHandler) Handle(data string) string {
    return fmt.Sprintf("DELETE: %s", data)
}

type HandlerFactory func() Handler

var handlers = map[string]HandlerFactory{
    http.MethodGet:    func() Handler { return &GetHandler{} },
    http.MethodPost:   func() Handler { return &PostHandler{} },
    http.MethodPut:    func() Handler { return &PutHandler{} },
    http.MethodDelete: func() Handler { return &DeleteHandler{} },
}

func CreateHandler(method string) (Handler, error) {
    factory, ok := handlers[method]
    if !ok {
        return nil, fmt.Errorf("unknown method: %s", method)
    }
    return factory(), nil
}

func main() {
    for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
        handler, _ := CreateHandler(method)
        fmt.Println(handler.Handle("resource"))
    }
}
```

### Ejercicio 3: Builder - Constructor de Solicitudes HTTP

**Objetivo:** Construir solicitudes HTTP complejas usando el patrón Builder.

**Requisitos:**
- Método fluido
- Headers personalizados
- Query parameters
- Body y validación
- Valores por defecto sensatos

**Solución esperada:**
```go
package main

import (
    "fmt"
    "net/url"
)

type HTTPRequest struct {
    method  string
    url     string
    headers map[string]string
    query   map[string]string
    body    string
}

type RequestBuilder struct {
    request *HTTPRequest
}

func NewRequestBuilder() *RequestBuilder {
    return &RequestBuilder{
        request: &HTTPRequest{
            method:  "GET",
            headers: make(map[string]string),
            query:   make(map[string]string),
        },
    }
}

func (rb *RequestBuilder) Method(m string) *RequestBuilder {
    rb.request.method = m
    return rb
}

func (rb *RequestBuilder) URL(u string) *RequestBuilder {
    rb.request.url = u
    return rb
}

func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
    rb.request.headers[key] = value
    return rb
}

func (rb *RequestBuilder) Query(key, value string) *RequestBuilder {
    rb.request.query[key] = value
    return rb
}

func (rb *RequestBuilder) Body(b string) *RequestBuilder {
    rb.request.body = b
    return rb
}

func (rb *RequestBuilder) Build() (*HTTPRequest, error) {
    if rb.request.url == "" {
        return nil, fmt.Errorf("URL is required")
    }
    return rb.request, nil
}

func (r *HTTPRequest) String() string {
    qp := url.Values{}
    for k, v := range r.query {
        qp.Add(k, v)
    }
    
    fullURL := r.url
    if len(r.query) > 0 {
        fullURL += "?" + qp.Encode()
    }
    
    s := fmt.Sprintf("%s %s\n", r.method, fullURL)
    for k, v := range r.headers {
        s += fmt.Sprintf("%s: %s\n", k, v)
    }
    if r.body != "" {
        s += "\n" + r.body
    }
    return s
}

func main() {
    req, _ := NewRequestBuilder().
        Method("POST").
        URL("https://api.example.com/data").
        Header("Content-Type", "application/json").
        Header("Authorization", "Bearer token").
        Query("limit", "10").
        Query("offset", "0").
        Body(`{"name":"John"}`).
        Build()
    
    fmt.Println(req)
}
```

### Ejercicio 4: Strategy - Algoritmo Intercambiable

**Objetivo:** Implementar estrategias de procesamiento intercambiables para análisis de datos.

**Requisitos:**
- Estrategias: Sum, Average, Min, Max
- Intercambiable en runtime
- Agregar nueva estrategia fácilmente
- Procesar datasets grandes

**Solución esperada:**
```go
package main

import (
    "fmt"
)

type DataStrategy interface {
    Process(data []int) int
}

type SumStrategy struct{}
func (ss *SumStrategy) Process(data []int) int {
    sum := 0
    for _, v := range data {
        sum += v
    }
    return sum
}

type AverageStrategy struct{}
func (as *AverageStrategy) Process(data []int) int {
    if len(data) == 0 {
        return 0
    }
    sum := 0
    for _, v := range data {
        sum += v
    }
    return sum / len(data)
}

type MinStrategy struct{}
func (ms *MinStrategy) Process(data []int) int {
    if len(data) == 0 {
        return 0
    }
    min := data[0]
    for _, v := range data[1:] {
        if v < min {
            min = v
        }
    }
    return min
}

type MaxStrategy struct{}
func (ms *MaxStrategy) Process(data []int) int {
    if len(data) == 0 {
        return 0
    }
    max := data[0]
    for _, v := range data[1:] {
        if v > max {
            max = v
        }
    }
    return max
}

type Analyzer struct {
    strategy DataStrategy
}

func (a *Analyzer) SetStrategy(s DataStrategy) {
    a.strategy = s
}

func (a *Analyzer) Analyze(data []int) int {
    return a.strategy.Process(data)
}

func main() {
    data := []int{10, 20, 30, 40, 50}
    analyzer := &Analyzer{}
    
    analyzer.SetStrategy(&SumStrategy{})
    fmt.Println("Sum:", analyzer.Analyze(data))
    
    analyzer.SetStrategy(&AverageStrategy{})
    fmt.Println("Average:", analyzer.Analyze(data))
    
    analyzer.SetStrategy(&MinStrategy{})
    fmt.Println("Min:", analyzer.Analyze(data))
    
    analyzer.SetStrategy(&MaxStrategy{})
    fmt.Println("Max:", analyzer.Analyze(data))
}
```

### Ejercicio 5: Repository - Abstracción de Datos

**Objetivo:** Implementar patrón Repository con múltiples implementaciones (SQL e In-Memory).

**Requisitos:**
- Interfaz Repository
- Implementación SQL
- Implementación In-Memory
- Idéntica interfaz, diferente almacenamiento
- Testeable

**Solución esperada:**
```go
package main

import (
    "errors"
    "fmt"
    "sync"
)

type Product struct {
    ID    int
    Name  string
    Price float64
}

type ProductRepository interface {
    GetByID(id int) (*Product, error)
    GetAll() []*Product
    Save(product *Product) error
    Delete(id int) error
}

// Implementación In-Memory
type InMemoryRepository struct {
    mu       sync.RWMutex
    products map[int]*Product
    nextID   int
}

func NewInMemoryRepository() *InMemoryRepository {
    return &InMemoryRepository{
        products: make(map[int]*Product),
        nextID:   1,
    }
}

func (ir *InMemoryRepository) GetByID(id int) (*Product, error) {
    ir.mu.RLock()
    defer ir.mu.RUnlock()
    
    p, ok := ir.products[id]
    if !ok {
        return nil, errors.New("product not found")
    }
    return p, nil
}

func (ir *InMemoryRepository) GetAll() []*Product {
    ir.mu.RLock()
    defer ir.mu.RUnlock()
    
    var products []*Product
    for _, p := range ir.products {
        products = append(products, p)
    }
    return products
}

func (ir *InMemoryRepository) Save(product *Product) error {
    ir.mu.Lock()
    defer ir.mu.Unlock()
    
    if product.ID == 0 {
        product.ID = ir.nextID
        ir.nextID++
    }
    ir.products[product.ID] = product
    return nil
}

func (ir *InMemoryRepository) Delete(id int) error {
    ir.mu.Lock()
    defer ir.mu.Unlock()
    
    if _, ok := ir.products[id]; !ok {
        return errors.New("product not found")
    }
    delete(ir.products, id)
    return nil
}

// Servicio que usa Repository
type ProductService struct {
    repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
    return &ProductService{repo: repo}
}

func (ps *ProductService) CreateProduct(name string, price float64) *Product {
    product := &Product{Name: name, Price: price}
    ps.repo.Save(product)
    return product
}

func (ps *ProductService) ListProducts() []*Product {
    return ps.repo.GetAll()
}

func main() {
    repo := NewInMemoryRepository()
    service := NewProductService(repo)
    
    service.CreateProduct("Laptop", 999.99)
    service.CreateProduct("Mouse", 29.99)
    service.CreateProduct("Keyboard", 79.99)
    
    products := service.ListProducts()
    for _, p := range products {
        fmt.Printf("ID: %d, Name: %s, Price: $%.2f\n", p.ID, p.Name, p.Price)
    }
}
```

---

## Resumen

Los patrones de diseño son herramientas poderosas cuando se usan apropiadamente. Go rechaza cierta complejidad innecesaria de lenguajes como Java, preferiendo soluciones simples y directas.

**Recordar:**
- ✓ Patrones cuando resuelven problemas reales
- ✗ Patrones como "magic bullets"
- Go idiomático > GoF patterns
- Simplicidad > Over-engineering

**Próximos pasos:**
1. Estudiar librerías populares y sus patterns
2. Practicar implementación
3. Conocer antipatterns comunes
4. Entender cuándo aplicar y cuándo no

