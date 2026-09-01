# Capítulo 31: Time - Manejo de tiempo y duraciones

## Introducción

El package `time` de Go proporciona funcionalidades para medir y mostrar tiempo. Es fundamental en casi cualquier aplicación real: desde registrar eventos con timestamps hasta implementar timeouts en operaciones de red, ejecutar tareas periódicas o simplemente dormir una goroutine. 

Go proporciona abstracciones elegantes para trabajar con tiempo que abstraen la complejidad de los relojes del sistema operativo, zonas horarias y precisión. A diferencia de otros lenguajes, Go trata el tiempo de forma segura en concurrencia, lo que es crítico en un sistema diseñado alrededor de goroutines.

Este capítulo explora cómo manejar tiempo de manera robusta y eficiente en Go.

---

## 31.1 ¿Qué es el Time Package?

### 31.1.1 Conceptos Fundamentales

El package `time` es la abstracción de Go sobre el concepto de tiempo. Proporciona:

1. **Instantes de Tiempo**: Representación de puntos específicos en el tiempo (type `Time`)
2. **Duraciones**: Representación de intervalos de tiempo (type `Duration`)
3. **Timers y Tickers**: Mecanismos para esperar y ejecutar código periódicamente
4. **Parsing y Formatting**: Conversión entre strings y objetos `Time`
5. **Zonas Horarias**: Manejo de diferentes zonas horarias del mundo

### 31.1.2 Arquitectura Conceptual

```
┌─────────────────────────────────────────────────────┐
│          Aplicación Go (Tiempo)                     │
├─────────────────────────────────────────────────────┤
│         Time Package (Abstracciones Go)             │
│  ┌──────────┬──────────┬──────────┬──────────────┐  │
│  │  Time    │ Duration │ Location │  Timer/Tick  │  │
│  └──────────┴──────────┴──────────┴──────────────┘  │
├─────────────────────────────────────────────────────┤
│    Runtime Go + Syscalls del SO                      │
│  ┌──────────┬──────────┬──────────┐                 │
│  │clock_now │timezone  │ scheduler│                 │
│  └──────────┴──────────┴──────────┘                 │
├─────────────────────────────────────────────────────┤
│   Reloj del Hardware + Datos de Zonas Horarias      │
└─────────────────────────────────────────────────────┘
```

### 31.1.3 Unix Time y Epoch

Go usa **Unix time** internamente, que es el número de segundos desde el 1 de enero de 1970 a las 00:00:00 UTC (Unix epoch).

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Unix epoch
	epoch := time.Unix(0, 0)
	fmt.Println("Epoch:", epoch) // 1970-01-01 00:00:00 +0000 UTC
	
	// Timestamp actual en segundos
	now := time.Now()
	fmt.Println("Ahora:", now)
	fmt.Println("Unix segundos:", now.Unix())
	fmt.Println("Unix nanosegundos:", now.UnixNano())
	
	// Convertir desde Unix timestamp
	ts := time.Unix(1609459200, 0) // 2021-01-01 00:00:00 UTC
	fmt.Println("Desde timestamp:", ts)
}
```

**Salida esperada:**
```
Epoch: 1970-01-01 00:00:00 +0000 UTC
Ahora: 2024-01-15 14:30:45.123456789 +0100 CET
Unix segundos: 1705337445
Unix nanosegundos: 1705337445123456789
Desde timestamp: 2021-01-01 00:00:00 +0000 UTC
```

### 31.1.4 Precisión y Monótonos

Go proporciona dos relojes:

1. **Wall Clock**: Reloj de pared que puede retroceder (cambios de hora, NTP)
2. **Monotonic Clock**: Reloj monótono que nunca retrocede (interno, para medir duraciones)

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// El tipo time.Time almacena ambos relojes
	now := time.Now()
	
	// El reloj monotónico es automático para duraciones
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	duration := time.Since(start) // Usa reloj monotónico
	
	fmt.Println("Duración medida:", duration)
	// El resultado es siempre positivo y exacto
}
```

### 31.1.5 Comparación: Go vs Otros Lenguajes

```go
// ╔═══════════════════════════════════════════════════════════╗
// ║  Go - Simple, seguro para concurrencia                  ║
// ╚═══════════════════════════════════════════════════════════╝
now := time.Now()                    // Instant actual
duration := 5 * time.Second          // Duración
timer := time.After(duration)        // Esperar con goroutine

// ╔═══════════════════════════════════════════════════════════╗
// ║  Java - Más verboso                                      ║
// ╚═══════════════════════════════════════════════════════════╝
// Instant now = Instant.now();
// Duration duration = Duration.ofSeconds(5);
// Thread.sleep(5000);

// ╔═══════════════════════════════════════════════════════════╗
// ║  Python - Diferente filosofía                            ║
// ╚═══════════════════════════════════════════════════════════╝
// import datetime
// now = datetime.datetime.now()
// time.sleep(5)

// ╔═══════════════════════════════════════════════════════════╗
// ║  Rust - Específico de plataforma                         ║
// ╚═══════════════════════════════════════════════════════════╝
// use std::time::SystemTime;
// let now = SystemTime::now();
```

---

## 31.2 Time Type - Instantes de Tiempo

### 31.2.1 ¿Qué es un Time?

`time.Time` representa un instante específico en el tiempo con precisión de nanosegundos. Es un struct que contiene:

```go
type Time struct {
	// Contiene campos privados
	// - Reloj de pared (segundos desde Unix epoch)
	// - Reloj monotónico (para medir duraciones)
	// - Ubicación (zona horaria)
}
```

### 31.2.2 Obtener la Hora Actual

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Obtener hora actual
	now := time.Now()
	fmt.Println("Hora actual:", now)
	
	// Hora actual en UTC
	utc := time.Now().UTC()
	fmt.Println("UTC:", utc)
	
	// Solo la fecha
	today := time.Now().Truncate(24 * time.Hour)
	fmt.Println("Hoy a las 00:00:00:", today)
	
	// Componentes individuales
	fmt.Println("Año:", now.Year())
	fmt.Println("Mes:", now.Month())        // time.January, time.February, etc.
	fmt.Println("Día:", now.Day())
	fmt.Println("Hora:", now.Hour())
	fmt.Println("Minuto:", now.Minute())
	fmt.Println("Segundo:", now.Second())
	fmt.Println("Nanosegundos:", now.Nanosecond())
	fmt.Println("Zona horaria:", now.Location())
}
```

**Salida esperada:**
```
Hora actual: 2024-01-15 14:30:45.123456789 +0100 CET
UTC: 2024-01-15 13:30:45.123456789 +0000 UTC
Hoy a las 00:00:00: 2024-01-15 00:00:00 +0100 CET
Año: 2024
Mes: January
Día: 15
Hora: 14
...
```

### 31.2.3 Crear Time desde Componentes

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Crear desde componentes (en UTC)
	t1 := time.Date(2024, time.January, 15, 14, 30, 45, 123456789, time.UTC)
	fmt.Println("Desde componentes:", t1)
	
	// Crear desde Unix timestamp (segundos)
	t2 := time.Unix(1705337445, 0)
	fmt.Println("Desde Unix (segundos):", t2)
	
	// Crear desde Unix timestamp (nanosegundos)
	t3 := time.Unix(0, 1705337445123456789)
	fmt.Println("Desde Unix (nanosegundos):", t3)
	
	// Hora de inicio de Go
	goStartTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Println("Lanzamiento de Go:", goStartTime)
	
	// Última medianoche
	now := time.Now()
	lastMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	fmt.Println("Última medianoche:", lastMidnight)
}
```

