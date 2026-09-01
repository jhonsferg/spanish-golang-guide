# Capítulo 38: HTTP avanzado - Servidores web y APIs robustas

## Índice del Capítulo 38

1. [38.1 Arquitectura de Servidores HTTP](#381-arquitectura-de-servidores-http)
2. [38.2 Enrutamiento Avanzado](#382-enrutamiento-avanzado)
3. [38.3 Middleware y Request Pipeline](#383-middleware-y-request-pipeline)
4. [38.4 Parsing de Requests](#384-parsing-de-requests)
5. [38.5 Escritura de Responses](#385-escritura-de-responses)
6. [38.6 Manejo de Errores HTTP](#386-manejo-de-errores-http)
7. [38.7 Autenticación y Autorización](#387-autenticación-y-autorización)
8. [38.8 CORS y Seguridad](#388-cors-y-seguridad)
9. [38.9 Graceful Shutdown](#389-graceful-shutdown)
10. [38.10 Testing de Servidores HTTP](#3810-testing-de-servidores-http)
11. [38.11 Buenas Prácticas y Patterns](#3811-buenas-prácticas-y-patterns)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 38.1 Arquitectura de Servidores HTTP

### El Ciclo de Vida de un Request HTTP

Entender cómo Go procesa requests HTTP es fundamental para construir APIs robustas:

```
┌─────────────────────────────────────────────────────────────┐
│                    Ciclo de Vida HTTP en Go                │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. LISTENING PHASE                                         │
│     └─ Server escucha en puerto (ej: 8080)                 │
│                                                              │
│  2. CONNECTION PHASE                                        │
│     └─ Cliente establece conexión TCP                      │
│     └─ Se crea un nuevo goroutine para cada conexión       │
│                                                              │
│  3. REQUEST READING PHASE                                   │
│     └─ Go lee: método, path, headers, body                 │
│     └─ Parsea según RFC 7231                               │
│                                                              │
│  4. ROUTING PHASE                                           │
│     └─ Se busca handler que coincida con el path           │
│     └─ Se ejecuta la cadena de middleware                  │
│                                                              │
│  5. HANDLER EXECUTION                                       │
│     └─ Handler procesa el request                          │
│     └─ Handler escribe response                            │
│                                                              │
│  6. RESPONSE WRITING PHASE                                  │
│     └─ Status line: HTTP/1.1 200 OK                        │
│     └─ Headers: Content-Type, Content-Length, etc          │
│     └─ Body (si aplica)                                    │
│                                                              │
│  7. CONNECTION CLOSE PHASE                                  │
│     └─ Keep-alive o close según headers                    │
│     └─ Conexión reutilizable o cerrada                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Estructura de un Handler en Go

Un **handler** es cualquier función que cumple la interfaz:

```go
type HandlerFunc func(ResponseWriter, *Request)
```

Esto significa que CUALQUIER función con esa firma es un handler válido:

```go
package main

import (
    "fmt"
    "net/http"
)

// Handler básico
func saludo(w http.ResponseWriter, r *http.Request) {
    // w es donde escribimos la respuesta
    // r contiene todo sobre el request
    fmt.Fprintf(w, "¡Hola, %s!", r.URL.Path)
}

// Handler con lógica más compleja
func apiUsuario(w http.ResponseWriter, r *http.Request) {
    // 1. Validar método HTTP
    if r.Method != http.MethodGet {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    // 2. Parsear parámetros
    id := r.URL.Query().Get("id")

    // 3. Procesar lógica
    usuario := buscarUsuario(id)
    if usuario == nil {
        http.Error(w, "Usuario no encontrado", http.StatusNotFound)
        return
    }

    // 4. Escribir respuesta
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"id":"%s","nombre":"%s"}`, usuario.ID, usuario.Nombre)
}

func main() {
    // Registrar handlers
    http.HandleFunc("/saludo", saludo)
    http.HandleFunc("/api/usuario", apiUsuario)

    // Iniciar servidor
    http.ListenAndServe(":8080", nil)
}
```

### La Interfaz http.Handler

La interfaz `http.Handler` es el corazón del sistema:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Esto significa que `http.HandlerFunc` es solo una conveniencia. Puedes crear tus propias estructuras:

```go
// Handler como estructura
type LoggingHandler struct {
    handler http.Handler
    logger  *log.Logger
}

// Implementar la interfaz
func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.logger.Printf("%s %s", r.Method, r.URL.Path)
    h.handler.ServeHTTP(w, r)
}

// Uso
mux := http.NewServeMux()
mux.HandleFunc("/api", apiHandler)

logged := &LoggingHandler{
    handler: mux,
    logger:  log.New(os.Stdout, "HTTP: ", 0),
}

http.ListenAndServe(":8080", logged)
```

### Goroutines y Concurrencia en Servidores

**Cada request se procesa en su propio goroutine**. Esto es automático:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Este código se ejecuta en un goroutine separado
    // Si aquí se lanzan 1000 requests simultáneos,
    // Go crea 1000 goroutines automáticamente
    
    fmt.Println(runtime.NumGoroutine()) // Ver cantidad de goroutines activos
}
```

**Implicaciones importantes:**

1. **Variable compartida entre goroutines debe ser thread-safe:**

```go
// ✗ PELIGROSO: Race condition
var contador int

func handler(w http.ResponseWriter, r *http.Request) {
    contador++  // Múltiples goroutines accediendo sin sincronización
}

// ✓ CORRECTO: Usar sync.Mutex
var (
    contador int
    mu       sync.Mutex
)

func handler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    contador++
    mu.Unlock()
}

// ✓ MEJOR: Usar canales
var contadorChan = make(chan int)

func handler(w http.ResponseWriter, r *http.Request) {
    contadorChan <- 1  // Enviar incremente a través del canal
}
```

2. **Context puede propagarse a cada request:**

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // El context tiene deadline, valores, y se puede cancelar
    select {
    case <-ctx.Done():
        // Request cancelado o deadline alcanzado
        return
    case resultado := <-procesarEnAsyncrono(ctx):
        fmt.Fprintf(w, "%v", resultado)
    }
}
```

### Estructura de Proyecto típico de API

```
mi-api/
├── main.go              # Punto de entrada
├── handlers/            
│   ├── users.go         # Handlers de usuarios
│   ├── products.go      # Handlers de productos
│   └── health.go        # Health checks
├── middleware/
│   ├── logging.go       # Middleware de logging
│   ├── auth.go          # Middleware de autenticación
│   └── cors.go          # Middleware de CORS
├── models/
│   └── user.go          # Estructuras de datos
├── db/
│   └── db.go            # Conexión a base de datos
└── config/
    └── config.go        # Configuración
```

---

## 38.2 Enrutamiento Avanzado

### http.ServeMux: El Enrutador Básico

Go proporciona `http.ServeMux` como enrutador integrado:

```go
mux := http.NewServeMux()

// Registrar rutas exactas
mux.HandleFunc("/", homeHandler)
mux.HandleFunc("/api/usuarios", usuariosHandler)

// Registrar con patrón (path matching)
mux.HandleFunc("/api/usuarios/", usuarioDetalleHandler)
```

**Matching de rutas:**

```go
// Ruta exacta
mux.HandleFunc("/api/usuarios", h1)
// ✓ GET /api/usuarios
// ✗ GET /api/usuarios/123

// Ruta con barra al final (subtree)
mux.HandleFunc("/api/usuarios/", h2)
// ✓ GET /api/usuarios/
// ✓ GET /api/usuarios/123
// ✓ GET /api/usuarios/123/posts

// Raíz (catch-all)
mux.HandleFunc("/", h3)
// ✓ GET /algo (si no hay otra ruta más específica)
```

**Extrayendo parámetros de URL:**

```go
func usuarioHandler(w http.ResponseWriter, r *http.Request) {
    // Con http.ServeMux, hay que parsear manualmente
    
    // Opción 1: Usando strings
    id := strings.TrimPrefix(r.URL.Path, "/usuarios/")
    
    // Opción 2: Usando Query parameters
    id := r.URL.Query().Get("id")
    
    // Opción 3: Usando path parameters (más moderno)
    // Requiere router personalizado
}
```

### Routers Personalizados: Patrón chi/gorilla/echo

Para aplicaciones más complejas, Go tiene varios routers populares. Aquí está el enfoque usando el patrón de middleware:

```go
package main

import (
    "net/http"
    "strings"
)

// Router personalizado simple
type Router struct {
    routes map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
    return &Router{
        routes: make(map[string]map[string]http.HandlerFunc),
    }
}

// Registrar ruta
func (r *Router) Register(method, path string, handler http.HandlerFunc) {
    if r.routes[method] == nil {
        r.routes[method] = make(map[string]http.HandlerFunc)
    }
    r.routes[method][path] = handler
}

// Implementar interfaz http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // Buscar handler
    if handlers, ok := r.routes[req.Method]; ok {
        if handler, ok := handlers[req.URL.Path]; ok {
            handler(w, req)
            return
        }
    }

    // No encontrado
    http.NotFound(w, req)
}

// Métodos de conveniencia
func (r *Router) Get(path string, handler http.HandlerFunc) {
    r.Register(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
    r.Register(http.MethodPost, path, handler)
}

// Uso
func main() {
    router := NewRouter()
    
    router.Get("/", homeHandler)
    router.Post("/usuarios", crearUsuarioHandler)
    
    http.ListenAndServe(":8080", router)
}
```

### Extracción de Parámetros Avanzada

```go
// Router con parámetros en path
type AdvancedRouter struct {
    patterns map[string]map[string]*pathPattern
}

type pathPattern struct {
    regex   string
    handler http.HandlerFunc
}

// Ejemplo: /usuarios/:id/posts/:postId
func (r *AdvancedRouter) Get(pattern string, handler http.HandlerFunc) {
    // Convertir /usuarios/:id a expresión regular
    // /usuarios/(?P<id>[^/]+)
}
```

### Comparación: Go vs Express vs FastAPI

```python
# EXPRESS (Node.js)
app.get('/usuarios/:id', (req, res) => {
    const id = req.params.id;
    res.json({id});
});

# FASTAPI (Python)
@app.get("/usuarios/{id}")
async def obtener_usuario(id: str):
    return {"id": id}

# GO - Sin framework
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/usuarios/", func(w http.ResponseWriter, r *http.Request) {
        id := strings.TrimPrefix(r.URL.Path, "/usuarios/")
        fmt.Fprintf(w, `{"id":"%s"}`, id)
    })
    http.ListenAndServe(":8080", mux)
}

# GO - Con framework (chi, gorilla, echo)
// Mucho más simple, pero requiere dependencia externa
```

---

## 38.3 Middleware y Request Pipeline

### Concepto: Pipeline de Middleware

El middleware es un patrón que permite procesar requests antes y después de los handlers:

```
Request
   │
   ├─→ [Middleware 1: Logger] ─→ Registra entrada
   │
   ├─→ [Middleware 2: Auth] ─→ Valida token JWT
   │
   ├─→ [Middleware 3: CORS] ─→ Añade headers CORS
   │
   ├─→ [HANDLER] ─→ Procesa lógica de negocio
   │
   ├─→ [Middleware 3: CORS] ─→ (en orden inverso)
   │
   ├─→ [Middleware 2: Auth] ─→
   │
   ├─→ [Middleware 1: Logger] ─→ Registra salida
   │
   └─→ Response

```

### Patrón: Cadena de Middleware

```go
// Función para encadenar middleware
func chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
    for _, m := range middlewares {
        handler = m(handler)
    }
    return handler
}

// Middleware de logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        inicio := time.Now()
        
        // Llamar al siguiente handler
        next(w, r)
        
        duracion := time.Since(inicio)
        log.Printf("%s %s completado en %v", r.Method, r.URL.Path, duracion)
    }
}

// Middleware de autenticación
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token == "" {
            http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
            return
        }
        
        if !validarToken(token) {
            http.Error(w, "Token inválido", http.StatusForbidden)
            return
        }
        
        // Token válido, continuar
        next(w, r)
    }
}

// Uso
func main() {
    mux := http.NewServeMux()
    
    // Sin middleware
    mux.HandleFunc("/public", chain(publicHandler))
    
    // Con middleware
    mux.HandleFunc("/protected", chain(protectedHandler, authMiddleware, loggingMiddleware))
    
    http.ListenAndServe(":8080", mux)
}
```

### Middleware Avanzado con ResponseWriter personalizado

A veces necesitas interceptar la respuesta (por ejemplo, para cambiar status code):

```go
// ResponseWriter personalizado que intercepta escrituras
type ResponseInterceptor struct {
    http.ResponseWriter
    statusCode int
    body       []byte
}

func (r *ResponseInterceptor) WriteHeader(code int) {
    r.statusCode = code
    r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseInterceptor) Write(b []byte) (int, error) {
    r.body = append(r.body, b...)
    return r.ResponseWriter.Write(b)
}

// Middleware que intercepta respuestas
func responseInterceptorMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        interceptor := &ResponseInterceptor{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        next(interceptor, r)
        
        // Hacer algo con status code y body
        log.Printf("Status: %d, Body size: %d bytes", interceptor.statusCode, len(interceptor.body))
    }
}
```

### Middleware Común: CORS

```go
// CORS Middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Permitir todos los orígenes (NO recomendado en producción)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        
        // Permitir métodos
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        
        // Permitir headers personalizados
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Tiempo de caché de CORS preflight
        w.Header().Set("Access-Control-Max-Age", "3600")
        
        // Manejar preflight request
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}
```

### Middleware Común: Timeout

```go
// Timeout middleware
func timeoutMiddleware(duracion time.Duration) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), duracion)
            defer cancel()
            
            // Crear nuevo request con contexto
            r = r.WithContext(ctx)
            
            // Canal para saber si handler terminó
            done := make(chan struct{})
            go func() {
                next(w, r)
                done <- struct{}{}
            }()
            
            select {
            case <-done:
                // Handler completó normalmente
            case <-ctx.Done():
                // Timeout
                http.Error(w, "Request timeout", http.StatusRequestTimeout)
            }
        }
    }
}
```

---

## 38.4 Parsing de Requests

### Tipos de Request Body

Go puede parsear automáticamente varios formatos:

```go
// 1. Query Parameters
func queryHandler(w http.ResponseWriter, r *http.Request) {
    // GET /search?q=golang&limit=10
    
    q := r.URL.Query().Get("q")           // "golang"
    limit := r.URL.Query().Get("limit")   // "10"
    
    // Obtener valor múltiple
    tags := r.URL.Query()["tag"]  // []string{"golang", "http"}
}

// 2. JSON Body
type Producto struct {
    Nombre string `json:"nombre"`
    Precio float64 `json:"precio"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    var producto Producto
    
    err := json.NewDecoder(r.Body).Decode(&producto)
    if err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }
    
    // Usar producto
    fmt.Printf("Producto: %s, Precio: $%.2f\n", producto.Nombre, producto.Precio)
}

// 3. Form Data (application/x-www-form-urlencoded)
func formHandler(w http.ResponseWriter, r *http.Request) {
    // POST con Content-Type: application/x-www-form-urlencoded
    
    r.ParseForm()
    
    nombre := r.FormValue("nombre")
    email := r.FormValue("email")
}
```

### Multipart File Upload

```go
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Límite de tamaño: 10MB
    r.ParseMultipartForm(10 << 20)
    
    // Obtener el archivo del form
    archivo, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error al obtener archivo", http.StatusBadRequest)
        return
    }
    defer archivo.Close()
    
    // Información del archivo
    fmt.Printf("Nombre: %s\n", handler.Filename)
    fmt.Printf("Tamaño: %d bytes\n", handler.Size)
    fmt.Printf("Tipo MIME: %s\n", handler.Header.Get("Content-Type"))
    
    // Guardar archivo
    rutaDestino := fmt.Sprintf("uploads/%s", handler.Filename)
    destino, err := os.Create(rutaDestino)
    if err != nil {
        http.Error(w, "Error al guardar", http.StatusInternalServerError)
        return
    }
    defer destino.Close()
    
    // Copiar contenido
    _, err = io.Copy(destino, archivo)
    if err != nil {
        http.Error(w, "Error al copiar", http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "Archivo guardado en: %s", rutaDestino)
}
```

### Validación de Input

```go
// Estructura con tags de validación
type Usuario struct {
    Nombre string `json:"nombre" validate:"required,min=3,max=50"`
    Email  string `json:"email" validate:"required,email"`
    Edad   int    `json:"edad" validate:"required,min=18,max=150"`
}

// Función de validación simple
func validarUsuario(u Usuario) error {
    if u.Nombre == "" || len(u.Nombre) < 3 {
        return errors.New("nombre debe tener al menos 3 caracteres")
    }
    
    if !strings.Contains(u.Email, "@") {
        return errors.New("email inválido")
    }
    
    if u.Edad < 18 || u.Edad > 150 {
        return errors.New("edad debe estar entre 18 y 150")
    }
    
    return nil
}

func crearUsuarioHandler(w http.ResponseWriter, r *http.Request) {
    var usuario Usuario
    
    err := json.NewDecoder(r.Body).Decode(&usuario)
    if err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }
    
    // Validar
    if err := validarUsuario(usuario); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Procesar usuario válido
}
```

### Content Negotiation

```go
func dataHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "nombre": "Juan",
        "edad":   30,
    }
    
    // Parsear Accept header
    tipoAceptado := r.Header.Get("Accept")
    
    if strings.Contains(tipoAceptado, "application/xml") {
        // Retornar XML
        w.Header().Set("Content-Type", "application/xml")
        xml.NewEncoder(w).Encode(data)
    } else if strings.Contains(tipoAceptado, "application/json") {
        // Retornar JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    } else if strings.Contains(tipoAceptado, "text/plain") {
        // Retornar texto
        w.Header().Set("Content-Type", "text/plain")
        fmt.Fprintf(w, "Nombre: %s, Edad: %d", data["nombre"], data["edad"])
    } else {
        // Default: JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    }
}
```

---

## 38.5 Escritura de Responses

### Headers HTTP

```go
func responseHandler(w http.ResponseWriter, r *http.Request) {
    // Establecer headers (ANTES de WriteHeader)
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Custom-Header", "valor")
    w.Header().Add("Set-Cookie", "session=abc123")  // Add permite múltiples valores
    
    // Escribir status code (solo una vez!)
    w.WriteHeader(http.StatusOK)
    
    // Escribir body
    fmt.Fprintf(w, `{"status":"ok"}`)
}
```

**Orden importante:**
1. Primero `w.Header().Set()` (múltiples veces)
2. Luego `w.WriteHeader()` (solo una vez)
3. Finalmente escribir body

### Status Codes Comunes

```go
// 2xx - Éxito
http.StatusOK                   // 200
http.StatusCreated              // 201
http.StatusAccepted             // 202
http.StatusNoContent            // 204

// 3xx - Redirección
http.StatusMovedPermanently    // 301
http.StatusFound                // 302
http.StatusTemporaryRedirect    // 307

// 4xx - Error del cliente
http.StatusBadRequest           // 400
http.StatusUnauthorized         // 401
http.StatusForbidden            // 403
http.StatusNotFound             // 404
http.StatusConflict             // 409
http.StatusUnprocessableEntity  // 422

// 5xx - Error del servidor
http.StatusInternalServerError  // 500
http.StatusNotImplemented       // 501
http.StatusServiceUnavailable   // 503
```

### Redirecciones

```go
func redirectHandler(w http.ResponseWriter, r *http.Request) {
    // Redirección simple (302)
    http.Redirect(w, r, "/nueva-url", http.StatusFound)
    
    // Redirección permanente (301)
    http.Redirect(w, r, "/nueva-url", http.StatusMovedPermanently)
}
```

### Streaming de Datos Grandes

```go
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    // Para archivos grandes, usar streaming en lugar de cargar todo en memoria
    
    archivo, err := os.Open("grande-archivo.zip")
    if err != nil {
        http.Error(w, "Archivo no encontrado", http.StatusNotFound)
        return
    }
    defer archivo.Close()
    
    // Información del archivo
    info, _ := archivo.Stat()
    tamaño := info.Size()
    
    // Headers para descarga
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", `attachment; filename="archivo.zip"`)
    w.Header().Set("Content-Length", strconv.FormatInt(tamaño, 10))
    
    // Copiar contenido directamente al response (streaming)
    io.Copy(w, archivo)
}
```

### JSON Response Helper

```go
// Función auxiliar para respuestas JSON
func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

// Uso
func productosHandler(w http.ResponseWriter, r *http.Request) {
    productos := []map[string]interface{}{
        {"id": 1, "nombre": "Laptop"},
        {"id": 2, "nombre": "Mouse"},
    }
    
    writeJSON(w, http.StatusOK, productos)
}
```

---

## 38.6 Manejo de Errores HTTP

### Estrategia de Errores en APIs

```
Error al procesar request
        │
        ├─ Validación fallida → 400 Bad Request
        │
        ├─ No autenticado → 401 Unauthorized
        │
        ├─ No autorizado → 403 Forbidden
        │
        ├─ Recurso no encontrado → 404 Not Found
        │
        ├─ Error en servidor → 500 Internal Server Error
        │
        └─ Servicio no disponible → 503 Service Unavailable
```

### Estructura de Error Estándar

```go
// Estructura de error API
type ErrorResponse struct {
    Codigo  string      `json:"codigo"`
    Mensaje string      `json:"mensaje"`
    Detalles interface{} `json:"detalles,omitempty"`
    Timestamp int64     `json:"timestamp"`
}

// Función auxiliar para errores
func writeError(w http.ResponseWriter, status int, codigo, mensaje string) {
    error := ErrorResponse{
        Codigo:    codigo,
        Mensaje:   mensaje,
        Timestamp: time.Now().UnixMilli(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(error)
}

// Uso
func usuarioHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    
    if id == "" {
        writeError(w, http.StatusBadRequest, "PARAM_MISSING", "Parámetro 'id' requerido")
        return
    }
    
    usuario, err := buscarUsuario(id)
    if err != nil {
        writeError(w, http.StatusNotFound, "USUARIO_NO_ENCONTRADO", "Usuario no existe")
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(usuario)
}
```

### Recovery de Panic en Handlers

```go
// Middleware para capturar panics
func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("PANIC: %v\n%s", err, debug.Stack())
                
                writeError(w, http.StatusInternalServerError, 
                    "INTERNAL_ERROR", "Error interno del servidor")
            }
        }()
        
        next(w, r)
    }
}
```

### Timeout y Context Cancellation

```go
func handlerConTimeout(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Operación con timeout
    resultado := make(chan interface{}, 1)
    
    go func() {
        // Simulación de operación lenta
        time.Sleep(2 * time.Second)
        resultado <- "datos"
    }()
    
    select {
    case res := <-resultado:
        json.NewEncoder(w).Encode(res)
    case <-ctx.Done():
        switch ctx.Err() {
        case context.Canceled:
            writeError(w, http.StatusBadRequest, "CANCELLED", "Request cancelado")
        case context.DeadlineExceeded:
            writeError(w, http.StatusGatewayTimeout, "TIMEOUT", "Request timeout")
        }
    }
}
```

---

## 38.7 Autenticación y Autorización

### API Key Authentication

```go
const ValidAPIKey = "sk-1234567890abcdef"

func validarAPIKey(apiKey string) bool {
    return apiKey == ValidAPIKey
}

func apiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        
        if apiKey == "" {
            writeError(w, http.StatusUnauthorized, "MISSING_KEY", "API Key requerida")
            return
        }
        
        if !validarAPIKey(apiKey) {
            writeError(w, http.StatusForbidden, "INVALID_KEY", "API Key inválida")
            return
        }
        
        next(w, r)
    }
}
```

### JWT Token Authentication

```go
import "github.com/golang-jwt/jwt/v5"

