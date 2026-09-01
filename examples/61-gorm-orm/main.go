// Ejemplo del Capítulo 61: GORM con un driver SQLite puro-Go (sin cgo) -
// modelos con asociaciones, AutoMigrate, Preload y soft delete, las
// piezas que más se usan de este ORM en el día a día.
package main

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Autor struct {
	gorm.Model
	Nombre string
	Libros []Libro // relación hasMany
}

type Libro struct {
	gorm.Model
	Titulo  string
	AutorID uint
}

func main() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&Autor{}, &Libro{}); err != nil {
		log.Fatal(err)
	}

	autor := Autor{
		Nombre: "Jorge Luis Borges",
		Libros: []Libro{
			{Titulo: "Ficciones"},
			{Titulo: "El Aleph"},
		},
	}
	// Create() inserta el autor Y sus libros asociados en una sola llamada.
	if err := db.Create(&autor).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("autor creado con ID=%d y %d libros\n", autor.ID, len(autor.Libros))

	// Preload carga la relación (evita el problema N+1 si se usa bien).
	var autorConLibros Autor
	if err := db.Preload("Libros").First(&autorConLibros, autor.ID).Error; err != nil {
		log.Fatal(err)
	}
	for _, l := range autorConLibros.Libros {
		fmt.Printf("  libro: %s\n", l.Titulo)
	}

	// Update parcial: solo cambia el campo indicado.
	db.Model(&autor).Update("Nombre", "Jorge Luis Borges (actualizado)")

	// Delete con gorm.Model es SOFT delete: pone deleted_at, no borra la fila.
	db.Delete(&Libro{}, "titulo = ?", "El Aleph")

	var librosActivos int64
	db.Model(&Libro{}).Where("autor_id = ?", autor.ID).Count(&librosActivos)
	fmt.Println("libros activos tras soft delete:", librosActivos)

	var totalConBorrados int64
	db.Unscoped().Model(&Libro{}).Where("autor_id = ?", autor.ID).Count(&totalConBorrados)
	fmt.Println("total incluyendo soft-deleted (Unscoped):", totalConBorrados)
}
