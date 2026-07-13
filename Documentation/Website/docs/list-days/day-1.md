---
slug: day-1
title: Day 1
---
# Go HTTP Fundamentals Study Notes

## 1. Request Lifecycle

``` mermaid
flowchart LR
A[Client/Browser] --> B[http.ListenAndServe]
B --> C[ServeMux]
C --> D[Handler]
D --> E[Read Request (r)]
D --> F[Write Response (w)]
F --> G[Client]
```

### Flow

1.  `http.ListenAndServe()` continuously listens for HTTP requests.
2.  It forwards each request to `ServeMux`.
3.  `ServeMux` matches the URL path.
4.  The matching handler executes business logic.
5.  The handler reads data from `*http.Request`.
6.  The handler writes headers, status code and body using
    `http.ResponseWriter`.

------------------------------------------------------------------------

## 2. http.Handler

``` go
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
```

Any type implementing `ServeHTTP` is an HTTP handler.

Example:

``` go
type HealthHandler struct{}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}
```

------------------------------------------------------------------------

## 3. http.HandlerFunc

Converts a normal function into an `http.Handler`.

``` go
func Health(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}
```

Registration:

``` go
mux.HandleFunc("/healthz", Health)
```

------------------------------------------------------------------------

## 4. ServeMux

Go's built-in router.

``` go
mux := http.NewServeMux()

mux.HandleFunc("/healthz", Health)
mux.HandleFunc("/readyz", Ready)
```

Routes requests based on URL path.

------------------------------------------------------------------------

## 5. Request Anatomy

``` go
func Handler(w http.ResponseWriter, r *http.Request) {

    r.Method
    r.URL.Path
    r.URL.Query()
    r.Header
    r.Body
    r.Cookie("session")
    r.Context()

}
```

  Field         Purpose
  ------------- -----------------------
  Method        GET, POST...
  URL.Path      Route
  URL.Query()   Query parameters
  Header        Request headers
  Body          Request payload
  Cookie        Read cookies
  Context       Cancellation, timeout

------------------------------------------------------------------------

## 6. Response Anatomy

``` go
w.Header().Set("Content-Type","application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status":"ok"}`))
```

Always:

1.  Set headers
2.  Set status code
3.  Write body

------------------------------------------------------------------------

## 7. HTTP Methods

  Method    Purpose                    Idempotent
  --------- -------------------------- ------------
  GET       Read                       ✅
  POST      Create                     ❌
  PUT       Replace entire resource    ✅
  PATCH     Partial update             Usually ❌
  DELETE    Delete                     ✅
  HEAD      Same as GET without body   ✅
  OPTIONS   Supported methods          ✅

Examples:

``` http
GET /users/10
POST /users
PUT /users/10
PATCH /users/10
DELETE /users/10
```

------------------------------------------------------------------------

## 8. Headers

Reading:

``` go
token := r.Header.Get("Authorization")
```

Writing:

``` go
w.Header().Set("Content-Type","application/json")
```

Common headers

  Header          Purpose
  --------------- -------------------
  Content-Type    Body format
  Accept          Expected response
  Authorization   Authentication
  User-Agent      Client info
  Cookie          Session
  Host            Requested host

------------------------------------------------------------------------

## 9. Query Parameters

    /search?q=golang&page=2

``` go
q := r.URL.Query().Get("q")
page := r.URL.Query().Get("page")
```

Path

    /search

Query

    q=golang&page=2

------------------------------------------------------------------------

## 10. Status Codes

### 2xx Success

-   200 OK
-   201 Created
-   202 Accepted
-   204 No Content

### 3xx Redirection

-   301 Moved Permanently
-   302 Found
-   304 Not Modified

### 4xx Client Errors

-   400 Bad Request
-   401 Unauthorized
-   403 Forbidden
-   404 Not Found
-   405 Method Not Allowed
-   409 Conflict
-   429 Too Many Requests

### 5xx Server Errors

-   500 Internal Server Error
-   502 Bad Gateway
-   503 Service Unavailable
-   504 Gateway Timeout

Go constants:

``` go
http.StatusOK
http.StatusCreated
http.StatusBadRequest
http.StatusNotFound
http.StatusInternalServerError
```

------------------------------------------------------------------------

## 11. Cookies

Read

``` go
cookie, err := r.Cookie("session")
```

Write

``` go
http.SetCookie(w,&http.Cookie{
 Name:"session",
 Value:"abc123",
})
```

Purpose

-   Login sessions
-   Preferences
-   Authentication

------------------------------------------------------------------------

## 12. Context

``` go
ctx := r.Context()
```

Uses

-   Cancellation
-   Timeouts
-   Deadlines
-   Request-scoped values

``` mermaid
flowchart LR
Client --> Gateway
Gateway --> BackendA
Gateway --> BackendB
Client -.disconnects.-> Gateway
Gateway --> CancelContext
CancelContext --> BackendA
CancelContext --> BackendB
```

------------------------------------------------------------------------

## 13. Middleware

Middleware wraps another handler.

``` mermaid
flowchart LR
Client --> Logging
Logging --> Auth
Auth --> Handler
```

Benefits

-   Logging
-   Authentication
-   Recovery
-   CORS
-   Rate limiting

Why `http.Handler` enables middleware:

Every middleware accepts and returns an `http.Handler`, allowing
handlers to be chained together.

------------------------------------------------------------------------

## 14. Go Modules

Initialize

``` bash
go mod init github.com/username/project
```

Useful commands

``` bash
go mod tidy
go mod download
go list -m all
```

------------------------------------------------------------------------

## 15. Table-Driven Tests

``` go
tests := []struct{
    name string
    want int
}{
    {"case1",1},
    {"case2",2},
}
```

Advantages

-   Less duplication
-   Easy to add cases
-   Standard Go style

------------------------------------------------------------------------

## 16. httptest

``` go
req := httptest.NewRequest(...)
rr := httptest.NewRecorder()

handler.ServeHTTP(rr, req)
```

Assertions

``` go
rr.Code
rr.Body.String()
rr.Header()
```

------------------------------------------------------------------------

## 17. Complete Mental Model

``` mermaid
flowchart TD
Client --> ListenAndServe
ListenAndServe --> ServeMux
ServeMux --> Middleware
Middleware --> Handler
Handler --> Request
Handler --> Response
Response --> Client
```

------------------------------------------------------------------------

## Cheat Sheet

  Item              Syntax
  ----------------- ---------------------------------
  Method            `r.Method`
  Path              `r.URL.Path`
  Query             `r.URL.Query().Get("key")`
  Header            `r.Header.Get("Authorization")`
  Cookie            `r.Cookie("session")`
  Context           `r.Context()`
  Response Header   `w.Header().Set()`
  Status            `w.WriteHeader()`
  Body              `w.Write()`
  Router            `http.NewServeMux()`
  Server            `http.ListenAndServe()`
  Test Request      `httptest.NewRequest()`
  Test Response     `httptest.NewRecorder()`



