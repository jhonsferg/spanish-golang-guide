# Capítulo 26: IO - Lectura y escritura de datos

## Índice
1. [¿Qué es el io Package?](#261-qué-es-el-io-package)
2. [Reader Interface](#262-reader-interface)
3. [Writer Interface](#263-writer-interface)
4. [ReadWriter Interface](#264-readwriter-interface)
5. [Closer Interface](#265-closer-interface)
6. [Copy y CopyN](#266-copy-y-copyn)
7. [Pipe](#267-pipe)
8. [MultiReader y MultiWriter](#268-multireader-y-multiwriter)
9. [ReadAll y ReadAtLeast](#269-readall-y-readatleast)
10. [Buffering con bufio](#2610-buffering-con-bufio)
11. [Buenas Prácticas y Patrones](#2611-buenas-prácticas-y-patrones)

---

## 26.1 ¿Qué es el io Package?

### 26.1.1 Conceptos Fundamentales

El package `io` es el corazón de las operaciones de entrada/salida en Go. Proporciona interfaces abstractas que definen cómo los datos pueden ser leídos y escritos, permitiendo que diferentes fuentes y destinos trabajen juntas sin acoplamiento.

**Características principales:**
- Interfaces simples y composables (Reader, Writer, Closer)
- Funciones utilitarias para copiar y manipular datos
- Soporte para streaming de datos
- Manejo eficiente de recursos

**Por qué interfaces en lugar de tipos concretos:**
```
Go prefiere interfaces pequeñas y específicas sobre tipos grandes.
Esto permite escribir código genérico que funciona con cualquier
implementador de la interfaz, sin importar el tipo concreto.
```

### 26.1.2 Flujo de Datos en Go

```
┌─────────────────────────────────────────────────────┐
│                 FLUJO DE DATOS GO                   │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Fuente (Reader)      →    Destino (Writer)        │
│  ┌──────────────┐          ┌──────────────┐        │
│  │  Archivo     │          │  Archivo     │        │
│  │  Conexión    │  datos   │  Conexión    │        │
│  │  Buffer      │    →     │  Buffer      │        │
│  │  Custom      │          │  Custom      │        │
│  └──────────────┘          └──────────────┘        │
│                                                     │
│  Con potencial transformación (Filter)             │
│  ┌──────────────┐      ┌──────────────┐           │
│  │  Reader      │  →   │ Compresión   │  →        │
│  │              │      │ Encriptación │           │
│  └──────────────┘      │ Filtros      │           │
│                        └──────────────┘           │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 26.1.3 Interfaces Principales

Go proporciona un conjunto mínimo de interfaces que componen todas las operaciones de I/O:

```go
// Reader: lectura de datos
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer: escritura de datos
type Writer interface {
    Write(p []byte) (n int, err error)
}

// Closer: liberar recursos
type Closer interface {
    Close() error
}

// ReadCloser: leer y cerrar
type ReadCloser interface {
    Reader
    Closer
}

// WriteCloser: escribir y cerrar
type WriteCloser interface {
    Writer
    Closer
}

// ReadWriter: leer y escribir
type ReadWriter interface {
    Reader
    Writer
}
```

### 26.1.4 Comparación con Otros Lenguajes

**Java Streams:**
```java
// Java - más complejo, múltiples métodos
InputStream is = new FileInputStream("file.txt");
byte[] buffer = new byte[1024];
int bytesRead = is.read(buffer);
```

**Go - minimalista:**
```go
// Go - una función, una responsabilidad
var r io.Reader // cualquier implementador funciona
n, err := r.Read(buffer)
```

**C++ iostream:**
```cpp
// C++ - sobrecarga de operadores, más verboso
std::ifstream file("file.txt");
char buffer[1024];
file.read(buffer, sizeof(buffer));
```

---

## 26.2 Reader Interface

### 26.2.1 Definición y Semántica

La interfaz `Reader` es la más fundamental:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**Semántica del método Read:**
- `p`: buffer donde escribir los datos leídos
- `n`: número de bytes realmente leídos (≤ len(p))
- `err`: error durante la lectura
- Si `n < len(p)`, probablemente EOF o error
- `n > 0` y `err == nil`: datos válidos
- `n > 0` y `err == EOF`: últimos datos antes del fin
- `n == 0` y `err != nil`: error sin datos

### 26.2.2 Implementadores Comunes de Reader

**Archivo:**
```go
package main

import (
	"fmt"
	"os"
)

func ejemploArchivoReader() {
	file, err := os.Open("datos.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Leídos %d bytes: %s\n", n, string(buffer[:n]))
}
```

**Strings (io.Reader sobre cadena):**
```go
package main

import (
	"fmt"
	"io"
	"strings"
)

func ejemploStringReader() {
	reader := strings.NewReader("Hola, Go!")
	buffer := make([]byte, 5)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Leído: %s\n", string(buffer[:n]))
	}
	// Salida:
	// Leído: Hola,
	// Leído:  Go!
}
```

**Bytes (buffer en memoria):**
```go
package main

import (
	"bytes"
	"fmt"
	"io"
)

func ejemploBytesReader() {
	data := []byte("Datos en memoria")
	reader := bytes.NewReader(data)

	buffer := make([]byte, 100)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Printf("Leído: %s\n", string(buffer[:n]))
}
```

**Conexión de red:**
```go
package main

import (
	"fmt"
	"io"
	"net"
)

func ejemploNetReader() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Printf("Recibido: %s\n", string(buffer[:n]))
}
```

### 26.2.3 Patrón de Lectura Seguro

```go
package main

import (
	"io"
	"os"
)

func leerCompleto(r io.Reader) ([]byte, error) {
	const maxBytes = 1024 * 1024 // 1 MB límite de seguridad

	var result []byte
	buffer := make([]byte, 32*1024) // Buffer de 32 KB

	for {
		n, err := r.Read(buffer)

		// Procesar datos leídos antes de verificar error
		if n > 0 {
			result = append(result, buffer[:n]...)

			// Verificar límite de seguridad
			if len(result) > maxBytes {
				return nil, io.ErrUnexpectedEOF
			}
		}

		// Verificar error
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func main() {
	file, _ := os.Open("documento.txt")
	defer file.Close()

	datos, err := leerCompleto(file)
	if err != nil {
		panic(err)
	}

	_ = datos // usar los datos
}
```

### 26.2.4 Implementar un Reader Personalizado

```go
package main

import (
	"io"
	"strings"
)

// RepetidoReader repite una cadena N veces
type RepetidoReader struct {
	contenido string
	repeticiones int
	posicion   int
	iteracion  int
}

func NewRepetidoReader(contenido string, repeticiones int) *RepetidoReader {
	return &RepetidoReader{
		contenido:   contenido,
		repeticiones: repeticiones,
	}
}

func (r *RepetidoReader) Read(p []byte) (int, error) {
	if r.iteracion >= r.repeticiones {
		return 0, io.EOF
	}

	// Crear un reader para la iteración actual
	reader := strings.NewReader(r.contenido)
	reader.Seek(int64(r.posicion), io.SeekStart)

	n, err := reader.Read(p)
	r.posicion += n

	// Si llegamos al final del contenido, pasar a la siguiente iteración
	if r.posicion >= len(r.contenido) {
		r.posicion = 0
		r.iteracion++
	}

	// Si es la última iteración y no hay más datos, retornar EOF
	if r.iteracion >= r.repeticiones && n == 0 {
		return 0, io.EOF
	}

	return n, err
}

func main() {
	reader := NewRepetidoReader("Hola ", 3)
	buffer := make([]byte, 20)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		println("Leído:", string(buffer[:n]))
	}
	// Salida: "Hola Hola Hola "
}
```

---

## 26.3 Writer Interface

### 26.3.1 Definición y Semántica

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**Semántica:**
- `p`: bytes a escribir
- `n`: número de bytes realmente escritos
- `err`: error durante la escritura
- Si `n < len(p)`, debe haber un error
- Si `err == nil`, se escribieron todos los bytes (n == len(p))

### 26.3.2 Implementadores Comunes de Writer

**Archivo:**
```go
package main

import (
	"os"
)

func escribirEnArchivo() {
	file, err := os.Create("salida.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := []byte("Hello, Go!")
	n, err := file.Write(data)
	if err != nil {
		panic(err)
	}

	if n != len(data) {
		panic("No se escribieron todos los bytes")
	}
}
```

**Bytes Buffer:**
```go
package main

import (
	"bytes"
	"fmt"
)

func escribirEnBuffer() {
	var buf bytes.Buffer

	// Buffer implementa Writer
	buf.Write([]byte("Primera línea\n"))
	buf.Write([]byte("Segunda línea\n"))

	fmt.Print(buf.String())
	// Salida:
	// Primera línea
	// Segunda línea
}
```

**Conexión de red:**
```go
package main

import (
	"net"
)

func escribirEnRed() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	mensaje := []byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")
	n, err := conn.Write(mensaje)
	if err != nil {
		panic(err)
	}

	if n != len(mensaje) {
		panic("No se escribió el mensaje completo")
	}
}
```

**Stdout/Stderr:**
```go
package main

import (
	"os"
)

func ejemploStdout() {
	os.Stdout.Write([]byte("Salida estándar\n"))
	os.Stderr.Write([]byte("Error estándar\n"))
}
```

### 26.3.3 Patrón de Escritura Segura

```go
package main

import (
	"io"
	"os"
)

func escribirCompleto(w io.Writer, datos []byte) error {
	escrito := 0

	for escrito < len(datos) {
		n, err := w.Write(datos[escrito:])

		if err != nil {
			return err
		}

		if n <= 0 {
			return io.ErrShortWrite
		}

		escrito += n
	}

	return nil
}

func main() {
	file, _ := os.Create("resultado.txt")
	defer file.Close()

	datos := []byte("Datos a escribir completamente")
	if err := escribirCompleto(file, datos); err != nil {
		panic(err)
	}
}
```

### 26.3.4 Implementar un Writer Personalizado

```go
package main

import (
	"fmt"
	"io"
	"strings"
)

// UpperWriter escribe datos en mayúsculas
type UpperWriter struct {
	dest io.Writer
}

func NewUpperWriter(dest io.Writer) *UpperWriter {
	return &UpperWriter{dest: dest}
}

func (uw *UpperWriter) Write(p []byte) (int, error) {
	// Convertir a mayúsculas
	upper := strings.ToUpper(string(p))
	return uw.dest.Write([]byte(upper))
}

// HexDumpWriter muestra datos en hexadecimal
type HexDumpWriter struct{}

func (hdw *HexDumpWriter) Write(p []byte) (int, error) {
	fmt.Printf("HEX: ")
	for _, b := range p {
		fmt.Printf("%02x ", b)
	}
	fmt.Println()
	return len(p), nil
}

func main() {
	// UpperWriter
	var buf strings.Builder
	upper := NewUpperWriter(&buf)
	upper.Write([]byte("hola mundo"))
	fmt.Println("Resultado:", buf.String()) // HOLA MUNDO

	// HexDumpWriter
	hex := &HexDumpWriter{}
	hex.Write([]byte("ABC"))
	// HEX: 41 42 43
}
```

---

## 26.4 ReadWriter Interface

### 26.4.1 Composición de Interfaces

ReadWriter combina Reader y Writer en una sola interfaz:

```go
type ReadWriter interface {
    Reader
    Writer
}
```

**Ventajas de la composición:**
- Una sola interfaz para operaciones bidireccionales
- Reduce parámetros en funciones
- Permite reutilizar código

### 26.4.2 Casos de Uso

**Conexiones bidireccionales:**
```go
package main

import (
	"io"
	"net"
)

func manejarConexion(conn io.ReadWriter) {
	// conn es bidireccional (TCP socket)
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	respuesta := buffer[:n]
	conn.Write(respuesta)
}

func servidor() {
	listener, _ := net.Listen("tcp", ":8080")
	conn, _ := listener.Accept()
	defer conn.Close()

	// conn implementa io.ReadWriter
	manejarConexion(conn)
}
```

**Procesos con stdio bidireccional:**
```go
package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

func ejecutarComandoInteractivo(comando string) error {
	cmd := exec.Command(comando)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	// Crear un tipo que combine stdin y stdout
	type bidireccional struct {
		io.Reader
		io.Writer
	}

	rw := &bidireccional{
		Reader: stdout,
		Writer: stdin,
	}

	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		println("Recibido:", scanner.Text())
	}

	return cmd.Start()
}
```

### 26.4.3 Implementar ReadWriter Personalizado

```go
package main

import (
	"bytes"
	"io"
	"sync"
)

// EchoReadWriter devuelve lo que recibe
type EchoReadWriter struct {
	buffer *bytes.Buffer
	mu     sync.Mutex
}

func NewEchoReadWriter() *EchoReadWriter {
	return &EchoReadWriter{
		buffer: &bytes.Buffer{},
	}
}

func (erw *EchoReadWriter) Read(p []byte) (int, error) {
	erw.mu.Lock()
	defer erw.mu.Unlock()
	return erw.buffer.Read(p)
}

func (erw *EchoReadWriter) Write(p []byte) (int, error) {
	erw.mu.Lock()
	defer erw.mu.Unlock()

	// Guardar en buffer y también hacer echo
	n, err := erw.buffer.Write(p)
	return n, err
}

// DuplexReader: leer de una fuente, escribir en otra
type DuplexReader struct {
	lector io.Reader
	escritor io.Writer
}

func NewDuplexReader(r io.Reader, w io.Writer) *DuplexReader {
	return &DuplexReader{lector: r, escritor: w}
}

func (dr *DuplexReader) Read(p []byte) (int, error) {
	n, err := dr.lector.Read(p)
	if n > 0 {
		dr.escritor.Write(p[:n])
	}
	return n, err
}

func (dr *DuplexReader) Write(p []byte) (int, error) {
	return dr.escritor.Write(p)
}
```

---

## 26.5 Closer Interface

### 26.5.1 Definición y Responsabilidad

```go
type Closer interface {
    Close() error
}
```

**Responsabilidades:**
- Liberar recursos (archivos, conexiones, memoria)
- Finalizar operaciones pendientes
- Retornar error si hay problemas cerrando
- Ser idempotente (seguro llamar múltiples veces)

### 26.5.2 Interfaces Compuestas con Closer

```go
type ReadCloser interface {
    Reader
    Closer
}

type WriteCloser interface {
    Writer
    Closer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

### 26.5.3 Patrón defer para Cierre Seguro

```go
package main

import (
	"io"
	"os"
)

func leerArchivoSeguro(ruta string) ([]byte, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	// Garantizar cierre incluso si hay pánico o retorno temprano
	defer file.Close()

	var resultado []byte
	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			resultado = append(resultado, buffer[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err // defer cierra el archivo
		}
	}

	return resultado, nil
} // file.Close() se ejecuta aquí
```

### 26.5.4 Capturar Errores de Close

```go
package main

import (
	"io"
	"os"
)

func escribirConValidacion(ruta string) error {
	file, err := os.Create(ruta)
	if err != nil {
		return err
	}

	// Capturar errores tanto de Write como de Close
	defer func() {
		if err := file.Close(); err != nil {
			println("Error al cerrar:", err.Error())
		}
	}()

	_, err = file.Write([]byte("Datos importantes"))
	return err
}

// Patrón mejorado con named return
func procesarConError(ruta string) (err error) {
	file, err := os.Open(ruta)
	if err != nil {
		return
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr // Propagar error de Close si no hay otro
		}
	}()

	// Procesamiento...
	return
}
```

### 26.5.5 Implementar Closer Personalizado

```go
package main

import (
	"fmt"
	"io"
)

// LimitedReader limitador con contador de lecturas
type LimitedReader struct {
	reader   io.Reader
	limite   int64
	leido    int64
	cerrado  bool
}

func NewLimitedReader(r io.Reader, limite int64) *LimitedReader {
	return &LimitedReader{
		reader: r,
		limite: limite,
	}
}

func (lr *LimitedReader) Read(p []byte) (int, error) {
	if lr.cerrado {
		return 0, io.ErrClosedPipe
	}

	if lr.leido >= lr.limite {
		return 0, io.EOF
	}

	// Limitar lectura al máximo permitido
	max := lr.limite - lr.leido
	if int64(len(p)) > max {
		p = p[:max]
	}

	n, err := lr.reader.Read(p)
	lr.leido += int64(n)

	return n, err
}

func (lr *LimitedReader) Close() error {
	fmt.Printf("Cerrado. Leído %d/%d bytes\n", lr.leido, lr.limite)
	lr.cerrado = true

	// Si el reader tiene Close, llamarlo
	if closer, ok := lr.reader.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
```

---

## 26.6 Copy y CopyN

### 26.6.1 io.Copy - Transferencia General

```go
func Copy(dst Writer, src Reader) (written int64, err error)
```

Copia datos de un Reader a un Writer hasta EOF:

```go
package main

import (
	"io"
	"os"
)

func copiarArchivo(origen, destino string) error {
	src, err := os.Open(origen)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy automáticamente usa buffering interno
	bytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}

	println("Copiados", bytes, "bytes")
	return nil
}

func main() {
	copiarArchivo("entrada.txt", "salida.txt")
}
```

**Caso: Descargar archivo de internet**

```go
package main

import (
	"io"
	"net/http"
	"os"
)

func descargarArchivo(url, rutaLocal string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(rutaLocal)
	if err != nil {
		return err
	}
	defer file.Close()

	// Usar io.Copy para transferir
	_, err = io.Copy(file, response.Body)
	return err
}

func main() {
	descargarArchivo("https://example.com/archivo.zip", "descargado.zip")
}
```

**Caso: Copiar con compresión**

```go
package main

import (
	"compress/gzip"
	"io"
	"os"
)

func comprimirArchivo(origen, destino string) error {
	src, _ := os.Open(origen)
	defer src.Close()

	dst, _ := os.Create(destino)
	defer dst.Close()

	// Crear writer comprimido
	writer := gzip.NewWriter(dst)
	defer writer.Close()

	// Copy automáticamente comprime
	_, err := io.Copy(writer, src)
	return err
}
```

### 26.6.2 io.CopyN - Copiar Cantidad Específica

```go
func CopyN(dst Writer, src Reader, n int64) (written int64, err error)
```

```go
package main

import (
	"io"
	"os"
)

func primerosNBytes(archivo string, n int64) ([]byte, error) {
	file, _ := os.Open(archivo)
	defer file.Close()

	// Leer solo los primeros n bytes
	return io.ReadAll(io.LimitReader(file, n))
}

func copiarLimitado(origen, destino string, maxBytes int64) error {
	src, _ := os.Open(origen)
	defer src.Close()

	dst, _ := os.Create(destino)
	defer dst.Close()

	// Copiar exactamente maxBytes (o EOF si es menos)
	written, err := io.CopyN(dst, src, maxBytes)
	if err != nil && err != io.EOF {
		return err
	}

	println("Se copiaron exactamente", written, "bytes")
	return nil
}
```

### 26.6.3 CopyBuffer - Control de Buffering

```go
package main

import (
	"io"
	"os"
)

func copiarConBufferPersonalizado() error {
	src, _ := os.Open("entrada.txt")
	defer src.Close()

	dst, _ := os.Create("salida.txt")
	defer dst.Close()

	// Buffer grande para archivos enormes
	buffer := make([]byte, 1024*1024) // 1 MB

	_, err := io.CopyBuffer(dst, src, buffer)
	return err
}

// Función personalizada con tracking
func copiarConProgreso(dst io.Writer, src io.Reader) error {
	buffer := make([]byte, 32*1024) // 32 KB
	totalCopiado := int64(0)

	for {
		n, err := src.Read(buffer)
		if n > 0 {
			written, writeErr := dst.Write(buffer[:n])
			totalCopiado += int64(written)

			// Mostrar progreso
			println("Copiados:", totalCopiado, "bytes")

			if writeErr != nil {
				return writeErr
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
```

---

## 26.7 Pipe

### 26.7.1 Concepto y Sintaxis

`io.Pipe()` crea un par de objetos conectados: Reader y Writer. Datos escritos en el Writer pueden ser leídos desde el Reader.

```go
func Pipe() (*PipeReader, *PipeWriter)
```

**Características:**
- Ambos extremos están en memoria (buffer pequeño)
- Write se bloquea si el Reader no mantiene ritmo
- Read espera si no hay datos
- Cerrar el Writer causa EOF en Reader
- Cerrar el Reader causa error en futuras escrituras

### 26.7.2 Comunicación Básica entre Goroutines

```go
package main

import (
	"fmt"
	"io"
)

func ejemploBasico() {
	reader, writer := io.Pipe()

	// Goroutine productor
	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Fprintf(writer, "Línea %d\n", i)
		}
		writer.Close() // Señalar fin
	}()

	// Lector principal
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Recibido: %s", string(buffer[:n]))
	}
}

func main() {
	ejemploBasico()
}
```

### 26.7.3 Patrón Productor-Consumidor Avanzado

```go
package main

import (
	"bufio"
	"fmt"
	"io"
)

// Procesador transforma datos pasando por pipe
type Procesador struct {
	nombre string
}

func (p *Procesador) Procesar(entrada io.Reader) io.Reader {
	lector, escritor := io.Pipe()

	go func() {
		defer escritor.Close()
		scanner := bufio.NewScanner(entrada)

		for scanner.Scan() {
			linea := scanner.Text()
			// Transformar datos
			transformada := fmt.Sprintf("[%s] %s\n", p.nombre, linea)
			fmt.Fprint(escritor, transformada)
		}
	}()

	return lector
}

func main() {
	// Crear cadena de procesadores
	entrada := strings.NewReader("Línea1\nLínea2\nLínea3\n")

	proc1 := &Procesador{nombre: "PROC1"}
	salida1 := proc1.Procesar(entrada)

	proc2 := &Procesador{nombre: "PROC2"}
	salida2 := proc2.Procesar(salida1)

	// Leer resultado final
	scanner := bufio.NewScanner(salida2)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	// [PROC2] [PROC1] Línea1
	// [PROC2] [PROC1] Línea2
	// [PROC2] [PROC1] Línea3
}
```

### 26.7.4 Pipe con Deadlock - Cómo Evitarlo

```go
package main

import (
	"io"
)

// ❌ DEADLOCK - Esto se bloquea
func deadlock() {
	r, w := io.Pipe()
	defer w.Close()
	defer r.Close()

	// Escribir en el mismo goroutine antes de leer
	// Si el buffer del pipe se llena, se bloquea para siempre
	for i := 0; i < 1000000; i++ {
		w.Write([]byte("dato"))
		// r.Read() nunca se ejecuta, hay deadlock
	}
}

// ✅ CORRECTO - Usar goroutines
func correcto() {
	r, w := io.Pipe()
	defer w.Close()
	defer r.Close()

	// Goroutine para escribir
	go func() {
		for i := 0; i < 100; i++ {
			w.Write([]byte("dato"))
		}
		w.Close()
	}()

	// Leer en paralelo
	buffer := make([]byte, 4)
	for {
		_, err := r.Read(buffer)
		if err == io.EOF {
			break
		}
	}
}
```

---

## 26.8 MultiReader y MultiWriter

### 26.8.1 MultiReader - Concatenar Lecturas

```go
func MultiReader(readers ...Reader) Reader
```

Crea un Reader que secuencialmente lee de múltiples Readers:

```go
package main

import (
	"fmt"
	"io"
	"strings"
)

func ejemploMultiReader() {
	// Tres fuentes de lectura
	r1 := strings.NewReader("Primera parte. ")
	r2 := strings.NewReader("Segunda parte. ")
	r3 := strings.NewReader("Tercera parte.")

	// Crear lector multi
	multi := io.MultiReader(r1, r2, r3)

	// Leer como si fuera una sola fuente
	buffer := make([]byte, 1024)
	n, _ := multi.Read(buffer)

	fmt.Println(string(buffer[:n]))
	// Salida: Primera parte. Segunda parte. Tercera parte.
}
```

**Caso: Unir múltiples archivos**

```go
package main

import (
	"io"
	"os"
)

func unirArchivos(archivos []string, salida string) error {
	readers := make([]io.Reader, len(archivos))

	// Abrir todos los archivos
	for i, archivo := range archivos {
		file, _ := os.Open(archivo)
		defer file.Close()
		readers[i] = file
	}

	// Usar MultiReader para leer todos secuencialmente
	multi := io.MultiReader(readers...)

	dst, _ := os.Create(salida)
	defer dst.Close()

	_, err := io.Copy(dst, multi)
	return err
}

func main() {
	unirArchivos([]string{"parte1.txt", "parte2.txt", "parte3.txt"}, "completo.txt")
}
```

### 26.8.2 MultiWriter - Escribir en Múltiples Destinos

```go
func MultiWriter(writers ...Writer) Writer
```

Crea un Writer que escribe a múltiples Writers:

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func ejemploMultiWriter() {
	var buf1, buf2 bytes.Buffer
	archivo, _ := os.Create("log.txt")
	defer archivo.Close()

	// Escribir a 3 destinos simultáneamente
	multi := io.MultiWriter(&buf1, &buf2, archivo)

	fmt.Fprintf(multi, "Esta línea va a 3 lugares\n")

	fmt.Println("Buffer 1:", buf1.String())
	fmt.Println("Buffer 2:", buf2.String())
	// También en archivo log.txt
}
```

**Caso: Logging a múltiples destinos**

```go
package main

import (
	"io"
	"os"
	"time"
)

type Logger struct {
	writer io.Writer
}

func NewLogger(writers ...io.Writer) *Logger {
	return &Logger{
		writer: io.MultiWriter(writers...),
	}
}

func (l *Logger) Log(mensaje string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	io.WriteString(l.writer, timestamp+" "+mensaje+"\n")
}

func main() {
	file, _ := os.Create("app.log")
	defer file.Close()

	// Log a stdout, stderr y archivo simultáneamente
	logger := NewLogger(os.Stdout, os.Stderr, file)

	logger.Log("Aplicación iniciada")
	logger.Log("Procesando datos")
	logger.Log("Aplicación finalizada")
}
```

### 26.8.3 Combinar MultiReader + MultiWriter

```go
package main

import (
	"io"
	"strings"
)

func transformarMultiples() {
	// Múltiples fuentes
	r1 := strings.NewReader("Entrada1 ")
	r2 := strings.NewReader("Entrada2 ")

	// Múltiples destinos
	var salida1, salida2 strings.Builder

	src := io.MultiReader(r1, r2)
	dst := io.MultiWriter(&salida1, &salida2)

	// Copiar de múltiples fuentes a múltiples destinos
	io.Copy(dst, src)

	println("Salida 1:", salida1.String())
	println("Salida 2:", salida2.String())
	// Ambos: "Entrada1 Entrada2 "
}
```

---

## 26.9 ReadAll y ReadAtLeast

### 26.9.1 ReadAll - Leer Todo el Contenido

```go
func ReadAll(r Reader) ([]byte, error)
```

Lee todo hasta EOF desde un Reader:

```go
package main

import (
	"io"
	"os"
)

func leerTodoArchivo(ruta string) ([]byte, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Leer el archivo completo en memoria
	contenido, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	println("Bytes leídos:", len(contenido))
	return contenido, nil
}

func main() {
	datos, _ := leerTodoArchivo("documento.txt")
	println("Contenido:", string(datos))
}
```

**Advertencia de seguridad:**

```go
package main

import (
	"io"
	"net/http"
)

func downloadConSeguridad(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// ❌ INSEGURO - Podría descargar 10 GB a memoria
	// datos, _ := io.ReadAll(resp.Body)

	// ✅ SEGURO - Limitar tamaño
	limitado := io.LimitReader(resp.Body, 10*1024*1024) // 10 MB máx
	datos, err := io.ReadAll(limitado)

	return datos, err
}
```

### 26.9.2 ReadAtLeast - Leer Mínimo de Bytes

```go
func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error)
```

Lee hasta llenar el buffer o completar mínimo requerido:

```go
package main

import (
	"fmt"
	"io"
	"strings"
)

func ejemploReadAtLeast() {
	reader := strings.NewReader("Hola Mundo")
	buffer := make([]byte, 100)

	// Leer al menos 5 bytes
	n, err := io.ReadAtLeast(reader, buffer, 5)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Leído %d bytes: %s\n", n, string(buffer[:n]))
	// Leído 11 bytes: Hola Mundo
}

// Caso: Validar protocolo
func leerCabecera() error {
	// Asegurar que el protocolo tiene al menos 8 bytes
	buffer := make([]byte, 1024)
	n, err := io.ReadAtLeast(os.Stdin, buffer, 8)

	if err == io.ErrShortBuffer {
		println("Cabecera muy corta")
		return err
	}

	if err == io.EOF {
		println("Conexión cerrada antes de cabecera completa")
		return err
	}

	println("Cabecera válida:", string(buffer[:n]))
	return nil
}
```

### 26.9.3 Patrones de Lectura con Validación

```go
package main

import (
	"io"
	"os"
)

func leerConValidacion(r io.Reader, minBytes, maxBytes int64) ([]byte, error) {
	// Limitar tamaño máximo
	limitado := io.LimitReader(r, maxBytes)

	datos, err := io.ReadAll(limitado)
	if err != nil {
		return nil, err
	}

	// Validar mínimo
	if int64(len(datos)) < minBytes {
		return nil, io.ErrUnexpectedEOF
	}

	return datos, nil
}

func main() {
	// Leer entre 10 bytes y 1 MB
	datos, err := leerConValidacion(os.Stdin, 10, 1024*1024)
	if err != nil {
		panic(err)
	}
	println("Datos válidos:", len(datos))
}
```

---

## 26.10 Buffering con bufio

### 26.10.1 bufio.Reader - Lectura Buffered

```go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func ejemploBufferedReader() {
	file, _ := os.Open("datos.txt")
	defer file.Close()

	// Crear reader buffered con buffer de 64 KB
	reader := bufio.NewReaderSize(file, 64*1024)

	// Leer línea a línea (eficiente)
	for {
		linea, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(linea)
	}
}

// Alternativa: Usar Scanner para mejor manejo de errores
func ejemploScanner() {
	file, _ := os.Open("datos.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
```

### 26.10.2 bufio.Writer - Escritura Buffered

```go
package main

import (
	"bufio"
	"os"
)

func ejemploBufferedWriter() {
	file, _ := os.Create("salida.txt")
	defer file.Close()

	// Writer buffered con buffer de 256 KB
	writer := bufio.NewWriterSize(file, 256*1024)

	// Escribir múltiples veces sin syscalls
	for i := 0; i < 1000; i++ {
		writer.WriteString("Línea de datos\n")
	}

	// Importante: Flush al terminar
	err := writer.Flush()
	if err != nil {
		panic(err)
	}
}

// Patrón con deferred flush
func escribirConFlush(ruta string) error {
	file, _ := os.Create(ruta)
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer func() {
		if err := writer.Flush(); err != nil {
			println("Error al flush:", err.Error())
		}
	}()

	// Escribir datos...
	writer.WriteString("Dato 1\n")
	writer.WriteString("Dato 2\n")

	return nil
}
```

### 26.10.3 bufio.ReadWriter - Bidireccional

```go
package main

import (
	"bufio"
	"io"
)

func ejemploReadWriter(rw io.ReadWriter) {
	// Crear reader y writer buffered
	buffered := bufio.NewReadWriter(
		bufio.NewReader(rw),
		bufio.NewWriter(rw),
	)

	// Leer línea
	linea, _ := buffered.ReadString('\n')
	println("Recibido:", linea)

	// Escribir respuesta
	buffered.WriteString("Respuesta\n")
	buffered.Flush()
}
```

### 26.10.4 Optimización de Buffer Size

```go
package main

import (
	"bufio"
	"io"
	"os"
	"time"
)

// Función para medir rendimiento
func benchmark(nombre string, bufferSize int, iter int) time.Duration {
	file, _ := os.Create("/tmp/test.txt")
	defer file.Close()

	reader := bufio.NewReaderSize(file, bufferSize)

	inicio := time.Now()
	for i := 0; i < iter; i++ {
		reader.ReadString('\n')
	}
	duracion := time.Since(inicio)

	println(nombre, ":", duracion.String())
	return duracion
}

func main() {
	// Comparar diferentes buffer sizes
	benchmark("Buffer 4KB", 4*1024, 100000)
	benchmark("Buffer 32KB", 32*1024, 100000)
	benchmark("Buffer 256KB", 256*1024, 100000)

	// Generalmente 32-64 KB es óptimo para disk I/O
}
```

---

## 26.11 Buenas Prácticas y Patrones

### 26.11.1 Arquitectura de I/O en Go

**Principios:**
1. Siempre usar Reader/Writer, nunca tipos concretos
2. Funciones aceptan interfaces, no structs
3. Cerrar recursos con defer
4. Verificar errores siempre

```go
package main

import (
	"io"
	"os"
)

// ❌ Evitar: tomar tipo concreto
func procesarMal(f *os.File) error {
	// Acoplado a *os.File
	return nil
}

// ✅ Correcto: aceptar interfaz
func procesarBien(r io.Reader) error {
	// Funciona con archivos, sockets, buffers, etc.
	return nil
}

// ✅ Mejor aún: separar responsabilidades
func procesarOptimo(r io.ReadCloser) (err error) {
	defer func() {
		if closeErr := r.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	// Procesamiento...
	return nil
}
```

### 26.11.2 Antipatrones Comunes

**Antipatrón 1: Ignorar errores parciales**

```go
// ❌ INCORRECTO
func copiarMal(dst io.Writer, src io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		dst.Write(buf[:n]) // ¿Y si Write falla parcialmente?
		if err != nil {
			break
		}
	}
}

// ✅ CORRECTO
func copiarBien(dst io.Writer, src io.Reader) error {
	return io.CopyBuffer(dst, src, make([]byte, 1024))
}
```

**Antipatrón 2: Buffer size inadecuado**

```go
// ❌ Muy pequeño - muchos syscalls
reader := bufio.NewReaderSize(file, 256) // 256 bytes

// ❌ Muy grande - desperdicia memoria
reader := bufio.NewReaderSize(file, 100*1024*1024) // 100 MB

// ✅ Óptimo - típicamente 32-64 KB
reader := bufio.NewReaderSize(file, 32*1024)
```

**Antipatrón 3: No validar tamaños**

```go
// ❌ INSEGURO - DoS posible
func servirArchivo(r io.Reader, w io.Writer) {
	io.Copy(w, r) // ¿Qué si es 10 GB?
}

// ✅ SEGURO - Limitar tamaño
func servirArchivoSeguro(r io.Reader, w io.Writer) error {
	limited := io.LimitReader(r, 100*1024*1024) // 100 MB máx
	_, err := io.Copy(w, limited)
	return err
}
```

### 26.11.3 Comparación: Lectura Eficiente vs Ineficiente

```go
package main

import (
	"bufio"
	"io"
	"os"
	"time"
)

// ❌ INEFICIENTE: un byte a la vez
func leerIneficiente(archivo string) {
	file, _ := os.Open(archivo)
	defer file.Close()

	for {
		b := make([]byte, 1)
		n, _ := file.Read(b)
		if n == 0 {
			break
		}
		// Procesar byte...
	}
}

// ⚠️ MEDIO: ReadAll sin límite
func leerMedio(archivo string) {
	file, _ := os.Open(archivo)
	defer file.Close()

	datos, _ := io.ReadAll(file)
	// Procesar todos los datos...
	_ = datos
}

// ✅ EFICIENTE: Scanner con buffering
func leerEficiente(archivo string) {
	file, _ := os.Open(archivo)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Procesar línea...
		_ = scanner.Text()
	}
}

// ✅ MÁS EFICIENTE: ReadAtLeast con buffer manual
func leerMasEficiente(archivo string) {
	file, _ := os.Open(archivo)
	defer file.Close()

	reader := bufio.NewReaderSize(file, 64*1024)
	buffer := make([]byte, 64*1024)

	for {
		n, err := reader.Read(buffer)
		if n == 0 && err == io.EOF {
			break
		}
		// Procesar buffer...
		_ = buffer[:n]
	}
}
```

### 26.11.4 Patrón: Middleware de I/O

```go
package main

import (
	"compress/gzip"
	"crypto/md5"
	"io"
	"os"
)

// WriterMiddleware permite encadenar transformaciones
type WriterMiddleware func(io.Writer) io.Writer

// Middleware de compresión
func Compress() WriterMiddleware {
	return func(w io.Writer) io.Writer {
		return gzip.NewWriter(w)
	}
}

// Middleware de hash
type HashWriter struct {
	dst  io.Writer
	hash io.Writer
}

func (hw *HashWriter) Write(p []byte) (int, error) {
	hw.hash.Write(p) // Calcular hash
	return hw.dst.Write(p) // Escribir destino
}

func Hash() WriterMiddleware {
	return func(w io.Writer) io.Writer {
		return &HashWriter{
			dst:  w,
			hash: md5.New(),
		}
	}
}

// Usar middleware
func procesarConMiddleware(src io.Reader) error {
	file, _ := os.Create("salida.gz")
	defer file.Close()

	// Aplicar middleware
	writer := file
	for _, mw := range []WriterMiddleware{Compress()} {
		writer = mw(writer)
	}

	_, err := io.Copy(writer, src)
	return err
}
```

### 26.11.5 Patrones de Error Handling

```go
package main

import (
	"errors"
	"io"
)

// Tipo personalizado para errores de I/O
type IOError struct {
	Op   string
	Path string
	Err  error
}

func (e *IOError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

// Usar en lectura
func leerConErroresMejor(ruta string) ([]byte, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, &IOError{"open", ruta, err}
	}
	defer file.Close()

	datos, err := io.ReadAll(file)
	if err != nil {
		return nil, &IOError{"read", ruta, err}
	}

	return datos, nil
}

// Envolver errores específicos
func manejarError(err error) {
	if errors.Is(err, io.EOF) {
		println("Fin de archivo")
	} else if errors.Is(err, io.ErrShortWrite) {
		println("Escritura incompleta")
	} else if errors.Is(err, io.ErrUnexpectedEOF) {
		println("EOF inesperado")
	}
}
```

---

# EJERCICIOS PROGRESIVOS

## Ejercicio 26.1: Leer Archivo con Reader

**Objetivo:** Implementar lectura completa de archivo usando Reader interface.

**Requisitos:**
- Abrir archivo especificado como argumento
- Leer en chunks de 4 KB
- Mostrar progreso en MB/s
- Validar sum SHA256 del contenido

```go
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	archivo := flag.String("file", "", "Archivo a leer")
	flag.Parse()

	if *archivo == "" {
		fmt.Println("Uso: go run solucion.go -file <archivo>")
		return
	}

	file, err := os.Open(*archivo)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// TODO: Implementar lectura con Reader
	// - Crear hash SHA256
	// - Leer en chunks de 4 KB
	// - Mostrar velocidad de lectura
	// - Imprimir hash final
}
```

---

## Ejercicio 26.2: Copiar y Validar Archivo

**Objetivo:** Implementar copia de archivo con validación de integridad.

**Requisitos:**
- Copiar archivo origen a destino
- Mostrar progreso durante copia
- Validar que se copiaron todos los bytes
- Detectar errores de escritura parcial

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	origen := flag.String("src", "", "Archivo origen")
	destino := flag.String("dst", "", "Archivo destino")
	flag.Parse()

	if *origen == "" || *destino == "" {
		fmt.Println("Uso: go run solucion.go -src <origen> -dst <destino>")
		return
	}

	// TODO: Implementar
	// - Abrir archivos
	// - Usar io.Copy
	// - Verificar coincidencia de tamaños
	// - Manejar errores de permisos/espacio
}
```

---

## Ejercicio 26.3: Pipe entre Goroutines

**Objetivo:** Implementar comunicación entre goroutines usando io.Pipe.

**Requisitos:**
- Productor genera números del 1-100
- Pipe transporta los datos
- Consumidor suma todos los números
- Usar canales para sincronización adicional

```go
package main

import (
	"fmt"
	"io"
)

func main() {
	// TODO: Implementar
	// - Crear pipe con io.Pipe()
	// - Goroutine productor que escribe números
	// - Goroutine consumidor que los lee y suma
	// - Mostrar suma final
	// - Manejar cierre seguro
}
```

---

## Ejercicio 26.4: Implementar Custom Reader

**Objetivo:** Crear un Reader personalizado que genere datos computados.

**Requisitos:**
- Reader que genera secuencia Fibonacci
- Limitar a N números máximo
- Implementar cierre graceful
- Validar contra implementación io.Reader

```go
package main

import (
	"io"
)

// FibonacciReader genera números Fibonacci
type FibonacciReader struct {
	// TODO: Definir campos
}

func NewFibonacciReader(max int) *FibonacciReader {
	// TODO: Implementar
	return nil
}

func (fr *FibonacciReader) Read(p []byte) (int, error) {
	// TODO: Implementar generación de Fibonacci
	// - Serializar como texto
	// - Retornar bytes leídos
	// - Retornar EOF cuando alcance máximo
	return 0, nil
}

func main() {
	// TODO: Usar FibonacciReader
	// - Leer todos los números
	// - Verificar valores correctos
	// - Validar EOF se retorna correctamente
}
```

---

## Ejercicio 26.5: Buffered I/O Optimizado

**Objetivo:** Implementar pipeline de lectura/escritura optimizado.

**Requisitos:**
- Leer archivo grande línea a línea
- Transformar (convertir a mayúsculas)
- Escribir a archivo salida
- Optimizar con buffer sizes apropiados
- Comparar rendimiento con diferentes buffer sizes

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func procesarConBuffer(entrada, salida string, bufSize int) (time.Duration, error) {
	inicio := time.Now()

	// TODO: Implementar
	// - Abrir entrada con bufio.NewReaderSize(bufSize)
	// - Leer línea a línea
	// - Convertir a mayúsculas
	// - Escribir con bufio.NewWriterSize(bufSize)
	// - Flush al terminar
	// - Retornar duración

	return time.Since(inicio), nil
}

func main() {
	entrada := flag.String("in", "", "Archivo entrada")
	salida := flag.String("out", "", "Archivo salida")
	flag.Parse()

	// TODO: Comparar rendimiento con:
	// - Buffer 4 KB
	// - Buffer 32 KB
	// - Buffer 256 KB
	// - Mostrar cuál es más rápido
}
```

---

## RESUMEN

El package `io` en Go proporciona abstracciones elegantes y eficientes para operaciones de entrada/salida:

**Puntos clave:**
- Interfaces pequeñas y composables (Reader, Writer, Closer)
- Desacoplamiento entre fuentes y destinos
- Funciones utilitarias (Copy, Pipe, MultiReader/Writer)
- Integración con buffering mediante `bufio`
- Manejo seguro de recursos con `defer`

**Cuándo usar qué:**
- **Reader**: para cualquier fuente de datos
- **Writer**: para cualquier destino
- **io.Copy**: transferencia simple de datos
- **io.Pipe**: comunicación entre goroutines
- **bufio**: optimizar I/O en bucles
- **MultiReader/Writer**: componer múltiples fuentes/destinos

El diseño de interfaces en Go permite escribir código genérico, testeable y eficiente que funciona con cualquier implementador sin cambios.