const jwtSecret = "mi-secreto-muy-seguro"

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// Generar token
func generarToken(userID int, username string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}

// Validar token
func validarToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })
    
    if err != nil || !token.Valid {
        return nil, err
    }
    
    return claims, nil
}

// Middleware JWT
func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Obtener token del header Authorization
        authHeader := r.Header.Get("Authorization")
        
        if authHeader == "" {
            writeError(w, http.StatusUnauthorized, "MISSING_TOKEN", "Token requerido")
            return
        }
        
        // Formato: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            writeError(w, http.StatusUnauthorized, "INVALID_FORMAT", "Formato inválido")
            return
        }
        
        claims, err := validarToken(parts[1])
        if err != nil {
            writeError(w, http.StatusForbidden, "INVALID_TOKEN", "Token inválido")
            return
        }
        
        // Guardar claims en contexto
        ctx := context.WithValue(r.Context(), "claims", claims)
        next(w, r.WithContext(ctx))
    }
}

// Uso en handler
func protectedHandler(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("claims").(*Claims)
    
    fmt.Fprintf(w, "Usuario: %s (ID: %d)", claims.Username, claims.UserID)
}
```

### Basic Auth

```go
func basicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        username, password, ok := r.BasicAuth()
        
        if !ok || !validarCredenciales(username, password) {
            w.Header().Set("WWW-Authenticate", `Basic realm="API"`)
            writeError(w, http.StatusUnauthorized, "INVALID_AUTH", "Credenciales inválidas")
            return
        }
        
        next(w, r)
    }
}

