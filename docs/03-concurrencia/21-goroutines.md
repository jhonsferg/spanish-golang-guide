# Capítulo 21: Goroutines - Concurrencia en Go

## Índice

1. [¿Qué es una Goroutine?](#211-qué-es-una-goroutine)
2. [Crear Goroutines](#212-crear-goroutines)
3. [Goroutines y Main](#213-goroutines-y-main)
4. [Ciclo de Vida de Goroutines](#214-ciclo-de-vida-de-goroutines)
5. [Goroutines Concurrentes](#215-goroutines-concurrentes)
6. [WaitGroup - Sincronización Básica](#216-waitgroup---sincronización-básica)
7. [Context - Control Avanzado](#217-context---control-avanzado)
8. [Goroutine Leaks](#218-goroutine-leaks)
9. [Debugging de Goroutines](#219-debugging-de-goroutines)
10. [Best Practices](#2110-best-practices)
11. [Buenas Prácticas y Antipatrones](#2111-buenas-prácticas-y-antipatrones)

---

## 21.1 ¿Qué es una Goroutine?

### 21.1.1 Concepto Fundamental

Una **goroutine** es la abstracción de Go para concurrencia ligera. Es una función que se ejecuta de forma concurrente con otras goroutines, pero con un costo de recursos significativamente menor que los hilos del sistema operativo (threads).

**Definición clave:**
> Una goroutine es una función que se ejecuta de forma concurrente e independiente con otras goroutines dentro del mismo programa, gestionada por el runtime de Go en lugar del sistema operativo.

Características fundamentales:

- **Ligeras**: Cuesta ~2KB de memoria (vs ~2MB de un thread del SO)
- **Concurrentes**: Se ejecutan simultáneamente
- **Gestionadas por Go**: El runtime, no el SO, controla su ejecución
- **No paralelas por defecto**: Muchas goroutines pueden ejecutarse en un único núcleo

### 21.1.2 Concurrencia vs Paralelismo

Es crucial entender la diferencia fundamental:

**Concurrencia:**

- Múltiples tareas se interleavan
- Una estructura de su ejecución
- "Hacer progreso en múltiples tareas" (no necesariamente simultáneamente)
- Resolución de problemas complejos

**Paralelismo:**

- Múltiples tareas se ejecutan en múltiples núcleos simultáneamente
- Propiedad del hardware
- "Hacer cosas simultáneamente"
- Mejora de rendimiento

**Go permite ambas:**

```
Concurrencia sin Paralelismo:
├─ Goroutine A ─┐
├─ Goroutine B ─┼─ CPU Core 1
├─ Goroutine C ─┘

Concurrencia CON Paralelismo:
├─ Goroutine A ─ CPU Core 1
├─ Goroutine B ─ CPU Core 2
└─ Goroutine C ─ CPU Core 3
```

### 21.1.3 Green Threads vs OS Threads

Go usa **green threads** (hilos verdes), no threads del sistema operativo:

| Aspecto | OS Thread | Green Thread (Goroutine) |
|--------|-----------|-------------------------|
| Memoria | ~2 MB | ~2 KB |
| Creación | Costosa (~1000 threads/core) | Económica (~100,000s) |
| Context Switch | Gestionado por SO | Gestionado por Go runtime |
| Scheduling | Preemptivo (SO decide) | Cooperativo + Preemptivo |
| Overhead | Alto | Muy bajo |
| Escalabilidad | Limitada | Excelente |

### 21.1.4 Ventajas de las Goroutines

**1. Escalabilidad**

```
Python/Java (OS threads): ~1000 threads simultáneos
Go (goroutines): ~100,000+ goroutines simultáneas
```

**2. Simplicidad de Sintaxis**
Go hace que la programación concurrente sea simple y elegante:

```go
// Java - verboso y complejo
Thread t = new Thread(() -> {
    // código
});
t.start();
t.join();

// Go - simple y directo
go funcionConcurrente()
```

**3. Gestión de Recursos**

- No hay que crear pools de threads
- No hay que gestionar límites
- El runtime optimiza automáticamente

**4. Mejor Composabilidad**

- Las goroutines se componen bien
- Canales para comunicación
- Context para cancelación

### 21.1.5 Casos de Uso Ideales

Goroutines son perfectas para:

1. **Operaciones I/O concurrentes**
   - Solicitudes HTTP simultáneas
   - Lectura/escritura de múltiples archivos
   - Operaciones de base de datos

2. **Servidores de red**
   - Cada conexión en una goroutine
   - Manejo de miles de clientes

3. **Procesamiento de eventos**
   - Event loops
   - Consumidores de eventos

4. **Tareas de larga duración**
   - Cálculos en paralelo
   - Procesamiento batch

### 21.1.6 M:N Scheduling - El Secreto de Go

Go implementa un modelo de scheduling **M:N**:

- **M** goroutines
- **N** OS threads (típicamente igual al número de CPUs)

```
Goroutines (M) ─┐
                ├─ Work Queue ─ Scheduler ─┐
                └───────────────┘          │
                                           ├─ OS Thread 1
                                           ├─ OS Thread 2
                                           ├─ OS Thread 3
                                           └─ OS Thread N
```

**Cómo funciona:**

1. Go crea un pequeño número de OS threads (~número de CPU cores)
2. Las M goroutines se distribuyen entre estos N threads
3. Cuando una goroutine se bloquea (I/O), el thread es liberado
4. Otra goroutine toma su lugar

Esto permite:

- Miles de goroutines con solo algunos threads
- Bloqueo sin desperdicio de recursos
- Eficiencia máxima en operaciones I/O

```go
// Ejemplo del scheduling M:N

// Cómo verlo en acción:
fmt.Println(runtime.NumGoroutine())  // Número de goroutines activas
fmt.Println(runtime.NumCPU())        // Número de CPUs
fmt.Println(runtime.GOMAXPROCS(-1))  // Max OS threads usados
```

---

## 21.2 Crear Goroutines

### 21.2.1 La Palabra Clave `go`

La forma más simple de crear una goroutine es preceder una llamada a función con la palabra clave `go`:

```go
package main

import "fmt"

func tarea() {
    fmt.Println("Ejecutándose en una goroutine")
}

func main() {
    // Llamada normal - síncrona
    tarea() // Espera a que termine

    // Llamada como goroutine - asíncrona
    go tarea() // No espera

    // Sin hacer nada aquí, main termina y el programa cierra
}
```

### 21.2.2 Diferencia: Llamada Normal vs Goroutine

```go
package main

import "fmt"

func operacion(num int) {
    fmt.Printf("Operación %d comenzó\n", num)
    time.Sleep(time.Second)
    fmt.Printf("Operación %d terminó\n", num)
}

func main() {
    // SÍNCRONO: espera a cada una
    operacion(1) // 1s
    operacion(2) // 1s
    // Total: 2 segundos

    // ASÍNCRONO: no espera
    go operacion(1)  // Inicia pero no espera
    go operacion(2)  // Inicia pero no espera
    // Sin espera explícita, main termina inmediatamente
}
```

### 21.2.3 Crear Goroutine con Funciones Anónimas

Las funciones anónimas son útiles para goroutines rápidas:

```go
package main

import "fmt"

func main() {
    // Goroutine con función anónima
    go func() {
        fmt.Println("Hola desde goroutine")
    }()

    // Con parámetros
    go func(nombre string) {
        fmt.Printf("Hola, %s\n", nombre)
    }("Carlos")

    // Con acceso a variables externas
    contador := 0
    go func() {
        contador++
        fmt.Printf("Contador: %d\n", contador)
    }()
}
```

### 21.2.4 Parámetros en Goroutines

**Importante:** Los parámetros se evalúan en el momento de crear la goroutine:

```go
package main

import "fmt"

func saludar(nombre string) {
    fmt.Printf("Hola, %s\n", nombre)
}

func main() {
    nombres := []string{"Alice", "Bob", "Charlie"}

    // ❌ INCORRECTO: Todos imprimirán "Charlie"
    for _, nombre := range nombres {
        go func() {
            fmt.Println(nombre) // Captura la variable, no el valor
        }()
    }

    // ✓ CORRECTO: Pasar como parámetro
    for _, nombre := range nombres {
        go func(n string) {
            fmt.Println(n) // Cada goroutine recibe su valor
        }(nombre)
    }

    // ✓ ALTERNATIVO: Copiar a variable local
    for _, nombre := range nombres {
        n := nombre // Copiar
        go func() {
            fmt.Println(n)
        }()
    }
}
```

### 21.2.5 Requisitos Sintácticos

Para crear una goroutine válida:

```go
// ✓ Válido: función sin retorno
go printf("mensaje")

// ✓ Válido: función con retorno (se ignora)
go strlen("texto")

// ✓ Válido: función anónima
go func() {}()

// ✓ Válido: método
go miStruct.metodo()

// ✓ Válido: closure
go func(x int) { fmt.Println(x) }(42)

// ❌ Inválido: no es una llamada a función
go 42
go x + y
go variable
```

---

## 21.3 Goroutines y Main

### 21.3.1 El Problema Fundamental

Cuando creas goroutines, la función `main` no espera a que terminen:

```go
package main

import (
    "fmt"
    "time"
)

func tarea() {
    time.Sleep(2 * time.Second)
    fmt.Println("Tarea completada")
}

func main() {
    go tarea()

    // main termina aquí
    // El programa cierra sin esperar la goroutine
    // Output: (vacío)
}
```

**¿Por qué?** Porque `main` es en sí misma una goroutine. Cuando termina, el programa cierra, incluso si hay otras goroutines activas.

### 21.3.2 Problema: Goroutines Huérfanas

```
main goroutine
├─ time.Sleep(2s) ← ESPERA 2 SEGUNDOS AQUÍ
└─ return ← TERMINA

Goroutine 1
└─ time.Sleep(2s) ← SIGUE EJECUTÁNDOSE
   fmt.Println() ← NUNCA LLEGA
```

### 21.3.3 Solución 1: Sleep en Main

La forma más primitiva (¡NO RECOMENDADA!):

```go
package main

import (
    "fmt"
    "time"
)

func tarea() {
    fmt.Println("Tarea iniciada")
    time.Sleep(1 * time.Second)
    fmt.Println("Tarea completada")
}

func main() {
    go tarea()

    // Esperar lo suficiente para que la goroutine termine
    time.Sleep(2 * time.Second)

    // Output:
    // Tarea iniciada
    // Tarea completada
}
```

**Problemas:**

- Tiempo arbitrario (¿2 segundos? ¿5?)
- Desperdicia CPU esperando
- No es confiable
- Código frágil

### 21.3.4 Solución 2: WaitGroup (Recomendado)

`sync.WaitGroup` es el mecanismo correcto:

```go
package main

import (
    "fmt"
    "sync"
)

func tarea(wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("Tarea completada")
}

func main() {
    var wg sync.WaitGroup

    wg.Add(1)          // Aumentar contador
    go tarea(&wg)      // Pasar WaitGroup

    wg.Wait()          // Esperar a que todas terminen

    fmt.Println("Todas las goroutines completadas")
}
```

### 21.3.5 Flujo de Ejecución

```
main goroutine
├─ wg.Add(1) ← contador = 1
├─ go tarea(&wg) ← inicia goroutine
│
├─ wg.Wait() ← BLOQUEA aquí
│
├─ Goroutine ejecuta tarea()
│  └─ wg.Done() ← contador = 0, desbloquea Wait()
│
└─ continue después de Wait() ← REANUDA aquí
   fmt.Println("Completadas")
```

### 21.3.6 Sincronización con Channels

Otra opción es usar channels (veremos más adelante):

```go
package main

import "fmt"

func tarea(done chan bool) {
    fmt.Println("Tarea ejecutándose")
    done <- true  // Señalizar finalización
}

func main() {
    done := make(chan bool)

    go tarea(done)

    <-done  // Esperar señal

    fmt.Println("Tarea completada")
}
```

---

## 21.4 Ciclo de Vida de Goroutines

### 21.4.1 Estados de una Goroutine

Una goroutine pasa por varios estados durante su existencia:

```
1. CREADA ─┐
           │
           ├─→ 2. RUNNABLE ─→ 3. RUNNING ─┐
           │                               │
           └────────────────────←──────────┤
                                           │
                            4. BLOCKED ←──┤
                                   ↓      │
                                   └──→ 5. RUNNABLE
                                           │
                                    6. TERMINATED ←─
```

**Explicación:**

1. **CREADA**: La goroutine ha sido creada con `go` pero no ha empezado
2. **RUNNABLE**: Está lista para ejecutarse (en la cola de scheduler)
3. **RUNNING**: Se está ejecutando en un OS thread
4. **BLOCKED**: Bloqueada (esperando I/O, channel, lock, etc.)
5. **TERMINATED**: Completó su ejecución

### 21.4.2 Estados en Detalle

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func tarea(id int) {
    fmt.Printf("[%d] Iniciando\n", id)
    time.Sleep(time.Second)
    fmt.Printf("[%d] Completada\n", id)
}

func main() {
    fmt.Printf("Goroutines activas al inicio: %d\n", runtime.NumGoroutine())

    go tarea(1)
    go tarea(2)

    fmt.Printf("Goroutines activas después de go: %d\n", runtime.NumGoroutine())

    time.Sleep(2 * time.Second)

    fmt.Printf("Goroutines activas al final: %d\n", runtime.NumGoroutine())
}

// Output:
// Goroutines activas al inicio: 1
// Goroutines activas después de go: 3
// [1] Iniciando
// [2] Iniciando
// [2] Completada
// [1] Completada
// Goroutines activas al final: 1
```

### 21.4.3 Tiempo de Vida Completo

```go
package main

import (
    "fmt"
    "time"
)

func trabajador(id int, duracion time.Duration) {
    fmt.Printf("[%d] Iniciando (%dms)\n", id, duracion.Milliseconds())

    // CREADA → RUNNING

    // Simular trabajo
    time.Sleep(duracion)

    // BLOCKED (durante Sleep) → RUNNING → TERMINATED

    fmt.Printf("[%d] Completada\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        id := i
        go func() {
            defer wg.Done()
            trabajador(id, time.Duration(id)*100*time.Millisecond)
        }()
    }

    wg.Wait()
    fmt.Println("Todas completadas")
}
```

### 21.4.4 Goroutines Que Nunca Terminan

```go
package main

import (
    "fmt"
    "time"
)

func serverEterno() {
    for {
        fmt.Println("Servidor ejecutando...")
        time.Sleep(time.Second)
    }
    // Nunca llega aquí (a menos que se cancele)
}

func main() {
    go serverEterno()

    time.Sleep(3 * time.Second)

    // El programa termina, serverEterno se detiene abruptamente
    fmt.Println("Main terminando")
}
```

---

## 21.5 Goroutines Concurrentes

### 21.5.1 Múltiples Goroutines Ejecutándose

El poder real de goroutines está en ejecutar múltiples tareas concurrentemente:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func descargarArchivo(nombre string, duracion time.Duration, wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Printf("Descargando: %s\n", nombre)
    time.Sleep(duracion)
    fmt.Printf("✓ %s completado\n", nombre)
}

func main() {
    var wg sync.WaitGroup

    start := time.Now()

    // Descargar 5 archivos
    archivos := []struct {
        nombre   string
        duracion time.Duration
    }{
        {"archivo1.txt", 1 * time.Second},
        {"archivo2.txt", 2 * time.Second},
        {"archivo3.txt", 1 * time.Second},
        {"archivo4.txt", 3 * time.Second},
        {"archivo5.txt", 1 * time.Second},
    }

    for _, archivo := range archivos {
        wg.Add(1)
        a := archivo
        go descargarArchivo(a.nombre, a.duracion, &wg)
    }

    wg.Wait()

    elapsed := time.Since(start)
    fmt.Printf("\nTiempo total: %.2fs\n", elapsed.Seconds())

    // Output:
    // Descargando: archivo1.txt
    // Descargando: archivo2.txt
    // Descargando: archivo3.txt
    // Descargando: archivo4.txt
    // Descargando: archivo5.txt
    // ✓ archivo1.txt completado
    // ✓ archivo3.txt completado
    // ✓ archivo5.txt completado
    // ✓ archivo2.txt completado
    // ✓ archivo4.txt completado
    //
    // Tiempo total: 3.00s (¡No es 8 segundos!)
}
```

**Análisis:**

- Síncrono: 1+2+1+3+1 = 8 segundos
- Concurrente: max(1,2,1,3,1) = 3 segundos
- **Mejora: 2.67x más rápido**

### 21.5.2 Goroutines con Trabajo Variable

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func procesarElemento(id int, wg *sync.WaitGroup) {
    defer wg.Done()

    duracion := time.Duration(rand.Intn(5)) * time.Second
    fmt.Printf("[%d] Procesando (%ds)\n", id, duracion)

    time.Sleep(duracion)

    fmt.Printf("[%d] ✓ Completado\n", id)
}

func main() {
    var wg sync.WaitGroup

    // Procesar 10 elementos concurrentemente
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        id := i
        go procesarElemento(id, &wg)
    }

    wg.Wait()
    fmt.Println("✓ Todos procesados")
}
```

### 21.5.3 Scheduling y Fairness

El scheduler de Go intenta ser justo, pero no garantiza orden:

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func mostrarGoroutineID(id int, wg *sync.WaitGroup) {
    defer wg.Done()

    // Obtener ID de la goroutine actual
    var gid uint64
    gid = uint64(id) // (En Go 1.17+, no hay API directa para goroutine ID)

    fmt.Printf("Goroutine %d en OS thread %d\n", gid, id)
}

func main() {
    var wg sync.WaitGroup
    numGoroutines := 8

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go mostrarGoroutineID(i, &wg)
    }

    wg.Wait()

    fmt.Printf("OS threads usados: %d\n", runtime.NumCPU())
}
```

### 21.5.4 Contención de Recursos

Cuando muchas goroutines acceden a recursos compartidos:

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var mu sync.Mutex
    contador := 0
    var wg sync.WaitGroup

    // 10 goroutines incrementando el contador
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            for j := 0; j < 1000; j++ {
                mu.Lock()
                contador++
                mu.Unlock()
            }
        }()
    }

    wg.Wait()

    fmt.Printf("Contador final: %d\n", contador)
    // Output: 10000 (¡Sin race condition!)
}
```

---

## 21.6 WaitGroup - Sincronización Básica

### 21.6.1 ¿Qué es WaitGroup?

`sync.WaitGroup` es un mecanismo simple de sincronización que permite esperar a que múltiples goroutines terminen:

```go
type WaitGroup struct {
    // Contador interno (no accesible directamente)
}

func (wg *WaitGroup) Add(delta int)    // Incrementar contador
func (wg *WaitGroup) Done()            // Decrementar contador
func (wg *WaitGroup) Wait()            // Bloquear hasta contador = 0
```

### 21.6.2 Funcionamiento Básico

```go
package main

import (
    "fmt"
    "sync"
)

func tarea(id int, wg *sync.WaitGroup) {
    defer wg.Done()  // Decrementar cuando termine

    fmt.Printf("Tarea %d ejecutándose\n", id)
}

func main() {
    var wg sync.WaitGroup

    // Indicar que esperaremos 3 goroutines
    wg.Add(3)

    go tarea(1, &wg)
    go tarea(2, &wg)
    go tarea(3, &wg)

    // Bloquear hasta que las 3 terminen
    wg.Wait()

    fmt.Println("✓ Todas completadas")
}
```

### 21.6.3 Patrón Típico

```go
package main

import (
    "fmt"
    "sync"
)

func procesarElemento(elemento string, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Procesando: %s\n", elemento)
}

func main() {
    var wg sync.WaitGroup
    elementos := []string{"A", "B", "C", "D", "E"}

    // Incrementar contador para cada goroutine
    wg.Add(len(elementos))

    // Crear goroutine para cada elemento
    for _, elem := range elementos {
        go procesarElemento(elem, &wg)
    }

    // Esperar a que todas terminen
    wg.Wait()
}
```

### 21.6.4 Add con Múltiples Goroutines

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    // Opción 1: Agregar una por una
    wg.Add(1)
    wg.Add(1)
    wg.Add(1)

    // Opción 2: Agregar en bloque
    var wg2 sync.WaitGroup
    wg2.Add(3)

    // Opción 3: Agregar dinámicamente en un loop
    var wg3 sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg3.Add(1)
        go func() {
            defer wg3.Done()
        }()
    }
}
```

### 21.6.5 Errores Comunes

**Error 1: No llamar a Done()**

```go
// ❌ INCORRECTO: Deadlock
func main() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        fmt.Println("Goroutine")
        // Olvidó llamar wg.Done()
    }()

    wg.Wait()  // ⏰ Espera eternamente
}
```

**Error 2: Add después de Wait**

```go
// ❌ INCORRECTO: Race condition
func main() {
    var wg sync.WaitGroup

    go func() {
        wg.Add(1)      // Demasiado tarde
        defer wg.Done()
    }()

    wg.Wait()  // Puede terminar antes de Add()
}
```

**Error 3: Add negativo**

```go
// ❌ INCORRECTO: Panic
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    wg.Done()
    wg.Done()  // ⚠️ Panic: sync: negative WaitGroup counter
}
```

### 21.6.6 Patrón Avanzado: WaitGroup Reutilizable

```go
package main

import (
    "fmt"
    "sync"
)

func procesarLote(items []string, numWorkers int) {
    var wg sync.WaitGroup

    // Canal para distribuir trabajo
    trabajos := make(chan string, len(items))

    // Iniciar workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            for trabajo := range trabajos {
                fmt.Printf("Worker %d: %s\n", id, trabajo)
            }
        }(i)
    }

    // Enviar trabajo
    for _, item := range items {
        trabajos <- item
    }
    close(trabajos)

    // Esperar a todos los workers
    wg.Wait()
}

func main() {
    items := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
    procesarLote(items, 3)
}
```

---

## 21.7 Context - Control Avanzado

### 21.7.1 ¿Qué es Context?

`context.Context` es un mecanismo poderoso para pasar información entre goroutines y controlar su ciclo de vida:

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

**Propósitos:**

- Pasar valores entre goroutines
- Establecer timeouts y deadlines
- Cancelación coordinada
- Rastreo de solicitudes (request tracing)

### 21.7.2 Context Raíz: Background y TODO

```go
package main

import (
    "context"
    "fmt"
)

func main() {
    // context.Background()
    // - Contexto raíz
    // - Nunca se cancela
    // - Sin deadline
    // - Usado al inicio del programa

    ctx1 := context.Background()
    fmt.Println("Background context:", ctx1)

    // context.TODO()
    // - Contexto raíz de relleno
    // - Usado cuando no sabes qué context usar
    // - Indica "necesito un context aquí, pero aún no sé cuál"

    ctx2 := context.TODO()
    fmt.Println("TODO context:", ctx2)

    // Generalmente: usa Background() en main
    ctx := context.Background()

    // Chequear si está cancelado
    fmt.Println("Error:", ctx.Err()) // nil (no cancelado)
    fmt.Println("Done channel:", ctx.Done())
}
```

### 21.7.3 Context con Timeout

Limitar el tiempo de ejecución:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func operacionLenta(ctx context.Context) {
    for i := 1; i <= 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Println("Cancelado:", ctx.Err())
            return
        default:
            fmt.Printf("Paso %d\n", i)
            time.Sleep(1 * time.Second)
        }
    }
    fmt.Println("✓ Completado")
}

func main() {
    // Crear context con timeout de 3 segundos
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()  // Liberar recursos

    operacionLenta(ctx)

    // Output:
    // Paso 1
    // Paso 2
    // Paso 3
    // Cancelado: context deadline exceeded
}
```

### 21.7.4 Context con Cancelación

Control manual de cancelación:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func servidor(ctx context.Context, nombre string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: Apagando\n", nombre)
            return
        default:
            fmt.Printf("%s: Ejecutando\n", nombre)
            time.Sleep(1 * time.Second)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    // Iniciar dos servidores
    go servidor(ctx, "API")
    go servidor(ctx, "Cache")

    // Dejar ejecutar 3 segundos
    time.Sleep(3 * time.Second)

    // Cancelar todos
    fmt.Println("Cancelando...")
    cancel()

    time.Sleep(1 * time.Second)
}
```

### 21.7.5 Context con Deadline

Similar a timeout pero con hora exacta:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func tarea(ctx context.Context) {
    deadline, ok := ctx.Deadline()

    if ok {
        fmt.Printf("Deadline: %v\n", deadline)
        fmt.Printf("Tiempo restante: %v\n", time.Until(deadline))
    }

    select {
    case <-ctx.Done():
        fmt.Println("Cancelado")
    case <-time.After(2 * time.Second):
        fmt.Println("Tarea completada")
    }
}

func main() {
    ctx, cancel := context.WithDeadline(
        context.Background(),
        time.Now().Add(5*time.Second),
    )
    defer cancel()

    tarea(ctx)
}
```

### 21.7.6 Context con Valores

Pasar datos entre goroutines:

```go
package main

import (
    "context"
    "fmt"
)

type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
)

func procesarPedido(ctx context.Context) {
    requestID := ctx.Value(requestIDKey)
    userID := ctx.Value(userIDKey)

    fmt.Printf("Pedido %v para usuario %v\n", requestID, userID)
}

func main() {
    ctx := context.Background()

    // Agregar valores
    ctx = context.WithValue(ctx, requestIDKey, "REQ-123")
    ctx = context.WithValue(ctx, userIDKey, "USER-456")

    procesarPedido(ctx)

    // Output: Pedido REQ-123 para usuario USER-456
}
```

### 21.7.7 Árbol de Context

Los contexts forman un árbol jerárquico:

```
Background()
├── WithTimeout(2s)
│   └── WithValue(userID)
│       └── WithCancel()
│
└── WithCancel()
    └── WithValue(requestID)
```

**Comportamiento:**

- Cancelar un context cancela todos sus hijos
- Los valores se heredan hacia abajo
- Los timeouts se heredan

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func trabajador(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: Cancelado\n", id)
            return
        default:
            fmt.Printf("Worker %d: Ejecutando\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // Context raíz
    rootCtx, rootCancel := context.WithCancel(context.Background())

    // Context hijo
    childCtx, childCancel := context.WithCancel(rootCtx)

    go trabajador(childCtx, 1)
    time.Sleep(1 * time.Second)

    // Cancelar solo el hijo
    fmt.Println("Cancelando hijo...")
    childCancel()
    time.Sleep(1 * time.Second)

    // Ambos se cancelan
    fmt.Println("Cancelando raíz...")
    rootCancel()

    time.Sleep(500 * time.Millisecond)
}
```

---

## 21.8 Goroutine Leaks

### 21.8.1 ¿Qué es una Goroutine Leak?

Un **goroutine leak** ocurre cuando una goroutine nunca termina, permanece activa indefinidamente:

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    fmt.Printf("Goroutines al inicio: %d\n", runtime.NumGoroutine())

    // Crear goroutine que nunca termina
    go func() {
        for {
            // Bucle infinito - nunca retorna
        }
    }()

    fmt.Printf("Goroutines después: %d\n", runtime.NumGoroutine())

    // Output:
    // Goroutines al inicio: 1
    // Goroutines después: 2
}
```

### 21.8.2 Causas Comunes

**1. Goroutine Bloqueada en Channel**

```go
// ❌ LEAK: Espera indefinida por dato
func main() {
    ch := make(chan int)
    go func() {
        valor := <-ch  // Espera dato que nunca llega
    }()
}
```

**2. Goroutine en Bucle Infinito Sin Cancelación**

```go
// ❌ LEAK: Bucle infinito sin forma de parar
func main() {
    go func() {
        for {
            // Continúa para siempre
        }
    }()
}
```

**3. Goroutine Esperando Lock**

```go
// ❌ LEAK: Deadlock de mutex
func main() {
    var mu sync.Mutex

    go func() {
        mu.Lock()
        mu.Lock()  // ⏰ Se bloquea aquí para siempre
        mu.Unlock()
        mu.Unlock()
    }()
}
```

**4. Timeout en Operación I/O**

```go
// ❌ LEAK: Conexión que nunca se cierra
func main() {
    go func() {
        conn, _ := net.Dial("tcp", "192.0.2.1:9999")
        // Nunca lee ni cierra
        // Goroutine bloqueada indefinidamente
    }()
}
```

### 21.8.3 Detección de Leaks

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    initialGoroutines := runtime.NumGoroutine()
    fmt.Printf("Goroutines iniciales: %d\n", initialGoroutines)

    // Ejecutar operación
    tarea()

    time.Sleep(100 * time.Millisecond)

    finalGoroutines := runtime.NumGoroutine()
    fmt.Printf("Goroutines finales: %d\n", finalGoroutines)

    if finalGoroutines > initialGoroutines {
        fmt.Printf("⚠️ LEAK DETECTADO: +%d goroutines\n",
            finalGoroutines-initialGoroutines)
    }
}

func tarea() {
    go func() {
        // Simular leak
        select {}  // Se bloquea para siempre
    }()
}
```

### 21.8.4 Prevención: Patrón Seguro

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

func procesador(ctx context.Context, trabajos <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()

    for {
        select {
        case <-ctx.Done():
            // ✓ Respeta cancelación
            return
        case trabajo, ok := <-trabajos:
            if !ok {
                // ✓ Canal cerrado
                return
            }
            fmt.Printf("Procesando: %d\n", trabajo)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    var wg sync.WaitGroup
    trabajos := make(chan int, 10)

    // Iniciar procesador
    wg.Add(1)
    go procesador(ctx, trabajos, &wg)

    // Enviar trabajos
    for i := 0; i < 5; i++ {
        trabajos <- i
    }
    close(trabajos)

    wg.Wait()
}
```

### 21.8.5 Patrón: Goroutine que Respeta Cancelación

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func servidor(ctx context.Context, puerto int) {
    for {
        select {
        case <-ctx.Done():
            // ✓ Cancelación capturada
            fmt.Printf("Servidor %d cerrando: %v\n", puerto, ctx.Err())
            return
        default:
            // ✓ Trabajo normal
            fmt.Printf("Servidor %d ejecutando\n", puerto)
            time.Sleep(time.Second)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go servidor(ctx, 8080)
    go servidor(ctx, 8081)

    time.Sleep(2 * time.Second)

    // ✓ Cancelar todos ordenadamente
    cancel()

    time.Sleep(500 * time.Millisecond)
}
```

---

## 21.9 Debugging de Goroutines

### 21.9.1 Contar Goroutines Activas

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func main() {
    fmt.Printf("Goroutines al inicio: %d\n", runtime.NumGoroutine())

    var wg sync.WaitGroup

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d: %d activas\n", id, runtime.NumGoroutine())
            time.Sleep(time.Second)
        }(i)
    }

    fmt.Printf("Después de crear: %d\n", runtime.NumGoroutine())

    wg.Wait()

    fmt.Printf("Goroutines al final: %d\n", runtime.NumGoroutine())
}
```

### 21.9.2 Stack Traces

Ver el estado de todas las goroutines:

```go
package main

import (
    "fmt"
    "runtime/debug"
)

func main() {
    go func() {
        for {
            // Bucle infinito
        }
    }()

    // Imprimir stack de todas las goroutines
    fmt.Println(debug.Stack())
}
```

**Output típico:**

```
goroutine 1 [runnable]:
main.main()
    /path/to/main.go:13

goroutine 2 [running]:
main.main.func1()
    /path/to/main.go:8
```

### 21.9.3 Profiling de Goroutines

```go
package main

import (
    "fmt"
    "os"
    "runtime"
    "runtime/pprof"
)

func trabajador(id int, done <-chan struct{}) {
    <-done
}

func main() {
    // Crear archivo de perfil
    f, _ := os.Create("goroutine.prof")
    defer f.Close()

    // Crear goroutines
    done := make(chan struct{})
    for i := 0; i < 100; i++ {
        go trabajador(i, done)
    }

    // Escribir perfil de goroutines
    pprof.Lookup("goroutine").WriteTo(f, 0)

    fmt.Printf("Goroutines activas: %d\n", runtime.NumGoroutine())
}
```

### 21.9.4 Trace para Análisis Detallado

```go
package main

import (
    "os"
    "runtime/trace"
    "sync"
    "time"
)

func main() {
    // Crear archivo de trace
    f, _ := os.Create("trace.out")
    defer f.Close()

    trace.Start(f)
    defer trace.Stop()

    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            time.Sleep(time.Millisecond)
        }(i)
    }

    wg.Wait()
}

// Ver con: go tool trace trace.out
```

### 21.9.5 Logging de Goroutines

```go
package main

import (
    "fmt"
    "log"
    "runtime"
    "sync"
    "time"
)

func tareaConLogging(id int, wg *sync.WaitGroup) {
    defer wg.Done()

    gid := runtime.NumGoroutine()  // Aproximado
    log.Printf("[G%d] Iniciando tarea %d", gid, id)

    time.Sleep(time.Second)

    log.Printf("[G%d] Completada tarea %d", gid, id)
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 3; i++ {
        wg.Add(1)
        go tareaConLogging(i, &wg)
    }

    wg.Wait()
    fmt.Printf("Goroutines finales: %d\n", runtime.NumGoroutine())
}
```

---

## 21.10 Best Practices

### 21.10.1 Nombrar Goroutines Lógicamente

```go
package main

import (
    "fmt"
    "sync"
)

func procesadorDatos(id int, items <-chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    for item := range items {
        fmt.Printf("[ProcesadorDatos-%d] %s\n", id, item)
    }
}

func generadorDatos(items chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    datos := []string{"A", "B", "C"}
    for _, d := range datos {
        items <- d
    }
    close(items)
}

func main() {
    var wg sync.WaitGroup
    items := make(chan string, 10)

    // Generador
    wg.Add(1)
    go generadorDatos(items, &wg)

    // Procesadores
    for i := 0; i < 2; i++ {
        wg.Add(1)
        go procesadorDatos(i, items, &wg)
    }

    wg.Wait()
}
```

### 21.10.2 Gestión de Ciclo de Vida Explícita

```go
package main

import (
    "context"
    "fmt"
    "sync"
)

type servicio struct {
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

func (s *servicio) Iniciar() {
    s.wg.Add(1)
    go s.correr()
}

func (s *servicio) correr() {
    defer s.wg.Done()
    fmt.Println("✓ Servicio iniciado")

    <-s.ctx.Done()

    fmt.Println("✗ Servicio detenido")
}

func (s *servicio) Detener() {
    s.cancel()
    s.wg.Wait()
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    srv := &servicio{ctx: ctx, cancel: cancel}
    srv.Iniciar()

    // Usar servicio...

    srv.Detener()  // Espera a que termine
}
```

### 21.10.3 Usar WaitGroup Consistentemente

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    // ✓ BUENA PRÁCTICA: Add antes de go
    tareas := []string{"T1", "T2", "T3"}

    for _, tarea := range tareas {
        wg.Add(1)  // ANTES de crear la goroutine

        t := tarea
        go func() {
            defer wg.Done()
            fmt.Printf("Ejecutando: %s\n", t)
        }()
    }

    wg.Wait()
}
```

### 21.10.4 Documentar Comportamiento Concurrente

```go
package main

// Procesador maneja múltiples solicitudes concurrentemente.
// Es seguro llamarlo desde múltiples goroutines.
// Se debe llamar Cerrar() para detener gracefully.
type Procesador struct {
    trabajos chan struct{}
    cancelar context.CancelFunc
}

// Procesar comienza el procesamiento en una nueva goroutine.
// Retorna inmediatamente.
func (p *Procesador) Procesar() {
    go p.ejecutar()
}

// Cerrar detiene el procesador gracefully.
// Bloquea hasta que el procesamiento termine.
func (p *Procesador) Cerrar() error {
    p.cancelar()
    return nil
}

func (p *Procesador) ejecutar() {
    // Implementación
}
```

### 21.10.5 Error Handling en Goroutines

```go
package main

import (
    "fmt"
    "sync"
)

type Resultado struct {
    ID    int
    Error error
}

func procesarConErrores(id int, resultados chan<- Resultado, wg *sync.WaitGroup) {
    defer wg.Done()

    // Simular operación que puede fallar
    err := error(nil)
    if id%2 == 0 {
        err = fmt.Errorf("error en tarea %d", id)
    }

    resultados <- Resultado{ID: id, Error: err}
}

func main() {
    var wg sync.WaitGroup
    resultados := make(chan Resultado, 5)

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go procesarConErrores(i, resultados, &wg)
    }

    // Goroutine para recolectar resultados
    go func() {
        wg.Wait()
        close(resultados)
    }()

    // Procesar resultados
    for r := range resultados {
        if r.Error != nil {
            fmt.Printf("✗ Tarea %d: %v\n", r.ID, r.Error)
        } else {
            fmt.Printf("✓ Tarea %d exitosa\n", r.ID)
        }
    }
}
```

---

## 21.11 Buenas Prácticas y Antipatrones

### 21.11.1 Comparación: Go vs Otros Lenguajes

**Go Goroutines vs Java Threads:**

```
// JAVA: Verboso y pesado
ExecutorService executor = Executors.newFixedThreadPool(10);
for (int i = 0; i < 1000; i++) {
    executor.submit(() -> {
        // tarea
    });
}
executor.shutdown();
executor.awaitTermination(1, TimeUnit.MINUTES);

// GO: Simple y elegante
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // tarea
    }()
}
wg.Wait()
```

**Escalabilidad:**

- Java: ~1000 threads en un servidor típico
- Go: ~100,000+ goroutines en el mismo servidor

### 21.11.2 Antipatrones Comunes

**Antipatrón 1: Ignorar el Retorno de Goroutine**

```go
// ❌ MALO: No hay forma de saber si la tarea se completó
go funcionImportante()

// ✓ BIEN: Señalizar finalización
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    funcionImportante()
}()
wg.Wait()
```

**Antipatrón 2: Goroutines sin Ciclo de Vida Controlado**

```go
// ❌ MALO: Goroutine eterno sin forma de parar
go func() {
    for {
        trabajo()
        time.Sleep(time.Second)
    }
}()

// ✓ BIEN: Respetar contexto
go func(ctx context.Context) {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            trabajo()
        }
    }
}(ctx)
```

**Antipatrón 3: No Sincronizar Acceso a Datos Compartidos**

```go
// ❌ MALO: Race condition
var contador int
go func() { contador++ }()
go func() { contador++ }()
// contador podría ser 1 o 2 (indefinido)

// ✓ BIEN: Usar sincronización
var mu sync.Mutex
var contador int
go func() {
    mu.Lock()
    defer mu.Unlock()
    contador++
}()
go func() {
    mu.Lock()
    defer mu.Unlock()
    contador++
}()
```

**Antipatrón 4: Crear Demasiadas Goroutines**

```go
// ❌ MALO: Crear una goroutine por elemento sin límite
for _, elemento := range millonesDeElementos {
    go procesarElemento(elemento)  // ⚠️ Millones de goroutines
}

// ✓ BIEN: Usar pool de workers
var wg sync.WaitGroup
trabajos := make(chan Elemento, 100)

// Crear N workers
for i := 0; i < 10; i++ {
    wg.Add(1)
    go worker(i, trabajos, &wg)
}

// Distribuir trabajo
for _, elemento := range millonesDeElementos {
    trabajos <- elemento
}
close(trabajos)
wg.Wait()
```

**Antipatrón 5: No Usar Canales para Sincronización**

```go
// ❌ MALO: Esperar arbitrariamente
go tarea()
time.Sleep(5 * time.Second)  // ¿Es suficiente?

// ✓ BIEN: Sincronización explicita
done := make(chan struct{})
go func() {
    defer close(done)
    tarea()
}()
<-done  // Esperar exactamente
```

### 21.11.3 Patrón: Worker Pool

Patrón muy común y eficiente:

```go
package main

import (
    "fmt"
    "sync"
)

type Trabajo struct {
    ID   int
    Dato string
}

type Resultado struct {
    ID     int
    Resultado string
    Error  error
}

func worker(id int, trabajos <-chan Trabajo, resultados chan<- Resultado, wg *sync.WaitGroup) {
    defer wg.Done()

    for trabajo := range trabajos {
        // Procesar trabajo
        resultado := Resultado{
            ID:        trabajo.ID,
            Resultado: fmt.Sprintf("Procesado por worker %d: %s", id, trabajo.Dato),
        }
        resultados <- resultado
    }
}

func main() {
    // Configuración
    numWorkers := 3
    numTrabajos := 10

    // Canales
    trabajos := make(chan Trabajo, numTrabajos)
    resultados := make(chan Resultado, numTrabajos)

    var wg sync.WaitGroup

    // Iniciar workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(i, trabajos, resultados, &wg)
    }

    // Enviar trabajos
    go func() {
        for i := 0; i < numTrabajos; i++ {
            trabajos <- Trabajo{ID: i, Dato: fmt.Sprintf("dato_%d", i)}
        }
        close(trabajos)
    }()

    // Recolectar resultados
    go func() {
        wg.Wait()
        close(resultados)
    }()

    // Procesar resultados
    for r := range resultados {
        fmt.Printf("[%d] %s\n", r.ID, r.Resultado)
    }
}
```

### 21.11.4 Patrón: Pipeline

Goroutines conectadas en cadena:

```go
package main

import (
    "fmt"
    "sync"
)

// Etapa 1: Generar números
func generador(max int) <-chan int {
    out := make(chan int)
    go func() {
        for i := 0; i < max; i++ {
            out <- i
        }
        close(out)
    }()
    return out
}

// Etapa 2: Duplicar
func duplicar(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * 2
        }
        close(out)
    }()
    return out
}

// Etapa 3: Elevar al cuadrado
func alCuadrado(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    // Conectar pipeline
    numeros := generador(5)
    duplicados := duplicar(numeros)
    cuadrados := alCuadrado(duplicados)

    // Consumir resultados
    for resultado := range cuadrados {
        fmt.Println(resultado)
    }

    // Output:
    // 0      (0*2)^2
    // 4      (1*2)^2
    // 16     (2*2)^2
    // 36     (3*2)^2
    // 64     (4*2)^2
}
```

### 21.11.5 Patrón: Fan-Out/Fan-In

Múltiples goroutines produciendo, una consumiendo:

```go
package main

import (
    "fmt"
    "sync"
)

// Fan-Out: Un canal se distribuye a múltiples goroutines
func distribuir(tareas <-chan int, numWorkers int) []<-chan int {
    canalesOut := make([]<-chan int, numWorkers)

    for i := 0; i < numWorkers; i++ {
        ch := make(chan int)
        canalesOut[i] = ch

        go func(out chan<- int) {
            defer close(out)
            for tarea := range tareas {
                out <- tarea * 2  // Procesar
            }
        }(ch)
    }

    return canalesOut
}

// Fan-In: Múltiples canales se combinan en uno
func combinar(canales ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    salida := make(chan int)

    for _, ch := range canales {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for valor := range c {
                salida <- valor
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(salida)
    }()

    return salida
}

func main() {
    // Crear tareas
    tareas := make(chan int, 10)
    go func() {
        for i := 0; i < 10; i++ {
            tareas <- i
        }
        close(tareas)
    }()

    // Distribuir a 3 workers
    canalesOut := distribuir(tareas, 3)

    // Combinar resultados
    resultados := combinar(canalesOut...)

    // Consumir
    for r := range resultados {
        fmt.Println(r)
    }
}
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: Hello desde Goroutine ⭐

Crear un programa que ejecute una tarea en una goroutine usando `WaitGroup`.

**Requisitos:**

- Crear una goroutine que imprima "Hola desde goroutine"
- Usar `sync.WaitGroup` para sincronización
- La goroutine debe completarse antes de que main termine

**Solución:**

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    wg.Add(1)

    go func() {
        defer wg.Done()
        fmt.Println("Hola desde goroutine")
    }()

    wg.Wait()
    fmt.Println("Main completado")
}

// Expected Output:
// Hola desde goroutine
// Main completado
```

---

### Ejercicio 2: Múltiples Tareas Concurrentes ⭐⭐

Crear 10 goroutines que ejecuten tareas concurrentemente con `WaitGroup`.

**Requisitos:**

- Crear 10 goroutines
- Cada goroutine simula trabajo con `time.Sleep`
- Usar `WaitGroup` para esperar todas
- Mostrar cuando cada tarea comienza y termina

**Solución:**

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func procesarTarea(id int, wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Printf("Tarea %d iniciada\n", id)
    time.Sleep(time.Duration(id*100) * time.Millisecond)
    fmt.Printf("Tarea %d completada\n", id)
}

func main() {
    var wg sync.WaitGroup

    start := time.Now()

    for i := 1; i <= 10; i++ {
        wg.Add(1)
        id := i
        go procesarTarea(id, &wg)
    }

    wg.Wait()

    elapsed := time.Since(start)
    fmt.Printf("Tiempo total: %.2fs\n", elapsed.Seconds())
}

// Expected: Todas las tareas ejecutándose aproximadamente en paralelo
// Tiempo total: ~1 segundo (no 5.5 segundos)
```

---

### Ejercicio 3: Timeouts con Context ⭐⭐

Crear una goroutine que respete timeout con `context.Context`.

**Requisitos:**

- Goroutine que realiza "trabajo" en un loop
- Usar `context.WithTimeout`
- Respetar la cancelación del context
- Mostrar cuántos pasos completó antes del timeout

**Solución:**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func trabajoLargo(ctx context.Context) {
    pasos := 0

    for i := 0; i < 100; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("Cancelado después de %d pasos: %v\n", pasos, ctx.Err())
            return
        default:
            pasos++
            fmt.Printf("Paso %d\n", pasos)
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Printf("Completado: %d pasos\n", pasos)
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    trabajoLargo(ctx)
}

// Expected: Aproximadamente 5 pasos antes del timeout
```

---

### Ejercicio 4: Cancelación Manual de Goroutines ⭐⭐⭐

Crear múltiples servidores que se pueden cancelar con `context.WithCancel`.

**Requisitos:**

- Crear 3 "servidores" concurrentes
- Usar `context.WithCancel` para controlar la ejecución
- Cancelar todos después de 2 segundos
- Mostrar que se cierran ordenadamente

**Solución:**

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

func servidor(ctx context.Context, nombre string, wg *sync.WaitGroup) {
    defer wg.Done()

    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            fmt.Printf("✗ %s apagado\n", nombre)
            return
        case <-ticker.C:
            fmt.Printf("✓ %s ejecutando\n", nombre)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    var wg sync.WaitGroup

    // Iniciar 3 servidores
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        nombre := fmt.Sprintf("Servidor-%d", i)
        go servidor(ctx, nombre, &wg)
    }

    // Esperar 2 segundos
    time.Sleep(2 * time.Second)

    // Cancelar todos
    fmt.Println("Cancelando servidores...")
    cancel()

    wg.Wait()
    fmt.Println("Todos los servidores cerrados")
}

// Expected: Servidores ejecutando, luego se cierran elegantemente
```

---

### Ejercicio 5: Generador - Patrón Productor ⭐⭐⭐⭐

Implementar patrón generador que produce valores en una goroutine.

**Requisitos:**

- Goroutine que genera números en un rango
- Usar channel para enviar valores
- Consumidor que procesa valores
- Implementar graceful shutdown
- Usar `sync.WaitGroup` para sincronización

**Solución:**

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Generador produce números del 1 al max
func generador(ctx context.Context, max int) <-chan int {
    out := make(chan int)

    go func() {
        defer close(out)

        for i := 1; i <= max; i++ {
            select {
            case <-ctx.Done():
                fmt.Println("Generador cancelado")
                return
            case out <- i:
                // Enviar número
            }
        }
    }()

    return out
}

// Procesador consume números del generador
func procesador(ctx context.Context, numeros <-chan int, id int, wg *sync.WaitGroup) {
    defer wg.Done()

    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Procesador %d terminado\n", id)
            return
        case num, ok := <-numeros:
            if !ok {
                fmt.Printf("Procesador %d: canal cerrado\n", id)
                return
            }
            fmt.Printf("Procesador %d: %d * 2 = %d\n", id, num, num*2)
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var wg sync.WaitGroup

    // Generador produce números del 1 al 20
    numeros := generador(ctx, 20)

    // Crear 2 procesadores
    for i := 1; i <= 2; i++ {
        wg.Add(1)
        id := i
        go procesador(ctx, numeros, id, &wg)
    }

    // Esperar 5 segundos, luego cancelar
    time.Sleep(5 * time.Second)
    fmt.Println("\nCancelando...")
    cancel()

    wg.Wait()
    fmt.Println("✓ Programa completado")
}

// Expected: Generador produce números, 2 procesadores los consumen en paralelo
```

---

## RESUMEN

Las **goroutines** son el corazón de la concurrencia en Go:

1. **Livianas**: ~2KB vs ~2MB de threads del SO
2. **Gestionadas**: El runtime de Go las gestiona eficientemente
3. **Simples**: Sintaxis elegante con `go` keyword
4. **Poderosas**: M:N scheduling para máxima eficiencia

**Conceptos clave:**

- **WaitGroup**: Sincronización de múltiples goroutines
- **Context**: Control de ciclo de vida, cancelación, timeouts
- **Leaks**: Evitar goroutines que nunca terminan
- **Patrones**: Worker pool, pipelines, fan-out/fan-in

**Best practices:**
✓ Sincronización explícita con WaitGroup o channels
✓ Respetar context para cancelación y timeouts
✓ Documentar comportamiento concurrente
✓ Usar pools de workers para limitar concurrencia
✓ Detectar goroutine leaks durante desarrollo

Go hace que la programación concurrente sea elegante, segura y eficiente.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/21-goroutines/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/21-goroutines):

```bash
cd examples/21-goroutines
go run .
```
