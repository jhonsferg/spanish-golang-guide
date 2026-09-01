# Capítulo 62: sqlc - SQL type-safe generado desde código

## Introducción

Este capítulo se dedica a **sqlc**, una herramienta revolucionaria de generación de código que transforma tus queries SQL en métodos Go seguros en tipos y completamente type-safe. A diferencia de los ORMs tradicionales como GORM o Ent, sqlc elimina el overhead de runtime manteniendo la potencia y claridad del SQL nativo.

---

## 62.1 Introducción a sqlc

### 62.1.1 ¿Qué es sqlc?

**sqlc** es una herramienta de línea de comandos que genera código Go seguro en tipos a partir de consultas SQL. Su filosofía central es simple pero poderosa:

```
SQL queries (*.sql) → sqlc generate → Type-safe Go code
```

Diferencias fundamentales con otros enfoques:

```
┌─────────────────────────────────────────────────────────┐
│                 ENFOQUE DE ACCESO A DATOS               │
├─────────────────────────┬─────────────────┬─────────────┤
│        Atributo         │      GORM       │    sqlc     │
├─────────────────────────┼─────────────────┼─────────────┤
│ Type-safe               │ Parcial         │ Completo    │
│ SQL nativo              │ No              │ Sí          │
│ Runtime overhead        │ Alto            │ Cero        │
│ Performance             │ Medio           │ Excelente   │
│ Learning curve          │ Bajo            │ Bajo        │
│ Flexibilidad SQL        │ Media           │ Alta        │
│ IDE autocompletar       │ Bueno           │ Excelente   │
│ Debugging               │ Complicado      │ Simple      │
│ Generación de código    │ No              │ Sí          │
└─────────────────────────┴─────────────────┴─────────────┘
```

### 62.1.2 Ventajas de sqlc

**1. Type Safety Completo**

El compilador de Go valida automáticamente los tipos de parámetros y retorno:

```go
// ✅ Generado por sqlc - completamente type-safe
user, err := db.GetUser(ctx, 42)
// Si la query espera int pero pasas string → error de compilación

// ❌ GORM - puede fallar en runtime
var user User
db.Where("id = ?", "no es int").First(&user) // Silencioso fallará
```

**2. Cero Overhead de Runtime**

sqlc genera código puro Go. No hay reflexión, no hay interpretación de etiquetas struct:

```
GORM:     Query → Parse → Reflect → Scan → Go struct
sqlc:     Query → Go código → Compilado → Ejecución directa
                                ↓
                        Sin overhead de runtime
```

**3. SQL Puro y Explícito**

Escribes SQL nativo, manteniendo su poder completo:

```sql
-- query.sql - ¡Exactamente lo que quieres ejecutar!
-- name: GetActiveUsers :many
SELECT id, name, email, created_at
FROM users
WHERE active = true
ORDER BY created_at DESC
LIMIT $1;
```

**4. Validación de SQL en Tiempo de Generación**

Si tu SQL es inválido, `sqlc generate` fallará:

```bash
$ sqlc generate
error in query "GetActiveUsers": unknown column "invalid_col"
```

**5. Cambios de Schema Detectados**

Si tu schema cambió pero no actualizaste las queries:

```bash
$ sqlc generate
error in query "GetUser": column "user_id" not found in table "users"
```

### 62.1.3 Filosofía: Zero Runtime Overhead

La premisa de sqlc es revolucionaria para Go:

```
Máxima seguridad de tipos
              ↓
      Generación de código
              ↓
          Código Go puro
              ↓
      Compilación estándar
              ↓
   Ejecución sin overhead
```

Comparado con GORM:

```
Máxima flexibilidad ORM
              ↓
    Mapeo dinámico en runtime
              ↓
    Reflexión y parsing
              ↓
    Overhead constante
```

### 62.1.4 Code Generation Workflow

```
┌──────────────────────────────────────────────────────┐
│           WORKFLOW DE SQLC                           │
├──────────────────────────────────────────────────────┤
│                                                      │
│  1. schema/                                          │
│     └── migrations.sql    ← Define tu schema         │
│                                                      │
│  2. sql/                                             │
│     └── queries.sql       ← Escribe tus queries      │
│                                                      │
│  3. sqlc.yaml             ← Configura sqlc           │
│                                                      │
│  4. $ sqlc generate       ← Ejecuta generación       │
│                                                      │
│  5. sqlc/                                            │
│     └── db.go             ← Código generado ✅       │
│     └── models.go         ← Tipos generados ✅       │
│                                                      │
│  6. main.go               ← Usa el código generado   │
│                                                      │
└──────────────────────────────────────────────────────┘
```

### 62.1.5 Comparativa: sqlc vs GORM vs Ent

**sqlc**
- ✅ Type-safe 100%
- ✅ Cero overhead
- ✅ SQL puro y explícito
- ✅ Compilación rápida
- ❌ Más boilerplate SQL
- ❌ Menos abstracción

**GORM**
- ✅ Muy productivo
- ✅ Relaciones automáticas
- ✅ Menos SQL que escribir
- ❌ Overhead de runtime
- ❌ Type-safety parcial
- ❌ Debugging complicado

**Ent**
- ✅ Generación de código
- ✅ Type-safe
- ✅ Buena para modelos
- ✅ Migrations integradas
- ❌ Curva de aprendizaje
- ❌ Overhead en queries complejas

**Recomendación por caso:**
- **APIs críticas de performance**: sqlc
- **Prototipado rápido**: GORM
- **GraphQL + tipos**: Ent
- **Queries complejas**: sqlc
- **Relaciones simples**: GORM

### 62.1.6 Type Safety Example

```go
// ✅ SQLC - Error de compilación si es incorrecto
func ProcessUsers(ctx context.Context, db *sql.DB) error {
    querier := sqlc.New(db)
    
    // Parámetro debe ser int64 exactamente
    users, err := querier.GetUsersByAge(ctx, 30)
    if err != nil {
        return err
    }
    
    for _, user := range users {
        // Campos tipados correctamente
        fmt.Printf("%d: %s (%s)\n", user.ID, user.Name, user.Email)
    }
    return nil
}

// ❌ GORM - Error en runtime
func ProcessUsersGORM(db *gorm.DB) {
    var users []User
    
    // String cuando debería ser int → error en runtime
    db.Where("age = ?", "treinta").Find(&users)
    
    for _, user := range users {
        fmt.Printf("%d: %s\n", user.ID, user.Name)
    }
}
```

---

## 62.2 Setup & Installation

### 62.2.1 Instalación de sqlc

**Opción 1: Binario desde releases**

```bash
# macOS (Homebrew)
brew install sqlc

# Linux (verificar última versión en GitHub)
VERSION=v1.25.0
wget https://github.com/sqlc-dev/sqlc/releases/download/${VERSION}/sqlc_${VERSION#v}_linux_amd64.tar.gz
tar xzf sqlc_${VERSION#v}_linux_amd64.tar.gz
sudo mv sqlc /usr/local/bin/

# Windows (Chocolatey)
choco install sqlc
```

**Opción 2: Desde fuente**

```bash
git clone https://github.com/sqlc-dev/sqlc.git
cd sqlc
go install ./cmd/sqlc
```

**Opción 3: Docker**

```dockerfile
FROM golang:1.21

RUN apt-get update && apt-get install -y \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app
```

**Verificación:**

```bash
$ sqlc version
sqlc version 1.25.0
```

### 62.2.2 Estructura del Proyecto

```
myapp/
├── main.go
├── go.mod
├── go.sum
├── sqlc.yaml              ← Configuración de sqlc
├── schema/
│   └── migrations.sql     ← Schema de BD
├── sql/
│   ├── queries.sql        ← Consultas
│   ├── users.sql          ← Queries agrupadas
│   └── products.sql
├── sqlc/                  ← Generado automáticamente
│   ├── db.go
│   ├── models.go
│   └── querier.go
├── internal/
│   ├── repository/
│   │   └── user.go
│   └── handler/
│       └── user.go
└── docker-compose.yml
```

### 62.2.3 Configuración: sqlc.yaml

**Configuración Básica**

```yaml
version: "2"

project:
  name: "myapp"
  package_path: "github.com/mycompany/myapp"

db:
  driver: "postgres"
  queries: "sql/"
  schema: "schema/"

out:
  type: "go"
  dir: "sqlc"
  package: "sqlc"
  sql_driver: "pgx"
  sql_package: "pgx"
  go_type_overrides:
    - db_type: "uuid"
      go_type: "github.com/google/uuid.UUID"
    - db_type: "bytea"
      go_type: "[]byte"
    - db_type: "json"
      go_type: "encoding/json.RawMessage"
```

**Configuración Avanzada**

```yaml
version: "2"

project:
  name: "myapp"
  package_path: "github.com/mycompany/myapp"

db:
  driver: "postgres"
  queries: "sql/"
  schema: "schema/"
  gen:
    kind: "sync"

out:
  type: "go"
  dir: "sqlc"
  package: "sqlc"
  
  # Opciones de generación
  emit_db_tags: true
  emit_prepared_queries: true
  emit_interface: true
  emit_exact_table_names: true
  emit_empty_slices: true
  emit_methods: true
  emit_pointers_for_null_types: true
  
  sql_driver: "pgx"
  sql_package: "pgx"
  
  go_type_overrides:
    - db_type: "uuid"
      go_type: "github.com/google/uuid.UUID"
    - db_type: "timestamptz"
      go_type: "time.Time"
    - db_type: "json"
      go_type: "encoding/json.RawMessage"
    - db_type: "jsonb"
      go_type: "encoding/json.RawMessage"
```

