# Capítulo 39: Database y SQL

## Introducción

El manejo de bases de datos es uno de los aspectos más críticos de cualquier aplicación moderna. Go proporciona el package `database/sql`, una abstracción elegante y poderosa que permite trabajar con cualquier base de datos relacional de manera uniforme. A diferencia de lenguajes como Python con SQLAlchemy o Java con JDBC, Go mantiene un enfoque minimalista pero completo que se alinea con su filosofía: ser simple, rápido y eficiente.

Este capítulo explora cómo Go maneja la persistencia de datos, desde conexiones básicas hasta transacciones complejas, pooling de conexiones y patrones de producción. Aprenderemos por qué el enfoque de Go a las bases de datos es diferente y por qué es tan efectivo en aplicaciones de alta escala.

---

## 39.1 El Package database/sql

### 39.1.1 Concepto de Abstracción de Base de Datos

El package `database/sql` en Go implementa una interfaz unificada para trabajar con múltiples sistemas de bases de datos. Su arquitectura es revolucionaria:

```
┌──────────────────────────────────────────────────────────┐
│          Código de Aplicación Go                         │
├──────────────────────────────────────────────────────────┤
│  database/sql (interfaz unificada del estándar)         │
│  ┌──────────────────────────────────────────────────┐   │
│  │ *sql.DB, *sql.Row, *sql.Rows, *sql.Tx, *sql.Stmt│   │
│  └──────────────────────────────────────────────────┘   │
├──────────────────────────────────────────────────────────┤
│        sql.Driver (implementación específica)            │
│  ┌──────────────────────────────────────────────────┐   │
│  │ github.com/lib/pq (PostgreSQL)                  │   │
│  │ github.com/go-sql-driver/mysql (MySQL)          │   │
│  │ modernc.org/sqlite (SQLite)                     │   │
│  │ github.com/jackc/pgx (Postgres - avanzado)      │   │
│  └──────────────────────────────────────────────────┘   │
├──────────────────────────────────────────────────────────┤
│     Base de Datos (PostgreSQL, MySQL, SQLite, etc.)     │
└──────────────────────────────────────────────────────────┘
```

### 39.1.2 Principios de Diseño

El design de `database/sql` sigue estos principios:

1. **Stateless**: Las conexiones son manejadas internamente por un pool
2. **Type-safe**: Golang es tipado estáticamente, los drivers deben respetar esto
3. **Minimalist**: Solo lo necesario, nada más
4. **Concurrent**: Diseñado para múltiples goroutines
5. **Portable**: El mismo código funciona con diferentes BD

### 39.1.3 Componentes Principales

```go
// DB es el punto de entrada principal - representa el pool de conexiones
type DB struct {
    // Campos privados
}

// Row representa una fila de un query
type Row struct {
    // Campos privados
}

// Rows representa múltiples filas de un query
type Rows struct {
    // Campos privados
}

// Tx representa una transacción
type Tx struct {
    // Campos privados
}

// Stmt representa una prepared statement
type Stmt struct {
    // Campos privados
}

// Driver es la interfaz que deben implementar los drivers
type Driver interface {
    Open(name string) (Conn, error)
}
```

### 39.1.4 Drivers Disponibles

| Driver | Paquete | BD | Características |
|--------|---------|-----|-----------------|
| MySQL | github.com/go-sql-driver/mysql | MySQL 5.7+ | Ligero, confiable |
| PostgreSQL | github.com/lib/pq | PostgreSQL 9.1+ | Maduro, completo |
| SQLite | modernc.org/sqlite | SQLite 3 | Sin dependencias C |
| SQL Server | github.com/denisenkom/go-mssqldb | SQL Server | Soporte completo |
| Oracle | github.com/sijms/go-ora | Oracle DB | Compatible con JDBC |
| PostgreSQL (pgx) | github.com/jackc/pgx | PostgreSQL 9.1+ | Alto rendimiento |

### 39.1.5 Comparación: Go sql vs SQLAlchemy (Python) vs JDBC (Java)

```go
// ╔═══════════════════════════════════════════════════════════╗
// ║              Go (Eficiente y Explícito)                  ║
// ╚═══════════════════════════════════════════════════════════╝
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
)

db, err := sql.Open("postgres", "dbname=mydb user=postgres password=secret")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

var name string
err = db.QueryRow("SELECT name FROM users WHERE id=$1", 1).Scan(&name)
if err != nil {
    log.Fatal(err)
}
fmt.Println(name)

// ╔═══════════════════════════════════════════════════════════╗
// ║         SQLAlchemy (Python - ORM abstracto)              ║
// ╚═══════════════════════════════════════════════════════════╝
/*
from sqlalchemy import create_engine, Column, Integer, String
from sqlalchemy.orm import sessionmaker

engine = create_engine("postgresql://user:password@localhost/mydb")
Session = sessionmaker(bind=engine)
session = Session()

user = session.query(User).filter(User.id == 1).first()
print(user.name)
session.close()
*/

// ╔═══════════════════════════════════════════════════════════╗
// ║         JDBC (Java - Verbose pero explícito)             ║
// ╚═══════════════════════════════════════════════════════════╝
/*
Class.forName("org.postgresql.Driver");
Connection conn = DriverManager.getConnection(
    "jdbc:postgresql://localhost:5432/mydb", 
    "user", 
    "password"
);

PreparedStatement pstmt = conn.prepareStatement(
    "SELECT name FROM users WHERE id=?"
);
pstmt.setInt(1, 1);
ResultSet rs = pstmt.executeQuery();
if (rs.next()) {
    String name = rs.getString("name");
    System.out.println(name);
}
rs.close();
pstmt.close();
conn.close();
*/
```

**Análisis comparativo:**

- **Go**: Balanceado entre simplicidad y control. No es un ORM completo, pero es transparente
- **SQLAlchemy**: Alto nivel de abstracción, pero puede ser lento en queries complejas
- **JDBC**: Muy verbose pero muy explícito, requiere gestión manual de recursos

---

## 39.2 Conectar a Bases de Datos

### 39.2.1 La función sql.Open()

