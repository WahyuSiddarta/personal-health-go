package middleware

import (
	"net/http"
	"sync"

	"github.com/WahyuSiddarta/be_saham_go/config"
	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// RateLimiter holds rate limiting configuration and state
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
	enabled  bool
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter() *RateLimiter {
	cfg := config.Get().RateLimit
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		r:        rate.Limit(cfg.RequestsPerSecond),
		b:        cfg.BurstSize,
		enabled:  cfg.Enabled,
	}
}

// getVisitor returns a rate limiter for the given IP address
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.mu.Lock()
		// Check again after acquiring write lock to avoid race condition
		if _, exists := rl.visitors[ip]; !exists {
			rl.visitors[ip] = limiter
		} else {
			limiter = rl.visitors[ip]
		}
		rl.mu.Unlock()
	}

	return limiter
}

// Middleware returns an Echo middleware function for rate limiting
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip rate limiting if disabled
			if !rl.enabled {
				return next(c)
			}

			// Get client IP address
			ip := c.RealIP()
			if ip == "" {
				ip = c.Request().RemoteAddr
			}

			// Get rate limiter for this IP
			limiter := rl.getVisitor(ip)

			// Check if request is allowed
			if !limiter.Allow() {
				// Log rate limit exceeded
				Logger.Warn().Str("ip", ip).Str("path", c.Request().URL.Path).Msg("Rate limit exceeded")

				return helper.ErrorResponse(c, http.StatusTooManyRequests, "Terlalu banyak permintaan, Silakan coba lagi beberapa saat", nil)
			}

			return next(c)
		}
	}
}

// CleanupVisitors removes old visitor entries (call this periodically to prevent memory leaks)
func (rl *RateLimiter) CleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// In a production environment, you might want to implement a more sophisticated cleanup
	// strategy based on last access time or use a LRU cache
	if len(rl.visitors) > 10000 {
		// Simple strategy: clear all when too many entries
		rl.visitors = make(map[string]*rate.Limiter)
	}
}

// LogRateLimitStatus logs the current rate limiting configuration status
func LogRateLimitStatus() {
	cfg := config.Get().RateLimit

	if cfg.Enabled {
		if Logger != nil {
			Logger.Info().
				Int("requests_per_second", cfg.RequestsPerSecond).
				Int("burst_size", cfg.BurstSize).
				Msg("Rate limiting enabled for all /api routes")
		}
	} else {
		if Logger != nil {
			Logger.Info().Msg("Rate limiting disabled")
		}
	}
}
