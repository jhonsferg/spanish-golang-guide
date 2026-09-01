# Capítulo 60: Chi router - Ruteo ligero y componible

## Índice Completo

1. [60.1 - Introducción a Chi](#601---introducción-a-chi)
2. [60.2 - Fundamentales de Chi](#602---fundamentales-de-chi)
3. [60.3 - Enrutamiento Avanzado](#603---enrutamiento-avanzado)
4. [60.4 - Sistema de Middleware](#604---sistema-de-middleware)
5. [60.5 - Request & Response](#605---request--response)
6. [60.6 - Composabilidad & Patrones](#606---composabilidad--patrones)
7. [60.7 - Validación & Manejo de Errores](#607---validación--manejo-de-errores)
8. [60.8 - Testing de Aplicaciones Chi](#608---testing-de-aplicaciones-chi)
9. [60.9 - Patrones Comunes](#609---patrones-comunes)
10. [60.10 - Performance & Comparación](#6010---performance--comparación)
11. [60.11 - Producción & Casos de Estudio](#6011---producción--casos-de-estudio)

---

## 60.1 - Introducción a Chi

### 60.1.1 ¿Qué es Chi?

**Chi** es un router HTTP minimalista y composable para Go, diseñado con la filosofía de: **"Small core, big ecosystem"**. A diferencia de frameworks como Gin o Echo que intentan proporcionar una solución todo-en-uno, Chi se enfoca en hacer bien UNA cosa: enrutar peticiones HTTP.

Chi fue creado por **Peter Bourgon** y el equipo de Gokit, y desde 2015 se ha convertido en el router preferido para aplicaciones Go que buscan composabilidad y flexibilidad.

**Características principales:**
- ✅ Router ultra-minimalista (~1000 líneas de código)
- ✅ Basado en Radix Tree para enrutamiento O(1)
- ✅ Middleware composable y modular
- ✅ Montaje de routers (router mounting)
- ✅ Compatible con `net/http` estándar
- ✅ Excelente rendimiento en producción
- ✅ Comunidad activa y ecosystem rico

### 60.1.2 Filosofía de Chi

```
┌──────────────────────────────────────────────────────┐
│           FILOSOFÍA DE CHI                           │
├──────────────────────────────────────────────────────┤
│ 1. Minimalismo                                       │
│    - Solo lo esencial para enrutamiento              │
│    - ~1000 líneas de código core                     │
│    - Cero dependencias externas                      │
│                                                      │
│ 2. Composabilidad                                    │
│    - Combinar routers es natural                     │
│    - Middleware componible                          │
│    - Handlers simples y puros                       │
│                                                      │
│ 3. Compatibilidad                                   │
│    - net/http 100% compatible                       │
│    - Funciona con cualquier middleware              │
│    - Fácil de integrar                              │
│                                                      │
│ 4. Performance                                      │
│    - Radix tree (O(1) routing)                      │
│    - Cero allocations en path hot-path              │
│    - Benchmarks: ~2x más rápido que algunos         │
│                                                      │
│ 5. Claridad                                         │
│    - Código legible y mantenible                    │
│    - Error handling explícito                       │
│    - Patterns claros                                │
└──────────────────────────────────────────────────────┘
```

### 60.1.3 ¿Cuándo usar Chi?

**✅ Usa Chi cuando:**

1. **Necesitas composabilidad**
   - Arquitecturas de microservicios
   - APIs modulares y escalables
   - Múltiples equipos trabajando

2. **Requieres control total**
   - No quieres "magia" de frameworks
   - Necesitas middleware personalizado
   - Decisiones explícitas en tu código

3. **Construyes APIs REST limpias**
   - HTTP estándar como está diseñado
   - Recursos claros y RESTful
   - JSON APIs simples

4. **Rendimiento es crítico**
   - Microservicios de bajo latency
   - APIs de alto volumen
   - Minimal overhead requerido

5. **Equipo prefiere Go idiomático**
   - Uso de `net/http` estándar
   - Interfaces simples
   - Código legible

**❌ Considera alternativas cuando:**

1. **Necesitas auto-generación de documentación**
   - OpenAPI/Swagger automático → Echo, Fiber
   - ORM integrado → Fiber
   - Validación automática de tipos

2. **Requieres "batteries included"**
   - HTML templates built-in → Fiber
   - Sesiones → Gin
   - ORM integrado

3. **Tu equipo prefiere patterns declarativos**
   - Tags y annotations → Gin, Fiber
   - Decoradores → Echo

### 60.1.4 Ecosistema Chi

Chi no incluye todo, pero tiene un ecosystem rico:

```go
// Middleware oficial y community
chi/middleware        // Logger, Recovery, Compressor, etc.
chi-render            // JSON/XML rendering
go-chi/chi-httprate   // Rate limiting
go-chi/cors           // CORS middleware
go-chi/prometheus     // Prometheus metrics
go-chi/jwtauth        // JWT authentication
go-chi/oauth          // OAuth middleware
```

### 60.1.5 Adopción en la Comunidad

Chi es usado en producción por:
- **Uber** (componentes internos)
- **Hashicorp** (APIs internas)
- **Countless startups** en la comunidad Go
- **Open source projects** populares

```
Estadísticas aproximadas (2024):
┌─────────────────────────────────┐
│ GitHub Stars: ~17,000+          │
│ Go package downloads: Millones  │
│ Versión: v5.x (estable)         │
│ Mantenimiento: Activo            │
└─────────────────────────────────┘
```

---

## 60.2 - Fundamentales de Chi

### 60.2.1 Instalación

```bash
# Instalación estándar
go get github.com/go-chi/chi/v5

# Verificar instalación
go mod tidy
```

### 60.2.2 Creación Básica de Router

```go
package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Crear router
	r := chi.NewRouter()
	
	// Definir rutas
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hola, Chi!"))
	})
	
	// Servidor
	http.ListenAndServe(":3000", r)
}
```

### 60.2.3 Estructura de Handlers

Chi usa la firma estándar de Go:

```go
// Firma de handler - igual a net/http
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Función helper
type HandlerFunc func(w http.ResponseWriter, r *http.Request)
```

**Ejemplo completo:**

```go
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// Procesar petición
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Mundo"
	}
	
	// Responder
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Hola, ` + name + `!"}`))
}

func main() {
	r := chi.NewRouter()
	r.Get("/hello", HelloHandler)
	http.ListenAndServe(":3000", r)
}
```

### 60.2.4 Métodos HTTP

Chi proporciona métodos convenientes para todos los métodos HTTP:

```go
r := chi.NewRouter()

// Métodos comunes
r.Get("/users", ListUsers)          // GET
r.Post("/users", CreateUser)        // POST
r.Put("/users/{id}", UpdateUser)    // PUT
r.Patch("/users/{id}", PartialUpdate) // PATCH
r.Delete("/users/{id}", DeleteUser) // DELETE
r.Head("/status", HealthCheck)      // HEAD
r.Options("/", Options)             // OPTIONS

// Manejo de múltiples métodos
r.MethodFunc("GET", "/data", GetData)
r.MethodFunc("POST", "/data", PostData)

// Cualquier método (catch-all)
r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Ruta no encontrada"))
})

// Método no permitido
r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Método no permitido"))
})
```

### 60.2.5 Listen and Serve

```go
// Básico
http.ListenAndServe(":3000", r)

// Con control de timeout
server := &http.Server{
	Addr:         ":3000",
	Handler:      r,
	ReadTimeout:  15 * time.Second,
	WriteTimeout: 15 * time.Second,
	IdleTimeout:  60 * time.Second,
}

if err := server.ListenAndServe(); err != nil {
	log.Fatal(err)
}

// Con graceful shutdown
go func() {
	if err := server.ListenAndServe(); err != nil && 
	   err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
}()

// Listening...
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
<-sigChan

ctx, cancel := context.WithTimeout(context.Background(), 
	5*time.Second)
defer cancel()
server.Shutdown(ctx)
```

### 60.2.6 Ejemplo Completo: Servidor Simple

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	
	// Middleware global
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// Rutas
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bienvenido a Chi",
			"time":    time.Now().String(),
		})
	})
	
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})
	
	log.Println("🚀 Servidor escuchando en :3000")
	http.ListenAndServe(":3000", r)
}
```

---

## 60.3 - Enrutamiento Avanzado

### 60.3.1 Parámetros de Ruta

Chi soporta parámetros dinámicos en rutas:

```go
r := chi.NewRouter()

// Parámetro simple
r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Usuario: " + id))
})

// Múltiples parámetros
r.Get("/posts/{postID}/comments/{commentID}", func(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	commentID := chi.URLParam(r, "commentID")
	w.Write([]byte("Post: " + postID + ", Comment: " + commentID))
})

// Parámetro con ruta
r.Get("/files/{filepath}", func(w http.ResponseWriter, r *http.Request) {
	filepath := chi.URLParam(r, "filepath")
	w.Write([]byte("Archivo: " + filepath))
})

// Números
r.Get("/api/v{version}/users", func(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	w.Write([]byte("API v" + version))
})
```

### 60.3.2 Rutas Wildcard

```go
// Wildcard (captura todo después)
r.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
	// r.RequestURI conserva el wildcard
	w.Write([]byte("Ruta: " + r.RequestURI))
})

// Uso práctico: Servir archivos estáticos
r.Get("/static/*", http.StripPrefix("/static", 
	http.FileServer(http.Dir("./static"))).ServeHTTP)

// Servir documentación
r.Get("/docs/*", http.StripPrefix("/docs",
	http.FileServer(http.Dir("./docs"))).ServeHTTP)
```

### 60.3.3 Patrones Regex

```go
// Aunque Chi no tiene regex built-in, puedes hacerlo
// manualmente para rutas específicas

// Opción 1: Validar dentro del handler
r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	// Validar que sea un número
	if !isNumeric(id) {
		http.Error(w, "ID debe ser numérico", http.StatusBadRequest)
		return
	}
	
	w.Write([]byte("Usuario válido: " + id))
})

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// Opción 2: Router separado para parámetros específicos
r.Route("/api/v{version}", func(r chi.Router) {
	// Solo números
	if version := chi.URLParam(r, "version"); isNumeric(version) {
		r.Get("/users", ListUsers)
	}
})
```

### 60.3.4 Query Parameters

```go
r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
	// Un parámetro
	query := r.URL.Query().Get("q")
	
	// Parámetro con valor por defecto
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	
	// Múltiples valores
	tags := r.URL.Query()["tag"]
	
	// Acceso a todos
	params := r.URL.Query()
	
	response := map[string]interface{}{
		"query": query,
		"page":  page,
		"tags":  tags,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
})

// Uso: /search?q=golang&page=2&tag=web&tag=api
```

### 60.3.5 Route Grouping (Group)

```go
r := chi.NewRouter()

// Agrupar rutas por prefijo
r.Route("/api", func(r chi.Router) {
	// Middleware solo para /api/*
	r.Use(AuthMiddleware)
	
	r.Route("/v1", func(r chi.Router) {
		r.Get("/users", ListUsers)
		r.Post("/users", CreateUser)
		r.Get("/users/{id}", GetUser)
		r.Delete("/users/{id}", DeleteUser)
	})
	
	r.Route("/v2", func(r chi.Router) {
		r.Get("/users", ListUsersV2)
		r.Post("/users", CreateUserV2)
	})
})

// Patrón de group adicional
r.Route("/admin", func(r chi.Router) {
	r.Use(AdminOnly)
	
	r.Get("/dashboard", AdminDashboard)
	r.Get("/stats", AdminStats)
})
```

**Resultado:**
```
GET    /api/v1/users         → ListUsers
POST   /api/v1/users         → CreateUser
GET    /api/v1/users/{id}    → GetUser
DELETE /api/v1/users/{id}    → DeleteUser
GET    /api/v2/users         → ListUsersV2
POST   /api/v2/users         → CreateUserV2
GET    /admin/dashboard      → AdminDashboard (requiere AdminOnly)
GET    /admin/stats          → AdminStats (requiere AdminOnly)
```

### 60.3.6 Route Mounting (Mount)

**El patrón más poderoso de Chi** - permite composición de routers:

```go
// Definir subrouter independiente
func createUsersRouter() chi.Router {
	r := chi.NewRouter()
	
	r.Get("/", ListUsers)
	r.Post("/", CreateUser)
	r.Get("/{id}", GetUser)
	r.Put("/{id}", UpdateUser)
	r.Delete("/{id}", DeleteUser)
	
	return r
}

// Router de posts
func createPostsRouter() chi.Router {
	r := chi.NewRouter()
	
	r.Get("/", ListPosts)
	r.Post("/", CreatePost)
	r.Get("/{id}", GetPost)
	r.Put("/{id}", UpdatePost)
	
	return r
}

// Router principal
func main() {
	r := chi.NewRouter()
	
	// Middleware global
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// Montar subrouters
	r.Mount("/api/users", createUsersRouter())
	r.Mount("/api/posts", createPostsRouter())
	
	// Pueden tener middleware específico
	r.Route("/admin", func(r chi.Router) {
		r.Use(AdminAuthMiddleware)
		r.Mount("/analytics", createAnalyticsRouter())
	})
	
	http.ListenAndServe(":3000", r)
}
```

**Ventajas del Mounting:**
- ✅ Separación de concerns
- ✅ Reutilización de subrouters
- ✅ Escalabilidad
- ✅ Testing independiente

### 60.3.7 Diagrama de Enrutamiento

```
┌─────────────────────────────────────────────────────────┐
│                   REQUEST FLOW                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. HTTP Request                                        │
│     ↓                                                   │
│  2. Chi Router                                          │
│     ├─ Match method (GET, POST, etc.)                  │
│     ├─ Match path (radix tree O(1))                    │
│     └─ Extract parameters                              │
│     ↓                                                   │
│  3. Middleware Stack (LIFO - Last In First Out)        │
│     ├─ middleware1                                     │
│     ├─ middleware2                                     │
│     └─ middleware3                                     │
│     ↓                                                   │
│  4. Handler                                            │
│     ├─ Process request                                 │
│     ├─ Access parameters                               │
│     └─ Write response                                  │
│     ↑                                                   │
│  5. Middleware unwinding (reverse order)               │
│     ├─ middleware3                                     │
│     ├─ middleware2                                     │
│     └─ middleware1                                     │
│     ↓                                                   │
│  6. HTTP Response                                      │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 60.4 - Sistema de Middleware

### 60.4.1 Firma y Concepto

```go
// Middleware es una función que envuelve handlers
type Middleware func(next http.Handler) http.Handler

// O usando chi
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Patrón general
func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Antes del handler
		fmt.Println("Antes del handler")
		
		// Llamar al siguiente
		next.ServeHTTP(w, r)
		
		// Después del handler
		fmt.Println("Después del handler")
	})
}
```

### 60.4.2 Middleware Integrado de Chi

```go
package middleware

// Logger - registra todas las peticiones
r.Use(middleware.Logger)

// Recoverer - recupera de panics
r.Use(middleware.Recoverer)

// Compressor - comprime respuestas (gzip, deflate)
r.Use(middleware.Compress(5))

// Timeout - establece timeout para peticiones
r.Use(middleware.Timeout(30 * time.Second))

// RequestID - añade ID único a cada petición
r.Use(middleware.RequestID)

// RealIP - obtiene la IP real del cliente
r.Use(middleware.RealIP)

// Heartbeat - endpoint healthcheck
r.Use(middleware.Heartbeat("/ping"))

// Strip slashes - normaliza rutas sin trailing slash
r.Use(middleware.StripSlashes)
```

### 60.4.3 Middleware Global y Local

```go
r := chi.NewRouter()

// Middleware GLOBAL - se aplica a todas las rutas
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)

