# Capítulo 58: Echo - Servidor web listo para producción

## Tabla de Contenidos
1. [Introducción a Echo](#introducción)
2. [Setup y Primeros Pasos](#setup)
3. [Routing & Path Matching](#routing)
4. [Middleware Ecosystem](#middleware)
5. [Request Binding & Validation](#binding)
6. [Response Rendering](#response)
7. [HTTP Features](#http-features)
8. [Advanced Topics](#advanced)
9. [Testing Echo Applications](#testing)
10. [Performance & Monitoring](#performance)
11. [Production Deployment & Case Studies](#production)

---

## 58.1 INTRODUCCIÓN A ECHO {#introducción}

### 58.1.1 Historia y Filosofía

Echo es un framework web ultrarrápido y minimalista para Go, lanzado en 2014. Su filosofía central se resume en:

- **Rendimiento**: Enrutamiento de O(log n) con algoritmo Radix Tree
- **Simplicidad**: API limpia y coherente
- **Extensibilidad**: Middleware pattern flexible
- **Productividad**: Built-in features para aplicaciones empresariales

```
┌─────────────────────────────────────────────────────┐
│          Evolución del Ecosistema Go Web             │
├─────────────────────────────────────────────────────┤
│ 2009 │ http.HandlerFunc (stdlib)                     │
│ 2011 │ gorilla/mux (terceros)                        │
│ 2014 │ Echo Framework (ligero + rápido)              │
│ 2015 │ Gin Framework (compatibilidad Martini)        │
│ 2020 │ Fiber Framework (inspirado Express.js)        │
│ 2023 │ Echo 4.x (HTTP/2, WebSocket nativo)           │
└─────────────────────────────────────────────────────┘
```

### 58.1.2 Features Principales

**Core Features:**
- Routing rápido con Radix Tree
- Middleware potente y composable
- Binding automático (JSON, XML, Form, Query)
- Validación integrada con `validator/v10`
- WebSocket y SSE nativo
- HTTP/2 Push support
- TLS automático con Let's Encrypt

**Enterprise Features:**
- Context management eficiente
- Logging estructurado
- Error handling sofisticado
- Graceful shutdown
- Health checks
- Métricas y profiling

### 58.1.3 Enterprise Adoption

Echo es utilizado en producción por:
- **Alibaba**: Servicios de backend a escala
- **Xiaomi**: APIs internas de IoT
- **SoundCloud**: Microservicios
- **Capital One**: Servicios financieros

### 58.1.4 Casos de Uso en Producción

**Ideal para:**
```go
// 1. REST APIs de alto rendimiento
GET /api/v1/users/:id
POST /api/v1/orders
PATCH /api/v1/products/:id

// 2. Microservicios internos
GET /health
POST /metrics
GET /swagger

// 3. Real-time applications
WebSocket /ws/notifications
SSE /stream/events

// 4. Gateway/Reverse Proxy
Middleware de autenticación global
Rate limiting
CORS management
```

**NO es ideal para:**
- Aplicaciones monolíticas muy grandes (mejor Gin o framework completo)
- Aplicaciones que necesitan ORM integrado (usar GORM + Echo)
- Proyectos que requieren scaffolding automático

---

## 58.2 SETUP Y PRIMEROS PASOS {#setup}

### 58.2.1 Instalación

```bash
# Crear módulo Go
mkdir myapp && cd myapp
go mod init github.com/usuario/myapp

# Instalar Echo
go get github.com/labstack/echo/v4

# Instalar middleware común
go get github.com/labstack/echo-contrib/
go get github.com/go-playground/validator/v10
```

### 58.2.2 Echo Instance - Creación

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	// 1. Crear instancia
	e := echo.New()

	// 2. Definir rutas
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "¡Hola, Echo!")
	})

	// 3. Iniciar servidor
	e.Logger.Fatal(e.Start(":1323"))
}
```

```bash
# Ejecutar
go run main.go

# Probar
curl http://localhost:1323
# Salida: ¡Hola, Echo!
```

### 58.2.3 Routing Básico

Echo soporta todos los métodos HTTP:

```go
e := echo.New()

// GET
e.GET("/users", getUsers)

// POST
e.POST("/users", createUser)

// PUT
e.PUT("/users/:id", updateUser)

// DELETE
e.DELETE("/users/:id", deleteUser)

// PATCH
e.PATCH("/users/:id", patchUser)

// HEAD
e.HEAD("/users", headUsers)

// OPTIONS
e.OPTIONS("/users", optionsUsers)

// Método genérico
e.Match([]string{"GET", "POST"}, "/flexible", handler)
```

### 58.2.4 Start Server - Configuración

```go
package main

import (
	"github.com/labstack/echo/v4"
	"time"
)

func main() {
	e := echo.New()

	// 1. Start básico
	e.Start(":1323")

	// 2. Con puerto customizado
	e.Start(":8080")

	// 3. Con configuración detallada
	e.Server.Addr = ":1323"
	e.Server.ReadTimeout = time.Second * 15
	e.Server.WriteTimeout = time.Second * 15
	e.Logger.Fatal(e.Start(":1323"))

	// 4. Escuchar en todas las interfaces
	e.Start(":0") // OS asigna puerto disponible

	// 5. HTTPS
	e.Logger.Fatal(e.StartTLS(":443", "cert.pem", "key.pem"))

	// 6. Graceful shutdown
	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("Servidor detenido")
		}
	}()

	// Graceful shutdown en Ctrl+C
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	<-sigterm

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
```

### 58.2.5 Contexto Echo

El `echo.Context` es la pieza central:

```go
func handler(c echo.Context) error {
	// Acceder a:
	c.Request()           // *http.Request
	c.Response()          // http.ResponseWriter
	c.Path()              // Ruta registrada
	c.Param("id")         // Parámetro de ruta
	c.QueryParam("name")  // Query string
	c.Bind(&user)         // Binding automático
	c.Validate(user)      // Validación
	c.JSON(200, data)     // Respuesta JSON
	c.Error(err)          // Error handling

	return nil
}
```

---

## 58.3 ROUTING & PATH MATCHING {#routing}

### 58.3.1 Definición de Rutas

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	// 1. Ruta simple
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Inicio")
	})

	// 2. Ruta con prefijo
	e.GET("/api/users", getUsers)

	// 3. Registro y ejecución
	e.Start(":1323")
}

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Lista de usuarios",
	})
}
```

### 58.3.2 Path Parameters

```go
e := echo.New()

// 1. Parámetro simple
e.GET("/users/:id", func(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, "Usuario "+id)
})

// 2. Múltiples parámetros
e.GET("/users/:id/posts/:postId", func(c echo.Context) error {
	userID := c.Param("id")
	postID := c.Param("postId")
	return c.JSON(http.StatusOK, map[string]string{
		"user": userID,
		"post": postID,
	})
})

// 3. Path catch-all (*)
e.GET("/files/*", func(c echo.Context) error {
	filepath := c.Param("*")
	return c.String(http.StatusOK, "Archivo: "+filepath)
})

// 4. Regex matching (opcional)
e.GET("/users/:id", func(c echo.Context) error {
	// Echo automáticamente valida el parámetro
	return c.String(http.StatusOK, "OK")
})
```

**Radix Tree - Internamente Echo organiza rutas:**

```
┌─────────────────────────────────────┐
│     Radix Tree para Routing          │
├─────────────────────────────────────┤
│                                     │
│          /                          │
│         / \                         │
│        u   f                        │
│       /     \                       │
│      s       i                      │
│     /         \                     │
│    e           l                    │
│   / \           \                   │
│  r   :id         e                  │
│  |    |          |                  │
│  s    {handler}  s                  │
│  |             \ |                  │
│  {list}         {handler}           │
│                                     │
└─────────────────────────────────────┘

GET /users → handler de listado O(1)
GET /users/:id → handler con ID O(log n)
GET /files/* → handler catch-all O(log n)
```

### 58.3.3 Query Parameters

```go
e := echo.New()

// 1. Query string simple
e.GET("/search", func(c echo.Context) error {
	name := c.QueryParam("name")
	page := c.QueryParam("page")

	return c.JSON(http.StatusOK, map[string]string{
		"name": name,
		"page": page,
	})
})

// 2. Query string con múltiples valores
e.GET("/tags", func(c echo.Context) error {
	// URL: /tags?tag=go&tag=web&tag=api
	tags := c.QueryParams()["tag"]
	return c.JSON(http.StatusOK, tags)
})

// 3. Query con valor por defecto
e.GET("/products", func(c echo.Context) error {
	category := c.QueryParam("category")
	if category == "" {
		category = "all"
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}

	return c.JSON(http.StatusOK, map[string]string{
		"category": category,
		"limit":    limit,
	})
})

// 4. Query param con validación
e.GET("/users", func(c echo.Context) error {
	page := c.QueryParam("page")
	pageNum := 1

	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}

	offset := (pageNum - 1) * 10
	return c.JSON(http.StatusOK, map[string]int{
		"offset": offset,
		"limit":  10,
	})
})
```

### 58.3.4 Path Groups

```go
e := echo.New()

// 1. Agrupar con prefijo
g := e.Group("/api")
g.GET("/users", getUsers)
g.GET("/posts", getPosts)
// Rutas: GET /api/users, GET /api/posts

// 2. Versioning de API
v1 := e.Group("/api/v1")
v2 := e.Group("/api/v2")

v1.GET("/users", getUsersV1)
v2.GET("/users", getUsersV2)

// 3. Admin routes con prefijo
admin := e.Group("/admin")
admin.Use(middleware.AuthMiddleware)
admin.GET("/users", adminGetUsers)
admin.DELETE("/users/:id", adminDeleteUser)

// 4. Nested groups
api := e.Group("/api")
api.GET("/status", getStatus)

v1 := api.Group("/v1")
v1.GET("/users", listUsers)
v1.POST("/users", createUser)

v2 := api.Group("/v2")
v2.GET("/users", listUsersV2)
v2.POST("/users", createUserV2)

// Resultado:
// GET /api/status
// GET /api/v1/users
// POST /api/v1/users
// GET /api/v2/users
// POST /api/v2/users
```

### 58.3.5 Middleware per Route

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()

	// 1. Middleware global
	e.Use(middleware.Logger())

	// 2. Middleware para grupo
	public := e.Group("")
	public.GET("/", Home)

	api := e.Group("/api")
	api.Use(middleware.AuthMiddleware)
	api.GET("/users", GetUsers)

	// 3. Middleware por ruta individual
	e.GET("/public", Home)
	e.GET("/protected", Protected,
		middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey: []byte("secret"),
		}),
	)

	// 4. Multiple middleware para ruta
	e.POST("/admin/users",
		CreateUser,
		middleware.AuthMiddleware,
		middleware.AdminMiddleware,
		middleware.LoggerMiddleware,
	)

	e.Start(":1323")
}

func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Home")
}

func GetUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"user1", "user2"})
}

