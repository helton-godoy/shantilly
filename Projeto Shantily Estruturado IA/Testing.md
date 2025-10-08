### Especificação Técnica: Estratégia de Qualidade e Testes

#### 1. Introdução e Objetivo

Este documento define a estratégia de testes canônica para o projeto Shantilly, atuando como um **Guard Rail** de qualidade para o desenvolvimento conduzido por Inteligência Artificial. A arquitetura de testes deve garantir a **coerência estrutural** e a **automação rigorosa da qualidade** desde o primeiro dia, utilizando o padrão de asserções **Testify** e uma estrutura de testes em camadas, alinhada com a arquitetura Elm (Bubble Tea).

#### 2. Meta Obrigatória de Cobertura e Aplicação em CI

A robustez do código e a manutenibilidade dependem diretamente da cobertura da suíte de testes.

##### 2.1. Meta de Cobertura
A meta obrigatória de cobertura de código para o projeto Shantilly é de **85% ou mais**. Esta métrica deve ser aplicada a todos os *Pull Requests* (PRs).

##### 2.2. Pipeline de Qualidade e Falha em CI
O processo de Integração Contínua (CI) deve ser configurado para tratar a cobertura como um requisito **fatal**. O comando primário de verificação de qualidade (`make test-race`) deve ser o porteiro (gatekeeper) que aciona a validação de cobertura.

*   **Ação Mandatória:** O pipeline de CI deve falhar explicitamente se o script de verificação (`./scripts/coverage-report.sh 85`) indicar cobertura abaixo da meta, após a execução dos testes.
*   **Comando de Execução:** O comando `make test-race` deve executar os testes unitários com o detector de *data race* ativado, o que é um requisito de qualidade Charm/Go.

#### 3. Padrão de Asserções (Testify)

O uso de asserções claras e legíveis é obrigatório para facilitar a auditoria e a manutenção dos testes. A biblioteca **`github.com/stretchr/testify`** deve ser utilizada, priorizando `assert` e `require`.

| Pacote Testify | Propósito | Regra de Uso |
| :--- | :--- | :--- |
| **`assert`** | Para asserções que devem continuar o teste mesmo em caso de falha (soft assertions). | Uso padrão para validação de comportamento e estado (ex: `assert.Equal`). |
| **`require`** | Para asserções que, se falharem, devem interromper o teste imediatamente (fatal checks). | Uso obrigatório em etapas de *setup* (ex: leitura de arquivo, inicialização de modelos) onde a falha impede o restante das asserções de serem executadas de forma válida.

##### 3.1. Exemplo de Uso de Assertions

A IA deve priorizar o uso de `assert.NoError` e `assert.Contains` para testes de tratamento de erros explícitos.

```go
import (
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
)

func TestComponentInitialization(t *testing.T) {
    // Uso de require para garantir que o setup crítico seja bem-sucedido
    // (ex: validação de YAML ou criação de um mock)
    // Se falhar, o teste para.
    configData := []byte("type: textinput, name: required_field, required: true")
    
    // Simulação de falha de I/O ou Unmarshal, verificando o erro com contexto.
    _, err := yaml.Unmarshal(configData, &cfg)
    require.NoError(t, err, "A decodificação do YAML deve funcionar")

    // Uso de assert para checar a lógica de validação:
    model, err := NewFormModel(cfg, theme)
    assert.NoError(t, err, "A criação do modelo não deve retornar erro")
    
    // Asserção de Erro Explícito com Contexto
    // CRÍTICO: Verificar se a função de I/O retorna a mensagem de contexto esperada.
    errReading := os.ReadFile("non_existent_file.yaml")
    assert.Error(t, errReading)
    // Verifica se o erro possui o contexto explícito (conforme fmt.Errorf("...: %w", err))
    assert.Contains(t, errReading.Error(), "erro ao ler o arquivo de configuração")
}
```

#### 4. Estratégia de Testes em Camadas

O projeto utiliza três camadas de testes para garantir a corretude em todos os níveis: Unitário, Integração e End-to-End (E2E).

##### 4.1. Testes Unitários (Pacotes `internal/`)
Foco na lógica de estado isolada e nos contratos de componentes.

1.  **Componentes (`internal/components/`)**:
    *   **Propósito:** Isolar e validar o comportamento de cada widget (ex: `TextInput`, `Slider`).
    *   **Estratégia:** Utilizar *Table-Driven Tests* e simulação de `tea.Msg` para acionar a função `Update()`. Os testes devem cobrir todos os estados: focado/não focado, válido/inválido (validação por `IsValid()` e `GetError()`), e a saída da `View()` (testes de *snapshot* de UI).
2.  **Modelos de Orquestração (`internal/models/`)**:
    *   **Propósito:** Validar a **cola** que une os componentes, focando na orquestração de estado, navegação de foco (`focusNext`, `focusPrev`), propagação de mensagens e validação agregada (`CanSubmit()`).
    *   **Estratégia:** Simular fluxos de usuário enviando sequências de `tea.KeyMsg` e verificar se o foco se move corretamente e se o estado do formulário/layout está correto.

##### 4.2. Testes de Integração (CLI e Contratos)
Os testes de integração validam a interação entre a CLI (`cmd/`) e os Modelos (`internal/models/`).

*   **Execução de Subprocesso:** Devem utilizar o pacote `os/exec` para invocar o binário `shantilly`.
*   **Contrato de Saída:** O teste de integração é crucial para validar o contrato de saída da aplicação, garantindo que o comando `shantilly form` produza um JSON válido e com os dados esperados.

##### 4.3. Coerência Estrutural (Testes de Ambiente e Estilo)
Para garantir a coerência estrutural exigida, os testes devem incluir mecanismos para contornar problemas de ambiente de TUI:

1.  **Consistência de Renderização:** Para evitar falhas em testes de `View()` em ambientes de CI sem TTY, deve-se forçar o perfil de cor **`TrueColor`** no `TestMain` dos pacotes de UI.

    ```go
    // Exemplo: internal/components/test_main.go
    func TestMain(m *testing.M) {
        // Forçar TrueColor para garantir que os códigos ANSI de cor sejam renderizados
        lipgloss.SetColorProfile(termenv.TrueColor) 
        os.Exit(m.Run())
    }
    ```

2.  **Resolução de Timeout E2E:** A lógica de execução TUI (`cmd/run.go`) deve ser capaz de lidar com ambientes não interativos (sem TTY) para evitar *timeouts* nos testes E2E/Integração, possivelmente injetando um tamanho de janela padrão quando `go-isatty` detectar um ambiente não-terminal.

#### 5. Execução e Fluxo de Validação

Os testes devem ser executados localmente antes de qualquer *commit* e de forma obrigatória no CI, conforme definido no `Makefile`.

| Comando | Propósito | Inclusão no Pipeline de Qualidade |
| :--- | :--- | :--- |
| `make test-race` | Executa testes unitários e de integração com detecção de *race conditions*. | **Obrigatório** (falha se houver *race* ou erro de teste). |
| `./scripts/coverage-report.sh 85` | Verifica se o resultado da cobertura atingiu o mínimo de 85%. | **Obrigatório** (falha se abaixo de 85%).
| `make lint` | Executa o `golangci-lint` (incluindo `errcheck` fatal). | **Obrigatório** (falha se houver avisos de linting).

A IA deve emular este pipeline localmente e priorizar a correção de qualquer falha antes de submeter o código.
