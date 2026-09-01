# Capítulo 30: Net - Programación de red

## Introducción

El package `net` de Go proporciona una interfaz portátil para la programación de red. Es el corazón de toda comunicación de red en Go, ofreciendo abstracciones de alto nivel para trabajar con conexiones TCP, UDP, HTTP y más. A diferencia de otros lenguajes donde la programación de red puede ser tediosa y propensa a errores, Go simplifica significativamente este proceso con su enfoque concurrente basado en goroutines.

Este capítulo explora cómo Go abstrae la complejidad de las conexiones de red, desde el protocolo TCP más bajo hasta clientes y servidores HTTP de alto nivel.

---

## 30.1 ¿Qué es el Net Package?

### 30.1.1 Conceptos Fundamentales

El `net` package es una abstracción unificada sobre los protocolos de red subyacentes del sistema operativo. Su propósito es:

1. **Portabilidad**: El mismo código funciona en Windows, Linux, macOS, etc.
2. **Simplicidad**: Interfaces limpias que ocultan la complejidad de syscalls
3. **Concurrencia**: Diseñado para goroutines, no threads
4. **Seguridad**: Manejo automático de recursos y timeouts

### 30.1.2 Arquitectura del Net Package

```
┌─────────────────────────────────────────────────────┐
│          Aplicación Go (HTTP, gRPC, etc.)          │
├─────────────────────────────────────────────────────┤
│            Net Package (Go Standard Library)        │
│  ┌──────────┬──────────┬─────────┬─────────────┐   │
│  │   TCP    │   UDP    │  HTTP   │   Lookups   │   │
│  └──────────┴──────────┴─────────┴─────────────┘   │
├─────────────────────────────────────────────────────┤
│       System Network Stack (SO - Windows/Linux)     │
├─────────────────────────────────────────────────────┤
│         Hardware (Tarjeta de red, puertos)         │
└─────────────────────────────────────────────────────┘
```

### 30.1.3 Interfaces Principales

```go
// Conn: interfaz para una conexión de red bidireccional
type Conn interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
    LocalAddr() Addr
    RemoteAddr() Addr
    SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
}

// Listener: interfaz para escuchar conexiones entrantes
type Listener interface {
    Accept() (Conn, error)
    Close() error
    Addr() Addr
}

// Addr: interfaz para direcciones de red
type Addr interface {
    Network() string  // "tcp", "udp", "ip", etc.
    String() string   // representación en string
}
```

### 30.1.4 Protocolos Soportados

Go soporta varios protocolos a través del net package:

| Protocolo | Tipo | Ejemplo | Uso |
|-----------|------|---------|-----|
| tcp | Conexión | "tcp" | Aplicaciones cliente-servidor confiables |
| tcp4 | Conexión | "tcp4" | IPv4 específicamente |
| tcp6 | Conexión | "tcp6" | IPv6 específicamente |
| udp | Sin conexión | "udp" | Aplicaciones de baja latencia |
| udp4 | Sin conexión | "udp4" | UDP sobre IPv4 |
| udp6 | Sin conexión | "udp6" | UDP sobre IPv6 |
| unix | Conexión | "unix" | IPC local (no en Windows) |
| unixgram | Sin conexión | "unixgram" | IPC local sin conexión |

### 30.1.5 Comparación: Go vs Otros Lenguajes

```go
// ╔═══════════════════════════════════════════════════════════╗
// ║              Go (Simple y Concurrente)                   ║
// ╚═══════════════════════════════════════════════════════════╝
conn, err := net.Dial("tcp", "example.com:80")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
buf := make([]byte, 1024)
n, _ := conn.Read(buf)
fmt.Println(string(buf[:n]))

// ╔═══════════════════════════════════════════════════════════╗
// ║         Java (Verbose, Basado en Threads)               ║
// ╚═══════════════════════════════════════════════════════════╝
/*
Socket socket = new Socket("example.com", 80);
InputStream in = socket.getInputStream();
OutputStream out = socket.getOutputStream();
out.write("GET / HTTP/1.1\r\n\r\n".getBytes());
byte[] buf = new byte[1024];
int n = in.read(buf);
System.out.println(new String(buf, 0, n));
socket.close();
*/

// ╔═══════════════════════════════════════════════════════════╗
// ║       Python (Simple pero con GIL - Global Lock)        ║
// ╚═══════════════════════════════════════════════════════════╝
/*
import socket
s = socket.socket()
s.connect(("example.com", 80))
s.send(b"GET / HTTP/1.1\r\n\r\n")
print(s.recv(1024).decode())
s.close()
*/
```

---

## 30.2 Conexiones TCP

### 30.2.1 Fundamentos de TCP

TCP (Transmission Control Protocol) es un protocolo orientado a conexión que garantiza:

- **Entrega garantizada**: Los datos llegan en orden y completos
- **Control de flujo**: Ambos extremos controlan la velocidad
- **Reordenamiento**: Los paquetes fuera de orden se reordenan automáticamente
- **Detección de errores**: Los datos dañados se retransmiten

### 30.2.2 Handshake TCP

```
┌─────────────┐                           ┌─────────────┐
│   Cliente   │                           │  Servidor   │
└─────────────┘                           └─────────────┘
      │                                           │
      │  SYN (SEQ=x)                             │
      ├──────────────────────────────────────────>│
      │                                           │
      │              SYN-ACK (SEQ=y, ACK=x+1)    │
      │<──────────────────────────────────────────┤
      │                                           │
      │  ACK (SEQ=x+1, ACK=y+1)                  │
      ├──────────────────────────────────────────>│
      │                                           │
      │          [Conexión establecida]          │
      ├◄─────── Datos pueden fluir ─────────────>│
```

### 30.2.3 Cliente TCP Básico

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // Conectarse a un servidor
    conn, err := net.Dial("tcp", "google.com:80")
    if err != nil {
        log.Fatalf("Error conectando: %v", err)
    }
    defer conn.Close()

    // Escribir datos
    fmt.Fprintf(conn, "GET / HTTP/1.0\r\nHost: google.com\r\n\r\n")

    // Leer respuesta
    buf := make([]byte, 512)
    n, err := conn.Read(buf)
    if err != nil {
        log.Fatalf("Error leyendo: %v", err)
    }

    fmt.Println("Respuesta recibida:")
    fmt.Println(string(buf[:n]))
}
```

### 30.2.4 Servidor TCP Básico (Echo Server)

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    // Copiar todo lo que llega de vuelta al cliente (echo)
    fmt.Printf("Cliente conectado: %s\n", conn.RemoteAddr())
    io.Copy(conn, conn)
}

func main() {
    // Escuchar en puerto 8080
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("Error escuchando: %v", err)
    }
    defer listener.Close()

    fmt.Println("Servidor escuchando en :8080")

    for {
        // Aceptar conexión entrante
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Error aceptando: %v", err)
            continue
        }

        // Manejar en goroutine
        go handleConnection(conn)
    }
}
```

