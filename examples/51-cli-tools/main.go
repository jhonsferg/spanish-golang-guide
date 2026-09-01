// Ejemplo del Capítulo 51: un CLI con subcomandos usando solo el
// paquete estándar `flag` - el patrón detrás de herramientas como
// `go` o `git` antes de necesitar una librería como cobra.
package main

import (
	"flag"
	"fmt"
	"os"
)

type tarea struct {
	ID        int
	Titulo    string
	Completa  bool
}

var tareas []tarea
var siguienteID = 1

func cmdAgregar(args []string) error {
	fs := flag.NewFlagSet("agregar", flag.ContinueOnError)
	titulo := fs.String("titulo", "", "título de la tarea (requerido)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *titulo == "" {
		return fmt.Errorf("agregar: -titulo es requerido")
	}
	tareas = append(tareas, tarea{ID: siguienteID, Titulo: *titulo})
	fmt.Printf("tarea #%d agregada: %s\n", siguienteID, *titulo)
	siguienteID++
	return nil
}

func cmdListar(args []string) error {
	fs := flag.NewFlagSet("listar", flag.ContinueOnError)
	soloPendientes := fs.Bool("pendientes", false, "mostrar solo tareas incompletas")
	if err := fs.Parse(args); err != nil {
		return err
	}
	for _, t := range tareas {
		if *soloPendientes && t.Completa {
			continue
		}
		estado := " "
		if t.Completa {
			estado = "x"
		}
		fmt.Printf("[%s] #%d %s\n", estado, t.ID, t.Titulo)
	}
	return nil
}

func despachar(comando string, args []string) error {
	switch comando {
	case "agregar":
		return cmdAgregar(args)
	case "listar":
		return cmdListar(args)
	default:
		return fmt.Errorf("comando desconocido: %s (usa: agregar | listar)", comando)
	}
}

func main() {
	// Simulamos una secuencia de invocaciones de línea de comandos, como si
	// el binario se hubiera llamado varias veces con distintos os.Args.
	invocaciones := [][]string{
		{"agregar", "-titulo=Escribir la guía"},
		{"agregar", "-titulo=Revisar ejemplos"},
		{"listar"},
	}

	for _, args := range invocaciones {
		fmt.Printf("$ cli %v\n", args)
		if err := despachar(args[0], args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}
