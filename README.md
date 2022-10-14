<img src="https://user-images.githubusercontent.com/24553367/195928856-08300a92-d243-4194-a664-1aefa977f56b.png" style="width: 450px;">

Whosbest API é um backend para a plataforma de competições e enquetes, Whosbest.
Trata-se de uma Web API construída em [Go](https://go.dev/) com [GraphQL](https://graphql.org/)
para as requisições e [WebSocket](https://developer.mozilla.org/pt-BR/docs/Web/API/WebSockets_API)
para uma análise em tempo real dos resultados. Além disso a persistência dos dados é feita em um
banco [PostgreSQL](https://www.postgresql.org/).
## Referências

 - [go-graphql](https://github.com/graphql-go/graphql)
 - [golang-jwt](https://github.com/golang-jwt/jwt)
 - [golang-migrate](https://github.com/golang-migrate/migrate)
 - [websocket](https://github.com/gorilla/websocket)

## Variáveis de Ambiente

Para rodar esse projeto, você vai precisar adicionar algumas variáveis de ambiente no seu `.env`. De forma geral,
recomendamos que o arquivo `.env.example` seja copiado, ele já fornece o mínimo necessário para iniciar o projeto,
não sendo necessário configurações adicionais.

## Rodando localmente

Clone o projeto

```bash
  git clone git@github.com:joaovicdsantos/whosbest-api.git
```

Entre no diretório do projeto

```bash
  cd whosbest-api
```

Execute as migrações

```bash
  docker compose --profile tools run migrate
```

Rode com docker compose

```bash
  docker compose up -d
```


## Documentação da API

#### Cadastrar

```http
  POST /register
```

| Body        | Tipo       | Descrição                           |
| :---------- | :--------- | :---------------------------------- |
| `username`  | `string`   | Nickname de usuário |
| `password`  | `string`   | Senha do usuário |

#### Logar

```http
  POST /login
```
| Body        | Tipo       | Descrição                           |
| :---------- | :--------- | :---------------------------------- |
| `username`  | `string`   | Nickname de usuário |
| `password`  | `string`   | Senha do usuário |


```http
  GET /graphql
```
| Body        | Tipo       | Descrição                           |
| :---------- | :--------- | :---------------------------------- |
| `query`     | `string`   | GraphQL query |

Além destes, há o endpoint relacioando ao WebSocket.
 ```http
  WEBSOCKET /ws
```
## Licença

[MIT](https://choosealicense.com/licenses/mit/)


## Autores

- [@joaovicdsantos](https://www.github.com/joaovicdsantos)

