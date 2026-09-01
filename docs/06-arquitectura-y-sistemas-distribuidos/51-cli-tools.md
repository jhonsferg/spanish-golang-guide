# Capítulo 51: CLI tools - Herramientas de línea de comandos

**Longitud:** ~1,500 líneas | **Complejidad:** Intermedio → Avanzado  
**Requisitos previos:** Cap. 3 (Fundamentos), Cap. 10 (I/O), Cap. 22 (Testing), Cap. 48 (Viper - Config)

---

## TABLA DE CONTENIDOS

1. [Introducción a CLI en Go](#51.1-introducción-a-cli-en-go)
2. [Flag Package (Estándar)](#51.2-flag-package-estándar)
3. [Cobra Framework](#51.3-cobra-framework)
4. [Urfave/CLI v2](#51.4-urfavecli-v2)
5. [Input/Output Avanzado](#51.5-inputoutput-avanzado)
6. [Formatting Output](#51.6-formatting-output)
7. [Configuration Files](#51.7-configuration-files)
8. [Error Handling y Logging](#51.8-error-handling-y-logging)
9. [Testing CLI](#51.9-testing-cli)
10. [Distribución y Packaging](#51.10-distribución-y-packaging)
11. [Buenas Prácticas y Case Studies](#51.11-buenas-prácticas-y-case-studies)

---

## 51.1 Introducción a CLI en Go

### 51.1.1 ¿Qué es una CLI?

Una **Command Line Interface (CLI)** es una aplicación que se ejecuta en terminal/consola y se controla mediante comandos de texto. A diferencia de:

- **GUI (Graphical User Interface):** Requiere ventanas, ratón, más recursos
- **Web UI:** Requiere servidor HTTP, navegador web
- **API (Application Programming Interface):** Interfaz programática para otros servicios

Las CLIs son ideales para:
- Automatización (scripts, CI/CD)
- Herramientas de desarrollo (compiladores, gestores de código)
- Administración de sistemas (DevOps)
- Procesos batch

### 51.1.2 ¿Por qué Go para CLI?

Go destaca en desarrollo de CLIs por:

```
┌─────────────────────────────────────────────────────┐
│         VENTAJAS DE GO PARA CLI                      │
├─────────────────────────────────────────────────────┤
│ ✅ Binario único (no requiere runtime/VM)            │
│ ✅ Startup ultrarrápido (ms, no segundos)            │
│ ✅ Cross-compilation trivial (GOOS, GOARCH)          │
│ ✅ Ecosistema de librerías CLI maduro                │
│ ✅ Concurrencia nativa (goroutines)                  │
│ ✅ Performance comparable a C/C++ (1.5-2x)           │
│ ✅ Tamaño binario pequeño (~10-50MB)                 │
│ ✅ Fácil testing (testable architecture)             │
└─────────────────────────────────────────────────────┘
```

**Comparativa de tamaño binario:**

| Lenguaje | Herramienta      | Tamaño   | Startup |
|----------|------------------|----------|---------|
| Go       | kubectl (1.29)   | 49 MB    | 50 ms   |
| Go       | Hugo             | 115 MB   | 30 ms   |
| Python   | pip (instalado)  | ~200 MB+ | 2000 ms |
| Node.js  | npm (instalado)  | ~300 MB+ | 1500 ms |
| Java     | Maven            | ~600 MB+ | 5000 ms |
| Ruby     | Rails CLI        | ~150 MB+ | 3000 ms |

### 51.1.3 Casos de Uso Reales en Producción

#### **kubectl** - Orquestación de Kubernetes
```bash
# Herramienta CLI escrita en Go (1.5M líneas)
$ kubectl get pods --namespace prod
$ kubectl apply -f deployment.yaml
```
- **Características:** Subcomandos complejos, config file hierarchy, plugins
- **Lecciones:** Context management, structured output, extensive help

#### **Hugo** - Static Site Generator
```bash
$ hugo new site myblog
$ hugo server --watch --port 1313
$ hugo build --minify
```
- **Características:** Múltiples flags, server embebido, watch mode
- **Lecciones:** Real-time interaction, verbose logging

#### **Docker CLI**
```bash
$ docker container run -it --name myapp alpine
$ docker compose up -d
```
- **Características:** Subcomandos, context switching, plugin system
- **Lecciones:** UX intuitivo, progressive disclosure

#### **Terraform CLI**
```bash
$ terraform plan -var="env=prod"
$ terraform apply -auto-approve
```
- **Características:** Complex state management, interactive modes, debug output

### 51.1.4 Ecosistema Go de Herramientas CLI

**Frameworks principales:**

1. **Flag (stdlib):** Minimalista, built-in
2. **Cobra:** Estándar de facto en Go (~40k GitHub stars)
3. **Urfave/CLI v2:** Ligero, funcional
4. **Click-like alternatives:** Kingpin, cli

**Librerías complementarias:**

- **Input:** promptui, survey, go-prompt
- **Output:** tablewriter, pterm, bubbletea (TUI)
- **Config:** viper (Cap. 48), env
- **Logging:** logrus, zap, go-log
- **Colors:** fatih/color, lipgloss, termenv

**Herramientas relacionadas:**

```go
// Validación de CLI
go install golang.org/x/lint/golint@latest

// Testing
go test ./...

// Build & Release
goreleaser
```

---

## 51.2 Flag Package (Estándar)

### 51.2.1 Parsing Básico de Flags

El paquete `flag` (stdlib) es suficiente para CLIs simples:

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // Definir flags
    namePtr := flag.String("name", "World", "name to greet")
    agePtr := flag.Int("age", 0, "age in years")
    verbosePtr := flag.Bool("verbose", false, "enable verbose output")

    flag.Parse()

    fmt.Printf("Hello, %s (age: %d)\n", *namePtr, *agePtr)
    if *verbosePtr {
        fmt.Println("Verbose mode enabled")
    }

    // Argumentos posicionales (sin flags)
    fmt.Printf("Remaining args: %v\n", flag.Args())
}
```

**Uso:**
```bash
$ go run main.go -name Alice -age 30 -verbose
Hello, Alice (age: 30)
Verbose mode enabled
Remaining args: []

$ go run main.go -name Bob file1.txt file2.txt
Hello, Bob (age: 0)
Remaining args: [file1.txt file2.txt]
```

### 51.2.2 Tipos de Datos Soportados

```go
package main

import (
    "flag"
    "fmt"
    "time"
)

func main() {
    // Tipos básicos
    stringFlag := flag.String("string", "default", "string value")
    intFlag := flag.Int("int", 0, "int value")
    int64Flag := flag.Int64("int64", 0, "int64 value")
    uintFlag := flag.Uint("uint", 0, "uint value")
    uint64Flag := flag.Uint64("uint64", 0, "uint64 value")
    floatFlag := flag.Float64("float", 0.0, "float value")
    boolFlag := flag.Bool("bool", false, "bool value")
    
    // Duración (formato: "300ms", "1h30m")
    durationFlag := flag.Duration("timeout", 10*time.Second, "timeout duration")

    flag.Parse()

    fmt.Printf("String: %s\n", *stringFlag)
    fmt.Printf("Int: %d\n", *intFlag)
    fmt.Printf("Float: %.2f\n", *floatFlag)
    fmt.Printf("Bool: %v\n", *boolFlag)
    fmt.Printf("Duration: %v\n", *durationFlag)
}
```

**Uso:**
```bash
$ go run main.go -string hello -int 42 -float 3.14 -bool -timeout 5s
String: hello
Int: 42
Float: 3.14
Bool: true
Duration: 5s
```

### 51.2.3 Custom Flag Types

Implementar `flag.Value` para tipos personalizados:

```go
package main

import (
    "flag"
    "fmt"
    "strings"
)

// Tipo personalizado: lista de strings
type StringSlice []string

func (s *StringSlice) String() string {
    return fmt.Sprintf("%v", *s)
}

func (s *StringSlice) Set(value string) error {
    *s = append(*s, value)
    return nil
}

// Tipo personalizado: host:puerto
type HostPort struct {
    Host string
    Port int
}

func (hp *HostPort) String() string {
    return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

func (hp *HostPort) Set(value string) error {
    parts := strings.Split(value, ":")
    if len(parts) != 2 {
        return fmt.Errorf("invalid host:port format: %s", value)
    }
    hp.Host = parts[0]
    _, err := fmt.Sscanf(parts[1], "%d", &hp.Port)
    return err
}

func main() {
    var tags StringSlice
    flag.Var(&tags, "tag", "tag (can be used multiple times)")

    hp := HostPort{Host: "localhost", Port: 8080}
    flag.Var(&hp, "server", "server host:port")

    flag.Parse()

    fmt.Printf("Tags: %v\n", tags)
    fmt.Printf("Server: %v\n", hp)
}
```

**Uso:**
```bash
$ go run main.go -tag prod -tag important -server example.com:443
Tags: [prod important]
Server: example.com:443
```

### 51.2.4 Subcomandos Básicos

Implementar subcomandos sin Cobra (usando flag.NewFlagSet):

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: myapp <command> [options]")
        fmt.Println("Commands: add, remove, list")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "add":
        addCmd := flag.NewFlagSet("add", flag.ExitOnError)
        name := addCmd.String("name", "", "item name")
        addCmd.Parse(os.Args[2:])
        fmt.Printf("Added: %s\n", *name)

    case "remove":
        removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
        id := removeCmd.Int("id", 0, "item id")
        removeCmd.Parse(os.Args[2:])
        fmt.Printf("Removed ID: %d\n", *id)

    case "list":
        listCmd := flag.NewFlagSet("list", flag.ExitOnError)
        sortBy := listCmd.String("sort", "name", "sort by: name|date")
        listCmd.Parse(os.Args[2:])
        fmt.Printf("List (sorted by %s)\n", *sortBy)

    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}
```

**Uso:**
```bash
$ go run main.go add -name "Task 1"
Added: Task 1

$ go run main.go list -sort date
List (sorted by date)

$ go run main.go remove -id 5
Removed ID: 5
```

### 51.2.5 Limitaciones del Flag Estándar

| Limitación | Impacto | Solución |
|-----------|---------|----------|
| Sin subcomandos jerarquizados | Difícil de escalar | Usar Cobra/Urfave |
| Sin auto-generated help | Help texto manual | Usar Cobra (auto help) |
| Sin validación de flags | Datos inválidos | Validar manualmente |
| Sin bash completion | UX pobre | Usar Cobra + completions |
| Sin command aliasing | Repetición de código | Urfave permite aliases |
| Parsing manual complejo | Propenso a errores | Cobra maneja complexity |

**Cuándo usar Flag estándar:**
- ✅ Herramientas muy simples (< 3 subcomandos)
- ✅ Máxima portabilidad (sin dependencias externas)
- ✅ Proyectos internos con pocos usuarios
- ❌ Herramientas complejas (kubectl, docker)
- ❌ Cuando necesitas UX profesional

---

## 51.3 Cobra Framework

### 51.3.1 Arquitectura de Cobra

**Estructura jerárquica:**

```
┌──────────────────────────────────────────┐
│         ROOT COMMAND                      │
│   (el programa principal: app)            │
├──────────────────────────────────────────┤
│                                           │
│  ├─ Persistent Flags (heredados por hijos)
│  └─ Action: Run()                        │
│                                           │
│      ├─ SUBCOMMAND 1 (app serve)          │
│      │   ├─ Local Flags                   │
│      │   └─ Action: Run()                 │
│      │                                    │
│      ├─ SUBCOMMAND 2 (app config)         │
│      │   ├─ Local Flags                   │
│      │   └─ Action: Run()                 │
│      │                                    │
│      └─ SUBCOMMAND 3 (app admin)          │
│          ├─ Sub-subcommand (app admin add)
│          └─ Action: Run()                 │
│                                           │
└──────────────────────────────────────────┘
```

### 51.3.2 Setup Básico de Cobra

**Instalación:**
```bash
go get -u github.com/spf13/cobra@latest
go install github.com/spf13/cobra-cli@latest
```

**Inicializar proyecto:**
```bash
cobra-cli init myapp
cd myapp
cobra-cli add serve
cobra-cli add config
cobra-cli add admin
cobra-cli add admin add  # sub-subcommand
```

**Estructura generada:**
```
myapp/
├── cmd/
│   ├── root.go      # Comando raíz
│   ├── serve.go     # Subcomando: serve
│   ├── config.go    # Subcomando: config
│   ├── admin.go     # Subcomando: admin
│   └── adminAdd.go  # Sub-subcomando: admin add
├── main.go
└── go.mod
```

**Ejemplo completo - main.go:**

```go
package main

import "myapp/cmd"

func main() {
    cmd.Execute()
}
```

**Ejemplo completo - cmd/root.go:**

```go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    verbose bool
    config  string
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "MyApp - A sample CLI application",
    Long:  `MyApp is a demonstration of Cobra framework for building CLI applications.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("MyApp v1.0.0")
        if verbose {
            fmt.Println("Verbose mode enabled")
            fmt.Printf("Config file: %s\n", config)
        }
    },
}

func init() {
    // Persistent flags (heredados por subcomandos)
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
        "enable verbose output")
    rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "config.yaml",
        "config file path")
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**Ejemplo completo - cmd/serve.go:**

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var (
    port int
    host string
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the server",
    Long:  `Start the application server on the specified host and port.`,
    PreRun: func(cmd *cobra.Command, args []string) {
        if verbose {
            fmt.Printf("Pre-run hook: preparing server (config: %s)\n", config)
        }
    },
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Starting server on %s:%d\n", host, port)
        if verbose {
            fmt.Println("Verbose: server started successfully")
        }
    },
    PostRun: func(cmd *cobra.Command, args []string) {
        if verbose {
            fmt.Println("Post-run hook: cleaning up")
        }
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)

    // Local flags (específicos de este comando)
    serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "server port")
    serveCmd.Flags().StringVarP(&host, "host", "h", "localhost", "server host")

    // Requerir flag
    serveCmd.MarkFlagRequired("port")
}
```

**Uso:**
```bash
$ go run main.go serve -p 9000 -h 0.0.0.0
Starting server on 0.0.0.0:9000

$ go run main.go -v serve -p 9000
Verbose mode enabled
Config file: config.yaml
Pre-run hook: preparing server (config: config.yaml)
Starting server on localhost:9000
Verbose: server started successfully
Post-run hook: cleaning up
```

### 51.3.3 Command Hooks (Pre/Post/Persistent)

Cobra proporciona múltiples puntos de ejecución:

```go
package cmd

import (
    "fmt"
    "log"

    "github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
    Use:   "setup",
    Short: "Setup command with all hooks",

    // Se ejecuta si el comando es válido pero ANTES del Run
    PreRunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("1. PreRunE - validar argumentos, inicializar")
        if len(args) == 0 {
            return fmt.Errorf("se requiere al menos un argumento")
        }
        return nil
    },

    // Alternativa no-error de PreRunE
    PreRun: func(cmd *cobra.Command, args []string) {
        fmt.Println("2. PreRun")
    },

    // Acción principal del comando
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("3. Run - procesando: %v\n", args)
    },

    // Se ejecuta DESPUÉS del Run (incluso si hay error)
    PostRun: func(cmd *cobra.Command, args []string) {
        fmt.Println("4. PostRun")
    },

    // Versión con manejo de errores
    PostRunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("5. PostRunE - cleanup")
        return nil
    },
}

// PersistentPreRun: se ejecuta en el comando actual y TODOS los hijos
var parentCmd = &cobra.Command{
    Use:   "parent",
    Short: "Parent command",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        fmt.Println("[Parent] PersistentPreRun - ejecutado antes que todos los hijos")
    },
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("[Parent] Run")
    },
}

var childCmd = &cobra.Command{
    Use:   "child",
    Short: "Child command",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("[Child] Run")
    },
}