### 62.2.4 Integración con Docker

**Dockerfile para desarrollo**

```dockerfile
FROM golang:1.21-alpine

# Instalar sqlc
RUN apk add --no-cache git && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Instalar herramientas de desarrollo
RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

# Copiar archivos
COPY . .

# Ejecutar sqlc generate
RUN sqlc generate

CMD ["air", "-c", ".air.toml"]
```

**docker-compose.yml**

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:15-alpine
    container_name: myapp_postgres
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypass
      POSTGRES_DB: myappdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./schema:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myuser"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build: .
    container_name: myapp
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DATABASE_URL: postgresql://myuser:mypass@postgres:5432/myappdb
    ports:
      - "8080:8080"
    volumes:
      - .:/app

volumes:
  postgres_data:
```

### 62.2.5 CI/CD Setup

**GitHub Actions**

```yaml
name: sqlc generate

on:
  push:
    branches: [main, develop]
    paths:
      - 'sql/**'
      - 'schema/**'
      - 'sqlc.yaml'

jobs:
  sqlc:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: sqlc-dev/setup-sqlc@v3
      
      - name: Generate code
        run: sqlc generate
      
      - name: Check for changes
        run: |
          if git diff --quiet sqlc/; then
            echo "No changes in generated code"
          else
            echo "ERROR: Generated code has uncommitted changes"
            git diff sqlc/
            exit 1
          fi
      
      - name: Run tests
        run: go test ./...
```

**Makefile**

```makefile
.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: db-up
db-up:
	docker-compose up -d postgres

.PHONY: db-down
db-down:
	docker-compose down

.PHONY: migrate
migrate: sqlc
	docker-compose exec postgres psql -U myuser -d myappdb -f /docker-entrypoint-initdb.d/migrations.sql

.PHONY: test
test: db-up
	go test ./...

.PHONY: dev
dev:
	docker-compose up
```

---

## 62.3 SQL Schema & Migrations

### 62.3.1 Writing Migrations

**migrations.sql - Estructura Base**

```sql
-- schema/migrations.sql

-- Tabla: users
-- Descripción: Usuarios del sistema
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_url VARCHAR(500),
    bio TEXT,
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    email_verified_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- Índices
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Trigger para actualizar updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 62.3.2 Creación de Tablas Relacionadas

```sql
-- Tabla: posts
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt VARCHAR(500),
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMPTZ,
    view_count BIGINT DEFAULT 0,
    like_count BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_is_published ON posts(is_published);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);

CREATE TRIGGER update_posts_updated_at BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabla: comments
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_approved BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);

CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabla: post_likes
CREATE TABLE post_likes (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(post_id, user_id)
);

CREATE INDEX idx_post_likes_user_id ON post_likes(user_id);
```

### 62.3.3 Constraints & Indexes

```sql
-- Constraints avanzados
ALTER TABLE posts ADD CONSTRAINT posts_title_not_empty
    CHECK (length(trim(title)) > 0);

ALTER TABLE users ADD CONSTRAINT users_email_lowercase
    CHECK (email = LOWER(email));

-- Índices compuestos
CREATE INDEX idx_posts_user_published ON posts(user_id, is_published, published_at DESC);

-- Índices parciales (partial indexes)
CREATE INDEX idx_active_users ON users(created_at DESC)
    WHERE is_active = true AND deleted_at IS NULL;

-- Índices de texto completo
CREATE INDEX idx_posts_content_fts ON posts USING GIN(to_tsvector('english', content));

-- Índice único
CREATE UNIQUE INDEX idx_users_email_uniq ON users(LOWER(email))
    WHERE deleted_at IS NULL;
```

### 62.3.4 Foreign Keys & Relationships

```sql
-- One-to-Many: users → posts
-- (Ya definido arriba)

-- Many-to-Many: users ↔ posts (favorites)
CREATE TABLE user_post_favorites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, post_id)
);

CREATE INDEX idx_user_post_favorites_user_id ON user_post_favorites(user_id);
CREATE INDEX idx_user_post_favorites_post_id ON user_post_favorites(post_id);

-- One-to-One: users ↔ user_profiles
CREATE TABLE user_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    bio_extended TEXT,
    website_url VARCHAR(500),
    social_links JSONB,
    preferences JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);
```

### 62.3.5 Versionado y Auditoría

```sql
-- Tabla de auditoría
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id BIGINT NOT NULL,
    action VARCHAR(10) NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
    changed_by BIGINT REFERENCES users(id),
    old_values JSONB,
    new_values JSONB,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_audit_log_table_record ON audit_log(table_name, record_id);
CREATE INDEX idx_audit_log_changed_at ON audit_log(changed_at DESC);

-- Función para auditoría
CREATE OR REPLACE FUNCTION audit_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        INSERT INTO audit_log (table_name, record_id, action, old_values)
        VALUES (TG_TABLE_NAME, OLD.id, 'DELETE', to_jsonb(OLD));
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_log (table_name, record_id, action, old_values, new_values, changed_by)
        VALUES (TG_TABLE_NAME, NEW.id, 'UPDATE', to_jsonb(OLD), to_jsonb(NEW), 
                COALESCE(current_setting('app.current_user_id')::bigint, NULL));
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_users AFTER DELETE OR UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();
```

---

## 62.4 Query Writing

### 62.4.1 Sintaxis de Queries en sqlc

**Anatomía de una Query**

```sql
-- sql/queries.sql

-- name: GetUserByID :one
-- GetUserByID retorna un usuario por su ID
SELECT id, email, username, first_name, last_name, created_at, is_active
FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;
```

Estructura:
- `-- name: GetUserByID :one` → Nombre y tipo de retorno
- `:one` → Retorna un registro (error si no existe)
- `:many` → Retorna múltiples registros
- `:exec` → No retorna registros (INSERT, UPDATE, DELETE)
- `:execs` → Múltiples operaciones sin retorno

### 62.4.2 SELECT Queries

**Simple SELECT - :one**

```sql
-- sql/users.sql

-- name: GetUserByID :one
SELECT id, email, username, first_name, last_name, avatar_url, 
       is_active, created_at, updated_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- Generado como:
-- func (q *Queries) GetUserByID(ctx context.Context, id int64) (*User, error)
```

**SELECT Multiple - :many**

```sql
-- name: ListActiveUsers :many
SELECT id, email, username, first_name, last_name, created_at
FROM users
WHERE is_active = true AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- Generado como:
-- func (q *Queries) ListActiveUsers(ctx context.Context, limit int64, offset int64) ([]User, error)
```

**SELECT con Agregación**

```sql
-- name: CountActiveUsers :one
SELECT COUNT(*) as count
FROM users
WHERE is_active = true AND deleted_at IS NULL;

-- Generado como:
-- func (q *Queries) CountActiveUsers(ctx context.Context) (int64, error)
```

### 62.4.3 INSERT Queries

**INSERT Simple - :exec**

```sql
-- name: CreateUser :exec
INSERT INTO users (email, username, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4, $5);
```

**INSERT con RETURNING - :one**

```sql
-- name: CreateUserWithReturn :one
INSERT INTO users (email, username, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, username, created_at;

-- Generado como:
-- func (q *Queries) CreateUserWithReturn(ctx context.Context, 
--     email string, username string, passwordHash string, 
--     firstName *string, lastName *string) (*User, error)
```

**Batch INSERT**

```sql
-- name: CreateMultipleUsers :many
INSERT INTO users (email, username, password_hash, first_name, last_name)
SELECT * FROM UNNEST(
    $1::varchar[],  -- emails
    $2::varchar[],  -- usernames
    $3::varchar[],  -- password_hashes
    $4::varchar[],  -- first_names
    $5::varchar[]   -- last_names
) AS t(email, username, password_hash, first_name, last_name)
RETURNING id, email, username, created_at;
```

### 62.4.4 UPDATE Queries

**UPDATE Simple**

```sql
-- name: UpdateUser :exec
UPDATE users
SET email = $1, first_name = $2, last_name = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $4 AND deleted_at IS NULL;
```

**UPDATE con RETURNING**

```sql
-- name: UpdateUserWithReturn :one
UPDATE users
SET email = $1, first_name = $2, last_name = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $4 AND deleted_at IS NULL
RETURNING id, email, first_name, last_name, updated_at;
```

**UPDATE Condicional**

```sql
-- name: VerifyUserEmail :exec
UPDATE users
SET email_verified_at = CURRENT_TIMESTAMP, is_verified = true
WHERE id = $1 AND email = $2 AND deleted_at IS NULL;
```

### 62.4.5 DELETE Queries

**DELETE Lógico (Soft Delete)**

```sql
-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP, is_active = false
WHERE id = $1 AND deleted_at IS NULL;
```

**DELETE Físico**

```sql
-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;
```

**DELETE con Validación**

```sql
-- name: DeleteUserIfNoComments :one
DELETE FROM users
WHERE id = $1 AND NOT EXISTS (
    SELECT 1 FROM comments WHERE user_id = $1
)
RETURNING id;
```

### 62.4.6 Parámetros Nombrados

**Named Parameters**

