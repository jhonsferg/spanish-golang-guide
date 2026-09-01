# Capítulo 40: Build tools y performance

## Introducción

La compilación y optimización de programas Go es un arte que va más allá de ejecutar `go build`. En este capítulo exploraremos el ecosistema completo de herramientas de compilación, profiling y optimización que Go proporciona para construir sistemas de alto rendimiento.

Go destaca por su velocidad de compilación y ejecución, pero optimizar un programa Go requiere comprensión profunda del compilador, el garbage collector y las herramientas de profiling. Abordaremos desde conceptos básicos de compilación condicional hasta técnicas avanzadas de profiling con flame graphs.

### Casos de Uso

- **Microservicios**: Latencias predecibles, startup rápido
- **Sistemas embebidos**: Compilación cruzada a múltiples arquitecturas
- **Servicios de baja latencia**: CPU profiling y optimización de hot paths
- **Aplicaciones con restricciones de memoria**: Memory profiling y tuning del GC
- **Herramientas CLI**: Build tags para características condicionales

---

## 40.1 Sistema de Construcción de Go

### 40.1.1 Fundamentos de go build

`go build` es la herramienta principal para compilar programas Go. Compila packages y sus dependencias, generando un ejecutable.

```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Build System!")
}
```

**Compilación básica:**

```bash
# Compilar en el directorio actual
go build

# Compilar especificando archivo de salida
go build -o myapp

# Compilar con información de depuración
go build -gcflags="all=-N -l"

# Compilar con indicadores del linker
go build -ldflags="-s -w"
```

**Indicadores importantes:**

- `-o`: Especifica nombre del archivo de salida
- `-v`: Verboso, muestra packages compilados
- `-x`: Imprime comandos ejecutados
- `-n`: Imprime comandos sin ejecutarlos
- `-race`: Habilita race detector
- `-trimpath`: Elimina rutas del sistema de build info
- `-a`: Fuerza recompilación de todos packages

```bash
# Build verboso con salida de comandos
go build -x -v -o app

# Build con race detector para aplicaciones concurrentes
go build -race -o app

# Build limpio (sin rutas de sistema)
go build -trimpath -o app
```

### 40.1.2 go install vs go build

`go install` compila e instala el package en `$GOPATH/bin` (o `$GOBIN`).

```go
// tools.go - Herramienta CLI
package main

import (
    "flag"
    "fmt"
)

func main() {
    name := flag.String("name", "World", "nombre para saludar")
    flag.Parse()
    fmt.Printf("Hello, %s!\n", *name)
}
```

```bash
# Build: genera ejecutable en directorio actual
go build -o mytool
./mytool -name Go

# Install: compila e instala en $GOBIN
go install
mytool -name Go

# Verificar ubicación
which mytool
```

### 40.1.3 go clean

Limpia archivos de construcción y cachés.

```bash
# Limpiar ejecutables compilados
go clean

# Limpiar todo incluyendo cachés
go clean -cache

# Limpiar modules cache
go clean -modcache

# Limpiar y mostrar comandos
go clean -x
```

### 40.1.4 Gestión de cachés

Go mantiene un caché de compilación para acelerar builds posteriores.

```bash
# Ver información del caché
go env GOCACHE

# Mostrar tamaño del caché
du -sh $(go env GOCACHE)

# Limpiar caché
go clean -cache

# Configurar directorio de caché personalizado
export GOCACHE=/custom/cache/path
go build
```

---

## 40.2 Build Tags (Compilación Condicional)

### 40.2.1 Etiquetas Básicas

Los build tags permiten compilación condicional basada en plataforma, arquitectura u otros criterios.

```go
// db_sqlite.go
//go:build sqlite
// +build sqlite

package database

import "database/sql"

func Connect() *sql.DB {
    // Implementación con SQLite
    return nil
}
```

```go
// db_postgres.go
//go:build postgres
// +build postgres

package database

import "database/sql"

func Connect() *sql.DB {
    // Implementación con PostgreSQL
    return nil
}
```

**Compilación condicional:**

```bash
# Compilar con SQLite
go build -tags sqlite -o app_sqlite

# Compilar con PostgreSQL
go build -tags postgres -o app_postgres

# Compilar con múltiples tags
go build -tags "postgres,json" -o app
```

### 40.2.2 Tags Específicos de Plataforma

Go soporta tags automáticos para SO y arquitectura.

```go
// logger_unix.go
//go:build linux || darwin
// +build linux darwin

package logger

import "log/syslog"

func initLogger() {
    // Inicializar logger con syslog en Unix
}
```

```go
// logger_windows.go
//go:build windows
// +build windows

package logger

import "golang.org/x/sys/windows/eventlog"

func initLogger() {
    // Inicializar logger con event log en Windows
}
```

### 40.2.3 Combinación Compleja de Tags

```go
// crypto_fast.go
//go:build (amd64 || arm64) && (linux || darwin)
// +build amd64 arm64
// +build linux darwin

package crypto

func Hash(data []byte) []byte {
    // Implementación optimizada con CPU features
    return nil
}
```

```go
// crypto_portable.go
//go:build !(amd64 || arm64) || (windows)
// +build !amd64,!arm64 windows

package crypto

func Hash(data []byte) []byte {
    // Implementación portable
    return nil
}
```

### 40.2.4 Debug Tags

```go
// debug.go
//go:build debug
// +build debug

package main

import "fmt"

const DEBUG = true

func debugLog(msg string) {
    if DEBUG {
        fmt.Println("[DEBUG]", msg)
    }
}
```

```go
// release.go
//go:build !debug
// +build !debug

package main

const DEBUG = false

func debugLog(msg string) {
    // No-op en release
}
```

```bash
# Build de depuración
go build -tags debug -o app_debug

# Build de release
go build -tags "" -o app_release
```

---

## 40.3 Versionado y Información de Build

### 40.3.1 Inyección de Versión con ldflags

`ldflags` permite inyectar valores en tiempo de compilación.

