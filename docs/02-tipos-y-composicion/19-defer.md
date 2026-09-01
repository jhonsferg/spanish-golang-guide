# Capítulo 19: Defer - Ejecución diferida y limpieza

## Introducción

El `defer` es uno de los mecanismos más poderosos y elegantes de Go para garantizar la ejecución de código de limpieza y finalización. A diferencia de otros lenguajes que utilizan bloques `finally`, constructores explícitos o context managers, Go proporciona una solución simple pero robusta: posponer la ejecución de una función hasta que la función que la contiene regrese.

Este capítulo te llevará desde los conceptos básicos de defer hasta patrones avanzados que garantizan la correcta liberación de recursos, manejo de transacciones y limpieza de estado. Comprenderás la semántica exacta de ejecución, los traps comunes y cómo evitarlos.

---

## 19.1 - ¿Qué es Defer?

### 19.1.1 Concepto Fundamental

`defer` es una instrucción que permite posponer la ejecución de una función hasta que la función envolvente termine. Es especialmente útil para garantizar que cierto código de limpieza siempre se ejecute, incluso si ocurren errores o retornos anticipados.

**Sintaxis básica:**

```go
defer funcionAlimpieza()
```

### 19.1.2 Motivación: Garantías de Limpieza

Antes de Go (y en muchos lenguajes):

```java
// Java: usar try-finally
FileReader reader = null;
try {
    reader = new FileReader("archivo.txt");
    // procesamiento
} finally {
    if (reader != null) {
        reader.close();
    }
}
```

Con Go y defer:

```go
// Go: simple y elegante
f, err := os.Open("archivo.txt")
if err != nil {
    return err
}
defer f.Close()
// procesamiento
// Close() se ejecutará automáticamente
```

### 19.1.3 Beneficios Principales

```go
package main

import (
    "fmt"
    "os"
)

func ejemploBeneficios() error {
    // Abrir archivo
    f, err := os.Open("datos.txt")
    if err != nil {
        return fmt.Errorf("no se pudo abrir: %w", err)
    }

    // Garantizar cierre INMEDIATAMENTE después de abrir
    defer f.Close()

    // Beneficio 1: Código más cercano al recurso
    // Beneficio 2: Imposible olvidar el cierre (si hay retorno anticipado)
    // Beneficio 3: Limpieza ocurre incluso con pánico

    if err := procesarArchivo(f); err != nil {
        return err  // f.Close() se ejecutará aquí
    }

    return nil  // f.Close() se ejecutará aquí también
}

func procesarArchivo(f *os.File) error {
    // procesamiento
    return nil
}
```

### 19.1.4 Comparativa: Go vs Otros Lenguajes

**Go - defer:**

```go
defer f.Close()
// Semántica clara, se ejecuta al salir de la función
```

**Python - context manager:**

```python
with open("archivo.txt") as f:
    # procesamiento
# Cierre garantizado al salir del with
```

**JavaScript - try-finally:**

```javascript
try {
    let f = fs.openSync("archivo.txt");
    // procesamiento
} finally {
    fs.closeSync(f);
}
```

**Ventaja de Go:** La simplicidad de `defer` reduce el código boilerplate y el riesgo de olvidar la limpieza.

---

## 19.2 - Semántica de Defer

### 19.2.1 Cuándo se Ejecuta

Defer no se ejecuta inmediatamente; se añade a una pila (stack) de funciones diferidas. Cuando la función envolvente regresa (via `return`, final de función, o pánico), todas las funciones diferidas se ejecutan en orden **LIFO** (Last In, First Out).

```go
package main

import "fmt"

func demoOrdenEjecucion() {
    fmt.Println("1. Inicio")

    defer fmt.Println("4. Defer 1 (primero añadido, último ejecutado)")
    defer fmt.Println("3. Defer 2 (segundo añadido)")
    defer fmt.Println("2. Defer 3 (tercero añadido, primero ejecutado)")

    fmt.Println("5. Fin de función")
}

// Salida:
// 1. Inicio
// 5. Fin de función
// 2. Defer 3 (tercero añadido, primero ejecutado)
// 3. Defer 2 (segundo añadido)
// 4. Defer 1 (primero añadido, último ejecutado)
```

### 19.2.2 Stack de Defers - Visualización

```
Ejecución de código:
┌─────────────────────────────────────┐
│ defer fmt.Println("A")              │  ← Se añade a la pila
│ defer fmt.Println("B")              │  ← Se añade a la pila
│ defer fmt.Println("C")              │  ← Se añade a la pila
│ (función termina)                   │
└─────────────────────────────────────┘

Stack de defers (LIFO):
┌─────────┐
│    C    │  ← Se ejecuta primero (LIFO)
├─────────┤
│    B    │  ← Se ejecuta segundo
├─────────┤
│    A    │  ← Se ejecuta último
└─────────┘

Orden de ejecución real: C, B, A
```

### 19.2.3 Ejecución Garantizada

Defer se ejecuta incluso si:

```go
func demoGarantías() {
    defer fmt.Println("Defer se ejecuta siempre")

    // Caso 1: Retorno normal
    if true {
        return  // Defer se ejecuta
    }

    fmt.Println("Esto no se ejecuta")
}

func demoPanico() {
    defer fmt.Println("Defer se ejecuta incluso con pánico")

    panic("Error crítico")  // Defer se ejecuta antes del pánico
}

func demoError() {
    defer fmt.Println("Defer se ejecuta con error")

    err := verificar()
    if err != nil {
        return  // Defer se ejecuta
    }
}

func verificar() error {
    return fmt.Errorf("algo falló")
}
```

### 19.2.4 Scope de Defer

Defer solo se ejecuta cuando se sale de la función que lo contiene:

```go
package main

import "fmt"

func demoScope() {
    fmt.Println("Función exterior inicio")

    defer func() {
        fmt.Println("Defer de exterior")
    }()

    func() {
        fmt.Println("Función interior")
        defer func() {
            fmt.Println("Defer de interior")
        }()
    }()  // Función interior termina aquí → defer interior se ejecuta

    fmt.Println("Volvimos a exterior")
}

// Salida:
// Función exterior inicio
// Función interior
// Defer de interior
// Volvimos a exterior
// Defer de exterior
```

---

## 19.3 - Defer en Funciones

### 19.3.1 Pattern: Acquire-Defer-Use

El patrón más común es: adquirir recurso, diferir liberación, usar recurso.