// Middleware LOCAL - solo para algunas rutas
r.Route("/api", func(r chi.Router) {
	r.Use(AuthMiddleware)
	r.Use(RateLimitMiddleware)
	
	r.Get("/users", GetUsers)
})

r.Route("/public", func(r chi.Router) {
	// Sin Auth ni RateLimit
	r.Get("/posts", GetPosts)
})
```

### 60.4.4 Crear Middleware Personalizado

**Ejemplo 1: Logging personalizado**

```go
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Registrar
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		// Ejecutar handler
		next.ServeHTTP(w, r)
		
		// Registrar tiempo
		duration := time.Since(start)
		log.Printf("Completado en %v", duration)
	})
}

// Usar
r.Use(LoggingMiddleware)
```

**Ejemplo 2: Autenticación**

```go
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		
		if token == "" {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}
		
		if !isValidToken(token) {
			http.Error(w, "Token inválido", http.StatusForbidden)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func isValidToken(token string) bool {
	// Validación real aquí
	return strings.HasPrefix(token, "Bearer ")
}
```

**Ejemplo 3: CORS**

```go
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", 
			"GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", 
			"Content-Type, Authorization")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
```

**Ejemplo 4: Rate Limiting**

```go
func RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := make(map[string]*time.Ticker)
	mu := &sync.Mutex{}
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		
		mu.Lock()
		if _, exists := limiter[ip]; !exists {
			limiter[ip] = time.NewTicker(1 * time.Second)
		}
		ticker := limiter[ip]
		mu.Unlock()
		
		select {
		case <-ticker.C:
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "Demasiadas peticiones", 
				http.StatusTooManyRequests)
		}
	})
}
```

### 60.4.5 Using() vs With()

```go
// Using() - Middleware para todas las rutas del router
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)

