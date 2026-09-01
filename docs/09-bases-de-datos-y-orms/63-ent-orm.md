# Capítulo 63: Ent - Framework de entidades

## Índice de Contenidos

1. [Introducción a Ent](#631-introducción-a-ent)
2. [Schema Definition](#632-schema-definition)
3. [Relationships](#633-relationships)
4. [Code Generation](#634-code-generation)
5. [CRUD Operations](#635-crud-operations)
6. [Advanced Querying](#636-advanced-querying)
7. [Hooks & Interceptors](#637-hooks--interceptors)
8. [Transactions](#638-transactions)
9. [Advanced Features](#639-advanced-features)
10. [Testing & Integration](#6310-testing--integration)
11. [Production & Case Studies](#6311-production--case-studies)

---

## 63.1 Introducción a Ent

### 63.1.1 ¿Qué es Ent?

**Ent** (Entity Framework) es un framework ORM moderno y potente para Go que utiliza un modelo **basado en grafos** para definir y gestionar esquemas de datos. A diferencia de otros ORMs tradicionales, Ent genera código type-safe completamente compilable que elimina reflexión en tiempo de ejecución.

```
┌─────────────────────────────────────────────────────────┐
│           Arquitectura de Ent Framework                 │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Schema Definition (Graph-based)                        │
│         ↓                                               │
│  Code Generation (Compile-time)                         │
│         ↓                                               │
│  Type-Safe Builders & Queries                           │
│         ↓                                               │
│  Database Execution (sqlc-like)                         │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Características principales:**

- **Type-safe por compilación**: Los errores de tipo se detectan en compilación
- **Graph-based schema**: Relaciones como nodos y aristas en un grafo
- **Code generation**: Genera todo el código necesario automáticamente
- **Zero reflection**: Sin reflexión en tiempo de ejecución
- **Builder pattern**: API fluida y expresiva
- **Cross-database**: PostgreSQL, MySQL, SQLite, MariaDB

### 63.1.2 Conceptos Fundamentales

#### Entity (Entidad)

Una **Entity** es una unidad de datos fundamental en Ent. Representa un nodo en el grafo de datos.

```
Entity: User
├── Fields: id, name, email, age
├── Edges: posts, comments
└── Constraints: email unique, age >= 0
```

#### Edge (Arista)

Una **Edge** define una relación entre entidades, formando las conexiones del grafo.

```
Edge: author (User -> Post)
├── Type: One-to-Many
├── Direction: Bidirectional
└── Constraints: onDelete=CASCADE
```

#### Field Types

Ent soporta tipos nativos con validación integrada:

```go
// Tipos básicos
Field().Int()          // int
Field().String()       // string
Field().Bool()         // bool
Field().Float()        // float64
Field().Time()         // time.Time
Field().UUID()         // uuid.UUID
Field().JSON()         // json.RawMessage
Field().Enum()         // Enumeración
Field().Bytes()        // []byte
```

### 63.1.3 Comparativa: Ent vs GORM vs sqlc

| Aspecto | Ent | GORM | sqlc |
|---------|-----|------|------|
| **Type Safety** | ⭐⭐⭐ Compilación | ⭐⭐ Runtime | ⭐⭐⭐ Compilación |
| **Learning Curve** | ⭐⭐⭐ Pronunciada | ⭐⭐ Suave | ⭐⭐⭐ Pronunciada |
| **API Design** | Builder fluido | Métodos encadenables | SQL raw |
| **Reflexión** | ❌ Ninguna | ⭐⭐⭐ Mucha | ❌ Ninguna |
| **Schema Validation** | ⭐⭐⭐ Integrada | ⭐⭐ Tags | ❌ Manual |
| **Migrations** | ⭐⭐⭐ Automáticas | ⭐⭐ Manual | ❌ Manual |
| **Hooks/Interceptors** | ⭐⭐⭐ Completos | ⭐⭐ Limitados | ❌ Ninguno |
| **Performance** | ⭐⭐⭐ Excelente | ⭐⭐ Bueno | ⭐⭐⭐ Excelente |
| **Graph Queries** | ⭐⭐⭐ Nativas | ⭐⭐ Posible | ⭐ SQL puro |
| **Production Ready** | ✅ Sí | ✅ Sí | ✅ Sí |

### 63.1.4 Cuándo Usar Ent

**✅ Usar Ent cuando:**

- Necesitas type-safety en compilación
- Trabajas con relaciones complejas (grafos)
- Requieres validación integrada
- Quieres evitar reflexión en tiempo real
- Tu modelo de datos es complejo y bien estructurado
- Performance crítica es importante

**❌ No usar Ent cuando:**

- Necesitas máxima flexibilidad ad-hoc
- Trabajas con queries muy dinámicas
- El schema cambia constantemente
- Necesitas legacy database mapping
- Performance es menos crítica que rapidez de desarrollo

### 63.1.5 Instalación y Configuración Inicial

```bash
# Crear proyecto Go
mkdir my-ent-app && cd my-ent-app
go mod init my-ent-app

# Instalar Ent
go get -u entgo.io/ent/cmd/ent

# Generar directorio de esquema
go run entgo.io/ent/cmd/ent init --template=dir User Post Comment

# Instalar driver de base de datos (ejemplo: PostgreSQL)
go get -u github.com/lib/pq
```

**Estructura de proyecto recomendada:**

```
my-ent-app/
├── ent/
│   ├── schema/
│   │   ├── user.go
│   │   ├── post.go
│   │   └── comment.go
│   ├── client.go
│   ├── migrate.go
│   └── ... (código generado)
├── main.go
├── go.mod
└── go.sum
```

### 63.1.6 Principios de Diseño de Ent

#### 1. Graph-Based Schema

El modelo de Ent piensa en datos como **nodos (entities) conectados por aristas (edges)**:

```
┌─────────────┐
│    User     │───── author ─────► ┌──────────┐
│  (id, name) │                    │   Post   │
└─────────────┘◄──── comments ─────┤(id, text)│
                                   └──────────┘
```

#### 2. Code Generation Over Configuration

Ent prefiere generar código completo a hacer reflexión:

```go
// ❌ Aproximación GORM (reflexión)
db.Where("age > ?", 18).Find(&users)

// ✅ Aproximación Ent (compilado)
users, _ := client.User.Query().Where(user.AgeGT(18)).All(ctx)
```

#### 3. Type Safety Total

Todo es tipado en compilación:

```go
// ✅ Compilado exitosamente
posts, _ := client.Post.Query().Where(post.TitleContains("Go")).All(ctx)

// ❌ Error de compilación (campo no existe)
posts, _ := client.Post.Query().Where(post.InvalidField("Go")).All(ctx)
```

### 63.1.7 Flujo de Desarrollo Típico

```
1. Definir Schema (Go code)
   └─► ent/schema/user.go

2. Ejecutar Generador
   └─► go generate ./ent

3. Usar Cliente Generado
   └─► client.User.Create()...

4. Migraciones Automáticas
   └─► client.Schema.Create(ctx)

5. Queries Type-Safe
   └─► client.User.Query().Where(...).All(ctx)
```

---

## 63.2 Schema Definition

### 63.2.1 Crear tu Primer Entity

Estructura básica de una entidad Ent:

```go
// ent/schema/user.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
    ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").NotEmpty(),
        field.String("email").Unique(),
        field.Int("age").Min(0).Max(150).Optional(),
    }
}

// Edges of the User.
func (User) Edges() []ent.Edge {
    return nil
}
```

**Generando código:**

```bash
go generate ./ent
```

Esto crea:
- `ent/client.go` - Cliente principal
- `ent/user.go` - Modelo de Usuario
- `ent/user_create.go` - Builder de creación
- `ent/user_query.go` - Builder de queries
- `ent/user_update.go` - Builder de actualización
- `ent/user_delete.go` - Builder de eliminación

### 63.2.2 Tipos de Campos

Ent soporta todos los tipos fundamentales con validadores integrados:

```go
package schema

import (
    "time"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/google/uuid"
)

type Product struct {
    ent.Schema
}

func (Product) Fields() []ent.Field {
    return []ent.Field{
        // String fields
        field.String("name").
            MaxLen(255).
            MinLen(1).
            NotEmpty(),
        
        // String con múltiples validadores
        field.String("slug").
            Unique().
            Immutable(). // No puede cambiar
            Match(regexp.MustCompile("^[a-z0-9-]+$")).
            Default(""),
        
        // Numeric fields
        field.Int("quantity").
            Min(0),
        
        field.Float("price").
            Min(0).
            Positive(),
        
        // Boolean
        field.Bool("active").
            Default(true),
        
        // Time fields
        field.Time("created_at").
            Default(time.Now).
            Immutable(),
        
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now),
        
        // UUID
        field.UUID("id", uuid.UUID{}).
            Default(uuid.New).
            Immutable(),
        
        // Bytes
        field.Bytes("data").
            Optional(),
        
        // Enum
        field.Enum("status").
            NamedValues(
                "active", "ACTIVE",
                "inactive", "INACTIVE",
                "deleted", "DELETED",
            ),
        
        // JSON
        field.JSON("metadata", map[string]interface{}{}).
            Default(map[string]interface{}{}),
    }
}

func (Product) Edges() []ent.Edge {
    return nil
}
```

### 63.2.3 Validadores y Constraints

Ent permite validación declarativa integrada:

```go
package schema

import (
    "regexp"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field {
        field.String("email").
            Unique().
            Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)),
        
        field.String("username").
            Unique().
            MinLen(3).
            MaxLen(32).
            Match(regexp.MustCompile("^[a-zA-Z0-9_]+$")).
            NotEmpty(),
        
        field.Int("age").
            Min(0).
            Max(150).
            Optional(),
        
        field.String("password").
            MinLen(8).
            MaxLen(255).
            Sensitive(), // No aparece en logs
        
        field.String("phone").
            Optional().
            Match(regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,}$`)),
        
        field.Bool("verified").
            Default(false),
    }
}

func (User) Indexes() []ent.Index {
    return []ent.Index{
        // Índice simple
        index.Fields("email"),
        
        // Índice compuesto
        index.Fields("email", "username"),
        
        // Índice único compuesto
        index.Fields("organization_id", "username").Unique(),
    }
}

func (User) Edges() []ent.Edge {
    return nil
}
```

### 63.2.4 Configuración Avanzada de Campos

```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type Article struct {
    ent.Schema
}

func (Article) Fields() []ent.Field {
    return []ent.Field{
        // Campo con valores por defecto
        field.String("title").
            NotEmpty(),
        
        // Campo inmutable (solo lectura después de creación)
        field.Time("published_at").
            Immutable().
            Optional(),
        
        // Campo sensible (no aparece en logs/dumps)
        field.String("api_key").
            Sensitive().
            Optional(),
        
        // Campo comentado
        field.Int("view_count").
            Default(0).
            Comment("Total number of views"),
        
        // Campo con valores por defecto variables
        field.String("slug").
            Default("").
            UpdateDefault(""), // Se resetea en updates
        
        // Campo con conversión de tipos
        field.String("tags").
            Default("[]").
            Optional(),
        
        // Campo con almacenamiento alternativo
        field.String("internal_id").
            StructTag(`db:"internal_id" json:"id"`).
            Optional(),
    }
}

func (Article) Edges() []ent.Edge {
    return nil
}
```

### 63.2.5 Migraciones Automáticas

Ent genera migraciones automáticas basadas en el schema:

```go
// main.go
package main

import (
    "context"
    "log"
    "my-ent-app/ent"
    _ "github.com/lib/pq"
)

func main() {
    client, err := ent.Open("postgres", 
        "host=localhost port=5432 user=postgres password=secret dbname=myapp sslmode=disable")
    if err != nil {
        log.Fatalf("failed opening connection to postgres: %v", err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    // Crear todas las tablas
    if err := client.Schema.Create(ctx); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }
    
    log.Println("Schema created successfully!")
}
```

### 63.2.6 Ejemplo Completo: Schema de Blog

```go
// ent/schema/user.go
package schema

import (
    "time"
    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("email").
            Unique().
            NotEmpty(),
        
        field.String("username").
            Unique().
            MinLen(3).
            MaxLen(32),
        
        field.String("display_name").
            NotEmpty(),
        
        field.String("bio").
            MaxLen(500).
            Optional(),
        
        field.String("avatar_url").
            Optional(),
        
        field.Enum("role").
            NamedValues(
                "admin", "ADMIN",
                "moderator", "MODERATOR",
                "user", "USER",
            ).
            Default("user"),
        
        field.Bool("verified").
            Default(false),
        
        field.Time("created_at").
            Default(time.Now).
            Immutable(),
        
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now),
    }
}

func (User) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("email"),
        index.Fields("username"),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("posts", Post.Type),
        edge.To("comments", Comment.Type),
        edge.To("followers", User.Type).
            From("following"),
    }
}
```

---

## 63.3 Relationships

### 63.3.1 Tipos de Relaciones

Ent soporta cuatro patrones de relaciones principales:

```
┌─────────────────────────────────────────────────────┐
│           Tipos de Relaciones en Ent                │
├─────────────────────────────────────────────────────┤
│                                                     │
│ One-to-One (O2O)                                    │
│ ┌─────────┐  has_profile  ┌────────┐               │
│ │  User   │─────────────► │ Profile│               │
│ └─────────┘               └────────┘               │
│                                                     │
│ One-to-Many (O2M)                                   │
│ ┌─────────┐  has_posts  ┌──────────┐               │
│ │  User   │────────────►│  Post    │               │
│ │   (1)   │             │  (Many)  │               │
│ └─────────┘             └──────────┘               │
│                                                     │
│ Many-to-Many (M2M)                                  │
│ ┌─────────┐  has_tags  ┌──────────┐                │
│ │  Post   │◄───────────┤  Tag     │                │
│ │(Many)   │────────────►│ (Many)   │                │
│ └─────────┘            └──────────┘                │
│                                                     │
│ Self-Reference                                      │
│ ┌──────────┐  followers  ┌──────────┐              │
│ │  User    │◄───────────│  User    │              │
│ │ following│────────────►│          │              │
│ └──────────┘            └──────────┘              │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 63.3.2 One-to-One Relationships

Relación 1:1 entre dos entidades:

```go
// ent/schema/user.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
        field.String("email").Unique(),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        // Un usuario tiene un perfil
        edge.To("profile", Profile.Type).
            Unique().
            Required(). // Relación obligatoria
            StorageKey(edge.Symbol("user_id")),
    }
}

// ent/schema/profile.go
type Profile struct {
    ent.Schema
}

func (Profile) Fields() []ent.Field {
    return []ent.Field{
        field.String("bio"),
        field.String("avatar_url").Optional(),
        field.Int("follower_count").Default(0),
    }
}

func (Profile) Edges() []ent.Edge {
    return []ent.Edge{
        // Relación inversa (back-reference)
        edge.From("user", User.Type).
            Ref("profile").
            Unique().
            Required(),
    }
}
```

**Uso:**

```go
// Crear usuario con perfil
profile := client.Profile.Create().
    SetBio("Software Engineer").
    SetAvatarURL("https://...").
    SaveX(ctx)

user := client.User.Create().
    SetName("Alice").
    SetEmail("alice@example.com").
    SetProfile(profile).
    SaveX(ctx)

// Consultar
profile, _ := user.Profile(ctx)
println(profile.Bio) // Software Engineer

user, _ := profile.User(ctx)
println(user.Email) // alice@example.com
```

### 63.3.3 One-to-Many Relationships

Relación 1:N entre entidades:

```go
// ent/schema/user.go
type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("username").Unique(),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        // Un usuario puede tener muchos posts
        edge.To("posts", Post.Type),
        edge.To("comments", Comment.Type),
    }
}

// ent/schema/post.go
type Post struct {
    ent.Schema
}

func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.String("title").NotEmpty(),
        field.String("content"),
        field.Time("created_at").Default(time.Now),
    }
}

func (Post) Edges() []ent.Edge {
    return []ent.Edge{
        // Un post pertenece a un usuario
        edge.From("author", User.Type).
            Ref("posts").
            Unique().
            Required(),
        
        edge.To("comments", Comment.Type),
    }
}

// ent/schema/comment.go
type Comment struct {
    ent.Schema
}

func (Comment) Fields() []ent.Field {
    return []ent.Field{
        field.String("text").NotEmpty(),
        field.Time("created_at").Default(time.Now),
    }
}

func (Comment) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("author", User.Type).
            Ref("comments").
            Unique().
            Required(),
        
        edge.From("post", Post.Type).
            Ref("comments").
            Unique().
            Required(),
    }
}
```

**Uso:**

```go
// Crear usuario
user := client.User.Create().
    SetUsername("alice").
    SaveX(ctx)

