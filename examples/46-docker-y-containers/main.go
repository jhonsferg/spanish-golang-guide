// Ejemplo del Capítulo 46: el tipo de servidor HTTP típico que se
// containeriza - puerto configurable por variable de entorno (el
// contrato usual en Docker/Kubernetes) y endpoint /healthz para
// healthchecks. Ver Dockerfile en este mismo directorio.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func puertoDesdeEnv() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "0" // 0 = el SO elige un puerto libre (útil para este demo)
}

func nuevoMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "servicio corriendo dentro (o fuera) de un contenedor")
	})
	return mux
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:"+puertoDesdeEnv())
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: nuevoMux()}
	go server.Serve(ln)
	defer server.Close()

	// En un contenedor real, aquí seguiría bloqueado sirviendo tráfico.
	// Para este ejemplo autocontenido, simulamos un healthcheck y salimos.
	resp, err := http.Get("http://" + ln.Addr().String() + "/healthz")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("healthcheck respondió con status:", resp.StatusCode)
}
