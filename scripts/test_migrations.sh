#!/bin/bash
# Script para testar migrations do zero
# Este script derruba todas as tabelas e reaplica as migrations
#
# ⚠️  IMPORTANTE: Este script precisa ser ATUALIZADO sempre que:
#    - Uma nova migration for criada
#    - Uma migration for renomeada
#    - Uma migration for removida
#    - A ordem das migrations for alterada
#
# 📝 Para atualizar: Adicione/remova/reordene as migrations no array "migrations" abaixo
# 🚀 Para executar: make test-migrations
#
# Última atualização: 30 de Janeiro de 2026
# Migrations registradas: 18

set -e  # Sair em caso de erro

echo "🧪 Iniciando teste de migrations..."
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

DB_USER="triiio_user"
DB_NAME="triiio_backend"
CONTAINER="triiio_db"

# Função para executar SQL
exec_sql() {
    docker exec -i $CONTAINER psql -U $DB_USER -d $DB_NAME -c "$1" 2>&1
}

# Função para executar arquivo SQL
exec_sql_file() {
    docker exec -i $CONTAINER psql -U $DB_USER -d $DB_NAME < "$1" 2>&1
}

echo "📋 1. Listando tabelas atuais..."
exec_sql "\dt" | grep "public |" || echo "Nenhuma tabela encontrada"
echo ""

echo "🗑️  2. Dropando todas as tabelas relacionadas a imóveis..."
exec_sql "DROP TABLE IF EXISTS imoveis CASCADE;"
exec_sql "DROP TABLE IF EXISTS anexos CASCADE;"
exec_sql "DROP TABLE IF EXISTS corretores_principais CASCADE;"
exec_sql "DROP TABLE IF EXISTS pacotes CASCADE;"
exec_sql "DROP TABLE IF EXISTS preco_vendas CASCADE;"
exec_sql "DROP TABLE IF EXISTS preco_alugueis CASCADE;"
exec_sql "DROP TABLE IF EXISTS plantas CASCADE;"
exec_sql "DROP TABLE IF EXISTS torres CASCADE;"
exec_sql "DROP TABLE IF EXISTS empreendimento_caracteristicas CASCADE;"
exec_sql "DROP TABLE IF EXISTS empreendimentos CASCADE;"
exec_sql "DROP TABLE IF EXISTS caracteristicas CASCADE;"
exec_sql "DROP TABLE IF EXISTS enderecos CASCADE;"
exec_sql "DROP TABLE IF EXISTS organizacoes CASCADE;"
echo -e "${GREEN}✅ Tabelas removidas${NC}"
echo ""

echo "🔄 3. Limpando schema_migrations..."
exec_sql "DELETE FROM schema_migrations WHERE version >= '20260113120000';"
echo -e "${GREEN}✅ Schema_migrations limpo${NC}"
echo ""

echo "📦 4. Aplicando migrations na ordem..."
echo ""

# Array de migrations na ordem correta
migrations=(
    "20260113120000_create_enderecos_table"
    "20260113120100_create_plantas_table"
    "20260113120200_create_organizacoes_table"
    "20260113120300_create_pacotes_table"
    "20260113120350_create_anexos_table"
    "20260113120400_create_corretores_principais_table"
    "20260113120500_create_caracteristicas_table"
    "20260113120600_create_preco_vendas_table"
    "20260113120700_create_preco_alugueis_table"
    "20260113120800_create_empreendimentos_table"
    "20260113120810_create_torres_table"
    "20260113120850_add_empreendimento_id_to_plantas_table"
    "20260113120900_create_empreendimento_caracteristicas_table"
    "20260113121000_create_imoveis_table"
    "20260113121200_add_foreign_keys_to_anexos"
    "20260114210019_alter_enderecos_estado_column"
    "20260114211500_add_id_integracao_to_related_tables"
    "20260129120100_alter_organizacoes_table"
)

failed=0
for migration in "${migrations[@]}"; do
    file="migrations/${migration}.up.sql"
    
    if [ -f "$file" ]; then
        echo -n "   Aplicando: $migration ... "
        if exec_sql_file "$file" > /dev/null 2>&1; then
            echo -e "${GREEN}✅${NC}"
            # Registrar na tabela schema_migrations
            version=$(echo $migration | cut -d'_' -f1)
            exec_sql "INSERT INTO schema_migrations (version, dirty) VALUES ('$version', false) ON CONFLICT DO NOTHING;" > /dev/null 2>&1
        else
            echo -e "${RED}❌ FALHOU${NC}"
            echo "Executando novamente para ver erro:"
            exec_sql_file "$file"
            failed=$((failed + 1))
            break
        fi
    else
        echo -e "${YELLOW}⚠️  Arquivo não encontrado: $file${NC}"
    fi
done

echo ""
echo "📊 5. Verificando resultado..."
echo ""

echo "Tabelas criadas:"
exec_sql "\dt" | grep "public |" | awk '{print "   - " $3}'
echo ""

echo "Estrutura da tabela imoveis:"
exec_sql "\d imoveis" | grep -E "corretor_principal_id|organizacao_id" || echo "   ✅ Sem organizacao_id (correto!)"
echo ""

if [ $failed -eq 0 ]; then
    echo -e "${GREEN}✅ SUCESSO! Todas as migrations foram aplicadas corretamente${NC}"
    exit 0
else
    echo -e "${RED}❌ FALHA! $failed migration(s) falharam${NC}"
    exit 1
fi
