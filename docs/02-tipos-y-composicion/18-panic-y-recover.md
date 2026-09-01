# Capítulo 18: Panic y recover

## Índice del Capítulo 18

1. [18.1 ¿Qué es Panic?](#181-qué-es-panic)
2. [18.2 Disparar Panic](#182-disparar-panic)
3. [18.3 Recuperación con Recover](#183-recuperación-con-recover)
4. [18.4 Defer y Panic](#184-defer-y-panic)
5. [18.5 Diferencia: Error vs Panic](#185-diferencia-error-vs-panic)
6. [18.6 Panic en Goroutines](#186-panic-en-goroutines)
7. [18.7 Patrones de Recuperación](#187-patrones-de-recuperación)
8. [18.8 Built-in Panics](#188-built-in-panics)
9. [18.9 Instrumentación: Capturar Stack Traces](#189-instrumentación-capturar-stack-traces)
10. [18.10 Testing Panics](#1810-testing-panics)
11. [18.11 Buenas Prácticas y Antipatrones](#1811-buenas-prácticas-y-antipatrones)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 18.1 ¿Qué es Panic?

### La Naturaleza del Panic

En Go, **panic** es un mecanismo especial para manejar situaciones **EXCEPCIONALES e IRRECUPERABLES** donde la ejecución normal no puede continuar. No es para errores normales; Go usa el tipo `error` para eso.

```
┌────────────────────────────────────────────────────────────┐
│         Manejo de Situaciones Anormales en Go             │
├────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────────┐      ┌──────────────────────┐  │
│  │   Error (error)      │      │   Panic              │  │
│  │                      │      │                      │  │
│  │  • Esperado          │      │  • Inesperado        │  │
│  │  • Recuperable       │      │  • IRRECUPERABLE     │  │
│  │  • Programador       │      │  • Programador o     │  │
│  │    decide acción     │      │    Especificación    │  │
│  │  • Ejemplo: archivo  │      │  • Ejemplo: assert   │  │
│  │    no encontrado     │      │    fallido           │  │
│  └──────────────────────┘      └──────────────────────┘  │
│                                                             │
└────────────────────────────────────────────────────────────┘
```

### Stack Unwinding

Cuando ocurre un **panic**:

1. **Se detiene la ejecución normal** de la función actual
2. **Se ejecutan todos los `defer`** en orden inverso (LIFO - Last In, First Out)
3. **Se propaga hacia arriba** en la pila de llamadas
4. **El programa termina** si no se recupera con `recover()`
5. **Se imprime un stack trace** antes de salir

```
Flujo de Ejecución Normal vs Panic:

NORMAL:
  main()
    └─ funcA() retorna normalmente
    └─ funcB() retorna normalmente
    └─ main() finaliza ✓

PANIC:
  main()
    └─ funcA()
    └─ funcB() ← panic aquí
       └─ [STOP] funcB
       └─ [defer B1] ejecuta
       └─ [defer B2] ejecuta
    └─ [defer A1] ejecuta
    └─ [defer A2] ejecuta
    └─ main() - no puede hacer nada
    └─ PROGRAMA TERMINA CON ERROR ✗
```

### Cuándo Usar Panic (Casos Excepcionales)

Panic es **SOLO** para situaciones donde NO es seguro continuar:

```go
package main

// ✓ CORRECTO: Assertion que nunca debería fallar
func calcularPromedio(calificaciones []float64) float64 {
    if len(calificaciones) == 0 {
        panic("calificaciones no puede estar vacío")
    }
    // ... continuar
}

// ✗ INCORRECTO: Usar panic para control de flujo
func buscarUsuario(id int) {
    usuario, err := db.GetUser(id)
    if err != nil {
        panic(err)  // MAL: esto es un error normal, no una excepción
    }
}

// ✓ CORRECTO: Programmer error que nunca debería ocurrir
func accederSlice(slice []int, indice int) int {
    if indice < 0 || indice >= len(slice) {
        panic(fmt.Sprintf("índice %d fuera de rango", indice))
    }
    return slice[indice]
}
```

### Panic NO es para Errores Normales

La filosofía de Go es:

> **"Errores son valores. Los panics son para lo inesperado."**

```go
// Comparación: Error vs Panic

// ✓ Error: Situación normal y esperada
file, err := os.Open("config.json")
if err != nil {
    log.Printf("no se pudo abrir config: %v", err)
    return nil
}

// ✓ Panic: Situación que NUNCA debería ocurrir
if version, err := strconv.Atoi(os.Getenv("VERSION")); err != nil {
    // Si VERSION no está configurada, es un problema de deployment
    // No tiene sentido continuar
    panic(fmt.Sprintf("VERSION env var inválida: %v", err))
}
```

### Filosofía: Por Qué Go No Usa Excepciones

A diferencia de Java, Python o C++, Go NO tiene un sistema de excepciones tradicional:

```
COMPARACIÓN CON OTROS LENGUAJES:

Java:
  try {
      conectarBaseDatos()
      procesarDatos()
  } catch (SQLException e) {
      // Cada línea podría lanzar una excepción diferente
      // Difícil saber cuál es de dónde
  }

Python:
  try:
      conectar_base_datos()
      procesar_datos()
  except DatabaseError as e:
      # Idem: difícil rastrear cuál línea falló

Go:
  if err := conectarBaseDatos(); err != nil {
      return fmt.Errorf("conectar DB: %w", err)
  }
  if err := procesarDatos(); err != nil {
      return fmt.Errorf("procesar datos: %w", err)
  }
  // Claro: sé exactamente dónde puede fallar cada línea
```

---

## 18.2 Disparar Panic

### Función `panic()`

La función built-in `panic()` toma un argumento de tipo `interface{}` (cualquier valor):

```go
package main

import "fmt"

func main() {
    // Panic con string
    panic("algo salió mal")
    
    // Panic con error
    panic(errors.New("error crítico"))
    
    // Panic con número
    panic(42)
    
    // Panic con struct personalizado
    type ErrorInfo struct {
        Mensaje string
        Código  int
    }
    panic(ErrorInfo{"error interno", 500})
}
```

### Mensaje de Panic en Diferentes Tipos

```go
package main

import "fmt"

func ejemploPanicTipos() {
    // Con string: mensaje claro
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recuperado: %v (tipo: %T)\n", r, r)
        }
    }()
    
    switch r := recover(); r {
    case nil:
        fmt.Println("Sin panic")
    case 1:
        panic("error string")
    case 2:
        panic(fmt.Errorf("error con formato"))
    case 3:
        panic(struct{msg string}{msg: "struct"})
    }
}

// Resultado del stack trace:
// panic: error string
// 
// goroutine 1 [running]:
// main.ejemploPanicTipos()
//     /ruta/al/archivo.go:25
// main.main()
//     /ruta/al/archivo.go:42
```

### Stack Trace Automático

Cuando panic ocurre, Go imprime automáticamente:

```
panic: índice fuera de rango

goroutine 1 [running]:
main.accederSlice({0xc00001a280, 0x3, 0x3}, -1)
    /home/usuario/proyecto/main.go:15 +0x44
main.main()
    /home/usuario/proyecto/main.go:30 +0x64
exit status 2
```

Información del stack trace:

```go
// Línea de panic
panic: índice fuera de rango

// Información de goroutines
goroutine 1 [running]  // ID y estado

// Stack de llamadas (más reciente primero)
main.accederSlice()           // Función donde ocurrió
    /ruta/archivo.go:15       // Archivo y línea
    +0x44                     // Offset en el código compilado

main.main()                   // Función que llamó
    /ruta/archivo.go:30
    +0x64

// Código de salida
exit status 2
```

### Crear Panics Informativos

```go
package main

import (
    "fmt"
    "runtime"
)

// Panic con contexto
func operacionCritica(config map[string]interface{}) {
    if config == nil {
        pc, file, line, _ := runtime.Caller(0)
        fn := runtime.FuncForPC(pc)
        panic(fmt.Sprintf(
            "%s [%s:%d] config no puede ser nil",
            fn.Name(), file, line,
        ))
    }
}

// Panic con información estructurada
type PanicInfo struct {
    Mensaje    string
    Función    string
    Archivo    string
    Línea      int
    Contexto   map[string]interface{}
}

func panicEstructurado(msg string, contexto map[string]interface{}) {
    pc, file, line, _ := runtime.Caller(1)
    fn := runtime.FuncForPC(pc)
    
    info := PanicInfo{
        Mensaje:  msg,
        Función:  fn.Name(),
        Archivo:  file,
        Línea:    line,
        Contexto: contexto,
    }
    
    panic(info)
}

// Uso
func procesar(datos interface{}) {
    if datos == nil {
        panicEstructurado("datos requerido", map[string]interface{}{
            "llamante": "procesar",
            "tipo":     "validación",
        })
    }
}
```

---

## 18.3 Recuperación con Recover

### Función `recover()`

`recover()` es la ÚNICA forma de capturar un panic. Características clave:

1. **Solo funciona dentro de `defer`**
2. **Retorna el valor pasado a `panic()`**
3. **Retorna `nil` si no hay panic**
4. **Detiene el stack unwinding**

```go
package main

import "fmt"

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("¡Panic capturado:", r)
        }
    }()
    
    panic("algo salió mal")
    // Línea nunca ejecutada sin defer/recover
}

// Output:
// ¡Panic capturado: algo salió mal
```

### Patrón de Recover Seguro

```go
package main

import (
    "fmt"
    "log"
)

// Patrón estándar para capturar panic
func ejecutarConRecuperación(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            // Convertir panic a error
            switch x := r.(type) {
            case string:
                err = fmt.Errorf("panic: %s", x)
            case error:
                err = fmt.Errorf("panic: %w", x)
            default:
                err = fmt.Errorf("panic: %v", x)
            }
            
            log.Printf("Recuperado de panic: %v", err)
        }
    }()
    
    fn()
    return nil
}

// Uso
func operacionRiesgosa() {
    panic("error inesperado")
}

func main() {
    if err := ejecutarConRecuperación(operacionRiesgosa); err != nil {
        fmt.Printf("Error controlado: %v\n", err)
    }
}
```

### Tipos de Valores en Recover

```go
package main

import "fmt"

func demostrarRecover() {
    defer func() {
        if r := recover(); r != nil {
            // Type assertion para manejar diferentes tipos
            switch v := r.(type) {
            case string:
                fmt.Printf("String panic: %s\n", v)
            case int:
                fmt.Printf("Int panic: %d\n", v)
            case error:
                fmt.Printf("Error panic: %v\n", v)
            case float64:
                fmt.Printf("Float panic: %f\n", v)
            default:
                fmt.Printf("Unknown panic type: %T, valor: %v\n", r, r)
            }
        }
    }()
    
    // Probar diferentes tipos
    // panic("string error")
    // panic(42)
    // panic(errors.New("error object"))
    // panic(3.14)
}
```

### Recovering en Múltiples Niveles

```go
package main

import "fmt"

func funcC() {
    fmt.Println("funcC: inicio")
    defer func() {
        fmt.Println("funcC: defer (recuperará el panic)")
        if r := recover(); r != nil {
            fmt.Printf("funcC: Recuperado panic: %v\n", r)
            // NO re-lanzar: El programa continúa
        }
    }()
    
    panic("panic desde funcC")
    fmt.Println("funcC: línea no ejecutada")
}

func funcB() {
    fmt.Println("funcB: inicio")
    defer func() {
        fmt.Println("funcB: defer (no hay panic)")
    }()
    
    funcC()
    fmt.Println("funcB: continuando después de funcC")
}

func funcA() {
    fmt.Println("funcA: inicio")
    defer func() {
        fmt.Println("funcA: defer")
    }()
    
    funcB()
    fmt.Println("funcA: fin")
}

func main() {
    funcA()
    fmt.Println("main: programa continuó normalmente")
}

// Output:
// funcA: inicio
// funcB: inicio
// funcC: inicio
// funcC: defer (recuperará el panic)
// funcC: Recuperado panic: panic desde funcC
// funcB: continuando después de funcC
// funcA: fin
// funcA: defer
// main: programa continuó normalmente
```

---

## 18.4 Defer y Panic

### Orden de Ejecución de Defer

Los `defer` se ejecutan en orden **LIFO** (Last In, First Out) incluso cuando hay panic:

```
Sin Panic:
  defer fmt.Println("1")
  defer fmt.Println("2")
  defer fmt.Println("3")
  // Output: 3, 2, 1

Con Panic (mismo orden):
  defer fmt.Println("1")
  defer fmt.Println("2")
  defer fmt.Println("3")
  panic("error")
  // Output: 3, 2, 1 (luego crash si no hay recover)
```

### Demostración Práctica

```go
package main

import "fmt"

func demoDeferPanic() {
    fmt.Println("Inicio")
    
    defer func() {
        fmt.Println("defer 1 - siempre se ejecuta")
    }()
    
    defer func() {
        fmt.Println("defer 2 - siempre se ejecuta")
    }()
    
    defer func() {
        fmt.Println("defer 3 - último registrado")
        if r := recover(); r != nil {
            fmt.Printf("  Recuperado: %v\n", r)
        }
    }()
    
    fmt.Println("Antes del panic")
    panic("¡Algo salió mal!")
    fmt.Println("Esta línea NO se ejecuta")
}

func main() {
    demoDeferPanic()
    fmt.Println("Programa continúa después del panic")
}

// Output:
// Inicio
// Antes del panic
// defer 3 - último registrado
//   Recuperado: ¡Algo salió mal!
// defer 2 - siempre se ejecuta
// defer 1 - siempre se ejecuta
// Programa continúa después del panic
```

### Defer SIEMPRE se Ejecuta

Esto es crítico: los defer se ejecutan incluso si:

- Ocurre panic
- Hay un return inesperado
- Hay break/continue
- La función termina normalmente

```go
package main

import "fmt"

func garantizarEjecución() string {
    defer func() {
        fmt.Println("defer: limpieza siempre ocurre")
    }()
    
    // Escenario 1: Return normal
    // return "valor"
    
    // Escenario 2: Panic
    panic("error")
    
    // Escenario 3: Nunca alcanzado
    return "nunca"
}

// defer se ejecuta en TODOS los casos
```

### Patrón: Cleanup Garantizado

```go
package main

import (
    "fmt"
    "os"
)

// Manejo de recursos con defer
func procesarArchivo(ruta string) error {
    archivo, err := os.Open(ruta)
    if err != nil {
        return err
    }
    // ✓ CORRECTO: defer asegura cierre
    defer archivo.Close()
    
    defer func() {
        fmt.Println("[LOG] Finalizando procesamiento")
    }()
    
    // Incluso si panic ocurre aquí:
    // 1. Se ejecuta defer de logging
    // 2. Se ejecuta defer de Close
    // 3. Error se propaga
    
    return nil
}

// Sin defer (peligro):
func procesarArchivoMal(ruta string) error {
    archivo, err := os.Open(ruta)
    if err != nil {
        return err
    }
    
    // ¿Qué pasa si panic aquí? El archivo NO se cierra
    panic("error durante lectura")
    
    archivo.Close()  // NUNCA SE EJECUTA
    return nil
}
```

### Patrones de Defer Efectivos

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Patrón 1: Medición de tiempo
func medirTiempo(nombre string) {
    inicio := time.Now()
    defer func() {
        duracion := time.Since(inicio)
        fmt.Printf("[%s] Duración: %v\n", nombre, duracion)
    }()
    
    // Operación que puede tomar tiempo
    time.Sleep(100 * time.Millisecond)
}

// Patrón 2: Locks garantizados
func operacionConcurrente(mu *sync.Mutex, dato *int) {
    mu.Lock()
    defer mu.Unlock()  // Siempre se libera
    
    // Incluso si panic aquí, el mutex se libera
    *dato++
    if *dato > 10 {
        panic("valor muy alto")
    }
}

// Patrón 3: Transacciones
func procesarTransaccion(db interface{}) error {
    // Iniciar transacción (imaginario)
    // tx := db.Begin()
    
    defer func() {
        if r := recover(); r != nil {
            // tx.Rollback()  // Revertir si panic
            fmt.Println("Transacción revertida por panic")
        } else if err := recover(); err != nil {
            // tx.Rollback()  // Revertir si error
        }
        // tx.Commit()  // Confirmar si todo OK
    }()
    
    // Operaciones críticas
    return nil
}
```

---

## 18.5 Diferencia: Error vs Panic

### Comparación Teórica

```
┌─────────────────────┬───────────────────┬──────────────────┐
│ Aspecto              │ Error             │ Panic            │
├─────────────────────┼───────────────────┼──────────────────┤
│ Naturaleza           │ Esperado          │ Inesperado       │
│ Control de flujo     │ Programador       │ Automático       │
│ Recuperación         │ Sí (manual)       │ Sí (automática)  │
│ Performance          │ Sin costo         │ Stack unwinding  │
│ Uso normal           │ Sí, siempre       │ SOLO excepciones │
│ Goroutine crash      │ No                │ Sí               │
└─────────────────────┴───────────────────┴──────────────────┘
```

### Ejemplos de Decisión: Error vs Panic

```go
package main

import (
    "fmt"
    "os"
)

// ✓ ERROR: Archivo no encontrado
func leerConfiguracion(ruta string) (map[string]string, error) {
    contenido, err := os.ReadFile(ruta)
    if err != nil {
        return nil, fmt.Errorf("no se pudo leer config: %w", err)
    }
    // ...
    return nil, nil
}

// ✓ PANIC: Programador olvidó inicializar
func procesarDatos(datos interface{}) {
    if datos == nil {
        panic("datos no puede ser nil - fallo de programador")
    }
}

// ✓ ERROR: Usuario proporciona entrada inválida
func parsearNumero(entrada string) (int, error) {
    // Usar strconv, que retorna error
    return strconv.Atoi(entrada)
}

// ✓ PANIC: Especificación del sistema violada
func asignarMemoria(tamaño int) {
    if tamaño < 0 {
        panic(fmt.Sprintf("tamaño no puede ser negativo: %d", tamaño))
    }
}

// ✓ ERROR: Operación normal que puede fallar
func conectarBaseDatos(host string, puerto int) (interface{}, error) {
    // La BD podría no estar disponible
    return nil, fmt.Errorf("no se pudo conectar a %s:%d", host, puerto)
}

// ✓ PANIC: Invariante de código violada
func accederArray(arr []int, indice int) int {
    if indice < 0 || indice >= len(arr) {
        panic("índice fuera de rango")
    }
    return arr[indice]
}
```

### Decisión de Diseño

La pregunta clave:

```
¿Es seguro CONTINUAR si esto falla?

       SÍ → Usar error
       │
       NO → ¿Es algo que el usuario/programa podría hacer?
           │
           SÍ → Considerar error también
           NO → Usar panic (programmer error)
```

Ejemplos de decisión:

```go
// Lectura de archivo de configuración
- Usuario olvida crear config.json → Error
- Sistema de archivos falla → Error
- Parser de JSON se corrompe → Error (o panic si es bug interno)

// Acceso a índice de array
- Índice fuera de rango (acceso directo) → Panic
- Búsqueda binaria no encuentra elemento → Error
- Tamaño del array es negativo → Panic

// Conexión a servicio
- Servicio no disponible → Error
- Socket corrupto (nunca debería pasar) → Panic
- Host inválido en especificación → Panic
```

### Interacción Error/Panic

```go
package main

import (
    "fmt"
)

// Función que puede fallar normalmente
func operacionNormal() error {
    return fmt.Errorf("operación falló")
}

// Función que envuelve operación normal
func wrapperConRecuperacion() (resultado interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            resultado = nil
            err = fmt.Errorf("panic en wrapper: %v", r)
        }
    }()
    
    // Esto es un error normal, NO un panic
    if err := operacionNormal(); err != nil {
        return nil, err  // Retornar error normalmente
    }
    
    // Si ocurre panic aquí, será capturado
    panic("algo excepcional")
    
    return "ok", nil
}

func main() {
    _, err := wrapperConRecuperacion()
    if err != nil {
        fmt.Printf("Error capturado: %v\n", err)
    }
}
```

---

## 18.6 Panic en Goroutines

### El Problema: Panic Mata la Goroutine

El comportamiento crítico de panic en goroutines:

```
┌─────────────────────────────────────────────────────────┐
│        Panic en Goroutines vs Main                      │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Panic en main()                                        │
│  └─ PROGRAMA TERMINA COMPLETAMENTE ✗                   │
│                                                          │
│  Panic en goroutine                                     │
│  └─ SOLO esa goroutine termina ✗                        │
│  └─ Otras goroutines continúan (¡sin supervisión!)     │
│  └─ Programa podría parecer funcionar                   │
│     pero está parcialmente muerto                       │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Demostración del Problema

```go
package main

import (
    "fmt"
    "time"
)

func goroutineConPanic() {
    for i := 1; i <= 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("Goroutine: %d\n", i)
        
        if i == 3 {
            panic("¡Error en goroutine!")
            // ← Panic aquí mata esta goroutine
        }
    }
}

func main() {
    fmt.Println("Main: inicio")
    
    go goroutineConPanic()
    
    // Main continúa ejecutándose
    for i := 1; i <= 5; i++ {
        time.Sleep(200 * time.Millisecond)
        fmt.Printf("Main: %d\n", i)
    }
    
    fmt.Println("Main: fin")
}

// Output:
// Main: inicio
// Goroutine: 1
// Main: 1
// Goroutine: 2
// Goroutine: 3
// panic: ¡Error en goroutine!  ← Solo la goroutine muere
// 
// Main: 2
// Main: 3
// Main: 4
// Main: 5
// Main: fin
```

### Aislamiento de Panic

**Las goroutines son aisladas en términos de panic:**

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    
    // Goroutine 1: con panic
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("G1: inicio")
        panic("error en G1")
        fmt.Println("G1: nunca se ejecuta")
    }()
    
    // Goroutine 2: sin afectar
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("G2: inicio")
        fmt.Println("G2: completada")
    }()
    
    wg.Wait()
    fmt.Println("Main: fin")
}

// Output:
// G1: inicio
// G2: inicio
// G2: completada
// panic: error en G1
// 
// [Program exits with panic]
```

### Supervisor de Goroutines

La solución: **recuperar panic en cada goroutine**:

```go
package main

import (
    "fmt"
    "log"
    "sync"
)

// Función wrapper que maneja panic
func ejecutarConSeguridad(nombre string, fn func()) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[%s] PANIC: %v", nombre, r)
            // Aquí se podría notificar, logging, etc.
        }
    }()
    
    fn()
}

func main() {
    var wg sync.WaitGroup
    
    // Goroutine con panic (supervisada)
    wg.Add(1)
    go func() {
        defer wg.Done()
        ejecutarConSeguridad("worker-1", func() {
            fmt.Println("Worker 1: procesando...")
            panic("algo salió mal")
        })
    }()
    
    // Goroutine normal (supervisada)
    wg.Add(1)
    go func() {
        defer wg.Done()
        ejecutarConSeguridad("worker-2", func() {
            fmt.Println("Worker 2: procesando...")
            fmt.Println("Worker 2: completado")
        })
    }()
    
    wg.Wait()
    fmt.Println("Main: todos los workers finalizaron")
}

// Output:
// Worker 1: procesando...
// Worker 2: procesando...
// Worker 2: completado
// [PANIC] worker-1: algo salió mal
// Main: todos los workers finalizaron ✓
```

### Pool de Goroutines Supervisadas

```go
package main

import (
    "fmt"
    "log"
    "sync"
    "time"
)

type WorkerPool struct {
    jobs    chan func()
    errCh   chan error
    wg      sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    return &WorkerPool{
        jobs:  make(chan func()),
        errCh: make(chan error, numWorkers),
    }
}

func (wp *WorkerPool) Start(numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for job := range wp.jobs {
        func() {
            defer func() {
                if r := recover(); r != nil {
                    err := fmt.Errorf("worker %d panic: %v", id, r)
                    log.Printf("ERROR: %v", err)
                    wp.errCh <- err
                }
            }()
            
            job()
        }()
    }
}

func (wp *WorkerPool) Submit(job func()) {
    wp.jobs <- job
}

func (wp *WorkerPool) Stop() {
    close(wp.jobs)
    wp.wg.Wait()
}

// Uso
func main() {
    pool := NewWorkerPool(3)
    pool.Start(3)
    
    // Enviar trabajos
    for i := 1; i <= 5; i++ {
        idx := i
        pool.Submit(func() {
            fmt.Printf("Job %d: ejecutando\n", idx)
            if idx == 3 {
                panic("fallo en job 3")
            }
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("Job %d: completado\n", idx)
        })
    }
    
    pool.Stop()
    fmt.Println("Pool: finalizado")
}
```

---

## 18.7 Patrones de Recuperación

### Patrón 1: Convertir Panic a Error

```go
package main

import (
    "fmt"
)

// Función que podría hacer panic
func puedeHacerPanic(shouldPanic bool) string {
    if shouldPanic {
        panic("algo salió mal")
    }
    return "éxito"
}

// Wrapper que convierte panic a error
func puedeHacerPanicSafe(shouldPanic bool) (resultado string, err error) {
    defer func() {
        if r := recover(); r != nil {
            resultado = ""
            switch x := r.(type) {
            case string:
                err = fmt.Errorf("panic: %s", x)
            case error:
                err = fmt.Errorf("panic: %w", x)
            default:
                err = fmt.Errorf("panic: %v", x)
            }
        }
    }()
    
    resultado = puedeHacerPanic(shouldPanic)
    return resultado, err
}

func main() {
    // Conversión a error
    resultado, err := puedeHacerPanicSafe(true)
    if err != nil {
        fmt.Printf("Error capturado: %v\n", err)
    }
}
```

### Patrón 2: Error Wrapping de Panics

```go
package main

import (
    "fmt"
    "runtime"
)

type PanicError struct {
    Mensaje string
    Stack   string
    Valor   interface{}
}

func (pe *PanicError) Error() string {
    return fmt.Sprintf("panic: %v\nstack:\n%s", pe.Valor, pe.Stack)
}

// Capturar panic con información completa
func capturaConContext(fn func()) error {
    defer func() {
        if r := recover(); r != nil {
            // Capturar stack trace
            buf := make([]byte, 4096)
            n := runtime.Stack(buf, false)
            
            perr := &PanicError{
                Valor: r,
                Stack: string(buf[:n]),
            }
            
            return  // Retornar error wrapping el panic
        }
    }()
    
    fn()
    return nil
}

func main() {
    err := capturaConContext(func() {
        panic("error crítico")
    })
    
    if err != nil {
        fmt.Printf("Error capturado:\n%v\n", err)
    }
}
```

### Patrón 3: Logging de Panics

```go
package main

import (
    "fmt"
    "log"
    "os"
)

// Logger que captura panics
func ejecutarConLogging(operacion func(), nombre string) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf(
                "[CRITICAL] Panic en %s: %v\nStack: %+v",
                nombre, r, r,
            )
            
            // Guardar en archivo de log
            file, _ := os.OpenFile(
                "panic.log",
                os.O_APPEND|os.O_CREATE|os.O_WRONLY,
                0644,
            )
            defer file.Close()
            
            fmt.Fprintf(
                file,
                "[%s] Panic: %v\n",
                nombre, r,
            )
        }
    }()
    
    operacion()
}