### 31.2.4 Comparación y Operaciones

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t1 := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	
	// Comparación
	fmt.Println("t1 == t2:", t1.Equal(t2))      // false
	fmt.Println("t1 < t2:", t1.Before(t2))      // true
	fmt.Println("t1 > t2:", t1.After(t2))       // false
	
	// Diferencia entre dos instantes
	diff := t2.Sub(t1)
	fmt.Println("Diferencia:", diff) // 2h0m0s
	
	// Sumar duración a un instante
	t3 := t1.Add(3 * time.Hour)
	fmt.Println("t1 + 3 horas:", t3)
	
	// Agregar una duración de tiempo calendario (cuidado con zonas horarias)
	t4 := t1.AddDate(0, 1, 0) // Sumar 1 mes
	fmt.Println("t1 + 1 mes:", t4)
	
	// Verificar si es zero value
	var zeroTime time.Time
	fmt.Println("Zero time:", zeroTime)
	fmt.Println("¿Es zero?:", zeroTime.IsZero()) // true
}
```

**Salida esperada:**
```
t1 == t2: false
t1 < t2: true
t1 > t2: false
Diferencia: 2h0m0s
t1 + 3 horas: 2024-01-15 13:00:00 +0000 UTC
t1 + 1 mes: 2024-02-15 10:00:00 +0000 UTC
Zero time: 0001-01-01 00:00:00 +0000 UTC
¿Es zero?: true
```

### 31.2.5 Información del Día

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	
	// Día de la semana
	fmt.Println("Día de la semana:", now.Weekday()) // Monday, Tuesday, etc.
	
	// Día del año (1-366)
	fmt.Println("Día del año:", now.YearDay())
	
	// Información de daylight saving
	zone, offset := now.Zone()
	fmt.Println("Zona:", zone)     // "CET", "EST", etc.
	fmt.Println("Offset (segundos):", offset)
	
	// Semana ISO 8601
	year, week := now.ISOWeek()
	fmt.Println("Año ISO:", year, "Semana:", week)
	
	// Hora del día como fracción
	h, m, s := now.Clock()
	fmt.Printf("Reloj: %02d:%02d:%02d\n", h, m, s)
}
```

---

## 31.3 Duration - Duraciones de Tiempo

### 31.3.1 ¿Qué es una Duration?

`time.Duration` es un número entero que representa un intervalo de tiempo con precisión de nanosegundos.

```go
type Duration int64 // Nanosegundos internamente
```

### 31.3.2 Creando Duraciones

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Duraciones predefinidas
	fmt.Println("1 segundo:", time.Second)        // 1s
	fmt.Println("1 minuto:", time.Minute)         // 1m0s
	fmt.Println("1 hora:", time.Hour)             // 1h0m0s
	fmt.Println("1 milisegundo:", time.Millisecond) // 1ms
	fmt.Println("1 microsegundo:", time.Microsecond) // 1µs
	fmt.Println("1 nanosegundo:", time.Nanosecond) // 1ns
	
	// Combinando unidades
	delay := 2*time.Hour + 30*time.Minute + 45*time.Second
	fmt.Println("2h 30m 45s:", delay)
	
	// Desde string (formato: "300ms", "1.5h", "2h45m")
	d, _ := time.ParseDuration("1h30m")
	fmt.Println("Desde string:", d)
	
	// Operaciones aritméticas
	d1 := 5 * time.Second
	d2 := 3 * time.Second
	fmt.Println("5s + 3s =", d1+d2)
	fmt.Println("5s - 3s =", d1-d2)
	fmt.Println("5s * 2 =", d1*2)
	fmt.Println("5s / 2 =", d1/2)
}
```

**Salida esperada:**
```
1 segundo: 1s
1 minuto: 1m0s
1 hora: 1h0m0s
1 milisegundo: 1ms
1 microsegundo: 1µs
1 nanosegundo: 1ns
2h 30m 45s: 2h30m45s
Desde string: 1h30m
5s + 3s = 8s
5s - 3s = 2s
5s * 2 = 10s
5s / 2 = 2s500ms
```

### 31.3.3 Conversión Entre Unidades

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	d := 2*time.Hour + 30*time.Minute + 45*time.Second + 500*time.Millisecond
	
	// Conversión a diferentes unidades
	fmt.Println("Nanosegundos:", d.Nanoseconds())      // 9045500000000
	fmt.Println("Microsegundos:", d.Microseconds())    // 9045500000
	fmt.Println("Milisegundos:", d.Milliseconds())     // 9045500
	fmt.Println("Segundos:", d.Seconds())              // 9045.5
	fmt.Println("Minutos:", d.Minutes())               // 150.75
	fmt.Println("Horas:", d.Hours())                   // 2.5125
	
	// Desglose manual
	h := d / time.Hour
	remainder := d % time.Hour
	m := remainder / time.Minute
	remainder = remainder % time.Minute
	s := remainder / time.Second
	
	fmt.Printf("Desglose: %d horas, %d minutos, %d segundos\n", h, m, s)
}
```

