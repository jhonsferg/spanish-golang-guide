// Ejemplo del Capítulo 19: defer, orden LIFO, evaluación de argumentos
// en el momento del defer, y el patrón de cierre de recursos.
package main

import "fmt"

func ordenLIFO() {
	fmt.Println("--- orden LIFO ---")
	for i := 1; i <= 3; i++ {
		defer fmt.Println("defer registrado con i =", i)
	}
	fmt.Println("fin de la función (los defer corren después, en orden inverso)")
}

type Recurso struct {
	Nombre string
}

func (r *Recurso) Cerrar() {
	fmt.Printf("cerrando recurso %q\n", r.Nombre)
}

func abrirRecurso(nombre string) *Recurso {
	fmt.Printf("abriendo recurso %q\n", nombre)
	return &Recurso{Nombre: nombre}
}

func usarRecursos() {
	fmt.Println("--- patrón de cierre garantizado ---")
	r1 := abrirRecurso("conexión-db")
	defer r1.Cerrar()

	r2 := abrirRecurso("archivo-log")
	defer r2.Cerrar()

	fmt.Println("usando ambos recursos...")
	// aunque hubiera un panic aquí, los defer igual se ejecutarían.
}

func main() {
	ordenLIFO()
	usarRecursos()
}
