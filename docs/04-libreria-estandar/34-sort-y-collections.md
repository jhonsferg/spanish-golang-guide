# Capítulo 34: Sort y collections

## Introducción

El package `sort` de Go proporciona herramientas poderosas para ordenar datos de cualquier tipo, desde slices primitivos hasta estructuras complejas. Además, Go ofrece el package `container` con implementaciones de estructuras de datos avanzadas como heaps y listas enlazadas. Este capítulo explora cómo ordenar eficientemente, buscar en datos ordenados y utilizar estructuras de datos especializadas para resolver problemas de rendimiento.

A diferencia de lenguajes como Python donde el ordenamiento es implícito, o C++ donde tienes múltiples opciones de STL, Go ofrece un balance entre flexibilidad y simplicidad mediante sus interfaces bien definidas.

---

## 34.1 Sort Package Fundamentos

### 34.1.1 La Interface sort.Interface

El corazón del ordenamiento en Go es la interface `sort.Interface`. Para que cualquier tipo sea ordenable, debe implementar tres métodos:

```go
type Interface interface {
    // Len devuelve la cantidad de elementos
    Len() int
    
    // Less devuelve true si el elemento en i es menor que en j
    Less(i, j int) bool
    
    // Swap intercambia los elementos en i y j
    Swap(i, j int)
}
```

Esta interfaz es minimalista pero poderosa. Define el contrato que todo algoritmo de ordenamiento en Go espera:

```go
package main

import (
    "fmt"
    "sort"
)

// PersonasPorEdad implementa sort.Interface
type PersonasPorEdad []struct {
    Nombre string
    Edad   int
}

func (p PersonasPorEdad) Len() int {
    return len(p)
}

func (p PersonasPorEdad) Less(i, j int) bool {
    return p[i].Edad < p[j].Edad
}

func (p PersonasPorEdad) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func main() {
    personas := PersonasPorEdad{
        {Nombre: "Alice", Edad: 30},
        {Nombre: "Bob", Edad: 25},
        {Nombre: "Charlie", Edad: 35},
    }
    
    sort.Sort(personas)
    
    for _, p := range personas {
        fmt.Printf("%s: %d años\n", p.Nombre, p.Edad)
    }
    // Output:
    // Bob: 25 años
    // Alice: 30 años
    // Charlie: 35 años
}
```

### 34.1.2 Cómo Funciona Internamente

Go utiliza **Quicksort** como algoritmo principal, con optimizaciones especiales:

```
┌────────────────────────────────────────────────────────┐
│         Algoritmo de Ordenamiento en Go                │
├────────────────────────────────────────────────────────┤
│                                                        │
│  1. Quicksort (principal)                             │
│     ├─ Partición rápida                               │
│     ├─ O(n log n) promedio, O(n²) peor caso          │
│     └─ Optimizado con fallback                        │
│                                                        │
│  2. Heapsort (fallback si recursión es profunda)      │
│     ├─ O(n log n) garantizado                         │
│     └─ Se activa si profundidad > 2*log(n)            │
│                                                        │
│  3. Insertion sort (para segmentos pequeños)          │
│     ├─ O(n²) pero con constantes bajas                │
│     └─ Se usa cuando len < 12                         │
│                                                        │
└────────────────────────────────────────────────────────┘
```

```go
// Ejemplo mostrando la complejidad
func benchmarkOrdenamiento(n int) {
    datos := make([]int, n)
    for i := 0; i < n; i++ {
        datos[i] = rand.Intn(1000)
    }
    
    // Mejor caso: O(n)
    sort.Ints(datos) // Ya está parcialmente ordenado
    
    // Peor caso: O(n log n) con fallback a heapsort
    // Go previene el caso patológico de O(n²)
}
```

### 34.1.3 Comparación: Go vs Otros Lenguajes

```go
// ╔═══════════════════════════════════════════════════════╗
// ║                  Go (InterfaceBasado)                ║
// ╚═══════════════════════════════════════════════════════╝
type Personas []Persona

func (p Personas) Len() int           { return len(p) }
func (p Personas) Less(i, j int) bool { return p[i].Edad < p[j].Edad }
func (p Personas) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

sort.Sort(personas)

// ╔═══════════════════════════════════════════════════════╗
// ║              C++ (Iteradores + Functores)            ║
// ╚═══════════════════════════════════════════════════════╝
// std::sort(personas.begin(), personas.end(),
//     [](const Persona& a, const Persona& b) {
//         return a.edad < b.edad;
//     });

// ╔═══════════════════════════════════════════════════════╗
// ║          Python (Funciones de comparación)            ║
// ╚═══════════════════════════════════════════════════════╝
// personas.sort(key=lambda p: p.edad)
```

### 34.1.4 La Interface sort.Sorter (Go 1.8+)

Go introdujo `sort.Slice` para evitar implementar la interfaz completa:

```go
personas := []Persona{
    {Nombre: "Alice", Edad: 30},
    {Nombre: "Bob", Edad: 25},
    {Nombre: "Charlie", Edad: 35},
}

// Sin implementar sort.Interface
sort.Slice(personas, func(i, j int) bool {
    return personas[i].Edad < personas[j].Edad
})

// Incluso más simple: sort.SliceStable (mantiene orden relativo)
sort.SliceStable(personas, func(i, j int) bool {
    return personas[i].Edad < personas[j].Edad
})
```

---

## 34.2 Ordenar Slices - Tipos Primitivos

