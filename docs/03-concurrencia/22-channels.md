# Capítulo 22: Channels - Comunicación entre goroutines

## Índice
1. [¿Qué es un Channel?](#221-qué-es-un-channel)
2. [Crear y Usar Channels](#222-crear-y-usar-channels)
3. [Canales Buffered vs Unbuffered](#223-canales-buffered-vs-unbuffered)
4. [Cerrar Channels](#224-cerrar-channels)
5. [Range over Channels](#225-range-over-channels)
6. [Comma-OK Pattern](#226-comma-ok-pattern)
7. [Direccionalidad de Channels](#227-direccionalidad-de-channels)
8. [Select Statement Básico](#228-select-statement-básico)
9. [Deadlocks y Common Mistakes](#229-deadlocks-y-common-mistakes)
10. [Patrones: Productor-Consumidor](#2210-patrones-productor-consumidor)
11. [Buenas Prácticas y Antipatrones](#2211-buenas-prácticas-y-antipatrones)

---

## 22.1 ¿Qué es un Channel?

### 22.1.1 Concepto Fundamental

Un **channel** en Go es un primitivo de concurrencia que permite la **comunicación segura entre goroutines**. Es el mecanismo principal de sincronización en Go basado en el paradigma CSP (Communicating Sequential Processes).

**Definición CSP (Communicating Sequential Processes):**
- Concepto desarrollado por Tony Hoare en 1978
- Define sistemas concurrentes como procesos independientes que se comunican
- En Go, las goroutines son procesos y los channels son el mecanismo de comunicación
- Principio: "No compartir memoria para comunicarse; comunicarse para compartir memoria"

```
Sin Channels (Compartir Memoria):
┌──────────────┐
│  Goroutine1  │
│  (Mutex)     │──→ Variable Compartida ←── │  Goroutine2  │
│              │                            │  (Mutex)     │
└──────────────┘                            └──────────────┘
Problema: Race conditions, deadlocks

Con Channels (Comunicación):
┌──────────────┐                            ┌──────────────┐
│  Goroutine1  │──→ [ Channel ] ──→│  Goroutine2  │
│  (Sender)    │                            │  (Receiver)  │
└──────────────┘                            └──────────────┘
Ventaja: Comunicación segura y ordenada
```

### 22.1.2 Analogía del Mundo Real

Imagina dos personas comunicándose a través de un tubo:

- **Enviar datos**: Una persona coloca un mensaje en el tubo (bloqueante si está lleno)
- **Recibir datos**: La otra persona toma un mensaje del tubo (bloqueante si está vacío)
- **Capacidad**: El tubo puede tener límite de mensajes simultáneamente
- **Cierre**: Cuando termina, la persona que envia puede cerrar el tubo

### 22.1.3 Características Clave de Channels

**1. Tipados**
- Cada channel transmite un tipo específico de dato
- Seguridad de tipos en tiempo de compilación

**2. Sincronización Implícita**
- Enviar y recibir son operaciones que sincronización automáticamente
- No requieren explícitamente locks o mutexes

**3. Ordenamiento Garantizado**
- Los datos se reciben en el orden que se enviaron (FIFO)
- Comportamiento predecible

**4. Bloqueo**
- Sender se bloquea si el channel está lleno
- Receiver se bloquea si el channel está vacío
- Esto permite control de flujo automático

**5. Seguridad de Concurrencia**
- Múltiples goroutines pueden enviar/recibir simultáneamente
- Go garantiza que no hay data races

### 22.1.4 Pipes vs Channels

| Aspecto | Pipes (Unix) | Channels (Go) |
|--------|-------------|---------------|
| Tipo | Byte stream | Datos tipados |
| Procesos | Procesos del SO | Goroutines |
| Sincronización | Automática por buffer | Configurable |
| Seguridad de tipos | No (bytes) | Sí (tipos genéricos) |
| Cierre | Automático al terminar | Explícito con close() |

### 22.1.5 Comparación con Otros Lenguajes

| Lenguaje | Mecanismo | Filosofía |
|----------|-----------|-----------|
| Go | Channels | CSP (comunicar para sincronizar) |
| Java | BlockingQueue, wait/notify | Compartir memoria (sincronización explícita) |
| Rust | mpsc, async/await | Seguridad de tipos + ownership |
| Python | Queue.Queue | Thread-safe pero con GIL |
| C++ | std::condition_variable | Mutex + locks |

---

## 22.2 Crear y Usar Channels

### 22.2.1 Sintaxis Básica

**Declaración:**
```go
// Channel sin inicializar (nil)
var ch chan int

// Channel inicializado con make() - unbuffered (capacidad 0)
ch := make(chan int)

// Channel buffered (capacidad 5)
ch := make(chan int, 5)

// Channel de strings
ch := make(chan string)

// Channel de interfaces (acepta cualquier tipo)
ch := make(chan interface{})
```

### 22.2.2 Operaciones Básicas

**Enviar (Send):**
```go
ch <- valor  // Enviar valor al channel
```

**Recibir (Receive):**
```go
valor := <-ch  // Recibir valor del channel

// Recibir descartando el valor
<-ch
```

**Regla de Direccionalidad:**
- `<-ch` : flecha apuntando al channel = recibir
- `ch <-` : flecha apuntando lejos del channel = enviar

### 22.2.3 Ejemplo Completo: Comunicación Simple

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Crear channel unbuffered de integers
	ch := make(chan int)

	// Lanzar goroutine que envía un valor
	go func() {
		fmt.Println("Sender: preparando valor...")
		time.Sleep(1 * time.Second)
		
		fmt.Println("Sender: enviando 42 al channel...")
		ch <- 42  // Enviar valor (se bloquea hasta que alguien reciba)
		
		fmt.Println("Sender: valor enviado!")
	}()

	// Main goroutine recibe
	fmt.Println("Main: esperando valor del channel...")
	valor := <-ch  // Recibir valor (se bloquea hasta que haya algo que recibir)
	fmt.Println("Main: recibido valor:", valor)
}

// Output:
// Main: esperando valor del channel...
// Sender: preparando valor...
// Sender: enviando 42 al channel...
// Sender: valor enviado!
// Main: recibido valor: 42
```

### 22.2.4 Tipos de Datos en Channels

```go
package main

func main() {
	// Channel de strings
	strCh := make(chan string)
	go func() {
		strCh <- "Hola, Go!"
	}()
	msg := <-strCh
	println(msg)

	// Channel de structs
	type Usuario struct {
		Nombre string
		Edad   int
	}
	
	userCh := make(chan Usuario)
	go func() {
		userCh <- Usuario{"Alice", 30}
	}()
	user := <-userCh
	println(user.Nombre, user.Edad)

	// Channel de slices
	sliceCh := make(chan []int)
	go func() {
		sliceCh <- []int{1, 2, 3, 4, 5}
	}()
	datos := <-sliceCh
	println(len(datos))

	// Channel de functions
	funcCh := make(chan func(int) int)
	go func() {
		funcCh <- func(x int) int { return x * 2 }
	}()
	fn := <-funcCh
	println(fn(5))  // Output: 10
}
```

### 22.2.5 Múltiples Receptores y Remitentes

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	// Múltiples goroutines enviando
	for i := 1; i <= 3; i++ {
		go func(id int) {
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Goroutine %d enviando...\n", id)
			ch <- id * 10
		}(i)
	}

	// Recibir múltiples valores
	for i := 0; i < 3; i++ {
		valor := <-ch
		fmt.Printf("Main recibió: %d\n", valor)
	}
}

// Output (orden puede variar):
// Goroutine 1 enviando...
// Main recibió: 10
// Goroutine 2 enviando...
// Main recibió: 20
// Goroutine 3 enviando...
// Main recibió: 30
```

---

## 22.3 Canales Buffered vs Unbuffered

### 22.3.1 Canales Unbuffered

**Definición:**
- Capacidad = 0
- Enviar y recibir deben estar listos simultáneamente
- El sender se bloquea hasta que hay un receiver
- El receiver se bloquea hasta que hay un sender

```go
ch := make(chan int)  // Unbuffered (capacidad 0)
```

**Diagrama de Bloqueo:**

```
Unbuffered Channel: make(chan int)

Scenario 1: Sender primero
┌─────────────┐
│  Sender     │
│  ch <- 42   │ ← BLOQUEADO (esperando receiver)
└─────────────┘
            ↓ (espera a receiver)
        [Channel]
            ↑ (cuando receiver llega)
┌─────────────┐
│  Receiver   │
│  <-ch       │ ← Se desbloquea
└─────────────┘

Scenario 2: Receiver primero
┌─────────────┐
│  Receiver   │
│  <-ch       │ ← BLOQUEADO (esperando sender)
└─────────────┘
            ↓ (espera a sender)
        [Channel]
            ↑ (cuando sender llega)
┌─────────────┐
│  Sender     │
│  ch <- 42   │ ← Se desbloquea
└─────────────┘
```

### 22.3.2 Canales Buffered

**Definición:**
- Capacidad > 0
- Enviar no se bloquea mientras haya espacio
- Recibir no se bloquea mientras haya datos
- Permite desacoplamiento temporal entre sender y receiver

```go
ch := make(chan int, 5)  // Buffered (capacidad 5)
```

**Diagrama de Bloqueo:**

```
Buffered Channel: make(chan int, 5)

Caso 1: Buffer con espacio
ch <- 42        // No se bloquea, datos en buffer
ch <- 43        // No se bloquea, datos en buffer
valor := <-ch   // Recibe 42, buffer tiene 43

Caso 2: Buffer lleno
ch <- 1         // OK
ch <- 2         // OK
ch <- 3         // OK
ch <- 4         // OK
ch <- 5         // OK
ch <- 6         // ← BLOQUEADO (buffer lleno)

Caso 3: Buffer vacío
valor := <-ch   // Recibe un valor
valor := <-ch   // Recibe otro valor
valor := <-ch   // ← BLOQUEADO (buffer vacío)
```

### 22.3.3 Comparación: Unbuffered vs Buffered

| Aspecto | Unbuffered | Buffered |
|--------|-----------|----------|
| Capacidad | 0 | > 0 |
| Sincronización | Directa | Desacoplada |
| Bloqueo sender | Si no hay receiver | Si buffer lleno |
| Bloqueo receiver | Si no hay sender | Si buffer vacío |
| Caso de uso | Handoff, sincronización | Productor-consumidor |
| Overhead | Menor | Mayor (más memoria) |

### 22.3.4 Ejemplo: Buffered vs Unbuffered

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// ===== UNBUFFERED =====
	fmt.Println("=== Unbuffered ===")
	
	unbufferedCh := make(chan int)
	
	go func() {
		// Sin receiver esperando, sender se bloquea aquí
		fmt.Println("Sender: enviando 1...")
		unbufferedCh <- 1
		fmt.Println("Sender: 1 enviado, continuando...")
		
		fmt.Println("Sender: enviando 2...")
		unbufferedCh <- 2
		fmt.Println("Sender: 2 enviado, continuando...")
	}()
	
	time.Sleep(2 * time.Second)  // Receiver tarda
	
	fmt.Println("Main: recibiendo...")
	fmt.Println("Main: recibió", <-unbufferedCh)
	fmt.Println("Main: recibió", <-unbufferedCh)
	
	time.Sleep(1 * time.Second)
	fmt.Println()
	
	// ===== BUFFERED =====
	fmt.Println("=== Buffered (capacidad 2) ===")
	
	bufferedCh := make(chan int, 2)
	
	go func() {
		// Sender no se bloquea, data va al buffer
		fmt.Println("Sender: enviando 1...")
		bufferedCh <- 1
		fmt.Println("Sender: 1 enviado, continuando...")
		
		fmt.Println("Sender: enviando 2...")
		bufferedCh <- 2
		fmt.Println("Sender: 2 enviado, continuando...")
		
		fmt.Println("Sender: terminado (sin bloqueo)")
	}()
	
	time.Sleep(2 * time.Second)  // Receiver tarda
	
	fmt.Println("Main: recibiendo...")
	fmt.Println("Main: recibió", <-bufferedCh)
	fmt.Println("Main: recibió", <-bufferedCh)
}

// Output:
// === Unbuffered ===
// Sender: enviando 1...
// Main: recibiendo...
// Sender: 1 enviado, continuando...
// Sender: enviando 2...
// Main: recibió 1
// Sender: 2 enviado, continuando...
// Main: recibió 2
//
// === Buffered (capacidad 2) ===
// Sender: enviando 1...
// Sender: 1 enviado, continuando...
// Sender: enviando 2...
// Sender: 2 enviado, continuando...
// Sender: terminado (sin bloqueo)
// Main: recibiendo...
// Main: recibió 1
// Main: recibió 2
```

### 22.3.5 Eligiendo Buffer Size

```go
package main

func main() {
	// Regla: Make() con buffer size es "aspirational"
	
	// 1. Para control de flujo: tan pequeño como sea posible
	//    Channel buffered con 1: productor envía al siguiente en pipelining
	resultCh := make(chan Result, 1)
	
	// 2. Para worker pools: workers * promedio de trabajos en cola
	numWorkers := 4
	workCh := make(chan Job, numWorkers*10)  // 40 trabajos en espera máximo
	
	// 3. Para multiplicación (fan-out): número de receptores
	numConsumers := 5
	fanoutCh := make(chan Data, numConsumers)
	
	// 4. Unbuffered para sincronización estricta
	syncCh := make(chan struct{})  // Común para señales, no datos
	
	// 5. Lo que casi nunca debes hacer:
	// ❌ tooBigCh := make(chan int, 1000000)  // Desperdicio de memoria
}

type Result struct{}
type Job struct{}
type Data struct{}
```

---

## 22.4 Cerrar Channels

### 22.4.1 La Función close()

```go
// Sintaxis
close(ch)

// Efectos:
// - Sender ya no puede enviar (panic si intenta)
// - Receiver puede seguir recibiendo valores en buffer
// - Receiver detecta cierre con comma-ok pattern
```

### 22.4.2 Regla de los 3: Quién Puede Cerrar

**Regla fundamental:**
- Solo el **sender** debe cerrar el channel
- Si múltiples senders, ninguno debe cerrar (o coordinar)
- El receiver **NUNCA** cierra

**Razón:** Si un receiver cierra o múltiples senders cierran, puede causar panic.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// ===== CORRECTO: Un sender que cierra =====
	fmt.Println("=== Correcto: Un sender ===")
	
	ch := make(chan int)
	
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)  // ✓ OK - el único sender cierra
	}()
	
	for val := range ch {
		fmt.Println("Recibido:", val)
	}
	fmt.Println("Channel cerrado detectado")
	fmt.Println()
	
	// ===== INCORRECTO: Receiver intenta cerrar =====
	fmt.Println("=== Incorrecto: Receiver cierra ===")
	
	ch2 := make(chan int)
	
	go func() {
		for i := 1; i <= 2; i++ {
			ch2 <- i
		}
	}()
	
	val := <-ch2
	fmt.Println("Recibido:", val)
	// close(ch2)  // ✗ PANIC - receiver nunca debe cerrar
}
```

### 22.4.3 Múltiples Senders: Patrón Sync.WaitGroup

Cuando hay múltiples senders y quieres cerrar el channel:

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	const numSenders = 3
	const numValues = 2
	
	// Buffer para evitar goroutine leak
	resultCh := make(chan int, numSenders*numValues)
	
	var wg sync.WaitGroup
	
	// Lanzar múltiples senders
	for i := 1; i <= numSenders; i++ {
		wg.Add(1)
		
		go func(senderID int) {
			defer wg.Done()
			
			for j := 1; j <= numValues; j++ {
				value := senderID*10 + j
				fmt.Printf("Sender %d enviando %d\n", senderID, value)
				resultCh <- value
			}
		}(i)
	}
	
	// Goroutine que espera a que terminen todos los senders
	go func() {
		wg.Wait()
		close(resultCh)  // ✓ Solo después de que terminaron todos
	}()
	
	// Recibir todos los valores
	for val := range resultCh {
		fmt.Println("Recibido:", val)
	}
	
	fmt.Println("Todos los senders terminaron")
}
```

### 22.4.4 Patrón: Coordinator for Closing

Para cerrar múltiples channels, usa un coordinator:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	const numWorkers = 3
	
	// Channels
	done := make(chan struct{})      // Signal para parar
	results := make(chan string, 10) // Resultados
	
	// Workers que responden al signal
	for id := 1; id <= numWorkers; id++ {
		go func(workerID int) {
			for {
				select {
				case <-done:  // Señal para terminar
					fmt.Printf("Worker %d terminando\n", workerID)
					return
				default:
					// Simular trabajo
					fmt.Printf("Worker %d trabajando\n", workerID)
					results <- fmt.Sprintf("Resultado %d", workerID)
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(id)
	}
	
	// Recopilar algunos resultados
	for i := 0; i < 5; i++ {
		fmt.Println("Main:", <-results)
	}
	
	// Señal para terminar todos los workers
	close(done)
	
	time.Sleep(1 * time.Second)
	fmt.Println("Todos los workers terminaron")
}
```

### 22.4.5 Comportamiento después de close()

```go
package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	
	fmt.Println("=== Recibir después de close ===")
	
	// Puedes seguir recibiendo valores del buffer
	fmt.Println("Valor 1:", <-ch)
	fmt.Println("Valor 2:", <-ch)
	
	// Después que se agota el buffer, recibe zero value
	fmt.Println("Valor 3:", <-ch)  // 0 (zero value)
	fmt.Println("Valor 4:", <-ch)  // 0 (zero value)
	
	fmt.Println()
	fmt.Println("=== Comma-ok pattern ===")
	
	ch2 := make(chan int, 1)
	ch2 <- 42
	close(ch2)
	
	if val, ok := <-ch2; ok {
		fmt.Println("Recibido:", val)
	} else {
		fmt.Println("Channel cerrado")
	}
	
	if val, ok := <-ch2; ok {
		fmt.Println("Recibido:", val)
	} else {
		fmt.Println("Channel cerrado")  // Imprime esto
	}
	
	fmt.Println()
	fmt.Println("=== Intentar enviar a channel cerrado ===")
	
	ch3 := make(chan int)
	close(ch3)
	
	// ❌ Esto causa panic:
	// ch3 <- 42  // panic: send on closed channel
}
```

---

## 22.5 Range over Channels

### 22.5.1 Iteración Básica con for-range

```go
// Sintaxis
for value := range ch {
	// value recibe cada elemento del channel
	// Loop termina cuando channel se cierra
}
```

**Cómo funciona:**
1. Bloquea esperando el siguiente valor del channel
2. Si hay un valor, lo asigna a `value`
3. Si el channel está cerrado, termina el loop
4. Repite hasta que channel se cierre

### 22.5.2 Ejemplo: Productor-Consumidor con Range

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Channel para números
	numCh := make(chan int)
	
	// Productor: genera números
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("Productor: generando %d\n", i)
			numCh <- i
			time.Sleep(300 * time.Millisecond)
		}
		close(numCh)  // Productor cierra cuando termina
	}()
	
	// Consumidor: procesa números con range
	fmt.Println("Consumidor: esperando números...")
	for num := range numCh {
		fmt.Printf("Consumidor: procesa %d\n", num)
	}
	
	fmt.Println("Consumidor: canal cerrado, terminando")
}

// Output:
// Consumidor: esperando números...
// Productor: generando 1
// Consumidor: procesa 1
// Productor: generando 2
// Consumidor: procesa 2
// ...
// Consumidor: canal cerrado, terminando
```

### 22.5.3 Detalles Importantes

**1. Range no funciona si el channel no se cierra:**

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	
	go func() {
		// Enviar solo 3 valores sin cerrar
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		// ❌ Sin close(ch), range seguirá esperando
	}()
	
	// Esto se BLOQUEA esperando más valores
	for num := range ch {
		fmt.Println("Recibido:", num)
	}
	
	fmt.Println("Esto nunca se imprime")
	
	// ⏰ Programa se queda esperando infinitamente
}
```

**2. Range recibe solo el valor, no el ok:**

```go
package main

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	
	// Con range solo obtienes el valor
	for val := range ch {
		println(val)  // 1, luego 2
	}
	
	// Si quieres el ok, usa receive explícito
	ch2 := make(chan int, 1)
	ch2 <- 42
	close(ch2)
	
	if val, ok := <-ch2; ok {
		println("Valor:", val, "OK:", ok)
	}
}
```

### 22.5.4 Patrones Comunes

**Patrón 1: Pipeline con Range**

```go
package main

import "fmt"

func main() {
	// Stage 1: Generar números
	numbers := make(chan int)
	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers)
	}()
	
	// Stage 2: Cuadrados (consume numbers, produce squares)
	squares := make(chan int)
	go func() {
		for num := range numbers {  // Automático: close(numbers)
			squares <- num * num
		}
		close(squares)
	}()
	
	// Stage 3: Consumir cuadrados
	for sq := range squares {
		fmt.Println(sq)  // 1, 4, 9, 16, 25
	}
}
```

**Patrón 2: Múltiples Consumidores**

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	tasks := make(chan int)
	
	// Lanzar múltiples workers
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			// Cada worker procesa tareas del canal
			for task := range tasks {
				fmt.Printf("Worker %d procesa tarea %d\n", workerID, task)
			}
			fmt.Printf("Worker %d terminó\n", workerID)
		}(w)
	}
	
	// Productor envia tareas
	for i := 1; i <= 10; i++ {
		tasks <- i
	}
	close(tasks)  // Workers se despiertan y terminan
	
	wg.Wait()
	fmt.Println("Todos los workers terminaron")
}
```

**Patrón 3: Range con índice (no típico pero posible)**

```go
package main

import (
	"fmt"
)

func main() {
	// Convertir channel a slice es una forma de contar
	values := make(chan int, 3)
	values <- 10
	values <- 20
	values <- 30
	close(values)
	
	// Iterar con índice manual
	index := 0
	for val := range values {
		fmt.Printf("[%d] = %d\n", index, val)
		index++
	}
	
	// O simplemente:
	index2 := 0
	for val := range values {
		_ = val  // Ya cerrado, no funciona dos veces
		index2++
	}
}
```

---

## 22.6 Comma-OK Pattern

### 22.6.1 Sintaxis y Significado

```go
// Recepciones siempre devuelven (value, ok)
value, ok := <-ch

// ok es true si valor es válido
// ok es false si channel está cerrado y no hay más datos
```

**Comparación con map:**

```go
// Map: m[key] devuelve (value, ok)
val, ok := m["key"]  // ok es false si key no existe

// Channel: <-ch devuelve (value, ok)
val, ok := <-ch      // ok es false si channel cerrado
```

### 22.6.2 Ejemplo: Detectar Cierre de Channel

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "primer mensaje"
		
		time.Sleep(500 * time.Millisecond)
		ch <- "segundo mensaje"
		
		time.Sleep(500 * time.Millisecond)
		close(ch)  // Cerrar después de enviar dos mensajes
	}()
	
	// Recibir con comma-ok para detectar cierre
	for {
		msg, ok := <-ch
		
		if !ok {
			fmt.Println("Channel cerrado, saliendo")
			break
		}
		
		fmt.Println("Recibido:", msg)
	}
}

// Output:
// Recibido: primer mensaje
// Recibido: segundo mensaje
// Channel cerrado, saliendo
```

### 22.6.3 Comparación: comma-ok vs range

```go
package main

import "fmt"

func demonstrateChannelReception() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	
	fmt.Println("=== Con range (automático) ===")
	// range automáticamente detecta cierre
	for val := range ch {
		fmt.Println(val)
	}
	
	fmt.Println()
	fmt.Println("=== Con comma-ok (manual) ===")
	
	ch2 := make(chan int, 2)
	ch2 <- 1
	ch2 <- 2
	close(ch2)
	
	// Comma-ok requiere loop manual
	for {
		val, ok := <-ch2
		if !ok {
			fmt.Println("Channel cerrado")
			break
		}
		fmt.Println(val)
	}
}

func main() {
	demonstrateChannelReception()
}
```

### 22.6.4 Casos de Uso del Comma-OK Pattern

**Caso 1: Timeout con select**

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "dato"
	}()
	
	// Esperar hasta 1 segundo
	select {
	case msg, ok := <-ch:
		if ok {
			fmt.Println("Recibido:", msg)
		} else {
			fmt.Println("Channel cerrado")
		}
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: no hay dato")
	}
}

