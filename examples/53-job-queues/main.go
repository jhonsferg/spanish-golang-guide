// Ejemplo del Capítulo 53: una cola de trabajos en memoria con
// reintentos automáticos — el patrón detrás de sistemas como Sidekiq o
// Asynq, simplificado a lo esencial con channels.
package main

import (
	"errors"
	"fmt"
	"sync"
)

type Trabajo struct {
	ID         int
	Intento    int
	MaxIntentos int
}

type Resultado struct {
	TrabajoID int
	Exito     bool
	Intentos  int
}

// procesar simula un trabajo que falla las dos primeras veces y luego
// tiene éxito — típico de llamadas a servicios externos poco confiables.
func procesar(t Trabajo) error {
	if t.Intento < 3 {
		return errors.New("servicio externo no disponible (simulado)")
	}
	return nil
}

func worker(id int, cola chan Trabajo, resultados chan<- Resultado, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range cola {
		t.Intento++
		err := procesar(t)
		if err != nil {
			if t.Intento < t.MaxIntentos {
				// Reencolar: en un sistema real esto llevaría backoff
				// exponencial y un límite de reintentos persistido.
				cola <- t
				continue
			}
			resultados <- Resultado{TrabajoID: t.ID, Exito: false, Intentos: t.Intento}
			continue
		}
		resultados <- Resultado{TrabajoID: t.ID, Exito: true, Intentos: t.Intento}
	}
}

func main() {
	const numTrabajos = 4
	const numWorkers = 2

	cola := make(chan Trabajo, numTrabajos*5) // espacio de sobra para reintentos
	resultados := make(chan Resultado, numTrabajos)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, cola, resultados, &wg)
	}

	for i := 1; i <= numTrabajos; i++ {
		cola <- Trabajo{ID: i, MaxIntentos: 3}
	}

	recibidos := 0
	for r := range resultados {
		fmt.Printf("trabajo #%d: éxito=%v tras %d intento(s)\n", r.TrabajoID, r.Exito, r.Intentos)
		recibidos++
		if recibidos == numTrabajos {
			close(cola)
			break
		}
	}

	wg.Wait()
}
