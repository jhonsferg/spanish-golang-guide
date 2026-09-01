# Capítulo 33: Reflect - Introspección y metaprogramación

## Introducción

El package `reflect` es una de las características más poderosas pero también más complejas de Go. Permite a tu código examinar y manipular tipos y valores en tiempo de ejecución, algo generalmente imposible en lenguajes compilados estáticos. Sin embargo, con gran poder viene una gran responsabilidad: el uso inadecuado de `reflect` puede hacer que el código sea lento, difícil de entender y propenso a errores en tiempo de ejecución.

Este capítulo te enseña cómo usar reflection correctamente: cuándo es apropiado, cuándo evitarlo, cómo implementarlo de forma segura, y cómo optimizar su rendimiento.

---

## 33.1 ¿Qué es Reflect?

### 33.1.1 Conceptos Fundamentales

`Reflect` es la capacidad de un programa para examinar su propia estructura en tiempo de ejecución. En Go, el package `reflect` proporciona:

1. **Introspección**: Examinar tipos, valores y estructura de datos en runtime
2. **Metaprogramación**: Generar o modificar código dinámicamente
3. **Serialización Genérica**: Implementar JSON, XML, ORMs sin especificar cada tipo
4. **Dispatch Dinámico**: Llamar métodos o funciones sin conocer el tipo exacto

### 33.1.2 ¿Por Qué No Siempre Usar Reflect?

Go es un lenguaje de tipos estáticos por una razón: la seguridad en tiempo de compilación. Reflect sacrifica esa seguridad por flexibilidad:

| Aspecto | Reflection | Tipado Estático |
|--------|-----------|-----------------|
| **Seguridad** | En runtime (errores en ejecución) | En compile-time |
| **Performance** | 100-1000x más lento | Nativo |
| **Debugging** | Difícil de seguir el código | Stack traces claros |
| **IDEs** | Autocompletado inefectivo | Autocompletado perfecto |
| **Refactoring** | Propenso a breakage silencioso | Detectado por compilador |

### 33.1.3 Casos de Uso Apropiados para Reflect

✅ **Buenas razones para usar reflect**:

1. **Frameworks generales**: JSON marshalers, ORMs, testing frameworks
2. **Serialización**: Convertir objetos a/desde diferentes formatos
3. **Dependency injection**: Inyectar dependencias automáticamente
4. **APIs públicas**: Librerías que necesitan soportar tipos arbitrarios
5. **Validación genérica**: Validadores que funcionan con cualquier struct

❌ **No uses reflect para**:

1. Lógica de negocio crítica en performance
2. Alternativa a interfaces bien diseñadas
3. Cuando puedas usar type switches
4. En loops críticos de performance

### 33.1.4 Arquitectura de Reflect

```
┌─────────────────────────────────────────────────────┐
│              Tipos y Valores en Runtime             │
├─────────────────────────────────────────────────────┤
│                   reflect Package                   │
│  ┌──────────┬──────────┬────────┬──────────────┐   │
│  │ Type     │ Value    │ Kind   │ Struct Tags  │   │
│  ├──────────┼──────────┼────────┼──────────────┤   │
│  │ Method   │ Field    │ Slice  │   Interface  │   │
│  └──────────┴──────────┴────────┴──────────────┘   │
├─────────────────────────────────────────────────────┤
│        Comparación con Java/Python Reflection       │
└─────────────────────────────────────────────────────┘
```

### 33.1.5 Comparación: Go vs Otros Lenguajes

```go
// ╔═══════════════════════════════════════════════════════════╗
// ║  Go: Rápido, Explícito pero Verboso                      ║
// ╚═══════════════════════════════════════════════════════════╝
package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
	Age  int
}

func inspectType(v interface{}) {
	t := reflect.TypeOf(v)
	fmt.Printf("Tipo: %v, Tipo Kind: %v\n", t, t.Kind())
}

// ╔═══════════════════════════════════════════════════════════╗
// ║  Python: Dinámico y Simple                               ║
// ╚═══════════════════════════════════════════════════════════╝
# Python
user = User("Alice", 30)
print(type(user))  # <class '__main__.User'>
print(user.__dict__)  # {'name': 'Alice', 'age': 30}

// ╔═══════════════════════════════════════════════════════════╗
// ║  Java: Complejo pero Poderoso                            ║
// ╚═══════════════════════════════════════════════════════════╝
// Java
Class<?> clazz = User.class;
Method[] methods = clazz.getMethods();
for (Method m : methods) {
    System.out.println(m.getName());
}
```

---

## 33.2 TypeOf y ValueOf

### 33.2.1 La Dualidad: Type vs Value

Reflection en Go se basa en dos conceptos centrales:

1. **`reflect.Type`**: Información estática del tipo (Kind, Fields, Methods)
2. **`reflect.Value`**: El valor actual junto con su información de tipo

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// Obtener Type y Value
	user := "Alice"
	
	t := reflect.TypeOf(user)      // *reflect.rtype (interno)
	v := reflect.ValueOf(user)     // reflect.Value (contiene tipo + valor)
	
	fmt.Println("Type:", t)              // string
	fmt.Println("Value:", v)             // Alice
	fmt.Println("Value interface:", v.Interface()) // Alice (como interface{})
}
```

### 33.2.2 reflect.Type - Información del Tipo

`reflect.Type` es una interface que describe la estructura del tipo:

```go
type Type interface {
	Align() int                                    // Alineación en memoria
	FieldAlign() int                               // Alineación de campos
	Method(int) Method                             // i-ésimo método
	MethodByName(string) (Method, bool)            // Método por nombre
	NumMethod() int                                // Número de métodos
	Name() string                                  // Nombre del tipo
	PkgPath() string                               // Path del package
	Size() uintptr                                 // Tamaño en bytes
	String() string                                // Representación string
	Kind() Kind                                    // Clasificación (struct, int, slice, etc.)
	Implements(u Type) bool                        // ¿Implementa interface?
	ConvertibleTo(u Type) bool                     // ¿Convertible a?
	Comparable() bool                              // ¿Es comparable?
	Bits() int                                     // Bits para tipos int/float
	ChanDir() ChanDir                              // Dirección del channel
	IsVariadic() bool                              // ¿Es función variadic?
	Elem() Type                                    // Tipo del elemento (ptr, slice, array)
	Field(i int) StructField                       // i-ésimo campo struct
	FieldByIndex(index []int) StructField          // Campo por índice anidado
	FieldByName(name string) (StructField, bool)   // Campo por nombre
	FieldByNameFunc(match func(string) bool) (StructField, bool)
	In(i int) Type                                 // Tipo i-ésimo parámetro
	NumIn() int                                    // Número de parámetros
	Out(i int) Type                                // Tipo i-ésimo retorno
	NumOut() int                                   // Número de retornos
	NumField() int                                 // Número de campos
}
```

### 33.2.3 Obtener Información de Tipos

```go
package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
	City string `json:"city"`
}

func (p *Person) Greet() {
	fmt.Printf("Hola, soy %s\n", p.Name)
}

func main() {
	p := Person{Name: "Alice", Age: 30, City: "Madrid"}
	t := reflect.TypeOf(p)
	
	// Información básica del tipo
	fmt.Printf("Nombre: %s\n", t.Name())                    // Person
	fmt.Printf("Kind: %v\n", t.Kind())                       // struct
	fmt.Printf("Size: %d bytes\n", t.Size())                 // Tamaño en memoria
	fmt.Printf("NumField: %d\n", t.NumField())               // 3
	fmt.Printf("NumMethod: %d\n", t.NumMethod())             // 0 (puntero tiene métodos)
	
	// Información de puntero
	pt := reflect.TypeOf(&p)
	fmt.Printf("Puntero Kind: %v\n", pt.Kind())              // ptr
	fmt.Printf("Puntero NumMethod: %d\n", pt.NumMethod())    // 1 (Greet)
	fmt.Printf("Elemento: %s\n", pt.Elem().Name())           // Person
}
```

### 33.2.4 reflect.Value - Obtener y Establecer Valores

```go
package main

import (
	"fmt"
	"reflect"
)

type Temperature struct {
	Celsius    float64
	Fahrenheit float64
}