// Crear posts
post1 := client.Post.Create().
    SetTitle("First Post").
    SetContent("Hello World").
    SetAuthor(user).
    SaveX(ctx)

post2 := client.Post.Create().
    SetTitle("Second Post").
    SetContent("Ent is awesome").
    SetAuthor(user).
    SaveX(ctx)

// Consultar todos los posts de un usuario
posts, _ := user.Posts(ctx)
fmt.Printf("User has %d posts\n", len(posts))

// Consultar posts con filtros
posts, _ := client.Post.Query().
    Where(post.HasAuthorWith(user.ID(user.ID))).
    All(ctx)
```

### 63.3.4 Many-to-Many Relationships

Relación M:N entre entidades:

```go
// ent/schema/post.go
type Post struct {
    ent.Schema
}

func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.String("title"),
        field.String("content"),
    }
}

func (Post) Edges() []ent.Edge {
    return []ent.Edge{
        // Un post puede tener muchas etiquetas
        edge.To("tags", Tag.Type),
    }
}

// ent/schema/tag.go
type Tag struct {
    ent.Schema
}

func (Tag) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Unique(),
        field.String("slug").Unique(),
    }
}

func (Tag) Edges() []ent.Edge {
    return []ent.Edge{
        // Una etiqueta está en muchos posts
        edge.From("posts", Post.Type).
            Ref("tags"),
    }
}
```

**Uso:**

```go
// Crear tags
golang := client.Tag.Create().
    SetName("Go").
    SetSlug("go").
    SaveX(ctx)

ent := client.Tag.Create().
    SetName("Ent").
    SetSlug("ent").
    SaveX(ctx)

database := client.Tag.Create().
    SetName("Database").
    SetSlug("database").
    SaveX(ctx)

// Crear post con tags
post := client.Post.Create().
    SetTitle("Building with Ent").
    SetContent("...").
    AddTags(golang, ent, database).
    SaveX(ctx)

// Consultar
tags, _ := post.Tags(ctx)
fmt.Printf("Post has %d tags\n", len(tags))

// Consultar posts por tag
posts, _ := golang.Posts(ctx)
fmt.Printf("%s has %d posts\n", golang.Name, len(posts))
```

### 63.3.5 Self-Referencing Relationships

Relaciones dentro de la misma entidad:

```go
// ent/schema/user.go
type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("username").Unique(),
        field.String("email").Unique(),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        // Un usuario sigue a muchos usuarios
        // Muchos usuarios son seguidos por otros
        edge.To("following", User.Type).
            From("followers"),
    }
}
```

**Uso:**

```go
// Crear usuarios
alice := client.User.Create().
    SetUsername("alice").
    SetEmail("alice@example.com").
    SaveX(ctx)

bob := client.User.Create().
    SetUsername("bob").
    SetEmail("bob@example.com").
    SaveX(ctx)

charlie := client.User.Create().
    SetUsername("charlie").
    SetEmail("charlie@example.com").
    SaveX(ctx)

// Alice sigue a Bob y Charlie
alice.Update().
    AddFollowing(bob, charlie).
    ExecX(ctx)

// Consultar a quién sigue Alice
following, _ := alice.QueryFollowing().All(ctx)
fmt.Printf("Alice follows %d users\n", len(following))