func Protected(c echo.Context) error {
	return c.String(http.StatusOK, "Protegido")
}

func CreateUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "created"})
}
```

### 58.3.6 Static File Serving

```go
e := echo.New()

// 1. Servir directorio completo
e.Static("/", "public")
// Archivos en public/ estarán disponibles en /

// 2. Servir un directorio específico
e.Static("/css", "assets/css")
e.Static("/js", "assets/js")
e.Static("/images", "assets/images")
// GET /css/style.css → assets/css/style.css
// GET /js/app.js → assets/js/app.js

// 3. Servir archivo único
e.File("/robots.txt", "public/robots.txt")
e.File("/sitemap.xml", "public/sitemap.xml")

// 4. Custom static handler
e.GET("/download/:file", func(c echo.Context) error {
	filename := c.Param("file")
	return c.File("downloads/" + filename)
})

// 5. Inline file serving
e.GET("/data", func(c echo.Context) error {
	return c.File("path/to/file.json")
})
```

### 58.3.7 Name Routes (Optional)

```go
e := echo.New()

// Asignar nombre a rutas
userRoute := e.GET("/users/:id", func(c echo.Context) error {
	return c.String(http.StatusOK, "Usuario")
})
userRoute.Name = "getUser"

// Generar URL (útil en templates)
e.GET("/", func(c echo.Context) error {
	url, _ := e.Reverse("getUser", "123")
	return c.String(http.StatusOK, "URL: "+url)
	// Salida: URL: /users/123
})
```

---

## 58.4 MIDDLEWARE ECOSYSTEM {#middleware}

### 58.4.1 Built-in Middleware

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// 1. Logger Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} | ${status} | ${method} ${path} | ${latency_human}` + "\n",
	}))

	// 2. Recover Middleware (maneja panics)
	e.Use(middleware.Recover())

	// 3. CORS Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Total-Count"},
		MaxAge:           300,
	}))

	// 4. Gzip Compression
	e.Use(middleware.Gzip())

	// 5. Request ID Middleware
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		TargetHeader: "X-Request-ID",
	}))

	// 6. Rate Limiter
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		100,
	)))

	// 7. Body Limit
	e.Use(middleware.BodyLimit("1M"))

	// 8. Request Timeout
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: "30s",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	e.Start(":1323")
}
```

### 58.4.2 Custom Middleware

```go
package middleware

import (
	"github.com/labstack/echo/v4"
	"time"
)

// 1. Middleware simple
func SimpleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Antes
		c.Set("user", "john")

		// Ejecutar siguiente middleware/handler
		err := next(c)

		// Después
		return err
	}
}