### 31.3.4 String y Parsing

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Duration a string
	d := 1*time.Hour + 30*time.Minute + 45*time.Second
	fmt.Println("String:", d.String()) // 1h30m45s
	
	// Parsing desde string
	valores := []string{"1s", "5m", "2h", "1.5h", "500ms", "1000000ns"}
	
	for _, v := range valores {
		d, err := time.ParseDuration(v)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Printf("%-10s = %v\n", v, d)
	}
	
	// Validación de formato
	_, err := time.ParseDuration("1 hour") // Error: formato inválido
	fmt.Println("Error esperado:", err)
}
```

### 31.3.5 Casos de Uso Comunes

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Timeout de operación
	operationTimeout := 30 * time.Second
	
	// Retry backoff exponencial
	retryDelays := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
	}
	
	// Rate limiting
	rateLimit := time.Second / 100 // 100 requests por segundo
	
	// Polling interval
	pollInterval := 5 * time.Second
	
	fmt.Println("Timeout:", operationTimeout)
	fmt.Println("Retry delays:", retryDelays)
	fmt.Println("Rate limit:", rateLimit)
	fmt.Println("Poll interval:", pollInterval)
	
	// Calcular time budget
	deadline := time.Now().Add(1 * time.Minute)
	remaining := time.Until(deadline)
	fmt.Println("Tiempo restante:", remaining)
}
```

---

## 31.4 Timers - Temporizadores

### 31.4.1 ¿Qué es un Timer?

Un `time.Timer` envía un evento a través de un canal después de una duración especificada. Es el mecanismo de Go para ejecutar código después de un cierto tiempo.

```go
type Timer struct {
	C <-chan time.Time  // Canal que recibe el evento
	// Campos privados...
}
```

### 31.4.2 Timer Básico

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Iniciando timer...")
	
	// Crear timer para 2 segundos
	timer := time.NewTimer(2 * time.Second)
	
	// Esperar el evento en el canal
	<-timer.C
	
	fmt.Println("¡Timer disparado después de 2 segundos!")
}
```

**Salida esperada:**
```
Iniciando timer...
(espera 2 segundos)
¡Timer disparado después de 2 segundos!
```

### 31.4.3 After - La Forma Simplificada

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Esperando con After...")
	
	// time.After es equivalente a time.NewTimer pero devuelve solo el canal
	<-time.After(2 * time.Second)
	
	fmt.Println("¡Hecho después de 2 segundos!")
	
	// Útil en select para timeouts
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout alcanzado")
	case result := <-someChannel:
		fmt.Println("Resultado recibido:", result)
	}
}

// Canal ficticio para ejemplo
var someChannel <-chan string = nil
```

### 31.4.4 Stop y Reset

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Crear timer
	timer := time.NewTimer(5 * time.Second)
	
	// Después de 1 segundo, detener el timer
	go func() {
		time.Sleep(1 * time.Second)
		if timer.Stop() {
			fmt.Println("Timer detenido antes de dispararse")
		}
	}()
	
	// Esperar resultado
	select {
	case <-timer.C:
		fmt.Println("Timer disparó")
	}
	
	// Reset - reutilizar el mismo timer
	timer.Reset(2 * time.Second)
	fmt.Println("Timer reiniciado para 2 segundos")
	<-timer.C
	fmt.Println("Timer disparó después de reset")
}
```

**Salida esperada:**
```
Timer detenido antes de dispararse
Timer reiniciado para 2 segundos
Timer disparó después de reset
```

### 31.4.5 AfterFunc - Ejecutar Función

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	
	// Ejecutar función después de 2 segundos
	timer := time.AfterFunc(2*time.Second, func() {
		fmt.Println("Función ejecutada después de 2 segundos")
		wg.Done()
	})
	
	fmt.Println("Esperando ejecución de función...")
	wg.Wait()
	fmt.Println("Completado")
	
	// Detener ejecución si no ha ocurrido
	timer.Stop()
}
```

### 31.4.6 Patrón: Timeout en Operaciones

```go
package main

import (
	"fmt"
	"time"
)

// Simula operación que tarda tiempo
func operacionLarga() <-chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "Resultado de operación"
	}()
	return ch
}

func main() {
	// Con timeout de 2 segundos
	select {
	case resultado := <-operacionLarga():
		fmt.Println("Éxito:", resultado)
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout: operación tardó demasiado")
	}
	
	// Con timeout de 5 segundos
	select {
	case resultado := <-operacionLarga():
		fmt.Println("Éxito:", resultado)
	case <-time.After(5 * time.Second):
		fmt.Println("Timeout: operación tardó demasiado")
	}
}
```

---

## 31.5 Tickers - Repetición Periódica

### 31.5.1 ¿Qué es un Ticker?

Un `time.Ticker` envía eventos periódicamente a través de un canal en intervalos regulares.

```go
type Ticker struct {
	C <-chan time.Time  // Canal que recibe eventos periódicamente
	// Campos privados...
}
```

### 31.5.2 Ticker Básico

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Ticker que dispara cada 500ms
	ticker := time.NewTicker(500 * time.Millisecond)
	
	// Ejecutar durante 2.5 segundos
	stop := time.After(2500 * time.Millisecond)
	
	for {
		select {
		case t := <-ticker.C:
			fmt.Println("Tick:", t.Format("15:04:05.000"))
		case <-stop:
			ticker.Stop()
			fmt.Println("Ticker detenido")
			return
		}
	}
}
```

**Salida esperada:**
```
Tick: 14:30:45.123
Tick: 14:30:45.623
Tick: 14:30:46.123
Tick: 14:30:46.623
Tick: 14:30:47.123
Ticker detenido
```

### 31.5.3 Tick - La Forma Simplificada

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// time.Tick devuelve solo el canal
	// Útil para loops simples
	
	c := time.Tick(100 * time.Millisecond)
	
	for i := 0; i < 5; i++ {
		fmt.Println("Tick:", <-c)
	}
	
	fmt.Println("Done")
	// Nota: time.Tick nunca puede ser detenido
	// Usar time.NewTicker si necesitas detener el ticker
}
```

### 31.5.4 Stop y Reset

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	
	// Dejar que dispare 3 veces
	for i := 0; i < 3; i++ {
		<-ticker.C
		fmt.Println("Tick", i+1)
	}
	
	// Detener el ticker
	ticker.Stop()
	fmt.Println("Ticker detenido")
	
	// Crear nuevo ticker con intervalo diferente
	ticker = time.NewTicker(500 * time.Millisecond)
	for i := 0; i < 2; i++ {
		<-ticker.C
		fmt.Println("Nuevo tick", i+1)
	}
	ticker.Stop()
}
```

### 31.5.5 Patrón: Task Scheduler

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// TaskScheduler ejecuta tareas periódicamente
type TaskScheduler struct {
	ticker *time.Ticker
	done   chan struct{}
	wg     sync.WaitGroup
}

func NewTaskScheduler(interval time.Duration) *TaskScheduler {
	return &TaskScheduler{
		ticker: time.NewTicker(interval),
		done:   make(chan struct{}),
	}
}

func (ts *TaskScheduler) Start(task func()) {
	ts.wg.Add(1)
	go func() {
		defer ts.wg.Done()
		for {
			select {
			case <-ts.ticker.C:
				task()
			case <-ts.done:
				return
			}
		}
	}()
}

func (ts *TaskScheduler) Stop() {
	ts.ticker.Stop()
	close(ts.done)
	ts.wg.Wait()
}

func main() {
	scheduler := NewTaskScheduler(1 * time.Second)
	
	counter := 0
	scheduler.Start(func() {
		counter++
		fmt.Printf("Tarea ejecutada: %d\n", counter)
	})
	
	// Dejar correr por 3.5 segundos
	time.Sleep(3500 * time.Millisecond)
	scheduler.Stop()
	fmt.Println("Scheduler detenido")
}
```

