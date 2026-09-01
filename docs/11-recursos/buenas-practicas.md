# Buenas prácticas - Cheatsheet

Recomendaciones consolidadas que aparecen dispersas a lo largo de los 66
capítulos, agrupadas para consulta rápida. Cada bloque enlaza al capítulo
donde se explica el porqué en profundidad.

## Manejo de errores

- Envuelve errores con contexto usando `%w`, no `%v`:
  `fmt.Errorf("leyendo config: %w", err)`.
- Compara con `errors.Is()` / `errors.As()`, nunca con `==` contra un error
  concreto ni con comparación de strings (`err.Error() == "..."`).
- Los errores son valores: puedes agruparlos con `errors.Join()` cuando una
  operación falla por más de una razón.
- No ignores errores con `_` salvo que el caso esté explícitamente
  documentado (p. ej. `defer resp.Body.Close()` cuando no hay nada más que
  hacer con ese error).
- Ver [Capítulo 17 - Manejo de errores](../02-tipos-y-composicion/17-manejo-de-errores.md).

## Concurrencia

- "No comuniques compartiendo memoria; comparte memoria comunicando" - usa
  channels para transferir *ownership* de datos, no para proteger memoria
  compartida que podrías simplemente no compartir.
- Toda goroutine que lanzas debe tener una forma clara de terminar. Si no
  puedes explicar cuándo termina, probablemente tienes un leak.
- Usa `context.Context` para cancelación y deadlines, pásalo como primer
  parámetro, nunca lo guardes en un struct.
- Con `sync.WaitGroup`, prefiere `wg.Go(func() { ... })` (Go 1.25+) a
  `Add`/`defer Done` manuales - ver [Go 1.25](../10-versiones-de-go/go-1-25.md).
- Corre siempre `go test -race` en CI; el race detector encuentra bugs que
  ningún code review va a atrapar a simple vista.
- Ver [Capítulo 21 - Goroutines](../03-concurrencia/21-goroutines.md),
  [Capítulo 25 - Patrones de concurrencia](../03-concurrencia/25-patrones-de-concurrencia.md).

## Diseño de APIs e interfaces

- Acepta interfaces, retorna structs concretos.
- Define interfaces del lado del consumidor, no del proveedor - una
  interfaz de un solo método bien nombrada (`io.Reader`, `io.Writer`) es
  más reusable que una interfaz grande "por si acaso".
- No expongas structs con muchos campos exportados como tu única API
  pública si vas a necesitar evolucionarla; considera opciones funcionales
  (`func WithTimeout(d time.Duration) Option`) para constructores con
  parámetros opcionales.
- Ver [Capítulo 13 - Interfaces](../02-tipos-y-composicion/13-interfaces.md).

## Estructura de proyecto

- `cmd/` para binarios, `internal/` para código que no debe importarse desde
  fuera del módulo, paquetes de dominio en la raíz o bajo un nombre
  descriptivo (no `utils/` ni `common/` como cajón de sastre).
- Un paquete, una responsabilidad. Si el nombre del paquete no te alcanza
  para describir qué hace en una frase, probablemente hace demasiado.
- Ver [Capítulo 20 - Paquetes](../03-concurrencia/20-paquetes.md),
  [Capítulo 45 - Clean architecture](../06-arquitectura-y-sistemas-distribuidos/45-clean-architecture.md).

## Testing

- Tests de tabla (`table-driven tests`) como default para funciones con
  varios casos - más fácil de extender que copiar-pegar `Test_X_CasoY`.
- `t.Parallel()` en tests que no comparten estado, para acelerar la suite.
- `testing/synctest` (Go 1.25+) para probar código con timers/tickers sin
  esperas reales - ver [Go 1.25](../10-versiones-de-go/go-1-25.md).
- Benchmarks con `b.Loop()` (Go 1.24+), no el patrón manual
  `for i := 0; i < b.N; i++`.
- Ver [Capítulo 36 - Testing](../05-produccion-y-herramientas/36-testing-package.md).

## Performance

- Mide antes de optimizar: `pprof` y benchmarks, no intuición.
- Evita allocaciones innecesarias en hot paths: reusa slices con `[:0]`,
  usa `sync.Pool` para objetos de vida corta y alta frecuencia.
- Compila con PGO (`-pgo=auto`) si tienes un perfil representativo de
  producción - mejoras de doble dígito porcentual son comunes sin tocar
  código.
- Ver [Capítulo 40 - Build tools y performance](../05-produccion-y-herramientas/40-build-tools-y-performance.md).

## Seguridad

- Nunca construyas SQL por concatenación de strings; usa placeholders
  (`database/sql`, `sqlc`, GORM con parámetros).
- Valida y sanea toda entrada externa antes de usarla en paths de
  filesystem - considera `os.Root` (Go 1.24+) para confinar el acceso.
- No captures parámetros `io.Reader` de aleatoriedad como "suficientemente
  aleatorios" para crypto salvo que sea explícitamente
  `crypto/rand`.
- Ver [Capítulo 35 - Hash y criptografía](../04-libreria-estandar/35-hash-y-criptografia.md).

## Herramientas del día a día

```bash
gofmt -l .          # ver qué archivos no están formateados
go vet ./...         # análisis estático incluido en el toolchain
go test -race ./...   # tests con detector de data races
go build -pgo=auto     # build con profile-guided optimization
go mod tidy          # sincroniza go.mod/go.sum con el código real
```
