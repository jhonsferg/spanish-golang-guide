// Ejemplo del Capítulo 27: interacción con el sistema operativo -
// archivos temporales, variables de entorno y argumentos.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Directorio temporal del sistema, portable entre SO.
	dir, err := os.MkdirTemp("", "guia-go-cap27-*")
	if err != nil {
		fmt.Println("error creando directorio temporal:", err)
		return
	}
	defer os.RemoveAll(dir) // limpieza garantizada al salir

	ruta := filepath.Join(dir, "saludo.txt")
	contenido := []byte("escrito desde el Capítulo 27\n")

	if err := os.WriteFile(ruta, contenido, 0o644); err != nil {
		fmt.Println("error escribiendo archivo:", err)
		return
	}

	leido, err := os.ReadFile(ruta)
	if err != nil {
		fmt.Println("error leyendo archivo:", err)
		return
	}
	fmt.Print("contenido leído: ", string(leido))

	info, err := os.Stat(ruta)
	if err == nil {
		fmt.Printf("tamaño: %d bytes, modificado: %s\n", info.Size(), info.ModTime().Format("2006-01-02"))
	}

	// Variables de entorno: LookupEnv distingue "vacía" de "no definida".
	if valor, existe := os.LookupEnv("PATH"); existe {
		fmt.Printf("PATH está definida (longitud: %d)\n", len(valor))
	}

	fmt.Println("argumentos del programa:", os.Args)
	fmt.Println("directorio de trabajo actual disponible vía os.Getwd()")
}