func validarCredenciales(username, password string) bool {
    // En producción: verificar contra base de datos con contraseña hasheada
    return username == "admin" && password == "secret123"
}
```

### Session Cookies

```go
var cookieStore *sessions.CookieStore

func init() {
    var err error
    cookieStore, err = sessions.NewCookieStore([]byte("secret-key-muy-segura"))
    if err != nil {
        panic(err)
    }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    
    // Validar credenciales
    if username != "admin" || password != "123456" {
        http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
        return
    }
    
    // Crear sesión
    session, _ := cookieStore.Get(r, "session")
    session.Values["user_id"] = 123
    session.Values["username"] = username
    
    session.Save(r, w)
    
    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func sessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := cookieStore.Get(r, "session")
        
        if auth, ok := session.Values["user_id"]; !ok || auth == nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        
        next(w, r)
    }
}
```

---

## 38.8 CORS y Seguridad

### CORS Completo

```go
type CORSConfig struct {
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    ExposedHeaders   []string
    MaxAge           int
    AllowCredentials bool
}

func corsHandler(config CORSConfig) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Verificar si origen está permitido
            allowedOrigin := false
            for _, o := range config.AllowedOrigins {
                if o == "*" || o == origin {
                    allowedOrigin = true
                    break
                }
            }
            
            if !allowedOrigin {
                next(w, r)
                return
            }
            
            // Headers CORS
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
            w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
            w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
            w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
            
            if config.AllowCredentials {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
            
            // Manejar preflight
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next(w, r)
        }
    }
}
```

### Headers de Seguridad

```go
func securityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Prevenir clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        
        // Prevenir MIME sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        
        // Habilitar XSS protection en navegadores antiguos
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        
        // HTTP Strict Transport Security (solo HTTPS)
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // No cachear datos sensibles
        w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        
        next(w, r)
    }
}
```

### CSRF Protection

```go
import "github.com/gorilla/csrf"

func csrfMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Obtener token CSRF del header o cookie
        csrfToken := csrf.Token(r)
        
        // Para formularios, incluir en template: {{ csrf.TemplateField }}
        
        // Para AJAX, comparar con header X-CSRF-Token
        if r.Method != http.MethodGet && r.Method != http.MethodOptions {
            headerToken := r.Header.Get("X-CSRF-Token")
            
            if headerToken != csrfToken {
                writeError(w, http.StatusForbidden, "INVALID_CSRF", "Token CSRF inválido")
                return
            }
        }
        
        next(w, r)
    }
}
```

### Rate Limiting

```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
    }
}

func (rl *RateLimiter) limitarMiddleware(requestsPerSecond float64) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // IP del cliente como identificador
            ip := r.RemoteAddr
            
            rl.mu.Lock()
            limiter, exists := rl.limiters[ip]
            if !exists {
                limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), int(requestsPerSecond))
                rl.limiters[ip] = limiter
            }
            rl.mu.Unlock()
            
            if !limiter.Allow() {
                writeError(w, http.StatusTooManyRequests, "RATE_LIMIT", "Demasiadas solicitudes")
                return
            }
            
            next(w, r)
        }
    }
}
```

---

## 38.9 Graceful Shutdown

### El Problema: Cerrar Limpiamente

Cuando apagamos un servidor HTTP, queremos:
1. No rechazar nuevos requests bruscamente
2. Permitir que requests en progreso terminen
3. Cerrar conexiones de base de datos
4. Liberar recursos correctamente

```
Shutdown Grácil vs Brusco:

SHUTDOWN BRUSCO:
Server
   │
   ├─ Request 1 (50% completado) ─── CORTADO ✗
   ├─ Request 2 (70% completado) ─── CORTADO ✗
   ├─ Conexión DB abierta ─────────── MEMORY LEAK ✗
   │
   └─ APAGADO (muchos errores)

SHUTDOWN GRÁCIL:
Server
   │
   ├─ [SEÑAL: SIGTERM]
   ├─ [Dejar de aceptar nuevos requests]
   │
   ├─ Request 1 (50% → 100% completado) ✓
   ├─ Request 2 (70% → 100% completado) ✓
   │
   ├─ Esperar timeout (ej: 30 segundos)
   ├─ [Cerrar conexión DB]
   ├─ [Liberar recursos]
   │
   └─ APAGADO (sin errores)
```

### Implementación: Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Crear servidor
    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)
    
    server := &http.Server{
        Addr:           ":8080",
        Handler:        mux,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    
    // Canal para errores del servidor
    serverErrors := make(chan error, 1)
    
    // Iniciar servidor en goroutine
    go func() {
        log.Printf("Servidor iniciado en %s", server.Addr)
        serverErrors <- server.ListenAndServe()
    }()
    
    // Canal para señales del SO
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case err := <-serverErrors:
        if err != http.ErrServerClosed {
            log.Fatalf("Error del servidor: %v", err)
        }
    case sig := <-quit:
        log.Printf("Señal recibida: %v", sig)
        
        // Contexto con timeout de 30 segundos
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        // Shutdown graceful
        if err := server.Shutdown(ctx); err != nil {
            log.Fatalf("Error al apagar servidor: %v", err)
        }
        
        log.Println("Servidor apagado correctamente")
    }
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    // Simular procesamiento
    time.Sleep(2 * time.Second)
    fmt.Fprintf(w, "OK")
}
```

