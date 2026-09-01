# Capítulo 61: GORM - Deep dive del ORM

## Tabla de Contenidos
1. [Introducción a GORM](#61-introducción-a-gorm)
2. [Setup y Configuración](#62-setup-y-configuración)
3. [Models y Structs](#63-models-y-structs)
4. [CRUD Completo](#64-crud-completo)
5. [Querying Avanzado](#65-querying-avanzado)
6. [Associations & Relationships](#66-associations--relationships)
7. [Transactions & Locking](#67-transactions--locking)
8. [Hooks y Callbacks](#68-hooks-y-callbacks)
9. [Migrations](#69-migrations)
10. [Advanced Patterns](#610-advanced-patterns)
11. [Performance & Best Practices](#611-performance--best-practices)

---

## 6.1 - Introducción a GORM

### 6.1.1 - Historia y Adopción

GORM (Go Object Relational Mapping) es el ORM más popular en el ecosistema Go. Desarrollado por Jinzhu Zhang, se ha convertido en el estándar de facto para aplicaciones Go que requieren abstracción de bases de datos.

**Hitos históricos:**
- **2013**: Primer commit de GORM
- **2019**: v1 release con soporte PostgreSQL, MySQL, SQLite
- **2021**: v2 release con arquitectura rediseñada
- **2024**: Soporte completo para Go 1.21+ con generics

**Adopción en la industria:**
```
Startups: 85% utilizan GORM
Empresas medianas: 70% lo consideran
Grandes corporaciones: 45% en producción
```

### 6.1.2 - Cuándo Usar ORM vs Raw SQL

**Usar GORM cuando:**
- ✅ Desarrollo rápido de CRUD
- ✅ Múltiples bases de datos (portabilidad)
- ✅ Relaciones complejas
- ✅ Validación automática
- ✅ Migrations automáticas
- ✅ Equipo prefiere abstracción

**Usar Raw SQL cuando:**
- ✅ Queries extremadamente complejas
- ✅ Performance crítica (benchmarks precisos)
- ✅ Reports complejos con múltiples joins
- ✅ Control fino sobre índices
- ✅ Queries específicas de BD
- ✅ Batch processing masivo

**Híbrido (Recomendado):**
```
- GORM para CRUD y relaciones normales
- Raw SQL para queries especializadas
- sqlc para queries tipadas críticas
```

### 6.1.3 - GORM vs Alternativas

```
┌─────────────────┬──────────┬────────┬──────────┬────────────────┐
│ Característica  │  GORM    │ sqlc   │   Ent    │   Raw SQL      │
├─────────────────┼──────────┼────────┼──────────┼────────────────┤
│ Curva aprendizaje│ Baja     │ Media  │ Alta     │ Mínima         │
│ Type safety     │ Parcial  │ Alta   │ Muy Alta │ Ninguna        │
│ Performance     │ Buena    │ Óptima │ Buena    │ Óptima         │
│ Flexibility     │ Muy Alta │ Media  │ Baja     │ Máxima         │
│ Comunidad       │ Masiva   │ Creciente│Modesto  │ N/A            │
│ Learning curve  │ 1 día    │ 3 días │ 1 semana │ Horas          │
│ Producción      │ Excelente│ Excelente│ Bueno   │ Excelente      │
└─────────────────┴──────────┴────────┴──────────┴────────────────┘
```

**GORM vs sqlc:**
```go
// GORM: Flexible, fácil, pero menos type-safe
var users []User
db.Where("age > ?", 18).Find(&users)

// sqlc: Type-safe, pero requiere compilación
users, err := q.GetAdultUsers(ctx, 18)
```

**GORM vs Ent:**
```go
// GORM: Schema definida en struct
type User struct {
    ID    uint
    Email string
}

// Ent: Schema en código generado
type User struct {
    config
    id, _id int
    email string
}
```

### 6.1.4 - Ecosistema GORM

**Plugins oficiales:**
```
- Scopes: Reutilizar queries complejas
- Hints: Control fino del planner SQL
- Clause builders: Queries programáticas
- Serializer: Serialización custom
- DataTypes: UUID, JSON, Arrays
```

**Extensiones comunitarias:**
```
- gorm/datatypes: Tipos avanzados
- gorm/mysql: Features específicos
- gorm/postgres: Features específicos
- gormwire: Inyección de dependencias
- gorm-logging: Logging avanzado
```

---

## 6.2 - Setup y Configuración

### 6.2.1 - Instalación

```bash
# Instalación básica
go get -u gorm.io/gorm

# Con driver PostgreSQL (recomendado)
go get -u gorm.io/driver/postgres

# Con driver MySQL
go get -u gorm.io/driver/mysql

# Con driver SQLite
go get -u gorm.io/driver/sqlite

# Con driver SQL Server
go get -u gorm.io/driver/sqlserver
```

**go.mod mínimo:**
```
module myapp
go 1.21
require (
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.7
    gorm.io/datatypes v1.2.0
)
```

### 6.2.2 - Conexión PostgreSQL

**Conexión básica:**
```go
package main

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    dsn := "host=localhost user=postgres password=secret " +
           "dbname=myapp port=5432 sslmode=disable"
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }
    
    // db está listo para usar
    fmt.Println("Conectado a PostgreSQL")
}
```

**Conexión con variables de entorno:**
```go
package config

import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

**Archivo .env:**
```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=myapp
DB_PORT=5432
```

### 6.2.3 - Conexión MySQL

```go
import "gorm.io/driver/mysql"

func initMySQL() (*gorm.DB, error) {
    dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

**Características MySQL específicas:**
```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// Usar AUTO_INCREMENT personalizado
type User struct {
    ID    uint   `gorm:"primaryKey;autoIncrement:false"`
    Email string
}
```

### 6.2.4 - Conexión SQLite

```go
import "gorm.io/driver/sqlite"

func initSQLite() (*gorm.DB, error) {
    // En memoria (perfecto para testing)
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    
    // O a archivo
    // db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

### 6.2.5 - Connection Pooling

**Configuración de pool:**
```go
package config

import (
    "database/sql"
    "time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
    dsn := "host=localhost user=postgres password=secret " +
           "dbname=myapp port=5432 sslmode=disable"
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Obtener la conexión SQL subyacente
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // Configurar pool
    sqlDB.SetMaxIdleConns(10)        // Conexiones inactivas máximas
    sqlDB.SetMaxOpenConns(100)       // Conexiones abiertas máximas
    sqlDB.SetConnMaxLifetime(time.Hour) // Vida útil máxima
    sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Inactivo máximo
    
    return db, nil
}
```

**Parámetros recomendados por escala:**
```
Desarrollo:
- MaxIdleConns: 5-10
- MaxOpenConns: 20-50

Producción (carga media):
- MaxIdleConns: 10-20
- MaxOpenConns: 50-100

Producción (alto tráfico):
- MaxIdleConns: 20-30
- MaxOpenConns: 100-300
```

### 6.2.6 - Logging Configuration

**Sin logging:**
```go
db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Silent),
})
```

**Logging detallado:**
```go
import (
    "gorm.io/gorm/logger"
    "log"
    "os"
    "time"
)

func InitDB() (*gorm.DB, error) {
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold: time.Second,    // Queries lentas
            LogLevel:      logger.Info,    // Info, Warn, Error
            Colorful:      true,           // SQL coloreado
        },
    )
    
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    
    return db, nil
}
```

**Logger personalizado:**
```go
type CustomLogger struct{}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
    return l
}

func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
    fmt.Printf("[INFO] %s %v\n", msg, data)
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
    fmt.Printf("[WARN] %s %v\n", msg, data)
}

func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
    fmt.Printf("[ERROR] %s %v\n", msg, data)
}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, 
    fc func() (string, int64), err error) {
    // Custom tracing
}
```

### 6.2.7 - Migration Setup

**Estructura recomendada:**
```
project/
├── main.go
├── config/
│   └── database.go
├── migrations/
│   ├── 001_init_users.go
│   ├── 002_add_posts.go
│   └── migrate.go
└── models/
    ├── user.go
    └── post.go
```

**migrations/migrate.go:**
```go
package migrations

import (
    "fmt"
    "gorm.io/gorm"
    "myapp/models"
)

func Run(db *gorm.DB) error {
    migrations := []func(*gorm.DB) error{
        MigrateUsers,
        MigratePosts,
        MigrateComments,
    }
    
    for _, migration := range migrations {
        if err := migration(db); err != nil {
            return fmt.Errorf("migration failed: %w", err)
        }
    }
    
    return nil
}

func MigrateUsers(db *gorm.DB) error {
    return db.AutoMigrate(&models.User{})
}

func MigratePosts(db *gorm.DB) error {
    return db.AutoMigrate(&models.Post{})
}

func MigrateComments(db *gorm.DB) error {
    return db.AutoMigrate(&models.Comment{})
}
```

---

## 6.3 - Models y Structs

### 6.3.1 - Definición Básica de Modelo

```go
package models

import "time"

type User struct {
    ID        uint            `gorm:"primaryKey"`
    Email     string          `gorm:"uniqueIndex;not null"`
    Name      string
    Age       int
    Active    bool            `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Especificar tabla personalizada
func (User) TableName() string {
    return "users"
}
```

### 6.3.2 - Struct Tags Completos

**Tags básicos:**
```go
type User struct {
    // primaryKey: Indica clave primaria
    ID uint `gorm:"primaryKey"`
    
    // column: Nombre de columna personalizado
    Email string `gorm:"column:user_email"`
    
    // type: Tipo de dato SQL específico
    Status string `gorm:"column:status;type:VARCHAR(50)"`
    
    // size: Tamaño máximo
    Bio string `gorm:"size:1000"`
    
    // not null: Constraint NOT NULL
    Name string `gorm:"not null"`
    
    // unique: Constraint UNIQUE
    Username string `gorm:"unique"`
    
    // uniqueIndex: Índice UNIQUE
    Phone string `gorm:"uniqueIndex:idx_phone"`
    
    // index: Índice simple
    Department string `gorm:"index"`
    
    // default: Valor por defecto
    Role string `gorm:"default:'user'"`
    
    // autoIncrement: AUTO_INCREMENT
    SerialNumber uint `gorm:"autoIncrement"`
    
    // precision: Para decimales
    Balance float64 `gorm:"type:decimal(10,2)"`
    
    // many2many: Relación many-to-many
    // Roles []Role `gorm:"many2many:user_roles;"`
}
```

### 6.3.3 - Field Options Avanzadas

```go
type Product struct {
    // autoIncrementIncrement: Incremento personalizado
    ID uint `gorm:"autoIncrement:100"`
    
    // serializer: Serialización custom
    Metadata map[string]interface{} `gorm:"serializer:json"`
    
    // scan: Permitir valores NULL
    Description *string `gorm:"type:text"`
    
    // comment: Comentario en BD
    Price float64 `gorm:"type:decimal(10,2);comment:Precio en USD"`
    
    // check: Constraint CHECK
    Age int `gorm:"check:age > 0"`
    
    // foreignKey: Llave foránea personalizada
    UserID uint `gorm:"foreignKey:UserID"`
    
    // references: Referencia personalizada
    User *User `gorm:"foreignKey:UserID;references:ID"`
    
    // constraint: DELETE CASCADE
    // User User `gorm:"constraint:OnDelete:CASCADE"`
    
    // embedded: Incrustar struct
    // Timestamps: Timestamps
}
```

### 6.3.4 - Types Soportados

```go
package models

import (
    "time"
    "database/sql"
    "github.com/lib/pq"
    "gorm.io/datatypes"
)

type SupportedTypes struct {
    // Números enteros
    TinyInt   int8
    SmallInt  int16
    Integer   int32
    BigInt    int64
    
    // Números flotantes
    Float32   float32
    Float64   float64
    Decimal   decimal.Decimal // Exactitud
    
    // Strings
    Char      [10]byte
    Varchar   string
    Text      string
    
    // Booleanos
    Bool      bool
    
    // Fechas y horas
    Date      time.Time
    DateTime  time.Time
    Timestamp time.Time
    
    // JSON
    JSONData  datatypes.JSONType
    JSONB     datatypes.JSONType `gorm:"type:jsonb"`
    
    // Arrays (PostgreSQL)
    IntArray  pq.Int64Array `gorm:"type:integer[]"`
    StrArray  pq.StringArray
    
    // Binario
    Data      []byte
    
    // UUID
    UUID      uuid.UUID `gorm:"type:uuid;primaryKey"`
    
    // Nullable
    OptString sql.NullString
    OptInt    sql.NullInt64
    OptTime   sql.NullTime
    
    // Pointer (permite NULL)
    PtrStr    *string
    PtrInt    *int
    PtrTime   *time.Time
}
```

### 6.3.5 - Primary Keys

**Auto-incremento:**
```go
type User struct {
    ID uint `gorm:"primaryKey;autoIncrement"`
    Email string
}
```

**UUID:**
```go
import "github.com/google/uuid"

type User struct {
    ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Email string
}
```

**Composite Primary Key:**
```go
type OrderItem struct {
    OrderID uint `gorm:"primaryKey"`
    ProductID uint `gorm:"primaryKey"`
    Quantity int
    Price float64
}
```

**Key personalizada:**
```go
type Account struct {
    AccountNumber string `gorm:"primaryKey;type:VARCHAR(20)"`
    Balance float64
}
```

### 6.3.6 - Timestamps Automáticos

**Incorporado automáticamente:**
```go
type User struct {
    ID        uint
    Email     string
    CreatedAt time.Time  // Auto-poblado en create
    UpdatedAt time.Time  // Auto-poblado en create y update
}
```

**Timestamps personalizados:**
```go
type User struct {
    ID        uint
    Email     string
    Created   time.Time `gorm:"autoCreateTime"`
    Updated   time.Time `gorm:"autoUpdateTime"`
    Deleted   *time.Time `gorm:"index"` // Soft delete
}
```

**Desactivar timestamps:**
```go
type User struct {
    ID    uint
    Email string `gorm:"autoCreateTime:false"`
}
```

### 6.3.7 - Associations Introducción

**Has One:**
```go
type User struct {
    ID       uint
    Profile  Profile  // Relación 1-1
}

type Profile struct {
    ID     uint
    Bio    string
    UserID uint
}
```

**Has Many:**
```go
type User struct {
    ID    uint
    Email string
    Posts []Post  // Relación 1-N
}

type Post struct {
    ID     uint
    Title  string
    UserID uint
}
```

**Belongs To:**
```go
type Post struct {
    ID     uint
    Title  string
    User   User  // Relación inversa
    UserID uint
}
```

**Many to Many:**
```go
type User struct {
    ID    uint
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID    uint
    Name  string
    Users []User `gorm:"many2many:user_roles;"`
}
```

---

## 6.4 - CRUD Completo

### 6.4.1 - Create (Insert)

**Insert simple:**
```go
func createUser(db *gorm.DB) {
    user := User{
        Email: "john@example.com",
        Name:  "John Doe",
        Age:   30,
    }
    
    result := db.Create(&user)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    
    fmt.Printf("User ID: %d\n", user.ID) // ID asignado automáticamente
}
```

**Insert con valores específicos:**
```go
user := User{Email: "jane@example.com"}
db.Model(&User{}).Create(&user)
```

**Batch insert:**
```go
users := []User{
    {Email: "user1@example.com", Name: "User 1"},
    {Email: "user2@example.com", Name: "User 2"},
    {Email: "user3@example.com", Name: "User 3"},
}

db.Create(&users)
```

**Insert ignorando campos:**
```go
db.Omit("Age", "Role").Create(&user)
```

**Insert con valores por defecto:**
```go
user := User{Email: "test@example.com"}
db.Create(&user)  // Active = true por defecto
```

**Validación antes de crear:**
```go
func (u *User) BeforeSave(tx *gorm.DB) error {
    if u.Email == "" {
        return errors.New("email requerido")
    }
    if u.Age < 0 {
        return errors.New("edad no válida")
    }
    return nil
}

db.Create(&user)  // Llama BeforeSave automáticamente
```

### 6.4.2 - Read/Find

**Find por ID:**
```go
var user User
db.First(&user, 1)  // ID = 1
// o
db.Find(&user, 1)

if errors.Is(db.Error, gorm.ErrRecordNotFound) {
    fmt.Println("User not found")
}
```

**Find todos:**
```go
var users []User
db.Find(&users)
```

**Find con condiciones:**
```go
var users []User
db.Where("age > ?", 18).Find(&users)

// Múltiples condiciones
db.Where("age > ? AND status = ?", 18, "active").Find(&users)
```

**Find con operadores:**
```go
// Igualdad
db.Where("email = ?", "john@example.com").Find(&users)

// Mayor que
db.Where("age > ?", 18).Find(&users)

// Like
db.Where("name LIKE ?", "%john%").Find(&users)

// IN
db.Where("status IN ?", []string{"active", "pending"}).Find(&users)

// Between
db.Where("age BETWEEN ? AND ?", 18, 65).Find(&users)

// NULL
db.Where("deleted_at IS NULL").Find(&users)

// NOT NULL
db.Where("deleted_at IS NOT NULL").Find(&users)
```

**Find con Maps:**
```go
var results []map[string]interface{}
db.Model(&User{}).Select("id", "email", "name").Find(&results)
```

### 6.4.3 - Update

**Update simple:**
```go
db.Model(&User{}).Where("id = ?", 1).Update("email", "new@example.com")
```

**Update múltiples campos:**
```go
db.Model(&User{}).Where("id = ?", 1).Updates(User{
    Email: "new@example.com",
    Age:   31,
})

// O con map (actualiza campos no nulos)
db.Model(&User{}).Where("id = ?", 1).Updates(map[string]interface{}{
    "email": "new@example.com",
    "age":   31,
})
```

**Update todos:**
```go
db.Model(&User{}).Update("status", "inactive")
```

**Update con estructura:**
```go
user := User{ID: 1, Email: "new@example.com"}
db.Model(&user).Updates(&user)
```

**Update selectivos:**
```go
db.Model(&User{}).Where("id = ?", 1).Select("email", "age").Updates(User{
    Email: "new@example.com",
    Age:   31,
})
```

**Update con expresiones SQL:**
```go
db.Model(&User{}).Update("age", gorm.Expr("age + ?", 1))
db.Model(&Order{}).Update("total", gorm.Expr("price * quantity"))
```

### 6.4.4 - Delete

**Delete simple:**
```go
db.Delete(&User{}, 1)  // ID = 1
```

**Delete con condición:**
```go
db.Where("email = ?", "old@example.com").Delete(&User{})
```

**Hard delete (permanente):**
```go
db.Unscoped().Delete(&User{}, 1)
```

**Delete todos:**
```go
db.Delete(&User{})  // Borra todos los registros
```

**Soft delete (por defecto con DeletedAt):**
```go
type User struct {
    ID        uint
    Email     string
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

db.Delete(&user)  // Marca como eliminado
// SELECT * FROM users WHERE deleted_at IS NULL  (automáticamente filtrado)
```

**Recuperar soft deleted:**
```go
db.Unscoped().Where("id = ?", 1).Find(&user)
```

**Purgar soft deleted:**
```go
db.Unscoped().Delete(&User{})  // Borra permanentemente
```

### 6.4.5 - Batch Operations

**Create batch:**
```go
users := []User{
    {Email: "user1@example.com"},
    {Email: "user2@example.com"},
}
result := db.CreateInBatches(users, 100)  // 100 por lote
```

**Update batch:**
```go
db.Model(&User{}).Where("status = ?", "pending").
    Update("status", "active")
```

**Delete batch:**
```go
db.Where("age > ?", 70).Delete(&User{})
```

**Batch processing:**
```go
var users []User
batchSize := 100

db.FindInBatches(batchSize, func(tx *gorm.DB, batch int) error {
    for _, user := range users {
        fmt.Printf("Processing user %d\n", user.ID)
    }
    return nil
}).Find(&users)
```

---

## 6.5 - Querying Avanzado

### 6.5.1 - Where Clauses Complejas

**AND/OR:**
```go
// AND (por defecto)
db.Where("age > ?", 18).Where("status = ?", "active").Find(&users)

// OR
db.Where("status = ? OR status = ?", "active", "pending").Find(&users)

// Combinado
db.Where("(age > ? AND status = ?) OR role = ?", 18, "active", "admin").Find(&users)
```

**IN con slices:**
```go
statuses := []string{"active", "pending"}
db.Where("status IN ?", statuses).Find(&users)

ages := []int{25, 30, 35}
db.Where("age IN ?", ages).Find(&users)
```

**NOT IN:**
```go
db.Where("status NOT IN ?", []string{"banned", "deleted"}).Find(&users)
```

**LIKE:**
```go
db.Where("email LIKE ?", "%@gmail.com").Find(&users)
db.Where("name LIKE ?", "%john%").Find(&users)
```

**NOT LIKE:**
```go
db.Where("email NOT LIKE ?", "%@spam.com").Find(&users)
```

**Condiciones dinámicas:**
```go
query := db.Model(&User{})

if email != "" {
    query = query.Where("email LIKE ?", "%"+email+"%")
}

if minAge > 0 {
    query = query.Where("age >= ?", minAge)
}

if status != "" {
    query = query.Where("status = ?", status)
}

query.Find(&users)
```

### 6.5.2 - Conditions Builder

**Reutilizar conditions:**
```go
isActive := db.Where("deleted_at IS NULL").Where("status = ?", "active")

var activeUsers []User
isActive.Find(&activeUsers)

var activeCount int64
isActive.Model(&User{}).Count(&activeCount)
```

**Conditions complejas:**
```go
cond := db.Where("age > ?", 18).
    Where("status IN ?", []string{"active", "pending"}).
    Where("(role = ? OR is_admin = ?)", "manager", true)

cond.Find(&users)
```

**Condicionales builder pattern:**
```go
type UserFilter struct {
    MinAge   int
    MaxAge   int
    Email    string
    Status   string
}

func (f UserFilter) Apply(db *gorm.DB) *gorm.DB {
    if f.MinAge > 0 {
        db = db.Where("age >= ?", f.MinAge)
    }
    if f.MaxAge > 0 {
        db = db.Where("age <= ?", f.MaxAge)
    }
    if f.Email != "" {
        db = db.Where("email LIKE ?", "%"+f.Email+"%")
    }
    if f.Status != "" {
        db = db.Where("status = ?", f.Status)
    }
    return db
}

filter := UserFilter{MinAge: 18, Status: "active"}
db.Scopes(filter.Apply).Find(&users)
```

### 6.5.3 - Select Específicos

**Seleccionar columnas específicas:**
```go
var users []User
db.Select("id", "email", "name").Find(&users)
```

**Select dinámico:**
```go
columns := []string{"id", "email", "name"}
db.Select(columns).Find(&users)
```

**Excluir columnas:**
```go
db.Omit("password", "api_key").Find(&users)
```

**Select con alias:**
```go
var results []map[string]interface{}
db.Model(&User{}).Select("id as user_id", "email as user_email").Find(&results)
```

**Select con expresiones:**
```go
db.Model(&User{}).Select("id", "email", "age + 1 as age_next_year").Find(&users)
```

### 6.5.4 - Ordering & Pagination

**Ordenar:**
```go
// Ascendente
db.Order("age ASC").Find(&users)
db.Order("age").Find(&users)  // ASC por defecto

// Descendente
db.Order("age DESC").Find(&users)

// Múltiples campos
db.Order("age DESC, email ASC").Find(&users)

// Raw SQL
db.Order(clause.OrderByColumn{
    Column: clause.Column{Name: "age"},
    Desc:   true,
}).Find(&users)
```

**Pagination:**
```go
// Offset/Limit
page := 2
pageSize := 20
offset := (page - 1) * pageSize

var users []User
db.Offset(offset).Limit(pageSize).Find(&users)

// Con orden
db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)
```

**Pagination helper:**
```go
func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if page <= 0 {
            page = 1
        }
        if pageSize <= 0 {
            pageSize = 10
        }
        
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

// Uso
db.Scopes(Paginate(1, 20)).Order("created_at DESC").Find(&users)
```

### 6.5.5 - Distinct & Group By

**Distinct:**
```go
var statuses []string
db.Model(&User{}).Distinct("status").Pluck("status", &statuses)
// Result: ["active", "inactive", "pending"]
```

**Group By:**
```go
var results []map[string]interface{}
db.Model(&User{}).
    Select("status", "COUNT(*) as count").
    Group("status").
    Find(&results)
```

**Group By con Having:**
```go
db.Model(&User{}).
    Select("status", "COUNT(*) as count").
    Group("status").
    Having("COUNT(*) > ?", 5).
    Find(&results)
```

**Agregaciones:**
```go
var count int64
db.Model(&User{}).Count(&count)

var sum int64
db.Model(&Order{}).Select("SUM(total)").Row().Scan(&sum)

var avg float64
db.Model(&User{}).Select("AVG(age)").Row().Scan(&avg)

var max int
db.Model(&User{}).Select("MAX(age)").Row().Scan(&max)
```

**Pluck:**
```go
var emails []string
db.Model(&User{}).Pluck("email", &emails)

// Con condición
db.Where("age > ?", 18).Pluck("email", &emails)
```

---

## 6.6 - Associations & Relationships

### 6.6.1 - Has One

```go
type User struct {
    ID      uint
    Email   string
    Profile Profile  // Relación has one
}

type Profile struct {
    ID     uint
    Bio    string
    UserID uint  // Foreign key
}

// Crear
user := User{Email: "john@example.com"}
user.Profile = Profile{Bio: "Developer"}
db.Create(&user)

// Leer
var user User
db.First(&user, 1)
db.Model(&user).Association("Profile").Find(&user.Profile)

// Update
db.Model(&user).Association("Profile").Replace(&Profile{Bio: "New Bio"})

// Delete (la relación)
db.Model(&user).Association("Profile").Delete()
```

**Configuración personalizada:**
```go
type User struct {
    ID      uint
    Profile Profile `gorm:"foreignKey:UserID;references:ID"`
}

type Profile struct {
    ID     uint
    UserID uint
    Bio    string
}
```

### 6.6.2 - Has Many

```go
type User struct {
    ID    uint
    Email string
    Posts []Post  // Relación has many
}

type Post struct {
    ID     uint
    Title  string
    UserID uint
}

// Crear con relación
user := User{Email: "john@example.com"}
user.Posts = []Post{
    {Title: "Post 1"},
    {Title: "Post 2"},
}
db.Create(&user)

// Leer
var user User
db.First(&user, 1)
db.Model(&user).Association("Posts").Find(&user.Posts)

// Append (agregar más posts)
db.Model(&user).Association("Posts").Append([]Post{
    {Title: "Post 3"},
})

// Replace (reemplazar todos)
db.Model(&user).Association("Posts").Replace([]Post{
    {Title: "New Post 1"},
    {Title: "New Post 2"},
})

// Delete (borrar asociación)
db.Model(&user).Association("Posts").Delete()
```

### 6.6.3 - Belongs To

```go
type Post struct {
    ID     uint
    Title  string
    UserID uint   // Foreign key
    User   User   // Relación belongs to
}

type User struct {
    ID    uint
    Email string
}

// Crear
post := Post{
    Title: "My Post",
    User: User{Email: "john@example.com"},
}
db.Create(&post)

// Leer
var post Post
db.First(&post, 1)
db.Model(&post).Association("User").Find(&post.User)

// Update
db.Model(&post).Update("user_id", 2)
```

### 6.6.4 - Many to Many

**Con tabla join automática:**
```go
type User struct {
    ID    uint
    Email string
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID    uint
    Name  string
    Users []User `gorm:"many2many:user_roles;"`
}

// Crear
user := User{Email: "john@example.com"}
user.Roles = []Role{
    {Name: "admin"},
    {Name: "user"},
}
db.Create(&user)

// Leer
var user User
db.First(&user, 1)
db.Model(&user).Association("Roles").Find(&user.Roles)

// Append
db.Model(&user).Association("Roles").Append(&Role{Name: "moderator"})

// Replace
db.Model(&user).Association("Roles").Replace([]Role{
    {Name: "user"},
})

// Delete
db.Model(&user).Association("Roles").Delete(&Role{Name: "admin"})
```

**Con tabla join personalizada:**
```go
type User struct {
    ID    uint
    Email string
    Roles []Role `gorm:"many2many:user_roles;foreignKey:UserID;references:RoleID;"`
}

type Role struct {
    ID    uint
    Name  string
}

type UserRole struct {
    UserID uint   `gorm:"primaryKey"`
    RoleID uint   `gorm:"primaryKey"`
    GrantedAt time.Time
}

// Crear con datos adicionales
db.Create(&UserRole{UserID: 1, RoleID: 1, GrantedAt: time.Now()})

// Consultar
var userRoles []UserRole
db.Where("user_id = ?", 1).Find(&userRoles)
```

### 6.5.5 - Preloading & Eager Loading

**Preload simple:**
```go
var user User
db.Preload("Posts").First(&user, 1)
// Ejecuta 2 queries: 1 para user, 1 para posts
```

**Preload múltiple:**
```go
db.Preload("Posts").
   Preload("Profile").
   Preload("Roles").
   First(&user, 1)
```

**Preload nested:**
```go
db.Preload("Posts.Comments").First(&user, 1)
```

**Preload con condiciones:**
```go
db.Preload("Posts", "status = ?", "published").
   First(&user, 1)
```

**Preload con orden:**
```go
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC").Limit(10)
}).First(&user, 1)
```

**Preload con select:**
```go
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
    return db.Select("id", "title", "user_id")
}).First(&user, 1)
```

**Select con Associations:**
```go
db.Select("id", "email").
   Preload("Posts", func(db *gorm.DB) *gorm.DB {
       return db.Select("id", "title")
   }).
   First(&user, 1)
```

### 6.6.6 - Lazy Loading

```go
var user User
db.First(&user, 1)

// Lazy load después
var posts []Post
db.Model(&user).Association("Posts").Find(&posts)
```

---

## 6.7 - Transactions & Locking

### 6.7.1 - Transaction Básica

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Todas las operaciones dentro de esta función son una transacción
    
    user := User{Email: "john@example.com", Age: 30}
    if err := tx.Create(&user).Error; err != nil {
        return err  // Rollback automático
    }
    
    post := Post{Title: "My Post", UserID: user.ID}
    if err := tx.Create(&post).Error; err != nil {
        return err  // Rollback automático
    }
    
    return nil  // Commit automático
})

if err != nil {
    log.Fatal(err)
}
```

**Transacción manual:**
```go
tx := db.BeginTx(context.Background(), &sql.TxOptions{
    Isolation: sql.LevelReadCommitted,
})

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&post).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

### 6.7.2 - Transacciones Anidadas

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Transacción externa
    
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    
    // Transacción anidada
    err := tx.Transaction(func(nestedTx *gorm.DB) error {
        if err := nestedTx.Create(&post).Error; err != nil {
            return err  // Rollback solo transacción anidada
        }
        return nil
    })
    
    if err != nil {
        return err  // Rollback toda la transacción
    }
    
    return nil
})
```

### 6.7.3 - Savepoints

```go
tx := db.BeginTx(context.Background(), nil)