### 30.2.5 Servidor TCP Avanzado con Logging

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
)

func handleClient(conn net.Conn) {
    defer conn.Close()

    localAddr := conn.LocalAddr().String()
    remoteAddr := conn.RemoteAddr().String()

    log.Printf("Nueva conexión de %s a %s", remoteAddr, localAddr)

    scanner := bufio.NewScanner(conn)

    for scanner.Scan() {
        line := scanner.Text()

        // Procesar comandos
        switch {
        case strings.HasPrefix(line, "ECHO "):
            msg := strings.TrimPrefix(line, "ECHO ")
            fmt.Fprintf(conn, "ECHO: %s\n", msg)

        case strings.HasPrefix(line, "TIME"):
            fmt.Fprintf(conn, "Hora del servidor\n")

        case strings.HasPrefix(line, "QUIT"):
            fmt.Fprintf(conn, "Adiós!\n")
            return

        default:
            fmt.Fprintf(conn, "Comando desconocido: %s\n", line)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Printf("Error scanner: %v", err)
    }
}

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    log.Println("Servidor TCP en :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Error accept: %v", err)
            continue
        }
        go handleClient(conn)
    }
}
```

### 30.2.6 Información de Conexión

```go
func inspectConnection(conn net.Conn) {
    // Dirección local
    fmt.Printf("Dirección local: %s\n", conn.LocalAddr())
    fmt.Printf("Tipo local: %s\n", conn.LocalAddr().Network())

    // Dirección remota
    fmt.Printf("Dirección remota: %s\n", conn.RemoteAddr())
    fmt.Printf("Tipo remoto: %s\n", conn.RemoteAddr().Network())
}
```

---

## 30.3 UDP - Protocolo sin Conexión

### 30.3.1 Fundamentos de UDP

UDP (User Datagram Protocol) es un protocolo sin conexión donde:

- **Sin garantía**: Los datagramas pueden perderse
- **Sin ordenamiento**: Los paquetes pueden llegar fuera de orden
- **Bajo overhead**: Mucho más rápido que TCP
- **Mejor para**: Streaming, gaming, DNS, VoIP

### 30.3.2 Cliente UDP

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // Resolver dirección UDP
    addr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
    if err != nil {
        log.Fatal(err)
    }

    // Conectarse (optativo en UDP)
    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Enviar datos
    msg := []byte("¿Quién soy?")
    _, err = conn.Write(msg)
    if err != nil {
        log.Fatal(err)
    }

    // Recibir respuesta
    buf := make([]byte, 1024)
    n, remoteAddr, err := conn.ReadFromUDP(buf)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Recibido de %s: %s\n", remoteAddr, string(buf[:n]))
}
```

### 30.3.3 Servidor UDP

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // Crear listener UDP
    addr, err := net.ResolveUDPAddr("udp", ":8888")
    if err != nil {
        log.Fatal(err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    fmt.Println("Servidor UDP escuchando en :8888")

    buf := make([]byte, 1024)

    for {
        // Leer datagrama
        n, remoteAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Printf("Error leyendo: %v", err)
            continue
        }

        fmt.Printf("De %s: %s\n", remoteAddr, string(buf[:n]))

        // Enviar respuesta
        response := []byte("Respuesta UDP")
        conn.WriteToUDP(response, remoteAddr)
    }
}
```

### 30.3.4 PacketConn - Interfaz Genérica UDP

```go
// PacketConn es la interfaz para conexiones sin conexión
type PacketConn interface {
    ReadFrom(b []byte) (n int, addr Addr, err error)
    WriteTo(b []byte, addr Addr) (n int, err error)
    Close() error
    LocalAddr() Addr
    SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
}
```

### 30.3.5 Multicast UDP

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // Dirección multicast (224.0.0.0 - 239.255.255.255)
    addr, err := net.ResolveUDPAddr("udp", "224.0.0.1:5353")
    if err != nil {
        log.Fatal(err)
    }

    // Escuchar multicast
    conn, err := net.ListenMulticastUDP("udp", nil, addr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    buf := make([]byte, 1024)

    for {
        n, srcAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Multicast de %s: %s\n", srcAddr, string(buf[:n]))
    }
}
```

---

## 30.4 HTTP Client

### 30.4.1 Fundamentos de HTTP

HTTP (HyperText Transfer Protocol) es un protocolo de aplicación que:

- Se construye sobre TCP
- Usa modelo cliente-servidor
- Es stateless (cada request es independiente)
- Soporta múltiples métodos: GET, POST, PUT, DELETE, etc.

### 30.4.2 HTTP GET Básico

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
)

func main() {
    // GET simple
    resp, err := http.Get("https://www.example.com")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.Status)
    fmt.Println("Status Code:", resp.StatusCode)

    // Leer body completo
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Body:", string(body[:100]))
}
```

### 30.4.3 HTTP POST con Datos

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "net/http"
)

func main() {
    // POST con JSON
    jsonData := []byte(`{"nombre":"Juan","edad":30}`)

    resp, err := http.Post(
        "https://httpbin.org/post",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### 30.4.4 Requests Personalizados

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
)

func main() {
    // Crear request personalizado
    req, err := http.NewRequest("GET", "https://api.github.com/users/golang", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Agregar headers
    req.Header.Add("Accept", "application/json")
    req.Header.Add("User-Agent", "Go-Client/1.0")

    // Enviar request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)

    // Ver headers de respuesta
    for k, v := range resp.Header {
        fmt.Printf("%s: %s\n", k, v)
    }

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### 30.4.5 Cliente HTTP Personalizado con Timeout

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func main() {
    // Cliente con timeout
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    req, _ := http.NewRequest("GET", "https://httpbin.org/delay/2", nil)

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)
}
```