```go
package main

import (
    "fmt"
    "os"
    "sync"
)

// Patrón 1: Archivos
func procesarArchivo(ruta string) error {
    f, err := os.Open(ruta)
    if err != nil {
        return err
    }
    defer f.Close()  // Cierre garantizado

    // Usar archivo
    buffer := make([]byte, 1024)
    _, err = f.Read(buffer)
    return err
}

// Patrón 2: Locks
func operacionCritica(mu *sync.Mutex) {
    mu.Lock()
    defer mu.Unlock()  // Desbloqueo garantizado

    // Sección crítica
    fmt.Println("Acceso exclusivo")
}

// Patrón 3: Recursos Customizados
type ConexionBD struct {
    nombre string
}

func (c *ConexionBD) Close() error {
    fmt.Printf("Cerrando conexión a %s\n", c.nombre)
    return nil
}

func abrirConexion(url string) (*ConexionBD, error) {
    conn := &ConexionBD{nombre: url}
    fmt.Printf("Abriendo conexión a %s\n", conn.nombre)
    return conn, nil
}

func consultarBD() error {
    conn, err := abrirConexion("localhost:5432")
    if err != nil {
        return err
    }
    defer conn.Close()  // Cierre garantizado

    // Consulta
    fmt.Println("Ejecutando consulta")
    return nil
}
```

### 19.3.2 Inicialización y Finalización

Usar defer para patrones init/fin:

```go
package main

import (
    "fmt"
    "time"
)

type Timer struct {
    nombre  string
    inicio  time.Time
    eventos []string
}

func (t *Timer) Iniciar() {
    t.inicio = time.Now()
    t.eventos = append(t.eventos, "Iniciado")
    fmt.Printf("[%s] Iniciado\n", t.nombre)
}

func (t *Timer) Finalizar() {
    duracion := time.Since(t.inicio)
    t.eventos = append(t.eventos, fmt.Sprintf("Finalizado en %v", duracion))
    fmt.Printf("[%s] Finalizado en %v\n", t.nombre, duracion)
}

func operacionLarga() {
    timer := &Timer{nombre: "OperacionLarga"}
    timer.Iniciar()
    defer timer.Finalizar()  // Garantizar finalización

    // Simular trabajo
    time.Sleep(100 * time.Millisecond)

    // Incluso con retorno anticipado:
    if true {
        return  // Finalización se ejecuta
    }
}
```

### 19.3.3 Cleanup en Caso de Error

```go
package main

import (
    "fmt"
    "os"
)

func crearArchivoTemporal() (string, error) {
    // Crear archivo temporal
    f, err := os.CreateTemp(".", "temp-*.txt")
    if err != nil {
        return "", err
    }
    defer f.Close()

    nombre := f.Name()
    fmt.Printf("Archivo temporal creado: %s\n", nombre)

    // Si algo falla, limpiar
    defer func() {
        if err != nil {
            os.Remove(nombre)
            fmt.Printf("Archivo temporal removido por error: %s\n", nombre)
        }
    }()

    // Escribir contenido
    _, err = f.WriteString("contenido")
    if err != nil {
        return "", err
    }

    return nombre, nil
}
```

---

## 19.4 - Múltiples Defers

### 19.4.1 Orden LIFO Detallado

```go
package main

import "fmt"

func demoMultiplesDefers() {
    fmt.Println("=== Inicio ===")

    // Stack de defers: [D3, D2, D1]
    defer fmt.Println("D1: Primero añadido, último ejecutado")

    // Stack de defers: [D3, D2, D1]
    defer fmt.Println("D2: Segundo añadido, segundo ejecutado")

    // Stack de defers: [D3, D2, D1]
    defer fmt.Println("D3: Tercero añadido, primero ejecutado")

    fmt.Println("=== Fin de código ===")
}

// Salida:
// === Inicio ===
// === Fin de código ===
// D3: Tercero añadido, primero ejecutado
// D2: Segundo añadido, segundo ejecutado
// D1: Primero añadido, último ejecutado
```

### 19.4.2 Patrones de Reversión

Los defers se usan frecuentemente para revertir operaciones en orden inverso:

```go
package main

import "fmt"

func demoReversion() {
    fmt.Println("Operación 1")
    defer func() {
        fmt.Println("Revertir operación 1")
    }()

    fmt.Println("Operación 2")
    defer func() {
        fmt.Println("Revertir operación 2")
    }()

    fmt.Println("Operación 3")
    defer func() {
        fmt.Println("Revertir operación 3")
    }()

    fmt.Println("Todas operaciones completadas")
}

// Salida:
// Operación 1
// Operación 2
// Operación 3
// Todas operaciones completadas
// Revertir operación 3
// Revertir operación 2
// Revertir operación 1
```

### 19.4.3 Visualización del Stack

```
Estado en diferentes puntos:

Punto 1: defer fmt.Println("A")
Stack: [A]

Punto 2: defer fmt.Println("B")
Stack: [B, A]  (B en el tope)

Punto 3: defer fmt.Println("C")
Stack: [C, B, A]  (C en el tope)

Ejecución (saliendo de función):
1. Pop C → ejecutar → Salida: C
2. Pop B → ejecutar → Salida: B
3. Pop A → ejecutar → Salida: A
```

### 19.4.4 Debugging de Múltiples Defers

```go
package main

import (
    "fmt"
    "log"
)

func funcionConMultiplesDefers() {
    log.Println("Función inicio")

    defer func() {
        log.Println("[defer 1] Limpieza de recurso A")
    }()

    defer func() {
        log.Println("[defer 2] Limpieza de recurso B")
    }()

    defer func() {
        log.Println("[defer 3] Limpieza de recurso C")
    }()

    log.Println("Cuerpo de función")
    log.Println("Función terminando")
}

// Para entender el orden de ejecución:
// 1. Lee los logs secuencialmente hasta "Función terminando"
// 2. Los defers se ejecutan en ORDEN INVERSO al que fueron añadidos
// 3. Usa defer cerca del recurso que protege para claridad
```

---

## 19.5 - Defer con Argumentos

### 19.5.1 Evaluación Temprana de Argumentos

Un aspecto crucial: los argumentos de una función deferada se evalúan INMEDIATAMENTE, no cuando se ejecuta.

```go
package main

import "fmt"

func demoEvaluacionTemprana() {
    x := 5

    defer func(valor int) {
        fmt.Printf("Defer recibió: %d\n", valor)
    }(x)

    x = 10
    fmt.Println("x ahora es:", x)
}

// Salida:
// x ahora es: 10
// Defer recibió: 5  ← No 10, porque se capturó x=5
```

**Visualización:**

```
Paso 1: x := 5
        x = 5

Paso 2: defer func(valor int) { ... }(x)
        Se evalúa x inmediatamente → 5 se pasa como argumento
        Se programa la ejecución diferida con valor=5

Paso 3: x = 10
        x = 10

Paso 4: Función termina
        Se ejecuta defer con el valor capturado (5), no el valor actual (10)
```

### 19.5.2 Captura de Variables

Comparar captura de argumentos vs captura de closures:

```go
package main

import "fmt"

func demoCaptura() {
    x := 5

    // Método 1: Argumento (evaluación temprana)
    defer func(v int) {
        fmt.Printf("Argumento: %d\n", v)
    }(x)

    // Método 2: Closure (referencia tardía)
    defer func() {
        fmt.Printf("Closure: %d\n", x)
    }()

    x = 10
}

// Salida:
// Closure: 10    ← Ve el valor actual
// Argumento: 5   ← Ve el valor capturado
```

