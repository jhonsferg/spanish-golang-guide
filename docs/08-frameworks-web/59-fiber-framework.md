# Capítulo 59: Fiber - Velocidad e inspiración en Express.js

## Una Guía Exhaustiva del Framework Web Más Rápido de Go

**Tabla de Contenidos**
- [59.1 - Introducción a Fiber](#591---introducción-a-fiber)
- [59.2 - Fiber Basics](#592---fiber-basics)
- [59.3 - Routing Deep Dive](#593---routing-deep-dive)
- [59.4 - Middleware](#594---middleware)
- [59.5 - Request Handling](#595---request-handling)
- [59.6 - Response Methods](#596---response-methods)
- [59.7 - Validation & Error Handling](#597---validation--error-handling)
- [59.8 - Performance & Optimization](#598---performance--optimization)
- [59.9 - Advanced Features](#599---advanced-features)
- [59.10 - Testing & Monitoring](#5910---testing--monitoring)
- [59.11 - Production & Case Studies](#5911---production--case-studies)

---

## 59.1 - Introducción a Fiber

### 59.1.1 ¿Qué es Fiber?

Fiber es un **framework web Go extremadamente rápido** inspirado en Express.js, diseñado para aprovechar al máximo el rendimiento de Go. Construido sobre **fasthttp** (no net/http), Fiber proporciona una API familiar para desarrolladores de Node.js mientras mantiene la velocidad y eficiencia de Go.

### 59.1.2 Historia y Evolución

```
Timeline de Fiber:
├── 2018: Primera versión (inspiración en Express.js)
├── 2019: Adopción de fasthttp
├── 2020: Introducción de Middleware ecosystem
├── 2021: Versión 2.0 con mejoras de performance
├── 2022: WebSocket soporte nativo
├── 2023: v2.50+ - 70k+ estrellas en GitHub
└── 2024: Maduración y adoption en producción
```

### 59.1.3 Por qué Fiber: Ventajas Clave

**Rendimiento Superior**
```
Benchmarks (requests/segundo):
┌─────────────────┬─────────────┬──────────┐
│ Framework       │ req/s       │ Latencia │
├─────────────────┼─────────────┼──────────┤
│ Fiber           │ 330,000     │ 0.14ms   │
│ Gin             │ 200,000     │ 0.23ms   │
│ Echo            │ 185,000     │ 0.27ms   │
│ Express.js      │ 95,000      │ 0.52ms   │
│ net/http        │ 110,000     │ 0.45ms   │
└─────────────────┴─────────────┴──────────┘

Consumo de memoria (10k conexiones):
Fiber:        ~15 MB
Gin:          ~22 MB
Echo:         ~25 MB
Express.js:   ~85 MB
```

**API Familiar para Desarrolladores Node.js**
```go
// Fiber (inspira en Express)
app.Get("/users", handler)
app.Post("/users", handler)

// vs net/http (más verbose)
http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" { /*...*/ }
})
```

**Características Modernas**
- Built-in middleware ecosystem
- Validación de esquemas
- Manejo de errores global
- WebSocket nativo
- Server-Sent Events (SSE)
- Rate limiting
- CORS
- CSRF protection

### 59.1.4 Filosofía de Fiber

```
Principios Fundamentales:

1. VELOCIDAD PRIMERO
   └─ Optimizado para throughput máximo
   
2. FAMILIARIDAD
   └─ API similar a Express.js para adopción rápida
   
3. SIMPLICIDAD
   └─ API mínima pero poderosa
   
4. FLEXIBILIDAD
   └─ Middleware system para extensiones
   
5. GOLANG-NATIVE
   └─ Aprovecha concurrencia con goroutines
```

### 59.1.5 ¿Cuándo Usar Fiber?

**✅ Use Fiber cuando:**
- Necesita máximo rendimiento HTTP
- Construye REST APIs de alto tráfico
- Requiere baja latencia
- Viene de Express.js
- Necesita WebSocket de alto rendimiento
- Presupuesto de recursos limitado

**⚠️ Considere alternativas cuando:**
- Necesita middleware de gorilla/mux específico
- Proyecto usa net/http estándar extensivamente
- Requiere compatibilidad 100% con stdlib
- Criptomonedas blockchain (use Gin/Echo más estables)

### 59.1.6 Comparación: Fiber vs Gin vs Echo

```go
// FIBER - Express-like
app := fiber.New()
app.Get("/", handler).Name("home")

// GIN - Similar pero con grouting más simple
r := gin.Default()
r.GET("/", handler)

// ECHO - Estructurado pero más verbose
e := echo.New()
e.GET("/", handler)

// net/http - Lo más simple pero sin convenios
http.HandleFunc("/", handler)
```

**Tabla Comparativa Detallada:**

```
┌─────────────────┬──────────┬──────────┬──────────┬──────────┐
│ Característica  │ Fiber    │ Gin      │ Echo     │ net/http │
├─────────────────┼──────────┼──────────┼──────────┼──────────┤
│ Rendimiento     │ ★★★★★   │ ★★★★☆   │ ★★★★☆   │ ★★★☆☆   │
│ Curva Aprendiz  │ ★★★★★   │ ★★★★☆   │ ★★★★☆   │ ★★☆☆☆   │
│ Comunidad       │ ★★★★☆   │ ★★★★★   │ ★★★★☆   │ ★★★★★   │
│ Documentación   │ ★★★★☆   │ ★★★★★   │ ★★★★★   │ ★★★★★   │
│ Middleware eco  │ ★★★★☆   │ ★★★★★   │ ★★★★☆   │ ★★☆☆☆   │
│ WebSocket       │ ★★★★☆   │ ★★★☆☆   │ ★★★☆☆   │ ★★★☆☆   │
│ Validación      │ ★★★★☆   │ ★★★☆☆   │ ★★★★☆   │ ☆☆☆☆☆   │
│ Producción      │ ★★★★★   │ ★★★★★   │ ★★★★★   │ ★★★★★   │
└─────────────────┴──────────┴──────────┴──────────┴──────────┘
```

### 59.1.7 Casos de Uso Reales

**1. APIs REST de Alto Tráfico**
```
Sistema de recomendaciones: 
- 100k req/s
- Fiber maneja con <5 instancias
- Gin/Echo requieren 8-10 instancias
- Ahorro: $$$$ en infraestructura
```

**2. Microservicios**
```
Malla de servicios con ~20 microservicios
Fiber reduce:
- Latencia P99: 50ms → 15ms
- Consumo RAM: 85% reducción
- Deploy: 2 instancias vs 5
```

**3. WebSocket en Tiempo Real**
```
Chat de 50k usuarios concurrentes:
- Fiber: 1 servidor
- Express: 20+ servidores
- Razón: Fasthttp + goroutines
```

### 59.1.8 Comunidad y Ecosistema

```
Estadísticas (2024):
- 70,000+ estrellas GitHub
- 20+ middlewares oficiales
- 100+ middlewares comunitarios
- Usado por: Uber, Alibaba, CloudFlare
- NPM equivalente: npm.pkg.github.com/gofiber
```

**Recursos Clave:**
- GitHub: github.com/gofiber/fiber
- Documentación: docs.gofiber.io
- Discord: 15k+ miembros
- Ejemplos: github.com/gofiber/recipes

---

## 59.2 - Fiber Basics

### 59.2.1 Instalación

```bash
# Instalación estándar
go get -u github.com/gofiber/fiber/v2

# Verificar instalación
go get -u github.com/gofiber/fiber/v2@latest

# Con módulos en go.mod
require github.com/gofiber/fiber/v2 v2.50.0
```

**go.mod típico para proyecto Fiber:**
```
module github.com/user/myapp

go 1.21

require (
    github.com/gofiber/fiber/v2 v2.50.0
    github.com/gofiber/template/html/v2 v2.0.5
    github.com/google/uuid v1.5.0
)
```

### 59.2.2 Primera Aplicación

**Ejemplo Mínimo:**
```go
package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()
    
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("¡Hola, Fiber!")
    })
    
    app.Listen(":3000")
}
```

**Ejecutar:**
```bash
go run main.go
# Escuchar en http://localhost:3000

# Test con curl
curl http://localhost:3000
# ¡Hola, Fiber!
```

### 59.2.3 Creación de Aplicación Avanzada

```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    // Configuración detallada de la aplicación
    config := fiber.Config{
        AppName:      "Mi Aplicación v1.0",
        ServerHeader: "Fiber",
        StrictRoute:  true,           // /users != /users/ (default false)
        CaseSensitive: false,         // /Users == /users (default false)
        Immutable:    false,          // Permitir modificación de handlers
        BodyLimit:    1024 * 1024,    // 1MB
        Prefork:      false,          // Multi-proceso (true en producción)
        GETOnly:      false,          // Solo permitir GET
        ErrorHandler: customErrorHandler,
    }
    
    app := fiber.New(config)
    
    // Middleware global
    app.Use(logger.New())
    app.Use(recover.New())
    
    // Rutas
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Servidor Fiber ejecutándose",
            "version": "1.0.0",
        })
    })
    
    // Manejo de 404
    app.Use(func(c *fiber.Ctx) error {
        return c.Status(404).JSON(fiber.Map{
            "error": "Ruta no encontrada",
        })
    })
    
    log.Fatal(app.Listen(":3000"))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
    code := 500
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }
    return c.Status(code).JSON(fiber.Map{
        "error":   err.Error(),
        "status":  code,
        "path":    c.Path(),
    })
}
```

### 59.2.4 Contexto de Fiber (fiber.Ctx)

**¿Qué es el Contexto?**
```
El contexto (c *fiber.Ctx) representa una solicitud HTTP individual.
Incluye toda la información de la solicitud y métodos para responder.

Flujo del Contexto:
Request HTTP
    ↓
Crear fiber.Ctx
    ↓
Pasar a middleware
    ↓
Pasar a handler
    ↓
Handler escribe respuesta
    ↓
Response HTTP
```

**Propiedades Importantes de fiber.Ctx:**

```go
func handler(c *fiber.Ctx) error {
    // Información de la solicitud
    method := c.Method()           // GET, POST, etc.
    path := c.Path()               // /users/123
    query := c.Query("name")       // Parámetro query
    param := c.Params("id")        // Parámetro de ruta
    header := c.Get("Content-Type") // Header
    ip := c.IP()                   // IP del cliente
    body := c.Body()               // Body sin procesar
    
    // Información del contexto
    baseURL := c.BaseURL()         // http://localhost:3000
    hostname := c.Hostname()       // localhost
    port := c.Port()               // "3000"
    protocol := c.Protocol()       // http
    
    // Métodos de respuesta
    c.SendString("texto")          // Response de texto
    c.JSON(map[string]interface{}{}) // Response JSON
    c.Status(200)                  // Establecer status code
    c.Redirect("/other")           // Redireccionar
    
    // Valores locales (compartir datos en middleware)
    c.Locals("userID", 123)        // Establecer
    userID := c.Locals("userID")   // Obtener (interface{})
    
    // Cookies
    c.Cookie(&fiber.Cookie{
        Name:  "session",
        Value: "abc123",
    })
    
    return nil
}
```

### 59.2.5 Configuración de Listen

```go
package main

import (
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()
    
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hola")
    })
    
    // Diferentes formas de escuchar
    
    // 1. Puerto simple
    app.Listen(":3000")
    
    // 2. Host específico
    app.Listen("127.0.0.1:3000")
    
    // 3. Con configuración TLS
    app.ListenTLS(":443", "cert.pem", "key.pem")
    
    // 4. Con listener personalizado
    ln, _ := net.Listen("tcp", ":3000")
    app.Listener(ln)
    
    // 5. Graceful shutdown
    go func() {
        time.Sleep(10 * time.Second)
        app.Shutdown()
    }()
    app.Listen(":3000")
}
```

### 59.2.6 Estructura de Proyecto Recomendada

```
proyecto-fiber/
├── main.go                 # Punto de entrada
├── config/
│   └── config.go          # Configuración
├── handlers/
│   ├── user.go           # Handlers de usuarios
│   ├── post.go           # Handlers de posts
│   └── auth.go           # Handlers de autenticación
├── middleware/
│   ├── auth.go           # Middleware de autenticación
│   ├── logger.go         # Middleware de logging
│   └── cors.go           # Middleware CORS
├── models/
│   ├── user.go           # Estructura User
│   └── post.go           # Estructura Post
├── routes/
│   ├── user.go           # Rutas de usuarios
│   ├── post.go           # Rutas de posts
│   └── routes.go         # Registro de todas las rutas
├── services/
│   ├── user_service.go   # Lógica de negocios
│   └── post_service.go
├── database/
│   └── db.go             # Conexión a BD
├── utils/
│   ├── validator.go      # Validadores
│   └── jwt.go            # Utilidades JWT
├── templates/
│   ├── index.html        # Templates HTML
│   └── user.html
├── static/
│   ├── css/
│   ├── js/
│   └── images/
├── tests/
│   ├── handler_test.go
│   ├── middleware_test.go
│   └── integration_test.go
├── docker/
│   └── Dockerfile
├── .env                    # Variables de entorno
├── go.mod
├── go.sum
├── Makefile              # Comandos útiles
└── README.md
```

**Makefile típico:**
```makefile
.PHONY: run test build clean docker-build

run:
	go run main.go

test:
	go test ./... -v -coverage

build:
	go build -o app

clean:
	rm -f app

docker-build:
	docker build -t myapp .

lint:
	golangci-lint run
```

---

## 59.3 - Routing Deep Dive

### 59.3.1 Routing Fundamentales

**Métodos HTTP Básicos:**
```go
app.Get("/users", getUsers)        // GET
app.Post("/users", createUser)     // POST
app.Put("/users/:id", updateUser)  // PUT
app.Delete("/users/:id", deleteUser) // DELETE
app.Patch("/users/:id", patchUser) // PATCH
app.Options("/users", handleOptions) // OPTIONS
app.Head("/users", headUsers)      // HEAD
app.Trace("/users", traceUsers)    // TRACE

// Todas las rutas (cualquier método)
app.All("/health", healthCheck)

// Rutas personalizado con método
app.Add("CUSTOM", "/api/custom", handler)
```

### 59.3.2 Parámetros de Ruta

```go
// Parámetro simple
app.Get("/users/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{"id": id})
})
// GET /users/123 → id = "123"

// Parámetros múltiples
app.Get("/users/:id/posts/:postId", func(c *fiber.Ctx) error {
    id := c.Params("id")
    postId := c.Params("postId")
    return c.JSON(fiber.Map{
        "userId": id,
        "postId": postId,
    })
})
// GET /users/123/posts/456 → id="123", postId="456"

// Parámetro opcional con +
app.Get("/files/:name+", func(c *fiber.Ctx) error {
    name := c.Params("name")
    return c.JSON(fiber.Map{"file": name})
})
// GET /files/path/to/file.txt → name = "path/to/file.txt"

// Obtener todos los parámetros
app.Get("/users/:id", func(c *fiber.Ctx) error {
    allParams := c.AllParams()
    return c.JSON(allParams)
})
```

### 59.3.3 Parámetros de Query

```go
// Parámetro query único
app.Get("/search", func(c *fiber.Ctx) error {
    query := c.Query("q")           // Valor simple
    limit := c.Query("limit", "10") // Con default
    
    return c.JSON(fiber.Map{
        "search": query,
        "limit": limit,
    })
})
// GET /search?q=golang&limit=20

// Múltiples parámetros query (mismo nombre)
app.Get("/tags", func(c *fiber.Ctx) error {
    tags := c.GetReqHeaders()["Tag"] // []string
    return c.JSON(fiber.Map{"tags": tags})
})

// Todos los query parameters
app.Get("/filter", func(c *fiber.Ctx) error {
    queries := c.Queries() // map[string]string
    return c.JSON(queries)
})
// GET /filter?name=john&age=30&city=NYC
// Output: {"name": "john", "age": "30", "city": "NYC"}

// Query parsing avanzado
app.Get("/items", func(c *fiber.Ctx) error {
    // ?sort=name:asc,date:desc
    sort := c.Query("sort")
    
    // ?page=1&size=10
    page := c.QueryInt("page", 1)    // Valor entero
    size := c.QueryInt("size", 10)
    
    // ?active=true
    active := c.QueryBool("active", false)
    
    return c.JSON(fiber.Map{
        "sort": sort,
        "page": page,
        "size": size,
        "active": active,
    })
})
```

### 59.3.4 Grouping de Rutas

```go
// Grupo simple
api := app.Group("/api")
api.Get("/status", func(c *fiber.Ctx) error {
    return c.SendString("API OK")
})
// GET /api/status

// Grupos anidados
v1 := app.Group("/api/v1")
users := v1.Group("/users")

users.Get("/", getAllUsers)           // GET /api/v1/users
users.Post("/", createUser)           // POST /api/v1/users
users.Get("/:id", getUser)            // GET /api/v1/users/:id
users.Put("/:id", updateUser)         // PUT /api/v1/users/:id
users.Delete("/:id", deleteUser)      // DELETE /api/v1/users/:id

// Grupos con middleware
admin := app.Group("/admin", authMiddleware, adminMiddleware)
admin.Get("/dashboard", adminDashboard)
admin.Get("/users", listAllUsers)

// Combinación de grupos y middleware
api.Use(corsMiddleware)
api.Use(rateLimitMiddleware)

protected := api.Group("/protected", authMiddleware)
protected.Get("/profile", userProfile)
protected.Post("/settings", updateSettings)
```

### 59.3.5 Rutas Regex (Expresiones Regulares)

```go
// Ruta simple regex
app.Get(`/:name([a-z]+)`, func(c *fiber.Ctx) error {
    name := c.Params("name")
    return c.JSON(fiber.Map{"name": name})
})
// GET /john → OK
// GET /john123 → Not match

// Regex complejo
app.Get(`/files/:name(\.[a-z]{2,4})`, func(c *fiber.Ctx) error {
    name := c.Params("name")
    return c.JSON(fiber.Map{"file": name})
})
// GET /files/.txt → OK
// GET /files/.verylongextension → Not match

// ID numérico
app.Get(`/users/:id(\d+)`, func(c *fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{"id": id})
})
// GET /users/123 → OK
// GET /users/abc → Not match

// UUID
app.Get(`/items/:uuid([0-9a-fA-F-]{36})`, func(c *fiber.Ctx) error {
    uuid := c.Params("uuid")
    return c.JSON(fiber.Map{"uuid": uuid})
})

// Combinaciones complejas
app.Get(`/:lang([a-z]{2})/posts/:id(\d+)`, func(c *fiber.Ctx) error {
    lang := c.Params("lang")
    id := c.Params("id")
    return c.JSON(fiber.Map{
        "lang": lang,
        "post_id": id,
    })
})
// GET /es/posts/123 → OK
// GET /english/posts/123 → Not match (lang > 2 chars)
```

### 59.3.6 Rutas Nombradas

```go
// Definir rutas con nombres
app.Get("/", homeHandler).Name("home")
app.Get("/about", aboutHandler).Name("about")
app.Get("/users/:id", userDetail).Name("user.detail")
app.Post("/users", createUser).Name("user.create")

// Usar nombres en respuestas
app.Get("/links", func(c *fiber.Ctx) error {
    // Generar URLs usando nombres
    homeUrl := c.GetRouteURL("home", fiber.Map{})
    aboutUrl := c.GetRouteURL("about", fiber.Map{})
    userUrl := c.GetRouteURL("user.detail", fiber.Map{
        "id": "123",
    })
    
    return c.JSON(fiber.Map{
        "home": homeUrl,
        "about": aboutUrl,
        "user": userUrl,
    })
})
// Output:
// {
//   "home": "/",
//   "about": "/about",
//   "user": "/users/123"
// }
```

### 59.3.7 Archivos Estáticos

```go
// Servir carpeta estática
app.Static("/public", "./public")
// GET /public/style.css → ./public/style.css
// GET /public/js/app.js → ./public/js/app.js

// Con configuración
app.Static("/assets", "./static", fiber.Static{
    Compress:      true,           // Comprimir archivos
    ByteRange:     true,           // Soportar Range requests
    Browse:        true,           // Listar directorio
    Index:         "index.html",   // Archivo por defecto
    CacheDuration: 1 * time.Hour,  // Cache en cliente
})

// Ruta personalizada para archivos
app.Get("/download/:file", func(c *fiber.Ctx) error {
    file := c.Params("file")
    return c.Download("./downloads/" + file)
})

// Servir un archivo específico
app.Get("/favicon.ico", func(c *fiber.Ctx) error {
    return c.SendFile("./public/favicon.ico")
})

// Combinación: estático + API
app.Static("/", "./public")        // Servir HTML/CSS/JS
app.Get("/api/data", apiHandler)   // API en /api
```

### 59.3.8 Router Avanzado: Wildcard y Fallback

```go
// Wildcard routes (deben estar al final)
app.Get("/api/*", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "message": "Ruta catch-all",
        "path": c.Path(),
    })
})

// Fallback 404
app.Use(func(c *fiber.Ctx) error {
    return c.Status(404).JSON(fiber.Map{
        "error": "No encontrado",
        "path": c.Path(),
        "method": c.Method(),
    })
})

// Prioridad de rutas
app.Get("/users/me", getCurrentUser)      // Específica (ejecuta primero)
app.Get("/users/:id", getUserById)        // General (ejecuta después)

// Orden:
// 1. Rutas exactas específicas
// 2. Rutas con parámetros
// 3. Rutas regex
// 4. Wildcard
// 5. Fallback 404
```

---

## 59.4 - Middleware

### 59.4.1 Concepto de Middleware

```
Flujo de Middleware:

Request HTTP
    ↓
┌─────────────────┐
│ Middleware 1    │
│ (Logger)        │
└──────┬──────────┘
       ↓
┌─────────────────┐
│ Middleware 2    │
│ (Auth)          │
└──────┬──────────┘
       ↓
┌─────────────────┐
│ Middleware 3    │
│ (RateLimit)     │
└──────┬──────────┘
       ↓
┌─────────────────┐
│ Handler         │
│ (Lógica)        │
└──────┬──────────┘
       ↓
Response HTTP
```

### 59.4.2 Middleware Básico

```go
// Middleware simple
func customMiddleware(c *fiber.Ctx) error {
    // Antes del handler
    println("Solicitud entrante:", c.Path())
    
    // Continuar al siguiente middleware/handler
    err := c.Next()
    
    // Después del handler
    println("Solicitud completada:", c.Path())
    
    return err
}

// Usar middleware
app.Use(customMiddleware)           // Global
app.Get("/", handler)

// Middleware específico para una ruta
app.Get("/admin", adminMiddleware, adminHandler)

// Múltiples middleware en una ruta
app.Post("/secure", 
    authMiddleware,
    validateMiddleware,
    secureHandler,
)

// Middleware en grupo
api := app.Group("/api", corsMiddleware, rateLimitMiddleware)
api.Get("/status", statusHandler)
```

### 59.4.3 Middleware Incorporado de Fiber

**Logger Middleware:**
```go
import "github.com/gofiber/fiber/v2/middleware/logger"

app.Use(logger.New(logger.Config{
    Format:        "[${ip}] ${status} - ${method} ${path}\n",
    TimeFormat:    "02/Jan/2006 15:04:05",
    TimeZone:      "America/New_York",
    Output:        os.Stdout,
}))

// Formato:
// [127.0.0.1] 200 - GET /users
// [192.168.1.1] 404 - POST /api/data
```

**Recover Middleware (Manejo de Panic):**
```go
import "github.com/gofiber/fiber/v2/middleware/recover"

app.Use(recover.New())  // Atrapa panics y evita crash

// Sin recover:
app.Get("/crash", func(c *fiber.Ctx) error {
    var arr []int
    return c.SendString(arr[100])  // Panic! ☠️
})

// Con recover: retorna 500 en lugar de crash
```

**CORS Middleware:**
```go
import "github.com/gofiber/fiber/v2/middleware/cors"

app.Use(cors.New(cors.Config{
    AllowOrigins:     "https://example.com,https://example2.com",
    AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
    AllowHeaders:     "Content-Type,Authorization",
    ExposeHeaders:    "Content-Length",
    MaxAge:           3600,
    AllowCredentials: true,
}))

// CORS completo (permitir todas las origins)
app.Use(cors.New())
```

**Rate Limiting:**
```go
import "github.com/gofiber/fiber/v2/middleware/limiter"
import "github.com/gofiber/fiber/v2/utils"

app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()  // Por IP
    },
    LimitReached: func(c *fiber.Ctx) error {
        return c.Status(429).JSON(fiber.Map{
            "error": "Demasiadas solicitudes",
        })
    },
}))

// Resultado: máximo 100 requests por minuto por IP
```

**Compression Middleware:**
```go
import "github.com/gofiber/fiber/v2/middleware/compress"

app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,  // LevelDefault, LevelBestSpeed, LevelBestCompression
}))

// Automáticamente comprime respuestas > 4KB
// Soporta gzip, deflate, brotli
```

**Request ID:**
```go
import "github.com/gofiber/fiber/v2/middleware/requestid"

app.Use(requestid.New())

app.Get("/", func(c *fiber.Ctx) error {
    id := c.Get(requestid.ConfigDefault.Header)
    return c.JSON(fiber.Map{"request_id": id})
})
// Output: {"request_id": "9f7b8b9e-8d9b-4b9b-9b9b-9b9b9b9b9b9b"}
```

**Timeout Middleware:**
```go
import "github.com/gofiber/fiber/v2/middleware/timeout"

app.Use(timeout.New(timeout.Config{
    Timeout: 5 * time.Second,
}))

// Si el handler tarda > 5s, retorna 408 Request Timeout
```

### 59.4.4 Custom Middleware Avanzado

**Middleware de Autenticación:**
```go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "strings"
)

func AuthMiddleware(c *fiber.Ctx) error {
    // Obtener token del header
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{
            "error": "Token no proporcionado",
        })
    }
    
    // Validar formato Bearer
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return c.Status(401).JSON(fiber.Map{
            "error": "Formato de token inválido",
        })
    }
    
    token := parts[1]
    
    // Validar token (ejemplo simple)
    userID, err := validateToken(token)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Token inválido",
        })
    }
    
    // Almacenar en contexto
    c.Locals("userID", userID)
    c.Locals("token", token)
    
    return c.Next()
}

func validateToken(token string) (int, error) {
    // Implementar validación JWT o similar
    // Por ahora, ejemplo simplificado
    if token == "" {
        return 0, fiber.NewError(401, "Token vacío")
    }
    return 123, nil
}
```

**Middleware de Logging Personalizado:**
```go
package middleware

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "time"
)

func CustomLogger(c *fiber.Ctx) error {
    start := time.Now()
    
    // Pasar al siguiente handler
    err := c.Next()
    
    // Calcular duración
    duration := time.Since(start)
    
    // Log
    statusCode := c.Response().StatusCode()
    color := getStatusColor(statusCode)
    
    fmt.Printf("%s[%d]%s %s %s en %v\n",
        color,
        statusCode,
        "\033[0m",  // Reset color
        c.Method(),
        c.Path(),
        duration,
    )
    
    return err
}

func getStatusColor(status int) string {
    switch {
    case status >= 200 && status < 300:
        return "\033[32m"  // Green
    case status >= 300 && status < 400:
        return "\033[36m"  // Cyan
    case status >= 400 && status < 500:
        return "\033[33m"  // Yellow
    default:
        return "\033[31m"  // Red
    }
}
```

**Middleware de Validación Global:**
```go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateContentType(c *fiber.Ctx) error {
    if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
        contentType := c.Get("Content-Type")
        if contentType != "application/json" {
            return c.Status(400).JSON(fiber.Map{
                "error": "Content-Type debe ser application/json",
            })
        }
    }
    return c.Next()
}

// Usar en rutas específicas
app.Post("/users", ValidateContentType, createUserHandler)
```

### 59.4.5 Orden de Ejecución de Middleware

```go
func main() {
    app := fiber.New()
    
    // 1. Middleware Global
    app.Use(func(c *fiber.Ctx) error {
        println("1. Middleware Global - ANTES")
        err := c.Next()
        println("1. Middleware Global - DESPUÉS")
        return err
    })
    
    // 2. Grupo con middleware
    api := app.Group("/api", func(c *fiber.Ctx) error {
        println("2. Grupo Middleware - ANTES")
        err := c.Next()
        println("2. Grupo Middleware - DESPUÉS")
        return err
    })
    
    // 3. Ruta específica con middleware
    api.Get("/users", func(c *fiber.Ctx) error {
        println("3. Middleware Ruta - ANTES")
        err := c.Next()
        println("3. Middleware Ruta - DESPUÉS")
        return err
    }, func(c *fiber.Ctx) error {
        println("4. HANDLER")
        return c.SendString("OK")
    })
    
    app.Listen(":3000")
}

// Orden de ejecución para GET /api/users:
// 1. Middleware Global - ANTES
// 2. Grupo Middleware - ANTES
// 3. Middleware Ruta - ANTES
// 4. HANDLER
// 3. Middleware Ruta - DESPUÉS
// 2. Grupo Middleware - DESPUÉS
// 1. Middleware Global - DESPUÉS
```

---

## 59.5 - Request Handling

### 59.5.1 Parsing JSON

```go
// Modelo de datos
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// Parsing JSON simple
app.Post("/users", func(c *fiber.Ctx) error {
    user := new(User)
    
    // Parsear JSON del body
    if err := c.BindJSON(user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Usuario recibido",
        "user": user,
    })
})

// Parsing JSON con manejo de errores personalizado
app.Post("/users", func(c *fiber.Ctx) error {
    user := new(User)
    
    if err := c.BindJSON(user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "JSON inválido",
            "details": err.Error(),
        })
    }
    
    // Validar datos
    if user.Name == "" || user.Email == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Nombre y email requeridos",
        })
    }
    
    return c.JSON(user)
})

// Parsing JSON sin struct (usando map)
app.Post("/flexible", func(c *fiber.Ctx) error {
    var data map[string]interface{}
    
    if err := c.BindJSON(&data); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    // Acceder a valores
    name := data["name"].(string)
    age := int(data["age"].(float64))
    
    return c.JSON(fiber.Map{
        "name": name,
        "age": age,
    })
})
```

### 59.5.2 Form Data

```go
// Modelo
type LoginForm struct {
    Username string `form:"username"`
    Password string `form:"password"`
    Remember bool   `form:"remember"`
}

// Parsing form-urlencoded
app.Post("/login", func(c *fiber.Ctx) error {
    form := new(LoginForm)
    
    if err := c.BindForm(form); err != nil {
        return c.Status(400).SendString("Error parsing form")
    }
    
    return c.JSON(fiber.Map{
        "username": form.Username,
        "remember": form.Remember,
    })
})

// HTML form correspondiente:
// <form method="POST" action="/login">
//   <input type="text" name="username">
//   <input type="password" name="password">
//   <input type="checkbox" name="remember">
//   <button type="submit">Login</button>
// </form>

// Obtener valores específicos
app.Post("/form", func(c *fiber.Ctx) error {
    username := c.FormValue("username")
    password := c.FormValue("password", "") // Con default
    
    return c.JSON(fiber.Map{
        "username": username,
        "password": password,
    })
})
```

### 59.5.3 URL Encoding

```go
type SearchQuery struct {
    Q      string `query:"q"`
    Limit  int    `query:"limit"`
    Offset int    `query:"offset"`
    Tags   []string `query:"tags"`
}

// Parsing query parameters
app.Get("/search", func(c *fiber.Ctx) error {
    query := new(SearchQuery)
    
    if err := c.BindQuery(query); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(query)
})

// GET /search?q=golang&limit=10&offset=0&tags=web&tags=api
// Output: {
//   "q": "golang",
//   "limit": 10,
//   "offset": 0,
//   "tags": ["web", "api"]
// }
```

### 59.5.4 File Uploads

**Archivo Único:**
```go
app.Post("/upload", func(c *fiber.Ctx) error {
    // Obtener archivo
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "No file provided",
        })
    }
    
    // Información del archivo
    filename := file.Filename
    size := file.Size
    
    // Validar tamaño (máx 10MB)
    if size > 10*1024*1024 {
        return c.Status(400).JSON(fiber.Map{
            "error": "File too large",
        })
    }
    
    // Guardar archivo
    path := fmt.Sprintf("./uploads/%s", filename)
    if err := c.SaveFile(file, path); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to save file",
        })
    }
    
    return c.JSON(fiber.Map{
        "filename": filename,
        "size": size,
        "path": path,
    })
})
```

**Múltiples Archivos:**
```go
app.Post("/upload-multiple", func(c *fiber.Ctx) error {
    // Obtener formulario
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).SendString("Error parsing form")
    }
    
    // Obtener archivos
    files := form.File["files"]
    
    uploaded := []string{}
    for _, file := range files {
        filename := file.Filename
        path := fmt.Sprintf("./uploads/%s", filename)
        
        if err := c.SaveFile(file, path); err != nil {
            continue
        }
        uploaded = append(uploaded, filename)
    }
    
    return c.JSON(fiber.Map{
        "uploaded": uploaded,
        "count": len(uploaded),
    })
})

// HTML correspondiente:
// <form method="POST" action="/upload-multiple" enctype="multipart/form-data">
//   <input type="file" name="files" multiple>
//   <button type="submit">Upload</button>
// </form>
```

### 59.5.5 Body Raw

```go
// Acceder a body sin procesar
app.Post("/raw", func(c *fiber.Ctx) error {
    body := c.Body()
    
    println("Raw body:", string(body))
    
    return c.JSON(fiber.Map{
        "bytes_received": len(body),
    })
})

// Body como string
app.Post("/text", func(c *fiber.Ctx) error {
    text := c.BodyParser(func(body []byte) error {
        println("Body:", string(body))
        return nil
    })
    
    return c.SendString("OK")
})
```

### 59.5.6 Cookies

**Lectura de Cookies:**
```go
app.Get("/", func(c *fiber.Ctx) error {
    // Obtener cookie específica
    sessionID := c.Cookies("sessionid")
    userPref := c.Cookies("theme", "dark") // Con default
    
    // Todas las cookies
    allCookies := c.Cookies()
    
    return c.JSON(fiber.Map{
        "sessionid": sessionID,
        "theme": userPref,
        "all": allCookies,
    })
})
```

**Escritura de Cookies:**
```go
app.Post("/login", func(c *fiber.Ctx) error {
    // Crear cookie
    c.Cookie(&fiber.Cookie{
        Name:     "sessionid",
        Value:    "abc123",
        Path:     "/",
        Domain:   "example.com",
        Expires:  time.Now().Add(24 * time.Hour),
        Secure:   true,        // Solo HTTPS
        HTTPOnly: true,        // No accesible desde JS
        SameSite: "Strict",    // CSRF protection
    })
    
    return c.JSON(fiber.Map{
        "message": "Cookie set",
    })
})

// Eliminar cookie
app.Post("/logout", func(c *fiber.Ctx) error {
    c.Cookie(&fiber.Cookie{
        Name:    "sessionid",
        Value:   "",
        Expires: time.Now().Add(-time.Hour),
    })
    
    return c.SendString("Logged out")
})
```

### 59.5.7 Headers

```go
// Lectura de headers
app.Get("/", func(c *fiber.Ctx) error {
    contentType := c.Get("Content-Type")
    userAgent := c.Get("User-Agent")
    auth := c.Get("Authorization")
    
    return c.JSON(fiber.Map{
        "content_type": contentType,
        "user_agent": userAgent,
        "auth": auth,
    })
})

// Todos los headers
app.Get("/headers", func(c *fiber.Ctx) error {
    headers := c.GetReqHeaders()
    return c.JSON(headers)
})

// Escritura de headers
app.Get("/download", func(c *fiber.Ctx) error {
    c.Set("Content-Type", "application/octet-stream")
    c.Set("Content-Disposition", `attachment; filename="file.pdf"`)
    c.Set("X-Custom-Header", "custom-value")
    
    return c.SendFile("./file.pdf")
})
```

---

## 59.6 - Response Methods

### 59.6.1 JSON Response

```go
app.Get("/api/user", func(c *fiber.Ctx) error {
    user := fiber.Map{
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
    }
    
    return c.JSON(user)
})

// Con indentación (pretty print)
app.Get("/pretty", func(c *fiber.Ctx) error {
    user := fiber.Map{"id": 1, "name": "John"}
    return c.JSON(user)  // Automáticamente indentado en desarrollo
})

// JSON con status code personalizado
app.Post("/users", func(c *fiber.Ctx) error {
    newUser := fiber.Map{
        "id": 123,
        "name": "Jane Doe",
    }
    
    return c.Status(201).JSON(newUser)  // 201 Created
})

// Array JSON
app.Get("/users", func(c *fiber.Ctx) error {
    users := []fiber.Map{
        {"id": 1, "name": "User 1"},
        {"id": 2, "name": "User 2"},
    }
    
    return c.JSON(users)
})
```

### 59.6.2 SendString (Texto Plano)

```go
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("¡Hola Mundo!")
})

// Con status code
app.Get("/error", func(c *fiber.Ctx) error {
    return c.Status(500).SendString("Error interno del servidor")
})

// Caracteres especiales
app.Get("/special", func(c *fiber.Ctx) error {
    text := "Héllo! ñ, 中文, 🔥"
    return c.SendString(text)
})
```

### 59.6.3 HTML Templates

**Instalación:**
```bash
go get -u github.com/gofiber/template/html/v2
```

**Configuración:**
```go
import "github.com/gofiber/template/html/v2"

func main() {
    // Crear engine HTML
    engine := html.New("./views", ".html")
    
    app := fiber.New(fiber.Config{
        Views: engine,
    })
    
    app.Get("/", func(c *fiber.Ctx) error {
        // Renderizar template
        return c.Render("index", fiber.Map{
            "Title": "Inicio",
            "Message": "Bienvenido a Fiber",
        })
    })
    
    app.Listen(":3000")
}
```

**Template HTML (views/index.html):**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <h1>{{ .Message }}</h1>
    <p>Fecha: {{ .Date }}</p>
</body>
</html>
```

**Más ejemplos:**
```go
// Template con loop
app.Get("/users", func(c *fiber.Ctx) error {
    users := []fiber.Map{
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
    }
    
    return c.Render("users", fiber.Map{
        "Users": users,
    })
})

// Template users.html:
// <ul>
//   {{range .Users}}
//     <li>{{.name}}</li>
//   {{end}}
// </ul>
```

### 59.6.4 File/Download

```go
// Descargar archivo
app.Get("/download/:file", func(c *fiber.Ctx) error {
    file := c.Params("file")
    return c.Download(fmt.Sprintf("./files/%s", file))
})

// Enviar archivo (inline)
app.Get("/file/:name", func(c *fiber.Ctx) error {
    name := c.Params("name")
    return c.SendFile(fmt.Sprintf("./files/%s", name))
})

// Con nombre personalizado
app.Get("/export", func(c *fiber.Ctx) error {
    c.Set("Content-Disposition", "attachment; filename=export.csv")
    return c.SendFile("./temp/data.csv")
})
```

### 59.6.5 Streaming

```go
import "io"

// Streaming de texto
app.Get("/stream", func(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("Connection", "keep-alive")
    
    for i := 0; i < 5; i++ {
        c.WriteString(fmt.Sprintf("data: Mensaje %d\n\n", i))
        time.Sleep(1 * time.Second)
    }
    
    return nil
})

// Streaming de archivo grande
app.Get("/stream-file", func(c *fiber.Ctx) error {
    file, _ := os.Open("./large-file.bin")
    defer file.Close()
    
    c.Set("Content-Type", "application/octet-stream")
    return c.SendStream(file)
})
```

### 59.6.6 Redirect

```go
// Redirección simple (302)
app.Get("/old-path", func(c *fiber.Ctx) error {
    return c.Redirect("/new-path")
})

// Redirección con status code
app.Get("/permanent", func(c *fiber.Ctx) error {
    return c.Redirect("/new", 301)  // Moved Permanently
})

// Redirección condicional
app.Get("/check", func(c *fiber.Ctx) error {
    if !isAuthorized(c) {
        return c.Redirect("/login")
    }
    return c.SendString("Bienvenido")
})

// Redirección a URL absoluta
app.Get("/external", func(c *fiber.Ctx) error {
    return c.Redirect("https://example.com")
})
```

---

## 59.7 - Validation & Error Handling

### 59.7.1 Schema Validation

**Instalación de Validator:**
```bash
go get -u github.com/go-playground/validator/v10
```

**Validación Básica:**
```go
import "github.com/go-playground/validator/v10"

type User struct {
    Name  string `json:"name" validate:"required,min=3,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,gte=18,lte=120"`
}

var validate = validator.New()

app.Post("/users", func(c *fiber.Ctx) error {
    user := new(User)
    
    if err := c.BindJSON(user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "JSON inválido",
        })
    }
    
    // Validar struct
    if err := validate.Struct(user); err != nil {
        errors := err.(validator.ValidationErrors)
        return c.Status(400).JSON(fiber.Map{
            "errors": formatValidationErrors(errors),
        })
    }
    
    return c.JSON(user)
})

func formatValidationErrors(errs validator.ValidationErrors) []fiber.Map {
    var result []fiber.Map
    for _, err := range errs {
        result = append(result, fiber.Map{
            "field": err.Field(),
            "tag": err.Tag(),
            "message": getErrorMessage(err),
        })
    }
    return result
}

func getErrorMessage(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return fmt.Sprintf("%s es requerido", err.Field())
    case "email":
        return fmt.Sprintf("%s debe ser un email válido", err.Field())
    case "min":
        return fmt.Sprintf("%s debe tener mínimo %s caracteres", err.Field(), err.Param())
    case "max":
        return fmt.Sprintf("%s debe tener máximo %s caracteres", err.Field(), err.Param())
    default:
        return fmt.Sprintf("%s es inválido", err.Field())
    }
}
```

**Tags de Validación Comunes:**
```go
type Product struct {
    // Strings
    Name     string `validate:"required"`                 // Requerido
    Slug     string `validate:"required,alphanum"`        // Alfanumérico
    URL      string `validate:"required,url"`             // URL válida
    Email    string `validate:"required,email"`           // Email válido
    Phone    string `validate:"required,e164"`            // Formato E.164
    Desc     string `validate:"max=500"`                  // Longitud máxima
    
    // Números
    Age      int    `validate:"gte=0,lte=150"`            // Entre 0 y 150
    Price    float64 `validate:"gt=0"`                    // Mayor que 0
    Rating   float32 `validate:"gte=0,lte=5"`             // Entre 0 y 5
    
    // Booleano
    Active   bool   `validate:"required"`                 // Requerido
    
    // Structs anidados
    Address  Address `validate:"required"`                // Validar anidado
    
    // Arrays
    Tags     []string `validate:"required,min=1,dive,required,min=3"`  // Array con validación
}
```

### 59.7.2 Custom Validators

```go
// Definir validador personalizado
func init() {
    validate.RegisterValidation("username", validateUsername)
    validate.RegisterValidation("zipcode", validateZipcode)
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // No permitir nombres reservados
    reserved := []string{"admin", "root", "system"}
    for _, r := range reserved {
        if username == r {
            return false
        }
    }
    return len(username) >= 3
}

func validateZipcode(fl validator.FieldLevel) bool {
    zipcode := fl.Field().String()
    // Validar formato CP España
    matched, _ := regexp.MatchString(`^[0-9]{5}$`, zipcode)
    return matched
}

// Usar en struct
type CreateAccount struct {
    Username string `validate:"required,username"`
    Zipcode  string `validate:"required,zipcode"`
}
```

### 59.7.3 Global Error Handler

```go
func main() {
    config := fiber.Config{
        ErrorHandler: globalErrorHandler,
    }
    
    app := fiber.New(config)
    app.Listen(":3000")
}

func globalErrorHandler(c *fiber.Ctx, err error) error {
    code := 500
    message := "Error interno"
    
    // Fiber errors
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
        message = e.Message
    }
    
    // Custom errors
    if customErr, ok := err.(*CustomError); ok {
        code = customErr.Code
        message = customErr.Message
    }
    
    return c.Status(code).JSON(fiber.Map{
        "error": message,
        "code": code,
        "timestamp": time.Now(),
    })
}

type CustomError struct {
    Code    int
    Message string
}

func (e *CustomError) Error() string {
    return e.Message
}
```

### 59.7.4 HTTP Status Codes

```go
// 2xx - Éxito
app.Get("/success", func(c *fiber.Ctx) error {
    return c.Status(200).SendString("OK")
})

app.Post("/created", func(c *fiber.Ctx) error {
    return c.Status(201).SendString("Creado")
})

// 3xx - Redirección
app.Get("/redirect", func(c *fiber.Ctx) error {
    return c.Status(302).Redirect("/new-location")
})

// 4xx - Error del cliente
app.Get("/error400", func(c *fiber.Ctx) error {
    return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
})

app.Get("/error401", func(c *fiber.Ctx) error {
    return c.Status(401).JSON(fiber.Map{"error": "No autorizado"})
})

app.Get("/error403", func(c *fiber.Ctx) error {
    return c.Status(403).JSON(fiber.Map{"error": "Prohibido"})
})

app.Get("/error404", func(c *fiber.Ctx) error {
    return c.Status(404).JSON(fiber.Map{"error": "No encontrado"})
})

app.Get("/error422", func(c *fiber.Ctx) error {
    return c.Status(422).JSON(fiber.Map{"error": "Datos inválidos"})
})

app.Get("/error429", func(c *fiber.Ctx) error {
    return c.Status(429).JSON(fiber.Map{"error": "Demasiadas solicitudes"})
})

// 5xx - Error del servidor
app.Get("/error500", func(c *fiber.Ctx) error {
    return c.Status(500).JSON(fiber.Map{"error": "Error interno"})
})

app.Get("/error503", func(c *fiber.Ctx) error {
    return c.Status(503).JSON(fiber.Map{"error": "Servicio no disponible"})
})
```

---

## 59.8 - Performance & Optimization

### 59.8.1 Benchmarks de Rendimiento

```
Comparativa de throughput (requests/segundo):

┌────────────────┬──────────────┬──────────┬─────────────┐
│ Framework      │ Throughput   │ Latencia │ Memoria     │
├────────────────┼──────────────┼──────────┼─────────────┤
│ Fiber          │ 330,000 req/s│ 0.14ms   │ 15 MB       │
│ Gin (std)      │ 200,000 req/s│ 0.23ms   │ 22 MB       │
│ Echo           │ 185,000 req/s│ 0.27ms   │ 25 MB       │
│ net/http       │ 110,000 req/s│ 0.45ms   │ 20 MB       │
│ Express (Node) │ 95,000 req/s │ 0.52ms   │ 85 MB       │
│ Flask (Python) │ 45,000 req/s │ 1.1ms    │ 45 MB       │
└────────────────┴──────────────┴──────────┴─────────────┘

Conexiones concurrentes (10,000):
Fiber:    1-2 instancias
Gin:      2-3 instancias
Echo:     2-3 instancias
Express:  10+ instancias
```

### 59.8.2 Optimizaciones de Fiber

**Configuración de Producción:**
```go
app := fiber.New(fiber.Config{
    // Network
    Network:           "tcp4",        // TCP4 solo (más rápido)
    Prefork:           true,          // Usar múltiples procesos
    MaxRequestsPerConn: 0,            // Sin límite
    
    // Performance
    BodyLimit:        1024 * 1024,    // 1MB
    Concurrency:      256 * 1024,     // Max goroutines simultáneas
    DisableKeepalive: false,          // Mantener conexiones vivas
    
    // Compresión
    CompressedFileSuffix: ".fiber.gz",
    
    // Errores
    DisableStartupMessage: false,
})

// Pooling de buffers (reduce GC)
app.Use(func(c *fiber.Ctx) error {
    // Fiber automáticamente maneja pooling
    return c.Next()
})
```

### 59.8.3 Memory Optimization

```go
// ❌ MAL - Crear slice grande cada request
app.Get("/data", func(c *fiber.Ctx) error {
    data := make([]byte, 1000000)  // 1MB cada request!
    return c.Send(data)
})

// ✅ BIEN - Reutilizar buffer
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1000000)
    },
}

app.Get("/data", func(c *fiber.Ctx) error {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    
    return c.Send(buf)
})

// ❌ MAL - Múltiples conversiones string
app.Get("/process", func(c *fiber.Ctx) error {
    for i := 0; i < 1000; i++ {
        str := fmt.Sprintf("Item %d", i)
        println(str)
    }
    return nil
})

// ✅ BIEN - Usar strings.Builder
app.Get("/process", func(c *fiber.Ctx) error {
    var sb strings.Builder
    for i := 0; i < 1000; i++ {
        fmt.Fprintf(&sb, "Item %d\n", i)
    }
    return c.SendString(sb.String())
})
```

### 59.8.4 Connection Pooling (Database)

```go
import "database/sql"

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", "connection_string")
    if err != nil {
        panic(err)
    }
    
    // Configuración de pool
    db.SetMaxOpenConns(25)    // Máximo conexiones abiertas
    db.SetMaxIdleConns(5)     // Máximo conexiones ociosas
    db.SetConnMaxLifetime(5 * time.Minute)  // Vida máxima conexión
    db.SetConnMaxIdleTime(2 * time.Minute)  // Tiempo máximo inactiva
}

app.Get("/users", func(c *fiber.Ctx) error {
    var users []User
    rows, err := db.Query("SELECT id, name FROM users")
    if err != nil {
        return c.Status(500).SendString("Error")
    }
    defer rows.Close()
    
    for rows.Next() {
        var user User
        rows.Scan(&user.ID, &user.Name)
        users = append(users, user)
    }
    
    return c.JSON(users)
})
```

### 59.8.5 Caching

```go
import "github.com/gofiber/fiber/v2/middleware/cache"

// Cache middleware
app.Use(cache.New(cache.Config{
    Expiration: 10 * time.Minute,
    CacheControl: true,
}))

// Cache específico por ruta
app.Get("/api/users", cache.New(cache.Config{
    Expiration:    5 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.OriginalURL()
    },
}), getUsersHandler)

// Invalidar cache manualmente
app.Delete("/users/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    // Eliminar usuario
    // ...
    // Invalidar cache
    c.Response().Header.Add("Clear-Cache", "true")
    return c.SendString("Deleted")
})
```

### 59.8.6 Fasthttp Internals

```
Ventajas de Fasthttp vs net/http:

1. ZERO-COPY I/O
   - Evita copiar datos innecesarios
   - Directa desde socket al buffer de aplicación

2. BUFFER POOLING
   - Reutiliza buffers (reduce GC)
   - Allocations ~90% menos

3. ASYNCHRONOUS PROCESSING
   - No bloquea al procesar requests
   - Goroutines eficientes

4. KEEPALIVE OPTIMIZATION
   - Reutiliza conexiones TCP
   - Reduce overhead de handshake

Arquitectura:
┌──────────────────────┐
│   Request HTTP       │
└──────────┬───────────┘
           ↓
┌──────────────────────┐
│ Fasthttp Parser      │
│ (Zero-copy)          │
└──────────┬───────────┘
           ↓
┌──────────────────────┐
│ Fiber Middleware     │
│ Pipeline             │
└──────────┬───────────┘
           ↓
┌──────────────────────┐
│ Handler Logic        │
└──────────┬───────────┘
           ↓
┌──────────────────────┐
│ Response Writer      │
│ (Zero-copy)          │
└──────────┬───────────┘
           ↓
┌──────────────────────┐
│ Response HTTP        │
└──────────────────────┘
```

### 59.8.7 Tuning de Parámetros

```go
app := fiber.New(fiber.Config{
    // Tune de red
    Network:              "tcp4",           // tcp, tcp4, tcp6, unix
    ReadTimeout:          30 * time.Second, // Timeout lectura
    WriteTimeout:         30 * time.Second, // Timeout escritura
    IdleTimeout:          5 * time.Minute,  // Timeout inactivo
    MaxRequestsPerConn:   0,                // 0 = sin límite
    
    // Tune de memoria
    BodyLimit:           4 * 1024 * 1024,   // 4MB
    Concurrency:         256 * 1024,        // Max goroutines
    
    // Tune de parseo
    JSONEncoder:         json.Marshal,      // Custom encoder
    JSONDecoder:         json.Unmarshal,    // Custom decoder
    
    // Tune de output
    NoDefaultContentType: false,             // Sin Content-Type default
    DisableKeepalive:    false,              // Mantener vivo
})
```

---

## 59.9 - Advanced Features

### 59.9.1 WebSocket

**Instalación:**
```bash
go get -u github.com/gofiber/websocket/v2
```

**Servidor WebSocket:**
```go
import "github.com/gofiber/websocket/v2"

// Upgrade a WebSocket
app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    // c.Locals puede ser usado para obtener datos
    id := c.Locals("id")
    
    for {
        messageType, message, err := c.ReadMessage()
        if err != nil {
            break
        }
        
        println("Mensaje recibido:", string(message))
        
        // Enviar echo back
        c.WriteMessage(messageType, message)
    }
}))
```

**Chat en Tiempo Real:**
```go
import "github.com/gofiber/websocket/v2"

