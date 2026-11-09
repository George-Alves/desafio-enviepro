# desafio-enviepro

CRM de estudo feito no Desafio EnviePro. Tem um menu lateral retrátil e uma aba **CRM** com:

- resumo de vendas do mês (total, ticket médio, quantidade) e clientes ativos;
- funil de vendas por etapa (Lead, Contato, Proposta, Negociação e Fechado);
- tabela de clientes com busca por nome, filtro por status e destaque em vermelho para quem está há mais de 15 dias sem contato.

Todos os números vêm do banco SQLite, passando pela API em Go. Nada está escrito direto na tela.

**Tecnologias:** Next.js 16 + React 19 + TypeScript + Tailwind (front-end), Go + Gin (back-end) e SQLite (banco).

## O que instalar antes

| Programa | Para quê | Como conferir |
|---|---|---|
| [Git](https://git-scm.com/downloads) | baixar o projeto | `git --version` |
| [Node.js](https://nodejs.org) (versão LTS) | rodar o front-end | `node -v` e `npm -v` |
| [Go](https://go.dev/dl/) 1.27 ou mais novo | rodar o back-end | `go version` |

Depois de instalar, feche e abra o terminal de novo. Não precisa instalar compilador C nem programa de banco de dados: o driver do SQLite é escrito só em Go e o banco é um arquivo.

## Como rodar

Você vai usar **dois terminais**, um para o back-end e outro para o front-end. Comece pelo back-end.

### 1. Baixar o projeto

```bash
git clone https://github.com/George-Alves/desafio-enviepro.git
cd desafio-enviepro
```

### 2. Back-end (terminal 1)

```bash
cd backend
go mod tidy
go run .
```

O banco `backend/crm.db` já vem no repositório com dados fictícios (20 clientes, 30 oportunidades e 25 vendas). Se o arquivo não existir ou estiver vazio, o programa cria e popula tudo sozinho. Quando aparecer `API rodando em http://localhost:8080`, deixe esse terminal aberto.

Para conferir, abra <http://localhost:8080/clientes> no navegador: deve aparecer uma lista em JSON.

### 3. Front-end (terminal 2)

```bash
cd frontend
npm install
npm run dev
```

Abra <http://localhost:3000> e clique em **CRM** no menu.

## Se a porta 8080 já estiver em uso

Se o `go run .` mostrar erro de porta, ou se o <http://localhost:8080/clientes> abrir uma página que não é uma lista JSON, outro programa do seu computador está usando a porta 8080. Use outra porta nos dois lados:

1. Back-end, no terminal 1:

   PowerShell:
   ```powershell
   $env:PORT="8081"; go run .
   ```
   Git Bash, macOS ou Linux:
   ```bash
   PORT=8081 go run .
   ```

2. Front-end: crie o arquivo `frontend/.env.local` com esta linha e reinicie o `npm run dev`:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8081
   ```

O front pode abrir em qualquer porta do `localhost`: se a 3000 estiver ocupada, o Next usa a 3001 sozinho e a API continua aceitando as chamadas.

## Dados fictícios

O arquivo `backend/crm.db` foi gerado pelo `seed.go` e está no repositório. As datas dos dados são calculadas a partir do dia em que o seed rodou (09/11/2026). Por isso, com o passar do tempo, o "resumo do mês" pode ficar vazio e os clientes aparecem com mais dias sem contato.

Para apagar tudo e gerar os dados de novo, com datas relativas a hoje, com o back-end parado:

```bash
cd backend
go run . seed
```

Para ver que a tela vem do banco, mude um dado (rodando o seed de novo ou editando o `crm.db`), recarregue a página do CRM e veja os números mudarem.

## Rotas da API

| Rota | O que devolve |
|---|---|
| `GET /clientes` | clientes com status, responsável, último contato, `dias_sem_contato` e `alerta` (mais de 15 dias) |
| `GET /funil` | quantidade e valor de oportunidades em cada etapa, na ordem do funil |
| `GET /vendas/resumo` | total, quantidade e ticket médio das vendas do mês atual |
| `GET /vendas` | lista de vendas, da mais recente para a mais antiga |

## Testes e verificações

```bash
cd backend
go test ./...
```

```bash
cd frontend
npx tsc --noEmit
npm run lint
npm run build
```

Os testes do back-end cobrem as regras de negócio com repositórios falsos (limite de 15 dias, funil com etapas vazias, ticket médio sem vendas), as 4 rotas com um banco de verdade e o CORS.

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
| `src/features/layout` | Menu lateral e a moldura das telas |
| `src/features/crm` | Painel do CRM: `components`, `hooks`, `services` e `types` |
| `src/components/ui` | Peças compartilhadas (`GlassCard`, `StatCard`, esqueleto de carregamento, bloco de erro) |
| `src/lib` | Cliente HTTP único (`api/client.ts`), React Query e formatadores |

### O que ficou de fora de propósito

Multi-tenant, Unit of Work, Outbox, idempotência, BFF e Strategy/Factory do padrão original não entram: este sistema tem um caso só e nenhum dinheiro em jogo, então seriam burocracia sem retorno.
