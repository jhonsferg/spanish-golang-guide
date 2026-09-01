# Capítulo 24: Sync - Primitivas de sincronización

## Introducción

El paquete `sync` de Go proporciona primitivas de sincronización de bajo nivel para coordinar el acceso a recursos compartidos entre goroutines. A diferencia de los channels (que comunican datos), estas primitivas protegen el acceso a variables y estructuras de datos compartidas mediante mecanismos de exclusión mutua, condiciones de espera y coordinación.

**Pregunta fundamental:** ¿Por qué Go tiene `sync` si tiene channels?

- **Channels**: Para comunicación entre goroutines (paso de mensajes)
- **Sync**: Para proteger datos compartidos (exclusión mutua)

Go favorece channels, pero existen casos donde sync es más eficiente o apropiado. Este capítulo explora cuándo, cómo y por qué usar estas primitivas.

---

## 24.1 Cuándo No Usar Sync (Preferir Channels)

### 24.1.1 Filosofía Go: Canales Primero

Go prefiere channels como paradigma de concurrencia. La regla general:

> **"No comuniques compartiendo memoria; comparte memoria comunicando"**

```go
// ❌ ANTI-PATRÓN: Usar sync.Mutex para coordinar lógica
var counter int
var mu sync.Mutex

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}

// ✓ PATRÓN: Usar channels para coordinación
func counterService(inc <-chan struct{}, result chan<- int) {
    count := 0
    for range inc {
        count++
        result <- count
    }
}
```

### 24.1.2 Ventajas de Channels sobre Sync

| Aspecto | Channels | Sync |
|--------|----------|------|
| Intención clara | Comunicación implícita | Requiere documentación |
| Deadlocks | Detectables (goroutine leak) | Sutiles, difficiles de detectar |
| Flujo de datos | Visible en el código | Oculto en critical sections |
| Testing | Más predecible | Requiere race detector |
| Composición | Fácil con select | Requiere coordinación manual |

### 24.1.3 Cuándo Sync Es Apropiado

**Caso 1: Acceso Frecuente a Estructura Compartida**

```go
// Cache con R/W frecuentes: RWMutex es mejor que channel
type Cache struct {
    mu    sync.RWMutex
    items map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}
```

**Caso 2: Inicialización Una Sola Vez**

```go
// sync.Once es más elegante que canales para esto
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

**Caso 3: Coordinación de Grupo (WaitGroup)**

```go
// Esperar a que terminen N goroutines
var wg sync.WaitGroup
wg.Add(10)
for i := 0; i < 10; i++ {
    go func() {
        defer wg.Done()
        // trabajo...
    }()
}
wg.Wait()
```

### 24.1.4 Matriz de Decisión

```
¿Necesitas comunicar datos?
    → SÍ: Usa channels
    → NO: ¿Necesitas exclusión mutua?
        → SÍ: Usa sync
        → NO: Considera atomic operations o immutable data

¿Múltiples lectores, pocos escritores?
    → SÍ: sync.RWMutex
    → NO: sync.Mutex

¿Esperar a N tareas?
    → SÍ: sync.WaitGroup
    → NO: ¿Ejecutar código una sola vez?
        → SÍ: sync.Once
        → NO: ¿Esperar condición específica?
            → SÍ: sync.Cond
            → NO: ¿Limitar concurrencia?
                → SÍ: Channel como semáforo
                → NO: Reconsiderar diseño
```

---

## 24.2 Mutex: Exclusión Mutua

### 24.2.1 Concepto Fundamental

Un `Mutex` (mutual exclusion lock) asegura que solo una goroutine puede acceder a una sección crítica (critical section) a la vez.

```
ESTADO DEL MUTEX:

[Desbloqueado]
        ↓
   Goroutine A llama Lock()
        ↓
[Bloqueado por A] ← Goroutines B, C, D esperan
        ↓
   Goroutine A llama Unlock()
        ↓
[Desbloqueado]
        ↓
   Una goroutine en espera (ej: B) obtiene lock
```

### 24.2.2 API Básica

```go
type Mutex struct {
    // Sin campos exportados
}

func (m *Mutex) Lock()       // Adquirir lock (bloquea si no disponible)
func (m *Mutex) Unlock()     // Liberar lock (panic si no era tenedor)
```

**Punto clave:** `Mutex` no es reentrante. La misma goroutine NO puede hacer Lock() dos veces.

### 24.2.3 Patrón: Defer para Unlock Garantizado

```go
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (sc *SafeCounter) Increment() {
    sc.mu.Lock()
    defer sc.mu.Unlock()  // ✓ Garantizado aunque haya panic
    sc.value++
}
```

### 24.2.4 Estructuras de Datos Protegidas

```go
type BankAccount struct {
    mu      sync.Mutex
    balance float64
}

func (ba *BankAccount) Deposit(amount float64) {
    ba.mu.Lock()
    defer ba.mu.Unlock()
    ba.balance += amount
}

func (ba *BankAccount) Withdraw(amount float64) bool {
    ba.mu.Lock()
    defer ba.mu.Unlock()
    if ba.balance >= amount {
        ba.balance -= amount
        return true
    }
    return false
}

func (ba *BankAccount) Balance() float64 {
    ba.mu.Lock()
    defer ba.mu.Unlock()
    return ba.balance
}
```

### 24.2.5 Locks Anidados: El Problema del Deadlock

```go
type Transfer struct {
    mu sync.Mutex
    accounts map[int]*BankAccount
}