### 34.2.1 Ordenar Números Enteros

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    numeros := []int{64, 34, 25, 12, 22, 11, 90}
    
    // Ordenamiento ascendente
    sort.Ints(numeros)
    fmt.Println("Ascendente:", numeros)
    // Output: Ascendente: [11 12 22 25 34 64 90]
    
    // Ordenamiento descendente usando Reverse
    sort.Sort(sort.Reverse(sort.IntSlice(numeros)))
    fmt.Println("Descendente:", numeros)
    // Output: Descendente: [90 64 34 25 22 12 11]
    
    // Búsqueda en slice ordenado
    idx := sort.SearchInts(numeros, 25)
    if idx < len(numeros) && numeros[idx] == 25 {
        fmt.Printf("Encontrado en índice: %d\n", idx)
    }
}
```

### 34.2.2 Ordenar Números Flotantes

```go
func main() {
    flotantes := []float64{3.14, 2.71, 1.41, 1.73, 2.23}
    
    sort.Float64s(flotantes)
    fmt.Println("Flotantes ordenados:", flotantes)
    // Output: Flotantes ordenados: [1.41 1.73 2.23 2.71 3.14]
    
    // Búsqueda binaria en flotantes
    idx := sort.SearchFloat64s(flotantes, 2.0)
    fmt.Printf("Posición para insertar 2.0: %d\n", idx)
    // Esto da la posición donde 2.0 debería ir para mantener orden
}
```

### 34.2.3 Ordenar Strings

```go
func main() {
    palabras := []string{"zebra", "apple", "banana", "cherry", "date"}
    
    // Ordenamiento alfabético
    sort.Strings(palabras)
    fmt.Println("Alfabético:", palabras)
    // Output: Alfabético: [apple banana cherry date zebra]
    
    // Verificar si está ordenado
    esOrdenado := sort.StringsAreSorted(palabras)
    fmt.Println("¿Está ordenado?:", esOrdenado)
    
    // Búsqueda binaria en strings
    idx := sort.SearchStrings(palabras, "cherry")
    if idx < len(palabras) && palabras[idx] == "cherry" {
        fmt.Printf("Encontrado 'cherry' en índice: %d\n", idx)
    }
}
```

### 34.2.4 Verificación de Orden

```go
func main() {
    numeros := []int{1, 2, 3, 4, 5}
    flotantes := []float64{1.1, 1.2, 1.3}
    palabras := []string{"a", "b", "c"}
    
    // Funciones de verificación
    fmt.Println("Ints ordenados:", sort.IntsAreSorted(numeros))
    fmt.Println("Float64s ordenados:", sort.Float64sAreSorted(flotantes))
    fmt.Println("Strings ordenados:", sort.StringsAreSorted(palabras))
    
    // Verificación genérica
    fmt.Println("Genérica:", sort.IsSorted(sort.IntSlice(numeros)))
}
```

### 34.2.5 Tabla Comparativa de Complejidad

| Función | Complejidad | Espacio | Estable | Caso de Uso |
|---------|-------------|---------|---------|------------|
| sort.Ints | O(n log n) | O(log n) | No | Enteros simples |
| sort.Float64s | O(n log n) | O(log n) | No | Flotantes simples |
| sort.Strings | O(n log n) | O(log n) | No | Strings simples |
| sort.Slice | O(n log n) | O(log n) | No | Tipos personalizados |
| sort.SliceStable | O(n log n) | O(n) | Sí | Orden relativo importante |

---

## 34.3 Ordenar Structs - Criterios Personalizados

### 34.3.1 Implementación Básica

```go
package main

import (
    "fmt"
    "sort"
)

type Empleado struct {
    ID       int
    Nombre   string
    Salario  float64
    Años     int
}

// Ordenar por salario (ascendente)
type PorSalario []Empleado

func (p PorSalario) Len() int {
    return len(p)
}

func (p PorSalario) Less(i, j int) bool {
    return p[i].Salario < p[j].Salario
}

func (p PorSalario) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

// Ordenar por antigüedad (años)
type PorAntigüedad []Empleado

func (p PorAntigüedad) Len() int {
    return len(p)
}

func (p PorAntigüedad) Less(i, j int) bool {
    return p[i].Años > p[j].Años // Descendente
}

func (p PorAntigüedad) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func main() {
    empleados := []Empleado{
        {1, "Alice", 50000, 5},
        {2, "Bob", 45000, 8},
        {3, "Charlie", 55000, 3},
        {4, "Diana", 52000, 6},
    }
    
    // Ordenar por salario
    sort.Sort(PorSalario(empleados))
    fmt.Println("Por salario:")
    for _, e := range empleados {
        fmt.Printf("  %s: $%.2f\n", e.Nombre, e.Salario)
    }
    
    // Ordenar por antigüedad
    sort.Sort(PorAntigüedad(empleados))
    fmt.Println("\nPor antigüedad (descendente):")
    for _, e := range empleados {
        fmt.Printf("  %s: %d años\n", e.Nombre, e.Años)
    }
}
```

### 34.3.2 Uso de sort.Slice para Mayor Flexibilidad

```go
func main() {
    empleados := []Empleado{
        {1, "Alice", 50000, 5},
        {2, "Bob", 45000, 8},
        {3, "Charlie", 55000, 3},
        {4, "Diana", 52000, 6},
    }
    
    // Múltiples ordenamientos sin crear tipos especiales
    
    // 1. Por salario descendente
    sort.Slice(empleados, func(i, j int) bool {
        return empleados[i].Salario > empleados[j].Salario
    })
    fmt.Println("Mayor salario primero:")
    for _, e := range empleados {
        fmt.Printf("  %s: $%.2f\n", e.Nombre, e.Salario)
    }
    
    // 2. Por nombre (alfabético)
    sort.Slice(empleados, func(i, j int) bool {
        return empleados[i].Nombre < empleados[j].Nombre
    })
    fmt.Println("\nAlfabético:")
    for _, e := range empleados {
        fmt.Printf("  %s\n", e.Nombre)
    }
    
    // 3. Por años de antigüedad, luego por salario (multi-criterio)
    sort.Slice(empleados, func(i, j int) bool {
        if empleados[i].Años != empleados[j].Años {
            return empleados[i].Años > empleados[j].Años
        }
        return empleados[i].Salario < empleados[j].Salario
    })
    fmt.Println("\nPor antigüedad, luego por salario:")
    for _, e := range empleados {
        fmt.Printf("  %s: %d años, $%.2f\n", e.Nombre, e.Años, e.Salario)
    }
}
```

### 34.3.3 Ordenamiento Estable

El ordenamiento estable mantiene el orden relativo de elementos iguales. Esto es crítico en ciertos casos:

```go
func main() {
    // Simulación: lista de compras con cantidades
    type Item struct {
        Nombre    string
        Cantidad  int
        Prioridad int
    }
    
    items := []Item{
        {"Pan", 2, 1},
        {"Leche", 1, 1},
        {"Queso", 3, 2},
        {"Huevos", 12, 1},
    }
    
    // sort.Slice NO es estable
    sort.Slice(items, func(i, j int) bool {
        return items[i].Prioridad < items[j].Prioridad
    })
    
    fmt.Println("Con sort.Slice (NO estable):")
    for _, item := range items {
        fmt.Printf("  %s (prioridad %d)\n", item.Nombre, item.Prioridad)
    }
    
    // Reiniciar
    items = []Item{
        {"Pan", 2, 1},
        {"Leche", 1, 1},
        {"Queso", 3, 2},
        {"Huevos", 12, 1},
    }
    
    // sort.SliceStable SÍ es estable
    sort.SliceStable(items, func(i, j int) bool {
        return items[i].Prioridad < items[j].Prioridad
    })
    
    fmt.Println("\nCon sort.SliceStable (ESTABLE):")
    for _, item := range items {
        fmt.Printf("  %s (prioridad %d)\n", item.Nombre, item.Prioridad)
    }
}
```

---

## 34.4 Search - Búsqueda Binaria

### 34.4.1 Búsqueda Binaria Básica

La búsqueda binaria es fundamental para trabajar con datos ordenados. Go proporciona tres funciones especializadas:

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    numeros := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
    
    // Buscar número que existe
    idx := sort.SearchInts(numeros, 12)
    fmt.Printf("Índice de 12: %d\n", idx)
    if idx < len(numeros) && numeros[idx] == 12 {
        fmt.Println("¡Encontrado!")
    }
    
    // Buscar número que NO existe
    idx = sort.SearchInts(numeros, 11)
    fmt.Printf("Índice para insertar 11: %d\n", idx)
    // SearchInts devuelve la posición donde el elemento debería ir
    
    // Búsqueda en strings
    palabras := []string{"apple", "banana", "cherry", "date", "elderberry"}
    idx = sort.SearchStrings(palabras, "cherry")
    fmt.Printf("Índice de 'cherry': %d\n", idx)
    
    // Búsqueda en flotantes
    valores := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
    idx = sort.SearchFloat64s(valores, 3.3)
    fmt.Printf("Índice de 3.3: %d\n", idx)
}
```