### 19.5.3 Errores Comunes con Argumentos

```go
package main

import (
    "fmt"
    "time"
)

// ❌ Incorrecto: Olvidar pasar el valor
func incorrecto() {
    tick := time.Now()
    defer func() {
        // tick fue evaluado aquí, no cuando se creó el defer
        duracion := time.Since(tick)
        fmt.Println("Duración:", duracion)
    }()
    time.Sleep(100 * time.Millisecond)
}

// ✓ Correcto: Pasar el instante de tiempo
func correcto() {
    defer func(t time.Time) {
        duracion := time.Since(t)
        fmt.Println("Duración:", duracion)
    }(time.Now())

    time.Sleep(100 * time.Millisecond)
}

// ❌ Incorrecto: Loop con captura (se verá después)
func bucleIncorrecto() {
    for i := 0; i < 3; i++ {
        defer func() {
            fmt.Printf("i = %d\n", i)  // Siempre imprimirá 3
        }()
    }
}

// ✓ Correcto: Pasar i como argumento
func bucleCorrecto() {
    for i := 0; i < 3; i++ {
        defer func(val int) {
            fmt.Printf("i = %d\n", val)  // Imprimirá 2, 1, 0
        }(i)
    }
}
```

---

## 19.6 - Closures en Defer

### 19.6.1 Captura de Variables por Referencia

```go
package main

import "fmt"

func demoClosure() {
    contador := 0

    // Este defer captura 'contador' por referencia
    defer func() {
        fmt.Printf("Contador final: %d\n", contador)
    }()

    contador = 5
    fmt.Printf("Contador actual: %d\n", contador)
}

// Salida:
// Contador actual: 5
// Contador final: 5
```

### 19.6.2 Closures para Manejo de Errores

```go
package main

import (
    "fmt"
    "sync"
)

func transaccion() error {
    var mu sync.Mutex
    mu.Lock()
    defer mu.Unlock()

    var err error

    // Usar defer con closure para capturar err
    defer func() {
        if err != nil {
            fmt.Printf("Error en transacción: %v\n", err)
        }
    }()

    err = operacion1()
    if err != nil {
        return err
    }

    err = operacion2()
    return err
}

func operacion1() error {
    return nil
}

func operacion2() error {
    return fmt.Errorf("falló operación 2")
}
```

### 19.6.3 Closures para Logging y Timing

```go
package main

import (
    "fmt"
    "time"
)

func medir(nombre string) func() {
    inicio := time.Now()
    fmt.Printf("[%s] Iniciado\n", nombre)

    // Retornar closure que será usado como defer
    return func() {
        duracion := time.Since(inicio)
        fmt.Printf("[%s] Completado en %v\n", nombre, duracion)
    }
}

func ejemploMedida() {
    defer medir("OperacionA")()
    defer medir("OperacionB")()

    time.Sleep(50 * time.Millisecond)
}

// Alternativa: inline closure
func ejemploInline() {
    defer func() {
        fmt.Println("[Ejemplo] Ejecutando cleanup")
    }()

    fmt.Println("[Ejemplo] Cuerpo de función")
}
```

### 19.6.4 Captura de Punteros

```go
package main

import "fmt"

type Estado struct {
    nombre string
    valor  int
}

func demoCapturaPunteros() {
    estado := &Estado{nombre: "test", valor: 0}

    defer func() {
        fmt.Printf("Estado final: %+v\n", estado)
    }()

    estado.valor = 42
    estado.nombre = "modificado"
}

// Salida: Estado final: &{nombre:modificado valor:42}
```

---

## 19.7 - Defer en Loops

### 19.7.1 El Trap Clásico

```go
package main

import (
    "fmt"
    "os"
)

// ❌ INCORRECTO: Defer dentro de loop
func procesarArchivosIncorrecto(archivos []string) {
    for _, archivo := range archivos {
        f, err := os.Open(archivo)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        // PROBLEMA: Todos los defers se ejecutan al final de la función
        // Los archivos no se cierran hasta que termina la función
        defer f.Close()

        // Procesamiento (pero archivo sigue abierto)
        procesar(f)
    }

    // Aquí se ejecutan TODOS los closes
}

func procesar(f *os.File) {
    fmt.Printf("Procesando: %s\n", f.Name())
}

// Resultado: Si hay muchos archivos, se pueden agotar los file descriptors
```

### 19.7.2 Solución 1: Función Anidada

```go
package main

import (
    "fmt"
    "os"
)

// ✓ CORRECTO: Envolver en función anidada
func procesarArchivosCorrectov1(archivos []string) {
    for _, archivo := range archivos {
        err := procesarUnArchivo(archivo)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        }
    }
}

func procesarUnArchivo(archivo string) error {
    f, err := os.Open(archivo)
    if err != nil {
        return err
    }

    // Defer dentro de función anidada se ejecuta al salir de la función
    defer f.Close()

    procesar2(f)
    return nil
}

func procesar2(f *os.File) {
    fmt.Printf("Procesando: %s\n", f.Name())
}
```

### 19.7.3 Solución 2: Closure Inmediato

```go
package main

import (
    "fmt"
    "os"
)

// ✓ CORRECTO: Closure anónimo inmediato
func procesarArchivosCorrectov2(archivos []string) {
    for _, archivo := range archivos {
        // Función anónima INMEDIATA que captura 'archivo'
        func(nombre string) {
            f, err := os.Open(nombre)
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                return
            }

            // Defer se ejecuta al salir del closure
            defer f.Close()

            procesar3(f)
        }(archivo)  // ← Ejecución inmediata
    }
}

func procesar3(f *os.File) {
    fmt.Printf("Procesando: %s\n", f.Name())
}
```

### 19.7.4 Comparativa de Soluciones

```
Trap Original:
┌─────────────────────────┐
│ Loop                    │
│  ├─ Abrir archivo 1     │
│  ├─ defer close 1       │
│  ├─ Abrir archivo 2     │
│  ├─ defer close 2       │
│  ├─ Abrir archivo 3     │
│  └─ defer close 3       │
│                         │
│ (Aquí se ejecutan)      │
│ ← close 3               │
│ ← close 2               │
│ ← close 1               │
└─────────────────────────┘

Solución 1 (función anidada):
┌─────────────────────────┐
│ Loop                    │
│  ├─ procesarUnArchivo() │
│  │   ├─ Abrir          │
│  │   ├─ defer close    │
│  │   ├─ Procesar       │
│  │   └─ close          │
│  ├─ procesarUnArchivo() │
│  │   ├─ Abrir          │
│  │   ├─ defer close    │
│  │   ├─ Procesar       │
│  │   └─ close          │
│  ...                    │
└─────────────────────────┘

Solución 2 (closure inmediato):
┌─────────────────────────┐
│ Loop                    │
│  ├─ func(archivo1)()    │
│  │   ├─ Abrir          │
│  │   ├─ defer close    │
│  │   ├─ Procesar       │
│  │   └─ close          │
│  ├─ func(archivo2)()    │
│  │   ├─ Abrir          │
│  │   ├─ defer close    │
│  │   ├─ Procesar       │
│  │   └─ close          │
│  ...                    │
└─────────────────────────┘
```