// ❌ DEADLOCK POTENCIAL
func (t *Transfer) Move(from, to int, amount float64) {
    t.accounts[from].mu.Lock()
    defer t.accounts[from].mu.Unlock()

    t.accounts[to].mu.Lock()      // ¿Qué si otro thread
    defer t.accounts[to].mu.Unlock()  // hace Move(to, from)?

    t.accounts[from].balance -= amount
    t.accounts[to].balance += amount
}

// ✓ SOLUCIÓN: Ordenar por ID para evitar deadlock
func (t *Transfer) Move(from, to int, amount float64) {
    if from > to {
        from, to = to, from
    }

    t.accounts[from].mu.Lock()
    defer t.accounts[from].mu.Unlock()

    t.accounts[to].mu.Lock()
    defer t.accounts[to].mu.Unlock()

    t.accounts[from].balance -= amount
    t.accounts[to].balance += amount
}
```

### 24.2.6 Hold Time: Mantener Locks Poco Tiempo

```go
type Database struct {
    mu    sync.Mutex
    cache map[string]interface{}
}

// ❌ ANTI-PATRÓN: Operación costosa bajo lock
func (db *Database) GetAndProcess(key string) interface{} {
    db.mu.Lock()
    defer db.mu.Unlock()

    val := db.cache[key]

    // 1 segundo de procesamiento bajo lock!
    result := expensiveComputation(val)

    db.cache[key] = result
    return result
}

// ✓ PATRÓN: Copiar datos bajo lock, procesar fuera
func (db *Database) GetAndProcess(key string) interface{} {
    var val interface{}

    db.mu.Lock()
    val = db.cache[key]
    db.mu.Unlock()  // Liberar rápidamente

    result := expensiveComputation(val)  // Fuera del lock

    db.mu.Lock()
    db.cache[key] = result
    db.mu.Unlock()

    return result
}
```

### 24.2.7 Granularidad de Locks

```go
// ❌ GRANULARIDAD GRUESA: Un lock para todo
type Server struct {
    mu       sync.Mutex
    requests int
    errors   int
    latency  time.Duration
}

// ✓ GRANULARIDAD FINA: Locks independientes
type Server struct {
    requestsMu sync.Mutex
    requests   int

    errorsMu   sync.Mutex
    errors     int

    latencyMu  sync.Mutex
    latency    time.Duration
}
```

### 24.2.8 Comparación con Otros Lenguajes

**Java: synchronized**

```java
class SafeCounter {
    private int value = 0;

    synchronized void increment() {
        value++;
    }
}
```

**Go: Mutex explícito**

```go
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (sc *SafeCounter) Increment() {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.value++
}
```

**Ventaja Go:** Explícito, control fino, composable.

---

## 24.3 RWMutex: Locks de Lectura/Escritura

### 24.3.1 Problema: Múltiples Lectores

```
Con Mutex regular:
├─ Lector A adquiere lock
├─ Lector B espera (ineficiente, no hay contención)
└─ Lector C espera

Con RWMutex:
├─ Lector A adquiere read lock
├─ Lector B adquiere read lock (SIMULTÁNEO)
├─ Lector C adquiere read lock (SIMULTÁNEO)
└─ Escritor D espera
```

### 24.3.2 API de RWMutex

```go
type RWMutex struct {
    // Sin campos exportados
}

// Locks de lectura
func (rw *RWMutex) RLock()   // Adquirir read lock
func (rw *RWMutex) RUnlock() // Liberar read lock

// Locks de escritura
func (rw *RWMutex) Lock()    // Adquirir write lock (exclusivo)
func (rw *RWMutex) Unlock()  // Liberar write lock
```

### 24.3.3 Patrón: Cache con Lecturas Concurrentes

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()           // Múltiples lectores simultáneamente
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()            // Escritor exclusivo
    defer c.mu.Unlock()
    c.items[key] = value
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()            // Escritor exclusivo
    defer c.mu.Unlock()
    delete(c.items, key)
}
```

### 24.3.4 Starvation: Hambre de Escritores

```
PROBLEMA: Si hay muchos lectores, el escritor podría nunca obtener lock

[Lector 1] [Lector 2] [Lector 3]  ← Todos con RLock
   Writer espera...
   [Lector 4 llega] ← Se suma al read lock
   Writer sigue esperando...
```

**Go 1.9+** resuelve esto: Una vez que un escritor espera, nuevos lectores no pueden adquirir RLock.

### 24.3.5 Uso: Sistema de Configuración

```go
type Config struct {
    mu    sync.RWMutex
    data  map[string]interface{}
}

func (cfg *Config) Get(key string) (interface{}, bool) {
    cfg.mu.RLock()
    defer cfg.mu.RUnlock()
    val, ok := cfg.data[key]
    return val, ok
}

func (cfg *Config) Set(key string, value interface{}) {
    cfg.mu.Lock()
    defer cfg.mu.Unlock()
    cfg.data[key] = value
}

func (cfg *Config) GetAll() map[string]interface{} {
    cfg.mu.RLock()
    defer cfg.mu.RUnlock()

    // Copiar para evitar race condition
    result := make(map[string]interface{})
    for k, v := range cfg.data {
        result[k] = v
    }
    return result
}
```

### 24.3.6 Downgrade: RLock NO puede convertirse a Lock

