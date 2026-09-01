# Capítulo 14: Embedding y composición avanzada

## Introducción

Go rechaza deliberadamente la herencia como mecanismo de reutilización de código. En su lugar, proporciona **embedding** (incrustación) de structs e interfaces como el principal patrón de composición. Este capítulo explora en profundidad cómo funciona el embedding, cómo resolver conflictos de nombres, cuándo usar embedding versus delegación explícita, y cómo implementar patrones arquitectónicos sofisticados basados en composición.

El embedding es más que syntax sugar: es una expresión de la filosofía de diseño de Go de favorecer la composición sobre la herencia, lo que resulta en código más flexible, mantenible y explícito.

---

## 14.1 Composición vs Herencia: Fundamentos Filosóficos

### 14.1.1 Por qué Go Rechaza la Herencia

La herencia clásica (de lenguajes orientados a objetos como Java o C++) presenta varios problemas fundamentales:

#### El Problema del Triángulo Mortal (Fragile Base Class Problem)

```
Lenguajes con herencia:

Padre
  ├─ métodos
  ├─ campos privados
  └─ detalles de implementación

     ↓ (herencia depende de detalles internos)

Hijo
  ├─ Debe respetar contrato de Padre
  ├─ Puede romperse si Padre cambia
  └─ Acoplamiento fuerte
```

En Java:

```java
// Cambio aparentemente inocuo en Padre
class Padre {
    protected int contador = 0;
    
    protected void incrementar() {
        contador++;
    }
}

// Hijo rompe si Padre cambia la lógica
class Hijo extends Padre {
    public void operacion() {
        incrementar(); // ¿Qué es exactamente "incrementar"?
    }
}
```

#### El Diamante de la Muerte (Diamond Problem)

```
      A
     / \
    B   C
     \ /
      D

¿Qué versión de método_x usa D? ¿De B, C, o A?
```

Go evita esto completamente: sin herencia = sin problema del diamante.

#### Jerarquías Profundas y Rigidez

Las jerarquías de herencia crean estructuras inflexibles:

```
Animal
  ├─ Mamífero
  │   ├─ Carnívoro
  │   │   ├─ Felino
  │   │   │   ├─ Gato
  │   │   │   └─ Tigre
  │   │   └─ Canino
  │   └─ Herbívoro
  └─ Ave
```

Si descubres que necesitas que `Gato` también sea `Doméstico`, toda la jerarquía sufre.

### 14.1.2 La Solución de Go: Composición y Interfaces

Go resolve esto con dos mecanismos simples:

**1. Composición: Reutilizar comportamiento embebiendo structs**

```go
// En lugar de: class Gato extends Felino extends Carnívoro extends Mamífero
// Go hace: Componer comportamientos específicos

type Mamífero struct {
    nombre string
    peso   float64
}

func (m *Mamífero) Respirar() {
    fmt.Println("Mamífero respira aire")
}

type Gato struct {
    Mamífero  // Embedding: Gato TIENE UN Mamífero
    rasguños  int
}

func (g *Gato) Maullar() {
    fmt.Println("Miau!")
}

// Gato automáticamente hereda los métodos de Mamífero
gato := &Gato{Mamífero: Mamífero{nombre: "Misi", peso: 4.5}}
gato.Respirar() // Funciona: Mamífero respira aire
gato.Maullar()  // Funciona: Miau!
```

**2. Interfaces: Polimorfismo sin jerarquía**

```go
// En lugar de una jerarquía Drawable -> Shape -> Circle
// Go usa interfaces basadas en métodos

type Dibujable interface {
    Dibujar()
    ObtenerÁrea() float64
}

type Círculo struct {
    radio float64
}

func (c Círculo) Dibujar() {
    fmt.Printf("Círculo de radio %.2f\n", c.radio)
}

func (c Círculo) ObtenerÁrea() float64 {
    return math.Pi * c.radio * c.radio
}

// Círculo satisface Dibujable sin declaración explícita
// Cualquier otro tipo que implemente los métodos también satisface Dibujable
```

### 14.1.3 Ventajas de la Composición

| Aspecto | Herencia | Composición (Go) |
|--------|----------|------------------|
| **Acoplamiento** | Fuerte (depende de internals) | Débil (contratos claros) |
| **Flexibilidad** | Rígida (estructura fija) | Flexible (combina libremente) |
| **Debugging** | Complejo (múltiples niveles) | Simple (responsabilidades claras) |
| **Reutilización** | Limitada a jerarquía | Mediante interfaces y composición |
| **Cambios** | Afecta toda la jerarquía | Aislados a componentes específicos |

### 14.1.4 Cuándo Usar Cada Patrón

```go
// ✅ Usa composición de STRUCTS cuando:
// - Quieres reutilizar código de otro struct
// - La relación es "TIENE UN" (has-a)

type Motor struct{ potencia int }
type Auto struct{
    Motor    // Auto TIENE UN Motor
    color string
}

// ✅ Usa composición de INTERFACES cuando:
// - Quieres polimorfismo
// - Múltiples tipos deben cumplir un contrato
// - La relación es "PUEDE SER" (is-a)

type Vehículo interface {
    Conducir()
    Frenar()
}

// ✅ Usa delegación explícita cuando:
// - Quieres control total sobre qué métodos exponer
// - Deseas agregar lógica adicional

type Proxy struct {
    real RealObject
}

func (p Proxy) Operación() {
    log.Println("Antes de la operación")
    p.real.Operación()
    log.Println("Después de la operación")
}
```

---

## 14.2 Embedding Básico: Campos Anónimos y Promoción

### 14.2.1 Campos Anónimos (Anonymous Fields)

Un **campo anónimo** es un campo sin nombre explícito:

```go
type Persona struct {
    nombre string
    edad   int
}

type Empleado struct {
    Persona     // Campo anónimo: tipo es el nombre
    empleadoID  string
    salario     float64
}

// Acceso a campos embebidos:
emp := Empleado{
    Persona: Persona{nombre: "Juan", edad: 30},
    empleadoID: "EMP001",
    salario: 50000,
}

// Campo anónimo se accede por tipo:
fmt.Println(emp.Persona.nombre)  // "Juan"
fmt.Println(emp.empleadoID)      // "EMP001"

// O promovido directamente (ver siguiente sección):
fmt.Println(emp.nombre)           // "Juan" (si usamos promoción)
```

### 14.2.2 Promoción de Campos (Field Promotion)

Cuando embedbas un struct, los campos del struct embebido se **promueven** y accesibles directamente:

```go
type Dirección struct {
    calle  string
    ciudad string
    país   string
}

type Persona struct {
    nombre string
    Dirección  // Campo anónimo
}

p := Persona{
    nombre: "Ana",
    Dirección: Dirección{
        calle: "Avenida Principal 123",
        ciudad: "Madrid",
        país: "España",
    },
}

// Los campos de Dirección se promueven:
fmt.Println(p.calle)   // "Avenida Principal 123" (acceso promovido)
fmt.Println(p.ciudad)  // "Madrid"
fmt.Println(p.país)    // "España"

// Pero también puedes acceder explícitamente:
fmt.Println(p.Dirección.calle)  // "Avenida Principal 123"
```

**Regla de Promoción:**

```
Un campo F de struct embebido E se promueve si:
1. F no es un campo de inicio implícito de E
2. F no se redeclara en el struct que embebe a E

Esto significa:
- Los campos promovidos se tratan como si fueran del struct exterior
- Aparecen en el inicializador de struct
- Se consideran para comparación y asignación
```

### 14.2.3 Promoción de Métodos (Method Promotion)

No solo se promueven campos; también se promueven **métodos**:

```go
type Vehículo struct {
    marca   string
    modelo  string
}

func (v *Vehículo) Conducir() {
    fmt.Printf("Conduciendo %s %s\n", v.marca, v.modelo)
}

func (v *Vehículo) Frenar() {
    fmt.Printf("%s %s frenando\n", v.marca, v.modelo)
}

type Auto struct {
    Vehículo          // Embed Vehículo
    numeroPuertas int
}

// Los métodos de Vehículo se promueven a Auto
auto := &Auto{
    Vehículo: Vehículo{marca: "Toyota", modelo: "Camry"},
    numeroPuertas: 4,
}

auto.Conducir()  // Funciona: "Conduciendo Toyota Camry"
auto.Frenar()    // Funciona: "Toyota Camry frenando"

// Bajo el capó: Go propaga implícitamente la llamada a Vehículo
// auto.Conducir() es equivalente a auto.Vehículo.Conducir()
```

