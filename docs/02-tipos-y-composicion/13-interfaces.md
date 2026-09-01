# Capítulo 13: Interfaces

## 13.1 ¿Qué es una Interface?

Una **interfaz** es un contrato que define un conjunto de métodos que un tipo debe implementar. Es el mecanismo fundamental de Go para lograr **polimorfismo** sin herencia de clases. Mientras que lenguajes como Java y C# usan jerarquías de herencia explícitas, Go adopta un enfoque más flexible y poderoso: **implementación implícita**.

### 13.1.1 El Concepto Fundamental

Una interfaz responde a una pregunta simple: "¿Qué puede hacer este tipo?" No "¿Qué es este tipo?". Esta diferencia filosófica es crucial para entender el diseño de Go.

```go
// En Java/C# se pregunta:
// Animal es un tipo base, Dog hereda de Animal
public class Animal { }
public class Dog extends Animal { }

// En Go se pregunta:
// ¿Qué tipos pueden hacer estas cosas?
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

### 13.1.2 Polimorfismo sin Herencia

El polimorfismo tradicional en lenguajes orientados a objetos se logra a través de herencia. Go logra lo mismo con **composición e interfaces**, sin herencia.

**Ventajas del enfoque de Go:**

1. **Desacoplamiento**: No necesitas conocer la jerarquía de tipos
2. **Flexibilidad**: Cualquier tipo existente puede "implementar" una interfaz sin modificación
3. **Composición sobre herencia**: Favorece la composición, que es más flexible
4. **No invasivo**: Las interfaces no se declaran en el tipo, se satisfacen automáticamente

### 13.1.3 ¿Por Qué Existen las Interfaces?

Las interfaces resuelven problemas concretos en programación:

- **Abstracción**: Separa la definición del comportamiento de la implementación
- **Testabilidad**: Permite mock objects y pruebas aisladas
- **Extensibilidad**: Nuevos tipos pueden satisfacer interfaces existentes sin tocar código anterior
- **APIs estables**: Define contratos públicos independientes de implementación

### 13.1.4 Casos de Uso Reales

```go
// 1. Manipulación de I/O genérico
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
    // Funciona con archivos, buffers, conexiones de red, etc.
    // Sin conocer el tipo específico
}

// 2. Serialización flexible
func Marshal(v interface{}) ([]byte, error) {
    // Maneja cualquier tipo
}

// 3. Patrones de diseño (Strategy, Adapter)
type Logger interface {
    Log(message string)
}

// Múltiples implementaciones: FileLogger, ConsoleLogger, RemoteLogger
```

---

## 13.2 Definición de Interfaces

### 13.2.1 Sintaxis Básica

Una interfaz en Go es declarada con la palabra clave `interface` y contiene un conjunto de firmas de método:

```go
type NombreInterface interface {
    Metodo1() TipoRetorno
    Metodo2(param TipoParam) (TipoRetorno, error)
    Metodo3()
}
```

**Características importantes:**

- Solo contiene **firmas de método**, nunca campos
- No tiene constructores ni destructores
- Puede estar vacía (`interface{}`)
- El orden de métodos es irrelevante para satisfacción

### 13.2.2 Ejemplo Progresivo

```go
package main

// Interfaz simple: un tipo que puede describirse a sí mismo
type Stringer interface {
    String() string
}

// Interfaz con múltiples métodos
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Interfaz con parámetros
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Interfaz con valores retorno múltiples
type JSONMarshaler interface {
    MarshalJSON() ([]byte, error)
}
```

### 13.2.3 Convenciones de Nombres

Go tiene convenciones específicas para nombres de interfaz:

- **Interfaz única**: Por lo general termina en `-er` si tiene un método (Reader, Writer, Closer)
- **Excepciones**: `io.ReaderWriter` es aceptable pero menos común
- **Descriptivas**: Usan sustantivos (Shape, Logger, Validator)
- **Minúsculas**: Comienzan con minúscula si son unexported

```go
// ✓ Bueno
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Validator interface {
    Validate() error
}

// ✗ Evitar
type DataReader interface {  // Demasiado genérico
    Read(p []byte) (n int, err error)
}

type IReader interface {  // Prefijo I innecesario
    Read(p []byte) (n int, err error)
}
```

### 13.2.4 Interfaz Estándar del Paquete fmt

Go incluye la interfaz `Stringer` en el paquete `fmt`:

```go
package fmt

type Stringer interface {
    String() string
}

// Cualquier tipo que implemente String() trabaja con fmt.Println
func Println(a ...interface{}) (n int, err error)
```

Ejemplo:

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d años)", p.Name, p.Age)
}

func main() {
    alice := Person{"Alice", 30}
    fmt.Println(alice)  // Usa automáticamente Person.String()
    // Output: Alice (30 años)
}
```

---

## 13.3 Implementación Implícita

### 13.3.1 El Paradigma: Implementación Automática

Esta es la característica que distingue más profundamente a Go de Java/C#. No declaras explícitamente que un tipo "implementa" una interfaz. La implementación es **implícita y automática**.

```go
// En Java (explícito):
public class MyWriter implements io.Writer {
    @Override
    public int write(byte[] p) throws IOException {
        // ...
    }
}

// En Go (implícito):
type MyWriter struct {}

func (w MyWriter) Write(p []byte) (int, error) {
    // ...
}

// MyWriter automáticamente satisface io.Writer
var w io.Writer = MyWriter{}  // ✓ Compila sin declaración explícita
```

### 13.3.2 Ventajas de la Implementación Implícita

