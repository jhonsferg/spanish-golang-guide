# Capítulo 36: Testing - Pruebas unitarias y benchmarks

## Introducción

El testing es uno de los pilares fundamentales del desarrollo de software confiable. Go incorpora soporte para testing directamente en su ecosistema mediante el paquete `testing` estándar, una filosofía pragmática que enfatiza la simplicidad y la efectividad sobre frameworks complejos.

A diferencia de lenguajes como Java (JUnit) o Python (pytest), Go adopta un enfoque minimalista pero poderoso: no necesitas aprender un framework complicado, sino que escribes tests directamente en el lenguaje mismo. Esta sección te enseñará cómo escribir tests efectivos, benchmarks, manejar mocks y lograr una cobertura de código completa.

---

## 36.1 ¿Qué es el Testing Package? Filosofía de Go

### 36.1.1 Principios de Testing en Go

Go promueve una filosofía de testing basada en estos principios:

1. **Tests son ciudadanos de primera clase**: No es un addon, sino parte integral del lenguaje
2. **Simplicidad**: El paquete `testing` es mínimo pero suficiente
3. **Parallelismo**: Tests pueden ejecutarse en paralelo por defecto
4. **Cobertura nativa**: Go proporciona herramientas de cobertura integradas
5. **Velocidad**: Tests rápidos incentivan escribir más tests

### 36.1.2 Comparación con Otros Lenguajes

| Aspecto | Go | Python (pytest) | Java (JUnit) |
|--------|-------|---------|----------|
| **Framework integrado** | Sí (testing) | No (pytest) | Sí (JUnit) |
| **Assertions** | Manual | Automáticas | Automáticas |
| **Setup/Teardown** | t.Cleanup() | fixtures | @Before/@After |
| **Paralelismo** | Nativo | Manual | Manual |
| **Benchmarks** | Integrado | Pytest-benchmark | JMH |
| **Curva aprendizaje** | Muy baja | Media | Alta |

### 36.1.3 Estructura de un Proyecto Testing

```
proyecto/
├── main.go
├── main_test.go          # Tests para main.go
├── utils.go
├── utils_test.go         # Tests para utils.go
├── api/
│   ├── handler.go
│   └── handler_test.go   # Tests para handler.go
└── go.mod
```

### 36.1.4 El Paquete testing

El paquete `testing` proporciona tipos y funciones clave:

```go
package testing

type T struct {
    // Métodos para tests normales
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Fatal(args ...interface{})
    Fatalf(format string, args ...interface{})
    Fail()
    FailNow()
    Log(args ...interface{})
    Logf(format string, args ...interface{})
    Cleanup(func())
    Run(name string, f func(*T)) bool
    // ... más métodos
}

type B struct {
    // Métodos para benchmarks
    ReportAllocs()
    ResetTimer()
    StopTimer()
    StartTimer()
    // ... más métodos
}
```

---

## 36.2 Escribir Unit Tests

### 36.2.1 Convención _test.go

Todo archivo de test en Go debe:

1. Terminar en `_test.go`
2. Estar en el mismo paquete que el código a probar
3. Contener funciones que comienzan con `Test` seguidas de una letra mayúscula

```go
// archivo_test.go
package mipaquete_test

import "testing"

func TestMiFuncion(t *testing.T) {
    // Test code
}
```

### 36.2.2 Estructura Básica de un Test

```go
package calculator

import "testing"

// Función a probar
func Add(a, b int) int {
    return a + b
}

// Test básico
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d, want %d", result, expected)
    }
}
```

### 36.2.3 Métodos de Error en *testing.T

Go no proporciona assertions automáticas. Debes verificar manualmente:

```go
func TestVariousErrors(t *testing.T) {
    // Error: registra el error pero continúa
    if 1 != 2 {
        t.Error("1 != 2")
    }
    t.Log("Este mensaje se imprime")
    
    // Errorf: como Error pero con formato
    if true != false {
        t.Errorf("true != false: %v", false)
    }
    
    // Fatal: registra y detiene el test inmediatamente
    if nil == nil {
        t.Fatal("nil == nil") // Test termina aquí
        t.Log("Esto no se ejecuta")
    }
    
    // Fatalf: como Fatal pero con formato
    if 1 > 2 {
        t.Fatalf("1 > 2: %d > %d", 1, 2)
    }
}
```

### 36.2.4 Logging en Tests

```go
func TestLogging(t *testing.T) {
    // Log básico
    t.Log("Información general")
    t.Logf("Valor: %d", 42)
    
    // Los logs solo se muestran si:
    // 1. El test falla
    // 2. Se ejecuta con -v (verbose)
    t.Log("Este log es visible con -v")
}
```

### 36.2.5 Cleanup (Limpieza)

Para tareas de cleanup después del test:

```go
func TestWithCleanup(t *testing.T) {
    // Simular apertura de recurso
    file := createTemporaryFile()
    t.Logf("Archivo creado: %s", file)
    
    // Registrar cleanup - se ejecuta al final del test
    t.Cleanup(func() {
        deleteFile(file)
        t.Logf("Archivo eliminado: %s", file)
    })
    
    // Test continúa normalmente
    if file == "" {
        t.Error("Archivo no fue creado")
    }
}
```

### 36.2.6 Skip en Tests

Saltar tests bajo ciertas condiciones:

```go
func TestCondicional(t *testing.T) {
    if testing.Short() {
        t.Skip("Saltando test largo en modo -short")
    }
    
    // Test que toma mucho tiempo
    for i := 0; i < 1000000; i++ {
        // operación intensiva
    }
}
```

---

## 36.3 Ejecutar Tests

### 36.3.1 Comandos Básicos

```bash
# Ejecutar todos los tests en el paquete actual
go test

# Ejecutar tests en todos los paquetes
go test ./...

# Ejecutar un test específico
go test -run TestAdd

# Modo verbose - muestra todos los logs
go test -v

# Tests cortos (salta tests marcados con testing.Short())
go test -short

# Ejecutar con timeout
go test -timeout 30s

# Ejecutar tests en paralelo (se establece automáticamente)
go test -parallel 4
```

### 36.3.2 Salida de go test

```
$ go test -v

=== RUN   TestAdd
    calculator_test.go:10: TestAdd iniciado
--- PASS: TestAdd (0.00s)
=== RUN   TestSubtract
    calculator_test.go:20: 5 - 3 = 2
--- PASS: TestSubtract (0.00s)
=== RUN   TestDivideByZero
    calculator_test.go:30: Error esperado: division by zero
--- PASS: TestDivideByZero (0.00s)
PASS
ok      github.com/ejemplo/calculator   0.456s
```

### 36.3.3 Flags Útiles

```bash
# Mostrar salida del test incluso si pasa
go test -v

# Ejecutar tests en secuencia (no paralelo)
go test -parallel 1

# Dejar de ejecutar en el primer fallo
go test -failfast

# Compilar pero no ejecutar
go test -c

# Ejecutar tests que coincidan con patrón regex
go test -run "TestAdd|TestSub"

# Ejecutar tests de un archivo específico
go test -run TestAdd archivo_test.go archivo.go
```

### 36.3.4 Profiling Básico

```bash
# Información de tiempo de ejecución
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Información de memoria
go test -memprofile=mem.prof
go tool pprof mem.prof
```

---

## 36.4 Subtests

### 36.4.1 ¿Por qué Subtests?

Los subtests permiten:

1. Organizar múltiples casos de prueba lógicamente
2. Ejecutar solo un subset de tests
3. Paralelismo granular
4. Mejor lectura de salida

### 36.4.2 Sintaxis Básica de Subtests

```go
func TestCalculator(t *testing.T) {
    t.Run("suma positivos", func(t *testing.T) {
        result := Add(2, 3)
        if result != 5 {
            t.Errorf("Add(2, 3) = %d, want 5", result)
        }
    })
    
    t.Run("suma con negativos", func(t *testing.T) {
        result := Add(-2, 3)
        if result != 1 {
            t.Errorf("Add(-2, 3) = %d, want 1", result)
        }
    })
    
    t.Run("resta", func(t *testing.T) {
        result := Subtract(5, 3)
        if result != 2 {
            t.Errorf("Subtract(5, 3) = %d, want 2", result)
        }
    })
}
```

Salida:

```
$ go test -v
=== RUN   TestCalculator
=== RUN   TestCalculator/suma_positivos
    calculator_test.go:10: ✓
--- PASS: TestCalculator/suma_positivos (0.00s)
=== RUN   TestCalculator/suma_con_negativos
    calculator_test.go:16: ✓
--- PASS: TestCalculator/suma_con_negativos (0.00s)
=== RUN   TestCalculator/resta
    calculator_test.go:22: ✓
--- PASS: TestCalculator/resta (0.00s)
--- PASS: TestCalculator (0.00s)
PASS
```

### 36.4.3 Ejecutar Subtests Específicos

```bash
# Ejecutar un subtest específico
go test -run TestCalculator/suma_positivos

# Ejecutar múltiples subtests con regex
go test -run TestCalculator/suma

# Ejecutar solo subtests, no el padre
go test -run TestCalculator/suma -v
```

### 36.4.4 Subtests Anidados

```go
func TestNestedSubtests(t *testing.T) {
    t.Run("operaciones básicas", func(t *testing.T) {
        t.Run("suma", func(t *testing.T) {
            if Add(2, 3) != 5 {
                t.Error("suma falló")
            }
        })
        
        t.Run("resta", func(t *testing.T) {
            if Subtract(5, 3) != 2 {
                t.Error("resta falló")
            }
        })
    })
    
    t.Run("operaciones avanzadas", func(t *testing.T) {
        t.Run("potencia", func(t *testing.T) {
            if Power(2, 3) != 8 {
                t.Error("potencia falló")
            }
        })
    })
}
```

