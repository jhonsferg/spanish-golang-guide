# Capítulo 43: WebSockets - Comunicación en tiempo real

## Introducción

WebSockets revolucionó la comunicación en tiempo real entre navegadores y servidores, eliminando el overhead de polling constante y permitiendo comunicación bidireccional full-duplex. En este capítulo exploraremos cómo implementar aplicaciones robustas, escalables y seguras con WebSockets en Go.

A diferencia de HTTP, que es unidireccional y sin estado, WebSockets establece una conexión persistente que permite a servidor y cliente enviarse mensajes en ambas direcciones sin necesidad de nuevos handshakes. Go es particularmente adecuado para WebSockets gracias a su modelo de concurrencia con goroutines, permitiendo manejar miles de conexiones simultáneas de forma eficiente.

### Casos de Uso

- **Chat en tiempo real**: Aplicaciones colaborativas con múltiples usuarios
- **Notificaciones push**: Alerts, mensajes, actualizaciones del sistema
- **Live streaming de datos**: Cotizaciones, sensores IoT, métricas
- **Juegos multiplayer**: Sincronización de estado en tiempo real
- **Colaboración en vivo**: Editores compartidos, whiteboards, presentaciones
- **Monitoreo y dashboards**: Actualizaciones de KPIs y métricas del sistema
- **IoT y telemetría**: Reporte continuo de eventos de dispositivos

---

## 43.1 ¿Qué son WebSockets?

### 43.1.1 Protocolo WebSocket

WebSocket es un protocolo que proporciona un canal de comunicación full-duplex sobre TCP. Comienza con un handshake HTTP upgrade y luego establece una conexión persistente.

**Fases de la conexión:**

```
1. Cliente → Servidor: GET /ws HTTP/1.1
                       Upgrade: websocket
                       Connection: Upgrade
                       Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==
                       Sec-WebSocket-Version: 13

2. Servidor → Cliente: HTTP/1.1 101 Switching Protocols
                       Upgrade: websocket
                       Connection: Upgrade
                       Sec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=

3. Conexión establecida: Frames bidireccionales
   - Cliente → Servidor: Mensaje de texto/binario
   - Servidor → Cliente: Mensaje de texto/binario
   - Ambas direcciones: Ping/Pong para heartbeat
```

### 43.1.2 Ventajas sobre Polling

**Polling tradicional:**

```
Intervalo: ~1 segundo
├─ 00:00:00 → Cliente: ¿Hay algo nuevo?
├─ 00:00:00 ← Servidor: No
├─ 00:00:01 → Cliente: ¿Hay algo nuevo?
├─ 00:00:01 ← Servidor: No
├─ 00:00:02 → Cliente: ¿Hay algo nuevo?
├─ 00:00:02 ← Servidor: Sí, nuevo mensaje [2KB headers + datos]
└─ Repetir...

Overhead: 36 requests por minuto × headers HTTP = ~3.6MB/hora
Latencia: Hasta 1 segundo de retraso
```

**WebSocket:**

```
Handshake inicial (~350 bytes)
└─ 00:00:02 ← Servidor: Nuevo mensaje (solo datos, sin headers)
   00:00:05 ← Servidor: Otro mensaje
   00:00:12 ← Cliente: Respuesta del usuario

Overhead: Minimal, solo frames binarios
Latencia: < 100ms típicamente
```

### 43.1.3 Estructura de Frames

WebSocket transmite datos en frames. Cada frame tiene:

```
Frame Format:
┌─────────────────────────────────────────────────────┐
│ FIN (1) | RSV (3) | Opcode (4) | MASK (1) │ Len (7)│  2 bytes
├─────────────────────────────────────────────────────┤
│         Extended payload length (0, 2, or 8 bytes)  │
├─────────────────────────────────────────────────────┤
│              Masking key (4 bytes si MASK=1)        │
├─────────────────────────────────────────────────────┤
│              Payload data (x bytes)                 │
└─────────────────────────────────────────────────────┘

Opcodes:
- 0x1: Frame de texto
- 0x2: Frame binario
- 0x8: Close
- 0x9: Ping
- 0xA: Pong
- 0x0: Continuación (para fragmentación)
```

### 43.1.4 Comparación: WebSocket vs Alternativas

```go
// Comparativa de tecnologías en tiempo real

// HTTP Polling
// Pros: Simple, compatible, stateless
// Cons: Latencia alta, overhead de headers, escalabilidad limitada
func pollServer() {
    // Cada segundo: GET /api/updates
    // Respuesta: {"updates": [...]}  // Con headers HTTP
}

// Long Polling
// Pros: Mejor latencia que polling, similar compatibilidad
// Cons: Usa más conexiones, complejidad de timeouts
func longPoll() {
    // GET /api/updates con timeout: espera hasta 30s
    // Si hay updates en 2s, responde inmediatamente
    // Si no hay, responde con timeout
}

// Server-Sent Events (SSE)
// Pros: Más eficiente que polling, HTTP estándar
// Cons: Unidireccional servidor→cliente
func serverSentEvents() {
    // GET /api/events
    // Content-Type: text/event-stream
    // Servidor puede enviar eventos ilimitados
}

// WebSocket
// Pros: Full-duplex, baja latencia, eficiente
// Cons: Requiere soporte del navegador/servidor
func webSocket() {
    // Upgrade HTTP → WebSocket
    // Bidireccional: cliente ↔ servidor
}
```

---

## 43.2 Gorilla WebSocket

### 43.2.1 Setup y Configuración Básica

`github.com/gorilla/websocket` es la librería más popular para WebSockets en Go.

```bash
go get github.com/gorilla/websocket
```

**Estructura básica:**

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

// Configurar upgrader con parámetros de seguridad
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // En producción: verificar Origin header
        // Por ahora, aceptar cualquiera (¡no usar en prod!)
        return true
    },
}

// Handler HTTP que actualiza a WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Actualizar conexión HTTP a WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Cliente conectado")

    // Leer un mensaje
    messageType, data, err := conn.ReadMessage()
    if err != nil {
        log.Println("Read error:", err)
        return
    }

    fmt.Printf("Recibido: %s\n", string(data))

    // Escribir respuesta
    err = conn.WriteMessage(messageType, data)
    if err != nil {
        log.Println("Write error:", err)
    }
}

func main() {
    http.HandleFunc("/ws", handleWebSocket)
    log.Println("Servidor WebSocket en :8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
```

**Cliente JavaScript:**

```javascript
// Conectar a WebSocket
const ws = new WebSocket('ws://localhost:8000/ws');

ws.onopen = () => {
    console.log('Conectado');
    ws.send('Hola desde JS');
};

ws.onmessage = (event) => {
    console.log('Mensaje del servidor:', event.data);
};

ws.onerror = (error) => {
    console.error('Error:', error);
};

ws.onclose = () => {
    console.log('Desconectado');
};
```

### 43.2.2 Configuración de Upgrader Segura

```go
package main

import (
    "log"
    "net/http"
    "net/url"
    "strings"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Validar origen para evitar CSRF
        origin := r.Header.Get("Origin")
        if origin == "" {
            return true // Conexiones sin Origin (mismo origin)
        }

        u, err := url.Parse(origin)
        if err != nil {
            return false
        }

        // Permitir solo dominios específicos
        allowedOrigins := map[string]bool{
            "http://localhost:3000": true,
            "https://app.example.com": true,
        }

        return allowedOrigins[u.Host]
    },
    HandshakeTimeout: 10 * time.Second,
}

func handleWS(w http.ResponseWriter, r *http.Request) {
    // Verificar método GET
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Verificar header Sec-WebSocket-Version
    if strings.ToLower(r.Header.Get("Connection")) != "upgrade" ||
        strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
        http.Error(w, "Invalid WebSocket request", http.StatusBadRequest)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade failed:", err)
        return
    }
    defer conn.Close()

    // Conexión establecida...
}