// 2. Middleware con parámetros
func TimingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		duration := time.Since(start)
		c.Response().Header().Set("X-Response-Time", duration.String())

		return err
	}
}

// 3. Middleware de autenticación
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return c.JSON(401, map[string]string{
				"error": "Token requerido",
			})
		}

		// Validar token (simplificado)
		if !isValidToken(token) {
			return c.JSON(401, map[string]string{
				"error": "Token inválido",
			})
		}

		c.Set("user_id", extractUserID(token))
		return next(c)
	}
}

// 4. Middleware de logging
func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		req := c.Request()
		res := c.Response()

		err := next(c)

		log.Printf(
			"%s | %s %s %d | %v\n",
			time.Now().Format("15:04:05"),
			req.Method,
			req.RequestURI,
			res.Status,
			time.Since(start),
		)

		return err
	}
}

// 5. Middleware de headers seguros
func SecurityMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		c.Response().Header().Set("Strict-Transport-Security",
			"max-age=31536000; includeSubDomains")

		return next(c)
	}
}

func isValidToken(token string) bool {
	// Implementar validación real
	return token != ""
}

func extractUserID(token string) string {
	// Implementar extracción real
	return "user123"
}
```

### 58.4.3 Middleware Chaining

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"time"
)

func main() {
	e := echo.New()

	// 1. Global middleware (ejecutados en orden)
	e.Use(middleware.Logger())      // 1er
	e.Use(middleware.Recover())     // 2do
	e.Use(SecurityMiddleware)       // 3ro

	// 2. Group middleware
	admin := e.Group("/admin")
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.AdminMiddleware)

	admin.GET("/users", getAdminUsers)
	admin.DELETE("/users/:id", deleteUser)

	// 3. Route-specific middleware
	e.GET("/public", publicHandler)
	e.GET("/protected", protectedHandler,
		middleware.AuthMiddleware,
		middleware.LoggerMiddleware,
	)

	// 4. Pipeline completo
	api := e.Group("/api")

	// Middleware global → Group middleware → Route middleware
	api.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		},
	))
	api.Use(middleware.RequestIDWithConfig(
		middleware.RequestIDConfig{
			TargetHeader: "X-Request-ID",
		},
	))

	api.GET("/data", func(c echo.Context) error {
		requestID := c.Request().Header.Get("X-Request-ID")
		return c.JSON(200, map[string]string{
			"request_id": requestID,
		})
	})

	e.Start(":1323")
}

func SecurityMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		return next(c)
	}
}

func getAdminUsers(c echo.Context) error {
	return c.JSON(200, []string{"admin1", "admin2"})
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(200, map[string]string{"deleted": id})
}

func publicHandler(c echo.Context) error {
	return c.String(200, "Público")
}

func protectedHandler(c echo.Context) error {
	return c.String(200, "Protegido")
}

// Visalizar flujo de middleware:
// GET /api/data
// ↓
// Global Logger Middleware
// ↓
// Global Recover Middleware
// ↓
// Global Security Middleware
// ↓
// Group CORS Middleware
// ↓
// Group RequestID Middleware
// ↓
// Handler: func(c echo.Context) error
// ↓
// Respuesta
```

### 58.4.4 Context Usage

```go
package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Middleware que guarda datos en context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", 123)
			c.Set("username", "john")
			return next(c)
		}
	})

	e.GET("/user", func(c echo.Context) error {
		// Recuperar datos del context
		userID := c.Get("user_id")
		username := c.Get("username")

		// Type assertion
		id, ok := userID.(int)
		if !ok {
			return c.JSON(400, map[string]string{"error": "invalid user_id"})
		}

		return c.JSON(200, map[string]interface{}{
			"id":       id,
			"username": username,
		})
	})

	e.Start(":1323")
}
```

### 58.4.5 Request/Response Middleware

```go
package middleware

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"io"
	"log"
)

// Middleware para inspeccionar request
func RequestInspectorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Leer y reguardar body
		bodyBytes, _ := io.ReadAll(c.Request().Body)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("Request: %s %s\n", c.Request().Method, c.Request().URL.Path)
		log.Printf("Headers: %v\n", c.Request().Header)
		log.Printf("Body: %s\n", string(bodyBytes))

		return next(c)
	}
}

// Middleware para inspeccionar response
func ResponseInspectorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Capturar respuesta
		originalWriter := c.Response().Writer

		// Custom writer para capturar response
		buf := new(bytes.Buffer)
		c.Response().Writer = &responseWriter{
			ResponseWriter: originalWriter,
			buf:            buf,
		}

		err := next(c)

		log.Printf("Response Status: %d\n", c.Response().Status)
		log.Printf("Response Body: %s\n", buf.String())

		return err
	}
}

type responseWriter struct {
	echo.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}
```

---

## 58.5 REQUEST BINDING & VALIDATION {#binding}

### 58.5.1 JSON Binding

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	e := echo.New()

	// 1. Binding JSON simple
	e.POST("/users", func(c echo.Context) error {
		user := new(User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid JSON",
			})
		}

		return c.JSON(http.StatusCreated, user)
	})

	// 2. Binding con validación
	e.POST("/users/validated", func(c echo.Context) error {
		user := new(User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if user.Name == "" || user.Email == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Name y Email son requeridos",
			})
		}

		return c.JSON(http.StatusCreated, user)
	})

	// 3. Binding a slice
	e.POST("/users/batch", func(c echo.Context) error {
		var users []User

		if err := c.Bind(&users); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusCreated, users)
	})

	// 4. Binding y BindJSON explícito
	e.POST("/users/explicit", func(c echo.Context) error {
		user := new(User)

		if err := c.BindJSON(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusCreated, user)
	})

	e.Start(":1323")
}
```

```bash
# Test
curl -X POST http://localhost:1323/users \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"John","email":"john@example.com"}'
```

### 58.5.2 Form Binding

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Remember string `form:"remember"`
}

func main() {
	e := echo.New()

	// 1. Form binding
	e.POST("/login", func(c echo.Context) error {
		form := new(LoginForm)

		if err := c.Bind(form); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, form)
	})

	// 2. Form file upload
	e.POST("/upload", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		// Guardar archivo
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Aquí guardar el archivo
		return c.JSON(http.StatusOK, map[string]string{
			"filename": file.Filename,
			"size":     string(rune(file.Size)),
		})
	})

	// 3. Multiple file upload
	e.POST("/upload-multiple", func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		files := form.File["files"]

		return c.JSON(http.StatusOK, map[string]interface{}{
			"count": len(files),
			"files": files,
		})
	})

	e.Start(":1323")
}
```

### 58.5.3 XML Binding

```go
package main

import (
	"encoding/xml"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Product struct {
	XMLName xml.Name `xml:"product"`
	ID      int      `xml:"id"`
	Name    string   `xml:"name"`
	Price   float64  `xml:"price"`
}

func main() {
	e := echo.New()

	// 1. XML Binding
	e.POST("/products", func(c echo.Context) error {
		product := new(Product)

		if err := c.BindXML(product); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.XML(http.StatusCreated, product)
	})

	e.Start(":1323")
}
```

