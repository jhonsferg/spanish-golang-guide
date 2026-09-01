# Capítulo 35: Hash y criptografía

## Introducción

La criptografía es el corazón de la seguridad moderna. Desde proteger contraseñas hasta garantizar la integridad de datos en tránsito, Go proporciona herramientas criptográficas robustas en sus paquetes `crypto/*`. Este capítulo explora cómo implementar seguridad real en tus aplicaciones Go.

**Nota importante:** La criptografía es un campo especializado. Si bien Go proporciona implementaciones seguras y probadas, siempre consulta con expertos en seguridad para aplicaciones críticas. No reinventes el rueda criptográfica.

---

## 35.1 Fundamentos de Criptografía

### 35.1.1 Conceptos Básicos

**Criptografía vs Criptología:**

- **Criptografía:** Técnica de escribir en secreto
- **Criptoanálisis:** Técnica de descifrar mensajes sin la clave
- **Criptología:** Combinación de ambas

**Tres pilares de la seguridad:**

```
[Confidencialidad] - Nadie más puede leer
[Integridad]       - No fue modificado
[Autenticidad]     - Sé quién lo envió
```

### 35.1.2 Hashing vs Encriptación

| Característica | Hashing | Encriptación |
|---|---|---|
| **Reversible** | No (unidireccional) | Sí (bidireccional) |
| **Entrada/Salida** | Cualquier tamaño → Tamaño fijo | Cualquier tamaño → Cualquier tamaño |
| **Uso principal** | Integridad, fingerprints | Confidencialidad |
| **Velocidad** | Rápido | Variable (depende del algoritmo) |
| **Clave** | No requiere | Sí requiere |
| **Ejemplo** | SHA256 | AES |

**Visualización:**

```
HASHING:
"Mi contraseña" → [SHA256] → "a7f3d4b2c1e9f6a..."
Irreversible. Mismo input = siempre mismo hash.

ENCRIPTACIÓN:
"Mi contraseña" + [Clave] → [AES] → "Φ≈∂ƒ˚∆˙©ƒ∆˚"
Reversible: "Φ≈∂ƒ˚∆˙©ƒ∆˚" + [Clave] → "Mi contraseña"
```

### 35.1.3 Casos de Uso Fundamentales

**Cuándo usar Hashing:**

- Almacenar contraseñas
- Verificar integridad de archivos
- Crear fingerprints de datos
- Implementar HMAC
- Estructuras de datos (hash tables, blockchain)

**Cuándo usar Encriptación:**

- Proteger datos en tránsito (HTTPS, TLS)
- Almacenamiento de datos sensibles
- Comunicación privada
- Cumplimiento normativo (GDPR, HIPAA)

**Cuándo usar Signing:**

- Autenticar mensajes
- Non-repudiation (prueba de envío)
- Verificar origen de datos
- Certificados digitales

### 35.1.4 Criptografía Simétrica vs Asimétrica

**Simétrica (una clave compartida):**

```
Alice                                    Bob
  |                                       |
  | Clave: "SecretoCompartido"           |
  |          ↓                            |
  | Plaintext → [AES] → Ciphertext -----→| Ciphertext → [AES] → Plaintext
  |                                       | Clave: "SecretoCompartido"
```

Pros: Rápido, simple
Cons: Distribución de clave, escalabilidad

**Asimétrica (par de claves):**

```
Alice                              Bob
  | Clave privada de Bob (pública) |
  |     ↓                          |
  | Plaintext → [RSA] → Ciphertext ----→| Ciphertext → [RSA privada] → Plaintext
  |                               |
  |                               | Clave privada: Bob
```

Pros: Distribución de clave segura, escalable
Cons: Más lento, computacionalmente intensivo

---

## 35.2 Funciones Hash

### 35.2.1 Conceptos de Hash

**Propiedades de un hash criptográfico seguro:**

1. **Determinismo:** Mismo input → siempre mismo output
2. **Eficiencia:** Rápido de calcular
3. **Avalanche effect:** Cambio mínimo en input → cambio máximo en output
4. **Unidireccionalidad:** Imposible obtener input del output
5. **Resistencia colisión:** Imposible encontrar dos inputs iguales
6. **Tamaño fijo:** Output siempre mismo tamaño

### 35.2.2 Familia SHA (Secure Hash Algorithm)

**Comparación de algoritmos:**

| Algoritmo | Tamaño | Velocidad | Seguridad | Recomendación |
|---|---|---|---|---|
| **MD5** | 128 bits | Muy rápido | ❌ Roto | NUNCA usar |
| **SHA1** | 160 bits | Rápido | ⚠️ Débil | Solo legacy |
| **SHA256** | 256 bits | Rápido | ✅ Seguro | **RECOMENDADO** |
| **SHA512** | 512 bits | Rápido | ✅ Muy seguro | Alto nivel seguridad |
| **SHA3-256** | 256 bits | Rápido | ✅ Seguro | Alternativa moderna |
| **SHA3-512** | 512 bits | Rápido | ✅ Muy seguro | Alternativa moderna |

### 35.2.3 Implementación Básica: SHA256

```go
package main

import (
    "crypto/sha256"
    "fmt"
    "io"
    "os"
)

// HashString calcula el SHA256 de una cadena
func HashString(s string) string {
    h := sha256.New()
    h.Write([]byte(s))
    return fmt.Sprintf("%x", h.Sum(nil))
}

// HashFile calcula el SHA256 de un archivo
func HashFile(filepath string) (string, error) {
    f, err := os.Open(filepath)
    if err != nil {
        return "", err
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// VerifyHash verifica que un hash coincida
func VerifyHash(data string, expectedHash string) bool {
    actualHash := HashString(data)
    // Usar comparación segura contra timing attacks
    return constantTimeCompare(actualHash, expectedHash)
}

// constantTimeCompare compara dos strings en tiempo constante
func constantTimeCompare(a, b string) bool {
    if len(a) != len(b) {
        return false
    }
    result := 0
    for i := 0; i < len(a); i++ {
        result |= int(a[i] ^ b[i])
    }
    return result == 0
}

func main() {
    // Ejemplo básico
    data := "Hola, mundo"
    hash := HashString(data)
    fmt.Printf("SHA256('%s') = %s\n", data, hash)

    // Verificar hash
    if VerifyHash(data, hash) {
        fmt.Println("✓ Hash verificado correctamente")
    }

    // Hash de archivo
    if file_hash, err := HashFile("ejemplo.txt"); err == nil {
        fmt.Printf("Hash del archivo: %s\n", file_hash)
    }
}
```

### 35.2.4 Otros Algoritmos Hash

