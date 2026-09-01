// Ejemplo del Capítulo 60: Chi - un router ligero y 100% compatible con
// net/http (cada handler sigue siendo un http.HandlerFunc normal), con
// su fuerte énfasis en composición de middleware por sub-router.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func jsonBody(s string) io.Reader { return strings.NewReader(s) }

type Autor struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

func responderJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func nuevoRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer, middleware.RequestID)

	autores := map[int]Autor{1: {1, "Rob Pike"}}
	siguienteID := 2

	// Un sub-router con su PROPIO middleware, montado bajo /api -
	// la composición que hace fuerte a Chi frente a un router plano.
	api := chi.NewRouter()
	api.Use(middleware.AllowContentType("application/json"))

	api.Get("/autores", func(w http.ResponseWriter, r *http.Request) {
		lista := make([]Autor, 0, len(autores))
		for _, a := range autores {
			lista = append(lista, a)
		}
		responderJSON(w, http.StatusOK, lista)
	})

	api.Post("/autores", func(w http.ResponseWriter, r *http.Request) {
		var a Autor
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
			return
		}
		a.ID = siguienteID
		autores[a.ID] = a
		siguienteID++
		responderJSON(w, http.StatusCreated, a)
	})

	api.Get("/autores/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
			return
		}
		a, ok := autores[id]
		if !ok {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "no encontrado"})
			return
		}
		responderJSON(w, http.StatusOK, a)
	})

	r.Mount("/api", api)
	return r
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := &http.Server{Handler: nuevoRouter()}
	go server.Serve(ln)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	base := "http://" + ln.Addr().String()

	// GET sin Content-Type: el middleware AllowContentType solo exige el
	// header en requests con body (POST), no en GET.
	respLista, _ := http.Get(base + "/api/autores")
	fmt.Println("GET /api/autores ->", respLista.StatusCode)
	respLista.Body.Close()

	respPost, _ := http.Post(base+"/api/autores", "application/json", jsonBody(`{"nombre":"Ken Thompson"}`))
	fmt.Println("POST /api/autores ->", respPost.StatusCode)
	respPost.Body.Close()

	respGet, _ := http.Get(base + "/api/autores/1")
	fmt.Println("GET /api/autores/1 ->", respGet.StatusCode)
	respGet.Body.Close()
}