// With() - Middleware solo para una ruta específica
r.With(middleware.Timeout(5*time.Second)).
	Get("/slow-endpoint", SlowHandler)

// Ejemplo más complejo
r.With(
	AuthMiddleware,
	RateLimitMiddleware,
	LoggingMiddleware,
).Post("/api/create", CreateHandler)

// Chain - crear cadena de middleware
chain := middleware.Chain(
	middleware.Logger,
	middleware.Recoverer,
	CustomMiddleware,
)

r.Use(chain)
```

### 60.4.6 Cadenas de Middleware

```go
// Crear helper para cadenas
type Chain struct {
	middlewares []func(http.Handler) http.Handler
}

func (c *Chain) Add(mw func(http.Handler) http.Handler) *Chain {
	c.middlewares = append(c.middlewares, mw)
	return c
}

func (c *Chain) Handler(h http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}

// Uso
chain := &Chain{}
chain.Add(middleware.Logger).
	Add(middleware.Recoverer).
	Add(AuthMiddleware)

r.Handle("/api/users", chain.Handler(GetUsers))
```

### 60.4.7 Middleware con Contexto

```go
// Pasar datos entre middleware
func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)
		
		// Guardar en contexto
		ctx := context.WithValue(r.Context(), "userID", userID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Acceder en handler
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	w.Write([]byte("Perfil de usuario: " + userID))
}
```

---

## 60.5 - Request & Response

### 60.5.1 Acceso a Datos de Request

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Método HTTP
	method := r.Method // GET, POST, etc.
	
	// URL y rutas
	path := r.URL.Path           // /users/123
	url := r.URL.String()        // URL completa
	rawQuery := r.URL.RawQuery   // q=search&page=1
	
	// Headers
	userAgent := r.Header.Get("User-Agent")
	contentType := r.Header.Get("Content-Type")
	authorization := r.Header.Get("Authorization")
	
	// Host
	host := r.Host
	scheme := r.Header.Get("X-Forwarded-Proto")
	
	// Cliente
	remoteAddr := r.RemoteAddr
}
```

