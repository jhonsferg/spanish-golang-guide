# Capítulo 47: Kubernetes básico - Orquestación de containers

## Introducción

Kubernetes (K8s) es la plataforma de orquestación de contenedores más adoptada en la industria. En este capítulo, exploraremos cómo Go se integra con Kubernetes para construir aplicaciones escalables, resilientes y autorreparables.

La orquestación de contenedores es esencial cuando tus aplicaciones Go se encuentran en múltiples contenedores que necesitan coordinarse, escalar y recuperarse automáticamente de fallos.

---

## 47.1 ¿Qué es Kubernetes?

### 47.1.1 Definición y Propósito

Kubernetes es un orquestador de contenedores de código abierto que automatiza el despliegue, escalado y operación de aplicaciones en contenedores. Proporciona:

- **Automatización**: Despliegues, actualizaciones y rollbacks automáticos
- **Escalado**: Ajuste horizontal y vertical de recursos
- **Autorrecuperación**: Reinicia contenedores fallidos, reemplaza nodos
- **Balanceo de carga**: Distribución automática de tráfico
- **Orquestación de almacenamiento**: Gestión de volúmenes persistentes
- **Secretos y configuración**: Manejo seguro de datos sensibles

### 47.1.2 Arquitectura General

```
┌─────────────────────────────────────────────────────┐
│          KUBERNETES CLUSTER                          │
│                                                      │
│  ┌──────────────────────────────────────────────┐  │
│  │         CONTROL PLANE (Master)                │  │
│  │                                               │  │
│  │ ┌─────────────┐ ┌─────────────┐              │  │
│  │ │ API Server  │ │ Controller  │              │  │
│  │ │             │ │ Manager     │              │  │
│  │ └─────────────┘ └─────────────┘              │  │
│  │                                               │  │
│  │ ┌─────────────┐ ┌─────────────┐              │  │
│  │ │ Scheduler   │ │ etcd        │              │  │
│  │ │             │ │ (database)  │              │  │
│  │ └─────────────┘ └─────────────┘              │  │
│  └──────────────────────────────────────────────┘  │
│                        │                            │
│  ┌────────────────────┼────────────────────────┐   │
│  │                    │                        │   │
│  ▼                    ▼                        ▼   │
│ ┌──────────┐      ┌──────────┐         ┌──────────┐│
│ │  Node 1  │      │  Node 2  │         │  Node N  ││
│ │          │      │          │         │          ││
│ │ ┌──────┐ │      │ ┌──────┐ │         │ ┌──────┐ ││
│ │ │ Pod  │ │      │ │ Pod  │ │   ...   │ │ Pod  │ ││
│ │ └──────┘ │      │ └──────┘ │         │ └──────┘ ││
│ │          │      │          │         │          ││
│ │ ┌──────┐ │      │ ┌──────┐ │         │          ││
│ │ │ Pod  │ │      │ │ Pod  │ │         │          ││
│ │ └──────┘ │      │ └──────┘ │         │          ││
│ └──────────┘      └──────────┘         └──────────┘│
│                                                    │
└─────────────────────────────────────────────────────┘
```

### 47.1.3 Comparación con Alternativas

| Característica | Kubernetes | Docker Swarm | Nomad |
|---|---|---|---|
| Complejidad | Alta | Baja | Media |
| Escalabilidad | Muy Alta (10,000+) | Media (10,000) | Alta |
| Ecosystem | Muy grande | Pequeño | Creciendo |
| Adoption | Dominante | Decreciente | Nicho |
| Curva aprendizaje | Pronunciada | Suave | Media |
| Multi-cloud | Nativo | Difícil | Excelente |
| Stateless workloads | Excelente | Bueno | Bueno |
| Stateful workloads | Muy bueno | Regular | Muy bueno |

### 47.1.4 Caso de Uso Real: Microservicios con Go

```go
// Servicio de usuarios
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/health", healthCheck)
    http.HandleFunc("/users", listUsers)

    log.Printf("Starting server on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, `{"status":"ok"}`)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, `[{"id":1,"name":"Alice"}]`)
}
```