```go
// main.go
package main

import (
    "flag"
    "fmt"
)

var (
    Version   = "dev"
    BuildDate = "unknown"
    GitCommit = "unknown"
)

func main() {
    version := flag.Bool("version", false, "mostrar versión")
    flag.Parse()

    if *version {
        fmt.Printf("Version: %s\n", Version)
        fmt.Printf("Build Date: %s\n", BuildDate)
        fmt.Printf("Git Commit: %s\n", GitCommit)
        return
    }

    fmt.Println("Running application...")
}
```

**Script de compilación:**

```bash
#!/bin/bash

VERSION=$(git describe --tags --always)
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD)

go build \
    -ldflags="-X main.Version=$VERSION \
              -X main.BuildDate=$BUILD_DATE \
              -X main.GitCommit=$GIT_COMMIT" \
    -o app

# Ejecutar
./app -version
```

### 40.3.2 Información de Build Detallada

Go proporciona información de build mediante `runtime/debug`.

```go
package main

import (
    "fmt"
    "runtime/debug"
)

func main() {
    info, _ := debug.ReadBuildInfo()

    fmt.Println("=== Build Info ===")
    fmt.Printf("Go Version: %s\n", info.GoVersion)

    for _, setting := range info.Settings {
        if setting.Key == "CGO_ENABLED" ||
           setting.Key == "GOOS" ||
           setting.Key == "GOARCH" {
            fmt.Printf("%s=%s\n", setting.Key, setting.Value)
        }
    }

    fmt.Println("\n=== Dependencies ===")
    for _, dep := range info.Deps {
        if dep.Replace != nil {
            fmt.Printf("%s -> %s (%s)\n",
                dep.Path, dep.Replace.Path, dep.Replace.Version)
        } else {
            fmt.Printf("%s %s\n", dep.Path, dep.Version)
        }
    }
}
```

### 40.3.3 Stripping y Optimización de Binarios

```bash
# Build estándar: ~10 MB
go build -o app_standard

# Build con información de depuración removida: ~5 MB
go build -ldflags="-s" -o app_stripped

# Build completamente optimizado: ~3 MB
go build -ldflags="-s -w" -o app_ultra

# Comparar tamaños
ls -lh app_*
```

**Explicación de ldflags:**

- `-s`: Strip symbol table
- `-w`: Strip DWARF debug info
- `-X`: Set package variable

---

## 40.4 Compilación Cruzada (Cross-Compilation)

### 40.4.1 Compilar para Múltiples Plataformas

```bash
#!/bin/bash
# build_all.sh - Script de compilación cruzada

APP_NAME="myapp"
VERSION="1.0.0"

# Definir targets
TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for TARGET in "${TARGETS[@]}"; do
    IFS='/' read -r OS ARCH <<< "$TARGET"

    OUTPUT="${APP_NAME}_${VERSION}_${OS}_${ARCH}"
    if [ "$OS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for $OS/$ARCH -> $OUTPUT"

    GOOS=$OS GOARCH=$ARCH go build \
        -ldflags="-X main.Version=$VERSION" \
        -o "$OUTPUT"

    if [ $? -eq 0 ]; then
        echo "✓ Success"
    else
        echo "✗ Failed"
    fi
done
```

### 40.4.2 Plataformas Soportadas

```bash
# Listar todos los GOOS/GOARCH soportados
go tool dist list

# Compilar para arquitecturas específicas
GOOS=linux GOARCH=amd64 go build    # Servidor x86-64
GOOS=linux GOARCH=arm64 go build    # Servidor ARM
GOOS=darwin GOARCH=arm64 go build   # MacBook M1/M2
GOOS=windows GOARCH=amd64 go build  # Windows x86-64
GOOS=linux GOARCH=riscv64 go build  # RISC-V
```

### 40.4.3 Compilación con CGO en Cross-Compilation

CGO complica la compilación cruzada. Usar `CGO_ENABLED=0` para builds puros de Go.

```bash
# Build puro de Go (sin CGO)
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o app.exe

# Build con CGO (requiere compilador C cruzado)
# En macOS, compilar para Linux con CGO:
# 1. Instalar osxcross o usar Docker
# 2. Configurar variables de compilación

# Alternativa: usar Docker
docker run --rm \
    -v $(pwd):/src \
    -w /src \
    golang:latest \
    sh -c 'GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -o app'
```

### 40.4.4 Verificar Binarios Compilados

```bash
# Ver información del ejecutable compilado
file app_linux_amd64
file app_darwin_arm64
file app_windows_amd64.exe

# Extraer información de versión
./app_linux_amd64 -version
./app_darwin_arm64 -version
./app_windows_amd64.exe -version

# Verificar tamaño en diferentes plataformas
for f in app_*; do
    echo "$f: $(du -h $f | cut -f1)"
done
```

---

## 40.5 Módulos y Gestión de Dependencias

### 40.5.1 go.mod y go.sum

El archivo `go.mod` declara module path y requisitos; `go.sum` contiene checksums.

```
// go.mod
module github.com/myuser/myapp

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/sirupsen/logrus v1.9.3
)

require (
    github.com/stretchr/testify v1.8.0 // indirect
)

exclude github.com/badlib/bad v1.0.0
```

**Gestión de versiones:**

```bash
# Obtener versión más reciente de un package
go get github.com/gorilla/mux

# Obtener versión específica
go get github.com/gorilla/mux@v1.8.0

# Obtener versión mínima
go get -u ./...

# Limpiar dependencias no usadas
go mod tidy

# Descargar dependencias
go mod download

# Verificar integridad
go mod verify

# Ver todas las versiones disponibles
go list -m -versions github.com/gorilla/mux
```

### 40.5.2 Vendoring

`go mod vendor` copia todas las dependencias al directorio `vendor/`.

```bash
# Crear directorio vendor
go mod vendor

# Compilar usando vendor
go build -mod=vendor -o app

# Verificar contenido de vendor
ls -la vendor/github.com/

# Usar vendor en CI/CD
git add vendor/
git commit -m "Add vendor dependencies"
```