### 14.2.4 Limitaciones de Promoción

La promoción tiene reglas importantes:

```go
type A struct {
    campo int
}

type B struct {
    A
    campo int  // Conflicto: B tiene su propio "campo"
}

b := B{}
// b.campo se refiere a B.campo, NO a A.campo
// Para acceder a A.campo DEBE ser: b.A.campo

fmt.Println(b.campo)    // Accede a B.campo
fmt.Println(b.A.campo)  // Accede a A.campo explícitamente
```

**Limitación importante: Receivers de Puntero vs Valor**

```go
type Punto struct {
    x, y int
}

// Método con receiver de puntero
func (p *Punto) Mover(dx, dy int) {
    p.x += dx
    p.y += dy
}

type Punto3D struct {
    Punto
    z int
}

// ✅ Funciona: receiver de puntero
p3d := &Punto3D{}
p3d.Mover(1, 2)  // OK: *Punto3D → *Punto

// ❌ No funciona: receiver de valor
p3d_val := Punto3D{}
// p3d_val.Mover(1, 2)  // ERROR: no se puede pasar Punto3D a *Punto
```

---

## 14.3 Promoción de Métodos: Comportamiento y Limitaciones

### 14.3.1 Cómo Funciona la Promoción de Métodos

```go
type Logger struct {
    nivel string
}

func (l *Logger) Log(msg string) {
    fmt.Printf("[%s] %s\n", l.nivel, msg)
}

type Aplicación struct {
    Logger          // Embed Logger
    nombre string
}

app := &Aplicación{
    Logger: Logger{nivel: "INFO"},
    nombre: "MiApp",
}

// Llamada promovida: app.Log() → app.Logger.Log()
app.Log("Iniciando aplicación")  // [INFO] Iniciando aplicación

// En realidad, Go expande esto a:
// func (a *Aplicación) Log(msg string) {
//     a.Logger.Log(msg)
// }
```

### 14.3.2 Cadena de Promoción

Los métodos se promueven a través de múltiples niveles de embedding:

```go
type Base struct {
    valor int
}

func (b *Base) Obtener() int {
    return b.valor
}

type Intermedio struct {
    Base
}

type Exterior struct {
    Intermedio
}

ext := &Exterior{
    Intermedio: Intermedio{
        Base: Base{valor: 42},
    },
}

// El método Obtener se promueve a través de Intermedio a Exterior
ext.Obtener()  // Retorna 42
// Go expande: ext.Intermedio.Base.Obtener()
```

### 14.3.3 Cuándo NO se Promueven los Métodos

```go
type A struct{}

func (a A) Método() {}        // receiver de valor: se promueve
func (a *A) MétodoPtr() {}    // receiver de puntero: se promueve

type B struct {
    A              // embedded valor
}

// Métodos promovidos en B (valor):
// - Método() ✅ (puedo acceder a B.A.Método())
// - MétodoPtr() ❌ (B es valor, A* requiere puntero)

type C struct {
    *A             // embedded puntero
}

// Métodos promovidos en C (puntero):
// - Método() ✅ (puedo acceder a (*C.A).Método())
// - MétodoPtr() ✅ (C es puntero a A)

var b B
b.Método()      // OK: B.A.Método()
// b.MétodoPtr() // ERROR: B no tiene puntero a A

var c C
c.Método()      // OK: (*C.A).Método() - Go auto-desreferencia
c.MétodoPtr()   // OK: C.A.MétodoPtr()
```

---

## 14.4 Embedding Múltiple: Composición de Múltiples Tipos

### 14.4.1 Embebiendo Múltiples Structs

Go permite embeber múltiples structs en uno solo:

```go
type Coordenada struct {
    x, y float64
}

type Dimensión struct {
    ancho, alto float64
}

type Forma struct {
    nombre string
}

type Rectángulo struct {
    Coordenada    // Posición
    Dimensión     // Tamaño
    Forma         // Información general
    color   string
}

rect := &Rectángulo{
    Coordenada: Coordenada{x: 10, y: 20},
    Dimensión:  Dimensión{ancho: 100, alto: 50},
    Forma:      Forma{nombre: "Rectángulo A"},
    color:      "azul",
}

// Acceso a campos promovidos:
fmt.Println(rect.x)      // 10 (de Coordenada)
fmt.Println(rect.ancho)  // 100 (de Dimensión)
fmt.Println(rect.nombre) // "Rectángulo A" (de Forma)
fmt.Println(rect.color)  // "azul"
```

### 14.4.2 Resolución de Nombres en Embedding Múltiple

Cuando múltiples structs embebidos tienen campos con el mismo nombre:

```go
type Persona struct {
    nombre string
}

type Libro struct {
    nombre string
    autor  string
}

type Colección struct {
    Persona
    Libro
}

// Conflicto: Ambos Persona y Libro tienen "nombre"
col := Colección{}
// col.nombre  // ERROR: ambiguous selector nombre

// Debes usar específicamente:
fmt.Println(col.Persona.nombre)  // Campo de Persona
fmt.Println(col.Libro.nombre)    // Campo de Libro

// Sin embargo, campos ÚNICOS se promueven sin conflicto:
fmt.Println(col.autor)  // Único en Libro, se promueve OK
```

**Diagrama de Resolución:**

```
Colección
├─ Persona
│  └─ nombre ✗ (conflicto)
├─ Libro
│  ├─ nombre ✗ (conflicto)
│  └─ autor ✓ (único, promovido)
└─ acceso a col.nombre ❌ AMBIGUO
   acceso a col.autor ✅ OK
```

### 14.4.3 Resolución de Conflictos

Estrategias para resolver conflictos:

```go
type Trabajador struct {
    nombre     string
    departamento string
}

type Voluntario struct {
    nombre string
    organización string
}

type Persona struct {
    Trabajador
    Voluntario
    // No hay forma de promover "nombre" automáticamente
}

p := Persona{
    Trabajador: Trabajador{nombre: "Juan", departamento: "IT"},
    Voluntario: Voluntario{nombre: "Juan", organización: "Cruz Roja"},
}

// ❌ p.nombre          // Ambiguo
// ❌ p.departamento    // Ambiguo (podría estar en Trabajador)
// ✅ p.Trabajador.nombre       // Explícito
// ✅ p.Voluntario.nombre       // Explícito
// ✅ p.Trabajador.departamento // Explícito

// Solución 1: Redefine en Persona para ser explícito
type PersonaMejorada struct {
    Trabajador
    Voluntario
    nombre string  // Sombrea los campos embebidos
}

// Solución 2: Usa métodos en lugar de campos dirección
type PersonaConMétodos struct {
    Trabajador
    Voluntario
}

func (p *PersonaConMétodos) NombreTrabajador() string {
    return p.Trabajador.nombre
}

func (p *PersonaConMétodos) NombreVoluntario() string {
    return p.Voluntario.nombre
}
```

---

## 14.5 Ambigüedad en Embeddings: Identificación y Resolución

### 14.5.1 Cuándo Surge la Ambigüedad

La ambigüedad ocurre cuando:

1. **Campos con mismo nombre en múltiples embeddings**
2. **Métodos con mismo nombre en múltiples embeddings**
3. **Un campo sombrea un método o viceversa**

```go
type A struct {
    valor int
}

func (a *A) Método() string {
    return "A"
}

type B struct {
    valor int
}

func (b *B) Método() string {
    return "B"
}

type C struct {
    A
    B
}

c := &C{
    A: A{valor: 1},
    B: B{valor: 2},
}

// Ambigüedad 1: campo
// fmt.Println(c.valor)  // ERROR: ambiguous selector valor

// Ambigüedad 2: método
// c.Método()  // ERROR: ambiguous selector Método

// Resolución: explícito
fmt.Println(c.A.valor)    // 1
fmt.Println(c.B.valor)    // 2
fmt.Println(c.A.Método()) // "A"
fmt.Println(c.B.Método()) // "B"
```