### 30.4.6 Reutilizar Conexiones

```go
package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "time"
)

func main() {
    // Cliente con connection pooling
    client := &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:       100,
            MaxIdleConnsPerHost: 100,
            IdleConnTimeout:    90 * time.Second,
            DialContext: (&net.Dialer{
                Timeout:   30 * time.Second,
                KeepAlive: 30 * time.Second,
            }).DialContext,
        },
    }

    // Las siguientes requests reutilizarán conexiones
    for i := 0; i < 5; i++ {
        resp, err := client.Get("https://httpbin.org/get")
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Request %d: %d\n", i+1, resp.StatusCode)
        resp.Body.Close()
    }
}
```

### 30.4.7 Manejo de Errores HTTP

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "net/url"
    "time"
)

func main() {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Get("https://invalid-domain-xyz.com")
    if err != nil {
        // Analizar tipo de error
        switch e := err.(type) {
        case net.Error:
            if e.Timeout() {
                fmt.Println("Timeout en conexión")
            } else {
                fmt.Println("Error de red:", e)
            }
        case *url.Error:
            fmt.Println("Error de URL:", e)
        default:
            fmt.Println("Error desconocido:", err)
        }
        return
    }
    defer resp.Body.Close()

    // Verificar status code
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        fmt.Printf("Error HTTP %d: %s\n", resp.StatusCode, string(body))
        return
    }
}
```

---

## 30.5 HTTP Server

### 30.5.1 Servidor HTTP Minimal

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    // Handler simple
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hola, %s!", r.URL.Path[1:])
    })

    log.Println("Servidor en :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 30.5.2 Múltiples Routes

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Bienvenido a Home"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "mensaje": "Esto es un API",
        "método":  r.Method,
    })
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/api", apiHandler)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 30.5.3 ServeMux - Router Personalizado

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()

    // Rutas exactas
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Inicio")
    })

    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Lista de usuarios")
    })

    // Rutas con prefijo (terminan con /)
    mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "API: %s", r.URL.Path)
    })

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### 30.5.4 Middleware

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// Middleware para logging
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Registrar antes
        log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

        // Llamar al siguiente handler
        next.ServeHTTP(w, r)

        // Registrar después
        log.Printf("Completado en %v\n", time.Since(start))
    })
}

// Middleware para autenticación
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Token requerido", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hola!")
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/hello", authMiddleware(http.HandlerFunc(helloHandler)))

    handler := loggingMiddleware(mux)

    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

### 30.5.5 Servidor con Métodos HTTP Específicos

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type Usuario struct {
    ID   int    `json:"id"`
    Nombre string `json:"nombre"`
}

func usuariosHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case http.MethodGet:
        usuarios := []Usuario{
            {1, "Juan"},
            {2, "María"},
        }
        json.NewEncoder(w).Encode(usuarios)

    case http.MethodPost:
        var u Usuario
        json.NewDecoder(r.Body).Decode(&u)
        fmt.Fprintf(w, `{"mensaje":"Usuario %s creado"}`, u.Nombre)

    case http.MethodDelete:
        fmt.Fprint(w, `{"mensaje":"Usuario eliminado"}`)

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func main() {
    http.HandleFunc("/usuarios", usuariosHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 30.5.6 Servidor con Context y Cancellación

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
    server := &http.Server{
        Addr: ":8080",
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Context con timeout
            ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
            defer cancel()

            select {
            case <-time.After(1 * time.Second):
                fmt.Fprint(w, "Respuesta después de 1s")
            case <-ctx.Done():
                http.Error(w, "Timeout", http.StatusRequestTimeout)
            }
        }),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
    }

    // Gratuito goroutine para el servidor
    go func() {
        log.Fatal(server.ListenAndServe())
    }()

    // Esperar señal de terminación
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(err)
    }
}
```

---

## 30.6 URL Parsing

### 30.6.1 Estructura de una URL

```
https://user:pass@example.com:8080/path?query=value&foo=bar#fragment

┌──────────────────────────────────────────────────────────────────┐
│                          URL Completa                            │
├──────────────────────────────────────────────────────────────────┤
│ https  = Scheme                                                  │
│ user:pass  = Userinfo (usuario y contraseña)                    │
│ example.com  = Host (hostname)                                  │
│ 8080  = Port                                                    │
│ /path  = Path                                                   │
│ query=value&foo=bar  = RawQuery (query string)                 │
│ fragment  = Fragment                                            │
└──────────────────────────────────────────────────────────────────┘
```

### 30.6.2 Parsear URLs

```go
package main

import (
    "fmt"
    "log"
    "net/url"
)

func main() {
    urlStr := "https://user:pass@example.com:8080/path/to/resource?key=value&foo=bar#section"

    // Parsear URL
    u, err := url.Parse(urlStr)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Scheme:", u.Scheme)
    fmt.Println("User:", u.User.Username())
    fmt.Println("Password:", u.User.String())
    fmt.Println("Host:", u.Host)
    fmt.Println("Hostname:", u.Hostname())
    fmt.Println("Port:", u.Port())
    fmt.Println("Path:", u.Path)
    fmt.Println("RawQuery:", u.RawQuery)
    fmt.Println("Fragment:", u.Fragment)
}
```

### 30.6.3 Query Parameters

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    u, _ := url.Parse("https://example.com?key=value&foo=bar&foo=baz")

    // Parsear query
    query := u.Query()

    fmt.Println("key:", query.Get("key"))
    fmt.Println("foo (singular):", query.Get("foo"))
    fmt.Println("foo (todos):", query["foo"])

    // Iterar
    for key, values := range query {
        for _, value := range values {
            fmt.Printf("%s = %s\n", key, value)
        }
    }
}
```

### 30.6.4 Construir URLs Dinámicamente

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    // Crear valores de query
    params := url.Values{}
    params.Add("nombre", "Juan")
    params.Add("edad", "30")
    params.Add("ciudad", "Madrid")
    params.Add("hobbies", "leer")
    params.Add("hobbies", "programar")

    // Construir URL
    u := &url.URL{
        Scheme:   "https",
        Host:     "api.example.com",
        Path:     "/usuarios",
        RawQuery: params.Encode(),
    }

    fmt.Println("URL:", u.String())
    // Output: https://api.example.com/usuarios?nombre=Juan&edad=30&ciudad=Madrid&hobbies=leer&hobbies=programar
}
```

