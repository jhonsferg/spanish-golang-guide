# Capítulo 64: Database drivers y SQL de bajo nivel

**Versión:** 1.0 | **Nivel:** Avanzado | **Requisitos Previos:** Capítulos 1-63, SQL básico, Concurrencia

---

## Tabla de Contenidos

1. Introducción a database/sql
2. Database Drivers
3. Connection Management
4. Basic Operations
5. Advanced SQL Patterns
6. Error Handling
7. Transactions
8. Performance Optimization
9. Testing Raw SQL
10. Scanning & Mapping
11. Production Patterns & Best Practices

---

## 64.1 - Introducción a database/sql

### 64.1.1 Standard Library Overview

La librería estándar de Go proporciona el paquete `database/sql` como una abstracción sobre drivers de bases de datos específicos. Este paquete define interfaces y tipos genéricos que permiten trabajar con diferentes motores de bases de datos usando una API consistente.

```

                    Aplicación Go                               │

                           │
            ┌──────────────┴──────────────┐
            │   database/sql (stdlib)     │
            │  - DB connection pooling    │
            │  - Query interface          │
            │  - Row scanning             │
            │  - Transaction support      │
            └──────────────┬──────────────┘
                           │
        ┌──────────┬───────┼────────┬──────────────────┐
        │          │       │        │                  │
    ┌──▼──┐ ┌──▼──┐ ┌──▼──┐           ┌────▼────┐     ┌──▼
     │ pq │    │pgx  │ │mysql│ │ odbc│  ......   │Drivers  │
     └─────┘    └─────┘ └─────┘ └─────┘           │Custom   │
        │          │       │        │              └─────────┘
        ▼          ▼       ▼        ▼
    PostgreSQL   PostgreSQL MySQL  ODBC
```

El paquete `database/sql` actúa como un adaptador universal, permitiendo que tu código se mantenga independiente del motor de base de datos específico.

**Características principales:**

- **Type safety:** Los tipos Go se mapean directamente a tipos de base de datos
- **Connection pooling:** Gestión automática del pool de conexiones
- **Prepared statements:** Seguridad contra inyección SQL
- **Scanning flexible:** Conversión automática de tipos
- **Context support:** Control de timeouts y cancelaciones
- **Driver agnostic:** Cambia de driver sin cambiar código

### 64.1.2 database/sql vs ORMs

```go
// ============ RAW SQL (database/sql) ============
row := db.QueryRowContext(ctx, 
    "SELECT id, name, email FROM users WHERE id = $1", userID)
var u User
if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
    return nil, err
}

// ============ ORM (GORM) ============
var u User
if err := db.WithContext(ctx).First(&u, userID).Error; err != nil {
    return nil, err
}

// ============ sqlc (Code generation) ============
// Genera el código basado en queries SQL verificadas
u, err := q.GetUserByID(ctx, userID)
```

**Comparación exhaustiva:**

| Aspecto | Raw SQL | GORM | sqlc |
|--------|---------|------|------|
| Control | 100% | Limitado | 100% |
| Performance | Excelente | Bueno | Excelente |
| Query optimization | Manual | Automático (puede fallar) | Manual (verificado) |
| Type safety | Go types | Runtime | Compile-time |
| Learning curve | Medio-Alto | Bajo | Medio |
| Flexibilidad | Máxima | Restringida | Máxima (con verificación) |
| Debugging | Directo | Complejo | Directo |
| Portabilidad | Drivers específicos | Múltiples drivers | Drivers específicos |
| JSON handling | Manual | Automático | Automático (tipado) |
| Migraciones | Manual | Automáticas | Manual |

**Cuándo usar cada uno:**

```

 Elige RAW SQL si:                            │

 • Queries complejas (CTEs, window functions)│
 • Performance crítica                        │
 • Necesitas control total sobre SQL          │
 • Equipo experto en SQL                      │
 • Database-specific features                 │
 • Microservicios con queries simples         │



 Elige ORM (GORM) si:                        │

 • Múltiples drivers de DB                    │
 • Modelo de datos complejo                   │
 • Prototipado rápido                        │
 Equipo sin expertise en SQL               │ 
 • CRUD straightforward                       │



 Elige sqlc si:                              │

 • Quieres type-safety en compile-time       │
 • SQL verificado automáticamente             │
 • Performance sin sacrificar seguridad       │
 • Workflow basado en queries SQL             │
 • Un único driver de DB                      │

```

### 64.1.3 When to Use Raw SQL

**Ventajas del Raw SQL:**

1. **Control Total:** Escribes exactamente qué SQL se ejecuta
2. **Performance:** Sin overhead de traducción de ORM
3. **Queries Complejas:** CTEs, window functions, lateral joins
4. **Database-Specific Features:** JSON operators, arrays, tipos personalizados
5. **Debugging:** El SQL es visible y fácil de optimizar
6. **Testing:** Queries aisladas son más fáciles de testear

**Desventajas del Raw SQL:**

1. **Mantenimiento:** Cambios en schema requieren actualizar queries
2. **Boilerplate:** Más código para scanning y mapping
3. **Error prone:** Riesgo de inyección SQL si no usas placeholders
4. **Portabilidad:** SQL específico del motor no es portable

**Ejemplo: Cuándo Raw SQL es imprescindible**

```go
// ❌ GORM no puede expresar esto fácilmente
// Necesitas CTEs, window functions y JSON aggregation simultáneamente

// ✅ Con raw SQL:
const complexQuery = `
WITH user_stats AS (
    SELECT 
        u.id,
        u.name,
        COUNT(DISTINCT o.id) as total_orders,
        SUM(o.total) FILTER (WHERE o.created_at > NOW() - INTERVAL '30 days') 
            as recent_spend,
        ROW_NUMBER() OVER (ORDER BY SUM(o.total) DESC) as rank
    FROM users u
    LEFT JOIN orders o ON o.user_id = u.id
    GROUP BY u.id, u.name
)
SELECT 
    s.id,
    s.name,
    s.total_orders,
    s.recent_spend,
    s.rank,
    json_agg(
        json_build_object('month', extract(month from o.created_at), 'total', o.total)
    ) as monthly_sales
FROM user_stats s
LEFT JOIN orders o ON o.user_id = s.id
WHERE s.rank <= 100
GROUP BY s.id, s.name, s.total_orders, s.recent_spend, s.rank
ORDER BY s.rank
`

rows, err := db.QueryContext(ctx, complexQuery)
// ... scanning lógica ...
```

### 64.1.4 Drivers Ecosystem

**Drivers Oficialmente Soportados:**

#### PostgreSQL

```go
// Driver pq (github.com/lib/pq)
import _ "github.com/lib/pq"

db, err := sql.Open("postgres", 
    "user=postgres password=pass dbname=mydb sslmode=disable")

// Driver pgx (github.com/jackc/pgx/v5)
import _ "github.com/jackc/pgx/v5/stdlib"

db, err := sql.Open("pgx", 
    "postgres://user:pass@localhost/mydb")
```

**Características de cada driver:**

| Feature | pq | pgx |
|---------|-----|-----|
| Velocidad | Buena | Excelente |
| Features | Completo | Más features |
| Streaming | Sí | Sí (mejor) |
| COPY support | Sí | Sí (mejor) |
| Custom types | Limitado | Excelente |
| Maturity | Muy maduro | Moderno, rápido |

#### MySQL

```go
// Driver github.com/go-sql-driver/mysql
import _ "github.com/go-sql-driver/mysql"

db, err := sql.Open("mysql", 
    "user:password@tcp(localhost:3306)/mydb?parseTime=true")
```

#### SQLite

```go
// Driver modernc.org/sqlite (puro Go)
import _ "modernc.org/sqlite"

db, err := sql.Open("sqlite", "file:mydb.db?cache=shared")

// o CGO-based (más rápido pero requiere compilador C)
import _ "github.com/mattn/go-sqlite3"

db, err := sql.Open("sqlite3", "file:mydb.db?cache=shared")
```

**Matriz de Decisión para Seleccionar Driver:**

```

 PostgreSQL                                                  │

 pq:  Estable, maduro, buen balance                         │
 pgx: Más fast, más features, recomendado para nuevo code   │



 MySQL                                                       │

 go-sql-driver/mysql: Estándar de facto, bien mantenido     │



 SQLite                                                      │

 modernc.org/sqlite: Go puro, portable                      │
 mattn/go-sqlite3:   Más rápido (requiere CGO)              │

```

---

## 64.2 - Database Drivers

### 64.2.1 PostgreSQL - lib/pq

**lib/pq** es el driver más antiguo y maduro para PostgreSQL en Go.

```go
package main

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq"
)

func main() {
    // Connection string completo
    connStr := "postgres://user:password@localhost:5432/mydb?sslmode=disable"
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Verificar conexión
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
}
```

**Características de lib/pq:**

- Connection string style: `postgres://user:pass@host:port/db`
- Soporta el protocolo v3 de PostgreSQL
- PlaceHolders: `$1, $2, ...` (numerados)
- Tipos especiales: ARRAY, JSON, JSONB, hstore

```go
// Trabajar con arrays PostgreSQL
var tags []string
err := db.QueryRow(
    "SELECT tags FROM articles WHERE id = $1", 
    articleID).Scan(pq.Array(&tags))

// Trabajar con JSON
var metadata map[string]interface{}
err := db.QueryRow(
    "SELECT metadata FROM users WHERE id = $1", 
    userID).Scan(pq.JSONArray(&metadata))

// UNNEST para expandir arrays
rows, err := db.Query(
    "SELECT unnest($1::text[]) as tag", 
    pq.Array([]string{"go", "rust", "python"}))
```