// Realizar operación
tx.Create(&user)

// Crear savepoint
txSave := tx.SavePoint("sp1")
tx.Create(&post)

if err := tx.Create(&comment).Error; err != nil {
    // Rollback al savepoint
    tx.RollbackTo("sp1")
    
    // Post aún existe, comment no fue creado
}

tx.Commit()
```

**Savepoint avanzado:**
```go
func transferMoney(db *gorm.DB, fromID, toID uint, amount float64) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Deducir de cuenta origen
        if err := tx.Model(&Account{}).
            Where("id = ?", fromID).
            Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return err
        }
        
        // Savepoint antes de acreditar
        txSp := tx.SavePoint("sp_credit")
        
        // Acreditar a cuenta destino
        if err := tx.Model(&Account{}).
            Where("id = ?", toID).
            Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
            
            txSp.RollbackTo("sp_credit")
            return err
        }
        
        return nil
    })
}
```

### 6.7.4 - Row-Level Locking

**SELECT FOR UPDATE:**
```go
var user User
db.Clauses(clause.Locking{Strength: "UPDATE"}).
    First(&user, 1)
// SELECT * FROM users WHERE id = 1 FOR UPDATE
```

**SELECT FOR UPDATE NOWAIT:**
```go
db.Clauses(clause.Locking{
    Strength: "UPDATE",
    Options:  "NOWAIT",
}).First(&user, 1)
```

**SELECT FOR SHARE:**
```go
db.Clauses(clause.Locking{Strength: "SHARE"}).
    First(&user, 1)