---

## 31.6 Parsing y Formatting - Convertir Fechas

### 31.6.1 Layouts de Formato

Go usa **layouts de referencia** en lugar de format strings como strftime. El formato de referencia es: `Mon Jan 2 15:04:05 MST 2006`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Date(2024, time.January, 15, 14, 30, 45, 0, time.UTC)
	
	// Layouts predefinidos
	fmt.Println("RFC3339:", t.Format(time.RFC3339))           // 2024-01-15T14:30:45Z
	fmt.Println("RFC822:", t.Format(time.RFC822))            // 15 Jan 24 14:30 UTC
	fmt.Println("RFC1123:", t.Format(time.RFC1123))          // Mon, 15 Jan 2024 14:30:45 UTC
	fmt.Println("RFC3339Nano:", t.Format(time.RFC3339Nano))  // 2024-01-15T14:30:45Z
	fmt.Println("UnixDate:", t.Format(time.UnixDate))        // Mon Jan 15 14:30:45 UTC 2024
	fmt.Println("RubyDate:", t.Format(time.RubyDate))        // Mon Jan 15 14:30:45 +0000 2024
	fmt.Println("ANSIC:", t.Format(time.ANSIC))              // Mon Jan 15 14:30:45 2024
}
```

### 31.6.2 Layouts Personalizados

La clave es usar la **fecha de referencia**: `Mon Jan 2 15:04:05 MST 2006`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Date(2024, time.January, 15, 14, 30, 45, 0, time.UTC)
	
	// Layouts personalizados usando la referencia
	layouts := map[string]string{
		"ISO Date":      "2006-01-02",                       // 2024-01-15
		"US Format":     "01/02/2006",                       // 01/15/2024
		"European":      "02-01-2006",                       // 15-01-2024
		"Long":          "2 January 2006 15:04:05",         // 15 January 2024 14:30:45
		"Short":         "02 Jan 06 15:04",                 // 15 Jan 24 14:30
		"ISO DateTime":  "2006-01-02T15:04:05",             // 2024-01-15T14:30:45
		"ISO DateTime Z":"2006-01-02T15:04:05Z07:00",       // 2024-01-15T14:30:45+00:00
	}
	
	for nombre, layout := range layouts {
		fmt.Printf("%-15s: %s\n", nombre, t.Format(layout))
	}
}
```

**Salida esperada:**
```
ISO Date       : 2024-01-15
US Format      : 01/15/2024
European       : 15-01-2024
Long           : 15 January 2024 14:30:45
Short          : 15 Jan 24 14:30
ISO DateTime   : 2024-01-15T14:30:45
ISO DateTime Z : 2024-01-15T14:30:45+00:00
```

### 31.6.3 Parsing desde String

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Parse con layout específico
	layouts := []struct {
		name   string
		layout string
		value  string
	}{
		{"ISO Date", "2006-01-02", "2024-01-15"},
		{"US Format", "01/02/2006", "01/15/2024"},
		{"European", "02-01-2006", "15-01-2024"},
		{"RFC3339", time.RFC3339, "2024-01-15T14:30:45Z"},
		{"Time only", "15:04:05", "14:30:45"},
	}
	
	for _, test := range layouts {
		t, err := time.Parse(test.layout, test.value)
		if err != nil {
			fmt.Printf("%s - Error: %v\n", test.name, err)
			continue
		}
		fmt.Printf("%s - %v\n", test.name, t)
	}
}
```

### 31.6.4 ParseInLocation - Controlar la Zona Horaria

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Sin zona horaria, asumir UTC
	t1, _ := time.Parse("2006-01-02 15:04:05", "2024-01-15 14:30:45")
	fmt.Println("UTC:", t1)
	
	// Sin zona horaria, asumir una zona específica
	ny, _ := time.LoadLocation("America/New_York")
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-15 14:30:45", ny)
	fmt.Println("Nueva York:", t2)
	
	// RFC3339 incluye zona horaria
	t3, _ := time.Parse(time.RFC3339, "2024-01-15T14:30:45+01:00")
	fmt.Println("Con zona:", t3)
	
	// Nota: el mismo tiempo en diferentes representaciones
	fmt.Println("T1 igual T3:", t1.Equal(time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)))
}
```

### 31.6.5 Formatos Comunes en Aplicaciones Reales

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	
	// Logs
	logFormat := "2006-01-02 15:04:05" // o time.DateTime
	fmt.Println("Log:", t.Format(logFormat))
	
	// API JSON
	jsonFormat := time.RFC3339Nano
	fmt.Println("JSON:", t.Format(jsonFormat))
	
	// Database (SQL)
	dbFormat := "2006-01-02 15:04:05"
	fmt.Println("DB:", t.Format(dbFormat))
	
	// CSV
	csvFormat := "2006-01-02"
	fmt.Println("CSV:", t.Format(csvFormat))
	
	// Timestamp para archivos
	fileFormat := "20060102_150405"
	fmt.Println("Archivo:", t.Format(fileFormat))
	
	// Día legible
	dayFormat := "Monday, January 02, 2006"
	fmt.Println("Legible:", t.Format(dayFormat))
}
```

---

## 31.7 Manejo de Zonas Horarias

### 31.7.1 Zonas Horarias Básicas

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// UTC (siempre disponible)
	utc := time.Now().UTC()
	fmt.Println("UTC:", utc)
	
	// Zona horaria local del sistema
	local := time.Now()
	fmt.Println("Local:", local)
	
	// Hora en ambas zonas
	t := time.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)
	fmt.Println("UTC:", t)
	fmt.Println("Local:", t.In(time.Local))
	
	// Convertir entre zonas
	ny, _ := time.LoadLocation("America/New_York")
	fmt.Println("Nueva York:", t.In(ny))
	
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	fmt.Println("Tokyo:", t.In(tokyo))
}
```

