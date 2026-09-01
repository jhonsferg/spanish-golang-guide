# Capítulo 9: Maps - Colecciones clave-valor

## Índice del Capítulo 9

1. [9.1 ¿Qué es un Map?](#91-qué-es-un-map)
2. [9.2 Declaración e Inicialización](#92-declaración-e-inicialización)
3. [9.3 Acceso a Elementos](#93-acceso-a-elementos)
4. [9.4 Modificación y Eliminación](#94-modificación-y-eliminación)
5. [9.5 La Prueba de Existencia (Comma OK)](#95-la-prueba-de-existencia-comma-ok)
6. [9.6 Iteración sobre Maps](#96-iteración-sobre-maps)
7. [9.7 Maps Anidados](#97-maps-anidados)
8. [9.8 Maps como Argumentos](#98-maps-como-argumentos)
9. [9.9 Maps vs Structs](#99-maps-vs-structs)
10. [9.10 Buenas Prácticas](#910-buenas-prácticas)

---

## 9.1 ¿Qué es un Map?

### Definición

Un **map** es una colección de pares **clave-valor**:

```
Map = Diccionario = Tabla Hash = Asociativo Array

Ejemplo real:
 Teléfono = Mapa de (Nombre → Número)
   "Juan" → "555-1234"
   "María" → "555-5678"
   "Carlos" → "555-9012"

 Inventario = Mapa de (Producto → Cantidad)
   "Manzanas" → 50
   "Naranjas" → 30
   "Plátanos" → 25

 Cache = Mapa de (URL → Respuesta)
    "http://api.com/data" → "JSON response"
```

### Características Clave

```

 Aspecto            │ Array/Slice    │ Map            │
clear
 Índice             │ Posición (0,1) │ Clave          │
 Acceso             │ O(1)           │ O(1)           │
 Orden              │ Sí             │ NO (aleatorio) │
 Duplicados         │ Permitido      │ No (claves)    │
 Tamaño             │ Dinámico       │ Dinámico       │
 Zero value         │ []              │ nil            │

```

### Por Qué Usar Maps

```
SIN Maps:
    nombreAño1 := "Juan"
    edadAño1 := 25
    nombreAño2 := "María"
    edadAño2 := 30

CON Maps:
    personas := map[string]int{
        "Juan": 25,
        "María": 30,
    }

Beneficios:
 Búsqueda rápida: personas["Juan"] → 25
 Dinámico: agregar nuevas personas fácilmente
 Legible: estructura clara de datos
 Escalable: cientos/miles de entradas
```

### Cómo Funciona Internamente

```
Go usa TABLA HASH (Hash Table) internamente:

Clave "Juan":
 Hash("Juan") = 12345
 Bucket = 12345 % num_buckets = 5
 Valor: 25 está en bucket 5

Ventaja:
 Acceso O(1) promedio (muy rápido)


 Hash Map Internamente:                      │
                                             │
 Bucket 0: []                                │
 Bucket 1: [("María", 30)]                   │
 Bucket 2: []                                │
 Bucket 3: [("Carlos", 35)]                  │
 Bucket 4: []                                │
 Bucket 5: [("Juan", 25)]                    │
 ...                                         │

```

---

## 9.2 Declaración e Inicialización

### Declaración Básica

```go
// Map vacío
var edades map[string]int        // nil map
fmt.Println(edades)              // map[]
fmt.Println(edades == nil)       // true
```

### Inicialización con Literals

```go
// Map con valores iniciales
edades := map[string]int{
    "Juan": 25,
    "María": 30,
    "Carlos": 35,
}

fmt.Println(edades)              // map[Juan:25 María:30 Carlos:35]
fmt.Println(edades["Juan"])      // 25
```

### Inicialización con make()

```go
// Map vacío con make
edades := make(map[string]int)
fmt.Println(edades)              // map[]
fmt.Println(len(edades))         // 0

// Map con capacidad estimada (optimización)
edades := make(map[string]int, 100)  // Espera ~100 entradas
```

### Tipos de Clave y Valor

```
TIPOS DE CLAVE permitidos:
 Cualquier tipo que sea COMPARABLE
 int, float64, string, bool
 Arrays (arrays tienen == definido)
 Structs (si todos los campos son comparables)
 Punteros
 ❌ NO: Slices, Maps, Canales (no comparables)

TIPOS DE VALOR:
 Cualquier tipo: int, string, interface{}, slices, structs, etc.

Ejemplos:
    map[string]string       // Teléfono: nombre → número
    map[int][]int          // Factores: número → lista de factores
    map[string]User        // Usuarios: ID → User struct
    map[string]interface{} // JSON: clave → cualquier valor
```

### Ejemplo: Map[string]interface{}

Almacena valores de cualquier tipo:

```go
datos := map[string]interface{}{
    "nombre": "Juan",
    "edad": 25,
    "activo": true,
    "salario": 50000.50,
}

fmt.Println(datos["nombre"])    // Juan (string)
fmt.Println(datos["edad"])      // 25 (int)
fmt.Println(datos["activo"])    // true (bool)
```

---

## 9.3 Acceso a Elementos

### Lectura Simple

```go
edades := map[string]int{
    "Juan": 25,
    "María": 30,
}

fmt.Println(edades["Juan"])     // 25
fmt.Println(edades["María"])    // 30
```

### Acceso a Clave Inexistente

```go
edades := map[string]int{
    "Juan": 25,
}

fmt.Println(edades["Carlos"])   // 0 (cero valor, sin panic)
```

**Importante:**
- NO hay panic si la clave no existe
- Retorna cero valor del tipo de valor
- Puede ser confuso: ¿existe clave con valor 0, o no existe?

---

## 9.4 Modificación y Eliminación

### Agregar/Modificar

```go
edades := map[string]int{
    "Juan": 25,
}

// Agregar nueva clave
edades["María"] = 30

// Modificar existente
edades["Juan"] = 26

fmt.Println(edades)
// map[Juan:26 María:30]
```

### Eliminar Elementos

```go
edades := map[string]int{
    "Juan": 25,
    "María": 30,
    "Carlos": 35,
}

// Eliminar clave
delete(edades, "María")

fmt.Println(edades)
// map[Juan:25 Carlos:35]

// Eliminar clave inexistente (seguro, no panic)
delete(edades, "Pedro")
```

### Limpiar Map (Vaciar todo)

```go
edades := map[string]int{
    "Juan": 25,
    "María": 30,
}

// NO existe clear(edades) en Go 1.22-
// Solución: crear nuevo map

edades = make(map[string]int)   // Vacío

// O mediante bucle (si necesitas Go < 1.22)
for k := range edades {
    delete(edades, k)
}
```

---

## 9.5 La Prueba de Existencia (Comma OK)

### ¿Cómo verificar si existe una clave?

Problema: Acceso a clave inexistente retorna cero valor (confuso)

```go
edades := map[string]int{
    "Juan": 25,
}

valor := edades["Carlos"]   // 0
// ¿Carlos no existe, o Carlos tiene 0 años?
```

### Solución: Comma OK

Go permite retornar segundo valor (existencia):

```go
edades := map[string]int{
    "Juan": 25,
}

// valor, existe := map[clave]
valor, existe := edades["Juan"]
if existe {
    fmt.Printf("Juan tiene %d años\n", valor)
} else {
    fmt.Println("Juan no existe")
}

// Output: Juan tiene 25 años
```

### Con Clave Inexistente

```go
edades := map[string]int{
    "Juan": 25,
}

valor, existe := edages["Carlos"]
if existe {
    fmt.Printf("Carlos existe: %d\n", valor)
} else {
    fmt.Println("Carlos no existe")
}

// Output: Carlos no existe

// Nota: valor es 0 (cero valor)
fmt.Println(valor)          // 0
```

### Idioma Go: Ignorar Valor

```go
_, existe := edades["Juan"]

if existe {
    fmt.Println("Juan existe")
}

// O más conciso:
if _, ok := edades["Juan"]; ok {
    fmt.Println("Juan existe")
}
```

---

## 9.6 Iteración sobre Maps

### Iteración Simple (clave y valor)

```go
edades := map[string]int{
    "Juan": 25,
    "María": 30,
    "Carlos": 35,
}

for nombre, edad := range edades {
    fmt.Printf("%s tiene %d años\n", nombre, edad)
}

// Output (orden ALEATORIO):
// María tiene 30 años
// Carlos tiene 35 años
// Juan tiene 25 años
```

**IMPORTANTE: Orden aleatorio**

Go deliberadamente itera en orden aleatorio en maps:

```go
// Si corres dos veces, el orden es diferente
for nombre := range edades {
    fmt.Println(nombre)
}

// Ejecución 1: María, Juan, Carlos
// Ejecución 2: Carlos, Juan, María
// Ejecución 3: Juan, Carlos, María
```

**¿Por qué?** Para evitar que programadores dependan del orden (que es no determinístico en hash maps).

### Solo Claves

```go
edades := map[string]int{
    "Juan": 25,
    "María": 30,
}

for nombre := range edades {
    fmt.Println(nombre)
}

// Output:
// María
// Juan (u otro orden)
```

### Solo Valores

```go
for _, edad := range edades {
    fmt.Println(edad)
}

// Output:
// 30
// 25 (u otro orden)
```

### Ordenar Claves antes de Iterar

Si necesitas orden específico:

```go
import "sort"

edades := map[string]int{
    "Juan": 25,
    "María": 30,
    "Carlos": 35,
}

// Extraer claves
claves := make([]string, 0, len(edades))
for k := range edades {
    claves = append(claves, k)
}

// Ordenar
sort.Strings(claves)

// Iterar en orden
for _, nombre := range claves {
    fmt.Printf("%s: %d\n", nombre, edades[nombre])
}

// Output:
// Carlos: 35
// Juan: 25
// María: 30
```

---

## 9.7 Maps Anidados

### Map Dentro de Map

```go
// Map de maps
ciudades := map[string]map[string]int{
    "España": map[string]int{
        "Madrid": 3200000,
        "Barcelona": 1600000,
    },
    "México": map[string]int{
        "Mexico City": 9000000,
        "Guadalajara": 5200000,
    },
}

// Acceso
fmt.Println(ciudades["España"]["Madrid"])    // 3200000
```

### Crear Maps Anidados Dinámicamente

```go
// Inicializa como nil
ciudades := make(map[string]map[string]int)

// ❌ ERROR: no puedes asignar a mapa anidado inexistente
ciudades["España"]["Madrid"] = 3200000    // PANIC!

// ✅ Correcto: crear mapa interior primero
ciudades["España"] = make(map[string]int)
ciudades["España"]["Madrid"] = 3200000

// Alternativa: Comprobar y crear
if ciudades["España"] == nil {
    ciudades["España"] = make(map[string]int)
}
ciudades["España"]["Madrid"] = 3200000
```

### Map Anidado Profundo

```go
// Map → Map → Map → Valor
config := map[string]map[string]map[string]string{
    "database": map[string]map[string]string{
        "mysql": map[string]string{
            "host": "localhost",
            "port": "3306",
            "user": "admin",
        },
    },
}

fmt.Println(config["database"]["mysql"]["host"])  // localhost
```

---

## 9.8 Maps como Argumentos

### Pasar Map a Función

```go
func mostrarEdades(edades map[string]int) {
    for nombre, edad := range edades {
        fmt.Printf("%s: %d\n", nombre, edad)
    }
}

edades := map[string]int{"Juan": 25, "María": 30}
mostrarEdades(edades)
```

### Maps son Referencias

Cambios en función afectan el map original:

```go
func incrementarEdades(edades map[string]int) {
    for k := range edades {
        edades[k]++  // Modifica original
    }
}

edades := map[string]int{"Juan": 25, "María": 30}
incrementarEdades(edades)
fmt.Println(edades)  // map[Juan:26 María:31]
```

### Retornar Map desde Función

```go
func crearEdades() map[string]int {
    return map[string]int{
        "Juan": 25,
        "María": 30,
    }
}

edades := crearEdades()
fmt.Println(edades)  // map[Juan:25 María:30]
```

---

## 9.9 Maps vs Structs

### Cuándo Usar Map

**Usa Map cuando:**
```
 Estructura de datos es dinámica/flexible
 Las claves son desconocidas en advance
 Puede haber muchas variaciones de claves
 Ejemplo: JSON deserialized, configuración dinámica
```

**Ejemplo de Map:**

```go
// Respuesta HTTP dinámica (claves desconocidas)
respuesta := map[string]interface{}{
    "status": 200,
    "data": []int{1, 2, 3},
    "error": nil,
}
```

### Cuándo Usar Struct

**Usa Struct cuando:**
```
 Estructura de datos es FIJA/conocida
 Campos específicos con tipos específicos
 Mejor para type safety
 Mejor para documentacinnn
```

**Ejemplo de Struct:**

```go
// Definición clara y segura
type Usuario struct {
    Nombre string
    Edad   int
    Email  string
    Activo bool
}

usuario := Usuario{
    Nombre: "Juan",
    Edad:   25,
    Email:  "juan@example.com",
    Activo: true,
}
```

### Comparación Directa

```go
// MAP: Flexible pero sin type safety
usuario := map[string]interface{}{
    "nombre": "Juan",
    "edad": 25,
}
// ❌ No hay validación en compile time
// ❌ Campo con typo no es detectado
usuario["nmbre"] = "Otro"  // Typo, pero válido

// STRUCT: Type-safe pero rígido
type Usuario struct {
    Nombre string
    Edad   int
}
// ✅ Solo estos campos específicos
// ✅ Typo es detectado en compile time
usuario.Nombre = "Juan"    // ✅ Correcto
usuario.nmbre = "Juan"     // ❌ ERROR de compilación
```

### Hybrid: Struct con Map

Combina lo mejor de ambos:

```go
type Usuario struct {
    Nombre string
    Edad   int
    Metadata map[string]interface{}  // Flexible
}

usuario := Usuario{
    Nombre: "Juan",
    Edad: 25,
    Metadata: map[string]interface{}{
        "ciudad": "Madrid",
        "intereses": []string{"Go", "Docker"},
        "premium": true,
    },
}

fmt.Println(usuario.Nombre)              // Juan (type-safe)
fmt.Println(usuario.Metadata["ciudad"])  // Madrid (flexible)
```

---

## 9.10 Buenas Prácticas

### Siempre Comprobar Existencia

```go
// ❌ Riesgo: ¿existe o valor es 0?
edad := edades["Juan"]
if edad > 18 {
    fmt.Println("Adulto")
}

// ✅ Seguro: comprobar existencia
if edad, ok := edades["Juan"]; ok && edad > 18 {
    fmt.Println("Adulto")
} else {
    fmt.Println("Menor o inexistente")
}
```

### Usar Valores Zero Apropiadamente

```go
// Si una clave no existe, obtiene cero valor
contador := map[string]int{}

contador["Juan"]++       // Juan: 0 + 1 = 1
contador["María"]++      // María: 0 + 1 = 1

fmt.Println(contador)    // map[Juan:1 María:1]
```

### Pre-asignar Capacity si Conoces Tamaño

```go
// ❌ Ineficiente
datos := make(map[string]int)
for i := 0; i < 1000; i++ {
    datos[fmt.Sprintf("key%d", i)] = i
}

// ✅ Eficiente
datos := make(map[string]int, 1000)
for i := 0; i < 1000; i++ {
    datos[fmt.Sprintf("key%d", i)] = i
}
```

### No Modificar durante Iteración

```go
// ❌ Peligroso (comportamiento indefinido)
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k := range m {
    delete(m, k)        // Modificar durante iteración
}

// ✅ Seguro: colecta claves primero
claves := make([]string, 0, len(m))
for k := range m {
    claves = append(claves, k)
}
for _, k := range claves {
    delete(m, k)
}
```

### Usar Nombres Significativos

```go
// ❌ Confuso
users := make(map[string]map[string]interface{})

// ✅ Claro
type User struct {
    Nombre string
    Edad   int
}
usersByID := make(map[int]User)

// O si es dinámico
userProfiles := make(map[string]map[string]interface{})
```

### Maps para Cache

```go
// Pattern común: Map como cache
type Cache struct {
    data map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    val, ok := c.data[key]
    return val, ok
}

func (c *Cache) Set(key string, val interface{}) {
    c.data[key] = val
}

func (c *Cache) Delete(key string) {
    delete(c.data, key)
}
```

### Maps para Deduplicación

```go
// Encontrar elementos únicos
numeros := []int{1, 2, 2, 3, 3, 3, 4}

seen := make(map[int]bool)
var unicos []int

for _, n := range numeros {
    if !seen[n] {
        seen[n] = true
        unicos = append(unicos, n)
    }
}

fmt.Println(unicos)  // [1 2 3 4]
```

### Usar map[string]bool como Set

Go no tiene tipo Set nativo, usa map como substituto:

```go
// Set de strings
visitado := make(map[string]bool)

// Agregar
visitado["Juan"] = true
visitado["María"] = true

// Comprobar
if visitado["Juan"] {
    fmt.Println("Juan ya visitado")
}

// Remover
delete(visitado, "María")
```

---

## Ejercicios del Capítulo 9

### Ejercicio 1: Gestor de Inventario

Crea programa que:
1. Use map[string]int para productos → cantidad
2. Permitir: agregar, remover, vender producto
3. Mostrar inventario total
4. Buscar producto por nombre
5. Decrementar cantidad al vender (validar que existe)

### Ejercicio 2: Analizador de Frecuencia

Crea programa que:
1. Reciba texto
2. Cuente frecuencia de cada palabra
3. Ignore mayúsculas/minúsculas
4. Muestre palabras ordenadas por frecuencia (más común primero)
5. Excluya palabras comunes (the, and, or)

### Ejercicio 3: Traductor de Idiomas

Crea programa que:
1. Use map[string]map[string]string para idiomas
2. Estructura: idioma → (palabra_origen → palabra_traducida)
3. Traducir frase palabra por palabra
4. Manejar palabras no traducidas
5. Agregar nuevas palabras dinámicamente

### Ejercicio 4: Caché Simple

Crea programa que:
1. Implemente caché con map[string]interface{}
2. Funciones: Get, Set, Delete, Clear, Exists
3. Simule cálculos costosos almacenando resultados
4. Muestre estadísticas: hits, misses
5. Ordena claves antes de imprimir

### Ejercicio 5: Análisis de JSON Dinámico

Crea programa que:
1. Decodifique JSON en map[string]interface{}
2. Navegue estructura (puede ser anidada)
3. Busque clave específica recursivamente
4. Muestre tipo de cada valor
5. Cuenta total de claves (profundidad arbitraria)

---

**Fin del Capítulo 9**

---