En Kubernetes, múltiples instancias de este servicio pueden ejecutarse simultáneamente, escalando automáticamente según la demanda.

---

## 47.2 Conceptos Básicos

### 47.2.1 Pod: La Unidad Más Pequeña

Un Pod es la unidad más pequeña desplegable en Kubernetes. Aunque generalmente contiene un contenedor, puede contener múltiples contenedores que comparten:

- Dirección IP
- Almacenamiento
- Namespace de red

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-go-pod
  namespace: default
  labels:
    app: hello-go
spec:
  containers:
  - name: app
    image: myregistry/hello-go:1.0
    ports:
    - containerPort: 8080
    env:
    - name: LOG_LEVEL
      value: "INFO"
```

### 47.2.2 Node: Máquinas de Trabajo

Cada Node es una máquina (virtual o física) que ejecuta Pods. Contiene:

- **kubelet**: Agente que asegura que los containers se ejecuten
- **Container runtime**: Docker, containerd, CRI-O
- **kube-proxy**: Maneja networking

```bash
# Inspeccionar nodos
kubectl get nodes
kubectl describe node node-1
kubectl logs node-name -n kube-system
```

### 47.2.3 Services: Descubrimiento de Servicios

Los Services proporcionan una dirección IP estable y un nombre DNS para acceder a Pods que pueden cambiar:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-go-service
spec:
  selector:
    app: hello-go
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

### 47.2.4 Namespaces: Aislamiento Lógico

Los Namespaces permiten dividir un cluster en múltiples entornos virtuales:

```bash
# Crear namespace
kubectl create namespace production
kubectl create namespace development

# Listar recursos en un namespace
kubectl get pods -n production

# Cambiar namespace por defecto
kubectl config set-context --current --namespace=production
```

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: microservices
  labels:
    environment: production
```

### 47.2.5 Labels y Selectors

Los labels son pares clave-valor que permiten organizar y seleccionar recursos:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-pod-v1
  labels:
    app: myapp
    version: "1.0"
    tier: frontend
    environment: production
spec:
  containers:
  - name: app
    image: myapp:1.0
```

Selección de recursos:

```bash
# Igualdad
kubectl get pods -l app=myapp

# Conjunto de valores
kubectl get pods -l "version in (1.0, 2.0)"

# Existencia
kubectl get pods -l "tier"

# Negación
kubectl get pods -l "environment!=development"
```

### 47.2.6 Annotations

Las Annotations almacenan metadata sin restricciones de claves/valores:

```yaml
metadata:
  annotations:
    description: "Production API server"
    contacts: "devops@company.com"
    backup-schedule: "daily-2am"
    git-commit: "a1b2c3d4e5f6"
```

---

## 47.3 Deployments

### 47.3.1 ¿Qué es un Deployment?

Un Deployment es un controlador que asegura que un número especificado de Pods estén siempre corriendo. Proporciona:

- Replicación de Pods
- Actualizaciones rolling
- Rollbacks automáticos
- Escalado horizontal

```
┌────────────────────────┐
│    Deployment          │
│  (spec: replicas: 3)   │
├────────────────────────┤
│  ReplicaSet            │
│  (versión actual)      │
├────────────────────────┤
│ Pod  │ Pod  │ Pod      │
│ v1.0 │ v1.0 │ v1.0    │
└────────────────────────┘
```

### 47.3.2 Manifest de Deployment Básico

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-go
  namespace: default
  labels:
    app: hello-go
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1           # Máximo 1 pod extra durante actualización
      maxUnavailable: 1     # Máximo 1 pod no disponible
  selector:
    matchLabels:
      app: hello-go
  template:
    metadata:
      labels:
        app: hello-go
        version: "1.0"
    spec:
      containers:
      - name: hello-go
        image: myregistry/hello-go:1.0
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        env:
        - name: PORT
          value: "8080"
        - name: LOG_LEVEL
          value: "INFO"
```

### 47.3.3 Rolling Updates

