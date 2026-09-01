# Capítulo 32: Math y números aleatorios

## Introducción

El análisis matemático es fundamental en Go para aplicaciones científicas, criptografía, simulaciones y procesamiento de datos. Este capítulo explora exhaustivamente el paquete `math`, operaciones con números aleatorios, gestión de precisión y la implementación de algoritmos numéricos robustos.

Go proporciona múltiples herramientas para trabajar con números:
- **math**: Funciones matemáticas estándar (IEEE 754)
- **math/rand**: PRNG (Pseudo Random Number Generator) rápido
- **crypto/rand**: Generador criptográficamente seguro
- **big**: Aritmética de precisión arbitraria
- **complex**: Soporte nativo para números complejos

---

## 32.1 Math Package Fundamentos

### 32.1.1 Constantes Matemáticas Importantes

Go define constantes matemáticas fundamentales en el paquete `math`:

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	// Constantes predefinidas
	fmt.Printf("π (Pi)           = %.15f\n", math.Pi)
	fmt.Printf("e (Euler)        = %.15f\n", math.E)
	fmt.Printf("φ (Golden Ratio) = %.15f\n", math.Phi)
	fmt.Printf("√2               = %.15f\n", math.Sqrt2)
	fmt.Printf("√3               = %.15f\n", math.Sqrt3)
	fmt.Printf("√5               = %.15f\n", math.Sqrt5)
	fmt.Printf("1/√2             = %.15f\n", math.SqrtE)
	fmt.Printf("ln(2)            = %.15f\n", math.Ln2)
	fmt.Printf("ln(10)           = %.15f\n", math.Ln10)
	fmt.Printf("log₁₀(e)         = %.15f\n", math.Log10E)
	fmt.Printf("log₂(e)          = %.15f\n", math.Log2E)

	// Valores especiales
	fmt.Println("\nValores especiales:")
	fmt.Printf("MaxFloat64 = %.2e\n", math.MaxFloat64)
	fmt.Printf("MinFloat64 = %.2e\n", math.MinFloat64)
	fmt.Printf("SmallestNonzeroFloat64 = %.2e\n", math.SmallestNonzeroFloat64)
	fmt.Printf("MaxInt64 = %d\n", math.MaxInt64)
	fmt.Printf("MinInt64 = %d\n", math.MinInt64)

	// Valores especiales IEEE 754
	fmt.Println("\nIEEE 754 especiales:")
	fmt.Printf("Inf (positivo)  = %v\n", math.Inf(1))
	fmt.Printf("Inf (negativo)  = %v\n", math.Inf(-1))
	fmt.Printf("NaN (Not a Number) = %v\n", math.NaN())
}
```

**Salida esperada:**
```
π (Pi)           = 3.141592653589793
e (Euler)        = 2.718281828459045
φ (Golden Ratio) = 1.618033988749895
√2               = 1.414213562373095
√3               = 1.732050807568877
...
```

### 32.1.2 Funciones Básicas del Paquete math

El paquete `math` proporciona funciones fundamentales para operaciones numéricas:

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	// Operaciones básicas
	fmt.Println("=== OPERACIONES BÁSICAS ===")
	fmt.Printf("Abs(-5.5) = %v\n", math.Abs(-5.5))
	fmt.Printf("Abs(3.14) = %v\n", math.Abs(3.14))

	// Mínimo y Máximo
	fmt.Println("\n=== MIN/MAX ===")
	fmt.Printf("Min(3.5, 2.1) = %v\n", math.Min(3.5, 2.1))
	fmt.Printf("Max(3.5, 2.1) = %v\n", math.Max(3.5, 2.1))

	// Módulo
	fmt.Println("\n=== MÓDULO ===")
	fmt.Printf("Mod(7.5, 2.3) = %v\n", math.Mod(7.5, 2.3))
	fmt.Printf("Remainder(7.5, 2.3) = %v\n", math.Remainder(7.5, 2.3))

	// Signo
	fmt.Println("\n=== SIGNO ===")
	fmt.Printf("Copysign(5.5, -1.0) = %v\n", math.Copysign(5.5, -1.0))
	fmt.Printf("Copysign(-5.5, 1.0) = %v\n", math.Copysign(-5.5, 1.0))

	// Comprobaciones IEEE 754
	fmt.Println("\n=== COMPROBACIONES ===")
	fmt.Printf("IsNaN(math.NaN()) = %v\n", math.IsNaN(math.NaN()))
	fmt.Printf("IsInf(math.Inf(1), 1) = %v\n", math.IsInf(math.Inf(1), 1))
	fmt.Printf("IsInf(math.Inf(-1), -1) = %v\n", math.IsInf(math.Inf(-1), -1))
}
```

### 32.1.3 Tabla Comparativa: Go vs NumPy vs C

| Operación | Go (math) | NumPy (Python) | C (libm) | Rendimiento |
|-----------|-----------|---|---------|-------------|
| Sin(x) | `math.Sin()` | `np.sin()` | `sin()` | Go ≈ C > NumPy |
| Cos(x) | `math.Cos()` | `np.cos()` | `cos()` | Go ≈ C > NumPy |
| Sqrt(x) | `math.Sqrt()` | `np.sqrt()` | `sqrt()` | Go ≈ C > NumPy |
| Exp(x) | `math.Exp()` | `np.exp()` | `exp()` | Go ≈ C > NumPy |
| Log(x) | `math.Log()` | `np.log()` | `log()` | Go ≈ C > NumPy |
| Pow(x,y) | `math.Pow()` | `np.power()` | `pow()` | Go ≈ C > NumPy |

---

## 32.2 Trigonometría

### 32.2.1 Funciones Trigonométricas Básicas

Las funciones trigonométricas trabajan con radianes, no grados:

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	// Conversión radianes <-> grados
	degreesToRadians := func(d float64) float64 { return d * math.Pi / 180 }
	radiansToDegrees := func(r float64) float64 { return r * 180 / math.Pi }

	fmt.Println("=== FUNCIONES TRIGONOMÉTRICAS BÁSICAS ===")

	// Ángulos comunes
	angles := map[string]float64{
		"0°":   degreesToRadians(0),
		"30°":  degreesToRadians(30),
		"45°":  degreesToRadians(45),
		"60°":  degreesToRadians(60),
		"90°":  degreesToRadians(90),
		"180°": degreesToRadians(180),
	}

	for label, rad := range angles {
		fmt.Printf("%s (%.4f rad): sin=%.4f, cos=%.4f, tan=%.4f\n",
			label, rad, math.Sin(rad), math.Cos(rad), math.Tan(rad))
	}

	// Verificación: sin²(x) + cos²(x) = 1
	fmt.Println("\n=== IDENTIDAD TRIGONOMÉTRICA ===")
	for _, angle := range []float64{30, 45, 60, 90} {
		rad := degreesToRadians(angle)
		sin := math.Sin(rad)
		cos := math.Cos(rad)
		result := sin*sin + cos*cos
		fmt.Printf("sin²(%v°) + cos²(%v°) = %.15f (expected 1.0)\n", angle, angle, result)
	}

	// Funciones inversas (producen radianes)
	fmt.Println("\n=== FUNCIONES INVERSAS ===")
	fmt.Printf("Asin(0.5) = %.4f rad = %.4f°\n", 
		math.Asin(0.5), radiansToDegrees(math.Asin(0.5)))
	fmt.Printf("Acos(0.5) = %.4f rad = %.4f°\n",
		math.Acos(0.5), radiansToDegrees(math.Acos(0.5)))
	fmt.Printf("Atan(1.0) = %.4f rad = %.4f°\n",
		math.Atan(1.0), radiansToDegrees(math.Atan(1.0)))
	fmt.Printf("Atan2(1.0, 1.0) = %.4f rad = %.4f°\n",
		math.Atan2(1.0, 1.0), radiansToDegrees(math.Atan2(1.0, 1.0)))

	// Funciones hiperbólicas
	fmt.Println("\n=== FUNCIONES HIPERBÓLICAS ===")
	x := 1.0
	fmt.Printf("Sinh(%.1f) = %.4f\n", x, math.Sinh(x))
	fmt.Printf("Cosh(%.1f) = %.4f\n", x, math.Cosh(x))
	fmt.Printf("Tanh(%.1f) = %.4f\n", x, math.Tanh(x))
}
```

### 32.2.2 Atan2 y Coordenadas Polares

`Atan2(y, x)` calcula el ángulo en radianes para las coordenadas (x, y):

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	// Convertir entre coordenadas cartesianas y polares
	cartesianToPolar := func(x, y float64) (r, theta float64) {
		r = math.Sqrt(x*x + y*y)
		theta = math.Atan2(y, x)
		return
	}

	polarToCartesian := func(r, theta float64) (x, y float64) {
		x = r * math.Cos(theta)
		y = r * math.Sin(theta)
		return
	}

	fmt.Println("=== CONVERSIÓN DE COORDENADAS ===")
	points := [][2]float64{
		{1, 0},    // 0°
		{1, 1},    // 45°
		{0, 1},    // 90°
		{-1, 1},   // 135°
		{-1, 0},   // 180°
		{-1, -1},  // 225°
		{0, -1},   // 270°
		{1, -1},   // 315°
	}

	for _, p := range points {
		r, theta := cartesianToPolar(p[0], p[1])
		degrees := theta * 180 / math.Pi
		fmt.Printf("Cartesian (%.1f, %.1f) -> Polar (r=%.4f, θ=%.2f°)\n", 
			p[0], p[1], r, degrees)
	}

	// Verificación inversa
	fmt.Println("\n=== VERIFICACIÓN INVERSA ===")
	r, theta := 5.0, math.Pi/4
	x, y := polarToCartesian(r, theta)
	r2, theta2 := cartesianToPolar(x, y)
	fmt.Printf("Original: r=%.4f, θ=%.4f\n", r, theta)
	fmt.Printf("Conversiones: x=%.4f, y=%.4f\n", x, y)
	fmt.Printf("Recuperado: r=%.4f, θ=%.4f\n", r2, theta2)
}
```

