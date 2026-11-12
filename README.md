# poc-go

API REST de usuários em Go + PostgreSQL, com arquitetura em camadas inspirada em Clean Architecture e padrão Repository/Service.

Projeto construído como estudo prático de Go, vindo de um background em Node.js, TypeScript, NestJS e SST (Serverless Stack Tool). A ideia é usar uma stack conhecida (API REST + banco relacional) como ponto de comparação para entender as diferenças de paradigma: tipagem estática, gerenciamento de módulos, tratamento de erros sem exceções e a forma mais explícita como Go lida com estrutura e dependências.

---

## Tecnologias

- **Go 1.22**
- **PostgreSQL**
- **chi** — router HTTP leve
- **pgx** — driver nativo para PostgreSQL
- **sqlc** — gera código Go type-safe a partir de queries SQL puras, evitando scans manuais e mantendo o schema do banco como fonte de verdade

---

## Estrutura do projeto

/cmd
  main.go              -> ponto de entrada, router e server
/internal
  /user
    entity.go          -> structs e tipos do domínio
    service.go         -> regras de negócio
    repository.go      -> acesso ao banco
    handler.go         -> adaptador HTTP
  /database
    db.go              -> conexão e configuração do banco
    models.go          -> structs geradas (sqlc) a partir do schema
    query.sql.go       -> queries geradas (sqlc) a partir do query.sql
  /validator
    /email
      email_test.go    -> testes do validador de email
      email.go         -> regras de validação de email
/db
  /migrations
    001_create_users.sql -> migration de criação da tabela users
  query.sql              -> queries SQL usadas pelo sqlc
  schema.sql             -> schema do banco usado pelo sqlc

Cada camada tem uma única responsabilidade e depende apenas de abstrações, nunca de implementações concretas das outras:

- O `handler` conhece o `service`, mas não sabe como os dados são persistidos.
- O `service` conhece o `repository` (via interface), mas não sabe como o SQL é executado.
- O `repository` executa o acesso ao banco, mas não decide regra de negócio.
- O `validator` não depende de nenhuma camada acima — só recebe e valida o dado.

---

## Como rodar

### 1. Criar o banco

No Beekeeper ou qualquer cliente SQL:

```sql
CREATE DATABASE poc_go;
```

### 2. Rodar a migration

Abra e execute o arquivo `db/migrations/001_create_users.sql` no seu cliente SQL.

### 3. Configurar as credenciais

Se o seu PostgreSQL usar configurações diferentes do padrão, exporte antes de rodar:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=sua_senha
export DB_NAME=poc_go
```

Se não exportar, o projeto usa `postgres` como valor padrão para todos os campos.

### 4. Instalar dependências e rodar

```bash
go mod tidy
go run cmd/main.go
```

Servidor sobe na porta `8080`.

---

## Endpoints

| Método | Rota | Descrição |
| --- | --- | --- |
| GET | /health | Health check |
| POST | /users | Cria usuário |
| GET | /users | Lista usuários com paginação |
| GET | /users/{id} | Busca usuário por ID |
| DELETE | /users/{id} | Soft delete de usuário, apenas preenchendo o campo deleted_at |

---

## Exemplos de requisição

### Criar usuário

```http
POST /users
Content-Type: application/json

{
  "name": "Nome",
  "email": "nome@email.com"
}
```

### Listar com paginação

```http
GET /users?page=1&limit=10
```

### Buscar por ID

```http
GET /users/1
```

### Soft delete

```http
DELETE /users/1
```

---

## O que aprendi neste projeto

- **Tipagem estática e structs em Go:** como modelar o domínio com segurança em tempo de compilação, sem os atalhos do TypeScript.
- **go.mod e go.sum:** o `go.mod` gerencia as dependências e suas versões; o `go.sum` garante a integridade criptográfica dos pacotes baixados.
- **Arquitetura em camadas:** separar handler, service e repository deixa cada parte testável de forma independente, seguindo a mesma lógica da arquitetura hexagonal que uso no dia a dia com Node.js.
- **Tratamento de erros em Go:** sem exceções, o erro é um valor retornado explicitamente. Isso obriga a tratar cada falha no lugar certo, o que torna o fluxo mais previsível.
- **sqlc:** gera código Go type-safe a partir de queries SQL puras (`query.sql` → `query.sql.go`), evitando escrever manualmente scans de `*sql.Rows` e reduzindo o risco de erros de tipagem entre banco e aplicação, sem o overhead de uma ORM tradicional.
- **Separação de validação (clean code):** extrair a validação de email para seu próprio pacote (`/validator/email`) manteve o `service` focado em regra de negócio, sem misturar responsabilidades. Deixa o código mais fácil de testar e reutilizar caso surjam outros validadores no projeto.
- **Testes unitários com o pacote `testing`:** aprendi o padrão table-driven test — organizar casos de teste em um slice de structs e rodar cada um como subteste com `t.Run`. Isso deixa fácil adicionar novos cenários sem duplicar código, e identifica exatamente qual caso falhou no output.

## Como testar

O projeto usa o pacote `testing` nativo do Go, com testes table-driven.

### Rodar todos os testes

```bash
go test ./...
```

### Rodar com mais detalhes (mostra cada subteste)

```bash
go test ./... -v
```

### Rodar testes de um pacote específico

```bash
go test ./internal/validator/email/...
```

### Ver cobertura de testes

```bash
go test ./... -cover
```