### 40.5.3 Trabajo con Módulos Locales

```go
// go.mod
module github.com/myuser/myapp

go 1.21

require (
    github.com/myuser/mylib v1.0.0
)

// Durante desarrollo local:
replace github.com/myuser/mylib => ../mylib
```

```bash
# Trabajar con módulo local
go work init
go work use ./myapp ./mylib

# El workspace permite editar mylib y ver cambios en myapp
# sin publicar nuevas versiones

# Limpiar workspace
go work edit -dropuse=./mylib
```

---

## 40.6 Profiling de CPU

### 40.6.1 CPU Profiling Básico

El CPU profiling mide qué funciones consumen más tiempo de CPU.

```go
package main

import (
    "fmt"
    "math"
    "runtime/pprof"
    "os"
)

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func computeExpensive() {
    for i := 0; i < 1000000; i++ {
        _ = math.Sqrt(float64(i))
    }
}

func main() {
    // Iniciar profiling de CPU
    f, _ := os.Create("cpu.prof")
    defer f.Close()
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Código a perfilar
    for i := 0; i < 100; i++ {
        computeExpensive()
    }

    for i := 30; i <= 35; i++ {
        fmt.Printf("fib(%d) = %d\n", i, fibonacci(i))
    }

    fmt.Println("Profiling completado. Ver cpu.prof")
}
```

**Analizar profiling:**

```bash
# Compilar y ejecutar con profiling
go run main.go

# Ver top hotspots
go tool pprof -top cpu.prof

# Interfaz interactiva
go tool pprof cpu.prof
# Dentro de pprof:
# (pprof) top
# (pprof) list fibonacci
# (pprof) web

# Generar gráfico (requiere graphviz)
go tool pprof -http=:8080 cpu.prof
```

### 40.6.2 Testing con Profiling

```go
// perf_test.go
package main

import (
    "testing"
    "os"
    "runtime/pprof"
)

func TestCPUProfile(t *testing.T) {
    f, _ := os.Create("test_cpu.prof")
    defer f.Close()

    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Test a ejecutar
    for i := 0; i < 100000; i++ {
        _ = fibonacci(25)
    }
}

func BenchmarkFibonacci(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        fibonacci(20)
    }
}
```

```bash
# Ejecutar tests con profiling
go test -cpuprofile=cpu.prof -bench=BenchmarkFibonacci

# Analizar
go tool pprof -http=:8080 cpu.prof
```

### 40.6.3 Flame Graphs

Flame graphs visualizan call stacks durante profiling.

```bash
# Generar flame graph (requiere flamegraph-go)
go tool pprof -http=:8080 cpu.prof

# O usando go-torch (deprecated pero aún funcional)
# go-torch --url=http://localhost:6060 --time=30 --file=flame.svg

# Alternativa: graphviz
go tool pprof -svg cpu.prof > graph.svg
open graph.svg
```

### 40.6.4 On-Demand Profiling con net/http/pprof

Para aplicaciones servidoras, usar `net/http/pprof` para profiling on-demand.

```go
package main

import (
    "fmt"
    "log"
    "math"
    "net/http"
    _ "net/http/pprof"  // Importar para habilitar pprof
    "time"
)

func computeLoop() {
    for {
        for i := 0; i < 100000000; i++ {
            _ = math.Sqrt(float64(i))
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from handler\n")
}

func main() {
    go computeLoop()

    http.HandleFunc("/", handler)

    fmt.Println("Server running on :6060")
    fmt.Println("Profiling available at:")
    fmt.Println("  http://localhost:6060/debug/pprof")
    fmt.Println("  http://localhost:6060/debug/pprof/profile?seconds=30")

    log.Fatal(http.ListenAndServe(":6060", nil))
}
```

**Capturar profiling en vivo:**

```bash
# Iniciar servidor
go run main.go

# En otra terminal: capturar 30 segundos de CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Dentro de pprof:
# (pprof) top
# (pprof) web

# Ver goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Ver heap
go tool pprof http://localhost:6060/debug/pprof/heap
```

---

## 40.7 Profiling de Memoria

### 40.7.1 Heap Profiling

El heap profiling identifica dónde se asigna memoria.

```go
package main

import (
    "fmt"
    "os"
    "runtime"
    "runtime/pprof"
)

func allocateSlice(size int) []int {
    s := make([]int, size)
    for i := range s {
        s[i] = i
    }
    return s
}

func leakyFunction(iterations int) [][]int {
    result := make([][]int, iterations)
    for i := 0; i < iterations; i++ {
        result[i] = allocateSlice(1000)
    }
    return result
}

func main() {
    f, _ := os.Create("mem.prof")
    defer f.Close()

    // Coleccionar garbage antes de profiling
    runtime.GC()

    pprof.WriteHeapProfile(f)

    // Asignar mucha memoria
    data := leakyFunction(1000)
    fmt.Printf("Allocated %d slices\n", len(data))

    // Profiling después de allocaciones
    pprof.WriteHeapProfile(f)
}
```

**Análisis de heap:**

```bash
# Compilar y ejecutar
go run main.go

# Ver allocaciones top
go tool pprof -alloc_space mem.prof
# (pprof) top

# Ver allocaciones en uso (in_use)
go tool pprof -alloc_objects mem.prof
# (pprof) top

# Interfaz web
go tool pprof -http=:8080 mem.prof
```

### 40.7.2 Escape Analysis

El escape analysis determina qué variables escapan al heap (requieren allocation).

```bash
# Ver escape analysis en tiempo de compilación
go build -gcflags="-m" main.go

# Más verbosidad
go build -gcflags="-m=4" main.go
```

**Ejemplo de escape analysis:**