```go
package main

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "fmt"
    "golang.org/x/crypto/sha3"
)

func DemostrarHashAlgoritmos(data string) {
    fmt.Println("=== Comparación de Algoritmos Hash ===\n")

    // MD5 (NO USAR en seguridad)
    md5_hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
    fmt.Printf("MD5:       %s (❌ INSEGURO)\n", md5_hash)

    // SHA1 (Debilitado)
    sha1_hash := fmt.Sprintf("%x", sha1.Sum([]byte(data)))
    fmt.Printf("SHA1:      %s (⚠️ DÉBIL)\n", sha1_hash)

    // SHA256 (Estándar)
    sha256_hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
    fmt.Printf("SHA256:    %s (✓ RECOMENDADO)\n", sha256_hash)

    // SHA512
    sha512_hash := fmt.Sprintf("%x", sha512.Sum512([]byte(data)))
    fmt.Printf("SHA512:    %s (✓ MÁXIMA SEGURIDAD)\n", sha512_hash[:64])
    fmt.Println("           (primeros 64 caracteres)")

    // SHA3-256
    sha3_256 := fmt.Sprintf("%x", sha3.Sum256([]byte(data)))
    fmt.Printf("SHA3-256:  %s (✓ MODERNO)\n", sha3_256)

    // SHA3-512
    sha3_512 := fmt.Sprintf("%x", sha3.Sum512([]byte(data)))
    fmt.Printf("SHA3-512:  %s (✓ MODERNO)\n", sha3_512[:64])
}

func main() {
    DemostrarHashAlgoritmos("Mensaje importante")
}
```

### 35.2.5 Propiedades Críticas: Avalanche Effect

```go
package main

import (
    "crypto/sha256"
    "fmt"
)

func mostrarAvalanche(data1, data2 string) {
    h1 := sha256.Sum256([]byte(data1))
    h2 := sha256.Sum256([]byte(data2))

    fmt.Printf("Input 1: %s\n", data1)
    fmt.Printf("Hash 1:  %x\n\n", h1)

    fmt.Printf("Input 2: %s\n", data2)
    fmt.Printf("Hash 2:  %x\n\n", h2)

    // Contar diferencias
    diferences := 0
    for i := 0; i < 32; i++ {
        if h1[i] != h2[i] {
            diferences++
        }
    }

    fmt.Printf("Bytes diferentes: %d/32\n", diferences)
    fmt.Printf("Cambios en bits (aproximado): ~%d/256\n", diferences*8/2)
}

func main() {
    // Pequeño cambio → Completamente diferente
    mostrarAvalanche("password", "passwor2")
}
```

---

## 35.3 HMAC - Códigos de Autenticación de Mensaje

### 35.3.1 ¿Qué es HMAC?

HMAC (Hash-based Message Authentication Code) combina:

- Un hash criptográfico
- Una clave secreta compartida

**Propósito:** Verificar integridad Y autenticidad

```
Sender: mensaje + clave secreta → [HMAC-SHA256] → código
        envía mensaje + código

Receiver: recibe mensaje + código
          mensaje + clave secreta → [HMAC-SHA256] → código calculado
          ¿código calculado == código recibido? → autenticidad verificada
```

### 35.3.2 Implementación HMAC

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "io"
)

// GenerarHMAC genera un código de autenticación HMAC-SHA256
func GenerarHMAC(mensaje string, clave string) string {
    h := hmac.New(sha256.New, []byte(clave))
    h.Write([]byte(mensaje))
    return fmt.Sprintf("%x", h.Sum(nil))
}

// VerificarHMAC verifica un código HMAC
func VerificarHMAC(mensaje string, codigo string, clave string) bool {
    codigoCalculado := GenerarHMAC(mensaje, clave)
    return hmac.Equal([]byte(codigoCalculado), []byte(codigo))
}

// HMAC para archivos
func GenerarHMACArchivo(filepath string, clave string) (string, error) {
    f, err := os.Open(filepath)
    if err != nil {
        return "", err
    }
    defer f.Close()

    h := hmac.New(sha256.New, []byte(clave))
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func main() {
    mensaje := "Transacción bancaria: $1000"
    claveSecreta := "MiClaveSecretaMuySegura"

    // Generar HMAC
    codigo := GenerarHMAC(mensaje, claveSecreta)
    fmt.Printf("Mensaje: %s\n", mensaje)
    fmt.Printf("HMAC:    %s\n\n", codigo)

    // Verificar HMAC válido
    if VerificarHMAC(mensaje, codigo, claveSecreta) {
        fmt.Println("✓ HMAC válido - Mensaje auténtico")
    }

    // Intentar modificar mensaje
    mensajeModificado := "Transacción bancaria: $10000"
    if !VerificarHMAC(mensajeModificado, codigo, claveSecreta) {
        fmt.Println("✗ HMAC inválido - Mensaje fue modificado")
    }
}
```

### 35.3.3 Casos de Uso HMAC

**Webhooks con validación:**

```go
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
    // Leer el HMAC enviado en header
    receivedHMAC := r.Header.Get("X-Hub-Signature-256")

    // Leer el cuerpo
    bodyBytes, _ := ioutil.ReadAll(r.Body)

    // Generar HMAC con clave secreta
    secretKey := os.Getenv("WEBHOOK_SECRET")
    expectedHMAC := GenerarHMAC(string(bodyBytes), secretKey)

    // Validar
    if !hmac.Equal([]byte(expectedHMAC), []byte(receivedHMAC)) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Procesar webhook...
}
```

**API authentication:**

```go
// Cliente
func HacerRequestAutenticado(url string, datos string, apiSecret string) {
    codigo := GenerarHMAC(datos, apiSecret)

    req, _ := http.NewRequest("POST", url, strings.NewReader(datos))
    req.Header.Set("X-API-Signature", codigo)

    // Enviar request...
}
```

---

## 35.4 Password Hashing - Protección de Contraseñas

### 35.4.1 ¿Por qué NO usar SHA directamente?

```go
// ❌ NUNCA hacer esto
contraseña := "MiContraseña123"
hash := sha256.Sum256([]byte(contraseña))  // ¡INSEGURO!

// Problemas:
// 1. Rainbow tables: precalcular millones de hashes
// 2. Sin salt: mismo usuario con otra BD = mismo hash
// 3. Muy rápido: fuerza bruta en segundos
// 4. Sin estirado: cualquiera puede intentar verificar
```

### 35.4.2 Bcrypt - El Estándar

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

// HashPassword genera un hash seguro de contraseña
func HashPassword(password string) (string, error) {
    // Cost debe estar entre 4 y 31
    // 10 = ~100ms en máquina típica
    // 12 = ~250ms (recomendado para login)
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password),
        12, // cost
    )
    return string(hashedPassword), err
}

// VerifyPassword verifica una contraseña contra su hash
func VerifyPassword(hashedPassword string, password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(hashedPassword),
        []byte(password),
    )
    return err == nil
}

func main() {
    // Registrar usuario
    contraseña := "MiContraseñaSegura123"
    hash, _ := HashPassword(contraseña)
    fmt.Printf("Hash original:    %s\n\n", hash)

    // Almacenar en BD: hash

    // Login - verificar
    if VerifyPassword(hash, contraseña) {
        fmt.Println("✓ Contraseña correcta - Login exitoso")
    }

    // Intentar contraseña incorrecta
    if !VerifyPassword(hash, "ContraseñaIncorrecta") {
        fmt.Println("✗ Contraseña incorrecta - Login fallido")
    }

    // Importante: cada hash es diferente (contiene salt aleatorio)
    hash2, _ := HashPassword(contraseña)
    fmt.Printf("\nHash diferente:   %s", hash2)
    fmt.Println("\n(Mismo input, diferente hash debido al salt)")
}
```