**1. Desacoplamiento Total**

```go
// paquete legacy (antiguo, no sabía de Reader)
type FileBuffer struct {
    data []byte
}

func (fb FileBuffer) Read(p []byte) (int, error) {
    copy(p, fb.data)
    return len(fb.data), nil
}

// paquete nuevo (años después)
type io.Reader interface {
    Read(p []byte) (n int, err error)
}

// FileBuffer automáticamente satisface io.Reader
// Sin tocar código anterior
```

**2. Flexibilidad para Adaptadores**

```go
// Tipo externo que no puedes modificar
type ThirdPartyLogger struct {
    // campos privados
}

func (tpl ThirdPartyLogger) LogMessage(msg string) {
    // implementación
}

// Crear adaptador sin modificar el tipo
type LoggerAdapter struct {
    thirdParty ThirdPartyLogger
}

func (la LoggerAdapter) Log(msg string) {  // Cumple interface Logger
    la.thirdParty.LogMessage(msg)
}
```

**3. No Requiere Overhead de Herencia**

```go
// En lenguajes con herencia:
class A { }
class B extends A { }  // Acoplamiento implícito
class C extends B { }  // Jerarquía profunda

// En Go:
type A struct { }
func (a A) Do() { }

type B struct { }
func (b B) Do() { }  // Independiente de A

var x Doer = A{}  // A satisface Doer
var y Doer = B{}  // B también satisface Doer
```

### 13.3.3 Comparativa: Java vs Go

```go
// JAVA - Hierárquico y acoplado
public interface Animal {
    void makeSound();
}

public class Dog implements Animal {
    public void makeSound() {
        System.out.println("Woof!");
    }
}

public class Cat implements Animal {
    public void makeSound() {
        System.out.println("Meow!");
    }
}

// El tipo DEBE implementar la interfaz
// Si cambias la interfaz, debes cambiar todas las implementaciones
```

```go
// GO - Estructural y desacoplado
type Animal interface {
    MakeSound() string
}

type Dog struct{}
func (d Dog) MakeSound() string { return "Woof!" }

type Cat struct{}
func (c Cat) MakeSound() string { return "Meow!" }

// Cualquier tipo con MakeSound() automáticamente es un Animal
// Las definiciones son completamente independientes
```

---

## 13.4 Interfaces Satisfechas Automáticamente

### 13.4.1 El Compilador Verifica Automáticamente

Go verifica en **tiempo de compilación** si un tipo satisface una interfaz. Esto ocurre cuando:

1. El compilador ve una asignación a una variable de tipo interfaz
2. Verifica que todos los métodos de la interfaz existan en el tipo
3. Las firmas coincidan exactamente

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type File struct {
    name string
}

// ✗ Error: File no tiene método Read
var r Reader = File{}

// Agregar el método
func (f File) Read(p []byte) (int, error) {
    // implementación
    return 0, nil
}

// ✓ Ahora compila
var r Reader = File{}
```

### 13.4.2 Conjuntos de Métodos (Method Sets)

Cada tipo tiene un **conjunto de métodos** (method set). Una interfaz es satisfecha cuando el method set contiene todos los métodos de la interfaz.

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}

type Buffer struct {
    data []byte
}

func (b Buffer) Write(p []byte) (int, error) {
    b.data = append(b.data, p...)
    return len(p), nil
}

// Method set de Buffer: {Write}
// Method set de Writer: {Write}
// Buffer satisface Writer ✓

var w Writer = Buffer{}  // ✓ Compila
```

### 13.4.3 Receptores por Valor vs Referencia

El tipo de receptor afecta qué interfaz puede satisfacer:

```go
type Incrementer interface {
    Increment()
}

type Counter struct {
    value int
}

// Receptor por valor
func (c Counter) Increment() {
    c.value++  // Solo modifica la copia, no el original
}

// Receptor por referencia
func (c *Counter) Decrement() {
    c.value--  // Modifica el original
}

func main() {
    c := Counter{0}
    
    // ✓ Funciona: Counter tiene el método Increment()
    var inc Incrementer = c
    
    // ✗ No funciona: Counter no tiene Decrement(), solo *Counter
    // var dec Incrementer = c  // Error: Counter doesn't have Decrement()
    
    // ✓ Funciona: *Counter tiene ambos
    var dec Incrementer = &c
}
```

**Regla importante**: Si todos los métodos de una interfaz usan receptores por referencia, solo `*Tipo` satisface la interfaz, no `Tipo`.

### 13.4.4 Tabla Visual: Satisfacción de Interfaces

```
┌─────────────────────────────────────────────────────┐
│  Interface Reader                                   │
│  ├─ Read([]byte) (int, error)                       │
└─────────────────────────────────────────────────────┘
          ▲         ▲         ▲         ▲
          │         │         │         │
      File    Buffer    Socket    Connection
      ✓       ✓         ✓         ✓
      
  Todos tienen método Read() → Todos satisfacen Reader
  
  Go verifica automáticamente esto en compilación
```

---

## 13.5 Valores de Interface

### 13.5.1 ¿Cómo Almacena Go los Valores de Interface?

Internamente, una interfaz en Go es un par de dos palabras (dos punteros):

```
┌──────────────────────────────────────────┐
│         Interface Value (16 bytes)       │
├──────────────────────────────────────────┤
│  type *itab           │  data pointer    │
│  (type descriptor)    │  (valor real)    │
└──────────────────────────────────────────┘
```

Donde:
- **itab**: Descriptor de tipo (qué tipo concreto se almacena)
- **data pointer**: Puntero al valor real