### 30.6.5 Encoding y Escaping

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    // URL encoding
    text := "¡Hola Mundo! ¿Cómo estás?"

    encoded := url.QueryEscape(text)
    fmt.Println("Encoded:", encoded)
    // Output: %C2%A1Hola+Mundo%21+%C2%BFc%C3%B3mo+est%C3%A1s%3F

    // Path encoding
    pathSegment := "mi/ruta especial"
    pathEncoded := url.PathEscape(pathSegment)
    fmt.Println("Path Encoded:", pathEncoded)
    // Output: mi%2Fruta%20especial

    // Decodificar
    decoded, _ := url.QueryUnescape(encoded)
    fmt.Println("Decoded:", decoded)
}
```

### 30.6.6 Casos de Uso Práctico

```go
package main

import (
    "fmt"
    "net/url"
)

func construirAPIRequest(base string, params map[string]string) string {
    u, _ := url.Parse(base)
    q := u.Query()

    for k, v := range params {
        q.Set(k, v)
    }

    u.RawQuery = q.Encode()
    return u.String()
}

func main() {
    // Construir URL de API con parámetros dinámicos
    apiURL := construirAPIRequest(
        "https://api.github.com/search/repositories",
        map[string]string{
            "q":     "language:go stars:>1000",
            "sort":  "stars",
            "order": "desc",
        },
    )

    fmt.Println(apiURL)
}
```

---

## 30.7 IP Addresses

### 30.7.1 Tipos de Direcciones IP

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Parsear IPv4
    ipv4 := net.ParseIP("192.168.1.1")
    fmt.Printf("IPv4: %s (tipo: %T)\n", ipv4, ipv4)

    // Parsear IPv6
    ipv6 := net.ParseIP("2001:db8::8a2e:370:7334")
    fmt.Printf("IPv6: %s\n", ipv6)

    // Validación
    if ipv4 == nil {
        fmt.Println("IP inválida")
    } else {
        fmt.Println("IP válida")
    }
}
```

### 30.7.2 Operaciones con IPs

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    ip := net.ParseIP("192.168.1.100")

    // Convertir a 4 bytes
    bytes4 := ip.To4()
    fmt.Printf("IPv4 bytes: %v\n", bytes4)

    // Convertir a 16 bytes
    bytes16 := ip.To16()
    fmt.Printf("IPv6 bytes: %v\n", bytes16)

    // Tipos de direcciones especiales
    fmt.Printf("IsLoopback: %v\n", ip.IsLoopback())           // false
    fmt.Printf("IsPrivate: %v\n", ip.IsPrivate())             // true
    fmt.Printf("IsGlobalUnicast: %v\n", ip.IsGlobalUnicast()) // false

    localhost := net.ParseIP("127.0.0.1")
    fmt.Printf("Localhost IsLoopback: %v\n", localhost.IsLoopback()) // true
}
```

### 30.7.3 CIDR - Notación de Redes

```
┌─────────────────────────────────────────────────┐
│    CIDR: Classless Inter-Domain Routing        │
├─────────────────────────────────────────────────┤
│ 192.168.1.0/24                                 │
│ ├─ 192.168.1.0 = Dirección de red              │
│ ├─ /24 = Máscara (primeros 24 bits)            │
│ └─ Hosts disponibles: 192.168.1.1 - 254       │
│                                                 │
│ 10.0.0.0/16                                   │
│ ├─ Hosts disponibles: 10.0.0.1 - 65534       │
│                                                 │
│ 172.16.0.0/12                                 │
│ ├─ Hosts disponibles: 172.16.0.1 - 172.31.255.254 │
└─────────────────────────────────────────────────┘
```

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // Parsear CIDR
    _, ipnet, err := net.ParseCIDR("192.168.1.0/24")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Red: %s\n", ipnet)
    fmt.Printf("Máscara: %s\n", ipnet.Mask)
    fmt.Printf("IP: %s\n", ipnet.IP)

    // Verificar si una IP está en la red
    testIP := net.ParseIP("192.168.1.100")
    if ipnet.Contains(testIP) {
        fmt.Println("192.168.1.100 está en la red")
    }

    testIP2 := net.ParseIP("10.0.0.1")
    if !ipnet.Contains(testIP2) {
        fmt.Println("10.0.0.1 NO está en la red")
    }
}
```

### 30.7.4 Máscaras de Red

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Crear máscara
    mask := net.CIDRMask(24, 32)
    fmt.Printf("Máscara /24: %s\n", mask)

    // Obtener información de la máscara
    ones, bits := mask.Size()
    fmt.Printf("Ones: %d, Total bits: %d\n", ones, bits)

    // Aplicar máscara a una IP
    ip := net.ParseIP("192.168.1.100")
    network := ip.Mask(mask)
    fmt.Printf("Red: %s\n", network)
}
```

### 30.7.5 Broadcast y Direcciones Especiales

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Dirección de loopback
    loopback := net.IPv4(127, 0, 0, 1)
    fmt.Printf("Loopback: %s\n", loopback)

    // Dirección 0.0.0.0 (cualquier interfaz)
    any := net.IPv4(0, 0, 0, 0)
    fmt.Printf("Any: %s\n", any)

    // Broadcast
    broadcast := net.IPv4(255, 255, 255, 255)
    fmt.Printf("Broadcast: %s\n", broadcast)

    // Verificar tipos
    fmt.Printf("Loopback IsLoopback: %v\n", loopback.IsLoopback())
    fmt.Printf("Any IsUnspecified: %v\n", any.IsUnspecified())
}
```

---

## 30.8 DNS Resolution

### 30.8.1 Conceptos de DNS

DNS (Domain Name System) traduce nombres de dominio en direcciones IP:

```
┌──────────────────────────────────────────────────┐
│  Usuario escribe: www.example.com en navegador  │
├──────────────────────────────────────────────────┤
│  Cliente DNS solicita a servidor DNS recursivo  │
├──────────────────────────────────────────────────┤
│  Servidor DNS pregunta a root nameserver        │
├──────────────────────────────────────────────────┤
│  Pregunta a TLD nameserver (.com)               │
├──────────────────────────────────────────────────┤
│  Pregunta a nameserver autoritativo             │
├──────────────────────────────────────────────────┤
│  Respuesta: 93.184.216.34                       │
└──────────────────────────────────────────────────┘
```