```go
// escape.go
package main

type Point struct {
    X, Y float64
}

// p escapa (retorna puntero)
func escapePointer() *Point {
    p := &Point{1, 2}
    return p
}

// p no escapa (retorna valor)
func noEscape() Point {
    p := Point{1, 2}
    return p
}

// slice escapa (pasa a fmt.Printf como interface{})
func escapeInterface() {
    nums := []int{1, 2, 3}
    println(nums)  // No escapa (println es built-in)
}
```

```bash
# Ver análisis de escape
go build -gcflags="-m" escape.go

# Output:
# ./escape.go:9:6: moved to heap: p
# ./escape.go:14:6: p does not escape
```

### 40.7.3 GC Tuning

Control del garbage collector para reducir pausas.

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "time"
)

func main() {
    // Ver configuración actual de GC
    fmt.Printf("GC Percentage: %d%%\n", debug.SetGCPercent(-1))

    // Restaurar y establecer valor personalizado
    debug.SetGCPercent(50)  // Ejecutar GC cuando heap crezca 50%

    // Ver métricas de GC
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("Alloc: %v MB\n", m.Alloc / 1024 / 1024)
    fmt.Printf("TotalAlloc: %v MB\n", m.TotalAlloc / 1024 / 1024)
    fmt.Printf("Sys: %v MB\n", m.Sys / 1024 / 1024)
    fmt.Printf("NumGC: %v\n", m.NumGC)
    fmt.Printf("PauseNs (last): %v ms\n", m.PauseNs[(m.NumGC+255)%256] / 1e6)

    // Forzar GC
    runtime.GC()

    time.Sleep(1 * time.Second)
}
```

**Variables de entorno para GC:**

```bash
# Desactivar GC (para profiling)
GOGC=off go run main.go

# Reducir frecuencia de GC (más throughput, pausas más largas)
GOGC=200 go run main.go

# Aumentar frecuencia de GC (menos latencia, menos throughput)
GOGC=25 go run main.go

# Ver trace de GC
GODEBUG=gctrace=1 go run main.go

# Output:
# gc 1 @0.521s 2%: 0.002+0.031+0.002 ms clock, 0.022+0.001+0.047 ms cpu, ...
```

### 40.7.4 Tracing

El tracing captura eventos de ejecución para análisis detallado.

```go
package main

import (
    "fmt"
    "os"
    "runtime/trace"
    "time"
)

func worker(id int) {
    for i := 0; i < 5; i++ {
        fmt.Printf("Worker %d: task %d\n", id, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Crear trace
    f, _ := os.Create("trace.out")
    defer f.Close()

    trace.Start(f)
    defer trace.Stop()

    // Código a trazar
    for i := 0; i < 3; i++ {
        go worker(i)
    }

    time.Sleep(2 * time.Second)
    fmt.Println("Done")
}
```

**Analizar trace:**

```bash
# Compilar y ejecutar
go run main.go

# Ver trace (requiere navegador)
go tool trace trace.out

# Se abre interfaz web con:
# - Timeline de goroutines
# - Eventos de GC
# - Eventos de red
# - Profiling integrado
```

---

## 40.8 Benchmarking

### 40.8.1 Escritura de Benchmarks

```go
// math_test.go
package math

import (
    "fmt"
    "testing"
)

func sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

func BenchmarkSum(b *testing.B) {
    numbers := make([]int, 1000)
    for i := range numbers {
        numbers[i] = i
    }

    b.ResetTimer()  // Excluir setup del timing

    for i := 0; i < b.N; i++ {
        sum(numbers)
    }
}

func BenchmarkSumParallel(b *testing.B) {
    numbers := make([]int, 1000)
    for i := range numbers {
        numbers[i] = i
    }

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            sum(numbers)
        }
    })
}

func BenchmarkSumVariableSizes(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
            numbers := make([]int, size)
            for i := range numbers {
                numbers[i] = i
            }

            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                sum(numbers)
            }
        })
    }
}
```

### 40.8.2 Ejecutar Benchmarks

```bash
# Ejecutar todos los benchmarks
go test -bench=. -benchmem

# Output:
# BenchmarkSum-8                  500000000      2 ns/op    0 B/op    0 allocs/op
# BenchmarkSumParallel-8         2000000000      0.5 ns/op  0 B/op    0 allocs/op

# Especificar benchmark
go test -bench=BenchmarkSum -benchmem

# Ejecutar con tiempo mínimo
go test -bench=. -benchtime=3s

# Guardar resultados baseline
go test -bench=. -benchmem > old.txt

# Comparar resultados
go test -bench=. -benchmem > new.txt
benchstat old.txt new.txt

# Verbose: mostrar iteraciones
go test -bench=. -v
```

### 40.8.3 Benchmarks con Allocations

```go
func BenchmarkAllocation(b *testing.B) {
    b.ReportAllocs()  // Reportar allocations

    for i := 0; i < b.N; i++ {
        _ = make([]int, 100)  // Allocation no necesaria
    }
}

func BenchmarkNoAllocation(b *testing.B) {
    b.ReportAllocs()

    slice := make([]int, 100)

    for i := 0; i < b.N; i++ {
        slice = slice[:0]  // Reutilizar slice
    }
}
```

```bash
# Ejecutar con reporte de allocations
go test -bench=. -benchmem

# Output:
# BenchmarkAllocation-8         100000000   10 ns/op  800 B/op  1 allocs/op
# BenchmarkNoAllocation-8       1000000000   1 ns/op    0 B/op  0 allocs/op
```

### 40.8.4 Detección de Regresión

```bash
# Herramienta benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Ejecutar benchmarks y guardar
go test -bench=. -benchmem > benchmarks.txt

# Comparar con baseline (después de cambios)
go test -bench=. -benchmem > benchmarks_new.txt
benchstat benchmarks.txt benchmarks_new.txt

# Output:
# name                    old time/op    new time/op    delta
# Sum/size_1000-8         2.00µs ± 2%    2.01µs ± 3%   ~     (p=0.432, n=5+5)
# Sum/size_10000-8        20.0µs ± 1%    20.1µs ± 2%   ~     (p=0.356, n=5+5)
```

---

## 40.9 Técnicas de Optimización

### 40.9.1 Inline y Compiler Optimizations

El compilador de Go puede inline funciones pequeñas.

```go
// sin inline explícito (compilador decide)
func add(a, b int) int {
    return a + b
}