```go
type Reader interface {
    Read(p []byte) (int, error)
}

type File struct { name string }
func (f File) Read(p []byte) (int, error) { /* ... */ }

func main() {
    f := File{"test.txt"}
    var r Reader = f
    
    // Internamente:
    // r.itab = descriptor de File
    // r.data = &f (puntero al valor File)
}
```

### 13.5.2 Nulidad y Validez

Una interfaz puede ser:

```go
type Writer interface {
    Write(p []byte) (int, error)
}

var w Writer
// w = nil (nil interface)
// w == nil → true

type File struct{}
func (f File) Write(p []byte) (int, error) { return len(p), nil }

f := File{}
w = f
// w != nil (interface contains File value)
// w == nil → false

var ptr *File = nil
w = ptr
// w != nil (interface contains *File, even though it points to nil)
// w == nil → false ← Sorpresa común

// Para comparar con nil en interfaz:
if w == nil || reflect.ValueOf(w).IsNil() {
    // Interface es realmente nula
}
```

### 13.5.3 Comparación de Interfaces

```go
type Comparable interface {
    Equal(other Comparable) bool
}

a1 := File{"a"}
a2 := File{"a"}
a3 := File{"b"}

// Comparación directa (si los valores son comparables)
i1 := Reader(a1)
i2 := Reader(a2)
i3 := Reader(a3)

i1 == i2  // true si File{"a"} == File{"a"}
i1 == i3  // false
```

### 13.5.4 Interfaz Vacía vs Interfaz Tipada

```go
var empty interface{}     // Almacena cualquier cosa
var reader io.Reader      // Almacena solo tipos que satisfacen Reader

// empty puede almacenar Reader
var r io.Reader
empty = r  // ✓

// reader no puede almacenar cualquier cosa
var something interface{} = 42
reader = something  // ✗ Error: cannot use interface{} as type Reader
```

---

## 13.6 Interfaces Vacías (interface{})

### 13.6.1 ¿Qué es interface{}?

`interface{}` es una interfaz **sin métodos**, lo que significa que **todos los tipos la satisfacen**. Es el equivalente a `Object` en Java o `Any` en otros lenguajes.

```go
type EmptyInterface interface {}

// Todos los tipos satisfacen interface{}
var x interface{}

x = 42              // int ✓
x = "hello"         // string ✓
x = []int{1, 2, 3}  // []int ✓
x = MyCustomType{}   // MyCustomType ✓
x = nil             // nil ✓
```

### 13.6.2 Cuándo Usar interface{}

**Casos apropiados:**

```go
// 1. Funciones realmente genéricas
func Println(a ...interface{}) (n int, err error) {
    // Acepta cualquier cantidad de valores de cualquier tipo
}

// 2. Almacenamiento genérico
type Database interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}) error
}

// 3. JSON donde no conoces la estructura
var data interface{}
json.Unmarshal(jsonBytes, &data)
// data podría ser map, []interface{}, string, number, etc.

// 4. APIs que requieren máxima flexibilidad
type Plugin interface {
    Execute(args ...interface{}) interface{}
}
```

### 13.6.3 Cuándo Evitar interface{}

**Anti-patrón: Sobrecarga de interface{}**

```go
// ✗ Mal: Pierde type safety
func Process(data interface{}) interface{} {
    // ¿Qué tipo se espera? Compilador no lo sabe
    // ¿Qué retorna? Misterio
    // Errores en tiempo de ejecución
}

// ✓ Bien: Usa interfaces específicas
func Process(reader io.Reader, writer io.Writer) error {
    // El compilador sabe exactamente qué hacer
    // Errores en tiempo de compilación
}
```

### 13.6.4 Rendimiento

`interface{}` tiene overhead:

```go
// Directamente con int
var x int = 42
x + 1  // Operación directa, sin indirección

// A través de interface{}
var empty interface{} = 42
result := empty.(int) + 1  // Type assertion + operación
// Más lento, requiere verificación de tipo

// En bucles críticos para rendimiento
for i := 0; i < 1000000; i++ {
    _ = interface{}(i)  // Evita esto
}
```

---

## 13.7 Type Assertions

### 13.7.1 Obtener el Valor Real

Cuando tienes un valor de interfaz, necesitas **type assertion** para acceder al tipo concreto:

```go
type Reader interface {
    Read(p []byte) (int, error)
}

var r Reader = File{}

// Type assertion: afirma que r contiene File
f := r.(File)  // Si r no es File, panic!

// Forma segura (comma-ok pattern)
f, ok := r.(File)
if !ok {
    fmt.Println("r no es File")
    return
}
// Usar f con seguridad
```

### 13.7.2 Syntax y Semántica

```go
// Acceso directo (inseguro)
value := interfaceValue.(ConcreteType)
// Si el tipo no coincide, panic en tiempo de ejecución

// Acceso seguro (comma-ok)
value, ok := interfaceValue.(ConcreteType)
if !ok {
    // Tipo incorrecto
}

// Con interface{}
var empty interface{} = "hello"

// Afirmar que es string
s, ok := empty.(string)
if ok {
    fmt.Println(s)  // hello
}

// Afirmar incorrectamente
n, ok := empty.(int)
// ok = false, n = 0, sin panic
```

### 13.7.3 Aplicaciones Prácticas

**Caso 1: Parsear JSON dinámico**

```go
var data interface{}
json.Unmarshal(jsonBytes, &data)

// Navegar estructura desconocida con seguridad
dataMap, ok := data.(map[string]interface{})
if ok {
    if name, ok := dataMap["name"].(string); ok {
        fmt.Println("Name:", name)
    }
}
```

