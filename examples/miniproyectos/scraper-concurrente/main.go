// Mini-proyecto de la Parte III (Concurrencia): un "scraper" concurrente
// con límite de peticiones simultáneas (rate limiting vía semáforo con
// channel) y agregación de resultados — el patrón real detrás de
// cualquier crawler o cliente que golpea una API con cuidado de no
// saturarla. Usa servidores HTTP locales para no depender de internet.
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Resultado struct {
	URL      string
	Bytes    int
	Duracion time.Duration
	Err      error
}

// scrapearTodas lanza una goroutine por URL, pero un semáforo con
// buffer de tamaño `maxConcurrentes` limita cuántas corren a la vez.
func scrapearTodas(urls []string, maxConcurrentes int) []Resultado {
	semaforo := make(chan struct{}, maxConcurrentes)
	resultados := make([]Resultado, len(urls))

	var wg sync.WaitGroup
	for i, url := range urls {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()

			semaforo <- struct{}{}        // adquiere un "slot"
			defer func() { <-semaforo }() // libera el slot al terminar

			inicio := time.Now()
			cliente := &http.Client{Timeout: 2 * time.Second}
			resp, err := cliente.Get(url)
			if err != nil {
				resultados[idx] = Resultado{URL: url, Err: err}
				return
			}
			defer resp.Body.Close()

			cuerpo, err := io.ReadAll(resp.Body)
			resultados[idx] = Resultado{
				URL:      url,
				Bytes:    len(cuerpo),
				Duracion: time.Since(inicio),
				Err:      err,
			}
		}(i, url)
	}
	wg.Wait()
	return resultados
}

// levantarServidorDemo simula páginas de distinto tamaño y latencia,
// para no depender de internet dentro de este ejemplo.
func levantarServidorDemo() (baseURL string, cerrar func()) {
	mux := http.NewServeMux()
	paginas := map[string]struct {
		tamano  int
		retraso time.Duration
	}{
		"/pagina-a": {tamano: 500, retraso: 10 * time.Millisecond},
		"/pagina-b": {tamano: 2000, retraso: 40 * time.Millisecond},
		"/pagina-c": {tamano: 100, retraso: 5 * time.Millisecond},
		"/pagina-d": {tamano: 1500, retraso: 25 * time.Millisecond},
		"/pagina-e": {tamano: 800, retraso: 15 * time.Millisecond},
	}
	for ruta, cfg := range paginas {
		tamano, retraso := cfg.tamano, cfg.retraso
		mux.HandleFunc(ruta, func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(retraso)
			w.Write(make([]byte, tamano))
		})
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := &http.Server{Handler: mux}
	go server.Serve(ln)
	return "http://" + ln.Addr().String(), func() { server.Close() }
}

func main() {
	base, cerrar := levantarServidorDemo()
	defer cerrar()

	urls := []string{
		base + "/pagina-a",
		base + "/pagina-b",
		base + "/pagina-c",
		base + "/pagina-d",
		base + "/pagina-e",
	}

	inicio := time.Now()
	resultados := scrapearTodas(urls, 2) // máximo 2 requests simultáneas
	duracionTotal := time.Since(inicio)

	sort.Slice(resultados, func(i, j int) bool { return resultados[i].URL < resultados[j].URL })

	totalBytes := 0
	for _, r := range resultados {
		if r.Err != nil {
			fmt.Printf("%s -> error: %v\n", r.URL, r.Err)
			continue
		}
		fmt.Printf("%s -> %d bytes en %v\n", r.URL, r.Bytes, r.Duracion.Round(time.Millisecond))
		totalBytes += r.Bytes
	}
	fmt.Printf("total: %d bytes de %d páginas en %v (con máx. 2 concurrentes)\n",
		totalBytes, len(urls), duracionTotal.Round(time.Millisecond))
}
