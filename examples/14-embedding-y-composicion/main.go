// Ejemplo del Capítulo 14: composición vía embedding de interfaces y
// structs, la alternativa de Go a la herencia clásica.
package main

import "fmt"

type Escritor interface {
	Escribir(msg string)
}

type Lector interface {
	Leer() string
}

// LectorEscritor se compone embebiendo otras dos interfaces.
type LectorEscritor interface {
	Lector
	Escritor
}

type Buffer struct {
	contenido string
}

func (b *Buffer) Escribir(msg string) { b.contenido += msg }
func (b *Buffer) Leer() string        { return b.contenido }

// Auditado embebe *Buffer: hereda Leer/Escribir "gratis" y además
// puede sobrescribir el comportamiento agregando logging.
type Auditado struct {
	*Buffer
	operaciones int
}

func (a *Auditado) Escribir(msg string) {
	a.operaciones++
	a.Buffer.Escribir(msg)
}

func usar(rw LectorEscritor) {
	rw.Escribir("hola ")
	rw.Escribir("mundo")
	fmt.Println("contenido:", rw.Leer())
}

func main() {
	b := &Buffer{}
	usar(b)

	a := &Auditado{Buffer: &Buffer{}}
	usar(a)
	fmt.Println("operaciones auditadas:", a.operaciones)
}
