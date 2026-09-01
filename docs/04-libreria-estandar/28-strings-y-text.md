# Capítulo 28: Strings y text - Manipulación de texto

## 📋 Tabla de Contenidos
1. [Strings Package - Overview](#281-strings-package---overview)
2. [Búsqueda y Localización](#282-búsqueda-y-localización)
3. [Transformación de Strings](#283-transformación-de-strings)
4. [Splitting y Joining](#284-splitting-y-joining)
5. [Sustitución de Strings](#285-sustitución-de-strings)
6. [Comparación de Strings](#286-comparación-de-strings)
7. [Regular Expressions](#287-regular-expressions)
8. [Submatch y Groups](#288-submatch-y-groups)
9. [Text Templates](#289-text-templates)
10. [HTML Templates](#2810-html-templates)
11. [Buenas Prácticas y Patrones](#2811-buenas-prácticas-y-patrones)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 28.1 Strings Package - Overview

### Conceptos Fundamentales

Go trata los strings como **secuencias inmutables de bytes**. Aunque en Go la mayoría de strings contienen texto UTF-8, internamente son tipos `[]byte`. Esta distinción es crucial para entender cómo funcionan las operaciones de strings.

```
┌─────────────────────────────────────────┐
│ String en Go (tipo immutable)            │
├─────────────────────────────────────────┤
│ • Tipo base: array de bytes             │
│ • Codificación: UTF-8 por defecto      │
│ • Longitud: número de bytes, no runes  │
│ • Acceso: [index] devuelve byte        │
│ • Rango: for...range itera runes       │
└─────────────────────────────────────────┘
```

### Package strings vs strconv

El package `strings` maneja operaciones **lingüísticas** (búsqueda, transformación, formatting):

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Operaciones comunes
	s := "Hello, World!"
	
	fmt.Println(strings.Contains(s, "World"))        // true
	fmt.Println(strings.ToUpper(s))                  // HELLO, WORLD!
	fmt.Println(strings.Index(s, "World"))           // 7
	fmt.Println(strings.Split(s, ","))              // [Hello  World!]
	fmt.Println(strings.Trim(s, "!"))               // Hello, World
}
```

El package `strconv` maneja **conversión de tipos**:

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Conversión numérica
	i, _ := strconv.Atoi("42")                      // string → int
	s := strconv.Itoa(42)                           // int → string
	f, _ := strconv.ParseFloat("3.14", 64)          // string → float64
	
	fmt.Println(i, s, f)
}
```

### Funciones Principales del Package strings

| Función | Propósito | Ejemplo |
|---------|-----------|---------|
| `Contains(s, substr string) bool` | Verifica existencia | `strings.Contains("hello", "ell")` |
| `Index(s, substr string) int` | Encuentra posición | `strings.Index("hello", "l")` |
| `Split(s, sep string) []string` | Divide en partes | `strings.Split("a,b,c", ",")` |
| `Join(a []string, sep string) string` | Une partes | `strings.Join([]string{"a","b"}, ",")` |
| `ToUpper(s string) string` | Convierte a mayúsculas | `strings.ToUpper("hello")` |
| `TrimSpace(s string) string` | Elimina espacios | `strings.TrimSpace(" hello ")` |
| `Replace(s, old, new string, n int) string` | Reemplaza n instancias | `strings.Replace("aaa", "a", "b", 2)` |
| `Count(s, substr string) int` | Cuenta ocurrencias | `strings.Count("hello", "l")` |

---

## 28.2 Búsqueda y Localización

### Contains: Búsqueda Básica

`Contains` verifica si un string contiene un substring:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "Go programming language"
	
	// Búsqueda case-sensitive
	fmt.Println(strings.Contains(text, "Go"))        // true
	fmt.Println(strings.Contains(text, "go"))        // false
	fmt.Println(strings.Contains(text, "gram"))      // true
	
	// Búsqueda de múltiples términos
	searchTerms := []string{"Go", "Python", "Java"}
	for _, term := range searchTerms {
		if strings.Contains(text, term) {
			fmt.Printf("%s encontrado\n", term)
		}
	}
}
```

### Index: Localización de Posición

`Index` devuelve el índice (posición en bytes) de la primera ocurrencia:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "The quick brown fox jumps"
	
	// Index y LastIndex
	fmt.Println(strings.Index(text, "o"))           // 12 (en "brown")
	fmt.Println(strings.LastIndex(text, "o"))       // 20 (en "fox")
	fmt.Println(strings.Index(text, "xyz"))         // -1 (no encontrado)
	
	// IndexAny: busca el primer byte de cualquiera en charset
	fmt.Println(strings.IndexAny(text, "aeiou"))    // 2 (primera vocal "e")
	
	// IndexRune: busca un rune específico
	fmt.Println(strings.IndexRune(text, 'o'))       // 12
	
	// Caso de uso: parsear hasta un delimitador
	line := "key=value&name=John"
	if idx := strings.Index(line, "="); idx >= 0 {
		key := line[:idx]
		rest := line[idx+1:]
		fmt.Printf("Key: %s, Rest: %s\n", key, rest)
	}
}
```

### HasPrefix y HasSuffix: Verificar Inicio/Fin

Funciones especializadas para verificar el principio y fin:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	filename := "document.pdf"
	url := "https://example.com"
	
	// Verificar extensión
	if strings.HasSuffix(filename, ".pdf") {
		fmt.Println("Es un PDF")
	}
	
	// Verificar protocolo
	if strings.HasPrefix(url, "https://") {
		fmt.Println("Conexión segura")
	}
	
	// Validar identificadores de API
	apiKey := "sk_live_1234567890"
	if strings.HasPrefix(apiKey, "sk_") {
		fmt.Println("Clave de API válida")
	}
	
	// Filtrar archivos
	extensions := []string{".go", ".txt", ".md"}
	files := []string{"main.go", "README.md", "photo.jpg"}
	
	for _, file := range files {
		for _, ext := range extensions {
			if strings.HasSuffix(file, ext) {
				fmt.Printf("✓ %s\n", file)
				break
			}
		}
	}
}
```

### Count: Contar Ocurrencias

Contar cuántas veces aparece un substring:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "The Mississippi River is mighty"
	
	// Contar substring
	fmt.Println(strings.Count(text, "i"))           // 4
	fmt.Println(strings.Count(text, "is"))          // 2
	fmt.Println(strings.Count(text, "Mississippi")) // 1
	
	// Caso de uso: validar densidad de palabras
	content := "hello hello world hello"
	word := "hello"
	count := strings.Count(content, word)
	density := float64(count) / float64(strings.Count(content, " ")+1)
	
	fmt.Printf("Palabra '%s' aparece %d veces\n", word, count)
	fmt.Printf("Densidad: %.2f%%\n", density*100)
}
```

---

## 28.3 Transformación de Strings

### ToUpper y ToLower: Cambio de Caso

Conversiones básicas de mayúsculas/minúsculas:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "Hello, Go Programming!"
	
	// Conversiones básicas
	fmt.Println(strings.ToUpper(text))              // HELLO, GO PROGRAMMING!
	fmt.Println(strings.ToLower(text))              // hello, go programming!
	
	// Comparación case-insensitive
	password := "MyPassword123"
	if strings.ToLower(password) == strings.ToLower("mypassword123") {
		fmt.Println("Las contraseñas coinciden")
	}
	
	// Generar búsqueda flexible
	searchQuery := "golang"
	documents := []string{"Golang Guide", "GO language", "GoLang Book"}
	
	for _, doc := range documents {
		if strings.Contains(strings.ToLower(doc), strings.ToLower(searchQuery)) {
			fmt.Printf("Encontrado: %s\n", doc)
		}
	}
}
```

### Title y Unicode Case Mapping

Go también proporciona transformaciones Unicode más complejas:

```go
package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	text := "hello world"
	
	// Title: primera letra de cada palabra en mayúscula
	fmt.Println(strings.Title(text))                // Hello World
	// Nota: strings.Title está deprecado en Go 1.18+, usar unicode/cases
	
	// Para Unicode correcto, usar strings.ToTitle
	fmt.Println(strings.ToTitle(text))              // HELLO WORLD
	
	// Procesar rune por rune para transformación personalizada
	title := titleCase(text)
	fmt.Println(title)                              // Hello World
}

// Capitalizar manualmente cada palabra
func titleCase(s string) string {
	runes := []rune(strings.ToLower(s))
	capitalize := true
	
	for i, r := range runes {
		if unicode.IsSpace(r) {
			capitalize = true
		} else if capitalize {
			runes[i] = unicode.ToUpper(r)
			capitalize = false
		}
	}
	
	return string(runes)
}
```

### TrimSpace, Trim, TrimPrefix, TrimSuffix

Eliminar caracteres de strings:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// TrimSpace: elimina espacios en blanco
	text := "  hello world  \n"
	fmt.Printf("'%s'\n", strings.TrimSpace(text))   // 'hello world'
	
	// Trim: elimina caracteres específicos de ambos lados
	path := "///home/user///"
	fmt.Println(strings.Trim(path, "/"))            // home/user
	
	// TrimLeft y TrimRight: solo un lado
	fmt.Println(strings.TrimLeft(path, "/"))        // home/user///
	fmt.Println(strings.TrimRight(path, "/"))       // ///home/user
	
	// TrimPrefix y TrimSuffix: elimina solo si coincide
	url := "https://example.com"
	domain := strings.TrimPrefix(url, "https://")
	fmt.Println(domain)                             // example.com
	
	filename := "report.pdf"
	name := strings.TrimSuffix(filename, ".pdf")
	fmt.Println(name)                               // report
	
	// Caso de uso: limpiar entrada de usuario
	userInput := "\n  hello world  \t"
	cleaned := strings.TrimSpace(userInput)
	fmt.Printf("Input limpiado: '%s'\n", cleaned)
	
	// Caso de uso: procesar paths
	paths := []string{"/usr/bin/", "/home/user", "/tmp///"}
	for _, p := range paths {
		clean := strings.Trim(p, "/")
		fmt.Printf("%s → %s\n", p, clean)
	}
}
```

---

## 28.4 Splitting y Joining

### Split: Dividir Strings

`Split` divide un string en una slice basado en un separador:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Split básico
	csv := "apple,banana,orange"
	parts := strings.Split(csv, ",")
	fmt.Println(parts)                              // [apple banana orange]
	
	// Split en palabras con Fields (automático para espacios en blanco)
	text := "The quick brown fox"
	words := strings.Fields(text)
	fmt.Println(words)                              // [The quick brown fox]
	
	// Fields es más inteligente: elimina espacios en blanco múltiples
	text2 := "  hello   world  \t\n  test  "
	fmt.Println(strings.Fields(text2))              // [hello world test]
	
	// Split con límite
	data := "a:b:c:d:e"
	fmt.Println(strings.SplitN(data, ":", 3))       // [a b c:d:e]
	fmt.Println(strings.Split(data, ":"))           // [a b c d e]
	
	// SplitAfter: incluye el separador
	fmt.Println(strings.SplitAfter(data, ":"))      // [a: b: c: d: e]
}
```

### Casos de Uso: Parsing CSV Manual

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Parsing simple de CSV (sin librerías especializadas)
	csvLine := "John,30,Engineer,São Paulo"
	fields := strings.Split(csvLine, ",")
	
	record := map[string]string{
		"name":     fields[0],
		"age":      fields[1],
		"role":     fields[2],
		"location": fields[3],
	}
	
	fmt.Println(record)
	
	// Procesar múltiples líneas
	csvData := `name,age,role
John,30,Engineer
Mary,28,Designer
Bob,35,Manager`
	
	lines := strings.Split(csvData, "\n")
	header := strings.Split(lines[0], ",")
	
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		values := strings.Split(line, ",")
		for i, value := range values {
			fmt.Printf("%s: %s, ", header[i], value)
		}
		fmt.Println()
	}
}
```

### Join: Unir Strings

`Join` combina una slice de strings con un separador:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Join básico
	parts := []string{"apple", "banana", "orange"}
	joined := strings.Join(parts, ", ")
	fmt.Println(joined)                             // apple, banana, orange
	
	// Construir URL
	pathSegments := []string{"", "api", "v1", "users", "123"}
	url := strings.Join(pathSegments, "/")
	fmt.Println(url)                                // /api/v1/users/123
	
	// Construir sentencia SQL (CUIDADO: SQL injection!)
	// ⚠️ ANTIPATRÓN: No hacer esto en producción
	ids := []string{"1", "2", "3"}
	// INCORRECTO: query := "SELECT * FROM users WHERE id IN (" + strings.Join(ids, ",") + ")"
	
	// Construir comando shell seguro
	cmd := "docker"
	args := []string{"run", "-d", "-p", "8080:80", "nginx"}
	fullCmd := cmd + " " + strings.Join(args, " ")
	fmt.Println(fullCmd)                            // docker run -d -p 8080:80 nginx
	
	// Serializar configuración
	config := []string{"debug=true", "port=8080", "host=localhost"}
	output := strings.Join(config, "; ")
	fmt.Println(output)                             // debug=true; port=8080; host=localhost
}
```

### Strings Builder: Construcción Eficiente

Para múltiples concatenaciones, usar `strings.Builder`:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// INCORRECTO: usando + (crea copias cada vez)
	var result string
	for i := 0; i < 1000; i++ {
		result += fmt.Sprintf("Line %d\n", i)      // ⚠️ O(n²) complejidad
	}
	
	// CORRECTO: usando strings.Builder
	var builder strings.Builder
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(&builder, "Line %d\n", i)      // O(n) complejidad
	}
	result := builder.String()
	fmt.Printf("Construido %d líneas\n", len(strings.Split(result, "\n")))
	
	// Ejemplo: generar HTML
	htmlBuilder := generateHTML([]string{"Home", "About", "Contact"})
	fmt.Println(htmlBuilder)
}

func generateHTML(items []string) string {
	var builder strings.Builder
	builder.WriteString("<ul>\n")
	
	for _, item := range items {
		builder.WriteString("  <li>" + item + "</li>\n")
	}
	
	builder.WriteString("</ul>")
	return builder.String()
}
```

---

## 28.5 Sustitución de Strings

### Replace y ReplaceAll

Reemplazar substrings:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "hello hello world hello"
	
	// Replace: reemplaza las primeras n ocurrencias
	fmt.Println(strings.Replace(text, "hello", "hi", 1))
	// Output: hi hello world hello
	
	fmt.Println(strings.Replace(text, "hello", "hi", 2))
	// Output: hi hi world hello
	
	// ReplaceAll: reemplaza todas las ocurrencias
	fmt.Println(strings.ReplaceAll(text, "hello", "hi"))
	// Output: hi hi world hi
	
	// Replace con -1 = ReplaceAll
	fmt.Println(strings.Replace(text, "hello", "hi", -1))
	// Output: hi hi world hi
}
```

### Repeat: Repetir Strings

Repetir un string n veces:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Repetir simple
	fmt.Println(strings.Repeat("*", 10))            // **********
	fmt.Println(strings.Repeat("ab", 3))            // ababab
	
	// Generar línea de separación
	separator := strings.Repeat("=", 50)
	fmt.Println(separator)
	
	// Indentación
	indent := strings.Repeat("  ", 4)               // 4 niveles (8 espacios)
	fmt.Println(indent + "Contenido indentado")
	
	// Generar patrón
	pattern := strings.Repeat("*-", 5) + "*"        // *-*-*-*-*-*
	fmt.Println(pattern)
	
	// Relleno dinámico
	func fillLine(text string, width int) {
		padding := width - len(text)
		if padding > 0 {
			filled := text + strings.Repeat(".", padding)
			fmt.Println(filled)
		}
	}("Name", 20)                                   // Name................
}
```

### Replacer: Reemplazos Múltiples Eficientes

Para múltiples reemplazos, `strings.Replacer` es más eficiente:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// INCORRECTO: múltiples Replace (ineficiente)
	text := "hello world, hello universe"
	// ⚠️ Cada Replace crea una nueva copia
	result := strings.ReplaceAll(text, "hello", "hi")
	result = strings.ReplaceAll(result, "world", "planeta")
	result = strings.ReplaceAll(result, "universe", "universo")
	fmt.Println(result)
	
	// CORRECTO: usar Replacer
	replacer := strings.NewReplacer(
		"hello", "hi",
		"world", "planeta",
		"universe", "universo",
	)
	result = replacer.Replace(text)
	fmt.Println(result)                             // hi planeta, hi universo
	
	// Ejemplo: templating simple
	template := "Hello @NAME, you have @COUNT messages"
	templateReplacer := strings.NewReplacer(
		"@NAME", "John",
		"@COUNT", "5",
	)
	rendered := templateReplacer.Replace(template)
	fmt.Println(rendered)                           // Hello John, you have 5 messages
	
	// Ejemplo: sanitizar HTML básico
	htmlReplacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"&", "&amp;",
	)
	userInput := "<script>alert('xss')</script>"
	safe := htmlReplacer.Replace(userInput)
	fmt.Println(safe)                               // &lt;script&gt;alert('xss')&lt;/script&gt;
}
```

---

## 28.6 Comparación de Strings

### Compare: Comparación Lexicográfica

Comparar strings de forma lexicográfica:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Compare: retorna -1 (menor), 0 (igual), 1 (mayor)
	fmt.Println(strings.Compare("apple", "apple"))  // 0
	fmt.Println(strings.Compare("apple", "banana")) // -1
	fmt.Println(strings.Compare("zebra", "apple"))  // 1
	
	// Ordenar strings personalizados
	words := []string{"zebra", "apple", "mango", "banana"}
	// Usar sort.Slice con Compare
}
```

### EqualFold: Comparación Case-Insensitive

Para comparación insensible a mayúsculas/minúsculas de forma correcta con Unicode:

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// EqualFold: comparación case-insensitive correcta con Unicode
	fmt.Println(strings.EqualFold("Hello", "hello"))       // true
	fmt.Println(strings.EqualFold("Straße", "STRASSE"))    // true (alemán)
	fmt.Println(strings.EqualFold("Μήλο", "μήλο"))         // true (griego)
	
	// Comparación de cabeceras HTTP (siempre case-insensitive)
	func headerMatch(name string) string {
		headers := map[string]string{
			"Content-Type":   "application/json",
			"Accept-Charset": "utf-8",
		}
		
		for key, value := range headers {
			if strings.EqualFold(name, key) {
				return value
			}
		}
		return ""
	}
	
	fmt.Println(headerMatch("content-type"))               // application/json
	fmt.Println(headerMatch("CONTENT-TYPE"))               // application/json
	fmt.Println(headerMatch("Content-Type"))               // application/json
}
```

### Casos de Uso: Búsqueda Case-Insensitive

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Búsqueda case-insensitive
	func caseInsensitiveSearch(haystack, needle string) bool {
		return strings.Contains(
			strings.ToLower(haystack),
			strings.ToLower(needle),
		)
	}
	
	documents := []string{"Go Language", "GOLANG", "golang guide"}
	search := "GoLang"
	
	for _, doc := range documents {
		if caseInsensitiveSearch(doc, search) {
			fmt.Printf("Encontrado: %s\n", doc)
		}
	}
	
	// Validar comandos case-insensitive
	func processCommand(cmd string) {
		switch strings.ToLower(strings.TrimSpace(cmd)) {
		case "quit":
			fmt.Println("Saliendo...")
		case "help":
			fmt.Println("Mostrando ayuda...")
		default:
			fmt.Println("Comando desconocido")
		}
	}
	
	processCommand("QUIT")
	processCommand("  Help  ")
}
```

