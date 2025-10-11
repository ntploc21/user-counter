package app

import (
	"fmt"
	"locntp-user-counter/config"
	"locntp-user-counter/internal/routes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
)

type Server struct {
	app *gin.Engine
}

func NewServer() *Server {
	// Set the release mode
	if config.GetAppConfig().IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		app: gin.Default(),
	}

	return s
}

// securityMiddleware sets the security headers and policies
func (s *Server) SecurityMiddleware() {
	// Set the security headers
	// X-Content-Type-Options: nosniff - Prevents browsers from MIME-sniffing a response away from the declared content-type
	// X-Frame-Options: DENY - Prevents clickjacking attacks
	// X-XSS-Protection: 1; mode=block - Prevents reflected XSS attacks
	// Strict-Transport-Security: max-age=31536000; includeSubDomains - Forces the browser to use HTTPS for the next year
	s.app.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Additional security headers
		c.Header(
			"Content-Security-Policy",
			"default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;",
		)
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		c.Next()
	})

	// Set the trusted proxies (default: Disable)
	err := s.app.SetTrustedProxies(nil)
	if err != nil {
		return
	}
}

func (s *Server) CorsMiddleware() {
	s.app.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		c.Header(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Accept, Authorization, X-Requested-With",
		)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// rateLimitMiddleware implements rate limiting
func (s *Server) RateLimitMiddleware(cacheLimiter *redis_rate.Limiter) {

	// If Redis rate limiter is provided, use it
	if cacheLimiter != nil {
		var requestRate = 10000
		s.app.Use(func(c *gin.Context) {
			context := c.Request.Context()
			res, err := cacheLimiter.Allow(context, c.ClientIP(), redis_rate.PerSecond(requestRate))
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if res.Allowed == 0 {
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}

			c.Next()
		})
		return
	}
	var requestRate int64 = 10000
	// Fallback to using in-memory rate limiter
	// Create a rate limiter: 10000 requests per second

	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  requestRate,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	s.app.Use(func(c *gin.Context) {
		context := c.Request.Context()
		limiterCtx, err := instance.Get(context, c.ClientIP())
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limiterCtx.Reached {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	})
}

// routeMiddleware sets up the routes and the corresponding handlers
func (s *Server) RouteHandler(db *gorm.DB, cache *redis.Client) {
	s.app = routes.SetupRouter(db, cache, s.app)
}

func (s *Server) StartServer() {
	host := config.GetAppConfig().Server.Host
	port := config.GetAppConfig().Server.Port

	// Start the server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      s.app,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("Server failed to start")
	}
}