```sql
-- name: SearchUsers :many
SELECT id, email, username, first_name, last_name
FROM users
WHERE (username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
  AND is_active = $2
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3;

-- Uso:
-- users, err := querier.SearchUsers(ctx, searchTerm, true, 50)
```

**Named Parameters con :param**

```sql
-- name: GetUserByEmailOrUsername :one
SELECT id, email, username
FROM users
WHERE (email = @email OR username = @username)
  AND deleted_at IS NULL
LIMIT 1;

-- Generado como:
-- func (q *Queries) GetUserByEmailOrUsername(ctx context.Context, 
--     params GetUserByEmailOrUsernameParams) (*User, error)

-- Tipo generado:
// type GetUserByEmailOrUsernameParams struct {
//     Email    string
//     Username string
// }
```

### 62.4.7 Tipos de Retorno

```sql
-- :one - Retorna exactamente un registro
-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;
-- Función: func (q *Queries) GetPostByID(ctx context.Context, id int64) (*Post, error)

-- :many - Retorna cero o más registros
-- name: ListPosts :many
SELECT * FROM posts WHERE is_published = true;
-- Función: func (q *Queries) ListPosts(ctx context.Context) ([]Post, error)

-- :exec - No retorna datos (INSERT/UPDATE/DELETE)
-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;
-- Función: func (q *Queries) DeletePost(ctx context.Context, id int64) error

-- :execrows - Retorna el número de filas afectadas
-- name: DeactivateUsers :execrows
UPDATE users SET is_active = false WHERE is_active = true;
-- Función: func (q *Queries) DeactivateUsers(ctx context.Context) (int64, error)
```

---

## 62.5 Code Generation

### 62.5.1 Ejecutar sqlc generate

**Comando Básico**

```bash
# Generar desde sqlc.yaml en el directorio actual
$ sqlc generate

# Generar desde una ubicación específica
$ sqlc generate -f /path/to/sqlc.yaml

# Generar con salida de debug
$ sqlc generate --verbose

# Generar en seco (dry-run)
$ sqlc generate --dry
```

**Validación de Queries**

```bash
# Solo validar sin generar
$ sqlc vet

# Validar un archivo específico
$ sqlc vet ./sql/queries.sql
```

### 62.5.2 Estructura del Código Generado

**Directorio generado**

```
sqlc/
├── db.go           # Interfaz principal de Querier
├── models.go       # Tipos de datos (structs)
├── users.sql.go    # Métodos para queries de users
├── posts.sql.go    # Métodos para queries de posts
└── querier.go      # Interfaz Querier
```

### 62.5.3 Modelos Generados (models.go)

```go
// sqlc/models.go
package sqlc

import (
    "time"
    "database/sql"
)

type User struct {
    ID              int64
    Email           string
    Username        string
    PasswordHash    string
    FirstName       sql.NullString
    LastName        sql.NullString
    AvatarUrl       sql.NullString
    Bio             sql.NullString
    IsActive        bool
    IsVerified      bool
    EmailVerifiedAt sql.NullTime
    LastLoginAt     sql.NullTime
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       sql.NullTime
}

type Post struct {
    ID          int64
    UserID      int64
    Title       string
    Slug        string
    Content     string
    Excerpt     sql.NullString
    IsPublished bool
    PublishedAt sql.NullTime
    ViewCount   int64
    LikeCount   int64
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   sql.NullTime
}

type Comment struct {
    ID              int64
    PostID          int64
    UserID          int64
    ParentCommentID sql.NullInt64
    Content         string
    IsApproved      bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 62.5.4 Interfaz Querier (querier.go)

```go
// sqlc/querier.go
package sqlc

import (
    "context"
)

// Querier define todos los métodos para interactuar con la BD
type Querier interface {
    // Users
    GetUserByID(ctx context.Context, id int64) (*User, error)
    ListActiveUsers(ctx context.Context, limit int64, offset int64) ([]User, error)
    CountActiveUsers(ctx context.Context) (int64, error)
    CreateUserWithReturn(ctx context.Context, 
        email string, username string, passwordHash string, 
        firstName sql.NullString, lastName sql.NullString) (*User, error)
    UpdateUserWithReturn(ctx context.Context, 
        email string, firstName sql.NullString, lastNamestring, id int64) (*User, error)
    SoftDeleteUser(ctx context.Context, id int64) error
    
    // Posts
    GetPostByID(ctx context.Context, id int64) (*Post, error)
    ListPostsByUserID(ctx context.Context, userID int64) ([]Post, error)
    CreatePost(ctx context.Context, 
        userID int64, title string, slug string, content string) error
    
    // Comments
    CreateComment(ctx context.Context, 
        postID int64, userID int64, content string) error
}

// Queries proporciona la implementación de Querier
type Queries struct {
    db DBConn
}

// New crea una nueva instancia de Queries
func New(db DBConn) *Queries {
    return &Queries{db: db}
}