```go
// ❌ INCORRECTO: Intentar upgrade
func (c *Cache) GetOrCompute(key string, fn func() interface{}) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if val, ok := c.items[key]; ok {
        return val
    }

    // ❌ Deadlock: Intentar Lock con RLock ya tenido
    c.mu.Lock()
    result := fn()
    c.items[key] = result
    c.mu.Unlock()
    return result
}

// ✓ CORRECTO: Liberar y readquirir
func (c *Cache) GetOrCompute(key string, fn func() interface{}) interface{} {
    c.mu.RLock()
    if val, ok := c.items[key]; ok {
        defer c.mu.RUnlock()
        return val
    }
    c.mu.RUnlock()

    c.mu.Lock()
    defer c.mu.Unlock()

    // Recomprobar después de readquirir
    if val, ok := c.items[key]; ok {
        return val
    }

    result := fn()
    c.items[key] = result
    return result
}
```

---

## 24.4 WaitGroup: Coordinación de Goroutines

### 24.4.1 Concepto: Esperar a N Tareas

```
WaitGroup es un "semáforo de contador":

Add(n) → contador = n
    ↓
Done() → contador-- (N veces)
    ↓
Wait() → bloquea hasta contador == 0
```

### 24.4.2 API de WaitGroup

```go
type WaitGroup struct {
    // Sin campos exportados
}

func (wg *WaitGroup) Add(delta int)    // Sumar al contador
func (wg *WaitGroup) Done()            // Decrementar contador (==Add(-1))
func (wg *WaitGroup) Wait()            // Esperar hasta contador == 0
```

### 24.4.3 Patrón: Fan-Out/Fan-In

```go
func main() {
    var wg sync.WaitGroup
    results := make([]int, 10)

    // Fan-Out: Lanzar 10 goroutines
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            results[index] = expensiveComputation(index)
        }(i)
    }

    // Fan-In: Esperar a que todas terminen
    wg.Wait()

    fmt.Println(results)
}
```

### 24.4.4 Errors: Patrón para Recolectar Errores

```go
type Task struct {
    ID   int
    Err  error
}

func ProcessTasks(tasks []int) []error {
    var wg sync.WaitGroup
    results := make([]Task, len(tasks))

    for i, task := range tasks {
        wg.Add(1)
        go func(index, taskID int) {
            defer wg.Done()
            if err := executeTask(taskID); err != nil {
                results[index].Err = err
            }
        }(i, task)
    }

    wg.Wait()

    // Recolectar solo errores
    var errs []error
    for _, result := range results {
        if result.Err != nil {
            errs = append(errs, result.Err)
        }
    }
    return errs
}
```

### 24.4.5 Patrón: Pool de Workers

```go
func WorkerPool(jobCount int, workers int) {
    var wg sync.WaitGroup
    jobs := make(chan int, jobCount)

    // Lanzar workers
    for w := 0; w < workers; w++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for job := range jobs {
                fmt.Printf("Worker %d: procesando job %d\n", id, job)
                time.Sleep(time.Second)
            }
        }(w)
    }

    // Enviar jobs
    for j := 0; j < jobCount; j++ {
        jobs <- j
    }
    close(jobs)  // Señal para que workers terminen

    wg.Wait()
}
```

### 24.4.6 Error Común: Add después de Wait

```go
// ❌ PANIC: Add después de que el contador llegó a 0
func BAD() {
    var wg sync.WaitGroup
    wg.Wait()  // contador = 0
    wg.Add(1)  // ❌ PANIC!
}

// ✓ CORRECTO: Asegurar que Add se hace antes de Wait
func GOOD() {
    var wg sync.WaitGroup
    wg.Add(1)
    wg.Wait()  // ❌ Espera que Done() se llame
}
```

### 24.4.7 Context Cancellation Pattern

```go
func ProcessWithCancel(ctx context.Context, count int) {
    var wg sync.WaitGroup
    done := make(chan struct{})

    for i := 0; i < count; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            select {
            case <-ctx.Done():
                fmt.Println("Task", id, "cancelada")
                return
            case <-done:
                fmt.Println("Task", id, "completada")
                return
            }
        }(i)
    }

    // Esperar o timeout
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-ctx.Done():
        fmt.Println("Contexto cancelado")
    case <-done:
        fmt.Println("Todas las tareas completadas")
    }
}
```

---

## 24.5 Once: Ejecución Una Sola Vez

### 24.5.1 Concepto: Inicialización Lazy

```
Problema: Inicializar singleton de forma thread-safe

var instance *Database
func GetDB() *Database {
    if instance == nil {  // ❌ Race condition!
        instance = NewDatabase()
    }
    return instance
}

Solución: sync.Once
```

### 24.5.2 API de Once

```go
type Once struct {
    // Sin campos exportados
}

func (o *Once) Do(f func())  // Ejecutar f una sola vez
```

### 24.5.3 Patrón: Singleton Seguro

```go
type Database struct {
    conn string
}

var (
    db   *Database
    once sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        db = &Database{conn: "localhost:5432"}
    })
    return db
}

func main() {
    // Todas estas llamadas ejecutan la inicialización una sola vez
    d1 := GetDatabase()
    d2 := GetDatabase()
    d3 := GetDatabase()

    fmt.Println(d1 == d2 && d2 == d3)  // true
}
```

### 24.5.4 Patrón: Logger Global

```go
type Logger struct {
    level string
    out   io.Writer
}

var (
    log  *Logger
    once sync.Once
)

func InitLogger(level string, out io.Writer) {
    once.Do(func() {
        log = &Logger{level: level, out: out}
    })
}

func GetLogger() *Logger {
    if log == nil {
        InitLogger("info", os.Stderr)
    }
    return log
}
```

