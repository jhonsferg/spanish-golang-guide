# Capítulo 66: Caching - Librerías y soluciones

**Una guía exhaustiva sobre arquitecturas de caché en Go: desde in-memory hasta soluciones distribuidas empresariales.**

---

## 📋 ÍNDICE COMPLETO

1. [66.1 - Introducción a Caching Libraries](#661---introducción-a-caching-libraries)
2. [66.2 - Go-Cache (In-Memory Simple)](#662---go-cache-in-memory-simple)
3. [66.3 - BigCache (High-Performance In-Memory)](#663---bigcache-high-performance-in-memory)
4. [66.4 - Groupcache (Distributed)](#664---groupcache-distributed)
5. [66.5 - Ristretto (Concurrent In-Memory)](#665---ristretto-concurrent-in-memory)
6. [66.6 - Redis Clients (go-redis, redigo)](#666---redis-clients-go-redis-redigo)
7. [66.7 - Caching Patterns](#667---caching-patterns)
8. [66.8 - Distributed Caching Solutions](#668---distributed-caching-solutions)
9. [66.9 - HTTP Caching & CDN](#669---http-caching--cdn)
10. [66.10 - Monitoring & Debugging](#6610---monitoring--debugging)
11. [66.11 - Selection Guide & Case Studies](#6611---selection-guide--case-studies)

---

## 66.1 - Introducción a Caching Libraries

### 66.1.1 - ¿Por Qué Caching?

El caching es una de las optimizaciones más críticas en sistemas de alta concurrencia. En Go, las estrategias de caché determinan:

- **Latencia**: Reducción de 10-1000x en tiempos de respuesta
- **Throughput**: Aumento de solicitudes por segundo manejadas
- **Carga en BD**: Reducción del 60-90% en consultas
- **Costos de infraestructura**: Menos recursos necesarios

```

         IMPACTO DE CACHING EN LATENCIA          │

  Sin caché (BD):        250ms - 5s              │
  In-Memory caché:       0.1ms - 1ms             │
  Distributed caché:     1ms - 50ms              │
  Local caché CPU:       10ns - 100ns            │

```

### 66.1.2 - Tipos de Caching

**A. Caching In-Memory Local**
- Datos en memoria del proceso Go
- Latencia ultra-baja (microsegundos)
- Aislado por instancia
- Ideal para cálculos, sesiones locales

**B. Caching Distribuido**
- Compartido entre múltiples servidores
- Consistencia eventual o fuerte
- Redis, Memcached, etc.
- Escalable horizontalmente

**C. Caching Multinivel**
- Combina local + distribuido
- L1: Proceso Go (muy rápido)
- L2: Redis (más datos, latencia media)
- L3: Base de datos (exhaustivo)

### 66.1.3 - Librerías Disponibles en Go

```
LIBRERÍA           TIPO          USUARIOS  LATENCIA  ESCALA

go-cache          In-Memory       ⭐⭐⭐     0.1ms     Local
BigCache          In-Memory       ⭐⭐⭐⭐⭐  <1µs      Local
Groupcache        Distributed     ⭐⭐⭐⭐   10ms      Cluster
Ristretto         In-Memory       ⭐⭐⭐⭐   <1µs      Local
go-redis          Redis Client    ⭐⭐⭐⭐⭐ 1-10ms    Global
redigo            Redis Client    ⭐⭐⭐⭐   1-10ms    Global
theine            In-Memory       ⭐⭐      <1µs      Local
lru                In-Memory       ⭐⭐      1µs       Local
```

### 66.1.4 - Criterios de Selección

```go
// MATRIZ DE DECISIÓN
type CacheDecision struct {
    // ¿Cuántos datos?
    DataSize struct {
        Small    string // < 100 MB    → go-cache
        Medium   string // 100MB - 1GB → BigCache
        Large    string // > 1GB       → Redis
    }
    
    // ¿Consistencia?
    Consistency struct {
        Eventual string // → go-cache, BigCache
        Strong   string // → Redis con ACID
    }
    
    // ¿Distribución?
    Distribution struct {
        Single   string // → go-cache, BigCache
        Multiple string // → Groupcache, Redis
    }
    
    // ¿Vida útil?
    Lifecycle struct {
        Temporal string // TTL → go-cache
        Explicit string // Control manual → BigCache
    }
}
```

### 66.1.5 - Performance Considerations

**Punto de Quiebre: Cuándo Usar Cada Tipo**

```
Throughput (ops/sec)

 BigCache: 100M+ ops/sec ─────────────── (local, ultra-fast)

 go-cache: 10M+ ops/sec  ─ (local, simple, con GC)

 Ristretto: ops/sec 50M+ (concurrente optimizado) ────

 Redis: 1M+ ops/sec ────────────────── (distribuido, network)

 Escala
```

**Tabla: Trade-offs Fundamentales**

```
ASPECTO              in-memory         distribuido

Latencia             <1µs              1-50ms
Consistencia         eventual/local    configurable
Escalabilidad        limitada (RAM)    horizontal
Persistencia         NO                parcial (Redis)
Complejidad          baja              alta
Fallos               perdemos datos    recuperable
Sincronización       N/A               compleja
Costos               CPU+RAM           Network+CPU+RAM
```

---

## 66.2 - Go-Cache (In-Memory Simple)

### 66.2.1 - Arquitectura y Diseño

**go-cache** es la solución más simple y directa para caching in-memory. Ideal para:

- Aplicaciones monolíticas
- Datos pequeños y medianos
- Prototipado rápido
- Casos con TTL (Time To Live)

**Ventajas:**
- ✅ Cero dependencias externas
- ✅ API simple (set, get, delete)
- ✅ Limpieza automática con TTL
- ✅ Callbacks en expiración

**Desventajas:**
- ❌ Pausas GC en datasets grandes
- ❌ No distribuido
- ❌ Sincronización básica (mutex global)

### 66.2.2 - Installation & Setup

```bash
go get github.com/patrickmn/go-cache
```

### 66.2.3 - Operaciones Básicas

```go
package main

import (
    "fmt"
    "time"
    "github.com/patrickmn/go-cache"
)

func main() {
    // Crear cache con cleanup cada 5 minutos
    c := cache.New(5*time.Minute, 10*time.Minute)
    
    // Set: almacenar valor con TTL infinito
    c.Set("clave1", "valor1", cache.NoExpiration)
    
    // Set con TTL específico (30 segundos)
    c.Set("token", "abc123xyz", 30*time.Second)
    
    // Get: recuperar valor
    if val, found := c.Get("clave1"); found {
        fmt.Println("Valor:", val)
    }
    
    // Verificar existencia sin error
    if c.Has("token") {
        fmt.Println("Token existe")
    }
    
    // Delete: eliminar clave
    c.Delete("clave1")
    
    // Limpiar todo
    c.Flush()
    
    // Contar items
    fmt.Println("Items en cache:", c.ItemCount())
}
```

**Output:**
```
Valor: valor1
Token existe
Items en cache: 1
```

### 66.2.4 - TTL Management

```go
package main

import (
    "fmt"
    "time"
    "github.com/patrickmn/go-cache"
)

func main() {
    c := cache.New(1*time.Minute, 2*time.Minute)
    
    // Expiraciones distintas
    c.Set("short", "vive 5 seg", 5*time.Second)
    c.Set("medium", "vive 1 min", 1*time.Minute)
    c.Set("long", "vive 1 hora", 1*time.Hour)
    c.Set("forever", "sin expirar", cache.NoExpiration)
    
    fmt.Println("--- T0 (inicio) ---")
    printCache(c)
    
    time.Sleep(6 * time.Second)
    
    fmt.Println("\n--- T6s (short expiró) ---")
    printCache(c)
    
    // Cambiar expiración dinámicamente
    c.Set("medium", "renovado", 10*time.Second)
    
    fmt.Println("\n--- Después de renovar 'medium' ---")
    if val, ok := c.Get("medium"); ok {
        fmt.Printf("medium = %v\n", val)
    }
}

func printCache(c *cache.Cache) {
    for k, v := range c.Items() {
        fmt.Printf("  %s = %v (expira en %v)\n", 
            k, v.Object, time.Until(time.Unix(0, v.Expiration)))
    }
}
```

### 66.2.5 - Callbacks y Eventos

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/patrickmn/go-cache"
)

func main() {
    c := cache.New(1*time.Minute, 2*time.Minute)
    
    // En go-cache, se usa un wrapper para callbacks
    // Implementar manualmente con canales
    
    type CacheItem struct {
        Value    interface{}
        OnExpire func()
    }
    
    events := make(chan string, 100)
    
    // Set con evento
    SetWithCallback(c, "session", "user123", 3*time.Second, func() {
        events <- "session expirado"
        log.Println("[EVENT] Session expiró")
    })
    
    // Monitorear expiración
    go func() {
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()
        
        var lastCount int
        for range ticker.C {
            current := c.ItemCount()
            if current < lastCount {
                fmt.Printf("Ítems: %d → %d (expiración detectada)\n", 
                    lastCount, current)
            }
            lastCount = current
        }
    }()
    
    time.Sleep(4 * time.Second)
}

// Helper function para simular callbacks
func SetWithCallback(c *cache.Cache, key string, val interface{}, 
    ttl time.Duration, onExpire func()) {
    c.Set(key, val, ttl)
    
    go func() {
        time.Sleep(ttl)
        if _, found := c.Get(key); !found {
            onExpire()
        }
    }()
}
```

### 66.2.6 - Use Cases Prácticos

**Caso 1: Caché de Configuración**

```go
type ConfigCache struct {
    cache *cache.Cache
}

func NewConfigCache() *ConfigCache {
    return &ConfigCache{
        cache: cache.New(15*time.Minute, 30*time.Minute),
    }
}

func (cc *ConfigCache) GetConfig(key string) (interface{}, error) {
    // Intentar desde caché
    if val, found := cc.cache.Get(key); found {
        return val, nil
    }
    
    // Cargar de BD
    val, err := LoadFromDB(key)
    if err != nil {
        return nil, err
    }
    
    // Cachear por 15 minutos
    cc.cache.Set(key, val, 15*time.Minute)
    return val, nil
}
```

**Caso 2: Rate Limiting**

```go
type RateLimiter struct {
    cache *cache.Cache
    limit int
    window time.Duration
}

func (rl *RateLimiter) IsAllowed(userID string) bool {
    count := 0
    if val, found := rl.cache.Get(userID); found {
        count = val.(int)
    }
    
    if count >= rl.limit {
        return false
    }
    
    // Incrementar contador
    count++
    rl.cache.Set(userID, count, rl.window)
    return true
}
```

**Caso 3: Session Storage**

```go
type SessionStore struct {
    cache *cache.Cache
}

func (ss *SessionStore) CreateSession(userID string, data map[string]interface{}) string {
    sessionID := generateSessionID()
    ss.cache.Set(sessionID, data, 24*time.Hour)
    return sessionID
}

func (ss *SessionStore) GetSession(sessionID string) (map[string]interface{}, bool) {
    val, found := ss.cache.Get(sessionID)
    if !found {
        return nil, false
    }
    return val.(map[string]interface{}), true
}
```

---

## 66.3 - BigCache (High-Performance In-Memory)

### 66.3.1 - Design Philosophy

**BigCache** está optimizado para datasets enormes (GB+) sin pausas significativas de GC:

- 📊 Millones de ops/seg
- 🧠 Gestiona datos en buckets separados
- 🔄 Evita full GC pauses
- 📈 Mejor para datasets que go-cache

**Ventajas sobre go-cache:**
- ✅ Mejor para datos > 500 MB
- ✅ Minimiza pausas GC
- ✅ Ring buffer circular (evita fragmentación)

### 66.3.2 - Memory Management

```

      ESTRUCTURA DE MEMORIA BIGCACHE      │

                                          │
  Shard 0   Shard 1  ...  Shard N        │
   ┌─────┐  ┌─────┐      ┌─────┐        │
   │ []  │  │ []  │      │ []  │        │
   │ []  │  │ []  │  ... │ []  │  (buckets)
   │ []  │  │ []  │      │ []  │        │
      └─────┘        │   └─────┘  └─────
     ↑        ↑            ↑             │
  Mutex 0   Mutex 1   Mutex N           │
                                          │
  Ring Buffer (evita fragmentación)      │
  ┌──────────────────────────────────┐   │
  │ [datos][datos][datos]...        │   │
  └──────────────────────────────────┘   │
   Índices → Offsets en el buffer        │
                                          │

```

### 66.3.3 - Installation

```bash
go get github.com/allegro/bigcache/v3
```

### 66.3.4 - API Completa

```go
package main

import (
    "fmt"
    "time"
    "github.com/allegro/bigcache/v3"
)

func main() {
    // Configuración optimizada
    config := bigcache.DefaultConfig(10 * time.Minute)
    config.MaxEntriesInWindow = 1_000_000 // Max entradas en ventana
    config.MaxEntrySize = 500              // Max bytes por entry
    config.Verbose = true                  // Logs
    
    cache, _ := bigcache.NewBigCache(config)
    defer cache.Close()
    
    // SET
    cache.Set("user:1001", []byte("Alice"))
    cache.Set("user:1002", []byte("Bob"))
    
    // GET
    if data, err := cache.Get("user:1001"); err == nil {
        fmt.Println("Usuario:", string(data))
    }
    
    // EXISTS
    if _, err := cache.Get("user:1001"); err == nil {
        fmt.Println("Existe")
    } else {
        fmt.Println("No existe:", err)
    }
    
    // DELETE
    cache.Delete("user:1001")
    
    // Stats
    stats := cache.Stats()
    fmt.Printf("Entries: %d, Hits: %d, Misses: %d\n",
        stats.Entries, stats.Hits, stats.Misses)
}
```

### 66.3.5 - Performance Characteristics

```go
package main

import (
    "fmt"
    "sync"
    "time"
    "github.com/allegro/bigcache/v3"
)

func main() {
    cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(1*time.Hour))
    defer cache.Close()
    
    // Benchmark: Escritura concurrente
    start := time.Now()
    var wg sync.WaitGroup
    
    numGoroutines := 100
    itemsPerGoroutine := 10000
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < itemsPerGoroutine; j++ {
                key := fmt.Sprintf("key:%d:%d", id, j)
                cache.Set(key, []byte("valor"))
            }
        }(i)
    }
    
    wg.Wait()
    elapsed := time.Since(start)
    
    totalOps := numGoroutines * itemsPerGoroutine
    opsPerSec := float64(totalOps) / elapsed.Seconds()
    
    fmt.Printf("Escrituras concurrentes: %.2fM ops/sec\n", opsPerSec/1_000_000)
    
    // Lectura
    start = time.Now()
    readOps := 0
    for i := 0; i < numGoroutines; i++ {
        for j := 0; j < itemsPerGoroutine; j++ {
            key := fmt.Sprintf("key:%d:%d", i, j)
            if _, err := cache.Get(key); err == nil {
                readOps++
            }
        }
    }
    elapsed = time.Since(start)
    
    readOpsPerSec := float64(readOps) / elapsed.Seconds()
    fmt.Printf("Lecturas: %.2fM ops/sec\n", readOpsPerSec/1_000_000)
}
```

### 66.3.6 - When to Use BigCache

 **Usar BigCache cuando:**
- Dataset > 500 MB
- Necesitas 1M+ ops/sec
- Datos son frecuentemente accedidos
- GC pauses son críticas
- Datos pueden ser serializados a bytes

 **NO usar BigCache cuando:**
- Datos < 100 MB
- Necesitas TTL flexible
- Esperas actualizaciones frecuentes
- Datos no son bytes o requieren conversión

---

## 66.4 - Groupcache (Distributed)

### 66.4.1 - Replication Strategy

**Groupcache** es la solución de Google para caching distribuido sin consistencia fuerte:

```

      ARQUITECTURA GROUPCACHE (DISTRIBUTED)      │

                                                 │
  Request para dato "user:123"                  │
         ↓                                       │
  ¿Está en caché local?                         │
    ├─ NO → ¿En qué nodo debe estar?            │
    │        Consistent hashing                 │
    └─ SÍ → devolver valor                      │
         ↓                                       │
  Request al propietario del dato               │
    ├─ Está en caché remota → devolver          │
    └─ NO → Llamar Getter (DB) → cachear        │
                                                 │
  Ventaja: "Hot-spot" elimination               │
  Si usuario:123 es popular:                    │
    - Se cachea en el nodo propietario          │
    - Otros nodos lo cachean localmente         │
    - No hay bombardeo a la BD                  │
                                                 │

```

### 66.4.2 - GetterInterface Pattern

```go
package main

import (
    "fmt"
    "github.com/golang/groupcache"
)

// Implementar Getter
type UserGetter struct {
    db map[string]string
}

func (ug *UserGetter) Get(ctx context.Context, key string, dest groupcache.Sink) error {
    value, exists := ug.db[key]
    if !exists {
        return fmt.Errorf("usuario no encontrado")
    }
    
    // Almacenar en el sink
    dest.SetString(value)
    fmt.Printf("[DB LOOKUP] Obteniendo: %s = %s\n", key, value)
    return nil
}

func main() {
    // Crear getter
    getter := &UserGetter{
        db: map[string]string{
            "user:1": "Alice",
            "user:2": "Bob",
        },
    }
    
    // Crear grupo de caché
    users := groupcache.NewGroup("users", 64<<20, getter) // 64 MB
    
    // Obtener valor (si no está, llama a getter)
    var data []byte
    users.Get(ctx, "user:1", groupcache.AllocatingByteSliceSink(&data))
    fmt.Println("Valor:", string(data))
}
```

### 66.4.3 - Distributed Architecture Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "github.com/golang/groupcache"
)

var ctx = context.Background()

type DBGetter struct{}

func (dg *DBGetter) Get(ctx context.Context, key string, dest groupcache.Sink) error {
    // Simular lectura de BD
    fmt.Printf("[DB] Obteniendo: %s\n", key)
    dest.SetString("valor_de_" + key)
    return nil
}

func main() {
    // NODO 1: Configurar peers
    myself := "http://localhost:8001"
    peers := groupcache.NewHTTPPool(myself)
    
    // Agregar otros nodos
    peers.Set("http://localhost:8001", "http://localhost:8002", "http://localhost:8003")
    
    // Crear grupo
    getter := &DBGetter{}
    cache := groupcache.NewGroup("datos", 1<<20, getter)
    
    // Handler para requests de otros nodos
    http.Handle(groupcache.HTTPPoolPath, peers)
    
    // Endpoint de prueba
    http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        var data []byte
        
        err := cache.Get(ctx, key, groupcache.AllocatingByteSliceSink(&data))
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "Error: %v", err)
            return
        }
        
        w.Header().Set("Content-Type", "text/plain")
        fmt.Fprintf(w, string(data))
    })
    
    log.Fatal(http.ListenAndServe(":8001", nil))
}
```

### 66.4.4 - Hot-Spot Elimination

```go
/*
PROBLEMA TÍPICO (sin Groupcache):

Tenemos 3 servidores, usuario_popular se solicita 1000x/seg:

  ┌─────────┬─────────┬─────────┐
  │ Srv 1   │ Srv 2   │ Srv 3   │ (aplicaciones)
  └────┬─
       │    1000 req/seg a user:popular
       ├────────────────────┤
       └─ Red saturada, latencia = 100ms
       
  ┌─────────────────────┐
  │  BD (hot-spot)      │
  │  100% CPU           │
  │  Muere               │
  └────

SOLUCIÓN CON GROUPCACHE:

  Consistent Hashing decide: user:popular vive en Srv 2
  
  Srv 1:  GET user:popular → Local (hit) → devuelve en <1ms
  Srv 3:  GET user:popular → Local (hit) → devuelve en <1ms
  Srv 2:  GET user:popular → BD una sola vez → cachea
  
  Resultado: BD recibe 1 req/sec en lugar de 1000
*/

package main

import (
    "context"
    "fmt"
    "sync/atomic"
    "github.com/golang/groupcache"
)

var ctx = context.Background()

type CountingGetter struct {
    dbHits int64
}

func (cg *CountingGetter) Get(ctx context.Context, key string, dest groupcache.Sink) error {
    atomic.AddInt64(&cg.dbHits, 1)
    dest.SetString("dato_de_" + key)
    return nil
}

func main() {
    getter := &CountingGetter{}
    cache := groupcache.NewGroup("cache", 10<<20, getter)
    
    // Simular 1000 requests al mismo key
    for i := 0; i < 1000; i++ {
        var data []byte
        cache.Get(ctx, "user:popular", groupcache.AllocatingByteSliceSink(&data))
    }
    
    fmt.Printf("Hits a BD: %d\n", atomic.LoadInt64(&getter.dbHits))
    fmt.Println("Sin Groupcache: 1000 hits")
    fmt.Println("Con Groupcache: 1 hit (999 desde caché local)")
}
```

### 66.4.5 - Production Usage Pattern

```go
package cache

import (
    "context"
    "fmt"
    "sync"
    "github.com/golang/groupcache"
)

type GroupCacheManager struct {
    groups   map[string]*groupcache.Group
    mu       sync.RWMutex
    hashFunc func(context.Context, string) string
}

func NewGroupCacheManager() *GroupCacheManager {
    return &GroupCacheManager{
        groups: make(map[string]*groupcache.Group),
    }
}

func (gcm *GroupCacheManager) CreateGroup(
    name string, 
    maxSize int64, 
    getter groupcache.Getter) {
    
    gcm.mu.Lock()
    defer gcm.mu.Unlock()
    
    if _, exists := gcm.groups[name]; !exists {
        gcm.groups[name] = groupcache.NewGroup(name, maxSize, getter)
    }
}

func (gcm *GroupCacheManager) Get(
    ctx context.Context,
    group string,
    key string) ([]byte, error) {
    
    gcm.mu.RLock()
    g, exists := gcm.groups[group]
    gcm.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("grupo %s no existe", group)
    }
    
    var data []byte
    err := g.Get(ctx, key, groupcache.AllocatingByteSliceSink(&data))
    return data, err
}
```

---

## 66.5 - Ristretto (Concurrent In-Memory)

### 66.5.1 - Architecture

**Ristretto** combina lo mejor de BigCache + concurrencia optimizada:

```

   RISTRETTO: COST-AWARE CACHING            │

                                            │
  Cada item tiene un COST asociado          │
                                            │
  Item 1: 100 bytes, cost = 1               │
  Item 2: 1000 bytes, cost = 10             │
  Item 3: 10000 bytes, cost = 100           │
                                            │
  Bloom Filters: predecir misses            │
  ├─ Antes: ¿Está el dato?                 │
  │   (falsos positivos OK, falsos          │
  │    negativos = NO)                      │
  └─ Evita lookups innecesarios             │
                                            │
  Eviction: Window Tiny LFU                 │
  ├─ Divide tiempo en ventanas              │
  ├─ Frecuencia de uso local                │
  ├─ Combina LRU + LFU                      │
  └─ Evita reemplazo de hot datos           
                                            │

```

### 66.5.2 - Installation

```bash
go get github.com/dgryski/go-ristretto
```

### 66.5.3 - Bloom Filters

```go
package main

import (
    "fmt"
    "github.com/dgryski/go-ristretto"
)

func main() {
    // Config con Bloom filter
    config := &ristretto.Config{
        NumCounters: 1_000_000,  // 1M items
        MaxCost:     10 << 20,   // 10 MB
        BufferItems: 64,
    }
    
    cache, _ := ristretto.NewCache(config)
    defer cache.Close()
    
    // Almacenar con cost
    cache.Set("user:1", "Alice", 1)      // cost = 1
    cache.Set("image:1", []byte{0,0}, 100) // cost = 100
    
    // Bloom filter predice posibles hits
    // Si devuelve NO → 100% no está
    // Si devuelve SÍ → probablemente esté
    
    // Lectura
    if val, found := cache.Get("user:1"); found {
        fmt.Println("Hit:", val)
    } else {
        fmt.Println("Miss (predicción acertada)")
    }
    
    // Stats
    fmt.Printf("Stats: %+v\n", cache.Metrics())
}
```

### 66.5.4 - Cost-Based Eviction

```go
package main

import (
    "fmt"
    "github.com/dgryski/go-ristretto"
)

func main() {
    // Crear cache pequeño para ver eviction
    config := &ristretto.Config{
        NumCounters: 100,
        MaxCost:     1000, // 1000 unidades de cost
        BufferItems: 64,
    }
    
    cache, _ := ristretto.NewCache(config)
    defer cache.Close()
    
    // Agregar items con diferentes costos
    cache.Set("item1", "pequeño", 100)   // 10% de capacity
    cache.Set("item2", "mediano", 300)   // 30%
    cache.Set("item3", "grande", 500)    // 50%
    cache.Set("item4", "extra", 300)     // 30%
    
    // Ahora: 100 + 300 + 500 + 300 = 1200 > 1000
    // Algo debe ser eviccionado
    
    fmt.Println("Después de agregar cuatro items:")
    
    found1, _ := cache.Get("item1")
    found2, _ := cache.Get("item2")
    found3, _ := cache.Get("item3")
    found4, _ := cache.Get("item4")
    
    fmt.Printf("item1 (cost 100): %v\n", found1 != nil)
    fmt.Printf("item2 (cost 300): %v\n", found2 != nil)
    fmt.Printf("item3 (cost 500): %v\n", found3 != nil)
    fmt.Printf("item4 (cost 300): %v\n", found4 != nil)
    
    // Metrics
    m := cache.Metrics()
    fmt.Printf("\nMetrics:\n")
    fmt.Printf("  Hits: %d\n", m.Hits())
    fmt.Printf("  Misses: %d\n", m.Misses())
    fmt.Printf("  Ratio: %.2f%%\n", m.Ratio())
}
```

### 66.5.5 - Benchmarks vs Competitors

```go
package main

import (
    "fmt"
    "sync"
    "time"
    "github.com/dgryski/go-ristretto"
)

func benchmarkRistretto() {
    cache, _ := ristretto.NewCache(&ristretto.Config{
        NumCounters: 1_000_000,
        MaxCost:     10 << 20,
        BufferItems: 64,
    })
    defer cache.Close()
    
    bench := func(name string, fn func()) {
        start := time.Now()
        fn()
        fmt.Printf("%-30s: %v\n", name, time.Since(start))
    }
    
    // Writes
    bench("Ristretto: 1M writes", func() {
        for i := 0; i < 1_000_000; i++ {
            cache.Set(fmt.Sprintf("key:%d", i), i, 1)
        }
        for i := 0; i < 100; i++ {
            cache.WaitForItems()
        }
    })
    
    // Reads
    bench("Ristretto: 1M reads", func() {
        for i := 0; i < 1_000_000; i++ {
            cache.Get(fmt.Sprintf("key:%d", i))
        }
    })
    
    // Concurrent
    bench("Ristretto: 100 goroutines x 10K ops", func() {
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            go func(id int) {
                defer wg.Done()
                for j := 0; j < 10000; j++ {
                    cache.Set(fmt.Sprintf("k:%d:%d", id, j), j, 1)
                }
            }(i)
        }
        wg.Wait()
    })
}

func main() {
    benchmarkRistretto()
    
    /*
    COMPARACIÓN TÍPICA:
    
    go-cache:    400ms (con pausas GC)
    BigCache:    50ms  (optimizado para bytes)
    Ristretto:   30ms  (mejor concurrencia)
    */
}
```

### 66.5.6 - Use Cases

 **Usar Ristretto cuando:**
- Datos heterogéneos (diferentes tamaños/costos)
- Alta concurrencia (muchas goroutines)
- Necesitas control fino sobre eviction
- Datos estructurados (no solo bytes)

---

## 66.6 - Redis Clients (go-redis, redigo)

### 66.6.1 - go-redis vs redigo Comparison

```
CRITERIO           go-redis           redigo

API                Modern, contextual  Tradicional
Connection Pool    Builtin, automático Manual
Streaming          ✅ Soportado        ❌ No
Context            ✅ Completo         ⚠️ Básico
Errores            Type assertion      error simple
Mantenimiento      ✅ Activo           ⚠️ Mantenido
Performance        ⭐⭐⭐⭐⭐         ⭐⭐
Tamaño             Más grandes         Más ligero
Documentación      ✅ Excelente        ✅ Buena
```

### 66.6.2 - go-redis Setup

```bash
go get github.com/redis/go-redis/v9
```

### 66.6.3 - Connection Pooling

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    
    // Conexión simple
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    
    // Pool automático (default: 10 conexiones)
    rdb = redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        PoolSize:     20,           // conexiones concurrentes
        MinIdleConns: 5,            // mínimas inactivas
        MaxRetries:   3,
    })
    
    // Test conexión
    pong, err := rdb.Ping(ctx).Result()
    fmt.Println("Ping:", pong, err)
    
    // Cluster (automático failover)
    clusterRdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "localhost:7000",
            "localhost:7001",
            "localhost:7002",
        },
    })
    
    defer clusterRdb.Close()
}
```

### 66.6.4 - Basic Operations

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer rdb.Close()
    
    // SET y GET
    err := rdb.Set(ctx, "user:1001", "Alice", 24*time.Hour).Err()
    if err != nil {
        panic(err)
    }
    
    val, err := rdb.Get(ctx, "user:1001").Result()
    fmt.Println("Usuario:", val)
    
    // DEL
    rdb.Del(ctx, "user:1001")
    
    // INCR (atomático)
    rdb.Set(ctx, "counter", "0", 0)
    count, _ := rdb.Incr(ctx, "counter").Result()
    fmt.Println("Counter:", count) // 1
    
    // LPUSH / RPOP (listas)
    rdb.RPush(ctx, "queue", "tarea1", "tarea2", "tarea3")
    item, _ := rdb.LPop(ctx, "queue").Result()
    fmt.Println("Item:", item)
    
    // HSET / HGET (hashes)
    rdb.HSet(ctx, "user:profile", "nombre", "Alice", "edad", 30)
    name, _ := rdb.HGet(ctx, "user:profile", "nombre").Result()
    fmt.Println("Nombre:", name)
    
    // SADD / SMEMBERS (sets)
    rdb.SAdd(ctx, "tags", "golang", "redis", "cache")
    members, _ := rdb.SMembers(ctx, "tags").Result()
    fmt.Println("Tags:", members)
    
    // ZADD / ZRANGE (sorted sets)
    rdb.ZAdd(ctx, "leaderboard",
        redis.Z{Score: 100, Member: "Alice"},
        redis.Z{Score: 80, Member: "Bob"},
    )
    top, _ := rdb.ZRevRange(ctx, "leaderboard", 0, 10).Result()
    fmt.Println("Top:", top)
}
```

### 66.6.5 - Pipelining

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer rdb.Close()
    
    // SIN pipelining: 100 commands = 100 round-trips (100ms @ 1ms/trip)
    // CON pipelining: 100 commands = 1 round-trip (~1ms)
    
    // Crear pipeline
    pipe := rdb.Pipeline()
    
    // Encolar comandos
    for i := 1; i <= 100; i++ {
        pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, 0)
    }
    
    // Ejecutar todo de una
    cmds, err := pipe.Exec(ctx)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Ejecutados %d comandos\n", len(cmds))
    
    // Pipelining con TxPipeline (transaccional)
    txPipe := rdb.TxPipeline()
    txPipe.Set(ctx, "account:alice", "100", 0)
    txPipe.Incr(ctx, "account:alice")  // Debe ser atómico
    
    _, err = txPipe.Exec(ctx)
    if err != nil {
        fmt.Println("Transacción fallida:", err)
    }
}
```

### 66.6.6 - Pub/Sub

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer rdb.Close()
    
    // SUSCRIPTOR
    go func() {
        pubsub := rdb.Subscribe(ctx, "notifications", "alerts")
        defer pubsub.Close()
        
        ch := pubsub.Channel()
        
        for msg := range ch {
            fmt.Printf("[%s] %s\n", msg.Channel, msg.Payload)
        }
    }()
    
    // PUBLICADOR
    for i := 1; i <= 5; i++ {
        rdb.Publish(ctx, "notifications", fmt.Sprintf("Notificación %d", i))
    }
}
```

