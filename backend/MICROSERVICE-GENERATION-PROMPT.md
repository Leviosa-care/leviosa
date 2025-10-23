# Microservice Generation Prompt

## Overview

Generate a Go microservice following the exact hexagonal architecture patterns used in the Leviosa care platform. This service should be production-ready with comprehensive testing, proper error handling, and integration with the shared core library.

## Architecture Requirements

### Core Architecture Pattern
**Hexagonal Architecture (Ports & Adapters)**
- **Domain Layer**: Pure business entities with NO external dependencies
- **Ports Layer**: Interfaces defining contracts for repositories and services
- **Application Layer**: Use cases and business workflows
- **Adapters Layer**: Infrastructure implementations

### Application Architecture: Microservice vs Monolith

**This architecture works for both microservices and monolithic applications:**

#### Microservice Structure (Single Domain)
```
service-name/
├── internal/
│   ├── domain/           # Single domain (e.g., authuser)
│   │   ├── user.go
│   │   ├── session.go
│   │   └── user_dto.go
│   ├── ports/            # Interfaces for this domain
│   └── adapters/         # Infrastructure for this domain
```

#### Monolith Structure (Multiple Domains)
```
monolith-app/
├── internal/
│   ├── domain/
│   │   ├── authuser/     # Auth user domain (same as microservice)
│   │   │   ├── user.go
│   │   │   ├── session.go
│   │   │   ├── user_dto.go
│   │   │   └── value_objects.go
│   │   ├── booking/      # Booking domain
│   │   │   ├── booking.go
│   │   │   ├── booking_dto.go
│   │   │   └── validation.go
│   │   ├── catalog/      # Catalog domain
│   │   │   ├── product.go
│   │   │   ├── category.go
│   │   │   └── product_dto.go
│   │   └── notification/ # Notification domain
│   │       ├── email.go
│   │       ├── sms.go
│   │       └── notification_dto.go
│   ├── ports/
│   │   ├── authuser/     # Interfaces for auth domain
│   │   ├── booking/      # Interfaces for booking domain
│   │   └── catalog/      # Interfaces for catalog domain
│   └── adapters/
│       ├── postgres/
│       │   ├── authuser/ # Auth repository implementations
│       │   ├── booking/  # Booking repository implementations
│       │   └── catalog/  # Catalog repository implementations
│       ├── http/
│       │   ├── authuser/ # Auth HTTP handlers
│       │   ├── booking/  # Booking HTTP handlers
│       │   └── catalog/  # Catalog HTTP handlers
│       └── rabbitmq/
│           ├── authuser/ # Auth message handlers
│           └── booking/  # Booking message handlers
```

**Key Principles for Both Architectures:**

1. **Domain Isolation**: Each domain maintains its own pure business logic
2. **Interface Separation**: Each domain has its own ports/interfaces
3. **Infrastructure Segregation**: Each domain has dedicated adapters
4. **Consistent Patterns**: Same validation, DTO, and repository patterns across all domains
5. **Shared Core**: Common utilities, error handling, and validation from shared core library

**Benefits:**
- **Microservices**: Easy to extract domains into separate services when needed
- **Monolith**: Organized structure that can be split into microservices later
- **Consistency**: Same patterns regardless of deployment model
- **Scalability**: Can scale individual domains independently in monolith or microservices

### Service Structure Template

```
service-name/
├── go.mod
├── makefile
├── main.go
├── internal/
│   ├── domain/           # Business entities (no dependencies)
│   │   ├── entity.go
│   │   ├── entity_test.go
│   │   ├── value_objects.go
│   │   └── enums.go
│   ├── ports/            # Interface definitions
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── aggregator.go
│   ├── application/      # Use cases and business logic
│   │   ├── service.go
│   │   ├── use_case_1.go
│   │   ├── use_case_2.go
│   │   └── aggregator/
│   │       ├── service.go
│   │       ├── workflow.go
│   │       └── workflow_test.go
│   └── adapters/         # Infrastructure implementations
│       ├── http/
│       │   ├── handlers.go
│       │   ├── routes.go
│       │   ├── middleware.go
│       │   └── helpers.go
│       ├── postgres/
│       │   ├── entity/
│       │   │   ├── repository.go
│       │   │   ├── queries.go
│       │   │   ├── main_test.go
│       │   │   └── *test.go
│       │   └── migrations/
│       ├── redis/
│       │   ├── cache/
│       │   │   ├── repository.go
│       │   │   ├── main_test.go
│       │   │   └── *test.go
│       │   └── session/
│       ├── rabbitmq/
│       │   ├── setup.go
│       │   ├── publisher.go
│       │   ├── consumer.go
│       │   └── handlers.go
│       └── external/
│           ├── stripe/
│           ├── s3/
│           └── vault/
├── test/
│   ├── integration/
│   │   ├── main_test.go
│   │   ├── entity_test.go
│   │   └── helpers_test.go
│   └── testdata/
│       ├── common.go
│       ├── database_operations.go
│       ├── http_helpers.go
│       ├── rabbitmq_helpers.go
│       └── entity_constructors.go
└── internal/migrations/
    └── *.sql
```

