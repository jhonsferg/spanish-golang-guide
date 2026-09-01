// Ejemplo del Capítulo 55 (versión condensada del proyecto integrado):
// un microservicio con apagado ordenado (graceful shutdown) — el patrón
// que necesita cualquier servicio que corra en Kubernetes/Docker para
// terminar en curso las requests activas antes de morir, en vez de
// cortarlas a mitad.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

func nuevoServidor(logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /trabajo-lento", func(w http.ResponseWriter, r *http.Request) {
		// Simula una request en curso que el shutdown ordenado debe
		// esperar a que termine, en vez de cortarla.
		time.Sleep(150 * time.Millisecond)
		fmt.Fprintln(w, "trabajo lento completado")
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &http.Server{Handler: mux}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		logger.Error("no se pudo escuchar", "error", err)
		return
	}

	server := nuevoServidor(logger)

	servidorTerminado := make(chan struct{})
	go func() {
		defer close(servidorTerminado)
		logger.Info("servidor iniciado", "addr", ln.Addr().String())
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("el servidor terminó con error", "error", err)
		}
	}()

	base := "http://" + ln.Addr().String()
	cliente := &http.Client{}

	// Lanza una request lenta EN VUELO justo antes de iniciar el shutdown,
	// para demostrar que graceful shutdown la espera en vez de cortarla.
	respuestaLenta := make(chan int, 1)
	go func() {
		resp, err := cliente.Get(base + "/trabajo-lento")
		if err != nil {
			respuestaLenta <- 0
			return
		}
		defer resp.Body.Close()
		respuestaLenta <- resp.StatusCode
	}()

	time.Sleep(20 * time.Millisecond) // deja que la request lenta arranque

	// --- Shutdown ordenado ---
	// En un servicio real, esto se dispara al recibir SIGTERM
	// (signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)).
	// Aquí lo disparamos directamente para que el ejemplo sea autocontenido.
	logger.Info("iniciando apagado ordenado")
	ctxApagado, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxApagado); err != nil {
		logger.Error("error durante el shutdown", "error", err)
	}

	<-servidorTerminado

	status := <-respuestaLenta
	logger.Info("resultado de la request lenta iniciada antes del shutdown", "status", status)
	if status != http.StatusOK {
		logger.Error("la request en vuelo NO se completó correctamente: el shutdown la cortó")
		return
	}
	logger.Info("shutdown ordenado completado sin cortar requests en curso")
}