### 66.6.7 - Cluster Support

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    
    // Cluster con 3 nodos
    clusterRdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "node1:6379",
            "node2:6379",
            "node3:6379",
        },
        ReadOnly: false, // escribe en primary
    })
    defer clusterRdb.Close()
    
    // Operaciones normales (automáticamente distribuidas)
    err := clusterRdb.Set(ctx, "user:100", "Alice", 0).Err()
    val, err := clusterRdb.Get(ctx, "user:100").Result()
    fmt.Println("Usuario:", val)
    
    // Pipeline en cluster
    pipe := clusterRdb.Pipeline()
    for i := 1; i <= 10; i++ {
        pipe.Set(ctx, fmt.Sprintf("k:%d", i), i, 0)
    }
    pipe.Exec(ctx)
    
    // Información del cluster
    info, _ := clusterRdb.ClusterInfo(ctx).Result()
    fmt.Println("Cluster info:", info)
}
```

---

## 66.7 - Caching Patterns

### 66.7.1 - Cache-Aside Pattern

```
            Application
                 |
         ┌───────┼───────┐
         |               |
         ▼               ▼
      Cache           Database
    (Redis)           (PostgreSQL)
    
    1. GET data → Try cache first
    2. If hit    → return data
    3. If miss   → query DB
    4. Store in cache
    5. Return data