func main() {
	temp := Temperature{Celsius: 25.0}
	v := reflect.ValueOf(&temp)
	
	// Acceder a valores
	fmt.Println("Value:", v)                       // &{25 0}
	fmt.Println("Elem:", v.Elem())                 // {25 0}
	fmt.Println("Interface:", v.Interface())       // &{25 0}
	
	// Convertir a tipo específico
	tempPtr := v.Interface().(*Temperature)
	fmt.Printf("Convertido: %v\n", tempPtr)
	
	// Modificar valores (REQUIERE PUNTERO)
	elemV := v.Elem()
	celsiusField := elemV.FieldByName("Celsius")
	
	if celsiusField.IsValid() && celsiusField.CanSet() {
		celsiusField.SetFloat(30.0)
		fmt.Printf("Actualizado: %v\n", temp)
	}
}
```

### 33.2.5 Kind vs Type

`Kind` es una clasificación simple del tipo, mientras que `Type` es más específico:

```go
package main

import (
	"fmt"
	"reflect"
)

type MyString string

func main() {
	// Type vs Kind
	s := MyString("hello")
	t := reflect.TypeOf(s)
	
	fmt.Printf("Type: %v\n", t)          // main.MyString
	fmt.Printf("Name: %s\n", t.Name())   // MyString
	fmt.Printf("Kind: %v\n", t.Kind())   // string (Kind siempre es simple)
	
	// Kind enumera tipos fundamentales
	kinds := []interface{}{
		42,           // int
		3.14,         // float64
		"hello",      // string
		[]int{1, 2},  // slice
		map[string]int{"a": 1}, // map
		struct{}{},   // struct
	}
	
	for _, k := range kinds {
		fmt.Printf("%T -> %v\n", k, reflect.TypeOf(k).Kind())
	}
	// output:
	// int -> int
	// float64 -> float64
	// string -> string
	// []int -> slice
	// map[string]int -> map
	// struct {} -> struct
}
```

---

## 33.3 Tipo de Datos - Introspección Profunda

### 33.3.1 Kind: Clasificación de Tipos

```go
const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)
```

### 33.3.2 Introspección de Structs

```go
package main

import (
	"fmt"
	"reflect"
)

type Product struct {
	ID       int       `json:"id" db:"id,pk"`
	Name     string    `json:"name" db:"name,notnull"`
	Price    float64   `json:"price" db:"price"`
	InStock  bool      `json:"-" db:"in_stock"`
	Tags     []string  `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"-"`
}

func inspectStruct(v interface{}) {
	t := reflect.TypeOf(v)
	
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	if t.Kind() != reflect.Struct {
		fmt.Println("No es un struct")
		return
	}
	
	fmt.Printf("Struct: %s\n", t.Name())
	fmt.Printf("Número de campos: %d\n\n", t.NumField())
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		
		fmt.Printf("Campo %d:\n", i)
		fmt.Printf("  Nombre: %s\n", field.Name)
		fmt.Printf("  Tipo: %v\n", field.Type)
		fmt.Printf("  Kind: %v\n", field.Type.Kind())
		fmt.Printf("  Tag (json): %s\n", field.Tag.Get("json"))
		fmt.Printf("  Tag (db): %s\n", field.Tag.Get("db"))
		fmt.Printf("  Exportado: %v\n", field.IsExported())
		fmt.Printf("  Embedded: %v\n", field.Anonymous)
		fmt.Println()
	}
}

func main() {
	p := Product{ID: 1, Name: "Laptop"}
	inspectStruct(p)
}
```

### 33.3.3 Acceder a Campos por Nombre

```go
package main

import (
	"fmt"
	"log"
	"reflect"
)

type Config struct {
	Host     string
	Port     int
	Debug    bool
	Timeout  int
}

func getFieldValue(obj interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(obj)
	
	// Si es puntero, desreferencia
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	
	// Buscar el campo
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		log.Fatalf("Campo no encontrado: %s", fieldName)
	}
	
	return field.Interface()
}

func main() {
	cfg := Config{
		Host:    "localhost",
		Port:    8080,
		Debug:   true,
		Timeout: 30,
	}
	
	fmt.Printf("Host: %v\n", getFieldValue(cfg, "Host"))
	fmt.Printf("Port: %v\n", getFieldValue(cfg, "Port"))
	fmt.Printf("Debug: %v\n", getFieldValue(cfg, "Debug"))
	
	// Accediendo con puntero
	fmt.Printf("Timeout: %v\n", getFieldValue(&cfg, "Timeout"))
}
```

### 33.3.4 Campos Anidados

```go
package main

import (
	"fmt"
	"reflect"
)

type Address struct {
	Street string
	City   string
}

type Company struct {
	Name    string
	Address Address  // Campo anidado
}

func printAllFields(v interface{}, prefix string) {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}
	
	if t.Kind() != reflect.Struct {
		return
	}
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)
		
		fullPath := prefix + field.Name
		
		if field.Type.Kind() == reflect.Struct {
			// Recursivamente procesar struct anidado
			printAllFields(value.Interface(), fullPath+".")
		} else {
			fmt.Printf("%s = %v\n", fullPath, value.Interface())
		}
	}
}

func main() {
	company := Company{
		Name: "TechCorp",
		Address: Address{
			Street: "Calle Principal 123",
			City:   "Madrid",
		},
	}
	
	printAllFields(company, "")
}
```

---

## 33.4 Valores - Obtener y Establecer

### 33.4.1 Métodos Get para reflect.Value

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// Diferentes tipos de valores
	values := []interface{}{
		42,
		3.14,
		"hello",
		true,
		[]int{1, 2, 3},
	}
	
	for _, v := range values {
		val := reflect.ValueOf(v)
		
		fmt.Printf("Valor: %v, Kind: %v\n", val.Interface(), val.Kind())
		
		// Métodos específicos por kind
		switch val.Kind() {
		case reflect.Int:
			fmt.Printf("  Int: %d\n", val.Int())
		case reflect.Float64:
			fmt.Printf("  Float: %f\n", val.Float())
		case reflect.String:
			fmt.Printf("  String: %s\n", val.String())
		case reflect.Bool:
			fmt.Printf("  Bool: %v\n", val.Bool())
		case reflect.Slice:
			fmt.Printf("  Len: %d, Cap: %d\n", val.Len(), val.Cap())
		}
	}
}
```

### 33.4.2 Métodos Set - Modificar Valores

Para modificar un valor, **SIEMPRE necesitas un puntero**:

```go
package main

import (
	"fmt"
	"reflect"
)

type Settings struct {
	MaxConnections int
	Timeout        float64
	Debug          bool
	Name           string
}

func setFieldValue(obj interface{}, fieldName string, newValue interface{}) error {
	v := reflect.ValueOf(obj)
	
	// IMPORTANTE: obj debe ser un puntero
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("obj debe ser un puntero")
	}
	
	// Desreferenciar el puntero
	v = v.Elem()
	
	// Obtener el campo
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("campo no encontrado: %s", fieldName)
	}
	
	// Verificar si se puede establecer
	if !field.CanSet() {
		return fmt.Errorf("no se puede establecer campo: %s", fieldName)
	}
	
	// Convertir el nuevo valor al tipo del campo
	newVal := reflect.ValueOf(newValue)
	
	if !newVal.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("tipo incompatible: %v vs %v", newVal.Type(), field.Type())
	}
	
	field.Set(newVal)
	return nil
}

func main() {
	settings := Settings{
		MaxConnections: 100,
		Timeout:        30.0,
		Debug:          false,
		Name:           "prod",
	}
	
	fmt.Printf("Antes: %+v\n", settings)
	
	setFieldValue(&settings, "MaxConnections", 200)
	setFieldValue(&settings, "Timeout", 60.0)
	setFieldValue(&settings, "Debug", true)
	setFieldValue(&settings, "Name", "staging")
	
	fmt.Printf("Después: %+v\n", settings)
}
```