```bash
# Actualizar imagen
kubectl set image deployment/hello-go \
  hello-go=myregistry/hello-go:2.0

# Ver historial de revisiones
kubectl rollout history deployment/hello-go

# Ver detalles de una revisión
kubectl rollout history deployment/hello-go --revision=2

# Rollback a revisión anterior
kubectl rollout undo deployment/hello-go

# Rollback a revisión específica
kubectl rollout undo deployment/hello-go --to-revision=1

# Pausar y reanudar rollout
kubectl rollout pause deployment/hello-go
kubectl rollout resume deployment/hello-go

# Esperar a que rollout se complete
kubectl rollout status deployment/hello-go -w
```

### 47.3.4 Escalado

```bash
# Escalar manualmente
kubectl scale deployment hello-go --replicas=5

# Escalado automático (Horizontal Pod Autoscaler)
kubectl autoscale deployment hello-go \
  --min=2 \
  --max=10 \
  --cpu-percent=80
```

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hello-go-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: hello-go
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### 47.3.5 Estrategias de Actualización

```yaml
# Blue-Green (dos versiones paralelas)
spec:
  strategy:
    type: Recreate  # Detiene todos, inicia nuevos

# Rolling Update (gradual)
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%

# Canary (con Flagger)
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
```

---

## 47.4 Services

### 47.4.1 Tipos de Services

#### ClusterIP (por defecto)

Expone el servicio solo dentro del cluster:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-go-internal
spec:
  type: ClusterIP
  selector:
    app: hello-go
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
```

```
┌─────────────────────────────────────┐
│         K8s Cluster                 │
│                                     │
│ Service: hello-go-internal:80       │
│ (IP fijo dentro del cluster)        │
│         ↓                           │
│ Pods con label app: hello-go        │
│ [Pod1] [Pod2] [Pod3]                │
└─────────────────────────────────────┘
```

#### NodePort

Expone el servicio en un puerto de cada Node:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-go-nodeport
spec:
  type: NodePort
  selector:
    app: hello-go
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080
```

Acceso: `<node-ip>:30080`

#### LoadBalancer

Crea un load balancer externo (requiere proveedor de cloud):

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-go-lb
spec:
  type: LoadBalancer
  selector:
    app: hello-go
  ports:
  - port: 80
    targetPort: 8080
```

Obtiene IP externa automáticamente.

### 47.4.2 Descubrimiento de Servicios

**DNS interno (automático):**

Desde un Pod, acceder a un servicio:

```bash
# Desde el mismo namespace
curl http://hello-go-service:80

# Desde otro namespace
curl http://hello-go-service.production.svc.cluster.local:80
```

**Inyección de variables de entorno:**

```bash
# Kubernetes automáticamente inyecta
HELLO_GO_SERVICE_HOST=10.0.0.1
HELLO_GO_SERVICE_PORT=80
HELLO_GO_SERVICE_PORT_HTTP=80
```

**En Go:**

```go
package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    // Kubernetes inyecta automáticamente
    serviceHost := os.Getenv("HELLO_GO_SERVICE_HOST")
    servicePort := os.Getenv("HELLO_GO_SERVICE_PORT")
    
    url := "http://" + serviceHost + ":" + servicePort
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    log.Printf("Response: %d", resp.StatusCode)
}
```

### 47.4.3 Endpoints y Service Discovery

```bash
# Ver endpoints detrás de un servicio
kubectl get endpoints hello-go-service
kubectl describe service hello-go-service
```

### 47.4.4 Headless Services

Para aplicaciones que necesitan comunicación peer-to-peer (StatefulSets):

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-go-headless
spec:
  clusterIP: None
  selector:
    app: hello-go
  ports:
  - port: 8080
    targetPort: 8080
```

---

## 47.5 ConfigMaps y Secrets

### 47.5.1 ConfigMaps: Configuración No Sensible

Almacena configuración en pares clave-valor:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
data:
  APP_NAME: "Hello Go API"
  LOG_LEVEL: "INFO"
  DATABASE_POOL_SIZE: "20"
  app.config: |
    [database]
    host=localhost
    port=5432
    timeout=30