var clients = make(map[*websocket.Conn]bool)

app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    clients[c] = true
    
    username := c.Query("username")
    
    // Anunciar entrada
    broadcast([]byte(username + " se ha conectado"))
    
    for {
        _, msg, err := c.ReadMessage()
        if err != nil {
            delete(clients, c)
            broadcast([]byte(username + " se desconectó"))
            break
        }
        
        // Broadcast del mensaje
        fullMsg := []byte(username + ": " + string(msg))
        broadcast(fullMsg)
    }
}))

func broadcast(msg []byte) {
    for client := range clients {
        client.WriteMessage(1, msg)
    }
}
```

**Cliente WebSocket (JavaScript):**
```javascript
const ws = new WebSocket("ws://localhost:3000/ws?username=John");

ws.onopen = () => {
    console.log("Conectado");
};

ws.onmessage = (event) => {
    console.log("Mensaje:", event.data);
};

ws.onerror = (error) => {
    console.error("Error:", error);
};

ws.onclose = () => {
    console.log("Desconectado");
};

// Enviar mensaje
function sendMessage(msg) {
    ws.send(msg);
}
```

### 59.9.2 Server-Sent Events (SSE)

```go
app.Get("/events", func(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("Connection", "keep-alive")
    c.Set("Access-Control-Allow-Origin", "*")
    
    // Simular envío de eventos
    for i := 0; i < 10; i++ {
        c.WriteString(fmt.Sprintf("data: Evento %d\n\n", i))
        c.Context().Response.Write(nil)  // Flush
        time.Sleep(1 * time.Second)
    }
    
    return nil
})

