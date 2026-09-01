# Capítulo 52: Caching - Sistemas de caché en producción

## Índice General
1. [52.1 - Introducción a Caching](#521---introducción-a-caching)
2. [52.2 - Memoria Local (In-Memory)](#522---memoria-local-in-memory)
3. [52.3 - Problemas de Cache](#523---problemas-de-cache-cache-invalidation)
4. [52.4 - Redis](#524---redis)
5. [52.5 - Redis Avanzado](#525---redis-avanzado)
6. [52.6 - Cache Patterns](#526---cache-patterns)
7. [52.7 - Distributed Caching](#527---distributed-caching)
8. [52.8 - Memcached](#528---memcached)
9. [52.9 - Caching de HTTP](#529---caching-de-http)
10. [52.10 - Testing y Monitoring](#5210---testing-y-monitoring)
11. [52.11 - Buenas Prácticas y Case Studies](#5211---buenas-prácticas-y-case-studies)

---

## 52.1 - INTRODUCCIÓN A CACHING

### 52.1.1 ¿Por Qué Cachear?

El caching es una de las técnicas más fundamentales para escalar sistemas. La premisa es simple: **guardar el resultado de operaciones costosas para evitar repetirlas**.

```
Latencia típica de acceso (2024):
L1 cache (CPU):         ~4 ciclos       (1 ns)
L2 cache (CPU):         ~10 ciclos      (3 ns)
L3 cache (CPU):         ~40 ciclos      (12 ns)
Memoria RAM:            ~250 ciclos     (100 ns)
SSD:                    ~5.000 ciclos   (1 μs)
Base de datos local:    ~1 ms
Base de datos remota:   ~10 ms
API de red:             ~100 ms
```

Un request sin cache que tarda 100ms tarda 1ms con cache en memoria. **100x más rápido**.

### 52.1.2 Niveles de Caché

```

                         Usuario (Navegador)                  │
                    [L0: Browser Cache]                        │

                           │ HTTP con ETag, Cache-Control
clear
                      CDN (Cloudflare, Akamai)                 │
                    [L1: CDN Cache]                            │

                           │ TCP/HTTP
clear
                    Aplicación (Go)                            │
  ┌─────────────────────────────────────────────────────────┐ │
  │       [L2: Application In-Memory Cache]                 │ │
  │  sync.Map / go-cache / local LRU                        │ │
  └──────────────────────┬──────────────────────────────────┘ │
                         │                                     │
  ┌──────────────────────▼──────────────────────────────────┐ │
  │       [L3: Distributed Cache (Redis/Memcached)]         │ │
  │  Shared across instances                                │ │
 │  └──────────────────────┬─────────────────────────────────
clear
                           │ TCP/Socket
clear
         Base de Datos (PostgreSQL, MySQL)                     │
                    [L4: DB Cache]                             │
         (Buffer pool, Query cache deprecated)                │

                           │
clear
           Storage (Disco físico, SSD)                         │
                    [L5: Disk]                                 │

```

- **L0 (Browser)**: Caché del navegador, controlada por headers HTTP
- **L1 (CDN)**: Contenido estático distribuido globalmente
- **L2 (Application)**: En memoria de la aplicación, rápido pero volátil
- **L3 (Redis/Memcached)**: Caché distribuida compartida entre instancias
- **L4 (Database)**: Buffer pool, índices, query cache
- **L5 (Disk)**: Almacenamiento persistente

### 52.1.3 Trade-offs: Complejidad vs Performance

| Aspecto | Caché Local | Redis | Memcached | HTTP Cache |
|---------|-------------|-------|-----------|-----------|
| Latencia | <1ms | 1-5ms | 1-5ms | 100-500ms |
| Consistencia | ⚠️ Eventual | ✅ Mejor | ⚠️ Eventual | ⚠️ Eventual |
| Escalabilidad | ❌ Limitada | ✅ Excelente | ✅ Excelente | ✅ Excelente |
| Persistencia | ❌ No | ✅ Sí | ❌ No | ⚠️ Parcial |
| Complejidad | ✅ Baja | ⚠️ Media | ✅ Baja | ⚠️ Media |
| Memory overhead | ⚠️ Aumenta | ✅ Centralizado | ✅ Centralizado | N/A |

### 52.1.4 Casos de Uso y Antipatterns

**✅ CASOS DE USO APROPIADOS:**
- Consultas frecuentes a base de datos (user profiles, settings)
- Cálculos computacionales costosos (reporte analíticos)
- Datos semi-estáticos (catálogos de productos)
- Sessions de usuario
- Rate limiting / throttling

**❌ ANTIPATTERNS COMUNES:**

```go
// ❌ MALO: No hay TTL, memory leak garantizado
func BadCache() {
    cache := make(map[string]interface{})
    cache["user:1000"] = userData // Nunca se borra
    // Después de millones de requests, OOM
}

// ✅ BIEN: TTL con evicción automática
func GoodCache() {
    cache := NewCacheWithTTL(time.Hour)
    cache.Set("user:1000", userData, time.Hour)
}

// ❌ MALO: Cache solo en la aplicación, inconsistente en cluster
func InconsistentCache() {
    // App A y App B tienen cachés diferentes
    // User actualiza en App A, App B devuelve dato viejo
}

// ✅ BIEN: Cache centralizado (Redis) para consistencia
func ConsistentCache() {
    // Todas las instancias consultan Redis
    // Un único source of truth
}

// ❌ MALO: Cache sin manejo de fallos
func CrashOnCacheFailure() {
    if err := redis.Get(key); err != nil {
        panic("Cache failed!") // ❌ Crash total
    }
}

// ✅ BIEN: Graceful degradation
func GracefulFailure() {
    val, err := redis.Get(key)
    if err != nil {
        log.Printf("Cache miss, using DB: %v", err)
        return db.Query(key) // Fallback a DB
    }
    return val
}
```

---

## 52.2 - MEMORIA LOCAL (IN-MEMORY)

### 52.2.1 sync.Map para Concurrencia

`sync.Map` es una solución de Go para mapas thread-safe sin locks explícitos.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// SimpleInMemoryCache usando sync.Map
type SimpleCache struct {
    data sync.Map
}

func (c *SimpleCache) Set(key string, value interface{}) {
    c.data.Store(key, value)
}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
    return c.data.Load(key)
}

func (c *SimpleCache) Delete(key string) {
    c.data.Delete(key)
}

// Conteo de operaciones
func (c *SimpleCache) Count() int {
    count := 0
    c.data.Range(func(key, value interface{}) bool {
        count++
        return true
    })
    return count
}

func ExampleSyncMap() {
    cache := &SimpleCache{}
    
    // Set/Get concurrente
    go func() {
        for i := 0; i < 1000; i++ {
            key := fmt.Sprintf("key:%d", i)
            cache.Set(key, fmt.Sprintf("value:%d", i))
        }
    }()
    
    go func() {
        for i := 0; i < 1000; i++ {
            key := fmt.Sprintf("key:%d", i)
            val, ok := cache.Get(key)
            if ok {
                _ = val.(string)
            }
        }
    }()
    
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Cache contiene %d items\n", cache.Count())
}

// Benchmark: sync.Map vs map con RWMutex
func BenchmarkSyncMapVsMutex() {
    // sync.Map
    syncMap := sync.Map{}
    start := time.Now()
    for i := 0; i < 100000; i++ {
        syncMap.Store(fmt.Sprintf("key:%d", i), i)
    }
    fmt.Printf("sync.Map: %v\n", time.Since(start))
    
    // map + RWMutex
    var mu sync.RWMutex
    regularMap := make(map[string]int)
    start = time.Now()
    for i := 0; i < 100000; i++ {
        mu.Lock()
        regularMap[fmt.Sprintf("key:%d", i)] = i
        mu.Unlock()
    }
    fmt.Printf("map + RWMutex: %v\n", time.Since(start))
    
    // Resultado: sync.Map es ~2-3x más rápido para reads + writes mixed
}
```

### 52.2.2 TTL y Evicción Manual

Para caché con expiración, necesitamos gestionar TTL manualmente en Go.

```go
package main

import (
    "sync"
    "time"
)

type CacheItem struct {
    Value      interface{}
    ExpiresAt  time.Time
}

type TTLCache struct {
    mu    sync.RWMutex
    items map[string]*CacheItem
    done  chan struct{}
}

func NewTTLCache() *TTLCache {
    c := &TTLCache{
        items: make(map[string]*CacheItem),
        done:  make(chan struct{}),
    }
    go c.cleanupExpired()
    return c
}

func (c *TTLCache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.items[key] = &CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.items[key]
    if !exists {
        return nil, false
    }
    
    // Check if expired
    if time.Now().After(item.ExpiresAt) {
        return nil, false
    }
    
    return item.Value, true
}

func (c *TTLCache) cleanupExpired() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            c.mu.Lock()
            now := time.Now()
            for key, item := range c.items {
                if now.After(item.ExpiresAt) {
                    delete(c.items, key)
                }
            }
            c.mu.Unlock()
        case <-c.done:
            return
        }
    }
}

func (c *TTLCache) Close() {
    close(c.done)
}

// Ejemplo de uso
func ExampleTTLCache() {
    cache := NewTTLCache()
    defer cache.Close()
    
    // Guardar con 5 segundos de TTL
    cache.Set("user:1", map[string]string{"name": "Alice"}, 5*time.Second)
    
    // Inmediatamente disponible
    val, ok := cache.Get("user:1")
    fmt.Printf("Inmediato: %v, %v\n", val, ok) // true
    
    // Después de 6 segundos, expirado
    time.Sleep(6 * time.Second)
    val, ok = cache.Get("user:1")
    fmt.Printf("Después de TTL: %v, %v\n", val, ok) // false
}
```

### 52.2.3 LRU (Least Recently Used)

Limitar caché por número máximo de items, removiendo los menos usados.

```go
package main

import (
    "container/list"
    "sync"
)

type LRUCache struct {
    capacity int
    cache    map[string]*list.Element
    list     *list.List
    mu       sync.RWMutex
}

type cacheNode struct {
    key   string
    value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*list.Element),
        list:     list.New(),
    }
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    elem, exists := c.cache[key]
    if !exists {
        return nil, false
    }
    
    // Mover a frente (más recientemente usado)
    c.list.MoveToFront(elem)
    return elem.Value.(*cacheNode).value, true
}

func (c *LRUCache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Si ya existe, actualizar y mover a frente
    if elem, exists := c.cache[key]; exists {
        c.list.MoveToFront(elem)
        elem.Value.(*cacheNode).value = value
        return
    }
    
    // Nuevo item
    elem := c.list.PushFront(&cacheNode{key: key, value: value})
    c.cache[key] = elem
    
    // Evictar LRU si se excede capacidad
    if c.list.Len() > c.capacity {
        evicted := c.list.Remove(c.list.Back())
        delete(c.cache, evicted.(*cacheNode).key)
    }
}

func ExampleLRU() {
    lru := NewLRUCache(3)
    
    lru.Set("a", 1)
    lru.Set("b", 2)
    lru.Set("c", 3)
    fmt.Println("Cache:", "a, b, c")
    
    lru.Set("d", 4) // Evicta 'a' (least recently used)
    
    _, exists := lru.Get("a")
    fmt.Printf("¿Existe 'a'?: %v\n", exists) // false
    
    val, exists := lru.Get("d")
    fmt.Printf("Valor de 'd': %v\n", val) // 4
}
```

### 52.2.4 Librería: github.com/patrickmn/go-cache

Librería popular para caché simple con TTL.

```go
package main

import (
    "fmt"
    "time"
    "github.com/patrickmn/go-cache"
)

func ExampleGoCache() {
    // Crear caché con default expiration 5 min, cleanup cada 10 min
    c := cache.New(5*time.Minute, 10*time.Minute)
    
    // Set
    c.Set("userID:123", map[string]interface{}{
        "name": "Alice",
        "role": "admin",
    }, cache.DefaultExpiration)
    
    // Set con TTL específico
    c.Set("session:abc123", "session_data", 30*time.Minute)
    
    // Get
    user, found := c.Get("userID:123")
    if found {
        fmt.Printf("Usuario: %v\n", user)
    }
    
    // Get y eliminar
    session, found := c.GetWithExpiration("session:abc123")
    if found {
        fmt.Printf("Session: %v\n", session)
    }
    
    // Iterar
    c.Items() // Retorna todos los items
    
    // Limpiar
    c.Flush()
}

// Uso en una aplicación real
type UserRepository struct {
    cache *cache.Cache
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        cache: cache.New(10*time.Minute, 15*time.Minute),
    }
}

func (ur *UserRepository) GetUser(userID string) (map[string]interface{}, error) {
    // Intentar obtener del caché
    if cached, found := ur.cache.Get("user:" + userID); found {
        return cached.(map[string]interface{}), nil
    }
    
    // No en caché, consultar DB
    user := map[string]interface{}{
        "id":    userID,
        "name":  "Alice",
        "email": "alice@example.com",
    }
    
    // Guardar en caché
    ur.cache.Set("user:"+userID, user, cache.DefaultExpiration)
    
    return user, nil
}

func (ur *UserRepository) InvalidateUser(userID string) {
    ur.cache.Delete("user:" + userID)
}
```

---

## 52.3 - PROBLEMAS DE CACHE (CACHE INVALIDATION)

### 52.3.1 TTL vs Event-Based

**TTL (Time-To-Live)**: Expiración automática tras N segundos.
**Event-Based**: Invalidación cuando ocurre un evento (ej: actualización de datos).

```go
package main

import (
    "sync"
    "time"
)

// TTL-based: Simple pero datos pueden estar desfasados
type TTLCacheStrategy struct {
    cache map[string]interface{}
    ttl   time.Duration
    mu    sync.RWMutex
}

func (c *TTLCacheStrategy) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = value
    
    // Auto-expire después de TTL
    go func() {
        time.Sleep(c.ttl)
        c.mu.Lock()
        delete(c.cache, key)
        c.mu.Unlock()
    }()
}

// Event-based: Más preciso pero requiere invalidación manual
type EventBasedCache struct {
    cache          map[string]interface{}
    mu             sync.RWMutex
    invalidateFunc func(key string)
}

func (c *EventBasedCache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = value
}

func (c *EventBasedCache) OnUserUpdate(userID string) {
    // Cuando el usuario se actualiza, invalida su caché
    key := "user:" + userID
    c.mu.Lock()
    delete(c.cache, key)
    c.mu.Unlock()
    
    // Notificar listeners
    if c.invalidateFunc != nil {
        c.invalidateFunc(key)
    }
}

// Hybrid: TTL + Event-based lo mejor de ambos mundos
type HybridCache struct {
    cache      map[string]interface{}
    expiresAt  map[string]time.Time
    mu         sync.RWMutex
    maxTTL     time.Duration
}

func (c *HybridCache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = value
    c.expiresAt[key] = time.Now().Add(c.maxTTL)
}

func (c *HybridCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    val, exists := c.cache[key]
    if !exists {
        return nil, false
    }
    
    if time.Now().After(c.expiresAt[key]) {
        return nil, false
    }
    
    return val, true
}

func (c *HybridCache) InvalidateOnEvent(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.cache, key)
    delete(c.expiresAt, key)
}
```

### 52.3.2 Cache Stampede (Thundering Herd)

Problema: Múltiples goroutines intentan recargar el mismo dato caducado simultáneamente.

```
Timeline del Cache Stampede:

T1: Cache expire para key="user:1"
    ├─ Request A: ¿Está en caché? NO
    ├─ Request B: ¿Está en caché? NO
    ├─ Request C: ¿Está en caché? NO
    ├─ Request D: ¿Está en caché? NO
    └─ Todas lanzan query a DB SIMULTÁNEAMENTE
       
DB recibe 4 queries idénticas → servidor se ahoga
```

**Solución 1: Locking**

```go
package main

import (
    "sync"
    "time"
)

type CacheWithLock struct {
    data  map[string]interface{}
    locks map[string]*sync.Mutex
    mu    sync.RWMutex
}

func (c *CacheWithLock) GetOrLoad(key string, loader func() (interface{}, error)) (interface{}, error) {
    // Obtener o crear lock para esta key
    c.mu.Lock()
    lock, exists := c.locks[key]
    if !exists {
        lock = &sync.Mutex{}
        c.locks[key] = lock
    }
    c.mu.Unlock()
    
    // Solo el primer goroutine carga datos, otros esperan
    lock.Lock()
    defer lock.Unlock()
    
    // Verificar nuevamente si alguien ya cargó
    c.mu.RLock()
    if val, exists := c.data[key]; exists {
        c.mu.RUnlock()
        return val, nil
    }
    c.mu.RUnlock()
    
    // Cargar datos
    val, err := loader()
    if err == nil {
        c.mu.Lock()
        c.data[key] = val
        c.mu.Unlock()
    }
    
    return val, err
}
```

**Solución 2: Probabilistic Cache Expiration**

```go
package main

import (
    "math/rand"
    "time"
)

type ProbabilisticCache struct {
    data      map[string]interface{}
    expiresAt map[string]time.Time
}

func (c *ProbabilisticCache) Set(key string, value interface{}, baseTTL time.Duration) {
    // En lugar de TTL fijo, usar rango probabilístico
    // TTL = baseTTL * random(0.8, 1.0)
    variance := 0.8 + rand.Float64()*0.2
    actualTTL := time.Duration(float64(baseTTL) * variance)
    
    c.data[key] = value
    c.expiresAt[key] = time.Now().Add(actualTTL)
}

// Resultado: Los items expiran en diferentes momentos,
// evitando que todos expiren simultáneamente
```

### 52.3.3 Cache Penetration

Problema: Búsquedas de datos inexistentes golpean constantemente la BD.

```
Escenario: Búsqueda de userID="99999" (no existe)

Request 1: Cache miss → DB query → No encontrado → Cache miss
Request 2: Cache miss → DB query → No encontrado → Cache miss
Request 3: Cache miss → DB query → No encontrado → Cache miss
... 1000 queries a BD por un dato que no existe
```

**Solución: Negative Caching**

```go
package main

import (
    "time"
)

type NegativeCacheEntry struct {
    Found bool
    Value interface{}
}

func (c *TTLCache) GetOrDefault(key string, loader func() (interface{}, error)) (interface{}, error) {
    // Intentar caché
    if cached, found := c.Get(key); found {
        if entry, ok := cached.(*NegativeCacheEntry); ok {
            if !entry.Found {
                return nil, nil // Dato no existe (cached)
            }
            return entry.Value, nil
        }
    }
    
    // Cargar de fuente
    val, err := loader()
    
    // Cachear resultado, aunque no exista
    if err == nil {
        c.Set(key, &NegativeCacheEntry{
            Found: val != nil,
            Value: val,
        }, 5*time.Minute)
    }
    
    return val, err
}
```

### 52.3.4 Cache Avalanche

Problema: Muchas claves expiran simultáneamente, causando pico de carga en BD.

```
Timeline de Cache Avalanche:

T=0:  Todos los items con TTL=60min fueron cacheados
T=60: TODOS expiran SIMULTÁNEAMENTE
      ├─ Millones de requests golpean BD
      ├─ BD se satura
      ├─ Respuesta lenta
      └─ Más timeouts → más reintentos → avalancha

Solución: Diversificar TTL
```

**Código:**

```go
package main

import (
    "math/rand"
    "time"
)

type AvalancheSafeCache struct {
    data      map[string]interface{}
    expiresAt map[string]time.Time
}

func (c *AvalancheSafeCache) Set(key string, value interface{}, baseTTL time.Duration) {
    // Agregar varianza aleatoria para evitar expiración simultánea
    variance := 0.9 + rand.Float64()*0.2 // 0.9 a 1.1
    actualTTL := time.Duration(float64(baseTTL) * variance)
    
    c.data[key] = value
    c.expiresAt[key] = time.Now().Add(actualTTL)
}

// Alternativa: Refresh ahead
func (c *AvalancheSafeCache) RefreshAhead(key string, value interface{}, ttl time.Duration) {
    actualTTL := ttl + time.Duration(rand.Intn(int(ttl/10))) // +0-10% de variance
    c.data[key] = value
    c.expiresAt[key] = time.Now().Add(actualTTL)
}
```

### 52.3.5 Diagrama ASCII de Sincronización

```
ESCENARIO: Cache Update en Sistema Distribuido

Instancia A              Redis Cache         Database
    │                        │                  │
    ├──────────── GET key ──────────────────────┤
    │                        │ NO ENCONTRADO     │
    ├──────────────────────────────────────────── QUERY ──────────────┤
    │                        │                  │ CALCULAR RESULTADO │
    │◄───────────────────────────────────────────── (2MB JSON)       │
    ├──────────── SET key (2MB) ───────
    │                        ├─ GUARDAR        │
    │                        ├─ TTL=3600s      │
    │
Instancia B (después de 100ms)
    │                        │                  │
    ├──────────── GET key ──────────────────────┤
 ENCONTRADO (2MB) ───────────┤    │◄────
    │ (Cache hit, sin golpear DB)               │
    │
Instancia C (después de 5000ms)
    │                        │                  │
    ├──────────── GET key ──────────────────────┤
    │                        │ EXPIRADO         │
    ├──────────────────────────────────────────── QUERY ──────────────┤
    │                        │                  │ RECALCULAR         │
```

---

## 52.4 - REDIS

### 52.4.1 Arquitectura Redis

Redis es un data store in-memory con múltiples data types:

```
Redis Internals:


   Client Connection (TCP)        │

               │
clear
   Redis Protocol Parser (RESP)   │  RESP = Redis Serialization Protocol
  Protocolo binario optimizado
   Command Queue (Single Thread)  │  Single-threaded pero rápido

               │
clear
   Data Structures                    │
  ├─ Strings (valores simple)         │
  ├─ Lists (linked list)              │
  ├─ Sets (colecciones)               │
  ├─ Sorted Sets (con score)          │
  ├─ Hashes (mapas)                   │
  ├─ HyperLogLog (cardinality)        │
  ├─ Streams (evento logs)            │
  └─ Geospatial (índices geo)         │

               │
clear
   Persistence                        │
  ├─ RDB (snapshots)                  │
  ├─ AOF (Append Only File)           │
  └─ Hybrid (RDB + AOF)               │

```

### 52.4.2 Setup con go-redis

```bash
go get github.com/redis/go-redis/v9
```

**Configuración básica:**

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func main() {
    // Conexión simple
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    defer client.Close()
    
    ctx := context.Background()
    
    // Ping
    pong, err := client.Ping(ctx).Result()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Redis: %s\n", pong)
    
    // Cluster
    clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "redis-node-1:6379",
            "redis-node-2:6379",
            "redis-node-3:6379",
        },
    })
    defer clusterClient.Close()
    
    // Sentinel (High Availability)
    sentinelClient := redis.NewSentinelClient(&redis.SentinelOptions{
        Addrs: []string{"sentinel1:26379", "sentinel2:26379"},
    })
    defer sentinelClient.Close()
    
    masterClient := sentinelClient.GetMasterClient(ctx, "mymaster")
    defer masterClient.Close()
}
```

### 52.4.3 Basic Operations

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

func BasicOperations() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // STRING
    client.Set(ctx, "username", "alice", 0)
    val, _ := client.Get(ctx, "username").Result()
    fmt.Printf("GET username: %s\n", val)
    
    // MGET/MSET (múltiples claves)
    client.MSet(ctx, "user:1:name", "Alice", "user:2:name", "Bob")
    vals, _ := client.MGet(ctx, "user:1:name", "user:2:name").Result()
    fmt.Printf("MGET: %v\n", vals)
    
    // INCR (contador)
    client.Set(ctx, "views:article:1", "0", 0)
    client.Incr(ctx, "views:article:1")
    client.Incr(ctx, "views:article:1")
    count, _ := client.Get(ctx, "views:article:1").Result()
    fmt.Printf("Views: %s\n", count)
    
    // APPEND
    client.Set(ctx, "message", "Hello", 0)
    client.Append(ctx, "message", " World")
    msg, _ := client.Get(ctx, "message").Result()
    fmt.Printf("Message: %s\n", msg)
    
    // LIST
    client.RPush(ctx, "queue", "task1", "task2", "task3")
    client.LPush(ctx, "queue", "urgent_task")
    len, _ := client.LLen(ctx, "queue").Result()
    fmt.Printf("Queue length: %d\n", len)
    
    // POP
    first, _ := client.LPop(ctx, "queue").Result()
    fmt.Printf("First task: %s\n", first)
    
    // HASH
    client.HSet(ctx, "user:100", map[string]interface{}{
        "name":  "Alice",
        "email": "alice@example.com",
        "role":  "admin",
    })
    
    email, _ := client.HGet(ctx, "user:100", "email").Result()
    fmt.Printf("User email: %s\n", email)
    
    all, _ := client.HGetAll(ctx, "user:100").Result()
    fmt.Printf("All fields: %v\n", all)
    
    // SET
    client.SAdd(ctx, "tags", "golang", "redis", "caching", "golang")
    members, _ := client.SMembers(ctx, "tags").Result()
    fmt.Printf("Tags: %v\n", members) // golang, redis, caching (sin duplicados)
    
    // ZSET (Sorted Set)
    client.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "Alice"})
    client.ZAdd(ctx, "leaderboard", redis.Z{Score: 85, Member: "Bob"})
    client.ZAdd(ctx, "leaderboard", redis.Z{Score: 92, Member: "Charlie"})
    
    // Top 3
    top, _ := client.ZRevRangeByScoreWithScores(ctx, "leaderboard", 
        &redis.ZRangeBy{Min: "-inf", Max: "+inf", Count: 3}).Result()
    for _, entry := range top {
        fmt.Printf("%v: %g\n", entry.Member, entry.Score)
    }
}

