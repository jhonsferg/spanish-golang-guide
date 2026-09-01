# Capítulo 16: Punteros - Conceptos profundos y aritmética

## Índice del Capítulo 16

1. [16.1 ¿Qué es un Puntero?](#161-qué-es-un-puntero)
2. [16.2 Operadores & y *](#162-operadores--y-)
3. [16.3 Nil Pointers](#163-nil-pointers)
4. [16.4 Punteros a Structs](#164-punteros-a-structs)
5. [16.5 Punteros a Arrays y Slices](#165-punteros-a-arrays-y-slices)
6. [16.6 Punteros en Funciones](#166-punteros-en-funciones)
7. [16.7 Punteros en Métodos](#167-punteros-en-métodos)
8. [16.8 Punteros e Interfaces](#168-punteros-e-interfaces)
9. [16.9 Unsafe: Aritmética de Punteros](#169-unsafe-aritmética-de-punteros)
10. [16.10 Comparación de Punteros](#1610-comparación-de-punteros)
11. [16.11 Buenas Prácticas y Antipatrones](#1611-buenas-prácticas-y-antipatrones)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 16.1 ¿Qué es un Puntero?

### La Memoria en Go

Antes de entender punteros, necesitas entender cómo Go maneja la memoria:

```
┌──────────────────────────────────────────────────────┐
│           MEMORIA EN TIEMPO DE EJECUCIÓN            │
├──────────────────────────────────────────────────────┤
│                                                       │
│  STACK (Pila)              HEAP (Montículo)         │
│  ┌──────────────────┐    ┌──────────────────┐      │
│  │ Variables Locales│    │ Objetos Grandes  │      │
│  │ (rápido acceso)  │    │ (larga vida útil)│      │
│  │                  │    │                  │      │
│  │ x = 5            │    │ slice: [1,2,3]   │      │
│  │ p = 0x2540001    │    │ struct{...}      │      │
│  │ (dirección)      │    │ (GC supervisa)   │      │
│  └──────────────────┘    └──────────────────┘      │
│                                                       │
│  Automático      Limpieza automática (GC)            │
│                                                       │
└──────────────────────────────────────────────────────┘
```

### Definición: ¿Qué es un Puntero?

Un **puntero** es una variable que almacena la **dirección de memoria** de otra variable:

```go
var x int = 42
var p *int = &x  // p almacena la dirección de x

// Visualización:
// x: 0x2540001 → [42]
// p: 0x2540010 → [0x2540001]  (la dirección de x)
```

### Concepto: Valor vs Referencia

```
┌────────────────────────────────────────────────────┐
│   PASO POR VALOR (Copia)                           │
├────────────────────────────────────────────────────┤
│                                                     │
│  func cambiar(x int) {                             │
│    x = 100                                         │
│  }                                                  │
│                                                     │
│  a := 5                                            │
│  cambiar(a)                                        │
│  fmt.Println(a)  // 5 (sin cambios)                │
│                                                     │
│  ✓ Original protegido                              │
│  ✗ Copia innecesaria (variables grandes)           │
│                                                     │
└────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────┐
│   PASO POR REFERENCIA (Puntero)                    │
├────────────────────────────────────────────────────┤
│                                                     │
│  func cambiar(p *int) {                            │
│    *p = 100                                        │
│  }                                                  │
│                                                     │
│  a := 5                                            │
│  cambiar(&a)                                       │
│  fmt.Println(a)  // 100 (cambió)                   │
│                                                     │
│  ✓ Sin copias costosas                             │
│  ✓ Mutación permitida                              │
│  ✗ Debes ser consciente de los cambios             │
│                                                     │
└────────────────────────────────────────────────────┘
```

### Go vs C: Punteros Más Seguros

Go tiene varias restricciones sobre punteros comparado con C:

```
┌────────────────────┬──────────────────┬──────────────────┐
│ Característica     │ C                │ Go               │
├────────────────────┼──────────────────┼──────────────────┤
│ Aritmética         │ ✓ Libre          │ ✗ Solo unsafe    │
│ Punteros genéricos │ ✓ void*          │ ✗ unsafe.Pointer │
│ Dangling pointers  │ ✓ Posible        │ ✗ Garbage collect│
│ Seguridad          │ ✗ Baja           │ ✓ Alta           │
│ Nil checks         │ ✗ Manual         │ ✓ Type system    │
└────────────────────┴──────────────────┴──────────────────┘
```

**Beneficios de los punteros "seguros" de Go:**

- No hay punteros colgantes (dangling pointers)
- No hay aritmética de punteros accidental
- Garbage collection automático
- Nil checks son obvios en el tipo

### Razones para Usar Punteros

```go
// 1. MUTACIÓN: Modificar variable original
func incrementar(p *int) {
    *p++
}

// 2. EFICIENCIA: Evitar copias de estructuras grandes
type Usuario struct {
    nombre string
    email  string
    datos  []byte
}

func procesarUsuario(u *Usuario) {  // Puntero, no copia
    // ...
}

// 3. INTERFACES: Algunas interfaces requieren punteros
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 4. CICLOS DE REFERENCIA: Listas, árboles, grafos
type Nodo struct {
    valor int
    siguiente *Nodo  // Debe ser puntero
    anterior *Nodo
}
```

---

## 16.2 Operadores & y *

### El Operador & (Address Of)

El operador `&` obtiene la **dirección de memoria** de una variable:

```go
x := 42
p := &x

fmt.Printf("Valor de x: %v\n", x)      // 42
fmt.Printf("Dirección de x: %v\n", p)  // 0xc0001a050
fmt.Printf("Tipo de p: %T\n", p)       // *int
```

**Reglas importantes:**

```go
// 1. Solo puedes obtener dirección de variables (lvalues)
x := 5
p := &x     // ✓ Válido

p := &(5)   // ✗ Error: 5 es una constante, no variable
p := &42    // ✗ Error: 42 es un literal

// 2. La dirección es diferente cada ejecución
func main() {
    x := 5
    fmt.Println(&x)  // Ejecuta 1: 0xc0001a050
                      // Ejecuta 2: 0xc0001a100 (diferente)
}

// 3. & funciona con cualquier tipo
var edad int = 25
var peso float64 = 75.5
var nombre string = "Alice"

p1 := &edad      // *int
p2 := &peso      // *float64
p3 := &nombre    // *string
```

### El Operador * (Dereference)

El operador `*` **desreferencia** un puntero (accede al valor que apunta):

```go
x := 42
p := &x

fmt.Println(*p)   // 42 (el valor que p apunta)
*p = 100          // Cambiar el valor a través del puntero

fmt.Println(x)    // 100 (x cambió)
```

**Distintos significados de `*` según contexto:**

```go
// 1. En declaración de tipo: define un puntero
var p *int        // "p es un puntero a int"

// 2. En expresión: desreferencia
x := 5
p := &x
fmt.Println(*p)   // Desreferencia: obtiene 5

// 3. En multiplicación: es multiplicación
a := 5
b := 3
c := a * b        // c = 15
```

### Ejemplo Práctico: Swap de Valores

```go
func swap(a, b *int) {
    temp := *a
    *a = *b
    *b = temp
}

func main() {
    x := 10
    y := 20

    fmt.Printf("Antes: x=%d, y=%d\n", x, y)  // Antes: x=10, y=20
    swap(&x, &y)
    fmt.Printf("Después: x=%d, y=%d\n", x, y) // Después: x=20, y=10
}
```

### Diagrama de Memoria Paso a Paso

```
PASO 1: Crear variables
┌─────────────┬─────────────┐
│ x (addr ?)  │ y (addr ?)  │
│ valor: 10   │ valor: 20   │
└─────────────┴─────────────┘

PASO 2: Obtener direcciones
┌──────────────────────────────┐
│ swap(&x, &y)                 │
│ a recibe: dirección de x     │
│ b recibe: dirección de y     │
└──────────────────────────────┘

PASO 3: Desreferenciar y cambiar
temp := *a          // temp = 10
*a = *b             // x ahora es 20
*b = temp           // y ahora es 10

RESULTADO:
┌─────────────┬─────────────┐
│ x (id: ?)   │ y (id: ?)   │
│ valor: 20   │ valor: 10   │
└─────────────┴─────────────┘
```

### Declaración de Punteros: Sintaxis

```go
// Forma 1: Declaración simple
var p *int

// Forma 2: Con inicialización
var p *int = &x

// Forma 3: Inferencia (recomendado)
p := &x

// Forma 4: Múltiples punteros
var p1, p2 *int = &x, &y

// Forma 5: Puntero a puntero
pp := &p  // pp es *(*int)
fmt.Println(**pp)  // Desreferenciar dos veces para obtener valor
```

**Importante:** Un puntero sin inicializar es `nil`:

```go
var p *int
fmt.Println(p)      // <nil>
fmt.Println(p == nil) // true

// PELIGRO: Desreferenciar nil causa pánico
// fmt.Println(*p)  // panic: runtime error: invalid memory address
```

---

## 16.3 Nil Pointers

### ¿Qué es Nil?

En Go, `nil` es el **valor cero** de punteros, canales, funciones, interfaces y slices. Representa "sin valor".

```go
var p *int
var q *string
var fn func()
var ch chan int

fmt.Println(p == nil)   // true
fmt.Println(q == nil)   // true
fmt.Println(fn == nil)  // true
fmt.Println(ch == nil)  // true
```

### Comparación entre Nil y Cero

```
┌─────────────────────────────────────────────────────┐
│           Nil vs Valor Cero                         │
├─────────────────────────────────────────────────────┤
│                                                      │
│  var x int          // x = 0 (valor cero)          │
│  var p *int         // p = nil (sin puntero)       │
│  var s string       // s = "" (valor cero)         │
│  var ps *string     // ps = nil (sin puntero)      │
│                                                      │
│  ┌──────────────────┬────────────────────────────┐ │
│  │ Tipo             │ Valor Nil/Cero             │ │
│  ├──────────────────┼────────────────────────────┤ │
│  │ *T (puntero)     │ nil                        │ │
│  │ []T (slice)      │ nil                        │ │
│  │ map[K]V (mapa)   │ nil                        │ │
│  │ chan T (canal)   │ nil                        │ │
│  │ func (función)   │ nil                        │ │
│  │ interface{}      │ nil                        │ │
│  │ int              │ 0                          │ │
│  │ string           │ ""                         │ │
│  │ bool             │ false                      │ │
│  └──────────────────┴────────────────────────────┘ │
│                                                      │
└─────────────────────────────────────────────────────┘
```

### Nil Checks: Seguridad de Punteros

**Problema:** Desreferenciar un puntero nil causa pánico:

```go
var p *int
fmt.Println(*p)  // panic: runtime error: invalid memory address
```

**Solución:** Siempre verificar nil antes de usar:

```go
func procesar(p *int) error {
    if p == nil {
        return errors.New("puntero es nil")
    }

    fmt.Println(*p)  // Seguro ahora
    return nil
}
```

### Patrón de Validación Común

```go
type Usuario struct {
    nombre string
    email  string
}

func mostrarUsuario(u *Usuario) {
    if u == nil {
        fmt.Println("Usuario es nil")
        return
    }

    fmt.Printf("Usuario: %s (%s)\n", u.nombre, u.email)
}

func main() {
    var u *Usuario        // nil
    mostrarUsuario(u)     // Imprime: Usuario es nil

    u = &Usuario{
        nombre: "Alice",
        email:  "alice@example.com",
    }
    mostrarUsuario(u)     // Imprime: Usuario: Alice (alice@example.com)
}
```

### Nil Slices vs Nil Pointers

Un aspecto confuso: los slices pueden ser nil sin ser punteros:

```go
var s []int          // nil slice
var p *[]int         // nil pointer to slice

fmt.Println(s == nil)    // true
fmt.Println(len(s))      // 0
fmt.Println(cap(s))      // 0

if s != nil {            // ✓ Seguro con nil
    s = append(s, 1)     // ✓ Funciona aunque sea nil
}

// Nil slices son seguros en muchas operaciones
```

### Inicialización Segura

```go
// ANTI-PATRÓN: Esto causa pánico
var p *int
var m map[string]string
if p != nil {
    *p = 42  // OK si p no es nil
}
if m != nil {
    m["key"] = "value"  // OK
}

// PATRÓN CORRECTO: Inicializar antes de usar
p := new(int)           // *p = 0
*p = 42                 // Seguro

m := make(map[string]string)
m["key"] = "value"      // Seguro

// O con valores iniciales
var s []int = make([]int, 0)
s = append(s, 1)        // Seguro
```

### Funciones que Retornan Punteros

A menudo devuelves un puntero o error:

```go
// Patrón: Puntero + error
func encontrarUsuario(id int) (*Usuario, error) {
    if id < 0 {
        return nil, errors.New("ID inválido")
    }

    // Buscar en base de datos
    // ...

    if noEncontrado {
        return nil, errors.New("usuario no encontrado")
    }

    return &Usuario{id: id, nombre: "Alice"}, nil
}

// Uso seguro
u, err := encontrarUsuario(1)
if err != nil {
    log.Fatal(err)
}

fmt.Println(u.nombre)  // Alice
```

---

## 16.4 Punteros a Structs

### Acceso a Campos a Través de Punteros

En Go, acceder a campos de una struct a través de un puntero es automático:

```go
type Persona struct {
    nombre string
    edad   int
}

p := &Persona{nombre: "Alice", edad: 30}

// Ambas formas son equivalentes:
fmt.Println((*p).nombre)  // Desreferenciar explícito
fmt.Println(p.nombre)     // Desreferenciar implícito (azúcar sintáctico)

// Modificar también funciona
(*p).edad = 31            // Explícito
p.edad = 31               // Implícito (recomendado)
```

### El Desazúcar Automático

Go automáticamente desreferencia punteros a structs:

```go
type Punto struct {
    x, y float64
}

pt := Punto{1, 2}
pp := &pt

// Estos son equivalentes:
_ = (*pp).x + (*pp).y       // Manual
_ = pp.x + pp.y             // Automático (Go lo desazúcara)

// Mismo para modificación:
(*pp).x = 10                // Manual
pp.x = 10                   // Automático
```

**¿Por qué esto es importante?**

Hace el código más limpio y legible. Sin esta característica, todos los accesos a structs a través de punteros sería tedioso.

### Inicialización de Structs con Punteros

```go
type Config struct {
    puerto    int
    host      string
    debug     bool
}

// Forma 1: Crear en stack, luego tomar dirección
c := Config{puerto: 8080, host: "localhost"}
p := &c

// Forma 2: Crear directamente en heap con &{}
p := &Config{
    puerto: 8080,
    host:   "localhost",
    debug:  true,
}

// Forma 3: Usar new() - todos los campos a cero
p := new(Config)
p.puerto = 8080           // ✓ Desazúcara automáticamente
p.host = "localhost"
```

### Ejemplos Prácticos

```go
// Ejemplo 1: Modificar en función
type Usuario struct {
    nombre string
    email  string
}

func actualizarEmail(u *Usuario, nuevoEmail string) {
    u.email = nuevoEmail
}

func main() {
    u := &Usuario{nombre: "Alice", email: "alice@old.com"}
    actualizarEmail(u, "alice@new.com")
    fmt.Println(u.email)  // alice@new.com
}

// Ejemplo 2: Métodos mutadores
func (u *Usuario) establecerNombre(nombre string) {
    u.nombre = nombre
}

// Ejemplo 3: Validación antes de modificar
func (c *Config) establecerPuerto(p int) error {
    if p < 1 || p > 65535 {
        return errors.New("puerto inválido")
    }
    c.puerto = p
    return nil
}
```

### Structs Anidadas con Punteros

```go
type Dirección struct {
    calle string
    ciudad string
}

type Persona struct {
    nombre   string
    dirección *Dirección  // Puntero a struct anidada
}

p := &Persona{
    nombre: "Alice",
    dirección: &Dirección{
        calle: "Calle Principal 123",
        ciudad: "Madrid",
    },
}

// Acceso multicapa
fmt.Println(p.dirección.ciudad)  // Madrid
```

### Nil Structs Pointer

```go
var p *Persona  // nil

if p != nil {
    fmt.Println(p.nombre)  // Seguro, no se ejecuta
}

// Esto causaría pánico:
// fmt.Println(p.nombre)  // panic: nil pointer dereference
```

---

## 16.5 Punteros a Arrays y Slices

### Puntero a Array vs Slice

Una distinción importante:

```
┌──────────────────────────────────────────────────┐
│  PUNTERO A ARRAY         vs     SLICE            │
├──────────────────────────────────────────────────┤
│                                                   │
│  []*int = puntero a   vs   []int = slice (ya     │
│          array estático        incluye puntero)  │
│                                                   │
│  *[5]int = puntero a  vs   []int = No necesitas  │
│         array de 5           el puntero          │
│                                                   │
└──────────────────────────────────────────────────┘
```

### Puntero a Array

```go
// Array estático
arr := [5]int{1, 2, 3, 4, 5}
p := &arr  // *[5]int

// Acceso automático a elementos
fmt.Println(p[0])      // 1 (desazúcara a (*p)[0])
fmt.Println(p[1])      // 2

// Modificación
p[0] = 100
fmt.Println(arr[0])    // 100 (el array original cambió)

// Iteración
for i, v := range p {  // p funciona en range
    fmt.Println(i, v)
}

// Conversión a slice
s := p[:]              // Convierte a slice
```

### Slices (Ya Son Punteros Internamente)

Los slices ya contienen un puntero internamente, así que **casi nunca necesitas *[]int**:

```go
s := []int{1, 2, 3, 4, 5}

// El slice ya es eficiente: contiene
// - Puntero al array subyacente
// - Longitud
// - Capacidad

// Pasar a función (ya es por referencia interna)
func modificarSlice(s []int) {
    s[0] = 100  // Modifica el slice original
}

modificarSlice(s)
fmt.Println(s[0])      // 100

// ANTI-PATRÓN: Evitar *[]int
func malo(ps *[]int) {         // ✗ Raro y confuso
    (*ps)[0] = 100
}

// PATRÓN CORRECTO: Solo []int
func correcto(s []int) {       // ✓ Claro
    s[0] = 100
}
```

### ¿Cuándo Usar Puntero a Array?

Usar puntero a array cuando:

1. **Necesitas cambiar la longitud del array** (aunque raro):

```go
func intercambiarArrays(p1, p2 *[3]int) {
    temp := p1
    p1 = p2  // Esto solo cambia el puntero local
    p2 = temp
    // No afecta a los arrays originales
}
```

2. **Pasas arrays muy grandes y quieres evitar copias**:

```go
// Forma 1: Copia el array completo
func procesarArray(arr [1000000]int) {}   // ✗ Copia

// Forma 2: Puntero, sin copias
func procesarArray(p *[1000000]int) {}    // ✓ Eficiente

// Forma 3: Lo mejor - usa slice
func procesarArray(s []int) {}            // ✓ Mejor aún
```

3. **Necesitas comparar arrays** (arrays grandes):

```go
arr1 := [5]int{1, 2, 3, 4, 5}
arr2 := [5]int{1, 2, 3, 4, 5}

if arr1 == arr2 {  // ✓ Compara contenido
    fmt.Println("Iguales")
}

// Con punteros (raro, pero posible)
if &arr1 == &arr2 {  // ✗ Nunca será true (direcciones diferentes)
    fmt.Println("Mismo array en memoria")
}
```

### Ejemplo: Modificar Array en Función

```go
// Forma 1: Puntero a array (tradicional)
func duplicarElementos(p *[5]int) {
    for i := range p {
        p[i] *= 2
    }
}

// Forma 2: Slice (moderno, recomendado)
func duplicarElementos(s []int) {
    for i := range s {
        s[i] *= 2
    }
}

func main() {
    // Con puntero a array
    arr := [5]int{1, 2, 3, 4, 5}
    duplicarElementos(&arr)
    fmt.Println(arr)  // [2 4 6 8 10]

    // Con slice
    s := []int{1, 2, 3, 4, 5}
    duplicarElementos(s)
    fmt.Println(s)    // [2 4 6 8 10]
}
```

---

## 16.6 Punteros en Funciones

### Paso por Referencia: Mutación

La razón principal para usar punteros en funciones es permitir **mutación**:

```go
// ANTI-PATRÓN: Paso por valor (copia)
func incrementarSinPuntero(x int) {
    x++  // Modifica la copia, no el original
}

func main() {
    a := 5
    incrementarSinPuntero(a)
    fmt.Println(a)  // 5 (sin cambios)
}

// PATRÓN CORRECTO: Paso por referencia (puntero)
func incrementarConPuntero(p *int) {
    (*p)++  // Modifica el original a través del puntero
}

func main() {
    a := 5
    incrementarConPuntero(&a)
    fmt.Println(a)  // 6 (cambió)
}
```

### Patrón: Inicialización de Estructuras

```go
type Configuración struct {
    puerto     int
    host       string
    maxClientes int
}

// Anti-patrón: Retornar struct (copia)
func crearConfig() Configuración {
    return Configuración{
        puerto: 8080,
        host: "localhost",
        maxClientes: 100,
    }
}

// Patrón correcto: Inicializar a través de puntero
func inicializarConfig(c *Configuración) error {
    if c == nil {
        return errors.New("configuración nil")
    }
    c.puerto = 8080
    c.host = "localhost"
    c.maxClientes = 100
    return nil
}

// Uso
var config Configuración
if err := inicializarConfig(&config); err != nil {
    log.Fatal(err)
}
```

### Patrón: Funciones que Retornan Punteros

```go
// Constructor: Crea y retorna puntero
func nuevoUsuario(nombre string) *Usuario {
    return &Usuario{
        nombre: nombre,
        creado: time.Now(),
    }
}

// Uso
u := nuevoUsuario("Alice")
fmt.Println(u.nombre)  // Alice

// Multiretorno: Puntero + error (muy común)
func obtenerUsuario(id int) (*Usuario, error) {
    if id <= 0 {
        return nil, errors.New("ID inválido")
    }
    // Buscar usuario...
    return &Usuario{id: id, nombre: "Bob"}, nil
}

// Uso seguro
u, err := obtenerUsuario(1)
if err != nil {
    log.Fatal(err)
}
fmt.Println(u.nombre)
```

### Eficiencia: Evitar Copias Innecesarias

```go
type DatosGrandes struct {
    datos [10000]int
    // Otros campos...
}

// ANTI-PATRÓN: Copia la estructura completa
func procesarSinPuntero(d DatosGrandes) int {
    total := 0
    for _, v := range d.datos {
        total += v
    }
    return total
}

// PATRÓN CORRECTO: Puntero, sin copia
func procesarConPuntero(p *DatosGrandes) int {
    total := 0
    for _, v := range p.datos {
        total += v
    }
    return total
}

func main() {
    d := DatosGrandes{}

    // Tiempo 1: 10ms (copia lenta)
    start := time.Now()
    for i := 0; i < 1000; i++ {
        procesarSinPuntero(d)
    }
    fmt.Println("Sin puntero:", time.Since(start))

    // Tiempo 2: 1ms (sin copias)
    start = time.Now()
    for i := 0; i < 1000; i++ {
        procesarConPuntero(&d)
    }
    fmt.Println("Con puntero:", time.Since(start))
}
```

### Validación de Argumentos Puntero

```go
func saludar(p *Persona) error {
    if p == nil {
        return errors.New("persona es nil")
    }

    if p.nombre == "" {
        return errors.New("nombre vacío")
    }

    fmt.Printf("Hola, %s!\n", p.nombre)
    return nil
}

func main() {
    err := saludar(nil)
    if err != nil {
        fmt.Println("Error:", err)  // Error: persona es nil
    }
}
```

---

## 16.7 Punteros en Métodos

### Receiver por Valor vs Pointer

Un aspecto crítico de Go es elegir entre receiver por valor o puntero:

```
┌──────────────────────────────────────┐
│  Receiver por Valor                  │
├──────────────────────────────────────┤
│  func (p Persona) String() string    │
│                                       │
│  ✓ No modifica el original           │
│  ✓ Funciona con valores y punteros   │
│  ✗ Copia la estructura               │
│                                       │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  Receiver por Puntero                │
├──────────────────────────────────────┤
│  func (p *Persona) Modificar()       │
│                                       │
│  ✓ Modifica el original              │
│  ✓ Sin copias                        │
│  ✗ Solo funciona con punteros        │
│                                       │
└──────────────────────────────────────┘
```

### Métodos Inmutables (Receiver por Valor)

```go
type Punto struct {
    x, y float64
}

// Receiver por valor: Cálcula distancia sin modificar
func (p Punto) Distancia() float64 {
    return math.Sqrt(p.x*p.x + p.y*p.y)
}

// Receiver por valor: Devuelve nuevo Punto
func (p Punto) Trasladar(dx, dy float64) Punto {
    return Punto{x: p.x + dx, y: p.y + dy}
}

func main() {
    pt := Punto{3, 4}

    fmt.Println(pt.Distancia())        // 5
    pt2 := pt.Trasladar(1, 1)
    fmt.Println(pt)                    // {3 4} (sin cambios)
    fmt.Println(pt2)                   // {4 5}
}
```

### Métodos Mutadores (Receiver por Puntero)

```go
type Saldo struct {
    cantidad float64
}

// Receiver por puntero: Modifica el original
func (s *Saldo) Depositar(monto float64) error {
    if monto < 0 {
        return errors.New("monto negativo")
    }
    s.cantidad += monto
    return nil
}

func (s *Saldo) Retirar(monto float64) error {
    if monto > s.cantidad {
        return errors.New("saldo insuficiente")
    }
    s.cantidad -= monto
    return nil
}

func main() {
    saldo := &Saldo{cantidad: 1000}

    saldo.Depositar(500)       // Modifica el original
    fmt.Println(saldo.cantidad) // 1500

    saldo.Retirar(200)         // Modifica el original
    fmt.Println(saldo.cantidad) // 1300
}
```

### Implicit Dereference: Azúcar Sintáctico

Go automáticamente dereferencia punteros para llamadas a métodos:

```go
type Persona struct {
    nombre string
}

func (p *Persona) SaludarPuntero() {
    fmt.Printf("Hola desde %s (puntero)\n", p.nombre)
}

func (p Persona) SaludarValor() {
    fmt.Printf("Hola desde %s (valor)\n", p.nombre)
}

func main() {
    // Con puntero
    p := &Persona{nombre: "Alice"}
    p.SaludarPuntero()      // ✓ Funciona
    p.SaludarValor()        // ✓ Go dereferencia automáticamente

    // Con valor
    v := Persona{nombre: "Bob"}
    v.SaludarPuntero()      // ✓ Go toma dirección automáticamente
    v.SaludarValor()        // ✓ Funciona

    // Equivalencias:
    // p.SaludarPuntero() ≡ (*p).SaludarPuntero()
    // p.SaludarValor() ≡ (*p).SaludarValor()
    // v.SaludarPuntero() ≡ (&v).SaludarPuntero()
}
```

### Regla: Elige Receiver por Puntero Cuando

**1. El método modifica el receptor:**

```go
func (p *Persona) CambiarNombre(nombre string) {
    p.nombre = nombre  // Modifica el original
}
```

**2. La estructura es grande (evitar copias):**

```go
type DatosGrandes struct {
    array [10000]int
}

// ✓ Correcto: Sin copias
func (d *DatosGrandes) Procesar() {}

// ✗ Ineficiente: Copia todo
func (d DatosGrandes) Procesar() {}
```

**3. Consistencia en el tipo:**

```go
type Usuario struct {
    nombre string
}

// Si algunos métodos usan puntero, todos deberían
func (u *Usuario) Saludar() { ... }
func (u *Usuario) Cambiar() { ... }

// Evitar mezclar para mantener consistencia
```

### Ejemplo Completo: Gestor de Tareas

```go
type Tarea struct {
    id       int
    titulo   string
    completa bool
}

type ListaTareas struct {
    tareas []*Tarea
}

// Constructor
func nuevaListaTareas() *ListaTareas {
    return &ListaTareas{
        tareas: make([]*Tarea, 0),
    }
}

// Método mutador: Añadir tarea
func (lt *ListaTareas) Agregar(titulo string) {
    t := &Tarea{
        id:     len(lt.tareas) + 1,
        titulo: titulo,
    }
    lt.tareas = append(lt.tareas, t)
}

// Método mutador: Marcar completa
func (lt *ListaTareas) Completar(id int) error {
    for _, t := range lt.tareas {
        if t.id == id {
            t.completa = true
            return nil
        }
    }
    return errors.New("tarea no encontrada")
}

// Método immutable: Contar
func (lt *ListaTareas) Contar() int {
    return len(lt.tareas)
}

// Método immutable: Obtener tarea
func (lt *ListaTareas) Obtener(id int) *Tarea {
    for _, t := range lt.tareas {
        if t.id == id {
            return t
        }
    }
    return nil
}

func main() {
    lista := nuevaListaTareas()
    lista.Agregar("Estudiar Go")
    lista.Agregar("Hacer ejercicio")

    fmt.Printf("Tareas: %d\n", lista.Contar())  // 2

    lista.Completar(1)

    for _, t := range lista.tareas {
        estado := "Pendiente"
        if t.completa {
            estado = "Completa"
        }
        fmt.Printf("[%d] %s - %s\n", t.id, t.titulo, estado)
    }
}
```

---

## 16.8 Punteros e Interfaces

### El Tipo interface{} (Interfaz Vacía)

La interfaz vacía `interface{}` puede contener cualquier valor:

```go
var x interface{} = 5
var y interface{} = "hola"
var z interface{} = []int{1, 2, 3}
var w interface{} = nil

fmt.Println(x)  // 5
fmt.Println(y)  // hola
fmt.Println(z)  // [1 2 3]
fmt.Println(w)  // <nil>
```

### Punteros en interface{}

```go
// interface{} puede contener un puntero
type Usuario struct {
    nombre string
}

var datos interface{} = &Usuario{nombre: "Alice"}

// Para acceder a los datos, necesitas type assertion
u := datos.(*Usuario)
fmt.Println(u.nombre)  // Alice
```

### Type Assertion con Punteros

```go
var datos interface{} = &Usuario{nombre: "Bob"}

// Forma 1: Sin verificación (pánico si falla)
u := datos.(*Usuario)
fmt.Println(u.nombre)  // Bob

// Forma 2: Con verificación (seguro)
u, ok := datos.(*Usuario)
if ok {
    fmt.Println(u.nombre)  // Bob
} else {
    fmt.Println("No es un *Usuario")
}

// Forma 3: Intentar diferentes tipos
switch v := datos.(type) {
case *Usuario:
    fmt.Println("Usuario:", v.nombre)
case *Persona:
    fmt.Println("Persona:", v.nombre)
case string:
    fmt.Println("String:", v)
default:
    fmt.Println("Tipo desconocido:", reflect.TypeOf(datos))
}
```

### Type Assertion: Valor vs Puntero

Diferencia importante:

```go
// Estructura
type Usuario struct {
    nombre string
}

// Caso 1: interface{} contiene valor
var datos interface{} = Usuario{nombre: "Alice"}

// Esto falla:
u1 := datos.(*Usuario)  // panic: interface conversion error

// Esto funciona:
u2 := datos.(Usuario)   // ✓

// Caso 2: interface{} contiene puntero
var datos interface{} = &Usuario{nombre: "Bob"}

// Esto funciona:
u3 := datos.(*Usuario)  // ✓

// Esto falla:
u4 := datos.(Usuario)   // panic: interface conversion error
```

### Interfaces Comunes que Requieren Punteros

```go
// io.Reader: Lee datos
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Método debe usar puntero para modificar buffer
func (f *File) Read(p []byte) (n int, err error) {
    // Modifica p, por eso necesita puntero
}

// json.Unmarshaler: Deserializa JSON
type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}

// Método debe usar puntero para modificar struct
func (u *Usuario) UnmarshalJSON(data []byte) error {
    // Modifica u, por eso necesita puntero
}
```

### Ejemplo Práctico: Procesador Genérico

```go
type Procesable interface {
    Procesar() error
}

type PDF struct {
    nombre string
}

func (p *PDF) Procesar() error {
    fmt.Printf("Procesando PDF: %s\n", p.nombre)
    return nil
}

type Imagen struct {
    nombre string
}

func (i *Imagen) Procesar() error {
    fmt.Printf("Procesando Imagen: %s\n", i.nombre)
    return nil
}

func procesarArchivos(archivos []Procesable) error {
    for _, archivo := range archivos {
        if err := archivo.Procesar(); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    archivos := []Procesable{
        &PDF{nombre: "documento.pdf"},
        &Imagen{nombre: "foto.jpg"},
        &PDF{nombre: "reporte.pdf"},
    }

    if err := procesarArchivos(archivos); err != nil {
        log.Fatal(err)
    }
}
```

---

## 16.9 Unsafe: Aritmética de Punteros

### ¿Qué es Unsafe?

El paquete `unsafe` permite operaciones peligrosas que Go normalmente protege:

```go
import "unsafe"

// Funciones principales:
// unsafe.Pointer     - Puntero genérico (como void* en C)
// unsafe.Sizeof()    - Tamaño de un tipo
// unsafe.Offsetof()  - Offset de un campo
// unsafe.Alignof()   - Alineación de un tipo
// uintptr            - Entero que representa dirección
```

### ADVERTENCIA: Usar Unsafe es Peligroso

```
⚠️  ADVERTENCIA ⚠️

unsafe se llama así POR UNA RAZÓN.

Usando unsafe puede:
✗ Corromper datos
✗ Causar pánicas impredecibles
✗ Violar la seguridad del programa
✗ Hacer código no portable (dependiente de arquitectura)

SOLO usa unsafe cuando:
1. Necesitas interoperar con C (cgo)
2. Necesitas performance extrema
3. No hay alternativa segura
4. REALMENTE sabes qué haces
```

### unsafe.Sizeof: Tamaño de Tipos

```go
import (
    "fmt"
    "unsafe"
)

func main() {
    // Tipos básicos
    fmt.Println(unsafe.Sizeof(int(0)))        // 8 (en 64-bit)
    fmt.Println(unsafe.Sizeof(float64(0)))    // 8
    fmt.Println(unsafe.Sizeof(bool(false)))   // 1

    // Structs
    type Usuario struct {
        id    int      // 8 bytes
        edad  int      // 8 bytes
        ativo bool     // 1 byte
    }

    fmt.Println(unsafe.Sizeof(Usuario{}))     // 24 (con padding)

    // Nota: El tamaño incluye padding para alineación
}
```

### unsafe.Offsetof: Posición de Campos

```go
type Persona struct {
    nombre string  // 16 bytes (string = puntero + len)
    edad   int     // 8 bytes
    activo bool    // 1 byte
}

func main() {
    p := Persona{}

    // Offset desde el inicio
    fmt.Println(unsafe.Offsetof(p.nombre))  // 0
    fmt.Println(unsafe.Offsetof(p.edad))    // 16
    fmt.Println(unsafe.Offsetof(p.activo))  // 24
}
```

### Aritmética de Punteros Básica

```go
import (
    "fmt"
    "unsafe"
)

func main() {
    // Array
    arr := [5]int{10, 20, 30, 40, 50}

    // Puntero al primer elemento
    p := &arr[0]
    fmt.Println(*p)  // 10

    // Aritmética con unsafe
    // Obtener el siguiente elemento
    p2 := (*int)(unsafe.Pointer(
        uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(int(0)),
    ))
    fmt.Println(*p2)  // 20

    // El tercer elemento
    p3 := (*int)(unsafe.Pointer(
        uintptr(unsafe.Pointer(p)) + 2*unsafe.Sizeof(int(0)),
    ))
    fmt.Println(*p3)  // 30
}
```

### Caso de Uso: Interoperación con C

```go
import (
    "fmt"
    "unsafe"
)

// Supongamos que C tiene esta estructura:
// typedef struct {
//     int id;
//     char name[256];
// } CUser;

type CUser struct {
    id   int32    // 4 bytes
    name [256]byte
}

func main() {
    // Crear instancia
    cuser := CUser{id: 42}
    copy(cuser.name[:], "Alice")

    // Obtener puntero para pasar a C
    ptr := unsafe.Pointer(&cuser)

    // En Go no podemos llamar C directamente aquí
    // Pero así se haría en cgo:
    // C.procesarUsuario((*C.CUser)(ptr))

    fmt.Println("Puntero para C:", ptr)
}
```

### Lectura de Bytes: Ejemplo Práctico

```go
import (
    "fmt"
    "unsafe"
)

func main() {
    x := 0x12345678

    // Obtener los bytes que forman el entero
    p := &x

    // Convertir a puntero de bytes
    bp := (*[4]byte)(unsafe.Pointer(p))

    fmt.Printf("Bytes de %x: ", x)
    for _, b := range bp {
        fmt.Printf("%02x ", b)
    }
    fmt.Println()

    // Salida en máquina little-endian:
    // Bytes de 12345678: 78 56 34 12
}
```

### Conversión entre Tipos (Muy Peligroso)

```go
func main() {
    // Convertir float64 a int64 de manera fea
    f := 3.14

    // Forma segura:
    i := int64(f)  // 3

    // Forma unsafe (cambiar representación binaria):
    fp := unsafe.Pointer(&f)
    ip := (*int64)(fp)

    fmt.Println("Seguro:", i)        // 3
    fmt.Println("Unsafe:", *ip)      // 4614256656552045286 (basura)

    // ✗ Los datos se interpretan de manera diferente
}
```

### Pool de Objetos: Caso Real

```go
import (
    "fmt"
    "sync"
    "unsafe"
)

// Pool preallocado de estructuras
type Objeto struct {
    id   int
    dato [1000]byte
}

type Pool struct {
    objects  []*Objeto
    current  int
    mu       sync.Mutex
}

func nuevoPool(tamaño int) *Pool {
    p := &Pool{
        objects: make([]*Objeto, tamaño),
    }

    // Preallocar objetos
    for i := 0; i < tamaño; i++ {
        p.objects[i] = &Objeto{id: i}
    }

    fmt.Printf("Pool: %d objetos prealloc, ~%d KB total\n",
        tamaño,
        (tamaño*int(unsafe.Sizeof(Objeto{})))/1024,
    )

    return p
}

func (p *Pool) Obtener() *Objeto {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.current >= len(p.objects) {
        p.current = 0  // Circular
    }

    obj := p.objects[p.current]
    p.current++

    return obj
}

func main() {
    pool := nuevoPool(100)

    for i := 0; i < 5; i++ {
        obj := pool.Obtener()
        fmt.Printf("Objeto %d\n", obj.id)
    }
}
```

---

## 16.10 Comparación de Punteros

### Comparación de Igualdad

```go
type Usuario struct {
    nombre string
}

func main() {
    u1 := Usuario{nombre: "Alice"}
    u2 := Usuario{nombre: "Alice"}

    p1 := &u1
    p2 := &u2
    p3 := &u1

    fmt.Println(p1 == p2)  // false (direcciones diferentes)
    fmt.Println(p1 == p3)  // true (misma dirección)
    fmt.Println(p1 != p2)  // true

    // Comparar con nil
    var p *Usuario
    fmt.Println(p == nil)      // true
    fmt.Println(p1 == nil)     // false
}
```

### Comparación de Contenido

```go
// Para comparar el CONTENIDO (no la dirección):
func main() {
    u1 := Usuario{nombre: "Alice"}
    u2 := Usuario{nombre: "Alice"}

    p1 := &u1
    p2 := &u2

    // Comparar punteros (dirección)
    fmt.Println(p1 == p2)      // false

    // Comparar contenido
    fmt.Println(*p1 == *p2)    // true

    // Para structs grandes, mejor usar reflect o métodos custom
    fmt.Println(reflect.DeepEqual(*p1, *p2))  // true
}
```

### Comparación antes de Usar

```go
func procesarUsuario(u *Usuario) error {
    // Verificación esencial
    if u == nil {
        return errors.New("usuario es nil")
    }

    // Ahora es seguro desreferenciar
    if u.nombre == "" {
        return errors.New("nombre vacío")
    }

    fmt.Println(u.nombre)
    return nil
}

func main() {
    var u *Usuario

    if err := procesarUsuario(u); err != nil {
        fmt.Println("Error:", err)  // Error: usuario es nil
    }
}
```

### Búsqueda de Puntero en Slice

```go
func encontrarUsuario(usuarios []*Usuario, id int) *Usuario {
    for _, u := range usuarios {
        if u.id == id {
            return u  // Retorna el puntero
        }
    }
    return nil  // No encontrado
}

func main() {
    usuarios := []*Usuario{
        &Usuario{id: 1, nombre: "Alice"},
        &Usuario{id: 2, nombre: "Bob"},
    }

    u := encontrarUsuario(usuarios, 1)
    if u != nil {
        fmt.Println(u.nombre)
    }
}
```

### Tipos Comparables

No todos los tipos son comparables:

```go
// Comparables
var p1 *int
var p2 *int
fmt.Println(p1 == p2)  // ✓ Válido

// No comparables
var s1 []int = []int{1, 2, 3}
var s2 []int = []int{1, 2, 3}
// fmt.Println(s1 == s2)  // ✗ Error: slice can only be compared to nil

// Workaround: Comparar punteros a arrays
arr1 := [3]int{1, 2, 3}
arr2 := [3]int{1, 2, 3}
p1 := &arr1
p2 := &arr2
fmt.Println(p1 == p2)  // false (direcciones diferentes)

// O usar reflect
fmt.Println(reflect.DeepEqual(s1, s2))  // true
```

---

## 16.11 Buenas Prácticas y Antipatrones

### Buena Práctica 1: Documentar Mutación

```go
// Documenta si la función modifica los argumentos
type Conexión struct {
    buffer []byte
}

// MALA: No está claro que p se modifica
func (p *Conexión) Leer() error {
    p.buffer = p.buffer[:100]
    // ...
    return nil
}

// BUENA: Está claro desde el nombre y doc
// // Leer lee datos y actualiza el buffer interno.
// // Nota: Modifica c en su lugar.
func (c *Conexión) Leer() error {
    c.buffer = c.buffer[:100]
    // ...
    return nil
}

// O retorna el nuevo valor
func (c *Conexión) Leer() ([]byte, error) {
    buffer := make([]byte, 100)
    // ...
    return buffer, nil
}
```

### Buena Práctica 2: Prefiere No Punteros por Defecto

```go
// ANTI-PATRÓN: Puntero innecesario
type Punto struct {
    x, y float64
}

// ✗ Evita punteros innecesarios
func (p *Punto) Distancia() float64 {
    return math.Sqrt(p.x*p.x + p.y*p.y)
}

// PATRÓN CORRECTO: Valor por defecto
func (p Punto) Distancia() float64 {
    return math.Sqrt(p.x*p.x + p.y*p.y)
}

// Solo usa puntero si NECESITAS:
// 1. Modificar el receptor
// 2. Evitar copia de tipo grande
func (p *Punto) Trasladar(dx, dy float64) {
    p.x += dx
    p.y += dy
}
```

### Buena Práctica 3: Nil Checks Explícitos

```go
func guardarEnArchivo(ruta string, datos *[]byte) error {
    // ✓ Verificar nil explícitamente
    if datos == nil {
        return errors.New("datos nil")
    }

    if len(*datos) == 0 {
        return errors.New("datos vacíos")
    }

    return os.WriteFile(ruta, *datos, 0644)
}

// Más limpio: Convertir a interfaz
func guardarEnArchivo(ruta string, datos interface{}) error {
    if datos == nil {
        return errors.New("datos nil")
    }

    // ...
    return nil
}
```

### Buena Práctica 4: Usar Constructores

```go
// ANTI-PATRÓN: Sin inicialización clara
type Config struct {
    puerto int
    host   string
    ssl    bool
}

// PATRÓN: Constructor explícito
func NewConfig() *Config {
    return &Config{
        puerto: 8080,
        host:   "localhost",
        ssl:    false,
    }
}

// Con opciones
func NewConfigWithPort(puerto int) *Config {
    c := NewConfig()
    c.puerto = puerto
    return c
}

func main() {
    cfg := NewConfig()
    fmt.Printf("Conectando a %s:%d\n", cfg.host, cfg.puerto)
}
```

### Buena Práctica 5: Evitar Copias Innecesarias

```go
// ANTI-PATRÓN: Copiar en cada llamada
type DatosGrandes struct {
    array [10000]int
}

func procesarDatos(datos DatosGrandes) {  // ✗ Copia
    // ...
}

// PATRÓN: Usar puntero o slice
func procesarDatos(datos *DatosGrandes) {  // ✓ Puntero
    // ...
}

// O si es un slice:
func procesarDatos(datos []int) {  // ✓ Slice (es puntero + meta)
    // ...
}
```

### Antipatrón 1: Pointer Receivers Innecesarios

```go
// ✗ MAL: Receiver es puntero aunque no se modifica
type Rectángulo struct {
    ancho, alto int
}

func (r *Rectángulo) Área() int {
    return r.ancho * r.alto
}

// ✓ CORRECTO: Receiver es valor
func (r Rectángulo) Área() int {
    return r.ancho * r.alto
}

// Ventajas:
// - Más claro (sin mutación)
// - Funciona con valores y punteros
// - Sem copias de tipos pequeños
```

### Antipatrón 2: Devolver Punteros a Variables Locales

```go
// En C, esto es PELIGROSO (dangling pointer)
// En Go, es SEGURO gracias a GC

func obtenerPersona() *Persona {
    p := Persona{nombre: "Alice"}  // Variable local
    return &p                       // ✓ Seguro en Go (GC lo moverá al heap)
}

// Go automáticamente promueve 'p' al heap
func main() {
    p := obtenerPersona()
    fmt.Println(p.nombre)  // Alice (funciona perfectamente)
}
```

### Antipatrón 3: Unsafe Excesivo

```go
// ✗ ANTI-PATRÓN: Usar unsafe cuando hay alternativa segura
func compararPunteros(p1, p2 *int) bool {
    return uintptr(unsafe.Pointer(p1)) == uintptr(unsafe.Pointer(p2))
}

// ✓ CORRECTO: Usar comparación de punteros
func compararPunteros(p1, p2 *int) bool {
    return p1 == p2
}

// ✗ ANTI-PATRÓN: Aritmética de punteros innecesaria
func siguienteElemento(arr *[10]int, idx int) *int {
    p := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&arr[0])) +
        uintptr(idx)*unsafe.Sizeof(int(0))))
    return p
}

// ✓ CORRECTO: Usar índices normales
func siguienteElemento(arr *[10]int, idx int) *int {
    return &arr[idx]
}
```

### Antipatrón 4: Confundir Punteros a Arrays con Slices

```go
// Datos
datos := []int{1, 2, 3, 4, 5}

// ✗ ANTI-PATRÓN: Puntero a slice innecesario
func procesarDatos(ps *[]int) {
    (*ps)[0] = 100  // Tedioso
}

// ✓ CORRECTO: El slice ya es eficiente
func procesarDatos(s []int) {
    s[0] = 100      // Limpio
}

// ✓ TAMBIÉN CORRECTO: Puntero a array cuando necesitas cambiar tamaño
arr := [10]int{1, 2, 3}
func redimensionar(pa *[10]int) {
    // Puedes usar pa
}
```

### Checklist de Buenas Prácticas

```
ANTES de usar puntero, pregúntate:

☐ ¿Necesito MODIFICAR el valor original?
   → SÍ: Usa puntero
   → NO: Continúa

☐ ¿Es la estructura GRANDE (>128 bytes)?
   → SÍ: Considera puntero (eficiencia)
   → NO: Continúa

☐ ¿Será NIL frecuentemente?
   → SÍ: Considera Optional pattern
   → NO: Continúa

☐ ¿Necesito ciclos de referencia (lista, árbol)?
   → SÍ: Necesitas punteros
   → NO: Probablemente no necesitas

DEFAULT: Usa valor, añade puntero solo si lo necesitas
```

---

## Ejercicios Progresivos

### Ejercicio 1: Swap de Variables

**Objetivo:** Implementar una función que intercambie dos variables usando punteros.

**Requisitos:**

- Función `Swap(a, b *int)` que intercambia valores
- Función `SwapStrings(a, b *string)` que intercambia strings
- Función genérica usando `interface{}`
- Pruebas con valores iniciales diferentes

**Código Base:**

```go
package main

import "fmt"

func Swap(a, b *int) {
    // TODO: Implementar
}

func SwapStrings(a, b *string) {
    // TODO: Implementar
}

// Genérica (más desafiante)
func SwapGeneric(a, b interface{}) {
    // TODO: Implementar (pista: necesitas reflection)
}

func main() {
    x, y := 10, 20
    fmt.Printf("Antes: x=%d, y=%d\n", x, y)
    Swap(&x, &y)
    fmt.Printf("Después: x=%d, y=%d\n", x, y)

    s1, s2 := "hola", "mundo"
    fmt.Printf("Antes: s1=%s, s2=%s\n", s1, s2)
    SwapStrings(&s1, &s2)
    fmt.Printf("Después: s1=%s, s2=%s\n", s1, s2)
}
```

**Solución esperada:** ~30 líneas

---

### Ejercicio 2: Gestor de Nodos de Lista Enlazada

**Objetivo:** Implementar una lista enlazada con inserción y eliminación.

**Requisitos:**

- Estructura `Nodo` con valor y puntero al siguiente
- Estructura `Lista` que gestiona el inicio
- Método `Agregar(valor int)`
- Método `Eliminar(valor int) bool`
- Método `Imprimir()` que muestra todos los valores
- Método `Buscar(valor int) *Nodo`

**Código Base:**

```go
package main

import "fmt"

type Nodo struct {
    valor     int
    siguiente *Nodo
}

type Lista struct {
    inicio *Nodo
}

func (l *Lista) Agregar(valor int) {
    // TODO: Implementar
}

func (l *Lista) Eliminar(valor int) bool {
    // TODO: Implementar
}

func (l *Lista) Imprimir() {
    // TODO: Implementar
}

func (l *Lista) Buscar(valor int) *Nodo {
    // TODO: Implementar
}

func main() {
    lista := &Lista{}

    lista.Agregar(10)
    lista.Agregar(20)
    lista.Agregar(30)

    lista.Imprimir()  // 10 -> 20 -> 30 -> nil

    lista.Eliminar(20)
    lista.Imprimir()  // 10 -> 30 -> nil

    nodo := lista.Buscar(30)
    if nodo != nil {
        fmt.Printf("Encontrado: %d\n", nodo.valor)
    }
}
```

**Solución esperada:** ~60 líneas

---

### Ejercicio 3: Modificador de Structs Complejos

**Objetivo:** Funciones que modifican structs complejos usando punteros.

**Requisitos:**

- Estructura `Empleado` con campos anidados
- Función `ActualizarSalario(e *Empleado, nuevo float64) error`
- Función `CambiarDepartamento(e *Empleado, dept string) error`
- Método `IncrementerSalario(porcentaje float64)`
- Validación de datos (salario > 0, departamento válido)

**Código Base:**

```go
package main

import (
    "errors"
    "fmt"
)

type Empleado struct {
    id          int
    nombre      string
    salario     float64
    departamento string
}

func ActualizarSalario(e *Empleado, nuevo float64) error {
    // TODO: Implementar con validación
}

func CambiarDepartamento(e *Empleado, dept string) error {
    // TODO: Implementar con validación
}

func (e *Empleado) Incrementar(porcentaje float64) error {
    // TODO: Implementar
}

func (e *Empleado) Mostrar() {
    fmt.Printf("[%d] %s - %s - $%.2f\n", e.id, e.nombre, e.departamento, e.salario)
}

func main() {
    emp := &Empleado{id: 1, nombre: "Alice", salario: 50000, departamento: "IT"}
    emp.Mostrar()

    if err := ActualizarSalario(emp, 55000); err != nil {
        fmt.Println("Error:", err)
    }
    emp.Mostrar()

    if err := emp.Incrementar(10); err != nil {
        fmt.Println("Error:", err)
    }
    emp.Mostrar()
}
```

**Solución esperada:** ~50 líneas

---

### Ejercicio 4: Unsafe - Lectura de Estructura en Bytes

**Objetivo:** Usar `unsafe` para leer la representación binaria de una estructura.

**Requisitos:**

- Leer tamaño de estructura con `unsafe.Sizeof`
- Leer offset de campos con `unsafe.Offsetof`
- Convertir struct a bytes usando `unsafe.Pointer`
- Mostrar representación hexadecimal
- IMPORTANTE: Documenta el peligro de unsafe

**Código Base:**

```go
package main

import (
    "fmt"
    "unsafe"
)

type Registro struct {
    id    int32    // 4 bytes
    precio float32  // 4 bytes
    activo bool     // 1 byte
}

func inspeccionarEstructura() {
    // TODO: Usar unsafe.Sizeof
    // TODO: Usar unsafe.Offsetof para cada campo
    // TODO: Mostrar tamaños y offsets
}

func estructuraABytes(r *Registro) []byte {
    // TODO: Convertir struct a bytes usando unsafe
    size := unsafe.Sizeof(*r)
    b := make([]byte, size)

    // PELIGRO: Esta es una operación unsafe
    // No uses esto en código de producción

    return b
}

func main() {
    inspeccionarEstructura()

    reg := &Registro{id: 42, precio: 19.99, activo: true}
    bytes := estructuraABytes(reg)

    fmt.Println("Bytes:", bytes)
    fmt.Printf("Tamaño: %d bytes\n", len(bytes))
}
```

**Solución esperada:** ~40 líneas

**Advertencia:** Este ejercicio enseña `unsafe`, pero incluir comentarios sobre los peligros es obligatorio.

---

### Ejercicio 5: Pool de Objetos Prealloc

**Objetivo:** Crear un gestor de pool de objetos para evitar allocaciones repetidas.

**Requisitos:**

- Estructura `Pool` que prealoca N objetos
- Estructura `Item` (simple, con datos simulados)
- Método `Obtener()` que devuelve un Item del pool
- Método `Devolver(item *Item)` que devuelve al pool
- Método `Estadísticas()` que muestra uso
- Usar goroutines y `sync.Mutex` para thread-safety
- Mostrar comparación de performance: con pool vs sin pool

**Código Base:**

```go
package main

import (
    "fmt"
    "sync"
    "time"
    "unsafe"
)

type Item struct {
    id   int
    dato [256]byte  // Simular datos grandes
}

type Pool struct {
    disponibles chan *Item
    total       int
    usado       int
    mu          sync.Mutex
}

func NewPool(tamaño int) *Pool {
    p := &Pool{
        disponibles: make(chan *Item, tamaño),
        total:       tamaño,
    }

    // TODO: Preallocar items

    return p
}

func (p *Pool) Obtener() *Item {
    // TODO: Obtener del pool (con timeout si es necesario)
    return nil
}

func (p *Pool) Devolver(item *Item) {
    // TODO: Devolver al pool
}

func (p *Pool) Estadísticas() {
    p.mu.Lock()
    defer p.mu.Unlock()

    fmt.Printf("Pool: Total=%d, Usado=%d, Disponibles=%d, "+
        "Tamaño=~%dKB\n",
        p.total, p.usado, len(p.disponibles),
        (p.total*int(unsafe.Sizeof(Item{})))/1024,
    )
}

func main() {
    pool := NewPool(100)
    pool.Estadísticas()

    // Simular múltiples goroutines
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            item := pool.Obtener()
            if item != nil {
                item.id = id
                time.Sleep(10 * time.Millisecond)
                pool.Devolver(item)
            }
        }(i)
    }

    wg.Wait()
    pool.Estadísticas()
}
```

**Solución esperada:** ~80 líneas

**Desafío adicional:** Implementar reordenamiento circular en el pool (reusar slots)

---

## Resumen de Conceptos Clave

| Concepto | Descripción | Cuándo Usar |
|----------|-------------|-----------|
| **Puntero** | Variable que almacena dirección de memoria | Mutación, eficiencia, estructuras enlazadas |
| **&** (Address-of) | Obtiene dirección de variable | Pasar por referencia |
| **\*** (Dereference) | Accede a valor de puntero | Leer/modificar a través de puntero |
| **nil** | Puntero sin valor | Indicar "sin objeto", requiere checks |
| **Pointer Receiver** | Método mutador con `*T` | Métodos que modifican struct |
| **Value Receiver** | Método inmutable con `T` | Métodos que solo leen |
| **Implicit Deref** | Go dereferencia automáticamente | Acceso a campos: `p.campo` vs `(*p).campo` |
| **Slice vs *Array** | Slice ya contiene puntero | Prefer `[]T` over `*[n]T` |
| **Nil Checks** | Verificar antes de usar | `if p != nil { ... }` |
| **unsafe** | Operaciones peligrosas | Solo cuando REALMENTE sea necesario |

---

## Conclusión

Los punteros son herramientas poderosas en Go que permiten:

1. **Mutación controlada** de datos
2. **Eficiencia** evitando copias innecesarias
3. **Estructuras de datos complejas** como listas y árboles
4. **Interoperación con C** a través de cgo

**Principios clave:**

- Usa punteros solo cuando sea necesario
- Verifica `nil` antes de desreferenciar
- Prefiere value receivers por defecto
- Evita `unsafe` a menos que sea absolutamente necesario
- Documenta la mutación claramente

Go balances seguridad con flexibilidad mejor que C, pero requiere disciplina. Los punteros "seguros" de Go (con garbage collection) hacen que sea mucho más fácil escribir código correcto.

---

**Última actualización:** 2024
**Versión:** 1.0
**Autor:** Guía Exhaustiva de Go en Español

---

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/16-punteros/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/16-punteros):

```bash
cd examples/16-punteros
go run .
```