### 14.5.2 Ambigüedad a Través de Múltiples Niveles

```go
type Base struct {
    dato int
}

type RamaIzquierda struct {
    Base
}

type RamaDerecha struct {
    Base
}

type Árbol struct {
    RamaIzquierda
    RamaDerecha
}

árbol := Árbol{
    RamaIzquierda: RamaIzquierda{Base: Base{dato: 1}},
    RamaDerecha:   RamaDerecha{Base: Base{dato: 2}},
}

// ábol.dato // ERROR: ambiguous
// Debe ser: árbol.RamaIzquierda.Base.dato o árbol.RamaDerecha.Base.dato

// Diagrama del problema (similar al Diamond Problem):
//         Base
//         /  \
//  RamaIzq    RamaDer
//         \  /
//          Árbol
```

### 14.5.3 Patrones para Evitar Ambigüedad

```go
// ❌ Malo: Ambigüedad potencial
type Malo struct {
    base1 Base
    base2 Base
}

// ✅ Bueno: Usar campos nombrados explícitamente
type Bueno struct {
    izquierda Base
    derecha   Base
}

// ✅ Bueno: Usar métodos para acceso
type Base struct {
    valor int
}

type Contenedor struct {
    primaria   Base
    secundaria Base
}

func (c Contenedor) ValorPrimario() int {
    return c.primaria.valor
}

func (c Contenedor) ValorSecundario() int {
    return c.secundaria.valor
}

// ✅ Bueno: Usar envoltorios (wrappers) tipados
type BaseIzquierda struct {
    Base
}

type BaseDerecha struct {
    Base
}

type Estructura struct {
    izq BaseIzquierda
    der BaseDerecha
}
```

---

## 14.6 Embedding de Interfaces: Composición Avanzada

### 14.6.1 Interfaces Compuestas

Las interfaces se pueden componer embebiendo otras interfaces:

```go
type Lector interface {
    Leer(p []byte) (n int, err error)
}

type Escritor interface {
    Escribir(p []byte) (n int, err error)
}

// Composición de interfaces
type LectorEscritor interface {
    Lector
    Escritor
}

// Un tipo que implementa ambas interfaces satisface LectorEscritor automáticamente
type Archivo struct {}

func (a Archivo) Leer(p []byte) (int, error) {
    // implementación
    return len(p), nil
}

func (a Archivo) Escribir(p []byte) (int, error) {
    // implementación
    return len(p), nil
}

// Archivo automáticamente satisface LectorEscritor
var rw LectorEscritor = Archivo{}
```

### 14.6.2 La Interfaz io.ReadWriter

Un ejemplo real de composición de interfaces:

```go
// En package "io":
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}

// Cualquier tipo que implemente Read y Write satisface ReadWriter
// Sin necesidad de herencia de clases
```

### 14.6.3 Embedding de Interfaces para Extensibilidad

```go
// Interfaz base
type Servidor interface {
    Iniciar() error
    Detener() error
}

// Extensión: servidor con logging
type ServidorConLog interface {
    Servidor
    ObtenerLogs() []string
}

// Implementación
type MiServidor struct {
    logs []string
}

func (s *MiServidor) Iniciar() error {
    s.logs = append(s.logs, "Servidor iniciado")
    return nil
}

func (s *MiServidor) Detener() error {
    s.logs = append(s.logs, "Servidor detenido")
    return nil
}

func (s *MiServidor) ObtenerLogs() []string {
    return s.logs
}

// MiServidor satisface ServidorConLog
var srv ServidorConLog = &MiServidor{}
```

### 14.6.4 Usando Embedding de Interfaces en Implementaciones

```go
// Patrón: Wrapper que embebe una interfaz
type HandlerConMiddleware interface {
    http.Handler  // Embebe el Handler de http
}

// Implementación con composición
type LoggingHandler struct {
    handler http.Handler
}

func (h LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("Petición: %s %s", r.Method, r.URL)
    h.handler.ServeHTTP(w, r)
}

// Uso
mux := http.NewServeMux()
originalHandler := mux

loggingHandler := LoggingHandler{handler: originalHandler}
// loggingHandler satisface http.Handler implícitamente
http.ListenAndServe(":8080", loggingHandler)
```

---

## 14.7 Shadowing y Override: Sobrescribir Métodos Embebidos

### 14.7.1 Qué es el Shadowing

**Shadowing** es cuando defines un método en el struct que embebe, sombra o reemplaza el método del struct embebido:

```go
type Animal struct {
    nombre string
}

func (a *Animal) Sonido() string {
    return "Sonido genérico"
}

type Perro struct {
    Animal
}

// Sin sobrescritura
perro1 := &Perro{Animal: Animal{nombre: "Rex"}}
fmt.Println(perro1.Sonido())  // "Sonido genérico"

// Con sobrescritura (shadowing)
func (p *Perro) Sonido() string {
    return "Guau!"
}

perro2 := &Perro{Animal: Animal{nombre: "Max"}}
fmt.Println(perro2.Sonido())  // "Guau!"

// Pero aún puedes acceder al método original:
fmt.Println(perro2.Animal.Sonido())  // "Sonido genérico"
```

### 14.7.2 Shadowin de Campos

```go
type Persona struct {
    nombre string
}

type Empleado struct {
    Persona
    nombre string  // Sombra el nombre de Persona
}

emp := Empleado{
    Persona: Persona{nombre: "Juan"},
    nombre: "Juan García",
}

fmt.Println(emp.nombre)          // "Juan García" (sombra)
fmt.Println(emp.Persona.nombre)  // "Juan" (original)
```

### 14.7.3 Patrones de Override Prácticos

```go
type ServidorBase struct {
    puerto int
}

func (s *ServidorBase) Iniciar() error {
    fmt.Printf("Iniciando servidor en puerto %d\n", s.puerto)
    return nil
}

type ServidorConConfiguracion struct {
    ServidorBase
    timeout time.Duration
}

// Override: método mejorado
func (s *ServidorConConfiguracion) Iniciar() error {
    // Llamar al método del padre primero
    if err := s.ServidorBase.Iniciar(); err != nil {
        return err
    }
    // Luego lógica adicional
    fmt.Printf("Configurando timeout de %v\n", s.timeout)
    return nil
}

// Uso
servidor := &ServidorConConfiguracion{
    ServidorBase: ServidorBase{puerto: 8080},
    timeout:      5 * time.Second,
}

servidor.Iniciar()  // Ejecuta la versión sobrescrita
```

### 14.7.4 Cuándo Evitar el Shadowing

```go
// ❌ Malo: Shadowing confuso
type Procesador struct {
    datos []int
}

func (p *Procesador) Procesar() {
    fmt.Println("Procesando datos base")
}

type ProcesadorAvanzado struct {
    Procesador
    datos []int  // Sombra el campo de Procesador
}

func (p *ProcesadorAvanzado) Procesar() {
    // ¿Es p.datos del Procesador o ProcesadorAvanzado?
    // Confuso para quien lee el código
}

// ✅ Bueno: Nombres explícitos
type ProcesadorBase struct {
    datosBase []int
}

type ProcesadorMejorado struct {
    base      ProcesadorBase
    datosExtra []int
}

func (p *ProcesadorMejorado) Procesar() {
    // Claramente se usan datosBase y datosExtra
}
```

---

## 14.8 Embedding Anidado: Múltiples Niveles de Composición

### 14.8.1 Estructura de Embedding Anidado

