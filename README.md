# desafio-enviepro

CRM de estudo feito no Desafio EnviePro: front-end em **Next.js + React**, back-end em **Go (Gin)** e banco **SQLite**.

> Status: em construção, um passo do tutorial por branch. As instruções completas de instalação e execução serão escritas no Passo 6.

## Arquitetura

O projeto segue a ideia do padrão de projetos do Atlantis, em tamanho reduzido para este desafio.

### Back-end (`backend/`): Clean Architecture com Portas e Adaptadores

Cada domínio (`cliente`, `funil`, `venda`) tem a mesma estrutura em `backend/internal/<domínio>/`:

| Pasta | Função |
|---|---|
| `domain/entity` | As "coisas" do negócio, sem nenhuma tecnologia |
| `domain/ports` | Interfaces do que a regra de negócio precisa (ex.: buscar clientes) |
| `app/usecase` | Regras de negócio; só conhecem as interfaces, nunca Gin nem SQLite |
| `adapter/handler` | Rotas HTTP (Gin) e o formato JSON da resposta |
| `adapter/repository` | Implementação das interfaces com SQL no SQLite |

Arquivos na raiz de `backend/`: `main.go` (liga as peças e sobe a API), `db.go` (conexão e tabelas) e `seed.go` (dados fictícios). `internal/shared/middleware` guarda o CORS.

### Front-end (`frontend/`): por features

| Pasta | Função |
|---|---|
| `src/app` | Rotas do Next.js; só compõem as telas |
| `src/features/<área>` | Tudo de uma área junto: `components`, `hooks`, `services`, `types` |
| `src/components/ui` | Peças compartilhadas (card com glassmorphism, tabela) |
| `src/lib` | Cliente HTTP único e configuração do React Query |

### O que ficou de fora de propósito

Multi-tenant, Unit of Work, Outbox, idempotência, BFF e Strategy/Factory do padrão original não entram: este sistema tem um caso só e nenhum dinheiro em jogo, então seriam burocracia sem retorno.
