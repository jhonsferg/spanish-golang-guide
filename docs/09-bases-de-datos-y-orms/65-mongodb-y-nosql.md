# Capítulo 65: MongoDB y NoSQL

## 65.1 Introducción a NoSQL

### 65.1.1 ¿Qué es NoSQL?

NoSQL (Not Only SQL) representa un cambio paradigmático en el almacenamiento de datos. A diferencia de las bases de datos relacionales que organizan información en tablas con esquemas rígidos, NoSQL ofrece modelos más flexibles y distribuidos.

**Características de NoSQL:**
- Esquemas flexibles (schemaless)
- Escalabilidad horizontal
- Replicación automática
- Mejor rendimiento en lectura/escritura masiva
- Distribución geográfica

### 65.1.2 Tipos de Bases de Datos NoSQL

```
┌─────────────────────────────────────────────────────┐
│                    NoSQL Landscape                   │
├─────────────────────────────────────────────────────┤
│                                                      │
│  Document:           Key-Value:    Graph:           │
│  • MongoDB           • Redis        • Neo4j          │
│  • Firebase          • Memcached    • ArangoDB       │
│  • CouchDB           • DynamoDB     • TigerGraph     │
│                                                      │
│  Column-Family:      Time-Series:                    │
│  • Cassandra         • InfluxDB                      │
│  • HBase             • Prometheus                    │
│                                                      │
└─────────────────────────────────────────────────────┘
```

### 65.1.3 SQL vs NoSQL Trade-offs

| Aspecto | SQL | NoSQL (Documentos) |
|---------|-----|-------------------|
| **Esquema** | Rígido, predefinido | Flexible, dinámico |
| **Escalabilidad** | Vertical (principalmente) | Horizontal |
| **Transacciones ACID** | Garantizadas | Limitadas (mejorando) |
| **Joins** | Nativos, optimizados | Costosos o via embedding |
| **Consistencia** | Fuerte (ACID) | Eventual (BASE) |
| **Velocidad lectura** | Moderada | Muy rápida (denormalizado) |
| **Volumen datos** | Millones de filas | Billones de documentos |
| **Complejidad queries** | Muy compleja | Moderada |
| **Replicación** | Configuración compleja | Integrada |

**Matriz de decisión:**

```
                ┌─ Muchas relaciones?
                │  ├─ Sí → SQL (PostgreSQL)
                │  └─ No → NoSQL
                │
        ├─ Datos estructurados?
        │  ├─ Sí → SQL
        │  └─ No → NoSQL
        │
        ├─ Escalabilidad crítica?
        │  ├─ Sí → NoSQL
        │  └─ Tal vez SQL
        │
        └─ ACID transacciones?
           ├─ Esencial → SQL
           └─ Opcional → NoSQL
```

### 65.1.4 Cuándo Usar MongoDB

**✅ Casos ideales:**
- Aplicaciones web modernas (contenido, perfiles)
- IoT y datos en tiempo real
- Análisis de Big Data
- Prototipado rápido (esquema evolutivo)
- Almacenamiento de documentos complejos
- Startups con crecimiento impredecible

**❌ Casos NO ideales:**
- Transacciones financieras complejas
- Datos altamente relacionados
- Reportes con joins complicados
- Integridad referencial crítica

### 65.1.5 Otras Bases de Datos NoSQL Relevantes

#### DynamoDB (AWS)
```go
// Totalmente manejada por AWS
// Facturación por throughput (reads/writes)
// Perfect para aplicaciones serverless
// Limitaciones: 400KB por item, 25 atributos máximo

// Ejemplo conceptual
type DynamoProduct struct {
    ID       string    `dynamodbav:"id"`
    Name     string    `dynamodbav:"name"`
    Price    float64   `dynamodbav:"price"`
    TTL      int64     `dynamodbav:"ttl"`  // auto-expire
}
```

#### Firestore (Google Cloud)
```go
// Tiempo real, sincronización automática
// Excelente para aplicaciones móviles
// Pricing: operaciones de lectura/escritura

type FirestoreUser struct {
    Name      string    `firestore:"name"`
    Email     string    `firestore:"email"`
    CreatedAt time.Time `firestore:"created_at,serverTimestamp"`
}
```

#### Cassandra
```go
// Orientada a columnas (no documentos)
// Altamente distribuida, sin punto único de fallo
// Ideal: Series de tiempo, logs masivos

type CassandraEvent struct {
    TimeID    time.Time
    UserID    string
    EventType string
    Metadata  map[string]interface{}
}
```

---

## 65.2 MongoDB Basics

### 65.2.1 Concepto de Collections y Documents

MongoDB organiza datos en **colecciones** (similar a tablas) que contienen **documentos** (similar a registros), pero con flexibilidad:

```
MongoDB Server
├── Database: "tienda"
│   ├── Collection: "productos"
│   │   ├── Document { _id: 1, nombre: "Laptop", precio: 999 }
│   │   ├── Document { _id: 2, nombre: "Mouse", precio: 29 }
│   │   └── Document { _id: 3, nombre: "Monitor", stock: 5 }
│   ├── Collection: "usuarios"
│   │   ├── Document { _id: ObjectId(...), email: "user@example.com" }
│   │   └── ...
│   └── Collection: "ordenes"
│
└── Database: "blog"
    ├── Collection: "posts"
    └── Collection: "comentarios"
```

### 65.2.2 Estructura de Documentos

Un documento MongoDB es esencialmente un JSON-like object:

```javascript
{
  "_id": ObjectId("507f1f77bcf86cd799439011"),
  "nombre": "Juan García",
  "email": "juan@example.com",
  "edad": 28,
  "activo": true,
  "direccion": {
    "calle": "Av. Principal 123",
    "ciudad": "Madrid",
    "codigo_postal": "28001"
  },
  "telefonos": ["91-234-5678", "600-123-456"],
  "metadata": {
    "creado": ISODate("2024-01-15T10:30:00Z"),
    "actualizado": ISODate("2024-01-20T14:45:00Z"),
    "tags": ["premium", "verificado"]
  }
}
```

**Ventajas de esta estructura:**
- Datos relacionados colocados juntos (embedding)
- Sin necesidad de joins complejos
- Flexible para evolucionar esquema

### 65.2.3 ObjectID

MongoDB genera automáticamente `_id` usando ObjectID, un identificador único de 12 bytes:

```
┌──────────────────────────────────────┐
│      ObjectID (12 bytes)              │
├──────────┬──────────┬──────┬──────────┤
│ Timestamp│Machine ID│Process│ Counter  │
│ 4 bytes  │ 3 bytes  │2bytes│ 3 bytes  │
└──────────┴──────────┴──────┴──────────┘

507f1f77bcf86cd799439011 =
│   │   │   │ │   │   │ │   │   │ │   │
Timestamp │ Machine └─ Counter
          Process
```

**Propiedades:**
- Único globally
- Ordenable cronológicamente
- Generado por el cliente (distribuido)
- Codificado en hexadecimal

```go
package main

import (
    "fmt"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func ExampleObjectID() {
    // Generar nuevo ObjectID
    id := primitive.NewObjectID()
    fmt.Println("ID:", id.Hex())
    
    // Obtener timestamp
    timestamp := id.Timestamp()
    fmt.Println("Creado:", timestamp)
    
    // Parsear string a ObjectID
    parsedID, _ := primitive.ObjectIDFromHexString("507f1f77bcf86cd799439011")
    fmt.Println("Parseado:", parsedID)
}
```

### 65.2.4 Tipos BSON

MongoDB internamente usa BSON (Binary JSON) para serialización. Los tipos soportados son:

| Tipo | Descripción | Ejemplo |
|------|-------------|---------|
| Double | Número flotante 64-bit | 3.14 |
| String | UTF-8 string | "Hola" |
| Object | Documento anidado | { campo: valor } |
| Array | Array de valores | [1, 2, 3] |
| BinaryData | Datos binarios | Buffer |
| ObjectID | ID único | ObjectId(...) |
| Boolean | true/false | true |
| Date | ISO 8601 date | ISODate(...) |
| Null | Nulo | null |
| Regex | Expresión regular | /pattern/flags |
| Int32 | Entero 32-bit | 42 |
| Int64 | Entero 64-bit | 9223372036854775807 |
| Timestamp | Timestamp MongoDB | Timestamp(time, inc) |
| Decimal128 | Decimal 128-bit | Decimal128(...) |

```go
package main

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func ExampleBSONTypes() {
    documento := map[string]interface{}{
        "nombre":    "Producto",
        "precio":    99.99,  // Double
        "cantidad":  int32(5),  // Int32
        "id":        primitive.NewObjectID(),  // ObjectID
        "activo":    true,  // Boolean
        "created":   time.Now(),  // Date
        "tags":      []string{"nuevo", "oferta"},  // Array
        "especif": bson.M{  // Nested Object
            "color": "rojo",
            "tamaño": "grande",
        },
    }
}
```

### 65.2.5 Arquitectura de MongoDB

MongoDB utiliza una arquitectura cliente-servidor con replicación integrada:

```
┌────────────────────────────────────────┐
│         MongoDB Architecture             │
├────────────────────────────────────────┤
│                                         │
│  Application Layer (Go Driver)          │
│         ↓                               │
│  Query Language Parser                  │
│         ↓                               │
│  ┌─────────────────────────────────┐   │
│  │   MongoDB Server (mongod)        │   │
│  ├─────────────────────────────────┤   │
│  │ Query Engine                    │   │
│  │ Index Manager                   │   │
│  │ Storage Engine (WiredTiger)     │   │
│  │ Replication & Sharding          │   │
│  └─────────────────────────────────┘   │
│         ↓                               │
│  ┌─────────────────────────────────┐   │
│  │   File System / OS Storage       │   │
│  │   • Data files (.wt)             │   │
│  │   • Index files                  │   │
│  │   • Journal                      │   │
│  └─────────────────────────────────┘   │
│                                         │
└────────────────────────────────────────┘
```

### 65.2.6 Replicación

**Replica Set:** Grupo de servidores MongoDB sincronizados:

