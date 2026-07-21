---
slug: day-6
title: Day 6
sidebar_position: 6
---
# Day 6 -- JWT Authentication & Authorization at the Gateway

## Learning Objectives

-   Understand JWT authentication and authorization.
-   Validate JWTs at the API Gateway.
-   Enforce route-level access control.
-   Implement secure identity forwarding.
-   Protect against common JWT attacks (e.g., Algorithm Confusion).

------------------------------------------------------------------------

# Authentication vs Authorization

Authentication answers **"Who are you?"**

If authentication fails, the server returns:
-   **401 Unauthorized**

Common failure reasons:
-   Missing token
-   Malformed token
-   Invalid signature
-   Expired token
-   Wrong issuer or audience

Authorization answers **"Are you allowed to do this?"**

If authorization fails, the server returns:
-   **403 Forbidden**

Example: A user has a valid JWT but their role is `user`, while the route requires `admin`.

``` mermaid
flowchart LR
    A[Request] --> B{Authenticated?}
    B -->|No| C[401 Unauthorized]
    B -->|Yes| D{Authorized?}
    D -->|No| E[403 Forbidden]
    D -->|Yes| F[Allow Request]
```

------------------------------------------------------------------------

# JWT Structure

A JWT consists of three parts:
`Header.Payload.Signature`

``` mermaid
flowchart LR
    A[Header] --> D[JWT]
    B[Payload / Claims] --> D
    C[Signature] --> D
```

## Header
Contains the metadata:
-   `alg` -- Algorithm (e.g., RS256)
-   `typ` -- Token type (JWT)

## Payload (Claims)
Contains the identity and metadata:
-   `exp` -- Expiration time
-   `nbf` -- Not valid before
-   `iss` -- Issuer
-   `aud` -- Audience
-   `role` -- Custom claim for authorization

## Signature
Generated using the private key to ensure the token hasn't been tampered with. It is verified using the public key.

------------------------------------------------------------------------

# RS256 Cryptography

RS256 uses asymmetric encryption.

-   **Private Key**: Signs the JWT. Only the **Auth Service** owns it.
-   **Public Key**: Verifies the JWT. The **Gateway** and backend services use it.

``` mermaid
flowchart LR
    A[Auth Service] -->|Private Key: Sign| B[JWT]
    B --> C[API Gateway]
    D[Public Key] -->|Verify| C
```

Important: If the public key is leaked, an attacker still **cannot** create valid JWTs because they lack the private key.

------------------------------------------------------------------------

# Gateway Authentication Flow

The Gateway intercepts every request to validate the identity of the caller.

``` mermaid
flowchart TD
    A[Client Request] --> B[Extract Authorization Header]
    B --> C{Bearer Token?}
    C -->|No| D[401 Unauthorized]
    C -->|Yes| E[Parse JWT]
    E --> F[Verify Signature with Public Key]
    F --> G[Validate Algorithm == RS256]
    G --> H[Validate Claims: exp, nbf, iss, aud]
    H --> I[Create Principal Object]
    I --> J[Inject Principal into Context]
    J --> K[Forward to Backend]
```

------------------------------------------------------------------------

# Bearer Token

Standard request format:
``` text
Authorization: Bearer <jwt-token>
```

The Gateway:
1.  Extracts the `Authorization` header.
2.  Removes the `Bearer ` prefix.
3.  Passes the remaining token to the JWT validator.

------------------------------------------------------------------------

# JWT Claims Validation

The Gateway must reject the token if any of these checks fail:
-   **Signature**: Must be cryptographically valid.
-   **Algorithm**: Must match the expected algorithm (RS256).
-   **Expiry (`exp`)**: Current time must be before expiry.
-   **Not Before (`nbf`)**: Current time must be after the "not before" time.
-   **Issuer (`iss`)**: Must match the trusted Auth Service issuer.
-   **Audience (`aud`)**: Must match the Gateway's expected audience.

------------------------------------------------------------------------

# Algorithm Confusion Attack

**The Danger:**
An attacker changes the header from `alg = RS256` (asymmetric) to `alg = HS256` (symmetric). If the server blindly trusts the header, it might try to verify the RS256 public key as an HS256 secret key, potentially allowing forged tokens.