### 35.4.3 Argon2 - Más Moderno

```go
package main

import (
    "encoding/base64"
    "fmt"
    "golang.org/x/crypto/argon2"
)

// HashPasswordArgon2 usa Argon2 para hashing
func HashPasswordArgon2(password string) string {
    // Configuración segura
    salt := make([]byte, 16)
    rand.Read(salt)  // Salt aleatorio

    // Argon2id es mejor que Argon2i o Argon2d
    hash := argon2.IDKey(
        []byte(password),
        salt,
        3,        // time (iteraciones)
        64*1024,  // memory (64 MB)
        4,        // parallelism
        32,       // keyLen (256 bits)
    )

    // Devolver salt + hash codificado
    return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=4$%s$%s",
        base64.RawStdEncoding.EncodeToString(salt),
        base64.RawStdEncoding.EncodeToString(hash),
    )
}
```

### 35.4.4 Comparación: Bcrypt vs Argon2

| Característica | Bcrypt | Argon2 |
|---|---|---|
| **Velocidad de hash** | Lento (~100-250ms) | Lento (~100-500ms) |
| **Resistencia GPU** | Buena | Excelente |
| **Resistencia timing attacks** | Sí | Sí |
| **Salt automático** | Sí | No (manual) |
| **Modernidad** | 2000s | 2015 (más moderno) |
| **Recomendación** | ✓ Seguro | ✓ Mejor opción |

### 35.4.5 Salt y Pepper

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

// GenerarSalt crea un salt criptográficamente seguro
func GenerarSalt(longitud int) (string, error) {
    salt := make([]byte, longitud)
    _, err := rand.Read(salt)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(salt), nil
}

// HashConSalt - Hash manual con salt
func HashConSalt(password string, salt string) string {
    h := sha256.New()
    h.Write([]byte(password + salt))
    return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
    pass := "MiContraseña"

    // Generar salt aleatorio
    salt, _ := GenerarSalt(16)
    hash := HashConSalt(pass, salt)

    fmt.Printf("Salt: %s\n", salt)
    fmt.Printf("Hash: %s\n", hash)

    // Para verificar: regenerar hash con mismo salt y comparar
    verificar := HashConSalt(pass, salt)
    fmt.Printf("Verificación: %s\n", verificar == hash)
}

// Pepper: constante en servidor (NO en BD)
const PEPPER = "mi-constante-secreta-en-env"
// Usar: HashConSalt(password + PEPPER, salt)
```

### 35.4.6 Contexto en Hashing

```go
package main

import (
    "golang.org/x/crypto/argon2"
)

// ContextoHash para diferente información
type ContextoHash struct {
    Username string
    Email    string
}

// HashConContexto vincula hash al usuario
func HashConContexto(password string, contexto ContextoHash) string {
    salt := []byte(contexto.Username + contexto.Email)
    hash := argon2.IDKey(
        []byte(password),
        salt[:16],
        3, 64*1024, 4, 32,
    )
    return fmt.Sprintf("%x", hash)
}

// Evita reutilizar hashes entre aplicaciones/usuarios
```

---

## 35.5 Encriptación Simétrica

### 35.5.1 AES - Advanced Encryption Standard

**AES:** Cifrado de bloques, 128/192/256 bits

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
)

// EncriptarAES encripta datos con AES
func EncriptarAES(plaintext string, clave string) (string, error) {
    // Clave debe ser 16, 24 o 32 bytes (128, 192, 256 bits)
    // Si es más corta, rellenarla; si es más larga, truncar
    if len(clave) > 32 {
        clave = clave[:32]
    }
    for len(clave) < 32 {
        clave += " "
    }

    bloque, err := aes.NewCipher([]byte(clave))
    if err != nil {
        return "", err
    }

    // Generar IV (Initialization Vector) aleatorio
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    // Crear stream con CBC
    stream := cipher.NewCBCEncrypter(bloque, iv)

    // Rellenar plaintext (PKCS7)
    plaintext_padded := pkcs7Pad(plaintext, aes.BlockSize)

    // Encriptar
    ciphertext := make([]byte, len(plaintext_padded))
    stream.CryptBlocks(ciphertext, []byte(plaintext_padded))

    // Prepend IV y codificar en base64
    resultado := append(iv, ciphertext...)
    return base64.StdEncoding.EncodeToString(resultado), nil
}

// DesencriptarAES desencripta datos AES
func DesencriptarAES(ciphertext64 string, clave string) (string, error) {
    // Rellenar clave
    if len(clave) > 32 {
        clave = clave[:32]
    }
    for len(clave) < 32 {
        clave += " "
    }

    // Decodificar base64
    ciphertext, err := base64.StdEncoding.DecodeString(ciphertext64)
    if err != nil {
        return "", err
    }

    // Extraer IV
    if len(ciphertext) < aes.BlockSize {
        return "", fmt.Errorf("ciphertext demasiado corto")
    }
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    // Crear cipher
    bloque, err := aes.NewCipher([]byte(clave))
    if err != nil {
        return "", err
    }

    // Desencriptar
    stream := cipher.NewCBCDecrypter(bloque, iv)
    plaintext := make([]byte, len(ciphertext))
    stream.CryptBlocks(plaintext, ciphertext)

    // Remover padding
    plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
    return string(plaintext), err
}

// PKCS7 padding
func pkcs7Pad(data string, blockSize int) string {
    padding := blockSize - (len(data) % blockSize)
    padChar := byte(padding)
    for i := 0; i < padding; i++ {
        data += string(padChar)
    }
    return data
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
    if len(data) == 0 {
        return nil, fmt.Errorf("datos vacíos")
    }
    padding := int(data[len(data)-1])
    if padding < 1 || padding > blockSize {
        return nil, fmt.Errorf("padding inválido")
    }
    return data[:len(data)-padding], nil
}

func main() {
    plaintext := "Mensaje secreto muy importante"
    clave := "MiClaveSecretaDe32BytesMuyBuena"

    // Encriptar
    encrypted, _ := EncriptarAES(plaintext, clave)
    fmt.Printf("Plaintext:  %s\n", plaintext)
    fmt.Printf("Encrypted:  %s\n\n", encrypted)

    // Desencriptar
    decrypted, _ := DesencriptarAES(encrypted, clave)
    fmt.Printf("Decrypted:  %s\n", decrypted)
    fmt.Printf("Coincide:   %v\n", plaintext == decrypted)
}
```

### 35.5.2 GCM - Galois/Counter Mode

**GCM:** Encriptación + Autenticación (AEAD)

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