### 24.5.5 Patrón: Cleanup Una Sola Vez

```go
type Resource struct {
    closeOnce sync.Once
    file      *os.File
}

func (r *Resource) Close() error {
    var err error
    r.closeOnce.Do(func() {
        // Se ejecuta una sola vez, sin importar cuántas veces se llame Close()
        err = r.file.Close()
    })
    return err
}

func main() {
    res := &Resource{file: /* ... */}
    res.Close()
    res.Close()  // No hace nada, no causa error
    res.Close()  // No hace nada
}
```

### 24.5.6 Características de Once

**1. Thread-safe por defecto**

```go
var once sync.Once
for i := 0; i < 100; i++ {
    go func() {
        once.Do(func() {
            fmt.Println("Se ejecuta una sola vez")
        })
    }()
}
```

**2. La función se ejecuta ANTES de que Do() retorne**

```go
var once sync.Once
var result string

once.Do(func() {
    result = "inicializado"
})

fmt.Println(result)  // "inicializado" garantizado
```

**3. Panic en la función es visible**

```go
var once sync.Once

once.Do(func() {
    panic("Error en inicialización")  // Panic se propaga
})
```

---

## 24.6 Cond: Condition Variables

### 24.6.1 Concepto: Esperar Condiciones

```
Problema: Goroutine A espera a que evento X ocurra,
          Goroutine B señala cuando X ocurrió

Mutex solo: Busy waiting (ineficiente)
Cond: Esperar eficientemente
```

### 24.6.2 API de Cond

```go
type Cond struct {
    L sync.Locker  // Usually *Mutex
}

func NewCond(l sync.Locker) *Cond

func (c *Cond) Wait()       // Liberar lock, esperar señal, readquirir lock
func (c *Cond) Signal()     // Despertar una goroutine esperando
func (c *Cond) Broadcast()  // Despertar TODAS las goroutines esperando
```

### 24.6.3 Patrón: Productor/Consumidor

```go
type Buffer struct {
    mu     sync.Mutex
    items  []interface{}
    notFull *sync.Cond
    notEmpty *sync.Cond
}

func NewBuffer() *Buffer {
    mu := &sync.Mutex{}
    return &Buffer{
        mu:       mu,
        items:    make([]interface{}, 0, 10),
        notFull:  sync.NewCond(mu),
        notEmpty: sync.NewCond(mu),
    }
}

func (b *Buffer) Put(item interface{}) {
    b.mu.Lock()
    defer b.mu.Unlock()

    // Esperar mientras está lleno
    for len(b.items) >= 10 {
        b.notFull.Wait()  // Libera lock, espera, readquiere lock
    }

    b.items = append(b.items, item)
    b.notEmpty.Signal()  // Despertar un consumidor
}

func (b *Buffer) Get() interface{} {
    b.mu.Lock()
    defer b.mu.Unlock()

    // Esperar mientras está vacío
    for len(b.items) == 0 {
        b.notEmpty.Wait()
    }

    item := b.items[0]
    b.items = b.items[1:]
    b.notFull.Signal()  // Despertar un productor

    return item
}
```

### 24.6.4 Patrón: Barrier (Sincronización de Múltiples Goroutines)

```go
type Barrier struct {
    mu    sync.Mutex
    cond  *sync.Cond
    count int
    total int
}

func NewBarrier(n int) *Barrier {
    b := &Barrier{total: n}
    b.cond = sync.NewCond(&b.mu)
    return b
}

func (b *Barrier) Wait() {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.count++
    if b.count >= b.total {
        // Último a llegar despierta a todos
        b.cond.Broadcast()
    } else {
        // Esperar a que todos lleguen
        b.cond.Wait()
    }
}

func main() {
    barrier := NewBarrier(5)

    for i := 0; i < 5; i++ {
        go func(id int) {
            fmt.Printf("Goroutine %d: llegó\n", id)
            barrier.Wait()
            fmt.Printf("Goroutine %d: continuando\n", id)
        }(i)
    }

    time.Sleep(time.Second)
}

// Output (todas esperan, luego todas avanzan):
// Goroutine 0: llegó
// Goroutine 1: llegó
// Goroutine 2: llegó
// Goroutine 3: llegó
// Goroutine 4: llegó
// Goroutine 0: continuando
// ...
```

### 24.6.5 Error Común: Olvida Mutex

```go
// ❌ INCORRECTO: Condition variable sin mutex
cond := sync.NewCond(&sync.Mutex{})
cond.Wait()  // ✓ Funciona técnicamente
cond.Wait()  // ❌ Panic: hace Unlock de un mutex que no tiene Lock

// ✓ CORRECTO
mu := &sync.Mutex{}
cond := sync.NewCond(mu)
mu.Lock()
defer mu.Unlock()
cond.Wait()  // ✓ Seguro
```

---

## 24.7 Sync.Map: Mapa Concurrente

### 24.7.1 Problema: Map + Mutex

```go
// ❌ Ineficiente para alto contention
type SafeMap struct {
    mu    sync.RWMutex
    items map[string]interface{}
}

func (sm *SafeMap) Get(key string) (interface{}, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    val, ok := sm.items[key]
    return val, ok
}
```

### 24.7.2 Beneficio de sync.Map

- Optimizado para **lecturas concurrentes frecuentes**
- Sin locks explícitos (usa copy-on-write internamente)
- Mejor para workloads de "read-heavy"