// TTL y Expiraciones
func TTLOperations() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // SET con expiración
    client.Set(ctx, "session:abc123", "user_data", 30*time.Minute)
    
    // Ver TTL restante
    ttl, _ := client.TTL(ctx, "session:abc123").Result()
    fmt.Printf("TTL: %v\n", ttl)
    
    // EXPIRE: establecer expiración en clave existente
    client.Set(ctx, "temp_data", "value", 0)
    client.Expire(ctx, "temp_data", 5*time.Minute)
    
    // EXPIREAT: expirar en fecha específica
    client.ExpireAt(ctx, "temp_data", time.Now().Add(1*time.Hour))
    
    // PERSIST: remover expiración
    client.Persist(ctx, "temp_data")
}

// Operaciones Atómicas
func AtomicOperations() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // GETSET: obtener valor anterior mientras se actualiza
    client.Set(ctx, "counter", "0", 0)
    prev, _ := client.GetSet(ctx, "counter", "1").Result()
    fmt.Printf("Previous value: %s\n", prev)
    
    // GETEX: obtener con opción de actualizar TTL
    client.Set(ctx, "key", "value", 0)
    val, _ := client.GetEx(ctx, "key", &redis.GetExOptions{
        EX: 30 * time.Second,
    }).Result()
    fmt.Printf("Value: %s\n", val)
}
```

### 52.4.4 Pub/Sub (Publicador/Suscriptor)

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

func PubSubExample() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // Suscriptor
    go func() {
        pubsub := client.Subscribe(ctx, "news", "alerts")
        defer pubsub.Close()
        
        ch := pubsub.Channel()
        for msg := range ch {
            fmt.Printf("[%s] %s\n", msg.Channel, msg.Payload)
        }
    }()
    
    // Publicador
    time.Sleep(1 * time.Second)
    client.Publish(ctx, "news", "Breaking news!")
    client.Publish(ctx, "alerts", "System alert!")
    
    time.Sleep(2 * time.Second)
}

// Pattern Subscription
func PatternSubscription() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // Suscribirse a patrón
    pubsub := client.PSubscribe(ctx, "user:*", "system:*")
    defer pubsub.Close()
    
    ch := pubsub.Channel()
    
    go func() {
        for msg := range ch {
            fmt.Printf("Pattern match [%s]: %s\n", msg.Channel, msg.Payload)
        }
    }()
    
    // Publicar mensajes
    client.Publish(ctx, "user:login", "user123 logged in")
    client.Publish(ctx, "system:status", "System healthy")
    client.Publish(ctx, "user:logout", "user456 logged out")
    
    time.Sleep(1 * time.Second)
}
```

