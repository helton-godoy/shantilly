# Shantilly: PRD Arquitetura Híbrida e Qualidade Go-Native

O PRD a seguir está formatado em Markdown/TXT, conforme o padrão de arquivos de contexto para o **TaskMaster IA**, sintetizando os requisitos de arquitetura e qualidade definidos no *Blueprint Técnico de Engenharia de Software Focado em Contratos e Qualidade Go-Native*.

O documento enfatiza a **coerência estrutural** e a **qualidade por design**, que são os principais "Guard Rails" para a IA, desviando dos problemas de desorganização e validação tardia do projeto inicial.

***

```markdown
# prd.txt

## Overview High Level

O projeto Shantilly é uma ferramenta CLI (Command Line Interface) de código aberto construída em Go que visa substituir utilitários TUI legados (como dialog/whiptail). Sua missão é permitir que administradores de sistemas e desenvolvedores criem **Interfaces de Usuário de Terminal (TUI) ricas e interativas** de forma declarativa, utilizando arquivos de configuração YAML/JSON.

Esta fase de desenvolvimento foca na implementação de uma **Arquitetura Híbrida CLI/Servidor (Wish)** e na correção de falhas de qualidade iniciais, priorizando a **maximalização da especificação declarativa e da automação rigorosa da qualidade**.

## Features Core

As funcionalidades centrais a serem entregues neste ciclo de desenvolvimento incluem:

1.  **Arquitetura Híbrida CLI/Servidor:** Implementação completa do modo servidor SSH (`shantilly serve`) utilizando o **Charm Wish**. O comando raiz (`shantilly`) deve usar Cobra para delegar comandos TUI locais (`form`, `menu`, `layout`, `tabs`) e o novo modo remoto (`serve`).
2.  **Conjunto de Componentes Completo (Fase 4):** Finalização da implementação dos seguintes componentes:
    *   `TextArea` (Reuso obrigatório de `bubbles/textarea`).
    *   `Slider` (Lógica customizada otimizada com Lip Gloss).
    *   `FilePicker` (Correção de API e reuso de `bubbles/filepicker/v2`).
    *   `RadioGroup` (Lógica customizada de cursor e View() finalizada).
    *   Funcionalidade de `Tabs` (Orquestração de foco e validação apenas para a aba ativa).
3.  **Contrato de Componentes Rígido:** Todo componente deve implementar a interface `Component` com métodos completos para `IsValid()`, `Value()`, `SetFocus()` e `GetError()`.
4.  **Qualidade por Design:** Aplicação de tratamento de erros explícito com contexto (`fmt.Errorf("...: %w", err)`) em todas as operações de I/O e validação, garantindo que o `errcheck` não gere avisos.

## Experiência do Usuário

### Personas
*   **Engenheiro DevOps/SysAdmin:** Utiliza o Shantilly para criar dashboards leves de monitoramento via SSH ou wizards de configuração para automação em shell scripts.

### Fluxo de Utilização
1.  O usuário define a TUI em um arquivo YAML (`config.yaml`).
2.  O usuário executa a TUI localmente (`shantilly form config.yaml`) ou se conecta ao servidor SSH (`ssh user@shantilly-server`).
3.  O usuário interage com a TUI usando a navegação por teclado (`Tab`/`Shift+Tab`) e teclas de movimento.
4.  O resultado dos dados coletados é serializado em JSON para uso em scripts downstream.

### Considerações de UI/UX
*   **Estilização Adaptativa:** O modo servidor (Wish) deve usar **Custom Renderers do Lip Gloss** para detectar o perfil de cor e o fundo (`HasDarkBackground`) do terminal do cliente SSH, garantindo renderização correta e legível em ambientes remotos.
*   **Layout Responsivo:** O LayoutModel deve gerenciar corretamente o redimensionamento (propagando `tea.WindowSizeMsg`) e a divisão de espaço usando `lipgloss.JoinHorizontal` e `lipgloss.JoinVertical`.
*   **Navegação e Foco:** A navegação entre componentes focáveis deve ser robusta, ignorando elementos estáticos.

## Arquitetura Técnica

### Componentes de Sistema e Integração
*   **Pilha:** Go-Native, Charm Ecosystem (Bubble Tea v2, Lip Gloss v2, Bubbles v2, Cobra).
*   **Camadas:** A arquitetura segue a separação rígida de responsabilidades:
    *   `cmd/`: Cobra CLI e lógica de inicialização de I/O (local ou remoto/Wish).
    *   `internal/models/`: Modelos de estado (FormModel, LayoutModel) — agnósticos ao tipo de I/O.
    *   `internal/components/`: Implementações de *widgets* que aderem ao Contrato `Component`.
*   **Integração Híbrida (Wish):** O `cmd/serve.go` deve iniciar o servidor Wish, que utiliza o `bubbletea middleware` para servir o `tea.Model` agnóstico a cada sessão SSH.

### Modelos de Dados e API
*   **Configuração:** Uso mandatório de `gopkg.in/yaml.v3` para parsing declarativo de `ComponentConfig`.
*   **Contrato de Qualidade:** A interface `Component` deve ser o contrato central, garantindo que `IsValid()` e `Value()` funcionem corretamente para agregação de estado no `FormModel`.

### Requerimentos de Qualidade e Infraestrutura (Guard Rails)
*   **Cobertura de Testes:** Meta obrigatória de **85% de cobertura**. O pipeline de CI (emulado localmente) deve falhar se essa meta não for atingida.
*   **Tratamento de Erros:** `errcheck` configurado como fatal. Uso obrigatório de *error wrapping* com contexto (`fmt.Errorf("...: %w", err)`) em I/O e *parsing*.
*   **Testes E2E:** Otimizar e corrigir o problema de *timeout* nos testes End-to-End, injetando tamanho de terminal (`tea.WithWindowSize`) em ambientes não-TTY.
*   **Estilização em Testes:** Forçar o perfil de cor `TrueColor` em `TestMain` nos pacotes de UI (`internal/components`, `internal/models`) para evitar inconsistências de `View()` em CI.

## Roadmap de Desenvolvimento

### Requerimentos do MVP (Foco no Escopo Atual)
O MVP deve focar na fundação de qualidade e nas funcionalidades centrais que demonstram o poder declarativo do Shantilly e sua capacidade híbrida:

1.  Finalizar a base de qualidade (Tratamento de Erros, Linting Fatal, Configuração de Cobertura de 85%).
2.  Completar a lógica de Orquestração (LayoutModel: Foco e Redimensionamento).
3.  Finalizar e Testar os Componentes Chave (`TextArea`, `Slider`, `FilePicker`, `RadioGroup`, `Tabs`).
4.  Implementar o subcomando `shantilly serve` (Wish Server) com Custom Renderers.
5.  Atingir e manter 85% de cobertura de testes na suíte Charm.

### Melhorias Futuras (Future Enhancements)
*   Implementação de busca (`List`) e paginação (`Paginator`) avançadas para o MenuModel.
*   Adicionar `Viewport` para lidar com logs de deploy ou textos longos.
*   Suporte a `Keybindings` configuráveis pelo usuário.
*   Implementação de um sistema de Plugins e Telemetria.

## Mitigações de Riscos

### Desafios Técnicos
*   **Risco de Regressão de Qualidade:** A correção de débitos técnicos (erros de `errcheck`) e o aumento de cobertura são complexos. **Mitigação:** O requisito de 85% de cobertura (verificado via CI) é a principal mitigação para evitar regressões.
*   **Complexidade de Estado MVU:** A gestão de foco e redimensionamento em layouts aninhados é complexa. **Mitigação:** Implementação de testes abrangentes para LayoutModel que validem a propagação de mensagens e foco.
*   **Risco de Segurança (Wish):** A exposição de um servidor SSH é intrinsecamente arriscada. **Mitigação:** A arquitetura Wish/Charm oferece uma superfície de ataque menor, pois não fornece acesso a *shell*. O `cmd/serve.go` deve implementar *middlewares* (logging, access control) para auditoria e restrição de acesso.

### Constraints de Resource e Appendix
*   **Constraint:** Todas as implementações devem ser em Go e alinhadas ao ecossistema Charm v2 (Bubble Tea, Lip Gloss, Bubbles).
*   **Appendix de Especificações Técnicas:** A IA deve utilizar os documentos de contexto fornecidos como guias de implementação obrigatórios, incluindo:
    *   `docs/dev/Architecture.md`
    *   `docs/dev/Server_Design.md`
    *   `docs/dev/Component_Contract.md`
    *   `docs/dev/Testing_v2.md`
    *   `docs/dev/Linting_Rules.md`

## Logical Dependency Chain (Cadeia de Dependências)

A ordem de implementação deve priorizar a Fundação de Qualidade e os Contratos Arquiteturais (Pilares 1 e 2) antes da Implementação de Features (Pilar 3) e da Integração de Rede (Pilar 4).

1.  **QUALIDADE E FUNDAÇÃO:**
    *   Atualizar dependências Charm para v2.
    *   Implementar Tratamento de Erros Explícito (%w) e corrigir avisos de `errcheck`.
    *   Corrigir o problema de estilo Lip Gloss nos testes (forçar TrueColor em `TestMain`).
2.  **CONTRATOS E ORQUESTRAÇÃO (Modelos):**
    *   Finalizar a lógica de Redimensionamento e Propagação de Foco no LayoutModel.
    *   Definir e aplicar o Contrato rígido da Interface `Component` (View, IsValid, Value).
3.  **IMPLEMENTAÇÃO DE COMPONENTES E TESTES UNITÁRIOS:**
    *   Implementar/Finalizar `TextArea`, `Slider`, `FilePicker`, `RadioGroup`, `Tabs`.
    *   Criar testes abrangentes para LayoutModel e Modelos de Orquestração (FormModel, MenuModel).
    *   Atingir 85% de cobertura nos pacotes `internal/components` e `internal/models`.
4.  **ARQUITETURA HÍBRIDA:**
    *   Implementar o subcomando `shantilly serve` (Wish Server) com Custom Renderers.
    *   Resolver e otimizar o *timeout* nos testes E2E/Integração (Injeção de `tea.WithWindowSize`).
5.  **FINALIZAÇÃO:**
    *   Executar testes de integração final e gerar relatório de cobertura (85%+).
    *   Atualizar a documentação do projeto com as novas funcionalidades e a arquitetura híbrida.
```