// EncriptarAES_GCM encripta con autenticación
func EncriptarAES_GCM(plaintext string, clave string) (string, error) {
    bloque, err := aes.NewCipher([]byte(clave[:32]))
    if err != nil {
        return "", err
    }

    // Crear GCM
    gcm, err := cipher.NewGCM(bloque)
    if err != nil {
        return "", err
    }

    // Generar nonce aleatorio
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    // Encriptar (incluye tag de autenticación)
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DesencriptarAES_GCM desencripta y verifica autenticación
func DesencriptarAES_GCM(ciphertext64 string, clave string) (string, error) {
    bloque, err := aes.NewCipher([]byte(clave[:32]))
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(bloque)
    if err != nil {
        return "", err
    }

    ciphertext, _ := base64.StdEncoding.DecodeString(ciphertext64)
    nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    return string(plaintext), err
}

func main() {
    plaintext := "Datos confidenciales"
    clave := "MiClaveSecretaDe32BytesMuyBuena"

    // Encriptar con autenticación
    encrypted, _ := EncriptarAES_GCM(plaintext, clave)
    fmt.Println("✓ Encrypted with authentication tag")

    // Desencriptar (si fue modificado, Open() falla)
    decrypted, err := DesencriptarAES_GCM(encrypted, clave)
    if err == nil {
        fmt.Printf("✓ Verified and decrypted: %s\n", decrypted)
    } else {
        fmt.Println("✗ Authentication failed - data tampered")
    }
}
```

### 35.5.3 Stream Ciphers vs Block Ciphers

```
BLOCK CIPHER (AES):
┌─────────────────────────┐
│ Plaintext: "Hola Mundo" │
│ Size: variado           │
└────────────┬────────────┘
             ↓
    ┌────────────────────┐
    │  PKCS7 Padding     │
    │ "Hola Mundo\x05"   │
    │ "..." (16 bytes)   │
    └────────────┬───────┘
             ↓
    ┌────────────────────┐
    │ Encriptar en bloques│
    │ 16 bytes at a time │
    └────────────┬───────┘
             ↓
    ┌────────────────────┐
    │ Ciphertext (16+ B) │
    └────────────────────┘

STREAM CIPHER (ChaCha20):
┌─────────────────────────┐
│ Plaintext: "Hola Mundo" │
│ (any size)              │
└────────────┬────────────┘
             ↓
    ┌────────────────────┐
    │ Generate keystream │
    │ XOR con plaintext  │
    └────────────┬───────┘
             ↓
    ┌────────────────────┐
    │ Ciphertext (same)  │
    │ (no padding)       │
    └────────────────────┘
```

### 35.5.4 ChaCha20-Poly1305

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "golang.org/x/crypto/chacha20poly1305"
    "io"
)

func EncriptarChaCha20(plaintext string, clave string) (string, error) {
    // Expandir/truncar clave a 32 bytes
    key := make([]byte, 32)
    copy(key, []byte(clave))

    cipher, _ := chacha20poly1305.New(key)

    // Generar nonce de 12 bytes
    nonce := make([]byte, 12)
    io.ReadFull(rand.Reader, nonce)

    ciphertext := cipher.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func main() {
    plaintext := "Mensaje cifrado"
    clave := "MiClaveSecreta32BytesMuyBuena!"

    encrypted, _ := EncriptarChaCha20(plaintext, clave)
    fmt.Println("ChaCha20-Poly1305 encrypted:", encrypted)
}
```

---

## 35.6 Encriptación Asimétrica

### 35.6.1 RSA - Rivest-Shamir-Adleman

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
)

// GenerarParRSA genera par de claves RSA
func GenerarParRSA(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, nil, err
    }
    return privateKey, &privateKey.PublicKey, nil
}