// Output:
// Timeout: no hay dato
```

**Caso 2: Procesamiento condicional**

```go
package main

import "fmt"

func processMessages(ch chan string) {
	for {
		msg, ok := <-ch
		
		// ok es false cuando channel está cerrado
		if !ok {
			fmt.Println("Channel cerrado, procesamiento terminado")
			return
		}
		
		// Procesar mensaje
		if msg == "exit" {
			fmt.Println("Señal de salida recibida")
			return
		}
		
		fmt.Printf("Procesando: %s\n", msg)
	}
}

func main() {
	ch := make(chan string)
	
	go processMessages(ch)
	
	ch <- "tarea1"
	ch <- "tarea2"
	ch <- "exit"
	
	close(ch)
}
```

**Caso 3: Verificar si hay dato sin bloquear**

```go
package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	ch <- 42
	
	// Intento de recibir sin bloquear (select con default)
	select {
	case val, ok := <-ch:
		if ok {
			fmt.Println("Valor disponible:", val)
		}
	default:
		fmt.Println("No hay valor disponible")
	}
	
	// Después de consumir
	<-ch
	
	select {
	case val, ok := <-ch:
		if ok {
			fmt.Println("Valor disponible:", val)
		}
	default:
		fmt.Println("No hay valor disponible")  // Imprime esto
	}
}
```

---

## 22.7 Direccionalidad de Channels

### 22.7.1 Tipos de Channels

**1. Bidireccionales (por defecto)**

```go
// Puede enviar y recibir
ch := make(chan int)