```go
type Contacto struct {
    email string
    telefono string
}

type Dirección struct {
    calle  string
    ciudad string
}

type Persona struct {
    nombre string
    Contacto   // Nivel 1
    Dirección  // Nivel 1
}

type Empleado struct {
    Persona    // Nivel 2: embebe Persona que embebe Contacto y Dirección
    empleadoID string
}

emp := &Empleado{
    Persona: Persona{
        nombre: "Ana",
        Contacto: Contacto{
            email: "ana@example.com",
            telefono: "555-1234",
        },
        Dirección: Dirección{
            calle: "Calle Principal 123",
            ciudad: "Madrid",
        },
    },
    empleadoID: "EMP001",
}

// Acceso a campos anidados (todos promovidos):
fmt.Println(emp.nombre)      // "Ana"
fmt.Println(emp.email)       // "ana@example.com"
fmt.Println(emp.calle)       // "Calle Principal 123"
fmt.Println(emp.empleadoID)  // "EMP001"

// Acceso explícito también funciona:
fmt.Println(emp.Persona.Contacto.email)  // "ana@example.com"
```

### 14.8.2 Resolución de Nombres en Embedding Anidado

```go
type A struct {
    valor int
}

type B struct {
    A
    otro string
}

type C struct {
    B
    otro string  // Sombrea B.otro
}

c := C{
    B: B{
        A: A{valor: 42},
        otro: "B",
    },
    otro: "C",
}

fmt.Println(c.valor)       // 42 (promovido de A a través de B)
fmt.Println(c.otro)        // "C" (sombra el otro de B)
fmt.Println(c.B.otro)      // "B" (acceso explícito)
fmt.Println(c.B.A.valor)   // 42 (acceso completamente explícito)

// Diagrama de resolución:
//     A.valor
//      ↓
//      B (promueve valor de A, otro de B)
//      ↓
//      C (sombrea otro de B)
//      ↓
// c.valor → B.A.valor ✅
// c.otro → C.otro (no B.otro) ✅
```

### 14.8.3 Limitaciones y Cuándo es Apropiado

```go
// ⚠️ Embedding profundo es difícil de seguir
type Level1 struct {
    data1 int
}

type Level2 struct {
    Level1
    data2 int
}

type Level3 struct {
    Level2
    data3 int
}

type Level4 struct {
    Level3
    data4 int
}

l4 := Level4{
    Level3: Level3{
        Level2: Level2{
            Level1: Level1{data1: 1},
            data2: 2,
        },
        data3: 3,
    },
    data4: 4,
}

// ✅ Funciona, pero:
// - Difícil de entender la estructura
// - Difícil de mantener
// - Fácil de introducir errores

// ✅ Mejor: Usa composición explícita cuando es profundo
type Mejor struct {
    level1 Level1
    data2 int
    data3 int
    data4 int
}
```

### 14.8.4 Métodos en Embedding Anidado

```go
type Motor struct {
    potencia int
}

func (m *Motor) Arrancar() {
    fmt.Printf("Motor de %dCV arrancando\n", m.potencia)
}

type Transmisión struct {
    Motor    // Embebe Motor
    velocidades int
}

func (t *Transmisión) CambiarVelocidad(v int) {
    fmt.Printf("Cambiando a velocidad %d\n", v)
}

type Auto struct {
    Transmisión  // Embebe Transmisión que embebe Motor
    marca string
}

// Método de Motor se promueve a través de Transmisión a Auto
auto := &Auto{
    Transmisión: Transmisión{
        Motor: Motor{potencia: 150},
        velocidades: 6,
    },
    marca: "Toyota",
}

// Promoción de métodos a través de múltiples niveles
auto.Arrancar()         // Motor.Arrancar()
auto.CambiarVelocidad(3)  // Transmisión.CambiarVelocidad()
fmt.Println(auto.marca)   // "Toyota"
fmt.Println(auto.potencia)  // 150 (promovido)
fmt.Println(auto.velocidades) // 6 (promovido)
```

---

## 14.9 Embedding vs Delegación: Cuándo Usar Cada Uno

### 14.9.1 Comparación Embedding vs Delegación

```go
type Búfer struct {
    datos []byte
}

func (b *Búfer) Leer(n int) []byte {
    return b.datos[:n]
}

func (b *Búfer) Escribir(p []byte) {
    b.datos = append(b.datos, p...)
}

// OPCIÓN 1: Embedding (promoción automática)
type LoggerConEmbedding struct {
    Búfer        // Los métodos de Búfer se promueven
}

logger1 := &LoggerConEmbedding{
    Búfer: Búfer{datos: []byte("inicio")},
}
logger1.Leer(3)  // Promovido: llama a Búfer.Leer


// OPCIÓN 2: Delegación (explícito)
type LoggerConDelegación struct {
    búfer Búfer   // Campo nombrado
}

logger2 := &LoggerConDelegación{
    búfer: Búfer{datos: []byte("inicio")},
}

func (l *LoggerConDelegación) Leer(n int) []byte {
    return l.búfer.Leer(n)  // Delegación explícita
}

logger2.Leer(3)  // Llama a nuestra delegación que llama a búfer.Leer
```

### 14.9.2 Cuándo Usar Embedding

| Criterio | Usar Embedding |
|----------|---|
| **Quieres todos los métodos** | ✅ Promueve todos |
| **Relación es "ES UN"** | ✅ Conceptualmente correcto |
| **Pocos métodos** | ✅ Sin ruido |
| **Quieres composición clara** | ✅ Expresa composición |

```go
// ✅ Embedding es apropiado aquí
type Vehículo struct {
    marca string
    modelo string
}

func (v *Vehículo) Conducir() { /* ... */ }
func (v *Vehículo) Frenar() { /* ... */ }

type Auto struct {
    Vehículo          // ✅ Auto ES UN Vehículo
    numeroPuertas int
}

// Los métodos Conducir y Frenar se promueven naturalmente
```

### 14.9.3 Cuándo Usar Delegación

| Criterio | Usar Delegación |
|----------|---|
| **Solo algunos métodos** | ✅ Control fino |
| **Quieres agregar lógica** | ✅ Interceptar llamadas |
| **Necesitas decoración** | ✅ Logging, métricas, etc |
| **Evitar violación de Liskov** | ✅ No "fingir" ser el tipo |

```go
// ✅ Delegación es apropiada aquí
type BaseDatos struct {
    conexión *sql.DB
}

func (db *BaseDatos) Conectar() error { /* ... */ }
func (db *BaseDatos) Ejecutar(query string) error { /* ... */ }

type BaseDatosConCache struct {
    bd BaseDatos
    cache map[string]interface{}
}

// Solo delegamos algunos métodos, otros los decoramos
func (c *BaseDatosConCache) Ejecutar(query string) error {
    if resultado, ok := c.cache[query]; ok {
        fmt.Println("Del cache")
        return nil
    }
    return c.bd.Ejecutar(query)
}

// No promovemos Conectar; tenemos control explícito
```

### 14.9.4 El Patrón Adapter

Usar delegación para convertir una interfaz en otra:

```go
// Interfaz 1: Antigua
type LectorViejo interface {
    ObtenerBytes() []byte
}

// Interfaz 2: Nueva
type LectorNuevo interface {
    Leer(p []byte) (n int, err error)
}

type Adaptador struct {
    viejo LectorViejo
}

// Adaptador implementa LectorNuevo delegando a LectorViejo
func (a *Adaptador) Leer(p []byte) (int, error) {
    datos := a.viejo.ObtenerBytes()
    n := copy(p, datos)
    return n, nil
}

// Ahora puedes usar LectorViejo donde se espera LectorNuevo
var nuevo LectorNuevo = &Adaptador{viejo: algunLectorViejo}
```

---

## 14.10 Patrones Arquitectónicos: Casos de Uso Avanzados

### 14.10.1 Patrón Middleware

Embebiendo handlers para crear un chain:

```go
type Handler interface {
    Manejar(ctx context.Context, req interface{}) (interface{}, error)
}

// Middleware base
type LoggingMiddleware struct {
    siguiente Handler
}

func (m *LoggingMiddleware) Manejar(ctx context.Context, req interface{}) (interface{}, error) {
    fmt.Printf("Iniciando: %v\n", req)
    defer fmt.Println("Finalizando")
    return m.siguiente.Manejar(ctx, req)
}

type AutenticaciónMiddleware struct {
    siguiente Handler
}

func (a *AutenticaciónMiddleware) Manejar(ctx context.Context, req interface{}) (interface{}, error) {
    fmt.Println("Verificando autenticación...")
    if !estaAutenticado(ctx) {
        return nil, fmt.Errorf("no autenticado")
    }
    return a.siguiente.Manejar(ctx, req)
}

type HandlerReal struct{}

func (h *HandlerReal) Manejar(ctx context.Context, req interface{}) (interface{}, error) {
    fmt.Println("Manejando solicitud real")
    return "resultado", nil
}

// Composición del chain
manejador := &LoggingMiddleware{
    siguiente: &AutenticaciónMiddleware{
        siguiente: &HandlerReal{},
    },
}

// Llamada
manejador.Manejar(context.Background(), "datos")
// Salida:
// Iniciando: datos
// Verificando autenticación...
// Manejando solicitud real
// Finalizando
```

### 14.10.2 Patrón Decorator

Usar embedding para decorar comportamiento:

```go
// Interfaz base
type Componente interface {
    Dibujar()
    ObtenerCosto() float64
}

// Implementación concreta
type Ventana struct {
    ancho, alto int
}

func (v *Ventana) Dibujar() {
    fmt.Printf("Dibujando ventana %dx%d\n", v.ancho, v.alto)
}

func (v *Ventana) ObtenerCosto() float64 {
    return 100.0
}

// Decorador 1: Borde
type ConBorde struct {
    Componente  // Embebe el componente a decorar
    ancho int
}

func (b *ConBorde) Dibujar() {
    fmt.Println("Dibujando borde")
    b.Componente.Dibujar()  // Llama al componente original
    fmt.Println("Borde dibujado")
}

func (b *ConBorde) ObtenerCosto() float64 {
    return b.Componente.ObtenerCosto() + 50.0
}

// Decorador 2: Sombra
type ConSombra struct {
    Componente  // Embebe el componente actual (que puede estar decorado)
    blur int
}

func (s *ConSombra) Dibujar() {
    fmt.Println("Aplicando sombra...")
    s.Componente.Dibujar()
    fmt.Println("Sombra aplicada")
}

func (s *ConSombra) ObtenerCosto() float64 {
    return s.Componente.ObtenerCosto() + 20.0
}

// Uso: Composición de decoradores
ventana := &Ventana{ancho: 200, alto: 150}

conBorde := &ConBorde{Componente: ventana, ancho: 5}
conBordeSombra := &ConSombra{Componente: conBorde, blur: 3}

conBordeSombra.Dibujar()  // Múltiples capas de decoración
fmt.Printf("Costo total: $%.2f\n", conBordeSombra.ObtenerCosto())  // 100 + 50 + 20

// Salida:
// Aplicando sombra...
// Dibujando borde
// Dibujando ventana 200x150
// Borde dibujado
// Sombra aplicada
// Costo total: $170.00
```

### 14.10.3 Patrón Chain of Responsibility

```go
type Solicitud struct {
    monto  float64
    requiere string  // "aprobación" | "validación" | etc
}

type Manejador interface {
    Procesar(solicitud *Solicitud) error
    SetSiguiente(manejador Manejador)
}

type ProcesadorBase struct {
    siguiente Manejador
}

func (p *ProcesadorBase) SetSiguiente(m Manejador) {
    p.siguiente = m
}

type ValidadorMonto struct {
    ProcesadorBase
}

func (v *ValidadorMonto) Procesar(solicitud *Solicitud) error {
    fmt.Println("Validando monto...")
    if solicitud.monto <= 0 {
        return fmt.Errorf("monto inválido")
    }
    if v.siguiente != nil {
        return v.siguiente.Procesar(solicitud)
    }
    return nil
}

type AprobadorGerente struct {
    ProcesadorBase
}

func (a *AprobadorGerente) Procesar(solicitud *Solicitud) error {
    fmt.Println("Aprobando por gerente...")
    if solicitud.monto > 10000 {
        return fmt.Errorf("requiere aprobación de director")
    }
    if a.siguiente != nil {
        return a.siguiente.Procesar(solicitud)
    }
    fmt.Println("✓ Solicitud aprobada")
    return nil
}

// Uso
validador := &ValidadorMonto{}
gerente := &AprobadorGerente{}
validador.SetSiguiente(gerente)

err := validador.Procesar(&Solicitud{monto: 5000})
if err != nil {
    fmt.Println("Error:", err)
}
// Salida:
// Validando monto...
// Aprobando por gerente...
// ✓ Solicitud aprobada
```

### 14.10.4 Patrón Observer con Embedding

```go
type Suscriptor interface {
    Actualizar(evento string, datos interface{})
}

type Publicador struct {
    suscriptores []Suscriptor
}

func (p *Publicador) Suscribir(s Suscriptor) {
    p.suscriptores = append(p.suscriptores, s)
}

func (p *Publicador) Notificar(evento string, datos interface{}) {
    for _, s := range p.suscriptores {
        s.Actualizar(evento, datos)
    }
}

type LoggerObservador struct {
    Publicador  // Embebe Publicador para heredar Suscribir/Notificar
}

func (l *LoggerObservador) Actualizar(evento string, datos interface{}) {
    fmt.Printf("[LOG] Evento: %s, Datos: %v\n", evento, datos)
}

type Notificador struct {
    Publicador
}

func (n *Notificador) Actualizar(evento string, datos interface{}) {
    fmt.Printf("[NOTIFICACIÓN] Evento: %s enviado a admin\n", evento)
}

// Uso
publicador := &LoggerObservador{}
publicador.Suscribir(&LoggerObservador{})
publicador.Suscribir(&Notificador{})

publicador.Notificar("usuarioCreado", map[string]string{
    "id": "123",
    "nombre": "Juan",
})
```

---

## 14.11 Buenas Prácticas y Antipatrones

### 14.11.1 Buenas Prácticas

#### 1. Mantén Embeddings Planos

```go
// ✅ Bueno: Estructura plana
type Empleado struct {
    Persona
    Detalles
    Contacto
}

// ❌ Malo: Embeding profundo innecesario
type Empleado struct {
    Departamento  // que embebe Empresa que embebe...
}
```

#### 2. Nombre los Campos Ambiguos

```go
// ❌ Malo: "nombre" es ambiguo
type Empleado struct {
    Persona
    Proyecto
}

// ✅ Bueno: Nombres explícitos cuando hay conflicto
type Empleado struct {
    nombrePersona string
    nombreProyecto string
    Persona
    Proyecto
}

// O mejor aún:
type Empleado struct {
    persona Persona
    proyecto Proyecto
}
```

#### 3. Usa Embedding para Relaciones "ES UN"

```go
// ✅ Bueno: Auto ES UN Vehículo
type Auto struct {
    Vehículo
}

// ❌ Malo: Auto NO ES UNA Rueda
type Auto struct {
    Rueda  // Confuso: Auto no es realmente una Rueda
}

// Mejor:
type Auto struct {
    ruedas [4]Rueda
}
```

#### 4. Documenta la Intención

```go
// ✅ Bueno: Documentación clara
type ServidorHTTP struct {
    // Embebe ServidorBase para heredar funcionalidad de servidor básica
    // y métodos como Iniciar, Detener.
    ServidorBase
    
    // Configuración HTTP específica
    puerto    int
    handlers  map[string]http.Handler
}
```

#### 5. Prefiere Composición Sobre Embedding Profundo

```go
// ❌ Antipatrón: Cadena de embeddings
type A struct{ data int }
type B struct{ A }
type C struct{ B }
type D struct{ C }
type E struct{ D }

// ✅ Mejor: Composición explícita
type E struct {
    a A
    data int
}
```

### 14.11.2 Antipatrones

#### Antipatrón 1: Diamond Problem Simulado

```go
// ❌ Simulando problemas de herencia múltiple
type Base struct {
    id int
}

type Izquierda struct {
    Base
}

type Derecha struct {
    Base
}

type Ambos struct {
    Izquierda
    Derecha  // ¿Quién es la Base "real"?
}

a := Ambos{}
// a.id // ERROR: ambiguous

// ✅ Solución: Usar nombres explícitos
type Ambos struct {
    izq Izquierda
    der Derecha
}
```

