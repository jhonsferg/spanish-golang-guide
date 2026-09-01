# Mini-proyecto: cliente HTTP que transforma JSON en un reporte

Un programa que consume una API JSON, transforma los datos crudos en un
resumen legible, y lo escribe tanto a stdout como a un archivo — el flujo
de "traer datos externos y hacer algo útil con ellos" que combina varios
paquetes de la [Parte IV](26-io-package.md) a la vez.

## Qué demuestra

- **`net/http`** para consumir una API (aquí, un servidor local de
  demostración, para no depender de internet).
- **`encoding/json`** para decodificar la respuesta a structs Go
  tipadas, incluyendo el patrón `EventoRaw` (con campos como llegan del
  JSON, p. ej. `Timestamp string`) → `Evento` (con el tipo que realmente
  quieres usar, `Cuando time.Time`) tras parsear.
- **`time.Parse`** con `time.RFC3339` para convertir timestamps de texto
  a `time.Time` manipulable.
- **`strings.Builder`** para construir el reporte de forma eficiente.
- **`os.CreateTemp`** para persistir el resultado a un archivo.

## Por qué separar "raw" de "parseado"

`EventoRaw` existe solo porque así llega el JSON de la API (timestamps
como string). Convertirlo a `Evento` con un `time.Time` real, en un paso
explícito, evita que el resto del programa tenga que lidiar con parseo
de fechas — un patrón que vale la pena en cualquier integración con una
API externa.

## Ejecutarlo

```bash
cd examples/miniproyectos/cliente-http-json
go run .
```

Código fuente: [`examples/miniproyectos/cliente-http-json/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/cliente-http-json)
