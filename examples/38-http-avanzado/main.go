// Ejemplo del Capítulo 38: servidor HTTP con routing por método,
// middleware encadenado y timeouts - usando solo net/http, sin
// frameworks (ver Parte VIII para Gin/Echo/Fiber/Chi).
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// conLogging es un middleware clásico: envuelve un handler y registra
// cada request antes de delegar al handler original.
func conLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s -> %v", r.Method, r.URL.Path, time.Since(inicio))
	})
}

// encadenar aplica varios middlewares en orden, de afuera hacia adentro.
func encadenar(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type Item struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

func nuevoMux() http.Handler {
	mux := http.NewServeMux()

	items := map[int]Item{1: {1, "teclado"}, 2: {2, "mouse"}}

	// Routing con método y wildcard, disponible desde Go 1.22.
	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		lista := make([]Item, 0, len(items))
		for _, it := range items {
			lista = append(lista, it)
		}
		json.NewEncoder(w).Encode(lista)
	})

	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var id int
		fmt.Sscanf(r.PathValue("id"), "%d", &id)
		item, ok := items[id]
		if !ok {
			http.Error(w, `{"error":"no encontrado"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(item)
	})

	return encadenar(mux, conLogging)
}

func main() {
	server := &http.Server{
		Handler:      nuevoMux(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go server.Serve(ln)
	defer server.Close()

	base := "http://" + ln.Addr().String()

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(base + "/items")
	if err != nil {
		log.Fatal(err)
	}
	var items []Item
	json.NewDecoder(resp.Body).Decode(&items)
	resp.Body.Close()
	fmt.Println("GET /items ->", items)

	resp2, err := client.Get(base + "/items/1")
	if err != nil {
		log.Fatal(err)
	}
	var item Item
	json.NewDecoder(resp2.Body).Decode(&item)
	resp2.Body.Close()
	fmt.Println("GET /items/1 ->", item)

	resp3, err := client.Get(base + "/items/999")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET /items/999 -> status", resp3.StatusCode)
	resp3.Body.Close()
}