func init() {
    rootCmd.AddCommand(setupCmd)
    rootCmd.AddCommand(parentCmd)
    parentCmd.AddCommand(childCmd)
}
```

**Orden de ejecución:**
```
parent child:
1. [Parent] PersistentPreRun  <- ejecutado primero
2. [Child] Run               <- ejecutado después
```

### 51.3.4 Context Propagation

Pasar datos entre comando raíz y subcomandos:

```go
package cmd

import (
    "context"
    "fmt"

    "github.com/spf13/cobra"
)

// Tipos personalizados para context
type AppContext struct {
    ConfigFile string
    Verbose    bool
    LogFile    string
}

var appCtx *AppContext

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "App with context",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Inicializar contexto raíz
        appCtx = &AppContext{
            ConfigFile: config,
            Verbose:    verbose,
            LogFile:    "/tmp/app.log",
        }
        if verbose {
            fmt.Printf("Context initialized: %+v\n", appCtx)
        }
    },
}

var dbCmd = &cobra.Command{
    Use:   "db",
    Short: "Database operations",
    Run: func(cmd *cobra.Command, args []string) {
        // Acceder al contexto desde el subcomando
        if appCtx == nil {
            fmt.Println("Error: context not initialized")
            return
        }
        fmt.Printf("DB operation with config: %s\n", appCtx.ConfigFile)
        fmt.Printf("Log file: %s\n", appCtx.LogFile)
    },
}

func init() {
    rootCmd.AddCommand(dbCmd)
    rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "config.yaml", "config file")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
}
```

### 51.3.5 Validación y Error Handling

```go
package cmd

import (
    "fmt"
    "os"
    "strconv"

    "github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
    Use:   "deploy <env> <version>",
    Short: "Deploy application",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) != 2 {
            return fmt.Errorf("requires exactly 2 arguments (env, version)")
        }
        // Validar primer argumento
        validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
        if !validEnvs[args[0]] {
            return fmt.Errorf("invalid environment: %s (must be dev|staging|prod)", args[0])
        }
        // Validar segundo argumento
        if _, err := strconv.Atoi(args[1]); err != nil {
            return fmt.Errorf("version must be a number, got: %s", args[1])
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Deploying v%s to %s\n", args[1], args[0])
    },
}

func init() {
    rootCmd.AddCommand(deployCmd)
}
```

**Uso:**
```bash
$ go run main.go deploy prod 2
Deploying v2 to prod

$ go run main.go deploy unknown 1
Error: invalid environment: unknown (must be dev|staging|prod)

$ go run main.go deploy prod abc
Error: version must be a number, got: abc
```

### 51.3.6 Cobra Completions

Auto-completar en bash, zsh, fish:

```go
// Generar completions
$ app completion bash > /etc/bash_completion.d/app

// En ~/.bashrc o ~/.zshrc
source /etc/bash_completion.d/app

// Uso
$ app [TAB][TAB]  # Muestra: serve deploy config admin
$ app serve -[TAB][TAB]  # Muestra: --port --host --help
```

Cobra genera automáticamente completions. Para customize:

```go
var serveCmd = &cobra.Command{
    Use: "serve",
    // ...
    ValidArgs: []string{"fast", "slow"},  // Argumentos válidos
}

