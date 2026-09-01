# Capítulo 41: Microservicios - Arquitectura distribuida

## Índice

1. [¿Qué son Microservicios?](#41-1-qué-son-microservicios)
2. [Comunicación entre Servicios](#41-2-comunicación-entre-servicios)
3. [gRPC Basics](#41-3-grpc-basics)
4. [gRPC Streaming](#41-4-grpc-streaming)
5. [Service Discovery](#41-5-service-discovery)
6. [Circuit Breaker](#41-6-circuit-breaker)
7. [Retry Logic](#41-7-retry-logic)
8. [API Gateway Pattern](#41-8-api-gateway-pattern)
9. [Distributed Tracing](#41-9-distributed-tracing)
10. [Transacciones Distribuidas](#41-10-transacciones-distribuidas)
11. [Buenas Prácticas y Antipatterns](#41-11-buenas-prácticas-y-antipatterns)

---

## 41.1 ¿Qué son Microservicios?

### 41.1.1 Definición y Concepto

Los **microservicios** son una arquitectura de software que estructura una aplicación como un conjunto de servicios pequeños, independientes y distribuidos que se comunican entre sí. Cada servicio es una unidad de negocio autónoma que puede ser desarrollada, desplegada y escalada de forma independiente.

### 41.1.2 Arquitectura Monolítica vs Microservicios

```
┌─────────────────────────────────────────┐
│          ARQUITECTURA MONOLÍTICA        │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────┐    │
│  │ UI / Presentación               │    │
│  └──────────────┬──────────────────┘    │
│  ┌──────────────┴──────────────────┐    │
│  │ Lógica de Negocio               │    │
│  │ - Usuarios                      │    │
│  │ - Productos                     │    │
│  │ - Órdenes                       │    │
│  │ - Pagos                         │    │
│  └──────────────┬──────────────────┘    │
│  ┌──────────────┴──────────────────┐    │
│  │ Base de Datos Compartida        │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│           ARQUITECTURA DE MICROSERVICIOS                       │
├────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  Servicio    │  │  Servicio    │  │  Servicio    │         │
│  │  Usuarios    │  │  Productos   │  │  Órdenes     │         │
│  │              │  │              │  │              │         │
│  │  Puerto 3001 │  │  Puerto 3002 │  │  Puerto 3003 │         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘         │
│  ┌──────┴───────┐  ┌──────┴───────┐  ┌──────┴───────┐         │
│  │  DB Usuarios │  │ DB Productos │  │  DB Órdenes  │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│         ▲                  ▲                  ▲                 │
│         └──────────────────┼──────────────────┘                │
│                      API Gateway                               │
│                      (Enrutamiento)                            │
└────────────────────────────────────────────────────────────────┘
```

### 41.1.3 Comparación Detallada

| Aspecto | Monolítico | Microservicios |
|---------|-----------|-----------------|
| **Desarrollo** | Un equipo en un codebase | Múltiples equipos independientes |
| **Despliegue** | Deploy de toda la aplicación | Deploy independiente por servicio |
| **Escalabilidad** | Escalar toda la aplicación | Escalar servicios específicos |
| **Tecnologías** | Un stack principal | Múltiples stacks permitidos |
| **Mantenimiento** | Código centralizado | Código distribuido |
| **Testing** | Tests integrados más simples | Tests distribuidos complejos |
| **Fallos** | Un fallo afecta todo | Un fallo afecta un servicio |
| **Datos** | Base de datos compartida | Bases de datos independientes |
| **Latencia** | Red local (bajo) | Red distribuida (alto) |
| **Complejidad** | Baja al inicio | Alta, pero escalable |

### 41.1.4 Ventajas de Microservicios

**1. Independencia Operacional**

- Cada equipo puede desarrollar su servicio sin sincronización
- Despliegues independientes sin afectar otros servicios
- Diferentes velocidades de desarrollo por equipo

**2. Escalabilidad Granular**

- Escalar solo los servicios que lo necesitan
- Optimizar recursos por servicio
- Mejor ROI en infraestructura

**3. Flexibilidad Tecnológica**

- Cada servicio elige su stack: Go, Java, Node.js, Python
- Adoptar nuevas tecnologías sin reescribir todo
- Usar la mejor herramienta para cada problema

**4. Tolerancia a Fallos**

- Fallo de un servicio no derriba todo el sistema
- Implementar circuit breakers y fallbacks
- Recuperación de fallos más granular

**5. Tiempo de Mercado (TTM)**

- Equipos independientes trabajan en paralelo
- Lanzar features sin sincronización
- Iterar rápidamente en cada servicio

### 41.1.5 Desventajas y Desafíos

**1. Complejidad Operacional**

```
- Monitoreo distribuido
- Múltiples plataformas de logging
- Trazabilidad compleja (tracing distribuido)
- Debugging de problemas distribuidos
```

**2. Consistencia de Datos**

- Transacciones distribuidas complejas
- Eventual consistency en lugar de ACID
- Necesidad de saga pattern

**3. Latencia de Red**

- Múltiples hops entre servicios
- Mayor latencia que llamadas locales
- Overhead de serialización

**4. Testing Distribuido**

- Testing de integración más complejo
- Necesidad de coordinar múltiples servicios
- Reproducción de bugs difícil

**5. Gestión de Dependencias**

- Control de versiones de APIs
- Backups de cambios de contrato
- Coordinación de cambios incompatibles

### 41.1.6 ¿Cuándo Usar Microservicios?

✅ **USAR MICROSERVICIOS CUANDO:**

- Aplicación es grande y compleja (> 100k LOC)
- Múltiples equipos de desarrollo (> 5 equipos)
- Diferentes dominios de negocio bien definidos
- Necesidad de escalar diferentes partes independientemente
- Despliegues frecuentes e independientes
- Tolerancia a fallos crítica
- Tecnologías diferentes por dominio

❌ **NO USAR MICROSERVICIOS CUANDO:**

- Aplicación pequeña o MVP
- Un único equipo de desarrollo
- Dominios débilmente acoplados
- Baja latencia es crítica
- Transacciones complejas entre servicios
- Recursos limitados para operaciones
- Equipo sin experiencia en sistemas distribuidos

### 41.1.7 Patrones de Microservicios

**1. Database per Service**
Cada microservicio tiene su propia base de datos. Evita acoplamiento de datos pero requiere sincronización.

**2. Event Sourcing**
Guardar todos los cambios de estado como eventos. Permite reproducir el estado y auditoría.

**3. CQRS (Command Query Responsibility Segregation)**
Separar lecturas de escrituras. Optimiza diferentes patrones de acceso.

**4. Saga Pattern**
Transacciones distribuidas mediante orquestación o coreografía de eventos.

### 41.1.8 Teorema CAP y Microservicios

El teorema CAP establece que un sistema distribuido solo puede garantizar 2 de 3 propiedades:

- **C (Consistencia)**: Todos los nodos ven los mismos datos
- **A (Disponibilidad)**: El sistema responde siempre
- **P (Tolerancia a Particiones)**: El sistema sigue funcionando si hay fallo de red

En microservicios típicamente se elige **AP** (Disponibilidad + Tolerancia):

- Sacrificamos Consistencia fuerte por Eventual Consistency
- Garantizamos disponibilidad del servicio
- Toleramos particiones de red

---

## 41.2 Comunicación entre Servicios

### 41.2.1 Patrones de Comunicación

```
┌──────────────────────────────────────────────────────┐
│     PATRONES DE COMUNICACIÓN ENTRE SERVICIOS        │
├──────────────────────────────────────────────────────┤
│ 1. SÍNCRONA                                          │
│    ├─ REST (HTTP)                                    │
│    ├─ gRPC                                           │
│    └─ SOAP / XML-RPC                                 │
│                                                      │
│ 2. ASÍNCRONA                                         │
│    ├─ Message Queue (RabbitMQ, Kafka)               │
│    ├─ Event Bus (Redis, NATS)                       │
│    └─ Pub/Sub (Google Pub/Sub, AWS SNS)             │
│                                                      │
│ 3. HÍBRIDA                                           │
│    └─ Request/Response + Event-driven                │
└──────────────────────────────────────────────────────┘
```

### 41.2.2 REST (HTTP) - Synchronous

REST es el patrón más común para comunicación síncrona.

**Ventajas:**

- Simple de implementar
- Well-established (HTTP/1.1, HTTP/2)
- Fácil debugging
- Amplio soporte de herramientas

**Desventajas:**

- Overhead de HTTP headers
- Latencia mayor que gRPC
- Menos eficiente en ancho de banda

```go
package main

import (
 "bytes"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
)

// Servicio de Usuario
type UserService struct {
 addr string
}

type User struct {
 ID    int    `json:"id"`
 Name  string `json:"name"`
 Email string `json:"email"`
}

func (us *UserService) GetUser(userID int) (*User, error) {
 resp, err := http.Get(fmt.Sprintf("http://%s/users/%d", us.addr, userID))
 if err != nil {
  return nil, fmt.Errorf("failed to get user: %w", err)
 }
 defer resp.Body.Close()

 if resp.StatusCode != http.StatusOK {
  body, _ := io.ReadAll(resp.Body)
  return nil, fmt.Errorf("user service returned %d: %s", resp.StatusCode, string(body))
 }

 var user User
 if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
  return nil, fmt.Errorf("failed to decode user: %w", err)
 }

 return &user, nil
}

func (us *UserService) CreateUser(user *User) (*User, error) {
 data, err := json.Marshal(user)
 if err != nil {
  return nil, err
 }

 resp, err := http.Post(
  fmt.Sprintf("http://%s/users", us.addr),
  "application/json",
  bytes.NewReader(data),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create user: %w", err)
 }
 defer resp.Body.Close()

 var created User
 if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
  return nil, fmt.Errorf("failed to decode created user: %w", err)
 }

 return &created, nil
}

// Servidor HTTP del servicio de usuario
func StartUserServer(addr string) {
 http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
   var user User
   if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
    http.Error(w, "Invalid request", http.StatusBadRequest)
    return
   }

   user.ID = 1 // Simulado
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(user)
  }
 })

 http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
   user := User{ID: 1, Name: "John", Email: "john@example.com"}
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(user)
  }
 })

 http.ListenAndServe(addr, nil)
}

func main() {
 // Iniciar servidor
 go StartUserServer(":3001")

 // Cliente
 userService := &UserService{addr: "localhost:3001"}
 user, err := userService.GetUser(1)
 if err != nil {
  fmt.Printf("Error: %v\n", err)
  return
 }

 fmt.Printf("Usuario obtenido: %+v\n", user)
}
```

### 41.2.3 Event-Driven Asynchronous

Comunicación mediante eventos. Desacopla servicios en el tiempo.

```go
package main

import (
 "encoding/json"
 "fmt"
 "sync"
)

// Event Bus simple con Go channels
type Event struct {
 Type      string          `json:"type"`
 Source    string          `json:"source"`
 Timestamp int64           `json:"timestamp"`
 Data      json.RawMessage `json:"data"`
}

type EventBus struct {
 subscribers map[string][]chan Event
 mu          sync.RWMutex
}

func NewEventBus() *EventBus {
 return &EventBus{
  subscribers: make(map[string][]chan Event),
 }
}

// Subscribe se suscribe a eventos de un tipo
func (eb *EventBus) Subscribe(eventType string) chan Event {
 eb.mu.Lock()
 defer eb.mu.Unlock()

 eventChan := make(chan Event, 100)
 eb.subscribers[eventType] = append(eb.subscribers[eventType], eventChan)
 return eventChan
}

// Publish publica un evento
func (eb *EventBus) Publish(event Event) {
 eb.mu.RLock()
 defer eb.mu.RUnlock()

 subscribers := eb.subscribers[event.Type]
 for _, sub := range subscribers {
  go func(ch chan Event) {
   ch <- event
  }(sub)
 }
}

// Ejemplo de uso
type OrderCreatedEvent struct {
 OrderID int    `json:"order_id"`
 UserID  int    `json:"user_id"`
 Amount  float64 `json:"amount"`
}

func ExampleEventDriven() {
 bus := NewEventBus()

 // Servicio de Pagos se suscribe a órdenes creadas
 paymentEvents := bus.Subscribe("order.created")
 go func() {
  for event := range paymentEvents {
   var orderEvent OrderCreatedEvent
   json.Unmarshal(event.Data, &orderEvent)
   fmt.Printf("[Pagos] Procesando pago para orden %d: $%.2f\n",
    orderEvent.OrderID, orderEvent.Amount)
  }
 }()

 // Servicio de Notificaciones se suscribe a órdenes creadas
 notifEvents := bus.Subscribe("order.created")
 go func() {
  for event := range notifEvents {
   var orderEvent OrderCreatedEvent
   json.Unmarshal(event.Data, &orderEvent)
   fmt.Printf("[Notificaciones] Enviando confirmación para orden %d\n",
    orderEvent.OrderID)
  }
 }()

 // Servicio de Órdenes publica evento
 orderEvent := OrderCreatedEvent{
  OrderID: 123,
  UserID:  1,
  Amount:  99.99,
 }
 data, _ := json.Marshal(orderEvent)
 bus.Publish(Event{
  Type:   "order.created",
  Source: "order-service",
  Data:   data,
 })
}
```

---

## 41.3 gRPC Basics

### 41.3.1 ¿Qué es gRPC?

gRPC es un framework de RPC (Remote Procedure Call) moderno:

- Usa HTTP/2 para transporte
- Protocol Buffers para serialización
- Bidireccional por defecto
- Tipado fuertemente
- Más rápido y eficiente que REST

### 41.3.2 Protocol Buffers (Protobuf)

File: `user.proto`

```protobuf
syntax = "proto3";

package userservice;

option go_package = "github.com/example/userservice/pb";

// Mensaje de usuario
message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
  int32 age = 4;
}

// Request para obtener usuario
message GetUserRequest {
  int32 id = 1;
}

// Response de operación de usuario
message UserResponse {
  User user = 1;
  string message = 2;
}

// Servicio de usuario
service UserService {
  rpc GetUser(GetUserRequest) returns (UserResponse);
  rpc CreateUser(User) returns (UserResponse);
  rpc UpdateUser(User) returns (UserResponse);
  rpc DeleteUser(GetUserRequest) returns (UserResponse);
}
```

### 41.3.3 Generar código Go

```bash
# Instalar protoc compiler y plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generar código Go
protoc --go_out=. --go-grpc_out=. user.proto
```

### 41.3.4 Implementar Servidor gRPC

```go
package main

import (
 "context"
 "fmt"
 "log"
 "net"

 "github.com/example/userservice/pb"
 "google.golang.org/grpc"
)

type UserServiceServer struct {
 pb.UnimplementedUserServiceServer
 users map[int32]*pb.User
}

func NewUserServiceServer() *UserServiceServer {
 return &UserServiceServer{
  users: make(map[int32]*pb.User),
 }
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
 user, ok := s.users[req.Id]
 if !ok {
  return &pb.UserResponse{
   Message: fmt.Sprintf("user %d not found", req.Id),
  }, nil
 }

 return &pb.UserResponse{
  User:    user,
  Message: "user found",
 }, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, user *pb.User) (*pb.UserResponse, error) {
 // Generar ID (simulado)
 user.Id = int32(len(s.users)) + 1
 s.users[user.Id] = user

 return &pb.UserResponse{
  User:    user,
  Message: "user created successfully",
 }, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, user *pb.User) (*pb.UserResponse, error) {
 if _, ok := s.users[user.Id]; !ok {
  return nil, fmt.Errorf("user %d not found", user.Id)
 }

 s.users[user.Id] = user
 return &pb.UserResponse{
  User:    user,
  Message: "user updated successfully",
 }, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
 if _, ok := s.users[req.Id]; !ok {
  return nil, fmt.Errorf("user %d not found", req.Id)
 }

 delete(s.users, req.Id)
 return &pb.UserResponse{
  Message: fmt.Sprintf("user %d deleted", req.Id),
 }, nil
}

func main() {
 // Crear listener
 lis, err := net.Listen("tcp", ":50051")
 if err != nil {
  log.Fatalf("Failed to listen: %v", err)
 }

 // Crear servidor gRPC
 s := grpc.NewServer()
 pb.RegisterUserServiceServer(s, NewUserServiceServer())

 log.Printf("Server listening at %v", lis.Addr())
 if err := s.Serve(lis); err != nil {
  log.Fatalf("Failed to serve: %v", err)
 }
}
```

### 41.3.5 Cliente gRPC

```go
package main

import (
 "context"
 "log"
 "time"

 "github.com/example/userservice/pb"
 "google.golang.org/grpc"
)

func main() {
 // Conectar al servidor
 conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
 if err != nil {
  log.Fatalf("Failed to dial: %v", err)
 }
 defer conn.Close()

 // Crear cliente
 client := pb.NewUserServiceClient(conn)

 // Crear usuario
 ctx, cancel := context.WithTimeout(context.Background(), time.Second)
 defer cancel()

 user := &pb.User{
  Name:  "John Doe",
  Email: "john@example.com",
  Age:   30,
 }

 createResp, err := client.CreateUser(ctx, user)
 if err != nil {
  log.Fatalf("Failed to create user: %v", err)
 }
 log.Printf("Create response: %v", createResp)

 // Obtener usuario
 ctx, cancel = context.WithTimeout(context.Background(), time.Second)
 defer cancel()

 getResp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
 if err != nil {
  log.Fatalf("Failed to get user: %v", err)
 }
 log.Printf("Get response: %v", getResp)
}
```

---

## 41.4 gRPC Streaming

### 41.4.1 Tipos de Streaming

```
┌─────────────────────────────────────────────────┐
│     TIPOS DE STREAMING EN gRPC                  │
├─────────────────────────────────────────────────┤
│ 1. Unary (request/response)                     │
│    Client ──req──> Server                       │
│    Client <─resp─ Server                        │
│                                                 │
│ 2. Server Streaming                             │
│    Client ──req──> Server                       │
│    Client <─resp1 Server                        │
│    Client <─resp2 Server                        │
│    Client <─resp3 Server                        │
│                                                 │
│ 3. Client Streaming                             │
│    Client ──req1──> Server                      │
│    Client ──req2──> Server                      │
│    Client ──req3──> Server                      │
│    Client <─resp─── Server                      │
│                                                 │
│ 4. Bidirectional Streaming                      │
│    Client ──req1──> Server                      │
│    Client ──req2──> Server <─resp1──── Client   │
│    Client ──req3──> Server <─resp2──── Client   │
│                     Server <─resp3──── Client   │
└─────────────────────────────────────────────────┘
```

### 41.4.2 Server Streaming

File: `chat.proto`

```protobuf
syntax = "proto3";

package chatservice;

option go_package = "github.com/example/chatservice/pb";

message Message {
  string sender = 1;
  string content = 2;
  int64 timestamp = 3;
}

message ChatRequest {
  string channel = 1;
  int32 limit = 2;
}

service ChatService {
  rpc GetMessages(ChatRequest) returns (stream Message);
  rpc SendMessage(Message) returns (Message);
}
```

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"

 "github.com/example/chatservice/pb"
 "google.golang.org/grpc"
)

type ChatServiceServer struct {
 pb.UnimplementedChatServiceServer
 messages map[string][]*pb.Message
 mu       sync.RWMutex
}

func NewChatServiceServer() *ChatServiceServer {
 return &ChatServiceServer{
  messages: make(map[string][]*pb.Message),
 }
}

// Server Streaming: El servidor envía múltiples mensajes
func (s *ChatServiceServer) GetMessages(req *pb.ChatRequest, stream grpc.ServerStream) error {
 s.mu.RLock()
 messages := s.messages[req.Channel]
 s.mu.RUnlock()

 for _, msg := range messages {
  if err := stream.SendMsg(msg); err != nil {
   return err
  }
 }

 return nil
}

func (s *ChatServiceServer) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
 s.mu.Lock()
 defer s.mu.Unlock()

 msg.Timestamp = time.Now().Unix()
 s.messages["general"] = append(s.messages["general"], msg)

 return msg, nil
}

// Cliente con Server Streaming
func ClientServerStreaming() {
 conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
 if err != nil {
  log.Fatalf("Failed to dial: %v", err)
 }
 defer conn.Close()

 client := pb.NewChatServiceClient(conn)

 // Solicitar stream de mensajes
 ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
 defer cancel()

 stream, err := client.GetMessages(ctx, &pb.ChatRequest{
  Channel: "general",
  Limit:   100,
 })
 if err != nil {
  log.Fatalf("Failed to get messages: %v", err)
 }

 for {
  msg, err := stream.Recv()
  if err != nil {
   break // EOF
  }
  fmt.Printf("Message from %s: %s\n", msg.Sender, msg.Content)
 }
}
```

### 41.4.3 Client Streaming

```protobuf
service SensorService {
  rpc RecordMetrics(stream Metric) returns (MetricsSummary);
}
```

```go
type MetricServiceServer struct {
 pb.UnimplementedSensorServiceServer
}

// Client Streaming: El cliente envía múltiples mensajes
func (s *MetricServiceServer) RecordMetrics(stream grpc.ClientStream) error {
 var count int32
 var total float64

 for {
  metric, err := stream.Recv()
  if err != nil {
   break // Client terminó stream
  }

  count++
  total += metric.Value
  fmt.Printf("Received metric: %s = %.2f\n", metric.Name, metric.Value)
 }

 avg := total / float64(count)
 return stream.SendAndClose(&pb.MetricsSummary{
  Count:   count,
  Average: avg,
  Total:   total,
 })
}

// Cliente
func ClientClientStreaming() {
 conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
 if err != nil {
  log.Fatalf("Failed to dial: %v", err)
 }
 defer conn.Close()

 client := pb.NewSensorServiceClient(conn)

 ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
 defer cancel()

 stream, err := client.RecordMetrics(ctx)
 if err != nil {
  log.Fatalf("Failed to record metrics: %v", err)
 }

 // Enviar múltiples métricas
 metrics := []pb.Metric{
  {Name: "cpu_usage", Value: 45.5},
  {Name: "memory_usage", Value: 72.3},
  {Name: "disk_usage", Value: 38.1},
 }

 for _, metric := range metrics {
  if err := stream.Send(&metric); err != nil {
   log.Fatalf("Failed to send metric: %v", err)
  }
 }

 summary, err := stream.CloseAndRecv()
 if err != nil {
  log.Fatalf("Failed to receive summary: %v", err)
 }

 fmt.Printf("Summary: Count=%d, Avg=%.2f, Total=%.2f\n",
  summary.Count, summary.Average, summary.Total)
}
```

### 41.4.4 Bidirectional Streaming

```protobuf
service ChatService {
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

```go
// Servidor con bidirectional streaming
func (s *ChatServiceServer) Chat(stream grpc.BidiStreamingServer) error {
 // En goroutine separada, esperar mensajes del cliente
 go func() {
  for {
   msg, err := stream.Recv()
   if err != nil {
    return
   }
   fmt.Printf("Server received from %s: %s\n", msg.From, msg.Content)

   // Simular respuesta
   response := &pb.ChatMessage{
    From:    "Server",
    To:      msg.From,
    Content: fmt.Sprintf("Echo: %s", msg.Content),
   }
   stream.Send(response)
  }
 }()

 // Este goroutine también puede enviar mensajes
 for i := 0; i < 5; i++ {
  time.Sleep(2 * time.Second)
  msg := &pb.ChatMessage{
   From:    "Server",
   To:      "Client",
   Content: fmt.Sprintf("Server message %d", i),
  }
  if err := stream.Send(msg); err != nil {
   return err
  }
 }

 return nil
}

// Cliente con bidirectional streaming
func ClientBidirectionalStreaming() {
 conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
 if err != nil {
  log.Fatalf("Failed to dial: %v", err)
 }
 defer conn.Close()

 client := pb.NewChatServiceClient(conn)

 ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
 defer cancel()

 stream, err := client.Chat(ctx)
 if err != nil {
  log.Fatalf("Failed to create chat stream: %v", err)
 }

 // Goroutine para enviar mensajes
 go func() {
  for i := 0; i < 5; i++ {
   msg := &pb.ChatMessage{
    From:    "Client",
    To:      "Server",
    Content: fmt.Sprintf("Client message %d", i),
   }
   stream.Send(msg)
   time.Sleep(1 * time.Second)
  }
  stream.CloseSend()
 }()

 // Goroutine para recibir mensajes
 go func() {
  for {
   msg, err := stream.Recv()
   if err != nil {
    break
   }
   fmt.Printf("Client received from %s: %s\n", msg.From, msg.Content)
  }
 }()

 // Esperar un poco para recibir respuestas
 time.Sleep(10 * time.Second)
}
```

---

## 41.5 Service Discovery

### 41.5.1 ¿Por qué Service Discovery?

En una arquitectura de microservicios, los servicios pueden estar en múltiples máquinas y cambiar de ubicación. Service Discovery permite que los servicios se encuentren entre sí dinámicamente.

```
┌────────────────────────────────────────────────────────┐
│          SERVICE DISCOVERY PATTERN                     │
├────────────────────────────────────────────────────────┤
│                                                        │
│  1. Servicio se registra en Service Registry          │
│  2. Otros servicios consultan el registro             │
│  3. Cuando servicio falla, se marca como unhealthy    │
│  4. Clientes son notificados de cambios               │
│                                                        │
│     Service A ──register──> ┌──────────────┐          │
│     Service B ──register──> │   Registry   │          │
│     Service C ──register──> │              │          │
│                             └──────────────┘          │
│                                    ▲                   │
│     Client ──query───────────────┘                    │
└────────────────────────────────────────────────────────┘
```

### 41.5.2 Implementar Service Registry Simple

```go
package main

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// ServiceInstance representa una instancia de servicio
type ServiceInstance struct {
 ServiceName string
 Host        string
 Port        int
 Metadata    map[string]string
 HealthCheck func() error
 LastHealth  time.Time
 Healthy     bool
}

// ServiceRegistry mantiene el registro de servicios
type ServiceRegistry struct {
 instances map[string][]*ServiceInstance
 mu        sync.RWMutex
 ticker    *time.Ticker
}

func NewServiceRegistry() *ServiceRegistry {
 sr := &ServiceRegistry{
  instances: make(map[string][]*ServiceInstance),
  ticker:    time.NewTicker(10 * time.Second),
 }

 // Goroutine de health check
 go sr.healthCheckLoop()

 return sr
}

// Register registra una instancia de servicio
func (sr *ServiceRegistry) Register(instance *ServiceInstance) error {
 sr.mu.Lock()
 defer sr.mu.Unlock()

 // Verificar salud antes de registrar
 if err := instance.HealthCheck(); err != nil {
  return fmt.Errorf("health check failed: %w", err)
 }

 instance.LastHealth = time.Now()
 instance.Healthy = true

 sr.instances[instance.ServiceName] = append(
  sr.instances[instance.ServiceName],
  instance,
 )

 fmt.Printf("Registered: %s at %s:%d\n", instance.ServiceName, instance.Host, instance.Port)
 return nil
}

// Deregister deregistra una instancia
func (sr *ServiceRegistry) Deregister(serviceName, host string, port int) {
 sr.mu.Lock()
 defer sr.mu.Unlock()

 instances := sr.instances[serviceName]
 for i, inst := range instances {
  if inst.Host == host && inst.Port == port {
   sr.instances[serviceName] = append(
    instances[:i],
    instances[i+1:]...,
   )
   fmt.Printf("Deregistered: %s at %s:%d\n", serviceName, host, port)
   return
  }
 }
}

// GetServiceInstances obtiene instancias de un servicio
func (sr *ServiceRegistry) GetServiceInstances(serviceName string) []*ServiceInstance {
 sr.mu.RLock()
 defer sr.mu.RUnlock()

 var healthy []*ServiceInstance
 for _, inst := range sr.instances[serviceName] {
  if inst.Healthy {
   healthy = append(healthy, inst)
  }
 }

 return healthy
}

// healthCheckLoop verifica periódicamente la salud de servicios
func (sr *ServiceRegistry) healthCheckLoop() {
 for range sr.ticker.C {
  sr.mu.Lock()

  for serviceName, instances := range sr.instances {
   for _, inst := range instances {
    err := inst.HealthCheck()
    wasHealthy := inst.Healthy

    if err != nil {
     inst.Healthy = false
     if wasHealthy {
      fmt.Printf("⚠️  %s at %s:%d is UNHEALTHY: %v\n",
       serviceName, inst.Host, inst.Port, err)
     }
    } else {
     inst.Healthy = true
     inst.LastHealth = time.Now()
     if !wasHealthy {
      fmt.Printf("✓ %s at %s:%d is HEALTHY again\n",
       serviceName, inst.Host, inst.Port)
     }
    }
   }
  }

  sr.mu.Unlock()
 }
}

// Ejemplo de uso
func ExampleServiceDiscovery() {
 registry := NewServiceRegistry()

 // Registrar usuario-service
 userService := &ServiceInstance{
  ServiceName: "user-service",
  Host:        "localhost",
  Port:        3001,
  Metadata:    map[string]string{"version": "1.0"},
  HealthCheck: func() error {
   // Simular health check
   return nil
  },
 }

 registry.Register(userService)

 // Registrar order-service
 orderService := &ServiceInstance{
  ServiceName: "order-service",
  Host:        "localhost",
  Port:        3002,
  Metadata:    map[string]string{"version": "2.0"},
  HealthCheck: func() error {
   return nil
  },
 }

 registry.Register(orderService)

 // Descubrir servicios
 userInstances := registry.GetServiceInstances("user-service")
 fmt.Printf("Descubiertos %d instancias de user-service\n", len(userInstances))

 for _, inst := range userInstances {
  fmt.Printf("  - %s:%d (v%s)\n", inst.Host, inst.Port, inst.Metadata["version"])
 }

 time.Sleep(15 * time.Second)
}
```

### 41.5.3 Load Balancing

```go
package main

import (
 "fmt"
 "math/rand"
)

type LoadBalancer interface {
 SelectInstance(instances []*ServiceInstance) *ServiceInstance
}

// Round Robin
type RoundRobinLB struct {
 current int
}

func (lb *RoundRobinLB) SelectInstance(instances []*ServiceInstance) *ServiceInstance {
 if len(instances) == 0 {
  return nil
 }
 instance := instances[lb.current%len(instances)]
 lb.current++
 return instance
}

// Random
type RandomLB struct{}

func (lb *RandomLB) SelectInstance(instances []*ServiceInstance) *ServiceInstance {
 if len(instances) == 0 {
  return nil
 }
 return instances[rand.Intn(len(instances))]
}

// Least Connections (simulado)
type LeastConnectionsLB struct {
 connections map[string]int
}

func NewLeastConnectionsLB() *LeastConnectionsLB {
 return &LeastConnectionsLB{
  connections: make(map[string]int),
 }
}

func (lb *LeastConnectionsLB) SelectInstance(instances []*ServiceInstance) *ServiceInstance {
 if len(instances) == 0 {
  return nil
 }

 minAddr := ""
 minConn := int(^uint(0) >> 1) // Max int

 for _, inst := range instances {
  addr := fmt.Sprintf("%s:%d", inst.Host, inst.Port)
  conn := lb.connections[addr]

  if conn < minConn {
   minConn = conn
   minAddr = addr
  }
 }

 // Encontrar instancia por dirección
 for _, inst := range instances {
  if fmt.Sprintf("%s:%d", inst.Host, inst.Port) == minAddr {
   return inst
  }
 }

 return instances[0]
}

func (lb *LeastConnectionsLB) RecordConnection(instance *ServiceInstance) {
 addr := fmt.Sprintf("%s:%d", instance.Host, instance.Port)
 lb.connections[addr]++
}

func (lb *LeastConnectionsLB) ReleaseConnection(instance *ServiceInstance) {
 addr := fmt.Sprintf("%s:%d", instance.Host, instance.Port)
 if lb.connections[addr] > 0 {
  lb.connections[addr]--
 }
}
```

---

## 41.6 Circuit Breaker

### 41.6.1 Patrón Circuit Breaker

```
┌─────────────────────────────────────────────────────────┐
│        ESTADOS DEL CIRCUIT BREAKER                      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│   CLOSED (Normal)                                      │
│   ├─ Permite requests                                  │
│   ├─ Si falla: incrementar contador                    │
│   └─ Si umbral: pasar a OPEN                           │
│                                                         │
│   OPEN (Fallando)                                      │
│   ├─ Rechaza requests                                  │
│   ├─ Devuelve error sin llamar servicio                │
│   └─ Después de timeout: pasar a HALF_OPEN            │
│                                                         │
│   HALF_OPEN (Recuperando)                              │
│   ├─ Permite requests limitados                        │
│   ├─ Si éxito: pasar a CLOSED                          │
│   └─ Si falla: pasar a OPEN                            │
│                                                         │
│   CLOSED ──(failures > threshold)──> OPEN              │
│     ▲                                      │             │
│     └──────(success)─ HALF_OPEN <──(timeout)           │
│                         │                               │
│                    (failure)                            │
│                         │                               │
│                         v                               │
│                       OPEN                              │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 41.6.2 Implementar Circuit Breaker

```go
package main

import (
 "errors"
 "fmt"
 "sync"
 "time"
)

type State string

const (
 StateClosed    State = "CLOSED"
 StateOpen      State = "OPEN"
 StateHalfOpen  State = "HALF_OPEN"
)

type CircuitBreaker struct {
 maxFailures      int
 timeout          time.Duration
 halfOpenRequests int

 state              State
 failureCount       int
 lastFailureTime    time.Time
 successCountInHO   int // Éxitos en estado HALF_OPEN
 mu                 sync.RWMutex
 onStateChange      func(State, State) // Callback al cambiar estado
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
 return &CircuitBreaker{
  maxFailures:      maxFailures,
  timeout:          timeout,
  halfOpenRequests: 3, // Requests permitidos en HALF_OPEN
  state:            StateClosed,
 }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
 cb.mu.Lock()
 defer cb.mu.Unlock()

 // Cambiar estado a HALF_OPEN si pasó el timeout
 if cb.state == StateOpen && time.Since(cb.lastFailureTime) > cb.timeout {
  cb.changeState(StateHalfOpen)
 }

 // En estado OPEN: rechazar request
 if cb.state == StateOpen {
  return errors.New("circuit breaker is OPEN")
 }

 // En estado HALF_OPEN: permitir solo algunos requests
 if cb.state == StateHalfOpen && cb.successCountInHO >= cb.halfOpenRequests {
  return errors.New("circuit breaker is HALF_OPEN (max requests reached)")
 }

 // Ejecutar función
 err := fn()

 if err != nil {
  cb.recordFailure()
 } else {
  cb.recordSuccess()
 }

 return err
}

func (cb *CircuitBreaker) recordFailure() {
 cb.failureCount++
 cb.lastFailureTime = time.Now()

 if cb.state == StateClosed && cb.failureCount >= cb.maxFailures {
  cb.changeState(StateOpen)
 } else if cb.state == StateHalfOpen {
  cb.changeState(StateOpen)
 }
}

func (cb *CircuitBreaker) recordSuccess() {
 if cb.state == StateClosed {
  cb.failureCount = 0
 } else if cb.state == StateHalfOpen {
  cb.successCountInHO++
  if cb.successCountInHO >= cb.halfOpenRequests {
   cb.changeState(StateClosed)
  }
 }
}

func (cb *CircuitBreaker) changeState(newState State) {
 oldState := cb.state
 cb.state = newState
 cb.failureCount = 0
 cb.successCountInHO = 0

 if cb.onStateChange != nil {
  go cb.onStateChange(oldState, newState)
 }

 fmt.Printf("🔄 Circuit Breaker: %s -> %s\n", oldState, newState)
}

func (cb *CircuitBreaker) State() State {
 cb.mu.RLock()
 defer cb.mu.RUnlock()
 return cb.state
}

// Ejemplo de uso
func ExampleCircuitBreaker() {
 cb := NewCircuitBreaker(3, 5*time.Second)
 cb.onStateChange = func(old, new State) {
  fmt.Printf("Estado cambió: %s → %s\n", old, new)
 }

 callCount := 0

 // Simular múltiples llamadas
 for i := 0; i < 15; i++ {
  err := cb.Call(func() error {
   callCount++
   if i < 4 { // Primeros 4 fallan
    return fmt.Errorf("service error %d", i)
   }
   return nil // Después funcionan
  })

  if err != nil {
   fmt.Printf("[%d] ❌ Error: %v (State: %s)\n", i, err, cb.State())
  } else {
   fmt.Printf("[%d] ✓ Success (State: %s)\n", i, cb.State())
  }

  time.Sleep(500 * time.Millisecond)
 }
}
```

### 41.6.3 Fallback Pattern

```go
type ServiceClient struct {
 cb       *CircuitBreaker
 fallback func() (interface{}, error)
}

func (sc *ServiceClient) GetUserWithFallback(userID int) (interface{}, error) {
 var result interface{}

 err := sc.cb.Call(func() error {
  // Intentar llamar al servicio
  user, err := sc.getUser(userID)
  result = user
  return err
 })

 if err != nil {
  fmt.Printf("Usando fallback para usuario %d\n", userID)
  return sc.fallback()
 }

 return result, nil
}

func (sc *ServiceClient) getUser(userID int) (interface{}, error) {
 // Llamar a servicio remoto
 return map[string]interface{}{
  "id":   userID,
  "name": "John Doe",
 }, nil
}
```

---

## 41.7 Retry Logic

### 41.7.1 Estrategias de Retry

```go
package main

import (
 "fmt"
 "math"
 "math/rand"
 "time"
)

type RetryPolicy struct {
 maxRetries    int
 initialDelay  time.Duration
 maxDelay      time.Duration
 multiplier    float64
 jitterFraction float64
}

// Linear Backoff: delay = initialDelay * attempt
func (rp *RetryPolicy) LinearBackoff(attempt int) time.Duration {
 delay := rp.initialDelay * time.Duration(attempt)
 if delay > rp.maxDelay {
  delay = rp.maxDelay
 }
 return delay
}

// Exponential Backoff: delay = initialDelay * (multiplier ^ attempt)
func (rp *RetryPolicy) ExponentialBackoff(attempt int) time.Duration {
 delay := time.Duration(float64(rp.initialDelay) *
  math.Pow(rp.multiplier, float64(attempt)))

 if delay > rp.maxDelay {
  delay = rp.maxDelay
 }

 return delay
}

// Exponential Backoff with Jitter
func (rp *RetryPolicy) ExponentialBackoffWithJitter(attempt int) time.Duration {
 baseDelay := rp.ExponentialBackoff(attempt)
 jitter := time.Duration(rand.Float64() *
  float64(baseDelay) *
  rp.jitterFraction)

 return baseDelay + jitter
}

// Retry ejecuta una función con reintentos
func (rp *RetryPolicy) Retry(fn func() error) error {
 var lastErr error

 for attempt := 0; attempt <= rp.maxRetries; attempt++ {
  err := fn()
  if err == nil {
   return nil // Éxito
  }

  lastErr = err

  if attempt < rp.maxRetries {
   delay := rp.ExponentialBackoffWithJitter(attempt)
   fmt.Printf("Intento %d falló: %v, reintentando en %v\n",
    attempt+1, err, delay)
   time.Sleep(delay)
  } else {
   fmt.Printf("Intento %d falló (final): %v\n", attempt+1, err)
  }
 }

 return lastErr
}

// Ejemplo
func ExampleRetryLogic() {
 policy := RetryPolicy{
  maxRetries:     4,
  initialDelay:   100 * time.Millisecond,
  maxDelay:       5 * time.Second,
  multiplier:     2.0,
  jitterFraction: 0.1,
 }

 attempt := 0
 err := policy.Retry(func() error {
  attempt++
  fmt.Printf("Intento #%d\n", attempt)

  if attempt < 3 {
   return fmt.Errorf("temporal failure")
  }

  return nil // Éxito en tercer intento
 })

 if err != nil {
  fmt.Printf("Final error: %v\n", err)
 } else {
  fmt.Println("✓ Success!")
 }
}
```

### 41.7.2 Condiciones de Reintento

```go
type RetryableError interface {
 error
 IsRetryable() bool
}

type TemporaryError struct {
 message string
}

func (e TemporaryError) Error() string   { return e.message }
func (e TemporaryError) IsRetryable() bool { return true }

type PermanentError struct {
 message string
}

func (e PermanentError) Error() string   { return e.message }
func (e PermanentError) IsRetryable() bool { return false }

// Retry solo si es retryable
func (rp *RetryPolicy) RetryOnCondition(fn func() error) error {
 for attempt := 0; attempt <= rp.maxRetries; attempt++ {
  err := fn()
  if err == nil {
   return nil
  }

  // Chequear si es retryable
  if retryable, ok := err.(RetryableError); ok && !retryable.IsRetryable() {
   return err // No reintentar
  }

  if attempt < rp.maxRetries {
   delay := rp.ExponentialBackoffWithJitter(attempt)
   time.Sleep(delay)
  } else {
   return err
  }
 }

 return nil
}
```

---

## 41.8 API Gateway Pattern

### 41.8.1 ¿Qué es un API Gateway?

El API Gateway es un servidor que actúa como único punto de entrada para todas las solicitudes del cliente. Es responsable de:

- Enrutamiento de requests a servicios
- Autenticación y autorización
- Rate limiting
- Transformación de requests/responses
- Cacheo
- Logging y monitoreo

```
┌─────────────┐
│   Clientes  │
└──────┬──────┘
       │
       v
┌─────────────────────────────────────┐
│       API GATEWAY                   │
├─────────────────────────────────────┤
│ ├─ Autenticación                    │
│ ├─ Rate Limiting                    │
│ ├─ Routing                          │
│ ├─ Load Balancing                   │
│ ├─ Cacheo                           │
│ └─ Transformación                   │
└──────┬────┬────┬──────────────────┘
       │    │    │
       v    v    v
   ┌────┐ ┌────┐ ┌────┐
   │ US │ │ OS │ │ PS │  (Microservicios)
   └────┘ └────┘ └────┘
```

### 41.8.2 Implementar API Gateway

```go
package main

import (
 "context"
 "fmt"
 "io"
 "net/http"
 "net/http/httputil"
 "net/url"
 "strings"
 "sync"
 "time"
)

type Route struct {
 Path    string
 Service string
 URL     string
}

type APIGateway struct {
 routes      map[string]*Route
 rateLimiter map[string]*RateLimiter
 cache       map[string]*CachedResponse
 mu          sync.RWMutex
}

type CachedResponse struct {
 Data      []byte
 ExpiresAt time.Time
}

type RateLimiter struct {
 requests []time.Time
 limit    int
 window   time.Duration
}

func NewAPIGateway() *APIGateway {
 return &APIGateway{
  routes:      make(map[string]*Route),
  rateLimiter: make(map[string]*RateLimiter),
  cache:       make(map[string]*CachedResponse),
 }
}

func (ag *APIGateway) RegisterRoute(route *Route) {
 ag.mu.Lock()
 defer ag.mu.Unlock()
 ag.routes[route.Path] = route
 ag.rateLimiter[route.Path] = &RateLimiter{
  requests: []time.Time{},
  limit:    100,
  window:   time.Minute,
 }
}

func (rl *RateLimiter) IsAllowed() bool {
 now := time.Now()

 // Limpiar requests viejos
 for i := 0; i < len(rl.requests); i++ {
  if now.Sub(rl.requests[i]) > rl.window {
   rl.requests = append(rl.requests[:i], rl.requests[i+1:]...)
   i--
  }
 }

 if len(rl.requests) >= rl.limit {
  return false
 }

 rl.requests = append(rl.requests, now)
 return true
}

func (ag *APIGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 // Autenticación
 if !ag.authenticate(r) {
  http.Error(w, "Unauthorized", http.StatusUnauthorized)
  return
 }

 // Buscar ruta
 route := ag.findRoute(r.URL.Path)
 if route == nil {
  http.Error(w, "Not Found", http.StatusNotFound)
  return
 }

 // Rate limiting
 if !ag.rateLimiter[route.Path].IsAllowed() {
  http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
  return
 }

 // Verificar cache
 cacheKey := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
 if cached := ag.getFromCache(cacheKey); cached != nil {
  w.Header().Set("X-Cache", "HIT")
  w.Write(cached)
  return
 }

 // Proxy a servicio
 ag.proxyRequest(w, r, route)
}

func (ag *APIGateway) authenticate(r *http.Request) bool {
 // Verificar token en Authorization header
 auth := r.Header.Get("Authorization")
 return strings.HasPrefix(auth, "Bearer ")
}

func (ag *APIGateway) findRoute(path string) *Route {
 ag.mu.RLock()
 defer ag.mu.RUnlock()

 for pattern, route := range ag.routes {
  if strings.HasPrefix(path, pattern) {
   return route
  }
 }
 return nil
}

func (ag *APIGateway) proxyRequest(w http.ResponseWriter, r *http.Request, route *Route) {
 targetURL, _ := url.Parse(route.URL)

 proxy := httputil.NewSingleHostReverseProxy(targetURL)

 // Modificar request
 r.Host = targetURL.Host
 r.RequestURI = ""

 // Agregar headers de tracing
 r.Header.Set("X-Gateway-Time", time.Now().Format(time.RFC3339))
 r.Header.Set("X-Original-URL", r.URL.String())

 proxy.ServeHTTP(w, r)
}

func (ag *APIGateway) getFromCache(key string) []byte {
 ag.mu.RLock()
 defer ag.mu.RUnlock()

 cached, ok := ag.cache[key]
 if !ok {
  return nil
 }

 if time.Now().After(cached.ExpiresAt) {
  return nil
 }

 return cached.Data
}

func (ag *APIGateway) cacheResponse(key string, data []byte, ttl time.Duration) {
 ag.mu.Lock()
 defer ag.mu.Unlock()

 ag.cache[key] = &CachedResponse{
  Data:      data,
  ExpiresAt: time.Now().Add(ttl),
 }
}

// Ejemplo
func ExampleAPIGateway() {
 gateway := NewAPIGateway()

 gateway.RegisterRoute(&Route{
  Path:    "/users",
  Service: "user-service",
  URL:     "http://localhost:3001",
 })

 gateway.RegisterRoute(&Route{
  Path:    "/orders",
  Service: "order-service",
  URL:     "http://localhost:3002",
 })

 fmt.Println("API Gateway iniciado en :8080")
 http.ListenAndServe(":8080", gateway)
}
```

---

## 41.9 Distributed Tracing

### 41.9.1 ¿Por qué Distributed Tracing?

En microservicios, un request atraviesa múltiples servicios. Distributed Tracing permite rastrear el flujo completo de un request a través del sistema.

```
┌──────────────┐
│ Client       │
└──────┬───────┘
       │ TraceID: abc123
       v
┌──────────────────┐
│ API Gateway      │ Span 1: 50ms
├──────────────────┤
│ SpanID: span1    │
└──────┬───────────┘
       │ TraceID: abc123
       v
┌──────────────────┐
│ User Service     │ Span 2: 30ms
├──────────────────┤
│ SpanID: span2    │
│ Parent: span1    │
└──────┬───────────┘
       │ TraceID: abc123
       v
┌──────────────────┐
│ Order Service    │ Span 3: 40ms
├──────────────────┤
│ SpanID: span3    │
│ Parent: span2    │
└──────────────────┘

Total: 120ms (50 + 30 + 40)
```

### 41.9.2 Implementar Simple Distributed Tracing

```go
package main

import (
 "context"
 "fmt"
 "time"
)

type TraceContext struct {
 TraceID string
 SpanID  string
 Parent  string
}

type Span struct {
 TraceID   string
 SpanID    string
 Parent    string
 Service   string
 Operation string
 StartTime time.Time
 EndTime   time.Time
 Tags      map[string]interface{}
 Logs      []string
}

type Tracer struct {
 spans []*Span
}

var traceStore = make(map[string]*Tracer)

func (t *Tracer) StartSpan(ctx context.Context, service, operation string) (*Span, context.Context) {
 traceID := ""
 parentSpan := ""

 // Extraer TraceID del contexto
 if val := ctx.Value("trace_id"); val != nil {
  traceID = val.(string)
 } else {
  traceID = generateID()
 }

 // Extraer parent span
 if val := ctx.Value("span_id"); val != nil {
  parentSpan = val.(string)
 }

 spanID := generateID()

 span := &Span{
  TraceID:   traceID,
  SpanID:    spanID,
  Parent:    parentSpan,
  Service:   service,
  Operation: operation,
  StartTime: time.Now(),
  Tags:      make(map[string]interface{}),
  Logs:      []string{},
 }

 t.spans = append(t.spans, span)

 // Crear nuevo contexto con span ID
 ctx = context.WithValue(ctx, "trace_id", traceID)
 ctx = context.WithValue(ctx, "span_id", spanID)

 return span, ctx
}

func (s *Span) End() {
 s.EndTime = time.Now()
 duration := s.EndTime.Sub(s.StartTime)
 fmt.Printf("[%s] %s.%s completed in %v\n",
  s.TraceID[:8], s.Service, s.Operation, duration)
}

func (s *Span) SetTag(key string, value interface{}) {
 s.Tags[key] = value
}

func (s *Span) Log(message string) {
 s.Logs = append(s.Logs, message)
}

func generateID() string {
 return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Ejemplo
func ExampleDistributedTracing() {
 tracer := &Tracer{}

 ctx := context.Background()

 // Simular API Gateway
 span1, ctx := tracer.StartSpan(ctx, "api-gateway", "handle_request")
 span1.SetTag("method", "GET")
 span1.SetTag("path", "/users/123")
 time.Sleep(50 * time.Millisecond)

 // Simular User Service
 span2, ctx := tracer.StartSpan(ctx, "user-service", "get_user")
 span2.SetTag("user_id", 123)
 time.Sleep(30 * time.Millisecond)
 span2.Log("User retrieved from cache")
 span2.End()

 // Simular Order Service
 span3, ctx := tracer.StartSpan(ctx, "order-service", "get_orders")
 span3.SetTag("user_id", 123)
 time.Sleep(40 * time.Millisecond)
 span3.Log("Orders fetched from database")
 span3.End()

 span1.End()

 // Mostrar trace completo
 fmt.Println("\n=== TRACE COMPLETO ===")
 for _, s := range tracer.spans {
  indent := "  "
  if s.Parent != "" {
   indent = "    "
  }
  fmt.Printf("%s[%s] %s.%s (%v)\n",
   indent, s.TraceID[:8], s.Service, s.Operation,
   s.EndTime.Sub(s.StartTime))
 }
}
```

---

## 41.10 Transacciones Distribuidas

### 41.10.1 Saga Pattern - Orquestación

En transacciones distribuidas, no podemos usar transacciones ACID tradicionales. El patrón Saga usa una serie de transacciones locales coordinadas por un orquestador.

```
┌───────────────────────────────────────────────────┐
│          SAGA PATTERN - ORQUESTACIÓN              │
├───────────────────────────────────────────────────┤
│                                                   │
│ Paso 1: Reservar Inventario                       │
│   ✓ Éxito → Paso 2                                │
│   ✗ Fallo → Compensar                             │
│                                                   │
│ Paso 2: Procesar Pago                             │
│   ✓ Éxito → Paso 3                                │
│   ✗ Fallo → Compensar Step 1 + Step 2             │
│                                                   │
│ Paso 3: Crear Orden                               │
│   ✓ Éxito → Completar                             │
│   ✗ Fallo → Compensar Step 1 + Step 2 + Step 3    │
│                                                   │
└───────────────────────────────────────────────────┘
```

### 41.10.2 Implementar Saga Orchestrator

```go
package main

import (
 "fmt"
 "time"
)

type SagaStep struct {
 Name        string
 Action      func() error
 Compensate  func() error
 StepNumber  int
}

type Saga struct {
 steps      []*SagaStep
 completed  []int // IDs de pasos completados
 failed     bool
 failedStep int
}

func NewSaga() *Saga {
 return &Saga{
  steps:     []*SagaStep{},
  completed: []int{},
 }
}

func (s *Saga) AddStep(name string, action func() error, compensate func() error) {
 step := &SagaStep{
  Name:       name,
  Action:     action,
  Compensate: compensate,
  StepNumber: len(s.steps),
 }
 s.steps = append(s.steps, step)
}

func (s *Saga) Execute() error {
 fmt.Println("▶ Iniciando Saga...")

 for i, step := range s.steps {
  fmt.Printf("\n[Paso %d] Ejecutando: %s\n", i+1, step.Name)

  err := step.Action()
  if err != nil {
   fmt.Printf("✗ Error en paso %d: %v\n", i+1, err)
   s.failed = true
   s.failedStep = i
   s.compensate()
   return err
  }

  fmt.Printf("✓ Completado: %s\n", step.Name)
  s.completed = append(s.completed, i)
 }

 fmt.Println("\n✓ Saga completada exitosamente")
 return nil
}

func (s *Saga) compensate() {
 fmt.Println("\n⟲ Compensando transacciones...")

 // Ejecutar compensaciones en orden inverso
 for i := len(s.completed) - 1; i >= 0; i-- {
  stepIdx := s.completed[i]
  step := s.steps[stepIdx]

  fmt.Printf("\n[Compensación %d] %s\n", stepIdx+1, step.Name)

  err := step.Compensate()
  if err != nil {
   fmt.Printf("✗ Error compensando: %v\n", err)
  } else {
   fmt.Printf("✓ Compensada: %s\n", step.Name)
  }
 }
}

// Ejemplo: Crear una orden
func ExampleOrderSaga() {
 saga := NewSaga()

 inventory := 10
 payment := 0.0

 // Paso 1: Reservar inventario
 saga.AddStep(
  "Reservar Inventario",
  func() error {
   if inventory <= 0 {
    return fmt.Errorf("sin inventario")
   }
   inventory--
   fmt.Printf("  Inventario: %d\n", inventory)
   return nil
  },
  func() error {
   inventory++
   fmt.Printf("  Inventario revertido: %d\n", inventory)
   return nil
  },
 )

 // Paso 2: Procesar pago
 saga.AddStep(
  "Procesar Pago",
  func() error {
   // Simular fallo ocasional
   if time.Now().UnixNano()%2 == 0 {
    return fmt.Errorf("pago rechazado")
   }
   payment = 99.99
   fmt.Printf("  Pago procesado: $%.2f\n", payment)
   return nil
  },
  func() error {
   payment = 0
   fmt.Println("  Pago revertido")
   return nil
  },
 )

 // Paso 3: Crear orden
 saga.AddStep(
  "Crear Orden",
  func() error {
   fmt.Println("  Orden creada #12345")
   return nil
  },
  func() error {
   fmt.Println("  Orden cancelada")
   return nil
  },
 )

 err := saga.Execute()
 if err != nil {
  fmt.Printf("\nSaga falló: %v\n", err)
 }
}
```

### 41.10.3 Saga Pattern - Coreografía (Event-Driven)

```go
type OrderCreatedEvent struct {
 OrderID int
 Items   int
 Amount  float64
}

type SagaCoreography struct {
 eventBus *EventBus
}

func NewSagaCoreography(bus *EventBus) *SagaCoreography {
 return &SagaCoreography{eventBus: bus}
}

func (sc *SagaCoreography) Start() {
 // Servicio de órdenes publica OrderCreated
 orders := sc.eventBus.Subscribe("order.created")
 go func() {
  for event := range orders {
   fmt.Println("📦 Orden creada, reservando inventario...")
   sc.eventBus.Publish(Event{
    Type: "inventory.reserved",
    Data: event.Data,
   })
  }
 }()

 // Servicio de inventario escucha OrderCreated y publica InventoryReserved
 inventory := sc.eventBus.Subscribe("inventory.reserved")
 go func() {
  for event := range inventory {
   fmt.Println("💳 Inventario reservado, procesando pago...")
   sc.eventBus.Publish(Event{
    Type: "payment.processed",
    Data: event.Data,
   })
  }
 }()

 // Servicio de pagos escucha InventoryReserved y publica PaymentProcessed
 payment := sc.eventBus.Subscribe("payment.processed")
 go func() {
  for event := range payment {
   fmt.Println("✓ Pago procesado, orden completada!")
   sc.eventBus.Publish(Event{
    Type: "order.completed",
    Data: event.Data,
   })
  }
 }()
}
```

---

## 41.11 Buenas Prácticas y Antipatterns

### 41.11.1 Versionamiento de APIs

```go
// v1 API
func (s *UserServiceV1) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
 user := &pb.User{
  Id:    req.Id,
  Name:  "John",
  Email: "john@example.com",
 }
 return &pb.UserResponse{User: user}, nil
}

// v2 API - Incluye campo adicional
func (s *UserServiceV2) GetUser(ctx context.Context, req *pb.GetUserRequestV2) (*pb.UserResponseV2, error) {
 user := &pb.UserV2{
  Id:       req.Id,
  Name:     "John",
  Email:    "john@example.com",
  Username: "johndoe", // Nuevo campo
 }
 return &pb.UserResponseV2{User: user}, nil
}

// Mantener ambas versiones durante transición
func main() {
 s := grpc.NewServer()
 pb.RegisterUserServiceServer(s, &UserServiceV1{})
 pb.RegisterUserServiceV2Server(s, &UserServiceV2{})
 // ...
}
```

### 41.11.2 Buenas Prácticas

**1. Database per Service**

```
✓ Cada servicio tiene su BD independiente
✓ Evita acoplamiento de datos
✓ Permite elegir BD óptima por servicio
```

**2. Idempotencia**

```go
// Usar request ID para garantizar idempotencia
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
 // Chequear si request ya fue procesado
 if order, exists := cache.Get(req.RequestId); exists {
  return order, nil
 }

 // Procesar orden
 order := s.processOrder(req)

 // Guardar en cache para futuros requests con mismo ID
 cache.Set(req.RequestId, order)

 return order, nil
}
```

**3. Timeouts**

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Si la operación toma más de 5s, se cancela automáticamente
resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
```

### 41.11.3 Antipatterns a Evitar

**❌ ANTIPATTERN 1: Too Fine-Grained Services**

```
Problema: Crear un servicio por cada entidad es excesivo
- Overhead de comunicación
- Dificultad en transacciones
- Complejidad innecesaria

✓ Solución: Agrupar servicios por dominio de negocio
  (ej: OrderService que incluya órdenes Y items)
```

**❌ ANTIPATTERN 2: Shared Database**

```
Problema: Múltiples servicios usando la misma BD
- Acoplamiento directo
- Imposible escalar independientemente
- Cambios de schema afectan a todos

✓ Solución: Database per Service
```

**❌ ANTIPATTERN 3: Distributed Transactions**

```
Problema: Usar transacciones ACID distribuidas
- Muy lento
- Deadlocks comunes
- Fallo del coordinador bloquea todo

✓ Solución: Eventual Consistency con Saga Pattern
```

**❌ ANTIPATTERN 4: Chatty Interfaces**

```
Problema: Múltiples requests para operaciones simples
GET /users/123    -> User
GET /users/123/profile  -> Profile
GET /users/123/settings -> Settings
Total: 3 requests para 1 operación

✓ Solución: Usar gRPC streaming o respuestas agregadas
```

**❌ ANTIPATTERN 5: No Monitoring**

```
Problema: Sin observabilidad en sistema distribuido
- Difícil encontrar donde está el problema
- No saber qué servicio falla

✓ Solución: Implementar
  - Distributed Tracing (Jaeger, Zipkin)
  - Logging centralizado (ELK, Loki)
  - Métricas (Prometheus, Grafana)
  - Health checks
```

### 41.11.4 Comparación: Go vs Java Spring Cloud vs Node.js

| Aspecto | Go | Java Spring Cloud | Node.js |
|---------|-------|------------------|---------|
| **Performance** | Excelente | Bueno | Moderado |
| **Startup time** | < 100ms | 3-10s | 500ms-2s |
| **Memory usage** | Bajo | Alto | Moderado |
| **Goroutines** | Millones | Threads (limitados) | Single-threaded |
| **gRPC support** | Nativo | Bueno | Bueno |
| **Ecosystem** | Creciendo | Maduro | Maduro |
| **Learning curve** | Moderada | Moderada | Baja |
| **Tipo del lenguaje** | Estáticamente tipado | Dinámicamente tipado | Dinámicamente tipado |

### 41.11.5 Resiliencia en Producción

```go
// Patrón completo de resiliencia
type ResilientClient struct {
 cb       *CircuitBreaker
 retry    *RetryPolicy
 timeout  time.Duration
}

func (rc *ResilientClient) Call(fn func() error) error {
 ctx, cancel := context.WithTimeout(context.Background(), rc.timeout)
 defer cancel()

 // Retry dentro de Circuit Breaker
 return rc.cb.Call(func() error {
  return rc.retry.Retry(fn)
 })
}

// Uso
client := &ResilientClient{
 cb:      NewCircuitBreaker(5, 30*time.Second),
 retry:   &RetryPolicy{maxRetries: 3, initialDelay: 100*time.Millisecond},
 timeout: 5 * time.Second,
}

err := client.Call(func() error {
 return callRemoteService()
})
```

---

## Ejercicios Progresivos

### Ejercicio 1: REST API Simple - Dos Servicios Comunicándose

**Objetivo**: Crear dos servicios REST que se comuniquen entre sí.

**Requisitos**:

- Servicio de Usuarios (puerto 3001)
- Servicio de Órdenes (puerto 3002)
- Órdenes llama a Usuarios para obtener datos del cliente
- Manejo de errores y timeouts

**Estructura esperada**:

```
ejercicio1/
├── user-service/
│   └── main.go
├── order-service/
│   └── main.go
└── README.md
```

---

### Ejercicio 2: gRPC Service - Proto + Server + Client

**Objetivo**: Implementar un servicio gRPC completo.

**Requisitos**:

- Definir `message.proto` con servicios gRPC
- Implementar servidor gRPC
- Implementar cliente gRPC
- Manejo de errors y context timeouts

**Estructura esperada**:

```
ejercicio2/
├── pb/
│   ├── message.proto
│   └── generated files
├── server/
│   └── main.go
├── client/
│   └── main.go
└── README.md
```

---

### Ejercicio 3: Circuit Breaker - Proteger contra Servicios Caídos

**Objetivo**: Implementar circuit breaker en cliente HTTP.

**Requisitos**:

- Implementar estados CLOSED/OPEN/HALF_OPEN
- Fallback cuando circuito está abierto
- Monitoreo de transiciones de estado
- Simular servicio caído y recuperación

**Pruebas**:

```
1. Servicio saludable → CLOSED
2. 5 fallos consecutivos → OPEN
3. Rechazar requests mientras OPEN
4. Después de timeout → HALF_OPEN
5. Si éxito en HALF_OPEN → CLOSED
```

---

### Ejercicio 4: Service Discovery - Registrar y Descubrir Servicios

**Objetivo**: Implementar un service registry simple.

**Requisitos**:

- Servicios se registran con su ubicación
- Health checks periódicos
- Descubrimiento dinámico
- Load balancing (Round-robin)

**Componentes**:

- Registry central
- Múltiples instancias de servicio
- Cliente que descubre y balancea carga

---

### Ejercicio 5: Sistema Completo - Múltiples Servicios con Resiliencia

**Objetivo**: Crear un sistema de compras con múltiples servicios resilientes.

**Arquitectura**:

```
API Gateway (8080)
├── Inventory Service (3001)
├── Payment Service (3002)
├── Order Service (3003)
└── Notification Service (3004)
```

**Requisitos**:

- Saga pattern para crear órdenes
- Circuit breakers en cada cliente
- Retry logic con exponential backoff
- Service discovery para ubicar servicios
- Distributed tracing
- Monitoreo de estados

**Flujo de Compra**:

```
1. Cliente POST /orders → API Gateway
2. Gateway valida (auth, rate limit)
3. Order Service reserva inventario
4. Payment Service procesa pago
5. Si éxito → Order confirmada
6. Si fallo → Compensar (revertir inventario)
7. Notification Service envía confirmación
```

---

## Conclusión

Los microservicios son una arquitectura poderosa pero compleja. Go ofrece excelentes herramientas para implementarlos:

- **gRPC** para comunicación eficiente
- **Context** para timeouts y cancelación
- **Goroutines** para concurrencia
- **Interfaces** para abstraer implementaciones

Recuerda:

- **No es todo o nada**: Monolitos bien diseñados son mejores que microservicios mal hechos
- **Complejidad operacional**: Requiere excelente monitoring y logging
- **Consistencia eventual**: Cambiar mentalidad de ACID a BASE
- **Testing distribuido**: Más importante que en monolitos

---

**Referencias y Recursos**:

- <https://microservices.io/>
- <https://grpc.io/>
- <https://www.infoq.com/microservices/>
- Libro: "Building Microservices" by Sam Newman

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/41-microservicios/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/41-microservicios):

```bash
cd examples/41-microservicios
go run .
```
