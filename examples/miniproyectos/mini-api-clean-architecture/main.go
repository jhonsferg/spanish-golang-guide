// Mini-proyecto de la Parte VI (Arquitectura y sistemas distribuidos):
// una API HTTP completa siguiendo Clean Architecture — dominio, casos de
// uso y repositorio (en memoria) totalmente desacoplados del transporte
// HTTP, expandiendo la idea del Capítulo 45 a un CRUD real. Ver
// Dockerfile en este directorio para el empaquetado de producción.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

func jsonReader(b []byte) io.Reader { return bytes.NewReader(b) }

// ============ Dominio ============

type Nota struct {
	ID       int    `json:"id"`
	Contenido string `json:"contenido"`
}

var ErrNotaNoEncontrada = errors.New("nota no encontrada")
var ErrContenidoVacio = errors.New("el contenido no puede estar vacío")

// ============ Puerto: lo que el caso de uso necesita ============

type RepositorioNotas interface {
	Crear(ctx context.Context, contenido string) (Nota, error)
	Listar(ctx context.Context) ([]Nota, error)
	Eliminar(ctx context.Context, id int) error
}

// ============ Casos de uso: reglas de negocio, sin saber qué es HTTP ============

type ServicioNotas struct {
	repo RepositorioNotas
}

func (s *ServicioNotas) CrearNota(ctx context.Context, contenido string) (Nota, error) {
	if contenido == "" {
		return Nota{}, ErrContenidoVacio
	}
	return s.repo.Crear(ctx, contenido)
}

func (s *ServicioNotas) ListarNotas(ctx context.Context) ([]Nota, error) {
	return s.repo.Listar(ctx)
}

func (s *ServicioNotas) EliminarNota(ctx context.Context, id int) error {
	return s.repo.Eliminar(ctx, id)
}

// ============ Adaptador: repositorio en memoria (intercambiable) ============

type repoEnMemoria struct {
	mu          sync.Mutex
	datos       map[int]Nota
	siguienteID int
}

func nuevoRepoEnMemoria() *repoEnMemoria {
	return &repoEnMemoria{datos: make(map[int]Nota), siguienteID: 1}
}

func (r *repoEnMemoria) Crear(ctx context.Context, contenido string) (Nota, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := Nota{ID: r.siguienteID, Contenido: contenido}
	r.datos[n.ID] = n
	r.siguienteID++
	return n, nil
}

func (r *repoEnMemoria) Listar(ctx context.Context) ([]Nota, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lista := make([]Nota, 0, len(r.datos))
	for _, n := range r.datos {
		lista = append(lista, n)
	}
	return lista, nil
}

func (r *repoEnMemoria) Eliminar(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.datos[id]; !ok {
		return ErrNotaNoEncontrada
	}
	delete(r.datos, id)
	return nil
}

// ============ Adaptador: transporte HTTP (solo traduce, no decide) ============

func nuevoMux(svc *ServicioNotas) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /notas", func(w http.ResponseWriter, r *http.Request) {
		var entrada struct {
			Contenido string `json:"contenido"`
		}
		json.NewDecoder(r.Body).Decode(&entrada)

		nota, err := svc.CrearNota(r.Context(), entrada.Contenido)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nota)
	})

	mux.HandleFunc("GET /notas", func(w http.ResponseWriter, r *http.Request) {
		notas, _ := svc.ListarNotas(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notas)
	})

	mux.HandleFunc("DELETE /notas/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		if err := svc.EliminarNota(r.Context(), id); err != nil {
			if errors.Is(err, ErrNotaNoEncontrada) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}

func main() {
	svc := &ServicioNotas{repo: nuevoRepoEnMemoria()}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: nuevoMux(svc)}
	go server.Serve(ln)
	defer server.Close()

	base := "http://" + ln.Addr().String()
	cliente := &http.Client{}

	crear := func(contenido string) *http.Response {
		body, _ := json.Marshal(map[string]string{"contenido": contenido})
		resp, err := cliente.Post(base+"/notas", "application/json", jsonReader(body))
		if err != nil {
			log.Fatal(err)
		}
		return resp
	}

	r1 := crear("Aplicar Clean Architecture al mini-proyecto")
	fmt.Println("POST /notas ->", r1.StatusCode)
	r1.Body.Close()

	rVacia := crear("")
	fmt.Println("POST /notas (vacía, rechazada por el caso de uso) ->", rVacia.StatusCode)
	rVacia.Body.Close()

	rLista, _ := cliente.Get(base + "/notas")
	var notas []Nota
	json.NewDecoder(rLista.Body).Decode(&notas)
	rLista.Body.Close()
	fmt.Printf("GET /notas -> %+v\n", notas)

	if len(notas) > 0 {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/notas/%d", base, notas[0].ID), nil)
		rDel, _ := cliente.Do(req)
		fmt.Println("DELETE /notas/{id} ->", rDel.StatusCode)
		rDel.Body.Close()
	}
}