```go
package main

import (
    "database/sql"
    _ "github.com/lib/pq"  // underscore: registra el driver
    "log"
)

func main() {
    // sql.Open NOT abre conexión, solo crea el pool
    db, err := sql.Open("postgres", "user=admin password=secret dbname=mydb sslmode=disable")
    if err != nil {
        log.Fatal("Error abriendo BD:", err)
    }
    defer db.Close()
    
    // Aquí ES cuando abrimos la conexión (lazy connect)
    if err := db.Ping(); err != nil {
        log.Fatal("Error conectando a BD:", err)
    }
    
    log.Println("Conectado exitosamente")
}
```

**Punto crítico**: `sql.Open()` NO abre una conexión inmediatamente. Solo crea el pool. La conexión real se abre cuando se ejecuta el primer query o cuando se llama a `Ping()`.

### 39.2.2 Connection Strings por Driver

```go
// PostgreSQL (lib/pq)
dsn := "user=postgres password=secret host=localhost port=5432 dbname=mydb sslmode=disable"
db, _ := sql.Open("postgres", dsn)

// MySQL (go-sql-driver/mysql)
dsn := "username:password@tcp(localhost:3306)/dbname"
db, _ := sql.Open("mysql", dsn)

// SQLite (sqlite3)
dsn := "file:mydb.db?cache=shared&mode=rwc"
db, _ := sql.Open("sqlite", dsn)

// SQL Server (go-mssqldb)
dsn := "server=localhost;user id=sa;password=secret;database=mydb"
db, _ := sql.Open("mssql", dsn)
```

### 39.2.3 Configuración de Pool de Conexiones

```go
package main

import (
    "database/sql"
    "time"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "...")
    defer db.Close()
    
    // Configuración del pool (ANTES de usar la BD)
    db.SetMaxOpenConns(25)      // máximo de conexiones abiertas
    db.SetMaxIdleConns(5)       // conexiones inactivas en el pool
    db.SetConnMaxLifetime(5 * time.Minute)    // vida máxima de una conexión
    db.SetConnMaxIdleTime(10 * time.Minute)   // tiempo máximo inactiva
    
    // Verificar que conecta
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
}
```

### 39.2.4 Manejo de Errores en Conexión

```go
package main

import (
    "database/sql"
    "errors"
    "log"
    "time"
)

func connectWithRetry(dsn string, maxRetries int) (*sql.DB, error) {
    var db *sql.DB
    var err error
    
    for i := 0; i < maxRetries; i++ {
        db, err = sql.Open("postgres", dsn)
        if err == nil {
            if err = db.Ping(); err == nil {
                return db, nil
            }
        }
        
        log.Printf("Intento %d falló: %v", i+1, err)
        time.Sleep(time.Second * time.Duration(i+1)) // backoff exponencial
    }
    
    return nil, errors.New("no se pudo conectar después de reintentos")
}

func main() {
    db, err := connectWithRetry("postgresql://...", 3)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
}
```

---

## 39.3 Ejecutar Queries

### 39.3.1 Query, QueryRow y Scan

```go
package main

import (
    "database/sql"
    "log"
)

type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

func main() {
    db, _ := sql.Open("sqlite", "users.db")
    defer db.Close()
    
    // Caso 1: UNA fila (QueryRow + Scan)
    var user User
    err := db.QueryRow(
        "SELECT id, name, email, age FROM users WHERE id = ?", 
        1,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    
    if err == sql.ErrNoRows {
        log.Println("Usuario no encontrado")
    } else if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Usuario: %+v\n", user)
    }
    
    // Caso 2: MÚLTIPLES filas (Query + Rows)
    rows, err := db.Query(
        "SELECT id, name, email, age FROM users WHERE age > ? ORDER BY name",
        18,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()  // IMPORTANTE: siempre cerrar rows
    
    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age); err != nil {
            log.Fatal(err)
        }
        users = append(users, u)
    }
    
    // Verificar errores de iteración
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Encontrados %d usuarios\n", len(users))
}
```

### 39.3.2 Scanning Básico

```go
package main

import "database/sql"

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Scan a variables simples
    var id int
    var name string
    var balance float64
    
    err := db.QueryRow(
        "SELECT id, name, balance FROM accounts WHERE id = ?", 1,
    ).Scan(&id, &name, &balance)
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Los tipos deben coincidir o ser convertibles
    // int/int64, string, float64/float32, []byte, time.Time son comunes
}
```

### 39.3.3 Patrones de Query Común

```go
package main

import (
    "database/sql"
    "log"
)

// Obtener un valor único
func getValue(db *sql.DB, query string, args ...interface{}) (interface{}, error) {
    var value interface{}
    err := db.QueryRow(query, args...).Scan(&value)
    return value, err
}

// Obtener una fila como map
func getRowAsMap(db *sql.DB, query string, args ...interface{}) (map[string]interface{}, error) {
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    if !rows.Next() {
        return nil, sql.ErrNoRows
    }
    
    values := make([]interface{}, len(columns))
    valuePtrs := make([]interface{}, len(columns))
    for i := range columns {
        valuePtrs[i] = &values[i]
    }
    
    if err := rows.Scan(valuePtrs...); err != nil {
        return nil, err
    }
    
    result := make(map[string]interface{})
    for i, col := range columns {
        result[col] = values[i]
    }
    return result, nil
}

// Obtener todas las filas como slice de maps
func getAllRowsAsMap(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    var results []map[string]interface{}
    for rows.Next() {
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range columns {
            valuePtrs[i] = &values[i]
        }
        
        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }
        
        entry := make(map[string]interface{})
        for i, col := range columns {
            entry[col] = values[i]
        }
        results = append(results, entry)
    }
    
    return results, rows.Err()
}

func main() {
    db, _ := sql.Open("sqlite", "test.db")
    defer db.Close()
    
    // Usar las funciones helper
    rowMap, _ := getRowAsMap(db, "SELECT * FROM users WHERE id = ?", 1)
    log.Printf("Row: %+v\n", rowMap)
    
    allRows, _ := getAllRowsAsMap(db, "SELECT * FROM users LIMIT 10")
    log.Printf("Rows: %+v\n", allRows)
}
```

---

## 39.4 Prepared Statements

### 39.4.1 ¿Por qué Prepared Statements?

Los prepared statements son fundamentales por dos razones críticas:

1. **Seguridad**: Previenen SQL injection
2. **Performance**: Se reutiliza el plan de ejecución

