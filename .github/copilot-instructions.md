# GitHub Copilot Instructions for GRAB (Go REST API Boilerplate)

**Version**: v2.0.0  
**Last Updated**: 2025-12-10  
**Purpose**: Developer-focused guidelines for building APIs with GRAB

---

## 📋 What is GRAB?

GRAB (Go REST API Boilerplate) is a production-ready Go REST API starter with:
- **Clean Architecture** (Handler → Service → Repository)
- **JWT Authentication** with refresh token rotation
- **Role-Based Access Control (RBAC)**
- **Database Migrations** (golang-migrate)
- **Docker-First Development**
- **89.81% Test Coverage**
- **Comprehensive Documentation**: https://vahiiiid.github.io/go-rest-api-docs/

---

## 🎯 Core Development Principles

### 1. **Environment Detection - Don't Hardcode Versions**
Instead of stating "Go 1.24" or "PostgreSQL 15", show how to check:

```bash
# Check Go version
go version

# Check Docker version
docker --version

# Check PostgreSQL version (inside container)
make exec-db
psql --version
```

### 2. **Docker-First Development**
- Developers run `make` commands on host
- **Makefile automatically detects** if Docker container is running
- Commands execute in container if available, host otherwise
- **No need to manually enter container** - the Makefile handles execution context

```bash
# Start containers first
make up

# Run tests (automatically in container if running)
make test

# Run linting (automatically in container if running)
make lint

# Apply lint fixes (automatically in container if running)
make lint-fix

# Generate Swagger docs (automatically in container if running)
make swag
```

### 3. **Clean Architecture Pattern**
Every domain follows this structure:

```
internal/
└── <domain>/
    ├── model.go       # Domain models (GORM)
    ├── dto.go         # Data Transfer Objects (API contracts)
    ├── repository.go  # Database access layer
    ├── service.go     # Business logic layer
    ├── handler.go     # HTTP handlers (Gin)
    └── *_test.go      # Tests for each layer
```

**Key Rules**:
- Handler → Service → Repository (never skip layers)
- No business logic in handlers
- No HTTP concerns in services
- Repository only talks to database

---

## 🚀 Common Development Tasks

### Adding a New Domain/Entity

**Example**: Adding a "Todo" entity

1. **Create directory structure**:
```bash
mkdir -p internal/todo
```

2. **Create model** (`internal/todo/model.go`):
```go
package todo

import (
    "time"
    "gorm.io/gorm"
)

type Todo struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Title       string         `gorm:"not null" json:"title"`
    Description string         `json:"description"`
    Completed   bool           `gorm:"default:false" json:"completed"`
    UserID      uint           `gorm:"not null" json:"user_id"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

3. **Create DTO** (`internal/todo/dto.go`):
```go
package todo

type CreateTodoRequest struct {
    Title       string `json:"title" binding:"required,min=3,max=200"`
    Description string `json:"description" binding:"max=1000"`
}

type UpdateTodoRequest struct {
    Title       string `json:"title" binding:"omitempty,min=3,max=200"`
    Description string `json:"description" binding:"omitempty,max=1000"`
    Completed   *bool  `json:"completed"`
}

type TodoResponse struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    UserID      uint      `json:"user_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

4. **Create repository** (`internal/todo/repository.go`):
```go
package todo

import (
    "context"
    "gorm.io/gorm"
)

type Repository interface {
    Create(ctx context.Context, todo *Todo) error
    FindByID(ctx context.Context, id uint) (*Todo, error)
    FindByUserID(ctx context.Context, userID uint) ([]Todo, error)
    Update(ctx context.Context, todo *Todo) error
    Delete(ctx context.Context, id uint) error
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

// Implementation methods...
```

5. **Create service** (`internal/todo/service.go`):
```go
package todo

import (
    "context"
    "go-rest-api-boilerplate/internal/errors"
)

type Service interface {
    CreateTodo(ctx context.Context, userID uint, req *CreateTodoRequest) (*TodoResponse, error)
    GetTodo(ctx context.Context, userID, todoID uint) (*TodoResponse, error)
    GetUserTodos(ctx context.Context, userID uint) ([]TodoResponse, error)
    UpdateTodo(ctx context.Context, userID, todoID uint, req *UpdateTodoRequest) (*TodoResponse, error)
    DeleteTodo(ctx context.Context, userID, todoID uint) error
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

// Implementation methods...
```