### 60.5.2 Query Parameters

```go
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Un valor
	query := r.URL.Query().Get("q")
	
	// Con valor por defecto
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	
	// Múltiples valores del mismo parámetro
	tags := r.URL.Query()["tag"]  // []string{"web", "api"}
	
	// Acceso directo
	q := r.URL.Query()
	for key, values := range q {
		for _, value := range values {
			log.Printf("%s = %s", key, value)
		}
	}
}
```

### 60.5.3 Path Parameters

```go
func UserDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Con Chi
	userID := chi.URLParam(r, "id")
	
	// Múltiples parámetros
	postID := chi.URLParam(r, "postID")
	commentID := chi.URLParam(r, "commentID")
}

// En main
r.Get("/users/{id}", UserDetailHandler)
r.Get("/posts/{postID}/comments/{commentID}", CommentHandler)
```

### 60.5.4 Body Parsing

Chi no proporciona automático binding, así que es manual (como debe ser):

```go
// JSON Request
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	
	// Decodificar
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	
	// Validar
	if req.Name == "" || req.Email == "" {
		http.Error(w, "Campos requeridos", http.StatusBadRequest)
		return
	}
	
	// Procesar
	user := User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	
	// Responder
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
```

**Form Data:**

```go
func FormHandler(w http.ResponseWriter, r *http.Request) {
	// Parsear form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parseando form", http.StatusBadRequest)
		return
	}
	
	// Acceder campos
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	// Archivos
	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, "Error en archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	log.Printf("Archivo: %s, Tamaño: %d bytes", 
		handler.Filename, handler.Size)
}
```

**Raw Body:**

```go
func RawBodyHandler(w http.ResponseWriter, r *http.Request) {
	// Leer cuerpo completo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo body", http.StatusBadRequest)
		return
	}
	
	// Usar body
	log.Printf("Body: %s", string(body))
}
```

### 60.5.5 Response Writing

```go
func ResponseExamples(w http.ResponseWriter, r *http.Request) {
	// Configurar headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Custom-Header", "value")
	
	// Status code (debe ser antes de Write)
	w.WriteHeader(http.StatusOK)
	
	// Escribir contenido
	w.Write([]byte("Hola"))
	
	// O usar fmt.Fprintf
	fmt.Fprintf(w, "Número: %d", 42)
	
	// JSON
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// Responder diferentes tipos
func JSONResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"status": "success",
		"data": map[string]string{
			"id":   "123",
			"name": "John",
		},
	}
	json.NewEncoder(w).Encode(data)
}

func XMLResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(struct {
		Status string `xml:"status"`
	}{Status: "ok"})
}

func PlainTextResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "Texto simple")
}

func HTMLResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<html><body><h1>Hola</h1></body></html>`
	w.Write([]byte(html))
}
```

### 60.5.6 Streaming Response

```go
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Streaming JSON
	w.Write([]byte("["))
	
	for i := 0; i < 100; i++ {
		if i > 0 {
			w.Write([]byte(","))
		}
		
		json.NewEncoder(w).Encode(map[string]int{
			"item": i,
		})
		
		// Flush después de cada item (si es necesario)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
	
	w.Write([]byte("]"))
}

// Server-Sent Events (SSE)
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	
	flusher := w.(http.Flusher)
	
	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "data: Mensaje %d\n\n", i)
		flusher.Flush()
		time.Sleep(1 * time.Second)
	}
}
```

---

## 60.6 - Composabilidad & Patrones

### 60.6.1 El Poder del Mounting

La verdadera fuerza de Chi es la composabilidad mediante mounting:

```go
// Cada module es independiente

// users/router.go
func NewUserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", ListUsers)
	r.Post("/", CreateUser)
	r.Get("/{id}", GetUser)
	r.Put("/{id}", UpdateUser)
	r.Delete("/{id}", DeleteUser)
	return r
}

// posts/router.go
func NewPostRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", ListPosts)
	r.Post("/", CreatePost)
	r.Get("/{id}", GetPost)
	return r
}

// comments/router.go
func NewCommentRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", ListComments)
	r.Post("/", CreateComment)
	return r
}

// main.go
func main() {
	r := chi.NewRouter()
	
	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// Mount
	r.Mount("/api/users", NewUserRouter())
	r.Mount("/api/posts", NewPostRouter())
	r.Mount("/api/comments", NewCommentRouter())
	
	http.ListenAndServe(":3000", r)
}
```

### 60.6.2 Inyección de Dependencias

```go
// Con interfaces
type UserStore interface {
	GetUser(id string) (User, error)
	SaveUser(user User) error
}

type UserHandler struct {
	store UserStore
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.store.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// Crear router
func NewUserRouter(store UserStore) chi.Router {
	h := &UserHandler{store: store}
	r := chi.NewRouter()
	r.Get("/{id}", h.GetUser)
	return r
}

// main.go
func main() {
	store := NewMySQLUserStore()
	r := chi.NewRouter()
	r.Mount("/users", NewUserRouter(store))
	http.ListenAndServe(":3000", r)
}
```

### 60.6.3 Composición de Middleware

```go
// Crear conjuntos de middleware reutilizables

// API middleware
func apiMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.Logger,
		middleware.Recoverer,
		AuthMiddleware,
		RateLimitMiddleware,
	}
}

// Admin middleware
func adminMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.Logger,
		middleware.Recoverer,
		StrictAuthMiddleware,
		AdminOnlyMiddleware,
	}
}

// Aplicar
func main() {
	r := chi.NewRouter()
	
	r.Route("/api", func(r chi.Router) {
		for _, mw := range apiMiddleware() {
			r.Use(mw)
		}
		r.Mount("/users", NewUserRouter())
	})
	
	r.Route("/admin", func(r chi.Router) {
		for _, mw := range adminMiddleware() {
			r.Use(mw)
		}
		r.Mount("/dashboard", NewAdminRouter())
	})
}
```

### 60.6.4 Handler Wrapping

```go
// Envolver handlers para agregar funcionalidad

func withTimer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("Handler tardó: %v", time.Since(start))
	})
}

func withRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal Server Error", 
					http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// Usar
r.Get("/fast", withTimer(withRecovery(
	http.HandlerFunc(FastHandler))))
```

### 60.6.5 Diagrama de Composición

```
┌──────────────────────────────────────────────────────────┐
│           ARQUITECTURA COMPOSABLE CON CHI               │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  ┌─── Aplicación Principal                              │
│  │                                                      │
│  ├─ Global Middleware                                  │
│  │  ├─ Logger                                          │
│  │  ├─ Recoverer                                       │
│  │  └─ RequestID                                       │
│  │                                                      │
│  ├─ /api (Route)                                       │
│  │  └─ API Middleware (Auth, RateLimit)                │
│  │     ├─ /users (Mount) → UserRouter                  │
│  │     │  └─ GET /{id}  → GetUser                      │
│  │     ├─ /posts (Mount) → PostRouter                  │
│  │     │  └─ GET /{id}  → GetPost                      │
│  │     └─ /comments (Mount) → CommentRouter            │
│  │        └─ GET /{id}  → GetComment                   │
│  │                                                      │
│  ├─ /admin (Route)                                     │
│  │  └─ Admin Middleware (StrictAuth)                   │
│  │     ├─ /dashboard (Mount) → AdminRouter             │
│  │     └─ /stats (Mount) → StatsRouter                 │
│  │                                                      │
│  └─ /health (Health Check)                            │
│     └─ GET /health → OK                                │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

---

## 60.7 - Validación & Manejo de Errores

### 60.7.1 Estrategias de Validación

```go
// Opción 1: Validación Manual Explícita
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	// Validación explícita
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name requerido")
		return
	}
	
	if !isValidEmail(req.Email) {
		respondError(w, http.StatusBadRequest, "Email inválido")
		return
	}
	
	if req.Age < 0 || req.Age > 150 {
		respondError(w, http.StatusBadRequest, "Age inválido")
		return
	}
	
	// Procesar
	user := createUser(req.Name, req.Email, req.Age)
	respondJSON(w, http.StatusCreated, user)
}

func isValidEmail(email string) bool {
	// Validación real
	return strings.Contains(email, "@")
}

// Opción 2: Usar library (go-playground/validator)
import "github.com/go-playground/validator/v10"

func CreateUserHandlerWithValidator(w http.ResponseWriter, r *http.Request) {
	type CreateUserReq struct {
		Name  string `json:"name" validate:"required,min=2,max=100"`
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age" validate:"required,min=0,max=150"`
	}
	
	var req CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	v := validator.New()
	if err := v.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	user := createUser(req.Name, req.Email, req.Age)
	respondJSON(w, http.StatusCreated, user)
}
```

### 60.7.2 Manejo de Errores

```go
// Opción 1: Errores Simples
type APIError struct {
	Code    int
	Message string
}

func handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	
	switch err.(type) {
	case *NotFoundError:
		respondError(w, http.StatusNotFound, err.Error())
	case *ValidationError:
		respondError(w, http.StatusBadRequest, err.Error())
	case *AuthError:
		respondError(w, http.StatusUnauthorized, err.Error())
	default:
		log.Printf("Unexpected error: %v", err)
		respondError(w, http.StatusInternalServerError, 
			"Internal Server Error")
	}
}

// Opción 2: Errores Estructurados
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func respondError(w http.ResponseWriter, code int, 
	err *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

// Uso
respondError(w, http.StatusBadRequest, &Error{
	Code:    "INVALID_EMAIL",
	Message: "Email format is invalid",
	Details: map[string]interface{}{
		"field": "email",
		"value": "invalid-email",
	},
})
```

### 60.7.3 Custom Error Types

```go
// Definir tipos de error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on %s: %s", 
		e.Field, e.Message)
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("Authentication error: %s", e.Message)
}

// Usar
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	user, err := findUser(id)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			respondError(w, http.StatusNotFound, 
				&Error{Code: "USER_NOT_FOUND", Message: err.Error()})
		}
		return
	}
	
	respondJSON(w, http.StatusOK, user)
}
```

### 60.7.4 Error Wrapping y Context

```go
// Go 1.13+ error wrapping
import "fmt"

func getUser(id string) (User, error) {
	user, err := db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		// Wrap error para contexto
		return User{}, fmt.Errorf("failed to query user %s: %w", id, err)
	}
	return user, nil
}