func init() {
    serveCmd.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string,
        toComplete string) ([]string, cobra.ShellCompDirective) {
        return []string{"config.yaml", "config.prod.yaml"}, cobra.ShellCompDirectiveDefault
    })
}
```

---

## 51.4 Urfave/CLI v2

### 51.4.1 Alternativa Minimalista a Cobra

Urfave es más ligero que Cobra, ideal para CLIs simples-medianas:

```bash
go get github.com/urfave/cli/v2
```

**Comparativa Cobra vs Urfave:**

| Aspecto | Cobra | Urfave |
|--------|-------|--------|
| Tamaño de import | ~500 KB | ~100 KB |
| Curva de aprendizaje | Media-Alta | Baja |
| Complejidad | Alta (OOP) | Baja (Funcional) |
| Subcomandos | Jerárquicos | Simples |
| Auto-completions | Sí (integrado) | Manual |
| Comunidad | Grande (kubectl) | Mediana |
| Mejor para | Apps grandes | Apps medianas |

### 51.4.2 Setup Básico

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "myapp",
        Usage: "A simple CLI app",
        Version: "1.0.0",

        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:    "verbose",
                Aliases: []string{"v"},
                Usage:   "enable verbose output",
            },
            &cli.StringFlag{
                Name:    "config",
                Aliases: []string{"c"},
                Usage:   "config file path",
                Value:   "config.yaml",
            },
        },

        Action: func(cCtx *cli.Context) error {
            if cCtx.Bool("verbose") {
                fmt.Println("Verbose mode enabled")
            }
            fmt.Printf("Config: %s\n", cCtx.String("config"))
            return nil
        },

        Commands: []*cli.Command{
            {
                Name:    "serve",
                Aliases: []string{"s"},
                Usage:   "start the server",
                Flags: []cli.Flag{
                    &cli.IntFlag{
                        Name:    "port",
                        Aliases: []string{"p"},
                        Value:   8080,
                        Usage:   "server port",
                    },
                },
                Action: func(cCtx *cli.Context) error {
                    fmt.Printf("Starting server on port %d\n", cCtx.Int("port"))
                    return nil
                },
            },
            {
                Name:  "config",
                Usage: "manage configuration",
                Subcommands: []*cli.Command{
                    {
                        Name:  "show",
                        Usage: "show current config",
                        Action: func(cCtx *cli.Context) error {
                            fmt.Println("Showing configuration...")
                            return nil
                        },
                    },
                    {
                        Name:  "edit",
                        Usage: "edit config",
                        Action: func(cCtx *cli.Context) error {
                            fmt.Println("Editing configuration...")
                            return nil
                        },
                    },
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
```

**Uso:**
```bash
$ go run main.go -v serve -p 9000
Verbose mode enabled
Config: config.yaml
Starting server on port 9000

$ go run main.go config show
Showing configuration...
```

### 51.4.3 Action Functions y Contexto

```go
package main

import (
    "fmt"
    "os"

    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name: "actions-demo",
        Commands: []*cli.Command{
            {
                Name: "deploy",
                Flags: []cli.Flag{
                    &cli.StringFlag{Name: "env"},
                    &cli.StringFlag{Name: "version"},
                },
                // Action recibe *cli.Context con acceso a flags, comandos, args
                Action: func(cCtx *cli.Context) error {
                    env := cCtx.String("env")
                    version := cCtx.String("version")
                    args := cCtx.Args().Slice()

                    fmt.Printf("Deploying v%s to %s\n", version, env)
                    if len(args) > 0 {
                        fmt.Printf("Extra args: %v\n", args)
                    }
                    return nil
                },
            },
            {
                Name: "process",
                Flags: []cli.Flag{
                    &cli.StringSliceFlag{Name: "tag"},
                    &cli.IntFlag{Name: "workers", Value: 4},
                },
                Action: func(cCtx *cli.Context) error {
                    tags := cCtx.StringSlice("tag")
                    workers := cCtx.Int("workers")

                    fmt.Printf("Processing with %d workers\n", workers)
                    fmt.Printf("Tags: %v\n", tags)
                    return nil
                },
            },
        },
    }

    app.Run(os.Args)
}
```

**Uso:**
```bash
$ go run main.go deploy --env prod --version 2.1 file.txt
Deploying v2.1 to prod
Extra args: [file.txt]

$ go run main.go process --tag critical --tag urgent --workers 8
Processing with 8 workers
Tags: [critical urgent]
```

### 51.4.4 Categorías y Grupos

Organizar comandos por categorías:

```go
package main

import (
    "fmt"
    "os"

    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name: "myapp",
        Commands: []*cli.Command{
            // Categoría: Database
            {
                Name:     "db:migrate",
                Category: "Database",
                Usage:    "run database migrations",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("Running migrations...")
                    return nil
                },
            },
            {
                Name:     "db:seed",
                Category: "Database",
                Usage:    "seed database with test data",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("Seeding database...")
                    return nil
                },
            },
            // Categoría: Server
            {
                Name:     "serve",
                Category: "Server",
                Usage:    "start server",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("Starting server...")
                    return nil
                },
            },
            {
                Name:     "serve:reload",
                Category: "Server",
                Usage:    "reload server",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("Reloading server...")
                    return nil
                },
            },
        },
    }

    app.Run(os.Args)
}
```

**Ayuda resultante:**
```
COMMANDS:
  Database:
    db:migrate  run database migrations
    db:seed     seed database with test data

  Server:
    serve       start server
    serve:reload reload server
```

### 51.4.5 Custom Help y Before/After

```go
package main

import (
    "fmt"
    "os"

    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name: "myapp",
        Usage: "My Application",
        HelpName: "myapp",

        // Ejecutado ANTES de procesar comandos
        Before: func(cCtx *cli.Context) error {
            fmt.Println("[BEFORE] Initializing app...")
            return nil
        },

        // Ejecutado DESPUÉS de procesar comandos
        After: func(cCtx *cli.Context) error {
            fmt.Println("[AFTER] Cleanup...")
            return nil
        },

        // Personalizar help text
        CustomAppHelpTemplate: `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} [OPTIONS] COMMAND [ARGS]

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}

COMMANDS:
   {{range .VisibleCommands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}

EXAMPLES:
   myapp serve --port 8080
   myapp config show --json
`,

        Commands: []*cli.Command{
            {
                Name: "serve",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("Serving...")
                    return nil
                },
            },
        },
    }

    app.Run(os.Args)
}
```

---

## 51.5 Input/Output Avanzado

### 51.5.1 Flag Parsing Avanzado

Combinar flags de múltiples formas:

```go
package main

import (
    "flag"
    "fmt"
    "log"
)

func main() {
    // Flags booleanos pueden ser usados de varias formas
    flag.Bool("debug", false, "debug mode")
    // Uso: -debug, -debug=true, -debug=false
    
    // String flags
    flag.String("output", "-", "output file (- = stdout)")
    // Uso: -output file.txt, -output=file.txt
    
    // Short form
    flag.String("o", "", "output file (shorthand)")
    // Uso: -o file.txt
    
    flag.Parse()

    // Acceder a argumentos posicionales
    fmt.Println("Positional args:", flag.Args())
}
```

**Uso avanzado:**
```bash
$ go run main.go -debug -output result.txt arg1 arg2
$ go run main.go -debug=true -o=result.txt
$ go run main.go -- -not-a-flag file.txt
```

### 51.5.2 Stdin/Stdout/Stderr

Leer y escribir en streams:

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

func main() {
    // === LECTURA ===
    
    // Leer línea completa desde stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter your name: ")
    name, _ := reader.ReadString('\n')
    fmt.Printf("Hello, %s", name)

    // Leer todo stdin
    input, _ := io.ReadAll(os.Stdin)
    fmt.Printf("Read %d bytes\n", len(input))

    // Leer línea a línea
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        if line == "quit" {
            break
        }
        fmt.Println("You said:", line)
    }

    // === ESCRITURA ===

    // Stdout
    fmt.Println("Normal output")
    fmt.Fprintf(os.Stdout, "To stdout\n")

    // Stderr (para errores/logging)
    fmt.Fprintf(os.Stderr, "Error message\n")
    fmt.Fprintln(os.Stderr, "Another error")

    // === REDIRECCIÓN ===
    
    // Abrir archivo para escritura
    f, _ := os.Create("output.txt")
    defer f.Close()

    // Escribir a archivo
    fmt.Fprintf(f, "Content to file\n")

    // Copy stdin a stdout (pipeline)
    io.Copy(os.Stdout, os.Stdin)
}
```

**Uso con pipes:**
```bash
$ echo "Hello" | go run main.go  # stdin -> programa -> stdout
$ go run main.go > output.txt    # stdout -> archivo
$ go run main.go 2> error.txt    # stderr -> archivo
$ go run main.go < input.txt     # archivo -> stdin
$ cmd1 | go run main.go | cmd2   # pipeline
```

### 51.5.3 Interactive Prompts

Usar `promptui` para input interactivo:

```bash
go get github.com/manifoldco/promptui
```

```go
package main

import (
    "fmt"

    "github.com/manifoldco/promptui"
)

