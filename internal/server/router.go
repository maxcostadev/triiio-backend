package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/vahiiiid/go-rest-api-boilerplate/internal/auth"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/config"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/health"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/middleware"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(h *Handlers, authService auth.Service, cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.New()

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	skipPaths := config.GetSkipPaths(cfg.App.Environment)
	loggerConfig := middleware.NewLoggerConfig(
		cfg.Logging.GetLogLevel(),
		skipPaths,
	)
	router.Use(middleware.Logger(loggerConfig))
	router.Use(errors.ErrorHandler())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))

	var checkers []health.Checker
	if cfg.Health.DatabaseCheckEnabled {
		dbChecker := health.NewDatabaseChecker(db)
		checkers = append(checkers, dbChecker)
	}
	healthService := health.NewService(checkers, cfg.App.Version, cfg.App.Environment)
	healthHandler := health.NewHandler(healthService)

	router.GET("/health", healthHandler.Health)
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rlCfg := cfg.Ratelimit
	if rlCfg.Enabled {
		router.Use(
			middleware.NewRateLimitMiddleware(
				rlCfg.Window,
				rlCfg.Requests,
				func(c *gin.Context) string {
					ip := c.ClientIP()
					if ip == "" {
						ip = c.GetHeader("X-Forwarded-For")
						if ip == "" {
							ip = c.GetHeader("X-Real-IP")
						}
						if ip == "" {
							ip = "unknown"
						}
					}
					return ip
				},
				nil,
			),
		)
	}

	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", h.User.Register)
			authGroup.POST("/login", h.User.Login)
			authGroup.POST("/refresh", h.User.RefreshToken)
			authGroup.POST("/logout", auth.AuthMiddleware(authService), h.User.Logout)
			authGroup.GET("/me", auth.AuthMiddleware(authService), h.User.GetMe)
		}

		// User endpoints - authenticated users can access their own resources
		usersGroup := v1.Group("/users")
		usersGroup.Use(auth.AuthMiddleware(authService))
		{
			usersGroup.GET("/:id", h.User.GetUser)
			usersGroup.PUT("/:id", h.User.UpdateUser)
			usersGroup.DELETE("/:id", h.User.DeleteUser)
		}

		// Admin endpoints - admin role required, following REST best practices
		adminGroup := v1.Group("/admin")
		adminGroup.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			// User management endpoints
			adminGroup.GET("/users", h.User.ListUsers)
			adminGroup.GET("/users/:id", h.User.GetUser)
			adminGroup.PUT("/users/:id", h.User.UpdateUser)
			adminGroup.DELETE("/users/:id", h.User.DeleteUser)
		}

		public := v1.Group("/sliders")
		{
			public.GET("", h.Sliders.ListSliders)
			public.GET("/location", h.Sliders.GetSliderByLocation)
			public.GET("/items/:item_id", h.Sliders.GetSliderItem)
			public.GET(":id", h.Sliders.GetSlider)
			public.GET("/:id/items", h.Sliders.GetSliderItems)
		}

		// Protected routes
		protected := v1.Group("/sliders")
		protected.Use(auth.AuthMiddleware(authService))
		{
			protected.POST("", h.Sliders.CreateSlider)
			protected.POST("/:id/items", h.Sliders.AddSliderItem)
			protected.PUT("/items/:item_id", h.Sliders.UpdateSliderItem)
			protected.DELETE("/items/:item_id", h.Sliders.DeleteSliderItem)

			protected.PUT("/:id", h.Sliders.UpdateSlider)
			protected.DELETE("/:id", h.Sliders.DeleteSlider)
		}

		// Imoveis endpoints
		imoveisPublic := v1.Group("/imoveis")
		{
			imoveisPublic.GET("", h.Imoveis.ListImoveis)
			imoveisPublic.GET("/:id", h.Imoveis.GetImovel)
			imoveisPublic.GET("/:id/anexos", h.Imoveis.GetAnexos)
			imoveisPublic.GET("/:id/caracteristicas", h.Imoveis.GetCaracteristicas)
		}

		imoveisProtected := v1.Group("/imoveis")
		imoveisProtected.Use(auth.AuthMiddleware(authService))
		{
			imoveisProtected.POST("", h.Imoveis.CreateImovel)
			imoveisProtected.POST("/import", h.Imoveis.ImportProperties)
			imoveisProtected.PUT("/:id", h.Imoveis.UpdateImovel)
			imoveisProtected.DELETE("/:id", h.Imoveis.DeleteImovel)
			imoveisProtected.POST("/:id/anexos", h.Imoveis.AddAnexo)
			imoveisProtected.POST("/:id/caracteristicas", h.Imoveis.AddCaracteristicas)
		}

		// Email endpoints - protected
		emailGroup := v1.Group("/emails")
		emailGroup.Use(auth.AuthMiddleware(authService))
		{
			emailGroup.POST("/send", h.Email.SendEmail)
			emailGroup.POST("/send-template", h.Email.SendTemplateEmail)
		}
	}

	return router
}