``` mermaid
sequenceDiagram
    Attacker->>Gateway: JWT (alg=HS256, signed with Public Key)
    Note over Gateway: Trusts Header 'alg'
    Gateway->>Gateway: Verify using Public Key as Secret
    Gateway-->>Attacker: 200 OK (Forged token accepted)
```

**The Solution:**
Always **pin the expected algorithm** (RS256). Never trust the `alg` field in the JWT header.

------------------------------------------------------------------------

# Audience (aud) and Issuer (iss)

-   **Audience (`aud`)**: Indicates who the token is intended for. If a token is issued for `mobile-app` but sent to `helix-gateway`, the Gateway should return **401 Unauthorized**.
-   **Issuer (`iss`)**: Identifies the authority that created the token. The Gateway validates that the issuer matches the configured trusted issuer (e.g., `auth-service`).

------------------------------------------------------------------------

# Principal

Instead of passing raw JWT strings throughout the system, the Gateway converts the validated claims into a typed **Principal** object.

``` go
type Principal struct {
    UserID string
    Email  string
    Role   string
}
```

The Principal is stored in the request context and used for authorization decisions.

------------------------------------------------------------------------

# Route Policies

Routes are categorized by their access requirements:

## Public
No JWT required.
-   `/login`
-   `/register`
-   `/healthz`

## Authenticated
Requires a valid JWT.
-   `/users/me`
-   `/orders`

## Role Protected
Requires a valid JWT AND a specific role.
-   `DELETE /admin/users` $\rightarrow$ Requires `Role = admin`

------------------------------------------------------------------------

# Identity Header Forwarding

The Gateway acts as a trust boundary. To prevent "Identity Spoofing", the Gateway must:
1.  **Strip** any client-supplied identity headers (e.g., `X-User-ID`).
2.  **Inject** trusted headers derived from the validated JWT.

``` mermaid
flowchart LR
    Client -->|X-User-ID: attacker| Gateway
    Gateway -->|Strip Untrusted Headers| Gateway
    Gateway -->|Inject X-User-ID: validated-id| Backend
```

Example forwarded headers:
-   `X-User-ID`
-   `X-User-Email`
-   `X-User-Role`

------------------------------------------------------------------------

# Key Rotation and Clock Skew

-   **Key Rotation**: The ability to replace keys without downtime. The system should support multiple valid public keys during a rotation period.
-   **Clock Skew**: Since servers' clocks aren't perfectly synced, the Gateway allows a small grace period (e.g., 1-2 minutes) when checking `exp` and `nbf`.

------------------------------------------------------------------------

# Implementation Details

**Library:** `github.com/golang-jwt/jwt/v5`

**Core Components:**
-   JWT Validator
-   Authentication Middleware
-   Route Policy Engine
-   Principal Context Injector
-   Identity Forwarding Logic

------------------------------------------------------------------------

# Verification Checklist

-   ✅ Missing token returns **401 Unauthorized**
-   ✅ Malformed token returns **401 Unauthorized**
-   ✅ Wrong signing algorithm returns **401 Unauthorized**
-   ✅ Wrong issuer returns **401 Unauthorized**
-   ✅ Wrong audience returns **401 Unauthorized**
-   ✅ Expired token returns **401 Unauthorized**
-   ✅ Valid token allows request
-   ✅ Forbidden role returns **403 Forbidden**

------------------------------------------------------------------------

# Scope Guard

Only the **API Gateway** performs JWT validation and route policy enforcement. Backend services should trust the identity headers forwarded by the Gateway and must not perform their own JWT validation.

------------------------------------------------------------------------

# Commit Message
`feat: enforce jwt route policies`

------------------------------------------------------------------------

# Key Takeaways

-   **Authentication** verifies identity (401 if fail).
-   **Authorization** verifies permissions (403 if fail).
-   **RS256** ensures only the Auth Service can issue tokens.
-   **Algorithm Pinning** prevents confusion attacks.
-   **Audience/Issuer** validation prevents token misuse across services.
-   **Principal** objects decouple identity from transport.
-   **Identity Forwarding** ensures backends only trust the Gateway.
