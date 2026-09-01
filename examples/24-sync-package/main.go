// Ejemplo del Capítulo 24: primitivas de sync — Mutex para proteger
// estado compartido, y Once para inicialización perezosa.
package main

import (
	"fmt"
	"sync"
)

// ContadorSeguro protege `valor` con un Mutex: sin él, incrementar desde
// múltiples goroutines a la vez es una data race (detectable con -race).
type ContadorSeguro struct {
	mu    sync.Mutex
	valor int
}

func (c *ContadorSeguro) Incrementar() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.valor++
}

func (c *ContadorSeguro) Valor() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.valor
}

var una sync.Once

func inicializarUnaVez() {
	una.Do(func() {
		fmt.Println("esto se imprime UNA sola vez, sin importar cuántas goroutines lo llamen")
	})
}

func main() {
	contador := &ContadorSeguro{}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Go(contador.Incrementar)
	}
	wg.Wait()
	fmt.Println("valor final tras 100 incrementos concurrentes:", contador.Valor())

	var wg2 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg2.Go(inicializarUnaVez)
	}
	wg2.Wait()
}
