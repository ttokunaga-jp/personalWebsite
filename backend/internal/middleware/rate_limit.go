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

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter enforces per-IP request throttling to mitigate abuse.
type RateLimiter struct {
	cfg            config.SecurityConfig
	visitors       sync.Map
	ttl            time.Duration
	trustedProxies []*net.IPNet
}

// NewRateLimiter constructs the limiter and wires cleanup into the Fx lifecycle.
func NewRateLimiter(lc fx.Lifecycle, cfg *config.AppConfig) *RateLimiter {
	return newRateLimiter(lc, cfg)
}

func newRateLimiter(lc fx.Lifecycle, cfg *config.AppConfig) *RateLimiter {
	if cfg == nil {
		return &RateLimiter{}
	}
	trusted := newTrustedProxyList(cfg.Server.TrustedProxies)
	rl := &RateLimiter{
		cfg:            cfg.Security,
		ttl:            10 * time.Minute,
		trustedProxies: trusted,
	}

	if lc != nil {
		ctx, cancel := context.WithCancel(context.Background())
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				go rl.cleanupLoop(ctx)
				return nil
			},
			OnStop: func(context.Context) error {
				cancel()
				return nil
			},
		})
	}

	return rl
}

// Handler returns the gin middleware.
func (r *RateLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if r == nil || r.cfg.RateLimitRequestsPerMinute <= 0 {
			c.Next()
			return
		}

		ip := resolveClientIP(c, r.trustedProxies)
		if strings.TrimSpace(ip) == "" {
			ip = "unknown"
		}
		if r.isWhitelisted(ip) {
			c.Next()
			return
		}

		limiter := r.getVisitor(ip)
		if !limiter.Allow() {
			err := errs.New(errs.CodeInvalidInput, http.StatusTooManyRequests, "rate limit exceeded", nil)
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

func (r *RateLimiter) isWhitelisted(ip string) bool {
	for _, allowed := range r.cfg.RateLimitWhitelist {
		if ip == allowed {
			return true
		}
	}
	return false
}

func (r *RateLimiter) getVisitor(ip string) *rate.Limiter {
	value, ok := r.visitors.Load(ip)
	if ok {
		v := value.(*visitor)
		v.lastSeen = time.Now()
		return v.limiter
	}

	limit := rate.Every(time.Minute / time.Duration(r.cfg.RateLimitRequestsPerMinute))
	burst := r.cfg.RateLimitBurst
	if burst <= 0 {
		burst = 1
	}
	limiter := rate.NewLimiter(limit, burst)
	visitor := &visitor{limiter: limiter, lastSeen: time.Now()}
	r.visitors.Store(ip, visitor)
	return limiter
}

func (r *RateLimiter) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
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

func (r *RateLimiter) cleanup() {
	now := time.Now()
	r.visitors.Range(func(key, value any) bool {
		v := value.(*visitor)
		if now.Sub(v.lastSeen) > r.ttl {
			r.visitors.Delete(key)
		}
		return true
	})
}