---

## 28.7 Regular Expressions

### Fundamentos: Sintaxis de Regex en Go

Go usa la librería `regexp` que implementa **RE2** (Google's regular expression library). Características:

```
┌────────────────────────────────────┐
│ Características del Regex en Go     │
├────────────────────────────────────┤
│ • Implementación: RE2 (Google)     │
│ • Evita: backtracking catastrófico │
│ • Tiempo: O(n) garantizado         │
│ • Limitaciones: sin lookbehind     │
│ • Ventaja: seguro en producción    │
└────────────────────────────────────┘
```

### Patrones Básicos

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Crear regex
	re := regexp.MustCompile(`\d+`)                 // Una o más dígitos
	
	// Operadores básicos
	patterns := map[string]string{
		`.` :          "Cualquier carácter",
		`\d` :         "Dígito [0-9]",
		`\w` :         "Palabra [a-zA-Z0-9_]",
		`\s` :         "Whitespace",
		`[abc]` :      "a, b, o c",
		`[^abc]` :     "No a, b, o c",
		`a*` :         "0 o más 'a'",
		`a+` :         "1 o más 'a'",
		`a?` :         "0 o 1 'a'",
		`a{3}` :       "Exactamente 3 'a'",
		`a{3,}` :      "3 o más 'a'",
		`a{3,5}` :     "Entre 3 y 5 'a'",
		`^` :          "Inicio de línea",
		`$` :          "Fin de línea",
		`|` :          "OR (alternancia)",
	}
	
	_ = patterns  // Referencia
	
	// Ejemplos prácticos
	text := "I have 3 apples and 25 oranges"
	
	// Encontrar números
	re = regexp.MustCompile(`\d+`)
	matches := re.FindAllString(text, -1)
	fmt.Println(matches)                            // [3 25]
	
	// Validar formato (email simple)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	emails := []string{"user@example.com", "invalid.email", "test@domain.co.uk"}
	for _, email := range emails {
		if emailRegex.MatchString(email) {
			fmt.Printf("✓ %s\n", email)
		} else {
			fmt.Printf("✗ %s\n", email)
		}
	}
}
```

### Compilación de Regex

Dos formas de crear regex:

```go
package main