func main() {
    // === PROMPT SIMPLE ===
    prompt := promptui.Prompt{
        Label: "Enter your name",
        Default: "Guest",
    }
    name, _ := prompt.Run()
    fmt.Printf("Hello, %s!\n", name)

    // === SELECT ===
    items := []string{"Development", "Staging", "Production"}
    selectPrompt := promptui.Select{
        Label: "Select environment",
        Items: items,
    }
    idx, env, _ := selectPrompt.Run()
    fmt.Printf("Selected index %d: %s\n", idx, env)

    // === CONFIRM ===
    confirmPrompt := promptui.Prompt{
        Label:     "Continue",
        IsConfirm: true,
    }
    result, _ := confirmPrompt.Run()
    if result == "y" {
        fmt.Println("Continuing...")
    }

    // === PASSWORD ===
    passwordPrompt := promptui.Prompt{
        Label: "Enter password",
        Mask:  '*',
    }
    password, _ := passwordPrompt.Run()
    fmt.Printf("Password length: %d\n", len(password))

    // === VALIDACIÓN ===
    validatePrompt := promptui.Prompt{
        Label: "Enter age",
        Validate: func(s string) error {
            if len(s) == 0 {
                return fmt.Errorf("age cannot be empty")
            }
            if _, err := fmt.Sscanf(s, "%d", new(int)); err != nil {
                return fmt.Errorf("not a valid number")
            }
            return nil
        },
    }
    age, _ := validatePrompt.Run()
    fmt.Printf("Age: %s\n", age)
}
```

**Ejecución interactiva:**
```
$ go run main.go
Enter your name: Alice
Hello, Alice!

Select environment:
  > Development
    Staging
    Production
Selected index 0: Development

Continue? (y/N) y
Continuing...

Enter password: ****
Password length: 8

Enter age: abc
Error: not a valid number
Enter age: 30
Age: 30
```

### 51.5.4 Progress Bars

Mostrar progreso de operaciones largas:

```bash
go get github.com/cheggaaa/pb/v3
```

```go
package main

import (
    "time"

    "github.com/cheggaaa/pb/v3"
)

func main() {
    // === PROGRESS BAR SIMPLE ===
    count := 100
    bar := pb.StartNew(count)
    for i := 0; i < count; i++ {
        bar.Increment()
        time.Sleep(50 * time.Millisecond)
    }
    bar.Finish()

    // === CON DESCRIPCIÓN ===
    bar = pb.Full.Start(50)
    bar.SetCurrent(0)
    for i := 0; i < 50; i++ {
        bar.Increment()
        time.Sleep(100 * time.Millisecond)
    }
    bar.Finish()

    // === MÚLTIPLES BARRAS ===
    bars := make([]*pb.ProgressBar, 3)
    for i := 0; i < 3; i++ {
        bars[i] = pb.StartNew(100)
    }

    for i := 0; i < 100; i++ {
        for _, bar := range bars {
            bar.Increment()
        }
        time.Sleep(50 * time.Millisecond)
    }

    for _, bar := range bars {
        bar.Finish()
    }

    // === CON FORMATO PERSONALIZADO ===
    tmpl := `{{ string . "prefix" }} {{ counters . }} {{ bar . "[" "=" (cycle . ">" ) " " "]" }} {{ percent . }} {{ string . "suffix" }}`
    bar = pb.ProgressBarTemplate(tmpl).Start(100)
    bar.Set("prefix", "Processing: ")
    bar.Set("suffix", " [remaining: " + bar.String() + "]")
    for i := 0; i < 100; i++ {
        bar.Increment()
        time.Sleep(50 * time.Millisecond)
    }
    bar.Finish()
}
```

**Salida:**
```
100 / 100 [████████████████████████████] 100 %

Processing:  25 / 100 [=====> ] 25 %
```

### 51.5.5 Spinners y Loaders

```bash
go get github.com/briandowns/spinner
```

```go
package main

import (
    "fmt"
    "time"

    "github.com/briandowns/spinner"
)

func main() {
    // === SPINNER SIMPLE ===
    s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
    s.Start()
    time.Sleep(3 * time.Second)
    s.Stop()
    fmt.Println("Done!")

    // === CON SUFIJO ===
    s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
    s.Suffix = " Loading... "
    s.Start()
    time.Sleep(2 * time.Second)
    s.Stop()
    fmt.Println("Finished!")

    // === CON COLOR ===
    s = spinner.New(spinner.CharSets[39], 100*time.Millisecond)
    s.Color("blue")
    s.Suffix = " Processing "
    s.Start()
    time.Sleep(2 * time.Second)
    s.Stop()
    fmt.Println("Complete!")
}
```

**Salida:**
```
⠙ Loading...  [después de 3s] Done!
⣾ Processing  [después de 2s] Finished!
🟦 Processing  [después de 2s] Complete!
```

---

## 51.6 Formatting Output

### 51.6.1 Plain Text Formatting

```go
package main

import "fmt"

func main() {
    // Ancho de columna
    fmt.Printf("%-20s | %10d | %8.2f\n", "Name", 42, 3.14)
    fmt.Printf("%-20s | %10d | %8.2f\n", "Product XYZ", 100, 29.99)

    // Padding y alineación
    fmt.Printf("%5d\n", 42)    // "   42" (derecha)
    fmt.Printf("%-5d\n", 42)   // "42   " (izquierda)
    fmt.Printf("%05d\n", 42)   // "00042" (ceros)

    // Strings
    fmt.Printf("%20s\n", "hello")   // "               hello" (derecha)
    fmt.Printf("%-20s\n", "hello")  // "hello               " (izquierda)
    fmt.Printf("%20.5s\n", "hello")  // "               hello" (truncado)

    // Valores booleanos
    fmt.Printf("%5v | %t\n", true, false)

    // Tipos (debugging)
    fmt.Printf("%T\n", 42)           // "int"
    fmt.Printf("%v\n", []int{1,2})   // "[1 2]"
    fmt.Printf("%#v\n", []int{1,2})  // "[]int{1, 2}"

    // Escaping
    fmt.Printf("%q\n", "hello\nworld")  // "hello\\nworld"
}
```

### 51.6.2 Tablas con Tabwriter

```go
package main

import (
    "fmt"
    "os"
    "text/tabwriter"
)

func main() {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

    // Encabezados
    fmt.Fprintln(w, "NAME\tAGE\tCITY\tOCCUPATION")
    fmt.Fprintln(w, "----\t---\t----\t----------")

    // Datos (con \t como separador)
    fmt.Fprintln(w, "Alice\t30\tNew York\tSoftware Engineer")
    fmt.Fprintln(w, "Bob\t25\tSan Francisco\tData Scientist")
    fmt.Fprintln(w, "Charlie\t35\tLondon\tDevOps Engineer")

    w.Flush()
}
```

**Salida:**
```
NAME     AGE  CITY           OCCUPATION
----     ---  ----           ----------
Alice    30   New York       Software Engineer
Bob      25   San Francisco  Data Scientist
Charlie  35   London         DevOps Engineer
```

### 51.6.3 Tablas Avanzadas

```bash
go get github.com/olekukonko/tablewriter
```

```go
package main

import (
    "os"

    "github.com/olekukonko/tablewriter"
)

func main() {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "Age", "City"})

    // Alineación
    table.SetColumnAlignment([]int{
        tablewriter.ALIGN_LEFT,
        tablewriter.ALIGN_CENTER,
        tablewriter.ALIGN_RIGHT,
    })

    // Datos
    table.Append([]string{"Alice", "30", "New York"})
    table.Append([]string{"Bob", "25", "London"})

    // Estilo
    table.SetBorders(tablewriter.Border{
        Left:   true,
        Right:  true,
        Top:    true,
        Bottom: true,
    })

    // Líneas automáticas
    table.SetAutoMergeCells(true)

    table.Render()
}
```

**Salida:**
```
+-------+-----+-----------+
| Name  | Age |      City |
+-------+-----+-----------+
| Alice |  30 |  New York |
| Bob   |  25 |    London |
+-------+-----+-----------+
```

### 51.6.4 JSON Output

Salida estructurada y filtrable:

```go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
)

type Server struct {
    Name    string   `json:"name"`
    Host    string   `json:"host"`
    Port    int      `json:"port"`
    Enabled bool     `json:"enabled"`
    Tags    []string `json:"tags"`
}

func main() {
    outputFormat := flag.String("format", "table", "output format: table|json")
    flag.Parse()

    servers := []Server{
        {"api-1", "10.0.1.1", 8080, true, []string{"prod", "api"}},
        {"api-2", "10.0.1.2", 8080, true, []string{"prod", "api"}},
        {"cache", "10.0.2.1", 6379, true, []string{"prod", "cache"}},
    }

    if *outputFormat == "json" {
        // Pretty JSON
        data, _ := json.MarshalIndent(servers, "", "  ")
        fmt.Println(string(data))
    } else if *outputFormat == "table" {
        // Tabla
        fmt.Printf("%-10s | %-12s | %-6s | %-8s | %s\n",
            "Name", "Host", "Port", "Enabled", "Tags")
        fmt.Println(string([]byte{'-', '-', '-', '-', '-', '-', '-', '-'}))
        for _, s := range servers {
            tags := ""
            for _, t := range s.Tags {
                tags += t + ","
            }
            fmt.Printf("%-10s | %-12s | %-6d | %-8v | %s\n",
                s.Name, s.Host, s.Port, s.Enabled, tags)
        }
    }
}
```

**Uso:**
```bash
$ go run main.go
Name       | Host         | Port   | Enabled  | Tags
--------

api-1      | 10.0.1.1     | 8080   | true     | prod,api,
api-2      | 10.0.1.2     | 8080   | true     | prod,api,
cache      | 10.0.2.1     | 6379   | true     | prod,cache,

$ go run main.go --format json
[
  {
    "name": "api-1",
    "host": "10.0.1.1",
    "port": 8080,
    "enabled": true,
    "tags": ["prod", "api"]
  },
  ...
]
```

### 51.6.5 Colores y Styling

```bash
go get github.com/fatih/color
go get github.com/charmbracelet/lipgloss
```

```go
package main

import (
    "fmt"

    "github.com/fatih/color"
    "github.com/charmbracelet/lipgloss"
)

