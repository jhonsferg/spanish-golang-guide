# Capítulo 45: Clean architecture - Diseño de sistemas escalables

## Introducción

Clean Architecture es una filosofía de diseño que enfatiza la independencia de frameworks, bases de datos y detalles de implementación. Propuesta por Robert C. Martin ("Uncle Bob"), esta arquitectura facilita la creación de sistemas mantenibles, testables y escalables.

En Go, la implementación de Clean Architecture es natural gracias al sistema de interfaces implícitas y la flexibilidad del lenguaje. Este capítulo te enseñará cómo estructurar aplicaciones Go para máxima escalabilidad.

---

## 45.1 Principios Fundamentales de Clean Architecture

### 45.1.1 Capas de la Arquitectura

Clean Architecture organiza el código en capas concéntricas, donde cada capa tiene responsabilidades específicas:

```
┌─────────────────────────────────────┐
│    Frameworks & Drivers Layer       │ (HTTP, BD, External Services)
│ ┌───────────────────────────────────┤
│ │  Interface Adapters Layer         │ (Controllers, Presenters, Gateways)
│ │ ┌─────────────────────────────────┤
│ │ │  Use Cases Layer                │ (Application Logic, Orchestration)
│ │ │ ┌───────────────────────────────┤
│ │ │ │  Entities Layer               │ (Business Logic, Domain Models)
│ │ │ │ (Core rules, value objects)   │
│ │ │ └───────────────────────────────┘
│ │ └─────────────────────────────────┘
│ └───────────────────────────────────┘
└─────────────────────────────────────┘
```

### 45.1.2 Regla de Dependencia

**La regla fundamental**: Las dependencias siempre apuntan hacia adentro. El código interno nunca debe conocer el código externo.

```
- Entities NO conocen Use Cases
- Use Cases NO conocen Interface Adapters
- Interface Adapters NO conocen Frameworks
```

### 45.1.3 Independencia de Frameworks

La lógica de negocio es agnóstica respecto a:

- Framework web (gin, echo, fiber)
- Base de datos (SQL, NoSQL)
- Frameworks de caché
- Sistemas de colas
- Servicios externos

### 45.1.4 Testabilidad

Cada capa se prueba independientemente:

```go
// Test de entidad sin dependencias externas
func TestUserEntity(t *testing.T) {
    user := domain.NewUser("john@example.com", "John")
    if !user.IsValid() {
        t.Fatal("User should be valid")
    }
}

// Test de use case con mocks
func TestCreateUserUseCase(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Save").Return(nil)
    
    useCase := usecases.NewCreateUserUseCase(mockRepo)
    err := useCase.Execute(&CreateUserRequest{})
    assert.NoError(t, err)
}
```

### 45.1.5 Mantenibilidad a Largo Plazo

Los beneficios crecen con el tiempo:

- **Fácil localización de bugs**: Cada capa tiene responsabilidades claras
- **Cambios sin fricción**: Modificar un framework no afecta lógica de negocio
- **Onboarding de nuevos desarrolladores**: Estructura predecible
- **Refactoring seguro**: Tests protegen cambios

---

## 45.2 Entities Layer - La Capa Más Interna

### 45.2.1 Definición y Responsabilidades

La capa Entities contiene las reglas de negocio fundamentales:

```go
package domain

import (
    "fmt"
    "regexp"
    "time"
)

// User es una entidad central del dominio
type User struct {
    ID        string
    Email     string
    Name      string
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
    Active    bool
}

// NewUser crea un usuario validado
func NewUser(email, name string) (*User, error) {
    if !isValidEmail(email) {
        return nil, fmt.Errorf("email inválido: %s", email)
    }
    
    if len(name) < 3 {
        return nil, fmt.Errorf("nombre debe tener al menos 3 caracteres")
    }
    
    return &User{
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Active:    true,
    }, nil
}

// IsValid verifica si la entidad es válida
func (u *User) IsValid() bool {
    return isValidEmail(u.Email) && 
           len(u.Name) >= 3 && 
           u.Email != "" &&
           u.Active
}

// ChangeEmail cambia el email validando
func (u *User) ChangeEmail(newEmail string) error {
    if !isValidEmail(newEmail) {
        return fmt.Errorf("nuevo email inválido")
    }
    u.Email = newEmail
    u.UpdatedAt = time.Now()
    return nil
}

// isValidEmail valida formato de email
func isValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}
```

### 45.2.2 Value Objects

Value Objects son objetos que se identifican por su valor, no por su identidad:

```go
package domain

import "fmt"

// Email es un value object
type Email struct {
    value string
}

func NewEmail(email string) (*Email, error) {
    if !isValidEmail(email) {
        return nil, fmt.Errorf("email inválido: %s", email)
    }
    return &Email{value: email}, nil
}

func (e *Email) Value() string {
    return e.value
}

func (e *Email) String() string {
    return e.value
}

// Dos emails con mismo valor son iguales
func (e *Email) Equals(other *Email) bool {
    return e.value == other.value
}

// Password es un value object
type Password struct {
    hash string
}

func NewPassword(plainText string) (*Password, error) {
    if len(plainText) < 8 {
        return nil, fmt.Errorf("contraseña debe tener mínimo 8 caracteres")
    }
    // En producción: usar bcrypt
    hash := hashPassword(plainText)
    return &Password{hash: hash}, nil
}

func (p *Password) Verify(plainText string) bool {
    return hashPassword(plainText) == p.hash
}

func hashPassword(text string) string {
    // Simplificado para ejemplo
    return fmt.Sprintf("hashed_%s", text)
}
```