### 34.4.2 Búsqueda Genérica con sort.Search

Para tipos personalizados, usa `sort.Search`:

```go
func main() {
    empleados := []Empleado{
        {1, "Alice", 45000, 3},
        {2, "Bob", 50000, 5},
        {3, "Charlie", 55000, 7},
        {4, "Diana", 60000, 9},
    }
    
    // Buscar por salario
    buscarSalario := 50000.0
    idx := sort.Search(len(empleados), func(i int) bool {
        return empleados[i].Salario >= buscarSalario
    })
    
    if idx < len(empleados) && empleados[idx].Salario == buscarSalario {
        fmt.Printf("Encontrado empleado con salario $%.2f: %s\n",
            buscarSalario, empleados[idx].Nombre)
    } else {
        fmt.Printf("Posición para insertar: %d\n", idx)
    }
    
    // Búsqueda por rango de edades
    // Encontrar primer empleado con ID >= 3
    idx = sort.Search(len(empleados), func(i int) bool {
        return empleados[i].ID >= 3
    })
    fmt.Printf("Primer empleado con ID >= 3 está en índice: %d\n", idx)
}
```

### 34.4.3 Complejidad y Rendimiento

```
┌────────────────────────────────────────────┐
│      Búsqueda Binaria vs Lineal             │
├────────────────────────────────────────────┤
│                                            │
│  Lineal (O(n)):                            │
│  Para n=1,000,000: ~500,000 comparaciones │
│  Para n=1,000,000,000: ~500,000,000       │
│                                            │
│  Binaria (O(log n)):                       │
│  Para n=1,000,000: ~20 comparaciones      │
│  Para n=1,000,000,000: ~30 comparaciones  │
│                                            │
│  Ganancia: 25,000x más rápida              │
│                                            │
└────────────────────────────────────────────┘
```

```go
func main() {
    // Generador de datos ordenados
    datos := make([]int, 1000000)
    for i := 0; i < len(datos); i++ {
        datos[i] = i * 2 // [0, 2, 4, 6, ...]
    }
    
    // Búsqueda de un elemento que no existe
    buscar := 1999999
    
    // Forma 1: Búsqueda binaria (eficiente)
    idx := sort.SearchInts(datos, buscar)
    // ~20 comparaciones
    
    // Forma 2: Búsqueda lineal (ineficiente)
    encontrado := false
    for i := 0; i < len(datos); i++ {
        if datos[i] == buscar {
            encontrado = true
            break
        }
    }
    // ~500,000 comparaciones
}
```

---

## 34.5 Reverse Ordering - Ordenamiento Inverso

### 34.5.1 El Wrapper sort.Reverse

Go proporciona un wrapper elegante para invertir el orden:

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    // Números en orden descendente
    numeros := []int{3, 1, 4, 1, 5, 9, 2, 6}
    
    // Opción 1: Ordenar ascendente, luego invertir
    sort.Ints(numeros)
    fmt.Println("Ascendente:", numeros)
    
    // Opción 2: Usar sort.Reverse directamente
    sort.Sort(sort.Reverse(sort.IntSlice(numeros)))
    fmt.Println("Descendente:", numeros)
    
    // Lo mismo funciona con strings
    palabras := []string{"zebra", "apple", "banana", "cherry"}
    sort.Sort(sort.Reverse(sort.StringSlice(palabras)))
    fmt.Println("Strings invertidos:", palabras)
    
    // Y con flotantes
    valores := []float64{3.14, 2.71, 1.41, 1.73}
    sort.Sort(sort.Reverse(sort.Float64Slice(valores)))
    fmt.Println("Flotantes invertidos:", valores)
}
```

### 34.5.2 Reverse con Structs Personalizados

```go
func main() {
    empleados := []Empleado{
        {1, "Alice", 50000, 5},
        {2, "Bob", 45000, 8},
        {3, "Charlie", 55000, 3},
        {4, "Diana", 52000, 6},
    }
    
    // Crear tipo para ordenamiento invertido
    type PorSalarioDesc []Empleado
    
    func (p PorSalarioDesc) Len() int {
        return len(p)
    }
    
    func (p PorSalarioDesc) Less(i, j int) bool {
        // Invertir la lógica: j < i
        return p[j].Salario < p[i].Salario
    }
    
    func (p PorSalarioDesc) Swap(i, j int) {
        p[i], p[j] = p[j], p[i]
    }
    
    sort.Sort(PorSalarioDesc(empleados))
    
    fmt.Println("Empleados por mayor salario:")
    for _, e := range empleados {
        fmt.Printf("  %s: $%.2f\n", e.Nombre, e.Salario)
    }
}
```

### 34.5.3 Reverse Genérico

```go
func main() {
    empleados := []Empleado{
        {1, "Alice", 50000, 5},
        {2, "Bob", 45000, 8},
        {3, "Charlie", 55000, 3},
    }
    
    // Ordenar por salario (ascendente)
    sort.Slice(empleados, func(i, j int) bool {
        return empleados[i].Salario < empleados[j].Salario
    })
    fmt.Println("Salario ascendente:")
    for _, e := range empleados {
        fmt.Printf("  %s: $%.2f\n", e.Nombre, e.Salario)
    }
    
    // Ordenar por salario (descendente)
    sort.Slice(empleados, func(i, j int) bool {
        return empleados[i].Salario > empleados[j].Salario  // Invertir condición
    })
    fmt.Println("Salario descendente:")
    for _, e := range empleados {
        fmt.Printf("  %s: $%.2f\n", e.Nombre, e.Salario)
    }
}
```

---

## 34.6 Heap Package - Montículos y Colas de Prioridad

### 34.6.1 Conceptos de Heap

Un heap es un árbol binario completo que satisface la propiedad heap:

```
Min-Heap (cada padre ≤ sus hijos):
          1
        /   \
       2     3
      / \   /
     4   5 6

Max-Heap (cada padre ≥ sus hijos):
          6
        /   \
       5     3
      / \   /
     4   1 2
```

Go implementa heaps como arrays donde:
- Padre de índice i: (i-1)/2
- Hijo izquierdo de i: 2*i+1
- Hijo derecho de i: 2*i+2

### 34.6.2 Implementar heap.Interface

```go
package main

import (
    "fmt"
    "container/heap"
)

// MinHeapInt implementa heap.Interface para enteros (Min-Heap)
type MinHeapInt []int

func (h MinHeapInt) Len() int {
    return len(h)
}

func (h MinHeapInt) Less(i, j int) bool {
    return h[i] < h[j]  // Min-heap
}