func main() {
    ejecutarConLogging(func() {
        fmt.Println("Operación: inicio")
        panic("algo inesperado")
    }, "procesamiento-datos")
}
```

### Patrón 4: Graceful Degradation

```go
package main

import (
    "fmt"
    "log"
)

type Servicio interface {
    Ejecutar() (interface{}, error)
}

// Servicio primario que podría fallar
type ServicioPrimario struct{}

func (s *ServicioPrimario) Ejecutar() (interface{}, error) {
    panic("servicio primario no disponible")
}

// Servicio de respaldo
type ServicioRespaldo struct{}

func (s *ServicioRespaldo) Ejecutar() (interface{}, error) {
    return "respuesta del respaldo", nil
}

// Ejecutar con fallback
func ejecutarConFallback(primario, respaldo Servicio) interface{} {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Primario falló: %v, usando respaldo", r)
        }
    }()
    
    resultado, _ := primario.Ejecutar()
    return resultado
}

func main() {
    primario := &ServicioPrimario{}
    respaldo := &ServicioRespaldo{}
    
    // Ejecutar primario
    resultado := ejecutarConFallback(primario, respaldo)
    fmt.Printf("Resultado: %v\n", resultado)
}
```

### Patrón 5: Re-lanzar Panic Selectivamente

```go
package main