```
┌────────────────────────────────────────┐
│           Replica Set Topology           │
├────────────────────────────────────────┤
│                                         │
│        ┌─────────────────┐              │
│        │    PRIMARY      │◄─ Escrituras │
│        │ (mongod:27017)  │              │
│        └────────┬────────┘              │
│                 │ Sincronización       │
│        ┌────────┴──────────┐            │
│        │                   │            │
│   ┌────▼───┐          ┌────▼───┐       │
│   │SECONDARY│          │SECONDARY       │
│   │(27018)  │          │(27019)  │     │
│   └─────────┘          └─────────┘     │
│   (lectura)             (lectura)      │
│                                         │
│  Oplog (registro de operaciones)       │
│  sincronizado entre todos              │
│                                         │
└────────────────────────────────────────┘
```

**Beneficios:**
- Alta disponibilidad (failover automático)
- Escalabilidad de lectura (réplicas secundarias)
- Durabilidad (datos en múltiples servidores)

### 65.2.7 Sharding

Para datos masivos, MongoDB distribuye colecciones entre múltiples servidores:

```
┌─────────────────────────────────────────┐
│         Sharding Architecture            │
├─────────────────────────────────────────┤
│                                         │
│  Client Application                    │
│         ↓                              │
│  ┌─────────────────────────┐          │
│  │  Mongos (Router)        │          │
│  │  Shard key: user_id     │          │
│  └────────┬────────┬───────┘          │
│           │        │                  │
│      ┌────▼──┐ ┌──▼────┐ ┌──────┐   │
│      │Shard 1│ │Shard 2│ │Shard 3   │
│      │ 0-100 │ │101-200│ │201-300   │
│      └───────┘ └───────┘ └──────┘   │
│                                         │
│  Config Servers (metadatos)            │
│  • Mapeo shard key → shard             │
│  • Rango de datos en cada shard        │
│                                         │
└─────────────────────────────────────────┘
```

**Tipos de Shard Key:**
- **Rango:** Datos contiguos en mismo shard (riesgo hot spot)
- **Hash:** Distribución uniforme (ideal)
- **Geoespacial:** Basado en ubicación

---

## 65.3 Driver Setup

### 65.3.1 Instalación del Driver

```bash
# Go MongoDB Driver oficial
go get go.mongodb.org/mongo-driver

# Versión específica
go get go.mongodb.org/mongo-driver@v1.12.1
```

**Dependencias principales:**
```go
import (
    "go.mongodb.org/mongo-driver/mongo"              // Cliente
    "go.mongodb.org/mongo-driver/mongo/options"      // Configuración
    "go.mongodb.org/mongo-driver/bson"               // Serialización
    "go.mongodb.org/mongo-driver/bson/primitive"     // Tipos BSON
)
```

### 65.3.2 Connection String

La URI de conexión especifica dónde y cómo conectarse:

```
mongodb://[usuario:contraseña@]host1[:puerto1][,host2[:puerto2],...][/database][?opciones]

Ejemplos:
┌────────────────────────────────────────────────────┐
│ mongodb://localhost:27017                          │
│ Local simple, sin autenticación                    │
│                                                     │
│ mongodb://user:pass@mongodb.example.com:27017      │
│ Conexión remota con autenticación                  │
│                                                     │
│ mongodb+srv://user:pass@cluster.mongodb.net        │
│ Atlas MongoDB (DNS SRV)                            │
│                                                     │
│ mongodb://host1,host2,host3/?replicaSet=rs0        │
│ Replica Set                                         │
└────────────────────────────────────────────────────┘
```

### 65.3.3 Configuración del Cliente

```go
package main

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, error) {
    // Crear contexto con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Configurar opciones del cliente
    opts := options.Client().
        ApplyURI("mongodb://localhost:27017").
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(5 * time.Minute)
    
    // Conectar
    client, err := mongo.Connect(ctx, opts)
    if err != nil {
        return nil, err
    }
    
    // Verificar conexión
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }
    
    return client, nil
}

// Uso recomendado: defer defer
func main() {
    client, _ := ConnectToMongoDB()
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        client.Disconnect(ctx)
    }()
}
```

### 65.3.4 Autenticación

```go
// Autenticación SCRAM (estándar)
opts := options.Client().
    ApplyURI("mongodb://usuario:contraseña@host:27017/database?authSource=admin")

// O configurar manualmente
auth := options.Credential{
    Username: "usuario",
    Password: "contraseña",
    AuthSource: "admin",  // Base de datos con usuario
}
opts := options.Client().
    SetAuth(auth)

// LDAP (Enterprise)
auth := options.Credential{
    Username:   "usuario@LDAP",
    Password:   "contraseña",
    AuthMechanism: "PLAIN",
}
```

### 65.3.5 Connection Pooling

MongoDB mantiene automáticamente un pool de conexiones:

```go
// Tamaño del pool
opts := options.Client().
    SetMaxPoolSize(100).        // Máximo
    SetMinPoolSize(10)          // Mínimo (siempre activas)

// Timeout de conexión
opts = options.Client().
    SetServerSelectionTimeout(5 * time.Second).
    SetMaxConnIdleTime(5 * time.Minute)

// Monitorar eventos de pool
opts = options.Client().
    SetMonitor(&event.CommandMonitor{
        Started: func(_ context.Context, evt *event.CommandStartedEvent) {
            log.Printf("Comando: %v", evt.CommandName)
        },
        Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
            log.Printf("Completado en: %d ns", evt.DurationNanos)
        },
    })
```

### 65.3.6 Estructura de Acceso a Datos

```go
package main

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
)

// Estructura global (mejor alternativa: inyección)
type Database struct {
    Client *mongo.Client
    DB     *mongo.Database
}

// Constructor recomendado
func NewDatabase(uri string) (*Database, error) {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    
    return &Database{
        Client: client,
        DB:     client.Database("miapp"),
    }, nil
}

// Método de acceso a colecciones
func (db *Database) GetUsersCollection() *mongo.Collection {
    return db.DB.Collection("usuarios")
}

// Cerrar conexión
func (db *Database) Disconnect(ctx context.Context) error {
    return db.Client.Disconnect(ctx)
}
```

---

## 65.4 Operaciones CRUD

### 65.4.1 Insert Operations

#### InsertOne
```go
package main

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
)

type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    Name  string             `bson:"name"`
    Email string             `bson:"email"`
    Age   int                `bson:"age"`
}

func (db *Database) InsertUser(ctx context.Context, user User) (string, error) {
    collection := db.DB.Collection("usuarios")
    
    result, err := collection.InsertOne(ctx, user)
    if err != nil {
        return "", err
    }
    
    // Obtener el ID generado
    id := result.InsertedID.(primitive.ObjectID)
    return id.Hex(), nil
}

// Uso
func ExampleInsertOne(db *Database) {
    ctx := context.Background()
    
    usuario := User{
        Name:  "Ana García",
        Email: "ana@example.com",
        Age:   28,
    }
    
    id, err := db.InsertUser(ctx, usuario)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Usuario insertado con ID: %s\n", id)
}
```

#### InsertMany
```go
func (db *Database) InsertUsers(ctx context.Context, users []User) ([]string, error) {
    collection := db.DB.Collection("usuarios")
    
    // Convertir a []interface{} (requerido por InsertMany)
    docs := make([]interface{}, len(users))
    for i, user := range users {
        docs[i] = user
    }
    
    // Insertar con opciones
    opts := options.InsertMany().SetOrdered(false)  // No ordenado = más rápido
    result, err := collection.InsertMany(ctx, docs, opts)
    if err != nil {
        return nil, err
    }
    
    // Extraer IDs
    ids := make([]string, len(result.InsertedIDs))
    for i, id := range result.InsertedIDs {
        ids[i] = id.(primitive.ObjectID).Hex()
    }
    
    return ids, nil
}

// Uso
func ExampleInsertMany(db *Database) {
    usuarios := []User{
        {Name: "Carlos", Email: "carlos@example.com", Age: 35},
        {Name: "María", Email: "maria@example.com", Age: 29},
        {Name: "Juan", Email: "juan@example.com", Age: 31},
    }
    
    ids, err := db.InsertUsers(context.Background(), usuarios)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Insertados %d usuarios: %v\n", len(ids), ids)
}
```

**Manejo de errores en InsertMany:**
```go
// insertMany no falla si algunos documentos fallan (con SetOrdered(false))
result, err := collection.InsertMany(ctx, docs, options.InsertMany().SetOrdered(false))

if err != nil {
    // Algunos documentos pueden haberse insertado
    writeErr, ok := err.(mongo.BulkWriteError)
    if ok {
        fmt.Printf("Insertados: %d\n", len(result.InsertedIDs))
        fmt.Printf("Errores: %v\n", writeErr.WriteErrors)
    }
}
```

### 65.4.2 Read Operations

#### FindOne
```go
func (db *Database) FindUserByID(ctx context.Context, id string) (*User, error) {
    collection := db.DB.Collection("usuarios")
    
    objID, _ := primitive.ObjectIDFromHexString(id)
    
    var user User
    err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("usuario no encontrado")
        }
        return nil, err
    }
    
    return &user, nil
}

func (db *Database) FindUserByEmail(ctx context.Context, email string) (*User, error) {
    collection := db.DB.Collection("usuarios")
    
    var user User
    err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

#### Find (múltiples documentos)
```go
func (db *Database) FindUsersOlderThan(ctx context.Context, age int) ([]User, error) {
    collection := db.DB.Collection("usuarios")
    
    filter := bson.M{"age": bson.M{"$gt": age}}
    
    // Opciones: ordenamiento, límite, skip
    opts := options.Find().
        SetSort(bson.M{"age": -1}).  // -1 descendente, 1 ascendente
        SetLimit(100)
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var users []User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    
    return users, nil
}