**Caso 2: Logger con información de tipo**

```go
type Logger interface {
    Log(msg string, args ...interface{})
}

func LogWithType(logger Logger, value interface{}) {
    var typeInfo string
    
    switch v := value.(type) {
    case string:
        typeInfo = "string"
    case int, int64:
        typeInfo = "integer"
    case float64:
        typeInfo = "float"
    default:
        typeInfo = "unknown"
    }
    
    logger.Log("Value: %v (%s)", value, typeInfo)
}
```

### 13.7.4 Errores Comunes

```go
// ✗ Error: Type assertion en interfaz incorrecta
type Reader interface {
    Read([]byte) (int, error)
}

var r Reader = File{}
conn, ok := r.(io.Writer)  // Puede no ser Writer
if !ok {
    // No fue Writer
}

// ✗ Panic: Sin comma-ok
var empty interface{} = 42
s := empty.(string)  // panic: interface conversion

// ✓ Correcto: Siempre usar comma-ok para valores no garantizados
s, ok := empty.(string)
if !ok {
    // Manejar error
}
```

---

## 13.8 Type Switches

### 13.8.1 Múltiples Type Assertions

Cuando necesitas comprobar múltiples tipos, `type switch` es más limpio que múltiples `if`:

```go
// ✗ Repetitivo
var x interface{} = 42

_, isInt := x.(int)
_, isStr := x.(string)
_, isFloat := x.(float64)

if isInt {
    // ...
} else if isStr {
    // ...
} else if isFloat {
    // ...
}

// ✓ Limpio con type switch
switch v := x.(type) {
case int:
    fmt.Printf("int: %d\n", v)
case string:
    fmt.Printf("string: %s\n", v)
case float64:
    fmt.Printf("float: %f\n", v)
default:
    fmt.Printf("unknown: %T\n", x)
}
```

### 13.8.2 Sintaxis y Semantica

```go
switch x := value.(type) {
case Type1:
    // x es Type1
case Type2:
    // x es Type2
case Type3, Type4:
    // x es Type3 o Type4
    // Pero no puedes distinguir cuál
default:
    // x es el tipo original
}
```

### 13.8.3 Ejemplo: Procesador de Múltiples Formatos

```go
func ProcessData(data interface{}) string {
    switch v := data.(type) {
    case string:
        return "String: " + v
        
    case int:
        return fmt.Sprintf("Int: %d", v)
        
    case float64:
        return fmt.Sprintf("Float: %.2f", v)
        
    case []interface{}:
        count := len(v)
        return fmt.Sprintf("Array with %d elements", count)
        
    case map[string]interface{}:
        keys := len(v)
        return fmt.Sprintf("Object with %d keys", keys)
        
    case nil:
        return "Null"
        
    default:
        return fmt.Sprintf("Unknown type: %T", v)
    }
}

func main() {
    fmt.Println(ProcessData("hello"))           // String: hello
    fmt.Println(ProcessData(42))                // Int: 42
    fmt.Println(ProcessData(3.14))              // Float: 3.14
    fmt.Println(ProcessData([]interface{}{1,2})) // Array with 2 elements
}
```

### 13.8.4 Type Switch con Interfaces

```go
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type ReadWriter interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
}

func HandleStream(stream interface{}) string {
    switch s := stream.(type) {
    case ReadWriter:
        return "Can read and write"
    case Reader:
        return "Can only read"
    case Writer:
        return "Can only write"
    default:
        return "Not a stream"
    }
}
```

### 13.8.5 Patrones Comunes

**Convertir interface{} a tipos específicos:**

```go
func ConvertToString(value interface{}) string {
    switch v := value.(type) {
    case string:
        return v
    case int:
        return strconv.Itoa(v)
    case float64:
        return strconv.FormatFloat(v, 'f', -1, 64)
    case bool:
        return strconv.FormatBool(v)
    case []byte:
        return string(v)
    default:
        return fmt.Sprintf("%v", v)
    }
}
```

---

## 13.9 Embedding de Interfaces

### 13.9.1 Composición de Interfaces

Una interfaz puede incluir otras interfaces, creando interfaces más complejas:

```go
// Interfaz base 1
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Interfaz base 2
type Closer interface {
    Close() error
}

// Interfaz compuesta (embedding)
type ReadCloser interface {
    Reader
    Closer
}

// ReadCloser requiere ambos métodos:
// - Read(p []byte) (n int, err error)
// - Close() error
```

### 13.9.2 Ventajas de Embedding de Interfaces

```go
// ✗ Sin embedding: Repetitivo
type FileOps interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
    Seek(int64, int) (int64, error)
}

// ✓ Con embedding: Claro y mantenible
type io.Reader interface {
    Read([]byte) (int, error)
}

type io.Writer interface {
    Write([]byte) (int, error)
}

type io.Closer interface {
    Close() error
}

type io.Seeker interface {
    Seek(int64, int) (int64, error)
}

type FileOps interface {
    io.Reader
    io.Writer
    io.Closer
    io.Seeker
}
```

### 13.9.3 Jerarquía de Interfaces Estándar

```
io.Reader
├─ io.ReadCloser (Reader + Closer)
├─ io.ReadSeeker (Reader + Seeker)
└─ io.ReadWriter (Reader + Writer)

io.Writer
├─ io.WriteCloser (Writer + Closer)
└─ io.WriteSeeker (Writer + Seeker)

io.ReadWriter
├─ io.ReadWriteCloser (Reader + Writer + Closer)
└─ io.ReadWriteSeeker (Reader + Writer + Seeker)
```