import (
    "fmt"
)

// Re-lanzar solo ciertos panics
func filtroDeRecuperacion(fn func()) error {
    defer func() {
        if r := recover(); r != nil {
            if val, ok := r.(string); ok && val == "erro_critico" {
                // Re-lanzar errores críticos
                panic(r)
            }
            
            // Otros panics: ignorar y continuar
            fmt.Printf("Panic ignorado: %v\n", r)
        }
    }()
    
    fn()
    return nil
}

func main() {
    // Panic que será ignorado
    filtroDeRecuperacion(func() {
        panic("error menor")
    })
    
    fmt.Println("Main: continuando después de error menor")
    
    // Panic que será re-lanzado
    filtroDeRecuperacion(func() {
        panic("erro_critico")  // Esto TERMINARÁ el programa
    })
}
```

---

## 18.8 Built-in Panics

### Panics Automáticos en Go

Go lanza panic automáticamente en situaciones muy específicas:

```
┌──────────────────────────┬──────────────────────────────┐
│ Condición                │ Panic                        │
├──────────────────────────┼──────────────────────────────┤
│ Index out of range       │ "index X out of range [Y]"  │
│ Nil pointer dereference  │ "runtime error: invalid     │
│                          │ memory address"             │
│ Slice append overflow    │ "slice bounds out of range" │
│ Divide by zero           │ "runtime error: integer     │
│                          │ divide by zero"             │
│ Nil map assignment       │ "assignment to entry in nil │
│                          │ map"                        │
│ Channel operations       │ "send on closed channel"    │
│ Type assertion fail      │ "interface conversion error"│
│ Stack overflow           │ "stack overflow"            │
└──────────────────────────┴──────────────────────────────┘
```

### Index Out of Range

```go
package main

