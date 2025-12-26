#!/bin/bash

# Script de setup do banco de dados para GastroGo
# Este script ajuda a configurar o banco de dados PostgreSQL

set -e

echo "🚀 Setup do Banco de Dados - GastroGo"
echo ""

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar se PostgreSQL está instalado
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ PostgreSQL não encontrado!${NC}"
    echo ""
    echo "Por favor, instale o PostgreSQL primeiro:"
    echo "  macOS:   brew install postgresql@15"
    echo "  Linux:   sudo apt install postgresql"
    echo "  Docker:  docker run --name gastrogo-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15"
    exit 1
fi

echo -e "${GREEN}✅ PostgreSQL encontrado${NC}"

# Configurações padrão
DB_NAME="${DB_NAME:-gastrogo}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

echo ""
echo "Configurações:"
echo "  Database: $DB_NAME"
echo "  User:     $DB_USER"
echo "  Host:     $DB_HOST"
echo "  Port:     $DB_PORT"
echo ""

# Verificar se o banco já existe
if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
    echo -e "${YELLOW}⚠️  Banco de dados '$DB_NAME' já existe${NC}"
    read -p "Deseja recriar? (isso apagará todos os dados!) [y/N]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Removendo banco existente..."
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "DROP DATABASE IF EXISTS $DB_NAME;"
        echo "Criando novo banco..."
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;"
        echo -e "${GREEN}✅ Banco criado${NC}"
    else
        echo "Usando banco existente..."
    fi
else
    echo "Criando banco de dados..."
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;" || {
        echo -e "${RED}❌ Erro ao criar banco. Verifique suas credenciais.${NC}"
        exit 1
    }
    echo -e "${GREEN}✅ Banco criado${NC}"
fi

# Verificar se migrate está instalado
if ! command -v migrate &> /dev/null; then
    echo ""
    echo -e "${YELLOW}⚠️  golang-migrate não encontrado!${NC}"
    echo ""
    echo "Por favor, instale o golang-migrate:"
    echo "  macOS:   brew install golang-migrate"
    echo "  Linux:   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz"
    exit 1
fi

echo -e "${GREEN}✅ golang-migrate encontrado${NC}"

# Construir URL de conexão
DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

echo ""
echo "Aplicando migrations..."
migrate -path db/migrations -database "$DB_URL" up

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Setup concluído com sucesso!${NC}"
    echo ""
    echo "Próximos passos:"
    echo "  1. Configure as variáveis de ambiente:"
    echo "     export DATABASE_URL=\"$DB_URL\""
    echo "  2. Execute a aplicação:"
    echo "     make run"
else
    echo ""
    echo -e "${RED}❌ Erro ao aplicar migrations${NC}"
    exit 1
fi