```

```go
package cache

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "time"
)

type User struct {
    ID   int
    Name string
}

func GetUserCacheAside(ctx context.Context, userID int) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", userID)
    
    // 1. Intentar desde caché
    if cached, err := getFromRedis(ctx, cacheKey); err == nil {
        var user User
        json.Unmarshal(cached, &user)
        fmt.Println("✅ Hit desde caché")
        return &user, nil
    }
    
    // 2. Miss - consultar BD
    fmt.Println("❌ Miss, consultando BD...")
    user, err := getUserFromDB(userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Cachear resultado
    data, _ := json.Marshal(user)
    setInRedis(ctx, cacheKey, data, 24*time.Hour)
    
    return user, nil
}

// Helpers simulados
func getFromRedis(ctx context.Context, key string) ([]byte, error) {
    return nil, errors.New("miss")
}

func setInRedis(ctx context.Context, key string, value []byte, ttl time.Duration) {}

func getUserFromDB(id int) (*User, error) {
    return &User{ID: id, Name: "Alice"}, nil
}
```

**Ventajas:**
 Simple de implementar
 Flexible (datos pueden no ser cacheados)

**Desventajas:**
 Responsabilidad del app (invalidación)
 Complejidad en aplicación

### 66.7.2 - Write-Through Pattern

```
            Application
                 |
         ┌───────┴────────┐
         |                |
         ▼                ▼
       Cache          Database
       (OK)           (OK)
    
    Síncronos ambos