```

**Locking en transacciones:**
```go
err := db.Transaction(func(tx *gorm.DB) error {
    var user User
    
    // Lock la fila
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&user, 1).Error; err != nil {
        return err
    }
    
    // Modificar con seguridad
    user.Balance -= 100
    return tx.Save(&user).Error
})
```

### 6.7.5 - Isolation Levels

```go
// READ UNCOMMITTED
tx := db.BeginTx(context.Background(), &sql.TxOptions{
    Isolation: sql.LevelReadUncommitted,
})

// READ COMMITTED (por defecto en PostgreSQL)
tx := db.BeginTx(context.Background(), &sql.TxOptions{
    Isolation: sql.LevelReadCommitted,
})

// REPEATABLE READ (por defecto en MySQL)
tx := db.BeginTx(context.Background(), &sql.TxOptions{
    Isolation: sql.LevelRepeatableRead,
})

// SERIALIZABLE
tx := db.BeginTx(context.Background(), &sql.TxOptions{
    Isolation: sql.LevelSerializable,
})
```

---

## 6.8 - Hooks y Callbacks

### 6.8.1 - Hook Lifecycle

```
CREATE:
    BeforeSave → BeforeCreate → [INSERT SQL] → AfterCreate → AfterSave

READ:
    BeforeFind → [SELECT SQL] → AfterFind

