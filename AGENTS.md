# GRAB - AI-Friendly Development Guide

**Version**: v2.0.0  
**Last Updated**: 2025-12-10  
**Purpose**: Universal AI assistant guidelines for GRAB (Go REST API Boilerplate)

> This file follows the OpenAI AGENTS.md standard and is compatible with all major AI coding assistants including GitHub Copilot, Cursor, Windsurf, JetBrains AI, and others.

---

## 📋 Project Overview

**GRAB (Go REST API Boilerplate)** is a production-ready Go REST API starter with Clean Architecture, comprehensive testing (89.81% coverage), and Docker-first development workflow.

### Technology Stack

- **Go**: Check version with `go version`
- **Gin**: HTTP router framework
- **GORM**: PostgreSQL ORM
- **PostgreSQL**: Check version with `make exec-db` then `psql --version`
- **Docker**: Check version with `docker --version`
- **golang-migrate**: Database migration tool
- **JWT**: Authentication with refresh token rotation
- **Air**: Hot-reload development
- **golangci-lint**: Code quality enforcement
- **Swagger**: OpenAPI documentation

### Documentation

- **Main Docs**: https://vahiiiid.github.io/go-rest-api-docs/
- **Repository**: https://github.com/vahiiiid/go-rest-api-boilerplate
- **Issues**: https://github.com/vahiiiid/go-rest-api-boilerplate/issues

---

## 🏗️ Architecture

### Clean Architecture Pattern

GRAB strictly follows Clean Architecture with clear layer separation:

```
Handler (HTTP) → Service (Business Logic) → Repository (Database)
```

**Domain Structure**:
```
internal/<domain>/
├── model.go       # GORM models with database tags
├── dto.go         # Data Transfer Objects (API contracts)
├── repository.go  # Database access interface + implementation
├── service.go     # Business logic interface + implementation
├── handler.go     # HTTP handlers with Gin + Swagger annotations
└── *_test.go      # Unit and integration tests
```

**Reference Implementation**: See `internal/user/` for complete domain example.

**Key Rules**:
- Handlers only handle HTTP concerns (bind, validate, respond)
- Services contain all business logic
- Repositories only interact with database
- Never skip layers or cross boundaries

---

## 🚀 Development Workflow

### Docker-First Approach

**Important**: Developers run `make` commands on the host machine. The Makefile automatically detects if Docker containers are running and executes commands in the appropriate context.

**No need to manually enter containers** - the Makefile handles this transparently.

```bash
# Start all containers
make up

# All commands below auto-detect Docker and execute accordingly
make test              # Run tests
make lint              # Run linting
make lint-fix          # Auto-fix linting issues
make swag              # Generate Swagger documentation
make migrate-up        # Apply database migrations
make logs              # View container logs
```

### Pre-Commit Checklist

Always run before committing:
```bash
make lint-fix    # Auto-fix linting issues
make lint        # Verify no remaining issues
make test        # Run all tests
make swag        # Update Swagger if API changed
```

---

## 📝 Common Tasks

### Adding a New Domain/Entity

**Step-by-step process**:

1. **Create directory structure**:
   ```bash
   mkdir -p internal/<domain>
   ```

2. **Create model** (`internal/<domain>/model.go`):
   ```go
   package <domain>
   
   import (
       "time"
       "gorm.io/gorm"
   )
   
   type <Entity> struct {
       ID        uint           `gorm:"primarykey" json:"id"`
       Name      string         `gorm:"not null" json:"name"`
       UserID    uint           `gorm:"not null" json:"user_id"`
       CreatedAt time.Time      `json:"created_at"`
       UpdatedAt time.Time      `json:"updated_at"`
       DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
   }
   ```

3. **Create DTOs** (`internal/<domain>/dto.go`):
   ```go
   package <domain>
   
   type Create<Entity>Request struct {
       Name string `json:"name" binding:"required,min=3,max=200"`
   }
   
   type Update<Entity>Request struct {
       Name string `json:"name" binding:"omitempty,min=3,max=200"`
   }
   
   type <Entity>Response struct {
       ID        uint      `json:"id"`
       Name      string    `json:"name"`
       UserID    uint      `json:"user_id"`
       CreatedAt time.Time `json:"created_at"`
   }
   ```