### 32.2.3 Aplicación: Rotación de Puntos

```go
package main

import (
	"fmt"
	"math"
)

func rotatePoint(x, y, angle float64) (float64, float64) {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return x*cos - y*sin, x*sin + y*cos
}

func main() {
	fmt.Println("=== ROTACIÓN DE PUNTOS ===")
	
	point := [2]float64{1, 0}
	angles := []struct {
		name  string
		angle float64
	}{
		{"90°", math.Pi / 2},
		{"180°", math.Pi},
		{"270°", 3 * math.Pi / 2},
		{"360°", 2 * math.Pi},
	}

	for _, a := range angles {
		x, y := rotatePoint(point[0], point[1], a.angle)
		fmt.Printf("Punto (%.1f, %.1f) rotado %s: (%.4f, %.4f)\n",
			point[0], point[1], a.name, x, y)
	}
}
```

---

## 32.3 Exponenciales y Logaritmos

### 32.3.1 Funciones Exponenciales

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== EXPONENCIALES ===")

	// Exp(x) = e^x
	fmt.Println("\nExp(x) - e elevado a x:")
	for i := -2; i <= 3; i++ {
		x := float64(i)
		fmt.Printf("e^%.1f = %.6f\n", x, math.Exp(x))
	}

	// Exp2(x) = 2^x
	fmt.Println("\nExp2(x) - 2 elevado a x:")
	for i := -2; i <= 5; i++ {
		x := float64(i)
		fmt.Printf("2^%.1f = %.6f\n", x, math.Exp2(x))
	}

	// Exp10(x) = 10^x (simulado con Pow)
	fmt.Println("\n10^x usando Pow(10, x):")
	for i := -2; i <= 3; i++ {
		x := float64(i)
		fmt.Printf("10^%.1f = %.6f\n", x, math.Pow(10, x))
	}

	// Expm1(x) = e^x - 1 (preciso para x pequeño)
	fmt.Println("\nExpm1(x) para x pequeño (evita pérdida de precisión):")
	smallValues := []float64{1e-10, 1e-5, 0.1}
	for _, x := range smallValues {
		naive := math.Exp(x) - 1
		accurate := math.Expm1(x)
		fmt.Printf("x=%.0e: Exp(x)-1=%.15f, Expm1(x)=%.15f\n", x, naive, accurate)
	}
}
```

### 32.3.2 Funciones Logarítmicas

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== LOGARITMOS ===")

	// Log(x) = ln(x) - logaritmo natural
	fmt.Println("\nLog(x) - logaritmo natural:")
	for i := 1; i <= 5; i++ {
		x := float64(i)
		fmt.Printf("ln(%.1f) = %.6f\n", x, math.Log(x))
	}

	// Log2(x) - logaritmo base 2
	fmt.Println("\nLog2(x) - logaritmo base 2:")
	for i := 1; i <= 8; i *= 2 {
		x := float64(i)
		fmt.Printf("log₂(%.1f) = %.6f\n", x, math.Log2(x))
	}

	// Log10(x) - logaritmo base 10
	fmt.Println("\nLog10(x) - logaritmo base 10:")
	for i := 1; i <= 100000; i *= 10 {
		x := float64(i)
		fmt.Printf("log₁₀(%.0f) = %.6f\n", x, math.Log10(x))
	}

	// Log1p(x) = ln(1+x) - preciso para x pequeño
	fmt.Println("\nLog1p(x) para x pequeño (evita pérdida de precisión):")
	smallValues := []float64{1e-10, 1e-5, 0.1}
	for _, x := range smallValues {
		naive := math.Log(1 + x)
		accurate := math.Log1p(x)
		fmt.Printf("x=%.0e: Log(1+x)=%.15f, Log1p(x)=%.15f\n", x, naive, accurate)
	}

	// Cambio de base: log_b(x) = ln(x) / ln(b)
	fmt.Println("\nCambio de base: log_3(27):")
	x, base := 27.0, 3.0
	result := math.Log(x) / math.Log(base)
	fmt.Printf("log₃(27) = %.6f (expected 3.0)\n", result)
}
```

### 32.3.3 Potencias y Raíces

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== POTENCIAS Y RAÍCES ===")

	// Pow(x, y) = x^y
	fmt.Println("\nPow(x, y) - x elevado a y:")
	fmt.Printf("2^8 = %.0f\n", math.Pow(2, 8))
	fmt.Printf("10^3 = %.0f\n", math.Pow(10, 3))
	fmt.Printf("2^0.5 = %.6f (equivalente a √2)\n", math.Pow(2, 0.5))

	// Sqrt(x) = √x
	fmt.Println("\nSqrt(x) - raíz cuadrada:")
	for i := 1; i <= 10; i++ {
		x := float64(i)
		fmt.Printf("√%.0f = %.6f\n", x, math.Sqrt(x))
	}

	// Cbrt(x) = ∛x (raíz cúbica)
	fmt.Println("\nCbrt(x) - raíz cúbica:")
	for i := 1; i <= 8; i++ {
		x := float64(i)
		fmt.Printf("∛%.0f = %.6f\n", x, math.Cbrt(x))
	}

	// Hypot(x, y) = √(x² + y²) - distancia euclidiana
	fmt.Println("\nHypot(x, y) - hipotenusa (√(x²+y²)):")
	for _, p := range [][2]float64{{3, 4}, {5, 12}, {1, 1}} {
		dist := math.Hypot(p[0], p[1])
		fmt.Printf("Hypot(%.0f, %.0f) = %.6f\n", p[0], p[1], dist)
	}
}
```

### 32.3.4 Tabla Comparativa: Funciones Especiales

| Función | Propósito | Caso de Uso |
|---------|----------|-----------|
| `Exp(x)` | e^x genérico | Crecimiento exponencial, distribuciones |
| `Expm1(x)` | e^x - 1 (preciso para x pequeño) | Análisis numérico de precisión |
| `Log(x)` | ln(x) - logaritmo natural | Transformaciones de datos |
| `Log1p(x)` | ln(1+x) (preciso para x pequeño) | Análisis numérico de precisión |
| `Pow(x,y)` | x^y genérico | Cálculos científicos |
| `Sqrt(x)` | √x - raíz cuadrada | Geometría, estadística |
| `Cbrt(x)` | ∛x - raíz cúbica | Volumen, conversiones |
| `Hypot(x,y)` | √(x²+y²) optimizado | Distancias, geometría |

---

## 32.4 Rounding Functions

### 32.4.1 Funciones de Redondeo

Go proporciona varias funciones para redondear números flotantes:

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	values := []float64{
		2.3, 2.5, 2.7,
		-2.3, -2.5, -2.7,
		0.1, 0.9, 1.5,
	}

	fmt.Println("=== COMPARATIVA DE REDONDEO ===")
	fmt.Println("Valor    | Ceil | Floor | Round | Trunc")
	fmt.Println("-".repeat(40))

	for _, v := range values {
		fmt.Printf("%7.1f | %4.0f | %5.0f | %5.0f | %5.0f\n",
			v, math.Ceil(v), math.Floor(v), 
			math.Round(v), math.Trunc(v))
	}

	// Detalle de cada función
	fmt.Println("\n=== DETALLE DE FUNCIONES ===")

	// Ceil - siempre redondea hacia arriba
	fmt.Println("\nCeil (hacia +∞):")
	fmt.Printf("Ceil(3.2) = %.0f\n", math.Ceil(3.2))
	fmt.Printf("Ceil(3.0) = %.0f\n", math.Ceil(3.0))
	fmt.Printf("Ceil(-3.2) = %.0f\n", math.Ceil(-3.2))

	// Floor - siempre redondea hacia abajo
	fmt.Println("\nFloor (hacia -∞):")
	fmt.Printf("Floor(3.8) = %.0f\n", math.Floor(3.8))
	fmt.Printf("Floor(3.0) = %.0f\n", math.Floor(3.0))
	fmt.Printf("Floor(-3.8) = %.0f\n", math.Floor(-3.8))

	// Round - redondea al entero más cercano (bankers' rounding)
	fmt.Println("\nRound (al más cercano, .5 hacia par):")
	fmt.Printf("Round(2.5) = %.0f\n", math.Round(2.5))  // 2 (par)
	fmt.Printf("Round(3.5) = %.0f\n", math.Round(3.5))  // 4 (par)
	fmt.Printf("Round(2.3) = %.0f\n", math.Round(2.3))  // 2
	fmt.Printf("Round(2.7) = %.0f\n", math.Round(2.7))  // 3

	// Trunc - elimina la parte decimal
	fmt.Println("\nTrunc (hacia 0):")
	fmt.Printf("Trunc(3.8) = %.0f\n", math.Trunc(3.8))
	fmt.Printf("Trunc(3.2) = %.0f\n", math.Trunc(3.2))
	fmt.Printf("Trunc(-3.8) = %.0f\n", math.Trunc(-3.8))
	fmt.Printf("Trunc(-3.2) = %.0f\n", math.Trunc(-3.2))
}
```