func main() {
    http.HandleFunc("/ws", handleWS)
    log.Println("Servidor seguro en :8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
```

### 43.2.3 Tipos de Datos Soportados

```go
import (
    "github.com/gorilla/websocket"
)

// WebSocket soporta dos tipos de frames:

// 1. Texto (UTF-8)
func sendTextMessage(conn *websocket.Conn, text string) error {
    return conn.WriteMessage(websocket.TextMessage, []byte(text))
}

// 2. Binario
func sendBinaryMessage(conn *websocket.Conn, data []byte) error {
    return conn.WriteMessage(websocket.BinaryMessage, data)
}

// Control frames (automáticos en muchos casos):
// - Ping: conn.WriteControl(websocket.PingMessage, ...)
// - Pong: automático en respuesta a ping
// - Close: conn.WriteMessage(websocket.CloseMessage, ...)
```

---

## 43.3 Lectura y Escritura de Mensajes

### 43.3.1 Bucle de Lectura

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func handleWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // Bucle de lectura: leer continuamente
    for {
        messageType, data, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error: %v", err)
            }
            return
        }

        fmt.Printf("Tipo: %d, Datos: %s\n", messageType, string(data))

        // Echo: devolver el mismo mensaje
        if err := conn.WriteMessage(messageType, data); err != nil {
            return
        }
    }
}

func main() {
    http.HandleFunc("/ws", handleWS)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
```

### 43.3.2 Lectura con Timeout

```go
package main

import (
    "fmt"
    "time"

    "github.com/gorilla/websocket"
)

func readWithTimeout(conn *websocket.Conn, timeout time.Duration) error {
    // Establecer deadline
    conn.SetReadDeadline(time.Now().Add(timeout))
    defer conn.SetReadDeadline(time.Time{})

    messageType, data, err := conn.ReadMessage()
    if err != nil {
        if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
            fmt.Println("Cliente desconectado")
            return err
        }
        if os.IsTimeout(err) {
            fmt.Println("Lectura expirada")
            return fmt.Errorf("timeout")
        }
        return err
    }

    fmt.Printf("Recibido: %s\n", string(data))
    return nil
}

func handleWSWithTimeout(conn *websocket.Conn) {
    for {
        err := readWithTimeout(conn, 30*time.Second)
        if err != nil {
            conn.Close()
            break
        }
    }
}
```

### 43.3.3 Escritura Segura (Mutex Protection)

En sistemas concurrentes, múltiples goroutines pueden intentar escribir simultáneamente. WebSocket requiere serialización:

```go
package main

import (
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type SafeConn struct {
    conn  *websocket.Conn
    write sync.Mutex
}

func NewSafeConn(conn *websocket.Conn) *SafeConn {
    return &SafeConn{conn: conn}
}

// Escribir de forma segura desde múltiples goroutines
func (sc *SafeConn) WriteMessage(messageType int, data []byte) error {
    sc.write.Lock()
    defer sc.write.Unlock()

    // Establecer deadline para no bloquear indefinidamente
    sc.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
    return sc.conn.WriteMessage(messageType, data)
}

func (sc *SafeConn) WriteJSON(v interface{}) error {
    sc.write.Lock()
    defer sc.write.Unlock()

    sc.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
    return sc.conn.WriteJSON(v)
}

func (sc *SafeConn) Close() error {
    sc.write.Lock()
    defer sc.write.Unlock()
    return sc.conn.Close()
}

// Uso:
func example(conn *websocket.Conn) {
    safe := NewSafeConn(conn)

    // Múltiples goroutines pueden escribir seguramente
    go func() {
        safe.WriteMessage(websocket.TextMessage, []byte("Mensaje 1"))
    }()

    go func() {
        safe.WriteMessage(websocket.TextMessage, []byte("Mensaje 2"))
    }()
}
```

### 43.3.4 WriteJSON y ReadJSON

```go
import (
    "encoding/json"

    "github.com/gorilla/websocket"
)

type Message struct {
    Type    string    `json:"type"`
    Content string    `json:"content"`
    User    string    `json:"user"`
    Time    time.Time `json:"time"`
}

func handleJSONWS(conn *websocket.Conn) {
    for {
        var msg Message
        
        // Leer y deserializar JSON
        err := conn.ReadJSON(&msg)
        if err != nil {
            break
        }

        fmt.Printf("Tipo: %s, Usuario: %s, Contenido: %s\n",
            msg.Type, msg.User, msg.Content)

        // Crear respuesta
        response := Message{
            Type:    "response",
            Content: "Mensaje recibido",
            Time:    time.Now(),
        }

        // Escribir y serializar JSON
        err = conn.WriteJSON(response)
        if err != nil {
            break
        }
    }
    conn.Close()
}
```

---

## 43.4 Control de Conexión

### 43.4.1 Ping/Pong Heartbeat

```go
package main

import (
    "fmt"
    "time"

    "github.com/gorilla/websocket"
)

func handleWSWithHeartbeat(conn *websocket.Conn) {
    // Configurar handlers para ping/pong
    conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    conn.SetPongHandler(func(string) error {
        fmt.Println("Pong recibido del cliente")
        conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    // Goroutine para enviar pings periódicamente
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            if err := conn.WriteControl(
                websocket.PingMessage,
                []byte("ping"),
                time.Now().Add(5*time.Second),
            ); err != nil {
                fmt.Println("Ping error:", err)
                conn.Close()
                return
            }
        }
    }()

    // Bucle de lectura normal
    for {
        messageType, data, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err,
                websocket.CloseGoingAway,
                websocket.CloseAbnormalClosure) {
                fmt.Println("Error:", err)
            }
            return
        }

        fmt.Printf("Recibido: %s\n", string(data))

        // Echo
        conn.WriteMessage(messageType, data)
    }
}
```

**Cliente JavaScript con heartbeat:**

```javascript
class WSClient {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.pingTimeout = null;
        this.reconnectDelay = 1000; // 1 segundo
        this.connect();
    }

    connect() {
        this.ws = new WebSocket(this.url);

        this.ws.onopen = () => {
            console.log('Conectado');
            this.resetPingTimeout();
        };

        this.ws.onmessage = (event) => {
            console.log('Mensaje:', event.data);
            this.resetPingTimeout();
        };

        this.ws.onpong = () => {
            console.log('Pong recibido');
            this.resetPingTimeout();
        };

        this.ws.onerror = () => {
            this.disconnect();
        };

        this.ws.onclose = () => {
            console.log('Desconectado');
            this.reconnect();
        };
    }

    resetPingTimeout() {
        clearTimeout(this.pingTimeout);
        this.pingTimeout = setTimeout(() => {
            console.log('No hay respuesta, reconectando...');
            this.disconnect();
            this.connect();
        }, 90000); // 90 segundos
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
        clearTimeout(this.pingTimeout);
    }

    reconnect() {
        setTimeout(() => {
            console.log('Reconectando...');
            this.connect();
        }, this.reconnectDelay);

        // Backoff exponencial
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
    }

    send(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(message);
        }
    }
}

// Uso:
const client = new WSClient('ws://localhost:8000/ws');
client.send('Hola');
```

### 43.4.2 Cierre Elegante (Graceful Close)

```go
package main

import (
    "fmt"
    "time"

    "github.com/gorilla/websocket"
)

func handleGracefulClose(conn *websocket.Conn) {
    defer conn.Close()

    // Escuchar mensajes
    for {
        messageType, data, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err,
                websocket.CloseGoingAway,
                websocket.CloseAbnormalClosure) {
                fmt.Println("Error:", err)
            }
            return
        }

        if string(data) == "quit" {
            // Cierre elegante
            msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Goodbye")
            conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
            return
        }

        conn.WriteMessage(messageType, data)
    }
}

// Función para cerrar servidor gracefully
func gracefulShutdown(server *http.Server, timeout time.Duration) {
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    <-stop // Esperar señal

    fmt.Println("Iniciando shutdown graceful...")

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        fmt.Println("Shutdown error:", err)
    }

    fmt.Println("Servidor cerrado")
}
```

### 43.4.3 Manejo de Errores de Conexión

```go
package main

import (
    "fmt"
    "log"

    "github.com/gorilla/websocket"
)

func classifyWSError(err error) string {
    if err == nil {
        return "no error"
    }

    if websocket.IsUnexpectedCloseError(err,
        websocket.CloseGoingAway,
        websocket.CloseAbnormalClosure) {
        return "unexpected close"
    }

    switch err {
    case websocket.ErrBadSubprotocol:
        return "bad subprotocol"
    case websocket.ErrBadFrameData:
        return "bad frame data"
    default:
        return fmt.Sprintf("other: %v", err)
    }
}

func handleWSWithErrorHandling(conn *websocket.Conn) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic recovered: %v", r)
        }
        conn.Close()
    }()

    conn.SetReadLimit(512 * 1024) // 512 KB max message

    for {
        messageType, data, err := conn.ReadMessage()
        if err != nil {
            errType := classifyWSError(err)
            log.Printf("Read error [%s]: %v", errType, err)
            return
        }

        if err := conn.WriteMessage(messageType, data); err != nil {
            errType := classifyWSError(err)
            log.Printf("Write error [%s]: %v", errType, err)
            return
        }
    }
}
```

---

## 43.5 Manejo de Mensajes

### 43.5.1 Protocolo de Mensajes Estructura

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

// Protocolo estándar
type WSMessage struct {
    ID        string          `json:"id,omitempty"`
    Type      string          `json:"type"` // "message", "join", "leave", etc.
    Sender    string          `json:"sender,omitempty"`
    Content   string          `json:"content,omitempty"`
    Data      json.RawMessage `json:"data,omitempty"`
    Timestamp time.Time       `json:"timestamp"`
}

// Versión con tipos específicos
type ChatMessage struct {
    ID      string    `json:"id"`
    Sender  string    `json:"sender"`
    Content string    `json:"content"`
    Time    time.Time `json:"time"`
}

type SystemEvent struct {
    Type    string `json:"type"` // "join", "leave", "error"
    Message string `json:"message"`
    Time    time.Time `json:"time"`
}

type MessageHandler struct{}

func (mh *MessageHandler) HandleMessage(msg WSMessage) error {
    switch msg.Type {
    case "chat":
        var chat ChatMessage
        if err := json.Unmarshal(msg.Data, &chat); err != nil {
            return err
        }
        fmt.Printf("Chat: %s - %s\n", chat.Sender, chat.Content)
        return nil

    case "system":
        var sys SystemEvent
        if err := json.Unmarshal(msg.Data, &sys); err != nil {
            return err
        }
        fmt.Printf("System: %s\n", sys.Message)
        return nil

    default:
        return fmt.Errorf("unknown message type: %s", msg.Type)
    }
}
```

### 43.5.2 Fragmentación de Mensajes

```go
package main

import (
    "fmt"
    "io"

    "github.com/gorilla/websocket"
)

// Enviar mensaje grande en fragmentos
func sendLargeMessage(conn *websocket.Conn, data []byte) error {
    const fragmentSize = 1024 * 64 // 64 KB

    // Primer fragmento
    if err := conn.WriteMessage(websocket.TextMessage, 
        data[:fragmentSize]); err != nil {
        return err
    }

    // Fragmentos intermedios
    for i := fragmentSize; i < len(data)-fragmentSize; i += fragmentSize {
        w, err := conn.NextWriter(websocket.TextMessage)
        if err != nil {
            return err
        }

        if _, err := w.Write(data[i : i+fragmentSize]); err != nil {
            w.Close()
            return err
        }

        if err := w.Close(); err != nil {
            return err
        }
    }

    // Último fragmento
    w, err := conn.NextWriter(websocket.TextMessage)
    if err != nil {
        return err
    }

    _, err = w.Write(data[len(data)-fragmentSize:])
    w.Close()
    return err
}

// Leer mensaje en streaming
func readLargeMessage(conn *websocket.Conn) ([]byte, error) {
    mt, r, err := conn.NextReader()
    if err != nil {
        return nil, err
    }

    if mt != websocket.TextMessage && mt != websocket.BinaryMessage {
        return nil, fmt.Errorf("invalid message type: %d", mt)
    }

    // Leer todo el contenido
    return io.ReadAll(r)
}
```

### 43.5.3 Manejo de Compresión

```go
import (
    "compress/flate"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(*http.Request) bool { return true },
}

func init() {
    // Habilitar compresión
    upgrader.HandshakeTimeout = 10 * time.Second

    // Configurar compresión (Go 1.11+)
    // Esto reduce tamaño de mensajes en ~80% para JSON
}

// Verificar si compresión está habilitada
func checkCompression(conn *websocket.Conn) {
    // La compresión es transparente en gorilla/websocket
    // pero puede verificarse en el contexto de la conexión
    fmt.Println("Compresión habilitada (transparente)")
}
```

---

## 43.6 Hub Pattern

### 43.6.1 Arquitectura Hub/Cliente

El patrón hub centraliza todas las conexiones y permite broadcasts eficientes:

```
┌─────────────────────────────────────────────┐
│              HUB (Central)                   │
│  - Manage clients                           │
│  - Route messages                           │
│  - Broadcast                                │
└──────────────────────────────────────────────┘
     ↑           ↑            ↑            ↑
  Client1     Client2      Client3     Client4
```

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "sync"

    "github.com/gorilla/websocket"
)

// Cliente representa una conexión WebSocket
type Client struct {
    ID   string
    send chan json.RawMessage // Canal para enviar mensajes
    conn *websocket.Conn
    hub  *Hub
}

// Hub gestiona todas las conexiones
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan json.RawMessage
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan json.RawMessage, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

// Iniciar el hub (goroutine principal)
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            fmt.Printf("Cliente registrado: %s (%d total)\n", 
                client.ID, len(h.clients))

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
            fmt.Printf("Cliente desregistrado: %s (%d total)\n",
                client.ID, len(h.clients))

        case msg := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- msg:
                default:
                    // Canal lleno, cliente desconectado
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

// Leer mensajes del cliente
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(512 * 1024)
    for {
        _, data, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err,
                websocket.CloseGoingAway,
                websocket.CloseAbnormalClosure) {
                log.Printf("Error: %v", err)
            }
            return
        }

        // Broadcast a todos los clientes
        c.hub.broadcast <- data
    }
}

// Escribir mensajes al cliente
func (c *Client) writePump() {
    defer c.conn.Close()

    for msg := range c.send {
        if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
}

// Handler HTTP
var hub = NewHub()

func init() {
    go hub.Run()
}

func handleWS(w http.ResponseWriter, r *http.Request) {
    var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    client := &Client{
        ID:   r.Header.Get("X-Client-ID"),
        send: make(chan json.RawMessage, 256),
        conn: conn,
        hub:  hub,
    }

    hub.register <- client

    go client.writePump()
    go client.readPump()
}

func main() {
    http.HandleFunc("/ws", handleWS)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
```

### 43.6.2 Broadcast Selectivo

```go
type MessageType string

const (
    BroadcastAll      MessageType = "broadcast"
    DirectMessage     MessageType = "direct"
    RoomBroadcast     MessageType = "room"
)

type DirectedMessage struct {
    Type      MessageType `json:"type"`
    From      string      `json:"from"`
    To        string      `json:"to,omitempty"` // Para direct messages
    Room      string      `json:"room,omitempty"` // Para room messages
    Content   string      `json:"content"`
}

type HubWithRooms struct {
    clients map[*Client]bool
    rooms   map[string]map[*Client]bool
    mu      sync.RWMutex
    
    broadcast chan *DirectedMessage
    register  chan *Client
    unregister chan *Client
}

func NewHubWithRooms() *HubWithRooms {
    return &HubWithRooms{
        clients:    make(map[*Client]bool),
        rooms:      make(map[string]map[*Client]bool),
        broadcast:  make(chan *DirectedMessage, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *HubWithRooms) Run() {
    for {
        select {
        case msg := <-h.broadcast:
            h.routeMessage(msg)
        }
    }
}

func (h *HubWithRooms) routeMessage(msg *DirectedMessage) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    data, _ := json.Marshal(msg)

    switch msg.Type {
    case BroadcastAll:
        // Enviar a todos
        for client := range h.clients {
            select {
            case client.send <- data:
            default:
                // Skip si el canal está lleno
            }
        }

    case DirectMessage:
        // Enviar a cliente específico
        for client := range h.clients {
            if client.ID == msg.To {
                client.send <- data
                break
            }
        }

    case RoomBroadcast:
        // Enviar a todos en la sala
        if room, ok := h.rooms[msg.Room]; ok {
            for client := range room {
                client.send <- data
            }
        }
    }
}

func (h *HubWithRooms) JoinRoom(client *Client, room string) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if _, ok := h.rooms[room]; !ok {
        h.rooms[room] = make(map[*Client]bool)
    }
    h.rooms[room][client] = true
}

func (h *HubWithRooms) LeaveRoom(client *Client, room string) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if room, ok := h.rooms[room]; ok {
        delete(room, client)
    }
}
```

---

## 43.7 Manejo de Errores

### 43.7.1 Errores Comunes

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "os"

    "github.com/gorilla/websocket"
)

func handleWSErrors(conn *websocket.Conn) {
    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            handleReadError(err)
            return
        }

        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            handleWriteError(err)
            return
        }
    }
}