4. **Create repository** (`internal/<domain>/repository.go`):
   ```go
   package <domain>
   
   import (
       "context"
       "gorm.io/gorm"
   )
   
   type Repository interface {
       Create(ctx context.Context, entity *<Entity>) error
       FindByID(ctx context.Context, id uint) (*<Entity>, error)
       Update(ctx context.Context, entity *<Entity>) error
       Delete(ctx context.Context, id uint) error
   }
   
   type repository struct {
       db *gorm.DB
   }
   
   func NewRepository(db *gorm.DB) Repository {
       return &repository{db: db}
   }
   ```

5. **Create service** (`internal/<domain>/service.go`):
   ```go
   package <domain>
   
   import "context"
   
   type Service interface {
       Create<Entity>(ctx context.Context, userID uint, req *Create<Entity>Request) (*<Entity>Response, error)
       Get<Entity>(ctx context.Context, userID, id uint) (*<Entity>Response, error)
       Update<Entity>(ctx context.Context, userID, id uint, req *Update<Entity>Request) (*<Entity>Response, error)
       Delete<Entity>(ctx context.Context, userID, id uint) error
   }
   
   type service struct {
       repo Repository
   }
   
   func NewService(repo Repository) Service {
       return &service{repo: repo}
   }
   ```

6. **Create handler** (`internal/<domain>/handler.go`):
   ```go
   package <domain>
   
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
   
   // @Summary Create entity
   // @Tags <domain>
   // @Accept json
   // @Produce json
   // @Security BearerAuth
   // @Param request body Create<Entity>Request true "Entity data"
   // @Success 201 {object} errors.Response{success=bool,data=<Entity>Response}
   // @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
   // @Router /api/v1/<domain> [post]
   func (h *Handler) Create<Entity>(c *gin.Context) {
       userID := contextutil.GetUserID(c)
       
       var req Create<Entity>Request
       if err := c.ShouldBindJSON(&req); err != nil {
           _ = c.Error(apiErrors.FromGinValidation(err))
           return
       }
       
       result, err := h.service.Create<Entity>(c.Request.Context(), userID, &req)
       if err != nil {
           _ = c.Error(apiErrors.InternalServerError(err))
           return
       }
       
       c.JSON(http.StatusCreated, apiErrors.Success(result))
   }
   ```

7. **Create database migration**:
   ```bash
   make migrate-create NAME=create_<table>_table
   ```
   
   Edit the generated `.up.sql` file:
   ```sql
   BEGIN;
   
   CREATE TABLE IF NOT EXISTS <table> (
       id SERIAL PRIMARY KEY,
       name VARCHAR(200) NOT NULL,
       user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       deleted_at TIMESTAMP
   );
   
   CREATE INDEX idx_<table>_user_id ON <table>(user_id);
   CREATE INDEX idx_<table>_deleted_at ON <table>(deleted_at);
   
   COMMIT;
   ```
   
   Edit the `.down.sql` file:
   ```sql
   BEGIN;
   
   DROP TABLE IF EXISTS <table>;
   
   COMMIT;
   ```

8. **Register routes** in `internal/server/router.go`:
   ```go
   // Initialize components
   <domain>Repo := <domain>.NewRepository(db)
   <domain>Service := <domain>.NewService(<domain>Repo)
   <domain>Handler := <domain>.NewHandler(<domain>Service)
   
   // Register routes (authenticated endpoints)
   <domain>Group := v1.Group("/<domain>")
   <domain>Group.Use(auth.AuthMiddleware(authService))
   {
       <domain>Group.POST("", <domain>Handler.Create<Entity>)
       <domain>Group.GET("/:id", <domain>Handler.Get<Entity>)
       <domain>Group.PUT("/:id", <domain>Handler.Update<Entity>)
       <domain>Group.DELETE("/:id", <domain>Handler.Delete<Entity>)
   }
   ```

9. **Write tests** for all layers

10. **Apply changes**:
    ```bash
    make migrate-up      # Apply migration
    make test            # Run tests
    make lint            # Check code quality
    make swag            # Update Swagger docs
    ```