import (
	"fmt"
	"regexp"
	"log"
)

func main() {
	// Opción 1: MustCompile (panics si hay error)
	// Usar cuando el patrón es una constante conocida
	re1 := regexp.MustCompile(`\d+`)
	fmt.Println(re1.FindString("abc123def"))        // 123
	
	// Opción 2: Compile (devuelve error)
	// Usar cuando el patrón viene de entrada del usuario
	pattern := `\d+`  // Podría venir de stdin
	re2, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Patrón inválido:", err)
	}
	fmt.Println(re2.FindString("abc123def"))        // 123
	
	// RENDIMIENTO: Compilar una vez, reutilizar muchas veces
	re := regexp.MustCompile(`^[a-z]+$`)
	
	// INCORRECTO: compilar cada vez (ineficiente)
	func isLowercase_bad(s string) bool {
		matched, _ := regexp.MatchString(`^[a-z]+$`, s)  // ⚠️ compila cada vez
		return matched
	}
	
	// CORRECTO: compilar una vez
	func isLowercase_good(s string) bool {
		return re.MatchString(s)                         // ✓ reutiliza compilación
	}
	
	fmt.Println(isLowercase_good("hello"))              // true
	fmt.Println(isLowercase_good("Hello"))              // false
}
```

### Métodos Principales de regexp

| Método | Propósito | Ejemplo |
|--------|-----------|---------|
| `MatchString(s)` | ¿Hay coincidencia? | `re.MatchString("abc123")` |
| `FindString(s)` | Primera coincidencia | `re.FindString("abc123def")` |
| `FindAllString(s, n)` | Todas las coincidencias | `re.FindAllString("a1b2c3", -1)` |
| `ReplaceAllString(s, repl)` | Reemplazar todas | `re.ReplaceAllString("a1a2", "x")` |
| `Split(s, n)` | Dividir por patrón | `re.Split("a1b2c", -1)` |

---

## 28.8 Submatch y Groups

### FindSubmatch: Extraer Grupos

Capturar partes específicas del texto:

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Patrón con grupos de captura (paréntesis)
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	
	email := "john@example.com"
	
	// FindStringSubmatch devuelve el match completo + grupos
	submatch := re.FindStringSubmatch(email)
	if submatch != nil {
		fmt.Println("Match completo:", submatch[0])  // john@example.com
		fmt.Println("Grupo 1 (usuario):", submatch[1]) // john
		fmt.Println("Grupo 2 (dominio):", submatch[2]) // example
		fmt.Println("Grupo 3 (TLD):", submatch[3])     // com
	}
	
	// Procesar múltiples emails
	emails := []string{
		"alice@company.org",
		"bob@domain.co.uk",
		"invalid.email",
	}
	
	for _, email := range emails {
		parts := re.FindStringSubmatch(email)
		if parts != nil {
			fmt.Printf("Usuario: %s, Dominio: %s.%s\n", parts[1], parts[2], parts[3])
		}
	}
}
```