### 31.7.2 LoadLocation - Cargar Zona Horaria

```go
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Cargar zona horaria por nombre IANA
	zones := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Europe/Paris",
		"Asia/Tokyo",
		"Australia/Sydney",
		"Pacific/Auckland",
	}
	
	for _, zoneName := range zones {
		loc, err := time.LoadLocation(zoneName)
		if err != nil {
			log.Fatal(err)
		}
		
		t := time.Now().In(loc)
		zone, offset := t.Zone()
		fmt.Printf("%-20s: %s (UTC%+d)\n", zoneName, zone, offset/3600)
	}
}
```

**Salida esperada:**
```
UTC                 : UTC (UTC+0)
America/New_York    : EST (UTC-5)
Europe/London       : GMT (UTC+0)
Europe/Paris        : CET (UTC+1)
Asia/Tokyo          : JST (UTC+9)
Australia/Sydney    : AEDT (UTC+11)
Pacific/Auckland    : NZDT (UTC+13)
```

### 31.7.3 FixedZone - Zona Horaria Fija

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Crear zona horaria con offset fijo
	// time.FixedZone(name, offset_en_segundos)
	
	ist := time.FixedZone("IST", 5*3600+30*60) // UTC+5:30
	cest := time.FixedZone("CEST", 2*3600)     // UTC+2
	
	t := time.Date(2024, time.January, 15, 14, 30, 0, 0, ist)
	fmt.Println("IST:", t)
	
	// Convertir a otra zona
	fmt.Println("CEST:", t.In(cest))
	
	// Información de zona
	zone, offset := t.Zone()
	fmt.Printf("Zona: %s, Offset: %d segundos\n", zone, offset)
}
```

### 31.7.4 DST - Daylight Saving Time

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Zona con DST
	ny, _ := time.LoadLocation("America/New_York")
	
	// Invierno (sin DST)
	winter := time.Date(2024, time.January, 15, 12, 0, 0, 0, ny)
	zone1, offset1 := winter.Zone()
	fmt.Printf("Enero: %s, offset: %d\n", zone1, offset1/3600)
	
	// Verano (con DST)
	summer := time.Date(2024, time.July, 15, 12, 0, 0, 0, ny)
	zone2, offset2 := summer.Zone()
	fmt.Printf("Julio: %s, offset: %d\n", zone2, offset2/3600)
	
	// Ambos en UTC
	fmt.Println("Diferencia UTC:", winter.Sub(summer).Hours(), "horas")
}
```

### 31.7.5 Antipatrones Comunes

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// ❌ ANTIPATRÓN: Ignorar zonas horarias
	t1, _ := time.Parse("2006-01-02 15:04:05", "2024-01-15 14:30:00")
	t2, _ := time.Parse("2006-01-02 15:04:05", "2024-01-15 14:30:00")
	// Ambas son UTC por defecto, pero si vinieron de diferentes fuentes,
	// puede no ser correcto
	
	// ✅ CORRECTO: Usar ParseInLocation o RFC3339
	ny, _ := time.LoadLocation("America/New_York")
	t3, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-15 14:30:00", ny)
	fmt.Println("Correcto:", t3)
	
	// ❌ ANTIPATRÓN: Asumir zona local en APIs
	// No guardar zone horaria en bases de datos
	
	// ✅ CORRECTO: Siempre almacenar en UTC
	now := time.Now().UTC()
	fmt.Println("Almacenar:", now.Format(time.RFC3339))
}
```

---

## 31.8 Medición de Tiempo - Benchmarking

### 31.8.1 Medir Tiempo Transcurrido

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Método 1: time.Now() al inicio y fin
	start := time.Now()
	
	// Código a medir
	time.Sleep(100 * time.Millisecond)
	
	elapsed := time.Since(start)
	fmt.Println("Tiempo transcurrido:", elapsed)
	fmt.Println("En segundos:", elapsed.Seconds())
	fmt.Println("En milisegundos:", elapsed.Milliseconds())
	
	// Método 2: Usando time.Until
	deadline := time.Now().Add(5 * time.Second)
	time.Sleep(2 * time.Second)
	remaining := time.Until(deadline)
	fmt.Println("Tiempo restante:", remaining)
}
```

### 31.8.2 Stopwatch - Estructura Auxiliar

```go
package main

import (
	"fmt"
	"time"
)

// Stopwatch mide intervalos de tiempo
type Stopwatch struct {
	start time.Time
	stop  time.Time
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{start: time.Now()}
}

func (s *Stopwatch) Stop() time.Duration {
	s.stop = time.Now()
	return s.Elapsed()
}

func (s *Stopwatch) Elapsed() time.Duration {
	if s.stop.IsZero() {
		return time.Since(s.start)
	}
	return s.stop.Sub(s.start)
}

func (s *Stopwatch) Reset() {
	s.start = time.Now()
	s.stop = time.Time{}
}

func main() {
	sw := NewStopwatch()
	
	// Simular operación
	time.Sleep(150 * time.Millisecond)
	elapsed := sw.Stop()
	
	fmt.Println("Tiempo:", elapsed)
	fmt.Printf("%.2f ms\n", elapsed.Seconds()*1000)
}
```

### 31.8.3 Benchmarking Simple

```go
package main

import (
	"fmt"
	"time"
)

// Función para medir
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func benchmark(name string, fn func(), iterations int) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	elapsed := time.Since(start)
	
	fmt.Printf("%s: %v (%.2f ns/op)\n",
		name,
		elapsed,
		float64(elapsed.Nanoseconds())/float64(iterations))
}

func main() {
	benchmark("fib(20)", func() { fibonacci(20) }, 1000)
	benchmark("fib(25)", func() { fibonacci(25) }, 100)
}
```

### 31.8.4 Comparación de Rendimiento

```go
package main

import (
	"fmt"
	"time"
)

func benchmarkOperation(name string, op func()) {
	start := time.Now()
	op()
	elapsed := time.Since(start)
	fmt.Printf("%-20s: %v\n", name, elapsed)
}

func main() {
	// Comparar diferentes operaciones
	benchmarkOperation("Sleep 10ms", func() {
		time.Sleep(10 * time.Millisecond)
	})
	
	benchmarkOperation("Create 1M ints", func() {
		for i := 0; i < 1000000; i++ {
			_ = i
		}
	})
	
	benchmarkOperation("String concat", func() {
		s := ""
		for i := 0; i < 10000; i++ {
			s += "x"
		}
	})
}
```

---

## 31.9 Sleep y Delays - Dormir Goroutines

