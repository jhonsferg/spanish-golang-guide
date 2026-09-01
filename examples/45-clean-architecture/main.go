// Ejemplo del Capítulo 45: Clean Architecture en miniatura - el caso de
// uso depende de una INTERFAZ de repositorio, no de una implementación
// concreta. Cambiar de memoria a una base de datos real no toca el
// dominio ni el caso de uso.
package main

import (
	"errors"
	"fmt"
)

// --- Dominio: sin dependencias externas, sin saber qué es HTTP o SQL ---

type Usuario struct {
	ID     int
	Nombre string
	Activo bool
}

var ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")

// --- Puerto (interfaz) que el caso de uso necesita ---

type RepositorioUsuarios interface {
	BuscarPorID(id int) (Usuario, error)
	Guardar(u Usuario) error
}

// --- Caso de uso: orquesta la lógica de negocio usando el puerto ---

type DesactivarUsuario struct {
	Repo RepositorioUsuarios
}

func (uc DesactivarUsuario) Ejecutar(id int) error {
	usuario, err := uc.Repo.BuscarPorID(id)
	if err != nil {
		return fmt.Errorf("desactivando usuario %d: %w", id, err)
	}
	usuario.Activo = false
	return uc.Repo.Guardar(usuario)
}

// --- Adaptador: implementación concreta del puerto, en memoria ---
// En producción sería un adaptador de PostgreSQL/GORM/sqlc - el caso de
// uso de arriba no cambiaría ni una línea.

type RepoUsuariosEnMemoria struct {
	datos map[int]Usuario
}

func NuevoRepoEnMemoria() *RepoUsuariosEnMemoria {
	return &RepoUsuariosEnMemoria{datos: map[int]Usuario{
		1: {ID: 1, Nombre: "Ada", Activo: true},
	}}
}

func (r *RepoUsuariosEnMemoria) BuscarPorID(id int) (Usuario, error) {
	u, ok := r.datos[id]
	if !ok {
		return Usuario{}, ErrUsuarioNoEncontrado
	}
	return u, nil
}

func (r *RepoUsuariosEnMemoria) Guardar(u Usuario) error {
	r.datos[u.ID] = u
	return nil
}

func main() {
	repo := NuevoRepoEnMemoria()
	caso := DesactivarUsuario{Repo: repo}

	if err := caso.Ejecutar(1); err != nil {
		fmt.Println("error:", err)
	}
	usuario, _ := repo.BuscarPorID(1)
	fmt.Printf("usuario tras el caso de uso: %+v\n", usuario)

	if err := caso.Ejecutar(999); err != nil {
		fmt.Println("error esperado:", err)
	}
}
