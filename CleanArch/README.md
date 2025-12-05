# 🧾 Desafio CleanArch

## 🚀 Pré-requisitos

Certifique-se de ter instalado:

- **Go** (versão recomendada: 1.24+)  
- **Docker** e **Docker Compose**  
- **Make**

---

## 📦 Configuração do Ambiente

Siga os passos abaixo após clonar o repositório:

### 1️⃣ Instalar dependências Go

Após o clone, execute:

```sh
go mod tidy
```

### 2️⃣ Subir os containers necessários

```sh
docker-compose up -d
```

Isso iniciará os serviços auxiliares exigidos pela aplicação (banco de dados e rabbitMQ).

### 3️⃣ Executar as migrations

```sh
make migrate
```

Esse comando criará automaticamente a(s) tabela(s) necessárias para o funcionamento da aplicação.

### ▶️ Executando a Aplicação

Navegue até o diretório principal do módulo:

```sh
cd cmd/ordersystem/
```

E execute:

```sh
go run main.go wire_gen.go
```

A aplicação estará pronta para receber requisições.

### 🌐 Testando a API (REST)

Você pode enviar requisições REST utilizando o arquivo:

- api/api.http

Ele contém exemplos prontos para uso em extensões como REST Client (VS Code)

### 🌐 Testando via GraphQL ou GRPC

- GraphQL

Para comunicação via GraphQL você pode utilizar o console playground do próprio GraphQl através do link http://localhost:8080/

```graphql
query queryOrders{
  orders{
    id
    Price
    Tax
    FinalPrice
  }
}
```

- GRPC

Você pode utilizar o próprio client Evans para as chamadas GRPC.