// Verificar tipo
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := getUser(id)
	
	if err != nil {
		log.Printf("Error: %v", err)
		respondError(w, http.StatusInternalServerError, 
			&Error{Code: "INTERNAL_ERROR", Message: "Failed to get user"})
		return
	}
	
	respondJSON(w, http.StatusOK, user)
}
```

### 60.7.5 Helpers para Response

```go
// Helpers para simplificar respuestas
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func respondSuccess(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func respondCreated(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// Uso simplificado
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := findUser(id)
	
	if err != nil {
		respondError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	
	respondSuccess(w, user)
}
```

---

## 60.8 - Testing de Aplicaciones Chi

### 60.8.1 Testing con httptest

Chi funciona perfecto con `net/http/httptest`:

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi/v5"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	return r
}

func TestGetRoot(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	setupRouter().ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	
	if w.Body.String() != "Hello" {
		t.Errorf("Expected 'Hello', got %s", w.Body.String())
	}
}
```

### 60.8.2 Testing Table-Driven

```go
func TestMultipleRoutes(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedCode   int
		expectedBody   string
	}{
		{
			name:         "GET root",
			method:       "GET",
			path:         "/",
			expectedCode: 200,
			expectedBody: "Welcome",
		},
		{
			name:         "GET user",
			method:       "GET",
			path:         "/users/123",
			expectedCode: 200,
			expectedBody: "User: 123",
		},
		{
			name:         "POST user",
			method:       "POST",
			path:         "/users",
			expectedCode: 201,
			expectedBody: "User created",
		},
		{
			name:         "Not found",
			method:       "GET",
			path:         "/notfound",
			expectedCode: 404,
		},
	}
	
	r := setupRouter()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			
			r.ServeHTTP(w, req)
			
			if w.Code != tt.expectedCode {
				t.Errorf("Expected %d, got %d", 
					tt.expectedCode, w.Code)
			}
			
			if tt.expectedBody != "" && 
			   w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'",
					tt.expectedBody, w.Body.String())
			}
		})
	}
}
```

### 60.8.3 Testing con JSON

```go
import (
	"bytes"
	"encoding/json"
)

func TestCreateUser(t *testing.T) {
	payload := map[string]string{
		"name":  "John",
		"email": "john@example.com",
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	setupRouter().ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
	
	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	
	if result["name"] != "John" {
		t.Error("Name not set correctly")
	}
}
```

### 60.8.4 Testing Middleware

```go
func TestAuthMiddleware(t *testing.T) {
	r := chi.NewRouter()
	r.Use(authMiddleware)
	
	r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	})
	
	// Sin token
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 sin token, got %d", w.Code)
	}
	
	// Con token válido
	req = httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 con token, got %d", w.Code)
	}
}
```

### 60.8.5 Testing de Integración

```go
func TestFullAPI(t *testing.T) {
	// Setup
	r := setupRouter()
	server := httptest.NewServer(r)
	defer server.Close()
	
	// Test 1: Create
	resp, err := http.Post(server.URL+"/users", 
		"application/json", 
		bytes.NewReader([]byte(`{"name":"John"}`)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Create failed: %d", resp.StatusCode)
	}
	
	// Test 2: List
	resp, err = http.Get(server.URL + "/users")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("List failed: %d", resp.StatusCode)
	}
	
	// Test 3: Get
	resp, err = http.Get(server.URL + "/users/1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Get failed: %d", resp.StatusCode)
	}
}
```

---

## 60.9 - Patrones Comunes

### 60.9.1 REST API Design

```go
package main

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = map[string]User{
	"1": {ID: "1", Name: "John", Email: "john@example.com"},
	"2": {ID: "2", Name: "Jane", Email: "jane@example.com"},
}

// Handlers
func ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if user.Name == "" || user.Email == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	
	user.ID = generateID()
	users[user.ID] = user
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	user.ID = id
	users[id] = user
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

// Router
func NewUserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", ListUsers)
	r.Post("/", CreateUser)
	r.Get("/{id}", GetUser)
	r.Put("/{id}", UpdateUser)
	r.Delete("/{id}", DeleteUser)
	return r
}

// Helpers
func generateID() string {
	return "user_" + randString(8)
}

func randString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.Mount("/api/users", NewUserRouter())
	
	log.Println("Server running on :3000")
	http.ListenAndServe(":3000", r)
}
```

### 60.9.2 Middleware de Autenticación

```go
// JWT Authentication Middleware
import "github.com/golang-jwt/jwt/v5"

const secretKey = "your-secret-key"

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}
		
		// Parsear token
		token, err := jwt.ParseWithClaims(authHeader, 
			&Claims{}, 
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})
		
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		
		// Pasar a contexto
		claims := token.Claims.(*Claims)
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

### 60.9.3 CORS Middleware

```go
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", 
			"GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", 
			"Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// O usar library
import "github.com/go-chi/cors"

r.Use(cors.Handler(cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	AllowCredentials: true,
	MaxAge:           300,
}))
```

### 60.9.4 Logging Middleware Avanzado

```go
type LogEntry struct {
	Method     string
	Path       string
	Status     int
	Duration   time.Duration
	RemoteAddr string
	UserAgent  string
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrapper para capturar status
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
		
		next.ServeHTTP(wrapped, r)
		
		entry := LogEntry{
			Method:     r.Method,
			Path:       r.URL.Path,
			Status:     wrapped.statusCode,
			Duration:   time.Since(start),
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.Header.Get("User-Agent"),
		}
		
		// Log
		log.Printf("[%s] %s %s -> %d (%v)",
			entry.Method, entry.Path, entry.RemoteAddr,
			entry.Status, entry.Duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}
```

### 60.9.5 Request ID Tracking

```go
import "github.com/google/uuid"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generar o usar existente
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Añadir a contexto
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		
		// Responder con el ID
		w.Header().Set("X-Request-ID", requestID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Usar en handlers
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)
	log.Printf("[%s] Processing request", requestID)
	w.Write([]byte("OK"))
}
```

---

## 60.10 - Performance & Comparación

### 60.10.1 Benchmarks de Rendimiento

```
┌─────────────────────────────────────────────────────────┐
│         BENCHMARKS DE ROUTERS GO (2024)                 │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ Operaciones por segundo (requests/sec):                │
│                                                         │
│ Chi v5              ████████████ 250,000 ops/sec       │
│ Echo v4             ████████████ 245,000 ops/sec       │
│ Fiber v3            ████████████ 240,000 ops/sec       │
│ Gin v1              ███████████  200,000 ops/sec       │
│ Stdlib http.Mux     ████████     120,000 ops/sec       │
│                                                         │
│ Memoria por instancia (MB):                            │
│                                                         │
│ Chi                 ██ 2.5 MB                          │
│ Echo                ███ 3.1 MB                         │
│ Fiber               ████ 4.2 MB                        │
│ Gin                 ███ 3.5 MB                         │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 60.10.2 Comparación de Frameworks

```go
┌──────────────────────────────────────────────────────────────┐
│         COMPARACIÓN: CHI vs GIN vs ECHO vs FIBER             │
├──────────────────────────────────────────────────────────────┤
│ CARACTERÍSTICA      │ CHI    │ GIN   │ ECHO  │ FIBER        │
├──────────────────────────────────────────────────────────────┤
│ Minimalismo         │ ✅✅✅  │ ⭐⭐  │ ⭐⭐  │ ⭐           │
│ Tamaño código       │ ✅✅✅  │ ⭐   │ ⭐⭐  │ ⭐           │
│ Composabilidad      │ ✅✅✅  │ ⭐   │ ⭐   │ ⭐           │
│ net/http compat.    │ ✅✅✅  │ ⭐   │ ✅✅  │ ⭐           │
│ Performance         │ ✅✅✅  │ ✅✅✅ │ ✅✅✅ │ ✅✅✅        │
│ Built-in features   │ ⭐    │ ✅✅✅ │ ✅✅  │ ✅✅         │
│ Learning curve      │ ✅✅✅  │ ✅✅  │ ✅✅  │ ⭐           │
│ Production ready    │ ✅✅✅  │ ✅✅✅ │ ✅✅✅ │ ✅✅         │
│ Comunidad           │ ✅✅   │ ✅✅✅ │ ✅✅  │ ⭐⭐         │
│ Middleware eco.     │ ✅✅   │ ✅✅✅ │ ✅✅✅ │ ✅✅         │
└──────────────────────────────────────────────────────────────┘

✅✅✅ = Excelente
✅✅   = Muy bien
✅     = Bien
⭐⭐⭐  = Débil
⭐⭐   = Aceptable
⭐    = Pobre
```

### 60.10.3 Código Comparativo

```go
// ============ CHI ============
func setupChi() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte("User: " + id))
	})
	
	return r
}

// ============ GIN ============
func setupGin() *gin.Engine {
	r := gin.Default()
	
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(200, "User: "+id)
	})
	
	return r
}

// ============ ECHO ============
func setupEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	
	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.String(200, "User: "+id)
	})
	
	return e
}

// ============ FIBER ============
func setupFiber() *fiber.App {
	app := fiber.New()
	
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.SendString("User: " + id)
	})
	
	return app
}
```

### 60.10.4 Cuándo Usar Chi vs Otros

```
┌─────────────────────────────────────────────────────────┐
│              MATRIZ DE DECISIÓN                         │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ USA CHI SI:                                             │
│ ✅ Necesitas máxima composabilidad                     │
│ ✅ Construyes microservicios                           │
│ ✅ Quieres máximo control                              │
│ ✅ Rendimiento es crítico                              │
│ ✅ Equipo prefiere Go idiomático                       │
│                                                         │
│ USA GIN SI:                                             │
│ ✅ Necesitas desarrollo rápido                         │
│ ✅ Validación automática de tipos                      │
│ ✅ ORM integrado requerido                             │
│ ✅ Comunidad y ejemplos importantes                    │
│                                                         │
│ USA ECHO SI:                                            │
│ ✅ Necesitas balance entre poder y simplicidad        │
│ ✅ Generación de OpenAPI requerida                    │
│ ✅ Middleware más avanzado built-in                   │
│                                                         │
│ USA FIBER SI:                                           │
│ ✅ Vienes de Node.js/Express                           │
│ ✅ Necesitas máximo rendimiento absoluto              │
│ ✅ Ecosystem Fiber es suficiente                      │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 60.11 - Producción & Casos de Estudio

### 60.11.1 Estructura Producción

```
proyecto/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── users.go
│   │   ├── posts.go
│   │   └── health.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logging.go
│   ├── models/
│   │   └── user.go
│   ├── storage/
│   │   └── database.go
│   └── router/
│       └── router.go
├── config/
│   └── config.go
├── docker/
│   └── Dockerfile
├── tests/
│   └── handlers_test.go
├── go.mod
└── go.sum
```

**main.go:**

```go
package main

import (
	"log"
	"net/http"
	"time"
	
	"myapi/internal/router"
	"myapi/config"
)

func main() {
	cfg := config.Load()
	
	r := router.New(cfg)
	
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && 
	   err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
```

**internal/router/router.go:**

```go
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"myapi/config"
	"myapi/internal/handlers"
	"myapi/internal/middleware/auth"
	"myapi/internal/middleware/cors"
)

func New(cfg *config.Config) chi.Router {
	r := chi.NewRouter()
	
	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler())
	
	// Health check
	r.Get("/health", handlers.Health)
	
	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Mount("/users", handlers.UserRouter())
		r.Mount("/posts", handlers.PostRouter())
	})
	
	// Admin routes (requieren auth)
	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.Middleware())
		r.Mount("/dashboard", handlers.AdminRouter())
	})
	
	return r
}
```

### 60.11.2 Docker Deployment

**Dockerfile:**

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -ldflags="-s -w" -o api ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/config ./config

EXPOSE 3000

CMD ["./api"]
```