```go
package main

import (
    "database/sql"
    "log"
)

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // ❌ INCORRECTO: SQL injection vulnerable
    userInput := "'; DROP TABLE users; --"
    query := "SELECT * FROM users WHERE name = '" + userInput + "'"
    // Esto ejecutaría: SELECT * FROM users WHERE name = ''; DROP TABLE users; --'
    
    // ✅ CORRECTO: Prepared statement
    stmt, err := db.Prepare("SELECT id, name, email FROM users WHERE name = ?")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    
    // Los ? son placeholders seguros
    rows, err := stmt.Query(userInput)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
}
```

### 39.4.2 Crear y Reutilizar Prepared Statements

```go
package main

import (
    "database/sql"
    "log"
)

type User struct {
    ID   int
    Name string
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Preparar una vez
    stmt, err := db.Prepare("SELECT id, name FROM users WHERE id = ?")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    
    // Ejecutar múltiples veces (eficiente)
    for i := 1; i <= 100; i++ {
        var user User
        err := stmt.QueryRow(i).Scan(&user.ID, &user.Name)
        if err == sql.ErrNoRows {
            log.Printf("Usuario %d no encontrado", i)
        } else if err != nil {
            log.Fatal(err)
        } else {
            log.Printf("Usuario: %+v\n", user)
        }
    }
}
```

### 39.4.3 Prepared Statements en Transacciones

```go
package main

import (
    "database/sql"
    "log"
)

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // Prepared statement dentro de transacción
    stmt, err := tx.Prepare("INSERT INTO logs (user_id, action) VALUES (?, ?)")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    
    // Insertar múltiples registros
    actions := []string{"login", "view_page", "logout"}
    for _, action := range actions {
        _, err := stmt.Exec(1, action)
        if err != nil {
            log.Fatal(err)
        }
    }
    
    if err := tx.Commit(); err != nil {
        log.Fatal(err)
    }
}
```

### 39.4.4 Named Parameters (Context-based)

```go
package main

import (
    "context"
    "database/sql"
    "log"
)

func main() {
    db, _ := sql.Open("postgres", "...")
    defer db.Close()
    
    // PostgreSQL soporta $1, $2, etc.
    ctx := context.Background()
    
    rows, err := db.QueryContext(ctx,
        "SELECT id, name FROM users WHERE age > $1 AND city = $2",
        18, "Madrid",
    )
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            log.Fatal(err)
        }
        log.Printf("%d: %s\n", id, name)
    }
}
```

---

## 39.5 Inserción, Actualización y Borrado

### 39.5.1 Exec, LastInsertId y RowsAffected

```go
package main

import (
    "database/sql"
    "log"
)

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // INSERT
    result, err := db.Exec(
        "INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
        "Alice", "alice@example.com", 30,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Obtener el ID insertado (solo en SQLite, MySQL)
    id, err := result.LastInsertId()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Inserción exitosa, ID: %d\n", id)
    
    // UPDATE
    result, err = db.Exec(
        "UPDATE users SET age = ? WHERE id = ?",
        31, id,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Obtener número de filas afectadas
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Filas actualizadas: %d\n", rowsAffected)
    
    // DELETE
    result, err = db.Exec(
        "DELETE FROM users WHERE age < ?",
        18,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    rowsDeleted, _ := result.RowsAffected()
    log.Printf("Filas eliminadas: %d\n", rowsDeleted)
}
```

### 39.5.2 Bulk Insert (Inserción masiva)

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "strings"
)

type User struct {
    Name  string
    Email string
    Age   int
}

func bulkInsert(db *sql.DB, users []User) error {
    if len(users) == 0 {
        return nil
    }
    
    // Construir query dinámica
    placeholders := []string{}
    args := []interface{}{}
    
    for i, user := range users {
        placeholders = append(placeholders, "(?, ?, ?)")
        args = append(args, user.Name, user.Email, user.Age)
    }
    
    query := "INSERT INTO users (name, email, age) VALUES " + 
             strings.Join(placeholders, ", ")
    
    result, err := db.Exec(query, args...)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    log.Printf("Insertados %d usuarios\n", rows)
    return nil
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    users := []User{
        {"Alice", "alice@example.com", 28},
        {"Bob", "bob@example.com", 35},
        {"Carol", "carol@example.com", 42},
    }
    
    if err := bulkInsert(db, users); err != nil {
        log.Fatal(err)
    }
}
```

### 39.5.3 Actualización Condicional

```go
package main

import (
    "database/sql"
    "log"
)

// Actualizar solo si cumple condiciones
func updateIfExists(db *sql.DB, userID int, newEmail string) (bool, error) {
    result, err := db.Exec(
        "UPDATE users SET email = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
        newEmail, userID,
    )
    if err != nil {
        return false, err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return false, err
    }
    
    return rows > 0, nil
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    updated, err := updateIfExists(db, 1, "newemail@example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    if updated {
        log.Println("Email actualizado")
    } else {
        log.Println("Usuario no encontrado")
    }
}
```

---

## 39.6 Transacciones ACID

### 39.6.1 Fundamentos de Transacciones

Las transacciones garantizan ACID:

- **Atomicidad**: Todo o nada
- **Consistencia**: BD pasa de un estado válido a otro válido
- **Aislamiento**: Cambios concurrentes no se interfieren
- **Durabilidad**: Una vez confirmada, persiste

```go
package main