## Implementation Patterns

### 1. Error Handling Strategy

**3-Layer Error Handling with Sentinel Errors**

```go
// Repository Layer - Use ClassifyPgError/ClassifyRedisError
result, err := db.Query(ctx, query, args...)
if err != nil {
    return nil, errs.ClassifyPgError("operation description", err)
}

// Service Layer - Handle sentinel errors with explicit switches
data, err := s.repo.SomeMethod(ctx, params)
if err != nil {
    switch {
    case errors.Is(err, errs.ErrRepositoryNotFound):
        // Handle not found case
    case errors.Is(err, errs.ErrConnectionFailure):
        return fmt.Errorf("business operation: %w", err)
    default:
        return fmt.Errorf("business operation: %w", err)
    }
}

// HTTP Handler Layer - Map to HTTP status codes
if err != nil {
    var statusCode int
    switch {
    case errors.Is(err, errs.ErrInvalidValue):
        statusCode = http.StatusBadRequest
    case errors.Is(err, errs.ErrConnectionFailure):
        statusCode = http.StatusServiceUnavailable
    default:
        statusCode = http.StatusInternalServerError
    }
    httpx.RespondWithError(w, err, statusCode)
    return
}
```

### 2. Testing Strategy

**Black-Box Integration Testing with Testcontainers**

```go
// Repository Unit Tests
func TestCreateEntity(t *testing.T) {
    ctx := context.Background()

    t.Run("should successfully create entity", func(t *testing.T) {
        // Arrange
        td.ClearEntityTable(t, ctx, testPool)
        entity := td.NewTestEntity("test-data")

        // Process encryption
        err := crypto.ProcessStruct(ctx, entity)
        require.NoError(t, err)

        // Act
        err = repo.CreateEntity(ctx, entity)

        // Assert
        require.NoError(t, err)
        exists, err := repo.ExistsByID(ctx, entity.ID)
        require.NoError(t, err)
        assert.True(t, exists)
    })
}

// Integration Tests
func TestCreateEntityAPI(t *testing.T) {
    ctx := context.Background()
    client := &http.Client{Timeout: 10 * time.Second}

    t.Run("should successfully create entity via API", func(t *testing.T) {
        // Clean state
        td.ClearAllTestData(t, ctx, testPool, testClient)

        // Test HTTP request
        request := domain.CreateEntityRequest{Name: "Test Entity"}
        req := td.NewCreateEntityRequest(t, ctx, testServerURL, request)
        resp, err := client.Do(req)
        require.NoError(t, err)

        // Assert HTTP response
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        // Verify database persistence
        entity := td.GetEntityFromDB(t, ctx, testPool)
        assert.Equal(t, "Test Entity", entity.Name)

        // Verify RabbitMQ message (if applicable)
        td.VerifyEntityCreatedMessage(t, testCh, entity)
    })
}
```

### 3. Encryption Requirements

**Mandatory Encryption with github.com/hengadev/encx**

```go
// All sensitive data must be encrypted
type User struct {
    ID         uuid.UUID `db:"id"`
    Email      string    `encx:"email" db:"email"`
    EmailHash  []byte    `encx:"hash_basic:email" db:"email_hash"`
    Password   string    `encx:"password" db:"password"`
    // ... other fields
}

// Process encryption before persistence
err := crypto.ProcessStruct(ctx, &entity)
require.NoError(t, err)

// Process decryption after retrieval
err := crypto.DecryptStruct(ctx, &entity)
require.NoError(t, err)
```

### 4. Service Construction Pattern

```go
// Service Constructor
type EntityService struct {
    repo   ports.EntityRepository
    cache  ports.EntityCache
    crypto encx.CryptoService
}

func New(repo ports.EntityRepository, cache ports.EntityCache, crypto encx.CryptoService) ports.EntityService {
    return &EntityService{
        repo:   repo,
        cache:  cache,
        crypto: crypto,
    }
}

// Aggregator Constructor
type AggregatorService struct {
    entityRepo ports.EntityRepository
    entitySvc ports.EntityService
    notifier  ports.NotificationService
}

func New(entityRepo ports.EntityRepository, entitySvc ports.EntityService, notifier ports.NotificationService) ports.AggregatorService {
    return &AggregatorService{
        entityRepo: entityRepo,
        entitySvc:  entitySvc,
        notifier:   notifier,
    }
}
```

### 4.1. Repository Construction Pattern

**Each Repository Must Follow This Structure**