### 13.9.4 Ejemplo Práctico

```go
package main

import "io"

// Interfaz base para repositorio
type DataStore interface {
    Save(key string, value []byte) error
    Load(key string) ([]byte, error)
}

// Interfaz para limpieza
type Closer interface {
    Close() error
}

// Interfaz compuesta
type PersistentStore interface {
    DataStore
    Closer
}

// Tipo concreto
type FileStore struct {
    dir string
}

func (fs FileStore) Save(key string, value []byte) error {
    // Implementación
    return nil
}

func (fs FileStore) Load(key string) ([]byte, error) {
    // Implementación
    return nil, nil
}

func (fs FileStore) Close() error {
    // Implementación
    return nil
}

// FileStore satisface PersistentStore
func main() {
    var store PersistentStore = FileStore{dir: "/data"}
    _ = store
}
```

---

## 13.10 Patrones Comunes

### 13.10.1 Patrón 1: Stringer

La interfaz `fmt.Stringer` es una de las más importantes:

```go
package fmt

type Stringer interface {
    String() string
}

// Útil para:
// 1. Imprimir valores personalizados
// 2. Logging
// 3. Depuración
```

**Ejemplo:**

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (age %d)", p.Name, p.Age)
}

func main() {
    alice := Person{"Alice", 30}
    fmt.Println(alice)           // Alice (age 30)
    fmt.Sprintf("%v", alice)     // Alice (age 30)
    log.Println(alice)           // Alice (age 30)
}
```

### 13.10.2 Patrón 2: Reader y Writer

El estándar de I/O genérico:

```go
package io

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Permite código genérico que funciona con:
// - Archivos
// - Buffers en memoria
// - Conexiones de red
// - Compresión
// - Encriptación
// - Etc.
```

**Ejemplo:**

```go
func CopyData(dst io.Writer, src io.Reader) (int64, error) {
    // Funciona con cualquier Reader/Writer
    return io.Copy(dst, src)
}

// Uso 1: Archivo a archivo
src, _ := os.Open("source.txt")
dst, _ := os.Create("dest.txt")
CopyData(dst, src)

// Uso 2: Memoria a archivo
buffer := bytes.Buffer{}
buffer.WriteString("hello")
dst, _ := os.Create("output.txt")
CopyData(dst, &buffer)

// Uso 3: Red a archivo
resp, _ := http.Get("https://example.com")
dst, _ := os.Create("download.html")
CopyData(dst, resp.Body)
```

### 13.10.3 Patrón 3: Marshaler

Para serialización personalizada:

```go
package encoding/json

type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```

**Ejemplo:**

```go
type Time struct {
    t time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
    // Serializar con formato personalizado
    return []byte(`"` + t.t.Format("2006-01-02") + `"`), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
    // Deserializar desde formato personalizado
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    parsed, err := time.Parse("2006-01-02", s)
    if err != nil {
        return err
    }
    t.t = parsed
    return nil
}
```

### 13.10.4 Patrón 4: Error

La interfaz más fundamental:

```go
package builtin

type error interface {
    Error() string
}

// Implementar error para tipos personalizados
type CustomError struct {
    Code    int
    Message string
}

func (e CustomError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func SomeFunction() error {
    return CustomError{500, "Internal error"}
}
```

### 13.10.5 Patrón 5: Predicado/Validator

```go
type Validator interface {
    Validate() error
}

type Config struct {
    Port int
    Host string
}

func (c Config) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }
    if c.Host == "" {
        return errors.New("host is required")
    }
    return nil
}

// Usar genéricamente
func Save(v Validator) error {
    if err := v.Validate(); err != nil {
        return err
    }
    // Guardar después de validación
    return nil
}
```

### 13.10.6 Patrón 6: Builder

```go
type QueryBuilder interface {
    Select(cols ...string) QueryBuilder
    Where(cond string) QueryBuilder
    OrderBy(field string) QueryBuilder
    Build() (string, error)
}

// Implementación que retorna la interfaz
type query struct {
    columns []string
    conds   []string
    orders  []string
}

func (q *query) Select(cols ...string) QueryBuilder {
    q.columns = append(q.columns, cols...)
    return q
}

func (q *query) Where(cond string) QueryBuilder {
    q.conds = append(q.conds, cond)
    return q
}

func (q *query) OrderBy(field string) QueryBuilder {
    q.orders = append(q.orders, field)
    return q
}

func (q *query) Build() (string, error) {
    // Construir query string
    return "SELECT ...", nil
}

// Uso fluido
func Query() QueryBuilder {
    return &query{}
}

func main() {
    sql, _ := Query().
        Select("id", "name").
        Where("age > 18").
        OrderBy("name").
        Build()
    fmt.Println(sql)
}
```

---

## 13.11 Buenas Prácticas

### 13.11.1 Interfaces Pequeñas

**Regla de Oro: Interfaces deben ser pequeñas**

```go
// ✗ Mala: Interfaz gigante, acoplada a detalles
type DataService interface {
    GetUser(id int) (*User, error)
    SaveUser(*User) error
    DeleteUser(id int) error
    GetPost(id int) (*Post, error)
    SavePost(*Post) error
    DeletePost(id int) error
    GetComments(postID int) ([]Comment, error)
    AddComment(Comment) error
    GetNotifications(userID int) ([]Notification, error)
    ClearNotifications(userID int) error
    AuthenticateUser(username, password string) (*User, error)
    LogActivity(action string) error
}