### Database Migrations

**Naming Convention**: `YYYYMMDDHHMMSS_verb_noun_table`

**Examples**:
- `20251025225126_create_users_table`
- `20251028000000_create_refresh_tokens_table`
- `20251210120000_add_avatar_to_users_table`
- `20251215143000_add_index_to_users_email`

**Commands**:
```bash
make migrate-create NAME=create_todos_table    # Create new migration
make migrate-up                                 # Apply all pending
make migrate-down                               # Rollback one
make migrate-status                             # Check status
make migrate-force VERSION=<version>           # Force version
```

**Best Practices**:
- Wrap in `BEGIN;` / `COMMIT;` transactions
- Use `IF NOT EXISTS` for safety
- Create indexes for foreign keys
- Create indexes for frequently queried columns
- Always write corresponding `.down.sql`
- Test rollback before committing

---

## 🔐 Authentication & Authorization

### Getting Current User

```go
import "github.com/vahiiiid/go-rest-api-boilerplate/internal/contextutil"

func (h *Handler) SomeHandler(c *gin.Context) {
    userID := contextutil.GetUserID(c)
    userEmail := contextutil.GetEmail(c)
    userName := contextutil.GetUserName(c)
    userRoles := contextutil.GetRoles(c)
    isAdmin := contextutil.IsAdmin(c)
    hasRole := contextutil.HasRole(c, "moderator")
    
    // Use user information...
}
```

### Protecting Routes

```go
import "github.com/vahiiiid/go-rest-api-boilerplate/internal/middleware"

// Require authentication (handled by auth package middleware)
// RequireAuth is from auth package, RequireRole/RequireAdmin are from middleware package

// Admin-only route
v1.Use(middleware.RequireAdmin()).
   POST("/admin/users", userHandler.CreateUser)

// Specific role required
v1.Use(middleware.RequireRole("admin")).
   POST("/admin/reports", reportHandler.CreateReport)

// Note: RequireRole and RequireAdmin already check authentication internally
```

---

## ❌ Error Handling

GRAB uses centralized error handling:

```go
import (
    "errors"
    apiErrors "github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
)

// Validation errors (automatic field extraction)
if err := c.ShouldBindJSON(&req); err != nil {
    _ = c.Error(apiErrors.FromGinValidation(err))
    return
}

// Standard errors
_ = c.Error(apiErrors.NotFound("Resource not found"))
_ = c.Error(apiErrors.Unauthorized("Authentication required"))
_ = c.Error(apiErrors.Forbidden("Access denied"))
_ = c.Error(apiErrors.BadRequest("Invalid request data"))
_ = c.Error(apiErrors.Conflict("Resource already exists"))

// Service/repository errors - check specific errors first
user, err := h.service.CreateUser(ctx, req)
if err != nil {
    // Check for known specific errors first
    if errors.Is(err, ErrEmailExists) {
        _ = c.Error(apiErrors.Conflict("Email already exists"))
        return
    }
    if errors.Is(err, ErrUserNotFound) {
        _ = c.Error(apiErrors.NotFound("User not found"))
        return
    }
    // Wrap unknown errors
    _ = c.Error(apiErrors.InternalServerError(err))
    return
}
```

---

## 🧪 Testing

### Test Structure

Use table-driven tests:

```go
func TestService_CreateEntity(t *testing.T) {
    tests := []struct {
        name        string
        userID      uint
        request     *CreateEntityRequest
        setupMocks  func(*MockRepository)
        expectError bool
        errorType   error
    }{
        {
            name:   "success",
            userID: 1,
            request: &CreateEntityRequest{Name: "Test"},
            setupMocks: func(m *MockRepository) {
                m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
            },
            expectError: false,
        },
        {
            name:        "validation_error",
            userID:      1,
            request:     &CreateEntityRequest{Name: ""},
            expectError: true,
        },
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
            result, err := service.CreateEntity(context.Background(), tt.userID, tt.request)
            
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

### Test Commands

```bash
make test              # Run all tests
make test-coverage     # Generate coverage report (opens in browser)
make test-verbose      # Run with verbose output
```

---

## 📚 Swagger/OpenAPI Documentation

### Annotations

```go
// @Summary Create user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 401 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
    // Handler implementation
}
```

### Update Documentation

```bash
make swag    # Regenerate Swagger docs