```go
// internal/adapters/postgres/entity/repository.go
package entityRepository

import (
    "context"

    "github.com/your-org/service-name/internal/ports"
    "github.com/jackc/pgx/v5/pgxpool"
    _ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
    pool   *pgxpool.Pool
    schema string
}

// New creates a new entity repository instance
func New(ctx context.Context, pool *pgxpool.Pool) ports.EntityRepository {
    return &Repository{
        pool:   pool,
        schema: "service_name", // Service-specific schema
    }
}

// Compile-time check to ensure Repository implements ports.EntityRepository
var _ ports.EntityRepository = (*Repository)(nil)
```

**Query Operations in Separate Files**

```go
// internal/adapters/postgres/entity/queries.go
package entityRepository

import (
    "context"
    "time"

    "github.com/your-org/service-name/internal/domain"
    "github.com/Leviosa-care/core/errs"
    "github.com/google/uuid"
)

func (r *Repository) CreateEntity(ctx context.Context, entity *domain.Entity) error {
    query := `
        INSERT INTO ` + r.schema + `.entities
        (id, name, description, email, email_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

    _, err := r.pool.Exec(ctx, query,
        entity.ID,
        entity.Name,
        entity.Description,
        entity.Email,
        entity.EmailHash,
        time.Now(),
        time.Now(),
    )

    if err != nil {
        return errs.ClassifyPgError("create entity", err)
    }

    return nil
}

func (r *Repository) GetEntityByID(ctx context.Context, id uuid.UUID) (*domain.Entity, error) {
    query := `
        SELECT id, name, description, email, email_hash, created_at, updated_at
        FROM ` + r.schema + `.entities
        WHERE id = $1`

    entity := &domain.Entity{}
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &entity.ID,
        &entity.Name,
        &entity.Description,
        &entity.Email,
        &entity.EmailHash,
        &entity.CreatedAt,
        &entity.UpdatedAt,
    )

    if err != nil {
        return nil, errs.ClassifyPgError("get entity by id", err)
    }

    return entity, nil
}
```

### 4.2. Port Interface Pattern

**Each Interface in Separate Files Under ports/**

```go
// internal/ports/entity_repository.go
package ports

import (
    "context"
    "time"

    "github.com/your-org/service-name/internal/domain"
    "github.com/google/uuid"
)

type EntityRepository interface {
    // CRUD Operations
    CreateEntity(ctx context.Context, entity *domain.Entity) error
    GetEntityByID(ctx context.Context, id uuid.UUID) (*domain.Entity, error)
    UpdateEntity(ctx context.Context, entity *domain.Entity) error
    DeleteEntity(ctx context.Context, id uuid.UUID) error

    // Query Operations
    GetAllEntities(ctx context.Context) ([]*domain.Entity, error)
    GetEntitiesByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Entity, error)
    SearchEntities(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Entity, error)

    // Existence Checks
    ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
    ExistsByName(ctx context.Context, name string) (bool, error)

    // Batch Operations
    BulkCreateEntities(ctx context.Context, entities []*domain.Entity) error
    BulkDeleteEntities(ctx context.Context, ids []uuid.UUID) error
}

// internal/ports/entity_service.go
package ports

import (
    "context"
    "time"

    "github.com/your-org/service-name/internal/domain"
    "github.com/google/uuid"
)

type EntityService interface {
    // Business Logic Operations
    CreateEntity(ctx context.Context, request *domain.CreateEntityRequest) (*domain.CreateEntityResponse, error)
    UpdateEntity(ctx context.Context, request *domain.UpdateEntityRequest) (*domain.Entity, error)
    DeleteEntity(ctx context.Context, request *domain.DeleteEntityRequest) error

    // Retrieval Operations
    GetEntity(ctx context.Context, request *domain.GetEntityRequest) (*domain.Entity, error)
    ListEntities(ctx context.Context, request *domain.ListEntitiesRequest) (*domain.ListEntitiesResponse, error)
    SearchEntities(ctx context.Context, request *domain.SearchEntitiesRequest) (*domain.SearchEntitiesResponse, error)

    // Validation Operations
    ValidateEntityAccess(ctx context.Context, userID, entityID uuid.UUID) (bool, error)
    ValidateEntityOwnership(ctx context.Context, userID, entityID uuid.UUID) (bool, error)
}

// internal/ports/entity_cache.go
package ports

import (
    "context"
    "time"

    "github.com/your-org/service-name/internal/domain"
    "github.com/google/uuid"
)