// EncriptarRSA encripta con clave pública
func EncriptarRSA(plaintext string, publicKey *rsa.PublicKey) (string, error) {
    ciphertext, err := rsa.EncryptOAEP(
        sha256.New(),
        rand.Reader,
        publicKey,
        []byte(plaintext),
        nil,
    )
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DesencriptarRSA desencripta con clave privada
func DesencriptarRSA(ciphertext64 string, privateKey *rsa.PrivateKey) (string, error) {
    ciphertext, _ := base64.StdEncoding.DecodeString(ciphertext64)

    plaintext, err := rsa.DecryptOAEP(
        sha256.New(),
        rand.Reader,
        privateKey,
        ciphertext,
        nil,
    )
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}

func main() {
    // Generar par de claves (2048 bits mínimo)
    privateKey, publicKey, _ := GenerarParRSA(2048)

    // Alice encripta con clave pública de Bob
    plaintext := "Mensaje solo para Bob"
    encrypted, _ := EncriptarRSA(plaintext, publicKey)
    fmt.Printf("Plaintext:  %s\n", plaintext)
    fmt.Printf("Encrypted:  %s...\n\n", encrypted[:50])

    // Bob desencripta con su clave privada
    decrypted, _ := DesencriptarRSA(encrypted, privateKey)
    fmt.Printf("Decrypted:  %s\n", decrypted)
    fmt.Printf("Coincide:   %v\n", plaintext == decrypted)
}
```

### 35.6.2 Exportar/Importar Claves RSA

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "os"
)

// GuardarClavePrivada guarda clave privada en archivo
func GuardarClavePrivada(filename string, key *rsa.PrivateKey) error {
    privateKey := &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(key),
    }

    file, _ := os.Create(filename)
    defer file.Close()
    pem.Encode(file, privateKey)
    return nil
}

// CargarClavePrivada carga clave privada de archivo
func CargarClavePrivada(filename string) (*rsa.PrivateKey, error) {
    keyData, _ := os.ReadFile(filename)
    block, _ := pem.Decode(keyData)
    return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// GuardarClavePublica guarda clave pública en archivo
func GuardarClavePublica(filename string, key *rsa.PublicKey) error {
    publicKeyBytes, _ := x509.MarshalPKIXPublicKey(key)
    publicKey := &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: publicKeyBytes,
    }

    file, _ := os.Create(filename)
    defer file.Close()
    pem.Encode(file, publicKey)
    return nil
}

// CargarClavePublica carga clave pública de archivo
func CargarClavePublica(filename string) (*rsa.PublicKey, error) {
    keyData, _ := os.ReadFile(filename)
    block, _ := pem.Decode(keyData)

    publicKeyInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
    return publicKeyInterface.(*rsa.PublicKey), nil
}

func main() {
    // Generar par
    privateKey, _, _ := rsa.GenerateKey(rand.Reader, 2048)

    // Guardar
    GuardarClavePrivada("private.pem", privateKey)
    GuardarClavePublica("public.pem", &privateKey.PublicKey)

    // Cargar
    loadedPrivate, _ := CargarClavePrivada("private.pem")
    loadedPublic, _ := CargarClavePublica("public.pem")

    fmt.Println("✓ Claves guardadas y cargadas exitosamente")
}
```

### 35.6.3 Comparación con Simétrica

```
SIMÉTRICA (AES):
┌─────────────┐
│ Clave única │
│  "Secret"   │
└────┬────┬──┘
     │    │
  Alice  Bob  (ambos tienen la misma clave)

Problema: ¿Cómo compartir la clave de forma segura?

ASIMÉTRICA (RSA):
┌──────────────────────┐
│ Par de claves        │
│ Privada: solo Bob    │
│ Pública: todos       │
└──────────────────────┘

Alice:              Bob:
Tiene public_bob    Tiene private_bob y public_bob
Encripta con        Desencripta con
public_bob    ←→    private_bob

No necesita compartir clave privada
```

---

## 35.7 Firmas Digitales

### 35.7.1 Conceptos: Sign & Verify

```go
package main

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
)

// FirmarMensaje firma un mensaje con clave privada
func FirmarMensaje(mensaje string, privateKey *rsa.PrivateKey) (string, error) {
    hash := sha256.Sum256([]byte(mensaje))

    signature, err := rsa.SignPSS(
        rand.Reader,
        privateKey,
        crypto.SHA256,
        hash[:],
        nil,
    )
    if err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(signature), nil
}

// VerificarFirma verifica firma con clave pública
func VerificarFirma(mensaje string, signature64 string, publicKey *rsa.PublicKey) bool {
    hash := sha256.Sum256([]byte(mensaje))
    signature, _ := base64.StdEncoding.DecodeString(signature64)

    err := rsa.VerifyPSS(
        publicKey,
        crypto.SHA256,
        hash[:],
        signature,
        nil,
    )

    return err == nil
}

func main() {
    // Generar claves
    privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
    publicKey := &privateKey.PublicKey

    // Alice firma un mensaje
    mensaje := "Confirmo pago de $1000"
    firma, _ := FirmarMensaje(mensaje, privateKey)

    fmt.Printf("Mensaje: %s\n", mensaje)
    fmt.Printf("Firma:   %s...\n\n", firma[:50])

    // Bob verifica la firma
    if VerificarFirma(mensaje, firma, publicKey) {
        fmt.Println("✓ Firma válida - Mensaje auténtico y no modificado")
    }

    // Intentar modificar mensaje
    mensajeModificado := "Confirmo pago de $10000"
    if !VerificarFirma(mensajeModificado, firma, publicKey) {
        fmt.Println("✗ Firma inválida - Mensaje fue modificado")
    }
}
```

### 35.7.2 ECDSA - Curvas Elípticas

```go
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "fmt"
    "math/big"
)

// GenerarClaveECDSA genera par ECDSA
func GenerarClaveECDSA() (*ecdsa.PrivateKey, error) {
    return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// FirmarECDSA firma con ECDSA
func FirmarECDSA(mensaje string, privateKey *ecdsa.PrivateKey) (string, error) {
    hash := sha256.Sum256([]byte(mensaje))

    r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
    if err != nil {
        return "", err
    }

    // Codificar firma (r, s)
    signature := append(r.Bytes(), s.Bytes()...)
    return fmt.Sprintf("%x", signature), nil
}

// VerificarECDSA verifica firma ECDSA
func VerificarECDSA(mensaje string, signatureHex string, publicKey *ecdsa.PublicKey) bool {
    hash := sha256.Sum256([]byte(mensaje))

    // Decodificar firma
    signatureBytes := make([]byte, hex.DecodedLen(len(signatureHex)))
    hex.Decode(signatureBytes, []byte(signatureHex))

    // Extraer r y s
    r := big.NewInt(0).SetBytes(signatureBytes[:len(signatureBytes)/2])
    s := big.NewInt(0).SetBytes(signatureBytes[len(signatureBytes)/2:])

    return ecdsa.Verify(publicKey, hash[:], r, s)
}

func main() {
    privateKey, _ := GenerarClaveECDSA()
    publicKey := &privateKey.PublicKey

    mensaje := "Contrato inteligente"
    firma, _ := FirmarECDSA(mensaje, privateKey)

    if VerificarECDSA(mensaje, firma, publicKey) {
        fmt.Println("✓ Firma ECDSA válida")
    }
}
```

### 35.7.3 Casos de Uso

```go
// Certificado digital
type Documento struct {
    Contenido string
    Firma     string
    Autoridad string
}

// Protocolo de no-rechazo
type Transaccion struct {
    ID       string
    Monto    float64
    Usuario  string
    Firma    string // Usuario firmó, no puede negar
    Timestamp int64
}

// Blockchain
type Bloque struct {
    Datos          string
    HashAnterior   string
    FirmaCreador   string // Solo creador puede crear bloque válido
}
```

---

## 35.8 Generación de Números Aleatorios Criptográficos

### 35.8.1 crypto/rand vs math/rand

```go
package main

import (
    "crypto/rand"
    "fmt"
    "math/rand"
    "time"
)

func DemostrarDiferencia() {
    // ❌ math/rand - NUNCA para criptografía
    mathRand := rand.New(rand.NewSource(time.Now().UnixNano()))
    weak := make([]byte, 16)
    for i := 0; i < 16; i++ {
        weak[i] = byte(mathRand.Intn(256))
    }
    fmt.Printf("math/rand (inseguro):   %x\n", weak)

    // ✓ crypto/rand - Correcto
    strong := make([]byte, 16)
    rand.Read(strong)
    fmt.Printf("crypto/rand (seguro):   %x\n", strong)

    // Diferencias:
    // math/rand:
    //   - Predecible si conoces la seed
    //   - Usa número anterior para generar siguiente
    //   - Adecuado para juegos, simulaciones
    //
    // crypto/rand:
    //   - Impredecible
    //   - Lee del sistema (/dev/urandom en Unix)
    //   - Debe usarse para seguridad
}

func main() {
    DemostrarDiferencia()
}
```

### 35.8.2 Generar Valores Aleatorios

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "math/big"
)

