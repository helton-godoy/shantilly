### Especificação Técnica: Regras de Estilização Remota (Lip Gloss e Wish)

#### 1. Introdução e Propósito

Este documento estabelece as regras obrigatórias para a estilização no projeto Shantilly, com foco crítico na **compatibilidade em arquiteturas híbridas CLI/Servidor (Wish)**. A principal diretriz é garantir que a renderização da Interface de Usuário de Terminal (TUI) seja visualmente correta e adaptativa, detectando e respeitando as capacidades de cor e o ambiente (fundo claro ou escuro) do terminal cliente, especialmente em sessões SSH.

A falha em implementar um sistema de renderização adaptativo resulta em artefatos visuais quebrados ou ausência de cores em ambientes remotos ou de CI, onde o Lip Gloss pode remover a saída de cor por não detectar um TTY.

#### 2. Uso Obrigatório de Custom Renderers no Modo Servidor

No modo servidor (`shantilly serve`), a arquitetura deve utilizar a funcionalidade de **Custom Renderers** do Lip Gloss. Isso garante que as definições de estilo sejam ativamente ligadas à sessão de I/O remota, permitindo a detecção correta do perfil de cor e do status de fundo.

##### 2.1. Regra de Inicialização

O `tea.Program` deve ser iniciado com um renderizador personalizado, o qual é construído usando a sessão SSH (`ssh.Session`) fornecida pelo *handler* do Wish.

```go
// Exemplo em cmd/serve.go (Handler da sessão Wish)
import (
    "github.com/charmbracelet/lipgloss/v2"
    "github.com/charmbracelet/ssh"
    // ...
)

func mySSHHandler(sess ssh.Session) {
    // CRÍTICO: Cria o Renderer, atrelando a saída à sessão SSH.
    renderer := lipgloss.NewRenderer(sess)
    
    // O Model do Shantilly deve ser criado e injetado com o tema ou estilos adaptativos gerados pelo renderer.
    initialModel := createShantillyModel(renderer) 

    p := renderer.NewProgram(initialModel)
    
    // ... p.Run() é executado, utilizando o renderer.
}
```

#### 3. Estratégia de Cores Adaptativas

A IA deve priorizar o uso de cores adaptativas para garantir que a interface seja legível tanto em terminais de fundo escuro quanto de fundo claro. O Custom Renderer resolve o problema de detecção.

##### 3.1. Uso de `lipgloss.AdaptiveColor`

As definições de estilo no pacote `internal/styles/theme.go` devem utilizar a estrutura `lipgloss.AdaptiveColor` sempre que a cor for crítica para o contraste.

*   Quando um `Style` é criado usando `renderer.NewStyle()` (como no exemplo acima), ele utiliza as informações do terminal remoto para resolver a cor adaptativa.

```go
// Exemplo: Estilo de Foco Adaptativo
// Define uma cor clara para background escuro e uma cor escura para background claro.
adaptiveFocusColor := lipgloss.AdaptiveColor{
    Light: "#7D56F4", // Cor do Terminal Light
    Dark:  "#AF87FF", // Cor do Terminal Dark
}

style := renderer.NewStyle().
    Background(adaptiveFocusColor)

// O renderer garante que o estado dark/light do terminal SSH seja corretamente detectado.
```

##### 3.2. Detecção de Fundo Escuro (`HasDarkBackground`)

O Custom Renderer é essencial porque, em um cenário de servidor-cliente, o `lipgloss.HasDarkBackground(in, out)` (que é o que o renderizador usa) opera sobre a I/O da sessão, permitindo que a aplicação saiba se o cliente SSH está usando um fundo escuro.

#### 4. Orientações para Compatibilidade e Qualidade de Estilo

Para manter a coerência arquitetural (Zero Tolerância a Incoerência Estrutural), a IA deve seguir estas regras:

##### 4.1. Centralização da Estilização
Toda a definição de estilo deve ser centralizada no pacote `internal/styles/` (ex: `theme.go`), garantindo que o `View()` de qualquer componente (em `internal/components/`) utilize o `Theme` injetado, promovendo consistência.

##### 4.2. Manutenção do Layout Responsivo
A lógica de layout nos modelos de orquestração (ex: `LayoutModel.View()`) deve continuar a depender das utilidades de composição do Lip Gloss:
*   Utilização obrigatória de `lipgloss.JoinHorizontal()` e `lipgloss.JoinVertical()` para estruturar painéis e colunas, garantindo que o layout reaja corretamente às mensagens de redimensionamento (`tea.WindowSizeMsg`).

##### 4.3. Prevenção de Inconsistências em Testes (`View()` Output)
Para evitar falhas em testes unitários que validam a saída `View()` de componentes (problema identificado na Fase 2), a IA deve garantir que os arquivos de *setup* de teste (`internal/components/test_main.go` e `internal/models/test_main.go`) **forcem um perfil de cor** antes da execução dos testes.

```go
// Exemplo: internal/components/test_main.go
import (
    "github.com/charmbracelet/lipgloss/v2"
    "github.com/muesli/termenv" // Requer importação
)

func TestMain(m *testing.M) {
    // Forçar TrueColor para garantir que os códigos ANSI de cor sejam gerados
    // e o teste de snapshot da View() não falhe em ambientes não TTY.
    lipgloss.SetColorProfile(termenv.TrueColor) 
    os.Exit(m.Run())
}
```

A conformidade com essas regras garante a renderização fiel da TUI, independentemente do ambiente de execução (local ou SSH).
