#!/bin/bash

# Shantilly - Sistema de Backup de Relatórios Dinâmicos
# Realiza backup automático e restauração de relatórios dinâmicos
# Mantém logs detalhados de todas as operações realizadas

set -e  # Parar em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
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

success() {
    echo -e "${PURPLE}✅${NC} $1"
}

# Verificar se está no diretório correto
if [[ ! -f "go.mod" ]]; then
    error "Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Configuração
REPORTS_DIR="docs/reports/dynamic"
BACKUP_DIR="docs/reports/backup"
LOG_DIR="docs/reports/logs"
BACKUP_NAME="dynamic-reports-backup-$(date +%Y%m%d-%H%M%S)"
MAX_BACKUPS=10

# Criar diretórios necessários
mkdir -p "$BACKUP_DIR"
mkdir -p "$LOG_DIR"

# Arquivo de log
LOG_FILE="$LOG_DIR/backup-$(date +%Y%m%d-%H%M%S).log"

# Função para logging detalhado
log_action() {
    local level="$1"
    local message="$2"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $message" | tee -a "$LOG_FILE"
}

# Função para verificar integridade dos arquivos
verify_integrity() {
    local dir="$1"
    local errors=0

    log_action "INFO" "Verificando integridade de $dir"

    # Verificar se todos os arquivos .md têm o formato correto
    while IFS= read -r file; do
        if [[ -f "$file" ]]; then
            # Verificar se tem cabeçalho obrigatório
            if ! grep -q "^# " "$file" || ! grep -q "^\\*\\*ID:\\*\\*" "$file"; then
                log_action "WARN" "Arquivo sem formato correto: $file"
                ((errors++))
            fi

            # Verificar se arquivo não está vazio
            if [[ ! -s "$file" ]]; then
                log_action "WARN" "Arquivo vazio: $file"
                ((errors++))
            fi
        fi
    done < <(find "$dir" -name "*.md" 2>/dev/null)

    if [[ $errors -eq 0 ]]; then
        log_action "INFO" "Verificação de integridade OK ($dir)"
        return 0
    else
        log_action "ERROR" "Encontrados $errors erros de integridade em $dir"
        return 1
    fi
}

# Função para calcular tamanho do backup
calculate_size() {
    local dir="$1"
    if [[ -d "$dir" ]]; then
        du -sh "$dir" 2>/dev/null | cut -f1 || echo "0B"
    else
        echo "0B"
    fi
}

# Função para criar backup completo
create_backup() {
    local backup_path="$BACKUP_DIR/$BACKUP_NAME"
    local temp_dir=$(mktemp -d)

    log_action "INFO" "Iniciando backup: $BACKUP_NAME"

    # Criar estrutura do backup
    mkdir -p "$temp_dir/$BACKUP_NAME"

    # Copiar relatórios dinâmicos
    if [[ -d "$REPORTS_DIR" ]]; then
        cp -r "$REPORTS_DIR" "$temp_dir/$BACKUP_NAME/" 2>/dev/null || {
            log_action "ERROR" "Falha ao copiar diretório de relatórios"
            rm -rf "$temp_dir"
            return 1
        }
    fi

    # Criar arquivo de metadados do backup
    cat > "$temp_dir/$BACKUP_NAME/backup-metadata.json" << EOF
{
    "backup_name": "$BACKUP_NAME",
    "created_at": "$(date -Iseconds)",
    "backup_type": "full",
    "source_directory": "$REPORTS_DIR",
    "total_files": "$(find $temp_dir/$BACKUP_NAME -type f | wc -l)",
    "total_size": "$(calculate_size "$temp_dir/$BACKUP_NAME")",
    "project": "$(basename $(pwd))",
    "backup_tool": "backup-dynamic-reports.sh"
}
EOF

    # Criar checksum para verificação de integridade
    if command -v sha256sum &> /dev/null; then
        find "$temp_dir/$BACKUP_NAME" -type f -exec sha256sum {} \; > "$temp_dir/$BACKUP_NAME/checksums.sha256"
        log_action "INFO" "Checksums calculados"
    fi

    # Compactar backup
    local tar_file="$backup_path.tar.gz"
    tar -czf "$tar_file" -C "$temp_dir" "$BACKUP_NAME" 2>/dev/null

    if [[ $? -eq 0 ]]; then
        local backup_size=$(calculate_size "$tar_file")
        log_action "INFO" "Backup criado com sucesso: $tar_file ($backup_size)"
        info "Backup criado: $tar_file"
        success "Tamanho do backup: $backup_size"
    else
        log_action "ERROR" "Falha ao criar arquivo de backup"
        rm -rf "$temp_dir"
        return 1
    fi

    # Limpar arquivos temporários
    rm -rf "$temp_dir"

    # Manter apenas os últimos N backups
    cleanup_old_backups

    return 0
}

