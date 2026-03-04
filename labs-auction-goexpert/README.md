## 🧪 How to Test the Application

This application can be tested **locally using Docker Compose**.

### 🐳 Local (Docker Compose)

#### Prerequisites
- Docker
- Docker Compose

---

#### Steps

1. Build and start the application:
   ```bash
   docker compose up --build
   ```
3. Make a request to the application locally to create an auction:

If you have the REST Client extension installed in VS Code, you can use the api/api.http file located at the project root to send requests. Alternatively, you may use your preferred API tool to make the requests as demonstrated below.

```bash
POST http://localhost:8080/auction HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "product_name":"Smartphone",
    "category":"Electronics",
    "description":"Smartphone with high performance",
    "condition":1
}
```

4. After creating the auction, you can retrieve the auction ID directly from MongoDB using a client of your preference.

5. Check the auction status.

Send a GET request as shown in the example below (you can also use the api.http file in the project root) using this ID.
After the duration configured in the .env file (AUCTION_DURATION) has elapsed, you will see that the auction status changes.

```bash
GET http://localhost:8080/auction/7ebb8193-06d8-4b60-acae-921d2926e0b5 HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

---

## ✅ Automated Test (Makefile)

There is an **integration test** that validates the auction is automatically marked as `Completed` after `AUCTION_DURATION`.

### Prerequisites
- Docker + Docker Compose (MongoDB)
- Make

### Run

Run everything with a single command:
```bash
make test
```