#### Antipatrón 2: Sobre-Embeding

```go
// ❌ Embeding innecesario
type Aplicación struct {
    Configuración
    Logger
    BaseDatos
    Servidor
    Caché
    Autenticación
    Autorización
    // ... 20 embeddings más
}

// ✅ Mejor: Organizar en componentes
type Aplicación struct {
    config      Configuración
    servicios   Servicios
    middleware  Middleware
}

type Servicios struct {
    bd   BaseDatos
    auth Autenticación
    cache Caché
}
```

#### Antipatrón 3: Shadowing Confuso

```go
// ❌ Shadowing que oculta la intención
type Transporte struct {
    velocidad float64
}

type Auto struct {
    Transporte
    velocidad float64  // ¿Cuál estamos usando?
}

// ✅ Mejor: Nombres descriptivos
type Auto struct {
    velocidadActual float64
    velocidadMáxima float64
}
```

#### Antipatrón 4: Embeding de Punteros Inconsistente

```go
// ❌ Inconsistente: algunos punteros, otros valores
type Contenedor struct {
    *A      // Puntero
    B       // Valor
    *C      // Puntero
}

// ✅ Consistente: todos valores o todos punteros
type Contenedor struct {
    a A
    b B
    c C
}

// O si necesitas punteros:
type Contenedor struct {
    a *A
    b *B
    c *C
}
```

#### Antipatrón 5: Interfaz Gorda Embebida

```go
// ❌ Embebiendo una interfaz gorda y no usándola
type ConexiónGrande interface {
    Conectar()
    Desconectar()
    Leer()
    Escribir()
    // ... 20 métodos más
}

type MiServicio struct {
    ConexiónGrande  // Pero solo necesito Leer y Escribir
}

// ✅ Mejor: Embebes solo lo que necesitas
type ConexiónSimple interface {
    Leer()
    Escribir()
}

type MiServicio struct {
    ConexiónSimple
}
```

### 14.11.3 Checklist de Diseño

Antes de embeber, pregúntate:

1. **¿Es una relación "ES UN"?** → Sí: embebe; No: compón explícitamente
2. **¿Voy a usar TODOS los métodos?** → Sí: embebe; No: compón
3. **¿Hay conflictos de nombres?** → Sí: reconsidera; No: procede
4. **¿El embedding es más claro que la delegación?** → Sí: embebe; No: delega
5. **¿Puedo explicar esta estructura a alguien nuevo?** → Sí: OK; No: simplifica

---

## Ejercicios Progresivos

### Ejercicio 1: Jerarquía de Vehículos

**Objetivo:** Crear una jerarquía de vehículos usando embedding sin herencia.

**Requisitos:**
- Struct `Vehículo` base con campos: marca, modelo, velocidadMáxima
- Métodos: Describir(), ObtenerVelocidadMáxima()
- Struct `Auto` embebiendo `Vehículo`, con campo adicional: numeroPuertas
- Struct `Moto` embebiendo `Vehículo`, con campo: tipo (Deportiva/Crucero)
- Struct `Camión` embebiendo `Vehículo`, con campo: capacidadCarga
- Cada tipo sobrescribe Describir() con descripción específica
- Función que recibe Vehículo por interfaz

**Solución:**

```go
package main

import (
    "fmt"
)

type Describible interface {
    Describir()
    ObtenerVelocidadMáxima() int
}

type Vehículo struct {
    marca           string
    modelo          string
    velocidadMáxima int
}

func (v *Vehículo) ObtenerVelocidadMáxima() int {
    return v.velocidadMáxima
}

func (v *Vehículo) Describir() {
    fmt.Printf("Vehículo: %s %s (Vel. Máx: %d km/h)\n",
        v.marca, v.modelo, v.velocidadMáxima)
}

type Auto struct {
    Vehículo
    numeroPuertas int
}

func (a *Auto) Describir() {
    fmt.Printf("Auto: %s %s (%d puertas, Vel. Máx: %d km/h)\n",
        a.marca, a.modelo, a.numeroPuertas, a.velocidadMáxima)
}

type Moto struct {
    Vehículo
    tipo string
}

func (m *Moto) Describir() {
    fmt.Printf("Moto %s: %s %s (Vel. Máx: %d km/h)\n",
        m.tipo, m.marca, m.modelo, m.velocidadMáxima)
}

type Camión struct {
    Vehículo
    capacidadCarga int
}

func (c *Camión) Describir() {
    fmt.Printf("Camión: %s %s (Carga: %d kg, Vel. Máx: %d km/h)\n",
        c.marca, c.modelo, c.capacidadCarga, c.velocidadMáxima)
}

func MostrarVehículos(vehiculos ...Describible) {
    for _, v := range vehiculos {
        v.Describir()
        fmt.Printf("  Velocidad máxima: %d km/h\n", v.ObtenerVelocidadMáxima())
    }
}

func main() {
    auto := &Auto{
        Vehículo: Vehículo{
            marca:           "Toyota",
            modelo:          "Camry",
            velocidadMáxima: 180,
        },
        numeroPuertas: 4,
    }

    moto := &Moto{
        Vehículo: Vehículo{
            marca:           "Harley-Davidson",
            modelo:          "Sportster",
            velocidadMáxima: 200,
        },
        tipo: "Crucero",
    }

    camión := &Camión{
        Vehículo: Vehículo{
            marca:           "Volvo",
            modelo:          "FH16",
            velocidadMáxima: 120,
        },
        capacidadCarga: 25000,
    }

    MostrarVehículos(auto, moto, camión)
}
```

### Ejercicio 2: Sistema de Eventos con Múltiples Tipos

**Objetivo:** Crear un sistema donde los eventos pueden ser de múltiples tipos simultáneamente.

**Requisitos:**
- Struct `EventoBase` con: ID, timestamp, descripción
- Struct `EventoConPrioridad` embebiendo `EventoBase`, con nivel
- Struct `EventoConCategía` embebiendo `EventoBase`, con categoría
- Struct `EventoCompleto` embebiendo tanto Prioridad como Categoría, resolver ambigüedad
- Métodos: ObtenerResumen(), EsUrgente(), PerteneceA()

**Solución:**

```go
package main

import (
    "fmt"
    "time"
)

type EventoBase struct {
    id          string
    timestamp   time.Time
    descripción string
}

func (e *EventoBase) ObtenerID() string {
    return e.id
}

func (e *EventoBase) ObtenerTimestamp() time.Time {
    return e.timestamp
}

type Prioridad struct {
    EventoBase
    nivel int  // 1 (baja) a 5 (crítica)
}

func (p *Prioridad) ObtenerResumen() string {
    return fmt.Sprintf("[Evento %s] %s [Prioridad %d]",
        p.id, p.descripción, p.nivel)
}

func (p *Prioridad) EsUrgente() bool {
    return p.nivel >= 4
}

type Categoría struct {
    EventoBase
    categoría string
}

func (c *Categoría) ObtenerResumen() string {
    return fmt.Sprintf("[Evento %s] %s [Categoría: %s]",
        c.id, c.descripción, c.categoría)
}

func (c *Categoría) PerteneceA(cat string) bool {
    return c.categoría == cat
}

type EventoCompleto struct {
    prioridad  Prioridad
    categoría  Categoría
    anotaciones string
}

func (ec *EventoCompleto) ObtenerID() string {
    return ec.prioridad.id
}

func (ec *EventoCompleto) ObtenerTimestamp() time.Time {
    return ec.prioridad.timestamp
}

func (ec *EventoCompleto) ObtenerResumen() string {
    return fmt.Sprintf("[ID: %s] %s | Prioridad: %d | Categoría: %s | Anotaciones: %s",
        ec.prioridad.id,
        ec.prioridad.descripción,
        ec.prioridad.nivel,
        ec.categoría.categoría,
        ec.anotaciones)
}

func (ec *EventoCompleto) EsUrgente() bool {
    return ec.prioridad.EsUrgente()
}

func (ec *EventoCompleto) PerteneceA(cat string) bool {
    return ec.categoría.PerteneceA(cat)
}

func main() {
    evento := EventoCompleto{
        prioridad: Prioridad{
            EventoBase: EventoBase{
                id:          "EVT001",
                timestamp:   time.Now(),
                descripción: "Error en producción",
            },
            nivel: 5,
        },
        categoría: Categoría{
            EventoBase: EventoBase{
                id:          "EVT001",
                timestamp:   time.Now(),
                descripción: "Error en producción",
            },
            categoría: "Sistema",
        },
        anotaciones: "Base de datos no responde",
    }

    fmt.Println(evento.ObtenerResumen())
    fmt.Printf("¿Es urgente? %v\n", evento.EsUrgente())
    fmt.Printf("¿Pertenece a 'Sistema'? %v\n", evento.PerteneceA("Sistema"))
}
```