### 31.9.1 time.Sleep - La Forma Básica

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Iniciando a:", time.Now().Format("15:04:05"))
	
	// Dormir por 2 segundos
	time.Sleep(2 * time.Second)
	
	fmt.Println("Despierto a:", time.Now().Format("15:04:05"))
	
	// Sleep no bloquea otras goroutines
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Goroutine 1 despierta")
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Goroutine 2 despierta")
	}()
	
	// Main goroutine duerme
	time.Sleep(3 * time.Second)
}
```

### 31.9.2 Exponential Backoff - Retry con Delay

```go
package main

import (
	"fmt"
	"math"
	"time"
)

// RetryWithBackoff reintenta una operación con backoff exponencial
func retryWithBackoff(maxRetries int, operation func() error) error {
	var err error
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = operation()
		if err == nil {
			return nil // Éxito
		}
		
		if attempt < maxRetries-1 {
			// Backoff exponencial: 2^attempt segundos, con máximo de 32s
			delay := time.Duration(math.Min(
				math.Pow(2, float64(attempt)),
				32,
			)) * time.Second
			
			fmt.Printf("Intento %d falló, esperando %v\n", attempt+1, delay)
			time.Sleep(delay)
		}
	}
	
	return err
}

// Simular operación que falla
var attemptCount = 0

func unreliableOperation() error {
	attemptCount++
	fmt.Printf("Intento %d...", attemptCount)
	if attemptCount < 3 {
		fmt.Println(" FALLO")
		return fmt.Errorf("error temporal")
	}
	fmt.Println(" ÉXITO")
	return nil
}

func main() {
	err := retryWithBackoff(5, unreliableOperation)
	if err != nil {
		fmt.Println("Error final:", err)
	}
}
```

### 31.9.3 Linear Backoff - Alternativa

```go
package main

import (
	"fmt"
	"time"
)

// LinearBackoff: delay = attempt * delayPerAttempt
func retryLinearBackoff(maxRetries int, baseDelay time.Duration, operation func() error) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		
		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			fmt.Printf("Reintentando en %v\n", delay)
			time.Sleep(delay)
		}
	}
	
	return fmt.Errorf("operación falló después de %d intentos", maxRetries)
}

func main() {
	counter := 0
	err := retryLinearBackoff(5, 100*time.Millisecond, func() error {
		counter++
		fmt.Println("Intento:", counter)
		if counter < 3 {
			return fmt.Errorf("fallo")
		}
		return nil
	})
	
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

### 31.9.4 Jitter - Evitar Thundering Herd

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// BackoffWithJitter: Exponential backoff + random jitter
func retryWithJitter(maxRetries int, operation func() error) error {
	rand.Seed(time.Now().UnixNano())
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		
		if attempt < maxRetries-1 {
			// Exponential base
			baseDelay := math.Pow(2, float64(attempt))
			
			// Agregar jitter aleatorio: +/- 25%
			jitterFactor := 0.75 + rand.Float64()*0.5 // [0.75, 1.25]
			delay := time.Duration(baseDelay*jitterFactor) * time.Second
			
			fmt.Printf("Esperando %v (attempt %d)\n", delay, attempt+1)
			time.Sleep(delay)
		}
	}
	
	return fmt.Errorf("falló después de %d intentos", maxRetries)
}

func main() {
	retryWithJitter(4, func() error {
		fmt.Println("Intentando...")
		return fmt.Errorf("fallo")
	})
}
```

### 31.9.5 Circuit Breaker Pattern

```go
package main

import (
	"fmt"
	"time"
)

// CircuitBreaker implementa el patrón circuit breaker
type CircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	failureCount     int
	lastFailureTime  time.Time
	state            string // "closed", "open", "half-open"
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: threshold,
		resetTimeout:     timeout,
		state:            "closed",
	}
}

func (cb *CircuitBreaker) Call(operation func() error) error {
	// Si está abierto, verificar si es hora de reintentar
	if cb.state == "open" {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = "half-open"
			cb.failureCount = 0
		} else {
			return fmt.Errorf("circuit breaker abierto")
		}
	}
	
	// Intentar operación
	err := operation()
	
	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		
		if cb.failureCount >= cb.failureThreshold {
			cb.state = "open"
			return fmt.Errorf("circuit breaker abierto después de %d fallos", cb.failureCount)
		}
		
		return err
	}
	
	// Éxito: resetear
	cb.failureCount = 0
	cb.state = "closed"
	return nil
}

func main() {
	cb := NewCircuitBreaker(3, 2*time.Second)
	
	// Simular operación que falla
	for i := 0; i < 5; i++ {
		err := cb.Call(func() error {
			return fmt.Errorf("servicio no disponible")
		})
		fmt.Printf("Intento %d: %v (estado: %s)\n", i+1, err, cb.state)
	}
}
```

---

## 31.10 Context Deadlines - Timeouts Estructurados

### 31.10.1 Context con Deadline

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Crear contexto con deadline
	ctx, cancel := context.WithDeadline(context.Background(),
		time.Now().Add(2*time.Second))
	defer cancel()
	
	// Realizar operación
	select {
	case <-doWork(ctx):
		fmt.Println("Trabajo completado a tiempo")
	case <-ctx.Done():
		fmt.Println("Deadline alcanzado:", ctx.Err())
	}
}

func doWork(ctx context.Context) <-chan bool {
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(3 * time.Second) // Tarda más que deadline
		ch <- true
	}()
	return ch
}
```

### 31.10.2 Context con Timeout

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// Simular operación que puede tardar
func fetchData(ctx context.Context, delayMS int) (string, error) {
	select {
	case <-time.After(time.Duration(delayMS) * time.Millisecond):
		return "datos", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	// Timeout de 500ms
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	
	// Operación que tarda 200ms (éxito)
	result, err := fetchData(ctx, 200)
	fmt.Println("Rápido:", result, err)
	
	// Operación que tarda 1000ms (timeout)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel2()
	
	result2, err2 := fetchData(ctx2, 1000)
	fmt.Println("Lento:", result2, err2)
}
```

### 31.10.3 Propagación de Deadlines

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// Función que respeta contexto
func processRequest(ctx context.Context, userID string) error {
	// Derivar contexto para suboperación
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	
	return fetchUserData(ctx, userID)
}

func fetchUserData(ctx context.Context, userID string) error {
	select {
	case <-time.After(2 * time.Second):
		return fmt.Errorf("datos fetched")
	case <-ctx.Done():
		return ctx.Err() // Contexto agotado
	}
}

func main() {
	// Request con deadline general de 500ms
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	
	err := processRequest(ctx, "user123")
	fmt.Println("Error:", err)
	// Output: Error: context deadline exceeded
}
```

### 31.10.4 HTTP Client con Timeout

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Crear cliente con timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	// O usar context con deadline
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	fmt.Println("Status:", resp.Status)
	resp.Body.Close()
}
```

---

## 31.11 Buenas Prácticas y Patrones

### 31.11.1 Awareness de Zonas Horarias

```go
package main