func handleReadError(err error) {
    // Cierre normal
    if websocket.IsCloseError(err, 
        websocket.CloseNormalClosure,
        websocket.CloseGoingAway,
        websocket.CloseAbnormalClosure) {
        log.Println("Cierre normal o inesperado")
        return
    }

    // Error de protocolo
    if websocket.IsUnexpectedCloseError(err,
        websocket.CloseGoingAway,
        websocket.CloseAbnormalClosure) {
        log.Println("Cierre inesperado:", err)
        return
    }

    // Timeout
    var netErr interface{ Timeout() bool }
    if errors.As(err, &netErr) && netErr.Timeout() {
        log.Println("Timeout de lectura")
        return
    }

    // Otros errores
    log.Println("Error de lectura:", err)
}

func handleWriteError(err error) {
    if websocket.IsCloseError(err) {
        log.Println("Conexión cerrada")
        return
    }

    if err == websocket.ErrCloseSent {
        log.Println("Ya se envió close message")
        return
    }

    log.Println("Error de escritura:", err)
}

// Recuperación de errores
func resilientHandler(conn *websocket.Conn) {
    defer conn.Close()

    maxRetries := 3
    retries := 0

    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            retries++
            if retries >= maxRetries {
                log.Printf("Max retries (%d) exceeded\n", maxRetries)
                return
            }

            log.Printf("Error (retry %d/%d): %v\n", retries, maxRetries, err)
            continue
        }

        // Reset contador en lectura exitosa
        retries = 0

        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            log.Println("Write failed:", err)
            return
        }
    }
}
```

### 43.7.2 Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
    state       string // "closed", "open", "half-open"
    mu          sync.Mutex
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       "closed",
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    // Estado abierto: rechazar llamadas
    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
            cb.failures = 0
        } else {
            return fmt.Errorf("circuit breaker open")
        }
    }

    // Intentar llamada
    err := fn()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }

    // Éxito: reset
    if cb.state == "half-open" {
        cb.state = "closed"
    }
    cb.failures = 0
    return nil
}

// Uso
func handleWSWithCircuitBreaker(conn *websocket.Conn) {
    cb := NewCircuitBreaker(3, 5*time.Second)

    for {
        err := cb.Call(func() error {
            _, data, err := conn.ReadMessage()
            if err != nil {
                return err
            }
            return conn.WriteMessage(websocket.TextMessage, data)
        })

        if err != nil {
            if err.Error() == "circuit breaker open" {
                log.Println("Circuit abierto, esperando...")
                time.Sleep(1 * time.Second)
            } else {
                log.Println("Error:", err)
                conn.Close()
                break
            }
        }
    }
}
```