UPDATE:
    BeforeSave → BeforeUpdate → [UPDATE SQL] → AfterUpdate → AfterSave

DELETE:
    BeforeDelete → [DELETE SQL] → AfterDelete
```

### 6.8.2 - Before/After Hooks

**BeforeSave (Create & Update):**
```go
func (u *User) BeforeSave(tx *gorm.DB) error {
    if u.Email == "" {
        return errors.New("email requerido")
    }
    
    // Normalizar email
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    
    return nil
}
```

**BeforeCreate:**
```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    
    return nil
}
```

**AfterCreate:**
```go
func (u *User) AfterCreate(tx *gorm.DB) error {
    // Log creación
    log.Printf("User %d created: %s\n", u.ID, u.Email)
    
    // Enviar email
    // SendWelcomeEmail(u.Email)
    
    return nil
}
```

**BeforeUpdate:**
```go
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // Validar cambios
    if tx.Statement.Changed("email") {
        if !isValidEmail(u.Email) {
            return errors.New("email no válido")
        }
    }
    return nil
}
```

**AfterUpdate:**
```go
func (u *User) AfterUpdate(tx *gorm.DB) error {
    // Log cambios
    for _, field := range tx.Statement.ChangedFields() {
        log.Printf("Field %s changed\n", field)
    }
    return nil
}
```

**BeforeDelete:**
```go
func (u *User) BeforeDelete(tx *gorm.DB) error {
    // Validar antes de borrar
    var postCount int64
    tx.Model(&Post{}).Where("user_id = ?", u.ID).Count(&postCount)
    
    if postCount > 0 {
        return errors.New("usuario tiene posts, no se puede borrar")
    }
    
    return nil
}
```

**AfterDelete:**
```go
func (u *User) AfterDelete(tx *gorm.DB) error {
    // Limpiar datos asociados
    return tx.Where("user_id = ?", u.ID).Delete(&Post{}).Error
}
```

**BeforeFind:**
```go
func (u *User) BeforeFind(tx *gorm.DB) error {
    // Aplicar automáticamente filtros
    if !tx.Statement.Changed("admin") {
        tx = tx.Where("deleted_at IS NULL")
    }
    return nil
}
```

**AfterFind:**
```go
func (u *User) AfterFind(tx *gorm.DB) error {
    // Desencriptar datos sensibles
    u.Password = ""  // Nunca devolver password
    return nil
}
```

### 6.8.3 - Hook Execution Order

**Orden de ejecución (Create):**
```
1. BeforeSave
2. BeforeCreate
3. SQL INSERT
4. AfterCreate
5. AfterSave
```

**Orden en transacciones:**
```go
db.Transaction(func(tx *gorm.DB) error {
    // Hooks ejecutados dentro de la transacción
    return tx.Create(&user).Error
})
```

### 6.8.4 - Validación con Hooks

```go
type User struct {
    ID       uint
    Email    string
    Age      int
    Password string
}

