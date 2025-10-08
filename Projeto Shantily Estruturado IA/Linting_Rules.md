### Especificação Técnica: Regras de Linting e Qualidade Go-Native

#### 1. Introdução e Propósito

Este documento define os padrões obrigatórios de formatação e tratamento de erros para o código-fonte do projeto Shantilly. Estes padrões servem como **Guard Rails** para a Inteligência Artificial (IA), garantindo que a qualidade e a coerência estrutural sejam priorizadas sobre a velocidade de implementação. O linter **golangci-lint** é a ferramenta de aplicação primária.

#### 2. Formatação e Coerência de Estilo (gofmt e goimports)

A formatação e a organização das importações devem ser consistentes em todo o repositório.

##### 2.1. Regra de Formatação Obrigatória

A IA deve garantir que todo o código-fonte esteja formatado de acordo com as convenções canônicas do Go.

*   **Ferramentas Mandatórias:** `gofmt` e `goimports`.
*   **Execução:** O comando `make fmt` deve ser executado antes de cada *commit* para garantir que todos os arquivos estejam em conformidade.

#### 3. Zero Tolerância a Erros Não Tratados (errcheck Fatal)

A negligência na verificação de erros (`errcheck`) é uma má prática que leva a bugs silenciosos e imprevisíveis. O tratamento explícito de todos os retornos de erro é mandatório.

##### 3.1. Configuração do Linter
A configuração do `golangci-lint` (via `.golangci.yml`) deve ser rigidamente definida para que o linter `errcheck` seja considerado uma falha fatal.

*   **Regra de CI:** Um *Pull Request* (PR) **não deve ser mesclado** se o `make lint` (que inclui `errcheck`) falhar.
*   **Ação Mandatória:** Toda chamada de função que retorna um `error` deve ter este valor explicitamente verificado e tratado. Isso se aplica especialmente a operações de I/O, manipulação de TUI e chamadas de CLI.

#### 4. Tratamento de Erros Explícito com Contexto (Error Wrapping)

O *error wrapping* (envolvimento de erro) é crucial para o *debugging* e a rastreabilidade, especialmente em fluxos complexos de CLI e execução TUI.

##### 4.1. Regra de Contextualização Obrigatória
Em qualquer função que retorne um erro vindo de uma camada inferior (I/O, bibliotecas externas, ou chamadas internas), o erro deve ser envolvido com uma mensagem que adicione contexto sobre o ponto de falha no fluxo de negócio.

*   **Padrão de Implementação:** Utilizar `fmt.Errorf("mensagem de contexto: %w", err)`. O `%w` é obrigatório para envolver o erro original, permitindo inspeções posteriores com `errors.Is` ou `errors.As`.

#### 5. Regras Específicas para Operações Críticas de I/O

As operações de entrada e saída são os pontos mais críticos onde o contexto de erro é vital para o diagnóstico.

##### 5.1. Leitura de Arquivos (os.ReadFile)
Ao ler arquivos de configuração (ex: YAML/JSON) no pacote `cmd/`, o erro original de I/O deve ser envolvido com o contexto do arquivo.

**Cenário (Cmd Layer):**

| Incorreto | Correto (Obrigatório) |
| :--- | :--- |
| ```go configData, err := os.ReadFile(filePath) if err != nil { return err }``` | ```go configData, err := os.ReadFile(filePath) if err != nil { return fmt.Errorf("erro ao ler o arquivo de configuração: %w", err) }``` |

##### 5.2. Deserialização de Configuração (yaml.Unmarshal)
A falha na decodificação do YAML (`gopkg.in/yaml.v3`) deve sempre ser envolta com contexto que indique qual tipo de configuração falhou.

**Cenário (Cmd Layer):**

| Incorreto | Correto (Obrigatório) |
| :--- | :--- |
| ```go if err := yaml.Unmarshal(data, &cfg); err != nil { return nil, err }``` | ```go if err := yaml.Unmarshal(data, &cfg); err != nil { return nil, fmt.Errorf("erro ao analisar o YAML de configuração: %w", err) }``` |

##### 5.3. Escrita de Saída (os.WriteFile, fmt.Fprintf)
Operações de escrita devem garantir que o erro seja tratado, especialmente no `cmd/` e `internal/models/` (ex: `ToJSON`).

*   **I/O CLI (Cmd/Root):** Em comandos CLI que realizam I/O para `stdout` ou `stderr` (ex: `fmt.Fprintf` para imprimir versão), o erro de escrita deve ser checado e, se necessário, o fluxo da função deve ser interrompido (`return`).
*   **Serialização de Saída (Model):**

```go
// Exemplo: shantilly_internal_models_form.go.txt (Trecho de ToJSON)
if err != nil {
    return fmt.Errorf("erro ao serializar dados: %w", err) // Adiciona contexto à falha
}
// O erro de os.WriteFile deve ser tratado na chamada (cmd/form.go)
return os.WriteFile(filePath, outData, 0600) // Este erro é propagado, mas precisa ser envolvido na camada cmd/
```

A IA deve garantir que os erros de `os.WriteFile` sejam envolvidos na camada CLI (`cmd/form.go`) para adicionar contexto de "falha ao salvar a saída".

##### 5.4. Lógica Interna de Componentes (internal/components/)
A inicialização de componentes que dependem de configurações complexas (como `TextInput` com `regexp`) ou que falham em validação interna devem retornar erros informativos.

*   **Exemplo (`internal/config/types.go`):** O método `Validate()` deve retornar um erro claro se o tipo de componente for inválido ou o `Name` estiver vazio.

```go
// Exemplo: Validação de Name em ComponentConfig
if c.Name == "" {
    return fmt.Errorf("nome do componente é obrigatório") // Erro explícito sem wrapping, pois é local
}
```
A adesão rigorosa a estas regras de linting e tratamento de erros garante a robustez exigida para o projeto Shantilly e facilita a depuração em cenários de execução híbrida (CLI local e Servidor Wish).
