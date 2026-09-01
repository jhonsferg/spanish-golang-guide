# Capítulo 27: OS - Interacción con el sistema operativo

## Introducción

El paquete `os` es fundamental en Go para cualquier aplicación que necesite interactuar con el sistema operativo subyacente. Proporciona una interfaz uniforme para trabajar con archivos, directorios, permisos, variables de entorno, argumentos de línea de comandos, señales y procesos. Este capítulo explora cómo Go abstrae las diferencias entre sistemas operativos (Windows, Linux, macOS) permitiendo escribir código portable.

---

## 27.1 ¿Qué es el os Package?

### Propósito y Alcance

El paquete `os` de Go es la puerta de entrada para la interacción de bajo nivel con el sistema operativo. Proporciona:

- **Operaciones con archivos**: abrir, crear, leer, escribir, eliminar
- **Operaciones con directorios**: listar, crear, traversar
- **Metadatos de archivos**: permisos, tamaño, tiempos de modificación
- **Variables de entorno**: acceder y modificar
- **Argumentos de línea de comandos**: parsing básico
- **Señales del SO**: manejo de SIGTERM, SIGINT, etc.
- **Ejecución de procesos**: crear y comunicarse con subprocesos
- **Información del sistema**: hostname, usuario actual, variables de entorno

### Abstracción Multiplataforma

Go se esfuerza por proporcionar la misma API en diferentes sistemas operativos:

```
┌─────────────────────────────────────────────────────────┐
│          Código Go (Independiente de SO)               │
│                                                         │
│  os.Open("archivo.txt")                                │
│  os.Mkdir("directorio")                                │
│  os.Setenv("VAR", "valor")                             │
└─────────────────────────────────────────────────────────┘
              ↓                    ↓                    ↓
       ┌─────────────┐    ┌─────────────┐    ┌──────────────┐
       │   Windows   │    │   Linux     │    │    macOS     │
       │   (NTFS)    │    │   (ext4)    │    │   (APFS)     │
       │   (WINAPI)  │    │   (POSIX)   │    │   (POSIX)    │
       └─────────────┘    └─────────────┘    └──────────────┘
```

### Comparación con Otros Lenguajes

**Go vs C (libc)**:

```go
// Go: Simple y seguro
file, err := os.Open("archivo.txt")
if err != nil {
    return err
}
defer file.Close()

// C: Manual y propenso a errores
FILE *file = fopen("archivo.txt", "r");
if (file == NULL) {
    perror("fopen");
    return -1;
}
// ... código ...
fclose(file);
```

**Go vs Python**:

```go
// Go: Control explícito, compilado
file, err := os.Open("archivo.txt")
if err != nil {
    return err
}
defer file.Close()

// Python: Más conciso, interpretado
with open("archivo.txt", "r") as file:
    # ...
```

### El Paquete filepath

Complementa a `os` para manejo portátil de rutas:

```go
import "path/filepath"

// Estas funciones manejan separadores correctamente en cada SO
path := filepath.Join("directorio", "archivo.txt")
// Windows: "directorio\archivo.txt"
// Unix:    "directorio/archivo.txt"
```

---

## 27.2 Operaciones Básicas con Archivos

### Abrir Archivos

#### Open: Lectura

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Abre en modo lectura (solo lectura)
    file, err := os.Open("datos.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    // file es un *os.File
    fmt.Println("Archivo abierto:", file.Name())
}
```

**Modos de apertura**:

- `os.Open()`: Solo lectura (equivalente a `O_RDONLY`)
- `os.OpenFile()`: Control total sobre flags y permisos

#### OpenFile: Control Total

```go
package main

import (
    "os"
)

func main() {
    // Abrir para lectura y escritura, crear si no existe
    file, err := os.OpenFile("datos.txt",
        os.O_RDWR|os.O_CREATE|os.O_TRUNC,
        0644)
    if err != nil {
        panic(err)
    }
    defer file.Close()
}
```

**Flags comunes**:

| Flag | Significado |
|------|-------------|
| `O_RDONLY` | Solo lectura |
| `O_WRONLY` | Solo escritura |
| `O_RDWR` | Lectura y escritura |
| `O_CREATE` | Crear si no existe |
| `O_EXCL` | Falla si existe (con CREATE) |
| `O_TRUNC` | Truncar a tamaño 0 |
| `O_APPEND` | Añadir al final |

### Crear Archivos

```go
// Método 1: OpenFile con flags
file, err := os.OpenFile("nuevo.txt",
    os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
    0644)
if err != nil {
    panic(err)
}
defer file.Close()

// Método 2: Create (simplificado, equivalente a O_CREATE|O_WRONLY|O_TRUNC)
file, err := os.Create("nuevo.txt")
if err != nil {
    panic(err)
}
defer file.Close()
```

### Lectura de Archivos

#### Lectura Manual

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("datos.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Buffer de lectura
    buffer := make([]byte, 1024)

    for {
        n, err := file.Read(buffer)
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            panic(err)
        }
        fmt.Print(string(buffer[:n]))
    }
}
```

#### Lectura Completa

```go
import "os"

// Leer archivo completo (cuidado con archivos grandes)
contenido, err := os.ReadFile("datos.txt")
if err != nil {
    panic(err)
}
// contenido es []byte
```

#### Lectura con Scanner

```go
import (
    "bufio"
    "os"
)

file, err := os.Open("datos.txt")
if err != nil {
    panic(err)
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    linea := scanner.Text()
    // Procesar línea
}
```