### Cleanup de Recursos

```go
// Estructura que maneja recursos
type App struct {
    server *http.Server
    db     *sql.DB
}

func (app *App) Shutdown(ctx context.Context) error {
    // Cerrar en orden correcto: últero en crear, primero en cerrar
    
    // 1. Apagar HTTP server
    if err := app.server.Shutdown(ctx); err != nil {
        return fmt.Errorf("error apagando servidor HTTP: %w", err)
    }
    
    // 2. Cerrar conexión DB
    if app.db != nil {
        if err := app.db.Close(); err != nil {
            return fmt.Errorf("error cerrando DB: %w", err)
        }
    }
    
    return nil
}

func main() {
    // Crear recursos
    db, _ := sql.Open("postgres", "...")
    
    app := &App{
        server: &http.Server{Addr: ":8080"},
        db:     db,
    }
    
    // Manejar shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-quit
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := app.Shutdown(ctx); err != nil {
            log.Fatalf("Error en shutdown: %v", err)
        }
    }()
    
    app.server.ListenAndServe()
}
```

### Connection Draining

```go
// Rastrear requests activos
type ConnectionDrainer struct {
    activeRequests int32
    mu             sync.Mutex
}

func (cd *ConnectionDrainer) Middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cd.mu.Lock()
        cd.activeRequests++
        cd.mu.Unlock()
        
        defer func() {
            cd.mu.Lock()
            cd.activeRequests--
            cd.mu.Unlock()
        }()
        
        next(w, r)
    }
}

func (cd *ConnectionDrainer) Wait(timeout time.Duration) error {
    deadline := time.Now().Add(timeout)
    
    for {
        cd.mu.Lock()
        active := cd.activeRequests
        cd.mu.Unlock()
        
        if active == 0 {
            return nil
        }
        
        if time.Now().After(deadline) {
            return fmt.Errorf("timeout esperando %d requests", active)
        }
        
        time.Sleep(100 * time.Millisecond)
    }
}
```

---

## 38.10 Testing de Servidores HTTP

### Testing con httptest

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthCheck(t *testing.T) {
    // Crear un request
    req := httptest.NewRequest("GET", "/health", nil)
    
    // Crear un recorder para capturar respuesta
    w := httptest.NewRecorder()
    
    // Ejecutar handler
    healthHandler(w, req)
    
    // Verificar status code
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
    
    // Verificar body
    if body := w.Body.String(); body != "OK" {
        t.Errorf("Expected 'OK', got '%s'", body)
    }
}

func TestGetUsuario(t *testing.T) {
    req := httptest.NewRequest("GET", "/usuarios?id=123", nil)
    w := httptest.NewRecorder()
    
    usuarioHandler(w, req)
    
    // Verificar JSON response
    if w.Header().Get("Content-Type") != "application/json" {
        t.Fatal("Content-Type no es JSON")
    }
    
    // Parsear respuesta
    var usuario Usuario
    json.NewDecoder(w.Body).Decode(&usuario)
    
    if usuario.ID != 123 {
        t.Errorf("Expected ID 123, got %d", usuario.ID)
    }
}
```

### Testing con Servidor Completo

```go
func TestAPIIntegration(t *testing.T) {
    // Crear servidor de prueba
    server := httptest.NewServer(http.HandlerFunc(apiHandler))
    defer server.Close()
    
    // Hacer request al servidor
    resp, err := http.Get(server.URL + "/usuarios")
    if err != nil {
        t.Fatal(err)
    }
    
    // Verificar respuesta
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected 200, got %d", resp.StatusCode)
    }
}
```

### Mocking de Dependencias

```go
// Interfaz para base de datos
type UserDB interface {
    GetUser(id string) (*User, error)
}

// Mock de UserDB
type MockUserDB struct {
    users map[string]*User
}

func (m *MockUserDB) GetUser(id string) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("no encontrado")
    }
    return user, nil
}

// Handler que usa la interfaz
func makeUserHandler(db UserDB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Query().Get("id")
        user, err := db.GetUser(id)
        if err != nil {
            http.Error(w, "No encontrado", http.StatusNotFound)
            return
        }
        json.NewEncoder(w).Encode(user)
    }
}

