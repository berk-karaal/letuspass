# LetusPass Frontend

- [LetusPass Frontend](#letuspass-frontend)
  - [Start Development Server](#start-development-server)
  - [Update Endpoint Functions and Types](#update-endpoint-functions-and-types)

## Start Development Server

Node version: `v20.15.0`

1. Install dependencies with `npm install`.
2. Create `.env` file. Use `.env.dist` as a template.
3. Make sure backend is up and running.
4. Start frontend server with `npm run dev`.

You can access the frontend server at `http://localhost:5173`.

## Update Endpoint Functions and Types

This project uses [orval](https://github.com/orval-labs/orval) to generate API client from OpenAPI
specification. To update the API endpoints and types, run the following command:

> Make sure that the backend server is running before running this command.

> This command will update `src/api/letuspass.ts` and `src/api/letuspass.schemas.ts` files if there
> are any changes in the OpenAPI specification. Do not manually edit these files.

```console
$ npm run orval
```