type EntityCache interface {
    // Cache Operations
    SetEntity(ctx context.Context, key string, entity *domain.Entity, ttl time.Duration) error
    GetEntity(ctx context.Context, key string) (*domain.Entity, error)
    DeleteEntity(ctx context.Context, key string) error

    // Search Cache
    SetSearchResult(ctx context.Context, key string, entities []*domain.Entity, ttl time.Duration) error
    GetSearchResult(ctx context.Context, key string) ([]*domain.Entity, error)

    // Cache Invalidation
    InvalidateUserEntities(ctx context.Context, userID uuid.UUID) error
    InvalidateEntity(ctx context.Context, entityID uuid.UUID) error
}
```

### 4.3. External Service Adapter Pattern

**Service Implementation in adapters/external/{service}/**

```go
// internal/adapters/external/stripe/service.go
package stripeService

import (
    "github.com/your-org/service-name/internal/ports"
    "github.com/stripe/stripe-go/v82"
)

// service handles Stripe operations
type service struct {
    client *stripe.Client
}

// Compile-time check to ensure service implements ports.StripeService
var _ ports.StripeService = (*service)(nil)

// NewService creates a new Stripe service instance
func NewService(apiKey, baseURL string) ports.StripeService {
    var client *stripe.Client

    if baseURL != "" {
        // Test environment with custom base URL
        backends := stripe.NewBackendsWithConfig(&stripe.BackendConfig{
            URL: &baseURL,
        })
        client = stripe.NewClient(apiKey, stripe.WithBackends(backends))
    } else {
        // Production environment
        client = stripe.NewClient(apiKey)
    }

    return &service{client: client}
}

// CreateCustomer implements ports.StripeService
func (s *service) CreateCustomer(ctx context.Context, request *domain.CreateCustomerRequest) (*domain.Customer, error) {
    // Implementation details...
}

// Additional service methods...
```

### 4.4. HTTP Helper Functions Pattern

**helpers.go File in Each HTTP Handler Package**

```go
// internal/adapters/http/entity/helpers.go
package entityHandler

import (
    "strings"
    "net/mail"
)

// maskEmail masks email for GDPR-compliant logging
func maskEmail(email string) string {
    if email == "" {
        return "[empty]"
    }

    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "[invalid]"
    }

    local, domain := parts[0], parts[1]
    if len(local) == 0 {
        return "[invalid]"
    }

    // Show first character + asterisks + last character if long enough
    if len(local) <= 2 {
        return string(local[0]) + "*@" + domain
    }

    return string(local[0]) + strings.Repeat("*", len(local)-2) + string(local[len(local)-1]) + "@" + domain
}

// validateEmail validates email format
func validateEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

// extractUserIDFromContext extracts user ID from request context
func extractUserIDFromContext(r *http.Request) (uuid.UUID, error) {
    userID, ok := r.Context().Value("user_id").(uuid.UUID)
    if !ok {
        return uuid.Nil, errs.NewUnauthenticatedErr("user not found in context")
    }
    return userID, nil
}

// parseUUID parses UUID from string parameter
func parseUUID(idStr string) (uuid.UUID, error) {
    id, err := uuid.Parse(idStr)
    if err != nil {
        return uuid.Nil, errs.NewInvalidValueErr("invalid UUID format")
    }
    return id, nil
}
```

### 4.5. Domain Validation Function Pattern

**Every Domain Type Must Have Associated Validation Functions**

```go
// internal/domain/value_objects.go
package domain

import (
    "context"
    "fmt"
    "strings"
    "regexp"

    "github.com/hengadev/errsx"
)

// Example: Name Validation
const (
    NameMinLength = 2
    NameMaxLength = 50
)

// Validation error keys and messages
const (
    NameRequiredKey     = "name_required"
    NameTooShortKey     = "name_too_short"
    NameTooLongKey      = "name_too_long"
    NameInvalidCharsKey = "name_invalid_chars"
)

const (
    NameRequiredMsg     = "is required"
    NameTooShortMsg     = "must be at least 2 characters long"
    NameTooLongMsg      = "must be no more than 50 characters long"
    NameInvalidCharsMsg = "contains invalid characters"
)

func validateName(name, fieldName string) error {
    var errs errsx.Map

    trimmed := strings.TrimSpace(name)

    if len(trimmed) == 0 {
        errs.Set(NameRequiredKey, fmt.Sprintf("%s %s", fieldName, NameRequiredMsg))
    }

    if len(trimmed) < NameMinLength {
        errs.Set(NameTooShortKey, fmt.Sprintf("%s %s", fieldName, NameTooShortMsg))
    }

    if len(trimmed) > NameMaxLength {
        errs.Set(NameTooLongKey, fmt.Sprintf("%s %s", fieldName, NameTooLongMsg))
    }

    // Reject dangerous patterns for GDPR compliance
    if strings.ContainsAny(trimmed, "<>;\"'&") {
        errs.Set(NameInvalidCharsKey, fmt.Sprintf("%s %s", fieldName, NameInvalidCharsMsg))
    }

    return errs.AsError()
}