// Funciones que aceptan bidireccionales
func process(ch chan int) {
	val := <-ch  // Recibir
	ch <- 42     // Enviar
}
```

**2. Send-only (solo envío)**

```go
// Solo para enviar
var sendCh chan<- int

sendCh = make(chan int)

// Operaciones válidas
sendCh <- 42

// Operaciones inválidas
// val := <-sendCh  // ✗ Compile error: receive from send-only channel
// close(sendCh)    // ✗ Compile error
```

**3. Receive-only (solo recepción)**

```go
// Solo para recibir
var recvCh <-chan int

recvCh = make(chan int)

// Operaciones válidas
val := <-recvCh

// Operaciones inválidas
// recvCh <- 42     // ✗ Compile error: send to receive-only channel
// close(recvCh)    // ✗ Compile error
```

### 22.7.2 Conversión Automática

```go
package main

import "fmt"

func sender(out chan<- int) {
	// Solo puede enviar
	for i := 1; i <= 3; i++ {
		out <- i
	}
	close(out)
}

func receiver(in <-chan int) {
	// Solo puede recibir
	for val := range in {
		fmt.Println("Recibido:", val)
	}
}

func main() {
	// Channel bidireccional
	ch := make(chan int)
	
	// Se convierte automáticamente a send-only
	go sender(ch)
	
	// Se convierte automáticamente a receive-only
	receiver(ch)
}

