---

# Product Service (Go)

A lightweight and clean REST API service for managing products.  
Built using Go, `database/sql`, PostgreSQL, and a layered architecture (handler → service → repository → infra).

This service is part of the **go-ecommerce monorepo**.

---

## Features

- CRUD product
- Category field (TEXT)
- Search by name or category
- Pagination (offset/limit)
- Indexing for performance
- PATCH for partial update
- Middlewares:
  - CORS
  - Rate Limit (in-memory)
  - Clean logging

---

## Tech Stack

- **Go** (pure `database/sql`)
- **PostgreSQL**
- **Layered Architecture**
- **Chi Router**
- **Shell scripts** for tooling / API test

---

## Project Structure

product-service/ ├── cmd/ │   └── main.go ├── internal/ │   ├── config/ │   ├── infra/ │   ├── middleware/ │   ├── models/ │   ├── repository/ │   ├── service/ │   └── handler/ ├── scripts/ ├── llm-context-engineering/ (ignored) └── bin/

---

## Environment Variables

Copy `.env.example` → `.env`:

DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=product_service

---

## Running

### 1. Start PostgreSQL  
Pastikan database sudah dibuat:

```sql
CREATE DATABASE product_service;

2. Build & Run

Gunakan script dev bawaan:

./scripts/dev.sh

Atau manual:

mkdir -p bin
go build -o bin/product-service ./cmd/main.go
./bin/product-service


---

API Endpoints

Products

Method	Endpoint	Description

GET	/products	List + search + pagination
GET	/products/{id}	Get detail
POST	/products	Create
PATCH	/products/{id}	Partial update
DELETE	/products/{id}	Delete



---

Scripts

scripts/api_test_manual.sh — manual curl tests

scripts/api_test_auto.sh — automated tests

scripts/dev.sh — rebuild + restart server

scripts/stop.sh — stop binary



---

Status

MVP features completed.
Next stage → adding tests, Redis caching, gRPC, Docker, deployment.


---

License

MIT

---
