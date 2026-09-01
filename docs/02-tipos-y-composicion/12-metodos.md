# Capítulo 12: Métodos - Funciones asociadas a tipos

## Introducción

En el capítulo anterior, explorar los structs como colecciones de campos tipados. Ahora, avanzamos hacia la **Programación Orientada a Objetos en Go** (OOP sin herencia) mediante **métodos**: funciones que pertenecen a un tipo específico.

Un **método** es una función con un **receiver**, es decir, un parámetro especial entre `func` y el nombre que vincula la función a un tipo. Esto permite escribir código más expresivo, reutilizable y modular.

```go
// Función regular
func distancia(x, y float64) float64 {
    return math.Sqrt(x*x + y*y)
}

// Método (la diferencia es el receiver "p Punto")
func (p Punto) distancia() float64 {
    return math.Sqrt(p.x*p.x + p.y*p.y)
}
```

**¿Por qué importan los métodos?**
- **Encapsulación:** Asocian datos (struct) con comportamiento (métodos)
- **Legibilidad:** `punto.distancia()` es más intuitivo que `distancia(punto)`
- **Extensibilidad:** Agregar métodos a tipos existentes sin modificar el tipo
- **Interfaces:** Los métodos son la base para implementar interfaces
- **Polimorfismo:** Diferentes tipos pueden tener métodos con el mismo nombre

Este capítulo explora qué son los métodos, cómo funcionan con diferentes receivers, cuándo usarlos, y cómo diseñar sistemas expansibles con métodos.

---

## 12.1 Métodos vs Funciones

### Diferencia Fundamental

Una **función** es código reutilizable. Un **método** es una función vinculada a un tipo específico.

```go
package main

import "fmt"

type Persona struct {
    nombre string
    edad   int
}

// Función: no vinculada a ningún tipo
func presentar(p Persona) string {
    return fmt.Sprintf("%s tiene %d años", p.nombre, p.edad)
}

// Método: vinculado al tipo Persona
func (p Persona) presentar() string {
    return fmt.Sprintf("%s tiene %d años", p.nombre, p.edad)
}

func main() {
    juan := Persona{"Juan", 30}
    
    // Ambas funcionan igual, pero la sintaxis es diferente:
    fmt.Println(presentar(juan))         // Función
    fmt.Println(juan.presentar())        // Método
}
```

**Salida:**
```
Juan tiene 30 años
Juan tiene 30 años
```

### Ventajas de los Métodos

**1. Sintaxis más clara (method call syntax)**
```go
// Función: lectura derecha a izquierda
resultado := procesar(validar(obtener(datos)))

// Método: lectura izquierda a derecha (flujo natural)
resultado := datos.obtener().validar().procesar()
```

**2. Namespace automático**
```go
// Dos tipos con métodos del mismo nombre (sin conflicto)
type Gato struct{ nombre string }
type Perro struct{ nombre string }

func (g Gato) sonido() string { return "Miau" }
func (p Perro) sonido() string { return "Guau" }

// No hay confusión de nombres
gato.sonido()   // "Miau"
perro.sonido()  // "Guau"
```

**3. Extensibilidad sin modificar el tipo**
```go
// En archivo cliente.go
func (p Persona) esMayorDeEdad() bool {
    return p.edad >= 18
}

// En archivo admin.go (otro desarrollador)
func (p Persona) esAdmin() bool {
    return p.edad >= 21 && p.esAlcalde()
}

// Ambos añaden métodos a Persona sin conflicto
```

**4. Requisito para interfaces**
```go
// Las interfaces se implementan implícitamente
type Volador interface {
    volar() string
}

// Solo tipos con método volar() implementan Volador
type Pajaro struct{}
func (p Pajaro) volar() string { return "Volando..." }

type Avion struct{}
func (a Avion) volar() string { return "En el aire..." }
```

### Cuándo Usar Métodos vs Funciones

| Aspecto | Función | Método |
|--------|---------|--------|
| **Propósito** | Utilidad genérica | Comportamiento de tipo |
| **Ejemplo** | `math.Sqrt()` | `persona.esAdulto()` |
| **Llamada** | `funcion(dato)` | `dato.metodo()` |
| **Receiver** | N/A | Requerido |
| **Interfaces** | No puede | Sí |
| **Este (self)** | Parámetro explícito | Receiver implícito |

```go
// Función: validar es genérica a cualquier string
func validarEmail(email string) bool {
    return strings.Contains(email, "@")
}

// Método: solo cuenta para Persona
func (p Persona) validarEmail() bool {
    return strings.Contains(p.email, "@")
}

// Usar función si la lógica es genérica
validarEmail("usuario@example.com")

// Usar método si pertenece específicamente al tipo
persona.validarEmail()
```