```bash
curl -X POST http://localhost:1323/products \
  -H "Content-Type: application/xml" \
  -d '<product><id>1</id><name>Laptop</name><price>999.99</price></product>'
```

### 58.5.4 Custom Binders

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type CustomData struct {
	Date  time.Time
	Count int
}

// Custom binder
type customBinder struct{}

func (cb *customBinder) Bind(req *http.Request, i interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	// Binding custom
	data := i.(*CustomData)
	dateStr := req.FormValue("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	data.Date = date
	return nil
}

func (cb *customBinder) BindPathParams(c echo.Context, i interface{}) error {
	return nil
}

func main() {
	e := echo.New()

	// Registrar custom binder
	e.Binder = &customBinder{}

	e.POST("/custom", func(c echo.Context) error {
		data := new(CustomData)

		if err := c.Bind(data); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, data)
	})

	e.Start(":1323")
}
```

### 58.5.5 Validation Integration

```go
package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type User struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=18,max=120"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func main() {
	e := echo.New()

	// Registrar validador
	e.Validator = &CustomValidator{
		validator: validator.New(),
	}

	// 1. Validación básica
	e.POST("/users", func(c echo.Context) error {
		user := new(User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := c.Validate(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusCreated, user)
	})

	// 2. Validación manual
	e.POST("/users-manual", func(c echo.Context) error {
		user := new(User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		// Validar manualmente
		if len(user.Name) < 3 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Name debe tener al menos 3 caracteres",
			})
		}

		if user.Age < 18 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Debe ser mayor de 18 años",
			})
		}

		return c.JSON(http.StatusCreated, user)
	})

	e.Start(":1323")
}
```

### 58.5.6 Error Formatting

```go
package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func formatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, ValidationError{
				Field: fieldError.Field(),
				Message: formatErrorMessage(fieldError),
			})
		}
	}

	return errors
}

func formatErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.ActualTag() {
	case "required":
		return fieldError.Field() + " es requerido"
	case "email":
		return fieldError.Field() + " debe ser un email válido"
	case "min":
		return fieldError.Field() + " debe tener mínimo " + fieldError.Param() + " caracteres"
	case "max":
		return fieldError.Field() + " debe tener máximo " + fieldError.Param() + " caracteres"
	default:
		return "Validación inválida en " + fieldError.Field()
	}
}