// Cliente (JavaScript)
// const eventSource = new EventSource("/events");
// eventSource.onmessage = (event) => {
//     console.log("SSE:", event.data);
// };
```

### 59.9.3 CORS

```go
import "github.com/gofiber/fiber/v2/middleware/cors"

// CORS abierto (desarrollo)
app.Use(cors.New())

// CORS restrictivo (producción)
app.Use(cors.New(cors.Config{
    AllowOrigins:     "https://example.com,https://app.example.com",
    AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
    AllowHeaders:     "Content-Type,Authorization,X-Custom-Header",
    ExposeHeaders:    "Content-Length,X-Total-Count",
    MaxAge:           3600,
    AllowCredentials: true,
}))

// CORS dinámico
app.Use(cors.New(cors.Config{
    AllowOriginFunc: func(origin string) bool {
        return strings.HasSuffix(origin, "example.com")
    },
}))
```

### 59.9.4 Rate Limiting

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"
import "github.com/gofiber/storage/redis"

// Rate limit simple (en memoria)
app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
}))

// Rate limit por IP
app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()  // Por IP del cliente
    },
    LimitReached: func(c *fiber.Ctx) error {
        return c.Status(429).JSON(fiber.Map{
            "error": "Demasiadas solicitudes",
            "retry_after": "60s",
        })
    },
}))

// Rate limit con Redis (distribuido)
// app.Use(limiter.New(limiter.Config{
//     Max:     100,
//     Expiration: 1 * time.Minute,
//     Storage:    redisStore,
// }))
```