### 33.4.3 Métodos Set Específicos

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// SetInt, SetFloat, SetString, SetBool para tipos específicos
	var container struct {
		Count   int
		Rating  float64
		Message string
		Active  bool
	}
	
	v := reflect.ValueOf(&container).Elem()
	
	// Usar Set Methods específicos
	v.FieldByName("Count").SetInt(42)
	v.FieldByName("Rating").SetFloat(4.5)
	v.FieldByName("Message").SetString("Hola!")
	v.FieldByName("Active").SetBool(true)
	
	fmt.Printf("%+v\n", container)
	// Output: {Count:42 Rating:4.5 Message:Hola! Active:true}
}
```

### 33.4.4 Elem y Indirection

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// Elem: desreferencia un puntero o consigue tipo del contenedor
	
	// Caso 1: Puntero a int
	x := 42
	ptrX := &x
	v := reflect.ValueOf(ptrX)
	
	fmt.Printf("v.Kind: %v\n", v.Kind())           // ptr
	fmt.Printf("v.Elem().Kind: %v\n", v.Elem().Kind()) // int
	fmt.Printf("v.Elem().Int(): %d\n", v.Elem().Int())  // 42
	
	// Caso 2: Slice
	nums := []int{1, 2, 3}
	v = reflect.ValueOf(nums)
	
	fmt.Printf("\nSlice elem type: %v\n", v.Type().Elem()) // int
	
	// Caso 3: Array
	arr := [3]int{10, 20, 30}
	v = reflect.ValueOf(arr)
	
	fmt.Printf("Array elem type: %v\n", v.Type().Elem()) // int
	
	// Caso 4: Map
	m := map[string]int{"a": 1}
	v = reflect.ValueOf(m)
	
	fmt.Printf("Map key type: %v\n", v.Type().Key())     // string
	fmt.Printf("Map value type: %v\n", v.Type().Elem())  // int
}
```

---

## 33.5 Structs Reflection - Trabajar con Estructuras

### 33.5.1 Iterar sobre Campos de un Struct

```go
package main

import (
	"fmt"
	"reflect"
)

type Article struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Author      string `json:"author" validate:"required"`
	Published   bool   `json:"published"`
	Views       int    `json:"views"`
}

func printStructInfo(obj interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	
	// Desreferenciar si es puntero
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	
	if t.Kind() != reflect.Struct {
		fmt.Println("No es un struct")
		return
	}
	
	fmt.Printf("=== %s ===\n", t.Name())
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		
		fmt.Printf("\n%d. %s\n", i+1, field.Name)
		fmt.Printf("   Type: %v\n", field.Type)
		fmt.Printf("   Value: %v\n", value.Interface())
		fmt.Printf("   Tag (json): %s\n", field.Tag.Get("json"))
		fmt.Printf("   Tag (validate): %s\n", field.Tag.Get("validate"))
	}
}

func main() {
	article := Article{
		Title:       "Reflection en Go",
		Description: "Una guía profunda",
		Author:      "Alice",
		Published:   true,
		Views:       1234,
	}
	
	printStructInfo(article)
}
```

### 33.5.2 Tag Parsing

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	ID       int    `json:"id" db:"id,pk" validate:"required"`
	Email    string `json:"email" db:"email,notnull,unique" validate:"required,email"`
	Name     string `json:"name" db:"name" validate:"required,min=3"`
	Password string `json:"-" db:"password,notnull"`
}

// Parsear tags customizados
func parseTag(tag string) map[string]string {
	result := make(map[string]string)
	
	if tag == "" {
		return result
	}
	
	parts := strings.Split(tag, ",")
	if len(parts) > 0 {
		result["name"] = parts[0]
		if len(parts) > 1 {
			result["options"] = strings.Join(parts[1:], ",")
		}
	}
	
	return result
}

func printTagInfo(obj interface{}) {
	t := reflect.TypeOf(obj)
	
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		
		fmt.Printf("\nCampo: %s\n", field.Name)
		
		// Tags principales
		for _, tagName := range []string{"json", "db", "validate"} {
			tag := field.Tag.Get(tagName)
			if tag != "" {
				parsed := parseTag(tag)
				fmt.Printf("  %s:\n", tagName)
				fmt.Printf("    name: %s\n", parsed["name"])
				fmt.Printf("    options: %s\n", parsed["options"])
			}
		}
	}
}

func main() {
	user := User{ID: 1, Email: "alice@example.com"}
	printTagInfo(user)
}
```

### 33.5.3 Copiar Datos entre Structs

```go
package main

import (
	"fmt"
	"reflect"
)

type PersonDTO struct {
	Name string
	Age  int
}

type Person struct {
	Name string
	Age  int
	ID   int // No existe en DTO
}

// Copiar campos comunes entre structs
func copyCommonFields(from interface{}, to interface{}) error {
	fromV := reflect.ValueOf(from)
	toV := reflect.ValueOf(to)
	
	if toV.Kind() != reflect.Ptr {
		return fmt.Errorf("'to' debe ser un puntero")
	}
	
	toV = toV.Elem()
	
	if fromV.Kind() == reflect.Ptr {
		fromV = fromV.Elem()
	}
	
	fromT := fromV.Type()
	toT := toV.Type()
	
	for i := 0; i < fromT.NumField(); i++ {
		fromField := fromT.Field(i)
		toField, exists := toT.FieldByName(fromField.Name)
		
		if !exists || fromField.Type != toField.Type {
			continue
		}
		
		toValue := toV.FieldByName(fromField.Name)
		if toValue.CanSet() {
			fromValue := fromV.Field(i)
			toValue.Set(fromValue)
		}
	}
	
	return nil
}

func main() {
	dto := PersonDTO{Name: "Alice", Age: 30}
	person := Person{ID: 100}
	
	copyCommonFields(dto, &person)
	
	fmt.Printf("Person: %+v\n", person)
	// Output: Person: {Name:Alice Age:30 ID:100}
}
```

### 33.5.4 Aplicar Cambios Genéricos a Structs

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Bio       string
}

// Aplicar función a todos los campos string
func applyToStringFields(obj interface{}, fn func(string) string) error {
	v := reflect.ValueOf(obj)
	
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("obj debe ser puntero")
	}
	
	v = v.Elem()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		
		if field.Kind() == reflect.String && field.CanSet() {
			original := field.String()
			modified := fn(original)
			field.SetString(modified)
		}
	}
	
	return nil
}

func main() {
	user := User{
		FirstName: "john",
		LastName:  "doe",
		Email:     "john@example.com",
		Bio:       "software engineer",
	}
	
	fmt.Printf("Antes: %+v\n", user)
	
	// Convertir a mayúsculas
	applyToStringFields(&user, strings.ToUpper)
	
	fmt.Printf("Después: %+v\n", user)
}
```

---

## 33.6 Methods Reflection - Llamar Métodos Dinámicamente

### 33.6.1 Introspección de Métodos

```go
package main

import (
	"fmt"
	"reflect"
)

type Calculator struct {
	LastResult float64
}

func (c *Calculator) Add(a, b float64) float64 {
	c.LastResult = a + b
	return c.LastResult
}

func (c *Calculator) Multiply(a, b float64) float64 {
	c.LastResult = a * b
	return c.LastResult
}

func (c *Calculator) GetResult() float64 {
	return c.LastResult
}

func inspectMethods(obj interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	
	fmt.Printf("Tipo: %v\n", t)
	fmt.Printf("Número de métodos: %d\n\n", t.NumMethod())
	
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methodValue := v.Method(i)
		
		fmt.Printf("Método %d: %s\n", i+1, method.Name)
		fmt.Printf("  Type: %v\n", method.Type)
		fmt.Printf("  Inputs: %d\n", method.Type.NumIn())  // Incluye receiver
		fmt.Printf("  Outputs: %d\n", method.Type.NumOut())
		
		for j := 0; j < method.Type.NumIn(); j++ {
			fmt.Printf("    In(%d): %v\n", j, method.Type.In(j))
		}
		
		for j := 0; j < method.Type.NumOut(); j++ {
			fmt.Printf("    Out(%d): %v\n", j, method.Type.Out(j))
		}
		fmt.Println()
	}
}

func main() {
	calc := &Calculator{}
	inspectMethods(calc)
}
```

### 33.6.2 Llamar Métodos Dinámicamente