func (h MinHeapInt) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *MinHeapInt) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *MinHeapInt) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func main() {
    h := &MinHeapInt{3, 2, 1, 7, 8}
    heap.Init(h)
    
    fmt.Println("Heap inicial:", *h)
    
    // Insertar elementos
    heap.Push(h, 4)
    heap.Push(h, 0)
    
    fmt.Println("Después de insertar 4 y 0:", *h)
    
    // Extraer en orden (menor primero)
    fmt.Println("\nExtrayendo elementos en orden:")
    for h.Len() > 0 {
        fmt.Printf("%d ", heap.Pop(h))
    }
    fmt.Println()
}
```

### 34.6.3 Max-Heap Implementación

```go
// MaxHeapInt implementa heap.Interface para enteros (Max-Heap)
type MaxHeapInt []int

func (h MaxHeapInt) Len() int {
    return len(h)
}

func (h MaxHeapInt) Less(i, j int) bool {
    return h[i] > h[j]  // Invertir para Max-heap
}

func (h MaxHeapInt) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *MaxHeapInt) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *MaxHeapInt) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func main() {
    h := &MaxHeapInt{3, 2, 1, 7, 8}
    heap.Init(h)
    
    fmt.Println("Max-Heap:", *h)
    
    heap.Push(h, 10)
    fmt.Println("Después de insertar 10:", *h)
    
    // Extraer en orden (mayor primero)
    fmt.Println("Extrayendo (mayor primero):")
    for h.Len() > 0 {
        fmt.Printf("%d ", heap.Pop(h))
    }
    fmt.Println()
}
```

### 34.6.4 Cola de Prioridad Práctica

```go
package main

import (
    "container/heap"
    "fmt"
)

type Tarea struct {
    Descripción string
    Prioridad   int  // Mayor = más urgente
    Índice      int  // Necesario para heap
}

type ColaTareas []*Tarea

func (c ColaTareas) Len() int { return len(c) }

func (c ColaTareas) Less(i, j int) bool {
    // Max-heap: tarea con mayor prioridad primero
    return c[i].Prioridad > c[j].Prioridad
}

func (c ColaTareas) Swap(i, j int) {
    c[i], c[j] = c[j], c[i]
    c[i].Índice = i
    c[j].Índice = j
}

func (c *ColaTareas) Push(x interface{}) {
    tarea := x.(*Tarea)
    tarea.Índice = len(*c)
    *c = append(*c, tarea)
}

func (c *ColaTareas) Pop() interface{} {
    old := *c
    n := len(old)
    tarea := old[n-1]
    tarea.Índice = -1
    *c = old[0 : n-1]
    return tarea
}

func main() {
    cola := &ColaTareas{}
    heap.Init(cola)
    
    // Agregar tareas
    tareas := []string{"Código", "Tests", "Documentación", "Deploy"}
    prioridades := []int{2, 3, 1, 4}
    
    for i, desc := range tareas {
        heap.Push(cola, &Tarea{desc, prioridades[i], 0})
    }
    
    fmt.Println("Procesando tareas por prioridad:")
    for cola.Len() > 0 {
        tarea := heap.Pop(cola).(*Tarea)
        fmt.Printf("  [Prioridad %d] %s\n", tarea.Prioridad, tarea.Descripción)
    }
}
```

### 34.6.5 Operaciones de Heap

```go
func main() {
    h := &MinHeapInt{5, 3, 7, 1, 9, 2}
    heap.Init(h)
    
    // 1. Init: Construye el heap (O(n))
    fmt.Println("Heap:", *h)
    
    // 2. Push: Agrega elemento y restaura propiedad (O(log n))
    heap.Push(h, 0)
    fmt.Println("Después Push(0):", *h)
    
    // 3. Pop: Extrae mínimo y restaura propiedad (O(log n))
    min := heap.Pop(h)
    fmt.Println("Pop devolvió:", min)
    fmt.Println("Heap actualizado:", *h)
    
    // 4. Remove: Elimina elemento en índice i (O(log n))
    if h.Len() > 0 {
        heap.Remove(h, 0)
        fmt.Println("Después Remove(0):", *h)
    }
    
    // 5. Fix: Restaura propiedad después de cambiar elemento (O(log n))
    (*h)[0] = 100
    heap.Fix(h, 0)
    fmt.Println("Después Fix(0) con valor 100:", *h)
}
```

---

## 34.7 Comparadores Personalizados - Ordenamiento Multi-Criterio

### 34.7.1 Ordenamiento por Múltiples Campos

```go
package main

import (
    "fmt"
    "sort"
)

type Producto struct {
    ID       int
    Nombre   string
    Precio   float64
    Stock    int
    Categoría string
}

func main() {
    productos := []Producto{
        {1, "Laptop", 999.99, 5, "Electrónica"},
        {2, "Mouse", 25.50, 100, "Accesorios"},
        {3, "Teclado", 75.00, 50, "Accesorios"},
        {4, "Monitor", 299.99, 15, "Electrónica"},
        {5, "Webcam", 85.00, 30, "Accesorios"},
    }
    
    // Ordenar por categoría, luego por precio (descendente)
    sort.Slice(productos, func(i, j int) bool {
        if productos[i].Categoría != productos[j].Categoría {
            return productos[i].Categoría < productos[j].Categoría
        }
        return productos[i].Precio > productos[j].Precio
    })
    
    fmt.Println("Ordenado por categoría y precio:")
    for _, p := range productos {
        fmt.Printf("  [%s] %s: $%.2f\n", p.Categoría, p.Nombre, p.Precio)
    }
    
    // Ordenar por disponibilidad (stock bajo primero), luego por nombre
    sort.Slice(productos, func(i, j int) bool {
        if productos[i].Stock != productos[j].Stock {
            return productos[i].Stock < productos[j].Stock
        }
        return productos[i].Nombre < productos[j].Nombre
    })
    
    fmt.Println("\nOdenado por stock bajo, luego por nombre:")
    for _, p := range productos {
        fmt.Printf("  %s (Stock: %d)\n", p.Nombre, p.Stock)
    }
}
```

### 34.7.2 Comparador Complejo con Lógica de Negocio

```go
type Usuario struct {
    ID           int
    Nombre       string
    Edad         int
    FechaRegistro string
    Premium      bool
}

