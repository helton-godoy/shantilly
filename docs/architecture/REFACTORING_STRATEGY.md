# Estratégia de Refatoração SHantilly 2026: Arquitetura Orientada a Studio

> **Status:** Em Progresso (Fase Híbrida)
> **Data:** 14 de Janeiro de 2026
> **Objetivo:** Elevar a qualidade de código (Premium Code Quality) e permitir a criação do SHantilly Studio através do desacoplamento do núcleo.

---

## 1. Visão Geral da Arquitetura

Para permitir que o **SHantilly Studio** e o **SHantilly CLI** coexistam com 100% de paridade visual e funcional, o código monolítico atual será refatorado em uma arquitetura de camadas baseada em **Biblioteca Compartilhada**.

### O Conceito "Core Library"

O código fonte será reorganizado para isolar a lógica de UI da lógica de parsing de texto.

```mermaid
graph TD
    subgraph "Aplicações (Consumers)"
        CLI[SHantilly CLI] -- Lê Stdin --> Lib
        Studio[SHantilly Studio] -- Interação Visual --> Lib
    end

    subgraph "SHantilly Shared Ecosystem"
        Lib[libs/SHantilly-ui]
    end

    CLI --> Lib
    Studio --> Lib
```

---

## 2. Estratégia de Build

Conforme o [ADR 0002](../adr/0002-cmake-primary-build.md), CMake é o sistema oficial de build, testes, cobertura e empacotamento.

1. **Build oficial (CMake):**
    - Configurado na raiz e nos subdiretórios `src/code/shantilly`, `libs/SHantilly-ui` e `tests`.
    - Produz o executável `build/bin/shantilly` e registra a suíte no CTest.
    - É o único build que pode fornecer evidência para o roadmap e a matriz de compatibilidade.

2. **QMake transitório:**
    - O arquivo `src/code/shantilly/SHantilly.pro` permanece apenas para comparação histórica.
    - Não define o comportamento esperado da CI e será removido após a migração da cobertura necessária.

---

## 3. Log de Migração de Componentes

Este registro rastreia quais componentes foram movidos do monólito (`src/code/shantilly`) para a biblioteca (`libs/SHantilly-ui`).

| Componente     | Data       | Motivo da Migração                                      | Dependências           |
| :------------- | :--------- | :------------------------------------------------------ | :--------------------- |
| **IconHelper** | 14/01/2026 | Utilitário independente necessário para UI.             | `Logger` (Moveu junto) |
| **Logger**     | 14/01/2026 | Dependência direta de `IconHelper` e usada globalmente. | Nenhuma (Qt Core)      |

---

## 4. Design Patterns Adotados

### 4.1. Builder Pattern (O Coração da Construção)

Atualmente, o SHantilly constrói widgets "on-the-fly" enquanto lê o texto. Isso impede o Studio de instanciar um widget sem simular um arquivo de texto.

**Solução:** Implementar o `SHantillyBuilder`.

- **Intenção:** Separar a construção de um objeto complexo da sua representação.
- **Aplicação:**
  - O CLI lê `add pushbutton "OK" btn_ok` -> Chama `builder->createButton("OK", "btn_ok")`.
  - O Studio recebe um Drag & Drop -> Chama `builder->createButton("OK", "btn_ok")`.
- **Benefício:** Garante que a inicialização de propriedades, estilos padrão e conexões de sinais internos sejam idênticas em ambos os apps.

### 4.2. Passive View (MVP - Model View Presenter)

Os Widgets Qt atuais contêm lógica de execução de shell script (ex: `system()`, `popen()`). Isso é uma violação do Princípio de Responsabilidade Única (SRP).

**Solução:** Tornar os Widgets "Passivos".

- **View (Widget):** Apenas exibe dados e emite sinais Qt puros (`signal: clicked()`). Não sabe o que é "Shell Script".
- **Presenter (Controller):**
  - No **CLI**: O Presenter conecta `clicked()` a uma função que escreve no stdout ou executa um comando.
  - No **Studio**: O Presenter conecta `clicked()` a uma função que seleciona o widget no Property Editor.

---

## 5. Estrutura de Diretórios Proposta

```text
/
├── CMakeLists.txt          # Build system mestre (Novo)
├── libs/
│   └── SHantilly-ui/         # A Biblioteca (Static/Shared Lib)
│       ├── include/        # Headers públicos (icon_helper.h, logger.h)
│       └── src/            # Implementação (icon_helper.cpp, logger.cpp)
├── src/
│   └── code/
│       └── SHantilly/        # Aplicação CLI (Legacy path mantido)
│           ├── SHantilly.pro # Build Legacy (Aponta para ../../../libs)
│           └── ...
└── apps/
    └── studio/             # O novo Editor (IDE)
```
