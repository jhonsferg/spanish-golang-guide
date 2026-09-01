# Capítulo 50: Distributed systems - Sistemas distribuidos a escala

## Índice del Capítulo 50

1. [50.1 Desafíos Fundamentales de Sistemas Distribuidos](#501-desafíos-fundamentales-de-sistemas-distribuidos)
2. [50.2 CAP Theorem - El Triángulo Imposible](#502-cap-theorem--el-triángulo-imposible)
3. [50.3 Consistencia Eventual y Causalidad](#503-consistencia-eventual-y-causalidad)
4. [50.4 Algoritmos de Consenso - Raft y Paxos](#504-algoritmos-de-consenso---raft-y-paxos)
5. [50.5 Transacciones Distribuidas](#505-transacciones-distribuidas)
6. [50.6 Message Queues y Event Sourcing](#506-message-queues-y-event-sourcing)
7. [50.7 Load Balancing - Distribución de Carga](#507-load-balancing--distribución-de-carga)
8. [50.8 Replicación de Datos](#508-replicación-de-datos)
9. [50.9 Sharding y Particionamiento](#509-sharding-y-particionamiento)
10. [50.10 Fault Tolerance - Tolerancia a Fallos](#5010-fault-tolerance--tolerancia-a-fallos)
11. [50.11 Buenas Prácticas y Trade-offs Operacionales](#5011-buenas-prácticas-y-trade-offs-operacionales)
12. [Ejercicios Progresivos](#ejercicios-progresivos)

---

## 50.1 Desafíos Fundamentales de Sistemas Distribuidos

### La Realidad: No Todo es Local

En programación local, tienes garantías:
- Tu función ejecuta en un procesador
- La memoria es accesible y rápida
- Si falla, todo falla (pero de forma clara)
- El tiempo es lineal y predecible

**En sistemas distribuidos, nada de esto es verdad:**

```
        Node A              Network              Node B
        
        f() {           [Packet Loss]           ready?
          ready?    ─→  [Latency]         ─→  process()
        }            ←─  [Corruption]     ←─   done!
                         [Delay: ???]
```

### Los Tres Problemas Fundamentales

#### 1. **Latencia - El Tiempo No Es Instantáneo**

En una máquina local:
```go
// Acceso a memoria: ~100 nanosegundos
x := data[i]
```

En una red:
```go
// Consulta a otro servidor: ~1-100 milisegundos
// = 1,000,000x más lento
response := callRemoteServer()
```

**Implicación crítica:** La latencia es tan grande que cambia completamente el diseño del sistema.

```
Timeline de latencias comparativas:
1 ns        ← Ciclo de CPU
10 ns       ← L1 cache hit
100 ns      ← L2 cache hit
1,000 ns    ← L3 cache hit
10,000 ns   ← Main memory access
1,000,000 ns = 1 ms ← Red local
10 ms       ← Internet típico
100 ms      ← Intercontinental
1,000 ms    ← Satellites (muy lento)
```

Esto significa que:
- Una consulta de red = 10,000,000 accesos a memoria
- Debes batching: No preguntes 1 cosa, pregunta 1000
- Debes caching: Evita preguntar si puedes

#### 2. **Partial Failures - Fallos Parciales**

En un programa local:
```go
// O funciona TODO o falla TODO
// Clear distinction
if x := riskyOperation(); x != nil {
    // Error, todo falla
}
```

En distribuido:
```go
// Node A envía a Node B
// Posibilidades:
// 1. Node B recibe y procesa ✓
// 2. Node B no recibe ✗
// 3. Node B procesa pero no envía respuesta ✗
// 4. Network falla a mitad ✗
// 5. Node B falla durante procesamiento ✗
// 6. Node B responde pero tú no recibes ✗

// NO SABES QUÉ PASÓ
```

**El problema del "ambigüedad":**
- Enviaste una transacción de dinero
- La conexión se cortó
- ¿Se procesó? ¿No?
- Tienes que asumir que PODRÍA haber sido procesada

#### 3. **Asynchronous Communication - No Hay Respuestas Garantizadas**

En local, la comunicación es síncrona:
```go
result := function()  // Esperas siempre
// result está garantizado
```

En distribuido:
```go
// Envías mensaje
// Esperas
// ... y esperas
// ¿Murió el servidor? ¿Está procesando?
// NO TIENES IDEA
```

### Teorema de Brewer: El CAP

**Teorema fundamental (2000):**
En todo sistema distribuido, solo puedes tener 2 de 3:

```
         Consistency (C)
                △
               /|\
              / | \
             /  |  \
            /   |   \
        AP /    |    \ CA
          /     |     \
         / Availab| \ 
        /________|____\
       Partition (P)
```

- **C (Consistency):** Todos leen el mismo valor siempre
- **A (Availability):** El sistema siempre responde
- **P (Partition Tolerance):** El sistema sigue funcionando si la red se particiona

**La verdad:**
- En internet, las particiones EXISTEN
- Luego, DEBES elegir entre C (SQL) o A (NoSQL)

---

## 50.2 CAP Theorem - El Triángulo Imposible

### Entender CAP: No es Binario

La mayoría piensa en CAP de forma binaria: "Elige 2"

Pero la realidad es más sutil:

```
Sistema actual: CP
- Cuando NO hay partición de red:
  → Tenemos Consistency Y Availability
  
- Cuando SÍ hay partición de red:
  → Elegimos entre:
    • Consistency: Rechazar operaciones (CP)
    • Availability: Permitir divergencia (AP)
```

### Patrón CA - Consistencia y Disponibilidad

**Requisito: Red perfecta (no realista)**

```go
// Ejemplo: Base de datos relacional tradicional

/*
Arquitectura:
  
  Client A        Client B
    |               |
    └───────┬───────┘
            |
      Database
      (monolítica)
      
Garantías:
- Todos ven el mismo estado (Consistency)
- Siempre responde (Availability)
- PERO: Si la red cae entre clients y DB → DOWN

Ejemplos: PostgreSQL standalone, Oracle, MySQL sin replicación
*/

// Código conceptual:
type CAsystem struct {
    db *postgres.DB
}

func (s *CAsystem) Transfer(from, to string, amount float64) error {
    return s.db.Transaction(func(tx *sql.Tx) error {
        // Atomic: TODO o NADA
        return nil
    })
}
```

**Problema práctico:**
```
Time:  0ms              100ms
       |                |
     User A          Network dies!
       |
    Request to DB
       |
       ├─→ Consistent write ✓
       │
       └─→ Error: DB unreachable ✗
           
Result: User sees error, goes to competitor
```

### Patrón CP - Consistencia y Tolerancia a Partición

**"Sacrificamos Availability"**

```go
/*
Arquitectura: Replicado pero conservador

  Client         Leader (Master)
    |                 |
    └─────────────────┘
          |
      Replicas
      (followers)
      
Si la red se particiona:
- Leader sigue disponible en su partición
- Followers en otra partición: LEEN SOLO (stale data)
- Nuevos writes: RECHAZADOS si no tienen quórum

Garantía: Nunca ves inconsistencia
Costo: Indisponibilidad durante partición
*/

type CPSystem struct {
    leader   *Node
    followers []*Node
    quorumSize int
}

func (s *CPSystem) Write(key string, value interface{}) error {
    // Replicamos a quórum ANTES de confirmar
    if !s.replicateToQuorum(key, value) {
        return errors.New("no quórum: operación rechazada")
    }
    return nil
}

// Ejemplos: Zookeeper, Etcd, HBase
```

**Ventaja para datos críticos:**
```
Banco:
- Transfers deben ser consistentes (CP)
- Es mejor perder momentánea disponibilidad
  que tener inconsistencia de dinero
```

### Patrón AP - Disponibilidad y Tolerancia a Partición

**"Sacrificamos Consistencia"**

```go
/*
Arquitectura: Replicado y optimista

  Client A           Client B
    |                  |
    ├─→ Replica A  ←─→  Replica B
    │       |            |
    │       └────────┬───┘
    │                |
    └──────────────────
    
Si red se particiona:
- A y B escriben a replicas diferentes
- Cada uno recibe respuesta OK ✓
- Cuando red se recupera:
  → Conflicto! A tiene X=1, B tiene X=2
  → "Eventual consistency" - eventualmente consistentes

Garantía: Siempre disponible
Costo: Posible inconsistencia temporal
*/

type APSystem struct {
    replicas map[string]*Node
    version  map[string]int64  // vector clock
}

func (s *APSystem) Write(key string, value interface{}) error {
    // Escribimos a NUESTRO nodo INMEDIATAMENTE
    // (sin esperar quórum)
    s.replicas["local"].Write(key, value)
    
    // Replicamos en background (best effort)
    go s.asyncReplicateToAll(key, value)
    
    return nil  // Siempre OK
}

// Ejemplos: Dynamo (AWS), Cassandra, Riak
```

**Ventaja para escalabilidad:**
```
Social Network:
- Cuando escribo un tweet:
  → Se replica inmediatamente a mi región
  → Mis followers lo ven al instante
  → Otras regiones lo reciben lentamente
  → Eventual consistency es aceptable
```

### Matriz de Trade-offs en CAP

```
                CP              AP
                
Consistency    FUERTE          EVENTUAL
               (C ≈ 1.0)        (C ≈ 0.95-0.99)

Availability   REDUCIDA         EXCELENTE
               (A ≈ 0.95)       (A ≈ 0.9999)

Partición      TOLERANTE        TOLERANTE
               (rechazo)        (divergencia)

Casos de uso   Bancos           Redes sociales
               Reservas         Caché distribuido
               Inventario       Recomendaciones

Go tools       Etcd             Consul
               Zookeeper        Redis
                                Cassandra
```

---

## 50.3 Consistencia Eventual y Causalidad

### El Concepto: No Ahora, Pero Pronto

**Consistencia Eventual = "Promesa futura"**

```go
// Escenario:
// 1. Escribo X = 1 en nodo A
// 2. Inmediatamente leo X en nodo B
// 3. Puedo obtener X = 0 (viejo) ← Esta es la eventual part
// 4. Leo nuevamente en 50ms
// 5. Ahora X = 1 ✓ ← Eventualmente consistente

// Garantía:
// "Si NO escribimos nada más,
//  eventualmente todos verán X = 1"

time := 0ms
write(X=1) → Nodo A
│
├─→ 10ms: read(X) en B → 0 (stale, ✗ inconsistente)
│
├─→ 30ms: read(X) en B → 1 ✓
│
├─→ 50ms: read(X) en C → 1 ✓
│
└─→ 100ms: Todos ven X = 1 (eventualmente consistente ✓)
```

### Vector Clocks - Causalidad sin Reloj Global

**El problema:** ¿Qué pasó primero sin un reloj global?

```go
// Reloj global no existe en distribuido
// Relojes delos servidores pueden estar desfasados

Servidor A: timestamp 10:00:00.001
Servidor B: timestamp 10:00:00.000  ← Un ms antes por error

// ¿Cuál operación fue primero?
```

**Solución: Vector Clocks**

```go
// Vector clock: Lista de (Node ID, contador)

type VectorClock map[string]int64

// Operación 1 en A:
A[A] = 1
VC1 := {A: 1}

// Operación 2 en B:
// B recibe VC1, actualiza:
VC1_received := {A: 1}
B[B] = 1
VC2 := {A: 1, B: 1}

// Comparación:
// VC1 < VC2 porque todos los contadores de VC1 ≤ VC2
// Luego: Op1 causalmente antes que Op2 ✓

// Operación 3 en A:
// A recibe VC2 = {A:1, B:1}
A[A] = 2
VC3 := {A: 2, B: 1}

// Ahora VC2 < VC3 (causalidad preservada)
```

**Implementación en Go:**

```go
package distributed

import (
    "fmt"
    "sync"
)

type VectorClock struct {
    clocks map[string]int64
    mu     sync.RWMutex
}

func NewVectorClock(nodeID string) *VectorClock {
    return &VectorClock{
        clocks: map[string]int64{nodeID: 0},
    }
}

// Incrementar reloj local
func (vc *VectorClock) Increment(nodeID string) {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    vc.clocks[nodeID]++
}

// Recibir evento de otro nodo
func (vc *VectorClock) Update(other *VectorClock) {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    
    other.mu.RLock()
    defer other.mu.RUnlock()
    
    for nodeID, timestamp := range other.clocks {
        if vc.clocks[nodeID] < timestamp {
            vc.clocks[nodeID] = timestamp
        }
    }
    
    // Incrementar propio
    if nodeID, exists := vc.clocks["self"]; exists {
        vc.clocks["self"] = nodeID + 1
    }
}

// Comparar causalidad
func (vc *VectorClock) Happens Before(other *VectorClock) bool {
    vc.mu.RLock()
    defer vc.mu.RUnlock()
    
    other.mu.RLock()
    defer other.mu.RUnlock()
    
    atLeastOneLess := false
    
    for nodeID, ts := range vc.clocks {
        otherTS := other.clocks[nodeID]
        if ts > otherTS {
            return false  // Not before
        }
        if ts < otherTS {
            atLeastOneLess = true
        }
    }
    
    return atLeastOneLess
}

// Concurrent si neither before other
func (vc *VectorClock) Concurrent(other *VectorClock) bool {
    return !vc.HappensBefore(other) && !other.HappensBefore(vc)
}

func (vc *VectorClock) String() string {
    vc.mu.RLock()
    defer vc.mu.RUnlock()
    return fmt.Sprintf("%v", vc.clocks)
}
```

### Quorum Reads - Tradeoff de Disponibilidad y Consistencia

**Idea:** Lee de múltiples réplicas, usa mayoritarios

```go
type QuorumStore struct {
    replicas     []*Node
    quorumSize   int
    replicationFactor int
}

// Quórum write: Escribe a R réplicas
func (qs *QuorumStore) Write(key string, value interface{}) error {
    successCount := 0
    errors := 0
    
    for _, replica := range qs.replicas {
        if replica.Write(key, value) == nil {
            successCount++
        } else {
            errors++
        }
    }
    
    // Si alcanzamos quórum, éxito
    if successCount >= qs.quorumSize {
        return nil
    }
    
    return fmt.Errorf("failed: %d/%d", errors, len(qs.replicas))
}

// Quórum read: Lee de R réplicas, retorna mayoritario
type ReadResult struct {
    Value   interface{}
    Version int64
}

func (qs *QuorumStore) Read(key string) (interface{}, error) {
    results := make([]ReadResult, 0)
    
    for _, replica := range qs.replicas[:qs.quorumSize] {
        val, version, err := replica.Read(key)
        if err == nil {
            results = append(results, ReadResult{val, version})
        }
    }
    
    // Retornar la versión más reciente (quórum de lecturas)
    if len(results) >= (qs.quorumSize / 2) {
        maxVersion := int64(-1)
        var maxValue interface{}
        
        for _, r := range results {
            if r.Version > maxVersion {
                maxVersion = r.Version
                maxValue = r.Value
            }
        }
        
        return maxValue, nil
    }
    
    return nil, errors.New("quórum read failed")
}

/*
Tuning de Quórum:

Si replicationFactor = 3:
- Write quórum = 2: Tolera 1 fallo
- Read quórum = 2: Lees valor consistente

Si replicationFactor = 5:
- Write quórum = 3: Tolera 2 fallos
- Read quórum = 3: Muy consistente pero más lento

Trade-off:
N = replication factor
W = write quorum (writes)
R = read quorum (reads)

Si W + R > N → Garantía de leer lo que escribiste
Si W + R ≤ N → Posible inconsistencia
*/
```

---

## 50.4 Algoritmos de Consenso - Raft y Paxos

### El Problema de Consenso

**Escenario:**
- 5 servidores
- 3 quieren hacer X
- 2 quieren hacer Y
- ¿Cómo llegan a acuerdo?
- ¿Cómo recuperarse si algunos fallan?

```go
// Sin consenso:
Server A: "Hagamos X"
Server B: "Hagamos Y"
Server C: "Creo que X, pero no estoy seguro"
Server D: "💥 CRASHED"
Server E: "💥 CRASHED"

Resultado: CAOS, inconsistencia

// Con consenso (Raft):
→ Se elige un líder
→ El líder propone: "Hagamos X"
→ Mayoría responde: "OK, X"
→ Confirmado: TODOS hacen X
```

### Raft: Consenso Simplificado

**Raft se basa en tres mecanismos:**

#### 1. **Leader Election - Elegir Líder**

```
Estados posibles: FOLLOWER, CANDIDATE, LEADER

Inicialmente: Todos FOLLOWER

Timeout sin heartbeat del líder:
  → FOLLOWER → CANDIDATE (me propongo)
  → Pido votos
  → Si gano mayoría → LEADER
  → Envío heartbeats cada 50ms

Si hay partición de red:
  → Viejo LEADER sigue enviando heartbeats en su partición
  → Nuevo CANDIDATE en otra partición gana votación
  → Dos "líderes" pero en particiones diferentes
  → Cuando se reconecten:
    → Líder con TERM más alto gana
    → Otro se demota a FOLLOWER
```

**Código Go - Leader Election Simplificado:**

```go
package raft

import (
    "math/rand"
    "time"
)

type RaftNode struct {
    id          string
    state       State  // FOLLOWER, CANDIDATE, LEADER
    currentTerm int64
    votedFor    string
    
    peers       []*RaftNode
    
    commitIndex int
    lastApplied int
    
    // Timers
    electionTimer  *time.Timer
    heartbeatTimer *time.Timer
}

type State int

const (
    FOLLOWER State = iota
    CANDIDATE
    LEADER
)

func (rn *RaftNode) electionTimeout() time.Duration {
    // Random entre 150-300ms
    return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

func (rn *RaftNode) startElection() {
    rn.state = CANDIDATE
    rn.currentTerm++
    rn.votedFor = rn.id
    
    // Pedir votos a peers
    voteCount := 1  // Me voto a mí mismo
    
    for _, peer := range rn.peers {
        if peer.id == rn.id {
            continue
        }
        
        // Enviar RequestVote RPC
        vote := rn.requestVote(peer)
        if vote {
            voteCount++
        }
    }
    
    // Si ganamos mayoría
    if voteCount > len(rn.peers)/2 {
        rn.state = LEADER
        rn.sendHeartbeats()
    } else {
        rn.state = FOLLOWER
    }
}

func (rn *RaftNode) requestVote(peer *RaftNode) bool {
    // RPC: ¿Me das tu voto?
    // Peer responde:
    // - SÍ si: currentTerm ≥ su term Y (no votó O votó por mí)
    // - NO en otro caso
    
    // Simplificado:
    if peer.currentTerm < rn.currentTerm {
        if peer.votedFor == "" || peer.votedFor == rn.id {
            return true  // Grant vote
        }
    }
    return false
}

func (rn *RaftNode) sendHeartbeats() {
    for _, peer := range rn.peers {
        if peer.id == rn.id {
            continue
        }
        go peer.AppendEntries(rn.currentTerm, rn.id)
    }
}

func (rn *RaftNode) AppendEntries(term int64, leaderID string) {
    if term > rn.currentTerm {
        rn.currentTerm = term
        rn.state = FOLLOWER
        rn.votedFor = ""
    }
    
    // Resetear election timer (líder está vivo)
    rn.electionTimer.Reset(rn.electionTimeout())
}
```

#### 2. **Log Replication - Sincronizar Logs**

```
Escenario con 3 servidores:

Servidor A (LEADER):
  Log: [1: X], [2: Y], [3: Z] ← Qué queremos que pase
  
Servidor B (FOLLOWER):
  Log: [1: X], [2: Y]          ← Retraso
  
Servidor C (FOLLOWER):
  Log: [1: X]                  ← Más retraso

Raft replicates:
  LEADER envía: "Aquí está [3: Z]"
  B recibe: [1: X], [2: Y], [3: Z] ✓
  C recibe: [1: X], [2: Y], [3: Z] ✓

Una vez que MAYORÍA tiene [3: Z]:
  → [3: Z] es COMMITTED
  → Se ejecuta en la state machine
```

**Código Go:**

```go
type LogEntry struct {
    Term  int64
    Index int64
    Command interface{}
}

type RaftNode struct {
    // ... (como antes)
    log []LogEntry
    
    // Leader state
    nextIndex  map[string]int64  // Siguiente entry a enviar
    matchIndex map[string]int64  // Última replicada
}

func (rn *RaftNode) AppendEntries(
    term int64,
    leaderID string,
    entries []LogEntry,
) {
    if term > rn.currentTerm {
        rn.currentTerm = term
        rn.state = FOLLOWER
    }
    
    // Agregar entries al log
    for _, entry := range entries {
        rn.log = append(rn.log, entry)
    }
    
    // Actualizar commit index
    if len(entries) > 0 && entries[len(entries)-1].Index > rn.commitIndex {
        rn.commitIndex = entries[len(entries)-1].Index
        rn.apply()  // Aplica a state machine
    }
}

// Replicar nueva entrada
func (rn *RaftNode) Replicate(command interface{}) error {
    if rn.state != LEADER {
        return errors.New("not leader")
    }
    
    entry := LogEntry{
        Term:    rn.currentTerm,
        Index:   int64(len(rn.log)),
        Command: command,
    }
    
    rn.log = append(rn.log, entry)
    
    // Enviar a todos los followers
    replicationCount := 1  // Yo lo tengo
    
    for _, peer := range rn.peers {
        if peer.id == rn.id {
            continue
        }
        
        idx := rn.nextIndex[peer.id]
        entries := rn.log[idx:]
        
        if peer.AppendEntries(rn.currentTerm, rn.id, entries) {
            rn.matchIndex[peer.id]++
            replicationCount++
        }
    }
    
    // Si mayoría tiene la entrada, commit
    if replicationCount > len(rn.peers)/2 {
        rn.commitIndex = entry.Index
        rn.apply()
        return nil
    }
    
    return errors.New("replication failed")
}
```

#### 3. **Safety - Garantías de Corrección**

```
Raft garantiza:
1. Election Safety: Máximo 1 líder por term
2. Log Matching: Si dos logs tienen entry con índice y term iguales,
                 todas las entries anteriores también son iguales
3. Leader Completeness: Si una entrada fue commited,
                        todos los líderes futuros la tendrán
4. State Machine Safety: Nunca se aplica dos entries diferentes
                        con el mismo índice
```

### Paxos: Más Poderoso Pero Complejo

```go
/*
Paxos vs Raft:

PAXOS:
+ Más general (no necesita líder)
+ Tolerancia bizantina (mentirosos)
- Muy complejo (difícil de implementar)
- 2 fases: Prepare y Accept

RAFT:
+ Más simple (basado en líder)
+ Fácil de entender y debuggear
+ Usado en etcd, Consul
- Requiere líder funcional
- 1 fase: Leader simplemente replica

Go community prefiere Raft
Ejemplo: github.com/hashicorp/raft
*/

// Paxos simplificado (conceptual):
type Proposer struct {
    proposalID int
    value      interface{}
}

// Fase 1: PREPARE
// Proposer propone: "¿Aceptarías mi propuesta?"
// Acceptor responde: "Sí, pero si recibo propuesta >99, la considero"

// Fase 2: ACCEPT
// Si mayoría dice "sí" en PREPARE:
//   Proposer envía: "Aquí está mi valor"
//   Acceptor acepta si no recibió propuesta con ID mayor

// Ventaja: Tolerancia a fallo de cualquier node
// Desventaja: 2 fases vs Raft's 1 fase
```

---

## 50.5 Transacciones Distribuidas

### El Problema: ACID en Múltiples Máquinas

```go
// Local (Easy):
db.Transaction(func(tx *Tx) {
    tx.Update("account_a", -100)
    tx.Update("account_b", +100)
    // Atomic: TODO o NADA
})

// Distribuido (Hard):
// account_a en Server A
// account_b en Server B

// Si ambos actualizan:
//   ✓ Ambos: Éxito
//   ✗ Uno sí, uno no: INCONSISTENCIA
//   ✗ Uno falla a mitad: ¿QUÉ HACEMOS?
```

### Two-Phase Commit (2PC)

**Fase 1: Prepare (Pregunta)**

```
Coordinator              Server A              Server B
    |                       |                      |
    ├─ "¿Puedes hacer X?" ──→|                      |
    |                       ├─ Locks resources      |
    |                       ├─ Simula ejecución     |
    |                       └─ "Sí, puedo" ──→|     |
    |                                         |     |
    |                    ← "Sí, puedo" ←──┤
    |<── "OK" ─────────────────────────┤
    |                                         |
```

**Fase 2: Commit (Ejecuta)**

```
Coordinator              Server A              Server B
    |                       |                      |
    ├─ "COMMIT" ────────────→|                      |
    |                       ├─ Ejecuta X            |
    |                       └─ "Done" ──→|          |
    |                                  ←─┤
    |<── "Committed" ←─────────────────┤
    |
```

**Implementación Go:**

```go
package dist

import (
    "context"
    "errors"
    "sync"
)

type Transaction2PC struct {
    id       string
    peers    []string
    state    TxState
    entries  map[string]interface{}
    locks    map[string]bool
}

type TxState int

const (
    INITIAL TxState = iota
    PREPARE
    COMMIT
    ABORT
)

func (tx *Transaction2PC) Prepare(ctx context.Context) error {
    tx.state = PREPARE
    
    // Fase 1: Preguntar a todos
    results := make(chan bool, len(tx.peers))
    wg := sync.WaitGroup{}
    
    for _, peer := range tx.peers {
        wg.Add(1)
        go func(p string) {
            defer wg.Done()
            if tx.canExecute(p) {
                tx.locks[p] = true  // Lock recurso
                results <- true
            } else {
                results <- false
            }
        }(peer)
    }
    
    wg.Wait()
    close(results)
    
    // Contar votos
    yesCount := 0
    for vote := range results {
        if vote {
            yesCount++
        } else {
            // Si UNO dice no, abort
            tx.Abort(ctx)
            return errors.New("prepare failed")
        }
    }
    
    if yesCount == len(tx.peers) {
        return tx.Commit(ctx)
    }
    
    return tx.Abort(ctx)
}

func (tx *Transaction2PC) Commit(ctx context.Context) error {
    tx.state = COMMIT
    
    // Fase 2: TODOS ejecutan
    errChan := make(chan error, len(tx.peers))
    wg := sync.WaitGroup{}
    
    for _, peer := range tx.peers {
        wg.Add(1)
        go func(p string) {
            defer wg.Done()
            if err := tx.executeOnPeer(p); err != nil {
                errChan <- err
            }
        }(peer)
    }
    
    wg.Wait()
    close(errChan)
    
    for err := range errChan {
        if err != nil {
            return err  // Alguno falló
        }
    }
    
    return nil
}

func (tx *Transaction2PC) Abort(ctx context.Context) error {
    tx.state = ABORT
    
    // Liberar locks
    for peer := range tx.locks {
        delete(tx.locks, peer)
    }
    
    return nil
}

func (tx *Transaction2PC) canExecute(peer string) bool {
    // Simular: ¿Puedo ejecutar esto?
    // En realidad: RPC al peer
    return true  // Simplificado
}

func (tx *Transaction2PC) executeOnPeer(peer string) error {
    // RPC: Ejecutar cambios
    return nil  // Simplificado
}

/*
Problemas con 2PC:

1. BLOCKING: Si Fase 1 OK pero Fase 2 falla...
   → Recursos quedan lockeados ← DEADLOCK POTENTIAL

2. Availability: Si coordinator cae en mitad:
   → ¿Commit o Abort? Nadie sabe
   → Timeout → Abort por defecto

3. Performance: N fases = N latencias de red

Uso real: Muy poco en sistemas distribuidos modernos
Razón: Demasiado lento y frágil
*/
```

### Sagas - Transacciones Distribuidas Modernas

**Idea:** En lugar de atomic transaction, ejecuta secuencia con compensations

```go
/*
Ejemplo: Transfer de dinero entre bancos

Approach tradicional (2PC):
  TransferAcceptWithdraw(from_account, 100)
    IF (ok)
      TransferDeposit(to_account, 100)
    ELSE
      Abort

Problema: Si paso 2 falla, paso 1 rollback se bloquea

Saga approach:
  Paso 1: Withdraw(from_account, 100)
          → Si falla: END (no iniciamos)
          → Si éxito: Continúa
          
  Paso 2: Deposit(to_account, 100)
          → Si falla: Ejecuta CompensateWithdraw(from_account, 100)
          → Si éxito: Continúa

Ventaja: Cada paso es independiente
*/

type Saga struct {
    steps []SagaStep
}

type SagaStep struct {
    name        string
    action      func() error
    compensation func() error
}

func (s *Saga) Execute() error {
    completed := []string{}
    
    for _, step := range s.steps {
        if err := step.action(); err != nil {
            // Falló, ejecutar compensations en reverso
            for i := len(completed) - 1; i >= 0; i-- {
                for _, step := range s.steps {
                    if step.name == completed[i] {
                        _ = step.compensation()  // Ignora errores en compensate
                        break
                    }
                }
            }
            return err
        }
        
        completed = append(completed, step.name)
    }
    
    return nil
}

// Uso:
transferSaga := &Saga{
    steps: []SagaStep{
        {
            name: "withdraw",
            action: func() error {
                return withdrawFromAccount("A", 100)
            },
            compensation: func() error {
                return depositToAccount("A", 100)  // Revert
            },
        },
        {
            name: "deposit",
            action: func() error {
                return depositToAccount("B", 100)
            },
            compensation: func() error {
                return withdrawFromAccount("B", 100)  // Revert
            },
        },
    },
}

if err := transferSaga.Execute(); err != nil {
    // Ya se compensó automáticamente
}
```

---

## 50.6 Message Queues y Event Sourcing

### Desacoplamiento con Message Queues

**Problema: Acoplamiento Síncrono**

```go
// API síncrona (acoplada):
func CreateOrder(order *Order) error {
    // 1. Crear orden en DB
    db.SaveOrder(order)
    
    // 2. Procesar pago (WAIT)
    if err := paymentService.Charge(order.Amount); err != nil {
        // Si falla, ¿ya guardamos la orden?
        return err
    }
    
    // 3. Enviar email (WAIT)
    emailService.Send(order.Email, "Gracias!")
    
    // 4. Actualizar inventario (WAIT)
    inventoryService.Decrease(order.Items)
    
    return nil
}

// Problema:
// - Si payment está lento: Toda operación es lenta
// - Si email falla: Toda operación falla
// - Si inventario está down: Orden no se completa
```

**Solución: Message Queue**

```go
type MessageQueue interface {
    Publish(topic string, message interface{}) error
    Subscribe(topic string, handler func(interface{})) error
}

func CreateOrder(order *Order, mq MessageQueue) error {
    // 1. Guardar orden
    db.SaveOrder(order)
    
    // 2. Publicar event (NO ESPERA)
    mq.Publish("order.created", OrderCreatedEvent{
        OrderID: order.ID,
        Amount:  order.Amount,
    })
    
    return nil  // Inmediato ✓
}

// Los servicios escuchan en background:

// Payment Service
func init(mq MessageQueue) {
    mq.Subscribe("order.created", func(msg interface{}) {
        event := msg.(OrderCreatedEvent)
        paymentService.Charge(event.Amount)
    })
}

// Email Service
func init(mq MessageQueue) {
    mq.Subscribe("order.created", func(msg interface{}) {
        event := msg.(OrderCreatedEvent)
        emailService.Send(event.Email, "Gracias!")
    })
}

// Inventory Service
func init(mq MessageQueue) {
    mq.Subscribe("order.created", func(msg interface{}) {
        event := msg.(OrderCreatedEvent)
        inventoryService.Decrease(event.Items)
    })
}

// Ventajas:
// - CreateOrder es rápido (O ms)
// - Servicios trabajan en paralelo
// - Si uno falla, otros continúan
// - Desacoplado y resiliente
```

**Implementación simple en Go:**

```go
package mq

import (
    "sync"
    "time"
)

type SimpleQueue struct {
    messages map[string][]interface{}
    handlers map[string][]func(interface{})
    mu       sync.RWMutex
}

func NewSimpleQueue() *SimpleQueue {
    return &SimpleQueue{
        messages: make(map[string][]interface{}),
        handlers: make(map[string][]func(interface{})),
    }
}

func (q *SimpleQueue) Publish(topic string, msg interface{}) error {
    q.mu.Lock()
    q.messages[topic] = append(q.messages[topic], msg)
    q.mu.Unlock()
    
    // Procesar suscriptores
    q.mu.RLock()
    handlers := q.handlers[topic]
    q.mu.RUnlock()
    
    for _, handler := range handlers {
        // Async processing
        go handler(msg)
    }
    
    return nil
}

func (q *SimpleQueue) Subscribe(topic string, handler func(interface{})) error {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    q.handlers[topic] = append(q.handlers[topic], handler)
    
    return nil
}

// Real-world example: Apache Kafka, RabbitMQ, AWS SQS
```

### Event Sourcing - Source of Truth

**Concepto: El event log es la fuente de verdad**

```go
/*
Enfoque tradicional (Estado):
  db.accounts: {id: 1, balance: 500}
  
  Problema:
  - ¿Cómo llegó a 500?
  - ¿Qué pasó antes?
  - Auditoría limitada

Event Sourcing (Eventos):
  event_log:
    - {timestamp: 10:00, account: 1, type: "deposit", amount: 1000}
    - {timestamp: 10:05, account: 1, type: "withdraw", amount: 200}
    - {timestamp: 10:10, account: 1, type: "withdraw", amount: 300}
    
  Estado derivado:
    balance = 1000 - 200 - 300 = 500
    
  Ventajas:
  - Auditoría completa
  - Reconstruir estado a cualquier momento
  - Debug facilitado
  - Replay de eventos
*/

type Event struct {
    ID        string
    Timestamp time.Time
    Type      string
    Data      map[string]interface{}
}

type EventStore struct {
    events []Event
    mu     sync.RWMutex
}

func (es *EventStore) Append(event Event) error {
    es.mu.Lock()
    defer es.mu.Unlock()
    
    event.Timestamp = time.Now()
    event.ID = generateID()  // UUID
    
    es.events = append(es.events, event)
    return nil
}

func (es *EventStore) GetEvents(filter func(Event) bool) []Event {
    es.mu.RLock()
    defer es.mu.RUnlock()
    
    var results []Event
    for _, event := range es.events {
        if filter(event) {
            results = append(results, event)
        }
    }
    return results
}

// Reconstruir estado desde events
type BankAccount struct {
    ID      string
    Balance float64
}

func (es *EventStore) ReplayAccount(accountID string) BankAccount {
    events := es.GetEvents(func(e Event) bool {
        return e.Data["account_id"] == accountID
    })
    
    acc := BankAccount{ID: accountID, Balance: 0}
    
    for _, event := range events {
        switch event.Type {
        case "deposit":
            acc.Balance += event.Data["amount"].(float64)
        case "withdraw":
            acc.Balance -= event.Data["amount"].(float64)
        }
    }
    
    return acc
}

// Ejemplo de uso:
es := &EventStore{}

es.Append(Event{
    Type: "deposit",
    Data: map[string]interface{}{
        "account_id": "acc1",
        "amount":     1000.0,
    },
})

es.Append(Event{
    Type: "withdraw",
    Data: map[string]interface{}{
        "account_id": "acc1",
        "amount":     200.0,
    },
})

// Reconstruir estado
account := es.ReplayAccount("acc1")
// Balance = 800
```

---

## 50.7 Load Balancing - Distribución de Carga

### Round-Robin Básico

```go
type RoundRobinBalancer struct {
    servers []*Server
    index   int
    mu      sync.Mutex
}

func (rb *RoundRobinBalancer) NextServer() *Server {
    rb.mu.Lock()
    defer rb.mu.Unlock()
    
    server := rb.servers[rb.index]
    rb.index = (rb.index + 1) % len(rb.servers)
    
    return server
}

/*
Problema: Servidores con diferentes capacidades

Server A: 32 cores
Server B: 2 cores

Round-robin envía igual carga a ambos
→ Server B está saturado
→ Server A tiene capacidad ociosa

Solución: Weighted round-robin
*/
```

### Consistent Hashing - Partición sin Rehashing

**Problema: Cuando agregas un servidor, todo se rehashea**

```go
/*
Hash tradicional:
  key = "user123"
  hash(key) % 3 = 1  → Server 1
  
  Añadimos Server 4:
  hash(key) % 4 = ?  → DIFERENTE SERVIDOR
  
  Resultado: Cache miss masivo (invalidation)

Consistent hashing:
  Coloca servidores en un círculo (ring)
  Busca el primer servidor en dirección horaria
  
  Al añadir servidor:
  - Solo las claves entre el nuevo y el anterior se mueven
  - ~1/N de las claves, no todas
*/

type ConsistentHash struct {
    servers   map[uint32]string  // hash → server_id
    ring      []uint32            // sorted hashes
    replication int              // virtual nodes
}

func (ch *ConsistentHash) AddServer(serverID string) {
    for i := 0; i < ch.replication; i++ {
        hash := hashFunc(fmt.Sprintf("%s:%d", serverID, i))
        ch.servers[hash] = serverID
        ch.ring = append(ch.ring, hash)
    }
    sort.Slice(ch.ring, func(i, j int) bool {
        return ch.ring[i] < ch.ring[j]
    })
}

func (ch *ConsistentHash) GetServer(key string) string {
    keyHash := hashFunc(key)
    
    // Buscar primer servidor ≥ keyHash
    for _, serverHash := range ch.ring {
        if serverHash >= keyHash {
            return ch.servers[serverHash]
        }
    }
    
    // Wrap around
    return ch.servers[ch.ring[0]]
}

/*
Ejemplo con 3 servidores:

Ring:
    0°: S1
  120°: S2
  240°: S3

key1 (hash=30°)  → S1 (primer servidor en dirección horaria)
key2 (hash=150°) → S2
key3 (hash=270°) → S3

Agregar S4 en 60°:
  key1 (30°)  → SIGUE EN S1 ✓
  key2 (150°) → SIGUE EN S2 ✓
  key3 (270°) → SIGUE EN S3 ✓
  
  Solo las claves entre S4 y S1 se mueven a S4
*/

package main

import (
    "fmt"
    "hash/fnv"
    "sort"
)

func hashFunc(key string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(key))
    return h.Sum32()
}

type ConsistentHashRing struct {
    nodes     map[uint32]string  // hash → nodeID
    sortedKeys []uint32           // sorted hash keys
    replicas  int
}

func NewConsistentHashRing(replicas int) *ConsistentHashRing {
    return &ConsistentHashRing{
        nodes:    make(map[uint32]string),
        replicas: replicas,
    }
}

func (chr *ConsistentHashRing) AddNode(nodeID string) {
    for i := 0; i < chr.replicas; i++ {
        hash := hashFunc(fmt.Sprintf("%s:%d", nodeID, i))
        chr.nodes[hash] = nodeID
        chr.sortedKeys = append(chr.sortedKeys, hash)
    }
    sort.Slice(chr.sortedKeys, func(i, j int) bool {
        return chr.sortedKeys[i] < chr.sortedKeys[j]
    })
}

func (chr *ConsistentHashRing) GetNode(key string) string {
    if len(chr.nodes) == 0 {
        return ""
    }
    
    hash := hashFunc(key)
    
    // Binary search para encontrar el primer nodo ≥ hash
    idx := sort.Search(len(chr.sortedKeys), func(i int) bool {
        return chr.sortedKeys[i] >= hash
    })
    
    if idx == len(chr.sortedKeys) {
        idx = 0  // Wrap around
    }
    
    return chr.nodes[chr.sortedKeys[idx]]
}
```

---

## 50.8 Replicación de Datos

### Master-Slave Replication

**Patrón: Un nodo escribe, múltiples leen**

```go
type MasterSlaveReplication struct {
    master  *Node
    slaves  []*Node
    binlog  []LogEntry  // Binary log
}

func (msr *MasterSlaveReplication) Write(key string, value interface{}) error {
    // 1. Escribir en master
    if err := msr.master.Write(key, value); err != nil {
        return err
    }
    
    // 2. Guardar en binlog
    entry := LogEntry{
        Command: fmt.Sprintf("SET %s = %v", key, value),
        Timestamp: time.Now(),
    }
    msr.binlog = append(msr.binlog, entry)
    
    // 3. Replicar a slaves en background
    go msr.replicateToSlaves(entry)
    
    return nil
}

func (msr *MasterSlaveReplication) Read(key string) (interface{}, error) {
    // Leer de master o slave (depende de política)
    // Opción 1: Leer de master siempre (consistencia fuerte)
    // Opción 2: Leer de slave (riesgo: stale data)
    
    return msr.master.Read(key)
}

func (msr *MasterSlaveReplication) replicateToSlaves(entry LogEntry) {
    for _, slave := range msr.slaves {
        go slave.ApplyLogEntry(entry)
    }
}

/*
Ventajas:
- Escrituras rápidas (solo master)
- Lecturas distribuidas (slaves)
- Recuperación simple

Desventajas:
- Slave lag (retraso en replicación)
- Master es SPOF (single point of failure)
- Lecturas pueden ser stale
*/

// Problema: Slave lag
// T=0: Write(X=1) en master
// T=5ms: Slave aún tiene X=0
// T=10ms: Write(X=2) en master
// T=50ms: Slave finalmente ve X=2 (¿Saltó X=1?)

type ReplicationLag struct {
    estimatedMs int64
}

func (msr *MasterSlaveReplication) EstimatedLag() time.Duration {
    // Calcula lag basado en lag de replicación
    // En producción: envía heartbeat y mide tiempo
    lastSlaveUpdate := time.Now()  // Simplificado
    return time.Since(lastSlaveUpdate)
}
```

### Master-Master Replication

**Patrón: Múltiples nodos pueden escribir**

```go
type MasterMasterReplication struct {
    nodes      []*Node
    vectorClocks map[string]*VectorClock
}

func (mmr *MasterMasterReplication) Write(
    nodeID string,
    key string,
    value interface{},
) error {
    // 1. Escribir en nodo local
    mmr.nodes[nodeID].Write(key, value)
    
    // 2. Incrementar vector clock
    mmr.vectorClocks[nodeID].Increment(nodeID)
    
    // 3. Replicar a otros nodos con vector clock
    for _, other := range mmr.nodes {
        if other.ID == nodeID {
            continue
        }
        go other.Replicate(key, value, mmr.vectorClocks[nodeID])
    }
    
    return nil
}

/*
Ventajas:
- Múltiples escrituras simultáneas
- No hay SPOF
- Alta disponibilidad

Desventajas:
- Conflictos de escritura (¿Qué valor es correcto?)
- Reconciliación compleja
- Causalidad difícil de mantener
*/

// Resolución de conflictos

type ConflictResolver interface {
    Resolve(a, b interface{}, ts_a, ts_b int64) interface{}
}

// Last-write-wins (simple pero puede perder datos)
type LastWriteWins struct{}

func (lww *LastWriteWins) Resolve(a, b interface{}, ts_a, ts_b int64) interface{} {
    if ts_a > ts_b {
        return a
    }
    return b
}

// Custom resolver (depende del tipo de dato)
type AccountBalance struct {
    Value int64
    Version int64
}

func ResolveBalance(a, b AccountBalance) AccountBalance {
    // Sumar en lugar de elegir uno
    return AccountBalance{
        Value:   a.Value + b.Value,
        Version: max(a.Version, b.Version) + 1,
    }
}
```

---

## 50.9 Sharding y Particionamiento

### Range-Based Sharding

```go
/*
Partición por rango:
  Shard 0: user_id 0-999
  Shard 1: user_id 1000-1999
  Shard 2: user_id 2000-2999
  
Cálculo: shard = user_id / 1000

Ventajas:
- Simple de implementar
- Fácil de encontrar datos

Desventajas:
- Desbalanceo (si IDs populares en un rango)
- Resharding complejo
*/

type RangeShardRouter struct {
    shards    map[int]*Shard
    rangeSize int64
}

func (rsr *RangeShardRouter) GetShard(userID int64) *Shard {
    shardID := userID / rsr.rangeSize
    return rsr.shards[int(shardID)]
}
```

### Hash-Based Sharding

```go
type HashShardRouter struct {
    shards    map[int]*Shard
    shardCount int
}

func (hsr *HashShardRouter) GetShard(key string) *Shard {
    hash := hashFunc(key)
    shardID := hash % uint32(hsr.shardCount)
    return hsr.shards[int(shardID)]
}

/*
Ventajas:
- Distribución uniforme
- Fácil agregar shards

Desventajas:
- Resharding requiere mover muchos datos
- No es range-query friendly
*/
```

### Consistent Hashing para Sharding

```go
type ShardingWithConsistentHash struct {
    ring *ConsistentHashRing
    shards map[string]*Shard
}

func (schr *ShardingWithConsistentHash) GetShard(key string) *Shard {
    shardID := schr.ring.GetNode(key)
    return schr.shards[shardID]
}

func (schr *ShardingWithConsistentHash) AddShard(shardID string) {
    schr.ring.AddNode(shardID)
    // Solo ~1/N de datos se mueven a nuevo shard
    schr.reshard()
}

/*
Resharding:
  Cuando añades novo shard:
  1. Nuevo shard toma datos entre él y el anterior
  2. Datos se migran en background
  3. Mientras: Doble lectura (local + viejo shard)
  4. Después: Limpia viejo shard
*/
```

### Problema: Hotspots

```go
/*
Si una shard recibe MUCHO tráfico:

Shard A: 10,000 req/s  ← Popular
Shard B: 1,000 req/s
Shard C: 1,000 req/s

Shard A está saturado, B y C ociosas

Soluciones:
1. Replicar Shard A (más lecturas)
2. Microsharding: Dividir Shard A en A1, A2, A3
3. Cache tier delante (Redis)
4. Cambiar clave de sharding (si posible)
*/

type HotspotHandler struct {
    primaryShards   map[int]*Shard
    replicaShards   map[int][]*Shard
    accessCounts    map[int]int64
}

func (hh *HotspotHandler) DetectHotspots() []int {
    var hotspots []int
    avgLoad := hh.computeAverageLoad()
    
    for shardID, count := range hh.accessCounts {
        if count > avgLoad*2 {  // 2x promedio
            hotspots = append(hotspots, shardID)
        }
    }
    
    return hotspots
}

func (hh *HotspotHandler) AddReplicas(shardID int) {
    primaryShard := hh.primaryShards[shardID]
    
    // Crear replicas
    for i := 0; i < 3; i++ {
        replica := primaryShard.CreateReplica()
        hh.replicaShards[shardID] = append(hh.replicaShards[shardID], replica)
    }
}
```

---

## 50.10 Fault Tolerance - Tolerancia a Fallos

### Retry con Exponential Backoff

```go
type RetryPolicy struct {
    maxRetries int
    baseDelay  time.Duration
    maxDelay   time.Duration
}

func (rp *RetryPolicy) Execute(fn func() error) error {
    var lastErr error
    
    for attempt := 0; attempt <= rp.maxRetries; attempt++ {
        if err := fn(); err == nil {
            return nil  // Success
        } else {
            lastErr = err
        }
        
        if attempt < rp.maxRetries {
            delay := rp.baseDelay * (1 << uint(attempt))  // 2^attempt
            
            // Cap máximo
            if delay > rp.maxDelay {
                delay = rp.maxDelay
            }
            
            // Jitter para evitar thundering herd
            jitter := time.Duration(rand.Int63n(int64(delay / 2)))
            time.Sleep(delay + jitter)
        }
    }
    
    return lastErr
}

// Uso:
policy := &RetryPolicy{
    maxRetries: 3,
    baseDelay:  100 * time.Millisecond,
    maxDelay:   10 * time.Second,
}

err := policy.Execute(func() error {
    return callRemoteAPI()
})
```

### Circuit Breaker

```go
type CircuitBreakerState int

const (
    CLOSED CircuitBreakerState = iota  // Normal
    OPEN                               // Falla detectada
    HALF_OPEN                          // Recuperación en progreso
)

type CircuitBreaker struct {
    state          CircuitBreakerState
    failureCount   int
    failureThreshold int
    successCount   int
    successThreshold int
    lastFailTime   time.Time
    timeout        time.Duration
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    switch cb.state {
    case CLOSED:
        return cb.callClosed(fn)
    case OPEN:
        return cb.callOpen(fn)
    case HALF_OPEN:
        return cb.callHalfOpen(fn)
    }
    return errors.New("unknown state")
}

func (cb *CircuitBreaker) callClosed(fn func() error) error {
    err := fn()
    
    if err != nil {
        cb.failureCount++
        cb.lastFailTime = time.Now()
        
        if cb.failureCount >= cb.failureThreshold {
            cb.state = OPEN
            return errors.New("circuit breaker opened")
        }
    } else {
        cb.failureCount = 0
    }
    
    return err
}

func (cb *CircuitBreaker) callOpen(fn func() error) error {
    // Si pasó timeout, intenta recuperación
    if time.Since(cb.lastFailTime) > cb.timeout {
        cb.state = HALF_OPEN
        cb.successCount = 0
        return cb.callHalfOpen(fn)
    }
    
    return errors.New("circuit breaker is open")
}

func (cb *CircuitBreaker) callHalfOpen(fn func() error) error {
    err := fn()
    
    if err != nil {
        cb.state = OPEN  // Aún falla, volver a OPEN
        cb.lastFailTime = time.Now()
        return err
    }
    
    cb.successCount++
    if cb.successCount >= cb.successThreshold {
        cb.state = CLOSED  // Recuperado
        cb.failureCount = 0
    }
    
    return nil
}
```

### Bulkheads - Aislamiento de Fallas

```go
/*
Bulkhead pattern: Aislar recursos

Sin bulkheads:
  Todos los requests usan thread pool compartido
  → Un servicio lento consume todos los threads
  → Otros servicios quedan sin threads
  → Cascading failure

Con bulkheads:
  Cada servicio tiene su thread pool
  → Servicio A lento → Pool A saturado
  → Servicio B no afectado
*/

type BulkheadExecutor struct {
    executors map[string]*executor.Executor
}

func (be *BulkheadExecutor) Execute(
    serviceName string,
    fn func() error,
) error {
    executor := be.executors[serviceName]
    
    // Enviar a este executor específico
    // Si está saturado, rechazo (no espero)
    return executor.ExecuteWithTimeout(fn, 5*time.Second)
}

// Uso:
be := &BulkheadExecutor{
    executors: map[string]*executor.Executor{
        "payment":   executor.New(10),  // 10 threads
        "shipping":  executor.New(5),
        "email":     executor.New(3),
    },
}

// Payment lenta no afecta shipping
be.Execute("payment", slowPaymentFn)
be.Execute("shipping", fastShippingFn)
```

### Timeouts

```go
// Timeout simple
func CallWithTimeout(fn func() error, timeout time.Duration) error {
    done := make(chan error, 1)
    
    go func() {
        done <- fn()
    }()
    
    select {
    case err := <-done:
        return err
    case <-time.After(timeout):
        return errors.New("timeout")
    }
}

// Timeout con contexto (Go idiomatic)
func CallWithContext(ctx context.Context, fn func() error) error {
    done := make(chan error, 1)
    
    go func() {
        done <- fn()
    }()
    
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Uso con timeout:
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := CallWithContext(ctx, func() error {
    return remoteService.Call()
})
```

---

## 50.11 Buenas Prácticas y Trade-offs Operacionales

### Monitoreo en Sistemas Distribuidos

```go
type DistributedMetrics struct {
    latencyHistogram  *prometheus.HistogramVec
    errorRate         prometheus.Gauge
    replicationLag    prometheus.Gauge
    quorumErrors      prometheus.Counter
}

func (dm *DistributedMetrics) RecordLatency(operation string, duration time.Duration) {
    dm.latencyHistogram.WithLabelValues(operation).Observe(duration.Seconds())
}

func (dm *DistributedMetrics) RecordError(serviceName string) {
    dm.errorRate.With(prometheus.Labels{"service": serviceName}).Add(1)
}

func (dm *DistributedMetrics) RecordReplicationLag(lagMs int64) {
    dm.replicationLag.Set(float64(lagMs))
}

/*
Métricas críticas en distribuido:

1. Latencia P99, P95, P50
   - P99: Peor 1% (SLA crítico)
   - P95: Peor 5% (usuario frustrado)
   - P50: Promedio (no significa mucho)

2. Error rate por servicio
   - ¿Qué servicio está fallando?

3. Replication lag
   - ¿Cuánto retraso en replicación?

4. Quorum availability
   - ¿Qué % de quórums tienen éxito?

5. Partition detection
   - ¿Cuándo la red se particiona?
*/
```

### Complejidad Operacional

```go
/*
Trade-off: Performance vs Complejidad

Opción 1: Monolítico centralizado
  + Simple de entender
  + Debug fácil
  + Sin problemas de distribución
  - No escala
  - Unavailability si cae

Opción 2: Distribuido
  + Escala a millones de usuarios
  + Alta disponibilidad
  - Complejidad: Múltiples fallos parciales
  - Debugging: Requiere herramientas especializadas
  - Operacionalización: Requiere equipo experto

Regla práctica:
- Si <10k usuarios/día: Monolítico
- Si <100k usuarios/día: Distribuido simple
- Si >100k usuarios/día: Distribuido complejo necesario
*/

// Ejemplo: API response times

// Monolítico:
// P99: 50ms
// P95: 30ms

// Distribuido con 3 réplicas:
// P99: 150ms (suma de latencias + network overhead)
// P95: 100ms

// Distribuido con quórum R/W:
// P99: 200ms (esperar a 2 de 3)
// P95: 150ms

// Tradeoff: Menos disponibilidad, más latencia
```

### Antipatterns Comunes

```go
/*
Antipattern 1: Asumir Consistencia Fuerte en Distribuido

INCORRECTO:
  Write(X=1)
  Sleep(1ms)
  value := Read(X)  // Espera 1 = Y? NO, podría ser 0

CORRECTO:
  Write(X=1)
  // Esperar confirmación de replicación
  WaitForReplication()
  value := Read(X)  // Ahora 1
*/

/*
Antipattern 2: Ignorar CAP Theorem

INCORRECTO:
  "Quiero C, A, y P al 100%"
  → Imposible

CORRECTO:
  "Mi caso de uso es: Bank (CP) o Social (AP)"
  → Elegir y diseñar alrededor
*/

/*
Antipattern 3: Asumir Que Network Es Rápido

INCORRECTO:
  Hacer 1000 RPCs secuenciales
  → 1000 * 10ms = 10 segundos

CORRECTO:
  Batching: 1 RPC con 1000 items
  → 1 * 10ms = 10ms
*/

/*
Antipattern 4: 2PC Sin Fallback

INCORRECTO:
  Do 2PC, si falla → DEADLOCK
  
CORRECTO:
  Sagas con compensations
  O aceptar eventual consistency
*/

// Checklist de diseño distribuido:
// ✓ ¿Elegimos CP o AP según el caso?
// ✓ ¿Tenemos plan para replication lag?
// ✓ ¿Batching en lugar de requests individuales?
// ✓ ¿Circuit breaker para fallos en cascada?
// ✓ ¿Monitoring de P99 latency?
// ✓ ¿Compensations si usamos Sagas?
// ✓ ¿Timeouts en todas las operaciones remotas?
// ✓ ¿Retries con exponential backoff?
```

---

## Ejercicios Progresivos

### Ejercicio 1: Replicación Master-Slave Simple

**Objetivo:** Implementar replicación básica con sincronización

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type DataStore struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (ds *DataStore) Set(key string, value interface{}) {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    ds.data[key] = value
}

func (ds *DataStore) Get(key string) interface{} {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    return ds.data[key]
}

type MasterSlaveSystem struct {
    master   *DataStore
    slaves   []*DataStore
    binlog   []LogEntry
}

type LogEntry struct {
    Operation string
    Key       string
    Value     interface{}
    Timestamp time.Time
}

func (mss *MasterSlaveSystem) Write(key string, value interface{}) {
    mss.master.Set(key, value)
    
    entry := LogEntry{
        Operation: "SET",
        Key:       key,
        Value:     value,
        Timestamp: time.Now(),
    }
    mss.binlog = append(mss.binlog, entry)
    
    // Replicate to slaves
    for _, slave := range mss.slaves {
        go slave.Set(key, value)
    }
}

func (mss *MasterSlaveSystem) Read(key string) interface{} {
    return mss.master.Get(key)
}

// Ejemplo de uso
func main() {
    master := &DataStore{data: make(map[string]interface{})}
    slave1 := &DataStore{data: make(map[string]interface{})}
    slave2 := &DataStore{data: make(map[string]interface{})}
    
    system := &MasterSlaveSystem{
        master:  master,
        slaves:  []*DataStore{slave1, slave2},
        binlog:  []LogEntry{},
    }
    
    system.Write("user:1", "Alice")
    time.Sleep(100 * time.Millisecond)  // Esperar replicación
    
    fmt.Println("Master:", system.master.Get("user:1"))
    fmt.Println("Slave1:", slave1.Get("user:1"))
    fmt.Println("Slave2:", slave2.Get("user:1"))
}
```

**Requisitos:**
- [ ] Master puede escribir
- [ ] Slaves replican cambios
- [ ] Logs se guardan
- [ ] Múltiples escrituras funcionan

---

### Ejercicio 2: Consistent Hashing

**Objetivo:** Implementar particionamiento de datos sin rehashing masivo

```go
package main

import (
    "fmt"
    "hash/fnv"
    "sort"
)

func hashFunc(key string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(key))
    return h.Sum32()
}

type Shard struct {
    id   string
    data map[string]interface{}
}

type ConsistentHash struct {
    nodes      map[uint32]string  // hash -> node_id
    sortedKeys []uint32
    replicas   int
}

func NewConsistentHash(replicas int) *ConsistentHash {
    return &ConsistentHash{
        nodes:    make(map[uint32]string),
        replicas: replicas,
    }
}

func (ch *ConsistentHash) AddNode(nodeID string) {
    for i := 0; i < ch.replicas; i++ {
        key := fmt.Sprintf("%s:%d", nodeID, i)
        hash := hashFunc(key)
        ch.nodes[hash] = nodeID
        ch.sortedKeys = append(ch.sortedKeys, hash)
    }
    sort.Slice(ch.sortedKeys, func(i, j int) bool {
        return ch.sortedKeys[i] < ch.sortedKeys[j]
    })
}

func (ch *ConsistentHash) GetNode(key string) string {
    if len(ch.nodes) == 0 {
        return ""
    }
    
    hash := hashFunc(key)
    idx := sort.Search(len(ch.sortedKeys), func(i int) bool {
        return ch.sortedKeys[i] >= hash
    })
    
    if idx == len(ch.sortedKeys) {
        idx = 0
    }
    
    return ch.nodes[ch.sortedKeys[idx]]
}

func main() {
    ch := NewConsistentHash(3)  // 3 virtual nodes
    
    ch.AddNode("server-1")
    ch.AddNode("server-2")
    ch.AddNode("server-3")
    
    // Distribuir claves
    keys := []string{"user:1", "user:2", "user:3", "cache:x", "cache:y"}
    for _, key := range keys {
        shard := ch.GetNode(key)
        fmt.Printf("%s -> %s\n", key, shard)
    }
}
```

**Requisitos:**
- [ ] Agregar nodos sin rehashear todo
- [ ] Distribución uniforme
- [ ] Virtual nodes para balance

---

### Ejercicio 3: Implementación Simple de Raft

**Objetivo:** Implementar leader election y log replication

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type RaftState int

const (
    FOLLOWER RaftState = iota
    CANDIDATE
    LEADER
)

type LogEntry struct {
    Term    int64
    Index   int64
    Command interface{}
}

type RaftNode struct {
    id           string
    state        RaftState
    currentTerm  int64
    votedFor     string
    log          []LogEntry
    commitIndex  int64
    
    peers        []*RaftNode
    mu           sync.RWMutex
}

func NewRaftNode(id string) *RaftNode {
    return &RaftNode{
        id:    id,
        state: FOLLOWER,
        log:   []LogEntry{},
    }
}

func (rn *RaftNode) Append(entry LogEntry) bool {
    rn.mu.Lock()
    defer rn.mu.Unlock()
    
    if rn.state != LEADER {
        return false
    }
    
    entry.Term = rn.currentTerm
    entry.Index = int64(len(rn.log))
    rn.log = append(rn.log, entry)
    
    rn.replicate()
    return true
}

func (rn *RaftNode) replicate() {
    // Enviar log entries a peers
    for _, peer := range rn.peers {
        go peer.Sync(rn.log)
    }
}

func (rn *RaftNode) Sync(entries []LogEntry) {
    rn.mu.Lock()
    defer rn.mu.Unlock()
    
    for _, entry := range entries {
        if entry.Index >= int64(len(rn.log)) {
            rn.log = append(rn.log, entry)
        }
    }
}

func (rn *RaftNode) ElectionTimeout() time.Duration {
    return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

func (rn *RaftNode) StartElection() {
    rn.mu.Lock()
    rn.state = CANDIDATE
    rn.currentTerm++
    rn.votedFor = rn.id
    rn.mu.Unlock()
    
    votes := 1
    for _, peer := range rn.peers {
        if peer.RequestVote(rn.currentTerm, rn.id) {
            votes++
        }
    }
    
    rn.mu.Lock()
    if votes > len(rn.peers)/2 && rn.state == CANDIDATE {
        rn.state = LEADER
        fmt.Printf("%s elected LEADER\n", rn.id)
    } else {
        rn.state = FOLLOWER
    }
    rn.mu.Unlock()
}

func (rn *RaftNode) RequestVote(term int64, candidateID string) bool {
    rn.mu.Lock()
    defer rn.mu.Unlock()
    
    if term > rn.currentTerm {
        rn.currentTerm = term
        rn.votedFor = ""
    }
    
    if term == rn.currentTerm && (rn.votedFor == "" || rn.votedFor == candidateID) {
        rn.votedFor = candidateID
        return true
    }
    
    return false
}

func main() {
    nodes := make([]*RaftNode, 3)
    for i := 0; i < 3; i++ {
        nodes[i] = NewRaftNode(fmt.Sprintf("node-%d", i))
    }
    
    // Conectar peers
    for i := 0; i < 3; i++ {
        nodes[i].peers = append(nodes[i].peers, nodes...)
    }
    
    // Iniciar elecciones
    for _, node := range nodes {
        go node.StartElection()
    }
    
    time.Sleep(1 * time.Second)
    
    for _, node := range nodes {
        fmt.Printf("%s: state=%v\n", node.id, node.state)
    }
}
```

**Requisitos:**
- [ ] Leader election funciona
- [ ] Log replication entre nodos
- [ ] Manejo de términos

---

### Ejercicio 4: Event Sourcing

**Objetivo:** Implementar source of truth basado en eventos

```go
package main

import (
    "fmt"
    "time"
)

type Event struct {
    ID        string
    Type      string
    Timestamp time.Time
    Data      map[string]interface{}
}

type EventStore struct {
    events []Event
}

func (es *EventStore) Append(eventType string, data map[string]interface{}) {
    es.events = append(es.events, Event{
        Type:      eventType,
        Timestamp: time.Now(),
        Data:      data,
    })
}

type Account struct {
    ID      string
    Balance float64
    Version int64
}

func (es *EventStore) ReplayAccount(accountID string) Account {
    acc := Account{ID: accountID}
    
    for _, event := range es.events {
        if event.Data["account_id"] != accountID {
            continue
        }
        
        switch event.Type {
        case "deposit":
            acc.Balance += event.Data["amount"].(float64)
        case "withdraw":
            acc.Balance -= event.Data["amount"].(float64)
        case "interest":
            acc.Balance *= (1 + event.Data["rate"].(float64))
        }
        acc.Version++
    }
    
    return acc
}

func main() {
    es := &EventStore{}
    
    // Simular transacciones
    es.Append("deposit", map[string]interface{}{
        "account_id": "acc1",
        "amount":     1000.0,
    })
    
    es.Append("withdraw", map[string]interface{}{
        "account_id": "acc1",
        "amount":     200.0,
    })
    
    es.Append("interest", map[string]interface{}{
        "account_id": "acc1",
        "rate":       0.01,  // 1%
    })
    
    account := es.ReplayAccount("acc1")
    fmt.Printf("Account %s: Balance = %.2f (v%d)\n", 
        account.ID, account.Balance, account.Version)
}
```

**Requisitos:**
- [ ] Guardar eventos
- [ ] Replay de estado
- [ ] Auditoría completa

---

### Ejercicio 5: Distributed System Completo

**Objetivo:** Integrar replicación, consenso, y fault tolerance

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type DistributedStore struct {
    nodes        []*RaftNode
    commitLog    []LogEntry
    currentLeader *RaftNode
    mu           sync.RWMutex
}

func (ds *DistributedStore) Write(key string, value interface{}) error {
    if ds.currentLeader.state != LEADER {
        return fmt.Errorf("not leader")
    }
    
    entry := LogEntry{
        Command: fmt.Sprintf("SET %s=%v", key, value),
    }
    
    if ds.currentLeader.Append(entry) {
        ds.commitLog = append(ds.commitLog, entry)
        return nil
    }
    
    return fmt.Errorf("replication failed")
}

func (ds *DistributedStore) HandleLeaderFailure() {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    
    // Trigger new election
    for _, node := range ds.nodes {
        if node != ds.currentLeader {
            node.StartElection()
        }
    }
    
    // Encontrar nuevo líder
    time.Sleep(500 * time.Millisecond)
    for _, node := range ds.nodes {
        if node.state == LEADER {
            ds.currentLeader = node
            break
        }
    }
}

func main() {
    // Crear 5 nodos
    nodes := make([]*RaftNode, 5)
    for i := 0; i < 5; i++ {
        nodes[i] = NewRaftNode(fmt.Sprintf("node-%d", i))
    }
    
    // Conectar
    for i := 0; i < 5; i++ {
        nodes[i].peers = append(nodes[i].peers, nodes...)
    }
    
    // Eleger líder
    nodes[0].StartElection()
    time.Sleep(500 * time.Millisecond)
    
    store := &DistributedStore{
        nodes:        nodes,
        currentLeader: nodes[0],
    }
    
    // Escribir data
    _ = store.Write("key1", "value1")
    fmt.Println("Written to leader")
    
    // Simular fallo de líder
    fmt.Println("Leader died!")
    store.HandleLeaderFailure()
    
    fmt.Printf("New leader: %s\n", store.currentLeader.id)
}
```

**Requisitos:**
- [ ] Replicación entre múltiples nodos
- [ ] Leader failover
- [ ] Recuperación automática
- [ ] Data consistency

---

## Resumen de Conceptos Clave

| Concepto | Aplicación | Trade-off |
|----------|-----------|----------|
| **CAP** | Elegir entre C, A, P | No puedes tener los 3 |
| **Raft** | Consenso y leader election | Simple pero requiere líder |
| **2PC** | Transacciones atómicas | Lento y puede bloquear |
| **Sagas** | Transacciones distribuidas | Eventual consistency |
| **Event Sourcing** | Auditoría completa | Storage overhead |
| **Consistent Hashing** | Sharding escalable | Resharding en background |
| **Circuit Breaker** | Fault tolerance | Latencia aumentada |
| **Retry + Backoff** | Resiliencia | Puede enmascarar problemas |

---

## Referencias y Lecturas Recomendadas

1. **"Designing Data-Intensive Applications"** - Martin Kleppmann
2. **"The Art of Distributed Systems"** - Leslie Lamport
3. **Etcd Documentation** - https://etcd.io
4. **Raft Consensus Algorithm** - https://raft.github.io
5. **Go Distributed Systems Libraries:**
   - github.com/hashicorp/raft
   - github.com/etcd-io/raft
   - github.com/go-zk/zk (Zookeeper client)

---

**Fin del Capítulo 50: Distributed Systems - Sistemas Distribuidos a Escala**

*Este capítulo proporciona una base sólida para entender sistemas distribuidos con Go, desde teoría fundamental hasta implementaciones prácticas. Los ejercicios progresivos permiten practicar conceptos en contextos realistas.*

---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/50-distributed-systems/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/50-distributed-systems):

```bash
cd examples/50-distributed-systems
go run .
```