### Ejercicio 3: Middleware HTTP con Composición

**Objetivo:** Crear un sistema de middleware HTTP usando embedding y composición de decoradores.

**Requisitos:**
- Interfaz `Handler` con método ServeHTTP
- Middleware de Logging: registra cada solicitud
- Middleware de Autenticación: verifica token
- Middleware de Compresión: comprime respuesta
- Middleware de Caché: cachea respuestas GET
- Composición múltiple de middleware

**Solución:**

```go
package main

import (
    "bytes"
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h(w, r)
}

// Middleware 1: Logging
type LoggingMiddleware struct {
    next http.Handler
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("[LOG] %s %s\n", r.Method, r.URL.Path)
    m.next.ServeHTTP(w, r)
    log.Printf("[LOG] Respuesta enviada\n")
}

// Middleware 2: Autenticación
type AutenticaciónMiddleware struct {
    next http.Handler
}

func (m *AutenticaciónMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("Authorization")
    if token == "" {
        http.Error(w, "No autorizado", http.StatusUnauthorized)
        return
    }
    log.Printf("[AUTH] Token válido: %s\n", token)
    m.next.ServeHTTP(w, r)
}

// Middleware 3: Caché simple
type Caché struct {
    cache map[string]string
    next  http.Handler
}

func NewCaché(next http.Handler) *Caché {
    return &Caché{
        cache: make(map[string]string),
        next:  next,
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
    body       bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    return rw.body.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
    rw.statusCode = statusCode
    rw.ResponseWriter.WriteHeader(statusCode)
}

func (m *Caché) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        if resp, ok := m.cache[r.URL.Path]; ok {
            log.Printf("[CACHE] HIT: %s\n", r.URL.Path)
            w.Header().Set("X-Cache", "HIT")
            fmt.Fprint(w, resp)
            return
        }
    }

    rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
    m.next.ServeHTTP(rw, r)

    if rw.statusCode == http.StatusOK && r.Method == "GET" {
        m.cache[r.URL.Path] = rw.body.String()
        log.Printf("[CACHE] MISS: cachada %s\n", r.URL.Path)
    }

    rw.ResponseWriter.Write(rw.body.Bytes())
}

// Middleware 4: Compresión
type CompresiónMiddleware struct {
    next http.Handler
}

func (m *CompresiónMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Accept-Encoding") != "gzip" {
        m.next.ServeHTTP(w, r)
        return
    }

    w.Header().Set("Content-Encoding", "gzip")
    gz := gzip.NewWriter(w)
    defer gz.Close()

    gzw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}
    m.next.ServeHTTP(gzw, r)
}

type gzipResponseWriter struct {
    http.ResponseWriter
    Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

// Handler real
func manejadorAPI(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Respuesta de API: %s", r.URL.Path)
}

func main() {
    handler := HandlerFunc(manejadorAPI)

    // Composición de middleware
    conCompresiónHandler := &CompresiónMiddleware{next: handler}
    conCachéHandler := NewCaché(conCompresiónHandler)
    conAuthHandler := &AutenticaciónMiddleware{next: conCachéHandler}
    conLogHandler := &LoggingMiddleware{next: conAuthHandler}

    http.Handle("/api", conLogHandler)
    log.Println("Servidor en :8080")
    // http.ListenAndServe(":8080", nil)
}
```

### Ejercicio 4: Árbol de Componentes UI

**Objetivo:** Crear un sistema de componentes UI donde cada componente puede embeber otros.

**Requisitos:**
- `Componente` base: renderText(), obtenerAltura(), obtenerAncho()
- Tipos: `Texto`, `Botón`, `Panel`, `Ventana`
- Algunos embedben otros (Ventana embebe Panel, Panel embebe múltiples componentes)
- Métodos especializados para cada tipo
- Renderizado recursivo

**Solución:**

```go
package main

import (
    "fmt"
    "strings"
)

type Componente interface {
    Renderizar(indent int)
    ObtenerAltura() int
    ObtenerAncho() int
}

type Texto struct {
    contenido string
}

func (t *Texto) Renderizar(indent int) {
    fmt.Printf("%s[Texto] %s\n", strings.Repeat("  ", indent), t.contenido)
}

func (t *Texto) ObtenerAltura() int {
    return 1
}

func (t *Texto) ObtenerAncho() int {
    return len(t.contenido)
}

type Botón struct {
    Texto        // Embebe Texto para heredar comportamiento
    acción string
}

func (b *Botón) Renderizar(indent int) {
    fmt.Printf("%s[Botón] %s (onClick: %s)\n",
        strings.Repeat("  ", indent), b.contenido, b.acción)
}

func (b *Botón) Click() {
    fmt.Printf("Ejecutando: %s\n", b.acción)
}

type Panel struct {
    titulo      string
    componentes []Componente
    ancho, alto int
}

func (p *Panel) Renderizar(indent int) {
    fmt.Printf("%s┌─ Panel: %s (%dx%d)\n",
        strings.Repeat("  ", indent), p.titulo, p.ancho, p.alto)
    for _, c := range p.componentes {
        c.Renderizar(indent + 1)
    }
    fmt.Printf("%s└─ Fin Panel\n", strings.Repeat("  ", indent))
}

func (p *Panel) AgregarComponente(c Componente) {
    p.componentes = append(p.componentes, c)
}

func (p *Panel) ObtenerAltura() int {
    return p.alto
}

func (p *Panel) ObtenerAncho() int {
    return p.ancho
}

type Ventana struct {
    Panel           // Embebe Panel
    título string
    minimizable bool
}

func (v *Ventana) Renderizar(indent int) {
    barra := "═"
    if v.minimizable {
        barra += " [_] [-] [X]"
    }
    fmt.Printf("%s╔%s╗ %s\n",
        strings.Repeat("  ", indent), strings.Repeat(barra, 20), v.título)

    // Renderizar componentes del panel
    for _, c := range v.componentes {
        c.Renderizar(indent + 1)
    }

    fmt.Printf("%s╚%s╝\n",
        strings.Repeat("  ", indent), strings.Repeat("═", 20))
}

func (v *Ventana) Minimizar() {
    fmt.Printf("Minimizando ventana: %s\n", v.título)
}

func (v *Ventana) Maximizar() {
    fmt.Printf("Maximizando ventana: %s\n", v.título)
}

func main() {
    // Crear componentes
    texto1 := &Texto{contenido: "Bienvenido"}
    botón1 := &Botón{Texto: Texto{contenido: "Aceptar"}, acción: "guardar()"}
    botón2 := &Botón{Texto: Texto{contenido: "Cancelar"}, acción: "cerrar()"}

    // Crear panel
    panelBotones := &Panel{
        titulo: "Controles",
        ancho:  40,
        alto:   10,
    }
    panelBotones.AgregarComponente(botón1)
    panelBotones.AgregarComponente(botón2)

    // Crear ventana que contiene el panel
    ventana := &Ventana{
        Panel: Panel{
            titulo:      "Formulario",
            ancho:       50,
            alto:        20,
            componentes: []Componente{texto1, panelBotones},
        },
        título:       "Aplicación Principal",
        minimizable:  true,
    }

    ventana.Renderizar(0)
    fmt.Println()
    ventana.Minimizar()
}
```

### Ejercicio 5: Sistema de Roles con Ambigüedad Resuelta