### 30.8.2 Lookup Básico

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // LookupHost - obtener IPs de un host
    addrs, err := net.LookupHost("golang.org")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Direcciones de golang.org:")
    for _, addr := range addrs {
        fmt.Println(" -", addr)
    }
}
```

### 30.8.3 Lookup de IPs

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // LookupIP - obtener IPs (más eficiente)
    ips, err := net.LookupIP("google.com")
    if err != nil {
        log.Fatal(err)
    }

    for _, ip := range ips {
        fmt.Printf("%s (%s)\n", ip, ip.String())
    }

    // Lookup de dirección inversa (IP a hostname)
    names, err := net.LookupAddr("8.8.8.8")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Nombres para 8.8.8.8:")
    for _, name := range names {
        fmt.Println(" -", name)
    }
}
```

### 30.8.4 Lookup de Registros DNS

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // MX records
    mxRecords, err := net.LookupMX("gmail.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Servidores de correo de gmail.com:")
    for _, mx := range mxRecords {
        fmt.Printf(" - %s (prioridad: %d)\n", mx.Host, mx.Preference)
    }

    // SRV records
    _, srvRecords, err := net.LookupSRV("http", "tcp", "example.com")
    if err != nil {
        // Puede no existir
        fmt.Println("No hay registros SRV")
    } else {
        for _, srv := range srvRecords {
            fmt.Printf(" - %s:%d\n", srv.Target, srv.Port)
        }
    }

    // NS records
    nsRecords, err := net.LookupNS("example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Nameservers de example.com:")
    for _, ns := range nsRecords {
        fmt.Println(" -", ns.Host)
    }

    // TXT records
    txtRecords, err := net.LookupTXT("example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Registros TXT de example.com:")
    for _, txt := range txtRecords {
        fmt.Println(" -", txt)
    }
}
```

### 30.8.5 DNS con Context y Timeout

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    // Resolver con timeout
    resolver := &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            dialer := &net.Dialer{
                Timeout: 3 * time.Second,
            }
            return dialer.DialContext(ctx, network, address)
        },
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Lookup con contexto
    addrs, err := resolver.LookupHost(ctx, "golang.org")
    if err != nil {
        log.Printf("Error en lookup: %v", err)
        return
    }

    fmt.Println("Direcciones:")
    for _, addr := range addrs {
        fmt.Println(" -", addr)
    }
}
```

### 30.8.6 Caché DNS Local

```go
package main

import (
    "fmt"
    "log"
    "net"
    "sync"
    "time"
)

type DNSCache struct {
    cache map[string][]net.IP
    mu    sync.RWMutex
}

func (dc *DNSCache) Lookup(hostname string) ([]net.IP, error) {
    dc.mu.RLock()
    if ips, ok := dc.cache[hostname]; ok {
        defer dc.mu.RUnlock()
        return ips, nil
    }
    dc.mu.RUnlock()

    // No en caché, hacer lookup
    ips, err := net.LookupIP(hostname)
    if err != nil {
        return nil, err
    }

    // Guardar en caché
    dc.mu.Lock()
    dc.cache[hostname] = ips
    dc.mu.Unlock()

    return ips, nil
}

func main() {
    cache := &DNSCache{cache: make(map[string][]net.IP)}

    ips1, _ := cache.Lookup("golang.org")
    fmt.Println("Primera búsqueda:", ips1)

    ips2, _ := cache.Lookup("golang.org")
    fmt.Println("Desde caché:", ips2)
}
```

---

## 30.9 Timeouts y Deadlines

### 30.9.1 Conceptos

Timeouts son críticos en programación de red. Sin ellos, un programa puede quedar congelado indefinidamente esperando una respuesta.

**Tipos de timeouts:**

- **Dial timeout**: Tiempo para conectarse
- **Read timeout**: Tiempo para leer datos
- **Write timeout**: Tiempo para escribir datos
- **Overall timeout**: Tiempo total de operación

### 30.9.2 Timeouts en TCP

```go
package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    // Dial con timeout
    conn, err := net.DialTimeout("tcp", "example.com:80", 5*time.Second)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Set deadline para lectura (5 segundos desde ahora)
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))

    // Set deadline para escritura
    conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

    // Operación que respeta el deadline
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
            fmt.Println("Lectura expiró")
        }
    } else {
        fmt.Printf("Leído %d bytes\n", n)
    }
}
```

### 30.9.3 Context Timeouts

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"
)

func dialWithContext(host string, port string) (net.Conn, error) {
    ctx, cancel := context.WithTimeout(
        context.Background(),
        5*time.Second,
    )
    defer cancel()

    dialer := &net.Dialer{
        Timeout: 5 * time.Second,
    }

    return dialer.DialContext(ctx, "tcp", net.JoinHostPort(host, port))
}

func main() {
    conn, err := dialWithContext("example.com", "80")
    if err != nil {
        if err == context.DeadlineExceeded {
            fmt.Println("Timeout en conexión")
        } else {
            log.Fatal(err)
        }
    }
    defer conn.Close()
}
```

### 30.9.4 HTTP Timeouts

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func main() {
    // Cliente con múltiples timeouts
    client := &http.Client{
        Timeout: 10 * time.Second, // Timeout total
    }

    // También se puede configurar más finamente
    tr := &http.Transport{
        DialTimeout:         5 * time.Second,
        TLSHandshakeTimeout: 5 * time.Second,
        ResponseHeaderTimeout: 5 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }
    client.Transport = tr

    // Hacer request
    resp, err := client.Get("https://httpbin.org/delay/3")
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)
}
```

### 30.9.5 Patrón: Retry con Backoff

```go
package main

import (
    "fmt"
    "log"
    "math"
    "math/rand"
    "net/http"
    "time"
)

func retryWithBackoff(maxRetries int, initialBackoff time.Duration) error {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    backoff := initialBackoff

    for attempt := 0; attempt < maxRetries; attempt++ {
        resp, err := client.Get("https://httpbin.org/status/500")
        if err == nil {
            resp.Body.Close()

            if resp.StatusCode < 500 {
                fmt.Printf("Éxito en intento %d\n", attempt+1)
                return nil
            }
        }

        fmt.Printf("Intento %d falló, esperando %v\n", attempt+1, backoff)

        time.Sleep(backoff)

        // Exponential backoff con jitter
        backoff = time.Duration(
            float64(backoff) * 1.5 +
                time.Duration(rand.Intn(100))*time.Millisecond,
        )
    }

    return fmt.Errorf("todas los intentos fallaron")
}

func main() {
    err := retryWithBackoff(3, 100*time.Millisecond)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 30.9.6 Connection Pooling y Reuso

```go
package main

