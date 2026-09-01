# Mini-proyecto: scraper concurrente con rate limiting

Un "scraper" que descarga varias páginas a la vez, pero nunca más de
`N` al mismo tiempo — el problema real detrás de cualquier cliente que
necesita ser rápido sin tumbar el servidor que consulta.

## Qué demuestra

- **Semáforo con channel con buffer**: `make(chan struct{}, maxConcurrentes)`
  limita cuántas goroutines pueden estar "trabajando" a la vez —
  adquirir es `semaforo <- struct{}{}`, liberar es `<-semaforo`.
- **`sync.WaitGroup`** para esperar a que todas las descargas terminen
  antes de reportar resultados.
- Escritura segura a un slice preasignado por índice (`resultados[idx] = ...`)
  en vez de un channel de resultados — válido porque cada goroutine
  escribe en su propia posición, sin solaparse con las demás.
- Un servidor HTTP local con páginas de distinto tamaño/latencia, para
  que el ejemplo sea reproducible sin depender de internet.

## El resultado que vale la pena mirar

Con `maxConcurrentes=2` y 5 páginas, el tiempo total es notablemente
menor que la suma de todas las latencias individuales, pero mayor que si
las 5 corrieran completamente en paralelo — exactamente el punto medio
que un rate limit real debería producir.

## Ejecutarlo

```bash
cd examples/miniproyectos/scraper-concurrente
go run .
```

Código fuente: [`examples/miniproyectos/scraper-concurrente/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/scraper-concurrente)