**Objetivo:** Crear un sistema donde un usuario puede tener múltiples roles, resolviendo conflictos de permisos.

**Requisitos:**
- `Usuario` base: id, nombre
- Interfaces: `Permisos`, `Responsabilidades`
- Tipos de rol: `Administrador`, `Empleado`, `Cliente`
- `UsuarioMultirrol` que es múltiples roles a la vez
- Resolver ambigüedad de permisos (qué rol prevalece)
- Métodos: ¿TienePèrmiso()?, ListarResponsabilidades()

**Solución:**

```go
package main

import (
    "fmt"
)

type Usuario struct {
    id    string
    nombre string
}

type Permisos interface {
    TienePèrmiso(acción string) bool
    ObtenerPermisosLista() []string
}

type Responsabilidades interface {
    ListarResponsabilidades() []string
    PuedeRevisar(item string) bool
}

type Administrador struct {
    Usuario
}

func (a *Administrador) TienePèrmiso(acción string) bool {
    permisos := map[string]bool{
        "leer":        true,
        "escribir":    true,
        "eliminar":    true,
        "administrar": true,
    }
    return permisos[acción]
}

func (a *Administrador) ObtenerPermisosLista() []string {
    return []string{"leer", "escribir", "eliminar", "administrar"}
}

func (a *Administrador) ListarResponsabilidades() []string {
    return []string{
        "Gestionar sistema",
        "Aprobar cambios críticos",
        "Gestionar usuarios",
    }
}

func (a *Administrador) PuedeRevisar(item string) bool {
    return true // Admin puede revisar todo
}

type Empleado struct {
    Usuario
    departamento string
}

func (e *Empleado) TienePèrmiso(acción string) bool {
    permisos := map[string]bool{
        "leer":     true,
        "escribir": true,
        "eliminar": false,
    }
    return permisos[acción]
}

func (e *Empleado) ObtenerPermisosLista() []string {
    return []string{"leer", "escribir"}
}

func (e *Empleado) ListarResponsabilidades() []string {
    return []string{
        fmt.Sprintf("Trabajar en %s", e.departamento),
        "Completar tareas asignadas",
    }
}

func (e *Empleado) PuedeRevisar(item string) bool {
    return item == "propios"
}

type Cliente struct {
    Usuario
}

func (c *Cliente) TienePèrmiso(acción string) bool {
    permisos := map[string]bool{
        "leer": true,
    }
    return permisos[acción]
}

func (c *Cliente) ObtenerPermisosLista() []string {
    return []string{"leer"}
}

func (c *Cliente) ListarResponsabilidades() []string {
    return []string{
        "Usar el servicio",
        "Pagar facturas",
    }
}

func (c *Cliente) PuedeRevisar(item string) bool {
    return false
}

// UsuarioMultirrol resuelve ambigüedad por prioridad
type UsuarioMultirrol struct {
    id    string
    nombre string
    admin   *Administrador
    empleado *Empleado
    cliente *Cliente
}

func (u *UsuarioMultirrol) TienePèrmiso(acción string) bool {
    // Prioridad: Admin > Empleado > Cliente
    if u.admin != nil && u.admin.TienePèrmiso(acción) {
        return true
    }
    if u.empleado != nil && u.empleado.TienePèrmiso(acción) {
        return true
    }
    if u.cliente != nil && u.cliente.TienePèrmiso(acción) {
        return true
    }
    return false
}

func (u *UsuarioMultirrol) ListarResponsabilidades() []string {
    var resp []string
    if u.admin != nil {
        resp = append(resp, u.admin.ListarResponsabilidades()...)
    }
    if u.empleado != nil {
        resp = append(resp, u.empleado.ListarResponsabilidades()...)
    }
    if u.cliente != nil {
        resp = append(resp, u.cliente.ListarResponsabilidades()...)
    }
    return resp
}

func (u *UsuarioMultirrol) ObtenerRoles() []string {
    var roles []string
    if u.admin != nil {
        roles = append(roles, "Administrador")
    }
    if u.empleado != nil {
        roles = append(roles, "Empleado")
    }
    if u.cliente != nil {
        roles = append(roles, "Cliente")
    }
    return roles
}

func VerificarAcceso(u interface{}, acción string) {
    switch usuario := u.(type) {
    case *Administrador:
        fmt.Printf("[%s] %s puede '%s': %v\n",
            usuario.nombre, "Administrador", acción, usuario.TienePèrmiso(acción))
    case *Empleado:
        fmt.Printf("[%s] %s puede '%s': %v\n",
            usuario.nombre, "Empleado", acción, usuario.TienePèrmiso(acción))
    case *Cliente:
        fmt.Printf("[%s] %s puede '%s': %v\n",
            usuario.nombre, "Cliente", acción, usuario.TienePèrmiso(acción))
    case *UsuarioMultirrol:
        fmt.Printf("[%s] Roles: %v puede '%s': %v\n",
            usuario.nombre, usuario.ObtenerRoles(), acción, usuario.TienePèrmiso(acción))
    }
}

func main() {
    // Usuarios individuales
    admin := &Administrador{Usuario: Usuario{id: "1", nombre: "Juan Admin"}}
    empleado := &Empleado{
        Usuario: Usuario{id: "2", nombre: "María Empleado"},
        departamento: "TI",
    }
    cliente := &Cliente{Usuario: Usuario{id: "3", nombre: "Pedro Cliente"}}

    // Verificar acceso individual
    fmt.Println("=== Acceso Individual ===")
    VerificarAcceso(admin, "administrar")
    VerificarAcceso(empleado, "escribir")
    VerificarAcceso(cliente, "eliminar")

    // Usuario multirrol
    fmt.Println("\n=== Usuario Multirrol ===")
    multirrol := &UsuarioMultirrol{
        id: "4",
        nombre: "Ana Multirrol",
        empleado: &Empleado{
            Usuario: Usuario{id: "4", nombre: "Ana"},
            departamento: "Ventas",
        },
        cliente: &Cliente{Usuario: Usuario{id: "4", nombre: "Ana"}},
    }

    VerificarAcceso(multirrol, "escribir")
    VerificarAcceso(multirrol, "eliminar")
    VerificarAcceso(multirrol, "administrar")

    fmt.Println("\nResponsabilidades:")
    for _, r := range multirrol.ListarResponsabilidades() {
        fmt.Printf("  • %s\n", r)
    }
}
```

---

## Resumen y Conclusiones

### Puntos Clave

1. **Go rechaza la herencia** para evitar complejidad, acoplamiento fuerte y el problema del diamante.

2. **Embedding es la composición de Go:**
   - Campos anónimos promocionan campos y métodos automáticamente
   - Proporciona un mecanismo ligero para compartir código
   - Requiere explicitación cuando hay conflictos

3. **Resolución de ambigüedad:**
   - Conflictos requieren acceso explícito
   - Shadowing permite override controlado
   - Múltiples embeddings de la misma base causa conflictos (similar al diamond problem)

4. **Interfase embedding para composición:**
   - Las interfaces pueden embeberse para crear interfaces compuestas
   - Permite polimorfismo sin jerarquía
   - Ejemplo: `io.ReadWriter = Reader + Writer`

5. **Patrones arquitectónicos con embedding:**
   - Middleware: composición de handlers
   - Decorators: embeber interfaz para decorar comportamiento
   - Chain of Responsibility: enlaces de procesadores

6. **Buenas prácticas:**
   - Mantén embeddings planos y simples
   - Usa embedding para relaciones "ES UN"
   - Prefiere composición explícita sobre embedding profundo
   - Documenta la intención del embedding
   - Resuelve conflictos de nombres proactivamente

7. **Cuándo embeber vs delegar:**
   - Embebe: relación "ES UN", quieres todos los métodos
   - Delega: necesitas control fino, lógica adicional, evitar ambigüedad

El embedding de Go es un diseño elegante que evita los problemas de herencia clásica mientras proporciona reutilización de código pragmática.


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/14-embedding-y-composicion/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/14-embedding-y-composicion):

```bash
cd examples/14-embedding-y-composicion
go run .
```
