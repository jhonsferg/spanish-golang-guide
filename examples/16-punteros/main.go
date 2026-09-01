// Ejemplo del Capítulo 16: punteros, por qué existen y cómo Go los usa
// sin aritmética de punteros (a diferencia de C).
package main

import "fmt"

type Saldo struct {
	Monto int
}

// retirar recibe un PUNTERO a Saldo: sin él, la modificación sería
// invisible fuera de la función (Go pasa todo por valor).
func retirar(s *Saldo, cantidad int) error {
	if cantidad > s.Monto {
		return fmt.Errorf("fondos insuficientes: tienes %d, pediste %d", s.Monto, cantidad)
	}
	s.Monto -= cantidad
	return nil
}

func main() {
	x := 10
	p := &x // p apunta a la dirección de x
	fmt.Println("valor de x:", x, "| dirección:", p, "| valor vía p:", *p)

	*p = 20 // modifica x indirectamente
	fmt.Println("x tras *p = 20:", x)

	cuenta := &Saldo{Monto: 100}
	if err := retirar(cuenta, 30); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("saldo tras retiro:", cuenta.Monto)

	if err := retirar(cuenta, 1000); err != nil {
		fmt.Println("error esperado:", err)
	}

	// new() reserva memoria y retorna un puntero al valor cero del tipo.
	entero := new(int)
	*entero = 42
	fmt.Println("puntero de new():", *entero)

	// Un puntero nil es válido de declarar, pero desreferenciarlo hace panic.
	var nulo *Saldo
	fmt.Println("puntero nil:", nulo == nil)
}