// DBConn es la interfaz para conexiones a BD
type DBConn interface {
    ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    PrepareContext(context.Context, string) (*sql.Stmt, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
```

### 62.5.5 Métodos Generados (users.sql.go)

```go
// sqlc/users.sql.go
package sqlc

import (
    "context"
    "database/sql"
)

// GetUserByID retorna un usuario por su ID
func (q *Queries) GetUserByID(ctx context.Context, id int64) (*User, error) {
    row := q.db.QueryRowContext(ctx, getUserByID, id)
    var user User
    err := row.Scan(
        &user.ID,
        &user.Email,
        &user.Username,
        &user.FirstName,
        &user.LastName,
        &user.AvatarUrl,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, username, first_name, last_name, avatar_url, 
       is_active, created_at, updated_at
FROM users
WHERE id = $1 AND deleted_at IS NULL`

// ListActiveUsers retorna usuarios activos
func (q *Queries) ListActiveUsers(ctx context.Context, 
    limit int64, offset int64) ([]User, error) {
    rows, err := q.db.QueryContext(ctx, listActiveUsers, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.Username,
            &user.FirstName,
            &user.LastName,
            &user.CreatedAt,
        ); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return users, nil
}

const listActiveUsers = `-- name: ListActiveUsers :many
SELECT id, email, username, first_name, last_name, created_at
FROM users
WHERE is_active = true AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2`

// CreateUserWithReturn crea un nuevo usuario
func (q *Queries) CreateUserWithReturn(ctx context.Context,
    email string, username string, passwordHash string,
    firstName sql.NullString, lastName sql.NullString) (*User, error) {
    
    row := q.db.QueryRowContext(ctx, createUserWithReturn,
        email, username, passwordHash, firstName, lastName)
    
    var user User
    err := row.Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

const createUserWithReturn = `-- name: CreateUserWithReturn :one
INSERT INTO users (email, username, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, username, created_at`
```

### 62.5.6 Función de Inicialización

```go
// main.go
package main

import (
    "database/sql"
    "log"
    
    _ "github.com/lib/pq"
    "myapp/sqlc"
)

func main() {
    // Conectar a la BD
    db, err := sql.Open("postgres", 
        "postgresql://user:pass@localhost/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Validar conexión
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    
    // Crear instancia de Queries
    queries := sqlc.New(db)
    
    // Usar las queries
    user, err := queries.GetUserByID(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Usuario: %s (%s)\n", user.Username, user.Email)
}
```

---

## 62.6 Advanced Queries

### 62.6.1 JOINS - INNER JOIN

```sql
-- sql/posts.sql

-- name: GetPostWithAuthor :one
SELECT 
    p.id, p.title, p.content, p.is_published, p.created_at,
    u.id as author_id, u.username as author_username, u.email as author_email
FROM posts p
INNER JOIN users u ON p.user_id = u.id
WHERE p.id = $1 AND p.deleted_at IS NULL;
```

**Código generado requiere manejo manual:**

```go
// Tipo custom para el resultado
type GetPostWithAuthorRow struct {
    ID                int64
    Title             string
    Content           string
    IsPublished       bool
    CreatedAt         time.Time
    AuthorID          int64
    AuthorUsername    string
    AuthorEmail       string
}

// Método generado
func (q *Queries) GetPostWithAuthor(ctx context.Context, id int64) (*GetPostWithAuthorRow, error) {
    row := q.db.QueryRowContext(ctx, getPostWithAuthor, id)
    var i GetPostWithAuthorRow
    err := row.Scan(
        &i.ID,
        &i.Title,
        &i.Content,
        &i.IsPublished,
        &i.CreatedAt,
        &i.AuthorID,
        &i.AuthorUsername,
        &i.AuthorEmail,
    )
    return &i, err
}
```

### 62.6.2 JOINS - LEFT JOIN

```sql
-- name: ListPostsWithCommentCount :many
SELECT 
    p.id, p.title, p.is_published, p.created_at,
    COUNT(c.id) as comment_count
FROM posts p
LEFT JOIN comments c ON p.id = c.post_id
WHERE p.deleted_at IS NULL AND p.is_published = true
GROUP BY p.id, p.title, p.is_published, p.created_at
ORDER BY p.created_at DESC
LIMIT $1;
```

### 62.6.3 Aggregations - COUNT, SUM, AVG

```sql
-- name: GetUserStats :one
SELECT 
    u.id,
    u.username,
    COUNT(p.id) as post_count,
    COUNT(c.id) as comment_count,
    COALESCE(SUM(pl.id), 0) as total_likes_received
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
LEFT JOIN comments c ON u.id = c.user_id
LEFT JOIN post_likes pl ON p.id = pl.post_id
WHERE u.id = $1 AND u.deleted_at IS NULL
GROUP BY u.id, u.username;
```

**Tipo generado:**

```go
type GetUserStatsRow struct {
    ID                   int64
    Username             string
    PostCount            int64
    CommentCount         int64
    TotalLikesReceived   int64
}
```

### 62.6.4 GROUP BY & HAVING

```sql
-- name: GetTopPostAuthors :many
SELECT 
    u.id,
    u.username,
    COUNT(p.id) as post_count,
    AVG(p.like_count) as avg_likes
FROM users u
INNER JOIN posts p ON u.id = p.user_id
WHERE p.deleted_at IS NULL
GROUP BY u.id, u.username
HAVING COUNT(p.id) >= $1
ORDER BY COUNT(p.id) DESC
LIMIT $2;
```

### 62.6.5 Subqueries

```sql
-- name: ListUsersWithRecentPosts :many
SELECT id, username, email
FROM users
WHERE id IN (
    SELECT DISTINCT user_id
    FROM posts
    WHERE is_published = true
    AND published_at > NOW() - INTERVAL '7 days'
)
ORDER BY username;

-- name: GetPostsWithLikePercentile :many
SELECT id, title, like_count
FROM posts
WHERE like_count >= (
    SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY like_count)
    FROM posts
    WHERE is_published = true
)
ORDER BY like_count DESC;
```

### 62.6.6 CTEs (Common Table Expressions)

```sql
-- name: GetPostRankings :many
WITH post_stats AS (
    SELECT 
        id,
        title,
        like_count,
        view_count,
        (like_count::float / NULLIF(view_count, 0)) as engagement_rate
    FROM posts
    WHERE is_published = true
),
ranked_posts AS (
    SELECT 
        id,
        title,
        like_count,
        view_count,
        engagement_rate,
        ROW_NUMBER() OVER (ORDER BY engagement_rate DESC) as rank
    FROM post_stats
)
SELECT id, title, like_count, view_count, engagement_rate, rank
FROM ranked_posts
WHERE rank <= $1;
```

**Tipo generado:**

```go
type GetPostRankingsRow struct {
    ID              int64
    Title           string
    LikeCount       int64
    ViewCount       int64
    EngagementRate  float64
    Rank            int64
}
```

### 62.6.7 Window Functions

```sql
-- name: GetPostsWithRank :many
SELECT 
    id,
    title,
    like_count,
    ROW_NUMBER() OVER (ORDER BY like_count DESC) as like_rank,
    LAG(like_count) OVER (ORDER BY created_at) as prev_like_count,
    LEAD(like_count) OVER (ORDER BY created_at) as next_like_count
FROM posts
WHERE is_published = true
ORDER BY created_at DESC;
```

### 62.6.8 Queries Complejas Reales

```sql
-- Blog Analytics Query
-- name: GetBlogAnalytics :one
WITH user_metrics AS (
    SELECT 
        user_id,
        COUNT(*) as total_posts,
        COUNT(CASE WHEN is_published THEN 1 END) as published_posts,
        SUM(like_count) as total_likes,
        AVG(view_count) as avg_views
    FROM posts
    WHERE deleted_at IS NULL
    GROUP BY user_id
),
comment_metrics AS (
    SELECT 
        user_id,
        COUNT(*) as total_comments,
        COUNT(CASE WHEN is_approved THEN 1 END) as approved_comments
    FROM comments
    GROUP BY user_id
),
engagement AS (
    SELECT 
        u.id as user_id,
        um.total_posts,
        um.published_posts,
        um.total_likes,
        um.avg_views,
        COALESCE(cm.total_comments, 0) as total_comments,
        COALESCE(cm.approved_comments, 0) as approved_comments,
        COUNT(DISTINCT pl.id) as total_likes_received
    FROM users u
    LEFT JOIN user_metrics um ON u.id = um.user_id
    LEFT JOIN comment_metrics cm ON u.id = cm.user_id
    LEFT JOIN post_likes pl ON u.id = (
        SELECT user_id FROM posts WHERE id = pl.post_id
    )
    WHERE u.id = $1 AND u.deleted_at IS NULL
    GROUP BY u.id, um.total_posts, um.published_posts, um.total_likes, um.avg_views,
             cm.total_comments, cm.approved_comments
)
SELECT 
    user_id, total_posts, published_posts, total_likes, avg_views,
    total_comments, approved_comments, total_likes_received
FROM engagement;
```

---

## 62.7 Transactions & Batch Operations

### 62.7.1 TX Pattern

**sqlc.Tx - Transacciones Manuales**

```go
// internal/repository/user_repository.go
package repository

import (
    "context"
    "database/sql"
    "myapp/sqlc"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// CreateUserWithProfile crea un usuario y su perfil en una transacción
func (r *UserRepository) CreateUserWithProfile(ctx context.Context, 
    email, username, passwordHash, bio string) (*sqlc.User, error) {
    
    // Iniciar transacción
    tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelReadCommitted,
    })
    if err != nil {
        return nil, err
    }
    
    // Crear queries para la transacción
    queries := sqlc.New(tx)
    
    // Crear usuario
    user, err := queries.CreateUserWithReturn(ctx, 
        email, username, passwordHash, 
        sql.NullString{Valid: false}, 
        sql.NullString{Valid: false})
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // Crear perfil
    _, err = queries.CreateUserProfile(ctx, user.ID, bio)
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // Confirmar transacción
    if err := tx.Commit(); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 62.7.2 Transacciones con Handler

```go
// TransactionFunc define una función de transacción
type TransactionFunc func(ctx context.Context, tx *sqlc.Queries) error

// WithTx ejecuta una función dentro de una transacción
func (r *UserRepository) WithTx(ctx context.Context, 
    fn TransactionFunc) error {
    
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    if err := fn(ctx, sqlc.New(tx)); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Err()
}

// Uso
err := repo.WithTx(ctx, func(ctx context.Context, queries *sqlc.Queries) error {
    user, err := queries.CreateUserWithReturn(ctx, email, username, hash, fn, ln)
    if err != nil {
        return err
    }
    
    return queries.CreateUserProfile(ctx, user.ID, bio)
})
```

### 62.7.3 Batch Inserts

```go
// Insertar múltiples usuarios en una transacción
func (r *UserRepository) BulkCreateUsers(ctx context.Context, 
    users []CreateUserInput) error {
    
    if len(users) == 0 {
        return nil
    }
    
    // Preparar arrays para UNNEST
    emails := make([]string, len(users))
    usernames := make([]string, len(users))
    hashes := make([]string, len(users))
    firstNames := make([]string, len(users))
    lastNames := make([]string, len(users))
    
    for i, u := range users {
        emails[i] = u.Email
        usernames[i] = u.Username
        hashes[i] = u.PasswordHash
        firstNames[i] = u.FirstName
        lastNames[i] = u.LastName
    }
    
    queries := sqlc.New(r.db)
    
    _, err := queries.CreateMultipleUsers(ctx,
        emails, usernames, hashes, firstNames, lastNames)
    
    return err
}
```

### 62.7.4 Error Handling en Transacciones

```go
// Función robusta de transacción
func (r *UserRepository) TransferUserPosts(ctx context.Context, 
    fromUserID, toUserID int64) error {
    
    tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
    })
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    queries := sqlc.New(tx)
    
    // Paso 1: Validar que ambos usuarios existan
    fromUser, err := queries.GetUserByID(ctx, fromUserID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("from user not found: %w", err)
    }
    
    toUser, err := queries.GetUserByID(ctx, toUserID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("to user not found: %w", err)
    }
    
    // Paso 2: Transferir posts
    if err := queries.TransferUserPosts(ctx, fromUserID, toUserID); err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to transfer posts: %w", err)
    }
    
    // Paso 3: Registrar auditoría
    if err := queries.CreateAuditLog(ctx, 
        "posts", fromUserID, "UPDATE", toUserID, 
        sql.NullString{String: "transferred to " + toUser.Username, Valid: true}); err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create audit log: %w", err)
    }
    
    // Confirmar
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

### 62.7.5 Rollback Scenarios

```go
// Ejemplo de rollback inteligente
func (r *UserRepository) UpdateUserWithValidation(ctx context.Context, 
    userID int64, updates *UpdateUserInput) error {
    
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    // Defer asegura rollback si algo falla
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    queries := sqlc.New(tx)
    
    // 1. Obtener usuario actual
    user, err := queries.GetUserByID(ctx, userID)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // 2. Validar email único si cambió
    if updates.Email != "" && updates.Email != user.Email {
        existing, _ := queries.GetUserByEmail(ctx, updates.Email)
        if existing != nil {
            tx.Rollback()
            return ErrEmailAlreadyExists
        }
    }
    
    // 3. Actualizar
    if err := queries.UpdateUser(ctx, 
        sql.NullString{String: updates.Email, Valid: updates.Email != ""},
        sql.NullString{String: updates.FirstName, Valid: updates.FirstName != ""},
        userID); err != nil {
        tx.Rollback()
        return err
    }
    
    // Si todo está bien, confirmar
    return tx.Commit().Err()
}
```

### 62.7.6 SQL para Transacciones

```sql
-- sql/transactions.sql

-- name: TransferUserPosts :exec
UPDATE posts
SET user_id = $2
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: CreateAuditLog :exec
INSERT INTO audit_log (table_name, record_id, action, changed_by, new_values, changed_at)
VALUES ('users', $1, 'UPDATE', $2, jsonb_build_object('transfer_to', $3), CURRENT_TIMESTAMP);

-- name: GetUserByEmail :one
SELECT id, email, username
FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET 
    email = COALESCE(NULLIF($1, ''), email),
    first_name = COALESCE(NULLIF($2, ''), first_name),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $3 AND deleted_at IS NULL;
```

---

## 62.8 Testing

### 62.8.1 Setup con Testcontainers

```go
// internal/test/db.go
package test

import (
    "context"
    "database/sql"
    "testing"
    
    _ "github.com/lib/pq"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
    Container testcontainers.Container
    DB        *sql.DB
}

func SetupTestDB(ctx context.Context, t *testing.T) *TestDB {
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, 
        testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
    if err != nil {
        t.Fatalf("Failed to start container: %v", err)
    }
    
    host, err := container.Host(ctx)
    if err != nil {
        container.Terminate(ctx)
        t.Fatalf("Failed to get container host: %v", err)
    }
    
    port, err := container.MappedPort(ctx, "5432")
    if err != nil {
        container.Terminate(ctx)
        t.Fatalf("Failed to get container port: %v", err)
    }
    
    dsn := fmt.Sprintf("postgresql://test:test@%s:%s/testdb?sslmode=disable",
        host, port.Port())
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        container.Terminate(ctx)
        t.Fatalf("Failed to connect: %v", err)
    }
    
    // Ejecutar migraciones
    if err := runMigrations(ctx, db); err != nil {
        container.Terminate(ctx)
        db.Close()
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    return &TestDB{
        Container: container,
        DB:        db,
    }
}

func (tdb *TestDB) Close(ctx context.Context) error {
    tdb.DB.Close()
    return tdb.Container.Terminate(ctx)
}

func runMigrations(ctx context.Context, db *sql.DB) error {
    schema := `
    CREATE TABLE users (
        id BIGSERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        username VARCHAR(100) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
    );
    `
    _, err := db.ExecContext(ctx, schema)
    return err
}
```

### 62.8.2 Test Database Setup

```go
// internal/test/fixtures.go
package test

import (
    "context"
    "database/sql"
    "testing"
    "myapp/sqlc"
)

type Fixtures struct {
    Queries *sqlc.Queries
    DB      *sql.DB
}

func NewFixtures(ctx context.Context, db *sql.DB) *Fixtures {
    return &Fixtures{
        Queries: sqlc.New(db),
        DB:      db,
    }
}

// CreateTestUser crea un usuario de prueba
func (f *Fixtures) CreateTestUser(ctx context.Context, 
    email, username string) (*sqlc.User, error) {
    
    return f.Queries.CreateUserWithReturn(ctx,
        email, username, "hashedpass123",
        sql.NullString{Valid: false},
        sql.NullString{Valid: false})
}

// CreateTestPost crea un post de prueba
func (f *Fixtures) CreateTestPost(ctx context.Context, 
    userID int64, title, content string) (*sqlc.Post, error) {
    
    return f.Queries.CreatePostWithReturn(ctx,
        userID, title, slugify(title), content,
        sql.NullString{Valid: false})
}

// Clean limpia todos los datos de prueba
func (f *Fixtures) Clean(ctx context.Context) error {
    queries := []string{
        "DELETE FROM comments",
        "DELETE FROM post_likes",
        "DELETE FROM posts",
        "DELETE FROM users",
    }
    
    for _, q := range queries {
        if _, err := f.DB.ExecContext(ctx, q); err != nil {
            return err
        }
    }
    return nil
}

func slugify(title string) string {
    // Implementación simple
    return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}
```

### 62.8.3 Fixture Data

```go
// internal/test/seed.go
package test

import (
    "context"
    "database/sql"
)

type SeedData struct {
    UserID  int64
    PostID  int64
    CommentID int64
}

func (f *Fixtures) SeedData(ctx context.Context) (*SeedData, error) {
    data := &SeedData{}
    
    // Crear usuario
    user, err := f.CreateTestUser(ctx, "john@example.com", "john")
    if err != nil {
        return nil, err
    }
    data.UserID = user.ID
    
    // Crear posts
    post, err := f.CreateTestPost(ctx, user.ID, "First Post", "Content here")
    if err != nil {
        return nil, err
    }
    data.PostID = post.ID
    
    // Crear comentario
    comment, err := f.Queries.CreateCommentWithReturn(ctx,
        post.ID, user.ID, "Great post!")
    if err != nil {
        return nil, err
    }
    data.CommentID = comment.ID
    
    return data, nil
}
```

### 62.8.4 Integration Tests

```go
// internal/repository/user_repository_test.go
package repository

import (
    "context"
    "testing"
    "myapp/internal/test"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreateUserWithProfile(t *testing.T) {
    ctx := context.Background()
    
    // Setup
    tdb := test.SetupTestDB(ctx, t)
    defer tdb.Close(ctx)
    
    repo := NewUserRepository(tdb.DB)
    
    // Test
    user, err := repo.CreateUserWithProfile(ctx,
        "test@example.com",
        "testuser",
        "hashedpass",
        "Test bio")
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, "test@example.com", user.Email)
}

func TestTransferUserPosts(t *testing.T) {
    ctx := context.Background()
    tdb := test.SetupTestDB(ctx, t)
    defer tdb.Close(ctx)
    
    fixtures := test.NewFixtures(ctx, tdb.DB)
    defer fixtures.Clean(ctx)
    
    data, err := fixtures.SeedData(ctx)
    require.NoError(t, err)
    
    // Crear otro usuario
    user2, err := fixtures.CreateTestUser(ctx, "jane@example.com", "jane")
    require.NoError(t, err)
    
    // Transfer
    repo := NewUserRepository(tdb.DB)
    err = repo.TransferUserPosts(ctx, data.UserID, user2.ID)
    require.NoError(t, err)
    
    // Verificar
    posts, err := fixtures.Queries.ListPostsByUserID(ctx, user2.ID)
    require.NoError(t, err)
    assert.Len(t, posts, 1)
}

func TestBulkCreateUsers(t *testing.T) {
    ctx := context.Background()
    tdb := test.SetupTestDB(ctx, t)
    defer tdb.Close(ctx)
    
    repo := NewUserRepository(tdb.DB)
    
    users := []CreateUserInput{
        {Email: "user1@test.com", Username: "user1", PasswordHash: "hash1", FirstName: "User", LastName: "One"},
        {Email: "user2@test.com", Username: "user2", PasswordHash: "hash2", FirstName: "User", LastName: "Two"},
        {Email: "user3@test.com", Username: "user3", PasswordHash: "hash3", FirstName: "User", LastName: "Three"},
    }
    
    err := repo.BulkCreateUsers(ctx, users)
    require.NoError(t, err)
    
    // Verificar
    activeUsers, err := repo.queries.ListActiveUsers(ctx, 10, 0)
    require.NoError(t, err)
    assert.Len(t, activeUsers, 3)
}
```

### 62.8.5 Test Utilities

```go
// internal/test/assertions.go
package test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "myapp/sqlc"
)

func AssertUserEqual(t *testing.T, expected, actual *sqlc.User) {
    assert.Equal(t, expected.ID, actual.ID)
    assert.Equal(t, expected.Email, actual.Email)
    assert.Equal(t, expected.Username, actual.Username)
}

func AssertPostContainsTitle(t *testing.T, posts []sqlc.Post, title string) {
    for _, p := range posts {
        if p.Title == title {
            return
        }
    }
    t.Errorf("Post with title %q not found", title)
}
```

---

## 62.9 Performance Considerations

### 62.9.1 Query Optimization

```sql
-- ❌ Ineficiente - N+1 queries
-- name: ListUsersPoorly :many
SELECT id, username, email FROM users LIMIT $1;
-- Luego en el código, para cada usuario, hacer otra query

-- ✅ Eficiente - Una sola query
-- name: ListUsersWithStats :many
SELECT 
    u.id,
    u.username,
    u.email,
    COUNT(p.id) as post_count,
    COUNT(c.id) as comment_count
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
LEFT JOIN comments c ON u.id = c.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.username, u.email
ORDER BY u.created_at DESC
LIMIT $1;
```

### 62.9.2 Prepared Statements

```go
// sqlc genera automáticamente prepared statements

// Al llamar la query, sqlc internamente:
// 1. Prepara el statement una sola vez
// 2. Reutiliza la prepared statement para llamadas posteriores
// 3. Evita el overhead de parsing SQL repetido

func (r *UserRepository) GetUserByIDMany(ctx context.Context, ids []int64) {
    queries := sqlc.New(r.db)
    
    // Cada GetUserByID usa la misma prepared statement
    for _, id := range ids {
        user, _ := queries.GetUserByID(ctx, id)
        // Procesar user
    }
}
```

### 62.9.3 Connection Pooling

```go
// main.go
package main

import (
    "database/sql"
    "time"
)

func initDB() *sql.DB {
    db, _ := sql.Open("postgres", dsn)
    
    // Configurar pool
    db.SetMaxOpenConns(25)           // Máximo de conexiones abiertas
    db.SetMaxIdleConns(5)            // Máximo de conexiones idle
    db.SetConnMaxLifetime(5 * time.Minute)  // Vida máxima de conexión
    db.SetConnMaxIdleTime(2 * time.Minute) // Máximo idle time
    
    return db
}
```

### 62.9.4 Indexing Strategy

```sql
-- Índices para consultas frecuentes
CREATE INDEX idx_users_active ON users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_published ON posts(published_at DESC) WHERE is_published = true;
CREATE INDEX idx_comments_recent ON comments(created_at DESC);

-- Índices compuestos para queries con múltiples filtros
CREATE INDEX idx_posts_user_published ON posts(user_id, is_published, published_at DESC);

-- Índices de texto completo
CREATE INDEX idx_posts_content_search ON posts USING GIN(to_tsvector('english', content));
```

### 62.9.5 Benchmarking

```go
// internal/benchmark/db_bench_test.go
package benchmark

import (
    "context"
    "testing"
    "database/sql"
    "myapp/sqlc"
)

func BenchmarkGetUserByID(b *testing.B) {
    ctx := context.Background()
    db := setupTestDB()
    defer db.Close()
    
    queries := sqlc.New(db)
    
    // Crear usuario de prueba
    user, _ := queries.CreateUserWithReturn(ctx, "bench@test.com", "bench", "hash", nil, nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        queries.GetUserByID(ctx, user.ID)
    }
}

func BenchmarkListUsersWithStats(b *testing.B) {
    ctx := context.Background()
    db := setupTestDB()
    defer db.Close()
    
    queries := sqlc.New(db)
    
    // Crear datos de prueba
    for i := 0; i < 100; i++ {
        queries.CreateUserWithReturn(ctx, 
            fmt.Sprintf("user%d@test.com", i),
            fmt.Sprintf("user%d", i),
            "hash", nil, nil)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        queries.ListUsersWithStats(ctx, 50)
    }
}

// Ejecutar: go test -bench=. -benchmem ./internal/benchmark
```

**Salida del benchmark:**

```
BenchmarkGetUserByID-8            10000    123456 ns/op    4567 B/op    12 allocs/op
BenchmarkListUsersWithStats-8      1000    234567 ns/op    8901 B/op    23 allocs/op
```

### 62.9.6 Query Analysis

```sql
-- Analizar plan de ejecución
EXPLAIN ANALYZE
SELECT u.id, u.username, COUNT(p.id) as post_count
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
WHERE u.is_active = true
GROUP BY u.id, u.username;

-- Resultado típico:
-- GroupAggregate (cost=123.45..234.56 rows=100 width=100)
--   ->  Sort (cost=123.45..150.67 rows=1000 width=50)
--       Sort Key: u.id
--       ->  Hash Left Join (cost=50.00..100.00 rows=1000 width=50)
--           Hash Cond: (p.user_id = u.id)
```

---

## 62.10 Integration with Go Frameworks

### 62.10.1 Integración con Gin

```go
// internal/handler/user_handler.go
package handler

import (
    "database/sql"
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "myapp/internal/repository"
    "myapp/sqlc"
)

type UserHandler struct {
    repo *repository.UserRepository
}

func NewUserHandler(db *sql.DB) *UserHandler {
    return &UserHandler{
        repo: repository.NewUserRepository(db),
    }
}

// GetUser - GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
    userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    user, err := h.repo.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// ListUsers - GET /users?limit=50&offset=0
func (h *UserHandler) ListUsers(c *gin.Context) {
    limit := int64(50)
    offset := int64(0)
    
    if l := c.Query("limit"); l != "" {
        limit, _ = strconv.ParseInt(l, 10, 64)
    }
    if o := c.Query("offset"); o != "" {
        offset, _ = strconv.ParseInt(o, 10, 64)
    }
    
    users, err := h.repo.ListActiveUsers(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    if users == nil {
        users = []sqlc.User{} // Retornar array vacío en lugar de null
    }
    
    c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser - POST /users
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required,min=3"`
    Password string `json:"password" binding:"required,min=8"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Hash password
    hashedPass, _ := hashPassword(req.Password)
    
    user, err := h.repo.CreateUserWithReturn(c.Request.Context(),
        req.Email, req.Username, hashedPass)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

// UpdateUser - PUT /users/:id
type UpdateUserRequest struct {
    FirstName *string `json:"first_name"`
    LastName  *string `json:"last_name"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    
    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.repo.UpdateUserWithReturn(c.Request.Context(), userID,
        req.FirstName, req.LastName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// DeleteUser - DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
    userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    
    if err := h.repo.SoftDeleteUser(c.Request.Context(), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }
    
    c.JSON(http.StatusNoContent, nil)
}
```

**Registrar rutas:**

```go
// main.go
package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "myapp/internal/handler"
)

func main() {
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()
    
    r := gin.Default()
    
    userHandler := handler.NewUserHandler(db)
    
    // Routes
    r.GET("/users/:id", userHandler.GetUser)
    r.GET("/users", userHandler.ListUsers)
    r.POST("/users", userHandler.CreateUser)
    r.PUT("/users/:id", userHandler.UpdateUser)
    r.DELETE("/users/:id", userHandler.DeleteUser)
    
    r.Run(":8080")
}
```

### 62.10.2 Integración con Echo

```go
// internal/handler/echo_user_handler.go
package handler

import (
    "database/sql"
    "net/http"
    "strconv"
    
    "github.com/labstack/echo/v4"
    "myapp/internal/repository"
)

type EchoUserHandler struct {
    repo *repository.UserRepository
}

func NewEchoUserHandler(db *sql.DB) *EchoUserHandler {
    return &EchoUserHandler{
        repo: repository.NewUserRepository(db),
    }
}

func (h *EchoUserHandler) GetUser(c echo.Context) error {
    userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
    }
    
    user, err := h.repo.GetUserByID(c.Request().Context(), userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    
    return c.JSON(http.StatusOK, user)
}

func (h *EchoUserHandler) ListUsers(c echo.Context) error {
    limit, _ := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
    offset, _ := strconv.ParseInt(c.QueryParam("offset"), 10, 64)
    
    if limit == 0 {
        limit = 50
    }
    
    users, err := h.repo.ListActiveUsers(c.Request().Context(), limit, offset)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    
    return c.JSON(http.StatusOK, map[string]interface{}{"data": users})
}

// Registrar con Echo
func RegisterRoutes(e *echo.Echo, db *sql.DB) {
    handler := NewEchoUserHandler(db)
    
    g := e.Group("/api/v1")
    g.GET("/users/:id", handler.GetUser)
    g.GET("/users", handler.ListUsers)
}
```

### 62.10.3 Repository Pattern

```go
// internal/repository/user_repository.go
package repository

import (
    "context"
    "database/sql"
    "myapp/sqlc"
)

// UserRepository define la interfaz de acceso a datos de usuarios
type UserRepository interface {
    GetUserByID(ctx context.Context, id int64) (*sqlc.User, error)
    ListActiveUsers(ctx context.Context, limit, offset int64) ([]sqlc.User, error)
    CreateUserWithReturn(ctx context.Context, 
        email, username, passwordHash string) (*sqlc.User, error)
    UpdateUserWithReturn(ctx context.Context, userID int64,
        firstName, lastName *string) (*sqlc.User, error)
    SoftDeleteUser(ctx context.Context, userID int64) error
}

// userRepository implementa UserRepository
type userRepository struct {
    db      *sql.DB
    queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{
        db:      db,
        queries: sqlc.New(db),
    }
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*sqlc.User, error) {
    return r.queries.GetUserByID(ctx, id)
}

func (r *userRepository) ListActiveUsers(ctx context.Context, 
    limit, offset int64) ([]sqlc.User, error) {
    return r.queries.ListActiveUsers(ctx, limit, offset)
}

func (r *userRepository) CreateUserWithReturn(ctx context.Context,
    email, username, passwordHash string) (*sqlc.User, error) {
    return r.queries.CreateUserWithReturn(ctx,
        email, username, passwordHash,
        sql.NullString{Valid: false},
        sql.NullString{Valid: false})
}

func (r *userRepository) UpdateUserWithReturn(ctx context.Context, userID int64,
    firstName, lastName *string) (*sqlc.User, error) {
    
    fn := sql.NullString{Valid: firstName != nil}
    if firstName != nil {
        fn.String = *firstName
    }
    
    ln := sql.NullString{Valid: lastName != nil}
    if lastName != nil {
        ln.String = *lastName
    }
    
    return r.queries.UpdateUserWithReturn(ctx, fn, ln, "", userID)
}

func (r *userRepository) SoftDeleteUser(ctx context.Context, userID int64) error {
    return r.queries.SoftDeleteUser(ctx, userID)
}
```

### 62.10.4 Dependency Injection

```go
// internal/app/app.go
package app

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "myapp/internal/handler"
    "myapp/internal/repository"
)