```

Usar ConfigMap en Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-go
spec:
  template:
    spec:
      containers:
      - name: app
        image: hello-go:1.0
        envFrom:
        - configMapRef:
            name: app-config
        volumeMounts:
        - name: config-file
          mountPath: /etc/config
      volumes:
      - name: config-file
        configMap:
          name: app-config
          items:
          - key: app.config
            path: app.conf
```

### 47.5.2 Secrets: Datos Sensibles

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-credentials
  namespace: default
type: Opaque
data:
  username: YWRtaW4=          # base64: admin
  password: cGFzc3dvcmQEyMw== # base64: password123
stringData:
  connection-string: "postgresql://admin:password123@postgres:5432/db"
```

Usar Secret en Deployment:

```yaml
spec:
  containers:
  - name: app
    image: hello-go:1.0
    env:
    - name: DB_USER
      valueFrom:
        secretKeyRef:
          name: db-credentials
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-credentials
          key: password
```

### 47.5.3 Gestión de Secrets en Go

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
)

func main() {
    // Método 1: Variables de entorno
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    
    // Método 2: Leer desde volumen montado
    secretPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
    token, err := ioutil.ReadFile(secretPath)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Token length: %d", len(token))
    
    // NUNCA loguear secretos
    log.Printf("User: %s", dbUser)
}
```

### 47.5.4 Tipos de Secrets

```yaml
# Opaque (default)
type: Opaque

# Docker Registry
type: kubernetes.io/dockercfg
data:
  .dockercfg: <base64-encoded-docker-config>

# Service Account Token
type: kubernetes.io/service-account-token

# Basic Auth
type: kubernetes.io/basic-auth
data:
  username: <base64>
  password: <base64>

# SSH Auth
type: kubernetes.io/ssh-auth
data:
  ssh-privatekey: <base64>

# TLS
type: kubernetes.io/tls
data:
  tls.crt: <base64>
  tls.key: <base64>
```

---

## 47.6 StatefulSets

### 47.6.1 ¿Cuándo usar StatefulSets?

StatefulSets para aplicaciones que necesitan:

- Identidad de red estable
- Almacenamiento persistente asociado
- Escalado ordenado (0, 1, 2 → no 1, 3, 2)
- Terminación ordenada

### 47.6.2 Manifest de StatefulSet

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres-cluster
spec:
  serviceName: postgres
  replicas: 3
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:13
        ports:
        - name: postgres
          containerPort: 5432
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 10Gi
```

### 47.6.3 Características de Identidad

```bash
# Los pods tienen nombres predecibles
postgres-0
postgres-1
postgres-2

# Y DNS predecible
postgres-0.postgres.default.svc.cluster.local
postgres-1.postgres.default.svc.cluster.local
```

### 47.6.4 Actualización de StatefulSets

```bash
# OnDelete (manual)
kubectl patch statefulset postgres -p \
  '{"spec":{"updateStrategy":{"type":"OnDelete"}}}'

# RollingUpdate (automático)
kubectl patch statefulset postgres -p \
  '{"spec":{"updateStrategy":{"type":"RollingUpdate"}}}'
```

---

## 47.7 Persistent Volumes

### 47.7.1 Conceptos: PV, PVC, Storage Class

```
┌──────────────────────────────────────┐
│   StorageClass                       │
│   (tipo de almacenamiento)           │
└──────────────────────────────────────┘
           ↑
           │ (define)
           │
┌──────────────────────────────────────┐
│   PersistentVolume (PV)              │
│   (recurso de almacenamiento)        │
└──────────────────────────────────────┘
           ↑
           │ (vinculado a)
           │
┌──────────────────────────────────────┐
│   PersistentVolumeClaim (PVC)        │
│   (solicitud de almacenamiento)      │
└──────────────────────────────────────┘
           ↑
           │ (usado por)
           │
    Pod / Deployment
```

### 47.7.2 StorageClass

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3
  iops: "3000"
  throughput: "125"
  fstype: ext4
allowVolumeExpansion: true
volumeBindingMode: WaitForFirstConsumer
```

### 47.7.3 PersistentVolume y PersistentVolumeClaim

