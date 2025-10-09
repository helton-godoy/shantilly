#!/bin/bash

# Shantilly - Sistema de Captura de Relatórios Dinâmicos
# Captura automaticamente relatórios gerados durante interações
# Mantém estrutura organizada e indexação automática

set -e  # Parar em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Funções de logging
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

debug() {
    echo -e "${CYAN}🔍${NC} $1"
}

# Verificar se está no diretório correto
if [[ ! -f "go.mod" ]]; then
    error "Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Configuração
REPORTS_DIR="docs/reports/dynamic"
CACHE_DIR="$REPORTS_DIR/cache"
INDEX_FILE="$REPORTS_DIR/dynamic-reports-index.md"
MAX_CACHE_SIZE=50

# Criar diretórios necessários
mkdir -p "$REPORTS_DIR"/{progress,analysis,status}
mkdir -p "$CACHE_DIR"

# Função para gerar ID único baseado no timestamp e hash do conteúdo
generate_id() {
    local content="$1"
    local timestamp=$(date -Iseconds)
    echo "report_$(echo -n "$timestamp" | sha256sum | cut -d' ' -f1 | cut -c1-16)"
}

# Função para detectar categoria baseada no conteúdo
detect_category() {
    local content="$1"

    if echo "$content" | grep -qi "progress\|desenvolvimento\|fase\|entrega"; then
        echo "progress"
    elif echo "$content" | grep -qi "analysis\|análise\|métricas\|cobertura"; then
        echo "analysis"
    elif echo "$content" | grep -qi "status\|estado\|situação"; then
        echo "status"
    else
        echo "general"
    fi
}

# Função para extrair contexto baseado no conteúdo
extract_context() {
    local content="$1"

    # Tentar identificar contexto baseado em palavras-chave
    if echo "$content" | grep -qi "macro tarefas\|fase\|entrega"; then
        echo "macro-tarefas"
    elif echo "$content" | grep -qi "documentação\|estrutura"; then
        echo "documentacao"
    elif echo "$content" | grep -qi "implementação\|sistema"; then
        echo "implementacao"
    elif echo "$content" | grep -qi "localização\|sistema implementado"; then
        echo "localizacao"
    elif echo "$content" | grep -qi "teste\|cobertura"; then
        echo "testes"
    elif echo "$content" | grep -qi "segurança"; then
        echo "seguranca"
    else
        echo "geral"
    fi
}

# Função para salvar relatório
save_report() {
    local content="$1"
    local title="${2:-Relatório Dinâmico}"
    local category=$(detect_category "$content")
    local context=$(extract_context "$content")
    local timestamp=$(date -Iseconds)
    local report_id=$(generate_id "$content")

    local filename="$report_id.md"
    local filepath="$REPORTS_DIR/$category/$filename"

    # Criar relatório com metadados
    cat > "$filepath" << EOF
# $title

**ID:** $report_id
**Categoria:** $category
**Contexto:** $context
**Timestamp:** $timestamp
**Gerado automaticamente por:** capture-dynamic-reports.sh

---

$content

---

*Este relatório foi capturado automaticamente em $(date '+%Y-%m-%d %H:%M:%S %Z')*
EOF

    info "Relatório salvo: $filepath"
    echo "$report_id|$timestamp|$category|$context|$title"
}

# Função para atualizar cache
update_cache() {
    local report_data="$1"
    local cache_file="$CACHE_DIR/latest-reports.cache"

    # Adicionar ao cache
    echo "$report_data" >> "$cache_file"

    # Manter apenas os últimos N relatórios
    local temp_file=$(mktemp)
    tail -n $MAX_CACHE_SIZE "$cache_file" > "$temp_file" 2>/dev/null || true
    mv "$temp_file" "$cache_file" 2>/dev/null || true

    debug "Cache atualizado"
}

