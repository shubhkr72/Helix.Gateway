package ratelimiter

const luaTokenBucket = `
local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local refill = tonumber(ARGV[2])

local now = redis.call("TIME")
local current = tonumber(now[1]) + tonumber(now[2]) / 1000000

local data = redis.call("HMGET", key, "tokens", "last")

local tokens = tonumber(data[1])
local last = tonumber(data[2])

if tokens == nil then
	tokens = capacity
	last = current
end

local elapsed = current - last

tokens = math.min(capacity, tokens + elapsed * refill)

local allowed = 0

if tokens >= 1 then
	tokens = tokens - 1
	allowed = 1
end

redis.call("HMSET",
	key,
	"tokens", tokens,
	"last", current
)

local ttl = math.ceil(capacity / refill * 2)

redis.call("EXPIRE", key, ttl)

return allowed
`