6. **Create handler** (`internal/todo/handler.go`):
```go
package todo

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/vahiiiid/go-rest-api-boilerplate/internal/contextutil"
    apiErrors "github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
)

type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

// @Summary Create todo
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTodoRequest true "Todo creation request"
// @Success 201 {object} errors.Response{success=bool,data=TodoResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/todos [post]
func (h *Handler) CreateTodo(c *gin.Context) {
    userID := contextutil.GetUserID(c)
    
    var req CreateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        _ = c.Error(apiErrors.FromGinValidation(err))
        return
    }
    
    todo, err := h.service.CreateTodo(c.Request.Context(), userID, &req)
    if err != nil {
        _ = c.Error(apiErrors.InternalServerError(err))
        return
    }
    
    c.JSON(http.StatusCreated, apiErrors.Success(todo))
}

// Additional handler methods...
```

7. **Create migration**:
```bash
make migrate-create NAME=create_todos_table
```

Edit the generated migration file:
```sql
-- migrations/YYYYMMDDHHMMSS_create_todos_table.up.sql
BEGIN;

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_deleted_at ON todos(deleted_at);

COMMIT;
```

```sql
-- migrations/YYYYMMDDHHMMSS_create_todos_table.down.sql
BEGIN;

DROP TABLE IF EXISTS todos;

COMMIT;
```

8. **Register routes** in `internal/server/router.go`:
```go
// Initialize todo components
todoRepo := todo.NewRepository(db)
todoService := todo.NewService(todoRepo)
todoHandler := todo.NewHandler(todoService)

// Register todo routes (authenticated endpoints)
todosGroup := v1.Group("/todos")
todosGroup.Use(auth.AuthMiddleware(authService))
{
    todosGroup.POST("", todoHandler.CreateTodo)
    todosGroup.GET("", todoHandler.GetUserTodos)
    todosGroup.GET("/:id", todoHandler.GetTodo)
    todosGroup.PUT("/:id", todoHandler.UpdateTodo)
    todosGroup.DELETE("/:id", todoHandler.DeleteTodo)
}
```

9. **Write tests** (`internal/todo/*_test.go`)

10. **Run migration and tests**:
```bash
make migrate-up
make test
make lint
make swag
```

---

### Creating Database Migrations

**Pattern**: `YYYYMMDDHHMMSS_verb_noun_table`

**Examples**:
- `20251025225126_create_users_table`
- `20251028000000_create_refresh_tokens_table`
- `20251122153000_create_roles_table`
- `20251210120000_add_avatar_to_users_table`

**⚠️ CRITICAL RULE: Development vs Production**

**DEVELOPMENT (Pre-Production)**:
- ✅ Modify original migration files for changes
- ✅ DROP and RECREATE tables
- ❌ Don't create new migrations for small changes

**PRODUCTION (After Deploy)**:
- ✅ ALWAYS create new migrations
- ❌ NEVER modify existing migrations
- ❌ NEVER drop tables

**Steps**:
```bash
# 1. Generate migration files
make migrate-create NAME=create_todos_table

# 2. Edit the .up.sql file
vim migrations/YYYYMMDDHHMMSS_create_todos_table.up.sql

# 3. Edit the .down.sql file (for rollback)
vim migrations/YYYYMMDDHHMMSS_create_todos_table.down.sql

# 4. Apply migration
make migrate-up

# 5. Verify
make migrate-status

# 6. If needed, rollback
make migrate-down
```

**Migration Best Practices**:
- Always wrap in `BEGIN;` / `COMMIT;` transactions
- Include indexes for foreign keys and frequently queried columns
- Use `IF NOT EXISTS` for safety
- Write corresponding `.down.sql` for every `.up.sql`
- Test rollback before committing

---

### Working with Authentication

**Protected Routes** (require valid JWT):
```go
import "github.com/vahiiiid/go-rest-api-boilerplate/internal/middleware"

// Authentication is typically handled by router setup
// Use middleware for role-based access control
```

**Role-Based Access**:
```go
import "github.com/vahiiiid/go-rest-api-boilerplate/internal/middleware"

// Admin-only route
v1.Use(middleware.RequireAdmin()).
   POST("/admin/users", userHandler.CreateUser)

// Specific role required
v1.Use(middleware.RequireRole("admin")).
   POST("/admin/reports", userHandler.CreateReport)
```