**Build y run:**

```bash
# Build imagen
docker build -t myapi:latest .

# Run contenedor
docker run -p 3000:3000 -e PORT=3000 myapi:latest

# Con docker-compose
docker-compose up -d
```

### 60.11.3 Caso de Estudio 1: Microservicio Simple

```go
// Simple microservice for user management

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	users = make(map[string]User)
	mu    sync.RWMutex
)

func listUsers(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	mu.Lock()
	defer mu.Unlock()
	
	if user.ID == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	
	users[user.ID] = user
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	mu.RLock()
	user, ok := users[id]
	mu.RUnlock()
	
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.Get("/users", listUsers)
	r.Post("/users", createUser)
	r.Get("/users/{id}", getUser)
	
	log.Println("🚀 Microservice running on :3000")
	http.ListenAndServe(":3000", r)
}
```

### 60.11.4 Caso de Estudio 2: API Gateway

```go
// Simple API Gateway aggregating multiple services

package main

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func createReverseProxy(target string) http.HandlerFunc {
	proxyURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// Gateway a diferentes servicios
	r.Route("/api", func(r chi.Router) {
		// Users service
		r.Mount("/users", chi.HandlerFunc(
			createReverseProxy("http://localhost:3001")))
		
		// Posts service
		r.Mount("/posts", chi.HandlerFunc(
			createReverseProxy("http://localhost:3002")))
		
		// Comments service
		r.Mount("/comments", chi.HandlerFunc(
			createReverseProxy("http://localhost:3003")))
	})
	
	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gateway OK"))
	})
	
	http.ListenAndServe(":3000", r)
}
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: REST API Simple ⭐

**Objetivo:** Crear una REST API completa con Chi

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// Implementar:
// 1. GET /todos - Listar todos
// 2. POST /todos - Crear todo
// 3. GET /todos/{id} - Obtener todo
// 4. PUT /todos/{id} - Actualizar todo
// 5. DELETE /todos/{id} - Eliminar todo
// 6. Usar middleware Logger y Recoverer

func main() {
	// Tu código aquí
}

// Solución esperada: ~50 líneas
```

**Solución:**

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	todos = make(map[string]Todo)
	mu    sync.RWMutex
	idGen int
)