import (
    "fmt"
    "net"
    "net/http"
    "time"
)

func main() {
    // Configurar transporte con pooling
    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,

        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
        IdleConnTimeout:     90 * time.Second,

        // Compresión
        DisableCompression: false,

        // Keep-alive en conexiones
        DisableKeepAlives: false,
    }

    client := &http.Client{
        Transport: transport,
        Timeout:   10 * time.Second,
    }

    // Hacer múltiples requests - reutilizará conexiones
    for i := 0; i < 5; i++ {
        resp, _ := client.Get("https://api.github.com/repos/golang/go")
        fmt.Printf("Request %d: %d\n", i+1, resp.StatusCode)
        resp.Body.Close()
    }
}
```

---

## 30.10 TLS/SSL

### 30.10.1 Conceptos de TLS

TLS (Transport Layer Security) proporciona:

- **Encriptación**: Datos no legibles en tránsito
- **Autenticación**: Verificar identidad del servidor
- **Integridad**: Detectar si datos fueron alterados

```
┌─────────────────────────────────────────────────┐
│           TLS Handshake (HTTPS)                │
├─────────────────────────────────────────────────┤
│ 1. Cliente hola -> Servidor                    │
│    ├─ Versión TLS                              │
│    ├─ Cipher suites soportados                 │
│    └─ Números aleatorios                       │
│                                                 │
│ 2. Servidor hola <- Cliente                    │
│    ├─ Certificado                              │
│    ├─ Cipher suite elegido                     │
│    └─ Parámetros Diffie-Hellman                │
│                                                 │
│ 3. Intercambio de claves criptográficas       │
│                                                 │
│ 4. Cambio a encriptación                      │
│    ├─ Todos los datos siguientes están        │
│    │  encriptados                              │
│    └─ Verificación de integridad               │
└─────────────────────────────────────────────────┘
```

### 30.10.2 HTTPS Client Básico

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
)

func main() {
    // HTTPS automáticamente verifica certificados
    resp, err := http.Get("https://www.google.com")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Status: %d, Bytes: %d\n", resp.StatusCode, len(body))
}
```

### 30.10.3 Ignorar Verificación de Certificados (INSEGURO)

```go
package main

import (
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "net/http"
)

func main() {
    // ⚠️ SOLO PARA DESARROLLO/TESTING ⚠️
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    }

    client := &http.Client{Transport: tr}
    resp, err := client.Get("https://expired.badssl.com/")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)
}
```

### 30.10.4 Servidor TLS

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Conexión HTTPS segura!")
    })

    // ListenAndServeTLS requiere archivos de certificado
    log.Println("HTTPS Server en :8443")
    log.Fatal(http.ListenAndServeTLS(
        ":8443",
        "cert.pem",      // Certificado
        "key.pem",       // Clave privada
        mux,
    ))
}
```

### 30.10.5 Generar Certificados Autofirmados

```bash
# Generar clave privada
openssl genrsa -out key.pem 2048

# Generar certificado autofirmado (válido 365 días)
openssl req -new -x509 -key key.pem -out cert.pem -days 365 \
    -subj "/C=ES/ST=Madrid/L=Madrid/O=MyOrg/CN=localhost"
```

### 30.10.6 Conexión TLS Personalizada

```go
package main

import (
    "crypto/tls"
    "fmt"
    "log"
)

func main() {
    // Configuración TLS personalizada
    tlsConfig := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        MaxVersion:               tls.VersionTLS13,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_AES_256_GCM_SHA384,
            tls.TLS_CHACHA20_POLY1305_SHA256,
            tls.TLS_AES_128_GCM_SHA256,
        },
    }

    // Conectar
    conn, err := tls.Dial("tcp", "www.google.com:443", tlsConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Información del certificado
    cert := conn.ConnectionState().PeerCertificates[0]
    fmt.Printf("Certificado CN: %s\n", cert.Subject.CommonName)
    fmt.Printf("Válido desde: %s\n", cert.NotBefore)
    fmt.Printf("Válido hasta: %s\n", cert.NotAfter)
    fmt.Printf("Versión TLS: %v\n", conn.ConnectionState().Version)
}
```

---

## 30.11 Buenas Prácticas y Patrones

### 30.11.1 Error Handling en Red

```go
package main

import (
    "fmt"
    "net"
    "net/http"
    "os"
    "syscall"
    "time"
)

func handleNetError(err error) {
    if err == nil {
        return
    }

    switch e := err.(type) {
    case net.Error:
        if e.Timeout() {
            fmt.Println("Timeout en operación de red")
        } else if e.Temporary() {
            fmt.Println("Error temporal, reintentar")
        } else {
            fmt.Println("Error permanente:", e)
        }

    case *net.OpError:
        fmt.Printf("Op: %s, Net: %s, Err: %v\n", e.Op, e.Net, e.Err)

    case syscall.Errno:
        fmt.Println("Error del SO:", e)

    default:
        fmt.Println("Error desconocido:", err)
    }
}

func main() {
    // Ejemplo: Error de conexión
    _, err := net.Dial("tcp", "invalid-host:80")
    handleNetError(err)
}
```

### 30.11.2 Patrón: Graceful Shutdown

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
    server := &http.Server{
        Addr: ":8080",
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprint(w, "Servidor corriendo")
        }),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Iniciar servidor en goroutine
    go func() {
        log.Println("Iniciando servidor en :8080")
        if err := server.ListenAndServe(); err != nil &&
            err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Esperar señal de terminación
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    log.Println("Recibida señal de terminación, shutting down...")

    // Graceful shutdown con timeout
    ctx, cancel := context.WithTimeout(
        context.Background(),
        30*time.Second,
    )
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Shutdown error: %v", err)
    }

    log.Println("Servidor detenido")
}
```

### 30.11.3 Patrón: Health Check

