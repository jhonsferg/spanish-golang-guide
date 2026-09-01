// Ejemplo del Capítulo 29: (de)serialización JSON con struct tags,
// campos opcionales y un tipo con Marshal/Unmarshal personalizados.
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Usuario struct {
	Nombre      string    `json:"nombre"`
	Email       string    `json:"email"`
	Edad        int       `json:"edad,omitempty"`
	Interno     string    `json:"-"` // nunca se serializa
	CreadoEn    time.Time `json:"creado_en"`
}

func main() {
	u := Usuario{
		Nombre:   "Grace Hopper",
		Email:    "grace@example.com",
		Interno:  "no debería aparecer en el JSON",
		CreadoEn: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	datos, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		fmt.Println("error al serializar:", err)
		return
	}
	fmt.Println(string(datos))

	// Deserializar JSON externo a una struct Go.
	entrada := `{"nombre":"Ada Lovelace","email":"ada@example.com","edad":36}`
	var u2 Usuario
	if err := json.Unmarshal([]byte(entrada), &u2); err != nil {
		fmt.Println("error al deserializar:", err)
		return
	}
	fmt.Printf("deserializado: %+v\n", u2)

	// Cuando no conoces la forma exacta del JSON, usa map[string]any.
	var generico map[string]any
	json.Unmarshal([]byte(`{"a":1,"b":[1,2,3],"c":{"d":true}}`), &generico)
	fmt.Println("JSON genérico:", generico)
}
