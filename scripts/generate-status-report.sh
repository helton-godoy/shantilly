#!/bin/bash

# Shantilly - Sistema de Relatórios Automatizados
# Gera relatórios completos de status do projeto

set -e  # Parar em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para logging
log() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

info() {
    echo -e "${GREEN}✓${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

# Verificar se está no diretório correto
if [[ ! -f "go.mod" ]]; then
    error "Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Criar diretórios necessários
mkdir -p docs/reports/status
mkdir -p docs/reports/analysis/coverage

# Arquivo de saída
OUTPUT_FILE="docs/reports/status/current-status.md"

log "Iniciando geração do relatório de status..."

# 1. GERAR TIMESTAMP
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S UTC%:z')
DATE_ISO=$(date -Iseconds)

# 2. MÉTRICAS DE COBERTURA
log "Gerando relatório de cobertura..."
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="docs/reports/analysis/coverage/coverage.html"

# Executar testes com cobertura
if go test -coverprofile="$COVERAGE_FILE" -covermode=count ./... 2>/dev/null; then
    # Calcular percentual de cobertura
    COVERAGE_PERCENT=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print substr($3, 1, length($3)-1)}')
    info "Cobertura de testes: $COVERAGE_PERCENT%"

    # Gerar HTML
    go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
    info "Relatório HTML gerado: $COVERAGE_HTML"
else
    warn "Falha ao gerar cobertura - pode haver testes com falha"
    COVERAGE_PERCENT="N/A"
fi

# 3. MÉTRICAS DE QUALIDADE (LINTING E COMPILAÇÃO)
log "Verificando métricas de qualidade..."

# Verificar linting
if golangci-lint run --config .golangci.yml ./... >/dev/null 2>&1; then
    LINT_STATUS="✓ Aprovado"
    LINT_DETAILS="Todas as verificações de linting passaram"
else
    LINT_STATUS="⚠ Problemas"
    LINT_DETAILS="Encontrados problemas de linting (ver detalhes no CI)"
fi

# Verificar compilação
if go build -o /tmp/shantilly-test ./cmd/shantilly 2>/dev/null; then
    BUILD_STATUS="✓ Sucesso"
    BUILD_DETAILS="Projeto compila sem erros"
else
    BUILD_STATUS="✗ Falha"
    BUILD_DETAILS="Projeto possui erros de compilação"
fi

# 4. CONTAGEM DE LINHAS DE CÓDIGO POR CATEGORIA
log "Contando linhas de código por categoria..."

count_lines() {
    local dir="$1"
    local category="$2"
    if [[ -d "$dir" ]]; then
        local count=$(find "$dir" -name "*.go" -exec wc -l {} + | tail -1 | awk '{print $1}' || echo "0")
        echo "$category: $count linhas"
    else
        echo "$category: 0 linhas (diretório não encontrado)"
    fi
}

COMPONENTS_COUNT=$(count_lines "internal/components" "Componentes")
MODELS_COUNT=$(count_lines "internal/models" "Modelos")
CONFIG_COUNT=$(count_lines "internal/config" "Configuração")
ERRORS_COUNT=$(count_lines "internal/errors" "Tratamento de Erros")
STYLES_COUNT=$(count_lines "internal/styles" "Estilos")
CMD_COUNT=$(count_lines "cmd" "Comandos CLI")
TOTAL_COUNT=$(find . -name "*.go" -exec wc -l {} + | tail -1 | awk '{print $1}' || echo "0")

# 5. MÉTRICAS DE TESTES
log "Analisando testes..."

TEST_COUNT=$(find . -name "*_test.go" -exec wc -l {} + | tail -1 | awk '{print $1}' || echo "0")
TEST_FILES=$(find . -name "*_test.go" | wc -l)

# 6. INFORMAÇÕES DO PROJETO
PROJECT_NAME=$(basename "$(pwd)")
GO_VERSION=$(go version | awk '{print $3}')
LAST_COMMIT=$(git log -1 --format="%H" 2>/dev/null || echo "N/A")
LAST_COMMIT_DATE=$(git log -1 --format="%ai" 2>/dev/null || echo "N/A")

# 7. GERAR RELATÓRIO MARKDOWN
log "Gerando relatório final..."

cat > "$OUTPUT_FILE" << EOF
# Relatório de Status - $PROJECT_NAME

**Gerado em:** $TIMESTAMP
**Versão:** $DATE_ISO

## 📊 Visão Geral

- **Projeto:** $PROJECT_NAME
- **Versão Go:** $GO_VERSION
- **Último Commit:** ${LAST_COMMIT:0:8}
- **Data do Commit:** $LAST_COMMIT_DATE

## 🧪 Cobertura de Testes

- **Percentual:** $COVERAGE_PERCENT%
- **Arquivo de Cobertura:** \`$COVERAGE_FILE\`
- **Relatório HTML:** \`$COVERAGE_HTML\`

## 🔍 Qualidade do Código

### Linting
**Status:** $LINT_STATUS
**Detalhes:** $LINT_DETAILS

### Compilação
**Status:** $BUILD_STATUS
**Detalhes:** $BUILD_DETAILS

## 📈 Métricas de Código

### Linhas de Código por Categoria
- $COMPONENTS_COUNT
- $MODELS_COUNT
- $CONFIG_COUNT
- $ERRORS_COUNT
- $STYLES_COUNT
- $CMD_COUNT

**Total de linhas:** $TOTAL_COUNT linhas de código Go

### Testes
- **Arquivos de teste:** $TEST_FILES arquivos
- **Linhas de código de teste:** $TEST_COUNT linhas

## 🏗️ Estrutura do Projeto

\`\`\`
├── internal/
│   ├── components/     ($COMPONENTS_COUNT)
│   ├── models/         ($MODELS_COUNT)
│   ├── config/         ($CONFIG_COUNT)
│   ├── errors/         ($ERRORS_COUNT)
│   └── styles/         ($STYLES_COUNT)
├── cmd/                ($CMD_COUNT)
└── Total: $TOTAL_COUNT linhas
\`\`\`

## 🚀 Como Executar

\`\`\`bash
# Executar todos os testes
make test

# Verificar cobertura
make coverage

# Executar linting
make lint

# Build completo
make build
\`\`\`

---

*Relatório gerado automaticamente por generate-status-report.sh*
EOF

info "Relatório gerado com sucesso: $OUTPUT_FILE"
info "Total de linhas analisadas: $TOTAL_COUNT"
info "Arquivos de teste encontrados: $TEST_FILES"

# 8. GERAR RESUMO EXECUTIVO
SUMMARY_FILE="docs/reports/status/summary.md"

cat > "$SUMMARY_FILE" << EOF
# Resumo Executivo - $(date '+%Y-%m-%d')

**Cobertura:** $COVERAGE_PERCENT% | **Build:** $BUILD_STATUS | **Lint:** $LINT_STATUS

**Total:** $TOTAL_COUNT linhas | **Testes:** $TEST_FILES arquivos

*Atualizado:* $TIMESTAMP
EOF

info "Resumo executivo gerado: $SUMMARY_FILE"

log "Relatório de status concluído com sucesso!"