Salida:

```
TestNestedSubtests
├── operaciones básicas
│   ├── suma
│   └── resta
└── operaciones avanzadas
    └── potencia
```

### 36.4.5 Paralelismo en Subtests

```go
func TestParallel(t *testing.T) {
    t.Run("test 1", func(t *testing.T) {
        t.Parallel()  // Ejecutar en paralelo
        // Test que no comparte estado
        result := Add(1, 2)
        if result != 3 {
            t.Error("test 1 falló")
        }
    })
    
    t.Run("test 2", func(t *testing.T) {
        t.Parallel()  // Ejecutar en paralelo
        result := Add(2, 2)
        if result != 4 {
            t.Error("test 2 falló")
        }
    })
}
```

---

## 36.5 Tabla-Driven Tests

### 36.5.1 El Patrón Tabla-Driven

Este es el patrón más común en Go para tests:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"suma positivos", 2, 3, 5},
        {"suma negativos", -2, -3, -5},
        {"suma mixtos", -2, 3, 1},
        {"suma con cero", 0, 5, 5},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", 
                    tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### 36.5.2 Ventajas del Patrón Tabla-Driven

1. **Separación clara**: Datos de prueba separados de la lógica
2. **Escalabilidad**: Fácil agregar nuevos casos
3. **Mantenibilidad**: Un solo punto de lógica de prueba
4. **Legibilidad**: Casos de prueba más claros

### 36.5.3 Tabla-Driven con Errores

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
        errMsg  string
    }{
        {"email válido", "user@example.com", false, ""},
        {"email vacío", "", true, "email is required"},
        {"sin @", "userexample.com", true, "invalid format"},
        {"múltiples @", "user@@example.com", true, "invalid format"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ValidateEmail() error = %v, wantErr %v", 
                    err, tt.wantErr)
            }
            if err != nil && err.Error() != tt.errMsg {
                t.Errorf("ValidateEmail() error message = %s, want %s",
                    err.Error(), tt.errMsg)
            }
        })
    }
}
```

### 36.5.4 Tabla-Driven Avanzada

Para casos más complejos:

```go
func TestParser(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantNodes int
        wantErr   bool
        setup     func()      // Setup opcional
        cleanup   func()      // Cleanup opcional
    }{
        {
            name:      "JSON válido",
            input:     `{"key": "value"}`,
            wantNodes: 1,
            wantErr:   false,
        },
        {
            name:      "JSON con configuración previa",
            input:     `[1, 2, 3]`,
            wantNodes: 3,
            wantErr:   false,
            setup: func() {
                initializeParser()
            },
            cleanup: func() {
                closeParser()
            },
        },
        {
            name:      "JSON inválido",
            input:     `{invalid}`,
            wantNodes: 0,
            wantErr:   true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                tt.setup()
                t.Cleanup(tt.cleanup)
            }
            
            nodes, err := ParseJSON(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ParseJSON() error = %v, wantErr %v", err, tt.wantErr)
            }
            if len(nodes) != tt.wantNodes {
                t.Errorf("ParseJSON() got %d nodes, want %d", 
                    len(nodes), tt.wantNodes)
            }
        })
    }
}
```

---

## 36.6 Testing Helpers

### 36.6.1 Funciones Auxiliares

Las funciones helper simplifican la escritura de tests:

```go
// helper.go en _test.go
package calculator_test

import "testing"

// Helper para validar resultado
func assertInt(t *testing.T, got, want int) {
    t.Helper()  // Indica que es una función auxiliar
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}

// Uso en test
func TestAddHelper(t *testing.T) {
    assertInt(t, Add(2, 3), 5)
    assertInt(t, Add(-2, 3), 1)
}
```

### 36.6.2 t.Helper() Explicado

`t.Helper()` es crucial para reports de error precisos:

```go
// Sin t.Helper() - reporta error en helper.go
func assertInt(t *testing.T, got, want int) {
    if got != want {
        t.Errorf("got %d, want %d", got, want)  // Línea 5 en helper.go
    }
}

// Con t.Helper() - reporta error en el test
func assertInt(t *testing.T, got, want int) {
    t.Helper()
    if got != want {
        t.Errorf("got %d, want %d", got, want)  // Línea 5 en calculator_test.go
    }
}
```

### 36.6.3 Helpers Comunes

```go
// Verificar igualdad genérica
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("assertEqual failed: got %v, want %v", got, want)
    }
}

// Verificar error esperado
func assertError(t *testing.T, err error, wantErr bool) {
    t.Helper()
    if (err != nil) != wantErr {
        t.Errorf("assertError failed: got error %v, wantErr %v", err, wantErr)
    }
}

// Verificar panic
func assertPanic(t *testing.T, f func()) {
    t.Helper()
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("assertPanic failed: expected panic but got none")
        }
    }()
    f()
}