// GenerarTokenSeguro genera token aleatorio
func GenerarTokenSeguro(longitud int) (string, error) {
    bytes := make([]byte, longitud)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerarNumeroDado genera número entre min y max
func GenerarNumeroDado(min, max int64) (int64, error) {
    nMax := big.NewInt(max - min)
    n, err := rand.Int(rand.Reader, nMax)
    if err != nil {
        return 0, err
    }
    return n.Int64() + min, nil
}

// GenerarUUID genera UUID v4
func GenerarUUID() string {
    b := make([]byte, 16)
    rand.Read(b)
    b[6] = (b[6] & 0x0f) | 0x40
    b[8] = (b[8] & 0x3f) | 0x80

    return fmt.Sprintf("%x-%x-%x-%x-%x",
        b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func main() {
    // Token para resetear contraseña
    token, _ := GenerarTokenSeguro(32)
    fmt.Printf("Reset token: %s\n", token)

    // Número aleatorio
    num, _ := GenerarNumeroDado(1, 100)
    fmt.Printf("Random 1-100: %d\n", num)

    // UUID
    uuid := GenerarUUID()
    fmt.Printf("UUID: %s\n", uuid)
}
```

### 35.8.3 Timing Attacks y Defensa

```go
package main

import (
    "crypto/subtle"
    "fmt"
    "time"
)

// ❌ Comparación insegura
func CompararInsegura(a, b string) bool {
    return a == b  // Tiempo de ejecución varía
}

// ✓ Comparación segura (tiempo constante)
func CompararSegura(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func midiendo() {
    correcto := "password123"
    incorrecto := "a"

    // Medir tiempo inseguro
    start := time.Now()
    for i := 0; i < 100000; i++ {
        CompararInsegura(correcto, incorrecto)
    }
    tiempoInseguro := time.Since(start)

    // Medir tiempo seguro
    start = time.Now()
    for i := 0; i < 100000; i++ {
        CompararSegura(correcto, incorrecto)
    }
    tiempoSeguro := time.Since(start)

    fmt.Printf("Insegura: %v\n", tiempoInseguro)
    fmt.Printf("Segura:   %v\n", tiempoSeguro)
    fmt.Println("\n✓ Segura siempre toma igual tiempo")
}

// ❌ Timing attack
func AtaqueDeTemporización() {
    contraseña := "MySecurePass123"

    // Atacante mide tiempo de respuesta
    intentos := []string{
        "a",           // ~0.001ms
        "Ma",          // ~0.002ms
        "My",          // ~0.003ms (más largo = más caracteres correctos)
        "MyS",         // ~0.004ms
        "MySecure...   // continuar
    }

    for _, intento := range intentos {
        start := time.Now()
        CompararInsegura(contraseña, intento)
        elapsed := time.Since(start)
        fmt.Printf("%s: %v ns\n", intento, elapsed.Nanoseconds())
    }
}

func main() {
    fmt.Println("=== Comparación de tiempos ===")
    midiendo()
}
```

---

## 35.9 Certificados y TLS

### 35.9.1 Certificados X.509

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "math/big"
    "os"
    "time"
)

// GenerarCertificadoAutofirmado crea un certificado self-signed
func GenerarCertificadoAutofirmado(
    nombre string,
    diasValido int,
) (*rsa.PrivateKey, []byte, error) {
    // Generar clave privada
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, nil, err
    }

    // Crear template
    template := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"Mi Empresa"},
            CommonName:   nombre,
        },
        NotBefore: time.Now(),
        NotAfter:  time.Now().AddDate(0, 0, diasValido),

        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
    }

    // Autofirmar
    certBytes, err := x509.CreateCertificate(
        rand.Reader,
        &template,
        &template,
        &privateKey.PublicKey,
        privateKey,
    )
    if err != nil {
        return nil, nil, err
    }

    return privateKey, certBytes, nil
}

// GuardarCertificado guarda certificado en archivo
func GuardarCertificado(filename string, certBytes []byte,
    privateKey *rsa.PrivateKey) error {

    // Guardar certificado
    certPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "CERTIFICATE",
        Bytes: certBytes,
    })

    // Guardar clave privada
    keyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })

    os.WriteFile(filename+".crt", certPEM, 0644)
    os.WriteFile(filename+".key", keyPEM, 0600)

    return nil
}

// AnalizarCertificado analiza un certificado
func AnalizarCertificado(certBytes []byte) error {
    cert, err := x509.ParseCertificate(certBytes)
    if err != nil {
        return err
    }

    fmt.Printf("Subject:     %s\n", cert.Subject)
    fmt.Printf("Issuer:      %s\n", cert.Issuer)
    fmt.Printf("Válido desde: %s\n", cert.NotBefore)
    fmt.Printf("Válido hasta: %s\n", cert.NotAfter)
    fmt.Printf("Serial:      %d\n", cert.SerialNumber)

    return nil
}

func main() {
    // Generar certificado
    privateKey, certBytes, _ := GenerarCertificadoAutofirmado("localhost", 365)

    // Analizar
    AnalizarCertificado(certBytes)

    // Guardar
    GuardarCertificado("server", certBytes, privateKey)

    fmt.Println("\n✓ Certificado generado: server.crt y server.key")
}
```

### 35.9.2 Servidor HTTPS

```go
package main

import (
    "net/http"
    "log"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("¡Conexión segura TLS!"))
    })

    // server.crt y server.key creados anteriormente
    log.Println("Servidor HTTPS en https://localhost:8443")
    log.Fatal(http.ListenAndServeTLS(
        ":8443",
        "server.crt",
        "server.key",
        nil,
    ))
}
```

### 35.9.3 Verificación de Cadena de Certificados

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
)

// VerificarCertificadoServidor verifica certificado TLS
func VerificarCertificadoServidor(hostname string) error {
    conn, err := tls.Dial("tcp", hostname+":443", nil)
    if err != nil {
        return err
    }
    defer conn.Close()

    // Obtener certificado
    cert := conn.ConnectionState().PeerCertificates[0]

    // Opciones de verificación
    opts := x509.VerifyOptions{
        DNSName: hostname,
        Roots:   nil, // Usar certificados del sistema
    }

    // Verificar
    if _, err := cert.Verify(opts); err != nil {
        return fmt.Errorf("certificado inválido: %v", err)
    }

    fmt.Printf("✓ Certificado válido para %s\n", hostname)
    fmt.Printf("  Válido hasta: %s\n", cert.NotAfter)

    return nil
}

func main() {
    VerificarCertificadoServidor("google.com")
}
```

---

## 35.10 Derivación de Claves

### 35.10.1 PBKDF2 - Password-Based KDF

```go
package main

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "golang.org/x/crypto/pbkdf2"
)

// DeriviarClaveContraseña derive una clave de una contraseña
func DeriviarClaveContraseña(contraseña string, salt []byte) string {
    // PBKDF2 aplica función hash múltiples veces
    clave := pbkdf2.Key(
        []byte(contraseña),
        salt,
        100000,         // iteraciones (más = más seguro pero más lento)
        32,             // longitud de salida (256 bits)
        sha256.New,
    )

    return base64.StdEncoding.EncodeToString(clave)
}

func main() {
    contraseña := "MiContraseña123"

    // Generar salt aleatorio
    salt := make([]byte, 16)
    rand.Read(salt)

    // Derivar clave
    clave := DeriviarClaveContraseña(contraseña, salt)

    fmt.Printf("Contraseña: %s\n", contraseña)
    fmt.Printf("Salt:       %x\n", salt)
    fmt.Printf("Clave:      %s\n", clave)

    // Mismo salt + contraseña = misma clave
    clave2 := DeriviarClaveContraseña(contraseña, salt)
    fmt.Printf("\nMisma contraseña + salt: %v\n", clave == clave2)
}
```

### 35.10.2 Scrypt

```go
package main

import (
    "encoding/base64"
    "fmt"
    "golang.org/x/crypto/scrypt"
)

// DeriviarConScrypt usa scrypt (más resistente a GPU)
func DeriviarConScrypt(contraseña string, salt []byte) (string, error) {
    clave, err := scrypt.Key(
        []byte(contraseña),
        salt,
        1<<15,  // N = 2^15 (memory parameter)
        8,      // r (blocksize parameter)
        1,      // p (parallelization parameter)
        32,     // keyLen
    )
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(clave), nil
}

func main() {
    salt := []byte("saltsaltsaltsalt")

    clave, _ := DeriviarConScrypt("contraseña", salt)
    fmt.Printf("Scrypt derivada: %s\n", clave)
}
```

### 35.10.3 Comparación de KDF

| Algoritmo | Resistencia GPU | Resistencia ASIC | Memoria | Velocidad | Recomendación |
|---|---|---|---|---|---|
| **PBKDF2** | Baja | Baja | Mínima | Rápido | Básico |
| **Scrypt** | Alta | Media | Alta | Lento | ✓ Bueno |
| **Argon2** | Muy alta | Muy alta | Variable | Lento | ✓ Mejor |
| **bcrypt** | Alta | Baja | Moderada | Lento | ✓ Muy bueno |

---

## 35.11 Buenas Prácticas y Patrones

### 35.11.1 Gestión de Claves

```go
package main

import (
    "fmt"
    "os"
)

// ❌ ANTIPATRONES
func AntiPatrones() {
    // 1. Hardcodear claves
    const API_KEY = "sk-1234567890abcdef"  // ¡NUNCA!

    // 2. Almacenar en archivos sin protección
    key := "SecretoImportante"
    os.WriteFile("secret.txt", []byte(key), 0644)  // ¡Legible para todos!

    // 3. Pasar claves en URLs
    url := "https://api.example.com?key=sk-1234567890"  // ¡En logs!
}

// ✓ PATRONES CORRECTOS
func PatronesCorrectos() {
    // 1. Usar variables de entorno
    apiKey := os.Getenv("API_KEY")

    // 2. Usar gestor de secretos
    // - HashiCorp Vault
    // - AWS Secrets Manager
    // - Azure Key Vault
    // - Google Secret Manager

    // 3. Pasar claves en headers
    // Authorization: Bearer sk-1234567890

    // 4. Rotar claves regularmente
    // - Cambiar cada 90 días
    // - Cuando un empleado se va
    // - Si hay sospecha de compromiso
}
```

### 35.11.2 Generación Segura de Salt

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

// GenerarSalt crea salt criptográficamente seguro
func GenerarSalt(longitud int) (string, error) {
    bytes := make([]byte, longitud)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}

func main() {
    salt, _ := GenerarSalt(16)
    fmt.Printf("Salt generado: %s\n", salt)
    fmt.Println("\n✓ Usar esto para cada contraseña de usuario")
}
```

### 35.11.3 Auditoría y Logging

```go
package main

import (
    "log"
    "time"
)

// EventoSeguridad registra eventos criptográficos importantes
type EventoSeguridad struct {
    Tipo      string    // "login_failed", "password_changed", etc
    Usuario   string
    Timestamp time.Time
    IP        string
    Detalles  string
}

func RegistrarEvento(evento EventoSeguridad) {
    log.Printf("[SECURITY] %s | Usuario: %s | IP: %s | %s\n",
        evento.Tipo,
        evento.Usuario,
        evento.IP,
        evento.Detalles,
    )
}

func main() {
    // Registrar intento fallido de login
    RegistrarEvento(EventoSeguridad{
        Tipo:      "login_failed",
        Usuario:   "user@example.com",
        Timestamp: time.Now(),
        IP:        "192.168.1.100",
        Detalles:  "Contraseña incorrecta",
    })
}
```

### 35.11.4 Rate Limiting contra Fuerza Bruta

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type RateLimiter struct {
    intentos map[string][]time.Time
    mu       sync.Mutex
    maxInt   int
    ventana  time.Duration
}

func NewRateLimiter(maxIntentosEnVentana int, ventana time.Duration) *RateLimiter {
    return &RateLimiter{
        intentos: make(map[string][]time.Time),
        maxInt:   maxIntentosEnVentana,
        ventana:  ventana,
    }
}

// PermitirIntento verifica si se permite otro intento
func (rl *RateLimiter) PermitirIntento(clave string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    ahora := time.Now()

    // Limpiar intentos antiguos
    var intentoValidos []time.Time
    for _, t := range rl.intentos[clave] {
        if ahora.Sub(t) < rl.ventana {
            intentoValidos = append(intentoValidos, t)
        }
    }
    rl.intentos[clave] = intentoValidos

    // Verificar si se permite
    if len(rl.intentos[clave]) < rl.maxInt {
        rl.intentos[clave] = append(rl.intentos[clave], ahora)
        return true
    }

    return false
}

func main() {
    limiter := NewRateLimiter(3, 1*time.Minute)

    // Simular intentos de login
    for i := 1; i <= 5; i++ {
        if limiter.PermitirIntento("user@example.com") {
            fmt.Printf("Intento %d: ✓ Permitido\n", i)
        } else {
            fmt.Printf("Intento %d: ✗ Bloqueado (rate limit)\n", i)
        }
    }
}
```

### 35.11.5 Encriptación End-to-End

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

// ClienteE2E implementa encriptación E2E
type ClienteE2E struct {
    claveLocal []byte  // Generada localmente, nunca se envía
}

// EncriptarPraEnviar encripta antes de enviar al servidor
func (c *ClienteE2E) EncriptarPraEnviar(mensaje string) (string, error) {
    bloque, _ := aes.NewCipher(c.claveLocal)
    gcm, _ := cipher.NewGCM(bloque)

    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)

    ciphertext := gcm.Seal(nonce, nonce, []byte(mensaje), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func main() {
    // Cliente crea su clave local
    claveLocal := make([]byte, 32)
    rand.Read(claveLocal)

    cliente := &ClienteE2E{claveLocal: claveLocal}

    // Encriptar mensaje
    cifrado, _ := cliente.EncriptarPraEnviar("Mensaje privado")

    // Enviar solo cifrado al servidor
    // Servidor NO puede leer

    // Receptor con misma clave local puede desencriptar
    fmt.Printf("Cifrado para servidor: %s\n", cifrado)
}
```

### 35.11.6 Verificación de Integridad

```go
package main

import (
    "crypto/md5"
    "crypto/sha256"
    "fmt"
    "io"
    "os"
)

// VerificarIntegridadArchivo compara hash de archivo
func VerificarIntegridadArchivo(filepath string, hashEsperado string) bool {
    f, _ := os.Open(filepath)
    defer f.Close()

    h := sha256.New()
    io.Copy(h, f)

    hashReal := fmt.Sprintf("%x", h.Sum(nil))
    return hashReal == hashEsperado
}

func main() {
    // Antes de distribuir archivo:
    // go run ./generar_archivo.go > archivo.bin
    // sha256sum archivo.bin > archivo.sha256

    // Usuario descarga y verifica:
    if VerificarIntegridadArchivo("archivo.bin", "abc123...") {
        fmt.Println("✓ Archivo íntegro")
    } else {
        fmt.Println("✗ Archivo corrupto o modificado")
    }
}
```

---

## Ejercicios Progresivos

### Ejercicio 1: Hash de Archivos

**Objetivo:** Crear herramienta para calcular y verificar SHA256 de archivos

```go
package main

import (
    "crypto/sha256"
    "flag"
    "fmt"
    "io"
    "os"
)

func calcularSHA256(filepath string) (string, error) {
    f, err := os.Open(filepath)
    if err != nil {
        return "", err
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }

    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func main() {
    verify := flag.String("verify", "", "Verificar hash")
    flag.Parse()

    if flag.NArg() == 0 {
        fmt.Println("Uso: hashfile <archivo> [-verify hash]")
        return
    }

    archivo := flag.Arg(0)
    hash, _ := calcularSHA256(archivo)

    fmt.Printf("%s: %s\n", archivo, hash)

    if *verify != "" {
        if hash == *verify {
            fmt.Println("✓ Hash coincide")
        } else {
            fmt.Println("✗ Hash no coincide")
        }
    }
}
```

**Requisitos:**

- Calcular SHA256 de cualquier archivo
- Comparar con hash conocido
- Manejar archivos grandes eficientemente
- Mostrar progreso para archivos > 100MB

---

### Ejercicio 2: Verificación HMAC

**Objetivo:** Implementar sistema de webhook con validación HMAC

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "os"
)

var WEBHOOK_SECRET = os.Getenv("WEBHOOK_SECRET")

func verificarHMAC(body []byte, firma string) bool {
    esperado := hmac.New(sha256.New, []byte(WEBHOOK_SECRET))
    esperado.Write(body)
    expected_sig := hex.EncodeToString(esperado.Sum(nil))

    return hmac.Equal([]byte(firma), []byte(expected_sig))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    firma := r.Header.Get("X-Signature")
    body, _ := io.ReadAll(r.Body)

    if !verificarHMAC(body, firma) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    fmt.Println("✓ Webhook verificado")
    w.WriteHeader(http.StatusOK)
}

func main() {
    http.HandleFunc("/webhook", webhookHandler)
    http.ListenAndServe(":8080", nil)
}
```

---

### Ejercicio 3: Password Hasher

**Objetivo:** Sistema seguro de autenticación con bcrypt

```go
package main

import (
    "bufio"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "os"
    "strings"
    "sync"
)

type UsuarioDB struct {
    hashes map[string]string
    mu     sync.RWMutex
}

func (db *UsuarioDB) Registrar(usuario, contraseña string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(contraseña), 12)
    if err != nil {
        return err
    }

    db.mu.Lock()
    db.hashes[usuario] = string(hash)
    db.mu.Unlock()

    return nil
}