### 64.2.2 PostgreSQL - pgx

**pgx** es un driver más moderno que ofrece mejor performance y más features.

```go
import (
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/stdlib"
)

// Para usar con database/sql
db, err := sql.Open("pgx", "postgres://user:pass@localhost/db")

// O usar pgx directamente (sin database/sql)
conn, err := pgx.Connect(ctx, "postgres://user:pass@localhost/db")
defer conn.Close(ctx)

// Ventajas de usar pgx directamente:
// - Mejor performance
// - Acceso a tipos nativos de PostgreSQL
// - Streaming de resultados
// - Mejor error handling
```

**Comparativa de Performance:**

```
Benchmark: 1000 queries SELECT * FROM usuarios

lib/pq:     234ms
pgx stdlib: 189ms
pgx native: 145ms

lib/pq:     100%
pgx stdlib: 81%
pgx native: 62%
```

**Usar pgx directamente (sin database/sql):**

```go
package main

import (
    "context"
    "github.com/jackc/pgx/v5"
)

type User struct {
    ID   int
    Name string
    Email string
}

func main() {
    conn, err := pgx.Connect(context.Background(), 
        "postgres://user:pass@localhost/mydb")
    if err != nil {
        panic(err)
    }
    defer conn.Close(context.Background())
    
    // QueryRow
    var u User
    err = conn.QueryRow(context.Background(),
        "SELECT id, name, email FROM users WHERE id = $1", 42).
        Scan(&u.ID, &u.Name, &u.Email)
    
    // Query con rows.Next()
    rows, err := conn.Query(context.Background(),
        "SELECT id, name FROM users")
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            panic(err)
        }
    }
    
    // Batch operations (muy rpido)
    batch := &pgx.Batch{}
    batch.Queue("INSERT INTO users(name) VALUES($1)", "Alice")
    batch.Queue("INSERT INTO users(name) VALUES($1)", "Bob")
    
    results := conn.SendBatch(context.Background(), batch)
    defer results.Close()
}
```

**Features únicos de pgx:**

```go
// Copy for bulk inserts (mucho más rápido)
rows := [][]interface{}{
    {1, "Alice", "alice@example.com"},
    {2, "Bob", "bob@example.com"},
}

_, err := conn.CopyFromRows(ctx, pgx.Identifier{"public", "users"},
    []string{"id", "name", "email"}, 
    pgx.CopyFromRows(rows))

// Prepared statements con mejor caching
stmt, err := conn.Prepare(ctx, "select_user", 
    "SELECT * FROM users WHERE id = $1")
row := conn.QueryRow(ctx, stmt, userID)

// Listen/Notify (pub/sub nativo)
conn, err := pgx.Connect(ctx, connString)
err = conn.Listen(ctx, "channel_name")

// En otra goroutine
notification, err := conn.WaitForNotification(ctx)
log.Println(notification.Payload) // Datos publicados
```

### 64.2.3 MySQL - go-sql-driver

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// Connection string
// [username[:password]@][protocol[(address)]]/dbname[?param=value]
db, err := sql.Open("mysql", 
    "root:password@tcp(localhost:3306)/mydb?parseTime=true")
```

**Parámetros importantes:**

```
parseTime=true           // Convertir DATE/DATETIME a time.Time
loc=Local                // Zona horaria para parsing
allowAllFiles=true       // Permitir LOAD DATA LOCAL
clientFoundRows=true     // affected rows vs found rows
multiStatements=true     // Permitir múltiples statements (cuidado con inyección)
```

```go
// Ejemplo completo con MySQL
package main

import (
    "database/sql"
    "log"
    "time"
    _ "github.com/go-sql-driver/mysql"
)

type Article struct {
    ID        int
    Title     string
    CreatedAt time.Time
}

func main() {
    db, err := sql.Open("mysql", 
        "user:pass@tcp(localhost:3306)/mydb?parseTime=true")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    var a Article
    row := db.QueryRow(
        "SELECT id, title, created_at FROM articles WHERE id = ?", 1)
    
    // Nota: MySQL usa ? para placeholders, no $1
    if err := row.Scan(&a.ID, &a.Title, &a.CreatedAt); err != nil {
        log.Fatal(err)
    }
    
    log.Printf("%+v\n", a)
}
```

### 64.2.4 SQLite

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"  // Puro Go, multiplataforma
)

// En memoria (para tests)
db, err := sql.Open("sqlite", ":memory:")

// Archivo
db, err := sql.Open("sqlite", "file:mydb.db?cache=shared&mode=rwc")
```

**Características de SQLite:**

```go
// SQLite es muy útil para:
// 1. Tests locales
// 2. Apps de una sola máquina
// 3. Embedded databases
// 4. Prototipos rápidos

// Pragma for optimization
_, err := db.Exec("PRAGMA journal_mode=WAL")      // Write-ahead logging
_, err := db.Exec("PRAGMA foreign_keys=ON")       // Enforce FK constraints
_, err := db.Exec("PRAGMA synchronous=NORMAL")    // Balance seguridad/velocidad
_, err := db.Exec("PRAGMA cache_size=10000")      // Buffer cache
```

**Cuidado: SQLite en modo concurrente**

```go
// SQLite tiene limitaciones de concurrencia
// Por defecto: solo 1 escritor a la vez

// Para mejor concurrencia, usa WAL mode:
_, err := db.Exec("PRAGMA journal_mode=WAL")

// Luego puedes tener múltiples writers simultáneamente
// pero con overhead adicional

// Mejor práctica: usa un pool de conexiones pequeño
db.SetMaxOpenConns(25)  // SQLite puede usar más connections con WAL
```

### 64.2.5 Driver Selection Criteria

**Matriz de Decisión:**

```

 Criterios para seleccionar un driver                         │

 1. COMPATIBILIDAD CON DATABASE/SQL                           │
    ✓ Todos los drivers principales la soportan              │
    ✗ Algunos drivers modernos tienen API nativa mejorada    │
                                                              │
 2. PERFORMANCE                                               │
    PostgreSQL: pgx (nativo) > pgx (stdlib) > lib/pq        │
    MySQL: go-sql-driver es decente                          │
    SQLite: mattn/go-sqlite3 > modernc.org/sqlite           │
                                                              │
 3. CARACTERÍSTICAS ESPECIALES                                │
    PostgreSQL: ARRAY, JSON, JSONB, hstore                   │
    MySQL: JSON (nivel básico)                               │
    SQLite: JSON1 extension (bsico)                         │
                                                              │
 4. PORTABILIDAD                                              │
    modernc.org/sqlite: Máxima (puro Go)                    │
    lib/pq, go-sql-driver: Alto (estándares ISO)             │
    pgx: Medio (PostgreSQL-specific pero bien documentado)   │
                                                              │
 5. MANTENIMIENTO                                             │
    lib/pq: Muy estable, cambios lentos                      │
    pgx: Activamente desarrollado, moderno                   │
    go-sql-driver: Bien mantenido, estable                   │
    modernc.org/sqlite: Activamente mantenido                │
                                                              │
 6. COMUNIDAD Y DOCUMENTACIÓN                                │
    lib/pq: Excelente, muchos recursos                       │
    go-sql-driver: Buena documentación                       │
    pgx: Documentación completa y ejemplos                   │

```

### 64.2.6 Custom Driver Development

Crear un driver personalizado permite soportar nuevas bases de datos.

```go
// Un driver debe implementar driver.Driver
package mydriver

import (
    "database/sql"
    "database/sql/driver"
)

type Driver struct{}

// Open abre una conexión a la base de datos
func (d *Driver) Open(name string) (driver.Conn, error) {
    return &Connection{
        // ... internals
    }, nil
}

// Implementar interface driver.Conn
type Connection struct{}

func (c *Connection) Prepare(query string) (driver.Stmt, error) {
    return &Statement{query: query}, nil
}

func (c *Connection) Close() error {
    return nil
}

func (c *Connection) Begin() (driver.Tx, error) {
    return &Transaction{}, nil
}

// Implementar driver.Stmt
type Statement struct {
    query string
}

func (s *Statement) Close() error { return nil }
func (s *Statement) NumInput() int { return -1 }

func (s *Statement) Exec(args []driver.Value) (driver.Result, error) {
    // Ejecutar query
    return &Result{}, nil
}

func (s *Statement) Query(args []driver.Value) (driver.Rows, error) {
    // Ejecutar query y retornar rows
    return &Rows{}, nil
}

// Registrar el driver
func init() {
    sql.Register("mydb", &Driver{})
}
```

**Ejemplo simplificado de driver en memoria:**

```go
package memdb

import (
    "database/sql"
    "database/sql/driver"
    "fmt"
    "sync"
)

type MemDriver struct {
    data map[string][]map[string]interface{}
    mu   sync.RWMutex
}

func (d *MemDriver) Open(name string) (driver.Conn, error) {
    return &MemConn{
        driver: d,
    }, nil
}

type MemConn struct {
    driver *MemDriver
}

func (c *MemConn) Prepare(query string) (driver.Stmt, error) {
    return &MemStmt{
        conn:  c,
        query: query,
    }, nil
}

func (c *MemConn) Close() error {
    return nil
}

func (c *MemConn) Begin() (driver.Tx, error) {
    return &MemTx{conn: c}, nil
}

// Más implementaciones...

// Usar el driver
func init() {
    sql.Register("memdb", &MemDriver{
        data: make(map[string][]map[string]interface{}),
    })
}

// En main:
db, _ := sql.Open("memdb", "")
```

---

## 64.3 - Connection Management

### 64.3.1 sql.Open & sql.DB