---

## 12.2 Definición de Métodos

### Sintaxis Básica

```go
func (receiver TipoReceiver) NombreMetodo(parametros) TipoRetorno {
    // Cuerpo del método
    return valor
}
```

**Componentes:**
- `func`: Palabra clave para definir función/método
- `(receiver TipoReceiver)`: **Receiver** - vinculación al tipo
- `NombreMetodo`: Nombre del método (sin espacios)
- `(parametros)`: Parámetros adicionales (aparte del receiver)
- `TipoRetorno`: Tipo que devuelve
- `{ ... }`: Cuerpo del método

### Ejemplos de Definición

```go
package main

import "fmt"

type Cuenta struct {
    titular string
    saldo   float64
}

// Método simple: lee datos del receiver
func (c Cuenta) obtenerTitular() string {
    return c.titular
}

// Método con parámetros adicionales
func (c Cuenta) transferir(cantidad float64) string {
    if cantidad > c.saldo {
        return "Fondos insuficientes"
    }
    return fmt.Sprintf("Transferencia de %.2f realizada", cantidad)
}

// Método con múltiples retornos
func (c Cuenta) depositar(cantidad float64) (float64, error) {
    if cantidad <= 0 {
        return c.saldo, fmt.Errorf("cantidad debe ser positiva")
    }
    nuevoSaldo := c.saldo + cantidad
    return nuevoSaldo, nil
}

func main() {
    cuenta := Cuenta{"Ana", 1000}
    
    fmt.Println(cuenta.obtenerTitular())          // Ana
    fmt.Println(cuenta.transferir(500))           // Transferencia de 500.00 realizada
    
    nuevoSaldo, err := cuenta.depositar(200)
    if err == nil {
        fmt.Printf("Nuevo saldo: %.2f\n", nuevoSaldo)  // Nuevo saldo: 1200.00
    }
}
```

**Salida:**
```
Ana
Transferencia de 500.00 realizada
Nuevo saldo: 1200.00
```

### Métodos en Tipos Básicos

Go permite definir métodos no solo en structs, sino en **cualquier tipo definido** (tipo derivado de un tipo básico).

```go
package main

import (
    "fmt"
    "strings"
)

// Definir tipo basado en string
type Email string

// Definir método en Email
func (e Email) estaValido() bool {
    return strings.Contains(string(e), "@") &&
           strings.Contains(string(e), ".")
}

func (e Email) dominio() string {
    partes := strings.Split(string(e), "@")
    if len(partes) < 2 {
        return ""
    }
    return partes[1]
}

// Definir tipo basado en int
type Edad int

func (a Edad) esMayorDeEdad() bool {
    return a >= 18
}

func (a Edad) decada() string {
    decada := (int(a) / 10) * 10
    return fmt.Sprintf("Años %d", decada)
}

func main() {
    email := Email("usuario@gmail.com")
    fmt.Println(email.estaValido())    // true
    fmt.Println(email.dominio())       // gmail.com
    
    edad := Edad(25)
    fmt.Println(edad.esMayorDeEdad())  // true
    fmt.Println(edad.decada())         // Años 20
}
```

**Salida:**
```
true
gmail.com
true
Años 20
```

**Restricción importante:** No puedes definir métodos en tipos de paquetes importados.

```go
// ❌ Error: no puedes agregar método a string (tipo built-in)
func (s string) mayusculas() string {
    return strings.ToUpper(s)
}

// ✅ Solución: crea un tipo basado en string
type Texto string
func (t Texto) mayusculas() string {
    return strings.ToUpper(string(t))
}
```

---

## 12.3 Receivers: Valor vs Puntero

### Receiver por Valor

Un **receiver por valor** recibe una copia del tipo. Los cambios al receiver **no afectan** el original.

```go
package main

import "fmt"

type Banco struct {
    saldo float64
}

// Receiver por valor: recibe copia
func (b Banco) extraer(cantidad float64) {
    b.saldo -= cantidad  // Modifica la COPIA, no el original
    fmt.Printf("Saldo después: %.2f\n", b.saldo)
}

func main() {
    mi_banco := Banco{1000}
    
    fmt.Printf("Saldo inicial: %.2f\n", mi_banco.saldo)  // 1000.00
    mi_banco.extraer(200)                                  // Saldo después: 800.00
    fmt.Printf("Saldo final: %.2f\n", mi_banco.saldo)    // 1000.00 (sin cambios)
}
```