### FindAllSubmatch: Múltiples Coincidencias con Grupos

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Extraer pares clave=valor
	re := regexp.MustCompile(`(\w+)=(\w+)`)
	text := "name=John age=30 city=NYC"
	
	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		fmt.Printf("Clave: %s, Valor: %s\n", match[1], match[2])
	}
	// Output:
	// Clave: name, Valor: John
	// Clave: age, Valor: 30
	// Clave: city, Valor: NYC
	
	// Mapeo automático
	data := parseKeyValues(text)
	fmt.Println(data)                               // map[age:30 city:NYC name:John]
}

func parseKeyValues(text string) map[string]string {
	re := regexp.MustCompile(`(\w+)=(\w+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	
	result := make(map[string]string)
	for _, match := range matches {
		result[match[1]] = match[2]
	}
	return result
}
```

### Named Groups: Grupos Nombrados

Para mayor claridad, usar nombres en lugar de índices:

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Patrón con grupos nombrados
	re := regexp.MustCompile(`(?P<date>\d{4}-\d{2}-\d{2}) (?P<time>\d{2}:\d{2}:\d{2}) (?P<level>\w+): (?P<message>.+)`)
	
	logLine := "2024-01-15 14:30:45 ERROR: Database connection failed"
	
	match := re.FindStringSubmatch(logLine)
	if match != nil {
		// Obtener nombres de grupos
		groupNames := re.SubexpNames()
		result := make(map[string]string)
		
		for i, name := range groupNames {
			if name != "" {
				result[name] = match[i]
			}
		}
		
		fmt.Println(result)
		// Output: map[date:2024-01-15 level:ERROR message:Database connection failed time:14:30:45]
		
		fmt.Printf("Fecha: %s\n", result["date"])
		fmt.Printf("Hora: %s\n", result["time"])
		fmt.Printf("Nivel: %s\n", result["level"])
		fmt.Printf("Mensaje: %s\n", result["message"])
	}
}
```

### Validación: Email y URL

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Regex para email (simplificado)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	// Regex para URL (básico)
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	
	// Regex para teléfono (formato: +XX XXX-XXXX-XXXX)
	phoneRegex := regexp.MustCompile(`^\+\d{1,3}\s\d{3}-\d{4}-\d{4}$`)
	
	testCases := map[string]struct {
		value string
		re    *regexp.Regexp
		name  string
	}{
		"email1": {"user@example.com", emailRegex, "email"},
		"email2": {"invalid.email", emailRegex, "email"},
		"url1":   {"https://example.com", urlRegex, "url"},
		"url2":   {"example.com", urlRegex, "url"},
		"phone1": {"+1 555-1234-5678", phoneRegex, "phone"},
		"phone2": {"555-1234-5678", phoneRegex, "phone"},
	}
	
	for testName, tc := range testCases {
		valid := tc.re.MatchString(tc.value)
		status := "✓"
		if !valid {
			status = "✗"
		}
		fmt.Printf("%s %s (%s): %s\n", status, testName, tc.name, tc.value)
	}
}
```

---

## 28.9 Text Templates

### Fundamentos: text/template Package

Templates permiten generar texto dinámico con lógica:

```go
package main

import (
	"fmt"
	"log"
	"text/template"
)

func main() {
	// Template simple
	tmpl, err := template.New("hello").Parse("Hello, {{.Name}}!")
	if err != nil {
		log.Fatal(err)
	}
	
	data := map[string]string{
		"Name": "World",
	}
	
	err = tmpl.Execute(os.Stdout, data)             // Hello, World!
	if err != nil {
		log.Fatal(err)
	}
}
```

### Sintaxis Básica de Templates

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

func main() {
	// Definir datos
	type Person struct {
		Name  string
		Age   int
		Email string
	}
	
	person := Person{"John Doe", 30, "john@example.com"}
	
	// Operaciones básicas
	templates := map[string]string{
		"simple":       "Hola {{.Name}}",
		"nested":       "{{.Name}} tiene {{.Age}} años",
		"field_access": "Email: {{.Email}}",
		"if":           "{{if gt .Age 18}}Mayor{{else}}Menor{{end}}",
		"and_or":       "{{if and .Name .Email}}Datos completos{{end}}",
		"with":         "{{with .Name}}El nombre es {{.}}{{end}}",
		"range":        "{{range .Hobbies}}Hobby: {{.}} {{end}}",
	}
	
	// Compilar y ejecutar
	for name, pattern := range templates {
		tmpl, err := template.New(name).Parse(pattern)
		if err != nil {
			log.Fatal(err)
		}
		
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, person)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Printf("%s: %s\n", name, buf.String())
	}
}
```

### Generación de Reportes

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"
)

func main() {
	type Report struct {
		Title   string
		Date    string
		Items   []string
		Total   float64
	}
	
	report := Report{
		Title: "Reporte de Ventas",
		Date:  time.Now().Format("2006-01-02"),
		Items: []string{"Producto A", "Producto B", "Producto C"},
		Total: 1500.50,
	}
	
	// Template con ciclos y formato
	tmplStr := `