import (
    "database/sql"
    "log"
)

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Iniciar transacción
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    
    // TODO algo dentro de la transacción
    _, err = tx.Exec("INSERT INTO logs (message) VALUES (?)", "Action 1")
    if err != nil {
        // Error: hacer rollback
        tx.Rollback()
        log.Fatal(err)
    }
    
    _, err = tx.Exec("INSERT INTO logs (message) VALUES (?)", "Action 2")
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    
    // Todo bien: hacer commit
    if err := tx.Commit(); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Transacción completada")
}
```

### 39.6.2 Transferencia entre Cuentas (Caso de Uso Real)

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
)

type Account struct {
    ID      int
    Balance float64
}

func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    // Versión con defer para garantizar rollback en caso de panic
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // Step 1: Leer balance de cuenta origen
    var fromBalance float64
    row := tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromID)
    if err := row.Scan(&fromBalance); err != nil {
        return fmt.Errorf("cuenta origen no encontrada: %w", err)
    }
    
    // Step 2: Verificar fondos suficientes
    if fromBalance < amount {
        return fmt.Errorf("fondos insuficientes: %.2f < %.2f", fromBalance, amount)
    }
    
    // Step 3: Restar de cuenta origen
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - ? WHERE id = ?",
        amount, fromID,
    )
    if err != nil {
        return err
    }
    
    // Step 4: Sumar a cuenta destino
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance + ? WHERE id = ?",
        amount, toID,
    )
    if err != nil {
        return err
    }
    
    // Step 5: Registrar transacción
    _, err = tx.Exec(
        "INSERT INTO transactions (from_account, to_account, amount, status) VALUES (?, ?, ?, ?)",
        fromID, toID, amount, "completed",
    )
    if err != nil {
        return err
    }
    
    // Step 6: Commit
    if err := tx.Commit(); err != nil {
        return err
    }
    
    log.Printf("Transferencia exitosa: %.2f de cuenta %d a %d\n", amount, fromID, toID)
    return nil
}

func main() {
    db, _ := sql.Open("sqlite", "bank.db")
    defer db.Close()
    
    if err := transferMoney(db, 1, 2, 100.50); err != nil {
        log.Printf("Transferencia fallida: %v\n", err)
    }
}
```

### 39.6.3 Niveles de Aislamiento

```go
package main

import (
    "database/sql"
    "database/sql/driver"
    "log"
)

// Go no abstrae esto bien, depende del driver
// PostgreSQL soporta BEGIN ISOLATION LEVEL

func main() {
    db, _ := sql.Open("postgres", "...")
    defer db.Close()
    
    // Transacción con aislamiento explícito
    tx, _ := db.Begin()
    
    // Ejecutar comando SQL directamente
    tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
    
    // Ahora la transacción está en el nivel más alto de aislamiento
    
    tx.Commit()
}
```

---

## 39.7 Connection Pooling

### 39.7.1 ¿Cómo Funciona el Pool?

```
┌─────────────────────────────────────────────────────────┐
│              Connection Pool en sql.DB                  │
├─────────────────────────────────────────────────────────┤
│  MaxOpenConns = 25 (máximo de conexiones abiertas)     │
│  MaxIdleConns = 5  (máximo de conexiones inactivas)    │
│  ┌──────────┬──────────┬──────────┬──────────┐          │
│  │  Conn 1  │  Conn 2  │  Conn 3  │  Conn 4  │  (activas)
│  └──────────┴──────────┴──────────┴──────────┘          │
│  ┌──────────┬──────────┐                                 │
│  │ Conn 5   │ Conn 6   │ ... (inactivas, esperando)    │
│  └──────────┴──────────┘                                 │
│                                                          │
│  Solicitud nueva → Espera idle conn → Si no hay: +1     │
└─────────────────────────────────────────────────────────┘
```

### 39.7.2 Configuración del Pool

```go
package main

import (
    "database/sql"
    "time"
)

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Configuraciones recomendadas
    
    // 1. Para aplicación web típica (100 RPS)
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    
    // 2. Para aplicación de alta carga (10,000 RPS)
    db.SetMaxOpenConns(100)
    db.SetMaxIdleConns(20)
    
    // 3. Para batch processing
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(2)
    
    // Ciclo de vida de conexiones
    db.SetConnMaxLifetime(15 * time.Minute)    // Reciclar después de 15 min
    db.SetConnMaxIdleTime(5 * time.Minute)     // Cerrar si inactiva 5 min
}
```

### 39.7.3 Monitoreo del Pool

```go
package main

import (
    "database/sql"
    "log"
    "time"
)

func monitorPool(db *sql.DB, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := db.Stats()
        log.Printf("Pool Stats:\n"+
            "  Open Connections: %d\n"+
            "  In Use: %d\n"+
            "  Idle: %d\n"+
            "  Wait Count: %d\n"+
            "  Wait Duration: %v\n"+
            "  Max Idle Closed: %d\n"+
            "  Max Lifetime Closed: %d\n",
            stats.OpenConnections,
            stats.InUse,
            stats.Idle,
            stats.WaitCount,
            stats.WaitDuration,
            stats.MaxIdleClosed,
            stats.MaxLifetimeClosed,
        )
    }
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    db.SetMaxOpenConns(10)
    
    // Monitorear en goroutine separada
    go monitorPool(db, 5*time.Second)
    
    // Simular carga
    for i := 0; i < 100; i++ {
        go func(id int) {
            _, _ = db.Exec("SELECT 1")
        }(i)
    }
    
    time.Sleep(30 * time.Second)
}
```

### 39.7.4 Anti-patrones Comunes

```go
package main

import (
    "database/sql"
    "log"
)

// ❌ ANTI-PATRÓN 1: Crear novo DB para cada query
func badExample1() {
    for i := 0; i < 1000; i++ {
        db, _ := sql.Open("sqlite", "mydb.db")  // MALO: crear nuevo
        db.QueryRow("SELECT 1")
        db.Close()
    }
}

// ✅ CORRECTO: Reutilizar DB global
var globalDB *sql.DB

func goodExample1() {
    for i := 0; i < 1000; i++ {
        globalDB.QueryRow("SELECT 1")  // Reutilizar
    }
}

// ❌ ANTI-PATRÓN 2: No cerrar Rows
func badExample2(db *sql.DB) {
    rows, _ := db.Query("SELECT * FROM users")
    for rows.Next() {
        // ... procesar
    }
    // Nunca cerrar rows
}

// ✅ CORRECTO: Siempre defer Close()
func goodExample2(db *sql.DB) {
    rows, _ := db.Query("SELECT * FROM users")
    defer rows.Close()
    for rows.Next() {
        // ... procesar
    }
}

// ❌ ANTI-PATRÓN 3: Ignoring RowsAffected
func badExample3(db *sql.DB) {
    result, _ := db.Exec("UPDATE users SET name = ? WHERE id = ?", "New Name", 999)
    // Ignora si realmente se actualizó
}

// ✅ CORRECTO: Verificar cambios
func goodExample3(db *sql.DB) {
    result, _ := db.Exec("UPDATE users SET name = ? WHERE id = ?", "New Name", 999)
    rows, _ := result.RowsAffected()
    if rows == 0 {
        log.Println("Usuario no encontrado")
    }
}
```