**Getting Current User**:
```go
import "github.com/vahiiiid/go-rest-api-boilerplate/internal/contextutil"

func (h *Handler) MyHandler(c *gin.Context) {
    userID := contextutil.GetUserID(c)
    userEmail := contextutil.GetEmail(c)
    userName := contextutil.GetUserName(c)
    userRoles := contextutil.GetRoles(c)
    isAdmin := contextutil.IsAdmin(c)
    hasRole := contextutil.HasRole(c, "moderator")
    
    // Use user information...
}
```

---

### Error Handling

Use the centralized error handling:

```go
import (
    "errors"
    apiErrors "github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
)

// Validation errors
if err := c.ShouldBindJSON(&req); err != nil {
    _ = c.Error(apiErrors.FromGinValidation(err))
    return
}

// Service errors - check specific errors first
todo, err := h.service.CreateTodo(ctx, userID, req)
if err != nil {
    // Check for known specific errors first
    if errors.Is(err, ErrTodoNotFound) {
        _ = c.Error(apiErrors.NotFound("Todo not found"))
        return
    }
    if errors.Is(err, ErrUnauthorized) {
        _ = c.Error(apiErrors.Unauthorized("Unauthorized access"))
        return
    }
    // Wrap unknown errors
    _ = c.Error(apiErrors.InternalServerError(err))
    return
}

c.JSON(http.StatusOK, apiErrors.Success(todo))
```

**Available Error Types**:
- `apiErrors.NotFound(message)` - 404 Not Found
- `apiErrors.Unauthorized(message)` - 401 Unauthorized
- `apiErrors.Forbidden(message)` - 403 Forbidden
- `apiErrors.BadRequest(message)` - 400 Bad Request
- `apiErrors.Conflict(message)` - 409 Conflict
- `apiErrors.InternalServerError(err)` - 500 Internal Server Error
- `apiErrors.TooManyRequests(retryAfter)` - 429 Too Many Requests

---

### Testing

**Test Structure**:
```go
func TestService_CreateTodo(t *testing.T) {
    tests := []struct {
        name        string
        userID      uint
        request     *CreateTodoRequest
        setupMocks  func(*MockRepository)
        expectError bool
        errorType   error
    }{
        {
            name:   "success",
            userID: 1,
            request: &CreateTodoRequest{
                Title: "Test Todo",
            },
            setupMocks: func(m *MockRepository) {
                m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
            },
            expectError: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := NewMockRepository(ctrl)
            if tt.setupMocks != nil {
                tt.setupMocks(mockRepo)
            }
            
            service := NewService(mockRepo)
            result, err := service.CreateTodo(context.Background(), tt.userID, tt.request)
            
            if tt.expectError {
                assert.Error(t, err)
                if tt.errorType != nil {
                    assert.Equal(t, tt.errorType, err)
                }
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
            }
        })
    }
}
```

**Run tests**:
```bash
make test              # Run all tests
make test-coverage     # Generate coverage report
make test-verbose      # Run with verbose output
```

---

### Updating Swagger Documentation

After adding/modifying endpoints:

```bash
# Regenerate Swagger docs
make swag

# View docs
open http://localhost:8080/swagger/index.html
```

**Swagger Annotations**:
```go
// @Summary Short description
// @Description Detailed description
// @Tags tag-name
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body CreateUserRequest true "User creation request"
// @Success 200 {object} UserResponse
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Validation error"
// @Router /api/v1/users [post]
```

---

## 📚 Out-of-the-Box Features

GRAB comes with these production-ready features:

1. **Authentication**: JWT with refresh tokens → [Docs](https://vahiiiid.github.io/go-rest-api-docs/AUTHENTICATION/)
2. **RBAC**: Role-based access control → [Docs](https://vahiiiid.github.io/go-rest-api-docs/RBAC/)
3. **Migrations**: Versioned database migrations → [Docs](https://vahiiiid.github.io/go-rest-api-docs/MIGRATIONS_GUIDE/)
4. **Health Checks**: `/health`, `/live`, `/ready` → [Docs](https://vahiiiid.github.io/go-rest-api-docs/HEALTH_CHECKS/)
5. **Rate Limiting**: Token bucket algorithm → [Docs](https://vahiiiid.github.io/go-rest-api-docs/RATE_LIMITING/)
6. **Structured Logging**: JSON logs with context → [Docs](https://vahiiiid.github.io/go-rest-api-docs/LOGGING/)
7. **API Response Format**: Standardized responses → [Docs](https://vahiiiid.github.io/go-rest-api-docs/API_RESPONSE_FORMAT/)
8. **Error Handling**: Centralized error management → [Docs](https://vahiiiid.github.io/go-rest-api-docs/ERROR_HANDLING/)
9. **Graceful Shutdown**: Clean server termination → [Docs](https://vahiiiid.github.io/go-rest-api-docs/GRACEFUL_SHUTDOWN/)
10. **Swagger/OpenAPI**: Auto-generated API docs → [Docs](https://vahiiiid.github.io/go-rest-api-docs/SWAGGER/)
11. **Context Helpers**: Request context utilities → [Docs](https://vahiiiid.github.io/go-rest-api-docs/CONTEXT_HELPERS/)

---

## 🔧 Configuration

Configuration uses YAML files + environment variables:

```yaml
# configs/config.yaml (base)
app:
  name: "GRAB API"
  environment: "development"

# Override with environment-specific files:
# - configs/config.development.yaml
# - configs/config.staging.yaml
# - configs/config.production.yaml

# Override individual values with env vars:
# DATABASE_PASSWORD=secret
# JWT_SECRET=secret
```

**Environment Variable Mapping**:
- `DATABASE_PASSWORD` → `database.password`
- `JWT_SECRET` → `jwt.secret`
- `APP_ENVIRONMENT` → `app.environment`

Full config guide: https://vahiiiid.github.io/go-rest-api-docs/CONFIGURATION/

---

## 🐳 Docker Commands

```bash
# Start all services
make up

# Stop all services
make down

# View logs
make logs

# Rebuild containers
make build

# Execute command in app container
make exec

# Execute command in db container
make exec-db

# Clean restart
make down && make up
```

---

## � Project-Specific Domain Rules

### Imoveis Domain Relationships

**Current Architecture** (as of January 2026):
```
Organizacao (1) → (N) CorretorPrincipal
                         ↓ (1)
                       Anexo (foto)
                         ↓ (1)
                       Imovel (N)
                         ↓ (N)
                       Anexo (fotos)
```

**Important Notes**:
- Imovel relates to `CorretorPrincipal`, NOT directly to `Organizacao`
- CorretorPrincipal has optional `FotoID` (points to anexos)
- Anexo can relate to: Imovel, Empreendimento, Planta (but NOT CorretorPrincipal - that's reversed)
- Use `Omit("FotoID")` when creating CorretorPrincipal to avoid FK violations

**Common Pitfalls**:
1. ❌ Don't use `imovel.OrganizacaoID` - use `imovel.CorretorPrincipalID`
2. ❌ Don't add `CorretorPrincipalID` to Anexo model - the relationship is reversed
3. ✅ Always add `TableName()` methods for Portuguese plural forms
4. ✅ Process relationships BEFORE creating/updating main entities

---

## �🔄 External Property Import System

### Overview
The external import system synchronizes properties from `dev-api-backend.pi8.com.br` with complete relationship management and smart anexos synchronization.

### Key Implementation Details

**Critical Success Factors**:

1. **Process ALL relationships BEFORE create/update**
   - Move relationship upsert logic BEFORE the if/else block
   - This ensures both create and update paths use the same relationship IDs
   - Example:
     ```go
     // BEFORE if/else
     var empreendimentoID uint
     if ext.Empreendimento != nil {
         empID, err := is.upsertEmpreendimento(ctx, ext.Empreendimento)
         empreendimentoID = empID
     }
     
     // THEN if/else for create vs update
     if isUpdate {
         updateReq.EmpreendimentoID = &empreendimentoID
     }
     ```

2. **Use DELETE + INSERT for anexos synchronization**
   - Prevents stale data (removed images on external API)
   - Ensures correct image count and order
   - Implementation:
     ```go
     // Delete all existing anexos
     db.Where("imovel_id = ?", imovelID).Delete(&Anexo{})
     
     // Insert current anexos from external API
     for _, imageURL := range imageURLs {
         // create anexo
     }
     ```

3. **Use Omit() for zero-value foreign keys**
   - Prevents FK constraint violations
   - Example: `db.Omit("FotoID").Create(&corretor)`

4. **Add TableName() for Portuguese plurals**
   - GORM uses English pluralization by default
   - Required for: `organizacoes`, `corretores_principais`

### Commands That Work

**Import Execution**:
```bash
# Run import in container
docker exec -it triiio_app go run cmd/importimoveis/main.go

# Or using Makefile
make import-properties
```

**Verification Queries**:
```bash
# Verify all relationships
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "
SELECT i.id, i.codigo, i.titulo, e.preco as preco_venda, 
       cp.nome as corretor, o.nome as organizacao 
FROM imoveis i 
LEFT JOIN preco_vendas e ON i.preco_venda_id = e.id 
LEFT JOIN corretores_principais cp ON i.corretor_principal_id = cp.id 
LEFT JOIN organizacoes o ON cp.organizacao_id = o.id 
WHERE i.deleted_at IS NULL 
ORDER BY i.id;"

# Count anexos per property
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "
SELECT imovel_id, COUNT(*) as total_anexos 
FROM anexos 
WHERE deleted_at IS NULL 
GROUP BY imovel_id 
ORDER BY imovel_id;"

# Overall counts
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "
SELECT 'Imóveis:' as tipo, COUNT(*)::text as total FROM imoveis WHERE deleted_at IS NULL
UNION ALL SELECT 'Corretores:', COUNT(*)::text FROM corretores_principais WHERE deleted_at IS NULL
UNION ALL SELECT 'Organizações:', COUNT(*)::text FROM organizacoes WHERE deleted_at IS NULL
UNION ALL SELECT 'Anexos:', COUNT(*)::text FROM anexos WHERE deleted_at IS NULL;"
```

### Expected Success Output

```
import completed: X created, Y updated, 0 failed
Synced N anexos for property ID Z
```

### Common Errors Solved

1. **FK violation on `foto_id`**: Use `Omit("FotoID")` when creating corretores
2. **Table not found**: Add `TableName()` method for Portuguese plurals
3. **Stale relationships**: Process ALL upserts before create/update
4. **Anexos not syncing**: Use DELETE + INSERT, not incremental add

---

## 🧪 Pre-Commit Checklist

Before committing code:

```bash
# 1. Fix linting issues automatically
make lint-fix

# 2. Check for remaining issues
make lint

# 3. Run all tests
make test

# 4. Update Swagger docs (if API changed)
make swag

# 5. Verify everything works
make up
curl http://localhost:8080/health
```

---

## 📖 Additional Resources

- **Documentation Site**: https://vahiiiid.github.io/go-rest-api-docs/
- **Main Repository**: https://github.com/vahiiiid/go-rest-api-boilerplate
- **Issues**: https://github.com/vahiiiid/go-rest-api-boilerplate/issues
- **Discussions**: https://github.com/vahiiiid/go-rest-api-boilerplate/discussions

---

## 💡 Tips for AI Assistants

- **Always reference existing patterns**: Look at `internal/user/` for domain structure examples
- **Follow Clean Architecture**: Never skip layers (Handler → Service → Repository)
- **Use context helpers**: `contextutil.GetUserID(c)` for authenticated user info
- **Minimal comments**: Write self-documenting code, comment WHY not WHAT
- **Test thoroughly**: Maintain 85%+ test coverage
- **Check Makefile**: All development commands are in `make help`
- **Read the docs**: Comprehensive guides at https://vahiiiid.github.io/go-rest-api-docs/
- **Import Logic**: Always upsert relationships BEFORE creating/updating main entity
- **Anexos Sync**: Use DELETE + INSERT pattern to ensure data freshness

### 🎯 MANDATORY: Document Every Learning

**At the end of EVERY task completion, you MUST:**

1. Update this file or AGENTS.md with what you learned
2. Add successful commands to "Commands That Work" sections
3. Document any errors encountered and their solutions
4. Add verification queries that proved the solution worked

**This is NOT optional** - it's how we prevent repeating the same debugging sessions.

**Good documentation example**:
```markdown
#### Problem: Anexos not syncing correctly
- **Root Cause**: Using incremental add instead of full sync
- **Solution**: DELETE all existing + INSERT current
- **Code**: `db.Where("imovel_id = ?", id).Delete(&Anexo{})`
- **Verification**: `SELECT COUNT(*) FROM anexos WHERE imovel_id = X`
- **Result**: Synced N anexos successfully
```

**Every error you fix today becomes knowledge that prevents tomorrow's errors.**

---

**Last Updated**: 2026-01-30  