### 32.4.2 Redondeo a Decimales Específicos

```go
package main

import (
	"fmt"
	"math"
)

// RoundTo redondea un float64 a n decimales
func roundTo(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}

// CeilTo redondea hacia arriba a n decimales
func ceilTo(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Ceil(value*multiplier) / multiplier
}

// FloorTo redondea hacia abajo a n decimales
func floorTo(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Floor(value*multiplier) / multiplier
}

func main() {
	value := math.Pi
	fmt.Printf("Valor original: %.15f\n\n", value)

	for decimals := 1; decimals <= 6; decimals++ {
		r := roundTo(value, decimals)
		c := ceilTo(value, decimals)
		f := floorTo(value, decimals)
		fmt.Printf("Decimales=%d: Round=%.6f, Ceil=%.6f, Floor=%.6f\n",
			decimals, r, c, f)
	}

	// Ejemplo práctico: redondeo de precios
	fmt.Println("\n=== REDONDEO DE PRECIOS ===")
	prices := []float64{19.994, 19.995, 19.996, 20.004}
	for _, price := range prices {
		rounded := roundTo(price, 2)
		fmt.Printf("$%.3f -> $%.2f\n", price, rounded)
	}
}
```

---

## 32.5 Números Complejos

### 32.5.1 Fundamentals de Números Complejos

Go tiene soporte nativo para números complejos con el tipo `complex64` y `complex128`:

```go
package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func main() {
	fmt.Println("=== NÚMEROS COMPLEJOS ===")

	// Creación de números complejos
	var z1 complex128 = 3 + 4i
	var z2 complex128 = 1 - 2i

	fmt.Printf("z1 = %v\n", z1)
	fmt.Printf("z2 = %v\n", z2)

	// Alternativa: usar complex()
	z3 := complex(2, 5)
	fmt.Printf("z3 = complex(2, 5) = %v\n", z3)

	// Componentes real e imaginario
	fmt.Println("\n=== COMPONENTES ===")
	fmt.Printf("z1 = %v\n", z1)
	fmt.Printf("Real(z1) = %v\n", real(z1))
	fmt.Printf("Imag(z1) = %v\n", imag(z1))

	// Operaciones aritméticas
	fmt.Println("\n=== OPERACIONES ARITMÉTICAS ===")
	fmt.Printf("z1 + z2 = %v + %v = %v\n", z1, z2, z1+z2)
	fmt.Printf("z1 - z2 = %v - %v = %v\n", z1, z2, z1-z2)
	fmt.Printf("z1 * z2 = %v * %v = %v\n", z1, z2, z1*z2)
	fmt.Printf("z1 / z2 = %v / %v = %v\n", z1, z2, z1/z2)

	// Función del paquete cmplx
	fmt.Println("\n=== FUNCIONES ESPECIALES (cmplx) ===")
	fmt.Printf("Abs(z1) = |%v| = %v\n", z1, cmplx.Abs(z1))
	fmt.Printf("Conj(z1) = %v* = %v\n", z1, cmplx.Conj(z1))
	fmt.Printf("Exp(z1) = e^%v = %v\n", z1, cmplx.Exp(z1))
	fmt.Printf("Log(z1) = ln(%v) = %v\n", z1, cmplx.Log(z1))
	fmt.Printf("Sqrt(z1) = √%v = %v\n", z1, cmplx.Sqrt(z1))
	fmt.Printf("Sin(z1) = sin(%v) = %v\n", z1, cmplx.Sin(z1))
	fmt.Printf("Cos(z1) = cos(%v) = %v\n", z1, cmplx.Cos(z1))

	// Conversión a forma polar
	fmt.Println("\n=== FORMA POLAR ===")
	abs := cmplx.Abs(z1)
	phase := cmplx.Phase(z1)
	fmt.Printf("z1 = %v\n", z1)
	fmt.Printf("Magnitud = %.4f\n", abs)
	fmt.Printf("Fase = %.4f rad = %.2f°\n", phase, phase*180/math.Pi)

	// Reconstrucción desde forma polar: z = r*e^(iθ)
	z_reconstructed := cmplx.Rect(abs, phase)
	fmt.Printf("Reconstruida = %.4f + %.4fi\n", real(z_reconstructed), imag(z_reconstructed))
}
```

### 32.5.2 Aplicación: Rotación con Números Complejos

```go
package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func main() {
	fmt.Println("=== ROTACIÓN USANDO NÚMEROS COMPLEJOS ===")

	// Punto en el plano como número complejo
	point := complex(1, 0) // (1, 0) en coordenadas cartesianas

	angles := []struct {
		name string
		rad  float64
	}{
		{"90°", math.Pi / 2},
		{"180°", math.Pi},
		{"270°", 3 * math.Pi / 2},
		{"360°", 2 * math.Pi},
	}

	for _, a := range angles {
		// Rotación: z' = z * e^(iθ)
		rotated := point * cmplx.Exp(complex(0, a.rad))
		
		fmt.Printf("Punto (1, 0) rotado %s:\n", a.name)
		fmt.Printf("  Resultado: (%.4f, %.4f)\n", real(rotated), imag(rotated))
	}

	// Ejemplo: cálculo de impedancia en circuitos AC
	fmt.Println("\n=== CIRCUITO AC: IMPEDANCIA ===")
	R := complex(10, 0)   // Resistencia: 10 Ω
	XL := complex(0, 5)   // Reactancia inductiva: 5 Ω
	XC := complex(0, -3)  // Reactancia capacitiva: -3 Ω

	Z := R + XL + XC     // Impedancia total
	fmt.Printf("Z = R + XL + XC = %v + %v + %v\n", R, XL, XC)
	fmt.Printf("Z_total = %v Ω\n", Z)
	fmt.Printf("|Z| = %.4f Ω (magnitud)\n", cmplx.Abs(Z))
	fmt.Printf("φ = %.4f rad = %.2f° (ángulo)\n", 
		cmplx.Phase(Z), cmplx.Phase(Z)*180/math.Pi)
}
```

---

## 32.6 Random Numbers - math/rand

### 32.6.1 Fundamentos de Generadores Pseudoaleatorios

Go proporciona `math/rand` para números pseudoaleatorios rápidos. Este es un PRNG (Pseudo Random Number Generator), no criptográficamente seguro:

```go
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("=== GENERADOR PSEUDOALEATORIO (math/rand) ===")

	// IMPORTANTE: Establecer seed para reproducibilidad
	rand.Seed(42)

	fmt.Println("\nSin seed (determinístico con seed 42):")
	for i := 0; i < 5; i++ {
		fmt.Printf("random número %d: %d\n", i+1, rand.Intn(100))
	}

	// Reset con la misma seed: produce los mismos números
	fmt.Println("\nReset con seed 42 (debería ser idéntico):")
	rand.Seed(42)
	for i := 0; i < 5; i++ {
		fmt.Printf("random número %d: %d\n", i+1, rand.Intn(100))
	}

	// Seed diferente: produce números diferentes
	fmt.Println("\nCon seed 123 (diferente):")
	rand.Seed(123)
	for i := 0; i < 5; i++ {
		fmt.Printf("random número %d: %d\n", i+1, rand.Intn(100))
	}
}
```

### 32.6.2 Funciones Principales del Paquete rand

```go
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	rand.Seed(42)

	fmt.Println("=== FUNCIONES PRINCIPALES DE math/rand ===")

	// Intn(n) - entero en [0, n)
	fmt.Println("\nIntn(n) - entero en [0, n):")
	fmt.Printf("Intn(100) = %d\n", rand.Intn(100))
	fmt.Printf("Intn(100) = %d\n", rand.Intn(100))
	fmt.Printf("Intn(100) = %d\n", rand.Intn(100))

	// Int63() - entero de 63 bits
	fmt.Println("\nInt63() - entero sin signo de 63 bits:")
	fmt.Printf("Int63() = %d\n", rand.Int63())

	// Int() - entero con signo (plataforma específica)
	fmt.Println("\nInt():")
	fmt.Printf("Int() = %d\n", rand.Int())

	// Float64() - float64 en [0.0, 1.0)
	fmt.Println("\nFloat64() - flotante en [0.0, 1.0):")
	for i := 0; i < 3; i++ {
		fmt.Printf("Float64() = %.6f\n", rand.Float64())
	}

	// Float32() - float32 en [0.0, 1.0)
	fmt.Println("\nFloat32() - flotante en [0.0, 1.0):")
	for i := 0; i < 3; i++ {
		fmt.Printf("Float32() = %.6f\n", rand.Float32())
	}

	// NormFloat64() - distribución normal N(0,1)
	fmt.Println("\nNormFloat64() - distribución normal estándar:")
	for i := 0; i < 5; i++ {
		fmt.Printf("NormFloat64() = %.4f\n", rand.NormFloat64())
	}

	// ExpFloat64() - distribución exponencial (λ=1)
	fmt.Println("\nExpFloat64() - distribución exponencial:")
	for i := 0; i < 5; i++ {
		fmt.Printf("ExpFloat64() = %.4f\n", rand.ExpFloat64())
	}

	// Perm(n) - permutación de [0, n)
	fmt.Println("\nPerm(10) - permutación de [0, 10):")
	fmt.Printf("%v\n", rand.Perm(10))

	// Shuffle(n, swap) - mezcla un slice
	fmt.Println("\nShuffle - mezclar slice:")
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})
	fmt.Printf("%v\n", numbers)
}
```

### 32.6.3 Generadores Independientes (Seguridad de Concurrencia)

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	fmt.Println("=== GENERADORES INDEPENDIENTES ===")

	// PROBLEMA: math/rand global no es seguro para goroutines
	// SOLUCIÓN: crear Rand independientes con mutex

	// Enfoque 1: Usar rand.New() con una fuente
	fmt.Println("\nEnfoque 1: rand.New() + sync.Mutex")
	source := rand.NewSource(42)
	generator := rand.New(source)

	for i := 0; i < 5; i++ {
		fmt.Printf("Generador local: %d\n", generator.Intn(100))
	}

	// Enfoque 2: Usar múltiples generadores en goroutines
	fmt.Println("\nEnfoque 2: Múltiples goroutines con Rand independientes")
	var wg sync.WaitGroup

	for g := 0; g < 3; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			gen := rand.New(rand.NewSource(int64(42 + goroutineID)))
			for i := 0; i < 3; i++ {
				fmt.Printf("Goroutine %d: %d\n", goroutineID, gen.Intn(100))
			}
		}(g)
	}

	wg.Wait()

	// ANTIPATRÓN: No hacer esto
	fmt.Println("\n❌ ANTIPATRÓN (no seguro para concurrencia):")
	fmt.Println("   rand.Seed(42)")
	fmt.Println("   // Múltiples goroutines usando rand global")
	fmt.Println("   // → Condiciones de carrera")

	// PATRÓN CORRECTO:
	fmt.Println("\n✓ PATRÓN CORRECTO:")
	fmt.Println("   for ... {")
	fmt.Println("       go func() {")
	fmt.Println("           gen := rand.New(rand.NewSource(seed))")
	fmt.Println("           // usar gen")
	fmt.Println("       }()")
	fmt.Println("   }")
}
```

---

## 32.7 Cryptographically Secure - crypto/rand

### 32.7.1 Generación Segura de Números Aleatorios

Para aplicaciones criptográficas, usar `crypto/rand` en lugar de `math/rand`:

```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("=== GENERACIÓN CRIPTOGRÁFICAMENTE SEGURA ===")

	// Read() - llenar buffer con bytes aleatorios seguros
	fmt.Println("\n1. Bytes aleatorios seguros:")
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Bytes aleatorios: %v\n", buffer)
	fmt.Printf("Hex: %s\n", hex.EncodeToString(buffer))
	fmt.Printf("Base64: %s\n", base64.StdEncoding.EncodeToString(buffer))

	// Generador de tokens seguros
	fmt.Println("\n2. Token de sesión (32 bytes):")
	token := make([]byte, 32)
	rand.Read(token)
	fmt.Printf("Token: %s\n", base64.URLEncoding.EncodeToString(token))

	// Entero aleatorio grande seguro
	fmt.Println("\n3. Entero grande aleatorio seguro:")
	max := big.NewInt(1000000)
	randomBig, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Entero aleatorio < 1,000,000: %v\n", randomBig)

	// Enteros seguros en rango específico
	fmt.Println("\n4. Entero aleatorio seguro en [1, 100]:")
	for i := 0; i < 5; i++ {
		// Generar en [0, 100), luego sumar 1 para [1, 101)
		randomInt, _ := rand.Int(rand.Reader, big.NewInt(100))
		fmt.Printf("  %d\n", randomInt.Int64()+1)
	}

	// Comparación: math/rand vs crypto/rand
	fmt.Println("\n=== COMPARATIVA ===")
	fmt.Println("math/rand:")
	fmt.Println("  ✓ Rápido")
	fmt.Println("  ✓ Determinístico (reproducible)")
	fmt.Println("  ✓ Suficiente para simulaciones")
	fmt.Println("  ✗ Predecible (no seguro para crypto)")
	fmt.Println()
	fmt.Println("crypto/rand:")
	fmt.Println("  ✓ Seguro (impredecible)")
	fmt.Println("  ✓ Para tokens, contraseñas, claves")
	fmt.Println("  ✗ Más lento que math/rand")
	fmt.Println("  ✗ No determinístico")
}
```

### 32.7.2 Generador de Tokens y UUIDs

```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// GenerateToken genera un token seguro en base64
func generateToken(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buffer), nil
}