# Função para restaurar backup
restore_backup() {
    local backup_file="$1"
    local restore_path="${2:-$REPORTS_DIR}"

    if [[ ! -f "$backup_file" ]]; then
        error "Arquivo de backup não encontrado: $backup_file"
        log_action "ERROR" "Arquivo de backup não encontrado: $backup_file"
        return 1
    fi

    log_action "INFO" "Iniciando restauração de $backup_file para $restore_path"

    # Verificar se arquivo de backup é válido
    if ! tar -tzf "$backup_file" >/dev/null 2>&1; then
        error "Arquivo de backup inválido ou corrompido"
        log_action "ERROR" "Arquivo de backup inválido: $backup_file"
        return 1
    fi

    # Criar backup do estado atual antes de restaurar (se existir)
    if [[ -d "$restore_path" ]]; then
        local pre_restore_backup="$BACKUP_DIR/pre-restore-$(date +%Y%m%d-%H%M%S).tar.gz"
        tar -czf "$pre_restore_backup" -C "$(dirname "$restore_path")" "$(basename "$restore_path")" 2>/dev/null || true
        log_action "INFO" "Backup pré-restauração criado: $pre_restore_backup"
    fi

    # Extrair backup
    if tar -xzf "$backup_file" -C "$(dirname "$restore_path")" 2>/dev/null; then
        log_action "INFO" "Restauração concluída com sucesso"
        success "Backup restaurado: $(basename "$backup_file")"
        info "Destino: $restore_path"
    else
        error "Falha na restauração do backup"
        log_action "ERROR" "Falha na extração do backup: $backup_file"
        return 1
    fi
}

# Função para listar backups disponíveis
list_backups() {
    log_action "INFO" "Listando backups disponíveis"

    if [[ ! -d "$BACKUP_DIR" ]]; then
        warn "Diretório de backups não encontrado"
        return 1
    fi

    echo -e "${BLUE}📦 Backups Disponíveis:${NC}"
    echo ""

    # Cabeçalho da tabela
    printf "%-30s %-15s %-10s %-20s\n" "NOME DO BACKUP" "DATA" "TAMANHO" "STATUS"
    echo "--------------------------------------------------------------------------------"

    # Listar backups ordenados por data (mais recentes primeiro)
    find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | sort -r | while read -r backup_file; do
        local filename=$(basename "$backup_file")
        local size=$(calculate_size "$backup_file")
        local date=$(stat -c %y "$backup_file" 2>/dev/null | cut -d'.' -f1 || echo "N/A")

        # Verificar integridade
        if tar -tzf "$backup_file" >/dev/null 2>&1; then
            local status="${GREEN}✓ OK${NC}"
        else
            local status="${RED}✗ Corrompido${NC}"
        fi

        printf "%-30s %-15s %-10s %-20s\n" "$filename" "$date" "$size" "$status"
    done

    local total_backups=$(find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | wc -l)
    echo ""
    info "Total de backups: $total_backups"
}

# Função para limpar backups antigos
cleanup_old_backups() {
    local max_backups=${1:-$MAX_BACKUPS}

    local backup_count=$(find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | wc -l)

    if [[ $backup_count -gt $max_backups ]]; then
        local to_remove=$((backup_count - max_backups))

        log_action "INFO" "Removendo $to_remove backups antigos (mantendo $max_backups)"

        # Remover backups mais antigos
        find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | sort | head -n $to_remove | while read -r old_backup; do
            local filename=$(basename "$old_backup")
            rm -f "$old_backup"
            log_action "INFO" "Backup removido: $filename"
        done

        info "Limpeza concluída"
    else
        debug "Nenhuma limpeza necessária ($backup_count/$max_backups backups)"
    fi
}

# Função para verificar backup específico
verify_backup() {
    local backup_file="$1"

    if [[ ! -f "$backup_file" ]]; then
        error "Arquivo de backup não encontrado: $backup_file"
        return 1
    fi

    log_action "INFO" "Verificando backup: $backup_file"

    # Teste básico de integridade
    if tar -tzf "$backup_file" >/dev/null 2>&1; then
        local temp_dir=$(mktemp -d)
        tar -xzf "$backup_file" -C "$temp_dir" >/dev/null 2>&1

        # Verificar estrutura esperada
        if [[ -f "$temp_dir/$(basename "$backup_file" .tar.gz)/backup-metadata.json" ]]; then
            success "Backup válido e íntegro"
            local metadata=$(cat "$temp_dir/$(basename "$backup_file" .tar.gz)/backup-metadata.json")
            echo -e "${BLUE}Metadados do backup:${NC}"
            echo "$metadata" | jq . 2>/dev/null || echo "$metadata"
        else
            warn "Backup sem metadados - pode estar corrompido"
            return 1
        fi

        rm -rf "$temp_dir"
        return 0
    else
        error "Backup corrompido ou inválido"
        return 1
    fi
}