// Verificar no es nil
func assertNotNil(t *testing.T, value interface{}) {
    t.Helper()
    if value == nil {
        t.Errorf("assertNotNil failed: value is nil")
    }
}
```

### 36.6.4 Test Fixtures

Datos reutilizables para tests:

```go
type TestFixture struct {
    User    *User
    Post    *Post
    Comment *Comment
}

func setupFixture(t *testing.T) *TestFixture {
    t.Helper()
    
    user := &User{ID: 1, Name: "John"}
    post := &Post{ID: 1, UserID: 1, Title: "Test"}
    comment := &Comment{ID: 1, PostID: 1, Text: "Good"}
    
    t.Cleanup(func() {
        deleteUser(user.ID)
        deletePost(post.ID)
        deleteComment(comment.ID)
    })
    
    return &TestFixture{
        User:    user,
        Post:    post,
        Comment: comment,
    }
}

func TestWithFixture(t *testing.T) {
    fixture := setupFixture(t)
    
    // Usar fixture en tests
    if fixture.User.Name != "John" {
        t.Error("user setup failed")
    }
}
```

---

## 36.7 Mocking y Stubs

### 36.7.1 Principio de Dependency Injection

Go usa interfaces para mocking:

```go
// Interfaz para abstracción
type Database interface {
    GetUser(id int) (*User, error)
    SaveUser(user *User) error
}

// Implementación real
type PostgresDB struct{}

func (db *PostgresDB) GetUser(id int) (*User, error) {
    // Conectar a PostgreSQL
}

// Mock para testing
type MockDB struct{}

func (db *MockDB) GetUser(id int) (*User, error) {
    return &User{ID: 1, Name: "Mock User"}, nil
}

// Función a probar - acepta interfaz
func GetUserWithCache(db Database, id int) *User {
    user, _ := db.GetUser(id)
    return user
}

// Test
func TestGetUserWithCache(t *testing.T) {
    mock := &MockDB{}
    user := GetUserWithCache(mock, 1)
    
    if user.Name != "Mock User" {
        t.Error("mock not working")
    }
}
```

### 36.7.2 Mock con Comportamiento Configurable

```go
type MockDB struct {
    GetUserFunc func(id int) (*User, error)
}

func (db *MockDB) GetUser(id int) (*User, error) {
    if db.GetUserFunc != nil {
        return db.GetUserFunc(id)
    }
    return nil, nil
}

func TestWithConfigurable(t *testing.T) {
    // Mock que retorna error
    mock := &MockDB{
        GetUserFunc: func(id int) (*User, error) {
            return nil, errors.New("not found")
        },
    }
    
    _, err := GetUserWithCache(mock, 999)
    if err == nil {
        t.Error("expected error for non-existent user")
    }
}
```

### 36.7.3 Mock Spy - Verificar Llamadas

```go
type SpyDB struct {
    GetUserCalls []int  // Track llamadas
}

func (db *SpyDB) GetUser(id int) (*User, error) {
    db.GetUserCalls = append(db.GetUserCalls, id)
    return &User{ID: id}, nil
}

func (db *SpyDB) SaveUser(user *User) error {
    return nil
}

func TestSpyMock(t *testing.T) {
    spy := &SpyDB{}
    
    GetUserWithCache(spy, 1)
    GetUserWithCache(spy, 2)
    GetUserWithCache(spy, 1)
    
    if len(spy.GetUserCalls) != 3 {
        t.Errorf("expected 3 calls, got %d", len(spy.GetUserCalls))
    }
    
    if spy.GetUserCalls[0] != 1 {
        t.Errorf("first call should be for id 1, got %d", spy.GetUserCalls[0])
    }
}
```

### 36.7.4 Stub HTTP

Para testing de APIs:

```go
import (
    "net/http"
    "net/http/httptest"
)

func TestAPIWithStub(t *testing.T) {
    // Crear servidor stub
    server := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"id": 1, "name": "John"}`))
        }))
    defer server.Close()
    
    // Usar servidor stub en cliente
    client := &http.Client{}
    resp, err := client.Get(server.URL + "/users/1")
    
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }
}
```

---

## 36.8 Benchmarks

### 36.8.1 Escribir Benchmarks

Los benchmarks miden performance:

```go
func BenchmarkAdd(b *testing.B) {
    // b.N es el número de iteraciones a ejecutar
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

func BenchmarkFib(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fib(20)
    }
}
```

### 36.8.2 Ejecutar Benchmarks

```bash
# Ejecutar benchmarks
go test -bench=.

# Benchmarks en el paquete actual y submódulos
go test -bench=. ./...

# Ejecutar un benchmark específico
go test -bench=BenchmarkAdd

# Verbose benchmark
go test -bench=. -v

# Incluir allocations
go test -bench=. -benchmem

# Duración de benchmark
go test -bench=. -benchtime=5s

# Comparar dos versiones
go test -bench=. -benchmem > old.txt
# hacer cambios
go test -bench=. -benchmem > new.txt
benchstat old.txt new.txt
```

### 36.8.3 Salida de Benchmark