// Output:
// Recibido: 1
// Recibido: 2
// Recibido: 3
```

### 22.7.3 Beneficios de la Direccionalidad

**1. Seguridad de Tipos**

```go
package main

import "fmt"

// API clara: solo envío
func generateNumbers(out chan<- int) {
	for i := 1; i <= 3; i++ {
		out <- i
	}
	close(out)
}

// API clara: solo recepción
func sumNumbers(in <-chan int) int {
	sum := 0
	for val := range in {
		sum += val
	}
	return sum
}

func main() {
	ch := make(chan int)
	
	go generateNumbers(ch)
	total := sumNumbers(ch)
	
	fmt.Println("Suma:", total)  // 6
}
```

**2. Prevención de Bugs**

```go
// ✗ SIN direccionalidad (errores en runtime)
func badFunction(ch chan int) {
	// ¿Qué se espera? ¿Enviar o recibir?
	// Ambos son posibles, confusión
}

// ✓ CON direccionalidad (errores en compilación)
func goodSender(out chan<- int) {
	// ✗ val := <-out  // Compile error - previene bug
}

func goodReceiver(in <-chan int) {
	// ✗ in <- 42      // Compile error - previene bug
}
```

### 22.7.4 Patrón: Pipeline Tipado

```go
package main

import (
	"fmt"
	"time"
)

