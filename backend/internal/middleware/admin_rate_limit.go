package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"golang.org/x/time/rate"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
)

type adminVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// AdminRateLimiter enforces per IP+session throttling for administrative routes.
type AdminRateLimiter struct {
	cfg            config.AdminRateLimitConfig
	visitors       sync.Map
	ttl            time.Duration
	trustedProxies []*net.IPNet
	enabled        bool
}

// NewAdminRateLimiter constructs the limiter and wires cleanup into the Fx lifecycle.
func NewAdminRateLimiter(lc fx.Lifecycle, cfg *config.AppConfig) *AdminRateLimiter {
	if cfg == nil {
		return &AdminRateLimiter{}
	}

	security := cfg.Security
	trusted := newTrustedProxyList(cfg.Server.TrustedProxies)
	limiter := &AdminRateLimiter{
		cfg:            security.AdminRateLimit,
		ttl:            10 * time.Minute,
		trustedProxies: trusted,
		enabled:        security.AdminRateLimit.Enabled && security.AdminRateLimit.RequestsPerMinute > 0,
	}

	if !limiter.enabled {
		return limiter
	}

	if lc != nil {
		ctx, cancel := context.WithCancel(context.Background())
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				go limiter.cleanupLoop(ctx)
				return nil
			},
			OnStop: func(context.Context) error {
				cancel()
				return nil
			},
		})
	}

	return limiter
}

// Handler returns a gin middleware applying rate limiting to admin endpoints.
func (r *AdminRateLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if r == nil || !r.enabled {
			c.Next()
			return
		}

		key := r.buildKey(c)
		limiter := r.getVisitor(key)
		if !limiter.Allow() {
			err := errs.New(errs.CodeInvalidInput, http.StatusTooManyRequests, "admin rate limit exceeded", nil)
			payload := gin.H{
				"error":   err.Code,
				"message": err.Message,
			}
			if requestID := c.Writer.Header().Get("X-Request-ID"); requestID != "" {
				payload["request_id"] = requestID
			}
			c.JSON(err.Status, payload)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (r *AdminRateLimiter) buildKey(c *gin.Context) string {
	ip := resolveClientIP(c, r.trustedProxies)
	if strings.TrimSpace(ip) == "" {
		ip = "unknown"
	}

	sessionID := "anonymous"
	if session, ok := GetSessionFromContext(c); ok && session != nil {
		if trimmed := strings.TrimSpace(session.ID); trimmed != "" {
			sessionID = trimmed
		}
	}

	return ip + "|" + sessionID
}

func (r *AdminRateLimiter) getVisitor(key string) *rate.Limiter {
	if key == "" {
		key = "unknown|anonymous"
	}

	if value, ok := r.visitors.Load(key); ok {
		entry := value.(*adminVisitor)
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limit := rate.Every(time.Minute / time.Duration(r.cfg.RequestsPerMinute))
	burst := r.cfg.Burst
	if burst <= 0 {
		burst = 1
	}

	limiter := rate.NewLimiter(limit, burst)
	r.visitors.Store(key, &adminVisitor{
		limiter:  limiter,
		lastSeen: time.Now(),
	})
	return limiter
}

func (r *AdminRateLimiter) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.cleanup()
		}
	}
}

func (r *AdminRateLimiter) cleanup() {
	now := time.Now()
	r.visitors.Range(func(key, value any) bool {
		entry := value.(*adminVisitor)
		if now.Sub(entry.lastSeen) > r.ttl {
			r.visitors.Delete(key)
		}
		return true
	})
}