// Consultar followers de Bob
followers, _ := bob.QueryFollowers().All(ctx)
fmt.Printf("Bob has %d followers\n", len(followers))

// Consultar mutualmente
isMutual := client.User.Query().
    Where(user.ID(alice.ID)).
    QueryFollowing().
    Where(user.ID(bob.ID)).
    Exist(ctx)

println(isMutual) // true si ambos se siguen mutuamente
```

### 63.3.6 Relaciones Bidireccionales

Relaciones con referencias inversas:

```go
// ent/schema/department.go
type Department struct {
    ent.Schema
}

func (Department) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Unique(),
        field.String("description").Optional(),
    }
}

func (Department) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("employees", Employee.Type),
    }
}

// ent/schema/employee.go
type Employee struct {
    ent.Schema
}

func (Employee) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
        field.String("title"),
        field.Float("salary"),
    }
}

func (Employee) Edges() []ent.Edge {
    return []ent.Edge{
        // Referencia inversa
        edge.From("department", Department.Type).
            Ref("employees").
            Unique().
            Required(),
    }
}
```

**Uso bidireccional:**

```go
// Crear desde ambas direcciones
dept := client.Department.Create().
    SetName("Engineering").
    SaveX(ctx)

emp := client.Employee.Create().
    SetName("Alice").
    SetTitle("Senior Engineer").
    SetSalary(150000).
    SetDepartment(dept).
    SaveX(ctx)

// Acceder desde ambos lados
dept2, _ := emp.Department(ctx)
println(dept2.Name) // Engineering

employees, _ := dept.Employees(ctx)
println(len(employees)) // 1
```

---

## 63.4 Code Generation

### 63.4.1 Cómo Funciona la Generación

```
┌──────────────────────────────────────────────────┐
│      Flujo de Generación de Código en Ent        │
├──────────────────────────────────────────────────┤
│                                                  │
│  1. Schema Files (user.go, post.go)              │
│     └─► Definiciones de entidades               │
│                                                  │
│  2. Ent Code Generator                           │
│     └─► go generate ./ent                        │
│                                                  │
│  3. Generados:                                   │
│     ├─► client.go        (Cliente principal)    │
│     ├─► user.go          (Modelo)               │
│     ├─► user_create.go   (Builder de creación)  │
│     ├─► user_query.go    (Builder de queries)   │
│     ├─► user_update.go   (Builder de update)    │
│     ├─► user_delete.go   (Builder de delete)    │
│     ├─► migrate/         (Migraciones)          │
│     └─► ... (más archivos)                      │
│                                                  │
│  4. Código Type-Safe Compilable                  │
│     └─► Sin reflexión en runtime                │
│                                                  │
└──────────────────────────────────────────────────┘
```

### 63.4.2 Configuración del Generador

Crear archivo `ent/config.go`:

```go
// +build ignore

package main

import (
    "log"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
)

func main() {
    opts := []entc.Option{
        entc.Default(),
    }
    
    if err := entc.Generate("./schema", &gen.Config{}, opts...); err != nil {
        log.Fatalf("running ent codegen: %v", err)
    }
}
```

En `ent/schema`:

```go
//go:generate go run config.go
```

Ejecutar:

```bash
go generate ./ent
```

### 63.4.3 Estructura de Código Generado

**Archivo generado: `ent/user.go`**

```go
// Code generated by entc, DO NOT EDIT.

package ent

import (
    "fmt"
    "time"
)

// User is the model entity for the User schema.
type User struct {
    // ID of the ent.
    ID int `json:"id,omitempty"`
    // Email holds the value of the "email" field.
    Email string `json:"email,omitempty"`
    // Username holds the value of the "username" field.
    Username string `json:"username,omitempty"`
    // Age holds the value of the "age" field.
    Age int `json:"age,omitempty"`
    // CreatedAt holds the value of the "created_at" field.
    CreatedAt time.Time `json:"created_at,omitempty"`
}

// String implements the fmt.Stringer interface.
func (u *User) String() string {
    var builder strings.Builder
    builder.WriteString("User(")
    builder.WriteString(fmt.Sprintf("id=%v", u.ID))
    builder.WriteString(", email=")
    builder.WriteString(u.Email)
    builder.WriteString(", username=")
    builder.WriteString(u.Username)
    if u.Age != 0 {
        builder.WriteString(fmt.Sprintf(", age=%v", u.Age))
    }
    builder.WriteString(")")
    return builder.String()
}

// Update returns a builder for updating this User.
// Note that you need to call Save and return its error.
func (u *User) Update() *UserUpdate {
    return (&UserClient{config: u.config}).UpdateOne(u)
}

// Delete returns a builder for deleting this User.
func (u *User) Delete() *UserDelete {
    return (&UserClient{config: u.config}).DeleteOne(u)
}
```

**Builders: `ent/user_create.go`**

```go
// UserCreate is the builder for creating a User entity.
type UserCreate struct {
    config
    mutation *UserMutation
}

// SetEmail sets the "email" field.
func (uc *UserCreate) SetEmail(s string) *UserCreate {
    uc.mutation.SetEmail(s)
    return uc
}

// SetUsername sets the "username" field.
func (uc *UserCreate) SetUsername(s string) *UserCreate {
    uc.mutation.SetUsername(s)
    return uc
}

// SetAge sets the "age" field.
func (uc *UserCreate) SetAge(i int) *UserCreate {
    uc.mutation.SetAge(i)
    return uc
}

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
    return uc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UserCreate) SaveX(ctx context.Context) *User {
    v, err := uc.Save(ctx)
    if err != nil {
        panic(err)
    }
    return v
}
```

### 63.4.4 Builder Pattern

El patrón builder es el corazón de la API de Ent:

```go
// ✅ Builder pattern fluido
user := client.User.Create().
    SetEmail("alice@example.com").
    SetUsername("alice").
    SetAge(30).
    SaveX(ctx)

// Equivalente a:
userBuilder := client.User.Create()
userBuilder.SetEmail("alice@example.com")
userBuilder.SetUsername("alice")
userBuilder.SetAge(30)
user := userBuilder.SaveX(ctx)

// Con manejo de errores
user, err := client.User.Create().
    SetEmail("bob@example.com").
    SetUsername("bob").
    Save(ctx)
if err != nil {
    log.Printf("error creating user: %v", err)
}
```

### 63.4.5 Query Builders

```go
// Query builder con métodos encadenables
users, err := client.User.Query().
    Where(
        user.AgeGT(18),
        user.EmailContains("@example.com"),
    ).
    Order(user.ByAge()).
    Limit(10).
    All(ctx)

// Con predicados complejos
posts, err := client.Post.Query().
    Where(
        post.Or(
            post.TitleContains("Go"),
            post.ContentContains("Go"),
        ),
        post.HasAuthorWith(
            user.AgeGT(25),
        ),
    ).
    WithAuthor().
    All(ctx)

// Ejemplo avanzado
users, _ := client.User.Query().
    Where(
        user.And(
            user.AgeGT(18),
            user.AgeLT(65),
        ),
    ).
    Order(
        user.ByAge(sql.OrderDesc()),
        user.ByEmail(),
    ).
    Offset(10).
    Limit(20).
    All(ctx)
```

---

## 63.5 CRUD Operations

### 63.5.1 Create (Crear)

Insertar nuevos registros con validación automática:

```go
package main