import "fmt"

func main() {
    // Panic: index out of range
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recuperado: %v\n", r)
        }
    }()
    
    arr := []int{1, 2, 3}
    
    fmt.Println(arr[0])  // OK: 1
    fmt.Println(arr[2])  // OK: 3
    fmt.Println(arr[3])  // PANIC: index 3 out of range [3]
}
```

### Nil Pointer Dereference

```go
package main

import "fmt"

type Usuario struct {
    Nombre string
    Edad   int
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recuperado: %v\n", r)
        }
    }()
    
    var u *Usuario  // nil pointer
    
    fmt.Println(u.Nombre)  // PANIC: runtime error: invalid memory address
}
```

### Divide by Zero

```go
package main

import "fmt"

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recuperado: %v\n", r)
        }
    }()
    
    dividendo := 10
    divisor := 0
    
    resultado := dividendo / divisor  // PANIC: divide by zero
    fmt.Println(resultado)
}
```

### Nil Map Assignment

```go
package main

import "fmt"

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recuperado: %v\n", r)
        }
    }()
    
    var m map[string]int  // nil map
    m["clave"] = 42       // PANIC: assignment to entry in nil map
}
```

### Send on Closed Channel

```go
package main

import "fmt"

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recuperado: %v\n", r)
        }
    }()
    
    ch := make(chan int)
    close(ch)
    
    ch <- 42  // PANIC: send on closed channel
}
```

### Type Assertion Panic

```go
package main