// ✓ Bien: Interfaces especializadas y reutilizables
type Repository interface {
    Get(id string) (interface{}, error)
    Save(item interface{}) error
    Delete(id string) error
}

type Authenticator interface {
    Authenticate(username, password string) (*User, error)
}

type Logger interface {
    Log(message string)
}
```

**Ventajas de interfaces pequeñas:**

1. Fáciles de implementar
2. Fáciles de testear
3. Más composables
4. Menos acoplamiento

### 13.11.2 Composición Sobre Herencia

```go
// ✗ Herencia simulada (acoplamiento fuerte)
type BaseService struct {
    db     *sql.DB
    logger Logger
}

type UserService struct {
    BaseService  // "Herencia" Go style
}

// ✓ Composición (acoplamiento débil)
type UserService struct {
    db     Repository
    logger Logger
}

func (us UserService) GetUser(id int) (*User, error) {
    return us.db.Get(id)
}
```

### 13.11.3 Depend on Interfaces, Not Implementations

```go
// ✗ Mal: Depende de implementación específica
func ProcessFile(file *os.File) error {
    // Solo funciona con *os.File
    data := make([]byte, 1024)
    n, err := file.Read(data)
    // ...
    return nil
}

// ✓ Bien: Depende de interfaz
func ProcessFile(reader io.Reader) error {
    // Funciona con cualquier Reader
    data := make([]byte, 1024)
    n, err := reader.Read(data)
    // ...
    return nil
}

// Ahora funciona con:
ProcessFile(os.File{})                      // Archivo
ProcessFile(bytes.Buffer{})                 // Buffer
ProcessFile(gzip.Reader{})                  // Comprimido
ProcessFile(crypto.Reader{})                // Encriptado
```

### 13.11.4 Aceptar Interfaces, Retornar Tipos Concretos

```go
// ✓ Patrón Go: Flexibilidad en entrada, control en salida
func NewServer(logger Logger) *Server {
    // Acepta cualquier Logger
    // Pero retorna *Server específico
}

func (s *Server) GetStatus() Status {
    // Retorna tipo concreto, no interfaz
    // Permite evitar allocations innecesarias
}

// ✗ Antipatrón: Retornar interfaz
func NewServer(logger Logger) Server {
    // Ineficiente: allocation, indirección
}
```

### 13.11.5 Evitar Interfaces Vacías Innecesarias

```go
// ✗ Antipatrón
func HandleData(data interface{}) {
    switch v := data.(type) {
    case string:
        // ...
    case int:
        // ...
    }
}

// ✓ Mejor: Interfaz específica
type DataHandler interface {
    Handle() error
}

func HandleData(handler DataHandler) error {
    return handler.Handle()
}
```

### 13.11.6 Naming Conventions

```go
// ✓ Buenas convenciones

// 1. Métodos simples: sufijo -er
type Reader interface { Read([]byte) (int, error) }
type Writer interface { Write([]byte) (int, error) }
type Closer interface { Close() error }

// 2. Predicados: prefijo Is- o Can-
type Validator interface { Validate() error }
type Comparable interface { Compare(other Comparable) int }

// 3. Sustantivos para comportamientos complejos
type Handler interface { Handle(interface{}) error }
type Dispatcher interface { Dispatch(interface{}) interface{} }

// ✗ Malas convenciones
type IReader interface { }        // Prefijo I innecesario
type ReadingInterface interface { }  // Sufijo Interface innecesario
type DOMManipulator interface { }   // Demasiado específico
```

### 13.11.7 Documentar Contratos

```go
// Reader es una interfaz para leer datos.
// Read lee hasta len(p) bytes del origen subyacente y los
// coloca en p. Retorna el número de bytes leídos (0 <= n <= len(p))
// y cualquier error encontrado.
// Si la fuente se agotó, n == 0 pero no hay error.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Validator verifica que un elemento es válido.
// Retorna nil si es válido, error describiendo problemas en caso contrario.
type Validator interface {
    Validate() error
}
```

### 13.11.8 Anti-patrones Comunes

**Anti-patrón 1: Type Assertion Innecesaria**

```go
// ✗ Mal
func Process(reader io.Reader) {
    file, ok := reader.(*os.File)
    if ok {
        // Procesar especialmente para File
    }
    // Procesar genéricamente
}

// ✓ Bien: Crear interfaz para comportamiento específico
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}

func Process(reader io.Reader) {
    if seeker, ok := reader.(Seeker); ok {
        // Usar Seek si disponible
    }
}
```

**Anti-patrón 2: Interfaz demasiado específica**

```go
// ✗ Mal: Interfaz para un único tipo
type MySQLDatabase interface {
    Query(sql string) (*sql.Rows, error)
}

// ✓ Bien: Interfaz genérica
type Database interface {
    Query(sql string) (*sql.Rows, error)
}
```

**Anti-patrón 3: Métodos que no se usan**

```go
// ✗ Mal: Interfaz con métodos que típicamente no se necesitan juntos
type HTTPServer interface {
    Start() error
    Stop() error
    ServeHTTP(w http.ResponseWriter, r *http.Request)
    GetMetrics() Metrics
    GetLogs() []string
    Restart() error
    ConfigReload() error
}

// ✓ Bien: Interfaces pequeñas y composables
type Server interface {
    Start() error
    Stop() error
}

