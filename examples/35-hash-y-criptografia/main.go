// Ejemplo del Capítulo 35: hashing con SHA-256, HMAC para verificar
// integridad/autenticidad, y generación de tokens criptográficamente
// seguros con crypto/rand.
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var claveSecreta = []byte("clave-secreta-de-ejemplo-no-usar-en-producción")

func firmar(mensaje string) string {
	mac := hmac.New(sha256.New, claveSecreta)
	mac.Write([]byte(mensaje))
	return hex.EncodeToString(mac.Sum(nil))
}

func verificar(mensaje, firma string) bool {
	esperada := firmar(mensaje)
	// hmac.Equal compara en tiempo constante: evita timing attacks.
	return hmac.Equal([]byte(esperada), []byte(firma))
}

func main() {
	// SHA-256 simple: útil para checksums, no para contraseñas.
	suma := sha256.Sum256([]byte("contenido a verificar"))
	fmt.Println("SHA-256:", hex.EncodeToString(suma[:]))

	mensaje := "transferir $100 a la cuenta 12345"
	firma := firmar(mensaje)
	fmt.Println("mensaje:", mensaje)
	fmt.Println("firma HMAC:", firma)

	fmt.Println("¿firma válida?", verificar(mensaje, firma))
	fmt.Println("¿firma válida tras alterar el mensaje?",
		verificar("transferir $100 a la cuenta 99999", firma))

	// Generación de tokens/salts aleatorios seguros para criptografía.
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		fmt.Println("error generando token:", err)
		return
	}
	fmt.Println("token aleatorio seguro:", hex.EncodeToString(token))
}