```go
package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "time"
)

type HealthChecker struct {
    dbConn net.Conn
    redisConn net.Conn
}

func (h *HealthChecker) checkDatabase() bool {
    // Verificar conexión a base de datos
    if h.dbConn == nil {
        return false
    }

    // Ping a la base de datos
    h.dbConn.SetDeadline(time.Now().Add(2 * time.Second))
    _, err := h.dbConn.Write([]byte("PING"))
    return err == nil
}

func (h *HealthChecker) handler(w http.ResponseWriter, r *http.Request) {
    status := map[string]bool{
        "database": h.checkDatabase(),
    }

    allHealthy := status["database"]

    if allHealthy {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"healthy"}`)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        fmt.Fprint(w, `{"status":"unhealthy"}`)
    }
}

func main() {
    hc := &HealthChecker{}
    http.HandleFunc("/health", hc.handler)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 30.11.4 Patrón: Rate Limiting

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
)

type RateLimiter struct {
    requestsPerSecond int
    mu                sync.Mutex
    lastRequest       time.Time
    tokenBucket       float64
}

func NewRateLimiter(rps int) *RateLimiter {
    return &RateLimiter{
        requestsPerSecond: rps,
        lastRequest:       time.Now(),
        tokenBucket:       float64(rps),
    }
}

func (rl *RateLimiter) Allow() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(rl.lastRequest).Seconds()
    rl.lastRequest = now

    // Agregar tokens
    rl.tokenBucket += elapsed * float64(rl.requestsPerSecond)
    if rl.tokenBucket > float64(rl.requestsPerSecond) {
        rl.tokenBucket = float64(rl.requestsPerSecond)
    }

    if rl.tokenBucket >= 1.0 {
        rl.tokenBucket -= 1.0
        return true
    }

    return false
}

func rateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            w.WriteHeader(http.StatusTooManyRequests)
            fmt.Fprint(w, "Rate limit exceeded")
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    limiter := NewRateLimiter(10) // 10 requests/segundo

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "OK")
    })

    handler := rateLimitMiddleware(limiter, mux)
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

### 30.11.5 Patrón: Circuit Breaker

```go
package main

import (
    "fmt"
    "log"
    "sync"
    "time"
)

type CircuitBreaker struct {
    mu           sync.RWMutex
    maxFailures  int
    timeout      time.Duration
    failures     int
    lastFailTime time.Time
    state        string // "closed", "open", "half-open"
}

const (
    StateClosed   = "closed"
    StateOpen     = "open"
    StateHalfOpen = "half-open"
)

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()

    if cb.state == StateOpen {
        if time.Since(cb.lastFailTime) > cb.timeout {
            cb.state = StateHalfOpen
            cb.failures = 0
        } else {
            cb.mu.Unlock()
            return fmt.Errorf("circuit breaker abierto")
        }
    }

    cb.mu.Unlock()

    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
        }

        return err
    }

    // Éxito
    cb.failures = 0
    cb.state = StateClosed

    return nil
}

func main() {
    cb := &CircuitBreaker{
        maxFailures: 3,
        timeout:     5 * time.Second,
        state:       StateClosed,
    }

    // Usar circuit breaker
    for i := 0; i < 10; i++ {
        err := cb.Call(func() error {
            if i < 5 {
                return fmt.Errorf("fallo simulado")
            }
            return nil
        })

        if err != nil {
            log.Printf("Intento %d: %v", i+1, err)
        } else {
            log.Printf("Intento %d: éxito", i+1)
        }

        time.Sleep(100 * time.Millisecond)
    }
}
```

### 30.11.6 Testing de Código de Red

```go
package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAPIHandler(t *testing.T) {
    // Crear servidor de test
    server := httptest.NewServer(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprint(w, "respuesta de test")
        }),
    )
    defer server.Close()

    // Hacer request al servidor de test
    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("esperaba 200, obtuve %d", resp.StatusCode)
    }
}

func TestTCPServer(t *testing.T) {
    // Crear listener de test
    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        t.Fatal(err)
    }
    defer listener.Close()

    go func() {
        conn, _ := listener.Accept()
        defer conn.Close()
        conn.Write([]byte("hola"))
    }()

    // Conectar como cliente
    conn, err := net.Dial("tcp", listener.Addr().String())
    if err != nil {
        t.Fatal(err)
    }
    defer conn.Close()

    buf := make([]byte, 100)
    n, _ := conn.Read(buf)

    if string(buf[:n]) != "hola" {
        t.Errorf("esperaba 'hola', obtuve '%s'", string(buf[:n]))
    }
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: TCP Server Simple - Echo Server

Implementa un servidor TCP que repita (echo) todo lo que reciba del cliente.

**Requisitos:**

- Escuchar en puerto 9000
- Aceptar múltiples conexiones simultáneas
- Para cada línea recibida, devolver: "ECHO: [línea recibida]"
- Comando especial "QUIT" para cerrar conexión
- Logging de conexiones

**Archivo:** `ejercicio1_echo.go`

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
)

// TODO: Implementar servidor echo con los requisitos anteriores
```

**Prueba desde terminal:**

```bash
go run ejercicio1_echo.go

# En otra terminal:
nc localhost 9000
hola
ECHO: hola
mundo
ECHO: mundo
QUIT
```

---

### Ejercicio 2: HTTP Client - Requests con Headers Personalizados

Crea un cliente HTTP que:

**Requisitos:**

- Hacer GET a múltiples URLs
- Agregar headers personalizados (User-Agent, Authorization)
- Manejo de errors (timeout, no alcanzable)
- Verificar status codes
- Mostrar headers de respuesta

**Archivo:** `ejercicio2_http_client.go`

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// TODO: Implementar cliente HTTP con requisitos
```

**Uso:**

```bash
go run ejercicio2_http_client.go
```

---

### Ejercicio 3: HTTP Server - Handlers y Rutas

Implementa un servidor HTTP con múltiples endpoints.

**Requisitos:**

- GET /usuarios - Listar usuarios (JSON)
- POST /usuarios - Crear usuario (JSON)
- GET /usuarios/:id - Obtener usuario específico
- DELETE /usuarios/:id - Eliminar usuario
- Middleware de logging
- Manejo de errores

**Archivo:** `ejercicio3_http_server.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type Usuario struct {
    ID   int    `json:"id"`
    Nombre string `json:"nombre"`
    Email  string `json:"email"`
}