### 59.9.5 JWT Authentication

```bash
go get -u github.com/golang-jwt/jwt/v5
```

**Generar Token:**
```go
import "github.com/golang-jwt/jwt/v5"

var secretKey = []byte("my-secret-key")

func generateToken(userID int) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "iat": time.Now().Unix(),
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

app.Post("/login", func(c *fiber.Ctx) error {
    // Validar credenciales
    userID := 123
    
    token, err := generateToken(userID)
    if err != nil {
        return c.Status(500).SendString("Error generando token")
    }
    
    return c.JSON(fiber.Map{
        "token": token,
    })
})
```

**Verificar Token (Middleware):**
```go
func authMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{
            "error": "No autorizado",
        })
    }
    
    token := strings.TrimPrefix(authHeader, "Bearer ")
    
    claims := jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Token inválido",
        })
    }
    
    c.Locals("userID", claims["sub"])
    return c.Next()
}

// Usar en ruta protegida
app.Get("/profile", authMiddleware, func(c *fiber.Ctx) error {
    userID := c.Locals("userID")
    return c.JSON(fiber.Map{
        "user_id": userID,
    })
})
```

---

## 59.10 - Testing & Monitoring

### 59.10.1 Unit Testing

```go
import (
    "testing"
    "github.com/gofiber/fiber/v2"
)

func TestGetUsers(t *testing.T) {
    app := fiber.New()
    
    app.Get("/users", func(c *fiber.Ctx) error {
        return c.JSON([]fiber.Map{
            {"id": 1, "name": "John"},
        })
    })
    
    // Test
    req := httptest.NewRequest("GET", "/users", nil)
    resp, _ := app.Test(req)
    
    // Assert
    if resp.StatusCode != 200 {
        t.Errorf("Expected 200, got %d", resp.StatusCode)
    }
}
```