---

## 52.5 - REDIS AVANZADO

### 52.5.1 Pipelines (Batch Commands)

Enviar múltiples comandos en una sola conexión (reduce latencia).

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func PipelineExample() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // Sin Pipeline (N roundtrips)
    start := time.Now()
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("key:%d", i)
        client.Set(ctx, key, fmt.Sprintf("value:%d", i), 0)
    }
    fmt.Printf("Sin Pipeline: %v\n", time.Since(start))
    
    // Con Pipeline (1 roundtrip)
    start = time.Now()
    pipe := client.Pipeline()
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("key:%d", i)
        pipe.Set(ctx, key, fmt.Sprintf("value:%d", i), 0)
    }
    _, err := pipe.Exec(ctx)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("Con Pipeline: %v\n", time.Since(start))
    // Típicamente 10-50x más rápido
}

// Conditional Pipeline
func ConditionalPipeline() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // Actualizar múltiples usuarios solo si existen
    pipe := client.Pipeline()
    userIDs := []string{"user:1", "user:2", "user:3"}
    
    for _, userID := range userIDs {
        pipe.Exists(ctx, userID)
    }
    
    results, _ := pipe.Exec(ctx)
    for i, result := range results {
        if val, ok := result.(*redis.IntCmd); ok && val.Val() > 0 {
            userID := userIDs[i]
            pipe.HSet(ctx, userID, "last_seen", time.Now().Unix())
        }
    }
    
    pipe.Exec(ctx)
}
```

### 52.5.2 Transactions (MULTI/EXEC)

Garantizar atomicidad en múltiples comandos.

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func TransactionExample() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // Transferencia de dinero: Alice → Bob (atomicidad)
    err := client.Watch(ctx, func(tx *redis.Tx) error {
        // Leer saldos
        aliceBalance, _ := tx.Get(ctx, "user:alice:balance").Int64()
        bobBalance, _ := tx.Get(ctx, "user:bob:balance").Int64()
        
        transferAmount := int64(100)
        
        // Verificar que Alice tenga suficientes fondos
        if aliceBalance < transferAmount {
            return fmt.Errorf("Insufficient funds")
        }
        
        // Ejecutar en transaction
        _, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
            pipe.Decrby(ctx, "user:alice:balance", transferAmount)
            pipe.Incrby(ctx, "user:bob:balance", transferAmount)
            return nil
        })
        return err
    }, "user:alice:balance")
    
    if err != nil {
        fmt.Printf("Transaction failed: %v\n", err)
    }
}

// Contador con WATCH (optimistic locking)
func CounterWithOptimisticLock() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    client.Set(ctx, "counter", "0", 0)
    
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        err := client.Watch(ctx, func(tx *redis.Tx) error {
            val, _ := tx.Get(ctx, "counter").Int64()
            
            _, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
                pipe.Set(ctx, "counter", val+1, 0)
                return nil
            })
            
            return err
        }, "counter")
        
        if err == redis.Nil {
            // Collision, reintentar
            continue
        }
        
        if err == nil {
            fmt.Printf("Counter incremented successfully\n")
            break
        }
    }
}
```

### 52.5.3 Lua Scripting

