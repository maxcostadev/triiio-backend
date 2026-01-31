# Referência de Comandos - Projeto Triiio API

> **IMPORTANTE**: Este documento registra comandos testados e que FUNCIONAM no ambiente atual.
> Consultá-lo ANTES de executar comandos evita erros repetidos.

---

## 🎯 PREMISSAS DE TRABALHO - LEIA ANTES DE QUALQUER MODIFICAÇÃO

### ⚠️ REGRA FUNDAMENTAL: VERIFICAR ANTES DE MODIFICAR

**SEMPRE** seguir este fluxo OBRIGATÓRIO antes de fazer qualquer modificação:

#### 1. ANÁLISE COMPLETA (Varredura de Arquivos)

```bash
# Antes de modificar QUALQUER código, SEMPRE executar:

# 1.1 Buscar TODAS as ocorrências da entidade/campo em TODO o projeto
grep -r "NomeDaEntidade" internal/ --include="*.go"
grep -r "nome_do_campo" internal/ --include="*.go"
grep -r "nome_do_campo" migrations/ --include="*.sql"

# 1.2 Buscar variações (CamelCase, snake_case, etc)
grep -r -i "organizacao" internal/ --include="*.go"
grep -r "OrganizacaoID\|organizacao_id\|Organizacao" internal/

# 1.3 Identificar TODOS os arquivos afetados
grep -l "OrganizacaoID" internal/**/*.go

# 1.4 Verificar dependências entre pacotes
grep -r "import.*imoveis" internal/
```

#### 2. MAPEAMENTO DE IMPACTO

**NUNCA** modificar apenas um arquivo. **SEMPRE** mapear:

- ✅ **Model** (`model.go`) - Definição da struct
- ✅ **DTO** (`dto.go`) - Request/Response structs
- ✅ **Repository** (`repository.go`) - Queries e mapeamentos
- ✅ **Service** (`service.go`) - Lógica de negócio
- ✅ **Handler** (`handler.go`) - Endpoints HTTP
- ✅ **Migrations** (`*.sql`) - Estrutura do banco
- ✅ **Tests** (`*_test.go`) - Testes unitários

#### 3. CHECKLIST DE VERIFICAÇÃO

Antes de confirmar QUALQUER modificação, verificar:

```bash
# ✅ 3.1 Compilação
go build ./internal/pacote/...

# ✅ 3.2 Erros no código
go vet ./internal/pacote/...

# ✅ 3.3 Buscar TODOs ou comentários que mencionam o campo
grep -r "TODO.*organizacao\|FIXME.*organizacao" internal/

# ✅ 3.4 Verificar se há outros serviços que usam a entidade
grep -r "Imovel\|ImovelResponse" internal/ --include="*.go" | grep -v "internal/imoveis"
```

#### 4. EXEMPLO REAL: Caso OrganizacaoID → CorretorPrincipalID

**❌ ERRO COMETIDO:**
```
Solicitação: "Corrigir erro em repository.go"
Ação: Modifiquei APENAS repository.go
Resultado: service.go continuou com OrganizacaoID, causando 10+ erros
```

**✅ FLUXO CORRETO:**
```bash
# 1. Buscar TODAS as ocorrências
grep -r "OrganizacaoID\|organizacao_id" internal/imoveis/

# Resultado encontrado:
# - internal/imoveis/model.go (struct Imovel)
# - internal/imoveis/dto.go (CreateImovelRequest, UpdateImovelRequest)
# - internal/imoveis/repository.go (mapeamento, queries)
# - internal/imoveis/service.go (CreateImovel, UpdateImovel, etc)
# - migrations/20260113121000_create_imoveis_table.up.sql

# 2. Modificar TODOS de uma vez
multi_replace_string_in_file com TODOS os arquivos

# 3. Verificar se compilou
go build ./internal/imoveis/...

# 4. Verificar erros restantes
get_errors filePaths=["/internal/imoveis"]
```

### 🔍 FERRAMENTAS DE ANÁLISE

#### Busca Semântica (Preferred)
```python
semantic_search(query="organizacao property relationship")
```

#### Busca por Padrão
```bash
grep_search(
    query="OrganizacaoID|organizacao_id",
    isRegexp=true,
    includePattern="internal/**/*.go"
)
```

#### Lista de Uso de Símbolo
```python
list_code_usages(
    symbolName="OrganizacaoID",
    filePaths=["/internal/imoveis/model.go"]
)
```

### 📋 TEMPLATE DE MODIFICAÇÃO

Ao receber solicitação de modificação, SEMPRE seguir:

```markdown
## Análise de Impacto: [Nome da Modificação]

### 1. Busca Completa
- [ ] Grep em internal/**/*.go
- [ ] Grep em migrations/*.sql
- [ ] Grep em testes
- [ ] Verificar imports

### 2. Arquivos Identificados
- [ ] model.go - linha X
- [ ] dto.go - linhas Y, Z
- [ ] repository.go - linhas A, B, C
- [ ] service.go - linhas D, E, F
- [ ] migrations/*.sql - arquivo X

### 3. Modificações Planejadas
- [ ] Arquivo 1: mudança X
- [ ] Arquivo 2: mudança Y
- [ ] Arquivo 3: mudança Z

### 4. Validação
- [ ] go build ./...
- [ ] get_errors
- [ ] Testes relevantes
```

