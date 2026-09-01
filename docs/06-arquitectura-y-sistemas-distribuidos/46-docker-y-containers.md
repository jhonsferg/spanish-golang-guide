# Capítulo 46: Docker y containers

## Índice de Contenidos

- [46.1 ¿Qué es Docker? Conceptos Fundamentales](#461-qué-es-docker-conceptos-fundamentales)
- [46.2 Dockerfile Basics: Construyendo Imágenes](#462-dockerfile-basics-construyendo-imágenes)
- [46.3 Optimización de Imágenes Docker](#463-optimización-de-imágenes-docker)
- [46.4 Build and Push: Creando y Distribuyendo](#464-build-and-push-creando-y-distribuyendo)
- [46.5 Go Specific: Optimizaciones para Go](#465-go-specific-optimizaciones-para-go)
- [46.6 Docker Compose: Orquestación Multi-Servicio](#466-docker-compose-orquestación-multi-servicio)
- [46.7 Networking en Docker](#467-networking-en-docker)
- [46.8 Volumes y Persistencia de Datos](#468-volumes-y-persistencia-de-datos)
- [46.9 Health Checks: Liveness y Readiness](#469-health-checks-liveness-y-readiness)
- [46.10 Security en Containers](#4610-security-en-containers)
- [46.11 Buenas Prácticas y CI/CD](#4611-buenas-prácticas-y-cicd)
- [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 46.1 ¿Qué es Docker? Conceptos Fundamentales

### 46.1.1 Container vs Máquina Virtual

Docker introduce una nueva forma de empaquetar y distribuir aplicaciones mediante **containers**. A diferencia de las máquinas virtuales tradicionales que virtualizan el hardware completo, los containers comparten el kernel del host OS pero mantienen aislamiento completo a nivel de procesos, memoria y sistema de archivos.

```
Máquinas Virtuales:
┌─────────────────────────────────────────┐
│ Hypervisor (KVM, VirtualBox, Hyper-V)  │
├──────────────┬──────────────┬───────────┤
│ VM1          │ VM2          │ VM3       │
├──────────────┼──────────────┼───────────┤
│ OS           │ OS           │ OS        │
│ Kernel       │ Kernel       │ Kernel    │
│ Runtime      │ Runtime      │ Runtime   │
│ App          │ App          │ App       │
└──────────────┴──────────────┴───────────┘

Containers Docker:
┌──────────────────────────────────────┐
│ Host OS Kernel (Linux)               │
├──────┬──────┬──────┬──────┬──────┬───┤
│ Con1 │ Con2 │ Con3 │ Con4 │ Con5 │..│
├──────┼──────┼──────┼──────┼──────┼───┤
│ FS   │ FS   │ FS   │ FS   │ FS   │...│
│ Proc │ Proc │ Proc │ Proc │ Proc │...│
│ App  │ App  │ App  │ App  │ App  │...│
└──────┴──────┴──────┴──────┴──────┴───┘
```

**Ventajas de Containers:**
- **Ligereza**: Arrancan en milisegundos vs segundos en VMs
- **Eficiencia de Recursos**: Múltiples containers comparten kernel
- **Portabilidad**: "Build once, run anywhere"
- **Escalabilidad**: Fácil crear y destruir instancias
- **Aislamiento**: Procesos, redes, filesystems independientes
- **Reproducibilidad**: Misma imagen = Mismo comportamiento

### 46.1.2 Imágenes Docker

Una **imagen Docker** es una plantilla inmutable que contiene todo lo necesario para ejecutar una aplicación: código, runtime, libraries, variables de entorno y configuración.

**Características clave:**
- **Inmutables**: Una vez construidas, no cambian
- **Capas (Layers)**: Compuestas de capas apiladas read-only
- **Versionables**: Se identifican por tags (v1.0, latest, sha256:...)
- **Reutilizables**: Se pueden compartir en registries (Docker Hub, ECR, etc.)
- **Inheritables**: Se construyen sobre base images

### 46.1.3 Concepto de Layers

Las imágenes Docker se construyen mediante layers apilados. Cada instrucción en un Dockerfile (FROM, RUN, COPY, etc.) crea un nuevo layer.

```
Dockerfile:
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY app /app
RUN chmod +x /app/myapp
ENTRYPOINT ["/app/myapp"]

Layers:
┌──────────────────────────────────┐ Layer 5 (Top)
│ ENTRYPOINT ["/app/myapp"]        │ - Writable layer
├──────────────────────────────────┤
│ chmod +x /app/myapp              │ Layer 4
├──────────────────────────────────┤
│ COPY app /app                    │ Layer 3
├──────────────────────────────────┤
│ apk add ca-certificates          │ Layer 2
├──────────────────────────────────┤
│ FROM alpine:3.18                 │ Layer 1 (Base)
└──────────────────────────────────┘
```

**Ventajas del sistema de layers:**
- **Caching**: Si un layer no cambió, Docker lo reutiliza
- **Eficiencia de almacenamiento**: Layers se comparten entre imágenes
- **Builds rápidos**: Solo reconstruye layers modificados
- **Control de cambios**: Permite trackear qué cambió en cada nivel

### 46.1.4 Container vs Imagen

| Aspecto | Imagen | Container |
|--------|--------|-----------|
| **Tipo** | Plantilla | Instancia en ejecución |
| **Estado** | Inmutable | Mutable (durante ejecución) |
| **Almacenamiento** | En disco | En memoria (+ volúmenes persistentes) |
| **Cantidad** | Una imagen → muchos containers | Cada uno es independiente |
| **Ciclo de vida** | Build → Push → Pull | Create → Start → Run → Stop → Remove |

### 46.1.5 Arquitectura de Docker

```
┌────────────────────────────────────────────────────────┐
│ Docker Client (docker CLI)                             │
│  - docker run, build, push, pull, etc.                 │
└────────────────┬─────────────────────────────────────┘
                 │ REST API
                 ↓
┌────────────────────────────────────────────────────────┐
│ Docker Daemon (dockerd)                                │
│  - Gestiona imágenes, containers, volumes, networks    │
├────────────────────────────────────────────────────────┤
│ Containerd (runtime)                                   │
│  - Gestiona ciclo de vida del container                │
├────────────────────────────────────────────────────────┤
│ runc (OCI runtime)                                     │
│  - Ejecuta el container realmente                      │
├────────────────────────────────────────────────────────┤
│ Linux Kernel                                           │
│  - cgroups, namespaces, networking                     │
└────────────────────────────────────────────────────────┘
```

---

## 46.2 Dockerfile Basics: Construyendo Imágenes

### 46.2.1 Instrucciones Fundamentales

#### FROM - Base Image

```dockerfile
FROM golang:1.21-alpine
# o especificar versión exacta
FROM golang:1.21.5-alpine3.18@sha256:abc123...
```

La instrucción `FROM` especifica la imagen base. Es la primera instrucción válida (excepto ARG). Seleccionar una buena base image es crucial para el tamaño final.

#### RUN - Ejecutar Comandos

```dockerfile
# Forma de shell (corre en /bin/sh -c)
RUN apk add --no-cache ca-certificates

# Forma exec (preferida, no dispara shell)
RUN ["apk", "add", "--no-cache", "ca-certificates"]

# Multi-línea con &&
RUN apk add --no-cache \
    ca-certificates \
    curl \
    git
```

**Nota**: Usar forma exec cuando sea posible para evitar PID 1 shell wrapping.

#### COPY vs ADD

```dockerfile
# COPY - Preferido para archivos locales
COPY . /app
COPY main.go /app/
COPY config/ /app/config/

# ADD - También descomprime y puede descargar URLs
ADD https://example.com/file.tar.gz /tmp/
ADD myapp.tar.gz /app/
```

**Mejor práctica**: Usar COPY excepto cuando necesites descomprimir o descargar URLs.

#### WORKDIR - Directorio de Trabajo

```dockerfile
WORKDIR /app
# Todos los comandos subsecuentes se ejecutan aquí
COPY . .
RUN go build -o myapp .
```

#### CMD - Comando Predeterminado

```dockerfile
# Forma exec (preferida)
CMD ["/app/myapp"]

# Forma shell
CMD /app/myapp

# Forma con ENTRYPOINT
CMD ["--port", "8080"]
```

CMD es el comando por defecto, pero puede ser sobrescrito: `docker run myimage arg1 arg2`

#### ENTRYPOINT - Punto de Entrada

```dockerfile
# Exec form (preferido)
ENTRYPOINT ["/app/myapp"]

# Shell form
ENTRYPOINT /app/myapp
```

**ENTRYPOINT vs CMD:**
- ENTRYPOINT es el comando que siempre se ejecuta
- CMD son argumentos por defecto (pueden sobrescribirse)

```dockerfile
ENTRYPOINT ["./myapp"]
CMD ["--port", "8080"]

# docker run myimage              → ejecuta: ./myapp --port 8080
# docker run myimage --port 9000  → ejecuta: ./myapp --port 9000
```

#### ENV - Variables de Entorno

```dockerfile
ENV GO111MODULE=on
ENV API_KEY=default_value
ENV PORT=8080

# Acceso desde dentro del container:
# os.Getenv("PORT") → "8080"
```

#### EXPOSE - Puertos

```dockerfile
EXPOSE 8080
EXPOSE 5432
```

EXPOSE documenta qué puertos escucha la app, pero no publica automáticamente. Necesitas `-p` al correr: `docker run -p 8080:8080 myapp`

#### ARG - Argumentos de Build

```dockerfile
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-alpine

ARG BUILD_TIME
ARG GIT_COMMIT
RUN echo "Built at ${BUILD_TIME}, commit ${GIT_COMMIT}"
```

Se pasan en build: `docker build --build-arg GO_VERSION=1.22 .`

#### USER - Usuario del Container

```dockerfile
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser
```

### 46.2.2 Orden de Instrucciones Óptimo

```dockerfile
# 1. Base image - Raramente cambia
FROM golang:1.21-alpine AS builder

# 2. ARGs de build (pueden cambiar)
ARG BUILD_TIME
ARG GIT_COMMIT

# 3. Instalar dependencias del SO (estables)
RUN apk add --no-cache ca-certificates

# 4. Copiar go.mod y go.sum (menos frecuente que código)
COPY go.mod go.sum ./

# 5. Descargar dependencias Go (cacheado)
RUN go mod download

# 6. Copiar código fuente (cambia frecuentemente)
COPY . .

# 7. Build
RUN CGO_ENABLED=0 go build -o app .

# 8. Stage final
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/app /app
USER nobody
ENTRYPOINT ["/app"]
```

### 46.2.3 Ejemplo Completo Go Application

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /src

# Instalar dependencias del SO
RUN apk add --no-cache git make

# Copiar módulos primero (mejor caching)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código
COPY . .

# Build static binary
ARG VERSION=dev
ARG COMMIT=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o myapp .

# Stage final - minimal
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiar solo el binario
COPY --from=builder /src/myapp .

# Usuario no-root
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app
USER app

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./myapp"]
CMD ["--port", "8080"]
```

---

## 46.3 Optimización de Imágenes Docker

### 46.3.1 Problema: Tamaños Inmanejables

```bash
# Sin optimizar
$ docker images | grep myapp
myapp    latest    1.2GB    # ¡Gigabyte completo!

# Con optimización
myapp    optimized 15MB     # 80x más pequeño
```

### 46.3.2 Multi-Stage Builds

La técnica más poderosa: usar múltiples stages, compilar en uno, ejecutar desde otro.

```dockerfile
# Stage 1: Builder - Todo lo necesario para compilar
FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app .

# Stage 2: Runtime - Apenas lo mínimo para ejecutar
FROM alpine:3.18

RUN apk add --no-cache ca-certificates
WORKDIR /app

# Solo copiar el binario compilado, nada más
COPY --from=builder /build/app .

ENTRYPOINT ["./app"]
```

**Comparación de tamaños:**

```
Single stage Dockerfile:
- golang:1.21-alpine (330MB) +
- código fuente +
- dependencias compiladas +
- herramientas de build
= ~850MB

Multi-stage con alpine:3.18:
- builder stage (330MB) - descartado
- alpine:3.18 (7MB) +
- app binary (20MB)
= ~27MB

¡Reducción de 30x!
```

### 46.3.3 Base Images Alternativas

```dockerfile
# Opción 1: Alpine (recomendado)
FROM alpine:3.18
# Ventaja: 7MB
# Desventaja: Usa musl libc, no glibc

# Opción 2: Debian slim
FROM debian:bookworm-slim
# Ventaja: glibc (compatible con más binarios)
# Desventaja: 80MB

# Opción 3: Distroless
FROM gcr.io/distroless/base-debian12
# Ventaja: ~30MB, solo lo básico
# Desventaja: Difícil de debuggear (sin shell)

# Opción 4: Scratch (vacío)
FROM scratch
# Ventaja: 0MB de base image
# Desventaja: Requiere binary completamente estático
```

### 46.3.4 Scratch Image para Go

Go puede compilar binarios completamente estáticos. Para máxima miniaturización:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilación completamente estática
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o app .

# Imagen final completamente vacía
FROM scratch

COPY --from=builder /src/app /app

EXPOSE 8080
ENTRYPOINT ["/app"]
```

**Resultado**: ~20MB para toda la imagen (solo el binario + herramientas mínimas de runtime de scratch)

### 46.3.5 Estrategias de Optimización

#### 1. Minimizar layers

```dockerfile
# ❌ Malo - 3 layers
RUN apk add curl
RUN apk add git
RUN apk add make

# ✅ Bueno - 1 layer
RUN apk add --no-cache \
    curl \
    git \
    make
```

#### 2. Limpiar caches después de instalar

```dockerfile
# Alpine
RUN apk add --no-cache package && \
    rm -rf /var/cache/apk/*

# Debian
RUN apt-get update && apt-get install -y package && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
```

#### 3. Usar .dockerignore

```
# .dockerignore
.git
.gitignore
node_modules/
*.log
.DS_Store
README.md
.vscode/
coverage/
.env
secrets/
```

#### 4. Binarios optimizados

```bash
# Sin optimizar
CGO_ENABLED=0 go build -o app .
# Resultado: ~20MB

# Optimizado con -ldflags
CGO_ENABLED=0 go build -ldflags="-w -s" -o app .
# -w: remover tabla de símbolos (debug)
# -s: remover tabla de strings
# Resultado: ~6MB
```

#### 5. UPX - Comprimir Binarios

```dockerfile
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache upx

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o app .
RUN upx --best --lzma app

FROM scratch
COPY --from=builder /src/app /app
ENTRYPOINT ["/app"]
```

**Resultado**: ~2-3MB comprimido. Nota: UPX ralentiza startup, no usar para apps críticas.

### 46.3.6 Análisis de Capas

```bash
$ docker history myapp:latest
IMAGE           CREATED        CREATED BY                                      SIZE
abc123def456    2 min ago      /bin/sh -c #(nop) ENTRYPOINT ["/app"]          0B
def456ghi789    2 min ago      /bin/sh -c apk add --no-cache ca-cert...       5.2MB
ghi789jkl012    2 min ago      /bin/sh -c #(nop) WORKDIR /app                 0B
jkl012mno345    2 min ago      /bin/sh -c #(nop) COPY file:app /app           15MB
```

Herramienta útil: `dive`

```bash
$ dive myapp:latest
# Interfaz interactiva para explorar layers
```

---

## 46.4 Build and Push: Creando y Distribuyendo

### 46.4.1 Docker Build

```bash
# Build básico
docker build -t myapp:1.0 .

# Build con Dockerfile específico
docker build -f Dockerfile.prod -t myapp:1.0 .

# Build con argumentos
docker build \
  --build-arg VERSION=1.0 \
  --build-arg COMMIT=abc123 \
  -t myapp:1.0 .

# Build con tags múltiples
docker build \
  -t myapp:1.0 \
  -t myapp:latest \
  -t registry.example.com/myapp:1.0 \
  .

# Build de un repositorio remoto
docker build https://github.com/user/repo.git#main:path/to/dockerfile
```

### 46.4.2 Build Context y Performance

```bash
# Build con contexto específico
docker build -t myapp:1.0 /path/to/context

# Build sin contexto (devuelve de stdin)
docker build -t myapp:1.0 - < Dockerfile
```

**Build context**: Directorio que se envía al daemon docker. Minimizar el contexto acelera builds.

```dockerfile
# Dockerfile con build-kit mejorado
DOCKER_BUILDKIT=1 docker build -t myapp:1.0 .
```

### 46.4.3 Image Naming y Tagging

```
registry.example.com/namespace/repository:tag

Ejemplos:
docker.io/library/golang:1.21-alpine          # Docker Hub oficial
docker.io/myuser/myapp:1.0                    # Docker Hub personal
gcr.io/myproject/myapp:v1.0                   # Google Container Registry
ECR: 123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:latest

Tags comunes:
- latest     # Última versión (inestable, usarla con cuidado)
- v1.0       # Versión semántica
- v1.0-rc1   # Release candidate
- main       # Rama de desarrollo
- sha-abc123 # Git commit SHA
```

### 46.4.4 Docker Registry - Registros Locales

```bash
# Iniciar un registry local (v2)
docker run -d -p 5000:5000 --name registry registry:2

# Taggear para registry local
docker tag myapp:1.0 localhost:5000/myapp:1.0

# Push a registry local
docker push localhost:5000/myapp:1.0

# Pull desde registry local
docker pull localhost:5000/myapp:1.0

# Ver imágenes en registry local
curl http://localhost:5000/v2/_catalog
```

### 46.4.5 Push a Docker Hub

```bash
# Login
docker login
# Ingresa credenciales o token

# Taggear con tu usuario
docker tag myapp:1.0 myuser/myapp:1.0
docker tag myapp:1.0 myuser/myapp:latest

# Push
docker push myuser/myapp:1.0
docker push myuser/myapp:latest

# Pull
docker pull myuser/myapp:1.0
```

### 46.4.6 Push a Google Container Registry (GCR)

```bash
# Configurar credenciales
gcloud auth configure-docker gcr.io

# Taggear
docker tag myapp:1.0 gcr.io/myproject/myapp:1.0

# Push
docker push gcr.io/myproject/myapp:1.0

# Pull
docker pull gcr.io/myproject/myapp:1.0
```

### 46.4.7 Push a AWS ECR

```bash
# Login a ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin \
  123456789.dkr.ecr.us-east-1.amazonaws.com

# Taggear
docker tag myapp:1.0 \
  123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:1.0

# Push
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:1.0
```

### 46.4.8 Scanning de Imágenes

```bash
# Trivy - Scanner opensource
trivy image myapp:1.0

# Docker Scout (integrado en algunas versiones)
docker scout cves myapp:1.0

# Salida:
# ✗ CRITICAL: CVE-2023-1234
#   alpine-base: 1.0-r10
#   Fix: 1.0-r11

# Filtrar solo críticas
trivy image --severity CRITICAL myapp:1.0
```

---

## 46.5 Go Specific: Optimizaciones para Go

### 46.5.1 Static Binaries en Go

Go puede compilar binarios completamente estáticos que no dependen de librerías del sistema:

```bash
# Binario dinámico (default)
go build -o app .
ldd ./app
# linux-vdso.so.1 => (0x...)
# libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x...)

# Binario estático
CGO_ENABLED=0 go build -o app .
ldd ./app
# not a dynamic executable

# Comprobación de dependencias
file ./app
# ELF 64-bit LSB executable, x86-64, statically linked
```

**Ventaja**: Funciona en cualquier Linux, sin importar las librerías instaladas.

### 46.5.2 Dockerfile Óptimo para Go

```dockerfile
FROM golang:1.21-alpine AS builder

# Instalar ca-certificates en builder (necesario para HTTPS)
RUN apk add --no-cache ca-certificates tzdata git

WORKDIR /src

# Copiar módulos para mejor caching
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Variables de build comunes
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Compilación estática y optimizada
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s \
      -X main.Version=${VERSION} \
      -X main.Commit=${COMMIT} \
      -X main.BuildTime=${BUILD_TIME}" \
    -o app .

# Stage 2: Runtime mínimo
FROM scratch

# Copiar CA certificates y timezone data desde builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar binario
COPY --from=builder /src/app /

EXPOSE 8080
ENTRYPOINT ["/app"]
```

**Flags explicados:**
- `CGO_ENABLED=0`: Deshabilitar CGO para binario estático
- `GOOS=linux GOARCH=amd64`: Compilación cruzada si es necesario
- `-trimpath`: Remover rutas absolutas (reproducibilidad)
- `-w -s`: Remover debug symbols y string tables
- `-X`: Inyectar variables de build (versión, commit, etc.)

### 46.5.3 Comparación Go vs Java vs Python

```dockerfile
# Go - Multistage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache ca-certificates
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app .

FROM scratch
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app /
ENTRYPOINT ["/app"]

# Tamaño final: ~15MB (solo binario + certs)

---

# Java - Multistage
FROM maven:3.9-eclipse-temurin-17 AS builder
COPY . .
RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-jre-alpine
COPY --from=builder /target/app.jar /app.jar
ENTRYPOINT ["java", "-jar", "/app.jar"]

# Tamaño final: ~500MB (JRE + JAR)

---

# Python - Multistage
FROM python:3.11-alpine AS builder
RUN apk add --no-cache gcc musl-dev
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

FROM python:3.11-alpine
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY . /app
WORKDIR /app
ENTRYPOINT ["python", "app.py"]

# Tamaño final: ~200MB (Python interpreter + packages)
```

**Ventajas de Go:**
- Binarios estáticos (scratch image posible)
- Tamaño final 20-30x más pequeño que Java/Python
- Arranque instantáneo
- Sin overhead de runtime

### 46.5.4 CGO Considerations

```go
// main.go - Código con CGO
package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Println(strings.ToUpper("hello"))
}
```

```bash
# Compilación normal con CGO
go build -o app .
# Funciona

# Sin CGO
CGO_ENABLED=0 go build -o app .
# También funciona (CGO solo necesario para C interop)

# Con librerias C (e.j., sqlite3)
import "github.com/mattn/go-sqlite3"

# CGO required:
go build -o app .                  # ✅ Funciona
CGO_ENABLED=0 go build -o app .   # ❌ Falla - requiere CGO
```

### 46.5.5 Go-Specific Patterns

```dockerfile
# Pattern 1: Copiar go.mod y go.sum primero
# Esto aprovecha el caching de layers
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./           # Layer 1 - Raramente cambia
RUN go mod download             # Layer 2 - Descarga dependencias

COPY . .                        # Layer 3 - Código (cambia frecuentemente)
RUN go build -o app .          # Layer 4 - Build

# Si cambias código, solo se reconstruyen layers 3 y 4
# Si cambias dependencias, todos se reconstruyen
# Si solo cambias Dockerfile format, se reutiliza el cache

# Pattern 2: Test en build
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY . .
RUN go test ./...              # Fallar en test = fallar en build
RUN go build -o app .

# Pattern 3: Live Reload Development
# Durante desarrollo, usar volume mounts
# docker run -v $(pwd):/src myapp-dev
# Dentro del container, usar "air" o "fswatch" para reload
```

---

## 46.6 Docker Compose: Orquestación Multi-Servicio

### 46.6.1 Archivo docker-compose.yml Básico

```yaml
version: '3.9'

services:
  web:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:password@db:5432/myapp
      - REDIS_URL=redis://cache:6379
    depends_on:
      - db
      - cache
    networks:
      - app-network

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=myapp
    volumes:
      - db-data:/var/lib/postgresql/data
    networks:
      - app-network

  cache:
    image: redis:7-alpine
    networks:
      - app-network

volumes:
  db-data:

networks:
  app-network:
    driver: bridge
```

### 46.6.2 Build en Compose

```yaml
services:
  web:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        VERSION: v1.0
        COMMIT: abc123

  # O simplemente:
  # build: .
```

### 46.6.3 Environment Variables

```yaml
services:
  web:
    # Inline
    environment:
      PORT: "8080"
      LOG_LEVEL: debug

    # Desde archivo .env
    env_file:
      - .env
      - .env.local

    # Desde archivo y override inline
    env_file:
      - .env
    environment:
      LOG_LEVEL: debug  # Sobrescribe .env si existe
```

**Archivo .env:**
```
DATABASE_URL=postgres://user:pass@localhost/db
API_KEY=secret123
DEBUG=true
```

### 46.6.4 Volúmenes en Compose

```yaml
services:
  web:
    volumes:
      # Named volume
      - app-data:/data

      # Bind mount
      - ./config:/app/config

      # Bind mount read-only
      - ./certs:/app/certs:ro

  db:
    volumes:
      - db-data:/var/lib/postgresql/data

volumes:
  # Named volumes deben declararse aquí
  app-data:
    driver: local

  db-data:
    driver: local
```

### 46.6.5 Networking en Compose

```yaml
version: '3.9'

services:
  web:
    # Accesible por nombre: "web:8080" desde otros servicios
    networks:
      - app-network
    ports:
      - "8080:8080"

  api:
    networks:
      - app-network
    # Accesible como: http://web:8080 (DNS automático)

  db:
    networks:
      - backend
    # No accesible desde "web" ya que está en otra red

networks:
  app-network:
    driver: bridge

  backend:
    driver: bridge
```

**DNS automático:**
```
En compose, los servicios se resuelven por nombre dentro de la misma red:
- web → Resuelve a IP de web
- db → Resuelve a IP de db
- web:8080 → Acceso directo
```

### 46.6.6 Comandos Docker Compose

```bash
# Iniciar servicios
docker-compose up                    # Foreground
docker-compose up -d                 # Detached
docker-compose up --build            # Rebuild antes de iniciar

# Parar servicios
docker-compose down                  # Parar y remover containers
docker-compose down -v               # Parar y remover con volúmenes
docker-compose stop                  # Parar sin remover
docker-compose start                 # Reactivar parados

# Logs
docker-compose logs                  # Todos los servicios
docker-compose logs -f               # Follow (tail)
docker-compose logs web              # Solo web
docker-compose logs --tail=50 db     # Últimas 50 líneas

# Ejecución
docker-compose exec web /bin/sh      # Ejecutar en container vivo
docker-compose run web go test ./... # Ejecutar comando una vez

# Info
docker-compose ps                    # Ver estado de servicios
docker-compose config                # Validar y ver config final
```

### 46.6.7 Ejemplo Completo: API + DB + Cache

```yaml
version: '3.9'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        VERSION: v1.0
    container_name: myapi
    ports:
      - "8080:8080"
    environment:
      PORT: "8080"
      DATABASE_URL: postgres://app:apppass@db:5432/myapp
      REDIS_URL: redis://cache:6379
      LOG_LEVEL: info
    depends_on:
      db:
        condition: service_healthy
      cache:
        condition: service_started
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 3s
      retries: 3

  db:
    image: postgres:15-alpine
    container_name: mydb
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: apppass
      POSTGRES_DB: myapp
    volumes:
      - db-data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  cache:
    image: redis:7-alpine
    container_name: mycache
    command: redis-server --appendonly yes
    volumes:
      - cache-data:/data
    networks:
      - app-network
    restart: unless-stopped

volumes:
  db-data:
  cache-data:

networks:
  app-network:
    driver: bridge
```

**Lanzar:**
```bash
docker-compose up -d
docker-compose logs -f api
docker-compose exec db psql -U app -d myapp -c "SELECT COUNT(*) FROM users;"
```

---

## 46.7 Networking en Docker

### 46.7.1 Tipos de Networks

```dockerfile
# 1. Bridge (default)
docker network create bridge-net
docker run --network bridge-net myapp

# 2. Host (comparte namespace de red del host)
docker run --network host myapp
# localhost:8080 = directamente 8080 en host

# 3. Overlay (para swarm/Kubernetes)
docker network create -d overlay swarm-net

# 4. None (sin red)
docker run --network none myapp
```

### 46.7.2 Bridge Network (Default)

```bash
# Crear red bridge personalizada
docker network create app-bridge

# Correr containers en la red
docker run --name web --network app-bridge -d myapp:latest
docker run --name db --network app-bridge -d postgres:15

# Dentro del container web, accede a db:
# curl http://db:5432
# → DNS resolve "db" a la IP del container db
```

**Beneficios vs bridge default:**
- DNS resolve por nombre de container
- Containers pueden conectarse/desconectarse dinámicamente
- Aislamiento entre redes

### 46.7.3 Host Network

```bash
# Comparte namespace de red del host
docker run --network host -p 8080:8080 myapp:latest

# ⚠️ WARNING: Sin aislamiento de red
# - container puede acceder a todos los puertos del host
# - requiere --cap-add para algunas operaciones
# - usado en casos específicos (monitoring, reverse proxy)
```

### 46.7.4 Port Mapping

```bash
# Mapear puerto específico
docker run -p 8080:8080 myapp

# Mapear múltiples puertos
docker run -p 8080:8080 -p 9090:9090 myapp

# Mapear a puerto arbitrario (ephemeral)
docker run -p 8080 myapp
# Acceso en puerto aleatorio, e.j. 32768

# Mapear solo localhost
docker run -p 127.0.0.1:8080:8080 myapp
# No accesible desde red

# Mapear todos los interfaces
docker run -p 0.0.0.0:8080:8080 myapp
# Accesible desde cualquier interfaz
```

### 46.7.5 Container Linking (Legacy - No Usar)

```bash
# DEPRECADO - No usar en código nuevo
docker run --name db -d postgres:15
docker run --link db:database -d myapp

# Usa networks en su lugar:
docker network create app
docker run --name db --network app -d postgres:15
docker run --name web --network app -d myapp
```

### 46.7.6 DNS en Docker

```yaml
# docker-compose.yml
services:
  web:
    networks:
      - app
    # Accesible como:
    # - web (resolución interna)
    # - web.app (FQDN)

  api:
    networks:
      - app
    # Desde aquí: curl http://web/
    #             curl http://web.app/

  # Custom DNS
  db:
    networks:
      - app
    dns:
      - 8.8.8.8
      - 8.8.4.4
```

### 46.7.7 Network Troubleshooting

```bash
# Ver redes
docker network ls
docker network inspect app-bridge

# Conectar container a red existente
docker network connect app-bridge running-container

# Desconectar
docker network disconnect app-bridge container-name

# Debuggear desde dentro de container
docker exec -it myapp /bin/sh

# Dentro del container:
ping db          # Verificar resolución DNS
curl http://db   # Verificar conectividad
netstat -tlnp    # Ver puertos abiertos
```

---

## 46.8 Volumes y Persistencia de Datos

### 46.8.1 Tipos de Mounts

```
┌─────────────────────┐
│ Named Volumes       │ Administrados por Docker
│ (/var/lib/docker/)  │ Ideal para datos persistentes
└─────────────────────┘

┌─────────────────────┐
│ Bind Mounts         │ Directorios del host
│ (./data)            │ Ideal para desarrollo
└─────────────────────┘

┌─────────────────────┐
│ tmpfs Mounts        │ En memoria
│ (--tmpfs /tmp)      │ Ideal para datos temporales
└─────────────────────┘
```

### 46.8.2 Named Volumes

```bash
# Crear volume
docker volume create db-data

# Usar en container
docker run -v db-data:/var/lib/postgresql/data postgres:15

# Inspeccionar
docker volume inspect db-data
# {
#     "Name": "db-data",
#     "Mountpoint": "/var/lib/docker/volumes/db-data/_data",
#     "Driver": "local"
# }

# Listar volúmenes
docker volume ls

# Eliminar
docker volume rm db-data

# Limpiar volúmenes sin usar
docker volume prune

# Dentro de docker-compose.yml
volumes:
  db-data:          # Volume named
    driver: local
```

### 46.8.3 Bind Mounts

```bash
# Desarrollo - Código local
docker run -v $(pwd)/src:/app/src myapp:latest

# Config files
docker run -v $(pwd)/config.yml:/etc/app/config.yml:ro postgres:15

# Read-only
docker run -v $(pwd)/certs:/app/certs:ro myapp:latest

# Absoluto vs relativo
docker run -v /absolute/path:/app/data postgres
docker run -v ./relative/path:/app/data postgres
```

**En docker-compose:**
```yaml
services:
  web:
    volumes:
      - ./src:/app/src              # Bind mount RW
      - ./config:/app/config:ro     # Bind mount RO
      - app-data:/data              # Named volume

volumes:
  app-data:
```

### 46.8.4 tmpfs Mounts

```bash
# Almacenamiento temporal en memoria
docker run --tmpfs /tmp myapp
docker run --tmpfs /run:rw,size=64M myapp

# En docker-compose:
services:
  app:
    tmpfs:
      - /tmp
      - /run:size=64M
```

**Uso:**
- Datos temporales (sesiones, cache)
- Archivos sensibles que no deben persistir
- Mejor performance que disco

### 46.8.5 Persistencia en Docker Compose

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:15
    volumes:
      - postgres-data:/var/lib/postgresql/data
    # Los datos persisten en postgres-data volume
    # Incluso si ejecutas docker-compose down

  api:
    build: .
    volumes:
      - ./logs:/app/logs    # Bind mount para ver logs
      - api-cache:/tmp/cache  # Named volume para cache

  redis:
    image: redis:7
    volumes:
      - redis-data:/data    # Datos persistentes
    command: redis-server --appendonly yes  # AOF enabled

volumes:
  postgres-data:
  api-cache:
  redis-data:
```

**Ciclo de vida:**
```bash
docker-compose up -d       # Crea containers y volumes
# Insertas datos en la DB

docker-compose down        # Parar containers
# Volúmenes persisten!

docker-compose up -d       # Reiniciar
# Datos siguen ahí
```

### 46.8.6 Backup y Restore de Volumes

```bash
# Backup
docker run --rm -v db-data:/data -v $(pwd):/backup \
  busybox tar czf /backup/db-backup.tar.gz -C /data .

# Restore
docker volume create db-data-restored
docker run --rm -v db-data-restored:/data -v $(pwd):/backup \
  busybox tar xzf /backup/db-backup.tar.gz -C /data

# Usando container como intermediario
docker exec container-name pg_dump -U user db > backup.sql
docker exec -i container-name psql -U user db < backup.sql
```

---

## 46.9 Health Checks: Liveness y Readiness

### 46.9.1 HEALTHCHECK en Dockerfile

```dockerfile
# Forma simple
HEALTHCHECK CMD curl -f http://localhost:8080/health || exit 1

# Con parámetros
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Explicación de parámetros:
# --interval=10s     → Ejecutar cada 10 segundos
# --timeout=3s       → Timeout de 3 segundos
# --start-period=5s  → Esperar 5 segundos antes de empezar a revisar
# --retries=3        → Después de 3 fallos = unhealthy

# Usando script
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s \
  CMD /bin/bash -c 'exec 3<>/dev/tcp/localhost/8080; printf "GET /health HTTP/1.1\r\nhost: localhost\r\n\r\n" >&3; timeout 2 cat <&3 | grep -q "200 OK"; exit $?'
```

### 46.9.2 Implementar Healthcheck en Go

```go
package main

import (
    "fmt"
    "net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
    // Verificar dependencias críticas
    if !isDBHealthy() {
        http.Error(w, "DB unhealthy", http.StatusServiceUnavailable)
        return
    }
    
    if !isCacheHealthy() {
        http.Error(w, "Cache unhealthy", http.StatusServiceUnavailable)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    // Similar pero más laxo - solo si está listo recibir tráfico
    if !canServeRequests() {
        http.Error(w, "Not ready", http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func main() {
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/ready", readinessHandler)
    http.ListenAndServe(":8080", nil)
}
```

### 46.9.3 Health Checks en Docker Compose

```yaml
services:
  api:
    build: .
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 3s
      start_period: 10s
      retries: 3

  db:
    image: postgres:15
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app -d myapp"]
      interval: 10s
      timeout: 5s
      start_period: 10s
      retries: 5

  redis:
    image: redis:7
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 3
```

### 46.9.4 Estados de Health

```bash
# Estados posibles:
# - starting    → Container acaba de arrancar
# - healthy     → Últimas pruebas exitosas
# - unhealthy   → Falló múltiples veces
# - none        → Sin healthcheck configurado

# Ver estado:
docker ps
# CONTAINER ID   IMAGE       STATUS
# abc123         myapp       Up 2 minutes (healthy)
# def456         mydb        Up 2 minutes (unhealthy)

# Ver logs de health:
docker inspect --format='{{json .State.Health}}' abc123 | jq
# {
#   "Status": "healthy",
#   "FailingStreak": 0,
#   "Log": [{...}]
# }
```

### 46.9.5 Liveness vs Readiness

```
LIVENESS CHECK:
- ¿El container aún vive?
- Si falla repetidamente → Kubernetes reinicia
- Ejemplo: ¿Puedo conectar a DB?

READINESS CHECK:
- ¿Está listo para recibir tráfico?
- Si falla → Kubernetes lo quita de load balancer
- Ejemplo: ¿He terminado la inicialización?

En Dockerfile:
HEALTHCHECK → Liveness check

En Kubernetes:
livenessProbe   → Restart si falla
readinessProbe  → Remove del service si falla
startupProbe    → Esperar antes de liveness
```

---

## 46.10 Security en Containers

### 46.10.1 Running as Non-Root

```dockerfile
# ❌ BAD - Running as root
FROM alpine:3.18
COPY app /
ENTRYPOINT ["/app"]

# ✅ GOOD - Create non-root user
FROM alpine:3.18

# Crear usuario
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

COPY --chown=appuser:appuser app /

USER appuser

ENTRYPOINT ["/app"]
```

**Verificar:**
```bash
docker run --rm myapp:latest id
# uid=1000(appuser) gid=1000(appuser) groups=1000(appuser)
```

### 46.10.2 Read-Only Root Filesystem

```dockerfile
# Dockerfile
FROM alpine:3.18
RUN mkdir -p /tmp /app
COPY app /app
USER appuser
ENTRYPOINT ["/app"]

# Docker run
docker run --read-only -v /tmp myapp:latest
# Solo /tmp es writable, rest es read-only
```

**En docker-compose:**
```yaml
services:
  app:
    read_only: true
    volumes:
      - /tmp   # Excepción para archivos temporales
```

### 46.10.3 Secrets Management

```dockerfile
# ❌ NO - Hardcoded secrets
ENV DB_PASSWORD=super_secret

# ✅ SI - Pasado en runtime
# docker build -t myapp .
# docker run -e DB_PASSWORD=secret myapp
```

**Mejor: Usar Docker Secrets (Swarm)**
```bash
echo "super_secret" | docker secret create db_password -

docker service create \
  --secret db_password \
  myapp:latest
```

**En aplicación Go:**
```go
package main

import (
    "io/ioutil"
    "os"
    "strings"
)

func getDBPassword() string {
    // Desde variable de entorno
    if pwd := os.Getenv("DB_PASSWORD"); pwd != "" {
        return pwd
    }

    // O desde archivo de secret
    if data, err := ioutil.ReadFile("/run/secrets/db_password"); err == nil {
        return strings.TrimSpace(string(data))
    }

    panic("DB_PASSWORD not set")
}
```

### 46.10.4 Image Scanning

```bash
# Trivy (recomendado)
trivy image myapp:latest
# Detecta CVEs en librerías

# Docker Scout
docker scout cves myapp:latest

# Snyk
snyk container test myapp:latest

# Aqua
aqua scan myapp:latest
```

**Ejemplo output:**
```
myapp:latest

Vulnerabilities
Total: 5

┌─────────────────────┬──────────┬──────────┐
│ Library             │ Severity │ CVE ID   │
├─────────────────────┼──────────┼──────────┤
│ openssl             │ CRITICAL │ CVE-2023 │
│ curl                │ HIGH     │ CVE-2024 │
└─────────────────────┴──────────┴──────────┘
```

### 46.10.5 Capability Dropping

```dockerfile
# ❌ Sin control
FROM alpine:3.18
COPY app /
ENTRYPOINT ["/app"]

# ✅ Con capabilities mínimas
FROM alpine:3.18
COPY app /
USER appuser
ENTRYPOINT ["/app"]

# docker run
docker run --cap-drop=ALL --cap-add=NET_BIND_SERVICE myapp
# Solo permite bind a puertos <1024
```

### 46.10.6 AppArmor y SELinux

```dockerfile
# Usar imagen base con AppArmor profile
FROM ubuntu:22.04

# Con AppArmor enabled en Docker:
docker run --security-opt apparmor=docker-default myapp

# Con SELinux:
docker run --security-opt label=type:svirt_apache_t myapp
```

### 46.10.7 Content Trust

```bash
# Firmar imágenes
export DOCKER_CONTENT_TRUST=1
docker push myuser/myapp:1.0
# Requiere keys criptográficas

# Verificar firma
docker pull --disable-content-trust=false myuser/myapp:1.0
```

---

## 46.11 Buenas Prácticas y CI/CD

### 46.11.1 Builds Reproducibles

```dockerfile
# Especificar tags exactos (no 'latest')
FROM golang:1.21.5-alpine3.18@sha256:abc123def456
FROM alpine:3.18@sha256:xyz789abc123

# Fijar versiones de paquetes
RUN apk add --no-cache \
    ca-certificates=20230506-r0 \
    curl=8.2.1-r0

# Usar digest exacto
RUN apk add --no-cache ca-certificates@20230506-r0
```

### 46.11.2 Versionado de Imágenes

```bash
# Semantic Versioning
docker tag myapp:latest myapp:v1.2.3
docker push myapp:v1.2.3
docker push myapp:v1.2    # Latest patch
docker push myapp:v1      # Latest minor
docker push myapp:latest  # Latest

# Git-based
docker build -t myapp:$(git rev-parse --short HEAD) .
docker push myapp:$(git rev-parse --short HEAD)

# Timestamp-based
docker build -t myapp:$(date +%Y%m%d-%H%M%S) .
```

### 46.11.3 CI/CD Pipeline (GitHub Actions)

```yaml
# .github/workflows/docker.yml
name: Docker Build and Push

on:
  push:
    branches: [main]
    tags: ['v*']

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          build-args: |
            VERSION=${{ github.ref_name }}
            COMMIT=${{ github.sha }}

      - name: Scan with Trivy
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
          format: 'sarif'
          output: 'trivy-results.sarif'

      - name: Upload Trivy results
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'
```

### 46.11.4 Pipeline GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - build
  - scan
  - push

variables:
  DOCKER_HOST: tcp://docker:2375
  REGISTRY: registry.gitlab.com
  IMAGE: $REGISTRY/$CI_PROJECT_PATH

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build
        --build-arg VERSION=$CI_COMMIT_TAG
        --build-arg COMMIT=$CI_COMMIT_SHA
        -t $IMAGE:$CI_COMMIT_SHA
        -t $IMAGE:latest
        .

scan:
  stage: scan
  image: aquasec/trivy:latest
  script:
    - trivy image --severity HIGH,CRITICAL $IMAGE:$CI_COMMIT_SHA

push:
  stage: push
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - echo $CI_REGISTRY_PASSWORD | docker login -u $CI_REGISTRY_USER --password-stdin $CI_REGISTRY
  script:
    - docker push $IMAGE:$CI_COMMIT_SHA
    - docker push $IMAGE:latest
  only:
    - tags
```

### 46.11.5 Buenas Prácticas - Resumen

| Práctica | ✅ Hacer | ❌ Evitar |
|----------|---------|----------|
| **Base Image** | alpine:3.18, debian:slim | ubuntu:latest, centos:latest |
| **User** | Crear usuario no-root | Ejecutar como root |
| **Secrets** | Variables de entorno, vault | Hardcoded en Dockerfile |
| **Layers** | Agrupar RUN, ordenar por frecuencia | Múltiples RUN innecesarios |
| **Caching** | COPY go.mod antes que código | COPY . como primer paso |
| **Tamaño** | Multi-stage, scratch, -ldflags="-w -s" | Todo en una imagen |
| **Scanning** | Trivy, Docker Scout | Ignorar vulnerabilidades |
| **Versionado** | Semver, digest exact | 'latest', tags flotantes |
| **Logs** | Stdout/stderr | Archivos en container |
| **Health** | HEALTHCHECK o probes K8s | Sin verificación |

### 46.11.6 Troubleshooting Común

```bash
# Image no construye
docker build --progress=plain -t myapp .
# Ver todos los pasos en detalle

# Container no arranca
docker logs container-name
docker run -it myapp:latest /bin/sh  # Shell interactivo

# Imagen muy grande
docker history myapp:latest
# Ver qué layer consume espacio

# Seguridad: Image contains secrets
trivy image myapp:latest
docker scout cves myapp:latest

# Performance: Builds lentos
DOCKER_BUILDKIT=1 docker build -t myapp .
# Usar BuildKit para mejores build cache
```

---

## Ejercicios Progresivos

### Ejercicio 1: Dockerfile Simple - Go App Básico

**Objetivo**: Crear un Dockerfile para una aplicación Go simple.

**Paso 1**: Crear estructura de proyecto

```bash
mkdir -p ejercicio1/cmd
cd ejercicio1
```

**Paso 2**: Crear archivo `main.go`

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from Container!\n")
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK\n")
    })

    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

**Paso 3**: Crear `Dockerfile`

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go build -o myapp .

EXPOSE 8080

ENTRYPOINT ["./myapp"]
```

**Paso 4**: Build y test

```bash
docker build -t myapp:1.0 .
docker run -p 8080:8080 myapp:1.0

# En otra terminal:
curl http://localhost:8080
curl http://localhost:8080/health
```

**Entregable**: Dockerfile que construya y ejecute exitosamente.

---

### Ejercicio 2: Multi-Stage Build - Optimizar Tamaño

**Objetivo**: Mejorar el Dockerfile anterior con multi-stage build.

**Paso 1**: Reemplazar Dockerfile

```dockerfile
# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o myapp .

# Stage 2: Runtime
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/myapp .

EXPOSE 8080

ENTRYPOINT ["./myapp"]
```

**Paso 2**: Comparar tamaños

```bash
# Dockerfile anterior (single-stage)
docker build -t myapp:single -f Dockerfile.single .
docker images | grep myapp

# Dockerfile nuevo (multi-stage)
docker build -t myapp:multi .
docker images | grep myapp

# Ver diferencia
docker images myapp:*
# myapp  multi      50MB
# myapp  single     850MB
```

**Entregable**: Multi-stage Dockerfile con imagen final <50MB.

---

### Ejercicio 3: Docker Compose - App + Database

**Objetivo**: Crear docker-compose.yml con API + PostgreSQL.

**Paso 1**: Crear estructura

```bash
mkdir -p ejercicio3
cd ejercicio3
mkdir -p cmd db
```

**Paso 2**: Crear aplicación Go mejorada (`main.go`)

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"
)

var db *sql.DB

func init() {
    // Conectar a PostgreSQL
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://app:apppass@db:5432/myapp?sslmode=disable"
    }

    var err error
    db, err = sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    }
}

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            fmt.Fprintf(w, "DB Error: %v\n", err)
            return
        }
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK\n")
    })

    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

**Paso 3**: Crear `go.mod`

```
module github.com/example/app

go 1.21

require github.com/lib/pq v1.10.9
```

**Paso 4**: Crear `Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
```

**Paso 5**: Crear `docker-compose.yml`

```yaml
version: '3.9'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://app:apppass@db:5432/myapp?sslmode=disable
    depends_on:
      db:
        condition: service_healthy
    networks:
      - app-network

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: apppass
      POSTGRES_DB: myapp
    volumes:
      - db-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

volumes:
  db-data:

networks:
  app-network:
```

**Paso 6**: Test

```bash
docker-compose up -d
docker-compose logs -f api
curl http://localhost:8080/health
docker-compose down -v
```

**Entregable**: Compose file funcional con API y DB conectadas.

---

### Ejercicio 4: Health Checks - Probes Configurados

**Objetivo**: Implementar health checks en Dockerfile y Compose.

**Paso 1**: Mejorar main.go con metrics

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

var startTime = time.Now()

func healthHandler(w http.ResponseWriter, r *http.Request) {
    uptime := time.Since(startTime).Seconds()
    
    // Verificar si está listo (después de 5 segundos)
    if uptime < 5 {
        w.WriteHeader(http.StatusServiceUnavailable)
        fmt.Fprintf(w, "Starting up\n")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status":"ok","uptime":%.1f}`, uptime)
}

func main() {
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "OK\n")
    })

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

**Paso 2**: Dockerfile con HEALTHCHECK

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates curl
COPY --from=builder /src/app /app

HEALTHCHECK --interval=5s --timeout=3s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080
ENTRYPOINT ["/app"]
```

**Paso 3**: Docker-compose con health checks

```yaml
version: '3.9'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 5s
      timeout: 3s
      start_period: 10s
      retries: 3
    networks:
      - app
    restart: on-failure

networks:
  app:
```

**Paso 4**: Observar health status

```bash
docker-compose up -d
sleep 15
docker-compose ps
# Ver estado: (starting), (healthy), o (unhealthy)

docker-compose logs api
docker inspect --format='{{json .State.Health}}' <container-id> | jq
```

**Entregable**: HEALTHCHECK implementado y funcionando correctamente.

---

### Ejercicio 5: CI/CD Pipeline - Automated Builds and Push

**Objetivo**: Crear pipeline CI/CD que construya y pushee imágenes automáticamente.

**Paso 1**: Configurar GitHub Actions

Crear archivo `.github/workflows/build.yml`:

```yaml
name: Build and Push Docker Image

on:
  push:
    branches:
      - main
    tags:
      - 'v*'
  pull_request:
    branches:
      - main

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Container Registry
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=sha,prefix={{branch}}-

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          build-args: |
            VERSION=${{ github.ref_name }}
            COMMIT=${{ github.sha }}
            BUILD_TIME=${{ github.event.head_commit.timestamp }}

      - name: Scan with Trivy
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
          format: 'sarif'
          output: 'trivy-results.sarif'
        if: github.event_name != 'pull_request'

      - name: Upload security results
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'
        if: github.event_name != 'pull_request'
```

**Paso 2**: Crear Dockerfile optimizado

```dockerfile
ARG GO_VERSION=1.21-alpine

FROM golang:${GO_VERSION} AS builder

WORKDIR /build

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s \
      -X main.Version=${VERSION} \
      -X main.Commit=${COMMIT} \
      -X main.BuildTime=${BUILD_TIME}" \
    -o app .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=builder /build/app .

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s \
  CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080

USER nobody

ENTRYPOINT ["./app"]
```

**Paso 3**: Versionar la aplicación en main.go

```go
package main

import (
    "fmt"
)

var (
    Version   = "dev"
    Commit    = "unknown"
    BuildTime = "unknown"
)

func init() {
    fmt.Printf("App Version: %s\n", Version)
    fmt.Printf("Commit: %s\n", Commit)
    fmt.Printf("Build Time: %s\n", BuildTime)
}
```

**Paso 4**: Push a repository

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/username/myapp.git
git push -u origin main

# Crear tag para release
git tag v1.0.0
git push origin v1.0.0
```

**Paso 5**: Verificar workflow

```bash
# En GitHub: Actions → Ver workflow en ejecución
# Verificar que:
# - Build completó exitosamente
# - Imagen fue pushed a ghcr.io
# - Scan de seguridad completó

# Pull imagen
docker pull ghcr.io/username/myapp:v1.0.0
docker run -p 8080:8080 ghcr.io/username/myapp:v1.0.0
```

**Entregable**: Pipeline CI/CD funcional que automatiza builds, scanning y push de imágenes Docker.

---

## Conclusión

Docker es fundamental en el desarrollo moderno de aplicaciones Go. Los containers proporcionan:
- **Consistencia**: Mismo comportamiento en dev, staging y prod
- **Escalabilidad**: Fácil horizontal scaling
- **Eficiencia**: Mejor utilización de recursos
- **Aislamiento**: Seguridad y estabilidad

Los próximos capítulos exploraremos Kubernetes, que orquesta containers a escala.

---

## Referencias Clave

- Docker Official Docs: https://docs.docker.com/
- Best Practices: https://docs.docker.com/develop/
- Go Docker: https://github.com/golang/go/wiki/Modules
- Trivy Security Scanner: https://github.com/aquasecurity/trivy
- Docker Compose Spec: https://github.com/compose-spec/compose-spec

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/46-docker-y-containers/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/46-docker-y-containers):

```bash
cd examples/46-docker-y-containers
go run .
```