---

## 19.8 - Defer y Valores de Retorno

### 19.8.1 Modificar Valores Retornados

Una característica poderosa: los defers pueden modificar los valores que una función retorna:

```go
package main

import "fmt"

// Usando named return values
func funcionConRetornoNombrado() (resultado string, err error) {
    defer func() {
        if err != nil {
            resultado = "error: " + resultado
        }
    }()

    // Asignación a variables nombradas
    resultado = "éxito"
    err = nil

    return  // Usa los valores nombrados
}

// Sin named return values (no se puede modificar)
func funcionSinRetornoNombrado() (string, error) {
    resultado := "éxito"
    err := error(nil)

    defer func() {
        // Aquí NO se puede modificar resultado porque es una copia
        _ = resultado
    }()

    return resultado, err
}

// Correcto: Usar punteros
func funcionConPuntero(resultado *string) error {
    defer func() {
        *resultado += " - modificado por defer"
    }()

    *resultado = "valor inicial"
    return nil
}
```

### 19.8.2 Patrón: Captura y Modificación de Error

```go
package main

import (
    "fmt"
)

// Patrón común: enriquecer errores
func operacionConEnriquecimiento() (err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("operación falló: %w", err)
        }
    }()

    // Aquí ocurre un error
    return fmt.Errorf("algo salió mal")
}

// Patrón: Logging de retorno
func funcionConLogging() (resultado int, err error) {
    defer func() {
        if err != nil {
            fmt.Printf("Función terminó con error: %v\n", err)
        } else {
            fmt.Printf("Función retornó: %d\n", resultado)
        }
    }()

    return 42, nil
}

// Patrón: Panic recovery
func funcionConRecupacion() (resultado string, err error) {
    defer func() {
        if p := recover(); p != nil {
            resultado = ""
            err = fmt.Errorf("recuperado de pánico: %v", p)
        }
    }()

    panic("algo grave")
}
```

### 19.8.3 Trampa: Modificar Retorno Sin Named Returns

```go
package main

import "fmt"

// ❌ Esto NO funciona
func noFunciona() string {
    resultado := "original"

    defer func() {
        resultado = "modificado"  // Modifica la copia, no el retorno
    }()

    return resultado
}

// ✓ Esto funciona
func funciona() (resultado string) {
    resultado = "original"

    defer func() {
        resultado = "modificado"  // Modifica la variable nombrada
    }()

    return
}

func demoModificacion() {
    fmt.Println("noFunciona:", noFunciona())  // original
    fmt.Println("funciona:", funciona())      // modificado
}
```

---

## 19.9 - Defer para Limpieza

### 19.9.1 Limpieza de Archivos

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func limpiarArchivos() error {
    // Lectura
    entrada, err := os.Open("entrada.txt")
    if err != nil {
        return err
    }
    defer entrada.Close()

    // Escritura
    salida, err := os.Create("salida.txt")
    if err != nil {
        return err
    }
    defer salida.Close()

    // Ambos archivos se cierran automáticamente al final

    // Procesamiento
    scanner := bufio.NewScanner(entrada)
    writer := bufio.NewWriter(salida)
    defer writer.Flush()

    for scanner.Scan() {
        writer.WriteString(scanner.Text() + "\n")
    }

    return scanner.Err()
}
```

### 19.9.2 Limpieza de Conexiones

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "sync"
)

type ConexionBD struct {
    db *sql.DB
}

func (c *ConexionBD) Close() error {
    return c.db.Close()
}

func (c *ConexionBD) Consultar(query string) ([]string, error) {
    filas, err := c.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer filas.Close()

    var resultados []string
    for filas.Next() {
        var valor string
        if err := filas.Scan(&valor); err != nil {
            return nil, err
        }
        resultados = append(resultados, valor)
    }

    return resultados, filas.Err()
}

func ejecutarOperacionBD() error {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        return err
    }

    conn := &ConexionBD{db: db}
    defer conn.Close()

    resultados, err := conn.Consultar("SELECT 1")
    if err != nil {
        return err
    }

    fmt.Println("Resultados:", resultados)
    return nil
}
```

### 19.9.3 Limpieza de Locks

```go
package main

import (
    "fmt"
    "sync"
)

type RecursoProtegido struct {
    mu    sync.Mutex
    datos map[string]int
}

func (r *RecursoProtegido) Leer(clave string) int {
    r.mu.Lock()
    defer r.mu.Unlock()

    return r.datos[clave]
}

func (r *RecursoProtegido) Escribir(clave string, valor int) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.datos[clave] = valor
    fmt.Printf("Escribió %s = %d\n", clave, valor)
}

func (r *RecursoProtegido) ActualizarMultiple(actualizaciones map[string]int) {
    r.mu.Lock()
    defer r.mu.Unlock()

    for k, v := range actualizaciones {
        r.datos[k] = v
    }
}

func demoBD() {
    recurso := &RecursoProtegido{datos: make(map[string]int)}

    // Lock se libera automáticamente
    recurso.Escribir("x", 10)

    valor := recurso.Leer("x")
    fmt.Println("Valor leído:", valor)
}
```

### 19.9.4 Limpieza de Transacciones

```go
package main

import (
    "fmt"
)

type Transaccion struct {
    id        string
    activa    bool
    cambios   []string
}

func (t *Transaccion) Comenzar() {
    t.activa = true
    t.cambios = []string{}
    fmt.Printf("[TX %s] Iniciada\n", t.id)
}

func (t *Transaccion) Registrar(cambio string) error {
    if !t.activa {
        return fmt.Errorf("transacción no activa")
    }
    t.cambios = append(t.cambios, cambio)
    return nil
}

func (t *Transaccion) Commit() error {
    if !t.activa {
        return fmt.Errorf("transacción no activa")
    }
    fmt.Printf("[TX %s] Commit con %d cambios\n", t.id, len(t.cambios))
    t.activa = false
    return nil
}

func (t *Transaccion) Rollback() error {
    if !t.activa {
        return fmt.Errorf("transacción no activa")
    }
    fmt.Printf("[TX %s] Rollback\n", t.id)
    t.cambios = []string{}
    t.activa = false
    return nil
}

func operacionConTransaccion() (err error) {
    tx := &Transaccion{id: "TX001"}
    tx.Comenzar()

    // Si ocurre error, hacer rollback
    defer func() {
        if err != nil {
            tx.Rollback()
        } else if tx.activa {
            tx.Commit()
        }
    }()

    if err = tx.Registrar("UPDATE usuarios"); err != nil {
        return
    }

    if err = tx.Registrar("UPDATE productos"); err != nil {
        return
    }

    return nil
}
```

---

## 19.10 - Patrones Avanzados

### 19.10.1 Patrón: Ensure (Always Execute)