import (
    "context"
    "log"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func main() {
    client, _ := ent.Open("postgres", dsn)
    defer client.Close()
    ctx := context.Background()
    
    // ✅ Creación simple
    user, err := client.User.Create().
        SetEmail("alice@example.com").
        SetUsername("alice").
        SetAge(30).
        Save(ctx)
    if err != nil {
        log.Printf("error creating user: %v", err)
    }
    
    // ✅ Crear con valores opcionales
    user2, _ := client.User.Create().
        SetEmail("bob@example.com").
        SetUsername("bob").
        SetAge(25).
        SetBio("Software engineer"). // campo opcional
        Save(ctx)
    
    // ✅ Crear múltiples
    users := []*ent.User{
        client.User.Create().
            SetEmail("charlie@example.com").
            SetUsername("charlie").
            SetAge(35).
            SaveX(ctx),
        client.User.Create().
            SetEmail("diana@example.com").
            SetUsername("diana").
            SetAge(28).
            SaveX(ctx),
    }
    
    // ✅ Creación con relaciones
    profile := client.Profile.Create().
        SetBio("Engineer").
        SaveX(ctx)
    
    user3 := client.User.Create().
        SetEmail("eve@example.com").
        SetUsername("eve").
        SetProfile(profile).
        SaveX(ctx)
}
```

### 63.5.2 Read (Leer)

Consultar registros con predicados type-safe:

```go
package main

import (
    "context"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func main() {
    client, _ := ent.Open("postgres", dsn)
    ctx := context.Background()
    
    // ✅ Obtener por ID
    u, _ := client.User.Get(ctx, 1)
    println(u.Email)
    
    // ✅ Query simple
    users, _ := client.User.Query().All(ctx)
    
    // ✅ Con filtros
    adults, _ := client.User.Query().
        Where(user.AgeGT(18)).
        All(ctx)
    
    // ✅ Con múltiples predicados
    results, _ := client.User.Query().
        Where(
            user.EmailContains("@example.com"),
            user.AgeGE(25),
            user.AgeLT(50),
        ).
        All(ctx)
    
    // ✅ Ordenamiento
    sorted, _ := client.User.Query().
        Order(user.ByAge()).
        All(ctx)
    
    // ✅ Reverse order
    descending, _ := client.User.Query().
        Order(user.ByAge(sql.OrderDesc())).
        All(ctx)
    
    // ✅ Obtener primero
    first, _ := client.User.Query().
        Order(user.ByEmail()).
        First(ctx)
    
    // ✅ Con Only (error si no existe)
    only, err := client.User.Query().
        Where(user.ID(1)).
        Only(ctx)
    
    // ✅ Contar
    count, _ := client.User.Query().
        Where(user.AgeGT(18)).
        Count(ctx)
    
    // ✅ Exist
    exists, _ := client.User.Query().
        Where(user.ID(999)).
        Exist(ctx)
}
```

### 63.5.3 Update (Actualizar)

Modificar registros existentes:

```go
package main

import (
    "context"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func main() {
    client, _ := ent.Open("postgres", dsn)
    ctx := context.Background()
    
    // ✅ Actualizar un registro
    u, _ := client.User.Get(ctx, 1)
    u, _ = u.Update().
        SetAge(31).
        SetBio("Updated bio").
        Save(ctx)
    
    // ✅ Actualizar múltiples (sin obtener primero)
    affected, _ := client.User.Update().
        Where(user.AgeGT(30)).
        SetBio("Senior developer").
        Exec(ctx)
    println(affected) // Número de registros actualizados
    
    // ✅ Actualizar por ID
    updated, _ := client.User.UpdateOneID(1).
        SetAge(32).
        SetEmail("newemail@example.com").
        Save(ctx)
    
    // ✅ Actualizar múltiples IDs
    affected, _ = client.User.Update().
        Where(user.IDIn(1, 2, 3)).
        SetRole("admin").
        Exec(ctx)
    
    // ✅ Incrementar/Decrementar
    client.User.Update().
        Where(user.ID(1)).
        AddAge(1). // Incrementar edad
        Exec(ctx)
    
    // ✅ Clear (limpiar campo opcional)
    client.User.Update().
        Where(user.ID(1)).
        ClearBio(). // Limpiar campo "bio"
        Exec(ctx)
    
    // ✅ Actualizar con relaciones
    dept, _ := client.Department.Get(ctx, 1)
    client.Employee.Update().
        Where(user.ID(1)).
        SetDepartment(dept).
        Exec(ctx)
}
```

### 63.5.4 Delete (Eliminar)

Eliminar registros:

```go
package main

import (
    "context"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func main() {
    client, _ := ent.Open("postgres", dsn)
    ctx := context.Background()
    
    // ✅ Eliminar un registro
    u, _ := client.User.Get(ctx, 1)
    client.User.DeleteOne(u).Exec(ctx)
    
    // ✅ Eliminar por ID
    client.User.DeleteOneID(1).Exec(ctx)
    
    // ✅ Eliminar múltiples
    affected, _ := client.User.Delete().
        Where(user.AgeGT(100)).
        Exec(ctx)
    println(affected) // Número eliminado
    
    // ✅ Eliminar todos
    affected, _ = client.User.Delete().
        Exec(ctx)
    
    // ✅ Eliminar por condiciones complejas
    affected, _ = client.User.Delete().
        Where(
            user.Or(
                user.EmailSuffix("@oldomain.com"),
                user.CreatedAtLT(time.Now().AddDate(-2, 0, 0)), // Más de 2 años
            ),
        ).
        Exec(ctx)
}
```

### 63.5.5 Batch Operations

Operaciones en lotes para mejorar performance:

```go
package main

import (
    "context"
    "my-ent-app/ent"
)

func batchCreate(client *ent.Client, ctx context.Context, count int) {
    // ✅ Crear múltiples en paralelo
    bulk := make([]*ent.UserCreate, count)
    
    for i := 0; i < count; i++ {
        bulk[i] = client.User.Create().
            SetEmail(fmt.Sprintf("user%d@example.com", i)).
            SetUsername(fmt.Sprintf("user%d", i)).
            SetAge(20 + i%50)
    }
    
    users, _ := client.User.CreateBulk(bulk...).Save(ctx)
    println(len(users), "users created")
}

func batchUpdate(client *ent.Client, ctx context.Context) {
    // ✅ Actualizar en lotes (más eficiente que queries)
    rows, _ := client.User.Update().
        Where(user.AgeGT(30)).
        SetRole("senior").
        Exec(ctx)
    println(rows, "updated")
}

func batchDelete(client *ent.Client, ctx context.Context) {
    // ✅ Eliminar en lotes
    users, _ := client.User.Query().
        Where(user.StatusEQ("inactive")).
        Limit(1000).
        All(ctx)
    
    ids := make([]int, len(users))
    for i, u := range users {
        ids[i] = u.ID
    }
    
    affected, _ := client.User.Delete().
        Where(user.IDIn(ids...)).
        Exec(ctx)
    
    println(affected, "deleted")
}
```

---

## 63.6 Advanced Querying

### 63.6.1 Where Predicates (Predicados)

Construcción type-safe de cláusulas WHERE:

```go
package main

import (
    "context"
    "my-ent-app/ent/user"
    "my-ent-app/ent/post"
)

func advancedQueries(ctx context.Context, client *ent.Client) {
    
    // ✅ Comparaciones numéricas
    users, _ := client.User.Query().Where(
        user.AgeEQ(30),      // Age == 30
        user.AgeGT(18),      // Age > 18
        user.AgeGTE(18),     // Age >= 18
        user.AgeLT(65),      // Age < 65
        user.AgeLTE(65),     // Age <= 65
        user.AgeNEQ(0),      // Age != 0
    ).All(ctx)
    
    // ✅ String operations
    results, _ := client.User.Query().Where(
        user.EmailEQ("alice@example.com"),
        user.EmailHasPrefix("alice"),
        user.EmailHasSuffix("@example.com"),
        user.EmailContains("example"),
        user.EmailIn("alice@ex.com", "bob@ex.com"),
        user.EmailNotIn("spam@ex.com"),
    ).All(ctx)
    
    // ✅ Operadores lógicos
    filtered, _ := client.User.Query().Where(
        user.And(
            user.AgeGT(18),
            user.AgeGT(65),
        ),
        user.Or(
            user.EmailContains("@gmail.com"),
            user.EmailContains("@hotmail.com"),
        ),
        user.Not(
            user.StatusEQ("banned"),
        ),
    ).All(ctx)
    
    // ✅ Predicados complejos
    complex, _ := client.User.Query().Where(
        user.Or(
            user.And(
                user.AgeGT(25),
                user.EmailContains("@company.com"),
            ),
            user.And(
                user.RoleEQ("admin"),
                user.VerifiedEQ(true),
            ),
        ),
    ).All(ctx)
    
    // ✅ NULL checks
    optional, _ := client.User.Query().Where(
        user.BioIsNil(),      // BIO es NULL
        user.BioNotNil(),     // BIO no es NULL
    ).All(ctx)
}
```

### 63.6.2 Ordering y Pagination

Ordenamiento y paginación:

```go
package main

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "my-ent-app/ent/user"
)

func pagination(ctx context.Context, client *ent.Client) {
    
    // ✅ Ordenamiento simple
    ascending, _ := client.User.Query().
        Order(user.ByEmail()).
        All(ctx)
    
    // ✅ Orden descendente
    descending, _ := client.User.Query().
        Order(user.ByAge(sql.OrderDesc())).
        All(ctx)
    
    // ✅ Ordenamiento múltiple
    sorted, _ := client.User.Query().
        Order(
            user.ByAge(sql.OrderDesc()),
            user.ByEmail(),
        ).
        All(ctx)
    
    // ✅ Paginación
    pageSize := 20
    pageNum := 1
    offset := (pageNum - 1) * pageSize
    
    page, _ := client.User.Query().
        Order(user.ByEmail()).
        Offset(offset).
        Limit(pageSize).
        All(ctx)
    
    println(len(page), "users in page")
    
    // ✅ Cursor-based pagination
    users, _ := client.User.Query().
        Order(user.ByID()).
        Limit(11).
        All(ctx)
    
    hasNext := len(users) > 10
    if hasNext {
        users = users[:10]
        nextCursor := users[len(users)-1].ID
        println("Next cursor:", nextCursor)
    }
}
```

### 63.6.3 Eager Loading (WithX)

Cargar relaciones automáticamente:

```go
package main

import (
    "context"
    "my-ent-app/ent"
    "my-ent-app/ent/post"
)

func eagerLoading(ctx context.Context, client *ent.Client) {
    
    // ❌ Sin eager loading (N+1 queries)
    posts, _ := client.Post.Query().All(ctx)
    for _, p := range posts {
        author, _ := p.QueryAuthor().Only(ctx)
        println(author.Username)
    }
    
    // ✅ Con eager loading (1 query)
    posts, _ = client.Post.Query().
        WithAuthor().
        All(ctx)
    for _, p := range posts {
        author := p.Edges.Author
        println(author.Username)
    }
    
    // ✅ Múltiples relaciones
    posts, _ = client.Post.Query().
        WithAuthor().
        WithComments().
        WithTags().
        All(ctx)
    
    for _, p := range posts {
        println("Post:", p.Title)
        println("Author:", p.Edges.Author.Username)
        println("Comments:", len(p.Edges.Comments))
        println("Tags:", len(p.Edges.Tags))
    }
    
    // ✅ Eager loading condicional
    posts, _ = client.Post.Query().
        WithAuthor(func(q *ent.UserQuery) {
            q.Where(user.VerifiedEQ(true))
        }).
        All(ctx)
    
    // ✅ Nested eager loading
    users, _ := client.User.Query().
        WithPosts(func(q *ent.PostQuery) {
            q.WithComments().WithTags()
        }).
        All(ctx)
}
```

### 63.6.4 Aggregations (Agregaciones)

Funciones de agregación:

```go
package main

import (
    "context"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func aggregations(ctx context.Context, client *ent.Client) {
    
    // ✅ Count
    count, _ := client.User.Query().Count(ctx)
    println("Total users:", count)
    
    // ✅ Count con filtro
    adults, _ := client.User.Query().
        Where(user.AgeGTE(18)).
        Count(ctx)
    
    // ✅ Min y Max
    minAge, _ := client.User.Query().Aggregate(
        ent.Min(user.FieldAge),
    ).Int(ctx)
    
    maxAge, _ := client.User.Query().Aggregate(
        ent.Max(user.FieldAge),
    ).Int(ctx)
    
    println("Age range:", minAge, "-", maxAge)
    
    // ✅ Sum y Avg
    totalAge, _ := client.User.Query().Aggregate(
        ent.Sum(user.FieldAge),
    ).Int(ctx)
    
    avgAge, _ := client.User.Query().Aggregate(
        ent.Avg(user.FieldAge),
    ).Float64(ctx)
    
    println("Average age:", avgAge)
    
    // ✅ Group by
    results, _ := client.User.Query().
        GroupBy(user.FieldRole).
        Aggregate(ent.Count()).
        All(ctx)
    
    for _, res := range results {
        role, _ := res.Scan()
        count, _ := res.Int()
        println(role, ":", count, "users")
    }
}
```

### 63.6.5 Graph Traversal

Traversar el grafo de datos:

```go
package main

import (
    "context"
    "my-ent-app/ent"
)

func graphTraversal(ctx context.Context, client *ent.Client) {
    
    // ✅ Traversal simple
    user, _ := client.User.Get(ctx, 1)
    
    // Todos los posts del usuario
    posts, _ := user.QueryPosts().All(ctx)
    
    // Todos los comentarios de los posts
    for _, post := range posts {
        comments, _ := post.QueryComments().All(ctx)
        println(post.Title, "has", len(comments), "comments")
    }
    
    // ✅ Traversal en cadena
    user, _ = client.User.Get(ctx, 1)
    authoredCommentAuthors, _ := user.
        QueryPosts().
        QueryComments().
        QueryAuthor().
        All(ctx)
    
    println(len(authoredCommentAuthors), "users commented on my posts")
    
    // ✅ Traversal bidireccional
    post, _ := client.Post.Get(ctx, 1)
    
    // Autor del post
    author, _ := post.QueryAuthor().Only(ctx)
    
    // Otros posts del autor
    otherPosts, _ := author.QueryPosts().
        Where(post.IDNEQ(post.ID)).
        All(ctx)
    
    println("Author has", len(otherPosts), "other posts")
}
```

---

## 63.7 Hooks & Interceptors

### 63.7.1 Hooks (Ganchos)

Ejecutar código antes y después de operaciones:

```go
// ent/schema/user.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("email"),
        field.String("username"),
    }
}

func (User) Hooks() []ent.Hook {
    return []ent.Hook{
        // Hook on create
        ent.On(
            ent.OpCreate|ent.OpUpdate,
            ValidateEmail(),
            NormalizeEmail(),
        ),
        // Hook on update
        ent.On(
            ent.OpUpdateOne,
            CheckPasswordUpdate(),
        ),
        // Hook on delete
        ent.On(
            ent.OpDelete,
            LogDeletion(),
        ),
    }
}

// Validador de email
func ValidateEmail() ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            if userMut, ok := m.(*ent.UserMutation); ok {
                if email, exists := userMut.Email(); exists {
                    if !isValidEmail(email) {
                        return nil, fmt.Errorf("invalid email: %s", email)
                    }
                }
            }
            return next.Mutate(ctx, m)
        })
    }
}

