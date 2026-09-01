// Ejemplo del Capítulo 48: capas de configuración con precedencia
// explícita - defaults < archivo < variables de entorno < flags - el
// orden que usan la mayoría de las apps de producción en Go.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Puerto     int    `json:"puerto"`
	Entorno    string `json:"entorno"`
	NivelLog   string `json:"nivel_log"`
}

func defaults() Config {
	return Config{Puerto: 8080, Entorno: "desarrollo", NivelLog: "info"}
}

func cargarDesdeArchivo(cfg *Config, ruta string) error {
	datos, err := os.ReadFile(ruta)
	if os.IsNotExist(err) {
		return nil // el archivo es opcional
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(datos, cfg)
}

func aplicarVariablesDeEntorno(cfg *Config) {
	if v := os.Getenv("APP_ENTORNO"); v != "" {
		cfg.Entorno = v
	}
	if v := os.Getenv("APP_NIVEL_LOG"); v != "" {
		cfg.NivelLog = v
	}
}

func aplicarFlags(cfg *Config, args []string) error {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	puerto := fs.Int("puerto", cfg.Puerto, "puerto HTTP")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg.Puerto = *puerto
	return nil
}

func main() {
	cfg := defaults()
	fmt.Printf("1. defaults:          %+v\n", cfg)

	archivoTemp, _ := os.CreateTemp("", "config-*.json")
	defer os.Remove(archivoTemp.Name())
	archivoTemp.WriteString(`{"puerto": 9090, "entorno": "staging"}`)
	archivoTemp.Close()

	if err := cargarDesdeArchivo(&cfg, archivoTemp.Name()); err != nil {
		fmt.Println("error leyendo archivo:", err)
		return
	}
	fmt.Printf("2. tras archivo:      %+v\n", cfg)

	os.Setenv("APP_NIVEL_LOG", "debug")
	aplicarVariablesDeEntorno(&cfg)
	fmt.Printf("3. tras env vars:     %+v\n", cfg)

	if err := aplicarFlags(&cfg, []string{"-puerto=3000"}); err != nil {
		fmt.Println("error parseando flags:", err)
		return
	}
	fmt.Printf("4. tras flags (gana): %+v\n", cfg)
}