```go
package main

import (
	"fmt"
	"reflect"
)

type Processor struct {
	name string
}

func (p *Processor) Process(data string) string {
	return fmt.Sprintf("[%s] %s", p.name, data)
}

func (p *Processor) ProcessInt(x int) int {
	return x * 2
}

func callMethod(obj interface{}, methodName string, args ...interface{}) (interface{}, error) {
	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)
	
	if !method.IsValid() {
		return nil, fmt.Errorf("método no encontrado: %s", methodName)
	}
	
	// Convertir argumentos a reflect.Value
	var reflectArgs []reflect.Value
	for _, arg := range args {
		reflectArgs = append(reflectArgs, reflect.ValueOf(arg))
	}
	
	// Validar número de argumentos
	if len(reflectArgs) != method.Type().NumIn() {
		return nil, fmt.Errorf("número de argumentos incorrecto: esperado %d, obtenido %d",
			method.Type().NumIn(), len(reflectArgs))
	}
	
	// Llamar el método
	results := method.Call(reflectArgs)
	
	// Retornar los resultados
	if len(results) == 0 {
		return nil, nil
	}
	
	if len(results) == 1 {
		return results[0].Interface(), nil
	}
	
	// Múltiples retornos
	var retvals []interface{}
	for _, r := range results {
		retvals = append(retvals, r.Interface())
	}
	return retvals, nil
}

func main() {
	proc := &Processor{name: "Processor1"}
	
	// Llamar Process
	result, err := callMethod(proc, "Process", "hello")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
	
	// Llamar ProcessInt
	result, err = callMethod(proc, "ProcessInt", 21)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
```

### 33.6.3 Llamar Métodos con Error Handling

```go
package main

import (
	"fmt"
	"reflect"
)

type DataStore struct{}

func (d *DataStore) Get(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("clave vacía")
	}
	return "valor", nil
}

func callMethodWithErrors(obj interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)
	
	if !method.IsValid() {
		return nil, fmt.Errorf("método no encontrado: %s", methodName)
	}
	
	var reflectArgs []reflect.Value
	for _, arg := range args {
		reflectArgs = append(reflectArgs, reflect.ValueOf(arg))
	}
	
	results := method.Call(reflectArgs)
	
	// Verificar si el último resultado es error
	var values []interface{}
	var lastErr error
	
	for i, r := range results {
		// Verificar si es un error
		if err, ok := r.Interface().(error); ok && i == len(results)-1 {
			lastErr = err
			continue
		}
		values = append(values, r.Interface())
	}
	
	return values, lastErr
}

func main() {
	ds := &DataStore{}
	
	// Llamada exitosa
	results, err := callMethodWithErrors(ds, "Get", "mykey")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Results:", results)
	}
	
	// Llamada con error
	results, err = callMethodWithErrors(ds, "Get", "")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Results:", results)
	}
}
```

---

## 33.7 Function Reflection - Introspección de Funciones

### 33.7.1 Obtener Información de Funciones

```go
package main

import (
	"fmt"
	"reflect"
)

func add(a, b int) int {
	return a + b
}

func process(name string, value int) (string, error) {
	if value < 0 {
		return "", fmt.Errorf("valor negativo")
	}
	return fmt.Sprintf("%s: %d", name, value), nil
}

func printFunctionSignature(fn interface{}) {
	t := reflect.TypeOf(fn)
	
	if t.Kind() != reflect.Func {
		fmt.Println("No es una función")
		return
	}
	
	fmt.Printf("Función: %v\n", t)
	fmt.Printf("  Inputs: %d\n", t.NumIn())
	
	for i := 0; i < t.NumIn(); i++ {
		fmt.Printf("    In(%d): %v\n", i, t.In(i))
	}
	
	fmt.Printf("  Outputs: %d\n", t.NumOut())
	
	for i := 0; i < t.NumOut(); i++ {
		fmt.Printf("    Out(%d): %v\n", i, t.Out(i))
	}
}

func main() {
	fmt.Println("=== add ===")
	printFunctionSignature(add)
	
	fmt.Println("\n=== process ===")
	printFunctionSignature(process)
}
```

### 33.7.2 Llamar Funciones Dinámicamente

```go
package main

import (
	"fmt"
	"reflect"
)

func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("división por cero")
	}
	return a / b, nil
}

func callFunction(fn interface{}, args ...interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(fn)
	t := v.Type()
	
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("no es una función")
	}
	
	if len(args) != t.NumIn() {
		return nil, fmt.Errorf("número de argumentos incorrecto")
	}
	
	// Convertir argumentos
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}
	
	// Llamar función
	results := v.Call(reflectArgs)
	
	// Convertir resultados
	var retvals []interface{}
	for _, r := range results {
		retvals = append(retvals, r.Interface())
	}
	
	return retvals, nil
}

func main() {
	// Llamar multiply
	results, _ := callFunction(multiply, 5, 3)
	fmt.Println("multiply(5, 3) =", results[0])
	
	// Llamar divide
	results, _ = callFunction(divide, 20, 4)
	fmt.Printf("divide(20, 4) = %v, error: %v\n", results[0], results[1])
	
	// Llamar divide con error
	results, _ = callFunction(divide, 20, 0)
	fmt.Printf("divide(20, 0) = %v, error: %v\n", results[0], results[1])
}
```

---

## 33.8 Type Assertions vs Reflect

### 33.8.1 Cuándo Usar Type Assertion

Type assertions son más rápidas y seguras en tiempo de compilación:

```go
package main

import (
	"fmt"
	"reflect"
	"time"
)

type Handler interface{}

// Procesar con Type Assertion (PREFERIDO)
func processWithTypeAssertion(h Handler) {
	switch v := h.(type) {
	case string:
		fmt.Println("String:", v)
	case int:
		fmt.Println("Int:", v)
	case error:
		fmt.Println("Error:", v)
	default:
		fmt.Println("Tipo desconocido")
	}
}

// Procesar con Reflection (MÁS LENTO)
func processWithReflection(h Handler) {
	t := reflect.TypeOf(h)
	fmt.Printf("Tipo: %v\n", t)
}

func benchmark() {
	iterations := 1_000_000
	
	// Benchmark Type Assertion
	start := time.Now()
	for i := 0; i < iterations; i++ {
		processWithTypeAssertion("test")
		processWithTypeAssertion(42)
		processWithTypeAssertion(fmt.Errorf("error"))
	}
	typeAssertionTime := time.Since(start)
	
	// Benchmark Reflection
	start = time.Now()
	for i := 0; i < iterations; i++ {
		processWithReflection("test")
		processWithReflection(42)
		processWithReflection(fmt.Errorf("error"))
	}
	reflectionTime := time.Since(start)
	
	fmt.Printf("Type Assertion:    %v (%d ops)\n", typeAssertionTime, iterations*3)
	fmt.Printf("Reflection:        %v (%d ops)\n", reflectionTime, iterations*3)
	fmt.Printf("Reflection es ~%.0fx más lento\n", float64(reflectionTime)/float64(typeAssertionTime))
}

func main() {
	fmt.Println("=== Type Assertion vs Reflection ===\n")
	
	fmt.Println("--- Performance ---")
	benchmark()
	
	fmt.Println("\n--- Comparación ---")
	fmt.Println("Type Assertion: Mejor para dispatcher/router")
	fmt.Println("Reflection: Necesario para frameworks genéricos")
}
```

### 33.8.2 Comparación: Type Assertion vs Reflect

| Aspecto | Type Assertion | Reflection |
|---------|---|---|
| **Velocidad** | ✅ Nativa (~2-3 ns) | ❌ 100-1000x más lento |
| **Seguridad compilación** | ✅ Detecta tipos en compile-time | ❌ En runtime solo |
| **Legibilidad** | ✅ Claro y directo | ❌ Complejo |
| **Debugging** | ✅ Stack trace útil | ❌ Difícil de seguir |
| **Flexibilidad** | ❌ Debe conocer tipos | ✅ Funciona con cualquier tipo |
| **Uso en frameworks** | ✅ Para router/dispatcher | ✅ Para serialización genérica |

### 33.8.3 Patrones Combinados

```go
package main

import (
	"fmt"
	"reflect"
)

type DataProcessor interface {
	Process(data interface{}) interface{}
}

// Usar Type Assertion primero (rápido), Reflection si es necesario (lento)
func smartProcess(data interface{}) interface{} {
	// Intentar type assertion primero (rápido)
	switch v := data.(type) {
	case string:
		return fmt.Sprintf("Procesando string: %s", v)
	case int:
		return v * 2
	case float64:
		return v * 1.5
	default:
		// Fallback a reflection para tipos complejos
		t := reflect.TypeOf(data)
		if t.Kind() == reflect.Struct {
			return fmt.Sprintf("Struct: %v", t.Name())
		}
		return fmt.Sprintf("Tipo desconocido: %v", t)
	}
}

func main() {
	results := []interface{}{
		smartProcess("hello"),
		smartProcess(21),
		smartProcess(3.14),
		smartProcess(struct{}{}),
	}
	
	for _, r := range results {
		fmt.Println(r)
	}
}
```