```
$ go test -bench=. -benchmem

BenchmarkAdd-8          2000000000   0.35 ns/op   0 B/op   0 allocs/op
BenchmarkFib-8             20000   52834 ns/op   0 B/op   0 allocs/op

ok      github.com/ejemplo/fib  2.345s
```

Interpretación:
- `BenchmarkAdd-8`: 8 CPUs
- `2000000000`: 2 mil millones de iteraciones (b.N)
- `0.35 ns/op`: 0.35 nanosegundos por operación
- `0 B/op`: 0 bytes asignados
- `0 allocs/op`: 0 asignaciones de memoria

### 36.8.4 Optimizar Benchmarks

```go
// ❌ MALO - Compilador optimiza el loop
func BenchmarkBadFib(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fib(5)
    }
}

// ✓ BUENO - Usar resultado para evitar optimización
func BenchmarkGoodFib(b *testing.B) {
    var result int
    for i := 0; i < b.N; i++ {
        result = Fib(5)
    }
    _ = result  // Usar resultado
}

// ✓ EXCELENTE - Resetear timer para setup
func BenchmarkWithSetup(b *testing.B) {
    // Setup costoso
    data := loadLargeDataset()
    
    b.ResetTimer()  // Reiniciar timer, descontar setup
    
    for i := 0; i < b.N; i++ {
        processData(data)
    }
}
```

### 36.8.5 ReportAllocs()

```go
func BenchmarkStringConcat(b *testing.B) {
    b.ReportAllocs()  // Incluir info de allocations
    
    for i := 0; i < b.N; i++ {
        s := "hello"
        s += "world"
    }
}

// Salida:
// BenchmarkStringConcat-8   1000000   1045 ns/op   32 B/op   2 allocs/op
```

### 36.8.6 Comparar Versiones

```go
// Función original - lenta
func SlowStringConcat(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s
    }
    return result
}

// Función optimizada - rápida
func FastStringConcat(strs []string) string {
    var buf strings.Builder
    for _, s := range strs {
        buf.WriteString(s)
    }
    return buf.String()
}

func BenchmarkSlowConcat(b *testing.B) {
    data := generateStringSlice(1000)
    for i := 0; i < b.N; i++ {
        SlowStringConcat(data)
    }
}

func BenchmarkFastConcat(b *testing.B) {
    data := generateStringSlice(1000)
    for i := 0; i < b.N; i++ {
        FastStringConcat(data)
    }
}

// Ejecutar:
// go test -bench=. -benchmem -benchstat=benchstat
```

---

## 36.9 Examples - Documentación Ejecutable

### 36.9.1 Escribir Examples

Examples son tests que también sirven como documentación:

```go
package calculator

import "fmt"

func ExampleAdd() {
    result := Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

func ExampleAdd_negative() {
    result := Add(-2, 3)
    fmt.Println(result)
    // Output: 1
}
```

### 36.9.2 Reglas de Examples

1. Función debe empezar con `Example`
2. Debe contener comentario `// Output:`
3. La salida exacta del programa debe coincidir

```go
// ✓ CORRECTO
func ExampleDiv() {
    fmt.Println(Div(10, 2))
    // Output: 5
}

// ❌ INCORRECTO - Sin Output
func ExampleDiv() {
    fmt.Println(Div(10, 2))
}

// ❌ INCORRECTO - Output no coincide
func ExampleDiv() {
    fmt.Println(Div(10, 2))
    // Output: 6
}
```

### 36.9.3 Examples para Paquetes

```go
package mylib

import "fmt"

func ExamplePackage() {
    // Ejemplo de cómo usar el paquete
    calculator := NewCalculator()
    result := calculator.Add(2, 3)
    fmt.Println("2 + 3 =", result)
    // Output: 2 + 3 = 5
}
```

### 36.9.4 Ejecutar Examples

```bash
# Ejecutar examples como tests
go test -run Example

# Ver si examples pasan
go test -run Example -v

# Ejecutar ejemplo específico
go test -run ExampleAdd

# Los examples aparecen en godoc
go doc -http=:6060
# Visitar http://localhost:6060
```

### 36.9.5 Documentación Generada

Examples aparecen automáticamente en `godoc`:

```
$ godoc github.com/usuario/calculadora Add
```

Muestra:

```
func Add(a, b int) int
    Add suma dos números.

    Example:
        result := Add(2, 3)
        fmt.Println(result)
        // Output: 5
```

---

## 36.10 Code Coverage

### 36.10.1 Medir Cobertura

```bash
# Generar reporte de cobertura
go test -cover

# Salida: coverage: 85.2% of statements

# Generar archivo de cobertura
go test -coverprofile=coverage.out

# Ver reporte en HTML
go tool cover -html=coverage.out

# Ver cobertura de archivo específico
go tool cover -html=coverage.out -o coverage.html
```

### 36.10.2 Cobertura por Función