**Método estático (pre-provisionado):**

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-app-data
spec:
  capacity:
    storage: 100Gi
  accessModes:
    - ReadWriteOnce
  storageClassName: standard
  nfs:
    server: nfs-server.example.com
    path: "/shared/app-data"
```

**PersistentVolumeClaim:**

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: app-data-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: fast-ssd
  resources:
    requests:
      storage: 50Gi
```

### 47.7.4 Usar PVC en Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-storage
spec:
  containers:
  - name: app
    image: hello-go:1.0
    volumeMounts:
    - name: app-storage
      mountPath: /data
  volumes:
  - name: app-storage
    persistentVolumeClaim:
      claimName: app-data-pvc
```

### 47.7.5 Modos de Acceso

- **ReadWriteOnce (RWO)**: Un nodo, lectura y escritura
- **ReadOnlyMany (ROX)**: Múltiples nodos, solo lectura
- **ReadWriteMany (RWX)**: Múltiples nodos, lectura y escritura

---

## 47.8 Health Probes

### 47.8.1 Liveness Probe

Reinicia el contenedor si considera que está "muerto":

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-go-healthy
spec:
  template:
    spec:
      containers:
      - name: app
        image: hello-go:1.0
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
```

### 47.8.2 Readiness Probe

Determina si el Pod está listo para recibir tráfico:

```yaml
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 1
```

### 47.8.3 Startup Probe

Para aplicaciones que tardan en iniciarse:

```yaml
startupProbe:
  httpGet:
    path: /startup
    port: 8080
  failureThreshold: 30
  periodSeconds: 1
```

### 47.8.4 Implementación en Go

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync/atomic"
    "time"
)

var ready int32 = 0

func main() {
    http.HandleFunc("/health", healthCheck)
    http.HandleFunc("/ready", readyCheck)
    http.HandleFunc("/startup", startupCheck)
    
    // Simular inicialización
    go func() {
        time.Sleep(2 * time.Second)
        atomic.StoreInt32(&ready, 1)
    }()
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, `{"status":"alive"}`)
}

