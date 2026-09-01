// Ejemplo del Capítulo 23: select para multiplexar channels, con
// timeout y un caso default no bloqueante.
package main

import (
	"fmt"
	"time"
)

func main() {
	rapido := make(chan string)
	lento := make(chan string)

	go func() {
		time.Sleep(20 * time.Millisecond)
		rapido <- "respuesta rápida"
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		lento <- "respuesta lenta"
	}()

	// select espera a que CUALQUIERA de los channels esté listo.
	for i := 0; i < 2; i++ {
		select {
		case msg := <-rapido:
			fmt.Println("recibido:", msg)
		case msg := <-lento:
			fmt.Println("recibido:", msg)
		case <-time.After(500 * time.Millisecond):
			fmt.Println("timeout esperando respuesta")
		}
	}

	// select con default: no bloquea si ningún channel está listo.
	senal := make(chan struct{})
	select {
	case <-senal:
		fmt.Println("señal recibida")
	default:
		fmt.Println("no hay señal todavía, seguimos sin bloquear")
	}
}
