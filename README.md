# LetusPass

- [LetusPass](#letuspass)
  - [Features](#features)
  - [Demo](#demo)
  - [Try on your local machine](#try-on-your-local-machine)
  - [Database ERD](#database-erd)
  - [Encryption System](#encryption-system)

<br>

LetusPass is a password manager application for teams or personal use. The name is
derived from "Let us pass".

**Backend stack:** [Go](https://go.dev/), [Gin](https://github.com/gin-gonic/gin),
[Gorm](https://github.com/go-gorm/gorm), [PostgreSQL](https://www.postgresql.org/),
[swaggo/swag](https://github.com/swaggo/swag), [zerolog](https://github.com/rs/zerolog) <br>
**Frontend stack:** [React](https://react.dev/), [Mantine](https://mantine.dev/), [React
Router](https://reactrouter.com/), [Redux](https://redux.js.org/), [React
Query](https://tanstack.com/query/latest/docs/framework/react/overview),
[Axios](https://github.com/axios/axios), [Orval](https://github.com/orval-labs/orval)

## Features
- End-to-end encryption
  - Every encryption and decryption is done on the client side. The server never
  sees saved credentials in unencrypted form.
- Shareble/Collabrative vaults
- Vault audit logs
- Permission management
- Mobile friendly UI
- OpenAPI/Swagger documentation
- Structured json logging

## Demo

Watch demo video on YouTube:

[![Demo video](https://img.youtube.com/vi/k7Nc9b3VfMY/0.jpg)](https://youtu.be/k7Nc9b3VfMY)

## Try on your local machine 

Simply just run `docker compose up` in the root directory.

The application will be available at `http://localhost:3000`.

You can access the OpenAPI documentation at `http://localhost:8080/swagger/index.html`.

You can stop the application with `Ctrl+C` and remove the containers using `docker compose down`
command.

## Database ERD
<div align="center">
  <img src="./assets/complete_erd.png" width="60%" height="60%">
</div>

## Encryption System

Every user has a private and a public key which derived from the user's password
and a salt. Public keys are uploaded to the server where other users can retrieve.

Every vault has a vault key which is used to encrypt and decrypt the vault items.
Vault keys are stored in encrypted form in the database. Vault keys need to be
decrypted before use. Decryption is done on the client side as well.

To understand how encryption and decryption work in the application, you can study below
graphs and flows. They only includes parts related to encryption/decryption
(permission management, audit logging, error handling, etc. are not included).


<details>
  <summary>Cryptography Related ERD</summary>

  <div align="center">
    <img src="./assets/exports/encryption-Crypto ERD.drawio.svg" width="75%">
  </div>

</details>

<details>
  <summary>Register Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Register.drawio.svg" width="50%">
  </div>
</details>

<details>
  <summary>Login Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Login.drawio.svg" width="50%">
  </div>
</details>

<details>
  <summary>Create Vault Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Create Vault.drawio.svg" width="75%">
  </div>
</details>


<details>
  <summary>Get Raw Vault Key Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Get Raw Vault Key.drawio.svg" width="75%">
  </div>
</details>

<details>
  <summary>Add User to Vault Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Add User to Vault.drawio.svg" width="75%">
  </div>
</details>

<details>
  <summary>Add New Vault Item to Vault Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Add New Vault Item To Vault.drawio.svg" width="60%">
  </div>
</details>


<details>
  <summary>Retrieve a Vault Item Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Retrieve a Vault Item.drawio.svg" width="60%">
  </div>
</details>


<details>
  <summary>Update a Vault Item Flow</summary>
  
  <div align="center">
    <img src="./assets/exports/encryption-Update Vault Item.drawio.svg" width="50%">
  </div>
</details>



---

Backend documentation at [backend/](./backend/README.md) <br>
Frontend documentation at [frontend/](./frontend/README.md)