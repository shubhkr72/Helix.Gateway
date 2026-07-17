---
slug: day-4
title: Day 4
sidebar_position: 4
---
# Day 4 - Middleware Fundamentals (Helix Gateway)

## Objectives

Learned:

-   Decorator-style middleware
-   `context.Context` and request-scoped values
-   Wrapping `http.ResponseWriter`
-   Middleware ordering
-   Structured JSON logging
-   High-cardinality vs low-cardinality fields
-   Request-ID middleware
-   Panic recovery middleware

------------------------------------------------------------------------

## 1. What is Middleware?

Middleware wraps an existing `http.Handler` to execute logic before
and/or after the next handler.

``` go
func Middleware(next http.Handler) http.Handler
```

-   Takes a handler as input.
-   Returns a new handler.
-   Enables composition.

Example chain:

``` text
Client
  │
  ▼
RequestID
  │
  ▼
Logging
  │
  ▼
Recovery
  │
  ▼
CORS
  │
  ▼
Router
  │
  ▼
Reverse Proxy
  │
  ▼
Backend
```

Execution:

-   Code before `next.ServeHTTP()` runs before the handler.
-   Code after `next.ServeHTTP()` runs after the handler.

------------------------------------------------------------------------

## 2. Decorator Pattern

Each middleware decorates the next handler.

``` go
handler := Logging(
    Recovery(
        RequestID(
            router,
        ),
    ),
)
```

Benefits:

-   Separation of concerns
-   Reusable logic
-   Clean handlers

------------------------------------------------------------------------

## 3. Context Values

Request-specific values belong in `context.Context`.

Never use global variables.

Create a value:

``` go
ctx := context.WithValue(r.Context(), requestIDKey, id)
```

Attach it:

``` go
r = r.WithContext(ctx)
next.ServeHTTP(w, r)
```

Read later:

``` go
id, _ := r.Context().Value(requestIDKey).(string)
```

Uses:

-   Request ID
-   Authentication info
-   Deadlines
-   Cancellation

------------------------------------------------------------------------

## 4. ResponseWriter Wrapping

`http.ResponseWriter` does not expose:

-   Status code
-   Bytes written

Wrap it.

``` go
type loggingResponseWriter struct {
    http.ResponseWriter
    status int
    bytes  int
}
```

Override:

-   `WriteHeader()`
-   `Write()`

Remember:

If the handler only calls:

``` go
w.Write(...)
```

Go automatically sends **200 OK**.

Initialize status to `http.StatusOK` or detect the first `Write()`.

------------------------------------------------------------------------

## 5. Middleware Ordering

Recommended order:

``` text
RequestID
    ↓
Logging
    ↓
Recovery
    ↓
CORS
    ↓
Router
```

Why?

-   RequestID available everywhere.
-   Recovery converts panics into HTTP 500.
-   Logging records the final response.
-   CORS adds response headers.

------------------------------------------------------------------------

## 6. Panic Recovery

Without custom recovery:

-   Go's `net/http` keeps the server alive.
-   Current request fails.
-   Default panic logging.
-   No structured JSON response.

Recovery middleware should:

-   `defer`
-   `recover()`
-   Log panic
-   Log stack trace (`runtime/debug.Stack()`)
-   Return JSON 500
-   Preserve Request-ID

Important:

`recover()` only works inside a deferred function during panic
unwinding.

------------------------------------------------------------------------

## 7. Structured Logging

Preferred over plain text.

Example:

``` json
{
  "timestamp":"2026-07-17T10:00:00Z",
  "request_id":"abc123",
  "method":"GET",
  "route_id":"users",
  "status":200,
  "bytes":512,
  "duration_ms":18,
  "backend":"http://localhost:9002",
  "client_ip":"192.168.1.15"
}
```

Benefits:

-   Easy filtering
-   Machine readable
-   Works well with Grafana, Loki, Elasticsearch, Datadog, Splunk

------------------------------------------------------------------------

## 8. Cardinality

## Low Cardinality (Good metric labels)

-   status_code
-   method
-   route
-   backend

## High Cardinality (Do NOT use as metric labels)

-   request_id
-   user_id
-   email
-   session_id
-   jwt
-   ip_address

Golden Rule:

> High-cardinality data belongs in logs, not metrics.

------------------------------------------------------------------------

## 9. Request-ID Middleware

Responsibilities:

1.  Read `X-Request-ID`
2.  Validate it
3.  Generate one if missing/invalid
4.  Store in context
5.  Echo in response header

Flow:

``` text
Request
   │
Read X-Request-ID
   │
Missing/Invalid?
   ├── Yes → Generate
   └── No  → Reuse
   │
Store in Context
   │
Echo Response Header
   │
Next Middleware
```

Benefits:

-   Correlate logs
-   Trace requests
-   Distributed tracing

------------------------------------------------------------------------

## 10. Recovery Middleware

Responsibilities:

-   Catch panics
-   Prevent broken responses
-   Return consistent JSON 500
-   Log stack traces
-   Keep server healthy

Never expose panic details to clients.

------------------------------------------------------------------------

## Key Takeaways

-   Middleware is a decorator around `http.Handler`.
-   Use `context.Context` for request-scoped data.
-   Wrap `ResponseWriter` to capture status and bytes.
-   Logging should surround Recovery.
-   Use JSON logs.
-   Keep high-cardinality fields out of metrics.
-   Recovery uses `defer` + `recover()`.
-   Echo Request-ID in responses.