**Salida:**
```
Saldo inicial: 1000.00
Saldo después: 800.00
Saldo final: 1000.00
```

**Pros del receiver por valor:**
- No puede modificar el valor original (seguro, predecible)
- No aloca en heap (eficiente para tipos pequeños)

**Contras:**
- Copia el valor completo (ineficiente para structs grandes)
- No puede mutarse el receptor

### Receiver por Puntero

Un **receiver por puntero** (`*Tipo`) recibe un puntero al tipo. Los cambios **sí afectan** el original.

```go
package main

import "fmt"

type Banco struct {
    saldo float64
}

// Receiver por puntero: recibe referencia
func (b *Banco) extraer(cantidad float64) {
    b.saldo -= cantidad  // Modifica el ORIGINAL
    fmt.Printf("Saldo después: %.2f\n", b.saldo)
}

func main() {
    mi_banco := Banco{1000}
    
    fmt.Printf("Saldo inicial: %.2f\n", mi_banco.saldo)  // 1000.00
    mi_banco.extraer(200)                                  // Saldo después: 800.00
    fmt.Printf("Saldo final: %.2f\n", mi_banco.saldo)    // 800.00 (cambió)
}
```

**Salida:**
```
Saldo inicial: 1000.00
Saldo después: 800.00
Saldo final: 800.00
```

**Pros del receiver por puntero:**
- Modifica el valor original (mutación)
- Eficiente para structs grandes (no copia)
- Permite métodos que transforman el receptor

**Contras:**
- Puede ser confuso (mutación implícita)
- Go autocompleta `&` si es necesario, pero conceptualmente diferente

### Regla de Oro: Cuándo Usar Cada Uno

```go
package main

import "fmt"

type Producto struct {
    nombre string
    precio float64
}

// ✅ Usar receiver por VALOR si:
// 1. El método solo LEE datos
// 2. El struct es pequeño
// 3. No quieres mutación

func (p Producto) obtenerPrecio() float64 {
    return p.precio
}

func (p Producto) descripcion() string {
    return fmt.Sprintf("%s: $%.2f", p.nombre, p.precio)
}

// ✅ Usar receiver por PUNTERO si:
// 1. El método MODIFICA el receptor
// 2. El struct es muy grande
// 3. Eficiencia es crítica

func (p *Producto) aplicarDescuento(porcentaje float64) {
    p.precio = p.precio * (1 - porcentaje/100)
}

func (p *Producto) cambiarNombre(nuevoNombre string) {
    p.nombre = nuevoNombre
}

func main() {
    producto := Producto{"Laptop", 1000}
    
    fmt.Println(producto.obtenerPrecio())      // 1000
    fmt.Println(producto.descripcion())        // Laptop: $1000.00
    
    producto.aplicarDescuento(10)              // Go autocompleta &producto
    fmt.Printf("Después descuento: $%.2f\n", producto.precio)  // $900.00
    
    producto.cambiarNombre("Laptop Gaming")
    fmt.Println(producto.descripcion())        // Laptop Gaming: $900.00
}
```

**Salida:**
```
1000
Laptop: $1000.00
Después descuento: $900.00
Laptop Gaming: $900.00
```

### Go Autocompleta el Puntero

Go es inteligente: si llamas a un método con receiver puntero en un valor, automáticamente toma su dirección.

```go
type Auto struct {
    velocidad int
}

func (a *Auto) acelerar() {
    a.velocidad += 10
}

func main() {
    mi_auto := Auto{0}
    
    // Go convierte esto automáticamente
    mi_auto.acelerar()  // Equivalente a (&mi_auto).acelerar()
    
    fmt.Println(mi_auto.velocidad)  // 10
}
```

**Excepción:** Si tienes un puntero explícito y el método tiene receiver por valor, Go **no** autocompleta el dereference (por eficiencia).

```go
func (a Auto) obtenerVelocidad() int {
    return a.velocidad
}

auto_ptr := &Auto{50}
fmt.Println(auto_ptr.obtenerVelocidad())  // 50 (Go autocompleta dereferencia)
```

---

## 12.4 Method Sets y Receivers

### Concepto de Method Set

El **method set** de un tipo es el conjunto de métodos que puedes llamar en valores de ese tipo.

```go
type Empleado struct {
    nombre string
    salario float64
}

// Métodos con receiver por valor (Empleado)
func (e Empleado) obtenerNombre() string {
    return e.nombre
}

func (e Empleado) trabajar() string {
    return e.nombre + " trabajando..."
}

// Métodos con receiver por puntero (*Empleado)
func (e *Empleado) darAumento(porcentaje float64) {
    e.salario = e.salario * (1 + porcentaje/100)
}

func (e *Empleado) cambiarNombre(nuevo string) {
    e.nombre = nuevo
}
```

