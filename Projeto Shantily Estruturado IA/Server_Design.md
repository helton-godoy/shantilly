### Especificação Técnica: Design do Servidor SSH (Wish)

#### 1. Introdução e Propósito

Este documento detalha o design arquitetural para a implementação do modo servidor (`shantilly serve` ou `shantilly daemon`) no projeto Shantilly, utilizando a biblioteca **Charm Wish**. O propósito é estender as capacidades do Shantilly para servir aplicações TUI remotamente através de SSH, mantendo a **coesão arquitetural** e garantindo que o núcleo da TUI (os `tea.Model`s em `internal/models/`) permaneça completamente **agnóstico ao I/O** (Input/Output).

O modo servidor deve ser uma **extensão limpa da CLI**, isolada no pacote `cmd/serve.go`, evitando modificações invasivas no *core* TUI.

#### 2. Arquitetura Híbrida e Separação de Responsabilidades

O fluxo híbrido separa claramente a Camada de Inicialização/Rede (responsabilidade do `cmd/serve.go`) da Camada de Estado/Lógica (responsabilidade de `internal/models/`):

| Camada | Pacote(s) | Responsabilidade |
| :--- | :--- | :--- |
| **I. Camada CLI/Rede** | `cmd/serve.go` | Define o comando Cobra, inicia o `wish.NewServer`, configura *middlewares* e lida com a injeção da sessão SSH no `tea.Program`. |
| **II. Camada de Orquestração** | `internal/models/` | Contém a lógica de estado do `FormModel`, `LayoutModel`, etc. É responsável pelo ciclo MVU (`Model`, `Update`, `View`) e navegação. |
| **III. Camada de Apresentação** | `internal/components/` | Contém a lógica de *widgets* e a interface `Component`. |

#### 3. Comando CLI e Inicialização do Servidor

O comando `serve` deve ser um subcomando Cobra, conforme definido no `cmd/root.go`. Sua função primária é configurar e iniciar o servidor Wish, que escutará conexões SSH.

##### 3.1. Estrutura do `cmd/serve.go`

A função `RunE` do `serveCmd` deve:

1.  Processar *flags* (ex: porta, endereço de *bind*, chaves SSH).
2.  Criar a função *Handler* principal que será executada para cada nova sessão SSH.
3.  Instanciar o servidor `wish.NewServer` com a cadeia de *middlewares* necessária.

##### 3.2. Exemplo: Inicialização do Servidor Wish

O *handler* deve injetar a sessão SSH no `tea.Program` e utilizar o **middleware bubbletea**.

```go
// cmd/serve.go (Estrutura Simplificada)

import (
	"github.com/charmbracelet/wish"
	// Importações de middleware de segurança, logging, e Bubble Tea
)

// Inicia o servidor SSH
func startServer(host string, port int) error {
    // Define a função de entrada da sessão SSH (Handler)
    handler := func(s ssh.Session) {
        // 1. Lógica de Negócio: Determinar qual modelo TUI rodar.
        // Por exemplo, carregar um `LayoutConfig` padrão ou verificar o comando
        // que o usuário tentou rodar via SSH.

        // 2. Criação do Modelo de Estado (Agnóstico ao I/O)
        // newModel() deve ser uma função que retorna um internal/models/XModel.
        initialModel := newModel() 

        // 3. Criação do Custom Renderer do Lip Gloss (Necessário para SSH)
        renderer := lipgloss.NewRenderer(s)

        // 4. Inicia o tea.Program com o middleware bubbletea
        p := renderer.NewProgram(initialModel) 
        
        // Configura o I/O do programa para a sessão SSH (isso é tratado pelo middleware bubbletea)
        if _, err := p.Run(); err != nil {
            // Lógica de tratamento de erro
        }
    }

    server, err := wish.NewServer(
        // Configuração de listen, chaves SSH, etc.
        &ssh.Server{Addr: fmt.Sprintf("%s:%d", host, port)}, 
        // Middlewares obrigatórios para a arquitetura Charm
        handler,
        // O middleware bubbletea envolve o handler e conecta I/O e eventos de resize
        bubbletea.Middleware(handler),
        // Adicionar middlewares de qualidade (conforme `docs/dev/Server_Design.md` exige)
        logging.Middleware(), 
        accesscontrol.Middleware(),
    )
    // ... lógica de listen
    return server.ListenAndServe()
}
```

#### 4. O Imperativo do Modelo Agnosticismo (`internal/models`)

A lógica de Model-View-Update (MVU) contida nos modelos em `internal/models/` (ex: `FormModel`, `LayoutModel`) **não deve ter conhecimento do ambiente de execução** (local CLI vs. sessão SSH).

*   **Regra de Implementação:** Os métodos `Init()`, `Update()`, e `View()` em `internal/models/` devem depender apenas de mensagens `tea.Msg` (como `tea.KeyMsg`, `tea.WindowSizeMsg`) e do estado interno do modelo.
*   **Motivação:** Isso garante que a lógica de negócio principal possa ser testada com a suíte unitária padrão (utilizando `teatest` ou *mocks*) e seja universalmente reutilizável pelo `cmd/run.go` (local) e `cmd/serve.go` (SSH).

#### 5. Estilização Remota Adaptativa (Custom Renderers)

A renderização correta de estilos e cores em sessões SSH é crítica para a qualidade. O **Lip Gloss** deve ser configurado para usar um **Custom Renderer** ligado à sessão SSH.

##### 5.1. Implementação do Renderer

No *handler* da sessão SSH (`cmd/serve.go`), o `tea.Program` e a criação de *styles* devem ser ligados a um renderizador que detecta o ambiente do cliente.

```go
// Sessão SSH é injetada no handler do Wish
func mySSHHandler(sess ssh.Session) {
    // CRÍTICO: Criar um renderizador para o cliente SSH.
    // Isso detecta o perfil de cor e o status de fundo escuro do cliente.
    renderer := lipgloss.NewRenderer(sess) 
    
    // Qualquer Lip Gloss Style criado a partir daqui será adaptativo:
    style := renderer.NewStyle().Background(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"})
    
    // O modelo inicial deve receber o tema adaptativo, que usa o Renderer
    initialModel := newModelWithAdaptiveTheme(renderer)

    p := renderer.NewProgram(initialModel)
    // ... rodar o programa
}
```

O uso de `lipgloss.NewRenderer(sess)` garante que o projeto Shantilly resolva o problema de compatibilidade e *styling* em cenários servidor-cliente.

#### 6. Middlewares Obrigatórios (Qualidade e Segurança)

Para um servidor acessível via rede, o uso de *middlewares* de qualidade é mandatório. A implementação do `cmd/serve.go` deve incluir explicitamente:

*   **Logging:** O middleware `logging.Middleware()` para rastrear conexões, desconexões e erros.
*   **Access Control:** O middleware `accesscontrol.Middleware()` para permitir a restrição de acesso a usuários específicos ou a comandos suportados. Isso limita a superfície de ataque ao TUI e não expõe um *shell* subjacente, o que é uma vantagem de segurança da arquitetura Wish sobre o OpenSSH.