// Alternativa: iterar
func (db *Database) FindUsersIterating(ctx context.Context, filter bson.M) ([]User, error) {
    collection := db.DB.Collection("usuarios")
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var users []User
    
    // Iterar resultado
    for cursor.Next(ctx) {
        var user User
        if err := cursor.Decode(&user); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, cursor.Err()
}
```

### 65.4.3 Update Operations

#### UpdateOne
```go
func (db *Database) UpdateUserAge(ctx context.Context, email string, newAge int) error {
    collection := db.DB.Collection("usuarios")
    
    filter := bson.M{"email": email}
    update := bson.M{
        "$set": bson.M{
            "age": newAge,
            "updated_at": time.Now(),
        },
    }
    
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return errors.New("usuario no encontrado")
    }
    
    fmt.Printf("Modificados: %d documentos\n", result.ModifiedCount)
    return nil
}
```

#### UpdateMany
```go
func (db *Database) UpdateUsersStatus(ctx context.Context, ageThreshold int, status string) error {
    collection := db.DB.Collection("usuarios")
    
    filter := bson.M{"age": bson.M{"$gte": ageThreshold}}
    update := bson.M{
        "$set": bson.M{"status": status},
        "$inc": bson.M{"update_count": 1},  // Incrementar contador
    }
    
    result, err := collection.UpdateMany(ctx, filter, update)
    if err != nil {
        return err
    }
    
    fmt.Printf("Documentos coincidentes: %d\n", result.MatchedCount)
    fmt.Printf("Documentos modificados: %d\n", result.ModifiedCount)
    return nil
}
```

**Operadores de actualización:**

| Operador | Descripción |
|----------|-------------|
| `$set` | Establecer campo |
| `$unset` | Remover campo |
| `$inc` | Incrementar |
| `$dec` | Decrementar |
| `$push` | Agregar a array |
| `$pop` | Remover de array |
| `$pull` | Remover elementos del array |
| `$addToSet` | Agregar a array (sin duplicados) |
| `$rename` | Renombrar campo |
| `$min`/`$max` | Solo si es menor/mayor |

#### Replace (Reemplazo completo)
```go
func (db *Database) ReplaceUser(ctx context.Context, id string, newUser User) error {
    collection := db.DB.Collection("usuarios")
    
    objID, _ := primitive.ObjectIDFromHexString(id)
    
    // ReplaceOne reemplaza TODO el documento (excepto _id)
    result, err := collection.ReplaceOne(
        ctx,
        bson.M{"_id": objID},
        newUser,
    )
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return errors.New("usuario no encontrado")
    }
    
    return nil
}
```

### 65.4.4 Delete Operations

#### DeleteOne
```go
func (db *Database) DeleteUserByID(ctx context.Context, id string) error {
    collection := db.DB.Collection("usuarios")
    
    objID, _ := primitive.ObjectIDFromHexString(id)
    
    result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        return err
    }
    
    if result.DeletedCount == 0 {
        return errors.New("usuario no encontrado")
    }
    
    return nil
}
```

#### DeleteMany
```go
func (db *Database) DeleteInactiveUsers(ctx context.Context, daysInactive int) (int64, error) {
    collection := db.DB.Collection("usuarios")
    
    // Usuarios sin actualizar hace más de X días
    threshold := time.Now().AddDate(0, 0, -daysInactive)
    
    filter := bson.M{
        "last_login": bson.M{"$lt": threshold},
    }
    
    result, err := collection.DeleteMany(ctx, filter)
    if err != nil {
        return 0, err
    }
    
    return result.DeletedCount, nil
}
```

### 65.4.5 Upsert (Insert o Update)

```go
func (db *Database) UpsertUser(ctx context.Context, user User) error {
    collection := db.DB.Collection("usuarios")
    
    opts := options.Update().SetUpsert(true)
    
    filter := bson.M{"email": user.Email}
    update := bson.M{
        "$set": user,
        "$setOnInsert": bson.M{"created_at": time.Now()},
    }
    
    result, err := collection.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        return err
    }
    
    if result.UpsertedID != nil {
        fmt.Println("Nuevo usuario creado")
    } else {
        fmt.Println("Usuario actualizado")
    }
    
    return nil
}
```

---

## 65.5 Querying Avanzado

### 65.5.1 Operadores de Comparación

```go
// MongoDB vs Go
/*
┌──────────────┬─────────────────┬────────────────────────┐
│ Operador     │ SQL             │ MongoDB                 │
├──────────────┼─────────────────┼────────────────────────┤
│ Igual        │ =               │ { field: value }        │
│ Mayor        │ >               │ { field: {$gt: value}}  │
│ Menor        │ <               │ { field: {$lt: value}}  │
│ Mayor/igual  │ >=              │ { field: {$gte: value}} │
│ Menor/igual  │ <=              │ { field: {$lte: value}} │
│ Diferente    │ !=              │ { field: {$ne: value}}  │
│ En           │ IN (...)        │ { field: {$in: [...]}} │
│ No en        │ NOT IN (...)    │ { field: {$nin: [...]}}│
└──────────────┴─────────────────┴────────────────────────┘
*/

// Ejemplos de comparación
filter := bson.M{
    "age": bson.M{"$gt": 18},              // age > 18
    "salary": bson.M{"$gte": 50000},       // salary >= 50000
    "department": bson.M{"$in": []string{"IT", "HR"}},
    "status": bson.M{"$ne": "inactive"},
}
```

### 65.5.2 Operadores Lógicos

```go
// $and, $or, $nor, $not

// Usuarios de IT O Finance con salario > 60000
filter := bson.M{
    "$or": []bson.M{
        {"department": "IT"},
        {"department": "Finance"},
    },
    "salary": bson.M{"$gt": 60000},
}

// NOT: edad mayor que 65 O menor que 18
filter = bson.M{
    "$nor": []bson.M{
        {"age": bson.M{"$gt": 65}},
        {"age": bson.M{"$lt": 18}},
    },
}

// AND explícito
filter = bson.M{
    "$and": []bson.M{
        {"active": true},
        {"verified": true},
        {"premium": true},
    },
}
```

### 65.5.3 Operadores de String

```go
// $regex: búsqueda de texto con expresiones regulares

// Usuarios cuyo nombre empieza con "An"
filter := bson.M{
    "name": bson.M{"$regex": "^An", "$options": "i"},  // i = case-insensitive
}

// Contiene "@gmail.com"
filter = bson.M{
    "email": bson.M{"$regex": "@gmail\\.com$"},
}

// $text: búsqueda de texto con índice (más eficiente)
// Requiere índice de texto: collection.Indexes().CreateOne(ctx, 
//   mongo.IndexModel{Keys: bson.D{{Key: "name", Value: "text"}}})
filter = bson.M{
    "$text": bson.M{"$search": "mongodb tutorial"},
}
```

### 65.5.4 Operadores de Array

```go
// $elemMatch: elemento que cumple todas las condiciones

type Order struct {
    ID    primitive.ObjectID `bson:"_id"`
    Items []struct {
        Product string
        Price   float64
    } `bson:"items"`
}

// Órdenes con al menos un item de más de $100
filter := bson.M{
    "items": bson.M{
        "$elemMatch": bson.M{
            "price": bson.M{"$gt": 100},
        },
    },
}

// $size: número de elementos
filter = bson.M{
    "items": bson.M{"$size": 3},  // Array exactamente 3 items
}

// $all: contiene todos los elementos
filter = bson.M{
    "tags": bson.M{"$all": []string{"python", "mongodb"}},
}
```

### 65.5.5 Proyección (Select en SQL)

```go
// Incluir/excluir campos en resultado

// Incluir solo name y email
opts := options.Find().SetProjection(bson.M{
    "name": 1,
    "email": 1,
})
// Resultado: solo estos campos + _id (siempre incluido)

// Excluir password
opts = options.Find().SetProjection(bson.M{
    "password": 0,
    "secret_key": 0,
})
// Resultado: todos excepto estos campos

// Excluir _id
opts = options.Find().SetProjection(bson.M{
    "_id": 0,
    "name": 1,
})

// Proyección de arrays
opts = options.Find().SetProjection(bson.M{
    "name": 1,
    "tags": bson.M{"$slice": 5},  // Solo primeros 5 elementos
})
```

### 65.5.6 Sorting, Limit, Skip

```go
// Ordenamiento
opts := options.Find().
    SetSort(bson.M{
        "age": -1,      // -1 descendente
        "name": 1,      // 1 ascendente
    })

// Límite
opts = options.Find().SetLimit(50)

// Skip (para paginación)
opts = options.Find().
    SetSkip(100).
    SetLimit(50)
    // Página 3 de 50 items: skip 100, limit 50

// Combinado: paginación eficiente
func (db *Database) GetUsersPaginated(ctx context.Context, page, pageSize int64) ([]User, error) {
    collection := db.DB.Collection("usuarios")
    
    skip := (page - 1) * pageSize
    
    opts := options.Find().
        SetSkip(skip).
        SetLimit(pageSize).
        SetSort(bson.M{"created_at": -1})
    
    cursor, _ := collection.Find(ctx, bson.M{}, opts)
    defer cursor.Close(ctx)
    
    var users []User
    cursor.All(ctx, &users)
    return users, nil
}
```

---

## 65.6 Aggregation Pipeline

### 65.6.1 Concepto y Arquitectura

La agregación es el procesamiento de datos en etapas. Cada etapa transforma datos:

```
Colección
    ↓
┌───────────────┐
│ $match        │  Filtrar documentos (como WHERE)
└───────────────┘
    ↓
┌───────────────┐
│ $project      │  Seleccionar/calcular campos
└───────────────┘
    ↓
┌───────────────┐
│ $group        │  Agrupar y agregar
└───────────────┘
    ↓
┌───────────────┐
│ $sort         │  Ordenar
└───────────────┘
    ↓
┌───────────────┐
│ $limit        │  Límite
└───────────────┘
    ↓