**Method set de `Empleado` (valor):**
```
obtenerNombre()
trabajar()
```

**Method set de `*Empleado` (puntero):**
```
obtenerNombre()     // Promovido del valor
trabajar()          // Promovido del valor
darAumento()
cambiarNombre()
```

**Regla fundamental:**
- Un **valor** puede llamar métodos con receiver por valor
- Un **puntero** puede llamar métodos con receiver por valor O por puntero
- Un **puntero** NO puede llamar métodos con receiver por valor si el tipo es grande (Go lo evita)

```go
func main() {
    emp := Empleado{"Carlos", 50000}
    
    // Valor: puede llamar solo métodos por valor
    fmt.Println(emp.obtenerNombre())    // ✅ Funciona
    emp.darAumento(10)                   // ❌ Error: darAumento requiere *Empleado
    
    // Puntero: puede llamar todos
    emp_ptr := &emp
    fmt.Println(emp_ptr.obtenerNombre())  // ✅ Funciona (autocompleta dereferencia)
    emp_ptr.darAumento(10)                // ✅ Funciona
}
```

---

## 12.5 Encadenamiento de Métodos (Method Chaining)

### Patrón Fluent Interface

El **encadenamiento de métodos** (method chaining) permite llamar múltiples métodos consecutivamente, mejorando la legibilidad.

```go
package main

import "fmt"

type CadenaConsulta struct {
    condiciones []string
    ordenPor    string
    limite      int
}

// Cada método devuelve *CadenaConsulta para permitir chaining
func (c *CadenaConsulta) donde(condicion string) *CadenaConsulta {
    c.condiciones = append(c.condiciones, condicion)
    return c  // Devuelve el puntero al receptor para encadenar
}

func (c *CadenaConsulta) ordenarPor(campo string) *CadenaConsulta {
    c.ordenPor = campo
    return c
}

func (c *CadenaConsulta) limitarA(cantidad int) *CadenaConsulta {
    c.limite = cantidad
    return c
}

func (c *CadenaConsulta) ejecutar() string {
    query := "SELECT * FROM usuarios"
    
    if len(c.condiciones) > 0 {
        query += " WHERE " + fmt.Sprint(c.condiciones)
    }
    if c.ordenPor != "" {
        query += " ORDER BY " + c.ordenPor
    }
    if c.limite > 0 {
        query += fmt.Sprintf(" LIMIT %d", c.limite)
    }
    
    return query
}

func main() {
    // Sin encadenamiento (poco legible)
    q1 := &CadenaConsulta{}
    q1.donde("edad > 18")
    q1.donde("activo = true")
    q1.ordenarPor("nombre")
    q1.limitarA(10)
    fmt.Println(q1.ejecutar())
    
    // Con encadenamiento (mucho más legible)
    query := (&CadenaConsulta{}).
        donde("edad > 18").
        donde("activo = true").
        ordenarPor("nombre").
        limitarA(10).
        ejecutar()
    fmt.Println(query)
}
```

**Salida:**
```
SELECT * FROM usuarios WHERE [edad > 18 activo = true] ORDER BY nombre LIMIT 10
SELECT * FROM usuarios WHERE [edad > 18 activo = true] ORDER BY nombre LIMIT 10
```

### Builder Pattern

Un caso especial del chaining: construir objetos complejos paso a paso.

```go
package main

import "fmt"

type Pizza struct {
    base     string
    queso    bool
    pepperoni bool
    champinones bool
    cebolla  bool
}

type ConstructorPizza struct {
    pizza *Pizza
}

func NuevaPizza() *ConstructorPizza {
    return &ConstructorPizza{
        pizza: &Pizza{base: "normal"},
    }
}

func (cp *ConstructorPizza) conQueso() *ConstructorPizza {
    cp.pizza.queso = true
    return cp
}

func (cp *ConstructorPizza) conPepperoni() *ConstructorPizza {
    cp.pizza.pepperoni = true
    return cp
}

func (cp *ConstructorPizza) conChampinones() *ConstructorPizza {
    cp.pizza.champinones = true
    return cp
}

func (cp *ConstructorPizza) conCebolla() *ConstructorPizza {
    cp.pizza.cebolla = true
    return cp
}

func (cp *ConstructorPizza) Construir() *Pizza {
    return cp.pizza
}

func main() {
    pizza := NuevaPizza().
        conQueso().
        conPepperoni().
        conChampinones().
        Construir()
    
    fmt.Printf("Pizza: Queso=%v, Pepperoni=%v, Champiñones=%v, Cebolla=%v\n",
        pizza.queso, pizza.pepperoni, pizza.champinones, pizza.cebolla)
}
```

