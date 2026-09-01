// Ejemplo del Capítulo 22: channels con y sin buffer, direccionalidad
// (chan<-, <-chan) y cierre de channels.
package main

import "fmt"

// generar solo puede ENVIAR al channel (chan<- int).
func generar(n int) <-chan int {
	salida := make(chan int)
	go func() {
		defer close(salida) // avisa a los receptores que no hay más valores
		for i := 1; i <= n; i++ {
			salida <- i
		}
	}()
	return salida
}

// cuadrados solo puede RECIBIR de entrada y ENVIAR a su salida.
func cuadrados(entrada <-chan int) <-chan int {
	salida := make(chan int)
	go func() {
		defer close(salida)
		for n := range entrada {
			salida <- n * n
		}
	}()
	return salida
}

func main() {
	// Channel sin buffer: el emisor bloquea hasta que alguien recibe.
	sinBuffer := make(chan string)
	go func() { sinBuffer <- "mensaje síncrono" }()
	fmt.Println(<-sinBuffer)

	// Channel con buffer: el emisor no bloquea hasta llenar el buffer.
	conBuffer := make(chan int, 2)
	conBuffer <- 1
	conBuffer <- 2
	fmt.Println("buffer lleno, len:", len(conBuffer), "cap:", cap(conBuffer))
	fmt.Println(<-conBuffer, <-conBuffer)

	// Pipeline: generar -> cuadrados, encadenando channels.
	for n := range cuadrados(generar(5)) {
		fmt.Println("cuadrado:", n)
	}

	// Leer de un channel cerrado retorna el valor cero + ok=false.
	cerrado := make(chan int)
	close(cerrado)
	v, ok := <-cerrado
	fmt.Println("de channel cerrado:", v, ok)
}