func (u *User) BeforeSave(tx *gorm.DB) error {
    // Validar email
    if !isValidEmail(u.Email) {
        return errors.New("email inválido")
    }
    
    // Validar edad
    if u.Age < 18 || u.Age > 120 {
        return errors.New("edad fuera de rango")
    }
    
    // Validar password longitud
    if len(u.Password) < 6 {
        return errors.New("password muy corto")
    }
    
    return nil
}

func isValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(pattern)
    return re.MatchString(email)
}
```

### 6.8.5 - Logging con Hooks

```go
type Audit struct {
    ID        uint
    Action    string    // CREATE, UPDATE, DELETE
    TableName string
    RecordID  uint
    Changes   string    // JSON
    CreatedAt time.Time
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    changes, _ := json.Marshal(u)
    
    audit := Audit{
        Action:    "CREATE",
        TableName: "users",
        RecordID:  u.ID,
        Changes:   string(changes),
    }
    
    return tx.Create(&audit).Error
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
    changeMap := make(map[string]interface{})
    
    for _, field := range tx.Statement.ChangedFields() {
        changeMap[field] = tx.Statement.Changes[field]
    }
    
    changes, _ := json.Marshal(changeMap)
    
    audit := Audit{
        Action:    "UPDATE",
        TableName: "users",
        RecordID:  u.ID,
        Changes:   string(changes),
    }
    
    return tx.Create(&audit).Error
}
```

---

## 6.9 - Migrations

### 6.9.1 - AutoMigrate

**Uso básico:**
```go
type User struct {
    ID    uint
    Email string `gorm:"unique"`
    Name  string
}

db.AutoMigrate(&User{})
```

**AutoMigrate múltiples modelos:**
```go
db.AutoMigrate(
    &User{},
    &Post{},
    &Comment{},
    &Tag{},
)
```

**AutoMigrate con opciones:**
```go
type User struct {
    ID    uint
    Email string
}

// Solo crear tabla si no existe
db.AutoMigrate(&User{})

// Con tabla existente, agregar nuevas columnas
db.AutoMigrate(&User{})

// Modificar tipos de columnas (cuidado)
db.Migrator().AlterColumn("email", "VARCHAR(255) NOT NULL")
```

### 6.9.2 - CreateTable, DropTable

**CreateTable:**
```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"unique;not null"`
    Name      string
    Age       int
    CreatedAt time.Time
    UpdatedAt time.Time
}

db.Migrator().CreateTable(&User{})

// Con tabla específica
db.Migrator().CreateTable(&User{}, "users_backup")

// Crear tabla sin índices
db.Migrator().CreateTable(&User{})
```

**DropTable:**
```go
db.Migrator().DropTable(&User{})

// Drop tabla específica
db.Migrator().DropTable("users_backup")

// Drop si existe
db.Migrator().DropTableIfExists(&User{})
```

**HasTable:**
```go
if db.Migrator().HasTable(&User{}) {
    fmt.Println("Tabla users existe")
}

if db.Migrator().HasTable("users") {
    fmt.Println("Tabla users existe")
}
```

### 6.9.3 - AlterTable Operations

**AddColumn:**
```go
type User struct {
    ID    uint
    Email string
    Phone string  // Nueva columna
}

db.Migrator().AddColumn(&User{}, "phone")

// Con tipo específico
db.Migrator().AddColumn(&User{}, "phone", "VARCHAR(20)")
```

**DropColumn:**
```go
db.Migrator().DropColumn(&User{}, "phone")
db.Migrator().DropColumn(&User{}, "phone", "email")  // Múltiples
```

**AlterColumn:**
```go
// Cambiar tipo de columna
db.Migrator().AlterColumn(&User{}, "email", "VARCHAR(255)")

// Cambiar nullable
db.Migrator().AlterColumn(&User{}, "email", "VARCHAR(255) NOT NULL")

// Cambiar default
db.Migrator().AlterColumn(&User{}, "status", "VARCHAR(20) DEFAULT 'active'")
```

**RenameColumn:**
```go
db.Migrator().RenameColumn(&User{}, "old_name", "new_name")
```

**HasColumn:**
```go
if db.Migrator().HasColumn(&User{}, "phone") {
    fmt.Println("Columna phone existe")
}
```

### 6.9.4 - Indexes & Constraints

**CreateIndex:**
```go
type User struct {
    ID    uint
    Email string
}

db.Migrator().CreateIndex(&User{}, "email")

// Índice compuesto
db.Migrator().CreateIndex(&User{}, "idx_email_status", "email", "status")

// Índice único
db.Migrator().CreateUniqueIndex(&User{}, "email")

// Con nombre personalizado
db.Migrator().CreateIndex(&User{}, &clause.Index{
    Name:   "idx_custom_email",
    Column: []clause.Column{{Name: "email"}},
})
```

**DropIndex:**
```go
db.Migrator().DropIndex(&User{}, "email")
db.Migrator().DropIndex(&User{}, "idx_email_status")
```

**HasIndex:**
```go
if db.Migrator().HasIndex(&User{}, "email") {
    fmt.Println("Índice email existe")
}
```

**AddConstraint:**
```go
// Constraint CHECK
db.Migrator().CreateConstraint(&User{}, "age")

// Foreign key
type Post struct {
    ID     uint
    UserID uint
}

db.Migrator().CreateConstraint(&Post{}, "fk_user")
```

### 6.9.5 - Foreign Keys

**AutoMigrate con FK:**
```go
type User struct {
    ID    uint
    Email string
}

type Post struct {
    ID     uint
    Title  string
    User   User
    UserID uint
}

db.AutoMigrate(&User{}, &Post{})
// Automáticamente crea FK
```

**FK personalizada:**
```go
type Post struct {
    ID     uint
    Title  string
    UserID uint `gorm:"foreignKey:UserID;references:ID"`
    User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
```

**FK con acciones:**
```go
type Post struct {
    ID     uint
    Title  string
    User   User `gorm:"constraint:OnDelete:CASCADE;OnUpdate:CASCADE"`
    UserID uint
}
```

**Opciones:**
```
OnDelete:
  - CASCADE: Borrar posts si se borra usuario
  - SET NULL: Poner NULL si se borra usuario
  - RESTRICT: No permitir borrar usuario con posts
  - NO ACTION: Similar a RESTRICT

OnUpdate:
  - CASCADE: Actualizar FK si cambia ID usuario
  - SET NULL: Poner NULL si cambia ID
```

---

## 6.10 - Advanced Patterns

### 6.10.1 - Soft Delete

**Implementación básica:**
```go
type User struct {
    ID        uint
    Email     string
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Soft delete
db.Delete(&user)  // No borra, solo marca como eliminado

// Query automáticamente excluye soft deleted
var users []User
db.Find(&users)  // Solo activos

// Incluir soft deleted
db.Unscoped().Find(&users)

// Solo soft deleted
db.Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
```

**Restaurar soft deleted:**
```go
db.Model(&User{}).Where("id = ?", 1).Update("deleted_at", nil)

// O con método
db.Model(&user).Update("deleted_at", nil)
```

**Purgar soft deleted:**
```go
// Borrar permanentemente
db.Unscoped().Delete(&user)

// Borrar permanentemente todos los soft deleted
db.Unscoped().Where("deleted_at IS NOT NULL").Delete(&User{})
```

### 6.10.2 - Polymorphism

```go
type Comment struct {
    ID               uint
    Content          string
    CommentableType  string  // "Post" o "Article"
    CommentableID    uint    // ID del objeto comentado
}

type Post struct {
    ID       uint
    Title    string
    Comments []Comment `gorm:"polymorphic:Commentable;"`
}

type Article struct {
    ID       uint
    Title    string
    Comments []Comment `gorm:"polymorphic:Commentable;"`
}

// Crear comentario en post
post := Post{ID: 1}
db.Model(&post).Association("Comments").Append(&Comment{Content: "Great post!"})

// Query comentarios de post
var comments []Comment
db.Where("commentable_type = ? AND commentable_id = ?", "Post", 1).Find(&comments)
```

### 6.10.3 - Custom Types

