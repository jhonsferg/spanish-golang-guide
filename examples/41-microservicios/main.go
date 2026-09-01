// Ejemplo del Capítulo 41: dos "microservicios" HTTP mínimos (catálogo y
// pedidos) donde uno llama al otro — usando solo net/http, para
// enfocarse en el patrón de comunicación, no en un framework.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

type Producto struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
}

type Pedido struct {
	ProductoID int     `json:"producto_id"`
	Producto   string  `json:"producto"`
	Total      float64 `json:"total"`
}

func arrancar(handler http.Handler) (baseURL string, cerrar func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: handler}
	go server.Serve(ln)
	return "http://" + ln.Addr().String(), func() { server.Close() }
}

// servicioCatalogo simula un microservicio independiente con su propia
// "base de datos" en memoria.
func servicioCatalogo() http.Handler {
	productos := map[int]Producto{
		1: {1, "Teclado mecánico", 89.90},
		2: {2, "Mouse ergonómico", 39.90},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /productos/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		fmt.Sscanf(r.PathValue("id"), "%d", &id)
		p, ok := productos[id]
		if !ok {
			http.Error(w, "producto no encontrado", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(p)
	})
	return mux
}

// servicioPedidos es OTRO microservicio: no comparte memoria con el
// catálogo, se comunica con él exclusivamente por HTTP.
func servicioPedidos(catalogoURL string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pedidos/{productoID}", func(w http.ResponseWriter, r *http.Request) {
		productoID := r.PathValue("productoID")

		resp, err := http.Get(catalogoURL + "/productos/" + productoID)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "no se pudo consultar el catálogo", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		var p Producto
		json.NewDecoder(resp.Body).Decode(&p)

		pedido := Pedido{ProductoID: p.ID, Producto: p.Nombre, Total: p.Precio}
		json.NewEncoder(w).Encode(pedido)
	})
	return mux
}

func main() {
	catalogoURL, cerrarCatalogo := arrancar(servicioCatalogo())
	defer cerrarCatalogo()

	pedidosURL, cerrarPedidos := arrancar(servicioPedidos(catalogoURL))
	defer cerrarPedidos()

	resp, err := http.Post(pedidosURL+"/pedidos/1", "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var pedido Pedido
	json.NewDecoder(resp.Body).Decode(&pedido)
	fmt.Printf("pedido creado vía comunicación entre microservicios: %+v\n", pedido)
}
