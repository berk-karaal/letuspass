# LetusPass

- [LetusPass](#letuspass)
  - [Features](#features)
  - [Demo](#demo)
  - [Try on your machine](#try-on-your-machine)
  - [Database ERD](#database-erd)
  - [Encryption System](#encryption-system)

<br>

LetusPass is a password manager application for teams or personal use. The name is
derived from "Let us pass".

**Backend stack:** Go, Gin, Gorm, PostgreSQL, swaggo/swag, zerolog <br>
**Frontend stack:** React, Mantine, React Router, Redux, React Query, Axios, Orval

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

[[ Demo video ]]

## Try on your machine

[[ Docker compose setup instructions ]]


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