# Função para estatísticas de backup
backup_stats() {
    log_action "INFO" "Gerando estatísticas de backup"

    local total_backups=$(find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | wc -l)
    local total_size=$(du -sh "$BACKUP_DIR" 2>/dev/null | cut -f1 || echo "0B")
    local oldest_backup=$(find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | sort | head -1 | xargs basename 2>/dev/null || echo "Nenhum")
    local newest_backup=$(find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | sort -r | head -1 | xargs basename 2>/dev/null || echo "Nenhum")

    echo -e "${BLUE}📊 Estatísticas de Backup:${NC}"
    echo ""
    echo "Total de backups: $total_backups"
    echo "Espaço ocupado: $total_size"
    echo "Backup mais antigo: $oldest_backup"
    echo "Backup mais recente: $newest_backup"
    echo "Localização: $BACKUP_DIR"
    echo ""

    # Verificar saúde dos backups
    local corrupted=0
    while IFS= read -r backup_file; do
        if ! tar -tzf "$backup_file" >/dev/null 2>&1; then
            ((corrupted++))
        fi
    done < <(find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null)

    if [[ $corrupted -eq 0 ]]; then
        success "Todos os backups estão íntegros"
    else
        warn "$corrupted backup(s) podem estar corrompidos"
    fi
}

# Função para restaurar do backup mais recente
restore_latest() {
    local latest_backup=$(find "$BACKUP_DIR" -name "*.tar.gz" -type f 2>/dev/null | sort -r | head -1)

    if [[ -z "$latest_backup" ]]; then
        error "Nenhum backup encontrado para restauração"
        return 1
    fi

    warn "Restaurando do backup mais recente: $(basename "$latest_backup")"
    restore_backup "$latest_backup" "$REPORTS_DIR"
}

# Função principal
main() {
    local action="$1"
    shift

    case "$action" in
        "create")
            log_action "INFO" "=== INICIANDO CRIAÇÃO DE BACKUP ==="
            verify_integrity "$REPORTS_DIR" || {
                error "Problemas de integridade detectados. Abortando backup."
                exit 1
            }
            create_backup
            log_action "INFO" "=== BACKUP CONCLUÍDO ==="
            ;;
        "restore")
            if [[ $# -eq 1 ]]; then
                restore_backup "$1"
            else
                error "Uso: $0 restore <arquivo_de_backup>"
                exit 1
            fi
            ;;
        "restore-latest")
            restore_latest
            ;;
        "list")
            list_backups
            ;;
        "verify")
            if [[ $# -eq 1 ]]; then
                verify_backup "$1"
            else
                error "Uso: $0 verify <arquivo_de_backup>"
                exit 1
            fi
            ;;
        "cleanup")
            cleanup_old_backups "$@"
            ;;
        "stats")
            backup_stats
            ;;
        *)
            echo "Sistema de Backup de Relatórios Dinâmicos"
            echo ""
            echo "Uso: $0 {create|restore|restore-latest|list|verify|cleanup|stats}"
            echo ""
            echo "Comandos:"
            echo "  create           Criar backup completo dos relatórios dinâmicos"
            echo "  restore FILE     Restaurar backup específico"
            echo "  restore-latest   Restaurar do backup mais recente"
            echo "  list             Listar todos os backups disponíveis"
            echo "  verify FILE      Verificar integridade de backup específico"
            echo "  cleanup [N]      Remover backups antigos (padrão: manter últimos 10)"
            echo "  stats            Mostrar estatísticas de backup"
            echo ""
            echo "Exemplos:"
            echo "  $0 create"
            echo "  $0 restore dynamic-reports-backup-20231009-120000.tar.gz"
            echo "  $0 restore-latest"
            echo "  $0 list"
            echo "  $0 verify dynamic-reports-backup-20231009-120000.tar.gz"
            echo "  $0 cleanup 5"
            echo "  $0 stats"
            exit 1
            ;;
    esac
}

# Executar função principal com logging
{
    main "$@"
} 2>&1 | tee -a "$LOG_FILE"

info "Operação concluída. Log detalhado: $LOG_FILE"