### Escritura de Archivos

#### Escritura Manual

```go
package main

import (
    "os"
)

func main() {
    file, err := os.Create("salida.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    datos := []byte("Hola, mundo!\n")
    n, err := file.Write(datos)
    if err != nil {
        panic(err)
    }
    // n es cantidad de bytes escritos
}
```

#### Escritura Completa

```go
import "os"

contenido := []byte("Contenido del archivo\n")
err := os.WriteFile("archivo.txt", contenido, 0644)
if err != nil {
    panic(err)
}
```

#### Escritura Bufferizada

```go
import (
    "bufio"
    "os"
)

file, err := os.Create("salida.txt")
if err != nil {
    panic(err)
}
defer file.Close()

writer := bufio.NewWriter(file)
defer writer.Flush()

writer.WriteString("Línea 1\n")
writer.WriteString("Línea 2\n")
```

### Eliminación de Archivos

```go
package main

import (
    "os"
)

func main() {
    err := os.Remove("archivo.txt")
    if err != nil {
        panic(err)
    }
}
```

### Manejo de Errores Comunes

```go
import (
    "errors"
    "os"
)

file, err := os.Open("archivo.txt")
if err != nil {
    if errors.Is(err, os.ErrNotExist) {
        // El archivo no existe
    } else if errors.Is(err, os.ErrPermission) {
        // Permiso denegado
    } else {
        // Otro error
    }
}
```

---

## 27.3 Permisos de Archivo

### Sistema de Permisos Unix

Go sigue el modelo de permisos Unix/POSIX incluso en Windows:

```
Bits de permiso: 0755

┌─────────────────────────────────────────────┐
│  0   (special bits: setuid, setgid, sticky) │
│  7   (propietario: rwx)                    │
│  5   (grupo: r-x)                          │
│  5   (otros: r-x)                          │
└─────────────────────────────────────────────┘

rwx = 4 + 2 + 1 = 7
r-- = 4
-w- = 2
--x = 1
--- = 0
```

### Permisos Comunes

```go
package main

import (
    "os"
)

const (
    // Lectura, escritura, ejecución para propietario
    // Lectura, ejecución para grupo y otros
    PermisosDirectorio = 0755

    // Lectura y escritura para propietario
    // Lectura para grupo y otros
    PermisosArchivo = 0644

    // Lectura, escritura, ejecución solo para propietario
    PermisosSecreto = 0700

    // Script ejecutable
    PermisosScript = 0755
)

func main() {
    // Crear archivo con permisos específicos
    os.WriteFile("config.txt", []byte("datos"), PermisosArchivo)

    // Crear directorio con permisos específicos
    os.Mkdir("privado", PermisosSecreto)
}
```

### Cambiar Permisos

```go
import "os"

// Cambiar permisos de un archivo
err := os.Chmod("archivo.txt", 0600)
if err != nil {
    panic(err)
}

// Cambiar propietario (solo Unix, requiere permisos elevados)
err = os.Chown("archivo.txt", uid, gid)
if err != nil {
    panic(err)
}
```

### Verificar Permisos

```go
import (
    "os"
    "fmt"
)

info, err := os.Stat("archivo.txt")
if err != nil {
    panic(err)
}

// Obtener modo (permisos)
mode := info.Mode()
fmt.Printf("Permisos: %o\n", mode.Perm()) // Imprime en octal

// Verificar si tiene permisos específicos
if mode.Perm()&(1<<6) != 0 { // Verificar lectura del propietario
    fmt.Println("Propietario puede leer")
}
```

### Diferencias Windows vs Unix

```
┌──────────────────────────────────────────────────────┐
│                   Windows (NTFS)                     │
├──────────────────────────────────────────────────────┤
│ • ACLs (Access Control Lists)                       │
│ • Atributos: Read-only, Hidden, System              │
│ • Herencia de permisos por defecto                  │
│ • Go abstrae a modelo Unix                          │
└──────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────┐
│                   Unix/Linux                         │
├──────────────────────────────────────────────────────┤
│ • Permisos rwx (9 bits)                             │
│ • setuid, setgid, sticky bits                       │
│ • Modelo de propietario/grupo/otros                 │
│ • Modelo nativo de Go                               │
└──────────────────────────────────────────────────────┘
```

---

## 27.4 Información de Archivos (FileInfo y Stat)

### Obtener Metadatos

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    info, err := os.Stat("archivo.txt")
    if err != nil {
        panic(err)
    }

    // FileInfo interface
    fmt.Println("Nombre:", info.Name())
    fmt.Println("Tamaño:", info.Size(), "bytes")
    fmt.Println("Modo:", info.Mode())
    fmt.Println("Es directorio:", info.IsDir())
    fmt.Println("Última modificación:", info.ModTime())

    // Obtener sys() para info específica del SO
    sys := info.Sys()
    fmt.Println("Sistema:", sys)
}
```

### Interfaz FileInfo

```go
// Contrato de FileInfo
type FileInfo interface {
    Name() string        // Nombre base del archivo
    Size() int64         // Tamaño en bytes
    Mode() FileMode      // Bits de permiso y tipo
    ModTime() time.Time  // Última modificación
    IsDir() bool         // ¿Es directorio?
    Sys() interface{}    // Info específica del SO
}
```

### FileMode (Tipo de Archivo)

```go
import (
    "os"
    "fmt"
)

info, _ := os.Stat("ruta")
mode := info.Mode()