// GenerateHexToken genera un token seguro en hexadecimal
func generateHexToken(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

// GenerateOTP genera un código OTP de n dígitos
func generateOTP(digits int) (string, error) {
	buffer := make([]byte, digits)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	
	otp := ""
	for _, b := range buffer {
		otp += fmt.Sprintf("%d", b%10)
	}
	return otp, nil
}

func main() {
	fmt.Println("=== GENERACIÓN DE TOKENS SEGUROS ===")

	// Tokens de sesión
	fmt.Println("\n1. Tokens de sesión (Base64):")
	for i := 0; i < 3; i++ {
		token, _ := generateToken(32)
		fmt.Printf("   %s\n", token)
	}

	// Tokens hexadecimales
	fmt.Println("\n2. Tokens hexadecimales (32 bytes = 256 bits):")
	for i := 0; i < 3; i++ {
		token, _ := generateHexToken(32)
		fmt.Printf("   %s\n", token)
	}

	// Códigos OTP
	fmt.Println("\n3. Códigos OTP de 6 dígitos:")
	for i := 0; i < 5; i++ {
		otp, _ := generateOTP(6)
		fmt.Printf("   %s\n", otp)
	}
}
```

---

## 32.8 Distribuciones Aleatorias

### 32.8.1 Distribución Normal

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// NormalDistribution genera números de una distribución normal
// usando el método Box-Muller
func normalsWithBoxMuller(n int, seed int64) []float64 {
	rand.Seed(seed)
	result := make([]float64, n)
	
	for i := 0; i < n; i += 2 {
		u1 := rand.Float64()
		u2 := rand.Float64()
		
		// Asegurar que u1 > 0 para evitar log(0)
		if u1 < 1e-10 {
			u1 = 1e-10
		}
		
		mag := math.Sqrt(-2.0 * math.Log(u1))
		z0 := mag * math.Cos(2.0*math.Pi*u2)
		
		if i < n {
			result[i] = z0
		}
		if i+1 < n {
			z1 := mag * math.Sin(2.0*math.Pi*u2)
			result[i+1] = z1
		}
	}
	
	return result
}

// estadísticas calcula media y desviación estándar
func statistics(data []float64) (mean, stddev float64) {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean = sum / float64(len(data))
	
	sumSq := 0.0
	for _, v := range data {
		diff := v - mean
		sumSq += diff * diff
	}
	stddev = math.Sqrt(sumSq / float64(len(data)))
	
	return
}

func main() {
	fmt.Println("=== DISTRIBUCIÓN NORMAL ===")

	// Generar con math/rand.NormFloat64()
	fmt.Println("\n1. Usando math/rand.NormFloat64() (N(0,1)):")
	rand.Seed(42)
	samples := make([]float64, 1000)
	for i := range samples {
		samples[i] = rand.NormFloat64()
	}
	
	mean, stddev := statistics(samples)
	fmt.Printf("Media = %.4f (esperada 0.0)\n", mean)
	fmt.Printf("Desv.Est. = %.4f (esperada 1.0)\n", stddev)

	// Generar con Box-Muller
	fmt.Println("\n2. Usando Box-Muller:")
	samplesBoxMuller := normalsWithBoxMuller(1000, 42)
	mean, stddev = statistics(samplesBoxMuller)
	fmt.Printf("Media = %.4f (esperada 0.0)\n", mean)
	fmt.Printf("Desv.Est. = %.4f (esperada 1.0)\n", stddev)

	// Distribución normal con media y desviación específicas
	// X ~ N(μ, σ) = μ + σ * Z, donde Z ~ N(0,1)
	fmt.Println("\n3. N(100, 15) - QI distribution:")
	mu, sigma := 100.0, 15.0
	rand.Seed(42)
	for i := 0; i < 10; i++ {
		value := mu + sigma*rand.NormFloat64()
		fmt.Printf("QI = %.1f\n", value)
	}

	// Histograma simple
	fmt.Println("\n4. Histograma de distribución normal:")
	rand.Seed(42)
	histogram := make(map[int]int)
	for i := 0; i < 10000; i++ {
		value := 100 + 15*rand.NormFloat64()
		bucket := int(value) / 5 * 5
		histogram[bucket]++
	}
	
	// Ordenar y mostrar
	var keys []int
	for k := range histogram {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	
	for _, k := range keys {
		if k >= 40 && k <= 160 {
			count := histogram[k]
			bar := ""
			for i := 0; i < count/100; i++ {
				bar += "█"
			}
			fmt.Printf("%3d-%3d: %s (%d)\n", k, k+5, bar, count)
		}
	}
}
```

### 32.8.2 Distribución Exponencial

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	fmt.Println("=== DISTRIBUCIÓN EXPONENCIAL ===")

	fmt.Println("\n1. Usando math/rand.ExpFloat64() (λ=1):")
	rand.Seed(42)
	sum := 0.0
	count := 10000
	for i := 0; i < count; i++ {
		sum += rand.ExpFloat64()
	}
	mean := sum / float64(count)
	fmt.Printf("Media = %.4f (esperada 1.0)\n", mean)

	// Exponencial con λ personalizado: X ~ Exp(λ)
	// P(X=x) = λ*e^(-λ*x)
	fmt.Println("\n2. Exponencial con λ=0.5:")
	lambda := 0.5
	rand.Seed(42)
	for i := 0; i < 10; i++ {
		value := rand.ExpFloat64() / lambda
		fmt.Printf("  Tiempo de espera = %.4f\n", value)
	}

	// Aplicación: Tiempo de llegada de eventos (proceso Poisson)
	fmt.Println("\n3. Tiempo entre eventos (Proceso Poisson, λ=2 eventos/min):")
	lambda = 2.0
	rand.Seed(42)
	time := 0.0
	eventNumber := 0
	for eventNumber < 10 {
		// Tiempo entre eventos ~ Exp(λ)
		timeBetweenEvents := rand.ExpFloat64() / lambda * 60 // en segundos
		time += timeBetweenEvents
		eventNumber++
		fmt.Printf("Evento %2d: t=%.1fs\n", eventNumber, time)
	}

	// Comparación: densidad de probabilidad teórica vs empírica
	fmt.Println("\n4. Densidad de probabilidad (λ=1):")
	lambda = 1.0
	rand.Seed(42)
	histogram := make(map[int]int)
	
	for i := 0; i < 10000; i++ {
		value := rand.ExpFloat64()
		bucket := int(value*2) // 0.0-0.5 -> 0, 0.5-1.0 -> 1, etc.
		if bucket < 10 {
			histogram[bucket]++
		}
	}

	for i := 0; i < 10; i++ {
		x := float64(i) / 2.0
		empirical := histogram[i] / 100.0
		theoretical := lambda * math.Exp(-lambda*x)
		bar := ""
		for j := 0; j < int(empirical*5); j++ {
			bar += "█"
		}
		fmt.Printf("[%.1f-%.1f) Emp=%.3f Teó=%.3f %s\n", x, x+0.5, empirical, theoretical, bar)
	}
}
```

### 32.8.3 Distribución de Poisson

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
)

// PoissonKnuth genera un número de Poisson(λ) usando el algoritmo de Knuth
func poissonKnuth(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	
	for p > L {
		k++
		p *= rand.Float64()
	}
	
	return k - 1
}

func main() {
	fmt.Println("=== DISTRIBUCIÓN DE POISSON ===")

	fmt.Println("\n1. Poisson(λ=3) usando Algoritmo de Knuth:")
	rand.Seed(42)
	lambda := 3.0
	sum := 0
	count := 10000
	
	for i := 0; i < count; i++ {
		sum += poissonKnuth(lambda)
	}
	
	empiricalMean := float64(sum) / float64(count)
	fmt.Printf("Media empírica = %.4f (esperada %.1f)\n", empiricalMean, lambda)

	// Distribución de probabilidad
	fmt.Println("\n2. Distribución P(X=k) para λ=3:")
	histogram := make(map[int]int)
	rand.Seed(42)
	
	for i := 0; i < 10000; i++ {
		k := poissonKnuth(3.0)
		if k < 15 {
			histogram[k]++
		}
	}

	for k := 0; k < 12; k++ {
		count := histogram[k]
		empirical := float64(count) / 100.0
		
		// Probabilidad teórica: P(X=k) = (e^-λ * λ^k) / k!
		factorial := 1.0
		for i := 1; i <= k; i++ {
			factorial *= float64(i)
		}
		theoretical := math.Exp(-lambda) * math.Pow(lambda, float64(k)) / factorial
		
		bar := ""
		for i := 0; i < count/50; i++ {
			bar += "█"
		}
		fmt.Printf("P(X=%2d) Emp=%.4f Teó=%.4f %s\n", k, empirical, theoretical, bar)
	}

	// Aplicación: Llamadas telefónicas, errores por página, etc.
	fmt.Println("\n3. Aplicación: Errores por línea de código (λ=2 errores/100 líneas):")
	rand.Seed(42)
	lambda = 2.0
	for i := 0; i < 5; i++ {
		errors := poissonKnuth(lambda)
		fmt.Printf("Archivo %d: %d errores detectados\n", i+1, errors)
	}
}
```

---

## 32.9 Big Numbers - Aritmética de Precisión Arbitraria

### 32.9.1 Trabajar con big.Int