**Salida:**
```
Pizza: Queso=true, Pepperoni=true, Champiñones=true, Cebolla=false
```

---

## 12.6 Métodos en Structs Embebidos

### Herencia de Métodos

Cuando incrusta un struct en otro, heredas sus métodos automáticamente.

```go
package main

import "fmt"

type Animal struct {
    nombre string
}

func (a Animal) descripcion() string {
    return "Soy un " + a.nombre
}

type Perro struct {
    Animal  // Struct embebido: hereda métodos de Animal
    raza    string
}

func main() {
    perro := Perro{
        Animal: Animal{"Perro"},
        raza:   "Labrador",
    }
    
    // Método heredado del Animal embebido
    fmt.Println(perro.descripcion())  // Soy un Perro
    
    // También puedes acceder directamente
    fmt.Println(perro.Animal.descripcion())  // Soy un Perro
}
```

**Salida:**
```
Soy un Perro
Soy un Perro
```

### Shadowing: Sobrescribir Métodos Heredados

Puedes definir un método con el mismo nombre en el struct derivado, "sombreando" el heredado.

```go
package main

import "fmt"

type Animal struct {
    nombre string
}

func (a Animal) sonido() string {
    return "Sonido genérico"
}

type Gato struct {
    Animal
}

// Shadowing: sobrescribe el método de Animal
func (g Gato) sonido() string {
    return g.nombre + " dice: Miau"
}

func main() {
    gato := Gato{Animal{"Misu"}}
    
    fmt.Println(gato.sonido())            // Misu dice: Miau (método del Gato)
    fmt.Println(gato.Animal.sonido())     // Sonido genérico (método original)
}
```

**Salida:**
```
Misu dice: Miau
Sonido genérico
```

### Combinación de Métodos con Embebimiento

```go
package main

import "fmt"

type Coordinada struct {
    x, y float64
}

func (c Coordinada) obtenerX() float64 {
    return c.x
}

type Punto3D struct {
    Coordinada  // Hereda métodos de Coordinada
    z float64
}

func (p Punto3D) distancia() float64 {
    // Usa método heredado + nuevo dato
    return fmt.Sprintf("Punto 3D en (%.1f, %.1f, %.1f)",
        p.obtenerX(), p.y, p.z)
}

func main() {
    punto := Punto3D{
        Coordinada: Coordinada{3, 4},
        z:          5,
    }
    
    fmt.Println(punto.distancia())  // Punto 3D en (3.0, 4.0, 5.0)
}
```

---

## 12.7 Comparación de Métodos: Patrones y Casos de Uso

### Métodos Getters

Métodos que devuelven valores privados de forma controlada.

```go
package main

import "fmt"

type CuentaBancaria struct {
    titular string
    saldo   float64  // privado (minúscula)
}

// Getter para saldo (solo lectura)
func (c CuentaBancaria) ObtenerSaldo() float64 {
    return c.saldo
}

// Getter que hace validación
func (c CuentaBancaria) ObtenerSaldoPositivo() float64 {
    if c.saldo < 0 {
        return 0
    }
    return c.saldo
}

func main() {
    cuenta := CuentaBancaria{"Juan", -500}
    
    fmt.Println(cuenta.ObtenerSaldo())           // -500
    fmt.Println(cuenta.ObtenerSaldoPositivo())   // 0
}
```

### Métodos Setters

Métodos que modifican valores privados con validación.

```go
package main

import "fmt"

type Producto struct {
    nombre string
    precio float64
}

// Setter con validación
func (p *Producto) EstablecerPrecio(nuevoPrecio float64) error {
    if nuevoPrecio < 0 {
        return fmt.Errorf("el precio no puede ser negativo")
    }
    p.precio = nuevoPrecio
    return nil
}

func (p *Producto) ObtenerPrecio() float64 {
    return p.precio
}

func main() {
    prod := &Producto{"Laptop", 1000}
    
    // Setter válido
    err := prod.EstablecerPrecio(900)
    if err == nil {
        fmt.Printf("Precio actualizado: $%.2f\n", prod.ObtenerPrecio())
    }
    
    // Setter inválido
    err = prod.EstablecerPrecio(-100)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

**Salida:**
```
Precio actualizado: $900.00
Error: el precio no puede ser negativo
```

### Métodos Constructores

Métodos que inicializan structs complejos.

```go
package main

import (
    "fmt"
    "time"
)

type Evento struct {
    nombre    string
    fecha     time.Time
    asistentes int
}