```

```go
func SetUserWriteThrough(ctx context.Context, user *User) error {
    // 1. Escribir en caché PRIMERO
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    data, _ := json.Marshal(user)
    
    if err := setInRedis(ctx, cacheKey, data, 24*time.Hour); err != nil {
        return fmt.Errorf("cache write failed: %w", err)
    }
    
    // 2. Escribir en BD
    if err := saveUserToDB(user); err != nil {
        // ROLLBACK: eliminar de caché si BD falló
        deleteFromRedis(ctx, cacheKey)
        return fmt.Errorf("db write failed: %w", err)
    }
    
    fmt.Println("✅ Datos consistentes (cache + db)")
    return nil
}
```

**Ventajas:**
 Datos siempre consistentes
 Sin complejidad de invalidación

**Desventajas:**
 Latencia (espera ambos)
 Si caché falla, falla todo

### 66.7.3 - Write-Behind Pattern

```
            Application
                 |
                 ▼
               Cache (retorna inmediato)
                 |
         ┌───────┘
         |
    [ASYNC]
         |
         ▼
           Database (delayed)
```

```go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

type WriteBackBuffer struct {
    mu      sync.Mutex
    queue   map[string]interface{}
    ticker  *time.Ticker
    done    chan struct{}
}

func NewWriteBackBuffer(flushInterval time.Duration) *WriteBackBuffer {
    wb := &WriteBackBuffer{
        queue:  make(map[string]interface{}),
        ticker: time.NewTicker(flushInterval),
        done:   make(chan struct{}),
    }
    
    // Goroutine que flush periódicamente
    go wb.flushPeriodically()
    return wb
}

