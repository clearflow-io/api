# ClearFlow Backend

ClearFlow is an intelligent expenses tracker designed to help users gain deep insights into their spending habits by automatically processing credit card bills and bank statements.

## 🌟 Vision

ClearFlow aims to eliminate the manual effort of tracking expenses. By simply uploading financial statements, users receive automated categorizations, trend analysis, and actionable insights to improve their financial health.

## 🏗️ Architecture

The project follows a layered architecture to ensure separation of concerns, testability, and maintainability.

```mermaid
graph TD
    Client[Client Browser/App] --> Router[Chi Router]
    subgraph internal_api [Internal API Layer]
        Router --> Auth[Clerk Auth Middleware]
        Auth --> Handlers[Handlers]
    end
    subgraph business_logic [Business Logic Layer]
        Handlers --> Services[Services]
    end
    subgraph data_layer [Data Layer]
        Services --> Repos[Repositories]
        Repos --> Jet[Go-Jet Generated Models]
        Jet --> DB[(PostgreSQL)]
    end
    subgraph config_utils [Config & Utils]
        Config[Config Loader]
        Utils[JSON & Validation Utils]
    end
    Handlers -.-> Utils
    Router -.-> Config
```

### Key Components

- **API Layer**: Built with `go-chi/chi`, providing a lightweight and high-performance RESTful interface. Includes global middleware for **Rate Limiting** (100 req/min) and **HTTPS Enforcement** (in non-local environments).
- **Authentication**: Secured by `Clerk`, providing robust session management and user identity.
- **Data Access**: Powered by `go-jet/jet`, ensuring type-safe SQL queries and seamless database mapping.
- **Validation**: Uses `go-playground/validator` for rigorous request data validation.
- **Configuration**: Environment-based configuration using `godotenv`.

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **SQL Builder/ORM**: Go-Jet
- **Auth**: Clerk
- **Logging**: Logrus

## 🚀 Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL (or use the provided Docker setup)

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/igorschechtel/clearflow-backend.git
   cd clearflow-backend
   ```

2. **Environment Variables**:
   Create a `.env` file in the root directory and populate it with your configuration:
   ```env
   DATABASE_URL=postgres://admin:password@localhost:5432/app_db?sslmode=disable
   PORT=8080
   ENV=local
   ALLOWED_ORIGINS=http://localhost:3000
   CLERK_SECRET_KEY=your_clerk_secret_key
   ```

3. **Start the Database**:
   ```bash
   docker-compose up -d
   ```

4. **Run Migrations**:
   ```bash
   make migrate-up
   ```

5. **Start the Server**:
   Using `Air` for live-reload:
   ```bash
   air
   ```
   Or standard Go run:
   ```bash
   go run main.go
   ```

## 📜 License

This project is licensed under the MIT License.