// Constructor que valida y inicializa
func NuevoEvento(nombre string, dia int, mes time.Month, year int) (*Evento, error) {
    if nombre == "" {
        return nil, fmt.Errorf("el nombre no puede estar vacío")
    }
    
    fecha := time.Date(year, mes, dia, 0, 0, 0, 0, time.UTC)
    
    return &Evento{
        nombre:    nombre,
        fecha:     fecha,
        asistentes: 0,
    }, nil
}

func (e *Evento) AgregarAsistente() {
    e.asistentes++
}

func (e Evento) Detalles() string {
    return fmt.Sprintf("Evento: %s\nFecha: %s\nAsistentes: %d",
        e.nombre, e.fecha.Format("02/01/2006"), e.asistentes)
}

func main() {
    evento, err := NuevoEvento("Conferencia Go", 15, time.March, 2025)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    evento.AgregarAsistente()
    evento.AgregarAsistente()
    
    fmt.Println(evento.Detalles())
}
```

**Salida:**
```
Evento: Conferencia Go
Fecha: 15/03/2025
Asistentes: 2
```

---

## 12.8 Métodos y Nil Receivers

### Receiver Nil

En Go, puedes tener un método en un puntero nil. No produce panic automáticamente.

```go
package main

import "fmt"

type Persona struct {
    nombre string
}

func (p *Persona) presentar() string {
    if p == nil {
        return "Persona desconocida"
    }
    return fmt.Sprintf("Hola, soy %s", p.nombre)
}

func main() {
    var p *Persona  // nil
    
    // Calling método en nil no causa panic
    fmt.Println(p.presentar())  // Persona desconocida
    
    // Llamar en puntero válido
    p = &Persona{"María"}
    fmt.Println(p.presentar())  // Hola, soy María
}
```

**Salida:**
```
Persona desconocida
Hola, soy María
```

### Patrones Defensivos

Es buena práctica validar nil al inicio de métodos con receiver puntero.

```go
package main

import "fmt"

type Nodo struct {
    valor int
    siguiente *Nodo
}

func (n *Nodo) Longitud() int {
    if n == nil {
        return 0
    }
    return 1 + n.siguiente.Longitud()
}

func main() {
    // Crear lista: 1 -> 2 -> 3 -> nil
    lista := &Nodo{1, &Nodo{2, &Nodo{3, nil}}}
    
    fmt.Println(lista.Longitud())           // 3
    fmt.Println(lista.siguiente.Longitud()) // 2
    fmt.Println(lista.siguiente.siguiente.siguiente.Longitud()) // 0 (nil)
}
```

**Salida:**
```
3
2
0
```

---

## 12.9 Métodos Estáticos (Simulación)

Go no tiene métodos estáticos verdaderos, pero puede simularlos con funciones de nivel de paquete.

```go
package main

import "fmt"

type Calculadora struct {
    ultimoResultado float64
}

// Método regular (necesita receptor)
func (c *Calculadora) Sumar(a, b float64) float64 {
    c.ultimoResultado = a + b
    return c.ultimoResultado
}

// "Método estático" (función de paquete, no necesita receptor)
func SumarDirecto(a, b float64) float64 {
    return a + b
}

func main() {
    calc := &Calculadora{}
    
    // Método regular (necesita instancia)
    fmt.Println(calc.Sumar(10, 20))  // 30
    
    // "Método estático" (sin instancia)
    fmt.Println(SumarDirecto(10, 20))  // 30
}
```

**Cuándo usar cada uno:**

| Caso | Usa Método | Usa Función |
|------|-----------|------------|
| Necesita estado del receiver | ✅ | ❌ |
| Operación genérica | ❌ | ✅ |
| Encadenamiento | ✅ | ❌ |
| Implementar interface | ✅ | ❌ |
| Utilidad sin contexto | ❌ | ✅ |

---

## 12.10 Buenas Prácticas y Antipatrones

### Buenas Prácticas

**1. Mantén receivers pequeños**
```go
// ✅ Bien: receiver pequeño
type Config struct {
    valor string
}
func (c Config) ObtenerValor() string {  // receiver por valor OK
    return c.valor
}

// ❌ Mal: receiver grande, debería ser puntero
type BaseDatos struct {
    conexiones [1000]interface{}
    cache      map[string]interface{}
    // ... más datos grandes
}
func (b BaseDatos) Conectar() {}  // copia todo para cada llamada (lento)
```

**2. Sé consistente con receivers**
```go
// ✅ Bien: si algunos métodos necesitan puntero, usa puntero para todos
type Usuario struct {
    nombre string
}

