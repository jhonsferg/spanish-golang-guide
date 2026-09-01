# Capítulo 5: Operadores y expresiones

## Índice del Capítulo 5

1. [5.1 ¿Qué es un Operador y una Expresión?](#51-qué-es-un-operador-y-una-expresión)
2. [5.2 Operadores Aritméticos](#52-operadores-aritméticos)
3. [5.3 Operadores de Comparación](#53-operadores-de-comparación)
4. [5.4 Operadores Lógicos](#54-operadores-lógicos)
5. [5.5 Operadores de Asignación](#55-operadores-de-asignación)
6. [5.6 Operadores Bit a Bit](#56-operadores-bit-a-bit)
7. [5.7 Precedencia y Asociatividad](#57-precedencia-y-asociatividad)
8. [5.8 Cortocircuito (Short-Circuit Evaluation)](#58-cortocircuito-short-circuit-evaluation)
9. [5.9 Incremento/Decremento](#59-incrementodecremento)
10. [5.10 Overflow, Underflow y Comportamiento Especial](#510-overflow-underflow-y-comportamiento-especial)

---

## 5.1 ¿Qué es un Operador y una Expresión?

### Definiciones

**Operador:** Símbolo que realiza una operación sobre uno o más valores

```go
+       // Operador: suma
*       // Operador: multiplicación
>       // Operador: mayor que
&&      // Operador: AND lógico
```

**Expresión:** Combinación de valores y operadores que produce un resultado

```go
3 + 5               // Expresión: produce 8
x > 10              // Expresión: produce true o false
nombre == "Juan"    // Expresión: produce true o false
```

**Comparación:**

```

 Concepto    │ Ejemplo        │
clear
 Operador    │ +              │
 Operandos   │ 3, 5           │
 Expresión   │ 3 + 5          │
 Resultado   │ 8              │

```

### Tipos de Operadores

```
GO TIENE 7 CATEGORÍAS DE OPERADORES:

1. ARITMÉTICOS         +, -, *, /, %
2. COMPARACIÓN        ==, !=, <, >, <=, >=
3. LÓGICOS            &&, ||, !
4. ASIGNACIÓN         =, +=, -=, *=, /=, %=, &=, |=, ^=, <<=, >>=
5. BIT A BIT          &, |, ^, <<, >>, &^
6. INCREMENTO/DECREM  ++, -- (solo postfix)
7. OTROS              & (address), * (dereference), . (selector)

Nos enfocaremos en 1-6.
```

---

## 5.2 Operadores Aritméticos

### Suma, Resta, Multiplicación

```go
a := 10
b := 3

suma := a + b           // 13
resta := a - b          // 7
multiplicacion := a * b // 30

fmt.Println(suma, resta, multiplicacion)
```

### División

```go
a := 10
b := 3

// División entera (trunca decimales)
division := a / b       // 3 (no 3.33)

// División con decimales
división_f := float64(a) / float64(b)  // 3.333...
```

**Importante: División por cero**

```go
x := 5 / 0      // ❌ PANIC: division by zero

// Tienes que validar
if b != 0 {
    x := a / b
} else {
    fmt.Println("Error: división por cero")
}
```

**Comparación con otros lenguajes:**

```
Go:         10 / 3 = 3   (enteros truncan)
Python 3:   10 / 3 = 3.333... (verdadera división)
C/C++:      10 / 3 = 3   (igual a Go, si ambos int)
Java:       10 / 3 = 3   (igual a Go, si ambos int)
```

### Módulo (Resto)

```go
a := 10
b := 3

resto := a % b          // 1 (porque 10 = 3*3 + 1)

// Usar para ciclos
for i := 0; i < 10; i++ {
    if i % 2 == 0 {
        fmt.Printf("%d es par\n", i)
    }
}
```

**Con números negativos:**

```go
// Resultado tiene signo del DIVIDENDO
fmt.Println(10 % 3)     // 1
fmt.Println(-10 % 3)    // -1 (negativo porque -10 es negativo)
fmt.Println(10 % -3)    // 1 (positivo porque 10 es positivo)
```

### Operaciones Unarias

```go
x := 5

+x      // 5 (unario positivo, raro)
-x      // -5 (unario negativo)
```

### Tipos Mixtos (Cuidado)

```go
var x int = 10
var y float64 = 3

z := x + y      // ❌ ERROR: cannot add int and float64

// Correcto: convertir
z := float64(x) + y     // ✅ 13.0
```

---

## 5.3 Operadores de Comparación

### Igualdad

```go
a := 5
b := 5
c := 10

a == b          // true
a == c          // false
a != c          // true
```

### Mayor/Menor

```go
a := 10
b := 5

a > b           // true
a < b           // false
a >= b          // true
a <= b          // false
```

### Con Strings

```go
s1 := "apple"
s2 := "banana"

s1 == s2        // false
s1 < s2         // true (comparación lexicográfica)
s1 <= s2        // true
```

**Lexicográfico significa: orden alfabético**

```go
"abc" < "abd"           // true (c < d)
"abc" < "abcd"          // true (más corto < más largo, si prefijo igual)
"ABC" < "abc"           // true (ASCII: mayúsculas < minsculas)
```

### Comparación de Valores Floating Point

```go
x := 0.1 + 0.2
y := 0.3

if x == y {
    fmt.Println("Iguales")
} else {
    fmt.Println("Diferentes")  // ← Se ejecuta (problemas flotantes)
}

// Correcto: comparar con tolerancia
epsilon := 0.0001
if math.Abs(x - y) < epsilon {
    fmt.Println("Aproximadamente iguales")
}
```

### Comparación de Pointers

```go
var x *int = nil
var y *int = nil

x == y          // true (ambos nil)
x == nil        // true
```

---

## 5.4 Operadores Lógicos

### AND (&&)

Ambos deben ser true:

```go
true && true    // true
true && false   // false
false && true   // false
false && false  // false
```

**Uso en condicionales:**

```go
edad := 25
licencia := true

if edad >= 18 && licencia {
    fmt.Println("Puedes conducir")
} else {
    fmt.Println("No puedes conducir")
}
```

### OR (||)

Al menos uno debe ser true:

```go
true || true    // true
true || false   // true
false || true   // true
false || false  // false
```

**Uso en condicionales:**

```go
es_admin := false
es_moderador := true

if es_admin || es_moderador {
    fmt.Println("Tiene acceso")
}
```

### NOT (!)

Invierte el valor:

```go
!true           // false
!false          // true
```

**Uso en condicionales:**

```go
activo := false

if !activo {
    fmt.Println("Usuario inactivo")
}

// Equivalente a:
if activo == false {
    fmt.Println("Usuario inactivo")
}
```

### Combinaciones

```go
edad := 25
ingresos := 50000
endeudado := true

// Múltiples condiciones
if edad >= 18 && ingresos > 30000 && !endeudado {
    fmt.Println("Aprobado para préstamo")
}

// Equivalente a
if (edad >= 18) && (ingresos > 30000) && (!endeudado) {
    fmt.Println("Aprobado para préstamo")
}
```

---

## 5.5 Operadores de Asignación

### Asignación Simple

```go
x := 5          // Asigna 5 a x
x = 10          // Reasigna 10 a x
```

### Asignación Compuesta

```go
x := 5

x += 3          // x = x + 3 = 8
x -= 2          // x = x - 2 = 6
x *= 2          // x = x * 2 = 12
x /= 3          // x = x / 3 = 4
x %= 3          // x = x % 3 = 1
```

**Tabla:**

```

 Operador    │ Equivalente       │
clear
 x += y      │ x = x + y        │
 x -= y      │ x = x - y        │
 x *= y      │ x = x * y        │
 x /= y      │ x = x / y        │
 x %= y      │ x = x % y        │
 x &= y      │ x = x & y        │
 x |= y      │ x = x | y        │
 x ^= y      │ x = x ^ y        │
 x <<= y     │ x = x << y       │
 x >>= y     │ x = x >> y       │

```

### Asignación Múltiple

```go
a, b := 1, 2            // a=1, b=2

// Intercambiar valores
a, b = b, a             // a=2, b=1

// Retornar múltiples valores
func divmod(a, b int) (int, int) {
    return a / b, a % b
}

cociente, resto := divmod(10, 3)  // cociente=3, resto=1
```

---

## 5.6 Operadores Bit a Bit

### AND de Bits (&)

Ambos bits deben ser 1:

```
  1101 (13)
& 1011 (11)
------
  1001 (9)
```

**En código:**

```go
x := 13  // 1101
y := 11  // 1011
z := x & y  // 1001 = 9

fmt.Println(z)  // 9
```

**Uso real: Verificar flags**

```go
// Definir permisos como bits
const (
    Lectura    = 1 << 0  // 0001 = 1
    Escritura  = 1 << 1  // 0010 = 2
    Ejecucion  = 1 << 2  // 0100 = 4
)

permisos := Lectura | Escritura  // 0011 = 3

// Verificar si tiene permiso de lectura
if permisos & Lectura != 0 {
    fmt.Println("Tiene lectura")
}
```

### OR de Bits (|)

Al menos uno debe ser 1:

```
  1101 (13)
| 1011 (11)
------
  1111 (15)
```

**En código:**

```go
x := 13  // 1101
y := 11  // 1011
z := x | y  // 1111 = 15

fmt.Println(z)  // 15
```

### XOR de Bits (^)

Bits deben ser diferentes:

```
  1101 (13)
^ 1011 (11)
------
  0110 (6)
```

**En código:**

```go
x := 13  // 1101
y := 11  // 1011
z := x ^ y  // 0110 = 6

fmt.Println(z)  // 6
```

**Uso real: Intercambiar sin variable temporal**

```go
a := 5      // 0101
b := 7      // 0111

a = a ^ b   // a = 0010 = 2
b = b ^ a   // b = 0101 = 5
a = a ^ b   // a = 0111 = 7

// Ahora a=7, b=5 (intercambiados)
```

### Desplazamiento Izquierdo (<<)

Corre bits a la izquierda, llena con 0 a la derecha:

```
1101 << 2 = 110100
```

**En código:**

```go
x := 5      // 0101
y := x << 2 // 10100 = 20

fmt.Println(y)  // 20
```

**Nota: Equivalente a multiplicar por 2^n**

```go
5 << 1      // 10 (5 * 2)
5 << 2      // 20 (5 * 4)
5 << 3      // 40 (5 * 8)
```

### Desplazamiento Derecho (>>)

Corre bits a la derecha, llena con 0 a la izquierda:

```
1101 >> 2 = 0011
```

**En código:**

```go
x := 20     // 10100
y := x >> 2 // 0101 = 5

fmt.Println(y)  // 5
```

**Nota: Equivalente a dividir por 2^n**

```go
20 >> 1     // 10 (20 / 2)
20 >> 2     // 5 (20 / 4)
20 >> 3     // 2 (20 / 8)
```

### AND NOT (&^)

Operador especial de Go (No existe en C/C++):

```
x &^ y   = x & ^y   (AND con NOT de y)
```

**Borra bits:**

```
  1101 (13)
&^0011 (3)
------
  1100 (12)
```

---

## 5.7 Precedencia y Asociatividad

### Tabla de Precedencia

```
ALTA PRECEDENCIA (se evalúa primero):
1. * / % << >> & &^        (Multiplicación, división, bits)
2. + -                      (Suma, resta)
3. == != < <= > >=          (Comparación)
4. &&                       (AND lógico)
BAJA PRECEDENCIA (se evalúa último):
5. ||                       (OR lógico)
```

**Ejemplo:**

```go
// Sin paréntesis
2 + 3 * 4           // 14 (no 20)
// Razón: * tiene mayor precedencia que +
// Se evalúa como: 2 + (3 * 4) = 2 + 12 = 14

// Con paréntesis (claridad)
(2 + 3) * 4         // 20
```

**Expresión compleja:**

```go
edad := 25
ingresos := 50000
tieneLicencia := true
es_criminales := false

// Difícil de leer
if edad > 18 && ingresos > 30000 || es_criminales && !tieneLicencia {
    // ...
}

// Mejor: usa paréntesis
if (edad > 18 && ingresos > 30000) || (es_criminales && !tieneLicencia) {
    // ...
}
```

### Asociatividad

Cuando dos operadores tienen igual precedencia, la **asociatividad** decide orden:

```
ASOCIATIVIDAD IZQUIERDA (más común):
x - y - z  =  (x - y) - z
10 - 5 - 2 =  (10 - 5) - 2 = 5 - 2 = 3

ASOCIATIVIDAD DERECHA (rara):
x ^ y ^ z  =  x ^ (y ^ z)
```

---

## 5.8 Cortocircuito (Short-Circuit Evaluation)

### AND Cortocircuito (&&)

Si el primer operando es false, el segundo NO se evalúa:

```go
x := false
y := func() bool {
    fmt.Println("y evaluada")   // NO se imprime
    return true
}()

if x && y {
    fmt.Println("Verdadero")    // NO se ejecuta
}

// Output: (nada)
// Razón: false && _ = false, no necesita evaluar y
```

**Ventaja: Evitar evaluaciones costosas**

```go
usuario := getUsuario()
if usuario != nil && usuario.IsAdmin() {
    // Si usuario es nil, IsAdmin() nunca se llama
    // Evita PANIC (llamar método en nil)
}
```

### OR Cortocircuito (||)

Si el primer operando es true, el segundo NO se evalúa:

```go
x := true
y := func() bool {
    fmt.Println("y evaluada")   // NO se imprime
    return false
}()

if x || y {
    fmt.Println("Verdadero")    // SÍ se ejecuta
}

// Output: Verdadero
// Razón: true || _ = true, no necesita evaluar y
```

**Ventaja: Proporcionar defaults**

```go
nombre := ""
nombreFinal := nombre || "Guest"  // ❌ ERROR: can't use || with strings

// Correcto con if
if nombre == "" {
    nombreFinal = "Guest"
}
```

---

## 5.9 Incremento/Decremento

### Posfix ++ y --

```go
x := 5

x++             // x = 6
x--             // x = 5

fmt.Println(x)  // 5
```

**Importante: SOLO postfix en Go**

```go
++x             // ❌ ERROR: Go no permite prefix ++ o --
--x             // ❌ ERROR

x++             // ✅ OK: postfix solo
x--             // ✅ OK: postfix solo
```

**NO retorna valor:**

```go
x := 5
y :=  ERROR: x++ no retorna valorx++        //

// Correcto
y := x
x++             // ✅ Dos líneas
```

**Diferencia con C++:**

```
C++:           Go:
i++            i++             // Postfix OK
++i            ++i             // ❌ ERROR
y = i++        y =  ERRORi++         //
i += 1         i += 1          // ✅ OK
```

### Uso en Bucles

```go
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// Equivalente a
for i := 0; i < 10; i += 1 {
    fmt.Println(i)
}
```

---

## 5.10 Overflow, Underflow y Comportamiento Especial

### Overflow de Enteros

Cuando un número excede el máximo de su tipo:

```go
var x int8 = 127     // Máximo para int8
x++                  // x = -128 (overflow silencioso)

fmt.Println(x)       // -128 ← ¡Cambió de 127 a -128!
```

**Diagrama:**

```
int8 es un círculo:

        0
   -1      1
-2           2
...         ...
   127    -128

Cuando sumas 1 a 127, wraps a -128 (overflow silencioso)
```

**Comportamiento por tipo:**

```go
// uint8: 0 a 255
var u uint8 = 255
u++             // u = 0 (overflow)

// int32: -2^31 a 2^31-1
var i int32 = 2147483647
i++             // i = -2147483648 (overflow)

// int (típicamente 64-bit en máquinas modernas)
var n int = 9_223_372_036_854_775_807
n++             // n = -9_223_372_036_854_775_808 (overflow)
```

### Prevenir Overflow

```go
// Validar antes de operación
x := int8(127)
const MaxInt8 = 127

if x < MaxInt8 {
    x++
} else {
    fmt.Println("Overflow detectado")
}

// O usar tipo más grande
x := int32(127)  // Más capacidad que int8
```

### Underflow de Enteros

Cuando un número baja del mínimo de su tipo:

```go
var x int8 = -128    // Mínimo para int8
x--                  // x = 127 (underflow)

fmt.Println(x)       // 127 ← ¡Cambió de -128 a 127!
```

### Comportamiento Especial de Floats

```go
// Infinito
var x float64 = 1.0 / 0.0       // +Inf
var y float64 = -1.0 / 0.0      // -Inf

// NaN (Not a Number)
var z float64 = 0.0 / 0.0       // NaN

// Comparación con infinito/NaN
fmt.Println(x == x)             // true
fmt.Println(z == z)             // false ← NaN nunca es igual a sí mismo
```

### División por Cero

```go
// Enteros: PANIC
x := 5 / 0      // panic: runtime error: integer divide by zero

// Floats: Infinito (sin panic)
y := 5.0 / 0.0  // +Inf (sin error)
```

---

## Ejercicios del Capítulo 5

### Ejercicio 1: Calculadora de Expresiones

Crea programa que:

1. Pida dos números
2. Pida operador (+, -, *, /, %)
3. Valide división por cero
4. Muestre resultado
5. Maneja tipos correctamente

### Ejercicio 2: Validador de Edad y Estado

Crea programa que:

1. Pida edad
2. Pida ingresos
3. Pida si tiene licencia
4. Use operadores lógicos para validar:
   - Puede trabajar si: edad >= 18
   - Puede conducir si: edad >= 18 && tiene_licencia
   - Puede solicitar préstamo si: edad >= 18 && ingresos > 50000
5. Muestre mensajes apropiados

### Ejercicio 3: Manipulador de Bits

Crea programa que:

1. Pida número (0-255)
2. Muestre su representación binaria
3. Realice operaciones:
   - Desplazamiento izquierdo
   - Desplazamiento derecho
   - AND, OR, XOR con otro número
4. Muestre resultados en binario

### Ejercicio 4: Comparador de Strings y Números

Crea programa que:

1. Pida dos valores (como strings)
2. Intente convertir a números
3. Si ambos son números: compara numéricamente
4. Si son strings: compara lexicográficamente
5. Muestra: mayor, menor, igual
6. Muestra diferencia de valor

### Ejercicio 5: Detector de Números Especiales

Crea programa que:

1. Pida número flotante
2. Valide si es:
   - Normal (número válido)
   - Infinito
   - NaN
3. Realiza operaciones:
   - Suma con infinito
   - Suma con NaN
   - División que produce infinito
4. Muestra comportamiento de cada uno

---

**Fin del Capítulo 5**

---

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/05-operadores-y-expresiones/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/05-operadores-y-expresiones):

```bash
cd examples/05-operadores-y-expresiones
go run .
```