import (
	"fmt"
	"time"
)

// ✅ CORRECTO: Siempre almacenar en UTC
func storeEvent(eventName string, timestamp time.Time) string {
	// Convertir a UTC antes de almacenar
	utcTime := timestamp.UTC()
	return fmt.Sprintf("%s happened at %s", eventName, utcTime.Format(time.RFC3339))
}

// ✅ CORRECTO: Mostrar en zona local del usuario
func displayToUser(utcTime time.Time, userTimezone string) string {
	loc, _ := time.LoadLocation(userTimezone)
	localTime := utcTime.In(loc)
	return fmt.Sprintf("Ocurrió a las %s en tu zona horaria", localTime.Format("15:04:05"))
}

func main() {
	now := time.Now()
	
	stored := storeEvent("login", now)
	fmt.Println(stored)
	
	display := displayToUser(now, "America/New_York")
	fmt.Println(display)
}
```

### 31.11.2 Precisión y Medición

```go
package main

import (
	"fmt"
	"time"
)

// ✅ CORRECTO: Usar time.Since para medir (reloj monotónico)
func measureOperation(operation func()) time.Duration {
	start := time.Now()
	operation()
	return time.Since(start)
}

// ❌ INCORRECTO: Usar time.Sub con hora actual (puede tener problemas de NTP)
func measureOperationBad(operation func()) time.Duration {
	start := time.Now()
	operation()
	end := time.Now()
	return end.Sub(start) // Funciona pero menos robusto
}

func main() {
	duration := measureOperation(func() {
		time.Sleep(100 * time.Millisecond)
	})
	
	fmt.Println("Operación tardó:", duration)
}
```

### 31.11.3 Testing con Tiempo

```go
package main

import (
	"fmt"
	"testing"
	"time"
)

// Función a probar que usa tiempo
func waitAndReturn(duration time.Duration) string {
	time.Sleep(duration)
	return "done"
}

// Usar time.After en tests para no esperar realmente
func TestWithTimeout(t *testing.T) {
	done := make(chan string, 1)
	
	go func() {
		done <- waitAndReturn(100 * time.Millisecond)
	}()
	
	select {
	case result := <-done:
		fmt.Println("Test passed:", result)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Test timeout")
		t.Fatal("operación tardó demasiado")
	}
}

// Usar mocks de tiempo para tests más rápidos
type MockClock struct {
	now time.Time
}

func (mc *MockClock) Now() time.Time {
	return mc.now
}

func main() {
	mock := &MockClock{now: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)}
	fmt.Println("Mock time:", mock.Now())
}
```

### 31.11.4 Manejo de Errores en Parsing

```go
package main

import (
	"fmt"
	"log"
	"time"
)

// ✅ CORRECTO: Validar error al parsear
func parseTime(dateString string) (time.Time, error) {
	// Usar layout específico, no adivinar
	return time.Parse("2006-01-02", dateString)
}

// ✅ CORRECTO: Proporcionar mensaje de error descriptivo
func parseTimeWithError(dateString string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"invalid date format '%s': expected YYYY-MM-DD, got error: %w",
			dateString, err)
	}
	return t, nil
}

func main() {
	// ✅ Validar al parsear
	dates := []string{"2024-01-15", "invalid", "2024/01/15"}
	
	for _, date := range dates {
		t, err := parseTimeWithError(date)
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		fmt.Println("Parsed:", t)
	}
}
```

### 31.11.5 Performance - Caché de Layouts

```go
package main

import (
	"fmt"
	"time"
)

// ❌ ANTIPATRÓN: Crear layout cada vez
func formatBad(times []time.Time) {
	for _, t := range times {
		// Esto es ineficiente si se hace muchas veces
		fmt.Println(t.Format("2006-01-02 15:04:05"))
	}
}

// ✅ CORRECTO: Reutilizar layout
const dateTimeLayout = "2006-01-02 15:04:05"

func formatGood(times []time.Time) {
	for _, t := range times {
		fmt.Println(t.Format(dateTimeLayout))
	}
}

func main() {
	times := []time.Time{
		time.Now(),
		time.Now().Add(1 * time.Hour),
		time.Now().Add(2 * time.Hour),
	}
	
	formatGood(times)
}
```

### 31.11.6 Sincronización con NTP

```go
package main

import (
	"fmt"
	"net"
	"time"
)

// Comprobar sincronización con servidor NTP
func checkTimeSync() {
	// En producción, usar "pool.ntp.org" o servidor específico
	// Este es solo un ejemplo conceptual
	
	// Go usa el reloj del sistema, que debería estar sincronizado
	localTime := time.Now()
	fmt.Println("Hora local:", localTime)
	
	// En aplicaciones críticas, verificar sincronización:
	// - Usar ntpd en Linux
	// - Usar w32tm en Windows
	// - Monitorear desviaciones
	
	_ = net.Dialer{} // Placeholder
}

func main() {
	checkTimeSync()
}
```

---

## Ejercicios Prácticos Progresivos

### Ejercicio 1: Timer Simple - Esperar con time.After

**Objetivo**: Crear una función que espere un tiempo determinado y muestre un mensaje.

```go
// TODO: Completar la implementación
package main

import (
	"fmt"
	"time"
)

// WaitAndPrint espera `duration` y luego imprime `message`
func waitAndPrint(duration time.Duration, message string) {
	// Tu código aquí
	// 1. Usar time.After para crear un canal que dispare después de `duration`
	// 2. Leer del canal con <-
	// 3. Imprimir el mensaje
}

func main() {
	fmt.Println("Iniciando...")
	waitAndPrint(2*time.Second, "¡2 segundos han pasado!")
	fmt.Println("Fin")
}
```

**Solución esperada:**
```go
func waitAndPrint(duration time.Duration, message string) {
	<-time.After(duration)
	fmt.Println(message)
}
```

---

### Ejercicio 2: Ticker - Ejecutar Código Periódicamente

**Objetivo**: Crear un conteo que imprime números cada segundo durante 5 segundos.

```go
// TODO: Completar la implementación
package main

import (
	"fmt"
	"time"
)

// CountWith Ticker imprime números del 1 al n, esperando interval entre cada uno
func countWithTicker(n int, interval time.Duration) {
	// Tu código aquí
	// 1. Crear un ticker con interval
	// 2. Crear un timer para detener después de n iteraciones
	// 3. Usar select para manejar ambos canales
	// 4. Detener el ticker cuando termines
}

