---
slug: day-5
title: Day 5
sidebar_position: 5
---
# Day 5 -- Auth Service and PostgreSQL

## Learning Objectives

-   Authentication vs Authorization
-   Password hashing
-   SQL migrations
-   Connection pools
-   Context-aware queries
-   JWT structure (Header, Claims, Signature)
-   JWT expiry, issuer, and audience

------------------------------------------------------------------------

# Authentication vs Authorization

Authentication answers **"Who are you?"**

``` text
POST /login
Email: john@example.com
Password: ********
```

If credentials are valid, the server returns a JWT.

Authorization answers **"What are you allowed to do?"**

``` text
DELETE /users/1
```

The server checks the user's role from the JWT.

``` mermaid
flowchart LR
A[User Login] --> B[Authentication]
B -->|JWT Issued| C[Authenticated User]
C --> D[Authorization]
D -->|role=admin| E[Allow]
D -->|role=user| F[403 Forbidden]
```

------------------------------------------------------------------------

# Password Hashing

Never store plaintext passwords.

``` go
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

Verify using:

``` go
bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
```

------------------------------------------------------------------------

# PostgreSQL

Store only authentication data.

Users table:

  Column          Type
  --------------- ----------------
  id              UUID/BIGSERIAL
  email           TEXT UNIQUE
  password_hash   TEXT
  role            TEXT
  created_at      TIMESTAMP
  updated_at      TIMESTAMP

------------------------------------------------------------------------

# SQL Migrations

    migrations/
    ├── 001_create_users.up.sql
    └── 001_create_users.down.sql

------------------------------------------------------------------------

# Connection Pool

``` go
db.SetMaxOpenConns(20)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)
```

``` mermaid
flowchart TD
A[Application] --> B[Connection Pool]
B --> C[(PostgreSQL)]
B --> D[(PostgreSQL)]
B --> E[(PostgreSQL)]
```

------------------------------------------------------------------------

# Context-aware Queries

Always use:

``` go
db.QueryContext(ctx, ...)
```

``` mermaid
sequenceDiagram
Client->>Server: Request
Server->>Database: QueryContext(ctx)
Client-->>Server: Disconnect
Server->>Database: Context Cancelled
Database-->>Server: Query Cancelled
```

------------------------------------------------------------------------

# JWT

JWT format:

    Header.Payload.Signature

``` mermaid
flowchart LR
A[Header] --> D[JWT]
B[Payload / Claims] --> D
C[Signature] --> D
```

## Header

``` json
{
  "alg":"RS256",
  "typ":"JWT"
}
```

## Claims

-   sub
-   email
-   role
-   iat
-   exp
-   iss
-   aud

Never place passwords inside JWT claims.

------------------------------------------------------------------------

# JWT Signing

``` mermaid
flowchart LR
A[Private Key<br/>Auth Service] -->|Sign| B[JWT]
B --> C[Gateway]
D[Public Key] -->|Verify| C
```

-   Private key remains only inside the auth service.
-   Gateway stores only the public key.

------------------------------------------------------------------------

# Build Tasks

## 1. Start PostgreSQL

``` bash
docker compose up -d postgres
```

## 2. Migration

Create:

    001_create_users.up.sql
    001_create_users.down.sql

## 3. Register Endpoint

    POST /register

Flow:

``` mermaid
flowchart TD
A[Validate Input]
-->B[Hash Password]
-->C[Insert User]
-->D[201 Created]
```

------------------------------------------------------------------------

## 4. Login Endpoint

    POST /login

``` mermaid
flowchart TD
A[Find Email]
-->B[Compare bcrypt Hash]
-->C[Generate JWT]
-->D[Return Access Token]
```

Response

``` json
{
  "access_token":"<jwt>"
}
```

------------------------------------------------------------------------

## 5. Seed Demo User

    make seed

Example:

-   Email: demo@example.com
-   Password: password123

------------------------------------------------------------------------

## 6. Project Structure

``` text
auth-service/
├── cmd/
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── jwt/
│   └── database/
├── migrations/
├── scripts/
├── keys/
├── Dockerfile
└── go.mod
```

------------------------------------------------------------------------

# Verification Checklist

-   ✅ Duplicate email returns **409 Conflict**
-   ✅ Invalid password returns **401 Unauthorized**
-   ✅ Successful login returns JWT
-   ✅ Passwords are hashed
-   ✅ Passwords are never logged
-   ✅ JWT expires
-   ✅ Private key is never committed
-   ✅ Database password is stored in environment variables or secrets

------------------------------------------------------------------------

# Scope Guard

Only the **Auth Service** uses PostgreSQL.

Do **not** store:

-   Gateway access logs
-   Gateway routing configuration

inside PostgreSQL during this phase.

