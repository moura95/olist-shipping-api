package server

import (
	"net/http"

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

//      @title                  Olist Api
//      @version                1.0
//      @description    Api Shipping
//      @termsOfService http://swagger.io/terms/

//      @license.name   Apache 2.0
//      @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host           localhost:8080
// @BasePath       /api/v1
func NewServer(cfg config.Config, store repository.Querier, log *zap.SugaredLogger) *Server {

	server := &Server{
		store:  &store,
		config: &cfg,
		logger: log,
	}
	var router *gin.Engine

	router = gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(middleware.RequestLogMiddleware(log))
	router.Use(middleware.ResponseLogMiddleware(log))

	// Middlewares existentes
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Init all Routers
	createRoutesV1(&store, server.config, router, log)

	server.router = router
	return server
}

func createRoutesV1(store *repository.Querier, cfg *config.Config, router *gin.Engine, log *zap.SugaredLogger) {
	// Criar services
	packageService := service.NewPackageService(*store, *cfg, log)

	// Criar handlers
	packageHandler := handler.NewPackageHandler(packageService, cfg, log)

	// Configurar rotas
	apiV1 := router.Group("/api/v1")
	{
		// Rotas de packages
		packages := apiV1.Group("/packages")
		{
			packages.GET("", packageHandler.List)
			packages.GET("/:id", packageHandler.GetByID)
			packages.POST("", packageHandler.Create)
			packages.PATCH("/:id/status", packageHandler.UpdateStatus)
			packages.POST("/:id/hire", packageHandler.HireCarrier)
			packages.DELETE("/:id", packageHandler.Delete)
		}

		// Rota de tracking
		apiV1.GET("/track/:tracking_code", packageHandler.GetByTrackingCode)

		// Rota de cotações
		apiV1.GET("/quotes", packageHandler.GetQuotes)

		// Rotas auxiliares
		carriers := apiV1.Group("/carriers")
		{
			carriers.GET("", packageHandler.ListCarriers)
		}

		states := apiV1.Group("/states")
		{
			states.GET("", packageHandler.ListStates)
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