// Test
func TestUserHandler(t *testing.T) {
    mockDB := &MockUserDB{
        users: map[string]*User{
            "1": {ID: "1", Name: "Juan"},
        },
    }
    
    handler := makeUserHandler(mockDB)
    
    req := httptest.NewRequest("GET", "/?id=1", nil)
    w := httptest.NewRecorder()
    
    handler(w, req)
    
    if w.Code != http.StatusOK {
        t.Fatal("Expected 200")
    }
}
```

### Testing de Middleware

```go
func TestAuthMiddleware(t *testing.T) {
    // Handler protegido
    protected := authMiddleware(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    
    // Test sin token
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    
    protected(w, req)
    
    if w.Code != http.StatusUnauthorized {
        t.Error("Esperaba 401 sin token")
    }
    
    // Test con token válido
    req = httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    w = httptest.NewRecorder()
    
    protected(w, req)
    
    if w.Code != http.StatusOK {
        t.Error("Esperaba 200 con token válido")
    }
}
```

---

## 38.11 Buenas Prácticas y Patterns

### RESTful API Conventions

```
Convenciones REST:

RECURSOS: SUSTANTIVOS (no verbos)
GET    /usuarios           ← Lista de usuarios
POST   /usuarios           ← Crear usuario
GET    /usuarios/:id       ← Obtener usuario específico
PUT    /usuarios/:id       ← Actualizar usuario completo
PATCH  /usuarios/:id       ← Actualizar parcialmente
DELETE /usuarios/:id       ← Eliminar usuario

NIVELES DE HTTP:
GET    /usuarios           → 200 OK
POST   /usuarios           → 201 Created
PUT    /usuarios/:id       → 200 OK o 204 No Content
DELETE /usuarios/:id       → 204 No Content
GET    /usuarios/999       → 404 Not Found
POST   /usuarios (inválido) → 400 Bad Request
```

### Versionado de APIs

```go
// Opción 1: URL path
// GET /v1/usuarios
// GET /v2/usuarios

mux.HandleFunc("/v1/usuarios", usuariosHandlerV1)
mux.HandleFunc("/v2/usuarios", usuariosHandlerV2)

// Opción 2: Accept header
// GET /usuarios
// Accept: application/vnd.api+json;version=2

func versionMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        accept := r.Header.Get("Accept")
        
        if strings.Contains(accept, "version=2") {
            // Usar v2
        } else {
            // Usar v1
        }
        
        next(w, r)
    }
}

// Opción 3: Query parameter
// GET /usuarios?version=2
version := r.URL.Query().Get("version")
```

### Pagination

```go
type PaginatedResponse struct {
    Data      interface{} `json:"data"`
    Page      int         `json:"page"`
    PageSize  int         `json:"page_size"`
    Total     int         `json:"total"`
    TotalPages int        `json:"total_pages"`
}