```go
// sql.Open NO abre una conexión inmediatamente
// Solo prepara el pool de conexiones
db, err := sql.Open("postgres", connStr)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Verificar que la conexión es válida
if err := db.Ping(); err != nil {
    log.Fatal("Database unavailable:", err)
}

// O con context y timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := db.PingContext(ctx); err != nil {
    log.Fatal("Database unavailable:", err)
}
```

**Estructura interna de sql.DB:**

```

         sql.DB (singleton)           │

 - Connection Pool                    │
 - MaxOpenConns: límite total         │
 - MaxIdleConns: conexiones idle      │
 - ConnMaxLifetime: edad máxima       │
 - ConnMaxIdleTime: idle máximo       │
 - Stats: métricas de pool            │

         │
         ├─ Conn 1 (idle)
         ├─ Conn 2 (idle)
         ├─ Conn 3 (in use)
         ├─ Conn 4 (in use)
         └─ Conn 5 (idle)
```

### 64.3.2 Connection Pooling

**Configuración del Pool:**

```go
db, err := sql.Open("postgres", connStr)
if err != nil {
    log.Fatal(err)
}

// CRÍTICO: Configurar el pool correctamente
db.SetMaxOpenConns(25)      // Máx conexiones abiertas
db.SetMaxIdleConns(5)       // Conexiones idle para reutilizar
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(10 * time.Minute)
```

**Cómo funciona el pool:**

```
SOLICITUD DE CONEXIÓN:
    │
    ├─ ¿Hay conexión idle disponible?
    │  └─ SÍ → Reutilizar
    │
    ├─ ¿Estamos bajo MaxOpenConns?
    │  └─ SÍ → Crear nueva conexión
    │
    └─ NO → Esperar (bloqueante)

LIBERACIÓN DE CONEXIÓN:
    │
    ├─ ¿Est bajo MaxIdleConns?
    │  └─ SÍ → Agregar al pool de idle
    │
    └─ NO → Cerrar la conexión
```

### 64.3.3 MaxOpenConns, MaxIdleConns

```go
// CONFIGURACIÓN POR CASO DE USO

// High-traffic web server
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)

// Batch processing
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(5)

// Background workers
db.SetMaxOpenConns(10)
db.SetMaxIdleConns(2)

// CLI tool
db.SetMaxOpenConns(5)
db.SetMaxIdleConns(1)
```

**Fórmula recomendada:**

```
MaxOpenConns = CPU_CORES * 4 + Extras para I/O bloqueante

Ejemplo: 8 cores
    MaxOpenConns = 8 * 4 + 4 = 36

MaxIdleConns = MaxOpenConns / 5 (como guía)
    MaxIdleConns = 36 / 5 = 7 (usar 7-10)
```

### 64.3.4 Connection Timeout

```go
// Timeout en el nivel de base de datos
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Query con timeout
row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

// Exec con timeout
_, err := db.ExecContext(ctx, 
    "INSERT INTO logs(message) VALUES($1)", msg)

// Timeout para toda una transacción
tx, err := db.BeginTx(ctx, &sql.TxOptions{})
```

### 64.3.5 Health Checks

```go
// Health check simple
func healthCheck(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    return db.PingContext(ctx)
}

// Health check con retry
func healthCheckWithRetry(db *sql.DB, maxRetries int) error {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        if err := db.PingContext(ctx); err == nil {
            cancel()
            return nil
        } else {
            lastErr = err
        }
        cancel()
        time.Sleep(time.Second)
    }
    return lastErr
}

// Health check periodico
func startHealthCheck(db *sql.DB, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        
        for range ticker.C {
            if err := healthCheck(db); err != nil {
                log.Printf("Health check failed: %v", err)
                // Alertar, escalar, etc.
            }
        }
    }()
}

// Stats del pool (debugging)
func printPoolStats(db *sql.DB) {
    stats := db.Stats()
    log.Printf("DB Stats: Open=%d, InUse=%d, Idle=%d, MaxOpen=%d",
        stats.OpenConnections,
        stats.InUse,
        stats.Idle,
        stats.MaxOpenConnections,
    )
}
```

---

## 64.4 - Basic Operations

### 64.4.1 Query Execution

**QueryRow (una fila):**

```go
var name string
err := db.QueryRow("SELECT name FROM users WHERE id = $1", id).Scan(&name)
if err == sql.ErrNoRows {
    return fmt.Errorf("user %d not found", id)
}
if err != nil {
    return err
}
```

**Query (múltiples filas):**

```go
rows, err := db.Query("SELECT id, name, email FROM users")
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var id int
    var name, email string
    
    if err := rows.Scan(&id, &name, &email); err != nil {
        return err
    }
    
    fmt.Printf("User: %d, %s, %s\n", id, name, email)
}

// Importante: Check error después del loop
if err = rows.Err(); err != nil {
    return err
}
```

**Con Context:**

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

row := db.QueryRowContext(ctx, 
    "SELECT name FROM users WHERE id = $1", id)
var name string
if err := row.Scan(&name); err != nil {
    return err
}
```

### 64.4.2 Scanning Results

```go
// Scanning básico
var id int
var name string
var email sql.NullString  // Para nullable columns

err := row.Scan(&id, &name, &email)

// Verificar null
if email.Valid {
    fmt.Println(email.String)
} else {
    fmt.Println("Email is NULL")
}
```

**Tipos nullables:**

```go
// sql.Null* para valores que pueden ser NULL

var (
    nullInt    sql.NullInt64
    nullString sql.NullString
    nullBool   sql.NullBool
    nullTime   sql.NullTime
)

row.Scan(&nullInt, &nullString, &nullBool, &nullTime)

if nullInt.Valid {
    fmt.Println(nullInt.Int64)
}
```

**Scanning en struct:**

```go
type User struct {
    ID    int
    Name  string
    Email sql.NullString
}

var u User
err := row.Scan(&u.ID, &u.Name, &u.Email)

// O crear helper function
func scanUser(row *sql.Row) (*User, error) {
    u := &User{}
    err := row.Scan(&u.ID, &u.Name, &u.Email)
    return u, err
}
```

### 64.4.3 Prepared Statements

**Crear y reutilizar:**

```go
stmt, err := db.Prepare("SELECT name FROM users WHERE id = $1")
if err != nil {
    return err
}
defer stmt.Close()

// Reutilizar statement múltiples veces
for _, id := range []int{1, 2, 3, 4, 5} {
    var name string
    err := stmt.QueryRow(id).Scan(&name)
    if err != nil {
        log.Printf("Error for id %d: %v", id, err)
    } else {
        fmt.Println(name)
    }
}
```

**Con Context:**

```go
stmt, err := db.PrepareContext(ctx, 
    "INSERT INTO users(name, email) VALUES($1, $2)")
if err != nil {
    return err
}
defer stmt.Close()

result, err := stmt.ExecContext(ctx, "Alice", "alice@example.com")
if err != nil {
    return err
}

lastID, err := result.LastInsertId()
```

**Performance:**

```
1000 queries sin prepared statement:
    - Parse query: 1000 veces
    - Plan: 1000 veces
    - Execute: 1000 veces
    Total: ~150ms

1000 queries con prepared statement:
    - Parse query: 1 vez
    - Plan: 1 vez
    - Execute: 1000 veces
    Total: ~50ms

Mejora: 3x más rápido
```

### 64.4.4 Exec for Modifications

```go
// INSERT
result, err := db.Exec(
    "INSERT INTO users(name, email) VALUES($1, $2)",
    "Bob", "bob@example.com")
if err != nil {
    return err
}

// LastInsertId (no funciona en PostgreSQL sin RETURNING)
id, err := result.LastInsertId()

// RowsAffected
rows, err := result.RowsAffected()
fmt.Printf("Inserted %d rows\n", rows)

// UPDATE
result, err := db.Exec(
    "UPDATE users SET email = $1 WHERE id = $2",
    "newemail@example.com", 42)

affected, _ := result.RowsAffected()
fmt.Printf("Updated %d rows\n", affected)

// DELETE
result, err := db.Exec("DELETE FROM users WHERE id = $1", 42)
```

**Nota sobre LastInsertId en PostgreSQL:**

```go
// ❌ LastInsertId no funciona directamente en PostgreSQL
id, err := result.LastInsertId()  // Retorna error

// ✅ Usar RETURNING
var id int
err := db.QueryRow(
    "INSERT INTO users(name) VALUES($1) RETURNING id",
    "Alice").Scan(&id)
```

---

## 64.5 - Advanced SQL Patterns

### 64.5.1 Complex Queries - Joins

```go
// INNER JOIN
const query = `
SELECT 
    u.id, u.name,
    COUNT(DISTINCT o.id) as order_count,
    SUM(o.total) as total_spent
FROM users u
INNER JOIN orders o ON o.user_id = u.id
GROUP BY u.id, u.name
HAVING SUM(o.total) > $1
ORDER BY total_spent DESC
`

rows, err := db.Query(query, 1000.0)
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    var orderCount int
    var totalSpent float64
    
    err := rows.Scan(&id, &name, &orderCount, &totalSpent)
    if err != nil {
        log.Println(err)
        continue
    }
    
    fmt.Printf("%s has %d orders totaling $%.2f\n",
        name, orderCount, totalSpent)
}
```

**LEFT JOIN con NULL handling:**

```go
const query = `
SELECT 
    u.id,
    u.name,
    COUNT(DISTINCT o.id)::INTEGER as order_count
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
GROUP BY u.id, u.name
ORDER BY u.name
`

type UserStats struct {
    ID         int
    Name       string
    OrderCount int  // Puede ser 0 (no NULL en este caso)
}

rows, err := db.Query(query)
defer rows.Close()