### 45.2.3 Agregados (Aggregates)

Agregados son conjuntos de entidades tratadas como una unidad:

```go
package domain

// Order es un agregado raíz
type Order struct {
    ID       string
    UserID   string
    Items    []*OrderItem
    Status   OrderStatus
    Total    float64
    CreatedAt time.Time
}

type OrderItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "PENDING"
    OrderStatusConfirmed OrderStatus = "CONFIRMED"
    OrderStatusShipped   OrderStatus = "SHIPPED"
    OrderStatusDelivered OrderStatus = "DELIVERED"
)

// AddItem agrega un producto al pedido
func (o *Order) AddItem(productID string, quantity int, price float64) error {
    if o.Status != OrderStatusPending {
        return fmt.Errorf("no se puede agregar items a pedido %s", o.Status)
    }
    
    o.Items = append(o.Items, &OrderItem{
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
    })
    
    o.recalculateTotal()
    return nil
}

// Confirm transiciona a estado confirmado
func (o *Order) Confirm() error {
    if o.Status != OrderStatusPending {
        return fmt.Errorf("solo pedidos PENDING pueden confirmarse")
    }
    if len(o.Items) == 0 {
        return fmt.Errorf("pedido debe tener items")
    }
    o.Status = OrderStatusConfirmed
    return nil
}

func (o *Order) recalculateTotal() {
    o.Total = 0
    for _, item := range o.Items {
        o.Total += float64(item.Quantity) * item.Price
    }
}
```

---

## 45.3 Use Cases Layer - Lógica de Aplicación

### 45.3.1 Responsabilidades de Use Cases

Use Cases orquestan la lógica de aplicación:

```go
package usecases

import "errors"

// CreateUserRequest es el request de entrada
type CreateUserRequest struct {
    Email    string `json:"email"`
    Name     string `json:"name"`
    Password string `json:"password"`
}

// CreateUserResponse es la respuesta
type CreateUserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

// CreateUserUseCase implementa la lógica de crear usuario
type CreateUserUseCase struct {
    userRepo     UserRepository
    emailService EmailService
}

// Interfaces que el use case requiere
type UserRepository interface {
    Save(user *domain.User) error
    FindByEmail(email string) (*domain.User, error)
}

type EmailService interface {
    SendWelcomeEmail(email, name string) error
}

// NewCreateUserUseCase crea el use case
func NewCreateUserUseCase(
    userRepo UserRepository,
    emailService EmailService,
) *CreateUserUseCase {
    return &CreateUserUseCase{
        userRepo:     userRepo,
        emailService: emailService,
    }
}

// Execute ejecuta el use case
func (uc *CreateUserUseCase) Execute(
    req *CreateUserRequest,
) (*CreateUserResponse, error) {
    // Validaciones
    if req.Email == "" || req.Name == "" {
        return nil, errors.New("email y nombre requeridos")
    }
    
    // Verificar que no exista
    existing, _ := uc.userRepo.FindByEmail(req.Email)
    if existing != nil {
        return nil, errors.New("usuario ya existe")
    }
    
    // Crear entidad
    user, err := domain.NewUser(req.Email, req.Name)
    if err != nil {
        return nil, err
    }
    
    // Guardar
    if err := uc.userRepo.Save(user); err != nil {
        return nil, err
    }
    
    // Enviar email (sin bloquear)
    go func() {
        uc.emailService.SendWelcomeEmail(user.Email, user.Name)
    }()
    
    return &CreateUserResponse{
        ID:    user.ID,
        Email: user.Email,
        Name:  user.Name,
    }, nil
}
```

### 45.3.2 Patrones de Use Cases

```go
package usecases

// Use Case simple: Lectura
type GetUserByIDUseCase struct {
    userRepo UserRepository
}

func (uc *GetUserByIDUseCase) Execute(userID string) (*domain.User, error) {
    return uc.userRepo.FindByID(userID)
}

// Use Case complejo: Transacción
type TransferMoneyUseCase struct {
    accountRepo   AccountRepository
    transactionSvc TransactionService
}

type TransferRequest struct {
    FromAccountID string
    ToAccountID   string
    Amount        float64
}

func (uc *TransferMoneyUseCase) Execute(req *TransferRequest) error {
    // Obtener cuentas
    from, err := uc.accountRepo.FindByID(req.FromAccountID)
    if err != nil {
        return err
    }
    
    to, err := uc.accountRepo.FindByID(req.ToAccountID)
    if err != nil {
        return err
    }
    
    // Validar
    if from.Balance < req.Amount {
        return errors.New("saldo insuficiente")
    }
    
    // Ejecutar transacción
    txn := uc.transactionSvc.Begin()
    
    from.Withdraw(req.Amount)
    to.Deposit(req.Amount)
    
    if err := txn.Commit(from, to); err != nil {
        txn.Rollback()
        return err
    }
    
    return nil
}
```

### 45.3.3 Inyección de Dependencias

```go
package usecases

import "go.uber.org/fx"

// Módulo de inyección de dependencias
var Module = fx.Options(
    fx.Provide(
        NewCreateUserUseCase,
        NewGetUserUseCase,
        NewUpdateUserUseCase,
        NewListUsersUseCase,
    ),
)

// El módulo inyecta automáticamente dependencias
type UseCaseContainer struct {
    CreateUser *CreateUserUseCase
    GetUser    *GetUserUseCase
    UpdateUser *UpdateUserUseCase
    ListUsers  *ListUsersUseCase
}
```

---

## 45.4 Interface Adapters - Adaptadores

### 45.4.1 Controllers (HTTP Handlers)

