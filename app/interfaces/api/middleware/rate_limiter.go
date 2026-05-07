package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"

	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
	"github.com/linzhengen/ddd-gin-admin/configs"
)

func RateLimiterMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := configs.C.RateLimiter
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	rc := configs.C.Redis
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": rc.Addr,
		},
		Password: rc.Password,
		DB:       cfg.RedisDB,
	})

	limiter := redis_rate.NewLimiter(ring)

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID := api.GetUserID(c)
		if userID != "" {
			limit := redis_rate.PerSecond(int(cfg.Count))
			res, err := limiter.Allow(c.Request.Context(), userID, limit)
			if err != nil {
				c.Next()
				return
			}
			if res.Allowed == 0 {
				h := c.Writer.Header()
				h.Set("X-RateLimit-Limit", strconv.FormatInt(cfg.Count, 10))
				h.Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
				retryAfterSec := int64(res.RetryAfter.Seconds())
				h.Set("X-RateLimit-Delay", strconv.FormatInt(retryAfterSec, 10))
				api.ResError(c, errors.ErrTooManyRequests)
				return
			}
		}

		c.Next()
	}
}
