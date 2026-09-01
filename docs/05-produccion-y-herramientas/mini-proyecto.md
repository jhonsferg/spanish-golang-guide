# Mini-proyecto: CLI con tests, benchmarks y profiling

Un CLI de estadísticas de texto (palabra más frecuente, longitud
promedio) con la lógica de negocio separada del `main()` en un paquete
propio — precisamente para que sea testeable sin tocar entrada/salida de
consola, el tema central de la [Parte V](36-testing-package.md).

## Qué demuestra

- **Separación CLI / lógica**: `textstats/` no importa `fmt` para
  imprimir nada; solo calcula y retorna valores. `main.go` es la única
  capa que sabe que existe una terminal.
- **Tests de tabla** (`TestContarPalabras`) cubriendo texto normal, texto
  vacío y texto con números.
- **Casos límite explícitos**: `TestPalabraMasFrecuente_MapaVacio` prueba
  qué pasa con un mapa vacío — el tipo de caso que se olvida si no se
  escribe a propósito.
- **Benchmark con `b.Loop()`** (Go 1.24+) para medir `ContarPalabras`
  sobre un texto repetido.

## Perfilar el benchmark con pprof

```bash
cd examples/miniproyectos/cli-con-tests/textstats
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof cpu.prof
```

Dentro de `pprof`, el comando `top` muestra qué funciones consumen más
CPU — en este caso, casi todo el tiempo debería estar en `strings.FieldsFunc`.

## Ejecutarlo

```bash
cd examples/miniproyectos/cli-con-tests
go run .
go test ./... -v -bench=.
```

Código fuente: [`examples/miniproyectos/cli-con-tests/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/cli-con-tests)