Ejecutar scripts Lua en Redis de forma atómica.

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func LuaScriptExample() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // Script: Incrementar si es menor que limit
    script := redis.NewScript(`
        local current = redis.call('get', KEYS[1])
        local limit = tonumber(ARGV[1])
        
        if not current then
            current = 0
        else
            current = tonumber(current)
        end
        
        if current < limit then
            redis.call('set', KEYS[1], current + 1)
            return current + 1
        else
            return -1
        end
    `)
    
    // Ejecutar script
    result, err := script.Run(ctx, client, []string{"counter"}, 10).Result()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("Result: %v\n", result)
}

// Rate Limiter usando Lua
func RateLimiterWithLua() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    // Script: Sliding window rate limiter
    rateLimiterScript := redis.NewScript(`
        local key = KEYS[1]
        local limit = tonumber(ARGV[1])
        local window = tonumber(ARGV[2])
        local now = tonumber(ARGV[3])
        
        local clearBefore = now - window
        redis.call('zremrangebyscore', key, 0, clearBefore)
        
        local current = redis.call('zcard', key)
        
        if current < limit then
            redis.call('zadd', key, now, now)
            redis.call('expire', key, window)
            return 1  -- Permitido
        else
            return 0  -- Rechazado
        end
    `)
    
    // Test rate limiting: 5 requests per 60 segundos
    userID := "user:123"
    limit := 5
    window := 60
    
    for i := 0; i < 10; i++ {
        result, _ := rateLimiterScript.Run(ctx, client, 
            []string{userID}, limit, window, time.Now().Unix()).Int64()
        
        if result == 1 {
            fmt.Printf("Request %d: Permitido\n", i+1)
        } else {
            fmt.Printf("Request %d: Rechazado (rate limit)\n", i+1)
        }
    }
}

// Distributed Lock con Lua
func DistributedLockWithLua() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer client.Close()
    ctx := context.Background()
    
    lockKey := "resource:lock"
    lockValue := "unique_token_123"
    
    // Adquirir lock
    acquireScript := redis.NewScript(`
        if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then
            redis.call('expire', KEYS[1], ARGV[2])
            return 1
        else
            return 0
        end
    `)
    
    result, _ := acquireScript.Run(ctx, client, 
        []string{lockKey}, lockValue, 30).Int64()
    
    if result == 1 {
        fmt.Println("Lock acquired!")
        
        // Hacer trabajo...
        
        // Liberar lock (solo si es el propietario)
        releaseScript := redis.NewScript(`
            if redis.call('get', KEYS[1]) == ARGV[1] then
                return redis.call('del', KEYS[1])
            else
                return 0  -- No es propietario
            end
        `)
        
        releaseScript.Run(ctx, client, []string{lockKey}, lockValue)
        fmt.Println("Lock released!")
    }
}
```

### 52.5.4 Cluster Mode

Redis Cluster para escalabilidad horizontal.

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func ClusterModeExample() {
    // Conectar a cluster
    clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "redis-node1:6379",
            "redis-node2:6379",
            "redis-node3:6379",
        },
        MaxRedirects: 8,
    })
    defer clusterClient.Close()
    
    ctx := context.Background()
    
    // Operaciones normales (automáticamente enrutadas)
    clusterClient.Set(ctx, "user:1000", "alice", 0)
    val, _ := clusterClient.Get(ctx, "user:1000").Result()
    fmt.Printf("User: %s\n", val)
    
    // MGET/MSET también funciona (pero requiere mismo slot)
    // Para claves en diferentes slots, usar pipeline
}

// Hash Tag para agrupar claves en mismo slot
func HashTagExample() {
    client := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{"redis-node1:6379", "redis-node2:6379"},
    })
    defer client.Close()
    
    ctx := context.Background()
    
    // Hash tag: {userId} asegura que ambas claves vayan al mismo slot
    userID := "user:123"
    userDataKey := fmt.Sprintf("{%s}:data", userID)
    userSettingsKey := fmt.Sprintf("{%s}:settings", userID)
    
    // Ahora estos MGET funcionarán (mismo slot)
    client.Set(ctx, userDataKey, "data_value", 0)
    client.Set(ctx, userSettingsKey, "settings_value", 0)
    
    vals, _ := client.MGet(ctx, userDataKey, userSettingsKey).Result()
    fmt.Printf("Values: %v\n", vals)
}

// Información del cluster
func ClusterInfo() {
    client := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{"redis-node1:6379"},
    })
    defer client.Close()
    
    ctx := context.Background()
    
    // Info del cluster
    info := client.ClusterInfo(ctx)
    fmt.Printf("Cluster state: %v\n", info)
    
    // Slots
    slots, _ := client.ClusterSlots(ctx).Result()
    for _, slot := range slots {
        fmt.Printf("Slot %d-%d: %v\n", slot.Start, slot.End, slot.Nodes)
    }
}
```

### 52.5.5 Sentinela para HA

Redis Sentinel para alta disponibilidad automática.

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func SentinelHA() {
    // Conectar a Sentinel
    sentinel := redis.NewSentinelClient(&redis.SentinelOptions{
        Addrs: []string{
            "sentinel1:26379",
            "sentinel2:26379",
            "sentinel3:26379",
        },
    })
    defer sentinel.Close()
    
    ctx := context.Background()
    
    // Obtener master
    masterAddr, err := sentinel.GetMasterAddrByName(ctx, "mymaster").Result()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Master: %s\n", masterAddr)
    
    // Conectar al master
    master := sentinel.GetMasterClient(ctx, "mymaster")
    defer master.Close()
    
    // Operaciones
    master.Set(ctx, "key", "value", 0)
    val, _ := master.Get(ctx, "key").Result()
    fmt.Printf("Value: %s\n", val)
    
    // Si master falla, Sentinel automáticamente promove un esclavo
}

// Monitoring de Sentinel
func SentinelMonitoring() {
    sentinel := redis.NewSentinelClient(&redis.SentinelOptions{
        Addrs: []string{"sentinel1:26379"},
    })
    defer sentinel.Close()
    
    ctx := context.Background()
    
    // Obtener estado del master
    masters, _ := sentinel.Masters(ctx).Result()
    for name, state := range masters {
        fmt.Printf("Master %s: %v\n", name, state)
    }
    
    // Obtener estado de esclavos
    slaves, _ := sentinel.Slaves(ctx, "mymaster").Result()
    for _, slave := range slaves {
        fmt.Printf("Slave: %v\n", slave)
    }
}
```

---

## 52.6 - CACHE PATTERNS

### 52.6.1 Cache-Aside (Lazy Loading)

El patrón más común: **aplicación es responsable de llenar el caché**.

```
Flujo:
1. Aplicación busca en caché
2. Si NO está (cache miss) → consulta BD
3. Guardar en caché
4. Devolver resultado

Ventaja: Simple, sin dependencia de caché
Desventaja: Primer acceso es lento (cold cache)
```

**Código:**

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

type UserRepository struct {
    redisClient *redis.Client
    db          *Database // Tu DB
}

func (ur *UserRepository) GetUser(ctx context.Context, userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // Paso 1: Intentar caché
    cachedJSON, err := ur.redisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        // Cache hit
        var user User
        json.Unmarshal([]byte(cachedJSON), &user)
        fmt.Printf("Cache hit: user %s\n", userID)
        return &user, nil
    }
    
    if err != redis.Nil {
        // Error de conexión a Redis, fallback a DB
        fmt.Printf("Redis error, using DB: %v\n", err)
        return ur.db.GetUser(userID)
    }
    
    // Paso 2: Cache miss, consultar BD
    fmt.Printf("Cache miss: user %s, querying DB\n", userID)
    user, err := ur.db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    // Paso 3: Guardar en caché
    userJSON, _ := json.Marshal(user)
    ur.redisClient.Set(ctx, cacheKey, userJSON, 1*time.Hour)
    
    return user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *User) error {
    // Actualizar BD
    err := ur.db.UpdateUser(user)
    if err != nil {
        return err
    }
    
    // Invalidar caché
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    ur.redisClient.Del(ctx, cacheKey)
    
    return nil
}
```

**Diagrama:**

```
REQUEST → [App] ──┬─ "Buscar en caché"
                  │
                  ├─ ¿Está en Redis? ──────────────┐
                  │     │ SÍ (hit)          ↓
                  │     └──────────────→ [RETURN] ✓ (1ms)
                  │     │ NO (miss)
     └──────────┐                  
                  │                ├─ Consultar DB (100ms)
                  │                ├─ Guardar en Redis (1ms)
                  │                └──────────────→ [RETURN] ✓
```

### 52.6.2 Write-Through

**Caché siempre actualizado**, pero requiere coordinación.

```
Flujo:
1. Escribir en caché
2. Escribir en BD (sincrónico)
3. Si BD falla, deshacer caché

Ventaja: Caché siempre consistente
Desventaja: Más lento (doble escritura)
```

**Código:**

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

func (ur *UserRepository) UpdateUserWriteThrough(ctx context.Context, user *User) error {
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    
    // Paso 1: Escribir en caché primero (rápido)
    userJSON, _ := json.Marshal(user)
    ur.redisClient.Set(ctx, cacheKey, userJSON, 1*time.Hour)
    
    // Paso 2: Escribir en BD (lento)
    err := ur.db.UpdateUser(user)
    if err != nil {
        // Revertir caché si DB falla
        fmt.Printf("DB error, reverting cache: %v\n", err)
        ur.redisClient.Del(ctx, cacheKey)
        return err
    }
    
    fmt.Printf("User %s updated (write-through)\n", user.ID)
    return nil
}
```

### 52.6.3 Write-Behind (Write-Back)

**Escritura asincrnnnica**: primero en caché, luego en BD de forma diferida.

```
Flujo:
1. Escribir en caché (inmediato)
2. Encolar actualización en BD (background job)
3. Aplicación continúa

Ventaja: Muy rápido para escrituras
Desventaja: Riesgo de inconsistencia si caché falla
```

**Código:**

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

type WriteBackCache struct {
    redis       *redis.Client
    db          *Database
    writeQueue  chan *PendingWrite
}

type PendingWrite struct {
    User      *User
    Timestamp time.Time
}

func (wbc *WriteBackCache) UpdateUserWriteBehind(ctx context.Context, user *User) error {
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    
    // Paso 1: Escribir en caché INMEDIATAMENTE
    userJSON, _ := json.Marshal(user)
    wbc.redis.Set(ctx, cacheKey, userJSON, 0)
    
    // Paso 2: Encolar para escribir en BD (background)
    wbc.writeQueue <- &PendingWrite{
        User:      user,
        Timestamp: time.Now(),
    }
    
    fmt.Printf("User %s cached, BD update queued\n", user.ID)
    return nil
}

// Background worker
func (wbc *WriteBackCache) ProcessWriteQueue() {
    for write := range wbc.writeQueue {
        retries := 0
        maxRetries := 3
        
        for retries < maxRetries {
            err := wbc.db.UpdateUser(write.User)
            if err == nil {
                fmt.Printf("User %s written to DB\n", write.User.ID)
                break
            }
            
            retries++
            time.Sleep(time.Duration(retries) * time.Second) // Exponential backoff
        }
        
        if retries == maxRetries {
            fmt.Printf("Failed to write user %s after %d retries\n", 
                write.User.ID, maxRetries)
            // Alertar administrador
        }
    }
}

func (wbc *WriteBackCache) Start() {
    go wbc.ProcessWriteQueue()
}
```