Resultado final
```

### 65.6.2 $match y $project

```go
func (db *Database) AggregateExample(ctx context.Context) []bson.M {
    collection := db.DB.Collection("usuarios")
    
    // Pipeline de agregación
    pipeline := mongo.Pipeline{
        // Etapa 1: Filtrar usuarios activos mayores de 25 años
        bson.D{
            {Key: "$match", Value: bson.D{
                {Key: "active", Value: true},
                {Key: "age", Value: bson.D{{Key: "$gt", Value: 25}}},
            }},
        },
        
        // Etapa 2: Proyectar solo campos necesarios
        bson.D{
            {Key: "$project", Value: bson.D{
                {Key: "name", Value: 1},
                {Key: "email", Value: 1},
                {Key: "age", Value: 1},
                {Key: "password", Value: 0},  // Excluir
                {Key: "age_group", Value: bson.D{  // Campo calculado
                    {Key: "$cond", Value: bson.A{
                        bson.D{{Key: "$gte", Value: bson.A{"$age", 30}}},
                        "Senior",
                        "Junior",
                    }},
                }},
            }},
        },
    }
    
    cursor, _ := collection.Aggregate(ctx, pipeline)
    defer cursor.Close(ctx)
    
    var results []bson.M
    cursor.All(ctx, &results)
    return results
}
```

### 65.6.3 $group y Acumuladores

```go
func (db *Database) GroupByDepartment(ctx context.Context) []bson.M {
    collection := db.DB.Collection("empleados")
    
    pipeline := mongo.Pipeline{
        bson.D{
            {Key: "$group", Value: bson.D{
                {Key: "_id", Value: "$department"},  // Agrupar por
                {Key: "total_empleados", Value: bson.D{  // $sum: 1 = contar
                    {Key: "$sum", Value: 1},
                }},
                {Key: "salario_promedio", Value: bson.D{
                    {Key: "$avg", Value: "$salary"},
                }},
                {Key: "salario_min", Value: bson.D{
                    {Key: "$min", Value: "$salary"},
                }},
                {Key: "salario_max", Value: bson.D{
                    {Key: "$max", Value: "$salary"},
                }},
                {Key: "nombres", Value: bson.D{  // Recopilar valores
                    {Key: "$push", Value: "$name"},
                }},
            }},
        },
        bson.D{
            {Key: "$sort", Value: bson.D{
                {Key: "total_empleados", Value: -1},
            }},
        },
    }
    
    cursor, _ := collection.Aggregate(ctx, pipeline)
    var results []bson.M
    cursor.All(ctx, &results)
    return results
}

/* Resultado:
{
  _id: "IT",
  total_empleados: 15,
  salario_promedio: 75000,
  salario_min: 50000,
  salario_max: 120000,
  nombres: ["Alice", "Bob", ...]
}
*/
```

**Acumuladores disponibles:**

| Operador | Descripción |
|----------|-------------|
| `$sum` | Suma valores |
| `$avg` | Promedio |
| `$min` | Mínimo |
| `$max` | Máximo |
| `$first` | Primer valor |
| `$last` | Último valor |
| `$push` | Array con todos los valores |
| `$addToSet` | Array con valores únicos |
| `$count` | Contar (en $group, usar `$sum: 1`) |

### 65.6.4 $lookup (Joins)

```go
func (db *Database) JoinUsersWithOrders(ctx context.Context) []bson.M {
    collection := db.DB.Collection("usuarios")
    
    // Simular JOIN: usuarios LEFT OUTER JOIN ordenes
    pipeline := mongo.Pipeline{
        bson.D{
            {Key: "$lookup", Value: bson.D{
                {Key: "from", Value: "ordenes"},          // Tabla a unir
                {Key: "localField", Value: "_id"},        // Campo local
                {Key: "foreignField", Value: "user_id"},  // Campo remoto
                {Key: "as", Value: "ordenes"},            // Nombre del array
            }},
        },
        bson.D{
            {Key: "$project", Value: bson.D{
                {Key: "name", Value: 1},
                {Key: "email", Value: 1},
                {Key: "ordenes_count", Value: bson.D{
                    {Key: "$size", Value: "$ordenes"},
                }},
            }},
        },
    }
    
    cursor, _ := collection.Aggregate(ctx, pipeline)
    var results []bson.M
    cursor.All(ctx, &results)
    return results
}

/* Resultado:
{
  _id: ObjectId(...),
  name: "Juan",
  email: "juan@example.com",
  ordenes_count: 3,
  ordenes: [
    { _id: ..., user_id: ..., total: 100 },
    ...
  ]
}
*/
```

### 65.6.5 $unwind

```go
// $unwind: expandir arrays en múltiples documentos

