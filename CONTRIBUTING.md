# Guia de Contribuição - SHantilly

Obrigado pelo interesse em contribuir com o SHantilly! Este documento define as diretrizes para garantir que as contribuições sejam integradas de forma suave e eficiente.

## 🚀 Como Começar

Antes de criar uma alteração, consulte o [roadmap](docs/ROADMAP.md), a [matriz de compatibilidade](docs/COMPATIBILITY.md), os [ADRs](docs/adr/) e o [handoff atual](docs/SESSION_HANDOFF.md). Contribuições devem pertencer à fase ativa ou justificar uma mudança de prioridade.

1. **Fork** o repositório no GitHub.
2. **Clone** seu fork localmente:

    ```bash
    git clone https://github.com/SEU_USUARIO/SHantilly.git
    cd SHantilly
    ```

3. Configure o ambiente de desenvolvimento (Docker recomendado):

    ```bash
    make build
    ```

## 🛠️ Padrões de Desenvolvimento

### Estilo de Código

Utilizamos ferramentas automáticas para manter a consistência do código. Por favor, não ignore os avisos do linter.

- **C++**: Obedece ao padrão definido em `.clang-format` (Google style com ajustes).
- **Linting**: Utilizamos `trunk`, `clang-tidy` e `cppcheck`.

Antes de enviar seu código, execute:

```bash
make format  # Formata o código
make lint    # Verifica problemas
```

### Mensagens de Commit

Seguimos a convenção [Conventional Commits](https://www.conventionalcommits.org/).

- `feat: adicionar novo widget Button`
- `fix: corrigir crash ao redimensionar janela`
- `docs: atualizar guia de instalação`
- `chore: atualizar dependências`

## 📦 Processo de Pull Request

1. Crie uma nova branch para sua feature ou correção:

    ```bash
    git checkout -b feat/minha-feature
    ```

2. Faça suas alterações e commits.
3. Garanta que os testes passem:

    ```bash
    make test
    ```

4. Envie para o seu fork:

    ```bash
    git push origin feat/minha-feature
    ```

5. Abra um Pull Request (PR) para a branch `main` do repositório oficial.
6. Preencha o template do PR com detalhes sobre o que foi alterado.

## 🧪 Testes

Contribuições sem testes podem ser rejeitadas. Se você adicionar uma nova funcionalidade, adicione um teste correspondente em `tests/`.

Mudanças herdadas do `dialogbox` exigem primeiro um teste de caracterização. Preserve sintaxe de comandos, protocolo `stdin/stdout`, opções CLI e códigos de saída. Atualize a matriz de compatibilidade na mesma contribuição.

## 📄 Licença

Ao contribuir para o SHantilly, você concorda que suas contribuições serão licenciadas sob a licença do projeto (GPLv3+).