func main() {
    // === FATIH/COLOR ===
    color.Red("Error: something went wrong")
    color.Green("Success: operation completed")
    color.Yellow("Warning: check this")
    color.Blue("Info: additional details")

    // Con formato
    color.Cyan("Deployed to: %s\n", "production")

    // Atributos
    bold := color.New(color.Bold)
    bold.Println("This is bold")

    underline := color.New(color.Underline)
    underline.Println("This is underlined")

    // Combinado
    multiColor := color.New(color.FgCyan, color.Bold, color.BgWhite)
    multiColor.Println("Cyan, bold, white background")

    // === LIPGLOSS (más sofisticado) ===
    
    // Estilos
    header := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        PaddingLeft(2).
        PaddingRight(2)

    success := lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Bold(true)

    error := lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true)

    fmt.Println(header.Render("DEPLOYMENT"))
    fmt.Println(success.Render("✓ Service deployed"))
    fmt.Println(error.Render("✗ Database connection failed"))

    // Layout
    left := lipgloss.NewStyle().
        Width(20).
        Align(lipgloss.Left).
        Render("Commands:")

    right := lipgloss.NewStyle().
        Width(30).
        Align(lipgloss.Left).
        Render("serve, config, admin")

    fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
}
```

**Salida:**
```
[MAGENTA] Error: something went wrong
[GREEN] Success: operation completed
[CYAN] Deployed to: production

┌────────────────────┐
│   DEPLOYMENT       │
└────────────────────┘
✓ Service deployed
✗ Database connection failed

Commands:          Commands: serve, config, admin
```

---

## 51.7 Configuration Files

### 51.7.1 Integración Viper con CLI

Viper proporciona jerarquía de configuración:

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/spf13/viper"
)

func main() {
    // Flags (máxima prioridad)
    configFile := flag.String("config", "config.yaml", "config file")
    logLevel := flag.String("log-level", "", "log level (debug|info|warn|error)")
    flag.Parse()

    // Viper setup
    if *configFile != "" {
        viper.SetConfigFile(*configFile)
    } else {
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
        viper.AddConfigPath(".")
        viper.AddConfigPath("/etc/myapp")
        viper.AddConfigPath(os.ExpandEnv("$HOME/.config/myapp"))
    }

    // Environment variables override
    viper.SetEnvPrefix("MYAPP")
    viper.AutomaticEnv()

    // Default values
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.host", "localhost")
    viper.SetDefault("database.driver", "postgres")
    viper.SetDefault("logging.level", "info")

    // Cargar configuración
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            log.Fatal("Error reading config:", err)
        }
        fmt.Println("No config file found, using defaults")
    }

    // Flags sobreescriben config file
    if *logLevel != "" {
        viper.Set("logging.level", *logLevel)
    }

    // Acceder a valores
    fmt.Printf("Server: %s:%d\n",
        viper.GetString("server.host"),
        viper.GetInt("server.port"))
    fmt.Printf("Database: %s\n", viper.GetString("database.driver"))
    fmt.Printf("Log Level: %s\n", viper.GetString("logging.level"))

    // Deserializar a struct
    var config AppConfig
    if err := viper.Unmarshal(&config); err != nil {
        log.Fatal("Error unmarshaling config:", err)
    }
    fmt.Printf("Config: %+v\n", config)
}

type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        Driver string `mapstructure:"driver"`
    } `mapstructure:"database"`
    Logging struct {
        Level string `mapstructure:"level"`
    } `mapstructure:"logging"`
}
```

**config.yaml:**
```yaml
server:
  host: 0.0.0.0
  port: 8080
database:
  driver: postgres
  host: localhost
  port: 5432
logging:
  level: debug
```

**Precedencia (mayor a menor):**
1. Flags CLI: `--log-level debug`
2. Environment variables: `MYAPP_LOGGING_LEVEL=debug`
3. Config file: `config.yaml`
4. Valores por defecto: `viper.SetDefault()`

### 51.7.2 Config File Hierarchy

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/viper"
)

func initConfig() error {
    // Buscar config en múltiples ubicaciones
    configPaths := []string{
        "./config.yaml",                                  // Local
        os.ExpandEnv("$HOME/.config/myapp/config.yaml"), // User home
        "/etc/myapp/config.yaml",                         // System
    }

    found := false
    for _, path := range configPaths {
        if _, err := os.Stat(path); err == nil {
            viper.SetConfigFile(path)
            fmt.Printf("Using config: %s\n", path)
            found = true
            break
        }
    }

    if !found {
        fmt.Println("No config file found, using defaults")
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
        viper.AddConfigPath(".")
    }

    viper.SetEnvPrefix("APP")
    viper.AutomaticEnv()

    viper.SetDefault("debug", false)
    viper.SetDefault("workers", 4)

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return err
        }
    }

    return nil
}