import "fmt"

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recuperado: %v\n", r)
        }
    }()
    
    var x interface{} = "string"
    
    // Panic en type assertion (sin ok)
    valor := x.(int)  // PANIC: interface conversion error
    fmt.Println(valor)
}
```

### Prevenir Built-in Panics

```go
package main

import "fmt"

// ✓ CORRECTO: Checks antes de acceso peligroso

func accesoSeguro(arr []int, indice int) (int, error) {
    if indice < 0 || indice >= len(arr) {
        return 0, fmt.Errorf("índice fuera de rango")
    }
    return arr[indice], nil
}

func dereferenciaSafe(ptr *int) (int, error) {
    if ptr == nil {
        return 0, fmt.Errorf("puntero es nil")
    }
    return *ptr, nil
}

func divisionSafe(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("división por cero")
    }
    return a / b, nil
}

func main() {
    // Uso seguro
    resultado, err := accesoSeguro([]int{1, 2, 3}, 10)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

---

## 18.9 Instrumentación: Capturar Stack Traces

### Función runtime.Stack

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            // Capturar stack trace completo
            buf := make([]byte, 4096)
            n := runtime.Stack(buf, false)
            stack := string(buf[:n])
            
            fmt.Printf("Panic: %v\n\n", r)
            fmt.Printf("Stack Trace:\n%s\n", stack)
        }
    }()
    
    funcionA()
}

func funcionA() {
    funcionB()
}

func funcionB() {
    funcionC()
}

func funcionC() {
    panic("error en funcionC")
}

// Output:
// Panic: error en funcionC
//
// Stack Trace:
// goroutine 1 [running]:
// main.main.func1()
//     /ruta/archivo.go:13 +0x44
// main.main()
//     /ruta/archivo.go:20 +0x84
// runtime.main()
//     /usr/lib/go/src/runtime/proc.go:250 +0x220
```

### Captura Selectiva de Frames

```go
package main

