# Shantilly

**Construtor de TUI Declarativo via YAML**

Shantilly é uma ferramenta CLI moderna em Go que permite criar Interfaces de Usuário de Terminal (TUI) ricas e interativas de forma declarativa, utilizando arquivos de configuração YAML.

[INSERIR DEMONSTRAÇÃO VISUAL AQUI (GIF/VÍDEO)]

## 🎯 Características

- **Declarativo**: Defina sua TUI em YAML, sem escrever código

- **Componentes Ricos**: TextInput, TextArea, Checkbox, RadioGroup, Slider

- **Validação Integrada**: Validação automática de campos obrigatórios, tamanhos, padrões

- **Layouts Flexíveis**: Horizontal e vertical, com redimensionamento responsivo

- **Estilização Adaptativa**: Suporte para terminais dark/light com Lip Gloss

- **Arquitetura MVU**: Construído sobre Bubble Tea (Arquitetura Elm)

## 🚀 Instalação

```
go install [github.com/helton/shantilly/cmd/shantilly@latest](https://github.com/helton/shantilly/cmd/shantilly@latest)
```

Ou clone e compile:

```
git clone [https://github.com/helton/shantilly.git](https://github.com/helton/shantilly.git)
cd shantilly
make build
```

## 📖 Uso Rápido

### Formulário Simples

Crie um arquivo `form.yaml`:

```
title: "Cadastro de Usuário"
description: "Formulário simples de cadastro"

components:
  - type: textinput
    name: username
    label: "Nome de usuário"
    required: true
    options:
      min_length: 3

  - type: textinput
    name: email
    label: "E-mail"
    required: true

  - type: checkbox
    name: terms
    label: "Aceito os termos"
    required: true
```

Execute:

```
shantilly form form.yaml
```

A saída será JSON:

```
{
  "username": "johndoe",
  "email": "john@example.com",
  "terms": true
}
```

### Layout Horizontal

```
layout: horizontal
components:
  - type: textinput
    name: host
    label: "Host"
  - type: slider
    name: port
    label: "Porta"
    options:
      min: 1024
      max: 65535
```

Execute:

```
shantilly layout layout.yaml
```

Este comando renderiza um layout TUI sem a lógica de um formulário, ideal para dashboards ou exibições estáticas.

Saída esperada:

```
+-----------------+ +-------------------------+
| Host            | | Porta                   |
|                 | |                         |
| [ input aqui ]  | | [========> slider < ]   |
+-----------------+ +-------------------------+
```

## 📦 Componentes Disponíveis

### TextInput

```
- type: textinput
  name: username
  label: "Nome de usuário"
  placeholder: "Digite aqui..."
  required: true
  options:
    min_length: 3
    max_length: 20
    pattern: '^[a-zA-Z0-9_]+$'
```

### TextArea

```
- type: textarea
  name: description
  label: "Descrição"
  required: false
  options:
    min_length: 10
    height: 5
    width: 50
```

### Checkbox

```
- type: checkbox
  name: agree
  label: "Concordo com os termos"
  required: true
  default: false
```

### RadioGroup

```
- type: radiogroup
  name: plan
  label: "Escolha um plano"
  required: true
  options:
    items:
      - id: free
        label: "Gratuito"
      - id: pro
        label: "Profissional"
```

### Slider

```
- type: slider
  name: volume
  label: "Volume"
  default: 50
  options:
    min: 0
    max: 100
    step: 5
    width: 30
```

Para ver esses componentes em ação, confira os exemplos completos na seção "🎨 Exemplos".

## 🛠️ Desenvolvimento

### Pré-requisitos

- Go 1.21+

- golangci-lint

### Compilação

```
make build
```

### Testes

```
make test          # Testes básicos
make test-race     # Com race detector
make coverage      # Cobertura (mínimo 85%)
```

### Qualidade

```
make fmt           # Formatação
make lint          # Linting (errcheck fatal)
make ci            # Pipeline completo
```

## 🏗️ Arquitetura

Shantilly segue o padrão Charm:

```
cmd/
├── shantilly/
│   ├── main.go
│   └── commands/
│       ├── root.go
│       ├── form.go
│       └── layout.go
internal/
├── components/      # Widgets (TextInput, Slider, etc.)
├── models/          # Orquestração (FormModel, LayoutModel)
├── config/          # Parsing YAML
└── styles/          # Temas Lip Gloss
```

### Princípios

- **MVU (Model-View-Update)**: Arquitetura Elm via Bubble Tea

- **Contrato de Componentes**: Interface rígida com 11 métodos

- **Tratamento de Erros**: Explícito com contexto (`fmt.Errorf(...: %w)`)

- **85% Cobertura**: Bloqueio em CI se abaixo da meta

## 📚 Documentação Completa

- [Arquitetura](https://www.google.com/search?q=docs/Architecture.md "null")

- [Contrato de Componentes](https://www.google.com/search?q=docs/Component_Contract.md "null")

- [Regras de Qualidade](https://www.google.com/search?q=docs/Linting_Rules.md "null")

- [Estratégia de Testes](https://www.google.com/search?q=docs/Testing.md "null")

## 🎨 Exemplos

Veja exemplos completos em [`docs/examples/`](https://www.google.com/search?q=docs/examples/ "null"):

- `simple-form.yaml`: Formulário de cadastro completo

- `horizontal-layout.yaml`: Dashboard com layout horizontal

- `vertical-layout.yaml`: Questionário de feedback

## 🗺️ Roadmap

- **SSH Ready**: Suporte para modo servidor (Wish), permitindo o acesso às TUIs via SSH.

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:

1. Fork o projeto

2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)

3. Commit suas mudanças usando Conventional Commits

4. Execute `make ci` para verificar qualidade

5. Push para a branch (`git push origin feature/AmazingFeature`)

6. Abra um Pull Request

## 📝 Licença

MIT License - veja [LICENSE](https://www.google.com/search?q=LICENSE "null") para detalhes.

## 🙏 Agradecimentos

Construído com:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea "null") - Framework TUI

- [Lip Gloss](https://github.com/charmbracelet/lipgloss "null") - Estilização

- [Bubbles](https://github.com/charmbracelet/bubbles "null") - Componentes

- [Cobra](https://github.com/spf13/cobra "null") - CLI Framework