# View at: http://localhost:8080/swagger/index.html
```

---

## ⚙️ Configuration

### Configuration Files

- `configs/config.yaml` - Base configuration
- `configs/config.development.yaml` - Development overrides
- `configs/config.staging.yaml` - Staging overrides
- `configs/config.production.yaml` - Production overrides

### Environment Variables

Override any config value with environment variables:

```bash
DATABASE_PASSWORD=secret      # Overrides database.password
JWT_SECRET=secret            # Overrides jwt.secret
APP_ENVIRONMENT=production   # Overrides app.environment
RATE_LIMIT_ENABLED=true      # Overrides ratelimit.enabled
```

**Full Configuration Guide**: https://vahiiiid.github.io/go-rest-api-docs/CONFIGURATION/

---

## 🎯 Out-of-the-Box Features

GRAB includes these production-ready features:

1. **JWT Authentication** - Access + refresh tokens with rotation ([Docs](https://vahiiiid.github.io/go-rest-api-docs/AUTHENTICATION/))
2. **RBAC** - Role-based access control ([Docs](https://vahiiiid.github.io/go-rest-api-docs/RBAC/))
3. **Database Migrations** - Versioned SQL migrations ([Docs](https://vahiiiid.github.io/go-rest-api-docs/MIGRATIONS_GUIDE/))
4. **Health Checks** - `/health`, `/live`, `/ready` endpoints ([Docs](https://vahiiiid.github.io/go-rest-api-docs/HEALTH_CHECKS/))
5. **Rate Limiting** - Token bucket algorithm ([Docs](https://vahiiiid.github.io/go-rest-api-docs/RATE_LIMITING/))
6. **Structured Logging** - JSON logs with context ([Docs](https://vahiiiid.github.io/go-rest-api-docs/LOGGING/))
7. **API Response Format** - Standardized responses ([Docs](https://vahiiiid.github.io/go-rest-api-docs/API_RESPONSE_FORMAT/))
8. **Error Handling** - Centralized error management ([Docs](https://vahiiiid.github.io/go-rest-api-docs/ERROR_HANDLING/))
9. **Graceful Shutdown** - Clean termination ([Docs](https://vahiiiid.github.io/go-rest-api-docs/GRACEFUL_SHUTDOWN/))
10. **Swagger/OpenAPI** - Auto-generated docs ([Docs](https://vahiiiid.github.io/go-rest-api-docs/SWAGGER/))
11. **Context Helpers** - Request utilities ([Docs](https://vahiiiid.github.io/go-rest-api-docs/CONTEXT_HELPERS/))

---

## 🔧 Quick Reference

### Essential Commands

| Task | Command |
|------|---------|
| Start development | `make up` |
| Stop containers | `make down` |
| Run tests | `make test` |
| Lint code | `make lint` |
| Fix linting | `make lint-fix` |
| Create migration | `make migrate-create NAME=<name>` |
| Apply migrations | `make migrate-up` |
| Rollback migration | `make migrate-down` |
| Migration status | `make migrate-status` |
| Update Swagger | `make swag` |
| View logs | `make logs` |
| Enter app container | `make exec` |
| Enter DB container | `make exec-db` |
| Clean restart | `make down && make up` |
| Health check | `curl localhost:8080/health` |
| View all commands | `make help` |

### Project Structure

```
go-rest-api-boilerplate/
├── .github/              # GitHub workflows, templates
├── .cursor/              # Cursor AI rules
├── .windsurf/            # Windsurf AI rules
├── api/                  # API documentation
├── cmd/                  # Application entry points
├── configs/              # Configuration files
├── internal/             # Application code
│   ├── auth/             # Authentication
│   ├── config/           # Config management
│   ├── contextutil/      # Context helpers
│   ├── db/               # Database setup
│   ├── errors/           # Error handling
│   ├── health/           # Health checks
│   ├── middleware/       # HTTP middleware
│   ├── migrate/          # Migration logic
│   ├── server/           # Router setup
│   └── user/             # User domain (reference)
├── migrations/           # SQL migration files
├── scripts/              # Helper scripts
├── tests/                # Integration tests
├── Dockerfile            # Docker image
├── docker-compose.yml    # Development compose
├── Makefile              # Development commands
├── AGENTS.md             # This file
└── README.md             # Project overview
```

---
## 📚 Migration History & Important Changes

### CorretorPrincipal Migration (January 2026)

**Breaking Change**: Migrated from `Organizacao` direct relationship to `CorretorPrincipal` relationship in Imovel.

**Architecture Change**:
```
# Before:
Imovel → Organizacao

