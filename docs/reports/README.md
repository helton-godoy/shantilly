# 📊 Sistema de Relatórios - Shantilly

Este documento descreve o sistema completo de relatórios implementado no projeto Shantilly, incluindo relatórios automatizados e dinâmicos.

## 📋 Visão Geral

O sistema de relatórios do Shantilly é composto por dois módulos principais:

1. **Relatórios Automatizados** - Scripts tradicionais que geram relatórios específicos
2. **Relatórios Dinâmicos** - Sistema avançado de captura automática de relatórios

## 🔧 Relatórios Automatizados

### Scripts Disponíveis

#### `scripts/generate-status-report.sh`
Gera relatório completo de status do projeto com métricas de qualidade.

```bash
# Executar diretamente
./scripts/generate-status-report.sh

# Ou via Makefile
make report
```

**Características:**
- ✅ Métricas de cobertura de testes
- ✅ Análise de qualidade (linting, compilação)
- ✅ Contagem de linhas de código por categoria
- ✅ Métricas de testes (arquivos e linhas)
- ✅ Informações do projeto (commits, versões)

#### `scripts/coverage-report.sh`
Gera relatórios detalhados de cobertura de testes.

```bash
# Gerar relatório com threshold específico
./scripts/coverage-report.sh 85
```

## 🚀 Sistema de Relatórios Dinâmicos

### Visão Geral

O sistema de relatórios dinâmicos captura automaticamente relatórios gerados durante interações, organizando-os por categoria e contexto com metadados completos.

### Scripts Principais

#### `scripts/capture-dynamic-reports.sh`
Script principal para captura e gerenciamento de relatórios dinâmicos.

```bash
# Captura interativa
./scripts/capture-dynamic-reports.sh capture "Título do Relatório"

# Captura via pipe
echo "Conteúdo do relatório" | ./scripts/capture-dynamic-reports.sh capture "Título"

# Listar relatórios
./scripts/capture-dynamic-reports.sh list [categoria]

# Buscar relatórios
./scripts/capture-dynamic-reports.sh search "palavra-chave"

# Atualizar índice
./scripts/capture-dynamic-reports.sh index

# Limpeza de cache
./scripts/capture-dynamic-reports.sh cleanup 30
```

**Características:**
- ✅ Captura automática com categorização inteligente
- ✅ Sistema de cache rotativo (50 relatórios)
- ✅ Indexação automática com metadados
- ✅ Detecção de contexto baseada em palavras-chave
- ✅ Geração de ID único (SHA256)
- ✅ Busca e listagem avançada

#### `scripts/backup-dynamic-reports.sh`
Sistema completo de backup e restauração de relatórios dinâmicos.

```bash
# Criar backup completo
./scripts/backup-dynamic-reports.sh create

# Listar backups disponíveis
./scripts/backup-dynamic-reports.sh list

# Verificar integridade de backup
./scripts/backup-dynamic-reports.sh verify arquivo.tar.gz

# Restaurar do backup mais recente
./scripts/backup-dynamic-reports.sh restore-latest

# Restaurar backup específico
./scripts/backup-dynamic-reports.sh restore arquivo.tar.gz

# Mostrar estatísticas
./scripts/backup-dynamic-reports.sh stats

# Limpar backups antigos
./scripts/backup-dynamic-reports.sh cleanup 5
```

**Características:**
- ✅ Backup completo com metadados
- ✅ Verificação de integridade automática
- ✅ Restauração granular e completa
- ✅ Logs detalhados de todas operações
- ✅ Rotação automática de backups
- ✅ Checksums para validação

### Integração com Makefile

Novos targets adicionados ao Makefile:

```bash
# Capturar relatório dinâmico interativo
make capture-reports

# Listar todos os relatórios dinâmicos
make list-dynamic-reports

# Criar backup dos relatórios dinâmicos
make backup-reports

# Restaurar do backup mais recente
make restore-reports

# Mostrar estatísticas do sistema de relatórios
make reports-stats
```

### Estrutura de Arquivos

```
docs/reports/
├── README.md                           # Esta documentação
├── dynamic/                           # Relatórios dinâmicos
│   ├── progress/                      # Relatórios de progresso
│   ├── analysis/                      # Análises e métricas
│   ├── status/                        # Status do projeto
│   ├── cache/                         # Cache de relatórios recentes
│   └── dynamic-reports-index.md       # Índice automático
├── backup/                           # Backups de relatórios
├── logs/                             # Logs de operações
├── status/                           # Relatórios de status automatizados
├── progress/                         # Relatórios de progresso manuais
├── analysis/                         # Análises técnicas
└── technical/                        # Documentação técnica
```

## 📊 Categorias de Relatórios

### Categorização Automática

O sistema utiliza palavras-chave para categorizar automaticamente os relatórios:

| Categoria | Palavras-chave | Descrição |
|:----------|:---------------|:----------|
| **progress** | progresso, desenvolvimento, fase, entrega, macro tarefas | Relatórios sobre andamento do projeto |
| **analysis** | análise, métricas, cobertura, performance, qualidade | Análises técnicas e métricas |
| **status** | status, estado, situação, bloqueadores | Status atual e problemas |
| **general** | Padrão quando não se encaixa nas categorias acima | Relatórios gerais |

### Extração de Contexto

Baseado no conteúdo, o sistema identifica contextos específicos:

- **macro-tarefas** - Análise de fases e entregas
- **documentacao** - Estrutura e organização de documentos
- **implementacao** - Desenvolvimento de sistemas
- **localizacao** - Mapeamento de localização de recursos
- **testes** - Cobertura e qualidade de testes
- **seguranca** - Problemas e soluções de segurança
- **geral** - Contexto padrão

## 🔍 Como Usar

### 1. Captura Básica

```bash
# Modo interativo
make capture-reports

# Via script direto
./scripts/capture-dynamic-reports.sh capture "Meu Relatório"
```

### 2. Listagem e Consulta

```bash
# Listar todos os relatórios
make list-dynamic-reports

# Listar por categoria específica
./scripts/capture-dynamic-reports.sh list progress

# Buscar por palavra-chave
./scripts/capture-dynamic-reports.sh search "teste"
```

### 3. Backup e Restauração

```bash
# Criar backup
make backup-reports

# Ver estatísticas
make reports-stats

# Listar backups disponíveis
./scripts/backup-dynamic-reports.sh list

# Restaurar se necessário
make restore-reports
```

### 4. Manutenção

```bash
# Atualizar índice manualmente
./scripts/capture-dynamic-reports.sh index

# Limpar cache antigo (30 dias)
./scripts/capture-dynamic-reports.sh cleanup 30

# Limpar backups antigos (manter últimos 5)
./scripts/backup-dynamic-reports.sh cleanup 5
```

## 📈 Recursos Avançados

### Sistema de Cache

- **Capacidade:** 50 relatórios recentes
- **Rotação:** FIFO (First In, First Out)
- **Performance:** Acesso otimizado aos relatórios mais recentes
- **Configuração:** Capacidade ajustável conforme necessidade

### Indexação Inteligente

- **Atualização automática:** Índice atualizado a cada captura
- **Metadados completos:** Timestamp, categoria, contexto, título
- **Estatísticas:** Contagem por categoria
- **Links diretos:** Navegação fácil entre relatórios

### Integração com Sistema Existente

- ✅ **Compatibilidade total** com `generate-status-report.sh`
- ✅ **Estrutura consistente** com relatórios existentes
- ✅ **Sem conflitos** entre sistemas automatizados e dinâmicos
- ✅ **Makefile integrado** com novos targets específicos

## 🛠️ Desenvolvimento e Manutenção

### Logs de Operações

Todos os logs são armazenados em:
```
docs/reports/logs/
├── backup-AAAAMMDD-HHMMSS.log    # Logs de backup
└── [outras operações de log]
```

### Verificação de Integridade

```bash
# Verificar integridade dos relatórios
./scripts/backup-dynamic-reports.sh verify arquivo.tar.gz

# Verificar estrutura completa
./scripts/backup-dynamic-reports.sh create  # Com verificação automática
```

### Monitoramento

```bash
# Estatísticas gerais
make reports-stats

# Verificar saúde do sistema
./scripts/backup-dynamic-reports.sh stats
```

## 🎯 Benefícios Alcançados

### Automatização Completa
- ✅ Zero intervenção manual necessária
- ✅ Captura automática durante interações
- ✅ Categorização inteligente baseada em IA

### Organização Avançada
- ✅ Estrutura clara e hierárquica por categoria/contexto
- ✅ Metadados completos para rastreabilidade total
- ✅ Indexação automática e sempre atualizada

### Performance Otimizada
- ✅ Sistema de cache para acesso rápido
- ✅ Processamento eficiente de texto
- ✅ Geração rápida de IDs únicos

### Confiabilidade
- ✅ Sistema de backup robusto com verificação
- ✅ Tratamento completo de erros
- ✅ Logs detalhados de todas operações

### Escalabilidade
- ✅ Suporte a crescimento ilimitado de relatórios
- ✅ Cache configurável conforme necessidade
- ✅ Busca eficiente em grandes volumes

## 📋 Status de Implementação

| Componente | Status | Descrição |
|:-----------|:-------|:----------|
| **Script de Captura** | ✅ **COMPLETO** | Todas funcionalidades implementadas |
| **Sistema de Backup** | ✅ **COMPLETO** | Backup e restauração funcionais |
| **Integração Makefile** | ✅ **COMPLETO** | Targets adicionados e testados |
| **Documentação** | ✅ **COMPLETO** | Documentação abrangente |
| **Relatórios de Exemplo** | ✅ **4 CAPTURADOS** | Exemplos práticos incluídos |

---

*Sistema de relatórios implementado em $(date '+%Y-%m-%d')*
*Versão: 1.0.0 - Produção*