func (db *Database) UnwindExample(ctx context.Context) {
    collection := db.DB.Collection("ordenes")
    
    // Antes: { _id: 1, items: [A, B, C] }
    // Después: 
    //   { _id: 1, items: A }
    //   { _id: 1, items: B }
    //   { _id: 1, items: C }
    
    pipeline := mongo.Pipeline{
        bson.D{
            {Key: "$unwind", Value: "$items"},
        },
        bson.D{
            {Key: "$group", Value: bson.D{
                {Key: "_id", Value: "$items"},
                {Key: "veces_vendido", Value: bson.D{
                    {Key: "$sum", Value: 1},
                }},
            }},
        },
    }
    
    cursor, _ := collection.Aggregate(ctx, pipeline)
    var results []bson.M
    cursor.All(ctx, &results)
    // Resultado: cada item con su contador de ventas
}
```

### 65.6.6 $facet (Búsquedas Facetadas)

```go
func (db *Database) FacetedSearch(ctx context.Context) bson.M {
    collection := db.DB.Collection("productos")
    
    pipeline := mongo.Pipeline{
        bson.D{
            {Key: "$match", Value: bson.D{
                {Key: "precio", Value: bson.D{{Key: "$gt", Value: 10}}},
            }},
        },
        bson.D{
            {Key: "$facet", Value: bson.D{
                // Múltiples análisis simultáneamente
                {Key: "por_categoria", Value: mongo.Pipeline{
                    bson.D{
                        {Key: "$group", Value: bson.D{
                            {Key: "_id", Value: "$category"},
                            {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
                        }},
                    },
                }},
                {Key: "rango_precio", Value: mongo.Pipeline{
                    bson.D{
                        {Key: "$group", Value: bson.D{
                            {Key: "_id", Value: "$price_range"},
                            {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
                        }},
                    },
                }},
            }},
        },
    }
    
    cursor, _ := collection.Aggregate(ctx, pipeline)
    var result bson.M
    cursor.Next(ctx)
    cursor.Decode(&result)
    return result
}
```

---

## 65.7 Índices y Performance

### 65.7.1 Creación de Índices

```go
func (db *Database) CreateIndexes(ctx context.Context) error {
    collection := db.DB.Collection("usuarios")
    
    indexModel := mongo.IndexModel{
        Keys: bson.D{{Key: "email", Value: 1}},
        Options: options.Index().SetUnique(true),  // Único
    }
    
    indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return err
    }
    
    fmt.Println("Índice creado:", indexName)
    return nil
}

// Múltiples índices
func (db *Database) CreateMultipleIndexes(ctx context.Context) error {
    collection := db.DB.Collection("ordenes")
    
    indexModels := []mongo.IndexModel{
        {
            Keys: bson.D{{Key: "user_id", Value: 1}},  // Índice simple
        },
        {
            Keys: bson.D{  // Índice compuesto
                {Key: "user_id", Value: 1},
                {Key: "created_at", Value: -1},
            },
        },
        {
            Keys: bson.D{{Key: "email", Value: "text"}},  // Índice de texto
        },
    }
    
    opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
    names, err := collection.Indexes().CreateMany(ctx, indexModels, opts)
    if err != nil {
        return err
    }
    
    for _, name := range names {
        fmt.Println("Creado:", name)
    }
    return nil
}
```

### 65.7.2 Tipos de Índices

| Tipo | Descripción | Caso de uso |
|------|-------------|------------|
| **Single Field** | Índice en un campo | Búsquedas comunes |
| **Compound** | Índice en múltiples campos | Queries con múltiples filtros |
| **Text** | Búsqueda de texto | Búsqueda por palabras clave |
| **Geospatial** | Índice espacial | Queries por ubicación |
| **Hashed** | Hash del valor | Sharding distribuido |
| **Sparse** | Solo documentos con campo | Optimización cuando muchos nulos |
| **TTL** | Expira documentos automáticamente | Sesiones, logs temporales |

### 65.7.3 Índice TTL (Time to Live)

```go
// Expirar documentos automáticamente

func (db *Database) CreateTTLIndex(ctx context.Context) error {
    collection := db.DB.Collection("sesiones")
    
    indexModel := mongo.IndexModel{
        Keys: bson.D{{Key: "created_at", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(3600),  // 1 hora
    }
    
    _, err := collection.Indexes().CreateOne(ctx, indexModel)
    return err
}

// Documentos en colección
type Session struct {
    ID        primitive.ObjectID `bson:"_id"`
    UserID    string             `bson:"user_id"`
    Token     string             `bson:"token"`
    CreatedAt time.Time          `bson:"created_at"`
}

// Se elimina automáticamente 1 hora después de CreatedAt
```

### 65.7.4 Explicar Queries

```go
func (db *Database) ExplainQuery(ctx context.Context) {
    collection := db.DB.Collection("usuarios")
    
    filter := bson.M{"age": bson.M{"$gt": 25}}
    
    // Obtener plan de ejecución
    var explainResult bson.M
    err := collection.Database().RunCommand(ctx, bson.M{
        "explain": bson.M{
            "find": "usuarios",
            "filter": filter,
        },
    }).Decode(&explainResult)
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Analizar resultado
    stats := explainResult["executionStats"].(bson.M)
    fmt.Printf("Documentos examinados: %v\n", stats["totalDocsExamined"])
    fmt.Printf("Documentos retornados: %v\n", stats["nReturned"])
    fmt.Printf("Eficiencia: %.2f%%\n", 
        float64(stats["nReturned"].(int32)) / 
        float64(stats["totalDocsExamined"].(int32)) * 100)
}
```

**Interpretación:**
- `nReturned` / `totalDocsExamined` ≈ 1.0 = Índice eficiente
- `stage` = "COLLECTION_SCAN" = No usa índice ⚠️
- `stage` = "IXSCAN" = Usa índice ✅

### 65.7.5 Best Practices de Performance

```
❌ Anti-pattern: Sin índices
func (db *Database) SlowQuery(ctx context.Context) {
    // Scans toda la colección
    collection.Find(ctx, bson.M{"email": email})
}

✅ Best practice: Con índice
func (db *Database) FastQuery(ctx context.Context) {
    // Índice creado: collection.Indexes().CreateOne(...,
    //   mongo.IndexModel{Keys: bson.D{{Key: "email", Value: 1}}})
    collection.Find(ctx, bson.M{"email": email})
}

❌ Anti-pattern: Índice en campo ineficaz
indexModel := mongo.IndexModel{
    Keys: bson.D{{Key: "status", Value: 1}},  // Solo 5 valores diferentes
}

✅ Best practice: Índice selectivo
indexModel := mongo.IndexModel{
    Keys: bson.D{{Key: "email", Value: 1}},  // Muchos valores únicos
    Options: options.Index().SetUnique(true),
}

❌ Anti-pattern: Demasiados índices
// 20+ índices ralentizan inserciones/actualizaciones

✅ Best practice: Índices estratégicos
// Solo en campos frecuentemente consultados

❌ Anti-pattern: Índice no selectivo
filter := bson.M{"active": true}  // 95% de documentos

✅ Best practice: Índice selectivo + query específica
indexModel := mongo.IndexModel{
    Keys: bson.D{
        {Key: "active", Value: 1},
        {Key: "created_at", Value: -1},
    },
    Options: options.Index().SetPartialFilterExpression(
        bson.M{"active": true},
    ),
}
```

---

## 65.8 Transacciones

### 65.8.1 Concepto de Sesiones y Transacciones

```
┌───────────────────────────────┐
│      Session (Sesión)          │
│                               │
│  ┌──────────────────────────┐ │
│  │  Transacción (opcional)  │ │
│  │  - Multi-document ACID   │ │
│  │  - Rollback disponible   │ │
│  └──────────────────────────┘ │
│                               │
│  Operaciones:                 │
│  • Todas dentro de sesión     │
│  • Con punto de consistencia  │
│  • Con isolamiento            │
└───────────────────────────────┘
```

### 65.8.2 Transacciones Multi-Documento

```go
package main

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
)

func TransferFunds(client *mongo.Client, fromID, toID string, amount float64) error {
    // Crear sesión
    session, err := client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(context.Background())
    
    // Ejecutar transacción
    return mongo.WithSession(context.Background(), session, func(sessionCtx mongo.SessionContext) error {
        // Iniciar transacción
        if err = session.StartTransaction(); err != nil {
            return err
        }
        
        ctx := context.Background()
        accountsCol := client.Database("banco").Collection("cuentas")
        
        // Operación 1: Restar dinero de cuenta origen
        _, err = accountsCol.UpdateOne(sessionCtx, bson.M{"_id": fromID}, bson.M{
            "$inc": bson.M{"saldo": -amount},
        })
        if err != nil {
            session.AbortTransaction(ctx)
            return err
        }
        
        // Operación 2: Agregar dinero a cuenta destino
        _, err = accountsCol.UpdateOne(sessionCtx, bson.M{"_id": toID}, bson.M{
            "$inc": bson.M{"saldo": amount},
        })
        if err != nil {
            session.AbortTransaction(ctx)
            return err
        }
        
        // Registrar transacción en log
        transactionsCol := client.Database("banco").Collection("transacciones")
        _, err = transactionsCol.InsertOne(sessionCtx, bson.M{
            "from": fromID,
            "to":   toID,
            "amount": amount,
            "timestamp": time.Now(),
        })
        if err != nil {
            session.AbortTransaction(ctx)
            return err
        }
        
        // Confirmar transacción
        return session.CommitTransaction(ctx)
    })
}

// Uso
func ExampleTransaction(client *mongo.Client) {
    err := TransferFunds(client, "account_001", "account_002", 100.50)
    if err != nil {
        log.Printf("Error en transacción: %v", err)
    } else {
        log.Println("Transferencia exitosa")
    }
}
```

### 65.8.3 Manejo de Errores en Transacciones

```go
func TransactionWithRetry(client *mongo.Client, operation func(mongo.SessionContext) error) error {
    for attempt := 0; attempt < 3; attempt++ {
        session, err := client.StartSession()
        if err != nil {
            return err
        }
        
        err = mongo.WithSession(context.Background(), session, func(ctx mongo.SessionContext) error {
            session.StartTransaction()
            
            // Ejecutar operación
            if err := operation(ctx); err != nil {
                session.AbortTransaction(context.Background())
                return err
            }
            
            // Intentar commit (puede fallar con transient error)
            return session.CommitTransaction(context.Background())
        })
        
        session.EndSession(context.Background())
        
        if err == nil {
            return nil  // Éxito
        }
        
        // Revisar si es error transitorio
        if mongo.IsNetworkError(err) || mongo.IsServerSelectionError(err) {
            log.Printf("Error transitorio, reintentando... (%d/3)", attempt+1)
            time.Sleep(time.Second * time.Duration(attempt+1))
            continue
        }
        
        return err  // Error no transitorio
    }
    
    return errors.New("transacción falló después de 3 intentos")
}
```

### 65.8.4 Garantías ACID

MongoDB proporciona garantías ACID en transacciones:

| Propiedad | Garantía |
|-----------|----------|
| **Atomicity** | Todas las operaciones o ninguna |
| **Consistency** | Datos válidos antes/después |
| **Isolation** | No ve cambios de otras transacciones |
| **Durability** | Cambios persisten tras crash |

```go
// Ejemplo: Guarantee de consistencia
func BookFlight(db *mongo.Database, session mongo.Session) error {
    ctx := context.Background()
    
    if err := session.StartTransaction(); err != nil {
        return err
    }
    
    flightsCol := db.Collection("vuelos")
    bookingsCol := db.Collection("reservas")
    
    // Verificar disponibilidad
    var flight bson.M
    err := flightsCol.FindOne(ctx, bson.M{"_id": "FL123"}).Decode(&flight)
    if err != nil {
        session.AbortTransaction(ctx)
        return errors.New("vuelo no encontrado")
    }
    
    available := flight["available_seats"].(int32)
    if available <= 0 {
        session.AbortTransaction(ctx)
        return errors.New("sin asientos disponibles")
    }
    
    // Decrementar asientos y crear reserva (ambos o ninguno)
    flightsCol.UpdateOne(ctx, bson.M{"_id": "FL123"}, 
        bson.M{"$inc": bson.M{"available_seats": -1}})
    
    bookingsCol.InsertOne(ctx, bson.M{
        "user_id": "user123",
        "flight": "FL123",
        "status": "confirmed",
    })
    
    return session.CommitTransaction(ctx)
}
```

---

## 65.9 Operaciones Bulk y Optimización

### 65.9.1 Bulk Write Operations

```go
package main

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
)

func BulkOperationsExample(db *mongo.Database, ctx context.Context) error {
    collection := db.Collection("usuarios")
    
    // Crear modelos de operaciones
    models := []mongo.WriteModel{
        // Insertar
        mongo.NewInsertOneModel().SetDocument(bson.M{
            "name": "Alice",
            "email": "alice@example.com",
        }),
        
        // Actualizar
        mongo.NewUpdateOneModel().
            SetFilter(bson.M{"email": "bob@example.com"}).
            SetUpdate(bson.M{"$set": bson.M{"status": "active"}}),
        
        // Reemplazar
        mongo.NewReplaceOneModel().
            SetFilter(bson.M{"_id": "charlie"}).
            SetReplacement(bson.M{
                "name": "Charlie",
                "email": "charlie@example.com",
            }).
            SetUpsert(true),
        
        // Eliminar
        mongo.NewDeleteOneModel().
            SetFilter(bson.M{"status": "inactive"}),
    }
    
    opts := options.BulkWrite().SetOrdered(false)  // No ordenado = paralelo
    result, err := collection.BulkWrite(ctx, models, opts)
    if err != nil {
        return err
    }
    
    fmt.Printf("Insertados: %d\n", result.InsertedCount)
    fmt.Printf("Modificados: %d\n", result.ModifiedCount)
    fmt.Printf("Eliminados: %d\n", result.DeletedCount)
    
    return nil
}
```

### 65.9.2 Procesamiento en Lotes

```go
func BatchProcessing(db *mongo.Database, ctx context.Context) error {
    collection := db.Collection("eventos")
    
    // Procesar millones de documentos eficientemente
    const batchSize = 1000
    
    filter := bson.M{"processed": false}
    opts := options.Find().SetBatchSize(batchSize)
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    var batch []bson.M
    
    for cursor.Next(ctx) {
        var doc bson.M
        cursor.Decode(&doc)
        
        batch = append(batch, doc)
        
        // Procesar cuando alcanza tamaño del lote
        if len(batch) >= batchSize {
            if err := processBatch(db, ctx, batch); err != nil {
                return err
            }
            batch = batch[:0]  // Limpiar
        }
    }
    
    // Procesar último lote
    if len(batch) > 0 {
        return processBatch(db, ctx, batch)
    }
    
    return nil
}

func processBatch(db *mongo.Database, ctx context.Context, batch []bson.M) error {
    collection := db.Collection("eventos")
    
    // Actualizar todos los documentos en lote
    ids := make([]interface{}, len(batch))
    for i, doc := range batch {
        ids[i] = doc["_id"]
    }
    
    _, err := collection.UpdateMany(ctx,
        bson.M{"_id": bson.M{"$in": ids}},
        bson.M{"$set": bson.M{"processed": true, "processed_at": time.Now()}},
    )
    
    return err
}
```

### 65.9.3 Optimización de Conexiones

```go
func OptimizedConnections(uri string) (*mongo.Client, error) {
    opts := options.Client().
        ApplyURI(uri).
        // Pool de conexiones
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(5 * time.Minute).
        // Timeouts
        SetServerSelectionTimeout(5 * time.Second).
        SetSocketTimeout(10 * time.Second).
        SetConnectTimeout(10 * time.Second).
        // Retry automático
        SetRetryWrites(true).
        SetRetryReads(true).
        // Compresión (reduce uso de red)
        SetCompressors([]string{"snappy", "zstd"})
    
    return mongo.Connect(context.Background(), opts)
}

// Monitorear pool
type PoolMonitor struct{}

func (pm *PoolMonitor) Event(evt *event.PoolEvent) {
    switch evt := evt.(type) {
    case *event.ConnectionCheckedInEvent:
        log.Println("Conexión devuelta al pool")
    case *event.ConnectionCheckedOutEvent:
        log.Println("Conexión sacada del pool")
    case *event.PoolClearedEvent:
        log.Println("Pool limpiado")
    }
}
```

### 65.9.4 Considerations de Memoria

```go
func MemoryEfficientQuery(db *mongo.Database, ctx context.Context) error {
    collection := db.Collection("productos")
    
    // ❌ Ineficiente: Cargar todo en memoria
    cursor, _ := collection.Find(ctx, bson.M{})
    defer cursor.Close(ctx)
    var allProducts []bson.M
    cursor.All(ctx, &allProducts)  // Todo en RAM
    
    // ✅ Eficiente: Iterar
    cursor, _ := collection.Find(ctx, bson.M{})
    defer cursor.Close(ctx)
    
    for cursor.Next(ctx) {
        var product bson.M
        cursor.Decode(&product)
        // Procesar sin acumular
    }
    
    return nil
}

// ✅ Proyección para reducir tamaño
opts := options.Find().SetProjection(bson.M{
    "name": 1,
    "price": 1,
    // Excluye campos grandes
    "description": 0,
    "image": 0,
})
cursor, _ := collection.Find(ctx, bson.M{}, opts)
```

---

## 65.10 Validación de Esquema

### 65.10.1 JSON Schema Validation

```go
func CreateCollectionWithValidation(db *mongo.Database, ctx context.Context) error {
    validator := bson.M{
        "$jsonSchema": bson.M{
            "bsonType": "object",
            "required": []string{"name", "email", "age"},
            "properties": bson.M{
                "name": bson.M{
                    "bsonType": "string",
                    "description": "Nombre del usuario",
                    "minLength": 2,
                    "maxLength": 100,
                },
                "email": bson.M{
                    "bsonType": "string",
                    "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
                    "description": "Email válido",
                },
                "age": bson.M{
                    "bsonType": "int",
                    "minimum": 0,
                    "maximum": 150,
                    "description": "Edad entre 0 y 150",
                },
                "tags": bson.M{
                    "bsonType": "array",
                    "items": bson.M{
                        "bsonType": "string",
                    },
                    "maxItems": 5,
                },
                "address": bson.M{
                    "bsonType": "object",
                    "properties": bson.M{
                        "street": bson.M{"bsonType": "string"},
                        "city": bson.M{"bsonType": "string"},
                        "zipcode": bson.M{"bsonType": "string"},
                    },
                },
            },
        },
    }
    
    opts := options.CreateCollection().SetValidator(validator)
    return db.CreateCollection(ctx, "usuarios", opts)
}

// Validación automática en inserciones
func TestValidation(db *mongo.Database, ctx context.Context) {
    collection := db.Collection("usuarios")
    
    // ✅ Válido
    collection.InsertOne(ctx, bson.M{
        "name": "Juan",
        "email": "juan@example.com",
        "age": 28,
    })
    
    // ❌ Inválido: falta campo required
    _, err := collection.InsertOne(ctx, bson.M{
        "name": "María",
        "age": 25,
        // falta email
    })
    // Error: Documento no cumple esquema
    
    // ❌ Inválido: age fuera de rango
    _, err = collection.InsertOne(ctx, bson.M{
        "name": "Carlos",
        "email": "carlos@example.com",
        "age": 200,
    })
    // Error: age > 150
}
```

### 65.10.2 Validación a Nivel de Aplicación

```go
package main

import (
    "errors"
    "regexp"
    "time"
)

type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    Name  string             `bson:"name"`
    Email string             `bson:"email"`
    Age   int                `bson:"age"`
    Created time.Time        `bson:"created_at,omitempty"`
}