---

## 39.8 Manejo de Errores

### 39.8.1 Tipos de Error Comunes

```go
package main

import (
    "database/sql"
    "errors"
    "log"
)

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Error 1: Fila no encontrada
    var name string
    err := db.QueryRow("SELECT name FROM users WHERE id = ?", 999).Scan(&name)
    if err == sql.ErrNoRows {
        log.Println("Usuario no encontrado")
    }
    
    // Error 2: Tipo incorrecto
    var age int
    err = db.QueryRow("SELECT age FROM users WHERE id = ?", 1).Scan(&age)
    if err != nil {
        log.Printf("Error de tipo: %v\n", err)
    }
    
    // Error 3: Conexión cerrada
    db.Close()
    _, err = db.Query("SELECT * FROM users")
    if err != nil {
        log.Printf("BD cerrada: %v\n", err)
    }
    
    // Error 4: SQL injection (con prepared statements)
    stmt, _ := db.Prepare("SELECT * FROM users WHERE name = ?")
    defer stmt.Close()
    rows, err := stmt.Query("'; DROP TABLE users; --")
    defer rows.Close()
    // Esto NO ejecutará el DROP TABLE porque está parametrizado
}
```

### 39.8.2 Retry Logic Robusto

```go
package main

import (
    "database/sql"
    "errors"
    "log"
    "math"
    "time"
)

// RetryConfig configura comportamiento de reintento
type RetryConfig struct {
    MaxRetries  int
    InitialWait time.Duration
    MaxWait     time.Duration
}

// RetryableFunc es una función que puede ser reintentada
type RetryableFunc func() error

// Retry ejecuta una función con lógica de reintento exponencial
func Retry(config RetryConfig, fn RetryableFunc) error {
    var lastErr error
    
    for attempt := 0; attempt < config.MaxRetries; attempt++ {
        if err := fn(); err != nil {
            lastErr = err
            
            // Determinar si es retryable
            if !isRetryable(err) {
                return err
            }
            
            // Esperar antes de reintentar
            wait := config.InitialWait * time.Duration(math.Pow(2, float64(attempt)))
            if wait > config.MaxWait {
                wait = config.MaxWait
            }
            
            log.Printf("Intento %d falló (%v), reintentando en %v\n", 
                attempt+1, err, wait)
            time.Sleep(wait)
        } else {
            return nil // Éxito
        }
    }
    
    return lastErr
}

// isRetryable determina si un error es retryable
func isRetryable(err error) bool {
    // Errors que NO debemos reintentar
    if errors.Is(err, sql.ErrNoRows) {
        return false
    }
    
    // Algunos errores de conexión son retryable
    if errors.Is(err, sql.ErrConnDone) {
        return true
    }
    
    return true
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    config := RetryConfig{
        MaxRetries:  3,
        InitialWait: 100 * time.Millisecond,
        MaxWait:     5 * time.Second,
    }
    
    err := Retry(config, func() error {
        rows, err := db.Query("SELECT * FROM users")
        if err != nil {
            return err
        }
        defer rows.Close()
        return nil
    })
    
    if err != nil {
        log.Fatal(err)
    }
}
```

### 39.8.3 Error Wrapping

```go
package main

import (
    "fmt"
    "log"
)

// CustomError envuelve errores de BD
type CustomError struct {
    Operation string
    Query     string
    Err       error
}

func (e *CustomError) Error() string {
    return fmt.Sprintf("BD %s falló: %v (query: %s)", e.Operation, e.Err, e.Query)
}

func executeQuery(query string) error {
    // Simular error
    return &CustomError{
        Operation: "SELECT",
        Query:     query,
        Err:       fmt.Errorf("timeout"),
    }
}

func main() {
    err := executeQuery("SELECT * FROM users")
    log.Println(err)
}
```

---

## 39.9 Scanning Avanzado

### 39.9.1 Tipos Null Seguros

```go
package main

import (
    "database/sql"
    "log"
)

type User struct {
    ID       int
    Name     string
    Email    sql.NullString    // Puede ser NULL
    Phone    sql.NullString    // Puede ser NULL
    Age      sql.NullInt64     // Puede ser NULL
    Balance  sql.NullFloat64   // Puede ser NULL
    LastSeen sql.NullTime      // Puede ser NULL
    IsActive sql.NullBool      // Puede ser NULL
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    var user User
    err := db.QueryRow(
        "SELECT id, name, email, phone, age, balance, last_seen, is_active FROM users WHERE id = ?",
        1,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Age, &user.Balance, &user.LastSeen, &user.IsActive)
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Usar valores nullos
    if user.Email.Valid {
        log.Printf("Email: %s\n", user.Email.String)
    } else {
        log.Println("Email no proporcionado")
    }
    
    if user.Age.Valid {
        log.Printf("Edad: %d\n", user.Age.Int64)
    }
}
```

### 39.9.2 Custom Scanners

```go
package main

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "log"
)

// JSONData es un tipo personalizado que puede scannear JSON
type JSONData map[string]interface{}

func (j JSONData) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("no se puede escanear de %T a JSONData", value)
    }
    return json.Unmarshal(bytes, &j)
}

func (j JSONData) Value() (driver.Value, error) {
    return json.Marshal(j)
}

type Config struct {
    ID    int
    Name  string
    Data  JSONData
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    var config Config
    err := db.QueryRow(
        "SELECT id, name, data FROM configs WHERE id = ?",
        1,
    ).Scan(&config.ID, &config.Name, &config.Data)
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Config: %+v\n", config)
    log.Printf("Data: %+v\n", config.Data)
}
```

### 39.9.3 Slice y Array Scanning

```go
package main

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "log"
)

// StringSlice es un slice de strings que puede scannear
type StringSlice []string

func (s StringSlice) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("no se puede escanear de %T a StringSlice", value)
    }
    return json.Unmarshal(bytes, (*[]string)(&s))
}

func (s StringSlice) Value() (driver.Value, error) {
    return json.Marshal(s)
}

type Post struct {
    ID    int
    Title string
    Tags  StringSlice  // JSON array en BD
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    var post Post
    err := db.QueryRow(
        "SELECT id, title, tags FROM posts WHERE id = ?",
        1,
    ).Scan(&post.ID, &post.Title, &post.Tags)
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Post: %s\n", post.Title)
    log.Printf("Tags: %v\n", post.Tags)
}
```

