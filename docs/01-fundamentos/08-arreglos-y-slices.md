# Capítulo 8: Arreglos y slices

## Índice del Capítulo 8

1. [8.1 ¿Qué son Arreglos y Slices?](#81-qué-son-arreglos-y-slices)
2. [8.2 Arreglos (Arrays) - Tamaño Fijo](#82-arreglos-arrays---tamaño-fijo)
3. [8.3 Slices - Dinámicos y Flexibles](#83-slices---dinámicos-y-flexibles)
4. [8.4 Declaración y Inicialización](#84-declaración-e-inicialización)
5. [8.5 Indexación y Slicing](#85-indexación-y-slicing)
6. [8.6 Len, Cap y make()](#86-len-cap-y-make)
7. [8.7 Append y Copy](#87-append-y-copy)
8. [8.8 Iteración sobre Colecciones](#88-iteración-sobre-colecciones)
9. [8.9 Slices Multidimensionales](#89-slices-multidimensionales)
10. [8.10 Buenas Prácticas](#810-buenas-prácticas)

---

## 8.1 ¿Qué son Arreglos y Slices?

### Definiciones

**Arreglo (Array):**

- Colección de elementos del **MISMO TIPO**
- Tamaño **FIJO**, conocido en compile time
- Elemento accesible por **ÍNDICE**

**Slice (Porción):**

- Colección de elementos del **MISMO TIPO**
- Tamaño **DINÁMICO**, puede crecer
- Elemento accesible por **ÍNDICE**
- Es una "vista" sobre un arreglo subyacente

### Comparación Rápida

```

 Aspecto     │ Array        │ Slice        │
clear
 Tamaño      │ FIJO         │ DINÁMICO     │
 Declaración │ [5]int       │ []int        │
 Crecimiento │ NO           │ S (append)  │
 Copia       │ Copia valor  │ Copia ref    │
 Uso         │ Raro         │ COMÚN        │
 Performance │ Más rápido   │ Ligeramente  │


REGLA DE ORO: Usa Slices 99% del tiempo
```

### Abstracción: Slice sobre Array

```
Array subyacente:
[10, 20, 30, 40, 50]

Slice puede apuntar a parte del array:
slice := array[1:4]    // [20, 30, 40]

El slice es una "vista" o "ventana" al array original.
Si modificas el slice, modificas el array.
```

---

## 8.2 Arreglos (Arrays) - Tamaño Fijo

### Declaración Básica

```go
var numeros [5]int      // Array de 5 enteros, inicializado a 0
var nombres [3]string   // Array de 3 strings, inicializado a ""

fmt.Println(numeros)    // [0 0 0 0 0]
fmt.Println(nombres)    // [  ]
```

### Inicialización Explícita

```go
// Con valores específicos
var colores [3]string = [3]string{"rojo", "verde", "azul"}

// Sintaxis corta
colores := [3]string{"rojo", "verde", "azul"}

// Sin especificar cantidad (infiere del número de elementos)
numeros := [...]int{10, 20, 30, 40}    // [4]int automáticamente

fmt.Println(len(colores))  // 3
fmt.Println(len(numeros))  // 4
```

### Acceso a Elementos

```go
numeros := [5]int{10, 20, 30, 40, 50}

fmt.Println(numeros[0])     // 10 (primer elemento)
fmt.Println(numeros[4])     // 50 (último elemento)
fmt.Println(numeros[-1])    // ❌ ERROR: índice negativo

numeros[2] = 99             // Modificar elemento
fmt.Println(numeros)        // [10 20 99 40 50]
```

### Iteración

```go
numeros := [5]int{10, 20, 30, 40, 50}

// Por índice
for i := 0; i < len(numeros); i++ {
    fmt.Println(i, numeros[i])
}

// Por range (retorna índice y valor)
for i, v := range numeros {
    fmt.Printf("Índice: %d, Valor: %d\n", i, v)
}
```

### Copia de Arreglos

```go
a := [3]int{1, 2, 3}
b := a                  // Copia el CONTENIDO completo
b[0] = 100

fmt.Println(a)          // [1 2 3] (sin cambios)
fmt.Println(b)          // [100 2 3]
```

### Arrays Multidimensionales

```go
// Matriz 2x3 (2 filas, 3 columnas)
matriz := [2][3]int{
    {1, 2, 3},
    {4, 5, 6},
}

fmt.Println(matriz[0][0])   // 1
fmt.Println(matriz[1][2])   // 6
```

### Limitaciones de Arrays

```go
// ❌ NO PUEDES cambiar tamaño
arr := [5]int{1, 2, 3, 4, 5}
arr = append(arr, 6)        // ❌ ERROR: no puedes append a array

// ❌ NO PUEDES comparar tamaños diferentes
a := [3]int{1, 2, 3}
b := [5]int{1, 2, 3, 4, 5}
if a == b { }               // ❌ ERROR: tipos diferentes [3]int vs [5]int

// ✅ PUEDES comparar arrays del mismo tamaño
c := [3]int{1, 2, 3}
d := [3]int{1, 2, 3}
if c == d {
    fmt.Println("Iguales")  // ✅ OK
}
```

---

## 8.3 Slices - Dinámicos y Flexibles

### ¿Por Qué Slices?

```
Arrays: Tamaño fijo, no puedes crecer

Ejemplo: Leer números de entrada
 ¿Cuántos números? No sabes
 No puedes hacer [?]int
 Necesitas Slice para tamaño desconocido
```

### Estructura Interna de un Slice

```go
// Un slice tiene 3 componentes internos:
type SliceHeader struct {
    Data uintptr    // Puntero al array subyacente
    Len  int        // Longitud actual
    Cap  int        // Capacidad (tamaño del array subyacente)
}

// Ejemplo visual:
slice := []int{10, 20, 30, 40, 50}

Slice Header:
 Data: apunta a [10, 20, 30, 40, 50]
 Len: 5 (tengo 5 elementos)
 Cap: 5 (el array tiene capacidad para 5)
```

### Declaración Básica

```go
// Slice vacío
var vacio []int             // nil slice

// Slice con valores
numeros := []int{10, 20, 30}

// Slice del array
arr := [5]int{10, 20, 30, 40, 50}
slice := arr[1:4]           // [20, 30, 40]

fmt.Println(vacio)          // []
fmt.Println(numeros)        // [10 20 30]
fmt.Println(slice)          // [20 30 40]
```

### Slice nil vs Slice Vacío

```go
var nilSlice []int          // nil slice (sin array subyacente)
var vacio []int = []int{}   // Slice vacío (array subyacente vacío)

fmt.Println(nilSlice == nil)        // true
fmt.Println(vacio == nil)           // false ← ¡Diferencia importante!

// Ambos tienen len=0 y cap=0
fmt.Println(len(nilSlice), cap(nilSlice))  // 0 0
fmt.Println(len(vacio), cap(vacio))        // 0 0

// Ambos funcionan igual en la mayoría de casos
for range nilSlice { }  // ✅ No hace nada
for range vacio { }     // ✅ No hace nada
```

### Modificar Slice Modifica Array Subyacente

```go
arr := [5]int{10, 20, 30, 40, 50}
slice := arr[1:4]           // [20, 30, 40]

slice[0] = 200              // Modifica slice

fmt.Println(slice)          // [200 30 40]
fmt.Println(arr)            // [10 200 30 40 50] ← ¡Array también cambió!
```

---

## 8.4 Declaración e Inicialización

### Formas de Crear un Slice

**Forma 1: Literal**

```go
numeros := []int{10, 20, 30, 40}
```

**Forma 2: make() - Tamaño y Capacidad**

```go
// make(tipo, len, cap)
slice := make([]int, 3)         // len=3, cap=3
fmt.Println(slice)              // [0 0 0]

// Con capacidad diferente a longitud
slice2 := make([]int, 3, 10)    // len=3, cap=10
fmt.Println(len(slice2))        // 3
fmt.Println(cap(slice2))        // 10
```

**Forma 3: Desde Array**

```go
arr := [5]int{10, 20, 30, 40, 50}
slice := arr[:]                 // Todos los elementos
slice2 := arr[1:4]              // Elementos 1, 2, 3
slice3 := arr[:3]               // Elementos 0, 1, 2
slice4 := arr[2:]               // Elementos 2 hasta el final
```

**Forma 4: Nil Slice**

```go
var slice []int                 // nil
```

### Comparación: make vs Literal

```go
// Literal: conoces valores
slice1 := []int{1, 2, 3}       // len=3, cap=3

// make: conoces tamaño, pero no valores
slice2 := make([]int, 3)       // len=3, cap=3, valores=0

// make: capacidad extra para crecer sin reasignación
slice3 := make([]int, 0, 10)   // len=0, cap=10 (vacío pero con espacio)
```

---

## 8.5 Indexacinnn y Slicing

### Indexación Simple

```go
numeros := []int{10, 20, 30, 40, 50}

fmt.Println(numeros[0])     // 10
fmt.Println(numeros[2])     // 30
fmt.Println(numeros[-1])    // ❌ ERROR: índice negativo

numeros[1] = 200
fmt.Println(numeros)        // [10 200 30 40 50]
```

### Operación de Slicing

```go
numeros := []int{10, 20, 30, 40, 50}

slice := numeros[start:end]
// Incluye start, EXCLUYE end
// Retorna nuevo slice

numeros[1:4]    // [20, 30, 40] (indices 1, 2, 3)
numeros[:3]     // [10, 20, 30] (desde inicio hasta 3)
numeros[2:]     // [30, 40, 50] (desde 2 hasta fin)
numeros[:]      // [10, 20, 30, 40, 50] (copia del slice completo)
```

### Slicing con Tres Parámetros (Capacidad)

```go
numeros := []int{10, 20, 30, 40, 50}

slice := numeros[1:4:4]
// [start:end:cap]
// len = end - start = 3
// cap = cap - start = 3

// Ejemplo más claro:
arr := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
slice := arr[2:5:7]
// len = 5 - 2 = 3
// cap = 7 - 2 = 5

fmt.Println(slice)          // [2 3 4]
fmt.Println(len(slice))     // 3
fmt.Println(cap(slice))     // 5
```

### Bounds Checking

```go
numeros := []int{10, 20, 30, 40, 50}

fmt.Println(numeros[0])     // 10 (válido)
fmt.Println(numeros[4])     // 50 (válido)
fmt.Println(numeros[5])     // ❌ PANIC: index out of range

numeros[1:6]                // ❌ PANIC: slice bounds out of range
```

---

## 8.6 Len, Cap y make()

### len() - Longitud

```go
numeros := []int{10, 20, 30}
fmt.Println(len(numeros))           // 3

vacio := []int{}
fmt.Println(len(vacio))             // 0

nilSlice := ([]int)(nil)
fmt.Println(len(nilSlice))          // 0 (no panic)
```

### cap() - Capacidad

```go
// make(tipo, len, cap)
slice := make([]int, 3, 10)

fmt.Println(len(slice))             // 3 (elementos actuales)
fmt.Println(cap(slice))             // 10 (espacio total disponible)

// Cuando haces append dentro de cap, no reasigna
slice = append(slice, 40)           // Len=4, cap=10 (sin reasignar)
```

### Reasignación de Slice

```go
slice := make([]int, 3, 3)  // len=3, cap=3

// Al append, se necesita más espacio
slice = append(slice, 40)   // len=4, cap=? (reasignación)

// Go crece típicamente: new_cap = old_cap * 2 (o más)
fmt.Println(len(slice))     // 4
fmt.Println(cap(slice))     // 6 o 8 (depende de implementación)
```

### Visualización de len vs cap

```
slice := make([]int, 3, 10)

len=3, cap=10:

 [X][X][X][ ][ ][ ][ ][ ][ ][ ] │
  ↑len=3→ ↑cap=10──────────────→│


Después de append(slice, 40):
len=4, cap=10:

 [X][X][X][X][ ][ ][ ][ ][ ][ ] │
  ↑len=4→ ↑cap=10──────────────→│


Después de append 7 más (total 11 elementos):
len=11, cap=? (reasignación):

 [X][X][X][X][X][X][X][X][X][X][X][ ]... │
  ↑len=11──────→ ↑cap=16 (o más)────────→│

```

---

## 8.7 Append y Copy

### Append - Añadir Elementos

```go
numeros := []int{10, 20, 30}

// Añadir un elemento
numeros = append(numeros, 40)
fmt.Println(numeros)        // [10 20 30 40]

// Añadir múltiples elementos
numeros = append(numeros, 50, 60, 70)
fmt.Println(numeros)        // [10 20 30 40 50 60 70]

// Append otro slice
otro := []int{80, 90}
numeros = append(numeros, otro...)  // ← El ... desempaca el slice
fmt.Println(numeros)        // [10 20 30 40 50 60 70 80 90]
```

### IMPORTANTE: Reasignación de Append

```go
slice := []int{1, 2, 3}
result := append(slice, 4)  // Puede o no ser el mismo slice

// REGLA: Siempre reasigna el resultado
slice = append(slice, 4)    // ✅ Correcto

// ❌ Riesgo de bug
append(slice, 4)            // Ignoras el resultado
fmt.Println(slice)          // Puede no estar actualizado
```

### Copy - Copiar Slice

```go
original := []int{1, 2, 3, 4, 5}

// Copiar todo
copia := make([]int, len(original))
copy(copia, original)       // copy(dst, src)
fmt.Println(copia)          // [1 2 3 4 5]

// Modificar copia no afecta original
copia[0] = 100
fmt.Println(original)       // [1 2 3 4 5] (sin cambios)
fmt.Println(copia)          // [100 2 3 4 5]
```

### Copy Parcial

```go
original := []int{1, 2, 3, 4, 5}

// Copiar solo primeros 3 elementos
copia := make([]int, 3)
copy(copia, original)
fmt.Println(copia)          // [1 2 3]

// Copy retorna número de elementos copiados
original2 := []int{10, 20}
copia2 := make([]int, 5)
n := copy(copia2, original2)    // n = 2
fmt.Println(copia2)             // [10 20 0 0 0]
fmt.Println(n)                  // 2
```

---

## 8.8 Iteración sobre Colecciones

### For con Índice

```go
numeros := []int{10, 20, 30}

for i := 0; i < len(numeros); i++ {
    fmt.Println(numeros[i])
}
```

### For Range - Índice y Valor

```go
numeros := []int{10, 20, 30}

for i, v := range numeros {
    fmt.Printf("Índice: %d, Valor: %d\n", i, v)
}

// Output:
// Índice: 0, Valor: 10
// Índice: 1, Valor: 20
// Índice: 2, Valor: 30
```

### For Range - Solo Valores

```go
numeros := []int{10, 20, 30}

for _, v := range numeros {
    fmt.Println(v)
}

// Output:
// 10
// 20
// 30
```

### For Range - Solo Índices

```go
numeros := []int{10, 20, 30}

for i := range numeros {
    fmt.Println(i)
}

// Output:
// 0
// 1
// 2
```

### Modificar durante Iteración

```go
numeros := []int{1, 2, 3, 4, 5}

// Modificar elemento (funciona)
for i, v := range numeros {
    numeros[i] = v * 2
}
fmt.Println(numeros)        // [2 4 6 8 10]

// ⚠️ Cuidado: append durante range
slice := []int{1, 2, 3}
for i, v := range slice {
    if i == 1 {
        slice = append(slice, 100)  // Agranda slice
    }
    fmt.Println(v)
}
// Output: 1, 2, 3 (el 100 NO se itera, range se ejecutó antes)
```

---

## 8.9 Slices Multidimensionales

### Slice 2D (Matriz)

```go
// Slice de slices
matriz := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

fmt.Println(matriz[0])      // [1 2 3]
fmt.Println(matriz[0][0])   // 1
fmt.Println(matriz[1][2])   // 6
```

### Crear Matriz 2D Dinámicamente

```go
filas := 3
columnas := 4

// Crear matriz 3x4
matriz := make([][]int, filas)
for i := 0; i < filas; i++ {
    matriz[i] = make([]int, columnas)
}

// Llenar
for i := 0; i < filas; i++ {
    for j := 0; j < columnas; j++ {
        matriz[i][j] = i*columnas + j
    }
}

fmt.Println(matriz)
// [[0 1 2 3] [4 5 6 7] [8 9 10 11]]
```

### Slice 3D (Tensor)

```go
// Slice de slice de slices
tensor := [][][]int{
    {
        {1, 2},
        {3, 4},
    },
    {
        {5, 6},
        {7, 8},
    },
}

fmt.Println(tensor[0][0][0])   // 1
fmt.Println(tensor[1][1][1])   // 8
```

### Jagged Arrays (Tamaños Diferentes)

```go
// En Go puedes tener filas de diferente tamaño
jagged := [][]int{
    {1, 2, 3},
    {4, 5},
    {6, 7, 8, 9},
}

fmt.Println(jagged[0])         // [1 2 3]
fmt.Println(jagged[1])         // [4 5]
fmt.Println(jagged[2])         // [6 7 8 9]
```

---

## 8.10 Buenas Prácticas

### Preferir Slices sobre Arrays

```go
// ❌ Raro usar arrays
func procesar(arr [5]int) {
}

// ✅ Usa slices
func procesar(slice []int) {
}
```

### Pre-asignar Slices si Conoces el Tamaño

```go
// ❌ Ineficiente (múltiples reasignaciones)
var resultado []int
for i := 0; i < 1000; i++ {
    resultado = append(resultado, i)
}

// ✅ Eficiente (pre-asignar)
resultado := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    resultado = append(resultado, i)
}
```

### Siempre Reasignar Append

```go
// ❌ Riesgo de bug
slice := []int{1, 2, 3}
append(slice, 4)            // Resultado ignorado
// slice podría no estar actualizado

// ✅ Correcto
slice := []int{1, 2, 3}
slice = append(slice, 4)    // Reasigna el resultado
```

### Usar Copy para Copias Reales

```go
// ❌ No copia, solo referencia
original := []int{1, 2, 3}
copia := original
copia[0] = 100
fmt.Println(original)       // [100 2 3] ← Modificó original!

// ✅ Copia real
original := []int{1, 2, 3}
copia := make([]int, len(original))
copy(copia, original)
copia[0] = 100
fmt.Println(original)       // [1 2 3] (sin cambios)
```

### Pasar Slices a Funciones

```go
// Slices son referencias, se pasan eficientemente
func procesar(datos []int) {
    // Cambios aquí afectan el original
    datos[0] = 100
}

numeros := []int{1, 2, 3}
procesar(numeros)
fmt.Println(numeros)        // [100 2 3]
```

### Ranges Seguros

```go
// ✅ Range es seguro con slices vacíos
vacio := []int{}
for range vacio {
    // No se ejecuta
}

// ✅ Range es seguro con nil slice
var nil_slice []int
for range nil_slice {
    // No se ejecuta
}
```

### Evitar Índices Hardcodeados

```go
// ❌ Mágico
elemento := datos[5]

// ✅ Descriptivo
const (
    PosicionImportante = 5
)
elemento := datos[PosicionImportante]

// O mejor, usar búsqueda
idx := findIndex(datos, valor)
if idx >= 0 {
    elemento := datos[idx]
}
```

### Slice vs Array - Cuándo Usar Cada Una

```
ARRAY:
 Tamaño FIJO y conocido
 Ejemplos: [12]Mes, [7]DiaSemana, [3]RGB
 Raro en práctica

SLICE:
 Tamaño VARIABLE
 Ejemplos: lista de usuarios, resultados de búsqueda
 99% de los casos
```

---

## Ejercicios del Capítulo 8

### Ejercicio 1: Manipulador de Números

Crea programa que:

1. Cree slice de 10 números aleatorios
2. Encuentre máximo y mínimo
3. Calcule promedio
4. Invierta el orden
5. Filtre números pares
6. Muestra cada resultado

### Ejercicio 2: Gestor de Tareas

Crea programa que:

1. Mantenga slice de tareas (strings)
2. Permitir agregar, remover, listar
3. Buscar tarea por nombre
4. Marcar tarea como completada (usando índice)
5. Pre-asignar slice si esperas 100 tareas

### Ejercicio 3: Matriz de Puntuaciones

Crea programa que:

1. Cree matriz 3x3 de puntuaciones
2. Calcule suma total
3. Calcule suma por fila
4. Calcule suma por columna
5. Encuentre puntuación máxima y su posición
6. Muestra matriz formateada

### Ejercicio 4: Transformador de Slices

Crea programa que:

1. Cree slice de números
2. Duplique cada elemento: [1,2,3] → [1,1,2,2,3,3]
3. Intercale dos slices: [1,3,5] y [2,4,6] → [1,2,3,4,5,6]
4. Rotación: [1,2,3,4,5] rotado 2 → [4,5,1,2,3]
5. Elimine duplicados: [1,1,2,3,3,3] → [1,2,3]

### Ejercicio 5: Comparador y Analizador

Crea programa que:

1. Cree dos slices de números
2. Determine si son iguales (mismo contenido)
3. Encuentre elementos en común
4. Encuentre elementos únicos de cada slice
5. Fusiónelos y ordene
6. Compare performance: append vs pre-asignación

---

**Fin del Capítulo 8**

---

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/08-arreglos-y-slices/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/08-arreglos-y-slices):

```bash
cd examples/08-arreglos-y-slices
go run .
```
