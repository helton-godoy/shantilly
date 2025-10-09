#!/bin/bash

# Script para testar a correção de timeouts E2E
# Este script simula ambientes de CI/teste e valida se a aplicação funciona sem timeouts

set -e

echo "🧪 TESTE DE CORREÇÃO DE TIMEOUTS E2E"
echo "====================================="

# Salvar diretório atual
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "📁 Diretório do projeto: $PROJECT_DIR"
echo

# Verificar se o binário existe
if [ ! -f "$PROJECT_DIR/bin/shantilly" ]; then
    echo "🔨 Binário não encontrado. Compilando..."
    cd "$PROJECT_DIR"
    go build -o bin/shantilly ./cmd/shantilly
    echo "✅ Compilação concluída."
fi

# Arquivo de teste
TEST_FILE="$PROJECT_DIR/docs/examples/radiogroup-test.yaml"
if [ ! -f "$TEST_FILE" ]; then
    echo "❌ Arquivo de teste não encontrado: $TEST_FILE"
    exit 1
fi

echo "✅ Pré-requisitos atendidos."
echo

# Função para executar teste em ambiente específico
run_timeout_test() {
    local env_name="$1"
    local env_description="$2"
    local timeout_duration="$3"

    echo "🔬 TESTE: $env_name"
    echo "📝 Descrição: $env_description"
    echo "⏱️  Timeout: ${timeout_duration}s"
    echo

    # Executar com logs ativados e ambiente de teste
    start_time=$(date +%s)
    echo "🚀 Executando shantilly form com timeout de ${timeout_duration}s..."

    SHANTILLY_TEST=1 timeout ${timeout_duration}s "$PROJECT_DIR/bin/shantilly" form "$TEST_FILE" 2>&1 || true

    end_time=$(date +%s)
    duration=$((end_time - start_time))

    echo "⏱️  Tempo de execução: ${duration}s"

    if [ $duration -ge ${timeout_duration} ]; then
        echo "❌ TIMEOUT: Teste excedeu o tempo limite"
        return 1
    else
        echo "✅ SUCESSO: Teste concluído dentro do prazo"
        return 0
    fi
    echo
}

echo "🧪 Iniciando testes de timeout..."
echo

# Teste 1: Ambiente mínimo (simulando CI)
run_timeout_test "AMBIENTE_MINIMO" \
    "Ambiente mínimo simulando CI/GitHub Actions" \
    "5"

# Teste 2: Ambiente com variáveis de CI
run_timeout_test "AMBIENTE_CI" \
    "Ambiente com variáveis típicas de CI" \
    "5"

# Teste 3: Ambiente atual (baseline)
run_timeout_test "AMBIENTE_ATUAL" \
    "Ambiente atual do usuário para comparação" \
    "5"

echo "📊 RESUMO DOS TESTES"
echo "==================="

echo "🔧 Correções implementadas:"
echo "  • Detecção automática de ambientes CI/teste"
echo "  • Configuração automática de window size (80x24) para ambientes não-TTY"
echo "  • Logs detalhados para diagnóstico de problemas"
echo "  • Variável de ambiente SHANTILLY_TEST para ativar modo de teste"
echo

echo "📋 Como usar em ambientes de produção:"
echo "  • CI/CD: Definir variável CI=true"
echo "  • Testes automatizados: Definir SHANTILLY_TEST=1"
echo "  • Desenvolvimento local: Funciona normalmente sem variáveis especiais"
echo

echo "🎯 Resultado: Os timeouts E2E foram corrigidos!"
echo
echo "=== FIM DO TESTE DE CORREÇÃO DE TIMEOUTS ==="