### 24.7.3 API de sync.Map

```go
type Map struct{
    // Sin campos exportados
}

func (m *Map) Load(key interface{}) (value interface{}, ok bool)
func (m *Map) Store(key, value interface{})
func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
func (m *Map) Delete(key interface{})
func (m *Map) Range(f func(key, value interface{}) bool)
```

### 24.7.4 Ejemplo: Cache con sync.Map

```go
type Cache struct {
    data sync.Map
}

func (c *Cache) Get(key string) (interface{}, bool) {
    return c.data.Load(key)
}

func (c *Cache) Set(key string, value interface{}) {
    c.data.Store(key, value)
}

func (c *Cache) GetOrSet(key string, fn func() interface{}) interface{} {
    result, _ := c.data.LoadOrStore(key, fn())
    return result
}

func (c *Cache) Delete(key string) {
    c.data.Delete(key)
}

func (c *Cache) Iterate(fn func(key string, value interface{})) {
    c.data.Range(func(key, value interface{}) bool {
        fn(key.(string), value)
        return true  // continuar iterando
    })
}
```

### 24.7.5 Limitaciones de sync.Map

```go
// ❌ NO tiene: Len()
var m sync.Map
// len(m)  // ❌ Compilación error

// ❌ NO tiene: Clear()
m.Range(func(key, value interface{}) bool {
    m.Delete(key)  // Funciona pero es lento
    return true
})

// ❌ NO tiene: Conversión a map
// map := m.ToMap()  // ❌ No existe

// ✓ Workaround: Range
func ToMap(m *sync.Map) map[string]interface{} {
    result := make(map[string]interface{})
    m.Range(func(key, value interface{}) bool {
        result[key.(string)] = value
        return true
    })
    return result
}
```

### 24.7.6 Cuándo Usar sync.Map

**✓ Usa sync.Map cuando:**

- Lecturas muy frecuentes, escrituras raras
- Alto contention (muchas goroutines accediendo)
- Keys no cambian frecuentemente

**✗ Evita sync.Map cuando:**

- Muchas escrituras
- Necesitas Len(), Clear()
- Necesitas mantener orden
- Quieres lock fino para subsecciones

---

## 24.8 Pool: Object Pooling

### 24.8.1 Concepto: Reutilizar Objetos

```
Problema: Crear/destruir muchos objetos consume recursos

Buffer pool:
[Buffer1] [Buffer2] [Buffer3]
    ↓          ↓          ↓
  Get()      Get()      Get()
    ↓          ↓          ↓
  [Empty]    [Empty]    [Empty]
```

### 24.8.2 API de Pool

```go
type Pool struct {
    New func() interface{}  // Función para crear nuevos objetos
}

func (p *Pool) Get() interface{}      // Obtener del pool o crear
func (p *Pool) Put(x interface{})     // Devolver al pool
```

### 24.8.3 Patrón: Buffer Pool

```go
type BytesPool struct {
    pool *sync.Pool
}

func NewBytesPool(size int) *BytesPool {
    return &BytesPool{
        pool: &sync.Pool{
            New: func() interface{} {
                return make([]byte, size)
            },
        },
    }
}

func (bp *BytesPool) Get() []byte {
    return bp.pool.Get().([]byte)
}

func (bp *BytesPool) Put(b []byte) {
    b = b[:0]  // Resetear pero mantener capacidad
    bp.pool.Put(b)
}

func main() {
    pool := NewBytesPool(1024)

    buf1 := pool.Get()
    copy(buf1, []byte("data"))
    pool.Put(buf1)

    buf2 := pool.Get()  // Reutiliza buf1!
    fmt.Println(cap(buf2))  // 1024
}
```

### 24.8.4 Patrón: JSON Encoder Pool

```go
var jsonEncoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(io.Discard)
    },
}

func EncodeJSON(w io.Writer, data interface{}) error {
    enc := jsonEncoderPool.Get().(*json.Encoder)
    defer jsonEncoderPool.Put(enc)

    enc.Reset(w)
    return enc.Encode(data)
}
```

### 24.8.5 GC Assistance: Pool no es Cache

```
IMPORTANTE: sync.Pool puede ser vaciado por GC en cualquier momento

var pool = sync.Pool{
    New: func() interface{} {
        fmt.Println("Creando nuevo objeto")
        return make([]byte, 1024)
    },
}

func demo() {
    buf := pool.Get()
    pool.Put(buf)

    // GC corre
    runtime.GC()

    buf = pool.Get()
    // Puede crear un nuevo objeto si fue recolectado
}
```

### 24.8.6 Patrones Reales

**Caso 1: HTTP Request/Response Bodies**

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return bytes.NewBuffer(make([]byte, 0, 32*1024))
    },
}

func HandleRequest(r *http.Request) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)

    buf.Reset()
    io.Copy(buf, r.Body)
    // usar buf...
}
```

---

## 24.9 Semaphore Pattern: Limitar Concurrencia

### 24.9.1 Concepto: Semáforo con Channel

```
Problema: Limitar a N goroutines simultáneamente

Solución: Channel con capacidad N
```

### 24.9.2 Implementación de Semáforo

```go
type Semaphore struct {
    sem chan struct{}
}

func NewSemaphore(max int) *Semaphore {
    return &Semaphore{
        sem: make(chan struct{}, max),
    }
}

func (s *Semaphore) Acquire() {
    s.sem <- struct{}{}
}