---

## 43.8 Escalado y Concurrencia

### 43.8.1 Múltiples Goroutines por Cliente

```go
type AdvancedClient struct {
    ID     string
    conn   *websocket.Conn
    hub    *Hub

    // Canales para diferentes operaciones
    send     chan json.RawMessage
    commands chan string
    events   chan Event

    // Control
    done   chan struct{}
    ticker *time.Ticker
}

func (c *AdvancedClient) Start() {
    // Lectura
    go c.readPump()
    // Escritura
    go c.writePump()
    // Procesamiento
    go c.processPump()
    // Heartbeat
    go c.heartbeatPump()
}

func (c *AdvancedClient) readPump() {
    defer c.Close()
    for {
        _, data, err := c.conn.ReadMessage()
        if err != nil {
            return
        }

        select {
        case c.commands <- string(data):
        case <-c.done:
            return
        }
    }
}

func (c *AdvancedClient) writePump() {
    defer c.Close()
    for {
        select {
        case msg := <-c.send:
            c.conn.WriteMessage(websocket.TextMessage, msg)
        case <-c.done:
            return
        }
    }
}

func (c *AdvancedClient) processPump() {
    defer c.Close()
    for {
        select {
        case cmd := <-c.commands:
            // Procesar comando
            result := c.handleCommand(cmd)
            c.send <- result

        case event := <-c.events:
            // Procesar evento
            c.handleEvent(event)

        case <-c.done:
            return
        }
    }
}

func (c *AdvancedClient) heartbeatPump() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            c.conn.WriteControl(
                websocket.PingMessage,
                []byte("ping"),
                time.Now().Add(5*time.Second),
            )
        case <-c.done:
            return
        }
    }
}

func (c *AdvancedClient) Close() {
    close(c.done)
    c.conn.Close()
}

func (c *AdvancedClient) handleCommand(cmd string) json.RawMessage {
    // Implementar lógica de comando
    return json.RawMessage(`{"status":"ok"}`)
}

func (c *AdvancedClient) handleEvent(e Event) {
    // Implementar manejo de eventos
}

type Event struct {
    Type string
    Data interface{}
}
```

