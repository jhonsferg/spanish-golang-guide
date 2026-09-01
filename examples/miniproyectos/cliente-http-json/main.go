// Mini-proyecto de la Parte IV (Librería estándar): un cliente que
// consume una API JSON, transforma los datos y genera un reporte de
// texto - combinando net/http, encoding/json, time, strings y os en un
// flujo real de "obtener datos externos y hacer algo útil con ellos".
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type EventoRaw struct {
	Usuario   string `json:"usuario"`
	Accion    string `json:"accion"`
	Timestamp string `json:"timestamp"` // RFC3339, como llega de la mayoría de APIs
}

type Evento struct {
	Usuario string
	Accion  string
	Cuando  time.Time
}

func levantarAPIDemo() (baseURL string, cerrar func()) {
	eventos := []EventoRaw{
		{"ana", "login", "2026-09-01T08:00:00Z"},
		{"beto", "login", "2026-09-01T08:05:00Z"},
		{"ana", "compra", "2026-09-01T08:10:00Z"},
		{"carla", "login", "2026-09-01T08:12:00Z"},
		{"ana", "logout", "2026-09-01T08:20:00Z"},
		{"beto", "compra", "2026-09-01T08:25:00Z"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /eventos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(eventos)
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := &http.Server{Handler: mux}
	go server.Serve(ln)
	return "http://" + ln.Addr().String(), func() { server.Close() }
}

func obtenerEventos(baseURL string) ([]Evento, error) {
	resp, err := http.Get(baseURL + "/eventos")
	if err != nil {
		return nil, fmt.Errorf("consultando API: %w", err)
	}
	defer resp.Body.Close()

	var raw []EventoRaw
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decodificando respuesta: %w", err)
	}

	eventos := make([]Evento, 0, len(raw))
	for _, r := range raw {
		cuando, err := time.Parse(time.RFC3339, r.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("timestamp inválido %q: %w", r.Timestamp, err)
		}
		eventos = append(eventos, Evento{Usuario: r.Usuario, Accion: r.Accion, Cuando: cuando})
	}
	return eventos, nil
}

// generarReporte transforma los eventos crudos en un resumen por usuario.
func generarReporte(eventos []Evento) string {
	porUsuario := make(map[string][]string)
	for _, e := range eventos {
		porUsuario[e.Usuario] = append(porUsuario[e.Usuario], e.Accion)
	}

	usuarios := make([]string, 0, len(porUsuario))
	for u := range porUsuario {
		usuarios = append(usuarios, u)
	}
	sort.Strings(usuarios)

	var sb strings.Builder
	sb.WriteString("Reporte de actividad\n")
	sb.WriteString("=====================\n")
	for _, u := range usuarios {
		fmt.Fprintf(&sb, "%s: %s (%d acciones)\n", u, strings.Join(porUsuario[u], " -> "), len(porUsuario[u]))
	}
	return sb.String()
}

func main() {
	base, cerrar := levantarAPIDemo()
	defer cerrar()

	eventos, err := obtenerEventos(base)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	reporte := generarReporte(eventos)
	fmt.Print(reporte)

	archivo, err := os.CreateTemp("", "reporte-*.txt")
	if err != nil {
		fmt.Println("error creando archivo:", err)
		return
	}
	defer os.Remove(archivo.Name())

	if _, err := archivo.WriteString(reporte); err != nil {
		fmt.Println("error escribiendo reporte:", err)
		return
	}
	archivo.Close()
	fmt.Println("\nreporte también escrito en:", archivo.Name())
}