### ⛔ ANTI-PADRÕES - NUNCA FAZER

1. ❌ Modificar apenas um arquivo quando solicitado "corrigir erro em X"
2. ❌ Assumir que grep inicial encontrou tudo
3. ❌ Ignorar variações de nomenclatura (CamelCase, snake_case, kebab-case)
4. ❌ Esquecer de verificar migrations
5. ❌ Modificar código sem compilar depois
6. ❌ Fazer mudanças sem buscar usages do símbolo
7. ❌ Confiar apenas em mensagens de erro do compilador

### ✅ BOAS PRÁTICAS OBRIGATÓRIAS

1. ✅ Sempre fazer varredura completa ANTES de modificar
2. ✅ Usar `multi_replace_string_in_file` para mudanças em múltiplos arquivos
3. ✅ Compilar após CADA mudança
4. ✅ Verificar erros com `get_errors` após modificações
5. ✅ Buscar variações de nomenclatura
6. ✅ Verificar migrations relacionadas
7. ✅ Documentar arquivos modificados
8. ✅ Validar que TODOS os erros relacionados foram corrigidos

---

## 🐘 PostgreSQL / Banco de Dados

### Credenciais (de .env)
```bash
DATABASE_HOST=triiio_db
DATABASE_PORT=5432
DATABASE_USER=triiio_user              # ⚠️ NÃO é "postgres"
DATABASE_PASSWORD='Soeusei2w@123&xJ'
DATABASE_NAME=triiio_backend           # ⚠️ NÃO é "grab"
DATABASE_SSLMODE=disable
```

### ✅ Comandos que FUNCIONAM

```bash
# Executar SQL via stdin
docker exec -i triiio_db psql -U triiio_user -d triiio_backend < script.sql

# Executar comando SQL direto
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "SELECT * FROM imoveis LIMIT 5;"

# Listar tabelas
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\dt"

# Descrever estrutura de tabela
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\d imoveis"

# Verificar se container está rodando
docker-compose ps | grep triiio_db
```

### ❌ Comandos que NÃO FUNCIONAM

```bash
# ERRADO: Usuário postgres não existe
docker exec -i triiio_db psql -U postgres -d triiio

# ERRADO: Database grab não existe
docker exec -i triiio_db psql -U triiio_user -d grab

# ERRADO: Makefile não tem este comando
make exec-db
```

## 🔄 Migrations

### ⚠️ PROBLEMA CONHECIDO
```bash
# make migrate-* não funciona sem Docker ou com variáveis de ambiente
make migrate-status  # ❌ ERRO: JWT_SECRET required
make migrate-up      # ❌ ERRO: JWT_SECRET required
```

### ✅ Comandos Funcionais

```bash
# Testar migrations do zero (⚠️  APENAS DESENVOLVIMENTO!)
make test-migrations

# Status das migrations (via tabela do banco)
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 10;"

# Listar migrations pendentes
ls -1 migrations/*.up.sql | sort
```

### ✅ Solução Alternativa - Via Docker

```bash
# Status das migrations
docker exec -i triiio_api go run cmd/migrate/main.go status

# Aplicar migrations
docker exec -i triiio_api go run cmd/migrate/main.go up

# Reverter 1 migration
docker exec -i triiio_api go run cmd/migrate/main.go down

# Criar nova migration
docker exec -i triiio_api go run cmd/migrate/main.go create nome_da_migration
```

### ✅ Solução Alternativa - Via golang-migrate CLI

```bash
# Se golang-migrate estiver instalado no host
migrate -path migrations -database "postgresql://triiio_user:Soeusei2w@123&xJ@localhost:5432/triiio_backend?sslmode=disable" status
migrate -path migrations -database "postgresql://triiio_user:Soeusei2w@123&xJ@localhost:5432/triiio_backend?sslmode=disable" up
migrate -path migrations -database "postgresql://triiio_user:Soeusei2w@123&xJ@localhost:5432/triiio_backend?sslmode=disable" down
```

### ✅ Solução Alternativa - Manualmente com psql

```bash
# Executar migrations manualmente (desenvolvimento/teste)
cd /home/maxcosta/triiio/triiio-api

# Aplicar migration específica
docker exec -i triiio_db psql -U triiio_user -d triiio_backend < migrations/20260113121000_create_imoveis_table.up.sql

# Reverter migration específica
docker exec -i triiio_db psql -U triiio_user -d triiio_backend < migrations/20260113121000_create_imoveis_table.down.sql
```

## 🐳 Docker

### Verificar containers rodando
```bash
docker-compose ps
# ou
docker ps | grep triiio
```

### Logs
```bash
# Ver logs do app
docker logs triiio_api -f

# Ver logs do banco
docker logs triiio_db -f
```

### Reiniciar serviços
```bash
# Reiniciar tudo
docker-compose restart

# Reiniciar apenas o app
docker-compose restart app

# Reiniciar apenas o banco
docker-compose restart db
```