Controllers adaptan requests HTTP a use cases:

```go
package adapters

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type UserController struct {
    createUserUseCase *usecases.CreateUserUseCase
    getUserUseCase    *usecases.GetUserUseCase
}

// CreateUser maneja POST /users
func (uc *UserController) CreateUser(c *gin.Context) {
    var req usecases.CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    resp, err := uc.createUserUseCase.Execute(&req)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, resp)
}

// GetUser maneja GET /users/:id
func (uc *UserController) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := uc.getUserUseCase.Execute(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

### 45.4.2 Presenters (Response Formatters)

Presenters transforman datos del dominio a DTOs:

```go
package adapters

type UserPresenter struct{}

// ToDTO transforma entidad a DTO
func (p *UserPresenter) ToDTO(user *domain.User) UserDTO {
    return UserDTO{
        ID:        user.ID,
        Email:     user.Email,
        Name:      user.Name,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
        Active:    user.Active,
    }
}

// ToList transforma slice de usuarios
func (p *UserPresenter) ToList(users []*domain.User) []UserDTO {
    dtos := make([]UserDTO, len(users))
    for i, u := range users {
        dtos[i] = p.ToDTO(u)
    }
    return dtos
}

type UserDTO struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
    Active    bool   `json:"active"`
}
```

### 45.4.3 Gateways (External Services)

Gateways adaptan servicios externos:

```go
package adapters

import (
    "net/http"
    "encoding/json"
)

type EmailGateway struct {
    apiURL string
    client *http.Client
}

type SendEmailRequest struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func (g *EmailGateway) SendWelcomeEmail(email, name string) error {
    req := SendEmailRequest{
        To:      email,
        Subject: "Bienvenido a nuestro servicio",
        Body:    "Hola " + name,
    }
    
    body, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST", g.apiURL+"/send", 
        bytes.NewBuffer(body))
    
    resp, err := g.client.Do(httpReq)
    if err != nil {
        return err
    }
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("error enviando email: %d", resp.StatusCode)
    }
    
    return nil
}
```

---

## 45.5 Frameworks & Drivers - Capa Externa

### 45.5.1 Implementación de Repositorios

Repositorios concretos usan base de datos específica:

```go
package infrastructure

import (
    "database/sql"
    "fmt"
)

type SQLUserRepository struct {
    db *sql.DB
}