import (
    "fmt"
    "runtime"
)

// Obtener información de la llamada actual
func informacionLlamada() (string, string, int) {
    pc, archivo, línea, ok := runtime.Caller(1)
    if !ok {
        return "", "", 0
    }
    
    fn := runtime.FuncForPC(pc)
    return fn.Name(), archivo, línea
}

func demostrar() {
    nombreFunc, archivo, línea := informacionLlamada()
    fmt.Printf("Función: %s\n", nombreFunc)
    fmt.Printf("Archivo: %s\n", archivo)
    fmt.Printf("Línea: %d\n", línea)
}

func main() {
    demostrar()
}
```

### Captura de Stack Trace Estructurado

```go
package main

import (
    "fmt"
    "runtime"
)

type StackFrame struct {
    Función string
    Archivo string
    Línea   int
}

// Capturar stack trace como estructura
func capturaStack() []StackFrame {
    var frames []StackFrame
    pcs := make([]uintptr, 32)
    
    // Obtener PC (Program Counter) de todas las funciones
    n := runtime.Callers(1, pcs)
    
    for _, pc := range pcs[:n] {
        fn := runtime.FuncForPC(pc)
        if fn == nil {
            continue
        }
        
        archivo, línea := fn.FileLine(pc)
        frames = append(frames, StackFrame{
            Función: fn.Name(),
            Archivo: archivo,
            Línea:   línea,
        })
    }
    
    return frames
}

func funcionA() {
    funcionB()
}

func funcionB() {
    frames := capturaStack()
    for i, frame := range frames {
        fmt.Printf("%d: %s (%s:%d)\n", i, frame.Función, frame.Archivo, frame.Línea)
    }
}

func main() {
    funcionA()
}
```

### Logging Decorado con Stack

```go
package main

import (
    "fmt"
    "log"
    "runtime"
)

// Logger que incluye información de llamada
type StackLogger struct {
    logger *log.Logger
}

func NewStackLogger() *StackLogger {
    return &StackLogger{
        logger: log.New(nil, "", log.LstdFlags),
    }
}

func (sl *StackLogger) LogPanic(mensaje string) {
    pc, archivo, línea, _ := runtime.Caller(1)
    fn := runtime.FuncForPC(pc)
    
    fmt.Printf("[PANIC] %s\n", mensaje)
    fmt.Printf("  Función: %s\n", fn.Name())
    fmt.Printf("  Ubicación: %s:%d\n", archivo, línea)
    
    // Stack completo
    buf := make([]byte, 2048)
    n := runtime.Stack(buf, false)
    fmt.Printf("  Stack:\n%s\n", string(buf[:n]))
}

func main() {
    logger := NewStackLogger()
    
    defer func() {
        if r := recover(); r != nil {
            logger.LogPanic(fmt.Sprintf("%v", r))
        }
    }()
    
    panic("error crítico")
}
```

---

## 18.10 Testing Panics

### Verificar Panics en Tests

```go
package main

import (
    "testing"
)

// Función que hace panic
func funcionQueHacePanic() {
    panic("error inesperado")
}

// Función que NO hace panic
func funcionNormal() string {
    return "ok"
}

// Test: Verificar que sí hace panic
func TestPanicOcurre(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            // Panic fue capturado: test pasa ✓
            t.Logf("Panic capturado como se esperaba: %v", r)
        } else {
            // No hubo panic: test falla ✗
            t.Error("Se esperaba panic pero no ocurrió")
        }
    }()
    
    funcionQueHacePanic()
    
    // Si llegamos aquí sin panic: falló
    t.Error("Esta línea no debería ejecutarse si hubo panic")
}

// Test: Verificar que NO hace panic
func TestNoPanic(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("Unexpected panic: %v", r)
        }
    }()
    
    resultado := funcionNormal()
    if resultado != "ok" {
        t.Errorf("Esperado 'ok', obtuve '%s'", resultado)
    }
}
```

### Tabla de Tests con Panics

```go
package main

import "testing"

func procesarValor(val int) int {
    if val < 0 {
        panic("valor negativo")
    }
    return val * 2
}

func TestProcesarValor(t *testing.T) {
    tests := []struct {
        nombre       string
        entrada      int
        debeHacerPan bool
    }{
        {"valor positivo", 5, false},
        {"valor cero", 0, false},
        {"valor negativo", -1, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.nombre, func(t *testing.T) {
            defer func() {
                r := recover()
                
                if tt.debeHacerPan && r == nil {
                    t.Error("Se esperaba panic")
                }
                if !tt.debeHacerPan && r != nil {
                    t.Errorf("Panic no esperado: %v", r)
                }
            }()
            
            procesarValor(tt.entrada)
        })
    }
}
```

### Subtests con Manejo de Panic

```go
package main

import "testing"

func TestConPanicProteccion(t *testing.T) {
    casos := []struct {
        nombre string
        datos  interface{}
        panic  bool
    }{
        {"datos válidos", "ok", false},
        {"datos inválidos", nil, true},
    }
    
    for _, caso := range casos {
        t.Run(caso.nombre, func(t *testing.T) {
            executarConProteccion := func() (resultado interface{}, panicked bool) {
                defer func() {
                    if r := recover(); r != nil {
                        panicked = true
                    }
                }()
                
                if caso.datos == nil {
                    panic("datos es nil")
                }
                resultado = caso.datos
                return
            }
            
            resultado, panicked := executarConProteccion()
            
            if caso.panic != panicked {
                t.Errorf("Esperado panic=%v, obtuve %v", caso.panic, panicked)
            }
            
            if !caso.panic && resultado != caso.datos {
                t.Errorf("Resultado incorrecto")
            }
        })
    }
}
```

### Funciones de Testing Helper

```go
package main

import "testing"

// Helper: Assert que una función hace panic
func AssertPanic(t *testing.T, fn func(), mensaje string) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("%s: Se esperaba panic pero no ocurrió", mensaje)
        }
    }()
    
    fn()
    t.Errorf("%s: Esta línea no debería ejecutarse", mensaje)
}

// Helper: Assert que una función NO hace panic
func AssertNoPanic(t *testing.T, fn func(), mensaje string) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("%s: Panic no esperado: %v", mensaje, r)
        }
    }()
    
    fn()
}

// Ejemplo de uso
func TestUsandoHelpers(t *testing.T) {
    AssertPanic(t, func() {
        panic("error")
    }, "debe hacer panic")
    
    AssertNoPanic(t, func() {
        x := 1 + 1
    }, "no debe hacer panic")
}
```

---

## 18.11 Buenas Prácticas y Antipatrones

### ✓ BUENAS PRÁCTICAS

#### 1. Usar Panic SOLO para lo Excepcional

```go
// ✓ CORRECTO
func parsearConfiguracion(json []byte) (Config, error) {
    // Error normal: archivo inválido
    return Config{}, fmt.Errorf("JSON inválido: %w", err)
}