// Verificar tipo
if mode.IsDir() {
    fmt.Println("Es directorio")
}
if mode.IsRegular() {
    fmt.Println("Es archivo regular")
}

// Verificar si es symlink
if mode&os.ModeSymlink != 0 {
    fmt.Println("Es enlace simbólico")
}

// Otros tipos especiales
if mode&os.ModeCharDevice != 0 {
    fmt.Println("Dispositivo de caracteres")
}
if mode&os.ModeSocket != 0 {
    fmt.Println("Socket")
}
```

### Comparación de Tiempos

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    info, _ := os.Stat("archivo.txt")
    modTime := info.ModTime()

    // Hace cuánto fue modificado
    edad := time.Since(modTime)
    fmt.Printf("Modificado hace: %v\n", edad)

    // Archivo más viejo que N horas
    if edad > 24*time.Hour {
        fmt.Println("Archivo muy viejo, considerar limpieza")
    }
}
```

### Estadísticas Específicas del SO

```go
package main

import (
    "fmt"
    "os"
    "syscall"
)

func main() {
    info, _ := os.Stat("archivo.txt")

    // En Unix
    if stat, ok := info.Sys().(*syscall.Stat_t); ok {
        fmt.Println("Inode:", stat.Ino)
        fmt.Println("UID:", stat.Uid)
        fmt.Println("GID:", stat.Gid)
        fmt.Println("Número de enlaces:", stat.Nlink)
    }
}
```

---

## 27.5 Operaciones con Directorios

### Crear Directorios

```go
package main

import (
    "os"
)

func main() {
    // Crear un solo directorio
    err := os.Mkdir("nuevo_directorio", 0755)
    if err != nil {
        panic(err)
    }

    // Crear directorios anidados (con MkdirAll)
    err = os.MkdirAll("ruta/profunda/con/multiples/niveles", 0755)
    if err != nil {
        panic(err)
    }
}
```

### Listar Archivos en un Directorio

#### ReadDir: Recomendado

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Retorna []DirEntry ordenada por nombre
    entries, err := os.ReadDir(".")
    if err != nil {
        panic(err)
    }

    for _, entry := range entries {
        info, _ := entry.Info()
        fmt.Printf("%-20s %10d bytes %v\n",
            entry.Name(),
            info.Size(),
            entry.IsDir())
    }
}
```

#### ListDir (antigua, usa Readdirnames)

```go
// NO RECOMENDADO - usa ReadDir en su lugar
file, err := os.Open(".")
if err != nil {
    panic(err)
}
defer file.Close()

names, err := file.Readdirnames(0) // 0 = todos
if err != nil {
    panic(err)
}
```

### Traversar Directorios (Recursivo)

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func traversar(dir string, nivel int) error {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        // Indentación para mostrar profundidad
        indent := ""
        for i := 0; i < nivel; i++ {
            indent += "  "
        }

        ruta := filepath.Join(dir, entry.Name())

        if entry.IsDir() {
            fmt.Printf("%s📁 %s/\n", indent, entry.Name())
            traversar(ruta, nivel+1) // Recursivo
        } else {
            info, _ := entry.Info()
            fmt.Printf("%s📄 %s (%d bytes)\n",
                indent, entry.Name(), info.Size())
        }
    }
    return nil
}

func main() {
    traversar(".", 0)
}
```

### Usar filepath.Walk para Traversal

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Contar archivos y directorios
        if info.IsDir() {
            fmt.Printf("Directorio: %s\n", path)
        } else {
            fmt.Printf("Archivo: %s (%d bytes)\n", path, info.Size())
        }

        return nil
    })
}
```

### Eliminar Directorios

```go
import "os"

// Eliminar directorio vacío
err := os.Remove("directorio_vacio")

// Eliminar directorio y su contenido recursivamente
err := os.RemoveAll("directorio_con_contenido")
```

### Directorio de Trabajo Actual

```go
import (
    "fmt"
    "os"
)

// Obtener directorio actual
wd, err := os.Getwd()
if err != nil {
    panic(err)
}
fmt.Println("Directorio actual:", wd)

// Cambiar directorio de trabajo
err = os.Chdir("..")
if err != nil {
    panic(err)
}
```

---

## 27.6 Manejo de Rutas (Paths)

### Rutas Absolutas vs Relativas

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Ruta relativa
    ruta := "datos/archivo.txt"
    fmt.Println("Relativa:", ruta)

    // Convertir a absoluta
    absoluta, err := filepath.Abs(ruta)
    if err != nil {
        panic(err)
    }
    fmt.Println("Absoluta:", absoluta)

    // Verificar si es absoluta
    fmt.Println("¿Es absoluta?:", filepath.IsAbs(absoluta))
}
```

### Operaciones con Rutas

#### Join: Concatenar componentes

```go
import "path/filepath"

// Maneja separadores automáticamente
ruta := filepath.Join("home", "usuario", "documentos", "archivo.txt")
// En Windows: "home\usuario\documentos\archivo.txt"
// En Unix:    "home/usuario/documentos/archivo.txt"
```

#### Dir: Directorio padre

```go
import (
    "fmt"
    "path/filepath"
)

ruta := filepath.Join("home", "usuario", "archivo.txt")
dir := filepath.Dir(ruta)
fmt.Println("Directorio:", dir)
// Output: "home/usuario"
```

#### Base: Nombre de archivo

```go
import (
    "fmt"
    "path/filepath"
)

ruta := "/home/usuario/archivo.txt"
nombre := filepath.Base(ruta)
fmt.Println("Nombre:", nombre)
// Output: "archivo.txt"
```