func main() {
	countWithTicker(5, 1*time.Second)
}
```

**Salida esperada:**
```
1
2
3
4
5
Completado
```

---

### Ejercicio 3: Formatter - Convertir Fechas a Múltiples Formatos

**Objetivo**: Crear una función que tome una fecha y la convierta a varios formatos.

```go
// TODO: Completar la implementación
package main

import (
	"fmt"
	"time"
)

// FormatMultiple convierte `t` a varios formatos comunes
func formatMultiple(t time.Time) map[string]string {
	// Tu código aquí
	// Retornar un mapa con claves: "ISO", "US", "European", "RFC3339", "Time"
	// Usar los layouts apropiados
	return nil // Reemplazar
}

func main() {
	t := time.Date(2024, time.January, 15, 14, 30, 45, 0, time.UTC)
	
	formats := formatMultiple(t)
	for format, value := range formats {
		fmt.Printf("%s: %s\n", format, value)
	}
}
```

**Salida esperada:**
```
ISO: 2024-01-15
US: 01/15/2024
European: 15-01-2024
RFC3339: 2024-01-15T14:30:45Z
Time: 14:30:45
```

---

### Ejercicio 4: Stopwatch - Medir Tiempo de Ejecución

**Objetivo**: Crear un Stopwatch que mida el tiempo de ejecución de una función.

```go
// TODO: Completar la implementación
package main

import (
	"fmt"
	"time"
)

// Stopwatch mide el tiempo transcurrido
type Stopwatch struct {
	// Tu código aquí
}

func NewStopwatch() *Stopwatch {
	// Tu código aquí
	return nil // Reemplazar
}

func (s *Stopwatch) Stop() time.Duration {
	// Tu código aquí
	return 0 // Reemplazar
}

func (s *Stopwatch) Elapsed() time.Duration {
	// Tu código aquí
	return 0 // Reemplazar
}

func main() {
	sw := NewStopwatch()
	time.Sleep(100 * time.Millisecond)
	elapsed := sw.Stop()
	
	fmt.Printf("Tiempo: %v\n", elapsed)
	fmt.Printf("Milisegundos: %.0f\n", elapsed.Seconds()*1000)
}
```

**Salida esperada:**
```
Tiempo: 100ms
Milisegundos: 100
```

---

### Ejercicio 5: Timezone Aware - Trabajar con Múltiples Zonas Horarias

**Objetivo**: Crear una función que muestre la hora en diferentes ciudades.

```go
// TODO: Completar la implementación
package main

import (
	"fmt"
	"time"
)

// TimeInCities retorna un mapa de ciudad -> hora local
func timeInCities(t time.Time, cities map[string]string) map[string]string {
	// Tu código aquí
	// cities: {"Nueva York": "America/New_York", "Tokio": "Asia/Tokyo", ...}
	// Retornar mapa con la hora local en cada ciudad
	return nil // Reemplazar
}

func main() {
	cities := map[string]string{
		"Nueva York":   "America/New_York",
		"Londres":      "Europe/London",
		"París":        "Europe/Paris",
		"Tokio":        "Asia/Tokyo",
		"Sydney":       "Australia/Sydney",
	}
	
	utcTime := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	times := timeInCities(utcTime, cities)
	
	for city, localTime := range times {
		fmt.Printf("%-15s: %s\n", city, localTime)
	}
}
```

**Salida esperada:**
```
Nueva York     : 2024-01-15 07:00:00 EST (UTC-5)
Londres        : 2024-01-15 12:00:00 GMT (UTC+0)
París          : 2024-01-15 13:00:00 CET (UTC+1)
Tokio          : 2024-01-15 21:00:00 JST (UTC+9)
Sydney         : 2024-01-15 23:00:00 AEDT (UTC+11)
```

---

## Resumen de Conceptos Clave

### Jerarquía de Abstracciones

```
┌─────────────────────────────────────────────────┐
│      time.Time (Instantes específicos)          │
│  - Now(), Unix(), Date()                        │
│  - Format(), Parse()                            │
│  - Add(), Sub(), Before(), After()              │
└─────────────────────────────────────────────────┘
        ↑ Se construyen sobre ↓
┌─────────────────────────────────────────────────┐
│   time.Duration (Intervalos)                    │
│  - Second, Minute, Hour, etc.                   │
│  - Milliseconds(), Seconds(), Hours()           │
│  - Operaciones aritméticas                      │
└─────────────────────────────────────────────────┘
        ↑ Con ↓
┌─────────────────────────────────────────────────┐
│   time.Location (Zonas horarias)                │
│  - UTC, Local                                   │
│  - LoadLocation()                               │
│  - FixedZone()                                  │
└─────────────────────────────────────────────────┘
        ↑ Habilitan ↓
┌─────────────────────────────────────────────────┐
│   Timers & Tickers (Eventos de tiempo)          │
│  - NewTimer(), After()                          │
│  - NewTicker(), Tick()                          │
│  - Canales para concurrencia                    │
└─────────────────────────────────────────────────┘
```

### Checklist de Buenas Prácticas

- ✅ **Almacenar siempre en UTC**, mostrar en zona local
- ✅ **Usar time.Since() para medir duraciones** (reloj monotónico)
- ✅ **Especificar layout explícitamente** en Parse/Format
- ✅ **Validar errores al parsear** fechas de entrada
- ✅ **Usar context.WithTimeout** para operaciones concurrentes
- ✅ **Detener Timers y Tickers** para evitar goroutine leaks
- ✅ **Manejar DST correctamente** (Go lo hace automáticamente)
- ✅ **Usar time.After y Ticker** en select para concurrencia elegante
- ✅ **Parsear con ParseInLocation si la zona no está en el string**
- ✅ **Testear con context deadlines** para operaciones time-sensitive

---

## Conclusión

El package `time` de Go proporciona abstracciones potentes y seguras para trabajar con tiempo en aplicaciones concurrentes. Su diseño favorece la claridad sobre la brevedad, lo que resulta en código más mantenible. Puntos clave:

1. **Time.Time** representa instantes con precisión de nanosegundos
2. **Duration** maneja intervalos de forma segura
3. **Timers y Tickers** integran tiempo con goroutines vía canales
4. **Zonas horarias** se manejan explícitamente (no implícitamente)
5. **Context** proporciona timeouts estructurados para operaciones distribuidas

Dominar el `time` package es esencial para construir aplicaciones Go robustas y eficientes.