```go
// Archivo: math.go
package math

func Add(a, b int) int {
    return a + b
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Archivo: math_test.go
package math_test

func TestAdd(t *testing.T) {
    // Cubre la rama feliz
    result := Add(2, 3)
    if result != 5 {
        t.Error("failed")
    }
}

// Falta: test para DivideByZero

// Ejecutar:
// go test -cover
// coverage: 75% of statements (solo Add está cubierto)
```

### 36.10.3 Aumentar Cobertura

```go
// Agregando test para Divide
func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    int
        want    int
        wantErr bool
    }{
        {"division normal", 10, 2, 5, false},
        {"division por cero", 10, 0, 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)
            if (err != nil) != tt.wantErr {
                t.Errorf("Divide() error = %v, wantErr %v", err, tt.wantErr)
            }
            if err == nil && got != tt.want {
                t.Errorf("Divide() = %d, want %d", got, tt.want)
            }
        })
    }
}

// Ahora: coverage: 100% of statements
```

### 36.10.4 Interpretar Reportes

En HTML, diferentes colores:
- **Verde**: Código cubierto
- **Rojo**: Código no cubierto
- **Gris**: Código no ejecutable

### 36.10.5 Mínimo de Cobertura

Configurar en CI/CD:

```bash
# Fallar si cobertura < 80%
go test -cover -coverprofile=coverage.out
coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
if (( $(echo "$coverage < 80" | bc -l) )); then
    echo "Coverage $coverage% is below threshold 80%"
    exit 1
fi
```

---

## 36.11 Buenas Prácticas y Patrones

### 36.11.1 Testing Principles

1. **AAA Pattern** (Arrange-Act-Assert):

```go
func TestUser(t *testing.T) {
    // Arrange - preparar
    user := &User{Name: "John", Age: 30}
    
    // Act - actuar
    result := user.IsAdult()
    
    // Assert - verificar
    if !result {
        t.Error("user should be adult")
    }
}
```

2. **One Assertion per Test (cuando sea posible)**:

```go
// ✓ MEJOR - Un concepto por test
func TestUserIsAdult(t *testing.T) {
    user := &User{Age: 30}
    if !user.IsAdult() {
        t.Error("30 should be adult")
    }
}

func TestUserIsNotAdult(t *testing.T) {
    user := &User{Age: 10}
    if user.IsAdult() {
        t.Error("10 should not be adult")
    }
}

// vs

// ❌ PEOR - Múltiples conceptos
func TestUser(t *testing.T) {
    user := &User{Age: 30}
    if !user.IsAdult() {
        t.Error("30 should be adult")
    }
    
    user2 := &User{Age: 10}
    if user2.IsAdult() {
        t.Error("10 should not be adult")
    }
}
```

### 36.11.2 TDD (Test-Driven Development)

```go
// 1. Escribir test que falla
func TestSortArray(t *testing.T) {
    input := []int{3, 1, 2}
    expected := []int{1, 2, 3}
    result := Sort(input)
    // Esta función no existe aún
}

// 2. Hacer que el test pase (simple)
func Sort(arr []int) []int {
    sort.Ints(arr)
    return arr
}

// 3. Refactorizar
// (optimización o mejora de código)
```

### 36.11.3 Antipatrones en Testing

❌ **No hacer:**

```go
// 1. Tests que dependen del orden
func TestA(t *testing.T) { setupState = true }
func TestB(t *testing.T) { // Depende de TestA
    if !setupState {
        t.Error("state not set")
    }
}

// 2. Tests aleatorios
func TestRandom(t *testing.T) {
    n := rand.Intn(100)
    if n != 42 {
        t.Error("unlucky")
    }
}

// 3. Testear implementación, no comportamiento
func TestPrivateField(t *testing.T) {
    obj := &MyStruct{}
    obj.privateField = 42  // ❌ No se puede acceder
    // Testea la implementación, no el comportamiento
}

// 4. Falta de assertions
func TestFunction(t *testing.T) {
    result := MyFunction()
    // Sin verificar resultado
}

// 5. Mensajes de error pobres
func TestAdd(t *testing.T) {
    if Add(2, 3) != 5 {
        t.Error("failed")  // ❌ Poco informativo
    }
}
```

✓ **Hacer:**

```go
// 1. Tests independientes
func TestValidationA(t *testing.T) {
    if !ValidateEmail("user@example.com") {
        t.Error("valid email rejected")
    }
}

func TestValidationB(t *testing.T) {
    if ValidateEmail("invalid") {
        t.Error("invalid email accepted")
    }
}

// 2. Tests determinísticos
func TestSorting(t *testing.T) {
    input := []int{3, 1, 2}
    expected := []int{1, 2, 3}
    if !reflect.DeepEqual(Sort(input), expected) {
        t.Error("sorting failed")
    }
}

// 3. Testear comportamiento observable
func TestUser(t *testing.T) {
    user := CreateUser("John", 30)
    if !user.IsAdult() {
        t.Error("30-year-old should be adult")
    }
}

// 4. Assertions significativas
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}
```