#### Ext: Extensión de archivo

```go
import (
    "fmt"
    "path/filepath"
)

ruta := "/home/usuario/documento.pdf"
ext := filepath.Ext(ruta)
fmt.Println("Extensión:", ext)
// Output: ".pdf"
```

#### Split: Separar directorio y base

```go
import (
    "path/filepath"
)

dir, file := filepath.Split("/home/usuario/archivo.txt")
// dir = "/home/usuario/"
// file = "archivo.txt"
```

#### Clean: Normalizar ruta

```go
import (
    "fmt"
    "path/filepath"
)

ruta := "/home/usuario/../usuario/./archivo.txt"
limpia := filepath.Clean(ruta)
fmt.Println(limpia)
// Output: "/home/usuario/archivo.txt"
```

### Separadores Específicos del SO

```go
import (
    "fmt"
    "path/filepath"
    "runtime"
)

fmt.Println("Separador:", string(filepath.Separator))
fmt.Println("Separador de lista:", string(filepath.ListSeparator))

// En Windows
if runtime.GOOS == "windows" {
    fmt.Println("Separator: \\")
    fmt.Println("ListSeparator: ;")
}

// En Unix
if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
    fmt.Println("Separator: /")
    fmt.Println("ListSeparator: :")
}
```

### Glob: Patrones de archivos

```go
package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    // Encontrar todos los archivos .go
    matches, err := filepath.Glob("*.go")
    if err != nil {
        panic(err)
    }

    for _, match := range matches {
        fmt.Println(match)
    }
}
```

### Enlaces Simbólicos

```go
import (
    "os"
    "path/filepath"
)

// Crear enlace simbólico
err := os.Symlink("archivo_original.txt", "enlace.txt")
if err != nil {
    panic(err)
}

// Leer destino del enlace
destino, err := os.Readlink("enlace.txt")
if err != nil {
    panic(err)
}

// Resolver enlace a ruta real
ruta, err := filepath.EvalSymlinks("enlace.txt")
if err != nil {
    panic(err)
}
```

---

## 27.7 Variables de Entorno

### Acceder a Variables de Entorno

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Obtener variable de entorno
    home := os.Getenv("HOME")
    fmt.Println("HOME:", home)

    // Obtener con valor por defecto
    user := os.Getenv("USER")
    if user == "" {
        user = "unknown"
    }
    fmt.Println("USER:", user)

    // Obtener con ok (más idiomático)
    user, ok := os.LookupEnv("USER")
    if !ok {
        user = "unknown"
    }
    fmt.Println("USER:", user)
}
```

### Establecer Variables de Entorno

```go
import "os"

// Establecer variable
err := os.Setenv("MI_VAR", "valor")
if err != nil {
    panic(err)
}

// Verificar que se estableció
valor := os.Getenv("MI_VAR")
fmt.Println(valor) // "valor"
```

### Listar Todas las Variables

```go
package main

import (
    "fmt"
    "os"
    "sort"
    "strings"
)

func main() {
    // Obtener todas las variables
    envs := os.Environ()
    // envs es []string con formato "KEY=VALUE"

    // Imprimir ordenadas
    sort.Strings(envs)
    for _, env := range envs {
        parts := strings.SplitN(env, "=", 2)
        fmt.Printf("%-20s = %s\n", parts[0], parts[1])
    }
}
```

### Usar Variables en Aplicaciones

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    // Configuración desde entorno
    dbHost := os.Getenv("DB_HOST")
    if dbHost == "" {
        dbHost = "localhost"
    }

    dbPort := 5432
    if portStr := os.Getenv("DB_PORT"); portStr != "" {
        if p, err := strconv.Atoi(portStr); err == nil {
            dbPort = p
        }
    }

    debugMode, _ := strconv.ParseBool(os.Getenv("DEBUG"))

    fmt.Printf("Conectando a %s:%d (debug=%v)\n",
        dbHost, dbPort, debugMode)
}
```

### Archivos .env (Pattern común)

```go
// load_env.go
package main

import (
    "bufio"
    "os"
    "strings"
)

func LoadEnv(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        linea := scanner.Text()

        // Ignorar comentarios y líneas vacías
        if strings.HasPrefix(linea, "#") || linea == "" {
            continue
        }

        parts := strings.SplitN(linea, "=", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            os.Setenv(key, value)
        }
    }
    return scanner.Err()
}
```

---

## 27.8 Argumentos de Línea de Comandos

### os.Args: Acceso Básico

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // os.Args es []string
    // os.Args[0] es el nombre del programa
    // os.Args[1:] son los argumentos

    fmt.Println("Programa:", os.Args[0])
    fmt.Println("Argumentos:", os.Args[1:])

    // Iterar sobre argumentos
    for i, arg := range os.Args[1:] {
        fmt.Printf("Argumento %d: %s\n", i+1, arg)
    }
}
```

**Ejecución**:

```bash
$ go run main.go hola mundo
Programa: /tmp/go-build.../main
Argumentos: [hola mundo]
```

### Flag Package: Parsing Estructurado

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    // Definir flags
    nombre := flag.String("nombre", "usuario", "Nombre del usuario")
    edad := flag.Int("edad", 0, "Edad del usuario")
    activo := flag.Bool("activo", true, "¿Está activo?")

    flag.Parse()

    fmt.Printf("Nombre: %s\n", *nombre)
    fmt.Printf("Edad: %d\n", *edad)
    fmt.Printf("Activo: %v\n", *activo)

    // Argumentos posicionales restantes
    args := flag.Args()
    fmt.Printf("Otros argumentos: %v\n", args)
}
```