```go
package main

import (
    "fmt"
)

// Garantizar que cierto código siempre se ejecute
func ensure(fn func()) {
    defer fn()
}

func ejemploEnsure() {
    contador := 0

    ensure(func() {
        fmt.Println("Limpieza, contador final:", contador)
    })

    contador = 5
    fmt.Println("Contador:", contador)
}

// Patrón con parámetros
func guardarRecurso(recurso interface{}) func() {
    return func() {
        fmt.Printf("Guardando/Liberando recurso: %T\n", recurso)
    }
}

func ejemploGuardarRecurso() {
    defer guardarRecurso("datos")()

    fmt.Println("Procesando recurso")
}
```

### 19.10.2 Patrón: Resource Guarding

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type PoolRecursos struct {
    mu       sync.Mutex
    libres   []interface{}
    en_uso   int
    maximo   int
}

func (p *PoolRecursos) Adquirir() (interface{}, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if len(p.libres) == 0 {
        if p.en_uso >= p.maximo {
            return nil, fmt.Errorf("pool lleno")
        }
        p.en_uso++
        return fmt.Sprintf("recurso_%d", p.en_uso), nil
    }

    recurso := p.libres[len(p.libres)-1]
    p.libres = p.libres[:len(p.libres)-1]
    return recurso, nil
}

func (p *PoolRecursos) Liberar(recurso interface{}) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.libres = append(p.libres, recurso)
}

func (p *PoolRecursos) WithRecurso(fn func(interface{}) error) error {
    recurso, err := p.Adquirir()
    if err != nil {
        return err
    }

    defer p.Liberar(recurso)

    return fn(recurso)
}