// Example: Custom Code Validation (OTP, etc.)
func ValidateCode(ctx context.Context, code string, expectedLength int) error {
    var errs errsx.Map

    // Check empty
    if code == "" {
        errs.Set("code_missing", "code is required")
    }

    // Check length
    if len(code) != expectedLength {
        errs.Set("invalid_length", fmt.Sprintf("code must be %d digits", expectedLength))
    }

    // Check numeric only
    for _, r := range code {
        if r < '0' || r > '9' {
            errs.Set("invalid_characters", "code must only contain digits")
            break
        }
    }

    return errs.AsError()
}

// Example: Complex Field Validation
func validateComplexField(value, fieldName string) error {
    var errs errsx.Map

    // Custom regex pattern
    pattern := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
    if !pattern.MatchString(value) {
        errs.Set("invalid_format", fmt.Sprintf("%s contains invalid characters", fieldName))
    }

    // Additional business logic validation
    if strings.Contains(strings.ToLower(value), "admin") {
        errs.Set("reserved_word", fmt.Sprintf("%s cannot contain reserved words", fieldName))
    }

    return errs.AsError()
}
```

### 4.6. DTO (Data Transfer Object) Pattern

**Request/Response DTOs with Validation Methods**

```go
// internal/domain/entity_dto.go
package domain

import (
    "context"
    "time"

    "github.com/Leviosa-care/core/validation"
    "github.com/hengadev/errsx"
    "github.com/google/uuid"
)

// Request DTOs - Input Validation
type CreateEntityRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Email       string `json:"email"`
    CategoryID  uuid.UUID `json:"category_id,omitempty"`
}

func (r *CreateEntityRequest) Valid(ctx context.Context) error {
    var errs errsx.Map

    // Validate each field using domain validation functions
    if err := validateName(r.Name, "name"); err != nil {
        errs.Set("name", err)
    }

    // Length validation for description
    if len(r.Description) > 1000 {
        errs.Set("description", "description must be no more than 1000 characters")
    }

    // Email validation using core validation
    if err := validation.ValidateEmail(r.Email); err != nil {
        errs.Set("email", err)
    }

    // UUID validation (optional field)
    if r.CategoryID != uuid.Nil {
        // Additional business logic for category validation could go here
        // For now, we just validate it's a valid UUID format
    }

    return errs.AsError()
}

type UpdateEntityRequest struct {
    ID          uuid.UUID `json:"id"`
    Name        *string  `json:"name,omitempty"`
    Description *string  `json:"description,omitempty"`
    Email       *string  `json:"email,omitempty"`
}

func (r *UpdateEntityRequest) Valid(ctx context.Context) error {
    var errs errsx.Map

    // ID is required for updates
    if r.ID == uuid.Nil {
        errs.Set("id", "entity ID is required")
    }

    // Validate optional fields only if provided
    if r.Name != nil {
        if err := validateName(*r.Name, "name"); err != nil {
            errs.Set("name", err)
        }
    }

    if r.Description != nil {
        if len(*r.Description) > 1000 {
            errs.Set("description", "description must be no more than 1000 characters")
        }
    }

    if r.Email != nil {
        if err := validation.ValidateEmail(*r.Email); err != nil {
            errs.Set("email", err)
        }
    }

    return errs.AsError()
}

type DeleteEntityRequest struct {
    ID uuid.UUID `json:"id"`
}

func (r *DeleteEntityRequest) Valid(ctx context.Context) error {
    var errs errsx.Map

    if r.ID == uuid.Nil {
        errs.Set("id", "entity ID is required")
    }

    return errs.AsError()
}