### 52.6.4 Refresh-Ahead

**Actualizar caché proactivamente** antes de que expire.

```
Flujo:
1. Detectar que caché está próximo a expirar
2. Actualizar en background (antes de expiración)
3. Usuario nunca experimenta cache miss

Ventaja: Transparente al usuario
Desventaja: Cálculo innecesario si no se accede
```

**Código:**

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

func (ur *UserRepository) GetUserWithRefreshAhead(ctx context.Context, userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // Obtener con TTL
    val, err := ur.redisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        // Cache hit
        var user User
        json.Unmarshal([]byte(val), &user)
        
        // Verificar si está próximo a expirar (< 5 min)
        ttl, _ := ur.redisClient.TTL(ctx, cacheKey).Result()
        if ttl < 5*time.Minute && ttl > 0 {
            // Actualizar en background
            go ur.refreshUserCache(ctx, userID)
        }
        
        return &user, nil
    }
    
    // Cache miss o error, cargar normalmente
    user, err := ur.db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    userJSON, _ := json.Marshal(user)
    ur.redisClient.Set(ctx, cacheKey, userJSON, 1*time.Hour)
    return user, nil
}

func (ur *UserRepository) refreshUserCache(ctx context.Context, userID string) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    user, err := ur.db.GetUser(userID)
    if err != nil {
        fmt.Printf("Error refreshing cache: %v\n", err)
        return
    }
    
    userJSON, _ := json.Marshal(user)
    ur.redisClient.Set(ctx, cacheKey, userJSON, 1*time.Hour)
    fmt.Printf("Refreshed cache for user %s\n", userID)
}
```

### 52.6.5 Diagramas de Patrones

```
CACHE-ASIDE:

 Request │

     │
     ├──→ Cache hit?
     │     YES ──→ Return cached ✓ (1ms)
     │     NO ──┐
     │          └──→ Query DB ──→ Cache ──→ Return (100ms)

WRITE-THROUGH:

 Write Op │

     │
     ├──→ Update Cache
     ├──→ Update DB (SYNC)
     └──→ Return ✓

WRITE-BEHIND:

 Write Op │

     │
     ├──→ Update Cache ✓ (INMEDIATO)
     ├──→ Queue to Background (ASYNC)
     └──→ Return
          │
          └──→ [Background: Update DB when possible]

REFRESH-AHEAD:

 Request │

     │
     ├──→ Cache hit + TTL < 5min?
     │     YES ──→ Return cached + Schedule refresh
     │     NO ──→ Query DB ──→ Cache ──→ Return
```

---

## 52.7 - DISTRIBUTED CACHING

### 52.7.1 Multi-Instance Consistency

Problema: Cuando hay múltiples instancias de aplicación, necesitamos sincronización.

```
Escenario SIN Distributed Cache:

Instance A           Instance B           Instance C
  Cache               Cache               Cache
         ┌─────────┐         ┌─────────┐
user: 100│ ≠       │user: 100│ ≠       │user: 100│
 (v1)    │         │  (v2)   │         │  (v3)   │
         └─────────┘         └─────────┘
   ↓                   ↓                   ↓
Diferentes versiones de datos!

SOLUCIÓN: Redis compartido
Instance A    Instance B    Instance C
   ↓             ↓             ↓
   └─────────────┴─────────────┘
           ↓
       Redis Cache
       ┌─────────┐
       │user: 100│ ÚNICA VERSIÓN
       └─────────┘
```

**Implementación:**

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

// Shared Redis client (todas las instancias usan el mismo Redis)
var redisClient = redis.NewClient(&redis.Options{
    Addr: "redis-master:6379",
})

type DistributedCache struct {
    redis *redis.Client
}

func NewDistributedCache() *DistributedCache {
    return &DistributedCache{redis: redisClient}
}

func (dc *DistributedCache) GetUser(ctx context.Context, userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // Todos los servidores consultan el MISMO Redis
    val, err := dc.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(val), &user)
        return &user, nil
    }
    
    // Cache miss en todas las instancias
    return db.GetUser(userID)
}

func (dc *DistributedCache) SetUser(ctx context.Context, user *User) error {
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    userJSON, _ := json.Marshal(user)
    
    // Escribir en Redis compartido
    return dc.redis.Set(ctx, cacheKey, userJSON, 1*time.Hour).Err()
}

func (dc *DistributedCache) InvalidateUser(ctx context.Context, userID string) error {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // Invalida para TODOS los servidores
    return dc.redis.Del(ctx, cacheKey).Err()
}
```

### 52.7.2 Coherencia Eventual

En sistemas distribuidos, la coherencia perfecta es imposible. Aceptamos **coherencia eventual**: los datos serán consistentes "eventualmente".

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

// Scenario: User updates name from Instance A
// Instance B debe enterarse

// Instance A: Actualiza
func UpdateUserNameInstanceA(userID, newName string) {
    ctx := context.Background()
    
    // 1. Actualizar DB
    db.UpdateUserName(userID, newName)
    
    // 2. Invalidar caché (Redis)
    redisClient.Del(ctx, fmt.Sprintf("user:%s", userID))
    
    // 3. Publicar evento
    redisClient.Publish(ctx, "user:updated", fmt.Sprintf(`{
        "event": "user_updated",
        "userID": "%s",
        "timestamp": %d
    }`, userID, time.Now().Unix()))
    
    fmt.Printf("Instance A: Updated user %s\n", userID)
}

// Instance B: Escucha eventos
func ListenForUserUpdatesInstanceB() {
    ctx := context.Background()
    pubsub := redisClient.Subscribe(ctx, "user:updated")
    defer pubsub.Close()
    
    ch := pubsub.Channel()
    for msg := range ch {
        var event map[string]interface{}
        json.Unmarshal([]byte(msg.Payload), &event)
        
        userID := event["userID"].(string)
        
        // Invalidar caché localmente
        cacheKey := fmt.Sprintf("user:%s", userID)
        redisClient.Del(ctx, cacheKey)
        
        fmt.Printf("Instance B: Invalidated cache for user %s\n", userID)
    }
}

// Timeline:
// T=0ms:   Instance A: Write DB + Delete Cache + Publish
// T=0ms:   Instance B: Cache still has old version (coherencia eventual)
// T=50ms:  Instance B: Receive event, invalidate cache
// T=50ms:  Next read on Instance B: Fresh data from DB
```

### 52.7.3 Cache Warming

Pre-cargar datos importantes en caché antes de que se necesiten.

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

func WarmCacheOnStartup() {
    ctx := context.Background()
    
    // Obtener datos críticos de BD
    topUsers, _ := db.GetTopUsers(1000)
    
    // Cargar en caché
    pipe := redisClient.Pipeline()
    
    for _, user := range topUsers {
        cacheKey := fmt.Sprintf("user:%s", user.ID)
        userJSON, _ := json.Marshal(user)
        pipe.Set(ctx, cacheKey, userJSON, 24*time.Hour)
    }
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        fmt.Printf("Error warming cache: %v\n", err)
        return
    }
    
    fmt.Printf("Cache warmed with %d top users\n", len(topUsers))
}

// Cache warming periódicamente
func PeriodicCacheWarmup() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        fmt.Println("Running periodic cache warmup...")
        WarmCacheOnStartup()
    }
}
```

### 52.7.4 Distributed Invalidation

Notificar a todas las instancias cuando cache debe ser invalidado.

```go
package main

import (
    "context"
    "fmt"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type InvalidationMessage struct {
    Keys      []string  `json:"keys"`
    Pattern   string    `json:"pattern"`  // Ej: "user:*"
    Timestamp time.Time `json:"timestamp"`
}

func BroadcastInvalidation(ctx context.Context, keys []string) error {
    msg := InvalidationMessage{
        Keys:      keys,
        Timestamp: time.Now(),
    }
    
    payload, _ := json.Marshal(msg)
    return redisClient.Publish(ctx, "cache:invalidation", string(payload)).Err()
}

func ListenForInvalidation(ctx context.Context) {
    pubsub := redisClient.Subscribe(ctx, "cache:invalidation")
    defer pubsub.Close()
    
    ch := pubsub.Channel()
    for msg := range ch {
        var invalidation InvalidationMessage
        json.Unmarshal([]byte(msg.Payload), &invalidation)
        
        // Invalidar localmente
        for _, key := range invalidation.Keys {
            redisClient.Del(ctx, key)
        }
        
        fmt.Printf("Invalidated %d keys\n", len(invalidation.Keys))
    }
}

// Pattern-based invalidation
func InvalidatePattern(ctx context.Context, pattern string) error {
    // Obtener todas las claves que coinciden
    var keys []string
    iter := redisClient.Scan(ctx, 0, pattern, 100).Iterator()
    
    for iter.Next(ctx) {
        keys = append(keys, iter.Val())
    }
    
    if len(keys) > 0 {
        redisClient.Del(ctx, keys...)
    }
    
    // Notificar a otros servidores
    msg := InvalidationMessage{
        Keys:      keys,
        Pattern:   pattern,
        Timestamp: time.Now(),
    }
    
    payload, _ := json.Marshal(msg)
    return redisClient.Publish(ctx, "cache:invalidation", string(payload)).Err()
}
```

---

## 52.8 - MEMCACHED

### 52.8.1 Protocol y Arquitectura

Memcached es más simple que Redis, optimizado para caché simple key-value.

```
Características:
- Solo strings (sin data types complejos)
- Más rápido que Redis para get/set simple
- Menos funcionalidades (no hay transactions, Pub/Sub)
- Bueno para caché puro (no persistencia)
```

**Comparativa:**

| Aspecto | Redis | Memcached |
|---------|-------|-----------|
| Data types | Strings, Lists, Sets, Zsets, Hashes | Solo Strings |
| Persistencia | Sí (RDB, AOF) | No |
| Transactions | Sí (MULTI/EXEC) | No |
| Pub/Sub | Sí | No |
| Lua Scripting | Sí | No |
| Memoria | Ms (features) | Menos |
| Latencia | 1-5ms | 0.5-2ms |
| Use case | General cache + DB | Caché puro |

### 52.8.2 go-memcache