**Uso**:

```bash
$ go run main.go -nombre Juan -edad 30 -activo=false archivo.txt
Nombre: Juan
Edad: 30
Activo: false
Otros argumentos: [archivo.txt]
```

### Flag Con Variables Existentes

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    config := struct {
        Nombre   string
        Edad     int
        Activo   bool
        Verbose  bool
    }{}

    // Asociar flags a variables existentes
    flag.StringVar(&config.Nombre, "nombre", "usuario", "Nombre")
    flag.IntVar(&config.Edad, "edad", 0, "Edad")
    flag.BoolVar(&config.Activo, "activo", true, "Activo")
    flag.BoolVar(&config.Verbose, "v", false, "Modo verbose")

    flag.Parse()

    fmt.Printf("Config: %+v\n", config)
}
```

### Custom Flag Parsing

```go
package main

import (
    "fmt"
    "os"
    "strings"
)

func parseArgs(args []string) map[string]string {
    result := make(map[string]string)

    for i := 0; i < len(args); i++ {
        if strings.HasPrefix(args[i], "-") {
            key := strings.TrimPrefix(args[i], "-")

            if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
                result[key] = args[i+1]
                i++
            } else {
                result[key] = "true"
            }
        }
    }

    return result
}

func main() {
    flags := parseArgs(os.Args[1:])
    fmt.Printf("Flags: %v\n", flags)
}
```

### Subcomandos

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Uso: programa <comando> [opciones]")
        os.Exit(1)
    }

    cmd := os.Args[1]

    switch cmd {
    case "crear":
        fs := flag.NewFlagSet("crear", flag.ExitOnError)
        nombre := fs.String("nombre", "", "Nombre de archivo")
        fs.Parse(os.Args[2:])
        fmt.Printf("Creando: %s\n", *nombre)

    case "listar":
        fmt.Println("Listando archivos...")

    case "eliminar":
        fs := flag.NewFlagSet("eliminar", flag.ExitOnError)
        archivo := fs.String("archivo", "", "Archivo a eliminar")
        fs.Parse(os.Args[2:])
        fmt.Printf("Eliminando: %s\n", *archivo)

    default:
        fmt.Printf("Comando desconocido: %s\n", cmd)
        os.Exit(1)
    }
}
```

---

## 27.9 Manejo de Señales del Sistema Operativo

### Concepto de Señales

Las señales son eventos asincronos que el SO envía a procesos:

```
Señal común    Descripción                 CTRL+X
──────────────────────────────────────────────────
SIGTERM        Terminar (limpiable)  
SIGINT         Interrupt (Ctrl+C)         Ctrl+C
SIGQUIT        Quit                        Ctrl+\
SIGHUP         Desconexión terminal  
SIGKILL        Matar (no se puede atrapar)
SIGUSR1        Usuario 1 (custom)  
SIGUSR2        Usuario 2 (custom)
```

### Capturar Señales

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Canal para recibir señales
    sigChan := make(chan os.Signal, 1)

    // Registrar qué señales queremos capturar
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    fmt.Println("Esperando señal... (Presiona Ctrl+C)")

    // Bloquear hasta recibir señal
    sig := <-sigChan
    fmt.Printf("\nSeñal recibida: %v\n", sig)
}
```

### Graceful Shutdown

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Goroutine que simula trabajo
    go trabajar()

    // Esperar señal
    sig := <-sigChan
    fmt.Printf("\nRecibida señal %v, finalizando...\n", sig)

    // Cleanup
    limpiar()

    fmt.Println("Programa finalizado")
}

func trabajar() {
    for i := 1; i <= 1000; i++ {
        fmt.Printf("Trabajando... paso %d\n", i)
        time.Sleep(1 * time.Second)
    }
}

func limpiar() {
    fmt.Println("Limpiando recursos...")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Limpieza completada")
}
```

### Multiple Señales

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    sigChan := make(chan os.Signal, 2)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    done := make(chan bool)

    // Iniciar goroutines de trabajo
    for i := 1; i <= 3; i++ {
        go func(id int) {
            for {
                select {
                case <-done:
                    fmt.Printf("Worker %d finalizando\n", id)
                    return
                default:
                    fmt.Printf("Worker %d trabajando\n", id)
                    time.Sleep(1 * time.Second)
                }
            }
        }(i)
    }

    // Esperar señal
    fmt.Println("Esperando señales... (Presiona Ctrl+C)")
    <-sigChan
    fmt.Println("\nSeñal recibida, deteniendo workers...")

    close(done) // Avisar a todos los workers
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Saliendo")
}
```

---

## 27.10 Procesos y Ejecución de Comandos

### Ejecutar Comandos Simples

```go
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    // Comando simple sin output capturado
    cmd := exec.Command("ls", "-la", "/tmp")

    // Ejecutar y esperar
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Capturar Output

```go
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    // Capturar stdout
    cmd := exec.Command("echo", "Hola, mundo!")
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Println("Output:", string(output))
}
```

### Capturar Output y Error

```go
package main

import (
    "bytes"
    "fmt"
    "os/exec"
)

func main() {
    cmd := exec.Command("ls", "/directorio_inexistente")

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        fmt.Println("Error:", err)
    }

    fmt.Println("STDOUT:", stdout.String())
    fmt.Println("STDERR:", stderr.String())
}
```