var users []UserStats
for rows.Next() {
    var u UserStats
    if err := rows.Scan(&u.ID, &u.Name, &u.OrderCount); err != nil {
        return err
    }
    users = append(users, u)
}

return rows.Err()
```

### 64.5.2 CTEs (WITH clauses)

```go
// Common Table Expressions para queries complejas
const query = `
WITH ranked_users AS (
    SELECT 
        id,
        name,
        email,
        created_at,
        ROW_NUMBER() OVER (ORDER BY created_at DESC) as rank
    FROM users
    WHERE active = true
),
top_spenders AS (
    SELECT 
        u.id,
        u.name,
        SUM(o.total) as total_spent
    FROM users u
    JOIN orders o ON o.user_id = u.id
    GROUP BY u.id, u.name
    ORDER BY total_spent DESC
    LIMIT 100
)
SELECT 
    ru.id,
    ru.name,
    ru.email,
    COALESCE(ts.total_spent, 0) as total_spent,
    ru.rank
FROM ranked_users ru
LEFT JOIN top_spenders ts ON ts.id = ru.id
ORDER BY ru.rank
LIMIT 50
`

type UserRanking struct {
    ID          int
    Name        string
    Email       string
    TotalSpent  float64
    Rank        int
}

rows, err := db.Query(query)
defer rows.Close()

var users []UserRanking
for rows.Next() {
    var u UserRanking
    if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.TotalSpent, &u.Rank); err != nil {
        return err
    }
    users = append(users, u)
}

return rows.Err()
```

### 64.5.3 Window Functions

```go
// Window functions para análisis avanzado
const query = `
SELECT 
    order_id,
    customer_id,
    amount,
    order_date,
    -- Running total
    SUM(amount) OVER (
        PARTITION BY customer_id 
        ORDER BY order_date 
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as running_total,
    
    -- Rank within customer
    RANK() OVER (
        PARTITION BY customer_id 
        ORDER BY amount DESC
    ) as amount_rank,
    
    -- Year-over-year comparison
    LAG(amount) OVER (
        PARTITION BY customer_id, DATE_TRUNC('month', order_date)::DATE
        ORDER BY order_date
    ) as previous_month_amount,
    
    -- Percentile
    PERCENT_RANK() OVER (
        ORDER BY amount
    ) as amount_percentile
FROM orders
WHERE order_date >= CURRENT_DATE - INTERVAL '1 year'
ORDER BY customer_id, order_date
`

type OrderAnalysis struct {
    OrderID        int
    CustomerID     int
    Amount         float64
    OrderDate      time.Time
    RunningTotal   float64
    AmountRank     int
    PrevMonthAmount sql.NullFloat64
    AmountPercentile float64
}

rows, err := db.Query(query)
defer rows.Close()

for rows.Next() {
    var o OrderAnalysis
    err := rows.Scan(
        &o.OrderID, &o.CustomerID, &o.Amount, &o.OrderDate,
        &o.RunningTotal, &o.AmountRank, &o.PrevMonthAmount, &o.AmountPercentile,
    )
    if err != nil {
        return err
    }
}
```

### 64.5.4 JSON Operators (PostgreSQL)

```go
// PostgreSQL JSON/JSONB operations
type Document struct {
    ID    int
    Title string
    Data  json.RawMessage  // Raw JSON
    Tags  []string
}

// JSON extract
const queryJSON = `
SELECT 
    id,
    title,
    data,
    data->'tags' as tags_array
FROM documents
WHERE data->>'type' = $1
`

// Scanning JSON
rows, err := db.Query(queryJSON, "article")
defer rows.Close()

for rows.Next() {
    var d Document
    var tagsJSON []byte
    
    err := rows.Scan(&d.ID, &d.Title, &d.Data, &tagsJSON)
    if err != nil {
        return err
    }
    
    // Unmarshaling JSON
    json.Unmarshal(tagsJSON, &d.Tags)
}

// JSON aggregation
const aggregateJSON = `
SELECT 
    category,
    jsonb_agg(jsonb_build_object(
        'id', id,
        'title', title,
        'price', price
    )) as items
FROM products
GROUP BY category
`

type CategoryProducts struct {
    Category string
    Items    json.RawMessage
}

rows, err := db.Query(aggregateJSON)
defer rows.Close()

for rows.Next() {
    var cp CategoryProducts
    if err := rows.Scan(&cp.Category, &cp.Items); err != nil {
        return err
    }
}
```

### 64.5.5 Array Types (PostgreSQL)

```go
import "github.com/lib/pq"

// Insertar arrays
tags := []string{"go", "database", "sql"}
_, err := db.Exec(
    "INSERT INTO articles(title, tags) VALUES($1, $2)",
    "My Article",
    pq.Array(tags),
)

// Querear arrays
var tags []string
err := db.QueryRow(
    "SELECT tags FROM articles WHERE id = $1",
    articleID,
).Scan(pq.Array(&tags))

// Array operations en SQL
const query = `
SELECT 
    id,
    title,
    tags,
    array_length(tags, 1) as tag_count,
    'go' = ANY(tags) as is_go_article
FROM articles
WHERE 'database' = ANY(tags)
ORDER BY array_length(tags, 1) DESC
`

type Article struct {
    ID          int
    Title       string
    Tags        pq.StringArray
    TagCount    sql.NullInt64
    IsGoArticle bool
}

rows, err := db.Query(query)
defer rows.Close()

for rows.Next() {
    var a Article
    if err := rows.Scan(&a.ID, &a.Title, &a.Tags, &a.TagCount, &a.IsGoArticle); err != nil {
        return err
    }
}
```

---

## 64.6 - Error Handling

### 64.6.1 sql.ErrNoRows

```go
// Distinguir entre "no encontrado" y error real
var name string
err := db.QueryRow(
    "SELECT name FROM users WHERE id = $1", userID,
).Scan(&name)

if err == sql.ErrNoRows {
    // Usuario no existe
    return fmt.Errorf("user %d not found", userID)
}
if err != nil {
    // Error real (conexión, etc)
    return fmt.Errorf("database error: %w", err)
}

// Usar nombre
fmt.Println(name)
```

### 64.6.2 Error Types

```go
// PostgreSQL specific errors
import "github.com/lib/pq"

result, err := db.Exec(
    "INSERT INTO users(email) VALUES($1) CONSTRAINT users_email_unique",
    email,
)

if err != nil {
    if pgErr, ok := err.(*pq.Error); ok {
        switch pgErr.Code {
        case "23505":  // unique_violation
            return fmt.Errorf("email %s already exists", email)
        case "23503":  // foreign_key_violation
            return fmt.Errorf("referenced record not found")
        case "23502":  // not_null_violation
            return fmt.Errorf("required field is missing")
        default:
            return fmt.Errorf("database error: %v", pgErr)
        }
    }
    return err
}

// MySQL error handling
import "github.com/go-sql-driver/mysql"

result, err := db.Exec("INSERT INTO users(email) VALUES(?)", email)
if err != nil {
    if me, ok := err.(*mysql.MySQLError); ok {
        switch me.Number {
        case 1062:  // Duplicate entry
            return fmt.Errorf("email %s already exists", email)
        default:
            return fmt.Errorf("database error: %v", me)
        }
    }
    return err
}
```

### 64.6.3 Wrapping Errors

```go
import "fmt"

func GetUser(db *sql.DB, id int) (*User, error) {
    row := db.QueryRow("SELECT id, name FROM users WHERE id = $1", id)
    
    var u User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("getting user: %w", 
                fmt.Errorf("user not found"))
        }
        return nil, fmt.Errorf("getting user: %w", err)
    }
    
    return &u, nil
}

// Usage
u, err := GetUser(db, 1)
if err != nil {
    log.Printf("Error: %v", err)  // Error: getting user: user not found
    
    // Unwrap para checking
    if errors.Is(err, sql.ErrNoRows) {
        // ...
    }
}
```

### 64.6.4 Retry Logic

```go
import (
    "github.com/cenkalti/backoff/v4"
)

func queryWithRetry(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
    var rows *sql.Rows
    
    // Exponential backoff retry
    err := backoff.Retry(func() error {
        var err error
        rows, err = db.Query(query, args...)
        return err
    }, backoff.NewExponentialBackOff())
    
    return rows, err
}

// Manual retry con jitter
func queryWithManualRetry(db *sql.DB, maxRetries int, 
    query string, args ...interface{}) (*sql.Rows, error) {
    
    var rows *sql.Rows
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        var err error
        rows, err = db.Query(query, args...)
        
        if err == nil {
            return rows, nil
        }
        
        lastErr = err
        
        // Solo retry en ciertos errores
        if !isRetryableError(err) {
            return nil, err
        }
        
        // Exponential backoff con jitter
        backoffDuration := time.Duration(math.Pow(2, float64(i))) * time.Second
        jitter := time.Duration(rand.Intn(100)) * time.Millisecond
        time.Sleep(backoffDuration + jitter)
    }
    
    return nil, fmt.Errorf("retry exceeded: %w", lastErr)
}

func isRetryableError(err error) bool {
    // Retry en timeouts y connection errors
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Timeout() {
        return true
    }
    
    errStr := err.Error()
    retryableErrors := []string{
        "connection refused",
        "connection reset",
        "i/o timeout",
        "EOF",
    }
    
    for _, pattern := range retryableErrors {
        if strings.Contains(strings.ToLower(errStr), pattern) {
            return true
        }
    }
    
    return false
}
```

### 64.6.5 Timeout Handling

```go
// Timeout a nivel de query
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

row := db.QueryRowContext(ctx, 
    "SELECT id FROM users WHERE id = $1", userID)