func (wb *WriteBackBuffer) Set(ctx context.Context, key string, value interface{}) error {
    // 1. Escribir en caché (retorna inmediato)
    wb.mu.Lock()
    wb.queue[key] = value
    wb.mu.Unlock()
    
    fmt.Printf("✅ Escrito en buffer: %s\n", key)
    return nil
}

func (wb *WriteBackBuffer) flushPeriodically() {
    for range wb.ticker.C {
        wb.flush()
    }
}

func (wb *WriteBackBuffer) flush() {
    wb.mu.Lock()
    defer wb.mu.Unlock()
    
    if len(wb.queue) == 0 {
        return
    }
    
    fmt.Printf("💾 Flushing %d items a BD...\n", len(wb.queue))
    
    for key, value := range wb.queue {
        data, _ := json.Marshal(value)
        // Simular escritura a BD
        fmt.Printf("  - %s → BD\n", key)
        delete(wb.queue, key)
    }
}

func (wb *WriteBackBuffer) Close() {
    wb.ticker.Stop()
    wb.flush()
    close(wb.done)
}
```

**Ventajas:**
 Latencia ultra-baja (caché inmediato)
 Escrituras en batch (eficiente)

**Desventajas:**
 Pérdida de datos si crash antes de flush
 Inconsistencia temporal

### 66.7.4 - Refresh-Ahead Pattern

```go
/*
En lugar de esperar a que expire TTL:
- Detectar uso frecuente
- Refrescar proactivamente ANTES de expiración
- Reducir misses
*/

package cache

import (
    "context"
    "fmt"
    "time"
)

type RefreshAheadCache struct {
    accessCounts map[string]int
    lastRefresh  map[string]time.Time
    threshold    int // accesos para hacer refresh
}

func (rac *RefreshAheadCache) Get(ctx context.Context, key string) (interface{}, error) {
    // Incrementar contador
    rac.accessCounts[key]++
    
    // Si fue accedido > threshold veces, refrescar proactivamente
    if rac.accessCounts[key] > rac.threshold {
        if time.Since(rac.lastRefresh[key]) > 10*time.Second {
            go func() {
                fmt.Printf("[REFRESH-AHEAD] Actualizando: %s\n", key)
                // Refrescar desde BD
                // setInCache(key, loadFromDB(key))
                rac.lastRefresh[key] = time.Now()
            }()
        }
    }
    
    // Devolver valor actual del caché
    return getFromCache(key)
}

// Helpers
func getFromCache(key string) (interface{}, error)   { return nil, nil }
```

### 66.7.5 - Time-Based vs Event-Based Invalidation

**Time-Based (TTL Simple):**
```go
// Expiración fija
rdb.Set(ctx, "config", "value", 5*time.Minute)
// Siempre expira en 5 min, aunque no cambie en BD
```

**Event-Based (Invalidación activa):**
```go
package cache

type EventInvalidator struct {
    pubsub chan string
}

// Cuando BD cambia, publicar evento
func (ei *EventInvalidator) OnUserUpdated(userID int) {
    ei.pubsub <- fmt.Sprintf("user:invalidate:%d", userID)
}

// Listener limpia caché
func (ei *EventInvalidator) Listen(ctx context.Context, rdb *redis.Client) {
    for event := range ei.pubsub {
        fmt.Printf("🗑️  Invalidando: %s\n", event)
        rdb.Del(ctx, event)
    }
}
```

**Hybrid (lo mejor):**
```go
// TTL + Event-based
rdb.Set(ctx, "user:1", "Alice", 1*time.Hour)  // TTL fallback

// Si user:1 se modifica en BD, invalidar inmediatamente
func UpdateUser(userID int) {
    // ... update en BD
    rdb.Del(ctx, fmt.Sprintf("user:%d", userID))  // Invalidar ahora
    // Próximo acceso recargará desde BD
}
```

---

## 66.8 - Distributed Caching Solutions

### 66.8.1 - Memcached vs Redis

```
CRITERIO             Memcached        Redis

Tipos de datos       Solo strings     ✅ 5+ tipos
TTL/Expiración       ✅ Sí             ✅ Sí
Persistencia         ❌ No             ✅ Opcional
Replicación          ❌ No             ✅ Sí (master-slave)
Cluster              ❌ Client-side    ✅ Nativo (v3+)
Pub/Sub              ❌ No             ✅ Sí
Lua Scripts          ❌ No             ✅ Sí (atómico)
Performance          ⭐⭐⭐⭐⭐        ⭐⭐⭐⭐
Memoria              Optimizada       Flexible
```

### 66.8.2 - Memcached Usage

```go
package main

import (
    "fmt"
    "github.com/bradfitz/gomemcache/memcache"
)

func main() {
    // Conectar
    mc, _ := memcache.New("localhost:11211")
    
    // SET
    item := &memcache.Item{
        Key:        "user:1",
        Value:      []byte("Alice"),
        Expiration: 3600, // segundos
    }
    mc.Set(item)
    
    // GET
    it, err := mc.Get("user:1")
    if err == nil {
        fmt.Println("User:", string(it.Value))
    }
    
    // DEL
    mc.Delete("user:1")
    
    // ADD (solo si no existe)
    mc.Add(item)
    
    // INCR
    mc.Increment("counter", 1)
}
```

### 66.8.3 - Redis Sentinel (High Availability)

```

     REDIS SENTINEL ARCHITECTURE              │

                                              │
  Sentinel 1  Sentinel 2  Sentinel 3         │
       |           |           |              │
              │clear
                   |                          │
         ┌─────────┼─────────┐               │
         |         |         |                │
         ▼         ▼         ▼               │
      Master   Slave 1   Slave 2            │
                                             │
  Monitor: si Master cae                     │
  └─ Promover Slave 1 a Master              │
  └─ Reconfigurar Slaves                    │
  └─ Notificar a aplicación                 │
                                             │

```

```go
package main

import (
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    // go-redis soporta Sentinel automáticamente
    rdb := redis.NewFailoverClient(&redis.FailoverOptions{
        MasterName: "mymaster",
        SentinelAddrs: []string{
            "sentinel1:26379",
            "sentinel2:26379",
            "sentinel3:26379",
        },
    })
    defer rdb.Close()
    
    // Funciona como cliente normal
    // Failover automático si master cae
}
```

### 66.8.4 - Redis Cluster (Horizontal Scaling)

```

    REDIS CLUSTER (Sharding)             │

                                         │
  Hash Slot Distribution:                │
  0-5460:   Nodo A                       │
  5461-10922: Nodo B                     │
  10923-16383: Nodo C                    │
                                         │
  GET key →                              │
  1. CRC16(key) % 16384 = slot           │
  2. Encontrar nodo que tiene slot       │
  3. Consultar ese nodo                  │
                                         │
  Rebalancing: agregar nodo D            │
  └─ Redistribuir slots                  │
  └─ Replicar datos                      │
  └─ Clientes se adaptan                 │
                                         │

```

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    
    // Cliente de cluster
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "node1:6379",
            "node2:6379",
            "node3:6379",
        },
    })
    defer rdb.Close()
    
    // API idéntica a cliente normal
    rdb.Set(ctx, "key1", "value1", 0)
    val, _ := rdb.Get(ctx, "key1").Result()
    fmt.Println(val)
    
    // Info del cluster
    info, _ := rdb.ClusterInfo(ctx).Result()
    fmt.Println("Cluster state:", info)
    
    // Conocer slots
    slots, _ := rdb.ClusterSlots(ctx).Result()
    for _, slot := range slots {
        fmt.Printf("Slot %d-%d: %v\n", slot.Start, slot.End, slot.Nodes)
    }
}
```

### 66.8.5 - DragonflyDB (Redis Alternative)

```
DragonflyDB vs Redis:
- Mejor performance (1.5-3x)
- Compatible con Redis protocol
- Mejor manejo de memory
- RESP3 protocol
- Menos overhead de replicación

Usar si:
 Performance crítica
 Quieres alternativa a Redis
 Compatible con go-redis
```

---

## 66.9 - HTTP Caching & CDN

### 66.9.1 - HTTP Cache Headers