// -gcflags="-l" desactiva inlining
go build -gcflags="-l" main.go

// Ver decisiones de inlining
go build -gcflags="-m" main.go
```

**Función que sí se inlines:**

```go
func fastAdd(a, b int) int {
    return a + b
}

func expensiveAdd(a, b int) int {
    // Función más grande, probablemente no se inline
    temp := a + 1
    temp += b + 1
    temp--
    return temp - 1
}

// En benchmark:
func BenchmarkInline(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fastAdd(1, 2)
    }
}
```

### 40.9.2 Evitar Allocations Innecesarias

```go
// ❌ Crea allocation en cada iteración
func inefficient() {
    for i := 0; i < 1000000; i++ {
        s := make([]byte, 1024)
        process(s)
    }
}

// ✓ Reutiliza buffer
func efficient() {
    buf := make([]byte, 1024)
    for i := 0; i < 1000000; i++ {
        buf = buf[:0]
        process(buf)
    }
}

// ✓ Usar pool de buffers
func withPool() {
    pool := sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)
        },
    }

    for i := 0; i < 1000000; i++ {
        buf := pool.Get().([]byte)
        process(buf)
        pool.Put(buf)
    }
}
```

### 40.9.3 String Concatenation

```go
// ❌ Lento: crea string en cada concatenación
func slowConcat(parts []string) string {
    result := ""
    for _, p := range parts {
        result = result + p
    }
    return result
}

// ✓ Rápido: usa strings.Builder
func fastConcat(parts []string) string {
    var b strings.Builder
    for _, p := range parts {
        b.WriteString(p)
    }
    return b.String()
}

func BenchmarkConcat(b *testing.B) {
    parts := []string{"hello", "world", "go", "lang"}
    b.Run("slow", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            slowConcat(parts)
        }
    })

    b.Run("fast", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            fastConcat(parts)
        }
    })
}
```

### 40.9.4 Interface{} vs Tipos Genéricos

```go
// Go 1.17: interface{} requiere runtime type switching
func processInterface(v interface{}) {
    switch v.(type) {
    case int:
        // ...
    case string:
        // ...
    }
}

// Go 1.18+: Generics es más eficiente
func processGeneric[T int | string](v T) {
    // Sin runtime overhead
}

func BenchmarkPolymorphism(b *testing.B) {
    b.Run("interface", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            processInterface(42)
        }
    })

    b.Run("generic", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            processGeneric(42)
        }
    })
}
```

### 40.9.5 Optimization Patterns

```go
// ✓ Batch processing para mejor cache locality
func batchProcess(data []int) int {
    total := 0
    const batchSize = 64

    for i := 0; i < len(data); i += batchSize {
        end := i + batchSize
        if end > len(data) {
            end = len(data)
        }

        for j := i; j < end; j++ {
            total += data[j]
        }
    }
    return total
}

// ✓ Reduce memory fragmentation
func objectPool[T any](size int, factory func() T) []T {
    return make([]T, size)
}

// ✓ Avoid lock contention
type Counter struct {
    mu    sync.Mutex
    count int
}

// ❌ Todos los goroutines compiten por el mismo lock
func (c *Counter) Increment() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}

// ✓ Usar sync/atomic para operaciones simples
type AtomicCounter struct {
    count atomic.Int64
}

func (c *AtomicCounter) Increment() {
    c.count.Add(1)
}
```

---

## 40.10 Debugging

### 40.10.1 Delve Debugger

Delve es el debugger principal para Go.

```bash
# Instalar Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Compilar con debug info
go build -gcflags="all=-N -l" -o app main.go

# Iniciar Delve
dlv exec ./app

# Comandos básicos en Delve:
# (dlv) break main.main
# (dlv) continue
# (dlv) next
# (dlv) step
# (dlv) print variable
# (dlv) whatis variable
```

### 40.10.2 Debugging con Breakpoints

```go
// main.go
package main

import "fmt"

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    for i := 0; i < 10; i++ {
        result := fibonacci(i)
        fmt.Printf("fib(%d) = %d\n", i, result)
    }
}
```

```bash
# Compilar sin optimizaciones
go build -gcflags="all=-N -l" -o app main.go

# Depurar con dlv
dlv debug ./app

# Dentro de dlv:
# (dlv) break main.fibonacci
# (dlv) condition 1 n == 5
# (dlv) continue
# (dlv) print n
# (dlv) goroutines
```

### 40.10.3 Debugging Concurrente

```go
// goroutines.go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < 5; i++ {
        fmt.Printf("Worker %d: %d\n", id, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
}
```

```bash
# Depurar goroutines
dlv debug ./goroutines.go

# Dentro de dlv:
# (dlv) goroutines           # Listar todas las goroutines
# (dlv) goroutine 1          # Cambiar a goroutine 1
# (dlv) goroutine 2 print id # Ejecutar comando en goroutine
# (dlv) break worker         # Breakpoint en función worker
```

### 40.10.4 Print Debugging vs Structured Logging

```go
// ❌ Print debugging (baja calidad)
fmt.Println("Debug: value =", v)

// ✓ Structured logging (mejor)
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("processing item",
    slog.Int("item_id", itemID),
    slog.Duration("duration", elapsed),
)
```

---

## 40.11 Buenas Prácticas y Patterns

### 40.11.1 Medir Antes de Optimizar

La regla de oro: **no optimices sin datos**.

```go
import (
    "fmt"
    "log"
    "time"
)

type Profiler struct {
    name  string
    start time.Time
}

func NewProfiler(name string) *Profiler {
    return &Profiler{
        name:  name,
        start: time.Now(),
    }
}

func (p *Profiler) Stop() {
    elapsed := time.Since(p.start)
    log.Printf("[%s] took %v\n", p.name, elapsed)
}