var id int
if err := row.Scan(&id); err != nil {
    if err == context.DeadlineExceeded {
        log.Println("Query timeout exceeded")
    } else if err == sql.ErrNoRows {
        log.Println("User not found")
    } else {
        log.Printf("Error: %v", err)
    }
    return err
}

// Timeout a nivel de conexión
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

// Todos los statements en la transacción hereda el timeout
```

---

## 64.7 - Transactions

### 64.7.1 BeginTx Pattern

```go
func TransferMoney(ctx context.Context, db *sql.DB, 
    fromUserID int, toUserID int, amount float64) error {
    
    tx, err := db.BeginTx(ctx, &sql.TxOptions{})
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // Debit from source account
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance - $1 WHERE user_id = $2",
        amount, fromUserID)
    if err != nil {
        return fmt.Errorf("debit: %w", err)
    }
    
    // Credit to destination account
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance + $1 WHERE user_id = $2",
        amount, toUserID)
    if err != nil {
        return fmt.Errorf("credit: %w", err)
    }
    
    // Record transaction
    _, err = tx.ExecContext(ctx,
        "INSERT INTO transaction_log(from_user, to_user, amount) VALUES($1, $2, $3)",
        fromUserID, toUserID, amount)
    if err != nil {
        return fmt.Errorf("record transaction: %w", err)
    }
    
    // Commit
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("commit: %w", err)
    }
    
    return nil
}
```

### 64.7.2 Commit/Rollback

```go
// Simple transaction
tx, err := db.Begin()
if err != nil {
    return err
}

// Ensure cleanup
rollback := true
defer func() {
    if rollback {
        tx.Rollback()
    }
}()

// Operaciones
_, err = tx.Exec("INSERT INTO users(name) VALUES($1)", "Alice")
if err != nil {
    return err
}

// Commit
if err = tx.Commit(); err != nil {
    return err
}
rollback = false  // Previene rollback en defer
```

### 64.7.3 Isolation Levels

```go
// PostgreSQL isolation levels
import "database/sql"

// READ UNCOMMITTED (debería ser READ COMMITTED en PostrgreSQL)
tx, _ := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelReadUncommitted,
})

// READ COMMITTED (default en PostgreSQL)
tx, _ := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelReadCommitted,
})

// REPEATABLE READ
tx, _ := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelRepeatableRead,
})

// SERIALIZABLE (máximo nivel, puede tener overhead)
tx, _ := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable,
})

// Read-only transaction (optimización)
tx, _ := db.BeginTx(ctx, &sql.TxOptions{
    ReadOnly: true,
})
```

```

 Isolation Levels y sus Problemas                             │

                                                              │
 READ UNCOMMITTED:                                            │
 • Dirty reads: SÍ                                            │
 • Non-repeatable reads: SÍ                                  │
 • Phantom reads: SÍ                                         │
 • Performance: Excelente                                    │
                                                              │
 READ COMMITTED (PostgreSQL default):                         │
 Dirty reads: NO                                           │ 
 • Non-repeatable reads: SÍ                                  │
 Phantom reads: SÍ                                         │ 
 • Performance: Muy bueno                                    │
                                                              │
 REPEATABLE READ:                                             │
 • Dirty reads: NO                                           │
 • Non-repeatable reads: NO                                  │
 • Phantom reads: SÍ (en teoría, PostgreSQL previene)        │
 • Performance: Bueno                                        │
                                                              │
 SERIALIZABLE:                                                
 • Dirty reads: NO                                           │
 • Non-repeatable NO reads:                                  
 • Phantom reads: NO                                         │
 • Performance: Puede lento ser                              
                                                              │

```

### 64.7.4 Savepoints

```go
// PostgreSQL savepoints
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// Crear savepoint
_, err = tx.Exec("SAVEPOINT sp1")
if err != nil {
    return err
}

// Operaciones
_, err = tx.Exec("INSERT INTO table1 VALUES(...)")
if err != nil {
    // Rollback al savepoint, no la transacción completa
    tx.Exec("ROLLBACK TO sp1")
    return err
}

// Más operaciones
_, err = tx.Exec("INSERT INTO table2 VALUES(...)")
if err != nil {
    return err
}

return tx.Commit().Error
```

### 64.7.5 Deadlock Handling

```go
func withDeadlockRetry(ctx context.Context, db *sql.DB, 
    fn func(*sql.Tx) error) error {
    
    const maxRetries = 3
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        tx, err := db.BeginTx(ctx, nil)
        if err != nil {
            return err
        }
        
        err = fn(tx)
        if err != nil {
            tx.Rollback()
            
            // Detectar deadlock (PostgreSQL error code 40P01)
            if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "40P01" {
                lastErr = err
                // Retry con backoff
                backoff := time.Duration(math.Pow(2, float64(i))) * 100 * time.Millisecond
                time.Sleep(backoff)
                continue
            }
            
            return err
        }
        
        if err = tx.Commit(); err != nil {
            // Si commit falló por deadlock, retry
            if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "40P01" {
                lastErr = err
                continue
            }
            return err
        }
        
        return nil
    }
    
    return fmt.Errorf("deadlock retry exceeded: %w", lastErr)
}

// Usage
err := withDeadlockRetry(ctx, db, func(tx *sql.Tx) error {
    // Transacción que podría causar deadlock
    _, err := tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", 100, 1)
    if err != nil {
        return err
    }
    _, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", 100, 2)
    return err
})
```

---

## 64.8 - Performance Optimization

### 64.8.1 Query Analysis (EXPLAIN)

```go
// EXPLAIN ANALYZE para entender query performance
query := `
EXPLAIN ANALYZE
SELECT u.id, u.name, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
GROUP BY u.id, u.name
ORDER BY order_count DESC
LIMIT 10
`

rows, err := db.Query(query)
defer rows.Close()

// Leer plan de ejecución
for rows.Next() {
    var line string
    if err := rows.Scan(&line); err != nil {
        log.Fatal(err)
    }
    fmt.Println(line)
}

// Ejemplo de output:
// Limit  (cost=45.23..45.24 rows=10 width=12) (actual time=2.123..2.145 rows=10 loops=1)
//   ->  Sort  (cost=45.23..45.24 rows=100 width=12) (actual time=2.120..2.125 rows=10 loops=1)
//        Sort Key: (count(o.id)) DESC
//        ->  GroupAggregate (cost=20.00..43.00 rows=100 width=12)
```

**Interpretación de EXPLAIN:**

```
cost=20.00..43.00  : Costo estimado (inicio..fin)
rows=100           : Filas estimadas
width=12           : Ancho promedio en bytes
actual time=2.1    : Tiempo real en ms
loops=1            : Cuántas veces se ejecutó
```

### 64.8.2 Index Strategies

```go
// Crear índices para optimizar queries comunes
_, err := db.Exec(`
CREATE INDEX idx_users_email ON users(email);
`)

// Índice compuesto para queries multi-column
_, err = db.Exec(`
CREATE INDEX idx_orders_user_date ON orders(user_id, created_at DESC);
`)

// Índice para búsqueda full-text (PostgreSQL)
_, err = db.Exec(`
CREATE INDEX idx_articles_title_tsvector ON articles 
USING GIN(to_tsvector('english', title));
`)

// Análisis sin crear índice (PostgreSQL)
err = db.QueryRow(`
SELECT * FROM pgstattuple('public.users')
`).Scan() // Ver estadísticas de tabla

// Verificar índices existentes
const checkIndexes = `
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'users'
`

rows, err := db.Query(checkIndexes)
defer rows.Close()

for rows.Next() {
    var name, def string
    if err := rows.Scan(&name, &def); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Index: %s\nDefinition: %s\n", name, def)
}
```

### 64.8.3 Connection Pooling Tuning

```go
// Benchmark para encontrar configuración óptima
package main

import (
    "database/sql"
    "fmt"
    "sync"
    "time"
)

func benchmarkPoolConfig(db *sql.DB, concurrency int) time.Duration {
    start := time.Now()
    var wg sync.WaitGroup
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                var x int
                db.QueryRow("SELECT 1").Scan(&x)
            }
        }()
    }
    
    wg.Wait()
    return time.Since(start)
}

func findOptimalPoolSize(db *sql.DB) {
    for maxConns := 5; maxConns <= 100; maxConns += 5 {
        db.SetMaxOpenConns(maxConns)
        duration := benchmarkPoolConfig(db, maxConns)
        fmt.Printf("MaxConns: %d, Duration: %v\n", maxConns, duration)
    }
}
```

### 64.8.4 Batch Operations

```go
// Batch inserts mucho más rápido
// ❌ Lento: 1000 queries individuales
for i := 0; i < 1000; i++ {
    db.Exec("INSERT INTO users(name) VALUES($1)", fmt.Sprintf("User %d", i))
}

// ✅ Rápido: Una sola query
const batchInsert = `
INSERT INTO users(name) VALUES
($1), ($2), ($3), ($4), ($5),
...
($1000)
`

// Con postgresql COPY (muy rápido)
import "github.com/lib/pq"

// Preparar datos
buffer := new(bytes.Buffer)
for i := 0; i < 1000; i++ {
    fmt.Fprintf(buffer, "%d\tUser %d\n", i, i)
}

// Copiar
_, err := db.Exec(pq.CopyIn("users", "id", "name"),
    // rows...
)

// O con pgx (aún más rápido)
rows := make([][]interface{}, 1000)
for i := 0; i < 1000; i++ {
    rows[i] = []interface{}{i, fmt.Sprintf("User %d", i)}
}

_, err := conn.CopyFromRows(ctx, pgx.Identifier{"users"}, 
    []string{"id", "name"}, pgx.CopyFromRows(rows))