func readyCheck(w http.ResponseWriter, r *http.Request) {
    if atomic.LoadInt32(&ready) == 1 {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"ready"}`)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        fmt.Fprint(w, `{"status":"not ready"}`)
    }
}

func startupCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, `{"status":"starting"}`)
}
```

---

## 47.9 Resource Management

### 47.9.1 Requests y Limits

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-go
spec:
  template:
    spec:
      containers:
      - name: app
        image: hello-go:1.0
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

- **Requests**: Garantizado
- **Limits**: Máximo permitido

### 47.9.2 Clases de Calidad (QoS)

**Guaranteed** (máxima prioridad):
```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "100m"
```

**Burstable**:
```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

**BestEffort** (sin garantía):
```yaml
# Sin recursos especificados
```

### 47.9.3 Vertical Pod Autoscaler

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: hello-go-vpa
spec:
  targetRef:
    apiVersion: "apps/v1"
    kind: Deployment
    name: hello-go
  updatePolicy:
    updateMode: "Auto"
```

### 47.9.4 Cálculo de Recursos

```bash
# Ver uso actual de recursos
kubectl top nodes
kubectl top pods

# Ver requests y limits
kubectl describe node node-1
```

---

## 47.10 Ingress

### 47.10.1 ¿Qué es Ingress?

Ingress proporciona acceso HTTP/HTTPS desde fuera del cluster, con:
- Routing basado en hostname
- Routing basado en rutas
- Terminación TLS
- Virtual hosting

```
┌─────────────────────┐
│  Internet Traffic   │
└──────────┬──────────┘
           │
    ┌──────▼──────┐
    │   Ingress   │
    │ Controller  │
    └──────┬──────┘
           │
    ┌──────▼──────────────────────┐
    │     Ingress Rules            │
    │ ├─ api.example.com → svc1   │
    │ └─ web.example.com → svc2   │
    └──────────────────────────────┘
```

### 47.10.2 Manifest de Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-go-ingress
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - api.example.com
    secretName: api-tls-cert
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: hello-go-service
            port:
              number: 80
  - host: web.example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 8080
      - path: /static
        pathType: Prefix
        backend:
          service:
            name: static-service
            port:
              number: 80
```

### 47.10.3 Tipos de Ingress Controllers

| Controller | Proveedor | Características |
|---|---|---|
| NGINX | Comunidad | Flexible, ampliamente usado |
| Istio | Google/comunidad | Service mesh, mTLS |
| Ambassador | Ambassador Labs | API-driven |
| Traefik | Traefik | Moderno, dinámico |
| AWS ALB | AWS | Integrado con AWS |

### 47.10.4 Routing Avanzado

```yaml
# Routing con expresiones regulares
- path: /api/v[0-9]+
  pathType: ImplementationSpecific

# Múltiples servicios
- path: /users
  backend:
    service:
      name: users-service
      port:
        number: 8080
- path: /products
  backend:
    service:
      name: products-service
      port:
        number: 8080
```

---

## 47.11 Buenas Prácticas y Patterns

### 47.11.1 Organización con Namespaces

```bash
# Estructura recomendada
kube-system          # Componentes de K8s
kube-public          # Info pública
kube-node-lease      # Control de nodos
production           # Aplicaciones en producción
staging              # Aplicaciones en testing
development          # Desarrollo local
```

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    environment: prod
---
apiVersion: v1
kind: Namespace
metadata:
  name: development
  labels:
    environment: dev
```

### 47.11.2 Role-Based Access Control (RBAC)

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: default
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: pod-reader
subjects:
- kind: ServiceAccount
  name: app-sa
  namespace: default
```

### 47.11.3 Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-hello-go
spec:
  podSelector:
    matchLabels:
      app: hello-go
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          role: frontend
    ports:
    - protocol: TCP
      port: 8080
```

### 47.11.4 Resource Quotas

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-quota
  namespace: development
spec:
  hard:
    requests.cpu: "10"
    requests.memory: "20Gi"
    limits.cpu: "20"
    limits.memory: "40Gi"
    pods: "100"
---
apiVersion: v1
kind: LimitRange
metadata:
  name: limit-range
  namespace: development
spec:
  limits:
  - max:
      cpu: "2"
      memory: "4Gi"
    min:
      cpu: "50m"
      memory: "64Mi"
    type: Container
```

### 47.11.5 Antipatterns

❌ **Hardcoded Configuration**:
```go
// MAL
const dbHost = "postgres.prod.example.com"
const dbPassword = "secretpassword123"
```

✅ **Usar ConfigMaps y Secrets**:
```go
// BIEN
dbHost := os.Getenv("DB_HOST")
dbPassword := os.Getenv("DB_PASSWORD")
```

❌ **Sin Resource Limits**:
```yaml
# MAL
spec:
  containers:
  - name: app
    image: myapp:1.0
```

✅ **Con Resource Limits**:
```yaml
# BIEN
spec:
  containers:
  - name: app
    image: myapp:1.0
    resources:
      requests:
        memory: "256Mi"
        cpu: "100m"
      limits:
        memory: "512Mi"
        cpu: "500m"
```

❌ **Single Replica en Producción**:
```yaml
# MAL
spec:
  replicas: 1
```

✅ **Alta Disponibilidad**:
```yaml
# BIEN
spec:
  replicas: 3
  affinity:
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - myapp
        topologyKey: kubernetes.io/hostname
```

### 47.11.6 Monitoring y Observabilidad

```yaml
apiVersion: v1
kind: ServiceMonitor
metadata:
  name: hello-go-monitor
spec:
  selector:
    matchLabels:
      app: hello-go
  endpoints:
  - port: metrics
    interval: 30s
```

En Go, exponer métricas con Prometheus:

```go
package main

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "log"
    "net/http"
)

var (
    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path"},
    )
)

func init() {
    prometheus.MustRegister(requestCount)
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        requestCount.WithLabelValues(r.Method, r.URL.Path).Inc()
        w.WriteHeader(http.StatusOK)
    })
    
    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 47.12 Ejercicios Progresivos

