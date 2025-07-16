# E‑Commerce with Multiple Categories

A production‑grade e‑commerce backend built in Go: type‑safe data layer, REST + gRPC streaming APIs, RBAC, Swagger docs, Paseto auth, and Docker support.

## 🔍 Features
- **PostgreSQL + sqlc**: compile‑time checked SQL, full unit tests + CI with Github Actions 
- **Gin REST API**: CRUD for products, users, orders, payments  
- **gRPC Streaming**: client‑streaming RPC for product creation with image chunking  
- **RBAC**: role‑based access control enforced in gRPC interceptors  
- **gRPC Gateway**: serve REST & gRPC on the same port  
- **Swagger UI**: auto‑generated OpenAPI docs
- **Paseto Auth**: secure, modern PASETO tokens (no JWT vulnerabilities)  
- **Dockerized**: one‑command spin‑up for dev & testing  

## 🛠️ Tech Stack
- **Language**: Go (1.24.4)  
- **Database**: PostgreSQL  
- **ORM/Queries**: [sqlc](https://github.com/kyleconroy/sqlc)  
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)  
- **RPC**: [gRPC](https://grpc.io/) + [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)  
- **Docs**: [Swagger](https://swagger.io/)  
- **Auth**: [Paseto](https://github.com/o1egl/paseto)  
- **Containerization**: Docker & Docker Compose  

## 🚀 Getting Started

### Prerequisites
- Go 1.24+  
- Docker & Docker Compose  
- PostgreSQL (or use the provided Docker setup)
