// Mini-proyecto de la Parte IX (Bases de datos y ORMs): el MISMO modelo
// (Producto: nombre + precio) creado y consultado con GORM, con el
// estilo de sqlc y con el estilo de Ent — para comparar directamente
// cuánto código y qué forma toma cada enfoque resolviendo lo mismo.
package main

import (
	"context"
	"fmt"

	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/mismo-modelo-tres-orms/entrepo"
	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/mismo-modelo-tres-orms/gormrepo"
	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/mismo-modelo-tres-orms/sqlcrepo"
)

var productosDemo = []struct {
	nombre string
	precio float64
}{
	{"Teclado mecánico", 89.90},
	{"Monitor 27\"", 249.90},
	{"Mouse ergonómico", 39.90},
}

func main() {
	fmt.Println("--- GORM ---")
	repoGorm, err := gormrepo.Nuevo()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, p := range productosDemo {
		if _, err := repoGorm.Crear(p.nombre, p.precio); err != nil {
			fmt.Println("error:", err)
			return
		}
	}
	listaGorm, _ := repoGorm.Listar()
	for _, p := range listaGorm {
		fmt.Printf("  %-20s $%.2f\n", p.Nombre, p.Precio)
	}

	fmt.Println("--- estilo sqlc ---")
	repoSqlc, err := sqlcrepo.Nuevo()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	ctx := context.Background()
	for _, p := range productosDemo {
		if _, err := repoSqlc.Crear(ctx, p.nombre, p.precio); err != nil {
			fmt.Println("error:", err)
			return
		}
	}
	listaSqlc, _ := repoSqlc.Listar(ctx)
	for _, p := range listaSqlc {
		fmt.Printf("  %-20s $%.2f\n", p.Nombre, p.Precio)
	}

	fmt.Println("--- estilo Ent ---")
	clienteEnt := entrepo.Nuevo()
	for _, p := range productosDemo {
		clienteEnt.CreateProducto().SetNombre(p.nombre).SetPrecio(p.precio).Save()
	}
	for _, p := range clienteEnt.QueryProductos() {
		fmt.Printf("  %-20s $%.2f\n", p.Nombre, p.Precio)
	}
}