### 36.11.4 Testing en CI/CD

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Run tests
        run: go test -race -cover ./...
      
      - name: Generate coverage
        run: go test -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v2
        with:
          files: ./coverage.out
```

### 36.11.5 Integración con Linters

```bash
# Usar gofmt, vet, golangci-lint
go fmt ./...
go vet ./...
golangci-lint run

# En tests
go test -vet=off  # Desactivar vet
```

### 36.11.6 Performance Testing Best Practices

```go
// 1. Tests significativos
func BenchmarkStringConcat(b *testing.B) {
    b.ReportAllocs()
    data := generateStrings(1000)  // Datos realistas
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        concatenateStrings(data)
    }
}

// 2. Evitar Dead Code Elimination
func BenchmarkWithCheck(b *testing.B) {
    var result int
    for i := 0; i < b.N; i++ {
        result = expensiveOperation()
    }
    _ = result  // Usar resultado
}

// 3. Comparar con baseline
// go test -bench=. -benchmem -count=5 > baseline.txt
// (hacer cambios)
// go test -bench=. -benchmem -count=5 > optimized.txt
// benchstat baseline.txt optimized.txt
```

---

## EJERCICIOS PROGRESIVOS

### Ejercicio 1: Unit Test Simple

**Objetivo**: Escribir tu primer test unitario

**Descripción**: Crea un archivo `calculator.go` con una función `Add` y un archivo `calculator_test.go` con tests.

```go
// calculator.go
package main

func Add(a, b int) int {
    return a + b
}

func Subtract(a, b int) int {
    return a - b
}
```

**Tu tarea**:
1. Escribe 3 tests para `Add` (casos positivos, negativos, cero)
2. Escribe 2 tests para `Subtract`
3. Ejecuta: `go test -v`
4. Verifica que todos pasan

**Solución esperada**:
```go
// calculator_test.go
package main

import "testing"

func TestAddPositive(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d, want 5", result)
    }
}

func TestAddNegative(t *testing.T) {
    result := Add(-2, -3)
    if result != -5 {
        t.Errorf("Add(-2, -3) = %d, want -5", result)
    }
}

func TestAddZero(t *testing.T) {
    result := Add(0, 5)
    if result != 5 {
        t.Errorf("Add(0, 5) = %d, want 5", result)
    }
}

func TestSubtract(t *testing.T) {
    result := Subtract(10, 3)
    if result != 7 {
        t.Errorf("Subtract(10, 3) = %d, want 7", result)
    }
}
```

---

### Ejercicio 2: Subtests y Paralelismo

**Objetivo**: Organizar tests con subtests y paralelismo

**Descripción**: Refactoriza los tests del ejercicio 1 usando `t.Run()` y `t.Parallel()`

**Tu tarea**:
1. Agrupa tests con `t.Run()`
2. Usa `t.Parallel()` para ejecutar en paralelo
3. Ejecuta: `go test -v -parallel 4`
4. Verifica que se ejecutan en paralelo

**Solución esperada**:
```go
// calculator_test.go (mejorado)
package main

import "testing"

func TestCalculatorOperations(t *testing.T) {
    t.Run("add operations", func(t *testing.T) {
        t.Parallel()
        
        t.Run("positive numbers", func(t *testing.T) {
            t.Parallel()
            if Add(2, 3) != 5 {
                t.Error("Add(2, 3) failed")
            }
        })
        
        t.Run("negative numbers", func(t *testing.T) {
            t.Parallel()
            if Add(-2, -3) != -5 {
                t.Error("Add(-2, -3) failed")
            }
        })
    })
    
    t.Run("subtract operations", func(t *testing.T) {
        t.Parallel()
        if Subtract(10, 3) != 7 {
            t.Error("Subtract(10, 3) failed")
        }
    })
}
```

---

### Ejercicio 3: Tabla-Driven Tests

**Objetivo**: Implementar tabla-driven tests para múltiples casos

**Descripción**: Crea un parser simple y usa tabla-driven tests

```go
// parser.go
package main

import "strconv"
import "strings"

func ParseNumbers(input string) ([]int, error) {
    parts := strings.Split(input, ",")
    var numbers []int
    for _, part := range parts {
        n, err := strconv.Atoi(strings.TrimSpace(part))
        if err != nil {
            return nil, err
        }
        numbers = append(numbers, n)
    }
    return numbers, nil
}
```

**Tu tarea**:
1. Implementa tabla-driven tests para `ParseNumbers`
2. Incluye casos: válidos, vacíos, con espacios, con errores
3. Ejecuta: `go test -v`

**Solución esperada**:
```go
// parser_test.go
package main

import (
    "testing"
    "reflect"
)

