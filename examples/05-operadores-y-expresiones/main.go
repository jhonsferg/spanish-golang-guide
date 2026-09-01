// Ejemplo del Capítulo 5: operadores aritméticos, de comparación,
// lógicos y a nivel de bits.
package main

import "fmt"

func main() {
	a, b := 17, 5

	fmt.Println("--- Aritméticos ---")
	fmt.Printf("%d + %d = %d\n", a, b, a+b)
	fmt.Printf("%d - %d = %d\n", a, b, a-b)
	fmt.Printf("%d * %d = %d\n", a, b, a*b)
	fmt.Printf("%d / %d = %d (división entera)\n", a, b, a/b)
	fmt.Printf("%d %% %d = %d (módulo)\n", a, b, a%b)

	fmt.Println("--- Comparación ---")
	fmt.Printf("%d == %d -> %v\n", a, b, a == b)
	fmt.Printf("%d > %d -> %v\n", a, b, a > b)

	fmt.Println("--- Lógicos (con cortocircuito) ---")
	verdadero, falso := true, false
	fmt.Printf("%v && %v -> %v\n", verdadero, falso, verdadero && falso)
	fmt.Printf("%v || %v -> %v\n", verdadero, falso, verdadero || falso)
	fmt.Printf("!%v -> %v\n", verdadero, !verdadero)

	fmt.Println("--- Bit a bit ---")
	fmt.Printf("%d & %d  = %d (AND)\n", a, b, a&b)
	fmt.Printf("%d | %d  = %d (OR)\n", a, b, a|b)
	fmt.Printf("%d ^ %d  = %d (XOR)\n", a, b, a^b)
	fmt.Printf("%d &^ %d = %d (AND NOT)\n", a, b, a&^b)
	fmt.Printf("%d << 2 = %d (shift izquierda)\n", a, a<<2)
	fmt.Printf("%d >> 2 = %d (shift derecha)\n", a, a>>2)
}