// TODO: Implementar servidor con endpoints RESTful
```

---

### Ejercicio 4: URL Parser - Extractor de Componentes

Crea un programa que analice URLs y extraiga componentes.

**Requisitos:**

- Parsear URLs complejas
- Extraer: scheme, host, port, path, query, fragment
- Construir URLs dinámicamente
- URL encoding/decoding
- Validar URLs

**Archivo:** `ejercicio4_url_parser.go`

```go
package main

import (
    "fmt"
    "net/url"
)

// TODO: Implementar parser de URLs
```

**Uso:**

```bash
go run ejercicio4_url_parser.go
```

---

### Ejercicio 5: TLS Connection - Cliente HTTPS Seguro

Implementa conexión TLS directa (no vía http.Client).

**Requisitos:**

- Conectar a servidor HTTPS
- Verificar certificado
- Mostrar información del certificado
- Enviar HTTP request manualmente
- Leer respuesta

**Archivo:** `ejercicio5_tls.go`

```go
package main

import (
    "crypto/tls"
    "fmt"
    "log"
)

// TODO: Implementar cliente TLS
```

---

## Antipatrones Comunes

### ❌ No Reusable HTTP Client

```go
// MALO: Cliente nuevo para cada request
for i := 0; i < 1000; i++ {
    resp, _ := http.Get(url)
    resp.Body.Close()
}
```

### ✅ Correcto: Reutilizar Client

```go
// BUENO: Reutilizar client
client := &http.Client{
    Timeout: 10 * time.Second,
}

for i := 0; i < 1000; i++ {
    resp, _ := client.Get(url)
    resp.Body.Close()
}
```

---

### ❌ Sin Timeouts

```go
// MALO: Puede quedar congelado indefinidamente
conn, _ := net.Dial("tcp", "host:port")
conn.Read(buf)
```

### ✅ Con Timeouts

```go
// BUENO: Timeout que evita bloqueos infinitos
conn, _ := net.DialTimeout("tcp", "host:port", 5*time.Second)
conn.SetReadDeadline(time.Now().Add(5*time.Second))
conn.Read(buf)
```

---

### ❌ Sin SSL en Producción

```go
// MALO: Ignorar errores de certificado
tr := &http.Transport{
    TLSClientConfig: &tls.Config{
        InsecureSkipVerify: true,
    },
}
```

### ✅ SSL Verificado

```go
// BUENO: Verificación estándar
client := &http.Client{}
resp, _ := client.Get("https://example.com")
```

---

### ❌ Resource Leaks

```go
// MALO: No cerrar bodies
for i := 0; i < 1000; i++ {
    resp, _ := client.Get(url)
    // Olvida cerrar resp.Body!
}
```

### ✅ Cleanup Correcto

```go
// BUENO: Siempre defer close
for i := 0; i < 1000; i++ {
    resp, _ := client.Get(url)
    defer resp.Body.Close()
}
```

---

## Diagrama: Flow de Conexión HTTP

```
┌──────────────┐
│   Cliente    │
└──────┬───────┘
       │
       │ 1. DNS Lookup
       ├─────────────────────┐
       │                     │
       │                  ┌──┴────────┐
       │                  │ DNS Server │
       │                  └──┬────────┘
       │                     │ IP: 93.184.216.34
       │<────────────────────┘
       │
       │ 2. TCP Handshake (SYN, SYN-ACK, ACK)
       │
       ├──────────────────────┐
       │                      │
       │                   ┌──┴──────────┐
       │                   │   Servidor   │
       │                   │  (:80)       │
       │                   └──┬──────────┘
       │<─────────────────────┘
       │
       │ 3. TLS Handshake (si HTTPS)
       │ ├─ ClientHello
       │ ├─ ServerHello + Certificate
       │ ├─ KeyExchange
       │ └─ Finished
       │
       │ 4. HTTP Request
       ├──────────────────────┐
       │ GET / HTTP/1.1        │
       │ Host: example.com     │
       │ ...                   │
       │                       │
       │                    ┌──┴──────────┐
       │                    │   Servidor   │
       │                    │  (procesa)   │
       │                    └──┬──────────┘
       │<─────────────────────┘
       │ 5. HTTP Response
       │ HTTP/1.1 200 OK
       │ ...
       │ [body]
       │
       │ 6. Close Connection
       ├──────────────────────┐
       │ FIN                   │
       │                       │
       │                    ┌──┴──────────┐
       │                    │   Servidor   │
       │                    │  (cierra)    │
       │                    └──┬──────────┘
       │<─────────────────────┘
       │ FIN-ACK
       │
       ✓ Conexión Cerrada
```

---

## Resumen de Conceptos

| Concepto | Protocolo | Uso | Características |
|----------|-----------|-----|-----------------|
| **TCP** | tcp:// | Conexiones confiables | Garantizado, ordenado, overhead |
| **UDP** | udp:// | Streaming, gaming, DNS | Rápido, sin garantía, sin conexión |
| **HTTP** | http:// | Web | Stateless, request-response |
| **HTTPS** | https:// | Web seguro | HTTP + TLS/SSL |
| **DNS** | Sobre UDP/TCP | Resolución de nombres | Distribuido, jerárquico |

---

## Referencias Útiles

- **Go Net Docs**: <https://pkg.go.dev/net>
- **Go HTTP Docs**: <https://pkg.go.dev/net/http>
- **Go TLS Docs**: <https://pkg.go.dev/crypto/tls>
- **RFC 793** (TCP): <https://tools.ietf.org/html/rfc793>
- **RFC 1035** (DNS): <https://tools.ietf.org/html/rfc1035>
- **RFC 7230-7237** (HTTP/1.1): <https://tools.ietf.org/html/rfc7230>

---

## Conclusión

El `net` package de Go proporciona herramientas poderosas para la programación de red. Sus características principales son:

1. **Simplicidad**: APIs limpias y fáciles de usar
2. **Concurrencia**: Goroutines, no threads pesados
3. **Eficiencia**: Reutilización de conexiones, connection pooling
4. **Robustez**: Manejo de errores, timeouts, graceful shutdown

Dominar estos conceptos es fundamental para construir aplicaciones distribuidas, APIs, microservicios y más.

El siguiente capítulo explorará funcionalidad más avanzada como WebSockets, gRPC y protocolos personalizados construidos sobre estas bases.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/30-net-package/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/30-net-package):

```bash
cd examples/30-net-package
go run .
```
