# Capítulo 20: Paquetes - Organización y reutilización de código

## Índice
1. [¿Qué es un Paquete?](#201-qué-es-un-paquete)
2. [Declaración de Paquete](#202-declaración-de-paquete)
3. [Importación de Paquetes](#203-importación-de-paquetes)
4. [Rutas de Importación](#204-rutas-de-importación)
5. [Visibilidad: Exportación e Importación](#205-visibilidad-exportación-e-importación)
6. [init() Function](#206-init-function)
7. [main Package](#207-main-package)
8. [Package Documentation](#208-package-documentation)
9. [Estructura de Directorios](#209-estructura-de-directorios)
10. [Circular Dependencies](#2010-circular-dependencies)
11. [Buenas Prácticas y Antipatrones](#2011-buenas-prácticas-y-antipatrones)

---

## 20.1 ¿Qué es un Paquete?

### 20.1.1 Concepto Fundamental

Un **paquete** en Go es la unidad fundamental de organización y reutilización de código. Es un mecanismo que permite:

- **Agrupar** código relacionado en unidades lógicas
- **Encapsular** funcionalidad interna de la pública
- **Reutilizar** código en múltiples programas
- **Gestionar** namespaces y evitar conflictos de nombres
- **Compartir** bibliotecas y módulos

Un paquete es simplemente un directorio que contiene uno o más archivos `.go`. Todos los archivos en el mismo directorio pertenecen al mismo paquete y comparten el mismo espacio de nombres.

```
proyecto/
├── main.go              // package main
├── utils/
│   ├── string.go        // package utils
│   └── math.go          // package utils (mismo paquete)
└── config/
    └── config.go        // package config
```

### 20.1.2 Diferencias con Otros Lenguajes

| Aspecto | Go | Python | Java |
|--------|-----|--------|------|
| Unidad | Directorio + namespace | Módulo/paquete | Nombre de paquete |
| Encapsulación | Primera letra (mayúscula/minúscula) | Convención por defecto | Palabras clave (private/public) |
| Nombre | Por directorio | Por nombre de módulo | Reverso de dominio |
| Import | Por ruta completa | Por nombre de módulo | Por nombre de paquete |
| Estandarización | Obligatoria | Flexible | Obligatoria |

### 20.1.3 Características Clave de los Paquetes en Go

**1. Namespace Global**
- Cada paquete define su propio namespace
- Los símbolos se distinguen por su paquete origen
- No hay conflictos de nombres entre paquetes

**2. Visibilidad Simple**
- Exportado: comienza con mayúscula (público)
- No exportado: comienza con minúscula (privado al paquete)
- No hay modificadores como `public`, `private`, `protected`

**3. Compilación Independiente**
- Go compila paquetes de forma independiente
- Las dependencias se resuelven en tiempo de compilación
- No hay runtime linking

**4. Inicialización Ordenada**
- Cada paquete puede tener función `init()`
- Se ejecutan en orden determinista
- Garantiza estado consistente

### 20.1.4 Estructura de un Paquete Mínimo

```go
// archivo: hello/hello.go
package hello

// Constante exportada
const Greeting = "Hello"

// Variable no exportada (privada)
var internalCount = 0

// Función exportada
func Greet(name string) string {
    internalCount++
    return Greeting + ", " + name
}

// Función no exportada (privada)
func incrementCount() {
    internalCount++
}
```

Uso:
```go
package main

import "miproyecto/hello"

func main() {
    msg := hello.Greet("World")
    println(msg)
    // hello.internalCount // Error: no se puede acceder
}
```

### 20.1.5 Importancia en el Ecosistema Go

Los paquetes son centrales en la filosofía de Go:

1. **Código Modular y Mantenible**: Divide tu programa en componentes independientes
2. **Reutilización**: Comparte código con otros desarrolladores
3. **Testing Aislado**: Prueba cada paquete por separado
4. **Escalabilidad**: Organiza proyectos grandes de forma coherente
5. **Documentación**: Estructura clara facilita documentación

---

## 20.2 Declaración de Paquete

### 20.2.1 Sentencia package

La declaración `package` debe ser la **primera línea** de código en todo archivo Go:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hola")
}
```

**Reglas obligatorias:**
- Debe estar presente en **todo** archivo `.go`
- Debe ser la **primera línea** ejecutable (después de comentarios)
- Todos los archivos en el mismo directorio deben tener el **mismo nombre de paquete**
- El nombre debe ser un identificador válido de Go

### 20.2.2 Convenciones de Nombres

**Reglas para nombres de paquetes:**

1. **Minúsculas**: siempre minúsculas (no CamelCase)
   ```go
   package utils     // ✓ Correcto
   package Utils     // ✗ Incorrecto
   package utilS     // ✗ Incorrecto
   ```

2. **Una palabra**: nombres cortos y descriptivos
   ```go
   package fmt       // ✓ format
   package bytes     // ✓ operaciones con bytes
   package strconv   // ✓ string conversion
   ```

3. **Sin guiones ni caracteres especiales**: sólo letras y números
   ```go
   package http2     // ✓ Válido
   package my_util   // ✗ Incorrecto (guión bajo no recomendado)
   package util-lib  // ✗ Inválido (guión no permitido)
   ```

4. **Descriptivos pero breves**: evita nombres genéricos
   ```go
   package parser    // ✓ Bueno
   package p         // ✗ Muy corto
   package util      // ✓ Aceptable
   package utils     // ✗ Menos específico
   ```

5. **Relacionados con el contenido**: el nombre debe reflejar el propósito
   ```go
   package json      // Contiene funciones para JSON
   package csv       // Contiene funciones para CSV
   package crypto    // Contiene funciones criptográficas
   ```

### 20.2.3 Excepciones: main y test

**Paquete main**
```go
package main

func main() {
    // Punto de entrada de la aplicación
}
```

El paquete `main` es especial:
- No se puede importar desde otros paquetes
- Debe contener una función `main()`
- Compilation ejecutable genera un binario

**Paquetes test**
```go
package mypackage_test

import (
    "testing"
    "miproyecto/mypackage"
)

func TestFunction(t *testing.T) {
    // Tests
}
```

Los tests pueden pertenecer al mismo paquete o a `_test`.

### 20.2.4 Múltiples Paquetes en Repositorio

Es común tener múltiples paquetes en un proyecto:

```
github.com/usuario/proyecto/
├── go.mod
├── main.go                    // package main
├── config/
│   └── config.go              // package config
├── database/
│   ├── db.go                  // package database
│   └── query.go               // package database
├── handlers/
│   ├── user.go                // package handlers
│   └── product.go             // package handlers
└── utils/
    ├── string.go              // package utils
    └── math.go                // package utils
```

Cada directorio define su propio paquete.

### 20.2.5 Aliases de Paquetes Internos

Aunque no es práctica común, puedes usar alias dentro del mismo paquete:

```go
package main

// No hay verdaderos "alias" a nivel de paquete
// Pero puedes crear aliases en código:
type (
    User = user.User
    Role = user.Role
)
```

---

## 20.3 Importación de Paquetes

### 20.3.1 Importación Simple

La forma más básica:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hola, Go!")
}
```

La sentencia `import` especifica:
- Ruta del paquete (desde `$GOPATH/src` o `$GOROOT/src`)
- Cómo se accede al paquete en el código

### 20.3.2 Múltiples Importaciones

**Forma explícita:**
```go
import "fmt"
import "os"
import "math"
```

**Forma agrupada (recomendada):**
```go
import (
    "fmt"
    "os"
    "math"
)
```

Gofmt reorganiza automáticamente los imports en este formato.

**Orden alfabético automático:**
```go
import (
    "fmt"           // stdlib
    "os"
    "path/filepath"
    
    "github.com/user/lib1"  // paquetes externos
    "github.com/user/lib2"
    
    "miproyecto/config"     // paquetes internos
    "miproyecto/utils"
)
```

Go automáticamente agrupa y ordena alfabéticamente (stdlib primero, luego externos).

### 20.3.3 Alias de Importación

Útil cuando hay conflictos de nombres o necesitas un alias descriptivo:

```go
package main

import (
    "fmt"
    fmt2 "github.com/other/fmt"  // Alias
)

func main() {
    fmt.Println("stdlib fmt")
    fmt2.Print("custom fmt")
}
```

**Casos de uso:**
- Evitar conflictos de nombres
- Aclarar el origen del paquete
- Nombres más cortos para paquetes complejos

```go
import (
    "database/sql"
    sqlite "github.com/mattn/go-sqlite3"
    postgres "github.com/lib/pq"
)
```

### 20.3.4 Dot Import (.)

Importa todos los símbolos públicos sin prefijo de paquete:

```go
package main

import . "fmt"

func main() {
    Println("Sin prefijo fmt.")  // En lugar de fmt.Println
}
```

**Desventajas:**
- Ambigüedad: no sabes de dónde viene el símbolo
- Dificulta el mantenimiento
- Conflictos potenciales

**Cuándo usar:**
- Tests únicamente
- Scripts simples y pequeños

**Ejemplo correcto (en tests):**
```go
package mypackage

import (
    . "testing"
    . "github.com/stretchr/testify/assert"
)

func TestSomething(t *T) {
    Equal(t, 5, 2+3)
}
```

### 20.3.5 Blank Import (_)

Importa un paquete sin usarlo directamente (por sus efectos secundarios):

```go
package main

import (
    "database/sql"
    _ "github.com/lib/pq"  // Driver PostgreSQL
)

func main() {
    db, _ := sql.Open("postgres", "...")
}
```

**Casos de uso:**
- Registrar drivers de base de datos
- Inicializar side-effects
- Forzar imports de dependencias

**Ejemplo: Registrar handlers HTTP**
```go
package main

import (
    "net/http"
    _ "miproyecto/handlers"  // Registra handlers en init()
)

func main() {
    http.ListenAndServe(":8080", nil)
}
```

En `handlers/handlers.go`:
```go
package handlers

import "net/http"

func init() {
    http.HandleFunc("/api/users", handleUsers)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

### 20.3.6 Manejo de Errores de Importación

**Import no utilizado:**
```go
import "fmt"  // Error: fmt no se usa

// Solución: eliminar o usar
import "fmt"
func main() {
    fmt.Println("Ahora se usa")
}
```

**Paquete no encontrado:**
```go
import "miproyecto/noexistente"  // Error en compilación
```

**Alias no utilizado:**
```go
import util "miproyecto/utils"  // Error si util no se usa
```

Go es estricto: **todo lo que importas debe ser usado**.

---

## 20.4 Rutas de Importación

### 20.4.1 Rutas Absolutas

En Go, todos los imports son **rutas absolutas** desde el módulo raíz:

```go
// NO existen imports relativos como en Python
import "../utils"           // ✗ No permitido
import "./config"           // ✗ No permitido
import "/absolute/path"     // ✗ No permitido

// Siempre absoluto desde el módulo
import "miproyecto/utils"        // ✓ Correcto
import "github.com/user/lib"     // ✓ Correcto
```

### 20.4.2 Go Modules (go.mod)

Go usa módulos para gestionar dependencias. El archivo `go.mod` define:

```
go.mod
├── module path (nombre único)
├── versión de Go
└── dependencias con versiones
```

**Archivo go.mod típico:**
```
module github.com/usuario/miproyecto

go 1.21

require (
    github.com/lib/pq v1.10.9
    github.com/google/uuid v1.3.0
)
```

**Instalación de dependencias:**
```bash
go get github.com/lib/pq              # Agrega a go.mod
go get -u github.com/lib/pq           # Actualiza a versión más reciente
go mod tidy                            # Limpia dependencias no utilizadas
```

### 20.4.3 Rutas de Paquetes Internos

Dentro de tu propio módulo:

```go
// En github.com/usuario/miproyecto
import (
    "github.com/usuario/miproyecto/config"
    "github.com/usuario/miproyecto/database"
)
```

**Estructura de proyecto:**
```
github.com/usuario/miproyecto/
├── go.mod
├── main.go
├── config/
│   └── config.go
└── database/
    └── db.go
```

**En main.go:**
```go
package main

import (
    "github.com/usuario/miproyecto/config"
    "github.com/usuario/miproyecto/database"
)

func main() {
    cfg := config.Load()
    db := database.Connect(cfg)
    // ...
}
```

### 20.4.4 Rutas de Paquetes Estándar

La biblioteca estándar de Go:

```go
import (
    "fmt"                    // fmt
    "os"                     // os
    "path/filepath"          // path subpackage
    "encoding/json"          // encoding subpackage
    "crypto/sha256"          // crypto subpackage
    "database/sql"           // database subpackage
)
```

### 20.4.5 Rutas de Módulos Externos

De GitHub y otros repositorios públicos:

```go
import (
    "github.com/lib/pq"                      // PostgreSQL driver
    "github.com/google/uuid"                 // Google UUID library
    "github.com/stretchr/testify/assert"     // Testify assertion library
    "gopkg.in/yaml.v3"                       // YAML v3
)
```

Estas se descargan automáticamente con `go get`.

### 20.4.6 Gestión de Versiones

Go Modules utiliza versionado semántico:

```go
// En go.mod
require (
    github.com/lib/pq v1.10.9      // Versión exacta
    github.com/google/uuid v1.3.0
)

// En código (sin cambios en imports)
import "github.com/lib/pq"
```

**Cambio de versión:**
```bash
go get github.com/lib/pq@v1.11.0     # Versión específica
go get github.com/lib/pq@latest      # Última versión
go get github.com/lib/pq@main        # Branch principal
```

### 20.4.7 Reemplazo de Módulos

Para development local o forks:

```
// go.mod
module github.com/usuario/miproyecto

require github.com/lib/pq v1.10.9

replace github.com/lib/pq => /local/path/to/pq
```

Útil para:
- Testing de cambios no publicados
- Trabajar con forks locales
- Monorepos

---

## 20.5 Visibilidad: Exportación e Importación

### 20.5.1 Principio de Visibilidad en Go

**Regla simple:** la primera letra del identificador determina la visibilidad.

```go
const (
    Exported   = "público (mayúscula)"    // Visible fuera del paquete
    notExported = "privado (minúscula)"    // Privado al paquete
)

func Visible() string {      // Exportado
    return notExported       // Privado, pero accesible aquí
}

func hidden() string {       // No exportado
    return "privado"
}
```

### 20.5.2 Exportación de Identificadores

**Qué se puede exportar:**

```go
package mypackage

// Constantes exportadas
const (
    MaxSize = 1000
    Version = "1.0.0"
)

// Variables exportadas
var (
    DefaultTimeout = 30
    GlobalConfig   *Config
)

// Tipos exportados
type User struct {
    ID   int
    Name string
}

// Funciones exportadas
func ProcessData(input string) string {
    return input
}

// Métodos exportados
func (u *User) String() string {
    return u.Name
}

// Interfaces exportadas
type Reader interface {
    Read() ([]byte, error)
}
```

**No exportados (privados al paquete):**

```go
package mypackage

const minSize = 1         // Privado
var internalState = 0     // Privado

type internalConfig struct {  // Privado
    timeout int
}

func helper() string {    // Privado
    return "interno"
}
```

### 20.5.3 Patrones de Acceso Controlado

**Getter para acceso controlado:**

```go
package config

// Privado
type database struct {
    host     string
    password string  // ¡Sensible!
}

// Exportado: acceso controlado
func (d *database) Host() string {
    return d.host
}

// NO exponemos el password
// func (d *database) Password() string {  // ¡No hacer!
```

**Constructor patrón:**

```go
package config

type Database struct {
    host     string      // privado
    port     int         // privado
    username string
}

// Constructor: validación y asignación segura
func NewDatabase(host string, port int) *Database {
    if port < 1 || port > 65535 {
        port = 5432  // default
    }
    return &Database{
        host:     host,
        port:     port,
        username: "admin",
    }
}

func (d *Database) Host() string {
    return d.host
}
```

**Builder pattern para configuración compleja:**

```go
type QueryBuilder struct {
    table    string
    where    []string
    orderBy  string
}

func Query() *QueryBuilder {
    return &QueryBuilder{}
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
    q.table = table
    return q
}

func (q *QueryBuilder) Where(condition string) *QueryBuilder {
    q.where = append(q.where, condition)
    return q
}

func (q *QueryBuilder) OrderBy(field string) *QueryBuilder {
    q.orderBy = field
    return q
}

// Uso:
query := Query().From("users").Where("age > 18").OrderBy("name")
```

### 20.5.4 Interfaces Públicas, Implementación Privada

Patrón importante en Go:

```go
package storage

// Interfaz pública: define el contrato
type Storage interface {
    Get(key string) (string, error)
    Set(key, value string) error
    Delete(key string) error
}

// Implementación privada
type memoryStorage struct {
    data map[string]string
}

// Constructor exportado devuelve interfaz
func NewMemoryStorage() Storage {
    return &memoryStorage{
        data: make(map[string]string),
    }
}

func (m *memoryStorage) Get(key string) (string, error) {
    val, ok := m.data[key]
    if !ok {
        return "", ErrNotFound
    }
    return val, nil
}

func (m *memoryStorage) Set(key, value string) error {
    m.data[key] = value
    return nil
}

func (m *memoryStorage) Delete(key string) error {
    delete(m.data, key)
    return nil
}
```

Uso:
```go
var store storage.Storage = storage.NewMemoryStorage()
store.Set("key", "value")
```

### 20.5.5 Documentación de Visibilidad

```go
// Esta función está exportada (mayúscula)
// GenerateToken crea un token JWT válido
func GenerateToken(userID string) (string, error) {
    return encodeToken(userID, time.Now().Add(24*time.Hour))
}

// Esta función es privada (minúscula)
// encodeToken hace el trabajo real
func encodeToken(userID string, expiry time.Time) (string, error) {
    // ...
}
```

---

## 20.6 init() Function

### 20.6.1 Qué es init()

`init()` es una función especial ejecutada automáticamente cuando:
- Se importa el paquete
- Se carga el programa

**Características:**
- Se ejecuta **antes** de `main()`
- No recibe argumentos
- No devuelve nada
- Se ejecuta automáticamente (no se llama manualmente)

```go
package main

import "fmt"

func init() {
    fmt.Println("1. Ejecuto primero")
}

func main() {
    fmt.Println("2. Ejecuto después")
}

// Salida:
// 1. Ejecuto primero
// 2. Ejecuto después
```

### 20.6.2 Orden de Ejecución

Cuando un programa comienza, Go ejecuta:

1. **Inicialización de variables globales** (en orden de declaración)
2. **Funciones init()** (en orden de declaración en el archivo)
3. **main()**

```go
package main

import "fmt"

var (
    step1 = print("A. Variable global\n")
)

func init() {
    fmt.Println("B. init() primera")
}

func init() {
    fmt.Println("C. init() segunda")
}

func main() {
    fmt.Println("D. main()")
}

// Salida:
// A. Variable global
// B. init() primera
// C. init() segunda
// D. main()
```

### 20.6.3 Múltiples init() en el Mismo Paquete

Un paquete puede tener múltiples funciones `init()`:

```go
// config/init.go
package config

import "log"

var ConfigPath string

func init() {
    log.Println("Inicializando config")
    ConfigPath = "/etc/app.conf"
}

// config/database.go
package config

import "log"

var DBConnection string

func init() {
    log.Println("Inicializando database")
    // Usa la variable ya inicializada
    if ConfigPath == "" {
        log.Fatal("ConfigPath debe estar inicializado")
    }
}
```

La ejecución es dentro del mismo archivo primero, luego orden alfabética.

### 20.6.4 Inicialización de Múltiples Paquetes

Cuando importas varios paquetes, cada uno ejecuta su `init()`:

```go
// main.go
package main

import (
    "miproyecto/config"
    "miproyecto/database"
    "miproyecto/logger"
)

func main() {
    // Todos los init() ya ejecutados
}
```

Orden de ejecución (aproximado):
1. `config.init()`
2. `database.init()`
3. `logger.init()`
4. `main()`

### 20.6.5 Usos Prácticos de init()

**1. Inicialización de registros**

```go
// handlers/routes.go
package handlers

import "net/http"

func init() {
    http.HandleFunc("/api/users", handleUsers)
    http.HandleFunc("/api/products", handleProducts)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}
```

**2. Carga de configuración**

```go
// config/config.go
package config

import (
    "flag"
    "log"
    "os"
)

var AppPort = ":8080"
var AppEnv = "development"

func init() {
    port := os.Getenv("PORT")
    if port != "" {
        AppPort = ":" + port
    }
    
    env := os.Getenv("ENV")
    if env != "" {
        AppEnv = env
    }
    
    log.Printf("App iniciada en puerto %s (env: %s)", AppPort, AppEnv)
}
```

**3. Setup de bases de datos**

```go
// database/db.go
package database

import (
    "database/sql"
    "log"
)

var DB *sql.DB

func init() {
    var err error
    DB, err = sql.Open("postgres", "user=app dbname=myapp")
    if err != nil {
        log.Fatal("No se conectó a la base de datos:", err)
    }
    
    if err := DB.Ping(); err != nil {
        log.Fatal("DB no responde:", err)
    }
    
    log.Println("Base de datos conectada")
}
```

**4. Registro de drivers**

```go
// main.go
package main

import (
    "database/sql"
    _ "github.com/lib/pq"  // Ejecuta init() del driver
)

func main() {
    db, _ := sql.Open("postgres", "...")
    // El driver está registrado por init()
}
```

### 20.6.6 Errores Comunes con init()

**Error 1: Dependencias en orden incorrecto**

```go
// ✗ Incorrecto: intenta usar variable aún no inicializada
package config

var (
    config = loadConfig()      // Problema: loadConfig() depende de env vars
    port   = config.Port
)

func loadConfig() *Config {
    // Error si se ejecuta antes de que env vars se carguen
    return &Config{}
}

// ✓ Correcto: usar init()
package config

var (
    config *Config
)

func init() {
    config = loadConfig()  // Ejecuta después de que todo esté listo
}
```

**Error 2: Ciclo de inicialización**

```go
// ✗ Incorrecto: A depende de B, B depende de A
// package a:
func init() {
    _ = b.GetValue()  // Llama a B
}

// package b:
func init() {
    _ = a.GetValue()  // Llama a A (¡ciclo!)
}
```

Solución: refactorizar dependencias.

---

## 20.7 main Package

### 20.7.1 Especificidad del Paquete main

El paquete `main` es especial en Go. Es el único paquete que:
- **Genera** un ejecutable binario
- **No se puede importar** desde otros paquetes
- **Debe contener** una función `main()`

```go
package main

func main() {
    println("Programa ejecutable")
}

// Compilado con: go build
// Genera: main (o main.exe en Windows)
```

### 20.7.2 Función main()

```go
package main

func main() {
    // Punto de entrada obligatorio
    // No puede ser importada
    // No puede devolver nada
    // Recibe argumentos por os.Args
}
```

**Acceso a argumentos de línea de comandos:**

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Argumentos:")
    for i, arg := range os.Args {
        fmt.Printf("  [%d]: %s\n", i, arg)
    }
}

// Ejecución: go run main.go arg1 arg2
// Salida:
// Argumentos:
//   [0]: /tmp/.../main
//   [1]: arg1
//   [2]: arg2
```

### 20.7.3 Estructura Típica de main Package

```go
package main

import (
    "flag"
    "log"
    "miproyecto/config"
    "miproyecto/server"
)

// Flags globales (parse de línea de comandos)
var (
    port    = flag.String("port", ":8080", "Puerto del servidor")
    debug   = flag.Bool("debug", false, "Modo debug")
    logFile = flag.String("log", "", "Archivo de log")
)

// init() para setup inicial
func init() {
    flag.Parse()
    
    if *logFile != "" {
        // Configurar logging a archivo
    }
}

// main() simple y limpia
func main() {
    cfg := config.Load(*debug)
    srv := server.New(cfg)
    
    if err := srv.Start(*port); err != nil {
        log.Fatal(err)
    }
}
```

### 20.7.4 Ventajas de Separar Lógica de main

**Estructura recomendada:**

```
proyecto/
├── main.go              // package main (simple)
├── config/
│   └── config.go        // package config (lógica)
├── server/
│   └── server.go        // package server (lógica)
└── handlers/
    └── handlers.go      // package handlers (lógica)
```

**main.go minimalista:**
```go
package main

import (
    "log"
    "miproyecto/config"
    "miproyecto/server"
)

func main() {
    cfg := config.New()
    srv := server.New(cfg)
    
    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}
```

**Ventajas:**
- Testeable: `config` y `server` se pueden testear
- No testeable: `main` es difícil de testear (es punto de entrada)
- Reutilizable: otros programas pueden importar `config`, `server`
- Mantenible: separación clara de responsabilidades

### 20.7.5 Programa Multi-binario (múltiples main packages)

Puedes tener múltiples puntos de entrada:

```
proyecto/
├── cmd/
│   ├── server/
│   │   └── main.go      // package main (servidor)
│   ├── cli/
│   │   └── main.go      // package main (herramienta CLI)
│   └── migrator/
│       └── main.go      // package main (migraciones)
├── config/
│   └── config.go        // package config
└── database/
    └── db.go            // package database
```

**Compilación:**
```bash
go build -o bin/server ./cmd/server
go build -o bin/cli ./cmd/cli
go build -o bin/migrator ./cmd/migrator
```

Cada binario es independiente pero comparten lógica de `config`, `database`.

---

## 20.8 Package Documentation

### 20.8.1 Comentarios de Documentación

Go automáticamente convierte comentarios en documentación usando `godoc`.

**Regla: El comentario inmediatamente antes de un identificador es su documentación.**

```go
package mypackage

// User representa a un usuario del sistema
type User struct {
    ID   int    // Identificador único
    Name string // Nombre completo del usuario
}

// NewUser crea un nuevo usuario con validación
func NewUser(name string) *User {
    return &User{
        Name: name,
    }
}

// Greet devuelve un saludo para el usuario
func (u *User) Greet() string {
    return "Hola, " + u.Name
}
```

### 20.8.2 Documentación de Paquetes

El comentario de documentación del **paquete** debe estar en un archivo, antes de la declaración `package`:

```go
// Package mypackage proporciona utilidades para gestionar usuarios
// y sus perfiles.
//
// Uso básico:
//
//  user := mypackage.NewUser("Alice")
//  greeting := user.Greet()
//
// Para operaciones avanzadas, ver la documentación en línea.
package mypackage
```

Esto aparece en `godoc`.

### 20.8.3 Estilo de Documentación en Go

**Reglas de formato:**

1. **Empieza con el nombre del identificador**
   ```go
   // User es un usuario del sistema
   // ✓ Correcto: empieza con el nombre
   
   // Este es un usuario
   // ✗ Incorrecto: no empieza con el nombre
   ```

2. **Oraciones completas**
   ```go
   // NewUser crea un nuevo usuario con la validación requerida
   // ✓ Correcto
   
   // Crear usuario
   // ✗ Menos formal
   ```

3. **Formato limpio**
   ```go
   // GenerateToken crea un token JWT para autenticación.
   // El token expira después de 24 horas.
   //
   // Parameters:
   //   userID: el ID del usuario
   //
   // Returns:
   //   el token JWT
   //   error si ocurre un problema
   func GenerateToken(userID string) (string, error)
   ```

### 20.8.4 Generación de Documentación con godoc

**Ver documentación en terminal:**
```bash
godoc miproyecto/mypackage
godoc -http=:6060              # Abre servidor web en localhost:6060
```

**Ejemplo de salida:**
```
PACKAGE

package mypackage

TYPES

type User struct {
    ID   int
    Name string
}
    User representa a un usuario del sistema

func NewUser(name string) *User
    NewUser crea un nuevo usuario con validación

func (u *User) Greet() string
    Greet devuelve un saludo para el usuario

FUNCTIONS

func GenerateToken(userID string) (string, error)
    GenerateToken crea un token JWT para autenticación.
```

### 20.8.5 Ejemplos en Documentación

```go
package mypackage

// ExampleUser muestra cómo usar el tipo User
func ExampleUser() {
    user := NewUser("Alice")
    fmt.Println(user.Greet())
    // Output: Hola, Alice
}

// ExampleUser_Greet muestra el método Greet
func ExampleUser_Greet() {
    user := NewUser("Bob")
    greeting := user.Greet()
    fmt.Println(greeting)
    // Output: Hola, Bob
}
```

Los ejemplos se ejecutan como tests (output debe coincidir).

### 20.8.6 Deprecación y Anotaciones

```go
package mypackage

// Deprecated: uso NewUserV2 en su lugar
func NewUser(name string) *User {
    return NewUserV2(name)
}

// TODO: refactorizar este método para mejorar performance
func (u *User) Process() error {
    // ...
}

// BUG(usuario): esta función no maneja valores negativos correctamente
func Calculate(value int) int {
    // ...
}
```

---

## 20.9 Estructura de Directorios

### 20.9.1 Estructura Estándar de Proyecto Go

```
github.com/usuario/myproject/
├── go.mod                    # Módulo definition
├── go.sum                    # Checksums de dependencias
├── README.md
├── LICENSE
├── Makefile                  # (opcional) tareas comunes
│
├── cmd/                      # Ejecutables/comandos
│   ├── server/
│   │   └── main.go
│   └── cli/
│       └── main.go
│
├── pkg/                      # Paquetes públicos (importables)
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── db.go
│   │   └── query.go
│   └── handlers/
│       ├── user.go
│       └── product.go
│
├── internal/                 # Paquetes privados (no importables)
│   ├── models/
│   │   └── models.go
│   └── utils/
│       └── utils.go
│
├── tests/                    # Tests de integración
│   ├── integration_test.go
│   └── fixtures/
│
└── docs/                     # Documentación
    └── api.md
```

### 20.9.2 Convención: cmd/ para ejecutables

El directorio `cmd/` es para aplicaciones ejecutables:

```
cmd/
├── server/
│   └── main.go              # Servidor principal
├── cli/
│   └── main.go              # Herramienta CLI
├── worker/
│   └── main.go              # Worker de background jobs
└── migrator/
    └── main.go              # Herramienta de migraciones
```

**Compilación:**
```bash
go build -o bin/server ./cmd/server
go build -o bin/cli ./cmd/cli
```

### 20.9.3 Convención: pkg/ para librerías públicas

Paquetes que **pueden ser importados** por otros proyectos:

```
pkg/
├── config/
│   └── config.go
├── database/
│   ├── db.go
│   └── query.go
└── api/
    └── client.go
```

Importable:
```go
import "github.com/usuario/myproject/pkg/config"
```

### 20.9.4 Convención: internal/ para código privado

Paquetes que **no pueden ser importados** (Go lo impide automáticamente):

```
internal/
├── models/
│   └── models.go
├── middleware/
│   └── auth.go
└── utils/
    └── utils.go
```

**Go rechaza:**
```go
// Desde otro proyecto:
import "github.com/usuario/myproject/internal/models"  // Error: 403
```

**Pero está permitido dentro del proyecto:**
```go
// En cmd/server/main.go:
import "github.com/usuario/myproject/internal/models"  // ✓ Permitido
```

### 20.9.5 Multiarchivo en un Paquete

Un paquete puede distribuirse en varios archivos:

```
database/
├── db.go            // package database
├── query.go         // package database
├── migration.go     // package database
└── db_test.go       // package database_test
```

**Reglas:**
- Todos deben ser `package database`
- Se importan como un único paquete
- Los tests van en `_test.go`

**Ejemplo: db.go**
```go
package database

import "database/sql"

var DB *sql.DB

func Connect(dsn string) error {
    var err error
    DB, err = sql.Open("postgres", dsn)
    return err
}
```

**Ejemplo: query.go**
```go
package database

type Query struct {
    table string
}

func QueryBuilder() *Query {
    return &Query{}
}

func (q *Query) From(table string) *Query {
    q.table = table
    return q
}
```

Importación simple:
```go
import "miproyecto/database"

func main() {
    database.Connect("...")
    q := database.QueryBuilder()
}
```

### 20.9.6 Ventajas de Estructura Clara

```
                    Estructura Clara
                          ↓
        ┌──────────────────┴──────────────────┐
        ↓                                      ↓
   Fácil de Navegar                    Fácil de Testear
        ↓                                      ↓
   Nuevos Desarrolladores            Aislamiento de Pruebas
   Rápida Onboarding                  Menos Dependencias
   Menos Confusión                    Tests más Rápidos
```

---

## 20.10 Circular Dependencies

### 20.10.1 Qué es una Dependencia Circular

Una dependencia circular ocurre cuando:
- El paquete A importa al paquete B
- El paquete B importa al paquete A

```
    A
    ↑
    |
    ↓
    B
    ↑
    |
    └─────────────→ A (ciclo)
```

### 20.10.2 Problema de Compilación

Go **rechaza** las dependencias circulares en tiempo de compilación:

**Código problemático:**

```go
// package a/a.go
package a

import "miproyecto/b"

func UsesB() {
    b.FunctionB()
}

// package b/b.go
package b

import "miproyecto/a"  // ¡Error! Ciclo

func FunctionB() {
    a.UsesB()  // Ciclo infinito
}
```

**Error de compilación:**
```
circular import not allowed
    miproyecto/a
    miproyecto/b
    miproyecto/a
```

### 20.10.3 Identificar Ciclos

**Caso 1: Ciclo directo**

```
package a → package b → package a (ciclo)
```

**Caso 2: Ciclo indirecto**

```
package a → package b → package c → package a (ciclo)
```

**Caso 3: Ciclo con múltiples niveles**

```
package handler → package service → package model → package handler
```

### 20.10.4 Soluciones

**Solución 1: Refactorizar en un único paquete**

```
// ✗ Antes: ciclo
a/a.go (usa b)
b/b.go (usa a)

// ✓ Después: un único paquete
combined/
├── a.go
├── b.go
└── combined.go  // package combined
```

**Solución 2: Extraer interfaz compartida**

```go
// ✗ Problemático:
// package a usa b, b usa a

// ✓ Solución: paquete de interfaces compartidas
package interfaces

type Handler interface {
    Handle() error
}

type Service interface {
    Process() error
}
```

Uso:

```go
// package a
package a

import "miproyecto/interfaces"

type Handler struct {
    service interfaces.Service
}

// package b
package b

import "miproyecto/interfaces"

type Service struct {
    // Accede a Handler a través de interfaz
}
```

**Solución 3: Invertir dependencia (Dependency Injection)**

```go
// ✗ Antes: b depende de a (cyclo potencial)
package b

import "miproyecto/a"

type Processor struct {
    validator a.Validator
}

// ✓ Después: inyectar dependencia
package b

type Processor struct {
    validator Validator  // Interfaz local
}

type Validator interface {
    Validate(input string) error
}
```

**Solución 4: Crear paquete mediador**

```
// ✗ Ciclo entre api y database

// ✓ Solución:
mediator/
├── mediator.go        // Coordina ambos

api/ (sin depender de database)
database/ (sin depender de api)
```

```go
// mediator/mediator.go
package mediator

import (
    "miproyecto/api"
    "miproyecto/database"
)

type Mediator struct {
    api      api.Handler
    database database.Store
}

func (m *Mediator) Process(request *api.Request) {
    data := m.database.Get(request.ID)
    m.api.Respond(data)
}
```

### 20.10.5 Buenas Prácticas para Evitar Ciclos

**Patrón: Arquitectura Layered**

```
┌─────────────────┐
│   main package  │
└────────┬────────┘
         ↓
┌─────────────────┐
│  cmd/cli        │
└────────┬────────┘
         ↓
┌─────────────────┐
│   handlers      │  ← Capa superior
├─────────────────┤
│   services      │  ← Capa media
├─────────────────┤
│   models        │  ← Capa inferior
├─────────────────┤
│   database      │  ← Datos
└─────────────────┘

Regla: Solo hacia abajo, nunca hacia arriba
```

---

## 20.11 Buenas Prácticas y Antipatrones

### 20.11.1 Buena Práctica: Nombres Descriptivos

```go
// ✓ Bueno: específico, descriptivo
package encoder
package parser
package tokenizer

// ✗ Malo: genérico, ambiguo
package util
package helper
package common
```

### 20.11.2 Buena Práctica: Cohesión Alta

Los miembros de un paquete deben estar altamente relacionados:

```go
// ✓ Bueno: cohesión alta
package json

type Decoder struct { /* ... */ }
func NewDecoder() *Decoder { /* ... */ }
func (d *Decoder) Decode() { /* ... */ }

// ✗ Malo: baja cohesión
package utils

type User struct { /* ... */ }
type Car struct { /* ... */ }
type APIClient struct { /* ... */ }
func RandomNumber() int { /* ... */ }
func ParseJSON(data string) { /* ... */ }
func WriteFile(path string) { /* ... */ }
// Todo lo que no sabía dónde poner va aquí
```

### 20.11.3 Buena Práctica: Bajo Acoplamiento

Los paquetes deben ser independientes:

```go
// ✓ Bueno: bajo acoplamiento (usa interfaces)
package service

type Database interface {
    Get(id string) (string, error)
}

type Service struct {
    db Database
}

func NewService(db Database) *Service {
    return &Service{db: db}
}

// ✗ Malo: alto acoplamiento (depende directamente)
package service

import "miproyecto/database"

type Service struct {
    db *database.PostgreSQL  // Acoplado a implementación específica
}
```

### 20.11.4 Buena Práctica: API Pública Clara

```go
// ✓ Bueno: interfaz clara
package storage

type Storage interface {
    Get(key string) ([]byte, error)
    Set(key string, value []byte) error
    Delete(key string) error
}

func NewMemoryStorage() Storage { /* ... */ }

// ✗ Malo: interfaz confusa
package storage

func GetFromDB(key string) string { /* ... */ }
func WriteToCache(key string, val string) { /* ... */ }
func RemoveFromStorage(id string) { /* ... */ }
func ValidateInput(data string) bool { /* ... */ }
func InternalHelper() { /* ... */ }
// No está claro cuál es la API pública
```

### 20.11.5 Buena Práctica: Documentación Completa

```go
// ✓ Bueno: bien documentado
package parser

// Parser procesa tokens de entrada
type Parser struct {
    tokens []Token
    pos    int
}

// NewParser crea un nuevo parser
func NewParser(input string) (*Parser, error) {
    tokens, err := tokenize(input)
    if err != nil {
        return nil, err
    }
    return &Parser{tokens: tokens}, nil
}

// Parse analiza los tokens y devuelve un AST
func (p *Parser) Parse() (*AST, error) {
    // ...
}

// ✗ Malo: sin documentación
package parser

type Parser struct {
    t []Token
    p int
}

func NewParser(i string) (*Parser, error) {
    // ...
}

func (p *Parser) Parse() (*AST, error) {
    // ...
}
```

### 20.11.6 Antipatrón: Paquete God

```go
// ✗ Antipatrón: paquete que hace todo
package app

type User struct { /* ... */ }
type Product struct { /* ... */ }
type Order struct { /* ... */ }

func ValidateUser(u *User) error { /* ... */ }
func ValidateProduct(p *Product) error { /* ... */ }
func ValidateOrder(o *Order) error { /* ... */ }

func SaveToDatabase(data interface{}) error { /* ... */ }
func SendEmail(recipient string, subject string) error { /* ... */ }
func LogEvent(event string) error { /* ... */ }

func ProcessPayment(amount float64) error { /* ... */ }
func GenerateReport() string { /* ... */ }

// ... 500 líneas más

// ✓ Solución: dividir en paquetes especializados
package domain
    type User struct { /* ... */ }
    type Product struct { /* ... */ }
    type Order struct { /* ... */ }

package validator
    func User(u *domain.User) error { /* ... */ }
    func Product(p *domain.Product) error { /* ... */ }

package storage
    func Save(data interface{}) error { /* ... */ }

package notification
    func SendEmail(recipient, subject string) error { /* ... */ }

package payment
    func Process(amount float64) error { /* ... */ }
```

### 20.11.7 Antipatrón: Importación Anidada

```go
// ✗ Malo: estructura confusa
package a
import "miproyecto/b"

package b
import "miproyecto/c"

package c
import "miproyecto/d"

package d
import "miproyecto/e"
// Difícil seguir la cadena de dependencias

// ✓ Bueno: dependencias claras
package handler
import (
    "miproyecto/service"
    "miproyecto/storage"
)

package service
import (
    "miproyecto/storage"
    "miproyecto/validator"
)

// Jerarquía clara: handler → service → storage
```

### 20.11.8 Buena Práctica: Versionado de Paquetes

```go
// Si necesitas cambios incompatibles:

// ✓ Opción 1: nuevo paquete con versión
package userv2

type User struct {
    id       int
    email    string
    verified bool
}

// ✓ Opción 2: branch/etiqueta en git
// v1.0.0
// v2.0.0 (cambios incompatibles)

// Uso en go.mod
// require github.com/usuario/proyecto v2.0.0
```

### 20.11.9 Checklist de Buenas Prácticas

```
□ Nombre descriptivo y minúsculo
□ API pública clara (exportados documentados)
□ Encapsulación (privados para detalles internos)
□ Sin dependencias circulares
□ Cohesión alta (miembros relacionados)
□ Bajo acoplamiento (usa interfaces)
□ Documentación completa (godoc-ready)
□ Tests incluidos (*_test.go)
□ No es un "paquete god"
□ Estructura de directorios clara
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: Crear Paquete Utils con Documentación

**Objetivo**: Crear un paquete `utils` con funciones útiles, documentadas y bien organizadas.

**Requisitos:**
- Crear directorio `utils/`
- Tres funciones exportadas: `Reverse()`, `ToUpperWords()`, `CountWords()`
- Documentación completa con comentarios
- Archivo `utils_test.go` con tests básicos

**Estructura:**
```
ejercicio1/
├── go.mod
├── main.go
└── utils/
    ├── string.go
    └── string_test.go
```

**Solución:**

```go
// go.mod
module ejercicio1

go 1.21
```

```go
// utils/string.go
package utils

import "strings"

// Reverse invierte una cadena de caracteres
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// ToUpperWords convierte cada palabra a mayúsculas
func ToUpperWords(s string) string {
    return strings.Title(s)
}

// CountWords cuenta la cantidad de palabras en una cadena
func CountWords(s string) int {
    words := strings.Fields(s)
    return len(words)
}
```

```go
// utils/string_test.go
package utils

import "testing"

func TestReverse(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"hello", "olleh"},
        {"Go", "oG"},
        {"", ""},
    }
    
    for _, tt := range tests {
        got := Reverse(tt.input)
        if got != tt.expected {
            t.Errorf("Reverse(%q) = %q, want %q", tt.input, got, tt.expected)
        }
    }
}
```

```go
// main.go
package main

import (
    "fmt"
    "ejercicio1/utils"
)

func main() {
    fmt.Println("Reverse:", utils.Reverse("Hello Go"))
    fmt.Println("Words:", utils.CountWords("The quick brown fox"))
}
```

### Ejercicio 2: Paquete Configuración con init()

**Objetivo**: Crear paquete `config` que carga configuración y utiliza `init()`.

**Requisitos:**
- Función `init()` que carga valores por defecto
- Constructor `New()` que valida configuración
- Implementar valores por defecto de environment variables
- Método para acceder (patrón getter privado)

**Estructura:**
```
ejercicio2/
├── go.mod
├── main.go
└── config/
    ├── config.go
    └── config_test.go
```

**Solución:**

```go
// config/config.go
package config

import (
    "fmt"
    "os"
    "log"
)

type Config struct {
    appName string
    port    string
    debug   bool
}

var (
    defaults *Config
)

func init() {
    log.Println("Inicializando configuración...")
    defaults = &Config{
        appName: "MyApp",
        port:    ":8080",
        debug:   false,
    }
}

func New() *Config {
    cfg := &Config{
        appName: getEnv("APP_NAME", defaults.appName),
        port:    getEnv("PORT", defaults.port),
        debug:   getEnv("DEBUG", "false") == "true",
    }
    return cfg
}

func (c *Config) AppName() string {
    return c.appName
}

func (c *Config) Port() string {
    return c.port
}

func (c *Config) IsDebug() bool {
    return c.debug
}

func getEnv(key, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultVal
}
```

```go
// main.go
package main

import (
    "fmt"
    "ejercicio2/config"
)

func main() {
    cfg := config.New()
    fmt.Printf("App: %s, Port: %s, Debug: %v\n", 
        cfg.AppName(), cfg.Port(), cfg.IsDebug())
}
```

### Ejercicio 3: Paquete Matemático Multi-archivo

**Objetivo**: Crear un paquete `math` con múltiples archivos sin exponer detalles internos.

**Requisitos:**
- Archivo `basic.go`: Funciones básicas exportadas
- Archivo `advanced.go`: Funciones avanzadas
- Archivo `internal.go`: Funciones privadas
- Implementar al menos 5 operaciones

**Estructura:**
```
ejercicio3/
├── go.mod
├── main.go
└── math/
    ├── basic.go
    ├── advanced.go
    ├── internal.go
    └── math_test.go
```

**Solución:**

```go
// math/basic.go
package math

// Add suma dos números
func Add(a, b float64) float64 {
    return a + b
}

// Subtract resta dos números
func Subtract(a, b float64) float64 {
    return a - b
}

// Multiply multiplica dos números
func Multiply(a, b float64) float64 {
    return a * b
}
```

```go
// math/advanced.go
package math

// Power calcula a elevado a la potencia b
func Power(a, b float64) float64 {
    result := 1.0
    for i := 0; i < int(b); i++ {
        result = Multiply(result, a)
    }
    return result
}

// Average calcula el promedio
func Average(numbers ...float64) float64 {
    if len(numbers) == 0 {
        return 0
    }
    sum := 0.0
    for _, n := range numbers {
        sum = Add(sum, n)
    }
    return Divide(sum, float64(len(numbers)))
}
```

```go
// math/internal.go
package math

// Divide es una función interna (privada)
func Divide(a, b float64) float64 {
    if b == 0 {
        return 0
    }
    return a / b
}

// Max es una función interna
func Max(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}
```

```go
// main.go
package main

import (
    "fmt"
    "ejercicio3/math"
)

func main() {
    fmt.Println("Add:", math.Add(5, 3))
    fmt.Println("Power:", math.Power(2, 8))
    fmt.Println("Average:", math.Average(10, 20, 30))
}
```

### Ejercicio 4: Importar con Alias y Resolver Conflictos

**Objetivo**: Trabajar con múltiples paquetes que tienen nombres similares, usando alias.

**Requisitos:**
- Crear paquetes con nombres potencialmente conflictivos
- Usar alias en imports para claridad
- Demostrar uso de múltiples paquetes similares

**Estructura:**
```
ejercicio4/
├── go.mod
├── main.go
├── parser/
│   └── parser.go
├── formatter/
│   └── formatter.go
└── serializer/
    └── serializer.go
```

**Solución:**

```go
// parser/parser.go
package parser

type Result struct {
    Data  string
    Error string
}

func Parse(input string) *Result {
    return &Result{Data: "Parsed: " + input}
}
```

```go
// formatter/formatter.go
package formatter

type Result struct {
    Formatted string
}

func Format(input string) *Result {
    return &Result{Formatted: "Formatted: " + input}
}
```

```go
// serializer/serializer.go
package serializer

func Serialize(data interface{}) string {
    return "Serialized"
}
```

```go
// main.go
package main

import (
    "fmt"
    p "ejercicio4/parser"      // alias
    f "ejercicio4/formatter"   // alias
    s "ejercicio4/serializer"  // alias
)

func main() {
    parsed := p.Parse("hello")
    fmt.Println(parsed.Data)
    
    formatted := f.Format("world")
    fmt.Println(formatted.Formatted)
    
    serialized := s.Serialize(parsed)
    fmt.Println(serialized)
}
```

### Ejercicio 5: Documentar Paquete con Ejemplos

**Objetivo**: Crear un paquete completamente documentado siguiendo estándares Go.

**Requisitos:**
- Package doc al inicio
- Documentación de cada función
- Archivo `example_test.go` con ejemplos
- Función `Example()` ejecutable

**Estructura:**
```
ejercicio5/
├── go.mod
├── main.go
└── greeting/
    ├── greeting.go
    └── example_test.go
```

**Solución:**

```go
// greeting/greeting.go
package greeting

import "strings"

// Package greeting proporciona funciones para generar saludos personalizados.
//
// Ejemplo básico:
//
//  msg := greeting.Hello("World")
//  fmt.Println(msg)  // Output: Hello, World!
package greeting

// Hello devuelve un saludo personalizado
func Hello(name string) string {
    return "Hello, " + name + "!"
}

// Greet crea un saludo formal
func Greet(firstName, lastName string) string {
    fullName := strings.Title(firstName) + " " + strings.Title(lastName)
    return "Greetings, " + fullName + "!"
}

// Goodbye devuelve un despido personalizado
func Goodbye(name string) string {
    return "Goodbye, " + name + "! See you soon."
}
```

```go
// greeting/example_test.go
package greeting

import (
    "fmt"
)

func ExampleHello() {
    msg := Hello("World")
    fmt.Println(msg)
    // Output: Hello, World!
}

func ExampleGreet() {
    msg := Greet("john", "doe")
    fmt.Println(msg)
    // Output: Greetings, John Doe!
}

func ExampleGoodbye() {
    msg := Goodbye("Alice")
    fmt.Println(msg)
    // Output: Goodbye, Alice! See you soon.
}
```

```go
// main.go
package main

import (
    "fmt"
    "ejercicio5/greeting"
)

func main() {
    fmt.Println(greeting.Hello("Go"))
    fmt.Println(greeting.Greet("jane", "smith"))
    fmt.Println(greeting.Goodbye("Bob"))
}
```

---

## RESUMEN Y PUNTOS CLAVE

### ✓ Conceptos Fundamentales

1. **Paquetes como Namespaces**: Organizan código, evitan conflictos
2. **Visibilidad Simple**: Primera letra (mayúscula/minúscula) determina todo
3. **No hay Imports Relativos**: Siempre rutas absolutas desde el módulo
4. **Encapsulación de Datos**: Usa funciones de acceso para datos privados
5. **Interfaces para Desacoplamiento**: Evita dependencias directas

### ✓ Mejores Prácticas

- ✓ Nombres descriptivos y minúsculos
- ✓ Documentación completa (godoc-ready)
- ✓ Cohesión alta en miembros del paquete
- ✓ Bajo acoplamiento entre paquetes
- ✓ Evitar ciclos de importación
- ✓ Usar `init()` para inicializaciones ordenadas
- ✓ Separar lógica de `main()`

### ✓ Antipatrones a Evitar

- ✗ Paquetes muy grandes ("paquete god")
- ✗ Nombres genéricos (util, helper, common)
- ✗ Alto acoplamiento entre paquetes
- ✗ Dependencias circulares
- ✗ Falta de documentación
- ✗ Mezcla de responsabilidades
- ✗ Exposición de detalles internos

---

## DIAGRAMA: Flujo de Importación y Visibilidad

```
Estructura del Proyecto:
┌─────────────────────────────────────────┐
│           github.com/user/app           │
├─────────────────────────────────────────┤
│ go.mod                                  │
│ cmd/                                    │
│ ├── main.go (package main)              │
│ pkg/                                    │
│ ├── config/                             │
│ │   └── config.go (package config)      │
│ ├── database/                           │
│ │   └── db.go (package database)        │
│ internal/                               │
│ └── models/                             │
│     └── models.go (package models)      │
└─────────────────────────────────────────┘

Importación:
┌──────────────────────────────────────────────────┐
│ En cmd/main.go:                                  │
│ import (                                         │
│   "github.com/user/app/pkg/config"   ✓ OK      │
│   "github.com/user/app/pkg/database" ✓ OK      │
│   "github.com/user/app/internal/..."✓ OK        │
│ )                                                │
└──────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────┐
│ En otro proyecto:                                │
│ import (                                         │
│   "github.com/user/app/pkg/config"   ✓ OK      │
│   "github.com/user/app/pkg/database" ✓ OK      │
│   "github.com/user/app/internal/..." ✗ BLOCKED │
│ )                                                │
└──────────────────────────────────────────────────┘

Visibilidad (dentro de un paquete):
┌────────────────────────────────────────┐
│ package mypackage                      │
├────────────────────────────────────────┤
│ Exportado (fuera del paquete):          │
│   const MAX = 100      ✓ Accesible     │
│   func Process() {}    ✓ Accesible     │
│   type User struct {}  ✓ Accesible     │
│                                        │
│ No Exportado (privado):                │
│   const min = 10       ✗ Privado      │
│   func helper() {}     ✗ Privado      │
│   type user struct {}  ✗ Privado      │
└────────────────────────────────────────┘
```

---

## REFERENCIAS Y RECURSOS

- [Go Packages Documentation](https://golang.org/doc/code#Workspaces)
- [Go Modules Reference](https://golang.org/ref/mod)
- [Standard Library](https://golang.org/pkg/)
- [Effective Go - Packages](https://golang.org/doc/effective_go#names)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

**Fin del Capítulo 20**

Ahora estás preparado para:
- Organizar tu código en paquetes eficientemente
- Manejar imports y visibilidad correctamente
- Documentar paquetes para uso público
- Evitar problemas comunes como ciclos de dependencia
- Crear una arquitectura escalable y mantenible


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/20-paquetes/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/20-paquetes):

```bash
cd examples/20-paquetes
go run .
```
