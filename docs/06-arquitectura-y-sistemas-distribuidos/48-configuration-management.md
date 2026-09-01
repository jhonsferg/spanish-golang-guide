# Capítulo 48: Configuration management

## Índice del Capítulo

1. [Configuración en Go](#481-configuración-en-go)
2. [Variables de Entorno](#482-variables-de-entorno)
3. [Viper Package](#483-viper-package)
4. [Configuración JSON](#484-configuración-json)
5. [Configuración YAML](#485-configuración-yaml)
6. [Jerarquía de Configuración](#486-jerarquía-de-configuración)
7. [Gestión de Secretos](#487-gestión-de-secretos)
8. [Feature Flags](#488-feature-flags)
9. [Validación de Configuración](#489-validación-de-configuración)
10. [Hot Reload](#4810-hot-reload)
11. [Buenas Prácticas](#4811-buenas-prácticas)

---

## 48.1 Configuración en Go

### Concepto Fundamental

La configuración es la información que cambia entre ambientes (desarrollo, testing, producción) sin modificar el código. Go proporciona múltiples fuentes para cargar configuración con diferentes niveles de prioritarios.

### Fuentes de Configuración

```
Flags CLI (máxima prioridad)
         ↑
    Variables de Entorno
         ↑
    Archivos de Configuración
         ↑
    Valores por Defecto (menor prioridad)
```

### 48.1.1 - Métodos de Carga

Go ofrece cuatro estrategias principales para cargar configuración:

#### 1. **Command-Line Flags**

Los flags son el método más directo, ideal para valores que cambian frecuentemente:

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // Definir flags
    host := flag.String("host", "localhost", "Dirección del servidor")
    port := flag.Int("port", 8080, "Puerto del servidor")
    debug := flag.Bool("debug", false, "Modo debug")
    
    flag.Parse()
    
    fmt.Printf("Host: %s\n", *host)
    fmt.Printf("Port: %d\n", *port)
    fmt.Printf("Debug: %v\n", *debug)
}
```

Uso: `./app -host=0.0.0.0 -port=3000 -debug`

#### 2. **Variables de Entorno**

Ideales para deployments en contenedores (Docker, Kubernetes):

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    host := os.Getenv("APP_HOST")
    port := os.Getenv("APP_PORT")
    
    if host == "" {
        host = "localhost"
    }
    if port == "" {
        port = "8080"
    }
    
    fmt.Printf("Host: %s, Port: %s\n", host, port)
}
```

#### 3. **Archivos de Configuración**

JSON, YAML o TOML para configuraciones complejas:

```json
{
    "server": {
        "host": "0.0.0.0",
        "port": 8080,
        "timeout": 30
    },
    "database": {
        "host": "localhost",
        "port": 5432,
        "name": "myapp"
    }
}
```

#### 4. **Configuración en Código**

Valores por defecto compilados o inicialización programática:

```go
type Config struct {
    Host    string
    Port    int
    Timeout int
}

var DefaultConfig = Config{
    Host:    "localhost",
    Port:    8080,
    Timeout: 30,
}
```

### 48.1.2 - Estructura de Configuración

Pattern recomendado en Go:

```go
type ServerConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    ReadTimeout  int    `json:"read_timeout"`
    WriteTimeout int    `json:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
    Database string `json:"database"`
    MaxConns int    `json:"max_connections"`
}

type AppConfig struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    LogLevel string         `json:"log_level"`
    Env      string         `json:"environment"`
}
```

### 48.1.3 - Loader Genérico

Crear un loader reutilizable:

```go
package config

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"
)

type Loader struct {
    configPath string
    env        string
}

func NewLoader() *Loader {
    return &Loader{
        env: getEnv("APP_ENV", "development"),
    }
}

func (l *Loader) Load() (*AppConfig, error) {
    // 1. Valores por defecto
    cfg := &AppConfig{
        Server: ServerConfig{
            Host: "localhost",
            Port: 8080,
        },
    }
    
    // 2. Cargar desde archivo
    configFile := l.configPath
    if configFile == "" {
        configFile = fmt.Sprintf("config.%s.json", l.env)
    }
    
    if err := l.loadFromFile(configFile, cfg); err != nil && !os.IsNotExist(err) {
        return nil, err
    }
    
    // 3. Override desde variables de entorno
    if host := os.Getenv("SERVER_HOST"); host != "" {
        cfg.Server.Host = host
    }
    if port := os.Getenv("SERVER_PORT"); port != "" {
        fmt.Sscanf(port, "%d", &cfg.Server.Port)
    }
    
    return cfg, nil
}

func (l *Loader) loadFromFile(path string, cfg *AppConfig) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, cfg)
}

func getEnv(key, defaultValue string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultValue
}
```

---

## 48.2 Variables de Entorno

### Concepto

Las variables de entorno son el estándar de facto en DevOps moderno, especialmente con Docker y Kubernetes. El principio 12-Factor establece que la configuración debe almacenarse en variables de entorno.

### 48.2.1 - Acceso Básico

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Lectura simple
    dbHost := os.Getenv("DB_HOST")
    fmt.Println("DB Host:", dbHost) // "" si no existe
    
    // Lectura con valor por defecto
    dbPort := getEnvOrDefault("DB_PORT", "5432")
    fmt.Println("DB Port:", dbPort)
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 48.2.2 - Conversión de Tipos

Go no tiene conversión automática para env vars (todas son strings):

```go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

// Parseo de int
func getEnvInt(key string, defaultValue int) int {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    
    intVal, err := strconv.Atoi(value)
    if err != nil {
        fmt.Printf("Error parsing %s as int: %v\n", key, err)
        return defaultValue
    }
    return intVal
}

// Parseo de bool
func getEnvBool(key string, defaultValue bool) bool {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    
    boolVal, err := strconv.ParseBool(value)
    if err != nil {
        fmt.Printf("Error parsing %s as bool: %v\n", key, err)
        return defaultValue
    }
    return boolVal
}

// Parseo de duration
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    
    dur, err := time.ParseDuration(value)
    if err != nil {
        fmt.Printf("Error parsing %s as duration: %v\n", key, err)
        return defaultValue
    }
    return dur
}

// Uso
func main() {
    port := getEnvInt("PORT", 8080)
    debug := getEnvBool("DEBUG", false)
    timeout := getEnvDuration("TIMEOUT", 30*time.Second)
    
    fmt.Printf("Port: %d, Debug: %v, Timeout: %v\n", port, debug, timeout)
}
```

### 48.2.3 - Validación Inmediata

```go
package config

import (
    "fmt"
    "os"
    "strconv"
)

func getEnvIntRequired(key string) (int, error) {
    value := os.Getenv(key)
    if value == "" {
        return 0, fmt.Errorf("variable de entorno requerida no encontrada: %s", key)
    }
    
    intVal, err := strconv.Atoi(value)
    if err != nil {
        return 0, fmt.Errorf("valor inválido para %s: %s", key, value)
    }
    
    if intVal <= 0 {
        return 0, fmt.Errorf("%s debe ser mayor a 0, recibido: %d", key, intVal)
    }
    
    return intVal, nil
}

// Uso
func main() {
    port, err := getEnvIntRequired("PORT")
    if err != nil {
        panic(err)
    }
    fmt.Println("Port:", port)
}
```

### 48.2.4 - Variables de Entorno en .env

Usar archivo `.env` para desarrollo local:

```
# .env
APP_NAME=MyApp
APP_ENV=development
APP_DEBUG=true

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapp_dev

REDIS_HOST=localhost
REDIS_PORT=6379

LOG_LEVEL=debug
```

Cargar con godotenv:

```go
package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "os"
)

func main() {
    // Cargar .env si existe (ignorar si no existe)
    godotenv.Load(".env")
    
    appName := os.Getenv("APP_NAME")
    dbHost := os.Getenv("DB_HOST")
    
    fmt.Printf("App: %s, DB: %s\n", appName, dbHost)
}
```

### 48.2.5 - Namespace de Variables

Patrón recomendado para aplicaciones monolíticas:

```
# Servidor
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_TIMEOUT=30s

# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=app
DB_PASSWORD=secret

# Cache
CACHE_TYPE=redis
CACHE_HOST=localhost
CACHE_PORT=6379

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## 48.3 Viper Package

### Concepto

Viper es la librería más popular para configuración en Go. Soporta múltiples formatos (JSON, YAML, TOML, HCL), jerarquía de configuración y override desde env vars.

### 48.3.1 - Instalación y Setup Básico

```bash
go get github.com/spf13/viper
```

Uso básico:

```go
package main

import (
    "fmt"
    "github.com/spf13/viper"
)

func main() {
    // Establecer valores por defecto
    viper.SetDefault("server.host", "localhost")
    viper.SetDefault("server.port", 8080)
    
    // Leer configuración
    host := viper.GetString("server.host")
    port := viper.GetInt("server.port")
    
    fmt.Printf("Host: %s, Port: %d\n", host, port)
}
```

### 48.3.2 - Carga de Archivos

```go
package config

import (
    "fmt"
    "github.com/spf13/viper"
    "path/filepath"
)

func LoadViperConfig(env string) error {
    // Configurar ubicación del archivo
    viper.SetConfigName(fmt.Sprintf("config.%s", env))
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.AddConfigPath(".")
    
    // Leer el archivo
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return fmt.Errorf("error leyendo configuración: %w", err)
        }
        // No encontrado, usar valores por defecto
    }
    
    fmt.Println("Configuración cargada desde:", viper.ConfigFileUsed())
    return nil
}

// Uso
func main() {
    if err := LoadViperConfig("production"); err != nil {
        panic(err)
    }
}
```

### 48.3.3 - Override desde Variables de Entorno

```go
package config

import (
    "github.com/spf13/viper"
    "strings"
)

func ConfigureEnvironmentOverride() {
    // Habilitar binding de env vars
    viper.SetEnvPrefix("APP")
    viper.AutomaticEnv()
    
    // Mapeo personalizado: env var a config path
    viper.BindEnv("database.host", "DB_HOST")
    viper.BindEnv("database.port", "DB_PORT")
    viper.BindEnv("database.user", "DB_USER")
    viper.BindEnv("database.password", "DB_PASSWORD")
    
    // Configurar sustitución de caracteres
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// Uso: APP_SERVER_PORT=3000 ./app
// Accede con: viper.GetInt("server.port") -> 3000
```

### 48.3.4 - Unmarshal a Struct

```go
package config

import (
    "fmt"
    "github.com/spf13/viper"
)

type ServerConfig struct {
    Host        string `mapstructure:"host"`
    Port        int    `mapstructure:"port"`
    ReadTimeout int    `mapstructure:"read_timeout"`
}

type AppConfig struct {
    Server   ServerConfig `mapstructure:"server"`
    LogLevel string       `mapstructure:"log_level"`
}

func LoadConfig() (*AppConfig, error) {
    viper.SetDefault("server.host", "localhost")
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("log_level", "info")
    
    var config AppConfig
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("error unmarshal config: %w", err)
    }
    
    return &config, nil
}

// Uso
func main() {
    cfg, err := LoadConfig()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Config: %+v\n", cfg)
}
```

### 48.3.5 - Watchers para Cambios

```go
package config

import (
    "fmt"
    "github.com/spf13/viper"
    "log"
)

func LoadConfigWithWatcher() {
    viper.SetConfigName("config.dev")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    
    if err := viper.ReadInConfig(); err != nil {
        log.Fatal(err)
    }
    
    // Configurar watcher
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        fmt.Println("Configuración modificada:", e.Name)
        
        var config AppConfig
        if err := viper.Unmarshal(&config); err != nil {
            fmt.Println("Error reloadear config:", err)
            return
        }
        
        // Aplicar nuevos valores
        fmt.Printf("Nueva configuración: %+v\n", config)
    })
}
```

---

## 48.4 Configuración JSON

### 48.4.1 - Parsing Básico

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Config struct {
    Server struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"server"`
    Database struct {
        Host     string `json:"host"`
        User     string `json:"user"`
        Password string `json:"password"`
    } `json:"database"`
}

func main() {
    data, err := os.ReadFile("config.json")
    if err != nil {
        panic(err)
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        panic(err)
    }
    
    fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
```

### 48.4.2 - Validación JSON

```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

func LoadAndValidateJSON(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("error leyendo archivo: %w", err)
    }
    
    // Validar JSON válido
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return nil, fmt.Errorf("JSON inválido: %w", err)
    }
    
    // Validar campos requeridos
    if _, ok := raw["server"]; !ok {
        return nil, fmt.Errorf("campo requerido 'server' no encontrado")
    }
    
    // Parse a struct
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("error parseando config: %w", err)
    }
    
    return &cfg, nil
}
```

### 48.4.3 - Valores por Defecto en JSON

```go
package config

import (
    "encoding/json"
)

type ServerConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

func (s *ServerConfig) UnmarshalJSON(data []byte) error {
    // Valores por defecto
    type Alias ServerConfig
    aux := &struct {
        *Alias
    }{
        Alias: &Alias{
            Host: "0.0.0.0",
            Port: 8080,
        },
    }
    
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    
    *s = ServerConfig(*aux.Alias)
    return nil
}
```

---

## 48.5 Configuración YAML

### 48.5.1 - Parsing YAML

```bash
go get gopkg.in/yaml.v3
```

```go
package main

import (
    "fmt"
    "os"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Pool     int    `yaml:"pool_size"`
    } `yaml:"database"`
}

func main() {
    data, _ := os.ReadFile("config.yaml")
    
    var cfg Config
    yaml.Unmarshal(data, &cfg)
    
    fmt.Printf("DB Pool: %d\n", cfg.Database.Pool)
}
```

### 48.5.2 - Estructura YAML Compleja

```yaml
# config.production.yaml
server:
  host: 0.0.0.0
  port: 8080
  timeouts:
    read: 30
    write: 30
    idle: 60
  tls:
    enabled: true
    cert_path: /etc/ssl/certs/server.crt
    key_path: /etc/ssl/private/server.key

database:
  host: db.example.com
  port: 5432
  user: app_user
  password: ${DB_PASSWORD}  # Sustituir desde env
  name: production_db
  pool:
    min: 5
    max: 20
    timeout: 30

logging:
  level: info
  format: json
  outputs:
    - stdout
    - file
  file:
    path: /var/log/app.log
    max_size_mb: 100
    max_backups: 10

features:
  auth:
    enabled: true
    provider: oauth2
  caching:
    enabled: true
    ttl: 3600
```

### 48.5.3 - Anidamiento y Referencias

```go
type ServerTimeouts struct {
    Read  int `yaml:"read"`
    Write int `yaml:"write"`
    Idle  int `yaml:"idle"`
}

type ServerTLS struct {
    Enabled  bool   `yaml:"enabled"`
    CertPath string `yaml:"cert_path"`
    KeyPath  string `yaml:"key_path"`
}

type ServerConfig struct {
    Host     string         `yaml:"host"`
    Port     int            `yaml:"port"`
    Timeouts ServerTimeouts `yaml:"timeouts"`
    TLS      ServerTLS      `yaml:"tls"`
}

type DatabasePool struct {
    Min     int `yaml:"min"`
    Max     int `yaml:"max"`
    Timeout int `yaml:"timeout"`
}

type DatabaseConfig struct {
    Host     string         `yaml:"host"`
    Port     int            `yaml:"port"`
    User     string         `yaml:"user"`
    Password string         `yaml:"password"`
    Name     string         `yaml:"name"`
    Pool     DatabasePool   `yaml:"pool"`
}

type AppConfig struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
}
```

---

## 48.6 Jerarquía de Configuración

### Concepto

La jerarquía establece la prioridad de configuración. Valores posteriores override anteriores:

```
1. Valores por Defecto (Código)
          ↓
2. Archivos de Configuración (config.yaml)
          ↓
3. Variables de Entorno (ENV VARS)
          ↓
4. Flags de CLI (máxima prioridad)
```

### 48.6.1 - Implementación de Jerarquía

```go
package config

import (
    "flag"
    "fmt"
    "os"
    "strconv"
    "github.com/spf13/viper"
)

type Config struct {
    Server struct {
        Host string
        Port int
    }
    Database struct {
        Host string
        Port int
    }
    LogLevel string
}

type ConfigLoader struct {
    configFile string
    env        string
}

func NewConfigLoader() *ConfigLoader {
    return &ConfigLoader{
        env: os.Getenv("APP_ENV"),
        if os.Getenv("APP_ENV") == "" {
            "development"
        },
    }
}

func (cl *ConfigLoader) Load() (*Config, error) {
    cfg := &Config{
        Server: struct {
            Host string
            Port int
        }{
            Host: "127.0.0.1",
            Port: 8080,
        },
        Database: struct {
            Host string
            Port int
        }{
            Host: "localhost",
            Port: 5432,
        },
        LogLevel: "info",
    }
    
    // 1. Cargar valores por defecto (ya asignados arriba)
    
    // 2. Override desde archivo
    if err := cl.loadFromFile(cfg); err != nil && !os.IsNotExist(err) {
        return nil, fmt.Errorf("error cargando config: %w", err)
    }
    
    // 3. Override desde variables de entorno
    cl.overrideFromEnv(cfg)
    
    // 4. Override desde flags CLI
    cl.overrideFromFlags(cfg)
    
    return cfg, nil
}

func (cl *ConfigLoader) loadFromFile(cfg *Config) error {
    if cl.configFile == "" {
        cl.configFile = fmt.Sprintf("config.%s.yaml", cl.env)
    }
    
    viper.SetConfigFile(cl.configFile)
    viper.SetConfigType("yaml")
    
    return viper.ReadInConfig()
}

func (cl *ConfigLoader) overrideFromEnv(cfg *Config) {
    if host := os.Getenv("SERVER_HOST"); host != "" {
        cfg.Server.Host = host
    }
    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil {
            cfg.Server.Port = port
        }
    }
    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        cfg.LogLevel = logLevel
    }
}

func (cl *ConfigLoader) overrideFromFlags(cfg *Config) {
    host := flag.String("server-host", cfg.Server.Host, "Host del servidor")
    port := flag.Int("server-port", cfg.Server.Port, "Puerto del servidor")
    logLevel := flag.String("log-level", cfg.LogLevel, "Nivel de logging")
    
    flag.Parse()
    
    cfg.Server.Host = *host
    cfg.Server.Port = *port
    cfg.LogLevel = *logLevel
}
```

### 48.6.2 - Diagrama de Precedencia

```
Archivo de Config: config.production.yaml
┌──────────────────────────────────────┐
│ server:                              │
│   host: prod.example.com (ARCHIVO)   │
│   port: 8080 (ARCHIVO)               │
└──────────────────────────────────────┘
         Override por ENV
         SERVER_HOST=0.0.0.0
┌──────────────────────────────────────┐
│ server:                              │
│   host: 0.0.0.0 (ENV)                │
│   port: 8080 (ARCHIVO)               │
└──────────────────────────────────────┘
         Override por FLAGS
         -server-port=9000
┌──────────────────────────────────────┐
│ server:                              │
│   host: 0.0.0.0 (ENV)                │
│   port: 9000 (FLAG)                  │
└──────────────────────────────────────┘
```

---

## 48.7 Gestión de Secretos

### 48.7.1 - Principio 12-Factor

El 12-Factor App establece que los secretos NO deben compilarse en el código ni guardarse en repositorio.

### Anti-patrón: Secrets en Código

```go
// ❌ NUNCA HACER ESTO
const (
    DatabasePassword = "super_secret_password_123"
    APIKey           = "sk_prod_abcd1234"
    JWTSecret        = "jwt_signing_key"
)
```

### 48.7.2 - Gestión Correcta con Variables de Entorno

```go
package config

import (
    "fmt"
    "os"
)

func GetDatabasePassword() (string, error) {
    pwd := os.Getenv("DB_PASSWORD")
    if pwd == "" {
        return "", fmt.Errorf("DB_PASSWORD no configurada")
    }
    return pwd, nil
}

func GetAPIKey() (string, error) {
    key := os.Getenv("API_KEY")
    if key == "" {
        return "", fmt.Errorf("API_KEY no configurada")
    }
    return key, nil
}

// Uso
func main() {
    pwd, err := GetDatabasePassword()
    if err != nil {
        panic(err)
    }
    fmt.Println("Database configured")
}
```

### 48.7.3 - Integración con HashiCorp Vault

```bash
go get github.com/hashicorp/vault/api
```

```go
package config

import (
    "fmt"
    "github.com/hashicorp/vault/api"
)

type VaultClient struct {
    client *api.Client
    path   string
}

func NewVaultClient(addr, token, secretPath string) (*VaultClient, error) {
    config := &api.Config{Address: addr}
    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("error creando vault client: %w", err)
    }
    
    client.SetToken(token)
    
    return &VaultClient{
        client: client,
        path:   secretPath,
    }, nil
}

func (vc *VaultClient) GetSecret(key string) (string, error) {
    secret, err := vc.client.Logical().Read(vc.path)
    if err != nil {
        return "", fmt.Errorf("error leyendo secreto: %w", err)
    }
    
    if secret == nil || secret.Data == nil {
        return "", fmt.Errorf("secreto no encontrado")
    }
    
    value, ok := secret.Data["data"].(map[string]interface{})[key]
    if !ok {
        return "", fmt.Errorf("clave %s no encontrada en secreto", key)
    }
    
    return value.(string), nil
}

// Uso
func main() {
    vault, err := NewVaultClient(
        "http://vault.example.com:8200",
        "hvs.token",
        "secret/data/myapp",
    )
    if err != nil {
        panic(err)
    }
    
    dbPassword, err := vault.GetSecret("db_password")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Password retrieved from Vault")
}
```

### 48.7.4 - Archivo de Secretos Protegido

Para desarrollo local, usar archivo `.env.local` (NO en git):

```
# .env (en git, valores placeholder)
DB_PASSWORD=CHANGE_ME
API_KEY=CHANGE_ME

# .env.local (NO en git, valores reales)
DB_PASSWORD=actualPassword123
API_KEY=sk_prod_real_key
```

Cargar con prioridad:

```go
package main

import (
    "github.com/joho/godotenv"
    "os"
)

func main() {
    // Cargar .env primero (valores por defecto/ejemplo)
    godotenv.Load(".env")
    
    // Override con .env.local si existe
    godotenv.Load(".env.local")
    
    // Las variables de entorno del sistema tienen prioridad máxima
    password := os.Getenv("DB_PASSWORD")
}
```

### 48.7.5 - Rotación de Secretos

Implementar reloader de secretos:

```go
package config

import (
    "fmt"
    "sync"
    "time"
)

type SecretManager struct {
    mu          sync.RWMutex
    secrets     map[string]string
    refreshFunc func() (map[string]string, error)
    ttl         time.Duration
    lastRefresh time.Time
}

func NewSecretManager(ttl time.Duration, 
    refreshFunc func() (map[string]string, error)) *SecretManager {
    return &SecretManager{
        secrets:     make(map[string]string),
        refreshFunc: refreshFunc,
        ttl:         ttl,
    }
}

func (sm *SecretManager) Get(key string) (string, error) {
    sm.mu.RLock()
    
    // Verificar si necesita refresh
    if time.Since(sm.lastRefresh) > sm.ttl {
        sm.mu.RUnlock()
        if err := sm.Refresh(); err != nil {
            return "", err
        }
        sm.mu.RLock()
    }
    
    value, ok := sm.secrets[key]
    sm.mu.RUnlock()
    
    if !ok {
        return "", fmt.Errorf("secreto %s no encontrado", key)
    }
    
    return value, nil
}

func (sm *SecretManager) Refresh() error {
    secrets, err := sm.refreshFunc()
    if err != nil {
        return fmt.Errorf("error refrescando secretos: %w", err)
    }
    
    sm.mu.Lock()
    sm.secrets = secrets
    sm.lastRefresh = time.Now()
    sm.mu.Unlock()
    
    return nil
}
```

---

## 48.8 Feature Flags

### Concepto

Feature flags permiten habilitar/deshabilitar funcionalidades sin deployar código, ideal para A/B testing y roll-outs graduales.

### 48.8.1 - Implementación Simple

```go
package features

import (
    "os"
    "strconv"
)

type FeatureFlags struct {
    flags map[string]bool
}

func New() *FeatureFlags {
    return &FeatureFlags{
        flags: map[string]bool{
            "new_dashboard":    false,
            "beta_api_v2":      false,
            "experimental_ui":  false,
            "advanced_search":  false,
        },
    }
}

func (ff *FeatureFlags) LoadFromEnv() {
    for key := range ff.flags {
        envKey := "FEATURE_" + strings.ToUpper(key)
        if val := os.Getenv(envKey); val != "" {
            ff.flags[key], _ = strconv.ParseBool(val)
        }
    }
}

func (ff *FeatureFlags) IsEnabled(feature string) bool {
    enabled, exists := ff.flags[feature]
    if !exists {
        return false
    }
    return enabled
}

// Uso
func main() {
    features := New()
    features.LoadFromEnv()
    
    if features.IsEnabled("new_dashboard") {
        // Cargar dashboard nuevo
    } else {
        // Cargar dashboard antiguo
    }
}
```

### 48.8.2 - Feature Flags Basados en Usuario

```go
package features

import (
    "math"
    "hash/fnv"
)

type UserFeatureFlags struct {
    flags       map[string]bool
    rolloutMap  map[string]int // Feature -> porcentaje (0-100)
}

func (uff *UserFeatureFlags) IsEnabledForUser(feature string, userID string) bool {
    // Si está globalmente deshabilitado, retornar false
    if !uff.flags[feature] {
        return false
    }
    
    rollout := uff.rolloutMap[feature]
    if rollout == 100 {
        return true // 100% rollout
    }
    if rollout == 0 {
        return false // 0% rollout
    }
    
    // Hash del userID para distribuir consistentemente
    hash := fnv.New32a()
    hash.Write([]byte(feature + userID))
    hashValue := hash.Sum32()
    
    return (hashValue % 100) < uint32(rollout)
}

// Uso: Feature flags para el 50% de usuarios
func main() {
    uff := &UserFeatureFlags{
        flags: map[string]bool{
            "new_algorithm": true,
        },
        rolloutMap: map[string]int{
            "new_algorithm": 50, // 50% de usuarios
        },
    }
    
    // Usuario "user123" puede obtener new_algorithm si hash cae en 50%
    if uff.IsEnabledForUser("new_algorithm", "user123") {
        // Usar nuevo algoritmo
    }
}
```

### 48.8.3 - Feature Flags con Configuración Persistente

```go
package features

import (
    "encoding/json"
    "os"
    "sync"
)

type FeatureFlagConfig struct {
    Feature   string `json:"feature"`
    Enabled   bool   `json:"enabled"`
    Rollout   int    `json:"rollout"` // 0-100
    Owner     string `json:"owner"`
    CreatedAt string `json:"created_at"`
}

type PersistentFeatureFlags struct {
    mu      sync.RWMutex
    config  map[string]FeatureFlagConfig
    file    string
}

func NewPersistentFlags(configFile string) (*PersistentFeatureFlags, error) {
    pff := &PersistentFeatureFlags{
        config: make(map[string]FeatureFlagConfig),
        file:   configFile,
    }
    
    if err := pff.Load(); err != nil {
        return nil, err
    }
    
    return pff, nil
}

func (pff *PersistentFeatureFlags) Load() error {
    data, err := os.ReadFile(pff.file)
    if err != nil {
        return err
    }
    
    pff.mu.Lock()
    defer pff.mu.Unlock()
    
    var flags []FeatureFlagConfig
    if err := json.Unmarshal(data, &flags); err != nil {
        return err
    }
    
    for _, flag := range flags {
        pff.config[flag.Feature] = flag
    }
    
    return nil
}

func (pff *PersistentFeatureFlags) Save() error {
    pff.mu.RLock()
    flags := make([]FeatureFlagConfig, 0, len(pff.config))
    for _, flag := range pff.config {
        flags = append(flags, flag)
    }
    pff.mu.RUnlock()
    
    data, err := json.MarshalIndent(flags, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(pff.file, data, 0644)
}

func (pff *PersistentFeatureFlags) IsEnabled(feature string) bool {
    pff.mu.RLock()
    defer pff.mu.RUnlock()
    
    cfg, exists := pff.config[feature]
    return exists && cfg.Enabled
}

func (pff *PersistentFeatureFlags) SetFeature(cfg FeatureFlagConfig) error {
    pff.mu.Lock()
    pff.config[cfg.Feature] = cfg
    pff.mu.Unlock()
    
    return pff.Save()
}
```

Archivo de configuración:

```json
[
  {
    "feature": "new_dashboard",
    "enabled": true,
    "rollout": 100,
    "owner": "frontend-team",
    "created_at": "2024-01-15T10:00:00Z"
  },
  {
    "feature": "beta_api_v2",
    "enabled": true,
    "rollout": 25,
    "owner": "backend-team",
    "created_at": "2024-01-10T14:30:00Z"
  }
]
```

---

## 48.9 Validación de Configuración

### 48.9.1 - Validación Básica

```go
package config

import (
    "fmt"
    "net"
    "strings"
)

type Config struct {
    Host     string
    Port     int
    Timeout  int
    LogLevel string
}

func (c *Config) Validate() error {
    // Validar Host
    if c.Host == "" {
        return fmt.Errorf("Host es requerido")
    }
    if !isValidHost(c.Host) {
        return fmt.Errorf("Host inválido: %s", c.Host)
    }
    
    // Validar Port
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("Puerto debe estar entre 1 y 65535, recibido: %d", c.Port)
    }
    
    // Validar Timeout
    if c.Timeout < 0 {
        return fmt.Errorf("Timeout no puede ser negativo: %d", c.Timeout)
    }
    
    // Validar LogLevel
    validLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    if !validLevels[strings.ToLower(c.LogLevel)] {
        return fmt.Errorf("LogLevel inválido: %s", c.LogLevel)
    }
    
    return nil
}

func isValidHost(host string) bool {
    // Validar hostname o IP
    if host == "localhost" || host == "0.0.0.0" || host == "127.0.0.1" {
        return true
    }
    return net.ParseIP(host) != nil
}
```

### 48.9.2 - Validación con go-playground/validator

```bash
go get github.com/go-playground/validator/v10
```

```go
package config

import (
    "github.com/go-playground/validator/v10"
)

type DatabaseConfig struct {
    Host     string `validate:"required,hostname_port"`
    Port     int    `validate:"required,min=1,max=65535"`
    User     string `validate:"required,min=1"`
    Password string `validate:"required,min=8"`
    Database string `validate:"required,min=1"`
    Pool     struct {
        Min int `validate:"required,min=1,max=100"`
        Max int `validate:"required,min=1,max=1000,gtfield=Min"`
    } `validate:"required"`
}

func (db *DatabaseConfig) Validate() error {
    validate := validator.New()
    return validate.Struct(db)
}

// Uso
func main() {
    cfg := &DatabaseConfig{
        Host:     "db.example.com:5432",
        Port:     5432,
        User:     "postgres",
        Password: "securepassword123",
        Database: "myapp",
        Pool: struct {
            Min int
            Max int
        }{Min: 5, Max: 20},
    }
    
    if err := cfg.Validate(); err != nil {
        fmt.Println("Validación fallida:", err)
    }
}
```

### 48.9.3 - Validación al Startup

```go
package main

import (
    "fmt"
    "os"
)

type AppConfig struct {
    Server   ServerConfig
    Database DatabaseConfig
}

func main() {
    // Cargar configuración
    cfg := loadConfig()
    
    // Validar al startup
    if err := cfg.Validate(); err != nil {
        fmt.Fprintf(os.Stderr, "Configuración inválida: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Configuración válida. Iniciando aplicación...")
}

func (ac *AppConfig) Validate() error {
    if err := ac.Server.Validate(); err != nil {
        return fmt.Errorf("validación server fallida: %w", err)
    }
    
    if err := ac.Database.Validate(); err != nil {
        return fmt.Errorf("validación database fallida: %w", err)
    }
    
    return nil
}
```

---

## 48.10 Hot Reload

### Concepto

Hot reload permite recargar configuración sin reiniciar la aplicación, mejorando la disponibilidad.

### 48.10.1 - Watcher de Archivos

```go
package config

import (
    "fmt"
    "github.com/fsnotify/fsnotify"
    "log"
    "sync"
)

type ConfigWatcher struct {
    mu       sync.RWMutex
    config   *AppConfig
    watcher  *fsnotify.Watcher
    filePath string
    onChange func(*AppConfig)
}

func NewConfigWatcher(filePath string, onChange func(*AppConfig)) (*ConfigWatcher, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, fmt.Errorf("error creando watcher: %w", err)
    }
    
    cw := &ConfigWatcher{
        watcher:  watcher,
        filePath: filePath,
        onChange: onChange,
    }
    
    if err := watcher.Add(filePath); err != nil {
        return nil, fmt.Errorf("error watching file: %w", err)
    }
    
    go cw.watch()
    
    return cw, nil
}

func (cw *ConfigWatcher) watch() {
    for {
        select {
        case event, ok := <-cw.watcher.Events:
            if !ok {
                return
            }
            
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("Detectado cambio en config:", event.Name)
                
                // Recargar configuración
                newConfig, err := loadConfigFromFile(cw.filePath)
                if err != nil {
                    log.Println("Error recargando config:", err)
                    continue
                }
                
                // Validar nueva configuración
                if err := newConfig.Validate(); err != nil {
                    log.Println("Configuración inválida:", err)
                    continue
                }
                
                // Actualizar
                cw.mu.Lock()
                cw.config = newConfig
                cw.mu.Unlock()
                
                // Notificar cambio
                if cw.onChange != nil {
                    cw.onChange(newConfig)
                }
            }
            
        case err, ok := <-cw.watcher.Errors:
            if !ok {
                return
            }
            log.Println("Error en watcher:", err)
        }
    }
}

func (cw *ConfigWatcher) Get() *AppConfig {
    cw.mu.RLock()
    defer cw.mu.RUnlock()
    return cw.config
}

func (cw *ConfigWatcher) Close() error {
    return cw.watcher.Close()
}
```

### 48.10.2 - Hot Reload Seguro

```go
package config

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type SafeConfig struct {
    mu          sync.RWMutex
    config      *AppConfig
    callbacks   []func(*AppConfig)
    graceful    bool
}

func (sc *SafeConfig) Update(newConfig *AppConfig) error {
    // Validar antes de aplicar
    if err := newConfig.Validate(); err != nil {
        return fmt.Errorf("configuración inválida: %w", err)
    }
    
    sc.mu.Lock()
    oldConfig := sc.config
    sc.config = newConfig
    sc.mu.Unlock()
    
    // Ejecutar callbacks de forma segura
    for _, cb := range sc.callbacks {
        go func(callback func(*AppConfig)) {
            defer func() {
                if r := recover(); r != nil {
                    fmt.Printf("Callback panic: %v\n", r)
                }
            }()
            callback(newConfig)
        }(cb)
    }
    
    // Log de cambios
    printConfigDiff(oldConfig, newConfig)
    
    return nil
}

func (sc *SafeConfig) Get() *AppConfig {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    return sc.config
}

func (sc *SafeConfig) OnChange(callback func(*AppConfig)) {
    sc.mu.Lock()
    sc.callbacks = append(sc.callbacks, callback)
    sc.mu.Unlock()
}

func printConfigDiff(old, new *AppConfig) {
    if old == nil {
        fmt.Println("Configuración inicial cargada")
        return
    }
    // Comparar campos y loguear cambios
}
```

### 48.10.3 - Graceful Transition

```go
type ServerConfig struct {
    // ... otros campos
    Handler      http.Handler `json:"-"`
}

// Recargar handler sin interrumpir requests en curso
func (sc *SafeConfig) OnConfigChange(newCfg *AppConfig) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Graceful shutdown del handler antiguo
    oldHandler := sc.config.Server.Handler
    
    // Crear nuevo handler
    newCfg.Server.Handler = createNewHandler(newCfg)
    
    // Transition: Los nuevos requests usan newHandler
    // Los requests antiguos continúan en oldHandler
    
    // Esperar a que terminen requests antiguos
    time.Sleep(5 * time.Second)
}
```

---

## 48.11 Buenas Prácticas

### 48.11.1 - Seguridad

#### ✅ **Hacer:**

```go
// ✅ Usar variables de entorno para secretos
password := os.Getenv("DB_PASSWORD")

// ✅ Validar configuración al startup
if err := cfg.Validate(); err != nil {
    log.Fatal(err)
}

// ✅ Loguear configuración sin secretos
logConfig := struct {
    Host string
    Port int
}{
    Host: cfg.Server.Host,
    Port: cfg.Server.Port,
}
log.Printf("Server config: %+v", logConfig)

// ✅ Usar archivos .env.local para desarrollo
godotenv.Load(".env.local")
```

#### ❌ **No Hacer:**

```go
// ❌ Hardcoded secrets
const DatabasePassword = "password123"

// ❌ Logging de secretos
log.Printf("Config: %+v", cfg) // Puede loguear password

// ❌ Secrets en git
// commit: DB_PASSWORD=secret123

// ❌ Archivos de configuración con secretos
// config.json: {"password": "actual_password"}
```

### 48.11.2 - Versionado de Configuración

```yaml
# config.v1.yaml
version: "1.0"
server:
  host: localhost
  port: 8080

# config.v2.yaml (incompatible change)
version: "2.0"
server:
  address: "localhost:8080"  # renamed from host/port

# Migration function
func migrateConfigV1ToV2(v1Config map[string]interface{}) map[string]interface{} {
    v2Config := make(map[string]interface{})
    server := v1Config["server"].(map[string]interface{})
    
    address := fmt.Sprintf("%s:%d",
        server["host"],
        server["port"],
    )
    
    v2Config["version"] = "2.0"
    v2Config["server"] = map[string]interface{}{
        "address": address,
    }
    
    return v2Config
}
```

### 48.11.3 - Documentación

```go
// config.go

// Config es la configuración principal de la aplicación.
//
// Soporta múltiples fuentes (en orden de prioridad):
// 1. Variables de entorno (máxima prioridad)
// 2. Archivos de configuración (config.{env}.yaml)
// 3. Valores por defecto (mínima prioridad)
//
// Ejemplo de uso:
//     cfg := Load()
//     if err := cfg.Validate(); err != nil {
//         log.Fatal(err)
//     }
type Config struct {
    // Server es la configuración del servidor HTTP
    Server ServerConfig `yaml:"server"`
    
    // Database es la configuración de la base de datos
    Database DatabaseConfig `yaml:"database"`
    
    // LogLevel puede ser: debug, info, warn, error
    LogLevel string `yaml:"log_level"`
}

// ServerConfig describe los parámetros del servidor
type ServerConfig struct {
    // Host es la dirección de bind (ej: 0.0.0.0, localhost)
    Host string `yaml:"host"`
    
    // Port es el puerto TCP (1-65535)
    Port int `yaml:"port"`
}
```

### 48.11.4 - Testing de Configuración

```go
package config

import (
    "os"
    "testing"
)

func TestLoadConfigFromEnv(t *testing.T) {
    // Setup
    os.Setenv("SERVER_HOST", "0.0.0.0")
    os.Setenv("SERVER_PORT", "3000")
    defer func() {
        os.Unsetenv("SERVER_HOST")
        os.Unsetenv("SERVER_PORT")
    }()
    
    cfg := Load()
    
    // Assertions
    if cfg.Server.Host != "0.0.0.0" {
        t.Errorf("Expected host 0.0.0.0, got %s", cfg.Server.Host)
    }
    if cfg.Server.Port != 3000 {
        t.Errorf("Expected port 3000, got %d", cfg.Server.Port)
    }
}

func TestConfigValidation(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: &Config{
                Server: ServerConfig{Host: "localhost", Port: 8080},
            },
            wantErr: false,
        },
        {
            name: "invalid port",
            config: &Config{
                Server: ServerConfig{Host: "localhost", Port: 99999},
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Unexpected error: %v", err)
            }
        })
    }
}
```

### 48.11.5 - Comparación con Otros Lenguajes

#### Python - python-decouple

```python
# Python
from decouple import config, Csv

DATABASE_URL = config('DATABASE_URL', default='postgresql://localhost')
DEBUG = config('DEBUG', default=False, cast=bool)
ALLOWED_HOSTS = config('ALLOWED_HOSTS', cast=Csv())

# Go equivalente
dbURL := os.Getenv("DATABASE_URL")
if dbURL == "" {
    dbURL = "postgresql://localhost"
}
```

#### Java - Spring Boot

```java
# application.properties
server.host=localhost
server.port=8080
spring.datasource.url=${DATABASE_URL}

# Go equivalente con Viper
viper.SetDefault("server.host", "localhost")
viper.SetDefault("server.port", 8080)
viper.BindEnv("spring.datasource.url", "DATABASE_URL")
```

#### Node.js - dotenv

```javascript
// Node.js
require('dotenv').config();
const dbUrl = process.env.DATABASE_URL;

// Go equivalente
godotenv.Load()
dbURL := os.Getenv("DATABASE_URL")
```

---

## Ejercicios Progresivos

### Ejercicio 1: Configuración Básica desde Variables de Entorno

**Objetivo:** Leer configuración básica del ambiente

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

// TODO: Implementar
// 1. Crear struct AppConfig con campos: Host, Port, Environment
// 2. Función LoadFromEnv() que lee:
//    - APP_HOST (default: localhost)
//    - APP_PORT (default: 8080)
//    - APP_ENV (default: development)
// 3. Función (c *AppConfig) Print() que muestra la configuración
// 4. main() que carga y imprime

// Solución esperada:
// $ APP_HOST=0.0.0.0 APP_PORT=3000 APP_ENV=staging ./app
// Output:
// Host: 0.0.0.0
// Port: 3000
// Environment: staging
```

### Ejercicio 2: Viper con YAML y Override de Ambiente

**Objetivo:** Cargar config desde YAML con override desde env vars

```go
package main

import (
    "fmt"
    "github.com/spf13/viper"
    "os"
)

// TODO: Implementar
// 1. Crear config.dev.yaml con estructura:
//    server:
//      host: localhost
//      port: 8080
//    database:
//      host: localhost
//      port: 5432
// 2. Usar Viper para cargar el archivo
// 3. Habilitar override desde env vars (SERVER_HOST, DB_PORT, etc)
// 4. Unmarshal a struct y imprimir

// Solución esperada:
// $ SERVER_HOST=prod.example.com go run main.go
// Output:
// Server Host: prod.example.com (from env, overrides file)
// Server Port: 8080 (from file)
// Database Host: localhost (from file)
// Database Port: 5432 (from file)
```

### Ejercicio 3: Secretos desde Archivo Protegido

**Objetivo:** Leer secretos de archivo .env.local sin cometerlo en git

```go
package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "os"
)

// TODO: Implementar
// 1. Crear .env con valores placeholder
// 2. Crear .env.local con valores reales (NO en git)
// 3. Implementar LoadSecrets() que:
//    - Carga .env primero
//    - Hace override con .env.local
// 4. Función GetDatabasePassword() que retorna password con validación
// 5. main() que intenta conectar a DB (simular)

// Solución esperada:
// $ git check-ignore .env.local  # Debe estar en .gitignore
// .env.local
// $ go run main.go
// Output:
// Database password loaded successfully
// Connecting to database...
```

### Ejercicio 4: Validación de Configuración al Startup

**Objetivo:** Validar configuración al iniciar y fallar gracefully

```go
package main

import (
    "fmt"
    "log"
    "net"
)

// TODO: Implementar
// 1. Struct Config con campos: Host, Port, Timeout, LogLevel
// 2. Método Validate() que verifica:
//    - Host no vacío y es hostname/IP válida
//    - Port entre 1-65535
//    - Timeout > 0
//    - LogLevel en [debug, info, warn, error]
// 3. main() que intenta cargar config inválida y falla
// 4. Mostrar error descriptivo

// Solución esperada:
// Config inválida:
// $ PORT=99999 go run main.go
// Output:
// Fatal: Configuration validation failed: Port 99999 out of valid range (1-65535)
```

### Ejercicio 5: Feature Flags Configurable

**Objetivo:** Sistema de feature flags que persiste en JSON

```go
package main

import (
    "encoding/json"
    "fmt"
)

// TODO: Implementar
// 1. Struct FeatureFlag con: Name, Enabled, Rollout, Owner
// 2. Struct FeatureFlagManager que:
//    - Carga flags desde flags.json
//    - Permite consultar si flag está enabled
//    - Permite habilitar/deshabilitar flags
//    - Persiste cambios en JSON
// 3. Método IsEnabledForUser(flagName, userID) que respeta rollout
// 4. CLI simple:
//    - get <flag> - muestra estado
//    - set <flag> true/false - activa/desactiva
//    - list - lista todos

// Solución esperada:
// $ go run main.go list
// Feature Flags:
// - new_dashboard: enabled=true, rollout=100%
// - beta_api: enabled=true, rollout=25%
//
// $ go run main.go set beta_api false
// Updated beta_api to disabled
//
// $ go run main.go get new_dashboard user123
// new_dashboard is ENABLED for user123
```

---

## Resumen del Capítulo

### Conceptos Clave

1. **Configuración en Go** - Múltiples fuentes: env vars, archivos, flags
2. **Jerarquía** - Prioridad: defaults < archivos < env < flags
3. **Viper** - Librería estándar para configuración compleja
4. **Secretos** - Nunca hardcodear, usar env vars o Vault
5. **Validación** - Fallar rápido al startup, no en runtime
6. **Hot Reload** - Recargar sin reiniciar (watchs de archivos)
7. **Feature Flags** - Habilitar features sin redeploy

### Comparación: Go vs Python vs Java

| Aspecto | Go | Python | Java |
|---------|-----|--------|------|
| Estándar | os.Getenv + Viper | python-decouple | Spring Boot |
| Secretos | Env vars + Vault | python-dotenv | Environment beans |
| Validación | Struct tags | pydantic | @Validated |
| Hot reload | fsnotify | watchdog | Spring Cloud Config |

### Anti-patrones a Evitar

❌ Hardcoded values  
❌ Secrets en código  
❌ Secrets en git  
❌ Sin validación  
❌ Valores inválidos al runtime  

### Best Practices

✅ Variables de entorno para configuración  
✅ Validación al startup  
✅ Logging sin secretos  
✅ .gitignore en .env.local  
✅ Feature flags para roll-outs  

---

**Fin del Capítulo 48: Configuration Management**

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/48-configuration-management/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/48-configuration-management):

```bash
cd examples/48-configuration-management
go run .
```