func main() {
    usuarios := []Usuario{
        {1, "Alice", 28, "2023-01-15", true},
        {2, "Bob", 35, "2022-06-20", false},
        {3, "Charlie", 28, "2023-02-10", true},
        {4, "Diana", 45, "2021-12-01", true},
        {5, "Eve", 28, "2023-01-20", false},
    }
    
    // Ordenamiento por prioridad de negocio:
    // 1. Premium primero
    // 2. Entre los mismos, más nuevos primero (mejor retención reciente)
    // 3. Si empate, por edad (usuarios más jóvenes tienden a gastar más)
    
    sort.Slice(usuarios, func(i, j int) bool {
        // Criterio 1: Premium
        if usuarios[i].Premium != usuarios[j].Premium {
            return usuarios[i].Premium // true > false
        }
        
        // Criterio 2: Fecha registro (más nuevos primero)
        if usuarios[i].FechaRegistro != usuarios[j].FechaRegistro {
            return usuarios[i].FechaRegistro > usuarios[j].FechaRegistro
        }
        
        // Criterio 3: Edad (más jóvenes primero)
        return usuarios[i].Edad < usuarios[j].Edad
    })
    
    fmt.Println("Usuarios ordenados por estrategia de negocio:")
    for _, u := range usuarios {
        premium := "Regular"
        if u.Premium {
            premium = "Premium"
        }
        fmt.Printf("  [%s] %s (Edad: %d, Registro: %s)\n",
            premium, u.Nombre, u.Edad, u.FechaRegistro)
    }
}
```

### 34.7.3 Comparador con Función Helper

```go
// Función que devuelve un comparador personalizado
func crearComparador(criterio string) func(i, j int) bool {
    return func(i, j int) bool {
        switch criterio {
        case "nombre":
            return productos[i].Nombre < productos[j].Nombre
        case "precio-asc":
            return productos[i].Precio < productos[j].Precio
        case "precio-desc":
            return productos[i].Precio > productos[j].Precio
        case "stock":
            return productos[i].Stock < productos[j].Stock
        default:
            return productos[i].ID < productos[j].ID
        }
    }
}

func main() {
    productos := []Producto{
        {1, "Laptop", 999.99, 5, "Electrónica"},
        {2, "Mouse", 25.50, 100, "Accesorios"},
        {3, "Teclado", 75.00, 50, "Accesorios"},
    }
    
    for _, criterio := range []string{"nombre", "precio-asc", "stock"} {
        sort.Slice(productos, crearComparador(criterio))
        fmt.Printf("Ordenado por %s:\n", criterio)
        for _, p := range productos {
            fmt.Printf("  %s\n", p.Nombre)
        }
        fmt.Println()
    }
}
```

---

## 34.8 Performance - Rendimiento y Optimización

### 34.8.1 Análisis de Algoritmos Utilizados

```
┌──────────────────────────────────────────────────────────┐
│         Algoritmos de Sorting en Go                      │
├──────────────────────────────────────────────────────────┤
│                                                          │
│ 1. QUICKSORT (Principal)                                │
│    - Promedio: O(n log n)                               │
│    - Espacio: O(log n) recursión                        │
│    - No estable por defecto                             │
│    - Ventaja: cache-friendly, rápido en práctica        │
│                                                          │
│ 2. HEAPSORT (Fallback)                                  │
│    - Garantizado: O(n log n)                            │
│    - Espacio: O(1)                                      │
│    - No estable                                         │
│    - Activado si profundidad > 2*log(n)                 │
│                                                          │
│ 3. INSERTION SORT (Segmentos pequeños)                  │
│    - Para n < 12: O(n²)                                 │
│    - Espacio: O(1)                                      │
│    - Muy rápido para datos pequeños                     │
│    - Parcialmente ordenado: O(n)                        │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

### 34.8.2 Benchmarking de Ordenamiento

```go
package main

import (
    "fmt"
    "math/rand"
    "sort"
    "testing"
    "time"
)

func benchmarkSort(nombre string, datos []int, b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Copiar datos para cada iteración
        copia := make([]int, len(datos))
        copy(copia, datos)
        sort.Ints(copia)
    }
}

func BenchmarkSort1K(b *testing.B) {
    datos := generarAleatorios(1000)
    benchmarkSort("1K", datos, b)
}

func BenchmarkSort10K(b *testing.B) {
    datos := generarAleatorios(10000)
    benchmarkSort("10K", datos, b)
}

func BenchmarkSort100K(b *testing.B) {
    datos := generarAleatorios(100000)
    benchmarkSort("100K", datos, b)
}

// Ejecutar con: go test -bench=. -benchmem
// Resultado esperado:
// BenchmarkSort1K-8      10000    100000 ns/op    0 B/op    0 allocs/op
// BenchmarkSort10K-8     1000   1000000 ns/op    0 B/op    0 allocs/op
// BenchmarkSort100K-8     100  10000000 ns/op    0 B/op    0 allocs/op

func generarAleatorios(n int) []int {
    datos := make([]int, n)
    for i := 0; i < n; i++ {
        datos[i] = rand.Intn(1000000)
    }
    return datos
}

func main() {
    // Comparación práctica de rendimiento
    fmt.Println("Comparación de rendimiento:")
    
    for _, n := range []int{1000, 10000, 100000} {
        datos := generarAleatorios(n)
        
        inicio := time.Now()
        sort.Ints(datos)
        duracion := time.Since(inicio)
        
        fmt.Printf("Ordenar %d elementos: %v\n", n, duracion)
    }
}
```

### 34.8.3 Optimizaciones Prácticas

```go
func main() {
    // Datos aleatorios
    aleatoria := generarAleatorios(100000)
    
    // Datos parcialmente ordenados
    parcial := make([]int, 100000)
    for i := 0; i < 100000; i++ {
        if i%10 == 0 {
            parcial[i] = rand.Intn(1000000)
        } else {
            parcial[i] = i
        }
    }
    
    // Datos casi ordenados
    sorted := make([]int, 100000)
    for i := 0; i < 100000; i++ {
        sorted[i] = i
    }
    // Desordenar 5 elementos
    for i := 0; i < 5; i++ {
        sorted[rand.Intn(100000)] = rand.Intn(1000000)
    }
    
    // Benchmark
    tiempo := func(nombre string, datos []int) {
        copia := make([]int, len(datos))
        copy(copia, datos)
        
        inicio := time.Now()
        sort.Ints(copia)
        duracion := time.Since(inicio)
        
        fmt.Printf("%s: %v\n", nombre, duracion)
    }
    
    tiempo("Aleatorio", aleatoria)
    tiempo("Parcial", parcial)
    tiempo("Casi ordenado", sorted)
    
    // Resultado: casi ordenado es ~100x más rápido
}
```

### 34.8.4 Estabilidad del Sort

```go
func main() {
    type Registro struct {
        Valor int
        Orden int // Posición original
    }
    
    datos := []Registro{
        {1, 0},
        {3, 1},
        {1, 2},
        {2, 3},
        {1, 4},
    }
    
    // sort.Slice NO es estable
    sort.Slice(datos, func(i, j int) bool {
        return datos[i].Valor < datos[j].Valor
    })
    
    fmt.Println("sort.Slice (NO estable):")
    for _, r := range datos {
        fmt.Printf("  Valor: %d, Orden original: %d\n", r.Valor, r.Orden)
    }
    
    // Reiniciar
    datos = []Registro{
        {1, 0},
        {3, 1},
        {1, 2},
        {2, 3},
        {1, 4},
    }
    
    // sort.SliceStable SÍ es estable
    sort.SliceStable(datos, func(i, j int) bool {
        return datos[i].Valor < datos[j].Valor
    })
    
    fmt.Println("\nsort.SliceStable (ESTABLE):")
    for _, r := range datos {
        fmt.Printf("  Valor: %d, Orden original: %d\n", r.Valor, r.Orden)
    }
}
```

---

## 34.9 Container/Ring - Lista Circular Doblemente Enlazada

### 34.9.1 Estructura de Ring