---

## 33.9 Casos de Uso - Aplicaciones Prácticas

### 33.9.1 Serialización Genérica (JSON-like)

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
)

// SimpleSerializer: Serializador genérico básico
type SimpleSerializer struct{}

func (s *SimpleSerializer) Serialize(obj interface{}) string {
	v := reflect.ValueOf(obj)
	t := v.Type()
	
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	
	if t.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", obj)
	}
	
	var parts []string
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		
		// Respetar json tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		
		fieldName := field.Name
		if jsonTag != "" && jsonTag != "omitempty" {
			fieldName = strings.Split(jsonTag, ",")[0]
		}
		
		parts = append(parts, fmt.Sprintf("\"%s\":%v", fieldName, value.Interface()))
	}
	
	return "{" + strings.Join(parts, ",") + "}"
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Secret string `json:"-"` // No serializar
}

func main() {
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com", Secret: "pass123"}
	
	s := SimpleSerializer{}
	fmt.Println(s.Serialize(user))
}
```

### 33.9.2 ORM Simple con Reflection

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Model struct {
	ID   int
	Name string
}

type Repository struct{}

// Convertir struct a INSERT statement
func (r *Repository) ToInsertSQL(obj interface{}) (string, []interface{}) {
	v := reflect.ValueOf(obj)
	t := v.Type()
	
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	
	var columns []string
	var values []interface{}
	var placeholders []string
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		
		// Usar tag db si existe
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}
		
		columnName := strings.Split(dbTag, ",")[0]
		columns = append(columns, columnName)
		values = append(values, value.Interface())
		placeholders = append(placeholders, "?")
	}
	
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		t.Name(),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	
	return sql, values
}

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func main() {
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	
	repo := &Repository{}
	sql, values := repo.ToInsertSQL(user)
	
	fmt.Println("SQL:", sql)
	fmt.Println("Values:", values)
}
```

### 33.9.3 Validator Genérico

```go
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Validator struct{}

func (v *Validator) Validate(obj interface{}) error {
	val := reflect.ValueOf(obj)
	typ := val.Type()
	
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}
	
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("no es un struct")
	}
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		
		tags := field.Tag.Get("validate")
		if tags == "" {
			continue
		}
		
		// Parsear reglas
		rules := strings.Split(tags, ",")
		for _, rule := range rules {
			if err := v.validateRule(field.Name, value, rule); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func (v *Validator) validateRule(fieldName string, value reflect.Value, rule string) error {
	rule = strings.TrimSpace(rule)
	
	switch {
	case rule == "required":
		if value.IsZero() {
			return fmt.Errorf("%s es requerido", fieldName)
		}
	case strings.HasPrefix(rule, "min="):
		min, _ := strconv.Atoi(rule[4:])
		if value.Kind() == reflect.String && len(value.String()) < min {
			return fmt.Errorf("%s debe tener al menos %d caracteres", fieldName, min)
		}
	case strings.HasPrefix(rule, "max="):
		max, _ := strconv.Atoi(rule[4:])
		if value.Kind() == reflect.String && len(value.String()) > max {
			return fmt.Errorf("%s no puede tener más de %d caracteres", fieldName, max)
		}
	}
	
	return nil
}

type User struct {
	Name     string `validate:"required,min=3,max=50"`
	Email    string `validate:"required"`
	Age      int    `validate:"required"`
}

func main() {
	validator := &Validator{}
	
	// Válido
	validUser := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	if err := validator.Validate(validUser); err == nil {
		fmt.Println("✓ Usuario válido")
	}
	
	// Inválido: nombre muy corto
	invalidUser := User{Name: "Al", Email: "al@example.com", Age: 25}
	if err := validator.Validate(invalidUser); err != nil {
		fmt.Println("✗", err)
	}
	
	// Inválido: campos vacíos
	emptyUser := User{}
	if err := validator.Validate(emptyUser); err != nil {
		fmt.Println("✗", err)
	}
}
```

### 33.9.4 Dependency Injection

```go
package main

import (
	"fmt"
	"reflect"
)

type Logger interface {
	Log(msg string)
}

type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(msg string) {
	fmt.Println("[LOG]", msg)
}

type Service struct {
	logger Logger `inject:"logger"`
	name   string
}

func (s *Service) Do() {
	s.logger.Log(fmt.Sprintf("Service %s is working", s.name))
}

type Injector struct {
	services map[string]interface{}
}

func (i *Injector) Register(name string, service interface{}) {
	i.services[name] = service
}

func (i *Injector) Inject(obj interface{}) error {
	val := reflect.ValueOf(obj)
	
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("obj debe ser un puntero")
	}
	
	val = val.Elem()
	typ := val.Type()
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("inject")
		
		if tag == "" {
			continue
		}
		
		service, ok := i.services[tag]
		if !ok {
			return fmt.Errorf("servicio no encontrado: %s", tag)
		}
		
		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			return fmt.Errorf("no se puede inyectar en %s", field.Name)
		}
		
		fieldVal.Set(reflect.ValueOf(service))
	}
	
	return nil
}

func main() {
	injector := &Injector{services: make(map[string]interface{})}
	injector.Register("logger", &ConsoleLogger{})
	
	service := &Service{name: "API"}
	injector.Inject(service)
	
	service.Do()
}
```

---

## 33.10 Performance y Limitaciones

### 33.10.1 Análisis de Performance

```go
package main

import (
	"fmt"
	"reflect"
	"time"
)

// Benchmark: acceso directo vs reflection
func benchmarkFieldAccess() {
	type Data struct {
		Field1 int
		Field2 string
		Field3 float64
	}
	
	data := Data{Field1: 42, Field2: "hello", Field3: 3.14}
	iterations := 10_000_000
	
	// Acceso directo
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = data.Field1
		_ = data.Field2
		_ = data.Field3
	}
	directTime := time.Since(start)
	
	// Con reflection
	v := reflect.ValueOf(data)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = v.FieldByName("Field1").Interface()
		_ = v.FieldByName("Field2").Interface()
		_ = v.FieldByName("Field3").Interface()
	}
	reflectTime := time.Since(start)
	
	fmt.Printf("Acceso directo:    %v\n", directTime)
	fmt.Printf("Reflection:        %v\n", reflectTime)
	fmt.Printf("Factor:            %.0fx\n", float64(reflectTime)/float64(directTime))
}

// Benchmark: llamar método directo vs reflection
func benchmarkMethodCall() {
	type Calculator struct{}
	
	func (c *Calculator) Add(a, b int) int {
		return a + b
	}
	
	calc := &Calculator{}
	iterations := 1_000_000
	
	// Llamada directa
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = calc.Add(5, 3)
	}
	directTime := time.Since(start)
	
	// Con reflection
	v := reflect.ValueOf(calc)
	method := v.MethodByName("Add")
	args := []reflect.Value{reflect.ValueOf(5), reflect.ValueOf(3)}
	
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = method.Call(args)
	}
	reflectTime := time.Since(start)
	
	fmt.Printf("Llamada directa:   %v\n", directTime)
	fmt.Printf("Reflection:        %v\n", reflectTime)
	fmt.Printf("Factor:            %.0fx\n", float64(reflectTime)/float64(directTime))
}

func main() {
	fmt.Println("=== Benchmark: Acceso a Campos ===")
	benchmarkFieldAccess()
	
	fmt.Println("\n=== Benchmark: Llamada de Métodos ===")
	benchmarkMethodCall()
}
```

### 33.10.2 Caché de Reflection

Para mejorar performance, cachear los resultados de reflection:

```go
package main

import (
	"fmt"
	"reflect"
	"sync"
)

// CachedType: cachear información de tipos
type CachedType struct {
	Type   reflect.Type
	Fields []reflect.StructField
	Methods map[string]reflect.Method
}

var typeCache = make(map[reflect.Type]*CachedType)
var typeCacheLock sync.RWMutex

func getCachedType(v interface{}) *CachedType {
	t := reflect.TypeOf(v)
	
	typeCacheLock.RLock()
	if cached, ok := typeCache[t]; ok {
		typeCacheLock.RUnlock()
		return cached
	}
	typeCacheLock.RUnlock()
	
	// Crear entrada en caché
	cached := &CachedType{
		Type:    t,
		Methods: make(map[string]reflect.Method),
	}
	
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	// Cachear campos
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			cached.Fields = append(cached.Fields, t.Field(i))
		}
	}
	
	// Cachear métodos
	for i := 0; i < reflect.TypeOf(v).NumMethod(); i++ {
		m := reflect.TypeOf(v).Method(i)
		cached.Methods[m.Name] = m
	}
	
	// Guardar en caché
	typeCacheLock.Lock()
	typeCache[t] = cached
	typeCacheLock.Unlock()
	
	return cached
}

type User struct {
	Name string
	Age  int
}

func main() {
	user := &User{Name: "Alice", Age: 30}
	
	// Primera vez: se crea la entrada en caché
	cached1 := getCachedType(user)
	fmt.Printf("Campos cacheados: %d\n", len(cached1.Fields))
	
	// Segunda vez: se obtiene del caché
	cached2 := getCachedType(user)
	fmt.Printf("Mismo objeto en caché: %v\n", cached1 == cached2)
}
```

### 33.10.3 Limitaciones de Reflect

```go
package main

import (
	"fmt"
	"reflect"
)

func demonstrateLimitations() {
	fmt.Println("=== Limitaciones de Reflection ===\n")
	
	// 1. No se puede crear tipos nuevos
	fmt.Println("1. No se puede crear tipos dinámicamente")
	fmt.Println("   ❌ reflect no puede crear un tipo T nuevo en runtime")
	
	// 2. No se puede acceder a variables locales de otra función
	fmt.Println("\n2. No se puede acceder a estado privado")
	var privateVar = 42
	v := reflect.ValueOf(privateVar)
	fmt.Printf("   Acceso a variable privada local: %v\n", v.Interface())
	fmt.Println("   ✓ Pero solo si la tienes en mano")
	
	// 3. Campos privados no son totalmente accesibles
	fmt.Println("\n3. Campos privados de structs")
	type privatStruct struct {
		public  int
		private int
	}
	
	ps := privatStruct{public: 1, private: 2}
	t := reflect.TypeOf(ps)
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := reflect.ValueOf(ps).Field(i)
		
		fmt.Printf("   Campo '%s': ", field.Name)
		if field.IsExported() {
			fmt.Printf("exportado (público), valor: %d\n", value.Int())
		} else {
			fmt.Printf("privado, no se puede usar CanSet()\n")
		}
	}
	
	// 4. Performance cost
	fmt.Println("\n4. Overhead de performance")
	fmt.Println("   Reflection es 100-1000x más lento que acceso directo")
	fmt.Println("   Usar solo cuando sea absolutamente necesario")
	
	// 5. Errores en runtime
	fmt.Println("\n5. Errores solo en runtime")
	fmt.Println("   Typos en nombres de métodos/campos causan panic")
	fmt.Println("   Compilador no puede detectarlos")
}

func main() {
	demonstrateLimitations()
}
```

---

## 33.11 Buenas Prácticas y Patrones

### 33.11.1 Patrones Seguros

```go
package main

import (
	"fmt"
	"reflect"
)

// PATRÓN 1: Cachear reflect.Type
type TypeInfo struct {
	Type   reflect.Type
	Fields map[string]int // Mapear nombre a índice
}

var typeInfoCache = make(map[string]*TypeInfo)

func getTypeInfo(obj interface{}) (*TypeInfo, error) {
	t := reflect.TypeOf(obj)
	name := t.String()
	
	// Buscar en caché
	if info, ok := typeInfoCache[name]; ok {
		return info, nil
	}
	
	// Crear entrada en caché
	info := &TypeInfo{
		Type:   t,
		Fields: make(map[string]int),
	}
	
	elem := t
	if t.Kind() == reflect.Ptr {
		elem = t.Elem()
	}
	
	if elem.Kind() == reflect.Struct {
		for i := 0; i < elem.NumField(); i++ {
			info.Fields[elem.Field(i).Name] = i
		}
	}
	
	typeInfoCache[name] = info
	return info, nil
}

// PATRÓN 2: Validar antes de usar
func safeSetField(obj interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(obj)
	
	// Validar que sea puntero
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("obj debe ser puntero")
	}
	
	v = v.Elem()
	
	// Validar que sea struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("obj debe ser struct")
	}
	
	// Buscar campo
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("campo no encontrado: %s", fieldName)
	}
	
	// Validar que se puede establecer
	if !field.CanSet() {
		return fmt.Errorf("no se puede establecer campo privado: %s", fieldName)
	}
	
	// Validar tipo
	newVal := reflect.ValueOf(value)
	if !newVal.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("tipo incompatible para %s: %v != %v",
			fieldName, newVal.Type(), field.Type())
	}
	
	field.Set(newVal)
	return nil
}

// PATRÓN 3: Usar interfaces cuando sea posible
type Processor interface {
	Process(data interface{}) interface{}
}

func processItems(items []interface{}, processor Processor) []interface{} {
	results := make([]interface{}, len(items))
	for i, item := range items {
		results[i] = processor.Process(item)
	}
	return results
}

func main() {
	type Config struct {
		Host string
		Port int
	}
	
	cfg := &Config{Host: "localhost", Port: 8080}
	
	// Usar patterns seguros
	if err := safeSetField(cfg, "Host", "0.0.0.0"); err != nil {
		fmt.Println("Error:", err)
	}
	
	fmt.Printf("Config: %+v\n", cfg)
}
```

### 33.11.2 Evitar Antipatrones

```go
package main

import (
	"fmt"
	"reflect"
)

func demonstrateAntipatterns() {
	fmt.Println("=== Antipatrones a Evitar ===\n")
	
	// ANTIPATRÓN 1: Usar reflect para lógica de negocio
	fmt.Println("❌ ANTIPATRÓN 1: Reflection para lógica core")
	fmt.Println("   Malo:")
	fmt.Println(`   if reflect.TypeOf(user).Name() == "Admin" { ... }`)
	fmt.Println("   Bueno:")
	fmt.Println("   if user.IsAdmin() { ... }")
	
	// ANTIPATRÓN 2: No cachear tipos
	fmt.Println("\n❌ ANTIPATRÓN 2: No cachear resultados de reflection")
	fmt.Println("   Malo:")
	fmt.Println(`
   for i := 0; i < 1000000; i++ {
       t := reflect.TypeOf(obj)  // Se repite innecesariamente
       field := t.FieldByName("Field")
   }`)
	fmt.Println("   Bueno: Cachear TypeOf fuera del loop")
	
	// ANTIPATRÓN 3: Ignorar errores
	fmt.Println("\n❌ ANTIPATRÓN 3: Ignorar validaciones")
	fmt.Println("   Malo:")
	fmt.Println(`
   field.Set(value)  // ¿Qué pasa si CanSet es false?`)
	fmt.Println("   Bueno:")
	fmt.Println(`
   if !field.CanSet() {
       return fmt.Errorf("no se puede establecer")
   }
   field.Set(value)`)
	
	// ANTIPATRÓN 4: Usar reflect cuando type assertion funciona
	fmt.Println("\n❌ ANTIPATRÓN 4: Reflection innecesaria")
	fmt.Println("   Malo:")
	fmt.Println(`   if reflect.TypeOf(v).Kind() == reflect.String { ... }`)
	fmt.Println("   Bueno:")
	fmt.Println(`   if str, ok := v.(string); ok { ... }`)
}

func main() {
	demonstrateAntipatterns()
}
```

### 33.11.3 Testing con Reflection

```go
package main

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	ID    int
	Name  string
	Email string
}

