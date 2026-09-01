# Mini-proyecto: conversor de unidades

Un CLI que convierte entre unidades de longitud, peso y temperatura,
usando **únicamente** lo cubierto en los capítulos 1-10: variables,
funciones, control de flujo, `strings` y `strconv` - sin structs ni
interfaces todavía (eso llega en la [Parte II](../02-tipos-y-composicion/mini-proyecto.md)).

## Qué demuestra

- Funciones como valores (`type conversion struct { ... factor func(v float64) float64 }`),
  para representar cada regla de conversión sin un `switch` gigante.
- Parseo de texto con `strings.Fields` y `strconv.ParseFloat`.
- Errores como valores (`fmt.Errorf` con `%w`) para reportar comandos o
  combinaciones de unidades inválidas.

## Cómo funciona

Cada comando tiene el formato `"<valor> <unidad> a <unidad>"`, por ejemplo
`"10 km a mi"`. `procesarComando` lo parsea, `convertir` busca la regla
correspondiente en las tablas de longitud/peso o cae al caso especial de
temperatura (que necesita una fórmula, no un factor multiplicativo).

## Ejecutarlo

```bash
cd examples/miniproyectos/conversor-unidades
go run .
```

Código fuente: [`examples/miniproyectos/conversor-unidades/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/conversor-unidades)