## 🔨 Build e Testes

### ✅ Go Build
```bash
# Compilar pacote específico
cd /home/maxcosta/triiio/triiio-api
go build ./internal/imoveis/...

# Compilar tudo
go build ./...

# Executar testes
go test ./internal/imoveis/... -v
```

### ❌ Erros Comuns

```bash
# ERRADO: Executar fora do diretório do projeto
go build ./internal/imoveis/...  # ❌ Se não estiver em /home/maxcosta/triiio/triiio-api

# CORRETO: Sempre fazer cd primeiro
cd /home/maxcosta/triiio/triiio-api && go build ./internal/imoveis/...
```

## 📝 Grep / Search

### ✅ Buscar em arquivos
```bash
# Buscar padrão em migrations
grep -r "organizacao_id" migrations/

# Buscar em arquivos .go
grep -r "OrganizacaoID" internal/

# Ver número da linha
grep -n "CorretorPrincipal" internal/imoveis/service.go

# Ver contexto (linhas antes/depois)
grep -A 5 -B 5 "pattern" file.go
```

## 🔍 Verificações Úteis

### Verificar estrutura do banco
```bash
# Listar todas as tabelas
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\dt" | grep -E "imovel|corretor|anexo|pacote|preco|planta|torre|empreend|caract|endereco|organiza"

# Ver colunas de uma tabela específica
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\d imoveis" | grep -E "organizacao|corretor"

# Contar registros
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "SELECT COUNT(*) FROM imoveis;"
```

### Verificar migrations aplicadas
```bash
# Via tabela schema_migrations (se existir)
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "SELECT * FROM schema_migrations ORDER BY version;"
```

### Verificar ordem de migrations
```bash
# Listar migrations em ordem
ls -1 migrations/*.up.sql | sort

# Filtrar por padrão
ls -1 migrations/202601*.up.sql | sort
```

### ⚠️ IMPORTANTE: Manter script de teste atualizado
```bash
# O arquivo scripts/test_migrations.sh contém um array hardcoded de migrations
# SEMPRE atualizar quando:
# - Criar nova migration
# - Renomear migration
# - Remover migration
# - Alterar ordem de migrations

# Localização: scripts/test_migrations.sh
# Executar com: make test-migrations
```

## 🚨 Checklist ANTES de Executar Comandos

1. **Banco de dados**: Sempre usar `triiio_user` e `triiio_backend`
2. **Container**: Verificar se está rodando com `docker-compose ps`
3. **Diretório**: Estar em `/home/maxcosta/triiio/triiio-api`
4. **Migrations**: Usar docker exec ou golang-migrate CLI, NÃO make commands
5. **Variáveis de ambiente**: Estão em `.env`, não tentar exportar manualmente

## 📊 Padrões de Nomenclatura

### Migrations
```
YYYYMMDDHHMMSS_verb_noun_table.up.sql
YYYYMMDDHHMMSS_verb_noun_table.down.sql

Exemplos:
20260113121000_create_imoveis_table.up.sql
20260113121200_add_foreign_keys_to_anexos.up.sql
20260113120400_create_corretores_principais_table.up.sql
```

### Ordem de criação de tabelas (dependências)
```
1. Tabelas independentes (enderecos, organizacoes, etc)
2. Tabelas com foreign keys simples
3. Tabelas que referenciam múltiplas outras (anexos sem FKs)
4. Tabelas finais (imoveis)
5. Adicionar foreign keys complexas (anexos FKs)
```

## 💡 Lições Aprendidas

1. **NÃO assumir nomes padrão** - Sempre verificar .env primeiro
2. **NÃO usar make para migrations** - Usar docker exec ou CLI direta
3. **Verificar container primeiro** - Evita "connection refused"
4. **Usar sed para edições** - Quando replace_string_in_file falha por tabs/espaços
5. **Verificar ordem de migrations** - Dependências devem vir antes
6. **Truncate com CASCADE** - Para lidar com foreign keys
7. **IF EXISTS em migrations** - Para idempotência

## 🔄 Fluxo de Teste de Migrations

```bash
# 1. Verificar estado atual
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\dt"

# 2. Fazer backup (opcional)
docker exec triiio_db pg_dump -U triiio_user triiio_backend > backup_$(date +%Y%m%d_%H%M%S).sql

# 3. Resetar banco (desenvolvimento apenas!)
docker exec -i triiio_db psql -U triiio_user -d triiio_backend < scripts/truncate_imoveis_tables.sql

# 4. Aplicar migrations manualmente uma por uma
for file in migrations/*.up.sql; do
    echo "Applying $file"
    docker exec -i triiio_db psql -U triiio_user -d triiio_backend < "$file"
    if [ $? -ne 0 ]; then
        echo "❌ Failed: $file"
        break
    fi
    echo "✅ Success: $file"
done

# 5. Verificar resultado
docker exec -i triiio_db psql -U triiio_user -d triiio_backend -c "\d imoveis"
```

---

**Última atualização**: 30 de Janeiro de 2026
**Mantido por**: AI Assistant (aprendizado contínuo)
