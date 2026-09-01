# Capítulo 10: Strings, runes y bytes

## Índice del Capítulo 10

1. [10.1 Las Tres Representaciones de Texto](#101-las-tres-representaciones-de-texto)
2. [10.2 Strings - Secuencias Inmutables](#102-strings---secuencias-inmutables)
3. [10.3 UTF-8 Encoding - El Estándar de Go](#103-utf-8-encoding---el-estándar-de-go)
4. [10.4 Runes - Caracteres Unicode](#104-runes---caracteres-unicode)
5. [10.5 Bytes - Datos Binarios](#105-bytes---datos-binarios)
6. [10.6 Conversiones entre String, Rune, Byte](#106-conversiones-entre-string-rune-byte)
7. [10.7 Paquete strings - Operaciones](#107-paquete-strings---operaciones)
8. [10.8 Paquete strconv - Conversiones](#108-paquete-strconv---conversiones)
9. [10.9 Expresiones Regulares](#109-expresiones-regulares)
10. [10.10 Buenas Prácticas con Texto](#1010-buenas-prácticas-con-texto)

---

## 10.1 Las Tres Representaciones de Texto

### Conceptos Fundamentales

Go tiene **tres formas diferentes** de representar texto:

```

 Tipo         │ Definición          │ Uso              │
clear
 string       │ Secuencia de BYTES  │ Texto (UTF-8)    │
 rune         │ Carácter UNICODE    │ Caracteres       │
 byte         │ Byte individual     │ Datos binarios   │
 []byte       │ Slice de bytes      │ Buffer binario   │

```

### Comparación Visual

```
Texto: "Hola 世"

String (internamente bytes en UTF-8):
 H    = 0x48     (1 byte)
 o    = 0x6F     (1 byte)
 l    = 0x6C     (1 byte)
 a    = 0x61     (1 byte)
 (space) = 0x20  (1 byte)
 世   = 0xE4 0xB8 0x96  (3 bytes)

Total: 9 bytes

Runes (caracteres Unicode):
 'H'   = U+0048
 'o'   = U+006F
 'l'   = U+006C
 'a'   = U+0061
 ' '   = U+0020
 '世'  = U+4E16

Total: 6 runes

Bytes (individuales):
[0x48, 0x6F, 0x6C, 0x61, 0x20, 0xE4, 0xB8, 0x96]
Total: 9 bytes
```

---

## 10.2 Strings - Secuencias Inmutables

### ¿Qué es un String?

Un string es una **secuencia inmutable de bytes** encodificados en UTF-8:

```go
s := "Hola"
fmt.Println(len(s))         // 4 (BYTES, no caracteres)
```

### Propiedades de Strings

**Inmutabilidad:**

```go
s := "Hola"
s[0] = 'J'              // ❌ ERROR: cannot assign to s[0]

// Correcto: crear nuevo string
s = "Jola"              // ✅ OK
```

**Concatenación:**

```go
s1 := "Hola"
s2 := "Mundo"
s3 := s1 + " " + s2     // "Hola Mundo"
```

**Indexación (retorna byte):**

```go
s := "Hola"
fmt.Println(s[0])       // 72 (ASCII value de 'H')
fmt.Println(s[0] == 'H')   // false (72 ≠ 'H')
fmt.Println(s[0] == byte('H'))  // true
```

### Literals de String

**String normal (interpretado):**

```go
s := "Hola\nMundo"     // \n se interpreta como newline
// Output:
// Hola
// Mundo
```

**Raw string (literal):**

```go
s := `Hola\nMundo`     // \n es LITERAL, no newline
// Output: Hola\nMundo
```

**Caracteres especiales:**

```go
s := "Línea 1\nLínea 2"      // Newline
s := "Tabulación\tAquí"       // Tab
s := "Comilla: \"Hola\""      // Escape comillas
s := "Barra invertida: \\"    // Escape barra
```

### Comparación de Strings

```go
s1 := "apple"
s2 := "apple"
s3 := "banana"

s1 == s2                // true
s1 == s3                // false
s1 < s3                 // true (lexicográfico: a < b)
s1 <= s2                // true
```

---

## 10.3 UTF-8 Encoding - El Estándar de Go

### ¿Por Qué UTF-8?

```
ASCII: Solo 127 caracteres (insuficiente para mundo)
Latin-1: 256 caracteres (aún insuficiente)
UTF-16: 2-4 bytes por carácter (ineficiente para ASCII)

UTF-8 (GOLD STANDARD):
 Compatible con ASCII (ASCII = UTF-8 para 0-127)
 1-4 bytes por carácter
 Autosincronizado (puedes buscar bytes válidos)
 Eficiente: caracteres comunes = 1-2 bytes
 Go usa UTF-8 por defecto
```

### Codificación UTF-8

```
U+0000 - U+007F (ASCII):
    1 byte:   0xxxxxxx
    Ejemplo: 'A' = 0x41

U+0080 - U+07FF:
    2 bytes:  110xxxxx 10xxxxxx
    Ejemplo: 'é' = 0xC3 0xA9

U+0800 - U+FFFF:
    3 bytes:  1110xxxx 10xxxxxx 10xxxxxx
    Ejemplo: '中' = 0xE4 0xB8 0xAD

U+10000 - U+10FFFF:
    4 bytes:  11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
    Ejemplo: '😀' = 0xF0 0x9F 0x98 0x80
```

### Implicación Práctica

```go
s := "Hola 中"

len(s)                  // 8 (bytes: H=1, o=1, l=1, a=1, sp=1, 中=3)
len([]rune(s))          // 6 (caracteres: 5 ASCII + 1 chino)

// ❌ Malo: iterar por índices
for i := 0; i < len(s); i++ {
    fmt.Printf("%c\n", s[i])  // Muestra bytes, no caracteres
}

// ✅ Bueno: iterar por runes
for _, r := range s {
    fmt.Printf("%c\n", r)     // Muestra caracteres correctos
}
```

---

## 10.4 Runes - Caracteres Unicode

### ¿Qué es un Rune?

```go
type rune = int32  // Alias para int32

// Rune representa un carácter Unicode
r := 'A'           // Rune literal
r2 := '中'         // Rune Unicode
r3 := '😀'         // Rune Emoji
```

### Conversión String → Runes

```go
s := "Hola 中国"

// Convertir a slice de runes
runes := []rune(s)

for i, r := range runes {
    fmt.Printf("Índice %d: %c (U+%04X)\n", i, r, r)
}

// Output:
// Índice 0: H (U+0048)
// Índice 1: o (U+006F)
// Índice 2: l (U+006C)
// Índice 3: a (U+0061)
// Índice 4:   (U+0020)
// Índice 5: 中 (U+4E2D)
// Índice 6: 国 (U+56FD)
```

### Iterar por Runes (Range)

```go
s := "Hola 中"

// range en string retorna (índice_byte, rune)
for i, r := range s {
    fmt.Printf("Byte %d: %c (tipo: %T)\n", i, r, r)
}

// Output:
// Byte 0: H (tipo: int32)
// Byte 1: o (tipo: int32)
// Byte 2: l (tipo: int32)
// Byte 3: a (tipo: int32)
// Byte 4:   (tipo: int32)
// Byte 5: 中 (tipo: int32)  ← Índice salta de 4 a 5 (aunque 中 es 3 bytes)
```

### Longitud Real (Caracteres)

```go
s := "Hello 世界😀"

fmt.Println(len(s))                 // 14 (bytes)
fmt.Println(len([]rune(s)))         // 8 (caracteres)
fmt.Println(utf8.RuneCountInString(s))  // 8 (caracteres, más eficiente)
```

---

## 10.5 Bytes - Datos Binarios

### ¿Qué es un Byte?

```go
type byte = uint8  // Alias para uint8

// Byte individual
b := byte(65)       // 65
fmt.Println(b)      // 65
fmt.Println(string(b))  // "A"
```

### Slice de Bytes

```go
// Slice de bytes
datos := []byte{72, 101, 108, 108, 111}  // "Hello"

// Convertir a string
s := string(datos)
fmt.Println(s)      // Hello

// Convertir string a bytes
s2 := "Hola"
b2 := []byte(s2)
fmt.Println(b2)     // [72 111 108 97]
```

### Modificar Bytes (NO Strings)

```go
// Strings son inmutables
s := "Hello"
s[0] = 'J'          // ❌ ERROR

// Bytes son mutables
b := []byte("Hello")
b[0] = 'J'          // ✅ OK
fmt.Println(string(b))  // Jello

// Esto es común: convertir a bytes, modificar, convertir a string
```

### Caso Real: Procesar Datos Binarios

```go
// Leer datos binarios de archivo
data := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}

// Procesar bytes
for i, b := range data {
    if b >= 0x41 && b <= 0x5A {
        fmt.Printf("Byte %d es mayúscula: %c\n", i, b)
    }
}
```

---

## 10.6 Conversiones entre String, Rune, Byte

### String ↔ []byte

```go
// String → []byte
s := "Hola"
b := []byte(s)          // [72 111 108 97]

// []byte → String
s2 := string(b)         // "Hola"
```

### String ↔ []rune

```go
// String → []rune
s := "Hola 中"
r := []rune(s)          // ['H', 'o', 'l', 'a', ' ', '中']

// []rune → String
s2 := string(r)         // "Hola 中"
```

### Byte vs Rune Conversión

```go
// Byte → Rune (extender)
b := byte(65)           // 65 (0x41)
r := rune(b)            // 65 (U+0041 'A')

// Rune → Byte (truncar si > 255)
r2 := rune(200)         // U+00C8
b2 := byte(r2)          // 200 (OK porque < 256)

r3 := rune(300)         // U+012C
b3 := byte(r3)          // 44 (truncado!)
```

### Comparación: Rendimiento

```
[]byte(s):   O(n) - copia datos
s[:]:        O(1) - no copia, vista

string(b):   O(n) - copia datos
b:           O(1) - ya está ahí

[]rune(s):   O(n) - decodifica UTF-8
s:           O(n) - acceder individual
```

---

## 10.7 Paquete strings - Operaciones

### Búsqueda y Contiene

```go
import "strings"

s := "Hola Mundo"

strings.Contains(s, "Mundo")        // true
strings.ContainsAny(s, "abc")       // true (contiene 'a')
strings.Index(s, "Mundo")           // 5
strings.LastIndex(s, "o")           // 8
strings.Count(s, "o")               // 2
```

### Transformación

```go
strings.ToUpper("hola")             // "HOLA"
strings.ToLower("HOLA")             // "hola"
strings.Title("hola mundo")         // "Hola Mundo" (deprecated, usar cases)
strings.Trim("  hola  ", " ")       // "hola"
strings.TrimPrefix("prefijo_hola", "prefijo_")  // "hola"
strings.TrimSuffix("hola_sufijo", "_sufijo")    // "hola"
```

### Split y Join

```go
// Dividir
partes := strings.Split("a,b,c", ",")   // ["a", "b", "c"]

// Unir
unidos := strings.Join([]string{"a", "b", "c"}, ",")  // "a,b,c"

// Split Fields (por espacios)
campos := strings.Fields("hola  mundo  test")  // ["hola", "mundo", "test"]
```

### Replace

```go
s := "hola hola mundo"

strings.Replace(s, "hola", "adiós", 1)      // "adiós hola mundo" (1 reemplazo)
strings.Replace(s, "hola", "adiós", -1)     // "adiós adiós mundo" (todos)
strings.ReplaceAll(s, "hola", "adiós")      // "adiós adiós mundo"
```

### Comparación

```go
strings.Compare("abc", "abd")       // -1 (abc < abd)
strings.Compare("abc", "abc")       // 0 (iguales)
strings.Compare("abd", "abc")       // 1 (abd > abc)

strings.EqualFold("HOLA", "hola")   // true (case-insensitive)
```

### Has Prefix/Suffix

```go
strings.HasPrefix("archivo.txt", "archivo")    // true
strings.HasSuffix("archivo.txt", ".txt")       // true
```

---

## 10.8 Paquete strconv - Conversiones

### String ↔ Números

```go
import "strconv"

// String → int
i, _ := strconv.Atoi("42")              // 42 (int)
i64, _ := strconv.ParseInt("42", 10, 64)    // 42 (int64)

// int → String
s := strconv.Itoa(42)                   // "42"
s64 := strconv.FormatInt(42, 10)        // "42"

// String → float
f, _ := strconv.ParseFloat("3.14", 64)  // 3.14 (float64)

// float → String
s := strconv.FormatFloat(3.14, 'f', 2, 64)  // "3.14"
```

### String ↔ Bool

```go
b, _ := strconv.ParseBool("true")       // true
b2, _ := strconv.ParseBool("1")         // true
b3, _ := strconv.ParseBool("false")     // false
b4, _ := strconv.ParseBool("0")         // false

// Bool → String
s := strconv.FormatBool(true)           // "true"
```

### Manejo de Errores

```go
i, err := strconv.Atoi("no es número")
if err != nil {
    fmt.Println("Error:", err)          // Error: strconv.Atoi: parsing "...": invalid syntax
}
```

---

## 10.9 Expresiones Regulares

### Introducción

```go
import "regexp"

// Compilar patrón
patron := regexp.MustCompile(`\d+`)     // Números

// Buscar
if patron.MatchString("abc123def") {
    fmt.Println("Contiene números")
}
```

### Operaciones Comunes

**MatchString - ¿Coincide el patrón?**

```go
patron := regexp.MustCompile(`^[a-z]+$`)  // Solo letras minúsculas

patron.MatchString("abc")       // true
patron.MatchString("abc123")    // false
```

**FindString - Primera coincidencia**

```go
patron := regexp.MustCompile(`\d+`)     // Números

patron.FindString("abc123def456")   // "123"
patron.FindString("abc")            // "" (no encontrado)
```

**FindAllString - Todas las coincidencias**

```go
patron := regexp.MustCompile(`\d+`)

matches := patron.FindAllString("abc123def456", -1)  // ["123", "456"]
```

**Replace - Reemplazar coincidencias**

```go
patron := regexp.MustCompile(`\d+`)
resultado := patron.ReplaceAllString("abc123def456", "X")  // "abcXdefX"
```

### Patrones Comunes

```go
// Email simple
regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// URL
regexp.MustCompile(`https?://[^\s]+`)

// Teléfono (formato 123-456-7890)
regexp.MustCompile(`\d{3}-\d{3}-\d{4}`)

// Hexadecimal
regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)

// IPv4
regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
```

---

## 10.10 Buenas Prácticas con Texto

### Iterar Correctamente

```go
// ❌ Malo: itera por bytes
s := "Hola 中"
for i := 0; i < len(s); i++ {
    fmt.Println(s[i])   // Muestra bytes, confuso
}

// ✅ Bueno: itera por runes
for _, r := range s {
    fmt.Println(string(r))  // Muestra caracteres
}
```

### Contar Caracteres

```go
// ❌ Malo
s := "Hola 中"
count := len(s)         // 8 (bytes, no caracteres)

// ✅ Bueno
count := len([]rune(s))  // 6 (caracteres reales)

// O más eficiente (Go 1.11+)
import "unicode/utf8"
count := utf8.RuneCountInString(s)  // 6
```

### Manipulación de Strings

```go
// Strings son inmutables, así que:
s := "Hello"

// ❌ Ineficiente: múltiples concatenaciones
for i := 0; i < 1000; i++ {
    s = s + "X"         // Crea nuevo string cada vez
}

// ✅ Eficiente: usar strings.Builder
import "strings"

var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString("X")
}
resultado := sb.String()
```

### Validación de Entrada

```go
import "regexp"

func esEmail(s string) bool {
    patron := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return patron.MatchString(s)
}

esEmail("usuario@example.com")  // true
esEmail("invalido@")            // false
```

### Normalizacinnn

```go
import "strings"
import "unicode"

// Minúsculas
s := strings.ToLower("HOLA")    // "hola"

// Remover espacios
s := strings.TrimSpace("  hola  ")  // "hola"

// Remover caracteres especiales
s := strings.Map(func(r rune) rune {
    if unicode.IsLetter(r) || unicode.IsDigit(r) {
        return r
    }
    return -1  // Remover
}, "Hol@-a")  // "Hola"
```

### Case-Insensitive

```go
import "strings"

s1 := "HOLA"
s2 := "hola"

if strings.EqualFold(s1, s2) {
    fmt.Println("Iguales (ignorando mayúsculas)")
}
```

### Slice de String

```go
s := "Hola Mundo"

// ❌ Peligroso si tiene UTF-8 no-ASCII
sub := s[0:4]       // "Hola" (OK porque son ASCII)

// Con caracteres multi-byte
s2 := "Hola 中国"
sub2 := s2[0:7]     // Puede cortarse un carácter UTF-8

// ✅ Mejor: convertir a runes
runes := []rune(s2)
sub3 := string(runes[0:4])  // "Hola" (correcto)
```

---

## Ejercicios del Capítulo 10

### Ejercicio 1: Analizador de Texto

Crea programa que:

1. Reciba texto
2. Cuente palabras, caracteres (real), bytes
3. Encuentre palabra más larga
4. Calcule promedio de longitud de palabra
5. Muestre histograma de frecuencia de caracteres

### Ejercicio 2: Formateador de Strings

Crea programa que:

1. Implemente función `center(s string, width int) string`
2. Implemente `pad(s string, char rune, width int) string`
3. Implemente `repeat(s string, times int) string`
4. Implemente `reverse(s string) string`
5. Implemente `removeVowels(s string) string`

### Ejercicio 3: Validador con Regex

Crea programa que valide:

1. Email válido
2. Teléfono formato 123-456-7890
3. Contraseña: 8+ caracteres, al menos 1 mayúscula, 1 minúscula, 1 número
4. URL vlida
5. Código hexadecimal (#RRGGBB)

### Ejercicio 4: Codificador/Decodificador

Crea programa que:

1. Implemente ROT13 (desplazar letras)
2. Implemente cifrado de sustitución simple
3. Implemente Base64 encoding/decoding
4. Manaje caracteres especiales y espacios
5. Valide entrada y salida

### Ejercicio 5: Procesador de Archivos de Texto

Crea programa que:

1. Lea línea de archivo
2. Cuente líneas, palabras, caracteres
3. Encuentre línea más larga
4. Reporte palabras únicas y su frecuencia
5. Busque patrón regex en archivo

---

**Fin del Capítulo 10**

---

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/10-strings-runes-y-bytes/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/10-strings-runes-y-bytes):

```bash
cd examples/10-strings-runes-y-bytes
go run .
```
