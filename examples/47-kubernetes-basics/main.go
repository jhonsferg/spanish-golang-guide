// Ejemplo del Capítulo 47: los tres endpoints que Kubernetes espera de
// un pod bien portado — liveness, readiness y startup probes — y una
// bandera de "listo" que se activa después de una inicialización
// simulada. Ver deployment.yaml en este directorio para el manifiesto
// que los consume.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

func main() {
	var listo atomic.Bool // false hasta que la "inicialización" termine

	mux := http.NewServeMux()
	mux.HandleFunc("GET /livez", func(w http.ResponseWriter, r *http.Request) {
		// Vivo mientras el proceso responda — si esto falla, Kubernetes
		// reinicia el pod.
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		// Listo para recibir tráfico — si falla, Kubernetes deja de
		// enrutar requests a este pod sin reiniciarlo.
		if listo.Load() {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: mux}
	go server.Serve(ln)
	defer server.Close()

	base := "http://" + ln.Addr().String()

	resp, _ := http.Get(base + "/readyz")
	fmt.Println("readyz antes de inicializar:", resp.StatusCode)
	resp.Body.Close()

	time.Sleep(20 * time.Millisecond) // simula carga de config, warm-up, etc.
	listo.Store(true)

	resp2, _ := http.Get(base + "/readyz")
	fmt.Println("readyz después de inicializar:", resp2.StatusCode)
	resp2.Body.Close()

	resp3, _ := http.Get(base + "/livez")
	fmt.Println("livez:", resp3.StatusCode)
	resp3.Body.Close()
}