func main() {
    if err := initConfig(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Debug: %v\n", viper.GetBool("debug"))
    fmt.Printf("Workers: %d\n", viper.GetInt("workers"))
}
```

**Búsqueda automática:**
```bash
# 1. Intenta ~/.config/myapp/config.yaml
$ myapp

# 2. Si existe APP_DEBUG=true, lo sobreescribe
$ APP_DEBUG=true myapp

# 3. Si pasas --config, usa ese archivo
$ myapp --config /custom/path.yaml
```

---

## 51.8 Error Handling y Logging

### 51.8.1 Errores con Contexto

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "os"
)

// Errores personalizados
type CLIError struct {
    Code    int
    Message string
    Err     error
}

func (e *CLIError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("CLI Error [%d]: %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("CLI Error [%d]: %s", e.Code, e.Message)
}

// Exitcodes estándar
const (
    ExitOK            = 0
    ExitGeneralError  = 1
    ExitMisuse        = 2
    ExitFileNotFound  = 3
    ExitValidation    = 4
    ExitUnavailable   = 5
    ExitInternalError = 6
)

func processFile(filename string) error {
    if filename == "" {
        return &CLIError{
            Code:    ExitValidation,
            Message: "filename is required",
        }
    }

    if _, err := os.Stat(filename); err != nil {
        return &CLIError{
            Code:    ExitFileNotFound,
            Message: fmt.Sprintf("file not found: %s", filename),
            Err:     err,
        }
    }

    return nil
}

func main() {
    if err := processFile(""); err != nil {
        cliErr := err.(*CLIError)
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(cliErr.Code)
    }

    if err := processFile("nonexistent.txt"); err != nil {
        cliErr := err.(*CLIError)
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(cliErr.Code)
    }

    fmt.Println("OK")
}
```

**Uso:**
```bash
$ go run main.go
CLI Error [4]: filename is required
$ echo $?
4

$ go run main.go
CLI Error [3]: file not found: nonexistent.txt: ...
$ echo $?
3
```

### 51.8.2 Structured Logging

```bash
go get github.com/sirupsen/logrus
```

```go
package main

import (
    "fmt"
    "os"

    log "github.com/sirupsen/logrus"
)

func init() {
    // Formato JSON para parsed logging
    log.SetFormatter(&log.JSONFormatter{
        PrettyPrint: true,
    })

    // Salida a stderr
    log.SetOutput(os.Stderr)

    // Nivel por defecto
    log.SetLevel(log.InfoLevel)
}

func main() {
    // Logging simple
    log.Info("Application started")
    log.Warn("This is a warning")
    log.Error("An error occurred")

    // Con contexto (fields)
    log.WithFields(log.Fields{
        "user":     "alice",
        "action":   "deploy",
        "env":      "production",
        "duration": "2.5s",
    }).Info("Deployment completed")

    // Errores con stack trace
    err := fmt.Errorf("connection failed")
    log.WithError(err).Error("Database error")

    // Logger para operaciones específicas
    deployLogger := log.WithFields(log.Fields{
        "operation": "deploy",
        "version":   "2.1.0",
    })
    deployLogger.Info("Starting deployment")
    deployLogger.Info("Deployment successful")
}
```

**Salida:**
```json
{
  "level": "info",
  "msg": "Application started",
  "time": "2024-01-15T10:30:00Z"
}
{
  "level": "info",
  "msg": "Deployment completed",
  "time": "2024-01-15T10:30:05Z",
  "user": "alice",
  "action": "deploy",
  "env": "production",
  "duration": "2.5s"
}
```

### 51.8.3 Log Levels Configurables

```go
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    log "github.com/sirupsen/logrus"
)

var (
    verbose   bool
    debug     bool
    logFormat string
)

func setupLogging() {
    // Formato
    if logFormat == "json" {
        log.SetFormatter(&log.JSONFormatter{})
    } else {
        log.SetFormatter(&log.TextFormatter{
            FullTimestamp: true,
            ForceColors:   true,
        })
    }

    log.SetOutput(os.Stderr)

    // Nivel
    if debug {
        log.SetLevel(log.DebugLevel)
    } else if verbose {
        log.SetLevel(log.InfoLevel)
    } else {
        log.SetLevel(log.WarnLevel)
    }
}

var rootCmd = &cobra.Command{
    Use:   "app",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        setupLogging()
    },
    Run: func(cmd *cobra.Command, args []string) {
        log.Debug("Debug message (only with --debug)")
        log.Info("Info message (only with --verbose or --debug)")
        log.Warn("Warning message (always shown)")
        log.Error("Error message (always shown)")
    },
}

func init() {
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
        "enable verbose output")
    rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false,
        "enable debug output")
    rootCmd.PersistentFlags().StringVarP(&logFormat, "log-format", "", "text",
        "log format: text|json")
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**Uso:**
```bash
$ go run main.go
Warning message (always shown)
Error message (always shown)

$ go run main.go -v
Info message (only with --verbose or --debug)
Warning message (always shown)
Error message (always shown)

$ go run main.go -d --log-format json
{"level":"debug","msg":"Debug message...","time":"..."}
{"level":"info","msg":"Info message...","time":"..."}
```

---

## 51.9 Testing CLI

### 51.9.1 Captura de Stdout/Stderr

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "testing"
)

// Función a testear
func printGreeting(name string, w io.Writer) {
    fmt.Fprintf(w, "Hello, %s!\n", name)
}

// Test con captura de stdout
func TestGreeting(t *testing.T) {
    var buf bytes.Buffer
    printGreeting("Alice", &buf)

    expected := "Hello, Alice!\n"
    if buf.String() != expected {
        t.Errorf("got %q, want %q", buf.String(), expected)
    }
}

// Capturar stdout del programa completo
func captureStdout(f func()) string {
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    f()

    w.Close()
    os.Stdout = old

    var buf bytes.Buffer
    io.Copy(&buf, r)
    return buf.String()
}

func TestMainOutput(t *testing.T) {
    output := captureStdout(func() {
        printGreeting("Bob", os.Stdout)
    })

    if output != "Hello, Bob!\n" {
        t.Errorf("got %q", output)
    }
}
```

**Ejecutar:**
```bash
$ go test -v
=== RUN   TestGreeting
--- PASS: TestGreeting (0.00s)
=== RUN   TestMainOutput
--- PASS: TestMainOutput (0.00s)
```

### 51.9.2 Table-Driven Tests

```go
package main

import (
    "fmt"
    "strings"
    "testing"
)

func formatOutput(name string, age int) string {
    return fmt.Sprintf("%s is %d years old", name, age)
}

func TestFormatOutput(t *testing.T) {
    tests := []struct {
        name     string
        age      int
        expected string
    }{
        {"Alice", 30, "Alice is 30 years old"},
        {"Bob", 25, "Bob is 25 years old"},
        {"", 0, " is 0 years old"},
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%s_%d", tt.name, tt.age), func(t *testing.T) {
            got := formatOutput(tt.name, tt.age)
            if got != tt.expected {
                t.Errorf("got %q, want %q", got, tt.expected)
            }
        })
    }
}
```

**Ejecutar:**
```bash
$ go test -v -run TestFormatOutput
=== RUN   TestFormatOutput
=== RUN   TestFormatOutput/Alice_30
--- PASS: TestFormatOutput/Alice_30 (0.00s)
=== RUN   TestFormatOutput/Bob_25
--- PASS: TestFormatOutput/Bob_25 (0.00s)
=== RUN   TestFormatOutput/_0
--- PASS: TestFormatOutput/_0 (0.00s)
```

### 51.9.3 Testing Cobra Commands

```go
package cmd

import (
    "bytes"
    "testing"

    "github.com/spf13/cobra"
)

func TestServeCommand(t *testing.T) {
    tests := []struct {
        name      string
        args      []string
        wantError bool
        wantOutput string
    }{
        {
            name:        "default",
            args:        []string{"serve"},
            wantError:   false,
            wantOutput: "Starting server on localhost:8080\n",
        },
        {
            name:        "custom port",
            args:        []string{"serve", "--port", "9000"},
            wantError:   false,
            wantOutput: "Starting server on localhost:9000\n",
        },
        {
            name:        "invalid port",
            args:        []string{"serve", "--port", "invalid"},
            wantError:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := &cobra.Command{}
            out := &bytes.Buffer{}
            cmd.SetOut(out)
            cmd.SetArgs(tt.args)

            // Tu comando aquí
            // err := cmd.Execute()

            // if (err != nil) != tt.wantError {
            //     t.Errorf("got error %v, wantError %v", err, tt.wantError)
            // }

            // if out.String() != tt.wantOutput {
            //     t.Errorf("got %q, want %q", out.String(), tt.wantOutput)
            // }
        })
    }
}
```

### 51.9.4 Integration Tests

```go
package main

import (
    "os"
    "os/exec"
    "strings"
    "testing"
)

func TestCLIIntegration(t *testing.T) {
    tests := []struct {
        name        string
        args        []string
        wantCode    int
        wantInOutput string
    }{
        {
            name:         "help",
            args:         []string{"--help"},
            wantCode:     0,
            wantInOutput: "Usage:",
        },
        {
            name:         "serve default",
            args:         []string{"serve"},
            wantCode:     0,
            wantInOutput: "Starting server",
        },
        {
            name:     "invalid command",
            args:     []string{"unknown"},
            wantCode: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command("go", append([]string{"run", "main.go"}, tt.args...)...)
            out, err := cmd.CombinedOutput()

            gotCode := cmd.ProcessState.ExitCode()
            if gotCode != tt.wantCode {
                t.Errorf("exit code: got %d, want %d", gotCode, tt.wantCode)
            }

            if tt.wantInOutput != "" && !strings.Contains(string(out), tt.wantInOutput) {
                t.Errorf("output %q does not contain %q", string(out), tt.wantInOutput)
            }
        })
    }
}
```

---

## 51.10 Distribución y Packaging

### 51.10.1 Cross-Compilation

Build para múltiples plataformas:

```bash
# Linux x86_64
GOOS=linux GOARCH=amd64 go build -o myapp-linux-amd64

# macOS x86_64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o myapp-darwin-amd64

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o myapp-darwin-arm64

# Windows x86_64
GOOS=windows GOARCH=amd64 go build -o myapp-windows-amd64.exe

# Linux ARM64 (Raspberry Pi, etc.)
GOOS=linux GOARCH=arm64 go build -o myapp-linux-arm64

# Crear distribución
mkdir -p dist
for os in linux darwin windows; do
    for arch in amd64 arm64; do
        [[ "$os" == "windows" ]] && ext=".exe" || ext=""
        GOOS=$os GOARCH=$arch go build -o dist/myapp-${os}-${arch}${ext}
    done
done
```

### 51.10.2 Versionado con ldflags

Inyectar información de versión en binario:

```go
// main.go
package main

import (
    "fmt"
)

var (
    Version = "dev"
    Build   = ""
    Date    = ""
)

func main() {
    fmt.Printf("MyApp %s (build: %s, date: %s)\n", Version, Build, Date)
}
```

**Build con versión:**
```bash
#!/bin/bash
VERSION="1.0.0"
BUILD=$(git rev-parse --short HEAD)
DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

go build \
    -ldflags "-X main.Version=$VERSION -X main.Build=$BUILD -X main.Date=$DATE" \
    -o myapp

$ ./myapp
MyApp 1.0.0 (build: a1b2c3d, date: 2024-01-15T10:30:00Z)
```

### 51.10.3 GoReleaser

Automatizar releases:

```bash
go install github.com/goreleaser/goreleaser@latest
```

**.goreleaser.yaml:**
```yaml
project_name: myapp

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -X main.Version={{ .Version }}
      - -X main.Build={{ .Commit }}
      - -X main.Date={{ .Date }}

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

release:
  github:
    owner: myuser
    name: myapp

changelog:
  sort: asc
```

**Uso:**
```bash
# Taggear release
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0

# Build y upload
goreleaser release --clean
```

### 51.10.4 Homebrew Distribution

Distribuir via Homebrew (macOS):

```bash
# Crear tap (repositorio)
mkdir homebrew-myapp
cd homebrew-myapp
```

**Formula (myapp.rb):**
```ruby
class Myapp < Formula
  desc "My awesome CLI application"
  homepage "https://github.com/myuser/myapp"
  url "https://github.com/myuser/myapp/releases/download/v1.0.0/myapp-darwin-amd64.tar.gz"
  sha256 "abc123def456..."
  version "1.0.0"

  def install
    bin.install "myapp"
  end

  test do
    system "#{bin}/myapp", "--version"
  end
end
```

**Uso:**
```bash
# Agregar tap
brew tap myuser/myapp

# Instalar
brew install myapp

# Desinstalar
brew uninstall myapp
```

### 51.10.5 Container Distribution

Distribuir en Docker:

**Dockerfile (multi-stage):**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=1.0.0" \
    -o myapp

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
ENTRYPOINT ["./myapp"]
CMD ["serve", "--port", "8080"]
```

**Build y push:**
```bash
docker build -t myuser/myapp:1.0.0 .
docker tag myuser/myapp:1.0.0 myuser/myapp:latest
docker push myuser/myapp:1.0.0
docker push myuser/myapp:latest
```

---

## 51.11 Buenas Prácticas y Case Studies

### 51.11.1 Diseño de UX para CLI

**✅ BEST PRACTICES:**

```
┌──────────────────────────────────────────────────────┐
│           UX PRINCIPLES PARA CLI                      │
├──────────────────────────────────────────────────────┤
│ 1. INTUITIVIDAD                                       │
│    ✅ app server start                               │
│    ❌ app s st                                       │
│                                                      │
│ 2. CONSISTENCIA                                       │
│    ✅ app config get   / app config set              │
│    ❌ app config get   / app cfg set                │
│                                                      │
│ 3. VERBOSIDAD CONTROLABLE                             │
│    ✅ app serve (silencio si es exitoso)             │
│    ✅ app serve -v (info detallada)                 │
│    ❌ app serve (muestra todo, difícil de parsear)   │
│                                                      │
│ 4. ERRORES ÚTILES                                    │
│    ✅ "Error: port must be 1-65535, got: 99999"    │
│    ❌ "Error: invalid port"                         │
│                                                      │
│ 5. AYUDA PROGRESIVA (Progressive Disclosure)         │
│    ✅ app help          (resumen corto)              │
│    ✅ app help serve    (específico del comando)    │
│    ✅ app serve --help  (con ejemplos)              │
│                                                      │
│ 6. SALIDA SCRIPTEABLE                                │
│    ✅ app list --json  (para parsear en scripts)    │
│    ✅ app list         (para humans readable)       │
│    ❌ app list         (output ambiguo)            │
└──────────────────────────────────────────────────────┘
```

**Implementar ayuda con ejemplos:**

```go
var serveCmd = &cobra.Command{
    Use:   "serve [FLAGS]",
    Short: "Start the application server",
    Long: `Start the application server on the specified host and port.

The server listens for HTTP requests and routes them accordingly.
Configuration can be provided via flags, environment variables, or config file.`,

    Example: `  # Start with default settings (localhost:8080)
  app serve

  # Start on specific port
  app serve --port 9000

  # Start with custom host and port
  app serve --host 0.0.0.0 --port 3000

  # Start in debug mode
  app serve --debug`,

    Run: func(cmd *cobra.Command, args []string) {
        // ...
    },
}
```

**Salida de ayuda:**
```
$ app serve --help
Start the application server on the specified host and port.

The server listens for HTTP requests and routes them accordingly.
Configuration can be provided via flags, environment variables, or config file.

Usage:
  app serve [FLAGS]

Flags:
  -h, --host string      server host (default "localhost")
  -p, --port int         server port (default 8080)
      --debug            enable debug mode
      --help             show this help message

Examples:
  # Start with default settings (localhost:8080)
  app serve

  # Start on specific port
  app serve --port 9000
```

### 51.11.2 Documentación - Man Pages

Generar man pages automaticamente:

```bash
go install github.com/spf13/cobra-cli@latest
cobra-cli generate-docs --doc-dir=docs
```

**Uso:**
```bash
man app              # Página del programa
man app-serve        # Subcomando
man app-config-show  # Sub-subcomando
```

### 51.11.3 Anti-patterns ❌ vs Best Practices ✅

| Anti-pattern | Problema | Best Practice |
|---|---|---|
| ❌ Parsing manual `os.Args` | Propenso a errores, difícil de mantener | ✅ Usar Cobra/Urfave |
| ❌ Hardcoded paths | No portátil, quebrantable | ✅ Usar viper + config dirs |
| ❌ No capturar stderr | Información de error se pierde | ✅ Separate stderr + structured logging |
| ❌ Salida no estructurada | Imposible parsear en scripts | ✅ JSON output + human readable |
| ❌ No validar entrada | Comportamiento impredecible | ✅ Validación con Args, Validate |
| ❌ No capturar exit codes | Imposible encadenar en scripts | ✅ Retornar códigos semánticos |
| ❌ Dependencias globales | Testing imposible | ✅ Inyectar dependencias |
| ❌ Sin tests | Regressions silenciosas | ✅ Tests integración + unit |
| ❌ Mensajes de error genéricos | Usuario no sabe qué pasó | ✅ Mensajes específicos + contexto |
| ❌ Sin versión visible | Usuario no sabe qué versión está usando | ✅ `--version` con ldflags |

### 51.11.4 Case Study: kubectl Clone (Simplificado)

Pequeño clon de kubectl que demuestra patrones:

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var config string

var rootCmd = &cobra.Command{
    Use:   "kctl",
    Short: "Kubernetes control-like tool",
    Long:  `A simplified kubernetes-like CLI tool for demonstration.`,
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        viper.SetConfigFile(config)
        viper.SetDefault("cluster", "local")
        viper.ReadInConfig()
    },
}

var getCmd = &cobra.Command{
    Use:   "get <resource> [flags]",
    Short: "Display resources",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        resource := args[0]
        namespace, _ := cmd.Flags().GetString("namespace")
        fmt.Printf("Getting %s in namespace %s\n", resource, namespace)

        if output, _ := cmd.Flags().GetString("output"); output == "json" {
            fmt.Println(`[{"name": "resource1", "status": "running"}]`)
        } else {
            fmt.Println("NAME       STATUS")
            fmt.Println("resource1  running")
        }
    },
}

var applyCmd = &cobra.Command{
    Use:   "apply -f <file>",
    Short: "Apply configuration from file",
    RunE: func(cmd *cobra.Command, args []string) error {
        file, _ := cmd.Flags().GetString("filename")
        if file == "" {
            return fmt.Errorf("flag --filename must be specified")
        }
        fmt.Printf("Applying %s\n", file)
        return nil
    },
}

var deleteCmd = &cobra.Command{
    Use:   "delete <resource>",
    Short: "Delete resources",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Deleted %s\n", args[0])
    },
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "",
        "config file (default: ~/.kctl/config)")

    // Get command
    getCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace")
    getCmd.Flags().StringP("output", "o", "table", "Output format: table|json|yaml")
    rootCmd.AddCommand(getCmd)

    // Apply command
    applyCmd.Flags().StringP("filename", "f", "", "Filename to apply")
    rootCmd.AddCommand(applyCmd)

    // Delete command
    rootCmd.AddCommand(deleteCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**Uso:**
```bash
$ kctl get pods
Getting pods in namespace default
NAME       STATUS
resource1  running

$ kctl get pods -n prod -o json
[{"name": "resource1", "status": "running"}]

$ kctl apply -f deployment.yaml
Applying deployment.yaml

$ kctl delete pods/mypod
Deleted pods/mypod
```

### 51.11.5 Case Study: Task Manager CLI

CLI interactiva para gestionar tareas:

```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
)

