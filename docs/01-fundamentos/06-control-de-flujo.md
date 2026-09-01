# Capítulo 6: Control de flujo

## Índice del Capítulo 6

1. [6.1 ¿Qué es Control de Flujo?](#61-qué-es-control-de-flujo)
2. [6.2 Sentencias if/else](#62-sentencias-ifelse)
3. [6.3 Sentencia switch](#63-sentencia-switch)
4. [6.4 Bucle for - Forma Tradicional](#64-bucle-for---forma-tradicional)
5. [6.5 Bucle for - Forma While](#65-bucle-for---forma-while)
6. [6.6 Bucle for - Range](#66-bucle-for---range)
7. [6.7 break y continue](#67-break-y-continue)
8. [6.8 Bucles Anidados y Labels](#68-bucles-anidados-y-labels)
9. [6.9 defer - Ejecución Diferida](#69-defer---ejecución-diferida)
10. [6.10 Patrones y Buenas Prácticas](#610-patrones-y-buenas-prácticas)

---

## 6.1 ¿Qué es Control de Flujo?

### Definición

Control de flujo es el **orden en que se ejecutan las sentencias** de tu programa:

```
Sin control de flujo:
    línea 1
    línea 2
    línea 3
    línea 4
    línea 5

Con control de flujo:
    línea 1
    SI condición: línea 2
    SI no: línea 3
    REPETIR línea 4 (mientras X)
    línea 5
```

### Diagrama General

```

 CONTROL DE FLUJO EN GO                       │

                                              │
 1. EJECUCIÓN SECUENCIAL (default)           │
    └─ línea por línea, en orden              │
                                              │
 2. BIFURCACIÓN (Decisiones)                 │
    ├─ if / else / else if                    │
 switch / case                          │    ├
    └─ Resultado: ejecuta UNA rama            │
                                              │
 3. REPETICIÓN (Bucles)                      │
    ├─ for (forma tradicional)                │
    ├─ for (mientras - while-like)            │
    ├─ for range (sobre collections)          │
    └─ for infinito                           │
                                              │
 4. SALTO (Control fino)                     │
    ├─ break (salir de bucle)                 │
    ├─ continue (siguiente iteración)         │
    ├─ goto (evitar en código real)           │
    └─ return (salir de función)              │
                                              │

```

---

## 6.2 Sentencias if/else

### Forma Básica

```go
if condición {
    // Se ejecuta si condición es true
}
```

**Ejemplo:**

```go
edad := 18

if edad >= 18 {
    fmt.Println("Eres adulto")
}
```

### if/else

```go
if condición {
    // Se ejecuta si es true
} else {
    // Se ejecuta si es false
}
```

**Ejemplo:**

```go
edad := 15

if edad >= 18 {
    fmt.Println("Eres adulto")
} else {
    fmt.Println("Eres menor")
}

// Output: Eres menor
```

### if/else if/else

```go
if condición1 {
    // Si condición1
} else if condición2 {
    // Si condición1 es false Y condición2 es true
} else if condición3 {
    // Si condición1 y 2 son false, Y condición3 es true
} else {
    // Si TODAS las anteriores son false
}
```

**Ejemplo:**

```go
calificacion := 75

if calificacion >= 90 {
    fmt.Println("A (Excelente)")
} else if calificacion >= 80 {
    fmt.Println("B (Bueno)")
} else if calificacion >= 70 {
    fmt.Println("C (Satisfactorio)")
} else if calificacion >= 60 {
    fmt.Println("D (Pasable)")
} else {
    fmt.Println("F (Suspendido)")
}

// Output: C (Satisfactorio)
```

### Alcance de Variables en if

```go
if x := 5; x > 3 {
    fmt.Println(x)      // ✅ OK: x existe en este scope
} else {
    fmt.Println(x)      // ✅ OK: x existe en else también
}

fmt.Println(x)          // ❌ ERROR: x no existe fuera del if
```

**Sintaxis:**

```go
if inicializacion; condición {
    // ...
}
```

**Ventaja: Variable local al if**

```go
// Malo: y existe fuera del if
y := getValue()
if y > 10 {
    fmt.Println("Sí")
}
// y sigue siendo accesible aquí

// Bueno: x solo existe en el if
if x := getValue(); x > 10 {
    fmt.Println("Sí")
}
// x no existe aquí
```

### Errores Comunes

**Paréntesis NO requeridos (diferencia con C/Java):**

```go
if edad > 18 {       // ✅ Correcto
    fmt.Println("OK")
}

if (edad > 18) {     // ✅ También OK, pero raro
    fmt.Println("OK")
}
```

**Else DEBE estar en la misma línea:**

```go
if edad > 18 {
    fmt.Println("Adulto")
}
else {              // ❌ ERROR: else en línea nueva
    fmt.Println("Menor")
}

// Correcto:
if edad > 18 {
    fmt.Println("Adulto")
} else {            // ✅ OK: } else {
    fmt.Println("Menor")
}
```

---

## 6.3 Sentencia switch

### Forma Básica

```go
switch valor {
case opcion1:
    // Se ejecuta si valor == opcion1
case opcion2:
    // Se ejecuta si valor == opcion2
default:
    // Se ejecuta si ningún case coincide
}
```

**Ejemplo:**

```go
dia := 3

switch dia {
case 1:
    fmt.Println("Lunes")
case 2:
    fmt.Println("Martes")
case 3:
    fmt.Println("Miércoles")
case 4:
    fmt.Println("Jueves")
case 5:
    fmt.Println("Viernes")
default:
    fmt.Println("Fin de semana")
}

// Output: Miércoles
```

### Fall-Through (break implícito)

Go tiene **break implícito** en cada case:

```go
// Comparación con C/C++ (sin break, fall-through)

C/C++:
    switch (x) {
        case 1:
            printf("Uno");
        case 2:
            printf("Dos");     // ← Se ejecuta aunque sea case 2
            // Falta break, cae al siguiente
    }

Go (DIFERENTE):
    switch x {
    case 1:
        fmt.Println("Uno")
    case 2:
        fmt.Println("Dos")     // ← Solo se ejecuta si x == 2
    }                          // ← Break implícito
```

**Si necesitas fall-through, usa `fallthrough`:**

```go
x := 2

switch x {
case 1:
    fmt.Println("Uno")
    fallthrough          // Continúa al siguiente case
case 2:
    fmt.Println("Dos")
    fallthrough          // Continúa al siguiente
case 3:
    fmt.Println("Tres")
}

// Output:
// Dos
// Tres
```

### Multiple Cases

```go
dia := "Sábado"

switch dia {
case "Lunes", "Martes", "Miércoles", "Jueves", "Viernes":
    fmt.Println("Es día laboral")
case "Sábado", "Domingo":
    fmt.Println("Es fin de semana")
}

// Output: Es fin de semana
```

### Switch sin Expresión (switch true)

```go
edad := 25
ingresos := 50000

switch {
case edad < 18:
    fmt.Println("Menor de edad")
case ingresos < 30000:
    fmt.Println("Ingresos bajos")
case edad >= 18 && ingresos >= 50000:
    fmt.Println("Adulto con buenos ingresos")
default:
    fmt.Println("Otra categoría")
}

// Output: Adulto con buenos ingresos
```

**Equivalente a if/else if, pero más legible:**

```go
// Versión if/else if (más verbose)
if edad < 18 {
    fmt.Println("Menor de edad")
} else if ingresos < 30000 {
    fmt.Println("Ingresos bajos")
} else if edad >= 18 && ingresos >= 50000 {
    fmt.Println("Adulto con buenos ingresos")
} else {
    fmt.Println("Otra categoría")
}

// Versión switch true (más concisa)
switch {
case edad < 18:
    fmt.Println("Menor de edad")
case ingresos < 30000:
    fmt.Println("Ingresos bajos")
case edad >= 18 && ingresos >= 50000:
    fmt.Println("Adulto con buenos ingresos")
default:
    fmt.Println("Otra categoría")
}
```

### Switch con Inicialización

```go
switch x := getValue(); x {
case 1:
    fmt.Println("Uno")
case 2:
    fmt.Println("Dos")
}

// x solo existe dentro del switch
```

---

## 6.4 Bucle for - Forma Tradicional

### Forma Básica

```go
for inicializacion; condicion; actualizacion {
    // Se ejecuta mientras condición sea true
}
```

**Ejemplo: Contar del 0 al 9**

```go
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// Output:
// 0
// 1
// ...
// 9
```

**Desglose:**

```
for i := 0;    i < 10;    i++ {
    │           │          │
    │           │          Actualización (después de cada iteración)
    │           Condición (se evalúa antes de cada iteración)
    Inicialización (se ejecuta una sola vez)
```

### Formas Alternativas

**Omitir inicialización:**

```go
i := 0
for ; i < 10; i++ {
    fmt.Println(i)
}
```

**Omitir actualización:**

```go
for i := 0; i < 10; {
    fmt.Println(i)
    i++              // Manual
}
```

**Omitir inicialización y actualización:**

```go
i := 0
for ; i < 10; {
    fmt.Println(i)
    i++
}

// Equivalente a while
for i < 10 {
    fmt.Println(i)
    i++
}
```

---

## 6.5 Bucle for - Forma While

Go NO tiene palabra clave `while`, pero for puede actuar como while:

```go
for condicion {
    // Se ejecuta mientras condición sea true
}
```

**Ejemplo: Contar mientras nmero < 5**

```go
i := 0

for i < 5 {
    fmt.Println(i)
    i++
}

// Output: 0, 1, 2, 3, 4
```

### Bucle Infinito

```go
for {
    fmt.Println("Infinito")
    // Salir con break
}
```

**Caso real: Menú interactivo**

```go
for {
    fmt.Println("1. Opción A")
    fmt.Println("2. Opción B")
    fmt.Println("3. Salir")
    fmt.Print("Elige: ")
    
    var opcion int
    fmt.Scanln(&opcion)
    
    switch opcion {
    case 1:
        fmt.Println("Elegiste A")
    case 2:
        fmt.Println("Elegiste B")
    case 3:
        fmt.Println("Saliendo...")
        break  // Sale del for
    default:
        fmt.Println("Opción inválida")
    }
}
```

---

## 6.6 Bucle for - Range

### Range sobre Slices

```go
numeros := []int{10, 20, 30, 40}

for i, valor := range numeros {
    fmt.Printf("Índice: %d, Valor: %d\n", i, valor)
}

// Output:
// Índice: 0, Valor: 10
// Índice: 1, Valor: 20
// Índice: 2, Valor: 30
// Índice: 3, Valor: 40
```

**Omitir índice:**

```go
for _, valor := range numeros {
    fmt.Println(valor)      // Solo valores
}
```

**Omitir valor:**

```go
for i := range numeros {
    fmt.Printf("Índice: %d\n", i)   // Solo índices
}
```

### Range sobre Strings

```go
s := "Hola"

// Por defecto: retorna rune e índice
for i, r := range s {
    fmt.Printf("Índice: %d, Rune: %c\n", i, r)
}

// Output:
// Índice: 0, Rune: H
// Índice: 1, Rune: o
// Índice: 2, Rune: l
// Índice: 3, Rune: a
```

### Range sobre Maps

```go
edades := map[string]int{
    "Juan": 25,
    "María": 30,
}

for nombre, edad := range edades {
    fmt.Printf("%s tiene %d años\n", nombre, edad)
}

// Output (orden aleatorio):
// Juan tiene 25 años
// María tiene 30 años
```

---

## 6.7 break y continue

### break - Salir del Bucle

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break       // Sale del for
    }
    fmt.Println(i)
}

// Output: 0, 1, 2, 3, 4
```

**Caso real: Buscar elemento**

```go
numeros := []int{10, 20, 30, 40, 50}
buscado := 30

for i, v := range numeros {
    if v == buscado {
        fmt.Printf("Encontrado en índice %d\n", i)
        break       // Sale cuando encuentra
    }
}
```

### continue - Siguiente Iteración

```go
for i := 0; i < 10; i++ {
    if i % 2 == 0 {
        continue    // Salta números pares
    }
    fmt.Println(i)
}

// Output: 1, 3, 5, 7, 9
```

**Comparación:**

```
break:
 Sale COMPLETAMENTE del bucle
 Ejecuta código después del for
 Tipicamente para búsqueda o error

continue:
 Salta ESTA ITERACIÓN
 Continúa con SIGUIENTE iteración
 Tipicamente para filtrado
```

### Bucles Anidados

```go
for i := 1; i <= 3; i++ {
    for j := 1; j <= 3; j++ {
        if j == 2 {
            continue  // Salta j==2 (inner loop)
        }
        fmt.Printf("(%d, %d) ", i, j)
    }
}

// Output: (1, 1) (1, 3) (2, 1) (2, 3) (3, 1) (3, 3)
```

---

## 6.8 Bucles Anidados y Labels

### Labels para break/continue

En bucles anidados, break/continue afecta SOLO al bucle más cercano:

```go
// Sin label: break sale solo del for j
for i := 1; i <= 3; i++ {
    for j := 1; j <= 3; j++ {
        if j == 2 {
            break   // Sale del for j (inner), NO del for i
        }
    }
}

// Con label: break sale del for i
Outer:
for i := 1; i <= 3; i++ {
    for j := 1; j <= 3; j++ {
        if j == 2 {
            break Outer  // Sale del bucle etiquetado
        }
    }
}
```

**Sintaxis:**

```go
Etiqueta:        // Nombre de la etiqueta (PascalCase)
for ... {
    for ... {
        break Etiqueta  // Usa la etiqueta
    }
}
```

**Caso real: Búsqueda en matriz 2D**

```go
matriz := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

buscado := 6

Buscar:
for i := 0; i < len(matriz); i++ {
    for j := 0; j < len(matriz[i]); j++ {
        if matriz[i][j] == buscado {
            fmt.Printf("Encontrado en [%d][%d]\n", i, j)
            break Buscar  // Sale de AMBOS bucles
        }
    }
}
```

---

## 6.9 defer - Ejecución Diferida

### ¿Qué es defer?

`defer` retrasa la ejecución de una sentencia hasta que la función retorna:

```go
func main() {
    fmt.Println("1")
    defer fmt.Println("3")   // Se ejecuta al final
    fmt.Println("2")
}

// Output:
// 1
// 2
// 3
```

### Orden de Ejecución (LIFO - Last In First Out)

```go
func main() {
    defer fmt.Println("A")   // Última en ejecutarse
    defer fmt.Println("B")   // Segunda en ejecutarse
    defer fmt.Println("C")   // Primera en ejecutarse
    fmt.Println("Inicio")
}

// Output:
// Inicio
// C
// B
// A

// Stack de defer:
// Push A → Push B → Push C → Function retorna → Pop C → Pop B → Pop A
```

### Casos de Uso: Limpieza de Recursos

**Cerrar archivo:**

```go
import "os"

func procesarArchivo(nombre string) {
    archivo, _ := os.Open(nombre)
    defer archivo.Close()      // Se ejecuta al final, incluso si hay error
    
    // Procesar archivo
    // ...
    
    // archivo.Close() se ejecuta automáticamente aquí
}
```

**Liberar conexión de BD:**

```go
func consultarBD() {
    db, _ := sql.Open("driver", "source")
    defer db.Close()           // Se ejecuta al final
    
    rows, _ := db.Query("SELECT ...")
    defer rows.Close()         // Se ejecuta al final
    
    // Usar rows
}
```

**Unlock de Mutex:**

```go
import "sync"

var mu sync.Mutex
var contador int

func incrementar() {
    mu.Lock()
    defer mu.Unlock()          // Se ejecuta al final, incluso si hay panic
    
    contador++
}
```

### defer con Funciones

```go
func ejecutar() {
    defer fmt.Println("Defer sin argumentos")
    
    // defer evalúa argumentos INMEDIATAMENTE
    x := 10
    defer fmt.Println("Valor de x:", x)  // Evalúa x aquí (10)
    
    x = 20
    fmt.Println("x dentro:", x)          // Imprime 20
    
    // Al salir:
    // Imprime "x dentro: 20"
    // Imprime "Valor de x: 10" (porque fue evaluado en defer)
    // Imprime "Defer sin argumentos"
}

// Output:
// x dentro: 20
// Valor de x: 10
// Defer sin argumentos
```

### defer con Funciones Anónimas

```go
x := 10

defer func() {
    fmt.Println("Valor de x:", x)  // Ve el valor actual (20)
}()

x = 20

// Output: Valor de x: 20
```

---

## 6.10 Patrones y Buenas Prácticas

### Preferir if/else sobre switch (cuando sea simple)

```go
// ❌ Excesivo
switch estado {
case "activo":
    fmt.Println("Activo")
case "inactivo":
    fmt.Println("Inactivo")
}

// ✅ Mejor
if estado == "activo" {
    fmt.Println("Activo")
} else {
    fmt.Println("Inactivo")
}
```

### Usar switch para casos múltiples

```go
// ✅ Ideal para switch
switch nivel {
case 1, 2, 3:
    fmt.Println("Principiante")
case 4, 5, 6:
    fmt.Println("Intermedio")
case 7, 8, 9, 10:
    fmt.Println("Avanzado")
}
```

### Early Exit (Salida Temprana)

```go
// ❌ Anidado profundo
func validar(x int) {
    if x > 0 {
        if x < 100 {
            if x % 2 == 0 {
                fmt.Println("Válido")
            } else {
                fmt.Println("Error: impar")
            }
        } else {
            fmt.Println("Error: muy grande")
        }
    } else {
        fmt.Println("Error: negativo")
    }
}

// ✅ Early exit (mejor)
func validar(x int) {
    if x <= 0 {
        fmt.Println("Error: negativo")
        return
    }
    if x >= 100 {
        fmt.Println("Error: muy grande")
        return
    }
    if x % 2 != 0 {
        fmt.Println("Error: impar")
        return
    }
    fmt.Println("Válido")
}
```

### Guard Clauses

```go
// ❌ Lógica invertida y anidada
if usuario != nil {
    if usuario.IsAdmin() {
        if usuario.HasPermission("delete") {
            // Hacer algo
        }
    }
}

// ✅ Guard clauses (más claro)
if usuario == nil {
    return
}
if !usuario.IsAdmin() {
    return
}
if !usuario.HasPermission("delete") {
    return
}
// Hacer algo
```

### Bucles Limpios

**Evitar manejo manual de índices:**

```go
// ❌ Complicado
elementos := []string{"a", "b", "c"}
for i := 0; i < len(elementos); i++ {
    fmt.Println(elementos[i])
}

// ✅ Range (más idiomático)
for _, e := range elementos {
    fmt.Println(e)
}
```

**Usar nombres descriptivos:**

```go
// ❌ i, j, k (sin contexto)
for i := 0; i < 10; i++ {
    // ...
}

// ✅ Nombre descriptivo
for intento := 0; intento < 10; intento++ {
    // ...
}

// O para bucles cortos, i está bien
for i := 0; i < len(lista); i++ {
    // ...
}
```

### defer para Cleanup

```go
func procesar() error {
    conexion := conectar()
    defer conexion.Cerrar()      // Siempre se ejecuta
    
    if err := validar(); err != nil {
        return err              // Aún ejecuta defer
    }
    
    procesar()
    return nil                  // Ejecuta defer al retornar
}
```

---

## Ejercicios del Capítulo 6

### Ejercicio 1: Calificador de Notas

Crea programa que:
1. Pida 5 calificaciones
2. Use switch para determinar letra (A, B, C, D, F)
3. Muestre promedio
4. Muestre calificación general
5. Use early exit si hay calificación negativa

### Ejercicio 2: Tabla de Multiplicar

Crea programa que:
1. Pida número
2. Imprima tabla de multiplicar (1-10)
3. Destaque múltiplos de 5 (usa continue)
4. Salta múltiplos de 3
5. Usa for y if/else

### Ejercicio 3: Búsqueda en Matriz

Crea programa que:
1. Cree matriz 3x3
2. Pida número a buscar
3. Use bucles anidados con label
4. Cuando encuentre, imprime posición y sale con break label
5. Si no encuentra, imprime mensaje

### Ejercicio 4: Validador Interactivo

Crea programa que:
1. Menú infinito (for infinito)
2. Pida opción:
   - 1: Validar edad
   - 2: Validar email
   - 3: Salir
3. Use switch y break para salir
4. Usa defer para imprimir "Programa terminado" al final

### Ejercicio 5: Contador con Restricciones

Crea programa que:
1. Cuente del 1 al 100
2. Salta números divisibles por 3 (continue)
3. Detiene en primer número > 50 divisible por 7 (break)
4. Imprime estadísticas al salir con defer
5. Usa labels si necesita control fino

---

**Fin del Capítulo 6**

---

