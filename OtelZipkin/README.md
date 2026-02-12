## 🧪 How to Test the Application

This application can be tested **locally using Docker Compose**.

---

### 🐳 Local (Docker Compose)

#### Prerequisites
- Docker
- Docker Compose
- A valid API key from **WeatherAPI**  
  👉 https://www.weatherapi.com/

---

#### Steps

1. Set the WeatherAPI key in your environment variables (`.env` file or `docker-compose.yml`).
2. Build and start the application:
   ```bash
   docker compose up --build
   ```
3. Make a request to the application locally:

If you have the REST Client extension installed in VS Code, you can use the api/api.http file located at the project root to send requests. Alternatively, you may use your preferred API tool to make the requests as demonstrated below.

```bash
POST http://localhost:8080/weather HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "cep":"13042710"
}
```

4. Access http://localhost:9411/ to see the traces/spans