# Capítulo 25: Patrones de concurrencia

## Índice

1. [Productor-Consumidor](#251-productor-consumidor)
2. [Pipeline](#252-pipeline)
3. [Fan-Out/Fan-In](#253-fan-outfan-in)
4. [Worker Pool](#254-worker-pool)
5. [Rate Limiting](#255-rate-limiting)
6. [Timeout Pattern](#256-timeout-pattern)
7. [Pub/Sub](#257-pubsub)
8. [Result Generator](#258-result-generator)
9. [Graceful Shutdown](#259-graceful-shutdown)
10. [Error Handling en Concurrencia](#2510-error-handling-en-concurrencia)
11. [Buenas Prácticas y Casos de Uso](#2511-buenas-prácticas-y-casos-de-uso)

---

## 25.1 Productor-Consumidor

### 25.1.1 Concepto Fundamental

El patrón **Productor-Consumidor** es uno de los más fundamentales en programación concurrente. Desacopla la producción de datos del consumo, permitiendo que ambas operaciones ocurran a velocidades diferentes.

**Características:**

- Un productor genera datos
- Un consumidor procesa datos
- Comunican a través de canales
- Buffer regulable para desacoplamiento temporal

**Ventajas:**

- Desacoplamiento de productor y consumidor
- Manejo automático de diferentes velocidades
- Simplifica el flujo de datos
- Facilita testing y mantenimiento

**Desventajas:**

- Requiere sincronización cuidadosa
- Posibles deadlocks si no se diseña bien
- Latencia añadida por los canales

### 25.1.2 Implementación Básica (Unbuffered)

```go
package main

import (
 "fmt"
 "sync"
)

func producerBasico(items chan<- int, wg *sync.WaitGroup) {
 defer wg.Done()
 for i := 1; i <= 5; i++ {
  fmt.Printf("Produciendo: %d\n", i)
  items <- i // Bloqueante hasta que alguien consuma
 }
 close(items)
}

func consumidorBasico(items <-chan int, wg *sync.WaitGroup) {
 defer wg.Done()
 for item := range items {
  fmt.Printf("Consumiendo: %d\n", item)
 }
}

func main() {
 items := make(chan int)
 var wg sync.WaitGroup

 wg.Add(2)
 go producerBasico(items, &wg)
 go consumidorBasico(items, &wg)

 wg.Wait()
 fmt.Println("Completado")
}
```

**Comportamiento:**

- El productor se bloquea hasta que el consumidor lee
- Sin buffer: acoplamiento temporal fuerte
- Ideal para transferencia síncrona de datos

### 25.1.3 Implementación con Buffer

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

func producerConBuffer(items chan<- int, count int, wg *sync.WaitGroup) {
 defer wg.Done()
 for i := 1; i <= count; i++ {
  items <- i * 2
  fmt.Printf("Producido: %d\n", i*2)
 }
 close(items)
}

func consumidorLento(items <-chan int, wg *sync.WaitGroup) {
 defer wg.Done()
 for item := range items {
  time.Sleep(100 * time.Millisecond) // Simula procesamiento lento
  fmt.Printf("Consumido: %d\n", item)
 }
}

func main() {
 // Buffer de 3 elementos
 items := make(chan int, 3)
 var wg sync.WaitGroup

 wg.Add(2)
 go producerConBuffer(items, 5, &wg)
 go consumidorLento(items, &wg)

 wg.Wait()
}
```

**Análisis:**

- Buffer permite que productor avance sin esperar
- Cuando buffer lleno: productor se bloquea
- Mejora throughput en sistemas con velocidades diferentes
- Tamaño de buffer crítico: muy pequeño = bloqueo, muy grande = memoria

### 25.1.4 Múltiples Productores y Consumidores

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

func productor(id int, items chan<- int, wg *sync.WaitGroup) {
 defer wg.Done()
 for i := 1; i <= 3; i++ {
  val := id*100 + i
  items <- val
  fmt.Printf("Productor %d: %d\n", id, val)
  time.Sleep(50 * time.Millisecond)
 }
}

func consumidor(id int, items <-chan int, wg *sync.WaitGroup) {
 defer wg.Done()
 for item := range items {
  fmt.Printf("Consumidor %d: %d\n", id, item)
  time.Sleep(100 * time.Millisecond)
 }
}

func main() {
 items := make(chan int, 5)
 var wg sync.WaitGroup

 // 3 productores
 for i := 1; i <= 3; i++ {
  wg.Add(1)
  go productor(i, items, &wg)
 }

 // 2 consumidores
 for i := 1; i <= 2; i++ {
  wg.Add(1)
  go consumidor(i, items, &wg)
 }

 // Goroutine para cerrar canal cuando todos produzcan
 go func() {
  var producWg sync.WaitGroup
  for i := 1; i <= 3; i++ {
   producWg.Add(1)
   go func(i int) {
    defer producWg.Done()
    productor(i+3, items, &producWg)
   }(i)
  }
  producWg.Wait()
  close(items)
 }()

 wg.Wait()
 fmt.Println("Completado")
}
```

### 25.1.5 Variación con Control de Flujo

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

type Task struct {
 ID   int
 Data string
}

type ProducerConsumer struct {
 tasks       chan Task
 maxBuffer   int
 processedCh chan Task
}

func NewProducerConsumer(maxBuffer int) *ProducerConsumer {
 return &ProducerConsumer{
  tasks:       make(chan Task, maxBuffer),
  maxBuffer:   maxBuffer,
  processedCh: make(chan Task, maxBuffer),
 }
}

func (pc *ProducerConsumer) Produce(count int) {
 go func() {
  for i := 1; i <= count; i++ {
   task := Task{ID: i, Data: fmt.Sprintf("Task-%d", i)}
   pc.tasks <- task
   fmt.Printf("[PROD] Enviado: ID=%d\n", task.ID)
   time.Sleep(50 * time.Millisecond)
  }
  close(pc.tasks)
 }()
}

func (pc *ProducerConsumer) Consume(numWorkers int) {
 var wg sync.WaitGroup
 for w := 0; w < numWorkers; w++ {
  wg.Add(1)
  go func(workerID int) {
   defer wg.Done()
   for task := range pc.tasks {
    fmt.Printf("[CONS %d] Procesando: %s\n", workerID, task.Data)
    time.Sleep(100 * time.Millisecond)
    pc.processedCh <- task
   }
  }(w)
 }

 go func() {
  wg.Wait()
  close(pc.processedCh)
 }()
}

func main() {
 pc := NewProducerConsumer(10)
 pc.Produce(5)
 pc.Consume(2)

 count := 0
 for range pc.processedCh {
  count++
 }
 fmt.Printf("Total procesado: %d\n", count)
}
```

---

## 25.2 Pipeline

### 25.2.1 Concepto Fundamental

Un **pipeline** es una serie de etapas conectadas en secuencia, donde la salida de una etapa es la entrada de la siguiente. Cada etapa procesa datos independientemente.

**Características:**

- Múltiples etapas secuenciales
- Datos fluyen de etapa a etapa
- Cada etapa puede procesar en paralelo
- Típicamente O(n) en tiempo total si bien diseñado

**Casos de uso:**

- Procesamiento de imágenes
- ETL (Extract-Transform-Load)
- Streaming de datos
- Análisis de logs
- Tuberías de compilación

### 25.2.2 Pipeline de 3 Etapas

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
  for i := 1; i <= max; i++ {
   fmt.Printf("[GEN] Generando: %d\n", i)
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
  for num := range in {
   result := num * 2
   fmt.Printf("[DUP] %d -> %d\n", num, result)
   out <- result
  }
  close(out)
 }()
 return out
}

// Etapa 3: Aumentar 10
func aumentar10(in <-chan int) <-chan int {
 out := make(chan int)
 go func() {
  for num := range in {
   result := num + 10
   fmt.Printf("[+10] %d -> %d\n", num, result)
   out <- result
  }
  close(out)
 }()
 return out
}

func main() {
 // Montar pipeline
 gen := generador(5)
 dup := duplicar(gen)
 final := aumentar10(dup)

 // Consumir resultados
 for result := range final {
  fmt.Printf("[OUT] Resultado final: %d\n", result)
 }
}
```

**Salida esperada:**

```
[GEN] Generando: 1
[DUP] 1 -> 2
[+10] 2 -> 12
[OUT] Resultado final: 12
...
```

### 25.2.3 Pipeline con Transformación Genérica

```go
package main

import (
 "fmt"
 "sync"
)

// Etapa genérica que acepta función de transformación
func etapaTransformar(
 in <-chan int,
 nombre string,
 transform func(int) int,
) <-chan int {
 out := make(chan int)
 go func() {
  for val := range in {
   result := transform(val)
   fmt.Printf("[%s] %d -> %d\n", nombre, val, result)
   out <- result
  }
  close(out)
 }()
 return out
}

func main() {
 gen := make(chan int)

 // Iniciar generador
 go func() {
  for i := 1; i <= 3; i++ {
   gen <- i
  }
  close(gen)
 }()

 // Construir pipeline dinámicamente
 resultado := etapaTransformar(gen, "CUADRADO", func(x int) int {
  return x * x
 })
 resultado = etapaTransformar(resultado, "MULTIPLICAR_3", func(x int) int {
  return x * 3
 })
 resultado = etapaTransformar(resultado, "RESTAR_5", func(x int) int {
  return x - 5
 })

 // Consumir
 for val := range resultado {
  fmt.Printf("Final: %d\n", val)
 }
}
```

### 25.2.4 Pipeline con Número de Etapas Variable

```go
package main

import (
 "fmt"
 "sync"
)

type Pipeline struct {
 etapas []func(<-chan int) <-chan int
}

func (p *Pipeline) Agregar(nombre string, transform func(int) int) {
 etapa := func(in <-chan int) <-chan int {
  out := make(chan int)
  go func() {
   for val := range in {
    result := transform(val)
    fmt.Printf("[%s] %d -> %d\n", nombre, val, result)
    out <- result
   }
   close(out)
  }()
  return out
 }
 p.etapas = append(p.etapas, etapa)
}

func (p *Pipeline) Ejecutar(in <-chan int) <-chan int {
 resultado := in
 for _, etapa := range p.etapas {
  resultado = etapa(resultado)
 }
 return resultado
}

func main() {
 pipeline := &Pipeline{}
 pipeline.Agregar("DOBLE", func(x int) int { return x * 2 })
 pipeline.Agregar("CUADRADO", func(x int) int { return x * x })
 pipeline.Agregar("RESTAR_1", func(x int) int { return x - 1 })

 // Entrada
 entrada := make(chan int, 5)
 for i := 1; i <= 5; i++ {
  entrada <- i
 }
 close(entrada)

 // Ejecutar y consumir
 salida := pipeline.Ejecutar(entrada)
 for val := range salida {
  fmt.Printf("Salida: %d\n", val)
 }
}
```

### 25.2.5 Pipeline con Fanout (Múltiples Salidas)

```go
package main

import (
 "fmt"
 "sync"
)

func generador(max int) <-chan int {
 out := make(chan int)
 go func() {
  for i := 1; i <= max; i++ {
   out <- i
  }
  close(out)
 }()
 return out
}

// Fan-out: divide entrada en múltiples salidas
func fanOut(in <-chan int, numSalidas int) []<-chan int {
 salidas := make([]<-chan int, numSalidas)
 for i := 0; i < numSalidas; i++ {
  salidas[i] = make(chan int, 1)
 }

 go func() {
  for val := range in {
   for _, ch := range salidas {
    ch.(chan int) <- val
   }
  }
  for _, ch := range salidas {
   close(ch.(chan int))
  }
 }()

 return salidas
}

func main() {
 entrada := generador(3)
 salidas := fanOut(entrada, 2)

 // Procesar en paralelo
 var wg sync.WaitGroup
 for i, salida := range salidas {
  wg.Add(1)
  go func(id int, ch <-chan int) {
   defer wg.Done()
   for val := range ch {
    fmt.Printf("Salida %d: %d\n", id, val)
   }
  }(i, salida)
 }

 wg.Wait()
}
```

---

## 25.3 Fan-Out/Fan-In

### 25.3.1 Concepto Fundamental

**Fan-Out**: Distribuir trabajo a múltiples workers en paralelo.
**Fan-In**: Recolectar resultados de múltiples sources en un único canal.

**Ventajas:**

- Paralelismo real para CPU-bound tasks
- Mejor utilización de recursos
- Escalabilidad automática
- Manejo eficiente de múltiples fuentes

**Casos de uso:**

- Procesamiento paralelo de datos
- Agregación de múltiples APIs
- Análisis distribuido
- Recolección de eventos

### 25.3.2 Fan-Out Básico

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

func trabajador(id int, trabajos <-chan int, resultados chan<- string) {
 for trabajo := range trabajos {
  fmt.Printf("Worker %d: procesando %d\n", id, trabajo)
  time.Sleep(time.Duration(trabajo*100) * time.Millisecond)
  resultados <- fmt.Sprintf("Worker %d procesó %d", id, trabajo)
 }
}

func main() {
 trabajos := make(chan int, 10)
 resultados := make(chan string)

 // Fan-out: crear 3 workers
 var wg sync.WaitGroup
 for w := 1; w <= 3; w++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   trabajador(id, trabajos, resultados)
  }(w)
 }

 // Enviar trabajos
 go func() {
  for i := 1; i <= 9; i++ {
   trabajos <- i
  }
  close(trabajos)
 }()

 // Esperar a que terminen workers
 go func() {
  wg.Wait()
  close(resultados)
 }()

 // Consumir resultados
 for resultado := range resultados {
  fmt.Printf("Resultado: %s\n", resultado)
 }
}
```

### 25.3.3 Fan-In Básico

```go
package main

import (
 "fmt"
 "sync"
 "time"
 "math/rand"
)

func productor(id int, out chan<- int, duracion time.Duration) {
 ticker := time.NewTicker(duracion)
 defer ticker.Stop()

 for i := 1; i <= 3; i++ {
  <-ticker.C
  val := id*1000 + i
  fmt.Printf("Productor %d: %d\n", id, val)
  out <- val
 }
}

func fanIn(
 productor1, productor2, productor3 <-chan int,
) <-chan int {
 out := make(chan int)
 var wg sync.WaitGroup

 multiplexear := func(in <-chan int) {
  defer wg.Done()
  for val := range in {
   out <- val
  }
 }

 wg.Add(3)
 go multiplexear(productor1)
 go multiplexear(productor2)
 go multiplexear(productor3)

 go func() {
  wg.Wait()
  close(out)
 }()

 return out
}

func main() {
 prod1 := make(chan int)
 prod2 := make(chan int)
 prod3 := make(chan int)

 go productor(1, prod1, 50*time.Millisecond)
 go productor(2, prod2, 75*time.Millisecond)
 go productor(3, prod3, 100*time.Millisecond)

 for val := range fanIn(prod1, prod2, prod3) {
  fmt.Printf("Fan-In recibió: %d\n", val)
 }
}
```

### 25.3.4 Fan-Out/Fan-In Combinado

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

type Trabajo struct {
 ID    int
 Valor int
}

type Resultado struct {
 TrabajoID int
 Valor     int
 Procesado int
}

func procesador(id int, trabajos <-chan Trabajo, resultados chan<- Resultado) {
 for trabajo := range trabajos {
  fmt.Printf("Procesador %d: trabajo %d\n", id, trabajo.ID)
  time.Sleep(100 * time.Millisecond)
  resultados <- Resultado{
   TrabajoID: trabajo.ID,
   Valor:     trabajo.Valor,
   Procesado: trabajo.Valor * 2,
  }
 }
}

func main() {
 trabajos := make(chan Trabajo, 20)
 resultados := make(chan Resultado)

 // Fan-out: 4 procesadores
 var wg sync.WaitGroup
 for p := 1; p <= 4; p++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   procesador(id, trabajos, resultados)
  }(p)
 }

 // Generar trabajos
 go func() {
  for i := 1; i <= 12; i++ {
   trabajos <- Trabajo{ID: i, Valor: i * 10}
  }
  close(trabajos)
 }()

 // Fan-in: recolectar resultados
 go func() {
  wg.Wait()
  close(resultados)
 }()

 // Consumir
 for resultado := range resultados {
  fmt.Printf(
   "Trabajo %d: %d * 2 = %d\n",
   resultado.TrabajoID,
   resultado.Valor,
   resultado.Procesado,
  )
 }
}
```

### 25.3.5 Ventajas vs Desventajas

| Aspecto | Ventajas | Desventajas |
|---------|----------|-------------|
| **Throughput** | Aumenta linealmente con workers | Overhead de sincronización |
| **Latencia** | Puede disminuir por paralelismo | Context switching si too many |
| **Recursos** | Aprovecha multi-core | Más goroutines = más memoria |
| **Escalabilidad** | Fácil agregar workers | Necesita load balancing |

---

## 25.4 Worker Pool

### 25.4.1 Concepto Fundamental

Un **Worker Pool** es un patrón donde un número fijo de workers procesan tareas de una cola compartida. Optimiza:

- **Throughput**: Constante independiente de tamaño de tarea
- **Latencia**: Predecible
- **Recursos**: Control de goroutines
- **Backpressure**: Sistema rechaza sobrecarga

**Diferencias con fan-out:**

- Fan-out: Crea workers dinámicamente
- Worker Pool: Pool fijo, reutilizable
- Pool: Mejor para sistemas de larga duración

### 25.4.2 Implementación Básica

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

type Tarea struct {
 ID   int
 Data string
}

type WorkerPool struct {
 workers   int
 tareas    chan Tarea
 resultados chan Resultado
 wg        sync.WaitGroup
}

type Resultado struct {
 TareaID int
 Salida  string
}

func NewWorkerPool(numWorkers int, bufferTareas int) *WorkerPool {
 return &WorkerPool{
  workers:    numWorkers,
  tareas:     make(chan Tarea, bufferTareas),
  resultados: make(chan Resultado, bufferTareas),
 }
}

func (wp *WorkerPool) Start() {
 for i := 1; i <= wp.workers; i++ {
  wp.wg.Add(1)
  go wp.worker(i)
 }
}

func (wp *WorkerPool) worker(id int) {
 defer wp.wg.Done()
 for tarea := range wp.tareas {
  fmt.Printf("Worker %d: procesando tarea %d\n", id, tarea.ID)
  time.Sleep(100 * time.Millisecond)

  wp.resultados <- Resultado{
   TareaID: tarea.ID,
   Salida:  fmt.Sprintf("Procesado: %s", tarea.Data),
  }
 }
}

func (wp *WorkerPool) EnviarTarea(tarea Tarea) {
 wp.tareas <- tarea
}

func (wp *WorkerPool) Cerrar() {
 close(wp.tareas)
 wp.wg.Wait()
 close(wp.resultados)
}

func (wp *WorkerPool) ObtenerResultados() <-chan Resultado {
 return wp.resultados
}

func main() {
 pool := NewWorkerPool(3, 10)
 pool.Start()

 // Enviar 10 tareas
 go func() {
  for i := 1; i <= 10; i++ {
   pool.EnviarTarea(Tarea{
    ID:   i,
    Data: fmt.Sprintf("tarea-%d", i),
   })
  }
 }()

 // Consumir resultados
 go func() {
  for resultado := range pool.ObtenerResultados() {
   fmt.Printf("Resultado: %s\n", resultado.Salida)
  }
 }()

 pool.Cerrar()
 fmt.Println("Pool cerrado")
}
```

### 25.4.3 Worker Pool con Prioridades

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

type Prioridad int

const (
 Baja    Prioridad = 3
 Normal  Prioridad = 2
 Alta    Prioridad = 1
)

type TareaPrioritizada struct {
 ID       int
 Datos    string
 Prioridad Prioridad
}

type ColaDetalles struct {
 baja   []TareaPrioritizada
 normal []TareaPrioritizada
 alta   []TareaPrioritizada
 mu     sync.Mutex
 cond   *sync.Cond
}

func NewColaDetalles() *ColaDetalles {
 cola := &ColaDetalles{
  baja:   make([]TareaPrioritizada, 0),
  normal: make([]TareaPrioritizada, 0),
  alta:   make([]TareaPrioritizada, 0),
 }
 cola.cond = sync.NewCond(&cola.mu)
 return cola
}

func (c *ColaDetalles) Agregar(tarea TareaPrioritizada) {
 c.mu.Lock()
 defer c.mu.Unlock()

 switch tarea.Prioridad {
 case Alta:
  c.alta = append(c.alta, tarea)
 case Normal:
  c.normal = append(c.normal, tarea)
 case Baja:
  c.baja = append(c.baja, tarea)
 }
 c.cond.Signal()
}

func (c *ColaDetalles) ObtenerSiguiente() (TareaPrioritizada, bool) {
 c.mu.Lock()
 defer c.mu.Unlock()

 for len(c.alta) == 0 && len(c.normal) == 0 && len(c.baja) == 0 {
  c.cond.Wait()
 }

 if len(c.alta) > 0 {
  tarea := c.alta[0]
  c.alta = c.alta[1:]
  return tarea, true
 }
 if len(c.normal) > 0 {
  tarea := c.normal[0]
  c.normal = c.normal[1:]
  return tarea, true
 }
 if len(c.baja) > 0 {
  tarea := c.baja[0]
  c.baja = c.baja[1:]
  return tarea, true
 }

 return TareaPrioritizada{}, false
}

func procesador(id int, cola *ColaDetalles, done <-chan struct{}, wg *sync.WaitGroup) {
 defer wg.Done()
 for {
  select {
  case <-done:
   fmt.Printf("Procesador %d: terminando\n", id)
   return
  default:
   tarea, ok := cola.ObtenerSiguiente()
   if !ok {
    time.Sleep(100 * time.Millisecond)
    continue
   }
   fmt.Printf(
    "Procesador %d: tarea %d (prioridad %d)\n",
    id, tarea.ID, tarea.Prioridad,
   )
   time.Sleep(50 * time.Millisecond)
  }
 }
}

func main() {
 cola := NewColaDetalles()
 done := make(chan struct{})
 var wg sync.WaitGroup

 // 2 procesadores
 for i := 1; i <= 2; i++ {
  wg.Add(1)
  go procesador(i, cola, done, &wg)
 }

 // Agregar tareas con diferentes prioridades
 for i := 1; i <= 6; i++ {
  var prio Prioridad
  switch i % 3 {
  case 0:
   prio = Alta
  case 1:
   prio = Normal
  default:
   prio = Baja
  }
  cola.Agregar(TareaPrioritizada{
   ID:        i,
   Datos:     fmt.Sprintf("tarea-%d", i),
   Prioridad: prio,
  })
 }

 time.Sleep(2 * time.Second)
 close(done)
 wg.Wait()
}
```

### 25.4.4 Worker Pool Dinámico (Escalable)

```go
package main

import (
 "fmt"
 "sync"
 "sync/atomic"
 "time"
)

type PoolDinamico struct {
 tareas        chan func()
 workersActivos int32
 maxWorkers    int
 mu            sync.Mutex
 wg            sync.WaitGroup
}

func NewPoolDinamico(minWorkers, maxWorkers int) *PoolDinamico {
 pool := &PoolDinamico{
  tareas:     make(chan func(), 100),
  maxWorkers: maxWorkers,
 }

 // Iniciar con workers mínimos
 for i := 0; i < minWorkers; i++ {
  pool.agregarWorker()
 }

 return pool
}

func (p *PoolDinamico) agregarWorker() {
 atomic.AddInt32(&p.workersActivos, 1)
 p.wg.Add(1)

 go func() {
  defer p.wg.Done()
  defer atomic.AddInt32(&p.workersActivos, -1)

  timeout := time.NewTimer(5 * time.Second)
  defer timeout.Stop()

  for {
   timeout.Reset(5 * time.Second)
   select {
   case tarea, ok := <-p.tareas:
    if !ok {
     return
    }
    tarea()
    timeout.Reset(5 * time.Second)
   case <-timeout.C:
    // Worker ocioso más de 5s: muere
    return
   }
  }
 }()
}

func (p *PoolDinamico) Enviar(tarea func()) {
 p.mu.Lock()
 active := atomic.LoadInt32(&p.workersActivos)
 p.mu.Unlock()

 if int(active) < p.maxWorkers && len(p.tareas) > int(active)*2 {
  p.agregarWorker()
 }

 p.tareas <- tarea
}

func (p *PoolDinamico) Esperar() {
 close(p.tareas)
 p.wg.Wait()
}

func main() {
 pool := NewPoolDinamico(2, 5)

 // Enviar 20 tareas
 for i := 1; i <= 20; i++ {
  id := i
  pool.Enviar(func() {
   fmt.Printf("Tarea %d\n", id)
   time.Sleep(200 * time.Millisecond)
  })
 }

 pool.Esperar()
 fmt.Println("Completado")
}
```

---

## 25.5 Rate Limiting

### 25.5.1 Concepto Fundamental

**Rate Limiting** controla la velocidad de operaciones para:

- Proteger sistemas de sobrecarga
- Prevenir DDoS
- Mantener SLA
- Controlar recursos
- Implementar throttling

**Estrategias:**

- Token Bucket: Permite ráfagas
- Sliding Window: Suavizado temporal
- Leaky Bucket: Control estricto
- Adaptive: Ajusta dinámicamente

### 25.5.2 Token Bucket Básico

```go
package main

import (
 "fmt"
 "time"
)

type TokenBucket struct {
 tokens    int
 capacity  int
 refillRate time.Duration
 mu        chan struct{}
}

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
 tb := &TokenBucket{
  tokens:    capacity,
  capacity:  capacity,
  refillRate: refillRate,
  mu:        make(chan struct{}, 1),
 }

 // Goroutine para rellenar tokens
 go func() {
  ticker := time.NewTicker(refillRate)
  defer ticker.Stop()

  for range ticker.C {
   tb.mu <- struct{}{} // Adquirir lock
   if tb.tokens < tb.capacity {
    tb.tokens++
    fmt.Printf("Token añadido, total: %d\n", tb.tokens)
   }
   <-tb.mu // Liberar lock
  }
 }()

 return tb
}

func (tb *TokenBucket) AcumularToken(n int) bool {
 tb.mu <- struct{}{} // Adquirir lock
 defer func() { <-tb.mu }() // Liberar lock

 if tb.tokens >= n {
  tb.tokens -= n
  return true
 }
 return false
}

func main() {
 tb := NewTokenBucket(5, 500*time.Millisecond)

 for i := 1; i <= 10; i++ {
  if tb.AcumularToken(1) {
   fmt.Printf("Solicitud %d: PERMITIDA\n", i)
  } else {
   fmt.Printf("Solicitud %d: DENEGADA\n", i)
  }
  time.Sleep(200 * time.Millisecond)
 }
}
```

### 25.5.3 Rate Limiter con time.Ticker

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

type RateLimiter struct {
 limiter <-chan time.Time
}

func NewRateLimiter(tasaOperacionesPorSegundo float64) *RateLimiter {
 intervalo := time.Duration(float64(time.Second) / tasaOperacionesPorSegundo)
 limiter := time.Tick(intervalo)
 return &RateLimiter{limiter: limiter}
}

func (rl *RateLimiter) Esperar() {
 <-rl.limiter
}

func main() {
 // 5 operaciones por segundo
 limiter := NewRateLimiter(5)

 for i := 1; i <= 15; i++ {
  limiter.Esperar()
  fmt.Printf("%d: %s\n", i, time.Now().Format("15:04:05.000"))
 }
}
```

### 25.5.4 Token Bucket Avanzado

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

type AdvancedTokenBucket struct {
 tokens        float64
 maxTokens     float64
 refillRate    float64 // tokens por segundo
 lastRefillTime time.Time
 mu            sync.Mutex
}

func NewAdvancedTokenBucket(maxTokens, refillRate float64) *AdvancedTokenBucket {
 return &AdvancedTokenBucket{
  tokens:         maxTokens,
  maxTokens:      maxTokens,
  refillRate:     refillRate,
  lastRefillTime: time.Now(),
 }
}

func (atb *AdvancedTokenBucket) Permitir(tokensNeeded float64) bool {
 atb.mu.Lock()
 defer atb.mu.Unlock()

 // Calcular tokens generados desde último acceso
 ahora := time.Now()
 tiempoTranscurrido := ahora.Sub(atb.lastRefillTime).Seconds()
 atb.tokens = min(
  atb.maxTokens,
  atb.tokens+tiempoTranscurrido*atb.refillRate,
 )
 atb.lastRefillTime = ahora

 if atb.tokens >= tokensNeeded {
  atb.tokens -= tokensNeeded
  return true
 }

 return false
}

func (atb *AdvancedTokenBucket) Esperar(tokensNeeded float64) {
 for {
  if atb.Permitir(tokensNeeded) {
   return
  }
  atb.mu.Lock()
  tiempoEspera := time.Duration(
   (tokensNeeded - atb.tokens) / atb.refillRate * float64(time.Second),
  )
  atb.mu.Unlock()
  time.Sleep(tiempoEspera)
 }
}

func min(a, b float64) float64 {
 if a < b {
  return a
 }
 return b
}

func main() {
 // 10 tokens por segundo, máximo 20
 bucket := NewAdvancedTokenBucket(20, 10)

 var wg sync.WaitGroup
 for i := 1; i <= 5; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   bucket.Esperar(1)
   fmt.Printf("Goroutine %d: %s\n", id, time.Now().Format("15:04:05.000"))
  }(i)
 }

 wg.Wait()
}
```

### 25.5.5 Backpressure

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

func productor(out chan<- int, limit <-chan struct{}) {
 for i := 1; i <= 10; i++ {
  <-limit // Esperar permiso
  out <- i
  fmt.Printf("Producido: %d\n", i)
  time.Sleep(100 * time.Millisecond)
 }
 close(out)
}

func consumidor(in <-chan int, limit chan<- struct{}, wg *sync.WaitGroup) {
 defer wg.Done()
 for val := range in {
  fmt.Printf("Consumiendo: %d\n", val)
  time.Sleep(200 * time.Millisecond)
  limit <- struct{}{} // Permitir siguiente producción
 }
}

func main() {
 out := make(chan int)
 limit := make(chan struct{}, 3) // Buffer de 3

 // Inicializar con 3 permisos
 for i := 0; i < 3; i++ {
  limit <- struct{}{}
 }

 var wg sync.WaitGroup
 wg.Add(1)
 go productor(out, limit)
 go consumidor(in: out, limit: limit, wg: &wg)

 wg.Wait()
}
```

---

## 25.6 Timeout Pattern

### 25.6.1 Concepto Fundamental

El patrón **Timeout** evita que goroutines esperen indefinidamente. Implementa:

- Deadlines
- Circuit breaker
- Recuperación de fallos
- Limitación de recursos

### 25.6.2 Timeout Básico con select

```go
package main

import (
 "fmt"
 "time"
)

func operacionLenta() <-chan string {
 result := make(chan string)
 go func() {
  time.Sleep(3 * time.Second)
  result <- "Completada"
 }()
 return result
}

func main() {
 select {
 case res := <-operacionLenta():
  fmt.Println("Resultado:", res)
 case <-time.After(1 * time.Second):
  fmt.Println("Timeout: operación tardó demasiado")
 }
}
```

### 25.6.3 Context con Timeout

```go
package main

import (
 "context"
 "fmt"
 "time"
)

func tareaConTiempo(ctx context.Context, id int) error {
 select {
 case <-time.After(2 * time.Second):
  fmt.Printf("Tarea %d: completada\n", id)
  return nil
 case <-ctx.Done():
  fmt.Printf("Tarea %d: cancelada\n", id)
  return ctx.Err()
 }
}

func main() {
 ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
 defer cancel()

 err := tareaConTiempo(ctx, 1)
 if err != nil {
  fmt.Println("Error:", err)
 }
}
```

### 25.6.4 Circuit Breaker

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

type Estado int

const (
 Cerrado Estado = iota
 Abierto
 Semiabiero
)

type CircuitBreaker struct {
 estado    Estado
 fallos    int
 exitos    int
 umbral    int
 mu        sync.RWMutex
 ultimoFallo time.Time
}

func NewCircuitBreaker(umbral int) *CircuitBreaker {
 return &CircuitBreaker{
  estado: Cerrado,
  umbral: umbral,
 }
}

func (cb *CircuitBreaker) Llamar(f func() error) error {
 cb.mu.Lock()
 defer cb.mu.Unlock()

 if cb.estado == Abierto {
  if time.Since(cb.ultimoFallo) > 30*time.Second {
   cb.estado = Semiabiero
   cb.exitos = 0
  } else {
   return fmt.Errorf("circuito abierto")
  }
 }

 err := f()

 if err != nil {
  cb.fallos++
  cb.ultimoFallo = time.Now()
  if cb.fallos >= cb.umbral {
   cb.estado = Abierto
  }
  return err
 }

 cb.fallos = 0
 if cb.estado == Semiabiero {
  cb.exitos++
  if cb.exitos >= 3 {
   cb.estado = Cerrado
  }
 }

 return nil
}

func main() {
 cb := NewCircuitBreaker(3)
 contador := 0

 for i := 1; i <= 10; i++ {
  err := cb.Llamar(func() error {
   contador++
   if contador <= 3 {
    return fmt.Errorf("fallo")
   }
   return nil
  })

  if err != nil {
   fmt.Printf("Intento %d: ERROR - %v\n", i, err)
  } else {
   fmt.Printf("Intento %d: ÉXITO\n", i)
  }

  time.Sleep(100 * time.Millisecond)
 }
}
```

---

## 25.7 Pub/Sub

### 25.7.1 Concepto Fundamental

**Pub/Sub** desacopla completamente productores (publishers) de consumidores (subscribers). Ventajas:

- Escalabilidad horizontal
- Desacoplamiento completo
- Fácil agregar nuevos suscriptores
- Broadcast de eventos

### 25.7.2 Pub/Sub Básico

```go
package main

import (
 "fmt"
 "sync"
)

type Evento struct {
 Tipo    string
 Datos   interface{}
}

type EventBus struct {
 suscriptores map[string][]chan Evento
 mu           sync.RWMutex
}

func NewEventBus() *EventBus {
 return &EventBus{
  suscriptores: make(map[string][]chan Evento),
 }
}

func (eb *EventBus) Suscribir(tipo string) <-chan Evento {
 eb.mu.Lock()
 defer eb.mu.Unlock()

 ch := make(chan Evento, 10)
 eb.suscriptores[tipo] = append(eb.suscriptores[tipo], ch)
 return ch
}

func (eb *EventBus) Publicar(evento Evento) {
 eb.mu.RLock()
 suscriptores, ok := eb.suscriptores[evento.Tipo]
 eb.mu.RUnlock()

 if !ok {
  return
 }

 for _, ch := range suscriptores {
  select {
  case ch <- evento:
  default:
   fmt.Println("Suscriptor sin espacio")
  }
 }
}

func main() {
 bus := NewEventBus()

 // Suscriptor 1
 eventos1 := bus.Suscribir("usuario")
 go func() {
  for evento := range eventos1 {
   fmt.Printf("Sub1: %v\n", evento)
  }
 }()

 // Suscriptor 2
 eventos2 := bus.Suscribir("usuario")
 go func() {
  for evento := range eventos2 {
   fmt.Printf("Sub2: %v\n", evento)
  }
 }()

 // Publicador
 for i := 1; i <= 3; i++ {
  bus.Publicar(Evento{
   Tipo:  "usuario",
   Datos: fmt.Sprintf("evento-%d", i),
  })
 }

 // (Nota: En una app real, usar sync.WaitGroup para esperar)
 fmt.Println("Publicado")
}
```

### 25.7.3 Pub/Sub Avanzado con Unsubscribe

```go
package main

import (
 "fmt"
 "sync"
)

type Suscriptor struct {
 ID   string
 ch   chan Evento
 quit chan struct{}
}

type EventBusAvanzado struct {
 suscriptores map[string][]*Suscriptor
 mu           sync.RWMutex
 publicar     chan Evento
}

func NewEventBusAvanzado() *EventBusAvanzado {
 eb := &EventBusAvanzado{
  suscriptores: make(map[string][]*Suscriptor),
  publicar:     make(chan Evento, 100),
 }
 go eb.demultiplexar()
 return eb
}

func (eb *EventBusAvanzado) demultiplexar() {
 for evento := range eb.publicar {
  eb.mu.RLock()
  suscriptores, ok := eb.suscriptores[evento.Tipo]
  eb.mu.RUnlock()

  if !ok {
   continue
  }

  for _, sub := range suscriptores {
   select {
   case sub.ch <- evento:
   case <-sub.quit:
    // Suscriptor fue desinscrito
   }
  }
 }
}

func (eb *EventBusAvanzado) Suscribir(tipo, id string) <-chan Evento {
 eb.mu.Lock()
 defer eb.mu.Unlock()

 sub := &Suscriptor{
  ID:   id,
  ch:   make(chan Evento, 10),
  quit: make(chan struct{}),
 }

 eb.suscriptores[tipo] = append(eb.suscriptores[tipo], sub)
 return sub.ch
}

func (eb *EventBusAvanzado) Desuscribir(tipo, id string) {
 eb.mu.Lock()
 defer eb.mu.Unlock()

 suscriptores, ok := eb.suscriptores[tipo]
 if !ok {
  return
 }

 for i, sub := range suscriptores {
  if sub.ID == id {
   close(sub.quit)
   eb.suscriptores[tipo] = append(
    suscriptores[:i],
    suscriptores[i+1:]...,
   )
   break
  }
 }
}

func (eb *EventBusAvanzado) Publicar(evento Evento) {
 eb.publicar <- evento
}

func main() {
 bus := NewEventBusAvanzado()

 ch1 := bus.Suscribir("news", "sub1")
 ch2 := bus.Suscribir("news", "sub2")

 go func() {
  for evento := range ch1 {
   fmt.Printf("Sub1: %s\n", evento.Datos)
  }
 }()

 go func() {
  for evento := range ch2 {
   fmt.Printf("Sub2: %s\n", evento.Datos)
  }
 }()

 bus.Publicar(Evento{"news", "breaking news"})
 bus.Desuscribir("news", "sub1")
 bus.Publicar(Evento{"news", "second news"})

 fmt.Println("Completado")
}
```

---

## 25.8 Result Generator

### 25.8.1 Concepto Fundamental

Un **Result Generator** es una goroutine que genera valores lazily (bajo demanda). Útil para:

- Secuencias infinitas
- Computación perezosa
- Generadores en pipeline
- Optimización de memoria

### 25.8.2 Generador Simple

```go
package main

import "fmt"

func generador(inicio, fin int) <-chan int {
 out := make(chan int)
 go func() {
  for i := inicio; i <= fin; i++ {
   out <- i
  }
  close(out)
 }()
 return out
}

func main() {
 for valor := range generador(1, 5) {
  fmt.Println(valor)
 }
}
```

### 25.8.3 Generador de Fibonacci

```go
package main

import "fmt"

func fibonacci() <-chan int {
 out := make(chan int)
 go func() {
  a, b := 0, 1
  for {
   out <- a
   a, b = b, a+b
  }
 }()
 return out
}

func main() {
 gen := fibonacci()
 for i := 0; i < 10; i++ {
  fmt.Println(<-gen)
 }
 close(gen) // Cerrar canal si necesario
}
```

### 25.8.4 Generador Filtrado

```go
package main

import "fmt"

func generador(inicio, fin int) <-chan int {
 out := make(chan int)
 go func() {
  for i := inicio; i <= fin; i++ {
   out <- i
  }
  close(out)
 }()
 return out
}

func filtrar(
 in <-chan int,
 predicado func(int) bool,
) <-chan int {
 out := make(chan int)
 go func() {
  for valor := range in {
   if predicado(valor) {
    out <- valor
   }
  }
  close(out)
 }()
 return out
}

func main() {
 // Números pares del 1 al 20
 pares := filtrar(
  generador(1, 20),
  func(n int) bool { return n%2 == 0 },
 )

 for n := range pares {
  fmt.Println(n)
 }
}
```

### 25.8.5 Generador de Pares Key-Value

```go
package main

import (
 "fmt"
 "sync"
)

type Par struct {
 Clave   string
 Valor   interface{}
}

func generadorDatos() <-chan Par {
 out := make(chan Par)
 go func() {
  datos := map[string]interface{}{
   "nombre": "Juan",
   "edad":   30,
   "ciudad": "Madrid",
  }

  for k, v := range datos {
   out <- Par{k, v}
  }
  close(out)
 }()
 return out
}

func main() {
 var wg sync.WaitGroup
 resultado := make(map[string]interface{})
 var mu sync.Mutex

 for par := range generadorDatos() {
  wg.Add(1)
  go func(p Par) {
   defer wg.Done()
   mu.Lock()
   resultado[p.Clave] = p.Valor
   mu.Unlock()
   fmt.Printf("Procesado: %s = %v\n", p.Clave, p.Valor)
  }(par)
 }

 wg.Wait()
 fmt.Printf("Resultado: %v\n", resultado)
}
```

---

## 25.9 Graceful Shutdown

### 25.9.1 Concepto Fundamental

**Graceful Shutdown** asegura que:

- Todas las tareas en progreso se completan
- Nuevas tareas se rechazan
- Recursos se liberan correctamente
- No hay corrupción de datos

### 25.9.2 Shutdown Básico

```go
package main

import (
 "fmt"
 "os"
 "os/signal"
 "sync"
 "syscall"
 "time"
)

func trabajador(id int, tareas <-chan int, wg *sync.WaitGroup) {
 defer wg.Done()
 for tarea := range tareas {
  fmt.Printf("Trabajador %d: procesando %d\n", id, tarea)
  time.Sleep(500 * time.Millisecond)
 }
 fmt.Printf("Trabajador %d: terminado\n", id)
}

func main() {
 tareas := make(chan int, 10)
 signals := make(chan os.Signal, 1)
 signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

 var wg sync.WaitGroup

 // 3 trabajadores
 for w := 1; w <= 3; w++ {
  wg.Add(1)
  go trabajador(w, tareas, &wg)
 }

 // Enviar tareas en goroutine
 go func() {
  for i := 1; i <= 20; i++ {
   select {
   case tareas <- i:
    fmt.Printf("Tarea %d encolada\n", i)
   case <-signals:
    fmt.Println("\nSeñal recibida, cerrando canal")
    close(tareas)
    return
   }
   time.Sleep(200 * time.Millisecond)
  }
  close(tareas)
 }()

 wg.Wait()
 fmt.Println("Shutdown completo")
}
```

### 25.9.3 Shutdown con Context

```go
package main

import (
 "context"
 "fmt"
 "os"
 "os/signal"
 "sync"
 "syscall"
 "time"
)

func servidor(ctx context.Context, wg *sync.WaitGroup) {
 defer wg.Done()

 for i := 1; i <= 10; i++ {
  select {
  case <-ctx.Done():
   fmt.Println("Servidor: shutdown recibido")
   return
  default:
   fmt.Printf("Servidor: operación %d\n", i)
   time.Sleep(500 * time.Millisecond)
  }
 }
}

func main() {
 ctx, cancel := context.WithCancel(context.Background())
 defer cancel()

 signals := make(chan os.Signal, 1)
 signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

 var wg sync.WaitGroup

 // Múltiples servidores
 for s := 1; s <= 3; s++ {
  wg.Add(1)
  go servidor(ctx, &wg)
 }

 // Goroutine para signal handler
 go func() {
  sig := <-signals
  fmt.Printf("\nSeñal recibida: %v\n", sig)
  cancel()
 }()

 wg.Wait()
 fmt.Println("Aplicación finalizada")
}
```

### 25.9.4 Shutdown con Timeout

```go
package main

import (
 "context"
 "fmt"
 "sync"
 "time"
)

type Aplicacion struct {
 tareas chan func()
 done   chan struct{}
}

func NewAplicacion() *Aplicacion {
 return &Aplicacion{
  tareas: make(chan func(), 100),
  done:   make(chan struct{}),
 }
}

func (a *Aplicacion) Iniciar(ctx context.Context) {
 go func() {
  for {
   select {
   case tarea := <-a.tareas:
    if tarea != nil {
     tarea()
    }
   case <-ctx.Done():
    close(a.done)
    return
   }
  }
 }()
}

func (a *Aplicacion) Encolar(tarea func()) error {
 select {
 case a.tareas <- tarea:
  return nil
 default:
  return fmt.Errorf("cola llena")
 }
}

func (a *Aplicacion) Cerrar(timeout time.Duration) error {
 ctx, cancel := context.WithTimeout(
  context.Background(),
  timeout,
 )
 defer cancel()

 // Detener aceptación de nuevas tareas
 close(a.tareas)

 // Esperar a que terminen o timeout
 select {
 case <-a.done:
  fmt.Println("Aplicación cerrada correctamente")
  return nil
 case <-ctx.Done():
  return fmt.Errorf("timeout en shutdown")
 }
}

func main() {
 app := NewAplicacion()
 ctx, cancel := context.WithCancel(context.Background())
 defer cancel()

 app.Iniciar(ctx)

 // Encolar tareas
 for i := 1; i <= 5; i++ {
  id := i
  app.Encolar(func() {
   fmt.Printf("Tarea %d ejecutada\n", id)
   time.Sleep(200 * time.Millisecond)
  })
 }

 // Simular fin de tareas
 time.Sleep(1500 * time.Millisecond)

 // Shutdown graceful
 err := app.Cerrar(5 * time.Second)
 if err != nil {
  fmt.Println("Error:", err)
 }
}
```

---

## 25.10 Error Handling en Concurrencia

### 25.10.1 Concepto Fundamental

En concurrencia, los errores necesitan propagarse sin bloquear el sistema. Estrategias:

- Canal de errores
- Agregación de errores
- Recuperación automática
- Logging centralizado

### 25.10.2 Canal de Errores

```go
package main

import (
 "fmt"
 "sync"
)

type Resultado struct {
 Valor int
 Error error
}

func procesador(
 id int,
 tareas <-chan int,
 resultados chan<- Resultado,
) {
 for tarea := range tareas {
  var resultado Resultado

  if tarea%3 == 0 {
   resultado = Resultado{Error: fmt.Errorf("tarea %d divisible por 3", tarea)}
  } else {
   resultado = Resultado{Valor: tarea * 2}
  }

  resultados <- resultado
 }
}

func main() {
 tareas := make(chan int, 5)
 resultados := make(chan Resultado, 5)

 // 2 procesadores
 var wg sync.WaitGroup
 for w := 1; w <= 2; w++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   procesador(w, tareas, resultados)
  }()
 }

 // Enviar tareas
 for i := 1; i <= 6; i++ {
  tareas <- i
 }
 close(tareas)

 // Esperar y procesar resultados
 go func() {
  wg.Wait()
  close(resultados)
 }()

 // Consumir
 for resultado := range resultados {
  if resultado.Error != nil {
   fmt.Printf("ERROR: %v\n", resultado.Error)
  } else {
   fmt.Printf("OK: %d\n", resultado.Valor)
  }
 }
}
```

### 25.10.3 Agregación de Errores

```go
package main

import (
 "fmt"
 "sync"
)

type ErroresAgregados struct {
 errores []error
 mu      sync.Mutex
}

func (ea *ErroresAgregados) Agregar(err error) {
 if err == nil {
  return
 }
 ea.mu.Lock()
 defer ea.mu.Unlock()
 ea.errores = append(ea.errores, err)
}

func (ea *ErroresAgregados) ObtenerTodos() []error {
 ea.mu.Lock()
 defer ea.mu.Unlock()
 copiar := make([]error, len(ea.errores))
 copy(copiar, ea.errores)
 return copiar
}

func (ea *ErroresAgregados) Hay() bool {
 ea.mu.Lock()
 defer ea.mu.Unlock()
 return len(ea.errores) > 0
}

func procesadorConError(
 id int,
 tareas <-chan int,
 errores *ErroresAgregados,
 wg *sync.WaitGroup,
) {
 defer wg.Done()

 for tarea := range tareas {
  if tarea%2 == 0 {
   errores.Agregar(fmt.Errorf(
    "procesador %d: número par %d",
    id, tarea,
   ))
  } else {
   fmt.Printf("Procesador %d: procesó %d\n", id, tarea)
  }
 }
}

func main() {
 tareas := make(chan int, 10)
 errores := &ErroresAgregados{}
 var wg sync.WaitGroup

 // 3 procesadores
 for p := 1; p <= 3; p++ {
  wg.Add(1)
  go procesadorConError(p, tareas, errores, &wg)
 }

 // Enviar tareas
 for i := 1; i <= 9; i++ {
  tareas <- i
 }
 close(tareas)

 wg.Wait()

 // Reportar errores
 if errores.Hay() {
  fmt.Println("Errores encontrados:")
  for _, err := range errores.ObtenerTodos() {
   fmt.Printf("  - %v\n", err)
  }
 }
}
```

### 25.10.4 Recuperación Automática

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

func tareaConReintentos(
 id int,
 maxReintentos int,
 operacion func() error,
) error {
 var err error

 for intento := 1; intento <= maxReintentos; intento++ {
  err = operacion()
  if err == nil {
   return nil
  }

  fmt.Printf(
   "Tarea %d: intento %d falló (%v)\n",
   id, intento, err,
  )

  if intento < maxReintentos {
   espera := time.Duration(intento*100) * time.Millisecond
   time.Sleep(espera)
  }
 }

 return fmt.Errorf("tarea %d falló después de %d reintentos", id, maxReintentos)
}

func main() {
 contador := 0
 operacion := func() error {
  contador++
  if contador < 3 {
   return fmt.Errorf("fallo temporal")
  }
  return nil
 }

 err := tareaConReintentos(1, 5, operacion)
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 } else {
  fmt.Println("Éxito después de reintentos")
 }
}
```

---

## 25.11 Buenas Prácticas y Casos de Uso

### 25.11.1 Matriz de Patrones

| Patrón | Caso de Uso | Ventajas | Desventajas |
|--------|-------------|----------|-------------|
| **Productor-Consumidor** | Desacoplamiento temporal | Simple, directo | Difícil escalar |
| **Pipeline** | Transformación en etapas | Flujo claro, composable | Latencia aditiva |
| **Fan-Out/Fan-In** | Paralelismo CPU-bound | Buen throughput | Sincronización compleja |
| **Worker Pool** | Tareas heterogéneas | Control de recursos | Setup inicial |
| **Rate Limiting** | Protección de sistema | Backpressure nativo | Overhead |
| **Timeout** | Operaciones con deadline | Recuperación rápida | Puede perder datos |
| **Pub/Sub** | Eventos distribuidos | Desacoplamiento total | Debugging difícil |

### 25.11.2 Criterios de Selección

```
¿Necesitas procesar items de una sola fuente?
├─ Sí: Productor-Consumidor
└─ No: ¿Muchas fuentes?
   ├─ Sí: Fan-In / Pub/Sub
   └─ No: ¿Múltiples etapas?
      ├─ Sí: Pipeline
      └─ No: Worker Pool

¿Necesitas limitar throughput?
└─ Sí: Rate Limiting

¿Necesitas deadline?
└─ Sí: Timeout Pattern / Context

¿Sistema distribuido?
└─ Sí: Pub/Sub, considerar Redis/NATS
```

### 25.11.3 Antipatrones Comunes

**1. Goroutines sin coordinación**

```go
// ❌ MALO: Goroutines sin wait
go func() { /* ... */ }()
// Función puede terminar antes de goroutines

// ✅ BUENO: Con sync.WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
 defer wg.Done()
 // ...
}()
wg.Wait()
```

**2. Sin backpressure**

```go
// ❌ MALO: Canal sin buffer → deadlock potencial
ch := make(chan int)
ch <- 1 // Bloqueado para siempre si nadie consume

// ✅ BUENO: Buffer o lectura
ch := make(chan int, 1)
ch <- 1
```

**3. Cerrando canal compartido**

```go
// ❌ MALO: Múltiples goroutines cerrando
close(ch) // panic si otro cierra

// ✅ BUENO: Un único responsable
// Solo el productor cierra
```

**4. Sin timeout**

```go
// ❌ MALO: Deadlock potencial
<-ch

// ✅ BUENO: Con deadline
select {
case <-ch:
case <-time.After(5 * time.Second):
 log.Fatal("Timeout")
}
```

### 25.11.4 Benchmarking de Patrones

```go
package main

import (
 "fmt"
 "testing"
 "time"
)

// Benchmark productor-consumidor
func BenchmarkProductorConsumidor(b *testing.B) {
 for i := 0; i < b.N; i++ {
  ch := make(chan int)
  go func() {
   for j := 0; j < 100; j++ {
    ch <- j
   }
   close(ch)
  }()

  for range ch {
  }
 }
}

// Benchmark worker pool
func BenchmarkWorkerPool(b *testing.B) {
 for i := 0; i < b.N; i++ {
  tareas := make(chan int, 10)
  resultados := make(chan int, 10)

  // 4 workers
  for w := 0; w < 4; w++ {
   go func() {
    for t := range tareas {
     resultados <- t * 2
    }
   }()
  }

  go func() {
   for j := 0; j < 100; j++ {
    tareas <- j
   }
   close(tareas)
  }()

  for j := 0; j < 100; j++ {
   <-resultados
  }
 }
}

func main() {
 fmt.Println("Benchmarks ejecutados manualmente:")
 fmt.Println("Ver output con: go test -bench=. -benchmem")
}
```

### 25.11.5 Casos de Uso Reales

**1. Web Server con Rate Limiting**

```go
type WebServer struct {
 limiter *RateLimiter
 workers int
 tareas  chan Request
}

func (ws *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 ws.limiter.Esperar()
 ws.tareas <- NewRequest(r)
}
```

**2. Data Pipeline**

```go
resultados := pipeline.
 Leer("datos.csv").
 Transformar(parseJSON).
 Filtrar(esValido).
 Mapear(calcular).
 Escribir("salida.csv")
```

**3. Event Distribution**

```go
type EventSystem struct {
 bus *EventBus
}

func (es *EventSystem) OnUserCreated(user User) {
 es.bus.Publicar(Evento{
  Tipo:  "user.created",
  Datos: user,
 })
}
```

---

## Ejercicios Progresivos

### Ejercicio 25.1: Pipeline Simple (3 Etapas)

**Objetivo**: Implementar un pipeline que:

- Etapa 1: Genera números del 1 al 10
- Etapa 2: Multiplica por 2
- Etapa 3: Suma 5
- Imprime el resultado final

**Requisitos:**

- Usar canales sin buffer
- Cada etapa en goroutine separada
- Cierre correcto de canales

**Solución esperada:**

```
1 -> 2 -> 7
2 -> 4 -> 9
...
10 -> 20 -> 25
```

### Ejercicio 25.2: Worker Pool

**Objetivo**: Crear un worker pool que:

- Procese 20 tareas
- Tenga 4 workers
- Registre inicio, duración y fin de cada tarea
- Implemente graceful shutdown

**Requisitos:**

- Usar sync.WaitGroup
- Buffer de tareas apropiado
- Recuperación de errores

**Salida esperada:**

```
Worker 1: Tarea 1 (100ms)
Worker 2: Tarea 2 (150ms)
...
Total: 20 tareas en ~500ms
```

### Ejercicio 25.3: Fan-Out/Fan-In

**Objetivo**: Sistema que:

- Distribuya trabajo a 5 procesadores
- Recolecte resultados
- Agregue estadísticas

**Requisitos:**

- Implementar fan-out sin buffer excesivo
- Usar fan-in genérico
- Reportar cantidad de errores

### Ejercicio 25.4: Rate Limiter

**Objetivo**: Limitar 10 solicitudes por segundo:

- Procesar 50 solicitudes
- Mostrar timestamps
- Verificar que respeta el límite

**Requisitos:**

- Usar token bucket
- Permitir ráfagas pequeñas
- Medir tiempo real

### Ejercicio 25.5: Pub/Sub con Múltiples Eventos

**Objetivo**: Sistema de eventos que:

- Tenga 3 tipos de eventos
- 4 suscriptores (cada uno a 1-2 tipos)
- Publique 10 eventos aleatorios
- Registre llegada a cada suscriptor

**Requisitos:**

- Sin deadlocks
- Manejo de desuscripción
- Validar que cada suscriptor recibe su tipo

---

**Fin del Capítulo 25**

Estos patrones son la base de sistemas Go escalables. La clave está en elegir el patrón correcto según:

1. Naturaleza de los datos (streaming, batch, evento)
2. Requisitos de latencia y throughput
3. Topología del sistema
4. Requisitos de recuperación

En producción, considera librerías como `gonum`, `akka`, o mensajería (RabbitMQ, Kafka) para problemas complejos.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/25-patrones-de-concurrencia/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/25-patrones-de-concurrencia):

```bash
cd examples/25-patrones-de-concurrencia
go run .
```