# After:
Imovel → CorretorPrincipal → Organizacao
         ↓
       Anexo (foto)
```

**Key Changes**:
1. Replaced `organizacao_id` with `corretor_principal_id` in imoveis table
2. Created `corretores_principais` table with relationship to `organizacoes`
3. Added `foto_id` to corretores_principais (references anexos)
4. Updated anexos table to support `corretor_principal_id`

**Migration Order**:
- Anexos created WITHOUT foreign keys initially (120350)
- CorretoresPrincipais created after anexos (120400)
- Foreign keys added to anexos in final migration (121200)

**Code Updates Required**:
- `repository.go`: Map `imovel.CorretorPrincipal` instead of `imovel.Organizacao`
- `service.go`: Replace all `OrganizacaoID` references with `CorretorPrincipalID`
- `dto.go`: Update request/response DTOs
- Methods renamed: `ListByOrganizacao` → `ListByCorretorPrincipal`

**TableName() Methods Required**:
```go
func (Organizacao) TableName() string {
    return "organizacoes"
}

func (CorretorPrincipal) TableName() string {
    return "corretores_principais"
}
```

**Testing After Migration**:
```bash
# 1. Drop and recreate database
make migrate-down  # or TRUNCATE tables
make migrate-up

# 2. Verify structure
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\d imoveis"
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\d corretores_principais"

# 3. Test relationships
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "
SELECT i.id, i.codigo, cp.nome as corretor, o.nome as organizacao 
FROM imoveis i 
LEFT JOIN corretores_principais cp ON i.corretor_principal_id = cp.id 
LEFT JOIN organizacoes o ON cp.organizacao_id = o.id;"
```

---
## � External Property Import System

### Overview

The import system synchronizes properties from an external API (`dev-api-backend.pi8.com.br`) with full relationship management and intelligent anexos synchronization.

### Key Features

1. **Complete Relationship Updates**: All relationships are updated/created on every import:
   - ✅ Empreendimento (Development/Building)
   - ✅ PrecoVenda (Selling Price)
   - ✅ PrecoAluguel (Rental Price)
   - ✅ CorretorPrincipal (Main Broker)
   - ✅ Organizacao (Organization)
   - ✅ Endereco (Address)

2. **Smart Anexos Sync**: DELETE + INSERT strategy ensures:
   - Removed images from external API are deleted locally
   - New images are automatically added
   - Image order is preserved
   - No duplicate anexos

### Import Command

```bash
# Inside Docker container
docker exec -it triiio_app go run cmd/importimoveis/main.go

