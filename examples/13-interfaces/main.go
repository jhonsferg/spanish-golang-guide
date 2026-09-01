// Ejemplo del Capítulo 13: interfaces satisfechas implícitamente y
// type switches.
package main

import (
	"fmt"
	"math"
)

type Figura interface {
	Area() float64
}

type Circulo struct{ Radio float64 }

func (c Circulo) Area() float64 { return math.Pi * c.Radio * c.Radio }

type Cuadrado struct{ Lado float64 }

func (s Cuadrado) Area() float64 { return s.Lado * s.Lado }

// describir acepta cualquier tipo que implemente Figura - Circulo y
// Cuadrado nunca declaran "implements Figura", simplemente tienen el método.
func describir(f Figura) string {
	switch v := f.(type) {
	case Circulo:
		return fmt.Sprintf("círculo de radio %.1f, área %.2f", v.Radio, v.Area())
	case Cuadrado:
		return fmt.Sprintf("cuadrado de lado %.1f, área %.2f", v.Lado, v.Area())
	default:
		return fmt.Sprintf("figura desconocida con área %.2f", v.Area())
	}
}

func main() {
	figuras := []Figura{
		Circulo{Radio: 2},
		Cuadrado{Lado: 3},
	}

	for _, f := range figuras {
		fmt.Println(describir(f))
	}

	// La interfaz vacía `any` acepta cualquier valor.
	var cualquiera any = 42
	if n, ok := cualquiera.(int); ok {
		fmt.Println("es un int:", n)
	}
}
