### Especificação Técnica: Contrato de Interface de Componentes (Component_Contract.md)

#### 1. Introdução e Propósito do Contrato

Este documento estabelece o contrato rígido da interface `Component` (`internal/components/component.go`). O objetivo deste contrato é garantir a **coerência estrutural** e a interoperabilidade de todos os *widgets* de Interface de Usuário de Terminal (TUI) no projeto Shantilly, prevenindo a introdução de lógica fragmentada ou incompleta.

A adesão rigorosa a este contrato é um **Guard Rail** mandatório para o desenvolvimento conduzido por IA. **Todo novo componente deve implementar completamente todos os 11 métodos listados** para ser reconhecido e orquestrado corretamente pelos modelos de estado (`FormModel`, `LayoutModel`, etc.).

#### 2. Definição Canônica da Interface `Component`

A interface `Component` estende o contrato base da arquitetura Bubble Tea (`tea.Model`) e adiciona métodos essenciais para gerenciamento de foco, validação e serialização de dados.

| Método | Tipo | Descrição Obrigatória (Propósito) |
| :--- | :--- | :--- |
| `Init()` | `tea.Cmd` | Inicializa o estado do componente. Retorna `nil` ou um `tea.Cmd` para operações assíncronas iniciais. |
| `Update(tea.Msg)` | `(tea.Model, tea.Cmd)` | **Núcleo do MVU.** Processa mensagens de entrada (`tea.Msg`) e retorna o novo `tea.Model` e comandos opcionais. |
| `View()` | `string` | **Núcleo do MVU.** Renderiza o estado atual do componente usando **Lipgloss**. |
| `Name()` | `string` | Retorna o identificador único (`name`) do componente, obrigatório para a serialização de dados. |
| `CanFocus()` | `bool` | Indica se o componente é interativo e pode receber foco (exclui componentes estáticos como `Text`). |
| `SetFocus(bool)` | `void` | Define o estado de foco do componente, utilizado para estilização e delegação de eventos. |
| **`IsValid()`** | `bool` | **Contrato de Qualidade.** Retorna `true` se o estado do componente (incluindo validações configuráveis) estiver válido. |
| **`GetError()`** | `string` | Retorna a mensagem de erro atual. Deve retornar uma *string* vazia se `IsValid()` for `true`. |
| `SetError(string)` | `void` | Define manualmente uma mensagem de erro (usado por validação de *business logic* de nível superior). |
| **`Value()`** | `interface{}` | **Contrato de Serialização.** Retorna o valor atual do componente (ex: `string`, `bool`, `float64`) para serialização JSON/YAML final. |
| `SetValue(interface{})` | `error` | Define o valor do componente de forma programática (usado para inicialização ou reset). |
| `Reset()` | `void` | Retorna o estado do componente (valor e erro) ao seu estado inicial. |

#### 3. Requisitos de Implementação Mandatórios

Todo componente deve aderir aos seguintes requisitos de MVU, Validação e Serialização:

##### 3.1. Lógica de Atualização (`Update()`)

O método `Update()` deve isolar a manipulação de teclas de navegação global (`Tab`, `Shift+Tab`) e delegar mensagens específicas do componente para o modelo interno (`bubbles` ou lógica customizada).

*   **Delegação:** O `Update()` do componente deve ser o único responsável por passar mensagens (ex: `tea.KeyMsg`) para o modelo `textinput.Model` interno.
*   **Foco/Estilo:** O método deve usar o estado `focused` (definido por `SetFocus(bool)`) para aplicar os estilos Lip Gloss corretos na `View()` (ex: `t.theme.InputFocused` vs `t.theme.Input`).

##### 3.2. Validação de Contrato (`IsValid()` e `GetError()`)

A coerência estrutural exige que a validação seja rigorosa e baseada na configuração declarativa (`ComponentConfig`).

1.  **Validação Requerida:** Se `config.Required` for `true`, `IsValid()` deve ser `false` se o `Value()` estiver vazio.
2.  **Validação de Opções:** A lógica de `IsValid()` deve verificar todas as regras de `Options` (ex: `min_length`, `max_length`, `pattern` em `TextInput`).
3.  **Localização do Erro:** O erro retornado por `GetError()` deve ser armazenado internamente no campo `errorMsg` do componente e deve ser exibido na `View()` quando `IsValid()` for `false`.

##### 3.3. Serialização de Valor (`Value()`)

O método `Value()` é o contrato de saída do componente para o `FormModel.ToJSON()`.

*   **Tipagem Rígida:** `Value()` deve retornar o tipo de dado esperado (ex: `string` para `textinput`, `bool` para `checkbox`, `float64` para `slider`) para garantir a correta serialização em JSON/YAML.

#### 4. Regras Anti-Fragmentação e Coerência Estrutural

Para garantir que a IA priorize a coerência sobre a velocidade, as seguintes regras são obrigatórias:

1.  **Implementação Completa da Interface:** O contrato `Component` deve ser implementado na íntegra. Não é permitido deixar métodos sem implementação (mesmo que com `return nil` ou corpo vazio, como em `SetError(string)`), e a implementação deve estar completa na primeira versão do componente.
2.  **Reuso Máximo de Bubbles:** Novos componentes interativos devem **maximizar o reuso** de implementações existentes da biblioteca `charmbracelet/bubbles` (como `textarea`, `filepicker`, `textinput`).
3.  **Redundância Arquitetural Proibida:** Componentes de orquestração de alto nível (ex: `Tabs` em `internal/components/`) **não devem duplicar** a lógica de orquestração de estado que pertence aos modelos em `internal/models/` (`TabsModel`, `FormModel`). A lógica complexa de navegação (`focusNext`, `focusPrev`) pertence ao modelo pai (Form ou Layout).

#### 5. Exemplos de Implementação (Adesão ao Contrato)

##### 5.1. Exemplo: `RadioGroup`

O `RadioGroup` utiliza lógica customizada e Lip Gloss para renderização. Sua conformidade com o contrato exige:

*   **`View()`:** Deve aplicar o estilo de foco (`rg.theme.Focused`) à linha do item que corresponde ao `rg.cursor`, e usar caracteres Unicode (ex: `(•)`) para indicar o item em `rg.selected`.
*   **`IsValid()`:** Deve verificar se `rg.selected` é diferente de `-1` se o campo for `required`.
*   **`Update()`:** Deve manipular `tea.KeyMsg` (`up`/`down`, `enter`/`space`) para alterar `rg.cursor` e `rg.selected`, e só deve permitir a alteração se o componente estiver `focused`.
*   **`Value()`:** Deve retornar o `ID` (string) do `RadioItem` selecionado, e não o índice.

##### 5.2. Exemplo: `Slider`

O `Slider` usa lógica customizada e o Lip Gloss para construir a barra de progresso interativa. Sua conformidade com o contrato exige:

*   **`Update()`:** Deve lidar com `tea.KeyLeft` e `tea.KeyRight` para ajustar o valor (`value`) e garantir que ele permaneça dentro dos limites (`min`/`max`) definidos na `config.Options`.
*   **`View()`:** Deve usar `lipgloss.Overlay` ou técnicas análogas do Lip Gloss para renderizar a barra de progresso e o valor atual de forma declarativa, aplicando o estilo de foco/erro conforme a necessidade.
*   **`Value()`:** Deve retornar o valor numérico (`float64` ou `int`) atual e não a representação em *string* da `View()`.