// Pipeline stage 1: Generar números
func numbers(out chan<- int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

// Pipeline stage 2: Multiplicar por 2
func double(in <-chan int, out chan<- int) {
	for num := range in {
		out <- num * 2
	}
	close(out)
}

// Pipeline stage 3: Sumar 10
func addTen(in <-chan int, out chan<- int) {
	for num := range in {
		out <- num + 10
	}
	close(out)
}

// Pipeline stage 4: Consumir
func display(in <-chan int) {
	for result := range in {
		fmt.Println(result)
	}
}

func main() {
	// Conectar pipeline
	nums := make(chan int)
	doubled := make(chan int)
	added := make(chan int)
	
	go numbers(nums)
	go double(nums, doubled)
	go addTen(doubled, added)
	
	display(added)  // Bloqueante, espera que se cierre
}

// Output:
// 12
// 14
// 16
// 18
// 20
```

---

## 22.8 Select Statement Básico

### 22.8.1 Sintaxis y Concepto

```go
// Select es como switch pero para channel operations
select {
case <-ch1:
	// ch1 tiene dato o cierra
case val := <-ch2:
	// ch2 envía val
case ch3 <- dato:
	// ch3 acepta dato
case <-time.After(timeout):
	// Timeout
default:
	// Ninguno de arriba está listo (no-blocking)
}
```

**Características:**
- Bloquea hasta que UNA operación esté lista
- Elige aleatoriamente si múltiples están listas (fairness)
- Default es NO-BLOCKING (ejecuta siempre si otros no están listos)
- Útil para multiplexing de múltiples channels

### 22.8.2 Ejemplo: Multiplexing Básico

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Dos channels con datos llegando en tiempos diferentes
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	// Goroutine 1: envía cada 1 segundo
	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(1 * time.Second)
			ch1 <- fmt.Sprintf("ch1: mensaje %d", i)
		}
		close(ch1)
	}()
	
	// Goroutine 2: envía cada 1.5 segundos
	go func() {
		for i := 1; i <= 2; i++ {
			time.Sleep(1500 * time.Millisecond)
			ch2 <- fmt.Sprintf("ch2: mensaje %d", i)
		}
		close(ch2)
	}()
	
	// Main: multiplexea ambos channels
	active := 2  // Cuántos channels activos
	
	for active > 0 {
		select {
		case msg, ok := <-ch1:
			if !ok {
				fmt.Println("ch1 cerrado")
				active--
				ch1 = nil  // Poner a nil para evitar panic en futuros select
			} else {
				fmt.Println("Recibido de ch1:", msg)
			}
			
		case msg, ok := <-ch2:
			if !ok {
				fmt.Println("ch2 cerrado")
				active--
				ch2 = nil
			} else {
				fmt.Println("Recibido de ch2:", msg)
			}
		}
	}
}

// Output (aprox):
// Recibido de ch1: ch1: mensaje 1
// Recibido de ch2: ch2: mensaje 1
// Recibido de ch1: ch1: mensaje 2
// Recibido de ch2: ch2: mensaje 2
// Recibido de ch1: ch1: mensaje 3
// ch1 cerrado
// ch2 cerrado
```

### 22.8.3 Select con Default (Non-blocking)

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	
	// Select con default: no se bloquea
	select {
	case val := <-ch:
		fmt.Println("Recibido:", val)
	default:
		fmt.Println("No hay dato disponible (non-blocking)")
	}
	
	// Enviar con non-blocking
	select {
	case ch <- 42:
		fmt.Println("Dato enviado")
	default:
		fmt.Println("No se puede enviar (channel bloqueado)")  // Imprime esto
	}
	
	// Con buffer, non-blocking send funciona
	ch2 := make(chan int, 1)
	
	select {
	case ch2 <- 42:
		fmt.Println("Dato enviado a ch2")
	default:
		fmt.Println("No se puede enviar a ch2")
	}
}

// Output:
// No hay dato disponible (non-blocking)
// No se puede enviar (channel bloqueado)
// Dato enviado a ch2
```

### 22.8.4 Select con Timeout

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	
	// Simular que el dato llega en 2 segundos
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "dato importante"
	}()
	
	// Esperar máximo 1 segundo
	select {
	case msg := <-ch:
		fmt.Println("Recibido:", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: no hay dato en 1 segundo")
	}
	
	fmt.Println()
	fmt.Println("=== Esperando suficiente tiempo ===")
	
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch2 <- "dato"
	}()
	
	// Esperar máximo 2 segundos
	select {
	case msg := <-ch2:
		fmt.Println("Recibido:", msg)
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout")
	}
}

// Output:
// Timeout: no hay dato en 1 segundo
//
// === Esperando suficiente tiempo ===
// Recibido: dato
```

### 22.8.5 Patrón: Fan-in (Multiplexing)

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Worker simula goroutines que hacen trabajo
func worker(id int, out chan<- string) {
	for {
		delay := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(delay)
		out <- fmt.Sprintf("Worker %d: completó tarea", id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	// Múltiples workers enviando al mismo channel
	work1 := make(chan string)
	work2 := make(chan string)
	work3 := make(chan string)
	
	go worker(1, work1)
	go worker(2, work2)
	go worker(3, work3)
	
	// Select multiplexea todos
	for i := 0; i < 6; i++ {
		select {
		case msg := <-work1:
			fmt.Println(msg)
		case msg := <-work2:
			fmt.Println(msg)
		case msg := <-work3:
			fmt.Println(msg)
		}
	}
}

// Output (orden variante):
// Worker 2: completó tarea
// Worker 1: completó tarea
// Worker 3: completó tarea
// ...
```

---

## 22.9 Deadlocks y Common Mistakes

### 22.9.1 Tipos de Deadlock

**Deadlock 1: Envío a Channel sin Receptor**

```go
package main

func main() {
	// ❌ DEADLOCK
	ch := make(chan int)
	
	// Intenta enviar, pero no hay quien reciba
	ch <- 42  // Se bloquea indefinidamente
	
	val := <-ch  // Nunca se alcanza
	println(val)
}

// Error: fatal error: all goroutines are asleep - deadlock!
```

**Solución:**

```go
package main

func main() {
	// ✓ CORRECTO
	ch := make(chan int)
	
	// Usar goroutine para evitar bloqueo
	go func() {
		ch <- 42
	}()
	
	val := <-ch
	println(val)  // 42
}
```

### 22.9.2 Deadlock 2: Recepción sin Remitente

```go
package main

func main() {
	// ❌ DEADLOCK
	ch := make(chan int)
	
	// Intenta recibir, pero no hay quien envíe
	val := <-ch  // Se bloquea indefinidamente
	println(val)
}

// Error: fatal error: all goroutines are asleep - deadlock!
```

**Solución:**

```go
package main

func main() {
	// ✓ CORRECTO
	ch := make(chan int)
	
	go func() {
		ch <- 42
	}()
	
	val := <-ch
	println(val)
}
```

### 22.9.3 Deadlock 3: Circular Dependencies

```go
package main

func main() {
	// ❌ DEADLOCK
	ch1 := make(chan int)
	ch2 := make(chan int)
	
	// Goroutine 1: espera a ch1, envia a ch2
	go func() {
		val := <-ch1
		ch2 <- val
	}()
	
	// Goroutine 2: espera a ch2, envia a ch1
	go func() {
		val := <-ch2
		ch1 <- val
	}()
	
	// Main espera indefinidamente
}

// Error: fatal error: all goroutines are asleep - deadlock!
```

**Solución: Inicializar una cadena**

```go
package main

import "fmt"

func main() {
	// ✓ CORRECTO
	ch1 := make(chan int)
	ch2 := make(chan int)
	
	go func() {
		val := <-ch1
		fmt.Println("Goroutine1 recibió:", val)
		ch2 <- val + 10
	}()
	
	go func() {
		val := <-ch2
		fmt.Println("Goroutine2 recibió:", val)
		ch1 <- val  // Esta línea nunca se alcanza (ya está en la cadena)
	}()
	
	// Iniciar la cadena
	ch1 <- 1
	
	// Esperar a que terminen
	fmt.Println("Goroutine1 envió:", <-ch2)
}
```

### 22.9.4 Deadlock 4: Enviar a Channel Cerrado

```go
package main

func main() {
	// ❌ PANIC
	ch := make(chan int)
	close(ch)
	
	// Enviar a canal cerrado causa panic
	ch <- 42  // panic: send on closed channel
}
```

**Solución:**

```go
package main

import "fmt"

func main() {
	// ✓ CORRECTO
	ch := make(chan int)
	
	go func() {
		ch <- 42
		close(ch)  // Sender cierra
	}()
	
	val := <-ch
	fmt.Println(val)
	
	// Verificar cierre
	if val, ok := <-ch; !ok {
		fmt.Println("Channel cerrado")
	}
}
```

### 22.9.5 Deadlock 5: Múltiples Goroutines esperando

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// ❌ DEADLOCK
	ch := make(chan int)
	
	var wg sync.WaitGroup
	
	// 3 goroutines esperando
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Todas esperan un valor
			val := <-ch
			fmt.Println("Recibido:", val)
		}()
	}
	
	// Enviar solo 1 valor
	ch <- 42
	
	wg.Wait()  // DEADLOCK: dos goroutines siguen esperando
}

// Error: fatal error: all goroutines are asleep - deadlock!
```

**Solución:**

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// ✓ CORRECTO: Usar broadcast (enviar a todos)
	ch := make(chan int, 3)  // Buffer para 3 valores
	
	var wg sync.WaitGroup
	
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := <-ch
			fmt.Println("Recibido:", val)
		}()
	}
	
	// Enviar un valor a cada goroutine
	for i := 0; i < 3; i++ {
		ch <- i
	}
	
	wg.Wait()
}
```

### 22.9.6 Debugging Deadlocks

```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// Técnica 1: Timeout global
	done := make(chan bool)
	
	go func() {
		// Código que podría deadlock
		ch := make(chan int)
		val := <-ch  // Potencial deadlock
		fmt.Println(val)
		
		done <- true
	}()
	
	select {
	case <-done:
		fmt.Println("Completado exitosamente")
	case <-time.After(2 * time.Second):
		fmt.Println("DEADLOCK DETECTADO después de 2 segundos")
		// Imprimir stack traces
		buf := make([]byte, 1<<20)
		stackSize := runtime.Stack(buf, true)
		fmt.Println("\nStack traces:")
		fmt.Println(string(buf[:stackSize]))
	}
}
```

---

## 22.10 Patrones: Productor-Consumidor

### 22.10.1 Patrón Clásico Productor-Consumidor

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Channel entre productor y consumidor
	items := make(chan int, 5)  // Buffered para desacoplamiento
	
	// Productor: genera datos
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("Productor: generando %d\n", i)
			items <- i * 10
			time.Sleep(100 * time.Millisecond)
		}
		close(items)  // Señal de fin
	}()
	
	// Consumidor: procesa datos
	fmt.Println("Consumidor: esperando items...")
	total := 0
	for item := range items {
		fmt.Printf("Consumidor: procesa %d\n", item)
		total += item
		time.Sleep(200 * time.Millisecond)  // Procesamiento lento
	}
	
	fmt.Printf("Consumidor: completado. Total = %d\n", total)
}
```