func (db *UsuarioDB) Verificar(usuario, contraseña string) bool {
    db.mu.RLock()
    hash, existe := db.hashes[usuario]
    db.mu.RUnlock()

    if !existe {
        return false
    }

    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(contraseña)) == nil
}

func main() {
    db := &UsuarioDB{hashes: make(map[string]string)}
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        scanner.Scan()
        cmd := strings.Fields(scanner.Text())

        if len(cmd) == 0 {
            continue
        }

        switch cmd[0] {
        case "register":
            if len(cmd) == 3 {
                db.Registrar(cmd[1], cmd[2])
                fmt.Println("✓ Usuario registrado")
            }
        case "login":
            if len(cmd) == 3 {
                if db.Verificar(cmd[1], cmd[2]) {
                    fmt.Println("✓ Login exitoso")
                } else {
                    fmt.Println("✗ Credenciales inválidas")
                }
            }
        case "exit":
            return
        }
    }
}
```

---

### Ejercicio 4: AES Encryption

**Objetivo:** Encriptar/desencriptar archivos con AES-GCM

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "flag"
    "fmt"
    "io"
    "os"
)

func encriptarArchivo(entrada, salida string, clave []byte) error {
    f, _ := os.Open(entrada)
    datos, _ := io.ReadAll(f)
    f.Close()

    bloque, _ := aes.NewCipher(clave)
    gcm, _ := cipher.NewGCM(bloque)

    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)

    cifrado := gcm.Seal(nonce, nonce, datos, nil)

    os.WriteFile(salida, cifrado, 0644)
    return nil
}

func desencriptarArchivo(entrada, salida string, clave []byte) error {
    cifrado, _ := os.ReadFile(entrada)

    bloque, _ := aes.NewCipher(clave)
    gcm, _ := cipher.NewGCM(bloque)

    nonce, cifrado := cifrado[:gcm.NonceSize()], cifrado[gcm.NonceSize():]

    datos, err := gcm.Open(nil, nonce, cifrado, nil)
    if err != nil {
        return fmt.Errorf("desencriptación falló: %v", err)
    }

    os.WriteFile(salida, datos, 0644)
    return nil
}

func main() {
    modo := flag.String("mode", "encrypt", "encrypt o decrypt")
    entrada := flag.String("in", "", "archivo entrada")
    salida := flag.String("out", "", "archivo salida")
    flag.Parse()

    clave := make([]byte, 32)
    rand.Read(clave)

    if *modo == "encrypt" {
        encriptarArchivo(*entrada, *salida, clave)
        fmt.Printf("✓ Encriptado: %s\n", *salida)
    } else {
        desencriptarArchivo(*entrada, *salida, clave)
        fmt.Printf("✓ Desencriptado: %s\n", *salida)
    }
}
```