func main() {
	e := echo.New()

	type User struct {
		Name  string `json:"name" validate:"required,min=3"`
		Email string `json:"email" validate:"required,email"`
	}

	e.Validator = &CustomValidator{
		validator: validator.New(),
	}

	e.POST("/users", func(c echo.Context) error {
		user := new(User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := c.Validate(user); err != nil {
			validationErrors := formatValidationErrors(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"errors": validationErrors,
			})
		}

		return c.JSON(http.StatusCreated, user)
	})

	e.Start(":1323")
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
```

---

## 58.6 RESPONSE RENDERING {#response}

### 58.6.1 JSON Responses

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	e := echo.New()

	// 1. JSON simple
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: "OK",
		})
	})

	// 2. JSON con datos
	e.GET("/users", func(c echo.Context) error {
		users := []map[string]string{
			{"id": "1", "name": "John"},
			{"id": "2", "name": "Jane"},
		}

		return c.JSON(http.StatusOK, Response{
			Status: "success",
			Data:   users,
		})
	})

	// 3. JSON con headers personalizados
	e.GET("/with-headers", func(c echo.Context) error {
		c.Response().Header().Set("X-Total-Count", "100")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Con headers",
		})
	})

	e.Start(":1323")
}
```

### 58.6.2 XML Responses

```go
e := echo.New()

type Item struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

e.GET("/items", func(c echo.Context) error {
	items := []Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
	}

	return c.XML(http.StatusOK, items)
})
```

### 58.6.3 HTML Templates

```go
package main

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Items []string
}

func main() {
	e := echo.New()

	// 1. Template simple
	t := template.Must(template.New("").Parse(`
		<h1>{{.Title}}</h1>
		<ul>
		{{range .Items}}
			<li>{{.}}</li>
		{{end}}
		</ul>
	`))

	e.Renderer = &TemplateRenderer{
		templates: t,
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "template", Page{
			Title: "Mi Página",
			Items: []string{"Item 1", "Item 2"},
		})
	})

	// 2. Template desde archivo
	htmlFile := template.Must(template.ParseFiles("views/index.html"))

	e.GET("/page", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.Start(":1323")
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w http.ResponseWriter, name string,
	data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
```

### 58.6.4 File Serving

```go
e := echo.New()

// 1. Servir archivo
e.GET("/download", func(c echo.Context) error {
	return c.File("path/to/file.pdf")
})

// 2. Servir con nombre personalizado
e.GET("/export", func(c echo.Context) error {
	return c.Attachment("path/to/data.csv", "export.csv")
})

// 3. Inline file
e.GET("/view", func(c echo.Context) error {
	return c.Inline("path/to/image.jpg", "image.jpg")
})
```

### 58.6.5 Streaming

```go
e := echo.New()

// 1. Stream simple
e.GET("/stream", func(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")

	for i := 0; i < 10; i++ {
		c.Response().Write([]byte("data: chunk " + string(rune(i)) + "\n\n"))
	}

	return nil
})

// 2. Server-Sent Events
e.GET("/events", func(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")

	done := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			c.Response().Write([]byte("data: " + string(rune(i)) + "\n\n"))
			time.Sleep(1 * time.Second)
		}
		done <- true
	}()

	<-done
	return nil
})
```

### 58.6.6 Custom Rendering

```go
type CustomRenderer struct{}

func (cr *CustomRenderer) Render(w http.ResponseWriter, name string,
	data interface{}, c echo.Context) error {

	// Renderizar customizado
	w.Header().Set("Content-Type", "application/custom")
	w.Write([]byte("Custom: " + name))

	return nil
}

e := echo.New()
e.Renderer = &CustomRenderer{}

e.GET("/custom", func(c echo.Context) error {
	return c.Render(http.StatusOK, "test", nil)
})
```

---

## 58.7 HTTP FEATURES {#http-features}

### 58.7.1 Content Negotiation

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Data struct {
	Message string `json:"message" xml:"message"`
}

func main() {
	e := echo.New()

	e.GET("/data", func(c echo.Context) error {
		data := Data{Message: "Hello"}

		// Content negotiation automática
		return c.JSON(http.StatusOK, data)
		// Responde con JSON si el client lo solicita
	})

	// Manual content negotiation
	e.GET("/flexible", func(c echo.Context) error {
		data := Data{Message: "Hello"}

		accept := c.Request().Header.Get("Accept")

		if accept == "application/xml" {
			return c.XML(http.StatusOK, data)
		}

		return c.JSON(http.StatusOK, data)
	})

	e.Start(":1323")
}
```

### 58.7.2 Compression

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// 1. Gzip compression
	e.Use(middleware.Gzip())

	// 2. Con configuración
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5, // 1-9, por defecto 6
	}))

	e.GET("/data", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"data": "Será comprimido automáticamente",
		})
	})

	e.Start(":1323")
}

// Test:
// curl -H "Accept-Encoding: gzip" http://localhost:1323/data
```

### 58.7.3 Range Requests

```go
e := echo.New()

e.GET("/video", func(c echo.Context) error {
	// HTTP 206 Partial Content
	// Útil para descargas de archivos grandes

	return c.File("video.mp4")
	// Echo maneja Range requests automáticamente
})
```

### 58.7.4 ETag & Cache-Control

```go
import "crypto/md5"

e := echo.New()

e.GET("/data", func(c echo.Context) error {
	data := "Este es el contenido"

	// 1. Generar ETag
	hash := md5.Sum([]byte(data))
	etag := fmt.Sprintf("%x", hash)

	// Comparar con cliente
	if c.Request().Header.Get("If-None-Match") == etag {
		return c.NoContent(http.StatusNotModified)
	}

	// 2. Establecer headers de cache
	c.Response().Header().Set("ETag", etag)
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	c.Response().Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	return c.String(http.StatusOK, data)
})
```

### 58.7.5 WebSocket Support

```go
package main

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	e := echo.New()

	e.GET("/ws", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		// Echo loop
		for {
			var msg string
			if err := ws.ReadMessage(&msg); err != nil {
				log.Println(err)
				break
			}

			if err := ws.WriteMessage(websocket.TextMessage,
				[]byte("Echo: "+msg)); err != nil {
				log.Println(err)
				break
			}
		}

		return nil
	})

	e.Start(":1323")
}

// Client test
// curl --include \
//   --no-buffer \
//   --header "Connection: Upgrade" \
//   --header "Upgrade: websocket" \
//   --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
//   --header "Sec-WebSocket-Version: 13" \
//   http://localhost:1323/ws
```

### 58.7.6 Server-Sent Events

```go
e := echo.New()

e.GET("/events", func(c echo.Context) error {
	// Headers de SSE
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	// Stream de eventos
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("data: Evento %d\n\n", i)
		c.Response().Write([]byte(msg))
		c.Response().Flush()

		time.Sleep(1 * time.Second)
	}

	return nil
})
```

---

## 58.8 ADVANCED TOPICS {#advanced}

### 58.8.1 Error Handling

```go
package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type CustomError struct {
	Code    int
	Message string
	Details interface{}
}

func (ce *CustomError) Error() string {
	return ce.Message
}

func main() {
	e := echo.New()

	// 1. Error handler global
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Error interno del servidor"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		if ce, ok := err.(*CustomError); ok {
			code = ce.Code
			message = ce.Message
		}

		c.JSON(code, map[string]interface{}{
			"error":   message,
			"status":  code,
		})
	}

	// 2. Lanzar errores
	e.GET("/error", func(c echo.Context) error {
		return &CustomError{
			Code:    http.StatusNotFound,
			Message: "Usuario no encontrado",
			Details: map[string]string{"id": "123"},
		}
	})

	// 3. Echo built-in errors
	e.GET("/notfound", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "Recurso no encontrado")
	})

	e.Start(":1323")
}
```

### 58.8.2 Custom Error Handlers

```go
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Error interno del servidor"
	var details interface{}

	// Tipo de error
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			message = m
		}
	} else if _, ok := err.(*json.SyntaxError); ok {
		code = http.StatusBadRequest
		message = "JSON inválido"
	} else if _, ok := err.(*strconv.NumError); ok {
		code = http.StatusBadRequest
		message = "Número inválido"
	} else {
		log.Printf("Error desconocido: %v", err)
	}

	// Respuesta de error estructurada
	resp := map[string]interface{}{
		"status":  code,
		"error":   message,
		"details": details,
	}

	// En desarrollo, incluir stack trace
	if os.Getenv("ENV") == "development" {
		resp["stack"] = fmt.Sprintf("%+v", err)
	}

	c.JSON(code, resp)
}

e := echo.New()
e.HTTPErrorHandler = customHTTPErrorHandler
```

### 58.8.3 Graceful Shutdown

```go
package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Iniciar servidor en goroutine
	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("Servidor cerrado")
		}
	}()

	// Esperar signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	println("Servidor cerrado correctamente")
}
```

### 58.8.4 Server Configuration

```go
e := echo.New()

// Configuración del servidor
e.Server.Addr = ":1323"
e.Server.ReadTimeout = time.Second * 15
e.Server.WriteTimeout = time.Second * 15
e.Server.IdleTimeout = time.Second * 60

// Logger
e.Logger.SetOutput(os.Stdout)
e.Logger.SetPrefix("[ECHO]")

// Debug mode
e.Debug = true

// Hide banner
e.HideBanner = true

e.Start(":1323")
```

### 58.8.5 TLS/HTTPS Setup

```go
package main

import (
	"crypto/tls"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "HTTPS OK")
	})

	// 1. StartTLS con archivos
	e.Logger.Fatal(e.StartTLS(
		":443",
		"path/to/cert.pem",
		"path/to/key.pem",
	))

	// 2. Con configuración TLS
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	e.Server.TLSConfig = tlsConfig
	e.Logger.Fatal(e.StartTLS(
		":443",
		"cert.pem",
		"key.pem",
	))
}
```

### 58.8.6 HTTP/2 Support

```go
// Echo soporta HTTP/2 automáticamente con StartTLS

e := echo.New()

// El servidor HTTP/2 está habilitado por defecto con TLS
e.Logger.Fatal(e.StartTLS(":443", "cert.pem", "key.pem"))

// Test:
// curl -I --http2 https://localhost:443/
```

---

## 58.9 TESTING ECHO APPLICATIONS {#testing}

### 58.9.1 httptest Integration

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGET(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	// Define handler
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	}

	// Ejecutar
	handler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello", w.Body.String())
}
```

### 58.9.2 Test Helpers

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
)

func setupTestServer() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Home")
	})

	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(200, map[string]string{
			"id": id,
		})
	})

	e.POST("/users", func(c echo.Context) error {
		return c.JSON(201, map[string]string{
			"status": "created",
		})
	})

	return e
}

func makeRequest(e *echo.Echo, method, path string) *http.Response {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	return w.Result()
}

// Uso en tests
func TestRouter(t *testing.T) {
	e := setupTestServer()

	// Test GET /
	resp := makeRequest(e, "GET", "/")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test GET /users/:id
	resp = makeRequest(e, "GET", "/users/123")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### 58.9.3 Mocking

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock de servicio
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(id string) (map[string]string, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestUserHandler(t *testing.T) {
	// Crear mock
	mockService := new(MockUserService)
	mockService.On("GetUser", "123").Return(
		map[string]string{"id": "123", "name": "John"},
		nil,
	)

	// Setup handler
	e := echo.New()
	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		user, _ := mockService.GetUser(id)
		return c.JSON(http.StatusOK, user)
	})

	// Test
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	// Assert
	mockService.AssertCalled(t, "GetUser", "123")
}
```