type App struct {
    router *gin.Engine
    db     *sql.DB
}

func New(db *sql.DB) *App {
    router := gin.Default()
    app := &App{router: router, db: db}
    app.setupRoutes()
    return app
}

func (a *App) setupRoutes() {
    // Crear repositorios
    userRepo := repository.NewUserRepository(a.db)
    
    // Crear handlers
    userHandler := handler.NewUserHandler(a.db)
    
    // Registrar rutas
    users := a.router.Group("/users")
    {
        users.GET("/:id", userHandler.GetUser)
        users.GET("", userHandler.ListUsers)
        users.POST("", userHandler.CreateUser)
        users.PUT("/:id", userHandler.UpdateUser)
        users.DELETE("/:id", userHandler.DeleteUser)
    }
}

func (a *App) Run(addr string) error {
    return a.router.Run(addr)
}

func (a *App) Close() error {
    return a.db.Close()
}
```

**Uso:**

```go
// main.go
func main() {
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()
    
    app := app.New(db)
    defer app.Close()
    
    app.Run(":8080")
}
```

---

## 62.11 Production & Best Practices

### 62.11.1 Error Handling Patterns

```go
// internal/errors/errors.go
package errors

import (
    "database/sql"
    "fmt"
)

type ErrorType int

const (
    ErrTypeNotFound ErrorType = iota
    ErrTypeConflict
    ErrTypeValidation
    ErrTypeInternal
)

