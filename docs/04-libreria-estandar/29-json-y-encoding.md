# Capítulo 29: JSON y encoding - Serialización de datos

## Índice
1. [JSON Encoding Fundamentos](#291-json-encoding-fundamentos)
2. [Marshal y MarshalIndent](#292-marshal-y-marshalindent)
3. [Unmarshal](#293-unmarshal)
4. [JSON Tags](#294-json-tags)
5. [Custom Types](#295-custom-types)
6. [Number Handling](#296-number-handling)
7. [Streaming JSON](#297-streaming-json)
8. [Validation](#298-validation)
9. [XML y CSV](#299-xml-y-csv)
10. [Base64 y Hex](#2910-base64-y-hex)
11. [Buenas Prácticas y Patrones](#2911-buenas-prácticas-y-patrones)

---

## 29.1 JSON Encoding Fundamentos

### 29.1.1 Introducción a JSON en Go

JSON (JavaScript Object Notation) es el formato más utilizado para intercambio de datos en APIs REST modernas. Go proporciona soporte nativo a través del paquete `encoding/json`.

**Características clave:**
- Serialización bidireccional (Marshal/Unmarshal)
- Mapeo automático entre tipos Go y JSON
- Soporte para tipos primitivos, estructuras y interfaces
- Streaming para archivos grandes
- Validación integrada de tipos

**Tipos soportados en Go:**
```
JSON null     → Go nil
JSON boolean  → Go bool
JSON number   → Go float64, int, int64, etc.
JSON string   → Go string
JSON array    → Go slice
JSON object   → Go map o struct
```

### 29.1.2 Especificación JSON y Tipos Go

**Mapeo de tipos fundamental:**

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Ejemplo struct {
	// Tipos soportados directamente
	Nombre      string    `json:"nombre"`           // string → JSON string
	Edad        int       `json:"edad"`             // int → JSON number
	Altura      float64   `json:"altura"`           // float64 → JSON number
	Activo      bool      `json:"activo"`           // bool → JSON boolean
	Etiquetas   []string  `json:"etiquetas"`        // []string → JSON array
	Metadata    map[string]interface{} `json:"metadata"` // map → JSON object
	Opcional    *string   `json:"opcional"`         // *string → JSON string o null
}

func main() {
	ejemplo := Ejemplo{
		Nombre:    "Juan",
		Edad:      30,
		Altura:    1.75,
		Activo:    true,
		Etiquetas: []string{"golang", "json"},
		Metadata: map[string]interface{}{
			"version": 1.0,
			"creado":  "2024-01-15",
		},
		Opcional: nil,
	}

	// Marshal a JSON
	jsonBytes, err := json.Marshal(ejemplo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("JSON: %s\n", string(jsonBytes))
	// Salida: JSON: {"nombre":"Juan","edad":30,"altura":1.75,"activo":true,"etiquetas":["golang","json"],"metadata":{"version":1,"creado":"2024-01-15"},"opcional":null}
}
```

### 29.1.3 Valores por Defecto y Inicialización

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Host     string
	Puerto   int
	Debug    bool
	Timeout  float64
	Tags     []string
	Extras   map[string]interface{}
}

func main() {
	// Valores por defecto de tipos Go
	var config Config
	fmt.Printf("Config vacío: %+v\n", config)
	// Config vacío: {Host: Puerto:0 Debug:false Timeout:0 Tags:[] Extras:map[]}

	// Al deserializar JSON incompleto, mantiene valores por defecto
	jsonData := []byte(`{"Host":"localhost"}`)
	err := json.Unmarshal(jsonData, &config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Config parcial: %+v\n", config)
	// Config parcial: {Host:localhost Puerto:0 Debug:false Timeout:0 Tags:[] Extras:map[]}
}
```

### 29.1.4 Comparación con Otros Lenguajes

**Java (Jackson):**
```java
ObjectMapper mapper = new ObjectMapper();
User user = new User("Juan", 30);
String json = mapper.writeValueAsString(user);
User deserialized = mapper.readValue(json, User.class);
```

**Python (json.dumps):**
```python
import json
user = {"name": "Juan", "age": 30}
json_str = json.dumps(user)
user_obj = json.loads(json_str)
```

**Go:**
```go
var user User
json.Unmarshal([]byte(`{"name":"Juan","age":30}`), &user)
jsonBytes, _ := json.Marshal(user)
```

**Diferencias:**
- Go requiere pasar punteros al Unmarshal
- Go necesita tipos estáticos (structs)
- Go es más seguro en tipos pero menos flexible
- Python es más dinámico, Java más verboso

---

## 29.2 Marshal y MarshalIndent

### 29.2.1 Marshal Básico

`Marshal` convierte una estructura Go a JSON compacto (sin espacios).

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Libro struct {
	Titulo  string
	Autor   string
	Paginas int
	Precio  float64
}

func main() {
	libro := Libro{
		Titulo:  "Clean Code",
		Autor:   "Robert Martin",
		Paginas: 464,
		Precio:  49.99,
	}

	// Marshal simple
	jsonBytes, err := json.Marshal(libro)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Salida: {"Titulo":"Clean Code","Autor":"Robert Martin","Paginas":464,"Precio":49.99}
}
```

**Características:**
- Retorna `[]byte` y `error`
- Todos los campos se capitalizan (públicos en Go)
- JSON es compacto (sin espacios innecesarios)
- Rápido y eficiente

### 29.2.2 MarshalIndent para Formato Legible

```go
package main

import (
	"encoding/json"
	"fmt"
)

type API struct {
	Version string `json:"version"`
	Status  string `json:"status"`
	Data    map[string]interface{} `json:"data"`
}

func main() {
	api := API{
		Version: "1.0",
		Status:  "activo",
		Data: map[string]interface{}{
			"usuarios": 1524,
			"uptime":   "99.9%",
		},
	}

	// MarshalIndent con 2 espacios
	jsonBytes, err := json.MarshalIndent(api, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	/* Salida:
	{
	  "version": "1.0",
	  "status": "activo",
	  "data": {
	    "usuarios": 1524,
	    "uptime": "99.9%"
	  }
	}
	*/
}
```

**Parámetros:**
- Primer parámetro: prefijo (ej: "", ">" para cada línea)
- Segundo parámetro: indentación (ej: "  ", "\t")

### 29.2.3 Marshaling de Colecciones

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// Slice de estructuras
	usuarios := []map[string]interface{}{
		{"id": 1, "nombre": "Alice", "admin": true},
		{"id": 2, "nombre": "Bob", "admin": false},
	}

	jsonBytes, _ := json.MarshalIndent(usuarios, "", "  ")
	fmt.Println(string(jsonBytes))
	/* Salida:
	[
	  {
	    "admin": true,
	    "id": 1,
	    "nombre": "Alice"
	  },
	  {
	    "admin": false,
	    "id": 2,
	    "nombre": "Bob"
	  }
	]
	*/

	// Map simple
	config := map[string]interface{}{
		"puerto":     8080,
		"host":       "localhost",
		"debug":      true,
		"timeout_ms": 5000,
	}

	jsonBytes, _ = json.MarshalIndent(config, "", "  ")
	fmt.Println(string(jsonBytes))
}
```

### 29.2.4 Errores Comunes en Marshal

```go
package main

import (
	"encoding/json"
	"fmt"
	"math"
)

func main() {
	// Error: canal no serializable
	ch := make(chan int)
	_, err := json.Marshal(ch)
	fmt.Println("Canal:", err) 
	// Error: json: unsupported type: chan int

	// Error: función no serializable
	fn := func() {}
	_, err = json.Marshal(fn)
	fmt.Println("Función:", err)
	// Error: json: unsupported type: func()

	// Error: NaN y Infinity
	valores := map[string]interface{}{
		"infinito": math.Inf(1),
		"nan":      math.NaN(),
	}
	_, err = json.Marshal(valores)
	fmt.Println("NaN/Inf:", err)
	// Error: json: unsupported value: Infinity

	// Nota: nil se serializa correctamente como "null"
	var ptr *string
	_, err = json.Marshal(ptr)
	fmt.Println("Nil error:", err) // nil error: <nil>
}
```

---

## 29.3 Unmarshal

### 29.3.1 Deserialización Básica

`Unmarshal` convierte JSON a estructuras Go.

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Producto struct {
	ID       int
	Nombre   string
	Precio   float64
	Disponible bool
}

func main() {
	jsonData := []byte(`{
		"ID": 101,
		"Nombre": "Laptop",
		"Precio": 999.99,
		"Disponible": true
	}`)

	var producto Producto
	err := json.Unmarshal(jsonData, &producto)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Producto: %+v\n", producto)
	// Producto: {ID:101 Nombre:Laptop Precio:999.99 Disponible:true}
}
```

**Reglas importantes:**
- REQUIERE un puntero (&) al destino
- Case-sensitive: el campo Go debe coincidir o usar tags
- Ignora campos JSON no mapeados
- Ignora campos Go sin equivalente en JSON

### 29.3.2 Type Conversion Automática

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Numeros struct {
	EnteroFromFloat float64 `json:"entero_float"`
	FloatFromInt    float64 `json:"float_int"`
	StringToInt     int     `json:"string_int"`
}

func main() {
	// JSON con tipos "incorrectos"
	jsonData := []byte(`{
		"entero_float": 42.0,
		"float_int": 100,
		"string_int": "999"
	}`)

	var nums Numeros
	err := json.Unmarshal(jsonData, &nums)
	
	// Conversión de tipos automática funciona
	fmt.Printf("%+v\n", nums)
	// {EnteroFromFloat:42 FloatFromInt:100 StringToInt:0}

	// Pero string a int falla (StringToInt será 0)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

### 29.3.3 Parcial Unmarshal

```go
package main

import (
	"encoding/json"
	"fmt"
)

type UsuarioCompleto struct {
	ID       int
	Nombre   string
	Email    string
	Edad     int
	Telefono string
}

type UsuarioMinimo struct {
	ID     int
	Nombre string
}

func main() {
	jsonData := []byte(`{
		"ID": 1,
		"Nombre": "Carlos",
		"Email": "carlos@example.com",
		"Edad": 28,
		"Telefono": "+34612345678"
	}`)

	// Deserializar solo en los campos que existen en la struct
	var usuario UsuarioMinimo
	err := json.Unmarshal(jsonData, &usuario)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Usuario mínimo: %+v\n", usuario)
	// Usuario mínimo: {ID:1 Nombre:Carlos}
	// Los campos Email, Edad, Telefono se ignoran

	// Útil para APIs que devuelven más datos de los que necesitas
}
```

### 29.3.4 Errores de Unmarshal

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Persona struct {
	Nombre string
	Edad   int
}

func main() {
	// JSON inválido
	_, err := json.Unmarshal([]byte(`{invalid}`), &Persona{})
	fmt.Println("JSON inválido:", err)
	// Error: invalid character 'i' looking for beginning of object key string

	// Tipo incorrecto
	datos := struct {
		Edad string
	}{}
	err = json.Unmarshal([]byte(`{"Edad": 30}`), &datos)
	fmt.Println("Tipo string cuando es number:", err)
	// Error: json: cannot unmarshal number into Go struct field .Edad of type string

	// No es un puntero
	// err := json.Unmarshal([]byte(`{}`), Persona{}) // COMPILARÍA UN ERROR
}
```

---

## 29.4 JSON Tags

### 29.4.1 Estructura de Tags

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Usuario struct {
	// Tag básico: renombra el campo en JSON
	ID       int    `json:"id"`
	
	// omitempty: omite si el valor es cero
	Nombre   string `json:"nombre,omitempty"`
	
	// ignore: ignora completamente este campo
	Contraseña string `json:"-"`
	
	// string: convierte a string en JSON
	Cantidad int64  `json:"cantidad,string"`
	
	// Combinado: múltiples opciones
	Actualizado *string `json:"actualizado,omitempty"`
}

func main() {
	usuario := Usuario{
		ID:       1,
		Nombre:   "Elena",
		Contraseña: "secreto123",
		Cantidad: 500,
	}

	jsonBytes, _ := json.MarshalIndent(usuario, "", "  ")
	fmt.Println(string(jsonBytes))
	/* Salida:
	{
	  "id": 1,
	  "nombre": "Elena",
	  "cantidad": "500"
	}
	*/
	// Nota: Contraseña y Actualizado se omiten
}
```

### 29.4.2 omitempty en Detalle

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Solicitud struct {
	// Siempre incluye
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	
	// Incluye solo si tiene valor
	Descripcion string `json:"descripcion,omitempty"`
	Prioridad   int    `json:"prioridad,omitempty"`
	Etiquetas   []string `json:"etiquetas,omitempty"`
	Asignado    *string `json:"asignado,omitempty"`
}

func main() {
	// Solicitud con algunos campos vacíos
	solicitud := Solicitud{
		ID:     42,
		Titulo: "Implementar cache",
		// Descripcion vacío (string zero value)
		// Prioridad 0 (int zero value)
		Etiquetas: nil, // nil slice
		Asignado:  nil, // nil pointer
	}

	jsonBytes, _ := json.MarshalIndent(solicitud, "", "  ")
	fmt.Println(string(jsonBytes))
	/* Salida:
	{
	  "id": 42,
	  "titulo": "Implementar cache"
	}
	*/
	// Los campos vacíos se omiten por omitempty

	// Slice vacío vs nil
	solicitud2 := Solicitud{
		ID:        43,
		Titulo:    "Otro",
		Etiquetas: []string{}, // Empty slice
	}

	jsonBytes, _ = json.MarshalIndent(solicitud2, "", "  ")
	fmt.Println(string(jsonBytes))
	// Slice vacío se omite también con omitempty
}
```

### 29.4.3 String Tag y Conversiones

```go
package main

import (
	"encoding/json"
	"fmt"
)

type API struct {
	// string: serializa como string incluso si es número
	RequestID   int64  `json:"request_id,string"`
	Timestamp   int64  `json:"timestamp,string"`
	Success     bool   `json:"success,string"` // "true" o "false"
	NumeroFloor float64 `json:"numero,string"`
}

func main() {
	api := API{
		RequestID:   9223372036854775807, // Max int64
		Timestamp:   1705334400,
		Success:     true,
		NumeroFloor: 3.14159,
	}

	jsonBytes, _ := json.MarshalIndent(api, "", "  ")
	fmt.Println(string(jsonBytes))
	/* Salida:
	{
	  "numero": "3.14159",
	  "request_id": "9223372036854775807",
	  "success": "true",
	  "timestamp": "1705334400"
	}
	*/

	// Al deserializar, se convierte automáticamente
	jsonData := []byte(`{"request_id":"12345","timestamp":"1234567890","success":"false","numero":"2.71"}`)
	var api2 API
	json.Unmarshal(jsonData, &api2)
	fmt.Printf("%+v\n", api2)
	// {RequestID:12345 Timestamp:1234567890 Success:false NumeroFloor:2.71}
}
```

### 29.4.4 Renaming Fields

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Empleado struct {
	// Go uses camelCase, JSON uses snake_case
	NombreCompleto string `json:"nombre_completo"`
	NumeroDocumento string `json:"numero_documento"`
	FechaContratacion string `json:"fecha_contratacion"`
	EsActivo bool `json:"es_activo"`
	
	// Puede haber campos sin mapeo
	DNI string `json:"-"` // No aparece en JSON
}

func main() {
	empleado := Empleado{
		NombreCompleto: "María García",
		NumeroDocumento: "12345678A",
		FechaContratacion: "2023-01-15",
		EsActivo: true,
		DNI: "12345678",
	}

	jsonBytes, _ := json.MarshalIndent(empleado, "", "  ")
	fmt.Println(string(jsonBytes))
	/* Salida:
	{
	  "es_activo": true,
	  "fecha_contratacion": "2023-01-15",
	  "nombre_completo": "María García",
	  "numero_documento": "12345678A"
	}
	*/
}
```

---

## 29.5 Custom Types

### 29.5.1 Implementar Marshaler

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Persona struct {
	Nombre string
	Nacimiento time.Time
}

// Implementar json.Marshaler
func (p Persona) MarshalJSON() ([]byte, error) {
	type Alias Persona
	return json.Marshal(&struct {
		Nacimiento string `json:"nacimiento"` // Custom format
		*Alias
	}{
		Nacimiento: p.Nacimiento.Format("2006-01-02"),
		Alias:      (*Alias)(&p),
	})
}

func main() {
	persona := Persona{
		Nombre:     "Jorge",
		Nacimiento: time.Date(1990, time.January, 15, 0, 0, 0, 0, time.UTC),
	}

	jsonBytes, _ := json.MarshalIndent(persona, "", "  ")
	fmt.Println(string(jsonBytes))
	/* Salida:
	{
	  "Nombre": "Jorge",
	  "nacimiento": "1990-01-15"
	}
	*/
}
```

### 29.5.2 Implementar Unmarshaler

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Evento struct {
	Titulo string
	Fecha  time.Time
}

// Implementar json.Unmarshaler
func (e *Evento) UnmarshalJSON(data []byte) error {
	type Alias Evento
	aux := &struct {
		Fecha string `json:"fecha"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse custom format
	t, err := time.Parse("2006-01-02", aux.Fecha)
	if err != nil {
		return err
	}

	e.Fecha = t
	return nil
}

func main() {
	jsonData := []byte(`{"Titulo":"Conferencia Go","fecha":"2024-06-15"}`)
	
	var evento Evento
	err := json.Unmarshal(jsonData, &evento)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("%+v\n", evento)
	// {Titulo:Conferencia Go Fecha:2024-06-15 00:00:00 +0000 UTC}
}
```

### 29.5.3 Custom Types con Validación

```go
package main

import (
	"encoding/json"
	"fmt"
	"errors"
)

type Email string

// Validar al unmarshaling
func (e *Email) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Validación simple
	if len(s) < 5 || !contains(s, "@") {
		return errors.New("email inválido")
	}

	*e = Email(s)
	return nil
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type Usuario struct {
	Nombre string `json:"nombre"`
	Email  Email  `json:"email"`
}

func main() {
	// Email válido
	jsonData1 := []byte(`{"nombre":"Ana","email":"ana@example.com"}`)
	var usuario1 Usuario
	err := json.Unmarshal(jsonData1, &usuario1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Usuario válido: %+v\n", usuario1)
	}

	// Email inválido
	jsonData2 := []byte(`{"nombre":"Bob","email":"bob"}`)
	var usuario2 Usuario
	err = json.Unmarshal(jsonData2, &usuario2)
	fmt.Println("Email inválido:", err)
	// Email inválido: email inválido
}
```

### 29.5.4 RawMessage para Datos Dinámicos

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Evento struct {
	Tipo   string          `json:"tipo"`
	Datos  json.RawMessage `json:"datos"` // JSON sin procesar
	Timestamp int64         `json:"timestamp"`
}

func main() {
	jsonData := []byte(`{
		"tipo": "usuario_creado",
		"datos": {"id": 123, "nombre": "Laura", "email": "laura@example.com"},
		"timestamp": 1705334400
	}`)

	var evento Evento
	err := json.Unmarshal(jsonData, &evento)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Tipo: %s\n", evento.Tipo)
	fmt.Printf("Datos raw: %s\n", string(evento.Datos))
	// Datos raw: {"id": 123, "nombre": "Laura", "email": "laura@example.com"}

	// Procesar datos específicamente según tipo
	if evento.Tipo == "usuario_creado" {
		var usuario map[string]interface{}
		json.Unmarshal(evento.Datos, &usuario)
		fmt.Printf("Usuario ID: %v\n", usuario["id"])
		// Usuario ID: 123
	}
}
```

---

## 29.6 Number Handling

### 29.6.1 Precision en Números

```go
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Finanzas struct {
	Precio    float64 `json:"precio"`
	Cantidad  int     `json:"cantidad"`
	Impuesto  float64 `json:"impuesto"`
}

func main() {
	// Problema: pérdida de precisión en float64
	jsonData := []byte(`{
		"precio": 19.99,
		"cantidad": 1000000000,
		"impuesto": 0.00000000001
	}`)

	var finanzas Finanzas
	json.Unmarshal(jsonData, &finanzas)
	
	fmt.Printf("Precio: %.20f\n", finanzas.Precio)
	fmt.Printf("Cantidad: %d\n", finanzas.Cantidad)
	fmt.Printf("Impuesto: %.20f\n", finanzas.Impuesto)
	// Precision limitada con float64
}
```

### 29.6.2 json.Number para Precisión Exacta

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Transaccion struct {
	ID        int    `json:"id"`
	Monto     json.Number `json:"monto"` // Mantiene precisión exacta
	Comisión  json.Number `json:"comision"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	jsonData := []byte(`{
		"id": 1,
		"monto": "12345.6789012345",
		"comision": "0.123",
		"timestamp": 1705334400
	}`)

	// Configurar Decoder para usar json.Number
	decoder := json.NewDecoder(nil)
	decoder.UseNumber()

	// Pero para Unmarshal, usamos UseNumber en un Decoder
	var transaccion Transaccion
	err := json.Unmarshal(jsonData, &transaccion)
	
	fmt.Printf("Monto: %s\n", transaccion.Monto)
	fmt.Printf("Comisión: %s\n", transaccion.Comisión)

	// Convertir json.Number
	monto, _ := transaccion.Monto.Float64()
	fmt.Printf("Monto como float64: %.10f\n", monto)

	// O mantener como string para procesamiento exacto
	montoStr := string(transaccion.Monto)
	fmt.Printf("Monto como string: %s\n", montoStr)
}
```

### 29.6.3 Decoder con UseNumber

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Medida struct {
	Nombre string      `json:"nombre"`
	Valor  json.Number `json:"valor"`
}

func main() {
	jsonStr := `{
		"nombre": "temperatura",
		"valor": 36.5
	}`

	// Opción 1: Sin UseNumber
	var medida1 Medida
	json.Unmarshal([]byte(jsonStr), &medida1)
	fmt.Println("Sin UseNumber:", medida1.Valor, "Tipo:", fmt.Sprintf("%T", medida1.Valor))

	// Opción 2: Con UseNumber
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.UseNumber()

	var medida2 Medida
	decoder.Decode(&medida2)
	fmt.Println("Con UseNumber:", medida2.Valor, "Tipo:", fmt.Sprintf("%T", medida2.Valor))

	// Conversión
	if i, err := medida2.Valor.Int64(); err == nil {
		fmt.Printf("Como int64: %d\n", i)
	}

	if f, err := medida2.Valor.Float64(); err == nil {
		fmt.Printf("Como float64: %f\n", f)
	}

	fmt.Printf("Como string: %s\n", string(medida2.Valor))
}
```

### 29.6.4 Errores en Conversiones Numéricas

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Valores struct {
	Numero int `json:"numero"`
}

func main() {
	// Error: string en campo numérico sin tag "string"
	jsonData := []byte(`{"numero": "123"}`)
	var valores Valores
	err := json.Unmarshal(jsonData, &valores)
	fmt.Println("String a int error:", err)
	// Error: json: cannot unmarshal string into Go struct field Valores.Numero of type int

	// Success: float válido a int
	jsonData2 := []byte(`{"numero": 456}`)
	err = json.Unmarshal(jsonData2, &valores)
	fmt.Println("Success:", err, "Valores:", valores)

	// Overflow
	jsonData3 := []byte(`{"numero": 99999999999999999999999999999}`)
	err = json.Unmarshal(jsonData3, &valores)
	fmt.Println("Overflow:", err)
}
```

---

## 29.7 Streaming JSON

### 29.7.1 Encoder/Decoder Streaming

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Evento struct {
	ID    int    `json:"id"`
	Tipo  string `json:"tipo"`
	Datos string `json:"datos"`
}

func main() {
	// Encoder: escribir múltiples objetos JSON
	var buffer bytes.Buffer

	encoder := json.NewEncoder(&buffer)

	eventos := []Evento{
		{ID: 1, Tipo: "login", Datos: "usuario1"},
		{ID: 2, Tipo: "logout", Datos: "usuario2"},
		{ID: 3, Tipo: "error", Datos: "conexión perdida"},
	}

	for _, evento := range eventos {
		encoder.Encode(evento) // Newline automático
	}

	fmt.Println("Encoded:")
	fmt.Println(buffer.String())
	/* Salida:
	{"id":1,"tipo":"login","datos":"usuario1"}
	{"id":2,"tipo":"logout","datos":"usuario2"}
	{"id":3,"tipo":"error","datos":"conexión perdida"}
	*/
}
```

### 29.7.2 Decoder Streaming Línea por Línea

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func main() {
	// Simular stream con múltiples objetos JSON por línea
	logData := `{"timestamp":"2024-01-15T10:30:00Z","level":"INFO","message":"Aplicación iniciada"}
{"timestamp":"2024-01-15T10:31:00Z","level":"DEBUG","message":"Conectando a BD"}
{"timestamp":"2024-01-15T10:32:00Z","level":"ERROR","message":"Error de conexión"}
{"timestamp":"2024-01-15T10:33:00Z","level":"INFO","message":"Reconectando..."}`

	decoder := json.NewDecoder(strings.NewReader(logData))

	var count int
	for decoder.More() {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Printf("[%s] %s: %s\n", entry.Level, entry.Timestamp, entry.Message)
		count++
	}

	fmt.Printf("\nTotal entradas: %d\n", count)
}
```

### 29.7.3 Procesamiento de JSON Grande

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Usuario struct {
	ID   int    `json:"id"`
	Nombre string `json:"nombre"`
}

func procesarUsuariosStream(jsonStream string) error {
	decoder := json.NewDecoder(strings.NewReader(jsonStream))

	// Expect JSON array start
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	fmt.Printf("Token inicial: %v\n", token)

	// Procesar objetos dentro del array
	for decoder.More() {
		var usuario Usuario
		if err := decoder.Decode(&usuario); err != nil {
			return err
		}

		// Procesar cada usuario (ej: guardar en BD, filtrar, etc)
		fmt.Printf("Procesando: %+v\n", usuario)
	}

	// Expect JSON array end
	token, err = decoder.Token()
	if err != nil {
		return err
	}
	fmt.Printf("Token final: %v\n", token)

	return nil
}

func main() {
	usuarios := `[
		{"id": 1, "nombre": "Alice"},
		{"id": 2, "nombre": "Bob"},
		{"id": 3, "nombre": "Charlie"}
	]`

	err := procesarUsuariosStream(usuarios)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

### 29.7.4 Token-based Parsing

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func analizarJSON(jsonStr string) {
	decoder := json.NewDecoder(strings.NewReader(jsonStr))

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch v := token.(type) {
		case json.Delim:
			fmt.Printf("Delimitador: %c\n", v)
		case string:
			fmt.Printf("String: %s\n", v)
		case float64:
			fmt.Printf("Número: %v\n", v)
		case bool:
			fmt.Printf("Boolean: %v\n", v)
		case nil:
			fmt.Println("Null")
		}
	}
}

func main() {
	jsonData := `{
		"nombre": "David",
		"edad": 30,
		"activo": true,
		"hobbies": ["lectura", "código"],
		"extra": null
	}`

	analizarJSON(jsonData)
}
```

---

## 29.8 Validation

### 29.8.1 Validación Manual Post-Unmarshal

```go
package main

import (
	"encoding/json"
	"fmt"
	"errors"
)

type Pedido struct {
	ID       int     `json:"id"`
	Cantidad int     `json:"cantidad"`
	Precio   float64 `json:"precio"`
	Email    string  `json:"email"`
}

// Validar después del unmarshal
func (p *Pedido) Validar() error {
	if p.ID <= 0 {
		return errors.New("ID debe ser positivo")
	}
	if p.Cantidad <= 0 {
		return errors.New("Cantidad debe ser positiva")
	}
	if p.Precio <= 0 {
		return errors.New("Precio debe ser positivo")
	}
	if len(p.Email) < 5 {
		return errors.New("Email inválido")
	}
	return nil
}

func main() {
	jsonData := []byte(`{
		"id": 1,
		"cantidad": -5,
		"precio": 99.99,
		"email": "usuario@example.com"
	}`)

	var pedido Pedido
	json.Unmarshal(jsonData, &pedido)

	if err := pedido.Validar(); err != nil {
		fmt.Println("Validación fallida:", err)
	}
}
```

### 29.8.2 Custom Unmarshaler con Validación

```go
package main

import (
	"encoding/json"
	"fmt"
	"errors"
	"strings"
)

type Repositorio struct {
	Nombre string
	Rama   string
}

func (r *Repositorio) UnmarshalJSON(data []byte) error {
	type Alias Repositorio
	aux := &struct {
		Nombre string `json:"nombre"`
		Rama   string `json:"rama"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Validaciones
	if strings.TrimSpace(aux.Nombre) == "" {
		return errors.New("nombre no puede estar vacío")
	}
	if aux.Rama != "main" && aux.Rama != "develop" && aux.Rama != "staging" {
		return errors.New("rama debe ser main, develop o staging")
	}

	r.Nombre = aux.Nombre
	r.Rama = aux.Rama
	return nil
}

func main() {
	// Válido
	json1 := []byte(`{"nombre":"myrepo","rama":"main"}`)
	var repo1 Repositorio
	err := json.Unmarshal(json1, &repo1)
	fmt.Println("Repo 1:", repo1, "Error:", err)

	// Inválido
	json2 := []byte(`{"nombre":"","rama":"main"}`)
	var repo2 Repositorio
	err = json.Unmarshal(json2, &repo2)
	fmt.Println("Error validación:", err)

	// Rama inválida
	json3 := []byte(`{"nombre":"myrepo","rama":"master"}`)
	var repo3 Repositorio
	err = json.Unmarshal(json3, &repo3)
	fmt.Println("Error rama:", err)
}
```

### 29.8.3 JSON Schema-like Validation

```go
package main

import (
	"encoding/json"
	"fmt"
	"errors"
)

type Configuracion struct {
	Puerto       int    `json:"puerto"`
	Host         string `json:"host"`
	Timeout      int    `json:"timeout"`
	MaxConnections int   `json:"max_connections"`
}

func (c *Configuracion) Validar() error {
	// Puerto
	if c.Puerto < 1024 || c.Puerto > 65535 {
		return errors.New("puerto debe estar entre 1024 y 65535")
	}

	// Host
	if c.Host == "" {
		return errors.New("host no puede estar vacío")
	}

	// Timeout
	if c.Timeout < 1000 || c.Timeout > 300000 {
		return errors.New("timeout debe estar entre 1000 y 300000 ms")
	}

	// MaxConnections
	if c.MaxConnections < 1 || c.MaxConnections > 10000 {
		return errors.New("max_connections debe estar entre 1 y 10000")
	}

	return nil
}

func main() {
	testCases := []string{
		`{"puerto":8080,"host":"localhost","timeout":30000,"max_connections":100}`,
		`{"puerto":80,"host":"example.com","timeout":30000,"max_connections":50}`, // puerto < 1024
		`{"puerto":9000,"host":"","timeout":30000,"max_connections":100}`,        // host vacío
	}

	for i, testCase := range testCases {
		var config Configuracion
		json.Unmarshal([]byte(testCase), &config)
		err := config.Validar()
		fmt.Printf("Caso %d: %v\n", i+1, err)
	}
}
```

---

## 29.9 XML y CSV

### 29.9.1 XML Marshaling

```go
package main

import (
	"encoding/xml"
	"fmt"
)

type Libro struct {
	XMLName xml.Name `xml:"libro"`
	ID      int      `xml:"id,attr"`
	Titulo  string   `xml:"titulo"`
	Autor   string   `xml:"autor"`
	Paginas int      `xml:"paginas"`
	Precio  float64  `xml:"precio"`
}

func main() {
	libro := Libro{
		ID:      1,
		Titulo:  "The Go Programming Language",
		Autor:   "Alan Donovan, Brian Kernighan",
		Paginas: 400,
		Precio:  54.99,
	}

	// Marshal
	xmlBytes, _ := xml.MarshalIndent(libro, "", "  ")
	fmt.Println(string(xmlBytes))
	/* Salida:
	<libro id="1">
	  <titulo>The Go Programming Language</titulo>
	  <autor>Alan Donovan, Brian Kernighan</autor>
	  <paginas>400</paginas>
	  <precio>54.99</precio>
	</libro>
	*/
}
```

### 29.9.2 XML Unmarshaling

```go
package main

import (
	"encoding/xml"
	"fmt"
)

type Persona struct {
	XMLName xml.Name `xml:"persona"`
	Nombre  string   `xml:"nombre"`
	Email   string   `xml:"email"`
	Telefono string  `xml:"telefono"`
}

func main() {
	xmlData := []byte(`
	<persona>
		<nombre>Elena Martín</nombre>
		<email>elena@example.com</email>
		<telefono>+34612345678</telefono>
	</persona>
	`)

	var persona Persona
	err := xml.Unmarshal(xmlData, &persona)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Persona: %+v\n", persona)
}
```

### 29.9.3 CSV con encoding/csv

```go
package main

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type Alumno struct {
	ID     string
	Nombre string
	Nota   string
}

func main() {
	// Escribir CSV
	csvData := [][]string{
		{"ID", "Nombre", "Nota"},
		{"1", "Alice", "9.5"},
		{"2", "Bob", "8.2"},
		{"3", "Charlie", "7.8"},
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	for _, record := range csvData {
		writer.Write(record)
	}
	writer.Flush()

	fmt.Println("CSV escrito:")
	fmt.Println(builder.String())

	// Leer CSV
	csvReader := csv.NewReader(strings.NewReader(builder.String()))
	records, _ := csvReader.ReadAll()

	fmt.Println("\nCSV leído:")
	for i, record := range records {
		if i == 0 {
			continue // Saltar header
		}
		alumno := Alumno{
			ID:     record[0],
			Nombre: record[1],
			Nota:   record[2],
		}
		fmt.Printf("%+v\n", alumno)
	}
}
```

### 29.9.4 Elección de Formatos

| Formato | Uso                     | Ventajas          | Desventajas |
|---------|-------------------------|-------------------|------------|
| JSON    | APIs REST, Config       | Legible, rápido   | Sem sem sem |
| XML     | SOAP, Legado            | Esquema, namespaces | Verboso    |
| CSV     | Datos tabulares         | Compacto, simple  | Sin tipos  |
| Protocol| Alto rendimiento        | Compacto, rápido  | Compilar   |

```go
package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

type Producto struct {
	XMLName xml.Name `xml:"producto"`
	ID      int      `xml:"id" json:"id"`
	Nombre  string   `xml:"nombre" json:"nombre"`
	Precio  float64  `xml:"precio" json:"precio"`
}

func main() {
	producto := Producto{ID: 1, Nombre: "Laptop", Precio: 999.99}

	// JSON
	jsonBytes, _ := json.Marshal(producto)
	fmt.Println("JSON:", string(jsonBytes))

	// XML
	xmlBytes, _ := xml.Marshal(producto)
	fmt.Println("XML:", string(xmlBytes))

	// CSV (manual)
	var csvBuilder strings.Builder
	csvWriter := csv.NewWriter(&csvBuilder)
	csvWriter.Write([]string{"ID", "Nombre", "Precio"})
	csvWriter.Write([]string{fmt.Sprint(producto.ID), producto.Nombre, fmt.Sprint(producto.Precio)})
	csvWriter.Flush()
	fmt.Println("CSV:", csvBuilder.String())
}
```

---

## 29.10 Base64 y Hex

### 29.10.1 Base64 Encoding

```go
package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	datos := "Hello, World! Esta es una prueba de encoding."

	// Standard encoding
	encoded := base64.StdEncoding.EncodeToString([]byte(datos))
	fmt.Println("Standard Base64:", encoded)
	// Standard Base64: SGVsbG8sIFdvcmxkISBFc3RhIGVzIHVuYSBwcnVlYmEgZGUgZW5jb2Rpbmcu

	// URL-safe encoding
	encodedURL := base64.URLEncoding.EncodeToString([]byte(datos))
	fmt.Println("URL-safe Base64:", encodedURL)

	// Decodificar
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	fmt.Println("Decodificado:", string(decoded))
}
```

### 29.10.2 Base64 en JSON

```go
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Documento struct {
	Nombre   string `json:"nombre"`
	Contenido string `json:"contenido"` // Base64
}

func main() {
	doc := Documento{
		Nombre:    "reporte.pdf",
		Contenido: base64.StdEncoding.EncodeToString([]byte("PDF binary data...")),
	}

	jsonBytes, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Println("JSON con Base64:")
	fmt.Println(string(jsonBytes))

	// Decodificar
	var doc2 Documento
	json.Unmarshal(jsonBytes, &doc2)
	contenido, _ := base64.StdEncoding.DecodeString(doc2.Contenido)
	fmt.Println("Contenido decodificado:", string(contenido))
}
```

### 29.10.3 Hex Encoding

```go
package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	datos := []byte("Secret message")

	// Hex encoding
	encoded := hex.EncodeToString(datos)
	fmt.Println("Hex:", encoded)
	// Hex: 5365637265742070657373616765

	// Hex decoding
	decoded, _ := hex.DecodeString(encoded)
	fmt.Println("Decodificado:", string(decoded))

	// Útil para checksums
	hash := hex.EncodeToString([]byte{0x1A, 0x2B, 0x3C, 0x4D})
	fmt.Println("Hash hex:", hash)
}
```

### 29.10.4 Encoding Chains

```go
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Token struct {
	Usuario string `json:"usuario"`
	Datos   string `json:"datos"` // Contenido base64
}

func main() {
	// Encriptar información
	info := "usuario=juan|id=123|rol=admin"

	// 1. JSON
	token := Token{
		Usuario: "juan",
		Datos:   base64.StdEncoding.EncodeToString([]byte(info)),
	}

	// 2. Marshal a JSON
	jsonBytes, _ := json.Marshal(token)
	fmt.Println("JSON:", string(jsonBytes))

	// 3. Base64 del JSON
	jsonBase64 := base64.StdEncoding.EncodeToString(jsonBytes)
	fmt.Println("JSON as Base64:", jsonBase64)

	// Desencriptar
	jsonFromBase64, _ := base64.StdEncoding.DecodeString(jsonBase64)
	var token2 Token
	json.Unmarshal(jsonFromBase64, &token2)

	datosDecodificados, _ := base64.StdEncoding.DecodeString(token2.Datos)
	fmt.Println("Datos decodificados:", string(datosDecodificados))
}
```

---

## 29.11 Buenas Prácticas y Patrones

### 29.11.1 Versionado de API

```go
package main

import (
	"encoding/json"
	"fmt"
)

// V1 API
type UsuarioV1 struct {
	ID   int    `json:"id"`
	Nombre string `json:"nombre"`
}

// V2 API - más campos
type UsuarioV2 struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// V3 API - nombres diferentes
type UsuarioV3 struct {
	ID        int    `json:"user_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type APIResponse struct {
	Version string          `json:"version"`
	Data    json.RawMessage `json:"data"`
}

func main() {
	// Soportar múltiples versiones
	responses := []string{
		`{"version":"1.0","data":{"id":1,"nombre":"Alice"}}`,
		`{"version":"2.0","data":{"id":1,"nombre":"Bob","email":"bob@ex.com","created_at":"2024-01-15"}}`,
		`{"version":"3.0","data":{"user_id":1,"full_name":"Charlie","email":"charlie@ex.com","created_at":"2024-01-15","updated_at":"2024-01-20"}}`,
	}

	for _, respStr := range responses {
		var resp APIResponse
		json.Unmarshal([]byte(respStr), &resp)

		switch resp.Version {
		case "1.0":
			var user UsuarioV1
			json.Unmarshal(resp.Data, &user)
			fmt.Printf("V1: %+v\n", user)
		case "2.0":
			var user UsuarioV2
			json.Unmarshal(resp.Data, &user)
			fmt.Printf("V2: %+v\n", user)
		case "3.0":
			var user UsuarioV3
			json.Unmarshal(resp.Data, &user)
			fmt.Printf("V3: %+v\n", user)
		}
	}
}
```

### 29.11.2 Compatibilidad Hacia Adelante

```go
package main

import (
	"encoding/json"
	"fmt"
)

// Nueva estructura con campos adicionales
type UsuarioModerno struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Telefono string `json:"telefono"`
	Extras   map[string]interface{} `json:"extras,omitempty"`
}

func main() {
	// Datos de versión antigua (sin teléfono ni extras)
	jsonAntiguo := []byte(`{"id":1,"nombre":"Diana","email":"diana@ex.com"}`)

	var usuario UsuarioModerno
	err := json.Unmarshal(jsonAntiguo, &usuario)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Usuario: %+v\n", usuario)
	// Usuario: {ID:1 Nombre:Diana Email:diana@ex.com Telefono: Extras:map[]}

	// Datos nuevos con campos adicionales
	jsonNuevo := []byte(`{
		"id":2,
		"nombre":"Elena",
		"email":"elena@ex.com",
		"telefono":"+34612345678",
		"extras":{"empresa":"TechCorp","departamento":"Engineering"}
	}`)

	err = json.Unmarshal(jsonNuevo, &usuario)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Usuario: %+v\n", usuario)
	// Usuario: {ID:2 Nombre:Elena Email:elena@ex.com Telefono:+34612345678 Extras:map[departamento:Engineering empresa:TechCorp]}
}
```

### 29.11.3 Error Handling y Logging

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type LogError struct {
	Timestamp string `json:"timestamp"`
	Error     string `json:"error"`
	Context   map[string]interface{} `json:"context"`
}

func procesarJSON(datos []byte, dest interface{}) error {
	if err := json.Unmarshal(datos, dest); err != nil {
		logErr := LogError{
			Timestamp: "2024-01-15T10:30:00Z",
			Error:     err.Error(),
			Context: map[string]interface{}{
				"input_size": len(datos),
				"dest_type": fmt.Sprintf("%T", dest),
			},
		}

		logJSON, _ := json.MarshalIndent(logErr, "", "  ")
		log.Printf("Error procesando JSON: %s", string(logJSON))
		return err
	}
	return nil
}

func main() {
	type Usuario struct {
		Nombre string `json:"nombre"`
		Edad   int    `json:"edad"`
	}

	// JSON inválido
	jsonInvalido := []byte(`{"nombre":"Frank","edad":"treinta"}`)
	var usuario Usuario

	if err := procesarJSON(jsonInvalido, &usuario); err != nil {
		fmt.Println("Fallido como se esperaba")
	}
}
```

### 29.11.4 Estructura de Datos Grandes

```go
package main

import (
	"encoding/json"
	"fmt"
)

// Usar índices para campos frecuentemente accesados
type Catalogo struct {
	Productos []Producto `json:"productos"`
	// Índice para búsqueda rápida
	productosIndex map[int]*Producto `json:"-"`
}

type Producto struct {
	ID    int    `json:"id"`
	Nombre string `json:"nombre"`
	Precio float64 `json:"precio"`
}

func (c *Catalogo) BuildIndex() {
	c.productosIndex = make(map[int]*Producto)
	for i := range c.Productos {
		c.productosIndex[c.Productos[i].ID] = &c.Productos[i]
	}
}

func (c *Catalogo) ObtenerProducto(id int) *Producto {
	return c.productosIndex[id]
}

func main() {
	jsonData := []byte(`{
		"productos":[
			{"id":1,"nombre":"Laptop","precio":1000},
			{"id":2,"nombre":"Mouse","precio":25},
			{"id":3,"nombre":"Teclado","precio":75}
		]
	}`)

	var catalogo Catalogo
	json.Unmarshal(jsonData, &catalogo)
	catalogo.BuildIndex()

	// Búsqueda rápida O(1)
	producto := catalogo.ObtenerProducto(2)
	fmt.Printf("Producto: %+v\n", producto)
	// Producto: &{ID:2 Nombre:Mouse Precio:25}
}
```

### 29.11.5 Seguridad y Validación

```go
package main

import (
	"encoding/json"
	"fmt"
	"errors"
	"strings"
)

type FormularioUsuario struct {
	Nombre   string `json:"nombre"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (f *FormularioUsuario) ValidarSeguridad() error {
	// Validar nombre
	if len(strings.TrimSpace(f.Nombre)) < 3 {
		return errors.New("nombre debe tener al menos 3 caracteres")
	}

	// Validar password
	if len(f.Password) < 8 {
		return errors.New("password debe tener al menos 8 caracteres")
	}

	// Validar bio (evitar inyección)
	if strings.Contains(f.Bio, "<script>") || strings.Contains(f.Bio, "javascript:") {
		return errors.New("bio contiene contenido malicioso")
	}

	// Limitar tamaño
	if len(f.Bio) > 1000 {
		return errors.New("bio demasiado larga")
	}

	return nil
}

func main() {
	testCases := []string{
		`{"nombre":"Jo","password":"123","bio":"Hola"}`,
		`{"nombre":"John","password":"securepass123","bio":"Soy programador"}`,
		`{"nombre":"Jane","password":"pass123","bio":"<script>alert('xss')</script>"}`,
	}

	for i, testCase := range testCases {
		var form FormularioUsuario
		json.Unmarshal([]byte(testCase), &form)

		err := form.ValidarSeguridad()
		if err != nil {
			fmt.Printf("Caso %d FALLIDO: %v\n", i+1, err)
		} else {
			fmt.Printf("Caso %d OK: %+v\n", i+1, form)
		}
	}
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Struct to JSON - Marshal/Unmarshal Básico

**Objetivo:** Practicar serialización y deserialización básica.

```go
package main

import (
	"encoding/json"
	"fmt"
)

// TODO: Define una estructura Libro con campos:
// - Titulo (string)
// - Autor (string)
// - AñoPublicacion (int)
// - Disponible (bool)
// - Precio (float64)

func main() {
	// TODO 1: Crear una instancia de Libro
	// libro := Libro{...}

	// TODO 2: Serializar a JSON (usa Marshal)
	// jsonBytes, err := json.Marshal(libro)

	// TODO 3: Imprimir el JSON
	// fmt.Println(string(jsonBytes))

	// TODO 4: Crear JSON manualmente
	// jsonData := []byte(`{...}`)

	// TODO 5: Deserializar a estructura (usa Unmarshal con puntero)
	// var libroDesde Libro
	// json.Unmarshal(jsonData, &libroDesde)

	// TODO 6: Imprimir estructura deserializada
	// fmt.Printf("%+v\n", libroDesde)
}
```

**Solución esperada:**
- Serializar y deserializar bidireccional funciona
- Los datos se preservan correctamente
- No hay errores de conversión de tipos

---

### Ejercicio 2: Custom Tags - Usando Struct Tags

**Objetivo:** Dominar struct tags y omitempty.

```go
package main

import (
	"encoding/json"
	"fmt"
)

// TODO: Define una estructura Empleado con tags JSON:
// - ID: campo id
// - Nombre: campo nombre
// - Departamento: campo departamento, omitempty
// - Salario: campo salario,string (convertir a string)
// - Gerente: campo gerente con puntero, omitempty
// - ApiKey: no debe aparecer en JSON (usar -)

func main() {
	// TODO 1: Crear empleado sin gerente
	// empleado1 := Empleado{ID: 1, Nombre: "Ana", Departamento: "", Salario: 50000}

	// TODO 2: Serializar a JSON formateado (usa MarshalIndent)
	// Departamento y Gerente no deben aparecer

	// TODO 3: Crear empleado con gerente
	// empleado2 con gerente asignado

	// TODO 4: Serializar el segundo
	// Los campos vacíos aún se omiten con omitempty

	// TODO 5: Deserializar JSON con campos faltantes
	// Debe funcionar sin problemas
}
```

---

### Ejercicio 3: Streaming - Procesar JSON Línea por Línea

**Objetivo:** Usar Decoder para procesar múltiples objetos JSON.

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaccion struct {
	ID        int     `json:"id"`
	Tipo      string  `json:"tipo"` // "credito" o "debito"
	Monto     float64 `json:"monto"`
	Concepto  string  `json:"concepto"`
}

func main() {
	// Simular un stream de transacciones (formato JSONL)
	jsonStream := `{"id":1,"tipo":"credito","monto":1000,"concepto":"Salario"}
{"id":2,"tipo":"debito","monto":500,"concepto":"Compra online"}
{"id":3,"tipo":"credito","monto":250,"concepto":"Reembolso"}
{"id":4,"tipo":"debito","monto":150,"concepto":"Suscripción"}`

	// TODO 1: Crear un Decoder desde el stream
	decoder := json.NewDecoder(strings.NewReader(jsonStream))

	// TODO 2: Iterar con decoder.More() y decoder.Decode()
	for decoder.More() {
		var trans Transaccion
		// TODO: Decode la transacción
		
		// TODO: Procesar (imprimir monto, calcular total por tipo, etc)
		fmt.Printf("%+v\n", trans)
	}

	// TODO 3: Calcular total de créditos y débitos
}
```

---

### Ejercicio 4: Validator - Implementar Validación Personalizada

**Objetivo:** Crear custom Unmarshaler con validación.

```go
package main

import (
	"encoding/json"
	"fmt"
	"errors"
)

type Rating int

// TODO: Implementar UnmarshalJSON para Rating
// Validar que el valor esté entre 1 y 5

type Resena struct {
	ID        int    `json:"id"`
	Autor     string `json:"autor"`
	Contenido string `json:"contenido"`
	Calificacion Rating `json:"calificacion"`
}

func main() {
	testCases := []string{
		`{"id":1,"autor":"Usuario1","contenido":"Excelente","calificacion":5}`,
		`{"id":2,"autor":"Usuario2","contenido":"Malo","calificacion":6}`, // Inválido
		`{"id":3,"autor":"Usuario3","contenido":"Normal","calificacion":3}`,
	}

	for i, jsonStr := range testCases {
		var resena Resena
		err := json.Unmarshal([]byte(jsonStr), &resena)
		if err != nil {
			fmt.Printf("Caso %d ERROR: %v\n", i+1, err)
		} else {
			fmt.Printf("Caso %d OK: %+v\n", i+1, resena)
		}
	}
}
```

---

### Ejercicio 5: Multi-format - Soportar JSON y XML

**Objetivo:** Serializar/deserializar en múltiples formatos.

```go
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Articulo struct {
	XMLName xml.Name `xml:"articulo"`
	ID      int      `xml:"id" json:"id"`
	Titulo  string   `xml:"titulo" json:"titulo"`
	Autor   string   `xml:"autor" json:"autor"`
	Fecha   string   `xml:"fecha" json:"fecha"`
	Tags    []string `xml:"tags>tag" json:"tags"`
}

func main() {
	// TODO 1: Crear instancia de Articulo
	// articulo := Articulo{...}

	// TODO 2: Serializar a JSON
	// jsonBytes, _ := json.MarshalIndent(articulo, "", "  ")
	// fmt.Println("JSON:")
	// fmt.Println(string(jsonBytes))

	// TODO 3: Serializar a XML
	// xmlBytes, _ := xml.MarshalIndent(articulo, "", "  ")
	// fmt.Println("\nXML:")
	// fmt.Println(string(xmlBytes))

	// TODO 4: Deserializar desde JSON
	// jsonData := []byte(...)
	// var artDesdeJSON Articulo
	// json.Unmarshal(...)

	// TODO 5: Deserializar desde XML
	// xmlData := []byte(...)
	// var artDesdeXML Articulo
	// xml.Unmarshal(...)

	// TODO 6: Comparar que ambos contienen los mismos datos
}
```

---

## Resumen de Conceptos Clave

| Concepto | Descripción | Ejemplo |
|----------|-------------|---------|
| **Marshal** | Serializar Go → JSON | `json.Marshal(obj)` |
| **Unmarshal** | Deserializar JSON → Go | `json.Unmarshal(data, &obj)` |
| **Tags** | Metadatos en struct | `` json:"campo,omitempty" `` |
| **Marshaler** | Interface customizado | `MarshalJSON()` |
| **Unmarshaler** | Interface desserialización | `UnmarshalJSON(data)` |
| **Streaming** | Procesar JSON grande | `Encoder/Decoder` |
| **json.Number** | Precisión exacta | Usar para moneda |
| **RawMessage** | JSON sin procesar | `json.RawMessage` |

---

## Comparativa: Go vs Otros Lenguajes

### Go (Recomendado para API REST)
```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
json.Marshal(user)
```

### Python (Dinámico)
```python
import json
user = {"name": "John", "age": 30}
json.dumps(user)
```

### Java (Verboso)
```java
ObjectMapper mapper = new ObjectMapper();
mapper.writeValueAsString(user);
```

---

## Casos de Uso Reales

1. **APIs REST**: Comunicación cliente-servidor con JSON
2. **Archivos de Configuración**: YAML/JSON para settings
3. **Bases de Datos**: Serializar objetos antes de guardar
4. **Streaming de Datos**: Procesar logs, eventos, transacciones
5. **GraphQL**: Transferencia de datos complejos

---

## Antipatrones a Evitar

❌ **No verificar errores en Unmarshal**
```go
json.Unmarshal(data, &obj) // Ignorar error es peligroso
```

✅ **Verificar siempre**
```go
if err := json.Unmarshal(data, &obj); err != nil {
    log.Fatal("Error:", err)
}
```

---

❌ **Usar float64 para dinero**
```go
type Precio struct {
    Monto float64 `json:"monto"` // Pérdida de precisión
}
```

✅ **Usar string o big.Decimal**
```go
type Precio struct {
    Monto string `json:"monto"` // Precisión exacta
}
```

---

❌ **Olvidar puntero en Unmarshal**
```go
json.Unmarshal(data, obj) // ERROR: debe ser &obj
```

✅ **Usar puntero**
```go
json.Unmarshal(data, &obj) // Correcto
```

---

## Conclusión

JSON es el formato estándar para APIs modernas en Go. Dominar `Marshal`/`Unmarshal`, tags y custom types permite crear sistemas robustos de serialización. Para proyectos grandes, considera usar librerías como `json-iterator` o `easyjson` para mejor performance.

**Próximos pasos:** Integrar JSON en un servidor HTTP REST con el paquete `net/http`.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/29-json-y-encoding/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/29-json-y-encoding):

```bash
cd examples/29-json-y-encoding
go run .
```