type Task struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Done bool   `json:"done"`
}

var tasks []Task
var dataFile string

func init() {
    home, _ := os.UserHomeDir()
    dataFile = filepath.Join(home, ".tasks.json")
    loadTasks()
}

func loadTasks() {
    data, _ := ioutil.ReadFile(dataFile)
    json.Unmarshal(data, &tasks)
}

func saveTasks() {
    data, _ := json.MarshalIndent(tasks, "", "  ")
    ioutil.WriteFile(dataFile, data, 0644)
}

var addCmd = &cobra.Command{
    Use:   "add",
    Short: "Add a new task",
    Run: func(cmd *cobra.Command, args []string) {
        prompt := promptui.Prompt{
            Label: "Task name",
        }
        name, _ := prompt.Run()

        task := Task{
            ID:   len(tasks) + 1,
            Name: name,
            Done: false,
        }
        tasks = append(tasks, task)
        saveTasks()
        fmt.Printf("Added task #%d: %s\n", task.ID, task.Name)
    },
}

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all tasks",
    Run: func(cmd *cobra.Command, args []string) {
        if len(tasks) == 0 {
            fmt.Println("No tasks")
            return
        }
        for _, t := range tasks {
            status := " "
            if t.Done {
                status = "✓"
            }
            fmt.Printf("[%s] #%d: %s\n", status, t.ID, t.Name)
        }
    },
}

var doneCmd = &cobra.Command{
    Use:   "done <id>",
    Short: "Mark task as done",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        var id int
        fmt.Sscanf(args[0], "%d", &id)
        for i := range tasks {
            if tasks[i].ID == id {
                tasks[i].Done = true
                saveTasks()
                fmt.Printf("Marked task #%d as done\n", id)
                return
            }
        }
        fmt.Println("Task not found")
    },
}

var rootCmd = &cobra.Command{
    Use:   "tasks",
    Short: "Task manager",
}

func main() {
    rootCmd.AddCommand(addCmd, listCmd, doneCmd)
    rootCmd.Execute()
}
```

**Uso:**
```bash
$ go run main.go add
Task name: Buy groceries
Added task #1: Buy groceries

$ go run main.go list
[ ] #1: Buy groceries

$ go run main.go done 1
Marked task #1 as done

$ go run main.go list
[✓] #1: Buy groceries
```

---

## 51.12 Ejercicios Progresivos

### EJERCICIO 1: CLI Básica con Flags (Conversor de Temperatura)

**Objetivo:** Aprender parsing básico de flags con flag package estándar

**Especificación:**
- Convertir entre °C, °F, K
- Flags: `--from`, `--to`, `--value`
- Soportar múltiples conversiones en un comando

```bash
# Uso esperado:
$ go run converter.go --from celsius --to fahrenheit --value 25
25°C = 77°F

$ go run converter.go --from celsius --to kelvin --value 0
0°C = 273.15K
```

**Solución:**
```go
package main

import (
    "flag"
    "fmt"
    "log"
    "math"
    "os"
    "strings"
)

func convertTemp(value float64, from, to string) (float64, error) {
    // Normalizar a Celsius primero
    var celsius float64

    switch strings.ToLower(from) {
    case "celsius", "c":
        celsius = value
    case "fahrenheit", "f":
        celsius = (value - 32) * 5 / 9
    case "kelvin", "k":
        celsius = value - 273.15
    default:
        return 0, fmt.Errorf("unknown temperature unit: %s", from)
    }

    // Convertir desde Celsius
    switch strings.ToLower(to) {
    case "celsius", "c":
        return celsius, nil
    case "fahrenheit", "f":
        return celsius*9/5 + 32, nil
    case "kelvin", "k":
        return celsius + 273.15, nil
    default:
        return 0, fmt.Errorf("unknown temperature unit: %s", to)
    }
}

func formatUnit(unit string) string {
    switch strings.ToLower(unit) {
    case "celsius", "c":
        return "°C"
    case "fahrenheit", "f":
        return "°F"
    case "kelvin", "k":
        return "K"
    default:
        return unit
    }
}

