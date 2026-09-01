// Package gormrepo implementa el repositorio de Producto con GORM. Ver
// sqlcrepo/ y entrepo/ para la MISMA operación con otras dos formas de
// acceder a datos en Go.
package gormrepo

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Producto struct {
	gorm.Model
	Nombre string
	Precio float64
}

type Repositorio struct {
	db *gorm.DB
}

func Nuevo() (*Repositorio, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Producto{}); err != nil {
		return nil, err
	}
	return &Repositorio{db: db}, nil
}

func (r *Repositorio) Crear(nombre string, precio float64) (Producto, error) {
	p := Producto{Nombre: nombre, Precio: precio}
	err := r.db.Create(&p).Error
	return p, err
}

func (r *Repositorio) Listar() ([]Producto, error) {
	var productos []Producto
	err := r.db.Order("precio desc").Find(&productos).Error
	return productos, err
}