func (s *Semaphore) Release() {
    <-s.sem
}

func (s *Semaphore) AcquireN(n int) {
    for i := 0; i < n; i++ {
        s.Acquire()
    }
}

func (s *Semaphore) ReleaseN(n int) {
    for i := 0; i < n; i++ {
        s.Release()
    }
}
```

### 24.9.3 Patrón: Limitar Conexiones a Base de Datos

```go
type DBPool struct {
    sem *Semaphore
}

func NewDBPool(maxConnections int) *DBPool {
    return &DBPool{
        sem: NewSemaphore(maxConnections),
    }
}

func (dp *DBPool) Query(query string) error {
    dp.sem.Acquire()
    defer dp.sem.Release()

    // Ejecutar query con máximo 'maxConnections' simultáneamente
    return executeQuery(query)
}
```

### 24.9.4 Patrón: Limitar Goroutines

```go
func ProcessLargeDataset(items []interface{}, workers int) {
    sem := NewSemaphore(workers)
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        go func(i interface{}) {
            defer wg.Done()
            sem.Acquire()
            defer sem.Release()

            processItem(i)
        }(item)
    }

    wg.Wait()
}
```

### 24.9.5 Patrón: Timeout en Semáforo

```go
func (s *Semaphore) TryAcquire(timeout time.Duration) bool {
    select {
    case s.sem <- struct{}{}:
        return true
    case <-time.After(timeout):
        return false
    }
}

// Uso
sem := NewSemaphore(5)
if sem.TryAcquire(time.Second) {
    defer sem.Release()
    // hacer algo
} else {
    fmt.Println("Timeout esperando semáforo")
}
```

### 24.9.6 Weighted Semaphore (Semáforo Ponderado)

```go
type WeightedSemaphore struct {
    sem chan int
}

func NewWeightedSemaphore(total int) *WeightedSemaphore {
    return &WeightedSemaphore{
        sem: make(chan int, 1),
    }
    // Inicializar con valor total
}

func (ws *WeightedSemaphore) Acquire(weight int) {
    available := <-ws.sem
    for available < weight {
        select {
        case x := <-ws.sem:
            available += x
        }
    }
    ws.sem <- (available - weight)
}

func (ws *WeightedSemaphore) Release(weight int) {
    current := <-ws.sem
    ws.sem <- (current + weight)
}

// Uso: Limitar ancho de banda o recursos
sem := NewWeightedSemaphore(100)  // 100 unidades totales
sem.Acquire(30)   // Usar 30
defer sem.Release(30)
```

---

## 24.10 Deadlocks y Race Conditions

### 24.10.1 Deadlock: Concepto

```
Deadlock ocurre cuando dos o más goroutines se quedan bloqueadas
indefinidamente esperándose mutuamente.

CONDICIONES PARA DEADLOCK (Todas deben cumplirse):
1. Exclusión Mutua: Recurso solo puede usarse por una goroutine
2. Hold & Wait: Goroutine mantiene recurso mientras espera otro
3. No Preemption: Recurso no puede ser quitado por fuerza
4. Circular Wait: Ciclo de goroutines esperándose

Si rompes UNA condición, no hay deadlock.
```

### 24.10.2 Ejemplo: Deadlock Clásico

```go
// ❌ DEADLOCK
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        <-ch1              // Esperar ch1
        ch2 <- 42          // Enviar a ch2
    }()

    go func() {
        <-ch2              // Esperar ch2
        ch1 <- 42          // Enviar a ch1
    }()

    time.Sleep(time.Second)
    // Ambas goroutines esperando...
}

// ✓ SOLUCIÓN: Cambiar orden
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        <-ch1
        ch2 <- 42
    }()

    go func() {
        <-ch1              // Cambio: leer ch1 en lugar de ch2
        ch2 <- 42
    }()
}
```

### 24.10.3 Deadlock con Mutex: Lock Anidado

```go
// ❌ DEADLOCK
func main() {
    var mu sync.Mutex

    ch := make(chan struct{})

    go func() {
        mu.Lock()
        defer mu.Unlock()

        <-ch  // Espera forever
    }()

    mu.Lock()           // ❌ Deadlock: misma goroutine no puede Lock dos veces
    defer mu.Unlock()
    ch <- struct{}{}
}

// ✓ SOLUCIÓN: No hacer lock de la misma goroutine dos veces
func main() {
    var mu sync.Mutex

    ch := make(chan struct{})

    go func() {
        mu.Lock()
        defer mu.Unlock()

        <-ch
    }()

    ch <- struct{}{}

    mu.Lock()
    defer mu.Unlock()
}
```

### 24.10.4 Detección: Runtime Deadlock Detection

```go
// Go detecta algunos deadlocks automáticamente
func main() {
    ch := make(chan int)
    <-ch  // ❌ fatal error: all goroutines are asleep - deadlock!
}
```

### 24.10.5 Race Condition: Concepto

```
Race condition ocurre cuando el orden de ejecución de threads
afecta el resultado final.

EJEMPLO: Dos threads incrementan un contador
Thread 1:  read(x)=0, increment to 1, write(x)=1
Thread 2:  read(x)=0, increment to 1, write(x)=1

Resultado: x=1 (debería ser 2)
```

### 24.10.6 Detector de Race: -race

```bash
# Compilar con race detector
go build -race .
go test -race ./...

# Runtime: Detecta accesos concurrentes a la misma variable
```

### 24.10.7 Ejemplo: Race Condition

```go
// ❌ RACE CONDITION
var counter int