### 59.10.2 Integration Testing

```go
func TestUserFlow(t *testing.T) {
    app := fiber.New()
    setupRoutes(app)
    
    // Crear usuario
    body := strings.NewReader(`{"name":"John","email":"john@example.com"}`)
    req := httptest.NewRequest("POST", "/users", body)
    req.Header.Set("Content-Type", "application/json")
    
    resp, _ := app.Test(req)
    if resp.StatusCode != 201 {
        t.Fatal("Failed to create user")
    }
    
    // Obtener usuarios
    req = httptest.NewRequest("GET", "/users", nil)
    resp, _ = app.Test(req)
    
    if resp.StatusCode != 200 {
        t.Fatal("Failed to get users")
    }
}
```

### 59.10.3 Health Checks

```go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "UP",
        "timestamp": time.Now(),
    })
})

app.Get("/health/detailed", func(c *fiber.Ctx) error {
    dbOK := checkDatabase()
    cacheOK := checkCache()
    
    status := "UP"
    if !dbOK || !cacheOK {
        status = "DEGRADED"
    }
    
    return c.JSON(fiber.Map{
        "status": status,
        "database": dbOK,
        "cache": cacheOK,
        "timestamp": time.Now(),
    })
})

func checkDatabase() bool {
    // Verificar conexión DB
    return true
}

func checkCache() bool {
    // Verificar conexión cache
    return true
}
```