```go
package main

import (
    "container/ring"
    "fmt"
)

func main() {
    // Crear un ring de 5 elementos
    r := ring.New(5)
    
    // Llenar con valores
    for i := 1; i <= 5; i++ {
        r.Value = i * 10
        r = r.Next()
    }
    
    // Iterar a través del ring
    fmt.Println("Valores en el ring:")
    r.Do(func(v interface{}) {
        fmt.Printf("  %d\n", v.(int))
    })
    
    // Ring es circular - puedes seguir avanzando indefinidamente
    fmt.Println("\nAvanzar 7 pasos (envuelve):")
    for i := 0; i < 7; i++ {
        fmt.Printf("  Paso %d: %d\n", i+1, r.Value.(int))
        r = r.Next()
    }
}
```

### 34.9.2 Operaciones en Ring

```go
func main() {
    // Crear ring con letras
    r := ring.New(4)
    vals := []string{"A", "B", "C", "D"}
    
    p := r
    for _, v := range vals {
        p.Value = v
        p = p.Next()
    }
    
    fmt.Println("Ring original:")
    r.Do(func(v interface{}) {
        fmt.Print(v.(string) + " ")
    })
    fmt.Println()
    
    // Linkear dos rings
    r2 := ring.New(2)
    r2.Value = "X"
    r2.Next().Value = "Y"
    
    fmt.Println("Ring 2 antes de linkear:")
    r2.Do(func(v interface{}) {
        fmt.Print(v.(string) + " ")
    })
    fmt.Println()
    
    // Link: conecta this.Next() con s
    r.Link(r2)
    
    fmt.Println("Después de linkear (size = 6):")
    r.Do(func(v interface{}) {
        fmt.Print(v.(string) + " ")
    })
    fmt.Println()
    
    // Unlink: extrae n elementos del ring
    r3 := r.Unlink(2)
    
    fmt.Println("Ring después de Unlink(2):")
    r.Do(func(v interface{}) {
        fmt.Print(v.(string) + " ")
    })
    fmt.Println()
    
    fmt.Println("Ring extraído (r3):")
    r3.Do(func(v interface{}) {
        fmt.Print(v.(string) + " ")
    })
    fmt.Println()
}
```

### 34.9.3 Casos de Uso de Ring

```go
// Caso 1: Buffer circular (Round-robin)
func main() {
    // Buffer de tamaño 3 para procesar en rotación
    buffer := ring.New(3)
    for i := 1; i <= 3; i++ {
        buffer.Value = fmt.Sprintf("Buffer[%d]", i)
        buffer = buffer.Next()
    }
    
    fmt.Println("Distribución round-robin de tareas:")
    for tarea := 1; tarea <= 7; tarea++ {
        fmt.Printf("  Tarea %d -> %s\n", tarea, buffer.Value.(string))
        buffer = buffer.Next()
    }
}

// Caso 2: Lista de reproducción circular
type Canción struct {
    Nombre   string
    Artista  string
}

func main() {
    playlist := ring.New(3)
    canciones := []Canción{
        {"Bohemian Rhapsody", "Queen"},
        {"Imagine", "John Lennon"},
        {"Stairway to Heaven", "Led Zeppelin"},
    }
    
    p := playlist
    for _, cancion := range canciones {
        p.Value = cancion
        p = p.Next()
    }
    
    fmt.Println("Reproducción en bucle (5 canciones):")
    for i := 0; i < 5; i++ {
        cancion := playlist.Value.(Canción)
        fmt.Printf("  Reproduciendo: %s - %s\n", cancion.Nombre, cancion.Artista)
        playlist = playlist.Next()
    }
}
```

---

## 34.10 Container/List - Lista Enlazada Doblemente Vinculada

### 34.10.1 Operaciones Básicas

```go
package main

import (
    "container/list"
    "fmt"
)

func main() {
    // Crear lista
    l := list.New()
    
    // Agregar elementos
    e1 := l.PushBack("A")  // Devuelve *list.Element
    e2 := l.PushBack("B")
    e3 := l.PushBack("C")
    
    fmt.Println("Lista después de PushBack (A, B, C):")
    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Print(e.Value.(string) + " ")
    }
    fmt.Println()
    
    // Agregar al frente
    l.PushFront("Z")
    
    fmt.Println("Después de PushFront(Z):")
    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Print(e.Value.(string) + " ")
    }
    fmt.Println()
    
    // Insertar en posición específica
    l.InsertAfter("B1", e1)   // Después de "A"
    l.InsertBefore("B2", e2)  // Antes de "B"
    
    fmt.Println("Después de inserciones:")
    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Print(e.Value.(string) + " ")
    }
    fmt.Println()
    
    // Eliminar elementos
    l.Remove(e3)  // Eliminar "C"
    
    fmt.Println("Después de Remove(C):")
    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Print(e.Value.(string) + " ")
    }
    fmt.Println()
    
    // Moverse hacia atrás
    fmt.Println("\nIteración hacia atrás:")
    for e := l.Back(); e != nil; e = e.Prev() {
        fmt.Print(e.Value.(string) + " ")
    }
    fmt.Println()
}
```

### 34.10.2 Operaciones Avanzadas

```go
func main() {
    l := list.New()
    
    // Agregar múltiples elementos
    for i := 1; i <= 5; i++ {
        l.PushBack(i)
    }
    
    // Longitud
    fmt.Printf("Longitud: %d\n", l.Len())
    
    // Buscar elemento
    var encontrado *list.Element
    for e := l.Front(); e != nil; e = e.Next() {
        if e.Value.(int) == 3 {
            encontrado = e
            break
        }
    }
    
    if encontrado != nil {
        fmt.Printf("Encontrado: %d\n", encontrado.Value.(int))
    }
    
    // Mover elemento
    if encontrado != nil {
        l.MoveToFront(encontrado)
        fmt.Print("Después de MoveToFront(3): ")
        for e := l.Front(); e != nil; e = e.Next() {
            fmt.Print(e.Value.(int), " ")
        }
        fmt.Println()
        
        l.MoveToBack(encontrado)
        fmt.Print("Después de MoveToBack(3): ")
        for e := l.Front(); e != nil; e = e.Next() {
            fmt.Print(e.Value.(int), " ")
        }
        fmt.Println()
    }
    
    // Vaciar lista
    l.Init()
    fmt.Printf("Longitud después de Init(): %d\n", l.Len())
}
```

### 34.10.3 Casos de Uso - LRU Cache