// ✓ CORRECTO
func inicializarSistema() {
    if os.Getenv("API_KEY") == "" {
        // Programmer error: configuración incompleta
        panic("API_KEY requerida pero no configurada")
    }
}
```

#### 2. Capturar Panic en Puntos de Entrada

```go
// ✓ CORRECTO: Capturar en goroutines
func main() {
    go ejecutarConRecuperacion(miTarea)
}

func ejecutarConRecuperacion(fn func()) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Tarea falló: %v", r)
        }
    }()
    
    fn()
}
```

#### 3. Defer para Cleanup Garantizado

```go
// ✓ CORRECTO: Defer garantiza limpieza
func procesarArchivo(ruta string) error {
    file, err := os.Open(ruta)
    if err != nil {
        return err
    }
    defer file.Close()  // Se ejecuta SIEMPRE
    
    // Incluso si ocurre panic aquí:
    datos := procesarDatos(file)
    
    return guardar(datos)
}
```

#### 4. Información Contextual en Panic

```go
// ✓ CORRECTO: Panic informativo
func validar(entrada interface{}) {
    if entrada == nil {
        pc, file, line, _ := runtime.Caller(1)
        fn := runtime.FuncForPC(pc)
        panic(fmt.Sprintf(
            "%s [%s:%d] entrada no puede ser nil",
            fn.Name(), file, line,
        ))
    }
}
```

#### 5. Distinguir Error vs Panic Claramente

```go
// ✓ CORRECTO: Distinción clara
func conectar(host string, puerto int) (Conn, error) {
    if puerto < 1 || puerto > 65535 {
        panic(fmt.Sprintf("puerto inválido: %d", puerto))  // Programmer error
    }
    
    if err := dial(host, puerto); err != nil {
        return nil, err  // Error normal
    }
    
    return conn, nil
}
```

### ✗ ANTIPATRONES

#### 1. Usar Panic para Control de Flujo

```go
// ✗ INCORRECTO: Panic como control de flujo
func buscar(items []string, objetivo string) int {
    for i, item := range items {
        if item == objetivo {
            return i
        }
    }
    panic("no encontrado")  // MAL: debería retornar error o valor especial
}

// ✓ CORRECTO
func buscar(items []string, objetivo string) (int, error) {
    for i, item := range items {
        if item == objetivo {
            return i, nil
        }
    }
    return -1, errors.New("no encontrado")
}
```

#### 2. Swallow Panics Silenciosamente

```go
// ✗ INCORRECTO: Ignorar panics
func operacion() {
    defer func() {
        recover()  // Ocultar problema sin logging
    }()
    
    alguna_operacion_critica()
}

// ✓ CORRECTO: Registrar y manejar
func operacion() error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("PANIC: %v", r)  // Registrar siempre
        }
    }()
    
    alguna_operacion_critica()
    return nil
}
```

#### 3. Panics en Goroutines sin Supervisión

```go
// ✗ INCORRECTO: Sin supervisión
go func() {
    tareaRiesgosa()  // Si hace panic, la goroutine muere silenciosamente
}()

// ✓ CORRECTO: Con supervisión
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Goroutine panic: %v", r)
        }
    }()
    tareaRiesgosa()
}()
```

#### 4. Panic sin Información

```go
// ✗ INCORRECTO: Vago
panic("error")

// ✓ CORRECTO: Informativo
panic(fmt.Sprintf("invariante violada: x=%d debe ser >0", x))
```

#### 5. No Usar Defer para Recursos

```go
// ✗ INCORRECTO: Riesgo de leak
file, _ := os.Open("datos.txt")
procesar(file)
file.Close()  // Si panic antes, NO se ejecuta

// ✓ CORRECTO: Defer garantiza cierre
file, _ := os.Open("datos.txt")
defer file.Close()  // Se ejecuta SIEMPRE
procesar(file)
```

### Matriz de Decisión

```
┌────────────────────────────────────┬──────────────┬──────────────┐
│ Situación                           │ Usar Error   │ Usar Panic   │
├────────────────────────────────────┼──────────────┼──────────────┤
│ Archivo no encontrado               │ ✓            │              │
│ Conexión a BD rechazada             │ ✓            │              │
│ Parámetro con rango inválido        │              │ ✓            │
│ Nil pointer en API interna          │              │ ✓            │
│ Índice fuera de rango en array      │              │ ✓            │
│ JSON inválido de usuario            │ ✓            │              │
│ Especificación incumplida           │              │ ✓            │
│ Invariante de código violada        │              │ ✓            │
│ Recurso agotado (memoria)           │              │ ✓            │
│ Llamada a función no registrada     │              │ ✓            │
│ Timeout de operación                │ ✓            │              │
│ Configuración del sistema incompleta│              │ ✓            │
└────────────────────────────────────┴──────────────┴──────────────┘
```

---

## Ejercicios Progresivos

### Ejercicio 1: Recuperación Simple

**Objetivo:** Escribir una función que hace panic y se recupera con defer.

```go
package main

import (
    "fmt"
)

// TODO: Implementar estas funciones
// 1. operacionRiesgosa() - Función que hace panic
// 2. ejecutarConRecuperacion() - Captura el panic y retorna error
// 3. main() - Demostrar el uso

// Requisitos:
// - operacionRiesgosa() debe hacer panic con mensaje
// - ejecutarConRecuperacion debe capturar y convertir a error
// - main debe manejar el error de forma segura

func main() {
    // Tu código aquí
}
```

**Solución esperada:**
- Función que hace panic
- Función wrapper con defer/recover
- Conversión de panic a error
- Output mostrando manejo seguro

---

### Ejercicio 2: Handler de Panics para Múltiples Operaciones

**Objetivo:** Crear un middleware que ejecuta múltiples operaciones de forma segura.

```go
package main

import (
    "fmt"
)

// TODO: Implementar
type OperacionHandler struct {
    operaciones []func() error
    resultados  []error
}

// 1. Agregar operaciones (que podrían hacer panic)
// 2. Ejecutar todas de forma segura
// 3. Reportar cuáles exitosas y cuáles fallaron

// Requisitos:
// - Debe ejecutar todas las operaciones incluso si algunas hacen panic
// - Debe convertir panic a error
// - Debe reportar estadísticas (exitosas, fallidas)

func main() {
    // Tu código aquí
}
```

**Solución esperada:**
- Handler que ejecuta múltiples tareas
- Recuperación individual de cada una
- Reporte de resultados
- Demostración con mix de éxito/panic

---

### Ejercicio 3: Cargar Configuración con Recuperación

**Objetivo:** Sistema de carga de configuración con panic recovery y defaults.

```go
package main

