# Capítulo 11: Structs - Estructuras de datos personalizadas

## Índice del Capítulo 11

1. [11.1 ¿Qué es un Struct?](#111-qué-es-un-struct)
2. [11.2 Declaración y Definición](#112-declaración-y-definición)
3. [11.3 Inicialización de Structs](#113-inicialización-de-structs)
4. [11.4 Acceso a Campos](#114-acceso-a-campos)
5. [11.5 Structs Anidados y Embedding](#115-structs-anidados-y-embedding)
6. [11.6 Etiquetas (Tags) de Struct](#116-etiquetas-tags-de-struct)
7. [11.7 Comparación y Copia](#117-comparación-y-copia)
8. [11.8 Punteros a Structs](#118-punteros-a-structs)
9. [11.9 Métodos y Receivers](#119-métodos-y-receivers)
10. [11.10 Buenas Prácticas](#1110-buenas-prácticas)

---

## 11.1 ¿Qué es un Struct?

### Definición Conceptual

Un **struct** es una colección de campos (variables) agrupados bajo un nombre:

```

 Struct = Tipo Compuesto │

 Agrupa datos relacionados│
 Cada campo tiene nombre  │
 Cada campo tiene tipo    │
 Inmutable (no puedes     │
  cambiar estructura)     │


Ejemplo real:
    Usuario:
    ├─ Nombre: string
    ├─ Edad: int
    ├─ Email: string
    └─ Activo: bool
```

### Comparación: Struct vs Map

```

 Aspecto      │ Map                 │ Struct           │
clear
 Estructura   │ Dinámica (runtime)  │ Fija (compile)   │
 Claves       │ Desconocidas        │ Conocidas        │
 Tipos        │ Puede variar        │ Tipado fuerte    │
 Performance  │ Acceso O(1)         │ Acceso O(1)      │
 Type safety  │ Débil               │ Fuerte           │
 Uso          │ JSON, config dinámica│ Datos tipados   │

```

### Por Qué Usar Structs

```
SIN Structs:
    nombre := "Juan"
    edad := 25
    email := "juan@example.com"
    activo := true

CON Structs:
    type Usuario struct {
        Nombre string
        Edad   int
        Email  string
        Activo bool
    }
    
    usuario := Usuario{
        Nombre: "Juan",
        Edad: 25,
        Email: "juan@example.com",
        Activo: true,
    }

Beneficios:
 Organización clara
 Type safety
 Documentación automática
 Métodos (next chapter)
 Serialización (JSON, XML)
```

---

## 11.2 Declaración y Definición

### Definición Básica

```go
type Usuario struct {
    Nombre string    // Campo 1
    Edad   int       // Campo 2
    Email  string    // Campo 3
    Activo bool      // Campo 4
}
```

**Componentes:**
```
type Usuario struct {
  ↑    ↑       ↑
  |    |       └─ Abre struct
  |    └───────── Nombre del tipo
  └────────────── Define nuevo tipo
```

### Campos de Mismo Tipo (Sintaxis Corta)

```go
// Verbose
type Punto struct {
    X float64
    Y float64
    Z float64
}

// Conciso
type Punto struct {
    X, Y, Z float64
}
```

### Campos Privados vs Públicos

**Privado (minúscula):**

```go
type Usuario struct {
    nombre string     // Solo accesible dentro del package
}

u := Usuario{"Juan"}
fmt.Println(u.nombre)   // ❌ ERROR si está en otro package
```

**Público (MAYÚSCULA):**

```go
type Usuario struct {
    Nombre string      // Accesible desde otros packages
}

u := Usuario{"Juan"}
fmt.Println(u.Nombre)  // ✅ OK
```

### Structs Anónimos (Sin Nombre de Tipo)

```go
// Uso local, sin nombrar el tipo
usuario := struct {
    Nombre string
    Edad   int
}{
    Nombre: "Juan",
    Edad: 25,
}

fmt.Println(usuario.Nombre)  // Juan
```

---

## 11.3 Inicialización de Structs

### Inicialización con Valores Zero

```go
type Usuario struct {
    Nombre string
    Edad   int
    Activo bool
}

var u Usuario    // Inicialización implícita

fmt.Println(u)   // {  0 false}
// Nombre = "", Edad = 0, Activo = false
```

### Inicialización con Nombres de Campos

```go
u := Usuario{
    Nombre: "Juan",
    Edad: 25,
    Activo: true,
}

fmt.Println(u.Nombre)  // Juan
```

### Inicialización Posicional (NO recomendado)

```go
// Depende del ORDEN de campos
u := Usuario{"Juan", 25, true}

// Problema: Si cambias orden de campos en struct, se rompe código
```

### Inicialización Parcial

```go
// Solo algunos campos
u := Usuario{
    Nombre: "Juan",
    Edad: 25,
    // Activo no especificado, es false (zero value)
}

fmt.Println(u.Activo)  // false
```

### Puntero a Struct

```go
p := &Usuario{
    Nombre: "Juan",
    Edad: 25,
}

// Acceso a través de puntero (Go autocompleta *)
fmt.Println(p.Nombre)      // Juan (no necesita *p.Nombre)
```

---

## 11.4 Acceso a Campos

### Lectura de Campos

```go
u := Usuario{
    Nombre: "Juan",
    Edad: 25,
}

fmt.Println(u.Nombre)  // Juan
fmt.Println(u.Edad)    // 25
```

### Modificación de Campos

```go
u := Usuario{
    Nombre: "Juan",
    Edad: 25,
}

u.Edad = 26             // Modificar
fmt.Println(u.Edad)    // 26
```

### A Través de Puntero

```go
u := Usuario{Nombre: "Juan"}
p := &u

// Go autocompleta el dereference
p.Nombre = "María"      // Equivalente a (*p).Nombre = "María"
fmt.Println(u.Nombre)  // María (se modificó el original)
```

### Verificación de Igualdad

```go
u1 := Usuario{"Juan", 25}
u2 := Usuario{"Juan", 25}
u3 := Usuario{"María", 30}

u1 == u2    // true (mismos valores)
u1 == u3    // false (valores diferentes)
```

---

## 11.5 Structs Anidados y Embedding

### Struct Dentro de Struct

```go
type Direccion struct {
    Calle   string
    Ciudad  string
    Pais    string
}

type Usuario struct {
    Nombre    string
    Edad      int
    Direccion Direccion   // Struct anidado
}

u := Usuario{
    Nombre: "Juan",
    Edad: 25,
    Direccion: Direccion{
        Calle: "Calle Principal 123",
        Ciudad: "Madrid",
        Pais: "España",
    },
}

fmt.Println(u.Direccion.Calle)  // Calle Principal 123
```

### Embedding (Composición)

Go permite "incrustación" de structs (no heredencia):

```go
type Persona struct {
    Nombre string
    Edad   int
}

type Empleado struct {
    Persona      // Embedding: los campos de Persona están aquí
    EmpleadoID   string
    Salario      float64
}

e := Empleado{
    Persona: Persona{
        Nombre: "Juan",
        Edad: 30,
    },
    EmpleadoID: "E001",
    Salario: 50000,
}

// Acceso directo (campos promovidos)
fmt.Println(e.Nombre)          // Juan (no e.Persona.Nombre)
fmt.Println(e.EmpleadoID)      // E001
```

**Diferencia:**

```go
// Composición (Go way)
e.Nombre                 // Acceso promovido

// vs Anidamiento
type Empleado struct {
    Persona Persona      // Campo con nombre
}
e.Persona.Nombre         // Acceso anidado
```

### Embedding Múltiple

```go
type Auditable struct {
    CreadoEn   time.Time
    ActualizadoEn time.Time
}

type Usuario struct {
    Nombre string
    Auditable      // Embedding múltiple
}

u := Usuario{
    Nombre: "Juan",
    Auditable: Auditable{
        CreadoEn: time.Now(),
    },
}

fmt.Println(u.CreadoEn)   // Acceso promovido
```

---

## 11.6 Etiquetas (Tags) de Struct

### ¿Qué son Tags?

Metadatos sobre campos, usados por librerías (JSON, XML, Database):

```go
type Usuario struct {
    Nombre string `json:"nombre" db:"nombre_usuario"`
    Edad   int    `json:"edad"`
    Email  string `json:"email" db:"correo_electronico"`
}
```

### Tag JSON (Común)

```go
type Usuario struct {
    Nombre string `json:"nombre"`
    Edad   int    `json:"edad"`
    Activo bool   `json:"activo"`
}

// Sin tags, se usa nombre de campo
// Con tags, se usa valor del tag
```

**Opciones de tag:**

```go
type Usuario struct {
    Nombre string `json:"nombre"`                  // Normal
    Oculto string `json:"-"`                       // Omitir en JSON
    Email  string `json:"email,omitempty"`         // Omitir si vacío
    ID     string `json:"id,string"`               // Convertir a string
}
```

### Tag Struct (Reflexión)

```go
import "reflect"

type Usuario struct {
    Nombre string `json:"nombre" doc:"Nombre del usuario"`
    Edad   int    `json:"edad" doc:"Edad en años"`
}

u := Usuario{Nombre: "Juan", Edad: 25}

// Acceder a tags en runtime
t := reflect.TypeOf(u)
for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fmt.Printf("Campo: %s, JSON tag: %s\n", 
        field.Name, 
        field.Tag.Get("json"))
}

// Output:
// Campo: Nombre, JSON tag: nombre
// Campo: Edad, JSON tag: edad
```

### Casos Comunes de Tags

```go
type Usuario struct {
    ID     int       `db:"id,primarykey" json:"id"`
    Nombre string    `db:"nombre" json:"nombre" validate:"required"`
    Email  string    `db:"email" json:"email" validate:"email"`
    Edad   int       `db:"edad" json:"edad" validate:"min=0,max=150"`
    Activo bool      `db:"activo" json:"activo"`
}

// Tags usados por:
// - db: bases de datos (SQL)
// - json: serialización JSON
// - validate: validación de campos
```

---

## 11.7 Comparación y Copia

### Comparación de Structs

```go
type Usuario struct {
    Nombre string
    Edad   int
}

u1 := Usuario{"Juan", 25}
u2 := Usuario{"Juan", 25}
u3 := Usuario{"María", 30}

u1 == u2    // true (mismos valores)
u1 != u3    // true (valores diferentes)
```

**Limitación: NO puedes comparar si hay slices o maps**

```go
type Grupo struct {
    Nombre  string
    Miembros []string  // Slice
}

g1 := Grupo{"A", []string{"Juan"}}
g2 := Grupo{"A", []string{"Juan"}}

// g1 == g2    // ❌ ERROR: cannot compare struct with []string
```

### Copia de Structs

```go
u1 := Usuario{"Juan", 25}
u2 := u1        // Copia el CONTENIDO completo

u2.Nombre = "María"

fmt.Println(u1.Nombre)  // Juan (sin cambios)
fmt.Println(u2.Nombre)  // María
```

**Copia vs Referencia:**

```go
// Copia (valor)
u1 := Usuario{"Juan", 25}
u2 := u1        // Copia

// Referencia (puntero)
u3 := Usuario{"Juan", 25}
p := &u3        // Puntero a u3
p.Nombre = "María"
fmt.Println(u3.Nombre)  // María (modificó original)
```

---

## 11.8 Punteros a Structs

### Creación de Puntero

```go
type Usuario struct {
    Nombre string
    Edad   int
}

u := Usuario{"Juan", 25}
p := &u         // Puntero a u

// o directamente
p2 := &Usuario{"Juan", 25}
```

### Acceso a Campos con Puntero

```go
u := Usuario{"Juan", 25}
p := &u

// Go autocompleta * (dereference)
fmt.Println(p.Nombre)      // Juan
fmt.Println((*p).Nombre)   // Juan (equivalente)

// Modificar
p.Edad = 26
fmt.Println(u.Edad)        // 26 (modificó el original)
```

### Puntero a Struct con make (Raro)

```go
// Casi nunca se usa make para structs
p := new(Usuario)      // Equivalente a &Usuario{}

p.Nombre = "Juan"
p.Edad = 25

fmt.Println(p)         // &{Juan 25}
```

---

## 11.9 Métodos y Receivers

### Concepto: Métodos

Un método es una función asociada a un tipo (struct):

```go
// Función regular
func sumar(a, b int) int {
    return a + b
}

// Método (función con receiver)
func (u Usuario) String() string {
    return fmt.Sprintf("%s (%d años)", u.Nombre, u.Edad)
}
```

### Definición de Método

```go
func (receiver TipoReceiver) NombreMetodo(parametros) TipoRetorno {
    // Cuerpo
}

Ejemplo:
func (u Usuario) GetEdad() int {
    return u.Edad
}

// receptor: u (de tipo Usuario)
// mtodo: GetEdad
// retorna: int
```

### Receiver por Valor vs Referencia

**Receiver por valor (copia):**

```go
func (u Usuario) SetNombre(nombre string) {
    u.Nombre = nombre   // Modifica la COPIA, no original
}

u := Usuario{"Juan", 25}
u.SetNombre("María")
fmt.Println(u.Nombre)   // Juan (no cambió)
```

**Receiver por referencia (puntero):**

```go
func (u *Usuario) SetNombre(nombre string) {
    u.Nombre = nombre   // Modifica el ORIGINAL
}

u := Usuario{"Juan", 25}
u.SetNombre("María")
fmt.Println(u.Nombre)   // María (cambió)
```

**Regla de oro:**
```
 Si el método MODIFICA el struct → receiver *Tipo
 Si el método SOLO LEE → receiver Tipo (o *Tipo si prefieres)
 Típicamente: usa *Tipo para coherencia
```

### Ejemplos Comunes

```go
type Usuario struct {
    Nombre string
    Edad   int
}

// Getter (lectura)
func (u Usuario) GetNombre() string {
    return u.Nombre
}

// Setter (modificación)
func (u *Usuario) SetNombre(nombre string) {
    u.Nombre = nombre
}

// Método que retorna algo calculado
func (u Usuario) EsAdulto() bool {
    return u.Edad >= 18
}

// Uso
u := Usuario{"Juan", 25}
fmt.Println(u.GetNombre())     // Juan
u.SetNombre("María")
fmt.Println(u.GetNombre())     // María
fmt.Println(u.EsAdulto())      // true
```

### Chaining de Métodos

```go
func (u *Usuario) SetNombre(nombre string) *Usuario {
    u.Nombre = nombre
    return u        // Retorna *Usuario para chaining
}

func (u *Usuario) SetEdad(edad int) *Usuario {
    u.Edad = edad
    return u
}

// Uso
u := &Usuario{}
u.SetNombre("Juan").SetEdad(25)
fmt.Println(u)  // &{Juan 25}
```

---

## 11.10 Buenas Prácticas

### Usar Structs para Datos Tipados

```go
// ❌ Malo: map sin estructura
usuario := map[string]interface{}{
    "nombre": "Juan",
    "edad": 25,
}

// ✅ Bueno: struct con tipos
type Usuario struct {
    Nombre string
    Edad   int
}
usuario := Usuario{"Juan", 25}
```

### Nombres Significativos

```go
// ❌ Confuso
type S struct {
    N string
    A int
}

// ✅ Claro
type Usuario struct {
    Nombre string
    Edad   int
}
```

### Inicialización Explícita

```go
// ❌ Posicional (frágil si cambias struct)
u := Usuario{"Juan", 25}

// ✅ Nombres de campos (robusto)
u := Usuario{
    Nombre: "Juan",
    Edad: 25,
}
```

### Usar Embedding para Composición

```go
// ❌ Anidamiento innecesario
type Empleado struct {
    Usuario Usuario
    Salario float64
}
e.Usuario.Nombre    // Acceso profundo

// ✅ Embedding
type Empleado struct {
    Usuario      // Campos promovidos
    Salario float64
}
e.Nombre    // Acceso directo
```

### Constructores para Validación

```go
// ✅ Constructor con validación
func NewUsuario(nombre string, edad int) (*Usuario, error) {
    if nombre == "" {
        return nil, errors.New("nombre no puede estar vacío")
    }
    if edad < 0 || edad > 150 {
        return nil, errors.New("edad inválida")
    }
    return &Usuario{nombre, edad}, nil
}

// Uso
u, err := NewUsuario("Juan", 25)
if err != nil {
    fmt.Println("Error:", err)
    return
}
```

### Documentar Structs Públicos

```go
// Usuario representa un usuario del sistema
type Usuario struct {
    // Nombre del usuario
    Nombre string
    
    // Edad en años
    Edad int
    
    // Dirección de email
    Email string
}
```

### Usar Métodos en Lugar de Funciones

```go
// ❌ Función global
func UsuarioString(u Usuario) string {
    return fmt.Sprintf("%s (%d)", u.Nombre, u.Edad)
}

// ✅ Método
func (u Usuario) String() string {
    return fmt.Sprintf("%s (%d)", u.Nombre, u.Edad)
}

u := Usuario{"Juan", 25}
fmt.Println(u)      // Usa método String() automáticamente
```

---

## Ejercicios del Capítulo 11

### Ejercicio 1: Gestor de Tareas

Crea programa que:
1. Defina struct Tarea (ID, Título, Descripción, Completada, FechaCreacion)
2. Implemente constructores NewTarea()
3. Implemente método String()
4. Implemente método Completar() que cambia estado
5. Implemente método Dias() que retorna días desde creación

### Ejercicio 2: Sistema de Cuentas Bancarias

Crea programa que:
1. Defina struct CuentaBancaria (Tituar, Saldo, NumeroCuenta, Tipo)
2. Implemente NewCuenta() con validaciones
3. Implemente métodos Depositar() y Retirar()
4. Implemente método Transferencia() entre cuentas
5. Implemente método Resumen() que muestra estado

### Ejercicio 3: Directorio de Contactos

Crea programa que:
1. Defina struct Contacto (Nombre, Email, Teléfono, Dirección)
2. Use embedding: Dirección struct dentro de Contacto
3. Implemente método String()
4. Implemente búsqueda por nombre/email
5. Implemente edición de contacto existente

### Ejercicio 4: Carrito de Compras

Crea programa que:
1. Defina struct Producto (ID, Nombre, Precio, Cantidad)
2. Defina struct Carrito (items []Producto)
3. Implemente método AgregarProducto()
4. Implemente método CalcularTotal()
5. Implemente método EliminarProducto()

### Ejercicio 5: Sistema de Vehículos

Crea programa que:
1. Defina struct Vehicle base (Marca, Modelo, Año)
2. Defina structs Auto, Moto, Camion con embedding
3. Implemente método String() en cada uno
4. Implemente método GetCapacidad()
5. Crea slice de vehículos y muestra información

---

**Fin del Capítulo 11**

---