```go
package main

import (
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("=== big.Int - NÚMEROS ENTEROS ARBITRARIOS ===")

	// Crear big.Int
	fmt.Println("\n1. Creación de big.Int:")

	// Desde cadena
	x := new(big.Int)
	x.SetString("123456789123456789123456789", 10)
	fmt.Printf("x = %v\n", x)

	// Desde int64
	y := big.NewInt(987654321)
	fmt.Printf("y = %v\n", y)

	// Operaciones básicas
	fmt.Println("\n2. Operaciones Aritméticas:")
	
	sum := new(big.Int).Add(x, y)
	fmt.Printf("x + y = %v\n", sum)

	diff := new(big.Int).Sub(x, y)
	fmt.Printf("x - y = %v\n", diff)

	prod := new(big.Int).Mul(x, y)
	fmt.Printf("x * y = %v\n", prod)

	quo := new(big.Int).Quo(x, y)
	rem := new(big.Int).Rem(x, y)
	fmt.Printf("x / y = %v (resto %v)\n", quo, rem)

	// Potencias
	fmt.Println("\n3. Potencias:")
	pow := new(big.Int).Exp(big.NewInt(2), big.NewInt(100), nil)
	fmt.Printf("2^100 = %v\n", pow)

	pow = new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil)
	fmt.Printf("10^50 = %v\n", pow)

	// Comparaciones
	fmt.Println("\n4. Comparaciones:")
	a := big.NewInt(100)
	b := big.NewInt(200)

	fmt.Printf("100 < 200? %v\n", a.Cmp(b) < 0)
	fmt.Printf("100 > 200? %v\n", a.Cmp(b) > 0)
	fmt.Printf("100 == 100? %v\n", a.Cmp(a) == 0)

	// Máximo y Mínimo
	fmt.Println("\n5. Máximo y Mínimo:")
	fmt.Printf("Max(100, 200) = %v\n", new(big.Int).Max(a, b))
	fmt.Printf("Min(100, 200) = %v\n", new(big.Int).Min(a, b))

	// Factorial con big.Int
	fmt.Println("\n6. Factorial(100):")
	factorial := func(n int64) *big.Int {
		result := big.NewInt(1)
		for i := int64(2); i <= n; i++ {
			result.Mul(result, big.NewInt(i))
		}
		return result
	}
	
	fact100 := factorial(100)
	fmt.Printf("100! = %v\n", fact100)
	fmt.Printf("Longitud: %d dígitos\n", len(fact100.String()))
}
```

### 32.9.2 Trabajar con big.Float

```go
package main

import (
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("=== big.Float - FLOTANTES DE PRECISIÓN ARBITRARIA ===")

	fmt.Println("\n1. Creación y precisión:")
	
	// Precisión: número de bits en la mantisa
	x := new(big.Float).SetPrec(100) // 100 bits de precisión
	x.SetString("3.141592653589793238462643383279502884197")
	fmt.Printf("π (100 bits) = %s\n", x.String())

	y := new(big.Float).SetPrec(200)
	y.SetString("2.718281828459045235360287471352662497757247093699959574966967627724076630353547594571382178525166427")
	fmt.Printf("e (200 bits) = %s\n", y.String())

	// Operaciones aritméticas
	fmt.Println("\n2. Operaciones Aritméticas:")
	
	a := new(big.Float).SetPrec(50).SetFloat64(1.0)
	b := new(big.Float).SetPrec(50).SetFloat64(3.0)
	
	div := new(big.Float).Quo(a, b)
	fmt.Printf("1/3 (50 bits) = %.30f\n", div)

	// Mayor precisión
	a2 := new(big.Float).SetPrec(200).SetFloat64(1.0)
	b2 := new(big.Float).SetPrec(200).SetFloat64(3.0)
	div2 := new(big.Float).Quo(a2, b2)
	fmt.Printf("1/3 (200 bits) = %.60f\n", div2)

	// Raíces cuadradas
	fmt.Println("\n3. Raíces cuadradas de precisión arbitraria:")
	two := new(big.Float).SetPrec(100).SetFloat64(2.0)
	sqrt2 := new(big.Float).Sqrt(two)
	fmt.Printf("√2 (100 bits) = %s\n", sqrt2.String())

	// Cálculo de π usando series
	fmt.Println("\n4. Cálculo de π usando la serie de Machin:")
	
	// π/4 = 4*arctan(1/5) - arctan(1/239)
	// π = 16*arctan(1/5) - 4*arctan(1/239)
	
	computePI := func(prec uint) *big.Float {
		// Implementación simplificada
		// En producción, usar bibliotecas especializadas
		return new(big.Float).SetPrec(prec).SetFloat64(3.141592653589793)
	}
	
	pi := computePI(100)
	fmt.Printf("π (aprox) = %s\n", pi.String())
}
```

### 32.9.3 Aplicación: Criptografía RSA (Simplificada)

```go
package main

import (
	"fmt"
	"math/big"
)

func gcd(a, b *big.Int) *big.Int {
	for b.Sign() != 0 {
		a, b = b, new(big.Int).Mod(a, b)
	}
	return a
}

func modInverse(a, m *big.Int) *big.Int {
	result := new(big.Int)
	result.ModInverse(a, m)
	return result
}

func main() {
	fmt.Println("=== RSA SIMPLIFICADO CON big.Int ===")

	// Seleccionar dos números primos grandes
	p := big.NewInt(61)
	q := big.NewInt(53)

	// Calcular n = p * q
	n := new(big.Int).Mul(p, q)
	fmt.Printf("n (modulo) = p * q = %v * %v = %v\n", p, q, n)

	// Calcular φ(n) = (p-1)(q-1)
	p_minus_1 := new(big.Int).Sub(p, big.NewInt(1))
	q_minus_1 := new(big.Int).Sub(q, big.NewInt(1))
	phi := new(big.Int).Mul(p_minus_1, q_minus_1)
	fmt.Printf("φ(n) = (p-1)(q-1) = %v\n", phi)

	// Seleccionar e tal que gcd(e, φ(n)) = 1
	e := big.NewInt(17)
	fmt.Printf("e = %v\n", e)
	fmt.Printf("gcd(e, φ(n)) = %v\n", gcd(e, phi))

	// Calcular d tal que e*d ≡ 1 (mod φ(n))
	d := modInverse(e, phi)
	fmt.Printf("d = %v\n", d)

	// Verificar: e*d mod φ(n) = 1
	verify := new(big.Int).Mul(e, d)
	verify.Mod(verify, phi)
	fmt.Printf("Verificación: e*d mod φ(n) = %v\n", verify)

	// Cifrado: C ≡ M^e (mod n)
	fmt.Println("\nCifrado/Descifrado:")
	message := big.NewInt(42)
	cipher := new(big.Int).Exp(message, e, n)
	fmt.Printf("Mensaje: %v\n", message)
	fmt.Printf("Cifrado: %v\n", cipher)

	// Descifrado: M ≡ C^d (mod n)
	decrypted := new(big.Int).Exp(cipher, d, n)
	fmt.Printf("Descifrado: %v\n", decrypted)
}
```

---

## 32.10 Numeric Stability - Estabilidad Numérica

