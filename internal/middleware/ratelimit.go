package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimit(
	cache *redis.Client,
	limit int,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			ip, _, err := net.SplitHostPort(
				r.RemoteAddr,
			)
			if err != nil {
				http.Error(
					w,
					"invalid ip",
					http.StatusBadRequest,
				)
				return
			}

			key := "rate_limit:" + ip

			ctx := context.Background()

			count, err := cache.Incr(
				ctx,
				key,
			).Result()

			if err != nil {
				http.Error(
					w,
					"internal error",
					http.StatusInternalServerError,
				)
				return
			}

			// первый запрос → ставим TTL
			if count == 1 {
				cache.Expire(
					ctx,
					key,
					time.Minute,
				)
			}

			// превышен лимит
			if count > int64(limit) {

				w.Header().Set(
					"Retry-After",
					strconv.Itoa(60),
				)

				http.Error(
					w,
					"too many requests",
					http.StatusTooManyRequests,
				)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