func ListTodos(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if req.Title == "" {
		http.Error(w, "Title required", http.StatusBadRequest)
		return
	}
	
	mu.Lock()
	defer mu.Unlock()
	
	idGen++
	id := fmt.Sprintf("todo_%d", idGen)
	
	todo := Todo{
		ID:        id,
		Title:     req.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	
	todos[id] = todo
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	mu.RLock()
	todo, ok := todos[id]
	mu.RUnlock()
	
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	mu.Lock()
	defer mu.Unlock()
	
	todo, ok := todos[id]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	
	if req.Title != "" {
		todo.Title = req.Title
	}
	todo.Completed = req.Completed
	
	todos[id] = todo
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	mu.Lock()
	defer mu.Unlock()
	
	if _, ok := todos[id]; !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	
	delete(todos, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.Get("/todos", ListTodos)
	r.Post("/todos", CreateTodo)
	r.Get("/todos/{id}", GetTodo)
	r.Put("/todos/{id}", UpdateTodo)
	r.Delete("/todos/{id}", DeleteTodo)
	
	log.Println("Server on :3000")
	http.ListenAndServe(":3000", r)
}
```

### Ejercicio 2: Router Composition ⭐⭐

**Objetivo:** Montar routers composables

```go
// Crear dos routers independientes:
// 1. UserRouter - /users endpoints
// 2. ProductRouter - /products endpoints
//
// Montarlos en router principal bajo /api
// Cada uno con su middleware específico

// Pasos:
// 1. Crear UserRouter()
// 2. Crear ProductRouter()
// 3. En main(), montar ambos
// 4. Tener handlers diferentes para cada uno
```

**Solución:**

```go
package main

import (
	"encoding/json" 
	"log"
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// User router
func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"user1", "user2"})
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	return r
}

// Product router
func ProductRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"prod1", "prod2"})
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	return r
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.Route("/api", func(r chi.Router) {
		r.Mount("/users", UserRouter())
		r.Mount("/products", ProductRouter())
	})
	
	log.Println("Server on :3000")
	http.ListenAndServe(":3000", r)
}
```

### Ejercicio 3: Middleware Chain ⭐⭐

**Objetivo:** Crear una cadena de middleware personalizado

```go
// Crear:
// 1. AuthMiddleware - valida token
// 2. LoggingMiddleware - registra peticiones
// 3. RateLimitMiddleware - limita requests
//
// Apilarlos en orden correcto
// Probar que funcionan juntos
```

**Solución:**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	
	"github.com/go-chi/chi/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := make(map[string]int)
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter[ip]++
		
		if limiter[ip] > 100 {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()
	
	r.Route("/api", func(r chi.Router) {
		r.Use(LoggingMiddleware)
		r.Use(RateLimitMiddleware)
		r.Use(AuthMiddleware)
		
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Protected resource"))
		})
	})
	
	log.Println("Server on :3000")
	http.ListenAndServe(":3000", r)
}
```

### Ejercicio 4: Microservicio Completo ⭐⭐⭐

**Objetivo:** Construir un microservicio real con Chi

```go
// Crear un servicio de libros con:
// 1. Endpoints CRUD completos
// 2. Validación de entrada
// 3. Manejo de errores
// 4. Middleware personalizado
// 5. Tests unitarios
```

Solución en siguiente sección.

### Ejercicio 5: API Gateway ⭐⭐⭐

**Objetivo:** Crear gateway que agrega múltiples servicios

```go
// Crear gateway que:
// 1. Enruta a servicio de usuarios
// 2. Enruta a servicio de posts
// 3. Enruta a servicio de comentarios
// 4. Añade rate limiting global
// 5. Logging centralizado
```

---

## ANTI-PATTERNS ❌ vs BEST PRACTICES ✅

### ❌ Anti-pattern: Monolithic Router

```go
// MAL - Todo en un solo lugar
func main() {
	r := chi.NewRouter()
	
	// Usuarios
	r.Get("/users", getUsers)
	r.Post("/users", createUser)
	r.Get("/users/{id}", getUser)
	
	// Posts
	r.Get("/posts", getPosts)
	r.Post("/posts", createPost)
	
	// Productos
	r.Get("/products", getProducts)
	
	// Comentarios
	r.Get("/comments", getComments)
	
	// ... 200 líneas más
	
	http.ListenAndServe(":3000", r)
}
```

### ✅ Best Practice: Composable Routers

```go
// BIEN - Modular y escalable
func main() {
	r := chi.NewRouter()
	
	r.Mount("/api/users", routes.UserRouter())
	r.Mount("/api/posts", routes.PostRouter())
	r.Mount("/api/products", routes.ProductRouter())
	r.Mount("/api/comments", routes.CommentRouter())
	
	http.ListenAndServe(":3000", r)
}
```

### ❌ Anti-pattern: Global State

```go
// MAL - Estado global compartido
var db *sql.DB
var cache map[string]string

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Usar db y cache global - difícil de testear
}
```

### ✅ Best Practice: Dependency Injection

```go
// BIEN - Inyección de dependencias
type Handler struct {
	db    *sql.DB
	cache Cache
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Usar h.db y h.cache - fácil de testear
}
```

### ❌ Anti-pattern: No Error Handling

```go
// MAL - Sin manejo de errores
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user) // ¿Qué si falla?
	db.SaveUser(user)                      // ¿Qué si falla?
	json.NewEncoder(w).Encode(user)        // ¿Qué si falla?
}
```

### ✅ Best Practice: Error Handling Explícito

```go
// BIEN - Manejo explícito de errores
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	if err := db.SaveUser(user); err != nil {
		respondError(w, http.StatusInternalServerError, 
			"Database error")
		return
	}
	
	respondJSON(w, http.StatusCreated, user)
}
```

### ❌ Anti-pattern: Middleware Monolítico

```go
// MAL - Middleware hace demasiado
func SuperMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging
		// Auth
		// CORS
		// Rate limiting
		// Validación
		// Transformación
		// ... 100 líneas más
		next.ServeHTTP(w, r)
	})
}
```

### ✅ Best Practice: Middleware Composable

```go
// BIEN - Cada middleware hace UNA cosa
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(CORSMiddleware)

r.Route("/api", func(r chi.Router) {
	r.Use(AuthMiddleware)
	r.Use(RateLimitMiddleware)
	
	// ... rutas
})
```

---

## CONCLUSIÓN

Chi es el router ideal para:
- ✅ Arquitecturas de microservicios
- ✅ Máxima composabilidad
- ✅ Performance critical applications
- ✅ Go idiomático y minimalista
- ✅ Equipos que valoran control y claridad

Con su filosofía "Small core, big ecosystem", Chi te proporciona exactamente lo que necesitas para enrutamiento HTTP, dejando todas las otras decisiones en tus manos. Esto lo hace perfecto para proyectos donde flexibilidad y composabilidad son prioritarias.

---

**Próximos pasos:**
- Explorar el ecosystem de Chi middleware
- Implementar patrones avanzados de composición
- Integrar con bases de datos y sistemas externos
- Containerizar aplicaciones Chi
- Desplegar en Kubernetes

---

*Capítulo 60 - Go: Guía Exhaustiva*  
*Chi Router: Lightweight & Composable Routing*  
*Versión 1.0 - Marzo 2024*
