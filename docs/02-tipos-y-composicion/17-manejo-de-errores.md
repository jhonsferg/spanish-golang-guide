# Capítulo 17: Manejo de errores

## Índice del Capítulo 17

1. [17.1 ¿Por Qué Go No Tiene Excepciones?](#171-¿por-qué-go-no-tiene-excepciones)
2. [17.2 La Interface error](#172-la-interface-error)
3. [17.3 Errores Simples: fmt.Errorf y errors.New](#173-errores-simples-fmterrorf-y-errorsnew)
4. [17.4 Patrones de Retorno de Errores](#174-patrones-de-retorno-de-errores)
5. [17.5 Comprobación de Errores](#175-comprobación-de-errores)
6. [17.6 Tipos de Errores Personalizados](#176-tipos-de-errores-personalizados)
7. [17.7 Error Wrapping](#177-error-wrapping)
8. [17.8 Sentinel Errors](#178-sentinel-errors)
9. [17.9 Custom Error Types Avanzados](#179-custom-error-types-avanzados)
10. [17.10 Patrones Avanzados](#1710-patrones-avanzados)
11. [17.11 Buenas Prácticas y Antipatrones](#1711-buenas-prácticas-y-antipatrones)
12. [17.12 Ejercicios Progresivos](#1712-ejercicios-progresivos)

---

## 17.1 ¿Por Qué Go No Tiene Excepciones?

### El Problema Histórico con Excepciones

Go deliberadamente **rechaza el modelo de excepciones** que popularizaron Java, C++ y Python. Esto no fue una omisión accidental, sino una decisión filosófica consciente. Para entender por qué, primero necesitas entender qué está mal con las excepciones.

**El Problema 1: El Flujo Invisible**

```java
// Java - Código aparentemente simple
try {
    result = readFile();           // ¿Qué errores puede lanzar?
    processData(result);           // ¿Y este?
    writeFile(result);             // ¿Y este?
    uploadToCloud();               // ¿Cuál captura qué?
} catch (IOException e) {
    // Aquí cae ALGO, pero... ¿qué exactamente?
    logger.error("Error", e);
}
```

**El problema:** No sabes qué línea lanzó la excepción. El flujo de control es **implícito y oculto**. Una excepción puede "saltar" varias funciones, y nunca sabrás dónde se originó realmente. El código aparenta ser lineal, pero en realidad tiene decenas de "salidas invisibles".

**El Problema 2: Las Excepciones Marcan Casos "Excepcionales"**

Pero en desarrollo real, los errores NO son excepcionales:

- Archivo no encontrado: común, no es una excepción
- Conexión rechazada: normal en redes, no es una excepción
- Usuario ingresa datos inválidos: esperado, no es una excepción

Las excepciones pretenden ser para eventos "verdaderamente excepcionales" (división por cero, overflow, corrupción de memoria). Pero se usan para CONTROL DE FLUJO normal. Esto es un abuso del mecanismo.

**El Problema 3: Las Excepciones Son Costosas**

```go
// Go - Si usara excepciones (hipotético)
// throw 10,000 veces por segundo = problema de performance
```

En sistemas de alta concurrencia (como Docker, Kubernetes, etc.), el manejo de errores ocurre CONSTANTEMENTE. Usar excepciones sería catastrófico para la performance.

**El Problema 4: Overhead Cognitivo**

```cpp
// C++ - ¿Qué función es segura para exceptions?
class DataProcessor {
    void process() noexcept;        // Segura
    void validate() throw();        // Segura
    void loadData() throw(Error);   // Lanza Error
    void save();                    // ¿Lanza qué?
};
```

Los desarrolladores deben entender qué funciones lanzan qué excepciones. Sin especificación explícita, es un juego de adivinanzas. Esto añade complejidad mental.

### La Filosofía de Go: Errores Como Valores

Go toma una aproximación radicalmente diferente: **los errores son valores normales**.

```go
// Go - Los errores son valores
result, err := readFile("data.txt")

// Exactamente lo mismo que:
result, count := strings.Cut("a:b", ":")

// Los errores son tan ordinarios como cualquier otro valor
```

**Implicaciones profundas:**

1. **El flujo es explícito:** Sabes exactamente dónde se chequea cada error
2. **Los errores son ordinarios:** No se abusa del mecanismo de excepciones
3. **Es eficiente:** Sin overhead de stack unwinding
4. **El código es claro:** Leer el código es ver el flujo completo

### Comparativa: Go vs Otros Lenguajes

#### Go vs Java (Excepciones)

```java
// Java - Excepciones implícitas
try {
    for (String file : files) {
        processFile(file);  // ¿Qué pasa si falla?
    }
} catch (IOException e) {
    System.out.println("Error");  // Demasiado tarde para saber dónde
}
```

```go
// Go - Errores explícitos
for _, file := range files {
    if err := processFile(file); err != nil {
        // Sé exactamente dónde falló
        log.Printf("Error procesando %s: %v", file, err)
        continue  // Decisión deliberada: seguir o parar
    }
}
```

#### Go vs Rust (Result Type)

```rust
// Rust - Result type fuerza manejo
fn readFile(name: &str) -> Result<String, Error> {
    // Compiler fuerza que el caller maneje el Result
}

// El caller DEBE hacer esto:
match readFile("data.txt") {
    Ok(data) => process(data),
    Err(e) => println!("Error: {}", e),
}

// O puedes propagar con ?
let data = readFile("data.txt")?;  // Si falla, retorna Err
```

```go
// Go - Similar pero más simple
data, err := readFile("data.txt")
if err != nil {
    return err  // Propaga el error
}
process(data)
```

#### Go vs C (errno)

```c
// C - errno es un hack
FILE *f = fopen("file.txt", "r");
if (f == NULL) {
    // Verifica el valor de retorno
    perror("fopen");  // errno se perdió si otra función lo cambió
}
```

```go
// Go - El error viaja junto con el valor
f, err := os.Open("file.txt")
if err != nil {
    // El error está aquí, no en una variable global
    log.Fatal(err)
}
```

### Ventajas de la Filosofía de Go

| Aspecto | Go | Excepciones |
|--------|-----|------------|
| Flujo de Control | Explícito | Implícito |
| Performance | Alta (sin overhead) | Menor (stack unwinding) |
| Claridad | Código claro | Código oculta paths |
| Uso en Loops | Natural | Anti-patrón |
| Errores Ordinarios | Soportados | "Abuso" del sistema |
| Debugging | Fácil (ve dónde falla) | Difícil (encuentra la fuente) |

### Desventajas y Críticas

Es justo mencionar que este enfoque también tiene costos:

**1. Verbosidad:**

```go
if err != nil {
    return err
}
if err != nil {
    return err
}
if err != nil {
    return err
}
// Más adelante en este capítulo: patrones para evitar esto
```

**2. Omisión de Errores:**

```go
_ = writeToDatabase()  // Ignorar errores es MUY fácil
```

**3. Menos Información Automática:**

```go
// En Java, la excepción tiene stack trace automático
// En Go, tú debes agregarlo (veremos cómo)
```

### La Verdad: Es una Elección

Go no dice "las excepciones son malas". Dice: "para nuestro caso de uso (sistemas distribuidos de alta performance), los errores como valores son mejor".

Si construyes una aplicación tipo "CRUD web", Java con excepciones funciona bien. Si construyes un contenedor (Docker), un orquestador (Kubernetes), o un servicio de redes, Go es superior porque los errores son constantes y necesitan performance.

---

## 17.2 La Interface error

### Definición de la Interface error

En Go, `error` es una interface. La interface **más simple del lenguaje**:

```go
// Definida en builtin (no necesita import)
type error interface {
    Error() string
}
```

¡Eso es todo! Un type que implementa `error` solo necesita un método: `Error()` que retorna `string`.

### Lo que significa

```go
// Esto es un error
err := errors.New("algo salió mal")

// Porque errors.New() retorna una struct que implementa error:
// type errorString struct {
//     s string
// }
// func (e *errorString) Error() string { return e.s }

// Puedes verificar si algo es un error
var err error = errors.New("test")
if err != nil {
    fmt.Println(err.Error())  // "test"
    fmt.Println(err)          // "test" (Go llama a Error() automáticamente)
}
```

### nil Es Especial

```go
var err error = nil

// En Go, nil tiene significado especial:
if err != nil {
    // No se ejecuta
}

// Porque nil es el "valor cero" de interfaces
// Significa: "sin error"
```

### El Patrón Idiomático

```go
// El patrón que verás en TODO código Go:
result, err := SomeOperation()
if err != nil {
    // Algo salió mal
    return err
}
// Aquí result es seguro para usar
```

Este patrón es tan fundamental que el Go compiler tiene optimizaciones especiales para él.

### ¿Por qué `interface{}` en lugar de una type?

Porque una interface permite que **cualquier type que implemente `Error()` sea un error**. Esto es crucial:

```go
// Un error simple
type SimpleError string
func (e SimpleError) Error() string { return string(e) }

// Un error complejo con contexto
type APIError struct {
    StatusCode int
    Message    string
    Endpoint   string
}
func (e APIError) Error() string {
    return fmt.Sprintf("[%d] %s at %s", e.StatusCode, e.Message, e.Endpoint)
}

// Ambos son `error`
var err1 error = SimpleError("oops")
var err2 error = APIError{500, "Server Error", "/api/users"}
```

### nil vs non-nil

Este es un punto crítico que causa bugs:

```go
// INCORRECTO - Bug común
var myErr *CustomError
if myErr != nil {  // VERDADERO, aunque myErr es "nil"
    // Se ejecuta!!!
}

// ¿Por qué? Porque myErr es nil, pero el type pointer (*CustomError) no es nil
// La interface error sigue teniendo un value: (type *CustomError, value nil)

// CORRECTO
var myErr *CustomError = nil
var err error = myErr
if err != nil {  // Ahora es FALSO
    // No se ejecuta
}

// CORRECTO - Siempre retorna explícitamente nil
func SomeFunc() error {
    var myErr *CustomError
    if condition {
        myErr = &CustomError{}
    }
    return myErr  // Automáticamente es nil si myErr es nil
}
```

---

## 17.3 Errores Simples: fmt.Errorf y errors.New

### errors.New - La Forma Clásica

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    // Forma 1: errors.New
    err := errors.New("archivo no encontrado")
    fmt.Println(err)  // "archivo no encontrado"

    // Forma 2: El equivalente manual
    // err := &errorString{s: "archivo no encontrado"}
}

type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

**Uso típico:**

```go
func OpenFile(name string) ([]byte, error) {
    if name == "" {
        return nil, errors.New("nombre de archivo vacío")
    }
    // Código real de apertura...
    return data, nil
}
```

### fmt.Errorf - Errores con Formato

`fmt.Errorf` es como `fmt.Sprintf` pero retorna un `error`:

```go
package main

import (
    "fmt"
)

func main() {
    username := "admin"
    id := 42

    // Sin formato
    err1 := fmt.Errorf("usuario inválido")

    // Con formato (mucho mejor)
    err2 := fmt.Errorf("usuario %q no encontrado", username)
    err3 := fmt.Errorf("usuario id=%d no tiene permisos", id)

    fmt.Println(err2)  // "usuario "admin" no encontrado"
    fmt.Println(err3)  // "usuario id=42 no tiene permisos"
}
```

**Por qué es mejor que errors.New:**

```go
// errors.New - sin contexto
return nil, errors.New("error al procesar")

// fmt.Errorf - con contexto
return nil, fmt.Errorf("error al procesar archivo %q: %v", filename, err)

// El segundo es infinitamente más útil para debugging
```

### Convenios de Mensajes de Error

Go tiene convenciones ESPECÍFICAS sobre cómo escribir mensajes de error:

```go
// ✗ INCORRECTO - Comienza con mayúscula
fmt.Errorf("Error al abrir archivo")

// ✗ INCORRECTO - Con punto al final
fmt.Errorf("error al abrir archivo.")

// ✓ CORRECTO - Minúscula, sin punto
fmt.Errorf("error al abrir archivo")

// ✓ CORRECTO - Excepción: cuando empieza con una constante
fmt.Errorf("LegacyError: something")

// ✓ CORRECTO - Empieza con nombre de variable
fmt.Errorf("user %q not found", username)
```

**Razón:** Porque los errores a menudo se encadenan:

```go
err := readFile(name)
if err != nil {
    // Sin mayúscula, se ve natural:
    return fmt.Errorf("no se pudo procesar datos: %w", err)
    // Output: "no se pudo procesar datos: error al abrir archivo"
}
```

### Ejemplo Práctico: Validador Básico

```go
package main

import (
    "fmt"
    "strings"
)

func ValidateEmail(email string) error {
    if email == "" {
        return fmt.Errorf("email vacío")
    }
    if !strings.Contains(email, "@") {
        return fmt.Errorf("email %q inválido: falta @", email)
    }
    parts := strings.Split(email, "@")
    if len(parts[1]) == 0 {
        return fmt.Errorf("email %q inválido: dominio vacío", email)
    }
    return nil
}

func main() {
    tests := []string{
        "",
        "user@",
        "user@domain.com",
    }

    for _, email := range tests {
        if err := ValidateEmail(email); err != nil {
            fmt.Printf("✗ %s: %v\n", email, err)
        } else {
            fmt.Printf("✓ %s válido\n", email)
        }
    }
}
```

Output:

```
✗ : email vacío
✗ user@: email "user@" inválido: dominio vacío
✓ user@domain.com válido
```

---

## 17.4 Patrones de Retorno de Errores

### Patrón 1: Retorno Simple

```go
// La forma más directa
func ReadFile(name string) ([]byte, error) {
    if name == "" {
        return nil, fmt.Errorf("nombre vacío")
    }
    // Intenta leer
    data, err := os.ReadFile(name)
    if err != nil {
        return nil, fmt.Errorf("no se pudo leer %q: %w", name, err)
    }
    return data, nil
}

// Uso
data, err := ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}
```

### Patrón 2: Retorno Múltiple

```go
// Cuando hay múltiples valores y uno es error
func ParseConfig(data []byte) (Config, error) {
    var cfg Config

    if len(data) == 0 {
        return cfg, fmt.Errorf("datos vacíos")
    }

    if err := json.Unmarshal(data, &cfg); err != nil {
        return cfg, fmt.Errorf("JSON inválido: %w", err)
    }

    return cfg, nil
}
```

**Importante:** El orden es VALUE, ERROR. Siempre.

```go
// ✓ CORRECTO
value, err := SomeFunc()

// ✗ INCORRECTO (Go de verdad)
err, value := SomeFunc()  // Rompe el idiom
```

### Patrón 3: Named Returns (Retornos Nombrados)

```go
// Named returns permiten inicializar variables de retorno
func ProcessFile(name string) (data []byte, err error) {
    // 'data' y 'err' existen desde el inicio
    // data = nil, err = nil

    f, err := os.Open(name)
    if err != nil {
        return  // Retorna (nil, err) implícitamente
    }
    defer f.Close()

    data, err = io.ReadAll(f)
    if err != nil {
        return  // Retorna (nil, err)
    }

    return  // Retorna (data, nil)
}
```

**Ventajas:**

- `return` vacío es claro en contexto
- Las variables tienen significado

**Desventajas:**

- Menos evidente qué se retorna
- Puede confundir si hay operaciones intermedias

**Cuándo usarlo:**

```go
// ✓ Usa named returns para funciones simples y claras
func SimpleFunc() (value int, err error) {
    // Claro qué se retorna
}

// ✗ Evita named returns para funciones complejas
func ComplexFunc() (a int, b string, c bool, err error) {
    // Confuso
    // Mejor: return a, b, c, err
}
```

### Patrón 4: Early Return

Este es el idiom preferido de Go para evitar anidación:

```go
// ✗ INCORRECTO - Anidamiento profundo
func Process(file string) error {
    data, err := os.ReadFile(file)
    if err == nil {
        if json.Unmarshal(data, &config) == nil {
            if validateConfig(&config) {
                return saveConfig(&config)
            } else {
                return fmt.Errorf("config inválida")
            }
        } else {
            return fmt.Errorf("JSON inválido")
        }
    } else {
        return fmt.Errorf("no se pudo leer")
    }
}

// ✓ CORRECTO - Early return
func Process(file string) error {
    data, err := os.ReadFile(file)
    if err != nil {
        return fmt.Errorf("no se pudo leer: %w", err)
    }

    if err := json.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("JSON inválido: %w", err)
    }

    if !validateConfig(&config) {
        return fmt.Errorf("config inválida")
    }

    return saveConfig(&config)
}
```

### Patrón 5: Retorno Implícito en Defer

```go
func CopyFile(src, dst string) (err error) {
    srcFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("no se pudo abrir fuente: %w", err)
    }
    defer srcFile.Close()  // Se cierra incluso si hay error

    dstFile, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("no se pudo crear destino: %w", err)
    }
    defer func() {
        if err != nil {
            _ = os.Remove(dst)  // Limpia si hubo error
        }
        dstFile.Close()
    }()

    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return fmt.Errorf("error copiando: %w", err)
    }

    return nil
}
```

---

## 17.5 Comprobación de Errores

### if err != nil - El Patrón Fundamental

```go
// La forma más común en Go
result, err := SomeOperation()
if err != nil {
    // Maneja el error
}
```

Este patrón se ve en BILLONES de líneas de código Go. Es el idiom.

### Variaciones

**1. Solo ignorar (cuando es seguro):**

```go
// A veces sabes que no puede fallar
_ = ioutil.WriteFile("/tmp/log.txt", []byte("..."), 0644)

// Mejor: usar explícitamente
ioutil.WriteFile("/tmp/log.txt", []byte("..."), 0644)  // Go lo sabe que ignoras
```

**2. Early return con error:**

```go
data, err := readData()
if err != nil {
    return err  // Propaga el error hacia arriba
}
```

**3. Log y continuar:**

```go
data, err := readData()
if err != nil {
    log.Println("warning:", err)  // Registra pero continúa
}
```

**4. Valor por defecto:**

```go
data, err := readConfig()
if err != nil {
    data = DefaultConfig  // Usa un valor por defecto
}
```

**5. Tratamiento específico:**

```go
data, err := readFile(name)
if err != nil {
    // Aquí es donde quiero SABER QUÉ error es
    if errors.Is(err, os.ErrNotExist) {
        return fmt.Errorf("archivo no existe: %s", name)
    }
    if errors.Is(err, os.ErrPermission) {
        return fmt.Errorf("sin permisos para leer: %s", name)
    }
    return fmt.Errorf("error desconocido: %w", err)
}
```

### Antipatrones

**❌ Ignorar errores silenciosamente:**

```go
// Malo - se perderá un error silenciosamente
_ = WriteDatabase(user)

// Mejor - sé explícito
if err := WriteDatabase(user); err != nil {
    return err
}
```

**❌ Asumir que no puede fallar:**

```go
// "Esta función no puede fallar"... 6 meses después:
// - Se actualiza a una versión que sí falla
// - El linter te señala que ignoras errores
// - Debuggear por qué el flujo es roto

// Siempre maneja
if err != nil {
    return err
}
```

**❌ Log y continuar cuando deberías fallar:**

```go
// Incluso si loggueas, continuaste con datos inválidos
if err := parseConfig(data); err != nil {
    log.Println("Config error:", err)
    // Pero continúas con config = {}  ← Peligroso
}
```

### Best Practices

**1. Error lo más pronto posible:**

```go
// ✓ BIEN
if err != nil {
    return err
}

// ✗ MALO - Esperas para chequear
result, err := SomeOp()
doSomethingElse()
if err != nil {  // Tarde
    return err
}
```

**2. Enriquece el error con contexto:**

```go
// ✗ Pobre contexto
if err != nil {
    return err  // "file not found" - no sé cuál archivo
}

// ✓ Buen contexto
if err != nil {
    return fmt.Errorf("no se pudo leer config %q: %w", filename, err)
}
```

**3. Sé específico con error types:**

```go
// ✓ Permite tratamiento específico
if errors.Is(err, ErrNotFound) {
    // Retrying
} else if errors.Is(err, ErrPermission) {
    // Reportar al admin
} else {
    // Log genérico
}
```

---

## 17.6 Tipos de Errores Personalizados

### Struct que Implementa error

Cualquier struct con un método `Error()` es un error:

```go
package main

import (
    "fmt"
)

// Define tu tipo de error
type ValidationError struct {
    Field   string
    Message string
}

// Implementa la interface error
func (e ValidationError) Error() string {
    return fmt.Sprintf("validación falla en campo %q: %s", e.Field, e.Message)
}

func ValidateName(name string) error {
    if name == "" {
        return ValidationError{Field: "name", Message: "no puede estar vacío"}
    }
    if len(name) < 3 {
        return ValidationError{Field: "name", Message: "debe tener al menos 3 caracteres"}
    }
    return nil
}

func main() {
    if err := ValidateName("ab"); err != nil {
        fmt.Println(err)  // "validación falla en campo "name": debe tener al menos 3 caracteres"
    }
}
```

### Error con Código de Error

```go
package main

import (
    "fmt"
)

// Tipo de error con código
type ErrorCode int

const (
    ErrUnknown ErrorCode = iota
    ErrNotFound
    ErrInvalidInput
    ErrPermissionDenied
)

type APIError struct {
    Code    ErrorCode
    Message string
    Details string
}

func (e APIError) Error() string {
    return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
}

func (e APIError) StatusCode() int {
    switch e.Code {
    case ErrNotFound:
        return 404
    case ErrInvalidInput:
        return 400
    case ErrPermissionDenied:
        return 403
    default:
        return 500
    }
}

func GetUser(id int) error {
    if id <= 0 {
        return APIError{
            Code:    ErrInvalidInput,
            Message: "Invalid user ID",
            Details: fmt.Sprintf("ID debe ser positivo, recibió: %d", id),
        }
    }
    return nil
}

func main() {
    err := GetUser(-1)
    if err != nil {
        apierr := err.(APIError)
        fmt.Println("Status:", apierr.StatusCode())  // 400
        fmt.Println("Error:", err)
    }
}
```

### Error con Stack Trace

```go
package main

import (
    "fmt"
    "runtime"
)

type DetailedError struct {
    Message string
    Stack   []string
}

func (e DetailedError) Error() string {
    return e.Message
}

func (e DetailedError) StackTrace() {
    for i, frame := range e.Stack {
        fmt.Printf("%d. %s\n", i+1, frame)
    }
}

func captureStack() []string {
    var stack []string
    pc := make([]uintptr, 10)
    n := runtime.Callers(2, pc)
    frames := runtime.CallersFrames(pc[:n])

    for {
        frame, more := frames.Next()
        stack = append(stack, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
        if !more {
            break
        }
    }
    return stack
}

func SomeFunc() error {
    return DetailedError{
        Message: "algo salió mal",
        Stack:   captureStack(),
    }
}

func main() {
    if err := SomeFunc(); err != nil {
        derr := err.(DetailedError)
        fmt.Println(err)
        derr.StackTrace()
    }
}
```

### Error con Timeout

```go
package main

import (
    "fmt"
    "time"
)

type TimeoutError struct {
    Operation string
    Duration  time.Duration
}

func (e TimeoutError) Error() string {
    return fmt.Sprintf("operación %q excedió timeout (%v)", e.Operation, e.Duration)
}

func (e TimeoutError) Timeout() bool {
    return true  // Implementa interface net.Error
}

func (e TimeoutError) Temporary() bool {
    return true  // Es un error temporal, puedes reintentar
}

func FetchData(timeout time.Duration) error {
    // Simula timeout
    return TimeoutError{
        Operation: "FetchData",
        Duration:  timeout,
    }
}

func main() {
    if err := FetchData(5 * time.Second); err != nil {
        // Puedes verificar si es timeout
        if te, ok := err.(interface{ Timeout() bool }); ok && te.Timeout() {
            fmt.Println("Timeout detectado, reintentando...")
        }
    }
}
```

---

## 17.7 Error Wrapping

### El Problema: Pérdida de Contexto

Antes de Go 1.13, los errores perdían contexto:

```go
// Go 1.12 y anterior
data, err := os.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("error inesperado")  // Perdió el error original
}

// El usuario ve: "error inesperado"
// Pero ¿cuál fue el error original? ¿No encontrado? ¿Sin permisos?
```

### Solución: Error Wrapping con %w (Go 1.13+)

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    data, err := os.ReadFile("missing.txt")
    if err != nil {
        // ✓ CORRECTO - Envuelve el error preservando el original
        wrappedErr := fmt.Errorf("no se pudo leer config: %w", err)
        fmt.Println(wrappedErr)  // "no se pudo leer config: open missing.txt: no such file or directory"

        // El error original está dentro
        fmt.Println("Error original:", errors.Unwrap(wrappedErr))
    }
}
```

### Cadena de Errores (Error Chain)

Los errores pueden anidarse múltiples veces:

```go
package main

import (
    "errors"
    "fmt"
)

func Level3() error {
    return fmt.Errorf("database connection failed")
}

func Level2() error {
    err := Level3()
    return fmt.Errorf("level2 error: %w", err)
}

func Level1() error {
    err := Level2()
    return fmt.Errorf("level1 error: %w", err)
}

func main() {
    err := Level1()

    // Ver la cadena completa
    fmt.Println(err)
    // Output: level1 error: level2 error: database connection failed

    // Unwrap accede solo al siguiente nivel
    fmt.Println(errors.Unwrap(err))
    // Output: level2 error: database connection failed

    // Para ver la cadena completa
    for err != nil {
        fmt.Println("  ->", err)
        err = errors.Unwrap(err)
    }
}
```

### errors.Is - Comparación en Cadena

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    // Simula un error
    err := fmt.Errorf("error al procesar: %w", os.ErrNotExist)

    // ✗ INCORRECTO - No funciona si está envuelto
    if err == os.ErrNotExist {
        fmt.Println("No encontrado")  // No se ejecuta
    }

    // ✓ CORRECTO - Busca en la cadena
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("No encontrado")  // Se ejecuta
    }
}
```

### errors.As - Type Assertion en Cadena

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

type CustomError struct {
    Code int
    Msg  string
}

func (e CustomError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

func SomeOp() error {
    return fmt.Errorf("operación falló: %w",
        CustomError{Code: 42, Msg: "invalid input"})
}

func main() {
    err := SomeOp()

    // errors.As busca en la cadena un tipo específico
    var customErr CustomError
    if errors.As(err, &customErr) {
        fmt.Printf("Código: %d, Mensaje: %s\n", customErr.Code, customErr.Msg)
        // Output: Código: 42, Mensaje: invalid input
    }

    // Lo mismo con errores estándar
    var pathErr *os.PathError
    if errors.As(err, &pathErr) {
        fmt.Println("Error de ruta:", pathErr)
    }
}
```

### Diagrama: Error Wrapping

```
┌─────────────────────────────────────────────────────┐
│ Level 1: "aplicación no pudo procesar"             │
│   └── Wrapped error                                │
│       ┌───────────────────────────────────────────┐
│       │ Level 2: "no se pudo escribir BD"         │
│       │   └── Wrapped error                       │
│       │       ┌───────────────────────────────────┐
│       │       │ Level 3: "conexión rechazada"     │
│       │       └───────────────────────────────────┘
│       └───────────────────────────────────────────┘
└─────────────────────────────────────────────────────┘

errors.Unwrap(L1) -> L2
errors.Unwrap(L2) -> L3
errors.Unwrap(L3) -> nil

errors.Is(L1, L3) -> true (búsqueda en cadena)
```

---

## 17.8 Sentinel Errors

### Qué son Sentinel Errors

Son errores específicos que se comparan por identidad (==), no por valor.

```go
package main

import (
    "errors"
    "io"
)

// Ejemplos de sentinel errors del stdlib
// io.EOF
// io.ErrUnexpectedEOF
// os.ErrNotExist
// os.ErrPermission

var (
    ErrUserNotFound     = errors.New("usuario no encontrado")
    ErrInvalidPassword  = errors.New("contraseña inválida")
    ErrPermissionDenied = errors.New("permiso denegado")
)

func AuthenticateUser(username, password string) error {
    if username == "" {
        return ErrInvalidPassword
    }
    // ... búsqueda en BD
    return ErrUserNotFound
}

func main() {
    err := AuthenticateUser("", "pass")

    // Los sentinels se comparan con ==
    if err == ErrInvalidPassword {
        println("Contraseña incorrecta")
    } else if err == ErrUserNotFound {
        println("Usuario no existe")
    }
}
```

### Comparación: == vs errors.Is

```go
package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

func main() {
    // Cuando NO está envuelto, ambos funcionan
    err := ErrNotFound

    if err == ErrNotFound {
        println("✓ == funciona")
    }

    if errors.Is(err, ErrNotFound) {
        println("✓ errors.Is funciona")
    }

    // Cuando está envuelto, solo errors.Is funciona
    wrappedErr := fmt.Errorf("operación falló: %w", ErrNotFound)

    if wrappedErr == ErrNotFound {
        println("✗ == no funciona con wrapped")
    }

    if errors.Is(wrappedErr, ErrNotFound) {
        println("✓ errors.Is funciona incluso con wrapped")
    }
}
```

### Cuándo Usar Sentinels

**✓ USA cuando:**

- El error es bien definido y no necesita información variable
- Necesitas que el código cliente pueda identificar exactamente qué pasó
- El error representa un caso específico (no encontrado, timeout, etc.)

```go
var (
    ErrNotFound = errors.New("no encontrado")
    ErrTimeout  = errors.New("timeout")
)

if errors.Is(err, ErrTimeout) {
    // Reintentar
}
```

**✗ EVITA cuando:**

- El error necesita información variable (usa structs)
- El error se puede envolver (usa tipos personalizados)
- Hay muchos sentinels (difícil de mantener)

```go
// ✗ Demasiados sentinels
var (
    ErrCode400 = errors.New("error 400")
    ErrCode401 = errors.New("error 401")
    ErrCode403 = errors.New("error 403")
    // ... 100+ más
)

// ✓ Mejor: usar un tipo
type APIError struct {
    Code    int
    Message string
}
```

### Sentinels en Stdlib

Go usa sentinels en muchos lugares:

```go
import (
    "io"
    "os"
)

// io.EOF - Fin de archivo
// io.ErrUnexpectedEOF - EOF cuando se esperaba más
// io.ErrShortWrite - Write escribió menos de lo esperado

// os.ErrNotExist - Archivo no existe
// os.ErrPermission - Sin permisos
// os.ErrClosed - Archivo/conexión cerrada

// usage:
if errors.Is(err, io.EOF) {
    // Fin normal del archivo
}

if errors.Is(err, os.ErrPermission) {
    // Sin permisos para acceder
}
```

---

## 17.9 Custom Error Types Avanzados

### Error con Contexto Estructurado

```go
package main

import (
    "fmt"
    "time"
)

type OperationError struct {
    Op        string    // Operación que falló
    Path      string    // Qué recurso
    Err       error     // Error subyacente
    Timestamp time.Time // Cuándo
    RetryCount int      // Intentos
}

func (e *OperationError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("operación %q en %q falló: %v (intento %d)",
            e.Op, e.Path, e.Err, e.RetryCount)
    }
    return fmt.Sprintf("operación %q en %q falló (intento %d)",
        e.Op, e.Path, e.RetryCount)
}

func (e *OperationError) IsRetryable() bool {
    // Algunos errores se pueden reintentar
    switch e.Err {
    case ErrTemporaryFailure, ErrTimeout:
        return true
    default:
        return false
    }
}

func (e *OperationError) ElapsedTime() time.Duration {
    return time.Since(e.Timestamp)
}

var (
    ErrTemporaryFailure = fmt.Errorf("fallo temporal")
    ErrTimeout          = fmt.Errorf("timeout")
)

func ProcessWithRetry(path string, maxRetries int) error {
    var lastErr error

    for attempt := 0; attempt <= maxRetries; attempt++ {
        if err := process(path); err != nil {
            lastErr = err
            continue
        }
        return nil
    }

    return &OperationError{
        Op:        "ProcessWithRetry",
        Path:      path,
        Err:       lastErr,
        Timestamp: time.Now(),
        RetryCount: maxRetries,
    }
}

func process(path string) error {
    // Simulación
    return ErrTemporaryFailure
}

func main() {
    if err := ProcessWithRetry("data.txt", 3); err != nil {
        opErr := err.(*OperationError)
        fmt.Println(err)
        fmt.Printf("Retryable: %v\n", opErr.IsRetryable())
        fmt.Printf("Duración: %v\n", opErr.ElapsedTime())
    }
}
```

### Error con Métodos Adicionales

```go
package main

import (
    "fmt"
)

type ValidationErrors struct {
    Errors map[string]string  // field -> error message
}

func (ve ValidationErrors) Error() string {
    if len(ve.Errors) == 0 {
        return "validación pasada"
    }

    msg := "errores de validación:\n"
    for field, err := range ve.Errors {
        msg += fmt.Sprintf("  %s: %s\n", field, err)
    }
    return msg
}

func (ve ValidationErrors) IsValid() bool {
    return len(ve.Errors) == 0
}

func (ve ValidationErrors) Get(field string) string {
    return ve.Errors[field]
}

func (ve ValidationErrors) Add(field, message string) {
    ve.Errors[field] = message
}

func ValidateUser(username, email string) error {
    errors := ValidationErrors{Errors: make(map[string]string)}

    if username == "" {
        errors.Add("username", "no puede estar vacío")
    } else if len(username) < 3 {
        errors.Add("username", "mínimo 3 caracteres")
    }

    if email == "" {
        errors.Add("email", "no puede estar vacío")
    } else if !isValidEmail(email) {
        errors.Add("email", "formato inválido")
    }

    if !errors.IsValid() {
        return errors
    }
    return nil
}

func isValidEmail(email string) bool {
    return len(email) > 5  // Validación simplificada
}

func main() {
    err := ValidateUser("ab", "invalid")
    if err != nil {
        fmt.Println(err)

        // Acceder a errores específicos
        if ve, ok := err.(ValidationErrors); ok {
            if usernameErr := ve.Get("username"); usernameErr != "" {
                fmt.Printf("Error de username: %s\n", usernameErr)
            }
        }
    }
}
```

### Error con Logging Integrado

```go
package main

import (
    "fmt"
    "log"
    "os"
)

type LoggedError struct {
    Message string
    Level   string  // DEBUG, INFO, WARN, ERROR
    Logger  *log.Logger
}

func (e LoggedError) Error() string {
    return e.Message
}

func (e LoggedError) Log() {
    if e.Logger != nil {
        e.Logger.Printf("[%s] %s", e.Level, e.Message)
    }
}

func NewLoggedError(msg string, level string) LoggedError {
    return LoggedError{
        Message: msg,
        Level:   level,
        Logger:  log.New(os.Stderr, "", log.LstdFlags),
    }
}

func RiskyOperation() error {
    if shouldFail := true; shouldFail {
        return NewLoggedError("operación falló por razones", "ERROR")
    }
    return nil
}

func main() {
    if err := RiskyOperation(); err != nil {
        logErr := err.(LoggedError)
        logErr.Log()  // Automáticamente registra
    }
}
```

---

## 17.10 Patrones Avanzados

### Error Handler Pattern

```go
package main

import (
    "fmt"
    "log"
)

// ErrorHandler define cómo manejar diferentes tipos de errores
type ErrorHandler interface {
    Handle(err error) error
}

// LoggingHandler registra el error y retorna
type LoggingHandler struct {
    Logger *log.Logger
}

func (h LoggingHandler) Handle(err error) error {
    h.Logger.Printf("Error: %v", err)
    return err
}

// RetryHandler reintentar hasta N veces
type RetryHandler struct {
    MaxRetries int
    Handler    ErrorHandler
}

func (h RetryHandler) Handle(err error) error {
    for i := 0; i < h.MaxRetries; i++ {
        if err == nil {
            return nil
        }
        fmt.Printf("Reintentando (%d/%d)...\n", i+1, h.MaxRetries)
    }
    if h.Handler != nil {
        return h.Handler.Handle(err)
    }
    return err
}

// Uso
func main() {
    handler := RetryHandler{
        MaxRetries: 3,
        Handler: LoggingHandler{
            Logger: log.New(os.Stderr, "", 0),
        },
    }

    var err error = fmt.Errorf("algo falló")
    handler.Handle(err)
}
```

### Defer-Based Error Recovery

```go
package main

import (
    "fmt"
)

func SafeOperation(shouldFail bool) (result string, err error) {
    defer func() {
        if r := recover(); r != nil {
            result = ""
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    if shouldFail {
        panic("operación crítica falló")
    }

    return "éxito", nil
}

func main() {
    result, err := SafeOperation(true)
    if err != nil {
        fmt.Println("Error:", err)
        fmt.Println("Result:", result)
    }
}
```

### Context-Based Error Handling

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type ErrorWithContext struct {
    Err     error
    Context context.Context
    Op      string
}

func (e ErrorWithContext) Error() string {
    deadline, ok := e.Context.Deadline()
    if ok {
        return fmt.Sprintf("%s: %v (deadline: %v)", e.Op, e.Err, deadline)
    }
    return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func LongRunningOp(ctx context.Context) error {
    select {
    case <-time.After(2 * time.Second):
        return nil
    case <-ctx.Done():
        return ErrorWithContext{
            Err:     ctx.Err(),
            Context: ctx,
            Op:      "LongRunningOp",
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    if err := LongRunningOp(ctx); err != nil {
        fmt.Println(err)
    }
}
```

### Multiple Error Accumulation

```go
package main

import (
    "fmt"
    "sync"
)

type MultiError struct {
    Errors []error
    mu     sync.Mutex
}

func (me *MultiError) Add(err error) {
    if err != nil {
        me.mu.Lock()
        me.Errors = append(me.Errors, err)
        me.mu.Unlock()
    }
}

func (me *MultiError) Error() string {
    me.mu.Lock()
    defer me.mu.Unlock()

    if len(me.Errors) == 0 {
        return "sin errores"
    }

    msg := fmt.Sprintf("%d errores:\n", len(me.Errors))
    for i, err := range me.Errors {
        msg += fmt.Sprintf("  %d. %v\n", i+1, err)
    }
    return msg
}

func (me *MultiError) HasError() bool {
    me.mu.Lock()
    defer me.mu.Unlock()
    return len(me.Errors) > 0
}

// Uso con goroutines
func ProcessItems(items []string) error {
    multi := &MultiError{Errors: []error{}}
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        go func(i string) {
            defer wg.Done()
            if err := processItem(i); err != nil {
                multi.Add(err)
            }
        }(item)
    }

    wg.Wait()

    if multi.HasError() {
        return multi
    }
    return nil
}

func processItem(item string) error {
    return fmt.Errorf("error procesando %s", item)
}

func main() {
    if err := ProcessItems([]string{"a", "b", "c"}); err != nil {
        fmt.Println(err)
    }
}
```

---

## 17.11 Buenas Prácticas y Antipatrones

### ✓ BUENAS PRÁCTICAS

**1. Siempre maneja errores explícitamente**

```go
// ✓ BIEN
data, err := readData()
if err != nil {
    return err
}

// ✗ MAL
data, _ := readData()  // Ignoras el error
```

**2. Añade contexto cuando envuelves errores**

```go
// ✗ Pobre
return fmt.Errorf("error: %w", err)

// ✓ Mejor
return fmt.Errorf("no se pudo procesar archivo %q: %w", filename, err)
```

**3. Usa errores de tipo cuando necesitas información adicional**

```go
// ✗ No puedes extraer código HTTP
return fmt.Errorf("HTTP error")

// ✓ Puedes extraer información
type HTTPError struct {
    Status int
    Body   string
}

func (e HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.Status, e.Body)
}
```

**4. Usa errors.Is y errors.As para comparación**

```go
// ✗ No funciona con wrapped errors
if err == io.EOF {
    // ...
}

// ✓ Funciona siempre
if errors.Is(err, io.EOF) {
    // ...
}
```

**5. No captureses sin reintentos**

```go
// ✗ Capturas pero luego ignoras
if err := risky(); err != nil {
    // Sin hacer nada
}

// ✓ Sé explícito
if err := risky(); err != nil {
    // Registra y continúa deliberadamente
    log.Println("warning:", err)
}
```

### ✗ ANTIPATRONES

**Antipatrón 1: Ignorar Errores**

```go
// ✗ MALO
_ = database.Save(user)

// El usuario se guardó o no? Nunca sabes
// 6 meses después: "¿Por qué no se guardó mi usuario?"
```

**Antipatrón 2: Errores Sin Contexto**

```go
// ✗ MALO
return fmt.Errorf("error")

// ✗ PEOR
return fmt.Errorf("unknown error")

// ✓ BIEN
return fmt.Errorf("no se pudo guardar usuario %q en BD: %w", user.ID, err)
```

**Antipatrón 3: Comparación Directa de Errores Envueltos**

```go
// ✗ MALO - No funciona si err está envuelto
if err == os.ErrNotExist {
    // ...
}

// ✓ BIEN
if errors.Is(err, os.ErrNotExist) {
    // ...
}
```

**Antipatrón 4: Panic en Código Producción**

```go
// ✗ MALO - El servidor cae completamente
if err != nil {
    panic(err)  // En un handler HTTP: cliente no recibe respuesta
}

// ✓ BIEN - Maneja gracefully
if err != nil {
    http.Error(w, "Internal Server Error", 500)
    log.Printf("Error: %v", err)
}
```

**Antipatrón 5: Errores sin retorno**

```go
// ✗ MALO - El caller no sabe si falló
func process(data []byte) {
    if err := validate(data); err != nil {
        log.Printf("Validación falló: %v", err)
        // Pero continúa como si hubiera pasado
    }
}

// ✓ BIEN
func process(data []byte) error {
    if err := validate(data); err != nil {
        return fmt.Errorf("validación falló: %w", err)
    }
    return nil
}
```

**Antipatrón 6: Cadenas de errores infinitas**

```go
// ✗ MALO - Demasiada envoltura
fmt.Errorf("level1: %w",
    fmt.Errorf("level2: %w",
        fmt.Errorf("level3: %w",
            fmt.Errorf("level4: %w",
                originalError))))

// ✓ MEJOR - Envuelve solo donde agrega valor
fmt.Errorf("no se pudo procesar: %w", err)
```

### Checklist: Error Handling Review

```
☐ ¿Se manejan todos los errores?
☐ ¿Se añade contexto al envolver?
☐ ¿Se usa %w en lugar de %v para errores?
☐ ¿Se usa errors.Is/errors.As para comparación?
☐ ¿Hay sentinels bien documentados?
☐ ¿Los custom types implementan error correctamente?
☐ ¿Se logean errores en lugares apropiados?
☐ ¿No hay ignorancia deliberada sin razón?
☐ ¿Los mensajes son claros y útiles?
☐ ¿Se evitan panics en código de producción?
```

---

## 17.12 Ejercicios Progresivos

### Ejercicio 1: Error Personalizado con Código

**Objetivo:** Crear un tipo de error que incluya un código de error.

**Requisitos:**

1. Define `ValidationError` con campos: `Field`, `Code` (int), `Message`
2. Implementa `Error()` en formato: `[CODE] FIELD: MESSAGE`
3. Crea función `ValidateAge(age int) error` que:
   - Retorna código 400 si age < 0
   - Retorna código 401 si age > 150
   - Retorna nil si es válido

**Plantilla:**

```go
package main

import "fmt"

type ValidationError struct {
    Field   string
    Code    int
    Message string
}

func (e ValidationError) Error() string {
    // TODO: Implementar formato [CODE] FIELD: MESSAGE
}

func ValidateAge(age int) error {
    // TODO: Implementar validación
}

func main() {
    tests := []int{-5, 25, 200}
    for _, age := range tests {
        if err := ValidateAge(age); err != nil {
            fmt.Println(err)
        } else {
            fmt.Printf("Age %d válido\n", age)
        }
    }
}
```

**Salida esperada:**

```
[400] age: edad no puede ser negativa
[401] age: edad no puede ser mayor a 150
Age 25 válido
```

---

### Ejercicio 2: Validador Múltiple con Errores Agrupados

**Objetivo:** Crear un validador que retorna múltiples errores.

**Requisitos:**

1. Define `MultiFieldError` con slice de errors
2. Implementa `Error()` que lista todos
3. Implementa método `IsValid() bool`
4. Crea `ValidateUser(name, email, phone string) error` que:
   - Valida que name no esté vacío
   - Valida que email contenga "@"
   - Valida que phone tenga 10+ dígitos
   - Retorna MultiFieldError con TODOS los errores encontrados

**Plantilla:**

```go
package main

import (
    "fmt"
    "strings"
)

type MultiFieldError struct {
    Errors []string
}

func (me MultiFieldError) Error() string {
    // TODO: Listar todos los errores
}

func (me MultiFieldError) IsValid() bool {
    // TODO: Implementar
}

func ValidateUser(name, email, phone string) error {
    // TODO: Validar y retornar MultiFieldError
}

func main() {
    err := ValidateUser("", "invalidemail", "123")
    if err != nil {
        fmt.Println(err)
    }
}
```

**Salida esperada:**

```
3 errores de validación:
 - name: no puede estar vacío
 - email: debe contener @
 - phone: debe tener al menos 10 dígitos
```

---

### Ejercicio 3: Lector de Archivo con Error Wrapping

**Objetivo:** Leer un archivo con error wrapping apropiado.

**Requisitos:**

1. Crea función `ReadConfig(filename string) (map[string]string, error)`
2. Implementa wrapping de errores en 3 niveles:
   - Nivel 1: Abre archivo
   - Nivel 2: Lee contenido
   - Nivel 3: Parsea líneas "key=value"
3. Cada nivel envuelve el error anterior con contexto
4. En main, usa `errors.Is()` para detectar `os.ErrNotExist`

**Plantilla:**

```go
package main

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "strings"
)

func ReadConfig(filename string) (map[string]string, error) {
    // TODO: Implementar con wrapping en 3 niveles
}

func main() {
    cfg, err := ReadConfig("noexiste.conf")
    if err != nil {
        // Verificar si es "no existe"
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("Archivo no encontrado")
        }
        fmt.Println("Full error:", err)
    }
}
```

**Salida esperada:**

```
Archivo no encontrado
Full error: no se pudo leer config: no se pudo abrir config.conf: open config.conf: no such file or directory
```

---

### Ejercicio 4: Stack de Errores

**Objetivo:** Crear una estructura que acumula errores durante operación batch.

**Requisitos:**

1. Define `ErrorStack` con slice de errors y método `Add(err error)`
2. Implementa `Error()` que lista todos en orden
3. Implementa `IsEmpty() bool` y `Count() int`
4. Crea función `ProcessBatch(items []string)` que:
   - Intenta procesar cada item
   - Si falla, añade error al stack sin parar
   - Retorna ErrorStack si hay errores

**Plantilla:**

```go
package main

import (
    "fmt"
    "sync"
)

type ErrorStack struct {
    errors []error
    mu     sync.Mutex
}

func (es *ErrorStack) Add(err error) {
    // TODO: Implementar con protección mutex
}

func (es *ErrorStack) Error() string {
    // TODO: Listar todos los errores
}

func (es *ErrorStack) IsEmpty() bool {
    // TODO: Implementar
}

func ProcessBatch(items []string) error {
    stack := &ErrorStack{}
    for _, item := range items {
        if item == "fail" {
            stack.Add(fmt.Errorf("error procesando %q", item))
        }
    }

    if !stack.IsEmpty() {
        return stack
    }
    return nil
}

func main() {
    err := ProcessBatch([]string{"ok", "fail", "ok", "fail"})
    if err != nil {
        fmt.Println(err)
    }
}
```

**Salida esperada:**

```
2 errores:
 1. error procesando "fail"
 2. error procesando "fail"
```

---

### Ejercicio 5: Handler Middleware de Errores

**Objetivo:** Crear un middleware que intercepta y maneja errores.

**Requisitos:**

1. Define interfaz `ErrorHandler` con método `Handle(err error) error`
2. Crea 3 handlers:
   - `LoggingHandler`: registra el error
   - `RetryHandler`: reintentar N veces
   - `FallbackHandler`: proporciona valor por defecto
3. Crea función `ProcessWithHandlers(op func() error, handlers ...ErrorHandler) error`
4. Encadena handlers (uno llama al siguiente)

**Plantilla:**

```go
package main

import (
    "fmt"
    "log"
)

type ErrorHandler interface {
    Handle(err error, next ErrorHandler) error
}

type LoggingHandler struct{}

func (h LoggingHandler) Handle(err error, next ErrorHandler) error {
    if err != nil {
        log.Printf("Error: %v", err)
    }
    if next != nil {
        return next.Handle(err, nil)
    }
    return err
}

// TODO: Implementar RetryHandler y FallbackHandler

func ProcessWithHandlers(op func() error, handlers ...ErrorHandler) error {
    // TODO: Implementar encadenamiento
}

func main() {
    callCount := 0
    operation := func() error {
        callCount++
        if callCount < 3 {
            return fmt.Errorf("intento %d falló", callCount)
        }
        return nil
    }

    // Debería reintentar 2 veces antes de exitir
    if err := ProcessWithHandlers(operation); err != nil {
        fmt.Println("Final error:", err)
    }
}
```

**Salida esperada:**

```
Error: intento 1 falló
Error: intento 2 falló
Success en intento 3
```

---

## Resumen: La Filosofía de Go para Errores

Go rechaza las excepciones y las reemplaza con **errores como valores**. Esta decisión:

1. **Hace explícito el flujo:** Ves exactamente dónde se manejan los errores
2. **Es eficiente:** Sin overhead de stack unwinding
3. **Es ordinario:** Los errores son tan comunes que no merecen mecanismo especial

La interface `error` es simple: solo necesita un método `Error() string`. Todo lo demás que viste (wrapping, sentinels, tipos personalizados) son patrones construidos sobre esta base.

**Regla de Oro:** Si tienes que preguntar "¿y si hay un error aquí?", ese es tu indicador de que deberías manejarlo explícitamente.

---

## Recursos Adicionales

- [Error Handling in Go (effective Go)](https://golang.org/doc/effective_go#errors)
- [Go 1.13 Error Wrapping](https://golang.org/doc/go1.13#error_wrapping)
- [Dave Cheney - Don't just check errors](https://dave.cheney.net/2014/12/24/inspecting-errors)
- [Rob Pike - Errors are values](https://go.dev/blog/errors-are-values)

---

**Fin del Capítulo 17**

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/17-manejo-de-errores/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/17-manejo-de-errores):

```bash
cd examples/17-manejo-de-errores
go run .
```
