# AGENTS.md - Mikmongo Project Guidelines

## Build, Lint, and Test Commands

### Build
```bash
go build ./...                    # Build all packages
go build -o bin/server ./cmd/server  # Build server binary
make build                        # Build Docker image
```

### Lint
```bash
golangci-lint run ./...           # Run linter
go fmt ./...                      # Format code
goimports -w .                    # Fix imports
```

### Test
```bash
# Run all unit tests
go test -v -race ./internal/... ./pkg/...

# Run single test file
go test -v -race ./internal/service/billing_service_test.go

# Run single test function
go test -v -race -run TestGenerateInvoice_NewSubscription ./internal/service/...

# Run single test with exact match
go test -v -race -run "^TestBillingIdempotency$" ./tests/integration/...

# Run integration tests (requires Docker)
go test -v -tags=integration ./tests/integration/...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Database Migrations
```bash
make migrate-up                   # Apply migrations
make migrate-down                 # Rollback last migration
make migrate-status               # Check migration status
make seed                         # Run seeds after migrations
make fresh                        # Reset + seed (DANGER: destroys data)
```

### Code Generation
```bash
make generate-mocks               # Generate mocks from interfaces
make model VAL=customer           # Generate model boilerplate
make service VAL=billing          # Generate service boilerplate
make handler VAL=billing          # Generate handler boilerplate
make repository VAL=customer      # Generate repository interface
make migration VAL=customers      # Generate migration file
```

## Code Style Guidelines

### Project Structure
```
cmd/                    # Application entrypoints (server, seed, migrate)
internal/               # Private application code
  handler/              # HTTP handlers (Gin framework)
  service/              # Business logic layer
  repository/           # Data access interfaces + implementations
    postgres/           # PostgreSQL implementations
  model/                # GORM models
  domain/               # Domain logic (DDD-style)
  scheduler/            # Cron jobs
  seeder/               # Database seeders
  migration/            # Goose migrations
  router/               # HTTP route definitions
pkg/                    # Public/reusable packages
tests/
  integration/          # Integration tests (require Docker)
  mocks/                # Generated mocks
```

### Imports Ordering
Three groups, separated by blank lines:
1. Standard library (context, fmt, time)
2. Third-party packages (github.com/..., gorm.io/...)
3. Local packages (mikmongo/...)

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/pkg/response"
)
```

### Naming Conventions

```go
// Structs/Types: PascalCase exported, camelCase internal
type BillingService struct { ... }           // Exported
type invoiceRepository struct { ... }        // Unexported implementation

// Interfaces: end with Repository/Service
type InvoiceRepository interface { ... }

// Methods: verb-first for actions
func (s *BillingService) GenerateInvoice(...)    // action
func (s *BillingService) List(ctx, limit, offset) // list with pagination

// Variables: short in small scopes, descriptive otherwise
inv, err := s.invoiceRepo.GetByID(ctx, id)    // short: clear scope
invoices, err := s.List(ctx, limit, offset)    // descriptive: returned
```

### Error Handling

```go
// Always wrap errors with context
if err != nil {
	return fmt.Errorf("operation failed: %w", err)
}

// Structured error messages
return nil, fmt.Errorf("subscription not found: %w", err)
return nil, fmt.Errorf("invalid plan ID: %w", err)
```

### Testing Conventions

```go
// Unit tests: same package with _test.go suffix
package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateInvoice_NewSubscription(t *testing.T) {
	svc, mocks := setupTest()
	result, err := svc.GenerateInvoice(ctx, subID, now)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// Test helper pattern
func newBillingServiceWithMocks() (*BillingService, *mocks.MockInvoiceRepository, ...) {
	invoiceRepo := &mocks.MockInvoiceRepository{}
	svc := NewBillingService(invoiceRepo, ...)
	return svc, invoiceRepo, ...
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
```

**Integration Tests:** Use `//go:build integration` tag
- MUST NOT use `t.Parallel()` (transactions share connection pool)
- Use `require` for setup, `assert` for assertions

### Handler Pattern

```go
func (h *BillingHandler) GetInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	invoice, err := h.service.GetInvoice(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, invoice)
}
```

### Service Pattern

```go
func (s *BillingService) GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}
	return invoice, nil
}
```

### GORM Model Pattern

```go
type Invoice struct {
	ID         string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL"`
	Status     string         `gorm:"type:varchar(20);not null;default:'draft'"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Customer   Customer       `gorm:"foreignKey:CustomerID"`
}
func (Invoice) TableName() string { return "invoices" }
```

### Context Usage
- Always pass `context.Context` as first parameter
- Use `c.Request.Context()` in handlers
- Use `r.db.WithContext(ctx)` in repositories

### Dependency Injection
```go
func NewBillingService(
	invoiceRepo repository.InvoiceRepository,
	subRepo repository.SubscriptionRepository,
) *BillingService {
	return &BillingService{
		invoiceRepo: invoiceRepo,
		subRepo:     subRepo,
	}
}

// Setter injection for circular dependencies
func (s *BillingService) SetNotificationService(n *NotificationService) {
	s.notificationSvc = n
}
```

## Key Dependencies
- **Web Framework:** gin-gonic/gin
- **ORM:** gorm.io/gorm with PostgreSQL driver
- **Testing:** stretchr/testify (assert, require, mock)
- **UUID:** google/uuid
- **Logging:** uber.org/zap
- **Migrations:** pressly/goose/v3
- **Queue:** rabbitmq/amqp091-go
- **Cache:** redis/go-redis/v9
- **Cron:** robfig/cron/v3