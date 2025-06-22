package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/docs"
	"github/moura95/olist-shipping-api/internal/handler"
	"github/moura95/olist-shipping-api/internal/middleware"
	"github/moura95/olist-shipping-api/internal/repository"
	"github/moura95/olist-shipping-api/internal/service"

	"go.uber.org/zap"
)

type Server struct {
	store  *repository.Querier
	router *gin.Engine
	config *config.Config
	logger *zap.SugaredLogger
}

// @title           Olist Shipping API
// @version         1.0
// @description     API para gerenciamento de entregas da Olist
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@olist.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func NewServer(cfg config.Config, store repository.Querier, log *zap.SugaredLogger) *Server {

	server := &Server{
		store:  &store,
		config: &cfg,
		logger: log,
	}
	var router *gin.Engine

	router = gin.Default()

	docs.SwaggerInfo.BasePath = ""

	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.RequestLogMiddleware(log))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("redirect")
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowHeaders("Content-Type")
	router.Use(cors.New(corsConfig))

	createRoutesV1(&store, server.config, router, log)

	server.router = router
	return server
}

func createRoutesV1(store *repository.Querier, cfg *config.Config, router *gin.Engine, log *zap.SugaredLogger) {
	packageService := service.NewPackageService(*store, *cfg, log)

	packageHandler := handler.NewPackageHandler(packageService, cfg, log)
	quoteHandler := handler.NewQuoteHandler(packageService, cfg, log)
	carrierHandler := handler.NewCarrierHandler(packageService, cfg, log)
	stateHandler := handler.NewStateHandler(packageService, cfg, log)

	apiV1 := router.Group("/api/v1")
	{
		packages := apiV1.Group("/packages")
		{
			packages.GET("", packageHandler.List)
			packages.GET("/:id", packageHandler.GetByID)
			packages.POST("", packageHandler.Create)
			packages.PATCH("/:id/status", packageHandler.UpdateStatus)
			packages.POST("/:id/hire", packageHandler.HireCarrier)
			packages.DELETE("/:id", packageHandler.Delete)
		}

		apiV1.GET("/packages/tracking/:tracking_code", packageHandler.GetByTrackingCode)

		apiV1.GET("/quotes", quoteHandler.GetQuotes)

		carriers := apiV1.Group("/carriers")
		{
			carriers.GET("", carrierHandler.List)
		}

		states := apiV1.Group("/states")
		{
			states.GET("", stateHandler.List)
		}
	}
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func RunGinServer(cfg config.Config, store repository.Querier, log *zap.SugaredLogger) {
	server := NewServer(cfg, store, log)

	_ = server.Start(cfg.HTTPServerAddress)
}