### 43.8.2 Buffering y Backpressure

```go
type BackpressureClient struct {
    conn       *websocket.Conn
    send       chan json.RawMessage
    maxBuffer  int
    dropping   bool
}

func NewBackpressureClient(conn *websocket.Conn, maxBuffer int) *BackpressureClient {
    return &BackpressureClient{
        conn:      conn,
        send:      make(chan json.RawMessage, maxBuffer),
        maxBuffer: maxBuffer,
    }
}

func (c *BackpressureClient) SendMessage(msg json.RawMessage) error {
    select {
    case c.send <- msg:
        c.dropping = false
        return nil
    default:
        // Buffer lleno
        if !c.dropping {
            log.Printf("Buffer completo, descartando mensajes")
            c.dropping = true
        }

        // Opción 1: Descartar
        return fmt.Errorf("buffer full")

        // Opción 2: Bloquear (causar backpressure)
        // c.send <- msg
        // return nil

        // Opción 3: Timeout
        // select {
        // case c.send <- msg:
        //     return nil
        // case <-time.After(100 * time.Millisecond):
        //     return fmt.Errorf("timeout")
        // }
    }
}

func (c *BackpressureClient) writePump() {
    for msg := range c.send {
        if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
}
```

### 43.8.3 Connection Pooling

```go
type ConnectionPool struct {
    maxConnections int
    activeConns    int
    waitList       chan struct{}
    mu             sync.Mutex
}

func NewConnectionPool(maxConnections int) *ConnectionPool {
    return &ConnectionPool{
        maxConnections: maxConnections,
        waitList:       make(chan struct{}, maxConnections),
    }
}

func (cp *ConnectionPool) Acquire() error {
    cp.mu.Lock()
    if cp.activeConns >= cp.maxConnections {
        cp.mu.Unlock()
        // Esperar slot disponible
        select {
        case <-cp.waitList:
        case <-time.After(30 * time.Second):
            return fmt.Errorf("connection pool timeout")
        }
        cp.mu.Lock()
    }
    cp.activeConns++
    cp.mu.Unlock()
    return nil
}

func (cp *ConnectionPool) Release() {
    cp.mu.Lock()
    cp.activeConns--
    cp.mu.Unlock()

    select {
    case cp.waitList <- struct{}{}:
    default:
    }
}

func (cp *ConnectionPool) Stats() (active int, waiting int) {
    cp.mu.Lock()
    active = cp.activeConns
    cp.mu.Unlock()
    waiting = len(cp.waitList)
    return
}
```

---

## 43.9 Load Balancing

### 43.9.1 Sticky Sessions

En entornos distribuidos, mantener un cliente con el mismo servidor:

```go
package main

import (
    "crypto/md5"
    "fmt"
    "hash"
)

type LoadBalancer struct {
    servers []*Server
}

type Server struct {
    ID       string
    Addr     string
    capacity int
    clients  int
}

func (lb *LoadBalancer) GetServerForClient(clientID string) *Server {
    // Hash consistente: mismo cliente siempre al mismo servidor
    h := md5.New()
    h.Write([]byte(clientID))
    hash := h.Sum(nil)
    hashInt := fmt.Sprintf("%d", hash[0])

    // Seleccionar servidor usando hash
    idx := int(hash[0]) % len(lb.servers)
    return lb.servers[idx]
}

func (lb *LoadBalancer) GetLeastLoadedServer() *Server {
    var leastLoaded *Server
    minClients := int(^uint(0) >> 1) // Max int

    for _, server := range lb.servers {
        if server.clients < minClients {
            minClients = server.clients
            leastLoaded = server
        }
    }

    return leastLoaded
}

// Uso en proxy WebSocket
func handleWSProxy(w http.ResponseWriter, r *http.Request) {
    clientID := r.Header.Get("X-Client-ID")

    server := lb.GetServerForClient(clientID)

    // Conectar al servidor backend
    backend := fmt.Sprintf("ws://%s/ws", server.Addr)

    // Proxy request...
}
```

### 43.9.2 Distribuir Conexiones

```go
type DistributedHub struct {
    nodes map[string]*Node
    mu    sync.RWMutex
}

type Node struct {
    ID      string
    clients map[*Client]bool
}

func (dh *DistributedHub) RegisterNode(id string) {
    dh.mu.Lock()
    defer dh.mu.Unlock()

    dh.nodes[id] = &Node{
        ID:      id,
        clients: make(map[*Client]bool),
    }
}

func (dh *DistributedHub) GetNodeStats() map[string]int {
    dh.mu.RLock()
    defer dh.mu.RUnlock()

    stats := make(map[string]int)
    for id, node := range dh.nodes {
        stats[id] = len(node.clients)
    }
    return stats
}

func (dh *DistributedHub) BroadcastAcrossNodes(msg []byte) {
    dh.mu.RLock()
    defer dh.mu.RUnlock()

    for _, node := range dh.nodes {
        for client := range node.clients {
            select {
            case client.send <- msg:
            default:
                // Skip
            }
        }
    }
}
```

### 43.9.3 Rebalanceo

```go
type Rebalancer struct {
    hub             *Hub
    maxClientsPerNode int
    ticker          *time.Ticker
}

func NewRebalancer(hub *Hub, maxClients int) *Rebalancer {
    return &Rebalancer{
        hub:               hub,
        maxClientsPerNode: maxClients,
        ticker:            time.NewTicker(1 * time.Minute),
    }
}

func (r *Rebalancer) Start() {
    for range r.ticker.C {
        r.rebalance()
    }
}

func (r *Rebalancer) rebalance() {
    // Verificar si algún nodo está sobrecargado
    stats := r.hub.GetStats() // Implementar esto

    for nodeID, clientCount := range stats {
        if clientCount > r.maxClientsPerNode {
            log.Printf("Nodo %s sobrecargado (%d clientes), migrando...",
                nodeID, clientCount)

            // Migrar clientes a otros nodos
            // (Implementación compleja: guardar estado, reconectar, etc.)
        }
    }
}
```

---

## 43.10 Seguridad

### 43.10.1 Autenticación

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/gorilla/websocket"
)

func authenticateClient(r *http.Request) (string, error) {
    // Opción 1: Token en header
    auth := r.Header.Get("Authorization")
    if auth == "" {
        return "", fmt.Errorf("missing authorization header")
    }

    parts := strings.Split(auth, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", fmt.Errorf("invalid authorization format")
    }

    token := parts[1]
    userID, err := validateToken(token)
    return userID, err
}

func validateToken(token string) (string, error) {
    // Implementar JWT validation
    // En producción, usar jwt-go o similar

    // Por ahora, simple validación
    if len(token) < 10 {
        return "", fmt.Errorf("invalid token")
    }

    // Extraer userID del token
    return "user123", nil
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        return origin == "https://app.example.com"
    },
}

func handleWSSecure(w http.ResponseWriter, r *http.Request) {
    // Verificar método GET
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Autenticar
    userID, err := authenticateClient(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        log.Println("Auth error:", err)
        return
    }

    // Upgrade
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    log.Printf("Usuario autenticado: %s", userID)

    // Asociar userID con conexión
    client := &Client{
        ID:   userID,
        conn: conn,
        hub:  hub,
    }

    // Procesar...
}
```

### 43.10.2 Rate Limiting

```go
import (
    "time"

    "golang.org/x/time/rate"
)