### 32.10.1 Pérdida de Precisión

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== ESTABILIDAD NUMÉRICA ===")

	// PROBLEMA 1: Sustracción de números cercanos
	fmt.Println("\n1. CANCELACIÓN CATASTRÓFICA - Sustracción cercana:")
	x := 1.0
	y := 1.0 + 1e-15

	diff := y - x
	fmt.Printf("x = %.20f\n", x)
	fmt.Printf("y = %.20f\n", y)
	fmt.Printf("y - x = %.20e (esperado 1e-15)\n", diff)
	fmt.Printf("Error relativo = %.2f%%\n", math.Abs(diff-1e-15)/1e-15*100)

	// PROBLEMA 2: Suma de números grandes y pequeños
	fmt.Println("\n2. Orden de suma - Números grandes + pequeños:")
	grande := 1e20
	pequeno := 1.0

	// Forma incorrecta: pequeño se pierde
	result1 := grande + pequeno
	fmt.Printf("(1e20 + 1.0) - 1e20 = %.1f (esperado 1.0)\n", result1-grande)

	// Forma correcta: acumular primero los pequeños
	result2 := pequeno + (1e20 + (-1e20))
	fmt.Printf("Reorganizando: %.1f\n", result2)

	// PROBLEMA 3: Fórmula numérica inestable
	fmt.Println("\n3. Ecuación cuadrática - Fórmula inestable:")
	// x² + 2x + 1 = 0 → x = -1 (doble raíz)
	a, b, c := 1.0, 2.0, 1.0

	discriminant := b*b - 4*a*c
	fmt.Printf("Discriminante = %.20f\n", discriminant)

	// Fórmula estándar (puede ser inestable)
	x1_naive := (-b + math.Sqrt(discriminant)) / (2 * a)
	x2_naive := (-b - math.Sqrt(discriminant)) / (2 * a)
	fmt.Printf("Fórmula estándar: x1=%.20f, x2=%.20f\n", x1_naive, x2_naive)

	// Fórmula más estable
	if b > 0 {
		x1_stable := (-b - math.Sqrt(discriminant)) / (2 * a)
		x2_stable := c / (a * x1_stable)
		fmt.Printf("Fórmula estable: x1=%.20f, x2=%.20f\n", x1_stable, x2_stable)
	}

	// PROBLEMA 4: Comparación de flotantes
	fmt.Println("\n4. Comparación de números flotantes:")
	a1 := 0.1
	a2 := 0.1 + 0.2
	a3 := 0.3

	fmt.Printf("0.1 + 0.2 = %.20f\n", a2)
	fmt.Printf("0.3       = %.20f\n", a3)
	fmt.Printf("¿Son iguales? %v\n", a2 == a3)

	// Comparación correcta con tolerancia
	epsilon := 1e-10
	fmt.Printf("¿Diferencia < 1e-10? %v\n", math.Abs(a2-a3) < epsilon)

	// Comparación relativa (mejor para números grandes)
	fmt.Println("\n5. Comparación relativa de flotantes:")
	x = 1e-10
	y = 2e-10
	
	absError := math.Abs(y - x)
	relError := absError / math.Max(math.Abs(x), math.Abs(y))
	
	fmt.Printf("x = %.2e, y = %.2e\n", x, y)
	fmt.Printf("Error absoluto = %.2e\n", absError)
	fmt.Printf("Error relativo = %.2e\n", relError)
	fmt.Printf("¿Próximos (rel < 1e-8)? %v\n", relError < 1e-8)
}
```

### 32.10.2 Overflow y Underflow

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== OVERFLOW Y UNDERFLOW ===")

	fmt.Println("\n1. Overflow en float64:")
	fmt.Printf("MaxFloat64 = %.2e\n", math.MaxFloat64)
	
	// Cuando multiplicamos números muy grandes
	a := math.MaxFloat64 / 2
	b := 3.0
	result := a * b
	fmt.Printf("(MaxFloat64/2) * 3 = %v\n", result)
	fmt.Printf("¿Es Infinito? %v\n", math.IsInf(result, 1))

	fmt.Println("\n2. Underflow en float64:")
	fmt.Printf("SmallestNonzeroFloat64 = %.2e\n", math.SmallestNonzeroFloat64)
	
	// Cuando dividimos números muy pequeños
	tiny := math.SmallestNonzeroFloat64 * 2
	divisor := 1e20
	result = tiny / divisor
	fmt.Printf("(tiny * 2) / 1e20 = %.2e\n", result)
	fmt.Printf("¿Es cero? %v\n", result == 0)

	fmt.Println("\n3. Evitar overflow con logaritmos:")
	// Problema: e^x puede overflow para x grande
	x := 1000.0
	fmt.Printf("Exp(1000) = %v (overflow)\n", math.Exp(x))

	// Solución: usar logaritmos
	// Calcular: log(e^a * e^b) = a + b
	a = 500.0
	b = 600.0
	logSum := a + b
	fmt.Printf("ln(e^500 * e^600) = %.0f\n", logSum)

	fmt.Println("\n4. Underflow en exp:")
	// Problema: e^(-x) puede underflow para x grande
	x = 1000.0
	fmt.Printf("Exp(-1000) = %.2e (underflow a zero)\n", math.Exp(-x))

	// Solución: e^(-x) = 1/e^x, pero usar lógica especial
	// O usar math.Exp1p, math.Log1p para casos cercanos a 0
	x = -700.0
	fmt.Printf("Exp(-700) = %.2e\n", math.Exp(x))
}
```

### 32.10.3 Funciones de Comparación Segura

```go
package main

import (
	"fmt"
	"math"
)

// AlmostEqual compara dos flotantes con tolerancia relativa
func almostEqual(a, b, relTol, absTol float64) bool {
	diff := math.Abs(a - b)
	return diff <= math.Max(relTol*math.Max(math.Abs(a), math.Abs(b)), absTol)
}

// ULP (Unit in the Last Place) distancia
func ulpDistance(a, b float64) int64 {
	// Simplificado - en producción usar math.Nextafter
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.MaxInt64
	}
	
	bitsA := math.Float64bits(a)
	bitsB := math.Float64bits(b)
	
	if bitsA > bitsB {
		return int64(bitsA - bitsB)
	}
	return int64(bitsB - bitsA)
}

func main() {
	fmt.Println("=== COMPARACIONES SEGURAS DE FLOTANTES ===")

	fmt.Println("\n1. Tolerancia Absoluta:")
	a, b := 0.0, 1e-10
	absTol := 1e-9
	fmt.Printf("¿%.0e ≈ %.0e (absTol=1e-9)? %v\n", a, b, almostEqual(a, b, 0, absTol))

	fmt.Println("\n2. Tolerancia Relativa:")
	a, b = 1e10, 1e10 + 1
	relTol := 1e-6
	fmt.Printf("¿%.0e ≈ %.0e (relTol=1e-6)? %v\n", a, b, almostEqual(a, b, relTol, 0))

	fmt.Println("\n3. Combinación:")
	a, b = 0.1, 0.2+0.1
	equal := almostEqual(a, b, 1e-10, 1e-15)
	fmt.Printf("¿%.20f ≈ %.20f? %v\n", a, b, equal)

	fmt.Println("\n4. Distancia ULP:")
	tests := [][2]float64{
		{0.1, 0.1},
		{0.1, 0.1+1e-15},
		{1.0, 1.0 + 1e-16},
	}

	for _, test := range tests {
		dist := ulpDistance(test[0], test[1])
		fmt.Printf("ULP(%.20f, %.20f) = %d\n", test[0], test[1], dist)
	}
}
```

---

## 32.11 Buenas Prácticas y Patrones

### 32.11.1 Seeding y Reproducibilidad

```go
package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("=== SEEDING Y REPRODUCIBILIDAD ===")

	// PATRÓN 1: Seeding determinístico para testing
	fmt.Println("\n1. Seeding determinístico (testing):")
	seed := int64(42)
	rand.Seed(seed)
	
	fmt.Println("   Primera ejecución con seed=42:")
	for i := 0; i < 3; i++ {
		fmt.Printf("     %d\n", rand.Intn(1000))
	}

	rand.Seed(seed)
	fmt.Println("   Segunda ejecución con seed=42 (idéntico):")
	for i := 0; i < 3; i++ {
		fmt.Printf("     %d\n", rand.Intn(1000))
	}

	// PATRÓN 2: Seeding no determinístico para producción
	fmt.Println("\n2. Seeding no determinístico (producción):")
	
	// Opción A: Usar time (rápido pero no seguro)
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("   Con time.Now(): %d\n", rand.Intn(1000))

	// Opción B: Usar crypto/rand (seguro)
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	fmt.Printf("   Con crypto/rand: %d\n", rand.Intn(1000))

	// PATRÓN 3: Generadores locales para concurrencia
	fmt.Println("\n3. Generadores independientes:")
	
	generators := make([]*rand.Rand, 3)
	for i := 0; i < 3; i++ {
		source := rand.NewSource(int64(100 + i))
		generators[i] = rand.New(source)
	}

	for i, gen := range generators {
		fmt.Printf("   Gen %d: %d, %d, %d\n", i, gen.Intn(100), gen.Intn(100), gen.Intn(100))
	}

	// PATRÓN 4: Reproducibilidad en distribuciones
	fmt.Println("\n4. Reproducibilidad de distribuciones:")
	seed = 777
	rand.Seed(seed)
	
	fmt.Println("   Primeros 5 de N(0,1):")
	for i := 0; i < 5; i++ {
		fmt.Printf("     %.4f\n", rand.NormFloat64())
	}

	// Reset y repetir
	rand.Seed(seed)
	fmt.Println("   Repetición con seed=777:")
	for i := 0; i < 5; i++ {
		fmt.Printf("     %.4f\n", rand.NormFloat64())
	}

	// ANTIPATRÓN: No hacer esto
	fmt.Println("\n❌ ANTIPATRÓN - No hacer esto:")
	fmt.Println("   - No usar time.Now() para seeds criptográficos")
	fmt.Println("   - No reutilizar math/rand en múltiples goroutines sin mutex")
	fmt.Println("   - No usar math/rand para tokens/keys (usar crypto/rand)")

	// PATRÓN CORRECTO: No hacer esto
	fmt.Println("\n✓ PATRÓN CORRECTO:")
	fmt.Println("   - Testing: Usar seed fijo")
	fmt.Println("   - Simulaciones: Usar seed de time.Now()")
	fmt.Println("   - Concurrencia: Crear Rand independientes")
	fmt.Println("   - Criptografía: Usar crypto/rand siempre")
}
```