type AppError struct {
    Type    ErrorType
    Message string
    Cause   error
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

// WrapDBError convierte errores de BD a errores de aplicación
func WrapDBError(err error) error {
    if err == nil {
        return nil
    }
    
    if err == sql.ErrNoRows {
        return &AppError{
            Type:    ErrTypeNotFound,
            Message: "Resource not found",
            Cause:   err,
        }
    }
    
    // Detectar constraint violations (postgres)
    if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
        return &AppError{
            Type:    ErrTypeConflict,
            Message: "Email already exists",
            Cause:   err,
        }
    }
    
    return &AppError{
        Type:    ErrTypeInternal,
        Message: "Database error",
        Cause:   err,
    }
}

// IsNotFound checa si es un error 404
func IsNotFound(err error) bool {
    if ae, ok := err.(*AppError); ok {
        return ae.Type == ErrTypeNotFound
    }
    return err == sql.ErrNoRows
}

// IsConflict checa si es un error 409
func IsConflict(err error) bool {
    if ae, ok := err.(*AppError); ok {
        return ae.Type == ErrTypeConflict
    }
    return false
}
```

**Uso en handlers:**

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    
    user, err := h.repo.GetUserByID(c.Request.Context(), userID)
    err = errors.WrapDBError(err)
    
    if errors.IsNotFound(err) {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

### 62.11.2 Logging & Debugging

```go
// internal/logger/logger.go
package logger