type HTTPHandler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Metrics interface {
    GetMetrics() Metrics
}
```

---

## 13.12 Ejercicios Progresivos

### Ejercicio 1: Sistema de Formas Geométricas

**Nivel**: Básico
**Objetivo**: Implementar múltiples tipos que satisfacen una interfaz común

Crea una interfaz `Forma` con métodos `Area()` e `Perímetro()`. Implementa esta interfaz para:
- `Círculo` (radio)
- `Rectángulo` (ancho, alto)
- `Triángulo` (a, b, c)

Luego crea una función que acepte `[]Forma` y calcule:
- Área total
- Perímetro total
- Forma con mayor área

```go
package main

import (
    "fmt"
    "math"
)

// TODO: Define la interfaz Forma

// TODO: Define estructura Círculo con radio

// TODO: Implementa Área() para Círculo

// TODO: Implementa Perímetro() para Círculo

// TODO: Define estructura Rectángulo con ancho y alto

// TODO: Implementa Área() para Rectángulo

// TODO: Implementa Perímetro() para Rectángulo

// TODO: Define estructura Triángulo con lados a, b, c

// TODO: Implementa Área() para Triángulo (fórmula de Heron)

// TODO: Implementa Perímetro() para Triángulo

// TODO: Función para calcular área total

// TODO: Función para encontrar forma con mayor área

func main() {
    formas := []Forma{
        Círculo{radio: 5},
        Rectángulo{ancho: 4, alto: 3},
        Triángulo{a: 3, b: 4, c: 5},
    }
    
    fmt.Printf("Área total: %.2f\n", AreaTotal(formas))
    // Área total: 95.87
    
    mayor := FormaConMayorArea(formas)
    fmt.Printf("Mayor forma: %.2f\n", mayor.Área())
    // Mayor forma: 78.54
}
```

**Solución esperada:**

```go
package main

import (
    "fmt"
    "math"
)

type Forma interface {
    Area() float64
    Perimetro() float64
}

type Circulo struct {
    radio float64
}

func (c Circulo) Area() float64 {
    return math.Pi * c.radio * c.radio
}

func (c Circulo) Perimetro() float64 {
    return 2 * math.Pi * c.radio
}

type Rectangulo struct {
    ancho, alto float64
}

func (r Rectangulo) Area() float64 {
    return r.ancho * r.alto
}

func (r Rectangulo) Perimetro() float64 {
    return 2 * (r.ancho + r.alto)
}

type Triangulo struct {
    a, b, c float64
}

func (t Triangulo) Area() float64 {
    s := (t.a + t.b + t.c) / 2
    return math.Sqrt(s * (s - t.a) * (s - t.b) * (s - t.c))
}

func (t Triangulo) Perimetro() float64 {
    return t.a + t.b + t.c
}

func AreaTotal(formas []Forma) float64 {
    total := 0.0
    for _, f := range formas {
        total += f.Area()
    }
    return total
}

func FormaConMayorArea(formas []Forma) Forma {
    if len(formas) == 0 {
        return nil
    }
    mayor := formas[0]
    for _, f := range formas[1:] {
        if f.Area() > mayor.Area() {
            mayor = f
        }
    }
    return mayor
}

func main() {
    formas := []Forma{
        Circulo{radio: 5},
        Rectangulo{ancho: 4, alto: 3},
        Triangulo{a: 3, b: 4, c: 5},
    }

    fmt.Printf("Área total: %.2f\n", AreaTotal(formas))
    mayor := FormaConMayorArea(formas)
    fmt.Printf("Mayor forma: %.2f\n", mayor.Area())
}
```

---

### Ejercicio 2: Sistema de Notificaciones

**Nivel**: Intermedio
**Objetivo**: Implementar múltiples tipos de notificadores

Crea una interfaz `Notificador` con método `Enviar(destinatario, mensaje string) error`.

Implementa para:
- `EmailNotificador`
- `SMSNotificador`
- `PushNotificador`

Crea un `DistribuidorNotificaciones` que acepte múltiples notificadores y envíe a todos. Maneja errores parciales (si 1 falla, continúa con otros).

```go
package main

import (
    "fmt"
)

// TODO: Interfaz Notificador

// TODO: Tipo EmailNotificador

// TODO: Método Enviar para EmailNotificador

// TODO: Tipo SMSNotificador

// TODO: Método Enviar para SMSNotificador

// TODO: Tipo PushNotificador

// TODO: Método Enviar para PushNotificador

type DistribuidorNotificaciones struct {
    notificadores []Notificador
}

// TODO: Método para agregar notificadores

// TODO: Método EnviarATodos que retorna errores por notificador

func main() {
    dist := DistribuidorNotificaciones{}
    dist.Agregar(EmailNotificador{})
    dist.Agregar(SMSNotificador{})
    dist.Agregar(PushNotificador{})
    
    errores := dist.EnviarATodos("usuario@example.com", "Hola")
    
    if errores != nil {
        fmt.Println("Errores:", errores)
    }
}
```

---

### Ejercicio 3: Logger Extensible

**Nivel**: Intermedio
**Objetivo**: Crear sistema de logging con múltiples backends

Crea interfaz `Escritor` con método `Escribir(msg string) error`.

Implementa:
- `LoggerConsola`
- `LoggerArchivo`
- `LoggerRed` (simular cliente HTTP)

Crea `LoggerMultiplex` que escriba a múltiples destinos en paralelo usando goroutines.

```go
package main

import (
    "fmt"
    "os"
)

// TODO: Interfaz Escritor

// TODO: LoggerConsola

// TODO: LoggerArchivo (mantener file handle)

// TODO: LoggerRed (simular con buffer)

type LoggerMultiplex struct {
    escritores []Escritor
}