```go
import "database/sql/driver"
import "encoding/json"

// Custom type
type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
    StatusPending  Status = "pending"
)

// Implementar Valuer
func (s Status) Value() (driver.Value, error) {
    return string(s), nil
}

// Implementar Scanner
func (s *Status) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion failed")
    }
    *s = Status(bytes)
    return nil
}

type User struct {
    ID     uint
    Status Status
}

// Uso
user := User{Status: StatusActive}
db.Create(&user)
```

**JSON custom type:**
```go
type Metadata map[string]interface{}

func (m Metadata) Value() (driver.Value, error) {
    return json.Marshal(m)
}

func (m *Metadata) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion failed")
    }
    return json.Unmarshal(bytes, &m)
}

type User struct {
    ID       uint
    Metadata Metadata `gorm:"type:jsonb"`
}
```

### 6.10.4 - JSON Fields

**JSONB en PostgreSQL:**
```go
import "gorm.io/datatypes"

type User struct {
    ID    uint
    Email string
    // Guardar como JSONB
    Profile datatypes.JSONType `gorm:"type:jsonb;default:'{}'::jsonb"`
}

user := User{
    Email: "john@example.com",
    Profile: datatypes.JSONType(`{"age": 30, "city": "NYC"}`),
}
db.Create(&user)
```

**JSON queries:**
```go
// PostgreSQL JSON operators
db.Where("profile->>'city' = ?", "NYC").Find(&users)

// MySQL JSON_EXTRACT
db.Where("JSON_EXTRACT(profile, '$.city') = ?", "NYC").Find(&users)
```

**Marshaling JSON:**
```go
type UserProfile struct {
    Age   int    `json:"age"`
    City  string `json:"city"`
    Job   string `json:"job"`
}

type User struct {
    ID       uint
    Email    string
    Profile  *UserProfile `gorm:"type:jsonb;serializer:json"`
}

// Usar directly
user := User{
    Email: "john@example.com",
    Profile: &UserProfile{
        Age:  30,
        City: "NYC",
        Job:  "Developer",
    },
}
db.Create(&user)

// Query
var user User
db.First(&user)
fmt.Printf("City: %s\n", user.Profile.City)
```

### 6.10.5 - Full-Text Search

**PostgreSQL FTS:**
```go
type Document struct {
    ID    uint
    Title string
    Body  string
    // Vector de búsqueda
    Vector datatypes.TSVector `gorm:"type:tsvector;index:,type:gin"`
}

// Crear función trigger para actualizar vector
db.Exec(`
    CREATE OR REPLACE FUNCTION documents_search_update() RETURNS trigger AS $$
    BEGIN
        NEW.vector := to_tsvector('english', NEW.title || ' ' || NEW.body);
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER documents_search_trigger
    BEFORE INSERT OR UPDATE ON documents
    FOR EACH ROW EXECUTE FUNCTION documents_search_update();
`)

// Búsqueda
var results []Document
db.Where("vector @@ plainto_tsquery('english', ?)", "golang").Find(&results)
```

**MySQL FULLTEXT:**
```go
type Article struct {
    ID    uint
    Title string `gorm:"index:,type:fulltext"`
    Body  string `gorm:"index:,type:fulltext"`
}

// Búsqueda
var articles []Article
db.Where("MATCH(title, body) AGAINST(? IN BOOLEAN MODE)", "+golang -rust").
    Find(&articles)
```

**SQLite MATCH:**
```go
// MATCH no es ideal en SQLite, usar LIKE
var articles []Article
db.Where("title LIKE ? OR body LIKE ?", "%golang%", "%golang%").
    Find(&articles)
```

---

## 6.11 - Performance & Best Practices

### 6.11.1 - Query Optimization

**Anti-pattern: N+1 Problem**
```go
// ❌ MALO: N+1 queries
var users []User
db.Find(&users)  // Query 1

for _, user := range users {
    var posts []Post
    db.Where("user_id = ?", user.ID).Find(&posts)  // N queries adicionales
}
```

**Solución con Preload:**
```go
// ✅ BUENO: Preload (2 queries total)
var users []User
db.Preload("Posts").Find(&users)

for _, user := range users {
    fmt.Printf("User %s has %d posts\n", user.Email, len(user.Posts))
}
```

**Solución con Joins:**
```go
// ✅ ALTERNATIVA: Join (1 query)
var results []struct {
    UserEmail string
    PostTitle string
}{} 
db.Model(&User{}).
    Select("users.email as user_email", "posts.title as post_title").
    Joins("LEFT JOIN posts ON users.id = posts.user_id").
    Find(&results)
```

### 6.11.2 - Batch Processing

**Batch insert:**
```go
// ❌ LENTO: Inserts individuales
for i := 0; i < 10000; i++ {
    db.Create(&User{Email: fmt.Sprintf("user%d@example.com", i)})
}

// ✅ RÁPIDO: Batch insert
users := make([]User, 10000)
for i := 0; i < 10000; i++ {
    users[i] = User{Email: fmt.Sprintf("user%d@example.com", i)}
}
db.CreateInBatches(users, 500)  // 500 registros por batch
```

**Batch update:**
```go
// ❌ LENTO
for _, id := range userIDs {
    db.Model(&User{}).Where("id = ?", id).Update("status", "inactive")
}

// ✅ RÁPIDO
db.Model(&User{}).Where("id IN ?", userIDs).Update("status", "inactive")
```

**Batch delete:**
```go
// ❌ LENTO
for _, id := range userIDs {
    db.Delete(&User{}, id)
}

// ✅ RÁPIDO
db.Where("id IN ?", userIDs).Delete(&User{})
```

### 6.11.3 - Indexes Strategy

**Índices esenciales:**
```go
type User struct {
    // Primary key (automático)
    ID uint `gorm:"primaryKey"`
    
    // Búsquedas frecuentes
    Email string `gorm:"uniqueIndex"`
    Status string `gorm:"index"`
    
    // Foreign keys
    DepartmentID uint `gorm:"index"`
    
    // Timestamps para queries de rango
    CreatedAt time.Time `gorm:"index"`
}
```

**Índices compuestos:**
```go
type Order struct {
    ID        uint
    UserID    uint
    Status    string
    CreatedAt time.Time
    
    // Índice para queries como (user_id, status)
    // `gorm:"index:idx_user_status,priority:1"`  // UserID
}

type OrderIdx struct {
    ID        uint
    UserID    uint `gorm:"index:idx_user_status,priority:1"`
    Status    string `gorm:"index:idx_user_status,priority:2"`
    CreatedAt time.Time
}
```

**Evitar índices innecesarios:**
```go
// ❌ Demasiados índices
type User struct {
    ID        uint   `gorm:"index"`
    Email     string `gorm:"index"`
    Name      string `gorm:"index"`
    Age       int    `gorm:"index"`
    Status    string `gorm:"index"`
    CreatedAt time.Time `gorm:"index"`
}

// ✅ Índices estratégicos
type User struct {
    ID        uint
    Email     string `gorm:"uniqueIndex"`  // Buscas frecuentes
    Name      string  // Rara vez como único criterio
    Age       int     // No indexar
    Status    string `gorm:"index"`  // Filtro común
    CreatedAt time.Time `gorm:"index"`  // Range queries
}
```

### 6.11.4 - Select Específicos

```go
// ❌ MALO: Traer todas las columnas
var users []User
db.Find(&users)  // SELECT * ...

// ✅ BUENO: Solo columnas necesarias
var emails []string
db.Model(&User{}).Where("status = ?", "active").
    Pluck("email", &emails)

// ✅ BUENO: Select específico
var users []struct {
    ID    uint
    Email string
}{} 
db.Model(&User{}).Select("id", "email").Find(&users)
```

### 6.11.5 - Connection Pooling Optimization

```go
func initOptimizedDB() (*gorm.DB, error) {
    dsn := "host=localhost user=postgres password=secret dbname=myapp port=5432 sslmode=disable"
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    
    sqlDB, _ := db.DB()
    
    // Configuración para alto rendimiento
    sqlDB.SetMaxIdleConns(25)      // Mantener conexiones disponibles
    sqlDB.SetMaxOpenConns(50)      // Límite máximo
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    sqlDB.SetConnMaxIdleTime(2 * time.Minute)
    
    return db, nil
}
```

### 6.11.6 - Scopes para Reutilizar Queries

```go
// Definir scopes
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("deleted_at IS NULL").Where("status = ?", "active")
}

func RecentUsers(days int) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("created_at > ?", time.Now().AddDate(0, 0, -days))
    }
}

func AdminUsers(db *gorm.DB) *gorm.DB {
    return db.Where("role = ?", "admin")
}

// Usar scopes
var activeAdmins []User
db.Scopes(ActiveUsers, AdminUsers).Find(&activeAdmins)

var recentUsers []User
db.Scopes(RecentUsers(7)).Find(&recentUsers)
```

### 6.11.7 - Pagination Eficiente

