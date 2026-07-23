---
slug: day-8
title: Day 8
sidebar_position: 8
---

# Day-8  Distributed Redis Token Bucket Rate Limiting

## Learning Objectives

-   Understand why in-memory rate limiting fails in distributed
    gateways.
-   Use Redis Lua scripts for atomic operations.
-   Implement a distributed token bucket.
-   Return standard HTTP 429 responses with rate-limit headers.
-   Apply fail-open/fail-closed policies.
-   Expire inactive limiter keys automatically.

------------------------------------------------------------------------

## Architecture

``` mermaid
flowchart LR
    C[Client] --> G[Gateway]
    G --> RL[Redis Lua Token Bucket]
    RL --> R[(Redis)]
    RL -->|Allowed| P[Reverse Proxy]
    RL -->|Denied| E[429 Too Many Requests]
```

------------------------------------------------------------------------

## Why Lua?

Without Lua:

``` text
Read tokens
↓
Check
↓
Write tokens
```

Two gateway instances can consume the same token.

With Lua:

``` mermaid
sequenceDiagram
    participant GW as Gateway
    participant Redis

    GW->>Redis: Execute Lua Script
    Note over Redis: Entire script runs atomically
    Redis-->>GW: allowed, remaining, retry_after, reset_after
```

Redis executes one Lua script at a time, preventing race conditions.

------------------------------------------------------------------------

## Token Bucket

-   Capacity = maximum burst.
-   Refill rate = tokens added per second.
-   Each request consumes one token.

``` mermaid
flowchart TD
    A[Request] --> B{Tokens Available?}
    B -->|Yes| C[Consume Token]
    C --> D[Forward Request]
    B -->|No| E[Return 429]
```

------------------------------------------------------------------------

## Lua Script Responsibilities

-   Read bucket state
-   Compute elapsed time
-   Refill tokens
-   Decide allow/deny
-   Update Redis atomically
-   Set TTL
-   Return:
    -   allowed
    -   remaining
    -   retry_after
    -   reset_after

------------------------------------------------------------------------

## Rate Limit Headers

-   X-RateLimit-Limit
-   X-RateLimit-Remaining
-   X-RateLimit-Reset
-   Retry-After (only on 429)

------------------------------------------------------------------------

## Failure Policies

``` mermaid
flowchart TD
    A[Redis Error] --> B{Failure Policy}
    B -->|fail_open| C[Continue Request]
    B -->|fail_closed| D[503 Service Unavailable]
```

-   Public routes → fail_open
-   Protected routes → fail_closed

------------------------------------------------------------------------

## Key Expiry

Redis keys expire automatically after inactivity to avoid unbounded
growth.

------------------------------------------------------------------------

## Concurrent Verification

Spawn many goroutines using the same key.

Example:

-   Capacity = 10
-   100 concurrent requests

Expected:

-   10 Allowed
-   90 Denied

Confirms Lua atomicity.

------------------------------------------------------------------------

## End-to-End Tests

-   Normal requests
-   Bucket exhaustion (429)
-   Token refill
-   Redis unavailable
-   Public route → fail_open
-   Protected route → fail_closed
-   Concurrent requests
-   TTL cleanup

------------------------------------------------------------------------

## Key Takeaways

-   Redis enables distributed rate limiting.
-   Lua guarantees atomic updates.
-   TTL prevents stale limiter keys.
-   Standard headers improve client behavior.
-   Failure policies keep the gateway resilient.
-   Concurrency tests validate correctness.