### 59.10.4 Metrics Collection

```bash
go get -u github.com/prometheus/client_golang
```

```go
import "github.com/prometheus/client_golang/prometheus"

// Definir métricas
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "path"},
    )
)

// Middleware de métricas
func metricsMiddleware(c *fiber.Ctx) error {
    start := time.Now()
    
    err := c.Next()
    
    duration := time.Since(start).Seconds()
    status := c.Response().StatusCode()
    
    httpRequestsTotal.WithLabelValues(
        c.Method(),
        c.Path(),
        fmt.Sprintf("%d", status),
    ).Inc()
    
    httpRequestDuration.WithLabelValues(
        c.Method(),
        c.Path(),
    ).Observe(duration)
    
    return err
}

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}
```

### 59.10.5 Logging

```go
import (
    "log"
    "os"
)

func setupLogging() {
    logFile, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    
    log.SetOutput(logFile)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
    setupLogging()
    
    app.Use(func(c *fiber.Ctx) error {
        log.Printf("[%s] %s %s\n", c.Method(), c.Path(), c.IP())
        return c.Next()
    })
    
    app.Listen(":3000")
}
```

---

## 59.11 - Production & Case Studies

### 59.11.1 Docker Deployment

**Dockerfile Multi-Stage:**
```dockerfile
# Build stage
FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

# Runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 3000
CMD ["./app"]
```