### 58.9.4 Table-Driven Tests

```go
func TestHandlers(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Home",
		},
		{
			name:           "GET /users/123",
			method:         http.MethodGet,
			path:           "/users/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /users",
			method:         http.MethodPost,
			path:           "/users",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "GET /notfound",
			method:         http.MethodGet,
			path:           "/notfound",
			expectedStatus: http.StatusNotFound,
		},
	}

	e := setupTestServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}
```

### 58.9.5 Full App Testing

```go
package main

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestApp() *echo.Echo {
	e := echo.New()

	e.POST("/users", func(c echo.Context) error {
		var user struct {
			Name string `json:"name"`
		}
		if err := c.Bind(&user); err != nil {
			return err
		}
		return c.JSON(201, user)
	})

	return e
}

func TestCreateUser(t *testing.T) {
	e := createTestApp()

	// Crear request con JSON
	body := []byte(`{"name":"John"}`)
	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	// Validar respuesta
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "John", resp["name"])
}
```

---

## 58.10 PERFORMANCE & MONITORING {#performance}

### 58.10.1 Performance Characteristics

```
┌────────────────────────────────────────────────────┐
│   Echo Performance vs Otros Frameworks (benchmarks) │
├────────────────────────────────────────────────────┤
│                                                    │
│ Throughput (requests/seg):                         │
│ Echo:      ████████████████████ 120,000            │
│ Gin:       ████████████████████ 118,000            │
│ Fiber:     ██████████████████ 110,000              │
│ Chi:       ███████████████ 95,000                  │
│ Gorilla:   ████████████ 75,000                     │
│                                                    │
│ Latency promedio (ms):                             │
│ Echo:      ▂ 0.85ms                                │
│ Gin:       ▂ 0.87ms                                │
│ Fiber:     ▃ 0.95ms                                │
│ Chi:       ▄ 1.2ms                                 │
│ Gorilla:   ▅ 1.5ms                                 │
│                                                    │
│ Memory (MB):                                       │
│ Echo:      ▁ 15MB                                  │
│ Gin:       ▂ 16MB                                  │
│ Fiber:     ▃ 18MB                                  │
│ Chi:       ▄ 22MB                                  │
│ Gorilla:   ▅ 28MB                                  │
│                                                    │
└────────────────────────────────────────────────────┘

Ventajas de Echo:
✓ Radix tree O(log n)
✓ Connection pooling eficiente
✓ GC optimizado
✓ Bajo overhead de middleware
```

### 58.10.2 Memory Profiling

```go
package main

import (
	"github.com/labstack/echo/v4"
	_ "net/http/pprof"
	"runtime"
	"runtime/pprof"
	"os"
)

func main() {
	// CPU profiling
	cpuProfile, _ := os.Create("cpu.prof")
	defer cpuProfile.Close()
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	e.Start(":1323")

	// Memory profiling
	memProfile, _ := os.Create("mem.prof")
	defer memProfile.Close()
	runtime.GC()
	pprof.WriteHeapProfile(memProfile)
}

// Análisis:
// go tool pprof cpu.prof
// go tool pprof mem.prof
// (pprof) top - muestra funciones con más uso
// (pprof) list functionName - detalles de función
```

### 58.10.3 Health Checks

```go
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthStatus struct {
	Status string `json:"status"`
	Uptime int64  `json:"uptime"`
}

func main() {
	e := echo.New()
	startTime := time.Now()

	// Health check
	e.GET("/health", func(c echo.Context) error {
		uptime := time.Since(startTime).Milliseconds()

		return c.JSON(http.StatusOK, HealthStatus{
			Status: "healthy",
			Uptime: uptime,
		})
	})

	// Readiness check (BD, etc.)
	e.GET("/ready", func(c echo.Context) error {
		// Verificar dependencias
		if !checkDatabase() {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ready",
		})
	})

	e.Start(":1323")
}

func checkDatabase() bool {
	// Implementar verificación real
	return true
}
```

### 58.10.4 Metrics Collection

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requests HTTP",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duración de requests en segundos",
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func main() {
	e := echo.New()

	// Middleware de métricas
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start).Seconds()
			requestCount.WithLabelValues(
				c.Request().Method,
				c.Request().URL.Path,
				string(rune(c.Response().Status)),
			).Inc()

			requestDuration.WithLabelValues(
				c.Request().Method,
				c.Request().URL.Path,
			).Observe(duration)

			return err
		}
	})

	// Endpoint de métricas
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Start(":1323")
}
```

### 58.10.5 Logging Integration

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	e := echo.New()

	// Logger
	logFile, err := os.OpenFile("app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	e.Logger.SetOutput(logFile)
	e.Logger.SetLevel(log.INFO)

	// Logger middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} | ${status} | ${method} ${path} ` +
			`| ${latency_human} | ${error}` + "\n",
		Output: logFile,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	e.Start(":1323")
}
```

---

## 58.11 PRODUCTION DEPLOYMENT & CASE STUDIES {#production}

### 58.11.1 Docker Deployment

**Dockerfile:**

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 1323

CMD ["./app"]
```

**docker-compose.yml:**

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "1323:1323"
    environment:
      - ENV=production
      - PORT=1323
    depends_on:
      - db
  
  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**Build y run:**

```bash
docker-compose up --build
```

### 58.11.2 Kubernetes Manifests

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: echo-app
  template:
    metadata:
      labels:
        app: echo-app
    spec:
      containers:
      - name: app
        image: echo-app:1.0
        ports:
        - containerPort: 1323
        livenessProbe:
          httpGet:
            path: /health
            port: 1323
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 1323
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: echo-app-service
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 1323
  selector:
    app: echo-app
```

### 58.11.3 Real-World Case Studies

**Case 1: REST API e-commerce (50k RPS)**

```go
// API con Echo
// - 100+ endpoints
// - Rate limiting
// - JWT authentication
// - Database connection pooling
// - Redis caching

e := echo.New()

// Middleware global
e.Use(middleware.Logger())
e.Use(middleware.Recover())
e.Use(middleware.CORSWithConfig(cors))
e.Use(middleware.JWTWithConfig(jwt))
e.Use(middleware.RateLimiter(limiter))