# Or using Makefile (if configured)
make import-properties
```

### Architecture

**File**: `internal/imoveis/import_service.go`

**Key Functions**:

1. **`ImportPublishedProperties()`** - Main entry point
   - Fetches list of published properties
   - Processes each property with upsert logic
   - Returns: `created count, updated count, failed count`

2. **`upsertImovelAndRelationships()`** - Core orchestration
   - **CRITICAL**: Processes relationships BEFORE create/update
   - Handles both new and existing properties
   - Ensures all FKs are set correctly

3. **`upsertEmpreendimento()`** - Building/Development
   - Searches by `id_integracao` (external ID)
   - Updates existing or creates new
   - Uses `Omit()` to avoid problematic date fields

4. **`upsertPrecoVenda()` & `upsertPrecoAluguel()`** - Pricing
   - Searches by `id_integracao`
   - Updates all pricing flags and amounts

5. **`upsertCorretorPrincipal()`** - Main Broker
   - Searches by `id_integracao`
   - Calls `upsertOrganizacao()` first
   - Uses `Omit("FotoID")` to prevent FK violations

6. **`upsertOrganizacao()`** - Organization
   - Searches by `nome` (name - unique)
   - Updates profile if changed

7. **`syncAnexosFromImages()`** - Image Synchronization
   - **Step 1**: DELETE all existing anexos for property
   - **Step 2**: INSERT all current anexos from external API
   - Guarantees sync with external data

### Implementation Pattern

```go
// CORRECT: Process relationships BEFORE create/update
func (is *importService) upsertImovelAndRelationships(...) {
    // 1. Upsert all relationships first
    var empreendimentoID uint
    if ext.Empreendimento != nil {
        empID, err := is.upsertEmpreendimento(ctx, ext.Empreendimento)
        if err == nil {
            empreendimentoID = empID
        }
    }
    
    // Same for PrecoVenda, PrecoAluguel, CorretorPrincipal...
    
    // 2. Create or Update imovel with relationship IDs
    if isUpdate {
        updateReq := &UpdateImovelRequest{
            EmpreendimentoID: &empreendimentoID,  // Use pointers for optional fields
            PrecoVendaID: &precoVendaID,
            // ...
        }
        imovelResp, err = is.service.UpdateImovel(ctx, imovelID, updateReq)
    } else {
        createReq := &CreateImovelRequest{
            EmpreendimentoID: empreendimentoID,  // uint directly for required fields
            // ...
        }
        imovelResp, err = is.service.CreateImovel(ctx, createReq)
    }
    
    // 3. Sync anexos (DELETE old + INSERT new)
    if err := is.syncAnexosFromImages(ctx, imovelID, ext.Imagens); err != nil {
        // Handle error
    }
}
```

### Common Issues & Solutions

#### 1. Foreign Key Constraint Violations

**Problem**: `insert or update on table "X" violates foreign key constraint`

**Solutions**:
- Use `Omit("FieldName")` when creating entities with optional FK fields that have zero values
- Example: `db.Omit("FotoID").Create(&corretor)` prevents inserting `foto_id = 0`
- For pointer fields, set to `nil` instead of zero value

#### 2. GORM Table Naming (Portuguese Plurals)

**Problem**: GORM creates wrong table names (`organizacaos` instead of `organizacoes`)

**Solution**: Add `TableName()` method
```go
func (Organizacao) TableName() string {
    return "organizacoes"
}

func (CorretorPrincipal) TableName() string {
    return "corretores_principais"
}
```

#### 3. Partial Relationship Updates

**Problem**: Only some relationships were updating, others stayed stale

**Solution**: Move ALL relationship upsert logic BEFORE the if/else block that handles create vs update. This ensures both code paths use the same relationship processing.

#### 4. Anexos Not Syncing with External API

**Problem**: Removed images stayed in database, new images not showing

**Solution**: Replaced incremental add logic with DELETE + INSERT strategy in `syncAnexosFromImages()`

### Testing Commands

```bash
# 1. Start containers
docker ps --format "{{.Names}}"  # Verify containers: triiio_app, triiio_db

# 2. Run import
docker exec -it triiio_app go run cmd/importimoveis/main.go

# 3. Verify results
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "
SELECT i.id, i.codigo, i.titulo, e.preco as preco_venda, 
       cp.nome as corretor, o.nome as organizacao 
FROM imoveis i 
LEFT JOIN preco_vendas e ON i.preco_venda_id = e.id 
LEFT JOIN corretores_principais cp ON i.corretor_principal_id = cp.id 
LEFT JOIN organizacoes o ON cp.organizacao_id = o.id 
WHERE i.deleted_at IS NULL 
ORDER BY i.id;"

# 4. Check anexos count
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "
SELECT imovel_id, COUNT(*) as total_anexos 
FROM anexos 
WHERE deleted_at IS NULL 
GROUP BY imovel_id 
ORDER BY imovel_id;"