```bash
go get github.com/bradfitz/gomemcache/memcache
```

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/bradfitz/gomemcache/memcache"
)

func MemcachedExample() {
    mc := memcache.New("localhost:11211")
    
    // SET
    err := mc.Set(&memcache.Item{
        Key:        "user:1",
        Value:      []byte(`{"name":"Alice","id":1}`),
        Expiration: 3600, // segundos
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // GET
    item, err := mc.Get("user:1")
    if err == memcache.ErrCacheMiss {
        fmt.Println("Cache miss")
    } else if err == nil {
        fmt.Printf("Found: %s\n", string(item.Value))
    }
    
    // DELETE
    mc.Delete("user:1")
    
    // MGET
    items, err := mc.GetMulti([]string{"user:1", "user:2", "user:3"})
    for key, item := range items {
        fmt.Printf("%s: %s\n", key, string(item.Value))
    }
    
    // ADD (solo si no existe)
    mc.Add(&memcache.Item{
        Key:   "session:abc",
        Value: []byte("session_data"),
    })
    
    // CAS (Compare-And-Swap, para atomicidad)
    item, _ := mc.Get("counter")
    // Modificar item...
    mc.CompareAndSwap(&memcache.Item{
        Key:   "counter",
        Value: []byte("new_value"),
        CAS:   item.CAS,
    })
    
    // FLUSH (limpiar todo)
    mc.FlushAll()
}
```

### 52.8.3 Comparativa Redis vs Memcached

**Cuando usar REDIS:**
- Necesitas data types complejos (Lists, Sets, Hashes)
- Persistence importante
- Transactions y Lua scripts
- Pub/Sub
- High Availability (Sentinel, Cluster)

**Cuando usar MEMCACHED:**
- Caché puro, solo key-value
- Máxima velocidad para operaciones simple
- Clustering distribuido (consistente hashing)
- Bajo overhead de memoria
- No necesitas persistencia

**Benchmark:**

```
1,000,000 operaciones:

Redis (Set + Get):     ~800ms (1.25M ops/sec)
Memcached (Set + Get): ~400ms (2.5M ops/sec)

Memcached es ~3x más rápido para get/set puro
Redis ofrece 10x más funcionalidades
```

### 52.8.4 Memcached Clustering

```go
package main

import (
    "fmt"
    "github.com/bradfitz/gomemcache/memcache"
)

func MemcachedCluster() {
    // Crear cliente con múltiples servidores
    mc := memcache.New(
        "memcached-node-1:11211",
        "memcached-node-2:11211",
        "memcached-node-3:11211",
    )
    
    // Memcached usa consistent hashing
    // Cada clave va automáticamente al mismo nodo
    
    item := &memcache.Item{
        Key:   "user:100",
        Value: []byte("alice"),
    }
    mc.Set(item)
    
    // La clave "user:100" siempre va al mismo nodo
    // Aunque agregues/remuevas nodos, mayoría de claves se quedan donde estaban
    
    retrieved, _ := mc.Get("user:100")
    fmt.Printf("Value: %s\n", string(retrieved.Value))
}
```

---

## 52.9 - CACHING DE HTTP

### 52.9.1 Headers de HTTP Caching

HTTP tiene mecanismos nativos de caché controlables desde el servidor.

```
Headers importantes:

Cache-Control: max-age=3600       # Caché durante 1 hora
Cache-Control: no-cache           # Validar con servidor antes de usar
Cache-Control: no-store           # Nunca cachear
Cache-Control: public             # Cacheable por navegadores y proxies
Cache-Control: private            # Solo navegador, no proxies
Cache-Control: must-revalidate    # Revalidar si está expirado

ETag: "abc123"                    # Hash del contenido
If-None-Match: "abc123"           # Si el ETag coincide, 304 Not Modified

Last-Modified: Mon, 01 Jan 2024 00:00:00 GMT
If-Modified-Since: Mon, 01 Jan 2024 00:00:00 GMT
```

**Go HTTP Server con caching:**

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func CacheableResponse(w http.ResponseWriter, r *http.Request) {
    data := `{"id":1,"name":"Alice","role":"admin"}`
    hash := calculateHash(data) // MD5 o SHA256
    
    // Establecer headers de caché
    w.Header().Set("Cache-Control", "public, max-age=3600")
    w.Header().Set("ETag", fmt.Sprintf(`"%s"`, hash))
    w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
    
    // Validar ETag (cliente envía If-None-Match)
    if match := r.Header.Get("If-None-Match"); match == fmt.Sprintf(`"%s"`, hash) {
        w.WriteHeader(http.StatusNotModified)
        return
    }
    
    // Validar Last-Modified
    if modifiedSince := r.Header.Get("If-Modified-Since"); modifiedSince != "" {
        w.WriteHeader(http.StatusNotModified)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, data)
}

func calculateHash(data string) string {
    hash := md5.Sum([]byte(data))
    return fmt.Sprintf("%x", hash)
}

func main() {
    http.HandleFunc("/api/users/1", CacheableResponse)
    http.ListenAndServe(":8080", nil)
}
```

### 52.9.2 HTTP Caching Layers

```
Browser Cache (LocalStorage, Service Workers)
         ↓
CDN Cache (Cloudflare, Akamai)
         ↓
Proxy Cache (nginx, varnish)
         ↓
Application HTTP Cache Middleware
         ↓
Application In-Memory Cache
         ↓
Database
```

### 52.9.3 Middleware de Caching

```go
package main

import (
    "crypto/md5"
    "fmt"
    "net/http"
    "time"
    "github.com/redis/go-redis/v9"
)

type CacheMiddleware struct {
    redis *redis.Client
    ttl   time.Duration
}

func (cm *CacheMiddleware) CacheHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Solo cachear GET requests
        if r.Method != http.MethodGet {
            next.ServeHTTP(w, r)
            return
        }
        
        ctx := r.Context()
        cacheKey := fmt.Sprintf("http:%s:%s", r.Method, r.URL.Path)
        
        // Intentar caché
        if cached, err := cm.redis.Get(ctx, cacheKey).Result(); err == nil {
            w.Header().Set("X-Cache", "HIT")
            w.Header().Set("Content-Type", "application/json")
            fmt.Fprint(w, cached)
            return
        }
        
        // Capturar response
        recorder := &responseRecorder{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        next.ServeHTTP(recorder, r)
        
        // Cachear respuesta
        if recorder.statusCode == http.StatusOK {
            cm.redis.Set(ctx, cacheKey, string(recorder.body), cm.ttl)
        }
        
        w.Header().Set("X-Cache", "MISS")
    })
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
    body       []byte
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
    rr.body = append(rr.body, b...)
    return rr.ResponseWriter.Write(b)
}

func main() {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    
    cm := &CacheMiddleware{
        redis: client,
        ttl:   1 * time.Hour,
    }
    
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, `{"data":"expensive computation result"}`)
    })
    
    http.Handle("/api/data", cm.CacheHandler(handler))
    http.ListenAndServe(":8080", nil)
}
```

### 52.9.4 CDN Integration

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

// Headers para CDN (Cloudflare, Akamai, CloudFront)
func CDNOptimizedResponse(w http.ResponseWriter, r *http.Request) {
    // Cache en CDN por 24 horas
    w.Header().Set("Cache-Control", "public, max-age=86400")
    
    // Permitir caching en CDN también
    w.Header().Set("Surrogate-Control", "public, max-age=86400")
    
    // Purge cache cuando sea necesario (Cloudflare API)
    // curl -X POST "https://api.cloudflare.com/client/v4/zones/{zone_id}/purge_cache" \
    //      -H "Authorization: Bearer {token}" \
    //      -d '{"files":["https://example.com/api/data"]}'
    
    w.Header().Set("Vary", "Accept-Encoding")
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, `{"cached":"by_cdn"}`)
}

// Invalidar cache en CDN
func InvalidateCDNCache(url string) error {
    // Cloudflare example
    // Implementar según el servicio CDN
    fmt.Printf("Invalidating CDN cache for: %s\n", url)
    return nil
}

func main() {
    http.HandleFunc("/api/public-data", CDNOptimizedResponse)
    http.ListenAndServe(":8080", nil)
}
```

---

## 52.10 - TESTING Y MONITORING

### 52.10.1 Mock Caches

```go
package main

import (
    "context"
    "sync"
)

// MockCache para testing
type MockCache struct {
    data     map[string]interface{}
    mu       sync.RWMutex
    callLog  []string
}

func NewMockCache() *MockCache {
    return &MockCache{
        data:    make(map[string]interface{}),
        callLog: make([]string, 0),
    }
}

func (mc *MockCache) Get(ctx context.Context, key string) (interface{}, error) {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    
    mc.callLog = append(mc.callLog, "GET:"+key)
    val, exists := mc.data[key]
    if exists {
        return val, nil
    }
    return nil, ErrCacheMiss
}

func (mc *MockCache) Set(ctx context.Context, key string, value interface{}) error {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    
    mc.callLog = append(mc.callLog, "SET:"+key)
    mc.data[key] = value
    return nil
}

func (mc *MockCache) Delete(ctx context.Context, key string) error {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    
    mc.callLog = append(mc.callLog, "DEL:"+key)
    delete(mc.data, key)
    return nil
}

func (mc *MockCache) GetCallLog() []string {
    mc.mu.RLock()
    defer mc.mu.RUnlock()
    return append([]string{}, mc.callLog...)
}

// Test
func TestUserRepository() {
    cache := NewMockCache()
    repo := NewUserRepository(cache)
    
    // Primer acceso: cache miss
    user, _ := repo.GetUser("user:1")
    
    // Segundo acceso: cache hit
    user, _ = repo.GetUser("user:1")
    
    log := cache.GetCallLog()
    // Verificar que se hizo SET seguido de GET
    assert(log[0] == "SET:user:1", "Primer call debe ser SET")
    assert(log[1] == "GET:user:1", "Segundo call debe ser GET")
}
```

### 52.10.2 Cache Hit Rates

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type CacheMetrics struct {
    hits       int64
    misses     int64
    evictions  int64
    mu         sync.RWMutex
}

func (cm *CacheMetrics) RecordHit() {
    atomic.AddInt64(&cm.hits, 1)
}

func (cm *CacheMetrics) RecordMiss() {
    atomic.AddInt64(&cm.misses, 1)
}

func (cm *CacheMetrics) RecordEviction() {
    atomic.AddInt64(&cm.evictions, 1)
}

func (cm *CacheMetrics) GetHitRate() float64 {
    hits := atomic.LoadInt64(&cm.hits)
    misses := atomic.LoadInt64(&cm.misses)
    total := hits + misses
    
    if total == 0 {
        return 0
    }
    return float64(hits) / float64(total) * 100
}

func (cm *CacheMetrics) Stats() map[string]interface{} {
    hits := atomic.LoadInt64(&cm.hits)
    misses := atomic.LoadInt64(&cm.misses)
    evictions := atomic.LoadInt64(&cm.evictions)
    total := hits + misses
    
    var hitRate float64
    if total > 0 {
        hitRate = float64(hits) / float64(total) * 100
    }
    
    return map[string]interface{}{
        "hits":       hits,
        "misses":     misses,
        "evictions":  evictions,
        "total":      total,
        "hit_rate":   fmt.Sprintf("%.2f%%", hitRate),
    }
}

// Endpoint para métricas
func MetricsEndpoint(metrics *CacheMetrics) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        stats := metrics.Stats()
        json.NewEncoder(w).Encode(stats)
    })
}
```

### 52.10.3 Memory Usage Monitoring

```go
package main

import (
    "fmt"
    "runtime"
)

func MonitorMemoryUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("=== Memory Stats ===\n")
    fmt.Printf("Alloc: %v MB\n", m.Alloc / 1024 / 1024)
    fmt.Printf("TotalAlloc: %v MB\n", m.TotalAlloc / 1024 / 1024)
    fmt.Printf("Sys: %v MB\n", m.Sys / 1024 / 1024)
    fmt.Printf("NumGC: %v\n", m.NumGC)
    fmt.Printf("Goroutines: %v\n", runtime.NumGoroutine())
}

// Monitoring periódico
func StartMemoryMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for range ticker.C {
        MonitorMemoryUsage()
    }
}

// Redis memory info
func MonitorRedisMemory(client *redis.Client) {
    ctx := context.Background()
    info, _ := client.Info(ctx, "memory").Result()
    fmt.Printf("Redis Memory Info:\n%s\n", info)
}
```

### 52.10.4 Performance Benchmarks

```go
package main

import (
    "context"
    "fmt"
    "testing"
    "time"
    "github.com/redis/go-redis/v9"
)

func BenchmarkLocalCacheVsRedis(b *testing.B) {
    ctx := context.Background()
    
    // Caché local
    localCache := make(map[string]interface{})
    
    // Redis
    redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    defer redisClient.Close()
    
    // Benchmark: Local cache
    b.Run("LocalCache", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            key := fmt.Sprintf("key:%d", i%1000)
            localCache[key] = "value"
        }
    })
    
    // Benchmark: Redis
    b.Run("Redis", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            key := fmt.Sprintf("key:%d", i%1000)
            redisClient.Set(ctx, key, "value", 0)
        }
    })
}

// Resultado típico:
// LocalCache:  1,000,000 ops en 10ms (100M ops/sec)
// Redis:       1,000,000 ops en 500ms (2M ops/sec)
// Local cache es 50x más rápido, pero no distribuido

func BenchmarkCachePatterns(b *testing.B) {
    // Cache-Aside vs Write-Through
    // Implementar benchmarks específicos
}
```

---

## 52.11 - BUENAS PRÁCTICAS Y CASE STUDIES

### 52.11.1 Cache Key Design

**✅ BUENO: Structured, namespaced keys**

```go
// Formato: "namespace:entity_type:entity_id:operation"

"user:profile:123"              // Perfil del usuario 123
"user:settings:123"             // Configuración del usuario 123
"post:detail:456"               // Detalles del post 456
"post:comments:456"             // Comentarios del post 456
"session:token:abc123xyz"       // Session token
"rate-limit:api:user:123:hour"  // Rate limit per hour
"leaderboard:weekly"            // Leaderboard semanal
```

**❌ MALO: Generic, unclear keys**

```go
"user"          // ¿Qué usuario?
"data:123"      // ¿Qué dato?
"cache1"        // ¿Qué?
"temp"          // Demasiado vago
```

**Implementación:**

```go
package main

const (
    UserPrefix           = "user"
    PostPrefix           = "post"
    SessionPrefix        = "session"
    RateLimitPrefix      = "ratelimit"
)

func buildCacheKey(parts ...string) string {
    key := ""
    for i, part := range parts {
        if i > 0 {
            key += ":"
        }
        key += part
    }
    return key
}

// Uso
userKey := buildCacheKey(UserPrefix, "profile", userID)
postKey := buildCacheKey(PostPrefix, "comments", postID)
```

### 52.11.2 TTL Strategy

| Tipo de Dato | TTL Recomendado | Razón |
|-------------|-----------------|-------|
| Perfil de usuario | 1 hora | Cambia raramente |
| Feed de usuario | 5 minutos | Actualización frecuente |
| Datos de sesión | 24 horas | Session expiration |
| Rate limit | 1 minuto | Precisión importante |
| Leaderboard | 5 minutos | Actualización cada poco |
| Configuración de app | 1 hora | Estable |
| Resultados de búsqueda | 30 minutos | Más variabilidad |

```go
var TTLConfig = map[string]time.Duration{
    "user:profile":    1 * time.Hour,
    "user:feed":       5 * time.Minute,
    "session":         24 * time.Hour,
    "ratelimit":       1 * time.Minute,
    "leaderboard":     5 * time.Minute,
    "config":          1 * time.Hour,
}

func GetTTL(key string) time.Duration {
    for prefix, ttl := range TTLConfig {
        if strings.HasPrefix(key, prefix) {
            return ttl
        }
    }
    return 10 * time.Minute // Default
}
```

### 52.11.3 Database + Cache Consistency

**Patrón Safe: Invalidate on Write**

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func UpdateUserProfile(ctx context.Context, userID string, data map[string]interface{}) error {
    // 1. Actualizar DB
    err := db.UpdateUser(userID, data)
    if err != nil {
        return fmt.Errorf("DB update failed: %w", err)
    }
    
    // 2. Invalidar caché
    cacheKey := fmt.Sprintf("user:profile:%s", userID)
    if err := redisClient.Del(ctx, cacheKey).Err(); err != nil {
        // Logging importante, pero no fallar el request
        log.Printf("Failed to invalidate cache: %v", err)
    }
    
    // 3. Notificar listeners
    redisClient.Publish(ctx, "user:updated", fmt.Sprintf(`{
        "userID": "%s",
        "event": "profile_updated",
        "timestamp": %d
    }`, userID, time.Now().Unix()))
    
    return nil
}

// ❌ MALO: No invalidar
func BadUpdateUser(ctx context.Context, userID string, data map[string]interface{}) error {
    return db.UpdateUser(userID, data)
    // ❌ Cache viejo devuelve datos desactualizados
}

// ✅ BIEN: Invalidar + Log + Notify
func GoodUpdateUser(ctx context.Context, userID string, data map[string]interface{}) error {
    return UpdateUserProfile(ctx, userID, data)
}
```

### 52.11.4 Error Handling (Cache Failures)

**✅ BIEN: Graceful Degradation**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/redis/go-redis/v9"
)

func GetUserWithFallback(ctx context.Context, userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // Intentar caché
    cached, err := redisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    if err != redis.Nil {
        // Error de conexión a Redis
        log.Printf("Redis error (cache unavailable): %v", err)
        // Continuar, iremos a DB
    }
    
    // Cache miss o error, consultar DB
    user, dbErr := db.GetUser(userID)
    if dbErr != nil {
        return nil, fmt.Errorf("failed to get user from DB: %w", dbErr)
    }
    
    // Intentar cachear (best-effort, no críticamente)
    go func() {
        userJSON, _ := json.Marshal(user)
        if err := redisClient.Set(context.Background(), 
            cacheKey, userJSON, 1*time.Hour).Err(); err != nil {
            log.Printf("Failed to cache user: %v", err)
        }
    }()
    
    return user, nil
}

// ❌ MALO: Crashing si cache falla
func BadGetUser(userID string) (*User, error) {
    val, err := redisClient.Get(context.Background(), "user:"+userID).Result()
    if err != nil {
        panic("Cache failed!") // ❌ CRASH TOTAL
    }
    var user User
    json.Unmarshal([]byte(val), &user)
    return &user, nil
}
```

### 52.11.5 Real-World Case Studies

**Case 1: Stripe - Payment Processing**

```
Desafío: Procesar $1B/día sin caché causaría:
- Millones de queries a DB por segundo
- Latencia inaceptable
- Cascada de fallos

Solución:
1. Redis + Cluster Mode
   - Caché de customer profiles (1 segundo TTL)
   - Caché de payment methods (1 hora TTL)
   - Caché de rate limits (1 minuto TTL)

2. Multi-level architecture:
   - L1: Application memory (microsegundos)
   - L2: Redis Cluster (milisegundos)
   - L3: PostgreSQL (centenas de ms)

3. Invalidation strategy:
   - Event-based: Cuando customer actualiza perfil
   - TTL: Expiración automática
   - Hybrid: TTL corto + event invalidation

Resultado: 99.99% uptime, latencia p99 < 100ms
```

**Case 2: Netflix - Content Delivery**

```
Desafío: Servir contenido a 200M usuarios globalmente
        2 petabytes de metadata

Solución:
1. Multi-tier caching:
   - Browser cache: 24 horas (imágenes, JS, CSS)
   - CDN cache: 1 hora (recomendaciones)
   - Application cache: 5 minutos (metadata)
   - Database cache: índices, buffer pools

2. Cache warming:
   - Pre-cargar datos populares en Redis
   - Top 1000 películas siempre cacheadas
   - Predictions ML para calentar caché

3. Invalidation:
   - Event-driven desde CMS
   - Invalidación en cascada por popularidad
   - Probabilistic expiration para evitar thundering herd

Resultado: 99.5% cache hit rate, mejor UX
```

**Case 3: GitHub - Repository Data**

```
Desafío: 100M+ repositorios, queries complejas
        Consistencia crítica (commits, issues)

Solución:
1. Cache-aside con strong consistency:
   - Verificar DB siempre para datos críticos
   - Caché solo para "nice-to-have" data
   - TTL corto: 1-5 minutos

2. Distributed cache:
   - Redis Cluster con replica
   - Automatic failover
   - Cache warming para repos populares

3. Instrumentation:
   - Monitorear hit rates por tipo de dato
   - Alert si hit rate < 80%
   - Trace cache misses a nivel de usuario

Resultado: Sub-100ms latency, 95%+ consistency
```

### 52.11.6 Antipatterns ❌ vs Best Practices ✅

```
ANTIPATTERN 1: No TTL
 cache[key] = value // Nunca expira
   Resultado: Memory leak, OOM después de millones de requests

 cache.Set(key, value, 1*time.Hour)
   Resultado: Memoria controlada, datos frescos


ANTIPATTERN 2: Same key format everywhere
 cache["user"] = userData
   cache["post"] = postData
   Namespacing confuso, colisiones potenciales

 cache["user:profile:123"] = userData
   cache["post:detail:456"] = postData
   Estructura clara, sin colisiones


ANTIPATTERN 3: Cache sin estrategia de invalidación
 Guardar datos en caché pero nunca invalidar
   Resultado: Datos desactualizados por horas

 Invalidar en cada write:
   - Direct: redis.Del(key)
   - Event-based: Publish invalidation event
   - TTL: Expiración automática


ANTIPATTERN 4: Cache failures causan crash
 val, err := redis.Get(key)
   if err != nil { panic(err) }
   Resultado: Cache falla → App falla

 Graceful degradation:
   val, err := redis.Get(key)
   if err != nil {
       log.Printf("Cache miss: %v", err)
       val = db.Query(key)
   }
   Resultado: Cache falla → App devuelve datos frescos de DB


ANTIPATTERN 5: Cachear absolutamente todo
 cache.Set("timestamp", time.Now())
   cache.Set("random_uuid", uuid.New())
   Resultado: Cache inútil, memoria desperdiciada

 Cachear selectivamente:
   - User profiles (reutilizado)
   - Post content (reutilizado)
   - NO: Timestamps actuales
   - NO: Valores únicos


ANTIPATTERN 6: TTL demasiado corto
 cache.Set(key, value, 100*time.Millisecond)
   Resultado: Cache hit rate ~0%, recursos desperdiciados

 TTL según frecuencia de cambio:
   - Datos estáticos: 24 horas
   - Datos semi-estáticos: 1 hora
   - Datos dinámicos: 5-10 minutos
   - Rate limit: 1 minuto


ANTIPATTERN 7: Caché inconsistente en cluster
 Cada instancia tiene su caché local
   Resultado: Datos diferentes en cada servidor

 Caché distribuida (Redis):
   Todos los servidores consultan el mismo Redis
   Resultado: Consistencia garantizada


ANTIPATTERN 8: No monitorear cache
 Caché silencioso, sin métricas
   Resultado: Hit rate desconocida, problemas no detectados

 Monitorear activamente:
   - Hit rate por tipo de dato
   - Memory usage
   - Eviction rate
   - Performance percentiles
```

### 52.11.7 Ejercicio 1: In-Memory Cache Simple

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Ejercicio 1: Implementar un User Store simple con caché
type User struct {
    ID    string
    Name  string
    Email string
}

type UserStore struct {
    cache map[string]*User
    mu    sync.RWMutex
}

func NewUserStore() *UserStore {
    return &UserStore{
        cache: make(map[string]*User),
    }
}

func (us *UserStore) Get(userID string) (*User, bool) {
    us.mu.RLock()
    defer us.mu.RUnlock()
    
    user, exists := us.cache[userID]
    return user, exists
}

func (us *UserStore) Set(userID string, user *User) {
    us.mu.Lock()
    defer us.mu.Unlock()
    
    us.cache[userID] = user
}

func (us *UserStore) Delete(userID string) {
    us.mu.Lock()
    defer us.mu.Unlock()
    
    delete(us.cache, userID)
}

// TODO en ejercicio: Agregá TTL automático usando goroutines
// TODO en ejercicio: Implementá LRU cuando se exceda 1000 items

func Exercise1() {
    store := NewUserStore()
    
    // Set
    store.Set("user:1", &User{
        ID:    "user:1",
        Name:  "Alice",
        Email: "alice@example.com",
    })
    
    // Get
    user, exists := store.Get("user:1")
    fmt.Printf("User exists: %v, data: %+v\n", exists, user)
    
    // Delete
    store.Delete("user:1")
    user, exists = store.Get("user:1")
    fmt.Printf("After delete - exists: %v\n", exists)
}
```

### 52.11.8 Ejercicio 2: Redis Integration

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

// Ejercicio 2: Counter con Redis + Expiry
func Exercise2() {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    defer client.Close()
    
    ctx := context.Background()
    
    // Contar visitas a un artículo
    articleID := "article:123"
    
    // Incrementar contador
    for i := 0; i < 10; i++ {
        count, err := client.Incr(ctx, articleID).Result()
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
        fmt.Printf("Visit %d: count = %d\n", i+1, count)
    }
    
    // Establecer expiración (1 hora)
    client.Expire(ctx, articleID, 1*time.Hour)
    
    // Obtener TTL restante
    ttl, _ := client.TTL(ctx, articleID).Result()
    fmt.Printf("TTL: %v\n", ttl)
    
    // TODO: Agregar rate limiting (máximo 100 views por minuto)
    // TODO: Agregar stats (views totales, visitors únicos)
}
```

### 52.11.9 Ejercicio 3: Cache-Aside Pattern

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

// Ejercicio 3: Database + Redis Cache-Aside
type Product struct {
    ID    string
    Name  string
    Price float64
}

type ProductService struct {
    redis *redis.Client
    db    *Database // Tu BD
}

func (ps *ProductService) GetProduct(ctx context.Context, productID string) (*Product, error) {
    cacheKey := fmt.Sprintf("product:%s", productID)
    
    // Paso 1: Intentar caché
    cached, err := ps.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var product Product
        json.Unmarshal([]byte(cached), &product)
        fmt.Printf("Cache HIT: %s\n", productID)
        return &product, nil
    }
    
    if err != redis.Nil {
        fmt.Printf("Redis error: %v\n", err)
    }
    
    // Paso 2: Cache miss, consultar BD
    fmt.Printf("Cache MISS: %s, querying DB\n", productID)
    product, err := ps.db.GetProduct(productID)
    if err != nil {
        return nil, err
    }
    
    // Paso 3: Guardar en caché
    productJSON, _ := json.Marshal(product)
    ps.redis.Set(ctx, cacheKey, productJSON, 1*time.Hour)
    
    return product, nil
}

