## 🧪 How to Test the Application

This application can be tested either **via Cloud Run** or **locally using Docker Compose**.

---

### ☁️ Cloud Run

Simply access the Cloud Run service URL and provide the `cep` as a query parameter.

**URL:** https://pos-go-expert-weather-mkgietjlja-uc.a.run.app

**Example:**
https://pos-go-expert-weather-mkgietjlja-uc.a.run.app/?cep=13042710


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
3. Access the application locally:

http://localhost:8080

Example:

http://localhost:8080/?cep=13042710
