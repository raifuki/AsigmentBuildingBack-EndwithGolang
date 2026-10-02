package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/raifuki/task-management/internal/cache"
	"github.com/raifuki/task-management/internal/config"
	"github.com/raifuki/task-management/internal/handlers"
	"github.com/raifuki/task-management/internal/middlewares"
	"github.com/raifuki/task-management/internal/ratelimit"
	"github.com/raifuki/task-management/internal/repositories"
	"github.com/raifuki/task-management/internal/services"
	"github.com/raifuki/task-management/pkg/jwt"
)

func SetupRouter(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.Logger())
	r.Use(middlewares.CORS())
	r.Use(middlewares.ErrorHandler())

	jwtMgr := jwt.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiredHours)
	cacheLayer := cache.NewRedisCache(rdb)
	redisLimiter := ratelimit.NewRedisLimiter(rdb)
	memoryLimiter := ratelimit.NewMemoryLimiter()
	limiter := ratelimit.NewFallbackLimiter(redisLimiter, memoryLimiter)

	globalCfg := ratelimit.Config{
		Capacity:   cfg.RateLimitGlobalCapacity,
		RefillRate: cfg.RateLimitGlobalRefill,
	}
	authCfg := ratelimit.Config{
		Capacity:   cfg.RateLimitAuthCapacity,
		RefillRate: cfg.RateLimitAuthRefill,
	}
	writeCfg := ratelimit.Config{
		Capacity:   cfg.RateLimitWriteCapacity,
		RefillRate: cfg.RateLimitWriteRefill,
	}

	userRepo := repositories.NewUserRepository(db)
	projRepo := repositories.NewProjectRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	cmtRepo := repositories.NewCommentRepository(db)

	authSvc := services.NewAuthService(userRepo, jwtMgr)
	projSvc := services.NewProjectService(projRepo, cacheLayer)
	taskSvc := services.NewTaskService(taskRepo, projRepo, cacheLayer)
	cmtSvc := services.NewCommentService(cmtRepo, taskRepo)

	authH := handlers.NewAuthHandler(authSvc)
	projH := handlers.NewProjectHandler(projSvc)
	taskH := handlers.NewTaskHandler(taskSvc)
	cmtH := handlers.NewCommentHandler(cmtSvc)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	api.Use(middlewares.RateLimitByIP(limiter, globalCfg, "global"))

	auth := api.Group("/auth")
	auth.Use(middlewares.RateLimitByIP(limiter, authCfg, "auth"))
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	protected := api.Group("")
	protected.Use(middlewares.Auth(jwtMgr))
	protected.Use(middlewares.RateLimitByUser(limiter, ratelimit.ReadLimit, "user"))
	{
		projects := protected.Group("/projects")
		{
			projects.POST("", middlewares.RateLimitByUser(limiter, writeCfg, "write"), projH.Create)
			projects.GET("", projH.List)
			projects.GET("/:id", projH.Get)
			projects.PUT("/:id", middlewares.RateLimitByUser(limiter, writeCfg, "write"), projH.Update)
			projects.DELETE("/:id", middlewares.RateLimitByUser(limiter, writeCfg, "write"), projH.Delete)
		}

		tasks := protected.Group("/tasks")
		{
			tasks.POST("/project/:projectId", middlewares.RateLimitByUser(limiter, writeCfg, "write"), taskH.Create)
			tasks.GET("", taskH.List)
			tasks.GET("/:id", taskH.Get)
			tasks.PUT("/:id", middlewares.RateLimitByUser(limiter, writeCfg, "write"), taskH.Update)
			tasks.DELETE("/:id", middlewares.RateLimitByUser(limiter, writeCfg, "write"), taskH.Delete)
		}

		comments := protected.Group("/comments")
		{
			comments.POST("/task/:taskId", middlewares.RateLimitByUser(limiter, writeCfg, "write"), cmtH.Create)
			comments.GET("/task/:taskId", cmtH.List)
			comments.GET("/:id", cmtH.Get)
			comments.PUT("/:id", middlewares.RateLimitByUser(limiter, writeCfg, "write"), cmtH.Update)
			comments.DELETE("/:id", middlewares.RateLimitByUser(limiter, writeCfg, "write"), cmtH.Delete)
		}
	}

	return r
}