func TestParseNumbers(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    []int
        wantErr bool
    }{
        {
            name:    "simple numbers",
            input:   "1,2,3",
            want:    []int{1, 2, 3},
            wantErr: false,
        },
        {
            name:    "numbers with spaces",
            input:   "1, 2, 3",
            want:    []int{1, 2, 3},
            wantErr: false,
        },
        {
            name:    "single number",
            input:   "42",
            want:    []int{42},
            wantErr: false,
        },
        {
            name:    "invalid number",
            input:   "1,abc,3",
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseNumbers(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseNumbers() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseNumbers() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

### Ejercicio 4: Mocking e Inyección de Dependencias

**Objetivo**: Implementar mocking usando interfaces

**Descripción**: Crea un servicio que usa una interfaz y escribe tests con mocks

```go
// user.go
package main

type User struct {
    ID   int
    Name string
}

type UserRepository interface {
    GetUser(id int) (*User, error)
}

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUserName(id int) (string, error) {
    user, err := s.repo.GetUser(id)
    if err != nil {
        return "", err
    }
    return user.Name, nil
}
```

**Tu tarea**:
1. Crea un `MockRepository` que implemente `UserRepository`
2. Escribe tests para `UserService.GetUserName()` usando el mock
3. Prueba el caso de éxito y error
4. Ejecuta: `go test -v`

**Solución esperada**:
```go
// user_test.go
package main

import (
    "errors"
    "testing"
)

type MockUserRepository struct {
    users map[int]*User
}

func (m *MockUserRepository) GetUser(id int) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("user not found")
    }
    return user, nil
}

func TestGetUserNameSuccess(t *testing.T) {
    mock := &MockUserRepository{
        users: map[int]*User{
            1: {ID: 1, Name: "John"},
        },
    }
    
    service := NewUserService(mock)
    name, err := service.GetUserName(1)
    
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if name != "John" {
        t.Errorf("GetUserName(1) = %s, want John", name)
    }
}

func TestGetUserNameNotFound(t *testing.T) {
    mock := &MockUserRepository{
        users: make(map[int]*User),
    }
    
    service := NewUserService(mock)
    _, err := service.GetUserName(999)
    
    if err == nil {
        t.Error("expected error for non-existent user")
    }
}
```

---

### Ejercicio 5: Benchmarks y Coverage

**Objetivo**: Escribir benchmarks y generar reporte de cobertura

**Descripción**: Compara performance de dos implementaciones

```go
// fibonacci.go
package main

// Implementación recursiva (lenta)
func FibRecursive(n int) int {
    if n <= 1 {
        return n
    }
    return FibRecursive(n-1) + FibRecursive(n-2)
}

// Implementación iterativa (rápida)
func FibIterative(n int) int {
    if n <= 1 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

**Tu tarea**:
1. Escribe tests unitarios para ambas funciones
2. Escribe benchmarks comparando performance
3. Ejecuta: `go test -bench=. -benchmem`
4. Genera cobertura: `go test -cover`
5. Reporta: ¿Cuál es más rápida?

**Solución esperada**:
```go
// fibonacci_test.go
package main

import "testing"

func TestFibRecursive(t *testing.T) {
    tests := []struct {
        n    int
        want int
    }{
        {0, 0},
        {1, 1},
        {6, 8},
        {10, 55},
    }
    
    for _, tt := range tests {
        if got := FibRecursive(tt.n); got != tt.want {
            t.Errorf("FibRecursive(%d) = %d, want %d", tt.n, got, tt.want)
        }
    }
}

func TestFibIterative(t *testing.T) {
    tests := []struct {
        n    int
        want int
    }{
        {0, 0},
        {1, 1},
        {6, 8},
        {10, 55},
    }
    
    for _, tt := range tests {
        if got := FibIterative(tt.n); got != tt.want {
            t.Errorf("FibIterative(%d) = %d, want %d", tt.n, got, tt.want)
        }
    }
}

func BenchmarkFibRecursive(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibRecursive(20)
    }
}

func BenchmarkFibIterative(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibIterative(20)
    }
}

// Ejecutar:
// $ go test -bench=. -benchmem
// BenchmarkFibRecursive-8      50000    21843 ns/op    0 B/op    0 allocs/op
// BenchmarkFibIterative-8   50000000       21.3 ns/op   0 B/op    0 allocs/op
// (La iterativa es ~1000x más rápida)
```

---

## Conclusión

El testing en Go es simple, efectivo y poderoso. Con el paquete `testing` integrado, puedes escribir tests desde el primer día sin dependencias externas. Recuerda:

1. **Simplicidad**: Go tests son simples intencionalmente
2. **Coverage**: Mide y mejora continuamente tu cobertura
3. **Velocidad**: Tests rápidos incentivan escribir más
4. **Paralelismo**: Aprovecha los multi-cores
5. **Mocking**: Usa interfaces para desacoplar código
6. **Benchmarks**: Mide performance de código crítico
7. **CI/CD**: Automatiza tests en cada cambio

¡Escribe tests confiables y mantén tu código confiable!

---

**Próximo Capítulo**: Concurrencia y Goroutines - Aprovecha la velocidad de Go

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/36-testing-package/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/36-testing-package):

```bash
cd examples/36-testing-package
go test -v -bench=. .
```