---

## 39.10 Context Integration

### 39.10.1 Context para Cancelación y Timeouts

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Context con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    rows, err := db.QueryContext(ctx,
        "SELECT id, name FROM users WHERE id > ? LIMIT 10",
        0,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        select {
        case <-ctx.Done():
            log.Println("Context cancelado")
            return
        default:
            var id int
            var name string
            if err := rows.Scan(&id, &name); err != nil {
                log.Fatal(err)
            }
            log.Printf("%d: %s\n", id, name)
        }
    }
}
```

### 39.10.2 Contextos Anidados

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

func queryUserWithContext(ctx context.Context, db *sql.DB, userID int) error {
    // Añadir timeout adicional si el context padre no lo tiene
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    var name string
    err := db.QueryRowContext(
        ctx,
        "SELECT name FROM users WHERE id = ?",
        userID,
    ).Scan(&name)
    
    if err == context.DeadlineExceeded {
        log.Println("Query excedió timeout")
    }
    
    return err
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Context padre con timeout más largo
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := queryUserWithContext(ctx, db, 1); err != nil {
        log.Fatal(err)
    }
}
```

### 39.10.3 Context con Valores

```go
package main

import (
    "context"
    "database/sql"
    "log"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func executeWithUserContext(ctx context.Context, db *sql.DB) {
    // Extraer valor del context
    userID := ctx.Value(UserIDKey).(int)
    
    var name string
    err := db.QueryRowContext(
        ctx,
        "SELECT name FROM users WHERE id = ? AND created_by = ?",
        1, userID,
    ).Scan(&name)
    
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Crear context con valores
    ctx := context.WithValue(context.Background(), UserIDKey, 42)
    
    executeWithUserContext(ctx, db)
}
```

---

## 39.11 Buenas Prácticas y Patterns

### 39.11.1 Migrations

```go
package main

import (
    "database/sql"
    "log"
)

// Schema define la estructura de la BD
var Schema = []string{
    `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        age INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`,
    `CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    )`,
    `CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        text TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    )`,
}

func runMigrations(db *sql.DB) error {
    for i, migration := range Schema {
        _, err := db.Exec(migration)
        if err != nil {
            log.Printf("Error en migración %d: %v\n", i+1, err)
            return err
        }
        log.Printf("Migración %d completada\n", i+1)
    }
    return nil
}

func main() {
    db, _ := sql.Open("sqlite", "myapp.db")
    defer db.Close()
    
    if err := runMigrations(db); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Base de datos inicializada")
}
```

### 39.11.2 Query Builder Pattern

```go
package main

import (
    "database/sql"
    "fmt"
    "strings"
)

// QueryBuilder construye queries dinámicos
type QueryBuilder struct {
    selectCols []string
    from       string
    whereCond  []string
    args       []interface{}
    orderBy    string
    limit      int
}

func NewQuery(from string) *QueryBuilder {
    return &QueryBuilder{
        from: from,
    }
}

func (qb *QueryBuilder) Select(cols ...string) *QueryBuilder {
    qb.selectCols = cols
    return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
    qb.whereCond = append(qb.whereCond, condition)
    qb.args = append(qb.args, args...)
    return qb
}

func (qb *QueryBuilder) OrderBy(col string) *QueryBuilder {
    qb.orderBy = col
    return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
    qb.limit = n
    return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
    // Construir SELECT
    if len(qb.selectCols) == 0 {
        qb.selectCols = []string{"*"}
    }
    query := "SELECT " + strings.Join(qb.selectCols, ", ") + " FROM " + qb.from
    
    // Construir WHERE
    if len(qb.whereCond) > 0 {
        query += " WHERE " + strings.Join(qb.whereCond, " AND ")
    }
    
    // Construir ORDER BY
    if qb.orderBy != "" {
        query += " ORDER BY " + qb.orderBy
    }
    
    // Construir LIMIT
    if qb.limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", qb.limit)
    }
    
    return query, qb.args
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    // Usar query builder
    qb := NewQuery("users").
        Select("id", "name", "email").
        Where("age > ?", 18).
        Where("is_active = ?", true).
        OrderBy("name").
        Limit(10)
    
    query, args := qb.Build()
    fmt.Println("Query:", query)
    fmt.Println("Args:", args)
    
    rows, _ := db.Query(query, args...)
    defer rows.Close()
}
```

### 39.11.3 Repository Pattern

```go
package main

import (
    "database/sql"
    "log"
)

type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

// UserRepository maneja todas las operaciones de usuario
type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id int) (*User, error) {
    var user User
    err := r.db.QueryRow(
        "SELECT id, name, email, age FROM users WHERE id = ?",
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) FindAll() ([]User, error) {
    rows, err := r.db.Query("SELECT id, name, email, age FROM users ORDER BY name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, rows.Err()
}

func (r *UserRepository) Create(user *User) (int, error) {
    result, err := r.db.Exec(
        "INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
        user.Name, user.Email, user.Age,
    )
    if err != nil {
        return 0, err
    }
    
    id, err := result.LastInsertId()
    return int(id), err
}

func (r *UserRepository) Update(user *User) error {
    result, err := r.db.Exec(
        "UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?",
        user.Name, user.Email, user.Age, user.ID,
    )
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}

func (r *UserRepository) Delete(id int) error {
    result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}

func main() {
    db, _ := sql.Open("sqlite", "mydb.db")
    defer db.Close()
    
    repo := NewUserRepository(db)
    
    // Crear usuario
    newUser := &User{Name: "Alice", Email: "alice@example.com", Age: 30}
    id, _ := repo.Create(newUser)
    log.Printf("Usuario creado con ID: %d\n", id)
    
    // Buscar usuario
    user, _ := repo.FindByID(id)
    log.Printf("Usuario encontrado: %+v\n", user)
    
    // Actualizar usuario
    user.Age = 31
    repo.Update(user)
    
    // Listar todos
    users, _ := repo.FindAll()
    log.Printf("Total usuarios: %d\n", len(users))
}
```

### 39.11.4 ORMs en Go