func usuariosConPaginacion(w http.ResponseWriter, r *http.Request) {
    // Parámetros
    page := 1
    pageSize := 10
    
    if p := r.URL.Query().Get("page"); p != "" {
        page, _ = strconv.Atoi(p)
    }
    if ps := r.URL.Query().Get("page_size"); ps != "" {
        pageSize, _ = strconv.Atoi(ps)
    }
    
    // Validar
    if pageSize > 100 {
        pageSize = 100 // Máximo
    }
    
    // Consultar BD
    usuarios, total, _ := buscarUsuariosConPaginacion(page, pageSize)
    
    response := PaginatedResponse{
        Data:       usuarios,
        Page:       page,
        PageSize:   pageSize,
        Total:      total,
        TotalPages: (total + pageSize - 1) / pageSize,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### Logging Estructurado

```go
import "log/slog"

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        inicio := time.Now()
        
        // Interceptar status code
        wrapped := &responseInterceptor{ResponseWriter: w}
        
        next(wrapped, r)
        
        duracion := time.Since(inicio)
        
        // Log estructurado
        slog.Info("HTTP request",
            "metodo", r.Method,
            "path", r.URL.Path,
            "status", wrapped.statusCode,
            "duracion_ms", duracion.Milliseconds(),
            "ip", r.RemoteAddr,
        )
    }
}
```

### Documentation: OpenAPI/Swagger

```go
// Comentarios que se pueden parsear para generar docs

// @Router /usuarios [get]
// @Summary Obtener lista de usuarios
// @Description Retorna paginada lista de usuarios del sistema
// @Param page query int false "Número de página"
// @Param page_size query int false "Tamaño de página"
// @Success 200 {object} []User
// @Failure 400 {object} ErrorResponse
func usuariosHandler(w http.ResponseWriter, r *http.Request) {
    // ...
}

// Herramientas como swag pueden generar documentación:
// go get -u github.com/swaggo/swag/cmd/swag
// swag init
```

### Timeouts Apropiados

```go
server := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    
    // Timeout para leer request completo
    ReadTimeout:       10 * time.Second,
    
    // Timeout para escribir response completa
    WriteTimeout:      10 * time.Second,
    
    // Timeout para idle connection
    IdleTimeout:       120 * time.Second,
    
    // Máximo de headers
    MaxHeaderBytes:    1 << 20, // 1MB
}

server.ListenAndServe()
```

### Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    
    failures int
    lastFail time.Time
    state    string // "closed", "open", "half-open"
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == "open" {
        if time.Since(cb.lastFail) > cb.timeout {
            cb.state = "half-open"
        } else {
            return errors.New("circuito abierto")
        }
    }
    
    err := fn()
    
    if err != nil {
        cb.failures++
        cb.lastFail = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }
    
    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Handler con Middleware - Logger y Timing

**Objetivo:** Crear middleware de logging que capture tiempo de ejecución y detalles del request.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// TODO: Implementar esto
// 1. Crear middleware que loguee: timestamp, método, path, IP, user-agent
// 2. Crear middleware que mida tiempo de ejecución
// 3. Encadenar ambos middlewares
// 4. Crear varios handlers de prueba
// 5. Verificar que se registren correctamente

// Puntos clave a considerar:
// - Orden de ejecución de middleware
// - Captura correcta de status code
// - Formato legible de logs
// - Cálculo preciso de duraciones
```

### Ejercicio 2: REST API - CRUD Endpoints con JSON

**Objetivo:** Crear una API RESTful completa con operaciones CRUD.

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "sync"
)

// TODO: Implementar una API de libros
// 1. Estructura Book con campos: ID, Title, Author, ISBN, Year
// 2. GET /books - Retorna lista de libros
// 3. POST /books - Crea nuevo libro (validar campos)
// 4. GET /books/:id - Obtiene libro específico
// 5. PUT /books/:id - Actualiza libro completo
// 6. DELETE /books/:id - Elimina libro
// 7. Usar sync.Map o map con mutex para almacenar en memoria
// 8. Retornar errores HTTP apropiados

// Puntos clave:
// - Validación de input
// - Status codes correctos
// - Manejo de IDs no encontrados
// - Parse/encode JSON
```

### Ejercicio 3: File Upload - Manejo de Uploads con Validación

**Objetivo:** Implementar upload de archivos con validaciones.

```go
package main

import (
    "fmt"
    "io"
    "mime"
    "net/http"
    "os"
)

// TODO: Implementar endpoint de upload
// 1. POST /upload - Acepta multipart/form-data
// 2. Validar tipo de archivo (solo PDF e imágenes)
// 3. Validar tamaño máximo (ej: 5MB)
// 4. Generar nombre único
// 5. Guardar en carpeta "uploads"
// 6. Retornar URL del archivo guardado
// 7. GET /files/:filename - Descargar archivo
// 8. DELETE /files/:filename - Eliminar archivo

// Puntos clave:
// - Validación MIME type
// - Prevención de ataques path traversal
// - Límite de tamaño
// - Manejo de errores en I/O
// - Nombre único: timestamp + hash o UUID
```

### Ejercicio 4: JWT Authentication - Autenticación basada en Tokens

**Objetivo:** Implementar autenticación JWT con login/protected endpoints.

```go
package main

import (
    "errors"
    "fmt"
    "net/http"
    "strings"
    "time"
)

// TODO: Implementar autenticación JWT
// 1. POST /login - Autentica usuario, retorna token
//    Aceptar: {username, password}
//    Retornar: {token, expires_in}
// 2. POST /refresh - Refrescar token expirado
// 3. GET /protected - Endpoint protegido con JWT
// 4. Implementar middleware que valide token
// 5. Tokens expiren en 1 hora
// 6. Refresh tokens válidos 7 días

// Puntos clave:
// - Generación de JWT
// - Validación de firma
// - Manejo de expiración
// - Refresh token flow
// - Almacenamiento seguro de secreto

// Usar: github.com/golang-jwt/jwt/v5
```

### Ejercicio 5: Graceful Shutdown - Server que Cierra Limpiamente

**Objetivo:** Implementar servidor con shutdown elegante.

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// TODO: Implementar servidor con graceful shutdown
// 1. Iniciar servidor HTTP
// 2. Rastrear requests activos
// 3. Capturar SIGINT/SIGTERM
// 4. Dejar de aceptar nuevos requests
// 5. Esperar a que se completen (máx 30 seg)
// 6. Cerrar recursos (DB, archivos)
// 7. Salir limpiamente

// Puntos clave:
// - Signal handling
// - Context con timeout
// - Connection draining
// - Orden de cleanup
// - Testing: enviar SIGTERM y verificar requests en progreso

// Simular requests largos:
// GET /slow - Duerme 5 segundos
// GET /health - Retorna OK inmediatamente

// Verificar que con timeout de 30 seg:
// - Requests de 5 seg se completan
// - Requests nuevos no se aceptan
// - Servidor se apaga ordenadamente
```

---

## Resumen de Conceptos Clave

### Arquitectura HTTP en Go

```
┌─────────────────────────────────────────────┐
│        Componentes de un Servidor HTTP       │
├─────────────────────────────────────────────┤
│                                              │
│  1. Listener (TCP)                          │
│     └─ Escucha conexiones en puerto         │
│                                              │
│  2. Connection Handler                      │
│     └─ Crea goroutine por conexión          │
│                                              │
│  3. Request Parser                          │
│     └─ Lee headers y body                   │
│                                              │
│  4. Router (ServeMux)                       │
│     └─ Mapea path a handler                 │
│                                              │
│  5. Middleware Chain                        │
│     └─ Procesa antes/después                │
│                                              │
│  6. Handler                                 │
│     └─ Lógica de negocio                    │
│                                              │
│  7. Response Writer                         │
│     └─ Serializa respuesta                  │
│                                              │
└─────────────────────────────────────────────┘
```

### Decisiones de Diseño

```
¿Usar http.ServeMux o router externo?
├─ ServeMux: Simple, built-in, suficiente para 80% de casos
└─ Externo (chi, gorilla): Más features, parámetros en path, más overhead

¿Cómo manejar middleware?
├─ Encadenamiento: Simples, transparente
├─ Inyección: Más flexible, más complejo
└─ Type wrapping: Balance entre flexibilidad y claridad

¿Cómo autenticar?
├─ API Key: Simple, sin estado
├─ JWT: Stateless, escalable
├─ Session: Con estado, más seguro para web
└─ Hybrid: JWT + refresh tokens

¿Graceful shutdown es obligatorio?
├─ Producción: SÍ, evita data loss
├─ Desarrollo: NO, no es crítico
└─ Testing: SÍ, para tests confiables
```

### Patrones de Código Go

```go
// Patrón 1: Middleware como función
func middlewareFunc(next Handler) Handler {
    return HandlerFunc(func(w ResponseWriter, r *Request) {
        // Before
        next.ServeHTTP(w, r)
        // after
    })
}

// Patrón 2: Dependency Injection
func newHandler(db Database, cache Cache) HandlerFunc {
    return func(w ResponseWriter, r *Request) {
        // Usar db y cache
    }
}

// Patrón 3: Error Wrapping
if err != nil {
    return fmt.Errorf("operación X falló: %w", err)
}

// Patrón 4: Context para valores
ctx := context.WithValue(r.Context(), "user", user)
r = r.WithContext(ctx)
```

---

## Referencias Internas del Capítulo

Este capítulo construye sobre:
- **Capítulo 8**: Interfaces y composición (http.Handler, middleware como composición)
- **Capítulo 9**: Channels y concurrencia (goroutines por request, context cancellation)
- **Capítulo 12**: Métodos y tipos (estructuras como handlers)
- **Capítulo 14**: JSON y serialización (parsing de requests, escritura de responses)
- **Capítulo 18**: Error handling (HTTP errors, recovery de panics)
- **Capítulo 28**: Context (propagación de valores y timeouts)

---

**Fin del Capítulo 38**

*Este capítulo proporciona el conocimiento profundo para construir servidores HTTP production-ready en Go, desde conceptos básicos hasta patrones avanzados de escalabilidad y seguridad.*

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/38-http-avanzado/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/38-http-avanzado):

```bash
cd examples/38-http-avanzado
go run .
```