func processData(data []int) {
    defer NewProfiler("processData").Stop()

    // Procesamiento...
    time.Sleep(100 * time.Millisecond)
}

// Output:
// [processData] took 100.123456ms
```

### 40.11.2 Build Pipeline en CI/CD

```yaml
# .github/workflows/build.yml
name: Build and Test

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build
        run: go build -v -o app

      - name: Run Tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Run Benchmarks
        run: go test -bench=. -benchmem > benchmarks.txt

      - name: Upload Benchmarks
        uses: actions/upload-artifact@v3
        with:
          name: benchmarks
          path: benchmarks.txt

      - name: Cross-Compile
        run: |
          GOOS=linux GOARCH=amd64 go build -o app_linux_amd64
          GOOS=darwin GOARCH=arm64 go build -o app_darwin_arm64
          GOOS=windows GOARCH=amd64 go build -o app_windows_amd64.exe

      - name: Upload Artifacts
        uses: actions/upload-artifact@v3
        with:
          name: binaries
          path: app_*
```

### 40.11.3 Performance Testing en CI/CD

```bash
#!/bin/bash
# scripts/benchmark.sh

# Ejecutar benchmarks
go test -bench=. -benchmem -run=^$ > current.txt

# Comparar con baseline
if [ -f baseline.txt ]; then
    benchstat baseline.txt current.txt

    # Fallar si hay regresión > 5%
    benchstat -threshold 0.05 baseline.txt current.txt
fi

# Guardar como nuevo baseline
cp current.txt baseline.txt
```

### 40.11.4 Profiling en Producción

Para aplicaciones en producción, usar configuración segura de profiling:

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
)

func main() {
    // Servidor principal
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    // Profiling en puerto separado con autenticación
    go func() {
        log.Fatal(http.ListenAndServe(":6060", mux))
    }()

    // ...
}

// En producción: proteger con firewall y autenticación
// iptables -I INPUT -p tcp --dport 6060 -s 10.0.0.0/8 -j ACCEPT
// iptables -I INPUT -p tcp --dport 6060 -j DROP
```

### 40.11.5 Comparación: Go vs C++ vs Rust

| Aspecto | Go | C++ | Rust |
|---------|----|----|------|
| **Compilación** | Muy rápida | Media | Lenta |
| **Ejecución** | Rápida | Muy rápida | Muy rápida |
| **GC Latency** | 1-100ms | N/A | N/A |
| **Memory Safety** | Parcial | No | Sí |
| **Profiling** | Excelente | Bueno | Bueno |
| **Curva aprendizaje** | Suave | Pronunciada | Pronunciada |

**Cuándo usar Go:**

- Microservicios, APIs, herramientas CLI
- Latencia <10ms aceptable
- Desarrollo rápido priorizado

**Cuándo usar C++:**

- Sistemas embebidos, gráficos
- Latencia < 1ms requerida
- Control total de memoria

**Cuándo usar Rust:**

- Sistemas críticos, kernels, firmware
- Memory safety crítica
- Performance comparable a C++

### 40.11.6 Antipatrones Comunes

```go
// ❌ Premature optimization
func prematureOptimization() {
    // Micro-optimizaciones sin benchmarks
    // Resultado: código complejo, poco beneficio
}

// ✓ Data-driven optimization
func datadriven() {
    // 1. Medir (benchmark, profile)
    // 2. Identificar bottleneck (top 20%)
    // 3. Optimizar (simple primero)
    // 4. Verificar ganancia
}

// ❌ Ignorar GC
var globalSlice []int

func appendLoop() {
    for i := 0; i < 1000000; i++ {
        globalSlice = append(globalSlice, i)
    }
}

// ✓ Pre-allocar
func efficientAppend() {
    slice := make([]int, 0, 1000000)
    for i := 0; i < 1000000; i++ {
        slice = append(slice, i)
    }
}

// ❌ Micro-optimizaciones que no importan
for i := 0; i < len(arr); i++ {
    // Premature optimization
}

// ✓ Enfoque pragmático
for _, item := range arr {
    process(item)
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Build Tags para Código Condicional

Crear una aplicación con diferentes backends de base de datos usando build tags.

```go
// go-build-tags/main.go
package main

import (
    "fmt"
    "log"
)

func main() {
    fmt.Println("=== Database Application ===")
    db, err := Connect()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Connected to: %s\n", db.Name())

    // Ejecutar query de prueba
    result, err := db.Query("SELECT 1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Query result: %v\n", result)
}
```

```go
// go-build-tags/db_sqlite.go
//go:build sqlite
// +build sqlite

package main

import "fmt"

type Database struct {
    driver string
}

func (d *Database) Name() string {
    return d.driver
}

func (d *Database) Query(sql string) (interface{}, error) {
    fmt.Printf("[SQLite] Executing: %s\n", sql)
    return 1, nil
}

func Connect() (*Database, error) {
    return &Database{driver: "SQLite"}, nil
}
```

```go
// go-build-tags/db_postgres.go
//go:build postgres
// +build postgres

package main

import "fmt"

type Database struct {
    driver string
}

func (d *Database) Name() string {
    return d.driver
}

func (d *Database) Query(sql string) (interface{}, error) {
    fmt.Printf("[PostgreSQL] Executing: %s\n", sql)
    return "PostgreSQL result", nil
}

func Connect() (*Database, error) {
    return &Database{driver: "PostgreSQL"}, nil
}
```

```bash
# Ejecutar
cd go-build-tags

# Con SQLite
go run -tags sqlite .
# Output: Connected to: SQLite

# Con PostgreSQL
go run -tags postgres .
# Output: Connected to: PostgreSQL

# Compilar para ambas plataformas
go build -tags sqlite -o app_sqlite .
go build -tags postgres -o app_postgres .
```

### Ejercicio 2: Cross-Compilation para Múltiples Plataformas

Crear script que compile para Windows, macOS y Linux.

```bash
#!/bin/bash
# go-cross-compile/build.sh