**docker-compose.yml:**
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/myapp
      - LOG_LEVEL=info
    depends_on:
      - db
  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_PASSWORD=password
    volumes:
      - db_data:/var/lib/postgresql/data
volumes:
  db_data:
```

### 59.11.2 Kubernetes Setup

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fiber-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: fiber-app
  template:
    metadata:
      labels:
        app: fiber-app
    spec:
      containers:
      - name: api
        image: myregistry/fiber-app:latest
        ports:
        - containerPort: 3000
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
---
apiVersion: v1
kind: Service
metadata:
  name: fiber-app-service
spec:
  selector:
    app: fiber-app
  ports:
  - port: 80
    targetPort: 3000
  type: LoadBalancer
```

### 59.11.3 Real-World Example: REST API Completa

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/go-playground/validator/v10"
    "log"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" validate:"required,min=3"`
    Email string `json:"email" validate:"required,email"`
}

var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
}

var validate = validator.New()
var nextID = 3

func main() {
    app := fiber.New(fiber.Config{
        AppName: "Fiber API v1.0",
    })
    
    // Middleware
    app.Use(cors.New())
    app.Use(logger.New())
    
    // Rutas
    api := app.Group("/api/v1")
    
    api.Get("/users", getUsers)
    api.Post("/users", createUser)
    api.Get("/users/:id", getUser)
    api.Put("/users/:id", updateUser)
    api.Delete("/users/:id", deleteUser)
    
    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "OK"})
    })
    
    log.Fatal(app.Listen(":3000"))
}

func getUsers(c *fiber.Ctx) error {
    return c.JSON(users)
}

func createUser(c *fiber.Ctx) error {
    user := new(User)
    if err := c.BindJSON(user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    
    if err := validate.Struct(user); err != nil {
        return c.Status(422).JSON(fiber.Map{"error": err.Error()})
    }
    
    user.ID = nextID
    nextID++
    users = append(users, *user)
    
    return c.Status(201).JSON(user)
}

func getUser(c *fiber.Ctx) error {
    id, _ := c.ParamsInt("id")
    for _, user := range users {
        if user.ID == id {
            return c.JSON(user)
        }
    }
    return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
}

func updateUser(c *fiber.Ctx) error {
    id, _ := c.ParamsInt("id")
    user := new(User)
    
    if err := c.BindJSON(user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    
    for i, u := range users {
        if u.ID == id {
            user.ID = id
            users[i] = *user
            return c.JSON(user)
        }
    }
    return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
}

func deleteUser(c *fiber.Ctx) error {
    id, _ := c.ParamsInt("id")
    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            return c.SendStatus(204)
        }
    }
    return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
}
```

### 59.11.4 Scaling Strategies

```
Estrategias de Escalado:

1. VERTICAL SCALING (más recursos en servidor)
   CPU: 1 core → 16 cores
   RAM: 1GB → 8GB
   Mejor para: Procesamiento intensivo

2. HORIZONTAL SCALING (más servidores)
   1 servidor → 10 servidores
   Load Balancer distribuye tráfico
   Mejor para: APIs sin estado

3. CACHING LAYER
   Redis para resultados frecuentes
   Reduce carga DB

   Topology:
   Clients → Load Balancer
              ├─ Fiber Instance 1 → Cache (Redis)
              ├─ Fiber Instance 2 → Cache (Redis)
              └─ Fiber Instance 3 → Cache (Redis)
                      ↓
                   Database

4. DATABASE OPTIMIZATION
   - Índices apropiados
   - Query optimization
   - Read replicas
   - Sharding

5. CDN PARA STATIC CONTENT
   - Distribuir assets globalmente
   - Reduce latencia
```

### 59.11.5 Troubleshooting

**Problema: High Memory Usage**
```go
// ❌ MAL - Leak de memoria
app.Get("/data", func(c *fiber.Ctx) error {
    go func() {
        data := make([]byte, 100*1024*1024)  // Nunca se libera
        time.Sleep(10 * time.Second)
    }()
    return c.SendString("OK")
})

// ✅ BIEN
app.Get("/data", func(c *fiber.Ctx) error {
    data := make([]byte, 100*1024*1024)
    defer func() { data = nil }()
    // Usar data
    return c.SendString("OK")
})
```

**Problema: Slow Requests**
```go
// Identificar bottleneck
app.Use(func(c *fiber.Ctx) error {
    start := time.Now()
    err := c.Next()
    
    duration := time.Since(start)
    if duration > 1*time.Second {
        log.Printf("SLOW REQUEST: %s %s took %v", 
            c.Method(), c.Path(), duration)
    }
    return err
})
```

**Problema: Connection Pool Exhausted**
```go
// Síntomas: Conexiones rechazadas
// Solución: Aumentar pool size
db.SetMaxOpenConns(50)  // De 25 a 50
db.SetMaxIdleConns(10)   // De 5 a 10
```

---

## Ejercicios Progresivos

### Ejercicio 1: Simple REST API

**Objetivo:** Crear una API CRUD de productos.

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

var products []Product
var nextID = 1

func main() {
    app := fiber.New()
    app.Use(logger.New())
    
    // Implementar CRUD aquí
    
    app.Listen(":3000")
}
```

**Tareas:**
1. ✅ GET /products - Listar todos
2. ✅ POST /products - Crear
3. ✅ GET /products/:id - Obtener uno
4. ✅ PUT /products/:id - Actualizar
5. ✅ DELETE /products/:id - Eliminar

### Ejercicio 2: Middleware Chain

**Objetivo:** Crear middleware encadenado de autenticación y logging.

```go
// Implementar middleware:
// 1. Logger personalizado
// 2. Autenticación simple
// 3. Rate limiter
// 4. CORS
```

### Ejercicio 3: Validation & Error Handling

**Objetivo:** Validar datos con formato de error consistente.

```go
// Requerimientos:
// 1. Validar email, edad, etc.
// 2. Mensajes de error en español
// 3. Error handler global
```

### Ejercicio 4: WebSocket Chat

**Objetivo:** Implementar chat en tiempo real con WebSocket.

```go
// Requerimientos:
// 1. Conexión WebSocket
// 2. Broadcast a múltiples clientes
// 3. Manejo de desconexiones
```

### Ejercicio 5: Production Deployment

**Objetivo:** Desplegar app en Docker + Kubernetes.

```
Incluir:
1. Dockerfile optimizado
2. docker-compose.yml
3. Health checks
4. Métricas Prometheus
5. Logs estructurados
```

---

## Diagramas ASCII

### Request Flow Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Request                             │
│                   GET /api/v1/users/123                        │
└──────────────────────────┬──────────────────────────────────────┘
                           ↓
            ┌──────────────────────────────┐
            │   Fiber Application Instance  │
            └──────────────────┬───────────┘
                              ↓
        ┌─────────────────────────────────────────┐
        │    Middleware Pipeline                   │
        │  ┌────────────────────────────────────┐ │
        │  │ 1. Logger Middleware               │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 2. CORS Middleware                │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 3. Auth Middleware                 │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 4. RateLimit Middleware            │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 5. Validation Middleware           │ │
        │  └────────┬───────────────────────────┘ │
        └───────────┼──────────────────────────────┘
                    ↓
          ┌─────────────────────┐
          │ Route Handler       │
          │ getUserById(123)    │
          └──────────┬──────────┘
                     ↓
          ┌─────────────────────┐
          │ Process Request     │
          │ - Query DB          │
          │ - Validate Data     │
          │ - Apply Logic       │
          └──────────┬──────────┘
                     ↓
          ┌─────────────────────┐
          │ Generate Response   │
          │ JSON Encode         │
          └──────────┬──────────┘
                     ↓
        ┌─────────────────────────────────────────┐
        │    Reverse Middleware Pipeline           │
        │  ┌────────────────────────────────────┐ │
        │  │ 5. Validation → Log (opcional)    │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 4. RateLimit → Track (opcional)   │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 3. Auth → Cleanup (opcional)      │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 2. CORS → Add Headers (opcional)  │ │
        │  └────────┬───────────────────────────┘ │
        │           ↓                              │
        │  ┌────────────────────────────────────┐ │
        │  │ 1. Logger → Log (200 OK, 45ms)   │ │
        │  └────────┬───────────────────────────┘ │
        └───────────┼──────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────────────────────┐