# 5. Verify counts
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "
SELECT 'Imóveis:' as tipo, COUNT(*)::text as total FROM imoveis 
WHERE deleted_at IS NULL
UNION ALL 
SELECT 'Corretores:', COUNT(*)::text FROM corretores_principais 
WHERE deleted_at IS NULL
UNION ALL 
SELECT 'Organizações:', COUNT(*)::text FROM organizacoes 
WHERE deleted_at IS NULL
UNION ALL 
SELECT 'Anexos:', COUNT(*)::text FROM anexos 
WHERE deleted_at IS NULL;"
```

### Expected Output

**Successful Import**:
```
import completed: X created, Y updated, 0 failed
Synced N anexos for property ID Z
```

**Database Verification**:
```
 id | codigo     | titulo  | preco_venda | corretor | organizacao 
----+------------+---------+-------------+----------+-------------
 11 | AP00000345 | Teste 4 |      750000 | PCA      | PCA
 12 | CA00000342 | Teste 1 |     1950000 | PCA      | PCA
...
```

### Development Workflow

1. **Compile Check**:
   ```bash
   cd /home/maxcosta/triiio/triiio-api && go build ./...
   ```

2. **Run Import** (in container):
   ```bash
   docker exec -it triiio_app go run cmd/importimoveis/main.go
   ```

3. **Verify Data** (see SQL commands above)

4. **Re-run Import** to test update logic:
   - Should show `0 created, N updated, 0 failed`
   - All relationships should be current
   - Anexos count should match external API

### Debugging Tips

- Check GORM logs in terminal output (shows SQL queries)
- Look for `UPDATE "table"` vs `INSERT INTO "table"` to verify upsert behavior
- Count anexos before/after to verify DELETE + INSERT
- Check `deleted_at IS NULL` in all queries (soft deletes)
- Verify `id_integracao` values match external API IDs

### Import Behavior

**First Import (Create)**:
- All properties are **created** since `id_integracao` doesn't exist
- All relationships (empreendimentos, precos) are **created**
- Attachments are **added**
- Result: `X created, 0 updated, 0 failed`

**Subsequent Imports (Update)**:
- Existing properties are **detected** by `id_integracao`
- Properties are **updated** with latest data from external API
- Relationships are **updated** with changed fields only
- Duplicate attachments are **skipped** (URL checking)
- Result: `0 created, X updated, 0 failed`

**What Gets Updated**:
- Business data: titulo, descricao, metragem, num_quartos, preco, etc.
- Relationships: endereco_id, empreendimento_id, preco_venda_id, corretor_principal_id, etc.
- updated_at timestamp (automatic)

**Preserved Fields**:
- created_at timestamp
- id (primary key)
- id_integracao (external ID mapping)

### External ID Mapping Strategy

All relationship entities from external API use `id_integracao` field:

**Format**: String representation of external ID (e.g., "514", "857")

**Entities with id_integracao**:
- Imovel (property)
- Empreendimento (enterprise/development)
- PrecoVenda (selling price)
- PrecoAluguel (rental price)
- CorretorPrincipal (main broker)

**Benefits**:
- Duplicate prevention: Check if record exists before creating
- Update synchronization: Update local data when external changes
- Traceability: Track which local records map to external API

**Upsert Logic**:
```go
// Check if exists by id_integracao
existing, err := repo.FindByIdIntegracao(ctx, fmt.Sprintf("%d", externalID))
if err == nil {
    // Update existing record
    existing.Field = newValue
    db.Save(&existing)
} else {
    // Create new record
    new := &Entity{
        IdIntegracao: fmt.Sprintf("%d", externalID),
        Field: value,
    }
    db.Create(&new)
}
```

---

## 📝 Sliders Domain

### Model Structure

**Slider**: Container for slider items
- ID, Name, Type (slideshow/carousel/static), Location
- Has many SliderItems

**SliderItem**: Individual slider entry
- ID, SliderID, ImageURL, LinkURL, Order, Tags
- **Content** (renamed from Caption on 2026-01-30)
- **Titulo** (renamed from Finalidade on 2026-01-30)

### Field Renaming (January 30, 2026)

**Changes**:
- `caption` → `content` - Content text for the slider item
- `finalidade` → `titulo` - Title/heading for the slider item

**Development Approach Used**:
- ✅ Modified original migration: `migrations/20260107205003_create_sliders_table.up.sql`
- ✅ Dropped and recreated tables with correct fields
- ✅ No additional migration files created (development phase)

**Files Modified**:
- `internal/sliders/model.go` - Model definition
- `internal/sliders/dto.go` - Request/Response DTOs
- `internal/sliders/service.go` - Business logic layer
- `internal/sliders/repository.go` - Database layer (Select fields)
- `migrations/20260107205003_create_sliders_table.up.sql` - Original migration updated
- `api/docs/*` - Swagger documentation regenerated

**Verification**:
```bash
# Check database structure
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\d slider_items"

# Verify columns exist
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'slider_items' AND column_name IN ('content', 'titulo');"
```

**Recreate Tables** (if needed in development):
```bash
# Drop tables
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "DROP TABLE IF EXISTS slider_items CASCADE; DROP TABLE IF EXISTS sliders CASCADE;"

# Recreate from migration
cd /home/maxcosta/triiio/triiio-api && docker exec -i triiio_db psql -U triiio_user -d triiio_backend < migrations/20260107205003_create_sliders_table.up.sql
```

---

## 💡 Best Practices for AI Assistants

1. **Reference Existing Code**: Always check `internal/user/` for patterns before creating new domains
2. **Follow Clean Architecture**: Never skip Handler → Service → Repository layers
3. **Use Context Helpers**: Import `contextutil` for user information from JWT
4. **Minimal Comments**: Write self-documenting code, comment WHY not WHAT
5. **Test Coverage**: Maintain 85%+ test coverage for all new code
6. **Check Makefile**: All development commands available in `make help`
7. **Read Documentation**: Comprehensive guides at https://vahiiiid.github.io/go-rest-api-docs/
8. **Version Checking**: Show commands to check versions, don't hardcode
9. **Docker-First**: Assume Docker containers running, use `make` commands
10. **Migration Naming**: Follow `YYYYMMDDHHMMSS_verb_noun_table` pattern
11. **Import Logic**: Always upsert relationships BEFORE creating/updating main entity
12. **Anexos Sync**: Use DELETE + INSERT pattern to ensure data freshness

### ⚠️ CRITICAL: Document Your Learning

**After completing ANY task, ALWAYS update this documentation with:**

1. **What worked**: Commands that succeeded, patterns that solved problems
2. **What failed**: Errors encountered and their solutions
3. **New patterns**: Any new architectural patterns or solutions discovered
4. **Gotchas**: Pitfalls to avoid in future similar tasks

**Why this matters**:
- Reduces repetition of same mistakes
- Builds project-specific knowledge base
- Improves context for future AI assistants
- Creates single source of truth for solutions

**Where to document**:
- Add to relevant section in AGENTS.md or .github/copilot-instructions.md
- Update "Common Issues & Solutions" sections
- Add to "Testing Commands" if verification steps were created
- Update "Expected Output" with real examples from your work

**Example**:
```markdown
### Problem: Foreign Key Violation on foto_id
**Solution**: Use `Omit("FotoID")` when creating CorretorPrincipal
**Command that worked**: `db.Omit("FotoID").Create(&corretor)`
**Verification**: `docker exec -i triiio_db psql ... "SELECT * FROM corretores_principais"`
```

---

## 🔗 Additional Resources

- **Full Documentation**: https://vahiiiid.github.io/go-rest-api-docs/
- **GitHub Repository**: https://github.com/vahiiiid/go-rest-api-boilerplate
- **Issue Tracker**: https://github.com/vahiiiid/go-rest-api-boilerplate/issues
- **Discussions**: https://github.com/vahiiiid/go-rest-api-boilerplate/discussions
- **Development Guide**: https://vahiiiid.github.io/go-rest-api-docs/DEVELOPMENT_GUIDE/
- **Quick Reference**: https://vahiiiid.github.io/go-rest-api-docs/QUICK_REFERENCE/

---

**Version**: v2.0.0  
**Last Updated**: 2025-12-10  
**Maintained By**: GRAB Contributors