#!/bin/bash
APP_NAME="todoapp"
VERSION="1.0.0"

# Limpiar builds anteriores
rm -f ${APP_NAME}_*

# Targets: OS/ARCH
TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for TARGET in "${TARGETS[@]}"; do
    IFS='/' read -r OS ARCH <<< "$TARGET"

    OUTPUT="${APP_NAME}_${VERSION}_${OS}_${ARCH}"
    [ "$OS" = "windows" ] && OUTPUT="${OUTPUT}.exe"

    echo "Building $TARGET -> $OUTPUT"

    GOOS=$OS GOARCH=$ARCH go build \
        -ldflags="-X main.Version=$VERSION -X main.OS=$OS -X main.ARCH=$ARCH -s -w" \
        -o "$OUTPUT" \
        . || exit 1

    # Mostrar info del binario
    file "$OUTPUT"
done

# Resumir builds
echo ""
echo "Build completed:"
ls -lh ${APP_NAME}_*
du -sh ${APP_NAME}_* | awk '{print $2": "$1}'
```

```go
// go-cross-compile/main.go
package main

import (
    "flag"
    "fmt"
    "runtime"
)

var (
    Version = "dev"
    OS      = "unknown"
    ARCH    = "unknown"
)

func main() {
    version := flag.Bool("version", false, "show version")
    flag.Parse()

    if *version {
        fmt.Printf("%s v%s (%s/%s)\n",
            "todoapp", Version,
            runtime.GOOS, runtime.GOARCH)
        return
    }

    fmt.Printf("Todo App running on %s/%s\n",
        runtime.GOOS, runtime.GOARCH)
}
```

```bash
cd go-cross-compile
chmod +x build.sh
./build.sh

# Verificar binarios
./todoapp_1.0.0_linux_amd64 -version
./todoapp_1.0.0_darwin_arm64 -version
# etc.
```

### Ejercicio 3: CPU Profiling y Optimización

Identificar y optimizar una función lenta.

```go
// go-cpu-profile/slow.go
package main

import (
    "fmt"
    "math"
    "time"
)

func slowFibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return slowFibonacci(n-1) + slowFibonacci(n-2)
}

func fastFibonacci(n int) int {
    if n <= 1 {
        return n
    }

    prev, curr := 0, 1
    for i := 2; i <= n; i++ {
        prev, curr = curr, prev+curr
    }
    return curr
}

func computeIntensive() float64 {
    result := 0.0
    for i := 0; i < 10000000; i++ {
        result += math.Sqrt(float64(i))
    }
    return result
}

func slowAlgorithm() {
    fmt.Println("Starting CPU profiling...")
    start := time.Now()

    // Fibonacci recursivo (lento)
    for i := 30; i <= 35; i++ {
        _ = slowFibonacci(i)
    }

    // Computación intensiva
    _ = computeIntensive()

    elapsed := time.Since(start)
    fmt.Printf("Slow algorithm took: %v\n", elapsed)
}

func fastAlgorithm() {
    fmt.Println("Starting optimized version...")
    start := time.Now()

    // Fibonacci iterativo (rápido)
    for i := 30; i <= 35; i++ {
        _ = fastFibonacci(i)
    }

    // Misma computación
    _ = computeIntensive()

    elapsed := time.Since(start)
    fmt.Printf("Fast algorithm took: %v\n", elapsed)
}
```

```go
// go-cpu-profile/main.go
package main

import (
    "flag"
    "os"
    "runtime/pprof"
)

func main() {
    optimize := flag.Bool("optimize", false, "use optimized version")
    profile := flag.Bool("profile", false, "enable CPU profiling")
    flag.Parse()

    if *profile {
        f, _ := os.Create("cpu.prof")
        defer f.Close()
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    if *optimize {
        fastAlgorithm()
    } else {
        slowAlgorithm()
    }
}
```

```bash
cd go-cpu-profile

# Ejecutar versión lenta
go run main.go slow.go
# Output: Slow algorithm took: ~2 seconds

# Con profiling
go run main.go slow.go -profile
go tool pprof -top cpu.prof

# Ejecutar versión optimizada
go run main.go slow.go -optimize
# Output: Fast algorithm took: ~50ms

# Comparación: 40x más rápido
```

### Ejercicio 4: Memory Profiling y Optimización

Detectar y reducir allocations innecesarias.

```go
// go-mem-profile/leaky.go
package main

import (
    "fmt"
    "os"
    "runtime"
    "runtime/pprof"
    "strings"
)

// ❌ Versión con muchas allocations
func leakyStringConcat(count int) string {
    result := ""
    for i := 0; i < count; i++ {
        result = result + fmt.Sprintf("Item %d\n", i)
    }
    return result
}

// ✓ Versión optimizada
func efficientStringConcat(count int) string {
    var b strings.Builder
    for i := 0; i < count; i++ {
        fmt.Fprintf(&b, "Item %d\n", i)
    }
    return b.String()
}

// ❌ Slice append sin pre-allocar
func leakySliceGrow(count int) []int {
    var nums []int
    for i := 0; i < count; i++ {
        nums = append(nums, i)
    }
    return nums
}

// ✓ Slice con pre-allocación
func efficientSliceGrow(count int) []int {
    nums := make([]int, 0, count)
    for i := 0; i < count; i++ {
        nums = append(nums, i)
    }
    return nums
}

func leaky(count int) {
    fmt.Println("=== Leaky (Many Allocations) ===")

    runtime.GC()
    pprof.WriteHeapProfile(os.Create("leaky_start.prof"))

    _ = leakyStringConcat(count)
    _ = leakySliceGrow(count)

    runtime.GC()
    pprof.WriteHeapProfile(os.Create("leaky_end.prof"))
}

func optimized(count int) {
    fmt.Println("=== Optimized (Fewer Allocations) ===")

    runtime.GC()
    pprof.WriteHeapProfile(os.Create("opt_start.prof"))

    _ = efficientStringConcat(count)
    _ = efficientSliceGrow(count)

    runtime.GC()
    pprof.WriteHeapProfile(os.Create("opt_end.prof"))
}
```

```go
// go-mem-profile/main.go
package main

import (
    "flag"
)

func main() {
    optimize := flag.Bool("optimize", false, "use optimized version")
    count := flag.Int("count", 10000, "number of iterations")
    flag.Parse()

    if *optimize {
        optimized(*count)
    } else {
        leaky(*count)
    }
}
```

```bash
cd go-mem-profile

# Ejecutar versión leaky
go run main.go -count 100000

# Analizar heap
go tool pprof leaky_end.prof
# (pprof) top -cum

# Ejecutar versión optimizada
go run main.go -optimize -count 100000

# Comparar
go tool pprof -http=:8080 leaky_end.prof
go tool pprof -http=:8080 opt_end.prof
```

### Ejercicio 5: Pipeline Completo (Build, Test, Profile, Optimize)

Crear una aplicación con pipeline de benchmarking y optimización.

```go
// go-full-pipeline/handler.go
package main

import (
    "fmt"
    "time"
)

type RequestHandler struct {
    name string
}

// ❌ Versión lenta
func (h *RequestHandler) ProcessSlow(data []int) int {
    result := 0
    for _, v := range data {
        for i := 0; i < 1000; i++ {
            result += v
        }
    }
    return result
}

// ✓ Versión optimizada
func (h *RequestHandler) ProcessFast(data []int) int {
    result := 0
    sum := 0
    for _, v := range data {
        sum += v
    }
    result = sum * 1000
    return result
}

func (h *RequestHandler) Process(data []int) int {
    return h.ProcessFast(data)
}
```

```go
// go-full-pipeline/handler_test.go
package main

import (
    "testing"
)

func BenchmarkProcessSlow(b *testing.B) {
    h := &RequestHandler{name: "slow"}
    data := make([]int, 100)
    for i := range data {
        data[i] = i
    }

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        h.ProcessSlow(data)
    }
}