```go
type LRUCache struct {
    Capacidad  int
    Cache      map[string]*list.Element
    ListaReciente *list.List
}

type Entrada struct {
    Clave   string
    Valor   interface{}
}

func NewLRUCache(capacidad int) *LRUCache {
    return &LRUCache{
        Capacidad:  capacidad,
        Cache:      make(map[string]*list.Element),
        ListaReciente: list.New(),
    }
}

func (lru *LRUCache) Get(clave string) interface{} {
    if elem, existe := lru.Cache[clave]; existe {
        // Mover a frente (más reciente)
        lru.ListaReciente.MoveToFront(elem)
        return elem.Value.(Entrada).Valor
    }
    return nil
}

func (lru *LRUCache) Put(clave string, valor interface{}) {
    if elem, existe := lru.Cache[clave]; existe {
        // Actualizar valor y mover a frente
        entrada := elem.Value.(Entrada)
        entrada.Valor = valor
        elem.Value = entrada
        lru.ListaReciente.MoveToFront(elem)
        return
    }
    
    // Nuevo elemento
    entrada := Entrada{clave, valor}
    elem := lru.ListaReciente.PushFront(entrada)
    lru.Cache[clave] = elem
    
    // Si excede capacidad, eliminar el menos recientemente usado
    if lru.ListaReciente.Len() > lru.Capacidad {
        elem := lru.ListaReciente.Back()
        lru.ListaReciente.Remove(elem)
        delete(lru.Cache, elem.Value.(Entrada).Clave)
    }
}

func main() {
    cache := NewLRUCache(3)
    
    cache.Put("a", 1)
    cache.Put("b", 2)
    cache.Put("c", 3)
    
    fmt.Println("Get('b'):", cache.Get("b"))
    
    cache.Put("d", 4)  // Elimina "a" (menos recientemente usado)
    
    fmt.Println("Get('a'):", cache.Get("a"))  // nil (fue eliminado)
    fmt.Println("Get('c'):", cache.Get("c"))  // 3
}
```

### 34.10.4 Casos de Uso - Cola (Deque)

```go
type Cola struct {
    items *list.List
}

func NewCola() *Cola {
    return &Cola{items: list.New()}
}

func (c *Cola) Enqueue(valor interface{}) {
    c.items.PushBack(valor)
}

func (c *Cola) Dequeue() interface{} {
    if c.items.Len() == 0 {
        return nil
    }
    elem := c.items.Front()
    c.items.Remove(elem)
    return elem.Value
}

func (c *Cola) IsEmpty() bool {
    return c.items.Len() == 0
}

func (c *Cola) Size() int {
    return c.items.Len()
}

func main() {
    cola := NewCola()
    
    cola.Enqueue("Cliente 1")
    cola.Enqueue("Cliente 2")
    cola.Enqueue("Cliente 3")
    
    fmt.Println("Procesando cola:")
    for !cola.IsEmpty() {
        cliente := cola.Dequeue()
        fmt.Printf("  Atendiendo: %s\n", cliente.(string))
    }
}
```

---

## 34.11 Buenas Prácticas y Patrones

### 34.11.1 Cuándo Usar Cada Estructura

```
┌─────────────────────────────────────────────────────────────┐
│           Comparativa de Estructuras de Datos               │
├──────────────┬──────────┬──────────┬──────────┬─────────────┤
│ Operación    │ Array    │ Slice    │ List     │ Heap        │
├──────────────┼──────────┼──────────┼──────────┼─────────────┤
│ Acceso       │ O(1)     │ O(1)     │ O(n)     │ O(n)        │
│ Búsqueda     │ O(log n) │ O(log n) │ O(n)     │ O(n)        │
│ Inserción    │ O(n)     │ O(n)     │ O(1)*    │ O(log n)    │
│ Eliminación  │ O(n)     │ O(n)     │ O(1)*    │ O(log n)    │
│ Espacio      │ Fijo     │ Dinámico │ Dinámico │ Dinámico    │
├──────────────┼──────────┼──────────┼──────────┼─────────────┤
│ Mejor para   │ Cache    │ General  │ Colas    │ Prioridades │
└──────────────┴──────────┴──────────┴──────────┴─────────────┘
* Con referencia al elemento
```

### 34.11.2 Selector de Algoritmo

```go
// Función auxiliar para elegir estrategia de ordenamiento
func Ordenar(datos interface{}, criterio string) {
    switch d := datos.(type) {
    case []int:
        if criterio == "rápido" {
            sort.Ints(d)
        } else if criterio == "estable" {
            sort.SliceStable(d, func(i, j int) bool {
                return d[i] < d[j]
            })
        }
    case []string:
        sort.Strings(d)
    case []float64:
        sort.Float64s(d)
    default:
        // Para tipos complejos, usar sort.Slice
        sort.Slice(d, func(i, j int) bool {
            // Comparación personalizada
            return false
        })
    }
}

func SeleccionarEstructura(caso string) string {
    switch caso {
    case "acceso_rapido":
        return "Array o Slice"
    case "insertaconFrecuente":
        return "List (LinkedList)"
    case "prioridades":
        return "Heap"
    case "circular":
        return "Ring"
    case "lru_cache":
        return "List + HashMap"
    default:
        return "Slice"
    }
}
```

### 34.11.3 Antipatrones Comunes

```go
// ❌ ANTIPATRÓN 1: Comparador incorrecto
func main() {
    datos := []int{3, 1, 4, 1, 5, 9}
    
    // Incorrecto: El comparador debe satisfacer:
    // 1. Irreflexividad: Less(i, i) debe ser false
    // 2. Transitividad: Si Less(i, j) y Less(j, k) entonces Less(i, k)
    
    // ❌ Esto causará comportamiento indefinido
    sort.Slice(datos, func(i, j int) bool {
        return datos[i] <= datos[j]  // Violó irreflexividad
    })
}

// ✓ CORRECTO
func main() {
    datos := []int{3, 1, 4, 1, 5, 9}
    
    // ✓ Correcto: Cumple propiedades de orden total
    sort.Slice(datos, func(i, j int) bool {
        return datos[i] < datos[j]  // Irreflexivo
    })
}

// ❌ ANTIPATRÓN 2: Heap property violado
type IncorrectHeap []int

func (h IncorrectHeap) Len() int { return len(h) }
func (h IncorrectHeap) Less(i, j int) bool {
    return h[i] < h[j]  // Min-heap correcto
}
func (h IncorrectHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}
func (h *IncorrectHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}
func (h *IncorrectHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// ❌ Esto violará la propiedad heap
func main() {
    h := &IncorrectHeap{1, 2, 3}
    // heap.Init() debería llamarse antes de otros operaciones
    heap.Push(h, 0)
    // Sin Init(), la propiedad heap no se garantiza
}

// ✓ CORRECTO
func main() {
    h := &IncorrectHeap{1, 2, 3}
    heap.Init(h)  // Construir el heap correctamente
    heap.Push(h, 0)
}

// ❌ ANTIPATRÓN 3: Modificar slice durante iteración
func main() {
    datos := []int{1, 2, 3, 4, 5}
    
    // ❌ Peligroso: modificar slice mientras se itera
    for i := 0; i < len(datos); i++ {
        if datos[i] == 3 {
            datos = append(datos[:i], datos[i+1:]...)  // Cambia length
        }
    }
}

// ✓ CORRECTO: usar variable temporal
func main() {
    datos := []int{1, 2, 3, 4, 5}
    
    // ✓ Mejor: construir nuevo slice
    var resultado []int
    for _, v := range datos {
        if v != 3 {
            resultado = append(resultado, v)
        }
    }
    datos = resultado
}
```

### 34.11.4 Comparación con Otros Lenguajes