### Pipes entre Procesos

```go
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    // Equivalente a: grep "go" < archivo.txt | wc -l

    grep := exec.Command("grep", "go")
    wc := exec.Command("wc", "-l")

    // Conectar salida de grep a entrada de wc
    pipe, err := grep.StdoutPipe()
    if err != nil {
        panic(err)
    }

    wc.Stdin = pipe

    // Capturar output de wc
    output, err := wc.Output()
    if err != nil {
        panic(err)
    }

    fmt.Println("Resultado:", string(output))
}
```

### Ambiente del Proceso

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    cmd := exec.Command("bash", "-c", "echo $MI_VAR")

    // Copiar ambiente actual
    cmd.Env = os.Environ()

    // Añadir variable específica
    cmd.Env = append(cmd.Env, "MI_VAR=valor_especial")

    output, _ := cmd.Output()
    fmt.Println("Output:", string(output))
}
```

### Directorio de Trabajo del Proceso

```go
import (
    "os"
    "os/exec"
)

cmd := exec.Command("ls", "-la")
cmd.Dir = "/tmp"  // Cambiar dir de trabajo

cmd.Run()
```

### Crear Proceso en Background

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    cmd := exec.Command("sleep", "10")

    // Iniciar sin esperar
    err := cmd.Start()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Proceso iniciado con PID: %d\n", cmd.Process.Pid)

    // Esperar en background
    go func() {
        err := cmd.Wait()
        if err != nil {
            fmt.Println("Error al esperar:", err)
        }
    }()

    fmt.Println("Continuando con el programa")
}
```

### Timeout en Procesos

```go
package main

import (
    "context"
    "fmt"
    "os/exec"
    "time"
)

func main() {
    // Crear contexto con timeout de 2 segundos
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, "sleep", "5")

    err := cmd.Run()
    if err != nil {
        fmt.Println("Error:", err)
        // Verifica si fue por timeout
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("¡Timeout!")
        }
    }
}
```

### Información del Proceso

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // PID del proceso actual
    fmt.Println("PID:", os.Getpid())

    // PID del proceso padre
    fmt.Println("PPID:", os.Getppid())

    // UID y GID (en Unix)
    fmt.Println("UID:", os.Getuid())
    fmt.Println("GID:", os.Getgid())

    // Hostname
    hostname, err := os.Hostname()
    if err != nil {
        panic(err)
    }
    fmt.Println("Hostname:", hostname)
}
```

---

## 27.11 Buenas Prácticas y Patrones

### Patrón Defer para Limpieza de Recursos

```go
package main

import (
    "fmt"
    "os"
)

func procesarArchivo(nombre string) error {
    file, err := os.Open(nombre)
    if err != nil {
        return err
    }
    // SIEMPRE usar defer para cerrar
    defer file.Close()

    // ... procesar archivo ...
    fmt.Println("Procesando:", file.Name())
    return nil
}

func main() {
    procesarArchivo("datos.txt")
    // El archivo se cierra automáticamente
}
```

### Validar Errores de Forma Exhaustiva

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("archivo.txt")
    if err != nil {
        // Distinguir tipos de error
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("Archivo no existe")
        } else if errors.Is(err, os.ErrPermission) {
            fmt.Println("Permiso denegado")
        } else {
            fmt.Printf("Error desconocido: %v\n", err)
        }
        return
    }
    defer file.Close()
}
```

### Funciones Helper para Operaciones Comunes

```go
package main

import (
    "os"
    "path/filepath"
)

// Verificar si archivo existe
func Existe(ruta string) bool {
    _, err := os.Stat(ruta)
    return err == nil
}

// Verificar si es directorio
func EsDirectorio(ruta string) bool {
    info, err := os.Stat(ruta)
    if err != nil {
        return false
    }
    return info.IsDir()
}

// Obtener tamaño de archivo en MB
func TamanoMB(ruta string) (float64, error) {
    info, err := os.Stat(ruta)
    if err != nil {
        return 0, err
    }
    return float64(info.Size()) / (1024 * 1024), nil
}

// Copiar archivo
func CopiarArchivo(src, dst string) error {
    contenido, err := os.ReadFile(src)
    if err != nil {
        return err
    }

    srcInfo, err := os.Stat(src)
    if err != nil {
        return err
    }

    return os.WriteFile(dst, contenido, srcInfo.Mode())
}
```

### Manejo Seguro de Rutas

```go
package main

import (
    "fmt"
    "path/filepath"
    "strings"
)

// Prevenir path traversal attacks
func ValidarRuta(basedir, userpath string) (string, error) {
    // Hacer absoluta desde basedir
    ruta := filepath.Join(basedir, userpath)

    // Resolver todos los symlinks y ..
    ruta, err := filepath.Abs(ruta)
    if err != nil {
        return "", err
    }

    // Verificar que sigue siendo dentro de basedir
    basedir, _ = filepath.Abs(basedir)
    if !strings.HasPrefix(ruta, basedir) {
        return "", fmt.Errorf("ruta fuera de permisos: %s", userpath)
    }

    return ruta, nil
}

func main() {
    // Permitir solo dentro de /home/usuario/datos
    ruta, err := ValidarRuta("/home/usuario/datos", "../../etc/passwd")
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Println("Ruta segura:", ruta)
}
```

### Patrón de Configuración

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

type Config struct {
    DataDir string
    LogFile string
    Debug   bool
}