```go
package main

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// Si necesitas un ORM completo (similar a SQLAlchemy), GORM es la opción estándar

type User struct {
    ID    uint
    Name  string
    Email string `gorm:"uniqueIndex"`
    Age   int
}

type Post struct {
    ID      uint
    Title   string
    Content string
    UserID  uint
}

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Migrar schema
    db.AutoMigrate(&User{}, &Post{})
    
    // Crear
    user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
    db.Create(&user)
    
    // Leer
    var result User
    db.First(&result, user.ID)
    log.Printf("Usuario: %+v\n", result)
    
    // Actualizar
    db.Model(&result).Update("age", 31)
    
    // Borrar
    db.Delete(&result)
}

// Nota: database/sql es para control fino
//       GORM es para comodidad y rapidez
```

---

## Ejercicios Progresivos

### Ejercicio 1: Conexión Simple y Listado

**Objetivo**: Conectar a BD y listar registros

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "modernc.org/sqlite"
)

func main() {
    // 1. Crear BD de prueba
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // 2. Crear tabla
    _, err = db.Exec(`
        CREATE TABLE products (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            price REAL NOT NULL,
            quantity INTEGER
        )
    `)
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. Insertar datos de prueba
    products := []struct {
        name     string
        price    float64
        quantity int
    }{
        {"Laptop", 999.99, 5},
        {"Mouse", 29.99, 50},
        {"Teclado", 79.99, 30},
    }
    
    for _, p := range products {
        _, err := db.Exec(
            "INSERT INTO products (name, price, quantity) VALUES (?, ?, ?)",
            p.name, p.price, p.quantity,
        )
        if err != nil {
            log.Fatal(err)
        }
    }
    
    // 4. Listar todos los productos
    rows, err := db.Query("SELECT id, name, price, quantity FROM products")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    fmt.Println("=== Productos ===")
    for rows.Next() {
        var id int
        var name string
        var price float64
        var quantity int
        
        if err := rows.Scan(&id, &name, &price, &quantity); err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("ID: %d | Nombre: %s | Precio: %.2f | Cantidad: %d\n",
            id, name, price, quantity)
    }
}
```

**Salida esperada:**
```
=== Productos ===
ID: 1 | Nombre: Laptop | Precio: 999.99 | Cantidad: 5
ID: 2 | Nombre: Mouse | Precio: 29.99 | Cantidad: 50
ID: 3 | Nombre: Teclado | Precio: 79.99 | Cantidad: 30
```

---

### Ejercicio 2: CRUD Operations

**Objetivo**: Implementar operaciones CRUD completas

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "modernc.org/sqlite"
)

type Article struct {
    ID       int
    Title    string
    Content  string
    Views    int
}

func main() {
    db, _ := sql.Open("sqlite", ":memory:")
    defer db.Close()
    
    // Crear tabla
    db.Exec(`
        CREATE TABLE articles (
            id INTEGER PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT,
            views INTEGER DEFAULT 0
        )
    `)
    
    // CREATE
    fmt.Println("=== CREATE ===")
    result, _ := db.Exec(
        "INSERT INTO articles (title, content, views) VALUES (?, ?, ?)",
        "Go Basics", "Introduction to Go...", 0,
    )
    id, _ := result.LastInsertId()
    fmt.Printf("Artículo creado con ID: %d\n\n", id)
    
    // READ
    fmt.Println("=== READ ===")
    var article Article
    db.QueryRow("SELECT id, title, content, views FROM articles WHERE id = ?", id).
        Scan(&article.ID, &article.Title, &article.Content, &article.Views)
    fmt.Printf("Artículo: %+v\n\n", article)
    
    // UPDATE
    fmt.Println("=== UPDATE ===")
    result, _ = db.Exec(
        "UPDATE articles SET views = views + 1 WHERE id = ?",
        id,
    )
    rows, _ := result.RowsAffected()
    fmt.Printf("Filas actualizadas: %d\n", rows)
    
    // Verificar actualización
    db.QueryRow("SELECT views FROM articles WHERE id = ?", id).Scan(&article.Views)
    fmt.Printf("Nuevas vistas: %d\n\n", article.Views)
    
    // DELETE
    fmt.Println("=== DELETE ===")
    result, _ = db.Exec("DELETE FROM articles WHERE id = ?", id)
    rows, _ = result.RowsAffected()
    fmt.Printf("Filas eliminadas: %d\n", rows)
}
```

---

### Ejercicio 3: Transacciones (Transferencia Bancaria)

**Objetivo**: Implementar transacciones ACID

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "modernc.org/sqlite"
)

func main() {
    db, _ := sql.Open("sqlite", ":memory:")
    defer db.Close()
    
    // Crear tabla
    db.Exec(`
        CREATE TABLE accounts (
            id INTEGER PRIMARY KEY,
            name TEXT,
            balance REAL
        )
    `)
    
    // Insertar datos de prueba
    db.Exec("INSERT INTO accounts (id, name, balance) VALUES (1, 'Alice', 1000)")
    db.Exec("INSERT INTO accounts (id, name, balance) VALUES (2, 'Bob', 500)")
    
    fmt.Println("=== Estado Inicial ===")
    printAccounts(db)
    
    // Transferencia con transacción
    fmt.Println("\n=== Transferencia: Alice -> Bob (200) ===")
    if err := transfer(db, 1, 2, 200); err != nil {
        fmt.Printf("Transferencia fallida: %v\n", err)
    } else {
        fmt.Println("Transferencia exitosa")
    }
    
    fmt.Println("\n=== Estado Final ===")
    printAccounts(db)
    
    // Intentar transferencia que falla
    fmt.Println("\n=== Transferencia fallida: Bob -> Alice (1000) ===")
    if err := transfer(db, 2, 1, 1000); err != nil {
        fmt.Printf("Transferencia rechazada: %v\n", err)
    }
    
    fmt.Println("\n=== Estado (sin cambios) ===")
    printAccounts(db)
}

func transfer(db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // Verificar fondos
    var balance float64
    err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromID).
        Scan(&balance)
    if err != nil {
        return err
    }
    
    if balance < amount {
        return fmt.Errorf("fondos insuficientes")
    }
    
    // Restar
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?",
        amount, fromID)
    if err != nil {
        return err
    }
    
    // Sumar
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?",
        amount, toID)
    if err != nil {
        return err
    }
    
    err = tx.Commit()
    return err
}