// Response DTOs - Output Formatting
type EntityResponse struct {
    ID          uuid.UUID  `json:"id"`
    Name        string     `json:"name"`
    Description string     `json:"description"`
    Email       string     `json:"email"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateEntityResponse struct {
    Entity EntityResponse `json:"entity"`
    Message string        `json:"message"`
}

type ListEntitiesRequest struct {
    Page     int               `json:"page,omitempty"`
    Limit    int               `json:"limit,omitempty"`
    Search   string            `json:"search,omitempty"`
    Filters  map[string]string `json:"filters,omitempty"`
    SortBy   string            `json:"sort_by,omitempty"`
    SortDesc bool              `json:"sort_desc,omitempty"`
}

func (r *ListEntitiesRequest) Valid(ctx context.Context) error {
    var errs errsx.Map

    // Validate pagination
    if r.Page < 1 {
        errs.Set("page", "page must be greater than 0")
    }

    if r.Limit < 1 || r.Limit > 100 {
        errs.Set("limit", "limit must be between 1 and 100")
    }

    // Validate search length
    if len(r.Search) > 100 {
        errs.Set("search", "search term must be no more than 100 characters")
    }

    // Validate sort field (whitelist allowed fields)
    allowedSortFields := []string{"name", "created_at", "updated_at"}
    if r.SortBy != "" {
        valid := false
        for _, field := range allowedSortFields {
            if r.SortBy == field {
                valid = true
                break
            }
        }
        if !valid {
            errs.Set("sort_by", "invalid sort field")
        }
    }

    return errs.AsError()
}

type ListEntitiesResponse struct {
    Entities []EntityResponse `json:"entities"`
    Total    int              `json:"total"`
    Page     int              `json:"page"`
    Limit    int              `json:"limit"`
    HasNext  bool             `json:"has_next"`
}

// Validation with Context Parameters
type ValidateCodeRequest struct {
    Email string `json:"email"`
    Code  string `json:"code"`
}

func (r ValidateCodeRequest) Valid(ctx context.Context, expectedLength int) error {
    var errs errsx.Map

    if err := validation.ValidateEmail(r.Email); err != nil {
        errs.Set("email", err)
    }

    if err := ValidateCode(ctx, r.Code, expectedLength); err != nil {
        errs.Set("code", err)
    }

    return errs.AsError()
}
```

**DTO Naming Conventions**

1. **Request DTOs**: `{Action}{Entity}Request` (e.g., `CreateEntityRequest`, `UpdateUserRequest`)
2. **Response DTOs**: `{Action}{Entity}Response` (e.g., `CreateEntityResponse`, `GetUserResponse`)
3. **List Responses**: `List{Entity}Response` with pagination info
4. **Validation Method**: `Valid(ctx context.Context, ...params) error`
5. **JSON Tags**: Always include `json:"field_name"` tags
6. **Optional Fields**: Use pointers for optional primitive types (`*string`, `*int`)
7. **UUID Fields**: Use `uuid.UUID` type with validation for required fields
8. **Time Fields**: Use `time.Time` for timestamps

### 5. HTTP Handler Pattern

```go
func (h *handler) CreateEntity(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    logger, err := ctxutil.GetLoggerFromContext(ctx)
    if err != nil {
        httpx.RespondWithError(w, err, http.StatusInternalServerError)
        return
    }

    var payload domain.CreateEntityRequest
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    if err := decoder.Decode(&payload); err != nil {
        logger.WarnContext(ctx, "Handler: Invalid JSON request body",
            "error", err,
            "operation", "create_entity",
            "method", r.Method,
            "path", r.URL.Path)
        httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
        return
    }

    if err := h.svc.CreateEntity(ctx, &payload); err != nil {
        // Detailed error logging with context
        var logLevel string
        var errorContext string
        switch {
        case errors.Is(err, errs.ErrInvalidValue):
            logLevel = "warn"
            errorContext = "invalid input validation"
        // ... other error cases
        }

        logger.ErrorContext(ctx, "Handler: Entity creation failed",
            "error_context", errorContext,
            "error", err)

        // Map to HTTP status code
        var statusCode int
        switch {
        case errors.Is(err, errs.ErrInvalidValue):
            statusCode = http.StatusBadRequest
        // ... other status mappings
        default:
            statusCode = http.StatusInternalServerError
        }

        httpx.RespondWithError(w, err, statusCode)
        return
    }

    logger.InfoContext(ctx, "Handler: Entity creation completed successfully")
    httpx.RespondWithJSON(w, struct {
        Message string `json:"message"`
        ID      string `json:"id"`
    }{
        Message: "Entity created successfully",
        ID:      payload.ID,
    }, http.StatusCreated)
}
```

### 5.1. Routing Pattern

**Each HTTP Handler Package Must Have a routes.go File with RegisterRoutes Function**

```go
// internal/adapters/http/entity/routes.go
package entityHandler

import (
    "net/http"

    "github.com/Leviosa-care/core/contracts/identity"
    mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
    // Middleware shortcuts for role-based access control
    RequireVisitor := h.authmw.RequireMinimumRole(identity.Visitor)
    RequireStandard := h.authmw.RequireMinimumRole(identity.Standard)
    RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)
    RequireAdmin := h.authmw.RequireAdmin

    // Public endpoints (no authentication required)
    // Creates a new entity with public access
    router.HandleFunc("POST /entities", mw.EnableCORS(h.CreateEntity))

    // Gets a list of all entities (public read access)
    router.HandleFunc("GET /entities", mw.EnableCORS(h.GetAllEntities))

    // Standard user endpoints (requires Standard role or higher)
    // Gets the current user's entities
    router.HandleFunc("GET /users/me/entities", RequireStandard(mw.EnableCORS(h.GetUserEntities)))

    // Updates an entity owned by the current user
    router.HandleFunc("PATCH /entities/{id}", RequireStandard(mw.EnableCORS(h.UpdateEntity)))

    // Deletes an entity owned by the current user
    router.HandleFunc("DELETE /entities/{id}", RequireStandard(mw.EnableCORS(h.DeleteEntity)))

    // Admin-only endpoints
    // Gets all entities in the system (admin access)
    router.HandleFunc("GET /admin/entities", RequireAdministrator(mw.EnableCORS(h.AdminGetAllEntities)))

    // Updates any entity in the system (admin access)
    router.HandleFunc("PATCH /admin/entities/{id}", RequireAdministrator(mw.EnableCORS(h.AdminUpdateEntity)))

    // Deletes any entity in the system (admin access)
    router.HandleFunc("DELETE /admin/entities/{id}", RequireAdministrator(mw.EnableCORS(h.AdminDeleteEntity)))

    // Updates entity status/role (admin only)
    router.HandleFunc("PATCH /admin/entities/{id}/status", RequireAdministrator(mw.EnableCORS(h.UpdateEntityStatus)))

    // Specialized endpoints
    // Search entities with filters
    router.HandleFunc("GET /entities/search", mw.EnableCORS(h.SearchEntities))

    // Get entity by UUID with detailed information
    router.HandleFunc("GET /entities/{id}", mw.EnableCORS(h.GetEntityByID))

    // Bulk operations on entities
    router.HandleFunc("POST /entities/bulk", RequireStandard(mw.EnableCORS(h.BulkCreateEntities)))

    router.HandleFunc("DELETE /entities/bulk", RequireStandard(mw.EnableCORS(h.BulkDeleteEntities)))
}
```

**Route Registration in Main Service**

```go
// In main.go or service initialization
func (h *handler) Routes() http.Handler {
    router := http.NewServeMux()

    // Register all handler routes
    h.entityHandler.RegisterRoutes(router)
    h.userHandler.RegisterRoutes(router)
    h.adminHandler.RegisterRoutes(router)

    // Apply global middleware
    return mw.Chain(
        mw.RequestID,
        mw.Logging,
        mw.Recovery,
        mw.Timeout(30*time.Second),
    )(router)
}
```

**Route Organization Principles**

1. **Group by Resource**: Organize routes by entity/resource type
2. **Consistent Patterns**: Use RESTful conventions (GET, POST, PATCH, DELETE)
3. **Role-Based Access**: Apply appropriate authentication middleware
4. **Clear Documentation**: Add comments explaining each endpoint's purpose
5. **Path Variables**: Use `{id}` for UUID-based resource identification
6. **Admin Separation**: Separate admin routes with `/admin/` prefix
7. **CORS Support**: Apply `mw.EnableCORS` to all routes
8. **Middleware Chaining**: Combine authentication and other middleware as needed

### 6. Repository Implementation Pattern

```go
// PostgreSQL Repository
func (r *repository) CreateEntity(ctx context.Context, entity *domain.Entity) error {
    query := `
        INSERT INTO entities (id, name, email, email_hash, created_at)
        VALUES ($1, $2, $3, $4, $5)`

    _, err := r.pool.Exec(ctx, query,
        entity.ID,
        entity.Name,
        entity.Email,
        entity.EmailHash,
        time.Now())

    if err != nil {
        return errs.ClassifyPgError("create entity", err)
    }

    return nil
}

// Redis Repository
func (r *repository) SetEntityCache(ctx context.Context, key string, entity *domain.Entity, ttl time.Duration) error {
    data, err := json.Marshal(entity)
    if err != nil {
        return errs.NewInvalidValueErr(fmt.Sprintf("marshal entity: %v", err))
    }

    err = r.client.Set(ctx, key, data, ttl).Err()
    if err != nil {
        return errs.ClassifyRedisError("set entity cache", err)
    }

    return nil
}
```

## Configuration Requirements

### Go Module Configuration

```go
// go.mod
module github.com/your-org/service-name

go 1.24.2

require (
    github.com/Leviosa-care/core v0.0.0
    github.com/hengadev/encx v0.5.3
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.7.5
    github.com/redis/go-redis/v9 v9.12.1
    github.com/rabbitmq/amqp091-go v1.10.0
    github.com/stretchr/testify v1.11.0
)
```

### Makefile Template

```makefile
# Test execution commands
.PHONY: test test-unit test-integration test-all test-verbose test-coverage test-smoke test-race benchmark
.PHONY: test-postgres test-redis test-unit-entity test-unit-cache
.PHONY: test-integration-entity test-integration-workflow

# Unit tests by adapter type
test-unit:
	go test -v -count=1 ./internal/adapters/...

test-postgres:
	go test -v -count=1 ./internal/adapters/postgres/...

test-redis:
	go test -v -count=1 ./internal/adapters/redis/...

# Unit tests by specific adapter
test-unit-entity:
	go test -v -count=1 ./internal/adapters/postgres/entity/...

test-unit-cache:
	go test -v -count=1 ./internal/adapters/redis/cache/...

# Integration tests by module
test-integration:
	go test -v -count=1 ./test/integration/...

test-integration-entity:
	go test -v -count=1 ./test/integration/entity/...

# All tests
test-all:
	go test -v -count=1 ./...

# Run specific test by name
test-unit-entity-%:
	go test -v -count=1 ./internal/adapters/postgres/entity/... -run $(TEST)

test-integration-entity-%:
	go test -v -count=1 ./test/integration/entity/... -run $(TEST)

# Advanced testing
test-verbose:
	go test -v -count=1 -json ./... | tee test-output.json

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-smoke:
	go test -v -count=1 -short ./...

test-race:
	go test -v -race -count=1 ./...

benchmark:
	go test -v -bench=. -benchmem ./...

# Default test command
test: test-unit
```

## Database Migration Pattern

```sql
-- Migration naming: {timestamp}_{service}_{action}_{entity}.sql
-- Example: 20250112153022_entity_create_table_entities.sql

CREATE TABLE entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email TEXT NOT NULL,
    email_hash BYTEA NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_entities_email_hash ON entities(email_hash);
CREATE INDEX idx_entities_created_at ON entities(created_at);
```

## Test Data Helpers

### Testdata Package Structure

```go
// test/testdata/common.go
func ClearEntityTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool)
func ClearAllTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, client *redis.Client)

// test/testdata/database_operations.go
func InsertEntity(t *testing.T, ctx context.Context, entity *domain.Entity, pool *pgxpool.Pool)
func GetEntityFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) *domain.Entity
func EntityExists(t *testing.T, ctx context.Context, id uuid.UUID, pool *pgxpool.Pool) bool

// test/testdata/entity_constructors.go
func NewTestEntity(name string) *domain.Entity
func NewTestEntityWithEncryption(name string, crypto encx.CryptoService) (*domain.Entity, error)

// test/testdata/http_helpers.go
func NewCreateEntityRequest(t *testing.T, ctx context.Context, baseURL string, payload domain.CreateEntityRequest) *http.Request
func ParseCreateEntityResponse(t *testing.T, resp *http.Response) (*domain.Entity, error)
func ParseErrorResponse(t *testing.T, resp *http.Response) (string, int)

// test/testdata/rabbitmq_helpers.go
func SetupEntityQueue(t *testing.T, ch *amqp.Channel)
func VerifyEntityCreatedMessage(t *testing.T, delivery <-chan amqp.Delivery, expected *domain.Entity)
```

## Main Service Entry Point

```go
// main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/Leviosa-care/core/ctxutil"
    "github.com/Leviosa-care/core/httpx"
    "github.com/Leviosa-care/core/logger"
    "github.com/hengadev/encx"
    "github.com/hengadev/encx/providers/hashicorpvault"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()

    // Initialize logger
    lg := logger.New()
    ctx = ctxutil.WithLogger(ctx, lg)

    // Initialize connections
    pool, err := initPostgres(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize PostgreSQL: %v", err)
    }
    defer pool.Close()

    client, err := initRedis(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize Redis: %v", err)
    }
    defer client.Close()

    // Initialize crypto
    crypto, err := initCrypto(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize crypto: %v", err)
    }

    // Initialize repositories
    entityRepo := postgres.NewEntityRepository(pool)
    entityCache := redis.NewEntityCache(client)

    // Initialize services
    entitySvc := entity.New(entityRepo, entityCache, crypto)
    aggregatorSvc := aggregator.New(entityRepo, entitySvc, notificationSvc)

    // Initialize HTTP handlers
    handler := http.New(entitySvc, aggregatorSvc)

    // Start HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: handler.Routes(),
    }

    go func() {
        lg.InfoContext(ctx, "Server starting", "addr", server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    lg.InfoContext(ctx, "Server shutting down")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Server shutdown failed: %v", err)
    }

    lg.InfoContext(ctx, "Server stopped")
}
```

## Quality Requirements

### Code Quality Standards
- **No External Dependencies in Domain Layer**: Domain entities must be pure
- **Comprehensive Error Handling**: All errors must be classified and handled properly
- **Full Test Coverage**: Unit tests for all repositories, integration tests for all workflows
- **Security-First**: All sensitive data must be encrypted using encx
- **Production Logging**: Structured logging with context at all levels
- **Graceful Degradation**: Proper handling of infrastructure failures

### Documentation Requirements
- **Godoc Comments**: All exported functions must have documentation
- **Error Documentation**: All error cases must be documented
- **API Documentation**: Clear request/response examples
- **Deployment Guide**: Environment setup and configuration

### Performance Requirements
- **Connection Pooling**: Proper database connection management
- **Caching Strategy**: Redis caching for frequently accessed data
- **Async Processing**: RabbitMQ for background operations
- **Timeout Handling**: Appropriate timeouts for all external calls
- **Resource Management**: Proper cleanup and defer patterns

This architecture ensures consistency, maintainability, and scalability across all microservices in the platform while maintaining high standards for security, testing, and observability.
