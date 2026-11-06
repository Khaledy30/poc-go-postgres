# poc-go

API REST de usuários em Go + PostgreSQL, com arquitetura em camadas inspirada em Clean Architecture e padrão Repository/Service.

Projeto construído como estudo prático de Go: tipagem estática, gerenciamento de módulos, estrutura modular e integração com banco relacional.

---

## Tecnologias

- **Go 1.22**
- **PostgreSQL**
- **chi** — router HTTP leve
- **pgx** — driver nativo para PostgreSQL

---

## Estrutura do projeto

```
/cmd
  main.go              -> ponto de entrada, router e server
/internal
  /user
    entity.go          -> structs e tipos do domínio
    service.go         -> regras de negócio
    repository.go      -> acesso ao banco
    handler.go         -> adaptador HTTP
/db
  /migrations
    001_create_users.sql
```

Cada camada tem uma única responsabilidade. O `service` não conhece HTTP. O `handler` não conhece banco. O `repository` não conhece regra de negócio.

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