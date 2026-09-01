# Capítulo 7: Funciones

## Índice del Capítulo 7

1. [7.1 ¿Qué es una Función?](#71-qué-es-una-función)
2. [7.2 Definición y Llamada Básica](#72-definición-y-llamada-básica)
3. [7.3 Parámetros](#73-parámetros)
4. [7.4 Valores de Retorno](#74-valores-de-retorno)
5. [7.5 Múltiples Retornos](#75-múltiples-retornos)
6. [7.6 Named Returns](#76-named-returns)
7. [7.7 Funciones Variádicas](#77-funciones-variádicas)
8. [7.8 Funciones Anónimas y Closures](#78-funciones-anónimas-y-closures)
9. [7.9 Recursión](#79-recursión)
10. [7.10 Buenas Prácticas y Patrones](#710-buenas-prácticas-y-patrones)

---

## 7.1 ¿Qué es una Función?

### Definición Conceptual

Una función es un **bloque de código reutilizable** que:

```
1. Recibe ENTRADAS (parámetros)
2. Realiza PROCESAMIENTO
3. Produce SALIDAS (retornos)
```

### Por Qué Usar Funciones

```
SIN funciones:
    fmt.Println("Juan")
    fmt.Println("María")
    fmt.Println("Carlos")

CON funciones:
    func saludar(nombre string) {
        fmt.Printf("Hola, %s!\n", nombre)
    }
    
    saludar("Juan")
    saludar("María")
    saludar("Carlos")

Beneficios:
 DRY: Don't Repeat Yourself
 Reutilizable
 Legible
 Testeable
 Mantenible
```

### Anatomía de una Función

```
func nombre(parametros) (retornos) {
    // Cuerpo de la función
}

Ejemplo:
func sumar(a int, b int) int {
    return a + b
}

Componentes:
 func: Palabra clave
 sumar: Nombre de función
 (a int, b int): Parámetros
 int: Tipo de retorno
 return a + b: Cuerpo y retorno
```

---

## 7.2 Definición y Llamada Básica

### Función sin Parámetros

```go
func saludar() {
    fmt.Println("Hola, Mundo!")
}

// Llamada
saludar()       // Output: Hola, Mundo!
```

### Función con Parámetros

```go
func saludar(nombre string) {
    fmt.Printf("Hola, %s!\n", nombre)
}

// Llamada
saludar("Juan")         // Output: Hola, Juan!
saludar("María")        // Output: Hola, María!
```

### Función con Retorno

```go
func sumar(a int, b int) int {
    return a + b
}

// Llamada
resultado := sumar(5, 3)
fmt.Println(resultado)      // 8
```

### Función sin Retorno Explícito

```go
func procesar() {
    fmt.Println("Procesando...")
    // Sin return, retorna implícitamente
}

procesar()      // Output: Procesando...
```

---

## 7.3 Parámetros

### Parámetros de Mismo Tipo

```go
// Verbose
func sumar(a int, b int, c int) int {
    return a + b + c
}

// Conciso (tipos iguales)
func sumar(a, b, c int) int {
    return a + b + c
}

// Ambas son equivalentes
```

### Parámetros de Tipos Diferentes

```go
func procesar(nombre string, edad int, activo bool) {
    fmt.Printf("%s, %d años, activo: %v\n", nombre, edad, activo)
}

// Llamada
procesar("Juan", 25, true)
```

### Parámetros por Referencia vs por Valor

**Por valor (copia):**

```go
func incrementar(x int) {
    x++  // Modifica la COPIA, no la original
}

num := 5
incrementar(num)
fmt.Println(num)    // 5 (sin cambios)
```

**Por referencia (puntero):**

```go
func incrementar(x *int) {
    *x++    // Modifica la ORIGINAL
}

num := 5
incrementar(&num)
fmt.Println(num)    // 6 (modificado)
```

### Parámetros Complejos (Slices)

**Slices son referencias implícitas:**

```go
func modificarSlice(s []int) {
    if len(s) > 0 {
        s[0] = 100  // Modifica el slice original
    }
}

numeros := []int{1, 2, 3}
modificarSlice(numeros)
fmt.Println(numeros)    // [100 2 3] (modificado)
```

---

## 7.4 Valores de Retorno

### Un Valor de Retorno

```go
func obtenerEdad() int {
    return 25
}

edad := obtenerEdad()
fmt.Println(edad)       // 25
```

### Sin Especificar Retorno (Void)

```go
func imprimir(mensaje string) {
    fmt.Println(mensaje)
    // No hay return
}

imprimir("Hola")        // Output: Hola
```

### Retorno Vacío Explícito

```go
func procesarConError() error {
    if todo_ok {
        return nil          // No hay error
    }
    return errors.New("Error ocurrió")
}
```

---

## 7.5 Múltiples Retornos

### Go Permite Múltiples Retornos

```go
func dividirConResto(dividendo int, divisor int) (int, int) {
    cociente := dividendo / divisor
    resto := dividendo % divisor
    return cociente, resto
}

// Llamada
q, r := dividirConResto(10, 3)
fmt.Println(q, r)       // 3 1
```

### Ignorar Retornos

```go
// Ignorar segundo retorno
q, _ := dividirConResto(10, 3)
fmt.Println(q)          // 3

// Ignorar todos excepto uno
_, r := dividirConResto(10, 3)
fmt.Println(r)          // 1
```

### Patrón Error Handling

```go
func leerArchivo(nombre string) (string, error) {
    contenido, err := ioutil.ReadFile(nombre)
    if err != nil {
        return "", err      // Retorna error
    }
    return string(contenido), nil
}

// Uso
contenido, err := leerArchivo("archivo.txt")
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(contenido)
```

---

## 7.6 Named Returns

### ¿Qué son?

Named returns nombran explícitamente los valores de retorno:

```go
// Sin nombres (usual)
func dividir(a, b int) (int, int) {
    return a / b, a % b
}

// Con nombres (menos común)
func dividir(a, b int) (cociente int, resto int) {
    cociente = a / b
    resto = a % b
    return
}

// Ambas equivalentes
```

### Ventajas

```go
// Documenta qué retorna
func obtenerDatos() (nombre string, edad int, error error) {
    // Es claro qué se retorna
    return "Juan", 25, nil
}

// Versus
func obtenerDatos() (string, int, error) {
    // ¿Qué es qué?
    return "Juan", 25, nil
}
```

### Retorno Implícito

```go
func calcular(x, y int) (suma int, resta int) {
    suma = x + y
    resta = x - y
    return      // Retorna implícitamente suma y resta
}

s, r := calcular(10, 3)
fmt.Println(s, r)   // 13 7
```

---

## 7.7 Funciones Variádicas

### ¿Qué es una Variádica?

Función que acepta número variable de argumentos:

```go
func sumar(numeros ...int) int {
    total := 0
    for _, n := range numeros {
        total += n
    }
    return total
}

// Llamadas
fmt.Println(sumar())                    // 0
fmt.Println(sumar(5))                   // 5
fmt.Println(sumar(5, 10))              // 15
fmt.Println(sumar(5, 10, 15, 20))      // 50
```

### Internamente es un Slice

```go
func imprimir(valores ...interface{}) {
    // valores es []interface{}
    for i, v := range valores {
        fmt.Printf("Posición %d: %v\n", i, v)
    }
}

imprimir("texto", 42, true, 3.14)
```

### Unpacking de Slices

```go
numeros := []int{1, 2, 3, 4, 5}

// Pasar slice a función variádica
fmt.Println(sumar(numeros...))      // 15
// El ... desempaca el slice en argumentos individuales
```

### Variádica no es Último Parámetro (Raro)

```go
// ❌ NO PERMITIDO
func procesar(obligatorio int, opcionales ...int, final string) {
}

// ✅ PERMITIDO
func procesar(obligatorio int, final string, opcionales ...int) {
}
```

---

## 7.8 Funciones Anónimas y Closures

### Funciones Anónimas

Funciones sin nombre, usualmente asignadas a variables:

```go
// Función anónima asignada a variable
saludar := func(nombre string) {
    fmt.Printf("Hola, %s!\n", nombre)
}

// Llamar función anónima
saludar("Juan")     // Output: Hola, Juan!
```

### Ejecutar Inmediatamente

```go
// Función anónima ejecutada inmediatamente
func() {
    fmt.Println("Ejecuta inmediatamente")
}()         // ← Paréntesis para ejecutar

// Con argumentos
func(nombre string) {
    fmt.Println("Hola,", nombre)
}("Juan")   // ← Paréntesis y argumentos
```

### Closures - Capturar Variables

```go
func crearContador() func() int {
    contador := 0
    
    // Función interna que captura 'contador'
    return func() int {
        contador++
        return contador
    }
}

// Uso
contar := crearContador()
fmt.Println(contar())   // 1
fmt.Println(contar())   // 2
fmt.Println(contar())   // 3

// Cada llamada a crearContador() crea nuevo contador
contar2 := crearContador()
fmt.Println(contar2())  // 1 (contador nuevo)
```

### Closure Captura por Referencia

```go
multiplicador := 2

func() {
    fmt.Println(multiplicador)  // Ve valor actual
}()

multiplicador = 5
func() {
    fmt.Println(multiplicador)  // Ve 5
}()
```

### Funciones como Argumentos

```go
// Función que acepta otra función
func aplicar(x int, operacion func(int) int) int {
    return operacion(x)
}

// Usar
resultado := aplicar(5, func(n int) int {
    return n * 2
})

fmt.Println(resultado)  // 10
```

### Map/Filter/Reduce

```go
// Map: transformar lista
func mapear(numeros []int, fn func(int) int) []int {
    resultado := make([]int, len(numeros))
    for i, n := range numeros {
        resultado[i] = fn(n)
    }
    return resultado
}

// Uso
nums := []int{1, 2, 3, 4, 5}
cuadrados := mapear(nums, func(n int) int {
    return n * n
})
fmt.Println(cuadrados)  // [1 4 9 16 25]
```

---

## 7.9 Recursión

### Definición

Función que se llama a sí misma:

```go
func factorial(n int) int {
    if n <= 1 {
        return 1            // Caso base
    }
    return n * factorial(n-1)  // Caso recursivo
}

fmt.Println(factorial(5))   // 120
```

### Visualización de Pila

```
factorial(5)
 5 * factorial(4)
 4 * factorial(3)
 3 * factorial(2)
 2 * factorial(1)
 1 (retorna 1)
 2 * 1 = 2
 3 * 2 = 6
 4 * 6 = 24
 5 * 24 = 120
```

### Caso Real: Traversal de Árbol

```go
type Nodo struct {
    Valor int
    Izq   *Nodo
    Der   *Nodo
}

func sumarArbol(nodo *Nodo) int {
    if nodo == nil {
        return 0
    }
    // Suma el nodo + izquierda + derecha
    return nodo.Valor + sumarArbol(nodo.Izq) + sumarArbol(nodo.Der)
}
```

### Cuidado: Stack Overflow

```go
// ❌ PELIGROSO - Sin caso base
func infinita() {
    infinita()      // Se llama a sí misma infinitamente
}

// ✅ SEGURO - Con caso base
func factorial(n int) int {
    if n <= 1 {
        return 1    // Caso base: detiene recursión
    }
    return n * factorial(n-1)
}
```

---

## 7.10 Buenas Prácticas y Patrones

### Nombres Descriptivos

```go
// ❌ Confuso
func f(x int) int {
    return x * x
}

// ✅ Claro
func calcularCuadrado(numero int) int {
    return numero * numero
}
```

### Una Sola Responsabilidad

```go
// ❌ Hace demasiado
func procesarYGuardar(datos []int) error {
    // Procesa datos
    // Valida datos
    // Transforma datos
    // Guarda en archivo
    // Registra en base de datos
    // Notifica por email
}

// ✅ Responsabilidades separadas
func procesar(datos []int) ([]int, error) {
    // Solo procesa
}

func guardar(datos []int) error {
    // Solo guarda
}

func notificar(mensaje string) error {
    // Solo notifica
}
```

### Parámetros Claramente Tipados

```go
// ❌ Ambiguo
func crear(str string, int int, bool bool) {
}

// ✅ Claro
func crearUsuario(nombre string, edad int, activo bool) {
}
```

### Retornar Errores Explícitamente

```go
// ❌ Ignora errores
func leerArchivo(nombre string) string {
    contenido, _ := ioutil.ReadFile(nombre)
    return string(contenido)
}

// ✅ Propaga errores
func leerArchivo(nombre string) (string, error) {
    contenido, err := ioutil.ReadFile(nombre)
    if err != nil {
        return "", err
    }
    return string(contenido), nil
}
```

### Early Return para Simplificar

```go
// ❌ Anidado profundo
func validar(usuario *Usuario) error {
    if usuario != nil {
        if usuario.Edad >= 18 {
            if len(usuario.Email) > 0 {
                return nil
            } else {
                return errors.New("Email vacío")
            }
        } else {
            return errors.New("Menor de edad")
        }
    } else {
        return errors.New("Usuario nulo")
    }
}

// ✅ Early return (más claro)
func validar(usuario *Usuario) error {
    if usuario == nil {
        return errors.New("Usuario nulo")
    }
    if usuario.Edad < 18 {
        return errors.New("Menor de edad")
    }
    if len(usuario.Email) == 0 {
        return errors.New("Email vacío")
    }
    return nil
}
```

### Defer para Cleanup

```go
func procesar(archivo string) error {
    f, err := os.Open(archivo)
    if err != nil {
        return err
    }
    defer f.Close()     // Se ejecuta al final, aunque haya error
    
    // Procesar archivo
    // Si hay error aquí, defer aún se ejecuta
    return nil
}
```

### Constantes sobre Números Mágicos

```go
// ❌ Número mágico
func validarEdad(edad int) bool {
    return edad >= 18
}

// ✅ Constante clara
const EdadMayoriaEdad = 18

func validarEdad(edad int) bool {
    return edad >= EdadMayoriaEdad
}
```

### Funciones Puras (Recomendadas)

```go
// ✅ Función pura (sin efectos secundarios)
func sumar(a, b int) int {
    return a + b
}

// Función con efectos secundarios (print, IO)
func sumarYImprimir(a, b int) int {
    resultado := a + b
    fmt.Println(resultado)      // Efecto secundario
    return resultado
}

// Mejor: separar lógica de presentación
func sumar(a, b int) int {
    return a + b
}

func mostrarResultado(resultado int) {
    fmt.Println(resultado)
}
```

### Funciones Pequeñas

```go
// ❌ Función muy larga (100+ líneas)
func procesarDatos() {
    // Demasiada lógica aquí
}

// ✅ Funciones pequeñas y enfocadas
func cargarDatos() ([]int, error) {
    // Solo carga
}

func validarDatos(datos []int) error {
    // Solo valida
}

func transformarDatos(datos []int) []int {
    // Solo transforma
}

func guardarDatos(datos []int) error {
    // Solo guarda
}

// Orquestar en función principal
func procesarDatos() error {
    datos, err := cargarDatos()
    if err != nil {
        return err
    }
    if err := validarDatos(datos); err != nil {
        return err
    }
    datos = transformarDatos(datos)
    return guardarDatos(datos)
}
```

---

## Ejercicios del Capítulo 7

### Ejercicio 1: Calculadora Funcional

Crea programa que:
1. Defina funciones: sumar, restar, multiplicar, dividir
2. Cada función tome dos números
3. Cree función "calcular" que acepte operación y dos números
4. Maneja división por cero
5. Usa menú para elegir operación

### Ejercicio 2: Procesador de Listas

Crea programa que:
1. Defina función "mapear(nums []int, fn func(int) int) []int"
2. Defina función "filtrar(nums []int, predicado func(int) bool) []int"
3. Defina función "reducir(nums []int, inicial int, fn func(int, int) int) int"
4. Demuestra cada una con ejemplos

### Ejercicio 3: Generador de Contadores

Crea programa que:
1. Función "crearContador(inicio int) func() int" que retorna closure
2. Cada closure incrementa su contador independientemente
3. Crea 3 contadores diferentes
4. Prueba que cada uno es independiente

### Ejercicio 4: Traverse de Estructuras

Crea programa que:
1. Defina estructura Carpeta con nombre, subcarpetas y archivos
2. Función recursiva para contar total de archivos
3. Función recursiva para buscar archivo por nombre
4. Función recursiva para mostrar estructura en árbol

### Ejercicio 5: Recursión vs Iteración

Crea programa que:
1. Fibonacci recursivo
2. Fibonacci iterativo
3. Mide tiempo de ejecución de ambos (para n=35)
4. Compara resultados y performance
5. Explica diferencia en eficiencia

---

**Fin del Capítulo 7**

---

