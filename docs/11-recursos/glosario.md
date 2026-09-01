# Glosario

Términos de Go que aparecen a lo largo de la guía, con una definición corta
y el capítulo donde se explican en detalle.

**Blank identifier** - El identificador `_`, usado para descartar un valor
que la sintaxis exige recibir pero que no vas a usar.

**Channel** - Tipo que permite enviar y recibir valores entre goroutines de
forma segura, con o sin buffer. Ver
[Capítulo 22](../03-concurrencia/22-channels.md).

**Composición** - El mecanismo de Go para reutilizar comportamiento vía
embedding de structs o interfaces, en lugar de herencia de clases. Ver
[Capítulo 14](../02-tipos-y-composicion/14-embedding-y-composicion.md).

**Defer** - Instrucción que difiere la ejecución de una llamada hasta que la
función que la contiene retorna. Ver
[Capítulo 19](../02-tipos-y-composicion/19-defer.md).

**Embedding** - Incluir un tipo (struct o interfaz) dentro de otro sin darle
nombre de campo, promoviendo sus métodos y campos al tipo contenedor.

**Generics / Parámetros de tipo** - Mecanismo introducido en Go 1.18 para
escribir funciones y tipos parametrizados por tipo, con constraints que
limitan qué tipos son válidos. Ver [Go 1.18](../10-versiones-de-go/go-1-18.md).

**Goroutine** - Función que se ejecuta de forma concurrente, gestionada por
el runtime de Go (no por el sistema operativo directamente). Mucho más
liviana que un thread de SO. Ver
[Capítulo 21](../03-concurrencia/21-goroutines.md).

**GOPATH** - Modelo de organización de código previo a los módulos de Go
(anterior a Go 1.11). Hoy en día obsoleto para la mayoría de proyectos.

**Interface** - Conjunto de firmas de métodos. En Go se satisfacen de forma
implícita: un tipo implementa una interfaz con solo tener los métodos
correctos, sin declararlo explícitamente. Ver
[Capítulo 13](../02-tipos-y-composicion/13-interfaces.md).

**Módulo (Go module)** - Unidad de versionado y distribución de código Go,
definida por un archivo `go.mod`. Ver
[Capítulo 20](../03-concurrencia/20-paquetes.md).

**Panic / Recover** - Mecanismo de Go para errores irrecuperables en tiempo
de ejecución (no para control de flujo normal, a diferencia de las
excepciones en otros lenguajes). Ver
[Capítulo 18](../02-tipos-y-composicion/18-panic-y-recover.md).

**Puntero** - Variable que almacena la dirección de memoria de otro valor.
Go no tiene aritmética de punteros como C. Ver
[Capítulo 16](../02-tipos-y-composicion/16-punteros.md).

**Receiver** - El parámetro que precede al nombre de una función para
declararla como método de un tipo (`func (s *Server) Start() error`). Ver
[Capítulo 12](../02-tipos-y-composicion/12-metodos.md).

**Select** - Instrucción que permite esperar sobre múltiples operaciones de
channel a la vez. Ver [Capítulo 23](../03-concurrencia/23-select.md).

**Slice** - Vista dinámica y redimensionable sobre un array subyacente; la
colección de uso más común en Go. Ver
[Capítulo 8](../01-fundamentos/08-arreglos-y-slices.md).

**Struct tag** - Metadata textual asociada a un campo de struct (p. ej.
`json:"name"`), leída en runtime vía reflection por librerías como
`encoding/json`. Ver [Capítulo 29](../04-libreria-estandar/29-json-y-encoding.md).

**Worker pool** - Patrón de concurrencia donde un número fijo de goroutines
consume tareas de un channel compartido. Ver
[Capítulo 25](../03-concurrencia/25-patrones-de-concurrencia.md).

**Zero value** - El valor por defecto que recibe una variable declarada sin
inicializar explícitamente (`0` para numéricos, `""` para strings, `nil`
para punteros/slices/maps/interfaces/channels/funciones).
