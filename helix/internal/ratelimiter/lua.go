package ratelimiter

const luaTokenBucket = `
local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local refill = tonumber(ARGV[2])

-- Current Redis server time
local now = redis.call("TIME")
local current = tonumber(now[1]) + tonumber(now[2]) / 1000000

-- Load bucket
local data = redis.call("HMGET", key, "tokens", "last")

local tokens = tonumber(data[1])
local last = tonumber(data[2])

if tokens == nil then
	tokens = capacity
	last = current
end

-- Refill tokens
local elapsed = current - last
tokens = math.min(capacity, tokens + elapsed * refill)

local allowed = 0
local retry_after = 0

if tokens >= 1 then
	tokens = tokens - 1
	allowed = 1
else
	retry_after = math.ceil((1 - tokens) / refill)
end

-- Remaining tokens after this request
local remaining = math.floor(tokens)

-- Time until bucket is completely full
local reset_after = math.ceil((capacity - tokens) / refill)

-- Save updated bucket
redis.call(
	"HMSET",
	key,
	"tokens", tokens,
	"last", current
)

-- Expire inactive buckets
local ttl = math.ceil((capacity / refill) * 2)
redis.call("EXPIRE", key, ttl)

return {
	allowed,
	remaining,
	retry_after,
	reset_after
}
`