// Validador personalizado
func (u *User) Validate() error {
    if u.Name == "" {
        return errors.New("nombre es requerido")
    }
    
    if len(u.Name) < 2 || len(u.Name) > 100 {
        return errors.New("nombre debe tener 2-100 caracteres")
    }
    
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(u.Email) {
        return errors.New("email inválido")
    }
    
    if u.Age < 0 || u.Age > 150 {
        return errors.New("edad debe estar entre 0 y 150")
    }
    
    return nil
}

// Uso en inserción
func (db *Database) InsertUserWithValidation(ctx context.Context, user User) error {
    // Validar antes de insertar
    if err := user.Validate(); err != nil {
        return err
    }
    
    user.Created = time.Now()
    
    collection := db.DB.Collection("usuarios")
    _, err := collection.InsertOne(ctx, user)
    return err
}
```

### 65.10.3 Validación de Datos Únicos

```go
func (db *Database) CreateUniqueIndexOnEmail(ctx context.Context) error {
    collection := db.DB.Collection("usuarios")
    
    indexModel := mongo.IndexModel{
        Keys: bson.D{{Key: "email", Value: 1}},
        Options: options.Index().
            SetUnique(true).
            SetSparse(true),  // Permite múltiples null
    }
    
    _, err := collection.Indexes().CreateOne(ctx, indexModel)
    return err
}

// Manejo de violación de unicidad
func (db *Database) InsertUserWithUniqueCheck(ctx context.Context, user User) (string, error) {
    collection := db.DB.Collection("usuarios")
    
    user.Created = time.Now()
    result, err := collection.InsertOne(ctx, user)
    
    if err != nil {
        // Verificar si es error de duplicado
        writeErr, ok := err.(mongo.WriteError)
        if ok && writeErr.Code == 11000 {
            return "", errors.New("email ya registrado")
        }
        return "", err
    }
    
    return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
```

---

## 65.11 Producción y Case Studies

### 65.11.1 Patrones de Manejo de Errores

```go
package main

import (
    "errors"
    "go.mongodb.org/mongo-driver/mongo"
)

// Tipos de errores MongoDB
func HandleMongoError(err error) string {
    if err == nil {
        return "Sin error"
    }
    
    // Error de conexión
    if mongo.IsNetworkError(err) {
        return "Error de red: intente nuevamente"
    }
    
    // Error de selección de servidor
    if mongo.IsServerSelectionError(err) {
        return "Servidor MongoDB no disponible"
    }
    
    // No encontrado
    if err == mongo.ErrNoDocuments {
        return "Documento no encontrado"
    }
    
    // Timeout
    if errors.Is(err, context.DeadlineExceeded) {
        return "Timeout en operación"
    }
    
    // Error de escritura
    writeErr, ok := err.(mongo.WriteError)
    if ok {
        if writeErr.Code == 11000 {
            return "Violación de constraint único"
        }
        return fmt.Sprintf("Error de escritura: %s", writeErr.Message)
    }
    
    // Bulk write error
    bulkErr, ok := err.(mongo.BulkWriteError)
    if ok {
        fmt.Printf("Errores en bulk: %v\n", bulkErr.WriteErrors)
        return "Algunos documentos fallaron"
    }
    
    return err.Error()
}

// Retry con backoff exponencial
func RetryOperation(operation func() error, maxRetries int) error {
    var lastErr error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Solo reintentar errores transitorios
        if !isTransientError(err) {
            return err
        }
        
        // Backoff exponencial
        backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
        time.Sleep(backoff)
    }
    
    return lastErr
}

func isTransientError(err error) bool {
    return mongo.IsNetworkError(err) || 
           mongo.IsServerSelectionError(err) ||
           errors.Is(err, context.DeadlineExceeded)
}
```

### 65.11.2 Logging y Monitoreo

```go
package main

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo/oplog"
)

type DatabaseMonitor struct {
    Logger Logger
}

func (m *DatabaseMonitor) SetupMonitoring(client *mongo.Client) {
    // Comando monitor
    cmdMonitor := &event.CommandMonitor{
        Started: func(_ context.Context, evt *event.CommandStartedEvent) {
            m.Logger.Debugf("Comando iniciado: %s (ns: %s)", evt.CommandName, evt.Namespace)
        },
        Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
            m.Logger.Debugf("Comando exitoso: %s (duración: %dns)", 
                evt.CommandName, evt.DurationNanos)
        },
        Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
            m.Logger.Warnf("Comando fallido: %s (razón: %s)", 
                evt.CommandName, evt.Failure)
        },
    }
    
    opts := options.Client().
        ApplyURI("mongodb://localhost:27017").
        SetMonitor(cmdMonitor)
    
    // Monitoring pool
    poolMonitor := &event.PoolMonitor{
        Event: func(evt *event.PoolEvent) {
            switch evt := evt.(type) {
            case *event.PoolCreatedEvent:
                m.Logger.Infof("Pool creado con tamaño máximo: %d", evt.Options.MaxPoolSize)
            case *event.ConnectionCheckedOutEvent:
                m.Logger.Debugf("Conexión reservada")
            case *event.ConnectionErrorOccurredEvent:
                m.Logger.Warnf("Error en conexión: %v", evt.Reason)
            }
        },
    }
}

// Health check
func (db *Database) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    if err := db.Client.Ping(ctx, nil); err != nil {
        return err
    }
    
    // Verificar colecciones críticas
    collections, err := db.DB.ListCollectionNames(ctx, bson.M{})
    if err != nil {
        return err
    }
    
    expectedCollections := []string{"usuarios", "ordenes", "productos"}
    for _, expected := range expectedCollections {
        found := false
        for _, col := range collections {
            if col == expected {
                found = true
                break
            }
        }
        if !found {
            return fmt.Errorf("colección faltante: %s", expected)
        }
    }
    
    return nil
}
```

### 65.11.3 Migrando de SQL a MongoDB

```go
// Estructura SQL → Document MongoDB

/*
SQL:
┌─────────────┐      ┌──────────────┐
│   users     │      │    posts     │
├─────────────┤      ├──────────────┤
│ id (PK)     │      │ id (PK)      │
│ name        │      │ title        │
│ email       │      │ content      │
└─────────────┘      │ user_id (FK) │
                     └──────────────┘

MongoDB (Denormalized):
{
  _id: ObjectId(...),
  name: "John",
  email: "john@example.com",
  posts: [
    { title: "Post 1", content: "...", created_at: ... },
    { title: "Post 2", content: "...", created_at: ... }
  ]
}

O Referencias (Normalized):
users: { _id: ..., name: "John", email: "..." }
posts: { _id: ..., title: "Post 1", user_id: ObjectId(...) }
*/