import (
    "fmt"
    "encoding/json"
)

// TODO: Implementar
type Config struct {
    Host     string
    Puerto   int
    Debug    bool
    MaxConn  int
}

// 1. cargarConfiguracion() - Lee JSON y hace panic si inválido
// 2. cargarConDefaults() - Wrapper que retorna config con valores por defecto si falla
// 3. Demostrar carga exitosa y fallida

// Requisitos:
// - Si JSON es inválido → panic
// - Si faltan campos obligatorios → panic
// - Recuperación retorna Config con valores por defecto
// - Logging de qué valores se usaron (originales o defaults)

func main() {
    // Tu código aquí
}
```

**Solución esperada:**
- Parser JSON que hace panic
- Wrapper con recuperación
- Valores por defecto apropiados
- Logging informativo

---

### Ejercicio 4: Protector de Goroutines (Supervisor)

**Objetivo:** Sistema que reinicia goroutines si hacen panic.

```go
package main

import (
    "fmt"
    "time"
)

// TODO: Implementar
type SupervisorGoroutine struct {
    tarea      func()
    maxReintentos int
    intentos   int
    // Más campos según necesites
}

// 1. Ejecutar tarea con supervisión
// 2. Si hace panic, reintentar hasta maxReintentos
// 3. Logging de cada intento y reinicio
// 4. Demostración con tarea que falla N veces antes de éxito

// Requisitos:
// - Capturar panics de goroutines
// - Reintentar automáticamente
// - Logging detallado
// - Fallar después de N intentos

func main() {
    // Tu código aquí
}
```

**Solución esperada:**
- Supervisor que monitorea goroutines
- Recuperación automática de panics
- Reintentos con límite
- Logging detallado de eventos

---

### Ejercicio 5: Sistema de Recuperación Multi-nivel

**Objetivo:** Sistema complejo con recuperación, logging y notificaciones.

```go
package main

import (
    "fmt"
    "log"
    "time"
)

// TODO: Implementar sistema completo

type EventoRecuperacion struct {
    Timestamp time.Time
    Mensaje   string
    Stack     string
    Tipo      string  // "warning", "error", "critical"
}

type SistemaRecuperacion struct {
    eventos     []EventoRecuperacion
    listeners   []func(EventoRecuperacion)
    logFile     string
}

// 1. Registrar handlers de notificación
// 2. Capturar panics con contexto completo
// 3. Registrar eventos
// 4. Notificar listeners
// 5. Guardar en archivo de log

// Requisitos:
// - Stack trace completo
// - Clasificación por severidad
// - Sistema de suscriptores (observadores)
// - Persistencia en archivo
// - Estadísticas de eventos

func main() {
    // Tu código aquí
}
```

**Solución esperada:**
- Sistema completo de recuperación
- Múltiples niveles de logging
- Pattern Observer implementado
- Stack traces capturados
- Persistencia de eventos
- Estadísticas

---

## Resumen del Capítulo

### Conceptos Clave

1. **Panic**: Mecanismo para situaciones EXCEPCIONALES e IRRECUPERABLES
2. **Recover**: ÚNICA forma de capturar panic (dentro de defer)
3. **Defer**: Se ejecuta SIEMPRE, incluso con panic
4. **Stack Unwinding**: Propagación de panic hacia arriba en la pila
5. **Error vs Panic**: Errores para lo esperado, panic para lo inesperado
6. **Goroutine Isolation**: Panic en goroutine NO mata todo el programa
7. **Built-in Panics**: Go lanza panic automáticamente en situaciones específicas

### Filosofía de Go

Go rechaza explícitamente el modelo de excepciones porque:

- **Errores son valores**: Se manejan explícitamente como retorno
- **Claridad**: Sé exactamente dónde puede fallar cada operación
- **Simpleza**: No hay try/catch/finally
- **Rendimiento**: Sin overhead de excepciones en camino normal

### Regla de Oro

```
┌─────────────────────────────────────────────────────────┐
│  REGLA DE ORO: Panic = Programmer Error                │
│                Error = Normal Failure                   │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Si el usuario/programa podría causar la situación:     │
│  → Usar error ✓                                         │
│                                                          │
│  Si es un fallo en la lógica interna del programa:      │
│  → Usar panic ✓                                         │
│                                                          │
│  Si es imposible determinar:                            │
│  → PREFIERE error (mejor ser defensivo)                │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## Apéndice: Comparativa con Otros Lenguajes

### Go vs Java

```
JAVA:
  try {
      conectar()
      procesar()
      guardar()
  } catch (SQLException e) {      // ¿De cuál línea?
      log.error(e)
  } catch (IOException e) {
      log.error(e)
  }
  
GO:
  if err := conectar(); err != nil {
      return fmt.Errorf("conectar: %w", err)  // Claro: de conectar()
  }
  if err := procesar(); err != nil {
      return fmt.Errorf("procesar: %w", err)  // Claro: de procesar()
  }
  if err := guardar(); err != nil {
      return fmt.Errorf("guardar: %w", err)   // Claro: de guardar()
  }
```

### Go vs Python

```
PYTHON:
  try:
      archivo = open("datos.txt")
      datos = json.load(archivo)
  except FileNotFoundError:
      usar_default()
  except JSONDecodeError:
      usar_default()
  finally:
      archivo.close()  # ¿Seguro si FileNotFoundError?
      
GO:
  archivo, err := os.Open("datos.txt")
  if err != nil {
      usar_default()
      return
  }
  defer archivo.Close()  // SIEMPRE se ejecuta
  
  var datos interface{}
  if err := json.NewDecoder(archivo).Decode(&datos); err != nil {
      usar_default()
      return
  }
```

### Go vs C (Signals)

```
C:
  signal(SIGSEGV, handler)  // No siempre es seguro
  // Comportamiento indefinido
  
GO:
  defer func() {
      if r := recover(); r != nil {
          // Seguro y controlado
      }
  }()
  // Operación riesgosa
```

---

## Checklist de Calidad

- [ ] Entendimiento de cuándo usar panic vs error
- [ ] Defers ejecutándose correctamente en orden LIFO
- [ ] Recuperación segura en goroutines
- [ ] Stack traces capturados e interpretados
- [ ] Pruebas incluyen verificación de panics
- [ ] No usar panic para control de flujo
- [ ] Información contextual en todos los panics
- [ ] Cleanup garantizado con defer en recursos
- [ ] Handlers de panic en puntos de entrada

---

**Fin del Capítulo 18**

*En el siguiente capítulo: CAPÍTULO 19: TESTING AVANZADO - Coverage, Benchmarks, y Fuzzing*

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/18-panic-y-recover/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/18-panic-y-recover):

```bash
cd examples/18-panic-y-recover
go run .
```