```go
package main

import (
    "net/http"
    "time"
)

func cacheableHandler(w http.ResponseWriter, r *http.Request) {
    // Público (caché puede guardarlo)
    w.Header().Set("Cache-Control", "public, max-age=3600")
    
    // Privado (navegador sí, CDN no)
    w.Header().Set("Cache-Control", "private, max-age=1800")
    
    // No cacheable
    w.Header().Set("Cache-Control", "no-cache, no-store")
    
    // Validar antes de usar
    w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Contenido"))
}

func etagHandler(w http.ResponseWriter, r *http.Request) {
    content := "Mi contenido"
    etag := generateETag(content)
    
    w.Header().Set("ETag", `"`+etag+`"`)
    
    // Si cliente envía If-None-Match
    if r.Header.Get("If-None-Match") == `"`+etag+`"` {
        w.WriteHeader(http.StatusNotModified)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(content))
}

func lastModifiedHandler(w http.ResponseWriter, r *http.Request) {
    lastMod := time.Now().UTC()
    w.Header().Set("Last-Modified", lastMod.Format(http.TimeFormat))
    
    // Si cliente envía If-Modified-Since
    ifModSince, _ := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))
    if lastMod.Before(ifModSince) {
        w.WriteHeader(http.StatusNotModified)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Contenido"))
}

func generateETag(content string) string {
    // Usar hash del contenido
    return hash(content)
}

func hash(s string) string {
    // Implementar
    return ""
}
```

### 66.9.2 - ETag & Conditional Requests

```
CLIENT                                     SERVER

Request:
GET /api/data


Response:
                                    ← 200 OK
                         ETag: "abc123"
                         Content: {...}

[Cliente guarda contenido + ETag]

Próxima request:
GET /api/data
If-None-Match: "abc123"


[Servidor compara ETag]

Si aún es válido:
                                    ← 304 Not Modified
                                       (sin cuerpo)

[Cliente usa versión en caché]
```

### 66.9.3 - Cache-Control Directives

```go
package main

import "net/http"

func demonstrateDirectives(w http.ResponseWriter, r *http.Request) {
    // max-age: segundos que puede servirse desde caché
    // 1 hora
    w.Header().Set("Cache-Control", "max-age=3600")
    
    // s-maxage: para caché compartido (CDN)
    // 1 día en CDN, 1 hora en navegador
    w.Header().Set("Cache-Control", "max-age=3600, s-maxage=86400")
    
    // public: cualquiera puede cachear
    w.Header().Set("Cache-Control", "public, max-age=3600")
    
    // private: solo navegador (no CDN ni proxies)
    w.Header().Set("Cache-Control", "private, max-age=1800")
    
    // no-cache: validar con servidor antes de usar
    w.Header().Set("Cache-Control", "no-cache")
    
    // no-store: nunca cachear (sensible)
    w.Header().Set("Cache-Control", "no-store")
    
    // must-revalidate: si expiró, DEBE revalidar
    w.Header().Set("Cache-Control", "max-age=3600, must-revalidate")
    
    // Combinado: frecuente en APIs
    // Público, 1 hora en CDN, 5 min en navegador, revalidar si expira
    w.Header().Set("Cache-Control", 
        "public, max-age=300, s-maxage=3600, must-revalidate")
}
```

### 66.9.4 - Cloudflare Integration

```go
package main

import (
    "fmt"
    "net/http"
)

func cfCacheHandler(w http.ResponseWriter, r *http.Request) {
    // Headers específicos de Cloudflare
    
    // Cache todo durante 1 hora (bypass auth)
    w.Header().Set("Cache-Control", "public, max-age=3600")
    
    // Cache status
    w.Header().Set("X-Frame-Options", "SAMEORIGIN")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    
    // Cloudflare Page Rules (desde CF, no desde código)
    // Pero podemos indicar intención con headers
    w.Header().Set("X-Cache-Status", "HIT") // informativo
    
    // Cache purge token (si necesitas limpiar caché)
    // POST /api/purge con headers especiales
    
    fmt.Fprint(w, "Contenido cacheado por Cloudflare")
}

// Para limpiar caché en Cloudflare (requiere API)
func purgeCloudflareCache(url string) error {
    // Requiere API key de Cloudflare
    // Example:
    req, _ := http.NewRequest("POST", 
        "https://api.cloudflare.com/client/v4/zones/{zone}/purge_cache", nil)
    
    req.Header.Set("X-Auth-Email", "your@email.com")
    req.Header.Set("X-Auth-Key", "your-api-key")
    
    client := &http.Client{}
    resp, _ := client.Do(req)
    fmt.Println("Purge response:", resp.Status)
    return nil
}
```

### 66.9.5 - Edge Caching Strategy

```

  EDGE CACHING MULTI-TIER ARCHITECTURE      │

                                            │
  Browser Cache (user machine)              │
  ├─ TTL: 5-30 min                          │
  ├─ Control: max-age                       │
  └─ Validación: ETag                       │
          ↓ Miss                            │
                                            │
  CDN Edge Cache (Cloudflare, Fastly, etc) │
  ├─ TTL: 1-24 horas                        │
  ├─ Control: s-maxage                      │
  ├─ Ubicación: global (pop nodes)          │
  └─ Hit rate: 95%+                         │
          ↓ Miss                            │
                                            │
  Origin Server Cache (Redis)               │
  ├─ TTL: flexible                          │
  ├─ Control: aplicación                    │
  └─ Hot data cerca de BD                   │
          ↓ Miss                            │
                                            │
  Database                                  │
  └─ Última consulta (lenta)                │
                                            │
  IMPACTO:                                  │
  ├─ Hit en browser: <10ms                  │
  ├─ Hit en CDN: 50-200ms                   │
  ├─ Hit en origen: 10-100ms                │
  └─ Hit en BD: 100-5000ms                  │
                                            │

```

---

## 66.10 - Monitoring & Debugging

### 66.10.1 - Cache Metrics

```go
package cache

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type CacheMetrics struct {
    hits        int64
    misses      int64
    sets        int64
    deletes     int64
    evictions   int64
    errors      int64
    totalSize   int64
    mu          sync.RWMutex
}

func (cm *CacheMetrics) RecordHit() {
    atomic.AddInt64(&cm.hits, 1)
}

func (cm *CacheMetrics) RecordMiss() {
    atomic.AddInt64(&cm.misses, 1)
}

func (cm *CacheMetrics) GetStats() map[string]interface{} {
    h := atomic.LoadInt64(&cm.hits)
    m := atomic.LoadInt64(&cm.misses)
    total := h + m
    
    hitRate := float64(0)
    if total > 0 {
        hitRate = float64(h) / float64(total) * 100
    }
    
    return map[string]interface{}{
        "hits":       h,
        "misses":     m,
        "hit_rate":   fmt.Sprintf("%.2f%%", hitRate),
        "total_ops":  total,
        "sets":       atomic.LoadInt64(&cm.sets),
        "deletes":    atomic.LoadInt64(&cm.deletes),
        "evictions":  atomic.LoadInt64(&cm.evictions),
        "errors":     atomic.LoadInt64(&cm.errors),
        "size_bytes": atomic.LoadInt64(&cm.totalSize),
    }
}

// Endpoint HTTP para monitoreo
func (cm *CacheMetrics) HTTPHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        stats := cm.GetStats()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(stats)
    })
}
```

### 66.10.2 - Hit Rate Monitoring

```go
package main

import (
    "fmt"
    "time"
)

type HitRateMonitor struct {
    window      time.Duration
    hitCounts   map[time.Time]int64
    missCounts  map[time.Time]int64
}

func (hrm *HitRateMonitor) RecordOperation(isHit bool) {
    now := time.Now().Truncate(time.Second)
    
    if isHit {
        hrm.hitCounts[now]++
    } else {
        hrm.missCounts[now]++
    }
}

func (hrm *HitRateMonitor) GetAverageHitRate() float64 {
    total := int64(0)
    hits := int64(0)
    
    for _, h := range hrm.hitCounts {
        hits += h
        total += h
    }
    
    for _, m := range hrm.missCounts {
        total += m
    }
    
    if total == 0 {
        return 0
    }
    
    return float64(hits) / float64(total) * 100
}

func (hrm *HitRateMonitor) PrintStats() {
    fmt.Printf("Hit rate: %.2f%%\n", hrm.GetAverageHitRate())
}
```

### 66.10.3 - Memory Profiling

```bash
# Generar perfil de memoria
go tool pprof http://localhost:6060/debug/pprof/heap

# En tu aplicación Go:
# Necesitas importar _ "net/http/pprof"
# y exponer en runtime

# Luego en REPL:
# (pprof) top     # Top consumers
# (pprof) list main.cacheFunction
# (pprof) web     # Generar gráfico
```

```go
package main

import (
    _ "net/http/pprof"
    "net/http"
)

func init() {
    // Exponer pprof en puerto 6060
    go http.ListenAndServe(":6060", nil)
}

func main() {
    // Tu aplicación aquí
    // http://localhost:6060/debug/pprof/
}
```