type RateLimitedClient struct {
    limiter *rate.Limiter
    conn    *websocket.Conn
}

func NewRateLimitedClient(conn *websocket.Conn, rps float64) *RateLimitedClient {
    // rps: requests per second
    return &RateLimitedClient{
        limiter: rate.NewLimiter(rate.Limit(rps), 10), // Burst de 10
        conn:    conn,
    }
}

func (c *RateLimitedClient) ReadMessage() ([]byte, error) {
    if !c.limiter.Allow() {
        return nil, fmt.Errorf("rate limit exceeded")
    }

    _, data, err := c.conn.ReadMessage()
    return data, err
}

// Implementar per-user rate limiting
type PerUserRateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rps      float64
}

func NewPerUserRateLimiter(rps float64) *PerUserRateLimiter {
    return &PerUserRateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rps:      rps,
    }
}

func (purl *PerUserRateLimiter) Allow(userID string) bool {
    purl.mu.RLock()
    limiter, exists := purl.limiters[userID]
    purl.mu.RUnlock()

    if !exists {
        limiter = rate.NewLimiter(rate.Limit(purl.rps), 10)

        purl.mu.Lock()
        purl.limiters[userID] = limiter
        purl.mu.Unlock()
    }

    return limiter.Allow()
}
```

### 43.10.3 DoS Prevention

```go
type DoSProtection struct {
    maxConnPerIP    int
    maxMsgsPerSec   int
    blacklist       map[string]time.Time
    connections     map[string]int
    messageRates    map[string]*rate.Limiter
    mu              sync.RWMutex
}

func NewDoSProtection(maxConnPerIP, maxMsgsPerSec int) *DoSProtection {
    return &DoSProtection{
        maxConnPerIP:  maxConnPerIP,
        maxMsgsPerSec: maxMsgsPerSec,
        blacklist:     make(map[string]time.Time),
        connections:   make(map[string]int),
        messageRates:  make(map[string]*rate.Limiter),
    }
}

func (dp *DoSProtection) CheckIP(ip string) error {
    dp.mu.Lock()
    defer dp.mu.Unlock()

    // ¿Está en blacklist?
    if expiry, ok := dp.blacklist[ip]; ok {
        if time.Now().Before(expiry) {
            return fmt.Errorf("IP blacklisted")
        }
        delete(dp.blacklist, ip)
    }

    // ¿Demasiadas conexiones?
    if dp.connections[ip] >= dp.maxConnPerIP {
        dp.blacklist[ip] = time.Now().Add(15 * time.Minute)
        return fmt.Errorf("too many connections from IP")
    }

    return nil
}

func (dp *DoSProtection) RecordConnection(ip string) {
    dp.mu.Lock()
    dp.connections[ip]++
    dp.mu.Unlock()
}

func (dp *DoSProtection) ReleaseConnection(ip string) {
    dp.mu.Lock()
    dp.connections[ip]--
    dp.mu.Unlock()
}

func (dp *DoSProtection) CheckMessageRate(ip string) bool {
    dp.mu.Lock()
    limiter, exists := dp.messageRates[ip]
    dp.mu.Unlock()

    if !exists {
        limiter = rate.NewLimiter(
            rate.Limit(dp.maxMsgsPerSec),
            dp.maxMsgsPerSec,
        )

        dp.mu.Lock()
        dp.messageRates[ip] = limiter
        dp.mu.Unlock()
    }

    return limiter.Allow()
}
```

---

## 43.11 Buenas Prácticas y Patterns

### 43.11.1 Testing WebSockets

```go
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gorilla/websocket"
)

func TestWebSocketEcho(t *testing.T) {
    // Crear servidor de prueba
    server := httptest.NewServer(http.HandlerFunc(handleWS))
    defer server.Close()

    // Convertir URL HTTP → WebSocket
    url := "ws" + strings.TrimPrefix(server.URL, "http")

    // Conectar cliente
    ws, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        t.Fatalf("Dial error: %v", err)
    }
    defer ws.Close()

    // Enviar mensaje
    testMsg := "Hello, WebSocket"
    if err := ws.WriteMessage(websocket.TextMessage, []byte(testMsg)); err != nil {
        t.Fatalf("Write error: %v", err)
    }

    // Leer respuesta
    _, msg, err := ws.ReadMessage()
    if err != nil {
        t.Fatalf("Read error: %v", err)
    }

    if string(msg) != testMsg {
        t.Errorf("Expected %q, got %q", testMsg, string(msg))
    }
}

func TestWebSocketJSON(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(handleJSONWS))
    defer server.Close()

    url := "ws" + strings.TrimPrefix(server.URL, "http")

    ws, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        t.Fatalf("Dial error: %v", err)
    }
    defer ws.Close()

    // Enviar JSON
    msg := map[string]string{
        "action": "ping",
        "data":   "test",
    }

    if err := ws.WriteJSON(msg); err != nil {
        t.Fatalf("WriteJSON error: %v", err)
    }

    // Leer respuesta
    var response map[string]string
    if err := ws.ReadJSON(&response); err != nil {
        t.Fatalf("ReadJSON error: %v", err)
    }

    if response["action"] != "pong" {
        t.Errorf("Unexpected response: %v", response)
    }
}

func TestWebSocketMultipleClients(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    server := httptest.NewServer(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            handleWSWithHub(hub, w, r)
        }),
    )
    defer server.Close()

    url := "ws" + strings.TrimPrefix(server.URL, "http")

    // Conectar 2 clientes
    ws1, _, _ := websocket.DefaultDialer.Dial(url, nil)
    defer ws1.Close()

    ws2, _, _ := websocket.DefaultDialer.Dial(url, nil)
    defer ws2.Close()

    // Cliente 1 envía mensaje
    if err := ws1.WriteMessage(websocket.TextMessage, []byte("broadcast")); err != nil {
        t.Fatal(err)
    }

    // Cliente 2 debe recibir
    _, msg, err := ws2.ReadMessage()
    if err != nil {
        t.Fatalf("Read error: %v", err)
    }

    if string(msg) != "broadcast" {
        t.Errorf("Expected broadcast, got %s", string(msg))
    }
}
```

### 43.11.2 Logging y Monitoring

```go
import (
    "log"
    "time"
)

type MetricsCollector struct {
    totalConnections   int64
    activeConnections  int64
    messagesReceived   int64
    messagesSent       int64
    bytesReceived      int64
    bytesSent          int64
    errors             int64

    mu sync.RWMutex
}

func (mc *MetricsCollector) RecordConnection() {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    mc.totalConnections++
    mc.activeConnections++
}

func (mc *MetricsCollector) RecordDisconnection() {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    mc.activeConnections--
}

func (mc *MetricsCollector) RecordMessageReceived(size int) {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    mc.messagesReceived++
    mc.bytesReceived += int64(size)
}

func (mc *MetricsCollector) RecordMessageSent(size int) {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    mc.messagesSent++
    mc.bytesSent += int64(size)
}

func (mc *MetricsCollector) RecordError() {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    mc.errors++
}