### Ejercicio 1: Pod Simple - Desplegar Aplicación Go

**Objetivo**: Crear y desplegar un Pod básico con una aplicación Go.

**Pasos**:

1. Crear aplicación Go:

```go
// main.go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello from Kubernetes Pod!")
    })
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

2. Crear Dockerfile:

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /build
COPY main.go .
RUN go build -o app main.go

FROM scratch
COPY --from=builder /build/app /app
EXPOSE 8080
CMD ["/app"]
```

3. Crear manifiesto Pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-go-pod
spec:
  containers:
  - name: app
    image: myregistry/hello-go:1.0
    ports:
    - containerPort: 8080
```

4. Comandos:

```bash
# Aplicar manifiesto
kubectl apply -f pod.yaml

# Ver logs
kubectl logs hello-go-pod

# Port-forward
kubectl port-forward hello-go-pod 8080:8080

# Probar
curl localhost:8080
```

### Ejercicio 2: Deployment - Replicas y Rolling Update

**Objetivo**: Crear un Deployment con múltiples replicas y realizar una actualización.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-go
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: hello-go
  template:
    metadata:
      labels:
        app: hello-go
    spec:
      containers:
      - name: app
        image: myregistry/hello-go:1.0
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
```

**Verificación**:

```bash
# Ver replicas
kubectl get deployment hello-go
kubectl get pods -l app=hello-go

# Actualizar imagen
kubectl set image deployment/hello-go \
  app=myregistry/hello-go:2.0

# Ver progreso
kubectl rollout status deployment/hello-go -w

# Ver historial
kubectl rollout history deployment/hello-go

# Rollback si es necesario
kubectl rollout undo deployment/hello-go
```

### Ejercicio 3: Service - Exposición Interna y Externa

**Objetivo**: Crear Services de diferentes tipos para acceder a Pods.

```yaml
---
# ClusterIP (interno)
apiVersion: v1
kind: Service
metadata:
  name: hello-go-internal
spec:
  type: ClusterIP
  selector:
    app: hello-go
  ports:
  - port: 80
    targetPort: 8080
---
# NodePort (acceso externo por puerto alto)
apiVersion: v1
kind: Service
metadata:
  name: hello-go-nodeport
spec:
  type: NodePort
  selector:
    app: hello-go
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080
---
# LoadBalancer (en cloud)
apiVersion: v1
kind: Service
metadata:
  name: hello-go-lb
spec:
  type: LoadBalancer
  selector:
    app: hello-go
  ports:
  - port: 80
    targetPort: 8080
```

**Verificación**:

```bash
# Ver servicios
kubectl get svc

# Test ClusterIP desde otro Pod
kubectl run -it busybox --image=busybox --restart=Never -- sh
# Dentro del container
wget -O- http://hello-go-internal:80

# Test NodePort
curl http://<node-ip>:30080

# Test LoadBalancer
kubectl get svc hello-go-lb -w
```

### Ejercicio 4: ConfigMap - Configuración Sin Hardcode

**Objetivo**: Externalizar configuración usando ConfigMaps.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_ENV: "production"
  LOG_LEVEL: "INFO"
  DB_HOST: "postgres.default.svc.cluster.local"
  app.conf: |
    [app]
    name=HelloGo
    version=1.0
    [logging]
    level=INFO
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-go-configurable
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-go
  template:
    metadata:
      labels:
        app: hello-go
    spec:
      containers:
      - name: app
        image: myregistry/hello-go:1.0
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: app-config
        volumeMounts:
        - name: config-file
          mountPath: /etc/config
      volumes:
      - name: config-file
        configMap:
          name: app-config
          items:
          - key: app.conf
            path: app.conf
```

**Aplicación Go que lee config**:

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

func main() {
    appEnv := os.Getenv("APP_ENV")
    logLevel := os.Getenv("LOG_LEVEL")
    
    config, _ := ioutil.ReadFile("/etc/config/app.conf")
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Env: %s, LogLevel: %s\n", appEnv, logLevel)
        fmt.Fprintf(w, "Config: %s\n", string(config))
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Ejercicio 5: Full Stack - App con Base de Datos

**Objetivo**: Desplegar aplicación completa con Database, Service Discovery y Persistencia.

```yaml
# Namespace
apiVersion: v1
kind: Namespace
metadata:
  name: app-stack