# Função para atualizar índice
update_index() {
    local index_temp=$(mktemp)

    cat > "$index_temp" << EOF
# 📊 Índice de Relatórios Dinâmicos

**Última atualização:** $(date '+%Y-%m-%d %H:%M:%S %Z')
**Total de relatórios:** $(find "$REPORTS_DIR" -name "*.md" -not -name "dynamic-reports-index.md" | wc -l)

## 📈 Estatísticas por Categoria

EOF

    # Estatísticas por categoria
    for category in progress analysis status general; do
        local count=$(find "$REPORTS_DIR/$category" -name "*.md" 2>/dev/null | wc -l)
        if [[ $count -gt 0 ]]; then
            echo "**$category:** $count relatórios" >> "$index_temp"
        fi
    done

    echo "" >> "$index_temp"
    echo "## 📋 Lista de Relatórios" >> "$index_temp"
    echo "" >> "$index_temp"
    echo "| ID | Data | Categoria | Contexto | Título |" >> "$index_temp"
    echo "|:---|:-----|:----------|:---------|:-------|" >> "$index_temp"

    # Listar relatórios ordenados por data (mais recentes primeiro)
    find "$REPORTS_DIR" -name "*.md" -not -name "dynamic-reports-index.md" -exec grep -l "^**ID:**" {} \; | \
    while IFS= read -r file; do
        local id=$(grep "^**ID:**" "$file" | head -1 | sed 's/\*\*ID:\*\* //')
        local timestamp=$(grep "^**Timestamp:**" "$file" | head -1 | sed 's/\*\*Timestamp:\*\* //')
        local category=$(grep "^**Categoria:**" "$file" | head -1 | sed 's/\*\*Categoria:\*\* //')
        local context=$(grep "^**Contexto:**" "$file" | head -1 | sed 's/\*\*Contexto:\*\* //')
        local title=$(grep "^# " "$file" | head -1 | sed 's/^# //')

        echo "| [$id]($category/$id.md) | $timestamp | $category | $context | $title |" >> "$index_temp"
    done | sort -r

    echo "" >> "$index_temp"
    echo "## 🔍 Últimos Relatórios Capturados" >> "$index_temp"
    echo "" >> "$index_temp"

    # Últimos 10 relatórios do cache
    if [[ -f "$CACHE_DIR/latest-reports.cache" ]]; then
        tail -10 "$CACHE_DIR/latest-reports.cache" | while IFS='|' read -r id timestamp category context title; do
            echo "- **$timestamp** - [$title]($category/$id.md) ($context)" >> "$index_temp"
        done
    fi

    echo "" >> "$index_temp"
    echo "---" >> "$index_temp"
    echo "" >> "$index_temp"
    echo "*Índice gerado automaticamente por capture-dynamic-reports.sh*" >> "$index_temp"

    mv "$index_temp" "$INDEX_FILE"
    info "Índice atualizado: $INDEX_FILE"
}

# Função para capturar relatório interativo
capture_interactive() {
    local title="$1"

    echo -e "${BLUE}📝 Capturando relatório dinâmico...${NC}"
    echo "Título: $title"
    echo ""
    echo -e "${CYAN}Digite o conteúdo do relatório (Ctrl+D para finalizar):${NC}"

    local content=""
    while IFS= read -r line; do
        content="$content$line"$'\n'
    done

    if [[ -n "$content" ]]; then
        local report_data=$(save_report "$content" "$title")
        update_cache "$report_data"
        update_index
        info "Relatório capturado com sucesso!"
    else
        warn "Conteúdo vazio, relatório não foi salvo"
    fi
}

# Função para capturar relatório programático
capture_programmatic() {
    local content="$1"
    local title="${2:-Relatório Dinâmico}"

    if [[ -n "$content" ]]; then
        local report_data=$(save_report "$content" "$title")
        update_cache "$report_data"
        update_index
        echo "$report_data"
    else
        error "Conteúdo vazio, relatório não foi salvo"
        return 1
    fi
}

# Função para listar relatórios por categoria
list_reports() {
    local category="$1"

    if [[ -n "$category" ]]; then
        find "$REPORTS_DIR/$category" -name "*.md" -exec basename {} \; | sort
    else
        find "$REPORTS_DIR" -name "*.md" -not -name "dynamic-reports-index.md" -exec basename {} \; | sort
    fi
}

# Função para buscar relatórios por contexto
search_reports() {
    local query="$1"

    find "$REPORTS_DIR" -name "*.md" -not -name "dynamic-reports-index.md" -exec grep -l "$query" {} \;
}

# Função para limpeza de cache antigo
cleanup_cache() {
    local cache_file="$CACHE_DIR/latest-reports.cache"
    local max_age_days=${1:-30}

    if [[ -f "$cache_file" ]]; then
        local temp_file=$(mktemp)
        local cutoff_date=$(date -d "$max_age_days days ago" +%s)

        while IFS='|' read -r id timestamp category context title; do
            local report_date=$(date -d "$(echo $timestamp | cut -d' ' -f1)" +%s 2>/dev/null || echo 0)
            if [[ $report_date -gt $cutoff_date ]]; then
                echo "$id|$timestamp|$category|$context|$title" >> "$temp_file"
            fi
        done < "$cache_file"

        mv "$temp_file" "$cache_file"
        debug "Cache limpo (relatórios mais antigos que $max_age_days dias removidos)"
    fi
}

# Função principal
main() {
    local action="$1"
    shift

    case "$action" in
        "capture")
            if [[ -t 0 ]]; then
                capture_interactive "$@"
            else
                # Modo pipe - capturar do stdin
                local content
                content=$(cat)
                capture_programmatic "$content" "$@"
            fi
            ;;
        "list")
            list_reports "$@"
            ;;
        "search")
            search_reports "$@"
            ;;
        "cleanup")
            cleanup_cache "$@"
            ;;
        "index")
            update_index
            ;;
        *)
            echo "Uso: $0 {capture|list|search|cleanup|index} [parâmetros]"
            echo ""
            echo "Exemplos:"
            echo "  echo 'Conteúdo do relatório' | $0 capture 'Título do Relatório'"
            echo "  $0 list [categoria]"
            echo "  $0 search 'palavra-chave'"
            echo "  $0 cleanup 30"
            echo "  $0 index"
            exit 1
            ;;
    esac
}

# Executar função principal
main "$@"