func Exercise3() {
    // TODO: Implementar
    // 1. Crear ProductService
    // 2. Getproduct 5 veces
    // 3. Verificar que primer acceso consulta DB
    // 4. Segundo acceso viene de caché
}
```

### 52.11.10 Ejercicio 4: Write-Through Pattern

```go
package main

// Ejercicio 4: Write-Through (caché siempre actualizado)
func (ps *ProductService) UpdateProduct(ctx context.Context, product *Product) error {
    cacheKey := fmt.Sprintf("product:%s", product.ID)
    
    // Paso 1: Actualizar caché primero
    productJSON, _ := json.Marshal(product)
    ps.redis.Set(ctx, cacheKey, productJSON, 1*time.Hour)
    
    // Paso 2: Actualizar BD (sincrónico)
    err := ps.db.UpdateProduct(product)
    if err != nil {
        // Revertir cach si falla
        ps.redis.Del(ctx, cacheKey)
        return err
    }
    
    fmt.Printf("Product %s updated (write-through)\n", product.ID)
    return nil
}

func Exercise4() {
    // TODO: Implementar
    // 1. Actualizar producto
    // 2. Verificar que caché se actualiza
    // 3. Verificar que BD se actualiza
    // 4. Simular fallo de BD y verificar rollback
}
```

### 52.11.11 Ejercicio 5: Production System

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync/atomic"
    "time"
    "github.com/redis/go-redis/v9"
)

// Ejercicio 5: Sistema de producción completo
type CacheService struct {
    redis   *redis.Client
    metrics *CacheMetrics
}

type CacheMetrics struct {
    hits       int64
    misses     int64
    evictions  int64
    errors     int64
}

func NewCacheService(redisAddr string) *CacheService {
    client := redis.NewClient(&redis.Options{Addr: redisAddr})
    
    return &CacheService{
        redis:   client,
        metrics: &CacheMetrics{},
    }
}

func (cs *CacheService) GetWithFallback(ctx context.Context, 
    key string, 
    loader func() (interface{}, error),
) (interface{}, error) {
    // Intentar caché
    val, err := cs.redis.Get(ctx, key).Result()
    if err == nil {
        atomic.AddInt64(&cs.metrics.hits, 1)
        return val, nil
    }
    
    if err != redis.Nil {
        // Error de conexión
        log.Printf("Redis error: %v", err)
        atomic.AddInt64(&cs.metrics.errors, 1)
    }
    
    // Cache miss, cargar datos
    atomic.AddInt64(&cs.metrics.misses, 1)
    
    data, err := loader()
    if err != nil {
        return nil, err
    }
    
    // Cachear
    dataJSON, _ := json.Marshal(data)
    cs.redis.Set(ctx, key, dataJSON, 1*time.Hour)
    
    return data, nil
}

func (cs *CacheService) MetricsHandler(w http.ResponseWriter, r *http.Request) {
    hits := atomic.LoadInt64(&cs.metrics.hits)
    misses := atomic.LoadInt64(&cs.metrics.misses)
    errors := atomic.LoadInt64(&cs.metrics.errors)
    
    total := hits + misses
    var hitRate float64
    if total > 0 {
        hitRate = float64(hits) / float64(total) * 100
    }
    
    metrics := map[string]interface{}{
        "hits":     hits,
        "misses":   misses,
        "errors":   errors,
        "hit_rate": fmt.Sprintf("%.2f%%", hitRate),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(metrics)
}

func Exercise5() {
    // TODO: Implementar
    // 1. Crear CacheService
    // 2. Setup HTTP server
    // 3. Agregar rutas:
    //    - GET /api/products/:id (con caché)
    //    - POST /api/products/:id (actualizar con invalidación)
    //    - GET /metrics (métricas)
    // 4. Test concurrencia
    // 5. Monitorear hit rate
    // 6. Implementar graceful degradation
}
```

---

## CONCLUSIÓN

El caching es un arte. Requiere entender:

1. **Arquitectura**: Múltiples niveles de caché
2. **Patrones**: Cache-aside, write-through, write-behind
3. **Problemas**: Stampede, penetration, avalanche
4. **Herramientas**: Redis, Memcached, HTTP headers
5. **Consistencia**: Balance entre performance y corrección
6. **Monitoreo**: Métricas, hit rates, latencia

**Reglas de Oro:**

 Siempre set TTL  
 Cachear selectivamente  
 Invalidar on write  
 Monitorear hit rates  
 Graceful degradation  
 Medir antes y después  

**Performance sin caché**: Inaceptable  
**Performance con caché bien diseñado**: Excelente  

El caché transforma una aplicación de "lenta y escalable" a "rápida y escalable". Úsalo sabiamente.

---

**Referencias:**
- Redis Official Documentation
- Stripe Engineering Blog: Caching at Scale
- Google SRE Book: Chapter on Caching
- Netflix Tech Blog: Cache Architecture


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/52-caching-system/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/52-caching-system):

```bash
cd examples/52-caching-system
go run .
```
