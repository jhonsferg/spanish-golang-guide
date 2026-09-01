# Capítulo 2: Instalación y configuración

## Índice del Capítulo 2

1. [2.1 Historia de Versiones de Go](#21-historia-de-versiones-de-go)
2. [2.2 Arquitectura del Go SDK](#22-arquitectura-del-go-sdk)
3. [2.3 Instalación por Sistema Operativo](#23-instalación-por-sistema-operativo)
4. [2.4 GOPATH vs Go Modules - Evolución](#24-gopath-vs-go-modules--evolución)
5. [2.5 Go Modules Detallado](#25-go-modules-detallado)
6. [2.6 Comandos go - Referencia Completa](#26-comandos-go--referencia-completa)
7. [2.7 Workspace de Go (Multi-módulo)](#27-workspace-de-go-multi-módulo)
8. [2.8 Configuración de Entorno](#28-configuración-de-entorno)
9. [2.9 IDEs y Herramientas](#29-ides-y-herramientas)
10. [2.10 Troubleshooting Común](#210-troubleshooting-común)

---

## 2.1 Historia de Versiones de Go

### Por Qué Importa la Historia de Versiones

Antes de instalar Go, es crucial entender DÓNDE viene y QUÉ ha cambiado. Go no siempre fue así. Las decisiones de versión reflejan la evolución del lenguaje basada en problemas reales.

### Era 1: 2009-2011 - Desarrollo Experimental (Pre-1.0)

**Contexto:**
Go fue anunciado públicamente en 2009, pero no fue una "1.0" formal hasta 2012. Los años 2009-2011 fueron experimentales.

**Características de Go pre-1.0:**

```
 Goroutines: primitivas
 Channels: versión temprana
 GC: No estaba completo
 Interfaces: No implícitas
 Mdddulos: NO existían (solo GOPATH)
 Standard library: Pequeña
 Compilador: Escrito en C, lento
```

**Problema:**
Los usuarios reportaban crashes, comportamiento impredecible, incompatibilidades. Google internamente usaba Go, pero no era estable para producción.

**Razón histórica:**
Go fue una apuesta. Google no sabía si funcionaría. Necesitaba validación en el mundo real antes de comprometerse a compatibilidad.

### Era 2: 2012-2014 - Go 1.0 a Go 1.3 (Lanzamiento Oficial)

**Go 1.0 (Marzo 2012)**

```
HITO CRÍTICO: Go1 Compatibility Promise


 A partir de Go 1.0, se garantiza que:       │
                                              │
 ✅ Código escrito en Go 1.0 funcionará en   │
    Go 1.1, 1.2, 1.3, etc.                  │
                                              │
 ✅ Cambios solo en versiones mayores        │
    (Go 2.0, 3.0) romperían compatibilidad  │
                                              │
 ✅ Garantía de 5+ años de soporte           │
                                              │

```

**Características de Go 1.0:**

```
 GC completo y funcional
 defer, panic, recover estables
 goroutines optimizadas
 channels tipados seguros
 Stdlib expandida
 Compilador: aún en C (lento)
 Documentación completa
 Promesa de compatibilidad
```

**Impacto:**
Desarrolladores podían CONFIAR en Go. Si escribías código en 1.0, no se rompería en 1.5.

**Go 1.1 - Go 1.3 (2013-2014):**

```
 Mejoras de GC (pausas más cortas)
 Optimizaciones del compilador
 Stack preemption (mejor multitarea)
 race detector (-race flag)
 Mejor soporte de plataformas
```

### Era 3: 2014-2018 - Maduración Gradual (Go 1.4 a Go 1.11)

**Go 1.4 (Diciembre 2014) - Reescritura del Compilador**

```
DECISIÓN CRUCIAL: Compilador reescrito en Go


 Antes (Go 1.0-1.3):                     │
 ├─ Compilador en C                      │
 ├─ Difícil de modificar                 │
 └─ Lento en desarrollo                  │
                                          │
 Después (Go 1.4+):                      │
 ├─ Compilador en Go                     │
 ├─ Fácil de modificar                   │
 ├─ Desarrolladores pueden contribuir    │
 └─ Bootstrap problemático resuelto      │

```

**Razón filosófica:**
Si Go es lo suficientemente bueno para tus apps, debería ser lo suficientemente bueno para Go mismo. Auto-bootstrapping.

**Go 1.5-1.11 (2015-2018):**

```
 Go 1.5: Compilador escrito en Go
 Go 1.6: HTTP/2 en stdlib
 Go 1.7: Context package (estándar)
 Go 1.8: Plugin support
 Go 1.9: Type aliases
 Go 1.10: Mejor cache de builds
 Go 1.11: Módulos beta, WebAssembly
```

### Era 4: 2018-2021 - Revolución de Módulos (Go 1.11 a Go 1.16)

**El Problema Pre-Módulos: GOPATH**

```
ANTES (GOPATH - 2009-2019):

$GOPATH/
 bin/              (ejecutables)
 pkg/              (compilados)
 src/
    ├── github.com/usuario/proyecto1/
    │   └── main.go
    ├── github.com/usuario/proyecto2/
    │   └── main.go
    └── golang.org/x/sys/

Problema:
 UN solo GOPATH para todos los proyectos
 Todos los proyectos usan ÚLTIMA versión de deps
 Imposible reproducibilidad
 Proyecto A necesita versión 1.0, Proyecto B necesita 2.0
   └─ ¡Conflicto sin resolver!
 Descargaba TODAS las versiones en mismo lugar
```

**Go 1.11 (Agosto 2018) - Módulos Beta**

```
REVOLUCIÓN: Go Modules


 Proyecto actual (no GOPATH)              │
                                           │
 proyecto/                                 │
 ├── go.mod          (qué necesitas)      │
 ├── go.sum          (qué descargar)      │
 ├── main.go                              │
 └── internal/       (código privado)     │
                                           │
 Beneficios:                              │
 ✅ Versionamiento semántico (v1.2.3)    │
 ✅ Reproducibilidad exacta                │
 ✅ Sin GOPATH requerido                   │
 ✅ Per-proyecto, no global                │
 ✅ Compatible con SemVer                  │

```

**Go 1.16 (Febrero 2021) - Módulos por Defecto**

```
DECISIÓN FINAL: GOPATH fue deprecado


 GO111MODULE comportamiento:             │
                                         │
 Go 1.11: GO111MODULE=off (GOPATH)     │
 Go 1.13: GO111MODULE=on (Módulos)     │
 Go 1.16: GOPATH + Módulos (ambos)     │
 Go 1.17+: Solo Módulos                 │
                                         │
 Efecto: GOPATH no es recomendado      │

```

### Era 5: 2021-2024 - Estabilización (Go 1.18 a Go 1.22)

**Go 1.18 (Marzo 2022) - GENERICS (Primera vez)**

```
HITO HISTÓRICO: Generics en Go

Después de 13 AÑOS sin generics, Go finalmente añadió soporte.

Antes (Go 1.17):
func Retornar(x interface{}) interface{} {
    return x
}

Después (Go 1.18):
func Retornar[T any](x T) T {
    return x
}
```

**Razón de espera:**
Go no quiso generics complejos como Java/C++. Esperó hasta poder implementarlos SIMPLEMENTE.

**Go 1.19-1.22 (2022-2024):**

```
 Go 1.19: Mejoras de performance
 Go 1.20: Mejores errores wrapping
 Go 1.21: Iterators, pruebas de performance
 Go 1.22: Range sobre integers, loop var fix
```

### Versión Actual Recomendada (2024)

**Go 1.22 (Febrero 2024)**

```
Estado actual:
 Estable y probado
 Todas las características modernas
 Performance óptima
 Soporte mínimo 5 años
 Recomendación: Instalar Go 1.22 o mayor
```

**Ciclo de soporte:**

```
Go 1.22 (Feb 2024) - Soporte hasta: Feb 2026+
Go 1.21 (Aug 2023) - Soporte hasta: Aug 2025+
Go 1.20 (Feb 2023) - Soporte terminó: Feb 2024
Go 1.19 (Sep 2022) - Soporte terminó: Sep 2023

Recomendación: Mantente en últimas 2-3 versiones
```

---

## 2.2 Arquitectura del Go SDK

### ¿Qu es el "Go SDK"?

SDK = Software Development Kit. El Go SDK es TODO lo que necesitas para compilar Go:

```
Go SDK (típicamente /usr/local/go en Unix, C:\Go en Windows)

 bin/
   ├── go              (El compilador + gestor)
   ├── gofmt           (Formateador de código)
   └── godoc           (Documentación)

 src/
   ├── runtime/        (Runtime de Go - memoria, GC, scheduler)
   ├── fmt/            (Paquete fmt - formatting)
   ├── io/             (Paquete io - Reader/Writer)
   ├── os/             (Paquete os - sistema operativo)
   ├── net/            (Paquete net - networking)
   ├── encoding/       (XML, JSON, base64...)
   ├── crypto/         (Criptografía)
   ├── sync/           (Sincronización)
   ├── context/        (Contexto)
   └── ... (100+ más paquetes)

 lib/
   ├── time/tzdata     (Información de timezones)
   └── go/src          (Fuente de stdlib)

 pkg/
    ├── linux_amd64/    (Stdlib compilada para Linux 64-bit)
    ├── darwin_arm64/   (Stdlib compilada para macOS ARM)
    ├── windows_amd64/  (Stdlib compilada para Windows 64-bit)
    └── ... (otras plataformas)
```

### El Corazón: El Runtime de Go

**¿Qué es el runtime?**

El runtime es el corazón de Go. Es código compiled que:

- Maneja memoria (allocation, garbage collection)
- Implementa el scheduler de goroutines
- Controla channels
- Maneja señales del SO
- Implementa concurrencia

```
Cada ejecutable de Go contiene el runtime.
Por eso los binarios de Go son "self-contained".
```

**Estructura interna del runtime:**

```
runtime/ (código fuente del runtime)

 malloc.go           (Memory allocator - cómo Go obtiene memoria)
 gc.go               (Garbage collector - cómo Go libera memoria)
 proc.go             (Processor - scheduler de goroutines)
 chan.go             (Implementación de channels)
 atomic.go           (Operaciones atómicas)
 signal.go           (Manejo de señales del SO)
 mgc.go              (Mark-and-sweep GC)
 syscall_*.s         (Llamadas al SO - assembler)
 ... (más internals)
```

### El Compilador: `go build`

**¿Cómo funciona la compilación?**

```
Flujo de compilación (simplificado):

1. Parse (Análisis sintáctico)
   └─ Leer código Go, crear AST (Abstract Syntax Tree)

2. Type Checking (Verificación de tipos)
   └─ Verificar que tipos son válidos
   └─ NO hay conversiones implícitas permitidas

3. SSA (Static Single Assignment)
   └─ Convertir a IR intermedio optimizado

4. Machine Code Generation (Generación de código máquina)
   └─ SSA → Assembly específico de arquitectura
   └─ Optimizaciones de compilador aplicadas

5. Linking (Enlazamiento)
   └─ Juntar archivos objeto
   └─ Resolver símbolos externos
   └─ Incrustrar runtime
   └─ Crear binario ejecutable final

TOTAL: Tipicamente 1-5 SEGUNDOS (mucho más rápido que C++)
```

**¿Por qué Go compila tan rápido?**

```
Factores de velocidad:

1. Sin preprocessor
   └─ C++ debe procesar #include recursivamente
   └─ Go no tiene esto

2. Resolución simple de dependencias
   └─ C++ busca en muchos directorios
   └─ Go: import "paquete" → archivo conocido

3. Sin templates complejos
   └─ C++ templates generan código por cada tipo
   └─ Go 1.18+ generics son simples

4. Compilador eficiente
   └─ Escrito en Go, altamente optimizado
   └─ Caché inteligente de builds

Resultado: Go compila 10-100x más rápido que C++
```

### Tooling Incluido

Go viene con herramientas integradas:

```
go build        Compilador
go run          Ejecutar directamente
go test         Testing framework
go fmt          Formateador de código
go vet          Análisis estático
go doc          Documentación
go get          Descargador de paquetes
go mod          Gestor de módulos
go tool         Utilidades diversas
go version      Ver versión de Go
go env          Ver variables de entorno
```

**Esto es IMPORTANTE:** No necesitas descargar herramientas externas. Vienen incluidas.

---

## 2.3 Instalación por Sistema Operativo

### Instalación en macOS

#### Opción 1: Homebrew (RECOMENDADO)

**Ventajas:**

- Más fácil
- Actualización automática
- Gestión de múltiples versiones

**Pasos:**

```bash
# 1. Instalar Homebrew si no lo tienes
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Instalar Go
brew install go

# 3. Verificar instalación
go version
# Output: go version go1.22.0 darwin/amd64
```

**Actualizar a última versión:**

```bash
brew upgrade go
```

#### Opción 2: Instalador Oficial DMG

**Pasos:**

```bash
# 1. Descargar desde https://go.dev/dl
# Seleccionar: go1.22.0.darwin-amd64.pkg

# 2. Ejecutar el instalador (GUI)
# Instala a: /usr/local/go

# 3. Verificar
go version
```

#### Opción 3: Instalación Manual (Avanzado)

```bash
# 1. Descargar
cd ~/Downloads
wget https://go.dev/dl/go1.22.0.darwin-amd64.tar.gz

# 2. Verificar integridad SHA256
sha256sum go1.22.0.darwin-amd64.tar.gz
# Comparar con https://go.dev/dl (lista de checksums)

# 3. Instalar
sudo rm -rf /usr/local/go  # Limpiar versión anterior
sudo tar -C /usr/local -xzf go1.22.0.darwin-amd64.tar.gz

# 4. Configurar PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
source ~/.zshrc

# 5. Verificar
go version
```

### Instalación en Linux (Ubuntu/Debian)

#### Opción 1: Gestor de Paquetes

```bash
# 1. Actualizar índice de paquetes
sudo apt update

# 2. Instalar Go
sudo apt install golang-go

# 3. Verificar (puede ser versión vieja)
go version
```

**Nota:** apt puede tener Go viejo. Si necesitas 1.22, usa Opción 2.

#### Opción 2: Instalación Manual (RECOMENDADO para 1.22)

```bash
# 1. Descargar
cd /tmp
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# 2. Verificar SHA256
sha256sum go1.22.0.linux-amd64.tar.gz
# Copiar checksum de https://go.dev/dl y comparar

# 3. Instalar
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# 4. Configurar PATH (si no está)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 5. Verificar
go version
which go  # Debe ser /usr/local/go/bin/go
```

#### Opción 3: Múltiples Versiones Simultáneamente

```bash
# Instalar Go 1.22
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
sudo mv /usr/local/go /usr/local/go1.22

# Instalar Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
sudo mv /usr/local/go /usr/local/go1.21

# Crear symlink a la versión default
sudo ln -s /usr/local/go1.22 /usr/local/go

# Switchear versión
sudo ln -sfn /usr/local/go1.21 /usr/local/go
go version  # Ahora 1.21
sudo ln -sfn /usr/local/go1.22 /usr/local/go
go version  # Ahora 1.22
```

### Instalación en Windows

#### Opción 1: Instalador MSI (RECOMENDADO)

**Pasos:**

```
1. Descargar https://go.dev/dl
   Seleccionar: go1.22.0.windows-amd64.msi

2. Ejecutar el instalador
   ├─ Acepta licencia
   ├─ Selecciona destino (default: C:\Program Files\Go)
   ├─ Selecciona componentes
   └─ Instala

3. Verifica (abre PowerShell NEW):
   go version
   # Output: go version go1.22.0 windows/amd64

4. IMPORTANTE: Cierra y reabre PowerShell
   El instalador actualiza PATH, necesitas nueva sesión
```

#### Opción 2: Instalación Manual

```powershell
# 1. Descargar (en PowerShell)
Invoke-WebRequest -Uri "https://go.dev/dl/go1.22.0.windows-amd64.zip" -OutFile "$env:TEMP\go1.22.0.zip"

# 2. Extraer (reemplaza versión anterior)
Remove-Item -Recurse "C:\Program Files\Go" -Force -ErrorAction SilentlyContinue
Expand-Archive -Path "$env:TEMP\go1.22.0.zip" -DestinationPath "C:\Program Files"

# 3. Añadir a PATH (requiere reboot o nueva sesión)
$env:Path += ";C:\Program Files\Go\bin"

# 4. Verificar (NEW PowerShell)
go version
```

#### Opción 3: Chocolatey (Si tienes)

```powershell
choco install golang

# Actualizar
choco upgrade golang
```

---

## 2.4 GOPATH vs Go Modules - Evolución

### La Era GOPATH (2009-2019)

**¿Qué es GOPATH?**

GOPATH es la ubicación donde Go descarga y almacena TODOS los paquetes:

```bash
# Típicamente:
export GOPATH=$HOME/go

# Estructura:
$HOME/go/
 bin/           (ejecutables compilados)
 pkg/           (paquetes compilados)
 src/
 github.com/  
    │   ├── usuario1/
    │   │   ├── proyecto1/
    │   │   │   └── main.go
    │   │   └── proyecto2/
       └── main.go    │  
    │   └── usuario2/
    │       └── proyecto3/
    │           └── main.go
    │
    └── golang.org/
        └── x/
            ├── sys/
            ├── net/
            └── ...
```

**El Problema Fundamental de GOPATH:**

```
GOPATH almacena TODOS los proyectos Y TODAS las dependencias
en el MISMO lugar.

Problema 1: Versionamiento
 Proyecto A necesita: github.com/pkg/db v1.0
 Proyecto B necesita: github.com/pkg/db v2.0
 GOPATH descarga: SOLO UNA versión
    └─ Un proyecto se rompe

Problema 2: Reproducibilidad
 Desarrollador A: go get (descarga v1.5)
 Desarrollador B: go get (descarga v1.6 - salió versión nueva)
 Mismo código: Comportamiento diferente
 Bug imposible de reproducir

Problema 3: Limpieza
 GOPATH crece indefinidamente
 Cada paquete que usaste: Permanece forever
 Gigabytes de código viejo acumulado
 Difícil limpiar sin romper algo
```

**Comparación Real:**

```
2007-2015: Todos los lenguajes tenían este problema
 Python: pip descargaba ÚLTIMA versión siempre
 Node.js: npm descargaba ÚLTIMA versión siempre
 Ruby: Gem descargaba ÚLTIMA versión siempre
 Go: GOPATH hacía lo mismo

2015+: Otros lenguajes resolvieron
 Python: virtualenv, venv
 Node.js: package-lock.json
 Ruby: Gemfile.lock
 Go: Esperó hasta Go 1.11 (vendoring, luego módulos)
```

### Transición a Go Modules (2018-2021)

**Go 1.11 (2018): Módulos Beta**

```
Nueva estructura (proyecto actual, NO GOPATH):

proyecto/
 go.mod              (definición de módulo)
 go.sum              (checksums de dependencias)
 main.go
 internal/           (código privado del proyecto)
 cmd/
    └── herramienta/
        └── main.go
```

**¿Qué es go.mod?**

```
module github.com/usuario/proyecto

go 1.22

require (
    github.com/google/uuid v1.3.0
    github.com/prometheus/client_golang v1.14.0
)

exclude (
    github.com/algo/problem v1.0.0  // No usar esta versión
)

retract (
    v1.0.0  // Esta versión tuvo un bug, no uses
    [v2.0.0, v2.0.5]  // Rango de versiones retraídas
)

replace (
    github.com/original/path => github.com/fork/path v1.0.0
)
```

**¿Qué es go.sum?**

```
github.com/google/uuid v1.3.0 h1:ADr5Z1zNQbMF...
github.com/google/uuid v1.3.0/go.mod h1:...
github.com/prometheus/client_golang v1.14.0 h1:PK1...
github.com/prometheus/client_golang v1.14.0/go.mod h1:...

Cada línea = hash criptográfico de una versión específica de un paquete.

Propósito:
 Verificar que el código descargado no fue modificado
 Garantizar que todos los desarrolladores descargan EXACTAMENTE lo mismo
 Proteger contra cambios maliciosos en repositorios
```

**Go 1.16 (2021): Módulos por Defecto**

```
DECISIÓN FINAL: GOPATH fue deprecado

Comportamiento:
 Go 1.16+: Si hay go.mod → módulos (ignorar GOPATH)
 Go 1.16+: Si NO hay go.mod → usa GOPATH (legacy)
 Go 1.17+: Preferencia absoluta a módulos
 Go 1.18+: GOPATH casi no se usa

Impacto:
 GOPATH es historia, olvídalo
```

---

## 2.5 Go Modules Detallado

### Crear un Proyecto con Módulos

**Paso 1: Crear directorio**

```bash
mkdir mi-proyecto
cd mi-proyecto
```

**Paso 2: Inicializar módulo**

```bash
go mod init github.com/usuario/mi-proyecto
```

**Genera archivo go.mod:**

```
module github.com/usuario/mi-proyecto

go 1.22
```

**¿Qué significa cada línea?**

```
module github.com/usuario/mi-proyecto
 "module": Nombre único de tu proyecto
 "github.com/usuario/mi-proyecto": Importpath completo
 Se usa en 'import' en otros proyectos:
   └─ import "github.com/usuario/mi-proyecto/pkg"

go 1.22
 "go": Versión mínima de Go requerida
 "1.22": Este código requiere Go 1.22 o mayor
 Si alguien con Go 1.21 intenta compilar: ERROR
```

### Añadir Dependencias

**Opción 1: Automático (Recomendado)**

```bash
# Editar main.go para usar un paquete
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "github.com/google/uuid"
)

func main() {
    id := uuid.New()
    fmt.Println(id)
}
EOF

# Ejecutar go run (descarga automáticamente)
go run main.go

# Resultado: go.mod actualizado
# ├─ require github.com/google/uuid v1.3.0
# └─ go.sum creado con hashes
```

**Opción 2: Manual con go get**

```bash
# Obtener versión específica
go get github.com/google/uuid@v1.3.0

# Obtener última versión
go get github.com/google/uuid@latest

# Obtener rama
go get github.com/google/uuid@main
```

### Comandos go mod Importantes

**Gestión de módulos:**

```bash
go mod init <module>              Crear nuevo módulo
go mod tidy                       Limpiar dependencias no usadas
go mod graph                      Ver árbol de dependencias
go mod verify                     Verificar checksums
go mod download <package>         Descargar sin compilar

go mod vendor                     Crear carpeta vendor/ (copiar deps)
go mod why <package>             Por qué está esta dependencia
```

**Gestión de versiones:**

```bash
go get -u ./...                  Actualizar todo a latest
go get <package>@latest          Actualizar a latest
go get <package>@v1.2.3          Bajar a versión específica
go get <package>@none            Remover dependencia
```

### Versioning Semántico

Go usa Semantic Versioning (SemVer):

```
v[MAYOR].[MENOR].[PATCH]

v1.2.3
 │ └─ PATCH: Bug fixes (v1.2.3 → v1.2.4)
 └─── MENOR: Nuevas features, backward compatible (v1.2.0 → v1.3.0)
 MAYOR: Cambios que rompen compatibilidad (v1.x.x → v2.0.0)

Reglas de compatibilidad:
 v1.2.2 → v1.2.3: Puedo actualizar sin miedo (bug fix)
 v1.2.0 → v1.3.0: Puedo actualizar (backward compatible)
 v1.x → v2.0.0: CUIDADO, puede romper código
```

---

## 2.6 Comandos go - Referencia Completa

### go run - Ejecutar Directamente

```bash
# Compilar y ejecutar sin crear binario
go run main.go

# Múltiples archivos
go run main.go config.go

# Con argumentos
go run main.go arg1 arg2
```

**Casos de uso:**

- Prototipado rápido
- Testing durante desarrollo
- Scripts de Go

### go build - Compilar

```bash
# Build en directorio actual
go build

# Especificar nombre de salida
go build -o miapp

# Build para otra plataforma
GOOS=linux GOARCH=amd64 go build

# Build optimizado para distribución
go build -ldflags="-s -w" -o miapp
 -s: Remover tabla de símbolos
 -w: Remover debugging info
 Resultado: Binario más pequeño
```

### go test - Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests de paquete específico
go test ./pkg/mypackage

# Con cobertura
go test -cover ./...

# Ver qué tests corren
go test -v ./...

# Benchmarks
go test -bench=. ./...

# Fuzzing (Go 1.18+)
go test -fuzz=FuzzTest -fuzztime=30s
```

### go fmt - Formateador

```bash
# Formatear archivos en directorio actual
go fmt ./...

# Mostrar cambios sin aplicar (-d = diff)
gofmt -d main.go

# Verificar formateado (para CI/CD)
gofmt -l . | grep -q . && exit 1  # Falla si hay archivos mal formateados
```

### go vet - Análisis Estático

```bash
# Buscar problemas comunes
go vet ./...

# Análisis específico
go vet -composites=false ./...

# Race detector (detecta data races)
go test -race ./...
```

### go get - Descargador de Paquetes

```bash
# Descargar/actualizar paquete
go get github.com/sirupsen/logrus

# Con versión especica
go get github.com/sirupsen/logrus@v1.8.1

# Actualizar todo
go get -u ./...

# Remover paquete
go get github.com/sirupsen/logrus@none

# Actualizar módulo de Go requerido
go get go@1.22
```

### go install - Instalar Ejecutables

```bash
# Instalar herramienta Go en $GOBIN (típicamente ~/go/bin)
go install github.com/cosmtrek/air@latest

# Luego usar
air

# Instalar binario ejecutable
go install ./cmd/myapp
```

### go list - Listar Información

```bash
# Listar paquetes en proyecto
go list ./...

# Listar dependencias
go list -m all

# Información detallada de módulo
go list -m -json all

# Findbugs
go list -json ./...
```

### go clean - Limpiar

```bash
# Remover binarios compilados
go clean

# Remover cache
go clean -cache

# Limpiar todo (caché, test binarios, etc)
go clean -cache -testcache
```

### go doc - Documentación

```bash
# Ver documentación de paquete
go doc fmt

# Ver documentación de función
go doc fmt.Println

# Servidor web de documentación (localhost:6060)
godoc -http=:6060
```

---

## 2.7 Workspace de Go (Multi-módulo)

### ¿Cuándo Necesitas Workspace?

Workspace es para proyectos que tienen múltiples módulos interdependientes:

```
Escenario típico:
 gateway/ (módulo 1: API Gateway)
 service1/ (módulo 2: Microservicio 1)
 service2/ (módulo 3: Microservicio 2)
 shared/ (módulo 4: Código compartido)

Todos dependen de shared.
Durante desarrollo, quieres cambios en shared reflejados inmediatamente.
```

### Crear Workspace

**Paso 1: Estructura de directorios**

```bash
mkdir workspace
cd workspace

# Crear módulos
mkdir -p {gateway,service1,service2,shared}

# Cada uno es un módulo Go
cd gateway && go mod init github.com/tu-org/gateway && cd ..
cd service1 && go mod init github.com/tu-org/service1 && cd ..
cd service2 && go mod init github.com/tu-org/service2 && cd ..
cd shared && go mod init github.com/tu-org/shared && cd ..
```

**Paso 2: Crear go.work**

```bash
cat > go.work << 'EOF'
go 1.22

use (
    ./gateway
    ./service1
    ./service2
    ./shared
)
EOF
```

**¿Qué hace go.work?**

```

 go.work dice al compilador:         │
                                      │
 "Cuando compiles en este workspace:│
  - gateway, service1, service2     │
  - Usan las VERSIONES de LOCALES  
    shared (no las de remote)       │
  - Cambios en shared son inmediatos│

```

### Usar Workspace en Imports

```go
// En gateway/main.go
package main

import (
    "github.com/tu-org/shared"
)

func main() {
    data := shared.GetData()
}
```

**El workspace automáticamente:**

- Busca shared/ localmente
- No descarga de GitHub
- Usa código local
- Cambios en shared son reflejados al compilar

---

## 2.8 Configuración de Entorno

### Variables de Entorno Importantes

**GOROOT - Dónde está instalado Go**

```bash
echo $GOROOT
# /usr/local/go (típicamente)

# Normalmente NOT necesitas configurar, el instalador lo hace
# Pero si instalaste manualmente:
export GOROOT=/ruta/a/go
```

**GOPATH - Ubicación de proyectos (Legacy)**

```bash
echo $GOPATH
# $HOME/go (típicamente)

# Ya no es recomendado con Módulos
# Pero algunos herramientas aún lo usan:
export GOPATH=$HOME/go
```

**GOBIN - Dónde instalar ejecutables**

```bash
echo $GOBIN
# $GOPATH/bin (por defecto)

# O configurar:
export GOBIN=$HOME/.local/bin
go install github.com/cosmtrek/air@latest
# Instala en ~/.local/bin
```

**GO111MODULE - Habilitar módulos (Legacy)**

```bash
echo $GO111MODULE
# (vacío o "on" en Go 1.16+)

# En Go 1.11-1.15:
export GO111MODULE=on  # Habilitar módulos experimentales

# En Go 1.16+: No es necesario, módulos son default
```

**GOOS - Sistema operativo target**

```bash
# Compilar para Linux
GOOS=linux go build

# Compilar para Windows
GOOS=windows go build

# Compilar para macOS
GOOS=darwin go build
```

**GOARCH - Arquitectura target**

```bash
# Compilar para 64-bit
GOARCH=amd64 go build

# Compilar para ARM (Raspberry Pi)
GOARCH=arm go build

# Compilar para ARM 64-bit
GOARCH=arm64 go build

# Compilar para 32-bit
GOARCH=386 go build

# Compilar para WebAssembly
GOOS=js GOARCH=wasm go build
```

**GOCACHE - Caché de compilación**

```bash
echo $GOCACHE
# ~/.cache/go-build (típicamente)

# Limpiar caché
go clean -cache

# Deshabilitar caché (para debugging)
GOCACHE=off go build
```

**CGO_ENABLED - Habilitar CGO (C interop)**

```bash
# Por defecto en Unix
CGO_ENABLED=1

# Deshabilitar (para binarios portátiles)
CGO_ENABLED=0 go build
```

### Configurar PATH

**Asegurar que /usr/local/go/bin está en PATH:**

```bash
# macOS/Linux (agregar a ~/.bashrc o ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin

# Windows PowerShell (permanente)
$env:Path += ";C:\Program Files\Go\bin"
[Environment]::SetEnvironmentVariable("Path", $env:Path, [EnvironmentVariableTarget]::User)
```

### Verificar Configuración

```bash
go env

# Output (ejemplo):
# GO111MODULE=on
# GOARCH=amd64
# GOBIN=
# GOCACHE=/home/user/.cache/go-build
# GOEXE=
# GOFLAGS=
# GOHOSTARCH=amd64
# GOHOSTOS=linux
# GOINSECURE=
# GOMODCACHE=/home/user/go/pkg/mod
# GONOPROXY=
# GOPATH=/home/user/go
# GOPROXY=https://proxy.golang.org,direct
# GOROOT=/usr/local/go
# GOSUMDB=sum.golang.org
# GOTMPDIR=
# GOTOOLCHAIN=auto
# GOWINERROR=
# GOWORK=
# GOFLAGS=
```

---

## 2.9 IDEs y Herramientas

### VS Code + Go Extension (RECOMENDADO)

**Instalación:**

```
1. Descargar VS Code: https://code.visualstudio.com
2. Instalar: code → Extensions → Buscar "Go"
3. Instalar: "Go" (official extension by Google)
4. Reiniciar VS Code
```

**Configuración (settings.json):**

```json
{
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    },
    "editor.defaultFormatter": "golang.go"
  },
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.useLanguageServer": true,
  "go.languageServerFlags": [
    "-rpc.trace"
  ],
  "gopls": {
    "gofumpt": true,
    "staticcheck": true,
    "usePlaceholders": true
  }
}
```

**Características:**

```
 Autocomplete (powered by gopls)
 Go to definition (F12)
 Find references (Ctrl+Shift+F10)
 Format on save (gofmt)
 Linting (golangci-lint)
 Debugging (Delve debugger)
 Testing (run tests desde editor)
```

### GoLand (JetBrains IDE)

**Características Premium:**

```
 IDE completo (no solo editor)
 Refactoring avanzado
 Integración Git excelente
 Debugging visual
 Database tools
 Costo: ~$150/año (30 días trial)
```

### Vim/Neovim + vim-go

**Setup:**

```vim
" En ~/.vimrc o init.vim:
Plug 'fatih/vim-go'

" Keybindings comunes:
" :GoRun           - Ejecutar
" :GoTest          - Tests
" :GoFmt           - Formatea
" :GoLint          - Lint
" :GoDef           - Ir a definición
```

### Otros Editores

```
Emacs:   go-mode.el
Sublime: GoSublime
Atom:    go-plus
```

---

## 2.10 Troubleshooting Común

### Problema: "go: command not found"

**Causa:** go no está en PATH

**Solución:**

```bash
# Verificar si Go está instalado
ls /usr/local/go/bin/go

# Si existe, añadir a PATH
export PATH=$PATH:/usr/local/go/bin

# Hacer permanente (macOS/Linux):
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# En Windows: Reiniciar PowerShell NUEVA sesión
```

### Problema: Versión antigua de Go

**Causa:** Múltiples instalaciones de Go

**Solución:**

```bash
# Encontrar todas las instalaciones
which -a go
find /usr -name go -type f 2>/dev/null

# Remover viejas
sudo rm -rf /usr/local/go
sudo rm -rf ~/go  # Si usas GOPATH legacy

# Instalar nueva
# (seguir pasos de instalación)

# Verificar
go version
```

### Problema: "cannot find module"

**Causa:** Dependencia no en go.mod

**Solución:**

```bash
# Ejecutar go mod tidy (descarga deps faltantes)
go mod tidy

# O manualmente
go get github.com/paquete/nombre
```

### Problema: GOPATH vs Modules confusión

**Síntoma:** Código funciona en un lado, no en otro

**Causa:** Mezcla de GOPATH y Modules

**Solución:**

```bash
# Verificar si estás en módulo
ls go.mod  # Si existe, eres en módulo

# Asegurar GO111MODULE correcto
export GO111MODULE=on

# Limpiar GOPATH legacy
rm -rf ~/go/src/github.com/...

# Usar módulos exclusivamente
```

### Problema: Binario lento

**Causa:** Includes debug info

**Solución:**

```bash
# Build optimizado
go build -ldflags="-s -w"
# -s: Sin tabla de símbolos
# -w: Sin info de debugging

# Ejemplo: Antes 15 MB, después 5 MB
```

### Problema: Cross-compile falla

**Causa:** CGO no disponible

**Solución:**

```bash
# Deshabilitar CGO
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Nota: Si tu código usa CGO, esto fallará
# CGO es para llamar código C desde Go
```

### Problema: Cache corrupto

**Causa:** Build caché corrupto

**Solución:**

```bash
# Limpiar cache
go clean -cache

# Más agresivo
go clean -cache -testcache -modcache

# Reconstruit todo
go build -a ./...  # -a = rebuild all
```

---

**Fin del Capítulo 2**

---