=====================================
{{.Title}}
Fecha: {{.Date}}
=====================================

Artículos:
{{range $i, $item := .Items}}
{{add $i 1}}. {{$item}}
{{end}}

Total: ${{printf "%.2f" .Total}}
====================================`
	
	// Funciones personalizadas
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	
	tmpl, err := template.New("report").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, report)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(buf.String())
}
```

### Sub-templates y Composición

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

func main() {
	// Definir múltiples templates relacionados
	tmplStr := `
{{define "header"}}
===== {{.Title}} =====
{{end}}

{{define "footer"}}
Generado: {{.Date}}
{{end}}

{{define "main"}}
{{template "header" .}}

Contenido: {{.Content}}

{{template "footer" .}}
{{end}}
`
	
	tmpl, err := template.New("base").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	
	data := map[string]string{
		"Title":   "Mi Documento",
		"Content": "Este es el contenido principal",
		"Date":    "2024-01-15",
	}
	
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "main", data)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(buf.String())
}
```

---

## 28.10 HTML Templates

### HTML Escaping Automático

El package `html/template` proporciona escaping automático para prevenir XSS:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"html/template"
)

func main() {
	tmplStr := `
<h1>{{.Title}}</h1>
<p>{{.Content}}</p>
<script>{{.Script}}</script>
`
	
	tmpl, err := template.New("page").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	
	data := map[string]string{
		"Title":   "My Page",
		"Content": "Hello <world>",
		"Script":  "alert('xss');",
	}
	
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal(err)
	}
	
	// Output:
	// <h1>My Page</h1>
	// <p>Hello &lt;world&gt;</p>
	// <script>alert(&#39;xss&#39;);</script>
	
	// Nota: html/template escapa automáticamente el contenido
}
```

### Diferencias: text/template vs html/template

```
┌──────────────────┬─────────────────────┬──────────────────────┐
│ Aspecto          │ text/template       │ html/template        │
├──────────────────┼─────────────────────┼──────────────────────┤
│ Escaping         │ No                  │ Automático (XSS)     │
│ Contexto         │ Genérico            │ HTML/JS/CSS/URI      │
│ Seguridad        │ Manual              │ Automática           │
│ Performance      │ Ligeramente más rápido │ Ligeramente más lento│
│ Uso típico       │ Reportes, configs   │ Web pages, emails    │
└──────────────────┴─────────────────────┴──────────────────────┘
```

### Generación de HTML Segura

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"html/template"
)

