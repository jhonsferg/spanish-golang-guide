// Ejemplo del Capítulo 49: métricas expuestas con expvar, el paquete de
// la librería estándar para publicar contadores y gauges vía HTTP
// (más simple que Prometheus para casos básicos, y sin dependencias).
package main

import (
	"expvar"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
)

var (
	solicitudesTotal = expvar.NewInt("solicitudes_total")
	erroresTotal     = expvar.NewInt("errores_total")
	tareasEnCurso    atomic.Int64
	gaugeTareas      = expvar.NewInt("tareas_en_curso")
)

func manejarSolicitud(exito bool) {
	solicitudesTotal.Add(1)
	if !exito {
		erroresTotal.Add(1)
	}

	tareasEnCurso.Add(1)
	gaugeTareas.Set(tareasEnCurso.Load())
	defer func() {
		tareasEnCurso.Add(-1)
		gaugeTareas.Set(tareasEnCurso.Load())
	}()
}

func main() {
	// expvar se registra automáticamente en /debug/vars al importarlo.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	server := &http.Server{Handler: http.DefaultServeMux}
	go server.Serve(ln)
	defer server.Close()

	manejarSolicitud(true)
	manejarSolicitud(true)
	manejarSolicitud(false)

	resp, err := http.Get("http://" + ln.Addr().String() + "/debug/vars")
	if err != nil {
		fmt.Println("error consultando métricas:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("status /debug/vars:", resp.StatusCode)
	fmt.Println("solicitudes_total:", solicitudesTotal.Value())
	fmt.Println("errores_total:", erroresTotal.Value())
	fmt.Println("tareas_en_curso (debe ser 0, todas terminaron):", gaugeTareas.Value())
}
