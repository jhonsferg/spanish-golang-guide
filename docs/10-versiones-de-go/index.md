# Versiones de Go

Go publica una nueva versión menor cada seis meses (febrero y agosto), con
compatibilidad hacia atrás como promesa central del lenguaje. Esta sección
resume qué trajo cada versión desde **Go 1.18** (la que introdujo generics,
el punto de inflexión de la era moderna de Go) hasta **Go 1.27**, la más
reciente.

Todo el contenido de esta sección proviene de las release notes oficiales en
[go.dev/doc](https://go.dev/doc/devel/release), no de memoria — cada página
resume lo verificable en la nota oficial de esa versión.

## Tabla resumen

| Versión | Fecha | Lo más destacado |
|---|---|---|
| [Go 1.18](go-1-18.md) | mar. 2022 | **Generics** (parámetros de tipo), fuzzing integrado, workspaces (`go.work`), paquete `netip` |
| [Go 1.19](go-1-19.md) | ago. 2022 | Soft memory limit (`GOMEMLIMIT`), tipos atómicos en `sync/atomic`, doc comments enriquecidos |
| [Go 1.20](go-1-20.md) | feb. 2023 | Conversión slice→array, `errors.Join`, PGO (profile-guided optimization) en preview |
| [Go 1.21](go-1-21.md) | ago. 2023 | `min`/`max`/`clear` built-in, `log/slog`, paquetes `slices`/`maps`/`cmp`, PGO en producción |
| [Go 1.22](go-1-22.md) | feb. 2024 | Variable de loop por iteración, `range` sobre enteros, routing por método/wildcards en `net/http`, `math/rand/v2` |
| [Go 1.23](go-1-23.md) | ago. 2024 | `range` sobre funciones (iteradores), paquete `iter`, `unique`, telemetría opt-in |
| [Go 1.24](go-1-24.md) | feb. 2025 | Generic type aliases, Swiss Tables para `map`, `os.Root` (FS sandboxing), `crypto/mlkem` (post-cuántico) |
| [Go 1.25](go-1-25.md) | ago. 2025 | `GOMAXPROCS` consciente de cgroups, GC "Green Tea" (experimental), `testing/synctest`, `sync.WaitGroup.Go` |
| [Go 1.26](go-1-26.md) | feb. 2026 | GC "Green Tea" por defecto, generics autorreferenciables, `go fix` como modernizador |
| [Go 1.27](go-1-27.md) | ago. 2026 | Métodos genéricos, `encoding/json/v2`, `crypto/mldsa` (post-cuántico), asignación de memoria más rápida |

## Cómo usar esta sección

Si vienes de una versión antigua de Go, lee las páginas en orden desde tu
versión actual hasta [Go 1.27](go-1-27.md) — cada una lista solo lo que
cambió *en esa versión*, no un acumulado. Si solo te interesa qué hay de
nuevo hoy, ve directo a [Go 1.27](go-1-27.md).