func main() {
	// Template HTML
	tmplStr := `
<!DOCTYPE html>
<html>
<head>
	<title>{{.PageTitle}}</title>
</head>
<body>
	<h1>{{.Heading}}</h1>
	
	<ul>
	{{range .Items}}
		<li>{{.Name}} - {{.Price}}</li>
	{{end}}
	</ul>
	
	<form action="{{.FormAction}}" method="{{.Method}}">
		<input type="hidden" name="csrf" value="{{.CSRF}}">
		{{range .Fields}}
			<label>{{.Label}}</label>
			<input type="{{.Type}}" name="{{.Name}}" value="{{.Value}}">
		{{end}}
	</form>
</body>
</html>
`
	
	type Item struct {
		Name  string
		Price string
	}
	
	type Field struct {
		Label string
		Type  string
		Name  string
		Value string
	}
	
	tmpl, err := template.New("page").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	
	data := map[string]interface{}{
		"PageTitle":  "Product Store",
		"Heading":    "Our <Featured> Products",
		"Items": []Item{
			{"Laptop", "$999"},
			{"Phone", "$599"},
			{"Tablet", "$399"},
		},
		"FormAction": "/search",
		"Method":     "GET",
		"CSRF":       "token123abc",
		"Fields": []Field{
			{"Search Query", "text", "q", "search <here>"},
			{"Category", "select", "cat", "electronics"},
		},
	}
	
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(buf.String())
	
	// Nota: Los caracteres < y > en valores se escapan automáticamente
	// search <here> → search &lt;here&gt;
}
```

---

## 28.11 Buenas Prácticas y Patrones

### UTF-8 Handling

Go maneja UTF-8 nativamente, pero es importante entender la diferencia entre bytes y runes:

```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	text := "Hello, 世界!"  // Contiene caracteres multibyte
	
	// Diferencia: len() cuenta bytes, no caracteres
	fmt.Printf("Bytes: %d\n", len(text))            // 13
	fmt.Printf("Caracteres: %d\n", utf8.RuneCountInString(text)) // 9
	
	// Iterar bytes: índices diretos
	for i := 0; i < len(text); i++ {
		fmt.Printf("Byte %d: %x\n", i, text[i])
	}
	
	// Iterar runes (caracteres): usar for...range
	for i, r := range text {
		fmt.Printf("Rune %d: %c (U+%04X)\n", i, r, r)
	}
	
	// Acceso seguro a rune
	runes := []rune(text)
	fmt.Printf("Primer carácter: %c\n", runes[0])
	fmt.Printf("Último carácter: %c\n", runes[len(runes)-1])
	
	// Subcadenas de runes
	sub := string(runes[0:5])                       // "Hello"
	fmt.Println(sub)
}
```

### Performance: Benchmarking

```go
package main

import (
	"strings"
	"testing"
)

// Comparar rendimiento de diferentes enfoques

// INCORRECTO: concatenación con +
func concatenate_bad(parts []string) string {
	var result string
	for _, p := range parts {
		result += p                                 // ⚠️ O(n²)
	}
	return result
}

// CORRECTO: strings.Builder
func concatenate_good(parts []string) string {
	var builder strings.Builder
	for _, p := range parts {
		builder.WriteString(p)                      // ✓ O(n)
	}
	return builder.String()
}

// CORRECTO: strings.Join (aún más simple)
func concatenate_best(parts []string) string {
	return strings.Join(parts, "")                  // ✓ O(n), idiomatic
}

// Benchmarks
func BenchmarkConcatenate_Bad(b *testing.B) {
	parts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		parts[i] = "x"
	}
	for i := 0; i < b.N; i++ {
		concatenate_bad(parts)
	}
}

func BenchmarkConcatenate_Good(b *testing.B) {
	parts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		parts[i] = "x"
	}
	for i := 0; i < b.N; i++ {
		concatenate_good(parts)
	}
}

func BenchmarkConcatenate_Best(b *testing.B) {
	parts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		parts[i] = "x"
	}
	for i := 0; i < b.N; i++ {
		concatenate_best(parts)
	}
}

// Resultados esperados:
// BenchmarkConcatenate_Bad-8      16 ns/op      // Muy lento
// BenchmarkConcatenate_Good-8      8 ns/op      // Rápido
// BenchmarkConcatenate_Best-8      7 ns/op      // Óptimo
```

### Testing de Strings

```go
package main

import (
	"regexp"
	"strings"
	"testing"
)

// Función a testear
func normalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return "", ErrInvalidEmail
	}
	
	return email, nil
}

type testCase struct {
	input    string
	expected string
	hasErr   bool
}

func TestNormalizeEmail(t *testing.T) {
	tests := []testCase{
		{"john@example.com", "john@example.com", false},
		{"JOHN@EXAMPLE.COM", "john@example.com", false},
		{"  john@example.com  ", "john@example.com", false},
		{"invalid", "", true},
		{"missing@domain", "", true},
		{"", "", true},
	}
	
	for _, test := range tests {
		result, err := normalizeEmail(test.input)
		
		if test.hasErr && err == nil {
			t.Errorf("Se esperaba error para %q", test.input)
		}
		
		if !test.hasErr && err != nil {
			t.Errorf("No se esperaba error para %q: %v", test.input, err)
		}
		
		if result != test.expected {
			t.Errorf("Para %q: se esperaba %q, se obtuvo %q", 
				test.input, test.expected, result)
		}
	}
}
```

### Antipatrones Comunes

```go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	// ❌ ANTIPATRÓN 1: Compilar regex cada vez
	func validateEmail_bad(email string) bool {
		// Se compila cada invocación - ineficiente
		matched, _ := regexp.MatchString(`^[\w.-]+@[\w.-]+\.\w+$`, email)
		return matched
	}
	
	// ✅ PATRÓN CORRECTO
	var emailRegex = regexp.MustCompile(`^[\w.-]+@[\w.-]+\.\w+$`)
	func validateEmail_good(email string) bool {
		return emailRegex.MatchString(email)
	}
	
	// ❌ ANTIPATRÓN 2: Concatenación repetida
	func buildList_bad(items []string) string {
		result := ""
		for _, item := range items {
			result += "- " + item + "\n"                // O(n²)
		}
		return result
	}
	
	// ✅ PATRÓN CORRECTO
	func buildList_good(items []string) string {
		var builder strings.Builder
		for _, item := range items {
			builder.WriteString("- ")
			builder.WriteString(item)
			builder.WriteString("\n")
		}
		return builder.String()
	}
	
	// ❌ ANTIPATRÓN 3: HTML injection
	func render_bad(userInput string) string {
		// No escapa el contenido - ¡XSS!
		return "<h1>" + userInput + "</h1>"
	}
	
	// ✅ PATRÓN CORRECTO: usar html/template
	import "html/template"
	func render_good(userInput string) string {
		t := template.Must(template.New("").Parse("<h1>{{.}}</h1>"))
		var buf strings.Builder
		t.Execute(&buf, userInput)
		return buf.String()
	}
	
	// ❌ ANTIPATRÓN 4: Comparación case-sensitive innecesaria
	func isSameCommand_bad(cmd1, cmd2 string) bool {
		return cmd1 == cmd2                           // Sensible a caso
	}
	
	// ✅ PATRÓN CORRECTO
	func isSameCommand_good(cmd1, cmd2 string) bool {
		return strings.EqualFold(cmd1, cmd2)         // Insensible a caso
	}
}
```

### Testing Table-Driven

```go
package main

