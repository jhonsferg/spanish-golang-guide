// Ejemplo del Capítulo 12: métodos con receiver de valor vs. de puntero,
// y por qué la elección importa.
package main

import "fmt"

type Contador struct {
	valor int
}

// Receiver de valor: opera sobre una COPIA; no muta el original.
func (c Contador) Duplicado() int {
	return c.valor * 2
}

// Receiver de puntero: opera sobre el original; puede mutarlo.
func (c *Contador) Incrementar() {
	c.valor++
}

type Rectangulo struct {
	Ancho, Alto float64
}

func (r Rectangulo) Area() float64 {
	return r.Ancho * r.Alto
}

func (r Rectangulo) Perimetro() float64 {
	return 2 * (r.Ancho + r.Alto)
}

func main() {
	c := Contador{valor: 5}
	fmt.Println("duplicado:", c.Duplicado())

	c.Incrementar() // Go toma la dirección automáticamente: (&c).Incrementar()
	c.Incrementar()
	fmt.Println("valor tras incrementar dos veces:", c.valor)

	r := Rectangulo{Ancho: 4, Alto: 3}
	fmt.Printf("área=%.1f perímetro=%.1f\n", r.Area(), r.Perimetro())
}