### 22.10.2 Patrón: Pipeline de Etapas

```go
package main

import (
	"fmt"
	"sync"
)

// Etapa 1: Generar números
func generate(out chan<- int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

// Etapa 2: Procesar (transformar)
func square(in <-chan int, out chan<- int) {
	for num := range in {
		out <- num * num
	}
	close(out)
}

// Etapa 3: Procesar (filtrar)
func filterEven(in <-chan int, out chan<- int) {
	for num := range in {
		if num%2 == 0 {
			out <- num
		}
	}
	close(out)
}

// Etapa 4: Consumir
func display(in <-chan int) {
	fmt.Println("Resultados finales:")
	for num := range in {
		fmt.Printf("  %d\n", num)
	}
}

func main() {
	// Conectar pipeline
	nums := make(chan int)
	squares := make(chan int)
	evens := make(chan int)
	
	go generate(nums)              // Etapa 1
	go square(nums, squares)       // Etapa 2
	go filterEven(squares, evens)  // Etapa 3
	display(evens)                 // Etapa 4 (bloqueante)
}

// Output:
// Resultados finales:
//   4
//   16
```

### 22.10.3 Patrón: Fan-Out/Fan-In

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Fan-out: distribuir trabajo a múltiples workers
func distribute(jobs <-chan int, workers int) <-chan int {
	results := make(chan int)
	
	// Lanzar workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			for job := range jobs {
				fmt.Printf("Worker %d procesa job %d\n", workerID, job)
				// Simular trabajo
				time.Sleep(500 * time.Millisecond)
				
				// Enviar resultado
				results <- job * 2
			}
		}(w)
	}
	
	// Cuando terminan todos los workers, cerrar results
	go func() {
		wg.Wait()
		close(results)
	}()
	
	return results
}

func main() {
	// Generar jobs
	jobs := make(chan int, 10)
	
	go func() {
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()
	
	// Fan-out a 3 workers, recopilar resultados (fan-in)
	results := distribute(jobs, 3)
	
	// Mostrar resultados
	fmt.Println("Resultados:")
	total := 0
	for result := range results {
		fmt.Printf("  Resultado: %d\n", result)
		total += result
	}
	fmt.Printf("Total: %d\n", total)
}
```

### 22.10.4 Patrón: Generador (Iterator Pattern)

```go
package main

import (
	"fmt"
)

// Generador que produce números en sequence
func fibonacci(n int) <-chan int {
	out := make(chan int)
	
	go func() {
		a, b := 0, 1
		for i := 0; i < n; i++ {
			out <- a
			a, b = b, a+b
		}
		close(out)
	}()
	
	return out
}

func main() {
	// Usar generador con range
	fmt.Println("Fibonacci (primeros 10):")
	for num := range fibonacci(10) {
		fmt.Printf("%d ", num)
	}
	fmt.Println()
	
	// Uso tipo filter
	fmt.Println("\nFilbonacci par (primeros 6):")
	count := 0
	for num := range fibonacci(20) {
		if num%2 == 0 && count < 6 {
			fmt.Printf("%d ", num)
			count++
		}
	}
	fmt.Println()
}

// Output:
// Fibonacci (primeros 10):
// 0 1 1 2 3 5 8 13 21 34
//
// Fibonacci par (primeros 6):
// 0 2 8 34 144 610
```

### 22.10.5 Patrón: Merge de Múltiples Channels

```go
package main

import (
	"fmt"
	"sync"
)

// Merge combina múltiples channels en uno
func merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	
	// Función que copia valores de un channel al output
	output := func(ch <-chan int) {
		defer wg.Done()
		for val := range ch {
			out <- val
		}
	}
	
	// Lanzar goroutine para cada channel
	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}
	
	// Cerrar output cuando todas las goroutines terminen
	go func() {
		wg.Wait()
		close(out)
	}()
	
	return out
}

func main() {
	// Crear múltiples productores
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	
	go func() {
		ch1 <- 1
		ch1 <- 4
		close(ch1)
	}()
	
	go func() {
		ch2 <- 2
		ch2 <- 5
		close(ch2)
	}()
	
	go func() {
		ch3 <- 3
		ch3 <- 6
		close(ch3)
	}()
	
	// Merge en uno
	results := merge(ch1, ch2, ch3)
	
	// Consumir
	fmt.Println("Merged results:")
	for val := range results {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}

// Output:
// Merged results:
// 1 4 2 5 3 6
// (orden puede variar)
```

---

## 22.11 Buenas Prácticas y Antipatrones

### 22.11.1 Reglas de Ownership (Quién Owna el Channel)

**Regla 1: Quien Crea Owna**

```go
package main

// ✓ CORRECTO: Creador owna y cierra
func goodProducerConsumer() <-chan int {
	// Esta función owna el channel
	out := make(chan int)
	
	go func() {
		for i := 1; i <= 3; i++ {
			out <- i
		}
		close(out)  // Owner cierra
	}()
	
	return out
}

func main() {
	// Receiver no puede cerrar
	for val := range goodProducerConsumer() {
		println(val)
	}
}
```

**Regla 2: Solo Sender Cierra**

```go
package main

import (
	"fmt"
	"sync"
)

// ❌ INCORRECTO: Múltiples senders tratando de cerrar
func badMultipleSenders() {
	ch := make(chan int)
	
	var wg sync.WaitGroup
	for s := 0; s < 3; s++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- s
			// close(ch)  // ✗ PANIC si múltiples cierran
		}()
	}
	
	go func() {
		wg.Wait()
		close(ch)  // ✓ Solo un closer
	}()
	
	for val := range ch {
		fmt.Println(val)
	}
}