// API routes
api := e.Group("/api/v1")
api.GET("/products", listProducts)        // 10k RPS
api.GET("/products/:id", getProduct)      // 8k RPS
api.POST("/orders", createOrder)          // 5k RPS
api.GET("/users/profile", getProfile)     // 3k RPS

// Cache
e.Use(cacheMiddleware)

// Total: 50k+ RPS manejados
e.Start(":1323")
```

**Resultados:**
- Latency: p99 < 100ms
- Throughput: 50k RPS sostenido
- Memory: 250MB
- CPU: 4 cores @ 60%

**Case 2: Microservicios internos**

```go
// API interna para comunicación entre servicios
// - 200+ rutas
// - gRPC compatibility
// - Health checks
// - Metrics

e := echo.New()

// Internal routes
internal := e.Group("/internal")
internal.GET("/health", health)
internal.GET("/metrics", metrics)
internal.POST("/batch", batchProcess)

// Service-to-service communication
e.POST("/rpc", rpcHandler)
e.GET("/stream/:resource", streamData)

// Result: < 50ms latency inter-service
e.Start(":1323")
```

**Case 3: WebSocket real-time notifications**

```go
// Notification hub
// - 1000+ conexiones simultáneas
// - Broadcasting
// - User-specific messages

var (
	clients = make(map[*Client]bool)
	broadcast = make(chan *Message, 256)
	register = make(chan *Client)
	unregister = make(chan *Client)
)

e.GET("/ws/:user_id", handleWebSocket)

// Hub goroutine
go runHub()

// Soporta 10k+ conexiones con 2 cores
```

### 58.11.4 Migration desde Gin/Fiber

**De Gin a Echo:**

```go
// Gin
router := gin.Default()
router.GET("/users/:id", func(c *gin.Context) {
	id := c.Param("id")
})

// Echo (muy similar)
e := echo.New()
e.GET("/users/:id", func(c echo.Context) error {
	id := c.Param("id")
	return nil
})

// Key differences:
// Gin: c.JSON(200, data)  → c.JSON(200, data)
// Echo: return c.JSON(200, data)  ← Echo requiere return
// Gin: c.Error(err)  →  Echo: return err
```

**De Fiber a Echo:**

```go
// Fiber
app := fiber.New()
app.Post("/users", func(c *fiber.Ctx) error {
	var user User
	c.BodyParser(&user)
	return c.JSON(user)
})

// Echo
e := echo.New()
e.POST("/users", func(c echo.Context) error {
	user := new(User)
	c.Bind(user)
	return c.JSON(200, user)
})

// Principales cambios:
// - BodyParser() → Bind()
// - Status() → statusCode como parámetro
// - SendFile() → File()
```

### 58.11.5 Troubleshooting

```go
// Problema 1: Connection pool exhausted
// Solución:
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// Problema 2: Memory leak en middleware
// Solución: Limpiar resources
func memoryLeakMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			// Cleanup
			c.Set("user", nil)
		}()
		return next(c)
	}
}

// Problema 3: Slow requests
// Solución: Timeout middleware
e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	Timeout: "30s",
}))

// Problema 4: High CPU usage
// Solución: CPU profiling
import _ "net/http/pprof"
// Acceder a http://localhost:1323/debug/pprof

// Problema 5: CORS errors
// Solución: Configurar CORS correctamente
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"https://example.com"},
	AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	AllowHeaders: []string{"Content-Type", "Authorization"},
}))
```

---

## 58.12 EJERCICIOS PROGRESIVOS {#ejercicios}

### Ejercicio 1: API REST Simple

**Objetivo:** Crear un CRUD básico para gestionar libros

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Year     int    `json:"year"`
	Pages    int    `json:"pages"`
}

var books = []Book{
	{ID: 1, Title: "Go Programming", Author: "John Doe", Year: 2020, Pages: 500},
	{ID: 2, Title: "Web Development", Author: "Jane Smith", Year: 2021, Pages: 450},
}

var nextID = 3

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Routes
	e.GET("/books", getBooks)
	e.GET("/books/:id", getBook)
	e.POST("/books", createBook)
	e.PUT("/books/:id", updateBook)
	e.DELETE("/books/:id", deleteBook)

	e.Start(":1323")
}

func getBooks(c echo.Context) error {
	return c.JSON(http.StatusOK, books)
}

func getBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, b := range books {
		if b.ID == id {
			return c.JSON(http.StatusOK, b)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "Libro no encontrado",
	})
}

func createBook(c echo.Context) error {
	book := new(Book)

	if err := c.Bind(book); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	book.ID = nextID
	nextID++
	books = append(books, *book)

	return c.JSON(http.StatusCreated, book)
}

func updateBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book := new(Book)

	if err := c.Bind(book); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	for i, b := range books {
		if b.ID == id {
			book.ID = id
			books[i] = *book
			return c.JSON(http.StatusOK, book)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "Libro no encontrado",
	})
}

func deleteBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			return c.JSON(http.StatusOK, map[string]string{
				"message": "Libro eliminado",
			})
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "Libro no encontrado",
	})
}
```

### Ejercicio 2: Custom Middleware

**Objetivo:** Crear middleware de autenticación y logging

```go
package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

// Middleware de autenticación
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Token requerido",
			})
		}

		// Validar token (simplificado)
		if !isValidToken(token) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Token inválido",
			})
		}

		c.Set("user_id", "123")
		return next(c)
	}
}

// Middleware de logging
func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()

		err := next(c)

		duration := time.Since(start)
		log.Printf(
			"%s | %s %s | %d | %v\n",
			time.Now().Format("15:04:05"),
			req.Method,
			req.URL.Path,
			c.Response().Status,
			duration,
		)

		return err
	}
}

// Middleware de rate limiting
func rateLimitMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	var count int
	var resetTime time.Time

	return func(c echo.Context) error {
		if time.Now().After(resetTime) {
			count = 0
			resetTime = time.Now().Add(1 * time.Minute)
		}

		if count >= 100 {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Rate limit excedido",
			})
		}

		count++
		return next(c)
	}
}

func main() {
	e := echo.New()

	// Aplicar middleware global
	e.Use(loggingMiddleware)

	// Rutas públicas
	e.POST("/login", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"token": "abc123",
		})
	})

	// Rutas protegidas
	api := e.Group("/api")
	api.Use(authMiddleware)

	api.GET("/profile", func(c echo.Context) error {
		userID := c.Get("user_id")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id": userID,
		})
	})

	e.Start(":1323")
}

func isValidToken(token string) bool {
	return token == "Bearer abc123"
}
```

### Ejercicio 3: Validation & Error Handling

**Objetivo:** Implementar validación completa con error handling sofisticado