func LoadConfig() *Config {
    config := &Config{
        DataDir: os.Getenv("APP_DATA_DIR"),
        LogFile: os.Getenv("APP_LOG_FILE"),
    }

    // Valores por defecto
    if config.DataDir == "" {
        home, _ := os.UserHomeDir()
        config.DataDir = filepath.Join(home, ".app", "data")
    }

    if config.LogFile == "" {
        home, _ := os.UserHomeDir()
        config.LogFile = filepath.Join(home, ".app", "app.log")
    }

    // Crear directorios si no existen
    os.MkdirAll(filepath.Dir(config.DataDir), 0755)
    os.MkdirAll(filepath.Dir(config.LogFile), 0755)

    return config
}

func main() {
    config := LoadConfig()
    fmt.Printf("Config: %+v\n", config)
}
```

### Antipatrones Comunes

```go
// ❌ ANTIPATRÓN 1: No cerrar archivos
func MaloNoClose(nombre string) string {
    file, _ := os.Open(nombre)
    // ¡Forgot defer file.Close()!
    contenido, _ := os.ReadFile(nombre)
    return string(contenido)
}

// ✅ CORRECTO: Cerrar con defer
func BuenoClose(nombre string) (string, error) {
    file, err := os.Open(nombre)
    if err != nil {
        return "", err
    }
    defer file.Close()

    contenido, err := os.ReadFile(nombre)
    return string(contenido), err
}

// ❌ ANTIPATRÓN 2: Ignorar errores de permisos
func MaloIgnorarPermisos(ruta string) {
    os.Chmod(ruta, 0777) // ¡Demasiado permisivo!
}

// ✅ CORRECTO: Permisos restrictivos
func BuenoPermisos(ruta string) {
    os.Chmod(ruta, 0644) // Usuario RW, otros R
}

// ❌ ANTIPATRÓN 3: Rutas hardcodeadas
func MaloRutasHardcoded() {
    file := os.Open("C:\\datos\\archivo.txt") // ¡Falla en Unix!
}

// ✅ CORRECTO: Usar filepath para portabilidad
func BuenoRutasPortables() {
    import "path/filepath"
    ruta := filepath.Join("datos", "archivo.txt")
    file, _ := os.Open(ruta)
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Implementar un `ls` Simple

**Objetivo**: Listar archivos en un directorio

```go
package main

import (
    "flag"
    "fmt"
    "os"
    "sort"
    "text/tabwriter"
)

func main() {
    // Flags
    long := flag.Bool("l", false, "Formato largo")
    flag.Parse()

    dir := "."
    if flag.NArg() > 0 {
        dir = flag.Arg(0)
    }

    entries, err := os.ReadDir(dir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    if *long {
        w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
        fmt.Fprintln(w, "TIPO\tPERMISOS\tTAMAÑO\tNOMBRE")

        for _, entry := range entries {
            info, _ := entry.Info()
            modo := info.Mode()

            tipo := "-"
            if entry.IsDir() {
                tipo = "d"
            }

            fmt.Fprintf(w, "%s\t%o\t%d\t%s\n",
                tipo, modo.Perm(), info.Size(), entry.Name())
        }
        w.Flush()
    } else {
        for _, entry := range entries {
            fmt.Println(entry.Name())
        }
    }
}
```

### Ejercicio 2: Copiar Archivo con Progreso

**Objetivo**: Copiar archivo mostrando progreso

```go
package main

import (
    "flag"
    "fmt"
    "io"
    "os"
)

func main() {
    flag.Parse()

    if flag.NArg() != 2 {
        fmt.Fprintf(os.Stderr, "Uso: %s <origen> <destino>\n", os.Args[0])
        os.Exit(1)
    }

    src := flag.Arg(0)
    dst := flag.Arg(1)

    // Abrir archivo origen
    srcFile, err := os.Open(src)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    defer srcFile.Close()

    // Obtener información del archivo origen
    srcInfo, err := srcFile.Stat()
    if err != nil {
        panic(err)
    }
    totalSize := srcInfo.Size()

    // Crear archivo destino
    dstFile, err := os.Create(dst)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    defer dstFile.Close()

    // Copiar con barra de progreso
    buffer := make([]byte, 1024*64)
    var copied int64

    for {
        n, err := srcFile.Read(buffer)
        if err != nil && err != io.EOF {
            panic(err)
        }

        if n > 0 {
            dstFile.Write(buffer[:n])
            copied += int64(n)

            // Mostrar progreso
            percent := (copied * 100) / totalSize
            fmt.Printf("\r[%-50s] %d%%",
                fmt.Sprintf("%s", "="), percent)
        }

        if err == io.EOF {
            break
        }
    }

    fmt.Println("\n¡Copia completada!")

    // Copiar permisos
    os.Chmod(dst, srcInfo.Mode())
}
```

### Ejercicio 3: Leer y Usar Variables de Entorno

**Objetivo**: Implementar un programa configurado por entorno

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

func main() {
    // Cargar configuración desde entorno
    appName := os.Getenv("APP_NAME")
    if appName == "" {
        appName = "MyApp"
    }

    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" {
        logLevel = "INFO"
    }

    cacheDir := os.Getenv("CACHE_DIR")
    if cacheDir == "" {
        home, _ := os.UserHomeDir()
        cacheDir = filepath.Join(home, ".cache", appName)
    }

    maxRetries, _ := strconv.Atoi(os.Getenv("MAX_RETRIES"))
    if maxRetries == 0 {
        maxRetries = 3
    }

    enableDebug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

    // Crear directorio cache si no existe
    os.MkdirAll(cacheDir, 0755)

    // Mostrar configuración
    fmt.Printf("Configuración:\n")
    fmt.Printf("  Aplicación: %s\n", appName)
    fmt.Printf("  Nivel de log: %s\n", logLevel)
    fmt.Printf("  Directorio cache: %s\n", cacheDir)
    fmt.Printf("  Max retries: %d\n", maxRetries)
    fmt.Printf("  Debug: %v\n", enableDebug)

    // Crear archivo de log
    logFile := filepath.Join(cacheDir, "app.log")
    file, err := os.OpenFile(logFile,
        os.O_CREATE|os.O_APPEND|os.O_WRONLY,
        0644)
    if err != nil {
        fmt.Println("Error creando log:", err)
        return
    }
    defer file.Close()

    // Escribir en log
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] %s: Started\n", timestamp, appName)
    file.WriteString(logEntry)

    fmt.Printf("\nLog escrito en: %s\n", logFile)
}
```

### Ejercicio 4: Ejecutar Comando y Capturar Output

**Objetivo**: Ejecutar comando del sistema y procesar resultado

```go
package main