```

### 64.8.5 Prepared Statements

```go
// Ya vimos esto en 64.4.3, pero repaso importante para performance

// ❌ Sin prepared statements (parse query cada vez)
for userID := 1; userID <= 1000; userID++ {
    db.QueryRow("SELECT name FROM users WHERE id = $1", userID)
}
// Tiempo: 150ms (1000 parses)

// ✅ Con prepared statements (parse 1 vez)
stmt, _ := db.Prepare("SELECT name FROM users WHERE id = $1")
defer stmt.Close()

for userID := 1; userID <= 1000; userID++ {
    stmt.QueryRow(userID)
}
// Tiempo: 50ms (1 parse + 1000 executes)

// Mejora: 3x más rápido
```

---

## 64.9 - Testing Raw SQL

### 64.9.1 Test Database Setup

```go
package data

import (
    "database/sql"
    "fmt"
    "testing"
    
    _ "github.com/lib/pq"
)

// Helper para crear DB de test
func setupTestDB(t *testing.T) *sql.DB {
    connStr := "postgres://test:test@localhost:5432/testdb"
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        t.Fatalf("Cannot connect to test database: %v", err)
    }
    
    // Verificar conexión
    if err := db.Ping(); err != nil {
        t.Fatalf("Test database unreachable: %v", err)
    }
    
    // Limpiar tables existentes
    cleanupSQL := `
    DROP TABLE IF EXISTS users CASCADE;
    DROP TABLE IF EXISTS orders CASCADE;
    `
    if _, err := db.Exec(cleanupSQL); err != nil {
        t.Fatalf("Failed to cleanup: %v", err)
    }
    
    // Crear schema
    schemaSQL := `
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );
    
    CREATE TABLE orders (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id),
        total DECIMAL(10,2) NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );
    `
    
    if _, err := db.Exec(schemaSQL); err != nil {
        t.Fatalf("Failed to create schema: %v", err)
    }
    
    t.Cleanup(func() {
        db.Close()
    })
    
    return db
}

// Test example
func TestGetUser(t *testing.T) {
    db := setupTestDB(t)
    
    // Insert test data
    _, err := db.Exec(
        "INSERT INTO users(name, email) VALUES($1, $2)",
        "Alice", "alice@example.com")
    if err != nil {
        t.Fatalf("Failed to insert test data: %v", err)
    }
    
    // Test the function
    var name string
    err = db.QueryRow("SELECT name FROM users WHERE email = $1",
        "alice@example.com").Scan(&name)
    
    if err != nil || name != "Alice" {
        t.Errorf("GetUser failed: got %s, want Alice", name)
    }
}
```

### 64.9.2 Testcontainers

```go
package data

import (
    "context"
    "database/sql"
    "fmt"
    "testing"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

// Container para PostgreSQL
func setupPostgresContainer(ctx context.Context, t *testing.T) *sql.DB {
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForListeningPort("5432/tcp"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, 
        testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
    if err != nil {
        t.Fatalf("Failed to create container: %v", err)
    }
    
    t.Cleanup(func() {
        container.Terminate(ctx)
    })
    
    // Get connection string
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "5432")
    
    connStr := fmt.Sprintf(
        "postgres://test:test@%s:%s/testdb?sslmode=disable",
        host, port.Port())
    
    db, _ := sql.Open("postgres", connStr)
    
    // Wait for DB to be ready
    for i := 0; i < 10; i++ {
        if err := db.Ping(); err == nil {
            return db
        }
    }
    
    t.Fatal("Database never became ready")
    return nil
}

func TestWithContainer(t *testing.T) {
    ctx := context.Background()
    db := setupPostgresContainer(ctx, t)
    defer db.Close()
    
    // Your tests here
}
```

### 64.9.3 Fixtures

```go
// fixture.go
package data

import (
    "database/sql"
    "testing"
)

type Fixture struct {
    DB *sql.DB
    t  *testing.T
}

func (f *Fixture) InsertUser(name, email string) int {
    var id int
    err := f.DB.QueryRow(
        "INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
        name, email).Scan(&id)
    
    if err != nil {
        f.t.Fatalf("Failed to insert user: %v", err)
    }
    
    return id
}

func (f *Fixture) InsertOrder(userID int, total float64) int {
    var id int
    err := f.DB.QueryRow(
        "INSERT INTO orders(user_id, total) VALUES($1, $2) RETURNING id",
        userID, total).Scan(&id)
    
    if err != nil {
        f.t.Fatalf("Failed to insert order: %v", err)
    }
    
    return id
}

func (f *Fixture) GetUserOrders(userID int) int {
    var count int
    err := f.DB.QueryRow(
        "SELECT COUNT(*) FROM orders WHERE user_id = $1", userID).Scan(&count)
    
    if err != nil {
        f.t.Fatalf("Failed to get order count: %v", err)
    }
    
    return count
}

// Usage in tests
func TestUserOrders(t *testing.T) {
    db := setupTestDB(t)
    fixture := &Fixture{DB: db, t: t}
    
    userID := fixture.InsertUser("Bob", "bob@example.com")
    fixture.InsertOrder(userID, 100.00)
    fixture.InsertOrder(userID, 200.00)
    
    if count := fixture.GetUserOrders(userID); count != 2 {
        t.Errorf("Expected 2 orders, got %d", count)
    }
}
```

### 64.9.4 Integration Testing

```go
package data

import (
    "context"
    "database/sql"
    "testing"
    "time"
)

// Integration test con transacción completa
func TestTransferMoney(t *testing.T) {
    db := setupTestDB(t)
    
    // Setup
    db.Exec("INSERT INTO users(name, email) VALUES($1, $2)", "Alice", "alice@example.com")
    db.Exec("INSERT INTO users(name, email) VALUES($1, $2)", "Bob", "bob@example.com")
    
    _, _ = db.Exec("ALTER TABLE users ADD COLUMN balance DECIMAL(10,2) DEFAULT 0")
    _, _ = db.Exec("UPDATE users SET balance = 1000 WHERE name = 'Alice'")
    _, _ = db.Exec("UPDATE users SET balance = 0 WHERE name = 'Bob'")
    
    // Test transaction
    ctx := context.Background()
    tx, _ := db.BeginTx(ctx, nil)
    
    tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1 WHERE name = $2", 100, "Alice")
    tx.ExecContext(ctx, "UPDATE users SET balance = balance + $1 WHERE name = $2", 100, "Bob")
    tx.Commit()
    
    // Verify
    var aliceBalance, bobBalance float64
    db.QueryRow("SELECT balance FROM users WHERE name = 'Alice'").Scan(&aliceBalance)
    db.QueryRow("SELECT balance FROM users WHERE name = 'Bob'").Scan(&bobBalance)
    
    if aliceBalance != 900 || bobBalance != 100 {
        t.Errorf("Transfer failed: Alice=%v, Bob=%v", aliceBalance, bobBalance)
    }
}

// Stress test con concurrencia
func TestConcurrentInserts(t *testing.T) {
    db := setupTestDB(t)
    
    numGoroutines := 10
    insertsPerGoroutine := 100
    
    done := make(chan error, numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            for j := 0; j < insertsPerGoroutine; j++ {
                _, err := db.Exec(
                    "INSERT INTO users(name, email) VALUES($1, $2)",
                    "User", "user@example.com")
                
                if err != nil {
                    done <- err
                    return
                }
            }
            done <- nil
        }(i)
    }
    
    for i := 0; i < numGoroutines; i++ {
        if err := <-done; err != nil {
            t.Fatalf("Insert error: %v", err)
        }
    }
    
    var count int
    db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    
    expected := numGoroutines * insertsPerGoroutine
    if count != expected {
        t.Errorf("Expected %d users, got %d", expected, count)
    }
}
```

### 64.9.5 Mocking SQL

```go
package data

import (
    "database/sql/driver"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
)

// Mock usando go-sqlmock
func TestQueryWithMock(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock: %v", err)
    }
    defer db.Close()
    
    // Expect a query
    rows := sqlmock.NewRows([]string{"id", "name", "email"}).
        AddRow(1, "Alice", "alice@example.com").
        AddRow(2, "Bob", "bob@example.com")
    
    mock.ExpectQuery("SELECT id, name, email FROM users").
        WillReturnRows(rows)
    
    // Execute query
    sqlRows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        t.Fatalf("Query failed: %v", err)
    }
    defer sqlRows.Close()
    
    // Verify
    users := []struct {
        ID    int
        Name  string
        Email string
    }{}
    
    for sqlRows.Next() {
        var u struct {
            ID    int
            Name  string
            Email string
        }
        sqlRows.Scan(&u.ID, &u.Name, &u.Email)
        users = append(users, u)
    }
    
    if len(users) != 2 {
        t.Errorf("Expected 2 users, got %d", len(users))
    }
    
    // Ensure expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

// Mock de transacción
func TestTransactionWithMock(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()
    
    // Expect BEGIN
    mock.ExpectBegin()
    
    // Expect INSERT
    mock.ExpectExec("INSERT INTO users").
        WithArgs("Alice", "alice@example.com").
        WillReturnResult(driver.NewResult(1, 1))
    
    // Expect COMMIT
    mock.ExpectCommit()
    
    // Run transaction
    tx, _ := db.Begin()
    tx.Exec("INSERT INTO users(name, email) VALUES(?, ?)", 
        "Alice", "alice@example.com")
    tx.Commit()
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Mock expectations failed: %v", err)
    }
}
```

---

## 64.10 - Scanning & Mapping

### 64.10.1 Row Scanning

```go
// Scanning básico
var id int
var name string
var email string

err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", 
    1).Scan(&id, &name, &email)