### 32.11.2 Testing con Números Aleatorios

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// Property-based test pattern
func testMathProperty(seed int64, iterations int, property func() bool) bool {
	rand.Seed(seed)
	for i := 0; i < iterations; i++ {
		if !property() {
			return false
		}
	}
	return true
}

func TestTrigonometricIdentity(t *testing.T) {
	// sin²(x) + cos²(x) = 1 siempre
	if !testMathProperty(42, 1000, func() bool {
		x := rand.Float64() * 2 * math.Pi
		result := math.Sin(x)*math.Sin(x) + math.Cos(x)*math.Cos(x)
		return math.Abs(result-1.0) < 1e-10
	}) {
		panic("Identidad trigonométrica fallida")
	}
	fmt.Println("✓ Identidad trigonométrica verificada")
}

func TestExponentialIdentity(t *testing.T) {
	// e^x * e^y = e^(x+y)
	if !testMathProperty(42, 1000, func() bool {
		x := (rand.Float64() - 0.5) * 10 // [-5, 5]
		y := (rand.Float64() - 0.5) * 10
		
		left := math.Exp(x) * math.Exp(y)
		right := math.Exp(x + y)
		
		return math.Abs(left-right) < math.Abs(left)*1e-10
	}) {
		panic("Identidad exponencial fallida")
	}
	fmt.Println("✓ Identidad exponencial verificada")
}

func main() {
	TestTrigonometricIdentity(nil)
	TestExponentialIdentity(nil)
}
```

### 32.11.3 Benchmarking de Funciones Matemáticas

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func benchmark(name string, iterations int, fn func()) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	elapsed := time.Since(start)
	
	nsPerOp := elapsed.Nanoseconds() / int64(iterations)
	fmt.Printf("%-20s: %6d ns/op (%v total)\n", name, nsPerOp, elapsed)
}

func main() {
	fmt.Println("=== BENCHMARK DE FUNCIONES MATEMÁTICAS ===\n")

	const iterations = 1000000
	rand.Seed(42)

	// Preparar datos
	values := make([]float64, 10000)
	for i := range values {
		values[i] = rand.Float64() * math.Pi * 2
	}
	idx := 0

	// Operaciones básicas
	benchmark("Sin", iterations, func() {
		math.Sin(values[idx%len(values)])
		idx++
	})

	idx = 0
	benchmark("Cos", iterations, func() {
		math.Cos(values[idx%len(values)])
		idx++
	})

	idx = 0
	benchmark("Sqrt", iterations, func() {
		math.Sqrt(values[idx%len(values)])
		idx++
	})

	idx = 0
	benchmark("Exp", iterations, func() {
		math.Exp(values[idx%len(values)] / 10)
		idx++
	})

	idx = 0
	benchmark("Log", iterations, func() {
		math.Log(values[idx%len(values)] + 1)
		idx++
	})

	idx = 0
	benchmark("Pow", iterations, func() {
		math.Pow(values[idx%len(values)], 2)
		idx++
	})

	// Funciones trigonométricas inversas
	idx = 0
	benchmark("Asin", iterations, func() {
		math.Asin(rand.Float64())
		idx++
	})

	idx = 0
	benchmark("Atan2", iterations, func() {
		math.Atan2(rand.Float64(), rand.Float64())
		idx++
	})

	// Random
	idx = 0
	benchmark("rand.Float64", iterations, func() {
		rand.Float64()
	})

	idx = 0
	benchmark("rand.Intn", iterations, func() {
		rand.Intn(1000000)
	})

	idx = 0
	benchmark("rand.NormFloat64", iterations, func() {
		rand.NormFloat64()
	})
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Calculadora Trigonométrica

Implementa una calculadora que convierta entre radianes y grados, y calcule seno, coseno y tangente:

```go
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Implementa aquí:
	// 1. Función para convertir grados a radianes
	// 2. Función para convertir radianes a grados
	// 3. Función que lea ángulos en grados
	// 4. Calcule e imprima sin(x), cos(x), tan(x)
	// 5. Lea desde stdin hasta que usuario ingrese "salir"
}

/* Ejemplo esperado:
Ingrese ángulo en grados (o "salir"): 45
45° = 0.7854 rad
sin(45°) = 0.7071
cos(45°) = 0.7071
tan(45°) = 1.0000

Ingrese ángulo en grados (o "salir"): 90
...
*/
```

---

### Ejercicio 2: Generador de Números Aleatorios Seeded

Crea un programa que genere secuencias reproducibles de números aleatorios:

```go
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// Implementa aquí:
	// 1. Función generarSecuencia(seed, cantidad) que retorna []int
	// 2. Generar 10 números aleatorios con seed=42
	// 3. Generar 10 números aleatorios con seed=123
	// 4. Generar nuevamente con seed=42 y verificar que sea idéntico
	// 5. Mostrar resultados en formato tabular
	
	// Esperado:
	// Seed 42: [...]
	// Seed 123: [...]
	// Seed 42 (repetido): [...] <- debe ser igual a Seed 42
}
```

---

### Ejercicio 3: Juego "Simón Dice" con Números Aleatorios

Implementa un juego interactivo donde el jugador debe adivinar una secuencia de números aleatorios:

```go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Implementa aquí:
	// 1. Generar secuencia creciente de números (1-10) aleatorios
	// 2. Mostrar cada número por 1 segundo
	// 3. Leer entrada del jugador
	// 4. Verificar si es correcto
	// 5. Aumentar dificultad (más números) si es correcto
	// 6. Mostrar "Game Over" con puntuación

	// Nivel 1: Memorizar 1 número
	// Nivel 2: Memorizar 2 números
	// Etc...
}
```

---

### Ejercicio 4: Método de Montecarlo para Calcular π

Estima el valor de π usando el método de Montecarlo (puntos aleatorios en cuadrado/círculo):

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	// Implementa aquí:
	// 1. Generar N puntos aleatorios en [0,1] x [0,1]
	// 2. Contar cuántos caen dentro del círculo (x²+y² ≤ 1)
	// 3. π ≈ 4 * (puntos_dentro / N)
	// 4. Mostrar convergencia a π con N = 1000, 10000, 100000, 1000000
	
	// Esperado:
	// N=1000:     π ≈ 3.1XX
	// N=10000:    π ≈ 3.14X
	// N=100000:   π ≈ 3.1415X
	// N=1000000:  π ≈ 3.14159X
	// Valor real: π ≈ 3.14159265
}
```

---

### Ejercicio 5: Generador de Tokens Seguros

Crea un sistema para generar tokens criptográficamente seguros para sesiones y contraseñas de recuperación:

```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	// Implementa aquí:
	// 1. Función generarToken(length int) string (base64)
	// 2. Función generarHexToken(length int) string
	// 3. Función generarOTP(digits int) string (código OTP)
	// 4. Función generarAPIKey() string (256 bits en hex)
	// 5. Generar 5 ejemplos de cada tipo
	
	// Salida esperada:
	// === Tokens de Sesión (Base64) ===
	// ...
	// === Tokens Hex (32 bytes) ===
	// ...
	// === Códigos OTP (6 dígitos) ===
	// ...
	// === API Keys (256-bit hex) ===
	// ...
}
```

---

## Resumen

Este capítulo ha cubierto exhaustivamente:

1. **math package**: Constantes, funciones trigonométricas, exponenciales y logarítmicas
2. **Rounding**: Ceil, Floor, Round, Trunc con aplicaciones prácticas
3. **Números complejos**: Operaciones, forma polar, rotaciones
4. **Random numbers**: math/rand para simulaciones, crypto/rand para seguridad
5. **Distribuciones**: Normal, Exponencial, Poisson con algoritmos implementados
6. **Big numbers**: Precisión arbitraria con big.Int y big.Float
7. **Estabilidad numérica**: Pérdida de precisión, overflow/underflow, comparaciones seguras
8. **Patrones**: Seeding, reproducibilidad, testing, benchmarking

Go proporciona herramientas robustas para análisis numérico, desde cálculos rápidos con float64 hasta aritmética de precisión arbitraria, combinadas con generadores aleatorios tanto rápidos como seguros.

---

## Referencias y Recursos Adicionales

- **Go math package**: https://golang.org/pkg/math/
- **Go math/rand**: https://golang.org/pkg/math/rand/
- **Go crypto/rand**: https://golang.org/pkg/crypto/rand/
- **IEEE 754 Floating Point**: https://en.wikipedia.org/wiki/IEEE_754
- **Numerical Stability**: https://en.wikipedia.org/wiki/Numerical_stability
- **PRNG vs CSPRNG**: https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator

---

**Fin del Capítulo 32**