func (mc *MetricsCollector) PrintStats() {
    mc.mu.RLock()
    defer mc.mu.RUnlock()

    log.Printf("=== WebSocket Metrics ===")
    log.Printf("Total Connections: %d", mc.totalConnections)
    log.Printf("Active Connections: %d", mc.activeConnections)
    log.Printf("Messages Received: %d", mc.messagesReceived)
    log.Printf("Messages Sent: %d", mc.messagesSent)
    log.Printf("Bytes Received: %.2f MB", float64(mc.bytesReceived)/1024/1024)
    log.Printf("Bytes Sent: %.2f MB", float64(mc.bytesSent)/1024/1024)
    log.Printf("Errors: %d", mc.errors)
}

// Monitoreo periódico
func (mc *MetricsCollector) StartMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        mc.PrintStats()
    }
}
```

### 43.11.3 Deployment

```yaml
# docker-compose.yml para WebSocket server
version: '3.8'

services:
  websocket:
    build: .
    ports:
      - "8000:8000"
    environment:
      - PORT=8000
      - LOG_LEVEL=info
      - MAX_CLIENTS_PER_NODE=10000
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 1G
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

```go
// Usando Redis para compartir estado entre instancias
package main

import (
    "context"

    "github.com/redis/go-redis/v9"
)

type DistributedHub struct {
    client *redis.Client
    nodeID string
}

func NewDistributedHub(nodeID string) *DistributedHub {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    return &DistributedHub{
        client: rdb,
        nodeID: nodeID,
    }
}

func (dh *DistributedHub) BroadcastMessage(ctx context.Context, msg []byte) error {
    // Publicar a todos los nodos vía Redis PubSub
    return dh.client.Publish(ctx, "websocket:broadcast", msg).Err()
}

func (dh *DistributedHub) Subscribe(ctx context.Context, userID string) {
    pubsub := dh.client.Subscribe(ctx, "websocket:broadcast")
    defer pubsub.Close()

    ch := pubsub.Channel()
    for msg := range ch {
        // Procesar mensaje
        _ = msg
    }
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Echo Server WebSocket Simple

Crear un servidor WebSocket que devuelve los mismos mensajes que recibe.

```bash
cd ejercicio-1-echo-server
```

**Requisitos:**
- Servidor HTTP que implementa WebSocket en `/ws`
- Leer mensajes del cliente
- Devolver exactamente lo mismo (echo)
- Manejo básico de errores
- Cliente HTML para probar

```go
// ejercicio-1-echo-server/main.go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func handleWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Cliente conectado")

    for {
        messageType, data, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
                fmt.Println("Conexión cerrada")
            }
            break
        }

        fmt.Printf("Recibido: %s\n", string(data))

        if err := conn.WriteMessage(messageType, data); err != nil {
            break
        }
    }
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Echo Server</title>
        </head>
        <body>
            <h1>WebSocket Echo Server</h1>
            <input id="input" type="text" placeholder="Mensaje">
            <button onclick="send()">Enviar</button>
            <div id="output"></div>

            <script>
                const ws = new WebSocket('ws://localhost:8000/ws');

                ws.onmessage = (event) => {
                    const output = document.getElementById('output');
                    output.innerHTML += '<p>Echo: ' + event.data + '</p>';
                };

                function send() {
                    const input = document.getElementById('input');
                    ws.send(input.value);
                    input.value = '';
                }
            </script>
        </body>
        </html>
    `)
}

func main() {
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/ws", handleWS)

    fmt.Println("Servidor en http://localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
```

```bash
cd ejercicio-1-echo-server
go run main.go
# Abrir http://localhost:8000 en el navegador
```

---

### Ejercicio 2: Chat Básico Multi-Cliente

Crear un chat donde todos los mensajes se difunden a todos los clientes.

```bash
cd ejercicio-2-chat-basico
```

**Requisitos:**
- Servidor con Hub que gestiona múltiples clientes
- Cada cliente envía un mensaje
- El servidor lo difunde a todos los demás
- Mostrar nombre del usuario
- Interfaz HTML mejorada

```go
// ejercicio-2-chat-basico/main.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

type Message struct {
    User    string `json:"user"`
    Content string `json:"content"`
    Type    string `json:"type"`
}

type Client struct {
    name string
    send chan Message
    conn *websocket.Conn
    hub  *Hub
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan Message
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

var hub = &Hub{
    clients:    make(map[*Client]bool),
    broadcast:  make(chan Message, 256),
    register:   make(chan *Client),
    unregister: make(chan *Client),
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            fmt.Printf("Usuario conectado: %s\n", client.name)

            // Notificar a todos
            h.broadcast <- Message{
                Type:    "system",
                Content: client.name + " se unió al chat",
            }

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()

            h.broadcast <- Message{
                Type:    "system",
                Content: client.name + " salió del chat",
            }

        case msg := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- msg:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (c *Client) readPump() {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()

    for {
        var msg Message
        if err := c.conn.ReadJSON(&msg); err != nil {
            return
        }

        msg.User = c.name
        msg.Type = "message"
        hub.broadcast <- msg
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()

    for msg := range c.send {
        c.conn.WriteJSON(msg)
    }
}

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func handleWS(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "Anónimo"
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    client := &Client{
        name: name,
        send: make(chan Message, 256),
        conn: conn,
        hub:  hub,
    }

    hub.register <- client

    go client.writePump()
    go client.readPump()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Chat WebSocket</title>
            <style>
                #chat { width: 400px; height: 300px; border: 1px solid #ccc; overflow-y: auto; margin-bottom: 10px; }
                .message { padding: 5px; border-bottom: 1px solid #eee; }
                .system { color: green; font-style: italic; }
            </style>
        </head>
        <body>
            <h1>Chat en Tiempo Real</h1>
            Tu nombre: <input id="name" type="text" placeholder="Tu nombre" value="Usuario">
            <button onclick="connect()">Conectar</button>
            <div id="chat"></div>
            <input id="message" type="text" placeholder="Mensaje" onkeypress="keypress(event)">

            <script>
                let ws;
                const chatDiv = document.getElementById('chat');

                function connect() {
                    const name = document.getElementById('name').value;
                    ws = new WebSocket('ws://localhost:8000/ws?name=' + encodeURIComponent(name));

                    ws.onmessage = (event) => {
                        const msg = JSON.parse(event.data);
                        const div = document.createElement('div');
                        div.className = 'message ' + (msg.type === 'system' ? 'system' : '');
                        div.textContent = msg.type === 'system' ? msg.content : msg.user + ': ' + msg.content;
                        chatDiv.appendChild(div);
                        chatDiv.scrollTop = chatDiv.scrollHeight;
                    };
                }

                function send() {
                    const input = document.getElementById('message');
                    ws.send(JSON.stringify({ content: input.value }));
                    input.value = '';
                }

                function keypress(e) {
                    if (e.key === 'Enter') send();
                }

                connect();
            </script>
        </body>
        </html>
    `)
}

