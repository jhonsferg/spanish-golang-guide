# Capítulo 23: Select - Multiplexing de channels

## Índice
1. [¿Qué es Select?](#231-qué-es-select)
2. [Sintaxis Básica](#232-sintaxis-básica)
3. [Select Bloqueante](#233-select-bloqueante)
4. [Default Case](#234-default-case)
5. [Timeouts](#235-timeouts)
6. [Cancelación](#236-cancelación)
7. [Select en Loops](#237-select-en-loops)
8. [Ordering de Select Cases](#238-ordering-de-select-cases)
9. [Select vs If](#239-select-vs-if)
10. [Patrones Avanzados](#2310-patrones-avanzados)
11. [Buenas Prácticas y Antipatrones](#2311-buenas-prácticas-y-antipatrones)

---

## 23.1 ¿Qué es Select?

### 23.1.1 Concepto Fundamental

El **`select` statement** es un mecanismo de control de flujo que permite a una goroutine esperar simultáneamente múltiples operaciones de comunicación (envío y recepción en channels). Es el equivalente de Go a las operaciones de multiplexing I/O en sistemas operativos.

El select resuelve un problema crítico en programación concurrente: **¿Cómo esperar múltiples eventos potenciales sin bloquear indefinidamente?**

```
┌─────────────────────────────────────────────┐
│         SELECT MULTIPLEXING MODEL           │
├─────────────────────────────────────────────┤
│                                             │
│  Channel 1 ─┐                              │
│             ├─ SELECT ─ Espera a que uno  │
│  Channel 2 ─┤         esté listo           │
│             │                             │
│  Channel 3 ─┘                              │
│                                             │
│  Timeout ──────┤ (mediante time.After)     │
│                                             │
└─────────────────────────────────────────────┘
```

### 23.1.2 Características Clave

**1. Multiplexing de Channels**
- Escucha múltiples channels simultáneamente
- Procesa el primero que tenga datos disponibles
- Sin polling activo (blocking eficiente)

**2. No Determinista**
- Si múltiples channels están listos, elige uno al azar
- Evita sesgo hacia canales específicos
- Garantiza fairness en sistemas de carga balanceada

**3. Bloqueante por Defecto**
- Bloquea la goroutine si ningún case está listo
- Se desbloquea cuando cualquier case pueda proceder
- Usa CPU zero mientras espera

**4. Timeout Native**
- Integración nativa con `time.After` para timeouts
- Cancelación con `context.Done()`
- Patrón elegante para operaciones temporizadas

### 23.1.3 Caso de Uso Real: Servidor Web Concurrente

```go
package main

import (
    "fmt"
    "time"
)

// Simulación de servidor que escucha múltiples fuentes
func serverMultiplexing() {
    requestsChan := make(chan string, 10)
    shutdownChan := make(chan bool)
    
    // Enviar solicitudes
    go func() {
        for i := 1; i <= 5; i++ {
            time.Sleep(500 * time.Millisecond)
            requestsChan <- fmt.Sprintf("Solicitud %d", i)
        }
    }()
    
    // Enviar señal de shutdown después de 3 segundos
    go func() {
        time.Sleep(3 * time.Second)
        shutdownChan <- true
    }()
    
    // Servidor multiplexando dos fuentes
    for {
        select {
        case req := <-requestsChan:
            fmt.Printf("Procesando: %s\n", req)
        case <-shutdownChan:
            fmt.Println("Servidor apagando...")
            return
        }
    }
}
```

### 23.1.4 Comparación con Otros Lenguajes

| Característica | Go Select | Rust Match | Java Select/Poll |
|---|---|---|---|
| Sintaxis | `select {...}` | `match` + `tokio::select!` | `Selector.select()` |
| Channels | Tipo nativo | Async/await | Selectable |
| Blocking | Implícito | Explícito | Explícito |
| Default | Con `default:` | Con `_` | N/A |
| Timeout | `time.After` | `tokio::time::timeout` | `SelectionKey.interestOps()` |
| Non-determinism | Nativo | Nativo | Nativo |
| Curva aprendizaje | Media | Alta | Alta |

---

## 23.2 Sintaxis Básica

### 23.2.1 Estructura General

```go
select {
case <-channel1:
    // Ejecuta si hay datos en channel1
    
case value := <-channel2:
    // Ejecuta si hay datos en channel2
    // value contiene el dato
    
case channel3 <- data:
    // Ejecuta si se puede enviar a channel3
    
case <-time.After(1 * time.Second):
    // Ejecuta después de 1 segundo si no hay otros casos listos
    
default:
    // Ejecuta si ningún case está listo (no bloqueante)
}
```

### 23.2.2 Casos de Uso Común

**Recibir de un Channel**
```go
select {
case msg := <-msgChan:
    fmt.Println("Mensaje:", msg)
}
```

**Enviar a un Channel**
```go
select {
case dataChan <- value:
    fmt.Println("Dato enviado")
}
```

**Múltiples Recepciones**
```go
select {
case req := <-requests:
    handleRequest(req)
    
case msg := <-messages:
    processMessage(msg)
    
case alert := <-alerts:
    handleAlert(alert)
}
```

**Mezcla de Envíos y Recepciones**
```go
select {
case value := <-inputChan:
    processValue(value)
    
case resultChan <- computeResult():
    fmt.Println("Resultado enviado")
}
```

### 23.2.3 Orden de Ejecución

El select tiene un orden de evaluación específico pero no determinista:

```go
// Compilación: orden izquierda a derecha
// Evaluación: ALEATORIA si múltiples están listos
select {
case <-chan1:  // Case 1 - Posición 1
    fmt.Println("Chan1")
    
case <-chan2:  // Case 2 - Posición 2
    fmt.Println("Chan2")
    
case <-chan3:  // Case 3 - Posición 3
    fmt.Println("Chan3")
}
```

Si chan1, chan2 y chan3 tienen datos, Go elige aleatoriamente cuál procesar. Esto es **intencional** para garantizar fairness.

### 23.2.4 Patrón General de Select Loop

```go
for {
    select {
    case job := <-jobQueue:
        fmt.Printf("Procesando: %v\n", job)
        
    case result := <-resultChannel:
        fmt.Printf("Resultado: %v\n", result)
        
    case <-stopChan:
        fmt.Println("Parando...")
        return
    }
}
```

---

## 23.3 Select Bloqueante

### 23.3.1 Comportamiento Bloqueante

El select **bloquea** la goroutine hasta que al menos uno de sus cases pueda proceder. No consume CPU mientras espera.

```go
package main

import (
    "fmt"
    "time"
)

func blockingSelectDemo() {
    chan1 := make(chan string)
    chan2 := make(chan string)
    
    // Goroutine que envía después de 2 segundos
    go func() {
        time.Sleep(2 * time.Second)
        chan1 <- "¡Datos en chan1!"
    }()
    
    // Select bloqueante - espera 2 segundos sin hacer nada
    fmt.Println("Esperando...")
    start := time.Now()
    
    select {
    case msg := <-chan1:
        fmt.Printf("Recibido: %s después de %v\n", 
            msg, time.Since(start))
        
    case msg := <-chan2:
        fmt.Printf("Recibido: %s\n", msg)
    }
}
```

Output:
```
Esperando...
Recibido: ¡Datos en chan1! después de 2.001234567s
```

### 23.3.2 Deadlock por Select Bloqueante

```go
// ❌ DEADLOCK: Select bloqueante en main sin goroutines
func deadlockExample() {
    chan1 := make(chan string)
    
    // Esto bloquea PARA SIEMPRE porque:
    // 1. chan1 no tiene datos
    // 2. No hay goroutines que envíen datos
    // 3. Select está esperando indefinidamente
    select {
    case msg := <-chan1:
        fmt.Println(msg)
    }
    
    // Nunca llega aquí
    fmt.Println("Fin")
}
```

### 23.3.3 Detectar Bloques Largos

```go
// ✓ CORRECTO: Usar timeout para detectar bloqueos
func selectWithTimeout() {
    chan1 := make(chan string)
    
    select {
    case msg := <-chan1:
        fmt.Println(msg)
        
    case <-time.After(3 * time.Second):
        fmt.Println("Timeout: nada llegó en 3 segundos")
    }
}
```

### 23.3.4 Select Vacío

```go
// ❌ ANTIPATRÓN: Select vacío bloquea para siempre
func emptySelectForever() {
    select {}  // ¡Bloquea indefinidamente!
}

// ✓ USO VÁLIDO: Mantener programa activo
func serverKeepAlive() {
    go func() {
        // Hacer cosas en background
        for {
            fmt.Println("Trabajando...")
            time.Sleep(1 * time.Second)
        }
    }()
    
    // Mantener main viva
    select {}
}
```

---

## 23.4 Default Case

### 23.4.1 Comportamiento No Bloqueante

El `default` case permite que select **no bloquee** si ningún otro case está listo:

```go
package main

import (
    "fmt"
    "time"
)

func defaultCaseDemo() {
    dataChan := make(chan int)
    
    // Enviar datos después de 2 segundos
    go func() {
        time.Sleep(2 * time.Second)
        dataChan <- 42
    }()
    
    // Verificar múltiples veces sin bloquear
    for i := 0; i < 5; i++ {
        select {
        case data := <-dataChan:
            fmt.Printf("Datos recibidos: %d\n", data)
            return
            
        default:
            fmt.Printf("No hay datos (intento %d)\n", i+1)
            time.Sleep(1 * time.Second)
        }
    }
}
```

Output:
```
No hay datos (intento 1)
No hay datos (intento 2)
Datos recibidos: 42
```

### 23.4.2 Non-blocking Send

```go
func nonBlockingSend() {
    resultChan := make(chan string)
    
    // Envío que no bloquea
    select {
    case resultChan <- "Resultado":
        fmt.Println("Resultado enviado")
        
    default:
        fmt.Println("No se pudo enviar (buffer lleno o sin receptores)")
    }
}
```

### 23.4.3 Patrón: Polling No Bloqueante

```go
func pollingPattern() {
    jobs := make(chan string, 10)
    results := make(chan string)
    
    // Llenar canal
    jobs <- "job1"
    jobs <- "job2"
    
    // Procesar con polling no bloqueante
    processed := 0
    for processed < 2 {
        select {
        case job := <-jobs:
            fmt.Printf("Procesando: %s\n", job)
            go func(j string) {
                time.Sleep(100 * time.Millisecond)
                results <- j + " completado"
            }(job)
            
        case result := <-results:
            fmt.Printf("Resultado: %s\n", result)
            processed++
            
        default:
            // Pequeña pausa si no hay nada
            time.Sleep(10 * time.Millisecond)
        }
    }
}
```

### 23.4.4 Patrones de Fallback

```go
// Patrón: Try-best-effort
func trySendWithFallback() {
    alertChan := make(chan string)
    logChan := make(chan string)
    
    alert := "CRÍTICO: Sistema degradado"
    
    select {
    case alertChan <- alert:
        fmt.Println("Alerta enviada por canal prioritario")
        
    default:
        // Fallback: usar log
        select {
        case logChan <- alert:
            fmt.Println("Alerta registrada en log")
        default:
            fmt.Println("Sistema saturado, descartando alerta")
        }
    }
}
```

---

## 23.5 Timeouts

### 23.5.1 Pattern: time.After

El patrón más común para implementar timeouts:

```go
package main

import (
    "fmt"
    "time"
)

func simpleTimeout() {
    resultChan := make(chan string)
    
    go func() {
        // Simular operación lenta
        time.Sleep(3 * time.Second)
        resultChan <- "Resultado"
    }()
    
    select {
    case result := <-resultChan:
        fmt.Println("Éxito:", result)
        
    case <-time.After(1 * time.Second):
        fmt.Println("❌ Timeout: operación tardó demasiado")
    }
}
```

### 23.5.2 Timeout con Retry

```go
func timeoutWithRetry() {
    for attempt := 1; attempt <= 3; attempt++ {
        fmt.Printf("Intento %d...\n", attempt)
        
        result := attemptWithTimeout(500 * time.Millisecond)
        if result != "" {
            fmt.Println("Éxito:", result)
            return
        }
        fmt.Println("Reintentando...")
    }
    fmt.Println("Fallido después de 3 intentos")
}

func attemptWithTimeout(timeout time.Duration) string {
    resultChan := make(chan string)
    
    go func() {
        // Operación que a veces falla
        if time.Now().UnixNano()%2 == 0 {
            resultChan <- "Operación exitosa"
        }
        // Si no envía nada, se produce timeout
    }()
    
    select {
    case result := <-resultChan:
        return result
        
    case <-time.After(timeout):
        return ""
    }
}
```

### 23.5.3 Timeout en Operaciones I/O

```go
func fetchDataWithTimeout() {
    dataChan := make(chan string)
    errorChan := make(chan error)
    
    go func() {
        // Simular fetch de API
        time.Sleep(2 * time.Second)
        dataChan <- "Datos del servidor"
    }()
    
    select {
    case data := <-dataChan:
        fmt.Println("Datos:", data)
        
    case err := <-errorChan:
        fmt.Println("Error:", err)
        
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout: Servidor no respondió")
    }
}
```

### 23.5.4 Pattern: time.Tick para Operaciones Periódicas

```go
func tickerPattern() {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    jobsChan := make(chan string)
    stopChan := make(chan bool)
    
    // Goroutine que envía jobs
    go func() {
        for i := 1; i <= 5; i++ {
            time.Sleep(200 * time.Millisecond)
            jobsChan <- fmt.Sprintf("Job %d", i)
        }
        stopChan <- true
    }()
    
    for {
        select {
        case <-ticker.C:
            fmt.Println("[Heartbeat] Sistema activo")
            
        case job := <-jobsChan:
            fmt.Printf("[Job] %s\n", job)
            
        case <-stopChan:
            fmt.Println("[Stop] Terminando")
            return
        }
    }
}
```

### 23.5.5 Timeout Escalonado

```go
func escalatingTimeout() {
    timeouts := []time.Duration{
        100 * time.Millisecond,
        500 * time.Millisecond,
        2 * time.Second,
    }
    
    for i, timeout := range timeouts {
        fmt.Printf("Intento %d con timeout %v\n", i+1, timeout)
        
        resultChan := make(chan string)
        
        go func() {
            time.Sleep(300 * time.Millisecond)
            resultChan <- "Éxito"
        }()
        
        select {
        case result := <-resultChan:
            fmt.Println("✓", result)
            return
            
        case <-time.After(timeout):
            if i < len(timeouts)-1 {
                fmt.Println("✗ Timeout, reintentando...")
            } else {
                fmt.Println("✗ Falló después de todos los intentos")
            }
        }
    }
}
```

---

## 23.6 Cancelación

### 23.6.1 Patrón: Done Channel

```go
package main

import (
    "fmt"
    "time"
)

func doneChannelPattern() {
    done := make(chan bool)
    resultChan := make(chan string)
    
    // Goroutine que puede ser cancelada
    go func() {
        for i := 0; i < 10; i++ {
            select {
            case <-done:
                fmt.Println("Goroutine cancelada")
                return
                
            default:
                fmt.Printf("Trabajando... %d\n", i)
                time.Sleep(500 * time.Millisecond)
                
                if i == 4 {
                    resultChan <- "Resultado a mitad"
                }
            }
        }
    }()
    
    select {
    case result := <-resultChan:
        fmt.Println("Recibido:", result)
        done <- true  // Cancelar goroutine
        
    case <-time.After(3 * time.Second):
        fmt.Println("Timeout, cancelando...")
        done <- true
    }
    
    time.Sleep(500 * time.Millisecond)
}
```

### 23.6.2 Context-Based Cancellation

```go
import "context"

func contextCancellation() {
    // Crear contexto con timeout
    ctx, cancel := context.WithTimeout(
        context.Background(), 
        2 * time.Second,
    )
    defer cancel()
    
    resultChan := make(chan string)
    
    go func() {
        for i := 0; i < 10; i++ {
            select {
            case <-ctx.Done():
                fmt.Println("Contexto cancelado:", ctx.Err())
                return
                
            default:
                time.Sleep(500 * time.Millisecond)
                fmt.Printf("Paso %d\n", i)
            }
        }
    }()
    
    select {
    case result := <-resultChan:
        fmt.Println("Resultado:", result)
        
    case <-ctx.Done():
        fmt.Println("Operación cancelada por timeout")
    }
}
```

### 23.6.3 Shutdown Graceful

```go
type Server struct {
    stopChan chan bool
    jobsChan chan string
}

func (s *Server) Start() {
    go func() {
        for {
            select {
            case <-s.stopChan:
                fmt.Println("Servidor detenido")
                return
                
            case job := <-s.jobsChan:
                s.processJob(job)
            }
        }
    }()
}

func (s *Server) Stop() {
    fmt.Println("Iniciando shutdown graceful...")
    
    // Dar tiempo para procesar jobs pendientes
    time.Sleep(1 * time.Second)
    
    s.stopChan <- true
    fmt.Println("Servidor parado")
}

func (s *Server) processJob(job string) {
    fmt.Printf("Procesando: %s\n", job)
    time.Sleep(200 * time.Millisecond)
}

func gracefulShutdownDemo() {
    server := &Server{
        stopChan: make(chan bool),
        jobsChan: make(chan string, 100),
    }
    
    server.Start()
    
    // Enviar jobs
    for i := 1; i <= 5; i++ {
        server.jobsChan <- fmt.Sprintf("Job %d", i)
    }
    
    time.Sleep(2 * time.Second)
    server.Stop()
}
```

### 23.6.4 Cancel al Primero en Completar

```go
func cancelOnFirstCompletion() {
    results := make(chan string)
    done := make(chan bool)
    
    // Múltiples goroutines competidoras
    for i := 1; i <= 3; i++ {
        go func(id int) {
            delay := time.Duration(id) * time.Second
            time.Sleep(delay)
            
            select {
            case <-done:
                fmt.Printf("Worker %d: Cancelado\n", id)
                
            default:
                results <- fmt.Sprintf("Worker %d completó", id)
            }
        }(i)
    }
    
    // Esperar al primero que termine
    select {
    case result := <-results:
        fmt.Println("✓", result)
        close(done)  // Cancelar otros
    }
    
    time.Sleep(2 * time.Second)
}
```

---

## 23.7 Select en Loops

### 23.7.1 For-Select Pattern

```go
package main

import (
    "fmt"
    "time"
)

func forSelectPattern() {
    inputChan := make(chan string, 10)
    outputChan := make(chan string)
    stopChan := make(chan bool)
    
    // Producer
    go func() {
        for i := 1; i <= 5; i++ {
            inputChan <- fmt.Sprintf("Item %d", i)
            time.Sleep(200 * time.Millisecond)
        }
        close(inputChan)
    }()
    
    // Consumer loop
    for {
        select {
        case item, ok := <-inputChan:
            if !ok {
                fmt.Println("Input cerrado, esperando outputs...")
                continue
            }
            fmt.Printf("Recibido: %s\n", item)
            go func(i string) {
                time.Sleep(100 * time.Millisecond)
                outputChan <- i + " procesado"
            }(item)
            
        case result := <-outputChan:
            fmt.Printf("Resultado: %s\n", result)
            
        case <-stopChan:
            fmt.Println("Deteniendo...")
            return
        }
    }
}
```

### 23.7.2 Worker Pool con Select

```go
type Task struct {
    ID   int
    Data string
}

type WorkerPool struct {
    workers int
    tasks   chan Task
    results chan string
    done    chan bool
}

func (wp *WorkerPool) Start() {
    for w := 0; w < wp.workers; w++ {
        go wp.worker(w)
    }
}

func (wp *WorkerPool) worker(id int) {
    for {
        select {
        case task, ok := <-wp.tasks:
            if !ok {
                fmt.Printf("Worker %d: Canal cerrado\n", id)
                return
            }
            result := fmt.Sprintf("Worker %d procesó Task %d: %s", 
                id, task.ID, task.Data)
            wp.results <- result
            
        case <-wp.done:
            fmt.Printf("Worker %d: Cancelado\n", id)
            return
        }
    }
}

func workerPoolDemo() {
    wp := &WorkerPool{
        workers: 3,
        tasks:   make(chan Task, 10),
        results: make(chan string),
        done:    make(chan bool),
    }
    
    wp.Start()
    
    // Enviar tareas
    go func() {
        for i := 1; i <= 10; i++ {
            wp.tasks <- Task{ID: i, Data: fmt.Sprintf("data_%d", i)}
        }
        close(wp.tasks)
    }()
    
    // Procesar resultados
    processed := 0
    for processed < 10 {
        select {
        case result := <-wp.results:
            fmt.Println(result)
            processed++
        }
    }
    
    // Señalizar fin
    for i := 0; i < 3; i++ {
        wp.done <- true
    }
}
```

### 23.7.3 Fan-in Pattern

```go
// Combinar múltiples canales en uno
func fanInPattern() {
    producer := func(name string) <-chan string {
        out := make(chan string)
        go func() {
            for i := 1; i <= 3; i++ {
                time.Sleep(time.Duration(i*100) * time.Millisecond)
                out <- fmt.Sprintf("%s: mensaje %d", name, i)
            }
            close(out)
        }()
        return out
    }
    
    // Fan-in: combinar múltiples canales
    merge := func(channels ...<-chan string) <-chan string {
        merged := make(chan string)
        go func() {
            count := len(channels)
            for count > 0 {
                for _, ch := range channels {
                    select {
                    case msg, ok := <-ch:
                        if ok {
                            merged <- msg
                        } else {
                            count--
                            continue
                        }
                    default:
                    }
                }
            }
            close(merged)
        }()
        return merged
    }
    
    // Crear productores
    ch1 := producer("A")
    ch2 := producer("B")
    ch3 := producer("C")
    
    // Combinar
    merged := merge(ch1, ch2, ch3)
    
    for msg := range merged {
        fmt.Println(msg)
    }
}
```

### 23.7.4 Fan-out Pattern

```go
// Distribuir trabajo a múltiples workers
func fanOutPattern() {
    tasks := []int{1, 2, 3, 4, 5}
    
    // Crear canales para cada worker
    results := make([]<-chan string, len(tasks))
    
    for i, task := range tasks {
        ch := make(chan string)
        go func(t int, c chan string) {
            time.Sleep(time.Duration(t*100) * time.Millisecond)
            c <- fmt.Sprintf("Resultado de tarea %d", t)
            close(c)
        }(task, ch)
        results[i] = ch
    }
    
    // Recolectar todos los resultados
    for i, resultCh := range results {
        select {
        case result := <-resultCh:
            fmt.Printf("Task %d: %s\n", i+1, result)
        }
    }
}
```

---

## 23.8 Ordering de Select Cases

### 23.8.1 No Determinismo por Diseño

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func nonDeterminismDemo() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    // Ambos listos simultáneamente
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "Mensaje de ch1"
        ch2 <- "Mensaje de ch2"
    }()
    
    // Ejecutar 10 veces - resultado diferente cada vez
    results := make(map[string]int)
    
    for trial := 0; trial < 10; trial++ {
        select {
        case msg := <-ch1:
            results[msg]++
            
        case msg := <-ch2:
            results[msg]++
        }
    }
    
    fmt.Println("Distribución de 10 ejecuciones:")
    for msg, count := range results {
        fmt.Printf("  %s: %d veces\n", msg, count)
    }
    // Output muestra ~50/50 split entre ch1 y ch2
}
```

### 23.8.2 Implicaciones para Testing

```go
// ❌ TEST FRÁGIL: Asume orden específico
func fragileTest(t *testing.T) {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    go func() {
        ch1 <- 1
        ch2 <- 2
    }()
    
    select {
    case v := <-ch1:
        if v != 1 {
            t.Fail()
        }
    case v := <-ch2:
        if v != 2 {
            t.Fail()  // ¡Puede fallar aleatoriamente!
        }
    }
}

// ✓ TEST ROBUSTO: Recopila ambos resultados
func robustTest(t *testing.T) {
    ch1 := make(chan int)
    ch2 := make(chan int)
    results := []int{}
    
    go func() {
        ch1 <- 1
        ch2 <- 2
    }()
    
    // Recolectar ambos valores
    for i := 0; i < 2; i++ {
        select {
        case v := <-ch1:
            results = append(results, v)
        case v := <-ch2:
            results = append(results, v)
        }
    }
    
    // Verificar contenido sin depender del orden
    if len(results) != 2 || (results[0] != 1 && results[0] != 2) {
        t.Fail()
    }
}
```

### 23.8.3 Fairness en Carga Equilibrada

```go
// Demostrar que el no-determinismo asegura fairness
func fairnessDemo() {
    heavy := make(chan int, 100)    // Muchos datos
    light := make(chan int, 10)     // Pocos datos
    
    // Llenar los canales
    for i := 0; i < 100; i++ {
        heavy <- i
    }
    for i := 0; i < 10; i++ {
        light <- i
    }
    
    heavyCount := 0
    lightCount := 0
    
    // Sin no-determinismo, se procesaría heavy primero
    // Con no-determinismo, se distribuye equitativamente
    for i := 0; i < 110; i++ {
        select {
        case <-heavy:
            heavyCount++
        case <-light:
            lightCount++
        }
    }
    
    fmt.Printf("Heavy: %d, Light: %d\n", heavyCount, lightCount)
    // Típicamente ~55 y ~55, no 100 y 10
}
```

### 23.8.4 Mitigación: Priority Patterns

Si necesitas prioridad, usa select anidado:

```go
// Patrón: Prioridad con select anidado
func priorityPattern() {
    critical := make(chan string)
    normal := make(chan string)
    background := make(chan string)
    
    // Procesar con prioridad explícita
    for {
        select {
        case task := <-critical:
            fmt.Printf("CRÍTICO: %s\n", task)
            
        default:
            // Luego evaluar normal
            select {
            case task := <-normal:
                fmt.Printf("Normal: %s\n", task)
                
            default:
                // Finalmente background
                select {
                case task := <-background:
                    fmt.Printf("Background: %s\n", task)
                    
                default:
                    return
                }
            }
        }
    }
}
```

---

## 23.9 Select vs If

### 23.9.1 Cuándo Usar Cada Uno

```
┌─────────────────────────────────────────────────────────┐
│           SELECT vs IF-ELSE para CHANNELS               │
├──────────────────────────────────┬──────────────────────┤
│ USE SELECT CUANDO:               │ USE IF-ELSE CUANDO:  │
├──────────────────────────────────┼──────────────────────┤
│ • Múltiples channels             │ • Un único channel   │
│ • Timeout necesario              │ • Sin timeout        │
│ • Operaciones concurrentes       │ • Lógica secuencial  │
│ • No determinismo aceptable      │ • Orden fijo         │
│ • Multiplexing real              │ • Validación simple  │
└──────────────────────────────────┴──────────────────────┘
```

### 23.9.2 Comparación Práctica

**Scenario: Recibir de un canal**

```go
// ✓ CON IF - Válido para un canal
func receiveWithIf() {
    ch := make(chan string)
    go func() {
        ch <- "Hola"
    }()
    
    // Bloquea hasta que hay datos
    msg := <-ch
    fmt.Println(msg)
}

// ✓ CON SELECT - Mejor para control avanzado
func receiveWithSelect() {
    ch := make(chan string)
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch <- "Hola"
    }()
    
    select {
    case msg := <-ch:
        fmt.Println(msg)
        
    case <-time.After(50 * time.Millisecond):
        fmt.Println("Timeout")
    }
}
```

**Scenario: Múltiples canales**

```go
// ❌ CON IF - Innecesariamente complejo
func multipleChannelsWithIf() {
    ch1 := make(chan string)
    ch2 := make(chan int)
    
    // No hay forma elegante de esperar ambos
    // Requeriría goroutines adicionales
}

// ✓ CON SELECT - Natural y limpio
func multipleChannelsWithSelect() {
    ch1 := make(chan string)
    ch2 := make(chan int)
    
    select {
    case msg := <-ch1:
        fmt.Println("String:", msg)
        
    case num := <-ch2:
        fmt.Println("Int:", num)
    }
}
```

**Scenario: Operaciones no bloqueantes**

```go
// ❌ CON IF - No se puede hacer
func nonBlockingWithIf() {
    ch := make(chan string)
    
    // Esto bloquea o panics
    // msg := <-ch  // bloquea
}

// ✓ CON SELECT - Default lo permite
func nonBlockingWithSelect() {
    ch := make(chan string)
    
    select {
    case msg := <-ch:
        fmt.Println("Mensaje:", msg)
        
    default:
        fmt.Println("No hay datos disponibles")
    }
}
```

### 23.9.3 Lógica Condicional

```go
// ✓ PATRÓN: Combinar select con if para lógica compleja
func advancedLogic() {
    dataChan := make(chan string)
    errorChan := make(chan error)
    
    go func() {
        // Simulación: 50% éxito, 50% error
        if time.Now().UnixNano()%2 == 0 {
            dataChan <- "Datos"
        } else {
            errorChan <- fmt.Errorf("fallo")
        }
    }()
    
    select {
    case data := <-dataChan:
        // Select elegir cuál channel tenía datos
        if len(data) > 0 {  // If para lógica dentro
            fmt.Println("Datos válidos:", data)
        }
        
    case err := <-errorChan:
        if err != nil {  // If para validación
            fmt.Println("Error:", err)
        }
        
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout")
    }
}
```

---

## 23.10 Patrones Avanzados

### 23.10.1 Rate Limiter

```go
package main

import (
    "fmt"
    "time"
)

func rateLimiterPattern() {
    // Limitar a 5 solicitudes por segundo
    ticker := time.NewTicker(200 * time.Millisecond)  // 1000/5
    defer ticker.Stop()
    
    requests := make(chan string, 20)
    results := make(chan string)
    
    // Productor: genera muchas solicitudes rápido
    go func() {
        for i := 1; i <= 20; i++ {
            requests <- fmt.Sprintf("Request %d", i)
        }
        close(requests)
    }()
    
    // Procesador rate-limited
    go func() {
        for req := range requests {
            <-ticker.C  // Esperar token
            results <- fmt.Sprintf("Procesado: %s", req)
        }
        close(results)
    }()
    
    // Consumer
    for result := range results {
        fmt.Println(result, "en", time.Now().Format("15:04:05"))
    }
}
```

### 23.10.2 Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    attempts  int
    maxAttempts int
    isOpen    bool
    reset     <-chan time.Time
}

func (cb *CircuitBreaker) attempt(fn func() error) error {
    if cb.isOpen {
        select {
        case <-cb.reset:
            cb.isOpen = false
            cb.attempts = 0
        default:
            return fmt.Errorf("circuit breaker abierto")
        }
    }
    
    err := fn()
    if err != nil {
        cb.attempts++
        if cb.attempts >= cb.maxAttempts {
            cb.isOpen = true
            ticker := time.NewTicker(2 * time.Second)
            cb.reset = ticker.C
        }
    } else {
        cb.attempts = 0
    }
    
    return err
}

func circuitBreakerDemo() {
    cb := &CircuitBreaker{maxAttempts: 3}
    
    for i := 0; i < 8; i++ {
        err := cb.attempt(func() error {
            if i < 4 {
                return fmt.Errorf("fallo simulado")
            }
            return nil
        })
        
        if err != nil {
            fmt.Printf("Intento %d: %v\n", i+1, err)
        } else {
            fmt.Printf("Intento %d: Éxito\n", i+1)
        }
        
        time.Sleep(500 * time.Millisecond)
    }
}
```

### 23.10.3 Bulkhead Pattern

Aislar recursos para evitar que un fallo cascade:

```go
type Bulkhead struct {
    semaphore chan bool
    handler   func(task string) error
}

func NewBulkhead(maxConcurrent int, handler func(string) error) *Bulkhead {
    return &Bulkhead{
        semaphore: make(chan bool, maxConcurrent),
        handler:   handler,
    }
}

func (b *Bulkhead) execute(task string) {
    select {
    case b.semaphore <- true:
        go func() {
            defer func() { <-b.semaphore }()
            
            err := b.handler(task)
            if err != nil {
                fmt.Printf("Task %s falló: %v\n", task, err)
            } else {
                fmt.Printf("Task %s completada\n", task)
            }
        }()
        
    default:
        fmt.Printf("Task %s rechazada (sin capacidad)\n", task)
    }
}

func bulkheadDemo() {
    bh := NewBulkhead(3, func(task string) error {
        time.Sleep(time.Second)
        return nil
    })
    
    // Enviar 10 tareas, pero solo 3 concurrentes
    for i := 1; i <= 10; i++ {
        bh.execute(fmt.Sprintf("Task %d", i))
    }
    
    time.Sleep(5 * time.Second)
}
```

### 23.10.4 Pipeline Pattern

```go
func pipelinePattern() {
    // Stage 1: Generar números
    numbers := make(chan int)
    go func() {
        for i := 1; i <= 5; i++ {
            numbers <- i
        }
        close(numbers)
    }()
    
    // Stage 2: Cuadrar números
    squares := make(chan int)
    go func() {
        for n := range numbers {
            squares <- n * n
        }
        close(squares)
    }()
    
    // Stage 3: Consumir resultados
    for sq := range squares {
        fmt.Println("Cuadrado:", sq)
    }
}
```

### 23.10.5 Timeout con Context

```go
import "context"

func operationWithContextTimeout(ctx context.Context, data string) (string, error) {
    resultChan := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        resultChan <- "Resultado: " + data
    }()
    
    select {
    case result := <-resultChan:
        return result, nil
        
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func contextTimeoutDemo() {
    // Timeout de 1 segundo
    ctx, cancel := context.WithTimeout(
        context.Background(),
        1 * time.Second,
    )
    defer cancel()
    
    result, err := operationWithContextTimeout(ctx, "datos")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println(result)
    }
}
```

---

## 23.11 Buenas Prácticas y Antipatrones

### 23.11.1 Buenas Prácticas

#### 1. Siempre Finalizar Loops Correctamente

```go
// ✓ CORRECTO: Detectar channel cerrado
func properChannelHandling() {
    ch := make(chan string)
    
    go func() {
        ch <- "dato1"
        ch <- "dato2"
        close(ch)  // Señalizar fin
    }()
    
    for msg := range ch {  // Detecta automáticamente close()
        fmt.Println(msg)
    }
}

// ❌ INCORRECTO: Panic si se recibe de channel cerrado
func improperChannelHandling() {
    ch := make(chan string)
    close(ch)
    
    msg := <-ch  // Panic!
}
```

#### 2. Usar Context para Cancelación

```go
// ✓ CORRECTO: Context para timeout y cancelación
func properCancellation(ctx context.Context) {
    resultChan := make(chan string)
    
    go func() {
        select {
        case <-ctx.Done():
            fmt.Println("Operación cancelada")
            return
        default:
            time.Sleep(2 * time.Second)
            resultChan <- "Resultado"
        }
    }()
    
    select {
    case result := <-resultChan:
        fmt.Println(result)
    case <-ctx.Done():
        fmt.Println("Cancelado por timeout")
    }
}

// ❌ INCORRECTO: Done channel manual (más verboso)
func improperCancellation() {
    done := make(chan bool)
    
    go func() {
        select {
        case <-done:
            return
        }
    }()
    
    done <- true
}
```

#### 3. Evitar Select Vacío en Producción

```go
// ✓ CORRECTO: Logging explícito
func properMainLoop() {
    go func() {
        for {
            time.Sleep(1 * time.Second)
            fmt.Println("[Background] Activo")
        }
    }()
    
    // Mantener programa activo
    select {}  // Necesario, documentado
}

// ❌ INCORRECTO: Sin goroutines activas
func improperMainLoop() {
    select {}  // Bloquea indefinidamente, sin propósito
}
```

#### 4. Buffering Apropiado

```go
// ✓ CORRECTO: Buffer adecuado para evitar deadlock
func properBuffering() {
    ch := make(chan int, 10)  // Buffer de 10
    
    for i := 0; i < 10; i++ {
        ch <- i  // No bloquea si buffer disponible
    }
}

// ⚠ CUIDADO: Sin buffer requiere receptor
func noBuffering() {
    ch := make(chan int)  // Sin buffer
    
    ch <- 1  // Bloquea hasta que alguien reciba
}
```

### 23.11.2 Antipatrones Comunes

#### Antipatrón 1: Múltiples Sends sin Sincronización

```go
// ❌ ANTIPATRÓN: Race condition
func racyPattern() {
    ch := make(chan int)
    
    go func() {
        ch <- 1  // Goroutine 1
    }()
    
    go func() {
        ch <- 2  // Goroutine 2
    }()
    
    // ¿Cuál se recibe? Indeterminista
    fmt.Println(<-ch)
    // Ambos se enviaron pero solo uno se recibió
}

// ✓ CORRECTO: Recolectar ambos
func properCollection() {
    ch := make(chan int)
    
    go func() {
        ch <- 1
    }()
    
    go func() {
        ch <- 2
    }()
    
    fmt.Println(<-ch)  // Recibir 1
    fmt.Println(<-ch)  // Recibir 2
}
```

#### Antipatrón 2: Deadlock por Select Bloqueante

```go
// ❌ ANTIPATRÓN: Deadlock indefinido
func deadlockSelect() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    select {
    case <-ch1:  // Nunca llegará datos
    case <-ch2:  // Nunca llegará datos
    }
    // Deadlock!
}

// ✓ CORRECTO: Timeout para detectar problemas
func properSelectWithTimeout() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    select {
    case <-ch1:
    case <-ch2:
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout: ningún canal respondió")
    }
}
```

#### Antipatrón 3: Cierre Incorrecto de Channels

```go
// ❌ ANTIPATRÓN: Enviar a channel cerrado (panic)
func improperClose() {
    ch := make(chan int)
    close(ch)
    ch <- 1  // ¡PANIC!
}

// ✓ CORRECTO: Coordinación clara
type Pool struct {
    out  chan int
    stop chan bool
}

func (p *Pool) worker() {
    for {
        select {
        case <-p.stop:
            close(p.out)
            return
        default:
            p.out <- 42
        }
    }
}
```

#### Antipatrón 4: Select sin Default en Loop Crítico

```go
// ❌ ANTIPATRÓN: Puede bloquear indefinidamente
func blockingLoop() {
    eventChan := make(chan string)
    
    for {
        select {
        case event := <-eventChan:
            fmt.Println(event)
        }
        // Sin default ni timeout, bloquea si no hay eventos
    }
}

// ✓ CORRECTO: Incluir timeout o default
func properLoop() {
    eventChan := make(chan string)
    
    for {
        select {
        case event := <-eventChan:
            fmt.Println(event)
            
        case <-time.After(30 * time.Second):
            fmt.Println("[Heartbeat] Activo")
        }
    }
}
```

### 23.11.3 Debugging de Select

```go
import "runtime"

// ✓ Herramienta: Detectar goroutines bloqueadas
func debugBlockedGoroutines() {
    before := runtime.NumGoroutine()
    fmt.Printf("Goroutines antes: %d\n", before)
    
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    // Goroutine bloqueada
    go func() {
        select {
        case <-ch1:
        case <-ch2:
        }
    }()
    
    after := runtime.NumGoroutine()
    fmt.Printf("Goroutines después: %d\n", after)
    
    if after > before {
        fmt.Println("⚠ Goroutine creada pero posiblemente bloqueada")
    }
}

// ✓ Logging para debug
func debugSelect() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "dato"
    }()
    
    fmt.Println("[DEBUG] Esperando select...")
    select {
    case msg := <-ch1:
        fmt.Printf("[DEBUG] ch1: %s\n", msg)
    case msg := <-ch2:
        fmt.Printf("[DEBUG] ch2: %s\n", msg)
    }
    fmt.Println("[DEBUG] Select completado")
}
```

### 23.11.4 Performance Considerations

```go
// ⚠ CUIDADO: Select es O(n) con n canales
// Para n > 10, considerar alternativas

// ✓ MEJOR: Combinar canales antes con fan-in
func efficientMultiplexing() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    ch3 := make(chan int)
    
    // Fan-in: combinar en un canal
    combined := make(chan int)
    
    go func() {
        select {
        case v := <-ch1:
            combined <- v
        case v := <-ch2:
            combined <- v
        case v := <-ch3:
            combined <- v
        }
    }()
    
    // Ahora solo selectar sobre un canal
    select {
    case v := <-combined:
        fmt.Println(v)
    }
}
```

---

## 23.12 Ejercicios Progresivos

### Ejercicio 1: Timeout Simple

**Objetivo:** Implementar una función que realiza una operación con timeout.

**Requisitos:**
- Crear una goroutine que envía un resultado después de 2 segundos
- Usar select con time.After para un timeout de 1 segundo
- Mostrar si fue exitoso o timeout

**Plantilla:**

```go
package main

import (
    "fmt"
    "time"
)

func fetchDataWithTimeout(timeout time.Duration) (string, error) {
    resultChan := make(chan string)
    
    go func() {
        // Simular operación lenta
        time.Sleep(2 * time.Second)
        resultChan <- "Datos obtenidos"
    }()
    
    // TODO: Implementar select con timeout
    select {
    // Completar aquí
    }
}

func main() {
    // Caso 1: Timeout de 1 segundo (debe fallar)
    result, err := fetchDataWithTimeout(1 * time.Second)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Resultado:", result)
    }
    
    // Caso 2: Timeout de 3 segundos (debe éxito)
    result, err = fetchDataWithTimeout(3 * time.Second)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Resultado:", result)
    }
}
```

**Solución esperada:**

```go
func fetchDataWithTimeout(timeout time.Duration) (string, error) {
    resultChan := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        resultChan <- "Datos obtenidos"
    }()
    
    select {
    case result := <-resultChan:
        return result, nil
        
    case <-time.After(timeout):
        return "", fmt.Errorf("timeout excedido")
    }
}
```

---

### Ejercicio 2: Múltiples Sources

**Objetivo:** Select que espera de 3+ fuentes diferentes.

**Requisitos:**
- Crear 3 goroutines que envían diferentes tipos de datos a intervalos variados
- Usar select para recibir de cualquiera que esté listo
- Procesar cada mensaje según su origen
- Detener cuando todas las fuentes se agotan

**Plantilla:**

```go
package main

import (
    "fmt"
    "time"
)

func multiplexingSources() {
    emailsChan := make(chan string)
    notificationsChan := make(chan string)
    alertsChan := make(chan string)
    
    // Goroutine 1: Emails cada 1 segundo
    go func() {
        for i := 1; i <= 3; i++ {
            time.Sleep(1 * time.Second)
            emailsChan <- fmt.Sprintf("Email %d", i)
        }
    }()
    
    // Goroutine 2: Notifications cada 1.5 segundos
    go func() {
        for i := 1; i <= 2; i++ {
            time.Sleep(1500 * time.Millisecond)
            notificationsChan <- fmt.Sprintf("Notif %d", i)
        }
    }()
    
    // Goroutine 3: Alerts cada 2 segundos
    go func() {
        time.Sleep(2 * time.Second)
        alertsChan <- "ALERT: Sistema"
        alertsChan <- "ALERT: Crítico"
    }()
    
    // TODO: Implementar select multiplexing
    // Procesar hasta recibir 7 mensajes totales
}

func main() {
    multiplexingSources()
}
```

**Resultado esperado:**
```
Email 1 recibido
Notif 1 recibido
Email 2 recibido
ALERT: Sistema recibido
Notif 2 recibido
Email 3 recibido
ALERT: Crítico recibido
```

---

### Ejercicio 3: Worker Loop con Timeout

**Objetivo:** Crear un worker que procesa tareas con timeout individual.

**Requisitos:**
- Implementar un type `Worker` que procesa tareas
- Cada tarea tiene un timeout de 500ms
- Si excede timeout, mostrar error y continuar
- Procesar 5 tareas con diferentes duraciones

**Plantilla:**

```go
package main

import (
    "fmt"
    "time"
)

type Task struct {
    ID       int
    Duration time.Duration
}

type Worker struct {
    tasksChan chan Task
}

func (w *Worker) Process(task Task, timeout time.Duration) error {
    resultChan := make(chan string)
    
    go func() {
        time.Sleep(task.Duration)
        resultChan <- fmt.Sprintf("Task %d completada", task.ID)
    }()
    
    // TODO: Implementar select con timeout
    return nil
}

func (w *Worker) Start() {
    for task := range w.tasksChan {
        err := w.Process(task, 500*time.Millisecond)
        if err != nil {
            fmt.Printf("❌ Task %d: %v\n", task.ID, err)
        } else {
            fmt.Printf("✓ Task %d: Completada\n", task.ID)
        }
    }
}

func main() {
    worker := &Worker{tasksChan: make(chan Task)}
    
    go worker.Start()
    
    // Enviar tareas con diferentes duraciones
    tasks := []Task{
        {ID: 1, Duration: 200 * time.Millisecond},
        {ID: 2, Duration: 600 * time.Millisecond},  // Timeout
        {ID: 3, Duration: 300 * time.Millisecond},
        {ID: 4, Duration: 700 * time.Millisecond},  // Timeout
        {ID: 5, Duration: 100 * time.Millisecond},
    }
    
    for _, task := range tasks {
        worker.tasksChan <- task
    }
    
    close(worker.tasksChan)
    time.Sleep(5 * time.Second)
}
```

---

### Ejercicio 4: Graceful Shutdown

**Objetivo:** Implementar shutdown graceful de un servidor concurrente.

**Requisitos:**
- Servidor procesa jobs de un canal
- Recibe señal de shutdown (SIGINT simulado)
- Completa jobs pendientes antes de cerrar
- Mostrar estado durante shutdown

**Plantilla:**

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Server struct {
    jobsChan  chan string
    stopChan  chan bool
    wg        sync.WaitGroup
}

func (s *Server) Start(numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        s.wg.Add(1)
        go s.worker(i)
    }
}

func (s *Server) worker(id int) {
    defer s.wg.Done()
    
    for {
        select {
        // TODO: Recibir jobs o señal de stop
        }
    }
}

func (s *Server) Stop() {
    fmt.Println("\n[SHUTDOWN] Cerrando servidor...")
    fmt.Println("[SHUTDOWN] Deteniendo workers...")
    
    // TODO: Señalizar stop a todos los workers
    // TODO: Esperar a que terminen
    // TODO: Mostrar confirmación
}

func (s *Server) Submit(job string) {
    select {
    case s.jobsChan <- job:
        // OK
    case <-time.After(100 * time.Millisecond):
        fmt.Printf("[SERVIDOR] Job rechazado (cola llena): %s\n", job)
    }
}

func main() {
    server := &Server{
        jobsChan: make(chan string, 5),
        stopChan: make(chan bool),
    }
    
    server.Start(2)
    
    // Simular jobs entrantes
    go func() {
        for i := 1; i <= 10; i++ {
            server.Submit(fmt.Sprintf("Job %d", i))
            time.Sleep(200 * time.Millisecond)
        }
    }()
    
    // Shutdown después de 3 segundos
    time.Sleep(3 * time.Second)
    server.Stop()
}
```

---

### Ejercicio 5: Rate Limiter con Select

**Objetivo:** Implementar un rate limiter que controla throughput usando time.Tick.

**Requisitos:**
- Limitar a 3 solicitudes por segundo
- Mostrar timestamps de cada solicitud procesada
- Usar time.Tick en select para throttling
- Procesar 10 solicitudes

**Plantilla:**

```go
package main

import (
    "fmt"
    "time"
)

type RateLimiter struct {
    requestsChan chan string
    ratePerSec   int
}

func (rl *RateLimiter) Start() {
    // Calcular intervalo entre solicitudes
    interval := time.Second / time.Duration(rl.ratePerSec)
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    processed := 0
    
    for {
        select {
        // TODO: Recibir solicitudes
        // TODO: Esperar tick para procesar
        // TODO: Mostrar timestamp
        }
    }
}

func main() {
    limiter := &RateLimiter{
        requestsChan: make(chan string, 20),
        ratePerSec:   3,  // 3 por segundo
    }
    
    // Procesar en background
    go func() {
        for i := 1; i <= 10; i++ {
            limiter.requestsChan <- fmt.Sprintf("Request %d", i)
            time.Sleep(50 * time.Millisecond)  // Enviar rápido
        }
        close(limiter.requestsChan)
    }()
    
    // Esto debería procesar 3 por segundo
    limiter.Start()
}
```

**Resultado esperado (aproximadamente):**
```
[00:00:00] Request 1 procesada (3 por segundo)
[00:00:00] Request 2 procesada
[00:00:00] Request 3 procesada
[00:00:01] Request 4 procesada  (esperar 1 segundo)
[00:00:01] Request 5 procesada
[00:00:01] Request 6 procesada
...
```

---

## Resumen del Capítulo

El **select statement** es un mecanismo fundamental para concurrencia eficiente en Go:

### Conceptos Clave

1. **Multiplexing**: Esperar múltiples canales sin polling
2. **Non-determinismo**: Si varios casos están listos, Go elige uno al azar
3. **Bloqueante por defecto**: Se bloquea hasta que algún case pueda proceder
4. **Default case**: Para operaciones no bloqueantes
5. **Timeouts**: `time.After` y context para operaciones temporizadas

### Patrones Importantes

- **for-select**: Loop continuo procesando eventos
- **Fan-in/Fan-out**: Combinar o distribuir trabajo
- **Rate limiting**: `time.Tick` para throttling
- **Graceful shutdown**: `done` channels y context
- **Circuit breaker**: Protección contra fallos cascada

### Buenas Prácticas

✓ Usar timeout para evitar bloqueos indefinidos
✓ Coordinar shutdown con goroutines
✓ Evitar select vacío sin contexto
✓ Preferir context para cancelación
✓ Testing considerando non-determinismo

### Antipatrones a Evitar

✗ Select sin mecanismo de timeout
✗ Enviar a channel cerrado
✗ Asumir orden específico de cases
✗ Deadlock por select bloqueante
✗ Cerrar channel desde multiple goroutines

---

## Referencias y Lecturas Adicionales

- **Effective Go**: https://golang.org/doc/effective_go#channels
- **Go Memory Model**: Explica synchronization primitives
- **Concurrency Patterns**: Rob Pike - Go Concurrency Patterns (video)
- **Context Package**: Standard library para cancellation
- **Channels vs Mutexes**: Cuándo usar cada uno

---

**Fin del Capítulo 23**

*Próximo: Capítulo 24 - Patrones Avanzados de Concurrencia*