import (
	"testing"
	"strings"
)

func TestStringOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
		expected bool
	}{
		{
			name:     "substring exists",
			input:    "hello world",
			contains: "world",
			expected: true,
		},
		{
			name:     "substring doesn't exist",
			input:    "hello world",
			contains: "xyz",
			expected: false,
		},
		{
			name:     "case sensitive",
			input:    "Hello World",
			contains: "hello",
			expected: false,
		},
		{
			name:     "empty substring",
			input:    "hello",
			contains: "",
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.Contains(tt.input, tt.contains)
			if result != tt.expected {
				t.Errorf("Contains(%q, %q) = %v, want %v",
					tt.input, tt.contains, result, tt.expected)
			}
		})
	}
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Parser Simple - Configuración INI

Crear un parser simple para archivos `.ini`:

```go
package main

import (
	"fmt"
	"strings"
)

// Parsear un archivo INI simple
// [section]
// key=value
// key2=value2

func parseINI(content string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	var currentSection string
	
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if line == "" || strings.HasPrefix(line, ";") {
			continue  // Skip empty lines and comments
		}
		
		// Detect section
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			result[currentSection] = make(map[string]string)
			continue
		}
		
		// Parse key=value
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			
			if currentSection != "" {
				result[currentSection][key] = value
			}
		}
	}
	
	return result
}

func main() {
	config := `
[database]
host=localhost
port=5432
name=myapp

[server]
address=0.0.0.0
port=8080
debug=true

[logging]
level=info
format=json
`
	
	parsed := parseINI(config)
	
	for section, values := range parsed {
		fmt.Printf("[%s]\n", section)
		for key, value := range values {
			fmt.Printf("  %s = %s\n", key, value)
		}
	}
}
```

### Ejercicio 2: Validador con Regex

Crear validadores para email, URL, teléfono y contraseña:

```go
package main

import (
	"fmt"
	"regexp"
)

type Validator struct {
	email    *regexp.Regexp
	url      *regexp.Regexp
	phone    *regexp.Regexp
	password *regexp.Regexp
}

func NewValidator() *Validator {
	return &Validator{
		email: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		url:   regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		phone: regexp.MustCompile(`^\+?1?\d{9,15}$`),
		// Contraseña: al menos 8 chars, mayúscula, minúscula, número
		password: regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`),
	}
}

func (v *Validator) ValidateEmail(email string) bool {
	return v.email.MatchString(email)
}

func (v *Validator) ValidateURL(url string) bool {
	return v.url.MatchString(url)
}

func (v *Validator) ValidatePhone(phone string) bool {
	return v.phone.MatchString(phone)
}

func (v *Validator) ValidatePassword(password string) bool {
	return v.password.MatchString(password)
}

func main() {
	v := NewValidator()
	
	tests := []struct {
		name  string
		value string
		fn    func(string) bool
	}{
		{"email", "user@example.com", v.ValidateEmail},
		{"email", "invalid.email", v.ValidateEmail},
		{"url", "https://example.com", v.ValidateURL},
		{"url", "example.com", v.ValidateURL},
		{"phone", "+1234567890", v.ValidatePhone},
		{"phone", "123abc", v.ValidatePhone},
		{"password", "SecurePass123", v.ValidatePassword},
		{"password", "weak", v.ValidatePassword},
	}
	
	for _, test := range tests {
		result := test.fn(test.value)
		status := "✓"
		if !result {
			status = "✗"
		}
		fmt.Printf("%s %s: %s\n", status, test.name, test.value)
	}
}
```

### Ejercicio 3: Generador de Reportes con Templates

Crear un generador de reportes de ventas:

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"
)

type Sale struct {
	Date      string
	Product   string
	Quantity  int
	Price     float64
}

type SalesReport struct {
	Title      string
	GeneratedAt string
	Sales      []Sale
	Total      float64
	Average    float64
}

func (s SalesReport) CalculateTotal() float64 {
	total := 0.0
	for _, sale := range s.Sales {
		total += sale.Price * float64(sale.Quantity)
	}
	return total
}

func (s SalesReport) CalculateAverage() float64 {
	if len(s.Sales) == 0 {
		return 0
	}
	return s.CalculateTotal() / float64(len(s.Sales))
}

const reportTemplate = `
╔════════════════════════════════════════════════════════════╗
║ {{.Title}}
║ Generado: {{.GeneratedAt}}
╠════════════════════════════════════════════════════════════╣

VENTAS:
{{range $i, $sale := .Sales}}
{{add $i 1}}. {{$sale.Product}}
   Fecha: {{$sale.Date}}
   Cantidad: {{$sale.Quantity}}
   Precio unitario: ${{printf "%.2f" $sale.Price}}
   Subtotal: ${{printf "%.2f" (multiply $sale.Price (itof $sale.Quantity))}}
{{end}}

════════════════════════════════════════════════════════════
RESUMEN:
• Total de ventas: ${{printf "%.2f" .Total}}
• Venta promedio: ${{printf "%.2f" .Average}}
• Número de transacciones: {{len .Sales}}
════════════════════════════════════════════════════════════
`

func main() {
	sales := []Sale{
		{"2024-01-15", "Laptop", 2, 999.99},
		{"2024-01-16", "Mouse", 10, 29.99},
		{"2024-01-17", "Keyboard", 5, 79.99},
	}
	
	report := SalesReport{
		Title:       "REPORTE DE VENTAS MENSUAL",
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Sales:       sales,
	}
	report.Total = report.CalculateTotal()
	report.Average = report.CalculateAverage()
	
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"multiply": func(a float64, b float64) float64 { return a * b },
		"itof": func(i int) float64 { return float64(i) },
	}
	
	tmpl, err := template.New("report").Funcs(funcMap).Parse(reportTemplate)
	if err != nil {
		log.Fatal(err)
	}
	
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, report)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(buf.String())
}
```

### Ejercicio 4: String Builder - Generador de SQL

Construir dinámicamente queries SQL manteniendo eficiencia:

```go
package main

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	builder strings.Builder
	params  []interface{}
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.builder.WriteString("SELECT ")
	qb.builder.WriteString(strings.Join(columns, ", "))
	qb.builder.WriteString(" ")
	return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.builder.WriteString("FROM ")
	qb.builder.WriteString(table)
	qb.builder.WriteString(" ")
	return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.builder.WriteString("WHERE ")
	qb.builder.WriteString(condition)
	qb.builder.WriteString(" ")
	qb.params = append(qb.params, args...)
	return qb
}

func (qb *QueryBuilder) OrderBy(columns ...string) *QueryBuilder {
	qb.builder.WriteString("ORDER BY ")
	qb.builder.WriteString(strings.Join(columns, ", "))
	qb.builder.WriteString(" ")
	return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
	qb.builder.WriteString(fmt.Sprintf("LIMIT %d ", n))
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	return strings.TrimSpace(qb.builder.String()), qb.params
}

func main() {
	query := NewQueryBuilder().
		Select("id", "name", "email").
		From("users").
		Where("age > ?", 18).
		Where("status = ?", "active").
		OrderBy("created_at DESC").
		Limit(10)
	
	sql, params := query.Build()
	fmt.Println("SQL:", sql)
	fmt.Println("Params:", params)
	
	// Output:
	// SQL: SELECT id, name, email FROM users WHERE age > ? WHERE status = ? ORDER BY created_at DESC LIMIT 10
	// Params: [18 active]
}
```

### Ejercicio 5: Replacer para Templating

Crear un sistema de templating simple con múltiples sustituciones:

```go
package main

import (
	"fmt"
	"strings"
	"time"
)

type EmailTemplate struct {
	subject string
	body    string
	data    map[string]string
}

func NewEmailTemplate(subject, body string) *EmailTemplate {
	return &EmailTemplate{
		subject: subject,
		body:    body,
		data:    make(map[string]string),
	}
}

func (et *EmailTemplate) Set(key, value string) *EmailTemplate {
	et.data[key] = value
	return et
}

func (et *EmailTemplate) Render() (string, string) {
	// Construir Replacer dinámicamente
	var pairs []string
	
	for key, value := range et.data {
		pairs = append(pairs, "{{"+key+"}}", value)
	}
	
	// Agregar variables especiales
	pairs = append(pairs, "{{today}}", time.Now().Format("2006-01-02"))
	pairs = append(pairs, "{{year}}", fmt.Sprintf("%d", time.Now().Year()))
	
	replacer := strings.NewReplacer(pairs...)
	
	return replacer.Replace(et.subject), replacer.Replace(et.body)
}

func main() {
	// Crear template de email
	template := NewEmailTemplate(
		"Bienvenido {{name}}",
		`Hola {{name}},

Gracias por registrarte en {{app_name}}.
Tu email: {{email}}
Fecha de registro: {{today}}

Saludos,
El equipo de {{app_name}}

© {{year}} {{company}}`,
	)
	
	// Configurar datos
	subject, body := template.
		Set("name", "Juan").
		Set("email", "juan@example.com").
		Set("app_name", "MyApp").
		Set("company", "Tech Company").
		Render()
	
	fmt.Println("ASUNTO:")
	fmt.Println(subject)
	fmt.Println("\nCUERPO:")
	fmt.Println(body)
	
	// Reutilizar template con diferentes datos
	fmt.Println("\n--- OTRO EMAIL ---\n")
	
	subject2, body2 := NewEmailTemplate(
		"Bienvenido {{name}}",
		`Hola {{name}},

Gracias por registrarte en {{app_name}}.
Tu email: {{email}}
Fecha de registro: {{today}}

Saludos,
El equipo de {{app_name}}

© {{year}} {{company}}`,
	).
		Set("name", "María").
		Set("email", "maria@example.com").
		Set("app_name", "MyApp").
		Set("company", "Tech Company").
		Render()
	
	fmt.Println("ASUNTO:")
	fmt.Println(subject2)
	fmt.Println("\nCUERPO:")
	fmt.Println(body2)
}
```

---

## Resumen de Comparaciones

### Go vs Otros Lenguajes

```
┌─────────────────┬────────────────┬────────────────┬────────────────┐
│ Característica  │ Go             │ Python         │ JavaScript     │
├─────────────────┼────────────────┼────────────────┼────────────────┤
│ Regex engine    │ RE2 (Google)   │ PCRE           │ V8             │
│ Performance     │ Muy rápido     │ Moderno        │ Moderno        │
│ Backtracking    │ No (seguro)    │ Sí (riesgo)    │ Sí (riesgo)    │
│ Templates       │ 2 packages     │ Jinja2         │ Mustache/EJS   │
│ Escaping HTML   │ Automático     │ Manual         │ Manual         │
│ Unicode         │ UTF-8 nativo   │ Nativo         │ Nativo         │
│ Inmutabilidad   │ Strings sí     │ Strings sí     │ Strings sí     │
└─────────────────┴────────────────┴────────────────┴────────────────┘
```

### Regex: RE2 vs PCRE

Go utiliza **RE2** (Google), que difiere de PCRE en algunas características:

```go
// Features en PCRE pero NO en RE2 (RE2 las evita por seguridad)
(?P<name>...)      // ✓ Named groups soportados
(?!...)            // ✗ Negative lookahead NO
(?<=...)           // ✗ Lookbehind NO
.*+                // ✗ Possessive quantifiers NO
\p{Greek}          // ✗ Unicode properties NO (parcial)

// Beneficios de RE2:
// • O(n) time complexity garantizado
// • No catastrophic backtracking
// • Seguro en entornos con entrada del usuario
// • Compilación determinística
```

---

## Consejo Final: Cuando Usar Qué

```
┌────────────────────────────────────────────────────────────────┐
│ USO                        │ HERRAMIENTA RECOMENDADA           │
├────────────────────────────────────────────────────────────────┤
│ Buscar/reemplazar simple   │ strings.Contains, strings.Replace  │
│ Parsing estructura conocida│ strings.Split, strings.Fields      │
│ Múltiples concatenaciones  │ strings.Builder o strings.Join     │
│ Búsqueda case-insensitive  │ strings.ToLower + Contains        │
│ Validación de formato      │ regexp (email, phone, etc)        │
│ Extraer datos              │ regexp.FindAllSubmatch            │
│ Generar texto estático     │ text/template                     │
│ Generar HTML dinámico      │ html/template (por seguridad)     │
│ Configuración              │ Parsers personalizados + strings   │
│ Reporte formateado         │ text/template + strings.Builder   │
└────────────────────────────────────────────────────────────────┘
```

---

## Recursos Adicionales

- **Package strings**: https://pkg.go.dev/strings
- **Package regexp**: https://pkg.go.dev/regexp  
- **RE2 Syntax**: https://github.com/google/re2/wiki/Syntax
- **text/template**: https://pkg.go.dev/text/template
- **html/template**: https://pkg.go.dev/html/template
- **Unicode en Go**: https://blog.golang.org/strings

---

**© 2024 - Guía Exhaustiva de Go - Capítulo 28**
*Última actualización: Enero 2024*