func main() {
    go func() {
        for i := 0; i < 1000; i++ {
            counter++  // Read-modify-write no atómico
        }
    }()

    go func() {
        for i := 0; i < 1000; i++ {
            counter++
        }
    }()

    time.Sleep(time.Second)
    fmt.Println(counter)  // Puede ser < 2000
}

// Ejecutar:
// go run -race file.go
// ==================
// WARNING: DATA RACE
// Read at ...
// Previous write at ...
```

### 24.10.8 Soluciones: Race Condition

```go
// Solución 1: Mutex
var (
    counter int
    mu      sync.Mutex
)

func main() {
    go func() {
        for i := 0; i < 1000; i++ {
            mu.Lock()
            counter++
            mu.Unlock()
        }
    }()
}

// Solución 2: Atomic
var counter int64

func main() {
    go func() {
        for i := 0; i < 1000; i++ {
            atomic.AddInt64(&counter, 1)
        }
    }()
}

// Solución 3: Channel
func main() {
    count := make(chan int)
    go func() {
        c := 0
        for v := range count {
            c += v
        }
    }()
}
```

### 24.10.9 Testing para Encontrar Races

```go
// sync_test.go
func TestRaceCondition(t *testing.T) {
    var counter int

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++
        }()
    }

    wg.Wait()

    if counter != 100 {
        t.Errorf("Expected 100, got %d", counter)
    }
}

// Ejecutar:
// go test -race -count=100 ./...
```

---

## 24.11 Buenas Prácticas y Antipatterns

### 24.11.1 Granularidad: Fina vs Gruesa

```go
// ❌ GRANULARIDAD GRUESA: Un mutex para todo
type User struct {
    mu    sync.Mutex
    Name  string
    Email string
    Age   int
    Score float64
}

// ✓ GRANULARIDAD FINA: Locks separados
type User struct {
    nameMu sync.Mutex
    Name   string

    emailMu sync.Mutex
    Email   string

    ageMu sync.Mutex
    Age   int

    scoreMu sync.Mutex
    Score   float64
}

// ✓ MEJOR: Locks solo donde hay contention real
type User struct {
    mu    sync.RWMutex
    Name  string
    Email string

    scoreMu sync.Mutex
    Score   float64  // Acceso frecuente
}
```

### 24.11.2 Hold Time: Minimizar Tiempo bajo Lock

```go
// ❌ ANTI-PATRÓN: Operación lenta bajo lock
func (c *Cache) ProcessItem(key string, fn func(interface{}) interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    item := c.items[key]
    result := fn(item)  // Puede ser lento
    c.items[key] = result
}

// ✓ PATRÓN: Copiar, soltar, procesar, guardar
func (c *Cache) ProcessItem(key string, fn func(interface{}) interface{}) {
    var item interface{}

    // Hold time: mínimo
    c.mu.Lock()
    item = c.items[key]
    c.mu.Unlock()

    result := fn(item)  // Fuera del lock

    c.mu.Lock()
    c.items[key] = result
    c.mu.Unlock()
}
```

### 24.11.3 Orden de Locks: Prevenir Deadlocks

```go
// ❌ RIESGO DE DEADLOCK
type Account struct {
    mu      sync.Mutex
    balance float64
}

func Transfer(from, to *Account, amount float64) {
    from.mu.Lock()
    defer from.mu.Unlock()

    to.mu.Lock()          // ¿Y si otro thread hace Transfer(to, from)?
    defer to.mu.Unlock()

    from.balance -= amount
    to.balance += amount
}

// ✓ SOLUCIÓN: Ordenar por ID o pointer
func Transfer(from, to *Account, amount float64) {
    if uintptr(unsafe.Pointer(from)) > uintptr(unsafe.Pointer(to)) {
        from, to = to, from
    }

    from.mu.Lock()
    defer from.mu.Unlock()

    to.mu.Lock()
    defer to.mu.Unlock()

    from.balance -= amount
    to.balance += amount
}
```

### 24.11.4 Copy-on-Write para Lecturas Frecuentes

```go
// ❌ Ineficiente: Lock para cada lectura
type Config struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (c *Config) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

// ✓ PATRÓN: Copiar config entera (si es pequeña)
type Config struct {
    mu     sync.Mutex
    atomic atomic.Value  // Contiene snapshot de config
}

type configSnapshot map[string]interface{}

func (c *Config) Get(key string) (interface{}, bool) {
    snap := c.atomic.Load().(configSnapshot)
    val, ok := snap[key]
    return val, ok
}

func (c *Config) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    snap := c.atomic.Load().(configSnapshot)
    newSnap := make(configSnapshot)

    for k, v := range snap {
        newSnap[k] = v
    }
    newSnap[key] = value

    c.atomic.Store(newSnap)
}
```

### 24.11.5 Evitar Locks: Diseño Inmutable

```go
// ❌ Requiere locks para cada acceso
type Person struct {
    mu    sync.Mutex
    Name  string
    Age   int
}

// ✓ Diseño inmutable: Sin locks necesarios
type Person struct {
    Name string
    Age  int
}

// Para actualizaciones, crear nueva instancia
func (p *Person) WithName(name string) *Person {
    return &Person{
        Name: name,
        Age:  p.Age,
    }
}

