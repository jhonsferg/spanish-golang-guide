# Capítulo 4: Variables, constantes y tipos básicos

## Índice del Capítulo 4

1. [4.1 ¿Qué es una Variable?](#41-qué-es-una-variable)
2. [4.2 Declaración de Variables](#42-declaración-de-variables)
3. [4.3 Cero Valores (Zero Values)](#43-cero-valores-zero-values)
4. [4.4 Tipos Básicos de Datos](#44-tipos-básicos-de-datos)
5. [4.5 Constantes](#45-constantes)
6. [4.6 Conversiones de Tipos](#46-conversiones-de-tipos)
7. [4.7 Strings y Runes](#47-strings-y-runes)
8. [4.8 Bytes y Runes Profundo](#48-bytes-y-runes-profundo)
9. [4.9 Paquete strconv - Conversiones](#49-paquete-strconv---conversiones)
10. [4.10 Buenas Prácticas de Variables](#410-buenas-prácticas-de-variables)

---

## 4.1 ¿Qué es una Variable?

### Definición Conceptual

Una variable es una **ubicación en memoria** que almacena un valor. Tiene tres propiedades:

```

 VARIABLE                             │

 1. NOMBRE                            │
    └─ Identificador: edad, nombre    │
                                      │
 2. TIPO                              │
    └─ Qué tipo de dato: int, string  │
                                      │
 3. VALOR                             │
    └─ Lo que almacena: 25, "Juan"    │
                                      │
 4. DIRECCIÓN (bonus)                 │
    └─ Dónde está en memoria: 0x2031a│

```

### Cómo Go Ve las Variables

```go
edad := 25
```

**Internamente:**

```

 Símbolo:  edad                      │

 Tipo:     int                       │
 Valor:    25                        │
 Ubicación: 0x2031a (en memoria RAM) │

```

### Por Qué Necesitas Tipos en Go

Go es **estrictamente tipado**. Esto significa:

```
Regla 1: Cada variable tiene un tipo
Regla 2: No hay conversiones implícitas
Regla 3: No puedes asignar un tipo a otro sin conversión explícita
```

**Comparación:**

```go
// Go (ESTRICTO)
var x int = 5
var y string = x    // ❌ ERROR: cannot use x (int) as string type

// Python (FLEXIBLE)
x = 5
y = x               # ✅ OK: 5 se convierte automáticamente

// C++ (INTERMEDIO)
int x = 5;
string y = x;       // ❌ ERROR (similar a Go)
string y = std::to_string(x);  // ✅ OK (conversión explícita)
```

**¿Por qué Go es estricto?**

```
Razones de diseño de Go:

1. SEGURIDAD
   └─ Previene bugs sutiles
 Si escribes código tipo-correcto, probablemente sea correcto   └

2. PERFORMANCE
   └─ Compilador sabe exactamente qué hacer
   └─ No necesita chequeos en runtime

3. CLARIDAD
   └─ Lector del código sabe qué tipo es cada variable
   └─ Menos sorpresas

Analogía: Es como un jardinero que dice "plantar rosas AQUÍ"
vs decir "plantar algo AQUÍ". Más específico = mejor.
```

---

## 4.2 Declaración de Variables

### Sintaxis Básica

**Opción 1: Declaración y Asignación**

```go
var nombre string = "Juan"
```

**Desglose:**
```
var      └─ Palabra clave "declare variable"
nombre   └─ Nombre de la variable
string   └─ Tipo explícito
=        └─ Asignar valor
"Juan"   └─ Valor inicial
```

**Opción 2: Inferencia de Tipos**

```go
var nombre = "Juan"
```

Go **infiere** que es string (porque "Juan" es una string literal).

**Opción 3: Short Declaration (MÁS COMÚN)**

```go
nombre := "Juan"
```

**Reglas:**
- `:=` solo funciona DENTRO de funciones
- `:=` es mucho más conciso que `var`
- Recomendado en 99% de casos

### Declaración Múltiple

**Variables múltiples de un tipo:**

```go
var x, y, z int
x = 1
y = 2
z = 3

// O más conciso
var a, b, c int = 1, 2, 3
```

**Variables múltiples de tipos diferentes:**

```go
var (
    nombre string = "Juan"
    edad   int    = 25
    altura float64 = 1.75
    activo bool   = true
)

// O con short declaration
nombre, edad := "Juan", 25  // Solo 2 variables simultáneamente
```

### Escopo de Variables

Variables solo existen dentro de su **scope** (alcance):

```go
package main

import "fmt"

var global string = "Soy global"  // Scope: TODO el package

func main() {
    var local string = "Soy local"     // Scope: dentro de main
    
    if true {
        var muyLocal string = "Scope: dentro del if"
        fmt.Println(muyLocal)
    }
    
    // fmt.Println(muyLocal)  // ❌ ERROR: muyLocal no existe aquí
}

func otraFuncion() {
    fmt.Println(global)        // ✅ OK: global es accesible
    // fmt.Println(local)       // ❌ ERROR: local es de main, no existe aquí
}
```

**Regla de Oro:**

```
Una variable solo existe desde su DECLARACIÓN
hasta el FIN del bloque en el que fue declarada.

Bloques:
 package (variables globales, accesibles por todo el package)
 function (variables de función)
 if / for / switch (variables de bloque)
 {} (bloque explícito)
```

### Shadowing (Sombra de Variables)

```go
var x int = 5

{
    var x int = 10      // NUEVA variable x, no la anterior
    fmt.Println(x)      // Output: 10
}

fmt.Println(x)          // Output: 5 (la x original)
```

**Advertencia:**

```
Shadowing es LEGAL pero CONFUSO. Evítalo en código real.

Ejemplo malo:
    func proceso() {
        user := "admin"
        {
            user := "guest"   // ❌ Confuso: parece igual, pero es otra
            fmt.Println(user) // Output: guest
        }
        fmt.Println(user)     // Output: admin ← ¿Esperas admin o guest?
    }

Mejor: Usa nombres diferentes
    func proceso() {
        user := "admin"
        {
            tempUser := "guest"
            fmt.Println(tempUser)
        }
        fmt.Println(user)
    }
```

---

## 4.3 Cero Valores (Zero Values)

### ¿Qué es un Cero Valor?

Cuando declaras una variable sin asignar valor, Go la inicializa a un **cero valor** (default):

```go
var edad int
fmt.Println(edad)           // 0

var nombre string
fmt.Println(nombre)         // "" (vacío)

var activo bool
fmt.Println(activo)         // false

var numero float64
fmt.Println(numero)         // 0
```

### Tabla de Cero Valores

```

 Tipo         │ Cero Valor  │
clear
 int (cualquier)   │ 0           │
 float (cualquier) │ 0           │
 string       │ ""          │
 bool         │ false       │
 pointer      │ nil         │
 slice        │ nil         │
 map          │ nil         │
 channel      │ nil         │
 interface{}  │ nil         │
 function     │ nil         │

```

### Por Qué Importa

**En Go, TODAS las variables están inicializadas:**

```go
var edad int           // No es basura, es 0
var nombre string      // No es basura, es ""
var activo bool        // No es basura, es false
```

**Comparación con otros lenguajes:**

```
C/C++:
    int age;           // ❌ Contiene "basura" (valor aleatorio)
    // Debes inicializar manualmente

Go:
    var age int        // ✅ Contiene 0 (seguro)
    // Go lo inicializa automáticamente
```

**Ventaja:**

```
No hay sorpresas. Go garantiza que toda variable tiene valor válido.
```

### Usar Cero Valores Estratégicamente

```go
// Usar cero valores es idiomatic en Go

// Ejemplo 1: Contador que empieza en 0
var contador int        // 0 - perfecto para contar

// Ejemplo 2: Flag que empieza en false
var procesado bool      // false - perfecto para flags

// Ejemplo 3: Text que empieza vacío
var mensaje string      // "" - perfecto para acumular texto
```

---

## 4.4 Tipos Básicos de Datos

### Categorías

```

 TIPOS BÁSICOS DE GO            │

                                │
 1. NÚMEROS ENTEROS             │
    ├─ int (plataforma)         │
    ├─ int8, int16, int32, int64│
    ├─ uint (sin signo)         │
    └─ uint8, uint16, uint32    │
                                │
 2. NÚMEROS DECIMALES           │
    ├─ float32                  │
    └─ float64                  │
                                │
 3. COMPLEJOS (raro)            │
    ├─ complex64                │
    └─ complex128               │
                                │
 4. TEXTO                       │
    ├─ string                   │
    └─ rune (carácter Unicode)  │
    └─ byte (byte simple)       │
                                │
 5. BOOLEANO                    │
    └─ bool                     │
                                │

```

### Enteros (Integers)

**int vs int64 vs int8:**

```go
// int - Tamaño depende de plataforma
var x int = 9_223_372_036_854_775_807  // 64-bit en máquinas modernas

// int64 - Exactamente 64 bits
var y int64 = 9_223_372_036_854_775_807

// int8 - Exactamente 8 bits
var z int8 = 127                        // Máximo: 127 (-128 a 127)
```

**Rango de valores:**

```

 Tipo        │ Mínimo                │ Máximo                │
clear
 int8        │ -128                  │ 127                   │
 int16       │ -32,768               │ 32,767                │
 int32       │ -2,147,483,648        │ 2,147,483,647         │
 int64       │ -9,223,372,036,854... │ 9,223,372,036,854...  │
 uint8       │ 0                     │ 255                   │
 uint16      │ 0                     │ 65,535                │
 uint32      │ 0                     │ 4,294,967,295         │
 uint64      │ 0                     │ 18,446,744,073,709... │

```

**¿Qué tipo usar?**

```
Regla 1: Usa 'int' por defecto
    var edad int = 25
    var contador int = 0

Regla 2: Usa int64 si necesitas números GRANDES
    var timestamp int64 = 1_234_567_890_123_456

Regla 3: Usa uint para números NO NEGATIVOS
    var edad uint = 25          // Nunca será negativo

Regla 4: Usa tipos específicos solo cuando sea NECESARIO
    var byte1 uint8 = 255       // Raro, usualmente para bytes
    var temperatura int16 = -40 // Raro
```

### Decimales (Floats)

```go
// float32 - 32 bits, menos preciso
var x float32 = 3.14159

// float64 - 64 bits, más preciso (DEFAULT)
var y float64 = 3.14159265358979

// Literal sin tipo es float64
pi := 3.14159265358979    // Tipo: float64
```

**Rango y precisión:**

```
float32:
 Rango: ~10^-38 a 10^38
 Precisión: ~7 dígitos decimales
 Usa cuando: Necesitas ocupar menos memoria

float64 (Recomendado):
 Rango: ~10^-308 a 10^308
 Precisión: ~15 dígitos decimales
 Usa cuando: (casi siempre)
```

**Advertencia: Comparaciones**

```go
x := 0.1 + 0.2
y := 0.3

if x == y {
    fmt.Println("Iguales")
} else {
    fmt.Println("Diferentes")  // ← Esto se ejecuta (bug flotante)
}

// Razón: Representación binaria de 0.1, 0.2, 0.3 no es exacta
```

### Booleanos

```go
// Solo dos valores: true y false
var activo bool = true
var eliminado bool = false

// Cero valor
var flag bool              // false

// Operadores booleanos
x := true && false         // false (AND)
y := true || false         // true (OR)
z := !true                 // false (NOT)
```

### Números con Separadores

Go permite separadores en números (Go 1.13+):

```go
// Más legible
var millones int = 1_000_000
var PI float64 = 3.141_592_653_589_793

// Hexadecimal
var hex int = 0xFF_FF_FF_FF

// Binario
var bits int = 0b1010_1100

// Output es igual
fmt.Println(1_000_000)     // 1000000
```

---

## 4.5 Constantes

### ¿Qué es una Constante?

Una constante es un valor que **no cambia**:

```go
const pi = 3.14159
const velocidadLuz = 299_792_458  // m/s

pi = 3.14  // ❌ ERROR: no puedes modificar constante
```

### Declaración

**Constante simple:**

```go
const nombre string = "Juan"
```

**Constantes múltiples:**

```go
const (
    pi     = 3.14159
    e      = 2.71828
    golden = 1.61803
)
```

**Sin tipo (untyped):**

```go
const x = 5        // Tipo se infiere cuando se usa
const y = 3.14     // Tipo se infiere cuando se usa
const z = "hola"   // Tipo se infiere cuando se usa
```

### Diferencia: Variable vs Constante

```

 Aspecto            │ Variable         │ Constante        │
clear
 Declaración        │ var x = 5        │ const x = 5      │
 Cambiar valor      │ x = 10 (✅ OK)   │ x = 10 (❌ ERROR) │
 Tipo               │ Tipado           │ Puede ser untyped│
 Inicialización     │ En runtime       │ En compile time  │
 Overhead           │ Ubicación RAM    │ Sin ubicación    │
 Performance        │ Acceso a memoria │ Inlined (rápido) │

```

### Constantes Untyped

```go
const x = 5          // untyped integer

// Puede usarse como cualquier tipo entero
var a int = x        // ✅ OK
var b int64 = x      // ✅ OK
var c int8 = x       // ✅ OK

// Pero no puedes cambiar tipo en runtime
var d float64 = x    // ✅ OK (compile time conversion)
```

### Iota - Enumeraciones

`iota` es una constante especial que incrementa automáticamente:

```go
const (
    Rojo int = iota    // 0
    Verde             // 1 (iota++ automático)
    Azul              // 2
    Amarillo          // 3
)

fmt.Println(Rojo)      // 0
fmt.Println(Azul)      // 2
```

**Casos de uso:**

```go
// Ejemplo: Estados de usuario
const (
    Inactivo int = iota  // 0
    Activo                // 1
    Bloqueado             // 2
    Suspendido            // 3
)

// Ejemplo: Permisos (flags)
const (
    Lectura uint = 1 << iota   // 1 (2^0)
    Escritura                   // 2 (2^1)
    Ejecucion                   // 4 (2^2)
    Borrar                      // 8 (2^3)
)

permisos := Lectura | Escritura  // 3 (puede leer y escribir)
```

---

## 4.6 Conversiones de Tipos

### Conversión Explícita

Go **requiere conversión explícita** entre tipos diferentes:

```go
var x int = 5
var y float64 = float64(x)      // Conversión explícita

var z int = int(y)              // De float a int

fmt.Println(x, y, z)            // 5 5 5
```

**Sintaxis general:**

```go
TipoDestino(valor)

// Ejemplos:
int32(x)
float64(x)
string(x)
bool(x)    // ❌ NO FUNCIONA: bool no tiene conversión
```

### Conversión de Números

**int ↔ float:**

```go
x := 10
f := float64(x)         // int → float64: 10.0

y := 3.14
i := int(y)             // float64 → int: 3 (se trunca)
```

**Diferentes tamaños de int:**

```go
var a int8 = 127
var b int16 = int16(a)  // int8 → int16: 127

var c int32 = 1000
var d int8 = int8(c)    // int32 → int8: -24 (overflow)
```

### Conversión a String

```go
x := 65
s := string(x)          // int → string: "A" (interpreta como rune)

f := 3.14
s2 := string(f)         // ❌ ERROR: no puedes convertir float a string

// Correcto: usar strconv
import "strconv"
s2 := strconv.FormatFloat(f, 'f', 2, 64)  // "3.14"
```

### Conversión de String

```go
s := "hello"

// string → []byte
b := []byte(s)          // [104 101 108 108 111]

// string → []rune
r := []rune(s)          // [104 101 108 108 111]

// string → int (requiere strconv)
i, err := strconv.Atoi("42")  // 42, nil
```

---

## 4.7 Strings y Runes

### ¿Qué es un String?

Un string en Go es una **secuencia inmutable de bytes**:

```go
s := "Hola"

// Internamente
// s = [72 111 108 97]  (bytes en ASCII/UTF-8)
// s = 'H' 'o' 'l' 'a'
```

**Propiedad: Inmutable**

```go
s := "Hola"
s[0] = 'J'      // ❌ ERROR: no puedes modificar strings

// Correcto: crear nuevo string
s = "Jola"      // ✅ OK: crear variable nueva
```

### Strings Literales

**Strings normales (interpretados):**

```go
s := "Hola\nMundo"      // Interpreta \n como newline

// Output:
// Hola
// Mundo
```

**Raw strings (literales):**

```go
s := `Hola\nMundo`      // \n es literal, no newline

// Output:
// Hola\nMundo
```

### Acceder a Caracteres

```go
s := "Hola"

// Por índice (retorna byte)
fmt.Println(s[0])       // 72 (valor ASCII de 'H')

// Convertir a char
fmt.Println(string(s[0]))   // "H"

// Rango de caracteres
fmt.Println(s[0:2])     // "Ho"
```

### ¿Qué es un Rune?

Un **rune** es un carácter Unicode (puede ser múltiples bytes):

```go
// String ASCII (1 byte por carácter)
s1 := "Hola"           // 4 bytes, 4 caracteres

// String con Unicode (múltiples bytes por carácter)
s2 := "Hola 世界"       // 12 bytes, 6 caracteres (世界 = 6 bytes)

// Acceder por bytes:
fmt.Println(len(s1))   // 4
fmt.Println(len(s2))   // 12 ← Cuenta BYTES, no caracteres

// Acceder por runes:
fmt.Println(len([]rune(s1)))   // 4
fmt.Println(len([]rune(s2)))   // 6 ← Cuenta CARACTERES
```

### Iterar String

**Por bytes (ASCII):**

```go
s := "Hola"

for i := 0; i < len(s); i++ {
    fmt.Printf("%c ", s[i])   // H o l a
}
```

**Por runes (Unicode correcto):**

```go
s := "Hola 世界"

for i, r := range s {
    fmt.Printf("Índice: %d, Rune: %c\n", i, r)
}

// Output:
// Índice: 0, Rune: H
// Índice: 1, Rune: o
// Índice: 2, Rune: l
// Índice: 3, Rune: a
// Índice: 4, Rune: (espacio)
// Índice: 5, Rune: 世
// Índice: 8, Rune: 界
```

---

## 4.8 Bytes y Runes Profundo

### Byte vs Rune

```

 Concepto │ Byte                │ Rune             │
clear
 Definición│ 8 bits (0-255)      │ 32 bits Unicode  │
 Alias    │ uint8               │ int32            │
 Uso      │ Datos binarios      │ Texto Unicode    │
 Ejemplo  │ 'A' = 65            │ '世' = 19990     │
 Rango    │ 0-255               │ 0-1,114,111      │

```

### Conversiones

```go
// byte ↔ string
b := byte(65)           // 65
s := string(b)          // "A"

// rune ↔ string
r := rune(65)           // 65
s := string(r)          // "A"

// Multiple runes → string
r1, r2, r3 := 'H', 'o', 'l'
s := string([]rune{r1, r2, r3})  // "Hol"
```

### UTF-8 Detallado

Go usa **UTF-8** internamente:

```
UTF-8 es variable-length encoding:

ASCII (1 byte):
    'A' = 0x41 = 01000001

Acentuados (2 bytes):
    'é' = 0xC3 0xA9

Asiáticos (3 bytes):
    '中' = 0xE4 0xB8 0xAD
    '世' = 0xE4 0xB8 0x96
    '界' = 0xE7 0x95 0x8C

Emojis (4 bytes):
    '😀' = 0xF0 0x9F 0x98 0x80
```

**Implicación:**

```go
s := "Hola😀"

// len(s) cuenta BYTES
fmt.Println(len(s))                 // 9 (H=1, o=1, l=1, a=1, 😀=4)

// len([]rune(s)) cuenta CARACTERES
fmt.Println(len([]rune(s)))         // 5
```

---

## 4.9 Paquete strconv - Conversiones

### Conversiones String ↔ Números

**String → int:**

```go
import "strconv"

s := "42"
i, err := strconv.Atoi(s)    // "ASCII to integer"
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println(i)           // 42 (type int)
}

// Versión genérica (permite especificar base)
i64, _ := strconv.ParseInt(s, 10, 64)  // base 10, 64 bits
i := int(i64)
```

**String → float:**

```go
s := "3.14"
f, _ := strconv.ParseFloat(s, 64)  // 3.14 (float64)

s2 := "2.5e2"
f2, _ := strconv.ParseFloat(s2, 64)  // 250.0
```

**String → bool:**

```go
b, _ := strconv.ParseBool("true")      // true
b2, _ := strconv.ParseBool("1")        // true
b3, _ := strconv.ParseBool("false")    // false
b4, _ := strconv.ParseBool("0")        // false
```

**int/float → String:**

```go
s := strconv.Itoa(42)               // "42"
s2 := strconv.FormatInt(42, 10)     // "42"
s3 := strconv.FormatFloat(3.14, 'f', 2, 64)  // "3.14"
```

### Manejo de Errores

```go
i, err := strconv.Atoi("abc")
if err != nil {
    fmt.Println("No es número válido")  // ← Esto se ejecuta
    fmt.Println(i)                      // 0 (valor zero)
}
```

---

## 4.10 Buenas Prácticas de Variables

### Nombrado (Naming)

**Regla 1: Nombres descriptivos**

```go
// ❌ Malo
var x int = 5
var y string = "Juan"

// ✅ Bueno
var edad int = 5
var nombre string = "Juan"
```

**Regla 2: Nombres cortos en contextos claros**

```go
// En bucle corto: ok usar i, j
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// Variable temporal que dura 2 líneas: ok
err := process()
if err != nil {
    return err
}
```

**Regla 3: MAYÚSCULAS para exportados**

```go
// Privado (dentro del package)
var contador int

// Público (otros packages pueden usar)
var Contador int
```

### Inicialización

**Preferencia: Short declaration**

```go
// ✅ Recomendado
nombre := "Juan"
edad := 25

// ❌ Verbose
var nombre string = "Juan"
var edad int = 25
```

**Excepto: Variables globales**

```go
// Global (fuera de funciones)
var globalCounter int = 0   // ← Debe ser 'var'

func main() {
    localCounter := 0       // ← Dentro puede ser short declaration
}
```

### Visibilidad Mínima

```go
// ❌ Exponer todo al package (malo)
var Global string = "valor"
var Another int = 5

// ✅ Solo lo que necesites exportado
var ConfigValue string = "privado"  // privado
var Config string = "publico"       // público
```

### Constantes sobre Variables

Usa constantes cuando el valor NO CAMBIA:

```go
// ❌ Variable innecesaria
var pi float64 = 3.14159

//  Constante
const pi = 3.14159
```

### Agrupar Relacionadas

```go
// ❌ Desordenado
var nombre string
var edad int
var ciudad string
var país string
var trabajando bool

// ✅ Agrupado
var (
    nombre string
    edad   int
    ciudad string
    país   string
    trabajando bool
)
```

---

## Ejercicios del Capítulo 4

### Ejercicio 1: Convertir Unidades

Crea programa que:
1. Pida distancia en metros
2. Convierta y muestre en: km, cm, mm
3. Use constantes para factores de conversión

### Ejercicio 2: Procesador de Texto

Crea programa que:
1. Pida un texto
2. Muestre: longitud, mayúsculas, minúsculas, invertido
3. Maneja Unicode (emojis, caracteres especiales)

### Ejercicio 3: Conversor de Bases Numéricas

Crea programa que:
1. Pida número decimal
2. Convierta a: binario, octal, hexadecimal
3. Usa fmt.Printf con diferentes verbos (%d, %b, %o, %x)

### Ejercicio 4: Validador de Tipos

Crea programa que:
1. Pida una cadena
2. Intente convertir a int, float, bool
3. Muestre qué conversión funcionó
4. Muestre errores si aplicable

---

**Fin del Capítulo 4**

---