### 66.10.4 - Debug Memory Leaks

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func detectMemoryLeak(interval time.Duration, threshold uint64) {
    var m runtime.MemStats
    var lastMem uint64
    
    for {
        time.Sleep(interval)
        
        runtime.ReadMemStats(&m)
        
        if m.Alloc > lastMem+threshold {
            fmt.Printf("Memory spike: %v MB → %v MB\n",
                lastMem/1024/1024, m.Alloc/1024/1024)
            
            // Generar snapshot
            runtime.GC()
            fmt.Printf("After GC: %v MB\n", m.Alloc/1024/1024)
        }
        
        lastMem = m.Alloc
    }
}

func PrintCacheStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Memoria caché estimada:\n")
    fmt.Printf("  Alloc: %v MB\n", m.Alloc/1024/1024)
    fmt.Printf("  TotalAlloc: %v MB\n", m.TotalAlloc/1024/1024)
    fmt.Printf("  Sys: %v MB\n", m.Sys/1024/1024)
    fmt.Printf("  Goroutines: %d\n", runtime.NumGoroutine())
}
```

---

## 66.11 - Selection Guide & Case Studies

### 66.11.1 - Decision Matrix

```
LIBRERÍA       SIZE  LATENCY  DISTRIBUTED  TTL  CONSISTENCY  COMPLEXITY

go-cache       ✅    ✅✅    ❌           ✅   Eventual     ⭐
BigCache       ✅✅  ✅      ❌           ❌   Eventual     ⭐⭐
Ristretto      ✅    ✅✅    ❌           ❌   Eventual     ⭐⭐⭐
Groupcache     ✅✅✅ ✅     ✅           ❌   Eventual     ⭐⭐⭐
go-redis       ✅✅✅ ✅✅   ✅✅         ✅   Configurable ⭐⭐
Memcached      ✅✅  ✅✅    ⚠️(manual)   ✅   Eventual     ⭐⭐

DECISION:
1. ¿Distribuido?
   - NO → in-memory (go-cache, BigCache, Ristretto)
   - SÍ → Groupcache, Redis

2. ¿Tamaño datos?
   - < 100 MB → go-cache
   - 100MB-1GB → BigCache
   - > 1GB → Redis + Cluster

3. ¿TTL importante?
   - SÍ → go-cache, Redis
   - NO → Ristretto, BigCache

4. ¿Concurrencia alta?
   - SÍ → Ristretto, Redis
   - NO → go-cache

5. ¿Performance crítica?
   - < 1µs → BigCache, Ristretto
   - 1-100µs → go-cache
   - 1-100ms → Redis
```

### 66.11.2 - Case Study 1: E-Commerce Platform

```
REQUISITOS:
- 1M+ usuarios activos
- Caché de: productos, carritos, sesiones
- Alta disponibilidad
- Consistencia importante

ARQUITECTURA:


        CLIENTE (Navegador)           │

                     │
        ┌─────────
        ▼                         ▼
   ┌─────────────┐        ┌─────────────┐
   │  Servidor 1 │        │  Servidor 2 │
   │ (Go + cache)│        │ (Go + cache)│
   └──────┬──────┘        └──────┬──────┘
          │                      │
          └──────────┬───────────┘
                     ▼
         ┌─────────────────────────┐
         │   Redis Cluster (3 nodos)
         │   - Productos (s-maxage)
         │   - Carritos (expire)
         │   - Sesiones (expire)
         └────────┬──
                  │
        ┌─────────┴────────────┐
        ▼                      ▼
    Database             Search (Elasticsearch)

ESTRATEGIA:
- Productos: go-cache local + Redis distribuido
  * Cache-Aside pattern
  * TTL: 1 hora
  * Invalidación: event-based (cuando admin actualiza)

- Carritos: Redis + Write-Behind
  * Caché 5min, flush a BD cada 10 seg
  * Persistencia en Redis (backup)

- Sesiones: Redis + Sentinel
  * Failover automático
  * TTL: 24 horas
  * Replicación a 2 slaves
```

```go
package ecommerce

import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
)

type Product struct {
    ID    int
    Name  string
    Price float64
}

type CatalogService struct {
    localCache  map[int]*Product // go-cache
    redis       *redis.Client
}

func (cs *CatalogService) GetProduct(ctx context.Context, id int) (*Product, error) {
    key := fmt.Sprintf("product:%d", id)
    
    // L1: Local cache
    if p, ok := cs.localCache[id]; ok {
        return p, nil
    }
    
    // L2: Redis (distribuido)
    val, err := cs.redis.Get(ctx, key).Result()
    if err == nil {
        var p Product
        json.Unmarshal([]byte(val), &p)
        cs.localCache[id] = &p
        return &p, nil
    }
    
    // L3: Database (miss)
    p, err := cs.getFromDB(id)
    if err != nil {
        return nil, err
    }
    
    // Cachear en ambos niveles
    data, _ := json.Marshal(p)
    cs.redis.Set(ctx, key, data, 1*time.Hour)
    cs.localCache[id] = p
    
    return p, nil
}
```

### 66.11.3 - Case Study 2: Real-Time Analytics

```
REQUISITOS:
- 100K eventos/seg
- Agregaciones rápidas
- Ventanas de tiempo (últimas 24h, 7d, 30d)
- Precisión importante

ARQUITECTURA:


 Event Stream │
 (100K/sec)   │

       │
       ▼
   ┌─────────────┐
   │   BigCache  │ (últimas 10 minutos)
   │ (buckets)   │
   └──────┬──────┘
          
          ▼
   ┌──────────────────┐
   │ Redis Aggregates │ (últimas 24h)
   │ - MIN           │
   │ - MAX           │
   │ - COUNT         │
   └──────┬───────────┘
          │
          ▼
   ┌──────────────┐
   │ TimeSeries DB│ (histórico)
   │ (InfluxDB)   │
   └──────────────┘

ESTRATEGIA:
- BigCache para hot data (eventos recientes)
- Redis sorted sets para agregaciones
- TimeSeries para histórico
- Cleanup automático por edad
```

```go
package analytics

import (
    "fmt"
    "time"
    "github.com/allegro/bigcache/v3"
    "github.com/redis/go-redis/v9"
)

type EventAggregator struct {
    bigcache *bigcache.BigCache
    redis    *redis.Client
}

func (ea *EventAggregator) RecordEvent(event map[string]interface{}) error {
    // 1. Cachear en BigCache (últimos 10 min)
    timestamp := time.Now()
    key := fmt.Sprintf("event:%d:%d", timestamp.Unix(), event["id"])
    data, _ := json.Marshal(event)
    ea.bigcache.Set(key, data)
    
    // 2. Agregar en Redis (últimas 24h)
    ea.redis.ZAdd(context.Background(),
        "events:24h",
        &redis.Z{
            Score:  float64(timestamp.Unix()),
            Member: string(data),
        },
    )
    
    return nil
}

func (ea *EventAggregator) GetStats() map[string]interface{} {
    ctx := context.Background()
    
    // Obtener count from Redis
    count, _ := ea.redis.ZCard(ctx, "events:24h").Result()
    
    return map[string]interface{}{
        "events_24h": count,
    }
}
```

### 66.11.4 - Migration Strategies

**Migración: go-cache → Redis**

```go
package migration

import (
    "context"
    "github.com/patrickmn/go-cache"
    "github.com/redis/go-redis/v9"
    "json"
)

func MigrateToRedis(ctx context.Context, 
    oldCache *cache.Cache,
    newRedis *redis.Client) {
    
    // Paso 1: Leer todo de go-cache
    items := oldCache.Items()
    
    // Paso 2: Escribir a Redis
    pipe := newRedis.Pipeline()
    
    for key, item := range items {
        // Convertir valor a JSON
        data, _ := json.Marshal(item.Object)
        
        // Calcular TTL restante
        ttl := time.Until(time.Unix(0, item.Expiration))
        if ttl < 0 {
            ttl = 1 * time.Second
        }
        
        pipe.Set(ctx, key, data, ttl)
    }
    
    pipe.Exec(ctx)
    
    // Paso 3: Dual-write nuevo código
    // Leer: redis (con fallback a go-cache)
    // Escribir: a ambos
    
    // Paso 4: Deprecated go-cache cuando hit-rate < 5%
}
```

### 66.11.5 - Best Practices ✅

```go
/*
 BEST PRACTICES
*/

// 1. Usar siempre time.Duration para TTL
cache.Set(key, value, 5*time.Minute)  // ✅ Correcto

// 2. Manejo de errores en operaciones distribuidas
val, err := redis.Get(ctx, key).Result()
if err == redis.Nil {
    // Cache miss - normal
    val, err = getFromDB()
} else if err != nil {
    // Error de network - considerar fallback
    log.Warn("Redis error, usando BD:", err)
}

// 3. Evitar thundering herd
// Si muchas requests llegan al mismo tiempo y miss:
// → Todos van a BD → sobrecarga

