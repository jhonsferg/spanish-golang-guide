// Ejemplo del Capítulo 21: lanzar goroutines y esperar a que terminen
// con sync.WaitGroup (usando WaitGroup.Go, disponible desde Go 1.25).
package main

import (
	"fmt"
	"sync"
	"time"
)

func trabajo(id int, resultados chan<- string) {
	time.Sleep(time.Duration(id) * 10 * time.Millisecond) // simula trabajo
	resultados <- fmt.Sprintf("goroutine %d terminó", id)
}

func main() {
	resultados := make(chan string, 5)

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		id := i // cada iteración ya tiene su propia variable desde Go 1.22
		wg.Go(func() {
			trabajo(id, resultados)
		})
	}

	// Goroutine auxiliar que cierra el channel cuando todo el trabajo terminó,
	// para poder usar `range` sobre `resultados` sin bloquear para siempre.
	go func() {
		wg.Wait()
		close(resultados)
	}()

	for r := range resultados {
		fmt.Println(r)
	}

	fmt.Println("todas las goroutines terminaron")
}
