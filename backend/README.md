# LetusPass Backend

- [LetusPass Backend](#letuspass-backend)
  - [Start Development Server](#start-development-server)
  - [Update Swagger Files](#update-swagger-files)

## Start Development Server

1. Install dependencies via `go mod download`.
2. Create `.env` file. Use `.env.dist` as a template.
3. Start the database with `docker compose up -d postgres`. (Run this command at the project root.)
4. Start backend server with `air` command. ([air](https://github.com/air-verse/air) should be
installed on your system.)

You can access the backend server at `http://localhost:8080`.

## Update Swagger Files

This project uses [swaggo/swag](https://github.com/swaggo/swag) to generate OpenAPI files. To update
the OpenAPI files, run the following command:

```console
$ just swag
```

> If [just](https://github.com/casey/just) is not installed on your system, you can find
> corresponding commands in the `Justfile`.

## Run Tests

```console
$ just test
```