```go
type PaginationResult struct {
    Total   int64
    Page    int
    PageSize int
    Data    interface{}
}

func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if page <= 0 {
            page = 1
        }
        if pageSize <= 0 {
            pageSize = 10
        }
        
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

func GetUsers(db *gorm.DB, page, pageSize int) (PaginationResult, error) {
    var users []User
    var total int64
    
    // Contar registros totales
    db.Model(&User{}).Count(&total)
    
    // Obtener página
    if err := db.Scopes(Paginate(page, pageSize)).Find(&users).Error; err != nil {
        return PaginationResult{}, err
    }
    
    return PaginationResult{
        Total:    total,
        Page:     page,
        PageSize: pageSize,
        Data:     users,
    }, nil
}
```

### 6.11.8 - Production Patterns

**Health Check:**
```go
func HealthCheck(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    return sqlDB.Ping()
}
```

**Graceful Shutdown:**
```go
func shutdown(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    return sqlDB.Close()
}

// Uso en main
defer shutdown(db)
```

**Retry Logic:**
```go
func WithRetry(fn func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        if err := fn(); err == nil {
            return nil
        } else {
            lastErr = err
            time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
        }
    }
    
    return lastErr
}

// Uso
err := WithRetry(func() error {
    return db.Create(&user).Error
}, 3)
```

**Context Timeout:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var users []User
db.WithContext(ctx).Find(&users)
```

### 6.11.9 - Troubleshooting

**Enable Query Logging:**
```go
import "gorm.io/gorm/logger"

db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Ver queries lentas
// Queries > 200ms aparecerán con [SLOW SQL]
```

**Analyze Query Plan:**
```go
// PostgreSQL
var plan []map[string]interface{}
db.Raw("EXPLAIN ANALYZE SELECT * FROM users WHERE age > ?", 18).
    Scan(&plan)

// MySQL
var plan []map[string]interface{}
db.Raw("EXPLAIN SELECT * FROM users WHERE age > ?", 18).
    Scan(&plan)
```

**Check Indexes:**
```go
// PostgreSQL
var indexes []map[string]interface{}
db.Raw("SELECT * FROM pg_indexes WHERE tablename = ?", "users").
    Scan(&indexes)

// MySQL
var indexes []map[string]interface{}
db.Raw("SHOW INDEX FROM users").Scan(&indexes)
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: Simple CRUD (Todos)

**Modelo:**
```go
package models

import "time"

type Todo struct {
    ID        uint
    Title     string
    Completed bool
    DueDate   *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Operaciones requeridas:**
```go
// 1. Create: Agregar nuevo todo
// 2. Read: Obtener todos, por ID, pendientes, completados
// 3. Update: Marcar como completado
// 4. Delete: Eliminar todo

// Bonus:
// - Soft delete
// - List paginado
// - Orden por fecha de creación
```

**Solución:**
```go
package main

import (
    "fmt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "time"
)

type Todo struct {
    ID        uint
    Title     string
    Completed bool
    DueDate   *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}

func main() {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&Todo{})
    
    // CREATE
    todo1 := Todo{Title: "Learn GORM", Completed: false}
    db.Create(&todo1)
    fmt.Printf("Created: %+v\n", todo1)
    
    // READ
    var todos []Todo
    db.Find(&todos)
    for _, t := range todos {
        fmt.Printf("Todo: %s (Completed: %v)\n", t.Title, t.Completed)
    }
    
    // UPDATE
    db.Model(&todo1).Update("completed", true)
    
    // DELETE
    db.Delete(&todo1)
    
    // Verify deleted
    db.Find(&todos)
    fmt.Printf("After delete: %d todos\n", len(todos))
}
```

### Ejercicio 2: Relationships (Blog + Posts + Comments)

**Modelos:**
```go
type User struct {
    ID    uint
    Email string
    Posts []Post
}

type Post struct {
    ID       uint
    Title    string
    Content  string
    UserID   uint
    Comments []Comment
}

type Comment struct {
    ID      uint
    Body    string
    PostID  uint
    UserID  uint
}
```

**Requerimientos:**
```
1. Crear usuario con múltiples posts
2. Crear posts con comentarios
3. Leer usuario con todos sus posts y comentarios
4. Actualizar post
5. Eliminar post (cascade elimina comentarios)
6. Preload para evitar N+1
```

**Solución:**
```go
func main() {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&User{}, &Post{}, &Comment{})
    
    // Create user with posts
    user := User{
        Email: "john@example.com",
        Posts: []Post{
            {
                Title:   "Post 1",
                Content: "Content 1",
                Comments: []Comment{
                    {Body: "Great!"},
                    {Body: "Thanks!"},
                },
            },
            {
                Title:   "Post 2",
                Content: "Content 2",
            },
        },
    }
    db.Create(&user)
    
    // Read with preload
    var loadedUser User
    db.Preload("Posts.Comments").First(&loadedUser, user.ID)
    
    fmt.Printf("User: %s\n", loadedUser.Email)
    for _, post := range loadedUser.Posts {
        fmt.Printf("  Post: %s (%d comments)\n", post.Title, len(post.Comments))
    }
}
```

### Ejercicio 3: Transactions & Hooks

**Requerimientos:**
```
1. Transfer money between accounts (transaction)
2. Log all changes (hooks)
3. Validate balance antes de transfer
4. Rollback si hay error
5. Usar BeforeUpdate para auditoría
```

**Solución:**
```go
type Account struct {
    ID      uint
    Email   string
    Balance float64
}

func (a *Account) BeforeUpdate(tx *gorm.DB) error {
    if a.Balance < 0 {
        return errors.New("balance cannot be negative")
    }
    return nil
}