// Helper para testing: verificar que struct tiene campos esperados
func TestStructFields(t *testing.T) {
	expectedFields := map[string]reflect.Kind{
		"ID":    reflect.Int,
		"Name":  reflect.String,
		"Email": reflect.String,
	}
	
	typ := reflect.TypeOf(User{})
	
	for name, expectedKind := range expectedFields {
		field, exists := typ.FieldByName(name)
		
		if !exists {
			t.Errorf("Campo esperado no encontrado: %s", name)
			continue
		}
		
		if field.Type.Kind() != expectedKind {
			t.Errorf("Tipo incorrecto para %s: esperado %v, obtenido %v",
				name, expectedKind, field.Type.Kind())
		}
	}
}

// Helper para testing: verificar que struct tiene métodos
func TestStructMethods(t *testing.T) {
	expectedMethods := []string{"String", "Validate"}
	
	typ := reflect.TypeOf((*User)(nil))
	
	for _, methodName := range expectedMethods {
		_, exists := typ.MethodByName(methodName)
		
		if !exists {
			t.Errorf("Método esperado no encontrado: %s", methodName)
		}
	}
}

// Helper para generar error descriptivo en tests
func compareStructs(t *testing.T, got, expected interface{}) {
	vGot := reflect.ValueOf(got)
	vExp := reflect.ValueOf(expected)
	tGot := vGot.Type()
	tExp := vExp.Type()
	
	if tGot != tExp {
		t.Errorf("Tipos diferentes: %v vs %v", tGot, tExp)
		return
	}
	
	if !reflect.DeepEqual(got, expected) {
		t.Logf("Diferencias:\n")
		for i := 0; i < tGot.NumField(); i++ {
			field := tGot.Field(i)
			gotVal := vGot.Field(i).Interface()
			expVal := vExp.Field(i).Interface()
			
			if !reflect.DeepEqual(gotVal, expVal) {
				t.Logf("  %s: %v != %v", field.Name, gotVal, expVal)
			}
		}
	}
}

func main() {
	fmt.Println("Testing con Reflection:")
	fmt.Println("- Validar estructura de structs")
	fmt.Println("- Verificar métodos requeridos")
	fmt.Println("- Comparar objetos complejos")
}
```

---

## Casos de Uso Reales

### Serialización JSON Avanzada
- Marshalling/unmarshalling basado en tags
- Soporte para tipos customizados
- Manejo de valores null opcionales

### ORMs (Object-Relational Mapping)
- Mapeo automático struct ↔ filas de BD
- Generación de queries SQL genéricas
- Escaneado de resultados de BD

### Testing Frameworks
- Reflection para descubrimiento automático de tests
- Assertions genéricas
- Mocking de interfaces

### Dependency Injection
- Inyección automática basada en tags
- Registro y resolución genérica de servicios
- Constructor injection automático

### APIs REST Genéricas
- Routing dinámico basado en métodos
- Serialización automática de responses
- Validación de requests genérica

---

## Resumen de Conceptos Clave

| Concepto | Descripción | Performance |
|----------|-------------|-------------|
| **reflect.TypeOf** | Obtener información de tipo | ~100 ns |
| **reflect.ValueOf** | Obtener valor + tipo | ~50 ns |
| **Value.FieldByName** | Acceso a campo por nombre | ~1000 ns |
| **Type.FieldByName** | Información de campo | ~500 ns |
| **Value.Call** | Llamar método dinámicamente | ~5000 ns |
| **Cachear TypeOf** | Reutilizar tipo en loops | 100x más rápido |
| **Type Assertion** | Alternativa rápida | ~5 ns |
| **reflect.DeepEqual** | Comparar valores | Variable |

---

## Comparativa: Go vs Otros Lenguajes

### Go (Reflection Conservador)
```go
t := reflect.TypeOf(value)
v := reflect.ValueOf(value)
// Explícito, performante, seguro
```

### Python (Reflection Dinámico)
```python
type(value).__name__  # Tipo
vars(value)           # Atributos
# Dinámico pero impredecible en performance
```

### Java (Reflection Detallado)
```java
Class<?> clazz = value.getClass();
Field[] fields = clazz.getDeclaredFields();
// Muy verbose pero muy flexible
```

### C# (Reflection Integrado)
```csharp
var type = value.GetType();
var properties = type.GetProperties();
// Similar a Java pero más limpio
```

---

## Antipatrones a Evitar

❌ **No cachear TypeOf en loops**
```go
for i := 0; i < 1000000; i++ {
    t := reflect.TypeOf(obj)  // Ineficiente!
}
```

✅ **Cachear fuera del loop**
```go
t := reflect.TypeOf(obj)
for i := 0; i < 1000000; i++ {
    // Usar t
}
```

---

❌ **No verificar CanSet antes de Set**
```go
field.Set(value)  // Puede panickear!
```

✅ **Siempre verificar primero**
```go
if field.CanSet() {
    field.Set(value)
}
```

---

❌ **Usar Reflection para todo**
```go
if reflect.TypeOf(user).Name() == "Admin" {
    // Lógica de negocio con reflection
}
```

✅ **Usar interfaces**
```go
type Admin interface {
    IsAdmin() bool
}
if u, ok := user.(Admin); ok && u.IsAdmin() {
    // Lógica de negocio
}
```

---

## Conclusión

`Reflect` es una herramienta poderosa pero debe usarse con cuidado. Es esencial para:
- Frameworks genéricos (JSON, ORM, DI)
- Librerías públicas que soportan tipos arbitrarios
- Testing frameworks

Pero evítala para lógica de negocio crítica. Recuerda:

1. **Performance**: 100-1000x más lento que código directo
2. **Seguridad**: Los errores ocurren en runtime
3. **Legibilidad**: El código reflection es difícil de seguir
4. **Caching**: Siempre cachear resultados de TypeOf/FieldByName
5. **Alternativas**: Considera type assertions e interfaces primero

**Próximos pasos:** Usa reflection prudentemente en tus proyectos, especialmente en capas de framework y serialización. En la lógica de aplicación, prefiere interfaces y type assertions.

---

# Ejercicios Progresivos

## Ejercicio 1: Struct Printer (Básico)

Crear una función `PrintStructAsTable` que use reflection para imprimir cualquier struct en formato tabla.

```go
package main

import (
	"fmt"
	"reflect"
)

type Product struct {
	ID       int
	Name     string
	Price    float64
	InStock  bool
}

func PrintStructAsTable(obj interface{}) {
	// TODO: Implementar
	// 1. Obtener el tipo y valor del objeto
	// 2. Iterar sobre todos los campos
	// 3. Imprimir nombre del campo y valor en formato tabla
	// 4. Manejar punteros correctamente
	
	// Salida esperada:
	// ┌──────────┬───────────────┐
	// │ Campo    │ Valor         │
	// ├──────────┼───────────────┤
	// │ ID       │ 1             │
	// │ Name     │ Laptop        │
	// │ Price    │ 999.99        │
	// │ InStock  │ true          │
	// └──────────┴───────────────┘
}

func main() {
	product := Product{
		ID:      1,
		Name:    "Laptop",
		Price:   999.99,
		InStock: true,
	}
	
	PrintStructAsTable(product)
	PrintStructAsTable(&product)  // También debe funcionar con punteros
}
```

---

## Ejercicio 2: Generic Deserializer (Intermedio)

Crear una función `MapToStruct` que convierta un `map[string]interface{}` a cualquier struct genéricamente.

```go
package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func MapToStruct(data map[string]interface{}, target interface{}) error {
	// TODO: Implementar
	// 1. Validar que target sea un puntero a struct
	// 2. Iterar sobre los campos del struct
	// 3. Buscar el valor en el map
	// 4. Convertir y asignar el valor al campo
	// 5. Retornar error si hay incompatibilidades de tipo
	
	// Casos a manejar:
	// - Campos que existen en el struct pero no en el map
	// - Campos que existen en el map pero no en el struct
	// - Conversiones de tipo (string -> int, etc.)
	
	return nil
}

func main() {
	data := map[string]interface{}{
		"ID":    1,
		"Name":  "Alice",
		"Email": "alice@example.com",
		"Age":   30,
	}
	
	var user User
	if err := MapToStruct(data, &user); err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	fmt.Printf("Usuario: %+v\n", user)
	
	// Caso con conversión: Age como string
	data2 := map[string]interface{}{
		"ID":    2,
		"Name":  "Bob",
		"Email": "bob@example.com",
		"Age":   "25",  // String que debe convertirse a int
	}
	
	var user2 User
	if err := MapToStruct(data2, &user2); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Usuario 2: %+v\n", user2)
	}
}
```

---

## Ejercicio 3: Struct Validator (Intermedio)

Crear un validador que use tags `validate` para validar structs automáticamente.

```go
package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Article struct {
	Title       string `validate:"required,min=5,max=100"`
	Description string `validate:"required,min=10"`
	Author      string `validate:"required,email"`
	Published   bool   `validate:""`
	Views       int    `validate:"min=0"`
}