```go
package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type User struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=18,max=120"`
	Phone    string `json:"phone" validate:"omitempty,len=10"`
	Password string `json:"password" validate:"required,min=8"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status string            `json:"status"`
	Errors []ValidationError `json:"errors"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func formatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   fieldError.Field(),
				Message: getErrorMessage(fieldError),
			})
		}
	}

	return errors
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.ActualTag() {
	case "required":
		return fe.Field() + " es requerido"
	case "email":
		return "Email inválido"
	case "min":
		return fe.Field() + " debe tener al menos " + fe.Param() + " caracteres"
	case "max":
		return fe.Field() + " no puede exceder " + fe.Param() + " caracteres"
	case "len":
		return fe.Field() + " debe tener exactamente " + fe.Param() + " caracteres"
	default:
		return "Validación inválida en " + fe.Field()
	}
}

func main() {
	e := echo.New()

	e.Validator = &CustomValidator{
		validator: validator.New(),
	}

	e.POST("/users", func(c echo.Context) error {
		user := new(User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "JSON inválido",
			})
		}

		if err := c.Validate(user); err != nil {
			errors := formatValidationErrors(err)
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Status: "error",
				Errors: errors,
			})
		}

		return c.JSON(http.StatusCreated, user)
	})

	e.Start(":1323")
}
```

### Ejercicio 4: WebSocket Server

**Objetivo:** Crear un servidor de chat con WebSocket

```go
package main

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Cliente conectado. Total: %d\n", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Cliente desconectado. Total: %d\n", len(h.clients))
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func main() {
	e := echo.New()
	hub := newHub()
	go hub.run()

	e.GET("/ws", func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
		}

		hub.register <- client

		go client.readPump()
		go client.writePump()

		return nil
	})

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
<!DOCTYPE html>
<html>
<body>
<input id="msg" type="text">
<button onclick="send()">Enviar</button>
<div id="output"></div>
<script>
const ws = new WebSocket('ws://localhost:1323/ws');
ws.onmessage = e => {
	document.getElementById('output').innerHTML += 
		'<p>' + e.data + '</p>';
};
function send() {
	const msg = document.getElementById('msg').value;
	ws.send(msg);
	document.getElementById('msg').value = '';
}
</script>
</body>
</html>
		`)
	})

	e.Start(":1323")
}
```

### Ejercicio 5: Production Deployment

**Objetivo:** Deployar una aplicación Echo con Docker, Kubernetes y CI/CD

**main.go:**

```go
package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Echo API v1.0")
	})

	// Graceful shutdown
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
```

**go.mod:**

```
module github.com/usuario/myapp

go 1.21

require github.com/labstack/echo/v4 v4.11.0
```

**Dockerfile:**

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:1323/health || exit 1

EXPOSE 1323
CMD ["./app"]
```

**Hacer build y deploy:**

```bash
# Build local
docker build -t echo-app:1.0 .

# Run
docker run -p 1323:1323 echo-app:1.0

# Deploy en Kubernetes
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Verificar
kubectl get pods
kubectl logs deployment/echo-app
```

---

## 58.13 COMPARACIONES TÉCNICAS

### Echo vs Gin

| Aspecto | Echo | Gin |
|--------|------|-----|
| Routing | Radix Tree O(log n) | Trie-based |
| Performance | 120k RPS | 118k RPS |
| Middleware | Muy flexible | Menos flexible |
| Validación | Built-in | Requiere integración |
| WebSocket | Nativo | Con librería |
| Async | Excelente | Muy bueno |
| Comunidad | Grande | Muy grande |
| Curva aprendizaje | Moderada | Fácil |

### Echo vs Fiber

| Aspecto | Echo | Fiber |
|--------|------|-------|
| Estilo | Go estándar | Express.js style |
| Performance | 120k RPS | 110k RPS |
| Memory | 15MB | 18MB |
| Context | echo.Context | *fiber.Ctx |
| Error handling | Devolver errores | c.Next() |
| Testing | Nativo httptest | Custom helpers |
| Database | GORM fácil | GORM fácil |
| Madurez | Muy madura | En crecimiento |

### Echo vs Chi

| Aspecto | Echo | Chi |
|--------|------|-----|
| Tamaño | ~500 KB | ~300 KB |
| Features | Muchas | Minimalista |
| Binding | Automático | Manual |
| Validación | Built-in | Manual |
| Performance | 120k RPS | 95k RPS |
| Uso | Producción | Ligero/API |
| Madurez | Enterprise-grade | Muy madura |

---

## 58.14 ANTI-PATTERNS ❌ vs BEST PRACTICES ✅

### Pattern 1: Binding & Validation

```go
// ❌ MAL: Manual parsing
e.POST("/users", func(c echo.Context) error {
	data, _ := ioutil.ReadAll(c.Request().Body)
	var user map[string]interface{}
	json.Unmarshal(data, &user)
	return c.JSON(200, user)
})

// ✅ BIEN: Echo binding
e.POST("/users", func(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(400, err)
	}
	if err := c.Validate(user); err != nil {
		return c.JSON(400, err)
	}
	return c.JSON(201, user)
})
```

### Pattern 2: Error Handling

```go
// ❌ MAL: Sin manejo de errores
e.GET("/users/:id", func(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(200, getUser(id))
})

// ✅ BIEN: Con manejo de errores
e.GET("/users/:id", func(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid ID"})
	}
	user, err := getUser(id)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, user)
})
```

### Pattern 3: Middleware

```go
// ❌ MAL: Lógica mezclada en handlers
e.GET("/protected", func(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(401, "Unauthorized")
	}
	// ... verificar token
	// ... lógica del handler
	return c.JSON(200, "OK")
})

// ✅ BIEN: Middleware separado
e.GET("/protected", protectedHandler,
	middleware.AuthMiddleware,
)

func protectedHandler(c echo.Context) error {
	return c.JSON(200, "OK")
}
```

### Pattern 4: Context Management

```go
// ❌ MAL: Guardar todo en context
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("user", user)
		c.Set("roles", roles)
		c.Set("permissions", perms)
		return next(c)
	}
})

// ✅ BIEN: Usar struct para datos relacionados
type UserContext struct {
	User        *User
	Roles       []string
	Permissions []string
}

e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("user_ctx", UserContext{
			User: user,
			Roles: roles,
			Permissions: perms,
		})
		return next(c)
	}
})
```

---

## CONCLUSIÓN

Echo es un framework enterprise-grade que combina:
- **Performance**: 120k+ RPS
- **Simplicidad**: API coherente y bien diseñada
- **Extensibilidad**: Middleware pattern flexible
- **Madureza**: Usado en producción por grandes empresas

Es la elección ideal para:
- APIs REST de alto rendimiento
- Microservicios internos
- Aplicaciones real-time con WebSocket
- Proyectos que requieren TLS/HTTP2 nativo

**Recursos adicionales:**
- Documentación: https://echo.labstack.com
- GitHub: https://github.com/labstack/echo
- Community: Discord, GitHub Discussions


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/58-echo-framework/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/58-echo-framework):

```bash
cd examples/58-echo-framework
go run .
```
