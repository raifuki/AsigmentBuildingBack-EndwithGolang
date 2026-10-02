package ratelimit

// Token Bucket Lua script
// KEYS[1] = bucket key (vd: "ratelimit:ip:127.0.0.1")
// ARGV[1] = capacity (số token tối đa)
// ARGV[2] = refill_rate (token/giây)
// ARGV[3] = now (unix timestamp dạng giây, có phần thập phân)
// ARGV[4] = requested (số token cần lấy, thường = 1)
//
// Trả về: {allowed(0/1), remaining_tokens, retry_after_seconds}
const tokenBucketScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

-- Lấy state từ Redis
local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])

-- Nếu chưa tồn tại → khởi tạo bucket đầy token
if tokens == nil then
    tokens = capacity
    last_refill = now
end

-- Tính token đã nạp thêm kể từ lần refill cuối
local elapsed = math.max(0, now - last_refill)
local refill = elapsed * refill_rate
tokens = math.min(capacity, tokens + refill)

local allowed = 0
local retry_after = 0

if tokens >= requested then
    tokens = tokens - requested
    allowed = 1
else
    -- Không đủ token → tính thời gian chờ
    local deficit = requested - tokens
    retry_after = math.ceil(deficit / refill_rate)
end

-- Cập nhật lại state
redis.call('HMSET', key, 'tokens', tokens, 'last_refill', now)

-- TTL = thời gian nạp đầy bucket từ 0 (để Redis tự cleanup)
local ttl = math.ceil(capacity / refill_rate) + 1
redis.call('EXPIRE', key, ttl)

return {allowed, math.floor(tokens), retry_after}
`