type ValidationError struct {
	Field   string
	Message string
}

func ValidateStruct(obj interface{}) []ValidationError {
	// TODO: Implementar
	// Tags soportados:
	// - required: campo no puede estar vacío
	// - min=N: mínima longitud (string) o valor (int)
	// - max=N: máxima longitud (string) o valor (int)
	// - email: validar formato email
	// - regex=pattern: validar con regex
	
	// Retornar lista de ValidationError con todos los errores encontrados
	
	return []ValidationError{}
}

func main() {
	// Válido
	validArticle := Article{
		Title:       "Reflection en Go",
		Description: "Una guía completa sobre reflection",
		Author:      "alice@example.com",
		Published:   true,
		Views:       100,
	}
	
	errors := ValidateStruct(validArticle)
	if len(errors) == 0 {
		fmt.Println("✓ Artículo válido")
	}
	
	// Inválido: título muy corto
	invalidArticle := Article{
		Title:       "Go",
		Description: "Corto",
		Author:      "invalid-email",
		Views:       -5,
	}
	
	errors = ValidateStruct(invalidArticle)
	for _, err := range errors {
		fmt.Printf("✗ %s: %s\n", err.Field, err.Message)
	}
}
```

---

## Ejercicio 4: Dynamic Method Caller (Avanzado)

Crear un sistema que permita llamar métodos dinámicamente por nombre, con validación de argumentos.

```go
package main

import (
	"fmt"
	"reflect"
)

type Calculator struct {
	LastResult float64
}

func (c *Calculator) Add(a, b float64) float64 {
	c.LastResult = a + b
	return c.LastResult
}

func (c *Calculator) Multiply(a, b, c float64) float64 {
	c.LastResult = a * b * c
	return c.LastResult
}

func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("división por cero")
	}
	c.LastResult = a / b
	return c.LastResult, nil
}

func (c *Calculator) GetResult() float64 {
	return c.LastResult
}

type MethodCall struct {
	Object     interface{}
	MethodName string
	Args       []interface{}
}

type MethodResult struct {
	ReturnValues []interface{}
	Error        error
}

func CallMethod(call MethodCall) MethodResult {
	// TODO: Implementar
	// 1. Validar que el método existe
	// 2. Validar número y tipos de argumentos
	// 3. Llamar el método
	// 4. Retornar valores o error
	
	// Casos especiales:
	// - Múltiples retornos
	// - Retorno con error
	// - Validación de tipos
	
	return MethodResult{Error: fmt.Errorf("no implementado")}
}

func main() {
	calc := &Calculator{}
	
	// Llamar Add
	result := CallMethod(MethodCall{
		Object:     calc,
		MethodName: "Add",
		Args:       []interface{}{10.0, 5.0},
	})
	
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
	} else {
		fmt.Println("Add(10, 5) =", result.ReturnValues[0])
	}
	
	// Llamar Divide con error
	result = CallMethod(MethodCall{
		Object:     calc,
		MethodName: "Divide",
		Args:       []interface{}{20.0, 0.0},
	})
	
	if result.Error != nil {
		fmt.Println("Divide error:", result.Error)
	}
	
	// Llamar con número de args incorrecto
	result = CallMethod(MethodCall{
		Object:     calc,
		MethodName: "Add",
		Args:       []interface{}{10.0},  // Falta un argumento
	})
	
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
	}
}
```

---

## Ejercicio 5: Mini ORM (Avanzado+)

Crear un mini ORM que use reflection para generar queries SQL y mappear resultados.

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	ID    int    `db:"id,pk"`
	Name  string `db:"name,notnull"`
	Email string `db:"email,unique"`
	Age   int    `db:"age"`
}

type Post struct {
	ID        int    `db:"id,pk"`
	UserID    int    `db:"user_id,fk"`
	Title     string `db:"title"`
	Content   string `db:"content"`
	Published bool   `db:"published"`
}

type QueryBuilder struct {
	table     string
	columns   []string
	values    []interface{}
	where     string
	whereArgs []interface{}
}

// GenerateInsertSQL: generar INSERT statement
func GenerateInsertSQL(obj interface{}) (string, []interface{}) {
	// TODO: Implementar
	// 1. Obtener nombre de tabla del struct
	// 2. Extraer campos con tag db
	// 3. Generar INSERT statement con placeholders
	// 4. Retornar SQL y valores
	
	// Ejemplo:
	// INSERT INTO user (id, name, email, age) VALUES (?, ?, ?, ?)
	
	return "", []interface{}{}
}

// GenerateSelectSQL: generar SELECT statement
func GenerateSelectSQL(obj interface{}, where map[string]interface{}) string {
	// TODO: Implementar
	// 1. Obtener columnas del struct
	// 2. Construir WHERE con condiciones AND
	// 3. Retornar SELECT statement
	
	// Ejemplo:
	// SELECT id, name, email, age FROM user WHERE name = ? AND age > ?
	
	return ""
}

// MapRowToStruct: mappear resultado de query a struct
func MapRowToStruct(columns []string, values []interface{}, target interface{}) error {
	// TODO: Implementar
	// 1. Iterar sobre columnas
	// 2. Encontrar campos correspondientes en struct
	// 3. Asignar valores
	
	return nil
}

// GetTableName: obtener nombre de tabla
func GetTableName(obj interface{}) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}

func main() {
	// Generar INSERT para User
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30}
	sql, values := GenerateInsertSQL(user)
	
	fmt.Println("INSERT SQL:", sql)
	fmt.Println("Values:", values)
	
	// Generar SELECT con WHERE
	selectSQL := GenerateSelectSQL(user, map[string]interface{}{
		"name": "Alice",
		"age":  30,
	})
	
	fmt.Println("\nSELECT SQL:", selectSQL)
	
	// Generar INSERT para Post
	post := Post{ID: 1, UserID: 1, Title: "Hello", Content: "World", Published: true}
	sql, values = GenerateInsertSQL(post)
	
	fmt.Println("\nINSERT POST SQL:", sql)
	fmt.Println("Values:", values)
}
```

---

## Resumen de Conceptos Clave

| Concepto | Descripción | Ejemplo |
|----------|-------------|---------|
| **reflect.TypeOf** | Obtener tipo estático | `reflect.TypeOf(var)` |
| **reflect.ValueOf** | Obtener valor runtime | `reflect.ValueOf(var)` |
| **Type.Kind()** | Clasificación de tipo | `int`, `string`, `struct` |
| **Type.NumField()** | Número de campos struct | Iterar con campo |
| **Type.Field(i)** | Información i-ésimo campo | Acceder a Name, Type, Tag |
| **Value.Field(i)** | Obtener valor i-ésimo campo | Leer/escribir valor |
| **Value.FieldByName** | Campo por nombre | Acceso dinámico |
| **Value.CanSet()** | ¿Se puede modificar? | Validar antes de Set |
| **Value.Set()** | Establecer valor | Modificar valor |
| **Type.Method(i)** | i-ésimo método | Información de método |
| **Value.MethodByName** | Método por nombre | Llamar dinámicamente |
| **Value.Call()** | Llamar función/método | Pasar argumentos como []Value |

---

## Buenas Prácticas Finales

1. **Cachear TypeOf**: No llamar en loops
2. **Validar tipos**: Siempre verificar Kind antes de operaciones
3. **Manejo de errores**: Reflection puede fallar en runtime
4. **Performance**: Considerar alternativas (type assertions, interfaces)
5. **Documentación**: Código reflection necesita comentarios explicativos
6. **Testing**: Usar DeepEqual para comparar en tests
7. **Interfaces primero**: Diseñar con interfaces, usar reflection solo cuando sea necesario

**Próximos pasos:** Integra reflection en proyectos reales para JSON APIs, ORMs simples, o frameworks de testing. ¡Pero usa con moderación!