func transferMoney(db *gorm.DB, fromID, toID uint, amount float64) error {
    return db.Transaction(func(tx *gorm.DB) error {
        var from, to Account
        
        tx.First(&from, fromID)
        if from.Balance < amount {
            return errors.New("insufficient funds")
        }
        
        tx.First(&to, toID)
        
        from.Balance -= amount
        to.Balance += amount
        
        if err := tx.Save(&from).Error; err != nil {
            return err
        }
        if err := tx.Save(&to).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### Ejercicio 4: Migrations & Schema

**Requerimientos:**
```
1. Create tabla users con índices
2. Add columna phone a users
3. Create tabla posts
4. Add foreign key posts->users
5. Create índice compuesto (user_id, status)
```

**Solución:**
```go
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Email    string `gorm:"uniqueIndex"`
    Name     string
    Phone    string `gorm:"index"`
}

type Post struct {
    ID     uint
    Title  string
    Status string `gorm:"index:idx_user_status,priority:2"`
    User   User   `gorm:"constraint:OnDelete:CASCADE"`
    UserID uint   `gorm:"index:idx_user_status,priority:1"`
}

func main() {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    
    // Auto migrate
    db.AutoMigrate(&User{}, &Post{})
    
    // Verify
    fmt.Println(db.Migrator().HasTable(&User{}))
    fmt.Println(db.Migrator().HasColumn(&User{}, "phone"))
    fmt.Println(db.Migrator().HasIndex(&Post{}, "idx_user_status"))
}
```

### Ejercicio 5: Complete Production App

**Requerimientos:**
```
1. API CRUD completo
2. Pagination
3. Filtering & Searching
4. Transactions
5. Logging & Auditoría
6. Soft delete
7. Preloading
8. Error handling
```

**Solución (Estructura):**
```
project/
├── main.go
├── config/database.go
├── models/
│   ├── user.go
│   ├── post.go
│   └── audit.go
├── handlers/
│   ├── user.go
│   └── post.go
├── migrations/
│   └── migrate.go
└── middleware/
    └── error.go
```

**main.go:**
```go
package main

import (
    "fmt"
    "net/http"
    "gorm.io/gorm"
)

var db *gorm.DB

func main() {
    // Initialize DB
    var err error
    db, err = InitDB()
    if err != nil {
        panic(err)
    }
    
    // Run migrations
    if err := RunMigrations(db); err != nil {
        panic(err)
    }
    
    // Setup routes
    http.HandleFunc("/users", ListUsers)
    http.HandleFunc("/users/:id", GetUser)
    http.HandleFunc("/users/create", CreateUser)
    
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
    page := r.URL.Query().Get("page")
    email := r.URL.Query().Get("email")
    
    var users []User
    query := db
    
    if email != "" {
        query = query.Where("email LIKE ?", "%"+email+"%")
    }
    
    query.Offset((toInt(page)-1)*10).Limit(10).Find(&users)
    
    // Return JSON
}
```

---

## COMPARATIVAS DETALLADAS

### GORM vs sqlc vs Ent

```
┌──────────────┬─────────────────────┬─────────────────────┬─────────────────────┐
│ Aspecto      │ GORM                │ sqlc                │ Ent                 │
├──────────────┼─────────────────────┼─────────────────────┼─────────────────────┤
│ Setup        │ Fácil (go get)      │ Requiere SQL manual │ Más complejo        │
│ Aprendizaje  │ 1-2 días            │ 2-3 días            │ 1-2 semanas         │
│ Flexibilidad │ Muy alta            │ Media               │ Baja                │
│ Type safety  │ Parcial             │ Total               │ Total               │
│ Performance  │ Bueno               │ Óptimo              │ Bueno               │
│ Comunidad    │ Grande              │ Mediana             │ Pequeña             │
│ Migrations   │ Integradas          │ Manual              │ Integradas          │
│ Hooks        │ Sí                  │ No                  │ Sí                  │
│ Testing      │ Fácil               │ Fácil               │ Más complejo        │
│ Producción   │ Excelente           │ Excelente           │ Bueno               │
└──────────────┴─────────────────────┴─────────────────────┴─────────────────────┘
```

### Benchmarks

```
Operación          GORM        sqlc        Raw SQL     Ent
─────────────────────────────────────────────────────────
Insert (1 row):    2.1ms       1.8ms       1.5ms       2.3ms
Select (1 row):    1.9ms       1.6ms       1.4ms       2.2ms
Update (1 row):    2.4ms       2.1ms       1.8ms       2.6ms
Delete (1 row):    2.2ms       2.0ms       1.7ms       2.4ms
Batch Insert (1k): 850ms       780ms       650ms       950ms
Complex Join:      4.2ms       3.8ms       3.5ms       5.1ms
```

---

## DIAGRAMAS

### Lifecycle CRUD

```
CREATE:
┌────────────────────────────────────────────────────┐
│                  User.Create(&user)                │
└────────────────────────────────────────────────────┘
              ↓
┌────────────────────────────────────────────────────┐
│               BeforeSave Hook                      │
│           - Validación de datos                   │
│           - Normalización                         │
└────────────────────────────────────────────────────┘
              ↓
┌────────────────────────────────────────────────────┐
│               BeforeCreate Hook                    │
│           - Hash de password                      │
│           - Generación de tokens                  │
└────────────────────────────────────────────────────┘
              ↓
┌────────────────────────────────────────────────────┐
│               [INSERT SQL]                         │
│      INSERT INTO users (email, name) VALUES...    │
└────────────────────────────────────────────────────┘
              ↓
┌────────────────────────────────────────────────────┐
│               AfterCreate Hook                     │
│           - Enviar email de bienvenida            │
│           - Logging                               │
└────────────────────────────────────────────────────┘
              ↓
┌────────────────────────────────────────────────────┐
│               AfterSave Hook                       │
│           - Actualizar cache                      │
│           - Notificaciones                        │
└────────────────────────────────────────────────────┘
```

### Association Types

```
Has One (1-1):
┌─────────────┐          ┌──────────────┐
│    User     │ 1    1   │   Profile    │
│  ────────   │◀─────────│  ──────────  │
│ ID          │          │ ID           │
│             │          │ UserID (FK)  │
└─────────────┘          └──────────────┘

Has Many (1-N):
┌─────────────┐          ┌──────────────┐
│    User     │ 1     N  │    Posts     │
│  ────────   │◀─────────│  ──────────  │
│ ID          │          │ ID           │
│             │          │ UserID (FK)  │
└─────────────┘          └──────────────┘

Many to Many (N-N):
┌─────────────┐     ┌────────────┐     ┌──────────┐
│    User     │────▶│ UserRoles  │◀────│  Role    │
│  ────────   │     │ ──────────│     │  ──────  │
│ ID          │     │ UserID    │     │ ID       │
│             │     │ RoleID    │     │          │
└─────────────┘     └────────────┘     └──────────┘
```

### Query Execution Flow

```
db.Where("age > ?", 18).
   Preload("Posts").
   Order("created_at DESC").
   Limit(10).
   Find(&users)
     ↓
┌─────────────────────────────────────┐
│  Build Query String                 │
│  SELECT * FROM users WHERE age > 18 │
│  ORDER BY created_at DESC           │
│  LIMIT 10                           │
└─────────────────────────────────────┘
     ↓
┌─────────────────────────────────────┐
│  Execute Main Query                 │
│  Get 10 users                       │
└─────────────────────────────────────┘
     ↓
┌─────────────────────────────────────┐
│  Preload Associations               │
│  For each user:                     │
│  SELECT * FROM posts WHERE user_id IN (...)
└─────────────────────────────────────┘
     ↓
┌─────────────────────────────────────┐
│  Return Result Set                  │
│  users with posts populated         │
└─────────────────────────────────────┘
```

---

## ANTI-PATTERNS vs BEST PRACTICES

### N+1 Problem

**❌ ANTI-PATTERN:**
```go
var users []User
db.Find(&users)  // Query 1

for _, user := range users {
    var posts []Post
    db.Where("user_id = ?", user.ID).Find(&posts)  // Queries 2-N+1
    // Total: 1 + N queries
}
```

**✅ BEST PRACTICE:**
```go
var users []User
db.Preload("Posts").Find(&users)  // Query 1 (users) + Query 2 (all posts)
// Total: 2 queries
```

### Loading All Columns

**❌ ANTI-PATTERN:**
```go
var users []User
db.Find(&users)  // SELECT * FROM users
// Trae password, api_key, sensible_data...
```

**✅ BEST PRACTICE:**
```go
var users []struct {
    ID    uint
    Email string
    Name  string
}{} 
db.Model(&User{}).Select("id", "email", "name").Find(&users)
// SELECT id, email, name FROM users
```

### Missing Indexes

**❌ ANTI-PATTERN:**
```go
type User struct {
    ID    uint
    Email string
    Age   int
    Phone string
}
// Sin índices - queries lentos
db.Where("email = ?", email).First(&user)  // Full table scan
```

**✅ BEST PRACTICE:**
```go
type User struct {
    ID    uint
    Email string `gorm:"uniqueIndex"`  // Índice único
    Age   int
    Phone string `gorm:"index"`        // Índice simple
}
// Queries rápidos
```

### No Pagination

**❌ ANTI-PATTERN:**
```go
var users []User
db.Find(&users)  // SELECT * FROM users
// Si hay 1 millón de usuarios, todo a memoria
```

**✅ BEST PRACTICE:**
```go
var users []User
page := 1
pageSize := 20
db.Offset((page-1)*pageSize).Limit(pageSize).Find(&users)
// SELECT * FROM users LIMIT 20 OFFSET 0
```

### No Transactions

**❌ ANTI-PATTERN:**
```go
// Transfer money
from.Balance -= 100
db.Save(&from)

to.Balance += 100
db.Save(&to)
// Si db.Save(&to) falla, dinero desaparece
```

**✅ BEST PRACTICE:**
```go
db.Transaction(func(tx *gorm.DB) error {
    from.Balance -= 100
    if err := tx.Save(&from).Error; err != nil {
        return err  // Rollback automático
    }
    
    to.Balance += 100
    if err := tx.Save(&to).Error; err != nil {
        return err  // Rollback automático
    }
    
    return nil  // Commit
})
```

---

## COMPARACIÓN: GORM vs Raw SQL

```go
// OPERACIÓN: Obtener usuarios activos de último mes

// Raw SQL
sql := `
    SELECT u.id, u.email, u.name, u.created_at
    FROM users u
    WHERE u.deleted_at IS NULL 
    AND u.status = 'active'
    AND u.created_at > NOW() - INTERVAL '30 days'
    ORDER BY u.created_at DESC
    LIMIT 20
`
rows, err := db.Raw(sql).Rows()
// Manual scanning, manual error handling

// GORM
var users []User
err := db.
    Where("deleted_at IS NULL").
    Where("status = ?", "active").
    Where("created_at > ?", time.Now().AddDate(0, 0, -30)).
    Order("created_at DESC").
    Limit(20).
    Find(&users).Error
// Automático scanning, automático error handling
```

---

## RECURSOS ADICIONALES

### Documentación oficial
- GitHub: https://github.com/go-gorm/gorm
- Docs: https://gorm.io/docs/

### Plugins útiles
```bash
go get -u gorm.io/datatypes
go get -u gorm.io/hints
go get -u gorm.io/clause_builder
```

### Testing con GORM
```go
import (
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestUserCreate(t *testing.T) {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&User{})
    
    user := User{Email: "test@example.com"}
    result := db.Create(&user)
    
    if result.Error != nil {
        t.Fatalf("failed to create user: %v", result.Error)
    }
    
    if user.ID == 0 {
        t.Fatal("user ID should not be 0")
    }
}
```

---

## CONCLUSIÓN

GORM es un ORM poderoso, flexible y bien mantenido que es la mejor opción para la mayoría de aplicaciones Go. Su curva de aprendizaje suave, comunidad activa y características extensivas lo hacen ideal para:

✅ Startups y MVPs
✅ Aplicaciones empresariales
✅ Prototipos rápidos
✅ Sistemas con relaciones complejas
✅ Negocios que requieren escalabilidad

Sin embargo, para casos muy específicos (raw performance crítica, queries ultra-complejas), considerar alternativas como sqlc o raw SQL.

**Tabla resumen de decisión:**

```
Caso de uso                      → Recomendación
─────────────────────────────────────────────────────
MVP rápido                       → GORM
Startup escalable                → GORM
Enterprise app                   → GORM + sqlc híbrido
Data warehouse                   → Raw SQL
Real-time analytics              → Raw SQL
Proyecto pequeño                 → SQLite + GORM
```


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/61-gorm-orm/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/61-gorm-orm):

```bash
cd examples/61-gorm-orm
go run .
```
