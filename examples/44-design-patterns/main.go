// Ejemplo del Capítulo 44: tres patrones de diseño clásicos, expresados
// de forma idiomática en Go - con funciones e interfaces pequeñas, sin
// imitar jerarquías de clases de otros lenguajes.
package main

import "fmt"

// --- Strategy: el algoritmo es un valor, no una jerarquía de clases ---

type EstrategiaDescuento func(precio float64) float64

func sinDescuento(precio float64) float64      { return precio }
func descuentoVIP(precio float64) float64      { return precio * 0.8 }
func descuentoLiquidacion(precio float64) float64 { return precio * 0.5 }

func aplicarDescuento(precio float64, estrategia EstrategiaDescuento) float64 {
	return estrategia(precio)
}

// --- Decorator: envolver un http.Handler-like con funcionalidad extra ---

type Notificador interface {
	Enviar(mensaje string) string
}

type NotificadorBase struct{}

func (NotificadorBase) Enviar(mensaje string) string { return "enviado: " + mensaje }

type ConTimestamp struct {
	Notificador
}

func (c ConTimestamp) Enviar(mensaje string) string {
	return "[2026-09-01] " + c.Notificador.Enviar(mensaje)
}

type ConMayusculas struct {
	Notificador
}

func (c ConMayusculas) Enviar(mensaje string) string {
	base := c.Notificador.Enviar(mensaje)
	resultado := make([]byte, len(base))
	for i := range base {
		b := base[i]
		if b >= 'a' && b <= 'z' {
			b -= 32
		}
		resultado[i] = b
	}
	return string(resultado)
}

// --- Observer: canal de eventos simple, sin librería externa ---

type Evento struct{ Tipo, Detalle string }

type Observador func(Evento)

type Publicador struct {
	observadores []Observador
}

func (p *Publicador) Suscribir(o Observador) {
	p.observadores = append(p.observadores, o)
}

func (p *Publicador) Publicar(e Evento) {
	for _, o := range p.observadores {
		o(e)
	}
}

func main() {
	fmt.Println("--- Strategy ---")
	precio := 100.0
	for nombre, estrategia := range map[string]EstrategiaDescuento{
		"sin descuento": sinDescuento,
		"VIP":           descuentoVIP,
		"liquidación":   descuentoLiquidacion,
	} {
		fmt.Printf("%s: %.2f\n", nombre, aplicarDescuento(precio, estrategia))
	}

	fmt.Println("--- Decorator ---")
	var n Notificador = NotificadorBase{}
	n = ConTimestamp{n}
	n = ConMayusculas{n}
	fmt.Println(n.Enviar("build completado"))

	fmt.Println("--- Observer ---")
	pub := &Publicador{}
	pub.Suscribir(func(e Evento) { fmt.Printf("log: %s - %s\n", e.Tipo, e.Detalle) })
	pub.Suscribir(func(e Evento) { fmt.Printf("métrica incrementada para: %s\n", e.Tipo) })
	pub.Publicar(Evento{Tipo: "pedido_creado", Detalle: "pedido #123"})
}