// TODO: Método Escribir que distribuye a todos los escritores

func main() {
    mult := LoggerMultiplex{}
    mult.AgregarEscritor(LoggerConsola{})
    mult.AgregarEscritor(LoggerArchivo{archivo: "app.log"})
    mult.AgregarEscritor(LoggerRed{endpoint: "http://logs.example.com"})
    
    err := mult.Escribir("Aplicación iniciada")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    }
}
```

---

### Ejercicio 4: Converter Genérico

**Nivel**: Avanzado
**Objetivo**: Usar type assertions con múltiples tipos

Crea interfaz `Convertible` con método `Convertir(formato string) (interface{}, error)`.

Implementa para diferentes tipos:
- `Documento` (a JSON, XML, PDF)
- `Imagen` (a PNG, JPG, WebP)
- `Tabla` (a CSV, JSON, HTML)

Crea función `ConvertirMuchosFormatos` que acepte `[]Convertible` y un formato, retornando resultados o errores.

```go
package main

import (
    "encoding/json"
    "fmt"
)

// TODO: Interfaz Convertible

// TODO: Estructura Documento

// TODO: Implementar Convertir para Documento

// TODO: Estructura Imagen

// TODO: Implementar Convertir para Imagen

// TODO: Estructura Tabla

// TODO: Implementar Convertir para Tabla

func ConvertirMuchosFormatos(items []Convertible, formato string) map[string]interface{} {
    resultados := make(map[string]interface{})
    
    for i, item := range items {
        result, err := item.Convertir(formato)
        // TODO: Guardar resultado o error
    }
    
    return resultados
}

func main() {
    items := []Convertible{
        Documento{titulo: "Reporte", contenido: "..."},
        Imagen{nombre: "foto.png"},
        Tabla{nombre: "datos", filas: 100},
    }
    
    resultados := ConvertirMuchosFormatos(items, "json")
    
    for clave, valor := range resultados {
        fmt.Printf("%s: %v\n", clave, valor)
    }
}
```

---

### Ejercicio 5: Sistema de Pagos

**Nivel**: Avanzado
**Objetivo**: Combinar múltiples técnicas (interfaces, type assertions, type switches)

Crea interfaz `Pagador` con método `Pagar(monto float64) error`.

Implementa:
- `TarjetaCredito` (simple)
- `Billetera` (requiere verificación de saldo)
- `Transferencia` (requiere código de autorización)

Crea función `ProcesarPago` que:
1. Acepta un `Pagador`
2. Calcula descuentos si el pagador es `Billetera` (-5%) o `Transferencia` (-10%)
3. Ejecuta el pago con monto ajustado
4. Registra en log qué tipo fue usado

```go
package main

import (
    "fmt"
)

// TODO: Interfaz Pagador

// TODO: TarjetaCredito

// TODO: Billetera

// TODO: Transferencia

type RegistroPago struct {
    Tipo        string
    MontoPagado float64
    Descuento   float64
    Exito       bool
    Error       string
}

func ProcesarPago(pagador Pagador, monto float64) RegistroPago {
    descuento := 0.0
    
    // TODO: Usar type switch para determinar descuento
    
    montoFinal := monto - (monto * descuento / 100)
    
    // TODO: Ejecutar pago
    
    // TODO: Retornar registro
    
    return RegistroPago{}
}

func main() {
    pagos := []Pagador{
        TarjetaCredito{numero: "1234-5678", cvv: "123"},
        Billetera{saldo: 1000},
        Transferencia{banco: "BanCo", cuenta: "123456"},
    }
    
    for _, pagador := range pagos {
        registro := ProcesarPago(pagador, 100)
        fmt.Printf("Pago: %+v\n", registro)
    }
}
```

**Salida esperada:**
```
Pago: {Tipo:TarjetaCredito MontoPagado:100 Descuento:0 Exito:true Error:}
Pago: {Tipo:Billetera MontoPagado:95 Descuento:5 Exito:true Error:}
Pago: {Tipo:Transferencia MontoPagado:90 Descuento:10 Exito:true Error:}
```

---

## 13.13 Resumen y Conclusión

Las interfaces son el corazón del diseño Go. Ofrecen:

1. **Polimorfismo** sin complejidad de herencia
2. **Desacoplamiento** a través de implementación implícita
3. **Flexibilidad** para adaptar y extender código
4. **Simplicidad** con interfaces pequeñas y enfocadas

**Puntos clave:**

- ✓ Mantén interfaces pequeñas (usualmente 1-3 métodos)
- ✓ Implementa implícitamente, sin declaraciones explícitas
- ✓ Acepta interfaces, retorna tipos concretos
- ✓ Usa type assertions y type switches con cuidado
- ✓ Compone interfaces para casos complejos
- ✓ Inspírate en interfaces estándar (io.Reader/Writer, fmt.Stringer)

Go no tiene "OOP" en sentido tradicional, pero logra polimorfismo real y elegante a través de interfaces. Dominar interfaces es dominar Go.

---

## Referencias y Lecturas Adicionales

- [Effective Go - Interfaces](https://golang.org/doc/effective_go#interfaces)
- [Go Code Review Comments - Interfaces](https://github.com/golang/go/wiki/CodeReviewComments#interfaces)
- [Standard Library: io package](https://pkg.go.dev/io)
- [Standard Library: fmt package](https://pkg.go.dev/fmt)
- [The Laws of Reflection](https://go.dev/blog/laws-of-reflection)


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/13-interfaces/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/13-interfaces):

```bash
cd examples/13-interfaces
go run .
```