// Normalizar email
func NormalizeEmail() ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            if userMut, ok := m.(*ent.UserMutation); ok {
                if email, exists := userMut.Email(); exists {
                    userMut.SetEmail(strings.ToLower(email))
                }
            }
            return next.Mutate(ctx, m)
        })
    }
}

// Logging de eliminación
func LogDeletion() ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            if userMut, ok := m.(*ent.UserMutation); ok {
                if id, exists := userMut.ID(); exists {
                    log.Printf("Deleting user %d", id)
                }
            }
            return next.Mutate(ctx, m)
        })
    }
}
```

### 63.7.2 Validators (Validadores)

Validación integrada en el schema:

```go
// ent/schema/product.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "fmt"
)

type Product struct {
    ent.Schema
}

func (Product) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            NotEmpty().
            MinLen(1).
            MaxLen(255),
        
        field.String("sku").
            Unique().
            Match(regexp.MustCompile(`^[A-Z0-9-]+$`)),
        
        field.Float("price").
            Min(0).
            Positive(),
        
        field.Int("stock").
            Min(0).
            Default(0),
        
        field.String("category").
            Enums("electronics", "books", "clothing", "food"),
    }
}

// Custom validator
func (Product) Validators() []func(Product) error {
    return []func(Product) error{
        func(p Product) error {
            if p.Price > 0 && p.Stock == 0 {
                return fmt.Errorf("available products must have stock > 0")
            }
            return nil
        },
    }
}
```

### 63.7.3 Middleware Pattern

Crear middleware reutilizable:

```go
package middleware

import (
    "context"
    "fmt"
    "log"
    "time"
    "my-ent-app/ent"
)

// LoggingMiddleware registra todas las operaciones
func LoggingMiddleware() ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            start := time.Now()
            op := m.Op()
            
            value, err := next.Mutate(ctx, m)
            
            duration := time.Since(start)
            status := "success"
            if err != nil {
                status = "error"
            }
            
            log.Printf("[%s] %s %s (%v)",
                op,
                m.Type(),
                status,
                duration,
            )
            
            return value, err
        })
    }
}

// AuditMiddleware registra cambios
func AuditMiddleware(auditLog chan string) ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            value, err := next.Mutate(ctx, m)
            
            if err == nil {
                msg := fmt.Sprintf("%s: %s (%s)",
                    time.Now().Format(time.RFC3339),
                    m.Op(),
                    m.Type(),
                )
                select {
                case auditLog <- msg:
                default:
                }
            }
            
            return value, err
        })
    }
}

// PermissionMiddleware verifica permisos
func PermissionMiddleware(userRole string) ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            if userRole != "admin" && m.Op().String() == "DELETE" {
                return nil, fmt.Errorf("permission denied: only admins can delete")
            }
            
            return next.Mutate(ctx, m)
        })
    }
}
```

---

## 63.8 Transactions

### 63.8.1 Patrón TX (Transacciones)

Ejecutar múltiples operaciones de forma atómica:

```go
package main

import (
    "context"
    "fmt"
    "my-ent-app/ent"
)

func transferFunds(ctx context.Context, client *ent.Client) error {
    
    // Crear transacción
    tx, err := client.Tx(ctx)
    if err != nil {
        return fmt.Errorf("starting transaction: %w", err)
    }
    
    // Operación 1: Restar de cuenta origen
    sender, err := tx.Account.Query().
        Where(account.ID(1)).
        ForUpdate(). // Lock para evitar race conditions
        Only(ctx)
    if err != nil {
        return rollback(tx, err)
    }
    
    if sender.Balance < 100 {
        return rollback(tx, fmt.Errorf("insufficient funds"))
    }
    
    if err := sender.Update().
        AddBalance(-100).
        Exec(ctx); err != nil {
        return rollback(tx, err)
    }
    
    // Operación 2: Sumar a cuenta destino
    receiver, err := tx.Account.Query().
        Where(account.ID(2)).
        ForUpdate().
        Only(ctx)
    if err != nil {
        return rollback(tx, err)
    }
    
    if err := receiver.Update().
        AddBalance(100).
        Exec(ctx); err != nil {
        return rollback(tx, err)
    }
    
    // Registrar transacción
    _, err = tx.Transaction.Create().
        SetFromAccount(sender).
        SetToAccount(receiver).
        SetAmount(100).
        Save(ctx)
    if err != nil {
        return rollback(tx, err)
    }
    
    // Commit
    return tx.Commit()
}

func rollback(tx *ent.Tx, err error) error {
    if rbErr := tx.Rollback(); rbErr != nil {
        return fmt.Errorf("rolling back transaction: %v (original error: %w)", rbErr, err)
    }
    return err
}
```

### 63.8.2 Transacciones Anidadas

```go
func complexTransaction(ctx context.Context, client *ent.Client) error {
    
    tx, err := client.Tx(ctx)
    if err != nil {
        return err
    }
    
    // Crear usuario
    user, err := tx.User.Create().
        SetEmail("alice@example.com").
        SetUsername("alice").
        Save(ctx)
    if err != nil {
        return rollback(tx, err)
    }
    
    // Crear múltiples posts
    posts := make([]*ent.Post, 3)
    for i := 0; i < 3; i++ {
        post, err := tx.Post.Create().
            SetTitle(fmt.Sprintf("Post %d", i+1)).
            SetContent("...").
            SetAuthor(user).
            Save(ctx)
        if err != nil {
            return rollback(tx, err)
        }
        posts[i] = post
    }
    
    // Agregar tags a posts
    for _, post := range posts {
        _, err := tx.Tag.Create().
            SetName("golang").
            AddPosts(post).
            Save(ctx)
        if err != nil {
            return rollback(tx, err)
        }
    }
    
    return tx.Commit()
}
```

### 63.8.3 Manejo de Errores en Transacciones

```go
func safeTransaction(ctx context.Context, client *ent.Client) (*ent.User, error) {
    
    tx, err := client.Tx(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %w", err)
    }
    
    user, err := tx.User.Create().
        SetEmail("bob@example.com").
        SetUsername("bob").
        Save(ctx)
    
    if err != nil {
        // Rollback automático en caso de error
        if rbErr := tx.Rollback(); rbErr != nil {
            return nil, fmt.Errorf("failed to rollback: %w (original: %w)", rbErr, err)
        }
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit: %w", err)
    }
    
    return user, nil
}
```

---

## 63.9 Advanced Features

### 63.9.1 Field Masks (Máscaras de Campos)

Actualizar solo campos específicos:

```go
package main

import (
    "context"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func updateWithMask(ctx context.Context, client *ent.Client) {
    
    // Actualizar solo email
    client.User.UpdateOneID(1).
        SetEmail("newemail@example.com").
        ClearAge(). // Limpiar opcional
        Exec(ctx)
    
    // Actualizar parcial con validación
    updates := map[string]interface{}{
        "email":     "alice@example.com",
        "username":  "alice_new",
    }
    
    builder := client.User.UpdateOneID(1)
    
    if email, ok := updates["email"]; ok {
        builder.SetEmail(email.(string))
    }
    
    if username, ok := updates["username"]; ok {
        builder.SetUsername(username.(string))
    }
    
    builder.Exec(ctx)
}
```

### 63.9.2 Privacy Rules (Reglas de Privacidad)

Control de acceso a nivel de schema:

```go
// ent/schema/post.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "context"
)