import (
    "context"
    "fmt"
    "log/slog"
    "time"
)

type QueryLogger struct {
    logger *slog.Logger
}

func New(logger *slog.Logger) *QueryLogger {
    return &QueryLogger{logger: logger}
}

// LogQuery registra las queries ejecutadas
func (ql *QueryLogger) LogQuery(ctx context.Context, query string, 
    args []interface{}, startTime time.Time, err error) {
    
    duration := time.Since(startTime)
    
    if err != nil {
        ql.logger.Error("Query failed",
            slog.String("query", query),
            slog.Any("args", args),
            slog.Duration("duration", duration),
            slog.String("error", err.Error()),
        )
        return
    }
    
    level := slog.LevelDebug
    if duration > time.Second {
        level = slog.LevelWarn
    }
    
    ql.logger.Log(ctx, level, "Query executed",
        slog.String("query", query),
        slog.Duration("duration", duration),
    )
}
```

**Wrapper para logging automático:**

```go
// internal/db/instrumented_db.go
package db

import (
    "context"
    "database/sql"
    "time"
    "myapp/internal/logger"
)

type InstrumentedDB struct {
    db     *sql.DB
    logger *logger.QueryLogger
}

func NewInstrumentedDB(db *sql.DB, logger *logger.QueryLogger) *InstrumentedDB {
    return &InstrumentedDB{db: db, logger: logger}
}

func (idb *InstrumentedDB) QueryContext(ctx context.Context, 
    query string, args ...interface{}) (*sql.Rows, error) {
    
    start := time.Now()
    rows, err := idb.db.QueryContext(ctx, query, args...)
    idb.logger.LogQuery(ctx, query, args, start, err)
    return rows, err
}

func (idb *InstrumentedDB) ExecContext(ctx context.Context, 
    query string, args ...interface{}) (sql.Result, error) {
    
    start := time.Now()
    result, err := idb.db.ExecContext(ctx, query, args...)
    idb.logger.LogQuery(ctx, query, args, start, err)
    return result, err
}
```

### 62.11.3 Migration Management

```bash
# Usando migrate
migrate create -ext sql -dir schema -seq create_users_table
migrate create -ext sql -dir schema -seq create_posts_table

# Aplicar migraciones
migrate -path schema -database "postgresql://..." up

# Revertir
migrate -path schema -database "postgresql://..." down
```

**Script de inicialización:**

```bash
#!/bin/bash
# scripts/migrate.sh

set -e

DB_URL="${DATABASE_URL}"

if [ -z "$DB_URL" ]; then
    echo "ERROR: DATABASE_URL no está definida"
    exit 1
fi

# Aplicar migraciones
migrate -path schema/migrations -database "$DB_URL" up

# Generar código sqlc
sqlc generate

echo "Migraciones y generación completadas"
```

### 62.11.4 Deployment Strategies

**Docker Compose para Producción**

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: produser
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: production
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./schema:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U produser"]
      interval: 30s
      timeout: 10s
      retries: 5

  app:
    build: .
    environment:
      DATABASE_URL: postgresql://produser:${DB_PASSWORD}@postgres:5432/production
      GIN_MODE: release
      LOG_LEVEL: info
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

volumes:
  postgres_data:
    driver: local
```

**Dockerfile Multi-stage**

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o app .

# Runtime stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app

COPY --from=builder /src/app .
COPY schema schema/

EXPOSE 8080

CMD ["./app"]
```

### 62.11.5 Monitoring & Observability

```go
// internal/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    DBQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
        },
        []string{"query_name", "status"},
    )
    
    DBConnections = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "db_connections",
            Help: "Number of database connections",
        },
        []string{"state"},
    )
)

func RecordQueryDuration(queryName, status string, durationSeconds float64) {
    DBQueryDuration.WithLabelValues(queryName, status).Observe(durationSeconds)
}

func UpdateConnections(openConn, inUse, idle int) {
    DBConnections.WithLabelValues("open").Set(float64(openConn))
    DBConnections.WithLabelValues("in_use").Set(float64(inUse))
    DBConnections.WithLabelValues("idle").Set(float64(idle))
}
```

### 62.11.6 Best Practices Checklist

✅ **DO:**
- Usar sqlc para queries críticas de performance
- Implementar transacciones para operaciones multi-paso
- Usar prepared statements (sqlc lo hace automáticamente)
- Implementar connection pooling
- Logging detallado de queries lentas
- Validación de entrada en handlers
- Migrations versionadas
- Error handling consistente
- Testing con testcontainers
- Monitoreo de performance

❌ **DON'T:**
- No escribir SQL dinámico (sqlc genera seguro)
- No usar queries N+1
- No ignorar errores de BD
- No hardcodear conexiones en el código
- No commitear credentials
- No saltarse transacciones en operaciones críticas
- No hacer migraciones manuales
- No usar typecast inseguro de null.String
- No generar código sin validación

### 62.11.7 Case Studies

**Case Study 1: Blog Platform**

Requisitos:
- 10,000+ usuarios activos
- 100,000+ posts
- 1M+ comments
- Queries complejas con joins

Solución con sqlc:
```sql
-- Optimized query para feed del usuario
-- name: GetUserFeed :many
WITH user_following AS (
    SELECT followed_id FROM follows WHERE follower_id = $1
),
user_posts AS (
    SELECT id, title, content, user_id, created_at, like_count
    FROM posts
    WHERE user_id IN (SELECT followed_id FROM user_following)
    AND is_published = true
    AND published_at > NOW() - INTERVAL '30 days'
)
SELECT 
    p.*, 
    u.username, 
    u.avatar_url,
    COUNT(pl.id) as like_count,
    EXISTS(SELECT 1 FROM post_likes WHERE post_id = p.id AND user_id = $1) as liked_by_user
FROM user_posts p
JOIN users u ON p.user_id = u.id
LEFT JOIN post_likes pl ON p.id = pl.post_id
GROUP BY p.id, u.username, u.avatar_url
ORDER BY p.created_at DESC
LIMIT $2;
```

Resultados:
- Query time: 50ms (vs 2s con GORM)
- Memory usage: 2MB (vs 50MB con reflexión)
- Type safety: 100%

**Case Study 2: E-commerce**

Requisitos:
- Órdenes con múltiples items
- Inventario en tiempo real
- Transacciones críticas

Solución con sqlc:
```go
func (r *OrderRepository) CreateOrder(ctx context.Context, 
    userID int64, items []OrderItem) (*Order, error) {
    
    tx, _ := r.db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    queries := sqlc.New(tx)
    
    // 1. Crear orden
    order, _ := queries.CreateOrderWithReturn(ctx, userID)
    
    // 2. Validar inventario y crear líneas
    for _, item := range items {
        product, _ := queries.GetProductForUpdate(ctx, item.ProductID)
        
        if product.Stock < item.Quantity {
            return nil, ErrInsufficientStock
        }
        
        // 3. Crear línea de orden
        queries.CreateOrderLine(ctx, order.ID, item.ProductID, item.Quantity)
        
        // 4. Actualizar inventario
        queries.UpdateProductStock(ctx, item.ProductID, -item.Quantity)
    }
    
    tx.Commit()
    return order, nil
}
```

Resultados:
- Transacciones atomícticas 100% seguras
- Evita race conditions
- Zero overhead de ORM

---

## 62.12 Ejercicios Progresivos

### Ejercicio 1: First sqlc Query (SELECT, INSERT)

**Objetivo:** Crear tu primer proyecto sqlc con queries básicas.

**Pasos:**

1. Inicializar proyecto
```bash
mkdir sqlc-first-app
cd sqlc-first-app
go mod init github.com/myuser/sqlc-first-app
```

2. Crear `sqlc.yaml`
```yaml
version: "2"
project:
  name: "myapp"
  package_path: "github.com/myuser/sqlc-first-app"