│                       HTTP Response                             │
│                    200 OK                                       │
│              {"id":123,"name":"John"}                           │
└─────────────────────────────────────────────────────────────────┘
```

### Performance Comparison Chart

```
Fiber vs Alternatives - Throughput (req/s)

Fiber        ████████████████████████████ 330,000 req/s
Gin          ██████████████████           200,000 req/s
Echo         █████████████████            185,000 req/s
net/http     ███████████                  110,000 req/s
Express      ██████████                   95,000 req/s


Memory Usage (10k connections)

Fiber        ███        15 MB
Gin          ██████     22 MB
Echo         ███████    25 MB
Express      ████████████████████████████ 85 MB


Latency (P99)

Fiber        ▌ 0.45ms
Gin          ██ 1.2ms
Echo         ███ 1.8ms
net/http     ████ 2.4ms
Express      ████████ 5.2ms
```

### Middleware Execution Order

```
REQUEST ENTERS APPLICATION

  ↓
┌─────────────────────────────────┐
│ Global Middleware #1            │ → Logs all requests
│ (app.Use())                     │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Global Middleware #2            │ → CORS handling
│ (app.Use())                     │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Group Middleware #1             │ → Rate limiting
│ (api.Use())                     │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Route Middleware #1             │ → Authentication
│ (app.Get(..., middleware))      │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Route Middleware #2             │ → Validation
│ (app.Get(..., middleware))      │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ HANDLER                         │ → Business Logic
│ (app.Get(..., handler))         │
└─────────────────────────────────┘
  ↓
  RETURN TO MIDDLEWARE (reversed order)
  ↓
┌─────────────────────────────────┐
│ Route Middleware #2 - After     │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Route Middleware #1 - After     │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Group Middleware #1 - After     │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Global Middleware #2 - After    │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Global Middleware #1 - After    │
└─────────────────────────────────┘
  ↓

RESPONSE SENT TO CLIENT
```

---

## Comparación: Fiber vs Gin vs Echo vs Express

```
┌──────────────────────┬──────────┬──────────┬──────────┬──────────┐
│ Criterio             │ Fiber    │ Gin      │ Echo     │ Express  │
├──────────────────────┼──────────┼──────────┼──────────┼──────────┤
│ Rendimiento          │ ★★★★★   │ ★★★★☆   │ ★★★★☆   │ ★★★☆☆   │
│ Velocidad Desarrollo │ ★★★★★   │ ★★★★★   │ ★★★★☆   │ ★★★★★   │
│ Curva Aprendizaje    │ ★★★★★   │ ★★★★★   │ ★★★★☆   │ ★★★★★   │
│ Documentación        │ ★★★★☆   │ ★★★★★   │ ★★★★★   │ ★★★★★   │
│ Comunidad            │ ★★★★☆   │ ★★★★★   │ ★★★★☆   │ ★★★★★   │
│ Middleware Builtin   │ ★★★★☆   │ ★★★★★   │ ★★★★☆   │ ★★★☆☆   │
│ WebSocket            │ ★★★★☆   │ ★★★☆☆   │ ★★★☆☆   │ ★★★★☆   │
│ Validación Builtin   │ ★★★★☆   │ ★★★☆☆   │ ★★★★☆   │ ☆☆☆☆☆   │
│ Idiomatic Go         │ ★★★☆☆   │ ★★★★☆   │ ★★★★☆   │ N/A      │
│ Producción Ready     │ ★★★★★   │ ★★★★★   │ ★★★★★   │ ★★★★★   │
└──────────────────────┴──────────┴──────────┴──────────┴──────────┘

CUÁNDO USAR CADA UNO:

Fiber:
  ✅ APIs de altísimo rendimiento
  ✅ Equipos que vienen de Node.js
  ✅ Proyectos greenfield sin dependencias
  ✅ Presupuesto de recursos limitado

Gin:
  ✅ Mejor balance rendimiento-madurez
  ✅ Proyectos que requieren máximo control
  ✅ Comunidad más grande
  ✅ Casos de uso estándar

Echo:
  ✅ APIs más estructuradas
  ✅ Proyectos que requieren validación builtin
  ✅ Equipos que prefieren arquitectura clara

Express.js:
  ✅ Desarrollo rápido en Node.js
  ✅ Máxima compatibilidad con ecosistema JS
  ✅ Full-stack JavaScript
```

---

## Anti-patterns ❌ vs Best Practices ✅

### 1. Manejo de Errores

```go
// ❌ ANTI-PATTERN: Ignorar errores
app.Get("/data", func(c *fiber.Ctx) error {
    result, _ := db.Query("SELECT * FROM users")  // Ignorar error!
    return c.JSON(result)
})

// ✅ BEST PRACTICE: Manejar errores
app.Get("/data", func(c *fiber.Ctx) error {
    result, err := db.Query("SELECT * FROM users")
    if err != nil {
        log.Printf("Database error: %v", err)
        return c.Status(500).JSON(fiber.Map{
            "error": "Error al obtener datos",
        })
    }
    return c.JSON(result)
})
```

### 2. Goroutines en Handlers

```go
// ❌ ANTI-PATTERN: Goroutine que no se espera
app.Post("/process", func(c *fiber.Ctx) error {
    go func() {
        // Esta goroutine puede causar leak
        expensiveOperation()
    }()
    return c.SendString("Procesando...")
})

// ✅ BEST PRACTICE: Procesar y responder
app.Post("/process", func(c *fiber.Ctx) error {
    result := expensiveOperation()  // Sincrónico
    return c.JSON(result)
})

// ✅ ALTERNATIVA: Job queue para procesos largos
app.Post("/async-process", func(c *fiber.Ctx) error {
    jobID := submitToQueue(&Task{...})
    return c.JSON(fiber.Map{
        "job_id": jobID,
        "status": "queued",
    })
})
```

### 3. Memory Leaks en WebSocket

```go
// ❌ ANTI-PATTERN: No limpiar clientes
var clients map[*websocket.Conn]bool

app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    clients[c] = true
    // Si no se ejecuta delete, memory leak
}))

// ✅ BEST PRACTICE: Limpiar en desconexión
app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    clients[c] = true
    defer func() {
        delete(clients, c)
        c.Close()
    }()
    
    for {
        _, msg, err := c.ReadMessage()
        if err != nil {
            break  // Trigger defer
        }
        broadcast(msg)
    }
}))
```

### 4. Blocking Operations

```go
// ❌ ANTI-PATTERN: Operación bloqueante en handler
app.Get("/file", func(c *fiber.Ctx) error {
    // Esto bloquea la goroutine
    time.Sleep(5 * time.Second)
    return c.SendString("OK")
})

// ✅ BEST PRACTICE: Context con timeout
app.Get("/file", func(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    result, err := fetchDataWithContext(ctx)
    if err != nil {
        return c.Status(504).SendString("Timeout")
    }
    return c.JSON(result)
})
```

### 5. Global State

```go
// ❌ ANTI-PATTERN: Variable global compartida
var globalUserCache = make(map[int]*User)

app.Get("/users/:id", func(c *fiber.Ctx) error {
    id, _ := c.ParamsInt("id")
    // Race condition con múltiples goroutines!
    return c.JSON(globalUserCache[id])
})

// ✅ BEST PRACTICE: Usar mutex o canales
var (
    userCacheMutex sync.RWMutex
    userCache      = make(map[int]*User)
)

app.Get("/users/:id", func(c *fiber.Ctx) error {
    id, _ := c.ParamsInt("id")
    
    userCacheMutex.RLock()
    user, ok := userCache[id]
    userCacheMutex.RUnlock()
    
    if !ok {
        return c.Status(404).SendString("No encontrado")
    }
    return c.JSON(user)
})
```

### 6. Validación

```go
// ❌ ANTI-PATTERN: Validación incompleta
app.Post("/users", func(c *fiber.Ctx) error {
    user := new(User)
    c.BindJSON(user)
    
    if user.Name != "" {  // Solo verificar si no está vacío
        // Guardar
    }
    return c.SendString("OK")
})

// ✅ BEST PRACTICE: Validación robusta
app.Post("/users", func(c *fiber.Ctx) error {
    user := new(User)
    if err := c.BindJSON(user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
    }
    
    if err := validate.Struct(user); err != nil {
        return c.Status(422).JSON(fiber.Map{"error": err.Error()})
    }
    
    // Validaciones de negocio
    if userExists(user.Email) {
        return c.Status(409).JSON(fiber.Map{"error": "Email duplicado"})
    }
    
    return c.JSON(user)
})
```

### 7. Database Connections

```go
// ❌ ANTI-PATTERN: Nueva conexión por request
app.Get("/data", func(c *fiber.Ctx) error {
    db, _ := sql.Open("postgres", connString)  // Caro!
    defer db.Close()
    
    rows, _ := db.Query("SELECT * FROM users")
    return c.JSON(rows)
})

// ✅ BEST PRACTICE: Connection pool global
var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", connString)
    if err != nil {
        panic(err)
    }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
}

app.Get("/data", func(c *fiber.Ctx) error {
    rows, _ := db.Query("SELECT * FROM users")
    return c.JSON(rows)
})
```

---

## Conclusión

Fiber es un framework moderno y extremadamente performante que ofrece:

1. **Velocidad**: 3x más rápido que Gin, 12x que Express
2. **Familiaridad**: API similar a Express.js
3. **Facilidad**: Desarrollo rápido y sencillo
4. **Producción**: Listo para ambientes de producción
5. **Comunidad**: Creciente y activa

Para proyectos que requieren máximo rendimiento con una API familiar, **Fiber es la opción ideal en Go**.

---

**Fin del Capítulo 59**

*Última actualización: 2024 | Fiber v2.50+*