type Post struct {
    ent.Schema
}

func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.String("title"),
        field.String("content"),
    }
}

// Privacy define reglas de acceso
func (Post) Policy() ent.Policy {
    return privacy.Policy{
        Mutation: privacy.MutationPolicy{
            // Solo el autor puede actualizar
            Rule: privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
                if m.Op().Is(ent.OpUpdateOne|ent.OpDelete) {
                    postMut := m.(*ent.PostMutation)
                    userID := auth.UserIDFromContext(ctx)
                    
                    authorID, _ := postMut.AuthorID()
                    if authorID != userID {
                        return privacy.Deny
                    }
                }
                return nil
            }),
        },
        Query: privacy.QueryPolicy{
            // Solo usuarios verificados pueden ver borradores
            Rule: privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
                postQ := q.(*ent.PostQuery)
                if isDraft := isDraftPost(postQ); isDraft {
                    userID := auth.UserIDFromContext(ctx)
                    if userID == 0 {
                        return privacy.Deny
                    }
                }
                return nil
            }),
        },
    }
}
```

### 63.9.3 Recursive Queries (Queries Recursivas)

Consultas jerárquicas:

```go
package main

import (
    "context"
    "my-ent-app/ent"
)

// Obtener comentarios anidados
func getCommentsTree(ctx context.Context, client *ent.Client, postID int) {
    
    post, _ := client.Post.Get(ctx, postID)
    
    // Obtener comentarios top-level
    topComments, _ := post.QueryComments().
        Where(comment.ParentIsNil()). // Sin padre
        All(ctx)
    
    // Para cada comentario, obtener respuestas
    for _, comment := range topComments {
        replies, _ := comment.QueryReplies().All(ctx)
        
        println(fmt.Sprintf("- %s (%d replies)", comment.Text, len(replies)))
        
        for _, reply := range replies {
            nestedReplies, _ := reply.QueryReplies().All(ctx)
            println(fmt.Sprintf("  - %s (%d replies)", reply.Text, len(nestedReplies)))
        }
    }
}

// Obtener árbol de categorías
func getCategoryTree(ctx context.Context, client *ent.Client, rootID int) {
    
    category, _ := client.Category.Get(ctx, rootID)
    printCategoryTree(ctx, client, category, 0)
}

func printCategoryTree(ctx context.Context, client *ent.Client, cat *ent.Category, depth int) {
    indent := strings.Repeat("  ", depth)
    println(indent + cat.Name)
    
    children, _ := cat.QueryChildren().All(ctx)
    for _, child := range children {
        printCategoryTree(ctx, client, child, depth+1)
    }
}
```

---

## 63.10 Testing & Integration

### 63.10.1 Configuración de Testing

```go
// tests/setup.go
package tests

import (
    "context"
    "log"
    "my-ent-app/ent"
    _ "github.com/mattn/go-sqlite3"
)

// NewTestClient crea un cliente de testing con SQLite en memoria
func NewTestClient(t *testing.T) *ent.Client {
    client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
    if err != nil {
        t.Fatalf("failed opening connection to sqlite: %v", err)
    }
    
    if err := client.Schema.Create(context.Background()); err != nil {
        t.Fatalf("failed creating schema resources: %v", err)
    }
    
    t.Cleanup(func() {
        client.Close()
    })
    
    return client
}

// SeedData proporciona datos de prueba
func SeedData(t *testing.T, client *ent.Client, ctx context.Context) {
    users := make([]*ent.User, 3)
    for i := 0; i < 3; i++ {
        user, err := client.User.Create().
            SetEmail(fmt.Sprintf("user%d@example.com", i)).
            SetUsername(fmt.Sprintf("user%d", i)).
            SetAge(20 + i*5).
            Save(ctx)
        if err != nil {
            t.Fatalf("failed seeding users: %v", err)
        }
        users[i] = user
    }
}
```

### 63.10.2 Unit Tests

```go
// user_test.go
package tests

import (
    "context"
    "testing"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func TestCreateUser(t *testing.T) {
    client := NewTestClient(t)
    ctx := context.Background()
    
    u, err := client.User.Create().
        SetEmail("test@example.com").
        SetUsername("testuser").
        SetAge(30).
        Save(ctx)
    
    if err != nil {
        t.Fatalf("failed creating user: %v", err)
    }
    
    if u.Email != "test@example.com" {
        t.Errorf("expected email test@example.com, got %s", u.Email)
    }
}

func TestUniqueConstraint(t *testing.T) {
    client := NewTestClient(t)
    ctx := context.Background()
    
    // Crear primer usuario
    _, err := client.User.Create().
        SetEmail("alice@example.com").
        SetUsername("alice").
        Save(ctx)
    if err != nil {
        t.Fatalf("failed creating first user: %v", err)
    }
    
    // Intentar crear otro con mismo email
    _, err = client.User.Create().
        SetEmail("alice@example.com").
        SetUsername("alice2").
        Save(ctx)
    
    if err == nil {
        t.Error("expected error for duplicate email, got nil")
    }
}

func TestQueryFilter(t *testing.T) {
    client := NewTestClient(t)
    ctx := context.Background()
    
    // Crear usuarios
    for i := 0; i < 5; i++ {
        client.User.Create().
            SetEmail(fmt.Sprintf("user%d@example.com", i)).
            SetUsername(fmt.Sprintf("user%d", i)).
            SetAge(20 + i).
            SaveX(ctx)
    }
    
    // Consultar adultos
    adults, _ := client.User.Query().
        Where(user.AgeGTE(21)).
        All(ctx)
    
    if len(adults) != 4 {
        t.Errorf("expected 4 adults, got %d", len(adults))
    }
}
```

### 63.10.3 Integration Tests

```go
// integration_test.go
package tests

import (
    "context"
    "testing"
)

func TestUserPostWorkflow(t *testing.T) {
    client := NewTestClient(t)
    ctx := context.Background()
    
    // Crear usuario
    user, _ := client.User.Create().
        SetEmail("author@example.com").
        SetUsername("author").
        SaveX(ctx)
    
    // Crear posts
    for i := 0; i < 3; i++ {
        client.Post.Create().
            SetTitle(fmt.Sprintf("Post %d", i)).
            SetContent("content").
            SetAuthor(user).
            SaveX(ctx)
    }
    
    // Verificar
    posts, _ := user.QueryPosts().All(ctx)
    if len(posts) != 3 {
        t.Errorf("expected 3 posts, got %d", len(posts))
    }
    
    // Actualizar post
    posts[0].Update().
        SetTitle("Updated Post").
        ExecX(ctx)
    
    // Verificar actualización
    updated, _ := client.Post.Get(ctx, posts[0].ID)
    if updated.Title != "Updated Post" {
        t.Errorf("expected 'Updated Post', got %s", updated.Title)
    }
}

func TestTransactionRollback(t *testing.T) {
    client := NewTestClient(t)
    ctx := context.Background()
    
    // Crear usuario inicial
    user, _ := client.User.Create().
        SetEmail("test@example.com").
        SetUsername("test").
        SaveX(ctx)
    
    // Intentar transacción con error
    tx, _ := client.Tx(ctx)
    
    _, err := tx.Post.Create().
        SetTitle("Post 1").
        SetContent("content").
        SetAuthor(user).
        Save(ctx)
    
    if err != nil {
        tx.Rollback()
    }
    
    // Verificar que no se creó
    posts, _ := client.Post.Query().All(ctx)
    if len(posts) != 0 {
        t.Errorf("expected 0 posts after rollback, got %d", len(posts))
    }
}
```

---

## 63.11 Production & Case Studies

### 63.11.1 Performance Optimization

#### Estrategias de Optimización:

```go
package main

import (
    "context"
    "my-ent-app/ent"
)

// 1. Índices apropiados
// En schema:
// index.Fields("email"),
// index.Fields("created_at"),
// index.Fields("status", "created_at"),

// 2. Eager loading para evitar N+1
// ❌ MALO: N+1 queries
users, _ := client.User.Query().All(ctx)
for _, u := range users {
    posts, _ := u.QueryPosts().All(ctx) // Extra query por usuario
}

// ✅ BUENO: 1 query
users, _ = client.User.Query().
    WithPosts().
    All(ctx)

// 3. Batch operations
// ❌ MALO: Múltiples queries
for i := 0; i < 1000; i++ {
    client.User.Create().
        SetEmail(fmt.Sprintf("user%d@example.com", i)).
        SaveX(ctx)
}

// ✅ BUENO: Batch insert
bulk := make([]*ent.UserCreate, 1000)
for i := 0; i < 1000; i++ {
    bulk[i] = client.User.Create().
        SetEmail(fmt.Sprintf("user%d@example.com", i))
}
client.User.CreateBulk(bulk...).SaveX(ctx)

// 4. Limit y Offset para paginación
users, _ := client.User.Query().
    Limit(20).
    Offset((pageNum-1)*20).
    All(ctx)

// 5. Select fields específicos (cuando sea posible)
// No siempre disponible, pero útil para queries complejas
users, _ := client.User.Query().
    Select(user.FieldID, user.FieldEmail).
    All(ctx)

// 6. Context con timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

users, err := client.User.Query().All(ctx)
if err == context.DeadlineExceeded {
    log.Println("Query timeout")
}
```

### 63.11.2 Real-World Example: Social Network

```go
// ent/schema/user.go
package schema

import (
    "time"
    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("username").Unique().MinLen(3),
        field.String("email").Unique(),
        field.String("bio").MaxLen(500).Optional(),
        field.Bytes("avatar").Optional(),
        field.Time("joined_at").Default(time.Now),
        field.Time("last_active").Default(time.Now).UpdateDefault(time.Now),
        field.Bool("verified").Default(false),
        field.Int("follower_count").Default(0),
    }
}

func (User) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("username"),
        index.Fields("email"),
        index.Fields("created_at"),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("posts", Post.Type),
        edge.To("likes", Like.Type),
        edge.To("following", User.Type).From("followers"),
        edge.To("messages", Message.Type),
    }
}