func (r *SQLUserRepository) Save(user *domain.User) error {
    query := `
        INSERT INTO users (id, email, name, password, created_at, updated_at, active)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    err := r.db.QueryRow(
        query,
        user.ID,
        user.Email,
        user.Name,
        user.Password,
        user.CreatedAt,
        user.UpdatedAt,
        user.Active,
    ).Scan()
    
    return err
}

func (r *SQLUserRepository) FindByEmail(email string) (*domain.User, error) {
    query := `SELECT id, email, name, password, created_at, updated_at, active FROM users WHERE email = $1`
    
    var user domain.User
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.Password,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.Active,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

### 45.5.2 Configuración de Frameworks

```go
package infrastructure

import (
    "github.com/gin-gonic/gin"
)

func SetupRoutes(
    engine *gin.Engine,
    userController *adapters.UserController,
) {
    // Grupo de usuarios
    users := engine.Group("/api/v1/users")
    {
        users.POST("", userController.CreateUser)
        users.GET("/:id", userController.GetUser)
        users.PUT("/:id", userController.UpdateUser)
        users.DELETE("/:id", userController.DeleteUser)
    }
}

func SetupDatabase(connString string) (*sql.DB, error) {
    db, err := sql.Open("postgres", connString)
    if err != nil {
        return nil, err
    }
    
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

### 45.5.3 Inyección en Main

```go
package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

func main() {
    // Configurar BD
    db, _ := infrastructure.SetupDatabase("postgres://user:pass@localhost/db")
    defer db.Close()
    
    // Repositorios concretos
    userRepo := infrastructure.NewSQLUserRepository(db)
    
    // Servicios
    emailGateway := adapters.NewEmailGateway("https://api.email.com")
    
    // Use cases
    createUserUC := usecases.NewCreateUserUseCase(userRepo, emailGateway)
    getUserUC := usecases.NewGetUserUseCase(userRepo)
    
    // Controllers
    userController := adapters.NewUserController(createUserUC, getUserUC)
    
    // Framework HTTP
    engine := gin.Default()
    infrastructure.SetupRoutes(engine, userController)
    
    engine.Run(":8080")
}
```

---

## 45.6 Dependency Inversion - Principio de Inversión

### 45.6.1 Problemas sin Inversión

```go
// ❌ MALO: Acoplamiento directo
package services

import "myapp/database"

type UserService struct {
    db *database.PostgresDB
}

// UserService está acoplado a PostgreSQL específicamente
```

### 45.6.2 Solución: Abstraer con Interfaces

```go
// ✅ BIEN: Abstraído con interfaces
package services

// UserRepository es la abstracción
type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id string) (*domain.User, error)
    Delete(id string) error
}

type UserService struct {
    repo UserRepository  // Depende de abstracción, no implementación
}

// Diferentes implementaciones pueden usarse
type PostgresUserRepository struct{ /* ... */ }
type MongoUserRepository struct{ /* ... */ }
type InMemoryUserRepository struct{ /* ... */ }
```

### 45.6.3 Inyección de Dependencias

```go
package main

// El main decide qué implementación usar
func main() {
    var repo services.UserRepository
    
    // Usa PostgreSQL en producción
    if os.Getenv("ENV") == "prod" {
        repo = infrastructure.NewPostgresUserRepository(db)
    } else {
        // Usa in-memory en tests
        repo = infrastructure.NewInMemoryUserRepository()
    }
    
    userService := services.NewUserService(repo)
    // userService funciona igual con cualquier implementación
}
```

### 45.6.4 Ventajas de Dependency Inversion

```go
// Test fácil
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Save(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func TestUserService(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Save", mock.Anything).Return(nil)
    
    service := services.NewUserService(mockRepo)
    err := service.CreateUser("john@example.com", "John")
    
    assert.NoError(t, err)
    mockRepo.AssertCalled(t, "Save", mock.Anything)
}
```

---

## 45.7 Domain-Driven Design (DDD)

### 45.7.1 Agregados y Raíces de Agregado

```go
package domain

// UserAggregate es la raíz del agregado
type UserAggregate struct {
    ID       string
    Email    *Email
    Name     string
    Profile  *UserProfile
    Roles    []*Role
}

// UserProfile es parte del agregado pero no es raíz
type UserProfile struct {
    Bio           string
    Avatar        string
    PhoneVerified bool
    PhoneNumber   string
}

// Role está en el agregado
type Role struct {
    ID    string
    Name  string
    Perms []string
}

// La raíz del agregado protege la invariante
func (ua *UserAggregate) AddRole(role *Role) error {
    // Validación: máximo 5 roles
    if len(ua.Roles) >= 5 {
        return fmt.Errorf("usuario no puede tener más de 5 roles")
    }
    
    // Validación: no roles duplicados
    for _, r := range ua.Roles {
        if r.ID == role.ID {
            return fmt.Errorf("rol ya asignado")
        }
    }
    
    ua.Roles = append(ua.Roles, role)
    return nil
}
```

### 45.7.2 Eventos de Dominio

```go
package domain

import "time"

// DomainEvent es la interfaz base
type DomainEvent interface {
    AggregateID() string
    EventType() string
    OccurredAt() time.Time
}

// UserCreatedEvent es un evento de dominio
type UserCreatedEvent struct {
    userID    string
    email     string
    name      string
    occurredAt time.Time
}

func (e *UserCreatedEvent) AggregateID() string {
    return e.userID
}

func (e *UserCreatedEvent) EventType() string {
    return "UserCreated"
}

func (e *UserCreatedEvent) OccurredAt() time.Time {
    return e.occurredAt
}

// El agregado emite eventos
type User struct {
    ID     string
    Email  string
    events []DomainEvent  // Eventos no persistidos
}

func (u *User) Create(email, name string) {
    u.ID = generateID()
    u.Email = email
    
    u.events = append(u.events, &UserCreatedEvent{
        userID:    u.ID,
        email:     email,
        name:      name,
        occurredAt: time.Now(),
    })
}

// Obtener y limpiar eventos
func (u *User) GetEvents() []DomainEvent {
    return u.events
}

func (u *User) ClearEvents() {
    u.events = nil
}
```

### 45.7.3 Repositorios de Agregados

```go
package domain

// UserRepository maneja agregados User completos
type UserRepository interface {
    Save(user *User) error
    FindByID(id string) (*User, error)
    FindByEmail(email string) (*User, error)
}

// La implementación perciste el agregado completo
type SQLUserAggregateRepository struct {
    db *sql.DB
}

func (r *SQLUserAggregateRepository) Save(user *User) error {
    // Guardar la raíz
    err := r.saveUser(user)
    if err != nil {
        return err
    }
    
    // Guardar eventos de dominio
    for _, event := range user.GetEvents() {
        r.saveEvent(event)
    }
    
    user.ClearEvents()
    return nil
}
```

### 45.7.4 Bounded Contexts

```go
package usercontext

// Un contexto delimitado para usuarios
type User struct {
    ID    string
    Email string
}

package ordercontext

// Otro contexto para órdenes - referencia a User diferente
type User struct {
    ID   string
    Name string
}

// Traducción entre contextos
package adapters

type UserTranslator struct{}

func (t *UserTranslator) UserFromOrderContextToUserContext(
    orderUser *ordercontext.User,
) *usercontext.User {
    return &usercontext.User{
        ID:    orderUser.ID,
        Email: lookupUserEmail(orderUser.ID),
    }
}
```

---

## 45.8 Separation of Concerns (Separación de Responsabilidades)

### 45.8.1 Single Responsibility Principle (SRP)

```go
// ❌ MALO: Múltiples responsabilidades
type UserManager struct {
    db *sql.DB
}

func (um *UserManager) CreateUser(name string) error {
    // Validación
    if len(name) < 3 {
        return errors.New("name too short")
    }
    
    // Guardar en BD
    query := "INSERT INTO users (name) VALUES ($1)"
    _, err := um.db.Exec(query, name)
    
    // Enviar email
    sendEmailTo("admin@example.com", "New user: "+name)
    
    return err
}

// ✅ BIEN: Cada clase una responsabilidad
type UserValidator struct{}

func (v *UserValidator) Validate(name string) error {
    if len(name) < 3 {
        return errors.New("nombre muy corto")
    }
    return nil
}

type UserRepository struct{ db *sql.DB }

func (r *UserRepository) Save(user *domain.User) error {
    query := "INSERT INTO users (name) VALUES ($1)"
    _, err := r.db.Exec(query, user.Name)
    return err
}

type NotificationService struct{}

func (ns *NotificationService) NotifyNewUser(user *domain.User) error {
    return sendEmailTo("admin@example.com", "Nuevo usuario: "+user.Name)
}
```

### 45.8.2 Cohesión Alta

```go
// ❌ BAJO: Elementos no relacionados juntos
package utils

func SendEmail(to, body string) error { /* ... */ }
func ValidateEmail(email string) bool { /* ... */ }
func GenerateID() string { /* ... */ }
func CalculateTax(amount float64) float64 { /* ... */ }

// ✅ ALTO: Elementos relacionados juntos
package email

func Send(to, body string) error { /* ... */ }
func ValidateAddress(email string) bool { /* ... */ }

package tax

func Calculate(amount float64) float64 { /* ... */ }

package identity

func GenerateID() string { /* ... */ }
```

### 45.8.3 Beneficios para Testing

```go
// Separación facilita tests
type OrderService struct {
    repo   OrderRepository
    tax    TaxCalculator
    payment PaymentProcessor
}

func TestOrderService_CreateOrder(t *testing.T) {
    // Mock cada dependencia independientemente
    mockRepo := &MockOrderRepository{}
    mockTax := &MockTaxCalculator{}
    mockPayment := &MockPaymentProcessor{}
    
    mockTax.On("Calculate").Return(10.0)
    mockPayment.On("Process").Return(nil)
    
    svc := OrderService{mockRepo, mockTax, mockPayment}
    
    // Test solo la lógica de OrderService
    err := svc.CreateOrder(&OrderRequest{Amount: 100})
    
    assert.NoError(t, err)
    mockTax.AssertCalled(t, "Calculate")
    mockPayment.AssertCalled(t, "Process")
}
```

---

## 45.9 Hexagonal Architecture (Ports & Adapters)

### 45.9.1 Concepto de Puertos

```
Exterior del Hexágono (Adapters)
┌─────────────────────────────────────┐
│                                     │
│   HTTP  │  gRPC  │  Message Queue   │
│   ───────────────────────────────   │
│                                     │
│   ◀─── Ports (Interfaces) ─────────  │
│                                     │
│  ┌───────────────────────────────┐  │
│  │   Application Core            │  │
│  │   (Business Logic)            │  │
│  └───────────────────────────────┘  │
│                                     │
│   ──────── Ports (Interfaces) ──────▶│
│                                     │
│  Database │ Cache  │ External APIs  │
│                                     │
└─────────────────────────────────────┘
Interior del Hexágono
```

### 45.9.2 Puertos de Entrada

```go
package application

// Puerto de entrada: HTTP
type UserInputPort interface {
    CreateUser(req *CreateUserRequest) (*CreateUserResponse, error)
    GetUser(id string) (*UserResponse, error)
}

// Puerto de entrada: Mensaje
type UserCommandPort interface {
    OnUserCreatedMessage(event []byte) error
}

// Implementación del puerto
type UserController struct {
    createUC *CreateUserUseCase
    getUC    *GetUserUseCase
}

func (uc *UserController) CreateUser(
    req *CreateUserRequest,
) (*CreateUserResponse, error) {
    return uc.createUC.Execute(req)
}

func (uc *UserController) GetUser(id string) (*UserResponse, error) {
    return uc.getUC.Execute(id)
}
```

### 45.9.3 Puertos de Salida

```go
package application

// Puertos de salida (dependencias)
type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id string) (*domain.User, error)
}

type EmailService interface {
    SendWelcomeEmail(email, name string) error
}

type CacheService interface {
    Set(key string, value interface{}, ttl time.Duration) error
    Get(key string) (interface{}, error)
}

// Adapters concretos pueden conectarse a diferentes tecnologías
type PostgresUserRepository struct{ /* SQL */ }
type MongoUserRepository struct{ /* NoSQL */ }
type RedisUserRepository struct{ /* Cache */ }

type GmailEmailService struct{ /* Google */ }
type SendgridEmailService struct{ /* Sendgrid */ }
type LocalEmailService struct{ /* Test */ }
```

### 45.9.4 Configuración Hexagonal

```go
package main

func main() {
    // Adapters de entrada (Driven)
    httpAdapter := adapters.NewHTTPAdapter()
    grpcAdapter := adapters.NewGrpcAdapter()
    
    // Adapters de salida (Driving)
    userRepo := infrastructure.NewPostgresUserRepository()
    emailSvc := infrastructure.NewGmailEmailService()
    cacheSvc := infrastructure.NewRedisCacheService()
    
    // Aplicación (core)
    createUserUC := usecases.NewCreateUserUseCase(userRepo, emailSvc)
    getUserUC := usecases.NewGetUserUseCase(userRepo, cacheSvc)
    
    // Conectar puertos al core
    httpAdapter.RegisterUseCase(createUserUC, getUserUC)
    grpcAdapter.RegisterUseCase(createUserUC, getUserUC)
    
    // Iniciar servidores
    go httpAdapter.Start()
    go grpcAdapter.Start()
}
```

---

## 45.10 Project Structure - Organización del Proyecto

### 45.10.1 Layout Estándar

```
myapp/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── domain/
│   │   ├── user.go              # Entidades
│   │   ├── order.go
│   │   └── events.go
│   ├── application/
│   │   ├── usecases/
│   │   │   ├── create_user.go
│   │   │   ├── get_user.go
│   │   │   └── transfer_money.go
│   │   └── ports/
│   │       ├── input.go
│   │       └── output.go
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── user_controller.go
│   │   │   ├── order_controller.go
│   │   │   └── router.go
│   │   ├── grpc/
│   │   │   └── user_service.go
│   │   └── presenter/
│   │       ├── user_presenter.go
│   │       └── order_presenter.go
│   └── infrastructure/
│       ├── repository/
│       │   ├── user_repository.go
│       │   ├── order_repository.go
│       │   └── cache_repository.go
│       ├── service/
│       │   ├── email_service.go
│       │   └── payment_service.go
│       └── config/
│           └── database.go
├── pkg/
│   └── shared/
│       ├── errors.go
│       └── validators.go
├── tests/
│   ├── integration/
│   │   └── user_test.go
│   └── unit/
│       └── domain_test.go
├── go.mod
├── go.sum
└── README.md
```

### 45.10.2 Estructura por Dominio (Alternative)

```
myapp/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── user/
│   │   ├── domain/
│   │   │   └── user.go
│   │   ├── application/
│   │   │   └── create_user.go
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── user_handler.go
│   │   │   └── repository/
│   │   │       └── postgres.go
│   │   └── service.go
│   ├── order/
│   │   ├── domain/
│   │   ├── application/
│   │   ├── adapter/
│   │   └── service.go
│   └── shared/
│       ├── database.go
│       └── config.go
└── tests/
```

### 45.10.3 Package Organization

```go
// internal/user/service.go - Punto de entrada del contexto
package user

import (
    "myapp/internal/user/application"
    "myapp/internal/user/adapter"
    "myapp/internal/user/domain"
)

type Service struct {
    createUserUC *application.CreateUserUseCase
    controller   *adapter.UserController
}

func NewService(db Database) *Service {
    repo := adapter.NewPostgresUserRepository(db)
    createUC := application.NewCreateUserUseCase(repo)
    controller := adapter.NewUserController(createUC)
    
    return &Service{
        createUserUC: createUC,
        controller:   controller,
    }
}

// Exponer solo lo necesario
func (s *Service) Controller() *adapter.UserController {
    return s.controller
}
```

---

## 45.11 Buenas Prácticas y Escalabilidad

### 45.11.1 Evolución de la Arquitectura

```
Fase 1: MVP
├── domain/
├── usecases/
├── adapter/
└── main.go

        ↓ (crece la complejidad)

Fase 2: Múltiples adaptadores
├── domain/
├── usecases/
├── adapter/
│   ├── http/
│   ├── grpc/
│   └── worker/
└── infrastructure/

        ↓ (múltiples equipos)

Fase 3: Microservicios
service-user/
├── domain/
├── usecases/
└── adapter/

service-order/
├── domain/
├── usecases/
└── adapter/

shared-libs/
```

### 45.11.2 Migraciones sin Dolor

```go
// Antes: Acoplado a PostgreSQL
type UserRepository interface {
    Save(user *domain.User, db *sql.DB) error
}

// Después: Agnóstico
type UserRepository interface {
    Save(user *domain.User) error
}

// Cambiar implementación es trivial
type OldPostgresRepo struct{ db *sql.DB }
type NewMongoRepo struct{ client *mongo.Client }

// Ambas satisfacen la interface
// No hay cambios en el código del dominio
```

### 45.11.3 Testing Pyramid

```
        /\          E2E (Pocos, Lentos)
       /  \         - Casos críticos
      /────\        - Flujos completos
     /  Integration\
    /  (Algunos)    \  - Interacción entre capas
   /──────────────────\
  /  Unit (Muchos)     \ - Lógica individual
 /─────────────────────\  - Rápidos, aislados
/________________________\
```

```go
// Tests de Unidad (rápidos, muchos)
func TestUserValidation(t *testing.T) {
    user, err := domain.NewUser("test@test.com", "Test")
    assert.NoError(t, err)
    assert.True(t, user.IsValid())
}

// Tests de Integración (intermedios)
func TestCreateUserUseCase_WithDB(t *testing.T) {
    db := setupTestDB()
    repo := NewPostgresUserRepository(db)
    uc := NewCreateUserUseCase(repo)
    
    resp, err := uc.Execute(&CreateUserRequest{...})
    assert.NoError(t, err)
}

// Tests E2E (pocos, lentos)
func TestUserCreationFlow(t *testing.T) {
    apiServer := startServer()
    client := newHTTPClient()
    
    resp := client.POST("/api/users", struct{...}{})
    assert.Equal(t, 201, resp.StatusCode)
}
```

### 45.11.4 Refactoring Seguro

```go
// Con Clean Architecture, refactoring es seguro
// Tests cubren cambios internos

// Antes
type UserRepository interface {
    Save(user *domain.User) error
}

// Necesitamos agregar caching
// Opción 1: Decorator pattern
type CachedUserRepository struct {
    repo   UserRepository
    cache  CacheService
}

func (c *CachedUserRepository) Save(user *domain.User) error {
    err := c.repo.Save(user)
    if err == nil {
        c.cache.Set(user.ID, user, 1*time.Hour)
    }
    return err
}

// Tests existentes pasan, no cambió la interface
// Solo implementación interna
```

### 45.11.5 Scaling de Equipos

```
Pequeño equipo (1-3):
- Estructura simple
- 1-2 servicios
- Repositorio monolítico

Equipo medio (4-10):
- Múltiples contextos de dominio
- Clara separación de responsabilidades
- Integración clara entre equipos

Equipo grande (10+):
- Microservicios por dominio
- Equipos autónomos
- Comunicación asincrónica
```

### 45.11.6 Anti-patrones a Evitar

```go
// ❌ ANTI-PATRÓN 1: Demasiadas capas
// domain → entities → services → application → usecase → 
// handler → presenter → formatter → serializer
// (Simplifica, las capas son conceptuales no físicas)

// ❌ ANTI-PATRÓN 2: Lógica en adaptadores
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Aquí no: lógica de negocio, validaciones complejas, etc.
    // Usar use cases
}

// ❌ ANTI-PATRÓN 3: Dependencies circulares
// domain conoce application
// application conoce adapter
// adapter conoce domain (circular!)

// ✅ BIEN: Inyectar abstracciones
type UserRepository interface { /* ... */ }
// domain no conoce implementation
// adapter implementa interface definida en domain

// ❌ ANTI-PATRÓN 4: Interfaces genéricas
type Repository interface {
    Get(id interface{}) (interface{}, error)
    Save(obj interface{}) error
}

// ✅ BIEN: Interfaces específicas del dominio
type UserRepository interface {
    GetByID(id string) (*User, error)
    SaveUser(user *User) error
}

// ❌ ANTI-PATRÓN 5: God Objects
type User struct {
    // 50 campos...
    // Métodos para manipular todo
}

// ✅ BIEN: Objetos pequeños, cohesivos
type UserProfile struct {
    Bio    string
    Avatar string
}

type UserSecurity struct {
    Password string
    MFACode  string
}
```

---

## EJERCICIOS PRÁCTICOS

### Ejercicio 1: Arquitectura Simple (3 Capas Básicas)

Implementa una aplicación simple de TODO con 3 capas:

**Requisitos:**
- Capa de Dominio: Entidad `Task`
- Capa de Aplicación: Use case `CreateTask`
- Capa de Presentación: HTTP handler

```go
package main

import (
    "fmt"
    "time"
)

// CAPA DE DOMINIO
type Task struct {
    ID        string
    Title     string
    Completed bool
    CreatedAt time.Time
}

func NewTask(title string) (*Task, error) {
    if len(title) == 0 {
        return nil, fmt.Errorf("título requerido")
    }
    return &Task{
        ID:        fmt.Sprintf("task_%d", time.Now().UnixNano()),
        Title:     title,
        Completed: false,
        CreatedAt: time.Now(),
    }, nil
}

// CAPA DE APLICACIÓN
type TaskRepository interface {
    Save(task *Task) error
    GetAll() []*Task
}

type CreateTaskUseCase struct {
    repo TaskRepository
}

func NewCreateTaskUseCase(repo TaskRepository) *CreateTaskUseCase {
    return &CreateTaskUseCase{repo: repo}
}

func (uc *CreateTaskUseCase) Execute(title string) (*Task, error) {
    task, err := NewTask(title)
    if err != nil {
        return nil, err
    }
    
    if err := uc.repo.Save(task); err != nil {
        return nil, err
    }
    
    return task, nil
}

// CAPA DE PRESENTACIÓN
type InMemoryTaskRepository struct {
    tasks map[string]*Task
}

func (r *InMemoryTaskRepository) Save(task *Task) error {
    r.tasks[task.ID] = task
    return nil
}

func (r *InMemoryTaskRepository) GetAll() []*Task {
    tasks := make([]*Task, 0, len(r.tasks))
    for _, t := range r.tasks {
        tasks = append(tasks, t)
    }
    return tasks
}

// Prueba
func main() {
    repo := &InMemoryTaskRepository{tasks: make(map[string]*Task)}
    uc := NewCreateTaskUseCase(repo)
    
    task, _ := uc.Execute("Aprender Clean Architecture")
    fmt.Printf("Tarea creada: %s (ID: %s)\n", task.Title, task.ID)
}
```

**Tu tarea:** Agrega:
- Método para marcar tarea como completada
- HTTP handler usando net/http
- Tests unitarios para la entidad

---

### Ejercicio 2: Separación de Lógica de Negocio

Refactoriza un servicio monolítico en Use Cases separados:

```go
package main

// ❌ ANTES: Todo mezclado
type OrderService struct {
    db Database
}

func (s *OrderService) ProcessOrder(req OrderRequest) {
    // Validar
    // Calcular impuesto
    // Procesar pago
    // Guardar en BD
    // Enviar email
    // Todo en un método gigante
}

// ✅ DESPUÉS: Use cases separados
// Tu tarea: Implementar
type CreateOrderUseCase struct {
    repo OrderRepository
}

type CalculateTaxUseCase struct {
    taxService TaxService
}

type ProcessPaymentUseCase struct {
    paymentGateway PaymentGateway
}

type SendOrderConfirmationUseCase struct {
    emailService EmailService
}

// Combinar en orquestador
type OrderOrchestrator struct {
    create CreateOrderUseCase
    tax    CalculateTaxUseCase
    pay    ProcessPaymentUseCase
    notify SendOrderConfirmationUseCase
}

func (o *OrderOrchestrator) Execute(req OrderRequest) error {
    // Ejecutar secuencialmente
    // Cada use case es testeable independientemente
}
```

**Tu tarea:**
- Implementa cada use case
- Crea tests para cada uno
- Implementa el orquestador

---

### Ejercicio 3: Abstracción con Interfaces

Haz que un servicio sea agnóstico de la base de datos:

```go
package main

// ✅ Define interfaz en dominio
type UserRepository interface {
    Save(user *User) error
    FindByID(id string) (*User, error)
}

// ✅ Implementación para PostgreSQL
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(user *User) error {
    // Query SQL
}

func (r *PostgresUserRepository) FindByID(id string) (*User, error) {
    // Query SQL
}

// Tu tarea: Implementa
type MongoUserRepository struct {
    client *mongo.Client
}

func (r *MongoUserRepository) Save(user *User) error {
    // Operación MongoDB
}

func (r *MongoUserRepository) FindByID(id string) (*User, error) {
    // Operación MongoDB
}

// Mock para tests
type MockUserRepository struct {
    SaveFunc    func(*User) error
    FindByIDFunc func(string) (*User, error)
}

// El mismo use case funciona con todos
type CreateUserUseCase struct {
    repo UserRepository
}

func (uc *CreateUserUseCase) Execute(user *User) error {
    return uc.repo.Save(user)
}

// Prueba
func main() {
    var repo UserRepository
    
    // Usa PostgreSQL
    repo = NewPostgresUserRepository(db)
    
    // O MongoDB
    repo = NewMongoUserRepository(client)
    
    // O Mock
    repo = &MockUserRepository{
        SaveFunc: func(u *User) error { return nil },
    }
    
    // El use case no sabe la diferencia
    uc := NewCreateUserUseCase(repo)
    uc.Execute(user)
}
```

---

### Ejercicio 4: Implementar Agregados (DDD)

Diseña un agregado `ShoppingCart`:

```go
package main

import "time"

// ✅ Raíz del agregado
type ShoppingCart struct {
    ID        string
    UserID    string
    Items     []*CartItem
    Total     float64
    Status    CartStatus
    CreatedAt time.Time
}

type CartItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

type CartStatus string

const (
    StatusActive    CartStatus = "ACTIVE"
    StatusAbandoned CartStatus = "ABANDONED"
    StatusCompleted CartStatus = "COMPLETED"
)

// Tu tarea: Implementa
func (c *ShoppingCart) AddItem(productID string, quantity int, price float64) error {
    // Validar que el carrito está activo
    // Validar cantidad > 0
    // Evitar duplicados
    // Recalcular total
}

func (c *ShoppingCart) RemoveItem(productID string) error {
    // Validar que existe
    // Remover
    // Recalcular total
}

func (c *ShoppingCart) Checkout() error {
    // Validar no vacío
    // Cambiar estado a COMPLETED
    // Emitir evento CheckoutCompleted
}

func (c *ShoppingCart) IsAbandoned(maxInactiveTime time.Duration) bool {
    // Verificar si ha estado ACTIVE más de maxInactiveTime
}

func (c *ShoppingCart) recalculateTotal() {
    c.Total = 0
    for _, item := range c.Items {
        c.Total += float64(item.Quantity) * item.Price
    }
}

// Tests
func TestShoppingCart(t *testing.T) {
    cart := &ShoppingCart{
        ID:     "cart_123",
        UserID: "user_456",
        Items:  make([]*CartItem, 0),
        Status: StatusActive,
    }
    
    // Tu tarea: Escribir tests
    // - Agregar item
    // - Remover item
    // - Validar total
    // - Validar estado
}
```

---

### Ejercicio 5: Proyecto Completo - Sistema de Reservas

Implementa un sistema de reservas de restaurante con Clean Architecture completa:

**Requisitos:**
1. Entidades: Restaurant, Reservation, User
2. Use Cases: CreateReservation, CancelReservation, ListAvailableTimes
3. Adaptadores: HTTP Controllers, PostgreSQL Repositories
4. Servicios: Email notification, Time availability check

```go
package main

// ESTRUCTURA ESPERADA
// cmd/api/main.go
// internal/domain/
//   ├── restaurant.go
//   ├── reservation.go
//   ├── user.go
//   └── events.go
// internal/application/
//   ├── usecases/
//   │   ├── create_reservation.go
//   │   ├── cancel_reservation.go
//   │   └── list_times.go
//   └── ports/
//       └── repositories.go
// internal/adapter/
//   ├── http/
//   │   └── reservation_handler.go
//   └── repository/
//       └── postgres_reservation.go
// internal/infrastructure/
//   └── config.go

package domain

import "time"

type Reservation struct {
    ID           string
    RestaurantID string
    UserID       string
    Datetime     time.Time
    PartySize    int
    Status       ReservationStatus
    CreatedAt    time.Time
}

type ReservationStatus string

const (
    StatusConfirmed ReservationStatus = "CONFIRMED"
    StatusCancelled ReservationStatus = "CANCELLED"
)

// Tu tarea:
// 1. Definir el dominio completo
// 2. Implementar use cases
// 3. Crear adaptadores HTTP
// 4. Implementar repositorios
// 5. Escribir tests (al menos 50% cobertura)
// 6. Documentar decisiones arquitectónicas

// Bonus: Implementar eventos de dominio para notificaciones
```

---

## Resumen y Conclusiones

Clean Architecture en Go proporciona:

✅ **Independencia de Frameworks**: Cambiar HTTP a gRPC es trivial  
✅ **Testabilidad**: Cada capa es testeable independientemente  
✅ **Mantenibilidad**: Código organizado, fácil de entender  
✅ **Escalabilidad**: Crece con el proyecto sin fricción  
✅ **Flexibilidad**: Implementaciones intercambiables  

### Reglas Clave

1. **Respetar la Regla de Dependencia**: Las dependencias apuntan hacia adentro
2. **Interfaces pequeñas**: Específicas del dominio, no genéricas
3. **Separar responsabilidades**: Cada componente una tarea
4. **Tests primero**: Arquitectura probada desde el inicio
5. **Documentar decisiones**: Registrar por qué elegiste cada estructura

### Próximos Pasos

- Aplicar en tu próximo proyecto
- Experimentar con Microservicios
- Estudiar Event Sourcing
- Aprender sobre CQRS
- Combinar con patrones de concurrencia de Go

Clean Architecture no es dogma, sino guidelines. Adapta a tu contexto, pero mantén los principios fundamentales.


---

## Ejecuta este ejemplo

Código fuente completo y verificado (`go build ./... && go vet ./...`) en
[`examples/45-clean-architecture/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/45-clean-architecture):

```bash
cd examples/45-clean-architecture
go run .
```
