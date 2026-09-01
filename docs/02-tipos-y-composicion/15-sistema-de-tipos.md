# Capítulo 15: Sistema de tipos - Definiciones, aliases y assertions

## Índice del Capítulo 15

1. [15.1 El Sistema de Tipos en Go](#151-el-sistema-de-tipos-en-go)
2. [15.2 Type Definitions vs Type Aliases](#152-type-definitions-vs-type-aliases)
3. [15.3 Definición de Tipos Personalizados](#153-definición-de-tipos-personalizados)
4. [15.4 Type Aliases - Go 1.9+](#154-type-aliases---go-19)
5. [15.5 Conversión de Tipos](#155-conversión-de-tipos)
6. [15.6 Compatibilidad de Tipos](#156-compatibilidad-de-tipos)
7. [15.7 Type Assertions](#157-type-assertions)
8. [15.8 Type Switches](#158-type-switches)
9. [15.9 Type Embedding](#159-type-embedding)
10. [15.10 Métodos con Tipos Diferentes](#1510-métodos-con-tipos-diferentes)
11. [15.11 Buenas Prácticas y Diseño de Tipos](#1511-buenas-prácticas-y-diseño-de-tipos)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 15.1 El Sistema de Tipos en Go

### La Filosofía de Tipificación en Go

Go es un lenguaje con **tipado estático fuerte e inferencia de tipos**. Esto significa:

1. **Tipado Estático Fuerte**: Cada variable tiene un tipo específico determinado en tiempo de compilación
2. **Inferencia Automática**: No necesitas declarar el tipo explícitamente; Go lo deduce
3. **Seguridad de Tipos**: Muchos errores se detectan en compilación, no en runtime
4. **Simplicidad**: El sistema de tipos es deliberadamente simple comparado con TypeScript o Rust

```go
// Tipado explícito
var edad int = 25

// Con inferencia
edad := 25  // Go deduce que es int

// Incorrecto: Go rechaza esto en compilación
edad = "veinticinco"  // Error: cannot use string as int
```

### Modelo de Tipos en Go: Identity vs Behavior

Go utiliza un modelo único de tipos llamado **"identity-based typing"** para tipos definidos por el usuario, pero **"structural typing"** para interfaces:

```
┌─────────────────────────────────────────────────┐
│       Sistema de Tipos en Go                    │
├─────────────────────────────────────────────────┤
│                                                  │
│  ┌──────────────────┐    ┌──────────────────┐  │
│  │  Tipos Básicos   │    │  Tipos de Usuario│  │
│  │  (int, string)   │    │  (Definidos)     │  │
│  │  Identity-based  │    │  Identity-based  │  │
│  └──────────────────┘    └──────────────────┘  │
│                                                  │
│  ┌──────────────────────────────────────────┐  │
│  │  Interfaces - Structural Typing          │  │
│  │  (Verifican métodos, no tipo específico) │  │
│  └──────────────────────────────────────────┘  │
│                                                  │
└─────────────────────────────────────────────────┘
```

### Ventajas del Sistema de Tipos de Go

**1. Seguridad en Compilación**

```go
func procesarEdad(edad int) string {
    return "Edad: " + edad  // Error en compilación ✓
}
```

**2. Sin Sorpresas en Runtime**

```go
// JavaScript: "5" + 3 = "53" (sorpresa!)
// Python: "5" + 3 = TypeError (error en runtime)
// Go: "5" + 3 = Error en compilación (esperado) ✓
```

**3. Documentación Automática**

```go
func guardarArchivo(contenido []byte, ruta string) error {
    // El tipo te dice exactamente qué espera
}
```

**4. Optimización del Compilador**

```go
// Go puede optimizar agresivamente porque conoce tipos exactos
for i := 0; i < 1_000_000; i++ {
    // El compilador sabe que i es int, puede optimizar loops
}
```

### Comparación con Otros Lenguajes

```
┌────────────┬──────────────┬────────────────┬─────────────────┐
│ Lenguaje   │ Tipificación │ Verificación   │ Flexibilidad    │
├────────────┼──────────────┼────────────────┼─────────────────┤
│ Go         │ Estática     │ Compilación    │ Alta (interfaces)
│ TypeScript │ Estática     │ Compilación    │ Muy Alta        │
│ Python     │ Dinámica     │ Runtime        │ Muy Alta        │
│ Rust       │ Estática     │ Compilación    │ Media           │
│ C++        │ Estática     │ Compilación    │ Muy Alta/Caótica│
└────────────┴──────────────┴────────────────┴─────────────────┘
```

**Go vs TypeScript**:

- Go es más simple, menos features
- TypeScript es más complejo, más seguridad de tipos
- Go confía en interfaces, TypeScript en tipos estructurales explícitos

**Go vs Python**:

- Go: errores en compilación, código más predecible
- Python: errores en runtime, código más flexible

**Go vs Rust**:

- Go: sistema de tipos más simple, más fácil de aprender
- Rust: sistema de tipos más complejo, más seguro (borrow checker)

---

## 15.2 Type Definitions vs Type Aliases

### El Dilema: ¿Cuándo Crear un Tipo Nuevo?

Antes de Go 1.9, solo existían **type definitions**. Go 1.9 introdujo **type aliases**. Son conceptos diferentes:

```
┌─────────────────────────────────────────────────────────┐
│  Type Definition (type X T)                             │
├─────────────────────────────────────────────────────────┤
│  • Crea un NUEVO tipo distinto                          │
│  • Go los considera tipos diferentes                    │
│  • Necesita conversión explícita entre ellos           │
│  • Más seguro, define semántica clara                  │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  Type Alias (type X = T)  [Go 1.9+]                    │
├─────────────────────────────────────────────────────────┤
│  • Crea un ALIAS del tipo existente                     │
│  • Go los considera el mismo tipo                       │
│  • Compatible sin conversión explícita                 │
│  • Útil para migración de código                       │
└─────────────────────────────────────────────────────────┘
```

### Type Definition: Crear Nuevos Tipos

```go
// Definición de tipo: Metros es un nuevo tipo basado en float64
type Metros float64

// Definición de tipo: Velocidad es un nuevo tipo basado en float64
type Velocidad float64

func main() {
    distancia := Metros(100.5)      // Conversión explícita
    velocidad := Velocidad(25.0)

    // ERROR: No puedes sumar directamente tipos diferentes
    // total := distancia + velocidad  // ✗ Compile error

    // Correcto: Conversión explícita
    total := Metros(float64(velocidad) + float64(distancia))
    fmt.Println(total)  // 125.5
}
```

**¿Por qué Go rechaza sumar Metros + Velocidad?**

Porque semánticamente no tiene sentido. Es una protección del lenguaje contra bugs sutiles:

```go
// Escenario sin type definitions (solo float64):
func calcularEnergia(distancia float64, velocidad float64) float64 {
    return distancia * velocidad  // ¡Puede estar invertido! 😱
}

// Con type definitions:
func calcularEnergia(distancia Metros, velocidad Velocidad) {
    // El compilador te lo permite explícitamente si haces conversión
    // Más seguro y documentado
}
```

### Type Alias: Simplemente Otro Nombre

```go
// Go 1.9+: Type Alias (no crea un tipo nuevo)
type MetrosAlias = float64

func main() {
    var valor1 MetrosAlias = 100.5
    var valor2 float64 = 50.0

    // Sí puedes sumar: son exactamente el mismo tipo
    suma := valor1 + valor2
    fmt.Println(suma)  // 150.5
}
```

### Comparación Directa

```go
// DEFINICIÓN (Go < 1.9 y después)
type Kilometros float64      // Nuevo tipo distinto

// ALIAS (Go >= 1.9)
type km = float64           // Otro nombre para float64

func main() {
    var km1 Kilometros = 100
    var km2 km = 50

    // km1 + km2 = ERROR ✗
    // km2 es compatible con float64: NO ERROR ✓

    var f float64 = km2      // Válido
    // var f2 float64 = km1  // ERROR: tipos incompatibles
}
```

### Historia de Go: Por Qué Necesitó Aliases

**Go 1.0 - Go 1.8:**

- Solo existían type definitions
- Problema: Refactorizar código era difícil
- Si renombrabas un tipo, debías convertir explícitamente en todos lados

```go
// Versión antigua
type SistemaViejo struct { }

// Quiero renombrarlo a SistemaModerno...
// Opción 1: Renombrar y romper todo código existente
// Opción 2: Duplicar código (Type Definition)
type SistemaModerno = SistemaViejo  // ✓ Go 1.9 lo permite
```

**Go 1.9:**

- Se añadieron type aliases
- Permite refactorización suave (deprecación gradual)

```go
// Código antiguo
type Usuario struct { ... }

// Nueva estructura reorganizada
type UsuarioV2 struct { ... }

// En la fase de transición:
type Usuario = UsuarioV2  // Alias temporal
// El código antiguo sigue funcionando sin cambios
```

---

## 15.3 Definición de Tipos Personalizados

### Tipos Basados en Tipos Primitivos

```go
// Tiempo en segundos
type Segundos int

// Dinero en centavos
type Centavos int64

// Temperatura en Celsius
type Celsius float32

// Estado de un proceso
type Estado string

func main() {
    tiempo := Segundos(3600)      // 1 hora
    dinero := Centavos(2999)      // 29.99
    temperatura := Celsius(25.5)
    estado := Estado("completado")

    fmt.Println(tiempo, dinero, temperatura, estado)
}
```

### Métodos en Tipos Personalizados

La verdadera potencia de type definitions es que puedes **añadir métodos** al nuevo tipo:

```go
type Segundos int

// Convertir a minutos (método)
func (s Segundos) AMinutos() int {
    return int(s) / 60
}

// Convertir a horas (método)
func (s Segundos) AHoras() int {
    return int(s) / 3600
}

// Validar si el tiempo es positivo (método)
func (s Segundos) EsValido() bool {
    return s > 0
}

func main() {
    tiempo := Segundos(7200)  // 2 horas

    fmt.Println(tiempo.AMinutos())  // 120
    fmt.Println(tiempo.AHoras())    // 2
    fmt.Println(tiempo.EsValido())  // true
}
```

### Tipos Basados en Structs

```go
// Tipo basado en struct
type Usuario struct {
    Nombre   string
    Email    string
    Edad     int
    Premium  bool
}

// Métodos específicos
func (u Usuario) EsAdulto() bool {
    return u.Edad >= 18
}

func (u Usuario) Saludar() string {
    return "Hola " + u.Nombre
}

// Receptor por valor vs referencia
func (u *Usuario) ModificarNombre(nuevoNombre string) {
    u.Nombre = nuevoNombre  // Requiere puntero para modificar
}

func main() {
    usuario := Usuario{"Juan", "juan@example.com", 25, false}

    fmt.Println(usuario.EsAdulto())     // true
    fmt.Println(usuario.Saludar())      // "Hola Juan"

    usuario.ModificarNombre("Carlos")
    fmt.Println(usuario.Nombre)         // "Carlos"
}
```

### Tipos Basados en Interfaces

```go
// Interfaz para procesar datos
type Procesador interface {
    Procesar(datos []byte) error
    ObtenerResultado() string
}

// Tipo que implementa la interfaz
type ProcesadorJSON struct {
    datos map[string]interface{}
}

func (p *ProcesadorJSON) Procesar(datos []byte) error {
    return json.Unmarshal(datos, &p.datos)
}

func (p *ProcesadorJSON) ObtenerResultado() string {
    bytes, _ := json.Marshal(p.datos)
    return string(bytes)
}
```

### Tipos Basados en Slices y Maps

```go
// Tipo basado en slice
type Numeros []int

// Métodos útiles
func (n Numeros) Suma() int {
    total := 0
    for _, num := range n {
        total += num
    }
    return total
}

func (n Numeros) Promedio() float64 {
    if len(n) == 0 {
        return 0
    }
    return float64(n.Suma()) / float64(len(n))
}

// Tipo basado en map
type Inventario map[string]int

func (i Inventario) Agregar(producto string, cantidad int) {
    i[producto] += cantidad
}

func (i Inventario) Total() int {
    total := 0
    for _, cantidad := range i {
        total += cantidad
    }
    return total
}

func main() {
    numeros := Numeros{1, 2, 3, 4, 5}
    fmt.Println(numeros.Suma())     // 15
    fmt.Println(numeros.Promedio()) // 3.0

    inv := Inventario{
        "lapices": 100,
        "cuadernos": 50,
    }
    inv.Agregar("lapices", 25)
    fmt.Println(inv.Total())  // 175
}
```

---

## 15.4 Type Aliases - Go 1.9+

### Cuándo Usar Type Aliases

**1. Migración Gradual de Código**

```go
// Versión antigua del código
type UsuarioID int

// Se descubre que es mejor usar string
// Versión nueva
type UsuarioIDV2 string

// Durante la transición (mantener compatibilidad):
type UsuarioID = UsuarioIDV2

// El código antiguo sigue funcionando sin cambios
func procesarUsuario(id UsuarioID) { ... }
```

**2. Simplificar Tipos Complejos**

```go
// Sin alias (verbose)
type Cache map[string][]interface{}

// Con alias (legible)
type Cache = map[string][]interface{}

// Uso
cache := Cache{"key": {"value1", "value2"}}
```

**3. Estándares de Paquetes**

```go
// Archivo: net/http/types.go
type Header = map[string][]string

// Uso en todo el paquete
headers := Header{"Content-Type": {"application/json"}}
```

### Diferencia Fundamental: Definition vs Alias

```go
// DEFINITION: Crea tipo nuevo, visible en godoc
type Usuario struct {
    Nombre string
}

// ALIAS: Solo renombra, más transparente
type UsuarioAlias = Usuario

func main() {
    u1 := Usuario{"Juan"}
    u2 := UsuarioAlias{"María"}

    // Tipos internamente idénticos
    // Pero semánticamente comunican diferente propósito
}
```

### Limitaciones de Type Aliases

```go
// Alias de tipo basado en primitivo
type MiInt = int

func (m MiInt) MiMetodo() {  // ✗ ERROR
    // No puedes añadir métodos a aliases de primitivos
}

// Por qué: El alias no es un tipo nuevo, es solo int
```

---

## 15.5 Conversión de Tipos

### Conversión Explícita entre Tipos

```go
type Metros float64
type Kilometros float64

func main() {
    // Conversión explícita
    distancia := Metros(100)

    // Convertir a Kilometros
    distanciaKm := Kilometros(float64(distancia) / 1000)
    fmt.Println(distanciaKm)  // 0.1

    // Dos pasos: Metros -> float64 -> Kilometros
}
```

**Fórmula de conversión:**

```
NuevoTipo(Expresión)
```

### Conversión entre Tipos Compatibles

```go
func main() {
    // int <-> float64
    var numeroInt int = 42
    var numeroFloat float64 = float64(numeroInt)
    var numeroIntAgain int = int(numeroFloat)

    // string <-> []byte
    texto := "Hola"
    bytes := []byte(texto)
    textoNuevo := string(bytes)

    // int <-> uint
    var num int = -5
    // var numUnsigned uint = uint(num)  // Cuidado con negativos
}
```

### Conversión de Tipos Definidos

```go
type Temperatura float64
type TemperaturaKelvin float64

const AbsoluteZeroCelsius = -273.15

func CelsiusAKelvin(c Temperatura) TemperaturaKelvin {
    return TemperaturaKelvin(float64(c) - AbsoluteZeroCelsius)
}

func main() {
    tempCelsius := Temperatura(0)
    tempKelvin := CelsiusAKelvin(tempCelsius)
    fmt.Println(tempKelvin)  // 273.15
}
```

### Reglas Importantes de Conversión

**1. No puedes convertir cualquier tipo a cualquier tipo**

```go
type Usuario struct { }
type Admin struct { }

func main() {
    u := Usuario{}
    // a := Admin(u)  // ✗ ERROR: No se pueden convertir structs
}
```

**2. Conversión de número a string**

```go
func main() {
    // ✗ INCORRECTO
    // texto := string(42)  // No funciona así

    // ✓ CORRECTO
    texto := strconv.Itoa(42)
    fmt.Println(texto)  // "42"
}
```

**3. Conversión de string a número**

```go
func main() {
    // string -> int
    num, err := strconv.Atoi("42")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(num, reflect.TypeOf(num))  // 42 int

    // string -> float64
    valor, err := strconv.ParseFloat("3.14", 64)
    // string -> int64
    numGrande, err := strconv.ParseInt("999999999", 10, 64)
}
```

### Casos de Uso Reales

**1. Manejo de Dinero (conversión de tipos)**

```go
type Centavos int64
type Dolares float64

func CentavosADolares(centavos Centavos) Dolares {
    return Dolares(float64(centavos) / 100.0)
}

func DolaresCentavos(dolares Dolares) Centavos {
    return Centavos(dolares * 100)
}

func main() {
    dinero := Centavos(2999)  // $29.99
    fmt.Println(CentavosADolares(dinero))  // 29.99
}
```

**2. Duración de Tiempo**

```go
type Horas int
type Minutos int

func HorasAMinutos(h Horas) Minutos {
    return Minutos(h) * 60
}

func main() {
    duracion := Horas(2)
    fmt.Println(HorasAMinutos(duracion))  // 120 minutos
}
```

---

## 15.6 Compatibilidad de Tipos

### Cuándo Dos Tipos Son Compatibles

```
┌─────────────────────────────────────────────────────┐
│  Reglas de Compatibilidad en Go                     │
├─────────────────────────────────────────────────────┤
│                                                      │
│  1. Tipos idénticos: int = int ✓                   │
│  2. Type alias: int = IntAlias (donde IntAlias=int)│
│  3. No entre type definitions diferentes            │
│     Metros ≠ Kilometros (ambos float64)            │
│  4. Especiales: nil, constantes sin tipo            │
│                                                      │
└─────────────────────────────────────────────────────┘
```

### Assignability Rules (Reglas de Asignación)

```go
func main() {
    // Básico: mismo tipo
    var a int = 5
    var b int = a  // ✓ Compatible

    // Sin type definition: son el mismo tipo
    x := 5
    var y int = x  // ✓ Compatible (int)

    // Con type definitions: tipos diferentes
    type Metros int
    var m Metros = 100
    var n int = m  // ✗ ERROR: tipos incompatibles

    // Conversión explícita
    n = int(m)  // ✓ Correcto
}
```

### Type Definitions vs Compatibilidad

```go
type Usuario struct {
    Nombre string
    Email  string
}

type UsuarioV2 struct {
    Nombre string
    Email  string
}

func main() {
    u1 := Usuario{"Juan", "juan@example.com"}

    // ✗ ERROR: Aunque tienen los mismos campos
    // u2 := UsuarioV2(u1)

    // ✓ Debes convertir campo a campo
    u2 := UsuarioV2{
        Nombre: u1.Nombre,
        Email:  u1.Email,
    }
}
```

### Interfaces y Compatibilidad

```go
// Interfaz definida
type Lector interface {
    Leer() ([]byte, error)
}

// Tipo que implementa la interfaz
type Archivo struct {
    ruta string
}

func (a Archivo) Leer() ([]byte, error) {
    return ioutil.ReadFile(a.ruta)
}

func procesarLector(l Lector) {
    datos, _ := l.Leer()
    fmt.Println(datos)
}

func main() {
    archivo := Archivo{"datos.txt"}

    // ✓ Compatible: Archivo implementa Lector
    procesarLector(archivo)

    // Structural typing: Go verifica métodos, no nombres explícitos
}
```

### Compatibilidad Numérica

```go
func main() {
    // Constantes sin tipo son flexibles
    const numero = 42

    var a int = numero      // ✓
    var b int64 = numero    // ✓
    var c float64 = numero  // ✓

    // Variables con tipo son estrictas
    var x int = 5
    var y int64 = x  // ✗ ERROR
    var z int64 = int64(x)  // ✓ Correcto
}
```

---

## 15.7 Type Assertions

### Qué es una Type Assertion

Una **type assertion** es una operación que dice "estoy seguro de que esta interfaz contiene este tipo específico". La sintaxis es `x.(T)`.

```go
var x interface{} = "hola"

// Type assertion sin verificación (peligroso)
texto := x.(string)  // Funciona porque x es string
fmt.Println(texto)   // "hola"

// Type assertion incorrecta (panic)
numero := x.(int)    // ✗ PANIC: x no es int
```

### The Comma-Ok Pattern (Verificación Segura)

```go
// INCORRECTO: Sin verificación
func procesar(v interface{}) {
    texto := v.(string)  // PANIC si no es string
}

// CORRECTO: Con verificación
func procesar(v interface{}) {
    texto, ok := v.(string)  // ok = true/false
    if !ok {
        fmt.Println("No es string")
        return
    }
    fmt.Println(texto)
}
```

**Patrón: comma-ok**

```
valor, ok := interfaz.(TipoEsperado)

Si ok=true:  valor contiene el datos, tipo correcto
Si ok=false: valor tiene valor cero del tipo, interfaz era diferente
```

### Ejemplos Prácticos

```go
// Guardar en interfaz
var dato interface{} = 42

// Tipo assertion sin verificación
fmt.Println(dato.(int))  // 42

// Tipo assertion con verificación
if num, ok := dato.(int); ok {
    fmt.Println("Es un número:", num)
} else {
    fmt.Println("No es un número")
}

// Con string
dato = "Hola"
if texto, ok := dato.(string); ok {
    fmt.Println("Es texto:", texto)
}
```

### Casos de Uso Reales

**1. Parsing JSON (Desconocemos el tipo)**

```go
import "encoding/json"

func main() {
    jsonData := []byte(`{
        "nombre": "Juan",
        "edad": 30,
        "activo": true
    }`)

    var data map[string]interface{}
    json.Unmarshal(jsonData, &data)

    // Extraer valores con type assertion
    if nombre, ok := data["nombre"].(string); ok {
        fmt.Println("Nombre:", nombre)
    }

    if edad, ok := data["edad"].(float64); ok {
        fmt.Println("Edad:", int(edad))
    }
}
```

**2. Manejo de Errores**

```go
func manejarError(err error) {
    // Verificar si es un tipo específico de error
    if fsErr, ok := err.(*os.PathError); ok {
        fmt.Println("Error de archivo:", fsErr.Path)
    } else if netErr, ok := err.(net.Error); ok {
        fmt.Println("Error de red:", netErr)
    } else {
        fmt.Println("Error genérico:", err)
    }
}
```

**3. Trabajar con Plugins o Extensiones**

```go
type Plugin interface {
    Execute(config map[string]interface{})
}

func ejecutarPlugin(plugin Plugin, config map[string]interface{}) {
    // Verificar si el plugin es de un tipo específico
    if debugPlugin, ok := plugin.(*DebugPlugin); ok {
        debugPlugin.SetLogLevel("DEBUG")
    }

    plugin.Execute(config)
}
```

### Antipatrones: Cuándo NO Usar Type Assertions

```go
// ✗ MAL: Type assertions excesivas
func procesarDatos(v interface{}) {
    if str, ok := v.(string); ok {
        // ...
    } else if num, ok := v.(int); ok {
        // ...
    } else if arr, ok := v.([]interface{}); ok {
        // ...
    } else if m, ok := v.(map[string]interface{}); ok {
        // ...
    }
}

// ✓ BIEN: Usar type switches (siguiente sección)
func procesarDatos(v interface{}) {
    switch v := v.(type) {
    case string:
        // ...
    case int:
        // ...
    case []interface{}:
        // ...
    case map[string]interface{}:
        // ...
    }
}
```

---

## 15.8 Type Switches

### Qué es un Type Switch

Un **type switch** es como un switch normal, pero en lugar de comparar valores, compara **tipos**.

```go
switch v := valor.(type) {
case int:
    fmt.Println("Es un int:", v)
case string:
    fmt.Println("Es un string:", v)
case bool:
    fmt.Println("Es un bool:", v)
default:
    fmt.Println("Tipo desconocido:", reflect.TypeOf(valor))
}
```

### Estructura Básica

```
┌─────────────────────────────────────────────┐
│  switch x := interfaz.(type) {              │
│      case Tipo1:                           │
│          // x es tipo Tipo1 aquí           │
│      case Tipo2:                           │
│          // x es tipo Tipo2 aquí           │
│      default:                              │
│          // x es del tipo original         │
│  }                                         │
└─────────────────────────────────────────────┘
```

### Ejemplos Progresivos

**Nivel 1: Tipos Básicos**

```go
func procesarValor(v interface{}) {
    switch val := v.(type) {
    case int:
        fmt.Printf("Número entero: %d\n", val)

    case float64:
        fmt.Printf("Número decimal: %.2f\n", val)

    case string:
        fmt.Printf("Texto: %s\n", val)

    case bool:
        if val {
            fmt.Println("Verdadero")
        } else {
            fmt.Println("Falso")
        }

    default:
        fmt.Println("Tipo desconocido")
    }
}

func main() {
    procesarValor(42)
    procesarValor(3.14)
    procesarValor("Hola")
    procesarValor(true)
}
```

**Nivel 2: Tipos Personalizados**

```go
type Usuario struct {
    Nombre string
}

type Producto struct {
    Nombre string
    Precio float64
}

func procesar(entidad interface{}) {
    switch e := entidad.(type) {
    case *Usuario:
        fmt.Println("Usuario:", e.Nombre)

    case *Producto:
        fmt.Printf("Producto: %s ($%.2f)\n", e.Nombre, e.Precio)

    case []int:
        fmt.Println("Slice de números:", e)

    default:
        fmt.Println("Entidad desconocida")
    }
}

func main() {
    procesar(&Usuario{"Juan"})
    procesar(&Producto{"Laptop", 999.99})
    procesar([]int{1, 2, 3, 4, 5})
}
```

**Nivel 3: Interfaces**

```go
type Imprimible interface {
    String() string
}

type Usuario struct {
    Nombre string
}

func (u Usuario) String() string {
    return "Usuario: " + u.Nombre
}

type Log struct {
    Mensaje string
}

func (l Log) String() string {
    return "LOG: " + l.Mensaje
}

func mostrar(v interface{}) {
    switch val := v.(type) {
    case Imprimible:
        fmt.Println(val.String())

    case fmt.Stringer:  // Go proporciona esta interfaz
        fmt.Println(val.String())

    default:
        fmt.Println(v)
    }
}

func main() {
    mostrar(Usuario{"María"})
    mostrar(Log{"Sistema iniciado"})
}
```

### Type Switch vs Type Assertion

```go
// TYPE ASSERTION: Una verificación
if str, ok := v.(string); ok {
    fmt.Println(str)
}

// TYPE SWITCH: Múltiples verificaciones elegantes
switch v := v.(type) {
case string:
    fmt.Println(v)
case int:
    fmt.Println(v * 2)
case []string:
    fmt.Println(strings.Join(v, ", "))
}
```

### Casos de Uso Reales

**1. Dispatcher de Eventos**

```go
type Evento interface{}

type EventoClick struct {
    X, Y int
}

type EventoTeclado struct {
    Tecla string
}

func procesarEvento(evento Evento) {
    switch e := evento.(type) {
    case EventoClick:
        fmt.Printf("Click en (%d, %d)\n", e.X, e.Y)

    case EventoTeclado:
        fmt.Printf("Tecla presionada: %s\n", e.Tecla)

    default:
        fmt.Println("Evento desconocido")
    }
}
```

**2. Parser de Configuración**

```go
func parsearConfiguracion(valor interface{}) string {
    switch v := valor.(type) {
    case string:
        return v

    case int:
        return strconv.Itoa(v)

    case bool:
        if v {
            return "true"
        }
        return "false"

    case []interface{}:
        items := make([]string, len(v))
        for i, item := range v {
            items[i] = parsearConfiguracion(item)
        }
        return "[" + strings.Join(items, ", ") + "]"

    default:
        return fmt.Sprint(v)
    }
}
```

---

## 15.9 Type Embedding

### Composición vs Herencia

Go rechaza deliberadamente la herencia. En su lugar, usa **composición**. Type embedding es una forma de lograr composición de tipos.

```
┌──────────────────────────────────────┐
│  Herencia (Otros lenguajes)          │
├──────────────────────────────────────┤
│  class Perro extends Animal {        │
│      ladrido()                       │
│  }                                   │
│                                      │
│  ✗ Jerarquía rígida                 │
│  ✗ Acoplamiento fuerte              │
│  ✗ Frágil base class problem        │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  Composición (Go style)              │
├──────────────────────────────────────┤
│  type Perro struct {                 │
│      Animal                          │
│  }                                   │
│                                      │
│  ✓ Flexible                         │
│  ✓ Acoplamiento débil               │
│  ✓ Composición explícita            │
└──────────────────────────────────────┘
```

### Type Embedding Básico

```go
// Tipo base
type Animal struct {
    Nombre string
    Edad   int
}

// Método en Animal
func (a Animal) Hablar() string {
    return a.Nombre + " hace un sonido"
}

// Tipo que embebe Animal
type Perro struct {
    Animal        // Embedding (campo sin nombre)
    Raza   string
}

func main() {
    perro := Perro{
        Animal: Animal{Nombre: "Rex", Edad: 3},
        Raza:   "Pastor Alemán",
    }

    // Acceso a campos del tipo embebido
    fmt.Println(perro.Nombre)      // "Rex"
    fmt.Println(perro.Edad)         // 3

    // Acceso a métodos del tipo embebido
    fmt.Println(perro.Hablar())     // "Rex hace un sonido"

    // Acceso explícito
    fmt.Println(perro.Animal.Nombre)  // "Rex"
}
```

### Diferencia: Embedding vs Campo Normal

```go
// ✗ INCORRECTO: Campo con nombre
type Gato struct {
    animal Animal  // Tiene nombre
}

g := Gato{
    animal: Animal{"Misi", 2},
}
g.Hablar()  // ✗ ERROR: No hereda métodos

// ✓ CORRECTO: Embedding (sin nombre)
type Gato struct {
    Animal  // Sin nombre
}

g := Gato{
    Animal: Animal{"Misi", 2},
}
g.Hablar()  // ✓ Funciona: Acceso directo a métodos
```

### Promoción de Métodos

```go
type Motor struct {
    RPM int
}

func (m Motor) Acelerar() {
    m.RPM += 1000
    fmt.Println("Motor acelerado:", m.RPM)
}

type Coche struct {
    Motor      // Método Acelerar() es promovido
    Marca string
}

func main() {
    coche := Coche{
        Motor: Motor{RPM: 0},
        Marca: "Toyota",
    }

    // El método Acelerar() está promovido
    coche.Acelerar()  // ✓ Funciona sin especificar Motor
}
```

### Shadowing: Sobrescribir Métodos Promovidos

```go
type Animal struct {
    Nombre string
}

func (a Animal) Sonido() string {
    return "Sonido genérico"
}

type Perro struct {
    Animal
}

// Sobrescribir método promovido
func (p Perro) Sonido() string {
    return p.Nombre + " ladra"
}

func main() {
    perro := Perro{Animal{"Rex"}}
    fmt.Println(perro.Sonido())          // "Rex ladra"
    fmt.Println(perro.Animal.Sonido())   // "Sonido genérico"
}
```

### Múltiple Embedding

```go
type Corredor struct {
    Velocidad float64
}

func (c Corredor) Correr() string {
    return fmt.Sprintf("Corre a %.1f km/h", c.Velocidad)
}

type Nadador struct {
    Profundidad float64
}

func (n Nadador) Nadar() string {
    return fmt.Sprintf("Nada a %.1f metros", n.Profundidad)
}

type Deportista struct {
    Corredor
    Nadador
    Nombre string
}

func main() {
    deportista := Deportista{
        Corredor: Corredor{20.5},
        Nadador:  Nadador{50.0},
        Nombre:   "Michael",
    }

    fmt.Println(deportista.Correr())     // Método de Corredor
    fmt.Println(deportista.Nadar())      // Método de Nadador
}
```

### Embedding de Punteros

```go
type Configuracion struct {
    Host string
    Port int
}

// Embeber un puntero
type Servidor struct {
    *Configuracion  // Puntero embebido
}

func main() {
    config := &Configuracion{"localhost", 8080}
    servidor := Servidor{config}

    fmt.Println(servidor.Host)      // "localhost"
    fmt.Println(servidor.Port)      // 8080
}
```

---

## 15.10 Métodos con Tipos Diferentes

### Receptores: Valor vs Puntero

```go
type Contador struct {
    valor int
}

// Receptor por VALOR (copia)
func (c Contador) Incrementar() {
    c.valor++  // Modifica la copia, no el original
}

// Receptor por PUNTERO (referencia)
func (c *Contador) IncrementarV2() {
    c.valor++  // Modifica el original
}

func main() {
    contador := Contador{0}

    contador.Incrementar()   // No efectivo
    fmt.Println(contador.valor)  // 0 (sin cambios)

    contador.IncrementarV2()  // Efectivo
    fmt.Println(contador.valor)  // 1 (cambió)
}
```

### El Method Set

Cada tipo tiene un **method set**: el conjunto de métodos que puede llamar.

```
┌──────────────────────────────────────────────┐
│  Method Set Rules                            │
├──────────────────────────────────────────────┤
│  Tipo T:      Puede llamar métodos (T)      │
│  Tipo T:      Puede llamar métodos (*T)     │
│  Tipo *T:     Puede llamar métodos (*T)    │
│  Tipo *T:     Puede llamar métodos (T)*    │
│               (* Go auto-dereferences)      │
│                                              │
│  Interfaz:    Require exactamente (T)      │
│               o exactamente (*T)            │
│               No ambos                      │
└──────────────────────────────────────────────┘
```

### Ejemplo Práctico: Interfaz de Lectura

```go
type Lector interface {
    Leer() ([]byte, error)
}

// Implementación con receptor por valor
type Buffer1 struct {
    datos []byte
}

func (b Buffer1) Leer() ([]byte, error) {
    return b.datos, nil
}

// Implementación con receptor por puntero
type Buffer2 struct {
    datos []byte
}

func (b *Buffer2) Leer() ([]byte, error) {
    return b.datos, nil
}

func procesarLector(l Lector) {
    datos, _ := l.Leer()
    fmt.Println(datos)
}

func main() {
    b1 := Buffer1{[]byte("Hola")}
    procesarLector(b1)  // ✓ Funciona (valor)

    b2 := Buffer2{[]byte("Mundo")}
    procesarLector(&b2)  // ✓ Funciona (puntero)

    // procesarLector(b2)  // ✗ ERROR: puntero esperado
}
```

### Casos de Uso: Valor vs Puntero

**Usar Receptor por VALOR cuando:**

```go
// El tipo es pequeño
type Punto struct {
    X, Y float64
}

func (p Punto) Distancia() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// El tipo no se modifica
type Usuario struct {
    Nombre string
}

func (u Usuario) Saludar() string {
    return "Hola " + u.Nombre
}
```

**Usar Receptor por PUNTERO cuando:**

```go
// Necesitas modificar el tipo
type Contador struct {
    valor int
}

func (c *Contador) Incrementar() {
    c.valor++
}

// El tipo es grande (evita copias)
type BaseDatos struct {
    Conexion *sql.DB
    Cache    map[string]string
    // ... muchos más campos
}

func (bd *BaseDatos) Guardar(clave string, valor string) error {
    bd.Cache[clave] = valor
    // ...
}
```

### Consistencia: Todos Valor o Todos Puntero

**Recomendación Go:** Si algún método de un tipo necesita un receptor puntero, **todos** los métodos deben ser receptores puntero.

```go
// ✓ BIEN: Consistente
type Carrito struct {
    items []string
}

func (c *Carrito) AgregarItem(item string) {
    c.items = append(c.items, item)
}

func (c *Carrito) ObtenerTotal() int {
    return len(c.items)
}

// ✗ EVITAR: Mezclar receptores
type Carrito2 struct {
    items []string
}

func (c Carrito2) ObtenerTotal() int {
    return len(c.items)
}

func (c *Carrito2) AgregarItem(item string) {
    c.items = append(c.items, item)
}
```

---

## 15.11 Buenas Prácticas y Diseño de Tipos

### Nomenclatura de Tipos

**Convenciones Go:**

```go
// ✓ BIEN: Nombres cortos, descriptivos, CamelCase
type Usuario struct { }
type Handler struct { }
type Procesador interface { }

// ✗ EVITAR: Demasiado corto
type U struct { }
type H struct { }

// ✗ EVITAR: Nombres de tipo en el nombre
type UsuarioStruct struct { }
type HandlerInterface interface { }

// ✗ EVITAR: Nombres genéricos
type Cosa struct { }
type Datos struct { }
```

### Cuándo Definir un Tipo Nuevo

**1. Semántica Clara**

```go
// ✓ BIEN: El nombre dice qué es
type Temperatura int
type Velocidad float64
type UsuarioID int

// ✗ EVITAR: Solo es un int
type Numero int
```

**2. Métodos Específicos**

```go
// ✓ BIEN: El tipo tiene métodos cohesivos
type Dinero int64

func (d Dinero) AñadirImpuesto(tasa float64) Dinero {
    return Dinero(float64(d) * (1 + tasa))
}

func (d Dinero) Formatear() string {
    return "$" + strconv.Itoa(int(d))
}
```

**3. Protección contra Errores**

```go
// ✓ BIEN: Previene bugs
type Email string
type Password string

// Diferentes tipos = imposible confundir
func ValidarLogin(email Email, password Password) error {
    // ...
}
```

### Antipatrones: Qué Evitar

**1. Abuso de interface{}**

```go
// ✗ MAL: Demasiado genérico
func Procesar(datos interface{}) interface{} {
    // Necesitas type assertions en todas partes
}

// ✓ BIEN: Tipos específicos
func Procesar(datos []byte) ([]byte, error) {
    // Claro qué espera y retorna
}
```

**2. Type Definitions Innecesarias**

```go
// ✗ MAL: No añade valor
type NumeroEntero int

// ✓ BIEN: Necesita semántica
type Edad int
```

**3. Jerarquías Profundas de Embedding**

```go
// ✗ MAL: Jerarquía confusa
type A struct { }
type B struct { A }
type C struct { B }
type D struct { C }

// ✓ BIEN: Composición clara
type D struct {
    a *A
    b *B
    c *C
}
```

### Diseño Basado en Comportamiento

Go favorece el diseño basado en **interfaces y comportamiento**, no en jerarquías de tipos.

```go
// ✗ EVITAR: Pensar en jerarquía
/*
         Animal (base)
           /  \
        Perro Gato
*/

// ✓ MEJOR: Pensar en interfaces
type Sonidoso interface {
    HacerSonido() string
}

type Comestible interface {
    Comer(comida string)
}

// Cualquier tipo puede implementar estas interfaces
type Perro struct { }
func (p Perro) HacerSonido() string { return "Guau" }
func (p Perro) Comer(comida string) { }
```

### Generics sin Generics (Antes de Go 1.18)

```go
// Antes de Go 1.18, ibas "al aire" con interface{}
type Contenedor struct {
    items []interface{}
}

func (c *Contenedor) Agregar(item interface{}) {
    c.items = append(c.items, item)
}

func (c *Contenedor) Obtener(i int) (interface{}, error) {
    if i < 0 || i >= len(c.items) {
        return nil, errors.New("índice fuera de rango")
    }
    return c.items[i], nil
}

// Uso:
contenedor := &Contenedor{}
contenedor.Agregar("texto")
contenedor.Agregar(42)

valor, _ := contenedor.Obtener(0)
texto := valor.(string)  // Type assertion necesaria
```

### Casos de Uso Reales

**1. Dinero Seguro**

```go
type Moneda string

const (
    USD Moneda = "USD"
    EUR Moneda = "EUR"
)

type Dinero struct {
    Cantidad int64   // Centavos para precisión
    Moneda   Moneda
}

func (d Dinero) AgregarImpuesto(tasa float64) Dinero {
    return Dinero{
        Cantidad: int64(float64(d.Cantidad) * (1 + tasa)),
        Moneda:   d.Moneda,
    }
}
```

**2. Tipos Seguros de ID**

```go
type UsuarioID int64
type ProductoID int64

// El compilador previene confusiones
func ObtenerUsuario(id UsuarioID) (*Usuario, error) {
    // ...
}

// Esto es un error en compilación:
// usuario := ObtenerUsuario(productoID)
```

**3. Estados y Enums**

```go
type EstadoOrden string

const (
    Pendiente EstadoOrden = "PENDIENTE"
    Enviada   EstadoOrden = "ENVIADA"
    Entregada EstadoOrden = "ENTREGADA"
)

type Orden struct {
    ID    int
    Estado EstadoOrden
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Sistema de Unidades de Medida

**Objetivo:** Crear tipos personalizados para distancia y tiempo con conversiones seguras.

**Tarea:**

1. Define tipos `Metros`, `Kilómetros`, `Segundos`, `Minutos`
2. Implementa métodos de conversión entre unidades
3. Calcula velocidad (distancia/tiempo)
4. Implementa validación (no valores negativos)

**Plantilla:**

```go
package main

import "fmt"

type Metros float64
type Kilómetros float64
type Segundos int
type Minutos int

// TODO: Implementar métodos de conversión
// MetrosAKilometros(), KilometrosAMetros(), etc.

// TODO: Calcular velocidad
// velocidad := distancia / tiempo

func main() {
    distancia := Metros(1500)
    tiempo := Segundos(300)

    // TODO: Mostrar distancia en km
    // TODO: Mostrar tiempo en minutos
    // TODO: Calcular y mostrar velocidad
}
```

**Requisitos:**

- [ ] 4 tipos definidos
- [ ] Métodos de conversión bidireccional
- [ ] Función de velocidad
- [ ] Validación de valores positivos
- [ ] Salida formateada correctamente

---

### Ejercicio 2: Sistema de Validación Segura

**Objetivo:** Crear tipos seguros basados en string con validación integrada.

**Tarea:**

1. Define tipos `Email`, `Password`, `Nombre`
2. Cada tipo debe tener métodos de validación
3. Implementa un tipo `Usuario` que combine todos
4. Valida antes de guardar en la "base de datos"

**Plantilla:**

```go
package main

import (
    "fmt"
    "strings"
)

type Email string
type Password string
type Nombre string

// TODO: Validadores
// (e Email) EsValido() bool
// (p Password) EsSegura() bool
// (n Nombre) EsValido() bool

type Usuario struct {
    Email    Email
    Password Password
    Nombre   Nombre
}

// TODO: (u *Usuario) Guardar() error

func main() {
    // TODO: Crear usuario con datos válidos
    // TODO: Intentar crear usuario con email inválido
    // TODO: Verificar password segura
}
```

**Requisitos:**

- [ ] Email valida formato @
- [ ] Password tiene mínimo 8 caracteres
- [ ] Nombre no está vacío
- [ ] Método Guardar() valida todo
- [ ] Manejo de errores con mensajes claros

---

### Ejercicio 3: Conversor Universal con Type Switches

**Objetivo:** Implementar un conversor que maneje múltiples tipos.

**Tarea:**

1. Crear función `Convertir` que acepta `interface{}`
2. Implementar conversiones: string ↔ int ↔ float ↔ bool
3. Usar type switches para detectar tipo origen
4. Retornar error si conversión es imposible

**Plantilla:**

```go
package main

import (
    "fmt"
    "strconv"
)

func ConvertirAString(v interface{}) (string, error) {
    // TODO: Type switch para convertir a string
    // string -> string
    // int -> string
    // float64 -> string
    // bool -> string
    // default -> error
}

func ConvertirAInt(v interface{}) (int, error) {
    // TODO: Conversiones a int
}

func main() {
    // TODO: Probar conversiones exitosas
    // TODO: Probar conversiones fallidas
}
```

**Requisitos:**

- [ ] 3 funciones de conversión (string, int, float64)
- [ ] Type switches con múltiples casos
- [ ] Manejo de errores
- [ ] Pruebas con valores válidos e inválidos

---

### Ejercicio 4: Tipos de Dominio - Dinero y Coordenadas

**Objetivo:** Crear tipos complejos con métodos y validación.

**Tarea:**

1. Tipo `Dinero` con cantidad y moneda
2. Tipo `Coordenadas` con latitud y longitud
3. Métodos: conversión de moneda, distancia entre puntos
4. Embedding de un tipo en otro (Ubicacion embebe Coordenadas)

**Plantilla:**

```go
package main

import (
    "fmt"
    "math"
)

type Moneda string
type Dinero struct {
    Cantidad int64
    Moneda   Moneda
}

// TODO: (d Dinero) EnOtraMoneda(nuevaMoneda Moneda, tasa float64) Dinero

type Coordenadas struct {
    Latitud  float64
    Longitud float64
}

// TODO: (c Coordenadas) DistanciaA(otra Coordenadas) float64

type Ubicacion struct {
    Coordenadas  // Embedding
    Nombre       string
}

func main() {
    // TODO: Crear dineros en diferentes monedas
    // TODO: Convertir entre monedas
    // TODO: Calcular distancia entre coordenadas
}
```

**Requisitos:**

- [ ] Tipos `Dinero` y `Coordenadas` con validación
- [ ] Métodos de conversión
- [ ] Cálculo de distancia haversine
- [ ] Embedding de `Coordenadas` en `Ubicacion`
- [ ] Pruebas con datos reales

---

### Ejercicio 5: Contenedor Genérico Seguro sin Genéricos

**Objetivo:** Implementar contenedor type-safe usando type assertions y switches.

**Tarea:**

1. Crear tipo `Contenedor` que almacena `interface{}`
2. Métodos `Agregar`, `Obtener` con type checking
3. Implementar iterador que valida tipos
4. Usar type switches para operaciones seguras

**Plantilla:**

```go
package main

import "fmt"

type Contenedor struct {
    items      []interface{}
    tipoEsperado string  // "int", "string", etc.
}

// TODO: func NewContenedor(tipoEsperado string) *Contenedor

// TODO: func (c *Contenedor) Agregar(item interface{}) error
//       (Verifica que item es del tipo esperado)

// TODO: func (c *Contenedor) Obtener(indice int) (interface{}, error)

// TODO: func (c *Contenedor) ObtenerComo(indice int, destino interface{}) error

// TODO: func (c *Contenedor) Iterar(f func(interface{}) error) error

func main() {
    // TODO: Crear contenedor de strings
    // TODO: Agregar strings válidos
    // TODO: Intentar agregar int (debe fallar)
    // TODO: Iterar y procesar elementos
}
```

**Requisitos:**

- [ ] Validación de tipos en Agregar()
- [ ] Type switches en Iterar()
- [ ] Mensajes de error descriptivos
- [ ] Manejo seguro con errores
- [ ] Pruebas de casos válidos e inválidos

---

## Resumen y Puntos Clave

1. **Sistema de Tipos**: Go es estáticamente tipado con inferencia automática
2. **Type Definitions**: Crean tipos nuevos; son distintos del tipo base
3. **Type Aliases**: Crean alias del tipo existente; compatibles sin conversión
4. **Conversión**: Explícita con sintaxis `NuevoTipo(valor)`
5. **Type Assertions**: `x.(T)` para verificar tipo en interfaz; usa comma-ok
6. **Type Switches**: `switch v := v.(type)` para múltiples tipos elegantemente
7. **Embedding**: Composición sobre herencia; promueve métodos automáticamente
8. **Method Set**: Determina qué métodos puede llamar cada tipo (valor vs puntero)
9. **Buenas Prácticas**: Nombres claros, semántica cohesiva, consistencia de receptores
10. **Go Philosophy**: Simplicidad, composición, interfaces para flexibilidad

---

**Siguiente capítulo:** CAPÍTULO 16: INTERFACES - CONTRATOS IMPLÍCITOS Y POLIMORFISMO EN GO

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/15-sistema-de-tipos/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/15-sistema-de-tipos):

```bash
cd examples/15-sistema-de-tipos
go run .
```