```go
// ╔═══════════════════════════════════════════════════════╗
// ║                    Go (Simple)                        ║
// ╚═══════════════════════════════════════════════════════╝
sort.Slice(empleados, func(i, j int) bool {
    return empleados[i].Salario < empleados[j].Salario
})

// ╔═══════════════════════════════════════════════════════╗
// ║              C++ (STL - Más verboso)                  ║
// ╚═══════════════════════════════════════════════════════╝
// std::sort(empleados.begin(), empleados.end(),
//     [](const Empleado& a, const Empleado& b) {
//         return a.salario < b.salario;
//     });

// ╔═══════════════════════════════════════════════════════╗
// ║          Python (Más conciso pero lento)              ║
// ╚═══════════════════════════════════════════════════════╝
// empleados.sort(key=lambda e: e.salario)

// ╔═══════════════════════════════════════════════════════╗
// ║           Java (Más boilerplate)                      ║
// ╚═══════════════════════════════════════════════════════╝
// Collections.sort(empleados,
//     (a, b) -> Double.compare(a.salario, b.salario));
```

---

## Ejercicios Progresivos

### Ejercicio 1: Ordenar Structs con Múltiples Criterios ⭐

```go
/*
Problema:
Dada una lista de estudiantes con campos (Nombre, Calificación, Fecha Inscripción),
implementa un sistema de ordenamiento que permita:

1. Por calificación (descendente), luego por nombre (alfabético)
2. Por fecha de inscripción (más nuevos primero)
3. Por calificación (ascendente) solo si la calificación >= 70

Requisitos:
- Implementar sort.Interface
- Usar sort.Slice para cada criterio
- Mostrar resultados de cada ordenamiento
*/

package main

import (
    "fmt"
    "sort"
    "time"
)

type Estudiante struct {
    Nombre       string
    Calificación float64
    FechaIns     time.Time
}

// TODO: Implementar sort.Interface
// TODO: Implementar cada criterio de ordenamiento
// TODO: Crear función main() con ejemplos

```

### Ejercicio 2: Búsqueda Binaria Personalizada ⭐⭐

```go
/*
Problema:
Implementar un sistema de búsqueda para encontrar el "rango" de elementos
dentro de un rango de valores en un array ordenado.

Requisitos:
- Buscar el primer elemento >= valor_mínimo
- Buscar el último elemento <= valor_máximo
- Retornar slice con elementos en el rango
- Comparar performance con búsqueda lineal

Ejemplo:
BuscarEnRango([1,2,3,4,5,6,7,8,9], 3, 7) → [3,4,5,6,7]
*/

package main

import (
    "fmt"
    "sort"
)

func BuscarEnRango(datos []int, min, max int) []int {
    // TODO: Implementar
    return nil
}

func main() {
    // TODO: Crear test cases
}
```

### Ejercicio 3: Priority Queue usando Heap ⭐⭐

```go
/*
Problema:
Implementar un sistema de eventos con prioridades. Cada evento tiene:
- ID (único)
- Descripción
- Prioridad (1-5, donde 5 es crítico)
- Timestamp

Requisitos:
- Implementar heap.Interface
- Push: Agregar evento con prioridad
- Pop: Procesar evento más urgente
- Update: Cambiar prioridad de evento existente
- Las tarea críticas se procesan primero

Ejemplo de uso:
queue.Push(Evento{1, "Login fallido", 5})
queue.Push(Evento{2, "Email sent", 1})
queue.Pop()  // Devuelve evento con ID 1
*/

package main

import (
    "container/heap"
    "fmt"
    "time"
)

type Evento struct {
    ID          int
    Descripción string
    Prioridad   int
    Timestamp   time.Time
    Índice      int
}

type EventoQueue []*Evento

// TODO: Implementar heap.Interface
// TODO: Implementar operación Update
// TODO: Crear función main() con ejemplos de uso
```

### Ejercicio 4: LRU Cache Extendido ⭐⭐

```go
/*
Problema:
Extender la implementación de LRU Cache para que incluya:
- Stats: Número de hits, misses y total de accesos
- Eviction Policy: Opción entre LRU, LFU (menos frecuente)
- TTL: Expiración automática de items
- Predictor: Sugerir si el próximo acceso será hit o miss

Requisitos:
- Usar container/list para LRU
- Mantener estadísticas
- Implementar TTL con goroutines
- Método Stats() que devuelva mapa con estadísticas
*/

package main

import (
    "container/list"
    "fmt"
    "time"
)

type LRUCacheAvanzado struct {
    Capacidad     int
    Cache         map[string]*list.Element
    ListaReciente *list.List
    // TODO: Agregar campos para stats y TTL
}

// TODO: Implementar métodos Get, Put, Stats
// TODO: Implementar expiración con timer
// TODO: Crear función main() con demo completo
```

### Ejercicio 5: Top K Elementos usando Heap ⭐⭐⭐

```go
/*
Problema:
Encontrar los K elementos más frecuentes en un stream de datos.
Este es un problema clásico en entrevistas técnicas.

Requisitos:
- Contar frecuencia de elementos
- Mantener heap de top K
- Eficiencia: O(n log k) tiempo
- Soportar valores negativos, strings, cualquier comparable

Bonus: Implementar para stream (datos que llegan continuamente)

Ejemplo:
TopK([1,1,1,2,2,3], 2) → [(1, frecuencia=3), (2, frecuencia=2)]
TopK(stream, 5) → Top 5 palabras más frecuentes en texto
*/

package main

import (
    "container/heap"
    "fmt"
    "sort"
)

type PareFreq struct {
    Elemento interface{}
    Frecuencia int
}

func TopKElementos(datos []interface{}, k int) []PareFreq {
    // TODO: Contar frecuencias
    // TODO: Usar heap para mantener top K
    // TODO: Retornar resultados ordenados
    return nil
}

// Bonus: versión para stream
type StreamTopK struct {
    K int
    // TODO: Agregar campos
}

// TODO: Implementar función main() con demo
```

---

## Resumen

Este capítulo ha explorado todo lo que necesitas saber sobre ordenamiento y estructuras de datos avanzadas en Go:

- **Sort**: Desde interfaces básicas hasta multi-criterio
- **Search**: Búsqueda binaria para máxima eficiencia
- **Heap**: Colas de prioridad y aplicaciones en tiempo real
- **Container**: Ring para circularidad, List para flexibilidad
- **Performance**: Algoritmos internos y benchmarking

Go provee herramientas simples pero poderosas. La clave es entender cuándo usar cada una y cómo optimizarlas para tu caso de uso específico.

### Puntos Clave

1. **sort.Interface** es el contrato fundamental (Len, Less, Swap)
2. **sort.Slice** es más práctico para la mayoría de casos
3. **Búsqueda binaria** es O(log n) vs O(n) lineal - usar siempre en datos ordenados
4. **Heap** es ideal para colas de prioridad y selección de top K
5. **List** es mejor que Array para inserciones/eliminaciones frecuentes
6. **Ring** para casos específicos como buffers circulares

¡Practica con los ejercicios para dominar estos conceptos!

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/34-sort-y-collections/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/34-sort-y-collections):

```bash
cd examples/34-sort-y-collections
go run .
```
