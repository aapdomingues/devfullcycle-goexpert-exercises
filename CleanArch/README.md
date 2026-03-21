# Desafio CleanArch

## Execucao

```sh
docker compose up --build
```

Esse comando:

- sobe MySQL e RabbitMQ
- aguarda os servicos ficarem disponiveis
- aplica as migrations automaticamente no startup da aplicacao Go
- inicia a aplicacao Go

## Portas

- REST: `http://localhost:8000/order`
- GraphQL Playground: `http://localhost:8080/`
- GraphQL Query endpoint: `http://localhost:8080/query`
- gRPC: `localhost:50051`
- RabbitMQ Management: `http://localhost:15672`

## Testes manuais

Use o arquivo [api/api.http](/CleanArch/api/api.http) para:

- criar uma order via REST
- listar as orders via REST

Exemplo de query GraphQL:

```graphql
query ListOrders {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

Exemplo no Evans para gRPC:

```txt
evans -r repl -p 50051
package pb
service OrderService
call ListOrders
```