---
# ConfigMap para la app
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: app-stack
data:
  DATABASE_URL: "postgres://postgres:password@postgres:5432/appdb"
---
# Secret para credenciales
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
  namespace: app-stack
type: Opaque
stringData:
  DB_PASSWORD: "password123"
---
# PersistentVolumeClaim para datos
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-data
  namespace: app-stack
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
---
# Service para PostgreSQL
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: app-stack
spec:
  clusterIP: None
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
---
# StatefulSet para PostgreSQL
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: app-stack
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:13
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: appdb
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: DB_PASSWORD
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: postgres-data
---
# Service para la aplicación
apiVersion: v1
kind: Service
metadata:
  name: hello-go-app
  namespace: app-stack
spec:
  type: NodePort
  selector:
    app: hello-go
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080
---
# Deployment de la aplicación Go
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-go
  namespace: app-stack
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: hello-go
  template:
    metadata:
      labels:
        app: hello-go
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - hello-go
              topologyKey: kubernetes.io/hostname
      containers:
      - name: app
        image: myregistry/hello-go:1.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: DATABASE_URL
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: DB_PASSWORD
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
# HorizontalPodAutoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hello-go-hpa
  namespace: app-stack
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: hello-go
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

**Aplicación Go completa**:

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync/atomic"
    "time"
    
    _ "github.com/lib/pq"
)

var (
    db    *sql.DB
    ready int32 = 0
)

func main() {
    // Conectar a base de datos
    dbURL := os.Getenv("DATABASE_URL")
    var err error
    db, err = sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer db.Close()
    
    if err = db.Ping(); err != nil {
        log.Fatalf("Failed to ping DB: %v", err)
    }
    
    atomic.StoreInt32(&ready, 1)
    
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/health", handleHealth)
    http.HandleFunc("/ready", handleReady)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, `{"status":"ok","timestamp":"`+time.Now().String()+`"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    if err := db.Ping(); err == nil {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"alive"}`)
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, `{"status":"dead"}`)
    }
}

func handleReady(w http.ResponseWriter, r *http.Request) {
    if atomic.LoadInt32(&ready) == 1 {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"ready"}`)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        fmt.Fprint(w, `{"status":"not ready"}`)
    }
}
```

---

## Resumen

En este capítulo, hemos explorado Kubernetes desde sus conceptos fundamentales hasta patrones avanzados:

1. **Conceptos básicos**: Pods, Nodes, Services, Namespaces
2. **Deployments**: Replicación y actualizaciones
3. **Services**: Descubrimiento y exposición
4. **ConfigMaps y Secrets**: Configuración segura
5. **StatefulSets**: Aplicaciones con estado
6. **Persistent Volumes**: Almacenamiento duradero
7. **Health Probes**: Confiabilidad
8. **Resource Management**: Optimización
9. **Ingress**: Acceso externo
10. **Buenas prácticas**: RBAC, Network Policies, monitoreo

Kubernetes permite que Go escale desde una sola instancia a miles de réplicas distribuidas automáticamente, proporcionando autorrecuperación, balanceo de carga y gestión de recursos sofisticada.

### Puntos Clave

✅ Los Pods son efímeros, los Services son estables  
✅ Deployments automatizan actualizaciones  
✅ StatefulSets para aplicaciones con estado  
✅ ConfigMaps para configuración no sensible  
✅ Secrets para datos sensibles (siempre)  
✅ Health probes para confiabilidad  
✅ Resource limits para estabilidad del cluster  
✅ RBAC para seguridad  
✅ Namespaces para aislamiento  

Dominar Kubernetes es esencial para desarrolladores Go en la era de la nube moderna.

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/47-kubernetes-basics/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/47-kubernetes-basics):

```bash
cd examples/47-kubernetes-basics
go run .
```