// ent/schema/post.go
type Post struct {
    ent.Schema
}

func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.String("content").NotEmpty(),
        field.Int("like_count").Default(0),
        field.Int("comment_count").Default(0),
        field.Time("created_at").Default(time.Now).Immutable(),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (Post) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("author", User.Type).Ref("posts").Required().Unique(),
        edge.To("comments", Comment.Type),
        edge.To("likes", Like.Type),
        edge.To("tags", Tag.Type),
    }
}

// Queries útiles
func (client *Client) GetUserFeed(ctx context.Context, userID int) ([]*Post, error) {
    return client.Post.Query().
        Where(post.HasAuthorWith(
            user.Or(
                user.ID(userID),
                user.HasFollowersWith(user.ID(userID)),
            ),
        )).
        WithAuthor().
        WithComments().
        WithLikes().
        Order(post.ByCreatedAtDesc()).
        Limit(20).
        All(ctx)
}

func (client *Client) GetTrendingPosts(ctx context.Context) ([]*Post, error) {
    week := time.Now().AddDate(0, 0, -7)
    return client.Post.Query().
        Where(post.CreatedAtGT(week)).
        Order(post.ByLikeCount()).
        Limit(10).
        All(ctx)
}
```

### 63.11.3 Migration Strategies

```go
// migrations/initial.sql
-- Crear versión inicial del schema
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

// migrations/add_posts.sql
-- Agregar tabla de posts
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
```

### 63.11.4 Best Practices

```
✅ BEST PRACTICES

1. Schema Design
   ├─ Define schemas claramente
   ├─ Use índices estratégicamente
   └─ Mantenga relaciones simples y claras

2. Query Optimization
   ├─ Use eager loading (WithX)
   ├─ Evite N+1 queries
   ├─ Implemente paginación
   └─ Use batch operations para múltiples

3. Validation
   ├─ Valide en schema level
   ├─ Use hooks para lógica compleja
   └─ Verifique restricciones únicas

4. Transactions
   ├─ Use para operaciones multi-entidad
   ├─ Siempre haga rollback en error
   └─ Mantenga transacciones cortas

5. Testing
   ├─ Use SQLite en memoria
   ├─ Test schema constraints
   ├─ Seed datos apropiadamente
   └─ Verifique edge cases

6. Production
   ├─ Use connection pooling
   ├─ Implemente logging
   ├─ Monitoree performance
   └─ Planee estrategia de backup
```

### 63.11.5 Troubleshooting

```go
// Problema: N+1 Queries
// ❌ MALO
for _, post := range posts {
    author, _ := post.QueryAuthor().Only(ctx)
}

// ✅ BUENO
posts, _ := client.Post.Query().
    WithAuthor().
    All(ctx)
for _, post := range posts {
    author := post.Edges.Author
}

// Problema: Validación fallida
// Solución: Revisar logs
err := client.User.Create().
    SetEmail("invalid-email").
    Exec(ctx)
if err != nil {
    log.Printf("validation error: %v", err)
}

// Problema: Race condition en transacción
// Solución: Use FOR UPDATE
user, _ := tx.User.Query().
    Where(user.ID(1)).
    ForUpdate(). // Lock exclusivo
    Only(ctx)

// Problema: Memory leaks
// Solución: Siempre cierre cliente
defer client.Close()

// Problema: Schema mismatch
// Solución: Regenerar código
go generate ./ent
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: Tu Primer Entity (Beginners)

**Objetivo:** Crear y usar una entidad simple

```go
// TAREA:
// 1. Crear schema de User con fields: email, username, age
// 2. Generar código
// 3. Crear 3 usuarios
// 4. Consultar usuarios > 18 años
// 5. Actualizar un usuario
// 6. Eliminar un usuario

// SOLUCIÓN:
package main

import (
    "context"
    "fmt"
    "my-ent-app/ent"
    "my-ent-app/ent/user"
)

func ejercicio1(client *ent.Client) {
    ctx := context.Background()
    
    // Crear schema
    // (Ver ent/schema/user.go)
    
    // Crear usuarios
    users := make([]*ent.User, 3)
    for i := 0; i < 3; i++ {
        u, _ := client.User.Create().
            SetEmail(fmt.Sprintf("user%d@example.com", i)).
            SetUsername(fmt.Sprintf("user%d", i)).
            SetAge(15 + i*5).
            Save(ctx)
        users[i] = u
    }
    
    // Consultar adultos
    adults, _ := client.User.Query().
        Where(user.AgeGTE(18)).
        All(ctx)
    fmt.Printf("Adults: %d\n", len(adults))
    
    // Actualizar
    users[0].Update().SetAge(21).ExecX(ctx)
    
    // Eliminar
    client.User.DeleteOne(users[2]).ExecX(ctx)
}
```

### Ejercicio 2: Relaciones (Intermediate)

**Objetivo:** Crear entidades con relaciones

```go
// TAREA:
// 1. Crear schema de Post con relación a User (author)
// 2. Crear usuario
// 3. Crear 5 posts del usuario
// 4. Consultar posts con eager loading
// 5. Agregar comentarios a posts

package main

func ejercicio2(client *ent.Client) {
    ctx := context.Background()
    
    // Crear usuario
    author, _ := client.User.Create().
        SetEmail("author@example.com").
        SetUsername("author").
        SaveX(ctx)
    
    // Crear posts
    for i := 1; i <= 5; i++ {
        client.Post.Create().
            SetTitle(fmt.Sprintf("Post %d", i)).
            SetContent(fmt.Sprintf("Content of post %d", i)).
            SetAuthor(author).
            SaveX(ctx)
    }
    
    // Consultar con eager loading
    posts, _ := client.Post.Query().
        WithAuthor().
        All(ctx)
    
    for _, p := range posts {
        fmt.Printf("%s by %s\n", p.Title, p.Edges.Author.Username)
    }
}
```

### Ejercicio 3: Queries Complejas (Intermediate)

**Objetivo:** Usar predicados y agregaciones avanzadas

```go
// TAREA:
// 1. Crear múltiples usuarios con edades variadas
// 2. Crear posts de diferentes autores
// 3. Consultar: posts de usuarios > 25 años
// 4. Contar posts por usuario
// 5. Ordenar posts por fecha descendente

package main

func ejercicio3(client *ent.Client) {
    ctx := context.Background()
    
    // Usuarios
    users := make([]*ent.User, 4)
    ages := []int{22, 28, 35, 19}
    for i := 0; i < 4; i++ {
        u, _ := client.User.Create().
            SetEmail(fmt.Sprintf("user%d@example.com", i)).
            SetUsername(fmt.Sprintf("user%d", i)).
            SetAge(ages[i]).
            SaveX(ctx)
        users[i] = u
    }
    
    // Posts
    for i, u := range users {
        for j := 0; j < 3; j++ {
            client.Post.Create().
                SetTitle(fmt.Sprintf("Post %d-%d", i, j)).
                SetContent("...").
                SetAuthor(u).
                SaveX(ctx)
        }
    }
    
    // Query: Posts de usuarios > 25
    posts, _ := client.Post.Query().
        Where(post.HasAuthorWith(user.AgeGT(25))).
        WithAuthor().
        Order(post.ByCreatedAtDesc()).
        All(ctx)
    
    fmt.Printf("Posts from authors > 25: %d\n", len(posts))
}
```

### Ejercicio 4: Hooks y Validación (Advanced)

**Objetivo:** Implementar validación y hooks personalizados

```go
// TAREA:
// 1. Implementar validador de email
// 2. Crear hook que normalice email a minúsculas
// 3. Crear hook que registre creaciones
// 4. Intentar crear usuario con email inválido
// 5. Verificar que se aplicaron hooks

package main

// Ver sección 63.7 para implementación completa
// Resumen:
// - Define validadores en schema
// - Define hooks en schema
// - Crea User con email
// - Verifica normalización
```

### Ejercicio 5: Red Social Completa (Expert)

**Objetivo:** Construcción completa de sistema de red social

```go
// TAREA:
// 1. Crear schema de User, Post, Comment, Tag
// 2. Relaciones: User-Post, User-Comment, Post-Comment, Post-Tag
// 3. Crear 3 usuarios
// 4. Cada usuario crea 2 posts
// 5. Agregar comentarios
// 6. Agregar tags a posts
// 7. Consultar feed de un usuario
// 8. Obtener posts trending

package main

func ejercicio5(client *ent.Client) {
    ctx := context.Background()
    
    // Crear usuarios
    users := make([]*ent.User, 3)
    for i := 0; i < 3; i++ {
        u, _ := client.User.Create().
            SetUsername(fmt.Sprintf("user%d", i)).
            SetEmail(fmt.Sprintf("user%d@example.com", i)).
            SaveX(ctx)
        users[i] = u
    }
    
    // Crear posts
    posts := make([]*ent.Post, 6)
    for i, u := range users {
        for j := 0; j < 2; j++ {
            p, _ := client.Post.Create().
                SetTitle(fmt.Sprintf("Post by %s #%d", u.Username, j+1)).
                SetContent("Post content...").
                SetAuthor(u).
                SaveX(ctx)
            posts[i*2+j] = p
        }
    }
    
    // Agregar tags
    golang, _ := client.Tag.Create().
        SetName("golang").
        SaveX(ctx)
    
    database, _ := client.Tag.Create().
        SetName("database").
        SaveX(ctx)
    
    posts[0].Update().AddTags(golang, database).ExecX(ctx)
    
    // Crear comentarios
    for i, p := range posts[:2] {
        for j := 0; j < 2; j++ {
            commenter := users[(i+1)%3]
            client.Comment.Create().
                SetText(fmt.Sprintf("Great post! %d", j)).
                SetAuthor(commenter).
                SetPost(p).
                SaveX(ctx)
        }
    }
    
    // Feed del usuario 0
    feed, _ := client.Post.Query().
        Where(post.Or(
            post.HasAuthorWith(user.ID(users[0].ID)),
            post.HasAuthorWith(
                user.IDIn(users[1].ID, users[2].ID),
            ),
        )).
        WithAuthor().
        WithComments().
        WithTags().
        Order(post.ByCreatedAtDesc()).
        All(ctx)
    
    fmt.Printf("User feed: %d posts\n", len(feed))
}
```