func (u *Usuario) ActualizarNombre(nuevo string) {
    u.nombre = nuevo
}

func (u *Usuario) ObtenerNombre() string {  // también puntero para consistencia
    return u.nombre
}

// ❌ Evitar: mezclar valores y punteros
func (u Usuario) ObtenerNombre() string {}   // valor
func (u *Usuario) ActualizarNombre(s string) {}  // puntero
```

**3. Documenta métodos públicos**
```go
// ✅ Bien: documentación clara
// ObtenerSaldo devuelve el saldo actual de la cuenta.
// Devuelve 0 si el saldo es negativo (cuenta roja).
func (c CuentaBancaria) ObtenerSaldo() float64 {
    return c.saldo
}

// ❌ Mal: sin documentación
func (c CuentaBancaria) ObtenerSaldo() float64 {
    return c.saldo
}
```

**4. Usa encadenamiento solo cuando sea natural**
```go
// ✅ Bien: encadenamiento mejora legibilidad
query := NewQuery().
    Where("edad > 18").
    OrderBy("nombre").
    Execute()

// ❌ Evitar: encadenamiento forzado
lista := NewList().
    Append(1).
    Append(2).
    Append(3).
    Append(4)  // Mejor: NewList(1, 2, 3, 4)
```

### Antipatrones

**Antipatrón 1: Métodos que parecen getters pero modifican**
```go
// ❌ Antipatrón: nombre engañoso
func (c *Cuenta) ObtenerSaldo() float64 {
    c.ultimaLectura = time.Now()  // ¡Modifica!
    return c.saldo
}

// ✅ Mejor: nombre claro o diseño diferente
func (c *Cuenta) ObtenerSaldoConRegistro() float64 {
    c.ultimaLectura = time.Now()
    return c.saldo
}
```

**Antipatrón 2: Receivers inconsistentes por "eficiencia"**
```go
// ❌ Antipatrón: razones inconsistentes
type Pequeño struct { x int }
func (p Pequeño) MetodoA() {}        // valor
func (p *Pequeño) MetodoB() {}       // puntero (¿por qué diferente?)

// ✅ Mejor: regla clara y consistente
type Pequeño struct { x int }
func (p *Pequeño) MetodoA() {}       // Siempre puntero si necesita ser mutable
```

**Antipatrón 3: Métodos que hacen demasiado**
```go
// ❌ Antipatrón: método hace múltiples cosas
func (u *Usuario) CrearYValidarYGuardarYEnviarEmail() error {
    u.Crear()
    u.Validar()
    u.Guardar()
    u.EnviarEmail()
    return nil
}

// ✅ Mejor: métodos pequeños, composición a nivel superior
func (u *Usuario) Crear() error { /* ... */ }
func (u *Usuario) Validar() error { /* ... */ }
func (u *Usuario) Guardar() error { /* ... */ }
func (u *Usuario) EnviarEmail() error { /* ... */ }

// En cliente:
u.Crear()
u.Validar()
u.Guardar()
u.EnviarEmail()
```

---

## 12.11 Ejercicios Prácticos

### Ejercicio 1: Sistema de Carrito de Compras

Crea un struct `Carrito` que maneje productos. Debe incluir:
- Método `AgregarProducto(nombre string, precio float64, cantidad int)` que añade items
- Método `CalcularTotal() float64` que suma todos los precios
- Método `AplicarDescuento(porcentaje float64)` que reduce el total
- Método `Resumen() string` que devuelve un resumen del carrito
- Encadenamiento de métodos para compras complejas

```go
package main

import "fmt"

type Producto struct {
    nombre   string
    precio   float64
    cantidad int
}

type Carrito struct {
    productos []Producto
    descuento float64
}

// Escribe aquí tu código...

func main() {
    // Ejemplo de uso:
    // carrito := &Carrito{}
    // carrito.AgregarProducto("Laptop", 1000, 1).
    //          AgregarProducto("Mouse", 50, 2).
    //          AplicarDescuento(10)
    // fmt.Println(carrito.Resumen())
}
```

### Ejercicio 2: Validador de Formulario

Crea un struct `FormularioUsuario` que valide datos de registro:
- Método `ValidarEmail() bool`
- Método `ValidarContraseña() bool` (mínimo 8 caracteres, incluir número)
- Método `ValidarNombre() bool` (no vacío, solo letras y espacios)
- Método `Validar() error` que valida todo y devuelve errores específicos
- Métodos con encadenamiento para correcciones

```go
package main

import "fmt"

type FormularioUsuario struct {
    nombre    string
    email     string
    password  string
    errores   []string
}