---

### Ejercicio 5: Verificación de Certificados

**Objetivo:** Cliente HTTPS que verifica certificados

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "flag"
    "fmt"
    "net/http"
    "time"
)

func verificarCertificado(hostname string) error {
    // Conectar con verificación de certificado
    dialConf := &tls.Dialer{
        Config: &tls.Config{
            ServerName: hostname,
        },
        Timeout: 10 * time.Second,
    }

    conn, err := dialConf.Dial("tcp", hostname+":443")
    if err != nil {
        return fmt.Errorf("conexión fallida: %v", err)
    }
    defer conn.Close()

    // Obtener certificado
    tlsConn := conn.(*tls.Conn)
    cert := tlsConn.ConnectionState().PeerCertificates[0]

    // Analizar
    fmt.Printf("Host:        %s\n", cert.Subject.CommonName)
    fmt.Printf("Válido desde: %s\n", cert.NotBefore)
    fmt.Printf("Válido hasta: %s\n", cert.NotAfter)
    fmt.Printf("Emisor:      %s\n", cert.Issuer.String())

    // Verificar fecha
    ahora := time.Now()
    if ahora.After(cert.NotAfter) {
        return fmt.Errorf("certificado expirado")
    }
    if ahora.Before(cert.NotBefore) {
        return fmt.Errorf("certificado aún no válido")
    }

    fmt.Println("✓ Certificado válido")
    return nil
}

func main() {
    hostname := flag.String("host", "google.com", "hostname")
    flag.Parse()

    if err := verificarCertificado(*hostname); err != nil {
        fmt.Printf("✗ Error: %v\n", err)
    }
}
```

---

## Resumen

La criptografía en Go es accesible pero requiere cuidado:

| Tarea | Usa | Nunca |
|---|---|---|
| **Hash de contraseñas** | bcrypt, argon2 | SHA256 directo |
| **Verificar integridad** | SHA256, SHA512 | MD5, SHA1 |
| **Cifrar datos** | AES-GCM | ECB |
| **Firmar mensajes** | RSA-PSS, ECDSA | Homemade |
| **Generar aleatorios** | crypto/rand | math/rand |
| **Almacenar claves** | Vault, env vars | hardcoded |

**Recuerda:**

- ✓ Usa bibliotecas estándar (crypto/*)
- ✓ Implementa rate limiting
- ✓ Audita eventos de seguridad
- ✓ Rota claves regularmente
- ✓ Consulta expertos si es crítico
- ✗ Nunca reinventes criptografía
- ✗ Nunca hardcodees secretos
- ✗ Nunca uses MD5 para seguridad

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/35-hash-y-criptografia/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/35-hash-y-criptografia):

```bash
cd examples/35-hash-y-criptografia
go run .
```