func main() {
    from := flag.String("from", "celsius", "source temperature unit")
    to := flag.String("to", "fahrenheit", "target temperature unit")
    value := flag.Float64("value", 0, "temperature value")
    precision := flag.Int("precision", 2, "decimal places")

    flag.Parse()

    if *value == 0 && *from == *to {
        fmt.Println("Usage: converter --from <unit> --to <unit> --value <number>")
        fmt.Println("Units: celsius|fahrenheit|kelvin")
        os.Exit(0)
    }

    result, err := convertTemp(*value, *from, *to)
    if err != nil {
        log.Fatal(err)
    }

    // Redondear
    multiplier := math.Pow(10, float64(*precision))
    result = math.Round(result*multiplier) / multiplier

    fmt.Printf("%g%s = %g%s\n",
        *value, formatUnit(*from),
        result, formatUnit(*to))
}
```

### EJERCICIO 2: Subcomandos con Cobra (Task Manager)

**Objetivo:** Implementar CLI con subcomandos jerárquicos

**Especificación:**
- Comandos: `add`, `list`, `complete`, `delete`
- Flags: `--priority` (high|medium|low), `--due`, `--filter`
- Persistir en JSON

**Estructura esperada:**
```
myapp/
├── cmd/
│   ├── root.go
│   ├── add.go
│   ├── list.go
│   ├── complete.go
│   └── delete.go
├── task/
│   └── task.go    (lógica de tareas)
└── main.go
```

**Solución:** (Ver ejercicio anterior simplificado en 51.11.5)

### EJERCICIO 3: Input Interactivo con Promptui

**Objetivo:** Crear form interactiva con validación

**Especificación:**
- Preguntar: nombre, email, edad, país (select)
- Validar email con regex
- Validar edad (18-120)
- Salvar resultado en JSON

```bash
$ go run form.go
Name: Alice
Email: alice@example.com
Age: 25
Country:
  > Argentina
    Brazil
    Chile
    ...
Data saved to form_response.json
```

**Solución:**
```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "regexp"
    "strconv"
    "time"

    "github.com/manifoldco/promptui"
)

type FormData struct {
    Name    string    `json:"name"`
    Email   string    `json:"email"`
    Age     int       `json:"age"`
    Country string    `json:"country"`
    Timestamp time.Time `json:"timestamp"`
}

func validateEmail(email string) error {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    if matched, _ := regexp.MatchString(pattern, email); !matched {
        return fmt.Errorf("invalid email format")
    }
    return nil
}

func validateAge(ageStr string) error {
    age, err := strconv.Atoi(ageStr)
    if err != nil {
        return fmt.Errorf("age must be a number")
    }
    if age < 18 || age > 120 {
        return fmt.Errorf("age must be between 18 and 120")
    }
    return nil
}

func main() {
    countries := []string{
        "Argentina", "Brazil", "Chile", "Colombia", "Mexico",
        "Peru", "Venezuela", "Uruguay", "Paraguay", "Bolivia",
    }

    // Nombre
    namePrompt := promptui.Prompt{
        Label: "Name",
        Validate: func(s string) error {
            if len(s) < 2 {
                return fmt.Errorf("name must be at least 2 characters")
            }
            return nil
        },
    }
    name, _ := namePrompt.Run()

    // Email
    emailPrompt := promptui.Prompt{
        Label:    "Email",
        Validate: validateEmail,
    }
    email, _ := emailPrompt.Run()

    // Edad
    agePrompt := promptui.Prompt{
        Label:    "Age",
        Validate: validateAge,
    }
    ageStr, _ := agePrompt.Run()
    age, _ := strconv.Atoi(ageStr)

    // País
    countrySelect := promptui.Select{
        Label: "Country",
        Items: countries,
    }
    _, country, _ := countrySelect.Run()

    // Guardar datos
    data := FormData{
        Name:      name,
        Email:     email,
        Age:       age,
        Country:   country,
        Timestamp: time.Now(),
    }

    jsonData, _ := json.MarshalIndent(data, "", "  ")
    ioutil.WriteFile("form_response.json", jsonData, 0644)

    fmt.Println("\n✓ Data saved to form_response.json")
    fmt.Println(string(jsonData))
}
```

### EJERCICIO 4: Output Formateado con Tablas y Colores

**Objetivo:** Crear herramienta de monitoreo con salida formateada

**Especificación:**
- Listar procesos/servicios con estado
- Tabla con colores (verde=running, rojo=stopped, amarillo=warning)
- Filtrar por nombre o estado
- Soportar JSON output

```bash
$ go run monitor.go
NAME       STATUS     CPU%   MEM%   UPTIME
api-1      RUNNING    12.5   156MB  2d 5h 32m
api-2      RUNNING    8.3    124MB  1d 3h 45m
cache      STOPPED    0.0    0 MB   -

$ go run monitor.go --filter running --json
[{"name":"api-1","status":"RUNNING",...}]
```

**Solución:**
```go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "text/tabwriter"
    "time"

    "github.com/fatih/color"
)

type Service struct {
    Name   string  `json:"name"`
    Status string  `json:"status"`
    CPU    float64 `json:"cpu_percent"`
    Memory string  `json:"memory"`
    Uptime string  `json:"uptime"`
}

var mockServices = []Service{
    {"api-1", "RUNNING", 12.5, "156MB", "2d 5h 32m"},
    {"api-2", "RUNNING", 8.3, "124MB", "1d 3h 45m"},
    {"cache", "STOPPED", 0.0, "0MB", "-"},
    {"db", "RUNNING", 5.2, "512MB", "5d 12h 15m"},
    {"queue", "WARNING", 88.9, "1024MB", "3h 20m"},
}

func getStatusColor(status string) *color.Color {
    switch status {
    case "RUNNING":
        return color.New(color.FgGreen)
    case "STOPPED":
        return color.New(color.FgRed)
    case "WARNING":
        return color.New(color.FgYellow)
    default:
        return color.New(color.FgWhite)
    }
}

func main() {
    filter := flag.String("filter", "", "filter by status (running|stopped|warning)")
    jsonOutput := flag.Bool("json", false, "output as JSON")
    flag.Parse()

    // Filtrar servicios
    var filtered []Service
    for _, s := range mockServices {
        if *filter == "" || s.Status == string(rune(*filter[0]-32))+*filter[1:] {
            filtered = append(filtered, s)
        }
    }

    if *jsonOutput {
        data, _ := json.MarshalIndent(filtered, "", "  ")
        fmt.Println(string(data))
        return
    }

    // Mostrar tabla
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "NAME\tSTATUS\tCPU%\tMEMORY\tUPTIME")
    fmt.Fprintln(w, "----\t------\t----\t------\t------")

    for _, s := range filtered {
        statusColor := getStatusColor(s.Status)
        fmt.Fprintf(w, "%s\t%s\t%g\t%s\t%s\n",
            s.Name,
            statusColor.Sprint(s.Status),
            s.CPU,
            s.Memory,
            s.Uptime)
    }

    w.Flush()
}
```

### EJERCICIO 5: CLI Distribuible con Versioning (Avanzado)

**Objetivo:** Crear CLI lista para distribución con versioning, auto-updates, docs

**Especificación:**
- Comando `version` mostrando información detallada
- Archivo `VERSION` que se lee en build
- Self-update mechanism (check updates en GitHub)
- Generate man pages automático
- Cross-compile para Linux/macOS/Windows
- Release en GitHub con goreleaser

**Estructura:**
```
myapp/
├── .goreleaser.yaml
├── VERSION
├── cmd/
│   └── ...
├── internal/
│   ├── update/
│   │   └── updater.go    (check updates)
│   └── version/
│       └── version.go    (version info)
└── Makefile
```

**Makefile:**
```makefile
VERSION := $(shell cat VERSION)
BUILD := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X myapp/internal/version.Version=$(VERSION) \
           -X myapp/internal/version.Build=$(BUILD) \
           -X myapp/internal/version.Date=$(DATE)

.PHONY: build release test

build:
    CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o myapp .

release:
    goreleaser release --clean

test:
    go test ./...
```

---

## CONCLUSIÓN

Este capítulo ha cubierto todo lo necesario para crear CLIs profesionales en Go:

**Hemos aprendido:**
- ✅ Frameworks (Flag, Cobra, Urfave)
- ✅ I/O avanzado (prompts, progress, colors)
- ✅ Formatting (tablas, JSON, colores)
- ✅ Configuration hierarchy (Viper)
- ✅ Error handling y logging estructurado
- ✅ Testing integración
- ✅ Distribución y packaging
- ✅ UX principles y case studies

**Próximos pasos:**
- Explorar bubbletea para TUIs interactivas
- Estudiar sistemas de plugins
- Investigar CLI en contexto de DevOps (kubernetes, terraform)

---

## REFERENCIAS Y RECURSOS

**Documentación Oficial:**
- [Flag Package](https://golang.org/pkg/flag/)
- [Cobra Documentation](https://cobra.dev/)
- [Urfave/CLI](https://cli.urfave.org/)

**Librerías Recomendadas:**
- [Viper](https://github.com/spf13/viper) - Configuration
- [Promptui](https://github.com/manifoldco/promptui) - Interactive prompts
- [Tablewriter](https://github.com/olekukonko/tablewriter) - Tables
- [Logrus](https://github.com/sirupsen/logrus) - Logging
- [GoReleaser](https://goreleaser.com/) - Automation

**Ejemplos Reales:**
- kubectl: github.com/kubernetes/kubectl
- Hugo: github.com/gohugoio/hugo
- Docker CLI: github.com/docker/cli
- Terraform: github.com/hashicorp/terraform

---

**Versión:** 1.0  
**Última actualización:** Enero 2024  
**Autor:** Guía exhaustiva de Go  
**Nivel:** Intermedio → Avanzado

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/51-cli-tools/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/51-cli-tools):

```bash
cd examples/51-cli-tools
go run .
```