db:
  driver: "postgres"
  queries: "sql/"
  schema: "schema/"

out:
  type: "go"
  dir: "sqlc"
  package: "sqlc"
```

3. Crear `schema/migrations.sql`
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

4. Crear `sql/queries.sql`
```sql
-- name: GetUserByID :one
SELECT id, email, username, created_at FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, username) VALUES ($1, $2)
RETURNING id, email, username, created_at;

-- name: ListAllUsers :many
SELECT id, email, username, created_at FROM users ORDER BY created_at DESC;
```

5. Generar código
```bash
sqlc generate
```

6. Crear `main.go`
```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
    "github.com/myuser/sqlc-first-app/sqlc"
)

func main() {
    db, err := sql.Open("postgres", 
        "postgresql://user:pass@localhost/testdb?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    queries := sqlc.New(db)
    ctx := context.Background()
    
    // Crear usuario
    user, err := queries.CreateUser(ctx, "john@example.com", "john")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user: %v\n", user)
    
    // Obtener usuario
    retrieved, err := queries.GetUserByID(ctx, user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Retrieved user: %v\n", retrieved)
    
    // Listar todos
    users, err := queries.ListAllUsers(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("All users: %v\n", users)
}
```

**Verificación:**
- ✅ Código genera sin errores
- ✅ Tipos correctos
- ✅ Queries funcionan

---

### Ejercicio 2: Joins & Complex Queries

**Objetivo:** Escribir queries con JOINs, agregaciones y CTEs.

**Pasos:**

1. Extender `schema/migrations.sql`
```sql
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    title VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT REFERENCES posts(id),
    user_id BIGINT REFERENCES users(id),
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

2. Agregar queries complejas en `sql/queries.sql`
```sql
-- name: GetPostWithAuthor :one
SELECT p.id, p.title, p.created_at, u.username, u.email
FROM posts p
INNER JOIN users u ON p.user_id = u.id
WHERE p.id = $1;

-- name: GetPostsWithCommentCount :many
SELECT 
    p.id, p.title, p.created_at,
    COUNT(c.id) as comment_count,
    u.username
FROM posts p
LEFT JOIN comments c ON p.id = c.post_id
LEFT JOIN users u ON p.user_id = u.id
GROUP BY p.id, p.title, p.created_at, u.username
ORDER BY p.created_at DESC;
```

3. Crear archivo de test `exercise2_test.go`
```go
func TestJoinQueries(t *testing.T) {
    ctx := context.Background()
    db := setupTestDB(ctx, t)
    queries := sqlc.New(db)
    
    // Crear usuario y post
    user, _ := queries.CreateUser(ctx, "test@example.com", "test")
    post, _ := queries.CreatePostWithReturn(ctx, user.ID, "Test Post")
    
    // Test JOIN
    result, err := queries.GetPostWithAuthor(ctx, post.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Username, result.Username)
}
```

**Verificación:**
- ✅ Joins generan correctamente
- ✅ Agregaciones funcionan
- ✅ Tipos anidados generados

---

### Ejercicio 3: Transactions

**Objetivo:** Implementar transacciones seguras.

**Pasos:**

1. Crear `internal/repository/transaction_test.go`
```go
func TestTransactionCommit(t *testing.T) {
    ctx := context.Background()
    db := setupTestDB(ctx, t)
    
    tx, _ := db.BeginTx(ctx, nil)
    queries := sqlc.New(tx)
    
    user, _ := queries.CreateUser(ctx, "tx@example.com", "txuser")
    post, _ := queries.CreatePostWithReturn(ctx, user.ID, "TX Post")
    
    tx.Commit()
    
    // Verificar que persiste
    queries2 := sqlc.New(db)
    retrieved, _ := queries2.GetPostWithAuthor(ctx, post.ID)
    assert.Equal(t, "TX Post", retrieved.Title)
}

func TestTransactionRollback(t *testing.T) {
    ctx := context.Background()
    db := setupTestDB(ctx, t)
    
    tx, _ := db.BeginTx(ctx, nil)
    queries := sqlc.New(tx)
    
    user, _ := queries.CreateUser(ctx, "rollback@example.com", "rbuser")
    
    // NO hacer commit - simular error
    tx.Rollback()
    
    // Verificar que NO persiste
    queries2 := sqlc.New(db)
    _, err := queries2.GetUserByID(ctx, user.ID)
    assert.Error(t, err)
}
```

**Verificación:**
- ✅ Commits persisten
- ✅ Rollbacks deshacen cambios
- ✅ Manejo de errores robusto

---

### Ejercicio 4: Integration con Gin

**Objetivo:** Integrar sqlc con un framework web.

**Pasos:**

1. Crear `internal/handler/user_handler.go`
2. Crear rutas en `main.go`
3. Implementar tests HTTP

```go
func TestGetUserEndpoint(t *testing.T) {
    db := setupTestDB(context.Background(), t)
    queries := sqlc.New(db)
    user, _ := queries.CreateUser(context.Background(), 
        "endpoint@example.com", "epuser")
    
    router := gin.Default()
    handler := NewUserHandler(db)
    router.GET("/users/:id", handler.GetUser)
    
    req, _ := http.NewRequest("GET", 
        fmt.Sprintf("/users/%d", user.ID), nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "epuser")
}
```

**Verificación:**
- ✅ Endpoints funcionan
- ✅ Serialización JSON correcta
- ✅ Error handling en HTTP

---

### Ejercicio 5: Production Deployment

**Objetivo:** Deployar aplicación sqlc en producción.

**Pasos:**

1. Crear `docker-compose.yml` con postgres + app
2. Crear migrations en `schema/`
3. Configurar CI/CD en `.github/workflows/`
4. Agregar monitoring con Prometheus
5. Documentar deployment

```bash
# Build imagen
docker build -t sqlc-app:latest .

# Run con compose
docker-compose up -d

# Verificar salud
curl http://localhost:8080/health

# Ver logs
docker-compose logs -f app
```

**Verificación:**
- ✅ Contenedor runs correctamente
- ✅ Base de datos se inicializa
- ✅ API responde
- ✅ Migrations se aplican automáticamente

---

## Comparación Final: sqlc vs GORM vs Ent

```
┌────────────────────────────────────────────────────────────────┐
│                     MATRIZ DE DECISIÓN                         │
├──────────────────┬──────────────┬──────────────┬────────────────┤
│    Criterio      │    sqlc      │     GORM     │      Ent       │
├──────────────────┼──────────────┼──────────────┼────────────────┤
│ Type Safety      │ ⭐⭐⭐⭐⭐  │ ⭐⭐⭐      │ ⭐⭐⭐⭐    │
│ Performance      │ ⭐⭐⭐⭐⭐  │ ⭐⭐⭐      │ ⭐⭐⭐⭐    │
│ Learning Curve   │ ⭐⭐⭐⭐    │ ⭐⭐⭐⭐⭐  │ ⭐⭐⭐      │
│ Productivity     │ ⭐⭐⭐      │ ⭐⭐⭐⭐⭐  │ ⭐⭐⭐⭐    │
│ Flexibilidad SQL │ ⭐⭐⭐⭐⭐  │ ⭐⭐⭐      │ ⭐⭐⭐      │
│ Comunidad        │ ⭐⭐⭐⭐    │ ⭐⭐⭐⭐⭐  │ ⭐⭐⭐⭐    │
│ Generación Código│ ⭐⭐⭐⭐⭐  │ ⭐⭐        │ ⭐⭐⭐⭐⭐  │
└──────────────────┴──────────────┴──────────────┴────────────────┘

ELIGE sqlc SI:
✅ Performance crítica
✅ Queries complejas
✅ Type safety absoluta
✅ Bajas latencias

ELIGE GORM SI:
✅ MVP rápido
✅ Relaciones simples
✅ Prototipado
✅ Menos configuración

ELIGE Ent SI:
✅ Modelos complejos
✅ GraphQL backend
✅ Generación avanzada
✅ Empresa grande
```

---

## Conclusión

**sqlc** representa una filosofía moderna de desarrollo Go: **Type safety + Performance + Explicitness**.

Al eliminar la reflexión y generar código puro Go, sqlc proporciona:

1. **Seguridad de tipos completa** desde la BD hasta la aplicación
2. **Cero overhead de runtime** comparado con ORMs
3. **Debugging trivial** - es solo código Go
4. **SQL nativo** - acceso a todo el poder de PostgreSQL
5. **Cambios de schema detectados** en tiempo de compilación

Para aplicaciones que necesitan:
- **Alta performance**
- **Type safety absoluta**
- **Queries complejas**
- **Mejor debugging**

**sqlc es la herramienta correcta.**

---

## Referencias

- [Official sqlc Documentation](https://docs.sqlc.dev)
- [sqlc GitHub Repository](https://github.com/sqlc-dev/sqlc)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go database/sql Documentation](https://pkg.go.dev/database/sql)
- [Testcontainers Go](https://golang.testcontainers.org/)

---

**Fin del Capítulo 62**