func main() {
	badMultipleSenders()
}
```

### 22.11.2 Tamaño Correcto de Buffer

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// ❌ ANTIPATRON: Buffer demasiado grande
	// tooBig := make(chan Task, 1000000)  // Desperdicio de memoria
	
	// ✓ BUENO: Buffer pequeño o calculado
	
	// Caso 1: Sin buffer (sincronización estricta)
	syncCh := make(chan struct{})  // 0 buffer
	
	// Caso 2: Buffer de 1 (patrón común en pipelining)
	pipelineCh := make(chan Result, 1)
	
	// Caso 3: Buffer para N workers
	numWorkers := 4
	workCh := make(chan Job, numWorkers)  // Suficiente para 1 job per worker
	
	// Caso 4: Timeout receiver
	fmt.Println("=== Buffer Sizing ===")
	
	// Productor rápido, consumidor lento
	results := make(chan int, 1)  // Pequeño buffer ayuda
	
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("Produce %d\n", i)
			results <- i
		}
		close(results)
	}()
	
	for val := range results {
		fmt.Printf("Consume %d\n", val)
		time.Sleep(100 * time.Millisecond)  // Lento
	}
}

type Task struct{}
type Result struct{}
type Job struct{}
```

### 22.11.3 Patrón: Closing Signal

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// ✓ BUENO: Usar channel de cierre para señal
	
	// Channel para work
	work := make(chan string)
	
	// Channel para señal de cierre
	done := make(chan struct{})
	
	// Workers
	var wg sync.WaitGroup
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		
		go func(id int) {
			defer wg.Done()
			
			for {
				select {
				case task, ok := <-work:
					if !ok {
						fmt.Printf("Worker %d: work channel cerrado\n", id)
						return
					}
					fmt.Printf("Worker %d procesa: %s\n", id, task)
				
				case <-done:
					fmt.Printf("Worker %d: recibió señal done\n", id)
					return
				}
			}
		}(w)
	}
	
	// Enviar trabajo
	for i := 1; i <= 3; i++ {
		work <- fmt.Sprintf("tarea-%d", i)
	}
	
	// Después de cierto tiempo, enviar señal de parada
	time.Sleep(1 * time.Second)
	
	close(done)  // Señal para que los workers se detengan
	
	wg.Wait()
	fmt.Println("Todos los workers terminaron")
}
```

### 22.11.4 Antipatrones Comunes

**Antipatrón 1: Compartir memoria para comunicarse**

```go
// ❌ ANTIPATRON: Usar Mutex en lugar de Channels
type SharedData struct {
	mu    sync.Mutex
	value int
}

// ✓ CORRECTO: Usar Channels para comunicación
// (El canal maneja la sincronización automáticamente)
func goodCommunication(out chan<- int) {
	out <- 42  // Seguro por diseño
}
```

**Antipatrón 2: Receiver que cierra**

```go
// ❌ ANTIPATRON
func badReceiverClose(ch <-chan int) {
	for val := range ch {
		println(val)
	}
	// close(ch)  // ✗ PANIC - Receiver no puede cerrar
}

// ✓ CORRECTO
func goodReceiverDone(ch <-chan int) {
	for val := range ch {
		// Channel puede estar cerrado, range detecto
		println(val)
	}
}
```

**Antipatrón 3: Omitir close()**

```go
// ❌ ANTIPATRON: No cerrar channel (goroutine leak)
func badNoClose() <-chan int {
	out := make(chan int)
	
	go func() {
		out <- 1
		out <- 2
		// ❌ Sin close(out), range nunca termina
	}()
	
	return out
}

// ✓ CORRECTO
func goodWithClose() <-chan int {
	out := make(chan int)
	
	go func() {
		defer close(out)  // Siempre cerrar
		out <- 1
		out <- 2
	}()
	
	return out
}
```

**Antipatrón 4: Selecting el mismo channel múltiples veces**

```go
// ❌ ANTIPATRON: Redundante
func badSelectRedundant(ch <-chan int) {
	select {
	case val := <-ch:
		println(val)
	case val := <-ch:  // ✗ Redundante, nunca se elige
		println(val)
	}
}

// ✓ CORRECTO: Si necesitas múltiples, usa range
func goodMultipleValues(ch <-chan int) {
	for val := range ch {
		println(val)  // Procesa todos
	}
}
```

### 22.11.5 Checklist de Buenas Prácticas

```go
/*
CHECKLIST: Canales en Go

1. CREATING (Creación)
   ✓ [ ] ¿Es channel unbuffered (0) o buffered (>0)?
   ✓ [ ] ¿El buffer size está justificado?
   ✓ [ ] ¿Es channel bidireccional, send-only, o receive-only?

2. SENDING (Envío)
   ✓ [ ] ¿Está en goroutine separada para evitar bloqueo?
   ✓ [ ] ¿Se puede recibir del otro lado del channel?
   ✓ [ ] ¿El channel no está cerrado?

3. RECEIVING (Recepción)
   ✓ [ ] ¿Usa range cuando sea posible (automático close detection)?
   ✓ [ ] ¿Si no usa range, verifica ok := <-ch?
   ✓ [ ] ¿Hay un timeout si puede esperar indefinidamente?

4. CLOSING (Cierre)
   ✓ [ ] ¿Solo el sender (creador) cierra?
   ✓ [ ] ¿Si múltiples senders, usa WaitGroup + coordinator?
   ✓ [ ] ¿Receiver nunca intenta cerrar?

5. SYNCHRONIZATION (Sincronización)
   ✓ [ ] ¿Hay deadlock potencial (sender/receiver esperando)?
   ✓ [ ] ¿Usa select con timeout si es necesario?
   ✓ [ ] ¿Evita múltiples goroutines leyendo del mismo channel?

6. CLEANUP (Limpieza)
   ✓ [ ] ¿Todos los goroutines de channels terminan?
   ✓ [ ] ¿No hay goroutine leaks (esperando channel cerrado)?
   ✓ [ ] Se usa WaitGroup donde es necesario?
*/

func checklistExample() {
	// Buen ejemplo que sigue el checklist
	
	numbers := make(chan int, 5)  // (1) Buffered, tamaño justificado
	
	go func() {  // (2) En goroutine separada
		for i := 1; i <= 10; i++ {
			numbers <- i
		}
		close(numbers)  // (4) Sender cierra
	}()
	
	// (3) Usar range
	sum := 0
	for num := range numbers {
		sum += num
	}
	
	println("Sum:", sum)  // (6) Se completó limpiamente
}
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: Hello Channel - Comunicación Simple

**Objetivo:** Entender básicos de channels: crear, enviar, recibir.

```go
package main

import (
	"fmt"
)

func main() {
	// TODO: Completar el programa
	
	// 1. Crear un channel de strings
	// greetingCh := make(chan string)
	
	// 2. Lanzar una goroutine que envía un saludo
	// go func() {
	//     greetingCh <- "Hola desde goroutine!"
	// }()
	
	// 3. Recibir el saludo en main
	// mensaje := <-greetingCh
	
	// 4. Imprimir el mensaje
	// fmt.Println(mensaje)
}

// Salida esperada:
// Hola desde goroutine!
```

**Solución:**

```go
package main

import (
	"fmt"
)

func main() {
	greetingCh := make(chan string)
	
	go func() {
		greetingCh <- "Hola desde goroutine!"
	}()
	
	mensaje := <-greetingCh
	fmt.Println(mensaje)
}
```

**Conceptos probados:**
- Sintaxis make(chan Type)
- Operadores <- para send/receive
- Goroutines y channels
- Bloqueo implícito

---

### Ejercicio 2: Buffered Channel - Control de Flujo