func main() {
    go hub.run()

    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/ws", handleWS)

    fmt.Println("Servidor en http://localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
```

---

### Ejercicio 3: Hub Pattern con Suscriptores

Implementar un sistema donde clientes pueden suscribirse a canales específicos.

```bash
cd ejercicio-3-hub-suscriptores
```

**Requisitos:**
- Implementar múltiples "canales" o "rooms"
- Los clientes se suscriben a canales específicos
- Mensajes se difunden solo a suscriptores del canal
- Comando para unirse/dejar un canal
- Mostrar lista de canales activos

```go
// ejercicio-3-hub-suscriptores/main.go (fragmento del hub pattern)
type Room struct {
    name    string
    clients map[*Client]bool
    mu      sync.RWMutex
}

type AdvancedHub struct {
    rooms      map[string]*Room
    broadcast  chan BroadcastMessage
    register   chan RoomRegistration
    unregister chan RoomRegistration
    mu         sync.RWMutex
}

type BroadcastMessage struct {
    room    string
    message Message
}

type RoomRegistration struct {
    client *Client
    room   string
    action string // "join" o "leave"
}

func (h *AdvancedHub) run() {
    for {
        select {
        case reg := <-h.register:
            h.mu.Lock()
            if _, ok := h.rooms[reg.room]; !ok {
                h.rooms[reg.room] = &Room{
                    name:    reg.room,
                    clients: make(map[*Client]bool),
                }
            }
            room := h.rooms[reg.room]
            h.mu.Unlock()

            room.mu.Lock()
            room.clients[reg.client] = true
            room.mu.Unlock()

            fmt.Printf("Usuario %s se unió a %s\n", reg.client.name, reg.room)

        case bcast := <-h.broadcast:
            h.mu.RLock()
            room, ok := h.rooms[bcast.room]
            h.mu.RUnlock()

            if !ok {
                continue
            }

            room.mu.RLock()
            for client := range room.clients {
                select {
                case client.send <- bcast.message:
                default:
                    close(client.send)
                    delete(room.clients, client)
                }
            }
            room.mu.RUnlock()
        }
    }
}
```

---

### Ejercicio 4: Streaming de Datos en Vivo

Crear un servidor que emite datos en tiempo real (simular sensor, stock prices, etc).

```bash
cd ejercicio-4-streaming-datos
```

**Requisitos:**
- Servidor que emite datos periódicamente
- Clientes se suscriben a streams específicos
- Mostrar datos en gráfico (usando Chart.js)
- Simular múltiples sensores/fuentes

```go
// ejercicio-4-streaming-datos/main.go
type DataPoint struct {
    Sensor    string    `json:"sensor"`
    Value     float64   `json:"value"`
    Timestamp time.Time `json:"timestamp"`
}

type Sensor struct {
    name      string
    value     float64
    mu        sync.RWMutex
    subscribers map[*Client]bool
}

func (s *Sensor) Update(v float64) {
    s.mu.Lock()
    s.value = v
    s.mu.Unlock()

    point := DataPoint{
        Sensor:    s.name,
        Value:     v,
        Timestamp: time.Now(),
    }

    data, _ := json.Marshal(point)

    s.mu.RLock()
    for client := range s.subscribers {
        client.send <- Message{
            Type:    "data",
            Content: string(data),
        }
    }
    s.mu.RUnlock()
}

func (s *Sensor) simulateData(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    baseValue := 20.0
    for range ticker.C {
        // Simular variación aleatoria
        offset := rand.Float64()*5 - 2.5
        s.Update(baseValue + offset)
    }
}
```

---

### Ejercicio 5: Chat Escalado con Múltiples Servidores

Implementar un sistema de chat que funciona con múltiples instancias del servidor usando Redis.

```bash
cd ejercicio-5-chat-escalado
docker-compose up
```

**Requisitos:**
- Múltiples instancias del servidor WebSocket
- Usar Redis para sincronizar mensajes entre servidores
- Load balancer (nginx) dirigiendo a diferentes servidores
- Persistencia de mensajes (opcional)
- Manejo de sesiones persistentes

```go
// ejercicio-5-chat-escalado/main.go
type DistributedChat struct {
    redisClient *redis.Client
    localHub    *Hub
    nodeID      string
}

func NewDistributedChat(nodeID string) *DistributedChat {
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
    })

    return &DistributedChat{
        redisClient: rdb,
        localHub:    &Hub{...},
        nodeID:      nodeID,
    }
}

func (dc *DistributedChat) publishMessage(ctx context.Context, msg Message) error {
    data, _ := json.Marshal(msg)
    return dc.redisClient.Publish(ctx, "chat:messages", data).Err()
}

func (dc *DistributedChat) subscribeMessages(ctx context.Context) {
    pubsub := dc.redisClient.Subscribe(ctx, "chat:messages")
    defer pubsub.Close()

    ch := pubsub.Channel()
    for msg := range ch {
        var parsed Message
        json.Unmarshal([]byte(msg.Payload), &parsed)

        // Difundir a clientes locales
        dc.localHub.broadcast <- parsed
    }
}
```

```yaml
# ejercicio-5-chat-escalado/docker-compose.yml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  chat1:
    build: .
    environment:
      - NODE_ID=node1
      - REDIS_ADDR=redis:6379
      - PORT=8001
    ports:
      - "8001:8001"
    depends_on:
      - redis

  chat2:
    build: .
    environment:
      - NODE_ID=node2
      - REDIS_ADDR=redis:6379
      - PORT=8002
    ports:
      - "8002:8002"
    depends_on:
      - redis

  chat3:
    build: .
    environment:
      - NODE_ID=node3
      - REDIS_ADDR=redis:6379
      - PORT=8003
    ports:
      - "8003:8003"
    depends_on:
      - redis

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - chat1
      - chat2
      - chat3
```

```nginx
# ejercicio-5-chat-escalado/nginx.conf
upstream chat {
    least_conn;
    server chat1:8001;
    server chat2:8002;
    server chat3:8003;
}

server {
    listen 80;
    server_name localhost;

    location / {
        proxy_pass http://chat;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## Conclusión

En este capítulo hemos cubierto el ecosistema completo de WebSockets en Go:

1. **Fundamentos**: Protocolo, ventajas sobre polling, estructura de frames
2. **Librería gorilla/websocket**: Setup, configuración, tipos de datos
3. **Lectura/Escritura**: Mensajes, JSON, timeouts, escritura segura
4. **Control de conexión**: Heartbeat, cierre elegante, manejo de errores
5. **Manejo de mensajes**: Protocolo, fragmentación, compresión
6. **Hub Pattern**: Broadcast, múltiples clientes, routing selectivo
7. **Errores**: Clasificación, recuperación, circuit breaker
8. **Escalado**: Goroutines, buffering, connection pooling
9. **Load Balancing**: Sticky sessions, distribución, rebalanceo
10. **Seguridad**: Autenticación, rate limiting, DoS prevention
11. **Prácticas**: Testing, monitoring, deployment

### Puntos Clave

- **Concurrencia**: Goroutines para cada cliente + hub central
- **Seguridad**: Validar origen, autenticar, limitar tasa
- **Escalabilidad**: Usar Redis para distribuir entre servidores
- **Monitoreo**: Recopilar métricas de conexiones y mensajes
- **Testing**: Usar httptest con WebSocket para pruebas
- **Performance**: Buffering prudente, heartbeats para detectar clientes muertos

Go ofrece herramientas excepcionales para construir sistemas WebSocket de tiempo real escalables, eficientes y seguros.

---

## Referencias y Recursos Adicionales

- [Gorilla WebSocket Documentation](https://github.com/gorilla/websocket)
- [RFC 6455 - The WebSocket Protocol](https://tools.ietf.org/html/rfc6455)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Writing Web Applications - golang.org](https://golang.org/doc/articles/wiki/)
- [net/http Server Shutdown](https://golang.org/pkg/net/http/#Server.Shutdown)
- [Redis Pub/Sub Go Client](https://github.com/redis/go-redis)
- [Performance Testing WebSockets](https://github.com/hashrocket/websocket-bench)
- [OWASP WebSocket Security](https://owasp.org/www-community/attacks/websocket)

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/43-websockets/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/43-websockets):

```bash
cd examples/43-websockets
go run .
```