if err != nil {
    log.Fatal(err)
}

fmt.Printf("User: %d, %s, %s\n", id, name, email)
```

### 64.10.2 Type Conversion

```go
// Conversiones automáticas
var (
    id       int       // INTEGER -> int
    name     string    // VARCHAR -> string
    balance  float64   // DECIMAL -> float64
    created  time.Time // TIMESTAMP -> time.Time
    active   bool      // BOOLEAN -> bool
)

err := db.QueryRow(`
SELECT id, name, balance, created_at, active 
FROM users WHERE id = $1`, 1).Scan(&id, &name, &balance, &created, &active)

// Conversiones con casting
const castQuery = `
SELECT 
    CAST(id AS TEXT) as id_text,
    CAST(balance AS INTEGER) as balance_int,
    EXTRACT(YEAR FROM created_at) as year_created
FROM users WHERE id = $1
`

var (
    idText      string  // CAST(id AS TEXT)
    balanceInt  int     // CAST(balance AS INTEGER)
    yearCreated int     // EXTRACT result
)

db.QueryRow(castQuery, 1).Scan(&idText, &balanceInt, &yearCreated)
```

### 64.10.3 Nullable Types

```go
// sql.Null* tipos para NULL-handling
var (
    nullInt     sql.NullInt64
    nullString  sql.NullString
    nullBool    sql.NullBool
    nullTime    sql.NullTime
    nullFloat   sql.NullFloat64
)

err := db.QueryRow(`
SELECT 
    id,
    phone,
    is_premium,
    last_login,
    discount_rate
FROM users WHERE id = $1`, 1).Scan(
    &nullInt,
    &nullString,
    &nullBool,
    &nullTime,
    &nullFloat,
)

// Verificar y usar
if nullString.Valid {
    fmt.Printf("Phone: %s\n", nullString.String)
} else {
    fmt.Println("No phone number")
}

if nullTime.Valid {
    fmt.Printf("Last login: %v\n", nullTime.Time)
}
```

### 64.10.4 Custom Scanners

```go
import "database/sql"

// Tipo customizado
type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
    StatusPending  Status = "pending"
)

// Implementar sql.Scanner interface
func (s *Status) Scan(value interface{}) error {
    switch v := value.(type) {
    case string:
        *s = Status(v)
    case []byte:
        *s = Status(v)
    case nil:
        *s = ""
    default:
        return fmt.Errorf("cannot scan %T into Status", v)
    }
    return nil
}

// Implementar driver.Valuer interface
func (s Status) Value() (driver.Value, error) {
    return string(s), nil
}

// Usar en scanning
var status Status
err := db.QueryRow("SELECT status FROM users WHERE id = $1", 1).Scan(&status)

// Aplicación
type User struct {
    ID     int
    Name   string
    Status Status  // Automáticamente scanneable
}

var u User
db.QueryRow("SELECT id, name, status FROM users WHERE id = $1", 1).Scan(
    &u.ID, &u.Name, &u.Status)
```

### 64.10.5 JSON Unmarshaling

```go
import "encoding/json"

type Product struct {
    ID   int
    Name string
    Specs map[string]interface{}  // JSON unmarshaling
    Tags []string                  // JSON array
}

// Scanning JSON
var specs json.RawMessage
err := db.QueryRow(
    "SELECT id, name, specs FROM products WHERE id = $1", 1).Scan(
    &p.ID, &p.Name, &specs)

if err != nil {
    return err
}

// Unmarshal
var specsMap map[string]interface{}
if err := json.Unmarshal(specs, &specsMap); err != nil {
    return err
}

// O unmarshaling directo
type ProductSpecs struct {
    Color   string `json:"color"`
    Size    string `json:"size"`
    Weight  float64 `json:"weight"`
}

var specs ProductSpecs
json.Unmarshal(specsBytes, &specs)
```

---

## 64.11 - Production Patterns & Best Practices

### 64.11.1 Migration Management

```go
// Schema migration pattern
package migrations

import (
    "database/sql"
    "embed"
    "fmt"
    "log"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Migration struct {
    Version string
    SQL     string
}

func LoadMigrations() ([]Migration, error) {
    entries, err := migrations.ReadDir("migrations")
    if err != nil {
        return nil, err
    }
    
    var migs []Migration
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        
        content, err := migrations.ReadFile("migrations/" + entry.Name())
        if err != nil {
            return nil, err
        }
        
        migs = append(migs, Migration{
            Version: entry.Name(),
            SQL:     string(content),
        })
    }
    
    return migs, nil
}

func ApplyMigrations(db *sql.DB) error {
    // Create migrations table
    _, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMP DEFAULT NOW()
    )
    `)
    if err != nil {
        return err
    }
    
    // Load and apply migrations
    migs, err := LoadMigrations()
    if err != nil {
        return err
    }
    
    for _, mig := range migs {
        // Check if already applied
        var exists bool
        err := db.QueryRow(
            "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)",
            mig.Version).Scan(&exists)
        
        if err != nil {
            return err
        }
        
        if !exists {
            log.Printf("Applying migration: %s", mig.Version)
            
            tx, err := db.Begin()
            if err != nil {
                return err
            }
            
            // Apply SQL
            if _, err := tx.Exec(mig.SQL); err != nil {
                tx.Rollback()
                return err
            }
            
            // Mark as applied
            _, err = tx.Exec(
                "INSERT INTO schema_migrations(version) VALUES($1)",
                mig.Version)
            if err != nil {
                tx.Rollback()
                return err
            }
            
            if err := tx.Commit(); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

### 64.11.2 Connection String Handling

```go
import "os"

// Load connection string from environment
func getConnString() string {
    // Preferencia: variable de env > default
    if conn := os.Getenv("DATABASE_URL"); conn != "" {
        return conn
    }
    
    // Build from components
    host := os.Getenv("DB_HOST")
    if host == "" {
        host = "localhost"
    }
    
    port := os.Getenv("DB_PORT")
    if port == "" {
        port = "5432"
    }
    
    user := os.Getenv("DB_USER")
    if user == "" {
        user = "postgres"
    }
    
    password := os.Getenv("DB_PASSWORD")
    if password == "" {
        password = "postgres"
    }
    
    dbname := os.Getenv("DB_NAME")
    if dbname == "" {
        dbname = "mydb"
    }
    
    sslMode := os.Getenv("DB_SSLMODE")
    if sslMode == "" {
        sslMode = "disable"  // Production: use "require"
    }
    
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        user, password, host, port, dbname, sslMode)
}

// Sensitive data handling
func sanitizeConnString(connStr string) string {
    // Remove password for logging
    re := regexp.MustCompile(`password=([^@&]+)`)
    return re.ReplaceAllString(connStr, "password=****")
}
```

### 64.11.3 Monitoring & Logging

```go
import (
    "log"
    "time"
)

type QueryLogger struct {
    db *sql.DB
}

// Log con timing
func (ql *QueryLogger) QueryWithLogging(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    start := time.Now()
    
    rows, err := ql.db.QueryContext(ctx, query, args...)
    
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("Query error after %v: %s [%v]", duration, query, args)
        return nil, err
    }
    
    if duration > 1*time.Second {
        log.Printf("SLOW QUERY (%v): %s", duration, query)
    }
    
    return rows, nil
}

// Pool statistics monitoring
func (ql *QueryLogger) LogPoolStats() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := ql.db.Stats()
        
        log.Printf("DB Pool Stats: OpenConns=%d, InUse=%d, Idle=%d, MaxOpen=%d, MaxLifetime=%v",
            stats.OpenConnections,
            stats.InUse,
            stats.Idle,
            stats.MaxOpenConnections,
            stats.MaxConnectionLifetime,
        )
    }
}

// Error categorization
func CategorizeError(err error) string {
    if err == sql.ErrNoRows {
        return "not_found"
    }
    
    if err == context.DeadlineExceeded {
        return "timeout"
    }
    
    if pgErr, ok := err.(*pq.Error); ok {
        return string(pgErr.Code)
    }
    
    return "unknown"
}
```

### 64.11.4 Case Studies - Raw SQL vs ORM

```

 CASE STUDY 1: Simple CRUD API                                 │

 Requisito: Basic GET, POST, PUT, DELETE endpoints             │
                                                                │
 RAW SQL:                                                       │
 ✓ Más boilerplate                                             │
 ✓ Mejor performance para simple queries                       │
 ✗ Más código para scanning                                    │
 Veredicto: ORM es mejor (menos código, performance suficiente)│
                                                                │
 GORM:                                                          │
 ✓ Mínimo código                                               │
 ✓ Type-safe (compile-time)                                   │
 ✓ Automático nil handling                                     │
 ✓ Migrations included                                         │
                                                                │



 CASE STUDY 2: Complex Analytics Queries                       │

 Requisito: CTEs, window functions, json aggregations          │
                                                                │
 RAW SQL:                                                       │
 ✓ Control total                                               │
 ✓ Optimal performance                                         │
 ✓ Database-specific features                                  │
 ✓ Easy debugging/optimization                                 │
 Veredicto: Raw SQL es ESENCIAL                                │
                                                                │
 GORM:                                                          │
 ✗ Difícil expresar lógica compleja                           │
 ✗ Puede generar ineficientes queries                          
 ✗ Debugging complejo                                          │
                                                                │



 CASE STUDY 3: Microservice with Single Responsibility         │

 Requisito: Endpoint específico, queries muy particulares      │
                                                                │
 Raw SQL:                                                       │
 ✓ Queries tailored para use case                             │
 ✓ Zero abstraction overhead                                   │
 ✓ Easy to understand without ORM knowledge                    │
 Veredicto: Raw SQL ideal para microservices                   │
                                                                │
 sqlc (Code generation):                                        │
 ✓ SQL verificado en compile-time                             │
 ✓ Type-safe con Zero runtime overhead                        │
 ✓ Fácil testing                                               │
 ✓ Mejor que raw SQL + GORM                                   │
                                                                │

