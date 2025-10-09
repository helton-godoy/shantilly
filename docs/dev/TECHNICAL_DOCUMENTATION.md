# Documentação Técnica Completa - Projeto SHantilly

## 📋 Tabela de Conteúdos

1. [Introdução e Visão Geral](#1-introdução-e-visão-geral)
2. [Arquitetura do Sistema](#2-arquitetura-do-sistema)
3. [Requisitos Funcionais e Não Funcionais](#3-requisitos-funcionais-e-não-funcionais)
4. [Casos de Uso](#4-casos-de-uso)
5. [Design de Componentes](#5-design-de-componentes)
6. [Fluxos de Processo](#6-fluxos-de-processo)
7. [Guias de Implementação e Instalação](#7-guias-de-implementação-e-instalação)
8. [Casos de Teste e Procedimentos de Validação](#8-casos-de-teste-e-procedimentos-de-validação)
9. [Instruções de Uso para Usuários Finais](#9-instruções-de-uso-para-usuários-finais)
10. [Planos de Manutenção e Atualização](#10-planos-de-manutenção-e-atualização)
11. [Análise de Riscos](#11-análise-de-riscos)
12. [Glossário de Termos Técnicos](#12-glossário-de-termos-técnicos)
13. [Apêndices e Referências](#13-apêndices-e-referências)

---

## 1. Introdução e Visão Geral

### 1.1 Propósito do Projeto

O **SHantilly** é a evolução moderna do dialog para Shell Scripts. Uma ferramenta inovadora desenvolvida em Go que permite criar Interfaces de Usuário de Terminal (TUI) ricas e interativas de forma declarativa através de arquivos de configuração YAML simples e intuitivos.

Assim como montar um cupcake com ingredientes pré-definidos, o Shantilly permite que você construa interfaces complexas para automação e scripts sem complicações. Nossa filosofia é transformar a complexidade do desenvolvimento de TUIs em uma experiência simples e agradável, permitindo que você se concentre na lógica da sua automação.

O projeto utiliza o ecossistema Charm, composto pelas bibliotecas Bubble Tea, Lip Gloss e Bubbles, para fornecer uma experiência de desenvolvimento simplificada e poderosa.

### 1.2 Nossa Missão

Nossa missão é empoderar administradores de sistemas, engenheiros de DevOps e desenvolvedores a construir TUIs elegantes e poderosas para seus scripts e ferramentas de forma declarativa. Queremos transformar a complexidade do desenvolvimento de TUIs em uma experiência simples e agradável, permitindo que você se concentre na lógica da sua automação.

### 1.3 Os Ingredientes Principais

O sabor único do Shantilly vem da combinação de ingredientes essenciais:

**Configuração Declarativa em YAML**: A base da nossa receita. Defina toda a estrutura da sua TUI, desde formulários a menus, em um arquivo YAML claro e legível, perfeito para ser chamado de um script bash ou zsh.

**Uma Confeitaria de Componentes**: Oferecemos uma seleção de componentes prontos para usar — de caixas de texto a seletores de arquivos — que são os blocos de construção para qualquer TUI.

**Estilo e Flexibilidade**: Assim como a cobertura de um cupcake, a aparência importa. Com um sistema de temas simples e layouts flexíveis, suas TUIs não serão apenas funcionais, mas também visualmente agradáveis.

### 1.3.1 Por Que Open Source?

Acreditamos que as melhores receitas são aquelas compartilhadas. Manter o Shantilly como um projeto de código aberto é nosso compromisso com a colaboração, transparência e o poder da comunidade. Queremos que todos possam usar, modificar e aprimorar a ferramenta, adicionando seu próprio "toque especial" à receita.

### 1.4 Objetivos Principais

- **Desenvolvimento Declarativo**: Permitir a criação de interfaces TUI complexas através de configuração YAML sem necessidade de programação tradicional
- **Arquitetura Híbrida**: Suporte tanto para execução local quanto remota via SSH (modo servidor)
- **Componentes Reutilizáveis**: Biblioteca extensiva de componentes TUI padronizados e validados
- **Qualidade de Produção**: Manutenção rigorosa de padrões de qualidade com cobertura de testes mínima de 85%

### 1.3 Escopo Delimitado

**Incluso:**

- Framework para criação de formulários declarativos via YAML
- Componentes interativos: TextInput, TextArea, Checkbox, RadioGroup, Slider, TextLabel
- Sistema de validação integrado com regras configuráveis
- Layouts flexíveis (horizontal/vertical) com redimensionamento responsivo
- Modo servidor SSH para acesso remoto às interfaces
- Sistema de temas e estilização adaptativa (dark/light mode)

**Excluído:**

- Interface gráfica desktop (GUI)
- Aplicações web ou mobile
- Banco de dados integrado
- Sistema de autenticação/autorização complexo

### 1.4 Público-Alvo

#### Desenvolvedores

- Desenvolvedores Go que precisam criar interfaces TUI rapidamente
- Equipes que utilizam infraestrutura declarativa (YAML/JSON)
- Desenvolvedores familiarizados com conceitos de UI/UX para terminais

#### Usuários Finais

- Administradores de sistemas que utilizam ferramentas CLI
- Operadores que precisam de interfaces interativas em ambientes remotos
- Usuários que preferem interfaces TUI ao invés de web/mobile

#### Mantenedores

- Equipe de desenvolvimento responsável pela evolução do projeto
- Contribuidores open-source
- Equipe de DevOps responsável pela infraestrutura e deployment

### 1.5 Versão Atual da Documentação

**Versão**: 1.0.0
**Data**: Outubro 2025
**Status**: Documentação técnica inicial completa

### 1.6 Histórico de Versões

| Versão | Data         | Descrição                                                            |
| ------ | ------------ | -------------------------------------------------------------------- |
| 1.0.0  | Outubro 2025 | Documentação técnica inicial completa baseada na implementação atual |

---

## 2. Arquitetura do Sistema

### 2.1 Visão Geral Arquitetural

O Shantilly adota uma arquitetura híbrida CLI/Servidor baseada no padrão **MVU (Model-View-Update)** do Bubble Tea, com separação rigorosa de responsabilidades conforme especificado em [`rules/Architecture.md`](rules/Architecture.md).

```mermaid
graph TB
    subgraph "Arquitetura Híbrida"
        CLI[CLI Local<br/>shantilly]
        SSH[Servidor SSH<br/>shantilly serve]

        subgraph "Camada de Aplicação"
            CLI --> CMD[cmd/ - Comandos Cobra]
            SSH --> CMD
            CMD --> MODELS[internal/models/ - Lógica MVU]
            CMD --> COMPONENTS[internal/components/ - Widgets]
        end

        subgraph "Camada de Apresentação"
            COMPONENTS --> STYLES[internal/styles/ - Temas]
        end

        subgraph "Camada de Configuração"
            CMD --> CONFIG[internal/config/ - YAML/JSON]
        end
    end
```

### 2.2 Tecnologias Empregadas

#### Frameworks e Bibliotecas Principais

- **Linguagem**: Go 1.24.4+
- **Framework TUI**: [Bubble Tea v2](https://github.com/charmbracelet/bubbletea/v2) - Arquitetura MVU
- **Estilização**: [Lip Gloss v2](https://github.com/charmbracelet/lipgloss/v2) - Estilos declarativos
- **Componentes**: [Bubbles v2](https://github.com/charmbracelet/bubbles/v2) - Widgets reutilizáveis
- **CLI Framework**: [Cobra v1.10.1](https://github.com/spf13/cobra) - Gerenciamento de comandos
- **Formatação**: [YAML v3](https://gopkg.in/yaml.v3) - Configuração declarativa

#### Ferramentas de Desenvolvimento

- **Linting**: golangci-lint com configuração específica
- **Testes**: Testify com cobertura mínima de 85%
- **Build**: Makefile com automação completa
- **CI/CD**: Pipeline com validação de qualidade

### 2.3 Padrões Arquiteturais Adotados

#### 2.3.1 Separação de Responsabilidades (SoC)

Conforme detalhado em [`rules/Architecture.md`](rules/Architecture.md), o projeto segue rigorosamente a estrutura de pacotes Go:

| Diretório              | Responsabilidade               | Regras Arquiteturais                                             |
| ---------------------- | ------------------------------ | ---------------------------------------------------------------- |
| `cmd/`                 | Ponto de entrada, comandos CLI | Apenas lógica de CLI e inicialização. Proibida lógica TUI.       |
| `internal/models/`     | Orquestração de estado         | Agnóstico ao I/O (local/SSH). Implementa padrão `tea.Model`.     |
| `internal/components/` | Widgets de UI                  | Implementa interface `Component`. Reutiliza componentes Bubbles. |
| `internal/styles/`     | Estilização declarativa        | Apenas definições Lip Gloss.                                     |
| `internal/config/`     | Parsing de configuração        | Validação de tipos e estruturas YAML/JSON.                       |

#### 2.3.2 Padrão MVU (Model-View-Update)

Baseado na arquitetura Elm, implementada pelo Bubble Tea:

```mermaid
statechart
    [*] --> Model : Inicialização
    Model --> Model : Update(msg)
    Model --> View : Renderização
    View --> [*] : tea.Quit
```

### 2.4 Integrações Externas e Dependências

#### 2.4.1 Dependências de Produção

```go
// go.mod - Principais dependências
github.com/charmbracelet/bubbletea/v2 v2.0.0-beta1  // Arquitetura MVU
github.com/charmbracelet/lipgloss/v2 v2.0.0-beta1   // Estilização
github.com/charmbracelet/bubbles/v2 v2.0.0-beta1    // Componentes base
github.com/spf13/cobra v1.10.1                      // CLI Framework
gopkg.in/yaml.v3 v3.0.1                           // Configuração YAML
```

#### 2.4.2 Dependências de Desenvolvimento

- **Testes**: `github.com/stretchr/testify v1.11.1`
- **Build**: Make, golangci-lint
- **Cobertura**: Ferramentas integradas no Makefile

---

## 3. Requisitos Funcionais e Não Funcionais

### 3.1 Requisitos Funcionais

#### 3.1.1 RF001 - Formulários Declarativos

**Descrição**: Permitir criação de formulários através de configuração YAML
**Prioridade**: Crítica
**Casos de Uso**:

- Definição de campos TextInput com validação
- Configuração de layouts horizontais/verticais
- Personalização de temas e estilos

**Critérios de Aceitação**:

- [ ] Arquivo YAML válido deve gerar formulário funcional
- [ ] Campos obrigatórios devem ser marcados como `required: true`
- [ ] Validação deve ocorrer em tempo real durante interação
- [ ] Submissão deve gerar JSON válido com dados do formulário

#### 3.1.2 RF002 - Componentes Interativos

**Descrição**: Implementar biblioteca completa de componentes TUI
**Prioridade**: Crítica

| Componente | Funcionalidade           | Validações Suportadas                 |
| ---------- | ------------------------ | ------------------------------------- |
| TextInput  | Campo de texto editável  | min_length, max_length, pattern       |
| TextArea   | Área de texto multilinha | height, width, min_length, max_length |
| Checkbox   | Seleção booleana         | required                              |
| RadioGroup | Seleção exclusiva        | required, items (id/label)            |
| Slider     | Controle deslizante      | min, max, step, width                 |
| TextLabel  | Texto estático           | -                                     |

#### 3.1.3 RF003 - Modo Servidor SSH

**Descrição**: Permitir acesso remoto às interfaces via SSH
**Prioridade**: Alta
**Critérios de Aceitação**:

- [ ] Comando `shantilly serve` deve iniciar servidor SSH
- [ ] Conexões SSH devem apresentar TUI interativa
- [ ] Múltiplas sessões simultâneas devem ser suportadas
- [ ] Renderização adaptativa baseada no terminal cliente

#### 3.1.4 RF004 - Sistema de Layouts

**Descrição**: Suporte a layouts flexíveis e responsivos
**Prioridade**: Média
**Critérios de Aceitação**:

- [ ] Layout horizontal deve organizar componentes lado a lado
- [ ] Layout vertical deve empilhar componentes
- [ ] Redimensionamento responsivo deve funcionar automaticamente
- [ ] Quebra de linha deve ocorrer em telas pequenas

### 3.2 Requisitos Não Funcionais

#### 3.2.1 Performance

**RNF001 - Tempo de Inicialização**

- Tempo máximo: 2 segundos para inicialização completa
- Métrica: Tempo desde execução do comando até apresentação da TUI

**RNF002 - Responsividade**

- Latência máxima: 100ms para resposta a eventos de teclado
- Métrica: Tempo entre keypress e atualização visual

**RNF003 - Consumo de Memória**

- Limite máximo: 50MB por sessão TUI
- Métrica: Uso de memória residente (RSS)

#### 3.2.2 Segurança

**RNF004 - Tratamento Seguro de SSH**

- Implementar middlewares de segurança apropriados
- Não expor shell subjacente através do servidor SSH
- Validar e sanitizar entradas de configuração

**RNF005 - Tratamento de Erros**

- Zero tolerância a erros não tratados (`errcheck` fatal)
- Todas as operações críticas devem ter tratamento de erro explícito
- Mensagens de erro devem incluir contexto significativo

#### 3.2.3 Usabilidade

**RNF006 - Acessibilidade**

- Suporte a temas dark/light automaticamente detectados
- Navegação consistente via Tab/Shift+Tab
- Indicadores visuais claros de foco e validação

**RNF007 - Internacionalização**

- Estrutura preparada para múltiplos idiomas
- Mensagens de erro e help text padronizadas

#### 3.2.4 Manutenibilidade

**RNF008 - Cobertura de Testes**

- Cobertura mínima: 85% em todos os pacotes
- Testes unitários obrigatórios para lógica crítica
- Testes de integração para fluxos end-to-end

**RNF009 - Qualidade de Código**

- Linting automático com golangci-lint
- Formatação obrigatória com `gofmt` e `goimports`
- Documentação obrigatória para APIs públicas

#### 3.2.5 Compatibilidade

**RNF010 - Suporte a Terminais**

- Compatível com terminais modernos (xterm, iTerm, etc.)
- Detecção automática de capacidades de cor
- Fallback adequado para terminais limitados

**RNF011 - Sistemas Operacionais**

- Suporte completo: Linux, macOS, Windows (via WSL)
- Binários nativos para plataformas suportadas

---

## 4. Casos de Uso

### 4.1 Scripts de Shell e Ferramentas DevOps Interativas

Transforme seus scripts bash em ferramentas interativas. Em vez de longos comandos com dezenas de flags, use o Shantilly para criar uma TUI que guia o usuário através das opções.

**Exemplo Prático: Script deploy.sh**

```bash
# deploy.sh - Script interativo para deploy
#!/bin/bash

# Chama o Shantilly com configuração de deploy
shantilly form deploy-config.yaml
```

```yaml
# deploy-config.yaml
title: "Deploy Application"
components:
  - type: radiogroup
    name: ambiente
    label: "Ambiente de Deploy"
    required: true
    options:
      items:
        - id: producao
          label: "Produção"
        - id: staging
          label: "Staging"

  - type: textinput
    name: versao
    label: "Versão para Deploy"
    required: true
    placeholder: "v1.2.3"

  - type: textarea
    name: notas
    label: "Notas do Deploy"
    required: false
    options:
      height: 3

  - type: checkbox
    name: backup
    label: "Criar backup antes do deploy?"
    required: false

  - type: spinner
    name: progresso
    label: "Iniciando deploy..."
```

**Componentes Utilizados**: menu, select, spinner, textarea (para exibir logs de deploy).

### 4.2 Assistentes de Configuração (Wizards) para Aplicações

Guie seus usuários passo a passo na configuração de um projeto ou ferramenta. Um wizard interativo é muito mais amigável do que pedir para editar um arquivo de configuração manualmente.

**Exemplo Prático: Comando my-app init**

```bash
# my-app init - Inicia assistente de configuração
#!/bin/bash

# Cria configuração inicial baseada nas respostas do usuário
result=$(shantilly form init-config.yaml)

# Processa o resultado JSON
echo $result | jq '.nome_projeto, .database_url, .porta'
```

```yaml
# init-config.yaml
title: "Configuração Inicial da Aplicação"
components:
  - type: textinput
    name: nome_projeto
    label: "Nome do Projeto"
    required: true
    options:
      min_length: 2

  - type: textinput
    name: database_url
    label: "URL do Banco de Dados"
    required: true
    placeholder: "postgresql://user:pass@host:5432/db"

  - type: textinput
    name: porta
    label: "Porta da Aplicação"
    required: true
    options:
      pattern: '^\d{4,5}$'

  - type: checkbox
    name: debug
    label: "Habilitar modo debug?"
    required: false

  - type: radiogroup
    name: ambiente
    label: "Ambiente"
    required: true
    options:
      items:
        - id: desenvolvimento
          label: "Desenvolvimento"
        - id: teste
          label: "Teste"
        - id: producao
          label: "Produção"
```

**Componentes Utilizados**: form, textinput, passwordinput, checkbox, radiogroup.

### 4.3 Dashboards de Monitoramento no Terminal

Por que sair do terminal para ver o status de suas aplicações? Crie dashboards leves que exibem métricas em tempo real, status de serviços ou o progresso de tarefas longas diretamente na sua sessão SSH.

**Exemplo Prático: Script monitor.sh**

```bash
# monitor.sh - Dashboard de monitoramento
#!/bin/bash

while true; do
    # Coleta métricas atuais
    cpu=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}')
    mem=$(free | awk 'NR==2{printf "%.1f", $3*100/$2 }')

    # Atualiza dados no dashboard
    shantilly layout monitor-dashboard.yaml --data "$(generate_metrics_json)"
    sleep 5
done
```

```yaml
# monitor-dashboard.yaml
layout: horizontal
components:
  - type: textlabel
    name: titulo_principal
    label: "🏠 DASHBOARD DE MONITORAMENTO"

  - type: layout
    layout: vertical
    components:
      - type: progress
        name: cpu_usage
        label: "CPU Usage"
        options:
          min: 0
          max: 100
          current: 45

      - type: progress
        name: memoria
        label: "Memory Usage"
        options:
          min: 0
          max: 100
          current: 67

  - type: layout
    layout: vertical
    components:
      - type: tabs
        name: status_servicos
        label: "Status dos Serviços"
        options:
          tabs:
            - label: "APIs"
              content: "✅ API v1 OK\n✅ API v2 OK\n❌ API v3 DOWN"
            - label: "Banco"
              content: "✅ PostgreSQL OK\n🔄 Redis SYNC"

      - type: textlabel
        name: ultima_atualizacao
        label: "Última atualização: 14:30:25"
```

**Componentes Utilizados**: layout (horizontal), tabs, progress.

### 4.4 Automação de Tarefas com Seleção de Arquivos

Dê superpoderes aos seus scripts de automação. Permita que o usuário selecione interativamente quais arquivos processar, qual tipo de relatório gerar ou em qual diretório salvar a saída.

**Exemplo Prático: Script de Backup**

```bash
# backup-script.sh - Seleciona pastas para backup
#!/bin/bash

# Permite usuário selecionar pastas via interface TUI
folders=$(shantilly form select-folders.yaml | jq -r '.pastas[]')

# Processa cada pasta selecionada
for folder in $folders; do
    echo "Fazendo backup de: $folder"
    # Lógica de backup aqui...
done
```

```yaml
# select-folders.yaml
title: "Seleção de Pastas para Backup"
components:
  - type: filepicker
    name: pastas
    label: "Selecione as pastas para backup"
    required: true
    options:
      mode: directory
      multiple: true

  - type: checkbox
    name: opcoes
    label: "Opções de Backup"
    options:
      items:
        - id: compressao
          label: "Compressão"
        - id: encriptacao
          label: "Encriptação"
        - id: verificacao
          label: "Verificação de integridade"

  - type: slider
    name: nivel_compressao
    label: "Nível de Compressão (1-9)"
    options:
      min: 1
      max: 9
      step: 1
      default: 6

  - type: radiogroup
    name: destino
    label: "Destino do Backup"
    required: true
    options:
      items:
        - id: local
          label: "Disco Local"
        - id: nuvem
          label: "Nuvem (S3)"
        - id: ambos
          label: "Ambos"
```

**Componentes Utilizados**: filepicker, checkbox (para opções), slider (para nível de compressão).

---

## 5. Design de Componentes

### 4.1 Contrato de Interface de Componentes

Conforme especificado em [`rules/Component_Contract.md`](rules/Component_Contract.md), todos os componentes devem implementar a interface `Component` que estende `tea.Model`:

```go
type Component interface {
    // MVU architecture methods
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string

    // Component-specific methods
    Name() string
    CanFocus() bool
    SetFocus(focused bool)
    IsValid() bool
    GetError() string
    SetError(msg string)
    Value() interface{}
    SetValue(value interface{}) error
    Reset()
}
```

### 4.2 Componentes Implementados

#### 4.2.1 TextInput

**Responsabilidades**:

- Campo de texto editável com cursor
- Validação de comprimento e padrões regex
- Suporte a placeholder e texto de ajuda

**Estados**:

- Normal: Campo editável sem foco
- Focused: Campo com indicador visual de foco
- Error: Campo com mensagem de erro visível
- Disabled: Campo não interativo

#### 4.2.2 TextArea

**Responsabilidades**:

- Área de texto multilinha com scroll
- Controle de altura e largura
- Validação de comprimento

**Características Especiais**:

- Scroll automático quando conteúdo excede altura
- Preservação de posição do cursor durante edição

#### 4.2.3 Checkbox

**Responsabilidades**:

- Seleção booleana (checked/unchecked)
- Indicador visual de estado
- Validação obrigatória quando configurado

#### 4.2.4 RadioGroup

**Responsabilidades**:

- Seleção exclusiva entre múltiplas opções
- Navegação via setas do teclado
- Indicadores visuais de seleção

**Estrutura de Dados**:

```yaml
- type: radiogroup
  name: plan
  options:
    items:
      - id: free
        label: "Gratuito"
      - id: pro
        label: "Profissional"
```

#### 4.2.5 Slider

**Responsabilidades**:

- Controle deslizante numérico
- Configuração de mínimo, máximo e passo
- Indicador visual de posição

#### 4.2.6 TextLabel

**Responsabilidades**:

- Exibição de texto estático
- Não interativo (não recebe foco)
- Formatação rica com Lip Gloss

### 4.3 Padrões de Design Aplicados

#### 4.3.1 Factory Pattern para Criação de Componentes

O padrão Factory é utilizado em [`internal/components/factory.go`](internal/components/factory.go) para centralizar a criação de componentes:

```go
func NewComponents(configs []config.ComponentConfig, theme *styles.Theme) ([]Component, error) {
    components := make([]Component, 0, len(configs))

    for _, cfg := range configs {
        comp, err := NewComponent(cfg, theme)
        if err != nil {
            return nil, fmt.Errorf("erro ao criar componente %s: %w", cfg.Name, err)
        }
        components = append(components, comp)
    }

    return components, nil
}
```

#### 4.3.2 Observer Pattern para Validação

Os modelos de orquestração observam mudanças nos componentes e propagam validações:

```go
// FormModel.Update() observa mudanças e valida
func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ... lógica de navegação

    // Propaga mensagem para componente focado
    if m.focusIndex >= 0 {
        updated, cmd := m.components[m.focusIndex].Update(msg)
        m.components[m.focusIndex] = updated.(Component)
        return m, cmd
    }

    return m, nil
}
```

### 4.4 Diagramas de Relacionamento

#### 4.4.1 Relacionamento entre Componentes e Modelos

```mermaid
classDiagram
    class Component {
        <<interface>>
        +Init() tea.Cmd
        +Update(tea.Msg) (tea.Model, tea.Cmd)
        +View() string
        +Name() string
        +CanFocus() bool
        +SetFocus(bool)
        +IsValid() bool
        +Value() interface{}
        +SetValue(interface{}) error
        +Reset()
    }

    class FormModel {
        +title string
        +components[] Component
        +focusIndex int
        +Update(tea.Msg) (tea.Model, tea.Cmd)
        +View() string
        +CanSubmit() bool
        +ToJSON() []byte
    }

    class TextInput {
        +cursor int
        +value string
        +Update(tea.Msg) (tea.Model, tea.Cmd)
        +View() string
        +IsValid() bool
    }

    Component <|.. TextInput
    FormModel --> Component : gerencia
```

---

## 5. Fluxos de Processo

### 5.1 Fluxo Principal: Execução de Formulário

#### 5.1.1 Fluxo Local (CLI)

```mermaid
sequenceDiagram
    participant U as Usuário
    participant CLI as shantilly form
    participant YAML as Config YAML
    participant Model as FormModel
    participant TUI as Bubble Tea

    U->>CLI: shantilly form config.yaml
    CLI->>YAML: Lê e valida configuração
    YAML->>CLI: Configuração válida
    CLI->>Model: NewFormModel(config, theme)
    Model->>Model: Inicializa componentes
    CLI->>TUI: tea.NewProgram(model).Run()
    TUI->>U: Apresenta TUI interativa
    U->>TUI: Interage (tab, input, enter)
    TUI->>Model: Propaga eventos
    Model->>TUI: Atualiza estado
    U->>TUI: Submete formulário (Enter)
    TUI->>CLI: Retorna JSON
    CLI->>U: Exibe resultado
```

#### 5.1.2 Fluxo Remoto (SSH)

```mermaid
sequenceDiagram
    participant C as Cliente SSH
    participant S as shantilly serve
    participant W as Wish Server
    participant Sess as Sessão SSH
    participant BT as Bubble Tea
    participant FM as FormModel

    C->>S: Conecta via SSH
    S->>W: wish.NewServer()
    W->>Sess: Cria nova sessão
    Sess->>BT: bubbletea.Middleware()
    Sess->>FM: Cria FormModel
    BT->>C: Apresenta TUI remota
    C->>BT: Interage com interface
    BT->>FM: Processa eventos
    FM->>BT: Atualiza display
    C->>BT: Submete formulário
    BT->>Sess: Encerra sessão
    Sess->>C: Retorna resultado
```

### 5.2 Tratamento de Erros

#### 5.2.1 Estratégia de Error Wrapping

Conforme [`rules/Linting_Rules.md`](rules/Linting_Rules.md), todos os erros devem incluir contexto significativo:

```go
// Padrão obrigatório
if err := yaml.Unmarshal(data, &cfg); err != nil {
    return nil, fmt.Errorf("erro ao analisar configuração YAML: %w", err)
}
```

#### 5.2.2 Tratamento em Diferentes Camadas

| Camada                 | Estratégia                       | Exemplo                               |
| ---------------------- | -------------------------------- | ------------------------------------- |
| `cmd/`                 | Envolver com contexto de I/O     | "erro ao ler arquivo de configuração" |
| `internal/models/`     | Envolver com contexto de negócio | "erro ao criar modelo de formulário"  |
| `internal/components/` | Erros locais sem wrapping        | "nome do componente é obrigatório"    |

### 5.3 Estilização Adaptativa

#### 5.3.1 Detecção de Ambiente

O sistema utiliza `lipgloss.NewRenderer()` para detectar automaticamente:

- Perfil de cor do terminal (TrueColor, 256 cores, 16 cores)
- Estado dark/light mode
- Capacidades de fundo transparente

#### 5.3.2 Aplicação de Temas

```go
// Custom renderer para SSH
renderer := lipgloss.NewRenderer(sess)
style := renderer.NewStyle().
    Background(lipgloss.AdaptiveColor{
        Light: "#FFFFFF",
        Dark:  "#000000",
    })
```

---

## 6. Guias de Implementação e Instalação

### 6.1 Pré-requisitos

#### 6.1.1 Ambiente de Desenvolvimento

- **Go**: Versão 1.24.4 ou superior
- **Sistema Operacional**: Linux, macOS, ou Windows (WSL)
- **Terminal**: Terminal moderno com suporte a cores (recomendado: xterm, iTerm2)
- **Git**: Para controle de versão

#### 6.1.2 Ferramentas de Desenvolvimento

```bash
# Instalação do Go (se necessário)
wget https://go.dev/dl/go1.24.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz

# Adicionar ao PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Instalar ferramentas de qualidade
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 6.2 Instalação

#### 6.2.1 Instalação via Go Install (Recomendado)

```bash
go install github.com/helton/shantilly/cmd/shantilly@latest
```

#### 6.2.2 Compilação Manual

```bash
# Clonar repositório
git clone https://github.com/helton/shantilly.git
cd shantilly

# Compilar binário
make build

# Instalar binário
sudo cp bin/shantilly /usr/local/bin/
```

### 6.3 Configuração de Ambiente

#### 6.3.1 Variáveis de Ambiente

```bash
# Opcional: Configurar diretório padrão para configurações
export SHANTILLY_CONFIG_DIR=~/.config/shantilly

# Opcional: Configurar porta padrão para servidor SSH
export SHANTILLY_SSH_PORT=2222
```

#### 6.3.2 Configuração SSH para Modo Servidor

Para utilizar o modo servidor, é necessário configurar chaves SSH:

```bash
# Gerar chaves SSH (se necessário)
ssh-keygen -t rsa -b 4096 -f ~/.ssh/shantilly_key

# Configurar servidor (exemplo)
shantilly serve --host 0.0.0.0 --port 2222 --key ~/.ssh/shantilly_key
```

### 6.4 Scripts de Automação

#### 6.4.1 Makefile - Comandos Disponíveis

```bash
make build          # Compilar binário
make test           # Executar testes básicos
make test-race      # Testes com detector de race conditions
make coverage       # Gerar relatório de cobertura
make fmt           # Formatação automática
make lint          # Análise de linting
make ci            # Pipeline completo de CI
make clean         # Limpar artefatos de build
make install       # Instalar binário
make uninstall     # Remover binário instalado
```

#### 6.4.2 Scripts de Desenvolvimento

- [`scripts/coverage-report.sh`](scripts/coverage-report.sh): Gera relatórios de cobertura
- [`scripts/terminal-diagnostic.sh`](scripts/terminal-diagnostic.sh): Diagnóstico de terminal
- [`scripts/clean-test.sh`](scripts/clean-test.sh): Limpeza de testes

### 6.5 Resolução de Problemas Comuns

#### 6.5.1 Problema: "command not found"

**Sintoma**: `shantilly: comando não encontrado`
**Solução**:

```bash
# Verificar instalação
which shantilly
# ou
ls -la ~/go/bin/shantilly

# Adicionar ao PATH se necessário
export PATH=$PATH:~/go/bin
```

#### 6.5.2 Problema: Cores não aparecem no terminal

**Sintoma**: Interface sem cores ou com cores inadequadas
**Solução**:

```bash
# Verificar capacidades do terminal
echo $TERM

# Forçar perfil de cor (em testes)
export TERM=xterm-256color
```

#### 6.5.3 Problema: Conexão SSH recusada

**Sintoma**: `ssh: connect to host ... port ... : Connection refused`
**Solução**:

```bash
# Verificar se servidor está rodando
ps aux | grep shantilly

# Iniciar servidor
shantilly serve --host 0.0.0.0 --port 2222

# Verificar porta
netstat -tlnp | grep 2222
```

---

## 7. Casos de Teste e Procedimentos de Validação

### 7.1 Estratégia de Testes

Conforme especificado em [`rules/Testing.md`](rules/Testing.md), o projeto adota uma estratégia de testes em camadas com cobertura mínima de 85%.

#### 7.1.1 Camadas de Teste

| Camada         | Pacote                                     | Ferramentas                  | Propósito                     |
| -------------- | ------------------------------------------ | ---------------------------- | ----------------------------- |
| **Unitários**  | `internal/components/`, `internal/models/` | Testify + Table-Driven Tests | Testar lógica isolada         |
| **Integração** | `cmd/`                                     | `os/exec` + Testify          | Testar comandos CLI completos |
| **E2E**        | Sistema completo                           | Subprocessos                 | Validar fluxos completos      |

### 7.2 Casos de Teste Unitários

#### 7.2.1 Testes de Componentes

**CT001 - TextInput: Validação de Comprimento**

```go
func TestTextInput_MinLengthValidation(t *testing.T) {
    config := config.ComponentConfig{
        Name: "username",
        Type: "textinput",
        Options: map[string]interface{}{
            "min_length": 3,
        },
    }

    comp, err := components.NewComponent(config, theme)
    require.NoError(t, err)

    // Testar valor muito curto
    err = comp.SetValue("ab")
    require.NoError(t, err)
    assert.False(t, comp.IsValid())
    assert.NotEmpty(t, comp.GetError())
}
```

**CT002 - RadioGroup: Seleção Exclusiva**

```go
func TestRadioGroup_ExclusiveSelection(t *testing.T) {
    // Teste de seleção única e comportamento exclusivo
    // Verifica se apenas uma opção pode estar selecionada
}
```

#### 7.2.2 Testes de Modelos

**CT003 - FormModel: Navegação de Foco**

```go
func TestFormModel_FocusNavigation(t *testing.T) {
    // Testar Tab/Shift+Tab entre componentes focusable
    // Verificar se foco se move corretamente
    // Validar que componentes não-focusable são pulados
}
```

**CT004 - FormModel: Submissão com Validação**

```go
func TestFormModel_SubmitWithValidation(t *testing.T) {
    // Criar formulário com campos obrigatórios
    // Testar submissão sem preencher campos obrigatórios
    // Verificar que submissão é bloqueada
    // Testar submissão com todos os campos válidos
}
```

### 7.3 Casos de Teste de Integração

#### 7.3.1 CTI001 - Comando Form Completo

```go
func TestCLI_FormCommand_CompleteFlow(t *testing.T) {
    // Criar arquivo YAML temporário
    yamlContent := `
title: "Test Form"
components:
  - type: textinput
    name: name
    required: true
`

    // Executar comando via os/exec
    cmd := exec.Command("shantilly", "form", "test.yaml")
    output, err := cmd.CombinedOutput()

    // Verificar saída JSON
    assert.NoError(t, err)
    assert.Contains(t, string(output), `"name"`)
}
```

#### 7.3.2 CTI002 - Modo Servidor SSH

```go
func TestCLI_ServeCommand_SSHConnection(t *testing.T) {
    // Iniciar servidor em background
    // Conectar via SSH usando biblioteca Go
    // Testar interação completa
    // Verificar que sessão é encerrada corretamente
}
```

### 7.4 Procedimentos de Validação

#### 7.4.1 Validação de Qualidade (QA)

**Passo 1: Execução de Testes**

```bash
# Executar suíte completa
make test-race

# Verificar cobertura
make coverage

# Validar linting
make lint
```

**Passo 2: Validação Manual**

- [ ] Testar todos os componentes em diferentes terminais
- [ ] Validar layouts horizontal/vertical em diferentes tamanhos de tela
- [ ] Testar modo servidor com múltiplas conexões SSH
- [ ] Verificar temas dark/light em diferentes ambientes

**Passo 3: Testes de Regressão**

- [ ] Executar testes após cada modificação significativa
- [ ] Validar que mudanças não quebram funcionalidades existentes
- [ ] Verificar compatibilidade com versões anteriores

#### 7.4.2 Métricas de Cobertura

**Meta Obrigatória**: 85% de cobertura em todos os pacotes

| Pacote                 | Cobertura Mínima | Estratégia                            |
| ---------------------- | ---------------- | ------------------------------------- |
| `internal/components/` | 90%              | Testes unitários para cada componente |
| `internal/models/`     | 90%              | Testes de orquestração e estado       |
| `internal/config/`     | 85%              | Testes de parsing e validação         |
| `cmd/`                 | 80%              | Testes de integração CLI              |

### 7.5 Exemplos Executáveis

#### 7.5.1 Teste Básico de Formulário

```bash
# Criar arquivo de teste
cat > test-form.yaml << 'EOF'
title: "Formulário de Teste"
components:
  - type: textinput
    name: nome
    label: "Nome"
    required: true
  - type: checkbox
    name: aceitar
    label: "Aceitar termos"
    required: true
EOF

# Executar formulário
shantilly form test-form.yaml
```

#### 7.5.2 Teste de Layout Horizontal

```bash
# Criar layout de teste
cat > test-layout.yaml << 'EOF'
layout: horizontal
components:
  - type: textinput
    name: esquerda
    label: "Campo Esquerdo"
  - type: textinput
    name: direita
    label: "Campo Direito"
EOF

# Executar layout
shantilly layout test-layout.yaml
```

---

## 8. Instruções de Uso para Usuários Finais

### 8.1 Tutoriais Iniciais

#### 8.1.1 Primeiro Formulário

Este tutorial demonstra como criar e executar seu primeiro formulário Shantilly.

**Passo 1: Criar arquivo de configuração**

```bash
cat > meu-primeiro-form.yaml << 'EOF'
title: "Meu Primeiro Formulário"
description: "Aprendendo a usar o Shantilly"

components:
  - type: textinput
    name: nome
    label: "Qual é o seu nome?"
    placeholder: "Digite seu nome completo"
    required: true
    options:
      min_length: 2

  - type: textinput
    name: email
    label: "E-mail"
    placeholder: "seu@email.com"
    required: true

  - type: checkbox
    name: newsletter
    label: "Gostaria de receber nossa newsletter?"
    required: false
EOF
```

**Passo 2: Executar formulário**

```bash
shantilly form meu-primeiro-form.yaml
```

**Passo 3: Interagir com a interface**

- Use **Tab** e **Shift+Tab** para navegar entre campos
- Digite normalmente nos campos de texto
- Use **Espaço** para marcar/desmarcar checkboxes
- Pressione **Enter** para submeter quando todos os campos obrigatórios estiverem preenchidos
- Use **Esc** ou **Ctrl+C** para sair

### 8.2 Funcionalidades Principais

#### 8.2.1 Navegação por Teclado

| Tecla         | Ação                | Contexto                       |
| ------------- | ------------------- | ------------------------------ |
| **Tab**       | Próximo campo       | Qualquer componente focado     |
| **Shift+Tab** | Campo anterior      | Qualquer componente focado     |
| **Enter**     | Submeter formulário | Quando todos os campos válidos |
| **Espaço**    | Selecionar opção    | Checkbox, RadioGroup           |
| **Setas**     | Navegar opções      | RadioGroup, Slider             |
| **Esc**       | Sair                | Qualquer momento               |
| **Ctrl+C**    | Forçar saída        | Qualquer momento               |

#### 8.2.2 Indicadores Visuais

**Estados dos Componentes:**

- **Campo normal**: Borda simples, texto padrão
- **Campo focado**: Borda destacada (tema ativo)
- **Campo com erro**: Borda vermelha com mensagem de erro
- **Campo válido**: Indicador verde (quando aplicável)

#### 8.2.3 Validações em Tempo Real

O sistema valida campos conforme você interage:

- **Campos obrigatórios**: Indicados com asterisco (*) ou texto claro
- **Comprimento mínimo/máximo**: Validado durante digitação
- **Padrões regex**: Aplicados quando campo perde foco
- **Dependências**: Validações condicionais entre campos

### 8.3 Configurações Personalizáveis

#### 8.3.1 Temas e Aparência

O Shantilly detecta automaticamente se você está usando um terminal dark ou light mode e ajusta as cores accordingly. Não é necessário configuração manual para temas básicos.

#### 8.3.2 Configurações Avançadas

```yaml
# Exemplo avançado com todas as opções
title: "Configuração Avançada"
description: "Demonstração de recursos avançados"

components:
  - type: textinput
    name: senha
    label: "Senha Segura"
    required: true
    options:
      min_length: 8
      max_length: 50
      pattern: '^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$'
      help: "Mínimo 8 caracteres, maiúscula, minúscula e número"

  - type: slider
    name: prioridade
    label: "Nível de Prioridade"
    default: 5
    options:
      min: 1
      max: 10
      step: 1
      width: 50

  - type: radiogroup
    name: categoria
    label: "Categoria"
    required: true
    options:
      items:
        - id: trabalho
          label: "Trabalho"
        - id: pessoal
          label: "Pessoal"
        - id: urgente
          label: "Urgente"
```

### 8.4 Resolução de Problemas do Usuário

#### 8.4.1 Problemas Comuns

**P: O formulário não aparece corretamente**

- R: Verifique se seu terminal suporta cores (echo $TERM)
- R: Tente redimensionar o terminal ou usar um terminal diferente

**P: Não consigo navegar entre campos**

- R: Use Tab (próximo) ou Shift+Tab (anterior)
- R: Verifique se há campos obrigatórios não preenchidos

**P: Mensagens de erro não aparecem**

- R: Certifique-se de que preencheu todos os campos obrigatórios
- R: Verifique se atendeu aos critérios de validação (comprimento, formato)

**P: A interface fica "quebrada" após redimensionar**

- R: O Shantilly deve se adaptar automaticamente, mas tente sair e executar novamente

#### 8.4.2 Obtendo Ajuda

```bash
# Ajuda geral
shantilly --help

# Ajuda específica de comando
shantilly form --help
shantilly serve --help

# Versão instalada
shantilly version
```

### 8.5 Exemplos Práticos

#### 8.5.1 Formulário de Contato

```yaml
# contato.yaml
title: "Fale Conosco"
components:
  - type: textinput
    name: nome
    label: "Nome"
    required: true
  - type: textinput
    name: email
    label: "E-mail"
    required: true
    options:
      pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
  - type: textarea
    name: mensagem
    label: "Mensagem"
    required: true
    options:
      height: 5
      min_length: 10
```

#### 8.5.2 Pesquisa de Satisfação

```yaml
# pesquisa.yaml
title: "Pesquisa de Satisfação"
components:
  - type: radiogroup
    name: experiencia
    label: "Como foi sua experiência?"
    required: true
    options:
      items:
        - id: excelente
          label: "Excelente"
        - id: boa
          label: "Boa"
        - id: regular
          label: "Regular"
        - id: ruim
          label: "Ruim"

  - type: slider
    name: recomendacao
    label: "Recomendaria para amigos?"
    options:
      min: 0
      max: 10
      step: 1

  - type: textarea
    name: sugestoes
    label: "Sugestões de melhoria"
    required: false
    options:
      height: 3
```

---

## 9. Planos de Manutenção e Atualização

### 9.1 Manutenção Preventiva

#### 9.1.1 Atualizações de Segurança

- **Dependências**: Atualizar bibliotecas críticas mensalmente
- **Vulnerabilidades**: Monitorar CVE das dependências Go
- **Modo Servidor**: Revisar configurações SSH regularmente

#### 9.1.2 Otimizações de Performance

- **Monitoramento**: Implementar métricas básicas de uso de memória e CPU
- **Profiling**: Análise periódica com `go tool pprof`
- **Benchmarking**: Testes de performance antes de releases maiores

### 9.2 Processos de Versionamento

#### 9.2.1 Versionamento Semântico

Seguindo [Semantic Versioning 2.0.0](https://semver.org/):

- **PATCH** (1.0.x): Correções de bugs e melhorias de segurança
- **MINOR** (1.x.0): Novas funcionalidades compatíveis com versões anteriores
- **MAJOR** (x.0.0): Mudanças que quebram compatibilidade

#### 9.2.2 Changelog e Release Notes

Cada release deve incluir:

- Lista detalhada de mudanças
- Instruções de migração (quando aplicável)
- Notas de compatibilidade
- Melhorias de performance

### 9.3 Estratégias de Backup/Rollback

#### 9.3.1 Backup de Configurações

```bash
# Backup de configurações personalizadas
cp ~/.shantilly/config.yaml ~/.shantilly/config.yaml.backup.$(date +%Y%m%d)
```

#### 9.3.2 Rollback de Versão

```bash
# Instalar versão específica
go install github.com/helton/shantilly/cmd/shantilly@v1.0.0

# Ou usar binário anterior
cp ~/backups/shantilly-v1.0.0 ~/bin/shantilly
```

### 9.4 Monitoramento e Logging

#### 9.4.1 Eventos Monitorados

- Inicialização de componentes
- Erros de validação
- Submissões de formulário
- Conexões SSH (modo servidor)

#### 9.4.2 Níveis de Log

```go
// Exemplo de implementação futura
logger := logrus.New()
logger.SetLevel(logrus.InfoLevel)

// Eventos críticos
logger.WithField("component", "FormModel").Error("Falha na validação")

// Eventos informativos
logger.WithField("form", "cadastro").Info("Formulário submetido com sucesso")
```

### 9.5 Ciclos de Atualização

#### 9.5.1 Patches (Mensal)

- Correções de bugs críticos
- Atualizações de segurança
- Melhorias de estabilidade

#### 9.5.2 Minor Releases (Trimestral)

- Novas funcionalidades compatíveis
- Melhorias de performance
- Expansão da biblioteca de componentes

#### 9.5.3 Major Releases (Anual)

- Refatoração arquitetural
- Mudanças de API/contrato
- Adoção de novas versões de dependências críticas

### 9.6 Responsabilidades de Equipe

#### 9.6.1 Equipe de Desenvolvimento

- **Líder Técnico**: Aprovação de mudanças arquiteturais
- **Desenvolvedores**: Implementação e testes de funcionalidades
- **DevOps**: Infraestrutura e deployment

#### 9.6.2 Processo de Code Review

- **Obrigatório**: Para todas as mudanças em produção
- **Critérios**: Cobertura de testes, linting, documentação
- **Aprovadores**: Pelo menos 2 membros da equipe técnica

---

## 10. Análise de Riscos

### 10.1 Riscos Técnicos

#### 10.1.1 RT001 - Dependências Externas

**Probabilidade**: Média | **Impacto**: Alto
**Descrição**: Falha em dependências críticas (Bubble Tea, Lip Gloss)
**Mitigações**:

- Monitoramento regular de versões
- Testes com múltiplas versões
- Estratégia de fallback quando possível

#### 10.1.2 RT002 - Compatibilidade de Terminais

**Probabilidade**: Baixa | **Impacto**: Médio
**Descrição**: Interfaces podem não funcionar em terminais antigos
**Mitigações**:

- Detecção automática de capacidades
- Fallback para modos limitados
- Documentação clara de requisitos

#### 10.1.3 RT003 - Vazamentos de Memória

**Probabilidade**: Baixa | **Impacto**: Alto
**Descrição**: Aplicação pode consumir memória excessiva em uso prolongado
**Mitigações**:

- Uso de `go tool pprof` para análise
- Testes de stress com múltiplas sessões
- Monitoramento de uso de memória

### 10.2 Riscos Operacionais

#### 10.2.1 RO001 - Configurações Inválidas

**Probabilidade**: Média | **Impacto**: Médio
**Descrição**: Arquivos YAML malformados podem causar crashes
**Mitigações**:

- Validação rigorosa de configuração
- Mensagens de erro claras e contextuais
- Exemplos bem documentados

#### 10.2.2 RO002 - Indisponibilidade do Servidor SSH

**Probabilidade**: Baixa | **Impacto**: Alto
**Descrição**: Modo servidor pode ficar indisponível
**Mitigações**:

- Monitoramento de processos
- Reinicialização automática
- Balanceamento de carga para alta disponibilidade

### 10.3 Riscos de Segurança

#### 10.3.1 RS001 - Exposição de Dados Sensíveis

**Probabilidade**: Baixa | **Impacto**: Crítico
**Descrição**: Dados de formulário podem ser expostos em logs
**Mitigações**:

- Sanitização de logs
- Não armazenar dados sensíveis em logs
- Configuração segura por padrão

#### 10.3.2 RS002 - Ataques SSH

**Probabilidade**: Média | **Impacto**: Alto
**Descrição**: Servidor SSH pode ser alvo de ataques
**Mitigações**:

- Uso de chaves SSH obrigatórias
- Limitação de usuários e comandos
- Middlewares de segurança apropriados

### 10.4 Matriz de Riscos

| Risco | Probabilidade | Impacto     | Prioridade | Status      |
| ----- | ------------- | ----------- | ---------- | ----------- |
| RT001 | Média (3)     | Alto (4)    | Alta       | Monitorando |
| RT002 | Baixa (2)     | Médio (3)   | Média      | Mitigado    |
| RT003 | Baixa (2)     | Alto (4)    | Média      | Monitorando |
| RO001 | Média (3)     | Médio (3)   | Média      | Mitigado    |
| RO002 | Baixa (2)     | Alto (4)    | Média      | Planejado   |
| RS001 | Baixa (2)     | Crítico (5) | Alta       | Mitigado    |
| RS002 | Média (3)     | Alto (4)    | Alta       | Mitigado    |

### 10.5 Estratégias de Contingência

#### 10.5.1 Plano de Recuperação de Desastres

- **Backup**: Configurações e dados críticos devem ser copiados diariamente
- **Redundância**: Múltiplas instâncias em modo servidor quando necessário
- **Monitoramento**: Alertas automáticos para falhas críticas

#### 10.5.2 Procedimentos de Emergência

1. **Detecção**: Sistemas de monitoramento identificam problemas
2. **Isolamento**: Isolar componentes afetados quando possível
3. **Recuperação**: Executar procedimentos de rollback automatizados
4. **Notificação**: Alertar equipe responsável
5. **Análise**: Investigar causa raiz após recuperação

---

## 11. Glossário de Termos Técnicos

### 11.1 Termos Gerais

**API (Application Programming Interface)**

- Interface que permite comunicação entre diferentes componentes de software
- No contexto do Shantilly: Interface `Component` que todos os widgets implementam

**CLI (Command Line Interface)**

- Interface de usuário baseada em texto operada através de comandos digitados
- Modo de execução local do Shantilly: `shantilly form arquivo.yaml`

**Framework**

- Estrutura de software que fornece funcionalidades genéricas reutilizáveis
- Exemplos no projeto: Cobra (CLI), Bubble Tea (TUI), Lip Gloss (estilos)

**Middleware**

- Software que atua como intermediário entre diferentes componentes
- No contexto SSH: `bubbletea.Middleware()` conecta TUI ao servidor Wish

**TUI (Terminal User Interface)**

- Interface gráfica operada dentro de um terminal de texto
- Característica principal do Shantilly, construída com Bubble Tea

### 11.2 Termos Específicos do Projeto

**Bubble Tea**

- Framework Go para construção de TUIs usando arquitetura MVU (Model-View-Update)
- Base arquitetural do Shantilly para gerenciamento de estado e interface

**Component Contract**

- Interface rígida que todos os componentes devem implementar
- Define 11 métodos obrigatórios incluindo `Update()`, `View()`, `IsValid()`, `Value()`

**Custom Renderer**

- Instância de `lipgloss.Renderer` ligada a uma sessão específica
- Necessário para detectar ambiente (dark/light) em conexões SSH

**Error Wrapping**

- Técnica de adicionar contexto a erros usando `fmt.Errorf("contexto: %w", err)`
- Padrão obrigatório no projeto para facilitar debugging

**Focus Management**

- Sistema de navegação entre componentes interativos usando Tab/Shift+Tab
- Gerenciado pelo `FormModel` ou `LayoutModel`

**Lip Gloss**

- Biblioteca para estilização declarativa de texto no terminal
- Responsável por cores, bordas e layouts visuais no Shantilly

**MVU (Model-View-Update)**

- Padrão arquitetural baseado na linguagem Elm
- Modelo: estado da aplicação
- View: representação visual
- Update: lógica de transição de estado

**Wish**

- Biblioteca Charm para servir aplicações Bubble Tea via SSH
- Permite execução remota de TUIs através de `shantilly serve`

### 11.3 Acrônimos e Abreviações

| Acrônimo | Expansão                   | Descrição                                        |
| -------- | -------------------------- | ------------------------------------------------ |
| **BTP**  | Bubble Tea Program         | Instância executável de uma aplicação Bubble Tea |
| **E2E**  | End-to-End                 | Testes que validam fluxo completo da aplicação   |
| **MVU**  | Model-View-Update          | Arquitetura de estado reativo                    |
| **SSH**  | Secure Shell               | Protocolo para conexão remota segura             |
| **TUI**  | Terminal User Interface    | Interface de usuário no terminal                 |
| **YAML** | YAML Ain't Markup Language | Formato de configuração legível                  |

### 11.4 Referências Externas

**Documentação Técnica**:

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [Lip Gloss Guide](https://github.com/charmbracelet/lipgloss) - Estilização
- [Cobra CLI Framework](https://cobra.dev/) - Gerenciamento de comandos
- [Semantic Versioning](https://semver.org/) - Padrão de versionamento

**Comunidades**:

- [Go Community](https://go.dev/) - Linguagem de programação
- [Charm Community](https://charm.sh/) - Ecossistema TUI

---

## 12. Apêndices e Referências

### 12.1 Apêndices

#### 12.1.1 Configurações de Exemplo

**Formulário Completo de Cadastro**:

```yaml
# cadastro-completo.yaml
title: "Cadastro de Novo Usuário"
description: "Formulário completo para criação de conta"

components:
  # Informações pessoais
  - type: textinput
    name: nome_completo
    label: "Nome Completo"
    required: true
    options:
      min_length: 2
      max_length: 100

  - type: textinput
    name: email
    label: "E-mail"
    required: true
    options:
      pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'

  - type: textinput
    name: telefone
    label: "Telefone"
    required: false
    options:
      pattern: '^\(\d{2}\)\s\d{4,5}-\d{4}$'

  # Configurações de segurança
  - type: textinput
    name: senha
    label: "Senha"
    required: true
    options:
      min_length: 8
      max_length: 50
      pattern: '^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$'

  - type: textinput
    name: confirmar_senha
    label: "Confirmar Senha"
    required: true

  # Preferências
  - type: checkbox
    name: newsletter
    label: "Receber newsletter por e-mail"
    required: false

  - type: radiogroup
    name: plano
    label: "Plano de Assinatura"
    required: true
    options:
      items:
        - id: basico
          label: "Básico (Gratuito)"
        - id: premium
          label: "Premium (R$ 29,90/mês)"
        - id: empresarial
          label: "Empresarial (Sob consulta)"

  # Feedback inicial
  - type: textarea
    name: expectativas
    label: "O que você espera da plataforma?"
    required: false
    options:
      height: 4
      min_length: 10
      max_length: 500

  - type: slider
    name: satisfacao_inicial
    label: "Nível de expectativa (1-10)"
    default: 7
    options:
      min: 1
      max: 10
      step: 1
```

#### 12.1.2 Logs de Amostra

**Log de Inicialização Normal**:

```
2025-10-08T19:50:00Z INFO shantilly starting version=1.0.0
2025-10-08T19:50:00Z INFO loaded configuration file=form.yaml components=6
2025-10-08T19:50:00Z INFO initialized TextInput component=name nome_completo
2025-10-08T19:50:00Z INFO initialized TextInput component=name email
2025-10-08T19:50:00Z INFO initialized TextInput component=name senha
2025-10-08T19:50:00Z INFO initialized RadioGroup component=name plano
2025-10-08T19:50:00Z INFO starting TUI program width=120 height=24
```

**Log de Erro de Validação**:

```
2025-10-08T19:51:00Z ERROR validation failed component=email error="formato de e-mail inválido"
2025-10-08T19:51:00Z WARN form submission blocked validation_errors=2
2025-10-08T19:51:05Z INFO user interaction field=email input_length=15
2025-10-08T19:51:10Z INFO validation passed component=email
```

### 12.2 Referências

#### 12.2.1 Documentos Internos

- [`rules/Architecture.md`](rules/Architecture.md) - Especificação arquitetural completa
- [`rules/Component_Contract.md`](rules/Component_Contract.md) - Contrato de componentes detalhado
- [`rules/Linting_Rules.md`](rules/Linting_Rules.md) - Regras de qualidade e linting
- [`rules/Server_Design.md`](rules/Server_Design.md) - Design do servidor SSH
- [`rules/Styling_Rules.md`](rules/Styling_Rules.md) - Regras de estilização remota
- [`rules/Testing.md`](rules/Testing.md) - Estratégia de testes e qualidade

#### 12.2.2 Documentos Externos

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss Styling Guide](https://github.com/charmbracelet/lipgloss)
- [Cobra CLI Framework](https://cobra.dev/)
- [YAML Specification](https://yaml.org/spec/)

### 12.3 Changelog da Documentação

| Versão | Data         | Mudanças                                         |
| ------ | ------------ | ------------------------------------------------ |
| 1.0.0  | Outubro 2025 | Criação inicial da documentação técnica completa |

### 12.4 Nossa Comunidade

O Shantilly é mais do que um projeto; é uma comunidade de desenvolvedores apaixonados por criar ferramentas melhores. Acreditamos que a colaboração é o ingrediente secreto que torna um projeto open source especial.

#### 12.4.1 Nossos Pontos de Encontro

**GitHub Discussions**: O melhor lugar para começar. Faça perguntas, compartilhe suas criações, proponha novas funcionalidades e participe de conversas sobre o futuro do projeto.

- **Link**: [github.com/helton/shantilly/discussions](https://github.com/helton/shantilly/discussions)

**GitHub Issues**: Encontrou um bug ou um comportamento inesperado? O rastreador de issues é o canal oficial para relatar problemas e acompanhar o progresso das correções.

- **Link**: [github.com/helton/shantilly/issues](https://github.com/helton/shantilly/issues)

**Discord (em breve!)**: Estamos planejando um espaço para conversas em tempo real, onde a comunidade possa interagir de forma mais direta. Fique de olho para o link do nosso servidor!

#### 12.4.2 Código de Conduta

Para garantir que nossa comunidade seja um espaço acolhedor e produtivo para todos, adotamos um Código de Conduta. Esperamos que todos os membros da comunidade o sigam em todas as interações. Queremos um ambiente onde todos se sintam seguros e respeitados.

- **Link**: [github.com/helton/shantilly/blob/main/CODE_OF_CONDUCT.md](https://github.com/helton/shantilly/blob/main/CODE_OF_CONDUCT.md)

### 12.5 Como Contribuir

O Shantilly é feito pela comunidade, para a comunidade. Cada contribuição, não importa o tamanho, ajuda a tornar nossa ferramenta ainda melhor.

#### 12.5.1 Formas de Ajudar

Qualquer pessoa pode contribuir, mesmo sem escrever uma única linha de código Go:

**Reportar Bugs**: Se algo não funciona como esperado, abrir uma issue detalhada é uma das contribuições mais valiosas que você pode fazer.

**Sugerir Melhorias**: Tem uma ideia para um novo componente ou uma funcionalidade que tornaria o Shantilly ainda mais poderoso? Abra uma issue ou inicie uma discussão!

**Melhorar a Documentação**: Encontrou uma explicação confusa ou um erro de digitação? Sugira uma melhoria! Uma documentação clara é fundamental para o sucesso de qualquer projeto.

**Escrever Código**: Se você se sente confortável com Go, adoraríamos receber seu Pull Request. Corrija um bug, implemente uma nova feature ou ajude a refatorar o código.

#### 12.5.2 Guia para Contribuidores de Código

Para garantir que o processo seja tranquilo para todos, preparamos um guia completo para quem deseja contribuir com código. Ele cobre nosso fluxo de trabalho, convenções de commit, como configurar o ambiente de desenvolvimento e executar os testes.

- **Link**: [github.com/helton/shantilly/blob/main/CONTRIBUTING.md](https://github.com/helton/shantilly/blob/main/CONTRIBUTING.md)

**Pronto para fazer sua primeira contribuição?** Dê uma olhada nas nossas issues marcadas como "good first issue". Elas são um ótimo ponto de partida!

### 12.6 Perguntas Frequentes (FAQ)

#### P: O que significa "Shantilly"?

R: É um trocadilho com "Chantilly". A ideia é que o projeto adiciona uma camada "doce" e fácil de usar sobre a complexidade do desenvolvimento de TUIs, assim como o chantilly cobre um bolo.

#### P: Qual a diferença entre o Shantilly e ferramentas como dialog ou whiptail?

R: Pense no Shantilly como a evolução do dialog. Enquanto o dialog é excelente para prompts simples e modais, o Shantilly permite criar interfaces muito mais ricas e modernas com múltiplos componentes, layouts e temas, tudo a partir de um arquivo YAML. Além disso, a saída de dados do Shantilly é estruturada (JSON/YAML), facilitando a integração com seus scripts.

#### P: O Shantilly é compatível com Windows, macOS e Linux?

R: Sim! Graças à sua base no framework Bubble Tea da Charm, o Shantilly é multiplataforma e funciona perfeitamente nos principais sistemas operacionais.

#### P: Como posso customizar a aparência da minha TUI?

R: Você pode definir cores, bordas e estilos para todos os componentes através de um arquivo de tema em YAML. É simples e permite que você adapte a TUI à identidade visual do seu projeto.

#### P: Preciso ser um expert em Go para usar o Shantilly?

R: Não. Você precisa de um conhecimento básico de Go para compilar e usar o binário, mas a criação das interfaces é feita de forma declarativa em YAML, o que é ideal para quem trabalha principalmente com shell scripts.

#### P: O projeto está sendo mantido ativamente?

R: Sim! O Shantilly é um projeto ativo e em constante evolução. Você pode acompanhar nosso progresso e as últimas atualizações diretamente no nosso repositório no GitHub.

#### P: Tenho uma ideia para um novo componente. Como posso sugeri-la?

R: Ótimo! A melhor forma de sugerir novas funcionalidades é abrir uma issue no GitHub ou iniciar uma conversa na nossa seção de Discussions.

### 12.7 Contatos para Suporte

#### 12.6.1 Equipe Técnica

- **E-mail**: <suporte@shantilly.dev>
- **GitHub Issues**: [github.com/helton/shantilly/issues](https://github.com/helton/shantilly/issues)
- **Discord**: [Servidor da comunidade Charm](https://charm.sh/)

#### 12.6.2 Recursos Comunitários

- **Documentação Oficial**: [pkg.go.dev/github.com/helton/shantilly](https://pkg.go.dev/github.com/helton/shantilly)
- **Exemplos**: [`docs/examples/`](docs/examples/)
- **Wiki**: [github.com/helton/shantilly/wiki](https://github.com/helton/shantilly/wiki)

### 12.5 Índice Remissivo

**A**

- AdaptiveColor, 98, 145, 168
- Análise de riscos, 198-205
- Apêndices, 206-214
- Arquitetura do sistema, 51-74

**B**

- Bubble Tea, 12, 52, 139, 168, 190
- Backup e rollback, 183

**C**

- Casos de teste, 127-153
- Changelog, 213
- CLI Framework, 52, 139
- Cobra, 52, 139, 190
- Component Contract, 17, 75-95, 139, 190
- Componentes disponíveis, 116-184

**D**

- Design de componentes, 75-95
- Desenvolvimento, 185-186

**E**

- Error Wrapping, 140, 168
- Estratégia de testes, 127-130
- Estilização, 98, 145

**F**

- Factory Pattern, 89
- Fluxos de processo, 96-123
- Formatação, 31

**G**

- Go Module, 49
- Guias de instalação, 124-126

**I**

- Instalação, 124-126
- Integração externa, 53

**L**

- Layouts, 17, 75, 116
- Licença, 287
- Linting, 30-31, 140
- Lip Gloss, 52, 98, 139, 145, 168, 190

**M**

- Manutenção, 178-186
- MVU (Model-View-Update), 52, 75, 139, 168
- Métodos obrigatórios, 17

**O**

- Objetivos, 10

**P**

- Padrões arquiteturais, 54-56
- Público-alvo, 11
- Pré-requisitos, 124

**Q**

- Qualidade de código, 31

**R**

- Requisitos funcionais, 57-61
- Requisitos não funcionais, 61-67
- Roadmap, 265
- Rollback, 183

**S**

- Segurança, 62
- Semantic Versioning, 181
- Server Design, 141, 190
- Shantilly, 1, 10, 49
- SSH, 17, 61, 141, 190
- Styling Rules, 141, 190

**T**

- Tecnologias empregadas, 52-53
- Testes, 30, 127-153, 186
- TUI, 52, 139, 168

**U**

- Uso básico, 115-184
- Usuário final, 154-177

**V**

- Validação, 17, 127-153
- Versionamento, 181-182
- Visão geral, 9-14

**W**

- Wish, 17, 61, 141, 168, 190

**Y**

- YAML, 12, 17, 52, 139, 190

---

**Fim da Documentação Técnica**

Esta documentação foi elaborada em conformidade com os padrões estabelecidos em [`rules/Architecture.md`](rules/Architecture.md), [`rules/Component_Contract.md`](rules/Component_Contract.md), [`rules/Linting_Rules.md`](rules/Linting_Rules.md), [`rules/Server_Design.md`](rules/Server_Design.md), [`rules/Styling_Rules.md`](rules/Styling_Rules.md) e [`rules/Testing.md`](rules/Testing.md), fornecendo cobertura completa e integrada de todos os aspectos técnicos do projeto Shantilly.
