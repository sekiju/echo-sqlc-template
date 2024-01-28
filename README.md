# echo-sqlc-template

## Stack
* Golang `1.21`
* Postgres `14`
* Self-written storage manager

## Modules
* [echo](https://echo.labstack.com/): web framework
* [pgx](https://github.com/jackc/pgx): PostgreSQL driver
* [sqlc](https://sqlc.dev/): type-safe SQL queries generator
* [tern](https://github.com/jackc/tern): migration tool for PostgreSQL
* [koanf](https://github.com/knadh/koanf): configuration management library
* [gomail](https://github.com/go-gomail/gomail): package for send emails

## Docker
```yaml
version: '3'
services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8000:8000"
    expose:
      - 8000
    environment:
      APPLICATION_ADDRESS: :8000
      APPLICATION_DEVELOPMENT: false
      DATABASE_URI: postgresql://postgres:postgres@postgres:5432/database
      STORAGE_HOST: http://localhost:8888
      STORAGE_DAEMON: http://localhost:8888
      SMTP_FROM: 
      SMTP_HOST:
      SMTP_USERNAME: 
      SMTP_PASSWORD: 
      SMTP_PORT: 
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:14.10-alpine3.18
    environment:
      POSTGRES_DB: database
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready -U postgres" ]
      interval: 10s
      timeout: 3s
      retries: 3
    expose:
      - 5432
    volumes:
      - postgres-data:/var/lib/postgresql/data

volumes:
  postgres-data:
```