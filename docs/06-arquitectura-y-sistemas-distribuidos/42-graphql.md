# Capítulo 42: GraphQL - APIs flexibles y eficientes

## Índice

1. [¿Qué es GraphQL?](#421-qué-es-graphql)
2. [Schema y Types](#422-schema-y-types)
3. [Queries](#423-queries)
4. [Mutations](#424-mutations)
5. [Subscriptions](#425-subscriptions)
6. [Resolvers](#426-resolvers)
7. [Framework gqlgen](#427-framework-gqlgen)
8. [Validación](#428-validación)
9. [Performance](#429-performance)
10. [Security](#4210-security)
11. [Buenas Prácticas y Patterns](#4211-buenas-prácticas-y-patterns)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 42.1 ¿Qué es GraphQL?

### Introducción a GraphQL

GraphQL (Graph Query Language) es un lenguaje y runtime para APIs que permite a los clientes solicitar exactamente los datos que necesitan. Fue desarrollado por Facebook en 2012 y liberado públicamente en 2015. A diferencia de REST, donde los endpoints retornan estructuras de datos predefinidas, GraphQL permite que el cliente especifique la forma exacta de los datos que desea.

**Filosofía central de GraphQL:**

- **Fuertemente tipado**: Cada campo tiene un tipo definido
- **Consultable**: Los clientes obtienen exactamente lo que piden
- **Evolutivo**: Las APIs pueden crecer sin romper clientes existentes
- **Autoexplicativo**: La documentación está integrada en el schema

### GraphQL vs REST

Comparación detallada entre los dos paradigmas:

```
┌─────────────────────────────────────────────────────────────────┐
│                    GraphQL vs REST                              │
├─────────────────────────┬───────────────────────────────────────┤
│ Aspecto                 │ REST            │ GraphQL             │
├─────────────────────────┼─────────────────┼─────────────────────┤
│ Fetching de datos       │ Endpoint fijo   │ Query personalizada  │
│ Over-fetching           │ Sí (frecuente)  │ No                  │
│ Under-fetching          │ Sí (N+1)        │ No (request graph)  │
│ Endpoints               │ Múltiples       │ Único               │
│ Caché HTTP              │ Fácil           │ Complejo            │
│ Curva de aprendizaje    │ Baja            │ Media               │
│ Complejidad cliente     │ Baja            │ Media-Alta          │
│ Type safety             │ Parcial         │ Total               │
└─────────────────────────┴─────────────────┴─────────────────────┘
```

### Problemas que resuelve GraphQL

**1. Over-fetching en REST**

```
Caso: Obtener nombre del usuario y su avatar

REST:
GET /api/users/123
Retorna: { id, name, email, phone, address, ..., avatar }
Problema: Datos innecesarios

GraphQL:
query {
  user(id: "123") {
    name
    avatar
  }
}
Respuesta exacta.
```

**2. Under-fetching (N+1 problem)**

```
Caso: Obtener usuario con sus posts y comentarios de cada post

REST:
GET /api/users/123                    (1 request)
GET /api/users/123/posts              (1 request)
GET /api/posts/1/comments             (1 request por post)
GET /api/posts/2/comments
...
Total: N+1 requests

GraphQL:
query {
  user(id: "123") {
    name
    posts {
      title
      comments {
        text
        author { name }
      }
    }
  }
}
Total: 1 request
```

### Cuándo usar GraphQL

**Ideal para:**

- APIs públicas con muchos clientes (web, mobile, desktop)
- Aplicaciones que necesitan fetching flexible de datos
- Microservicios que requieren agregación de datos
- Real-time applications con websockets
- Aplicaciones con clientes débiles (IoT, smartwatches)

**REST sigue siendo mejor para:**

- Archivos binarios grandes (descargas)
- APIs simples y lineales
- Cuando el caché HTTP es crítico
- Sistemas con muy baja latencia requerida

### Ecosistema GraphQL en Go

Librerías principales:

```go
// 1. gqlgen - Code generation (recomendado)
github.com/99designs/gqlgen

// 2. graphql-go - Referencia
github.com/graphql-go/graphql

// 3. ent - ORM con GraphQL integrado
entgo.io

// 4. hasura - GraphQL engine
github.com/hasura/graphql-engine

// 5. graph-gophers - Experimental
github.com/graph-gophers/graphql-go
```

---

## 42.2 Schema y Types

### Estructura del Schema GraphQL

El schema es la descripción de todos los tipos, queries y mutations disponibles. En GraphQL SDL (Schema Definition Language), se define así:

```graphql
# Schema básico
type Query {
  hello: String
  user(id: ID!): User
  posts(first: Int): [Post!]!
}

type User {
  id: ID!
  name: String!
  email: String!
  posts: [Post!]!
}

type Post {
  id: ID!
  title: String!
  content: String
  author: User!
}

type Mutation {
  createUser(name: String!): User!
}
```

### Scalar Types (Tipos Escalares)

Los tipos escalares son los bloques básicos que no pueden dividirse más:

```graphql
# Scalar types en GraphQL

# Enteros
type User {
  age: Int              # -2,147,483,648 a 2,147,483,647
  views: Int            # Típicamente para contadores
}

# Números decimales
type Product {
  price: Float          # 32-bit floating point
  weight: Float
}

# Texto
type Article {
  title: String         # UTF-8 string
  content: String
}

# Booleano
type Feature {
  enabled: Boolean
  beta: Boolean
}

# ID
type Node {
  id: ID!              # Identificador único (string o número)
}

# Custom scalars
scalar DateTime
scalar JSON
scalar Upload

type BlogPost {
  createdAt: DateTime!
  metadata: JSON
  attachment: Upload
}
```

**Modificadores de tipos:**

```graphql
type User {
  id: ID!              # No nulo (required)
  name: String!        # String no nulo
  email: String        # Nullable
  tags: [String!]!     # Array no nulo, items no nulos
  friends: [User]      # Array nullable, items nullable
}
```

### Object Types

Los tipos objeto son las estructuras complejas:

```graphql
type User {
  id: ID!
  name: String!
  email: String!
  profile: Profile!          # Type anidado
  posts: [Post!]!            # Relación con lista
  createdAt: String!
}

type Profile {
  bio: String
  avatar: String
  website: String
}

type Post {
  id: ID!
  title: String!
  content: String!
  author: User!
  comments: [Comment!]!
  likeCount: Int!
  createdAt: String!
}

type Comment {
  id: ID!
  text: String!
  author: User!
  post: Post!
}
```

### Interfaces

Las interfaces definen un conjunto de campos que múltiples tipos deben implementar:

```graphql
interface Node {
  id: ID!
}

interface Timestamped {
  createdAt: String!
  updatedAt: String!
}

type User implements Node & Timestamped {
  id: ID!
  name: String!
  createdAt: String!
  updatedAt: String!
}

type Post implements Node & Timestamped {
  id: ID!
  title: String!
  content: String!
  createdAt: String!
  updatedAt: String!
}

# Query que puede retornar diferentes tipos
type Query {
  node(id: ID!): Node
  timeline: [Timestamped!]!
}
```

### Enums

Enumeraciones para valores predefinidos:

```graphql
enum UserRole {
  ADMIN
  USER
  MODERATOR
  GUEST
}

enum PostStatus {
  DRAFT
  PUBLISHED
  ARCHIVED
  DELETED
}

enum SortOrder {
  ASC
  DESC
}

type User {
  id: ID!
  name: String!
  role: UserRole!
}

type Post {
  id: ID!
  title: String!
  status: PostStatus!
}

type Query {
  posts(
    first: Int
    sortBy: String
    order: SortOrder
  ): [Post!]!
}
```

### Unions

Uniones para representar múltiples tipos posibles:

```graphql
type User {
  id: ID!
  name: String!
}

type Bot {
  id: ID!
  name: String!
  version: String!
}

union Account = User | Bot

type Query {
  account(id: ID!): Account
}

# En queries, usamos type guards
query {
  account(id: "1") {
    ... on User {
      name
    }
    ... on Bot {
      name
      version
    }
  }
}
```

### Input Types

Tipos para parámetros de entrada (mutations):

```graphql
input CreateUserInput {
  name: String!
  email: String!
  role: UserRole!
}

input UpdateUserInput {
  name: String
  email: String
  role: UserRole
}

input CreatePostInput {
  title: String!
  content: String!
  tags: [String!]
  status: PostStatus
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User
  createPost(input: CreatePostInput!): Post!
}

# Usage en queries
mutation {
  createUser(input: {
    name: "John"
    email: "john@example.com"
    role: USER
  }) {
    id
    name
  }
}
```

### Schema completo de ejemplo

```graphql
scalar DateTime

enum UserRole {
  ADMIN
  USER
  MODERATOR
}

enum PostStatus {
  DRAFT
  PUBLISHED
  ARCHIVED
}

interface Node {
  id: ID!
}

interface Timestamped {
  createdAt: DateTime!
  updatedAt: DateTime!
}

type User implements Node & Timestamped {
  id: ID!
  name: String!
  email: String!
  role: UserRole!
  bio: String
  avatar: String
  posts: [Post!]!
  followers: [User!]!
  following: [User!]!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Post implements Node & Timestamped {
  id: ID!
  title: String!
  content: String!
  status: PostStatus!
  author: User!
  comments: [Comment!]!
  likes: [User!]!
  likeCount: Int!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Comment implements Node & Timestamped {
  id: ID!
  text: String!
  author: User!
  post: Post!
  likes: [User!]!
  likeCount: Int!
  createdAt: DateTime!
  updatedAt: DateTime!
}

input CreateUserInput {
  name: String!
  email: String!
  bio: String
}

input CreatePostInput {
  title: String!
  content: String!
}

input CreateCommentInput {
  text: String!
  postId: ID!
}

type Query {
  user(id: ID!): User
  users(first: Int, after: String): [User!]!
  post(id: ID!): Post
  posts(first: Int, status: PostStatus): [Post!]!
  me: User
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, name: String, bio: String): User
  createPost(input: CreatePostInput!): Post!
  publishPost(id: ID!): Post
  createComment(input: CreateCommentInput!): Comment!
  likePost(id: ID!): Post!
}

type Subscription {
  postCreated: Post!
  commentAdded(postId: ID!): Comment!
  userFollowed: User!
}
```

---

## 42.3 Queries

### Estructura básica de Queries

Las queries en GraphQL permiten consultar datos del servidor. A diferencia de REST GET, las queries de GraphQL permiten especificar exactamente qué campos necesitas:

```graphql
# Query simple
query {
  user(id: "1") {
    name
    email
  }
}

# Query con múltiples campos
query {
  user(id: "1") {
    name
    email
    posts {
      title
      createdAt
    }
  }
}

# Query con argumentos
query {
  users(first: 10) {
    id
    name
    email
  }

  posts(first: 5, status: PUBLISHED) {
    title
    author {
      name
    }
  }
}
```

### Arguments (Argumentos)

Los argumentos permiten pasar parámetros a los campos:

```graphql
# Argumentos simples
type Query {
  user(id: ID!): User              # Argumento requerido
  users(first: Int): [User!]!      # Argumento opcional
}

# Múltiples argumentos
type Query {
  posts(
    first: Int
    after: String
    status: PostStatus
    authorId: ID
  ): [Post!]!
}

# Ejemplo de uso
query {
  posts(
    first: 20
    after: "cursor123"
    status: PUBLISHED
    authorId: "user5"
  ) {
    id
    title
    author { name }
  }
}

# Argumentos complejos (input types)
input PostFilter {
  status: PostStatus
  authorId: ID
  createdAfter: DateTime
}

type Query {
  posts(
    first: Int
    after: String
    filter: PostFilter
  ): [Post!]!
}

query {
  posts(first: 20, filter: {
    status: PUBLISHED
    authorId: "user5"
  }) {
    id
    title
  }
}
```

### Aliases (Alias)

Los aliases permiten renombrar campos en la respuesta:

```graphql
# Sin aliases - Error si se repite el campo
query {
  user(id: "1") {
    id
  }
  user(id: "2") {  # Error: duplicate field
    id
  }
}

# Con aliases - Válido
query {
  user1: user(id: "1") {
    id
    name
  }
  user2: user(id: "2") {
    id
    name
  }
}

# Respuesta:
{
  "user1": { "id": "1", "name": "Alice" },
  "user2": { "id": "2", "name": "Bob" }
}

# Aliases para renombrar campos
query {
  user(id: "1") {
    userId: id
    userName: name
    userEmail: email
  }
}

# Respuesta:
{
  "user": {
    "userId": "1",
    "userName": "Alice",
    "userEmail": "alice@example.com"
  }
}
```

### Fragments

Los fragments reutilizables definen un conjunto de campos que se pueden usar en múltiples queries:

```graphql
# Definir fragment
fragment UserInfo on User {
  id
  name
  email
  avatar
}

fragment PostPreview on Post {
  id
  title
  excerpt: content
  createdAt
}

# Usar fragments en queries
query {
  user(id: "1") {
    ...UserInfo
    posts {
      ...PostPreview
      author {
        ...UserInfo
      }
    }
  }
}

# Fragments con argumentos
fragment UserWithPosts($count: Int) on User {
  id
  name
  posts(first: $count) {
    title
  }
}

query($userCount: Int) {
  user(id: "1") {
    ...UserWithPosts(count: $userCount)
  }
}
```

### Variables

Las variables permiten parametrizar queries:

```graphql
# Query con variables
query GetUser($id: ID!, $first: Int) {
  user(id: $id) {
    name
    email
    posts(first: $first) {
      title
    }
  }
}

# Variables JSON:
{
  "id": "1",
  "first": 10
}

# Múltiples queries con variables
query GetUserAndPosts($userId: ID!, $postCount: Int!) {
  user(id: $userId) {
    id
    name
    posts(first: $postCount) {
      title
      author { name }
    }
  }

  featuredPosts: posts(first: 5, status: PUBLISHED) {
    title
    likeCount
  }
}

# Variables:
{
  "userId": "1",
  "postCount": 20
}
```

### Directives

Las directivas controlan el comportamiento de la ejecución:

```graphql
# @include - Incluir campo si la condición es true
query GetUser($id: ID!, $includeEmail: Boolean!) {
  user(id: $id) {
    name
    email @include(if: $includeEmail)
  }
}

# Variables: { "id": "1", "includeEmail": true }
# Resultado: { "name": "Alice", "email": "alice@example.com" }

# @skip - Saltar campo si la condición es true
query GetUser($id: ID!, $skipEmail: Boolean!) {
  user(id: $id) {
    name
    email @skip(if: $skipEmail)
  }
}

# Custom directives (definidas en server)
directive @auth(role: UserRole!) on FIELD_DEFINITION
directive @cache(ttl: Int!) on FIELD_DEFINITION

type Query {
  adminPanel: String @auth(role: ADMIN)
  posts: [Post!]! @cache(ttl: 3600)
}
```

### Introspection

GraphQL permite consultar el schema en tiempo de ejecución:

```graphql
# Obtener todos los tipos
query {
  __schema {
    types {
      name
      kind
      description
    }
  }
}

# Obtener detalles de un tipo
query {
  __type(name: "User") {
    name
    kind
    description
    fields {
      name
      type {
        name
        kind
      }
    }
  }
}

# En Go, usar instrospection para validación
package main

import "github.com/graphql-go/graphql"

func GetSchema(schema *graphql.Schema) {
  introspection := graphql.IntrospectionFromSchema(schema)
  // Usar para validación, documentación, etc.
}
```

---

## 42.4 Mutations

### Concepto de Mutations

Las mutations modifican datos en el servidor. Son explícitamente diferentes a las queries para evitar efectos secundarios accidentales:

```graphql
# Estructura básica
type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User
  deleteUser(id: ID!): Boolean!
}

# Usar mutation
mutation {
  createUser(input: {
    name: "John"
    email: "john@example.com"
  }) {
    id
    name
    email
  }
}

# Respuesta:
{
  "createUser": {
    "id": "123",
    "name": "John",
    "email": "john@example.com"
  }
}
```

### Input Types

Los input types definen la estructura de datos de entrada:

```graphql
input CreateUserInput {
  name: String!
  email: String!
  password: String!
  bio: String
  avatar: String
}

input UpdateUserInput {
  name: String
  bio: String
  avatar: String
  role: UserRole
}

input CreatePostInput {
  title: String!
  content: String!
  status: PostStatus
  tags: [String!]
}

input UpdatePostInput {
  title: String
  content: String
  status: PostStatus
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User
  createPost(input: CreatePostInput!): Post!
  updatePost(id: ID!, input: UpdatePostInput!): Post
}
```

### Operaciones CRUD

CRUD completo con mutations:

```graphql
# Create
mutation CreateNewUser {
  createUser(input: {
    name: "Alice"
    email: "alice@example.com"
    password: "secure123"
  }) {
    id
    name
    email
  }
}

# Read (query, no mutation)
query GetUser {
  user(id: "1") {
    id
    name
    email
  }
}

# Update
mutation UpdateUserEmail {
  updateUser(id: "1", input: {
    email: "newemail@example.com"
  }) {
    id
    email
    updatedAt
  }
}

# Delete
mutation DeleteUser {
  deleteUser(id: "1")
}

# Respuesta: { "deleteUser": true }
```

### Múltiples mutations

Ejecutar varias mutations en una sola petición:

```graphql
mutation CreateAndPublishPost {
  # Primera mutation
  createPost: createPost(input: {
    title: "New Post"
    content: "Content here"
    status: DRAFT
  }) {
    id
    status
  }

  # Segunda mutation (se ejecuta después)
  publishPost: publishPost(id: "123") {
    id
    status
    publishedAt
  }

  # Tercera mutation
  likePost: likePost(id: "123") {
    likeCount
  }
}
```

**Importante:** Las mutations se ejecutan secuencialmente, no en paralelo.

### Errores en Mutations

Manejo de errores:

```graphql
type Mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
}

type CreateUserPayload {
  user: User
  errors: [UserError!]!
  success: Boolean!
}

type UserError {
  field: String!
  message: String!
}

# Uso:
mutation {
  createUser(input: { name: "", email: "invalid" }) {
    user {
      id
      name
    }
    errors {
      field
      message
    }
    success
  }
}

# Respuesta:
{
  "createUser": {
    "user": null,
    "errors": [
      { "field": "name", "message": "Name cannot be empty" },
      { "field": "email", "message": "Invalid email format" }
    ],
    "success": false
  }
}
```

### Transacciones en Mutations

Patrón para operaciones transaccionales:

```graphql
type Mutation {
  transferMoney(
    from: ID!
    to: ID!
    amount: Float!
  ): TransferPayload!
}

type TransferPayload {
  fromAccount: Account!
  toAccount: Account!
  transaction: Transaction!
  success: Boolean!
}

mutation {
  transferMoney(from: "acc1", to: "acc2", amount: 100) {
    fromAccount { balance }
    toAccount { balance }
    transaction {
      id
      amount
      timestamp
    }
    success
  }
}
```

---

## 42.5 Subscriptions

### WebSockets y Tiempo Real

Las subscriptions permiten que el servidor envíe datos al cliente en tiempo real:

```graphql
type Subscription {
  # Nueva notificación
  notificationAdded: Notification!

  # Cambios en un recurso específico
  postUpdated(id: ID!): Post!

  # Stream de eventos
  userFollowed: User!

  # Actualizaciones de un comentario
  commentAdded(postId: ID!): Comment!
}

# Estructura WebSocket GraphQL
# 1. Conexión: {"type": "connection_init"}
# 2. Subscription: {
#     "id": "1",
#     "type": "start",
#     "payload": {
#       "query": "subscription { postUpdated(id: \"1\") { id title } }"
#     }
#   }
# 3. Servidor envía datos: {
#     "id": "1",
#     "type": "data",
#     "payload": { "data": { "postUpdated": {...} } }
#   }
```

### Implementación básica

```go
// subscription.go
package main

import (
 "context"
 "github.com/99designs/gqlgen/graphql"
)

type Subscription struct {
 // Canal para broadcast de eventos
 events chan *Post
}

// Resolver para postUpdated
func (r *subscriptionResolver) PostUpdated(
 ctx context.Context,
 id string,
) (<-chan *Post, error) {

 posts := make(chan *Post, 100)

 go func() {
  defer close(posts)

  for {
   select {
   case <-ctx.Done():
    return
   case post := <-r.events:
    if post.ID == id {
     posts <- post
    }
   }
  }
 }()

 return posts, nil
}

// Publicar evento de actualización
func (r *mutationResolver) UpdatePost(
 ctx context.Context,
 id string,
 input UpdatePostInput,
) (*Post, error) {

 post := &Post{
  ID:      id,
  Title:   input.Title,
  Content: input.Content,
 }

 // Notificar suscriptores
 r.subscription.events <- post

 return post, nil
}
```

### Server-Sent Events (SSE)

Alternativa a WebSockets para subscriptions:

```go
// sse_subscription.go
package main

import (
 "encoding/json"
 "fmt"
 "net/http"
 "time"
)

type EventBroker struct {
 clients map[chan interface{}]bool
 publish chan interface{}
}

func NewEventBroker() *EventBroker {
 broker := &EventBroker{
  clients: make(map[chan interface{}]bool),
  publish: make(chan interface{}),
 }

 go broker.run()
 return broker
}

func (b *EventBroker) run() {
 for {
  select {
  case client := <-b.subscribe:
   b.clients[client] = true
  case client := <-b.unsubscribe:
   delete(b.clients, client)
  case event := <-b.publish:
   // Enviar a todos los clientes
   for client := range b.clients {
    client <- event
   }
  }
 }
}

func (b *EventBroker) ServeSSE(w http.ResponseWriter, r *http.Request) {
 flusher, ok := w.(http.Flusher)
 if !ok {
  http.Error(w, "Streaming not supported", http.StatusInternalServerError)
  return
 }

 w.Header().Set("Content-Type", "text/event-stream")
 w.Header().Set("Cache-Control", "no-cache")
 w.Header().Set("Connection", "keep-alive")

 events := make(chan interface{})
 b.subscribe <- events
 defer func() { b.unsubscribe <- events }()

 for {
  select {
  case <-r.Context().Done():
   return
  case event := <-events:
   data, _ := json.Marshal(event)
   fmt.Fprintf(w, "data: %s\n\n", string(data))
   flusher.Flush()
  }
 }
}
```

### Casos de uso reales

```go
// chat_subscription.go - Chat en tiempo real
type ChatMessage struct {
 ID        string    `json:"id"`
 Author    *User     `json:"author"`
 Text      string    `json:"text"`
 Timestamp time.Time `json:"timestamp"`
}

type Subscription struct {
 messages chan *ChatMessage
}

func (r *subscriptionResolver) MessageAdded(
 ctx context.Context,
 roomID string,
) (<-chan *ChatMessage, error) {

 messages := make(chan *ChatMessage, 100)

 go func() {
  defer close(messages)
  for {
   select {
   case <-ctx.Done():
    return
   case msg := <-r.subscription.messages:
    if msg.RoomID == roomID {
     messages <- msg
    }
   }
  }
 }()

 return messages, nil
}

// notifications_subscription.go - Notificaciones
type Notification struct {
 ID        string    `json:"id"`
 Type      string    `json:"type"` // new_follower, post_liked, etc
 User      *User     `json:"user"`
 Data      string    `json:"data"`
 CreatedAt time.Time `json:"created_at"`
}

func (r *subscriptionResolver) NotificationReceived(
 ctx context.Context,
) (<-chan *Notification, error) {

 // Obtener ID del usuario del contexto
 userID := ctx.Value("user_id").(string)

 notifications := make(chan *Notification, 50)

 go func() {
  defer close(notifications)
  for {
   select {
   case <-ctx.Done():
    return
   case notif := <-r.notificationBroker.Subscribe(userID):
    notifications <- notif
   }
  }
 }()

 return notifications, nil
}
```

---

## 42.6 Resolvers

### Concepto de Resolvers

Los resolvers son funciones que resuelven los valores de los campos en un schema GraphQL. Cada campo tiene un resolver asociado:

```
Query
├── user(id) ──> Resolver que busca en BD
└── posts() ──> Resolver que lista posts

User
├── name ──> Resolver que retorna nombre
├── email ──> Resolver que retorna email
└── posts ──> Resolver que busca posts del usuario (N+1!)
```

### Function-based Resolvers

```go
// resolver.go
package main

import (
 "context"
 "database/sql"
 "errors"
)

// QueryResolver implementa Query
type QueryResolver struct {
 db *sql.DB
}

// User resolver
func (r *QueryResolver) User(
 ctx context.Context,
 id string,
) (*User, error) {
 user := &User{}

 err := r.db.QueryRowContext(ctx,
  "SELECT id, name, email FROM users WHERE id = ?",
  id,
 ).Scan(&user.ID, &user.Name, &user.Email)

 if err == sql.ErrNoRows {
  return nil, errors.New("user not found")
 }
 return user, err
}

// Posts resolver
func (r *QueryResolver) Posts(
 ctx context.Context,
 first *int,
) ([]*Post, error) {
 limit := 10
 if first != nil {
  limit = *first
 }

 rows, err := r.db.QueryContext(ctx,
  "SELECT id, title, content, author_id FROM posts LIMIT ?",
  limit,
 )
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 posts := []*Post{}
 for rows.Next() {
  post := &Post{}
  err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID)
  if err != nil {
   return nil, err
  }
  posts = append(posts, post)
 }

 return posts, nil
}

// UserResolver para campos de User
type UserResolver struct {
 db *sql.DB
}

// Posts resolver en User
func (r *UserResolver) Posts(
 ctx context.Context,
 user *User,
) ([]*Post, error) {

 rows, err := r.db.QueryContext(ctx,
  "SELECT id, title, content FROM posts WHERE author_id = ?",
  user.ID,
 )
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 posts := []*Post{}
 for rows.Next() {
  post := &Post{}
  err := rows.Scan(&post.ID, &post.Title, &post.Content)
  if err != nil {
   return nil, err
  }
  posts = append(posts, post)
 }

 return posts, nil
}
```

### Field Resolvers

Los field resolvers son para campos específicos:

```go
// field_resolvers.go
package main

import "context"

type User struct {
 ID    string
 Name  string
 Email string
 // No incluir PostCount en struct
}

// Post count como resolver, no campo
func (r *UserResolver) PostCount(
 ctx context.Context,
 user *User,
) (int, error) {
 var count int

 err := r.db.QueryRowContext(ctx,
  "SELECT COUNT(*) FROM posts WHERE author_id = ?",
  user.ID,
 ).Scan(&count)

 return count, err
}

// Campos derivados
func (r *UserResolver) DisplayName(
 ctx context.Context,
 user *User,
) (string, error) {
 // Lógica: si name está vacío, usar email
 if user.Name == "" {
  return user.Email, nil
 }
 return user.Name, nil
}

// Avatar URL con transformación
func (r *UserResolver) AvatarURL(
 ctx context.Context,
 user *User,
) (string, error) {
 // Generar URL gravatar o CDN
 hash := md5.Sum([]byte(user.Email))
 return fmt.Sprintf("https://www.gravatar.com/avatar/%x", hash), nil
}

// Autorización como resolver
func (r *UserResolver) Email(
 ctx context.Context,
 user *User,
) (string, error) {
 // Solo mostrar email al usuario mismo o admin
 currentUser := ctx.Value("user").(*User)

 if currentUser.ID == user.ID || currentUser.Role == "ADMIN" {
  return user.Email, nil
 }

 return "", errors.New("unauthorized")
}
```

### Middleware en Resolvers

```go
// middleware.go
package main

import (
 "context"
 "log"
 "time"
)

// Middleware de timing
func TimingMiddleware(
 next graphql.FieldResolver,
) graphql.FieldResolver {
 return func(ctx context.Context, field *graphql.Field) (interface{}, error) {
  start := time.Now()

  result, err := next(ctx, field)

  duration := time.Since(start)
  log.Printf("Field %s took %v", field.Name, duration)

  return result, err
 }
}

// Middleware de logging
func LoggingMiddleware(
 next graphql.FieldResolver,
) graphql.FieldResolver {
 return func(ctx context.Context, field *graphql.Field) (interface{}, error) {
  log.Printf("Resolving field: %s.%s", field.ObjectName, field.Name)
  return next(ctx, field)
 }
}

// Middleware de caché
type CacheMiddleware struct {
 cache map[string]interface{}
 ttl   map[string]time.Time
}

func (cm *CacheMiddleware) Middleware(
 next graphql.FieldResolver,
) graphql.FieldResolver {
 return func(ctx context.Context, field *graphql.Field) (interface{}, error) {
  // Generar clave de caché
  key := field.Name

  // Verificar caché
  if cached, ok := cm.cache[key]; ok {
   if expiry, ok := cm.ttl[key]; ok && time.Now().Before(expiry) {
    log.Printf("Cache HIT: %s", key)
    return cached, nil
   }
  }

  // Ejecutar resolver
  result, err := next(ctx, field)

  // Cachear resultado
  cm.cache[key] = result
  cm.ttl[key] = time.Now().Add(5 * time.Minute)

  return result, err
 }
}

// Aplicar middleware
func (r *Resolver) User(ctx context.Context, id string) (*User, error) {
 // El middleware se aplica automáticamente en gqlgen
 // mediante hooks de ejecución
 user := &User{}
 // ... lógica de búsqueda
 return user, nil
}
```

### Error Handling en Resolvers

```go
// error_handling.go
package main

import (
 "context"
 "errors"
 "github.com/graphql-go/graphql"
)

// Custom error
type ResolverError struct {
 Code    string `json:"code"`
 Message string `json:"message"`
 Field   string `json:"field"`
}

// Error handler
func (r *QueryResolver) User(
 ctx context.Context,
 id string,
) (*User, error) {
 if id == "" {
  return nil, &graphql.Error{
   Message: "User ID cannot be empty",
   Extensions: map[string]interface{}{
    "code": "INVALID_ID",
    "field": "id",
   },
  }
 }

 user := &User{}
 err := r.db.QueryRowContext(ctx,
  "SELECT id, name FROM users WHERE id = ?",
  id,
 ).Scan(&user.ID, &user.Name)

 if err == sql.ErrNoRows {
  return nil, &graphql.Error{
   Message: "User not found",
   Extensions: map[string]interface{}{
    "code": "NOT_FOUND",
   },
  }
 }

 if err != nil {
  return nil, &graphql.Error{
   Message: "Database error",
   Extensions: map[string]interface{}{
    "code": "DATABASE_ERROR",
    "original": err.Error(),
   },
  }
 }

 return user, nil
}
```

---

## 42.7 Framework gqlgen

### Introducción a gqlgen

gqlgen es el framework recomendado para GraphQL en Go. Genera código automáticamente a partir del schema SDL:

```bash
# Instalar gqlgen
go install github.com/99designs/gqlgen@latest

# Inicializar proyecto
gqlgen init
```

### Estructura del proyecto

```
project/
├── graph/
│   ├── schema.graphqls      # Schema SDL
│   ├── generated.go         # Generado automáticamente
│   └── resolver.go          # Tus resolvers
├── server.go                # Servidor HTTP
└── go.mod
```

### Definir Schema

```graphql
# graph/schema.graphqls
scalar DateTime

type User {
  id: ID!
  name: String!
  email: String!
  posts: [Post!]!
  createdAt: DateTime!
}

type Post {
  id: ID!
  title: String!
  content: String!
  author: User!
  createdAt: DateTime!
}

input CreateUserInput {
  name: String!
  email: String!
}

input CreatePostInput {
  title: String!
  content: String!
  authorId: ID!
}

type Query {
  user(id: ID!): User
  posts(first: Int): [Post!]!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  createPost(input: CreatePostInput!): Post!
}
```

### Generar código

```bash
# Generar basado en schema
go run github.com/99designs/gqlgen generate

# Resultado: graph/generated.go con tipos y stubs
```

### Implementar Resolvers

```go
// graph/resolver.go
package graph

import (
 "context"
 "time"
)

type Resolver struct {
 db *sql.DB
}

// User query
func (r *queryResolver) User(
 ctx context.Context,
 id string,
) (*User, error) {
 user := &User{}

 err := r.db.QueryRowContext(ctx,
  "SELECT id, name, email, created_at FROM users WHERE id = ?",
  id,
 ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

 return user, err
}

// Posts query
func (r *queryResolver) Posts(
 ctx context.Context,
 first *int,
) ([]*Post, error) {

 limit := 10
 if first != nil && *first > 0 {
  limit = *first
 }

 rows, err := r.db.QueryContext(ctx,
  "SELECT id, title, content, author_id, created_at FROM posts LIMIT ?",
  limit,
 )
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 var posts []*Post
 for rows.Next() {
  post := &Post{}
  err := rows.Scan(&post.ID, &post.Title, &post.Content,
   &post.AuthorID, &post.CreatedAt)
  if err != nil {
   return nil, err
  }
  posts = append(posts, post)
 }

 return posts, nil
}

// Mutation: CreateUser
func (r *mutationResolver) CreateUser(
 ctx context.Context,
 input CreateUserInput,
) (*User, error) {

 user := &User{
  ID:        uuid.New().String(),
  Name:      input.Name,
  Email:     input.Email,
  CreatedAt: time.Now(),
 }

 _, err := r.db.ExecContext(ctx,
  "INSERT INTO users (id, name, email, created_at) VALUES (?, ?, ?, ?)",
  user.ID, user.Name, user.Email, user.CreatedAt,
 )

 return user, err
}

// User.posts - Resolver para relación
func (r *userResolver) Posts(
 ctx context.Context,
 obj *User,
) ([]*Post, error) {

 rows, err := r.db.QueryContext(ctx,
  "SELECT id, title, content, author_id, created_at FROM posts WHERE author_id = ?",
  obj.ID,
 )
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 var posts []*Post
 for rows.Next() {
  post := &Post{}
  err := rows.Scan(&post.ID, &post.Title, &post.Content,
   &post.AuthorID, &post.CreatedAt)
  if err != nil {
   return nil, err
  }
  posts = append(posts, post)
 }

 return posts, nil
}

// Post.author - Resolver para relación inversa
func (r *postResolver) Author(
 ctx context.Context,
 obj *Post,
) (*User, error) {

 user := &User{}
 err := r.db.QueryRowContext(ctx,
  "SELECT id, name, email, created_at FROM users WHERE id = ?",
  obj.AuthorID,
 ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

 return user, err
}
```

### Servidor gqlgen

```go
// server.go
package main

import (
 "database/sql"
 "log"
 "net/http"
 "os"

 "github.com/99designs/gqlgen/graphql/handler"
 "github.com/99designs/gqlgen/graphql/playground"
 "yourmodule/graph"
 _ "github.com/mattn/go-sqlite3"
)

func main() {
 // Conectar a BD
 db, err := sql.Open("sqlite3", "test.db")
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 // Crear resolver
 resolver := &graph.Resolver{db: db}

 // Crear servidor GraphQL
 srv := handler.NewDefaultServer(
  graph.NewExecutableSchema(
   graph.Config{Resolvers: resolver},
  ),
 )

 // Rutas
 http.Handle("/query", srv)
 http.Handle("/", playground.Handler("GraphQL playground", "/query"))

 log.Println("Servidor en http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 42.8 Validación

### Input Validation

```go
// validation.go
package main

import (
 "errors"
 "regexp"
 "strings"
 "unicode"
)

// ValidateCreateUserInput valida entrada de usuario
func ValidateCreateUserInput(input CreateUserInput) error {
 // Nombre no vacío
 if strings.TrimSpace(input.Name) == "" {
  return errors.New("name cannot be empty")
 }

 // Nombre longitud
 if len(input.Name) < 2 || len(input.Name) > 100 {
  return errors.New("name must be 2-100 characters")
 }

 // Email válido
 if !isValidEmail(input.Email) {
  return errors.New("invalid email format")
 }

 // Password strength
 if !isStrongPassword(input.Password) {
  return errors.New("password must have uppercase, lowercase, number, special char")
 }

 return nil
}

func isValidEmail(email string) bool {
 re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
 return re.MatchString(email)
}

func isStrongPassword(password string) bool {
 if len(password) < 8 {
  return false
 }

 hasUpper := false
 hasLower := false
 hasNumber := false
 hasSpecial := false

 for _, r := range password {
  if unicode.IsUpper(r) {
   hasUpper = true
  } else if unicode.IsLower(r) {
   hasLower = true
  } else if unicode.IsDigit(r) {
   hasNumber = true
  } else if unicode.IsPunct(r) || unicode.IsSymbol(r) {
   hasSpecial = true
  }
 }

 return hasUpper && hasLower && hasNumber && hasSpecial
}

// En resolver
func (r *mutationResolver) CreateUser(
 ctx context.Context,
 input CreateUserInput,
) (*User, error) {

 // Validar entrada
 if err := ValidateCreateUserInput(input); err != nil {
  return nil, &graphql.Error{
   Message: err.Error(),
   Extensions: map[string]interface{}{
    "code": "VALIDATION_ERROR",
   },
  }
 }

 // ... crear usuario
 return user, nil
}
```

### Custom Validation Rules

```go
// custom_rules.go
package main

import (
 "context"
 "fmt"
)

type ValidationRule interface {
 Validate(value interface{}) error
}

// Rule: única email
type UniqueEmailRule struct {
 db *sql.DB
}

func (r *UniqueEmailRule) Validate(value interface{}) error {
 email := value.(string)

 var count int
 err := r.db.QueryRow(
  "SELECT COUNT(*) FROM users WHERE email = ?",
  email,
 ).Scan(&count)

 if err != nil {
  return err
 }

 if count > 0 {
  return fmt.Errorf("email already exists")
 }

 return nil
}

// Rule: usuario existe
type UserExistsRule struct {
 db *sql.DB
}

func (r *UserExistsRule) Validate(value interface{}) error {
 userID := value.(string)

 var id string
 err := r.db.QueryRow(
  "SELECT id FROM users WHERE id = ?",
  userID,
 ).Scan(&id)

 if err == sql.ErrNoRows {
  return fmt.Errorf("user not found")
 }

 return err
}

// Rule: post pertenece a usuario
type PostOwnershipRule struct {
 db *sql.DB
}

func (r *PostOwnershipRule) Validate(
 ctx context.Context,
 postID string,
 userID string,
) error {

 var authorID string
 err := r.db.QueryRowContext(ctx,
  "SELECT author_id FROM posts WHERE id = ?",
  postID,
 ).Scan(&authorID)

 if err == sql.ErrNoRows {
  return fmt.Errorf("post not found")
 }

 if authorID != userID {
  return fmt.Errorf("not authorized to modify this post")
 }

 return nil
}
```

### Error Responses

```go
// error_response.go
package main

import (
 "github.com/graphql-go/graphql"
)

type ErrorResponse struct {
 Code    string      `json:"code"`
 Message string      `json:"message"`
 Details interface{} `json:"details,omitempty"`
}

// GraphQL error formatter
func ErrorFormatter(err error) *graphql.Error {
 var code string
 var details interface{}

 switch err {
 case ErrNotFound:
  code = "NOT_FOUND"
 case ErrUnauthorized:
  code = "UNAUTHORIZED"
 case ErrValidation:
  code = "VALIDATION_ERROR"
 default:
  code = "INTERNAL_ERROR"
 }

 return &graphql.Error{
  Message: err.Error(),
  Extensions: map[string]interface{}{
   "code":    code,
   "details": details,
  },
 }
}

// Batch errors
type BatchErrors struct {
 errors []error
}

func (be *BatchErrors) Add(err error) {
 be.errors = append(be.errors, err)
}

func (be *BatchErrors) HasErrors() bool {
 return len(be.errors) > 0
}

func (be *BatchErrors) Error() string {
 msg := ""
 for _, err := range be.errors {
  msg += err.Error() + "; "
 }
 return msg
}
```

---

## 42.9 Performance

### N+1 Problem

El problema N+1 ocurre cuando resolver un campo requiere queries adicionales:

```go
// ❌ Problema: N+1 queries
func (r *userResolver) Posts(
 ctx context.Context,
 user *User,
) ([]*Post, error) {

 // 1 query por usuario
 rows, err := r.db.QueryContext(ctx,
  "SELECT id, title FROM posts WHERE author_id = ?",
  user.ID,
 )
 // Si hay 100 usuarios, 100 queries!

 // ...
 return posts, nil
}

// Queries ejecutadas:
// SELECT * FROM users;                           // 1 query
// SELECT * FROM posts WHERE author_id = 'user1' // 100 queries
// SELECT * FROM posts WHERE author_id = 'user2'
// ...
// Total: 101 queries (N+1)
```

### DataLoader para N+1

```go
// dataloader.go
package main

import (
 "context"
 "fmt"
 "strings"

 "github.com/graph-gophers/dataloader/v7"
)

// PostLoader carga múltiples posts de una sola vez
type PostLoader struct {
 db *sql.DB
}

func (pl *PostLoader) LoadAll(
 ctx context.Context,
 keys dataloader.Keys,
) []*dataloader.Result {

 // Convertir keys a strings
 userIDs := make([]string, len(keys))
 for i, key := range keys {
  userIDs[i] = key.String()
 }

 // Una sola query para todos
 placeholders := strings.Repeat("?,", len(userIDs))
 placeholders = placeholders[:len(placeholders)-1]

 query := fmt.Sprintf(
  "SELECT author_id, id, title FROM posts WHERE author_id IN (%s)",
  placeholders,
 )

 args := make([]interface{}, len(userIDs))
 for i, id := range userIDs {
  args[i] = id
 }

 rows, err := pl.db.QueryContext(ctx, query, args...)
 if err != nil {
  return []*dataloader.Result{{Data: nil, Error: err}}
 }
 defer rows.Close()

 // Agrupar posts por usuario
 postsByUser := make(map[string][]*Post)
 for rows.Next() {
  var authorID, id, title string
  if err := rows.Scan(&authorID, &id, &title); err != nil {
   continue
  }

  post := &Post{ID: id, Title: title, AuthorID: authorID}
  postsByUser[authorID] = append(postsByUser[authorID], post)
 }

 // Retornar en orden original
 results := make([]*dataloader.Result, len(keys))
 for i, key := range keys {
  userID := key.String()
  results[i] = &dataloader.Result{
   Data:  postsByUser[userID],
   Error: nil,
  }
 }

 return results
}

// Usar en resolver
func (r *userResolver) Posts(
 ctx context.Context,
 user *User,
) ([]*Post, error) {

 // Obtener loader del contexto
 loader := ctx.Value("post_loader").(*dataloader.Loader)

 // Queue la carga
 thunk := loader.Load(ctx, dataloader.StringKey(user.ID))

 // Wait para resultado
 result, err := thunk()
 if err != nil {
  return nil, err
 }

 return result.([]*Post), nil
}

// Setup en servidor
func setupDataloaders(db *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  ctx := context.WithValue(
   r.Context(),
   "post_loader",
   dataloader.NewBatchedLoader(
    &PostLoader{db: db},
   ),
  )

  // Pasar contexto con loaders a GraphQL handler
 }
}
```

### Caching

```go
// caching.go
package main

import (
 "context"
 "time"

 "github.com/redis/go-redis/v9"
)

type CachedResolver struct {
 db    *sql.DB
 redis *redis.Client
}

// User con caché
func (r *CachedResolver) User(
 ctx context.Context,
 id string,
) (*User, error) {

 // Intenta caché
 cacheKey := fmt.Sprintf("user:%s", id)
 val, err := r.redis.Get(ctx, cacheKey).Result()

 if err == nil {
  // Cache HIT
  user := &User{}
  json.Unmarshal([]byte(val), user)
  return user, nil
 }

 // Cache MISS - query BD
 user := &User{}
 err = r.db.QueryRowContext(ctx,
  "SELECT id, name, email FROM users WHERE id = ?",
  id,
 ).Scan(&user.ID, &user.Name, &user.Email)

 if err != nil {
  return nil, err
 }

 // Cachear por 1 hora
 data, _ := json.Marshal(user)
 r.redis.Set(ctx, cacheKey, data, 1*time.Hour)

 return user, nil
}

// Invalidar caché después de mutation
func (r *CachedResolver) UpdateUser(
 ctx context.Context,
 id string,
 input UpdateUserInput,
) (*User, error) {

 user := &User{ID: id}
 // ... actualizar en BD

 // Invalidar caché
 r.redis.Del(ctx, fmt.Sprintf("user:%s", id))

 return user, nil
}

// Caché estratégico para queries pesadas
func (r *CachedResolver) Posts(
 ctx context.Context,
 filter *PostFilter,
) ([]*Post, error) {

 // Generar clave de caché basada en filtro
 cacheKey := generateCacheKey("posts", filter)

 val, err := r.redis.Get(ctx, cacheKey).Result()
 if err == nil {
  var posts []*Post
  json.Unmarshal([]byte(val), &posts)
  return posts, nil
 }

 // Query y cachear
 posts, err := r.fetchPostsFromDB(ctx, filter)
 if err != nil {
  return nil, err
 }

 data, _ := json.Marshal(posts)
 r.redis.Set(ctx, cacheKey, data, 5*time.Minute)

 return posts, nil
}
```

---

## 42.10 Security

### Authentication

```go
// auth.go
package main

import (
 "context"
 "errors"
 "fmt"
 "strings"

 "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
 UserID string `json:"user_id"`
 Email  string `json:"email"`
 Role   string `json:"role"`
 jwt.RegisteredClaims
}

// Middleware de autenticación
func AuthMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

  // Obtener token del header
  authHeader := r.Header.Get("Authorization")
  if authHeader == "" {
   next.ServeHTTP(w, r)
   return
  }

  // Parse "Bearer <token>"
  parts := strings.Split(authHeader, " ")
  if len(parts) != 2 || parts[0] != "Bearer" {
   http.Error(w, "Invalid auth header", http.StatusUnauthorized)
   return
  }

  token, err := jwt.ParseWithClaims(parts[1], &Claims{},
   func(token *jwt.Token) (interface{}, error) {
    // Verificar algoritmo
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
     return nil, errors.New("invalid signing method")
    }
    return []byte("secret_key"), nil
   },
  )

  if err != nil || !token.Valid {
   http.Error(w, "Invalid token", http.StatusUnauthorized)
   return
  }

  claims := token.Claims.(*Claims)

  // Pasar usuario al contexto
  ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
  ctx = context.WithValue(ctx, "user_role", claims.Role)

  next.ServeHTTP(w, r.WithContext(ctx))
 })
}

// Obtener usuario del contexto
func GetUserID(ctx context.Context) (string, error) {
 userID, ok := ctx.Value("user_id").(string)
 if !ok {
  return "", errors.New("unauthorized")
 }
 return userID, nil
}

func GetUserRole(ctx context.Context) string {
 role, ok := ctx.Value("user_role").(string)
 if !ok {
  return "guest"
 }
 return role
}
```

### Authorization

```go
// authorization.go
package main

import (
 "context"
 "errors"
)

type Permission struct {
 Resource string
 Action   string
}

var rolePermissions = map[string][]Permission{
 "ADMIN": {
  {Resource: "users", Action: "read"},
  {Resource: "users", Action: "create"},
  {Resource: "users", Action: "update"},
  {Resource: "users", Action: "delete"},
  {Resource: "posts", Action: "read"},
  {Resource: "posts", Action: "delete"},
 },
 "USER": {
  {Resource: "posts", Action: "read"},
  {Resource: "posts", Action: "create"},
  {Resource: "comments", Action: "read"},
  {Resource: "comments", Action: "create"},
 },
 "GUEST": {
  {Resource: "posts", Action: "read"},
 },
}

// Verificar permisos
func CheckPermission(ctx context.Context, resource, action string) error {
 role := GetUserRole(ctx)

 permissions := rolePermissions[role]
 for _, p := range permissions {
  if p.Resource == resource && p.Action == action {
   return nil
  }
 }

 return errors.New("permission denied")
}

// Usar en resolver
func (r *mutationResolver) CreatePost(
 ctx context.Context,
 input CreatePostInput,
) (*Post, error) {

 if err := CheckPermission(ctx, "posts", "create"); err != nil {
  return nil, err
 }

 // ... crear post
 return post, nil
}

// Autorización a nivel de objeto
func (r *mutationResolver) DeletePost(
 ctx context.Context,
 id string,
) (bool, error) {

 userID, err := GetUserID(ctx)
 if err != nil {
  return false, err
 }

 // Verificar que es propietario
 var authorID string
 err = r.db.QueryRowContext(ctx,
  "SELECT author_id FROM posts WHERE id = ?",
  id,
 ).Scan(&authorID)

 if err != nil {
  return false, err
 }

 if authorID != userID && GetUserRole(ctx) != "ADMIN" {
  return false, errors.New("not authorized")
 }

 _, err = r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", id)
 return err == nil, err
}
```

### Rate Limiting

```go
// rate_limit.go
package main

import (
 "context"
 "fmt"
 "net/http"
 "time"

 "github.com/redis/go-redis/v9"
)

type RateLimiter struct {
 redis *redis.Client
 limit int
 ttl   time.Duration
}

func NewRateLimiter(redis *redis.Client) *RateLimiter {
 return &RateLimiter{
  redis: redis,
  limit: 100,      // 100 requests
  ttl:   time.Minute, // per minute
 }
}

// Middleware GraphQL rate limit
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

  // Obtener identificador del cliente
  clientID := getClientID(r)
  key := fmt.Sprintf("ratelimit:%s", clientID)

  // Incrementar contador
  count, err := rl.redis.Incr(r.Context(), key).Result()
  if err != nil {
   http.Error(w, "Rate limiter error", http.StatusInternalServerError)
   return
  }

  // Establecer TTL en primer request
  if count == 1 {
   rl.redis.Expire(r.Context(), key, rl.ttl)
  }

  // Verificar límite
  if count > int64(rl.limit) {
   w.Header().Set("Retry-After", fmt.Sprintf("%d", rl.ttl.Seconds()))
   http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
   return
  }

  // Headers informativos
  w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
  w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.limit-int(count)))

  next.ServeHTTP(w, r)
 })
}

func getClientID(r *http.Request) string {
 // Usar token si está autenticado
 if userID, ok := r.Context().Value("user_id").(string); ok {
  return userID
 }

 // Sino, usar IP
 return r.RemoteAddr
}
```

---

## 42.11 Buenas Prácticas y Patterns

### Schema Design

```graphql
# ✅ Buen diseño
type User implements Node {
  id: ID!
  name: String!
  email: String!
  bio: String
  posts(first: Int, after: String): PostConnection!
  followers(first: Int): UserConnection!
  createdAt: DateTime!
}

type PostConnection {
  edges: [PostEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type PostEdge {
  node: Post!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

# ❌ Malo: Sin paginación
type User {
  id: ID!
  allPosts: [Post!]!  # Puede ser miles!
}
```

### Relay Cursor Connections

Patrón estándar para paginación:

```go
// relay.go
package main

import (
 "encoding/base64"
 "fmt"
 "strconv"
 "strings"
)

type PageInfo struct {
 HasNextPage     bool    `json:"hasNextPage"`
 HasPreviousPage bool    `json:"hasPreviousPage"`
 StartCursor     *string `json:"startCursor"`
 EndCursor       *string `json:"endCursor"`
}

type Edge struct {
 Node   interface{} `json:"node"`
 Cursor string      `json:"cursor"`
}

type Connection struct {
 Edges    []*Edge   `json:"edges"`
 PageInfo *PageInfo `json:"pageInfo"`
}

// Codificar cursor (offset:count)
func EncodeCursor(offset, count int) string {
 data := fmt.Sprintf("%d:%d", offset, count)
 return base64.StdEncoding.EncodeToString([]byte(data))
}

// Decodificar cursor
func DecodeCursor(cursor string) (offset, count int, err error) {
 data, err := base64.StdEncoding.DecodeString(cursor)
 if err != nil {
  return
 }

 parts := strings.Split(string(data), ":")
 if len(parts) != 2 {
  return 0, 0, fmt.Errorf("invalid cursor format")
 }

 offset, _ = strconv.Atoi(parts[0])
 count, _ = strconv.Atoi(parts[1])
 return
}

// Construir conexión
func BuildConnection(
 totalCount int,
 items []interface{},
 offset int,
 first int,
) *Connection {

 edges := []*Edge{}
 for i, item := range items {
  cursor := EncodeCursor(offset+i, 1)
  edges = append(edges, &Edge{
   Node:   item,
   Cursor: cursor,
  })
 }

 var startCursor, endCursor *string
 if len(edges) > 0 {
  startCursor = &edges[0].Cursor
  endCursor = &edges[len(edges)-1].Cursor
 }

 return &Connection{
  Edges: edges,
  PageInfo: &PageInfo{
   HasNextPage:     offset+first < totalCount,
   HasPreviousPage: offset > 0,
   StartCursor:     startCursor,
   EndCursor:       endCursor,
  },
 }
}
```

### Versioning

```graphql
# Estrategia 1: Schema versioning
type Query {
  # Deprecated, usar postsV2
  posts: [Post!]! @deprecated(reason: "Use postsV2 instead")
  postsV2(first: Int): PostConnection!
}

# Estrategia 2: Field versioning
type User {
  id: ID!
  name: String!
  email: String! @deprecated(reason: "Use contact.email")
  contact: Contact!
}

type Contact {
  email: String!
  phone: String
}

# Estrategia 3: Enum versionining para cambios menores
enum PostStatus {
  DRAFT
  PUBLISHED
  ARCHIVED
  DELETED @deprecated(reason: "Use ARCHIVED")
}
```

### Testing GraphQL

```go
// graphql_test.go
package main

import (
 "context"
 "testing"

 "github.com/stretchr/testify/assert"
 "github.com/99designs/gqlgen/client"
 "github.com/99designs/gqlgen/graphql/handler"
)

type GetUserResponse struct {
 User *User `json:"user"`
}

func TestQueryUser(t *testing.T) {
 c := client.New(
  handler.NewDefaultServer(
   NewExecutableSchema(Config{
    Resolvers: &mockResolver{},
   }),
  ),
 )

 var resp GetUserResponse

 c.MustPost(`
  query {
   user(id: "1") {
    id
    name
    email
   }
  }
 `, &resp)

 assert.Equal(t, "1", resp.User.ID)
 assert.Equal(t, "Alice", resp.User.Name)
}

func TestCreateUserMutation(t *testing.T) {
 c := client.New(
  handler.NewDefaultServer(
   NewExecutableSchema(Config{
    Resolvers: &mockResolver{},
   }),
  ),
 )

 var resp struct {
  CreateUser *User `json:"createUser"`
 }

 c.MustPost(`
  mutation {
   createUser(input: {name: "Bob", email: "bob@example.com"}) {
    id
    name
    email
   }
  }
 `, &resp)

 assert.NotEmpty(t, resp.CreateUser.ID)
 assert.Equal(t, "Bob", resp.CreateUser.Name)
}

type mockResolver struct{}

func (r *mockResolver) Query() QueryResolver {
 return &mockQueryResolver{}
}

type mockQueryResolver struct{}

func (qr *mockQueryResolver) User(ctx context.Context, id string) (*User, error) {
 return &User{
  ID:    id,
  Name:  "Alice",
  Email: "alice@example.com",
 }, nil
}
```

### Antipatterns a Evitar

```graphql
# ❌ ANTIPATTERN 1: Campos sin limitación de tamaño
type Query {
  allUsers: [User!]!  # Puede retornar millones!
  allPosts: [Post!]!  # Sin paginación = problema
}

# ✅ Correcto: Paginación requerida
type Query {
  users(first: Int!, after: String): UserConnection!
  posts(first: Int!, after: String): PostConnection!
}

# ❌ ANTIPATTERN 2: Relaciones circulares sin límite
type User {
  posts: [Post!]!
}

type Post {
  author: User!
  comments: [Comment!]!
}

type Comment {
  author: User!
  post: Post!
}

# Query que explota:
query {
  user {
    posts {
      comments {
        author {
          posts {  # Profundidad ilimitada!
            comments {
              author { ... }
            }
          }
        }
      }
    }
  }
}

# ✅ Correcto: Query depth limiting
directive @depth(max: Int!) on FIELD_DEFINITION

type User {
  posts: [Post!]! @depth(max: 3)
}

# ❌ ANTIPATTERN 3: Errores sin contexto
{
  "errors": [
    {"message": "error"}  # ¿Qué error? ¿Dónde?
  ]
}

# ✅ Correcto: Errores informativos
{
  "errors": [
    {
      "message": "User not found",
      "extensions": {
        "code": "NOT_FOUND",
        "field": "user.id",
        "providedValue": "invalid-id"
      }
    }
  ]
}

# ❌ ANTIPATTERN 4: Sin validación de entrada
mutation {
  createUser(input: { name: "" }) { id }  # Válido? Inválido?
}

# ✅ Correcto: Input types fuertemente validados
input CreateUserInput {
  name: String!      # No nulo, mínimo 2 chars
  email: String!     # Validar formato email
  password: String!  # Validar fuerza
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Schema Simple - Definir Types Básicos

**Objetivo:** Crear un schema GraphQL para una API de tareas.

**Requisitos:**

- Definir type `Task` con campos: id, title, description, completed, dueDate
- Definir type `Category` con campos: id, name, color
- Crear relación Task -> Category (muchas tareas por categoría)
- Definir Query root con campos: tasks, task(id), categories

**Archivo: `schema-ejercicio1.graphql`**

```graphql
scalar DateTime

type Category {
  id: ID!
  name: String!
  color: String!
  tasks: [Task!]!
}

type Task {
  id: ID!
  title: String!
  description: String
  completed: Boolean!
  dueDate: DateTime
  category: Category!
}

type Query {
  task(id: ID!): Task
  tasks(categoryId: ID, completed: Boolean): [Task!]!
  categories: [Category!]!
}
```

**Solución esperada:** Schema válido compilable con `gqlgen validate schema.graphql`

---

### Ejercicio 2: Queries y Mutations - Operaciones CRUD

**Objetivo:** Implementar operaciones CRUD completas.

**Requisitos:**

- Agregar tipos de entrada: CreateTaskInput, UpdateTaskInput
- Implementar mutations: createTask, updateTask, deleteTask, createCategory
- Todos los campos requeridos en inputs
- Validar que IDs existan antes de actualizar/eliminar

**Archivo: `schema-ejercicio2.graphql`**

```graphql
# Completar schema anterior

input CreateTaskInput {
  title: String!
  description: String
  categoryId: ID!
}

input UpdateTaskInput {
  title: String
  description: String
  completed: Boolean
  categoryId: ID
}

type Mutation {
  createTask(input: CreateTaskInput!): Task!
  updateTask(id: ID!, input: UpdateTaskInput!): Task
  deleteTask(id: ID!): Boolean!
  createCategory(name: String!, color: String!): Category!
}
```

**Prueba:**

```
mutation {
  createTask(input: {title: "Tarea 1", categoryId: "cat1"}) {
    id
    title
  }
}
```

---

### Ejercicio 3: Resolvers - Implementar Lógica de Negocio

**Objetivo:** Implementar resolvers en Go para el schema del ejercicio 2.

**Requisitos:**

- Implementar resolver de Query.tasks con filtrado opcional
- Implementar resolver de Mutation.createTask con generación de ID
- Implementar resolver de relación Task.category
- Manejar errores apropiadmente

**Archivo: `resolver-ejercicio3.go`**

```go
// Estructura base
package main

import (
 "context"
 "errors"
 "github.com/google/uuid"
)

var tasks = make(map[string]*Task)
var categories = make(map[string]*Category)

type queryResolver struct{}

func (r *queryResolver) Tasks(
 ctx context.Context,
 categoryID *string,
 completed *bool,
) ([]*Task, error) {

 var result []*Task

 for _, task := range tasks {
  if categoryID != nil && task.CategoryID != *categoryID {
   continue
  }
  if completed != nil && task.Completed != *completed {
   continue
  }
  result = append(result, task)
 }

 return result, nil
}

type mutationResolver struct{}

func (r *mutationResolver) CreateTask(
 ctx context.Context,
 input CreateTaskInput,
) (*Task, error) {

 // Validar que categoría existe
 if _, ok := categories[input.CategoryID]; !ok {
  return nil, errors.New("category not found")
 }

 task := &Task{
  ID:         uuid.New().String(),
  Title:      input.Title,
  Description: input.Description,
  Completed:  false,
  CategoryID: input.CategoryID,
 }

 tasks[task.ID] = task
 return task, nil
}
```

---

### Ejercicio 4: Subscriptions - Actualizaciones en Tiempo Real

**Objetivo:** Agregar subscriptions para notificaciones en tiempo real.

**Requisitos:**

- Definir Subscription con: taskCreated, taskUpdated(id), taskDeleted
- Implementar broadcast de eventos cuando se crea/actualiza/elimina tarea
- Usar canales de Go para comunicación

**Archivo: `subscription-ejercicio4.go`**

```go
// Estructura para broadcast
package main

import "context"

type EventBroker struct {
 subscribers map[string][]chan *Task
}

func (eb *EventBroker) Subscribe(ctx context.Context, eventType string) chan *Task {
 ch := make(chan *Task, 10)
 eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
 return ch
}

func (eb *EventBroker) Publish(eventType string, task *Task) {
 for _, ch := range eb.subscribers[eventType] {
  select {
  case ch <- task:
  default:
   // Buffer lleno, ignorar
  }
 }
}

type subscriptionResolver struct {
 broker *EventBroker
}

func (r *subscriptionResolver) TaskCreated(
 ctx context.Context,
) (<-chan *Task, error) {
 return r.broker.Subscribe(ctx, "task_created"), nil
}

// Modificar mutationResolver para publicar eventos
func (r *mutationResolver) CreateTask(
 ctx context.Context,
 input CreateTaskInput,
) (*Task, error) {

 task := &Task{...}

 // Publicar evento
 r.broker.Publish("task_created", task)

 return task, nil
}
```

---

### Ejercicio 5: Full API - GraphQL Completo con Validación

**Objetivo:** Crear una API GraphQL completa y funcional.

**Requisitos:**

- Usar gqlgen con base de datos SQLite
- Implementar validación de entrada (título no vacío, fecha válida)
- Agregar autenticación (token JWT simple)
- Incluir paginación con relay cursors
- Rate limiting simple

**Archivo: `main-ejercicio5.go`**

```go
package main

import (
 "database/sql"
 "log"
 "net/http"
 "os"

 "github.com/99designs/gqlgen/graphql/handler"
 "github.com/99designs/gqlgen/graphql/playground"
 _ "github.com/mattn/go-sqlite3"
)

func main() {
 // Conectar a BD
 db, err := sql.Open("sqlite3", "tasks.db")
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 // Crear tablas
 createTables(db)

 // Resolver con BD
 resolver := &Resolver{db: db}

 // Handler GraphQL
 srv := handler.NewDefaultServer(
  NewExecutableSchema(Config{Resolvers: resolver}),
 )

 // Middleware
 http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
  // Auth check
  if token := r.Header.Get("Authorization"); token == "" {
   http.Error(w, "Missing auth token", http.StatusUnauthorized)
   return
  }

  srv.ServeHTTP(w, r)
 })

 http.Handle("/", playground.Handler("Tasks API", "/query"))

 log.Println("Server running on :8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTables(db *sql.DB) {
 schema := `
 CREATE TABLE IF NOT EXISTS categories (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  color TEXT NOT NULL
 );

 CREATE TABLE IF NOT EXISTS tasks (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  completed BOOLEAN DEFAULT 0,
  category_id TEXT NOT NULL,
  due_date TEXT,
  FOREIGN KEY (category_id) REFERENCES categories(id)
 );
 `

 if _, err := db.Exec(schema); err != nil {
  log.Fatal(err)
 }
}

// Validation functions
func validateTaskInput(input CreateTaskInput) error {
 if input.Title == "" {
  return errors.New("title cannot be empty")
 }
 if len(input.Title) > 200 {
  return errors.New("title too long")
 }
 return nil
}
```

---

## Resumen de Conceptos Clave

| Concepto | Descripción |
|----------|-------------|
| **Schema** | Definición de tipos, queries y mutations |
| **Query** | Lectura de datos (GET) |
| **Mutation** | Escritura de datos (POST/PUT/DELETE) |
| **Subscription** | Actualizaciones en tiempo real (WebSocket) |
| **Resolver** | Función que resuelve un campo |
| **Scalar** | Tipos primitivos (String, Int, Float, Boolean, ID) |
| **Input Type** | Estructura para parámetros de mutaciones |
| **Interface** | Contrato que múltiples tipos deben implementar |
| **Union** | Tipo que puede ser uno de varios tipos |
| **DataLoader** | Solución para N+1 problem |
| **Relay** | Patrón estándar para paginación |

---

## Comparación: REST vs GraphQL en Go

```go
// REST
GET /api/users/123                 // Obtener usuario
GET /api/users/123/posts           // Obtener posts
GET /api/posts/1/comments          // N+1!
Total: N+1 requests

// GraphQL
query {
  user(id: "123") {
    name
    posts {
      title
      comments {
        text
      }
    }
  }
}
Total: 1 request
```

---

## Herramientas Recomendadas

1. **gqlgen** - Generación de código
2. **go-graphql-playground** - IDE interactivo
3. **dataloader** - Prevenir N+1 queries
4. **redis** - Caché
5. **golang-jwt** - Autenticación

---

## Referencias Útiles

- [GraphQL oficial](https://graphql.org)
- [gqlgen documentation](https://gqlgen.com)
- [GraphQL best practices](https://graphql.org/learn/best-practices)
- [How GraphQL will Shape the Future of APIs](https://www.apollographql.com)

---

**Capítulo completado.** Este capítulo ha cubierto GraphQL comprehensivamente: desde conceptos básicos hasta implementación en producción con Go, incluyendo patterns, security y performance optimization.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/42-graphql/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/42-graphql):

```bash
cd examples/42-graphql
go run .
```
