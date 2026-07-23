# Status da Implementação

**Atualizado em:** 22/07/2026

**Fase ativa:** [Fase 0 — Baseline](ROADMAP.md)

Este documento resume apenas fatos verificados. O detalhamento funcional e as evidências ficam na [matriz de compatibilidade](COMPATIBILITY.md).

## Estado atual

| Área | Estado | Evidência |
| --- | --- | --- |
| Build CMake/Qt6 | Passa no contêiner Debian 13 | `make build_internal` |
| Testes unitários | 3 casos de modelos de produção | `unit_tests` |
| CLI `--help` e `--version` | Compatível | Testes CTest dedicados |
| Entrada `stdin` e saída de `query` | Compatível no cenário mínimo | `compatibility_stdin_query` |
| Eventos interativos e widgets originais | Ainda não caracterizados | Matriz de compatibilidade |
| Extensões calendar, table e chart | Implementadas, ainda não verificadas | Sem evidência automatizada |

## Arquitetura em transição

- O ponto de entrada CLI preservado em `legacy/v1_monolith` está ativo temporariamente para manter o contrato do `dialogbox`.
- `SHantilly.cc` ainda concentra criação, mutação e relatório de widgets.
- `legacy/v2_incomplete` é compilado por compatibilidade transitória, mas não representa a arquitetura alvo.
- `libs/SHantilly-ui` contém widgets, configurações, temas, ícones e o builder em evolução.
- CMake é o sistema oficial; arquivos qmake permanecem somente como material transitório.

## Próximos critérios

1. Corrigir todas as instruções e testes que ainda apontam para caminhos antigos.
2. Adicionar caracterização para opções CLI, comandos, layouts e widgets originais.
3. Remover o V2 incompleto dos alvos ativos após cobrir o comportamento necessário.
4. Extrair protocolo, parser, execução e UI conforme os ADRs, mantendo os testes verdes.
