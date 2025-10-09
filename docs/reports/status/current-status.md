# Relatório de Status - shantilly

**Gerado em:** 2025-10-09 08:14:33 UTC-04:00
**Versão:** 2025-10-09T08:14:33-04:00

## 📊 Visão Geral

- **Projeto:** shantilly
- **Versão Go:** go1.24.4
- **Último Commit:** 1081fbac
- **Data do Commit:** 2025-10-08 15:01:38 -0400

## 🧪 Cobertura de Testes

- **Percentual:** N/A%
- **Arquivo de Cobertura:** `coverage.out`
- **Relatório HTML:** `docs/reports/analysis/coverage/coverage.html`

## 🔍 Qualidade do Código

### Linting
**Status:** ⚠ Problemas
**Detalhes:** Encontrados problemas de linting (ver detalhes no CI)

### Compilação
**Status:** ✓ Sucesso
**Detalhes:** Projeto compila sem erros

## 📈 Métricas de Código

### Linhas de Código por Categoria
- Componentes: 7189 linhas
- Modelos: 2938 linhas
- Configuração: 1588 linhas
- Tratamento de Erros: 805 linhas
- Estilos: 669 linhas
- Comandos CLI: 261 linhas

**Total de linhas:** 13450 linhas de código Go

### Testes
- **Arquivos de teste:** 13 arquivos
- **Linhas de código de teste:** 5889 linhas

## 🏗️ Estrutura do Projeto

```
├── internal/
│   ├── components/     (Componentes: 7189 linhas)
│   ├── models/         (Modelos: 2938 linhas)
│   ├── config/         (Configuração: 1588 linhas)
│   ├── errors/         (Tratamento de Erros: 805 linhas)
│   └── styles/         (Estilos: 669 linhas)
├── cmd/                (Comandos CLI: 261 linhas)
└── Total: 13450 linhas
```

## 🚀 Como Executar

```bash
# Executar todos os testes
make test

# Verificar cobertura
make coverage

# Executar linting
make lint

# Build completo
make build
```

---

*Relatório gerado automaticamente por generate-status-report.sh*