import (
    "bytes"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
)

func main() {
    timeout := flag.Duration("timeout", 5*time.Second, "Timeout de ejecución")
    flag.Parse()

    if flag.NArg() == 0 {
        fmt.Fprintf(os.Stderr, "Uso: %s [opciones] <comando>\n", os.Args[0])
        os.Exit(1)
    }

    cmdArgs := flag.Args()
    cmd := cmdArgs[0]
    args := cmdArgs[1:]

    // Crear comando
    execCmd := exec.Command(cmd, args...)

    // Capturar stdout y stderr
    var stdout, stderr bytes.Buffer
    execCmd.Stdout = &stdout
    execCmd.Stderr = &stderr

    // Ejecutar con timeout
    done := make(chan error, 1)
    go func() {
        done <- execCmd.Run()
    }()

    // Esperar resultado o timeout
    select {
    case err := <-done:
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        }

        // Mostrar output
        if output := stdout.String(); output != "" {
            fmt.Println("STDOUT:")
            fmt.Println(output)
        }

        if output := stderr.String(); output != "" {
            fmt.Println("STDERR:")
            fmt.Println(output)
        }

    case <-time.After(*timeout):
        fmt.Fprintf(os.Stderr, "¡Timeout! Comando no terminó en %v\n", *timeout)
        if execCmd.Process != nil {
            execCmd.Process.Kill()
        }
        os.Exit(1)
    }
}
```

### Ejercicio 5: Graceful Shutdown con Signal Handling

**Objetivo**: Manejar señales del SO para terminar limpiamente

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

type Application struct {
    mu    sync.Mutex
    stop  chan struct{}
    wg    sync.WaitGroup
    count int
}

func (app *Application) Worker(id int) {
    defer app.wg.Done()

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-app.stop:
            fmt.Printf("Worker %d: Recibida señal de parada\n", id)
            return

        case <-ticker.C:
            app.mu.Lock()
            app.count++
            count := app.count
            app.mu.Unlock()

            fmt.Printf("Worker %d: Ciclo %d\n", id, count)
        }
    }
}

func (app *Application) Start(numWorkers int) {
    fmt.Printf("Iniciando aplicación con %d workers\n", numWorkers)

    for i := 1; i <= numWorkers; i++ {
        app.wg.Add(1)
        go app.Worker(i)
    }
}

func (app *Application) Stop() {
    fmt.Println("\nDeteniendo aplicación...")
    close(app.stop)
    app.wg.Wait()
    fmt.Println("Aplicación detenida correctamente")
}

func main() {
    app := &Application{
        stop: make(chan struct{}),
    }

    // Registrar handler de señales
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Iniciar aplicación
    app.Start(3)

    // Esperar señal
    sig := <-sigChan
    fmt.Printf("\nSeñal recibida: %v\n", sig)

    // Detener limpiamente
    app.Stop()
}
```

---

## Resumen y Puntos Clave

### Lo Más Importante

1. **Siempre usar `defer` para cerrar archivos y recursos**
2. **Validar errores de forma específica** usando `errors.Is()`
3. **Usar `filepath` para portabilidad** en diferentes sistemas operativos
4. **Manejar señales correctamente** para terminación limpia
5. **Validar entrada de usuario** antes de usar en rutas

### Diferencias Clave Go vs Otros Lenguajes

| Aspecto | Go | Python | C |
|--------|----|---------|----|
| Abstracción SO | Alta | Alta | Baja |
| API uniforme | Sí | Sí | No |
| Manejo de errores | explícito | excepciones | código retorno |
| Recursos | defer | context manager | manual |
| Procesos | os/exec | subprocess | fork/exec |

### Recursos Adicionales

- Documentación oficial: <https://pkg.go.dev/os>
- Documentación filepath: <https://pkg.go.dev/path/filepath>
- Documentación exec: <https://pkg.go.dev/os/exec>

---

**Fin del Capítulo 27**

Hemos cubierto completamente el paquete `os` de Go, desde operaciones básicas con archivos hasta manejo avanzado de procesos y señales. Este conocimiento es fundamental para desarrollar herramientas CLI, scripts de deployment y cualquier aplicación que interactúe con el sistema operativo.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/27-os-package/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/27-os-package):

```bash
cd examples/27-os-package
go run .
```
