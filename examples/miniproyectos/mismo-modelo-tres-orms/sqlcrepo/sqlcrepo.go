// Package sqlcrepo implementa el repositorio de Producto en el estilo
// que `sqlc generate` produciría a partir de SQL crudo: sin ORM, con
// Scan explícito. Ver gormrepo/ y entrepo/ para la misma operación con
// otras dos formas de acceder a datos.
package sqlcrepo

import (
	"context"
	"database/sql"

	// Se usa el mismo driver subyacente que gormrepo (glebarez/sqlite lo
	// envuelve) para que ambos registren el driver "sqlite" una sola vez
	// al correr en el mismo binario — ver mismo-modelo-tres-orms/main.go.
	_ "github.com/glebarez/go-sqlite"
)

type Producto struct {
	ID     int64
	Nombre string
	Precio float64
}

type Repositorio struct {
	db *sql.DB
}

func Nuevo() (*Repositorio, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE productos (
		id     INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		precio REAL NOT NULL
	)`); err != nil {
		return nil, err
	}
	return &Repositorio{db: db}, nil
}

func (r *Repositorio) Crear(ctx context.Context, nombre string, precio float64) (Producto, error) {
	row := r.db.QueryRowContext(ctx,
		"INSERT INTO productos (nombre, precio) VALUES (?, ?) RETURNING id, nombre, precio",
		nombre, precio)
	var p Producto
	err := row.Scan(&p.ID, &p.Nombre, &p.Precio)
	return p, err
}

func (r *Repositorio) Listar(ctx context.Context) ([]Producto, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, nombre, precio FROM productos ORDER BY precio DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Precio); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, rows.Err()
}