func MigrateUserData(sqlDB *sql.DB, mongodb *mongo.Database, ctx context.Context) error {
    // Leer de SQL
    rows, err := sqlDB.Query("SELECT id, name, email FROM users")
    if err != nil {
        return err
    }
    defer rows.Close()
    
    collection := mongodb.Collection("usuarios")
    var docs []interface{}
    
    for rows.Next() {
        var id int
        var name, email string
        
        if err := rows.Scan(&id, &name, &email); err != nil {
            return err
        }
        
        doc := bson.M{
            "sql_id": id,  // Guardar ID original por si acaso
            "name":   name,
            "email":  email,
        }
        docs = append(docs, doc)
    }
    
    // Insertar en MongoDB
    _, err = collection.InsertMany(ctx, docs)
    return err
}
```

### 65.11.4 Caso de Estudio: Aplicación E-commerce

```go
// Estructura de datos optimizada para e-commerce

type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    SKU         string             `bson:"sku"`
    Name        string             `bson:"name"`
    Description string             `bson:"description"`
    Price       float64            `bson:"price"`
    Stock       int                `bson:"stock"`
    Category    string             `bson:"category"`
    Images      []string           `bson:"images"`
    Reviews     []Review           `bson:"reviews"`
    CreatedAt   time.Time          `bson:"created_at"`
}

type Review struct {
    UserID    string    `bson:"user_id"`
    Rating    int       `bson:"rating"`     // 1-5
    Comment   string    `bson:"comment"`
    CreatedAt time.Time `bson:"created_at"`
}

type Order struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    UserID    string             `bson:"user_id"`
    Items     []OrderItem        `bson:"items"`
    Total     float64            `bson:"total"`
    Status    string             `bson:"status"`  // pending, paid, shipped, delivered
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

type OrderItem struct {
    ProductID   string  `bson:"product_id"`
    Name        string  `bson:"name"`
    Quantity    int     `bson:"quantity"`
    UnitPrice   float64 `bson:"unit_price"`
    Subtotal    float64 `bson:"subtotal"`
}

// Queries optimizadas
func (db *Database) GetProductsByCategory(ctx context.Context, category string) ([]Product, error) {
    collection := db.DB.Collection("productos")
    
    // Crear índice si no existe
    collection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{{Key: "category", Value: 1}},
    })
    
    cursor, err := collection.Find(ctx, bson.M{
        "category": category,
        "stock":    bson.M{"$gt": 0},
    })
    if err != nil {
        return nil, err
    }
    
    var products []Product
    cursor.All(ctx, &products)
    return products, nil
}