```

### 64.11.5 Troubleshooting

**Problema: Conexiones agotadas**

```go
// Síntomas: "too many connections" error

// Solución 1: Revisar pool configuration
if err := db.Ping(); err != nil {
    return err
}

stats := db.Stats()
log.Printf("Active connections: %d / %d", stats.OpenConnections, stats.MaxOpenConnections)

// Solución 2: Reducir MaxOpenConns
db.SetMaxOpenConns(25)  // Fue 100

// Solución 3: Implementar connection waiting
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
row := db.QueryRowContext(ctx, "SELECT 1")
```

**Problema: Queries lentas**

```go
// Síntomas: Timeouts, usuario slow queries

// Paso 1: Analizar query
const analyzeQuery = `
EXPLAIN ANALYZE
SELECT * FROM orders
WHERE created_at > NOW() - INTERVAL '30 days'
AND status = 'completed'
`

// Paso 2: Verificar índices
const checkIndices = `
SELECT indexname FROM pg_indexes
WHERE tablename = 'orders'
`

// Paso 3: Crear índice si falta
_, _ = db.Exec(`
CREATE INDEX CONCURRENTLY idx_orders_created_status
ON orders(created_at DESC, status)
WHERE status = 'completed'
`)

// Paso 4: Verify improvement
const timeQuery = `
SELECT EXTRACT(MILLISECOND FROM (
    SELECT COUNT(*) FROM orders
    WHERE created_at > NOW() - INTERVAL '30 days'
))
`
```

**Problema: Deadlocks**

```go
// Síntomas: "deadlock detected" errors

// Solución 1: Visualizar transacciones en conflicto
const findDeadlocks = `
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocking_locks.pid AS blocking_pid,
    blocked_activity.query AS blocked_query,
    blocking_activity.query AS blocking_query
FROM pg_locks AS blocked_locks
JOIN pg_statements AS blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_locks AS blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
    AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
    AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
    AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
    AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
    AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
    AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
    AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
    AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
    AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
    AND blocking_locks.pid != blocked_locks.pid
JOIN pg_statements AS blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted
`

// Solución 2: Reducir transaction scope
// ❌ Largo
tx, _ := db.Begin()
doComplexComputations()  // Largas operaciones
tx.Exec(...)
tx.Commit()

// ✅ Corto
doComplexComputations()  // Fuera de tx
tx, _ := db.Begin()
tx.Exec(...)
tx.Commit()

// Solución 3: Consistent ordering
// ❌ Potencial deadlock
tx1: UPDATE users SET balance = balance - 100 WHERE id = 1
tx1: UPDATE users SET balance = balance + 100 WHERE id = 2

tx2: UPDATE users SET balance = balance + 50 WHERE id = 2
tx2: UPDATE users SET balance = balance - 50 WHERE id = 1

// ✅ Sin deadlock (mismo orden)
tx1: UPDATE users SET balance = ... WHERE id IN (1, 2) ORDER BY id

// Solución 4: Retry con backoff exponencial
err := withDeadlockRetry(ctx, db, transactionFn)
```

---

## Ejercicios Progresivos

### Ejercicio 1: Basic CRUD with database/sql

**Requisito:** Implementar operaciones CRUD completas usando raw SQL.

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

type User struct {
    ID    int
    Name  string
    Email string
}

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", "postgres://user:pass@localhost/testdb")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create table
    db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL
    )
    `)
}

// CREATE
func CreateUser(name, email string) (*User, error) {
    var id int
    err := db.QueryRow(
        "INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
        name, email).Scan(&id)
    
    if err != nil {
        return nil, err
    }
    
    return &User{ID: id, Name: name, Email: email}, nil
}

// READ
func GetUser(id int) (*User, error) {
    var u User
    err := db.QueryRow(
        "SELECT id, name, email FROM users WHERE id = $1", id).
        Scan(&u.ID, &u.Name, &u.Email)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    return &u, err
}

// UPDATE
func UpdateUser(id int, name, email string) error {
    result, err := db.Exec(
        "UPDATE users SET name = $1, email = $2 WHERE id = $3",
        name, email, id)
    
    if err != nil {
        return err
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}

// DELETE
func DeleteUser(id int) error {
    result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
    
    if err != nil {
        return err
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}

func main() {
    defer db.Close()
    
    // Test CRUD
    u, _ := CreateUser("Alice", "alice@example.com")
    fmt.Printf("Created: %+v\n", u)
    
    u, _ = GetUser(u.ID)
    fmt.Printf("Retrieved: %+v\n", u)
    
    UpdateUser(u.ID, "Alice Updated", "alice.new@example.com")
    u, _ = GetUser(u.ID)
    fmt.Printf("Updated: %+v\n", u)
    
    DeleteUser(u.ID)
    _, err := GetUser(u.ID)
    fmt.Printf("Deleted: err=%v\n", err)
}
```

### Ejercicio 2: Prepared Statements & Performance

**Requisito:** Comparar performance con y sin prepared statements.

```go
// TODO: Implementar benchmark
// Insertar 10,000 usuarios y medir tiempo con/sin prepared statements
// Resultado esperado: 3x más rápido con prepared statements
```

### Ejercicio 3: Complex Queries with Joins

**Requisito:** Ejecutar queries complejas con múltiples joins y aggregations.

```go
// TODO: Crear schema con users, orders, items
// Query: Top 10 usuarios por dinero gastado en últimos 30 días
// Incluir: COUNT(orders), SUM(total), AVG(order_value)
```

### Ejercicio 4: Transactions & Error Handling

**Requisito:** Implementar transacción segura con proper error handling.

```go
// TODO: Transacción de transferencia de dinero
// - Debit from account A
// - Credit to account B
// - Log transaction
// - Retry on deadlock
// - Rollback on error
```

### Ejercicio 5: Production App Monitoring

**Requisito:** App con connection pooling, query logging, y health checks.

```go
// TODO: Servidor HTTP con:
// - Connection pool optimizado
// - Query logging con timing
// - Pool stats monitoring
// - Health check endpoint
// - Graceful shutdown
```

---

## Resumen de Best Practices

```
 HAGA:

1. Usar Prepared Statements SIEMPRE
   db.Prepare() o batch para seguridad contra inyección SQL

2. Configurar Connection Pool
   MaxOpenConns = CPU_CORES * 4
   MaxIdleConns = MaxOpenConns / 5

3. Manejar Errores Específicamente
   if err == sql.ErrNoRows { }
   if pgErr, ok := err.(*pq.Error) { }

4. Usar Context para Timeouts
   ctx, cancel := context.WithTimeout(...)
   defer cancel()

5. Defer db.Close()
   defer db.Close()

6. Log Slow Queries
   track query duration, alert if > 1s

7. Test with Real Database
   TestContainers > Fixtures > Mocks

8. Monitor Pool Stats
   OpenConnections, InUse, Idle metrics


 NO HAGA:

1. String concatenation para queries
   ❌ "SELECT * FROM users WHERE id = " + userID
   ✅ "SELECT * FROM users WHERE id = $1", userID

2. Ignorar sql.ErrNoRows
   ❌ if err != nil { return err }
   ✅ if err == sql.ErrNoRows { /* 404 */ }

3. Reabrir conexiones constantemente
   ❌ for { db, _ := sql.Open(...); db.Query(...) }
   ✅ Usar singleton sql.DB con pooling

4. No configurar MaxOpenConns
   ❌ Dejar default ilimitado
   ✅ Explícitamente establecer límite

5. No usar transactions para operaciones relacionadas
   ❌ Múltiples queries independientes
   ✅ Agrupar en transacción

6. No testear con database real
   ❌ Solo mocks
   ✅ Integration tests con DB real
```

---

## Diagramas de Arquitectura

```
FLUJO DE QUERY EXECUTION:

User Code
    │
    ├─ db.Query("SELECT * FROM users WHERE id = ?", 1)
    │
    ▼
Connection Pool Manager
    │
    ├─ ¿Conexión idle disponible?
    │  └─ SÍ → Reutilizar
    │  └─ NO → Crear nueva (si < MaxOpenConns)
    │
    ▼
Network → PostgreSQL Server
    │
    ├─ Parse Query
    ├─ Plan Query
    ├─ Execute
    │
    ▼
Results
    │
    ├─ rows.Next()
    ├─ rows.Scan(&...variables...)
    │
    ▼
Cleanup
    │
    ├─ rows.Close()
    ├─ Retorna conexión al pool


TRANSACTION LIFECYCLE:

BEGIN TRANSACTION
    │
    ├─ Isolation Level: READ COMMITTED
    │
    ├─ Statement 1: UPDATE ...
    │  └─ Si error → ROLLBACK (fin)
    │
    ├─ Statement 2: INSERT ...
 Si error → ROLLBACK (fin)    │  
    │
    ├─ Statement 3: UPDATE ...
    │  └─ Si error → ROLLBACK (fin)
    │
    ├─ COMMIT
    │  └─ Si error → ROLLBACK (fin)
    │
    ▼
CONFIRMADO O ANULADO
```

---

**Fin del CAPÍTULO 64**

Word Count: ~1,950 líneas | Tamaño: ~45 KB


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/64-database-drivers-y-raw-sql/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/64-database-drivers-y-raw-sql):

```bash
cd examples/64-database-drivers-y-raw-sql
go run .
```