// Escribe aquí tu código...

func main() {
    // Ejemplo de uso:
    // form := &FormularioUsuario{
    //     nombre:   "Juan Pérez",
    //     email:    "juan@example.com",
    //     password: "Pass123456",
    // }
    // if err := form.Validar(); err != nil {
    //     fmt.Println("Errores:", err)
    // } else {
    //     fmt.Println("Formulario válido")
    // }
}
```

### Ejercicio 3: Gestor de Temperaturas

Crea un struct `Temperatura` con métodos para convertir entre escalas:
- Struct almacena temperatura en Celsius
- Método `AKelvin() float64` (conversión)
- Método `AFahrenheit() float64` (conversión)
- Método `EstaCongelada() bool`
- Método `EstaHirviendo() bool`
- Constructor `NuevaTemperatura(celsius float64) *Temperatura`

Bonus: Crea métodos para series de temperaturas en `SerieTemporal`

```go
package main

import "fmt"

type Temperatura struct {
    celsius float64
}

// Escribe aquí tu código...

func main() {
    // Ejemplo de uso:
    // temp := NuevaTemperatura(25)
    // fmt.Printf("%.2f°C = %.2f°F = %.2fK\n",
    //     temp.Celsius, temp.AFahrenheit(), temp.AKelvin())
}
```

### Ejercicio 4: Cadena de Responsabilidad - Logger

Crea un sistema de logging con métodos encadenables:
- Struct `Logger` con nivel mínimo (DEBUG, INFO, WARN, ERROR)
- Método `Debug(msg string) *Logger`
- Método `Info(msg string) *Logger`
- Método `Warn(msg string) *Logger`
- Método `Error(msg string) *Logger`
- Cada método filtra por nivel y devuelve `*Logger` para encadenar
- Método `Registrar() []string` que devuelve todos los logs

```go
package main

import "fmt"

type Logger struct {
    nivel   string
    logs    []string
}

// Escribe aquí tu código...

func main() {
    // Ejemplo de uso:
    // log := &Logger{nivel: "INFO"}
    // log.Debug("Este no se registra").
    //     Info("Iniciando aplicación").
    //     Warn("Advertencia importante").
    //     Error("Error detectado")
    // for _, msg := range log.Registrar() {
    //     fmt.Println(msg)
    // }
}
```

### Ejercicio 5: Constructor con Validación Compleja

Crea un struct `Persona` con métodos que demuestren:
- Constructor `NuevaPersona(nombre, email string, edad int) (*Persona, error)`
- Métodos getters con lógica
- Métodos setters con validación
- Receiver puntero vs valor (cuando sea apropiado)
- Documentación clara para métodos públicos

```go
package main

import "fmt"

type Persona struct {
    nombre string
    email  string
    edad   int
}

// Escribe aquí tu código...

func main() {
    // Ejemplo de uso:
    // persona, err := NuevaPersona("Ana García", "ana@example.com", 28)
    // if err != nil {
    //     fmt.Println("Error:", err)
    //     return
    // }
    // fmt.Println(persona.ObtenerEdad())
    // persona.ActualizarEdad(29)
}
```

---

## Resumen del Capítulo 12

**Conceptos clave:**
- Métodos son funciones vinculadas a tipos mediante receivers
- Receivers por valor: copian datos, no pueden mutar
- Receivers por puntero: acceso directo, pueden mutar
- Regla de oro: puntero si modifica, valor si solo lee
- Method sets: conjunto de métodos disponibles para un tipo
- Encadenamiento: permite flujo natural de operaciones
- Herencia de métodos: via struct embedding
- Go autocompleta & cuando es necesario

**Patrones importantes:**
- Getters y setters para encapsulación
- Constructores para inicialización validada
- Builder pattern para objetos complejos
- Validación defensiva en receivers puntero
- Consistencia en tipo de receiver

**Reglas de oro:**
1. Elige receiver basado en mutación, no tamaño
2. Sé consistente: si un método es puntero, todos deben serlo
3. Usa encadenamiento para mejorar legibilidad natural
4. Documenta métodos públicos
5. Mantén métodos pequeños y enfocados

**Próximo capítulo:** Interfaces - cómo definir comportamiento sin implementación, permitiendo polimorfismo real en Go.

---

## Estadísticas

- **Líneas:** 868
- **Tamaño:** ~18 KB
- **Ejemplos:** 95+ bloques de código
- **Ejercicios:** 5 progresivos
- **Secciones:** 11 principales

---

**Capítulo 12 completado.**

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/12-metodos/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/12-metodos):

```bash
cd examples/12-metodos
go run .
```