// Análisis: Productos más vendidos
func (db *Database) TopSellingProducts(ctx context.Context, limit int) ([]bson.M, error) {
    collection := db.DB.Collection("ordenes")
    
    pipeline := mongo.Pipeline{
        bson.D{{Key: "$unwind", Value: "$items"}},
        bson.D{
            {Key: "$group", Value: bson.D{
                {Key: "_id", Value: "$items.product_id"},
                {Key: "total_vendido", Value: bson.D{{Key: "$sum", Value: "$items.quantity"}}},
                {Key: "ingresos", Value: bson.D{{Key: "$sum", Value: "$items.subtotal"}}},
            }},
        },
        bson.D{{Key: "$sort", Value: bson.D{{Key: "total_vendido", Value: -1}}}},
        bson.D{{Key: "$limit", Value: int64(limit)}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    
    var results []bson.M
    cursor.All(ctx, &results)
    return results, nil
}
```

### 65.11.5 Backup y Recuperación

```go
// Backup con mongodump
func BackupDatabase(dbName string, outputPath string) error {
    cmd := exec.Command("mongodump",
        "--db", dbName,
        "--out", outputPath,
        "--uri", "mongodb://localhost:27017",
    )
    
    return cmd.Run()
}

// Restore con mongorestore
func RestoreDatabase(inputPath string, dbName string) error {
    cmd := exec.Command("mongorestore",
        "--uri", "mongodb://localhost:27017",
        "--db", dbName,
        inputPath,
    )
    
    return cmd.Run()
}

// Exportar a JSON (para análisis)
func ExportToJSON(dbName, collName, outputFile string) error {
    cmd := exec.Command("mongoexport",
        "--uri", "mongodb://localhost:27017",
        "--db", dbName,
        "--collection", collName,
        "--out", outputFile,
        "--pretty",
    )
    
    return cmd.Run()
}

// Importar desde JSON
func ImportFromJSON(dbName, collName, inputFile string) error {
    cmd := exec.Command("mongoimport",
        "--uri", "mongodb://localhost:27017",
        "--db", dbName,
        "--collection", collName,
        "--file", inputFile,
        "--jsonArray",  // Si es array de objetos
    )
    
    return cmd.Run()
}
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: CRUD Básico - Gestión de Usuarios

```go
package main

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    Name  string             `bson:"name"`
    Email string             `bson:"email"`
    Age   int                `bson:"age"`
}

func Exercise1() {
    client, _ := mongo.Connect(context.Background(), 
        options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(context.Background())
    
    collection := client.Database("ejercicios").Collection("usuarios")
    
    // 1. Insertar usuarios
    users := []interface{}{
        User{Name: "Ana", Email: "ana@example.com", Age: 28},
        User{Name: "Bob", Email: "bob@example.com", Age: 32},
        User{Name: "Carlos", Email: "carlos@example.com", Age: 25},
    }
    
    result, _ := collection.InsertMany(context.Background(), users)
    fmt.Printf("Insertados: %v\n", result.InsertedIDs)
    
    // 2. Buscar todos
    cursor, _ := collection.Find(context.Background(), bson.M{})
    var allUsers []User
    cursor.All(context.Background(), &allUsers)
    fmt.Printf("Total usuarios: %d\n", len(allUsers))
    
    // 3. Buscar por email
    var user User
    collection.FindOne(context.Background(), 
        bson.M{"email": "ana@example.com"}).Decode(&user)
    fmt.Printf("Encontrado: %s (edad: %d)\n", user.Name, user.Age)
    
    // 4. Actualizar edad
    collection.UpdateOne(context.Background(),
        bson.M{"email": "bob@example.com"},
        bson.M{"$set": bson.M{"age": 33}})
    fmt.Println("Bob actualizado a edad 33")
    
    // 5. Eliminar usuario
    collection.DeleteOne(context.Background(),
        bson.M{"email": "carlos@example.com"})
    fmt.Println("Carlos eliminado")
}
```

### Ejercicio 2: Queries Avanzadas y Filtros

```go
func Exercise2() {
    client, _ := mongo.Connect(context.Background(),
        options.Client().ApplyURI("mongodb://localhost:27017"))
    collection := client.Database("ejercicios").Collection("productos")
    
    // Crear documentos
    collection.InsertMany(context.Background(), []interface{}{
        bson.M{"nombre": "Laptop", "precio": 999, "stock": 5, "categoria": "IT"},
        bson.M{"nombre": "Mouse", "precio": 29, "stock": 50, "categoria": "IT"},
        bson.M{"nombre": "Monitor", "precio": 299, "stock": 10, "categoria": "IT"},
        bson.M{"nombre": "Teclado", "precio": 89, "stock": 20, "categoria": "IT"},
        bson.M{"nombre": "Silla", "precio": 249, "stock": 8, "categoria": "Oficina"},
    })
    
    // Queries
    fmt.Println("1. Productos con precio entre 50 y 300:")
    cursor, _ := collection.Find(context.Background(), bson.M{
        "precio": bson.M{"$gte": 50, "$lte": 300},
    })
    var productos []bson.M
    cursor.All(context.Background(), &productos)
    for _, p := range productos {
        fmt.Printf("  - %s: $%v\n", p["nombre"], p["precio"])
    }
    
    fmt.Println("\n2. Productos con stock bajo (< 15):")
    cursor, _ = collection.Find(context.Background(), bson.M{
        "stock": bson.M{"$lt": 15},
    })
    cursor.All(context.Background(), &productos)
    for _, p := range productos {
        fmt.Printf("  - %s: %d unidades\n", p["nombre"], p["stock"])
    }
    
    fmt.Println("\n3. Productos de IT ordenados por precio (descendente):")
    opts := options.Find().SetSort(bson.M{"precio": -1})
    cursor, _ = collection.Find(context.Background(), 
        bson.M{"categoria": "IT"}, opts)
    cursor.All(context.Background(), &productos)
    for _, p := range productos {
        fmt.Printf("  - %s: $%v\n", p["nombre"], p["precio"])
    }
}
```

### Ejercicio 3: Agregación y Pipeline

```go
func Exercise3() {
    client, _ := mongo.Connect(context.Background(),
        options.Client().ApplyURI("mongodb://localhost:27017"))
    collection := client.Database("ejercicios").Collection("ventas")
    
    // Datos de ejemplo
    collection.InsertMany(context.Background(), []interface{}{
        bson.M{"producto": "A", "cantidad": 10, "precio": 50, "mes": "enero"},
        bson.M{"producto": "B", "cantidad": 5, "precio": 100, "mes": "enero"},
        bson.M{"producto": "A", "cantidad": 15, "precio": 50, "mes": "febrero"},
        bson.M{"producto": "C", "cantidad": 20, "precio": 25, "mes": "febrero"},
        bson.M{"producto": "B", "cantidad": 8, "precio": 100, "mes": "marzo"},
    })
    
    // Agregación: Ventas totales por producto
    pipeline := mongo.Pipeline{
        bson.D{{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$producto"},
            {Key: "cantidad_total", Value: bson.D{{Key: "$sum", Value: "$cantidad"}}},
            {Key: "ingresos", Value: bson.D{
                {Key: "$sum", Value: bson.D{
                    {Key: "$multiply", Value: bson.A{"$cantidad", "$precio"}},
                }},
            }},
        }}},
        bson.D{{Key: "$sort", Value: bson.D{{Key: "ingresos", Value: -1}}}},
    }
    
    cursor, _ := collection.Aggregate(context.Background(), pipeline)
    var resultados []bson.M
    cursor.All(context.Background(), &resultados)
    
    fmt.Println("Ventas por producto:")
    for _, r := range resultados {
        fmt.Printf("%s: %d unidades, $%.2f en ingresos\n",
            r["_id"], r["cantidad_total"], r["ingresos"])
    }
}
```

### Ejercicio 4: Transacciones

```go
func Exercise4() {
    client, _ := mongo.Connect(context.Background(),
        options.Client().ApplyURI("mongodb://localhost:27017"))
    db := client.Database("ejercicios")
    
    // Crear colecciones
    db.Collection("cuentas").InsertMany(context.Background(), []interface{}{
        bson.M{"_id": "cuenta1", "saldo": 1000},
        bson.M{"_id": "cuenta2", "saldo": 500},
    })
    
    // Transacción: transferencia
    session, _ := client.StartSession()
    defer session.EndSession(context.Background())
    
    err := mongo.WithSession(context.Background(), session, func(ctx mongo.SessionContext) error {
        session.StartTransaction()
        
        cuentas := db.Collection("cuentas")
        
        // Restar de cuenta1
        _, err := cuentas.UpdateOne(ctx, 
            bson.M{"_id": "cuenta1"},
            bson.M{"$inc": bson.M{"saldo": -200}})
        if err != nil {
            session.AbortTransaction(ctx)
            return err
        }
        
        // Sumar a cuenta2
        _, err = cuentas.UpdateOne(ctx,
            bson.M{"_id": "cuenta2"},
            bson.M{"$inc": bson.M{"saldo": 200}})
        if err != nil {
            session.AbortTransaction(ctx)
            return err
        }
        
        return session.CommitTransaction(ctx)
    })
    
    if err != nil {
        fmt.Printf("Error en transacción: %v\n", err)
    } else {
        fmt.Println("Transferencia exitosa")
        
        // Verificar saldos
        var cuenta1, cuenta2 bson.M
        db.Collection("cuentas").FindOne(context.Background(), 
            bson.M{"_id": "cuenta1"}).Decode(&cuenta1)
        db.Collection("cuentas").FindOne(context.Background(),
            bson.M{"_id": "cuenta2"}).Decode(&cuenta2)
        
        fmt.Printf("Cuenta 1: $%v\n", cuenta1["saldo"])
        fmt.Printf("Cuenta 2: $%v\n", cuenta2["saldo"])
    }
}
```

### Ejercicio 5: Aplicación Completa - Blog

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type BlogPost struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Title     string             `bson:"title"`
    Content   string             `bson:"content"`
    Author    string             `bson:"author"`
    Tags      []string           `bson:"tags"`
    Views     int                `bson:"views"`
    CreatedAt time.Time          `bson:"created_at"`
    Comments  []Comment          `bson:"comments"`
}

type Comment struct {
    Author  string    `bson:"author"`
    Text    string    `bson:"text"`
    Rating  int       `bson:"rating"`
    Created time.Time `bson:"created_at"`
}

type BlogDB struct {
    db *mongo.Database
}

func Exercise5() {
    client, _ := mongo.Connect(context.Background(),
        options.Client().ApplyURI("mongodb://localhost:27017"))
    blog := &BlogDB{db: client.Database("blog")}
    
    // Crear índices
    blog.db.Collection("posts").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
        {Keys: bson.D{{Key: "title", Value: "text"}}},
        {Keys: bson.D{{Key: "tags", Value: 1}}},
        {Keys: bson.D{{Key: "created_at", Value: -1}}},
    })
    
    // Insertar posts
    posts := []interface{}{
        BlogPost{
            Title: "MongoDB Basics",
            Content: "Introducción a MongoDB...",
            Author: "Juan",
            Tags: []string{"mongodb", "database", "nosql"},
            Views: 150,
            CreatedAt: time.Now().AddDate(0, 0, -5),
            Comments: []Comment{
                {Author: "Ana", Text: "Excelente post!", Rating: 5, Created: time.Now().AddDate(0, 0, -3)},
            },
        },
        BlogPost{
            Title: "Go with MongoDB",
            Content: "Usando driver de MongoDB en Go...",
            Author: "María",
            Tags: []string{"go", "mongodb", "programming"},
            Views: 200,
            CreatedAt: time.Now().AddDate(0, 0, -2),
        },
    }
    
    result, _ := blog.db.Collection("posts").InsertMany(context.Background(), posts)
    fmt.Printf("Posts insertados: %d\n", len(result.InsertedIDs))
    
    // Buscar posts por tag
    fmt.Println("\nPosts con tag 'mongodb':")
    cursor, _ := blog.db.Collection("posts").Find(context.Background(),
        bson.M{"tags": "mongodb"})
    var foundPosts []BlogPost
    cursor.All(context.Background(), &foundPosts)
    for _, p := range foundPosts {
        fmt.Printf("- %s por %s (%d views)\n", p.Title, p.Author, p.Views)
    }
    
    // Posts más populares
    fmt.Println("\nPosts ordenados por popularidad:")
    opts := options.Find().SetSort(bson.M{"views": -1})
    cursor, _ = blog.db.Collection("posts").Find(context.Background(), bson.M{}, opts)
    cursor.All(context.Background(), &foundPosts)
    for _, p := range foundPosts {
        fmt.Printf("- %s (%d views)\n", p.Title, p.Views)
    }
    
    // Agregar comentario
    fmt.Println("\nAñadiendo comentario...")
    blog.db.Collection("posts").UpdateOne(context.Background(),
        bson.M{"title": "Go with MongoDB"},
        bson.M{"$push": bson.M{"comments": Comment{
            Author: "Carlos",
            Text: "Muy útil, gracias!",
            Rating: 4,
            Created: time.Now(),
        }}})
    
    // Estadísticas
    fmt.Println("\nEstadísticas:")
    pipeline := mongo.Pipeline{
        bson.D{{Key: "$group", Value: bson.D{
            {Key: "_id", Value: nil},
            {Key: "total_posts", Value: bson.D{{Key: "$sum", Value: 1}}},
            {Key: "promedio_views", Value: bson.D{{Key: "$avg", Value: "$views"}}},
            {Key: "total_comentarios", Value: bson.D{
                {Key: "$sum", Value: bson.D{{Key: "$size", Value: "$comments"}}},
            }},
        }}},
    }
    
    cursor, _ = blog.db.Collection("posts").Aggregate(context.Background(), pipeline)
    var stats []bson.M
    cursor.All(context.Background(), &stats)
    if len(stats) > 0 {
        s := stats[0]
        fmt.Printf("Total posts: %v\n", s["total_posts"])
        fmt.Printf("Promedio views: %.1f\n", s["promedio_views"])
        fmt.Printf("Total comentarios: %v\n", s["total_comentarios"])
    }
}
```

---

## COMPARACIONES IMPORTANTES

### MongoDB vs SQL (PostgreSQL)

```
┌──────────────────────────┬─────────────────┬──────────────────┐
│ Aspecto                  │ PostgreSQL (SQL)│ MongoDB (NoSQL)   │
├──────────────────────────┼─────────────────┼──────────────────┤
│ Modelo de datos          │ Tablas + Filas  │ Documentos JSON   │
│ Esquema                  │ Rígido          │ Flexible          │
│ Relaciones               │ Foreign Keys    │ References/Embed  │
│ Queries múltiples tablas │ JOINs           │ Aggregation       │
│ ACID                     │ Garantizados    │ Multi-doc (nueva) │
│ Escalabilidad            │ Vertical        │ Horizontal        │
│ Índices                  │ B-tree          │ B-tree/Hashed     │
│ Duplicación              │ Normalizado     │ Denormalizado     │
│ Replicación              │ Configurado     │ Integrado         │
│ Cuota aprox.             │ $400/año        │ Gratuito (OSS)    │
└──────────────────────────┴─────────────────┴──────────────────┘
```

### MongoDB vs Redis

```
┌─────────────────┬─────────────────────┬──────────────────────┐
│ Característica  │ Redis               │ MongoDB              │
├─────────────────┼─────────────────────┼──────────────────────┤
│ Tipo            │ In-memory cache     │ Document database    │
│ Persistencia    │ Opcional (RDB/AOF)  │ Persistente (default)│
│ Estructura      │ Key-value           │ Documentos           │
│ Tamaño datos    │ Todo en RAM         │ Disco (+ cache)      │
│ Complejidad     │ Estructuras simples │ Documentos complejos │
│ Queries         │ Clave exacta        │ Filtros avanzados    │
│ Velocidad       │ Microsegundos       │ Milisegundos         │
│ Transacciones   │ Limited (Lua)       │ Multi-document ACID  │
│ Replicación     │ Master-slave        │ Replica Set          │
└─────────────────┴─────────────────────┴──────────────────────┘
```

---

## CONCLUSIONES Y MEJORES PRÁCTICAS

### ✅ Mejores Prácticas

1. **Índices Estratégicos:** Crear solo en campos consultados frecuentemente
2. **Denormalización Moderada:** Balancear entre performance y consistencia
3. **Validación Dual:** MongoDB + Aplicación
4. **Monitoreo:** Usar explain() y logs de comando
5. **Transacciones:** Para operaciones críticas multi-documento
6. **Connection Pooling:** Configurar min/max pool size
7. **Error Handling:** Implementar retry logic para fallos transitorios
8. **Backup Regular:** Usar mongodump o Atlas backups
9. **Escalabilidad:** Planificar sharding desde el inicio
10. **Testing:** Usar testcontainers para MongoDB en tests

### ❌ Anti-patrones

1. **Sin índices:** Causa collection scans lentos
2. **Documentos enormes:** > 16MB causa problemas
3. **Demasiados índices:** Ralentiza inserciones
4. **Sin validación:** Datos inconsistentes
5. **Queries sin límites:** Puede causar OOM
6. **Connections sin pool:** Agotamiento de recursos
7. **Ignorar sharding:** Problemas a escala
8. **Sin manejo de errores:** Fallos silenciosos

---

**FIN DEL CAPÍTULO 65**

Tamaño total: ~2,050 líneas | ~46 KB | 70% teoría + 30% código

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/65-mongodb-y-nosql/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/65-mongodb-y-nosql):

```bash
cd examples/65-mongodb-y-nosql
go run .
```
