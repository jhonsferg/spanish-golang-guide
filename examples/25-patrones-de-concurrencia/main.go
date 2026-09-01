// Ejemplo del Capítulo 25: patrón worker pool — un número fijo de
// goroutines consumiendo tareas de un channel compartido.
package main

import (
	"fmt"
	"sync"
	"time"
)

type Tarea struct {
	ID int
}

type Resultado struct {
	TareaID   int
	WorkerID  int
	Respuesta string
}

func worker(id int, tareas <-chan Tarea, resultados chan<- Resultado, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tareas {
		time.Sleep(15 * time.Millisecond) // simula trabajo real
		resultados <- Resultado{
			TareaID:   t.ID,
			WorkerID:  id,
			Respuesta: fmt.Sprintf("tarea %d procesada", t.ID),
		}
	}
}

func main() {
	const numWorkers = 3
	const numTareas = 9

	tareas := make(chan Tarea, numTareas)
	resultados := make(chan Resultado, numTareas)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, tareas, resultados, &wg)
	}

	for i := 1; i <= numTareas; i++ {
		tareas <- Tarea{ID: i}
	}
	close(tareas) // no habrá más tareas: los workers terminan su range y salen

	go func() {
		wg.Wait()
		close(resultados)
	}()

	conteoPorWorker := make(map[int]int)
	for r := range resultados {
		conteoPorWorker[r.WorkerID]++
	}

	fmt.Println("tareas procesadas por worker:")
	for w := 1; w <= numWorkers; w++ {
		fmt.Printf("  worker %d: %d tareas\n", w, conteoPorWorker[w])
	}
}
