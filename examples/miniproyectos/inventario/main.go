// Mini-proyecto de la Parte II (Tipos y composición): un sistema de
// inventario con productos físicos y digitales que comparten
// comportamiento vía interfaces y embedding — sin ningún framework,
// solo lo cubierto en los capítulos 11-19.
package main

import (
	"errors"
	"fmt"
)

// --- Interfaces: el contrato que cualquier producto debe cumplir ---

type Vendible interface {
	Precio() float64
	Descripcion() string
}

type ConStock interface {
	Vendible
	Reservar(cantidad int) error
	StockDisponible() int
}

// --- Tipo base, embebido en los productos concretos ---

type productoBase struct {
	Nombre       string
	PrecioBase   float64
}

func (p productoBase) Precio() float64 { return p.PrecioBase }

// --- Producto físico: tiene stock limitado ---

type ProductoFisico struct {
	productoBase
	Stock int
}

func (p ProductoFisico) Descripcion() string {
	return fmt.Sprintf("%s (físico) - $%.2f", p.Nombre, p.Precio())
}

var ErrStockInsuficiente = errors.New("stock insuficiente")

func (p *ProductoFisico) Reservar(cantidad int) error {
	if cantidad > p.Stock {
		return fmt.Errorf("reservando %q: %w (disponible: %d)", p.Nombre, ErrStockInsuficiente, p.Stock)
	}
	p.Stock -= cantidad
	return nil
}

func (p ProductoFisico) StockDisponible() int { return p.Stock }

// --- Producto digital: stock "infinito", no implementa ConStock ---

type ProductoDigital struct {
	productoBase
	TamanoMB int
}

func (p ProductoDigital) Descripcion() string {
	return fmt.Sprintf("%s (digital, %dMB) - $%.2f", p.Nombre, p.TamanoMB, p.Precio())
}

// --- Inventario: opera sobre las interfaces, no sobre tipos concretos ---

type Inventario struct {
	productos []Vendible
}

func (inv *Inventario) Agregar(p Vendible) {
	inv.productos = append(inv.productos, p)
}

func (inv *Inventario) ValorTotal() float64 {
	total := 0.0
	for _, p := range inv.productos {
		total += p.Precio()
	}
	return total
}

func (inv *Inventario) Listar() {
	for _, p := range inv.productos {
		fmt.Println(" -", p.Descripcion())
	}
}

// procesarVenta solo puede vender productos que implementen ConStock:
// el type assertion es la forma idiomática de pedir "más capacidad".
func procesarVenta(inv *Inventario, nombre string, cantidad int) error {
	for _, p := range inv.productos {
		conStock, esFisico := p.(ConStock)
		if !esFisico {
			continue
		}
		if desc := p.Descripcion(); len(desc) > 0 && contieneNombre(desc, nombre) {
			return conStock.Reservar(cantidad)
		}
	}
	return fmt.Errorf("producto %q no encontrado o no tiene stock gestionable", nombre)
}

func contieneNombre(descripcion, nombre string) bool {
	return len(descripcion) >= len(nombre) && descripcion[:len(nombre)] == nombre
}

func main() {
	inv := &Inventario{}

	teclado := &ProductoFisico{productoBase: productoBase{Nombre: "Teclado", PrecioBase: 89.90}, Stock: 5}
	inv.Agregar(teclado)
	inv.Agregar(&ProductoFisico{productoBase: productoBase{Nombre: "Mouse", PrecioBase: 39.90}, Stock: 2})
	inv.Agregar(ProductoDigital{productoBase: productoBase{Nombre: "Curso de Go", PrecioBase: 29.90}, TamanoMB: 850})

	fmt.Println("--- inventario inicial ---")
	inv.Listar()
	fmt.Printf("valor total: $%.2f\n", inv.ValorTotal())

	fmt.Println("--- procesando venta de 3 teclados ---")
	if err := procesarVenta(inv, "Teclado", 3); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("venta exitosa, stock restante:", teclado.StockDisponible())
	}

	fmt.Println("--- procesando venta imposible de 10 teclados ---")
	if err := procesarVenta(inv, "Teclado", 10); err != nil {
		fmt.Println("error esperado:", err)
	}
}