func BenchmarkProcessFast(b *testing.B) {
    h := &RequestHandler{name: "fast"}
    data := make([]int, 100)
    for i := range data {
        data[i] = i
    }

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        h.ProcessFast(data)
    }
}

func TestProcess(t *testing.T) {
    h := &RequestHandler{name: "test"}
    data := []int{1, 2, 3, 4, 5}

    result := h.Process(data)
    expected := 15 * 1000

    if result != expected {
        t.Errorf("Expected %d, got %d", expected, result)
    }
}
```

```bash
#!/bin/bash
# go-full-pipeline/pipeline.sh

cd go-full-pipeline

echo "=== Step 1: Build ==="
go build -o app -ldflags="-s -w"
echo "✓ Build successful"

echo ""
echo "=== Step 2: Test ==="
go test -v
echo "✓ Tests passed"

echo ""
echo "=== Step 3: Benchmark (Baseline) ==="
go test -bench=. -benchmem > baseline.txt
cat baseline.txt

echo ""
echo "=== Step 4: Profile & Optimize ==="
# Simular "optimización"
# (en real, se editaría el código)

echo ""
echo "=== Step 5: Benchmark (After Optimization) ==="
go test -bench=. -benchmem > optimized.txt
cat optimized.txt

echo ""
echo "=== Step 6: Compare Results ==="
if command -v benchstat &> /dev/null; then
    benchstat baseline.txt optimized.txt
else
    echo "benchstat not installed. Install with:"
    echo "go install golang.org/x/perf/cmd/benchstat@latest"
fi

echo ""
echo "✓ Pipeline complete"
```

```bash
cd go-full-pipeline
chmod +x pipeline.sh
./pipeline.sh

# Output esperado:
# === Step 1: Build ===
# ✓ Build successful
# === Step 2: Test ===
# --- PASS: TestProcess (0.00s)
# ✓ Tests passed
# === Step 3: Benchmark (Baseline) ===
# BenchmarkProcessSlow-8  1000 1000000 ns/op  0 B/op  0 allocs/op
# BenchmarkProcessFast-8  100000 10000 ns/op  0 B/op  0 allocs/op
```

---

## Conclusión

En este capítulo hemos cubierto el ecosistema completo de herramientas de build y performance en Go:

1. **Build System**: Entender `go build`, `go install`, y gestión de cachés
2. **Build Tags**: Compilación condicional para múltiples plataformas
3. **Versionado**: Inyectar información en tiempo de compilación
4. **Cross-Compilation**: Compilar para diferentes arquitecturas
5. **Modules**: Gestión moderna de dependencias
6. **CPU Profiling**: Identificar bottlenecks de rendimiento
7. **Memory Profiling**: Detectar y reducir allocations
8. **Benchmarking**: Medir rendimiento de forma reproducible
9. **Optimización**: Técnicas prácticas de performance
10. **Debugging**: Herramientas para diagnóstico
11. **Prácticas**: Data-driven optimization y CI/CD

### Puntos Clave

- **Mide primero**: No optimices sin datos de profiling
- **Automatiza**: CI/CD debe incluir builds cross-platform y benchmarks
- **Itera**: Build → Test → Profile → Optimize → Repeat
- **Monitorea**: En producción, captura profiling data de forma segura
- **Compara**: Usa benchstat para detectar regresiones

Go ofrece herramientas excepcionales para performance tunning. Con disciplina en medición y optimización data-driven, construirás sistemas Go de alto rendimiento.

---

## Referencias y Recursos Adicionales

- [Go Build Documentation](https://golang.org/cmd/go/)
- [Go Code Review Comments - Performance](https://github.com/golang/go/wiki/CodeReviewComments)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Delve Debugger](https://github.com/go-delve/delve)
- [Go Performance Book](https://golang.org/doc/diagnostics)
- [benchstat tool](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/40-build-tools-y-performance/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/40-build-tools-y-performance):

```bash
cd examples/40-build-tools-y-performance
go test -v -bench=. .
```
