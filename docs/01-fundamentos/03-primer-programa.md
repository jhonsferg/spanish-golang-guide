# Capítulo 3: Tu primer programa - De Hello World a los conceptos fundamentales

## Índice del Capítulo 3

1. [3.1 La Anatomía de un Programa Go](#31-la-anatomía-de-un-programa-go)
2. [3.2 El Programa Más Simple: Hello World](#32-el-programa-más-simple-hello-world)
3. [3.3 ¿Qué Sucede Cuando Ejecutas tu Programa?](#33-qué-sucede-cuando-ejecutas-tu-programa)
4. [3.4 La Function main()](#34-la-function-main)
5. [3.5 El Package main](#35-el-package-main)
6. [3.6 Imports y la Standard Library](#36-imports-y-la-standard-library)
7. [3.7 Compilar vs Ejecutar](#37-compilar-vs-ejecutar)
8. [3.8 Variables de Entorno en Ejecución](#38-variables-de-entorno-en-ejecución)
9. [3.9 Argumentos de Línea de Comando](#39-argumentos-de-línea-de-comando)
10. [3.10 Tu Primer Programa Real: Aplicación Interactiva](#310-tu-primer-programa-real-aplicación-interactiva)

---

## 3.1 La Anatomía de un Programa Go

### Estructura Mínima de Go

Todo programa Go tiene esta estructura:

```
archivo.go
 package declaration
 imports (opcional)
 constants (opcional)
 variables (opcional)
 types (opcional)
 functions
 main function (si es ejecutable)
```

### Ejemplo Anatómico Completo

```go
// 1. Package declaration (OBLIGATORIO)
package main

// 2. Imports (OPCIONAL)
import (
    "fmt"
    "os"
)

// 3. Constants (OPCIONAL)
const Version = "1.0.0"

// 4. Variables globales (NO RECOMENDADO)
var GlobalCounter int

// 5. Types (OPCIONAL)
type Usuario struct {
    Nombre string
    Edad   int
}

// 6. Funciones helper (OPCIONAL)
func saludar(nombre string) string {
    return fmt.Sprintf("Hola, %s!", nombre)
}

// 7. Main function (OBLIGATORIA en ejecutables)
func main() {
    mensaje := saludar("Juan")
    fmt.Println(mensaje)
}
```

**Orden de las cosas en Go:**

```
Go NO importa mucho el orden, EXCEPTO:

 package: SIEMPRE primera línea (si existe)
 imports: SIEMPRE despus de package
 El resto: Puedes ordenar como quieras
   ├─ Constants → Variables → Types → Functions
   ├─ O Functions → Constants → Variables
   └─ Importa legibilidad, no el lenguaje
```

### Reglas Críticas

**1. Un archivo = Un package**

```
Cada archivo .go pertenece a exactamente un package:

main.go:
    package main  ← ESTE archivo es parte del package main

helper.go:
    package main  ← ESTE archivo TAMBIÉN es parte del package main

NO es válido:
    main.go:    package main
    helper.go:  package otro  ← ¡ERROR! Todos deben ser mismo package
```

**2. Múltiples archivos, mismo package**

```
CORRECTO - Mismos package en múltiples archivos:

proyecto/
 main.go          (package main)
 helpers.go       (package main)
 database.go      (package main)
 go.mod

Cuando compilas:
    go build

Go automáticamente:
    ├─ Lee TODOS los archivos .go
    ├─ Los agrupa por package
    ├─ Los compila juntos
    └─ Crea 1 binario ejecutable
```

**3. Visibilidad de símbolos (Exportados)**

```go
// En Go, la visibilidad es por MAYÚSCULA:

// PRIVADO (accesible solo dentro del package)
func saludar() { }        ← Comienza con minúscula
var contador int          ← Privado
type usuario struct { }   ← Privado

// PÚBLICO (exportado, accesible desde otros packages)
func Saludar() { }        ← Comienza con MAYÚSCULA
var Contador int          ← Público
type Usuario struct { }   ← Público

Impacto:
 Si importas github.com/usuario/pkg
   ├─ Puedes usar: pkg.Saludar()
   ├─ NO puedes usar: pkg.saludar()
   └─ Si intentas: ERROR de compilación
```

---

## 3.2 El Programa Más Simple: Hello World

### El Clásico

```go
package main

import "fmt"

func main() {
    fmt.Println("Hola, Mundo!")
}
```

**Ejecutarlo:**

```bash
# Opción 1: Ejecutar directamente (compilación interna)
go run main.go
# Output: Hola, Mundo!

# Opción 2: Compilar primero, luego ejecutar
go build
./main           # En macOS/Linux
./main.exe       # En Windows

# Output: Hola, Mundo!
```

### ¿Por Qué Este Programa Funciona?

Análisis línea por línea:

```go
package main              // 1. Declara que este archivo es el package main
                          //    "main" es especial: denota un programa ejecutable

import "fmt"              // 2. Importa el package fmt (formatting)
                          //    fmt provee funciones para imprimir

func main() {             // 3. Define función main()
                          //    DEBE existir exactamente una en cada ejecutable
                          //    DEBE ser en package main
                          //    NO recibe argumentos
                          //    NO retorna valores

    fmt.Println(...)      // 4. Llama función Println del package fmt
                          //    Imprime texto + newline
}
```

### Variaciones Comunes

**Múltiples salidas:**

```go
package main

import "fmt"

func main() {
    fmt.Println("Línea 1")
    fmt.Println("Línea 2")
    fmt.Println("Línea 3")
}

// Output:
// Línea 1
// Línea 2
// Línea 3
```

**Print vs Println vs Printf:**

```go
package main

import "fmt"

func main() {
    // Print: sin newline
    fmt.Print("A")
    fmt.Print("B")
    fmt.Print("C")
    // Output: ABC

    fmt.Println()  // Salto de línea

    // Println: con newline
    fmt.Println("Hola")
    fmt.Println("Mundo")
    // Output:
    // Hola
    // Mundo

    fmt.Println()  // Salto de línea

    // Printf: formateado (C-style)
    nombre := "Juan"
    edad := 30
    fmt.Printf("Nombre: %s, Edad: %d\n", nombre, edad)
    // Output: Nombre: Juan, Edad: 30
}
```

---

## 3.3 ¿Qué Sucede Cuando Ejecutas tu Programa?

### Con `go run` (Ejecución Directa)

```
$ go run main.go


 Paso 1: PARSE (Análisis Sintáctico)                 │
 └─ Lee main.go                                      │
 └─ Construye Abstract Syntax Tree (AST)             │
 └─ Si hay errores de sintaxis: FALLA aquí           │

 Paso 2: TYPE CHECKING (Verificación de Tipos)      │
 └─ Verifica que tipos sean válidos                  │
 └─ fmt.Println existe y acepta string? ✓            │
 Si hay error de tipos: FALLA aquí                │

 Paso 3: COMPILATION (Compilación)                   │
 └─ Convierte Go a código máquina intermedio (SSA)   │
 └─ Aplica optimizaciones                            │
 └─ Genera assembly para tu plataforma               │

 Paso 4: LINKING (Enlazamiento)                      │
 └─ Junta archivos objeto                            │
 └─ Incrustra runtime de Go                          │
 └─ Crea binario temporal en memoria                 │

 Paso 5: EXECUTION (Ejecución)                       │
 └─ Ejecuta binario temporal                         │
 └─ Runtime inicializa (scheduler, GC, etc)          │
 └─ Llama main()                                     │
 └─ Imprime: "Hola, Mundo!"                          │
 └─ main() retorna                                   │
 └─ Runtime cleanup                                  │
 └─ Programa termina                                 │


TOTAL: ~1 SEGUNDO (primer run, sin caché)
```

### Con `go build` + Ejecución Manual

```
$ go build
$ ./main


 go build:                                            │
 └─ Pasos 1-4: PARSE → TYPE CHECK → COMPILE → LINK   │
 └─ Crea binario PERMANENTE: ./main                  │
 └─ Guarda en disco                                  │
                                                      │
 Luego ./main:                                       │
 └─ SO carga binario desde disco a memoria           │
 └─ Ejecuta (sin compilación)                        │
 └─ Runtime inicializa                               │
 └─ Llama main()                                     │
 └─ Termina                                          │
                                                      │
 VENTAJA: Si ejecutas ./main 100 veces, no           │
          recompila. Muy más rápido.                 │

```

### El Runtime de Go

Cuando tu programa comienza, ANTES de que `main()` sea llamada:

```go
// Esto sucede ANTES de main()

func init() {
    // Go ejecuta TODOS los init() de todos los packages
    // Orden: Dependencias primero

    // Tipicamente usado para:
    // ├─ Inicializar variables globales complejas
    // ├─ Registrar drivers (database/sql)
    // ├─ Validación de configuración
    // └─ Setup de paquetes
}

func main() {
    // AHORA sí, tu código corre
}
```

**Orden de inicialización completo:**

```
1. Program Start
   ├─ Load binario en memoria
   ├─ Runtime de Go se inicializa
   │   ├─ Memory allocator
   │   ├─ Garbage collector
   │   ├─ Goroutine scheduler
   │   └─ Signal handlers

2. Package Initialization (packages importados)
   ├─ Variables globales se declaran
   ├─ init() se ejecuta (en orden de import)
   └─ Esto sucede para TODOS los packages

3. main() Initialization
   ├─ init() de main se ejecuta (si existe)
   └─ main() se ejecuta

4. Execution
 Tu código Go   └

5. Termination
   ├─ main() retorna
   ├─ deferred functions se ejecutan (reverse order)
   ├─ OS cleanup
   └─ Process termina con exit code

Total: Típicamente <100ms de overhead
```

---

## 3.4 La Function main()

### Definición Exacta

```go
func main() {
    // Tu código aquí
}
```

**Reglas estrictas:**

```
 DEBE llamarse exactamente "main" (minúsculas)
 NO puede tener argumentos
 NO puede retornar valores
 DEBE estar en package main
 Debe existir exactamente UNA en todo el programa
 Es el entry point de ejecución
 Si retorna (o llega al final), programa termina
```

### Comparación con Otros Lenguajes

```
Go:          func main() { }
C:           int main() { }
Java:        public static void main(String[] args) { }
Python:      if __name__ == "__main__": pass
JavaScript:  // No tiene main, solo ejecuta

Diferencia: Go no permite argumentos a main()
           Para argumentos, usa os.Args
```

### Retornar Códigos de Salida

Go no tiene `return` en main, pero puedes usar `os.Exit`:

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Haciendo algo...")

    // Si todo OK
    fmt.Println("Éxito!")
    os.Exit(0)      // 0 = OK

    // Si hay error
    fmt.Println("Error ocurrió")
    os.Exit(1)      // 1 = Error genérico
}
```

**Códigos de salida estándar:**

```
0   - Éxito
1   - Error genérico
2   - Mal uso del programa
127 - Comando no encontrado
128 - Señal recibida (+ número de señal)
```

### Ejemplo: Programa que valida

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Uso: programa <nombre>")
        os.Exit(2)      // Mal uso
    }

    nombre := os.Args[1]

    if nombre == "" {
        fmt.Println("Error: nombre vacío")
        os.Exit(1)      // Error
    }

    fmt.Printf("Hola, %s!\n", nombre)
    os.Exit(0)          // OK
}
```

---

## 3.5 El Package main

### ¿Qué es un Package?

Package es una unidad de organización de código:

```go
// Un package es UNA carpeta con MÚLTIPLES archivos .go
// que comparten el MISMO nombre de package

proyecto/
 main.go          // package main
 helpers.go       // package main
 database.go      // package main
 go.mod

Cuando haces: go build
 Go lee TODOS los archivos .go
 Los agrupa por package (todos son "main" aquí)
 Los compila juntos
 Crea 1 binario
```

### Package main vs Otros Packages

**Package main (Ejecutable):**

```

 package main                        │

 Características:                    │
 ├─ DEBE contener func main()        │
 ├─ Se compila a BINARIO ejecutable  │
 ├─ Puedes ejecutar: ./binario       │
 └─ No puedes importar desde otro    │

```

**Otros packages (Librería):**

```go
package math

// NO puede tener func main()

func Add(a, b int) int {
    return a + b
}

// Se compila a LIBRERÍA (.a file)
// Otros packages pueden importar:
//     import "github.com/usuario/math"
//     result := math.Add(1, 2)
```

### Estructura de Carpetas

**Simple (1 package en raíz):**

```
proyecto/
 go.mod
 main.go          (package main)
 helpers.go       (package main)
 database.go      (package main)

Compilar: go build
Resultado: ./proyecto (binario)
```

**Compleja (múltiples packages):**

```
proyecto/
 go.mod                    (module github.com/usuario/proyecto)
 main.go                   (package main)

 internal/
   ├── database/
   │   ├── db.go            (package database)
   │   └── query.go         (package database)
   │
   └── utils/
       ├── string.go        (package utils)
       └── math.go          (package utils)

 pkg/
    └── api/
        ├── handler.go       (package api)
        └── server.go        (package api)

Import paths:
 github.com/usuario/proyecto/internal/database
 github.com/usuario/proyecto/internal/utils
 github.com/usuario/proyecto/pkg/api

En main.go:
    import (
        "github.com/usuario/proyecto/internal/database"
        "github.com/usuario/proyecto/pkg/api"
    )
```

---

## 3.6 Imports y la Standard Library

### ¿Qué es un Import?

Import permite usar código de otro package:

```go
import "fmt"

func main() {
    // Ahora puedes usar todo del package fmt
    fmt.Println("Hola")
}
```

### Tipos de Imports

**Import simple:**

```go
import "fmt"

// Uso: fmt.Println()
```

**Import múltiple:**

```go
import (
    "fmt"
    "os"
    "strings"
)

// Uso:
// fmt.Println()
// os.Exit()
// strings.ToUpper()
```

**Import con alias:**

```go
import (
    f "fmt"           // Alias f para fmt
    "os"
)

func main() {
    f.Println("Hola")  // Usando alias
    os.Exit(0)
}
```

**Import anónimo (side effects):**

```go
import (
    "database/sql"
    _ "github.com/lib/pq"  // Ejecuta init() de pq, pero no usa funciones
)

// Se usa típicamente para registrar drivers
```

### Standard Library Esencial

Go viene con una stdlib vasta. Las principales para principiantes:

**fmt - Formatting**

```go
import "fmt"

fmt.Println(...)        // Imprime con newline
fmt.Printf("%s", "x")   // Formato tipo C
fmt.Sprintf(...)        // Retorna string formateado
```

**os - Sistema Operativo**

```go
import "os"

os.Exit(0)              // Terminar programa
os.Args                 // Argumentos de línea de comando
os.Getenv("HOME")       // Variable de entorno
os.Getwd()              // Directorio actual
```

**strings - Strings**

```go
import "strings"

strings.ToUpper("hola")         // HOLA
strings.ToLower("HOLA")         // hola
strings.Split("a,b,c", ",")     // []string{"a", "b", "c"}
strings.Contains("hola", "ol")  // true
```

**io - Input/Output**

```go
import "io"

io.Copy(dst, src)       // Copiar bytes de src a dst
io.ReadAll(r)           // Leer todo de un reader
// Interfaces: Reader, Writer, Closer
```

**math - Matemáticas**

```go
import "math"

math.Sqrt(16)           // 4.0
math.Pi                 // 3.14159...
math.Max(1, 2)          // 2
math.Min(1, 2)          // 1
```

**time - Tiempo**

```go
import "time"

time.Now()                          // Hora actual
time.Sleep(1 * time.Second)         // Dormir
time.ParseDuration("5s")            // Duración
time.Date(2024, 1, 15, 10, 30, 0)  // Crear fecha
```

### Path Completo de Imports

Cuando escribes `import "fmt"`, Go busca en:

```
1. Go root: /usr/local/go/src/fmt
2. Go path: $GOPATH/src/fmt (legacy)
3. Go modules cache: $GOMODCACHE (típicamente ~/.go/pkg/mod)
4. External: github.com/usuario/package
```

---

## 3.7 Compilar vs Ejecutar

### go run: Todo en Uno

```bash
go run main.go

# Internamente:
# 1. Compila a binario temporal
# 2. Lo ejecuta
# 3. Borra el binario
```

**Ventajas:**

- Rápido para testing
- No deja archivos temporales
- Perfecto para scripts

**Desventajas:**

- Compila cada vez (lento si se ejecuta múltiples veces)
- No creas binario distribuible

### go build: Crear Ejecutable

```bash
go build

# Genera: ./main (en macOS/Linux) o ./main.exe (Windows)
```

**Personalizar nombre:**

```bash
go build -o mi-app
# Genera: ./mi-app

# Con ruta
go build -o bin/miapp
# Genera: ./bin/miapp
```

**Compilar para otra plataforma:**

```bash
# Linux 64-bit
GOOS=linux GOARCH=amd64 go build

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o app.exe

# macOS ARM (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build
```

### Comparación de Performance

```
                     │ go run │ go build (first) │ go build (cached) │
clear
 Compilación         │ Sí     │ Sí               │ Solo cambios      │
 Tiempo (cold start) │ 2-3s   │ 2-3s             │ 0.5s              │
 Ejecución           │ 0.1s   │ 0.1s             │ 0.1s              │
 Binario generado    │ No     │ Sí               │ Sí                │
 Ideal para          │ Dev    │ Distribución     │ Redistribución    │

```

### Compilación con Optimizaciones

```bash
# Build pequeño (remover debug info)
go build -ldflags="-s -w"

# Ejemplo:
# Sin flags: 5 MB
# Con flags: 2 MB

# Compilación rápida (sin optimizaciones)
go build -gcflags="all=-N"
```

---

## 3.8 Variables de Entorno en Ejecución

### Acceder a Variables de Entorno

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Obtener variable de entorno
    home := os.Getenv("HOME")      // Retorna string, vacío si no existe
    fmt.Println("HOME:", home)

    // Obtener variable O valor default
    user, ok := os.LookupEnv("USER")
    if !ok {
        user = "desconocido"
    }
    fmt.Println("USER:", user)
}
```

**Ejecutar con variables de entorno:**

```bash
# Pasar variable inline
HOME=/root go run main.go

# O configurar primero
export CUSTOM_VAR="valor"
go run main.go
```

### Variables de Entorno Comunes de Go

```bash
# Compilación
GOOS=linux              # Sistema operativo target
GOARCH=amd64            # Arquitectura target
GO111MODULE=on          # Habilitar módulos
CGO_ENABLED=0           # Deshabilitar C interop

# Ejemplo: Cross-compile para Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app-linux

# Runtime
GOGC=50                 # Garbage collection (% de crecimiento heap antes de GC)
GOMAXPROCS=4            # Max CPUs a usar (goroutines)
GODEBUG=gctrace=1       # Debug output de GC
```

---

## 3.9 Argumentos de Línea de Comando

### Acceder a Argumentos

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // os.Args es []string
    // os.Args[0] es el nombre del programa
    // os.Args[1], os.Args[2], etc. son argumentos

    fmt.Println("Nombre del programa:", os.Args[0])
    fmt.Println("Total de argumentos:", len(os.Args)-1)
    fmt.Println("Todos:", os.Args)
}
```

**Ejecutar:**

```bash
go run main.go arg1 arg2 arg3

# Output:
# Nombre del programa: /tmp/main
# Total de argumentos: 3
# Todos: [/tmp/main arg1 arg2 arg3]
```

### Procesar Argumentos

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Uso: suma <num1> <num2>")
        os.Exit(1)
    }

    // Convertir string a int
    num1, _ := strconv.Atoi(os.Args[1])
    num2, _ := strconv.Atoi(os.Args[2])

    resultado := num1 + num2
    fmt.Printf("%d + %d = %d\n", num1, num2, resultado)
}
```

**Ejecutar:**

```bash
go build -o suma
./suma 5 3
# Output: 5 + 3 = 8
```

### Banderas (Flags)

Go tiene paquete `flag` para banderas:

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // Definir flags
    nombre := flag.String("nombre", "Mundo", "Nombre a saludar")
    verbose := flag.Bool("verbose", false, "Modo verbose")

    // Parsear argumentos
    flag.Parse()

    fmt.Printf("Hola, %s!\n", *nombre)

    if *verbose {
        fmt.Println("Modo verbose habilitado")
    }
}
```

**Ejecutar:**

```bash
go run main.go -nombre Juan -verbose
# Output:
# Hola, Juan!
# Modo verbose habilitado
```

---

## 3.10 Tu Primer Programa Real: Aplicación Interactiva

### Requisitos

Crear una aplicación que:

1. Pida nombre al usuario
2. Pida edad del usuario
3. Valide que edad sea número
4. Imprima información

### Versión 1: Básica

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    // Pedir nombre
    fmt.Print("¿Cuál es tu nombre? ")
    nombre, _ := reader.ReadString('\n')
    nombre = strings.TrimSpace(nombre)  // Remover newline

    // Pedir edad
    fmt.Print("¿Cuál es tu edad? ")
    edadStr, _ := reader.ReadString('\n')
    edadStr = strings.TrimSpace(edadStr)

    // Convertir edad a número
    edad, err := strconv.Atoi(edadStr)
    if err != nil {
        fmt.Println("Error: edad debe ser un número")
        os.Exit(1)
    }

    // Mostrar información
    fmt.Printf("\nHola, %s!\n", nombre)
    fmt.Printf("Tienes %d años.\n", edad)

    // Calcular algo
    if edad >= 18 {
        fmt.Println("Eres mayor de edad.")
    } else {
        fmt.Println("Eres menor de edad.")
    }
}
```

**Ejecutar:**

```bash
go build -o app
./app

# Interacción:
# ¿Cuál es tu nombre? Juan
# ¿Cuál es tu edad? 25
#
# Hola, Juan!
# Tienes 25 años.
# Eres mayor de edad.
```

### Versión 2: Mejorada con Validación

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    // Pedir nombre
    nombre := leerLinea(reader, "¿Cuál es tu nombre? ")
    if nombre == "" {
        fmt.Println("Error: nombre no puede estar vacío")
        os.Exit(1)
    }

    // Pedir edad válida
    edad := leerEdad(reader)

    // Mostrar información
    mostrarInfo(nombre, edad)
}

func leerLinea(reader *bufio.Reader, prompt string) string {
    fmt.Print(prompt)
    linea, _ := reader.ReadString('\n')
    return strings.TrimSpace(linea)
}

func leerEdad(reader *bufio.Reader) int {
    for {
        edadStr := leerLinea(reader, "¿Cuál es tu edad? ")

        edad, err := strconv.Atoi(edadStr)
        if err != nil {
            fmt.Println("Error: debe ser un número entero")
            continue
        }

        if edad < 0 || edad > 150 {
            fmt.Println("Error: edad debe estar entre 0 y 150")
            continue
        }

        return edad
    }
}

func mostrarInfo(nombre string, edad int) {
    fmt.Printf("\n=== INFORMACIÓN ===\n")
    fmt.Printf("Nombre: %s\n", nombre)
    fmt.Printf("Edad: %d años\n", edad)

    if edad >= 18 {
        fmt.Println("Estado: Mayor de edad")
    } else {
        fmt.Println("Estado: Menor de edad")
    }

    generacion := determinarGeneracion(edad)
    fmt.Printf("Generación: %s\n", generacion)
}

func determinarGeneracion(edad int) string {
    if edad < 13 {
        return "Niño"
    } else if edad < 20 {
        return "Adolescente"
    } else if edad < 65 {
        return "Adulto"
    } else {
        return "Adulto Mayor"
    }
}
```

**Ejecutar:**

```bash
go build -o app
./app

# Interacción:
# ¿Cuál es tu nombre? Juan
# ¿Cuál es tu edad? abc
# Error: debe ser un número entero
# ¿Cuál es tu edad? 25
#
# === INFORMACIÓN ===
# Nombre: Juan
# Edad: 25 años
# Estado: Mayor de edad
# Generación: Adulto
```

### Conceptos Usados (Adelanto)

```
Esto que escribiste usa:

 bufio.Reader           (Capítulo 5: Lectura de entrada)
 strings.TrimSpace      (Capítulo 4: Strings)
 strconv.Atoi          (Capítulo 4: Conversiones)
 if/else               (Capítulo 6: Control de flujo)
 funciones             (Capítulo 7)
 strings               (Capítulo 4)
 manejo de errores     (Capítulo 16)
 os.Exit               (Capítulo 3: Tu primer programa)

Los próximos capítulos profundizarán en cada uno.
```

---

## Ejercicios del Capítulo 3

### Ejercicio 1: Calculadora Simple

Crea programa que:

1. Pida 2 números
2. Pida operación (+, -, *, /)
3. Muestre resultado
4. Valide división por cero

### Ejercicio 2: Conversor de Temperaturas

Crea programa que:

1. Pida temperatura en Celsius
2. Muestre equivalentes en Fahrenheit y Kelvin
3. Use funciones helper

### Ejercicio 3: Validador de Contraseña

Crea programa que:

1. Pida contraseña al usuario
2. Valide:
   - Mínimo 8 caracteres
   - Contenga al menos 1 número
   - Contenga al menos 1 mayúscula
3. Muestre mensaje de validez

### Ejercicio 4: Adivina el Número

Crea programa que:

1. Genere número aleatorio 1-100
2. Pida al usuario que adivine
3. Diga "más alto" o "más bajo"
4. Cuente intentos
5. Felicite cuando acierta

---

**Fin del Capítulo 3**

---

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/03-primer-programa/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/03-primer-programa):

```bash
cd examples/03-primer-programa
go run .
```