---

## RESUMEN COMPARATIVO

### Ent vs GORM vs sqlc

```
╔════════════════════╦═════════════════╦═════════════════╦════════════════╗
║ Característica     ║ Ent             ║ GORM            ║ sqlc           ║
╠════════════════════╬═════════════════╬═════════════════╬════════════════╣
║ Type Safety        ║ ✅ Compilación ║ ⚠️ Runtime       ║ ✅ Compilación║
║ Learning Curve     ║ ⭐⭐⭐ (3/5)   ║ ⭐⭐ (2/5)      ║ ⭐⭐⭐ (3/5) ║
║ Query Building     ║ ✅ Builders     ║ ✅ Chains       ║ SQL manual     ║
║ Performance        ║ ⭐⭐⭐ Óptimo  ║ ⭐⭐ Bueno      ║ ⭐⭐⭐ Óptimo ║
║ Validation         ║ ✅ Integrada   ║ ⚠️ Tag-based   ║ ❌ Manual      ║
║ Schema Migration   ║ ✅ Automática  ║ ✅ Automática   ║ Manual         ║
║ Hooks/Middleware   ║ ✅ Completos   ║ ⚠️ Limitados    ║ ❌ No          ║
║ Graph Operations   ║ ✅ Nativas     ║ ⚠️ Posible      ║ ❌ SQL raw     ║
║ Reflection         ║ ❌ Ninguna     ║ ⭐⭐⭐ Mucha    ║ ❌ Ninguna    ║
║ Flexibility        ║ ⭐⭐⭐ Alta    ║ ⭐⭐⭐ Máxima   ║ ⭐⭐ Limitada ║
║ Community          ║ ⭐⭐⭐ Creciente║ ⭐⭐⭐ Grande   ║ ⭐⭐ Pequeña  ║
║ Production Ready   ║ ✅ Sí          ║ ✅ Sí           ║ ✅ Sí          ║
╚════════════════════╩═════════════════╩═════════════════╩════════════════╝
```

### Cuándo usar cada uno:

**ENT es mejor para:**
- Modelos complejos con múltiples relaciones
- Type safety en compilación
- Validación integrada
- Aplicaciones medianas a grandes

**GORM es mejor para:**
- Prototipado rápido
- Máxima flexibilidad
- Migraciones complejas
- Equipos que prefieren reflexión

**sqlc es mejor para:**
- Máximo control sobre SQL
- Performance crítica
- Queries muy específicas
- Equipos que dominan SQL

---

## CONCLUSIÓN

Ent representa una evolución moderna en ORMs para Go, combinando:

1. **Type Safety**: Detección de errores en compilación
2. **Graph Thinking**: Relaciones como nodos y aristas
3. **Code Generation**: Elimina reflexión innecesaria
4. **Developer Experience**: API fluida e intuitiva
5. **Production Ready**: Usado en proyectos grandes

**Recomendación**: Ent es ideal para aplicaciones Go enterprise que requieren robustez, mantenibilidad y escalabilidad.

---

---

## REFERENCIA RÁPIDA

### Commands Básicos

```bash
# Instalar Ent
go get -u entgo.io/ent/cmd/ent

# Inicializar esquema
go run entgo.io/ent/cmd/ent init User Post

# Generar código
go generate ./ent

# Ejecutar migrations
go run main.go
```

### Snippets Más Usados

```go
// Crear
user := client.User.Create().
    SetEmail("alice@example.com").
    SaveX(ctx)

// Leer
user, _ := client.User.Get(ctx, 1)
users, _ := client.User.Query().All(ctx)

// Actualizar
user.Update().SetAge(31).ExecX(ctx)

// Eliminar
client.User.DeleteOne(user).ExecX(ctx)

// Query con filtro
users, _ := client.User.Query().
    Where(user.AgeGT(18)).
    All(ctx)

// Eager loading
posts, _ := client.Post.Query().
    WithAuthor().
    WithComments().
    All(ctx)

// Transacción
tx, _ := client.Tx(ctx)
defer tx.Rollback() // Safe rollback
// operaciones...
tx.Commit()
```

### Operadores de Predicado

| Operador | SQL | Ejemplo |
|----------|-----|---------|
| EQ | = | `user.EmailEQ("alice@ex.com")` |
| NEQ | != | `user.AgeNEQ(18)` |
| GT | > | `user.AgeGT(18)` |
| GTE | >= | `user.AgeGTE(18)` |
| LT | < | `user.AgeLT(65)` |
| LTE | <= | `user.AgeLTE(65)` |
| In | IN | `user.StatusIn("active", "pending")` |
| NotIn | NOT IN | `user.StatusNotIn("banned")` |
| HasPrefix | LIKE 'x%' | `user.EmailHasPrefix("alice")` |
| HasSuffix | LIKE '%x' | `user.EmailHasSuffix("@ex.com")` |
| Contains | LIKE '%x%' | `user.EmailContains("example")` |
| IsNil | IS NULL | `user.BioIsNil()` |
| NotNil | IS NOT NULL | `user.BioNotNil()` |

### Field Types

```go
field.Bool()           // boolean
field.Int()            // int
field.Float()          // float64
field.String()         // string
field.Bytes()          // []byte
field.UUID()           // uuid.UUID
field.Time()           // time.Time
field.JSON()           // json.RawMessage
field.Enum()           // enumeración
field.Other()          // tipos personalizados
```

### Validadores Comunes

```go
field.String("name").
    NotEmpty().              // no vacío
    MinLen(1).               // mínimo 1 carácter
    MaxLen(255).             // máximo 255 caracteres
    Unique().                // valor único
    Immutable().             // no modificable
    Match(regexp.MustCompile("^[a-z]+$")).
    Default("").             // valor por defecto

field.Int("age").
    Min(0).                  // mínimo 0
    Max(150).                // máximo 150
    Positive().              // valor positivo
    Nonnegative().           // no negativo

field.Time("created_at").
    Default(time.Now).       // ahora
    UpdateDefault(time.Now). // resetear en update
    Immutable()              // no modificable
```

### Operadores Lógicos

```go
// AND
user.And(
    user.AgeGT(18),
    user.EmailContains("@example.com"),
)

// OR
user.Or(
    user.StatusEQ("admin"),
    user.StatusEQ("moderator"),
)

// NOT
user.Not(
    user.StatusEQ("banned"),
)

// Combinaciones
user.And(
    user.AgeGT(18),
    user.Or(
        user.EmailContains("@gmail.com"),
        user.EmailContains("@hotmail.com"),
    ),
)
```

### Relaciones Quick Setup

```go
// One-to-One
edge.To("profile", Profile.Type).Unique().Required()

// One-to-Many
edge.To("posts", Post.Type)
edge.From("author", User.Type).Ref("posts").Unique().Required()

// Many-to-Many
edge.To("tags", Tag.Type)
edge.From("posts", Post.Type).Ref("tags")

// Self-Referencing
edge.To("following", User.Type).From("followers")
```

### Testing Pattern

```go
func TestFeature(t *testing.T) {
    client := ent.Open("sqlite3", "file:ent?mode=memory")
    defer client.Close()
    
    ctx := context.Background()
    client.Schema.CreateX(ctx)
    
    // test code
}
```

### Error Handling

```go
// Solo valor
u, err := client.User.Create().
    SetEmail("test@ex.com").
    Save(ctx)
if err != nil {
    log.Fatal(err)
}

// Panic version (cuidado)
u := client.User.Create().
    SetEmail("test@ex.com").
    SaveX(ctx)

// Rollback en transacción
if err != nil {
    if rbErr := tx.Rollback(); rbErr != nil {
        log.Fatalf("rollback error: %v", rbErr)
    }
    return fmt.Errorf("operation failed: %w", err)
}
```

### Performance Checklist

- [ ] Índices en campos consultados frecuentemente
- [ ] Eager loading (WithX) en queries
- [ ] Batch operations para inserción masiva
- [ ] Paginación en queries grandes
- [ ] Connection pooling configurado
- [ ] Queries con timeout
- [ ] Prepared statements (automático en Ent)

### Recursos Útiles

- **Documentación oficial**: https://entgo.io
- **Ejemplos**: https://github.com/ent/ent/tree/master/examples
- **Discord**: https://discord.gg/qZmPgTE6RQ
- **Issues**: https://github.com/ent/ent/issues
- **Changelog**: https://github.com/ent/ent/releases

---

**Fin del Capítulo 63**

*Última actualización: 2024*
*Versión de Ent referenciada: v0.12+*

---

## Apéndice: Glosario de Términos

**Entity**: Unidad fundamental de datos, representa una tabla en la BD.

**Edge**: Relación entre entidades, similar a una clave foránea.

**Schema**: Definición de estructura de datos en Go.

**Mutation**: Operación que modifica la base de datos (Create, Update, Delete).

**Builder Pattern**: API fluida que encadena métodos.

**Code Generation**: Proceso de crear código automáticamente desde schemas.

**Type-Safe**: Sistema de tipos que detecta errores en compilación.

**Graph-Based**: Modelo donde datos son nodos conectados por aristas.

**Eager Loading**: Cargar relaciones automáticamente en una sola query.

**N+1 Problem**: Anti-patrón donde múltiples queries se ejecutan innecesariamente.

**Field Mask**: Actualizar solo campos específicos.

**Hook**: Función ejecutada antes/después de operación.

**Transaction**: Conjunto de operaciones ejecutadas atómicamente.

**Predicate**: Condición de filtrado en queries (Where clause).

**Aggregation**: Operación que resume datos (Count, Sum, Avg, Min, Max).

---

*Documento exhaustivo compilado en 2024*
*Cobertura: 11 secciones, 50+ subsecciones, 56+ ejemplos de código*
*Para Go developers nivel intermediate+*