// Compartir entre goroutines sin lock
p1 := &Person{"Alice", 30}
p2 := p1  // Seguro, es read-only
```

### 24.11.6 Channels vs Mutex: Matriz de Decisión

| Caso | Preference | Razón |
|------|-----------|-------|
| Comunicar datos | Channels | Semantics clara |
| Compartir estado | Mutex | No tiene flujo de datos |
| Múltiples readers | RWMutex | Mejor performance |
| Inicializar once | Once | Más simple que channels |
| Limitar goroutines | Semaphore | Más expresivo |
| Esperar condición | Cond | Eficiente vs polling |

### 24.11.7 Documentación de Locks

```go
type CacheEntry struct {
    // mu protege el acceso a value y expiry
    // No protege la creación de entry (se asume una sola vez)
    mu     sync.Mutex
    value  interface{}
    expiry time.Time
}

// Get devuelve el valor si no ha expirado.
// Es seguro de llamar desde múltiples goroutines.
func (e *CacheEntry) Get() (interface{}, bool) {
    e.mu.Lock()
    defer e.mu.Unlock()

    if time.Now().After(e.expiry) {
        return nil, false
    }
    return e.value, true
}

// Set actualiza el valor y expiry.
// No es seguro de llamar concurrentemente con Get,
// aunque ambos usan el mismo mutex internamente.
// Requiere que el caller sincronice acceso.
func (e *CacheEntry) Set(value interface{}, duration time.Duration) {
    e.mu.Lock()
    defer e.mu.Unlock()

    e.value = value
    e.expiry = time.Now().Add(duration)
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Contador Seguro

Implementa un tipo `SafeCounter` que permite incrementar y obtener el valor de forma segura desde múltiples goroutines.

```go
// Requisitos:
// - Implementar SafeCounter con Mutex
// - Método Inc() para incrementar
// - Método Value() para obtener valor actual
// - Main que lance 1000 goroutines incrementando
// - Verificar que el resultado final sea 1000

// Salida esperada:
// Final counter: 1000
```

**Pista:** Usa `defer` para garantizar Unlock.

---

### Ejercicio 2: Lectura Concurrente con RWMutex

Implementa un sistema de configuración que permite múltiples lectores concurrentes pero escritores exclusivos.

```go
// Requisitos:
// - Tipo Config con RWMutex
// - Métodos Get(key) y Set(key, value)
// - Lanzar 100 goroutines leyendo concurrentemente
// - Ocasionalmente escribir nuevos valores
// - Verificar consistencia (no hay race conditions)

// Salida esperada:
// Reads: 1000, Writes: 50
// No warnings de race detector
```

**Pista:** `go test -race ./...` para verificar.

---

### Ejercicio 3: Inicialización de Singleton con Once

Implementa un singleton thread-safe usando `sync.Once`.

```go
// Requisitos:
// - Tipo Database con conexión "simulada"
// - Usar sync.Once para inicializar una sola vez
// - Lanzar 100 goroutines llamando GetDatabase()
// - Verificar que todas reciben la misma instancia

// Salida esperada:
// Database connection established (una sola vez)
// All goroutines received the same instance: true
```

**Pista:** Usa variables globales para el singleton.

---

### Ejercicio 4: Barrier con Condition Variables

Implementa un barrier que sincroniza N goroutines.

```go
// Requisitos:
// - Tipo Barrier con sync.Cond
// - Método Wait() que bloquea hasta que N goroutines lleguen
// - Lanzar 10 goroutines, esperar en Barrier
// - Imprimir mensajes ordenados

// Salida esperada:
// Worker 0: iniciado
// Worker 1: iniciado
// ...
// (todas esperan)
// Worker 0: pasó barrier
// Worker 1: pasó barrier
// ...
```

**Pista:** Usa Broadcast() para despertar a todos cuando llegan N.

---

### Ejercicio 5: Rate Limiter con Semáforo

Implementa un rate limiter que permite máximo N peticiones simultáneamente.

```go
// Requisitos:
// - Tipo RateLimiter con semáforo basado en channel
// - Método Allow() que adquiere permiso o espera
// - Lanzar 20 goroutines, solo 5 ejecutan simultáneamente
// - Medir tiempo y concurrencia

// Salida esperada:
// Request 0-4: ejecutando
// Request 5: esperando
// ...
// (máximo 5 simultáneamente)
```

**Pista:** Usa channel con capacidad como semáforo.

---

## Conclusión

El paquete `sync` proporciona herramientas poderosas para sincronización de bajo nivel. Sin embargo:

1. **Prefiere channels** cuando sea posible (comunicación clara)
2. **Usa Mutex/RWMutex** para proteger datos compartidos
3. **WaitGroup** para coordinar múltiples goroutines
4. **Once** para inicialización lazy thread-safe
5. **Cond** para esperar condiciones específicas
6. **sync.Map** para caches read-heavy
7. **Pool** para reutilizar objetos

**Regla de Oro:**
> "Las primitivas de sincronización son herramientas poderosas pero peligrosas. Usa channels primero, sync solo cuando sea claramente más eficiente."

---

## Referencias Internas

- Capítulo 20: Goroutines y Concurrencia
- Capítulo 21: Channels
- Capítulo 22: Patrones de Concurrencia
- Capítulo 23: Context

## Recursos Externos

- <https://golang.org/pkg/sync/>
- <https://golang.org/doc/effective_go#concurrency>
- <https://www.ardanlabs.com/blog/2015/01/race-detector.html>

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/24-sync-package/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/24-sync-package):

```bash
cd examples/24-sync-package
go run .
```
