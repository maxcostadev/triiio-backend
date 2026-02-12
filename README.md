<div align="center">

![TRIIIO Logo](https://www.triiio.com.br/wp-content/uploads/2018/12/logo-triiiio-site-1.jpg)

# TRIIIO Backend API

API REST de gerenciamento de imóveis com integração externa, pronta para produção.

*Backend robusto construído com Go, Clean Architecture, autenticação JWT, RBAC e sincronização inteligente com APIs externas.*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)

**[🚀 Quick Start](#-quick-start)** • **[📖 Documentação](#-documentação)** • **[🏠 Funcionalidades](#-funcionalidades-principais)**

</div>

---

## 🏠 Funcionalidades Principais

### Sistema de Gerenciamento de Imóveis
- **CRUD Completo** de imóveis com relacionamentos complexos
- **Importação Inteligente** da API externa (dev-api-backend.pi8.com.br)
- **Sistema de Anexos** com gerenciamento de imagens
- **Endereços Geolocalizados** com integração
- **Empreendimentos e Plantas** associados aos imóveis
- **Preços de Venda e Aluguel** com múltiplas condições

### Arquitetura e Segurança
- **Clean Architecture** (Handler → Service → Repository)
- **JWT Authentication** com refresh token rotation
- **RBAC** (Role-Based Access Control)
- **Rate Limiting** e proteção contra abusos
- **Structured Logging** com rastreamento de requisições
- **Validação de Dados** em todas as camadas

### Infraestrutura
- **Docker-First Development** com hot-reload em 2 segundos
- **PostgreSQL 15** com migrações versionadas
- **Health Checks** Kubernetes-ready
- **Swagger/OpenAPI** para documentação interativa
- **Graceful Shutdown** para deploys sem downtime

---

## 🎯 Por Que TRIIIO Backend?

API REST completa para gestão de imóveis com sincronização externa, construída com as melhores práticas do ecossistema Go.

```bash
make up          # Inicia containers com hot-reload
make migrate-up  # Aplica migrações do banco
make import-properties  # Importa imóveis da API externa
```

**Recursos Implementados:**

✅ **Clean Architecture** — Separação clara de responsabilidades  
✅ **Importação Externa** — Sistema inteligente de sincronização com API externa  
✅ **Mapeamento id_integracao** — Rastreamento de registros externos sem duplicação  
✅ **JWT Authentication** — OAuth 2.0 compliant com refresh tokens  
✅ **RBAC** — Controle de acesso baseado em roles  
✅ **Migrações Versionadas** — PostgreSQL com controle total de schema  
✅ **Swagger/OpenAPI** — Documentação interativa auto-gerada  
✅ **Logging Estruturado** — JSON logs com request IDs  
✅ **Error Handling** — Respostas padronizadas e machine-readable  
✅ **Docker Production-Ready** — Multi-stage builds otimizados  
✅ **Health Checks** — Kubernetes-ready probes  
✅ **Hot-Reload** — Desenvolvimento ágil com Air (2 segundos!)  

### 🏆 Seguindo Padrões Go

Arquitetura baseada em **[official Go project layout](https://go.dev/doc/modules/layout)** e **[golang-standards/project-layout](https://github.com/golang-standards/project-layout)**.

### 🎯 Ideal Para

- 🏢 **Gestão Imobiliária** — Sistema completo de cadastro e sincronização  
- 🔄 **Integração de APIs** — Importação e sincronização de dados externos  
- 📊 **Dados Relacionais Complexos** — Imóveis, endereços, preços, anexos  
- 🚀 **Produção** — Pronto para deploy com Docker e Kubernetes

---

## 🚀 Quick Start

Inicie a API em **menos de 2 minutos**:

### Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)

### Setup Rápido ⚡

```bash
# 1. Clone o repositório
git clone <seu-repositorio>
cd triiio-backend

# 2. Inicie os containers
make up

# 3. Aplique as migrações
make migrate-up

# 4. Verifique o status
make migrate-status
```

**🎉 Pronto!** Sua API está rodando em:

- **API Base URL:** http://localhost:8080/api/v1
- **Swagger UI:** http://localhost:8080/swagger/index.html
- **Health Checks:** http://localhost:8080/health
  - Liveness: http://localhost:8080/health/live
  - Readiness: http://localhost:8080/health/ready

### Configuração Inicial

**Criar Usuário Admin:**

```bash
make create-admin              # Interativo: solicita email, nome, senha
make promote-admin ID=1        # Promove usuário existente a admin
```

**Importar Imóveis da API Externa:**

```bash
make import-properties         # Sincroniza com dev-api-backend.pi8.com.br
```

### Containers Docker

O projeto utiliza os seguintes containers:
- **triiio_app** - Aplicação Go com hot-reload
- **triiio_db** - PostgreSQL 15

### Banco de Dados

- **Host:** localhost (ou triiio_db dentro do Docker)
- **Porta:** 5432
- **Database:** triiio_backend
- **Usuário:** triiio_user
- **Senha:** Configurada no arquivo `.env`

---

## ✨ Testando a API

### Documentação Interativa com Swagger

Abra [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) para explorar e testar todos os endpoints de forma interativa.

### Coleção do Postman

Importe a coleção pré-configurada localizada em `api/postman_collection.json` com exemplos de requisições e testes prontos.

---

## 💎 Diferenciais do TRIIIO Backend

Most boilerplates give you code. **GRAB gives you a professional development workflow.**

#### 🔐 Authentication That Actually Works

- **OAuth 2.0 BCP compliant** — JWT-based auth (HS256) with refresh token rotation and automatic reuse detection
- **Enhanced security** — Refresh tokens with family tracking, secure token invalidation, and breach detection
- **Context helpers** — Type-safe user extraction (no more casting nightmares)
- **Password security** — Bcrypt hashing with best-practice cost factor
- **Rate limiting** — Token-bucket protection against abuse built-in

#### 🔑 Role-Based Access Control (RBAC)

- **Many-to-many architecture** — Flexible roles system with extensible permissions
- **Secure admin CLI** — Interactive admin creation with strong password enforcement (no defaults in code)
- **JWT-integrated authorization** — Roles embedded in tokens for server-side validation
- **Protected endpoints** — Middleware-based access control (RequireRole, RequireAdmin)
- **Three-endpoint pattern** — `/auth/me` (current user), `/users/:id` (specific), `/users` (admin list)
- **Paginated user management** — Admin-only user listing with filtering and search

#### 🏠 Smart External API Integration

- **id_integracao mapping system** — Track external API records with unique identifiers (`EXT_{id}` format)
- **Upsert logic** — Automatically creates new records or updates existing ones based on external ID mapping
- **Intelligent duplicate prevention** — No more duplicate data on re-imports
- **Relationship synchronization** — Updates empreendimentos, prices, and all related entities
- **Attachment deduplication** — URL-based checking prevents duplicate images
- **Robust error handling** — Foreign key constraint management, date field validation
- **Preserves audit trails** — Updates business fields while preserving `created_at` timestamps

**Example Import Flow:**
```bash
# First import: Creates 100 properties
docker exec triiio_app go run cmd/importimoveis/main.go
# Result: 100 created, 0 updated, 0 failed

# Re-import com dados atualizados: Atualiza propriedades existentes
make import-properties
# Resultado: 0 criados, 100 atualizados, 0 falhas
```

**Como funciona:**
- Primeira importação: Cria todos os imóveis (X criados, 0 atualizados)
- Importações subsequentes: Atualiza dados existentes (0 criados, X atualizados)
- Mapeamento por `id_integracao` evita duplicação
- Sincroniza relacionamentos: empreendimentos, preços, endereços, anexos
- Anexos são sincronizados com DELETE + INSERT para garantir consistência

#### 🗄️ Database Setup That Doesn't Fight You

- **PostgreSQL + GORM** — Production-grade ORM with relationship support
- **golang-migrate** — Industry-standard migrations with timestamp versioning
- **Complete migration CLI** — Create, apply, rollback with ease

  ```bash
  make migrate-create NAME=add_posts_table  # Create with timestamp
  make migrate-up                            # Apply all pending
  make migrate-down                          # Rollback last (safe)
  make migrate-down STEPS=3                  # Rollback multiple
  make migrate-status                        # Check current version
  make migrate-goto VERSION=<timestamp>      # Jump to specific version
  ```

- **Safety features** — Confirmation prompts, dirty state detection
- **Transaction support** — BEGIN/COMMIT wrappers for data integrity
- **Connection pooling** — Configured for performance out of the box

#### 🐳 Docker That Saves Your Sanity

- **2-second hot-reload** — Powered by Air, actually works in Docker
- **One command to rule them all** — `make quick-start` handles everything
- **Development & production** — Separate optimized configs
- **Multi-stage builds** — Tiny production images (~20MB)

#### 🏥 Production-Grade Health Checks

- **Kubernetes-ready probes** — Liveness (`/health/live`) and readiness (`/health/ready`) endpoints
- **Database health monitoring** — Response time tracking with pass/warn/fail thresholds
- **RFC-compliant responses** — Following IETF draft standards for health check format
- **Zero-downtime deployments** — Smart readiness checks for load balancer integration
- **Extensible architecture** — Easy to add custom health checkers (Redis, external APIs, etc.)

#### 📚 Documentação

- **Swagger Auto-gerado** — API explorer interativo em `/swagger/index.html`
- **Coleção Postman** — Importe e teste imediatamente de `api/postman_collection.json`

#### 🧪 Tests That Give You Confidence

- **Comprehensive coverage** — Handlers, services, and repositories all tested
- **In-memory SQLite** — No external dependencies for tests
- **Table-driven tests** — Go idiomatic testing patterns
- **CI/CD ready** — GitHub Actions configured and working

#### 📦 Standardized API Responses

- **Consistent envelope format** — All responses wrapped in `{success, data, error, meta}` structure
- **JSend-inspired design** — Industry best practice for API response formatting
- **Type-safe responses** — Predictable structure for frontend integration
- **Metadata support** — Pagination, timestamps, request IDs built-in

#### ⚠️ Error Handling That Makes Sense

- **Structured API errors** — Machine-readable codes (NOT_FOUND, VALIDATION_ERROR, etc.)
- **Detailed error info** — Code, message, details, timestamp, path, request ID
- **Validation details** — Clear field-level error messages for bad requests
- **Centralized middleware** — Single error handler for consistent responses
- **Rate limit errors** — Includes `retry_after` for proper backoff logic

#### 🏗️ Architecture That Scales

- **Clean layers** — Handler → Service → Repository (no shortcuts)
- **Dependency injection** — Proper DI, easy to mock and test
- **Domain-driven** — Organize by feature, not by layer
- **Official Go layout** — Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

---

## 🛠️ Desenvolvimento

### Comandos Principais

#### Docker
```bash
make up              # Inicia containers com hot-reload
make down            # Para containers
make restart         # Reinicia containers
make logs            # Visualiza logs do app
make build           # Reconstrói containers
```

#### Migrações de Banco
```bash
make migrate-create NAME=nome_da_migration  # Cria nova migration
make migrate-up                              # Aplica migrations pendentes
make migrate-down                            # Rollback da última migration
make migrate-down STEPS=3                    # Rollback de 3 migrations
make migrate-status                          # Status atual
make migrate-goto VERSION=20260113120000     # Vai para versão específica
```

#### Testes e Qualidade
```bash
make test              # Executa todos os testes
make test-coverage     # Gera relatório de cobertura
make lint              # Verifica qualidade do código
make lint-fix          # Corrige problemas automaticamente
```

#### Documentação
```bash
make swag              # Gera documentação Swagger
```

#### Gerenciamento de Imóveis
```bash
make import-properties # Importa imóveis da API externa
```

#### Admin
```bash
make create-admin           # Cria novo usuário admin
make promote-admin ID=123   # Promove usuário existente
```

**O que você tem:**

- 🔥 **Hot-reload** — Mudanças refletem em ~2 segundos (Air)
- 📦 **Volume mounts** — Edite no IDE, roda no container
- 🗄️ **PostgreSQL** — Banco na rede Docker interna
- 📚 **Ferramentas pré-instaladas** — Não precisa instalar Go no host

### Estrutura de Migrações

Migrações seguem o padrão `YYYYMMDDHHMMSS_acao_tabela`:

**Exemplos:**
- `20260113120000_create_enderecos_table`
- `20260113120400_create_corretores_principais_table`
- `20260113120900_create_imoveis_table`

**Ordem de criação (importante devido a foreign keys):**
1. Tabelas base (users, roles, organizacoes)
2. Tabelas de relacionamento simples (enderecos, plantas)
3. Anexos (sem FKs iniciais)
4. Corretores principais
5. Preços (venda/aluguel)
6. Empreendimentos
7. Imóveis (referencia todas as anteriores)
8. Foreign keys adicionais em anexos

### Sem Docker

Precisa de Go 1.21+ instalado:

```bash
make build-binary    # Compila para bin/server
make run-binary      # Compila e executa (requer PostgreSQL local)
```

---

## 🚢 Deployment

### Production-Ready From Day One

GRAB includes optimized production builds:

```bash
make docker-up-prod  # Start production containers
```

**What's included:**

- ✅ Multi-stage Docker builds (minimal image size)
- ✅ Production-grade health checks (liveness & readiness probes)
- ✅ Environment-based configuration
- ✅ No development dependencies
- ✅ Production logging

### Deploy Anywhere

Ready for:

- **AWS ECS/Fargate** — Container orchestration
- **Google Cloud Run** — Serverless containers
- **DigitalOcean App Platform** — Platform-as-a-service
- **Kubernetes** — Self-managed orchestration
- **Any VPS** — Using Docker Compose

---

## 📖 Documentação

### Swagger/OpenAPI Interativo

Acesse a documentação interativa em:

**http://localhost:8080/swagger/index.html**

Teste todos os endpoints diretamente pelo navegador.

### Postman Collection

Importe a coleção pré-configurada de `api/postman_collection.json` com exemplos de requisições e testes.

### Estrutura do Projeto

```
triiio-backend/
├── cmd/                      # Entry points da aplicação
│   ├── server/              # Servidor principal
│   ├── migrate/             # CLI de migrações
│   ├── createadmin/         # CLI de criação de admin
│   └── importimoveis/       # Importador de imóveis
├── internal/                # Código da aplicação
│   ├── auth/                # Autenticação JWT
│   ├── config/              # Configuração
│   ├── contextutil/         # Helpers de contexto
│   ├── db/                  # Setup do banco
│   ├── errors/              # Tratamento de erros
│   ├── health/              # Health checks
│   ├── middleware/          # Middlewares HTTP
│   ├── server/              # Setup do router
│   ├── user/                # Domínio de usuários
│   ├── imoveis/             # Domínio de imóveis
│   └── sliders/             # Domínio de sliders
├── migrations/              # Migrações SQL
├── configs/                 # Arquivos de configuração
├── api/                     # Documentação da API
├── scripts/                 # Scripts auxiliares
├── Dockerfile               # Imagem Docker
├── docker-compose.yml       # Compose de desenvolvimento
├── docker-compose.prod.yml  # Compose de produção
├── Makefile                 # Comandos de desenvolvimento
└── README.md                # Este arquivo
```

### Domínios Implementados

#### Imóveis (`internal/imoveis/`)
Sistema completo de gerenciamento de imóveis com:
- CRUD de imóveis
- Importação de API externa
- Gerenciamento de anexos (imagens)
- Endereços geolocalizados
- Empreendimentos e plantas
- Preços de venda e aluguel
- Corretores principais
- Características e pacotes

#### Usuários (`internal/user/`)
- Autenticação com JWT
- Registro e login
- Perfil de usuário
- RBAC (controle de acesso)

#### Sliders (`internal/sliders/`)
- Gerenciamento de sliders
- Itens de slider
- Suporte a diferentes tipos

### Clean Architecture

Cada domínio segue a estrutura:

```
internal/<dominio>/
├── model.go       # Modelos GORM
├── dto.go         # Data Transfer Objects
├── repository.go  # Camada de acesso a dados
├── service.go     # Lógica de negócio
├── handler.go     # Handlers HTTP (Gin)
└── *_test.go      # Testes
```

**Fluxo:** Handler → Service → Repository → Database

---

## 🤝 Contribuindo

Contribuições são bem-vindas! Veja [CONTRIBUTING.md](CONTRIBUTING.md) para:

- Guias de estilo de código
- Processo de pull request
- Requisitos de testes
- Convenções de commit

### Checklist Antes de Commitar

```bash
make lint-fix    # Corrige problemas automaticamente
make lint        # Verifica qualidade do código
make test        # Executa todos os testes
make swag        # Atualiza documentação (se API mudou)
```

---

## � Tecnologias Utilizadas

- **[Go](https://go.dev/)** — Linguagem de programação
- **[Gin](https://github.com/gin-gonic/gin)** — Framework web rápido
- **[GORM](https://github.com/go-gorm/gorm)** — ORM friendly para desenvolvedores
- **[PostgreSQL](https://www.postgresql.org/)** — Banco de dados relacional
- **[golang-migrate](https://github.com/golang-migrate/migrate)** — Migrações de banco
- **[Viper](https://github.com/spf13/viper)** — Gerenciamento de configuração
- **[golang-jwt](https://github.com/golang-jwt/jwt)** — Implementação JWT
- **[swaggo](https://github.com/swaggo/swag)** — Gerador de documentação Swagger
- **[Air](https://github.com/air-verse/air)** — Hot-reload para desenvolvimento
- **[Docker](https://www.docker.com/)** — Containerização

---

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

---

## 💬 Suporte

Para dúvidas ou problemas, entre em contato com a equipe de desenvolvimento.

---

<div align="center">

**Desenvolvido com ❤️ para TRIIIO**

</div>