// Solución: lock mutuo o probabilístico
func GetWithLocking(ctx context.Context, key string) {
    // Obtener lock distribuido
    lock := redis.NewScript(`
        if redis.call('SET', KEYS[1], ARGV[1], 'NX', 'EX', 30) then
            return 1
        else
            return 0
        end
    `)
    
    ok, _ := lock.Run(ctx, redis.Client, []string{key + ":lock"}, uuid()).Result()
    if ok == 1 {
        // Solo este proceso consulta BD
        val := getFromDB()
        redis.Set(ctx, key, val, 1*time.Hour)
    } else {
        // Otros esperan
        time.Sleep(100 * time.Millisecond)
        redis.Get(ctx, key)
    }
}

// 4. Invalidación explícita
// Nunca confiar solo en TTL para datos importantes
func UpdateUser(userID int, data *User) error {
    if err := saveUserToDB(data); err != nil {
        return err
    }
    
    // Invalidar caché INMEDIATAMENTE
    redis.Del(ctx, fmt.Sprintf("user:%d", userID))
    
    return nil
}

// 5. Monitoreo del hit-rate
// < 50%: probablemente datos no cacheables
// > 95%: puede ser problema de freshness
```

### 66.11.6 - Anti-Patterns ❌

```go
/*
 ANTI-PATTERNS A EVITAR
*/

// ❌ 1. Cache sin límite de memoria
cache := make(map[string]interface{})
for {
    cache[generateKey()] = largeData  // OOM
}

// ✅ Usar caché con eviction policy
// BigCache, Ristretto, Redis con maxmemory

// ❌ 2. No validar datos en caché
cachedUser, _ := cache.Get("user:1")
user := cachedUser.(User)  // Panic si tipo incorrecto!

// ✅ Validar tipos
if val, ok := cache.Get("user:1"); ok {
    user, ok := val.(User)
    if !ok {
        // Tipo incorrecto, invalidar
        cache.Delete("user:1")
        user, _ = getFromDB()
    }
}

// ❌ 3. TTL muy largo para datos que cambian
redis.Set(ctx, "user:1", userData, 30*time.Day)  // ¿Qué si usuario se actualiza?

// ✅ Balance: TTL + invalidación event-based
redis.Set(ctx, "user:1", userData, 1*time.Hour)
// + invalidar en UpdateUser()

// ❌ 4. Comprimir sin medir overhead
data, _ := json.Marshal(user)
compressed := gzip.Compress(data)  // ¿Vale la pena?

// ✅ Medir
// Si compression overhead (CPU) > latencia savedahuguería
// → No comprimir

// ❌ 5. Caché centralizado sin replicación
redis := NewSingleRedis(":6379")  // Punto único de falla

// ✅ Usar replicación/cluster
redis := NewClusterRedis([3 nodes])  // Failover automático
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: go-cache Simple Cache

**Objetivo:** Implementar un caché de sesiones simple con TTL.

```go
package main

import (
    "fmt"
    "time"
    "github.com/patrickmn/go-cache"
)

func main() {
    // TODO: Implementar un session store con go-cache
    // 1. Crear cache con cleanup cada 5 minutos
    // 2. CreateSession: genera ID único, guarda datos
    // 3. GetSession: retorna datos si existen
    // 4. DeleteSession: limpia sesión
    // 5. ListActiveSessions: cuenta sesiones activas
    
    // Requisitos:
    // - Sesiones expiran en 30 minutos
    // - Mostrar hit rate de accesos
    // - Validar antes de cada operación
    
    // Resultado esperado:
    // CreateSession() → "session_123abc"
    // GetSession("session_123abc") → {user: "alice", rol: "admin"}
    // ActiveSessions: 5
    // Hit rate: 85%
}

// SOLUCIÓN:
type SessionStore struct {
    cache *cache.Cache
    hits  int64
    misses int64
}

func (ss *SessionStore) CreateSession(data map[string]interface{}) string {
    sessionID := generateID()
    ss.cache.Set(sessionID, data, 30*time.Minute)
    return sessionID
}

func (ss *SessionStore) GetSession(sessionID string) (map[string]interface{}, bool) {
    val, found := ss.cache.Get(sessionID)
    if found {
        ss.hits++
    } else {
        ss.misses++
    }
    return val.(map[string]interface{}), found
}
```

### Ejercicio 2: BigCache High-Performance

**Objetivo:** Cachear millones de items sin pausas GC.

```go
package main

import (
    "fmt"
    "github.com/allegro/bigcache/v3"
)

func main() {
    // TODO: Crear caché de 100M items
    // 1. Medir escritura: 1M items/segundo
    // 2. Medir lectura: 5M lookups
    // 3. Mostrar hit rate
    // 4. Monitorear memoria
    // 5. Comparar con go-cache
    
    // Resultado esperado:
    // BigCache: 100M+ ops/sec, < 100ms GC pause
    // go-cache: 10M ops/sec, > 500ms GC pause
}
```

### Ejercicio 3: Redis Distributed Cache

**Objetivo:** Implementar caché distribuido con Redis cluster.

```go
package main

import (
    "context"
    "github.com/redis/go-redis/v9"
)

func main() {
    // TODO: Crear sistema de caché distribuido
    // 1. Conectar a Redis cluster
    // 2. Implementar Get/Set/Delete con error handling
    // 3. Configurar pipelining para batch ops
    // 4. Implementar rate limiting por usuario
    // 5. Monitorear hit rate
    
    // Resultado esperado:
    // 1000 usuarios, 10K ops/sec
    // 90% hit rate
    // Failover automático si nodo cae
}
```

### Ejercicio 4: Cache with Compression

**Objetivo:** Cachear datos comprimidos para ahorrar memoria.

```go
package main

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
)

func main() {
    // TODO: Implementar caché con compresión
    // 1. Serializar structs a JSON
    // 2. Comprimir con gzip
    // 3. Almacenar en Redis
    // 4. Medir overhead (CPU time vs memory saved)
    // 5. Decidir si compresión vale la pena
    
    // Estructura de prueba:
    type LargeData struct {
        ID       int
        Name     string
        Email    string
        Profile  string // 10KB de texto
        Metadata map[string]interface{}
    }
    
    // Resultado esperado:
    // Sin compresión: 10KB por item, 100ms/op
    // Con compresión: 2KB por item, 110ms/op
    // Vale la pena si tenemos >1M items
}
```

### Ejercicio 5: Production System (Todos los Patrones)

**Objetivo:** Diseñar un sistema de caché multi-nivel para aplicación real.

```go
package main

import (
    "context"
    "github.com/patrickmn/go-cache"
    "github.com/allegro/bigcache/v3"
    "github.com/redis/go-redis/v9"
)

type ProductCatalogSystem struct {
    l1Cache   *cache.Cache      // go-cache para hot products
    l2Cache   *bigcache.BigCache // todos los productos
    l3Cache   *redis.Client      // distribuido
}

func main() {
    // TODO: Implementar sistema de 3 niveles
    // 1. L1 (Local/Hot): top 100 productos, TTL 5 min
    // 2. L2 (BigCache): todos los productos, TTL 1 hora
    // 3. L3 (Redis): distribuido, TTL 24 horas
    //
    // Requisitos:
    // - Cache-Aside pattern para L1
    // - Write-Through para updates
    // - Invalidación event-based
    // - Monitoreo de hit rate en cada nivel
    // - Fallback automático si nivel falla
    //
    // Pruebas:
    // - 10K req/sec
    // - 95% hit rate esperado
    // - < 50ms latencia P99
    // - Failover si Redis cae
}
```

---

## RESUMEN COMPARATIVO

```

         MATRIZ DE DECISIÓN (RESUMEN EJECUTIVO)              │

                                                             │
  Datos < 100MB                                             │
  └─ Baja latencia crítica      → BigCache, Ristretto       │
  └─ TTl + Simple               → go-cache                  │
                                                             │
  Datos 100MB - 1GB                                         │
  └─ Concentrado en server único → BigCache                │
  └─ Distribuido (multi-server)  → Groupcache             │
                                                             │
  Datos > 1GB                                               │
  └─ Todos a Redis cluster                                  │
                                                             │
  HIGH CONCURRENCY                                          │
  └─ >100K ops/sec → Ristretto, BigCache                   │
  └─ >1M ops/sec   → Redis cluster                         │
                                                             │
  DISTRIBUTED                                               │
  └─ Mismo datacenter   → Groupcache                       │
  └─ Multi-datacenter   → Redis sentinel/cluster            │
                                                             │

```

---

## CONCLUSIONES

### Recomendaciones Finales

1. **Comienza simple**: go-cache para prototipado
2. **Escala vertical**: BigCache cuando necesites performance
3. **Escala horizontal**: Redis cuando requieras distribución
4. **Monitorea siempre**: hit rate, latencia, memoria
5. **Invalida proactivamente**: no solo confíes en TTL

### Próximos Pasos

- Estudiar Redis persistence (RDB, AOF)
- Configurar replicación y failover
- Implementar backup strategies
- Profiling y tuning específico

**FIN DEL CAPÍTULO 66**

---

**Próximos capítulos:** Deployment, Performance Tuning, Advanced Patterns

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/66-caching-libraries/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/66-caching-libraries):

```bash
cd examples/66-caching-libraries
go run .
```