func printAccounts(db *sql.DB) {
    rows, _ := db.Query("SELECT id, name, balance FROM accounts ORDER BY id")
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        var balance float64
        rows.Scan(&id, &name, &balance)
        fmt.Printf("%d. %s: $%.2f\n", id, name, balance)
    }
}
```

---

### Ejercicio 4: Prepared Statements

**Objetivo**: Usar prepared statements para seguridad y performance

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "modernc.org/sqlite"
)

func main() {
    db, _ := sql.Open("sqlite", ":memory:")
    defer db.Close()
    
    // Crear tabla
    db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY,
            username TEXT UNIQUE,
            email TEXT,
            age INTEGER
        )
    `)
    
    // Prepared statement para inserción
    fmt.Println("=== Prepared Statement para INSERT ===")
    insertStmt, _ := db.Prepare(
        "INSERT INTO users (username, email, age) VALUES (?, ?, ?)")
    defer insertStmt.Close()
    
    // Usar múltiples veces
    data := []struct {
        user  string
        email string
        age   int
    }{
        {"alice", "alice@example.com", 28},
        {"bob", "bob@example.com", 35},
        {"carol", "carol@example.com", 42},
    }
    
    for _, d := range data {
        _, err := insertStmt.Exec(d.user, d.email, d.age)
        if err != nil {
            log.Printf("Error: %v\n", err)
        } else {
            fmt.Printf("Insertado: %s\n", d.user)
        }
    }
    
    // Prepared statement para búsqueda
    fmt.Println("\n=== Prepared Statement para SELECT ===")
    searchStmt, _ := db.Prepare(
        "SELECT id, username, email, age FROM users WHERE age > ? ORDER BY age")
    defer searchStmt.Close()
    
    // Usar múltiples veces con diferentes parámetros
    for _, minAge := range []int{25, 30, 40} {
        fmt.Printf("\nUsuarios mayores a %d años:\n", minAge)
        rows, _ := searchStmt.Query(minAge)
        for rows.Next() {
            var id int
            var username, email string
            var age int
            rows.Scan(&id, &username, &email, &age)
            fmt.Printf("  %s (%d) - %s\n", username, age, email)
        }
        rows.Close()
    }
}
```

---

### Ejercicio 5: Connection Pooling y Monitoreo

**Objetivo**: Configurar y monitorear pool de conexiones

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "sync"
    "time"
    _ "modernc.org/sqlite"
)

func main() {
    db, _ := sql.Open("sqlite", ":memory:")
    defer db.Close()
    
    // Crear tabla
    db.Exec(`
        CREATE TABLE tasks (
            id INTEGER PRIMARY KEY,
            name TEXT,
            duration_ms INTEGER
        )
    `)
    
    // Configurar pool
    db.SetMaxOpenConns(5)
    db.SetMaxIdleConns(2)
    db.SetConnMaxLifetime(1 * time.Minute)
    
    fmt.Println("=== Configuración del Pool ===")
    fmt.Println("MaxOpenConns: 5")
    fmt.Println("MaxIdleConns: 2")
    
    // Simular carga concurrente
    fmt.Println("\n=== Ejecutando 50 queries concurrentes ===")
    var wg sync.WaitGroup
    results := make(chan string, 50)
    
    start := time.Now()
    
    for i := 1; i <= 50; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // Simular query
            var result int
            err := db.QueryRow("SELECT ?", id).Scan(&result)
            if err != nil {
                results <- fmt.Sprintf("Error en query %d: %v", id, err)
            } else {
                results <- fmt.Sprintf("Query %d completado", id)
            }
        }(i)
    }
    
    wg.Wait()
    close(results)
    
    elapsed := time.Since(start)
    fmt.Printf("\nTiempo total: %v\n", elapsed)
    
    // Mostrar estadísticas
    fmt.Println("\n=== Estadísticas del Pool ===")
    stats := db.Stats()
    fmt.Printf("Conexiones abiertas: %d\n", stats.OpenConnections)
    fmt.Printf("En uso: %d\n", stats.InUse)
    fmt.Printf("Inactivas: %d\n", stats.Idle)
    fmt.Printf("Esperas: %d\n", stats.WaitCount)
    fmt.Printf("Tiempo total de espera: %v\n", stats.WaitDuration)
}
```

---

## Resumen de Conceptos Clave

| Concepto | Descripción | Ejemplo |
|----------|-------------|---------|
| **sql.Open()** | Crea pool, no abre conexión | `db, _ := sql.Open("sqlite", "file.db")` |
| **Query()** | Para múltiples filas | `rows, _ := db.Query("SELECT * FROM users")` |
| **QueryRow()** | Para una fila | `db.QueryRow("SELECT name FROM users WHERE id=?", 1)` |
| **Exec()** | Para INSERT/UPDATE/DELETE | `db.Exec("INSERT INTO users...")` |
| **Prepared Statements** | Seguros y rápidos | `stmt, _ := db.Prepare("SELECT * FROM users WHERE id=?")` |
| **Transacciones** | ACID guarantee | `tx, _ := db.Begin(); ... tx.Commit()` |
| **Pool de Conexiones** | Gestión automática | `db.SetMaxOpenConns(25)` |
| **sql.NullXxx** | Tipos seguros para NULL | `var email sql.NullString` |
| **Context** | Cancelación y timeouts | `db.QueryContext(ctx, query)` |
| **Repository Pattern** | Abstracción de acceso a datos | Implementar CRUD methods |

---

## Conclusión

El package `database/sql` de Go proporciona una abstracción elegante y poderosa para trabajar con bases de datos. Su enfoque minimalista pero completo permite escribir código seguro, eficiente y mantenible. Combinándolo con patrones como Repository y Query Builder, se pueden construir aplicaciones de datos sólidas y escalables.

La clave está en:
1. Usar **prepared statements** para seguridad y performance
2. Gestionar **transacciones** adecuadamente
3. Configurar el **pool de conexiones** según la carga
4. Manejar **errores** explícitamente
5. Aplicar **patterns** como Repository para mejor organización


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/39-database-y-sql/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/39-database-y-sql):

```bash
cd examples/39-database-y-sql
go run .
```