**Objetivo:** Entender diferencia entre buffered y unbuffered.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// TODO: Completar el programa
	
	// Crear channel buffered con capacidad 3
	// dataCh := make(chan int, 3)
	
	// Enviar 3 valores SIN bloqueo (buffer tiene espacio)
	// dataCh <- 10
	// dataCh <- 20
	// dataCh <- 30
	
	// Esto SÍ se bloquearía (4to valor, buffer lleno)
	// dataCh <- 40  // Bloquea aquí
	
	// Resolver: antes de enviar el 4to, recibir uno
	// go func() {
	//     time.Sleep(1 * time.Second)
	//     val := <-dataCh  // Libera espacio en buffer
	//     fmt.Println("Recibido:", val)
	// }()
	
	// Ahora puedo enviar sin bloqueo
	// dataCh <- 40
	
	// TODO: Recibir los otros 3 valores
	// for i := 0; i < 3; i++ {
	//     fmt.Println("Recibido:", <-dataCh)
	// }
}

// Salida esperada:
// Recibido: 10
// Recibido: 20
// Recibido: 30
// Recibido: 40
```

**Solución:**

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	dataCh := make(chan int, 3)
	
	dataCh <- 10
	dataCh <- 20
	dataCh <- 30
	
	go func() {
		time.Sleep(1 * time.Second)
		val := <-dataCh
		fmt.Println("Recibido:", val)
	}()
	
	dataCh <- 40
	
	for i := 0; i < 3; i++ {
		fmt.Println("Recibido:", <-dataCh)
	}
}
```

**Conceptos probados:**
- Buffer size y capacidad
- Bloqueo condicional
- Non-blocking send/receive con buffer

---

### Ejercicio 3: Pipeline - Múltiples Etapas

**Objetivo:** Entender pipelines: generar → transformar → consumir.

```go
package main

import (
	"fmt"
)

// TODO: Implementar estas funciones

// generate: produce números del 1 al 5, luego cierra
// func generate(out chan<- int) {
// }

// square: recibe números, envía sus cuadrados
// func square(in <-chan int, out chan<- int) {
// }

// print: recibe números y los imprime
// func print(in <-chan int) {
// }

func main() {
	// TODO: Conectar el pipeline
	
	// nums := make(chan int)
	// squares := make(chan int)
	
	// go generate(nums)
	// go square(nums, squares)
	// print(squares)
}

// Salida esperada:
// 1
// 4
// 9
// 16
// 25
```

**Solución:**

```go
package main

import (
	"fmt"
)

func generate(out chan<- int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

func square(in <-chan int, out chan<- int) {
	for num := range in {
		out <- num * num
	}
	close(out)
}

func print(in <-chan int) {
	for result := range in {
		fmt.Println(result)
	}
}

func main() {
	nums := make(chan int)
	squares := make(chan int)
	
	go generate(nums)
	go square(nums, squares)
	print(squares)
}
```

**Conceptos probados:**
- Direccionalidad de channels (send-only, receive-only)
- Múltiples etapas
- close() y range
- for-range para iteración

---

### Ejercicio 4: Fan-Out/Fan-In - Distribución de Trabajo

**Objetivo:** Distribuir tareas a múltiples workers y recopilar resultados.

```go
package main

import (
	"fmt"
	"sync"
)

// TODO: Implementar worker
// worker: procesa un job (multiplica por 2) y envía resultado
// func worker(id int, jobs <-chan int, results chan<- int) {
// }

func main() {
	// TODO: Implementar fan-out/fan-in
	
	// Crear channels
	// jobs := make(chan int, 10)
	// results := make(chan int, 10)
	
	// Lanzar 3 workers con WaitGroup
	// var wg sync.WaitGroup
	// numWorkers := 3
	
	// for w := 1; w <= numWorkers; w++ {
	//     wg.Add(1)
	//     go func(id int) {
	//         defer wg.Done()
	//         worker(id, jobs, results)
	//     }(w)
	// }
	
	// Enviar jobs
	// for i := 1; i <= 10; i++ {
	//     jobs <- i
	// }
	// close(jobs)
	
	// Esperar a que terminen y cerrar results
	// go func() {
	//     wg.Wait()
	//     close(results)
	// }()
	
	// Recopilar resultados
	// for result := range results {
	//     fmt.Println("Resultado:", result)
	// }
}

// Salida esperada:
// Resultado: 2
// Resultado: 4
// (... 10 resultados, orden puede variar)
```

**Solución:**

```go
package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d procesa job %d\n", id, job)
		results <- job * 2
	}
}

func main() {
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	
	var wg sync.WaitGroup
	numWorkers := 3
	
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}
	
	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	close(jobs)
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	for result := range results {
		fmt.Println("Resultado:", result)
	}
}
```

**Conceptos probados:**
- sync.WaitGroup
- Múltiples goroutines enviando al mismo channel
- Fan-out de jobs
- Fan-in de resultados

---

### Ejercicio 5: Generador con Range - Iterator Pattern

**Objetivo:** Crear un generador que produce secuencia infinita (pero controlada).

```go
package main

import (
	"fmt"
)

// TODO: Implementar generador de Fibonacci
// fibonacci: devuelve channel que genera primeros N números Fibonacci
// func fibonacci(n int) <-chan int {
// }

func main() {
	// TODO: Usar el generador con range
	
	// Generar primeros 10 números Fibonacci
	// fmt.Println("Fibonacci (primeros 10):")
	// for num := range fibonacci(10) {
	//     fmt.Printf("%d ", num)
	// }
	// fmt.Println()
	
	// Bonus: Filtrar solo números pares
	// fmt.Println("\nFibonacci pares (primeros 6):")
	// count := 0
	// for num := range fibonacci(20) {
	//     if num%2 == 0 {
	//         fmt.Printf("%d ", num)
	//         count++
	//         if count == 6 {
	//             break
	//         }
	//     }
	// }
	// fmt.Println()
}

// Salida esperada:
// Fibonacci (primeros 10):
// 0 1 1 2 3 5 8 13 21 34
//
// Fibonacci pares (primeros 6):
// 0 2 8 34 144 610
```

**Solución:**

```go
package main

import (
	"fmt"
)

func fibonacci(n int) <-chan int {
	out := make(chan int)
	
	go func() {
		defer close(out)  // Cerrar cuando termina
		
		a, b := 0, 1
		for i := 0; i < n; i++ {
			out <- a
			a, b = b, a+b
		}
	}()
	
	return out
}

func main() {
	fmt.Println("Fibonacci (primeros 10):")
	for num := range fibonacci(10) {
		fmt.Printf("%d ", num)
	}
	fmt.Println()
	
	fmt.Println("\nFibonacci pares (primeros 6):")
	count := 0
	for num := range fibonacci(20) {
		if num%2 == 0 {
			fmt.Printf("%d ", num)
			count++
			if count == 6 {
				break
			}
		}
	}
	fmt.Println()
}
```

**Conceptos probados:**
- Generadores que devuelven `<-chan`
- defer close()
- range para iterar generadores
- break para terminar iteración temprano

---

## REFERENCIAS Y RECURSOS

**Documentación Oficial:**
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Channel axioms](https://github.com/golang/go/wiki/CodeReviewComments#chan)

**Libros Recomendados:**
- "The Go Programming Language" - Donovan & Kernighan (Capítulo 8)
- "Concurrency in Go" - Katherine Cox-Buday

**Patrones:**
- Rob Pike's "Concurrency Patterns" talks
- Go Concurrency Patterns: Pipeline, Fan-out/Fan-in
- Advanced Go Concurrency Patterns

---

**Fin del Capítulo 22: CHANNELS - COMUNICACIÓN ENTRE GOROUTINES**

*Este capítulo cubre el paradigma CSP de Go, permitiendo escribir programas concurrentes seguros sin race conditions. Los channels son la piedra angular de la filosofía "share memory by communicating" de Go.*
