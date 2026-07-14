---
slug: day-3
title: Day 3
sidebar_position: 3
---
# Day 3 -- Declarative Routing and Configuration

> Study notes for building a configurable API Gateway.

## Topics Covered

-   Declarative routing
-   Hardcoded vs configuration-driven routing
-   Longest prefix matching
-   Path boundary validation
-   `strip_prefix`
-   YAML configuration
-   Go configuration structs
-   YAML decoding
-   Syntax vs semantic validation
-   Fail-fast startup
-   Ambiguous route detection
-   404 vs 503 responses
-   Table-driven testing
-   Interview notes

------------------------------------------------------------------------

## Declarative Routing

Instead of hardcoding routes in Go:

``` go
if strings.HasPrefix(path, "/users") {
    proxy(usersBackend)
}
```

Store routes in YAML:

``` yaml
routes:
  - id: users
    prefix: /users
    backend: http://localhost:9002
```

Benefits:

-   No code changes
-   No recompilation
-   Easier maintenance
-   Production standard

------------------------------------------------------------------------

## Longest Prefix Matching

Configured routes:

``` text
/
/api
/api/v1
```

Request:

``` text
/api/v1/users
```

Winner:

``` text
/api/v1
```

Always choose the **most specific** (longest) matching prefix.

``` mermaid
flowchart TD
    A["/api/v1/users"]
    A --> B["/"]
    A --> C["/api"]
    A --> D["/api/v1"]
    D --> E["Selected"]
```

------------------------------------------------------------------------

## Path Boundary Rules

Route:

``` text
/users
```

Matches:

-   `/users`
-   `/users/`
-   `/users/42`

Does NOT match:

-   `/users-old`
-   `/users123`

After the prefix, the next character must be either `/` or
end-of-string.

------------------------------------------------------------------------

## strip_prefix

Configuration:

``` yaml
prefix: /orders
strip_prefix: true
```

Request:

``` text
/orders/123/items
```

Backend receives:

``` text
/123/items
```

If the rewritten path becomes empty, rewrite it to `/`.

------------------------------------------------------------------------

## YAML Configuration

``` yaml
server:
  port: 8080

routes:
  - id: users
    prefix: /users
    backend: http://localhost:9002
    strip_prefix: false
```

Go structs:

``` go
type Config struct {
    Server ServerConfig `yaml:"server"`
    Routes []Route      `yaml:"routes"`
}
```

------------------------------------------------------------------------

## Validation

Validate after decoding:

-   Unique route IDs
-   Valid backend URLs
-   Non-empty backends
-   Positive timeouts
-   No duplicate/ambiguous routes

### Syntax vs Semantic

  Syntax            Semantic
  ----------------- ----------------
  YAML format       Business rules
  Missing colon     Duplicate IDs
  Bad indentation   Empty backend

------------------------------------------------------------------------

## Fail-Fast Startup

``` mermaid
flowchart TD
A[Read YAML] --> B[Decode]
B --> C[Validate]
C -->|Valid| D[Start Gateway]
C -->|Invalid| E[Exit]
```

------------------------------------------------------------------------

## Route Selection

``` mermaid
flowchart TD
Start --> Match
Match --> LongestPrefix
LongestPrefix --> Rewrite
Rewrite --> Proxy
```

------------------------------------------------------------------------

## Responses

**404 Not Found**

Returned when no route matches.

**503 Service Unavailable**

Returned when the matched route has no healthy backend.

------------------------------------------------------------------------

## Testing

Table-driven tests should cover:

-   `/users`
-   `/users/`
-   `/users/42`
-   `/users-old`
-   Longest-prefix selection
-   `strip_prefix`
-   Invalid configuration

------------------------------------------------------------------------

## Interview Points

-   Why declarative routing?
-   Why longest-prefix matching?
-   Why `HasPrefix` alone is unsafe?
-   Difference between syntax and semantic validation.
-   Why fail-fast startup?