func demoPiscina() {
    pool := &PoolRecursos{libres: []interface{}{}, maximo: 5}

    // Usar patrón WithRecurso
    err := pool.WithRecurso(func(recurso interface{}) error {
        fmt.Println("Usando:", recurso)
        time.Sleep(10 * time.Millisecond)
        return nil
    })

    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### 19.10.3 Patrón: Cleanup Chains

```go
package main

import (
    "fmt"
)

// Cadena de cleanup automática
type CleanupChain struct {
    handlers []func() error
}

func (cc *CleanupChain) Add(handler func() error) {
    cc.handlers = append(cc.handlers, handler)
}

func (cc *CleanupChain) ExecuteAll() error {
    // Ejecutar en orden inverso
    for i := len(cc.handlers) - 1; i >= 0; i-- {
        if err := cc.handlers[i](); err != nil {
            fmt.Printf("Error en cleanup: %v\n", err)
            // Continuar con los demás, pero registrar error
        }
    }
    return nil
}

func operacionMultiRecurso() error {
    cleanup := &CleanupChain{}

    // Adquirir recurso 1
    recurso1 := abrirRecurso("recurso1")
    cleanup.Add(func() error {
        return cerrarRecurso(recurso1)
    })

    // Adquirir recurso 2
    recurso2 := abrirRecurso("recurso2")
    cleanup.Add(func() error {
        return cerrarRecurso(recurso2)
    })

    // Adquirir recurso 3
    recurso3 := abrirRecurso("recurso3")
    cleanup.Add(func() error {
        return cerrarRecurso(recurso3)
    })

    // Garantizar cleanup
    defer cleanup.ExecuteAll()

    // Usar recursos
    fmt.Println("Usando:", recurso1, recurso2, recurso3)

    return nil
}

func abrirRecurso(nombre string) string {
    fmt.Printf("Abriendo %s\n", nombre)
    return nombre
}

func cerrarRecurso(nombre string) error {
    fmt.Printf("Cerrando %s\n", nombre)
    return nil
}
```

### 19.10.4 Patrón: Timing y Profiling

```go
package main

import (
    "fmt"
    "time"
)

func profile(nombre string) func() {
    inicio := time.Now()
    fmt.Printf("┌─ [%s] Inicio\n", nombre)

    return func() {
        duracion := time.Since(inicio)
        fmt.Printf("└─ [%s] Fin en %v\n", nombre, duracion)
    }
}

func operacionConProfiling() {
    defer profile("operación completa")()

    defer profile("paso 1")()
    simularTrabajo(50 * time.Millisecond)

    defer profile("paso 2")()
    simularTrabajo(75 * time.Millisecond)

    defer profile("paso 3")()
    simularTrabajo(100 * time.Millisecond)
}

func simularTrabajo(duracion time.Duration) {
    time.Sleep(duracion)
}

// Salida:
// ┌─ [operación completa] Inicio
// ┌─ [paso 1] Inicio
// └─ [paso 1] Fin en 50.1ms
// ┌─ [paso 2] Inicio
// └─ [paso 2] Fin en 75.2ms
// ┌─ [paso 3] Inicio
// └─ [paso 3] Fin en 100.3ms
// └─ [operación completa] Fin en 225.6ms
```

---

## 19.11 - Buenas Prácticas y Antipatrones

### 19.11.1 Buena Práctica: Defer Inmediatamente

```go
package main

import (
    "fmt"
    "os"
)

// ✓ BUENA: Defer inmediatamente después de adquirir
func buena_practica() error {
    f, err := os.Open("archivo.txt")
    if err != nil {
        return err
    }
    defer f.Close()  // ← Aquí, cerca del recurso

    // Procesamiento
    return nil
}

// ❌ MALA: Defer alejado del recurso
func mala_practica() error {
    f, err := os.Open("archivo.txt")
    if err != nil {
        return err
    }

    // Mucho código...
    _, err = f.ReadAt(make([]byte, 10), 0)
    if err != nil {
        return err
    }

    // Defer alejado del recurso
    defer f.Close()

    return nil
}

// Razón: El defer debe estar cerca para que sea obvio
// qué recurso se está protegiendo
```

### 19.11.2 Buena Práctica: Manejo de Errores

```go
package main

import (
    "fmt"
)

// ✓ BUENA: Capturar y propagar errores
func manejoErroresCorrecto() (err error) {
    recurso := abrirRecurso2()
    defer func() {
        if cerr := cerrarRecurso2(recurso); cerr != nil && err == nil {
            err = cerr
        }
    }()

    if err = procesarRecurso(recurso); err != nil {
        return
    }

    return nil
}

// ❌ MALA: Ignorar errores en defer
func manejoErroresIncorrecto() error {
    recurso := abrirRecurso2()
    defer cerrarRecurso2(recurso)  // Error se ignora

    return procesarRecurso(recurso)
}

func abrirRecurso2() string {
    return "recurso"
}

func cerrarRecurso2(r string) error {
    fmt.Printf("Cerrando %s\n", r)
    return nil
}

func procesarRecurso(r string) error {
    return nil
}
```

### 19.11.3 Antipatrón: Defer en Loops

```go
package main

// ❌ ANTIPATRÓN: Ya visto, pero recapitular
func antipatron_loop() {
    for i := 0; i < 1000; i++ {
        f := abrirArchivo2()
        defer f.Close()  // ← NO HACER ESTO
    }
    // Todos los 1000 closes se ejecutan al final
    // Posible exhaustión de file descriptors
}

// ✓ CORRECTO: Usar función anidada
func patron_correcto_loop() {
    for i := 0; i < 1000; i++ {
        procesarArchivoCorrecto()
    }
}

func procesarArchivoCorrecto() {
    f := abrirArchivo2()
    defer f.Close()  // Se ejecuta al salir de esta función

    // Procesamiento
}

func abrirArchivo2() *MockFile {
    return &MockFile{nombre: "archivo"}
}

type MockFile struct {
    nombre string
}

func (m *MockFile) Close() error {
    return nil
}
```

### 19.11.4 Antipatrón: Defer Sin Pair

```go
package main

import (
    "fmt"
)

// ❌ ANTIPATRÓN: Defer sin recurso correspondiente
func antipatron_sin_pair() {
    defer fmt.Println("cleanup")  // ¿Qué se está limpiando?

    // Código que no tiene relación con el defer
    x := 5
    y := 10
    fmt.Println(x + y)
}

// ✓ CORRECTO: Defer para algo específico
func patron_correcto_pair() {
    mu := &mockMutex{}
    mu.Lock()
    defer mu.Unlock()  // Claro qué se está limpiando

    // Sección crítica
    fmt.Println("Sección crítica")
}

type mockMutex struct{}

func (m *mockMutex) Lock()   {}
func (m *mockMutex) Unlock() {}
```

### 19.11.5 Antipatrón: Defer con Efectos Secundarios Confusos

```go
package main

import (
    "fmt"
)

// ❌ ANTIPATRÓN: Defer hace algo inesperado
func antipatron_confuso() (resultado int) {
    resultado = 10

    defer func() {
        resultado += 5  // Confuso: ¿por qué suma?
    }()

    return resultado  // Retorna 15, no 10
}

// ✓ CORRECTO: Defer tiene propósito claro
func patron_correcto_claro() (resultado int, err error) {
    resultado = 10

    defer func() {
        if err != nil {
            resultado = -1  // Claro: si hay error, resultado es -1
        }
    }()

    return resultado, err
}
```

### 19.11.6 Testing con Defer

```go
package main

import (
    "fmt"
    "testing"
)

func TestConDefer(t *testing.T) {
    // Usar defer para limpieza de test
    defer func() {
        fmt.Println("Limpieza de test")
    }()

    // Setup
    datos := setup()
    defer cleanup(datos)

    // Test
    resultado := procesar(datos)
    if resultado != esperado {
        t.Errorf("Resultado incorrecto: %v", resultado)
    }
}

func setup() string {
    return "datos"
}

func cleanup(d string) {
    fmt.Printf("Limpiando: %s\n", d)
}

func procesar(d string) int {
    return 42
}

const esperado = 42
```

---

## Ejercicios Prácticos

### Ejercicio 1: Manejo de Archivos con Defer

**Objetivo:** Crear una función que procese múltiples archivos, usando defer para garantizar el cierre.

**Archivo:** `ejercicio1_archivos.go`

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// Ejercicio 1: Procesar archivo con defer
//
// Requisitos:
// 1. Crear una función que lea líneas de un archivo
// 2. Usar defer para garantizar cierre del archivo
// 3. Contar líneas que contienen una palabra clave
// 4. Retornar count y error
// 5. Garantizar que el archivo se cierre incluso con error
//
// Implementar:
func ContarLineasConPalabra(ruta, palabra string) (int, error) {
    // COMPLETAR
    return 0, nil
}

// Solución esperada:
func ContarLineasConPalabraSol(ruta, palabra string) (int, error) {
    archivo, err := os.Open(ruta)
    if err != nil {
        return 0, err
    }
    defer archivo.Close()  // Garantizar cierre

    scanner := bufio.NewScanner(archivo)
    contador := 0

    for scanner.Scan() {
        linea := scanner.Text()
        if strings.Contains(strings.ToLower(linea), strings.ToLower(palabra)) {
            contador++
        }
    }

    return contador, scanner.Err()
}

// Casos de prueba:
func TestEjercicio1() {
    // Crear archivo de prueba
    contenido := "Go es un lenguaje moderno\nGo tiene defer\nRust también existe\nGo es excelente"
    err := os.WriteFile("test1.txt", []byte(contenido), 0644)
    if err != nil {
        panic(err)
    }
    defer os.Remove("test1.txt")

    // Test 1: Buscar "Go"
    count, err := ContarLineasConPalabraSol("test1.txt", "Go")
    if err != nil {
        fmt.Println("Error:", err)
    }
    if count != 3 {
        fmt.Printf("Error: esperaba 3, obtuve %d\n", count)
    } else {
        fmt.Println("✓ Test 1 pasado: 3 líneas con 'Go'")
    }

    // Test 2: Buscar "Rust"
    count, err = ContarLineasConPalabraSol("test1.txt", "Rust")
    if count != 1 {
        fmt.Printf("Error: esperaba 1, obtuve %d\n", count)
    } else {
        fmt.Println("✓ Test 2 pasado: 1 línea con 'Rust'")
    }

    // Test 3: Archivo inexistente
    count, err = ContarLineasConPalabraSol("inexistente.txt", "test")
    if err == nil {
        fmt.Println("Error: debería haber error para archivo inexistente")
    } else {
        fmt.Println("✓ Test 3 pasado: error para archivo inexistente")
    }
}
```

---

### Ejercicio 2: Transacción con Rollback

**Objetivo:** Implementar un sistema de transacciones que use defer para commit/rollback automático.

**Archivo:** `ejercicio2_transaccion.go`

```go
package main

import (
    "fmt"
)

// Ejercicio 2: Transacción con defer para commit/rollback
//
// Requisitos:
// 1. Crear estructura Transaccion con estado
// 2. Implementar Comenzar(), Registrar(), Commit(), Rollback()
// 3. Usar defer para garantizar commit o rollback
// 4. Si ocurre error, hacer rollback automático
// 5. Si no hay error, hacer commit automático
//
// Crear un tipo Transaccion:
type Transaccion struct {
    id     string
    activa bool
    datos  map[string]interface{}
    logs   []string
}

// Implementar métodos:
func (t *Transaccion) Comenzar() {
    // COMPLETAR
}

func (t *Transaccion) Registrar(clave string, valor interface{}) error {
    // COMPLETAR
    return nil
}

func (t *Transaccion) Commit() error {
    // COMPLETAR
    return nil
}

func (t *Transaccion) Rollback() error {
    // COMPLETAR
    return nil
}

// Solución esperada:
type TransaccionSol struct {
    id     string
    activa bool
    datos  map[string]interface{}
    logs   []string
}

func (t *TransaccionSol) Comenzar() {
    t.activa = true
    t.datos = make(map[string]interface{})
    t.logs = append(t.logs, "TX INICIO")
    fmt.Printf("[%s] Transacción iniciada\n", t.id)
}

func (t *TransaccionSol) Registrar(clave string, valor interface{}) error {
    if !t.activa {
        return fmt.Errorf("transacción no activa")
    }
    t.datos[clave] = valor
    t.logs = append(t.logs, fmt.Sprintf("REGISTRAR: %s=%v", clave, valor))
    return nil
}

func (t *TransaccionSol) Commit() error {
    if !t.activa {
        return fmt.Errorf("transacción no activa")
    }
    fmt.Printf("[%s] COMMIT con %d cambios\n", t.id, len(t.datos))
    t.activa = false
    return nil
}

func (t *TransaccionSol) Rollback() error {
    if !t.activa {
        return fmt.Errorf("transacción no activa")
    }
    fmt.Printf("[%s] ROLLBACK\n", t.id)
    t.datos = make(map[string]interface{})
    t.activa = false
    return nil
}

// Patrón de uso con defer:
func EjecutarTransaccion(fn func(*TransaccionSol) error) error {
    tx := &TransaccionSol{id: "TX001"}
    tx.Comenzar()

    var err error
    defer func() {
        if err != nil {
            tx.Rollback()
        } else if tx.activa {
            tx.Commit()
        }
    }()

    err = fn(tx)
    return err
}

// Casos de prueba:
func TestEjercicio2() {
    // Test 1: Transacción exitosa
    fmt.Println("\n=== Test 1: Transacción exitosa ===")
    err := EjecutarTransaccion(func(tx *TransaccionSol) error {
        tx.Registrar("usuario_id", 123)
        tx.Registrar("nombre", "Juan")
        return nil
    })
    if err == nil {
        fmt.Println("✓ Test 1 pasado: Transacción completada")
    }

    // Test 2: Transacción con error (rollback)
    fmt.Println("\n=== Test 2: Transacción con error ===")
    err = EjecutarTransaccion(func(tx *TransaccionSol) error {
        tx.Registrar("producto_id", 456)
        return fmt.Errorf("error en la transacción")
    })
    if err != nil {
        fmt.Println("✓ Test 2 pasado: Rollback automático")
    }
}
```

---

### Ejercicio 3: Mutex y Locks

**Objetivo:** Usar defer para liberar locks automáticamente en operaciones concurrentes.

**Archivo:** `ejercicio3_mutex.go`

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Ejercicio 3: Usar defer para liberar mutex
//
// Requisitos:
// 1. Crear estructura CuentaBancaria con saldo y mutex
// 2. Implementar Depositar(cantidad) con lock/unlock via defer
// 3. Implementar Retirar(cantidad) con lock/unlock via defer
// 4. Implementar Saldo() que devuelva saldo actual
// 5. Garantizar que nunca hay race condition
//
// Estructura base:
type CuentaBancaria struct {
    mu    sync.Mutex
    saldo float64
}

// Implementar métodos:
func (c *CuentaBancaria) Depositar(cantidad float64) error {
    // COMPLETAR
    return nil
}

func (c *CuentaBancaria) Retirar(cantidad float64) error {
    // COMPLETAR
    return nil
}

func (c *CuentaBancaria) Saldo() float64 {
    // COMPLETAR
    return 0
}

// Solución esperada:
type CuentaBancariaSol struct {
    mu    sync.Mutex
    saldo float64
}

func (c *CuentaBancariaSol) Depositar(cantidad float64) error {
    if cantidad <= 0 {
        return fmt.Errorf("cantidad debe ser positiva")
    }

    c.mu.Lock()
    defer c.mu.Unlock()

    c.saldo += cantidad
    fmt.Printf("Depositado: +%.2f, saldo: %.2f\n", cantidad, c.saldo)
    return nil
}

func (c *CuentaBancariaSol) Retirar(cantidad float64) error {
    if cantidad <= 0 {
        return fmt.Errorf("cantidad debe ser positiva")
    }

    c.mu.Lock()
    defer c.mu.Unlock()

    if c.saldo < cantidad {
        return fmt.Errorf("saldo insuficiente")
    }

    c.saldo -= cantidad
    fmt.Printf("Retirado: -%.2f, saldo: %.2f\n", cantidad, c.saldo)
    return nil
}

func (c *CuentaBancariaSol) Saldo() float64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.saldo
}

// Casos de prueba:
func TestEjercicio3() {
    cuenta := &CuentaBancariaSol{}

    // Test 1: Depósito simple
    err := cuenta.Depositar(100)
    if err != nil {
        fmt.Println("Error:", err)
    }

    saldo := cuenta.Saldo()
    if saldo != 100 {
        fmt.Printf("Error: saldo esperado 100, obtuve %.2f\n", saldo)
    } else {
        fmt.Println("✓ Test 1 pasado: Depósito correcto")
    }

    // Test 2: Retiro válido
    err = cuenta.Retirar(30)
    if err != nil {
        fmt.Println("Error:", err)
    }

    saldo = cuenta.Saldo()
    if saldo != 70 {
        fmt.Printf("Error: saldo esperado 70, obtuve %.2f\n", saldo)
    } else {
        fmt.Println("✓ Test 2 pasado: Retiro correcto")
    }

    // Test 3: Retiro con saldo insuficiente
    err = cuenta.Retirar(100)
    if err == nil {
        fmt.Println("Error: debería fallar con saldo insuficiente")
    } else {
        fmt.Println("✓ Test 3 pasado: Retiro rechazado correctamente")
    }

    // Test 4: Operaciones concurrentes
    fmt.Println("\n=== Test 4: Operaciones concurrentes ===")
    cuenta2 := &CuentaBancariaSol{}
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            cuenta2.Depositar(1)
        }()
    }

    wg.Wait()
    finalSaldo := cuenta2.Saldo()
    if finalSaldo != 100 {
        fmt.Printf("Error: saldo esperado 100, obtuve %.2f\n", finalSaldo)
    } else {
        fmt.Println("✓ Test 4 pasado: Concurrencia sin race conditions")
    }
}
```

---

### Ejercicio 4: Stack Trace y Logging

**Objetivo:** Usar defer para logging automático de entrada/salida y timing.

**Archivo:** `ejercicio4_logging.go`

```go
package main

import (
    "fmt"
    "time"
)

// Ejercicio 4: Crear decorator para logging y timing
//
// Requisitos:
// 1. Crear función Trace que retorne una función para defer
// 2. Loguear entrada con nombre de función y argumentos
// 3. Loguear salida con tiempo transcurrido
// 4. Registrar nivel de profundidad (indentación)
// 5. Mostrar stack trace visual
//
// Implementar:
var profundidad int

func Trace(nombre string) func() {
    // COMPLETAR
    return func() {}
}

// Solución esperada:
func TraceSol(nombre string) func() {
    profundidad++
    indentacion := ""
    for i := 0; i < profundidad-1; i++ {
        indentacion += "  "
    }

    inicio := time.Now()
    fmt.Printf("%s→ %s\n", indentacion, nombre)

    return func() {
        duracion := time.Since(inicio)
        fmt.Printf("%s← %s (%v)\n", indentacion, nombre, duracion)
        profundidad--
    }
}

// Funciones para demostración:
func funcionA() {
    defer TraceSol("A")()
    time.Sleep(10 * time.Millisecond)
    funcionB()
}

func funcionB() {
    defer TraceSol("B")()
    time.Sleep(20 * time.Millisecond)
    funcionC()
}

func funcionC() {
    defer TraceSol("C")()
    time.Sleep(15 * time.Millisecond)
}

// Casos de prueba:
func TestEjercicio4() {
    fmt.Println("=== Test 4: Stack Trace y Timing ===")
    profundidad = 0
    funcionA()
    fmt.Println("✓ Test 4: Trace de ejecución completado")
}
```

---

### Ejercicio 5: Resource Pool

**Objetivo:** Implementar un sistema de pool de recursos con adquisición/liberación automática.

**Archivo:** `ejercicio5_pool.go`

```go
package main

import (
    "fmt"
    "sync"
)

// Ejercicio 5: Pool de recursos con defer
//
// Requisitos:
// 1. Crear estructura Pool que gestione recursos
// 2. Implementar Adquirir() que saca un recurso del pool
// 3. Implementar Liberar(recurso) que devuelve al pool
// 4. Implementar WithRecurso(fn) que usa defer para liberación automática
// 5. Evitar race conditions
//
// Estructura base:
type Pool struct {
    mu       sync.Mutex
    libres   []interface{}
    en_uso   int
    maximo   int
    factory  func() interface{}
}

// Implementar métodos:
func (p *Pool) Adquirir() (interface{}, error) {
    // COMPLETAR
    return nil, nil
}

func (p *Pool) Liberar(recurso interface{}) {
    // COMPLETAR
}

func (p *Pool) WithRecurso(fn func(interface{}) error) error {
    // COMPLETAR
    return nil
}

// Solución esperada:
type PoolSol struct {
    mu       sync.Mutex
    libres   []interface{}
    en_uso   int
    maximo   int
    factory  func() interface{}
}

func NewPoolSol(maximo int, factory func() interface{}) *PoolSol {
    return &PoolSol{
        libres:  make([]interface{}, 0, maximo),
        maximo:  maximo,
        factory: factory,
    }
}

func (p *PoolSol) Adquirir() (interface{}, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if len(p.libres) > 0 {
        recurso := p.libres[len(p.libres)-1]
        p.libres = p.libres[:len(p.libres)-1]
        return recurso, nil
    }

    if p.en_uso >= p.maximo {
        return nil, fmt.Errorf("pool agotado")
    }

    p.en_uso++
    return p.factory(), nil
}

func (p *PoolSol) Liberar(recurso interface{}) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.libres = append(p.libres, recurso)
}

func (p *PoolSol) WithRecurso(fn func(interface{}) error) error {
    recurso, err := p.Adquirir()
    if err != nil {
        return err
    }

    defer p.Liberar(recurso)

    return fn(recurso)
}

// Casos de prueba:
func TestEjercicio5() {
    fmt.Println("=== Test 5: Resource Pool ===")

    contador := 0
    factory := func() interface{} {
        contador++
        return fmt.Sprintf("recurso_%d", contador)
    }

    pool := NewPoolSol(3, factory)

    // Test 1: Adquirir y liberar simple
    fmt.Println("\n--- Test 1: Adquisición simple ---")
    err := pool.WithRecurso(func(r interface{}) error {
        fmt.Printf("Usando: %v\n", r)
        return nil
    })

    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("✓ Test 1 pasado: Recurso adquirido y liberado")
    }

    // Test 2: Múltiples adquisiciones
    fmt.Println("\n--- Test 2: Múltiples adquisiciones ---")
    for i := 0; i < 3; i++ {
        err := pool.WithRecurso(func(r interface{}) error {
            fmt.Printf("Usando: %v\n", r)
            return nil
        })
        if err != nil {
            fmt.Println("Error:", err)
        }
    }
    fmt.Println("✓ Test 2 pasado: 3 recursos manejados")

    // Test 3: Pool agotado
    fmt.Println("\n--- Test 3: Pool agotado ---")

    // Adquirir todos
    recursos := make([]interface{}, 0)
    for i := 0; i < 3; i++ {
        r, err := pool.Adquirir()
        if err != nil {
            fmt.Println("Error:", err)
        }
        recursos = append(recursos, r)
    }

    // Intentar adquirir cuando está agotado
    _, err = pool.Adquirir()
    if err == nil {
        fmt.Println("Error: debería fallar cuando pool está agotado")
    } else {
        fmt.Println("✓ Test 3 pasado: Pool agotado correctamente")
    }

    // Liberar y probar disponibilidad
    for _, r := range recursos {
        pool.Liberar(r)
    }

    // Debería funcionar de nuevo
    _, err = pool.Adquirir()
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("✓ Test 3b pasado: Recurso disponible después de liberar")
    }
}

// Función main para ejecutar todos los tests
func RunAllTests() {
    fmt.Println("╔═══════════════════════════════════╗")
    fmt.Println("║     EJERCICIOS DE DEFER EN GO     ║")
    fmt.Println("╚═══════════════════════════════════╝")

    fmt.Println("\n" + "="*35)
    TestEjercicio1()

    fmt.Println("\n" + "="*35)
    TestEjercicio2()

    fmt.Println("\n" + "="*35)
    TestEjercicio3()

    fmt.Println("\n" + "="*35)
    TestEjercicio4()

    fmt.Println("\n" + "="*35)
    TestEjercicio5()

    fmt.Println("\n" + "="*35)
    fmt.Println("✓ Todos los ejercicios completados")
}
```

---

## Conclusión

El `defer` es un mecanismo fundamental en Go que simplifica enormemente el manejo de recursos y la garantía de limpieza. Key takeaways:

1. **Defer se ejecuta al salir de la función**, en orden LIFO, garantizado incluso con pánico.

2. **Usar inmediatamente después de adquirir recursos**, manteniendo el código relacionado junto.

3. **Evitar defer en loops** - usar funciones anidadas en su lugar.

4. **Argumentos se evalúan inmediatamente** - importante para capturar valores en el momento correcto.

5. **Closures capturan referencias**, útil para modificar estado o capturar contexto.

6. **Named returns permiten modificación desde defer**, patrón poderoso para enriquecimiento de errores.

7. **Patrones avanzados** como Ensure, Resource Guarding y Cleanup Chains generalizan casos de uso comunes.

El dominio de defer es esencial para escribir código Go robusto, seguro y limpio.

---

## Referencias y Lecturas Adicionales

- [Effective Go - Defer, Panic, and Recover](https://golang.org/doc/effective_go#defer)
- [Go Blog - Defer, Panic, and Recover](https://blog.golang.org/defer-panic-and-recover)
- [Go Specification - Defer Statement](https://golang.org/ref/spec#Defer_statements)

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/19-defer/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/19-defer):

```bash
cd examples/19-defer
go run .
```
