# Golang Project

This is a Golang project that utilizes the Gin and gRPC frameworks, as well as PostgreSQL and Redis databases.

## Prerequisites

Make sure you have the following dependencies installed on your system:

- Docker
- Go
- Migrate (https://github.com/golang-migrate/migrate)
- SQLC (https://github.com/kyleconroy/sqlc)
- Evans (https://github.com/ktr0731/evans)
- Statik (https://github.com/rakyll/statik)
- Protoc (Protocol Buffers compiler) and the necessary plugins

## Getting Started

1. Clone this repository.

2. Set up the PostgreSQL database using Docker:

   ```bash
   make postgres
   ```

3. Create the required PostgreSQL database:

   ```bash
   make createdb
   ```

4. Run the database migrations:

   ```bash
   make migrateup
   ```

5. Generate the SQL code from the queries using SQLC:

   ```bash
   make sqlc
   ```

6. Start the server:

   ```bash
   make server
   ```

## Makefile Commands

- `make network`: Create a Docker network for the project.
- `make postgres`: Start the PostgreSQL database container.
- `make createdb`: Create the PostgreSQL database.
- `make dropdb`: Drop the PostgreSQL database.
- `make migrateup`: Run database migrations.
- `make migrateup1`: Run a specific database migration version.
- `make migratedown`: Rollback database migrations.
- `make migratedown1`: Rollback a specific database migration version.
- `make db_docs`: Build database documentation.
- `make db_schema`: Generate database schema SQL file.
- `make sqlc`: Generate Go code from SQL queries using SQLC.
- `make test`: Run tests.
- `make server`: Start the server.
- `make mock`: Generate mock DB for testing.
- `make proto`: Generate protocol buffer files and Swagger documentation.
- `make evans`: Start the Evans gRPC client.
- `make redis`: Start the Redis container.

Please ensure that you have the required dependencies installed before executing the above commands.
